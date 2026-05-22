// dragon_turtle_handler.go — 龍龜不死 Boss 系統 handler（DAY-186）
// 業界依據：Royal Fishing JILI「Immortal Boss mechanic — Golden Toad and Ancient Crocodile
// bosses appear randomly and award consecutive wins ranging from 50X to 150X until they
// leave the screen. This creates extended winning sequences impossible in standard fish games.」
// T144 龍龜不死 Boss 機制：
//   1. 龍龜出現在場上，不會被擊破（HP=99999，永遠不死）
//   2. 每次命中都給獎勵（50-150x betLevel），不扣 HP
//   3. 龍龜在場上移動，直到 Lifetime（30秒）結束離開畫面
//   4. 全服廣播每次命中，讓所有玩家看到「有人在打龍龜」
//   5. 龍龜離開時廣播總結（總命中數/總獎勵）
// 設計差異：
//   - 與 DAY-129 不死 BOSS（隱形/每次射擊有機率命中）不同，
//     龍龜是「可見目標物」（在場上移動），玩家需要主動瞄準
//   - 與普通 BOSS（需要擊破）完全不同，龍龜是「持續收割型」
//   - 玩家不需要擊破，只要命中就有獎勵，製造「穩定收益」的安心感
//   - 全服共享龍龜，所有玩家都可以打，製造「搶打龍龜」的競爭感
package game

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// DragonTurtleHitRewardMin 每次命中最低獎勵倍率
	DragonTurtleHitRewardMin = 50
	// DragonTurtleHitRewardMax 每次命中最高獎勵倍率
	DragonTurtleHitRewardMax = 150
	// DragonTurtleCooldownSec 全服冷卻時間（秒）
	DragonTurtleCooldownSec = 60
	// DragonTurtleAnnounceMinHits 全服公告最低命中數
	DragonTurtleAnnounceMinHits = 5
	// DragonTurtleHitBroadcastIntervalMs 命中廣播節流間隔（毫秒）
	DragonTurtleHitBroadcastIntervalMs = 200
)

// dragonTurtleManager 龍龜不死 Boss 管理器（全服共享）
type dragonTurtleManager struct {
	mu          sync.Mutex
	isActive    bool      // 龍龜是否在場上
	instanceID  string    // 當前龍龜的 InstanceID
	totalHits   int       // 全服總命中數
	totalReward int       // 全服總獎勵
	lastHitAt   time.Time // 上次廣播時間（節流用）
	cooldownEnd time.Time // 冷卻結束時間
}

// newDragonTurtleManager 建立龍龜不死 Boss 管理器
func newDragonTurtleManager() *dragonTurtleManager {
	return &dragonTurtleManager{}
}

// isDragonTurtle 判斷是否為龍龜不死 Boss
func isDragonTurtle(defID string) bool {
	return defID == "T144"
}

// notifyDragonTurtleSpawn 龍龜生成時通知（由 spawnTarget 呼叫）
func (g *Game) notifyDragonTurtleSpawn(instanceID string, x, y float64) {
	g.DragonTurtle.mu.Lock()
	g.DragonTurtle.isActive = true
	g.DragonTurtle.instanceID = instanceID
	g.DragonTurtle.totalHits = 0
	g.DragonTurtle.totalReward = 0
	g.DragonTurtle.lastHitAt = time.Time{}
	g.DragonTurtle.mu.Unlock()

	log.Printf("[DragonTurtle] 龍龜不死 Boss 出現 instance=%s pos=(%.0f,%.0f)", instanceID, x, y)

	// 全服廣播龍龜出現
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgDragonTurtle,
		Payload: ws.DragonTurtlePayload{
			Phase:      "turtle_appear",
			InstanceID: instanceID,
			X:          x,
			Y:          y,
		},
	})

	// 全服公告
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "dragon_turtle_appear",
			"content":    "🐢 龍龜不死 Boss 出現！命中即可獲得 50-150x 獎勵！",
			"color":      "#4CAF50",
			"duration":   5,
			"priority":   6,
		},
	})
}

