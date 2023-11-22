package actions

import (
	"encoding/json"
	"time"

	containerService "github.com/containrrr/watchtower/pkg/container"
	"github.com/containrrr/watchtower/pkg/device"
	"github.com/containrrr/watchtower/pkg/filters"
	"github.com/containrrr/watchtower/pkg/types"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func GetDeviceInfo(client containerService.Client) types.Device {
	device, err := device.MakeDevice()
	if err != nil {
		log.Error(err)
	}
	containers, _ := client.ListContainers(filters.NoFilter)
	for _, container := range containers {
		// Get watchtower release
		if container.IsWatchtower() && container.HasImageInfo() {
			device.SupervisorRelease = container.ImageInfo().ID
			continue
		}
	}
	return *device
}

func BroadcastHardwareStatus(conn *websocket.Conn, client containerService.Client, freq float64) {
	for {
		resources, err := device.GetHardwareStatus()
		if err != nil {
			return
		}
		data, _ := json.Marshal(resources)
		conn.WriteMessage(websocket.TextMessage, data)
		time.Sleep(time.Duration(1/freq) * time.Second)
	}
}
