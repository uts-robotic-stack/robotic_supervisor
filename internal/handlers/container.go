package handlers

import (
	"net/http"
	"sync"

	"github.com/containrrr/watchtower/internal/actions"
	"github.com/containrrr/watchtower/pkg/container"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type ContainerHandler struct {
	client        container.Client
	logsFrequency float64
	wsClients     ClientList
	sync.Mutex
}

func NewContainerHandler(client container.Client, logFreq float64) *ContainerHandler {
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

func (h *ContainerHandler) HandleContainerStart(c *gin.Context) {
	log.Info("Received HTTP request to start container")
	// This should not be service body
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.Info("Running service handle")
	services := []container.Service{}
	rawServiceData := body["services"].(map[string]interface{})

	// Extract services from the raw data
	for name, config := range rawServiceData {
		service := container.MakeService(config.(map[string]interface{}), name)
		services = append(services, service)
	}

	// Execute actions on the services
	for _, service := range services {
		if err := actions.StartContainer(h.client, &service); err != nil {
			log.Error(err)
		}
	}
	c.JSON(http.StatusOK, nil)
}

func (h *ContainerHandler) HandleContainerStop(c *gin.Context) {
	log.Info("Received HTTP request to stop container")
	// This should not be service body
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.Info("Running service handle")
	services := []container.Service{}
	rawServiceData := body["services"].(map[string]interface{})

	// Extract services from the raw data
	for name, config := range rawServiceData {
		service := container.MakeService(config.(map[string]interface{}), name)
		services = append(services, service)
	}

	// Execute actions on the services
	for _, service := range services {
		if err := actions.StartContainer(h.client, &service); err != nil {
			log.Error(err)
		}
	}
	c.JSON(http.StatusOK, nil)
}

func (h *ContainerHandler) HandleContainerInspect(c *gin.Context) {
	log.Info("Received HTTP request to inspect container")
	containerName := c.Query("container")
	output, err := actions.InspectContainer(h.client, containerName)
	if err != nil {
		log.Error(err)
	}
	c.JSON(http.StatusOK, output)
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
