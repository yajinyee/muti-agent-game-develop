// thunderbolt_lobster_handler.go — 雷霆龍蝦免費射擊系統 handler（DAY-150）
// 業界依據：royalfishingsite.com 2026「Thunderbolt Lobster feature —
// 15 seconds of free play followed by automatic shooting」
// 設計：T114 雷霆龍蝦擊破後觸發「免費射擊模式」
// 15 秒內所有子彈不扣費，Server 自動幫玩家以當前 betLevel 射擊
// 每 0.5 秒自動射擊一次，優先選高倍率目標
// 全服廣播：讓其他玩家看到「有人觸發了雷霆龍蝦免費射擊」
package game

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"

	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/combat"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

const (
	ThunderboltLobsterDefID       = "T114"
	ThunderboltLobsterDuration    = 15 * time.Second  // 免費射擊持續時間
	ThunderboltLobsterShotInterval = 500 * time.Millisecond // 自動射擊間隔
	ThunderboltLobsterCooldown    = 60 * time.Second  // 每個玩家的冷卻時間
	ThunderboltLobsterMaxShots    = 30                // 最多射擊次數（15s / 0.5s = 30）
)

// thunderboltLobsterSession 免費射擊 session
type thunderboltLobsterSession struct {
	PlayerID  string
	StartAt   time.Time
	EndAt     time.Time
	ShotCount int
	KillCount int
	TotalReward int
}

// thunderboltLobsterManager 管理所有玩家的免費射擊 session
type thunderboltLobsterManager struct {
	mu       sync.Mutex
	sessions map[string]*thunderboltLobsterSession // playerID → session
	cooldowns map[string]time.Time                 // playerID → 冷卻結束時間
}

func newThunderboltLobsterManager() *thunderboltLobsterManager {
	return &thunderboltLobsterManager{
		sessions:  make(map[string]*thunderboltLobsterSession),
		cooldowns: make(map[string]time.Time),
	}
}

// CanTrigger 判斷玩家是否可以觸發免費射擊
func (m *thunderboltLobsterManager) CanTrigger(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, active := m.sessions[playerID]; active {
		return false
	}
	if cd, ok := m.cooldowns[playerID]; ok && time.Now().Before(cd) {
		return false
	}
	return true
}

// StartSession 開始免費射擊 session
func (m *thunderboltLobsterManager) StartSession(playerID string) *thunderboltLobsterSession {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess := &thunderboltLobsterSession{
		PlayerID: playerID,
		StartAt:  time.Now(),
		EndAt:    time.Now().Add(ThunderboltLobsterDuration),
	}
	m.sessions[playerID] = sess
	return sess
}

// RecordShot 記錄一次自動射擊
func (m *thunderboltLobsterManager) RecordShot(playerID string, isKill bool, reward int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok {
		return
	}
	sess.ShotCount++
	if isKill {
		sess.KillCount++
		sess.TotalReward += reward
	}
}

// EndSession 結束 session，設定冷卻
func (m *thunderboltLobsterManager) EndSession(playerID string) *thunderboltLobsterSession {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[playerID]
	if !ok {
		return nil
	}
	delete(m.sessions, playerID)
	m.cooldowns[playerID] = time.Now().Add(ThunderboltLobsterCooldown)
	return sess
}

// IsActive 判斷玩家是否在免費射擊模式中
func (m *thunderboltLobsterManager) IsActive(playerID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.sessions[playerID]
	return ok
}

// RemovePlayer 玩家離線時清理
func (m *thunderboltLobsterManager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, playerID)
}

// isThunderboltLobster 判斷是否為雷霆龍蝦目標
func isThunderboltLobster(defID string) bool {
	return defID == ThunderboltLobsterDefID
}

