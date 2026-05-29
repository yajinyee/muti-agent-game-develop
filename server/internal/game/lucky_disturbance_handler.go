// lucky_disturbance_handler.go — T218 幸運擾動魚
// 設計：Disturbance System（Fisch Roblox，2026-01 最新）
//       活躍度越高倍率越高：統計觸發前 30 秒的擊破數，每次擊破 +1 擾動值
//       擾動值 ≥ 30 → 最高倍率 ×50.0，全服 ×39.5 加成 79 秒（超越 T217 的 ×39.0）
//       觸發率：0.0014%（最稀有）；個人冷卻 295 秒；全服冷卻 355 秒
//       業界依據：Fisch「Disturbance System」活躍度驅動稀有魚生成（2026-01）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyDisturbanceManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
}

func newLuckyDisturbanceManager() *luckyDisturbanceManager {
	return &luckyDisturbanceManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyDisturbanceFish(defID string) bool {
	return defID == "T218"
}

func (m *luckyDisturbanceManager) tryLuckyDisturbanceFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(355 * time.Second)
	m.personalCD[p.ID] = now.Add(295 * time.Second)
	m.mu.Unlock()

	// 計算擾動值（基於玩家最近擊破數，最高 30）
	disturbance := p.RecentKills
	if disturbance > 30 {
		disturbance = 30
	}
	if disturbance < 1 {
		disturbance = 1
	}

	// 擾動值越高，倍率越高（線性插值：1→×5.0，30→×50.0）
	disturbMult := 5.0 + float64(disturbance-1)*(50.0-5.0)/29.0

	g.broadcast(protocol.Envelope{
		Type: "lucky_disturbance",
		Payload: map[string]interface{}{
			"event":        "disturbance_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"disturbance":  disturbance,
			"disturb_mult": disturbMult,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🌊⚡ 擾動爆發！%s 觸發擾動魚！擾動值 %d，倍率 ×%.1f！", p.GetDisplayName(), disturbance, disturbMult), "critical", "#00CED1")
	log.Printf("[LuckyDisturbance] %s 觸發擾動魚（擾動值 %d，倍率 ×%.1f）", p.GetDisplayName(), disturbance, disturbMult)

	go func() {
		time.Sleep(500 * time.Millisecond)

		// 對場上所有目標施加擾動傷害
		disturbed := g.applyDisturbanceEffect(disturbMult)

		if disturbance >= 20 {
			// 高擾動值 → 全服加成
			globalBoostMult := 39.5
			globalBoostSecs := 79
			g.broadcast(protocol.Envelope{
				Type: "lucky_disturbance",
				Payload: map[string]interface{}{
					"event":             "disturbance_perfect",
					"trigger_id":        p.ID,
					"trigger_name":      p.GetDisplayName(),
					"disturbance":       disturbance,
					"disturbed_count":   disturbed,
					"disturb_mult":      disturbMult,
					"global_boost_mult": globalBoostMult,
					"global_boost_secs": globalBoostSecs,
				},
			})
			g.sendAnnounce(fmt.Sprintf("🌊✨ 完美擾動！擾動值 %d，影響 %d 個目標！全服 ×%.1f 加成 %d 秒！", disturbance, disturbed, globalBoostMult, globalBoostSecs), "critical", "#00CED1")
			log.Printf("[LuckyDisturbance] 完美擾動！擾動值 %d，影響 %d 個目標，全服 ×%.1f 加成 %d 秒（超越 T217 的 ×39.0）", disturbance, disturbed, globalBoostMult, globalBoostSecs)
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_disturbance",
				Payload: map[string]interface{}{
					"event":           "disturbance_end",
					"disturbance":     disturbance,
					"disturbed_count": disturbed,
					"disturb_mult":    disturbMult,
				},
			})
		}
	}()
	return true
}

// applyDisturbanceEffect 執行擾動效果（對場上所有目標施加倍率傷害）
func (g *Game) applyDisturbanceEffect(disturbMult float64) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	disturbed := 0
	for id, t := range g.targets {
		if t.HP <= 0 {
			continue
		}
		reward := int(float64(t.Def.Multiplier) * disturbMult)
		// HP -60%（擾動傷害）
		damage := int(float64(t.HP) * 0.60)
		t.HP -= damage
		if t.HP <= 0 {
			t.HP = 0
		}
		disturbed++

		if t.HP <= 0 {
			delete(g.targets, id)
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgTargetKill,
				Payload: map[string]interface{}{
					"id":     id,
					"reward": reward,
				},
			})
		} else {
			g.broadcast(protocol.Envelope{
				Type: protocol.MsgTargetUpdate,
				Payload: map[string]interface{}{
					"id":     id,
					"hp":     t.HP,
					"max_hp": t.Def.HP,
				},
			})
		}

		g.broadcast(protocol.Envelope{
			Type: "lucky_disturbance",
			Payload: map[string]interface{}{
				"event":        "disturbance_hit",
				"target_id":    id,
				"damage":       damage,
				"reward":       reward,
				"disturb_mult": disturbMult,
			},
		})
	}
	return disturbed
}
