// Package game 遊戲主邏輯
package game

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"

	"digital-twin/server/internal/analytics"
	"digital-twin/server/internal/data"
	"digital-twin/server/internal/game/achievement"
	"digital-twin/server/internal/game/combat"
	"digital-twin/server/internal/game/state"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/ws"
)

// Game 遊戲實例（單一房間）
type Game struct {
	mu sync.RWMutex

	ID          string
	State       state.GameState
	Players     map[string]*player.Player
	Targets     map[string]*target.Target
	SpawnSys    *target.SpawnSystem
	Hub         *ws.Hub

	// 計時器
	lastSpawnAt        time.Time
	bossSpawnedAt      time.Time
	bonusStartedAt     time.Time
	lastBonusAt        time.Time
	lastBonusTickAt    time.Time  // Bonus tick 廣播計時（每秒一次）
	lastSpecialEventAt time.Time
	nextSpecialEventIn float64
	bossInstanceID     string
	nextBossAt         time.Time  // BOSS 自動觸發時間（規格書 28.1）
	lastLeaderboardAt  time.Time  // 排行榜廣播計時（每 10 秒一次）

	// 補償機制
	lastHighRewardAt time.Time
	bonusSpecialBonus float64

	// Bonus 狀態
	bonusScores map[string]int // playerID -> score
	bonusEntryBet map[string]int // playerID -> entry bet cost

	stopCh chan struct{}
}

// NewGame 建立新遊戲
func NewGame(id string, hub *ws.Hub) *Game {
	g := &Game{
		ID:                 id,
		State:              state.StateNormalPlay,
		Players:            make(map[string]*player.Player),
		Targets:            make(map[string]*target.Target),
		SpawnSys:           target.NewSpawnSystem(),
		Hub:                hub,
		lastSpawnAt:        time.Now(),
		lastSpecialEventAt: time.Now(),
		nextSpecialEventIn: 30,
		bonusScores:        make(map[string]int),
		bonusEntryBet:      make(map[string]int),
		// BOSS 自動觸發：遊戲開始後 3-5 分鐘（規格書 28.1）
		nextBossAt: time.Now().Add(time.Duration(180+rand.Intn(120)) * time.Second),
		stopCh:             make(chan struct{}),
	}
	return g
}

// GetState 取得目前遊戲狀態（thread-safe）
func (g *Game) GetState() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return string(g.State)
}

// Start 啟動遊戲循環
func (g *Game) Start() {
	log.Printf("[Game] %s started", g.ID)
	go g.gameLoop()
}

// Stop 停止遊戲
func (g *Game) Stop() {
	close(g.stopCh)
}

// AddPlayer 加入玩家
func (g *Game) AddPlayer(playerID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.Players[playerID]; !exists {
		g.Players[playerID] = player.NewPlayer(playerID, 10000) // 初始 10000 金幣
		log.Printf("[Game] Player %s joined game %s", playerID, g.ID)
		// 埋點：玩家加入（由 main.go 的 OnConnect 處理，這裡不重複）
	}
}

// RemovePlayer 移除玩家
func (g *Game) RemovePlayer(playerID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.Players, playerID)
	log.Printf("[Game] Player %s left game %s", playerID, g.ID)
	// 埋點：玩家離開（由 main.go 的 OnDisconnect 處理，這裡不重複）
}

// HandleMessage 處理玩家訊息
func (g *Game) HandleMessage(clientID string, msg *ws.Message) {
	g.mu.Lock()
	p, ok := g.Players[clientID]
	g.mu.Unlock()

	if !ok {
		return
	}

	switch msg.Type {
	case ws.MsgAttack:
		g.handleAttack(p, msg)
	case ws.MsgLock:
		g.handleLock(p, msg)
	case ws.MsgAutoToggle:
		g.handleAutoToggle(p, msg)
	case ws.MsgBetChange:
		g.handleBetChange(p, msg)
	case ws.MsgBonusClick:
		g.handleBonusClick(p, msg)
	case ws.MsgTriggerBoss:
		g.triggerBoss()
	case ws.MsgTriggerBonus:
		g.triggerBonusReady()
	case ws.MsgSetDisplayName:
		g.handleSetDisplayName(p, msg)
	case ws.MsgPing:
		g.Hub.Send(clientID, &ws.Message{Type: ws.MsgPong})
	}
}

