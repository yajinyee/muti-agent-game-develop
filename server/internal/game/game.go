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
	"digital-twin/server/internal/game/dailyboss"
	"digital-twin/server/internal/game/event"
	"digital-twin/server/internal/game/friend"
	"digital-twin/server/internal/game/guild"
	"digital-twin/server/internal/game/guildwar"
	"digital-twin/server/internal/game/jackpot"
	"digital-twin/server/internal/game/mission"
	"digital-twin/server/internal/game/season"
	"digital-twin/server/internal/game/state"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/game/tournament"
	"digital-twin/server/internal/game/vip"
	"digital-twin/server/internal/game/referral"
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
	missionMgr  *mission.Manager    // 每日任務管理器（DAY-037）
	jackpotMgr  *jackpot.Manager    // Progressive Jackpot 管理器（DAY-048）
	jackpotHist *jackpot.History    // Jackpot 中獎歷史（DAY-048e）
	tournamentMgr *tournament.Tournament // 週賽管理器（DAY-066）
	Season      *season.Manager     // 賽季通行證管理器（DAY-072）
	Friends     *friend.Manager     // 好友系統管理器（DAY-073）
	Guild       *guild.Manager      // 公會系統管理器（DAY-074）
	GuildWar    *guildwar.Manager   // 公會戰管理器（DAY-076）
	DailyBoss   *dailyboss.Manager  // 每日 BOSS 挑戰管理器（DAY-077）
	VIP         *vip.Manager        // VIP 等級管理器（DAY-078）
	Event       *event.Manager      // 限時活動管理器（DAY-079）
	Referral    *referral.Manager   // 推薦碼管理器（DAY-082）

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
	lastTournamentAt   time.Time  // 週賽排名廣播計時（每 30 秒一次，DAY-066）
	lastGuildWarAt     time.Time  // 公會戰廣播計時（每 60 秒一次，DAY-076）
	lastDailyBossAt    time.Time  // 每日 BOSS 廣播計時（每 30 秒一次，DAY-077）
	lastEventAt        time.Time  // 限時活動廣播計時（每 30 秒一次，DAY-079）

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
		tournamentMgr:      tournament.New(),
		Season:             season.New(),
		Friends:            friend.New(),
		Guild:              guild.New(),
		GuildWar:           guildwar.New(),
		DailyBoss:          dailyboss.New(),
		VIP:                vip.New(),
		Event:              event.New(30 * time.Minute),
		Referral:           referral.NewManager(),
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

