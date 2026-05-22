// ice_phoenix_handler.go — 冰鳳凰覺醒 BOSS 系統（DAY-200）
// 業界依據：Royal Fishing JILI「Ice Phoenix Awaken Feature — fixed jackpot mechanic
// that awards up to 300x the bet when players eliminate the Ice Phoenix boss.
// Multicoloured phoenix (blue, pink, purple, orange) with magical aura.
// Awaken Boss with 30x basic multiplier. Power Up attack delivers 6x-10x boost
// for rewards up to 300 times bet.」
//
// 設計：擊破 T158 後觸發「冰鳳凰覺醒」：
//   1. 基礎獎勵：30x betLevel（固定）
//   2. Power Up 攻擊：隨機選 3-5 個目標，每個 6-10x 倍率（70% 擊破機率）
//   3. 5% 機率觸發「冰霜爆發」：全場所有目標受到冰霜攻擊（50% 擊破機率，0.60x 倍率）
//   4. 全服廣播冰鳳凰覺醒，讓所有玩家看到冰霜特效
//
// 設計差異：
//   - 與鳳凰魚涅槃（全場爆炸，一次性）不同，冰鳳凰是「覺醒 BOSS + Power Up 攻擊」雙段式
//   - 與神秘龍魚（8波攻擊）不同，冰鳳凰是「精準 Power Up」（3-5個目標，高倍率）
//   - 「冰霜爆發」（5%）讓玩家有「說不定這次全場清場」的期待感
//   - 業界最高 300x 的設計：30x 基礎 + 10x × 5 目標 + 冰霜爆發加成
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

// 冰鳳凰常數
const (
	IcePhoenixCooldownSec    = 50    // 全服冷卻 50 秒
	IcePhoenixBaseReward     = 30    // 基礎獎勵 30x betLevel
	IcePhoenixPowerUpMin     = 3     // Power Up 最少目標數
	IcePhoenixPowerUpMax     = 5     // Power Up 最多目標數
	IcePhoenixPowerUpMultMin = 6.0   // Power Up 倍率最小值
	IcePhoenixPowerUpMultMax = 10.0  // Power Up 倍率最大值
	IcePhoenixPowerUpChance  = 0.70  // Power Up 擊破機率 70%
	IcePhoenixFrostChance    = 0.05  // 冰霜爆發機率 5%
	IcePhoenixFrostKillChance = 0.50 // 冰霜爆發擊破機率 50%
	IcePhoenixFrostMult      = 0.60  // 冰霜爆發獎勵倍率 0.60x
)

// icePhoenixManager 冰鳳凰管理器（全服共享）
type icePhoenixManager struct {
	mu          sync.Mutex
	isActive    bool
	cooldownEnd time.Time
}

func newIcePhoenixManager() *icePhoenixManager {
	return &icePhoenixManager{}
}

// isIcePhoenix 判斷是否為冰鳳凰（T158）
func isIcePhoenix(defID string) bool {
	return defID == "T158"
}

// tryIcePhoenixAwaken 擊破 T158 後觸發冰鳳凰覺醒
func (g *Game) tryIcePhoenixAwaken(p *player.Player, multiplier float64) {
	mgr := g.IcePhoenix
	mgr.mu.Lock()

	// 全服冷卻檢查
	if mgr.isActive || time.Now().Before(mgr.cooldownEnd) {
		mgr.mu.Unlock()
		return
	}
	mgr.isActive = true
	mgr.mu.Unlock()

	log.Printf("[IcePhoenix] player=%s triggered awaken (betLevel=%d)", p.ID, p.BetLevel)

	// 全服廣播：冰鳳凰覺醒開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgIcePhoenix,
		Payload: ws.IcePhoenixPayload{
			Event:       "awaken_start",
			KillerName:  p.DisplayName,
			BaseReward:  IcePhoenixBaseReward * p.BetLevel,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, IcePhoenixBaseReward*p.BetLevel, map[string]string{
		"message": fmt.Sprintf("❄️ %s 觸發冰鳳凰覺醒！Power Up 攻擊開始！", p.DisplayName),
	})
	g.broadcastAnnouncement(ann)

	// 給基礎獎勵
	baseReward := IcePhoenixBaseReward * p.BetLevel
	g.mu.Lock()
	if pp, ok := g.Players[p.ID]; ok {
		pp.Coins += baseReward
	}
	g.mu.Unlock()

	// 啟動 Power Up 攻擊 goroutine
	go g.runIcePhoenixPowerUp(p, baseReward)
}

