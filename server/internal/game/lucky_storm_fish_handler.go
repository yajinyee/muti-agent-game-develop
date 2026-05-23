// lucky_storm_fish_handler.go — 幸運風暴魚系統（DAY-230）
// 業界原創「風暴旋轉+位置混亂」機制
//
// 設計：擊破 T188 後在場上建立「風暴中心」（持續 10 秒）：
//   - 風暴建立在場景中央附近（隨機偏移），半徑 320px
//   - 風暴範圍內所有目標每 1.5 秒被「風暴旋轉」（隨機傳送到風暴範圍內的新位置）
//   - 風暴範圍內目標被擊破：獎勵 ×2.5 倍率加成（乘法）
//   - 10 秒後「風暴爆發」：範圍內所有目標 80% 擊破機率（0.75x 倍率，全服共享）
//   - 個人冷卻 22 秒；全服冷卻 35 秒
//
// 設計差異：
//   - 與黑洞魚（DAY-221，重力吸引+奇點爆炸，HP 損失）不同，風暴魚是「位置混亂+旋轉」，讓玩家有「要趁目標在風暴範圍內快打」的緊迫感
//   - 與傳送魚（DAY-223，全場瞬間移動，傳送混亂加成）不同，風暴魚是「局部風暴範圍內旋轉」，讓玩家有「要站在風暴範圍內打」的空間策略感
//   - 「風暴旋轉」讓目標在範圍內隨機移動，製造「混亂爽感」
//   - 「風暴爆發」讓玩家有「等待→爆發」的高潮設計
//   - 全服廣播風暴位置讓所有玩家都往同一個地方打，製造「全服聚焦」的社交感
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyStormPersonalCD    = 22 * time.Second // 個人冷卻
	LuckyStormGlobalCD      = 35 * time.Second // 全服冷卻
	LuckyStormDuration      = 10 * time.Second // 風暴持續時間
	LuckyStormRotateInterval = 1500 * time.Millisecond // 旋轉間隔（1.5 秒）
	LuckyStormRadius        = 320.0            // 風暴半徑（px）
	LuckyStormKillMult      = 2.5              // 風暴範圍內擊破倍率
	LuckyStormBlastChance   = 0.80             // 風暴爆發擊破機率
	LuckyStormBlastMult     = 0.75             // 風暴爆發倍率
)

// luckyStormFishManager 幸運風暴魚管理器
type luckyStormFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldown map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 風暴狀態
	active     bool
	stormX     float64   // 風暴中心 X
	stormY     float64   // 風暴中心 Y
	activeUntil time.Time
	instanceID string
}

func newLuckyStormFishManager() *luckyStormFishManager {
	return &luckyStormFishManager{
		personalCooldown: make(map[string]time.Time),
	}
}

// isLuckyStormFish 判斷是否為幸運風暴魚
func isLuckyStormFish(defID string) bool {
	return defID == "T188"
}

// isInStorm 判斷目標是否在風暴範圍內（供 handleKill 使用）
func (g *Game) isInStorm(targetID string) bool {
	mgr := g.LuckyStormFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if !mgr.active || time.Now().After(mgr.activeUntil) {
		return false
	}
	g.mu.RLock()
	t, ok := g.Targets[targetID]
	g.mu.RUnlock()
	if !ok {
		return false
	}
	dx := t.X - mgr.stormX
	dy := t.Y - mgr.stormY
	return math.Sqrt(dx*dx+dy*dy) <= LuckyStormRadius
}

// getLuckyStormMultiplier 取得風暴範圍內擊破倍率（供 handleKill 使用）
func (g *Game) getLuckyStormMultiplier(targetID string) float64 {
	if g.isInStorm(targetID) {
		return LuckyStormKillMult
	}
	return 1.0
}

// tryLuckyStormFish 擊破 T188 後觸發風暴（供 handleKill 使用）
func (g *Game) tryLuckyStormFish(p *player.Player) {
	mgr := g.LuckyStormFish
	mgr.mu.Lock()

	// 全服冷卻檢查
	if time.Now().Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}
	// 個人冷卻檢查
	if until, ok := mgr.personalCooldown[p.ID]; ok && time.Now().Before(until) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldown[p.ID] = time.Now().Add(LuckyStormPersonalCD)
	mgr.globalCooldownUntil = time.Now().Add(LuckyStormGlobalCD)

	// 建立風暴中心（場景中央附近隨機偏移）
	centerX := 500.0 + (rand.Float64()-0.5)*200.0 // 400-600
	centerY := 300.0 + (rand.Float64()-0.5)*100.0 // 250-350

	instanceID := fmt.Sprintf("storm_%d", time.Now().UnixNano())
	mgr.active = true
	mgr.stormX = centerX
	mgr.stormY = centerY
	mgr.activeUntil = time.Now().Add(LuckyStormDuration)
	mgr.instanceID = instanceID
	mgr.mu.Unlock()

	log.Printf("[LuckyStorm] player=%s triggered storm at (%.0f,%.0f) instance=%s",
		p.ID, centerX, centerY, instanceID)

	// 全服廣播風暴開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyStormFish,
		Payload: ws.LuckyStormFishPayload{
			Event:       "storm_start",
			InstanceID:  instanceID,
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			StormX:      centerX,
			StormY:      centerY,
			Radius:      LuckyStormRadius,
			DurationSec: int(LuckyStormDuration.Seconds()),
			KillMult:    LuckyStormKillMult,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyStormFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌪️ %s 觸發風暴！風暴範圍內目標擊破獲得 ×%.1f！",
			p.DisplayName, LuckyStormKillMult),
		"color": "#1ABC9C",
	})
	g.broadcastAnnouncement(ann)

	// 啟動風暴 goroutine
	go g.runLuckyStorm(instanceID, centerX, centerY)
}

