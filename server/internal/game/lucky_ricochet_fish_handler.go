// lucky_ricochet_fish_handler.go — 幸運反彈魚系統（DAY-220）
// 業界原創「子彈反彈」機制
//
// 設計：擊破 T178 後觸發「反彈模式」（8秒）：
//   - 玩家的每次射擊在命中目標後，子彈會「反彈」到最近的另一個目標
//   - 反彈範圍：第1跳 200px，第2跳 150px，第3跳 100px（最多 3 跳）
//   - 每次反彈命中：60% 擊破機率，0.55x 倍率
//   - 個人冷卻 18 秒；全服廣播反彈開始/每次反彈/結束
//
// 設計差異：
//   - 與閃電鰻（連鎖跳躍，隨機目標）不同，反彈魚是「射擊觸發的即時反彈」，
//     讓玩家感受到「每一槍都有額外效果」的爽感
//   - 「反彈範圍遞減」讓玩家有「要把魚聚集在一起才能最大化反彈效果」的策略感
//   - 「8 秒反彈模式」讓玩家有「趕快在 8 秒內多打幾槍」的緊迫感
//   - 全服廣播讓其他玩家看到「有人觸發了反彈模式」，製造羨慕感
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	LuckyRicochetPersonalCD = 18 * time.Second // 個人冷卻
	LuckyRicochetDuration   = 8 * time.Second  // 反彈模式持續時間
	LuckyRicochetMaxBounces = 3                // 最大反彈次數
	LuckyRicochetKillChance = 0.60             // 反彈命中擊破機率
	LuckyRicochetKillMult   = 0.55             // 反彈擊破倍率
)

// ricochetRanges 每跳的反彈範圍（px）
var ricochetRanges = []float64{200.0, 150.0, 100.0}

// luckyRicochetFishManager 幸運反彈魚管理器
type luckyRicochetFishManager struct {
	mu sync.Mutex

	// 個人冷卻
	personalCooldowns map[string]time.Time

	// 當前反彈模式狀態（per player）
	activeSessions map[string]*ricochetSession
}

// ricochetSession 單個玩家的反彈模式 session
type ricochetSession struct {
	playerID  string
	expiresAt time.Time
	bounceCount int // 本次射擊已反彈次數（每次射擊重置）
}

func newLuckyRicochetFishManager() *luckyRicochetFishManager {
	return &luckyRicochetFishManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*ricochetSession),
	}
}

// isLuckyRicochetFish 判斷是否為幸運反彈魚
func isLuckyRicochetFish(defID string) bool {
	return defID == "T178"
}

// isRicochetActive 判斷玩家是否在反彈模式中（供 handleAttack 使用）
func (g *Game) isRicochetActive(playerID string) bool {
	mgr := g.LuckyRicochetFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	sess, ok := mgr.activeSessions[playerID]
	if !ok {
		return false
	}
	if time.Now().After(sess.expiresAt) {
		delete(mgr.activeSessions, playerID)
		return false
	}
	return true
}

// tryLuckyRicochetFish 擊破 T178 後觸發反彈模式（供 handleKill 使用）
func (g *Game) tryLuckyRicochetFish(p *player.Player) {
	mgr := g.LuckyRicochetFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok {
		if time.Now().Before(cd) {
			mgr.mu.Unlock()
			return
		}
	}

	// 設定個人冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyRicochetPersonalCD)

	// 啟動反彈模式
	expiresAt := time.Now().Add(LuckyRicochetDuration)
	mgr.activeSessions[p.ID] = &ricochetSession{
		playerID:  p.ID,
		expiresAt: expiresAt,
	}
	mgr.mu.Unlock()

	// 全服廣播：反彈模式開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyRicochetFish,
		Payload: ws.LuckyRicochetFishPayload{
			Event:       "ricochet_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyRicochetDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyRicochetFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("%s 觸發反彈模式！8 秒內每槍最多反彈 3 次！", p.DisplayName),
		"color":   "#FF8C00",
	})
	g.broadcastAnnouncement(ann)

	log.Printf("[LuckyRicochet] player=%s activated ricochet mode for %v", p.ID, LuckyRicochetDuration)

	// 啟動計時器，到期後廣播結束
	go func() {
		time.Sleep(LuckyRicochetDuration)
		mgr.mu.Lock()
		delete(mgr.activeSessions, p.ID)
		mgr.mu.Unlock()

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyRicochetFish,
			Payload: ws.LuckyRicochetFishPayload{
				Event:      "ricochet_end",
				PlayerID:   p.ID,
				PlayerName: p.DisplayName,
			},
		})
		log.Printf("[LuckyRicochet] player=%s ricochet mode ended", p.ID)
	}()
}

// doRicochetBounce 執行反彈（供 handleAttack 在命中後呼叫）
// hitX, hitY 是命中位置；excludeID 是剛被命中的目標（不再反彈到它）
func (g *Game) doRicochetBounce(p *player.Player, hitX, hitY float64, excludeID string, bounceNum int) {
	if bounceNum >= LuckyRicochetMaxBounces {
		return
	}

	// 確認反彈模式仍有效
	if !g.isRicochetActive(p.ID) {
		return
	}

	bounceRange := ricochetRanges[bounceNum]

	// 找最近的目標（在反彈範圍內，排除剛被命中的目標）
	g.mu.RLock()
	var nearestID string
	var nearestDist float64 = math.MaxFloat64
	for id, t := range g.Targets {
		if id == excludeID || t.HP <= 0 {
			continue
		}
		dx := t.X - hitX
		dy := t.Y - hitY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= bounceRange && dist < nearestDist {
			nearestDist = dist
			nearestID = id
		}
	}
	g.mu.RUnlock()

	if nearestID == "" {
		return
	}

	// 執行反彈命中
	g.mu.Lock()
	bounceTarget := g.Targets[nearestID]
	if bounceTarget == nil || bounceTarget.HP <= 0 {
		g.mu.Unlock()
		return
	}

	killed := false
	reward := 0
	if rand.Float64() < LuckyRicochetKillChance {
		// 擊破
		bounceTarget.HP = 0
		killed = true
		avgBet := 1
		if len(g.Players) > 0 {
			total := 0
			for _, pl := range g.Players {
				betDef := data.GetBetDef(pl.BetLevel)
				if betDef != nil {
					total += betDef.BetCost
				}
			}
			avgBet = total / len(g.Players)
		}
		reward = int(bounceTarget.Multiplier * float64(avgBet) * LuckyRicochetKillMult)
		delete(g.Targets, nearestID)
	}
	bounceX := bounceTarget.X
	bounceY := bounceTarget.Y
	g.mu.Unlock()

	// 給觸發玩家獎勵
	if killed && reward > 0 {
		p.AddCoins(reward)
	}

	// 廣播反彈命中
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyRicochetFish,
		Payload: ws.LuckyRicochetFishPayload{
			Event:      "ricochet_bounce",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			BounceNum:  bounceNum + 1,
			TargetID:   nearestID,
			Killed:     killed,
			Reward:     reward,
			X:          bounceX,
			Y:          bounceY,
		},
	})

	log.Printf("[LuckyRicochet] player=%s bounce#%d: target=%s killed=%v reward=%d",
		p.ID, bounceNum+1, nearestID, killed, reward)

	// 繼續反彈（遞迴，下一跳）
	if killed {
		go g.doRicochetBounce(p, bounceX, bounceY, nearestID, bounceNum+1)
	}
}
