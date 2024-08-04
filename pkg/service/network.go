package service

type NetworkConfig struct {
	IPAddress string `json:"ip_address"`
	Gateway   string `json:"gateway"`
	Subnet    string `json:"subnet"`
}
