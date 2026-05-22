// phoenix_fish_handler.go — 鳳凰魚涅槃重生系統 handler（DAY-185）
// 業界依據：Ocean King 3 Plus「Phoenix Fish — when defeated, the Phoenix Fish triggers a
// rebirth explosion that deals massive damage to all fish on screen, with the Phoenix
// rising from the ashes to grant a 30-second luck boost」
// 擊破 T143 後觸發「涅槃爆炸」：
//   1. 場上所有目標受到爆炸傷害（普通 80% 擊破，特殊 50%，BOSS 20%）
//   2. 爆炸後「鳳凰重生」：全服 30 秒內所有擊破獎勵 +30%（加法疊加）
//   3. 全服廣播爆炸效果，讓所有玩家看到「鳳凰在全場爆炸」
// 設計差異：
//   - 與隕石魚（隨機選目標/數量驅動）不同，鳳凰魚是「全場同時爆炸」（一次性清場），
//     讓玩家感受到「鳳凰涅槃，全場燃燒」的壯觀感
//   - 與黃金鯊魚（全服 ×1.5，12秒）不同，鳳凰魚的重生加成是「+30%（加法）」且持續 30 秒，
//     讓玩家在爆炸後的「重生期」內持續享受加成
//   - 兩段式設計（爆炸→重生）讓玩家有「先爽一波，再持續爽」的雙重滿足感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// PhoenixFishNormalKillChance 普通目標爆炸擊破機率（80%）
	PhoenixFishNormalKillChance = 0.80
	// PhoenixFishSpecialKillChance 特殊目標爆炸擊破機率（50%）
	PhoenixFishSpecialKillChance = 0.50
	// PhoenixFishBossKillChance BOSS 爆炸擊破機率（20%）
	PhoenixFishBossKillChance = 0.20
	// PhoenixFishRewardMult 爆炸擊破獎勵倍率
	PhoenixFishRewardMult = 0.55
	// PhoenixFishRebirthBoost 重生加成（+30%，加法）
	PhoenixFishRebirthBoost = 0.30
	// PhoenixFishRebirthDurationSec 重生加成持續時間（秒）
	PhoenixFishRebirthDurationSec = 30
	// PhoenixFishCooldownSec 全服冷卻時間（秒）
	PhoenixFishCooldownSec = 45
	// PhoenixFishAnnounceMinKills 全服公告最低擊破數
	PhoenixFishAnnounceMinKills = 5
)

// phoenixFishManager 鳳凰魚管理器（全服共享）
type phoenixFishManager struct {
	mu          sync.Mutex
	isActive    bool      // 爆炸是否在進行
	rebirthEnd  time.Time // 重生加成結束時間
	cooldownEnd time.Time
}

// newPhoenixFishManager 建立鳳凰魚管理器
func newPhoenixFishManager() *phoenixFishManager {
	return &phoenixFishManager{}
}

// isPhoenixFish 判斷是否為鳳凰魚（T143）
func isPhoenixFish(defID string) bool {
	return defID == "T143"
}

// isOnCooldownPhoenix 檢查是否在全服冷卻中
func (m *phoenixFishManager) isOnCooldownPhoenix() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return time.Now().Before(m.cooldownEnd)
}

// activatePhoenix 激活鳳凰爆炸
func (m *phoenixFishManager) activatePhoenix() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isActive {
		return false
	}
	m.isActive = true
	m.cooldownEnd = time.Now().Add(time.Duration(PhoenixFishCooldownSec) * time.Second)
	return true
}

// deactivatePhoenix 結束爆炸，激活重生加成
func (m *phoenixFishManager) deactivatePhoenix() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isActive = false
	m.rebirthEnd = time.Now().Add(time.Duration(PhoenixFishRebirthDurationSec) * time.Second)
}

// getRebirthBoost 取得重生加成（供 handleKill 使用）
// 重生期間回傳 0.30，否則回傳 0.0
func (m *phoenixFishManager) getRebirthBoost() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if time.Now().Before(m.rebirthEnd) {
		return PhoenixFishRebirthBoost
	}
	return 0.0
}

