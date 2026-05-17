// Package player 管理玩家狀態
package player

import (
	"sync"
	"time"

	"digital-twin/server/internal/data"
)

// Player 玩家狀態
type Player struct {
	mu sync.RWMutex

	ID          string
	Coins       int
	BetLevel    int
	LaborValue  int // 勞動值 0-100
	IsAuto      bool
	LockTargetID string // 鎖定目標 InstanceID
	SessionStart time.Time
	LastAttackAt time.Time

	// 統計
	TotalBet    int
	TotalReward int
	AttackCount int
	KillCount   int
}

// NewPlayer 建立新玩家
func NewPlayer(id string, initialCoins int) *Player {
	return &Player{
		ID:           id,
		Coins:        initialCoins,
		BetLevel:     1,
		LaborValue:   0,
		IsAuto:       false,
		SessionStart: time.Now(),
	}
}

// GetBetDef 取得目前投注定義
func (p *Player) GetBetDef() *data.BetDef {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return data.GetBetDef(p.BetLevel)
}

// GetCharacter 取得目前角色
func (p *Player) GetCharacter() *data.CharacterDef {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return data.GetCharacterByBetLevel(p.BetLevel)
}

// CanAttack 是否可以攻擊（金幣足夠）
func (p *Player) CanAttack() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	bet := data.GetBetDef(p.BetLevel)
	return p.Coins >= bet.BetCost
}

// DeductBet 扣除投注金額，回傳是否成功
func (p *Player) DeductBet() (int, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	bet := data.GetBetDef(p.BetLevel)
	if p.Coins < bet.BetCost {
		return 0, false
	}
	p.Coins -= bet.BetCost
	p.TotalBet += bet.BetCost
	p.AttackCount++
	p.LastAttackAt = time.Now()
	return bet.BetCost, true
}

// AddReward 增加獎勵
func (p *Player) AddReward(amount int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Coins += amount
	p.TotalReward += amount
}

// AddLaborValue 增加勞動值，回傳是否觸發 Bonus
func (p *Player) AddLaborValue(amount int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.LaborValue += amount
	if p.LaborValue >= data.LaborValueMax {
		p.LaborValue = data.LaborValueMax
		return true
	}
	return false
}

// ResetLaborValue 重置勞動值（Bonus 觸發後）
func (p *Player) ResetLaborValue() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.LaborValue = 0
}

// SetBetLevel 切換投注等級
func (p *Player) SetBetLevel(level int) bool {
	if level < 1 || level > 10 {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.BetLevel = level
	return true
}

// SetLock 設定鎖定目標
func (p *Player) SetLock(targetID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.LockTargetID = targetID
}

// SetAuto 設定自動攻擊
func (p *Player) SetAuto(enabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.IsAuto = enabled
}

// Snapshot 取得玩家狀態快照（用於傳送給 Client）
func (p *Player) Snapshot() PlayerSnapshot {
	p.mu.RLock()
	defer p.mu.RUnlock()
	char := data.GetCharacterByBetLevel(p.BetLevel)
	bet := data.GetBetDef(p.BetLevel)
	return PlayerSnapshot{
		ID:              p.ID,
		Coins:           p.Coins,
		BetLevel:        p.BetLevel,
		BetCost:         bet.BetCost,
		CharacterID:     char.ID,
		CharacterName:   char.Name,
		LaborValue:      p.LaborValue,
		IsAuto:          p.IsAuto,
		LockTargetID:    p.LockTargetID,
		ProjectileSpeed: bet.ProjectileSpeed,
		FireRate:        bet.FireRate,
	}
}

// PlayerSnapshot 玩家狀態快照
type PlayerSnapshot struct {
	ID              string  `json:"id"`
	Coins           int     `json:"coins"`
	BetLevel        int     `json:"bet_level"`
	BetCost         int     `json:"bet_cost"`
	CharacterID     string  `json:"character_id"`
	CharacterName   string  `json:"character_name"`
	LaborValue      int     `json:"labor_value"`
	IsAuto          bool    `json:"is_auto"`
	LockTargetID    string  `json:"lock_target_id"`
	ProjectileSpeed float64 `json:"projectile_speed"`
	FireRate        float64 `json:"fire_rate"`
}
