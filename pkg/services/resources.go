package services

import (
	"strconv"
	"strings"
)

type ServicePort struct {
	Target   string
	Protocol string
	HostIp   string
	HostPort string
}

type ServiceResources struct {
	CPUPeriod         int64
	CPUQuota          int64
	CpusetCpus        string
	CpusetMems        string
	MemoryLimit       int64
	MemoryReservation int64
	MemorySwap        int64
	MemorySwappiness  int64
	OomKillDisable    bool
}

func MakeDependsOnFromDict(config map[string]interface{}) []string {
	output := make([]string, 0)
	if dependsOnOpt, exist := config["depends_on"].([]interface{}); exist {
		for _, dependsOn := range dependsOnOpt {
			output = append(output, dependsOn.(string))
		}
	}
	return output
}

func MakeDeployResourcesFromDict(config map[string]interface{}) ServiceResources {
	resources := ServiceResources{}
	if deployOpt, exist := config["deploy"].(map[string]interface{}); exist {
		if resourcesOpt, exist := deployOpt["resources"].(map[string]interface{}); exist {
			limitOpt := resourcesOpt["limits"].(map[string]interface{})
			// CPU usage
			var cpuPeriod float64 = 100000                                   // Default value of 100000
			cpuQuota, _ := strconv.ParseFloat(limitOpt["cpus"].(string), 64) // Combination of period and quota to determine cpu limitation
			resources.CPUQuota = int64(cpuQuota * cpuPeriod)
			resources.CPUPeriod = int64(cpuPeriod)

			// Memory usage
			memoryInString := limitOpt["memory"].(string)
			memory, _ := strconv.ParseInt(memoryInString[:len(memoryInString)-1], 10, 64)
			suffix := string(memoryInString[len(memoryInString)-1])
			switch {
			case suffix == "k" || suffix == "K":
				memory = memory * 1024
			case suffix == "m" || suffix == "M":
				memory = memory * 1048576
			case suffix == "g" || suffix == "G":
				memory = memory * 1073741824
			}
			resources.MemoryLimit = memory
		}
	}
	return resources
}

func MakePortBindingFromDict(config map[string]interface{}) []ServicePort {
	ports := make([]ServicePort, 0)
	if portOpt, exist := config["ports"].([]interface{}); exist {
		for _, portData := range portOpt {
			// We need to properly split the string to port and host ip address
			splittedPort := strings.Split(portData.(string), ":")
			port := ServicePort{
				Target:   splittedPort[0],
				Protocol: "tcp",
				HostIp:   "0.0.0.0",
				HostPort: splittedPort[1],
			}
			ports = append(ports, port)
		}
	}
	return ports
}
