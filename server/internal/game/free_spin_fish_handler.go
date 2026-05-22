// free_spin_fish_handler.go — 自由旋轉魚免費射擊系統（DAY-204）
// 業界依據：Galaxsys King of Ocean 2026「Free Spin Fish triggers bonus rounds」
// 「Free Spin Fish, Captain Fish, and Money Fish trigger bonus rounds, extra multipliers, and instant payouts.」
//
// 設計：擊破 T162 後觸發「個人免費射擊模式」（10秒，不扣費）：
//   1. 系統每 0.6 秒自動選最高價值目標射擊（80% 擊破機率，0.80x 倍率）
//   2. 每擊破一個目標 +1 秒（最多延長到 20 秒）
//   3. 個人冷卻 30 秒（不影響其他玩家）
//   4. 結束時廣播結算（擊破數/總獎勵/延長秒數）
//
// 設計差異（與 T157 雷霆龍蝦 V2 的區別）：
//   - T157：全服共享砲台，觸發後全服所有玩家都看到自動射擊，0.75x 倍率
//   - T162：個人免費射擊，只有觸發者受益，0.80x 倍率（更高），不扣費（真正免費）
//   - 「不扣費」是核心設計：讓玩家感受到「白嫖的爽感」
//   - 個人冷卻讓每個玩家都有機會觸發，不會被一個玩家壟斷
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

// 自由旋轉魚常數（DAY-204）
const (
	FreeSpinFishCooldownSec  = 30   // 個人冷卻 30 秒
	FreeSpinFishBaseDuration = 10.0 // 基礎持續時間 10 秒
	FreeSpinFishMaxDuration  = 20.0 // 最大持續時間 20 秒
	FreeSpinFishInterval     = 600  // 自動射擊間隔 600ms
	FreeSpinFishKillChance   = 0.80 // 擊破機率 80%
	FreeSpinFishMult         = 0.80 // 獎勵倍率 0.80x（不扣費，純獎勵）
	FreeSpinFishExtendSec    = 1.0  // 每次擊破延長 1 秒
)

// freeSpinFishSession 自由旋轉魚個人 session
type freeSpinFishSession struct {
	playerID    string    // 觸發者玩家 ID
	betLevel    int       // 觸發時的 betLevel（用於計算獎勵）
	endTime     time.Time // 結束時間（可延長）
	killCount   int       // 擊破數
	totalReward int       // 總獎勵
	extendSec   float64   // 已延長秒數
}

// freeSpinFishManager 自由旋轉魚管理器（個人 session）
type freeSpinFishManager struct {
	mu       sync.Mutex
	sessions map[string]*freeSpinFishSession // playerID → session
	cooldowns map[string]time.Time           // playerID → 冷卻結束時間
}

func newFreeSpinFishManager() *freeSpinFishManager {
	return &freeSpinFishManager{
		sessions:  make(map[string]*freeSpinFishSession),
		cooldowns: make(map[string]time.Time),
	}
}

// isFreeSpinFish 判斷是否為自由旋轉魚（T162，DAY-204）
func isFreeSpinFish(defID string) bool {
	return defID == "T162"
}

// tryFreeSpinFishMode 擊破 T162 後觸發個人免費射擊模式
func (g *Game) tryFreeSpinFishMode(p *player.Player, multiplier float64) {
	mgr := g.FreeSpinFish
	mgr.mu.Lock()

	// 個人冷卻檢查
	if cd, ok := mgr.cooldowns[p.ID]; ok && time.Now().Before(cd) {
		mgr.mu.Unlock()
		return
	}
	// 已有活躍 session 則不重複觸發
	if _, ok := mgr.sessions[p.ID]; ok {
		mgr.mu.Unlock()
		return
	}

	session := &freeSpinFishSession{
		playerID: p.ID,
		betLevel: p.BetLevel,
		endTime:  time.Now().Add(FreeSpinFishBaseDuration * time.Second),
	}
	mgr.sessions[p.ID] = session
	mgr.mu.Unlock()

	log.Printf("[FreeSpinFish] player=%s triggered free spin mode (betLevel=%d)", p.ID, p.BetLevel)

	// 廣播：免費射擊開始（個人）
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgFreeSpinFish,
		Payload: ws.FreeSpinFishPayload{
			Event:       "free_spin_start",
			PlayerID:    p.ID,
			Duration:    FreeSpinFishBaseDuration,
			MaxDuration: FreeSpinFishMaxDuration,
		},
	}); err != nil {
		log.Printf("[FreeSpinFish] send free_spin_start error: %v", err)
	}

	// 全服廣播（讓其他玩家知道有人觸發了免費射擊）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgFreeSpinFish,
		Payload: ws.FreeSpinFishPayload{
			Event:      "free_spin_broadcast",
			PlayerName: p.DisplayName,
		},
	})

	// 啟動免費射擊 goroutine
	go g.runFreeSpinFishMode(p, session)
}

