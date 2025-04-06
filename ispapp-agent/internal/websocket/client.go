package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/lxzan/gws"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024
)

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// Client is a WebSocket client
type Client struct {
	url           string
	conn          *gws.Conn
	log           *logrus.Logger
	send          chan Message
	handlers      map[string]MessageHandler
	reconnectWait time.Duration
	authToken     string
	deviceID      string
	mutex         sync.Mutex
	ctx           context.Context
	cancel        context.CancelFunc
}

// MessageHandler defines a function that handles specific message types
type MessageHandler func(message Message) error

// NewClient creates a new WebSocket client
func NewClient(url string, log *logrus.Logger, deviceID, authToken string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		url:           url,
		log:           log,
		send:          make(chan Message, 256),
		handlers:      make(map[string]MessageHandler),
		reconnectWait: 5 * time.Second,
		authToken:     authToken,
		deviceID:      deviceID,
		ctx:           ctx,
		cancel:        cancel,
	}
}

// RegisterHandler registers a handler for a specific message type
func (c *Client) RegisterHandler(messageType string, handler MessageHandler) {
	c.handlers[messageType] = handler
}

// ClientHandler implements the gws event handler interface
type ClientHandler struct {
	client *Client
}

// OnClose handles connection close events
func (h *ClientHandler) OnClose(socket *gws.Conn, err error) {
	h.client.log.Infof("WebSocket connection closed: %v", err)
	h.client.mutex.Lock()
	h.client.conn = nil
	h.client.mutex.Unlock()
}

// OnPong handles pong responses from the server
func (h *ClientHandler) OnPong(socket *gws.Conn, payload []byte) {
	// Just log for debugging if needed
	h.client.log.Debugf("Received pong from server")
}

// OnOpen handles successful connection establishment
func (h *ClientHandler) OnOpen(socket *gws.Conn) {
	h.client.log.Info("WebSocket connection established")
}

// OnPing handles ping requests from the server
func (h *ClientHandler) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.WritePong(payload)
}

// OnMessage handles incoming messages
func (h *ClientHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()

	// Parse the message
	var msg Message
	if err := json.Unmarshal(message.Data.Bytes(), &msg); err != nil {
		h.client.log.Errorf("Failed to unmarshal message: %v", err)
		return
	}

	// Handle the message
	if handler, ok := h.client.handlers[msg.Type]; ok {
		if err := handler(msg); err != nil {
			h.client.log.Errorf("Error handling message type '%s': %v", msg.Type, err)
		}
	} else {
		h.client.log.Warnf("No handler registered for message type: %s", msg.Type)
	}
}

// Connect establishes a WebSocket connection
func (c *Client) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn != nil {
		return nil
	}

	c.log.Infof("Connecting to WebSocket server: %s", c.url)

	// Add auth token and device ID to header
	header := http.Header{}
	if c.authToken != "" {
		header.Add("Authorization", "Bearer "+c.authToken)
	}
	if c.deviceID != "" {
		header.Add("X-Device-ID", c.deviceID)
	}

	// Configure client options
	clientOptions := &gws.ClientOption{
		Addr:                c.url,
		RequestHeader:       header,
		HandshakeTimeout:    writeWait,
		WriteMaxPayloadSize: maxMessageSize,
	}

	// Create event handler
	eventHandler := &ClientHandler{client: c}

	// Connect to the WebSocket server
	conn, _, err := gws.NewClient(eventHandler, clientOptions)
	if err != nil {
		return fmt.Errorf("websocket dial error: %v", err)
	}

	c.conn = conn
	c.log.Info("WebSocket connection created")

	// Start the read loop in a goroutine
	go conn.ReadLoop()

	// Start the send goroutine
	go c.handleOutgoingMessages()

	return nil
}

// Disconnect closes the WebSocket connection
func (c *Client) Disconnect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn == nil {
		return nil
	}

	c.log.Info("Disconnecting from WebSocket server")

	// Close the connection
	c.conn.WriteClose(1000, []byte("normal closure"))
	c.conn = nil

	return nil
}

// Send sends a message to the server
func (c *Client) Send(messageType string, content interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn == nil {
		return fmt.Errorf("connection not established")
	}

	c.send <- Message{
		Type:    messageType,
		Content: content,
	}

	return nil
}

// Start begins the WebSocket client with auto-reconnection
func (c *Client) Start() error {
	// Initial connection
	if err := c.Connect(); err != nil {
		c.log.Warnf("Initial WebSocket connection failed: %v", err)
	}

	// Auto-reconnect goroutine
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				if c.conn == nil {
					c.log.Info("Attempting to reconnect...")
					if err := c.Connect(); err != nil {
						c.log.Warnf("Reconnection failed: %v", err)
						time.Sleep(c.reconnectWait)
					}
				}
				time.Sleep(time.Second)
			}
		}
	}()

	return nil
}

// Stop shuts down the WebSocket client
func (c *Client) Stop() error {
	c.cancel()
	return c.Disconnect()
}

// handleOutgoingMessages processes messages from the send channel
func (c *Client) handleOutgoingMessages() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case message, ok := <-c.send:
			if !ok {
				return
			}

			c.mutex.Lock()
			if c.conn == nil {
				c.mutex.Unlock()
				continue
			}

			data, err := json.Marshal(message)
			if err != nil {
				c.log.Errorf("Failed to marshal message: %v", err)
				c.mutex.Unlock()
				continue
			}

			err = c.conn.WriteClose(1000, data)
			if err != nil {
				c.log.Errorf("Failed to write message: %v", err)
				c.conn = nil // Mark connection as closed
			}
			c.mutex.Unlock()
		}
	}
}
