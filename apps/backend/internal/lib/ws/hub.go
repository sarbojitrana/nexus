package ws

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 8192
	sendBuffer     = 16
)

func NewUpgrader(allowedOrigins []string) websocket.Upgrader {
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = true
	}

	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == "" || allowed[origin]
		},
	}
}

type Event struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type Client struct {
	UserID string
	conn   *websocket.Conn
	send   chan []byte
	hub    *Hub
}

func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		UserID: userID,
		conn:   conn,
		send:   make(chan []byte, sendBuffer),
		hub:    hub,
	}
}

func (c *Client) Send(event any) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	select {
	case c.send <- data:
	default: // slow consumer, drop rather than block the hub
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadPump(handle func(raw []byte)) {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		handle(raw)
	}
}

type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]bool // userID -> connections (multiple tabs/devices)
	logger  *zerolog.Logger

	OnConnect    func(userID string)
	OnDisconnect func(userID string)
}

func NewHub(logger *zerolog.Logger) *Hub {
	return &Hub{
		clients: make(map[string]map[*Client]bool),
		logger:  logger,
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	if h.clients[c.UserID] == nil {
		h.clients[c.UserID] = make(map[*Client]bool)
	}
	wentOnline := len(h.clients[c.UserID]) == 0
	h.clients[c.UserID][c] = true
	h.mu.Unlock()

	if wentOnline && h.OnConnect != nil {
		h.OnConnect(c.UserID)
	}
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	wentOffline := false
	if conns, ok := h.clients[c.UserID]; ok {
		delete(conns, c)
		if len(conns) == 0 {
			delete(h.clients, c.UserID)
			wentOffline = true
		}
	}
	h.mu.Unlock()

	close(c.send)

	if wentOffline && h.OnDisconnect != nil {
		h.OnDisconnect(c.UserID)
	}
}

func (h *Hub) IsOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[userID]) > 0
}

func (h *Hub) SendToUser(userID string, event any) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients[userID] {
		select {
		case c.send <- data:
		default:
		}
	}
}

func (h *Hub) SendToUsers(userIDs []string, event any) {
	for _, id := range userIDs {
		h.SendToUser(id, event)
	}
}
