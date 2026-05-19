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
	"digital-twin/server/internal/game/jackpot"
	"digital-twin/server/internal/game/mission"
	"digital-twin/server/internal/game/state"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/player"
	"digital-twin/server/internal/store"
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
	store       store.Store // 玩家狀態持久化（DAY-026）
	initialCoins int        // 玩家初始金幣（從 config 傳入）
	missionMgr  *mission.Manager // 每日任務管理器（DAY-037）
	jackpotMgr  *jackpot.Manager  // Progressive Jackpot 管理器（DAY-048）
	jackpotHist *jackpot.History  // Jackpot 中獎歷史（DAY-048e）

	// 計時器
	lastSpawnAt        time.Time
	bossSpawnedAt      time.Time
	bonusStartedAt     time.Time
	lastBonusAt        time.Time
	lastBonusTickAt    time.Time  // Bonus tick 廣播計時（每秒一次）
	lastSpecialEventAt time.Time
	nextSpecialEventIn float64
	lastJackpotSaveAt  time.Time  // Jackpot 狀態儲存計時（每 30 秒，DAY-049d）
	bossInstanceID     string
	nextBossAt         time.Time  // BOSS 自動觸發時間（規格書 28.1）
	lastLeaderboardAt  time.Time  // 排行榜廣播計時（每 10 秒一次）
	lastJackpotAt      time.Time  // Jackpot 廣播計時（每 5 秒一次，DAY-048）

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
	return NewGameWithStore(id, hub, nil, 10000)
}

// NewGameWithStore 建立新遊戲（帶 Store 和初始金幣設定）
func NewGameWithStore(id string, hub *ws.Hub, s store.Store, initialCoins int) *Game {
	g := &Game{
		ID:                 id,
		State:              state.StateNormalPlay,
		Players:            make(map[string]*player.Player),
		Targets:            make(map[string]*target.Target),
		SpawnSys:           target.NewSpawnSystem(),
		Hub:                hub,
		store:              s,
		initialCoins:       initialCoins,
		missionMgr:         mission.NewManager(),
		jackpotMgr:         jackpot.NewManager(),
		jackpotHist:        jackpot.NewHistory(10),
		lastSpawnAt:        time.Now(),
		lastSpecialEventAt: time.Now(),
		nextSpecialEventIn: 30,
		bonusScores:        make(map[string]int),
		bonusEntryBet:      make(map[string]int),
		// BOSS 自動觸發：遊戲開始後 3-5 分鐘（規格書 28.1）
		nextBossAt: time.Now().Add(time.Duration(180+rand.Intn(120)) * time.Second),
		stopCh:     make(chan struct{}),
	}
	return g
}

// GetState 取得目前遊戲狀態（thread-safe）
func (g *Game) GetState() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return string(g.State)
}

// GetActiveTargetCount 取得當前活躍目標物數量（thread-safe）
// 供 /metrics 端點使用，讓 Grafana 能監控目標物數量
func (g *Game) GetActiveTargetCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.Targets)
}

// GetJackpotSnapshot 取得 Jackpot 池快照（thread-safe，DAY-048）
// 供 /metrics 端點使用，讓 Grafana 能監控 Jackpot 池金額
func (g *Game) GetJackpotSnapshot() map[string]int {
	snap := g.jackpotMgr.GetSnapshot()
	return map[string]int{
		"mini":  snap[jackpot.LevelMini],
		"major": snap[jackpot.LevelMajor],
		"grand": snap[jackpot.LevelGrand],
	}
}

// GetJackpotHistory 取得最近 Jackpot 中獎記錄（DAY-048e）
func (g *Game) GetJackpotHistory(n int) []jackpot.WinRecord {
	return g.jackpotHist.GetRecent(n)
}

// GetJackpotDailyStats 取得今日 Jackpot 統計（DAY-049）
func (g *Game) GetJackpotDailyStats() jackpot.DailyStats {
	return g.jackpotHist.GetDailyStats()
}

// Start 啟動遊戲循環
func (g *Game) Start() {
	log.Printf("[Game] %s started", g.ID)
	// 恢復 Jackpot 池狀態（DAY-049d）
	g.loadJackpotState()
	go g.gameLoop()
}

// Stop 停止遊戲
func (g *Game) Stop() {
	close(g.stopCh)
}

