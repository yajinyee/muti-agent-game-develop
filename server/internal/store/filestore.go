// Package store — FileStore JSON 檔案持久化（DAY-098）
// 不需要 Redis，Server 重啟後自動恢復所有玩家狀態
// 儲存路徑：data/players/<playerID>.json
package store

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// FullPlayerState 完整玩家持久化狀態（DAY-098）
// 包含所有子系統的資料，Server 重啟後完整恢復
type FullPlayerState struct {
	// 基礎資料
	PlayerID       string    `json:"player_id"`
	DisplayName    string    `json:"display_name"`
	Coins          int       `json:"coins"`
	MaxCoins       int       `json:"max_coins"`
	BetLevel       int       `json:"bet_level"`
	WeaponLevel    int       `json:"weapon_level"`
	KillCount      int       `json:"kill_count"`
	TotalBet       int       `json:"total_bet"`
	TotalReward    int       `json:"total_reward"`
	AttackCount    int       `json:"attack_count"`
	SessionScore   int       `json:"session_score"`
	RoomDifficulty string    `json:"room_difficulty"`
	LastSeen       time.Time `json:"last_seen"`

	// 登入資訊（DAY-065）
	LastLoginDate  string `json:"last_login_date"`
	LoginStreak    int    `json:"login_streak"`
	MaxLoginStreak int    `json:"max_login_streak"`

	// 砲台外觀（DAY-071）
	EquippedSkin string   `json:"equipped_skin"`
	OwnedSkins   []string `json:"owned_skins"`

	// VIP 系統（DAY-078）
	VIPTotalSpend   int       `json:"vip_total_spend"`
	VIPLevel        int       `json:"vip_level"`
	VIPLastWeeklyAt time.Time `json:"vip_last_weekly_at"`

	// 賽季通行證（DAY-072）
	SeasonPoints  int   `json:"season_points"`
	SeasonLevel   int   `json:"season_level"`
	SeasonClaimed []int `json:"season_claimed"`

	// 魚類圖鑑（DAY-081）
	CodexEntries []CodexEntryState `json:"codex_entries"`

	// 成就系統（DAY-100）
	Achievements []AchievementState `json:"achievements"`
	UnlockedTitles []TitleState     `json:"unlocked_titles"`
	ActiveTitle    string           `json:"active_title"`

	// 玩家統計（DAY-096）
	StatsTotalSessions   int     `json:"stats_total_sessions"`
	StatsTotalPlayTime   int64   `json:"stats_total_play_time"`
	StatsTotalShots      int     `json:"stats_total_shots"`
	StatsTotalKills      int     `json:"stats_total_kills"`
	StatsTotalBet        int     `json:"stats_total_bet"`
	StatsTotalReward     int     `json:"stats_total_reward"`
	StatsTotalBonuses    int     `json:"stats_total_bonuses"`
	StatsTotalBossKills  int     `json:"stats_total_boss_kills"`
	StatsBestMultiplier  float64 `json:"stats_best_multiplier"`
	StatsBestStreak      int     `json:"stats_best_streak"`
	StatsBestSession     int     `json:"stats_best_session"`
	StatsBestBonus       int     `json:"stats_best_bonus"`
	StatsMaxCoins        int     `json:"stats_max_coins"`
	StatsJackpotWins     int     `json:"stats_jackpot_wins"`
	StatsJackpotMini     int     `json:"stats_jackpot_mini"`
	StatsJackpotMinor    int     `json:"stats_jackpot_minor"`
	StatsJackpotMajor    int     `json:"stats_jackpot_major"`
	StatsJackpotGrand    int     `json:"stats_jackpot_grand"`
	StatsJackpotPayout   int     `json:"stats_jackpot_payout"`
	StatsHitCount        int     `json:"stats_hit_count"`
	StatsMissCount       int     `json:"stats_miss_count"`
	StatsFirstPlayAt     time.Time `json:"stats_first_play_at"`
	StatsLastPlayAt      time.Time `json:"stats_last_play_at"`
}

