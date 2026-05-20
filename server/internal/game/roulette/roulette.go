// Package roulette 雙層倍率輪盤系統（DAY-113）
// 參考 JILI Royal Fishing 的 ChainLong King 雙層輪盤機制
// 內圈 × 外圈 = 最終倍率，最高 1000x
package roulette

import (
	"math/rand"
	"sync"
	"time"
)

// Segment 輪盤格子
type Segment struct {
	Multiplier float64 `json:"multiplier"`
	Label      string  `json:"label"`
	Color      string  `json:"color"`
	Weight     int     `json:"weight"`
}

// InnerSegments 內圈格子（8格，1x-10x）
var InnerSegments = []Segment{
	{Multiplier: 1, Label: "1x", Color: "#9E9E9E", Weight: 30},
	{Multiplier: 2, Label: "2x", Color: "#4CAF50", Weight: 25},
	{Multiplier: 3, Label: "3x", Color: "#8BC34A", Weight: 18},
	{Multiplier: 4, Label: "4x", Color: "#FFC107", Weight: 12},
	{Multiplier: 5, Label: "5x", Color: "#FF9800", Weight: 8},
	{Multiplier: 6, Label: "6x", Color: "#FF5722", Weight: 4},
	{Multiplier: 8, Label: "8x", Color: "#E91E63", Weight: 2},
	{Multiplier: 10, Label: "10x", Color: "#FFD700", Weight: 1},
}

// OuterSegments 外圈格子（12格，1x-100x）
var OuterSegments = []Segment{
	{Multiplier: 1, Label: "1x", Color: "#9E9E9E", Weight: 25},
	{Multiplier: 2, Label: "2x", Color: "#4CAF50", Weight: 20},
	{Multiplier: 3, Label: "3x", Color: "#8BC34A", Weight: 15},
	{Multiplier: 5, Label: "5x", Color: "#FFC107", Weight: 12},
	{Multiplier: 8, Label: "8x", Color: "#FF9800", Weight: 8},
	{Multiplier: 10, Label: "10x", Color: "#FF5722", Weight: 7},
	{Multiplier: 15, Label: "15x", Color: "#E91E63", Weight: 5},
	{Multiplier: 20, Label: "20x", Color: "#9C27B0", Weight: 4},
	{Multiplier: 30, Label: "30x", Color: "#673AB7", Weight: 2},
	{Multiplier: 50, Label: "50x", Color: "#3F51B5", Weight: 1},
	{Multiplier: 80, Label: "80x", Color: "#2196F3", Weight: 1},
	{Multiplier: 100, Label: "100x", Color: "#FFD700", Weight: 1},
}

// TriggerCondition 觸發條件
type TriggerCondition struct {
	DefID    string  // 目標 DefID
	Chance   float64 // 觸發機率
	MinBet   int     // 最低投注等級（0=無限制）
}

// TriggerConditions 各目標觸發雙層輪盤的條件
var TriggerConditions = []TriggerCondition{
	{DefID: "B001", Chance: 1.0, MinBet: 0},  // BOSS 必定觸發
	{DefID: "T103", Chance: 0.05, MinBet: 5}, // 流星 5% 機率（高投注）
	{DefID: "T104", Chance: 0.08, MinBet: 5}, // 金草 8% 機率（高投注）
	{DefID: "T105", Chance: 0.10, MinBet: 7}, // 金幣魚 10% 機率（高投注）
}

// SpinResult 單次旋轉結果
type SpinResult struct {
	SegmentIndex int     `json:"segment_index"`
	Multiplier   float64 `json:"multiplier"`
	Label        string  `json:"label"`
	Color        string  `json:"color"`
}

// RouletteResult 雙層輪盤完整結果
type RouletteResult struct {
	Inner         SpinResult `json:"inner"`
	Outer         SpinResult `json:"outer"`
	FinalMult     float64    `json:"final_mult"`     // 內圈 × 外圈
	BaseReward    int        `json:"base_reward"`
	FinalReward   int        `json:"final_reward"`
	IsJackpot     bool       `json:"is_jackpot"`     // 最終倍率 >= 500x
	IsMegaWin     bool       `json:"is_mega_win"`    // 最終倍率 >= 100x
}

// Session 一次輪盤遊戲的狀態
type Session struct {
	ID          string     `json:"id"`
	PlayerID    string     `json:"player_id"`
	TargetDefID string     `json:"target_def_id"`
	BaseReward  int        `json:"base_reward"`
	StartedAt   time.Time  `json:"started_at"`
	Result      *RouletteResult `json:"result,omitempty"`
	Resolved    bool       `json:"resolved"`
}

var (
	innerTotalWeight int
	outerTotalWeight int
)

