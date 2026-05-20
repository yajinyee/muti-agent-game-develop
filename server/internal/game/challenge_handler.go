// challenge_handler.go — 隱藏挑戰系統 handler（DAY-085）
package game

import (
	"log"

	"digital-twin/server/internal/game/challenge"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// ChallengeUnlockedPayload 挑戰解鎖通知（Server → Client）
// 注意：此 payload 直接定義在 handler 中，避免 protocol.go 過大
type ChallengeUnlockedPayload struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Reward      int    `json:"reward"`
	TitleID     string `json:"title_id,omitempty"`
	IsHidden    bool   `json:"was_hidden"` // 是否是隱藏挑戰（解鎖前不知道）
}

// notifyChallengeKill 在擊破目標後檢查挑戰（由 handleKill 呼叫）
func (g *Game) notifyChallengeKill(p *player.Player, defID string, multiplier float64, reward int) {
	if g.Challenge == nil {
		return
	}

	// 記錄擊破事件（速度/倍率挑戰）
	unlocked := g.Challenge.RecordKill(p.ID, defID, multiplier)

	// 記錄金幣（財富挑戰）
	coinUnlocked := g.Challenge.RecordCoins(p.ID, reward)
	unlocked = append(unlocked, coinUnlocked...)

	// 發送解鎖通知
	for _, def := range unlocked {
		g.sendChallengeUnlocked(p, def)
	}

	// 檢查全類型挑戰
	allTypes := []string{"T001", "T002", "T003", "T004", "T005", "T006",
		"T101", "T102", "T103", "T104", "T105"}
	if allDef := g.Challenge.RecordAllTypes(p.ID, allTypes); allDef != nil {
		g.sendChallengeUnlocked(p, allDef)
	}
}

// notifyChallengeStreak 在連擊更新後檢查挑戰（由 notifyStreakKill 呼叫）
func (g *Game) notifyChallengeStreak(p *player.Player, streak int) {
	if g.Challenge == nil {
		return
	}
	unlocked := g.Challenge.RecordStreak(p.ID, streak)
	for _, def := range unlocked {
		g.sendChallengeUnlocked(p, def)
	}
}

// notifyChallengeBoss 在 BOSS 擊殺後檢查挑戰（由 handleBossKill 呼叫）
func (g *Game) notifyChallengeBoss(p *player.Player) {
	if g.Challenge == nil {
		return
	}
	if def := g.Challenge.TryUnlock(p.ID, challenge.ChallengeBossFirst); def != nil {
		g.sendChallengeUnlocked(p, def)
	}
}

// notifyChallengeWheel 在轉盤觸發後檢查挑戰（由 notifyWheelKill 呼叫）
func (g *Game) notifyChallengeWheel(p *player.Player, multiplier float64) {
	if g.Challenge == nil {
		return
	}
	if multiplier >= 100 {
		if def := g.Challenge.TryUnlock(p.ID, challenge.ChallengeWheelMax); def != nil {
			g.sendChallengeUnlocked(p, def)
		}
	}
}

// notifyChallengeJackpot 在 Jackpot 中獎後檢查挑戰（由 jackpot_handler 呼叫）
func (g *Game) notifyChallengeJackpot(playerID string) {
	if g.Challenge == nil {
		return
	}
	g.mu.RLock()
	p := g.Players[playerID]
	g.mu.RUnlock()
	if p == nil {
		return
	}
	if def := g.Challenge.TryUnlock(playerID, challenge.ChallengeJackpot); def != nil {
		g.sendChallengeUnlocked(p, def)
	}
}

// sendChallengeUnlocked 發送挑戰解鎖通知給玩家，並發放獎勵
func (g *Game) sendChallengeUnlocked(p *player.Player, def *challenge.ChallengeDef) {
	if def == nil {
		return
	}

	// 自動發放獎勵（不需要玩家手動領取）
	reward := g.Challenge.ClaimReward(p.ID, def.ID)
	if reward > 0 {
		p.AddCoins(reward)
	}

	// 發送解鎖通知
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgChallengeUnlocked,
		Payload: ChallengeUnlockedPayload{
			ID:          string(def.ID),
			Name:        def.Name,
			Description: def.Description,
			Icon:        def.Icon,
			Reward:      reward,
			TitleID:     def.TitleID,
			IsHidden:    def.IsHidden,
		},
	}); err != nil {
		log.Printf("[Challenge] send unlock error: %v", err)
	}

	log.Printf("[Challenge] player=%s unlocked challenge=%s reward=%d", p.ID, def.ID, reward)
}
