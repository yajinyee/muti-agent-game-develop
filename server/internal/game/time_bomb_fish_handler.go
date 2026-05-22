// time_bomb_fish_handler.go — 時間炸彈魚系統 handler（DAY-189）
// 業界靈感：Ocean King 系列炸彈魚概念 + 倒數計時緊張感設計
// T147 時間炸彈魚出現後，螢幕顯示 10 秒倒數計時：
//   - 倒數結束前玩家擊破 → 「拆彈成功」：全服 +25% 加成持續 15 秒
//   - 倒數結束無人擊破 → 「炸彈爆炸」：全場目標 80% 擊破機率（0.5x 倍率）
// 設計差異：
//   - 與連鎖爆炸魚（被動觸發）不同，時間炸彈魚是「主動倒數」，製造「搶時間」的緊張感
//   - 與鳳凰魚（擊破後爆炸）不同，時間炸彈魚是「不擊破才爆炸」，玩家需要主動介入
//   - 「拆彈成功」的加成獎勵讓玩家有「英雄感」，「炸彈爆炸」讓玩家有「清場爽感」
//   - 兩種結果都有獎勵，但方式不同，製造「要不要拆彈」的策略決策
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

const (
	timeBombCountdown      = 10 * time.Second // 倒數時間
	timeBombCooldown       = 40 * time.Second // 全服冷卻
	timeBombKillChance     = 0.80             // 爆炸擊破機率（80%）
	timeBombRewardMult     = 0.50             // 爆炸獎勵倍率（0.5x）
	timeBombDefuseBonus    = 0.25             // 拆彈成功加成（+25%）
	timeBombDefuseDuration = 15 * time.Second // 拆彈加成持續時間
)

// timeBombManager 時間炸彈魚系統管理器（全服共享）
type timeBombManager struct {
	mu           sync.Mutex
	isActive     bool
	instanceID   string
	defuseEnd    time.Time // 拆彈加成結束時間
	lastCooldown time.Time
	stopTimer    chan struct{}
}

func newTimeBombManager() *timeBombManager {
	return &timeBombManager{
		stopTimer: make(chan struct{}, 1),
	}
}

// isTimeBombFish 判斷是否為時間炸彈魚
func isTimeBombFish(defID string) bool {
	return defID == "T147"
}

// canTrigger 是否可以觸發（冷卻檢查）
func (m *timeBombManager) canTrigger() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isActive {
		return false
	}
	return time.Since(m.lastCooldown) >= timeBombCooldown
}

// startCountdown 開始倒數
func (m *timeBombManager) startCountdown(instanceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isActive = true
	m.instanceID = instanceID
	m.lastCooldown = time.Now()
	// 重置 stop channel
	select {
	case <-m.stopTimer:
	default:
	}
}

// defuse 拆彈成功，回傳是否成功（防止重複觸發）
func (m *timeBombManager) defuse() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.isActive {
		return false
	}
	m.isActive = false
	m.defuseEnd = time.Now().Add(timeBombDefuseDuration)
	// 發送停止訊號
	select {
	case m.stopTimer <- struct{}{}:
	default:
	}
	return true
}

// explode 炸彈爆炸，回傳是否成功（防止重複觸發）
func (m *timeBombManager) explode() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.isActive {
		return false
	}
	m.isActive = false
	return true
}

// getDefuseBoost 取得拆彈加成（供 handleKill 使用）
func (m *timeBombManager) getDefuseBoost() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if time.Now().Before(m.defuseEnd) {
		return timeBombDefuseBonus
	}
	return 0.0
}

// isDefuseActive 是否在拆彈加成期間
func (m *timeBombManager) isDefuseActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return time.Now().Before(m.defuseEnd)
}

