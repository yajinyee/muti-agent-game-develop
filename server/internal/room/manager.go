// Package room 多房間管理（DAY-019）
// 支援多個獨立遊戲房間同時運行，每個房間有獨立的遊戲狀態和玩家列表
package room

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Config 房間設定
type Config struct {
	Name        string  // 房間顯示名稱
	MinBetLevel int     // 最低投注等級（1-10）
	MaxBetLevel int     // 最高投注等級（1-10）
	MaxPlayers  int     // 最大玩家數（建議 8-16）
	Theme       string  // 主題（chiikawa/default）
	RTPTarget   float64 // 目標 RTP（0.92-0.96）
}

// Info 房間資訊（供 HTTP API 回傳）
type Info struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	PlayerCount int       `json:"player_count"`
	MaxPlayers  int       `json:"max_players"`
	MinBetLevel int       `json:"min_bet_level"`
	MaxBetLevel int       `json:"max_bet_level"`
	Theme       string    `json:"theme"`
	CreatedAt   time.Time `json:"created_at"`
	IsFull      bool      `json:"is_full"`
}

// Room 單一遊戲房間
type Room struct {
	ID        string
	Config    Config
	Players   map[string]bool // playerID → 是否在線
	CreatedAt time.Time
	mu        sync.RWMutex
}

// PlayerCount 目前玩家數
func (r *Room) PlayerCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.Players)
}

// IsFull 是否已滿
func (r *Room) IsFull() bool {
	return r.PlayerCount() >= r.Config.MaxPlayers
}

// AddPlayer 加入玩家
func (r *Room) AddPlayer(playerID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.Players) >= r.Config.MaxPlayers {
		return false
	}
	r.Players[playerID] = true
	return true
}

// RemovePlayer 移除玩家
func (r *Room) RemovePlayer(playerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Players, playerID)
}

// HasPlayer 是否包含玩家
func (r *Room) HasPlayer(playerID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Players[playerID]
}

// ToInfo 轉換為 API 回傳格式
func (r *Room) ToInfo() Info {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return Info{
		ID:          r.ID,
		Name:        r.Config.Name,
		PlayerCount: len(r.Players),
		MaxPlayers:  r.Config.MaxPlayers,
		MinBetLevel: r.Config.MinBetLevel,
		MaxBetLevel: r.Config.MaxBetLevel,
		Theme:       r.Config.Theme,
		CreatedAt:   r.CreatedAt,
		IsFull:      len(r.Players) >= r.Config.MaxPlayers,
	}
}

// Manager 多房間管理器
type Manager struct {
	rooms    map[string]*Room
	mu       sync.RWMutex
	maxRooms int
}

// NewManager 建立房間管理器，並預建預設房間
func NewManager() *Manager {
	m := &Manager{
		rooms:    make(map[string]*Room),
		maxRooms: 20,
	}

	// 預建 3 個預設房間（低/中/高投注）
	m.createRoom("room-001", Config{
		Name:        "初心者房間",
		MinBetLevel: 1,
		MaxBetLevel: 4,
		MaxPlayers:  16,
		Theme:       "chiikawa",
		RTPTarget:   0.92,
	})
	m.createRoom("room-002", Config{
		Name:        "一般房間",
		MinBetLevel: 3,
		MaxBetLevel: 7,
		MaxPlayers:  12,
		Theme:       "chiikawa",
		RTPTarget:   0.94,
	})
	m.createRoom("room-003", Config{
		Name:        "高手房間",
		MinBetLevel: 6,
		MaxBetLevel: 10,
		MaxPlayers:  8,
		Theme:       "chiikawa",
		RTPTarget:   0.96,
	})

	log.Printf("[RoomManager] Initialized with %d default rooms", len(m.rooms))
	return m
}

// createRoom 內部建立房間（不加鎖，由呼叫者負責）
func (m *Manager) createRoom(id string, cfg Config) *Room {
	room := &Room{
		ID:        id,
		Config:    cfg,
		Players:   make(map[string]bool),
		CreatedAt: time.Now(),
	}
	m.rooms[id] = room
	return room
}

