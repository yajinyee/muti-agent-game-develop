// thunderbolt_lobster_v2_handler.go — 雷霆龍蝦免費射擊系統 V2（DAY-199）
// 業界依據：Royal Fishing JILI「Thunderbolt Lobster feature — provides 15 seconds of free play
// followed by automatic shooting from the Thunderbolt Turret.
// Players can earn extra seconds during this period to extend gameplay and increase reward potential.
// This feature is triggered when players shoot the Thunderbolt Lobster,
// which explodes, releasing a burst of energy that powers up the game.」
//
// 設計：擊破 T157 後觸發「雷霆砲台模式」（15秒，全服共享）：
//   1. 系統自動每 0.5 秒選最高價值目標射擊（85% 擊破機率，0.75x 倍率）
//   2. 每擊破一個目標 +0.5 秒（最多延長到 30 秒）
//   3. 全服廣播「雷霆砲台啟動」，讓所有玩家看到自動射擊效果
//   4. 結束時廣播結算（擊破數/總獎勵/延長秒數）
//
// 設計差異（與 DAY-150 T114 的區別）：
//   - DAY-150 T114：個人觸發，個人免費射擊（per-player session）
//   - DAY-199 T157：全服共享，觸發後全服所有玩家都看到砲台在自動射擊
//   - 「延長時間」機制讓玩家感受到「打得越準，免費射擊越久」的技巧感
//   - 全服廣播讓觸發者有「我在幫全服打魚」的英雄感
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

// 雷霆龍蝦 V2 常數（全服共享版，DAY-199）
const (
	TBLobsterV2CooldownSec  = 45   // 全服冷卻 45 秒
	TBLobsterV2BaseDuration = 15.0 // 基礎持續時間 15 秒
	TBLobsterV2MaxDuration  = 30.0 // 最大持續時間 30 秒
	TBLobsterV2Interval     = 500  // 自動射擊間隔 500ms
	TBLobsterV2KillChance   = 0.85 // 擊破機率 85%
	TBLobsterV2Mult         = 0.75 // 獎勵倍率 0.75x
	TBLobsterV2ExtendSec    = 0.5  // 每次擊破延長 0.5 秒
)

// tbLobsterV2Session 雷霆龍蝦 V2 會話（全服共享）
type tbLobsterV2Session struct {
	killerID    string    // 觸發者玩家 ID
	killerName  string    // 觸發者顯示名稱
	betLevel    int       // 觸發者 betLevel（用於計算獎勵）
	endTime     time.Time // 結束時間（可延長）
	killCount   int       // 擊破數
	totalReward int       // 總獎勵
	extendSec   float64   // 已延長秒數
}

// tbLobsterV2Manager 雷霆龍蝦 V2 管理器（全服共享）
type tbLobsterV2Manager struct {
	mu          sync.Mutex
	session     *tbLobsterV2Session
	cooldownEnd time.Time
}

func newTBLobsterV2Manager() *tbLobsterV2Manager {
	return &tbLobsterV2Manager{}
}

// isThunderboltLobsterV2 判斷是否為雷霆龍蝦 V2（T157，DAY-199）
func isThunderboltLobsterV2(defID string) bool {
	return defID == "T157"
}

// tryThunderboltLobsterFreePlay 擊破 T157 後觸發雷霆砲台模式（全服共享）
func (g *Game) tryThunderboltLobsterFreePlay(p *player.Player, multiplier float64) {
	mgr := g.ThunderboltLobsterV2
	mgr.mu.Lock()

	// 全服冷卻檢查
	if mgr.session != nil || time.Now().Before(mgr.cooldownEnd) {
		mgr.mu.Unlock()
		return
	}

	session := &tbLobsterV2Session{
		killerID:   p.ID,
		killerName: p.DisplayName,
		betLevel:   p.BetLevel,
		endTime:    time.Now().Add(TBLobsterV2BaseDuration * time.Second),
	}
	mgr.session = session
	mgr.mu.Unlock()

	log.Printf("[TBLobsterV2] player=%s triggered free play (betLevel=%d)", p.ID, p.BetLevel)

	// 全服廣播：雷霆砲台啟動
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgThunderboltLobster,
		Payload: ws.ThunderboltLobsterPayload{
			Event:       "turret_start",
			KillerName:  p.DisplayName,
			Duration:    TBLobsterV2BaseDuration,
			MaxDuration: TBLobsterV2MaxDuration,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, 0, map[string]string{
		"message": fmt.Sprintf("⚡ %s 觸發雷霆砲台！自動射擊 15 秒！", p.DisplayName),
	})
	g.broadcastAnnouncement(ann)

	// 啟動自動射擊 goroutine
	go g.runTBLobsterV2Turret(session)
}