// tryTimeBombFishSpawn T147 生成時觸發倒數（由 spawnTarget 呼叫）
func (g *Game) tryTimeBombFishSpawn(instanceID string) {
	if g.TimeBomb == nil {
		return
	}
	if !g.TimeBomb.canTrigger() {
		return
	}

	g.TimeBomb.startCountdown(instanceID)

	log.Printf("[TimeBomb] spawned: instance=%s, countdown=%ds", instanceID, int(timeBombCountdown.Seconds()))

	// 全服廣播：時間炸彈魚出現（倒數開始）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTimeBombFish,
		Payload: ws.TimeBombFishPayload{
			Phase:      "bomb_appear",
			InstanceID: instanceID,
			Countdown:  int(timeBombCountdown.Seconds()),
			Message:    fmt.Sprintf("💣 時間炸彈魚出現！%d 秒內擊破可拆彈！否則全場爆炸！", int(timeBombCountdown.Seconds())),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventBossWarning, "時間炸彈魚", 0, map[string]string{
		"message": fmt.Sprintf("💣 時間炸彈魚降臨！%d 秒倒數！快去擊破拆彈！", int(timeBombCountdown.Seconds())),
	})
	g.broadcastAnnouncement(ann)

	// 啟動倒數 goroutine
	go g.runTimeBombCountdown(instanceID)
}

// runTimeBombCountdown 執行倒數計時（goroutine）
func (g *Game) runTimeBombCountdown(instanceID string) {
	// 每秒廣播倒數更新
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	remaining := int(timeBombCountdown.Seconds())

	for {
		select {
		case <-ticker.C:
			remaining--
			if remaining <= 0 {
				// 倒數結束，觸發爆炸
				g.onTimeBombExplode(instanceID)
				return
			}
			// 廣播倒數更新（每秒）
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgTimeBombFish,
				Payload: ws.TimeBombFishPayload{
					Phase:      "bomb_tick",
					InstanceID: instanceID,
					Countdown:  remaining,
				},
			})

		case <-g.TimeBomb.stopTimer:
			log.Printf("[TimeBomb] countdown stopped (defused): instance=%s", instanceID)
			return

		case <-g.stopCh:
			return
		}
	}
}

// notifyTimeBombDefuse 玩家擊破時間炸彈魚（拆彈成功）（由 handleKill 呼叫）
func (g *Game) notifyTimeBombDefuse(p *player.Player, instanceID string, baseMult float64) {
	if g.TimeBomb == nil {
		return
	}

	if !g.TimeBomb.defuse() {
		return // 已爆炸或已被拆彈
	}

	// 基礎擊破獎勵
	baseReward := int(baseMult * float64(p.BetLevel))
	p.AddCoins(baseReward)

	log.Printf("[TimeBomb] defused by player=%s, base_reward=%d, defuse_bonus=+%.0f%%",
		p.ID, baseReward, timeBombDefuseBonus*100)

	// 全服廣播：拆彈成功
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTimeBombFish,
		Payload: ws.TimeBombFishPayload{
			Phase:        "bomb_defused",
			InstanceID:   instanceID,
			KillerID:     p.ID,
			KillerName:   p.DisplayName,
			BaseReward:   baseReward,
			NewBalance:   p.GetCoins(),
			BonusPct:     int(timeBombDefuseBonus * 100),
			BonusDuration: int(timeBombDefuseDuration.Seconds()),
			Message:      fmt.Sprintf("💚 %s 拆彈成功！全服 +%d%% 加成持續 %d 秒！", p.DisplayName, int(timeBombDefuseBonus*100), int(timeBombDefuseDuration.Seconds())),
		},
	})

	// 全服公告
	ann := g.Announce.Create(announce.EventMegaWin, p.DisplayName, baseReward, map[string]string{
		"message": fmt.Sprintf("💚 %s 拆彈成功！全服獲得 +%d%% 加成！", p.DisplayName, int(timeBombDefuseBonus*100)),
	})
	g.broadcastAnnouncement(ann)

	// 15 秒後廣播加成結束
	go func() {
		time.Sleep(timeBombDefuseDuration)
		if !g.TimeBomb.isDefuseActive() {
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgTimeBombFish,
				Payload: ws.TimeBombFishPayload{
					Phase:      "defuse_end",
					InstanceID: instanceID,
					Message:    "💚 拆彈加成結束",
				},
			})
		}
	}()
}

