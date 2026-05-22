// drill_lobster_handler.go — 鑽頭龍蝦穿透爆炸系統（DAY-195）
// 業界依據：Royal Fishing JILI「Drill Bit Lobster (80X) — fires a penetrating drill that passes
// through multiple fish before self-detonating, capturing everything in the explosion radius.
// Mechanical marvel with penetrating drill projectiles.」
// 設計：擊破 T153 後發射「穿透鑽頭」：
//   1. 鑽頭從擊破位置出發，沿隨機方向穿透（每 150ms 前進一步）
//   2. 穿透路徑上的目標：80% 擊破機率（0.70x 倍率）
//   3. 穿透最多 5 個目標後，在終點「自爆」（300px 半徑爆炸）
//   4. 爆炸範圍內所有目標：75% 擊破機率（0.65x 倍率）
//   5. 全服廣播每步穿透和最終爆炸
// 設計差異：
//   - 與連鎖爆炸魚（靜態爆炸，BFS 擴散）不同，鑽頭龍蝦是「動態移動的穿透彈」，
//     讓玩家看到「鑽頭在場上移動穿透」的動態視覺
//   - 與隕石魚（從天而降，隨機目標）不同，鑽頭龍蝦是「從擊破點出發，沿路徑穿透」，
//     讓玩家感受到「鑽頭是從我打的那條魚身上射出來的」的因果感
//   - 穿透 + 爆炸雙段式設計：穿透製造「一路收割」的爽感，爆炸製造「最後一擊清場」的高潮感
//   - 鑽頭方向隨機（8方向），讓每次觸發都有不同的視覺路徑
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

// 鑽頭龍蝦常數
const (
	DrillLobsterCooldownSec    = 25    // 全服冷卻 25 秒
	DrillLobsterPenetrateMax   = 5     // 最多穿透 5 個目標
	DrillLobsterPenetrateChance = 0.80 // 穿透擊破機率 80%
	DrillLobsterPenetrateMult  = 0.70  // 穿透獎勵倍率 0.70x
	DrillLobsterExplodeRadius  = 300.0 // 爆炸半徑 300px
	DrillLobsterExplodeChance  = 0.75  // 爆炸擊破機率 75%
	DrillLobsterExplodeMult    = 0.65  // 爆炸獎勵倍率 0.65x
	DrillLobsterStepPx         = 120.0 // 每步移動距離 120px
	DrillLobsterStepMs         = 150   // 每步間隔 150ms
)

// drillLobsterManager 鑽頭龍蝦管理器（全服冷卻）
type drillLobsterManager struct {
	mu          sync.Mutex
	isActive    bool
	cooldownEnd time.Time
}

func newDrillLobsterManager() *drillLobsterManager {
	return &drillLobsterManager{}
}

// isDrillBitLobster 判斷是否為鑽頭龍蝦（T153，DAY-195）
func isDrillBitLobster(defID string) bool {
	return defID == "T153"
}

// isDrillLobster 判斷是否為鑽頭龍蝦（T106，DAY-142）
func isDrillLobster(defID string) bool {
	return defID == "T106"
}

