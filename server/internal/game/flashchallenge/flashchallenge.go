// Package flashchallenge 閃電挑戰系統（DAY-123）
// 業界依據：Infingame（2026-05-19）確認 Challenges 工具是 2026 年最熱門留存機制
// 限時 90 秒的特殊目標挑戰，完成獎勵豐厚，全服可見增加社交競爭感
package flashchallenge

import (
	"math/rand"
	"sync"
	"time"
)

// ChallengeState 挑戰狀態
type ChallengeState string

const (
	StateIdle    ChallengeState = "idle"    // 無挑戰
	StateActive  ChallengeState = "active"  // 挑戰進行中
	StateSuccess ChallengeState = "success" // 挑戰成功
	StateFailed  ChallengeState = "failed"  // 挑戰失敗（時間到）
)

// ChallengeType 挑戰類型
type ChallengeType string

const (
	TypeKillCount    ChallengeType = "kill_count"    // 擊破指定數量目標
	TypeKillSpecific ChallengeType = "kill_specific" // 擊破指定種類目標
	TypeKillStreak   ChallengeType = "kill_streak"   // 達到指定連擊數
	TypeKillBoss     ChallengeType = "kill_boss"      // 擊殺 BOSS
	TypeHighMult     ChallengeType = "high_mult"      // 獲得高倍率擊破（≥10x）
)

// ChallengeDef 挑戰定義
type ChallengeDef struct {
	Type        ChallengeType `json:"type"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Icon        string        `json:"icon"`
	Color       string        `json:"color"`
	Target      int           `json:"target"`       // 目標數量
	TargetDefID string        `json:"target_def_id"` // 指定目標種類（TypeKillSpecific 用）
	Duration    int           `json:"duration"`     // 持續秒數
	BaseReward  int           `json:"base_reward"`  // 基礎獎勵（金幣）
	BonusReward int           `json:"bonus_reward"` // 超額完成獎勵
	Weight      int           `json:"-"`            // 觸發權重
}

// availableChallenges 所有可能的挑戰定義
var availableChallenges = []ChallengeDef{
	{
		Type:        TypeKillCount,
		Title:       "⚡ 閃電獵殺",
		Description: "90秒內擊破 15 個目標",
		Icon:        "⚡",
		Color:       "#FFD700",
		Target:      15,
		Duration:    90,
		BaseReward:  3000,
		BonusReward: 1500,
		Weight:      30,
	},
	{
		Type:        TypeKillCount,
		Title:       "🔥 狂熱模式",
		Description: "60秒內擊破 10 個目標",
		Icon:        "🔥",
		Color:       "#FF6B35",
		Target:      10,
		Duration:    60,
		BaseReward:  2000,
		BonusReward: 1000,
		Weight:      25,
	},
	{
		Type:        TypeKillSpecific,
		Title:       "🍄 蘑菇獵人",
		Description: "90秒內擊破 5 個巨大蘑菇",
		Icon:        "🍄",
		Color:       "#8B4513",
		Target:      5,
		TargetDefID: "T006",
		Duration:    90,
		BaseReward:  4000,
		BonusReward: 2000,
		Weight:      15,
	},
	{
		Type:        TypeKillSpecific,
		Title:       "🪙 金幣狂潮",
		Description: "120秒內擊破 3 條金幣魚",
		Icon:        "🪙",
		Color:       "#FFD700",
		Target:      3,
		TargetDefID: "T105",
		Duration:    120,
		BaseReward:  8000,
		BonusReward: 4000,
		Weight:      10,
	},
	{
		Type:        TypeKillStreak,
		Title:       "💥 連擊大師",
		Description: "90秒內達到 10 連擊",
		Icon:        "💥",
		Color:       "#FF4500",
		Target:      10,
		Duration:    90,
		BaseReward:  5000,
		BonusReward: 2500,
		Weight:      15,
	},
	{
		Type:        TypeHighMult,
		Title:       "✨ 高倍獵手",
		Description: "60秒內獲得 3 次 10x 以上擊破",
		Icon:        "✨",
		Color:       "#9B59B6",
		Target:      3,
		Duration:    60,
		BaseReward:  6000,
		BonusReward: 3000,
		Weight:      5,
	},
}

// PlayerProgress 玩家在當前挑戰的進度
type PlayerProgress struct {
	PlayerID    string
	PlayerName  string
	Progress    int       // 當前進度
	Completed   bool      // 是否已完成
	CompletedAt time.Time // 完成時間
}

// Challenge 當前進行中的挑戰
type Challenge struct {
	Def       ChallengeDef
	State     ChallengeState
	StartedAt time.Time
	EndsAt    time.Time
	Progress  map[string]*PlayerProgress // playerID -> 進度
}

// GetTimeLeft 取得剩餘秒數
func (c *Challenge) GetTimeLeft() int {
	if c.State != StateActive {
		return 0
	}
	left := time.Until(c.EndsAt).Seconds()
	if left < 0 {
		return 0
	}
	return int(left)
}

// IsExpired 是否已超時
func (c *Challenge) IsExpired() bool {
	return c.State == StateActive && time.Now().After(c.EndsAt)
}

// GetTopPlayers 取得進度前 5 名玩家
func (c *Challenge) GetTopPlayers() []*PlayerProgress {
	players := make([]*PlayerProgress, 0, len(c.Progress))
	for _, p := range c.Progress {
		players = append(players, p)
	}
	// 按進度排序（完成者優先，再按完成時間）
	for i := 0; i < len(players); i++ {
		for j := i + 1; j < len(players); j++ {
			a, b := players[i], players[j]
			if !a.Completed && b.Completed {
				players[i], players[j] = b, a
			} else if a.Completed && b.Completed {
				if a.CompletedAt.After(b.CompletedAt) {
					players[i], players[j] = b, a
				}
			} else if a.Progress < b.Progress {
				players[i], players[j] = b, a
			}
		}
	}
	if len(players) > 5 {
		return players[:5]
	}
	return players
}

// Snapshot 挑戰快照（用於廣播）
type Snapshot struct {
	State       ChallengeState `json:"state"`
	Type        ChallengeType  `json:"type"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Icon        string         `json:"icon"`
	Color       string         `json:"color"`
	Target      int            `json:"target"`
	TargetDefID string         `json:"target_def_id"`
	Duration    int            `json:"duration"`
	TimeLeft    int            `json:"time_left"`
	BaseReward  int            `json:"base_reward"`
	BonusReward int            `json:"bonus_reward"`
	TopPlayers  []PlayerSnap   `json:"top_players"`
}

