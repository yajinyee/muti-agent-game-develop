// Package stats 玩家個人統計系統（DAY-096）
// 追蹤玩家的詳細遊戲統計，供個人統計面板顯示
package stats

import (
	"sync"
	"time"
)

// PlayerStats 玩家詳細統計
type PlayerStats struct {
	mu sync.RWMutex

	// 基礎統計
	TotalSessions  int     `json:"total_sessions"`   // 總遊戲場次
	TotalPlayTime  int64   `json:"total_play_time"`  // 總遊戲時間（秒）
	TotalShots     int     `json:"total_shots"`      // 總射擊次數
	TotalKills     int     `json:"total_kills"`      // 總擊破次數
	TotalBet       int     `json:"total_bet"`        // 總投注金幣
	TotalReward    int     `json:"total_reward"`     // 總獲得金幣
	TotalBonuses   int     `json:"total_bonuses"`    // 觸發 Bonus 次數
	TotalBossKills int     `json:"total_boss_kills"` // 擊殺 BOSS 次數

	// 最佳記錄
	BestMultiplier  float64 `json:"best_multiplier"`   // 最高單次倍率
	BestStreak      int     `json:"best_streak"`       // 最高連擊數
	BestSessionScore int    `json:"best_session_score"` // 單場最高得分
	BestBonusReward int     `json:"best_bonus_reward"` // 單次 Bonus 最高獎勵
	MaxCoins        int     `json:"max_coins"`         // 歷史最高金幣

	// Jackpot 統計
	JackpotWins     int `json:"jackpot_wins"`      // 總 Jackpot 中獎次數
	JackpotMiniWins int `json:"jackpot_mini_wins"` // Mini Jackpot 中獎次數
	JackpotMinorWins int `json:"jackpot_minor_wins"` // Minor Jackpot 中獎次數
	JackpotMajorWins int `json:"jackpot_major_wins"` // Major Jackpot 中獎次數
	JackpotGrandWins int `json:"jackpot_grand_wins"` // Grand Jackpot 中獎次數
	TotalJackpotPayout int `json:"total_jackpot_payout"` // 總 Jackpot 獲得金幣

	// 命中率統計
	HitCount   int `json:"hit_count"`   // 命中次數（攻擊到目標）
	MissCount  int `json:"miss_count"`  // 未命中次數

	// 時間統計
	FirstPlayAt time.Time `json:"first_play_at"` // 首次遊戲時間
	LastPlayAt  time.Time `json:"last_play_at"`  // 最後遊戲時間

	// 當前 Session 統計（不持久化）
	sessionStart time.Time
}

// NewPlayerStats 建立新的玩家統計
func NewPlayerStats() *PlayerStats {
	now := time.Now()
	return &PlayerStats{
		FirstPlayAt: now,
		LastPlayAt:  now,
	}
}

// StartSession 開始新的遊戲場次
func (s *PlayerStats) StartSession() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalSessions++
	s.sessionStart = time.Now()
	s.LastPlayAt = time.Now()
}

// EndSession 結束遊戲場次，記錄時長
func (s *PlayerStats) EndSession() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.sessionStart.IsZero() {
		elapsed := int64(time.Since(s.sessionStart).Seconds())
		s.TotalPlayTime += elapsed
		s.sessionStart = time.Time{}
	}
}

// RecordShot 記錄一次射擊
func (s *PlayerStats) RecordShot(betCost int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalShots++
	s.TotalBet += betCost
}

// RecordKill 記錄一次擊破
func (s *PlayerStats) RecordKill(multiplier float64, reward int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalKills++
	s.TotalReward += reward
	s.HitCount++
	if multiplier > s.BestMultiplier {
		s.BestMultiplier = multiplier
	}
}

// RecordMiss 記錄一次未命中（攻擊但目標未死）
func (s *PlayerStats) RecordMiss() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.MissCount++
}

// RecordStreak 記錄連擊數
func (s *PlayerStats) RecordStreak(streak int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if streak > s.BestStreak {
		s.BestStreak = streak
	}
}

// RecordBonus 記錄 Bonus 觸發
func (s *PlayerStats) RecordBonus(reward int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalBonuses++
	if reward > s.BestBonusReward {
		s.BestBonusReward = reward
	}
}

// RecordBossKill 記錄 BOSS 擊殺
func (s *PlayerStats) RecordBossKill() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalBossKills++
}

// RecordJackpot 記錄 Jackpot 中獎
func (s *PlayerStats) RecordJackpot(level string, amount int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.JackpotWins++
	s.TotalJackpotPayout += amount
	switch level {
	case "mini":
		s.JackpotMiniWins++
	case "minor":
		s.JackpotMinorWins++
	case "major":
		s.JackpotMajorWins++
	case "grand":
		s.JackpotGrandWins++
	}
}

// RecordSessionScore 記錄本場得分（結束時呼叫）
func (s *PlayerStats) RecordSessionScore(score int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if score > s.BestSessionScore {
		s.BestSessionScore = score
	}
}

// UpdateMaxCoins 更新歷史最高金幣
func (s *PlayerStats) UpdateMaxCoins(coins int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if coins > s.MaxCoins {
		s.MaxCoins = coins
	}
}

// GetHitRate 計算命中率（0.0-1.0）
func (s *PlayerStats) GetHitRate() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	total := s.HitCount + s.MissCount
	if total == 0 {
		return 0.0
	}
	return float64(s.HitCount) / float64(total)
}

