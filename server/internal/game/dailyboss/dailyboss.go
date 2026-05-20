// Package dailyboss 每日特殊 BOSS 挑戰系統（DAY-077）
// 每天 UTC+8 00:00 重置一個特殊 BOSS，全服玩家合力貢獻傷害
// 擊殺後所有貢獻者按比例分配獎勵（最高 50000 金幣）
// 未擊殺則次日重置（難度降低 20%）
package dailyboss

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// BossStatus 每日 BOSS 狀態
type BossStatus string

const (
	BossStatusActive  BossStatus = "active"   // 挑戰中
	BossStatusDefeated BossStatus = "defeated" // 已擊殺
	BossStatusExpired BossStatus = "expired"  // 時間到未擊殺
)

// BossType 每日 BOSS 類型
type BossType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	BaseHP      int    `json:"base_hp"`      // 基礎 HP
	BaseReward  int    `json:"base_reward"`  // 基礎獎勵池（全服分配）
	Color       string `json:"color"`        // 顯示顏色
}

// DailyBossTypes 每日 BOSS 類型池（輪流出現）
var DailyBossTypes = []BossType{
	{
		ID:          "daily_boss_dragon",
		Name:        "海龍王",
		Icon:        "🐉",
		Description: "傳說中的海底霸主，擊敗牠可獲得豐厚寶藏",
		BaseHP:      10000,
		BaseReward:  100000,
		Color:       "#FF4444",
	},
	{
		ID:          "daily_boss_kraken",
		Name:        "深海巨怪",
		Icon:        "🦑",
		Description: "來自深淵的恐怖生物，觸手可以橫掃整片海域",
		BaseHP:      8000,
		BaseReward:  80000,
		Color:       "#8844FF",
	},
	{
		ID:          "daily_boss_whale",
		Name:        "金色鯨魚",
		Icon:        "🐋",
		Description: "全身閃耀金光的神秘鯨魚，傳說牠的身體裡裝滿了金幣",
		BaseHP:      6000,
		BaseReward:  60000,
		Color:       "#FFD700",
	},
	{
		ID:          "daily_boss_shark",
		Name:        "暗影鯊魚",
		Icon:        "🦈",
		Description: "速度極快的黑色鯊魚，只有最快的討伐者才能獲得最多獎勵",
		BaseHP:      5000,
		BaseReward:  50000,
		Color:       "#444466",
	},
	{
		ID:          "daily_boss_turtle",
		Name:        "古代神龜",
		Icon:        "🐢",
		Description: "擁有千年壽命的神龜，背上的龜殼藏著無盡的智慧與財富",
		BaseHP:      12000,
		BaseReward:  120000,
		Color:       "#44AA44",
	},
	{
		ID:          "daily_boss_jellyfish",
		Name:        "電光水母",
		Icon:        "🪼",
		Description: "發出耀眼電光的巨型水母，觸碰牠的觸手會受到強烈電擊",
		BaseHP:      7000,
		BaseReward:  70000,
		Color:       "#44DDFF",
	},
	{
		ID:          "daily_boss_crab",
		Name:        "霸王螃蟹",
		Icon:        "🦀",
		Description: "巨大的紅色螃蟹，鉗子力量足以夾碎任何東西",
		BaseHP:      9000,
		BaseReward:  90000,
		Color:       "#FF6622",
	},
}

// Contribution 玩家貢獻記錄
type Contribution struct {
	PlayerID    string `json:"player_id"`
	DisplayName string `json:"display_name"`
	Damage      int    `json:"damage"`
	Reward      int    `json:"reward"` // 結算後填入
}

