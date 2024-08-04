package service

import "fmt"

type VolumeConfig struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func FormatVolumes(volumes []VolumeConfig) []string {
	var formatted []string
	for _, v := range volumes {
		formatted = append(formatted, fmt.Sprintf("%s:%s", v.Source, v.Target))
	}
	return formatted
}