// Start 啟動遊戲循環
func (g *Game) Start() {
	log.Printf("[Game] %s started", g.ID)
	// 恢復 Jackpot 池狀態（DAY-049d）
	g.loadJackpotState()
	// 啟動連擊超時檢查（DAY-083）
	g.startStreakTicker()
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
		var savedState *store.PlayerState
		if g.store != nil {
			if saved, err := g.store.LoadPlayer(playerID); err == nil && saved != nil {
				savedState = saved
				p.Coins = int(saved.Coins)
				p.MaxCoins = int(saved.MaxCoins)
				p.KillCount = saved.KillCount
				if saved.BetLevel >= 1 && saved.BetLevel <= 10 {
					p.BetLevel = saved.BetLevel
				}
				if saved.DisplayName != "" {
					p.DisplayName = saved.DisplayName
				}
				// 恢復登入資訊（DAY-065）
				p.LoginStreak = saved.LoginStreak
				p.MaxLoginStreak = saved.MaxLoginStreak
				p.LastLoginDate = saved.LastLoginDate
				// 恢復外觀資訊（DAY-071）
				if saved.EquippedSkin != "" {
					p.EquippedSkin = saved.EquippedSkin
				}
				if len(saved.OwnedSkins) > 0 {
					p.OwnedSkins = saved.OwnedSkins
				}
				log.Printf("[Game] Player %s restored: coins=%d, kills=%d, streak=%d, skin=%s", playerID, p.Coins, p.KillCount, p.LoginStreak, p.EquippedSkin)
			}
		}

		g.Players[playerID] = p
		log.Printf("[Game] Player %s joined game %s", playerID, g.ID)

		// 非同步發送任務列表 + 每日登入獎勵 + 賽季快照（連線後立即讓玩家看到）
		go func() {
			time.Sleep(200 * time.Millisecond) // 等待連線穩定
			g.sendMissionUpdate(playerID)
			g.checkAndSendDailyBonus(playerID, savedState)
			// 發送賽季通行證快照（DAY-072）
			g.mu.RLock()
			pp := g.Players[playerID]
			g.mu.RUnlock()
			if pp != nil {
				g.sendSeasonUpdate(pp)
				// 發送好友列表（DAY-073）
				g.sendFriendList(pp)
				// 通知好友上線（DAY-073）
				g.notifyFriendsOnline(playerID, pp.DisplayName)
				// 更新公會在線狀態 + 發送公會資訊（DAY-074）
				g.Guild.SetOnlineStatus(playerID, true)
				g.sendGuildUpdate(pp)
				// 發送 VIP 狀態（DAY-078）
				g.sendVIPUpdate(pp)
				// 發送限時活動狀態（DAY-079）
				g.sendEventUpdate(pp)
				// 發送圖鑑狀態（DAY-081）
				g.sendCodexUpdate(pp.ID)
				// 發送推薦碼資訊（DAY-082）
				g.sendReferralInfo(pp)
			}
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
			PlayerID:       p.ID,
			DisplayName:    p.DisplayName,
			Coins:          int64(p.Coins),
			Labor:          p.LaborValue,
			BetLevel:       p.BetLevel,
			SessionScore:   int64(p.SessionScore),
			MaxCoins:       int64(p.MaxCoins),
			KillCount:      p.KillCount,
			RoomID:         g.ID,
			LoginStreak:    p.LoginStreak,
			MaxLoginStreak: p.MaxLoginStreak,
			LastLoginDate:  p.LastLoginDate,
			// 外觀資訊（DAY-071）
			EquippedSkin: p.EquippedSkin,
			OwnedSkins:   p.OwnedSkins,
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

	// 通知好友下線（DAY-073）
	if p != nil {
		go g.notifyFriendsOffline(playerID, p.DisplayName)
		// 更新公會在線狀態（DAY-074）
		g.Guild.SetOnlineStatus(playerID, false)
	}
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
	case ws.MsgUpgradeWeapon:
		// 武器升級（DAY-067）
		g.handleUpgradeWeapon(p, msg)
	case ws.MsgSetTitle:
		// 設定顯示稱號（DAY-068）
		g.handleSetTitle(p, msg)
	case ws.MsgBuySkin:
		// 購買砲台外觀（DAY-071）
		g.handleBuySkin(p, msg)
	case ws.MsgEquipSkin:
		// 裝備砲台外觀（DAY-071）
		g.handleEquipSkin(p, msg)
	case ws.MsgClaimSeasonLevel:
		// 領取賽季等級獎勵（DAY-072）
		g.handleClaimSeasonLevel(p, msg)
	case ws.MsgSendFriendRequest:
		// 發送好友請求（DAY-073）
		g.handleSendFriendRequest(p, msg)
	case ws.MsgAcceptFriendRequest:
		// 接受好友請求（DAY-073）
		g.handleAcceptFriendRequest(p, msg)
	case ws.MsgRejectFriendRequest:
		// 拒絕好友請求（DAY-073）
		g.handleRejectFriendRequest(p, msg)
	case ws.MsgRemoveFriend:
		// 移除好友（DAY-073）
		g.handleRemoveFriend(p, msg)
	case ws.MsgGetFriendList:
		// 查詢好友列表（DAY-073）
		g.handleGetFriendList(p, msg)
	// 公會系統（DAY-074）
	case ws.MsgCreateGuild:
		g.handleCreateGuild(p, msg)
	case ws.MsgJoinGuild:
		g.handleJoinGuild(p, msg)
	case ws.MsgLeaveGuild:
		g.handleLeaveGuild(p, msg)
	case ws.MsgKickGuildMember:
		g.handleKickGuildMember(p, msg)
	case ws.MsgPromoteGuildMember:
		g.handlePromoteGuildMember(p, msg)
	case ws.MsgGetGuildInfo:
		g.handleGetGuildInfo(p, msg)
	case ws.MsgGetGuildList:
		g.handleGetGuildList(p, msg)
	case ws.MsgGuildChat:
		g.handleGuildChat(p, msg)
	// 公會戰系統（DAY-076）
	case ws.MsgGetGuildWarStatus:
		g.handleGetGuildWarStatus(p)
	// 每日 BOSS 挑戰（DAY-077）
	case ws.MsgGetDailyBoss:
		g.handleGetDailyBoss(p)
	case ws.MsgDailyBossAttack:
		g.handleDailyBossAttack(p, msg)
	// VIP 等級系統（DAY-078）
	case ws.MsgGetVIPStatus:
		g.handleGetVIPStatus(p)
	case ws.MsgClaimVIPWeekly:
		g.handleClaimVIPWeekly(p)
	// 限時活動系統（DAY-079）
	case ws.MsgGetEventStatus:
		g.handleGetEventStatus(p)
	// 魚類圖鑑系統（DAY-081）
	case ws.MsgGetCodex:
		g.handleGetCodex(p.ID)
	// 推薦碼系統（DAY-082）
	case ws.MsgGetReferralInfo:
		g.handleGetReferralInfo(p)
	case ws.MsgUseReferralCode:
		g.handleUseReferralCode(p, msg)
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
		PlayerID:       p.ID,
		TargetID:       targetID,
		BetLevel:       p.BetLevel,
		IsAuto:         p.IsAuto,
		IsLock:         targetID != "",
		ClickX:         payload.ClickX,
		ClickY:         payload.ClickY,
		WeaponPowerMod: p.GetWeaponPowerMod(), // 武器攻擊力加成（DAY-067）
		EventKillAdd:   g.getEventKillChanceAdd(), // 限時活動擊破率加成（DAY-079）
	}

	result := combat.ProcessAttack(req, t)

	// Progressive Jackpot 貢獻（DAY-048）：每次攻擊抽取 0.5% 進入 Jackpot 池
	if jackpotWin := g.jackpotMgr.Contribute(betCost, p.ID); jackpotWin != nil {
		g.handleJackpotWin(p, jackpotWin)
	}

	// VIP 消費記錄 + 金幣返還（DAY-078）
	g.notifyVIPSpend(p, betCost)

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

	// 發放獎勵（套用限時活動倍率，DAY-079）
	eventRewardMult := g.Event.GetRewardMult()
	finalReward := result.Reward
	if eventRewardMult > 1.0 {
		finalReward = int(float64(result.Reward) * eventRewardMult)
	}
	// 套用連擊倍率（DAY-083）
	streakMult := g.notifyStreakKill(p)
	if streakMult > 1.0 {
		finalReward = int(float64(finalReward) * streakMult)
	}
	rewardUnlocks := p.AddReward(finalReward)
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
		g.sendAchievements(p.ID, []*achievement.AchievementUnlock{u})
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
			Reward:     finalReward, // 套用活動倍率後的獎勵（DAY-079）
			LaborGain:  result.LaborGain,
			KillerID:   p.ID,
			Quality:    string(t.Quality),
		},
	})

	// legendary 品質目標擊破：10% 機率觸發提前 BOSS（DAY-070）
	// 參考 Fishing Frenzy Chapter 3 的 S 級魚召喚 Boss Fish 機制
	if t.Quality == target.QualityLegendary && t.Def.Type != data.TargetTypeBoss {
		g.mu.RLock()
		currentState := g.State
		g.mu.RUnlock()
		if currentState == state.StateNormalPlay && rand.Intn(10) == 0 {
			// 廣播特殊 BOSS 召喚通知
			g.Hub.Broadcast(&ws.Message{
				Type: ws.MsgBossEvent,
				Payload: ws.BossEventPayload{
					Event: "legendary_summon",
				},
			})
			// 縮短 BOSS 觸發冷卻（讓 BOSS 更快出現）
			g.mu.Lock()
			g.nextBossAt = time.Now().Add(5 * time.Second) // 5 秒後觸發 BOSS
			g.mu.Unlock()
		}
	}

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
			g.sendAchievements(p.ID, []*achievement.AchievementUnlock{u})
		}
		g.triggerBonusReady()
	}

	// 週賽積分：擊破目標（DAY-066）
	g.tournamentMgr.AddPoints(p.ID, p.DisplayName, tournament.PointKill, result.Multiplier)
	// 賽季積分同步（DAY-072）：擊破積分 = max(1, floor(multiplier))
	killPts := int(result.Multiplier)
	if killPts < 1 {
		killPts = 1
	}
	newLevels := g.addSeasonPoints(p.ID, killPts)
	g.checkSeasonLevelNotify(p, newLevels)
	// 公會任務進度：擊破目標（DAY-074）
	guildID := g.Guild.GetPlayerGuildID(p.ID)
	if guildID != "" {
		completedTasks := g.Guild.UpdateTaskProgress(p.ID, guild.TaskKillTargets, 1)
		g.notifyGuildTaskComplete(guildID, completedTasks)
	}
	// 公會戰積分：擊殺（DAY-076）
	go g.notifyGuildWarKill(p.ID, int(result.Multiplier))
	// 每日 BOSS 傷害貢獻（DAY-077）
	go g.notifyDailyBossKill(p, result.Multiplier)
	// 魚類圖鑑：記錄擊破（DAY-081）
	g.notifyCodexKill(p.ID, t.DefID, result.Multiplier)
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

	// 生成新目標（套用限時活動 SpawnMult，DAY-079）
	eventSpawnMult := g.getEventSpawnMult()
	effectiveMaxTargets := int(float64(data.MaxTargetsOnScreen) * eventSpawnMult)
	effectiveInterval := data.SpawnInterval / eventSpawnMult
	if now.Sub(g.lastSpawnAt).Seconds() >= effectiveInterval &&
		targetCount < effectiveMaxTargets {
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
	// 週賽排名廣播（每 30 秒，DAY-066）
	shouldBroadcastTournament := now.Sub(g.lastTournamentAt) >= 30*time.Second
	if shouldBroadcastTournament {
		g.lastTournamentAt = now
	}
	// 公會戰廣播（每 60 秒，DAY-076）
	shouldBroadcastGuildWar := now.Sub(g.lastGuildWarAt) >= 60*time.Second
	if shouldBroadcastGuildWar {
		g.lastGuildWarAt = now
	}
	// 每日 BOSS 廣播（每 30 秒，DAY-077）
	shouldBroadcastDailyBoss := now.Sub(g.lastDailyBossAt) >= 30*time.Second
	if shouldBroadcastDailyBoss {
		g.lastDailyBossAt = now
	}
	// 限時活動廣播（每 30 秒，DAY-079）
	shouldBroadcastEvent := now.Sub(g.lastEventAt) >= 30*time.Second
	if shouldBroadcastEvent {
		g.lastEventAt = now
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
	// 週賽排名廣播（每 30 秒，DAY-066）
	if shouldBroadcastTournament {
		go g.broadcastTournament()
	}
	// 公會戰廣播（每 60 秒，DAY-076）
	if shouldBroadcastGuildWar {
		go g.broadcastGuildWar()
	}
	// 每日 BOSS 廣播（每 30 秒，DAY-077）
	if shouldBroadcastDailyBoss {
		go g.broadcastDailyBoss()
	}
	// 限時活動 Tick + 廣播（每 30 秒，DAY-079）
	if shouldBroadcastEvent {
		go g.tickAndBroadcastEvent()
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
			InstanceID:   instanceID,
			DefID:        def.ID,
			Name:         def.Name,
			Type:         string(def.Type),
			X:            x,
			Y:            y,
			HP:           def.HP,
			MaxHP:        def.HP,
			Speed:        def.Speed,
			Lifetime:     def.Lifetime,
			Behavior:     def.SpecialBehavior,
			Multiplier:   t.Multiplier,
			Quality:      string(t.Quality),
			QualityColor: t.QualityColor,
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
			InstanceID:   instanceID,
			DefID:        def.ID,
			Name:         def.Name,
			Type:         string(def.Type),
			X:            x,
			Y:            y,
			HP:           def.HP,
			MaxHP:        def.HP,
			Speed:        def.Speed,
			Lifetime:     def.Lifetime,
			Behavior:     def.SpecialBehavior,
			Multiplier:   t.Multiplier,
			Quality:      string(t.Quality),
			QualityColor: t.QualityColor,
		},
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

// handleUpgradeWeapon 處理武器升級（DAY-067）
func (g *Game) handleUpgradeWeapon(p *player.Player, msg *ws.Message) {
	var payload ws.UpgradeWeaponPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}
	if payload.WeaponLevel < 1 || payload.WeaponLevel > 3 {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "invalid_weapon", Message: "武器等級需在 1-3 之間"},
		})
		return
	}
	if !p.UpgradeWeapon(payload.WeaponLevel) {
		return
	}
	log.Printf("[Game] Player %s upgraded weapon to LV%d", p.ID, payload.WeaponLevel)
	g.sendPlayerUpdate(p)
}

