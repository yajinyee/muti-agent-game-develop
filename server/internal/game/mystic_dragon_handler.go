// mystic_dragon_handler.go — 神秘龍魚八波攻擊系統（DAY-197）
// 業界依據：Ocean King 3「Mystic Dragon — Catch this fish to get 8 waves and have more
// chances to kill any fish on the screen.」
// 設計：擊破 T155 後觸發「八波龍息攻擊」：
//   1. 每波（共 8 波）：全場隨機選 3-5 個目標，每個目標 65% 擊破機率（0.55x 倍率）
//   2. 每波間隔 800ms，讓玩家看到「龍息一波一波掃過全場」的壯觀視覺
//   3. 第 8 波（最終波）：「龍怒爆發」— 全場所有目標 85% 擊破機率（0.70x 倍率）
//   4. 全服共享獎勵（按玩家數平均分配）
//   5. 全服廣播每波結果，讓所有玩家看到「龍息在掃場」
// 設計差異：
//   - 與閃電魚自動連鎖（每 0.5 秒單目標，8 秒）不同，神秘龍魚是「每波多目標」（3-5 個），
//     讓玩家感受到「龍息是面狀攻擊，不是點狀攻擊」
//   - 與鳳凰魚涅槃（一次性全場爆炸）不同，神秘龍魚是「8 波漸進式攻擊」，
//     有節奏感，讓玩家期待「下一波會打到哪些魚」
//   - 第 8 波「龍怒爆發」是全場清場，製造「最後一波最爽」的高潮設計
//   - 全服共享獎勵讓所有玩家都受益，製造「大家一起爽」的社群感
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// 神秘龍魚常數
const (
	MysticDragonCooldownSec    = 40    // 全服冷卻 40 秒
	MysticDragonWaves          = 8     // 共 8 波
	MysticDragonWaveIntervalMs = 800   // 每波間隔 800ms
	MysticDragonWaveTargetMin  = 3     // 每波最少目標數
	MysticDragonWaveTargetMax  = 5     // 每波最多目標數
	MysticDragonWaveChance     = 0.65  // 每波擊破機率 65%
	MysticDragonWaveMult       = 0.55  // 每波獎勵倍率 0.55x
	MysticDragonFinalChance    = 0.85  // 第 8 波（龍怒爆發）擊破機率 85%
	MysticDragonFinalMult      = 0.70  // 第 8 波獎勵倍率 0.70x
)

// mysticDragonManager 神秘龍魚管理器（全服冷卻）
type mysticDragonManager struct {
	mu          sync.Mutex
	isActive    bool
	cooldownEnd time.Time
}

func newMysticDragonManager() *mysticDragonManager {
	return &mysticDragonManager{}
}

// isMysticDragon 判斷是否為神秘龍魚（T155）
func isMysticDragon(defID string) bool {
	return defID == "T155"
}

