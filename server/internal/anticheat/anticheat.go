// Package anticheat 玩家異常行為偵測系統（DAY-105）
// 偵測 bot 攻擊、RTP 異常、金幣暴增等可疑行為
// 觸發警告後記錄到 log，並透過 Admin Dashboard 顯示
package anticheat

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// AlertLevel 警告等級
type AlertLevel string

const (
	AlertInfo     AlertLevel = "info"     // 資訊（記錄但不處理）
	AlertWarning  AlertLevel = "warning"  // 警告（記錄 + 通知）
	AlertCritical AlertLevel = "critical" // 嚴重（記錄 + 通知 + 可能封鎖）
)

// AlertType 警告類型
type AlertType string

const (
	AlertHighRTP        AlertType = "high_rtp"         // RTP 異常偏高
	AlertBotAttack      AlertType = "bot_attack"        // 攻擊頻率異常（疑似 bot）
	AlertCoinSpike      AlertType = "coin_spike"        // 金幣短時間暴增
	AlertBonusAbuse     AlertType = "bonus_abuse"       // Bonus 觸發頻率異常
	AlertJackpotAbuse   AlertType = "jackpot_abuse"     // Jackpot 中獎頻率異常
)

// Alert 警告記錄
type Alert struct {
	ID          string     `json:"id"`
	PlayerID    string     `json:"player_id"`
	DisplayName string     `json:"display_name"`
	Type        AlertType  `json:"type"`
	Level       AlertLevel `json:"level"`
	Message     string     `json:"message"`
	Value       float64    `json:"value"`    // 觸發值（如 RTP 200%）
	Threshold   float64    `json:"threshold"` // 門檻值
	CreatedAt   time.Time  `json:"created_at"`
	Resolved    bool       `json:"resolved"`
}

// PlayerRecord 玩家行為記錄（滑動視窗）
type PlayerRecord struct {
	PlayerID    string
	DisplayName string

	// 攻擊頻率追蹤（最近 10 秒）
	AttackTimes []time.Time

	// RTP 追蹤（最近 1000 次攻擊）
	TotalBet    int64
	TotalReward int64
	AttackCount int64

	// 金幣追蹤（最近 60 秒）
	CoinsHistory []coinSnapshot

	// Bonus 追蹤（最近 10 分鐘）
	BonusTimes []time.Time

	// Jackpot 追蹤（最近 1 小時）
	JackpotTimes []time.Time

	// 最後警告時間（防止重複警告）
	LastAlertAt map[AlertType]time.Time
}

type coinSnapshot struct {
	Coins int64
	At    time.Time
}

// Manager 異常偵測管理器
type Manager struct {
	mu      sync.RWMutex
	records map[string]*PlayerRecord // playerID → record
	alerts  []*Alert                 // 最近 100 條警告
	alertID int
}

// Thresholds 偵測門檻
const (
	ThresholdRTP          = 2.5   // RTP > 250% 觸發警告（需要 > 100 次攻擊）
	ThresholdAttackPerSec = 8.0   // 每秒攻擊 > 8 次觸發警告
	ThresholdCoinSpike    = 50000 // 60 秒內金幣增加 > 50000 觸發警告
	ThresholdBonusPerMin  = 5     // 10 分鐘內 Bonus > 5 次觸發警告
	ThresholdJackpotPerHr = 3     // 1 小時內 Jackpot > 3 次觸發警告
	MinAttacksForRTP      = 100   // RTP 計算需要至少 100 次攻擊
	AlertCooldown         = 5 * time.Minute // 同類型警告冷卻時間
)

// New 建立新的異常偵測管理器
func New() *Manager {
	return &Manager{
		records: make(map[string]*PlayerRecord),
		alerts:  make([]*Alert, 0, 100),
	}
}

// EnsureRecord 確保玩家記錄存在
func (m *Manager) EnsureRecord(playerID, displayName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.records[playerID]; !ok {
		m.records[playerID] = &PlayerRecord{
			PlayerID:    playerID,
			DisplayName: displayName,
			LastAlertAt: make(map[AlertType]time.Time),
		}
	}
}

