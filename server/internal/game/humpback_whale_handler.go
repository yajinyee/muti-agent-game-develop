// humpback_whale_handler.go — 座頭鯨覺醒系統（DAY-203）
// 業界依據：Royal Fishing JILI「Humpback Whale offers 90-150x with 15x base multiplier.
// Awaken Boss mechanic — when defeated, triggers a powerful wave attack that sweeps the screen.
// The Humpback Whale's signature breach mechanic creates massive splash zones.」
//
// 設計：擊破 T161 後觸發「鯨歌覺醒」：
//   1. 基礎獎勵：15x betLevel（立即給予）
//   2. 鯨歌波浪攻擊：3 波，每波 3 個目標（65% 擊破機率，0.60x 倍率），每波間隔 1 秒
//   3. 5% 機率觸發「深海巨浪」：全場所有目標（60% 擊破機率，0.65x 倍率）
//   4. 最高組合：150x betLevel
//
// 設計差異：
//   - 與冰鳳凰（Power Up 精準攻擊 3-5 個目標）不同，座頭鯨是「波浪式掃場」（3 波 × 3 個目標），
//     讓玩家感受到「鯨魚在場上掀起巨浪」的壯觀感
//   - 與神秘龍魚（8 波，每波 3-5 個目標）不同，座頭鯨是「3 波快速掃場」，節奏更緊湊
//   - 「深海巨浪」（5%）讓玩家有「說不定這次全場清場」的期待感
//   - 基礎獎勵 15x 確保玩家擊破後立刻有回報，不會空手
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

// 座頭鯨覺醒常數
const (
	HumpbackWhaleCooldownSec  = 45    // 全服冷卻 45 秒
	HumpbackWhaleBaseReward   = 15    // 基礎獎勵 15x betLevel
	HumpbackWhaleWaves        = 3     // 波浪攻擊 3 波
	HumpbackWhaleTargetsPerWave = 3   // 每波 3 個目標
	HumpbackWhaleWaveIntervalMs = 1000 // 每波間隔 1 秒
	HumpbackWhaleWaveChance   = 0.65  // 波浪擊破機率 65%
	HumpbackWhaleWaveMult     = 0.60  // 波浪獎勵倍率 0.60x
	HumpbackWhaleTidalChance  = 0.05  // 深海巨浪觸發機率 5%
	HumpbackWhaleTidalKillChance = 0.60 // 深海巨浪擊破機率 60%
	HumpbackWhaleTidalMult    = 0.65  // 深海巨浪獎勵倍率 0.65x
	HumpbackWhaleTidalIntervalMs = 50 // 深海巨浪每個目標間隔 50ms
)

// humpbackWhaleManager 座頭鯨覺醒管理器（全服共享）
type humpbackWhaleManager struct {
	mu          sync.Mutex
	isActive    bool
	cooldownEnd time.Time
}

func newHumpbackWhaleManager() *humpbackWhaleManager {
	return &humpbackWhaleManager{}
}

// isHumpbackWhale 判斷是否為座頭鯨（T161）
func isHumpbackWhale(defID string) bool {
	return defID == "T161"
}

// tryHumpbackWhaleAwaken 擊破 T161 後觸發鯨歌覺醒
func (g *Game) tryHumpbackWhaleAwaken(p *player.Player, mult float64) {
	mgr := g.HumpbackWhale
	mgr.mu.Lock()

	// 全服冷卻檢查
	if mgr.isActive || time.Now().Before(mgr.cooldownEnd) {
		mgr.mu.Unlock()
		return
	}
	mgr.isActive = true
	mgr.mu.Unlock()

	log.Printf("[HumpbackWhale] player=%s triggered awaken mult=%.1f", p.ID, mult)

	// 基礎獎勵：15x betLevel（立即給予）
	baseReward := HumpbackWhaleBaseReward * p.BetLevel
	g.mu.Lock()
	if pp, ok := g.Players[p.ID]; ok {
		pp.Coins += baseReward
	}
	g.mu.Unlock()

	// 廣播覺醒開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgHumpbackWhale,
		Payload: ws.HumpbackWhalePayload{
			Event:       "awaken_start",
			KillerName:  p.DisplayName,
			BaseReward:  baseReward,
			WaveCount:   HumpbackWhaleWaves,
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, baseReward, map[string]string{
		"message": fmt.Sprintf("🐋 %s 觸發座頭鯨覺醒！鯨歌波浪攻擊開始！", p.DisplayName),
	})
	g.broadcastAnnouncement(ann)

	// 執行鯨歌波浪攻擊
	go g.runHumpbackWhaleWaves(p, baseReward)
}

