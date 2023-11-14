package types

type Device struct {
	Status            string
	Uuid              string
	Type              string
	OnlineDuration    int64
	OsType            string
	DeviceRole        string
	InternetStatus    string
	SupervisorRelease string
}
