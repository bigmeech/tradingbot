package clients

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketClient manages a WebSocket connection with automatic reconnection, ping/pong handling, and rate limiting.
type WebSocketClient struct {
	conn               *websocket.Conn
	url                string
	connectionLifetime time.Duration
	pingInterval       time.Duration
	pongTimeout        time.Duration
	rateLimit          int
	messageTicker      *time.Ticker
	maxStreams         int
	activeStreams      int
	mu                 sync.Mutex
	stopCh             chan struct{}
	reconnectCh        chan struct{}
}

// NewWebSocketClient initializes a new WebSocketClient with configuration options.
func NewWebSocketClient(url string, connectionLifetime, pingInterval, pongTimeout time.Duration, rateLimit, maxStreams int) *WebSocketClient {
	client := &WebSocketClient{
		url:                url,
		connectionLifetime: connectionLifetime,
		pingInterval:       pingInterval,
		pongTimeout:        pongTimeout,
		rateLimit:          rateLimit,
		messageTicker:      time.NewTicker(time.Second / time.Duration(rateLimit)),
		maxStreams:         maxStreams,
		stopCh:             make(chan struct{}),
		reconnectCh:        make(chan struct{}),
	}
	go client.manageConnection()
	return client
}

// Connect opens a WebSocket connection and sets up the ping/pong handlers.
func (c *WebSocketClient) Connect() error {
	var err error
	c.conn, _, err = websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	// Set up WebSocket ping/pong handlers
	c.conn.SetPongHandler(func(string) error {
		return c.conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(c.pongTimeout))
	})
	return nil
}

func (c *WebSocketClient) GetConnectionUrl() string {
	return c.url
}

// manageConnection handles automatic reconnection, ping/pong, and connection lifetime limits.
func (c *WebSocketClient) manageConnection() {
	go c.pingHandler()
	go c.reconnectionHandler()
}

// pingHandler sends periodic pings to keep the connection alive.
func (c *WebSocketClient) pingHandler() {
	ticker := time.NewTicker(c.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Ping failed, reconnecting:", err)
				c.reconnect()
			}
			c.mu.Unlock()
		case <-c.stopCh:
			return
		}
	}
}

// reconnectionHandler handles reconnection when the connection lifetime expires.
func (c *WebSocketClient) reconnectionHandler() {
	ticker := time.NewTicker(c.connectionLifetime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Reconnecting WebSocket due to connection lifetime expiry")
			c.reconnect()
		case <-c.stopCh:
			return
		}
	}
}

// reconnect safely reconnects the WebSocket connection.
func (c *WebSocketClient) reconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	_ = c.conn.Close()
	if err := c.Connect(); err != nil {
		log.Println("Failed to reconnect:", err)
	} else {
		log.Println("Reconnected successfully")
	}
}

// StartStreaming manages stream subscriptions, enforcing a limit on the maximum number of streams.
func (c *WebSocketClient) StartStreaming(handler func(string, []byte)) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.activeStreams >= c.maxStreams {
		return fmt.Errorf("max streams limit reached")
	}

	c.activeStreams++
	go func() {
		for {
			select {
			case <-c.stopCh:
				return
			case <-c.messageTicker.C:
				_, message, err := c.conn.ReadMessage()
				if err != nil {
					log.Println("Error reading WebSocket message:", err)
					c.reconnect()
					continue
				}
				handler(c.url, message)
			}
		}
	}()
	return nil
}

// StopStreaming decreases the active stream count and closes the connection if no streams remain.
func (c *WebSocketClient) StopStreaming() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.activeStreams > 0 {
		c.activeStreams--
		if c.activeStreams == 0 {
			close(c.stopCh)
			return c.conn.Close()
		}
	}
	return nil
}
