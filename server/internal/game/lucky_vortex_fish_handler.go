// lucky_vortex_fish_handler.go — 幸運漩渦魚系統（DAY-234）
// 業界原創「漩渦旋轉+高倍率加成」機制
//
// 設計：擊破 T192 後觸發「漩渦模式」（8 秒）：
//   - 場景中央半徑 300px 內的所有目標每 2 秒被「漩渦旋轉」（繞中心旋轉 45 度）
//   - 漩渦模式期間擊破任何目標獲得 ×2.2 倍率加成（乘法）
//   - 8 秒後「漩渦爆發」：漩渦範圍內所有目標 70% 擊破機率（0.75x 倍率，全服共享）
//   - 個人冷卻 20 秒；全服冷卻 30 秒
//
// 設計差異：
//   - 與磁力魚（DAY-232，向中央聚集）不同，漩渦魚是「繞圈旋轉」，
//     讓玩家有「要跟著旋轉方向瞄準」的動態策略感
//   - 「漩渦旋轉」讓目標繞中心移動，製造「旋轉魚群」的視覺爽感
//   - ×2.2 倍率加成（全場有效，不限漩渦範圍）讓玩家有「趕快在 8 秒內多打」的緊迫感
//   - 「漩渦爆發」讓玩家有「等待→爆發」的高潮設計
//   - 全服廣播漩渦讓所有玩家都往中央打，製造「全服聚焦」的社交感
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
	LuckyVortexPersonalCD     = 20 * time.Second          // 個人冷卻
	LuckyVortexGlobalCD       = 30 * time.Second          // 全服冷卻
	LuckyVortexDuration       = 8 * time.Second           // 漩渦模式持續時間
	LuckyVortexRotateInterval = 2000 * time.Millisecond   // 每次旋轉間隔
	LuckyVortexKillBoost      = 2.2                       // 漩渦模式期間擊破倍率加成（乘法）
	LuckyVortexBlastChance    = 0.70                      // 漩渦爆發擊破機率
	LuckyVortexBlastMult      = 0.75                      // 漩渦爆發倍率
	LuckyVortexRadius         = 300.0                     // 漩渦範圍（半徑 px）
	LuckyVortexRotateDeg      = 45.0                      // 每次旋轉角度（度）
	LuckyVortexCenterX        = 500.0                     // 場景中央 X
	LuckyVortexCenterY        = 300.0                     // 場景中央 Y
)

// luckyVortexFishManager 幸運漩渦魚管理器
type luckyVortexFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 漩渦模式狀態
	active      bool
	activeUntil time.Time
	instanceID  string
}

func newLuckyVortexFishManager() *luckyVortexFishManager {
	return &luckyVortexFishManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyVortexFish 判斷是否為幸運漩渦魚
func isLuckyVortexFish(defID string) bool {
	return defID == "T192"
}

// isVortexActive 判斷漩渦模式是否啟動中（供 handleKill 使用）
func (g *Game) isVortexActive() bool {
	mgr := g.LuckyVortexFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if !mgr.active {
		return false
	}
	if time.Now().After(mgr.activeUntil) {
		mgr.active = false
		return false
	}
	return true
}

// getLuckyVortexBoost 取得漩渦模式倍率加成（供 handleKill 使用）
func (g *Game) getLuckyVortexBoost() float64 {
	if g.isVortexActive() {
		return LuckyVortexKillBoost
	}
	return 1.0
}

// tryLuckyVortexFish 擊破 T192 後觸發漩渦模式（供 handleKill 使用）
func (g *Game) tryLuckyVortexFish(p *player.Player) {
	mgr := g.LuckyVortexFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.personalCooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}
	// 全服冷卻檢查
	if time.Now().Before(mgr.globalCooldownUntil) {
		mgr.mu.Unlock()
		return
	}
	// 已有漩渦模式啟動中
	if mgr.active && time.Now().Before(mgr.activeUntil) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyVortexPersonalCD)
	mgr.globalCooldownUntil = time.Now().Add(LuckyVortexGlobalCD)

	// 啟動漩渦模式
	mgr.active = true
	mgr.activeUntil = time.Now().Add(LuckyVortexDuration)
	instanceID := fmt.Sprintf("vortex_%d", time.Now().UnixNano())
	mgr.instanceID = instanceID
	mgr.mu.Unlock()

	log.Printf("[LuckyVortex] player=%s activated vortex mode for %v", p.ID, LuckyVortexDuration)

	// 全服廣播：漩渦模式開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyVortexFish,
		Payload: ws.LuckyVortexFishPayload{
			Event:       "vortex_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyVortexDuration.Seconds()),
			KillBoost:   LuckyVortexKillBoost,
			CenterX:     LuckyVortexCenterX,
			CenterY:     LuckyVortexCenterY,
			Radius:      LuckyVortexRadius,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyVortexFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🌀 %s 觸發漩渦！目標開始旋轉，×%.1f 倍率加成！",
			p.DisplayName, LuckyVortexKillBoost),
		"color": "#16A085",
	})
	g.broadcastAnnouncement(ann)

	// 啟動漩渦旋轉 goroutine
	go g.runLuckyVortex(p, instanceID)
}

