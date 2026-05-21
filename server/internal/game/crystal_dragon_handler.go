// crystal_dragon_handler.go — 水晶龍收集大獎系統 handler（DAY-153）
// 業界依據：jiligames.com JILI Flying Dragon 2026「collect crystals to get the grand prize!
// Kill the Underworld Dragon and win the prize!」
// 設計：T117 水晶龍擊破後掉落水晶碎片，全服玩家共同收集水晶（目標 50 個），
// 達到目標後觸發「地獄龍大獎」，按貢獻比例分配獎勵（最高 200x betLevel）
// 社交設計：全服合作收集，讓所有玩家都有參與感，增加留存
package game

import (
	"fmt"
	"log"
	"sort"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/crystaldragon"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	CrystalDragonDefID = "T117"
)

// isCrystalDragon 判斷是否為水晶龍目標
func isCrystalDragon(defID string) bool {
	return defID == CrystalDragonDefID
}

// notifyCrystalDragonKill 水晶龍被擊破後處理（由 handleKill 呼叫）
func (g *Game) notifyCrystalDragonKill(p *player.Player, t *target.Target) {
	if g.CrystalDragon == nil {
		return
	}
	if g.CrystalDragon.IsOnCooldown() {
		log.Printf("[CrystalDragon] on cooldown, skip crystal drop for player=%s", p.ID)
		return
	}

	// 增加水晶
	newTotal, triggered := g.CrystalDragon.AddCrystals(
		p.ID, p.DisplayName, p.BetLevel, crystaldragon.CrystalPerKill,
	)

	snap := g.CrystalDragon.GetSnapshot()

	// 廣播水晶掉落
	msg := fmt.Sprintf("💎 %s 擊破水晶龍！掉落 %d 個水晶！全服進度：%d/%d",
		p.DisplayName, crystaldragon.CrystalPerKill, newTotal, crystaldragon.CrystalGoal)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCrystalDragonDrop,
		Payload: ws.CrystalDragonDropPayload{
			KillerID:      p.ID,
			KillerName:    p.DisplayName,
			CrystalsGain:  crystaldragon.CrystalPerKill,
			TotalCrystals: newTotal,
			Goal:          crystaldragon.CrystalGoal,
			Progress:      snap.Progress,
			Message:       msg,
		},
	})

	log.Printf("[CrystalDragon] player=%s killed crystal dragon, crystals=%d/%d triggered=%v",
		p.ID, newTotal, crystaldragon.CrystalGoal, triggered)

	// 達到目標，觸發地獄龍大獎
	if triggered {
		go g.triggerHellDragonReward()
	}
}

// triggerHellDragonReward 觸發地獄龍大獎（全服廣播 + 發放獎勵）
func (g *Game) triggerHellDragonReward() {
	if g.CrystalDragon == nil {
		return
	}

	contributors := g.CrystalDragon.TriggerHellDragon()
	if contributors == nil || len(contributors) == 0 {
		return
	}

	log.Printf("[CrystalDragon] Hell Dragon triggered! %d contributors", len(contributors))

	// 計算總水晶數（用於比例計算）
	totalCrystals := 0
	for _, c := range contributors {
		totalCrystals += c.Crystals
	}

	// 按貢獻排序（多的在前）
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Crystals > contributors[j].Crystals
	})

	// 發放獎勵並建立廣播列表
	entries := make([]ws.CrystalContributorEntry, 0, len(contributors))
	totalReward := 0

	g.mu.RLock()
	players := g.Players
	g.mu.RUnlock()

	for _, c := range contributors {
		reward := crystaldragon.CalcReward(c, totalCrystals)
		totalReward += reward

		entries = append(entries, ws.CrystalContributorEntry{
			PlayerID:   c.PlayerID,
			PlayerName: c.PlayerName,
			Crystals:   c.Crystals,
			Reward:     reward,
		})

		// 發放獎勵給在線玩家
		if p, ok := players[c.PlayerID]; ok {
			p.AddCoins(reward)
			log.Printf("[CrystalDragon] reward player=%s crystals=%d reward=%d",
				c.PlayerID, c.Crystals, reward)
		}
	}

	// 廣播地獄龍大獎
	topName := ""
	if len(contributors) > 0 {
		topName = contributors[0].PlayerName
	}
	msg := fmt.Sprintf("🐉 地獄龍降臨！%s 貢獻最多水晶！全服玩家共獲得 %d 金幣！",
		topName, totalReward)

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCrystalDragonReward,
		Payload: ws.CrystalDragonRewardPayload{
			Contributors: entries,
			TotalReward:  totalReward,
			Message:      msg,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventCrystalDragon, topName, totalReward, map[string]string{
		"total": fmt.Sprintf("%d", totalReward),
	})
	g.broadcastAnnouncement(ann)

	log.Printf("[CrystalDragon] Hell Dragon reward distributed: total=%d to %d players",
		totalReward, len(contributors))
}

// tickCrystalDragonDecay 水晶衰減檢查（由 game loop 每 30 秒呼叫）
func (g *Game) tickCrystalDragonDecay() {
	if g.CrystalDragon == nil {
		return
	}
	if decayed := g.CrystalDragon.CheckDecay(); decayed {
		snap := g.CrystalDragon.GetSnapshot()
		// 廣播進度更新（衰減後）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgCrystalDragonUpdate,
			Payload: ws.CrystalDragonUpdatePayload{
				TotalCrystals: snap.TotalCrystals,
				Goal:          snap.Goal,
				Progress:      snap.Progress,
			},
		})
	}
}

// sendCrystalDragonStatus 發送水晶狀態給新加入的玩家（登入時呼叫）
func (g *Game) sendCrystalDragonStatus(p *player.Player) {
	if g.CrystalDragon == nil {
		return
	}
	snap := g.CrystalDragon.GetSnapshot()
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgCrystalDragonStatus,
		Payload: ws.CrystalDragonStatusPayload{
			TotalCrystals: snap.TotalCrystals,
			Goal:          snap.Goal,
			Progress:      snap.Progress,
			CooldownSecs:  snap.CooldownSecs,
		},
	})
}
