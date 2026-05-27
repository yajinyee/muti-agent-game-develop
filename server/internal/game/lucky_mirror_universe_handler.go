// lucky_mirror_universe_handler.go — T186 幸運鏡像宇宙魚
// 業界依據：Royal Fishing「Mirror Fish」+ 量子糾纏概念
// 設計：擊破後開啟鏡像宇宙 25 秒，複製場上最強 3 個目標（HP 50%，獎勵 ×2.0）
//       25 秒內全部擊破 → 鏡像完美：全服 ×17.0 加成 36 秒
//       個人冷卻 105 秒；全服冷卻 165 秒
package game

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyMirrorUniverseManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	perfectBoost *mirrorUniversePerfectBoost
}

type mirrorUniversePerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyMirrorUniverseManager() *luckyMirrorUniverseManager {
	return &luckyMirrorUniverseManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyMirrorUniverseFish(defID string) bool {
	return defID == "T186"
}

func (m *luckyMirrorUniverseManager) getMirrorUniverseMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyMirrorUniverseManager) tryLuckyMirrorUniverseFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(105 * time.Second)
	m.globalCD = now.Add(165 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_mirror_universe",
		Payload: map[string]interface{}{
			"event":        "mirror_universe_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     25,
		},
	})
	g.sendAnnounce(fmt.Sprintf("🪞✨ 鏡像宇宙！%s 開啟鏡像宇宙！複製最強 3 個目標！獎勵 ×2.0！", p.GetDisplayName()), "critical", "#1A237E")
	log.Printf("[LuckyMirrorUniverse] %s 觸發鏡像宇宙魚", p.GetDisplayName())

	go func() {
		time.Sleep(500 * time.Millisecond)

		// 找出場上倍率最高的 3 個目標，複製（HP 50%，獎勵 ×2.0）
		g.mu.Lock()
		type targetInfo struct {
			id   string
			mult float64
		}
		var targets []targetInfo
		for id, t := range g.targets {
			if t.HP > 0 {
				targets = append(targets, targetInfo{id: id, mult: t.Def.Multiplier})
			}
		}
		sort.Slice(targets, func(i, j int) bool {
			return targets[i].mult > targets[j].mult
		})
		mirrorCount := 3
		if len(targets) < mirrorCount {
			mirrorCount = len(targets)
		}
		mirroredIDs := make([]string, 0, mirrorCount)
		for i := 0; i < mirrorCount; i++ {
			t := g.targets[targets[i].id]
			if t != nil {
				// 複製目標：HP 設為 50%
				t.HP = int(float64(t.Def.HP) * 0.5)
				mirroredIDs = append(mirroredIDs, t.InstanceID)
				g.broadcast(protocol.Envelope{
					Type: "target_update",
					Payload: map[string]interface{}{
						"id":          t.InstanceID,
						"hp":          t.HP,
						"max_hp":      t.Def.HP,
						"is_mirrored": true,
						"mirror_mult": 2.0,
					},
				})
			}
		}
		g.mu.Unlock()

		g.broadcast(protocol.Envelope{
			Type: "lucky_mirror_universe",
			Payload: map[string]interface{}{
				"event":        "mirror_universe_active",
				"mirrored_ids": mirroredIDs,
				"mirror_count": mirrorCount,
				"mirror_mult":  2.0,
				"duration":     25,
				"trigger_id":   p.ID,
				"trigger_name": p.GetDisplayName(),
			},
		})

		// 等待 25 秒後判定完美
		time.Sleep(25 * time.Second)

		// 判定：鏡像目標是否全部被擊破
		g.mu.RLock()
		allKilled := true
		for _, id := range mirroredIDs {
			if t, ok := g.targets[id]; ok && t.HP > 0 {
				allKilled = false
				break
			}
		}
		g.mu.RUnlock()

		if allKilled && mirrorCount >= 3 {
			boostMult := 17.0
			boostSecs := 36
			m.mu.Lock()
			m.perfectBoost = &mirrorUniversePerfectBoost{
				mult:      boostMult,
				expiresAt: time.Now().Add(time.Duration(boostSecs) * time.Second),
			}
			m.mu.Unlock()
			g.broadcast(protocol.Envelope{
				Type: "lucky_mirror_universe",
				Payload: map[string]interface{}{
					"event":        "mirror_perfect",
					"boost_mult":   boostMult,
					"boost_secs":   boostSecs,
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
				},
			})
			g.sendAnnounce(fmt.Sprintf("🪞🏆 鏡像完美！%s 全部擊破！全服 ×%.1f 加成 %d 秒！", p.GetDisplayName(), boostMult, boostSecs), "critical", "#0D47A1")
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_mirror_universe",
				Payload: map[string]interface{}{
					"event":        "mirror_universe_end",
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
				},
			})
		}
	}()
	return true
}
