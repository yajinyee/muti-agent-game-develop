// fire_storm_fish_handler.go — 火焰風暴魚系統 handler（DAY-176）
// 業界依據：Ocean King 3 Plus「Fire Storm feature — triggers a fire storm that burns multiple
// fish simultaneously, creating chain combustion across the screen」
// 擊破 T134 後觸發「火焰風暴」：場上隨機 4-8 個目標被火焰標記，
// 15 秒內逐一燃燒擊破（每 1.5 秒一個），獎勵 0.6x 倍率
// 設計差異：與漩渦魚（吸入基礎目標）不同，火焰風暴是「隨機標記任意目標」（包含特殊目標），
// 且有「燃燒蔓延」的視覺過程（每 1.5 秒一個），製造「火焰逐漸蔓延」的戲劇感
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
	// FireStormCooldownSec 全服冷卻時間（秒）
	FireStormCooldownSec = 30
	// FireStormDurationSec 火焰風暴持續時間（秒）
	FireStormDurationSec = 15
	// FireStormBurnIntervalMs 每個目標燃燒間隔（ms）
	FireStormBurnIntervalMs = 1500
	// FireStormMinTargets 最少標記目標數
	FireStormMinTargets = 4
	// FireStormMaxTargets 最多標記目標數
	FireStormMaxTargets = 8
	// FireStormKillRewardMult 燃燒擊破獎勵倍率（比直接擊破低，平衡 RTP）
	FireStormKillRewardMult = 0.60
	// FireStormAnnounceThreshold 全服公告門檻（燃燒目標數）
	FireStormAnnounceThreshold = 5
)

// fireStormManager 火焰風暴魚管理器（全服共享冷卻）
type fireStormManager struct {
	mu         sync.Mutex
	isActive   bool
	cooldownAt time.Time
}

// newFireStormManager 建立火焰風暴魚管理器
func newFireStormManager() *fireStormManager {
	return &fireStormManager{}
}

// isFireStormFish 判斷是否為火焰風暴魚（T134）
func isFireStormFish(defID string) bool {
	return defID == "T134"
}

// canTrigger 檢查是否可以觸發（全服冷卻 + 非活躍中）
func (m *fireStormManager) canTrigger() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isActive {
		return false
	}
	return time.Now().After(m.cooldownAt)
}

// setActive 設定活躍狀態
func (m *fireStormManager) setActive(active bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isActive = active
	if !active {
		m.cooldownAt = time.Now().Add(time.Duration(FireStormCooldownSec) * time.Second)
	}
}

// tryFireStormFish 擊破 T134 後觸發火焰風暴（DAY-176）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryFireStormFish(p *player.Player) {
	if !g.FireStorm.canTrigger() {
		return
	}
	g.FireStorm.setActive(true)
	defer g.FireStorm.setActive(false)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 收集場上存活目標（排除 T134 自身和 BOSS）
	g.mu.RLock()
	var candidates []string
	for tid, t := range g.Targets {
		if t.HP > 0 && t.Def.Type != data.TargetTypeBoss && t.Def.ID != "T134" {
			candidates = append(candidates, tid)
		}
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		return
	}

	// 隨機選取 4-8 個目標
	rng.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	count := FireStormMinTargets + rng.Intn(FireStormMaxTargets-FireStormMinTargets+1)
	if count > len(candidates) {
		count = len(candidates)
	}
	burnTargets := candidates[:count]

	log.Printf("[FireStorm] player=%s triggered, burning %d targets", p.ID, count)

	// 廣播火焰風暴開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgFireStormFish,
		Payload: ws.FireStormFishPayload{
			Phase:       "fire_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			TargetCount: count,
			TargetIDs:   burnTargets,
			DurationSec: FireStormDurationSec,
		},
	})

	// 全服公告
	if count >= FireStormAnnounceThreshold {
		g.announceFireStormFish(p.DisplayName, count)
	}

	// 逐一燃燒擊破（每 1.5 秒一個）
	totalReward := 0
	burnedCount := 0
	for _, tid := range burnTargets {
		time.Sleep(time.Duration(FireStormBurnIntervalMs) * time.Millisecond)

		g.mu.Lock()
		t, ok := g.Targets[tid]
		if !ok || t.HP <= 0 {
			g.mu.Unlock()
			// 廣播目標已消失（跳過）
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgFireStormFish,
				Payload: ws.FireStormFishPayload{
					Phase:    "fire_burn",
					TargetID: tid,
					Skipped:  true,
				},
			})
			continue
		}

		// 計算獎勵
		midMult := (t.Def.MultiplierMin + t.Def.MultiplierMax) / 2.0
		reward := int(midMult * float64(p.BetLevel) * FireStormKillRewardMult)
		if reward < 1 {
			reward = 1
		}

		// 擊破目標
		t.HP = 0
		delete(g.Targets, tid)
		g.mu.Unlock()

		// 發放獎勵給觸發玩家
		p.AddReward(reward)
		totalReward += reward
		burnedCount++

		// 廣播燃燒擊破（全服）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgFireStormFish,
			Payload: ws.FireStormFishPayload{
				Phase:    "fire_burn",
				TargetID: tid,
				Reward:   reward,
				Skipped:  false,
			},
		})

		// 廣播目標消失
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetKill,
			Payload: ws.TargetKillPayload{
				InstanceID: tid,
				KillerID:   p.ID,
				Reward:     reward,
			},
		})
	}

	// 廣播火焰風暴結束（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgFireStormFish,
		Payload: ws.FireStormFishPayload{
			Phase:       "fire_end",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			BurnedCount: burnedCount,
			TotalReward: totalReward,
		},
	})

	log.Printf("[FireStorm] player=%s done: burned=%d totalReward=%d",
		p.ID, burnedCount, totalReward)
}

// announceFireStormFish 全服公告火焰風暴魚（DAY-176）
func (g *Game) announceFireStormFish(playerName string, count int) {
	color := "#FF4500" // 橙紅色（火焰感）
	if count >= 7 {
		color = "#FF0000" // 大規模火焰用純紅色
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "fire_storm_fish",
			"message":    fmt.Sprintf("🔥 %s 觸發火焰風暴！%d 個目標正在燃燒！", playerName, count),
			"color":      color,
			"duration":   5.0,
			"priority":   3,
		},
	})
}