// tryPhoenixFishRebirth 擊破 T143 後觸發涅槃爆炸（DAY-185）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryPhoenixFishRebirth(p *player.Player, triggerID string, triggerX, triggerY float64) {
	// 全服冷卻檢查
	if g.PhoenixFish.isOnCooldownPhoenix() {
		return
	}
	if !g.PhoenixFish.activatePhoenix() {
		return // 已有其他爆炸在進行
	}

	log.Printf("[PhoenixFish] player=%s triggered phoenix rebirth", p.ID)

	// 廣播涅槃爆炸開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgPhoenixFish,
		Payload: ws.PhoenixFishPayload{
			Phase:      "phoenix_explode",
			TriggerID:  triggerID,
			TriggerX:   triggerX,
			TriggerY:   triggerY,
			KillerID:   p.ID,
			KillerName: p.DisplayName,
		},
	})

	// 短暫延遲讓 Client 播放爆炸動畫
	time.Sleep(400 * time.Millisecond)

	// 收集場上所有目標
	g.mu.RLock()
	type candidate struct {
		instanceID string
		defID      string
		x, y       float64
		multiplier float64
		targetType data.TargetType
	}
	var candidates []candidate
	for id, t := range g.Targets {
		if id == triggerID || t.HP <= 0 {
			continue
		}
		candidates = append(candidates, candidate{
			t.InstanceID, t.DefID, t.X, t.Y, t.Multiplier, t.Def.Type,
		})
	}
	g.mu.RUnlock()

	totalReward := 0
	totalKills := 0

	// 對所有目標進行爆炸傷害
	for _, dt := range candidates {
		// 判斷擊破機率
		var killChance float64
		switch dt.targetType {
		case data.TargetTypeBoss:
			killChance = PhoenixFishBossKillChance
		case data.TargetTypeSpecial:
			killChance = PhoenixFishSpecialKillChance
		default:
			killChance = PhoenixFishNormalKillChance
		}

		if rand.Float64() >= killChance {
			continue // 未擊破
		}

		// 擊破目標
		g.mu.Lock()
		t, ok := g.Targets[dt.instanceID]
		if !ok || t.HP <= 0 {
			g.mu.Unlock()
			continue
		}
		reward := int(float64(p.BetLevel) * dt.multiplier * PhoenixFishRewardMult)
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
	}

	if totalReward > 0 {
		p.AddReward(totalReward)
	}

	// 結束爆炸，激活重生加成
	g.PhoenixFish.deactivatePhoenix()

	// 廣播鳳凰重生（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgPhoenixFish,
		Payload: ws.PhoenixFishPayload{
			Phase:       "phoenix_rebirth",
			TriggerID:   triggerID,
			TotalKills:  totalKills,
			TotalReward: totalReward,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			BoostPct:    int(PhoenixFishRebirthBoost * 100),
			BoostSec:    PhoenixFishRebirthDurationSec,
		},
	})

	// 個人結果通知
	if totalReward > 0 {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgReward,
			Payload: ws.RewardPayload{
				Source:     "phoenix_fish",
				Amount:     totalReward,
				Multiplier: float64(totalKills),
				NewBalance: p.Coins,
			},
		})
	}

	// 全服公告：擊破 ≥ 5 個
	if totalKills >= PhoenixFishAnnounceMinKills {
		g.announcePhoenixFish(p.DisplayName, totalKills, totalReward)
	}

	log.Printf("[PhoenixFish] player=%s kills=%d total_reward=%d rebirth_boost=+30%% for %ds",
		p.ID, totalKills, totalReward, PhoenixFishRebirthDurationSec)

	// 30 秒後廣播重生加成結束
	time.Sleep(time.Duration(PhoenixFishRebirthDurationSec) * time.Second)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgPhoenixFish,
		Payload: ws.PhoenixFishPayload{
			Phase:      "rebirth_end",
			KillerID:   p.ID,
			KillerName: p.DisplayName,
		},
	})
}

// getPhoenixRebirthBoost 取得鳳凰重生加成（供 handleKill 使用）
func (g *Game) getPhoenixRebirthBoost() float64 {
	if g.PhoenixFish == nil {
		return 0.0
	}
	return g.PhoenixFish.getRebirthBoost()
}

// announcePhoenixFish 全服公告鳳凰魚涅槃重生（DAY-185）
func (g *Game) announcePhoenixFish(playerName string, kills, reward int) {
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "phoenix_fish",
			"message":    fmt.Sprintf("🔥 %s 的鳳凰魚涅槃爆炸！擊破 %d 個目標！獲得 %d 金幣！全服重生加成 +30%% 持續 30 秒！", playerName, kills, reward),
			"color":      "#FF4400", // 火焰橙紅色
			"duration":   6.0,
			"priority":   4,
		},
	})
}
