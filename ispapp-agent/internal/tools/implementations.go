package tools

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

// DefaultUbusCaller implements the UbusCaller interface
type DefaultUbusCaller struct {
	socket    string
	timeout   time.Duration
	connected bool
}

// NewDefaultUbusCaller creates a new UBus caller with the given socket path and timeout
func NewDefaultUbusCaller(socket string, timeoutSeconds int) *DefaultUbusCaller {
	return &DefaultUbusCaller{
		socket:    socket,
		timeout:   time.Duration(timeoutSeconds) * time.Second,
		connected: true,
	}
}

// Call invokes an UBus method and returns the result
func (u *DefaultUbusCaller) Call(service, method string, args map[string]interface{}) (map[string]interface{}, error) {
	if !u.IsConnected() {
		return nil, errors.New("not connected to UBus")
	}

	// This is a simplified implementation that would typically use ubus_cli or a native Go UBus client
	cmd := exec.Command("ubus", "call", service, method)

	// Convert args to JSON or command line parameters
	if len(args) > 0 {
		argStrs := []string{}
		for k, v := range args {
			argStrs = append(argStrs, fmt.Sprintf("%s=%v", k, v))
		}
		cmd.Args = append(cmd.Args, strings.Join(argStrs, " "))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ubus call failed: %w", err)
	}

	// Parse the output into a map (simplified)
	result := make(map[string]interface{})
	result["output"] = string(output)

	return result, nil
}

// CallList invokes the UBus list method for a service
func (u *DefaultUbusCaller) CallList(service string) (map[string]interface{}, error) {
	if !u.IsConnected() {
		return nil, errors.New("not connected to UBus")
	}

	cmd := exec.Command("ubus", "list", service)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ubus list failed: %w", err)
	}

	// Parse the output into a map (simplified)
	result := make(map[string]interface{})
	result["methods"] = strings.Split(strings.TrimSpace(string(output)), "\n")

	return result, nil
}

// IsConnected returns true if the UBus connection is active
func (u *DefaultUbusCaller) IsConnected() bool {
	return u.connected
}

// Reconnect attempts to reconnect to UBus
func (u *DefaultUbusCaller) Reconnect() error {
	// Simulate reconnection
	u.connected = true
	return nil
}

// DefaultCommandRunner implements the CommandRunner interface
type DefaultCommandRunner struct {
	workDir string
}

// NewDefaultCommandRunner creates a new command runner with optional working directory
func NewDefaultCommandRunner(workDir string) *DefaultCommandRunner {
	return &DefaultCommandRunner{
		workDir: workDir,
	}
}

// Run executes a command and returns its output
func (c *DefaultCommandRunner) Run(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)

	if c.workDir != "" {
		cmd.Dir = c.workDir
	}

	return cmd.CombinedOutput()
}

// RunWithTimeout executes a command with a timeout
func (c *DefaultCommandRunner) RunWithTimeout(timeout int, command string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)

	if c.workDir != "" {
		cmd.Dir = c.workDir
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return out.Bytes(), fmt.Errorf("command timed out after %d seconds", timeout)
		}
		return out.Bytes(), err
	}

	return out.Bytes(), nil
}

// DefaultFileReader implements the FileReader interface
type DefaultFileReader struct{}

// NewDefaultFileReader creates a new file reader
func NewDefaultFileReader() *DefaultFileReader {
	return &DefaultFileReader{}
}

// ReadFile reads the entire contents of a file
func (f *DefaultFileReader) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// ReadFileLines reads a file and returns it as lines
func (f *DefaultFileReader) ReadFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// FileExists checks if a file exists
func (f *DefaultFileReader) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ReadDirNames reads the names of all files in a directory
func (f *DefaultFileReader) ReadDirNames(path string) ([]string, error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	return names, nil
}