// runFreeSpinFishMode 執行免費射擊模式（goroutine）
func (g *Game) runFreeSpinFishMode(p *player.Player, session *freeSpinFishSession) {
	ticker := time.NewTicker(FreeSpinFishInterval * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 檢查是否超時
			mgr := g.FreeSpinFish
			mgr.mu.Lock()
			s, ok := mgr.sessions[p.ID]
			if !ok || s != session {
				mgr.mu.Unlock()
				return
			}
			if time.Now().After(s.endTime) {
				// 超時結束
				killCount := s.killCount
				totalReward := s.totalReward
				extendSec := s.extendSec
				delete(mgr.sessions, p.ID)
				mgr.cooldowns[p.ID] = time.Now().Add(FreeSpinFishCooldownSec * time.Second)
				mgr.mu.Unlock()
				g.finalizeFreeSpinFish(p, killCount, totalReward, extendSec)
				return
			}
			mgr.mu.Unlock()

			// 執行一次免費射擊
			g.doFreeSpinFishShot(p, session)

		case <-g.stopCh:
			return
		}
	}
}

// doFreeSpinFishShot 執行一次免費射擊（選最高價值目標）
func (g *Game) doFreeSpinFishShot(p *player.Player, session *freeSpinFishSession) {
	// 選最高價值目標（前 30% 中隨機選）
	g.mu.RLock()
	type targetInfo struct {
		id   string
		mult int
		x    float64
		y    float64
	}
	var candidates []targetInfo
	for _, t := range g.Targets {
		if t.HP > 0 && t.Def.Type != "boss" {
			candidates = append(candidates, targetInfo{
				id:   t.InstanceID,
				mult: int(t.Multiplier),
				x:    t.X,
				y:    t.Y,
			})
		}
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		return
	}

	// 按倍率排序，選前 30%
	for i := 0; i < len(candidates)-1; i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].mult > candidates[i].mult {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}
	topN := len(candidates) / 3
	if topN < 1 {
		topN = 1
	}
	chosen := candidates[rand.Intn(topN)]

	// 80% 擊破機率
	killed := rand.Float64() < FreeSpinFishKillChance
	reward := 0

	if killed {
		// 計算獎勵（不扣費，純獎勵）
		reward = int(float64(chosen.mult) * float64(session.betLevel) * FreeSpinFishMult)

		// 給予獎勵
		g.mu.Lock()
		p.Coins += reward
		g.mu.Unlock()

		// 更新 session
		mgr := g.FreeSpinFish
		mgr.mu.Lock()
		if s, ok := mgr.sessions[p.ID]; ok && s == session {
			s.killCount++
			s.totalReward += reward
			// 延長時間
			remaining := time.Until(s.endTime).Seconds()
			newRemaining := remaining + FreeSpinFishExtendSec
			if newRemaining > FreeSpinFishMaxDuration {
				newRemaining = FreeSpinFishMaxDuration
			}
			s.endTime = time.Now().Add(time.Duration(newRemaining * float64(time.Second)))
			s.extendSec += FreeSpinFishExtendSec
		}
		mgr.mu.Unlock()
	}

	// 廣播射擊結果（個人）
	remaining := time.Until(session.endTime).Seconds()
	if remaining < 0 {
		remaining = 0
	}
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgFreeSpinFish,
		Payload: ws.FreeSpinFishPayload{
			Event:       "free_spin_shot",
			TargetID:    chosen.id,
			TargetX:     chosen.x,
			TargetY:     chosen.y,
			Killed:      killed,
			Reward:      reward,
			KillCount:   session.killCount,
			Remaining:   remaining,
		},
	}); err != nil {
		log.Printf("[FreeSpinFish] send free_spin_shot error: %v", err)
	}
}

// finalizeFreeSpinFish 結算免費射擊模式
func (g *Game) finalizeFreeSpinFish(p *player.Player, killCount, totalReward int, extendSec float64) {
	log.Printf("[FreeSpinFish] player=%s ended: kills=%d reward=%d extend=%.1fs",
		p.ID, killCount, totalReward, extendSec)

	// 廣播結算（個人）
	if err := g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgFreeSpinFish,
		Payload: ws.FreeSpinFishPayload{
			Event:       "free_spin_end",
			KillCount:   killCount,
			TotalReward: totalReward,
			ExtendSec:   extendSec,
		},
	}); err != nil {
		log.Printf("[FreeSpinFish] send free_spin_end error: %v", err)
	}

	// 全服公告（≥5 個擊破時）
	if killCount >= 5 {
		ann := g.Announce.Create(announce.EventBigWin, p.DisplayName, totalReward, map[string]string{
			"message": formatFreeSpinFishAnnounce(p.DisplayName, killCount, totalReward),
			"color":   freeSpinFishColor(killCount),
		})
		g.broadcastAnnouncement(ann)
	}
}

// formatFreeSpinFishAnnounce 格式化全服公告文字
func formatFreeSpinFishAnnounce(playerName string, killCount, totalReward int) string {
	switch {
	case killCount >= 15:
		return fmt.Sprintf("🌀 %s 自由旋轉魚免費射擊擊破 %d 個目標！獲得 %d 金幣！", playerName, killCount, totalReward)
	case killCount >= 10:
		return fmt.Sprintf("🌀 %s 免費射擊擊破 %d 個目標！獲得 %d 金幣！", playerName, killCount, totalReward)
	default:
		return fmt.Sprintf("🌀 %s 免費射擊獲得 %d 金幣！", playerName, totalReward)
	}
}

// freeSpinFishColor 依擊破數決定公告顏色
func freeSpinFishColor(killCount int) string {
	switch {
	case killCount >= 15:
		return "#FF4500" // 橙紅色（超多擊破）
	case killCount >= 10:
		return "#FFD700" // 金色
	default:
		return "#00CED1" // 青色（自由旋轉魚主題色）
	}
}
