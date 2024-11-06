package actions

import (
	"bytes"
	"errors"
	"io"
	"time"

	containerService "github.com/dkhoanguyen/watchtower/pkg/container"
	"github.com/dkhoanguyen/watchtower/pkg/filters"
	"github.com/dkhoanguyen/watchtower/pkg/types"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func InspectContainer(client containerService.Client, name string) (dockerTypes.ContainerJSON, error) {
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
		return dockerTypes.ContainerJSON{}, errors.New("cannot find container")
	}
	return *container.ContainerInfo(), nil
}

// broadcastLogs reads logs from a Docker container and sends them to the WebSocket connection.
func GetLogs(client containerService.Client, name string) ([]byte, error) {
	output := make([]byte, 0)
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
		return output, errors.New("cannot find container")
	}

	logs, err := client.StreamLogs(container, false, "100")
	if err != nil {
		return output, err
	}
	defer logs.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, logs)
	if err != nil {
		return output, err
	}
	return buf.Bytes(), nil
}

// broadcastLogs reads logs from a Docker container and sends them to the WebSocket connection.
func BroadcastLogs(conn *websocket.Conn, client containerService.Client, name string, freq float64) {
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

	defer func() {
		if err := conn.Close(); err != nil {
			log.Error("Unable to close websocket connection")
		}
		log.Info("Connection closed")
	}()
	var buf bytes.Buffer

	ticker := time.NewTicker(time.Duration((1/freq)*1000) * time.Millisecond) // Adjust the duration as needed
	defer ticker.Stop()

	for range ticker.C {
		logs, err := client.StreamLogs(container, false, "1")
		if err != nil {
			return
		}

		_, err = io.Copy(&buf, logs)
		if err != nil {
			logs.Close()
			return
		}
		logs.Close()

		err = conn.WriteMessage(websocket.TextMessage, buf.Bytes())
		if err != nil {
			logs.Close()
			log.Error("Error writing websocket message")
			return
		}
		buf.Reset()
	}
}