// CodexEntryState 圖鑑條目持久化狀態
type CodexEntryState struct {
	TargetID      string    `json:"target_id"`
	Unlocked      bool      `json:"unlocked"`
	UnlockedAt    time.Time `json:"unlocked_at,omitempty"`
	KillCount     int       `json:"kill_count"`
	MaxMultiplier float64   `json:"max_multiplier"`
}

// AchievementState 成就持久化狀態
type AchievementState struct {
	ID         string    `json:"id"`
	UnlockedAt time.Time `json:"unlocked_at"`
}

// TitleState 稱號持久化狀態
type TitleState struct {
	ID string `json:"id"`
}

// FileStore JSON 檔案持久化 Store（DAY-098）
// 每個玩家一個 JSON 檔案，儲存在 dataDir/players/ 目錄
type FileStore struct {
	dataDir string
	mu      sync.RWMutex
	// 記憶體快取（避免頻繁讀檔）
	cache map[string]*FullPlayerState
	// 排行榜（記憶體）
	leaderboard map[string]int64
}

// NewFileStore 建立 FileStore，dataDir 為資料目錄（如 "data"）
func NewFileStore(dataDir string) (*FileStore, error) {
	playersDir := filepath.Join(dataDir, "players")
	if err := os.MkdirAll(playersDir, 0755); err != nil {
		return nil, fmt.Errorf("create players dir: %w", err)
	}
	fs := &FileStore{
		dataDir:     dataDir,
		cache:       make(map[string]*FullPlayerState),
		leaderboard: make(map[string]int64),
	}

	// 統計已有的玩家資料數量
	entries, err := os.ReadDir(playersDir)
	playerCount := 0
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
				playerCount++
			}
		}
	}
	log.Printf("[FileStore] Initialized at %s (%d existing players)", playersDir, playerCount)
	return fs, nil
}

// playerPath 取得玩家 JSON 檔案路徑
func (fs *FileStore) playerPath(playerID string) string {
	// 安全化 playerID（避免路徑穿越）
	safe := filepath.Base(playerID)
	if safe == "." || safe == ".." {
		safe = "unknown"
	}
	return filepath.Join(fs.dataDir, "players", safe+".json")
}

// SaveFull 儲存完整玩家狀態
func (fs *FileStore) SaveFull(state *FullPlayerState) error {
	state.LastSeen = time.Now()

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal player state: %w", err)
	}

	path := fs.playerPath(state.PlayerID)
	// 先寫暫存檔再 rename（原子操作，避免寫到一半 crash）
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("write player file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("rename player file: %w", err)
	}

	// 更新快取
	fs.mu.Lock()
	cp := *state
	fs.cache[state.PlayerID] = &cp
	fs.mu.Unlock()

	return nil
}

// LoadFull 讀取完整玩家狀態
func (fs *FileStore) LoadFull(playerID string) (*FullPlayerState, error) {
	// 先查快取
	fs.mu.RLock()
	if cached, ok := fs.cache[playerID]; ok {
		fs.mu.RUnlock()
		cp := *cached
		return &cp, nil
	}
	fs.mu.RUnlock()

	// 讀檔
	path := fs.playerPath(playerID)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil // 新玩家，不是錯誤
	}
	if err != nil {
		return nil, fmt.Errorf("read player file: %w", err)
	}

	var state FullPlayerState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("unmarshal player state: %w", err)
	}

	// 更新快取
	fs.mu.Lock()
	cp := state
	fs.cache[playerID] = &cp
	fs.mu.Unlock()

	return &state, nil
}

// SavePlayer 實作 Store 介面（向下相容）
func (fs *FileStore) SavePlayer(state *PlayerState) error {
	full := &FullPlayerState{
		PlayerID:       state.PlayerID,
		DisplayName:    state.DisplayName,
		Coins:          int(state.Coins),
		MaxCoins:       int(state.MaxCoins),
		BetLevel:       state.BetLevel,
		KillCount:      state.KillCount,
		SessionScore:   int(state.SessionScore),
		LastLoginDate:  state.LastLoginDate,
		LoginStreak:    state.LoginStreak,
		MaxLoginStreak: state.MaxLoginStreak,
		EquippedSkin:   state.EquippedSkin,
		OwnedSkins:     state.OwnedSkins,
	}
	return fs.SaveFull(full)
}