// tryThunderboltLobster 擊破雷霆龍蝦後觸發免費射擊模式（由 handleKill 呼叫）
func (g *Game) tryThunderboltLobster(p *player.Player, killedInstanceID string, killedX, killedY float64) {
	if !g.ThunderboltLobster.CanTrigger(p.ID) {
		return
	}

	sess := g.ThunderboltLobster.StartSession(p.ID)
	log.Printf("[ThunderboltLobster] player=%s triggered free shooting (15s, 30 shots)", p.ID)

	// 廣播免費射擊開始（全服可見）
	activateMsg := &ws.Message{
		Type: ws.MsgThunderboltLobsterActivate,
		Payload: ws.ThunderboltLobsterActivatePayload{
			TriggerID:    killedInstanceID,
			TriggerX:     killedX,
			TriggerY:     killedY,
			KillerID:     p.ID,
			KillerName:   p.DisplayName,
			Duration:     int(ThunderboltLobsterDuration.Seconds()),
			ShotInterval: int(ThunderboltLobsterShotInterval.Milliseconds()),
			Message:      fmt.Sprintf("⚡ %s 觸發了雷霆龍蝦！免費射擊 15 秒！", p.DisplayName),
		},
	}
	g.Hub.Broadcast(activateMsg)

	// 全服公告
	ann := g.Announce.Create(announce.EventThunderboltLobster, p.DisplayName, ThunderboltLobsterMaxShots, nil)
	g.broadcastAnnouncement(ann)

	// 啟動免費射擊 goroutine
	go g.runThunderboltLobsterShots(p, sess)
}

// runThunderboltLobsterShots 執行免費射擊循環
func (g *Game) runThunderboltLobsterShots(p *player.Player, sess *thunderboltLobsterSession) {
	ticker := time.NewTicker(ThunderboltLobsterShotInterval)
	defer ticker.Stop()

	shotIndex := 0
	for {
		select {
		case <-ticker.C:
			if time.Now().After(sess.EndAt) || shotIndex >= ThunderboltLobsterMaxShots {
				// 時間到或達到最大射擊次數，結束
				g.endThunderboltLobsterSession(p)
				return
			}

			// 選擇最佳目標（優先高倍率）
			target := g.selectThunderboltTarget()
			if target == nil {
				// 沒有目標，繼續等待
				shotIndex++
				continue
			}

			// 執行免費射擊（不扣費）
			req := combat.AttackRequest{
				PlayerID:       p.ID,
				TargetID:       target.InstanceID,
				BetLevel:       p.BetLevel,
				IsAuto:         true,
				IsLock:         true,
				WeaponPowerMod: p.GetWeaponPowerMod(),
				EventKillAdd:   g.getEventKillChanceAdd(),
			}

			g.mu.RLock()
			t := g.Targets[target.InstanceID]
			g.mu.RUnlock()

			if t == nil {
				shotIndex++
				continue
			}

			result := combat.ProcessAttack(req, t)

			isKill := false
			reward := 0
			if result.IsKill {
				g.mu.Lock()
				tCheck, exists := g.Targets[t.InstanceID]
				if exists && tCheck.IsAlive {
					isKill = true
					reward = result.Reward
					tCheck.IsAlive = false
					tCheck.HP = 0
					delete(g.Targets, t.InstanceID)
					g.mu.Unlock()
					// 發放獎勵（免費射擊，不扣費）
					p.AddCoins(reward)
					// 廣播目標消失
					g.Hub.Broadcast(&ws.Message{
						Type: ws.MsgTargetKill,
						Payload: ws.TargetKillPayload{
							InstanceID: t.InstanceID,
							DefID:      t.DefID,
							Multiplier: float64(t.Multiplier),
							Reward:     reward,
							LaborGain:  0,
							KillerID:   p.ID,
							Quality:    "normal",
						},
					})
				} else {
					g.mu.Unlock()
				}
			}

			// 記錄射擊
			g.ThunderboltLobster.RecordShot(p.ID, isKill, reward)
			shotsLeft := ThunderboltLobsterMaxShots - shotIndex - 1

			// 廣播自動射擊（全服可見）
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgThunderboltLobsterShot,
				Payload: ws.ThunderboltLobsterShotPayload{
					KillerID:    p.ID,
					KillerName:  p.DisplayName,
					TargetID:    t.InstanceID,
					TargetDefID: t.DefID,
					TargetName:  t.Def.Name,
					TargetX:     t.X,
					TargetY:     t.Y,
					IsKill:      isKill,
					Reward:      reward,
					Multiplier:  float64(t.Multiplier),
					ShotIndex:   shotIndex,
					ShotsLeft:   shotsLeft,
				},
			})

			// 更新玩家狀態（每 5 次更新一次，避免過於頻繁）
			if shotIndex%5 == 0 {
				g.sendPlayerUpdate(p)
			}

			shotIndex++

		case <-g.stopCh:
			return
		}
	}
}