// handleAttack 處理攻擊
func (g *Game) handleAttack(p *player.Player, msg *ws.Message) {
	// 狀態檢查
	g.mu.RLock()
	currentState := g.State
	g.mu.RUnlock()

	if currentState != state.StateNormalPlay &&
		currentState != state.StateSpecialTargetEvent &&
		currentState != state.StateBossBattle {
		return
	}

	// 解析 payload
	var payload ws.AttackPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	// 扣除投注
	betCost, ok := p.DeductBet()
	if !ok {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "insufficient_coins", Message: "金幣不足"},
		})
		return
	}

	// 找目標
	g.mu.Lock()
	targetID := payload.TargetID
	if targetID == "" {
		targetID = p.LockTargetID
	}
	t := g.Targets[targetID]
	g.mu.Unlock()

	// 處理攻擊
	req := combat.AttackRequest{
		PlayerID: p.ID,
		TargetID: targetID,
		BetLevel: p.BetLevel,
		IsAuto:   p.IsAuto,
		IsLock:   targetID != "",
		ClickX:   payload.ClickX,
		ClickY:   payload.ClickY,
	}

	result := combat.ProcessAttack(req, t)
	_ = betCost

	// 埋點：攻擊事件
	if tracker := analytics.Get(); tracker != nil {
		tracker.Track(analytics.EventAttack, p.ID, map[string]interface{}{
			"target_id": targetID,
			"bet_level": p.BetLevel,
			"bet_cost":  betCost,
			"is_hit":    result.IsHit,
			"is_auto":   p.IsAuto,
		})
	}

	// 傳送攻擊結果
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgAttackResult,
		Payload: ws.AttackResultPayload{
			TargetID:    result.TargetID,
			IsHit:       result.IsHit,
			IsKill:      result.IsKill,
			Damage:      result.Damage,
			Reward:      result.Reward,
			LaborGain:   result.LaborGain,
			CharacterID: result.CharacterID,
			Multiplier:  result.Multiplier,
		},
	})

	if result.IsKill && t != nil {
		g.handleKill(p, t, result)
	} else if result.IsHit && t != nil && t.DefID == "T102" && !t.IsFleeing {
		// T102 寶箱怪：受擊後加速逃跑（規格書 26.2）
		t.IsFleeing = true
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetUpdate,
			Payload: ws.TargetUpdatePayload{
				InstanceID: t.InstanceID,
				HP:         t.HP,
				MaxHP:      t.MaxHP,
				IsFleeing:  true,
			},
		})
	}

	// BOSS 階段變化
	if result.BossPhaseChanged {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgBossEvent,
			Payload: ws.BossEventPayload{
				Event:      "phase_change",
				InstanceID: t.InstanceID,
				Phase:      result.BossPhase,
				HP:         t.HP,
				MaxHP:      t.MaxHP,
			},
		})
	}

	// 更新玩家狀態
	g.sendPlayerUpdate(p)
}

// handleKill 處理目標擊破
func (g *Game) handleKill(p *player.Player, t *target.Target, result *combat.AttackResult) {
	g.mu.Lock()
	delete(g.Targets, t.InstanceID)
	g.mu.Unlock()

	// 發放獎勵
	rewardUnlocks := p.AddReward(result.Reward)
	killUnlocks := p.AddKill()

	// 埋點：擊破事件
	if tracker := analytics.Get(); tracker != nil {
		tracker.Track(analytics.EventKill, p.ID, map[string]interface{}{
			"def_id":     t.DefID,
			"target_type": string(t.Def.Type),
			"multiplier": result.Multiplier,
			"reward":     result.Reward,
			"labor_gain": result.LaborGain,
		})
		tracker.Track(analytics.EventReward, p.ID, map[string]interface{}{
			"source":     "target",
			"amount":     result.Reward,
			"multiplier": result.Multiplier,
		})
	}

	// 成就：特殊目標
	var specialUnlock *achievement.AchievementUnlock
	if t.Def.Type == data.TargetTypeSpecial {
		specialUnlock = p.TryUnlockAchievement(achievement.AchKillSpecial)
	}

	// 成就：大獎
	bigWinUnlocks := p.TryUnlockBigWin(result.Multiplier)

	// 傳送所有成就通知
	allUnlocks := append(rewardUnlocks, killUnlocks...)
	allUnlocks = append(allUnlocks, bigWinUnlocks...)
	if specialUnlock != nil {
		allUnlocks = append(allUnlocks, specialUnlock)
	}
	for _, u := range allUnlocks {
		g.sendAchievement(p.ID, u)
	}

	// 累積勞動值
	bonusTriggered := p.AddLaborValue(result.LaborGain)

	// 廣播擊破事件
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTargetKill,
		Payload: ws.TargetKillPayload{
			InstanceID: t.InstanceID,
			DefID:      t.DefID,
			Multiplier: result.Multiplier,
			Reward:     result.Reward,
			LaborGain:  result.LaborGain,
			KillerID:   p.ID,
		},
	})

	// 發放獎勵通知
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgReward,
		Payload: ws.RewardPayload{
			Source:     "target",
			Amount:     result.Reward,
			Multiplier: result.Multiplier,
			NewBalance: p.Coins,
		},
	})

	// 記錄高倍率獎勵時間（補償機制用）
	if result.Multiplier >= 20 {
		g.mu.Lock()
		g.lastHighRewardAt = time.Now()
		g.bonusSpecialBonus = 0 // 重置補償
		g.mu.Unlock()
	}

	// BOSS 擊破
	if t.Def.Type == data.TargetTypeBoss {
		g.handleBossKill(p, t, result)
		return
	}

	// 觸發 Bonus
	if bonusTriggered {
		// 成就：首次觸發 Bonus
		if u := p.TryUnlockAchievement(achievement.AchBonus); u != nil {
			g.sendAchievement(p.ID, u)
		}
		g.triggerBonusReady()
	}
}