// runIcePhoenixPowerUp 冰鳳凰 Power Up 攻擊
func (g *Game) runIcePhoenixPowerUp(p *player.Player, baseReward int) {
	// 決定 Power Up 目標數（3-5 個）
	targetCount := IcePhoenixPowerUpMin + rand.Intn(IcePhoenixPowerUpMax-IcePhoenixPowerUpMin+1)

	// 收集候選目標（排除 BOSS 和幻影分身）
	g.mu.RLock()
	type targetInfo struct {
		id   string
		mult float64
	}
	var candidates []targetInfo
	for _, t := range g.Targets {
		if t.HP <= 0 || t.DefID == "B001" || isGhostFishClone(t.DefID) {
			continue
		}
		candidates = append(candidates, targetInfo{id: t.InstanceID, mult: t.Multiplier})
	}
	g.mu.RUnlock()

	// 隨機選取目標（不重複）
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	if len(candidates) > targetCount {
		candidates = candidates[:targetCount]
	}

	// 執行 Power Up 攻擊（每個目標間隔 400ms）
	powerUpKills := 0
	powerUpReward := 0
	powerUpResults := make([]ws.IcePhoenixPowerUpResult, 0, len(candidates))

	for i, c := range candidates {
		time.Sleep(400 * time.Millisecond)

		// 6-10x 倍率（隨機）
		mult := IcePhoenixPowerUpMultMin + rand.Float64()*(IcePhoenixPowerUpMultMax-IcePhoenixPowerUpMultMin)
		killed := rand.Float64() < IcePhoenixPowerUpChance

		var reward int
		if killed {
			g.mu.Lock()
			t, ok := g.Targets[c.id]
			if ok && t.HP > 0 {
				reward = int(mult * float64(p.BetLevel))
				if reward < 1 {
					reward = 1
				}
				delete(g.Targets, c.id)
				powerUpKills++
				powerUpReward += reward
				// 給觸發者獎勵
				if pp, ok2 := g.Players[p.ID]; ok2 {
					pp.Coins += reward
				}
			} else {
				killed = false
			}
			g.mu.Unlock()
		}

		result := ws.IcePhoenixPowerUpResult{
			TargetID: c.id,
			Mult:     mult,
			Killed:   killed,
			Reward:   reward,
		}
		powerUpResults = append(powerUpResults, result)

		// 廣播每次 Power Up 攻擊
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgIcePhoenix,
			Payload: ws.IcePhoenixPayload{
				Event:       "power_up_shot",
				ShotIndex:   i + 1,
				TotalShots:  len(candidates),
				PowerUpResult: result,
			},
		})
	}

	// 檢查是否觸發冰霜爆發（5%）
	frostKills := 0
	frostReward := 0
	hasFrost := rand.Float64() < IcePhoenixFrostChance

	if hasFrost {
		time.Sleep(600 * time.Millisecond)

		// 全服廣播：冰霜爆發開始
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgIcePhoenix,
			Payload: ws.IcePhoenixPayload{
				Event: "frost_burst_start",
			},
		})

		// 全場所有目標受到冰霜攻擊
		g.mu.RLock()
		var allTargets []string
		for id, t := range g.Targets {
			if t.HP > 0 && t.DefID != "B001" && !isGhostFishClone(t.DefID) {
				allTargets = append(allTargets, id)
			}
		}
		g.mu.RUnlock()

		for _, id := range allTargets {
			time.Sleep(60 * time.Millisecond)
			if rand.Float64() < IcePhoenixFrostKillChance {
				g.mu.Lock()
				t, ok := g.Targets[id]
				if ok && t.HP > 0 {
					r := int(t.Multiplier * float64(p.BetLevel) * IcePhoenixFrostMult)
					if r < 1 {
						r = 1
					}
					delete(g.Targets, id)
					frostKills++
					frostReward += r
					if pp, ok2 := g.Players[p.ID]; ok2 {
						pp.Coins += r
					}
				}
				g.mu.Unlock()
			}
		}

		// 廣播冰霜爆發結果
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgIcePhoenix,
			Payload: ws.IcePhoenixPayload{
				Event:       "frost_burst_result",
				FrostKills:  frostKills,
				FrostReward: frostReward,
			},
		})

		// 全服公告（≥3 個冰霜擊破）
		if frostKills >= 3 {
			ann2 := g.Announce.Create(announce.EventMegaWin, p.DisplayName, frostReward, map[string]string{
				"message": fmt.Sprintf("❄️💥 冰霜爆發！擊破 %d 個目標！獎勵 %d 金幣！", frostKills, frostReward),
			})
			g.broadcastAnnouncement(ann2)
		}
	}

	// 最終結算
	totalReward := baseReward + powerUpReward + frostReward
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgIcePhoenix,
		Payload: ws.IcePhoenixPayload{
			Event:         "awaken_result",
			KillerName:    p.DisplayName,
			BaseReward:    baseReward,
			PowerUpKills:  powerUpKills,
			PowerUpReward: powerUpReward,
			FrostKills:    frostKills,
			FrostReward:   frostReward,
			TotalReward:   totalReward,
			HasFrost:      hasFrost,
		},
	})

	// 全服公告（≥100x 才公告）
	if totalReward >= 100*p.BetLevel {
		ann3 := g.Announce.Create(announce.EventMegaWin, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("❄️✨ %s 冰鳳凰覺醒！總獎勵 %d 金幣！", p.DisplayName, totalReward),
		})
		g.broadcastAnnouncement(ann3)
	}

	// 重置管理器
	mgr := g.IcePhoenix
	mgr.mu.Lock()
	mgr.isActive = false
	mgr.cooldownEnd = time.Now().Add(IcePhoenixCooldownSec * time.Second)
	mgr.mu.Unlock()

	log.Printf("[IcePhoenix] awaken complete: powerUp=%d/%d frost=%d/%d total=%d",
		powerUpKills, len(candidates), frostKills, 0, totalReward)
}
