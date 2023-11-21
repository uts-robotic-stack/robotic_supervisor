package handlers

import (
	"time"

	"github.com/containrrr/watchtower/internal/actions"
	"github.com/containrrr/watchtower/pkg/container"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type ContainerHandler struct {
	Client *container.Client
	Lock   chan bool
}

var upgrader = websocket.Upgrader{}

func (h *ContainerHandler) HandleStreamLogs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	containerName := c.Query("container")
	log.Info("Requested container: " + containerName)
	if err := actions.StreamLogs(*h.Client, containerName, true, 60*time.Second, conn); err != nil {
		log.Error(err)
	}
	// select {
	// case chanValue := <-h.Lock:
	// 	defer func() {
	// 		h.Lock <- chanValue
	// 	}()

	// default:
	// 	log.Info("Skipped. Another docker process is already running.")
	// 	// c.JSON(http.StatusConflict, "Request dropped. Another docker process is already running.")
	// }
}