// onTimeBombExplode 炸彈爆炸（倒數結束）
func (g *Game) onTimeBombExplode(instanceID string) {
	if g.TimeBomb == nil {
		return
	}

	if !g.TimeBomb.explode() {
		return // 已被拆彈
	}

	log.Printf("[TimeBomb] EXPLODED: instance=%s", instanceID)

	// 廣播爆炸開始
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTimeBombFish,
		Payload: ws.TimeBombFishPayload{
			Phase:      "bomb_explode",
			InstanceID: instanceID,
			Message:    "💥 時間炸彈爆炸！全場目標受到爆炸傷害！",
		},
	})

	// 等待爆炸動畫（500ms）
	time.Sleep(500 * time.Millisecond)

	// 對全場目標造成傷害
	g.mu.Lock()
	type explodeEntry struct {
		instanceID string
		defID      string
		mult       float64
		killed     bool
		reward     int
	}
	var entries []explodeEntry
	totalReward := 0
	killCount := 0

	// 計算全服平均 betLevel
	avgBetLevel := 1
	if len(g.Players) > 0 {
		total := 0
		for _, p := range g.Players {
			total += p.BetLevel
		}
		avgBetLevel = total / len(g.Players)
		if avgBetLevel < 1 {
			avgBetLevel = 1
		}
	}

	for id, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		if t.DefID == "B001" || t.DefID == "T144" || t.DefID == "T147" {
			continue // 不炸 BOSS、龍龜、自己
		}

		entry := explodeEntry{
			instanceID: id,
			defID:      t.DefID,
			mult:       t.Multiplier,
		}

		if rand.Float64() < timeBombKillChance {
			reward := int(t.Multiplier * float64(avgBetLevel) * timeBombRewardMult)
			if reward < 1 {
				reward = 1
			}
			entry.killed = true
			entry.reward = reward
			totalReward += reward
			killCount++
			t.HP = 0
			delete(g.Targets, id)

			// 廣播目標被炸死
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgTargetKill,
				Payload: ws.TargetKillPayload{
					InstanceID: id,
					KillerID:   "time_bomb",
					Reward:     0, // 爆炸獎勵統一在結果廣播
					Multiplier: t.Multiplier,
				},
			})
		}
		entries = append(entries, entry)
	}

	// 把獎勵分配給所有在線玩家（平均分配）
	if totalReward > 0 && len(g.Players) > 0 {
		rewardPerPlayer := totalReward / len(g.Players)
		if rewardPerPlayer < 1 {
			rewardPerPlayer = 1
		}
		for _, p := range g.Players {
			p.AddCoins(rewardPerPlayer)
		}
	}
	g.mu.Unlock()

	log.Printf("[TimeBomb] exploded: killed=%d, total_reward=%d", killCount, totalReward)

	// 廣播爆炸結果
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTimeBombFish,
		Payload: ws.TimeBombFishPayload{
			Phase:       "bomb_result",
			InstanceID:  instanceID,
			KillCount:   killCount,
			TotalReward: totalReward,
			Message:     fmt.Sprintf("💥 炸彈爆炸！擊破 %d 個目標，全服共享 %d 金幣！", killCount, totalReward),
		},
	})

	// 全服公告（≥5 個擊破）
	if killCount >= 5 {
		ann := g.Announce.Create(announce.EventMegaWin, "時間炸彈", totalReward, map[string]string{
			"message": fmt.Sprintf("💥 時間炸彈爆炸！擊破 %d 個目標！全服共享 %d 金幣！", killCount, totalReward),
		})
		g.broadcastAnnouncement(ann)
	}
}

// getTimeBombDefuseBoost 取得拆彈加成（供 handleKill 使用）
func (g *Game) getTimeBombDefuseBoost() float64 {
	if g.TimeBomb == nil {
		return 0.0
	}
	return g.TimeBomb.getDefuseBoost()
}
