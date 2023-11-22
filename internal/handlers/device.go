package handlers

import (
	"net/http"
	"time"

	"github.com/containrrr/watchtower/internal/actions"
	"github.com/containrrr/watchtower/pkg/container"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type DeviceHandler struct {
	Client                  container.Client
	HardwareStatusFrequency float64
}

func (d *DeviceHandler) HandleGetDeviceInfo(c *gin.Context) {
	log.Info("Received HTTP request to get device-info")
	output := actions.GetDeviceInfo(d.Client)
	c.JSON(http.StatusOK, output)
}

func (d *DeviceHandler) HandlerWSHardwareStatus(c *gin.Context) {
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
	go func() {
		actions.BroadcastHardwareStatus(conn, d.Client)
		time.Sleep(time.Duration(1/d.HardwareStatusFrequency) * time.Millisecond)
	}()
	select {}
}
