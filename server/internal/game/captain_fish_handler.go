// captain_fish_handler.go — 船長魚全服競速模式 handler（DAY-163）
// 業界依據：King of Ocean 2026（Galaxsys）「Captain Fish trigger bonus rounds」
// 擊破船長魚後觸發「全服競速模式」，30 秒內全服玩家競爭擊破最多目標
// 第一名獲得額外大獎（betLevel × 30），第二名 × 15，第三名 × 8
// 設計：競技型社交機制（與水晶龍的合作型形成對比），製造「全服競爭」的緊張爽感
package game

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// captainFishManager 船長魚競速模式管理器（全服共享）
type captainFishManager struct {
	mu          sync.Mutex
	isActive    bool
	startAt     time.Time
	duration    float64
	triggerID   string
	triggerName string
	scores      map[string]*captainRaceScore // playerID -> score
	cooldownEnd time.Time
}

type captainRaceScore struct {
	PlayerID    string
	PlayerName  string
	KillCount   int
	TotalReward int
	BetLevel    int
}

func newCaptainFishManager() *captainFishManager {
	return &captainFishManager{
		scores: make(map[string]*captainRaceScore),
	}
}

// isCaptainFish 判斷是否為船長魚
func isCaptainFish(defID string) bool {
	return defID == "T123"
}

// IsCaptainRaceActive 查詢競速模式是否活躍（供其他系統使用）
func (g *Game) IsCaptainRaceActive() bool {
	if g.CaptainFish == nil {
		return false
	}
	g.CaptainFish.mu.Lock()
	defer g.CaptainFish.mu.Unlock()
	if !g.CaptainFish.isActive {
		return false
	}
	return time.Since(g.CaptainFish.startAt).Seconds() < g.CaptainFish.duration
}

// recordCaptainRaceKill 記錄競速模式中的擊破（由 handleKill 呼叫）
func (g *Game) recordCaptainRaceKill(p *player.Player, reward int) {
	if g.CaptainFish == nil {
		return
	}
	g.CaptainFish.mu.Lock()
	defer g.CaptainFish.mu.Unlock()

	if !g.CaptainFish.isActive {
		return
	}
	if time.Since(g.CaptainFish.startAt).Seconds() >= g.CaptainFish.duration {
		return
	}

	score, ok := g.CaptainFish.scores[p.ID]
	if !ok {
		score = &captainRaceScore{
			PlayerID:   p.ID,
			PlayerName: p.DisplayName,
			BetLevel:   p.BetLevel,
		}
		g.CaptainFish.scores[p.ID] = score
	}
	score.KillCount++
	score.TotalReward += reward

	// 廣播即時排名更新（每 5 次擊破廣播一次，避免過於頻繁）
	if score.KillCount%5 == 0 {
		go g.broadcastCaptainRaceUpdate()
	}
}

// broadcastCaptainRaceUpdate 廣播競速排名更新
func (g *Game) broadcastCaptainRaceUpdate() {
	if g.CaptainFish == nil {
		return
	}
	g.CaptainFish.mu.Lock()
	entries := g.buildCaptainRaceEntries()
	remaining := g.CaptainFish.duration - time.Since(g.CaptainFish.startAt).Seconds()
	if remaining < 0 {
		remaining = 0
	}
	g.CaptainFish.mu.Unlock()

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCaptainFishRace,
		Payload: ws.CaptainFishRacePayload{
			Phase:         "race_update",
			RemainingTime: remaining,
			Entries:       entries,
		},
	})
}

// buildCaptainRaceEntries 建立排名列表（需在 mu.Lock 下呼叫）
func (g *Game) buildCaptainRaceEntries() []ws.CaptainRaceEntry {
	entries := make([]ws.CaptainRaceEntry, 0, len(g.CaptainFish.scores))
	for _, s := range g.CaptainFish.scores {
		entries = append(entries, ws.CaptainRaceEntry{
			PlayerID:    s.PlayerID,
			PlayerName:  s.PlayerName,
			KillCount:   s.KillCount,
			TotalReward: s.TotalReward,
		})
	}
	// 按擊破數排序（降序）
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].KillCount != entries[j].KillCount {
			return entries[i].KillCount > entries[j].KillCount
		}
		return entries[i].TotalReward > entries[j].TotalReward
	})
	// 加入排名
	for i := range entries {
		entries[i].Rank = i + 1
	}
	return entries
}

