package service

type ServiceMap struct {
	Services map[string]Service `json:"services"`
}

type ServiceIDMap struct {
	ServiceID map[string]string `json:"services"`
}

type Service struct {
	Action          string            `json:"action"`
	Image           string            `json:"image"`
	ContainerID     string            `json:"container_id"`
	Name            string            `json:"name"`
	Tty             bool              `json:"tty"`
	Privileged      bool              `json:"privileged"`
	Restart         string            `json:"restart"`
	Network         string            `json:"network"`
	NetworkSettings NetworkConfig     `json:"network_settings"`
	Mounts          []MountConfig     `json:"mounts"`
	EnvVars         map[string]string `json:"env_vars"`
	Volumes         []VolumeConfig    `json:"volumes"`
	Command         []string          `json:"command"`
	Resources       ResourceConfig    `json:"resources"`
	Sysctls         map[string]string `json:"sysctls"`
}