// selectThunderboltTarget 選擇最佳射擊目標（優先高倍率，其次快要離開畫面的）
func (g *Game) selectThunderboltTarget() *thunderboltTargetCandidate {
	g.mu.RLock()
	defer g.mu.RUnlock()

	type candidate struct {
		instanceID string
		defID      string
		name       string
		multiplier float64
		x, y       float64
		score      float64
	}

	candidates := make([]candidate, 0, len(g.Targets))
	for _, t := range g.Targets {
		if t.HP <= 0 {
			continue
		}
		score := t.Multiplier * 2.0
		// 快要離開畫面的目標加分
		if t.X < 200 {
			score += 10.0
		}
		candidates = append(candidates, candidate{
			instanceID: t.InstanceID,
			defID:      t.DefID,
			name:       t.Def.Name,
			multiplier: float64(t.Multiplier),
			x:          t.X,
			y:          t.Y,
			score:      score,
		})
	}

	if len(candidates) == 0 {
		return nil
	}

	// 按分數排序，選最高分
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	// 從前 3 名中隨機選一個（避免總是打同一個目標）
	topN := 3
	if len(candidates) < topN {
		topN = len(candidates)
	}
	chosen := candidates[rand.Intn(topN)]

	return &thunderboltTargetCandidate{
		InstanceID: chosen.instanceID,
		DefID:      chosen.defID,
		Name:       chosen.name,
		Multiplier: chosen.multiplier,
		X:          chosen.x,
		Y:          chosen.y,
	}
}

// thunderboltTargetCandidate 目標候選
type thunderboltTargetCandidate struct {
	InstanceID string
	DefID      string
	Name       string
	Multiplier float64
	X, Y       float64
}

// endThunderboltLobsterSession 結束免費射擊 session
func (g *Game) endThunderboltLobsterSession(p *player.Player) {
	sess := g.ThunderboltLobster.EndSession(p.ID)
	if sess == nil {
		return
	}

	log.Printf("[ThunderboltLobster] player=%s free shooting ended: shots=%d kills=%d reward=%d",
		p.ID, sess.ShotCount, sess.KillCount, sess.TotalReward)

	// 最終更新玩家狀態
	g.sendPlayerUpdate(p)

	// 廣播免費射擊結束（全服可見）
	endMsg := &ws.Message{
		Type: ws.MsgThunderboltLobsterEnd,
		Payload: ws.ThunderboltLobsterEndPayload{
			KillerID:    p.ID,
			KillerName:  p.DisplayName,
			TotalShots:  sess.ShotCount,
			TotalKills:  sess.KillCount,
			TotalReward: sess.TotalReward,
			NewBalance:  p.Coins,
			Message:     fmt.Sprintf("⚡ %s 的雷霆龍蝦免費射擊結束！擊破 %d 個目標，獲得 %d 金幣！", p.DisplayName, sess.KillCount, sess.TotalReward),
		},
	}
	g.Hub.Broadcast(endMsg)

	// 全服公告：擊破 ≥5 個目標時廣播
	if sess.KillCount >= 5 {
		g.announceThunderboltLobsterResult(p.DisplayName, sess.KillCount, sess.TotalReward)
	}
}

// announceThunderboltLobsterResult 全服公告雷霆龍蝦結果
func (g *Game) announceThunderboltLobsterResult(playerName string, kills, reward int) {
	ann := g.Announce.Create(announce.EventThunderboltLobsterResult, playerName, kills, map[string]string{
		"reward": fmt.Sprintf("%d", reward),
	})
	g.broadcastAnnouncement(ann)
}
