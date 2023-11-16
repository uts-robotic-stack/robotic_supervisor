package actions

import (
	containerService "github.com/containrrr/watchtower/pkg/container"
	"github.com/containrrr/watchtower/pkg/device"
	"github.com/containrrr/watchtower/pkg/filters"
	"github.com/containrrr/watchtower/pkg/types"
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
