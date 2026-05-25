// lucky_drill_torpedo_handler.go — T113 幸運鑽頭魚雷魚系統
// server-event-agent 負責維護
// 業界依據：Royal Fishing Jili「Drill Torpedo — orange mechanical lobster shoots penetrating drill
//           through multiple fish, self-explodes at end of trajectory to capture everything in blast radius」
// 設計：擊破 T113 後，發射「鑽頭魚雷」穿透場上最多 5 個目標（每個 HP -60%）；
//       魚雷飛行結束後在終點爆炸（r=180px，HP -40%）；
//       穿透目標數 × 1.2 = 累積倍率（最高 ×6.0）；
//       若穿透 ≥ 4 個目標 → 「完美穿透」：全服 ×2.2 加成 6 秒；
//       個人冷卻 18 秒；全服冷卻 30 秒
package game

import (
	"log"
	"math"
	"math/rand"
	"time"

	"chiikawa-game/internal/protocol"
)

// luckyDrillTorpedoManager 管理鑽頭魚雷系統
type luckyDrillTorpedoManager struct {
	playerCooldowns map[string]time.Time
	globalCooldown  time.Time
	// 完美穿透全服加成
	perfectBoost *drillPerfectBoost
}

type drillPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyDrillTorpedoManager() *luckyDrillTorpedoManager {
	return &luckyDrillTorpedoManager{
		playerCooldowns: make(map[string]time.Time),
	}
}

// isLuckyDrillTorpedoFish 判斷是否為鑽頭魚雷魚
func isLuckyDrillTorpedoFish(defID string) bool {
	return defID == "T113"
}

// getDrillPerfectMult 取得完美穿透全服倍率（供 handleKill 使用）
func (m *luckyDrillTorpedoManager) getDrillPerfectMult() float64 {
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyDrillTorpedoManager) canTrigger(playerID string) bool {
	now := time.Now()
	if now.Before(m.globalCooldown) {
		return false
	}
	if cd, ok := m.playerCooldowns[playerID]; ok {
		if now.Before(cd) {
			return false
		}
	}
	return true
}

// tryLuckyDrillTorpedo 嘗試觸發鑽頭魚雷
func (g *Game) tryLuckyDrillTorpedo(playerID string, killerName string) {
	m := g.luckyDrillTorpedo
	if !m.canTrigger(playerID) {
		return
	}

	now := time.Now()
	m.playerCooldowns[playerID] = now.Add(18 * time.Second)
	m.globalCooldown = now.Add(30 * time.Second)

	// 廣播觸發
	g.hub.Broadcast(protocol.MsgLuckyDrillTorpedo, protocol.LuckyDrillTorpedoPayload{
		Event:       "trigger",
		TriggerID:   playerID,
		TriggerName: killerName,
	})
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  "🚀 " + killerName + " 發射鑽頭魚雷！",
		Priority: "high",
		Color:    "#FF6B35",
	})

	go g.runDrillTorpedo(playerID, killerName)
}

