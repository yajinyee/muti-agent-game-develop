// lucky_skill_chain_handler.go — T192 幸運技能連鎖魚
// 設計：技能連鎖 25 秒，每次擊破提升技能等級（Lv.1-10）
//       Lv.1=×2.0, Lv.5=×10.0, Lv.10=×25.0
//       達到 Lv.10 → 全服 ×20.0 加成 38 秒（新最高）
//       否則結束時全服 ×19.5 加成 35 秒
//       觸發率：0.035%；個人冷卻 125 秒；全服冷卻 185 秒
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"chiikawa-game/internal/protocol"
)

type luckySkillChainManager struct {
	mu         sync.Mutex
	personalCD map[string]time.Time
	globalCD   time.Time
	chainBoost *skillChainPerfectBoost
	// 活躍的技能連鎖 session
	activeSession *skillChainSession
}

type skillChainPerfectBoost struct {
	mult      float64
	expiresAt time.Time
}

type skillChainSession struct {
	playerID  string
	level     int
	expiresAt time.Time
}

func newLuckySkillChainManager() *luckySkillChainManager {
	return &luckySkillChainManager{
		personalCD: make(map[string]time.Time),
	}
}

func isLuckySkillChainFish(defID string) bool {
	return defID == "T192"
}

func (m *luckySkillChainManager) getSkillChainMult() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.chainBoost != nil && time.Now().Before(m.chainBoost.expiresAt) {
		return m.chainBoost.mult
	}
	return 1.0
}

// 技能等級對應倍率
func skillLevelToMult(level int) float64 {
	switch {
	case level >= 10:
		return 25.0
	case level >= 9:
		return 20.0
	case level >= 8:
		return 16.0
	case level >= 7:
		return 13.0
	case level >= 6:
		return 11.0
	case level >= 5:
		return 10.0
	case level >= 4:
		return 7.0
	case level >= 3:
		return 5.0
	case level >= 2:
		return 3.0
	default:
		return 2.0
	}
}

func (m *luckySkillChainManager) tryLuckySkillChainFish(g *Game, p *Player) bool {
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
	m.personalCD[p.ID] = now.Add(125 * time.Second)
	m.globalCD = now.Add(185 * time.Second)

	// 建立 session
	session := &skillChainSession{
		playerID:  p.ID,
		level:     1,
		expiresAt: now.Add(25 * time.Second),
	}
	m.activeSession = session
	m.mu.Unlock()

	g.broadcast(protocol.Envelope{
		Type: "lucky_skill_chain",
		Payload: map[string]interface{}{
			"event":        "skill_chain_start",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"duration":     25,
			"level":        1,
			"level_mult":   skillLevelToMult(1),
		},
	})
	g.sendAnnounce(fmt.Sprintf("🔗✨ 技能連鎖！%s 啟動技能連鎖！25 秒內擊破提升等級！Lv.10 → 全服 ×20.0！", p.GetDisplayName()), "critical", "#1A237E")
	log.Printf("[LuckySkillChain] %s 觸發技能連鎖魚", p.GetDisplayName())

	go func() {
		time.Sleep(25 * time.Second)

		m.mu.Lock()
		finalLevel := 1
		if m.activeSession != nil && m.activeSession.playerID == p.ID {
			finalLevel = m.activeSession.level
			m.activeSession = nil
		}
		m.mu.Unlock()

		var globalBoostMult float64
		var globalBoostSecs int
		isPerfect := finalLevel >= 10

		if isPerfect {
			// 達到 Lv.10 → 全服 ×20.0 加成 38 秒（新最高）
			globalBoostMult = 20.0
			globalBoostSecs = 38
		} else {
			// 未達 Lv.10 → 全服 ×19.5 加成 35 秒
			globalBoostMult = 19.5
			globalBoostSecs = 35
		}

		m.mu.Lock()
		m.chainBoost = &skillChainPerfectBoost{
			mult:      globalBoostMult,
			expiresAt: time.Now().Add(time.Duration(globalBoostSecs) * time.Second),
		}
		m.mu.Unlock()

		finalMult := skillLevelToMult(finalLevel)
		g.broadcast(protocol.Envelope{
			Type: "lucky_skill_chain",
			Payload: map[string]interface{}{
				"event":              "skill_chain_complete",
				"trigger_id":         p.ID,
				"trigger_name":       p.GetDisplayName(),
				"final_level":        finalLevel,
				"final_mult":         finalMult,
				"is_perfect":         isPerfect,
				"global_boost_mult":  globalBoostMult,
				"global_boost_secs":  globalBoostSecs,
			},
		})

		if isPerfect {
			g.sendAnnounce(fmt.Sprintf("🔗🏆 完美連鎖！%s 達到 Lv.10！×%.1f！全服 ×%.1f 加成 %d 秒！（新最高）",
				p.GetDisplayName(), finalMult, globalBoostMult, globalBoostSecs), "critical", "#0D47A1")
		} else {
			g.sendAnnounce(fmt.Sprintf("🔗 技能連鎖結束！%s 達到 Lv.%d（×%.1f）！全服 ×%.1f 加成 %d 秒！",
				p.GetDisplayName(), finalLevel, finalMult, globalBoostMult, globalBoostSecs), "high", "#1565C0")
		}
	}()
	return true
}

// onKillDuringSkillChain — 擊破時提升技能等級
func (m *luckySkillChainManager) onKillDuringSkillChain(g *Game, p *Player) {
	m.mu.Lock()
	if m.activeSession == nil || m.activeSession.playerID != p.ID {
		m.mu.Unlock()
		return
	}
	if time.Now().After(m.activeSession.expiresAt) {
		m.activeSession = nil
		m.mu.Unlock()
		return
	}
	if m.activeSession.level < 10 {
		m.activeSession.level++
	}
	level := m.activeSession.level
	m.mu.Unlock()

	levelMult := skillLevelToMult(level)
	g.broadcast(protocol.Envelope{
		Type: "lucky_skill_chain",
		Payload: map[string]interface{}{
			"event":        "skill_chain_level_up",
			"trigger_id":   p.ID,
			"trigger_name": p.GetDisplayName(),
			"level":        level,
			"level_mult":   levelMult,
		},
	})
}

func (m *luckySkillChainManager) isSkillChainActive(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.activeSession != nil &&
		m.activeSession.playerID == playerID &&
		time.Now().Before(m.activeSession.expiresAt)
}