// PlayerSnap 玩家進度快照
type PlayerSnap struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	Progress   int    `json:"progress"`
	Completed  bool   `json:"completed"`
}

// Manager 閃電挑戰管理器
type Manager struct {
	mu              sync.RWMutex
	current         *Challenge
	rng             *rand.Rand
	lastTriggerAt   time.Time
	minIntervalSecs int // 兩次挑戰最短間隔（秒）
}

// New 建立閃電挑戰管理器
func New() *Manager {
	return &Manager{
		rng:             rand.New(rand.NewSource(time.Now().UnixNano())),
		minIntervalSecs: 300, // 最短 5 分鐘間隔
	}
}

// ShouldTrigger 是否應該觸發新挑戰（由 game loop 定期呼叫）
// triggerType: "random"（隨機觸發）/ "boss"（BOSS 擊殺後觸發）
func (m *Manager) ShouldTrigger(triggerType string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 已有進行中的挑戰
	if m.current != nil && m.current.State == StateActive {
		return false
	}

	// 距離上次觸發太近
	if time.Since(m.lastTriggerAt).Seconds() < float64(m.minIntervalSecs) {
		return false
	}

	switch triggerType {
	case "boss":
		return true // BOSS 擊殺後必定觸發
	case "random":
		return m.rng.Float64() < 0.15 // 15% 機率（每次 game loop 呼叫）
	}
	return false
}

// StartChallenge 開始新挑戰
// 回傳挑戰快照（nil = 無法開始）
func (m *Manager) StartChallenge() *Snapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 選擇挑戰類型（加權隨機）
	def := m.pickChallenge()
	if def == nil {
		return nil
	}

	m.current = &Challenge{
		Def:       *def,
		State:     StateActive,
		StartedAt: time.Now(),
		EndsAt:    time.Now().Add(time.Duration(def.Duration) * time.Second),
		Progress:  make(map[string]*PlayerProgress),
	}
	m.lastTriggerAt = time.Now()

	return m.buildSnapshot()
}

