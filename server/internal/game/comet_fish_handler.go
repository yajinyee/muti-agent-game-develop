// comet_fish_handler.go — 彗星魚連鎖爆炸系統（DAY-206）
// 業界依據：Ocean King 3 Plus「Comet Fish — streaks across the screen leaving a trail of explosions,
// each explosion has a chance to capture fish in its radius. The comet leaves destruction in its wake.」
//
// 設計：T164 彗星魚出現後，沿隨機弧線軌跡飛越全場（1.5秒），
//   1. 沿途每 200ms 留下一個「彗星爆炸點」（共 7 個）
//   2. 每個爆炸點 200px 半徑，70% 擊破機率，0.65x 倍率
//   3. 最終在終點「超新星爆炸」（400px 半徑，80% 擊破機率，0.75x 倍率）
//   4. 玩家擊破彗星魚本身可以「提前引爆」（立即觸發超新星）
//   5. 全服冷卻 40 秒（防止頻繁觸發）
//
// 設計差異（與其他爆炸系統的區別）：
//   - 連環炸彈蟹（T159）：靜態多點爆炸，玩家擊破後觸發
//   - 鑽頭龍蝦（T153）：穿透移動，沿路徑擊破
//   - 彗星魚（T164）：「動態軌跡」— 玩家看到彗星飛越全場，沿途爆炸，最後超新星
//   - 「提前引爆」機制讓玩家有「要不要現在打」的策略決策
//   - 超新星爆炸範圍是普通爆炸的 2 倍，製造「最後一擊最爽」的高潮感
package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// 彗星魚常數（DAY-206）
const (
	CometFishGlobalCooldownSec   = 40    // 全服冷卻 40 秒
	CometFishTrailCount          = 7     // 軌跡爆炸點數量
	CometFishTrailInterval       = 200   // 軌跡爆炸間隔 200ms
	CometFishTrailRadius         = 200.0 // 軌跡爆炸半徑 200px
	CometFishTrailKillChance     = 0.70  // 軌跡爆炸擊破機率 70%
	CometFishTrailMult           = 0.65  // 軌跡爆炸獎勵倍率 0.65x
	CometFishSupernovaRadius     = 400.0 // 超新星爆炸半徑 400px
	CometFishSupernovaKillChance = 0.80  // 超新星擊破機率 80%
	CometFishSupernovaMult       = 0.75  // 超新星獎勵倍率 0.75x
)

// cometFishManager 彗星魚管理器（全服冷卻）
type cometFishManager struct {
	mu         sync.Mutex
	isActive   bool      // 是否正在飛行
	cooldownAt time.Time // 全服冷卻結束時間
}

func newCometFishManager() *cometFishManager {
	return &cometFishManager{}
}

// isCometFish 判斷是否為彗星魚（T164，DAY-206）
func isCometFish(defID string) bool {
	return defID == "T164"
}

// notifyCometFishSpawn 彗星魚生成時觸發軌跡飛行（由 spawnTarget 呼叫）
func (g *Game) notifyCometFishSpawn(t *target.Target) {
	mgr := g.CometFish
	mgr.mu.Lock()

	// 全服冷卻檢查
	if time.Now().Before(mgr.cooldownAt) {
		mgr.mu.Unlock()
		return
	}
	if mgr.isActive {
		mgr.mu.Unlock()
		return
	}
	mgr.isActive = true
	mgr.cooldownAt = time.Now().Add(CometFishGlobalCooldownSec * time.Second)
	mgr.mu.Unlock()

	// 生成彗星軌跡點（弧線路徑）
	trailPoints := generateCometTrail(t.X, t.Y)

	// 廣播彗星出現
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCometFish,
		Payload: ws.CometFishPayload{
			Event:       "comet_appear",
			InstanceID:  t.InstanceID,
			StartX:      t.X,
			StartY:      t.Y,
			TrailPoints: trailPoints,
			TrailCount:  len(trailPoints),
		},
	})

	log.Printf("[CometFish] comet appeared at (%.0f,%.0f), trail=%d points", t.X, t.Y, len(trailPoints))

	// 啟動軌跡爆炸 goroutine
	go g.runCometFishTrail(t.InstanceID, trailPoints)
}

