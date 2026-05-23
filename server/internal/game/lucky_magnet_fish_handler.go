// lucky_magnet_fish_handler.go — 幸運磁力魚系統（DAY-232）
// 業界原創「磁力聚集+磁力爆發」機制
//
// 設計：擊破 T190 後觸發「磁力場」（12 秒）：
//   - 場上所有目標物被「磁力吸引」，每 1.5 秒向場景中央移動（聚集效果）
//   - 磁力場期間擊破任何目標獲得 ×1.8 倍率加成（乘法）
//   - 12 秒後「磁力爆發」：所有聚集在中央區域（半徑 200px）的目標 75% 擊破機率（0.80x 倍率，全服共享）
//   - 個人冷卻 20 秒；全服冷卻 30 秒
//
// 設計差異：
//   - 與黑洞魚（DAY-221，重力傷害+奇點爆炸，HP 損失）不同，磁力魚是「聚集移動+磁力爆發」，
//     讓玩家有「等目標聚集再打」的策略感
//   - 「磁力吸引」讓目標物緩慢移動到中央，製造「魚群聚集」的視覺爽感
//   - 「磁力爆發」讓玩家有「等待→爆發」的高潮設計
//   - 全服廣播磁力場讓所有玩家都往中央打，製造「全服聚焦」的社交感
//   - 與傳送魚（DAY-223，瞬間移動）不同，磁力魚是「緩慢聚集」，讓玩家有「看著魚群聚集」的期待感
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
	LuckyMagnetPersonalCD    = 20 * time.Second // 個人冷卻
	LuckyMagnetGlobalCD      = 30 * time.Second // 全服冷卻
	LuckyMagnetDuration      = 12 * time.Second // 磁力場持續時間
	LuckyMagnetPullInterval  = 1500 * time.Millisecond // 每次磁力吸引間隔
	LuckyMagnetKillBoost     = 1.8              // 磁力場期間擊破倍率加成（乘法）
	LuckyMagnetBlastChance   = 0.75             // 磁力爆發擊破機率
	LuckyMagnetBlastMult     = 0.80             // 磁力爆發倍率
	LuckyMagnetBlastRadius   = 200.0            // 磁力爆發範圍（中央半徑 px）
	LuckyMagnetPullStep      = 60.0             // 每次磁力吸引移動距離（px）
	LuckyMagnetCenterX       = 500.0            // 場景中央 X
	LuckyMagnetCenterY       = 300.0            // 場景中央 Y
)

// luckyMagnetFishManager 幸運磁力魚管理器
type luckyMagnetFishManager struct {
	mu sync.Mutex

	// 個人冷卻（playerID → cooldownUntil）
	personalCooldowns map[string]time.Time

	// 全服冷卻
	globalCooldownUntil time.Time

	// 磁力場狀態
	active      bool
	activeUntil time.Time
	instanceID  string // 用於區分不同次觸發
}

func newLuckyMagnetFishManager() *luckyMagnetFishManager {
	return &luckyMagnetFishManager{
		personalCooldowns: make(map[string]time.Time),
	}
}

// isLuckyMagnetFish 判斷是否為幸運磁力魚
func isLuckyMagnetFish(defID string) bool {
	return defID == "T190"
}

// isMagnetActive 判斷磁力場是否啟動中（供 handleKill 使用）
func (g *Game) isMagnetActive() bool {
	mgr := g.LuckyMagnetFish
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

// getLuckyMagnetBoost 取得磁力場倍率加成（供 handleKill 使用）
// 磁力場啟動中時回傳 1.8，否則回傳 1.0
func (g *Game) getLuckyMagnetBoost() float64 {
	if g.isMagnetActive() {
		return LuckyMagnetKillBoost
	}
	return 1.0
}

// tryLuckyMagnetFish 擊破 T190 後觸發磁力場（供 handleKill 使用）
func (g *Game) tryLuckyMagnetFish(p *player.Player) {
	mgr := g.LuckyMagnetFish
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
	// 已有磁力場啟動中
	if mgr.active && time.Now().Before(mgr.activeUntil) {
		mgr.mu.Unlock()
		return
	}

	// 設定冷卻
	mgr.personalCooldowns[p.ID] = time.Now().Add(LuckyMagnetPersonalCD)
	mgr.globalCooldownUntil = time.Now().Add(LuckyMagnetGlobalCD)

	// 啟動磁力場
	mgr.active = true
	mgr.activeUntil = time.Now().Add(LuckyMagnetDuration)
	instanceID := fmt.Sprintf("magnet_%d", time.Now().UnixNano())
	mgr.instanceID = instanceID
	mgr.mu.Unlock()

	log.Printf("[LuckyMagnet] player=%s activated magnet field for %v", p.ID, LuckyMagnetDuration)

	// 全服廣播：磁力場開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMagnetFish,
		Payload: ws.LuckyMagnetFishPayload{
			Event:       "magnet_start",
			PlayerID:    p.ID,
			PlayerName:  p.DisplayName,
			DurationSec: int(LuckyMagnetDuration.Seconds()),
			KillBoost:   LuckyMagnetKillBoost,
			CenterX:     LuckyMagnetCenterX,
			CenterY:     LuckyMagnetCenterY,
			BlastRadius: LuckyMagnetBlastRadius,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventLuckyMagnetFish, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🧲 %s 觸發磁力場！所有目標向中央聚集，×%.1f 倍率加成！",
			p.DisplayName, LuckyMagnetKillBoost),
		"color": "#3498DB",
	})
	g.broadcastAnnouncement(ann)

	// 啟動磁力場 goroutine
	go g.runLuckyMagnetField(p, instanceID)
}

