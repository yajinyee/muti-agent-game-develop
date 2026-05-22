// blackhole_handler.go — 黑洞漩渦武器系統 handler（DAY-166）
// 業界依據：
//   - Ocean King 3 2026 Vortex 機制 — 放置後吸引周圍目標向中心移動，最終爆炸擊破
//   - Black Hole Fishing 2026（Steam）— 用黑洞吸魚的核心玩法，2026 年最新趨勢
//   - 設計差異：與炸彈（即時爆炸）不同，黑洞有 3 秒「吸引期」，讓玩家看到目標被吸入的過程
//     費用 10x betLevel（介於魚雷 6x 和軌道炮 15x 之間），吸引半徑 300px（比炸彈 200px 大 50%）
package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"digital-twin/server/internal/game/specialweapon"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// handleBlackHoleFire 黑洞漩渦武器發射（DAY-166）
// 流程：扣費 → 廣播 black_hole_place → 等 1.5s（吸引期）→ 廣播 black_hole_suck → 等 1.5s（爆炸期）→ 擊破目標 → 廣播 result
func (g *Game) handleBlackHoleFire(p *player.Player, cx, cy float64) {
	// 計算費用（10x betLevel）
	cost := p.BetLevel * specialweapon.BlackHoleCostMultiplier
	if p.Coins < cost {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "insufficient_coins", Message: fmt.Sprintf("金幣不足，黑洞費用 %d 金幣", cost)},
		})
		// 退還充能
		g.SpecialWeapon.AddCharge(p.ID, specialweapon.WeaponBlackHole)
		return
	}

	// 扣除費用
	p.Coins -= cost
	if p.Coins < 0 {
		p.Coins = 0
	}

	log.Printf("[BlackHole] player=%s placed at (%.0f, %.0f), cost=%d, balance=%d",
		p.ID, cx, cy, cost, p.Coins)

	// Phase 1：廣播黑洞放置（全服看到黑洞出現）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBlackHoleResult,
		Payload: ws.BlackHoleResultPayload{
			ShooterID:   p.ID,
			ShooterName: p.DisplayName,
			Phase:       "black_hole_place",
			CenterX:     cx,
			CenterY:     cy,
			Radius:      specialweapon.BlackHoleRadius,
			Cost:        cost,
		},
	})

	// 等待吸引期（1.5 秒）— 讓 Client 播放目標被吸入的動畫
	time.Sleep(1500 * time.Millisecond)

	// 計算吸引範圍內的目標
	g.mu.RLock()
	targets := make([]specialweapon.TargetPos, 0, len(g.Targets))
	for _, t := range g.Targets {
		if t.HP > 0 && t.DefID != "B001" { // BOSS 不受黑洞影響
			targets = append(targets, specialweapon.TargetPos{
				InstanceID: t.InstanceID,
				X:          t.X,
				Y:          t.Y,
				Multiplier: t.Multiplier,
			})
		}
	}
	g.mu.RUnlock()

	hitIDs := specialweapon.CalcBlackHoleTargets(cx, cy, targets)

	// Phase 2：廣播吸入數量（讓 Client 顯示「吸入 N 個目標」）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBlackHoleResult,
		Payload: ws.BlackHoleResultPayload{
			ShooterID:   p.ID,
			ShooterName: p.DisplayName,
			Phase:       "black_hole_suck",
			CenterX:     cx,
			CenterY:     cy,
			Radius:      specialweapon.BlackHoleRadius,
			SuckedCount: len(hitIDs),
			Cost:        cost,
		},
	})

	// 等待爆炸期（1.5 秒）— 讓 Client 播放目標聚集到中心的動畫
	time.Sleep(1500 * time.Millisecond)

	// Phase 3：爆炸，擊破目標
	hitEntries := make([]ws.BlackHoleKillEntry, 0, len(hitIDs))
	totalReward := 0

	g.mu.Lock()
	for _, id := range hitIDs {
		t, ok := g.Targets[id]
		if !ok || t.HP <= 0 {
			continue
		}

		// 依目標類型決定擊破機率
		var killChance float64
		switch {
		case t.Multiplier >= 30.0:
			killChance = specialweapon.BlackHoleBossKillChance
		case t.Multiplier >= 10.0:
			killChance = specialweapon.BlackHoleSpecialKillChance
		default:
			killChance = specialweapon.BlackHoleNormalKillChance
		}

		if rand.Float64() >= killChance {
			continue // 未擊破
		}

		// 計算獎勵（黑洞獎勵 = 目標倍率 × betLevel × 0.70，比炸彈 0.5 高，因為費用更高）
		baseReward := int(float64(p.BetLevel) * t.Multiplier * 0.70)
		if baseReward < 1 {
			baseReward = 1
		}

		hitEntries = append(hitEntries, ws.BlackHoleKillEntry{
			InstanceID: t.InstanceID,
			DefID:      t.DefID,
			Multiplier: t.Multiplier,
			Reward:     baseReward,
		})
		totalReward += baseReward

		// 擊破目標
		t.HP = 0
		delete(g.Targets, id)
	}
	g.mu.Unlock()

	// 發放獎勵
	if totalReward > 0 {
		p.Coins += totalReward
		if p.Coins > p.MaxCoins {
			p.MaxCoins = p.Coins
		}
	}

	// 廣播被擊破的目標（讓 Client 播放死亡動畫）
	for _, entry := range hitEntries {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetKill,
			Payload: ws.TargetKillPayload{
				InstanceID: entry.InstanceID,
				KillerID:   p.ID,
				Reward:     entry.Reward,
				Multiplier: entry.Multiplier,
			},
		})
	}

	// Phase 4：廣播最終結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBlackHoleResult,
		Payload: ws.BlackHoleResultPayload{
			ShooterID:   p.ID,
			ShooterName: p.DisplayName,
			Phase:       "result",
			CenterX:     cx,
			CenterY:     cy,
			Radius:      specialweapon.BlackHoleRadius,
			SuckedCount: len(hitIDs),
			HitTargets:  hitEntries,
			TotalReward: totalReward,
			NewBalance:  p.Coins,
			Cost:        cost,
		},
	})

	// 全服公告：擊破 ≥4 個目標時廣播
	if len(hitEntries) >= 4 {
		g.announceBlackHole(p.DisplayName, len(hitEntries), totalReward)
	}

	log.Printf("[BlackHole] player=%s sucked=%d killed=%d reward=%d net=%d",
		p.ID, len(hitIDs), len(hitEntries), totalReward, totalReward-cost)
}

// announceBlackHole 全服公告黑洞漩渦大豐收（DAY-166）
func (g *Game) announceBlackHole(playerName string, killCount, reward int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "black_hole_vortex",
			"message":    fmt.Sprintf("🌀 %s 的黑洞漩渦吸入並擊破了 %d 個目標，獲得 %d 金幣！", playerName, killCount, reward),
			"color":      "#6600CC",
			"duration":   4.0,
			"priority":   3,
		},
	})
}
