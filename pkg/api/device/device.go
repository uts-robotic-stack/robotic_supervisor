package device

import (
	"encoding/json"
	"net/http"

	"github.com/containrrr/watchtower/pkg/types"
	log "github.com/sirupsen/logrus"
)

var (
	lock chan bool
)

// New is a factory function creating a new  Handler instance
func New(
	getFunc func() types.Device,
	updateLock chan bool) *Handler {
	if updateLock != nil {
		lock = updateLock
	} else {
		lock = make(chan bool, 1)
		lock <- true
	}
	return &Handler{
		GetFunc: getFunc,
		Path:    "/api/v1/watchtower/device-info",
	}
}

// Handler is an API handler used for triggering container update scans
type Handler struct {
	GetFunc func() types.Device
	Path    string
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
		log.Info("Received HTTP request to get device-info")
		output := handle.GetFunc()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)

	default:
		log.Info("Skipped. Another docker process is already running.")
		http.Error(w, "Request dropped. Another docker process is already running.", http.StatusConflict)
	}
}
