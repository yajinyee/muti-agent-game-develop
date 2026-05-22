// golden_wave_fish_handler.go — 黃金波浪魚全場倍率衝擊系統（DAY-207）
// 業界依據：Ocean King 4 Brand New World（2025 最新版）
// 「Golden Wave Fish — triggers a golden tidal wave that sweeps across the entire screen,
//  temporarily boosting all multipliers by 2x for 8 seconds while simultaneously
//  dealing damage to all fish in its path.」
//
// 設計：擊破 T165 後觸發「黃金波浪」：
//   1. 波浪從左到右掃過全場（1.2秒），每 150ms 擊破一列目標（70% 機率，0.60x 倍率）
//   2. 波浪掃過後，全服所有玩家獲得「黃金加成」：×2.0 倍率，持續 8 秒
//   3. 黃金加成期間，所有擊破獎勵翻倍（乘法，不是加法）
//   4. 全服冷卻 50 秒
//
// 設計差異（與其他加成系統的區別）：
//   - 搖滾骷髏安可（+30% 加法）：加法加成，效果較小
//   - 幸運草魚（+50% 加法）：加法加成，效果中等
//   - 黃金波浪（×2.0 乘法）：乘法加成，效果最強（所有獎勵翻倍）
//   - 「波浪掃場 + 立即 2x 倍率」雙重效果，爆發感最強
//   - 乘法加成讓高倍率目標的獎勵更爆炸（50x 目標 → 100x 效果）
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
	"digital-twin/server/internal/ws"
)

// 黃金波浪魚常數（DAY-207）
const (
	GoldenWaveFishCooldownSec  = 50    // 全服冷卻 50 秒
	GoldenWaveFishWaveColumns  = 8     // 波浪分 8 列（每列 160px，1280px 全場）
	GoldenWaveFishColInterval  = 150   // 每列間隔 150ms（8列 × 150ms = 1.2秒）
	GoldenWaveFishKillChance   = 0.70  // 波浪擊破機率 70%
	GoldenWaveFishWaveMult     = 0.60  // 波浪擊破獎勵倍率 0.60x
	GoldenWaveFishBoostMult    = 2.0   // 黃金加成倍率 ×2.0（乘法）
	GoldenWaveFishBoostSec     = 8     // 黃金加成持續時間 8 秒
	GoldenWaveFishColWidth     = 160.0 // 每列寬度 160px（1280 / 8）
)

// goldenWaveFishManager 黃金波浪魚管理器（全服冷卻 + 加成狀態）
type goldenWaveFishManager struct {
	mu         sync.Mutex
	isActive   bool      // 是否正在波浪中
	boostEnd   time.Time // 黃金加成結束時間
	cooldownAt time.Time // 全服冷卻結束時間
}

func newGoldenWaveFishManager() *goldenWaveFishManager {
	return &goldenWaveFishManager{}
}

// isGoldenWaveFish 判斷是否為黃金波浪魚（T165，DAY-207）
func isGoldenWaveFish(defID string) bool {
	return defID == "T165"
}

// getGoldenWaveBoost 取得黃金波浪加成倍率（供 handleKill 使用）
// 黃金加成期間回傳 2.0（乘法），否則回傳 1.0（無加成）
func (g *Game) getGoldenWaveBoost() float64 {
	mgr := g.GoldenWaveFish
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if time.Now().Before(mgr.boostEnd) {
		return GoldenWaveFishBoostMult
	}
	return 1.0
}

// tryGoldenWaveFish 擊破 T165 後觸發黃金波浪
func (g *Game) tryGoldenWaveFish(p *target.Target) {
	mgr := g.GoldenWaveFish
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
	mgr.cooldownAt = time.Now().Add(GoldenWaveFishCooldownSec * time.Second)
	mgr.mu.Unlock()

	log.Printf("[GoldenWave] triggered by killing T165 at (%.0f,%.0f)", p.X, p.Y)

	// 廣播波浪開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGoldenWaveFish,
		Payload: ws.GoldenWaveFishPayload{
			Event:      "wave_start",
			Columns:    GoldenWaveFishWaveColumns,
			BoostMult:  GoldenWaveFishBoostMult,
			BoostSec:   GoldenWaveFishBoostSec,
		},
	})

	// 啟動波浪掃場 goroutine
	go g.runGoldenWaveSweep()
}