// runLuckyMagnetField 磁力場主循環（goroutine）
func (g *Game) runLuckyMagnetField(p *player.Player, instanceID string) {
	ticker := time.NewTicker(LuckyMagnetPullInterval)
	defer ticker.Stop()

	endTimer := time.NewTimer(LuckyMagnetDuration)
	defer endTimer.Stop()

	pullCount := 0

	for {
		select {
		case <-ticker.C:
			// 確認 instanceID 仍有效
			g.LuckyMagnetFish.mu.Lock()
			if g.LuckyMagnetFish.instanceID != instanceID {
				g.LuckyMagnetFish.mu.Unlock()
				return
			}
			g.LuckyMagnetFish.mu.Unlock()

			pullCount++
			movedCount := g.doMagnetPull(instanceID, pullCount)
			log.Printf("[LuckyMagnet] pull#%d: moved %d targets toward center", pullCount, movedCount)

		case <-endTimer.C:
			// 磁力場結束，觸發磁力爆發
			g.LuckyMagnetFish.mu.Lock()
			if g.LuckyMagnetFish.instanceID != instanceID {
				g.LuckyMagnetFish.mu.Unlock()
				return
			}
			g.LuckyMagnetFish.active = false
			g.LuckyMagnetFish.mu.Unlock()

			log.Printf("[LuckyMagnet] field ended, triggering magnet blast")
			g.doMagnetBlast(p, instanceID)
			return
		}
	}
}

// doMagnetPull 執行磁力吸引（所有目標向中央移動）
func (g *Game) doMagnetPull(instanceID string, pullNum int) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	movedCount := 0
	type movedTarget struct {
		ID string
		X  float64
		Y  float64
	}
	var movedTargets []movedTarget

	for id, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		// 計算向中央的方向向量
		dx := LuckyMagnetCenterX - t.X
		dy := LuckyMagnetCenterY - t.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < 20.0 {
			// 已在中央附近，不再移動
			continue
		}

		// 移動一步（最多移動 LuckyMagnetPullStep px，不超過中央）
		step := math.Min(LuckyMagnetPullStep, dist)
		t.X += (dx / dist) * step
		t.Y += (dy / dist) * step

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

		movedTargets = append(movedTargets, movedTarget{ID: id, X: t.X, Y: t.Y})
		movedCount++
	}

	if movedCount == 0 {
		return 0
	}

	// 廣播磁力吸引（含移動後的位置）
	type targetPos struct {
		ID string  `json:"id"`
		X  float64 `json:"x"`
		Y  float64 `json:"y"`
	}
	positions := make([]targetPos, 0, len(movedTargets))
	for _, mt := range movedTargets {
		positions = append(positions, targetPos{ID: mt.ID, X: mt.X, Y: mt.Y})
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMagnetFish,
		Payload: ws.LuckyMagnetFishPayload{
			Event:       "magnet_pull",
			PullNum:     pullNum,
			MovedCount:  movedCount,
			Positions:   positions,
		},
	})

	return movedCount
}

// doMagnetBlast 磁力爆發（中央區域所有目標 75% 擊破機率）
func (g *Game) doMagnetBlast(p *player.Player, instanceID string) {
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

	// 找出中央區域內的所有目標
	for id, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		dx := t.X - LuckyMagnetCenterX
		dy := t.Y - LuckyMagnetCenterY
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist > LuckyMagnetBlastRadius {
			continue
		}

		// 75% 擊破機率
		killed := rand.Float64() < LuckyMagnetBlastChance
		reward := 0

		if killed {
			t.HP = 0
			reward = int(t.Multiplier * float64(1) * LuckyMagnetBlastMult)
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

	// 全服共享獎勵（按在線玩家數平均分配）
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

	log.Printf("[LuckyMagnet] blast: killed=%d totalReward=%d", killedCount, totalReward)

	// 廣播磁力爆發結果
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
		Type: ws.MsgLuckyMagnetFish,
		Payload: ws.LuckyMagnetFishPayload{
			Event:        "magnet_blast",
			KilledCount:  killedCount,
			TotalReward:  totalReward,
			BlastResults: payloadResults,
		},
	})

	// 磁力爆發結束廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLuckyMagnetFish,
		Payload: ws.LuckyMagnetFishPayload{
			Event: "magnet_end",
		},
	})

	// 全服公告（≥3 個擊破才公告）
	if killedCount >= 3 {
		color := "#3498DB"
		if killedCount >= 6 {
			color = "#1ABC9C"
		}
		ann := g.Announce.Create(announce.EventLuckyMagnetFish, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🧲 磁力爆發！擊破 %d 個目標，全服共享 %d 獎勵！",
				killedCount, totalReward),
			"color": color,
		})
		g.broadcastAnnouncement(ann)
	}
}
