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
			deviceSubgroup.GET("/hardware-status", deviceHandler.HandlerWSHardwareStatus)
		}

		watchtowerSubgroup := v1.Group("/robotics_supervisor")
		{
			watchtowerSubgroup.POST("/update", watchtowerHandler.HandlePostUpdate)
			watchtowerSubgroup.POST("/download", watchtowerHandler.HandlePostDownload)
			watchtowerSubgroup.GET("/log-stream", containerHandler.HandleWSLogs)
			watchtowerSubgroup.GET("/list", containerHandler.HandleContainerStart)
			watchtowerSubgroup.POST("/start", containerHandler.HandleContainerStart)
			watchtowerSubgroup.POST("/stop", containerHandler.HandleContainerStop)
			watchtowerSubgroup.GET("/inspect", containerHandler.HandleContainerInspect)
		}
	}
}
