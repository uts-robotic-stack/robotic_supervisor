package device

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
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

type ConnectionStatus int

const (
	StatusNoConnection ConnectionStatus = iota
	StatusDNSOnly
	StatusHTTPOnly
	StatusFullConnection
)

func (s ConnectionStatus) String() string {
	switch s {
	case StatusNoConnection:
		return "No internet connection. DNS resolution and HTTP requests both failed."
	case StatusDNSOnly:
		return "Internet connection is limited. DNS resolution works but HTTP requests failed. Possible firewall or HTTP restriction."
	case StatusHTTPOnly:
		return "Partial internet connection. HTTP works but DNS is limited or unreliable."
	case StatusFullConnection:
		return "Internet connection is active with full HTTP and DNS access."
	default:
		return "Unknown connection status."
	}
}

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
	devices, err := GetDevices()
	if err != nil {

	}

	serialDevices, err := GetSerialDevices()
	if err != nil {
		serialDevices = make([]string, 0)
	}
	status, _ := CheckInternetConnection()

	return &types.Device{
		Type:            deviceType,
		LastOn:          time.Now().Format(time.RFC3339),
		OnDuration:      onlineDuration,
		SoftwareVersion: "",
		IpAddress:       devices,
		InternetStatus:  int(status),
		Fleet:           "UTS Mechatronics Lab",
		SerialDevices:   serialDevices,
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

// CheckInternetConnection determines the internet connection status.
func CheckInternetConnection() (ConnectionStatus, string) {
	client := http.Client{
		Timeout: 5 * time.Second, // Set a timeout to avoid hanging
	}

	// List of URLs to check for HTTP access
	urls := []string{
		"http://google.com",
		"http://example.com",
		"http://cloudflare.com",
	}

	// Test HTTP connectivity
	for _, url := range urls {
		resp, err := client.Get(url)
		if err == nil {
			// Close response body if no error
			defer resp.Body.Close()
			return StatusFullConnection, StatusFullConnection.String() + fmt.Sprintf(" Successfully connected to %s", url)
		}
	}

	// If HTTP fails, check DNS connectivity
	if checkDNS("8.8.8.8") {
		return StatusDNSOnly, StatusDNSOnly.String()
	}

	return StatusNoConnection, StatusNoConnection.String()
}

// checkDNS performs a DNS lookup to test basic connectivity.
func checkDNS(dnsServer string) bool {
	conn, err := net.DialTimeout("udp", dnsServer+":53", 3*time.Second)
	if err != nil {
		return false // DNS lookup failed
	}
	conn.Close()
	return true // DNS lookup succeeded
}

// GetSerialDevices lists all serial devices in the /dev directory
func GetSerialDevices() ([]string, error) {
	// Define the regular expression to match ttyUSBx and ttyACMx devices
	serialDevicePattern := regexp.MustCompile(`^tty(USB|ACM)\d+$`)

	// Read the /dev directory
	files, err := os.ReadDir("/dev")
	if err != nil {
		return nil, err
	}

	// Filter files matching the pattern
	var serialDevices []string
	for _, file := range files {
		if serialDevicePattern.MatchString(file.Name()) {
			serialDevices = append(serialDevices, file.Name())
		}
	}

	return serialDevices, nil
}
