// Package respin 實作 Rapid Respin 系統（DAY-121）
// 業界依據：Reflex Gaming Big Game Fishing Rapid Riches（2026-05-14）
// Rapid Respin 是 2026 年捕魚機最新熱門機制，讓玩家在短時間內連續觸發高倍率獎勵
package respin

import (
	"sync"
	"time"
)

// 觸發機率（依投注等級）
const (
	BaseChance    = 0.04  // 基礎觸發機率 4%（LV1-4）
	MidChance     = 0.06  // 中等觸發機率 6%（LV5-7）
	HighChance    = 0.08  // 高等觸發機率 8%（LV8-10）
	ChainWindow   = 10 * time.Second // 連鎖觸發視窗（10 秒內再次擊破可連鎖）
	MaxChain      = 5     // 最大連鎖次數
	CooldownTime  = 30 * time.Second // 玩家冷卻時間（觸發後 30 秒不再觸發）
)

// ChainMultipliers 連鎖倍率（第 1-5 次）
var ChainMultipliers = []float64{1.0, 1.5, 2.0, 3.0, 5.0}

// Session 單次 Rapid Respin 連鎖 session
type Session struct {
	PlayerID    string
	ChainCount  int       // 當前連鎖次數（0-based，0=第一次）
	StartedAt   time.Time // session 開始時間
	LastRespinAt time.Time // 最後一次 respin 時間
	Active      bool
}

// GetCurrentMult 取得當前連鎖倍率
func (s *Session) GetCurrentMult() float64 {
	if s.ChainCount < 0 {
		return 1.0
	}
	if s.ChainCount >= len(ChainMultipliers) {
		return ChainMultipliers[len(ChainMultipliers)-1]
	}
	return ChainMultipliers[s.ChainCount]
}

// CanChain 判斷是否可以繼續連鎖
func (s *Session) CanChain() bool {
	if !s.Active {
		return false
	}
	if s.ChainCount >= MaxChain-1 {
		return false
	}
	return time.Since(s.LastRespinAt) <= ChainWindow
}

// Manager Rapid Respin 管理器
type Manager struct {
	mu       sync.Mutex
	sessions map[string]*Session // playerID -> session
	cooldowns map[string]time.Time // playerID -> 冷卻結束時間
}

// New 建立新的 Rapid Respin 管理器
func New() *Manager {
	return &Manager{
		sessions:  make(map[string]*Session),
		cooldowns: make(map[string]time.Time),
	}
}

// ShouldTrigger 判斷是否應該觸發 Rapid Respin
// betLevel: 1-10
// 回傳 (shouldTrigger bool, isChain bool, chainCount int)
func (m *Manager) ShouldTrigger(playerID string, betLevel int, randFloat float64) (bool, bool, int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 檢查是否在冷卻中
	if cd, ok := m.cooldowns[playerID]; ok {
		if time.Now().Before(cd) {
			return false, false, 0
		}
	}

	// 檢查是否有進行中的 session（連鎖觸發）
	if sess, ok := m.sessions[playerID]; ok && sess.Active {
		if sess.CanChain() {
			// 連鎖觸發機率更高（基礎 × 2）
			chance := m.getChance(betLevel) * 2.0
			if randFloat < chance {
				sess.ChainCount++
				sess.LastRespinAt = time.Now()
				// 達到最大連鎖，標記結束（但本次仍觸發）
				if sess.ChainCount >= MaxChain-1 {
					sess.Active = false
					m.cooldowns[playerID] = time.Now().Add(CooldownTime)
				}
				return true, true, sess.ChainCount
			}
		} else {
			// 連鎖視窗過期，結束 session，設定冷卻
			sess.Active = false
			m.cooldowns[playerID] = time.Now().Add(CooldownTime)
		}
		// session 已結束（無論是視窗過期還是達到最大連鎖），不走新觸發
		if !sess.Active {
			return false, false, 0
		}
	}

	// 新觸發
	chance := m.getChance(betLevel)
	if randFloat < chance {
		sess := &Session{
			PlayerID:     playerID,
			ChainCount:   0,
			StartedAt:    time.Now(),
			LastRespinAt: time.Now(),
			Active:       true,
		}
		m.sessions[playerID] = sess
		// 設定冷卻（連鎖結束後才冷卻，這裡先不設，在 EndSession 設）
		return true, false, 0
	}

	return false, false, 0
}

// EndSession 結束玩家的 Respin session（連鎖視窗過期或達到最大連鎖）
func (m *Manager) EndSession(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sess, ok := m.sessions[playerID]; ok {
		sess.Active = false
	}
	// 設定冷卻
	m.cooldowns[playerID] = time.Now().Add(CooldownTime)
}

// GetSession 取得玩家當前 session（唯讀）
func (m *Manager) GetSession(playerID string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sess, ok := m.sessions[playerID]; ok && sess.Active {
		return sess
	}
	return nil
}

// RemovePlayer 清理玩家資料（離線時呼叫）
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, playerID)
	delete(m.cooldowns, playerID)
}

// getChance 依投注等級取得觸發機率
func (m *Manager) getChance(betLevel int) float64 {
	switch {
	case betLevel >= 8:
		return HighChance
	case betLevel >= 5:
		return MidChance
	default:
		return BaseChance
	}
}
