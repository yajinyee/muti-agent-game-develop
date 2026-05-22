// freeze_bomb_handler.go — 冰凍炸彈魚系統（DAY-170）
// 業界依據：King of Ocean 2026「The freezing blast pauses an entire school for a few seconds —
// useful when a high-tier creature is escaping the frame.」
// 設計：擊破 T128 冰凍炸彈魚後，場上所有特殊目標（T101-T127）被冰凍 6 秒，停止移動
// 設計差異：與黃金海龜（全場時間停止 8 秒）不同，冰凍炸彈魚只凍結特殊目標
// 讓玩家集中火力打高價值目標，視覺上特殊目標變成冰藍色，普通目標繼續移動
package game

import (
	"log"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// freezeBombManager 冰凍炸彈魚管理器
type freezeBombManager struct {
	mu          sync.Mutex
	isActive    bool      // 是否正在冰凍中
	expiresAt   time.Time // 冰凍結束時間
	cooldown    time.Time // 全服冷卻
	frozenCount int       // 被冰凍的目標數
	FreezeSec   int       // 冰凍持續時間（秒）
	CooldownSec int       // 全服冷卻（秒）
}

func newFreezeBombManager() *freezeBombManager {
	return &freezeBombManager{
		FreezeSec:   6,  // 冰凍 6 秒
		CooldownSec: 25, // 25 秒全服冷卻
	}
}

// isFreezeBomb 判斷是否為冰凍炸彈魚
func isFreezeBomb(defID string) bool {
	return defID == "T128"
}

// IsSpecialFrozen 查詢特殊目標是否被冰凍（供 Client 端查詢）
func (m *freezeBombManager) IsSpecialFrozen() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isActive && time.Now().Before(m.expiresAt)
}

// tryFreezeBomb 冰凍炸彈魚擊破後觸發冰凍
func (g *Game) tryFreezeBomb(p *player.Player, instanceID string, fx, fy float64) {
	g.FreezeBomb.mu.Lock()
	if time.Now().Before(g.FreezeBomb.cooldown) {
		g.FreezeBomb.mu.Unlock()
		return
	}
	if g.FreezeBomb.isActive {
		g.FreezeBomb.mu.Unlock()
		return
	}
	g.FreezeBomb.isActive = true
	g.FreezeBomb.expiresAt = time.Now().Add(time.Duration(g.FreezeBomb.FreezeSec) * time.Second)
	g.FreezeBomb.cooldown = time.Now().Add(time.Duration(g.FreezeBomb.CooldownSec) * time.Second)
	g.FreezeBomb.mu.Unlock()

	// 收集場上所有特殊目標
	g.mu.RLock()
	var frozenIDs []string
	var frozenEntries []ws.FreezeBombEntry
	for _, t := range g.Targets {
		if t.Def != nil && t.Def.Type == data.TargetTypeSpecial {
			frozenIDs = append(frozenIDs, t.InstanceID)
			frozenEntries = append(frozenEntries, ws.FreezeBombEntry{
				InstanceID: t.InstanceID,
				DefID:      t.DefID,
				X:          t.X,
				Y:          t.Y,
			})
		}
	}
	g.mu.RUnlock()

	frozenCount := len(frozenEntries)

	g.FreezeBomb.mu.Lock()
	g.FreezeBomb.frozenCount = frozenCount
	g.FreezeBomb.mu.Unlock()

	// 廣播冰凍開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgFreezeBomb,
		Payload: ws.FreezeBombPayload{
			Phase:        "freeze_start",
			TriggerID:    p.ID,
			TriggerName:  p.DisplayName,
			FreezeX:      fx,
			FreezeY:      fy,
			FrozenCount:  frozenCount,
			DurationSec:  g.FreezeBomb.FreezeSec,
			FrozenTargets: frozenEntries,
		},
	})

	if frozenCount > 0 {
		log.Printf("[FreezeBomb] player=%s froze %d special targets for %ds",
			p.ID, frozenCount, g.FreezeBomb.FreezeSec)
	}

	// 全服公告：≥3 個特殊目標被冰凍時廣播
	if frozenCount >= 3 {
		g.announceFreezeBomb(p.DisplayName, frozenCount)
	}

	// 等待冰凍結束
	go func() {
		time.Sleep(time.Duration(g.FreezeBomb.FreezeSec) * time.Second)

		g.FreezeBomb.mu.Lock()
		g.FreezeBomb.isActive = false
		g.FreezeBomb.mu.Unlock()

		// 廣播冰凍結束（全服）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgFreezeBomb,
			Payload: ws.FreezeBombPayload{
				Phase:       "freeze_end",
				TriggerID:   p.ID,
				TriggerName: p.DisplayName,
				FrozenCount: frozenCount,
				DurationSec: g.FreezeBomb.FreezeSec,
			},
		})
	}()
}

// announceFreezeBomb 全服公告冰凍炸彈魚
func (g *Game) announceFreezeBomb(playerName string, frozenCount int) {
	ann := g.Announce.Create(announce.EventFreezeBomb, playerName, frozenCount, nil)
	g.broadcastAnnouncement(ann)
}
