package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lxzan/gws"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second
)

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// MessageHandler is a function that processes incoming messages
type MessageHandler func(msg Message) error

// Client handles WebSocket communication
type Client struct {
	url           string
	deviceID      string
	log           *logrus.Logger
	conn          *websocket.Conn
	handlers      map[string]MessageHandler
	handlersMu    sync.RWMutex
	sendMu        sync.Mutex
	closeChan     chan struct{}
	closeOnce     sync.Once
	send          chan Message
	reconnectWait time.Duration
	authToken     string
	mutex         sync.Mutex
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewClient creates a new WebSocket client
func NewClient(url string, deviceID string, log *logrus.Logger) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		url:           url,
		deviceID:      deviceID,
		log:           log,
		handlers:      make(map[string]MessageHandler),
		closeChan:     make(chan struct{}),
		send:          make(chan Message, 256),
		reconnectWait: 5 * time.Second,
		authToken:     "",
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Connect establishes a connection to the WebSocket server
func (c *Client) Connect() error {
	c.log.Debugf("Connecting to %s", c.url)

	// Create dialer with reasonable timeouts
	dialer := websocket.Dialer{
		Proxy:            websocket.DefaultDialer.Proxy,
		HandshakeTimeout: 10 * time.Second,
	}

	// Add auth headers if needed
	headers := http.Header{
		"X-Device-ID": {c.deviceID},
	}

	conn, _, err := dialer.Dial(c.url, headers)
	if err != nil {
		return fmt.Errorf("websocket dial error: %w", err)
	}

	c.conn = conn

	// Start reader in a goroutine
	go c.readPump()

	// Start outgoing message handler
	go c.handleOutgoingMessages()

	return nil
}

// OnMessage registers a handler for a specific message type
func (c *Client) OnMessage(msgType string, handler MessageHandler) {
	c.handlersMu.Lock()
	defer c.handlersMu.Unlock()

	c.handlers[msgType] = handler
}

// SendJSON sends a JSON message with the specified type and data
func (c *Client) SendJSON(msgType string, data interface{}) error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	message := map[string]interface{}{
		"type": msgType,
		"data": data,
	}

	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	return c.conn.WriteJSON(message)
}

// Ping sends a ping message to the server
func (c *Client) Ping() error {
	return c.SendJSON("ping", map[string]interface{}{
		"timestamp": time.Now().Unix(),
	})
}

// Disconnect closes the WebSocket connection
func (c *Client) Disconnect() {
	c.mutex.Lock()
	c.closeOnce.Do(func() {
		c.log.Debug("Closing WebSocket connection")

		if c.conn != nil {
			// Send close message
			c.conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			c.conn.Close()
		}

		close(c.closeChan)
	})
	c.log.Info("WebSocket connection closed")
	c.mutex.Unlock()
}

// readPump processes incoming messages
func (c *Client) readPump() {
	defer c.Disconnect()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
			) {
				c.log.Errorf("WebSocket read error: %v", err)
			}
			break
		}

		// Process message
		c.handleMessage(message)
	}
}

// handleMessage processes incoming WebSocket messages
func (c *Client) handleMessage(message []byte) {
	var msg struct {
		Type string                 `json:"type"`
		Data map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Errorf("Failed to parse message: %v", err)
		return
	}

	c.log.Debugf("Received message type: %s", msg.Type)

	c.handlersMu.RLock()
	handler, exists := c.handlers[msg.Type]
	c.handlersMu.RUnlock()

	if exists {
		go func() {
			message := Message{Type: msg.Type, Content: msg.Data}
			if err := handler(message); err != nil {
				c.log.Errorf("Error handling message type '%s': %v", msg.Type, err)
			}
		}() // Process message in a goroutine
	} else {
		c.log.Warnf("No handler for message type: %s", msg.Type)
	}
}

// RegisterHandler registers a handler for a specific message type
func (c *Client) RegisterHandler(messageType string, handler MessageHandler) {
	c.handlersMu.Lock()
	defer c.handlersMu.Unlock()
	c.handlers[messageType] = handler
}

// ClientHandler implements the gws event handler interface
type ClientHandler struct {
	client *Client
}

// OnClose handles connection close events
func (h *ClientHandler) OnClose(socket *gws.Conn, err error) {
	h.client.log.Infof("WebSocket connection closed: %v", err.Error())
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
func (c *Client) Stop() {
	c.cancel()
	c.Disconnect()
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

			err = c.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				c.log.Errorf("Failed to write message: %v", err)
				c.conn = nil // Mark connection as closed
			}
			c.mutex.Unlock()
		}
	}
}

// Authenticate authenticates the client with the WebSocket server
func (c *Client) Authenticate(login, key string) error {
	c.log.Infof("Authenticating with login: %s", login)
	authPayload := map[string]interface{}{
		"action": "auth",
		"payload": map[string]string{
			"login": login,
			"key":   key,
		},
	}

	err := c.Send("auth", authPayload)
	if err != nil {
		c.log.Errorf("Authentication failed: %v", err)
	}
	return err
}

// RefreshToken refreshes the authentication token
func (c *Client) RefreshToken(refreshToken string) error {
	refreshPayload := map[string]interface{}{
		"action": "refreshToken",
		"payload": map[string]string{
			"refreshToken": refreshToken,
		},
	}

	return c.Send("refreshToken", refreshPayload)
}

// ReportEvent sends an event to the WebSocket server
func (c *Client) ReportEvent(eventType string, data interface{}) error {
	eventPayload := map[string]interface{}{
		"action":  eventType,
		"payload": data,
	}

	return c.Send(eventType, eventPayload)
}

// SelfHeal attempts to reconnect and reauthenticate if the connection is lost
func (c *Client) SelfHeal() {
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				if c.conn == nil {
					c.log.Warn("Connection lost. Attempting to reconnect...")
					if err := c.Connect(); err != nil {
						c.log.Errorf("Reconnection failed: %v", err)
						time.Sleep(c.reconnectWait)
						continue
					}

					// Reauthenticate after reconnecting
					if err := c.Authenticate(c.deviceID, c.authToken); err != nil {
						c.log.Errorf("Reauthentication failed: %v", err)
					}
				}
				time.Sleep(time.Second)
			}
		}
	}()
}
