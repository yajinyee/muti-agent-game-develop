// Package ws — WebSocket Hub
// server-infra-agent 負責維護
package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 50 * time.Second
	maxMessageSize = 65536
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client 代表一個 WebSocket 連線
type Client struct {
	ID   string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// Hub 管理所有連線
type Hub struct {
	mu         sync.RWMutex
	clients    map[string]*Client
	OnConnect    func(clientID string)
	OnDisconnect func(clientID string)
	OnMessage    func(clientID string, msgType string, payload json.RawMessage)
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

// ServeWS 升級 HTTP 連線為 WebSocket
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request, clientID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[Hub] upgrade error: %v", err)
		return
	}
	c := &Client{ID: clientID, hub: h, conn: conn, send: make(chan []byte, 256)}
	h.mu.Lock()
	h.clients[clientID] = c
	h.mu.Unlock()

	if h.OnConnect != nil {
		h.OnConnect(clientID)
	}
	go c.writePump()
	c.readPump()
}

// Send 傳送訊息給指定玩家
func (h *Hub) Send(clientID string, msgType string, payload interface{}) {
	h.mu.RLock()
	c, ok := h.clients[clientID]
	h.mu.RUnlock()
	if !ok {
		return
	}
	data, err := json.Marshal(map[string]interface{}{"type": msgType, "payload": payload})
	if err != nil {
		return
	}
	select {
	case c.send <- data:
	default:
		log.Printf("[Hub] send buffer full for %s", clientID)
	}
}

// Broadcast 廣播給所有玩家
func (h *Hub) Broadcast(msgType string, payload interface{}) {
	data, err := json.Marshal(map[string]interface{}{"type": msgType, "payload": payload})
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, c := range h.clients {
		select {
		case c.send <- data:
		default:
		}
	}
}

// PlayerCount 回傳目前連線數
func (h *Hub) PlayerCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// PlayerIDs 回傳所有玩家 ID
func (h *Hub) PlayerIDs() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	ids := make([]string, 0, len(h.clients))
	for id := range h.clients {
		ids = append(ids, id)
	}
	return ids
}

func (c *Client) readPump() {
	defer func() {
		c.hub.mu.Lock()
		delete(c.hub.clients, c.ID)
		c.hub.mu.Unlock()
		c.conn.Close()
		if c.hub.OnDisconnect != nil {
			c.hub.OnDisconnect(c.ID)
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var env struct {
			Type    string          `json:"type"`
			Payload json.RawMessage `json:"payload"`
		}
		if err := json.Unmarshal(msg, &env); err != nil {
			continue
		}
		if c.hub.OnMessage != nil {
			c.hub.OnMessage(c.ID, env.Type, env.Payload)
		}
	}
}

func (c *Client) writePump() {
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