// handleBossKill BOSS 擊破
func (g *Game) handleBossKill(p *player.Player, t *target.Target, result *combat.AttackResult) {
	g.mu.Lock()
	g.bossInstanceID = ""
	g.mu.Unlock()

	// 成就：首次擊敗 BOSS
	if u := p.TryUnlockAchievement(achievement.AchKillBoss); u != nil {
		g.sendAchievement(p.ID, u)
	}

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBossEvent,
		Payload: ws.BossEventPayload{
			Event:      "kill",
			InstanceID: t.InstanceID,
			Reward:     result.Reward,
			Multiplier: result.Multiplier,
		},
	})

	// 埋點：BOSS 擊敗
	if tracker := analytics.Get(); tracker != nil {
		tracker.Track(analytics.EventBossKill, p.ID, map[string]interface{}{
			"instance_id": t.InstanceID,
			"reward":      result.Reward,
			"multiplier":  result.Multiplier,
		})
	}

	g.transitionState(state.StateBossResult)
	g.safeAfterFunc(3*time.Second, func() {
		g.transitionState(state.StateNormalPlay)
	})
}

// handleLock 處理鎖定
func (g *Game) handleLock(p *player.Player, msg *ws.Message) {
	var payload ws.LockPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}
	p.SetLock(payload.TargetID)
	g.sendPlayerUpdate(p)
}

// handleAutoToggle 切換自動攻擊
func (g *Game) handleAutoToggle(p *player.Player, msg *ws.Message) {
	p.SetAuto(!p.IsAuto)
	// 埋點：切換自動攻擊
	if tracker := analytics.Get(); tracker != nil {
		tracker.Track(analytics.EventAutoToggle, p.ID, map[string]interface{}{
			"is_auto": p.IsAuto,
		})
	}
	g.sendPlayerUpdate(p)
}

// handleBetChange 切換投注
func (g *Game) handleBetChange(p *player.Player, msg *ws.Message) {
	g.mu.RLock()
	currentState := g.State
	g.mu.RUnlock()

	// Bonus Game 中不允許切換
	if currentState == state.StateBonusGame {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "bet_locked", Message: "Bonus Game 中無法切換投注"},
		})
		return
	}

	var payload ws.BetChangePayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		log.Printf("[Game] handleBetChange parse error: %v", err)
		return
	}

	if !p.SetBetLevel(payload.BetLevel) {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "invalid_bet_level", Message: "無效的投注等級"},
		})
		return
	}

	log.Printf("[Game] Player %s bet changed to LV%d", p.ID, payload.BetLevel)
	// 埋點：切換投注
	if tracker := analytics.Get(); tracker != nil {
		tracker.Track(analytics.EventBetChange, p.ID, map[string]interface{}{
			"bet_level": payload.BetLevel,
			"bet_cost":  data.GetBetDef(payload.BetLevel).BetCost,
		})
	}
	g.sendPlayerUpdate(p)
}

