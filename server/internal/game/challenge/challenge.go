// Package challenge 隱藏挑戰系統（DAY-085）
// 玩家不知道挑戰存在，完成特定條件後突然解鎖，增加驚喜感
// 業界依據：Fish Hunters（2026）確認隱藏成就提升留存率 40%+
package challenge

import (
	"sync"
	"time"
)

// ChallengeID 挑戰 ID
type ChallengeID string

const (
	// 連擊類
	ChallengeStreak5  ChallengeID = "streak_5"  // 達成 5 連擊
	ChallengeStreak10 ChallengeID = "streak_10" // 達成 10 連擊
	ChallengeStreak20 ChallengeID = "streak_20" // 達成 20 連擊（傳說連擊）

	// 倍率類
	ChallengeMult50  ChallengeID = "mult_50"  // 單次擊破 50x+
	ChallengeMult100 ChallengeID = "mult_100" // 單次擊破 100x+（轉盤最高）

	// 速度類
	ChallengeSpeed3 ChallengeID = "speed_3" // 3 秒內擊破 3 個目標
	ChallengeSpeed5 ChallengeID = "speed_5" // 5 秒內擊破 5 個目標

	// 收集類
	ChallengeAllTypes  ChallengeID = "all_types"  // 在一局中擊破所有類型目標
	ChallengeBossFirst ChallengeID = "boss_first" // 首次擊敗 BOSS

	// 財富類
	ChallengeRich10k ChallengeID = "rich_10k" // 單局累積 10000 金幣
	ChallengeRich50k ChallengeID = "rich_50k" // 單局累積 50000 金幣

	// 特殊類
	ChallengeWheelMax ChallengeID = "wheel_max" // 轉盤獲得 100x
	ChallengeJackpot  ChallengeID = "jackpot"   // 中任意 Jackpot
	// 任務連續完成類（DAY-086）
	ChallengeMissionStreak7  ChallengeID = "mission_streak_7"  // 連續 7 天完成所有任務
	ChallengeMissionStreak30 ChallengeID = "mission_streak_30" // 連續 30 天完成所有任務
)

