// rock_skeleton_handler.go — 搖滾骷髏演唱會系統 handler（DAY-192）
// 業界依據：JILI 2026「Rock Skeleton Concert — Rock Skeleton and Super Awakening Performance,
// you can get a large bonus of up to 3,000 times」
// T150 搖滾骷髏魚機制：
//   1. 擊破 T150 後觸發「演唱會模式」（15 秒）
//   2. 每 1 秒「音符炸彈」隨機命中 2-4 個目標（70% 擊破機率，0.60x 倍率）
//   3. 第 10 秒觸發「超級覺醒高潮」：全場所有目標 HP 降低 70%，持續 5 秒
//   4. 演唱會結束後：≥10 個擊破 → 全服 +30% 加成 10 秒（安可獎勵）
//   5. 全服廣播每次音符炸彈，讓所有玩家看到「演唱會在全場轟炸」
// 設計差異：
//   - 與閃電魚（每 0.5 秒單目標，8 秒）不同，搖滾骷髏是「每 1 秒多目標（2-4 個）」，
//     節奏更有音樂感，且有「高潮覺醒」雙段式設計
//   - 與鳳凰魚（一次性全場爆炸）不同，搖滾骷髏是「持續 15 秒的演唱會」，
//     有節奏感和高潮設計，讓玩家感受到「演唱會從開場到高潮到安可」的完整體驗
//   - 「超級覺醒高潮」讓第 10-15 秒的音符炸彈效率大幅提升（HP 降低 70%），
//     製造「演唱會越到後面越爽」的正向反饋
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	// RockSkeletonDurationSec 演唱會持續時間（秒）
	RockSkeletonDurationSec = 15
	// RockSkeletonIntervalMs 每次音符炸彈間隔（ms）
	RockSkeletonIntervalMs = 1000
	// RockSkeletonNotesPerBeat 每次音符炸彈命中目標數（2-4 個）
	RockSkeletonNotesMin = 2
	RockSkeletonNotesMax = 4
	// RockSkeletonKillChance 音符炸彈擊破機率（70%）
	RockSkeletonKillChance = 0.70
	// RockSkeletonRewardMult 音符炸彈擊破獎勵倍率
	RockSkeletonRewardMult = 0.60
	// RockSkeletonAwakeningBeat 超級覺醒高潮觸發時間（第 10 秒）
	RockSkeletonAwakeningBeat = 10
	// RockSkeletonAwakeningHPReduction 超級覺醒 HP 降低比例（70%）
	RockSkeletonAwakeningHPReduction = 0.70
	// RockSkeletonEncoreMinKills 安可獎勵最低擊破數
	RockSkeletonEncoreMinKills = 10
	// RockSkeletonEncoreBonus 安可全服加成（+30%）
	RockSkeletonEncoreBonus = 0.30
	// RockSkeletonEncoreDurationSec 安可加成持續時間（秒）
	RockSkeletonEncoreDurationSec = 10
	// RockSkeletonCooldownSec 全服冷卻時間（秒）
	RockSkeletonCooldownSec = 45
	// RockSkeletonAnnounceMinKills 全服公告最低擊破數
	RockSkeletonAnnounceMinKills = 8
)

// rockSkeletonManager 搖滾骷髏演唱會管理器（全服共享）
type rockSkeletonManager struct {
	mu          sync.Mutex
	isActive    bool
	encoreEnd   time.Time // 安可加成結束時間
	cooldownEnd time.Time
}

// newRockSkeletonManager 建立搖滾骷髏演唱會管理器
func newRockSkeletonManager() *rockSkeletonManager {
	return &rockSkeletonManager{}
}

// isRockSkeleton 判斷是否為搖滾骷髏魚（T150）
func isRockSkeleton(defID string) bool {
	return defID == "T150"
}

// isOnCooldownRS 檢查是否在全服冷卻中
func (m *rockSkeletonManager) isOnCooldownRS() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return time.Now().Before(m.cooldownEnd)
}

// activateRS 激活演唱會
func (m *rockSkeletonManager) activateRS() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isActive {
		return false
	}
	m.isActive = true
	m.cooldownEnd = time.Now().Add(time.Duration(RockSkeletonCooldownSec) * time.Second)
	return true
}

// deactivateRS 結束演唱會
func (m *rockSkeletonManager) deactivateRS() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isActive = false
}

// activateEncore 激活安可加成
func (m *rockSkeletonManager) activateEncore() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.encoreEnd = time.Now().Add(time.Duration(RockSkeletonEncoreDurationSec) * time.Second)
}

// getRockSkeletonEncoreBoost 取得安可加成（供 handleKill 使用）
// 回傳加成值（0.0 = 無加成，0.30 = +30%）
func (g *Game) getRockSkeletonEncoreBoost() float64 {
	g.RockSkeleton.mu.Lock()
	defer g.RockSkeleton.mu.Unlock()
	if time.Now().Before(g.RockSkeleton.encoreEnd) {
		return RockSkeletonEncoreBonus
	}
	return 0.0
}