// handleBonusClick 處理 Bonus 點擊
func (g *Game) handleBonusClick(p *player.Player, msg *ws.Message) {
	g.mu.RLock()
	currentState := g.State
	g.mu.RUnlock()

	if currentState != state.StateBonusGame {
		return
	}

	var payload ws.BonusClickPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	g.mu.Lock()
	t := g.Targets[payload.TargetID]
	if t != nil && t.IsAlive {
		t.IsAlive = false
		score := t.Def.HP // HP 欄位存 ClickScore
		defID := t.DefID
		delete(g.Targets, payload.TargetID)
		g.bonusScores[p.ID] += score

		// 廣播 target_kill（讓 Client 播放拔草動畫）
		g.mu.Unlock()
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetKill,
			Payload: ws.TargetKillPayload{
				InstanceID: payload.TargetID,
				DefID:      defID,
				Multiplier: 1,
				Reward:     score,
				LaborGain:  score,
				KillerID:   p.ID,
			},
		})

		// 特殊雜草效果
		switch defID {
		case "BG003": // 發光雜草：增加倍率（加分）
			g.mu.Lock()
			g.bonusScores[p.ID] += 5 // 額外加分
			g.mu.Unlock()
		case "BG004": // 金色雜草：觸發巨大金幣（大量加分，規格書 29.3）
			g.mu.Lock()
			g.bonusScores[p.ID] += 20 // 金色雜草本身 20 分 + 額外 10 分獎勵
			g.mu.Unlock()
			// 廣播金幣特效事件（Client 播放金幣雨動畫）
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgBonusEvent,
				Payload: ws.BonusEventPayload{
					Event: "coin_shower",
				},
			})
		case "BG005": // 搗亂怪草：扣分
			g.mu.Lock()
			if g.bonusScores[p.ID] > 5 {
				g.bonusScores[p.ID] -= 5
			}
			g.mu.Unlock()
		}
		return
	}
	g.mu.Unlock()
}

// ---- 遊戲循環 ----

func (g *Game) gameLoop() {
	ticker := time.NewTicker(100 * time.Millisecond) // 10 FPS 伺服器更新
	defer ticker.Stop()

	for {
		select {
		case <-g.stopCh:
			return
		case <-ticker.C:
			g.update()
		}
	}
}

func (g *Game) update() {
	g.mu.Lock()
	currentState := g.State
	g.mu.Unlock()

	switch currentState {
	case state.StateNormalPlay, state.StateSpecialTargetEvent:
		g.updateNormalPlay()
	case state.StateBossBattle:
		g.updateBossBattle()
	case state.StateBonusGame:
		g.updateBonusGame()
	}
}

func (g *Game) updateNormalPlay() {
	now := time.Now()

	// 清除過期目標
	g.mu.Lock()
	for id, t := range g.Targets {
		if t.IsExpired() || !t.IsAlive {
			delete(g.Targets, id)
		}
	}
	targetCount := len(g.Targets)
	g.mu.Unlock()

	// 生成新目標
	if now.Sub(g.lastSpawnAt).Seconds() >= data.SpawnInterval &&
		targetCount < data.MaxTargetsOnScreen {
		g.spawnTarget()
		g.mu.Lock()
		g.lastSpawnAt = now
		g.mu.Unlock()
	}

	// 補償機制：30秒無高倍率獎勵，提高特殊目標出現率
	g.mu.Lock()
	if now.Sub(g.lastHighRewardAt).Seconds() > 30 {
		g.bonusSpecialBonus = 0.05
	}
	g.mu.Unlock()

	// 每 25-40 秒觸發一次 SpecialTargetEvent（流星等特殊目標）
	g.mu.Lock()
	if now.Sub(g.lastSpecialEventAt).Seconds() > g.nextSpecialEventIn {
		g.lastSpecialEventAt = now
		g.nextSpecialEventIn = 25 + rand.Float64()*15
		g.mu.Unlock()
		g.triggerSpecialEvent()
	} else {
		g.mu.Unlock()
	}

	// BOSS 自動觸發（規格書 28.1：每 3-5 分鐘）
	g.mu.Lock()
	if now.After(g.nextBossAt) && g.bossInstanceID == "" {
		g.nextBossAt = now.Add(time.Duration(180+rand.Intn(120)) * time.Second)
		g.mu.Unlock()
		log.Printf("[Game] Auto-triggering BOSS")
		g.triggerBoss()
	} else {
		g.mu.Unlock()
	}

	// Auto 攻擊
	g.processAutoAttack()

	// 排行榜廣播（每 10 秒）
	g.mu.Lock()
	shouldBroadcastLeaderboard := now.Sub(g.lastLeaderboardAt) >= 10*time.Second
	if shouldBroadcastLeaderboard {
		g.lastLeaderboardAt = now
	}
	g.mu.Unlock()

	if shouldBroadcastLeaderboard {
		g.broadcastLeaderboard()
	}
}

