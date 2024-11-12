package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dkhoanguyen/watchtower/internal/actions"
	containerService "github.com/dkhoanguyen/watchtower/pkg/container"
	"github.com/dkhoanguyen/watchtower/pkg/filters"
	"github.com/dkhoanguyen/watchtower/pkg/service"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ContainerHandler struct {
	client        containerService.Client
	logsFrequency float64
	wsClients     ClientList
	sync.Mutex
}

func NewContainerHandler(client containerService.Client, logFreq float64) *ContainerHandler {
	return &ContainerHandler{
		client:        client,
		logsFrequency: logFreq,
		wsClients:     make(ClientList),
	}
}

// addClient will add clients to our clientList
func (h *ContainerHandler) addClient(client *Client) {
	// Lock so we can manipulate
	h.Lock()
	defer h.Unlock()
	// Add Client
	h.wsClients[client] = true
}

// removeClient will remove the client and clean up
func (h *ContainerHandler) removeClient(client *Client) {
	h.Lock()
	defer h.Unlock()
	// Check if Client exists, then delete it
	if _, ok := h.wsClients[client]; ok {
		// close connection
		client.connection.Close()
		// remove
		delete(h.wsClients, client)
	}
}

// Handle create (equivalent to load)
func (h *ContainerHandler) HandleContainerStart(c *gin.Context) {
	log.Info("Received HTTP request to create container")

	var srvMap service.ServiceMap
	if err := c.ShouldBindJSON(&srvMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Response
	resp := service.ServiceIDMap{ServiceID: make(map[string]string)}
	for serviceName, serviceReq := range srvMap.Services {
		config := &container.Config{
			Image:  serviceReq.Image.Name,
			Tty:    serviceReq.Tty,
			Env:    service.FormatEnvVars(serviceReq.EnvVars),
			Cmd:    serviceReq.Command,
			Labels: serviceReq.Labels,
		}

		hostConfig := &container.HostConfig{
			Privileged:  serviceReq.Privileged,
			NetworkMode: container.NetworkMode(serviceReq.Network),
			Mounts:      service.FormatMounts(serviceReq.Mounts),
			Binds:       service.FormatVolumes(serviceReq.Volumes),
		}

		if serviceReq.Restart != "" {
			hostConfig.RestartPolicy = container.RestartPolicy{Name: serviceReq.Restart}
		}

		networkingConfig := &network.NetworkingConfig{}
		if serviceReq.NetworkSettings.IPAddress != "" {
			networkingConfig.EndpointsConfig = map[string]*network.EndpointSettings{
				serviceReq.Network: {
					IPAMConfig: &network.EndpointIPAMConfig{
						IPv4Address: serviceReq.NetworkSettings.IPAddress,
					},
					Gateway: serviceReq.NetworkSettings.Gateway,
				},
			}
		}

		id, err := h.client.StartContainer(
			serviceName, *config, *hostConfig, *networkingConfig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create container: %v", err)})
			return
		}
		resp.ServiceID[serviceName] = id.ShortID()
	}
	c.JSON(http.StatusOK, resp)
}

// Handle stop
func (h *ContainerHandler) HandleContainerStop(c *gin.Context) {
	log.Info("Received HTTP request to stop container")
	var srvIDMap service.ServiceIDMap
	if err := c.ShouldBindJSON(&srvIDMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	resp := make(map[string]bool)
	failedService := ""

	// List all containers
	containers, _ := h.client.ListContainers(filters.NoFilter)

	for serviceName := range srvIDMap.ServiceID {
		for _, cnt := range containers {
			if cnt.Name()[1:] == serviceName {
				err := h.client.StopContainer(cnt, time.Second)
				resp[serviceName] = true
				if err != nil {
					log.Error(err)
					resp[serviceName] = false
					failedService += serviceName + " "
				}
				break
			}
		}
	}

	if failedService != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to stop " + failedService})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Handle logs
func (h *ContainerHandler) HandleWSLogs(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := NewWSClient(conn, h)
	h.addClient(client)

	containerName := c.Query("container")
	go client.readMessages()
	go client.broadcastLogs(containerName)
}

func (h *ContainerHandler) HandlerContainerLogs(c *gin.Context) {
	log.Info("Received HTTP request to get container logs")
	containerName := c.Query("container")
	output, err := actions.GetLogs(h.client, containerName)
	if err != nil {
		log.Error(err)
	}
	c.JSON(http.StatusOK, output)
}

// Handle inspect
func (h *ContainerHandler) HandleContainerInspect(c *gin.Context) {
	log.Info("Received HTTP request to inspect container")
	c.JSON(http.StatusOK, nil)
}

func (h *ContainerHandler) HandleGetAllContainers(c *gin.Context) {
	log.Info("Received HTTP request to get all containers")
	containers, _ := h.client.ListContainers(filters.NoFilter)
	containerList := service.ServiceMap{
		Services: make(map[string]service.Service),
	}
	for _, cnt := range containers {
		containerDetails := service.Service{
			EnvVars: make(map[string]string),
			Sysctls: make(map[string]string),
			Labels:  make(map[string]string),
		}

		containerDetails.Name = strings.ReplaceAll(cnt.Name(), "/", "")
		containerDetails.Command = cnt.GetCreateConfig().Cmd
		containerDetails.ContainerID = cnt.ContainerInfo().ID
		containerDetails.Status = cnt.ContainerInfo().State.Status

		containerDetails.Image.Name = cnt.ContainerInfo().Config.Image
		containerDetails.Image.ID = cnt.ContainerInfo().Image
		containerDetails.Image.Created = cnt.ImageInfo().Created

		containerDetails.Labels = cnt.GetCreateConfig().Labels
		containerDetails.Privileged = cnt.GetCreateHostConfig().Privileged
		containerDetails.Resources.CPU = cnt.GetCreateHostConfig().Resources.CPUShares
		containerDetails.Resources.Memory = cnt.GetCreateHostConfig().Memory
		containerDetails.Sysctls = cnt.GetCreateHostConfig().Sysctls
		containerDetails.Restart = cnt.GetCreateHostConfig().RestartPolicy.Name

		// Get env vars and convert from list of string to map[string]string
		for _, envVar := range cnt.ContainerInfo().Config.Env {
			parts := strings.SplitN(envVar, "=", 2)
			if len(parts) == 2 {
				key := strings.ToUpper(parts[0])      // Capitalize the key
				value := parts[1]                     // Get the value
				containerDetails.EnvVars[key] = value // Add to map
			}
		}
		containerList.Services[strings.ReplaceAll(cnt.Name(), "/", "")] = containerDetails
	}
	c.JSON(http.StatusOK, containerList)
}

func (h *ContainerHandler) HandleGetDefaultServices(c *gin.Context) {
	log.Info("Received HTTP request to get default services")

	// Obtain default services
	// In the future this should be in a redis instance
	data, err := os.ReadFile("/config/default_services.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// Unmarshal YAML data into Go struct
	var services service.ServiceMap
	err = yaml.Unmarshal(data, &services)
	if err != nil {
		log.Fatalf("Unable to read settings.yaml to obtain default services: %v", err)
	}
	c.JSON(http.StatusOK, services)
}

func (h *ContainerHandler) HandleGetExcludedServices(c *gin.Context) {
	log.Info("Received HTTP request to get excluded services")

	// Obtain default services
	// In the future this should be in a redis instance
	data, err := os.ReadFile("/config/excluded_services.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// Unmarshal YAML data into Go struct
	var services map[string][]string
	err = yaml.Unmarshal(data, &services)
	if err != nil {
		log.Fatalf("Unable to read settings.yaml to obtain default services: %v", err)
	}
	c.JSON(http.StatusOK, services["services"])
}