// tryRockSkeletonConcert 擊破 T150 後觸發搖滾骷髏演唱會（DAY-192）
// 由 handleKill 呼叫（在 goroutine 中執行）
func (g *Game) tryRockSkeletonConcert(p *player.Player, triggerID string, triggerX, triggerY float64) {
	// 全服冷卻檢查
	if g.RockSkeleton.isOnCooldownRS() {
		return
	}
	if !g.RockSkeleton.activateRS() {
		return // 已有演唱會在進行
	}
	defer g.RockSkeleton.deactivateRS()

	log.Printf("[RockSkeleton] player=%s triggered concert", p.ID)

	// 廣播演唱會開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRockSkeletonConcert,
		Payload: ws.RockSkeletonConcertPayload{
			Phase:       "concert_start",
			TriggerID:   triggerID,
			TriggerX:    triggerX,
			TriggerY:    triggerY,
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			DurationSec: RockSkeletonDurationSec,
		},
	})

	totalReward := 0
	totalKills := 0
	awakeningTriggered := false

	for beat := 1; beat <= RockSkeletonDurationSec; beat++ {
		time.Sleep(time.Duration(RockSkeletonIntervalMs) * time.Millisecond)

		// 第 10 秒：觸發超級覺醒高潮
		if beat == RockSkeletonAwakeningBeat && !awakeningTriggered {
			awakeningTriggered = true
			g.triggerRockSkeletonAwakening(p, triggerID)
		}

		// 音符炸彈：隨機選 2-4 個目標
		notesCount := RockSkeletonNotesMin + rand.Intn(RockSkeletonNotesMax-RockSkeletonNotesMin+1)

		g.mu.RLock()
		type candidate struct {
			instanceID string
			defID      string
			x, y       float64
			multiplier float64
		}
		var candidates []candidate
		for id, t := range g.Targets {
			if id == triggerID || !t.IsAlive || t.DefID == "B001" {
				continue
			}
			candidates = append(candidates, candidate{
				instanceID: t.InstanceID,
				defID:      t.DefID,
				x:          t.X,
				y:          t.Y,
				multiplier: t.Multiplier,
			})
		}
		g.mu.RUnlock()

		if len(candidates) == 0 {
			continue
		}

		// 隨機選取目標（不重複）
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})
		if notesCount > len(candidates) {
			notesCount = len(candidates)
		}
		selected := candidates[:notesCount]

		// 廣播音符炸彈（全服）
		noteTargetIDs := make([]string, 0, len(selected))
		noteTargetXs := make([]float64, 0, len(selected))
		noteTargetYs := make([]float64, 0, len(selected))
		for _, c := range selected {
			noteTargetIDs = append(noteTargetIDs, c.instanceID)
			noteTargetXs = append(noteTargetXs, c.x)
			noteTargetYs = append(noteTargetYs, c.y)
		}

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgRockSkeletonConcert,
			Payload: ws.RockSkeletonConcertPayload{
				Phase:         fmt.Sprintf("note_%d", beat),
				Beat:          beat,
				NoteTargetIDs: noteTargetIDs,
				NoteTargetXs:  noteTargetXs,
				NoteTargetYs:  noteTargetYs,
				KillerID:      p.ID,
				IsAwakening:   awakeningTriggered,
			},
		})

		// 對選中目標執行音符炸彈
		beatKills := 0
		beatReward := 0
		for _, c := range selected {
			if rand.Float64() >= RockSkeletonKillChance {
				continue // 未命中
			}

			g.mu.Lock()
			t, ok := g.Targets[c.instanceID]
			if !ok || !t.IsAlive {
				g.mu.Unlock()
				continue
			}
			reward := int(float64(p.BetLevel) * c.multiplier * RockSkeletonRewardMult)
			if reward < 1 {
				reward = 1
			}
			t.IsAlive = false
			delete(g.Targets, c.instanceID)
			g.mu.Unlock()

			totalReward += reward
			totalKills++
			beatKills++
			beatReward += reward

			// 廣播目標擊破
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgTargetKill,
				Payload: ws.TargetKillPayload{
					InstanceID: c.instanceID,
					KillerID:   p.ID,
					Reward:     reward,
					Multiplier: c.multiplier,
				},
			})

			log.Printf("[RockSkeleton] beat=%d target=%s mult=%.0f reward=%d",
				beat, c.instanceID, c.multiplier, reward)
		}

		// 廣播本拍結果（有擊破才廣播）
		if beatKills > 0 {
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgRockSkeletonConcert,
				Payload: ws.RockSkeletonConcertPayload{
					Phase:       "beat_result",
					Beat:        beat,
					BeatKills:   beatKills,
					BeatReward:  beatReward,
					TotalKills:  totalKills,
					KillerID:    p.ID,
					IsAwakening: awakeningTriggered,
				},
			})
		}
	}

	// 發放總獎勵
	if totalReward > 0 {
		p.AddCoins(totalReward)
		g.sendPlayerUpdate(p)
	}

	// 安可獎勵：≥10 個擊破 → 全服 +30% 加成 10 秒
	if totalKills >= RockSkeletonEncoreMinKills {
		g.RockSkeleton.activateEncore()

		// 廣播安可加成開始（全服）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgRockSkeletonConcert,
			Payload: ws.RockSkeletonConcertPayload{
				Phase:           "encore_start",
				TotalKills:      totalKills,
				TotalReward:     totalReward,
				KillerID:        p.ID,
				KillerName:      p.DisplayName,
				EncoreDuration:  RockSkeletonEncoreDurationSec,
				EncoreBonus:     RockSkeletonEncoreBonus,
			},
		})

		log.Printf("[RockSkeleton] encore activated! kills=%d", totalKills)

		// 10 秒後廣播安可結束
		go func() {
			time.Sleep(time.Duration(RockSkeletonEncoreDurationSec) * time.Second)
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgRockSkeletonConcert,
				Payload: ws.RockSkeletonConcertPayload{
					Phase:      "encore_end",
					KillerID:   p.ID,
					KillerName: p.DisplayName,
				},
			})
		}()
	} else {
		// 廣播演唱會結束（無安可）
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgRockSkeletonConcert,
			Payload: ws.RockSkeletonConcertPayload{
				Phase:       "concert_end",
				TotalKills:  totalKills,
				TotalReward: totalReward,
				KillerID:    p.ID,
				KillerName:  p.DisplayName,
			},
		})
	}

	// 個人結果通知
	if totalReward > 0 {
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgReward,
			Payload: ws.RewardPayload{
				Source:     "rock_skeleton_concert",
				Amount:     totalReward,
				Multiplier: float64(totalKills),
				NewBalance: p.Coins,
			},
		})
	}

	// 全服公告：擊破 ≥ 8 個
	if totalKills >= RockSkeletonAnnounceMinKills {
		g.announceRockSkeletonConcert(p.DisplayName, totalKills, totalReward, totalKills >= RockSkeletonEncoreMinKills)
	}

	log.Printf("[RockSkeleton] player=%s concert_end kills=%d total_reward=%d encore=%v",
		p.ID, totalKills, totalReward, totalKills >= RockSkeletonEncoreMinKills)
}

