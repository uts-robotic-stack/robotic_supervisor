package services

import (
	"fmt"
	"strings"
)

type ServiceBuild struct {
	Context    string
	Dockerfile string
	Args       map[string]*string
}

const (
	RestartAlways        = "always"
	RestartOnFailure     = "on-failure"
	RestartNoRetry       = "no"
	RestartUnlessStopped = "unless-stopped"
)

const (
	ActionUnload = "unload"
	ActionRun    = "start"
	ActionStop   = "stop"
)

type ShellCommand []string
type Labels map[string]string

type Image struct {
	ID      string
	Name    string
	Tag     string
	Created string
}

func MakeBuildOptFromDict(config map[string]interface{}, path string) ServiceBuild {
	output := ServiceBuild{}
	buildOpt := config["build"].(map[string]interface{})
	output.Context = path
	output.Dockerfile = buildOpt["dockerfile"].(string)

	// Only extract build arg if arg exists
	if buildArgs, exist := buildOpt["args"].([]interface{}); exist {
		formattedArg := make(map[string]*string)
		for _, arg := range buildArgs {
			if _, ok := arg.(string); ok {
				splittedString := strings.Split(arg.(string), "=")
				key := splittedString[0]
				value := arg.(string)[len(key+"="):]
				formattedArg[key] = &value
			}
		}
		output.Args = formattedArg
	}

	return output
}

func MakeContainerNameFromDict(config map[string]interface{}) string {
	return config["container_name"].(string)
}

func MakeActionFromDict(config map[string]interface{}) string {
	return config["action"].(string)
}

func MakeCommandFromDict(config map[string]interface{}, cmdType string) ShellCommand {
	output := make(ShellCommand, 0)
	if cmdOpt, exist := config[cmdType].([]interface{}); exist {
		for _, args := range cmdOpt {
			output = append(output, args.(string))
		}
	} else {
		return nil
	}
	return output
}

func MakeEnviromentFromDict(config map[string]interface{}) []string {
	env := make([]string, 0)
	// Environment variables
	if envVarsOpt, exist := config["environment"].([]interface{}); exist {
		for _, envVars := range envVarsOpt {
			env = append(env, envVars.(string))
		}
	}
	return env
}

func MakeExtraHosts(config map[string]interface{}, hostname string) []string {
	hosts := make([]string, 0)
	if extraHostsOpt, exist := config["extra_hosts"].([]interface{}); exist {
		for _, host := range extraHostsOpt {
			hosts = append(hosts, host.(string))
		}

	}
	// Default expose host machine
	hosts = append(hosts, fmt.Sprintf("%s:127.0.0.1", hostname))
	return hosts
}

func MakeRestartOptFromDict(config map[string]interface{}) string {
	output := "no"
	if restartOpt, exist := config["restart"].(string); exist {
		output = restartOpt
	}
	return output
}

func MakePrivilegedFromDict(config map[string]interface{}) bool {
	if privileged, exist := config["privileged"].(bool); exist {
		return privileged
	}
	return false
}

func MakeTTYFromDict(config map[string]interface{}) bool {
	if tty, exist := config["tty"].(bool); exist {
		return tty
	}
	return false
}

func MakeImageFromDict(config map[string]interface{}) string {
	output := ""
	if image, exist := config["image"].(string); exist {
		output = image
	}
	return output
}
