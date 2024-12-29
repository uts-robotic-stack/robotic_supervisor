package cmd

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dkhoanguyen/watchtower/internal/actions"
	"github.com/dkhoanguyen/watchtower/internal/api"
	"github.com/dkhoanguyen/watchtower/internal/flags"
	"github.com/dkhoanguyen/watchtower/internal/handlers"
	"github.com/dkhoanguyen/watchtower/internal/meta"
	"github.com/dkhoanguyen/watchtower/internal/middleware"
	"github.com/dkhoanguyen/watchtower/pkg/container"
	"github.com/dkhoanguyen/watchtower/pkg/filters"
	"github.com/dkhoanguyen/watchtower/pkg/metrics"
	"github.com/dkhoanguyen/watchtower/pkg/notifications"
	"github.com/dkhoanguyen/watchtower/pkg/service"
	t "github.com/dkhoanguyen/watchtower/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

var (
	client            container.Client
	scheduleSpec      string
	cleanup           bool
	noRestart         bool
	noPull            bool
	monitorOnly       bool
	enableLabel       bool
	disableContainers []string
	notifier          t.Notifier
	timeout           time.Duration
	lifecycleHooks    bool
	rollingRestart    bool
	scope             string
	labelPrecedence   bool
	redisAddr         string
)

var rootCmd = NewRootCommand()

// NewRootCommand creates the root command for watchtower
func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "watchtower",
		Short: "Automatically updates running Docker containers",
		Long: `
	Watchtower automatically updates running Docker containers whenever a new image is released.
	More information available at https://github.com/dkhoanguyen/watchtower/.
	`,
		Run:    Run,
		PreRun: PreRun,
		Args:   cobra.ArbitraryArgs,
	}
}

func init() {
	flags.SetDefaults()
	flags.RegisterDockerFlags(rootCmd)
	flags.RegisterSystemFlags(rootCmd)
	flags.RegisterNotificationFlags(rootCmd)
}

