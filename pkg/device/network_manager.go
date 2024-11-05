package device

import (
	"fmt"
	"log"

	"github.com/godbus/dbus/v5"
)

const (
	networkManagerService    = "org.freedesktop.NetworkManager"
	networkManagerPath       = "/org/freedesktop/NetworkManager"
	networkManagerInterface  = "org.freedesktop.NetworkManager"
	networkManagerDeviceType = "org.freedesktop.NetworkManager.Device"
	deviceIP4ConfigInterface = "org.freedesktop.NetworkManager.IP4Config"
)

// NetworkDevice represents a network device with its associated IP addresses and other details.
type NetworkDevice struct {
	Path        dbus.ObjectPath
	DeviceType  string
	IPAddresses []string
}

// GetDevices retrieves all network devices and their associated IP addresses.
func GetDevices() ([]NetworkDevice, error) {
	// Connect to the system bus
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the system bus: %v", err)
	}

	// Get the NetworkManager object
	nmObj := conn.Object(networkManagerService, dbus.ObjectPath(networkManagerPath))

	// Call the method to get all devices
	var devicePaths []dbus.ObjectPath
	err = nmObj.Call(fmt.Sprintf("%s.GetDevices", networkManagerInterface), 0).Store(&devicePaths)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %v", err)
	}

	// Collect information on each device
	var devices []NetworkDevice
	for _, devicePath := range devicePaths {
		device, err := getDeviceInfo(conn, devicePath)
		if err != nil {
			log.Printf("Failed to get information for device %s: %v", devicePath, err)
			continue
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// getDeviceInfo retrieves information about a single device.
func getDeviceInfo(conn *dbus.Conn, devicePath dbus.ObjectPath) (NetworkDevice, error) {
	deviceObj := conn.Object(networkManagerService, devicePath)

	// Get the device type
	var deviceType string
	err := deviceObj.Call(fmt.Sprintf("%s.DeviceType", networkManagerDeviceType), 0).Store(&deviceType)
	if err != nil {
		return NetworkDevice{}, fmt.Errorf("failed to get device type: %v", err)
	}

	// Get the IP4Config object path for the device
	var ip4ConfigPath dbus.ObjectPath
	err = deviceObj.Call(fmt.Sprintf("%s.IP4Config", networkManagerDeviceType), 0).Store(&ip4ConfigPath)
	if err != nil {
		return NetworkDevice{
			Path:       devicePath,
			DeviceType: deviceType,
		}, nil // Device may not have an IP4Config
	}

	// Get IP addresses from the IP4Config object
	ip4ConfigObj := conn.Object(networkManagerService, ip4ConfigPath)
	var ipAddresses []map[string]dbus.Variant
	err = ip4ConfigObj.Call(fmt.Sprintf("%s.Addresses", deviceIP4ConfigInterface), 0).Store(&ipAddresses)
	if err != nil {
		return NetworkDevice{}, fmt.Errorf("failed to get IP addresses: %v", err)
	}

	// Extract IP addresses as strings
	var ips []string
	for _, addr := range ipAddresses {
		ip := addr["address"].Value().(string)
		ips = append(ips, ip)
	}

	return NetworkDevice{
		Path:        devicePath,
		DeviceType:  deviceType,
		IPAddresses: ips,
	}, nil
}