// notifyDragonTurtleHit 龍龜被命中時處理（由 handleAttack 呼叫）
// 龍龜不死：命中不扣 HP，直接給獎勵
func (g *Game) notifyDragonTurtleHit(p *player.Player, instanceID string) {
	g.DragonTurtle.mu.Lock()
	if !g.DragonTurtle.isActive || g.DragonTurtle.instanceID != instanceID {
		g.DragonTurtle.mu.Unlock()
		return
	}

	// 計算命中獎勵（50-150x betLevel）
	rewardMult := DragonTurtleHitRewardMin + rand.Intn(DragonTurtleHitRewardMax-DragonTurtleHitRewardMin+1)
	reward := rewardMult * p.BetLevel

	g.DragonTurtle.totalHits++
	g.DragonTurtle.totalReward += reward

	totalHits := g.DragonTurtle.totalHits
	totalReward := g.DragonTurtle.totalReward

	// 節流：200ms 內不重複廣播（防止多人同時命中造成廣播風暴）
	now := time.Now()
	shouldBroadcast := now.Sub(g.DragonTurtle.lastHitAt) >= DragonTurtleHitBroadcastIntervalMs*time.Millisecond
	if shouldBroadcast {
		g.DragonTurtle.lastHitAt = now
	}
	g.DragonTurtle.mu.Unlock()

	// 發放獎勵給命中玩家
	p.AddCoins(reward)
	g.sendPlayerUpdate(p)

	// 廣播命中事件（節流，全服可見）
	if shouldBroadcast {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgDragonTurtle,
			Payload: ws.DragonTurtlePayload{
				Phase:       "turtle_hit",
				InstanceID:  instanceID,
				HitterID:    p.ID,
				HitterName:  p.DisplayName,
				HitReward:   reward,
				HitMult:     rewardMult,
				TotalHits:   totalHits,
				TotalReward: totalReward,
			},
		})
	}

	// 個人命中回饋（不節流，讓玩家立刻看到獎勵）
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgDragonTurtle,
		Payload: ws.DragonTurtlePayload{
			Phase:      "my_hit",
			InstanceID: instanceID,
			HitReward:  reward,
			HitMult:    rewardMult,
			TotalHits:  totalHits,
		},
	}); err != nil {
		log.Printf("[DragonTurtle] send my_hit error: %v", err)
	}

	log.Printf("[DragonTurtle] player=%s hit reward=%d (×%d) totalHits=%d",
		p.ID, reward, rewardMult, totalHits)
}

// notifyDragonTurtleLeave 龍龜離開畫面時處理（由 gameLoop 目標超時移除時呼叫）
func (g *Game) notifyDragonTurtleLeave(instanceID string) {
	g.DragonTurtle.mu.Lock()
	if !g.DragonTurtle.isActive || g.DragonTurtle.instanceID != instanceID {
		g.DragonTurtle.mu.Unlock()
		return
	}
	totalHits := g.DragonTurtle.totalHits
	totalReward := g.DragonTurtle.totalReward
	g.DragonTurtle.isActive = false
	g.DragonTurtle.instanceID = ""
	g.DragonTurtle.cooldownEnd = time.Now().Add(DragonTurtleCooldownSec * time.Second)
	g.DragonTurtle.mu.Unlock()

	log.Printf("[DragonTurtle] 龍龜不死 Boss 離開 instance=%s totalHits=%d totalReward=%d",
		instanceID, totalHits, totalReward)

	// 全服廣播龍龜離開
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgDragonTurtle,
		Payload: ws.DragonTurtlePayload{
			Phase:       "turtle_leave",
			InstanceID:  instanceID,
			TotalHits:   totalHits,
			TotalReward: totalReward,
		},
	})

	// 全服公告（≥5 次命中才公告）
	if totalHits >= DragonTurtleAnnounceMinHits {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgAnnouncement,
			Payload: map[string]interface{}{
				"event_type": "dragon_turtle_leave",
				"content":    "🐢 龍龜不死 Boss 離開！全服共命中 " + dragonTurtleItoa(totalHits) + " 次，總獎勵 " + dragonTurtleItoa(totalReward) + " 金幣！",
				"color":      "#FFD700",
				"duration":   6,
				"priority":   5,
			},
		})
	}
}

// dragonTurtleItoa 整數轉字串（避免 import strconv，使用獨立命名防止衝突）
func dragonTurtleItoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	buf := make([]byte, 20)
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
