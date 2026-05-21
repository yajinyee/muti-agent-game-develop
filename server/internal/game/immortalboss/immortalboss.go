// Package immortalboss — 不死 BOSS 連勝系統（DAY-129）
// 業界依據：JILI Royal Fishing 2026 Immortal Boss
// Golden Toad 和 Ancient Crocodile 隨機出現，每次命中都給獎勵（50x-150x），
// 直到牠們自己離開畫面，製造「連續獲勝序列」的爽感。
package immortalboss

import (
	"math/rand"
	"sync"
	"time"
)

// BossType 不死 BOSS 類型
type BossType string

const (
	BossGoldenToad       BossType = "golden_toad"       // 金蟾蜍：50x-120x
	BossAncientCrocodile BossType = "ancient_crocodile" // 古鱷魚：60x-150x
)

// BossDef 不死 BOSS 定義
type BossDef struct {
	ID          BossType
	Name        string
	Icon        string
	MinMult     float64 // 最小倍率
	MaxMult     float64 // 最大倍率
	Duration    float64 // 在場時間（秒）
	TriggerRate float64 // 觸發機率（每次 spawnTarget 時）
	Color       string  // 顯示顏色
}

// BossDefs 不死 BOSS 定義表
var BossDefs = map[BossType]*BossDef{
	BossGoldenToad: {
		ID:          BossGoldenToad,
		Name:        "金蟾蜍",
		Icon:        "🐸",
		MinMult:     50.0,
		MaxMult:     120.0,
		Duration:    25.0, // 25 秒後離開
		TriggerRate: 0.008, // 0.8% 機率觸發
		Color:       "#FFD700",
	},
	BossAncientCrocodile: {
		ID:          BossAncientCrocodile,
		Name:        "古鱷魚",
		Icon:        "🐊",
		MinMult:     60.0,
		MaxMult:     150.0,
		Duration:    20.0, // 20 秒後離開
		TriggerRate: 0.005, // 0.5% 機率觸發
		Color:       "#228B22",
	},
}

// HitRecord 命中記錄
type HitRecord struct {
	PlayerID   string
	PlayerName string
	Multiplier float64
	Reward     int
	HitAt      time.Time
}

// Session 不死 BOSS 出現 session
type Session struct {
	InstanceID  string
	Def         *BossDef
	SpawnedAt   time.Time
	ExpiresAt   time.Time
	HitCount    int
	TotalReward int
	Hits        []HitRecord
}

// IsExpired 是否已過期
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// GetRemainingSeconds 剩餘秒數
func (s *Session) GetRemainingSeconds() float64 {
	remaining := time.Until(s.ExpiresAt).Seconds()
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Snapshot 快照
type Snapshot struct {
	Active          bool
	InstanceID      string
	BossType        BossType
	BossName        string
	BossIcon        string
	BossColor       string
	MinMult         float64
	MaxMult         float64
	HitCount        int
	TotalReward     int
	RemainingSeconds float64
}

// Manager 不死 BOSS 管理器
type Manager struct {
	mu      sync.RWMutex
	session *Session // 當前 session（nil = 無活躍 BOSS）
	lastAt  time.Time // 上次觸發時間（冷卻用）
	cooldown time.Duration
}

// New 建立新管理器
func New() *Manager {
	return &Manager{
		cooldown: 3 * time.Minute, // 3 分鐘冷卻
	}
}

// ShouldTrigger 判斷是否觸發不死 BOSS
// 在 spawnTarget 時呼叫，機率觸發
func (m *Manager) ShouldTrigger() (bool, *BossDef) {
	m.mu.RLock()
	hasActive := m.session != nil && !m.session.IsExpired()
	lastAt := m.lastAt
	m.mu.RUnlock()

	// 已有活躍 BOSS，不重複觸發
	if hasActive {
		return false, nil
	}

	// 冷卻中
	if time.Since(lastAt) < m.cooldown {
		return false, nil
	}

	// 隨機選擇 BOSS 類型並判斷觸發
	bossList := []*BossDef{BossDefs[BossGoldenToad], BossDefs[BossAncientCrocodile]}
	rand.Shuffle(len(bossList), func(i, j int) { bossList[i], bossList[j] = bossList[j], bossList[i] })

	for _, def := range bossList {
		if rand.Float64() < def.TriggerRate {
			return true, def
		}
	}
	return false, nil
}

// StartSession 開始不死 BOSS session
func (m *Manager) StartSession(instanceID string, def *BossDef) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	s := &Session{
		InstanceID: instanceID,
		Def:        def,
		SpawnedAt:  now,
		ExpiresAt:  now.Add(time.Duration(def.Duration * float64(time.Second))),
		Hits:       make([]HitRecord, 0, 20),
	}
	m.session = s
	m.lastAt = now
	return s
}

// RecordHit 記錄命中，回傳本次倍率和獎勵
func (m *Manager) RecordHit(instanceID string, playerID string, playerName string, betCost int) (float64, int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session == nil || m.session.InstanceID != instanceID {
		return 0, 0, false
	}
	if m.session.IsExpired() {
		return 0, 0, false
	}

	def := m.session.Def
	// 隨機倍率（MinMult ~ MaxMult）
	mult := def.MinMult + rand.Float64()*(def.MaxMult-def.MinMult)
	// 四捨五入到整數倍率
	multInt := int(mult)
	if multInt < int(def.MinMult) {
		multInt = int(def.MinMult)
	}
	reward := betCost * multInt

	hit := HitRecord{
		PlayerID:   playerID,
		PlayerName: playerName,
		Multiplier: float64(multInt),
		Reward:     reward,
		HitAt:      time.Now(),
	}
	m.session.Hits = append(m.session.Hits, hit)
	m.session.HitCount++
	m.session.TotalReward += reward

	return float64(multInt), reward, true
}

// CheckExpiry 檢查是否過期，回傳過期的 session（若有）
func (m *Manager) CheckExpiry() *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session == nil {
		return nil
	}
	if m.session.IsExpired() {
		expired := m.session
		m.session = nil
		return expired
	}
	return nil
}

// GetSnapshot 取得快照
func (m *Manager) GetSnapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.session == nil || m.session.IsExpired() {
		return Snapshot{Active: false}
	}
	return Snapshot{
		Active:           true,
		InstanceID:       m.session.InstanceID,
		BossType:         m.session.Def.ID,
		BossName:         m.session.Def.Name,
		BossIcon:         m.session.Def.Icon,
		BossColor:        m.session.Def.Color,
		MinMult:          m.session.Def.MinMult,
		MaxMult:          m.session.Def.MaxMult,
		HitCount:         m.session.HitCount,
		TotalReward:      m.session.TotalReward,
		RemainingSeconds: m.session.GetRemainingSeconds(),
	}
}

// IsActive 是否有活躍的不死 BOSS
func (m *Manager) IsActive() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.session != nil && !m.session.IsExpired()
}

// GetActiveInstanceID 取得活躍 BOSS 的 instanceID
func (m *Manager) GetActiveInstanceID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.session == nil || m.session.IsExpired() {
		return ""
	}
	return m.session.InstanceID
}
