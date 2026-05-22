// lightning_auto_chain_handler.go — 閃電魚自動連鎖系統 handler（DAY-183）
// 業界依據：Ocean King 3 Monster Awaken「Lightning Fish — Catching a Lightning Fish will
// trigger a Lightning Chain. Lightning Chain will continue to catch fish automatically
// until time runs out.」
// 擊破 T141 後觸發「閃電自動連鎖」：
//   1. 系統自動每 0.5 秒選一個隨機目標發射閃電
//   2. 持續 8 秒（共最多 16 次自動攻擊）
//   3. 每次有 65% 擊破機率，獎勵 0.60x 倍率
//   4. 全服廣播每次自動攻擊，讓所有玩家看到「閃電在自動收割」
// 設計差異：
//   - 與 T103 閃電鰻（手動觸發，5跳）不同，閃電魚是「全自動持續連鎖」（8秒/16次），
//     玩家不需要操作，純粹享受「自動收割」的爽感
//   - 與 T139 雷霆鯊魚（手動跳躍，20跳）不同，閃電魚是「時間驅動」（每 0.5 秒一次），
//     讓玩家感受到「閃電在持續不斷地攻擊」
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// LightningAutoChainDurationSec 自動連鎖持續時間（秒）
	LightningAutoChainDurationSec = 8
	// LightningAutoChainIntervalMs 每次自動攻擊間隔（ms）
	LightningAutoChainIntervalMs = 500
	// LightningAutoChainKillChance 每次攻擊擊破機率（65%）
	LightningAutoChainKillChance = 0.65
	// LightningAutoChainRewardMult 擊破獎勵倍率
	LightningAutoChainRewardMult = 0.60
	// LightningAutoChainCooldownSec 全服冷卻時間（秒）
	LightningAutoChainCooldownSec = 30
	// LightningAutoChainAnnounceMinKills 全服公告最低擊破數
	LightningAutoChainAnnounceMinKills = 6
)

// lightningAutoChainManager 閃電魚自動連鎖管理器（全服共享）
type lightningAutoChainManager struct {
	mu          sync.Mutex
	isActive    bool
	cooldownEnd time.Time
}

// newLightningAutoChainManager 建立閃電魚自動連鎖管理器
func newLightningAutoChainManager() *lightningAutoChainManager {
	return &lightningAutoChainManager{}
}

// isLightningAutoFish 判斷是否為閃電魚（T141）
func isLightningAutoFish(defID string) bool {
	return defID == "T141"
}

// isOnCooldown 檢查是否在全服冷卻中
func (m *lightningAutoChainManager) isOnCooldown() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return time.Now().Before(m.cooldownEnd)
}

// activate 激活自動連鎖
func (m *lightningAutoChainManager) activate() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isActive {
		return false
	}
	m.isActive = true
	m.cooldownEnd = time.Now().Add(time.Duration(LightningAutoChainCooldownSec) * time.Second)
	return true
}

// deactivate 結束自動連鎖
func (m *lightningAutoChainManager) deactivate() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isActive = false
}

