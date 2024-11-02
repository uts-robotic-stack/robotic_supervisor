package handlers

import (
	"net/http"

	"github.com/dkhoanguyen/watchtower/internal/actions"
	"github.com/dkhoanguyen/watchtower/pkg/container"
	"github.com/gin-gonic/gin"
	"github.com/godbus/dbus/v5"
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

	go actions.BroadcastHardwareStatus(conn, d.Client, d.HardwareStatusFrequency)
}

func (d *DeviceHandler) HandleShutdown(c *gin.Context) {
	log.Info("Received HTTP request to shutdown device")

	conn, err := dbus.SystemBus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to system bus"})
		return
	}

	// Call the systemd logind interface to power off
	systemd := conn.Object("org.freedesktop.login1", "/org/freedesktop/login1")
	call := systemd.Call("org.freedesktop.login1.Manager.PowerOff", 0, true)
	if call.Err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shutdown the device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device is shutting down"})
}

// RestartHandler allows authenticated users to restart the Raspberry Pi device via D-Bus.
func (d *DeviceHandler) HandleRestart(c *gin.Context) {
	conn, err := dbus.SystemBus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to system bus"})
		return
	}

	// Call the systemd logind interface to reboot
	systemd := conn.Object("org.freedesktop.login1", "/org/freedesktop/login1")
	call := systemd.Call("org.freedesktop.login1.Manager.Reboot", 0, true)
	if call.Err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restart the device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device is restarting"})
}