// LoadPlayer 實作 Store 介面（向下相容）
func (fs *FileStore) LoadPlayer(playerID string) (*PlayerState, error) {
	full, err := fs.LoadFull(playerID)
	if err != nil || full == nil {
		return nil, err
	}
	return &PlayerState{
		PlayerID:       full.PlayerID,
		DisplayName:    full.DisplayName,
		Coins:          int64(full.Coins),
		MaxCoins:       int64(full.MaxCoins),
		BetLevel:       full.BetLevel,
		KillCount:      full.KillCount,
		SessionScore:   int64(full.SessionScore),
		LastLoginDate:  full.LastLoginDate,
		LoginStreak:    full.LoginStreak,
		MaxLoginStreak: full.MaxLoginStreak,
		EquippedSkin:   full.EquippedSkin,
		OwnedSkins:     full.OwnedSkins,
		LastSeen:       full.LastSeen,
	}, nil
}

// DeletePlayer 刪除玩家資料
func (fs *FileStore) DeletePlayer(playerID string) error {
	fs.mu.Lock()
	delete(fs.cache, playerID)
	fs.mu.Unlock()

	path := fs.playerPath(playerID)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete player file: %w", err)
	}
	return nil
}

// GetTopPlayers 取得排行榜前 N 名
func (fs *FileStore) GetTopPlayers(n int) ([]*PlayerState, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	type entry struct {
		playerID string
		score    int64
	}
	entries := make([]entry, 0, len(fs.leaderboard))
	for pid, score := range fs.leaderboard {
		entries = append(entries, entry{pid, score})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].score > entries[j].score
	})

	result := make([]*PlayerState, 0, n)
	for i, e := range entries {
		if i >= n {
			break
		}
		if cached, ok := fs.cache[e.playerID]; ok {
			result = append(result, &PlayerState{
				PlayerID:    cached.PlayerID,
				DisplayName: cached.DisplayName,
				Coins:       int64(cached.Coins),
				MaxCoins:    int64(cached.MaxCoins),
				KillCount:   cached.KillCount,
				SessionScore: int64(cached.SessionScore),
			})
		}
	}
	return result, nil
}

// UpdateLeaderboard 更新排行榜分數
func (fs *FileStore) UpdateLeaderboard(playerID string, score int64) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	if score > fs.leaderboard[playerID] {
		fs.leaderboard[playerID] = score
	}
	return nil
}

// SetJSON 儲存任意 JSON 資料
func (fs *FileStore) SetJSON(key string, value interface{}, _ time.Duration) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("SetJSON marshal: %w", err)
	}
	// 儲存到 dataDir/kv/ 目錄
	kvDir := filepath.Join(fs.dataDir, "kv")
	if err := os.MkdirAll(kvDir, 0755); err != nil {
		return fmt.Errorf("create kv dir: %w", err)
	}
	safe := filepath.Base(key)
	path := filepath.Join(kvDir, safe+".json")
	return os.WriteFile(path, data, 0644)
}

// GetJSON 讀取任意 JSON 資料
func (fs *FileStore) GetJSON(key string, dest interface{}) error {
	kvDir := filepath.Join(fs.dataDir, "kv")
	safe := filepath.Base(key)
	path := filepath.Join(kvDir, safe+".json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		return fmt.Errorf("read kv file: %w", err)
	}
	return json.Unmarshal(data, dest)
}

// Close 關閉（FileStore 不需要特別關閉）
func (fs *FileStore) Close() error {
	return nil
}

// IsRedis 是否使用 Redis
func (fs *FileStore) IsRedis() bool {
	return false
}
