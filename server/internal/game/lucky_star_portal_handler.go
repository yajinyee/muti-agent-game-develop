// lucky_star_portal_handler.go — T166 幸運星際門戶魚
// 業界依據：Fishing Carnival「portal teleport mechanic」
// 設計：擊破後開啟星際門戶，隨機傳送 5 個目標到高密度區域（中央），傳送後 HP -50%
//       全部傳送 → 完美門戶：全服 ×5.5 加成 12 秒
//       個人冷卻 55 秒；全服冷卻 85 秒
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyStarPortalManager struct {
	mu           sync.Mutex
	personalCD   map[string]time.Time
	globalCD     time.Time
	perfectBoost *starPortalPerfectBoost
}

type starPortalPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

func newLuckyStarPortalManager() *luckyStarPortalManager {
	return &luckyStarPortalManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckyStarPortalFish(defID string) bool {
	return defID == "T166"
}

func (m *luckyStarPortalManager) getStarPortalPerfectMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.perfectBoost != nil && time.Now().Before(m.perfectBoost.expiresAt) {
		return m.perfectBoost.mult
	}
	return 1.0
}

func (m *luckyStarPortalManager) tryLuckyStarPortalFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(55 * time.Second)
	m.globalCD = now.Add(85 * time.Second)
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_star_portal",
		Payload: map[string]interface{}{
			"event":        "portal_open",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
		},
	})
	g.sendAnnounce(fmt.Sprintf("🌌✨ 星際門戶！%s 開啟傳送門！5 個目標被傳送！", p.GetDisplayName()), "high", "#7B1FA2")
	log.Printf("[LuckyStarPortal] %s 觸發星際門戶魚", p.GetDisplayName())

	go func() {
		time.Sleep(500 * time.Millisecond)

		// 傳送 5 個目標到中央高密度區域，HP -50%
		teleportCount := g.applyStarPortalTeleport(p)

		g.broadcast(protocol.Envelope{
			Type: "lucky_star_portal",
			Payload: map[string]interface{}{
				"event":           "portal_teleport",
				"teleport_count":  teleportCount,
				"trigger_id":      p.ID,
			},
		})

		time.Sleep(1 * time.Second)

		// 判定是否完美（傳送 ≥5 個）
		if teleportCount >= 5 {
			m.mu.Lock()
			m.perfectBoost = &starPortalPerfectBoost{
				mult:      5.5,
				expiresAt: time.Now().Add(12 * time.Second),
			}
			m.mu.Unlock()

			g.broadcast(protocol.Envelope{
				Type: "lucky_star_portal",
				Payload: map[string]interface{}{
					"event":        "portal_perfect",
					"boost_mult":   5.5,
					"boost_secs":   12,
					"trigger_id":   p.ID,
					"trigger_name": p.GetDisplayName(),
				},
			})
			g.sendAnnounce(fmt.Sprintf("🌌💫 完美門戶！%s 傳送 %d 個目標！全服 ×5.5 加成 12 秒！", p.GetDisplayName(), teleportCount), "critical", "#E040FB")
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_star_portal",
				Payload: map[string]interface{}{
					"event":          "portal_end",
					"teleport_count": teleportCount,
					"trigger_id":     p.ID,
				},
			})
		}
	}()

	return true
}

// applyStarPortalTeleport — 傳送最多 5 個目標到中央，HP -50%
func (g *Game) applyStarPortalTeleport(p *Player) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 收集所有非 BOSS 目標
	var candidates []*Target
	for _, t := range g.targets {
		if t.Def.Type != "boss" {
			candidates = append(candidates, t)
		}
	}

	// 隨機選最多 5 個
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	count := 5
	if len(candidates) < count {
		count = len(candidates)
	}

	// 中央高密度區域座標（畫面中央）
	centerX := 640.0
	centerY := 360.0

	teleportCount := 0
	for i := 0; i < count; i++ {
		t := candidates[i]
		// 傳送到中央附近（隨機偏移 ±100px）
		t.X = centerX + (rand.Float64()-0.5)*200
		t.Y = centerY + (rand.Float64()-0.5)*200
		// HP -50%
		damage := t.MaxHP / 2
		t.HP -= damage
		if t.HP < 1 {
			t.HP = 1
		}
		teleportCount++

		// 廣播位置更新
		g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
			InstanceID: t.InstanceID,
			HP:         t.HP,
			MaxHP:      t.MaxHP,
			X:          t.X,
			Y:          t.Y,
		})
	}
	return teleportCount
}
