// Package ws WebSocket 連線管理
package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024 // 512KB

	// Rate limiting：每個客戶端每秒最多 30 條訊息（token bucket）
	// 捕魚機正常操作：攻擊 ~10/s，投注切換 ~2/s，ping ~1/s → 30/s 足夠
	rateLimitPerSecond = 30
	rateLimitBurst     = 60 // 允許短暫爆發（最多 2 秒的量）
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	EnableCompression: true, // permessage-deflate，減少頻寬（skill-go-websocket-scalability）
	CheckOrigin: func(r *http.Request) bool {
		return true // Prototype 允許所有來源
	},
}

// ClientRole 客戶端角色
type ClientRole string

const (
	RolePlayer    ClientRole = "player"    // 一般玩家（可攻擊、投注）
	RoleSpectator ClientRole = "spectator" // 觀戰者（只讀，收廣播但不能發遊戲指令）
)

// rateLimiter Token bucket 速率限制器（per-client）
type rateLimiter struct {
	tokens    float64
	maxTokens float64
	refillRate float64 // tokens per second
	lastRefill time.Time
	mu         sync.Mutex
}

// newRateLimiter 建立速率限制器
func newRateLimiter(perSecond, burst int) *rateLimiter {
	return &rateLimiter{
		tokens:     float64(burst),
		maxTokens:  float64(burst),
		refillRate: float64(perSecond),
		lastRefill: time.Now(),
	}
}

// Allow 嘗試消耗一個 token，回傳是否允許
func (r *rateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastRefill).Seconds()
	r.lastRefill = now

	// 補充 token
	r.tokens += elapsed * r.refillRate
	if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}

	if r.tokens >= 1.0 {
		r.tokens -= 1.0
		return true
	}
	return false
}

// Client WebSocket 客戶端
type Client struct {
	ID       string
	Hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	PlayerID string
	Role     ClientRole   // 客戶端角色（DAY-023）
	limiter  *rateLimiter // per-client 速率限制（DAY-036）
}

// Hub 管理所有 WebSocket 連線
type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client

	// 訊息處理回調
	OnMessage    func(clientID string, msg *Message)
	OnConnect    func(clientID string)
	OnDisconnect func(clientID string)

	// 訊息吞吐量計數器（原子操作，供 /metrics 使用）
	MsgReceived atomic.Int64 // 收到的訊息總數
	MsgSent     atomic.Int64 // 發送的訊息總數
	MsgDropped  atomic.Int64 // 丟棄的訊息總數（buffer full + rate limit）
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

	log.Printf("[WS] Client connected: %s (role=%s)", client.ID, client.Role)
	if h.OnConnect != nil && client.Role == RolePlayer {
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

	log.Printf("[WS] Client disconnected: %s (role=%s)", client.ID, client.Role)
	if h.OnDisconnect != nil && client.Role == RolePlayer {
		h.OnDisconnect(client.ID)
	}
}

// SpectatorCount 目前觀戰者數量
func (h *Hub) SpectatorCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	count := 0
	for _, c := range h.clients {
		if c.Role == RoleSpectator {
			count++
		}
	}
	return count
}

// PlayerCount 目前玩家數量（不含觀戰者）
func (h *Hub) PlayerCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	count := 0
	for _, c := range h.clients {
		if c.Role == RolePlayer {
			count++
		}
	}
	return count
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
		h.MsgSent.Add(1)
	default:
		h.MsgDropped.Add(1)
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
			h.MsgSent.Add(1)
		default:
			h.MsgDropped.Add(1)
			log.Printf("[WS] Broadcast buffer full for client %s", client.ID)
		}
	}
}

// ClientCount 目前總連線數（玩家 + 觀戰者）
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// ServeWS 處理 WebSocket 升級請求（一般玩家）
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request, clientID string) {
	h.serveWSWithRole(w, r, clientID, RolePlayer)
}

// ServeSpectatorWS 處理觀戰者 WebSocket 升級請求（DAY-023）
func (h *Hub) ServeSpectatorWS(w http.ResponseWriter, r *http.Request, clientID string) {
	h.serveWSWithRole(w, r, clientID, RoleSpectator)
}

// serveWSWithRole 內部：依角色建立 WebSocket 連線
func (h *Hub) serveWSWithRole(w http.ResponseWriter, r *http.Request, clientID string, role ClientRole) {
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
		Role:     role,
		limiter:  newRateLimiter(rateLimitPerSecond, rateLimitBurst),
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

		// 觀戰者只允許 ping，其他遊戲指令一律忽略（DAY-023）
		if c.Role == RoleSpectator && msg.Type != MsgPing {
			log.Printf("[WS] Spectator %s tried to send %s, ignored", c.ID, msg.Type)
			continue
		}

		// Rate limiting：超過速率限制時丟棄訊息（DAY-036）
		// ping 訊息豁免（避免影響心跳機制）
		if msg.Type != MsgPing && !c.limiter.Allow() {
			log.Printf("[WS] Rate limit exceeded for client %s, dropping %s", c.ID, msg.Type)
			h.MsgDropped.Add(1)
			continue
		}

		h.MsgReceived.Add(1)
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
