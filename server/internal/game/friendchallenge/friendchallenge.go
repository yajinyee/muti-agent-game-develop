// Package friendchallenge 好友挑戰系統（DAY-102）
// 好友間 1v1 挑戰：3 分鐘內比較分數，勝者獲得全部賭注
// 業界依據：ourculturemag.com 2026-03 確認 group competitions 是 2026 年 iGaming 核心趨勢
package friendchallenge

import (
	"sync"
	"time"
)

// ChallengeStatus 挑戰狀態
type ChallengeStatus string

const (
	StatusPending   ChallengeStatus = "pending"   // 等待對方接受
	StatusActive    ChallengeStatus = "active"    // 進行中
	StatusCompleted ChallengeStatus = "completed" // 已完成
	StatusDeclined  ChallengeStatus = "declined"  // 已拒絕
	StatusExpired   ChallengeStatus = "expired"   // 已過期（未接受）
)

const (
	ChallengeDuration = 3 * time.Minute // 挑戰持續時間
	ChallengeStake    = 1000            // 每人賭注金幣
	PendingTimeout    = 30 * time.Second // 等待接受超時
)

// Challenge 挑戰記錄
type Challenge struct {
	ID          string          `json:"id"`
	ChallengerID string         `json:"challenger_id"`
	ChallengerName string       `json:"challenger_name"`
	ChallengedID string         `json:"challenged_id"`
	ChallengedName string       `json:"challenged_name"`
	Status      ChallengeStatus `json:"status"`
	Stake       int             `json:"stake"`        // 每人賭注
	StartAt     time.Time       `json:"start_at"`
	EndAt       time.Time       `json:"end_at"`
	CreatedAt   time.Time       `json:"created_at"`

	// 分數（進行中/完成後）
	ChallengerScore int `json:"challenger_score"`
	ChallengedScore int `json:"challenged_score"`

	// 結果
	WinnerID    string `json:"winner_id,omitempty"`
	WinnerName  string `json:"winner_name,omitempty"`
	Prize       int    `json:"prize,omitempty"` // 勝者獲得的金幣
}

// IsExpired 是否已過期
func (c *Challenge) IsExpired() bool {
	if c.Status == StatusPending {
		return time.Now().After(c.CreatedAt.Add(PendingTimeout))
	}
	if c.Status == StatusActive {
		return time.Now().After(c.EndAt)
	}
	return false
}

// TimeRemaining 剩餘時間（秒）
func (c *Challenge) TimeRemaining() int {
	if c.Status != StatusActive {
		return 0
	}
	remaining := time.Until(c.EndAt)
	if remaining < 0 {
		return 0
	}
	return int(remaining.Seconds())
}

// Manager 好友挑戰管理器
type Manager struct {
	mu         sync.RWMutex
	challenges map[string]*Challenge // challengeID → Challenge
	// playerID → 當前進行中的挑戰 ID（每人同時只能有一個）
	activeChallenges map[string]string
}

// New 建立新的挑戰管理器
func New() *Manager {
	return &Manager{
		challenges:       make(map[string]*Challenge),
		activeChallenges: make(map[string]string),
	}
}

