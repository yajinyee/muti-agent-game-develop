// lucky_tnt_bonus_handler.go — T217 幸運 TNT 爆炸魚
// 設計：TNT Bonus 機制（BGaming Fishing Club 2，2026-04 最新）
//       水下大爆炸：全場 HP -80%，每個目標獎勵 ×100.0
//       全場清空 → 全服 ×39.0 加成 78 秒（超越 T216 的 ×38.5）
//       觸發率：0.0016%（最稀有）；個人冷卻 290 秒；全服冷卻 350 秒
//       業界依據：BGaming「Fishing Club 2」TNT Bonus（×100 stake，2026-04）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyTNTBonusManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
}

func newLuckyTNTBonusManager() *luckyTNTBonusManager {
	return &luckyTNTBonusManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyTNTBonusFish(defID string) bool {
	return defID == "T217"
}

func (m *luckyTNTBonusManager) tryLuckyTNTBonusFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(350 * time.Second)
	m.personalCD[p.ID] = now.Add(290 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_tnt_bonus",
		Payload: map[string]interface{}{
			"event":        "tnt_countdown",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"tnt_mult":     100.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("💣🌊 TNT 引爆！%s 觸發 TNT 爆炸魚！水下大爆炸，全場 HP -80%%，每個獎勵 ×100.0！", p.GetDisplayName()), "critical", "#FF4500")
	log.Printf("[LuckyTNTBonus] %s 觸發 TNT 爆炸魚（水下大爆炸 ×100.0）", p.GetDisplayName())

	go func() {
		// 3 秒倒數
		for i := 3; i >= 1; i-- {
			time.Sleep(1 * time.Second)
			g.broadcast(protocol.Envelope{
				Type: "lucky_tnt_bonus",
				Payload: map[string]interface{}{
					"event":     "tnt_tick",
					"countdown": i,
				},
			})
		}

		// 爆炸！全場 HP -80%
		blasted := g.applyTNTExplosion(100.0)

		g.broadcast(protocol.Envelope{
			Type: "lucky_tnt_bonus",
			Payload: map[string]interface{}{
				"event":         "tnt_explode",
				"trigger_id":    p.ID,
				"trigger_name":  p.GetDisplayName(),
				"blasted_count": blasted,
				"tnt_mult":      100.0,
			},
		})

		if blasted >= 3 {
			globalBoostMult := 39.0
			globalBoostSecs := 78
			g.broadcast(protocol.Envelope{
				Type: "lucky_tnt_bonus",
				Payload: map[string]interface{}{
					"event":             "tnt_perfect",
					"trigger_id":        p.ID,
					"trigger_name":      p.GetDisplayName(),
					"blasted_count":     blasted,
					"global_boost_mult": globalBoostMult,
					"global_boost_secs": globalBoostSecs,
				},
			})
			g.sendAnnounce(fmt.Sprintf("💣✨ TNT 完美爆炸！炸毀 %d 個目標！全服 ×%.1f 加成 %d 秒！", blasted, globalBoostMult, globalBoostSecs), "critical", "#FF4500")
			log.Printf("[LuckyTNTBonus] TNT 完美爆炸！炸毀 %d 個目標，全服 ×%.1f 加成 %d 秒（超越 T216 的 ×38.5）", blasted, globalBoostMult, globalBoostSecs)
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_tnt_bonus",
				Payload: map[string]interface{}{
					"event":         "tnt_end",
					"blasted_count": blasted,
				},
			})
		}
	}()
	return true
}

// applyTNTExplosion 執行 TNT 爆炸（全場 HP -80%，每個目標獎勵 tntMult）
func (g *Game) applyTNTExplosion(tntMult float64) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	blasted := 0
	for id, t := range g.targets {
		if t.HP <= 0 {
			continue
		}
		// HP -80%
		damage := int(float64(t.HP) * 0.80)
		t.HP -= damage
		if t.HP <= 0 {
			t.HP = 0
		}

		reward := int(float64(t.Def.Multiplier) * tntMult)
		blasted++

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
			Type: "lucky_tnt_bonus",
			Payload: map[string]interface{}{
				"event":     "tnt_blast",
				"target_id": id,
				"damage":    damage,
				"reward":    reward,
				"tnt_mult":  tntMult,
			},
		})
	}
	return blasted
}
