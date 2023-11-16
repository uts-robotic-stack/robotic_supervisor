package container

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	lock chan bool
)

// New is a factory function creating a new  Handler instance
func New(
	postFunc func(map[string]interface{}),
	getFunc func() []string,
	updateLock chan bool) *Handler {
	if updateLock != nil {
		lock = updateLock
	} else {
		lock = make(chan bool, 1)
		lock <- true
	}

	return &Handler{
		PostFunc: postFunc,
		GetFunc:  getFunc,
		Path:     "/api/v1/watchtower/container",
	}
}

// Handler is an API handler used for triggering container update scans
type Handler struct {
	PostFunc func(map[string]interface{})
	GetFunc  func() []string
	Path     string
}

// Handle is the actual http.Handle function doing all the heavy lifting
func (handle *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	select {
	case chanValue := <-lock:
		defer func() {
			lock <- chanValue
		}()
		if r.Method == http.MethodPost {
			log.Info("Received HTTP request to start/stop container")
			w.Header().Set("Content-Type", "application/json")

			var reqBody map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			handle.PostFunc(reqBody)
		} else if r.Method == http.MethodGet {
			log.Info("Received HTTP request to get all container")
			output := handle.GetFunc()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(output)
		}

	default:
		log.Info("Skipped. Another docker process is already running.")
		http.Error(w, "Request dropped. Another docker process is already running.", http.StatusConflict)
	}
}

// Websocket for streaming logs

type WSHandler struct {
	fn   func(string, *websocket.Conn)
	Path string
}

func NewWSHandler(handlerFunc func(string, *websocket.Conn)) *WSHandler {
	return &WSHandler{
		fn:   handlerFunc,
		Path: "/api/v1/watchtower/logs",
	}
}

var upgrader = websocket.Upgrader{}

func (handle *WSHandler) Handle(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	// Extract info
	containerName := r.URL.Query().Get("container_name")

	// Handle info
	handle.fn(containerName, conn)
}
