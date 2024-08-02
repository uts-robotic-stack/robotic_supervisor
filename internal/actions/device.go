package actions

import (
	"encoding/json"
	"time"

	container "github.com/containrrr/watchtower/pkg/container"
	"github.com/containrrr/watchtower/pkg/device"
	"github.com/containrrr/watchtower/pkg/filters"
	"github.com/containrrr/watchtower/pkg/types"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func GetDeviceInfo(client container.Client) types.Device {
	device, err := device.MakeDevice()
	if err != nil {
		log.Error(err)
	}
	containers, _ := client.ListContainers(filters.NoFilter)
	for _, container := range containers {
		// Get supervisor release
		if container.IsWatchtower() && container.HasImageInfo() {
			device.SoftwareVersion = container.ImageInfo().ID
			break
		}
	}
	return *device
}

func BroadcastHardwareStatus(conn *websocket.Conn, client container.Client, freq float64) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Error("Unable to close websocket connection")
		}
		log.Info("Connection closed")
	}()

	for {
		resources, err := device.GetHardwareStatus()
		if err != nil {
			return
		}
		data, _ := json.Marshal(resources)
		err = conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			return
		}
		time.Sleep(time.Duration(1/freq) * time.Second)
	}
}
