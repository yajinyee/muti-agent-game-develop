// abyss_vortex_handler.go — 深淵漩渦魚系統（DAY-202）
// 業界依據：Ocean King 2「Vortex Fish — sucks all fish of the same species into a whirlpool.
// Catching a Vortex Fish will suck all fish of the same species in the area into a whirlpool.」
// + SteamDB OceanFest 2026「Abyssal Vortex (Depth 3, persistent whirlpool)」
//
// 設計：擊破 T160 後在擊破位置生成「深淵漩渦」（持續 5 秒）：
//   1. 漩渦每 0.5 秒「吸引脈衝」：場上所有目標向漩渦中心移動（位置更新廣播）
//   2. 被吸入漩渦中心（100px 內）的目標：80% 擊破機率，0.70x 倍率
//   3. 漩渦結束後「深淵爆炸」：300px 半徑，60% 擊破機率，0.55x 倍率
//   4. 全服廣播漩渦位置、每次脈衝結果、最終爆炸結果
//
// 設計差異：
//   - 與漩渦魚（直接擊破所有基礎目標）不同，深淵漩渦是「物理吸引」，
//     讓玩家看到「魚被吸向漩渦中心」的動態視覺過程
//   - 與連鎖爆炸魚（靜態爆炸）不同，深淵漩渦是「持續 5 秒的動態場景」，
//     有「吸引→聚集→爆炸」的完整敘事弧
//   - 「吸引脈衝」讓玩家感受到「漩渦在把魚吸過來」的物理感
//   - 最終爆炸是「聚集後的清場」，讓玩家有「等待→爆發」的高潮感
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

// 深淵漩渦魚常數
const (
	AbyssVortexCooldownSec  = 40    // 全服冷卻 40 秒
	AbyssVortexDuration     = 5     // 漩渦持續 5 秒
	AbyssVortexPulseMs      = 500   // 脈衝間隔 500ms
	AbyssVortexPullRadius   = 500.0 // 吸引半徑 500px（場上大部分目標都會被吸引）
	AbyssVortexKillRadius   = 100.0 // 擊破半徑 100px（進入中心才會被擊破）
	AbyssVortexKillChance   = 0.80  // 中心擊破機率 80%
	AbyssVortexKillMult     = 0.70  // 中心擊破倍率 0.70x
	AbyssVortexBlastRadius  = 300.0 // 最終爆炸半徑 300px
	AbyssVortexBlastChance  = 0.60  // 最終爆炸擊破機率 60%
	AbyssVortexBlastMult    = 0.55  // 最終爆炸倍率 0.55x
	AbyssVortexPullSpeed    = 180.0 // 吸引速度 180px/pulse（每次脈衝移動距離）
)

// abyssVortexManager 深淵漩渦魚管理器（全服共享）
type abyssVortexManager struct {
	mu          sync.Mutex
	isActive    bool
	cooldownEnd time.Time
}

func newAbyssVortexManager() *abyssVortexManager {
	return &abyssVortexManager{}
}

// isAbyssVortexFish 判斷是否為深淵漩渦魚（T160）
func isAbyssVortexFish(defID string) bool {
	return defID == "T160"
}

// tryAbyssVortexPull 擊破 T160 後觸發深淵漩渦
func (g *Game) tryAbyssVortexPull(p *player.Player, killX, killY float64) {
	mgr := g.AbyssVortex
	mgr.mu.Lock()

	// 全服冷卻檢查
	if mgr.isActive || time.Now().Before(mgr.cooldownEnd) {
		mgr.mu.Unlock()
		return
	}
	mgr.isActive = true
	mgr.mu.Unlock()

	log.Printf("[AbyssVortex] player=%s triggered vortex at (%.0f, %.0f)", p.ID, killX, killY)

	// 全服廣播：漩渦開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAbyssVortex,
		Payload: ws.AbyssVortexPayload{
			Event:      "vortex_start",
			KillerName: p.DisplayName,
			VortexX:    killX,
			VortexY:    killY,
			Duration:   AbyssVortexDuration,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌀 %s 觸發深淵漩渦！所有目標正在被吸入！", p.DisplayName),
	})
	g.broadcastAnnouncement(ann)

	// 執行漩渦主循環
	go g.runAbyssVortexLoop(p, killX, killY)
}