// AddPlayer 加入玩家（從 Store 恢復狀態，若無則建立新玩家）
func (g *Game) AddPlayer(playerID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.Players[playerID]; !exists {
		p := player.NewPlayer(playerID, g.initialCoins)

		// 從 Store 恢復玩家狀態（若有）
		if g.store != nil {
			if saved, err := g.store.LoadPlayer(playerID); err == nil && saved != nil {
				p.Coins = int(saved.Coins)
				p.MaxCoins = int(saved.MaxCoins)
				p.KillCount = saved.KillCount
				if saved.BetLevel >= 1 && saved.BetLevel <= 10 {
					p.BetLevel = saved.BetLevel
				}
				if saved.DisplayName != "" {
					p.DisplayName = saved.DisplayName
				}
				log.Printf("[Game] Player %s restored: coins=%d, kills=%d", playerID, p.Coins, p.KillCount)
			}
		}

		g.Players[playerID] = p
		log.Printf("[Game] Player %s joined game %s", playerID, g.ID)

		// 非同步發送任務列表（連線後立即讓玩家看到今日任務）
		go func() {
			time.Sleep(200 * time.Millisecond) // 等待連線穩定
			g.sendMissionUpdate(playerID)
		}()
	}
}

// RemovePlayer 移除玩家（儲存狀態到 Store）
func (g *Game) RemovePlayer(playerID string) {
	g.mu.Lock()
	p := g.Players[playerID]
	delete(g.Players, playerID)
	g.mu.Unlock()

	// 儲存玩家狀態到 Store（讓下次加入時能恢復）
	if g.store != nil && p != nil {
		state := &store.PlayerState{
			PlayerID:     p.ID,
			DisplayName:  p.DisplayName,
			Coins:        int64(p.Coins),
			Labor:        p.LaborValue,
			BetLevel:     p.BetLevel,
			SessionScore: int64(p.SessionScore),
			MaxCoins:     int64(p.MaxCoins),
			KillCount:    p.KillCount,
			RoomID:       g.ID,
		}
		if err := g.store.SavePlayer(state); err != nil {
			log.Printf("[Game] Failed to save player %s: %v", playerID, err)
		} else {
			log.Printf("[Game] Player %s saved: coins=%d", playerID, p.Coins)
		}
		// 更新排行榜
		g.store.UpdateLeaderboard(playerID, int64(p.SessionScore))
	}

	log.Printf("[Game] Player %s left game %s", playerID, g.ID)
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
	case ws.MsgGetMissions:
		// 查詢任務列表（DAY-037）
		g.sendMissionUpdate(clientID)
	case ws.MsgClaimMission:
		// 領取任務獎勵（DAY-037）
		g.handleClaimMission(p, msg)
	case ws.MsgClientPerf:
		// Client 端效能數據上報（DAY-045）
		g.handleClientPerf(clientID, msg)
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

	// Progressive Jackpot 貢獻（DAY-048）：每次攻擊抽取 0.5% 進入 Jackpot 池
	if jackpotWin := g.jackpotMgr.Contribute(betCost, p.ID); jackpotWin != nil {
		g.handleJackpotWin(p, jackpotWin)
	}

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

	// 連擊系統（DAY-022）：2 秒內連續擊破觸發 Combo
	comboCount, laborBonus := p.AddKillCombo()
	if comboCount >= 2 {
		// 廣播連擊事件（只傳給觸發者，避免干擾其他玩家）
		g.Hub.Send(p.ID, &ws.Message{
			Type: ws.MsgComboEvent,
			Payload: ws.ComboEventPayload{
				ComboCount: comboCount,
				LaborBonus: laborBonus,
				PlayerID:   p.ID,
			},
		})
		// 任務進度：連擊達人（DAY-038）
		// 每次達成 2+ 連擊，累積 +1（不是 +comboCount，避免連擊串讓任務太快完成）
		go g.updateMissionProgress(p.ID, mission.MissionCombo, 1)
	}

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

	// 任務進度更新（DAY-037）
	go func() {
		// 擊破目標任務
		g.updateMissionProgress(p.ID, mission.MissionKillTargets, 1)
		// 累積金幣任務
		g.updateMissionProgress(p.ID, mission.MissionEarnCoins, result.Reward)
		// 高倍率目標任務（30x+）
		if result.Multiplier >= 30.0 {
			g.updateMissionProgress(p.ID, mission.MissionKillHighMult, 1)
		}
	}()

	// 累積勞動值（加入 Combo 加成）
	laborGain := result.LaborGain
	if laborBonus > 0 {
		laborGain = int(float64(laborGain) * (1.0 + laborBonus))
	}
	bonusTriggered := p.AddLaborValue(laborGain)

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

	// 任務進度：擊敗 BOSS（DAY-037）
	go g.updateMissionProgress(p.ID, mission.MissionKillBoss, 1)

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
	// Jackpot 廣播（每 5 秒，DAY-048）
	shouldBroadcastJackpot := now.Sub(g.lastJackpotAt) >= 5*time.Second
	if shouldBroadcastJackpot {
		g.lastJackpotAt = now
	}
	// Jackpot 狀態儲存（每 30 秒，DAY-049d）
	shouldSaveJackpot := now.Sub(g.lastJackpotSaveAt) >= 30*time.Second
	if shouldSaveJackpot {
		g.lastJackpotSaveAt = now
	}
	g.mu.Unlock()

	if shouldBroadcastLeaderboard {
		g.broadcastLeaderboard()
	}
	if shouldBroadcastJackpot {
		g.broadcastJackpot()
	}
	// Jackpot 狀態儲存（每 30 秒，非同步，DAY-049d）
	if shouldSaveJackpot {
		go g.saveJackpotState()
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
			Multiplier: t.Multiplier,
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
			Multiplier: t.Multiplier,
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

	// BOSS HP 依玩家平均 bet 等級縮放（DAY-044b）
	// 設計原則：玩家在 60 秒內有 ~50% 機率打死 BOSS
	// 期望攻擊次數 = fire_rate × 60 × 玩家數
	// BOSS HP = 期望攻擊次數 × bet_cost × 0.5
	g.mu.RLock()
	avgBetLevel := 5
	playerCount := len(g.Players)
	if playerCount > 0 {
		total := 0
		for _, p := range g.Players {
			total += p.BetLevel
		}
		avgBetLevel = total / playerCount
		if avgBetLevel < 1 {
			avgBetLevel = 1
		}
	}
	g.mu.RUnlock()

	betDef := data.GetBetDef(avgBetLevel)
	// 單人：fire_rate × 60 × bet_cost × 0.5
	// 多人：每個玩家都在打，HP 乘以玩家數（但有上限避免太難）
	effectivePlayers := playerCount
	if effectivePlayers < 1 {
		effectivePlayers = 1
	}
	if effectivePlayers > 4 {
		effectivePlayers = 4 // 最多 4 人效果，避免 HP 過高
	}
	bossHP := int(betDef.FireRate * 60 * float64(betDef.BetCost) * 0.5 * float64(effectivePlayers))
	if bossHP < 100 {
		bossHP = 100 // 最低 100 HP
	}
	if bossHP > 10000 {
		bossHP = 10000 // 最高 10000 HP
	}

	// 建立動態 HP 的 BOSS def（不修改原始 def）
	bossDef := *def
	bossDef.HP = bossHP

	t := target.NewTarget(instanceID, &bossDef, 1100, 360)

	g.mu.Lock()
	g.Targets[instanceID] = t
	g.bossInstanceID = instanceID
	g.bossSpawnedAt = time.Now()
	g.mu.Unlock()

	log.Printf("[Game] BOSS spawned: HP=%d (avgBetLV=%d, players=%d)", bossHP, avgBetLevel, playerCount)

	g.transitionState(state.StateBossBattle)
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgBossEvent,
		Payload: ws.BossEventPayload{
			Event:      "spawn",
			InstanceID: instanceID,
			HP:         bossHP,
			MaxHP:      bossHP,
		},
	})

	// 埋點：BOSS 生成
	if tracker := analytics.Get(); tracker != nil {
		tracker.Track(analytics.EventBossSpawn, "system", map[string]interface{}{
			"instance_id": instanceID,
			"boss_def":    "B001",
			"hp":          bossHP,
			"avg_bet_lv":  avgBetLevel,
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

		// 任務進度：完成 Bonus Game（DAY-037）
		go func(pid string, r int) {
			g.updateMissionProgress(pid, mission.MissionPlayBonus, 1)
			g.updateMissionProgress(pid, mission.MissionEarnCoins, r)
		}(playerID, reward)
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
// 評分維度：倍率、HP 殘量、目標存活時間（存活越久代表快要消失）
func (g *Game) scoreTarget(t *target.Target) float64 {
	score := 0.0

	// 1. 倍率加分（高倍率目標優先）
	score += t.Multiplier * 2.0

	// 2. HP 殘量加分（HP 越低越優先，快要擊破的目標）
	score += (1.0 - t.HPPercent()) * 30.0

	// 3. 存活時間加分（存活越久代表快要超時消失，優先打掉）
	// 存活超過 50% Lifetime 開始加分，超過 80% 大幅加分
	elapsed := time.Since(t.SpawnedAt).Seconds()
	if t.Def.Lifetime > 0 {
		lifeRatio := elapsed / t.Def.Lifetime
		if lifeRatio > 0.8 {
			score += 40.0 // 快要消失，最高優先
		} else if lifeRatio > 0.5 {
			score += 15.0 // 過半存活時間
		}
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

// GetMissionResetAt 取得每日任務下次重置時間（thread-safe）
func (g *Game) GetMissionResetAt() time.Time {
	return g.missionMgr.ResetAt()
}

// SpectatorSnapshot 觀戰快照（DAY-023）：包含遊戲狀態 + 所有目標 + 排行榜
type SpectatorSnapshot struct {
	GameState   string                  `json:"game_state"`
	Targets     []ws.TargetSpawnPayload `json:"targets"`
	Leaderboard ws.LeaderboardPayload   `json:"leaderboard"`
	PlayerCount int                     `json:"player_count"`
	Timestamp   int64                   `json:"timestamp"`
}

// GetSpectatorSnapshot 取得觀戰快照（供觀戰者連線時初始化用）
func (g *Game) GetSpectatorSnapshot() SpectatorSnapshot {
	g.mu.RLock()
	defer g.mu.RUnlock()

	targets := make([]ws.TargetSpawnPayload, 0, len(g.Targets))
	for _, t := range g.Targets {
		if !t.IsAlive {
			continue
		}
		targets = append(targets, ws.TargetSpawnPayload{
			InstanceID: t.InstanceID,
			DefID:      t.DefID,
			Name:       t.Def.Name,
			Type:       string(t.Def.Type),
			X:          t.X,
			Y:          t.Y,
			HP:         t.HP,
			MaxHP:      t.MaxHP,
			Speed:      t.Def.Speed,
			Lifetime:   t.Def.Lifetime,
			Behavior:   t.Def.SpecialBehavior,
			Multiplier: t.Multiplier,
		})
	}

	return SpectatorSnapshot{
		GameState:   string(g.State),
		Targets:     targets,
		Leaderboard: g.GetLeaderboardData(),
		PlayerCount: len(g.Players),
		Timestamp:   time.Now().UnixMilli(),
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

// ── 每日任務系統（DAY-037）──────────────────────────────────────

// sendMissionUpdate 傳送任務列表給指定玩家
func (g *Game) sendMissionUpdate(playerID string) {
	statuses := g.missionMgr.GetPlayerMissions(playerID)
	payloads := make([]ws.MissionPayload, 0, len(statuses))
	for _, s := range statuses {
		payloads = append(payloads, ws.MissionPayload{
			ID:            s.Mission.ID,
			Name:          s.Mission.Name,
			Description:   s.Mission.Description,
			Icon:          s.Mission.Icon,
			Target:        s.Mission.Target,
			Current:       s.Progress.Current,
			Completed:     s.Progress.Completed,
			RewardClaimed: s.Progress.RewardClaimed,
			Reward:        s.Mission.Reward,
		})
	}
	g.Hub.Send(playerID, &ws.Message{
		Type: ws.MsgMissionUpdate,
		Payload: ws.MissionUpdatePayload{
			PlayerID:      playerID,
			Missions:      payloads,
			ResetAt:       g.missionMgr.ResetAt().UnixMilli(),
			ResetTimezone: "UTC+8",
		},
	})
}

// updateMissionProgress 更新任務進度並通知玩家
// 由各遊戲事件（擊殺、BOSS、Bonus）呼叫
func (g *Game) updateMissionProgress(playerID string, mType mission.MissionType, amount int) {
	completed := g.missionMgr.UpdateProgress(playerID, mType, amount)

	// 通知任務完成
	for _, m := range completed {
		log.Printf("[Mission] Player %s completed: %s", playerID, m.Name)
		g.Hub.Send(playerID, &ws.Message{
			Type: ws.MsgMissionComplete,
			Payload: ws.MissionCompletePayload{
				MissionID: m.ID,
				Name:      m.Name,
				Icon:      m.Icon,
				Reward:    m.Reward,
			},
		})
	}

	// 更新任務進度（有變化才發送）
	if amount > 0 {
		g.sendMissionUpdate(playerID)
	}
}

// handleClaimMission 處理領取任務獎勵
func (g *Game) handleClaimMission(p *player.Player, msg *ws.Message) {
	var payload ws.ClaimMissionPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}
	if payload.MissionID == "" {
		return
	}

	reward := g.missionMgr.ClaimReward(p.ID, payload.MissionID)
	if reward <= 0 {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "mission_not_claimable", Message: "任務未完成或已領取"},
		})
		return
	}

	// 發放獎勵
	p.AddReward(reward)
	log.Printf("[Mission] Player %s claimed reward %d for mission %s", p.ID, reward, payload.MissionID)

	// 通知玩家
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgReward,
		Payload: ws.RewardPayload{
			Source:     "mission",
			Amount:     reward,
			Multiplier: 1.0,
			NewBalance: p.Coins,
		},
	})
	g.sendPlayerUpdate(p)
	g.sendMissionUpdate(p.ID)
}

// handleClientPerf 處理 Client 端效能數據上報（DAY-045）
// Client 每 30 秒發送一次，Server 記錄並暴露到 /metrics
// 同時檢查高延遲玩家並輸出警告 log
func (g *Game) handleClientPerf(clientID string, msg *ws.Message) {
	var payload ws.ClientPerfPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	// 更新 Hub 中的 Client 效能快照
	g.Hub.UpdateClientPerf(clientID, payload.FPS, payload.MemoryMB, payload.DrawCalls, payload.Quality)

	// 高延遲警告（DAY-045）：Client 端 ping > 200ms 輸出警告 log
	// 這讓運維人員能識別網路品質差的玩家
	if payload.PingMs > 200 {
		log.Printf("[PerfAlert] High latency player %s: ping=%dms fps=%.1f quality=%s",
			clientID, payload.PingMs, payload.FPS, payload.Quality)
	}

	// 低 FPS 警告：Client 端 FPS < 20 輸出警告 log
	if payload.FPS > 0 && payload.FPS < 20 {
		log.Printf("[PerfAlert] Low FPS player %s: fps=%.1f memory=%.1fMB drawcalls=%d quality=%s",
			clientID, payload.FPS, payload.MemoryMB, payload.DrawCalls, payload.Quality)
	}
}

// handleJackpotWin 處理 Jackpot 中獎（DAY-048）
func (g *Game) handleJackpotWin(p *player.Player, win *jackpot.JackpotWin) {
	// 發放獎勵給中獎玩家
	p.AddReward(win.Amount)

	// 取得顯示名稱
	displayName := p.DisplayName
	if displayName == "" {
		displayName = p.ID[:8]
	}

	// 記錄到歷史（DAY-048e）
	g.jackpotHist.Add(win, displayName)

	// 廣播中獎通知給所有玩家
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgJackpotWin,
		Payload: ws.JackpotWinPayload{
			Level:      string(win.Level),
			Amount:     win.Amount,
			WinnerID:   p.ID,
			WinnerName: displayName,
			NewBalance: p.Coins,
		},
	})

	// 更新中獎玩家的狀態
	g.sendPlayerUpdate(p)

	log.Printf("[Jackpot] %s won %s jackpot: %d coins (player: %s)",
		p.ID, win.Level, win.Amount, displayName)
}