// handleSetTitle 處理設定顯示稱號（DAY-068）
func (g *Game) handleSetTitle(p *player.Player, msg *ws.Message) {
	var payload ws.SetTitlePayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}
	if !p.SetTitle(achievement.TitleID(payload.TitleID)) {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: ws.ErrorPayload{Code: "title_not_unlocked", Message: "尚未解鎖此稱號"},
		})
		return
	}
	log.Printf("[Game] Player %s set title to %s", p.ID, payload.TitleID)
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
			// 稱號（DAY-068）
			TitleID:    snapshots[i].TitleID,
			TitleName:  snapshots[i].TitleName,
			TitleIcon:  snapshots[i].TitleIcon,
			TitleColor: snapshots[i].TitleColor,
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

// broadcastTournament 廣播週賽排名給所有玩家（每 30 秒，DAY-066）
// 每個玩家收到的訊息包含自己的排名和積分
func (g *Game) broadcastTournament() {
	snap := g.tournamentMgr.GetSnapshot()

	// 取得所有在線玩家 ID
	g.mu.RLock()
	playerIDs := make([]string, 0, len(g.Players))
	for id := range g.Players {
		playerIDs = append(playerIDs, id)
	}
	g.mu.RUnlock()

	// 轉換 RankEntry → TournamentRankEntry
	rankings := make([]ws.TournamentRankEntry, len(snap.Rankings))
	for i, r := range snap.Rankings {
		rankings[i] = ws.TournamentRankEntry{
			Rank:        r.Rank,
			PlayerID:    r.PlayerID,
			DisplayName: r.DisplayName,
			Points:      r.Points,
			Prize:       r.Prize,
			PrizeLabel:  r.PrizeLabel,
		}
	}

	// 對每個玩家個別發送（包含自己的排名）
	for _, pid := range playerIDs {
		rank, points := g.tournamentMgr.GetPlayerRank(pid)

		// 標記自己
		personalRankings := make([]ws.TournamentRankEntry, len(rankings))
		copy(personalRankings, rankings)
		for i := range personalRankings {
			personalRankings[i].IsSelf = (personalRankings[i].PlayerID == pid)
		}

		g.Hub.Send(pid, &ws.Message{
			Type: ws.MsgTournamentUpdate,
			Payload: ws.TournamentUpdatePayload{
				WeekStart:    snap.WeekStart,
				WeekEnd:      snap.WeekEnd,
				SecondsLeft:  snap.SecondsLeft,
				Rankings:     personalRankings,
				TotalPlayers: snap.TotalPlayers,
				PlayerRank:   rank,
				PlayerPoints: points,
			},
		})
	}
}

