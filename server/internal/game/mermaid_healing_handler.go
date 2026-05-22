// mermaid_healing_handler.go — 美人魚治癒系統 handler（DAY-178）
// 業界依據：Ocean King 3 Plus「The Mermaid feature — catching the Mermaid triggers a healing
// event, restoring coins to the player and granting a brief luck boost」
// 擊破 T136 後觸發「美人魚治癒」：
//   1. 為觸發玩家恢復 betLevel × 15 金幣（治癒機制）
//   2. 觸發後 20 秒內擊破獎勵 +20% 幸運加成（全服共享）
//   3. 全服廣播「美人魚降臨」
// 設計差異：唯一的「治癒型」目標，讓玩家在連續失敗後有「回血」的機會
package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// MermaidHealCooldownSec 個人冷卻時間（秒）
	MermaidHealCooldownSec = 30
	// MermaidHealMultiplier 治癒金幣倍率（betLevel × 此值）
	MermaidHealMultiplier = 15
	// MermaidLuckBoostDurationSec 幸運加成持續時間（秒）
	MermaidLuckBoostDurationSec = 20
	// MermaidLuckBoostPercent 幸運加成百分比（+20%）
	MermaidLuckBoostPercent = 0.20
	// MermaidAnnounceThreshold 全服公告門檻（治癒金幣數）
	MermaidAnnounceThreshold = 100
)

// mermaidManager 美人魚治癒管理器
type mermaidManager struct {
	mu       sync.Mutex
	cooldown map[string]time.Time // playerID → 冷卻結束時間
	// 幸運加成狀態（全服共享）
	luckActive   bool
	luckExpiresAt time.Time
}

// newMermaidManager 建立美人魚治癒管理器
func newMermaidManager() *mermaidManager {
	return &mermaidManager{
		cooldown: make(map[string]time.Time),
	}
}

// isMermaid 判斷是否為美人魚（T136）
func isMermaid(defID string) bool {
	return defID == "T136"
}

// isOnCooldown 檢查玩家是否在冷卻中
func (m *mermaidManager) isOnCooldown(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	cd, ok := m.cooldown[playerID]
	if !ok {
		return false
	}
	return time.Now().Before(cd)
}

// setCooldown 設定玩家冷卻
func (m *mermaidManager) setCooldown(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cooldown[playerID] = time.Now().Add(time.Duration(MermaidHealCooldownSec) * time.Second)
}

// activateLuckBoost 激活全服幸運加成
func (m *mermaidManager) activateLuckBoost() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.luckActive = true
	m.luckExpiresAt = time.Now().Add(time.Duration(MermaidLuckBoostDurationSec) * time.Second)
}

// getMermaidLuckBoost 取得幸運加成（供 handleKill 使用）
// 回傳額外加成比例（0.0 = 無加成，0.20 = +20%）
func (m *mermaidManager) getMermaidLuckBoost() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.luckActive {
		return 0.0
	}
	if time.Now().After(m.luckExpiresAt) {
		m.luckActive = false
		return 0.0
	}
	return MermaidLuckBoostPercent
}

// tryMermaidHealing 擊破 T136 後觸發美人魚治癒（DAY-178）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryMermaidHealing(p *player.Player) {
	// 個人冷卻檢查
	if g.Mermaid.isOnCooldown(p.ID) {
		return
	}
	g.Mermaid.setCooldown(p.ID)

	// 計算治癒金幣
	healAmount := p.BetLevel * MermaidHealMultiplier
	if healAmount < 1 {
		healAmount = 1
	}

	// 發放治癒金幣
	p.AddReward(healAmount)

	log.Printf("[Mermaid] player=%s healed=%d coins (betLevel=%d)",
		p.ID, healAmount, p.BetLevel)

	// 廣播美人魚治癒（個人）
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgMermaidHealing,
		Payload: ws.MermaidHealingPayload{
			Phase:      "heal_start",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			HealAmount: healAmount,
			NewBalance: p.Coins,
			LuckBoostDurationSec: MermaidLuckBoostDurationSec,
		},
	})

	// 全服廣播：美人魚降臨
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgMermaidHealing,
		Payload: ws.MermaidHealingPayload{
			Phase:      "heal_broadcast",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			HealAmount: healAmount,
		},
	})

	// 激活全服幸運加成
	g.Mermaid.activateLuckBoost()

	// 廣播幸運加成開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgMermaidHealing,
		Payload: ws.MermaidHealingPayload{
			Phase:                "luck_start",
			PlayerID:             p.ID,
			PlayerName:           p.DisplayName,
			LuckBoostPercent:     MermaidLuckBoostPercent,
			LuckBoostDurationSec: MermaidLuckBoostDurationSec,
		},
	})

	// 全服公告
	if healAmount >= MermaidAnnounceThreshold {
		g.announceMermaidHealing(p.DisplayName, healAmount)
	}

	// 等待幸運加成結束
	time.Sleep(time.Duration(MermaidLuckBoostDurationSec) * time.Second)

	// 廣播幸運加成結束（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgMermaidHealing,
		Payload: ws.MermaidHealingPayload{
			Phase: "luck_end",
		},
	})

	log.Printf("[Mermaid] player=%s luck boost ended", p.ID)
}

// getMermaidLuckBoost 取得美人魚幸運加成（供 handleKill 使用）
func (g *Game) getMermaidLuckBoost() float64 {
	return g.Mermaid.getMermaidLuckBoost()
}

// announceMermaidHealing 全服公告美人魚治癒（DAY-178）
func (g *Game) announceMermaidHealing(playerName string, healAmount int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "mermaid_healing",
			"message":    fmt.Sprintf("🧜 %s 遇見美人魚！恢復了 %d 金幣！全服幸運加成 +20%%！", playerName, healAmount),
			"color":      "#00CED1", // 深青色（美人魚感）
			"duration":   4.0,
			"priority":   2,
		},
	})
}