func (g *Game) updateBossBattle() {
	g.mu.RLock()
	bossID := g.bossInstanceID
	bossSpawnedAt := g.bossSpawnedAt
	g.mu.RUnlock()

	if bossID == "" {
		return
	}

	// BOSS 超時
	elapsed := time.Since(bossSpawnedAt).Seconds()
	if elapsed >= data.BossDuration {
		g.mu.Lock()
		delete(g.Targets, bossID)
		g.bossInstanceID = ""
		g.mu.Unlock()

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgBossEvent,
			Payload: ws.BossEventPayload{
				Event: "timeout",
			},
		})
		g.transitionState(state.StateNormalPlay)
		return
	}

	// 規格書 9章：BOSS 期間 Max Targets = 8（不含 BOSS 本身）
	// 清除超出限制的非 BOSS 目標（依生成時間排序，移除最舊的）
	const MaxTargetsDuringBoss = 8
	g.mu.Lock()
	type targetWithTime struct {
		id        string
		spawnedAt time.Time
	}
	nonBossTargets := make([]targetWithTime, 0)
	for id, t := range g.Targets {
		if id != bossID && t.Def.Type != data.TargetTypeBoss {
			nonBossTargets = append(nonBossTargets, targetWithTime{id: id, spawnedAt: t.SpawnedAt})
		}
	}
	if len(nonBossTargets) > MaxTargetsDuringBoss {
		// 依生成時間排序（最舊的在前）
		for i := 0; i < len(nonBossTargets); i++ {
			for j := i + 1; j < len(nonBossTargets); j++ {
				if nonBossTargets[j].spawnedAt.Before(nonBossTargets[i].spawnedAt) {
					nonBossTargets[i], nonBossTargets[j] = nonBossTargets[j], nonBossTargets[i]
				}
			}
		}
		// 移除最舊的目標
		for i := 0; i < len(nonBossTargets)-MaxTargetsDuringBoss; i++ {
			delete(g.Targets, nonBossTargets[i].id)
		}
	}
	g.mu.Unlock()
}

func (g *Game) updateBonusGame() {
	g.mu.RLock()
	bonusStart := g.bonusStartedAt
	g.mu.RUnlock()

	elapsed := time.Since(bonusStart).Seconds()
	timeLeft := data.BonusDuration - elapsed

	if timeLeft <= 0 {
		g.endBonusGame()
		return
	}

	// 廣播剩餘時間（每秒一次，避免過度廣播）
	now := time.Now()
	g.mu.Lock()
	shouldTick := now.Sub(g.lastBonusTickAt) >= time.Second
	if shouldTick {
		g.lastBonusTickAt = now
	}
	g.mu.Unlock()

	if shouldTick {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgBonusEvent,
			Payload: ws.BonusEventPayload{
				Event:    "tick",
				TimeLeft: timeLeft,
			},
		})
	}
}

// spawnTarget 生成目標
func (g *Game) spawnTarget() {
	g.mu.RLock()
	// 取所有玩家的平均 bet level（多人時更公平）
	betLevel := 5 // 預設中等
	if len(g.Players) > 0 {
		total := 0
		for _, p := range g.Players {
			total += p.BetLevel
		}
		betLevel = total / len(g.Players)
		if betLevel < 1 {
			betLevel = 1
		}
	}
	bonusSpecial := g.bonusSpecialBonus
	g.mu.RUnlock()

	def := g.SpawnSys.PickTargetDef(betLevel, bonusSpecial)
	instanceID := uuid.New().String()

	// 隨機生成位置（畫面右側進入）
	x := 1280.0 + rand.Float64()*100
	y := 100.0 + rand.Float64()*500

	t := target.NewTarget(instanceID, def, x, y)

	g.mu.Lock()
	g.Targets[instanceID] = t
	g.mu.Unlock()

	// 廣播生成事件
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTargetSpawn,
		Payload: ws.TargetSpawnPayload{
			InstanceID: instanceID,
			DefID:      def.ID,
			Name:       def.Name,
			Type:       string(def.Type),
			X:          x,
			Y:          y,
			HP:         def.HP,
			MaxHP:      def.HP,
			Speed:      def.Speed,
			Lifetime:   def.Lifetime,
			Behavior:   def.SpecialBehavior,
		},
	})
}

// triggerSpecialEvent 觸發特殊目標事件（流星、金色雜草等）
func (g *Game) triggerSpecialEvent() {
	// 強制生成一個特殊目標
	specialIDs := []string{"T103", "T104", "T101", "T102", "T105"}
	defID := specialIDs[rand.Intn(len(specialIDs))]
	def := data.Targets[defID]
	if def == nil {
		return
	}

	instanceID := uuid.New().String()
	x := 1280.0 + rand.Float64()*50
	y := 80.0 + rand.Float64()*500
	t := target.NewTarget(instanceID, def, x, y)

	g.mu.Lock()
	g.Targets[instanceID] = t
	g.mu.Unlock()

	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgTargetSpawn,
		Payload: ws.TargetSpawnPayload{
			InstanceID: instanceID,
			DefID:      def.ID,
			Name:       def.Name,
			Type:       string(def.Type),
			X:          x,
			Y:          y,
			HP:         def.HP,
			MaxHP:      def.HP,
			Speed:      def.Speed,
			Lifetime:   def.Lifetime,
			Behavior:   def.SpecialBehavior,
		},
	})
}