// RecordKill 記錄玩家擊破（由 handleKill 呼叫）
// 回傳：(進度更新, 是否完成, 是否首次完成)
func (m *Manager) RecordKill(playerID, playerName, targetDefID string, mult float64, streak int) (int, bool, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.current == nil || m.current.State != StateActive {
		return 0, false, false
	}
	if m.current.IsExpired() {
		m.current.State = StateFailed
		return 0, false, false
	}

	// 確保玩家有進度記錄
	pp, ok := m.current.Progress[playerID]
	if !ok {
		pp = &PlayerProgress{
			PlayerID:   playerID,
			PlayerName: playerName,
		}
		m.current.Progress[playerID] = pp
	}

	// 已完成的玩家不再計算
	if pp.Completed {
		return pp.Progress, true, false
	}

	def := m.current.Def
	shouldCount := false

	switch def.Type {
	case TypeKillCount:
		shouldCount = true
	case TypeKillSpecific:
		shouldCount = targetDefID == def.TargetDefID
	case TypeKillStreak:
		shouldCount = streak >= def.Target
		if shouldCount {
			pp.Progress = def.Target // 直接設為完成
			pp.Completed = true
			pp.CompletedAt = time.Now()
			return pp.Progress, true, true
		}
		return pp.Progress, false, false
	case TypeHighMult:
		shouldCount = mult >= 10.0
	case TypeKillBoss:
		shouldCount = targetDefID == "B001"
	}

	if shouldCount {
		pp.Progress++
		if pp.Progress >= def.Target {
			pp.Completed = true
			pp.CompletedAt = time.Now()
			return pp.Progress, true, true
		}
	}

	return pp.Progress, false, false
}

// CheckExpiry 檢查挑戰是否超時（由 game loop 定期呼叫）
// 回傳：是否剛剛超時（需要廣播結果）
func (m *Manager) CheckExpiry() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.current == nil || m.current.State != StateActive {
		return false
	}
	if m.current.IsExpired() {
		m.current.State = StateFailed
		return true
	}
	return false
}

// GetSnapshot 取得當前挑戰快照
func (m *Manager) GetSnapshot() *Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.buildSnapshot()
}

// GetPlayerProgress 取得玩家進度
func (m *Manager) GetPlayerProgress(playerID string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.current == nil {
		return 0, false
	}
	pp, ok := m.current.Progress[playerID]
	if !ok {
		return 0, false
	}
	return pp.Progress, pp.Completed
}

// CalcReward 計算玩家獎勵
// progress: 玩家完成進度；completed: 是否完成
func (m *Manager) CalcReward(progress int, completed bool) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.current == nil {
		return 0
	}
	def := m.current.Def
	if completed {
		return def.BaseReward + def.BonusReward
	}
	// 部分完成：按比例給安慰獎（最少 10%）
	if progress <= 0 || def.Target <= 0 {
		return 0
	}
	ratio := float64(progress) / float64(def.Target)
	if ratio < 0.1 {
		return 0
	}
	return int(float64(def.BaseReward) * ratio * 0.5) // 安慰獎 = 基礎獎勵 × 進度比例 × 50%
}

// GetCurrentDef 取得當前挑戰定義（nil = 無挑戰）
func (m *Manager) GetCurrentDef() *ChallengeDef {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current == nil {
		return nil
	}
	def := m.current.Def
	return &def
}

// IsActive 是否有進行中的挑戰
func (m *Manager) IsActive() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.current != nil && m.current.State == StateActive && !m.current.IsExpired()
}

// ---- 內部輔助函數 ----

// pickChallenge 加權隨機選擇挑戰
func (m *Manager) pickChallenge() *ChallengeDef {
	totalWeight := 0
	for _, c := range availableChallenges {
		totalWeight += c.Weight
	}
	roll := m.rng.Intn(totalWeight)
	cumulative := 0
	for i := range availableChallenges {
		cumulative += availableChallenges[i].Weight
		if roll < cumulative {
			def := availableChallenges[i]
			return &def
		}
	}
	def := availableChallenges[0]
	return &def
}

// buildSnapshot 建立快照（需持有鎖）
func (m *Manager) buildSnapshot() *Snapshot {
	if m.current == nil {
		return &Snapshot{State: StateIdle}
	}
	c := m.current
	topPlayers := c.GetTopPlayers()
	snaps := make([]PlayerSnap, 0, len(topPlayers))
	for _, p := range topPlayers {
		snaps = append(snaps, PlayerSnap{
			PlayerID:   p.PlayerID,
			PlayerName: p.PlayerName,
			Progress:   p.Progress,
			Completed:  p.Completed,
		})
	}
	return &Snapshot{
		State:       c.State,
		Type:        c.Def.Type,
		Title:       c.Def.Title,
		Description: c.Def.Description,
		Icon:        c.Def.Icon,
		Color:       c.Def.Color,
		Target:      c.Def.Target,
		TargetDefID: c.Def.TargetDefID,
		Duration:    c.Def.Duration,
		TimeLeft:    c.GetTimeLeft(),
		BaseReward:  c.Def.BaseReward,
		BonusReward: c.Def.BonusReward,
		TopPlayers:  snaps,
	}
}
