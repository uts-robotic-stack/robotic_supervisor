package services

import "github.com/docker/docker/api/types/network"

type ServiceNetwork struct {
	Name    string
	Aliases []string
	IPv4    string
	IPv6    string
}

type Network struct {
	Name           string `json:"name"`
	ID             string
	CheckDuplicate bool
	Labels         Labels
	Internal       bool
	Attachable     bool
	Driver         string `json:"driver"`
	Ipam           network.IPAM
	EnableIPv6     bool
}

func MakeNetworksFromDict(config map[string]interface{}) []ServiceNetwork {
	network := make([]ServiceNetwork, 0)
	if networkOpts, exist := config["networks"].(map[string]interface{}); exist {
		for name, networkData := range networkOpts {
			network = append(network, ServiceNetwork{
				Name: name,
				IPv4: networkData.(map[string]interface{})["ipv4_address"].(string),
			})
		}
	}
	return network
}
