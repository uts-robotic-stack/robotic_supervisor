package device

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/dkhoanguyen/watchtower/pkg/types"
	"github.com/godbus/dbus/v5"
)

const (
	networkManagerService    = "org.freedesktop.NetworkManager"
	networkManagerPath       = "/org/freedesktop/NetworkManager"
	networkManagerInterface  = "org.freedesktop.NetworkManager"
	networkManagerDeviceType = "org.freedesktop.NetworkManager.Device"
	deviceIP4ConfigInterface = "org.freedesktop.NetworkManager.IP4Config"
)

// GetDevices retrieves all network devices and their associated IP addresses.
func GetDevices() (map[string]types.NetworkDevice, error) {
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

	devices := make(map[string]types.NetworkDevice)
	for _, devicePath := range devicePaths {
		device, err := getDeviceInfo(conn, devicePath)
		if err != nil {
			continue
		}
		// Assuming device.DeviceName is unique and can be used as the key
		devices[device.DeviceName] = device
	}
	return devices, nil
}

// getDeviceInfo retrieves information about a single device.
func getDeviceInfo(conn *dbus.Conn, devicePath dbus.ObjectPath) (types.NetworkDevice, error) {
	deviceObj := conn.Object(networkManagerService, devicePath)

	// Retrieve the DeviceType property using GetProperty
	var deviceType uint32
	err := deviceObj.Call("org.freedesktop.DBus.Properties.Get", 0, networkManagerDeviceType, "DeviceType").Store(&deviceType)
	if err != nil {
		return types.NetworkDevice{}, fmt.Errorf("failed to get device type: %v", err)
	}

	// Retrieve the IP4Config path for the device
	var ip4ConfigPath dbus.ObjectPath
	err = deviceObj.Call("org.freedesktop.DBus.Properties.Get", 0, networkManagerDeviceType, "Ip4Config").Store(&ip4ConfigPath)
	if err != nil || ip4ConfigPath == "" {
		return types.NetworkDevice{}, fmt.Errorf("failed to get device interface name: %v", err)
	}

	// Retrieve the Interface property (e.g., "eth0" or "wlan0")
	var interfaceName string
	err = deviceObj.Call("org.freedesktop.DBus.Properties.Get", 0, networkManagerDeviceType, "Interface").Store(&interfaceName)
	if err != nil {
		return types.NetworkDevice{}, fmt.Errorf("failed to get device interface name: %v", err)
	}

	// Retrieve the Addresses property from IP4Config using GetProperty
	ip4ConfigObj := conn.Object(networkManagerService, ip4ConfigPath)
	var ipAddresses [][]uint32
	err = ip4ConfigObj.Call("org.freedesktop.DBus.Properties.Get", 0, deviceIP4ConfigInterface, "Addresses").Store(&ipAddresses)
	if err != nil {
		return types.NetworkDevice{}, fmt.Errorf("failed to get IP addresses: %v", err)
	}
	// Convert each IP address from uint32 to a human-readable format
	var ip = ""
	for _, addr := range ipAddresses {
		if len(addr) > 0 {
			ip = convertUint32ToIP(addr[0])
		}
	}

	return types.NetworkDevice{
		Path:       devicePath,
		DeviceType: types.NewDeviceType(int(deviceType)).String(),
		DeviceName: interfaceName,
		IPAddress:  ip,
	}, nil
}

// convertUint32ToIP converts a uint32 IP address to its string representation.
func convertUint32ToIP(ipUint32 uint32) string {
	ipBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(ipBytes, ipUint32)
	return net.IP(ipBytes).String()
}