// notifyCometFishKill 玩家擊破彗星魚時「提前引爆」超新星
func (g *Game) notifyCometFishKill(p *player.Player, t *target.Target) {
	mgr := g.CometFish
	mgr.mu.Lock()
	if !mgr.isActive {
		mgr.mu.Unlock()
		return
	}
	// 標記為非活躍（防止 runCometFishTrail 繼續執行超新星）
	mgr.isActive = false
	mgr.mu.Unlock()

	log.Printf("[CometFish] player=%s triggered early supernova at (%.0f,%.0f)", p.ID, t.X, t.Y)

	// 廣播提前引爆
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCometFish,
		Payload: ws.CometFishPayload{
			Event:      "early_supernova",
			InstanceID: t.InstanceID,
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			X:          t.X,
			Y:          t.Y,
		},
	})

	// 立即觸發超新星
	go g.doCometFishSupernova(p, t.X, t.Y, true)
}

// runCometFishTrail 執行彗星軌跡爆炸（goroutine）
func (g *Game) runCometFishTrail(instanceID string, trailPoints []ws.CometPoint) {
	totalTrailReward := 0
	totalTrailKills := 0

	for i, pt := range trailPoints {
		time.Sleep(CometFishTrailInterval * time.Millisecond)

		// 檢查是否已被提前引爆
		g.CometFish.mu.Lock()
		active := g.CometFish.isActive
		g.CometFish.mu.Unlock()
		if !active {
			log.Printf("[CometFish] trail interrupted at point %d (early supernova)", i)
			return
		}

		// 執行軌跡爆炸
		kills, reward := g.doCometFishTrailBlast(pt.X, pt.Y)
		totalTrailKills += kills
		totalTrailReward += reward

		// 廣播軌跡爆炸結果
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgCometFish,
			Payload: ws.CometFishPayload{
				Event:      "trail_blast",
				InstanceID: instanceID,
				X:          pt.X,
				Y:          pt.Y,
				BlastIndex: i + 1,
				KillCount:  kills,
				Reward:     reward,
			},
		})
	}

	// 所有軌跡爆炸完成，觸發超新星
	g.CometFish.mu.Lock()
	active := g.CometFish.isActive
	if active {
		g.CometFish.isActive = false
	}
	g.CometFish.mu.Unlock()

	if !active {
		return // 已被提前引爆
	}

	// 超新星在最後一個軌跡點
	lastPt := trailPoints[len(trailPoints)-1]
	g.doCometFishSupernova(nil, lastPt.X, lastPt.Y, false)

	log.Printf("[CometFish] trail complete: kills=%d reward=%d", totalTrailKills, totalTrailReward)
}

// doCometFishTrailBlast 執行單個軌跡爆炸點
func (g *Game) doCometFishTrailBlast(cx, cy float64) (kills, reward int) {
	type candidate struct {
		instanceID string
		multiplier float64
	}

	g.mu.RLock()
	var candidates []candidate
	for _, t := range g.Targets {
		if !t.IsAlive || isCometFish(t.DefID) {
			continue
		}
		dx := t.X - cx
		dy := t.Y - cy
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= CometFishTrailRadius {
			candidates = append(candidates, candidate{t.InstanceID, t.Multiplier})
		}
	}
	betLevel := 1
	for _, p := range g.Players {
		betLevel = p.BetLevel
		break
	}
	g.mu.RUnlock()

	for _, c := range candidates {
		if rand.Float64() >= CometFishTrailKillChance {
			continue
		}
		rewardAmt := int(float64(betLevel) * c.multiplier * CometFishTrailMult)
		if rewardAmt < 1 {
			rewardAmt = 1
		}
		g.mu.Lock()
		if tgt, ok := g.Targets[c.instanceID]; ok && tgt.IsAlive {
			tgt.IsAlive = false
			tgt.HP = 0
			delete(g.Targets, c.instanceID)
			kills++
			reward += rewardAmt
			g.distributeRewardToAll(rewardAmt)
		}
		g.mu.Unlock()
	}
	return
}

