// lucky_fishing_net_handler.go — T216 幸運漁網魚
// 設計：Fishing Net 機制（BGaming Fishing Club 2，2026-04 最新）
//       撒網捕獲全場所有目標，每個目標獎勵 ×60.0
//       全部捕獲 → 全服 ×38.5 加成 77 秒（超越 T215 的 ×38.0）
//       觸發率：0.0018%（最稀有）；個人冷卻 285 秒；全服冷卻 345 秒
//       業界依據：BGaming「Fishing Club 2」Fishing Net Bonus（×60 stake，2026-04）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyFishingNetManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
}

func newLuckyFishingNetManager() *luckyFishingNetManager {
	return &luckyFishingNetManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyFishingNetFish(defID string) bool {
	return defID == "T216"
}

func (m *luckyFishingNetManager) tryLuckyFishingNetFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(345 * time.Second)
	m.personalCD[p.ID] = now.Add(285 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_fishing_net",
		Payload: map[string]interface{}{
			"event":        "net_cast",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"net_mult":     60.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🎣🌊 漁網撒出！%s 觸發漁網魚！全場所有目標被捕獲，每個獎勵 ×60.0！", p.GetDisplayName()), "critical", "#1E90FF")
	log.Printf("[LuckyFishingNet] %s 觸發漁網魚（全場捕獲 ×60.0）", p.GetDisplayName())

	go func() {
		time.Sleep(1 * time.Second)

		// 撒網捕獲全場所有目標
		caught := g.applyFishingNetCatch(60.0)

		g.broadcast(protocol.Envelope{
			Type: "lucky_fishing_net",
			Payload: map[string]interface{}{
				"event":        "net_haul",
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
				"caught_count": caught,
				"net_mult":     60.0,
			},
		})

		if caught >= 5 {
			// 完美漁網：捕獲 5 個以上 → 全服加成
			globalBoostMult := 38.5
			globalBoostSecs := 77
			g.broadcast(protocol.Envelope{
				Type: "lucky_fishing_net",
				Payload: map[string]interface{}{
					"event":             "net_perfect",
					"trigger_id":        p.ID,
					"trigger_name":      p.GetDisplayName(),
					"caught_count":      caught,
					"global_boost_mult": globalBoostMult,
					"global_boost_secs": globalBoostSecs,
				},
			})
			g.sendAnnounce(fmt.Sprintf("🎣✨ 完美漁網！捕獲 %d 個目標！全服 ×%.1f 加成 %d 秒！", caught, globalBoostMult, globalBoostSecs), "critical", "#1E90FF")
			log.Printf("[LuckyFishingNet] 完美漁網！捕獲 %d 個目標，全服 ×%.1f 加成 %d 秒（超越 T215 的 ×38.0）", caught, globalBoostMult, globalBoostSecs)
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_fishing_net",
				Payload: map[string]interface{}{
					"event":        "net_end",
					"caught_count": caught,
				},
			})
		}
	}()
	return true
}

// applyFishingNetCatch 執行漁網捕獲（捕獲場上所有目標，每個獎勵 netMult）
func (g *Game) applyFishingNetCatch(netMult float64) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	caught := 0
	for id, t := range g.targets {
		if t.HP <= 0 {
			continue
		}
		reward := int(float64(t.Def.Multiplier) * netMult)
		t.HP = 0
		delete(g.targets, id)
		caught++

		g.broadcast(protocol.Envelope{
			Type: "lucky_fishing_net",
			Payload: map[string]interface{}{
				"event":     "net_catch",
				"target_id": id,
				"reward":    reward,
				"net_mult":  netMult,
			},
		})
		g.broadcast(protocol.Envelope{
			Type: protocol.MsgTargetKill,
			Payload: map[string]interface{}{
				"id":     id,
				"reward": reward,
			},
		})
	}
	return caught
}