// DailyBoss 每日 BOSS 實例
type DailyBoss struct {
	DateID      string     `json:"date_id"`      // 格式：2026-05-20
	BossType    BossType   `json:"boss_type"`
	MaxHP       int        `json:"max_hp"`
	CurrentHP   int        `json:"current_hp"`
	Status      BossStatus `json:"status"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       time.Time  `json:"end_at"`       // 次日 00:00
	DefeatedAt  *time.Time `json:"defeated_at"`
	Contributions map[string]*Contribution `json:"contributions"` // playerID → contribution
	TotalDamage int        `json:"total_damage"`
	RewardPool  int        `json:"reward_pool"`  // 實際獎勵池（依 HP 損失比例）
	DifficultyMod float64  `json:"difficulty_mod"` // 難度修正（連續未擊殺時降低）
}

// Manager 每日 BOSS 管理器
type Manager struct {
	mu          sync.RWMutex
	current     *DailyBoss
	history     []*DailyBoss // 最近 7 天歷史
	consecutiveFails int     // 連續未擊殺天數（用於降低難度）
}

// New 建立新的每日 BOSS 管理器
func New() *Manager {
	m := &Manager{
		history: make([]*DailyBoss, 0, 7),
	}
	m.spawnTodayBoss()
	return m
}

// spawnTodayBoss 生成今日 BOSS（非 thread-safe）
func (m *Manager) spawnTodayBoss() {
	now := time.Now()
	dateID := getDateID(now)
	start, end := getDayRange(now)

	// 依日期選擇 BOSS 類型（確保每天不同）
	dayOfYear := now.YearDay()
	bossType := DailyBossTypes[dayOfYear%len(DailyBossTypes)]

	// 難度修正：連續未擊殺時降低 20%/次，最多降低 60%
	diffMod := 1.0 - float64(m.consecutiveFails)*0.2
	if diffMod < 0.4 {
		diffMod = 0.4
	}

	maxHP := int(float64(bossType.BaseHP) * diffMod)
	rewardPool := int(float64(bossType.BaseReward) * diffMod)

	m.current = &DailyBoss{
		DateID:        dateID,
		BossType:      bossType,
		MaxHP:         maxHP,
		CurrentHP:     maxHP,
		Status:        BossStatusActive,
		StartAt:       start,
		EndAt:         end,
		Contributions: make(map[string]*Contribution),
		RewardPool:    rewardPool,
		DifficultyMod: diffMod,
	}
}

// AddDamage 玩家對每日 BOSS 造成傷害
// 回傳 (isDefeated, reward) — reward > 0 表示 BOSS 被擊殺，玩家獲得的獎勵
func (m *Manager) AddDamage(playerID, displayName string, damage int) (bool, int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.current == nil || m.current.Status != BossStatusActive {
		return false, 0
	}

	// 確保不超過剩餘 HP
	if damage > m.current.CurrentHP {
		damage = m.current.CurrentHP
	}

	// 記錄貢獻
	if _, ok := m.current.Contributions[playerID]; !ok {
		m.current.Contributions[playerID] = &Contribution{
			PlayerID:    playerID,
			DisplayName: displayName,
		}
	}
	m.current.Contributions[playerID].Damage += damage
	m.current.TotalDamage += damage
	m.current.CurrentHP -= damage

	// 檢查是否擊殺
	if m.current.CurrentHP <= 0 {
		m.current.CurrentHP = 0
		now := time.Now()
		m.current.DefeatedAt = &now
		m.current.Status = BossStatusDefeated
		m.consecutiveFails = 0

		// 計算每人獎勵
		m.settleRewards()

		// 取得觸發擊殺的玩家獎勵
		reward := 0
		if c, ok := m.current.Contributions[playerID]; ok {
			reward = c.Reward
		}
		return true, reward
	}

	return false, 0
}

// settleRewards 結算獎勵（非 thread-safe）
// 按貢獻比例分配獎勵池
func (m *Manager) settleRewards() {
	if m.current.TotalDamage == 0 {
		return
	}

	for _, c := range m.current.Contributions {
		ratio := float64(c.Damage) / float64(m.current.TotalDamage)
		c.Reward = int(float64(m.current.RewardPool) * ratio)
		// 最低保底 100 金幣（有貢獻就有獎勵）
		if c.Reward < 100 {
			c.Reward = 100
		}
	}
}

// CheckAndReset 檢查是否需要重置（每分鐘呼叫）
// 回傳 (needsReset, oldBoss) — needsReset=true 表示已重置
func (m *Manager) CheckAndReset() (bool, *DailyBoss) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.current == nil {
		m.spawnTodayBoss()
		return true, nil
	}

	now := time.Now()
	if now.Before(m.current.EndAt) {
		return false, nil
	}

	// 時間到，未擊殺
	old := m.current
	if old.Status == BossStatusActive {
		old.Status = BossStatusExpired
		m.consecutiveFails++
	}

	// 保留歷史
	m.history = append(m.history, old)
	if len(m.history) > 7 {
		m.history = m.history[1:]
	}

	// 生成新 BOSS
	m.spawnTodayBoss()
	return true, old
}

// GetSnapshot 取得當前 BOSS 快照（thread-safe）
func (m *Manager) GetSnapshot() *DailyBoss {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.current
}

// GetTopContributors 取得前 N 名貢獻者（thread-safe）
func (m *Manager) GetTopContributors(n int) []*Contribution {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.current == nil {
		return nil
	}

	contribs := make([]*Contribution, 0, len(m.current.Contributions))
	for _, c := range m.current.Contributions {
		contribs = append(contribs, c)
	}

	sort.Slice(contribs, func(i, j int) bool {
		return contribs[i].Damage > contribs[j].Damage
	})

	if n > len(contribs) {
		n = len(contribs)
	}
	return contribs[:n]
}

// GetPlayerContribution 取得玩家的貢獻（thread-safe）
func (m *Manager) GetPlayerContribution(playerID string) *Contribution {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.current == nil {
		return nil
	}
	return m.current.Contributions[playerID]
}

// GetHistory 取得歷史記錄（thread-safe）
func (m *Manager) GetHistory() []*DailyBoss {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*DailyBoss, len(m.history))
	copy(result, m.history)
	return result
}

// GetHPPercent 取得 HP 百分比（thread-safe）
func (m *Manager) GetHPPercent() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.current == nil || m.current.MaxHP == 0 {
		return 0
	}
	return float64(m.current.CurrentHP) / float64(m.current.MaxHP)
}

// GetStatus 取得狀態（thread-safe）
func (m *Manager) GetStatus() BossStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.current == nil {
		return BossStatusExpired
	}
	return m.current.Status
}

// GetDateID 取得今日 ID（thread-safe）
func (m *Manager) GetDateID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.current == nil {
		return ""
	}
	return m.current.DateID
}

// GetConsecutiveFails 取得連續未擊殺天數（thread-safe）
func (m *Manager) GetConsecutiveFails() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.consecutiveFails
}

// FormatHP 格式化 HP 顯示（如 8500/10000）
func (b *DailyBoss) FormatHP() string {
	return fmt.Sprintf("%d/%d", b.CurrentHP, b.MaxHP)
}

// getDateID 取得日期 ID（格式：2026-05-20）
func getDateID(t time.Time) string {
	loc := time.FixedZone("UTC+8", 8*60*60)
	t8 := t.In(loc)
	return fmt.Sprintf("%04d-%02d-%02d", t8.Year(), int(t8.Month()), t8.Day())
}

// getDayRange 取得今日的開始和結束時間（UTC+8）
func getDayRange(t time.Time) (start, end time.Time) {
	loc := time.FixedZone("UTC+8", 8*60*60)
	t8 := t.In(loc)
	start = time.Date(t8.Year(), t8.Month(), t8.Day(), 0, 0, 0, 0, loc).UTC()
	end = time.Date(t8.Year(), t8.Month(), t8.Day()+1, 0, 0, 0, 0, loc).UTC()
	return
}