// doCometFishSupernova 執行超新星爆炸
func (g *Game) doCometFishSupernova(triggerer *player.Player, cx, cy float64, isEarly bool) {
	type candidate struct {
		instanceID string
		multiplier float64
	}

	g.mu.RLock()
	var candidates []candidate
	for _, t := range g.Targets {
		if !t.IsAlive || isCometFish(t.DefID) {
			continue
		}
		dx := t.X - cx
		dy := t.Y - cy
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= CometFishSupernovaRadius {
			candidates = append(candidates, candidate{t.InstanceID, t.Multiplier})
		}
	}
	betLevel := 1
	for _, p := range g.Players {
		betLevel = p.BetLevel
		break
	}
	g.mu.RUnlock()

	kills := 0
	reward := 0

	for _, c := range candidates {
		if rand.Float64() >= CometFishSupernovaKillChance {
			continue
		}
		rewardAmt := int(float64(betLevel) * c.multiplier * CometFishSupernovaMult)
		if rewardAmt < 1 {
			rewardAmt = 1
		}
		g.mu.Lock()
		if tgt, ok := g.Targets[c.instanceID]; ok && tgt.IsAlive {
			tgt.IsAlive = false
			tgt.HP = 0
			delete(g.Targets, c.instanceID)
			kills++
			reward += rewardAmt
			g.distributeRewardToAll(rewardAmt)
		}
		g.mu.Unlock()
	}

	// 廣播超新星結果
	playerID := ""
	playerName := ""
	if triggerer != nil {
		playerID = triggerer.ID
		playerName = triggerer.DisplayName
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCometFish,
		Payload: ws.CometFishPayload{
			Event:      "supernova",
			X:          cx,
			Y:          cy,
			PlayerID:   playerID,
			PlayerName: playerName,
			KillCount:  kills,
			Reward:     reward,
			IsEarly:    isEarly,
		},
	})

	log.Printf("[CometFish] supernova at (%.0f,%.0f): kills=%d reward=%d early=%v", cx, cy, kills, reward, isEarly)

	// 全服公告（≥5 個擊破）
	if kills >= 5 {
		color := "#FFD700"
		if kills >= 10 {
			color = "#FF8C00"
		}
		if kills >= 15 {
			color = "#FF4500"
		}
		msg := fmt.Sprintf("☄️ 彗星超新星爆炸！擊破 %d 個目標！獎勵 %d 金幣！", kills, reward)
		if isEarly && playerName != "" {
			msg = fmt.Sprintf("☄️ %s 提前引爆彗星！超新星擊破 %d 個目標！", playerName, kills)
		}
		ann := g.Announce.Create(announce.EventMegaWin, playerName, reward, map[string]string{
			"message": msg,
			"color":   color,
		})
		g.broadcastAnnouncement(ann)
	}
}

// generateCometTrail 生成彗星弧線軌跡點（二次貝茲曲線）
// 從起始位置沿弧線飛越全場，生成 CometFishTrailCount 個爆炸點
func generateCometTrail(startX, startY float64) []ws.CometPoint {
	const (
		screenW = 1280.0
		screenH = 720.0
	)

	// 隨機終點（對角線方向，確保飛越全場）
	endX := screenW - startX + rand.Float64()*200 - 100
	if endX < 100 {
		endX = 100
	}
	if endX > screenW-100 {
		endX = screenW - 100
	}
	endY := rand.Float64() * screenH

	// 弧線控制點（向上弧，製造「彗星飛越」感）
	ctrlX := (startX+endX)/2 + (rand.Float64()-0.5)*400
	ctrlY := (startY+endY)/2 - 200 - rand.Float64()*200

	points := make([]ws.CometPoint, CometFishTrailCount)
	for i := 0; i < CometFishTrailCount; i++ {
		t := float64(i+1) / float64(CometFishTrailCount)
		// 二次貝茲曲線公式
		x := (1-t)*(1-t)*startX + 2*(1-t)*t*ctrlX + t*t*endX
		y := (1-t)*(1-t)*startY + 2*(1-t)*t*ctrlY + t*t*endY
		points[i] = ws.CometPoint{X: x, Y: y}
	}
	return points
}

// distributeRewardToAll 將獎勵分配給所有玩家（全服共享）
func (g *Game) distributeRewardToAll(amount int) {
	if amount <= 0 {
		return
	}
	playerCount := len(g.Players)
	if playerCount == 0 {
		return
	}
	perPlayer := amount / playerCount
	if perPlayer < 1 {
		perPlayer = 1
	}
	for _, p := range g.Players {
		p.AddCoins(perPlayer)
	}
}
