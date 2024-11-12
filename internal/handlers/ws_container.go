package handlers

import (
	"bufio"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/dkhoanguyen/watchtower/pkg/filters"
	"github.com/dkhoanguyen/watchtower/pkg/types"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// ClientList is a map used to help manage a map of clients
type ClientList map[*Client]bool

// Client is a websocket client, basically a frontend visitor
type Client struct {
	// the websocket connection
	connection *websocket.Conn
	handler    *ContainerHandler
	egress     chan []byte
}

var (
	// pongWait is how long we will await a pong response from client
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval  = 100 * time.Millisecond
	writeInterval = 100 * time.Millisecond
)

func NewWSClient(conn *websocket.Conn, handler *ContainerHandler) *Client {
	return &Client{
		connection: conn,
		handler:    handler,
		egress:     make(chan []byte),
	}
}

func (c *Client) broadcastLogs(containerName string) {
	pingTicker := time.NewTicker(pingInterval)
	writeTicker := time.NewTicker(writeInterval)

	containers, _ := c.handler.client.ListContainers(filters.NoFilter)
	var container types.Container
	foundContainer := false
	for _, cnt := range containers {
		if cnt.Name()[1:] == containerName {
			container = cnt
			foundContainer = true
		}
	}
	if !foundContainer {
		return
	}
	fmt.Println(container.Name())
	// var buf bytes.Buffer
	logs, err := c.handler.client.StreamLogs(container, true, "20")
	if err != nil {
		return
	}
	defer logs.Close()
	scanner := bufio.NewScanner(logs)

	healthy := true

	defer func() {
		// Graceful Close the Connection once this
		// function is done
		pingTicker.Stop()
		writeTicker.Stop()
		healthy = false
		c.handler.removeClient(c)
	}()

	go func() {
		for scanner.Scan() && healthy {
			message := scanner.Text()

			// Split the message into lines
			lines := strings.Split(message, "\n")

			// Write each line to WebSocket
			for _, line := range lines {
				c.egress <- []byte(line)
				time.Sleep(time.Duration(100 * time.Millisecond))
			}

		}
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				// Manager has closed this connection channel, so communicate that to frontend
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					// Log that the connection is closed and the reason
					log.Println("connection closed: ", err)
				}
				// Return to close the goroutine
				return
			}
			message = sanitizeAndClean(message)
			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Error(err)
			}
		case <-pingTicker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Error("writemsg: ", err)
				return // return to break this goroutine triggeing cleanup
			}
		}
	}
}

// readMessages will start the client to read messages and handle them
// appropriatly.
// This is suppose to be ran as a goroutine
func (c *Client) readMessages() {
	defer func() {
		// Graceful Close the Connection once this
		// function is done
		log.Info("Close connection from read side")
		c.handler.removeClient(c)
	}()

	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	// Configure how to handle Pong responses
	c.connection.SetPongHandler(c.pongHandler)

	// Loop Forever
	for {
		// ReadMessage is used to read the next message in queue
		// in the connection
		_, _, err := c.connection.ReadMessage()

		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break // Break the loop to close conn & Cleanup
		}
	}
}

// pongHandler is used to handle PongMessages for the Client
func (c *Client) pongHandler(pongMsg string) error {
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

// sanitizeAndClean removes all non-printable and non-UTF-8 characters from the input byte slice
func sanitizeAndClean(data []byte) []byte {
	var sanitized strings.Builder
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		if r == utf8.RuneError && size == 1 {
			// Invalid UTF-8 byte, skip it
			data = data[1:]
			continue
		}
		// Check if rune is printable ASCII or common UTF-8
		if r >= 32 && r <= 126 || r == '\n' || r == '\t' {
			sanitized.WriteRune(r)
		}
		// Move to next rune
		data = data[size:]
	}
	return []byte(sanitized.String())
}
