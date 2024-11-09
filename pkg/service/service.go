package service

type ServiceMap struct {
	Services map[string]Service `yaml:"services" json:"services"`
}

type ServiceIDMap struct {
	ServiceID map[string]string `yaml:"services" json:"services"`
}

type Service struct {
	Action          string            `yaml:"action" json:"action"`
	Image           Image             `yaml:"image" json:"image"`
	ContainerID     string            `yaml:"container_id" json:"container_id"`
	Name            string            `yaml:"name" json:"name"`
	Tty             bool              `yaml:"tty" json:"tty"`
	Privileged      bool              `yaml:"privileged" json:"privileged"`
	Restart         string            `yaml:"restart" json:"restart"`
	Network         string            `yaml:"network" json:"network"`
	NetworkSettings NetworkConfig     `yaml:"network_settings" json:"network_settings"`
	Mounts          []MountConfig     `yaml:"mounts" json:"mounts"`
	EnvVars         map[string]string `yaml:"env_vars" json:"env_vars"`
	Volumes         []VolumeConfig    `yaml:"volumes" json:"volumes"`
	Command         []string          `yaml:"command" json:"command"`
	Resources       ResourceConfig    `yaml:"resources" json:"resources"`
	Sysctls         map[string]string `yaml:"sysctls" json:"sysctls"`
	Status          string            `yaml:"status" json:"status"`
	Labels          map[string]string `yaml:"labels" json:"labels"`
}

type Image struct {
	Name string `yaml:"name" json:"name"`
	ID   string `yaml:"id" json:"id"`
}