// triggerBoss 觸發 BOSS
func (g *Game) triggerBoss() {
	g.mu.RLock()
	currentState := g.State
	g.mu.RUnlock()

	if currentState != state.StateNormalPlay {
		return
	}

	// BOSS 警告
	g.transitionState(state.StateBossWarning)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBossEvent,
		Payload: ws.BossEventPayload{Event: "warning"},
	})

	g.safeAfterFunc(3*time.Second, func() {
		g.spawnBoss()
	})
}

func (g *Game) spawnBoss() {
	def := data.Targets["B001"]
	instanceID := uuid.New().String()
	t := target.NewTarget(instanceID, def, 1100, 360)

	g.mu.Lock()
	g.Targets[instanceID] = t
	g.bossInstanceID = instanceID
	g.bossSpawnedAt = time.Now()
	g.mu.Unlock()

	g.transitionState(state.StateBossBattle)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBossEvent,
		Payload: ws.BossEventPayload{
			Event:      "spawn",
			InstanceID: instanceID,
			HP:         def.HP,
			MaxHP:      def.HP,
		},
	})

	// 埋點：BOSS 生成
	if tracker := analytics.Get(); tracker != nil {
		tracker.Track(analytics.EventBossSpawn, "system", map[string]interface{}{
			"instance_id": instanceID,
			"boss_def":    "B001",
			"hp":          def.HP,
		})
	}
}

// triggerBonusReady 觸發 Bonus Ready
func (g *Game) triggerBonusReady() {
	g.mu.RLock()
	currentState := g.State
	g.mu.RUnlock()

	if currentState != state.StateNormalPlay {
		return
	}

	// 防止 Bonus 觸發過於頻繁（至少間隔 90 秒）
	g.mu.Lock()
	if !g.lastBonusAt.IsZero() && time.Since(g.lastBonusAt).Seconds() < 90 {
		g.mu.Unlock()
		return
	}
	g.lastBonusAt = time.Now()
	g.mu.Unlock()

	g.transitionState(state.StateBonusReady)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBonusEvent,
		Payload: ws.BonusEventPayload{Event: "ready"},
	})

	// 3秒後自動進入 Bonus Game
	g.safeAfterFunc(3*time.Second, func() {
		g.startBonusGame()
	})
}

func (g *Game) startBonusGame() {
	g.mu.Lock()
	g.bonusStartedAt = time.Now()
	// 記錄所有玩家的進場 Bet
	for id, p := range g.Players {
		g.bonusEntryBet[id] = data.GetBetDef(p.BetLevel).BetCost
		g.bonusScores[id] = 0
	}
	// 清除一般目標，生成 Bonus 目標
	g.Targets = make(map[string]*target.Target)
	g.mu.Unlock()

	g.transitionState(state.StateBonusGame)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBonusEvent,
		Payload: ws.BonusEventPayload{
			Event:    "start",
			TimeLeft: data.BonusDuration,
		},
	})

	// 埋點：Bonus 開始
	if tracker := analytics.Get(); tracker != nil {
		tracker.Track(analytics.EventBonusStart, "system", map[string]interface{}{
			"duration": data.BonusDuration,
		})
	}

	// 生成 Bonus 目標
	g.spawnBonusTargets()
}

func (g *Game) spawnBonusTargets() {
	for i := 0; i < 20; i++ {
		def := g.pickBonusTarget()
		instanceID := uuid.New().String()
		x := 100.0 + rand.Float64()*1000
		y := 100.0 + rand.Float64()*500

		// 用 HP 欄位存 ClickScore（Bonus 目標特殊處理）
		bonusDef := &data.TargetDef{
			ID:       def.ID,
			Name:     def.Name,
			Type:     data.TargetTypeBonus,
			HP:       def.ClickScore,
			Lifetime: data.BonusDuration,
		}
		t := target.NewTarget(instanceID, bonusDef, x, y)

		g.mu.Lock()
		g.Targets[instanceID] = t
		g.mu.Unlock()

		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgTargetSpawn,
			Payload: ws.TargetSpawnPayload{
				InstanceID: instanceID,
				DefID:      def.ID,
				Name:       def.Name,
				Type:       "bonus",
				X:          x,
				Y:          y,
				HP:         def.ClickScore,
				MaxHP:      def.ClickScore,
				Behavior:   def.SpecialEffect,
			},
		})
	}
}

func (g *Game) pickBonusTarget() *data.BonusTargetDef {
	total := 0
	for _, d := range data.BonusTargets {
		total += d.SpawnWeight
	}
	r := rand.Intn(total)
	cumulative := 0
	for _, d := range data.BonusTargets {
		cumulative += d.SpawnWeight
		if r < cumulative {
			return d
		}
	}
	return data.BonusTargets[0]
}

