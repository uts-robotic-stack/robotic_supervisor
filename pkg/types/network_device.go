package types

import (
	"fmt"

	"github.com/godbus/dbus/v5"
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

type NetworkDevice struct {
	// Ignore Path as this is only for internal use
	Path       dbus.ObjectPath `json:"-" yaml:"-"`
	DeviceType string          `json:"device_type" yaml:"device_type"`
	DeviceName string          `json:"device_name" yaml:"device_name"`
	IPAddress  string          `json:"ip_address" yaml:"ip_address"`
}
