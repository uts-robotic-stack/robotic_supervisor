package api

import (
	"github.com/containrrr/watchtower/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetRoutes sets up the API routes.
func SetRoutes(router *gin.Engine,
	deviceHandler *handlers.DeviceHandler,
	watchtowerHandler *handlers.WatchtowerHandler,
	containerHandler *handlers.ContainerHandler) {

	v1 := router.Group("/api/v1")
	{
		deviceSubgroup := v1.Group("/device")
		{
			deviceSubgroup.GET("/info", deviceHandler.HandleGetDeviceInfo)
		}

		watchtowerSubgroup := v1.Group("/watchtower")
		{
			watchtowerSubgroup.POST("/update", watchtowerHandler.HandlePostUpdate)
			watchtowerSubgroup.POST("/download", watchtowerHandler.HandlePostDownload)
			watchtowerSubgroup.GET("/logs", containerHandler.HandlerLogs)
		}
	}
}