// RemoveRecord 移除玩家記錄（離線時）
func (m *Manager) RemoveRecord(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.records, playerID)
}

// RecordAttack 記錄攻擊事件，回傳是否觸發警告
func (m *Manager) RecordAttack(playerID string, betCost int) *Alert {
	m.mu.Lock()
	defer m.mu.Unlock()

	rec, ok := m.records[playerID]
	if !ok {
		return nil
	}

	now := time.Now()
	rec.AttackTimes = append(rec.AttackTimes, now)
	rec.TotalBet += int64(betCost)
	rec.AttackCount++

	// 清理 10 秒前的攻擊記錄
	cutoff := now.Add(-10 * time.Second)
	filtered := rec.AttackTimes[:0]
	for _, t := range rec.AttackTimes {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	rec.AttackTimes = filtered

	// 偵測攻擊頻率異常：最近 10 秒內超過 80 次（= 8次/秒）
	// 使用計數而非時間差，避免同一時間點 duration=0 的問題
	if len(rec.AttackTimes) > int(ThresholdAttackPerSec*10) {
		aps := float64(len(rec.AttackTimes)) / 10.0
		return m.createAlert(rec, AlertBotAttack, AlertWarning,
			fmt.Sprintf("攻擊頻率異常：10秒內 %d 次（門檻 %.0f 次/秒）", len(rec.AttackTimes), ThresholdAttackPerSec),
			aps, ThresholdAttackPerSec)
	}
	return nil
}

// RecordReward 記錄獎勵事件，回傳是否觸發 RTP 警告
func (m *Manager) RecordReward(playerID string, reward int) *Alert {
	m.mu.Lock()
	defer m.mu.Unlock()

	rec, ok := m.records[playerID]
	if !ok {
		return nil
	}

	rec.TotalReward += int64(reward)

	// 需要至少 100 次攻擊才計算 RTP
	if rec.AttackCount < MinAttacksForRTP || rec.TotalBet == 0 {
		return nil
	}

	rtp := float64(rec.TotalReward) / float64(rec.TotalBet)
	if rtp > ThresholdRTP {
		return m.createAlert(rec, AlertHighRTP, AlertWarning,
			fmt.Sprintf("RTP 異常：%.1f%%（門檻 %.0f%%，基於 %d 次攻擊）",
				rtp*100, ThresholdRTP*100, rec.AttackCount),
			rtp, ThresholdRTP)
	}
	return nil
}

// RecordCoins 記錄金幣變化，回傳是否觸發暴增警告
func (m *Manager) RecordCoins(playerID string, coins int64) *Alert {
	m.mu.Lock()
	defer m.mu.Unlock()

	rec, ok := m.records[playerID]
	if !ok {
		return nil
	}

	now := time.Now()
	rec.CoinsHistory = append(rec.CoinsHistory, coinSnapshot{Coins: coins, At: now})

	// 清理 60 秒前的記錄
	cutoff := now.Add(-60 * time.Second)
	filtered := rec.CoinsHistory[:0]
	for _, s := range rec.CoinsHistory {
		if s.At.After(cutoff) {
			filtered = append(filtered, s)
		}
	}
	rec.CoinsHistory = filtered

	// 計算 60 秒內金幣增量
	if len(rec.CoinsHistory) >= 2 {
		oldest := rec.CoinsHistory[0].Coins
		newest := rec.CoinsHistory[len(rec.CoinsHistory)-1].Coins
		delta := newest - oldest
		if delta > ThresholdCoinSpike {
			return m.createAlert(rec, AlertCoinSpike, AlertWarning,
				fmt.Sprintf("金幣暴增：60秒內增加 %d 金幣（門檻 %d）", delta, ThresholdCoinSpike),
				float64(delta), ThresholdCoinSpike)
		}
	}
	return nil
}

// RecordBonus 記錄 Bonus 觸發，回傳是否觸發頻率警告
func (m *Manager) RecordBonus(playerID string) *Alert {
	m.mu.Lock()
	defer m.mu.Unlock()

	rec, ok := m.records[playerID]
	if !ok {
		return nil
	}

	now := time.Now()
	rec.BonusTimes = append(rec.BonusTimes, now)

	// 清理 10 分鐘前的記錄
	cutoff := now.Add(-10 * time.Minute)
	filtered := rec.BonusTimes[:0]
	for _, t := range rec.BonusTimes {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	rec.BonusTimes = filtered

	if len(rec.BonusTimes) > ThresholdBonusPerMin {
		return m.createAlert(rec, AlertBonusAbuse, AlertWarning,
			fmt.Sprintf("Bonus 頻率異常：10分鐘內觸發 %d 次（門檻 %d 次）",
				len(rec.BonusTimes), ThresholdBonusPerMin),
			float64(len(rec.BonusTimes)), ThresholdBonusPerMin)
	}
	return nil
}

// RecordJackpot 記錄 Jackpot 中獎，回傳是否觸發頻率警告
func (m *Manager) RecordJackpot(playerID string) *Alert {
	m.mu.Lock()
	defer m.mu.Unlock()

	rec, ok := m.records[playerID]
	if !ok {
		return nil
	}

	now := time.Now()
	rec.JackpotTimes = append(rec.JackpotTimes, now)

	// 清理 1 小時前的記錄
	cutoff := now.Add(-1 * time.Hour)
	filtered := rec.JackpotTimes[:0]
	for _, t := range rec.JackpotTimes {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	rec.JackpotTimes = filtered

	if len(rec.JackpotTimes) > ThresholdJackpotPerHr {
		return m.createAlert(rec, AlertJackpotAbuse, AlertCritical,
			fmt.Sprintf("Jackpot 頻率異常：1小時內中獎 %d 次（門檻 %d 次）",
				len(rec.JackpotTimes), ThresholdJackpotPerHr),
			float64(len(rec.JackpotTimes)), ThresholdJackpotPerHr)
	}
	return nil
}

// createAlert 建立警告（非 thread-safe，需在 mu.Lock 內呼叫）
func (m *Manager) createAlert(rec *PlayerRecord, alertType AlertType, level AlertLevel, msg string, value, threshold float64) *Alert {
	// 冷卻時間檢查
	if last, ok := rec.LastAlertAt[alertType]; ok {
		if time.Since(last) < AlertCooldown {
			return nil
		}
	}
	rec.LastAlertAt[alertType] = time.Now()

	m.alertID++
	alert := &Alert{
		ID:          fmt.Sprintf("alert-%06d", m.alertID),
		PlayerID:    rec.PlayerID,
		DisplayName: rec.DisplayName,
		Type:        alertType,
		Level:       level,
		Message:     msg,
		Value:       value,
		Threshold:   threshold,
		CreatedAt:   time.Now(),
	}

	// 保留最近 100 條警告
	m.alerts = append(m.alerts, alert)
	if len(m.alerts) > 100 {
		m.alerts = m.alerts[1:]
	}

	log.Printf("[AntiCheat] %s %s: %s", level, rec.PlayerID, msg)
	return alert
}

// GetAlerts 取得最近警告列表
func (m *Manager) GetAlerts(limit int) []*Alert {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if limit <= 0 || limit > len(m.alerts) {
		limit = len(m.alerts)
	}
	// 回傳最新的 limit 條（倒序）
	result := make([]*Alert, limit)
	for i := 0; i < limit; i++ {
		result[i] = m.alerts[len(m.alerts)-1-i]
	}
	return result
}

// GetAlertCount 取得警告統計
func (m *Manager) GetAlertCount() (total, critical, warning int) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total = len(m.alerts)
	for _, a := range m.alerts {
		switch a.Level {
		case AlertCritical:
			critical++
		case AlertWarning:
			warning++
		}
	}
	return total, critical, warning
}

// GetPlayerRTP 取得玩家當前 RTP（用於 Admin Dashboard）
func (m *Manager) GetPlayerRTP(playerID string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	rec, ok := m.records[playerID]
	if !ok || rec.TotalBet == 0 {
		return 0
	}
	return float64(rec.TotalReward) / float64(rec.TotalBet)
}
