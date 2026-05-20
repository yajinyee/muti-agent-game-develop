// Package referral 推薦碼系統（DAY-082）
// 玩家可以生成推薦碼，邀請新玩家加入，雙方各得獎勵
package referral

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ReferralRecord 推薦記錄
type ReferralRecord struct {
	ReferrerID  string    `json:"referrer_id"`
	RefereeID   string    `json:"referee_id"`
	Code        string    `json:"code"`
	CreatedAt   time.Time `json:"created_at"`
	Rewarded    bool      `json:"rewarded"`    // 是否已發放獎勵
	RewardedAt  time.Time `json:"rewarded_at,omitempty"`
}

// PlayerReferralInfo 玩家推薦資訊
type PlayerReferralInfo struct {
	PlayerID      string    `json:"player_id"`
	MyCode        string    `json:"my_code"`        // 我的推薦碼
	UsedCode      string    `json:"used_code"`      // 我使用的推薦碼（空=未使用）
	ReferredBy    string    `json:"referred_by"`    // 推薦我的玩家 ID
	ReferralCount int       `json:"referral_count"` // 我成功推薦的人數
	TotalReward   int       `json:"total_reward"`   // 累計推薦獎勵
}

// Manager 推薦碼管理器
type Manager struct {
	mu sync.RWMutex

	// code -> referrerID（推薦碼對應的推薦人）
	codes map[string]string

	// playerID -> PlayerReferralInfo
	players map[string]*PlayerReferralInfo

	// 推薦記錄列表
	records []*ReferralRecord
}

// 獎勵設定
const (
	ReferrerReward = 1000 // 推薦人獎勵（每成功推薦一人）
	RefereeReward  = 500  // 被推薦人獎勵（首次使用推薦碼）
	MaxReferrals   = 20   // 每人最多推薦人數（防止濫用）
)

// NewManager 建立推薦碼管理器
func NewManager() *Manager {
	return &Manager{
		codes:   make(map[string]string),
		players: make(map[string]*PlayerReferralInfo),
		records: make([]*ReferralRecord, 0),
	}
}

// GetOrCreateCode 取得或建立玩家的推薦碼
func (m *Manager) GetOrCreateCode(playerID string) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	info := m.getOrCreateInfo(playerID)
	if info.MyCode != "" {
		return info.MyCode
	}

	// 生成唯一推薦碼（6位英數字）
	code := m.generateUniqueCode()
	info.MyCode = code
	m.codes[code] = playerID
	return code
}

// generateUniqueCode 生成唯一推薦碼（需在鎖內呼叫）
func (m *Manager) generateUniqueCode() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // 去掉容易混淆的字元
	for {
		code := make([]byte, 6)
		for i := range code {
			code[i] = chars[rand.Intn(len(chars))]
		}
		codeStr := string(code)
		if _, exists := m.codes[codeStr]; !exists {
			return codeStr
		}
	}
}

// UseCode 使用推薦碼，回傳 (referrerID, error)
// 條件：玩家未使用過推薦碼，推薦碼有效，不能推薦自己
func (m *Manager) UseCode(playerID string, code string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	info := m.getOrCreateInfo(playerID)

	// 已使用過推薦碼
	if info.UsedCode != "" {
		return "", fmt.Errorf("already used a referral code")
	}

	// 推薦碼不存在
	referrerID, ok := m.codes[code]
	if !ok {
		return "", fmt.Errorf("invalid referral code")
	}

	// 不能推薦自己
	if referrerID == playerID {
		return "", fmt.Errorf("cannot use your own referral code")
	}

	// 推薦人已達上限
	referrerInfo := m.getOrCreateInfo(referrerID)
	if referrerInfo.ReferralCount >= MaxReferrals {
		return "", fmt.Errorf("referrer has reached maximum referrals")
	}

	// 記錄使用
	info.UsedCode = code
	info.ReferredBy = referrerID

	// 更新推薦人統計
	referrerInfo.ReferralCount++
	referrerInfo.TotalReward += ReferrerReward

	// 更新被推薦人統計
	info.TotalReward += RefereeReward

	// 建立推薦記錄
	record := &ReferralRecord{
		ReferrerID: referrerID,
		RefereeID:  playerID,
		Code:       code,
		CreatedAt:  time.Now(),
		Rewarded:   true,
		RewardedAt: time.Now(),
	}
	m.records = append(m.records, record)

	return referrerID, nil
}

// GetInfo 取得玩家推薦資訊
func (m *Manager) GetInfo(playerID string) *PlayerReferralInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	info := m.getOrCreateInfo(playerID)
	// 確保有推薦碼
	if info.MyCode == "" {
		code := m.generateUniqueCode()
		info.MyCode = code
		m.codes[code] = playerID
	}

	// 複製一份
	copy := *info
	return &copy
}

// GetRecentRecords 取得最近的推薦記錄（最多 10 筆）
func (m *Manager) GetRecentRecords() []*ReferralRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()

	n := len(m.records)
	if n > 10 {
		n = 10
	}
	result := make([]*ReferralRecord, n)
	copy(result, m.records[len(m.records)-n:])
	return result
}

// getOrCreateInfo 取得或建立玩家推薦資訊（需在鎖內呼叫）
func (m *Manager) getOrCreateInfo(playerID string) *PlayerReferralInfo {
	if info, ok := m.players[playerID]; ok {
		return info
	}
	info := &PlayerReferralInfo{
		PlayerID: playerID,
	}
	m.players[playerID] = info
	return info
}

// GetStats 取得全局推薦統計
func (m *Manager) GetStats() (totalCodes int, totalReferrals int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.codes), len(m.records)
}
