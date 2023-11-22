package types

type Device struct {
	Status            string `json:"status"`
	Uuid              string `json:"uuid"`
	Type              string `json:"device_type"`
	OnlineDuration    int64  `json:"online_duration"`
	OsType            string `json:"os_type"`
	DeviceRole        string `json:"device_role"`
	InternetStatus    string `json:"internet_status"`
	SupervisorRelease string `json:"supervisor_release"`
}

type HardwareStatus struct {
	Cpu         float64 `json:"cpu"`
	Ram         float64 `json:"ram"`
	Temperature float64 `json:"temperature"`
	Storage     float64 `json:"storage"`
}