// runLuckyVortex 漩渦主循環（goroutine）
func (g *Game) runLuckyVortex(p *player.Player, instanceID string) {
	ticker := time.NewTicker(LuckyVortexRotateInterval)
	defer ticker.Stop()

	endTimer := time.NewTimer(LuckyVortexDuration)
	defer endTimer.Stop()

	rotateCount := 0

	for {
		select {
		case <-ticker.C:
			// 確認 instanceID 仍有效
			g.LuckyVortexFish.mu.Lock()
			if g.LuckyVortexFish.instanceID != instanceID {
				g.LuckyVortexFish.mu.Unlock()
				return
			}
			g.LuckyVortexFish.mu.Unlock()

			rotateCount++
			rotatedCount := g.doVortexRotate(instanceID, rotateCount)
			log.Printf("[LuckyVortex] rotate#%d: rotated %d targets", rotateCount, rotatedCount)

		case <-endTimer.C:
			// 漩渦模式結束，觸發漩渦爆發
			g.LuckyVortexFish.mu.Lock()
			if g.LuckyVortexFish.instanceID != instanceID {
				g.LuckyVortexFish.mu.Unlock()
				return
			}
			g.LuckyVortexFish.active = false
			g.LuckyVortexFish.mu.Unlock()

			log.Printf("[LuckyVortex] vortex ended, triggering vortex blast")
			g.doVortexBlast(p, instanceID)
			return
		}
	}
}

// doVortexRotate 執行漩渦旋轉（漩渦範圍內目標繞中心旋轉 45 度）
func (g *Game) doVortexRotate(instanceID string, rotateNum int) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	rotatedCount := 0
	type rotatedTarget struct {
		ID string
		X  float64
		Y  float64
	}
	var rotatedTargets []rotatedTarget

	// 旋轉角度（弧度）
	angleRad := LuckyVortexRotateDeg * math.Pi / 180.0

	for id, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		// 計算到中央的距離
		dx := t.X - LuckyVortexCenterX
		dy := t.Y - LuckyVortexCenterY
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist > LuckyVortexRadius || dist < 10.0 {
			// 不在漩渦範圍內，或已在中央
			continue
		}

		// 繞中心旋轉 45 度（順時針）
		cosA := math.Cos(angleRad)
		sinA := math.Sin(angleRad)
		newDX := dx*cosA - dy*sinA
		newDY := dx*sinA + dy*cosA

		t.X = LuckyVortexCenterX + newDX
		t.Y = LuckyVortexCenterY + newDY

		// 邊界限制
		if t.X < 50 {
			t.X = 50
		}
		if t.X > 950 {
			t.X = 950
		}
		if t.Y < 50 {
			t.Y = 50
		}
		if t.Y > 550 {
			t.Y = 550
		}

		rotatedTargets = append(rotatedTargets, rotatedTarget{ID: id, X: t.X, Y: t.Y})
		rotatedCount++
	}

	if rotatedCount == 0 {
		return 0
	}

	// 廣播漩渦旋轉（含旋轉後的位置）
	type targetPos struct {
		ID string  `json:"id"`
		X  float64 `json:"x"`
		Y  float64 `json:"y"`
	}
	positions := make([]targetPos, 0, len(rotatedTargets))
	for _, rt := range rotatedTargets {
		positions = append(positions, targetPos{ID: rt.ID, X: rt.X, Y: rt.Y})
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyVortexFish,
		Payload: ws.LuckyVortexFishPayload{
			Event:         "vortex_rotate",
			RotateNum:     rotateNum,
			RotatedCount:  rotatedCount,
			Positions:     positions,
		},
	})

	return rotatedCount
}

