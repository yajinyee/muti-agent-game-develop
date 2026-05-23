// lucky_clone_fish_handler.go — 幸運分身魚系統（DAY-242）
// 業界原創「三方向同時射擊」機制
//
// 設計：擊破 T200 後觸發「分身模式」（8 秒）：
//   - 玩家的每次射擊同時產生 2 個「分身子彈」，分別向左右各偏移 30 度飛出
//   - 分身子彈命中目標：60% 擊破機率，×0.7 倍率（個人獎勵）
//   - 分身子彈搜尋範圍：偏移方向 300px 內最近目標
//   - 個人冷卻 20 秒；全服廣播分身模式開始/每次分身命中/結束
//
// 設計差異：
//   - 與迴旋鏢魚（T189，折返穿透）不同，分身魚是「同時三方向射擊」，讓玩家有「一槍打三個方向」的爽感
//   - 與充能魚（T183，累積爆發）不同，分身魚是「持續輔助」，每槍都有額外收益
//   - 「偏移 30 度」讓分身子彈打到不同目標，最大化覆蓋範圍
//   - 「60% 擊破機率 × 0.7 倍率」確保 RTP 平衡（期望值 = 0.6 × 0.7 × 2 = 0.84 額外收益/槍）
//   - 全服廣播讓其他玩家看到「有人觸發了分身模式」，製造羨慕感
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
	LuckyCloneFishPersonalCD  = 20 * time.Second // 個人冷卻
	LuckyCloneFishDuration    = 8 * time.Second  // 分身模式持續時間
	LuckyCloneFishAngleDeg    = 30.0             // 分身子彈偏移角度（度）
	LuckyCloneFishKillChance  = 0.60             // 分身子彈擊破機率
	LuckyCloneFishKillMult    = 0.70             // 分身子彈擊破倍率
	LuckyCloneFishSearchRange = 300.0            // 分身子彈搜尋範圍（px）
)

// cloneFishSession 單個玩家的分身模式 session
type cloneFishSession struct {
	playerID  string
	expiresAt time.Time
}

// luckyCloneFishManager 幸運分身魚管理器
type luckyCloneFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 當前分身模式狀態（per player）
	activeSessions map[string]*cloneFishSession
}

func newLuckyCloneFishManager() *luckyCloneFishManager {
	return &luckyCloneFishManager{
		personalCooldowns: make(map[string]time.Time),
		activeSessions:    make(map[string]*cloneFishSession),
	}
}

// isLuckyCloneFish 判斷是否為幸運分身魚
func isLuckyCloneFish(defID string) bool {
	return defID == "T200"
}

// isCloneModeActive 判斷玩家是否在分身模式中（供 handleAttack 使用）
func (g *Game) isCloneModeActive(playerID string) bool {
	mgr := g.LuckyCloneFish
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

// tryLuckyCloneFish 擊破 T200 後觸發分身模式
func (g *Game) tryLuckyCloneFish(p *player.Player) {
	mgr := g.LuckyCloneFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}

	// 設定個人冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyCloneFishPersonalCD)

	// 啟動分身模式
	expiresAt := time.Now().Add(LuckyCloneFishDuration)
	mgr.activeSessions[p.ID] = &cloneFishSession{
		playerID:  p.ID,
		expiresAt: expiresAt,
	}
	mgr.mu.Unlock()

	log.Printf("[LuckyClone] player=%s activated clone mode for %v", p.ID, LuckyCloneFishDuration)

	// 個人訊息：分身模式開始
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgLuckyCloneFish,
		Payload: ws.LuckyCloneFishPayload{
			Event:       "clone_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyCloneFishDuration.Seconds()),
			AngleDeg:    LuckyCloneFishAngleDeg,
			KillChance:  LuckyCloneFishKillChance,
			KillMult:    LuckyCloneFishKillMult,
		},
	})

	// 全服廣播：分身模式開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyCloneFish,
		Payload: ws.LuckyCloneFishPayload{
			Event:       "clone_broadcast",
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyCloneFishDuration.Seconds()),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyCloneFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("👥 %s 觸發分身模式！%d 秒內每槍同時發射 3 個方向！",
			p.DisplayName, int(LuckyCloneFishDuration.Seconds())),
		"color": "#8E44AD",
	})
	g.broadcastAnnouncement(ann)

	// 啟動計時器，到期後廣播結束
	go func() {
		time.Sleep(LuckyCloneFishDuration)
		mgr.mu.Lock()
		delete(mgr.activeSessions, p.ID)
		mgr.mu.Unlock()

		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgLuckyCloneFish,
			Payload: ws.LuckyCloneFishPayload{
				Event:      "clone_end",
				PlayerID:   p.ID,
				PlayerName: p.DisplayName,
			},
		})
		log.Printf("[LuckyClone] player=%s clone mode ended", p.ID)
	}()
}

