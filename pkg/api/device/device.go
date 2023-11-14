package device

// New is a factory function creating a new  Handler instance
func New(
	postFunc func(map[string]interface{}),
	getFunc func() []string,
	updateLock chan bool) *Handler {

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
