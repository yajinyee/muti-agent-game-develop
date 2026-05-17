// Package ws WebSocket 連線管理
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
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024 // 512KB
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	EnableCompression: true, // permessage-deflate，減少頻寬（skill-go-websocket-scalability）
	CheckOrigin: func(r *http.Request) bool {
		return true // Prototype 允許所有來源
	},
}

// Client WebSocket 客戶端
type Client struct {
	ID       string
	Hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	PlayerID string
}

// Hub 管理所有 WebSocket 連線
type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client

	// 訊息處理回調
	OnMessage func(clientID string, msg *Message)
	OnConnect func(clientID string)
	OnDisconnect func(clientID string)
}

// NewHub 建立新 Hub
func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

// Register 註冊客戶端
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	h.clients[client.ID] = client
	h.mu.Unlock()

	log.Printf("[WS] Client connected: %s", client.ID)
	if h.OnConnect != nil {
		h.OnConnect(client.ID)
	}
}

// Unregister 移除客戶端
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	if _, ok := h.clients[client.ID]; ok {
		delete(h.clients, client.ID)
		close(client.send)
	}
	h.mu.Unlock()

	log.Printf("[WS] Client disconnected: %s", client.ID)
	if h.OnDisconnect != nil {
		h.OnDisconnect(client.ID)
	}
}

// Send 傳送訊息給指定客戶端
func (h *Hub) Send(clientID string, msg *Message) error {
	h.mu.RLock()
	client, ok := h.clients[clientID]
	h.mu.RUnlock()

	if !ok {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case client.send <- data:
	default:
		log.Printf("[WS] Send buffer full for client %s, dropping message", clientID)
	}
	return nil
}

// Broadcast 廣播訊息給所有客戶端
func (h *Hub) Broadcast(msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[WS] Broadcast marshal error: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		select {
		case client.send <- data:
		default:
			log.Printf("[WS] Broadcast buffer full for client %s", client.ID)
		}
	}
}

// ClientCount 目前連線數
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// ServeWS 處理 WebSocket 升級請求
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request, clientID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:       clientID,
		Hub:      h,
		conn:     conn,
		send:     make(chan []byte, 256),
		PlayerID: clientID,
	}

	h.Register(client)

	go client.writePump()
	go client.readPump(h)
}

// readPump 讀取客戶端訊息
func (c *Client) readPump(h *Hub) {
	defer func() {
		h.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WS] Read error for %s: %v", c.ID, err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Printf("[WS] Parse error for %s: %v", c.ID, err)
			continue
		}

		if h.OnMessage != nil {
			h.OnMessage(c.ID, &msg)
		}
	}
}

// writePump 傳送訊息給客戶端
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 每個訊息獨立傳送，避免 JSON 合併造成解析問題
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}

			// 批次傳送佇列中的訊息（各自獨立 frame）
			n := len(c.send)
			for i := 0; i < n; i++ {
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				w.Write(<-c.send)
				if err := w.Close(); err != nil {
					return
				}
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
