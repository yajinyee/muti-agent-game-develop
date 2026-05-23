// lucky_boomerang_fish_handler.go — 幸運迴旋鏢魚系統（DAY-231）
// 業界原創「迴旋鏢來回穿透」機制
//
// 設計：擊破 T189 後觸發「迴旋鏢模式」（10 秒）：
//   - 玩家的每次射擊發射「迴旋鏢子彈」（替代普通子彈）
//   - 迴旋鏢沿射擊方向飛出，命中目標後「折返」，沿反方向繼續飛行
//   - 折返後再次命中目標繼續折返，最多來回 3 次（去程 + 回程 + 再去程）
//   - 每次命中：70% 擊破機率，0.65x 倍率（個人獎勵）
//   - 個人冷卻 18 秒；全服廣播迴旋鏢開始/每次命中/結束
//
// 設計差異：
//   - 與反彈魚（DAY-220，命中後跳到最近目標，範圍遞減）不同，迴旋鏢是「直線來回穿透」，
//     讓玩家有「要瞄準一排魚」的策略感
//   - 「折返」讓玩家感受到「一槍打多個目標」的爽感
//   - 「最多 3 次折返」確保 RTP 平衡，不會無限穿透
//   - 「10 秒迴旋鏢模式」讓玩家有「趕快在 10 秒內多打幾槍」的緊迫感
//   - 全服廣播讓其他玩家看到「有人觸發了迴旋鏢模式」，製造羨慕感
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
	LuckyBoomerangPersonalCD  = 18 * time.Second // 個人冷卻
	LuckyBoomerangDuration    = 10 * time.Second // 迴旋鏢模式持續時間
	LuckyBoomerangMaxBounces  = 3                // 最大折返次數（去+回+去）
	LuckyBoomerangKillChance  = 0.70             // 每次命中擊破機率
	LuckyBoomerangKillMult    = 0.65             // 每次命中擊破倍率
	LuckyBoomerangSearchRange = 250.0            // 折返後搜尋下一個目標的範圍（px）
)

// luckyBoomerangFishManager 幸運迴旋鏢魚管理器
type luckyBoomerangFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 當前迴旋鏢模式狀態（per player）
	activeSessions map[string]*boomerangSession
}

// boomerangSession 單個玩家的迴旋鏢模式 session
type boomerangSession struct {
	playerID  string
	expiresAt time.Time
}

func newLuckyBoomerangFishManager() *luckyBoomerangFishManager {
	return &luckyBoomerangFishManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*boomerangSession),
	}
}

// isLuckyBoomerangFish 判斷是否為幸運迴旋鏢魚
func isLuckyBoomerangFish(defID string) bool {
	return defID == "T189"
}

// isBoomerangActive 判斷玩家是否在迴旋鏢模式中（供 handleAttack 使用）
func (g *Game) isBoomerangActive(playerID string) bool {
	mgr := g.LuckyBoomerangFish
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

// tryLuckyBoomerangFish 擊破 T189 後觸發迴旋鏢模式（供 handleKill 使用）
func (g *Game) tryLuckyBoomerangFish(p *player.Player) {
	mgr := g.LuckyBoomerangFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}

	// 設定個人冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyBoomerangPersonalCD)

	// 啟動迴旋鏢模式
	expiresAt := time.Now().Add(LuckyBoomerangDuration)
	mgr.activeSessions[p.ID] = &boomerangSession{
		playerID:  p.ID,
		expiresAt: expiresAt,
	}
	mgr.mu.Unlock()

	log.Printf("[LuckyBoomerang] player=%s activated boomerang mode for %v", p.ID, LuckyBoomerangDuration)

	// 全服廣播：迴旋鏢模式開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyBoomerangFish,
		Payload: ws.LuckyBoomerangFishPayload{
			Event:       "boomerang_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyBoomerangDuration.Seconds()),
			MaxBounces:  LuckyBoomerangMaxBounces,
			KillChance:  LuckyBoomerangKillChance,
			KillMult:    LuckyBoomerangKillMult,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyBoomerangFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🪃 %s 觸發迴旋鏢模式！10 秒內每槍最多折返 %d 次！",
			p.DisplayName, LuckyBoomerangMaxBounces),
		"color": "#E67E22",
	})
	g.broadcastAnnouncement(ann)

	// 啟動計時器，到期後廣播結束
	go func() {
		time.Sleep(LuckyBoomerangDuration)
		mgr.mu.Lock()
		delete(mgr.activeSessions, p.ID)
		mgr.mu.Unlock()

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgLuckyBoomerangFish,
			Payload: ws.LuckyBoomerangFishPayload{
				Event:      "boomerang_end",
				PlayerID:   p.ID,
				PlayerName: p.DisplayName,
			},
		})
		log.Printf("[LuckyBoomerang] player=%s boomerang mode ended", p.ID)
	}()
}