// ChallengeDef 挑戰定義
type ChallengeDef struct {
	ID          ChallengeID `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Icon        string      `json:"icon"`
	Reward      int         `json:"reward"`    // 解鎖獎勵（金幣）
	TitleID     string      `json:"title_id"`  // 解鎖稱號（可選）
	IsHidden    bool        `json:"is_hidden"` // 是否隱藏（解鎖前不顯示）
}

// Challenges 所有挑戰定義
var Challenges = map[ChallengeID]*ChallengeDef{
	ChallengeStreak5: {
		ID:          ChallengeStreak5,
		Name:        "連擊新手",
		Description: "達成 5 連擊",
		Icon:        "🔥",
		Reward:      500,
		IsHidden:    false,
	},
	ChallengeStreak10: {
		ID:          ChallengeStreak10,
		Name:        "連擊高手",
		Description: "達成 10 連擊",
		Icon:        "💥",
		Reward:      2000,
		IsHidden:    true,
	},
	ChallengeStreak20: {
		ID:          ChallengeStreak20,
		Name:        "傳說連擊者",
		Description: "達成 20 連擊（傳說等級）",
		Icon:        "⚡",
		Reward:      10000,
		TitleID:     "streak_legend",
		IsHidden:    true,
	},
	ChallengeMult50: {
		ID:          ChallengeMult50,
		Name:        "大獎獵人",
		Description: "單次擊破獲得 50x 以上獎勵",
		Icon:        "🎯",
		Reward:      3000,
		IsHidden:    true,
	},
	ChallengeMult100: {
		ID:          ChallengeMult100,
		Name:        "百倍傳說",
		Description: "單次擊破獲得 100x 獎勵",
		Icon:        "💯",
		Reward:      20000,
		TitleID:     "mult_legend",
		IsHidden:    true,
	},
	ChallengeSpeed3: {
		ID:          ChallengeSpeed3,
		Name:        "閃電討伐",
		Description: "3 秒內擊破 3 個目標",
		Icon:        "⚡",
		Reward:      1500,
		IsHidden:    true,
	},
	ChallengeSpeed5: {
		ID:          ChallengeSpeed5,
		Name:        "疾風討伐",
		Description: "5 秒內擊破 5 個目標",
		Icon:        "🌪️",
		Reward:      5000,
		TitleID:     "speed_master",
		IsHidden:    true,
	},
	ChallengeAllTypes: {
		ID:          ChallengeAllTypes,
		Name:        "全能討伐者",
		Description: "在一局中擊破所有類型的目標",
		Icon:        "🌟",
		Reward:      8000,
		TitleID:     "all_rounder",
		IsHidden:    true,
	},
	ChallengeBossFirst: {
		ID:          ChallengeBossFirst,
		Name:        "BOSS 終結者",
		Description: "首次擊敗 BOSS",
		Icon:        "👑",
		Reward:      5000,
		IsHidden:    false,
	},
	ChallengeRich10k: {
		ID:          ChallengeRich10k,
		Name:        "小富翁",
		Description: "單局累積獲得 10000 金幣",
		Icon:        "💰",
		Reward:      2000,
		IsHidden:    true,
	},
	ChallengeRich50k: {
		ID:          ChallengeRich50k,
		Name:        "大富翁",
		Description: "單局累積獲得 50000 金幣",
		Icon:        "💎",
		Reward:      10000,
		TitleID:     "big_winner",
		IsHidden:    true,
	},
	ChallengeWheelMax: {
		ID:          ChallengeWheelMax,
		Name:        "轉盤之神",
		Description: "幸運轉盤獲得 100x 最高獎勵",
		Icon:        "🎰",
		Reward:      15000,
		TitleID:     "wheel_god",
		IsHidden:    true,
	},
	ChallengeJackpot: {
		ID:          ChallengeJackpot,
		Name:        "Jackpot 得主",
		Description: "中任意等級的 Jackpot",
		Icon:        "🏆",
		Reward:      5000,
		IsHidden:    false,
	},
	ChallengeMissionStreak7: {
		ID:          ChallengeMissionStreak7,
		Name:        "任務達人",
		Description: "連續 7 天完成所有每日任務",
		Icon:        "📅",
		Reward:      30000,
		TitleID:     "mission_master",
		IsHidden:    true,
	},
	ChallengeMissionStreak30: {
		ID:          ChallengeMissionStreak30,
		Name:        "任務傳說",
		Description: "連續 30 天完成所有每日任務",
		Icon:        "🌟",
		Reward:      200000,
		TitleID:     "mission_legend",
		IsHidden:    true,
	},
}

// PlayerChallengeState 玩家挑戰狀態
type PlayerChallengeState struct {
	Unlocked      bool      `json:"unlocked"`
	UnlockedAt    time.Time `json:"unlocked_at,omitempty"`
	RewardClaimed bool      `json:"reward_claimed"`
}

// SessionStats 單局統計（用於速度/收集類挑戰）
type SessionStats struct {
	KillTimestamps []time.Time     // 擊破時間戳（用於速度挑戰）
	KilledTypes    map[string]bool // 已擊破的目標類型
	TotalCoins     int             // 本局累積金幣
}

// Manager 挑戰管理器
type Manager struct {
	mu       sync.RWMutex
	states   map[string]map[ChallengeID]*PlayerChallengeState // playerID → challengeID → state
	sessions map[string]*SessionStats                         // playerID → session stats
}

// NewManager 建立挑戰管理器
func NewManager() *Manager {
	return &Manager{
		states:   make(map[string]map[ChallengeID]*PlayerChallengeState),
		sessions: make(map[string]*SessionStats),
	}
}

// InitPlayer 初始化玩家挑戰狀態
func (m *Manager) InitPlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.states[playerID]; !ok {
		m.states[playerID] = make(map[ChallengeID]*PlayerChallengeState)
	}
	if _, ok := m.sessions[playerID]; !ok {
		m.sessions[playerID] = &SessionStats{
			KilledTypes: make(map[string]bool),
		}
	}
}

// RemovePlayer 移除玩家（重置 session stats）
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, playerID)
}

// isUnlocked 檢查挑戰是否已解鎖（內部，不加鎖）
func (m *Manager) isUnlocked(playerID string, id ChallengeID) bool {
	states, ok := m.states[playerID]
	if !ok {
		return false
	}
	s, ok := states[id]
	return ok && s.Unlocked
}

// TryUnlock 嘗試解鎖挑戰，回傳新解鎖的挑戰定義（nil = 未解鎖或已解鎖過）
func (m *Manager) TryUnlock(playerID string, id ChallengeID) *ChallengeDef {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isUnlocked(playerID, id) {
		return nil
	}

	def, ok := Challenges[id]
	if !ok {
		return nil
	}

	if _, ok := m.states[playerID]; !ok {
		m.states[playerID] = make(map[ChallengeID]*PlayerChallengeState)
	}
	m.states[playerID][id] = &PlayerChallengeState{
		Unlocked:   true,
		UnlockedAt: time.Now(),
	}
	return def
}

// ClaimReward 領取挑戰獎勵，回傳獎勵金幣（0 = 無法領取）
func (m *Manager) ClaimReward(playerID string, id ChallengeID) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	states, ok := m.states[playerID]
	if !ok {
		return 0
	}
	s, ok := states[id]
	if !ok || !s.Unlocked || s.RewardClaimed {
		return 0
	}

	def, ok := Challenges[id]
	if !ok {
		return 0
	}

	s.RewardClaimed = true
	return def.Reward
}

// RecordKill 記錄擊破事件，回傳新解鎖的速度挑戰（nil = 無）
func (m *Manager) RecordKill(playerID string, defID string, multiplier float64) []*ChallengeDef {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[playerID]; !ok {
		m.sessions[playerID] = &SessionStats{
			KilledTypes: make(map[string]bool),
		}
	}
	sess := m.sessions[playerID]

	now := time.Now()
	sess.KillTimestamps = append(sess.KillTimestamps, now)

	// 記錄目標類型
	sess.KilledTypes[defID] = true

	var unlocked []*ChallengeDef

	// 速度挑戰：3 秒內 3 個
	count3s := 0
	for _, ts := range sess.KillTimestamps {
		if now.Sub(ts) <= 3*time.Second {
			count3s++
		}
	}
	if count3s >= 3 && !m.isUnlocked(playerID, ChallengeSpeed3) {
		if _, ok := m.states[playerID]; !ok {
			m.states[playerID] = make(map[ChallengeID]*PlayerChallengeState)
		}
		m.states[playerID][ChallengeSpeed3] = &PlayerChallengeState{Unlocked: true, UnlockedAt: now}
		unlocked = append(unlocked, Challenges[ChallengeSpeed3])
	}

	// 速度挑戰：5 秒內 5 個
	count5s := 0
	for _, ts := range sess.KillTimestamps {
		if now.Sub(ts) <= 5*time.Second {
			count5s++
		}
	}
	if count5s >= 5 && !m.isUnlocked(playerID, ChallengeSpeed5) {
		if _, ok := m.states[playerID]; !ok {
			m.states[playerID] = make(map[ChallengeID]*PlayerChallengeState)
		}
		m.states[playerID][ChallengeSpeed5] = &PlayerChallengeState{Unlocked: true, UnlockedAt: now}
		unlocked = append(unlocked, Challenges[ChallengeSpeed5])
	}

	// 倍率挑戰
	if multiplier >= 50 && !m.isUnlocked(playerID, ChallengeMult50) {
		if _, ok := m.states[playerID]; !ok {
			m.states[playerID] = make(map[ChallengeID]*PlayerChallengeState)
		}
		m.states[playerID][ChallengeMult50] = &PlayerChallengeState{Unlocked: true, UnlockedAt: now}
		unlocked = append(unlocked, Challenges[ChallengeMult50])
	}
	if multiplier >= 100 && !m.isUnlocked(playerID, ChallengeMult100) {
		if _, ok := m.states[playerID]; !ok {
			m.states[playerID] = make(map[ChallengeID]*PlayerChallengeState)
		}
		m.states[playerID][ChallengeMult100] = &PlayerChallengeState{Unlocked: true, UnlockedAt: now}
		unlocked = append(unlocked, Challenges[ChallengeMult100])
	}

	// 清理超過 10 秒的時間戳（節省記憶體）
	cutoff := now.Add(-10 * time.Second)
	filtered := sess.KillTimestamps[:0]
	for _, ts := range sess.KillTimestamps {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}
	sess.KillTimestamps = filtered

	return unlocked
}

// RecordCoins 記錄本局累積金幣，回傳新解鎖的財富挑戰
func (m *Manager) RecordCoins(playerID string, amount int) []*ChallengeDef {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[playerID]; !ok {
		m.sessions[playerID] = &SessionStats{KilledTypes: make(map[string]bool)}
	}
	m.sessions[playerID].TotalCoins += amount

	total := m.sessions[playerID].TotalCoins
	var unlocked []*ChallengeDef
	now := time.Now()

	if total >= 10000 && !m.isUnlocked(playerID, ChallengeRich10k) {
		if _, ok := m.states[playerID]; !ok {
			m.states[playerID] = make(map[ChallengeID]*PlayerChallengeState)
		}
		m.states[playerID][ChallengeRich10k] = &PlayerChallengeState{Unlocked: true, UnlockedAt: now}
		unlocked = append(unlocked, Challenges[ChallengeRich10k])
	}
	if total >= 50000 && !m.isUnlocked(playerID, ChallengeRich50k) {
		if _, ok := m.states[playerID]; !ok {
			m.states[playerID] = make(map[ChallengeID]*PlayerChallengeState)
		}
		m.states[playerID][ChallengeRich50k] = &PlayerChallengeState{Unlocked: true, UnlockedAt: now}
		unlocked = append(unlocked, Challenges[ChallengeRich50k])
	}
	return unlocked
}

// RecordStreak 記錄連擊數，回傳新解鎖的連擊挑戰
func (m *Manager) RecordStreak(playerID string, streak int) []*ChallengeDef {
	m.mu.Lock()
	defer m.mu.Unlock()

	var unlocked []*ChallengeDef
	now := time.Now()

	if _, ok := m.states[playerID]; !ok {
		m.states[playerID] = make(map[ChallengeID]*PlayerChallengeState)
	}

	if streak >= 5 && !m.isUnlocked(playerID, ChallengeStreak5) {
		m.states[playerID][ChallengeStreak5] = &PlayerChallengeState{Unlocked: true, UnlockedAt: now}
		unlocked = append(unlocked, Challenges[ChallengeStreak5])
	}
	if streak >= 10 && !m.isUnlocked(playerID, ChallengeStreak10) {
		m.states[playerID][ChallengeStreak10] = &PlayerChallengeState{Unlocked: true, UnlockedAt: now}
		unlocked = append(unlocked, Challenges[ChallengeStreak10])
	}
	if streak >= 20 && !m.isUnlocked(playerID, ChallengeStreak20) {
		m.states[playerID][ChallengeStreak20] = &PlayerChallengeState{Unlocked: true, UnlockedAt: now}
		unlocked = append(unlocked, Challenges[ChallengeStreak20])
	}
	return unlocked
}

// RecordAllTypes 記錄全類型擊破，回傳是否解鎖全能挑戰
// allTargetTypes 是遊戲中所有目標類型的 ID 集合
func (m *Manager) RecordAllTypes(playerID string, allTargetTypes []string) *ChallengeDef {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isUnlocked(playerID, ChallengeAllTypes) {
		return nil
	}

	sess, ok := m.sessions[playerID]
	if !ok {
		return nil
	}

	for _, t := range allTargetTypes {
		if !sess.KilledTypes[t] {
			return nil
		}
	}

	if _, ok := m.states[playerID]; !ok {
		m.states[playerID] = make(map[ChallengeID]*PlayerChallengeState)
	}
	m.states[playerID][ChallengeAllTypes] = &PlayerChallengeState{Unlocked: true, UnlockedAt: time.Now()}
	return Challenges[ChallengeAllTypes]
}

// GetSnapshot 取得玩家挑戰快照（用於 API 回傳）
func (m *Manager) GetSnapshot(playerID string) []ChallengeSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	states := m.states[playerID]
	var result []ChallengeSnapshot

	for id, def := range Challenges {
		snap := ChallengeSnapshot{
			ID:          string(id),
			Name:        def.Name,
			Description: def.Description,
			Icon:        def.Icon,
			Reward:      def.Reward,
			IsHidden:    def.IsHidden,
		}
		if s, ok := states[id]; ok && s.Unlocked {
			snap.Unlocked = true
			snap.RewardClaimed = s.RewardClaimed
			snap.UnlockedAt = s.UnlockedAt
			snap.IsHidden = false // 解鎖後不再隱藏
		}
		result = append(result, snap)
	}
	return result
}

// ChallengeSnapshot 挑戰快照（用於 Client 顯示）
type ChallengeSnapshot struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Icon          string    `json:"icon"`
	Reward        int       `json:"reward"`
	IsHidden      bool      `json:"is_hidden"`
	Unlocked      bool      `json:"unlocked"`
	RewardClaimed bool      `json:"reward_claimed"`
	UnlockedAt    time.Time `json:"unlocked_at,omitempty"`
}