// GetTournamentSnapshot 取得週賽快照（供 HTTP 端點使用，DAY-066）
func (g *Game) GetTournamentSnapshot() tournament.Snapshot {
	return g.tournamentMgr.GetSnapshot()
}

// PlayerProfile 玩家個人資料（DAY-069）
type PlayerProfile struct {
	PlayerID    string                 `json:"player_id"`
	DisplayName string                 `json:"display_name"`
	Coins       int                    `json:"coins"`
	KillCount   int                    `json:"kill_count"`
	SessionScore int                   `json:"session_score"`
	MaxCoins    int                    `json:"max_coins"`
	WeaponLevel int                    `json:"weapon_level"`
	LoginStreak int                    `json:"login_streak"`
	MaxLoginStreak int                 `json:"max_login_streak"`
	// 稱號
	TitleID    string                  `json:"title_id"`
	TitleName  string                  `json:"title_name"`
	TitleIcon  string                  `json:"title_icon"`
	TitleColor string                  `json:"title_color"`
	// 成就
	AchievementCount int               `json:"achievement_count"`
	Achievements     []AchievementInfo `json:"achievements"`
	// 週賽
	TournamentPoints int               `json:"tournament_points"`
	TournamentRank   int               `json:"tournament_rank"`
	// 時間戳
	Timestamp int64                    `json:"timestamp"`
}