// broadcastJackpot 廣播 Jackpot 池當前金額（每 5 秒，DAY-048）
func (g *Game) broadcastJackpot() {
	snap := g.jackpotMgr.GetSnapshot()
	g.Hub.Broadcast(&ws.Message{
		Type: ws.MsgJackpotUpdate,
		Payload: ws.JackpotUpdatePayload{
			Mini:  snap[jackpot.LevelMini],
			Major: snap[jackpot.LevelMajor],
			Grand: snap[jackpot.LevelGrand],
		},
	})
}

// saveJackpotState 儲存 Jackpot 池狀態到 Store（DAY-049d）
// 每 30 秒自動呼叫，確保 Server 重啟後能恢復 Jackpot 池
func (g *Game) saveJackpotState() {
	if g.store == nil {
		return
	}
	state := g.jackpotMgr.SaveState()
	key := "jackpot_state:" + g.ID
	if err := g.store.SetJSON(key, state, 7*24*time.Hour); err != nil {
		log.Printf("[Jackpot] Failed to save state: %v", err)
	}
}

// loadJackpotState 從 Store 恢復 Jackpot 池狀態（DAY-049d）
// 在 Game 啟動時呼叫
func (g *Game) loadJackpotState() {
	if g.store == nil {
		return
	}
	key := "jackpot_state:" + g.ID
	var state jackpot.PoolState
	if err := g.store.GetJSON(key, &state); err != nil {
		// 找不到或解析失敗，使用預設值（正常情況）
		return
	}
	g.jackpotMgr.LoadState(state)
	log.Printf("[Jackpot] Restored state: mini=%d major=%d grand=%d",
		state.Mini, state.Major, state.Grand)
}
