package handlers

// Handler defines the interface for all agent handlers
type Handler interface {
	// Name returns the handler's name
	Name() string
	
	// Start initializes and starts the handler
	Start() error
	
	// Stop shuts down the handler gracefully
	Stop() error
}

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	name string
}

// Name returns the handler name
func (h *BaseHandler) Name() string {
	return h.name
}

// NewBaseHandler creates a new base handler with the given name
func NewBaseHandler(name string) BaseHandler {
	return BaseHandler{name: name}
}
