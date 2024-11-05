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
	NM_DEVICE_TYPE_UNKNOWN    DeviceType = 0  // unknown device
	NM_DEVICE_TYPE_ETHERNET   DeviceType = 1  // a wired ethernet device
	NM_DEVICE_TYPE_WIFI       DeviceType = 2  // an 802.11 WiFi device
	NM_DEVICE_TYPE_UNUSED1    DeviceType = 3  // not used
	NM_DEVICE_TYPE_UNUSED2    DeviceType = 4  // not used
	NM_DEVICE_TYPE_BT         DeviceType = 5  // a Bluetooth device supporting PAN or DUN access protocols
	NM_DEVICE_TYPE_OLPC_MESH  DeviceType = 6  // an OLPC XO mesh networking device
	NM_DEVICE_TYPE_WIMAX      DeviceType = 7  // an 802.16e Mobile WiMAX broadband device
	NM_DEVICE_TYPE_MODEM      DeviceType = 8  // a modem supporting analog telephone, CDMA/EVDO, GSM/UMTS, or LTE network access protocols
	NM_DEVICE_TYPE_INFINIBAND DeviceType = 9  // an IP-over-InfiniBand device
	NM_DEVICE_TYPE_BOND       DeviceType = 10 // a bond master interface
	NM_DEVICE_TYPE_VLAN       DeviceType = 11 // an 802.1Q VLAN interface
	NM_DEVICE_TYPE_ADSL       DeviceType = 12 // ADSL modem
	NM_DEVICE_TYPE_BRIDGE     DeviceType = 13 // a bridge master interface
	NM_DEVICE_TYPE_GENERIC    DeviceType = 14 // generic support for unrecognized device types
	NM_DEVICE_TYPE_TEAM       DeviceType = 15 // a team master interface
	NM_DEVICE_TYPE_TUN        DeviceType = 16 // a TUN or TAP interface
	NM_DEVICE_TYPE_IP_TUNNEL  DeviceType = 17 // an IP tunnel interface
	NM_DEVICE_TYPE_MACVLAN    DeviceType = 18 // a MACVLAN interface
	NM_DEVICE_TYPE_VXLAN      DeviceType = 19 // a VXLAN interface
	NM_DEVICE_TYPE_VETH       DeviceType = 20 // a VETH interface
)

func (dt DeviceType) String() string {
	switch dt {
	case NM_DEVICE_TYPE_UNKNOWN:
		return "Unknown device"
	case NM_DEVICE_TYPE_ETHERNET:
		return "Wired Ethernet device"
	case NM_DEVICE_TYPE_WIFI:
		return "WiFi device"
	case NM_DEVICE_TYPE_BT:
		return "Bluetooth device"
	case NM_DEVICE_TYPE_OLPC_MESH:
		return "OLPC XO mesh networking device"
	case NM_DEVICE_TYPE_WIMAX:
		return "Mobile WiMAX broadband device"
	case NM_DEVICE_TYPE_MODEM:
		return "Modem"
	case NM_DEVICE_TYPE_INFINIBAND:
		return "IP-over-InfiniBand device"
	case NM_DEVICE_TYPE_BOND:
		return "Bond master interface"
	case NM_DEVICE_TYPE_VLAN:
		return "802.1Q VLAN interface"
	case NM_DEVICE_TYPE_ADSL:
		return "ADSL modem"
	case NM_DEVICE_TYPE_BRIDGE:
		return "Bridge master interface"
	case NM_DEVICE_TYPE_GENERIC:
		return "Generic device"
	case NM_DEVICE_TYPE_TEAM:
		return "Team master interface"
	case NM_DEVICE_TYPE_TUN:
		return "TUN or TAP interface"
	case NM_DEVICE_TYPE_IP_TUNNEL:
		return "IP tunnel interface"
	case NM_DEVICE_TYPE_MACVLAN:
		return "MACVLAN interface"
	case NM_DEVICE_TYPE_VXLAN:
		return "VXLAN interface"
	case NM_DEVICE_TYPE_VETH:
		return "VETH interface"
	default:
		return fmt.Sprintf("Unknown device type (%d)", dt)
	}
}

func NewDeviceType(num int) DeviceType {
	switch num {
	case 0:
		return NM_DEVICE_TYPE_UNKNOWN
	case 1:
		return NM_DEVICE_TYPE_ETHERNET
	case 2:
		return NM_DEVICE_TYPE_WIFI
	case 3:
		return NM_DEVICE_TYPE_UNUSED1
	case 4:
		return NM_DEVICE_TYPE_UNUSED2
	case 5:
		return NM_DEVICE_TYPE_BT
	case 6:
		return NM_DEVICE_TYPE_OLPC_MESH
	case 7:
		return NM_DEVICE_TYPE_WIMAX
	case 8:
		return NM_DEVICE_TYPE_MODEM
	case 9:
		return NM_DEVICE_TYPE_INFINIBAND
	case 10:
		return NM_DEVICE_TYPE_BOND
	case 11:
		return NM_DEVICE_TYPE_VLAN
	case 12:
		return NM_DEVICE_TYPE_ADSL
	case 13:
		return NM_DEVICE_TYPE_BRIDGE
	case 14:
		return NM_DEVICE_TYPE_GENERIC
	case 15:
		return NM_DEVICE_TYPE_TEAM
	case 16:
		return NM_DEVICE_TYPE_TUN
	case 17:
		return NM_DEVICE_TYPE_IP_TUNNEL
	case 18:
		return NM_DEVICE_TYPE_MACVLAN
	case 19:
		return NM_DEVICE_TYPE_VXLAN
	case 20:
		return NM_DEVICE_TYPE_VETH
	default:
		return NM_DEVICE_TYPE_UNKNOWN
	}
}

// NetworkDevice represents a network device with its associated IP addresses and other details.
type NetworkDevice struct {
	Path       dbus.ObjectPath
	DeviceType string
	DeviceName string
	IPAddress  string
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
		// return nil, err
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
	var ip = ""
	for _, addr := range ipAddresses {
		if len(addr) > 0 {
			ip = convertUint32ToIP(addr[0])
		}
	}

	return NetworkDevice{
		Path:       devicePath,
		DeviceType: NewDeviceType(int(deviceType)).String(),
		IPAddress:  ip,
	}, nil
}

// convertUint32ToIP converts a uint32 IP address to its string representation.
func convertUint32ToIP(ipUint32 uint32) string {
	ipBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(ipBytes, ipUint32)
	return net.IP(ipBytes).String()
}
