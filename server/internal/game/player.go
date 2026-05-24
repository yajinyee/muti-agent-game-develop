// Package game — 玩家狀態
// server-core-agent 負責維護
package game

import (
	"chiikawa-game/internal/data"
)

const InitialCoins = 10000

// Player 代表一個玩家的遊戲狀態
type Player struct {
	ID             string
	Coins          int
	BetLevel       int
	IsAuto         bool
	LockTargetID   string
	LaborValue     int
	EntryBetCost   int // Bonus 進入時的 BetCost
}

func NewPlayer(id string) *Player {
	return &Player{
		ID:       id,
		Coins:    InitialCoins,
		BetLevel: 1,
	}
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