func init() {
	for _, s := range InnerSegments {
		innerTotalWeight += s.Weight
	}
	for _, s := range OuterSegments {
		outerTotalWeight += s.Weight
	}
}

// Manager 雙層輪盤管理器
type Manager struct {
	mu       sync.Mutex
	rng      *rand.Rand
	sessions map[string]*Session // playerID → active session
}

// NewManager 建立新管理器
func NewManager() *Manager {
	return &Manager{
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
		sessions: make(map[string]*Session),
	}
}

// ShouldTrigger 判斷是否觸發雙層輪盤
func (m *Manager) ShouldTrigger(defID string, betLevel int) bool {
	for _, cond := range TriggerConditions {
		if cond.DefID == defID {
			if betLevel < cond.MinBet {
				return false
			}
			m.mu.Lock()
			defer m.mu.Unlock()
			return m.rng.Float64() < cond.Chance
		}
	}
	return false
}

// StartSession 開始一次輪盤遊戲，回傳 session ID
func (m *Manager) StartSession(playerID, targetDefID string, baseReward int) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	sessionID := generateID()
	session := &Session{
		ID:          sessionID,
		PlayerID:    playerID,
		TargetDefID: targetDefID,
		BaseReward:  baseReward,
		StartedAt:   time.Now(),
	}
	m.sessions[playerID] = session
	return session
}

// Resolve 執行輪盤旋轉並計算結果
func (m *Manager) Resolve(playerID string) (*RouletteResult, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[playerID]
	if !ok || session.Resolved {
		return nil, false
	}

	// 旋轉內圈
	inner := m.spinInner()
	// 旋轉外圈
	outer := m.spinOuter()

	finalMult := inner.Multiplier * outer.Multiplier
	finalReward := int(float64(session.BaseReward) * finalMult)

	result := &RouletteResult{
		Inner:       inner,
		Outer:       outer,
		FinalMult:   finalMult,
		BaseReward:  session.BaseReward,
		FinalReward: finalReward,
		IsJackpot:   finalMult >= 500,
		IsMegaWin:   finalMult >= 100,
	}

	session.Result = result
	session.Resolved = true
	delete(m.sessions, playerID)

	return result, true
}

// GetSession 取得玩家當前的輪盤 session
func (m *Manager) GetSession(playerID string) (*Session, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[playerID]
	return s, ok
}

// HasActiveSession 玩家是否有進行中的輪盤
func (m *Manager) HasActiveSession(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[playerID]
	return ok && !s.Resolved
}

// CancelSession 取消輪盤（玩家離線時）
func (m *Manager) CancelSession(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, playerID)
}

// spinInner 旋轉內圈
func (m *Manager) spinInner() SpinResult {
	r := m.rng.Intn(innerTotalWeight)
	cumulative := 0
	for i, s := range InnerSegments {
		cumulative += s.Weight
		if r < cumulative {
			return SpinResult{
				SegmentIndex: i,
				Multiplier:   s.Multiplier,
				Label:        s.Label,
				Color:        s.Color,
			}
		}
	}
	return SpinResult{SegmentIndex: 0, Multiplier: InnerSegments[0].Multiplier, Label: InnerSegments[0].Label, Color: InnerSegments[0].Color}
}

// spinOuter 旋轉外圈
func (m *Manager) spinOuter() SpinResult {
	r := m.rng.Intn(outerTotalWeight)
	cumulative := 0
	for i, s := range OuterSegments {
		cumulative += s.Weight
		if r < cumulative {
			return SpinResult{
				SegmentIndex: i,
				Multiplier:   s.Multiplier,
				Label:        s.Label,
				Color:        s.Color,
			}
		}
	}
	return SpinResult{SegmentIndex: 0, Multiplier: OuterSegments[0].Multiplier, Label: OuterSegments[0].Label, Color: OuterSegments[0].Color}
}

// generateID 生成唯一 ID
func generateID() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = chars[r.Intn(len(chars))]
	}
	return string(b)
}

// GetInnerSegments 取得內圈格子定義（供 Client 顯示用）
func GetInnerSegments() []Segment {
	return InnerSegments
}

// GetOuterSegments 取得外圈格子定義（供 Client 顯示用）
func GetOuterSegments() []Segment {
	return OuterSegments
}

// ExpectedRTP 計算雙層輪盤的期望倍率（用於 RTP 驗證）
func ExpectedRTP() float64 {
	var innerExpected float64
	for _, s := range InnerSegments {
		innerExpected += s.Multiplier * float64(s.Weight) / float64(innerTotalWeight)
	}
	var outerExpected float64
	for _, s := range OuterSegments {
		outerExpected += s.Multiplier * float64(s.Weight) / float64(outerTotalWeight)
	}
	return innerExpected * outerExpected
}