// tryDrillLobsterChain T106 鑽頭龍蝦連帶效果（DAY-142）
// 擊破後觸發穿透鑽頭，沿水平方向穿透所有目標，到達邊緣後爆炸
func (g *Game) tryDrillLobsterChain(p *player.Player, instanceID string, startX, startY float64) {
	// 使用 T153 的穿透爆炸邏輯（相同機制，T106 是舊版本）
	// T106 方向固定為水平向右（業界原版設計）
	g.mu.RLock()
	type drillT106Info struct {
		instanceID string
		multiplier float64
		x, y       float64
	}
	var targets []drillT106Info
	for _, t := range g.Targets {
		if t.InstanceID == instanceID {
			continue
		}
		// 水平方向：Y 座標相近（±60px），X 在右側
		if t.X > startX && math.Abs(t.Y-startY) < 60 {
			targets = append(targets, drillT106Info{
				instanceID: t.InstanceID,
				multiplier: t.Multiplier,
				x:          t.X,
				y:          t.Y,
			})
		}
	}
	g.mu.RUnlock()

	// 按 X 座標排序（從近到遠）
	for i := 0; i < len(targets)-1; i++ {
		for j := i + 1; j < len(targets); j++ {
			if targets[j].x < targets[i].x {
				targets[i], targets[j] = targets[j], targets[i]
			}
		}
	}

	totalKills := 0
	totalReward := 0

	// 穿透最多 5 個目標
	for i, t := range targets {
		if i >= 5 {
			break
		}
		time.Sleep(150 * time.Millisecond)
		if rand.Float64() < 0.80 {
			r := int(t.multiplier * float64(p.BetLevel) * 0.70)
			g.mu.Lock()
			if _, ok := g.Targets[t.instanceID]; ok {
				delete(g.Targets, t.instanceID)
				p.Coins += r
				totalKills++
				totalReward += r
			}
			g.mu.Unlock()
		}
	}

	// 廣播結果
	if totalKills > 0 {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgDrillLobster,
			Payload: ws.DrillLobsterPayload{
				Phase:       "drill_result",
				TriggerID:   instanceID,
				KillerID:    p.ID,
				KillerName:  p.DisplayName,
				TotalKills:  totalKills,
				TotalReward: totalReward,
			},
		})
	}

	log.Printf("[DrillLobster T106] player=%s kills=%d reward=%d", p.ID, totalKills, totalReward)
}