// runDrillTorpedo 執行鑽頭魚雷邏輯
func (g *Game) runDrillTorpedo(playerID string, killerName string) {
	time.Sleep(400 * time.Millisecond)

	g.mu.Lock()
	p, ok := g.players[playerID]
	if !ok {
		g.mu.Unlock()
		return
	}
	betCost := p.GetBetDef().BetCost

	// 選取穿透路徑上的目標（最多 5 個，按 X 座標排序模擬直線穿透）
	type tInfo struct {
		id string
		x  float64
	}
	var candidates []tInfo
	for id, t := range g.targets {
		if t.Def.Type == "boss" {
			continue
		}
		candidates = append(candidates, tInfo{id: id, x: t.X})
	}
	g.mu.Unlock()

	// 隨機選起始方向（從左到右或從右到左）
	if rand.Intn(2) == 0 {
		// 從左到右：選 X 最小的 5 個
		for i := 0; i < len(candidates)-1; i++ {
			for j := i + 1; j < len(candidates); j++ {
				if candidates[i].x > candidates[j].x {
					candidates[i], candidates[j] = candidates[j], candidates[i]
				}
			}
		}
	} else {
		// 從右到左：選 X 最大的 5 個
		for i := 0; i < len(candidates)-1; i++ {
			for j := i + 1; j < len(candidates); j++ {
				if candidates[i].x < candidates[j].x {
					candidates[i], candidates[j] = candidates[j], candidates[i]
				}
			}
		}
	}

	maxPenetrate := 5
	if len(candidates) < maxPenetrate {
		maxPenetrate = len(candidates)
	}
	selected := candidates[:maxPenetrate]

	// 逐一穿透（每 350ms 一個）
	hitTargets := []string{}
	accumMult := 1.0
	totalReward := 0

	for i, c := range selected {
		time.Sleep(350 * time.Millisecond)

		g.mu.Lock()
		t, exists := g.targets[c.id]
		if !exists {
			g.mu.Unlock()
			continue
		}

		// HP -60%（穿透傷害）
		damage := int(float64(t.MaxHP) * 0.6)
		t.HP -= damage
		if t.HP < 1 {
			t.HP = 1
		}

		g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
			InstanceID: t.InstanceID,
			HP:         t.HP,
			MaxHP:      t.MaxHP,
			X:          t.X,
			Y:          t.Y,
		})

		hitTargets = append(hitTargets, c.id)
		accumMult += 1.2
		if accumMult > 6.0 {
			accumMult = 6.0
		}
		penReward := int(float64(betCost) * 1.2)
		totalReward += penReward

		g.hub.Broadcast(protocol.MsgLuckyDrillTorpedo, protocol.LuckyDrillTorpedoPayload{
			Event:        "penetrate",
			TriggerID:    playerID,
			TriggerName:  killerName,
			HitTargets:   hitTargets,
			PenetrateCnt: i + 1,
			AccumMult:    accumMult,
			TotalReward:  totalReward,
		})
		g.mu.Unlock()
	}

	// 終點爆炸（AOE r=180px）
	time.Sleep(500 * time.Millisecond)

	g.mu.Lock()
	// 選爆炸中心（最後一個穿透目標附近，或場地中心）
	explodeX := GameWidth / 2.0
	explodeY := GameHeight / 2.0
	if len(selected) > 0 {
		lastID := selected[len(selected)-1].id
		if lt, exists := g.targets[lastID]; exists {
			explodeX = lt.X
			explodeY = lt.Y
		}
	}

	// AOE 傷害
	aoeHits := []string{}
	for id, t := range g.targets {
		dx := t.X - explodeX
		dy := t.Y - explodeY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= 180 {
			damage := int(float64(t.MaxHP) * 0.4)
			t.HP -= damage
			if t.HP < 1 {
				t.HP = 1
			}
			g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
				InstanceID: t.InstanceID,
				HP:         t.HP,
				MaxHP:      t.MaxHP,
				X:          t.X,
				Y:          t.Y,
			})
			aoeHits = append(aoeHits, id)
		}
	}

	aoeReward := int(float64(betCost) * float64(len(aoeHits)) * 0.5)
	totalReward += aoeReward

	g.hub.Broadcast(protocol.MsgLuckyDrillTorpedo, protocol.LuckyDrillTorpedoPayload{
		Event:       "explode",
		TriggerID:   playerID,
		TriggerName: killerName,
		ExplodeX:    explodeX,
		ExplodeY:    explodeY,
		HitTargets:  aoeHits,
		TotalReward: totalReward,
		AccumMult:   accumMult,
	})

	// 結算
	p2, ok2 := g.players[playerID]
	if ok2 {
		p2.AddCoins(totalReward)
		g.sendPlayerUpdate(playerID)
	}

	// 完美穿透（穿透 ≥ 4 個）
	isPerfect := len(hitTargets) >= 4
	if isPerfect {
		g.luckyDrillTorpedo.perfectBoost = &drillPerfectBoost{
			mult:      2.2,
			expiresAt: time.Now().Add(6 * time.Second),
		}
		g.hub.Broadcast(protocol.MsgLuckyDrillTorpedo, protocol.LuckyDrillTorpedoPayload{
			Event:       "perfect",
			TriggerID:   playerID,
			TriggerName: killerName,
			TotalReward: totalReward,
			AccumMult:   2.2,
		})
		g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
			Message:  "🚀💥 完美穿透！" + killerName + " 全服 ×2.2 加成 6 秒！",
			Priority: "high",
			Color:    "#FF4500",
		})
		// 6 秒後清除加成
		go func() {
			time.Sleep(6 * time.Second)
			g.mu.Lock()
			g.luckyDrillTorpedo.perfectBoost = nil
			g.mu.Unlock()
			g.hub.Broadcast(protocol.MsgLuckyDrillTorpedo, protocol.LuckyDrillTorpedoPayload{
				Event:       "perfect_end",
				TriggerID:   playerID,
				TriggerName: killerName,
			})
		}()
	}
	g.mu.Unlock()

	log.Printf("[DrillTorpedo] Player %s: penetrate=%d, aoe=%d, mult=%.1f, reward=%d, perfect=%v",
		playerID, len(hitTargets), len(aoeHits), accumMult, totalReward, isPerfect)
}