// tryCaptainFishRace 嘗試觸發船長魚競速模式（擊破 T123 後呼叫）
func (g *Game) tryCaptainFishRace(p *player.Player, instanceID string, x, y float64) {
	if g.CaptainFish == nil {
		return
	}

	const duration = 30.0 // 30 秒競速
	const cooldown = 120.0 // 120 秒冷卻（競速模式影響大，冷卻更長）

	g.CaptainFish.mu.Lock()
	if g.CaptainFish.isActive {
		g.CaptainFish.mu.Unlock()
		return // 已在競速模式中
	}
	if time.Now().Before(g.CaptainFish.cooldownEnd) {
		g.CaptainFish.mu.Unlock()
		return // 冷卻中
	}

	// 啟動競速模式
	g.CaptainFish.isActive = true
	g.CaptainFish.startAt = time.Now()
	g.CaptainFish.duration = duration
	g.CaptainFish.triggerID = p.ID
	g.CaptainFish.triggerName = p.DisplayName
	g.CaptainFish.scores = make(map[string]*captainRaceScore) // 清空舊分數
	g.CaptainFish.mu.Unlock()

	log.Printf("[CaptainFish] player=%s triggered race mode for %.0fs (global)", p.ID, duration)

	// 廣播競速開始（全服）
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgCaptainFishRace,
		Payload: ws.CaptainFishRacePayload{
			TriggerID:     instanceID,
			TriggerX:      x,
			TriggerY:      y,
			KillerID:      p.ID,
			KillerName:    p.DisplayName,
			Phase:         "race_start",
			DurationSecs:  duration,
			RemainingTime: duration,
			Entries:       []ws.CaptainRaceEntry{},
		},
	})

	// 全服公告
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgAnnouncement,
		Payload: map[string]interface{}{
			"event_type": "captain_fish_race",
			"message":    fmt.Sprintf("⚓ %s 擊破船長魚！全服競速開始！30 秒內擊破最多目標獲得大獎！", p.DisplayName),
			"color":      "#4488FF",
			"duration":   5.0,
			"priority":   5,
		},
	})

	// 等待競速結束
	go func() {
		time.Sleep(time.Duration(duration * float64(time.Second)))

		g.CaptainFish.mu.Lock()
		g.CaptainFish.isActive = false
		g.CaptainFish.cooldownEnd = time.Now().Add(time.Duration(cooldown * float64(time.Second)))
		entries := g.buildCaptainRaceEntries()
		g.CaptainFish.mu.Unlock()

		// 發放獎勵
		g.distributeCaptainRaceRewards(entries)

		// 廣播競速結束
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgCaptainFishRace,
			Payload: ws.CaptainFishRacePayload{
				KillerID:   p.ID,
				KillerName: p.DisplayName,
				Phase:      "race_end",
				Entries:    entries,
			},
		})

		log.Printf("[CaptainFish] race ended, %d participants, cooldown %.0fs", len(entries), cooldown)
	}()
}

// distributeCaptainRaceRewards 發放競速獎勵
func (g *Game) distributeCaptainRaceRewards(entries []ws.CaptainRaceEntry) {
	// 獎勵倍率：第1名 ×30，第2名 ×15，第3名 ×8
	rewardMults := []int{30, 15, 8}

	g.mu.RLock()
	players := g.Players
	g.mu.RUnlock()

	for i, entry := range entries {
		if i >= len(rewardMults) {
			break
		}
		p, ok := players[entry.PlayerID]
		if !ok {
			continue
		}
		bonus := p.BetLevel * rewardMults[i]
		p.AddReward(bonus)

		log.Printf("[CaptainFish] rank=%d player=%s bonus=%d (betLevel=%d × %d)",
			i+1, p.ID, bonus, p.BetLevel, rewardMults[i])

		// 個人通知
		_ = g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgCaptainFishRace,
			Payload: ws.CaptainFishRacePayload{
				Phase:       "race_reward",
				KillerID:    p.ID,
				KillerName:  p.DisplayName,
				Entries:     entries,
				MyRank:      i + 1,
				MyBonus:     bonus,
				MyKillCount: entry.KillCount,
			},
		})

		// 第一名全服公告
		if i == 0 && entry.KillCount > 0 {
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgAnnouncement,
				Payload: map[string]interface{}{
					"event_type": "captain_fish_winner",
					"message":    fmt.Sprintf("🏆 %s 競速第一！擊破 %d 個目標，獲得 %d 金幣獎勵！", p.DisplayName, entry.KillCount, bonus),
					"color":      "#FFD700",
					"duration":   5.0,
					"priority":   4,
				},
			})
		}
	}
}