// doCloneShots 執行分身子彈（供 handleAttack 在射擊後呼叫）
// originX, originY 是射擊起點（砲台位置，通常在左側）
// targetX, targetY 是主要目標位置
// excludeID 是主要目標（避免重複命中）
func (g *Game) doCloneShots(p *player.Player, originX, originY, targetX, targetY float64, excludeID string) {
	// 確認分身模式仍有效
	if !g.isCloneModeActive(p.ID) {
		return
	}

	// 計算主要射擊方向
	dx := targetX - originX
	dy := targetY - originY
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist < 1.0 {
		// 無法計算方向，使用預設向右
		dx = 1.0
		dy = 0.0
		dist = 1.0
	}
	mainDirX := dx / dist
	mainDirY := dy / dist

	// 計算兩個偏移方向（+30度 和 -30度）
	angleRad := LuckyCloneFishAngleDeg * math.Pi / 180.0

	// 左偏 +30 度
	leftDirX := mainDirX*math.Cos(angleRad) - mainDirY*math.Sin(angleRad)
	leftDirY := mainDirX*math.Sin(angleRad) + mainDirY*math.Cos(angleRad)

	// 右偏 -30 度
	rightDirX := mainDirX*math.Cos(-angleRad) - mainDirY*math.Sin(-angleRad)
	rightDirY := mainDirX*math.Sin(-angleRad) + mainDirY*math.Cos(-angleRad)

	// 發射兩個分身子彈
	go g.fireCloneShot(p, originX, originY, leftDirX, leftDirY, excludeID, "left")
	go g.fireCloneShot(p, originX, originY, rightDirX, rightDirY, excludeID, "right")
}

// fireCloneShot 發射單個分身子彈
func (g *Game) fireCloneShot(p *player.Player, originX, originY, dirX, dirY float64, excludeID string, side string) {
	// 在偏移方向上搜尋最近的目標
	g.mu.RLock()
	var nearestID string
	var nearestDist float64 = math.MaxFloat64
	var nearestX, nearestY float64

	for id, t := range g.Targets {
		if id == excludeID || t.HP <= 0 || t.DefID == "B001" {
			continue
		}
		dx := t.X - originX
		dy := t.Y - originY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > LuckyCloneFishSearchRange {
			continue
		}
		// 確認目標在偏移方向上（點積 > 0.25，即夾角 < 75 度）
		if dist > 0 {
			dotProduct := (dx/dist)*dirX + (dy/dist)*dirY
			if dotProduct < 0.25 {
				continue
			}
		}
		if dist < nearestDist {
			nearestDist = dist
			nearestID = id
			nearestX = t.X
			nearestY = t.Y
		}
	}
	g.mu.RUnlock()

	if nearestID == "" {
		// 沒有找到目標，分身子彈消失
		return
	}

	// 執行命中
	g.mu.Lock()
	cloneTarget := g.Targets[nearestID]
	if cloneTarget == nil || cloneTarget.HP <= 0 {
		g.mu.Unlock()
		return
	}

	killed := false
	reward := 0

	if rand.Float64() < LuckyCloneFishKillChance {
		// 擊破
		cloneTarget.HP = 0
		killed = true

		// 計算獎勵（個人獎勵）
		betDef := data.GetBetDef(p.BetLevel)
		betCost := 1
		if betDef != nil {
			betCost = betDef.BetCost
		}
		reward = int(cloneTarget.Multiplier * float64(betCost) * LuckyCloneFishKillMult)
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

	log.Printf("[LuckyClone] player=%s clone_shot side=%s target=%s killed=%v reward=%d",
		p.ID, side, nearestID, killed, reward)

	// 廣播分身子彈命中
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyCloneFish,
		Payload: ws.LuckyCloneFishPayload{
			Event:      "clone_hit",
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			Side:       side,
			TargetID:   nearestID,
			Killed:     killed,
			Reward:     reward,
			X:          nearestX,
			Y:          nearestY,
			DirX:       dirX,
			DirY:       dirY,
		},
	})

	// 若擊破，廣播目標死亡
	if killed {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetKill,
			Payload: ws.TargetKillPayload{
				InstanceID: nearestID,
				DefID:      cloneTarget.DefID,
				KillerID:   p.ID,
				Reward:     reward,
			},
		})
	}
}