// AchievementInfo 成就資訊（用於 Profile）
type AchievementInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	UnlockedAt  int64  `json:"unlocked_at"` // Unix ms
}

// GetPlayerProfile 取得玩家個人資料（供 /profile HTTP 端點使用）
func (g *Game) GetPlayerProfile(playerID string) (*PlayerProfile, bool) {
	g.mu.RLock()
	p, ok := g.Players[playerID]
	g.mu.RUnlock()
	if !ok {
		return nil, false
	}

	snap := p.Snapshot()
	leaderSnap := p.LeaderboardSnapshot()

	// 成就列表
	achUnlocked := p.GetAchievements()
	loginStreak, maxLoginStreak := p.GetLoginInfo()

	achInfos := make([]AchievementInfo, 0, len(achUnlocked))
	for _, u := range achUnlocked {
		achInfos = append(achInfos, AchievementInfo{
			ID:         string(u.ID),
			Name:       u.Name,
			Icon:       u.Icon,
			UnlockedAt: u.UnlockedAt.UnixMilli(),
		})
	}

	// 週賽排名
	tournSnap := g.tournamentMgr.GetSnapshot()
	tournPoints := 0
	tournRank := 0
	for _, entry := range tournSnap.Rankings {
		if entry.PlayerID == playerID {
			tournPoints = entry.Points
			tournRank = entry.Rank
			break
		}
	}

	return &PlayerProfile{
		PlayerID:       playerID,
		DisplayName:    snap.DisplayName,
		Coins:          snap.Coins,
		KillCount:      snap.KillCount,
		SessionScore:   snap.SessionScore,
		MaxCoins:       leaderSnap.MaxCoins,
		WeaponLevel:    snap.WeaponLevel,
		LoginStreak:    loginStreak,
		MaxLoginStreak: maxLoginStreak,
		TitleID:        snap.TitleID,
		TitleName:      snap.TitleName,
		TitleIcon:      snap.TitleIcon,
		TitleColor:     snap.TitleColor,
		AchievementCount: len(achUnlocked),
		Achievements:   achInfos,
		TournamentPoints: tournPoints,
		TournamentRank:   tournRank,
		Timestamp:      time.Now().UnixMilli(),
	}, true
}

