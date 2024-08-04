package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/dkhoanguyen/watchtower/internal/actions"
	containerService "github.com/dkhoanguyen/watchtower/pkg/container"
	"github.com/dkhoanguyen/watchtower/pkg/service"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
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
func (h *ContainerHandler) HandleContainerCreate(c *gin.Context) {
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
			Image: serviceReq.Image,
			Tty:   serviceReq.Tty,
			Env:   service.FormatEnvVars(serviceReq.EnvVars),
			Cmd:   serviceReq.Command,
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

		id, err := h.client.CreateContainer(
			serviceName, *config, *hostConfig, *networkingConfig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create container: %v", err)})
			return
		}
		resp.ServiceID[serviceName] = id.ShortID()
	}

	c.JSON(http.StatusOK, resp)
}

// Handle start
func (h *ContainerHandler) HandleContainerStart(c *gin.Context) {
	log.Info("Received HTTP request to start container")
	var srvIDMap service.ServiceIDMap
	if err := c.ShouldBindJSON(&srvIDMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	fmt.Println(srvIDMap)
	resp := make(map[string]bool)
	failedService := ""
	for serviceName, serviceID := range srvIDMap.ServiceID {
		fmt.Println(serviceID)
		err := h.client.StartContainerByID(serviceID)
		if err != nil {
			log.Error(err)
			resp[serviceName] = false
			failedService += serviceName + " "
			continue
		}
		resp[serviceName] = true
	}
	if failedService != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to start " + failedService})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Handle stop

// Handle remove (equivalen to unload)

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