// triggerRockSkeletonAwakening 觸發超級覺醒高潮（第 10 秒）
// 全場所有目標 HP 降低 70%，持續 5 秒（讓後半段音符炸彈更容易擊破）
func (g *Game) triggerRockSkeletonAwakening(p *player.Player, triggerID string) {
	g.mu.Lock()
	var awakenedTargets []string
	for id, t := range g.Targets {
		if id == triggerID || !t.IsAlive || t.DefID == "B001" {
			continue
		}
		// HP 降低 70%（最少 1）
		newHP := int(float64(t.HP) * (1.0 - RockSkeletonAwakeningHPReduction))
		if newHP < 1 {
			newHP = 1
		}
		t.HP = newHP
		awakenedTargets = append(awakenedTargets, id)
	}
	g.mu.Unlock()

	log.Printf("[RockSkeleton] awakening triggered! affected=%d targets", len(awakenedTargets))

	// 廣播超級覺醒高潮（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgRockSkeletonConcert,
		Payload: ws.RockSkeletonConcertPayload{
			Phase:           "awakening",
			AwakenedTargets: awakenedTargets,
			AwakenedCount:   len(awakenedTargets),
			KillerID:        p.ID,
			KillerName:      p.DisplayName,
			Message:         "🎸 超級覺醒！全場目標 HP 降低 70%！",
		},
	})

	// 全服公告（覺醒時）
	if len(awakenedTargets) >= 3 {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgAnnouncement,
			Payload: map[string]interface{}{
				"event_type": "rock_skeleton_awakening",
				"message":    fmt.Sprintf("🎸💀 %s 的搖滾骷髏演唱會進入超級覺醒！%d 個目標 HP 降低 70%%！", p.DisplayName, len(awakenedTargets)),
				"color":      "#FF00FF", // 洋紅色（搖滾感）
				"duration":   4.0,
				"priority":   4,
			},
		})
	}
}

// announceRockSkeletonConcert 全服公告搖滾骷髏演唱會（DAY-192）
func (g *Game) announceRockSkeletonConcert(playerName string, kills, reward int, hasEncore bool) {
	icon := "🎸💀"
	msg := fmt.Sprintf("%s %s 的搖滾骷髏演唱會！擊破 %d 個目標！獲得 %d 金幣！",
		icon, playerName, kills, reward)
	color := "#FF6600" // 橙色
	if hasEncore {
		msg = fmt.Sprintf("🎸💀🎵 %s 的搖滾骷髏演唱會安可！擊破 %d 個目標！全服 +30%% 加成 10 秒！",
			playerName, kills)
		color = "#FF00FF" // 洋紅色（安可更特別）
	}
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "rock_skeleton_concert",
			"message":    msg,
			"color":      color,
			"duration":   5.0,
			"priority":   3,
		},
	})
}