// GetAllPlayerProfiles 取得所有在線玩家的個人資料摘要（供 /profiles 端點）
func (g *Game) GetAllPlayerProfiles() []PlayerProfile {
	g.mu.RLock()
	playerIDs := make([]string, 0, len(g.Players))
	for id := range g.Players {
		playerIDs = append(playerIDs, id)
	}
	g.mu.RUnlock()

	profiles := make([]PlayerProfile, 0, len(playerIDs))
	for _, id := range playerIDs {
		if profile, ok := g.GetPlayerProfile(id); ok {
			profiles = append(profiles, *profile)
		}
	}
	return profiles
}

// GetPlayerCodexSnapshot 取得玩家圖鑑快照（供 /codex HTTP 端點使用，DAY-081）
func (g *Game) GetPlayerCodexSnapshot(playerID string) (interface{}, bool) {
	g.mu.RLock()
	p, ok := g.Players[playerID]
	g.mu.RUnlock()
	if !ok || p.Codex == nil {
		return nil, false
	}
	entries := p.Codex.GetSnapshot()
	unlocked, total := p.Codex.GetStats()
	return map[string]interface{}{
		"player_id":      playerID,
		"entries":        entries,
		"unlocked_count": unlocked,
		"total_count":    total,
		"is_complete":    p.Codex.IsComplete(),
	}, true
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

// sendAchievements 批次傳送成就解鎖通知，並檢查稱號解鎖（DAY-068）
func (g *Game) sendAchievements(playerID string, unlocks []*achievement.AchievementUnlock) {
	if len(unlocks) == 0 {
		return
	}
	g.mu.RLock()
	p, ok := g.Players[playerID]
	g.mu.RUnlock()

	for _, u := range unlocks {
		g.sendAchievement(playerID, u)
		// 檢查是否解鎖新稱號
		if ok {
			titleDef := p.OnAchievementUnlocked(u.ID)
			if titleDef != nil {
				log.Printf("[Title] Player %s unlocked title: %s (%s)", playerID, titleDef.Name, titleDef.ID)
				g.Hub.Send(playerID, &ws.Message{
					Type: ws.MsgTitleUnlocked,
					Payload: ws.TitleUnlockedPayload{
						TitleID:     string(titleDef.ID),
						TitleName:   titleDef.Name,
						TitleIcon:   titleDef.Icon,
						TitleColor:  titleDef.Color,
						Description: titleDef.Description,
					},
				})
			}
		}
	}
}

// ── 每日任務系統（DAY-037）已移至 mission_handler.go ──────────────────────────────────────

// ── 效能上報（DAY-045）已移至 perf_handler.go ──────────────────────────────────────

// ── Jackpot 系統（DAY-048）已移至 jackpot_handler.go ──────────────────────────────────────

// handleBuySkin 處理購買砲台外觀（DAY-071）
func (g *Game) handleBuySkin(p *player.Player, msg *ws.Message) {
	var payload ws.BuySkinPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	// 找到外觀定義
	var skinDef *ws.SkinDef
	for i := range ws.AvailableSkins {
		if ws.AvailableSkins[i].ID == payload.SkinID {
			skinDef = &ws.AvailableSkins[i]
			break
		}
	}
	if skinDef == nil {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: map[string]string{"message": "外觀不存在"},
		})
		return
	}

	// 嘗試購買
	if !p.BuySkin(skinDef.ID, skinDef.Price) {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: map[string]string{"message": "金幣不足或已擁有此外觀"},
		})
		return
	}

	// 購買成功：自動裝備
	p.EquipSkin(skinDef.ID)
	equippedSkin, ownedSkins := p.GetSkinInfo()

	// 通知玩家
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgSkinUpdate,
		Payload: ws.SkinUpdatePayload{
			PlayerID:     p.ID,
			EquippedSkin: equippedSkin,
			OwnedSkins:   ownedSkins,
			NewBalance:   p.Coins,
		},
	})

	// 更新玩家狀態（讓 Client 看到新金幣餘額）
	g.sendPlayerUpdate(p)

	log.Printf("[Skin] 玩家 %s 購買外觀 %s（%s），剩餘金幣 %d",
		p.ID, skinDef.ID, skinDef.Name, p.Coins)
}

// handleEquipSkin 處理裝備砲台外觀（DAY-071）
func (g *Game) handleEquipSkin(p *player.Player, msg *ws.Message) {
	var payload ws.EquipSkinPayload
	if err := remarshal(msg.Payload, &payload); err != nil {
		return
	}

	if !p.EquipSkin(payload.SkinID) {
		g.Hub.Send(p.ID, &ws.Message{
			Type:    ws.MsgError,
			Payload: map[string]string{"message": "未擁有此外觀"},
		})
		return
	}

	equippedSkin, ownedSkins := p.GetSkinInfo()
	g.Hub.Send(p.ID, &ws.Message{
		Type: ws.MsgSkinUpdate,
		Payload: ws.SkinUpdatePayload{
			PlayerID:     p.ID,
			EquippedSkin: equippedSkin,
			OwnedSkins:   ownedSkins,
			NewBalance:   0, // 裝備操作不改變金幣
		},
	})

	g.sendPlayerUpdate(p)
	log.Printf("[Skin] 玩家 %s 裝備外觀 %s", p.ID, payload.SkinID)
}