// doBoomerangBounce 執行迴旋鏢折返（供 handleAttack 在命中後呼叫）
// hitX, hitY 是命中位置；dirX, dirY 是當前飛行方向（單位向量）；excludeID 是剛被命中的目標
func (g *Game) doBoomerangBounce(p *player.Player, hitX, hitY float64, dirX, dirY float64, excludeID string, bounceNum int) {
	if bounceNum >= LuckyBoomerangMaxBounces {
		return
	}

	// 確認迴旋鏢模式仍有效
	if !g.isBoomerangActive(p.ID) {
		return
	}

	// 折返方向（反轉 X 方向，模擬迴旋鏢折返）
	returnDirX := -dirX
	returnDirY := dirY * 0.3 // Y 方向略微偏移，讓軌跡更自然

	// 正規化方向向量
	length := math.Sqrt(returnDirX*returnDirX + returnDirY*returnDirY)
	if length > 0 {
		returnDirX /= length
		returnDirY /= length
	}

	// 在折返方向上搜尋最近的目標
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
		if dist > LuckyBoomerangSearchRange {
			continue
		}
		// 確認目標在折返方向上（點積 > 0.3，即夾角 < 72 度）
		if dist > 0 {
			dotProduct := (dx/dist)*returnDirX + (dy/dist)*returnDirY
			if dotProduct < 0.3 {
				continue
			}
		}
		if dist < nearestDist {
			nearestDist = dist
			nearestID = id
		}
	}
	g.mu.RUnlock()

	if nearestID == "" {
		// 沒有找到目標，迴旋鏢消失
		log.Printf("[LuckyBoomerang] player=%s bounce#%d: no target found, boomerang fades",
			p.ID, bounceNum+1)
		return
	}

	// 執行命中
	g.mu.Lock()
	bounceTarget := g.Targets[nearestID]
	if bounceTarget == nil || bounceTarget.HP <= 0 {
		g.mu.Unlock()
		return
	}

	killed := false
	reward := 0
	bounceX := bounceTarget.X
	bounceY := bounceTarget.Y

	if rand.Float64() < LuckyBoomerangKillChance {
		// 擊破
		bounceTarget.HP = 0
		killed = true

		// 計算獎勵（個人獎勵，基於觸發玩家的 betCost）
		betDef := data.GetBetDef(p.BetLevel)
		betCost := 1
		if betDef != nil {
			betCost = betDef.BetCost
		}
		reward = int(bounceTarget.Multiplier * float64(betCost) * LuckyBoomerangKillMult)
		if reward < 1 {
			reward = 1
		}
		delete(g.Targets, nearestID)
	}
	g.mu.Unlock()

	// 給觸發玩家獎勵
	if killed && reward > 0 {
		p.AddCoins(reward)
	}

	log.Printf("[LuckyBoomerang] player=%s bounce#%d: target=%s killed=%v reward=%d dir=(%.2f,%.2f)",
		p.ID, bounceNum+1, nearestID, killed, reward, returnDirX, returnDirY)

	// 廣播迴旋鏢命中
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyBoomerangFish,
		Payload: ws.LuckyBoomerangFishPayload{
			Event:      "boomerang_hit",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			BounceNum:  bounceNum + 1,
			TargetID:   nearestID,
			Killed:     killed,
			Reward:     reward,
			X:          bounceX,
			Y:          bounceY,
			DirX:       returnDirX,
			DirY:       returnDirY,
		},
	})

	// 繼續折返（遞迴，下一跳）
	if killed {
		go g.doBoomerangBounce(p, bounceX, bounceY, returnDirX, returnDirY, nearestID, bounceNum+1)
	}
}
