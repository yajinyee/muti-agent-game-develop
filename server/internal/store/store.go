// Package store 玩家狀態持久化層（DAY-026）
// 支援 Redis 持久化和記憶體降級兩種模式
// 設計原則：Redis 不可用時自動降級到記憶體模式，不中斷服務
package store

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// PlayerState 玩家持久化狀態
type PlayerState struct {
	PlayerID    string    `json:"player_id"`
	DisplayName string    `json:"display_name"`
	Coins       int64     `json:"coins"`
	Labor       int       `json:"labor"`
	BetLevel    int       `json:"bet_level"`
	SessionScore int64    `json:"session_score"`
	MaxCoins    int64     `json:"max_coins"`
	KillCount   int       `json:"kill_count"`
	RoomID      string    `json:"room_id"`
	LastSeen    time.Time `json:"last_seen"`
}

// Store 玩家狀態儲存介面
// 兩種實作：RedisStore（有 Redis）和 MemoryStore（降級模式）
type Store interface {
	// SavePlayer 儲存玩家狀態（玩家離開或定期儲存時呼叫）
	SavePlayer(state *PlayerState) error

	// LoadPlayer 讀取玩家狀態（玩家加入時呼叫）
	// 找不到時回傳 nil, nil（不是錯誤）
	LoadPlayer(playerID string) (*PlayerState, error)

	// DeletePlayer 刪除玩家狀態（玩家帳號刪除時）
	DeletePlayer(playerID string) error

	// GetTopPlayers 取得排行榜前 N 名
	GetTopPlayers(n int) ([]*PlayerState, error)

	// UpdateLeaderboard 更新排行榜分數
	UpdateLeaderboard(playerID string, score int64) error

	// Close 關閉連線
	Close() error

	// IsRedis 是否使用 Redis（用於日誌顯示）
	IsRedis() bool
}

// New 建立 Store 實例，自動選擇 Redis 或記憶體模式
// redisURL 為空或連線失敗時，自動降級到記憶體模式
func New(redisURL string) Store {
	if redisURL == "" {
		log.Println("[Store] REDIS_URL not set, using in-memory store (player state will not persist)")
		return NewMemoryStore()
	}

	store, err := NewRedisStore(redisURL)
	if err != nil {
		log.Printf("[Store] Redis connection failed: %v", err)
		log.Println("[Store] Falling back to in-memory store (player state will not persist)")
		return NewMemoryStore()
	}

	log.Printf("[Store] Connected to Redis: %s", redisURL)
	return store
}

// MemoryStore 記憶體模式（降級用，Server 重啟後狀態丟失）
type MemoryStore struct {
	players     map[string]*PlayerState
	leaderboard map[string]int64 // playerID → score
	mu          sync.RWMutex
}

// NewMemoryStore 建立記憶體 Store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		players:     make(map[string]*PlayerState),
		leaderboard: make(map[string]int64),
	}
}

func (m *MemoryStore) SavePlayer(state *PlayerState) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	state.LastSeen = time.Now()
	// 深拷貝避免外部修改
	copy := *state
	m.players[state.PlayerID] = &copy
	return nil
}

func (m *MemoryStore) LoadPlayer(playerID string) (*PlayerState, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	state, ok := m.players[playerID]
	if !ok {
		return nil, nil // 找不到不是錯誤
	}
	copy := *state
	return &copy, nil
}

func (m *MemoryStore) DeletePlayer(playerID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.players, playerID)
	return nil
}

func (m *MemoryStore) GetTopPlayers(n int) ([]*PlayerState, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 簡單排序：找前 N 名
	type entry struct {
		playerID string
		score    int64
	}
	entries := make([]entry, 0, len(m.leaderboard))
	for pid, score := range m.leaderboard {
		entries = append(entries, entry{pid, score})
	}
	// 排序（降序）
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].score > entries[i].score {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	result := make([]*PlayerState, 0, n)
	for i, e := range entries {
		if i >= n {
			break
		}
		if state, ok := m.players[e.playerID]; ok {
			copy := *state
			result = append(result, &copy)
		}
	}
	return result, nil
}

func (m *MemoryStore) UpdateLeaderboard(playerID string, score int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if score > m.leaderboard[playerID] {
		m.leaderboard[playerID] = score
	}
	return nil
}

func (m *MemoryStore) Close() error {
	return nil
}

func (m *MemoryStore) IsRedis() bool {
	return false
}

// RedisStore Redis 模式（生產環境用）
// 注意：需要安裝 github.com/redis/go-redis/v9
// 目前為骨架實作，待 Phase 2 完整實作
type RedisStore struct {
	url string
	// client *redis.Client  ← Phase 2 取消注釋
}

// NewRedisStore 建立 Redis Store
// Phase 2 實作：連線 Redis，驗證連線
func NewRedisStore(url string) (*RedisStore, error) {
	// TODO Phase 2：
	// opt, err := redis.ParseURL(url)
	// if err != nil {
	//     return nil, fmt.Errorf("invalid redis URL: %w", err)
	// }
	// client := redis.NewClient(opt)
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := client.Ping(ctx).Err(); err != nil {
	//     return nil, fmt.Errorf("redis ping failed: %w", err)
	// }
	// return &RedisStore{url: url, client: client}, nil

	// 目前骨架：直接回傳錯誤，讓呼叫者降級到 MemoryStore
	return nil, fmt.Errorf("redis store not yet implemented (Phase 2), use REDIS_URL='' to use memory store")
}

func (r *RedisStore) SavePlayer(state *PlayerState) error {
	// TODO Phase 2：HMSET player:{id} ...
	return fmt.Errorf("not implemented")
}

func (r *RedisStore) LoadPlayer(playerID string) (*PlayerState, error) {
	// TODO Phase 2：HGETALL player:{id}
	return nil, fmt.Errorf("not implemented")
}

func (r *RedisStore) DeletePlayer(playerID string) error {
	// TODO Phase 2：DEL player:{id}
	return fmt.Errorf("not implemented")
}

func (r *RedisStore) GetTopPlayers(n int) ([]*PlayerState, error) {
	// TODO Phase 2：ZREVRANGE leaderboard:daily:{date} 0 n-1 WITHSCORES
	return nil, fmt.Errorf("not implemented")
}

func (r *RedisStore) UpdateLeaderboard(playerID string, score int64) error {
	// TODO Phase 2：ZADD leaderboard:daily:{date} score playerID
	return fmt.Errorf("not implemented")
}

func (r *RedisStore) Close() error {
	// TODO Phase 2：client.Close()
	return nil
}

func (r *RedisStore) IsRedis() bool {
	return true
}