// GetRTP 計算實際 RTP（TotalReward / TotalBet）
func (s *PlayerStats) GetRTP() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.TotalBet == 0 {
		return 0.0
	}
	return float64(s.TotalReward) / float64(s.TotalBet)
}

// GetCurrentSessionTime 取得當前 Session 已進行時間（秒）
func (s *PlayerStats) GetCurrentSessionTime() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.sessionStart.IsZero() {
		return 0
	}
	return int64(time.Since(s.sessionStart).Seconds())
}

// Snapshot 取得統計快照（用於傳送給 Client）
func (s *PlayerStats) Snapshot() StatsSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hitRate := 0.0
	total := s.HitCount + s.MissCount
	if total > 0 {
		hitRate = float64(s.HitCount) / float64(total)
	}

	rtp := 0.0
	if s.TotalBet > 0 {
		rtp = float64(s.TotalReward) / float64(s.TotalBet)
	}

	currentSessionTime := int64(0)
	if !s.sessionStart.IsZero() {
		currentSessionTime = int64(time.Since(s.sessionStart).Seconds())
	}

	return StatsSnapshot{
		TotalSessions:      s.TotalSessions,
		TotalPlayTimeSec:   s.TotalPlayTime + currentSessionTime,
		TotalShots:         s.TotalShots,
		TotalKills:         s.TotalKills,
		TotalBet:           s.TotalBet,
		TotalReward:        s.TotalReward,
		TotalBonuses:       s.TotalBonuses,
		TotalBossKills:     s.TotalBossKills,
		BestMultiplier:     s.BestMultiplier,
		BestStreak:         s.BestStreak,
		BestSessionScore:   s.BestSessionScore,
		BestBonusReward:    s.BestBonusReward,
		MaxCoins:           s.MaxCoins,
		JackpotWins:        s.JackpotWins,
		JackpotMiniWins:    s.JackpotMiniWins,
		JackpotMinorWins:   s.JackpotMinorWins,
		JackpotMajorWins:   s.JackpotMajorWins,
		JackpotGrandWins:   s.JackpotGrandWins,
		TotalJackpotPayout: s.TotalJackpotPayout,
		HitRate:            hitRate,
		RTP:                rtp,
		FirstPlayAt:        s.FirstPlayAt.UnixMilli(),
		LastPlayAt:         s.LastPlayAt.UnixMilli(),
	}
}

// StatsSnapshot 統計快照（用於 WebSocket 傳送）
type StatsSnapshot struct {
	TotalSessions      int     `json:"total_sessions"`
	TotalPlayTimeSec   int64   `json:"total_play_time_sec"`
	TotalShots         int     `json:"total_shots"`
	TotalKills         int     `json:"total_kills"`
	TotalBet           int     `json:"total_bet"`
	TotalReward        int     `json:"total_reward"`
	TotalBonuses       int     `json:"total_bonuses"`
	TotalBossKills     int     `json:"total_boss_kills"`
	BestMultiplier     float64 `json:"best_multiplier"`
	BestStreak         int     `json:"best_streak"`
	BestSessionScore   int     `json:"best_session_score"`
	BestBonusReward    int     `json:"best_bonus_reward"`
	MaxCoins           int     `json:"max_coins"`
	JackpotWins        int     `json:"jackpot_wins"`
	JackpotMiniWins    int     `json:"jackpot_mini_wins"`
	JackpotMinorWins   int     `json:"jackpot_minor_wins"`
	JackpotMajorWins   int     `json:"jackpot_major_wins"`
	JackpotGrandWins   int     `json:"jackpot_grand_wins"`
	TotalJackpotPayout int     `json:"total_jackpot_payout"`
	HitRate            float64 `json:"hit_rate"`
	RTP                float64 `json:"rtp"`
	FirstPlayAt        int64   `json:"first_play_at_ms"`
	LastPlayAt         int64   `json:"last_play_at_ms"`
}

// LoadState 從持久化資料恢復統計狀態（DAY-098）
func (s *PlayerStats) LoadState(
	totalSessions int, totalPlayTime int64, totalShots int, totalKills int,
	totalBet int, totalReward int, totalBonuses int, totalBossKills int,
	bestMultiplier float64, bestStreak int, bestSession int, bestBonus int, maxCoins int,
	jackpotWins int, jackpotMini int, jackpotMinor int, jackpotMajor int, jackpotGrand int, jackpotPayout int,
	hitCount int, missCount int, firstPlayAt time.Time, lastPlayAt time.Time,
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalSessions = totalSessions
	s.TotalPlayTime = totalPlayTime
	s.TotalShots = totalShots
	s.TotalKills = totalKills
	s.TotalBet = totalBet
	s.TotalReward = totalReward
	s.TotalBonuses = totalBonuses
	s.TotalBossKills = totalBossKills
	s.BestMultiplier = bestMultiplier
	s.BestStreak = bestStreak
	s.BestSessionScore = bestSession
	s.BestBonusReward = bestBonus
	s.MaxCoins = maxCoins
	s.JackpotWins = jackpotWins
	s.JackpotMiniWins = jackpotMini
	s.JackpotMinorWins = jackpotMinor
	s.JackpotMajorWins = jackpotMajor
	s.JackpotGrandWins = jackpotGrand
	s.TotalJackpotPayout = jackpotPayout
	s.HitCount = hitCount
	s.MissCount = missCount
	if !firstPlayAt.IsZero() {
		s.FirstPlayAt = firstPlayAt
	}
	if !lastPlayAt.IsZero() {
		s.LastPlayAt = lastPlayAt
	}
}
