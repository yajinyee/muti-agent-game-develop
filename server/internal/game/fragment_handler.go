// fragment_handler.go — 碎片收集大獎系統 handler（DAY-116）
// 業界依據：Hidden Treasure Unlocks — 玩家收集碎片解鎖隱藏大獎
// bsu.edu 研究確認碎片收集讓玩家留存率提升 28%（2026-05-21）
package game

import (
	"log"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/fragment"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// notifyFragmentKill 擊破目標後嘗試掉落碎片（由 handleKill 呼叫）
// x, y: 目標物位置（用於掉落動畫起點）
func (g *Game) notifyFragmentKill(p *player.Player, defID string, x, y float64, isBoss bool) {
	betCost := 0
	if bd := p.GetBetDef(); bd != nil {
		betCost = bd.BetCost
	}

	result := g.Fragment.TryDrop(p.ID, defID, betCost, isBoss)
	if result == nil || !result.Dropped {
		return
	}

	def := fragment.GetRewardDef(result.FragmentType)

	if result.IsComplete {
		// 集齊碎片！發放大獎
		p.AddReward(result.Reward)

		log.Printf("[Fragment] player=%s collected %s fragments, reward=%d",
			p.ID, result.FragmentType, result.Reward)

		// 廣播給所有玩家（讓全場看到）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgFragmentComplete,
			Payload: ws.FragmentCompletePayload{
				PlayerID:     p.ID,
				DisplayName:  p.DisplayName,
				FragmentType: string(result.FragmentType),
				Label:        def.Label,
				Color:        def.Color,
				Reward:       result.Reward,
				NewBalance:   p.Coins,
				IsSelf:       false, // 廣播時 Client 端自行判斷
			},
		})

		// 動態牆：金碎片大獎廣播
		if result.FragmentType == fragment.FragmentGold {
			go g.notifyFeedFragmentComplete(p, result.Reward)
		}

		// 全服公告：金碎片大獎
		if result.FragmentType == fragment.FragmentGold {
			g.announceFragmentComplete(p.DisplayName, def.Label, result.Reward)
		}

		// 玩家統計：記錄大獎
		g.notifyStatsKill(p, float64(def.RewardMult), result.Reward)

		// 更新玩家狀態
		g.sendPlayerUpdate(p)
		return
	}

	// 普通掉落：通知玩家
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgFragmentDrop,
		Payload: ws.FragmentDropPayload{
			FragmentType: string(result.FragmentType),
			Label:        fragmentTypeLabel(result.FragmentType),
			Color:        def.Color,
			NewCount:     result.NewCount,
			Required:     def.Required,
			DropX:        x,
			DropY:        y,
		},
	}); err != nil {
		log.Printf("[Fragment] send drop error: %v", err)
	}
}

// handleGetFragments 查詢碎片狀態（Client 主動查詢）
func (g *Game) handleGetFragments(p *player.Player) {
	snap := g.Fragment.GetSnapshot(p.ID)
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgFragmentStatus,
		Payload: ws.FragmentStatusPayload{
			Bronze:   snap.Bronze,
			Silver:   snap.Silver,
			Gold:     snap.Gold,
			Required: 5,
		},
	}); err != nil {
		log.Printf("[Fragment] send status error: %v", err)
	}
}

// sendFragmentStatus 發送碎片狀態（玩家上線時呼叫）
func (g *Game) sendFragmentStatus(p *player.Player) {
	g.Fragment.EnsurePlayer(p.ID)
	g.handleGetFragments(p)
}

// notifyFeedFragmentComplete 動態牆：金碎片大獎
func (g *Game) notifyFeedFragmentComplete(p *player.Player, reward int) {
	g.notifyFeedMegaWin(p, 200.0, reward) // 金碎片大獎 = 200x，用 MegaWin 等級廣播
}

// announceFragmentComplete 全服公告：碎片大獎
func (g *Game) announceFragmentComplete(displayName, label string, reward int) {
	ann := g.Announce.Create(announce.EventBigWin, displayName, reward, map[string]string{
		"multiplier": label,
	})
	g.broadcastAnnouncement(ann)
}

// fragmentTypeLabel 碎片類型顯示名稱
func fragmentTypeLabel(ft fragment.FragmentType) string {
	switch ft {
	case fragment.FragmentBronze:
		return "銅碎片"
	case fragment.FragmentSilver:
		return "銀碎片"
	case fragment.FragmentGold:
		return "金碎片"
	default:
		return "碎片"
	}
}
