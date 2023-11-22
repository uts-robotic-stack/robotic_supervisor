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
}

func (d *DeviceHandler) HandleGetDeviceInfo(c *gin.Context) {
	log.Info("Received HTTP request to get device-info")
	output := actions.GetDeviceInfo(*d.Client)
	c.JSON(http.StatusOK, output)
}
