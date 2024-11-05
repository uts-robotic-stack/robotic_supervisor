package device

import (
	"time"

	"github.com/dkhoanguyen/watchtower/pkg/types"
	"github.com/shirou/gopsutil/v4/host"
)

const (
	Unknown  = "Unknown"
	Rpi3     = "Raspberry Pi 3B+"
	Rpi4_2GB = "Raspberry Pi 4 - 2GB"
	Rpi4_4GB = "Raspberry Pi 4 - 4GB"
	Rpi4_8GB = "Raspberry Pi 4 - 8GB"
	Rpi5_4GB = "Raspberry Pi 5 - 4GB"
	Rpi5_8GB = "Raspberry Pi 5 - 8Gb"
)

// TODO: Come up with a better way to handle this
func MakeDevice() (*types.Device, error) {

	deviceType, err := getDeviceType()
	if err != nil {
		deviceType = Unknown
	}

	onlineDuration, err := getUptimeSeconds()
	if err != nil {
		onlineDuration = 0
	}

	return &types.Device{
		Type:            deviceType,
		LastOn:          time.Now().Format(time.RFC3339),
		OnDuration:      onlineDuration,
		SoftwareVersion: "",
		IpAddress:       "192.168.0.1",
		Fleet:           "UTS Mechatronics Lab",
	}, err
}

func getDeviceType() (string, error) {
	return Rpi4_8GB, nil
}

func getUptimeSeconds() (int64, error) {
	uptimeSeconds, err := host.Uptime()
	if err != nil {
		return 0, err
	}
	return int64(uptimeSeconds), nil
}
