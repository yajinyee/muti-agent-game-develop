// tutorial_handler.go — 新手引導系統（DAY-115）
// 新玩家第一次進入遊戲時，顯示互動式引導
// 引導步驟：1.射擊 2.切換投注 3.觸發Bonus 4.完成（給新手獎勵）
// 業界依據：optikpi.com（2026）確認 onboarding 可提升留存率 50%
package game

import (
	"log"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// TutorialStep 引導步驟
type TutorialStep int

const (
	TutorialStepShoot      TutorialStep = 1 // 步驟1：點擊射擊
	TutorialStepBetChange  TutorialStep = 2 // 步驟2：切換投注等級
	TutorialStepBonus      TutorialStep = 3 // 步驟3：觸發 Bonus（勞動值滿）
	TutorialStepComplete   TutorialStep = 4 // 步驟4：完成引導
)

// TutorialReward 完成新手引導的獎勵
const TutorialReward = 5000 // 5000 金幣

// checkAndSendTutorial 檢查是否需要發送新手引導（由 AddPlayer 呼叫）
func (g *Game) checkAndSendTutorial(p *player.Player) {
	if p.TutorialCompleted {
		return
	}
	// 新玩家（從未攻擊過）才顯示引導
	if p.AttackCount > 0 {
		// 老玩家但未完成引導（可能是舊資料），直接標記完成
		p.TutorialCompleted = true
		return
	}

	log.Printf("[Tutorial] Sending tutorial to new player %s", p.ID)
	g.sendTutorialStep(p, TutorialStepShoot)
}

// sendTutorialStep 發送引導步驟
func (g *Game) sendTutorialStep(p *player.Player, step TutorialStep) {
	var title, desc, highlight, action string
	var arrowX, arrowY float64

	switch step {
	case TutorialStepShoot:
		title = "👆 步驟 1/3：射擊目標"
		desc = "點擊畫面中的目標物來射擊！\n擊破目標可以獲得金幣獎勵。"
		highlight = "game_area"
		action = "shoot"
		arrowX = 640
		arrowY = 360
	case TutorialStepBetChange:
		title = "💰 步驟 2/3：調整投注"
		desc = "點擊 + / - 按鈕調整投注等級。\n投注越高，獎勵越豐厚！"
		highlight = "bet_buttons"
		action = "bet_change"
		arrowX = 200
		arrowY = 680
	case TutorialStepBonus:
		title = "🌿 步驟 3/3：觸發 Bonus"
		desc = "繼續射擊，勞動值滿後自動觸發 Bonus Game！\n或使用 Buy Bonus 直接購買觸發。"
		highlight = "labor_bar"
		action = "bonus"
		arrowX = 640
		arrowY = 30
	case TutorialStepComplete:
		title = "🎉 新手引導完成！"
		desc = "恭喜！你已掌握基本操作。\n獲得新手獎勵：🪙 5000 金幣！"
		highlight = "none"
		action = "complete"
		arrowX = 0
		arrowY = 0
	}

	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgTutorialStep,
		Payload: ws.TutorialStepPayload{
			Step:      int(step),
			TotalStep: 3,
			Title:     title,
			Desc:      desc,
			Highlight: highlight,
			Action:    action,
			ArrowX:    arrowX,
			ArrowY:    arrowY,
		},
	}); err != nil {
		log.Printf("[Tutorial] send step error: %v", err)
	}
}

// handleTutorialAction 處理玩家完成引導步驟（Client → Server）
func (g *Game) handleTutorialAction(p *player.Player, msg *ws.Message) {
	if p.TutorialCompleted {
		return
	}

	var payload ws.TutorialActionPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	switch payload.Action {
	case "shoot_done":
		// 完成射擊步驟，進入投注步驟
		g.sendTutorialStep(p, TutorialStepBetChange)
	case "bet_done":
		// 完成投注步驟，進入 Bonus 步驟
		g.sendTutorialStep(p, TutorialStepBonus)
	case "bonus_done", "skip":
		// 完成 Bonus 步驟或跳過，發放獎勵並完成引導
		g.completeTutorial(p)
	}
}

// completeTutorial 完成新手引導，發放獎勵
func (g *Game) completeTutorial(p *player.Player) {
	if p.TutorialCompleted {
		return
	}

	p.TutorialCompleted = true
	p.AddCoins(TutorialReward)

	log.Printf("[Tutorial] Player %s completed tutorial, reward=%d", p.ID, TutorialReward)

	// 發送完成通知
	g.sendTutorialStep(p, TutorialStepComplete)

	// 發送獎勵通知
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgReward,
		Payload: ws.RewardPayload{
			Source:     "tutorial",
			Amount:     TutorialReward,
			Multiplier: 1.0,
			NewBalance: p.GetCoins(),
		},
	})

	// 動態牆：新手完成引導（common 成就事件）
	go g.notifyFeedAchievement(p, "完成新手引導", "🎓", "common")
}

// notifyTutorialShoot 在玩家第一次射擊後推進引導（由 handleAttack 呼叫）
func (g *Game) notifyTutorialShoot(p *player.Player) {
	if p.TutorialCompleted || p.AttackCount != 1 {
		return
	}
	// 第一次射擊完成，推進到投注步驟
	g.sendTutorialStep(p, TutorialStepBetChange)
}

// notifyTutorialBetChange 在玩家切換投注後推進引導（由 handleBetChange 呼叫）
func (g *Game) notifyTutorialBetChange(p *player.Player) {
	if p.TutorialCompleted {
		return
	}
	// 切換投注完成，推進到 Bonus 步驟
	g.sendTutorialStep(p, TutorialStepBonus)
}

// notifyTutorialBonus 在 Bonus 觸發後完成引導（由 triggerBonusReady 呼叫）
func (g *Game) notifyTutorialBonus(p *player.Player) {
	if p.TutorialCompleted {
		return
	}
	g.completeTutorial(p)
}
