package handlers

import (
	"net/http"
	"sync"

	"github.com/containrrr/watchtower/internal/actions"
	"github.com/containrrr/watchtower/pkg/container"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ContainerHandler struct {
	Client      container.Client
	Connections map[*websocket.Conn]struct{}
	sync.Mutex
}

func (h *ContainerHandler) HandlerLogs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	containerName := c.Query("container")
	// Iterate through each container and retrieve its logs
	go actions.BroadcastLogs(conn, h.Client, containerName)

	select {}
}