// runAbyssVortexLoop 漩渦主循環（goroutine）
func (g *Game) runAbyssVortexLoop(p *player.Player, vortexX, vortexY float64) {
	pulseCount := AbyssVortexDuration * 1000 / AbyssVortexPulseMs // 10 次脈衝
	totalKills := 0
	totalReward := 0

	for pulse := 0; pulse < pulseCount; pulse++ {
		time.Sleep(AbyssVortexPulseMs * time.Millisecond)

		pulseKills := 0
		pulseReward := 0

		// 吸引脈衝：移動所有目標向漩渦中心
		g.mu.Lock()
		type targetMove struct {
			id   string
			newX float64
			newY float64
			dist float64
			mult float64
		}
		var moves []targetMove

		for _, t := range g.Targets {
			if t.HP <= 0 || t.DefID == "B001" || isGhostFishClone(t.DefID) {
				continue
			}
			dx := vortexX - t.X
			dy := vortexY - t.Y
			dist := math.Sqrt(dx*dx + dy*dy)

			// 只吸引吸引半徑內的目標
			if dist > AbyssVortexPullRadius {
				continue
			}
			if dist < 1.0 {
				dist = 1.0
			}

			// 計算新位置（向漩渦中心移動）
			moveStep := AbyssVortexPullSpeed
			if dist < moveStep {
				moveStep = dist // 不超過漩渦中心
			}
			ratio := moveStep / dist
			newX := t.X + dx*ratio
			newY := t.Y + dy*ratio

			moves = append(moves, targetMove{
				id:   t.InstanceID,
				newX: newX,
				newY: newY,
				dist: dist,
				mult: t.Multiplier,
			})
		}

		// 更新目標位置並檢查是否進入擊破範圍
		type killedTarget struct {
			id     string
			mult   float64
			reward int
		}
		var killed []killedTarget

		for _, mv := range moves {
			t, ok := g.Targets[mv.id]
			if !ok || t.HP <= 0 {
				continue
			}
			// 更新位置
			t.X = mv.newX
			t.Y = mv.newY

			// 檢查是否進入擊破範圍（100px 內）
			newDist := math.Sqrt((vortexX-mv.newX)*(vortexX-mv.newX) + (vortexY-mv.newY)*(vortexY-mv.newY))
			if newDist <= AbyssVortexKillRadius {
				if rand.Float64() < AbyssVortexKillChance {
					r := int(mv.mult * float64(p.BetLevel) * AbyssVortexKillMult)
					if r < 1 {
						r = 1
					}
					delete(g.Targets, mv.id)
					killed = append(killed, killedTarget{id: mv.id, mult: mv.mult, reward: r})
					pulseKills++
					pulseReward += r
					totalKills++
					totalReward += r
					// 給觸發者獎勵
					if pp, ok2 := g.Players[p.ID]; ok2 {
						pp.Coins += r
					}
				}
			}
		}
		g.mu.Unlock()

		// 廣播脈衝結果（包含目標移動資訊）
		pulseNum := pulse + 1
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgAbyssVortex,
			Payload: ws.AbyssVortexPayload{
				Event:       "vortex_pulse",
				VortexX:     vortexX,
				VortexY:     vortexY,
				PulseNum:    pulseNum,
				PulseKills:  pulseKills,
				PulseReward: pulseReward,
				TotalKills:  totalKills,
			},
		})
	}

	// 漩渦結束：深淵爆炸
	time.Sleep(200 * time.Millisecond) // 短暫停頓，讓玩家感受到「漩渦消失→爆炸」

	blastKills := 0
	blastReward := 0

	g.mu.Lock()
	type blastTarget struct {
		id   string
		mult float64
	}
	var inBlast []blastTarget
	for _, t := range g.Targets {
		if t.HP <= 0 || t.DefID == "B001" || isGhostFishClone(t.DefID) {
			continue
		}
		dx := vortexX - t.X
		dy := vortexY - t.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= AbyssVortexBlastRadius {
			inBlast = append(inBlast, blastTarget{id: t.InstanceID, mult: t.Multiplier})
		}
	}

	for _, bt := range inBlast {
		if rand.Float64() < AbyssVortexBlastChance {
			t, ok := g.Targets[bt.id]
			if ok && t.HP > 0 {
				r := int(bt.mult * float64(p.BetLevel) * AbyssVortexBlastMult)
				if r < 1 {
					r = 1
				}
				delete(g.Targets, bt.id)
				blastKills++
				blastReward += r
				totalKills++
				totalReward += r
				if pp, ok2 := g.Players[p.ID]; ok2 {
					pp.Coins += r
				}
			}
		}
	}
	g.mu.Unlock()

	// 廣播最終爆炸
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAbyssVortex,
		Payload: ws.AbyssVortexPayload{
			Event:       "vortex_blast",
			VortexX:     vortexX,
			VortexY:     vortexY,
			BlastKills:  blastKills,
			BlastReward: blastReward,
		},
	})

	// 廣播最終結算
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAbyssVortex,
		Payload: ws.AbyssVortexPayload{
			Event:       "vortex_result",
			KillerName:  p.DisplayName,
			TotalKills:  totalKills,
			TotalReward: totalReward,
		},
	})

	// 全服公告（≥5 個擊破才公告）
	if totalKills >= 5 {
		ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("🌀💥 深淵漩渦！吸入並擊破 %d 個目標！獎勵 %d 金幣！", totalKills, totalReward),
		})
		g.broadcastAnnouncement(ann)
	}

	// 重置管理器
	mgr := g.AbyssVortex
	mgr.mu.Lock()
	mgr.isActive = false
	mgr.cooldownEnd = time.Now().Add(AbyssVortexCooldownSec * time.Second)
	mgr.mu.Unlock()

	log.Printf("[AbyssVortex] complete: pulseKills=%d blastKills=%d totalReward=%d",
		totalKills-blastKills, blastKills, totalReward)
}
