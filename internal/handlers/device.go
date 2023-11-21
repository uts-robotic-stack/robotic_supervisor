package handlers

import (
	"net/http"

	"github.com/containrrr/watchtower/internal/actions"
	"github.com/containrrr/watchtower/pkg/container"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type DeviceHandler struct {
	Client *container.Client
	Lock   chan bool
}

func (d *DeviceHandler) HandleGetDeviceInfo(c *gin.Context) {
	select {
	case chanValue := <-d.Lock:
		defer func() {
			d.Lock <- chanValue
		}()
		log.Info("Received HTTP request to get device-info")
		output := actions.GetDeviceInfo(*d.Client)
		c.JSON(http.StatusOK, output)

	default:
		log.Info("Skipped. Another docker process is already running.")
		c.JSON(http.StatusConflict, "Request dropped. Another docker process is already running.")
	}
}
