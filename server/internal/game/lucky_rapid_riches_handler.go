// lucky_rapid_riches_handler.go — T220 幸運快速暴富魚
// 設計：Rapid Riches 機制（Reflex Gaming「Big Game Fishing Rapid Riches」，2026-05 最新）
//       5 秒內快速連擊：每次擊破 +1 連擊，每次連擊獎勵 ×200.0
//       連擊 ≥ 10 → 全服 ×41.0 加成 82 秒（新史上最高，超越 T219 的 ×40.0）
//       觸發率：0.001%（最稀有）；個人冷卻 310 秒；全服冷卻 370 秒
//       業界依據：Reflex Gaming「Big Game Fishing Rapid Riches」快速獎勵機制（2026-05）
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckyRapidRichesManager struct {
	globalCD   time.Time
	mu         sync.Mutex
	personalCD map[string]time.Time
	// 追蹤活躍的 Rapid Riches 狀態
	activeSession map[string]*rapidRichesSession
}

type rapidRichesSession struct {
	PlayerID  string
	StartTime time.Time
	HitCount  int
	mu        sync.Mutex
}

func newLuckyRapidRichesManager() *luckyRapidRichesManager {
	return &luckyRapidRichesManager{
		personalCD:    make(map[string]time.Time),
		activeSession: make(map[string]*rapidRichesSession),
	}
}

func isLuckyRapidRichesFish(defID string) bool {
	return defID == "T220"
}

func (m *luckyRapidRichesManager) tryLuckyRapidRichesFish(g *Game, p *Player) bool {
	m.mu.Lock()
	now := time.Now()
	if now.Before(m.globalCD) || now.Before(m.personalCD[p.ID]) {
		m.mu.Unlock()
		return false
	}
	m.globalCD = now.Add(370 * time.Second)
	m.personalCD[p.ID] = now.Add(310 * time.Second)

	// 建立 Rapid Riches 會話
	session := &rapidRichesSession{
		PlayerID:  p.ID,
		StartTime: now,
		HitCount:  0,
	}
	m.activeSession[p.ID] = session
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_rapid_riches",
		Payload: map[string]interface{}{
			"event":        "rapid_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     5,
			"hit_mult":     200.0,
		},
	})
	g.sendAnnounce(fmt.Sprintf("💰⚡ 快速暴富！%s 觸發快速暴富魚！5 秒內快速連擊，每次 ×200.0！", p.GetDisplayName()), "critical", "#FFD700")
	log.Printf("[LuckyRapidRiches] %s 觸發快速暴富魚（5 秒快速連擊 ×200.0）", p.GetDisplayName())

	go func() {
		// 5 秒快速連擊時間
		time.Sleep(5 * time.Second)

		m.mu.Lock()
		sess, ok := m.activeSession[p.ID]
		if ok {
			delete(m.activeSession, p.ID)
		}
		m.mu.Unlock()

		if !ok {
			return
		}

		sess.mu.Lock()
		hitCount := sess.HitCount
		sess.mu.Unlock()

		if hitCount >= 10 {
			// 完美快速暴富：10 次以上連擊 → 全服加成
			globalBoostMult := 41.0
			globalBoostSecs := 82
			g.broadcast(protocol.Envelope{
				Type: "lucky_rapid_riches",
				Payload: map[string]interface{}{
					"event":             "rapid_perfect",
					"trigger_id":        p.ID,
					"trigger_name":      p.GetDisplayName(),
					"hit_count":         hitCount,
					"global_boost_mult": globalBoostMult,
					"global_boost_secs": globalBoostSecs,
				},
			})
			g.sendAnnounce(fmt.Sprintf("💰🌟 完美快速暴富！%d 次連擊！全服 ×%.1f 加成 %d 秒！（新史上最高）", hitCount, globalBoostMult, globalBoostSecs), "critical", "#FFD700")
			log.Printf("[LuckyRapidRiches] 完美快速暴富！%d 次連擊，全服 ×%.1f 加成 %d 秒（新史上最高，超越 T219 的 ×40.0）", hitCount, globalBoostMult, globalBoostSecs)
		} else {
			g.broadcast(protocol.Envelope{
				Type: "lucky_rapid_riches",
				Payload: map[string]interface{}{
					"event":     "rapid_end",
					"hit_count": hitCount,
				},
			})
			g.sendAnnounce(fmt.Sprintf("💰 快速暴富結束！%d 次連擊（需要 10 次觸發完美）", hitCount), "normal", "#FFD700")
		}
	}()
	return true
}

// onRapidRichesHit 在 Rapid Riches 活躍期間，每次擊破都觸發此函數
func (m *luckyRapidRichesManager) onRapidRichesHit(g *Game, playerID string, targetID string, targetMult int) {
	m.mu.Lock()
	sess, ok := m.activeSession[playerID]
	m.mu.Unlock()

	if !ok {
		return
	}

	// 檢查是否在 5 秒內
	if time.Since(sess.StartTime) > 5*time.Second {
		return
	}

	sess.mu.Lock()
	sess.HitCount++
	hitCount := sess.HitCount
	sess.mu.Unlock()

	reward := targetMult * 200

	g.broadcast(protocol.Envelope{
		Type: "lucky_rapid_riches",
		Payload: map[string]interface{}{
			"event":     "rapid_hit",
			"player_id": playerID,
			"target_id": targetID,
			"hit_count": hitCount,
			"reward":    reward,
			"hit_mult":  200.0,
		},
	})
	log.Printf("[LuckyRapidRiches] 快速連擊 #%d！目標 %s，獎勵 %d（×200.0）", hitCount, targetID, reward)
}