// Execute the root func and exit in case of errors
func Execute() {
	rootCmd.AddCommand(notifyUpgradeCommand)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// PreRun is a lifecycle hook that runs before the command is executed.
func PreRun(cmd *cobra.Command, _ []string) {
	f := cmd.PersistentFlags()
	flags.ProcessFlagAliases(f)
	if err := flags.SetupLogging(f); err != nil {
		log.Fatalf("Failed to initialize logging: %s", err.Error())
	}

	scheduleSpec, _ = f.GetString("schedule")

	flags.GetSecretsFromFiles(cmd)
	cleanup, noRestart, monitorOnly, timeout = flags.ReadFlags(cmd)

	if timeout < 0 {
		log.Fatal("Please specify a positive value for timeout value.")
	}

	enableLabel, _ = f.GetBool("label-enable")
	disableContainers, _ = f.GetStringSlice("disable-containers")
	lifecycleHooks, _ = f.GetBool("enable-lifecycle-hooks")
	rollingRestart, _ = f.GetBool("rolling-restart")
	scope, _ = f.GetString("scope")
	labelPrecedence, _ = f.GetBool("label-take-precedence")

	if scope != "" {
		log.Debugf(`Using scope %q`, scope)
	}

	// configure environment vars for client
	err := flags.EnvConfig(cmd)
	if err != nil {
		log.Fatal(err)
	}

	noPull, _ = f.GetBool("no-pull")
	includeStopped, _ := f.GetBool("include-stopped")
	includeRestarting, _ := f.GetBool("include-restarting")
	reviveStopped, _ := f.GetBool("revive-stopped")
	removeVolumes, _ := f.GetBool("remove-volumes")
	warnOnHeadPullFailed, _ := f.GetString("warn-on-head-failure")

	if monitorOnly && noPull {
		log.Warn("Using `WATCHTOWER_NO_PULL` and `WATCHTOWER_MONITOR_ONLY` simultaneously might lead to no action being taken at all. If this is intentional, you may safely ignore this message.")
	}

	client = container.NewClient(container.ClientOptions{
		IncludeStopped:    includeStopped,
		ReviveStopped:     reviveStopped,
		RemoveVolumes:     removeVolumes,
		IncludeRestarting: includeRestarting,
		WarnOnHeadFailed:  container.WarningStrategy(warnOnHeadPullFailed),
	})

	notifier = notifications.NewNotifier(cmd)
	notifier.AddLogHook()

	// Populate redis if data does not exist
}

// Run is the main execution flow of the command
func Run(c *cobra.Command, names []string) {
	filter, filterDesc := filters.BuildFilter(names, disableContainers, enableLabel, scope)
	runOnce, _ := c.PersistentFlags().GetBool("run-once")
	// enableUpdateAPI, _ := c.PersistentFlags().GetBool("http-api-update")
	// enableMetricsAPI, _ := c.PersistentFlags().GetBool("http-api-metrics")
	// unblockHTTPAPI, _ := c.PersistentFlags().GetBool("http-api-periodic-polls")
	apiToken, _ := c.PersistentFlags().GetString("http-api-token")
	healthCheck, _ := c.PersistentFlags().GetBool("health-check")
	port, _ := c.PersistentFlags().GetString("port")
	updateOnStartup, _ := c.PersistentFlags().GetBool("update-on-startup")
	redisAddr, _ = c.PersistentFlags().GetString("redis-addr")

	if healthCheck {
		// health check should not have pid 1
		if os.Getpid() == 1 {
			time.Sleep(1 * time.Second)
			log.Fatal("The health check flag should never be passed to the main watchtower container process")
		}
		os.Exit(0)
	}

	if rollingRestart && monitorOnly {
		log.Fatal("Rolling restarts is not compatible with the global monitor only flag")
	}

	awaitDockerClient()

	if err := actions.CheckForSanity(client, filter, rollingRestart); err != nil {
		logNotifyExit(err)
	}

	if runOnce {
		writeStartupMessage(c, time.Time{}, filterDesc)
		runUpdatesWithNotifications(filter)
		notifier.Close()
		os.Exit(0)
		return
	}

	if err := actions.CheckForMultipleWatchtowerInstances(client, cleanup, scope); err != nil {
		logNotifyExit(err)
	}

	// The lock is shared between the scheduler and the HTTP API. It only allows one update to run at a time.
	clientLock := make(chan bool, 1)
	clientLock <- true

	// Create a new Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	// Use default recovery
	router.Use(gin.Recovery())
	// Add CORS middleware
	router.Use(middleware.CORSMiddleware())
	// Add authentication
	router.Use(middleware.AuthMiddleware(apiToken))
	// Add logging
	router.Use(middleware.Logger())

	// Create Redis handler
	redisHandler := handlers.NewRedisHandler("10.211.55.7:6379", "", 0)

	// Check for data
	ctx := context.Background()
	data, err := redisHandler.Get(ctx, "default_services")
	if err == nil || data == "" {
		log.Info("Populating redis with default services")
		// Load default services to redis db
		yamlData, err := os.ReadFile("/config/default_services.yaml")
		if err != nil {
			log.Error(err)
		}
		var defaultServices service.ServiceMap
		if err = yaml.Unmarshal(yamlData, &defaultServices); err != nil {
			log.Error(err)
		}
		var jsonData []byte
		if jsonData, err = json.Marshal(defaultServices); err != nil {
			log.Error(err)
		}
		if err = redisHandler.Set(ctx, "default_services", jsonData); err != nil {
			log.Error(err)
		}
	}

	data, err = redisHandler.Get(ctx, "excluded_services")
	if err == nil || data == "" {
		log.Info("Populating redis with excluded services")
		yamlData, err := os.ReadFile("/config/excluded_services.yaml")
		if err != nil {
			log.Error(err)
		}
		var excludedServices map[string][]string
		err = yaml.Unmarshal(yamlData, &excludedServices)
		if err != nil {
			log.Error(err)
		}
		var jsonData []byte
		if jsonData, err = json.Marshal(excludedServices["services"]); err != nil {
			log.Error(err)
		}
		if err = redisHandler.Set(ctx, "excluded_services", jsonData); err != nil {
			log.Error(err)
		}
	}

	// Create handlers
	watchtowerHandler := handlers.WatchtowerHandler{
		Client:            &client,
		Filter:            filter,
		Notifier:          notifier,
		ScheduleSpec:      scheduleSpec,
		Cleanup:           cleanup,
		NoRestart:         noRestart,
		NoPull:            noPull,
		MonitorOnly:       monitorOnly,
		EnableLabel:       enableLabel,
		DisableContainers: disableContainers,
		Timeout:           timeout,
		LifecycleHooks:    lifecycleHooks,
		RollingRestart:    rollingRestart,
		Scope:             scope,
		LabelPrecedence:   labelPrecedence,
		Lock:              clientLock,
	}

	deviceHandler := handlers.DeviceHandler{
		Client:                  client,
		HardwareStatusFrequency: 0.1, // Once every 10 seconds
	}

	containerHandler := handlers.NewContainerHandler(client, 1, redisHandler)
	userHandler := handlers.NewUserHandler()

	// Set routes
	api.SetRoutes(router, &deviceHandler, &watchtowerHandler, containerHandler, userHandler)

	log.Infof("Serving api at port %v", port)
	// Start api
	go func() {
		router.Run(":" + port)
	}()

	// Run update once startup to check and download updates from the cloud
	if updateOnStartup {
		runCheckForUpdates(filter)
		metric := runUpdatesWithNotifications(filter)
		metrics.RegisterScan(metric)
	}

	if err := runChecksOnSchedule(c, filter, filterDesc, clientLock); err != nil {
		log.Error(err)
	}

	os.Exit(1)
}

func logNotifyExit(err error) {
	log.Error(err)
	notifier.Close()
	os.Exit(1)
}

func awaitDockerClient() {
	log.Debug("Sleeping for a second to ensure the docker api client has been properly initialized.")
	time.Sleep(1 * time.Second)
}

func formatDuration(d time.Duration) string {
	sb := strings.Builder{}

	hours := int64(d.Hours())
	minutes := int64(math.Mod(d.Minutes(), 60))
	seconds := int64(math.Mod(d.Seconds(), 60))

	if hours == 1 {
		sb.WriteString("1 hour")
	} else if hours != 0 {
		sb.WriteString(strconv.FormatInt(hours, 10))
		sb.WriteString(" hours")
	}

	if hours != 0 && (seconds != 0 || minutes != 0) {
		sb.WriteString(", ")
	}

	if minutes == 1 {
		sb.WriteString("1 minute")
	} else if minutes != 0 {
		sb.WriteString(strconv.FormatInt(minutes, 10))
		sb.WriteString(" minutes")
	}

	if minutes != 0 && (seconds != 0) {
		sb.WriteString(", ")
	}

	if seconds == 1 {
		sb.WriteString("1 second")
	} else if seconds != 0 || (hours == 0 && minutes == 0) {
		sb.WriteString(strconv.FormatInt(seconds, 10))
		sb.WriteString(" seconds")
	}

	return sb.String()
}

func writeStartupMessage(c *cobra.Command, sched time.Time, filtering string) {
	noStartupMessage, _ := c.PersistentFlags().GetBool("no-startup-message")
	enableUpdateAPI, _ := c.PersistentFlags().GetBool("http-api-update")

	var startupLog *log.Entry
	if noStartupMessage {
		startupLog = notifications.LocalLog
	} else {
		startupLog = log.NewEntry(log.StandardLogger())
		// Batch up startup messages to send them as a single notification
		notifier.StartNotification()
	}

	startupLog.Info("Watchtower ", meta.Version)

	notifierNames := notifier.GetNames()
	if len(notifierNames) > 0 {
		startupLog.Info("Using notifications: " + strings.Join(notifierNames, ", "))
	} else {
		startupLog.Info("Using no notifications")
	}

	startupLog.Info(filtering)

	if !sched.IsZero() {
		until := formatDuration(time.Until(sched))
		startupLog.Info("Scheduling first run: " + sched.Format("2006-01-02 15:04:05 -0700 MST"))
		startupLog.Info("Note that the first check will be performed in " + until)
	} else if runOnce, _ := c.PersistentFlags().GetBool("run-once"); runOnce {
		startupLog.Info("Running a one time update.")
	} else {
		startupLog.Info("Periodic runs are not enabled.")
	}

	if enableUpdateAPI {
		// TODO: make listen port configurable
		startupLog.Info("The HTTP API is enabled at :8080.")
	}

	if !noStartupMessage {
		// Send the queued up startup messages, not including the trace warning below (to make sure it's noticed)
		notifier.SendNotification(nil)
	}

	if log.IsLevelEnabled(log.TraceLevel) {
		startupLog.Warn("Trace level enabled: log will include sensitive information as credentials and tokens")
	}
}

func runChecksOnSchedule(c *cobra.Command, filter t.Filter, filtering string, lock chan bool) error {
	if lock == nil {
		lock = make(chan bool, 1)
		lock <- true
	}

	scheduler := cron.New()
	err := scheduler.AddFunc(
		scheduleSpec,
		func() {
			v := <-lock
			defer func() { lock <- v }()
			// Check for updates from registry and from local devices
			runCheckForUpdates(filter)
		})

	if err != nil {
		return err
	}

	writeStartupMessage(c, scheduler.Entries()[0].Schedule.Next(time.Now()), filtering)

	scheduler.Start()

	// Graceful shut-down on SIGINT/SIGTERM
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	signal.Notify(interrupt, syscall.SIGTERM)

	<-interrupt
	scheduler.Stop()
	log.Info("Waiting for running update to be finished...")
	<-lock
	return nil
}

func runCheckForUpdates(filter t.Filter) {
	updateParams := t.UpdateParams{
		Filter:          filter,
		Cleanup:         cleanup,
		NoRestart:       noRestart,
		Timeout:         timeout,
		MonitorOnly:     monitorOnly,
		LifecycleHooks:  lifecycleHooks,
		RollingRestart:  rollingRestart,
		LabelPrecedence: labelPrecedence,
	}
	// Check for updates from registry first
	if updateAvailable, err := actions.CheckForNewUpdateFromRegistry(client, updateParams); err != nil {
		log.Error(err)
	} else if updateAvailable {
		log.Info("Updates available from registry. Attempting to pull updates now...")
		err := actions.DownloadUpdate(client, updateParams)
		if err != nil {
			log.Error(err)
		}
	} else if !updateAvailable {
		log.Debug("Updates not available from upstream")
	} else {
		log.Debug("Unable to check for update from upstream registry")
	}
}

func runUpdatesWithNotifications(filter t.Filter) *metrics.Metric {
	notifier.StartNotification()
	updateParams := t.UpdateParams{
		Filter:          filter,
		Cleanup:         cleanup,
		NoRestart:       noRestart,
		Timeout:         timeout,
		MonitorOnly:     monitorOnly,
		LifecycleHooks:  lifecycleHooks,
		RollingRestart:  rollingRestart,
		LabelPrecedence: labelPrecedence,
		NoPull:          noPull,
	}
	// Run and check for updated on the cloud. Do not attempt to load local image
	log.Info("Update requested. Updating...")
	result, err := actions.Update(client, updateParams)
	if err != nil {
		log.Error(err)
	}
	notifier.SendNotification(result)
	metricResults := metrics.NewMetric(result)
	notifications.LocalLog.WithFields(log.Fields{
		"Scanned": metricResults.Scanned,
		"Updated": metricResults.Updated,
		"Failed":  metricResults.Failed,
	}).Info("Session done")
	return metricResults
}
