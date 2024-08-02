package services

import "strings"

const (
	VolumeTypeBind  = "bind"
	VolumeTypeMount = "mount"
)

type Volume struct {
	Name string `json:"name"`
}

type ServiceVolume struct {
	Type        string
	Source      string
	Destination string
	Option      string
}

func MakeVolumesFromDict(config map[string]interface{}) []ServiceVolume {
	volumes := make([]ServiceVolume, 0)
	if volumeOpt, exist := config["volumes"].([]interface{}); exist {
		fromStringToVolume := func(volStr string) ServiceVolume {
			separateValues := strings.Split(volStr, ":")
			if len(separateValues) >= 2 {
				return ServiceVolume{
					Type:        VolumeTypeBind,
					Source:      separateValues[0],
					Destination: separateValues[1],
				}
			}
			return ServiceVolume{}
		}
		for _, volData := range volumeOpt {
			volumes = append(volumes, fromStringToVolume(volData.(string)))
		}
	}
	return volumes
}