// runGoldenWaveSweep 執行黃金波浪掃場（goroutine）
func (g *Game) runGoldenWaveSweep() {
	totalKills := 0
	totalReward := 0

	for col := 0; col < GoldenWaveFishWaveColumns; col++ {
		time.Sleep(GoldenWaveFishColInterval * time.Millisecond)

		// 計算本列的 X 範圍
		colMinX := float64(col) * GoldenWaveFishColWidth
		colMaxX := colMinX + GoldenWaveFishColWidth

		// 擊破本列目標
		kills, reward := g.doGoldenWaveColumnBlast(colMinX, colMaxX)
		totalKills += kills
		totalReward += reward

		// 廣播本列結果
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgGoldenWaveFish,
			Payload: ws.GoldenWaveFishPayload{
				Event:     "wave_column",
				ColIndex:  col,
				ColX:      colMinX + GoldenWaveFishColWidth/2,
				KillCount: kills,
				Reward:    reward,
			},
		})
	}

	// 波浪掃場完成，啟動黃金加成
	g.GoldenWaveFish.mu.Lock()
	g.GoldenWaveFish.isActive = false
	g.GoldenWaveFish.boostEnd = time.Now().Add(GoldenWaveFishBoostSec * time.Second)
	g.GoldenWaveFish.mu.Unlock()

	log.Printf("[GoldenWave] sweep complete: kills=%d reward=%d, boost ×%.1f for %ds",
		totalKills, totalReward, GoldenWaveFishBoostMult, GoldenWaveFishBoostSec)

	// 廣播黃金加成開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGoldenWaveFish,
		Payload: ws.GoldenWaveFishPayload{
			Event:      "boost_start",
			BoostMult:  GoldenWaveFishBoostMult,
			BoostSec:   GoldenWaveFishBoostSec,
			TotalKills: totalKills,
			TotalReward: totalReward,
		},
	})

	// 全服公告
	color := "#FFD700"
	if totalKills >= 10 {
		color = "#FF8C00"
	}
	msg := fmt.Sprintf("🌊 黃金波浪！擊破 %d 個目標！全服 ×%.0f 倍率加成 %d 秒！",
		totalKills, GoldenWaveFishBoostMult, GoldenWaveFishBoostSec)
	ann := g.Announce.Create(announce.EventMegaWin, "", totalReward, map[string]string{
		"message": msg,
		"color":   color,
	})
	g.broadcastAnnouncement(ann)

	// 等待加成結束，廣播結束訊息
	time.Sleep(GoldenWaveFishBoostSec * time.Second)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGoldenWaveFish,
		Payload: ws.GoldenWaveFishPayload{
			Event: "boost_end",
		},
	})
	log.Printf("[GoldenWave] boost ended")
}

// doGoldenWaveColumnBlast 執行單列波浪擊破
func (g *Game) doGoldenWaveColumnBlast(minX, maxX float64) (kills, reward int) {
	type candidate struct {
		instanceID string
		multiplier float64
	}

	g.mu.RLock()
	var candidates []candidate
	for _, t := range g.Targets {
		if !t.IsAlive || isGoldenWaveFish(t.DefID) {
			continue
		}
		if t.X >= minX && t.X < maxX {
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
		if rand.Float64() >= GoldenWaveFishKillChance {
			continue
		}
		rewardAmt := int(float64(betLevel) * c.multiplier * GoldenWaveFishWaveMult)
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
			// 全服共享獎勵
			g.distributeRewardToAll(rewardAmt)
		}
		g.mu.Unlock()
	}
	return
}

// goldenWaveScreenWidth 用於計算列位置（避免 math import 警告）
var _ = math.Sqrt // 確保 math 被使用
