package tools

import (
	"net/http"
	"os"
)

// UbusCaller provides access to UBus services
type UbusCaller interface {
	// Call invokes an UBus method and returns the result
	Call(service, method string, args map[string]interface{}) (map[string]interface{}, error)

	// CallList invokes the UBus list method for a service
	CallList(service string) (map[string]interface{}, error)

	// IsConnected returns true if the UBus connection is active
	IsConnected() bool

	// Reconnect attempts to reconnect to UBus
	Reconnect() error
}

// CommandRunner executes system commands
type CommandRunner interface {
	// Run executes a command and returns its output
	Run(command string, args ...string) ([]byte, error)

	// RunWithTimeout executes a command with a timeout
	RunWithTimeout(timeout int, command string, args ...string) ([]byte, error)
}

// FileReader reads from the filesystem
type FileReader interface {
	// ReadFile reads the entire contents of a file
	ReadFile(path string) ([]byte, error)

	// ReadFileLines reads a file and returns it as lines
	ReadFileLines(path string) ([]string, error)

	// FileExists checks if a file exists
	FileExists(path string) bool

	// ReadDirNames reads the names of all files in a directory
	ReadDirNames(path string) ([]string, error)
}

// FileWriter writes to the filesystem
type FileWriter interface {
	// WriteFile writes data to a file, creating it if it doesn't exist
	WriteFile(path string, data []byte, perm os.FileMode) error

	// AppendToFile appends data to an existing file
	AppendToFile(path string, data []byte) error

	// MkdirAll creates a directory and all necessary parents
	MkdirAll(path string, perm os.FileMode) error

	// Remove deletes a file or empty directory
	Remove(path string) error

	// RemoveAll deletes a file or directory and any children
	RemoveAll(path string) error
}

// HttpClient performs HTTP operations
type HttpClient interface {
	// Get performs an HTTP GET request
	Get(url string) ([]byte, error)

	// Post performs an HTTP POST request with the given data
	Post(url string, contentType string, body []byte) ([]byte, error)

	// PostForm submits form data to the specified URL
	PostForm(url string, data map[string]string) ([]byte, error)

	// Do performs a custom HTTP request
	Do(req *http.Request) ([]byte, error)
}

// ConfigManager handles application configuration
type ConfigManager interface {
	// GetString retrieves a string configuration value
	GetString(key string) string

	// GetInt retrieves an integer configuration value
	GetInt(key string) int

	// GetBool retrieves a boolean configuration value
	GetBool(key string) bool

	// Set sets a configuration value
	Set(key string, value interface{}) error

	// Save persists configuration changes
	Save() error

	// Load reloads configuration from storage
	Load() error
}

// Logger provides logging capabilities
type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	// WithFields returns a logger with contextual fields attached
	WithFields(fields map[string]interface{}) Logger
}