// runHumpbackWhaleWaves 執行鯨歌波浪攻擊（goroutine）
func (g *Game) runHumpbackWhaleWaves(p *player.Player, baseReward int) {
	totalReward := baseReward
	totalKills := 0
	hasTidal := rand.Float64() < HumpbackWhaleTidalChance

	for wave := 1; wave <= HumpbackWhaleWaves; wave++ {
		time.Sleep(HumpbackWhaleWaveIntervalMs * time.Millisecond)

		// 選取本波目標（隨機選 3 個存活目標）
		g.mu.RLock()
		type waveTarget struct {
			id   string
			mult float64
		}
		var candidates []waveTarget
		for _, t := range g.Targets {
			if t.HP <= 0 || t.DefID == "B001" || isGhostFishClone(t.DefID) {
				continue
			}
			candidates = append(candidates, waveTarget{id: t.InstanceID, mult: t.Multiplier})
		}
		g.mu.RUnlock()

		// 隨機打亂並取前 N 個
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})
		if len(candidates) > HumpbackWhaleTargetsPerWave {
			candidates = candidates[:HumpbackWhaleTargetsPerWave]
		}

		// 對選中目標執行波浪攻擊
		waveKills := 0
		waveReward := 0
		for _, ct := range candidates {
			if rand.Float64() < HumpbackWhaleWaveChance {
				g.mu.Lock()
				t, ok := g.Targets[ct.id]
				if ok && t.HP > 0 {
					r := int(ct.mult * float64(p.BetLevel) * HumpbackWhaleWaveMult)
					if r < 1 {
						r = 1
					}
					delete(g.Targets, ct.id)
					waveKills++
					waveReward += r
					totalKills++
					totalReward += r
					if pp, ok2 := g.Players[p.ID]; ok2 {
						pp.Coins += r
					}
				}
				g.mu.Unlock()
			}
		}

		// 廣播本波結果
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgHumpbackWhale,
			Payload: ws.HumpbackWhalePayload{
				Event:      "wave_attack",
				WaveNum:    wave,
				WaveKills:  waveKills,
				WaveReward: waveReward,
				TotalKills: totalKills,
			},
		})
	}

	// 深海巨浪（5% 機率）
	tidalKills := 0
	tidalReward := 0
	if hasTidal {
		time.Sleep(500 * time.Millisecond) // 短暫停頓，讓玩家感受到「巨浪即將來臨」

		// 廣播深海巨浪開始
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgHumpbackWhale,
			Payload: ws.HumpbackWhalePayload{
				Event: "tidal_wave_start",
			},
		})

		// 全服公告
		ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, 0, map[string]string{
			"message": fmt.Sprintf("🌊 深海巨浪！%s 的座頭鯨掀起滔天巨浪！", p.DisplayName),
		})
		g.broadcastAnnouncement(ann)

		// 對全場所有目標執行深海巨浪
		g.mu.RLock()
		type tidalTarget struct {
			id   string
			mult float64
		}
		var allTargets []tidalTarget
		for _, t := range g.Targets {
			if t.HP <= 0 || t.DefID == "B001" || isGhostFishClone(t.DefID) {
				continue
			}
			allTargets = append(allTargets, tidalTarget{id: t.InstanceID, mult: t.Multiplier})
		}
		g.mu.RUnlock()

		for _, tt := range allTargets {
			time.Sleep(HumpbackWhaleTidalIntervalMs * time.Millisecond)
			if rand.Float64() < HumpbackWhaleTidalKillChance {
				g.mu.Lock()
				t, ok := g.Targets[tt.id]
				if ok && t.HP > 0 {
					r := int(tt.mult * float64(p.BetLevel) * HumpbackWhaleTidalMult)
					if r < 1 {
						r = 1
					}
					delete(g.Targets, tt.id)
					tidalKills++
					tidalReward += r
					totalKills++
					totalReward += r
					if pp, ok2 := g.Players[p.ID]; ok2 {
						pp.Coins += r
					}
				}
				g.mu.Unlock()
			}
		}

		// 廣播深海巨浪結果
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgHumpbackWhale,
			Payload: ws.HumpbackWhalePayload{
				Event:       "tidal_wave_result",
				TidalKills:  tidalKills,
				TidalReward: tidalReward,
			},
		})

		// 全服公告（≥3 個擊破才公告）
		if tidalKills >= 3 {
			ann2 := g.Announce.Create(announce.EventMegaWin, p.DisplayName, tidalReward, map[string]string{
				"message": fmt.Sprintf("🌊💥 深海巨浪！擊破 %d 個目標！", tidalKills),
			})
			g.broadcastAnnouncement(ann2)
		}
	}

	// 最終結算廣播
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgHumpbackWhale,
		Payload: ws.HumpbackWhalePayload{
			Event:       "awaken_result",
			KillerName:  p.DisplayName,
			BaseReward:  baseReward,
			TotalKills:  totalKills,
			TotalReward: totalReward,
			HasTidal:    hasTidal,
		},
	})

	// 全服公告（≥100x 才公告）
	if totalReward >= p.BetLevel*100 {
		ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, totalReward, map[string]string{
			"message": fmt.Sprintf("🐋✨ 座頭鯨覺醒！%s 獲得 %d 金幣！", p.DisplayName, totalReward),
		})
		g.broadcastAnnouncement(ann)
	}

	// 重置管理器
	mgr := g.HumpbackWhale
	mgr.mu.Lock()
	mgr.isActive = false
	mgr.cooldownEnd = time.Now().Add(HumpbackWhaleCooldownSec * time.Second)
	mgr.mu.Unlock()

	log.Printf("[HumpbackWhale] complete: waveKills=%d tidalKills=%d totalReward=%d hasTidal=%v",
		totalKills-tidalKills, tidalKills, totalReward, hasTidal)
}
