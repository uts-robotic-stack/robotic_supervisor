package services

import (
	"fmt"

	"github.com/docker/docker/api/types/mount"
)

type MountConfig struct {
	Type   string `json:"type"`
	Source string `json:"source"`
	Target string `json:"target"`
}

func formatEnvVars(envVars map[string]string) []string {
	var formatted []string
	for key, value := range envVars {
		formatted = append(formatted, fmt.Sprintf("%s=%s", key, value))
	}
	return formatted
}

func formatMounts(mounts []MountConfig) []mount.Mount {
	var formatted []mount.Mount
	for _, m := range mounts {
		formatted = append(formatted, mount.Mount{
			Type:   mount.Type(m.Type),
			Source: m.Source,
			Target: m.Target,
		})
	}
	return formatted
}