// CreateChallenge 發起挑戰
// 回傳 (challenge, error_code)
func (m *Manager) CreateChallenge(challengerID, challengerName, challengedID, challengedName string) (*Challenge, string) {
	if challengerID == challengedID {
		return nil, "self_challenge"
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 檢查是否已有進行中的挑戰
	if _, ok := m.activeChallenges[challengerID]; ok {
		return nil, "already_in_challenge"
	}
	if _, ok := m.activeChallenges[challengedID]; ok {
		return nil, "opponent_in_challenge"
	}

	id := challengerID + ":" + challengedID + ":" + time.Now().Format("150405")
	c := &Challenge{
		ID:             id,
		ChallengerID:   challengerID,
		ChallengerName: challengerName,
		ChallengedID:   challengedID,
		ChallengedName: challengedName,
		Status:         StatusPending,
		Stake:          ChallengeStake,
		CreatedAt:      time.Now(),
	}
	m.challenges[id] = c
	return c, ""
}

// AcceptChallenge 接受挑戰
func (m *Manager) AcceptChallenge(challengeID, playerID string) (*Challenge, string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	c, ok := m.challenges[challengeID]
	if !ok {
		return nil, "not_found"
	}
	if c.ChallengedID != playerID {
		return nil, "not_your_challenge"
	}
	if c.Status != StatusPending {
		return nil, "invalid_status"
	}
	if c.IsExpired() {
		c.Status = StatusExpired
		return nil, "expired"
	}

	// 開始挑戰
	now := time.Now()
	c.Status = StatusActive
	c.StartAt = now
	c.EndAt = now.Add(ChallengeDuration)

	// 標記雙方為進行中
	m.activeChallenges[c.ChallengerID] = challengeID
	m.activeChallenges[c.ChallengedID] = challengeID

	return c, ""
}

// DeclineChallenge 拒絕挑戰
func (m *Manager) DeclineChallenge(challengeID, playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	c, ok := m.challenges[challengeID]
	if !ok || c.ChallengedID != playerID || c.Status != StatusPending {
		return false
	}
	c.Status = StatusDeclined
	return true
}

// AddScore 增加玩家分數（擊破目標時呼叫）
func (m *Manager) AddScore(playerID string, points int) *Challenge {
	m.mu.Lock()
	defer m.mu.Unlock()

	challengeID, ok := m.activeChallenges[playerID]
	if !ok {
		return nil
	}
	c, ok := m.challenges[challengeID]
	if !ok || c.Status != StatusActive {
		return nil
	}

	if playerID == c.ChallengerID {
		c.ChallengerScore += points
	} else if playerID == c.ChallengedID {
		c.ChallengedScore += points
	}
	return c
}

// GetActiveChallenge 取得玩家當前進行中的挑戰
func (m *Manager) GetActiveChallenge(playerID string) *Challenge {
	m.mu.RLock()
	defer m.mu.RUnlock()

	challengeID, ok := m.activeChallenges[playerID]
	if !ok {
		return nil
	}
	return m.challenges[challengeID]
}

// CheckAndFinish 檢查並結算到期的挑戰
// 回傳已結算的挑戰列表
func (m *Manager) CheckAndFinish() []*Challenge {
	m.mu.Lock()
	defer m.mu.Unlock()

	var finished []*Challenge
	for _, c := range m.challenges {
		if c.Status == StatusActive && time.Now().After(c.EndAt) {
			m.finishChallenge(c)
			finished = append(finished, c)
		}
		// 清理過期的 pending 挑戰
		if c.Status == StatusPending && c.IsExpired() {
			c.Status = StatusExpired
		}
	}
	return finished
}

// finishChallenge 結算挑戰（非 thread-safe，呼叫前需持有鎖）
func (m *Manager) finishChallenge(c *Challenge) {
	c.Status = StatusCompleted

	// 決定勝者
	if c.ChallengerScore > c.ChallengedScore {
		c.WinnerID = c.ChallengerID
		c.WinnerName = c.ChallengerName
		c.Prize = c.Stake * 2
	} else if c.ChallengedScore > c.ChallengerScore {
		c.WinnerID = c.ChallengedID
		c.WinnerName = c.ChallengedName
		c.Prize = c.Stake * 2
	} else {
		// 平局：各退回賭注
		c.WinnerID = ""
		c.WinnerName = "平局"
		c.Prize = c.Stake // 各退回
	}

	// 清除進行中標記
	delete(m.activeChallenges, c.ChallengerID)
	delete(m.activeChallenges, c.ChallengedID)
}

// ForceFinish 強制結算指定玩家的挑戰（玩家離線時呼叫）
func (m *Manager) ForceFinish(playerID string) *Challenge {
	m.mu.Lock()
	defer m.mu.Unlock()

	challengeID, ok := m.activeChallenges[playerID]
	if !ok {
		return nil
	}
	c, ok := m.challenges[challengeID]
	if !ok || c.Status != StatusActive {
		return nil
	}

	// 離線玩家視為棄賽（對方勝）
	if playerID == c.ChallengerID {
		c.WinnerID = c.ChallengedID
		c.WinnerName = c.ChallengedName
	} else {
		c.WinnerID = c.ChallengerID
		c.WinnerName = c.ChallengerName
	}
	c.Status = StatusCompleted
	c.Prize = c.Stake * 2

	delete(m.activeChallenges, c.ChallengerID)
	delete(m.activeChallenges, c.ChallengedID)
	return c
}

// GetChallengeByID 取得挑戰記錄
func (m *Manager) GetChallengeByID(challengeID string) *Challenge {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.challenges[challengeID]
}

// IsInChallenge 檢查玩家是否在挑戰中
func (m *Manager) IsInChallenge(playerID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.activeChallenges[playerID]
	return ok
}
