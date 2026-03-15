package api

import (
	"github.com/dkhoanguyen/watchtower/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetRoutes sets up the API routes.
func SetRoutes(router *gin.Engine,
	deviceHandler *handlers.DeviceHandler,
	watchtowerHandler *handlers.WatchtowerHandler,
	containerHandler *handlers.ContainerHandler,
	userHandler *handlers.UserHandler) {

	v1 := router.Group("/api/v1")
	{
		deviceSubgroup := v1.Group("/device")
		{
			deviceSubgroup.GET("/info", deviceHandler.HandleGetDeviceInfo)
			deviceSubgroup.GET("/hardware-status", deviceHandler.HandlerWSHardwareStatus)
			deviceSubgroup.POST("/shutdown", deviceHandler.HandleShutdown)
			deviceSubgroup.POST("/restart", deviceHandler.HandleRestart)
		}

		watchtowerSubgroup := v1.Group("/supervisor")
		{
			watchtowerSubgroup.POST("/update", watchtowerHandler.HandlePostUpdate)
			watchtowerSubgroup.POST("/download", watchtowerHandler.HandlePostDownload)
			watchtowerSubgroup.GET("/auto-update", watchtowerHandler.HandleAutoUpdateStatus)
			watchtowerSubgroup.POST("/auto-update", watchtowerHandler.HandleAutoUpdateToggle)
			watchtowerSubgroup.GET("/log-stream", containerHandler.HandleWSLogs)
			watchtowerSubgroup.GET("/log", containerHandler.HandlerContainerLogs)
			watchtowerSubgroup.POST("/load", containerHandler.HandleContainerCreate)
			watchtowerSubgroup.POST("/run", containerHandler.HandleContainerRun)
			watchtowerSubgroup.POST("/load-run", containerHandler.HandleContainerCreate)
			watchtowerSubgroup.POST("/stop-unload", containerHandler.HandleContainerStop)
			watchtowerSubgroup.GET("/all", containerHandler.HandleGetAllContainers)
			watchtowerSubgroup.GET("/default", containerHandler.HandleGetDefaultServices)
			watchtowerSubgroup.GET("/excluded", containerHandler.HandleGetExcludedServices)
		}

		updatesSubgroup := v1.Group("/updates")
		{
			updatesSubgroup.POST("/apply", watchtowerHandler.HandlePostUpdate)
			updatesSubgroup.POST("/download", watchtowerHandler.HandlePostDownload)
			updatesSubgroup.GET("/auto-update", watchtowerHandler.HandleAutoUpdateStatus)
			updatesSubgroup.POST("/auto-update", watchtowerHandler.HandleAutoUpdateToggle)
		}

		containersSubgroup := v1.Group("/containers")
		{
			containersSubgroup.POST("", containerHandler.HandleContainerCreate)
			containersSubgroup.POST("/run", containerHandler.HandleContainerRun)
			containersSubgroup.POST("/stop", containerHandler.HandleContainerStop)
			containersSubgroup.GET("", containerHandler.HandleGetAllContainers)
			containersSubgroup.GET("/inspect", containerHandler.HandleContainerInspect)
			containersSubgroup.GET("/:name/logs", containerHandler.HandlerContainerLogs)
			containersSubgroup.GET("/:name/logs/stream", containerHandler.HandleWSLogs)
		}

		servicesSubgroup := v1.Group("/services")
		{
			servicesSubgroup.GET("/default", containerHandler.HandleGetDefaultServices)
			servicesSubgroup.GET("/excluded", containerHandler.HandleGetExcludedServices)
		}

		signInSubgroup := v1.Group("/signin")
		{
			signInSubgroup.POST("", userHandler.HandleUserSignIn)
			signInSubgroup.GET("/role", userHandler.HandleRole)
		}
	}
}
