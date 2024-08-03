package services

const (
	Unknown  = "Unknown"
	Unloaded = "Unload"
	Stopped  = "Stopped"
	Running  = "Running"
)

type Service struct {
	Name          string            `json:"name"`
	Status        string            `json:"status"`
	Action        string            `json:"action"`
	Hostname      string            `json:"hostname"`
	User          string            `json:"user"`
	CapAdd        []string          `json:"cap_add"`
	CapDrop       []string          `json:"cap_drop"`
	BuildOpt      ServiceBuild      `json:"build_opt"`
	CgroupParent  string            `json:"cgroup_parent"`
	Command       ShellCommand      `json:"command"`
	ContainerName string            `json:"container_name"`
	Domainname    string            `json:"domain_name"`
	DependsOn     []string          `json:"depends_on"`
	Devices       []string          `json:"devices"`
	EntryPoint    ShellCommand      `json:"entrypoint"`
	Environment   []string          `json:"environment"`
	EnvFile       []string          `json:"env_file"`
	Expose        []string          `json:"expose"`
	ExtraHosts    []string          `json:"extra_hosts"`
	IpcMode       string            `json:"ipc_mode"`
	Resources     ServiceResources  `json:"resources"`
	Networks      []ServiceNetwork  `json:"networks"`
	NetworkMode   string            `json:"network_mode"`
	Ports         []ServicePort     `json:"ports"`
	Privileged    bool              `json:"privileged"`
	Sysctls       map[string]string `json:"sysctls"`
	Restart       string            `json:"restart"`
	Tty           bool              `json:"tty"`
	Volumes       []ServiceVolume   `json:"volumes"`
	WorkingDir    string            `json:"working_dir"`
	Image         string            `json:"image"`
}

func (s Service) GetName() string {
	return s.Name
}

func (s Service) GetStatus() string {
	return s.Status
}

func MakeServiceFromDict(config map[string]interface{},
	name string) Service {
	output := Service{}
	output.Name = name
	// Image
	output.Image = MakeImageFromDict(config)
	// Action
	output.Action = MakeActionFromDict(config)
	// Container Name
	output.ContainerName = MakeContainerNameFromDict(config)
	// Commands
	output.Command = MakeCommandFromDict(config, "command")
	// Dependencies
	output.DependsOn = MakeDependsOnFromDict(config)
	// Deployment and Resources
	output.Resources = MakeDeployResourcesFromDict(config)
	// Entrypoint
	output.EntryPoint = MakeCommandFromDict(config, "entrypoint")
	// Environment Variables
	output.Environment = MakeEnviromentFromDict(config)
	// Network
	output.Networks = MakeNetworksFromDict(config)
	// Ports
	output.Ports = MakePortBindingFromDict(config)
	// Privileged
	output.Privileged = MakePrivilegedFromDict(config)
	// Restart
	output.Restart = MakeRestartOptFromDict(config)
	// TTY
	output.Tty = MakeTTYFromDict(config)
	// Volumes
	output.Volumes = MakeVolumesFromDict(config)
	return output
}

func MakeCommonSettings() []string {
	settings := make([]string, 0)
	settings = append(settings, "ROS_MASTER_URI=\"http://localhost:11311\"")
	settings = append(settings, "ROS_IP=\"192.168.27.1\"")
	return settings
}
