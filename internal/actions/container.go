package actions

import (
	"strconv"
	"strings"

	containerService "github.com/containrrr/watchtower/pkg/container"
	"github.com/containrrr/watchtower/pkg/filters"
	"github.com/containrrr/watchtower/pkg/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/gorilla/websocket"
)

func RunContainer(client containerService.Client, service *containerService.Service) error {
	// Create config
	containerConfig, networkConfig, hostConfig := makeContainerCreateOptions(service, nil)
	_, err := client.StartContainer(
		service.ContainerName, containerConfig, hostConfig, networkConfig)
	return err
}

func StopContainer(client containerService.Client, service *containerService.Service) error {
	containers, _ := client.ListContainers(filters.NoFilter)
	for _, container := range containers {
		// Skip if watchtower
		if container.IsWatchtower() {
			continue
		}
		if container.Name()[1:] == service.ContainerName {
			// 10 seconds stop timeout
			err := client.StopContainer(container, 10)
			return err
		}
	}
	return nil
}

// broadcastLogs reads logs from a Docker container and sends them to the WebSocket connection.
func BroadcastLogs(conn *websocket.Conn, client containerService.Client, name string) {
	containers, _ := client.ListContainers(filters.NoFilter)
	var container types.Container
	foundContainer := false
	for _, cnt := range containers {
		if cnt.Name()[1:] == name {
			container = cnt
			foundContainer = true
		}
	}
	if !foundContainer {
		return
	}

	logs, err := client.StreamLogs(container, true)
	if err != nil {
		return
	}
	defer logs.Close()

	buf := make([]byte, 4096)
	for {
		n, err := logs.Read(buf)
		if err != nil {
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, buf[:n])
		if err != nil {
			return
		}
	}
}

func makeContainerCreateOptions(
	service *containerService.Service,
	network *containerService.Network) (container.Config, network.NetworkingConfig, container.HostConfig) {
	containerConfig := makeContainerConfig(service)
	networkConfig := makeNetworkConfig(service)
	hostConfig := makeHostConfig(service)
	return containerConfig, networkConfig, hostConfig
}

func makeContainerConfig(service *containerService.Service) container.Config {
	return container.Config{
		Hostname:   service.Hostname,
		Domainname: service.Domainname,
		User:       service.User,
		Tty:        service.Tty,
		Cmd:        strslice.StrSlice(service.Command),
		Entrypoint: strslice.StrSlice(service.EntryPoint),
		Image:      service.Image,
		WorkingDir: service.WorkingDir,
		StopSignal: "SIGTERM",
		Env:        service.Environment,
	}
}

func makeNetworkConfig(service *containerService.Service) network.NetworkingConfig {
	// If the current working environment is dev-related
	// the we fuse the service network with host settings
	endPointConfig := map[string]*network.EndpointSettings{}
	return network.NetworkingConfig{
		EndpointsConfig: endPointConfig,
	}
}

func makeHostConfig(service *containerService.Service) container.HostConfig {
	// Prepare binding
	extraHost := make([]string, 0)
	return container.HostConfig{
		AutoRemove:    false,
		Binds:         prepareVolumeBinding(service),
		CapAdd:        service.CapAdd,
		CapDrop:       service.CapDrop,
		ExtraHosts:    extraHost,
		NetworkMode:   container.NetworkMode("host"),
		RestartPolicy: getRestartPolicy(service),
		LogConfig: container.LogConfig{
			Type: "json-file",
		},
		IpcMode:      container.IpcMode(service.IpcMode),
		PortBindings: getPortBinding(service),
		Resources:    getResouces(service),
		Sysctls:      service.Sysctls,
		Privileged:   service.Privileged,
	}
}

func prepareVolumeBinding(service *containerService.Service) []string {
	output := []string{}
	for _, volume := range service.Volumes {
		if len(volume.Source) > 0 && len(volume.Destination) > 0 {
			bindMount := volume.Source + ":" + volume.Destination
			if len(volume.Option) > 0 {
				bindMount = bindMount + ":" + volume.Option
			}
			output = append(output, bindMount)
		}
	}
	return output
}

func getRestartPolicy(service *containerService.Service) container.RestartPolicy {
	restart := container.RestartPolicy{}
	if service.Restart != "" {
		split := strings.Split(service.Restart, ":")
		var attemps int
		if len(split) > 1 {
			attemps, _ = strconv.Atoi(split[1])
		}
		restart.Name = split[0]
		restart.MaximumRetryCount = attemps
	}
	return restart

}

func getPortBinding(service *containerService.Service) nat.PortMap {
	bindingMap := nat.PortMap{}
	for _, port := range service.Ports {
		p := nat.Port(port.Target + "/" + port.Protocol)
		bind := bindingMap[p]
		binding := nat.PortBinding{
			HostIP:   port.HostIp,
			HostPort: port.HostPort,
		}
		bind = append(bind, binding)
		bindingMap[p] = bind
	}
	return bindingMap
}

func getResouces(service *containerService.Service) container.Resources {
	serviceResources := service.Resources
	deviceMappingList := []container.DeviceMapping{}
	for _, device := range service.Devices {
		deviceSplit := strings.Split(device, ":")
		deviceMapping := container.DeviceMapping{
			CgroupPermissions: "rwm",
		}
		switch len(deviceSplit) {
		case 3:
			deviceMapping.CgroupPermissions = deviceSplit[2]
			fallthrough
		case 2:
			deviceMapping.PathInContainer = deviceSplit[1]
			fallthrough
		case 1:
			deviceMapping.PathInContainer = deviceSplit[0]
		}
		deviceMappingList = append(deviceMappingList, deviceMapping)
	}

	resources := container.Resources{
		CgroupParent:   service.CgroupParent,
		OomKillDisable: &serviceResources.OomKillDisable,
		Devices:        deviceMappingList,
		CPUPeriod:      serviceResources.CPUPeriod,
		CPUQuota:       serviceResources.CPUQuota,
		CpusetCpus:     serviceResources.CpusetCpus,
		Memory:         serviceResources.MemoryLimit,
	}

	return resources
}
