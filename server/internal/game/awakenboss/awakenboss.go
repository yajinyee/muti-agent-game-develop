// Package awakenboss — 覺醒 BOSS 系統（DAY-130）
// 業界依據：JILI Royal Fishing 2026 Awaken Boss
// 覺醒 BOSS 比不死 BOSS 更強（90x-200x），
// 並有 Power Up 攻擊機制（6x-10x 加成），
// 需要多次命中才能觸發 Power Up，製造「蓄力爆發」的爽感。
package awakenboss

import (
	"math/rand"
	"sync"
	"time"
)

// BossType 覺醒 BOSS 類型
type BossType string

const (
	BossAwakenDragon  BossType = "awaken_dragon"  // 覺醒龍：90x-180x
	BossIcePhoenix    BossType = "ice_phoenix"     // 冰鳳凰：120x-300x（最高倍率）
	BossHumpbackWhale BossType = "humpback_whale"  // 座頭鯨：90x-150x（DAY-158，業界依據：royal-fishing.uk 2026）
	BossLegendDragon  BossType = "legend_dragon"   // 傳說龍：120x-200x（DAY-158，業界依據：royal-fishing.uk 2026）
)

// BossDef 覺醒 BOSS 定義
type BossDef struct {
	ID              BossType
	Name            string
	Icon            string
	MinMult         float64 // 基礎最小倍率
	MaxMult         float64 // 基礎最大倍率
	PowerUpMinMult  float64 // Power Up 最小加成倍率
	PowerUpMaxMult  float64 // Power Up 最大加成倍率
	PowerUpThreshold int    // 觸發 Power Up 所需命中次數
	Duration        float64 // 在場時間（秒）
	TriggerRate     float64 // 觸發機率
	Color           string  // 顯示顏色
}

// BossDefs 覺醒 BOSS 定義表
var BossDefs = map[BossType]*BossDef{
	BossAwakenDragon: {
		ID:              BossAwakenDragon,
		Name:            "覺醒龍",
		Icon:            "🐉",
		MinMult:         90.0,
		MaxMult:         180.0,
		PowerUpMinMult:  6.0,
		PowerUpMaxMult:  10.0,
		PowerUpThreshold: 5, // 5 次命中觸發 Power Up
		Duration:        30.0,
		TriggerRate:     0.003, // 0.3% 機率觸發
		Color:           "#FF4500",
	},
	BossIcePhoenix: {
		ID:              BossIcePhoenix,
		Name:            "冰鳳凰",
		Icon:            "🦅",
		MinMult:         120.0,
		MaxMult:         300.0,
		PowerUpMinMult:  8.0,
		PowerUpMaxMult:  10.0,
		PowerUpThreshold: 8, // 8 次命中觸發 Power Up（更難但更強）
		Duration:        25.0,
		TriggerRate:     0.001, // 0.1% 機率觸發（稀有）
		Color:           "#00BFFF",
	},
	// DAY-158：座頭鯨（業界依據：royal-fishing.uk 2026「Humpback Whale offers 90-150x with 15x base multiplier」）
	// 座頭鯨是 Royal Fishing 的標誌性覺醒 BOSS，15x 基礎倍率，Power Up 後最高 150x
	BossHumpbackWhale: {
		ID:              BossHumpbackWhale,
		Name:            "座頭鯨",
		Icon:            "🐋",
		MinMult:         90.0,
		MaxMult:         150.0,
		PowerUpMinMult:  6.0,
		PowerUpMaxMult:  8.0,
		PowerUpThreshold: 6, // 6 次命中觸發 Power Up
		Duration:        35.0, // 在場時間較長（鯨魚移動慢）
		TriggerRate:     0.002, // 0.2% 機率觸發
		Color:           "#4169E1",
	},
	// DAY-158：傳說龍（業界依據：royal-fishing.uk 2026「Legend Dragon reaches 120-200x from 20x base」）
	// 傳說龍是 Royal Fishing 最強的覺醒 BOSS，20x 基礎倍率，Power Up 後最高 200x
	BossLegendDragon: {
		ID:              BossLegendDragon,
		Name:            "傳說龍",
		Icon:            "🐲",
		MinMult:         120.0,
		MaxMult:         200.0,
		PowerUpMinMult:  8.0,
		PowerUpMaxMult:  10.0,
		PowerUpThreshold: 10, // 10 次命中觸發 Power Up（最難但最強）
		Duration:        20.0, // 在場時間較短（稀有感）
		TriggerRate:     0.0008, // 0.08% 機率觸發（最稀有）
		Color:           "#9400D3",
	},
}

// HitRecord 命中記錄
type HitRecord struct {
	PlayerID   string
	PlayerName string
	Multiplier float64
	Reward     int
	IsPowerUp  bool
	HitAt      time.Time
}