// runLuckyStorm 風暴主 goroutine：每 1.5 秒旋轉目標，10 秒後爆發
func (g *Game) runLuckyStorm(instanceID string, centerX, centerY float64) {
	ticker := time.NewTicker(LuckyStormRotateInterval)
	defer ticker.Stop()

	deadline := time.Now().Add(LuckyStormDuration)
	rotateCount := 0

	for {
		select {
		case <-ticker.C:
			if time.Now().After(deadline) {
				ticker.Stop()
				goto blast
			}
			rotateCount++
			g.doStormRotate(instanceID, centerX, centerY, rotateCount)
		}
	}

blast:
	// 確認風暴仍然是同一個 instance
	mgr := g.LuckyStormFish
	mgr.mu.Lock()
	if mgr.instanceID != instanceID {
		mgr.mu.Unlock()
		return
	}
	mgr.active = false
	mgr.mu.Unlock()

	g.doStormBlast(instanceID, centerX, centerY)
}

// doStormRotate 執行一次風暴旋轉（隨機傳送範圍內目標到新位置）
func (g *Game) doStormRotate(instanceID string, centerX, centerY float64, rotateCount int) {
	g.mu.Lock()
	type movedTarget struct {
		ID   string  `json:"id"`
		NewX float64 `json:"new_x"`
		NewY float64 `json:"new_y"`
	}
	moved := make([]movedTarget, 0)

	for id, t := range g.Targets {
		if !t.IsAlive || t.Def.Type == "boss" {
			continue
		}
		dx := t.X - centerX
		dy := t.Y - centerY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > LuckyStormRadius {
			continue
		}
		// 在風暴範圍內隨機新位置（保持在範圍內）
		angle := rand.Float64() * 2 * math.Pi
		r := rand.Float64() * LuckyStormRadius * 0.85 // 不要太靠邊緣
		newX := centerX + r*math.Cos(angle)
		newY := centerY + r*math.Sin(angle)
		// 邊界限制
		if newX < 80 {
			newX = 80
		}
		if newX > 920 {
			newX = 920
		}
		if newY < 60 {
			newY = 60
		}
		if newY > 540 {
			newY = 540
		}
		t.X = newX
		t.Y = newY
		moved = append(moved, movedTarget{ID: id, NewX: newX, NewY: newY})
	}
	g.mu.Unlock()

	if len(moved) == 0 {
		return
	}

	log.Printf("[LuckyStorm] rotate #%d: moved %d targets", rotateCount, len(moved))

	// 廣播旋轉結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyStormFish,
		Payload: ws.LuckyStormFishPayload{
			Event:       "storm_rotate",
			InstanceID:  instanceID,
			RotateCount: rotateCount,
			MovedCount:  len(moved),
			Targets:     moved,
		},
	})
}

// doStormBlast 風暴爆發：範圍內所有目標 80% 擊破機率
func (g *Game) doStormBlast(instanceID string, centerX, centerY float64) {
	log.Printf("[LuckyStorm] blast triggered instance=%s", instanceID)

	g.mu.Lock()
	type blastTarget struct {
		ID      string  `json:"id"`
		Killed  bool    `json:"killed"`
		Reward  int     `json:"reward"`
	}
	blastResults := make([]blastTarget, 0)
	totalReward := 0
	killedCount := 0

	for id, t := range g.Targets {
		if !t.IsAlive || t.Def.Type == "boss" {
			continue
		}
		dx := t.X - centerX
		dy := t.Y - centerY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > LuckyStormRadius {
			continue
		}
		// 80% 擊破機率
		if rand.Float64() < LuckyStormBlastChance {
			reward := int(float64(t.Def.MultiplierMax) * LuckyStormBlastMult)
			t.IsAlive = false
			t.HP = 0
			totalReward += reward
			killedCount++
			blastResults = append(blastResults, blastTarget{ID: id, Killed: true, Reward: reward})
		} else {
			blastResults = append(blastResults, blastTarget{ID: id, Killed: false})
		}
	}
	g.mu.Unlock()

	log.Printf("[LuckyStorm] blast killed=%d totalReward=%d", killedCount, totalReward)

	// 廣播風暴爆發
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyStormFish,
		Payload: ws.LuckyStormFishPayload{
			Event:       "storm_blast",
			InstanceID:  instanceID,
			KilledCount: killedCount,
			TotalReward: totalReward,
			BlastMult:   LuckyStormBlastMult,
			Targets:     blastResults,
		},
	})

	// 全服公告（≥3 個擊破時）
	if killedCount >= 3 {
		color := "#1ABC9C"
		if killedCount >= 6 {
			color = "#00FF7F"
		}
		ann := g.Announce.Create(announce.EventLuckyStormFish, "", 0, map[string]string{
			"message": fmt.Sprintf("🌪️ 風暴爆發！%d 個目標被摧毀！獎勵 %d！",
				killedCount, totalReward),
			"color": color,
		})
		g.broadcastAnnouncement(ann)
	}
}
