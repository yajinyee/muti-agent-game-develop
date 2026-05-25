// Package game — 玩家狀態
// server-core-agent 負責維護
package game

import (
	"time"

	"chiikawa-game/internal/data"
)

const InitialCoins = 10000

// ComboLevel Combo 等級定義
type ComboLevel struct {
	Hits     int
	MultBonus float64
	Name     string
}

var ComboLevels = []ComboLevel{
	{5,  0.1, "COMBO x5"},
	{10, 0.2, "COMBO x10"},
	{20, 0.5, "COMBO x20"},
	{30, 1.0, "COMBO x30 MAX"},
}

// Player 代表一個玩家的遊戲狀態
type Player struct {
	ID             string
	DisplayName    string // 玩家顯示名稱（可自訂）
	Coins          int
	BetLevel       int
	IsAuto         bool
	LockTargetID   string
	LaborValue     int
	EntryBetCost   int // Bonus 進入時的 BetCost
	// Combo 系統
	ComboCount     int
	LastHitTime    time.Time
	ComboMultBonus float64
}

func NewPlayer(id string) *Player {
	// 預設顯示名稱：ID 前 8 碼
	displayName := id
	if len(id) > 8 {
		displayName = id[:8]
	}
	return &Player{
		ID:          id,
		DisplayName: displayName,
		Coins:       InitialCoins,
		BetLevel:    1,
	}
}

// GetDisplayName 取得玩家顯示名稱
func (p *Player) GetDisplayName() string {
	if p.DisplayName != "" {
		return p.DisplayName
	}
	if len(p.ID) > 8 {
		return p.ID[:8]
	}
	return p.ID
}

func (p *Player) GetBetDef() data.BetLevel {
	return data.GetBetLevel(p.BetLevel)
}

func (p *Player) GetCharacterID() string {
	return p.GetBetDef().CharacterID
}

func (p *Player) GetCharacterName() string {
	switch p.GetCharacterID() {
	case "hachiware":
		return "Hachiware"
	case "usagi":
		return "Usagi"
	default:
		return "Chiikawa"
	}
}

func (p *Player) AddCoins(amount int) {
	p.Coins += amount
}

func (p *Player) SpendCoins(amount int) bool {
	if p.Coins < amount {
		return false
	}
	p.Coins -= amount
	return true
}

func (p *Player) AddLabor(amount int) bool {
	p.LaborValue += amount
	if p.LaborValue > 100 {
		p.LaborValue = 100
	}
	return p.LaborValue >= 100
}

func (p *Player) ResetLabor() {
	p.LaborValue = 0
}

// AddCombo 增加 Combo 計數，回傳是否達到新 Combo 等級
func (p *Player) AddCombo() (newLevel bool, levelName string, multBonus float64) {
	const comboTimeout = 3.0 // 3 秒內沒有命中則重置
	now := time.Now()
	if !p.LastHitTime.IsZero() && now.Sub(p.LastHitTime).Seconds() > comboTimeout {
		p.ComboCount = 0
		p.ComboMultBonus = 0
	}
	p.ComboCount++
	p.LastHitTime = now

	// 檢查是否達到新等級
	for i := len(ComboLevels) - 1; i >= 0; i-- {
		lvl := ComboLevels[i]
		if p.ComboCount == lvl.Hits {
			p.ComboMultBonus = lvl.MultBonus
			return true, lvl.Name, lvl.MultBonus
		}
	}
	// 更新當前 Combo 加成
	p.ComboMultBonus = 0
	for _, lvl := range ComboLevels {
		if p.ComboCount >= lvl.Hits {
			p.ComboMultBonus = lvl.MultBonus
		}
	}
	return false, "", p.ComboMultBonus
}

// ResetCombo 重置 Combo
func (p *Player) ResetCombo() {
	p.ComboCount = 0
	p.ComboMultBonus = 0
}

// GetComboMultBonus 取得當前 Combo 倍率加成
func (p *Player) GetComboMultBonus() float64 {
	const comboTimeout = 3.0
	if !p.LastHitTime.IsZero() && time.Since(p.LastHitTime).Seconds() > comboTimeout {
		p.ComboCount = 0
		p.ComboMultBonus = 0
	}
	return p.ComboMultBonus
}
