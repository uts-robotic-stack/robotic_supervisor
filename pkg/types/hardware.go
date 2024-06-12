package types

type HardwareStatus struct {
	Cpu            float64 `json:"cpu"`
	Ram            float64 `json:"ram"`
	Temperature    float64 `json:"temperature"`
	Storage        float64 `json:"storage"`
	StartupTime    float64 `json:"startup_time"`
	UpTime         float64 `json:"uptime"`
	BatteryLevel   float64 `json:"battery"`
	NetworkTraffic float64 `json:"network_traffic"`
	InternetStatus float64 `json:"internet_status"`
}