// Session 覺醒 BOSS session
type Session struct {
	InstanceID      string
	Def             *BossDef
	SpawnedAt       time.Time
	ExpiresAt       time.Time
	HitCount        int       // 總命中次數
	PowerUpCount    int       // Power Up 觸發次數
	HitsSincePowerUp int      // 上次 Power Up 後的命中次數
	TotalReward     int
	Hits            []HitRecord
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

// GetPowerUpProgress 取得 Power Up 進度（0.0-1.0）
func (s *Session) GetPowerUpProgress() float64 {
	if s.Def.PowerUpThreshold <= 0 {
		return 0
	}
	progress := float64(s.HitsSincePowerUp) / float64(s.Def.PowerUpThreshold)
	if progress > 1.0 {
		return 1.0
	}
	return progress
}

// Snapshot 快照
type Snapshot struct {
	Active           bool
	InstanceID       string
	BossType         BossType
	BossName         string
	BossIcon         string
	BossColor        string
	MinMult          float64
	MaxMult          float64
	PowerUpMinMult   float64
	PowerUpMaxMult   float64
	PowerUpThreshold int
	HitCount         int
	PowerUpCount     int
	PowerUpProgress  float64
	TotalReward      int
	RemainingSeconds float64
}

// Manager 覺醒 BOSS 管理器
type Manager struct {
	mu       sync.RWMutex
	session  *Session
	lastAt   time.Time
	cooldown time.Duration
}

// New 建立新管理器
func New() *Manager {
	return &Manager{
		cooldown: 5 * time.Minute, // 5 分鐘冷卻（比不死 BOSS 更長）
	}
}

// ShouldTrigger 判斷是否觸發覺醒 BOSS
func (m *Manager) ShouldTrigger() (bool, *BossDef) {
	m.mu.RLock()
	hasActive := m.session != nil && !m.session.IsExpired()
	lastAt := m.lastAt
	m.mu.RUnlock()

	if hasActive {
		return false, nil
	}
	if time.Since(lastAt) < m.cooldown {
		return false, nil
	}

	// 隨機選擇 BOSS 類型（DAY-158：加入座頭鯨和傳說龍）
	bossList := []*BossDef{BossDefs[BossAwakenDragon], BossDefs[BossIcePhoenix], BossDefs[BossHumpbackWhale], BossDefs[BossLegendDragon]}
	rand.Shuffle(len(bossList), func(i, j int) { bossList[i], bossList[j] = bossList[j], bossList[i] })

	for _, def := range bossList {
		if rand.Float64() < def.TriggerRate {
			return true, def
		}
	}
	return false, nil
}

// StartSession 開始覺醒 BOSS session
func (m *Manager) StartSession(instanceID string, def *BossDef) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	s := &Session{
		InstanceID: instanceID,
		Def:        def,
		SpawnedAt:  now,
		ExpiresAt:  now.Add(time.Duration(def.Duration * float64(time.Second))),
		Hits:       make([]HitRecord, 0, 30),
	}
	m.session = s
	m.lastAt = now
	return s
}

// RecordHit 記錄命中，回傳（倍率, 獎勵, 是否 Power Up, 是否成功）
func (m *Manager) RecordHit(instanceID string, playerID string, playerName string, betCost int) (float64, int, bool, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session == nil || m.session.InstanceID != instanceID {
		return 0, 0, false, false
	}
	if m.session.IsExpired() {
		return 0, 0, false, false
	}

	def := m.session.Def
	m.session.HitCount++
	m.session.HitsSincePowerUp++

	// 判斷是否觸發 Power Up
	isPowerUp := m.session.HitsSincePowerUp >= def.PowerUpThreshold
	if isPowerUp {
		m.session.PowerUpCount++
		m.session.HitsSincePowerUp = 0
	}

	// 計算倍率
	var mult float64
	if isPowerUp {
		// Power Up 倍率：基礎倍率 × Power Up 加成
		baseMult := def.MinMult + rand.Float64()*(def.MaxMult-def.MinMult)
		powerMult := def.PowerUpMinMult + rand.Float64()*(def.PowerUpMaxMult-def.PowerUpMinMult)
		mult = baseMult * powerMult
	} else {
		mult = def.MinMult + rand.Float64()*(def.MaxMult-def.MinMult)
	}

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
		IsPowerUp:  isPowerUp,
		HitAt:      time.Now(),
	}
	m.session.Hits = append(m.session.Hits, hit)
	m.session.TotalReward += reward

	return float64(multInt), reward, isPowerUp, true
}

// CheckExpiry 檢查是否過期
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
		PowerUpMinMult:   m.session.Def.PowerUpMinMult,
		PowerUpMaxMult:   m.session.Def.PowerUpMaxMult,
		PowerUpThreshold: m.session.Def.PowerUpThreshold,
		HitCount:         m.session.HitCount,
		PowerUpCount:     m.session.PowerUpCount,
		PowerUpProgress:  m.session.GetPowerUpProgress(),
		TotalReward:      m.session.TotalReward,
		RemainingSeconds: m.session.GetRemainingSeconds(),
	}
}

// IsActive 是否有活躍的覺醒 BOSS
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
