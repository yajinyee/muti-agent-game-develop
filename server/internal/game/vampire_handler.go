// vampire_handler.go — 吸血鬼成長倍率系統 handler（DAY-152）
// 業界依據：jiligames.com 2026「The explicit multiplier of vampires increases the more you fight,
// and there is a chance that you can enter the multiplier mode, up to X5.」
// 設計：T116 吸血鬼每次被命中倍率增加，達到閾值後進入「倍率模式」（全服廣播），最終擊破獲得最高倍率獎勵
// 命中 1-4 次：倍率 ×1.0（基礎）
// 命中 5-9 次：倍率 ×2.0（覺醒）
// 命中 10-14 次：倍率 ×3.5（狂暴）
// 命中 ≥15 次：倍率 ×5.0（血月模式，全服廣播）
package game

import (
	"fmt"
	"log"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	VampireDefID = "T116"

	// 倍率成長閾值
	VampirePhase1Hits = 5  // 5 次命中進入覺醒（×2.0）
	VampirePhase2Hits = 10 // 10 次命中進入狂暴（×3.5）
	VampirePhase3Hits = 15 // 15 次命中進入血月（×5.0，全服廣播）

	// 各階段倍率加成
	VampirePhase1Mult = 2.0 // 覺醒
	VampirePhase2Mult = 3.5 // 狂暴
	VampirePhase3Mult = 5.0 // 血月
)

// isVampire 判斷是否為吸血鬼目標
func isVampire(defID string) bool {
	return defID == VampireDefID
}

// getVampireMultBonus 根據命中次數取得倍率加成
func getVampireMultBonus(hitCount int) float64 {
	switch {
	case hitCount >= VampirePhase3Hits:
		return VampirePhase3Mult
	case hitCount >= VampirePhase2Hits:
		return VampirePhase2Mult
	case hitCount >= VampirePhase1Hits:
		return VampirePhase1Mult
	default:
		return 1.0
	}
}

// getVampirePhaseName 取得階段名稱
func getVampirePhaseName(hitCount int) string {
	switch {
	case hitCount >= VampirePhase3Hits:
		return "🩸 血月模式"
	case hitCount >= VampirePhase2Hits:
		return "🦇 狂暴模式"
	case hitCount >= VampirePhase1Hits:
		return "🌙 覺醒模式"
	default:
		return "😴 沉睡"
	}
}

// notifyVampireHit 吸血鬼被命中時呼叫（由 handleAttack 呼叫，命中但未擊破時）
// 更新倍率並廣播階段變化
func (g *Game) notifyVampireHit(p *player.Player, t *target.Target) {
	if t == nil || !isVampire(t.DefID) {
		return
	}

	prevHits := t.HitCount - 1 // HitCount 已在 combat 中增加
	currHits := t.HitCount

	prevPhase := getVampirePhase(prevHits)
	currPhase := getVampirePhase(currHits)

	// 計算新倍率（基礎倍率 × 成長加成）
	newMultBonus := getVampireMultBonus(currHits)
	newMult := t.Def.MultiplierMin * newMultBonus

	// 更新目標倍率
	g.mu.Lock()
	if liveTarget, ok := g.Targets[t.InstanceID]; ok {
		liveTarget.Multiplier = newMult
	}
	g.mu.Unlock()

	// 廣播倍率更新（讓所有玩家看到吸血鬼在成長）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgVampireGrow,
		Payload: ws.VampireGrowPayload{
			InstanceID:  t.InstanceID,
			HitCount:    currHits,
			MultBonus:   newMultBonus,
			NewMult:     newMult,
			PhaseName:   getVampirePhaseName(currHits),
			PhaseChanged: currPhase != prevPhase,
		},
	})

	// 階段變化時廣播特殊訊息
	if currPhase != prevPhase {
		log.Printf("[Vampire] instance=%s phase changed: %s → %s (hits=%d, mult=%.1fx)",
			t.InstanceID, getVampirePhaseName(prevHits), getVampirePhaseName(currHits), currHits, newMultBonus)

		// 血月模式：全服廣播 + 公告
		if currPhase == 3 {
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgVampireBloodMoon,
				Payload: ws.VampireBloodMoonPayload{
					InstanceID: t.InstanceID,
					HitCount:   currHits,
					MultBonus:  VampirePhase3Mult,
					NewMult:    newMult,
					Message:    fmt.Sprintf("🩸 吸血鬼進入血月模式！倍率 ×%.0f！誰能擊破它？", VampirePhase3Mult),
				},
			})
			ann := g.Announce.Create(announce.EventVampireBloodMoon, "", int(newMult), map[string]string{
				"mult": fmt.Sprintf("%.0f", VampirePhase3Mult),
			})
			g.broadcastAnnouncement(ann)
		}
	}
}

// getVampirePhase 取得階段編號（0=沉睡, 1=覺醒, 2=狂暴, 3=血月）
func getVampirePhase(hitCount int) int {
	switch {
	case hitCount >= VampirePhase3Hits:
		return 3
	case hitCount >= VampirePhase2Hits:
		return 2
	case hitCount >= VampirePhase1Hits:
		return 1
	default:
		return 0
	}
}

// notifyVampireKill 吸血鬼被擊破時呼叫（由 handleKill 呼叫）
// 廣播最終倍率和獎勵
func (g *Game) notifyVampireKill(p *player.Player, t *target.Target, finalReward int) {
	if t == nil || !isVampire(t.DefID) {
		return
	}

	multBonus := getVampireMultBonus(t.HitCount)
	phaseName := getVampirePhaseName(t.HitCount)

	log.Printf("[Vampire] player=%s killed vampire: hits=%d, mult=%.1fx, reward=%d",
		p.ID, t.HitCount, multBonus, finalReward)

	// 廣播吸血鬼擊破（全服可見）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgVampireKilled,
		Payload: ws.VampireKilledPayload{
			InstanceID:  t.InstanceID,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			HitCount:    t.HitCount,
			MultBonus:   multBonus,
			FinalMult:   t.Multiplier,
			FinalReward: finalReward,
			PhaseName:   phaseName,
			Message:     fmt.Sprintf("🦇 %s 擊破了吸血鬼！命中 %d 次，%s，獲得 %d 金幣！", p.DisplayName, t.HitCount, phaseName, finalReward),
		},
	})

	// 血月模式擊破：全服公告
	if t.HitCount >= VampirePhase3Hits {
		ann := g.Announce.Create(announce.EventVampireKill, p.DisplayName, finalReward, map[string]string{
			"hits": fmt.Sprintf("%d", t.HitCount),
			"mult": fmt.Sprintf("%.0f", multBonus),
		})
		g.broadcastAnnouncement(ann)
	}
}
