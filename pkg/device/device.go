package device

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/containrrr/watchtower/pkg/types"
)

const (
	Offline      = "Offline"
	Online       = "Online"
	Restarting   = "Restarting"
	ShuttingDown = "Shutting down"
	Unknown      = "Unknown"
)

const (
	HasInternet = "Connected"
	NoInternet  = "Disconnected"
)

const (
	Rpi3 = "Raspberry Pi 3B+"
	Rpi4 = "Raspberry Pi 4"
)

const (
	Raspbian = "Raspbian"
)

const (
	Develop = "develop"
	Nightly = "nightly"
	UAT     = "uat"
	Prod    = "production"
)

// TODO: Come up with a better way to handle this
func MakeDevice() (*types.Device, error) {
	status, err := getStatus()
	if err != nil {
		status = Unknown
	}
	uuid, err := getUUID()
	if err != nil {
		uuid = "defaults"
	}

	deviceType, err := getDeviceType()
	if err != nil {
		deviceType = Rpi4
	}

	onlineDuration, err := getUptimeInHours()
	if err != nil {
		onlineDuration = 0
	}

	osType, err := getOsType()
	if err != nil {
		osType = Raspbian
	}

	deviceRole, err := getDeviceRole()
	if err != nil {
		deviceRole = Develop
	}

	internetStatus, err := getInternetStatus()
	if err != nil {
		internetStatus = NoInternet
	}

	return &types.Device{
		Status:         status,
		Uuid:           uuid,
		Type:           deviceType,
		OnlineDuration: onlineDuration,
		OsType:         osType,
		DeviceRole:     deviceRole,
		InternetStatus: internetStatus,
	}, err
}

func getStatus() (string, error) {
	return Online, nil
}

func getUUID() (string, error) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", nil
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "Serial" {
			// The serial number is the second field
			return fields[2], nil
		}
	}
	return "", nil
}

func getDeviceType() (string, error) {
	return Rpi4, nil
}

func getUptimeInHours() (int64, error) {
	cmd := exec.Command("uptime", "-s")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	uptimeString := strings.TrimSpace(string(output))
	uptime, err := time.Parse("2006-01-02 15:04:05", uptimeString)
	if err != nil {
		return 0, err
	}

	return int64(time.Since(uptime).Hours()), nil
}

func getOsType() (string, error) {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			// Extract the distribution identifier
			return strings.Trim(line[3:], `"`), nil
		}
	}

	return "", fmt.Errorf("distribution identifier not found in /etc/os-release")
}

func getDeviceRole() (string, error) {
	return Develop, nil
}

func getInternetStatus() (string, error) {
	// Define a timeout for the HTTP request
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	// Make a GET request to a reliable external server (e.g., Google's public DNS)
	resp, err := client.Get("http://clients3.google.com/generate_204")
	if err != nil {
		return NoInternet, err
	}
	defer resp.Body.Close()

	return HasInternet, nil
}
