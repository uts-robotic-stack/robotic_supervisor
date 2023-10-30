package container

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var (
	lock chan bool
)

// New is a factory function creating a new  Handler instance
func New(handleFunc func(map[string]interface{}), updateLock chan bool) *Handler {
	if updateLock != nil {
		lock = updateLock
	} else {
		lock = make(chan bool, 1)
		lock <- true
	}

	return &Handler{
		fn:   handleFunc,
		Path: "/watchtower/v1/container",
	}
}

// Handler is an API handler used for triggering container update scans
type Handler struct {
	fn   func(map[string]interface{})
	Path string
}

// Handle is the actual http.Handle function doing all the heavy lifting
func (handle *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	log.Info("Received HTTP request to start/stop container")
	w.Header().Set("Content-Type", "application/json")

	var reqBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	select {
	case chanValue := <-lock:
		defer func() {
			lock <- chanValue
		}()
		handle.fn(reqBody)
	default:
		log.Info("Skipped. Another update already running.")
	}

}
