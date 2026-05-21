// drill_lobster_handler.go — 鑽頭龍蝦連帶效果 handler（DAY-142）
// 業界依據：Royal Fishing JILI 2026「Drill Bit Lobster (80X) — fires a penetrating drill
// through multiple fish before self-detonating, capturing everything in blast radius」
// 擊破 T106 後觸發穿透鑽頭：沿水平方向穿透所有目標，到達邊緣後爆炸，連帶擊破爆炸範圍內目標
package game

import (
	"fmt"
	"log"
	"math"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// DrillPenetrateYRange 鑽頭穿透 Y 軸容差（px）— 水平穿透範圍
	DrillPenetrateYRange = 80.0
	// DrillExplosionRadius 鑽頭到達邊緣後爆炸半徑（px）
	DrillExplosionRadius = 180.0
	// DrillRewardMult 鑽頭連帶擊破獎勵倍率（比直接擊破低，平衡 RTP）
	DrillRewardMult = 0.55
	// DrillPenetrateDelayMs 穿透動畫延遲（ms）
	DrillPenetrateDelayMs = 60
)

// isDrillLobster 判斷是否為鑽頭龍蝦（T106）
func isDrillLobster(defID string) bool {
	return defID == "T106"
}

// drillTarget 鑽頭連帶目標（內部結構）
type drillTarget struct {
	instanceID string
	defID      string
	x, y       float64
	multiplier float64
}