// doVortexBlast 漩渦爆發（漩渦範圍內所有目標 70% 擊破機率）
func (g *Game) doVortexBlast(p *player.Player, instanceID string) {
	g.mu.Lock()

	type blastResult struct {
		TargetID string
		Killed   bool
		Reward   int
		X        float64
		Y        float64
	}

	var results []blastResult
	totalReward := 0
	killedCount := 0

	// 找出漩渦範圍內的所有目標
	for id, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		dx := t.X - LuckyVortexCenterX
		dy := t.Y - LuckyVortexCenterY
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist > LuckyVortexRadius {
			continue
		}

		// 70% 擊破機率
		killed := rand.Float64() < LuckyVortexBlastChance
		reward := 0

		if killed {
			t.HP = 0
			reward = int(t.Multiplier * float64(1) * LuckyVortexBlastMult)
			if reward < 1 {
				reward = 1
			}
			totalReward += reward
			killedCount++
			delete(g.Targets, id)
		}

		results = append(results, blastResult{
			TargetID: id,
			Killed:   killed,
			Reward:   reward,
			X:        t.X,
			Y:        t.Y,
		})
	}
	g.mu.Unlock()

	// 全服共享獎勵
	if totalReward > 0 {
		g.mu.RLock()
		playerCount := len(g.Players)
		players := make([]*player.Player, 0, playerCount)
		for _, pl := range g.Players {
			players = append(players, pl)
		}
		g.mu.RUnlock()

		if playerCount > 0 {
			share := totalReward / playerCount
			if share < 1 {
				share = 1
			}
			for _, pl := range players {
				pl.AddCoins(share)
			}
		}
	}

	log.Printf("[LuckyVortex] blast: killed=%d totalReward=%d", killedCount, totalReward)

	// 廣播漩渦爆發結果
	type blastResultPayload struct {
		TargetID string  `json:"target_id"`
		Killed   bool    `json:"killed"`
		Reward   int     `json:"reward"`
		X        float64 `json:"x"`
		Y        float64 `json:"y"`
	}
	payloadResults := make([]blastResultPayload, 0, len(results))
	for _, r := range results {
		payloadResults = append(payloadResults, blastResultPayload{
			TargetID: r.TargetID,
			Killed:   r.Killed,
			Reward:   r.Reward,
			X:        r.X,
			Y:        r.Y,
		})
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyVortexFish,
		Payload: ws.LuckyVortexFishPayload{
			Event:        "vortex_blast",
			KilledCount:  killedCount,
			TotalReward:  totalReward,
			BlastResults: payloadResults,
		},
	})

	// 漩渦結束廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyVortexFish,
		Payload: ws.LuckyVortexFishPayload{
			Event: "vortex_end",
		},
	})

	// 全服公告（≥3 個擊破才公告）
	if killedCount >= 3 {
		color := "#16A085"
		if killedCount >= 6 {
			color = "#1ABC9C"
		}
		ann := g.Announce.Create(announce.EventLuckyVortexFish, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🌀 漩渦爆發！擊破 %d 個目標，全服共享 %d 獎勵！",
				killedCount, totalReward),
			"color": color,
		})
		g.broadcastAnnouncement(ann)
	}
}
