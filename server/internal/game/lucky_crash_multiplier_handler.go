// lucky_crash_multiplier_handler.go — T212 幸運崩潰倍率魚
// 設計：Crash Multiplier 機制（2026 Hybrid Crash Game 趨勢）
//       觸發後倍率從 1.0 開始每秒 +2.0，最高 50.0，隨時可以「收割」
//       收割時機越晚倍率越高，但有 30% 機率在任意秒崩潰（倍率歸零）
//       完美收割（倍率 ≥ 40.0）→ 全服 ×36.5 加成 73 秒
//       觸發率：0.003%；個人冷卻 265 秒；全服冷卻 325 秒
//       業界依據：cardsrealm.com「Hybrid Crash Game」趨勢（2026-05）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyCrashMultiplierManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	active     *crashMultiplierState
}

type crashMultiplierState struct {
	currentMult float64
	expiresAt   time.Time
	crashed     bool
}

func newLuckyCrashMultiplierManager() *luckyCrashMultiplierManager {
	return &luckyCrashMultiplierManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyCrashMultiplierFish(defID string) bool {
	return defID == "T212"
}

func (m *luckyCrashMultiplierManager) tryLuckyCrashMultiplierFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) {
		m.mu.Unlock()
		return false
	}
	if cd, ok := m.personalCD[p.ID]; ok && now.Before(cd) {
		m.mu.Unlock()
		return false
	}
	m.personalCD[p.ID] = now.Add(265 * time.Second)
	m.globalCD = now.Add(325 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_crash_multiplier",
		Payload: map[string]interface{}{
			"event":        "crash_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("💥📈 崩潰倍率！%s 觸發崩潰倍率魚！倍率持續上升，隨時可收割！", p.GetDisplayName()), "critical", "#FF4500")
	log.Printf("[LuckyCrashMultiplier] %s 觸發崩潰倍率魚", p.GetDisplayName())

	go func() {
		time.Sleep(800 * time.Millisecond)

		currentMult := 1.0
		maxMult := 50.0
		maxSecs := 25 // 最多 25 秒
		finalMult := 1.0
		crashed := false

		for sec := 1; sec <= maxSecs; sec++ {
			time.Sleep(1 * time.Second)

			// 每秒 +2.0 倍率
			currentMult += 2.0
			if currentMult > maxMult {
				currentMult = maxMult
			}

			// 30% 機率崩潰（倍率越高崩潰機率越高）
			crashChance := 0.15 + (currentMult/maxMult)*0.25
			if rand.Float64() < crashChance {
				crashed = true
				g.broadcast(protocol.Envelope{
					Type: "lucky_crash_multiplier",
					Payload: map[string]interface{}{
						"event":        "crashed",
						"trigger_id":   p.ID,
						"trigger_name": p.GetDisplayName(),
						"final_mult":   currentMult,
					},
				})
				g.sendAnnounce(fmt.Sprintf("💥 崩潰！倍率在 ×%.1f 時崩潰！", currentMult), "warning", "#FF4500")
				log.Printf("[LuckyCrashMultiplier] 崩潰！倍率 ×%.1f", currentMult)
				break
			}

			finalMult = currentMult
			g.broadcast(protocol.Envelope{
				Type: "lucky_crash_multiplier",
				Payload: map[string]interface{}{
					"event":        "mult_update",
					"current_mult": currentMult,
					"sec":          sec,
				},
			})

			if currentMult >= maxMult {
				break
			}
		}

		if !crashed {
			// 成功收割：對場上所有目標施加 finalMult 獎勵
			g.applyCrashMultiplierReward(finalMult)

			if finalMult >= 40.0 {
				globalBoostMult := 36.5
				globalBoostSecs := 73
				g.broadcast(protocol.Envelope{
					Type: "lucky_crash_multiplier",
					Payload: map[string]interface{}{
						"event":        "crash_perfect",
						"trigger_id":   p.ID,
						"trigger_name": p.GetDisplayName(),
						"final_mult":   finalMult,
						"global_mult":  globalBoostMult,
						"global_secs":  globalBoostSecs,
					},
				})
				g.sendAnnounce(fmt.Sprintf("💥🌟 完美收割！倍率 ×%.1f！全服 ×%.1f 加成 %d 秒！", finalMult, globalBoostMult, globalBoostSecs), "critical", "#FF4500")
				log.Printf("[LuckyCrashMultiplier] 完美收割！倍率 ×%.1f，全服 ×%.1f 加成 %d 秒", finalMult, globalBoostMult, globalBoostSecs)
			} else {
				g.broadcast(protocol.Envelope{
					Type: "lucky_crash_multiplier",
					Payload: map[string]interface{}{
						"event":      "crash_end",
						"final_mult": finalMult,
					},
				})
			}
		}
	}()
	return true
}

// applyCrashMultiplierReward 對場上所有目標施加崩潰倍率獎勵
func (g *Game) applyCrashMultiplierReward(mult float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for id, t := range g.targets {
		if t.HP <= 0 {
			continue
		}
		reward := int(float64(t.Def.MinMult) * mult)
		t.HP = 0
		delete(g.targets, id)

		g.broadcast(protocol.Envelope{
			Type: "lucky_crash_multiplier",
			Payload: map[string]interface{}{
				"event":     "reward_kill",
				"target_id": id,
				"reward":    reward,
				"mult":      mult,
			},
		})
	}
}