// tryDrillLobsterChain 擊破 T106 後觸發穿透鑽頭連帶效果（DAY-142）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryDrillLobsterChain(p *player.Player, triggerID string, triggerX, triggerY float64) {
	// 收集穿透路徑上的目標（水平方向，Y 軸 ±80px）
	g.mu.RLock()
	var penetrateTargets []drillTarget

	for id, t := range g.Targets {
		if id == triggerID || t.HP <= 0 || t.DefID == "B001" {
			continue
		}
		if math.Abs(t.Y-triggerY) <= DrillPenetrateYRange {
			penetrateTargets = append(penetrateTargets, drillTarget{
				instanceID: t.InstanceID,
				defID:      t.DefID,
				x:          t.X,
				y:          t.Y,
				multiplier: t.Multiplier,
			})
		}
	}
	g.mu.RUnlock()

	// 廣播鑽頭發射（讓 Client 播放穿透動畫）
	penetrateIDs := make([]string, 0, len(penetrateTargets))
	for _, dt := range penetrateTargets {
		penetrateIDs = append(penetrateIDs, dt.instanceID)
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgDrillLobsterChain,
		Payload: ws.DrillLobsterChainPayload{
			TriggerID:    triggerID,
			TriggerX:     triggerX,
			TriggerY:     triggerY,
			Phase:        "drill_start",
			PenetrateIDs: penetrateIDs,
		},
	})

	// 分批穿透（每 60ms 一個，製造連續穿透感）
	totalReward := 0
	var killedEntries []ws.DrillKillEntry

	for i, dt := range penetrateTargets {
		time.Sleep(time.Duration(DrillPenetrateDelayMs) * time.Millisecond)

		g.mu.Lock()
		t, ok := g.Targets[dt.instanceID]
		if !ok || t.HP <= 0 {
			g.mu.Unlock()
			continue
		}
		reward := int(float64(p.BetLevel) * dt.multiplier * DrillRewardMult)
		if reward < 1 {
			reward = 1
		}
		t.HP = 0
		delete(g.Targets, dt.instanceID)
		g.mu.Unlock()

		totalReward += reward
		killedEntries = append(killedEntries, ws.DrillKillEntry{
			InstanceID: dt.instanceID,
			DefID:      dt.defID,
			Multiplier: dt.multiplier,
			Reward:     reward,
			Phase:      "penetrate",
		})

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetKill,
			Payload: ws.TargetKillPayload{
				InstanceID: dt.instanceID,
				KillerID:   p.ID,
				Reward:     reward,
				Multiplier: dt.multiplier,
			},
		})

		log.Printf("[DrillLobster] penetrate[%d] target=%s mult=%.0f reward=%d",
			i, dt.instanceID, dt.multiplier, reward)
	}

	// 鑽頭到達邊緣，爆炸！
	explodeX := 1280.0
	if triggerX > 640 {
		explodeX = 0.0
	}
	explodeY := triggerY

	time.Sleep(200 * time.Millisecond)

	// 收集爆炸範圍內的目標
	g.mu.RLock()
	var explodeTargets []drillTarget
	for id, t := range g.Targets {
		if t.HP <= 0 || t.DefID == "B001" {
			continue
		}
		dx := t.X - explodeX
		dy := t.Y - explodeY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= DrillExplosionRadius {
			explodeTargets = append(explodeTargets, drillTarget{
				instanceID: id,
				defID:      t.DefID,
				x:          t.X,
				y:          t.Y,
				multiplier: t.Multiplier,
			})
		}
	}
	g.mu.RUnlock()

	// 廣播爆炸開始
	explodeIDs := make([]string, 0, len(explodeTargets))
	for _, dt := range explodeTargets {
		explodeIDs = append(explodeIDs, dt.instanceID)
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgDrillLobsterChain,
		Payload: ws.DrillLobsterChainPayload{
			TriggerID:  triggerID,
			TriggerX:   explodeX,
			TriggerY:   explodeY,
			Phase:      "explosion",
			ExplodeIDs: explodeIDs,
		},
	})

	time.Sleep(100 * time.Millisecond)

	// 爆炸擊破
	for _, dt := range explodeTargets {
		g.mu.Lock()
		t, ok := g.Targets[dt.instanceID]
		if !ok || t.HP <= 0 {
			g.mu.Unlock()
			continue
		}
		reward := int(float64(p.BetLevel) * dt.multiplier * DrillRewardMult)
		if reward < 1 {
			reward = 1
		}
		t.HP = 0
		delete(g.Targets, dt.instanceID)
		g.mu.Unlock()

		totalReward += reward
		killedEntries = append(killedEntries, ws.DrillKillEntry{
			InstanceID: dt.instanceID,
			DefID:      dt.defID,
			Multiplier: dt.multiplier,
			Reward:     reward,
			Phase:      "explosion",
		})

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetKill,
			Payload: ws.TargetKillPayload{
				InstanceID: dt.instanceID,
				KillerID:   p.ID,
				Reward:     reward,
				Multiplier: dt.multiplier,
			},
		})
	}

	if totalReward <= 0 {
		return
	}

	// 發放總獎勵
	p.AddReward(totalReward)

	// 廣播鑽頭連帶結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgDrillLobsterChain,
		Payload: ws.DrillLobsterChainPayload{
			TriggerID:     triggerID,
			TriggerX:      triggerX,
			TriggerY:      triggerY,
			Phase:         "result",
			KilledTargets: killedEntries,
			TotalReward:   totalReward,
			KillerID:      p.ID,
			KillerName:    p.DisplayName,
		},
	})

	// 個人結果通知
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgReward,
		Payload: ws.RewardPayload{
			Source:     "drill_lobster",
			Amount:     totalReward,
			Multiplier: float64(len(killedEntries)),
			NewBalance: p.Coins,
		},
	})

	// 全服公告：連帶擊破 ≥3 個目標
	if len(killedEntries) >= 3 {
		g.announceDrillLobsterChain(p.DisplayName, len(killedEntries), totalReward)
	}

	log.Printf("[DrillLobster] player=%s penetrate=%d explode=%d total_reward=%d",
		p.ID, len(penetrateTargets), len(explodeTargets), totalReward)
}

// announceDrillLobsterChain 全服公告鑽頭龍蝦連帶效果（DAY-142）
func (g *Game) announceDrillLobsterChain(playerName string, killCount int, reward int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "drill_lobster_chain",
			"message":    fmt.Sprintf("🦞 %s 的鑽頭龍蝦連帶擊破 %d 個目標！獲得 %d 金幣！", playerName, killCount, reward),
			"color":      "#FF6B35",
			"duration":   4.0,
			"priority":   2,
		},
	})
}
