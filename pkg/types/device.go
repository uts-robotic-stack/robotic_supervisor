package types

type Device struct {
	Type            string `json:"device_type"`
	OnDuration      int64  `json:"on_duration"`
	LastOn          string `json:"last_on"`
	SoftwareVersion string `json:"software_version"`
	IpAddress       string `json:"ip_address"`
	Fleet           string `json:"fleet"`
}
