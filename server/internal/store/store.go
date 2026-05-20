// Package store 玩家狀態持久化層（DAY-026/028）
// 支援 Redis 持久化和記憶體降級兩種模式
// 設計原則：Redis 不可用時自動降級到記憶體模式，不中斷服務
package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
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
	// 每日登入獎勵（DAY-065）
	LastLoginDate string `json:"last_login_date"` // "2026-05-20"（UTC+8）
	LoginStreak   int    `json:"login_streak"`    // 連續登入天數
	MaxLoginStreak int   `json:"max_login_streak"` // 歷史最高連續天數
	// 砲台外觀（DAY-071）
	EquippedSkin string   `json:"equipped_skin"` // 當前裝備外觀
	OwnedSkins   []string `json:"owned_skins"`   // 已擁有外觀列表
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

	// SetJSON 儲存任意 JSON 資料（通用 key-value，DAY-049d）
	SetJSON(key string, value interface{}, ttl time.Duration) error

	// GetJSON 讀取任意 JSON 資料（通用 key-value，DAY-049d）
	// 找不到時回傳 ErrNotFound
	GetJSON(key string, dest interface{}) error

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
	cp := *state
	m.players[state.PlayerID] = &cp
	return nil
}

func (m *MemoryStore) LoadPlayer(playerID string) (*PlayerState, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	state, ok := m.players[playerID]
	if !ok {
		return nil, nil // 找不到不是錯誤
	}
	cp := *state
	return &cp, nil
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

	type entry struct {
		playerID string
		score    int64
	}
	entries := make([]entry, 0, len(m.leaderboard))
	for pid, score := range m.leaderboard {
		entries = append(entries, entry{pid, score})
	}
	// 排序（降序）
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].score > entries[j].score
	})

	result := make([]*PlayerState, 0, n)
	for i, e := range entries {
		if i >= n {
			break
		}
		if state, ok := m.players[e.playerID]; ok {
			cp := *state
			result = append(result, &cp)
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

// SetJSON 記憶體模式：儲存 JSON 資料（DAY-049d）
func (m *MemoryStore) SetJSON(key string, value interface{}, _ time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("SetJSON marshal: %w", err)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	// 用 players map 的 key 空間儲存（前綴區分）
	// 注意：這裡借用 players map 儲存任意 JSON，key 不會和 player_id 衝突（因為有前綴）
	m.players[key] = &PlayerState{PlayerID: key, DisplayName: string(data)}
	return nil
}

// GetJSON 記憶體模式：讀取 JSON 資料（DAY-049d）
func (m *MemoryStore) GetJSON(key string, dest interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	entry, ok := m.players[key]
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}
	return json.Unmarshal([]byte(entry.DisplayName), dest)
}

// RedisStore Redis 模式（生產環境用）
// 使用 github.com/redis/go-redis/v9
// Key 設計：
//   player:{id}              → Hash（玩家狀態 JSON）
//   leaderboard:daily:{date} → Sorted Set（分數排行）
//   player:{id}:ttl          → 7 天 TTL
type RedisStore struct {
	client *redis.Client
	url    string
}

const (
	playerKeyPrefix    = "player:"
	leaderboardKeyFmt  = "leaderboard:daily:%s"
	playerTTL          = 7 * 24 * time.Hour
	leaderboardTTL     = 30 * 24 * time.Hour
	redisTimeout       = 3 * time.Second
)

// NewRedisStore 建立 Redis Store，驗證連線
func NewRedisStore(url string) (*RedisStore, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("invalid redis URL: %w", err)
	}

	client := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &RedisStore{client: client, url: url}, nil
}

func (r *RedisStore) playerKey(playerID string) string {
	return playerKeyPrefix + playerID
}

func (r *RedisStore) leaderboardKey() string {
	return fmt.Sprintf(leaderboardKeyFmt, time.Now().Format("2006-01-02"))
}

func (r *RedisStore) SavePlayer(state *PlayerState) error {
	state.LastSeen = time.Now()
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal player state: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()

	key := r.playerKey(state.PlayerID)
	pipe := r.client.Pipeline()
	pipe.Set(ctx, key, data, playerTTL)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("redis save player %s: %w", state.PlayerID, err)
	}
	return nil
}

func (r *RedisStore) LoadPlayer(playerID string) (*PlayerState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()

	data, err := r.client.Get(ctx, r.playerKey(playerID)).Bytes()
	if err == redis.Nil {
		return nil, nil // 找不到不是錯誤
	}
	if err != nil {
		return nil, fmt.Errorf("redis load player %s: %w", playerID, err)
	}

	var state PlayerState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("unmarshal player state: %w", err)
	}
	return &state, nil
}

func (r *RedisStore) DeletePlayer(playerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()

	if err := r.client.Del(ctx, r.playerKey(playerID)).Err(); err != nil {
		return fmt.Errorf("redis delete player %s: %w", playerID, err)
	}
	return nil
}

func (r *RedisStore) GetTopPlayers(n int) ([]*PlayerState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()

	// ZREVRANGE 取前 N 名（分數降序）
	results, err := r.client.ZRevRangeWithScores(ctx, r.leaderboardKey(), 0, int64(n-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("redis get top players: %w", err)
	}

	players := make([]*PlayerState, 0, len(results))
	for _, z := range results {
		playerID := z.Member.(string)
		state, err := r.LoadPlayer(playerID)
		if err != nil {
			log.Printf("[Store] Warning: failed to load player %s from leaderboard: %v", playerID, err)
			continue
		}
		if state == nil {
			// 排行榜有記錄但玩家資料已過期，建立最小記錄
			state = &PlayerState{
				PlayerID:     playerID,
				SessionScore: int64(z.Score),
			}
		}
		players = append(players, state)
	}
	return players, nil
}

func (r *RedisStore) UpdateLeaderboard(playerID string, score int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()

	key := r.leaderboardKey()
	// ZADD NX GT：只在新分數更高時更新（Redis 6.2+）
	// 降級方案：先取現有分數，比較後決定是否更新
	existing, err := r.client.ZScore(ctx, key, playerID).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("redis get score: %w", err)
	}

	if err == redis.Nil || float64(score) > existing {
		pipe := r.client.Pipeline()
		pipe.ZAdd(ctx, key, redis.Z{Score: float64(score), Member: playerID})
		pipe.Expire(ctx, key, leaderboardTTL)
		if _, err := pipe.Exec(ctx); err != nil {
			return fmt.Errorf("redis update leaderboard: %w", err)
		}
	}
	return nil
}

func (r *RedisStore) Close() error {
	return r.client.Close()
}

func (r *RedisStore) IsRedis() bool {
	return true
}

// SetJSON Redis 模式：儲存任意 JSON 資料（DAY-049d）
func (r *RedisStore) SetJSON(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("SetJSON marshal: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	return r.client.Set(ctx, key, data, ttl).Err()
}

// GetJSON Redis 模式：讀取任意 JSON 資料（DAY-049d）
func (r *RedisStore) GetJSON(key string, dest interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return fmt.Errorf("redis get %s: %w", key, err)
	}
	return json.Unmarshal(data, dest)
}
