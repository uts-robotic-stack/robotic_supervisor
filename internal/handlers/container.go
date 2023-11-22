package handlers

import (
	"net/http"
	"sync"

	"github.com/containrrr/watchtower/internal/actions"
	"github.com/containrrr/watchtower/pkg/container"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ContainerHandler struct {
	Client        container.Client
	LogsFrequency float64
	sync.Mutex
}

func NewContainerHandler(client container.Client) *ContainerHandler {
	return &ContainerHandler{
		Client: client,
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
	defer conn.Close()

	containerName := c.Query("container")
	// Iterate through each container and retrieve its logs
	go func() {
		actions.BroadcastLogs(conn, h.Client, containerName, h.LogsFrequency)
	}()

	select {}
}