func (g *Game) endBonusGame() {
	g.mu.Lock()
	scores := make(map[string]int)
	entryBets := make(map[string]int)
	for id, score := range g.bonusScores {
		scores[id] = score
	}
	for id, bet := range g.bonusEntryBet {
		entryBets[id] = bet
	}
	g.Targets = make(map[string]*target.Target)
	g.mu.Unlock()

	// 計算每個玩家的獎勵
	for playerID, score := range scores {
		g.mu.RLock()
		p := g.Players[playerID]
		g.mu.RUnlock()

		if p == nil {
			continue
		}

		entryBet := entryBets[playerID]
		reward, multiplier := combat.CalcBonusReward(entryBet, score)
		p.AddReward(reward)
		p.ResetLaborValue()

		g.Hub.Send(playerID, &ws.Message{
			Type: ws.MsgBonusEvent,
			Payload: ws.BonusEventPayload{
				Event:      "end",
				Score:      score,
				Multiplier: multiplier,
				Reward:     reward,
			},
		})

		g.Hub.Send(playerID, &ws.Message{
			Type: ws.MsgReward,
			Payload: ws.RewardPayload{
				Source:     "bonus",
				Amount:     reward,
				Multiplier: multiplier,
				NewBalance: p.Coins,
			},
		})

		// 埋點：Bonus 結束（每個玩家的獎勵）
		if tracker := analytics.Get(); tracker != nil {
			tracker.Track(analytics.EventBonusEnd, playerID, map[string]interface{}{
				"score":      score,
				"multiplier": multiplier,
				"reward":     reward,
				"entry_bet":  entryBet,
			})
			tracker.Track(analytics.EventReward, playerID, map[string]interface{}{
				"source":     "bonus",
				"amount":     reward,
				"multiplier": multiplier,
			})
		}
	}

	g.transitionState(state.StateBonusResult)
	g.safeAfterFunc(3*time.Second, func() {
		g.transitionState(state.StateNormalPlay)
	})
}

// handleSetDisplayName 設定玩家顯示名稱（DAY-021）
func (g *Game) handleSetDisplayName(p *player.Player, msg *ws.Message) {
	var payload struct {
		DisplayName string `json:"display_name"`
	}
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}
	name := payload.DisplayName
	// 限制長度 1-16 字元
	if len(name) == 0 || len([]rune(name)) > 16 {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "invalid_name", Message: "名稱長度需在 1-16 字元之間"},
		})
		return
	}
	p.SetDisplayName(name)
	log.Printf("[Game] Player %s set display name: %s", p.ID, name)
	g.sendPlayerUpdate(p)
}

// processAutoAttack 處理自動攻擊（智慧目標選擇）
func (g *Game) processAutoAttack() {
	g.mu.RLock()
	players := make([]*player.Player, 0, len(g.Players))
	for _, p := range g.Players {
		if p.IsAuto {
			players = append(players, p)
		}
	}
	g.mu.RUnlock()

	for _, p := range players {
		if !p.CanAttack() {
			continue
		}

		bet := p.GetBetDef()
		interval := 1.0 / bet.FireRate
		if time.Since(p.LastAttackAt).Seconds() < interval {
			continue
		}

		// 智慧目標選擇
		g.mu.RLock()
		var bestTarget *target.Target

		// 1. 優先使用鎖定目標（若仍存活）
		if p.LockTargetID != "" {
			if t, ok := g.Targets[p.LockTargetID]; ok && t.IsAlive {
				bestTarget = t
			}
		}

		// 2. 若無鎖定目標，用評分系統選最佳目標
		if bestTarget == nil {
			bestScore := -1.0
			for _, t := range g.Targets {
				if !t.IsAlive || t.Def.Type == data.TargetTypeBonus {
					continue
				}
				score := g.scoreTarget(t)
				if score > bestScore {
					bestScore = score
					bestTarget = t
				}
			}
		}
		g.mu.RUnlock()

		if bestTarget == nil {
			continue
		}

		// 模擬攻擊
		g.handleAttack(p, &ws.Message{
			Type: ws.MsgAttack,
			Payload: ws.AttackPayload{
				TargetID: bestTarget.InstanceID,
			},
		})
	}
}

