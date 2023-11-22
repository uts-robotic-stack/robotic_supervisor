package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/containrrr/watchtower/internal/actions"
	"github.com/containrrr/watchtower/pkg/container"
	"github.com/containrrr/watchtower/pkg/filters"
	"github.com/containrrr/watchtower/pkg/metrics"
	"github.com/containrrr/watchtower/pkg/notifications"
	"github.com/containrrr/watchtower/pkg/types"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type WatchtowerHandler struct {
	Client            *container.Client
	Filter            types.Filter
	Notifier          types.Notifier
	ScheduleSpec      string
	Cleanup           bool
	NoRestart         bool
	NoPull            bool
	MonitorOnly       bool
	EnableLabel       bool
	DisableContainers []string
	Timeout           time.Duration
	LifecycleHooks    bool
	RollingRestart    bool
	Scope             string
	LabelPrecedence   bool
	Lock              chan bool
}

func (w *WatchtowerHandler) HandlePostUpdate(c *gin.Context) {
	select {
	case chanValue := <-w.Lock:
		defer func() {
			w.Lock <- chanValue
		}()
		log.Info("Received HTTP request to apply updates")
		imagesParams := c.Query("images")
		images := strings.Split(imagesParams, ",")

		// By default do not apply any filter
		filter := filters.NoFilter

		// If POST has any image then apply those images only
		if len(images) > 0 {
			filter = filters.FilterByImage(images, w.Filter)
		}
		w.Notifier.StartNotification()
		updateParams := types.UpdateParams{
			Filter:          filter,
			Cleanup:         w.Cleanup,
			NoRestart:       w.NoRestart,
			Timeout:         w.Timeout,
			MonitorOnly:     w.MonitorOnly,
			LifecycleHooks:  w.LifecycleHooks,
			RollingRestart:  w.RollingRestart,
			LabelPrecedence: w.LabelPrecedence,
			NoPull:          w.NoPull,
		}
		// Run and check for updated on the cloud. Do not attempt to load local image
		log.Info("Update requested. Updating...")
		result, err := actions.Update(*w.Client, updateParams)
		if err != nil {
			log.Error(err)
		}
		w.Notifier.SendNotification(result)
		metricResults := metrics.NewMetric(result)
		notifications.LocalLog.WithFields(log.Fields{
			"Scanned": metricResults.Scanned,
			"Updated": metricResults.Updated,
			"Failed":  metricResults.Failed,
		}).Info("Session done")
		c.JSON(http.StatusOK, metricResults)

	default:
		log.Info("Skipped. Another update process is already running.")
		c.JSON(http.StatusConflict, "Request dropped. Another update process is already running.")
	}
}

func (w *WatchtowerHandler) HandlePostDownload(c *gin.Context) {
	select {
	case chanValue := <-w.Lock:
		defer func() {
			w.Lock <- chanValue
		}()
		log.Info("Received HTTP request to download updates")
		imagesParams := c.Query("images")
		log.Info("Requested images: " + imagesParams)

		images := strings.Split(imagesParams, ",")

		// By default do not apply any filter
		filter := filters.NoFilter

		// If POST has any image then apply those images only
		if len(images) > 0 {
			filter = filters.FilterByImage(images, w.Filter)
		}
		w.Notifier.StartNotification()
		updateParams := types.UpdateParams{
			Filter:          filter,
			Cleanup:         w.Cleanup,
			NoRestart:       w.NoRestart,
			Timeout:         w.Timeout,
			MonitorOnly:     w.MonitorOnly,
			LifecycleHooks:  w.LifecycleHooks,
			RollingRestart:  w.RollingRestart,
			LabelPrecedence: w.LabelPrecedence,
			NoPull:          w.NoPull,
		}
		// Download updates
		err := actions.DownloadUpdate(*w.Client, updateParams)
		if err != nil {
			log.Error(err)
		}
		c.JSON(http.StatusOK, nil)

	default:
		log.Info("Skipped. Another download process is already running.")
		c.JSON(http.StatusConflict, "Request dropped. Another download process is already running.")
	}
}