// tryDrillLobsterPenetrate 擊破 T153 後觸發穿透爆炸
func (g *Game) tryDrillLobsterPenetrate(p *player.Player, instanceID string, startX, startY float64) {
	mgr := g.DrillLobster
	mgr.mu.Lock()

	// 全服冷卻檢查
	if mgr.isActive || time.Now().Before(mgr.cooldownEnd) {
		mgr.mu.Unlock()
		return
	}
	mgr.isActive = true
	mgr.mu.Unlock()

	defer func() {
		mgr.mu.Lock()
		mgr.isActive = false
		mgr.cooldownEnd = time.Now().Add(DrillLobsterCooldownSec * time.Second)
		mgr.mu.Unlock()
	}()

	// 隨機選擇 8 方向之一
	directions := [][2]float64{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{0.707, 0.707}, {-0.707, 0.707}, {0.707, -0.707}, {-0.707, -0.707},
	}
	dir := directions[rand.Intn(len(directions))]
	dirX, dirY := dir[0], dir[1]

	// local struct for target info
	type drillTargetInfo struct {
		instanceID string
		defID      string
		x, y       float64
		multiplier float64
		hp         int
	}

	// 廣播鑽頭出發（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgDrillLobster,
		Payload: ws.DrillLobsterPayload{
			Phase:      "drill_start",
			TriggerID:  instanceID,
			KillerID:   p.ID,
			KillerName: p.DisplayName,
			StartX:     startX,
			StartY:     startY,
			DirX:       dirX,
			DirY:       dirY,
		},
	})

	log.Printf("[DrillLobster] player=%s triggered drill at (%.0f,%.0f) dir=(%.2f,%.2f)",
		p.ID, startX, startY, dirX, dirY)

	// 穿透階段：沿方向移動，每步檢查目標
	curX, curY := startX, startY
	penetrateCount := 0
	totalKills := 0
	totalReward := 0
	hitTargetIDs := make(map[string]bool) // 避免重複命中

	for step := 1; step <= DrillLobsterPenetrateMax; step++ {
		time.Sleep(DrillLobsterStepMs * time.Millisecond)

		curX += dirX * DrillLobsterStepPx
		curY += dirY * DrillLobsterStepPx

		// 找最近的目標（50px 範圍內）
		g.mu.RLock()
		var nearestTarget *drillTargetInfo
		nearestDist := 50.0
		for _, t := range g.Targets {
			if hitTargetIDs[t.InstanceID] {
				continue
			}
			dx := t.X - curX
			dy := t.Y - curY
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < nearestDist {
				nearestDist = dist
				nearestTarget = &drillTargetInfo{
					instanceID: t.InstanceID,
					defID:      t.DefID,
					x:          t.X,
					y:          t.Y,
					multiplier: t.Multiplier,
					hp:         t.HP,
				}
			}
		}
		g.mu.RUnlock()

		var killID, killName string
		var killMult float64
		var stepReward int
		isKill := false

		if nearestTarget != nil {
			hitTargetIDs[nearestTarget.instanceID] = true
			penetrateCount++

			// 80% 擊破機率
			if rand.Float64() < DrillLobsterPenetrateChance {
				isKill = true
				stepReward = int(nearestTarget.multiplier * float64(p.BetLevel) * DrillLobsterPenetrateMult)
				totalKills++
				totalReward += stepReward
				killID = nearestTarget.instanceID
				killName = nearestTarget.defID
				killMult = nearestTarget.multiplier

				// 擊破目標
				g.mu.Lock()
				if _, ok := g.Targets[nearestTarget.instanceID]; ok {
					delete(g.Targets, nearestTarget.instanceID)
				}
				g.mu.Unlock()

				// 給予獎勵
				g.mu.Lock()
				p.Coins += stepReward
				g.mu.Unlock()
			}
		}

		// 廣播每步穿透（全服）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgDrillLobster,
			Payload: ws.DrillLobsterPayload{
				Phase:      fmt.Sprintf("drill_%d", step),
				TriggerID:  instanceID,
				StepIndex:  step,
				CurX:       curX,
				CurY:       curY,
				IsKill:     isKill,
				KilledID:   killID,
				KilledName: killName,
				KilledMult: killMult,
				StepReward: stepReward,
				TotalKills: totalKills,
			},
		})

		_ = killName
	}

	// 爆炸階段：在終點爆炸
	time.Sleep(200 * time.Millisecond)
	explodeKills, explodeReward := g.doDrillLobsterExplosion(p, curX, curY, hitTargetIDs)
	totalKills += explodeKills
	totalReward += explodeReward

	// 廣播爆炸（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgDrillLobster,
		Payload: ws.DrillLobsterPayload{
			Phase:         "drill_explode",
			TriggerID:     instanceID,
			CurX:          curX,
			CurY:          curY,
			ExplodeRadius: DrillLobsterExplodeRadius,
			ExplodeKills:  explodeKills,
			ExplodeReward: explodeReward,
		},
	})

	// 廣播結果（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgDrillLobster,
		Payload: ws.DrillLobsterPayload{
			Phase:          "drill_result",
			TriggerID:      instanceID,
			KillerID:       p.ID,
			KillerName:     p.DisplayName,
			PenetrateCount: penetrateCount,
			TotalKills:     totalKills,
			TotalReward:    totalReward,
		},
	})

	// 全服公告（≥4 個擊破）
	if totalKills >= 4 {
		ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("🦞 %s 鑽頭龍蝦穿透爆炸！擊破 %d 個目標！獎勵 %d 金幣！",
				p.DisplayName, totalKills, totalReward),
		})
		g.broadcastAnnouncement(ann)
	}

	log.Printf("[DrillLobster] done: penetrate=%d totalKills=%d totalReward=%d",
		penetrateCount, totalKills, totalReward)
}

// doDrillLobsterExplosion 在終點爆炸，擊破範圍內目標
func (g *Game) doDrillLobsterExplosion(p *player.Player, cx, cy float64, skipIDs map[string]bool) (int, int) {
	type explodeTarget struct {
		instanceID string
		multiplier float64
	}

	g.mu.RLock()
	var targets []explodeTarget
	for _, t := range g.Targets {
		if skipIDs[t.InstanceID] {
			continue
		}
		dx := t.X - cx
		dy := t.Y - cy
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= DrillLobsterExplodeRadius {
			targets = append(targets, explodeTarget{
				instanceID: t.InstanceID,
				multiplier: t.Multiplier,
			})
		}
	}
	g.mu.RUnlock()

	kills := 0
	reward := 0
	for _, t := range targets {
		if rand.Float64() < DrillLobsterExplodeChance {
			r := int(t.multiplier * float64(p.BetLevel) * DrillLobsterExplodeMult)
			g.mu.Lock()
			if _, ok := g.Targets[t.instanceID]; ok {
				delete(g.Targets, t.instanceID)
				p.Coins += r
				kills++
				reward += r
			}
			g.mu.Unlock()
		}
	}
	return kills, reward
}