// GetRoom 取得指定房間
func (m *Manager) GetRoom(roomID string) (*Room, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.rooms[roomID]
	return r, ok
}

// GetOrDefault 取得指定房間，不存在時回傳 room-001
func (m *Manager) GetOrDefault(roomID string) *Room {
	if roomID == "" {
		roomID = "room-001"
	}
	m.mu.RLock()
	r, ok := m.rooms[roomID]
	m.mu.RUnlock()
	if !ok {
		// 找不到指定房間，回傳人數最少的房間
		return m.FindLeastPopulated()
	}
	return r
}

// FindLeastPopulated 找人數最少且未滿的房間
func (m *Manager) FindLeastPopulated() *Room {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var best *Room
	for _, r := range m.rooms {
		if r.IsFull() {
			continue
		}
		if best == nil || r.PlayerCount() < best.PlayerCount() {
			best = r
		}
	}
	if best == nil {
		// 所有房間都滿了，回傳 room-001（讓 Server 決定是否拒絕）
		return m.rooms["room-001"]
	}
	return best
}

// CreateRoom 動態建立新房間
func (m *Manager) CreateRoom(cfg Config) (*Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.rooms) >= m.maxRooms {
		return nil, fmt.Errorf("max rooms (%d) reached", m.maxRooms)
	}

	id := fmt.Sprintf("room-%03d", len(m.rooms)+1)
	// 確保 ID 不重複
	for {
		if _, exists := m.rooms[id]; !exists {
			break
		}
		id = fmt.Sprintf("room-%s", time.Now().Format("150405"))
	}

	room := m.createRoom(id, cfg)
	log.Printf("[RoomManager] Created room %s: %s", id, cfg.Name)
	return room, nil
}

// DeleteRoom 刪除房間（需要先踢出所有玩家）
func (m *Manager) DeleteRoom(roomID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	r, ok := m.rooms[roomID]
	if !ok {
		return fmt.Errorf("room %s not found", roomID)
	}
	if r.PlayerCount() > 0 {
		return fmt.Errorf("room %s still has %d players", roomID, r.PlayerCount())
	}

	delete(m.rooms, roomID)
	log.Printf("[RoomManager] Deleted room %s", roomID)
	return nil
}

// ListRooms 列出所有房間資訊
func (m *Manager) ListRooms() []Info {
	m.mu.RLock()
	defer m.mu.RUnlock()

	infos := make([]Info, 0, len(m.rooms))
	for _, r := range m.rooms {
		infos = append(infos, r.ToInfo())
	}
	return infos
}

// JoinRoom 玩家加入房間
func (m *Manager) JoinRoom(roomID, playerID string) (*Room, error) {
	room, ok := m.GetRoom(roomID)
	if !ok {
		return nil, fmt.Errorf("room %s not found", roomID)
	}
	if !room.AddPlayer(playerID) {
		return nil, fmt.Errorf("room %s is full (%d/%d)",
			roomID, room.PlayerCount(), room.Config.MaxPlayers)
	}
	log.Printf("[RoomManager] Player %s joined room %s (%d/%d)",
		playerID, roomID, room.PlayerCount(), room.Config.MaxPlayers)
	return room, nil
}

// LeaveRoom 玩家離開房間
func (m *Manager) LeaveRoom(roomID, playerID string) {
	room, ok := m.GetRoom(roomID)
	if !ok {
		return
	}
	room.RemovePlayer(playerID)
	log.Printf("[RoomManager] Player %s left room %s (%d/%d)",
		playerID, roomID, room.PlayerCount(), room.Config.MaxPlayers)
}

// FindPlayerRoom 找玩家所在的房間
func (m *Manager) FindPlayerRoom(playerID string) *Room {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, r := range m.rooms {
		if r.HasPlayer(playerID) {
			return r
		}
	}
	return nil
}
