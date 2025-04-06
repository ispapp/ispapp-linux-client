package handlers

// Handler interface for all agent components
type Handler interface {
	Name() string
	Start() error
	Stop() error
}
