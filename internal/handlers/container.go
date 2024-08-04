package handlers

import (
	"net/http"
	"sync"

	"github.com/dkhoanguyen/watchtower/internal/actions"
	containerService "github.com/dkhoanguyen/watchtower/pkg/container"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type ContainerHandler struct {
	client        containerService.Client
	logsFrequency float64
	wsClients     ClientList
	sync.Mutex
}

func NewContainerHandler(client containerService.Client, logFreq float64) *ContainerHandler {
	return &ContainerHandler{
		client:        client,
		logsFrequency: logFreq,
		wsClients:     make(ClientList),
	}
}

// addClient will add clients to our clientList
func (h *ContainerHandler) addClient(client *Client) {
	// Lock so we can manipulate
	h.Lock()
	defer h.Unlock()
	// Add Client
	h.wsClients[client] = true
}

// removeClient will remove the client and clean up
func (h *ContainerHandler) removeClient(client *Client) {
	h.Lock()
	defer h.Unlock()
	// Check if Client exists, then delete it
	if _, ok := h.wsClients[client]; ok {
		// close connection
		client.connection.Close()
		// remove
		delete(h.wsClients, client)
	}
}

func (h *ContainerHandler) HandleWSLogs(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := NewWSClient(conn, h)
	h.addClient(client)

	containerName := c.Query("container")
	go client.readMessages()
	go client.broadcastLogs(containerName)
}

func (h *ContainerHandler) HandleContainerInspect(c *gin.Context) {
	log.Info("Received HTTP request to inspect container")
	// containerName := c.Query("container")
	// containerJSON, err := actions.InspectContainer(h.client, containerName)
	// if err != nil {
	// 	log.Error(err)
	// }
	// config := containerJSON.Config
	// hostConfig := containerJSON.HostConfig
	// networkSettings := containerJSON.NetworkSettings
	// service := &srv.Service{
	// 	Name:          containerJSON.Name,
	// 	Status:        containerJSON.State.Status,
	// 	Action:        "start", // assuming the action is to start the container
	// 	Hostname:      config.Hostname,
	// 	User:          config.User,
	// 	CapAdd:        hostConfig.CapAdd,
	// 	CapDrop:       hostConfig.CapDrop,
	// 	CgroupParent:  hostConfig.CgroupParent,
	// 	Command:       srv.ShellCommand(config.Cmd),
	// 	ContainerName: containerJSON.Name,
	// 	DomainName:    config.Domainname,
	// 	Environment:   config.Env,
	// 	Privileged:    hostConfig.Privileged,
	// 	Restart:       hostConfig.RestartPolicy.Name,
	// 	Tty:           config.Tty,
	// 	WorkingDir:    config.WorkingDir,
	// 	Image:         config.Image,
	// }

	// // // Fill resources
	// // if hostConfig.Resources != (container.Resources{}) {
	// // 	service.Resources = srv.ServiceResources{
	// // 		CPUPeriod:         hostConfig.Resources.CPUPeriod,
	// // 		CPUQuota:          hostConfig.Resources.CPUQuota,
	// // 		CpusetCpus:        hostConfig.Resources.CpusetCpus,
	// // 		CpusetMems:        hostConfig.Resources.CpusetMems,
	// // 		MemoryLimit:       hostConfig.Resources.Memory,
	// // 		MemoryReservation: hostConfig.Resources.MemoryReservation,
	// // 		MemorySwap:        hostConfig.Resources.MemorySwap,
	// // 		MemorySwappiness:  *hostConfig.Resources.MemorySwappiness,
	// // 		OomKillDisable:    *hostConfig.Resources.OomKillDisable,
	// // 	}
	// // }

	// // Fill networks
	// for networkName, networkConfig := range networkSettings.Networks {
	// 	service.Networks = append(service.Networks, srv.ServiceNetwork{
	// 		Name:    networkName,
	// 		Aliases: networkConfig.Aliases,
	// 		IPv4:    networkConfig.IPAddress,
	// 		IPv6:    networkConfig.GlobalIPv6Address,
	// 	})
	// }

	// // // Fill ports
	// // for _, port := range config.ExposedPorts {
	// // 	for _, binding := range networkSettings.Ports[port] {
	// // 		service.Ports = append(service.Ports, srv.ServicePort{
	// // 			Target:   port.Port(),
	// // 			Protocol: port.Proto(),
	// // 			HostIp:   binding.HostIP,
	// // 			HostPort: binding.HostPort,
	// // 		})
	// // 	}
	// // }

	// // Fill volumes
	// for _, mount := range hostConfig.Mounts {
	// 	service.Volumes = append(service.Volumes, srv.ServiceVolume{
	// 		Type:        string(mount.Type),
	// 		Source:      mount.Source,
	// 		Destination: mount.Target,
	// 		Option:      "",
	// 	})
	// }

	c.JSON(http.StatusOK, nil)
}

func (h *ContainerHandler) HandlerContainerLogs(c *gin.Context) {
	log.Info("Received HTTP request to get container logs")
	containerName := c.Query("container")
	output, err := actions.GetLogs(h.client, containerName)
	if err != nil {
		log.Error(err)
	}
	c.JSON(http.StatusOK, output)
}