// tryLightningAutoChain 擊破 T141 後觸發閃電自動連鎖（DAY-183）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryLightningAutoChain(p *player.Player, triggerID string, triggerX, triggerY float64) {
	// 全服冷卻檢查
	if g.LightningAutoChain.isOnCooldown() {
		return
	}
	if !g.LightningAutoChain.activate() {
		return // 已有其他連鎖在進行
	}
	defer g.LightningAutoChain.deactivate()

	log.Printf("[LightningAutoChain] player=%s triggered auto chain", p.ID)

	// 廣播自動連鎖開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLightningAutoChain,
		Payload: ws.LightningAutoChainPayload{
			Phase:       "chain_start",
			TriggerID:   triggerID,
			TriggerX:    triggerX,
			TriggerY:    triggerY,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			DurationSec: LightningAutoChainDurationSec,
		},
	})

	totalReward := 0
	totalKills := 0
	totalAttacks := 0
	maxAttacks := LightningAutoChainDurationSec * 1000 / LightningAutoChainIntervalMs

	for attack := 0; attack < maxAttacks; attack++ {
		time.Sleep(time.Duration(LightningAutoChainIntervalMs) * time.Millisecond)

		// 從全場隨機選一個存活目標
		g.mu.RLock()
		var candidates []struct {
			instanceID string
			defID      string
			x, y       float64
			multiplier float64
		}
		for id, t := range g.Targets {
			if id == triggerID || t.HP <= 0 || t.DefID == "B001" {
				continue
			}
			candidates = append(candidates, struct {
				instanceID string
				defID      string
				x, y       float64
				multiplier float64
			}{t.InstanceID, t.DefID, t.X, t.Y, t.Multiplier})
		}
		g.mu.RUnlock()

		if len(candidates) == 0 {
			break // 沒有目標，結束
		}

		// 隨機選一個目標
		dt := candidates[rand.Intn(len(candidates))]
		totalAttacks++

		// 廣播自動攻擊（全服）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLightningAutoChain,
			Payload: ws.LightningAutoChainPayload{
				Phase:      fmt.Sprintf("auto_%d", attack+1),
				TargetID:   dt.instanceID,
				TargetX:    dt.x,
				TargetY:    dt.y,
				AttackNum:  attack + 1,
				KillerID:   p.ID,
			},
		})

		// 65% 機率擊破
		if rand.Float64() >= LightningAutoChainKillChance {
			continue
		}

		// 擊破目標
		g.mu.Lock()
		t, ok := g.Targets[dt.instanceID]
		if !ok || t.HP <= 0 {
			g.mu.Unlock()
			continue
		}
		reward := int(float64(p.BetLevel) * dt.multiplier * LightningAutoChainRewardMult)
		if reward < 1 {
			reward = 1
		}
		t.HP = 0
		delete(g.Targets, dt.instanceID)
		g.mu.Unlock()

		totalReward += reward
		totalKills++

		// 廣播目標擊破
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetKill,
			Payload: ws.TargetKillPayload{
				InstanceID: dt.instanceID,
				KillerID:   p.ID,
				Reward:     reward,
				Multiplier: dt.multiplier,
			},
		})

		log.Printf("[LightningAutoChain] auto[%d] target=%s mult=%.0f reward=%d",
			attack+1, dt.instanceID, dt.multiplier, reward)
	}

	if totalReward > 0 {
		p.AddReward(totalReward)
	}

	// 廣播自動連鎖結果（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLightningAutoChain,
		Payload: ws.LightningAutoChainPayload{
			Phase:       "result",
			TriggerID:   triggerID,
			TotalAttacks: totalAttacks,
			TotalKills:  totalKills,
			TotalReward: totalReward,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
		},
	})

	// 個人結果通知
	if totalReward > 0 {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgReward,
			Payload: ws.RewardPayload{
				Source:     "lightning_auto_chain",
				Amount:     totalReward,
				Multiplier: float64(totalKills),
				NewBalance: p.Coins,
			},
		})
	}

	// 全服公告：擊破 ≥ 6 個
	if totalKills >= LightningAutoChainAnnounceMinKills {
		g.announceLightningAutoChain(p.DisplayName, totalAttacks, totalKills, totalReward)
	}

	log.Printf("[LightningAutoChain] player=%s attacks=%d kills=%d total_reward=%d",
		p.ID, totalAttacks, totalKills, totalReward)
}

// announceLightningAutoChain 全服公告閃電魚自動連鎖（DAY-183）
func (g *Game) announceLightningAutoChain(playerName string, attacks, kills, reward int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "lightning_auto_chain",
			"message":    fmt.Sprintf("⚡ %s 的閃電魚自動連鎖！%d 次攻擊擊破 %d 個目標！獲得 %d 金幣！", playerName, attacks, kills, reward),
			"color":      "#FFFF00", // 亮黃色（閃電感）
			"duration":   5.0,
			"priority":   3,
		},
	})
}