// tryMysticDragonWaves 擊破 T155 後觸發八波龍息攻擊
func (g *Game) tryMysticDragonWaves(p *player.Player, instanceID string) {
	mgr := g.MysticDragon
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
		mgr.cooldownEnd = time.Now().Add(MysticDragonCooldownSec * time.Second)
		mgr.mu.Unlock()
	}()

	// 廣播龍息攻擊開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgMysticDragon,
		Payload: ws.MysticDragonPayload{
			Phase:      "dragon_start",
			TriggerID:  instanceID,
			KillerID:   p.ID,
			KillerName: p.DisplayName,
			TotalWaves: MysticDragonWaves,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("🐉 %s 擊破神秘龍魚！八波龍息攻擊開始！", p.DisplayName),
	})
	g.broadcastAnnouncement(ann)

	log.Printf("[MysticDragon] player=%s triggered 8-wave attack", p.ID)

	totalKills := 0
	totalReward := 0
	avgBet := g.getAverageBetLevelForAnglerfish() // 複用平均 betLevel 計算

	// 執行 8 波攻擊
	for wave := 1; wave <= MysticDragonWaves; wave++ {
		time.Sleep(MysticDragonWaveIntervalMs * time.Millisecond)

		isFinalWave := wave == MysticDragonWaves
		waveKills, waveReward := g.doMysticDragonWave(p, wave, isFinalWave, avgBet)
		totalKills += waveKills
		totalReward += waveReward

		// 廣播每波結果（全服）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgMysticDragon,
			Payload: ws.MysticDragonPayload{
				Phase:        fmt.Sprintf("wave_%d", wave),
				TriggerID:    instanceID,
				WaveIndex:    wave,
				WaveKills:    waveKills,
				WaveReward:   waveReward,
				TotalKills:   totalKills,
				IsFinalWave:  isFinalWave,
			},
		})

		// 第 8 波（龍怒爆發）全服公告
		if isFinalWave && waveKills >= 3 {
			ann2 := g.Announce.Create(announce.EventMegaWin, p.DisplayName, waveReward, map[string]string{
				"message": fmt.Sprintf("🐉💥 龍怒爆發！第 8 波擊破 %d 個目標！獎勵 %d 金幣！",
					waveKills, waveReward),
			})
			g.broadcastAnnouncement(ann2)
		}
	}

	// 廣播最終結果（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgMysticDragon,
		Payload: ws.MysticDragonPayload{
			Phase:       "dragon_result",
			TriggerID:   instanceID,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			TotalKills:  totalKills,
			TotalReward: totalReward,
		},
	})

	// 全服公告（≥8 個擊破）
	if totalKills >= 8 {
		ann3 := g.Announce.Create(announce.EventMegaWin, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("🐉✨ %s 神秘龍魚八波攻擊！擊破 %d 個目標！全服共享 %d 金幣！",
				p.DisplayName, totalKills, totalReward),
		})
		g.broadcastAnnouncement(ann3)
	}

	log.Printf("[MysticDragon] done: totalKills=%d totalReward=%d", totalKills, totalReward)
}

// doMysticDragonWave 執行單波龍息攻擊
func (g *Game) doMysticDragonWave(p *player.Player, waveIdx int, isFinalWave bool, avgBet int) (int, int) {
	type waveTarget struct {
		instanceID string
		multiplier float64
	}

	// 收集場上所有目標
	g.mu.RLock()
	var allTargets []waveTarget
	for _, t := range g.Targets {
		allTargets = append(allTargets, waveTarget{
			instanceID: t.InstanceID,
			multiplier: t.Multiplier,
		})
	}
	g.mu.RUnlock()

	if len(allTargets) == 0 {
		return 0, 0
	}

	// 決定本波目標數
	var selectedTargets []waveTarget
	if isFinalWave {
		// 第 8 波：全場所有目標
		selectedTargets = allTargets
	} else {
		// 普通波：隨機選 3-5 個目標
		targetCount := MysticDragonWaveTargetMin + rand.Intn(MysticDragonWaveTargetMax-MysticDragonWaveTargetMin+1)
		if targetCount > len(allTargets) {
			targetCount = len(allTargets)
		}
		// 隨機打亂後取前 N 個
		perm := rand.Perm(len(allTargets))
		for i := 0; i < targetCount; i++ {
			selectedTargets = append(selectedTargets, allTargets[perm[i]])
		}
	}

	// 決定本波擊破機率和倍率
	killChance := MysticDragonWaveChance
	killMult := MysticDragonWaveMult
	if isFinalWave {
		killChance = MysticDragonFinalChance
		killMult = MysticDragonFinalMult
	}

	kills := 0
	reward := 0
	playerCount := g.getPlayerCount()

	for _, t := range selectedTargets {
		if rand.Float64() < killChance {
			r := int(t.multiplier * float64(avgBet) * killMult)
			g.mu.Lock()
			if _, ok := g.Targets[t.instanceID]; ok {
				delete(g.Targets, t.instanceID)
				kills++
				reward += r
				// 全服共享獎勵
				for _, pl := range g.Players {
					pl.Coins += r / maxInt(1, playerCount)
				}
			}
			g.mu.Unlock()
		}
	}

	return kills, reward
}

// getPlayerCount 取得當前玩家數（thread-safe）
func (g *Game) getPlayerCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.Players)
}

// maxInt 輔助函數（避免與 max 衝突）
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