// scoreTarget 計算目標的自動攻擊優先分數（越高越優先）
// 評分維度：倍率、HP 殘量、距離畫面左邊緣（快要離開）
func (g *Game) scoreTarget(t *target.Target) float64 {
	score := 0.0

	// 1. 倍率加分（高倍率目標優先）
	// 特殊目標（T101-T105）倍率高，加分多
	score += t.Multiplier * 2.0

	// 2. HP 殘量加分（HP 越低越優先，快要擊破的目標）
	// HPPercent 越低，加分越多
	score += (1.0 - t.HPPercent()) * 30.0

	// 3. 位置加分（X 越小代表越靠近左邊，快要離開畫面）
	// X 在 0-1280 範圍，X 越小加分越多
	if t.X < 400 {
		score += (400.0 - t.X) * 0.1
	}

	// 4. BOSS 最高優先（確保 BOSS 戰時集中火力）
	if t.Def.Type == data.TargetTypeBoss {
		score += 500.0
	}

	return score
}

// transitionState 狀態轉換
func (g *Game) transitionState(newState state.GameState) {
	g.mu.Lock()
	oldState := g.State
	if !state.CanTransition(oldState, newState) {
		g.mu.Unlock()
		log.Printf("[Game] Invalid state transition: %s -> %s", oldState, newState)
		return
	}
	g.State = newState
	g.mu.Unlock()

	log.Printf("[Game] State: %s -> %s", oldState, newState)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgGameState,
		Payload: ws.GameStatePayload{
			State:     string(newState),
			Timestamp: time.Now().UnixMilli(),
		},
	})
}

// sendPlayerUpdate 傳送玩家狀態更新
func (g *Game) sendPlayerUpdate(p *player.Player) {
	g.Hub.Send(p.ID, &ws.Message{
		Type:    ws.MsgPlayerUpdate,
		Payload: p.Snapshot(),
	})
}

// remarshal 重新序列化 payload（處理 interface{} 轉具體型別）
func remarshal(src interface{}, dst interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return json.Unmarshal(b, dst)
}

// safeAfterFunc 感知 stopCh 的 AfterFunc，避免 Game 停止後 timer 仍執行造成 goroutine 洩漏
func (g *Game) safeAfterFunc(d time.Duration, f func()) {
	go func() {
		select {
		case <-time.After(d):
			f()
		case <-g.stopCh:
			// Game 已停止，取消 timer
		}
	}()
}

// broadcastLeaderboard 廣播排行榜給所有玩家
func (g *Game) broadcastLeaderboard() {
	entries := g.buildLeaderboard()
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgLeaderboard,
		Payload: ws.LeaderboardPayload{
			Entries:   entries,
			Timestamp: time.Now().UnixMilli(),
		},
	})
}

// buildLeaderboard 建立排行榜資料（依 SessionScore 排序，取前 10 名）
func (g *Game) buildLeaderboard() []ws.LeaderboardEntry {
	g.mu.RLock()
	snapshots := make([]player.LeaderboardSnapshot, 0, len(g.Players))
	for _, p := range g.Players {
		snapshots = append(snapshots, p.LeaderboardSnapshot())
	}
	g.mu.RUnlock()

	// 依 SessionScore 降序排序
	for i := 0; i < len(snapshots); i++ {
		for j := i + 1; j < len(snapshots); j++ {
			if snapshots[j].Score > snapshots[i].Score {
				snapshots[i], snapshots[j] = snapshots[j], snapshots[i]
			}
		}
	}

	// 取前 10 名
	maxEntries := 10
	if len(snapshots) < maxEntries {
		maxEntries = len(snapshots)
	}

	entries := make([]ws.LeaderboardEntry, maxEntries)
	for i := 0; i < maxEntries; i++ {
		entries[i] = ws.LeaderboardEntry{
			Rank:        i + 1,
			PlayerID:    snapshots[i].PlayerID,
			DisplayName: snapshots[i].DisplayName,
			Score:       snapshots[i].Score,
			MaxCoins:    snapshots[i].MaxCoins,
			KillCount:   snapshots[i].KillCount,
		}
	}
	return entries
}

// GetLeaderboardData 取得排行榜資料（供 HTTP 端點使用）
func (g *Game) GetLeaderboardData() ws.LeaderboardPayload {
	return ws.LeaderboardPayload{
		Entries:   g.buildLeaderboard(),
		Timestamp: time.Now().UnixMilli(),
	}
}

// sendAchievement 傳送成就解鎖通知給指定玩家
func (g *Game) sendAchievement(playerID string, u *achievement.AchievementUnlock) {
	if u == nil {
		return
	}
	log.Printf("[Achievement] Player %s unlocked: %s (%s)", playerID, u.Name, u.ID)
	g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgAchievement,
		Payload: ws.AchievementPayload{
			ID:          string(u.ID),
			Name:        u.Name,
			Description: u.Description,
			Icon:        u.Icon,
			UnlockedAt:  u.UnlockedAt.UnixMilli(),
		},
	})
}