// runTBLobsterV2Turret 雷霆砲台自動射擊主循環
func (g *Game) runTBLobsterV2Turret(session *tbLobsterV2Session) {
	ticker := time.NewTicker(TBLobsterV2Interval * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mgr := g.ThunderboltLobsterV2
			mgr.mu.Lock()
			if mgr.session != session {
				mgr.mu.Unlock()
				return
			}
			// 檢查是否超時
			if time.Now().After(session.endTime) {
				mgr.session = nil
				mgr.cooldownEnd = time.Now().Add(TBLobsterV2CooldownSec * time.Second)
				killCount := session.killCount
				totalReward := session.totalReward
				extendSec := session.extendSec
				mgr.mu.Unlock()

				// 廣播結算
				g.Hub.Broadcast(&ws.Message{
					Type: ws.MsgThunderboltLobster,
					Payload: ws.ThunderboltLobsterPayload{
						Event:       "turret_end",
						KillCount:   killCount,
						TotalReward: int64(totalReward),
						ExtendSec:   extendSec,
					},
				})

				// 全服公告（≥5 擊破才公告）
				if killCount >= 5 {
					ann := g.Announce.Create(announce.EventMegaWin, "雷霆砲台", totalReward, map[string]string{
						"message": fmt.Sprintf("⚡ 雷霆砲台結束！共擊破 %d 個目標！獎勵 %d 金幣！", killCount, totalReward),
					})
					g.broadcastAnnouncement(ann)
				}
				return
			}
			mgr.mu.Unlock()

			// 執行一次自動射擊
			g.doTBLobsterV2Shot(session)

		case <-g.stopCh:
			return
		}
	}
}

// doTBLobsterV2Shot 執行一次雷霆砲台自動射擊
func (g *Game) doTBLobsterV2Shot(session *tbLobsterV2Session) {
	g.mu.RLock()

	type targetInfo struct {
		id   string
		mult float64
	}
	var candidates []targetInfo

	for _, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		// 跳過 BOSS（避免干擾主要遊戲流程）
		if t.DefID == "B001" {
			continue
		}
		// 跳過幽靈魚幻影分身
		if isGhostFishClone(t.DefID) {
			continue
		}
		candidates = append(candidates, targetInfo{
			id:   t.InstanceID,
			mult: t.Multiplier,
		})
	}
	g.mu.RUnlock()

	if len(candidates) == 0 {
		return
	}

	// 從前 40% 高倍率目標中隨機選一個（偏向高價值）
	topN := len(candidates) * 2 / 5
	if topN < 1 {
		topN = 1
	}
	// 簡單選擇排序找前 topN 個高倍率目標
	for i := 0; i < topN; i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].mult > candidates[i].mult {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}
	chosen := candidates[rand.Intn(topN)]

	// 85% 擊破機率
	killed := rand.Float64() < TBLobsterV2KillChance

	g.mu.Lock()
	t, ok := g.Targets[chosen.id]
	if !ok || t.HP <= 0 {
		g.mu.Unlock()
		return
	}

	var reward int
	if killed {
		reward = int(t.Multiplier * float64(session.betLevel) * TBLobsterV2Mult)
		if reward < 1 {
			reward = 1
		}
		delete(g.Targets, chosen.id)
	}
	g.mu.Unlock()

	if killed {
		mgr := g.ThunderboltLobsterV2
		mgr.mu.Lock()
		if mgr.session == session {
			session.killCount++
			session.totalReward += reward
			session.extendSec += TBLobsterV2ExtendSec
			// 延長時間（不超過最大值）
			remaining := time.Until(session.endTime).Seconds()
			if remaining+TBLobsterV2ExtendSec <= TBLobsterV2MaxDuration {
				session.endTime = session.endTime.Add(
					time.Duration(TBLobsterV2ExtendSec * float64(time.Second)),
				)
			}
			killCount := session.killCount
			totalReward := session.totalReward
			remaining2 := time.Until(session.endTime).Seconds()
			mgr.mu.Unlock()

			// 廣播射擊結果
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgThunderboltLobster,
				Payload: ws.ThunderboltLobsterPayload{
					Event:       "turret_shot",
					TargetID:    chosen.id,
					Killed:      true,
					Reward:      int64(reward),
					KillCount:   killCount,
					TotalReward: int64(totalReward),
					Remaining:   remaining2,
				},
			})

			// 給觸發者獎勵
			g.mu.Lock()
			if p2, ok2 := g.Players[session.killerID]; ok2 {
				p2.Coins += reward
			}
			g.mu.Unlock()

			// 全服公告（每 10 次擊破公告一次）
			if killCount%10 == 0 {
				ann := g.Announce.Create(announce.EventMegaWin, "雷霆砲台", totalReward, map[string]string{
					"message": fmt.Sprintf("⚡ 雷霆砲台已擊破 %d 個目標！", killCount),
				})
				g.broadcastAnnouncement(ann)
			}
		} else {
			mgr.mu.Unlock()
		}
	} else {
		// 未擊破，廣播射擊嘗試
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgThunderboltLobster,
			Payload: ws.ThunderboltLobsterPayload{
				Event:    "turret_shot",
				TargetID: chosen.id,
				Killed:   false,
			},
		})
	}
}
