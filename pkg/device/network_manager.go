package device

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/godbus/dbus/v5"
)

const (
	networkManagerService    = "org.freedesktop.NetworkManager"
	networkManagerPath       = "/org/freedesktop/NetworkManager"
	networkManagerInterface  = "org.freedesktop.NetworkManager"
	networkManagerDeviceType = "org.freedesktop.NetworkManager.Device"
	deviceIP4ConfigInterface = "org.freedesktop.NetworkManager.IP4Config"
)

type DeviceType int

const (
	NM_DEVICE_TYPE_UNKNOWN     DeviceType = 0  // unknown device
	NM_DEVICE_TYPE_ETHERNET    DeviceType = 1  // a wired ethernet device
	NM_DEVICE_TYPE_WIFI        DeviceType = 2  // an 802.11 WiFi device
	NM_DEVICE_TYPE_UNUSED1     DeviceType = 3  // not used
	NM_DEVICE_TYPE_UNUSED2     DeviceType = 4  // not used
	NM_DEVICE_TYPE_BT          DeviceType = 5  // a Bluetooth device supporting PAN or DUN access protocols
	NM_DEVICE_TYPE_OLPC_MESH   DeviceType = 6  // an OLPC XO mesh networking device
	NM_DEVICE_TYPE_WIMAX       DeviceType = 7  // an 802.16e Mobile WiMAX broadband device
	NM_DEVICE_TYPE_MODEM       DeviceType = 8  // a modem supporting analog telephone, CDMA/EVDO, GSM/UMTS, or LTE network access protocols
	NM_DEVICE_TYPE_INFINIBAND  DeviceType = 9  // an IP-over-InfiniBand device
	NM_DEVICE_TYPE_BOND        DeviceType = 10 // a bond master interface
	NM_DEVICE_TYPE_VLAN        DeviceType = 11 // an 802.1Q VLAN interface
	NM_DEVICE_TYPE_ADSL        DeviceType = 12 // ADSL modem
	NM_DEVICE_TYPE_BRIDGE      DeviceType = 13 // a bridge master interface
	NM_DEVICE_TYPE_GENERIC     DeviceType = 14 // generic support for unrecognized device types
	NM_DEVICE_TYPE_TEAM        DeviceType = 15 // a team master interface
	NM_DEVICE_TYPE_TUN         DeviceType = 16 // a TUN or TAP interface
	NM_DEVICE_TYPE_IP_TUNNEL   DeviceType = 17 // an IP tunnel interface
	NM_DEVICE_TYPE_MACVLAN     DeviceType = 18 // a MACVLAN interface
	NM_DEVICE_TYPE_VXLAN       DeviceType = 19 // a VXLAN interface
	NM_DEVICE_TYPE_VETH        DeviceType = 20 // a VETH interface
)


// NetworkDevice represents a network device with its associated IP addresses and other details.
type NetworkDevice struct {
	Path        dbus.ObjectPath
	DeviceType  string
	DeviceName 	string
	IPAddresses string
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

	// Retrieve the DeviceType property using GetProperty
	var deviceType uint32
	err := deviceObj.Call("org.freedesktop.DBus.Properties.Get", 0, networkManagerDeviceType, "DeviceType").Store(&deviceType)
	if err != nil {
		return NetworkDevice{}, fmt.Errorf("failed to get device type: %v", err)
	}

	// Retrieve the IP4Config path for the device
	var ip4ConfigPath dbus.ObjectPath
	err = deviceObj.Call("org.freedesktop.DBus.Properties.Get", 0, networkManagerDeviceType, "Ip4Config").Store(&ip4ConfigPath)
	if err != nil || ip4ConfigPath == "" {
		return err, nil
	}

	// Retrieve the Interface property (e.g., "eth0" or "wlan0")
	var interfaceName string
	err = deviceObj.Call("org.freedesktop.DBus.Properties.Get", 0, networkManagerDeviceType, "Interface").Store(&interfaceName)
	if err != nil {
		return NetworkDevice{}, fmt.Errorf("failed to get device interface name: %v", err)
	}
	fmt.Println(interfaceName)
	// Retrieve the Addresses property from IP4Config using GetProperty
	ip4ConfigObj := conn.Object(networkManagerService, ip4ConfigPath)
	var ipAddresses [][]uint32
	err = ip4ConfigObj.Call("org.freedesktop.DBus.Properties.Get", 0, deviceIP4ConfigInterface, "Addresses").Store(&ipAddresses)
	if err != nil {
		return NetworkDevice{}, fmt.Errorf("failed to get IP addresses: %v", err)
	}
	// Convert each IP address from uint32 to a human-readable format
	var ips []string
	for _, addr := range ipAddresses {
		if len(addr) > 0 {
			ip := convertUint32ToIP(addr[0])
			ips = append(ips, ip)
		}
	}

	data := make([]string,0)
	return NetworkDevice{
		Path:        devicePath,
		DeviceType:  "",
		IPAddresses: data,
	}, nil
}

// convertUint32ToIP converts a uint32 IP address to its string representation.
func convertUint32ToIP(ipUint32 uint32) string {
	ipBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(ipBytes, ipUint32)
	return net.IP(ipBytes).String()
}