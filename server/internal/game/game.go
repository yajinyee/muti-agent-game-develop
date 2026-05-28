// Package game — 遊戲主循環與狀態機
// server-core-agent 負責維護
package game

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	"chiikawa-game/internal/data"
	"chiikawa-game/internal/protocol"
	"chiikawa-game/internal/ws"
)

// GameState 遊戲狀態
type GameState string

const (
	StateNormalPlay  GameState = "normal_play"
	StateBossWarning GameState = "boss_warning"
	StateBossBattle  GameState = "boss_battle"
	StateBossResult  GameState = "boss_result"
	StateBonusReady  GameState = "bonus_ready"
	StateBonusGame   GameState = "bonus_game"
	StateBonusResult GameState = "bonus_result"
)

// Game 遊戲實例
type Game struct {
	mu      sync.RWMutex
	hub     *ws.Hub
	state   GameState
	players map[string]*Player
	targets map[string]*Target
	boss    *Target

	// 計時器
	spawnTimer    float64
	bossTimer     float64 // BOSS 戰剩餘時間
	bossWarningAt time.Time
	nextBossIn    float64 // 下次 BOSS 觸發倒數（秒）
	bonusPlayer   string  // 觸發 Bonus 的玩家 ID
	bonusTimer    float64
	bonusScore    map[string]int // playerID -> score
	bonusWeedHP   map[string]int // weedID -> hp（BG002 需要連點2次）
	bossPhase2    bool           // BOSS Phase 2 是否已觸發
	bossPhase3    bool           // BOSS Phase 3 是否已觸發（絕望模式，HP ≤ 20%）

	// 幸運特殊魚系統
	luckyChainLightning  *luckyChainLightningManager
	luckyCrabTorpedo     *luckyCrabTorpedoManager
	luckyVortex          *luckyVortexManager
	luckyGoldenDragon    *luckyGoldenDragonManager
	luckyThunderLobster  *luckyThunderLobsterManager
	luckyAwakenedPhoenix *luckyAwakenedPhoenixManager
	luckyShockwaveBomb   *luckyShockwaveBombManager
	luckyDrillTorpedo    *luckyDrillTorpedoManager
	luckyTimeFreeze      *luckyTimeFreezeManager
	luckyChainExplosion  *luckyChainExplosionManager
	// DAY-295 新增
	luckyChainLongKing  *luckyChainLongKingManager
	luckyDragonShotgun  *luckyDragonShotgunManager
	luckyRocketCannon   *luckyRocketCannonManager
	luckyDeepWhirlpool  *luckyDeepWhirlpoolManager
	luckyVampireMult    *luckyVampireMultManager
	// DAY-296 新增
	luckyMirrorFish   *luckyMirrorFishManager
	luckyGoldenRain   *luckyGoldenRainManager
	luckyFreezeBomb   *luckyFreezeBombManager
	luckyThunderStorm *luckyThunderStormManager
	luckyLuckyWheel   *luckyLuckyWheelManager

	// DAY-301 新增
	luckyJackpotFish *luckyJackpotFishManager
	luckyCoopFish    *luckyCoopFishManager
	luckyTimeWarp    *luckyTimeWarpManager

	// DAY-302 新增
	luckyChainMeteor *luckyChainMeteorManager

	// DAY-303 新增
	luckyCrashFish *luckyCrashFishManager

	// DAY-304 新增
	luckyElectricEel  *luckyElectricEelManager
	luckyAnglerFish   *luckyAnglerFishManager
	luckyBlackHole    *luckyBlackHoleManager
	luckyBountyHunter *luckyBountyHunterManager
	luckyTsunami      *luckyTsunamiManager

	// DAY-305 新增
	luckyDragonWrathV2 *luckyDragonWrathV2Manager
	luckyHumpbackWhale *luckyHumpbackWhaleManager
	luckyLegendDragon  *luckyLegendDragonManager
	luckyGuildWar      *luckyGuildWarManager
	luckyQualityFish   *luckyQualityFishManager

	// DAY-306 新增
	luckyTornado      *luckyTornadoManager
	luckyEarthquake   *luckyEarthquakeManager
	luckyVolcano      *luckyVolcanoManager
	luckyCosmicRay    *luckyCosmicRayManager
	luckyDivineDragon *luckyDivineDragonManager

	// DAY-307 新增
	luckyQuantum  *luckyQuantumManager
	luckySupernova *luckySupernovaManager
	luckyInfinite *luckyInfiniteManager
	luckyGenesis  *luckyGenesisManager
	luckyRebirth  *luckyRebirthManager

	// DAY-308 新增
	luckyAwakenedCroc *luckyAwakenedCrocManager
	luckyVampireV2    *luckyVampireV2Manager
	luckySuperAwaken  *luckySuperAwakenManager
	luckyGiantPrize   *luckyGiantPrizeManager
	luckyImmortalBoss *luckyImmortalBossManager

	// DAY-309 新增
	luckyIcePhoenix       *luckyIcePhoenixManager
	luckyDragonFury       *luckyDragonFuryManager
	luckyMultCascade      *luckyMultCascadeManager
	luckyAwakenBossV2     *luckyAwakenBossV2Manager
	luckyUltimateJudgment *luckyUltimateJudgmentManager

	// DAY-310 新增
	luckyComboBurst      *luckyComboBurstManager
	luckyTimeBomb        *luckyTimeBombManager
	luckyElementalFusion *luckyElementalFusionManager
	luckyTreasureHunter  *luckyTreasureHunterManager
	luckyMythAwaken      *luckyMythAwakenManager

	// DAY-312 新增
	luckyStarPortal    *luckyStarPortalManager
	luckyDragonSoul    *luckyDragonSoulManager
	luckySpacetimeRift *luckySpacetimeRiftManager
	luckyHolyJudgment  *luckyHolyJudgmentManager
	luckyBigBang       *luckyBigBangManager

	// DAY-313 新增
	luckyJackpotPool *luckyJackpotPoolManager

	// DAY-314 新增
	luckyMultiverse  *luckyMultiverseManager
	luckyTimeLoop    *luckyTimeLoopManager
	luckyFateWheel   *luckyFateWheelManager
	luckyDivineRealm *luckyDivineRealmManager
	luckyFinalPower  *luckyFinalPowerManager

	// DAY-315 新增
	luckyMutation    *luckyMutationManager
	luckyArcticStorm *luckyArcticStormManager
	luckyFisherWild  *luckyFisherWildManager
	luckyRiskLevel   *luckyRiskLevelManager
	luckyCosmicPulse *luckyCosmicPulseManager

	// DAY-316 新增
	luckyMirrorUniverse   *luckyMirrorUniverseManager
	luckyGravityField     *luckyGravityFieldManager
	luckyTimeAcceleration *luckyTimeAccelerationManager
	luckyNebulaVortex     *luckyNebulaVortexManager
	luckyCosmicJudgment   *luckyCosmicJudgmentManager

	// DAY-317 新增
	luckyPvpBattle       *luckyPvpBattleManager
	luckySkillChain      *luckySkillChainManager
	luckyGlobalExplosion *luckyGlobalExplosionManager
	luckySpacetimeFold   *luckySpacetimeFoldManager
	luckyCosmicEnd       *luckyCosmicEndManager

	lastTick time.Time
}

func NewGame(hub *ws.Hub) *Game {
	g := &Game{
		hub:        hub,
		state:      StateNormalPlay,
		players:    make(map[string]*Player),
		targets:    make(map[string]*Target),
		bonusScore: make(map[string]int),
		bonusWeedHP: make(map[string]int),
		bossPhase2: false,
		bossPhase3: false,
		lastTick:   time.Now(),

		// 初始化幸運特殊魚系統
		luckyChainLightning:  newLuckyChainLightningManager(),
		luckyCrabTorpedo:     newLuckyCrabTorpedoManager(),
		luckyVortex:          newLuckyVortexManager(),
		luckyGoldenDragon:    newLuckyGoldenDragonManager(),
		luckyThunderLobster:  newLuckyThunderLobsterManager(),
		luckyAwakenedPhoenix: newLuckyAwakenedPhoenixManager(),
		luckyShockwaveBomb:   newLuckyShockwaveBombManager(),
		luckyDrillTorpedo:    newLuckyDrillTorpedoManager(),
		luckyTimeFreeze:      newLuckyTimeFreezeManager(),
		luckyChainExplosion:  newLuckyChainExplosionManager(),
		// DAY-295 新增
		luckyChainLongKing:  newLuckyChainLongKingManager(),
		luckyDragonShotgun:  newLuckyDragonShotgunManager(),
		luckyRocketCannon:   newLuckyRocketCannonManager(),
		luckyDeepWhirlpool:  newLuckyDeepWhirlpoolManager(),
		luckyVampireMult:    newLuckyVampireMultManager(),
		// DAY-296 新增
		luckyMirrorFish:   newLuckyMirrorFishManager(),
		luckyGoldenRain:   newLuckyGoldenRainManager(),
		luckyFreezeBomb:   newLuckyFreezeBombManager(),
		luckyThunderStorm: newLuckyThunderStormManager(),
		luckyLuckyWheel:   newLuckyLuckyWheelManager(),

		// DAY-301 新增
		luckyJackpotFish: newLuckyJackpotFishManager(),
		luckyCoopFish:    newLuckyCoopFishManager(),
		luckyTimeWarp:    newLuckyTimeWarpManager(),

		// DAY-302 新增
		luckyChainMeteor: newLuckyChainMeteorManager(),

		// DAY-303 新增
		luckyCrashFish: newLuckyCrashFishManager(),

		// DAY-304 新增
		luckyElectricEel:  newLuckyElectricEelManager(),
		luckyAnglerFish:   newLuckyAnglerFishManager(),
		luckyBlackHole:    newLuckyBlackHoleManager(),
		luckyBountyHunter: newLuckyBountyHunterManager(),
		luckyTsunami:      newLuckyTsunamiManager(),

		// DAY-305 新增
		luckyDragonWrathV2: newLuckyDragonWrathV2Manager(),
		luckyHumpbackWhale: newLuckyHumpbackWhaleManager(),
		luckyLegendDragon:  newLuckyLegendDragonManager(),
		luckyGuildWar:      newLuckyGuildWarManager(),
		luckyQualityFish:   newLuckyQualityFishManager(),

		// DAY-306 新增
		luckyTornado:      newLuckyTornadoManager(),
		luckyEarthquake:   newLuckyEarthquakeManager(),
		luckyVolcano:      newLuckyVolcanoManager(),
		luckyCosmicRay:    newLuckyCosmicRayManager(),
		luckyDivineDragon: newLuckyDivineDragonManager(),

		// DAY-307 新增
		luckyQuantum:   newLuckyQuantumManager(),
		luckySupernova: newLuckySupernovaManager(),
		luckyInfinite:  newLuckyInfiniteManager(),
		luckyGenesis:   newLuckyGenesisManager(),
		luckyRebirth:   newLuckyRebirthManager(),

		// DAY-308 新增
		luckyAwakenedCroc: newLuckyAwakenedCrocManager(),
		luckyVampireV2:    newLuckyVampireV2Manager(),
		luckySuperAwaken:  newLuckySuperAwakenManager(),
		luckyGiantPrize:   newLuckyGiantPrizeManager(),
		luckyImmortalBoss: newLuckyImmortalBossManager(),

		// DAY-309 新增
		luckyIcePhoenix:       newLuckyIcePhoenixManager(),
		luckyDragonFury:       newLuckyDragonFuryManager(),
		luckyMultCascade:      newLuckyMultCascadeManager(),
		luckyAwakenBossV2:     newLuckyAwakenBossV2Manager(),
		luckyUltimateJudgment: newLuckyUltimateJudgmentManager(),

		// DAY-310 新增
		luckyComboBurst:      newLuckyComboBurstManager(),
		luckyTimeBomb:        newLuckyTimeBombManager(),
		luckyElementalFusion: newLuckyElementalFusionManager(),
		luckyTreasureHunter:  newLuckyTreasureHunterManager(),
		luckyMythAwaken:      newLuckyMythAwakenManager(),

		// DAY-312 新增
		luckyStarPortal:    newLuckyStarPortalManager(),
		luckyDragonSoul:    newLuckyDragonSoulManager(),
		luckySpacetimeRift: newLuckySpacetimeRiftManager(),
		luckyHolyJudgment:  newLuckyHolyJudgmentManager(),
		luckyBigBang:       newLuckyBigBangManager(),

		// DAY-313 新增
		luckyJackpotPool: newLuckyJackpotPoolManager(),

		// DAY-314 新增
		luckyMultiverse:  newLuckyMultiverseManager(),
		luckyTimeLoop:    newLuckyTimeLoopManager(),
		luckyFateWheel:   newLuckyFateWheelManager(),
		luckyDivineRealm: newLuckyDivineRealmManager(),
		luckyFinalPower:  newLuckyFinalPowerManager(),

		// DAY-315 新增
		luckyMutation:    newLuckyMutationManager(),
		luckyArcticStorm: newLuckyArcticStormManager(),
		luckyFisherWild:  newLuckyFisherWildManager(),
		luckyRiskLevel:   newLuckyRiskLevelManager(),
		luckyCosmicPulse: newLuckyCosmicPulseManager(),

		// DAY-316 新增
		luckyMirrorUniverse:   newLuckyMirrorUniverseManager(),
		luckyGravityField:     newLuckyGravityFieldManager(),
		luckyTimeAcceleration: newLuckyTimeAccelerationManager(),
		luckyNebulaVortex:     newLuckyNebulaVortexManager(),
		luckyCosmicJudgment:   newLuckyCosmicJudgmentManager(),

		// DAY-317 新增
		luckyPvpBattle:       newLuckyPvpBattleManager(),
		luckySkillChain:      newLuckySkillChainManager(),
		luckyGlobalExplosion: newLuckyGlobalExplosionManager(),
		luckySpacetimeFold:   newLuckySpacetimeFoldManager(),
		luckyCosmicEnd:       newLuckyCosmicEndManager(),
	}
	g.nextBossIn = 180 + rand.Float64()*120 // 3-5 分鐘
	return g
}

// Start 啟動遊戲循環
func (g *Game) Start() {
	go g.loop()
}

// AddPlayer 玩家加入
func (g *Game) AddPlayer(id string) {
	g.mu.Lock()
	p := NewPlayer(id)
	g.players[id] = p
	g.mu.Unlock()

	log.Printf("[Game] Player joined: %s", id)
	g.sendPlayerUpdate(id)
	g.hub.Send(id, protocol.MsgGameState, protocol.GameStatePayload{
		State:     string(g.state),
		Timestamp: time.Now().UnixMilli(),
	})
	// 把現有目標物同步給新玩家
	g.mu.RLock()
	for _, t := range g.targets {
		g.hub.Send(id, protocol.MsgTargetSpawn, g.targetSpawnPayload(t))
	}
	if g.boss != nil {
		g.hub.Send(id, protocol.MsgTargetSpawn, g.targetSpawnPayload(g.boss))
	}
	g.mu.RUnlock()
}

// RemovePlayer 玩家離開
func (g *Game) RemovePlayer(id string) {
	g.mu.Lock()
	delete(g.players, id)
	g.mu.Unlock()
	log.Printf("[Game] Player left: %s", id)
}

// HandleMessage 處理玩家訊息
func (g *Game) HandleMessage(clientID string, msgType string, payload json.RawMessage) {
	switch msgType {
	case protocol.MsgAttack:
		var req protocol.AttackRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			return
		}
		g.handleAttack(clientID, req)
	case protocol.MsgAutoToggle:
		g.handleAutoToggle(clientID)
	case protocol.MsgBetChange:
		var req protocol.BetChangeRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			return
		}
		g.handleBetChange(clientID, req.BetLevel)
	case protocol.MsgLock:
		var req protocol.LockRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			return
		}
		g.handleLock(clientID, req.TargetID)
	case protocol.MsgBonusClick:
		var req protocol.BonusClickRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			return
		}
		g.handleBonusClick(clientID, req)
	case protocol.MsgPing:
		g.hub.Send(clientID, protocol.MsgPong, map[string]interface{}{})
	case protocol.MsgTriggerBoss:
		g.triggerBoss()
	case protocol.MsgTriggerBonus:
		g.triggerBonus(clientID)
	case protocol.MsgCollectGoldenCoin:
		var req protocol.CollectGoldenCoinRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			return
		}
		g.collectGoldenCoin(clientID, req.CoinID)
	case protocol.MsgSetDisplayName:
		var req protocol.SetDisplayNameRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			return
		}
		g.handleSetDisplayName(clientID, req.Name)
	case "crash_harvest":
		// T130 崩潰魚：玩家點擊收割
		g.handleCrashHarvest(clientID)
	}
}

// ── 主循環 ────────────────────────────────────────────────────

func (g *Game) loop() {
	ticker := time.NewTicker(50 * time.Millisecond) // 20 ticks/sec
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		delta := now.Sub(g.lastTick).Seconds()
		g.lastTick = now
		g.tick(delta)
	}
}

func (g *Game) tick(delta float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	switch g.state {
	case StateNormalPlay:
		g.tickNormal(delta)
	case StateBossWarning:
		// 警告 3 秒後進入 BOSS 戰
		if time.Since(g.bossWarningAt).Seconds() >= 3 {
			g.spawnBoss()
		}
	case StateBossBattle:
		g.tickBoss(delta)
	case StateBonusGame:
		g.tickBonus(delta)
	}
}

func (g *Game) tickNormal(delta float64) {
	// 目標物生成
	g.spawnTimer += delta
	if g.spawnTimer >= SpawnInterval {
		g.spawnTimer = 0
		if len(g.targets) < MaxTargets {
			g.spawnTarget()
		}
	}

	// 移除超時目標物
	for id, t := range g.targets {
		if t.IsExpired() {
			delete(g.targets, id)
			g.hub.Broadcast(protocol.MsgTargetKill, protocol.TargetKillPayload{
				InstanceID: t.InstanceID,
				DefID:      t.Def.ID,
				Multiplier: t.Multiplier,
				Reward:     0,
				LaborGain:  0,
				KillerID:   "",
			})
		}
	}

	// BOSS 觸發計時
	g.nextBossIn -= delta
	if g.nextBossIn <= 0 && len(g.players) > 0 {
		g.startBossWarning()
	}

	// AUTO 射擊
	for _, p := range g.players {
		if p.IsAuto {
			g.autoFire(p)
		}
	}
}

func (g *Game) tickBoss(delta float64) {
	if g.boss == nil {
		g.setState(StateNormalPlay)
		return
	}
	g.bossTimer -= delta
	if g.bossTimer <= 0 {
		// BOSS 超時
		g.boss = nil
		g.setState(StateNormalPlay)
		g.hub.Broadcast(protocol.MsgBossEvent, protocol.BossEventPayload{
			Event: "timeout",
		})
		g.nextBossIn = 180 + rand.Float64()*120
	}
}

func (g *Game) tickBonus(delta float64) {
	g.bonusTimer -= delta
	g.hub.Broadcast(protocol.MsgBonusEvent, protocol.BonusEventPayload{
		Event:    "tick",
		TimeLeft: g.bonusTimer,
	})
	if g.bonusTimer <= 0 {
		g.endBonus()
	}
}

// ── 目標物生成 ────────────────────────────────────────────────

func (g *Game) spawnTarget() {
	// 取平均 BetLevel 決定難度
	avgBet := 1
	if len(g.players) > 0 {
		total := 0
		for _, p := range g.players {
			total += p.BetLevel
		}
		avgBet = total / len(g.players)
	}

	def := pickTargetDef(avgBet)
	t := NewTarget(def, SpawnX, spawnY())
	g.targets[t.InstanceID] = t
	g.hub.Broadcast(protocol.MsgTargetSpawn, g.targetSpawnPayload(t))
}

func (g *Game) targetSpawnPayload(t *Target) protocol.TargetSpawnPayload {
	return protocol.TargetSpawnPayload{
		InstanceID: t.InstanceID,
		DefID:      t.Def.ID,
		Name:       t.Def.Name,
		Type:       string(t.Def.Type),
		X:          t.X,
		Y:          t.Y,
		HP:         t.HP,
		MaxHP:      t.MaxHP,
		Speed:      t.Def.Speed,
		Lifetime:   t.Def.Lifetime,
		Behavior:   string(t.Def.Behavior),
		Multiplier: t.Multiplier,
	}
}

// ── 攻擊處理 ──────────────────────────────────────────────────

func (g *Game) handleAttack(playerID string, req protocol.AttackRequest) {
	g.mu.Lock()
	defer g.mu.Unlock()

	p, ok := g.players[playerID]
	if !ok {
		return
	}

	bet := p.GetBetDef()
	if !p.SpendCoins(bet.BetCost) {
		g.hub.Send(playerID, protocol.MsgError, protocol.ErrorPayload{
			Code:    "insufficient_coins",
			Message: "Not enough coins",
		})
		return
	}

	// Jackpot 貢獻（每次下注）
	g.luckyJackpotFish.ContributeBet(bet.BetCost)
	// DAY-313 Progressive Jackpot Pool 貢獻
	g.luckyJackpotPool.onShot(g, playerID, bet.BetCost)

	// 找目標
	var t *Target
	if req.TargetID != "" {
		t = g.targets[req.TargetID]
		if t == nil && g.boss != nil && g.boss.InstanceID == req.TargetID {
			t = g.boss
		}
	}

	if t == nil {
		// 未命中
		g.hub.Send(playerID, protocol.MsgAttackResult, protocol.AttackResultPayload{
			TargetID:    req.TargetID,
			IsHit:       false,
			CharacterID: p.GetCharacterID(),
			PosX:        req.ClickX,
			PosY:        req.ClickY,
		})
		return
	}

	// 命中
	isKill := t.TryKill(bet.BetCost)
	result := protocol.AttackResultPayload{
		TargetID:    t.InstanceID,
		IsHit:       true,
		IsKill:      isKill,
		Damage:      bet.AttackPower,
		CharacterID: p.GetCharacterID(),
		Multiplier:  t.Multiplier,
		PosX:        t.X,
		PosY:        t.Y,
	}

	if isKill {
		// Combo 系統：命中時增加 Combo
		_, _, comboBonus := p.AddCombo()
		effectiveMult := t.Multiplier * (1.0 + comboBonus)

		// DAY-301 全服加成倍率疊加
		jackpotBoost := g.luckyJackpotFish.getGrandBoostMult()
		coopBoost := g.luckyCoopFish.getCoopBoostMult()
		warpDmgMult := g.luckyTimeWarp.getTimeWarpDamageMult()
		if jackpotBoost > 1.0 {
			effectiveMult *= jackpotBoost
		}
		if coopBoost > 1.0 {
			effectiveMult *= coopBoost
		}
		if warpDmgMult > 1.0 {
			effectiveMult *= warpDmgMult
		}
		// DAY-302 連鎖隕石完美加成
		chainMeteorBoost := g.luckyChainMeteor.getChainMeteorPerfectMult()
		if chainMeteorBoost > 1.0 {
			effectiveMult *= chainMeteorBoost
		}
		// DAY-303 崩潰魚完美收割加成
		crashPerfectBoost := g.luckyCrashFish.getCrashPerfectMult()
		if crashPerfectBoost > 1.0 {
			effectiveMult *= crashPerfectBoost
		}
		// DAY-304 新增全服加成
		eelSuperBoost := g.luckyElectricEel.getEelSuperMult()
		if eelSuperBoost > 1.0 {
			effectiveMult *= eelSuperBoost
		}
		anglerPerfectBoost := g.luckyAnglerFish.getAnglerPerfectMult()
		if anglerPerfectBoost > 1.0 {
			effectiveMult *= anglerPerfectBoost
		}
		blackHoleSingularityBoost := g.luckyBlackHole.getBlackHoleSingularityMult()
		if blackHoleSingularityBoost > 1.0 {
			effectiveMult *= blackHoleSingularityBoost
		}
		bountyPerfectBoost := g.luckyBountyHunter.getBountyPerfectMult()
		if bountyPerfectBoost > 1.0 {
			effectiveMult *= bountyPerfectBoost
		}
		tsunamiPerfectBoost := g.luckyTsunami.getTsunamiPerfectMult()
		if tsunamiPerfectBoost > 1.0 {
			effectiveMult *= tsunamiPerfectBoost
		}
		// DAY-305 新增全服加成
		dragonWrathV2Boost := g.luckyDragonWrathV2.getDragonWrathV2PerfectMult()
		if dragonWrathV2Boost > 1.0 {
			effectiveMult *= dragonWrathV2Boost
		}
		whaleSongBoost := g.luckyHumpbackWhale.getWhaleSongMult()
		if whaleSongBoost > 1.0 {
			effectiveMult *= whaleSongBoost
		}
		legendDragonBoost := g.luckyLegendDragon.getLegendDragonRageMult()
		if legendDragonBoost > 1.0 {
			effectiveMult *= legendDragonBoost
		}
		guildWarBoost := g.luckyGuildWar.getGuildWarVictoryMult()
		if guildWarBoost > 1.0 {
			effectiveMult *= guildWarBoost
		}
		qualityLegendaryBoost := g.luckyQualityFish.getQualityLegendaryMult()
		if qualityLegendaryBoost > 1.0 {
			effectiveMult *= qualityLegendaryBoost
		}
		// DAY-306 新增全服加成
		tornadoPerfectBoost := g.luckyTornado.getTornadoPerfectMult()
		if tornadoPerfectBoost > 1.0 {
			effectiveMult *= tornadoPerfectBoost
		}
		earthquakePerfectBoost := g.luckyEarthquake.getEarthquakePerfectMult()
		if earthquakePerfectBoost > 1.0 {
			effectiveMult *= earthquakePerfectBoost
		}
		volcanoPerfectBoost := g.luckyVolcano.getVolcanoPerfectMult()
		if volcanoPerfectBoost > 1.0 {
			effectiveMult *= volcanoPerfectBoost
		}
		cosmicRayPerfectBoost := g.luckyCosmicRay.getCosmicRayPerfectMult()
		if cosmicRayPerfectBoost > 1.0 {
			effectiveMult *= cosmicRayPerfectBoost
		}
		divineDragonPerfectBoost := g.luckyDivineDragon.getDivineDragonPerfectMult()
		if divineDragonPerfectBoost > 1.0 {
			effectiveMult *= divineDragonPerfectBoost
		}
		// DAY-307 新增全服加成
		quantumCollapseBoost := g.luckyQuantum.getQuantumPerfectMult()
		if quantumCollapseBoost > 1.0 {
			effectiveMult *= quantumCollapseBoost
		}
		supernovaPerfectBoost := g.luckySupernova.getSupernovaPerfectMult()
		if supernovaPerfectBoost > 1.0 {
			effectiveMult *= supernovaPerfectBoost
		}
		supernovaMultBoost := g.luckySupernova.getSupernovaMultBoost()
		if supernovaMultBoost > 1.0 {
			effectiveMult *= supernovaMultBoost
		}
		infinitePerfectBoost := g.luckyInfinite.getInfinitePerfectMult()
		if infinitePerfectBoost > 1.0 {
			effectiveMult *= infinitePerfectBoost
		}
		genesisPerfectBoost := g.luckyGenesis.getGenesisPerfectMult()
		if genesisPerfectBoost > 1.0 {
			effectiveMult *= genesisPerfectBoost
		}
		rebirthPerfectBoost := g.luckyRebirth.getRebirthPerfectMult()
		if rebirthPerfectBoost > 1.0 {
			effectiveMult *= rebirthPerfectBoost
		}
		// 重生魚：重生目標擊破倍率加成
		rebirthKillMult := g.luckyRebirth.getRebirthKillMult(t.InstanceID)
		if rebirthKillMult > 1.0 {
			effectiveMult *= rebirthKillMult
		}
		// DAY-308 新增全服加成
		awakenedCrocBoost := g.luckyAwakenedCroc.getCrocPerfectMult()
		if awakenedCrocBoost > 1.0 {
			effectiveMult *= awakenedCrocBoost
		}
		vampireV2Boost := g.luckyVampireV2.getVampireV2PerfectMult()
		if vampireV2Boost > 1.0 {
			effectiveMult *= vampireV2Boost
		}
		vampireV2KillMult := g.luckyVampireV2.getVampireV2KillMult(playerID)
		if vampireV2KillMult > 1.0 {
			effectiveMult *= vampireV2KillMult
		}
		superAwakenBoost := g.luckySuperAwaken.getSuperAwakenPerfectMult()
		if superAwakenBoost > 1.0 {
			effectiveMult *= superAwakenBoost
		}
		giantPrizeBoost := g.luckyGiantPrize.getGiantPrizePerfectMult()
		if giantPrizeBoost > 1.0 {
			effectiveMult *= giantPrizeBoost
		}
		immortalBossBoost := g.luckyImmortalBoss.getImmortalBossPerfectMult()
		if immortalBossBoost > 1.0 {
			effectiveMult *= immortalBossBoost
		}
		// DAY-309 新增全服加成
		icePhoenixBoost := g.luckyIcePhoenix.getIcePhoenixPerfectMult()
		if icePhoenixBoost > 1.0 {
			effectiveMult *= icePhoenixBoost
		}
		dragonFuryBoost := g.luckyDragonFury.getDragonFuryPerfectMult()
		if dragonFuryBoost > 1.0 {
			effectiveMult *= dragonFuryBoost
		}
		multCascadeBoost := g.luckyMultCascade.getMultCascadePerfectMult()
		if multCascadeBoost > 1.0 {
			effectiveMult *= multCascadeBoost
		}
		// 倍率瀑布：個人擊破倍率加成
		multCascadeKillBonus := g.luckyMultCascade.getMultCascadeKillBonus(playerID)
		if multCascadeKillBonus > 1.0 {
			effectiveMult *= multCascadeKillBonus
		}
		awakenBossV2Boost := g.luckyAwakenBossV2.getAwakenBossV2PerfectMult()
		if awakenBossV2Boost > 1.0 {
			effectiveMult *= awakenBossV2Boost
		}
		ultimateJudgmentBoost := g.luckyUltimateJudgment.getUltimateJudgmentPerfectMult()
		if ultimateJudgmentBoost > 1.0 {
			effectiveMult *= ultimateJudgmentBoost
		}
		// DAY-310 新增全服加成
		comboBurstPerfectBoost := g.luckyComboBurst.getComboBurstPerfectMult()
		if comboBurstPerfectBoost > 1.0 {
			effectiveMult *= comboBurstPerfectBoost
		}
		// 連擊爆發：個人連擊倍率加成
		comboBurstKillMult := g.luckyComboBurst.onKillDuringComboBurst(g, p)
		if comboBurstKillMult > 1.0 {
			effectiveMult *= comboBurstKillMult
		}
		timeBombPerfectBoost := g.luckyTimeBomb.getTimeBombPerfectMult()
		if timeBombPerfectBoost > 1.0 {
			effectiveMult *= timeBombPerfectBoost
		}
		// 時間炸彈：擊破時增加能量
		g.luckyTimeBomb.onKillDuringTimeBomb(g, p)
		elementalFusionPerfectBoost := g.luckyElementalFusion.getElementalFusionPerfectMult()
		if elementalFusionPerfectBoost > 1.0 {
			effectiveMult *= elementalFusionPerfectBoost
		}
		treasureHunterPerfectBoost := g.luckyTreasureHunter.getTreasureHunterPerfectMult()
		if treasureHunterPerfectBoost > 1.0 {
			effectiveMult *= treasureHunterPerfectBoost
		}
		// 寶藏獵人：寶藏目標倍率加成
		treasureMult := g.luckyTreasureHunter.getTreasureMult(playerID, t.InstanceID)
		if treasureMult > 1.0 {
			effectiveMult *= treasureMult
			g.luckyTreasureHunter.onTreasureKilled(g, p, t.InstanceID)
		}
		mythAwakenPerfectBoost := g.luckyMythAwaken.getMythAwakenPerfectMult()
		if mythAwakenPerfectBoost > 1.0 {
			effectiveMult *= mythAwakenPerfectBoost
		}
		// 神話覺醒：全場目標倍率 ×3.0
		mythMult := g.luckyMythAwaken.getMythMult()
		if mythMult > 1.0 {
			effectiveMult *= mythMult
			g.luckyMythAwaken.onKillDuringMyth(playerID)
		}
		// DAY-312 新增全服加成
		starPortalPerfectBoost := g.luckyStarPortal.getStarPortalPerfectMult()
		if starPortalPerfectBoost > 1.0 {
			effectiveMult *= starPortalPerfectBoost
		}
		dragonSoulPerfectBoost := g.luckyDragonSoul.getDragonSoulPerfectMult()
		if dragonSoulPerfectBoost > 1.0 {
			effectiveMult *= dragonSoulPerfectBoost
		}
		spacetimeRiftPerfectBoost := g.luckySpacetimeRift.getSpacetimeRiftPerfectMult()
		if spacetimeRiftPerfectBoost > 1.0 {
			effectiveMult *= spacetimeRiftPerfectBoost
		}
		holyJudgmentPerfectBoost := g.luckyHolyJudgment.getHolyJudgmentPerfectMult()
		if holyJudgmentPerfectBoost > 1.0 {
			effectiveMult *= holyJudgmentPerfectBoost
		}
		bigBangPerfectBoost := g.luckyBigBang.getBigBangPerfectMult()
		if bigBangPerfectBoost > 1.0 {
			effectiveMult *= bigBangPerfectBoost
		}
		// DAY-314 新增全服加成
		multiversePerfectBoost := g.luckyMultiverse.getMultiversePerfectMult()
		if multiversePerfectBoost > 1.0 {
			effectiveMult *= multiversePerfectBoost
		}
		timeLoopPerfectBoost := g.luckyTimeLoop.getTimeLoopPerfectMult()
		if timeLoopPerfectBoost > 1.0 {
			effectiveMult *= timeLoopPerfectBoost
		}
		// 時間迴圈：個人擊破倍率加成
		timeLoopKillMult := g.luckyTimeLoop.getTimeLoopMult()
		if timeLoopKillMult > 1.0 {
			effectiveMult *= timeLoopKillMult
		}
		fateWheelPerfectBoost := g.luckyFateWheel.getFateWheelPerfectMult()
		if fateWheelPerfectBoost > 1.0 {
			effectiveMult *= fateWheelPerfectBoost
		}
		divineRealmPerfectBoost := g.luckyDivineRealm.getDivineRealmPerfectMult()
		if divineRealmPerfectBoost > 1.0 {
			effectiveMult *= divineRealmPerfectBoost
		}
		finalPowerBoost := g.luckyFinalPower.getFinalPowerMult()
		if finalPowerBoost > 1.0 {
			effectiveMult *= finalPowerBoost
		}
		// DAY-315 新增全服加成
		mutationBoost := g.luckyMutation.getMutationMult()
		if mutationBoost > 1.0 {
			effectiveMult *= mutationBoost
		}
		arcticStormBoost := g.luckyArcticStorm.getArcticStormMult()
		if arcticStormBoost > 1.0 {
			effectiveMult *= arcticStormBoost
		}
		fisherWildBoost := g.luckyFisherWild.getFisherWildMult()
		if fisherWildBoost > 1.0 {
			effectiveMult *= fisherWildBoost
		}
		riskLevelBoost := g.luckyRiskLevel.getRiskLevelMult()
		if riskLevelBoost > 1.0 {
			effectiveMult *= riskLevelBoost
		}
		cosmicPulseBoost := g.luckyCosmicPulse.getCosmicPulseMult()
		if cosmicPulseBoost > 1.0 {
			effectiveMult *= cosmicPulseBoost
		}
		// DAY-316 新增全服加成
		mirrorUniverseBoost := g.luckyMirrorUniverse.getMirrorUniverseMult()
		if mirrorUniverseBoost > 1.0 {
			effectiveMult *= mirrorUniverseBoost
		}
		gravityFieldBoost := g.luckyGravityField.getGravityFieldMult()
		if gravityFieldBoost > 1.0 {
			effectiveMult *= gravityFieldBoost
		}
		timeAccelerationBoost := g.luckyTimeAcceleration.getTimeAccelerationMult()
		if timeAccelerationBoost > 1.0 {
			effectiveMult *= timeAccelerationBoost
		}
		nebulaVortexBoost := g.luckyNebulaVortex.getNebulaVortexMult()
		if nebulaVortexBoost > 1.0 {
			effectiveMult *= nebulaVortexBoost
		}
		cosmicJudgmentBoost := g.luckyCosmicJudgment.getCosmicJudgmentMult()
		if cosmicJudgmentBoost > 1.0 {
			effectiveMult *= cosmicJudgmentBoost
		}
		// DAY-317 新增全服加成
		pvpBattleBoost := g.luckyPvpBattle.getPvpBattleMult()
		if pvpBattleBoost > 1.0 {
			effectiveMult *= pvpBattleBoost
		}
		skillChainBoost := g.luckySkillChain.getSkillChainMult()
		if skillChainBoost > 1.0 {
			effectiveMult *= skillChainBoost
		}
		globalExplosionBoost := g.luckyGlobalExplosion.getGlobalExplosionMult()
		if globalExplosionBoost > 1.0 {
			effectiveMult *= globalExplosionBoost
		}
		spacetimeFoldBoost := g.luckySpacetimeFold.getSpacetimeFoldMult()
		if spacetimeFoldBoost > 1.0 {
			effectiveMult *= spacetimeFoldBoost
		}
		cosmicEndBoost := g.luckyCosmicEnd.getCosmicEndMult()
		if cosmicEndBoost > 1.0 {
			effectiveMult *= cosmicEndBoost
		}
		// 龍捲風期間擊破計數
		if g.luckyTornado.isTornadoActive() {
			g.luckyTornado.notifyTornadoKill(g, playerID)
		}

		reward := int(float64(bet.BetCost) * effectiveMult)
		result.Reward = reward
		result.LaborGain = t.Def.LaborGain

		p.AddCoins(reward)
		laborFull := p.AddLabor(t.Def.LaborGain)

		// 廣播擊破
		if t.Def.Type == data.TypeBoss {
			// BOSS Phase 2：HP 降到 50% 以下時觸發
			if !g.bossPhase2 && t.HPPercent() < 0.5 {
				g.bossPhase2 = true
				g.hub.Broadcast(protocol.MsgBossEvent, protocol.BossEventPayload{
					Event:      "phase_change",
					InstanceID: t.InstanceID,
					HP:         t.HP,
					MaxHP:      t.MaxHP,
				})
				log.Printf("[Game] Boss Phase 2 triggered! HP: %d/%d", t.HP, t.MaxHP)
			}
			// BOSS Phase 3：HP 降到 20% 以下時觸發（絕望模式）
			if !g.bossPhase3 && t.HPPercent() < 0.2 {
				g.bossPhase3 = true
				g.hub.Broadcast(protocol.MsgBossEvent, protocol.BossEventPayload{
					Event:      "phase_change",
					Phase:      3,
					InstanceID: t.InstanceID,
					HP:         t.HP,
					MaxHP:      t.MaxHP,
				})
				log.Printf("[Game] Boss Phase 3 triggered (絕望模式)! HP: %d/%d", t.HP, t.MaxHP)
			}
			g.handleBossKill(playerID, p, t, reward)
		} else {
			delete(g.targets, t.InstanceID)
			g.hub.Broadcast(protocol.MsgTargetKill, protocol.TargetKillPayload{
				InstanceID: t.InstanceID,
				DefID:      t.Def.ID,
				Multiplier: t.Multiplier,
				Reward:     reward,
				LaborGain:  t.Def.LaborGain,
				KillerID:   playerID,
			})

			// 幸運特殊魚觸發
			killerName := p.GetDisplayName() // 使用顯示名稱
			switch {
			case isLuckyChainLightningFish(t.Def.ID):
				g.tryLuckyChainLightning(playerID, killerName)
			case isLuckyCrabTorpedoFish(t.Def.ID):
				g.tryLuckyCrabTorpedo(playerID, killerName)
			case isLuckyVortexFish(t.Def.ID):
				g.tryLuckyVortex(playerID, killerName)
			case isLuckyGoldenDragonFish(t.Def.ID):
				g.tryLuckyGoldenDragon(playerID, killerName)
			case isLuckyThunderLobsterFish(t.Def.ID):
				g.tryLuckyThunderLobster(playerID, killerName)
			case isLuckyAwakenedPhoenixFish(t.Def.ID):
				g.tryLuckyAwakenedPhoenix(playerID, killerName)
			case isLuckyShockwaveBombFish(t.Def.ID):
				g.tryLuckyShockwaveBomb(playerID, killerName)
			case isLuckyDrillTorpedoFish(t.Def.ID):
				g.tryLuckyDrillTorpedo(playerID, killerName)
			case isLuckyTimeFreezeFish(t.Def.ID):
				g.tryLuckyTimeFreeze(playerID, killerName)
			case isLuckyChainExplosionFish(t.Def.ID):
				g.tryLuckyChainExplosion(playerID, killerName)
			// DAY-295 新增
			case isLuckyChainLongKingFish(t.Def.ID):
				g.luckyChainLongKing.tryLuckyChainLongKing(g, playerID, killerName)
			case isLuckyDragonShotgunFish(t.Def.ID):
				g.luckyDragonShotgun.tryLuckyDragonShotgun(g, playerID, killerName)
			case isLuckyRocketCannonFish(t.Def.ID):
				g.luckyRocketCannon.tryLuckyRocketCannon(g, playerID, killerName)
			case isLuckyDeepWhirlpoolFish(t.Def.ID):
				g.luckyDeepWhirlpool.tryLuckyDeepWhirlpool(g, playerID, killerName)
			case isLuckyVampireMultFish(t.Def.ID):
				g.luckyVampireMult.tryLuckyVampireMult(g, playerID, killerName)
			// DAY-296 新增
			case isLuckyMirrorFish(t.Def.ID):
				g.tryLuckyMirrorFish(playerID, killerName)
			case isLuckyGoldenRainFish(t.Def.ID):
				g.tryLuckyGoldenRain(playerID, killerName)
			case isLuckyFreezeBombFish(t.Def.ID):
				g.tryLuckyFreezeBomb(playerID, killerName)
			case isLuckyThunderStormFish(t.Def.ID):
				g.tryLuckyThunderStorm(playerID, killerName)
			case isLuckyLuckyWheelFish(t.Def.ID):
				g.tryLuckyLuckyWheel(playerID, killerName)
			// DAY-301 新增
			case isLuckyJackpotFish(t.Def.ID):
				g.tryLuckyJackpotFish(playerID, killerName)
			case isLuckyCoopFish(t.Def.ID):
				g.tryLuckyCoopFish(playerID, killerName)
			case isLuckyTimeWarpFish(t.Def.ID):
				g.tryLuckyTimeWarp(playerID, killerName)
			// DAY-302 新增
			case isLuckyChainMeteorFish(t.Def.ID):
				g.tryLuckyChainMeteor(playerID, killerName)
			// DAY-303 新增
			case isLuckyCrashFish(t.Def.ID):
				g.tryLuckyCrashFish(playerID, killerName)
			// DAY-304 新增
			case isLuckyElectricEelFish(t.Def.ID):
				g.tryLuckyElectricEelFish(playerID, killerName)
			case isLuckyAnglerFish(t.Def.ID):
				g.tryLuckyAnglerFish(playerID, killerName)
			case isLuckyBlackHoleFish(t.Def.ID):
				g.tryLuckyBlackHoleFish(playerID, killerName)
			case isLuckyBountyHunterFish(t.Def.ID):
				g.tryLuckyBountyHunterFish(playerID, killerName)
			case isLuckyTsunamiFish(t.Def.ID):
				g.tryLuckyTsunamiFish(playerID, killerName)
			// DAY-305 新增
			case isLuckyDragonWrathV2Fish(t.Def.ID):
				g.tryLuckyDragonWrathV2Fish(playerID, killerName)
			case isLuckyHumpbackWhaleFish(t.Def.ID):
				g.tryLuckyHumpbackWhaleFish(playerID, killerName)
			case isLuckyLegendDragonFish(t.Def.ID):
				g.tryLuckyLegendDragonFish(playerID, killerName)
			case isLuckyGuildWarFish(t.Def.ID):
				g.tryLuckyGuildWarFish(playerID, killerName)
			case isLuckyQualityFish(t.Def.ID):
				g.tryLuckyQualityFish(playerID, killerName)
			// DAY-306 新增
			case isLuckyTornadoFish(t.Def.ID):
				g.luckyTornado.tryLuckyTornadoFish(g, playerID, killerName)
			case isLuckyEarthquakeFish(t.Def.ID):
				g.luckyEarthquake.tryLuckyEarthquakeFish(g, playerID, killerName)
			case isLuckyVolcanoFish(t.Def.ID):
				g.luckyVolcano.tryLuckyVolcanoFish(g, playerID, killerName)
			case isLuckyCosmicRayFish(t.Def.ID):
				g.luckyCosmicRay.tryLuckyCosmicRayFish(g, playerID, killerName)
			case isLuckyDivineDragonFish(t.Def.ID):
				g.luckyDivineDragon.tryLuckyDivineDragonFish(g, playerID, killerName)
			// DAY-307 新增
			case isLuckyQuantumFish(t.Def.ID):
				g.luckyQuantum.tryLuckyQuantumFish(g, playerID, killerName)
			case isLuckySupernovaFish(t.Def.ID):
				g.luckySupernova.tryLuckySupernovaFish(g, playerID, killerName)
			case isLuckyInfiniteFish(t.Def.ID):
				g.luckyInfinite.tryLuckyInfiniteFish(g, playerID, killerName)
			case isLuckyGenesisFish(t.Def.ID):
				g.luckyGenesis.tryLuckyGenesisFish(g, playerID, killerName)
			case isLuckyRebirthFish(t.Def.ID):
				g.luckyRebirth.tryLuckyRebirthFish(g, playerID, killerName)
			// DAY-308 新增
			case isLuckyAwakenedCrocFish(t.Def.ID):
				g.luckyAwakenedCroc.tryLuckyAwakenedCrocFish(g, p)
			case isLuckyVampireV2Fish(t.Def.ID):
				g.luckyVampireV2.tryLuckyVampireV2Fish(g, p)
			case isLuckySuperAwakenFish(t.Def.ID):
				g.luckySuperAwaken.tryLuckySuperAwakenFish(g, p)
			case isLuckyGiantPrizeFish(t.Def.ID):
				g.luckyGiantPrize.tryLuckyGiantPrizeFish(g, p)
			case isLuckyImmortalBossFish(t.Def.ID):
				g.luckyImmortalBoss.tryLuckyImmortalBossFish(g, p)
			// DAY-309 新增
			case isLuckyIcePhoenixFish(t.Def.ID):
				g.luckyIcePhoenix.tryLuckyIcePhoenixFish(g, p)
			case isLuckyDragonFuryFish(t.Def.ID):
				g.luckyDragonFury.tryLuckyDragonFuryFish(g, p)
			case isLuckyMultCascadeFish(t.Def.ID):
				g.luckyMultCascade.tryLuckyMultCascadeFish(g, p)
			case isLuckyAwakenBossV2Fish(t.Def.ID):
				g.luckyAwakenBossV2.tryLuckyAwakenBossV2Fish(g, p)
			case isLuckyUltimateJudgmentFish(t.Def.ID):
				g.luckyUltimateJudgment.tryLuckyUltimateJudgmentFish(g, p)
			// DAY-310 新增
			case isLuckyComboBurstFish(t.Def.ID):
				g.luckyComboBurst.tryLuckyComboBurstFish(g, p)
			case isLuckyTimeBombFish(t.Def.ID):
				g.luckyTimeBomb.tryLuckyTimeBombFish(g, p)
			case isLuckyElementalFusionFish(t.Def.ID):
				g.luckyElementalFusion.tryLuckyElementalFusionFish(g, p)
			case isLuckyTreasureHunterFish(t.Def.ID):
				g.luckyTreasureHunter.tryLuckyTreasureHunterFish(g, p)
			case isLuckyMythAwakenFish(t.Def.ID):
				g.luckyMythAwaken.tryLuckyMythAwakenFish(g, p)
			// DAY-312 新增
			case isLuckyStarPortalFish(t.Def.ID):
				g.luckyStarPortal.tryLuckyStarPortalFish(g, p)
			case isLuckyDragonSoulFish(t.Def.ID):
				g.luckyDragonSoul.tryLuckyDragonSoulFish(g, p)
			case isLuckySpacetimeRiftFish(t.Def.ID):
				g.luckySpacetimeRift.tryLuckySpacetimeRiftFish(g, p)
			case isLuckyHolyJudgmentFish(t.Def.ID):
				g.luckyHolyJudgment.tryLuckyHolyJudgmentFish(g, p)
			case isLuckyBigBangFish(t.Def.ID):
				g.luckyBigBang.tryLuckyBigBangFish(g, p)
			// DAY-313 新增 Progressive Jackpot 系列
			case isLuckyJackpotPoolFish(t.Def.ID):
				g.luckyJackpotPool.onTargetKill(g, p, t.Def.ID)
			// DAY-314 新增
			case isLuckyMultiverseFish(t.Def.ID):
				g.luckyMultiverse.tryLuckyMultiverseFish(g, p)
			case isLuckyTimeLoopFish(t.Def.ID):
				g.luckyTimeLoop.tryLuckyTimeLoopFish(g, p)
			case isLuckyFateWheelFish(t.Def.ID):
				g.luckyFateWheel.tryLuckyFateWheelFish(g, p)
			case isLuckyDivineRealmFish(t.Def.ID):
				g.luckyDivineRealm.tryLuckyDivineRealmFish(g, p)
			case isLuckyFinalPowerFish(t.Def.ID):
				g.luckyFinalPower.tryLuckyFinalPowerFish(g, p)
			// DAY-315 新增
			case isLuckyMutationFish(t.Def.ID):
				g.luckyMutation.tryLuckyMutationFish(g, p)
			case isLuckyArcticStormFish(t.Def.ID):
				g.luckyArcticStorm.tryLuckyArcticStormFish(g, p)
			case isLuckyFisherWildFish(t.Def.ID):
				g.luckyFisherWild.tryLuckyFisherWildFish(g, p)
			case isLuckyRiskLevelFish(t.Def.ID):
				g.luckyRiskLevel.tryLuckyRiskLevelFish(g, p)
			case isLuckyCosmicPulseFish(t.Def.ID):
				g.luckyCosmicPulse.tryLuckyCosmicPulseFish(g, p)
			// DAY-316 新增
			case isLuckyMirrorUniverseFish(t.Def.ID):
				g.luckyMirrorUniverse.tryLuckyMirrorUniverseFish(g, p)
			case isLuckyGravityFieldFish(t.Def.ID):
				g.luckyGravityField.tryLuckyGravityFieldFish(g, p)
			case isLuckyTimeAccelerationFish(t.Def.ID):
				g.luckyTimeAcceleration.tryLuckyTimeAccelerationFish(g, p)
			case isLuckyNebulaVortexFish(t.Def.ID):
				g.luckyNebulaVortex.tryLuckyNebulaVortexFish(g, p)
			case isLuckyCosmicJudgmentFish(t.Def.ID):
				g.luckyCosmicJudgment.tryLuckyCosmicJudgmentFish(g, p)
			// DAY-317 新增
			case isLuckyPvpBattleFish(t.Def.ID):
				g.luckyPvpBattle.tryLuckyPvpBattleFish(g, p)
			case isLuckySkillChainFish(t.Def.ID):
				g.luckySkillChain.tryLuckySkillChainFish(g, p)
			case isLuckyGlobalExplosionFish(t.Def.ID):
				g.luckyGlobalExplosion.tryLuckyGlobalExplosionFish(g, p)
			case isLuckySpacetimeFoldFish(t.Def.ID):
				g.luckySpacetimeFold.tryLuckySpacetimeFoldFish(g, p)
			case isLuckyCosmicEndFish(t.Def.ID):
				g.luckyCosmicEnd.tryLuckyCosmicEndFish(g, p)
			}
			if g.luckyChainExplosion.isChainExplosionActive(playerID) {
				g.notifyChainExplosionKill(playerID, killerName, t.X, t.Y)
			}
			// 凍結期間擊破計數
			if g.luckyTimeFreeze.isTimeFreezeActive() {
				g.luckyTimeFreeze.notifyFreezeKill(playerID)
			}
			// 吸血鬼模式：每次擊破吸收倍率
			if g.luckyVampireMult.isVampireActive(playerID) {
				g.luckyVampireMult.notifyVampireKill(g, playerID)
			}
			if g.luckyAwakenedPhoenix.isAwakenedPhoenixActive(playerID) {
				powerUpMult, powerReward, isDone := g.luckyAwakenedPhoenix.consumeAwakenedPhoenixShot(playerID, true, bet.BetCost)
				if powerReward > 0 {
					p.AddCoins(powerReward)
					reward += powerReward
					result.Reward = reward
					g.notifyAwakenedPhoenixShot(playerID, killerName, powerUpMult, powerReward, isDone)
				}
			}
			// 全服合作魚：擊破計數
			if g.luckyCoopFish.isCoopActive() {
				g.notifyCoopKill(playerID, killerName, false)
			}
			// 時間扭曲魚：擊破計數
			if g.luckyTimeWarp.isTimeWarpActive() {
				g.luckyTimeWarp.notifyWarpKill(playerID)
			}
			// DAY-304 賞金獵人：賞金目標擊破通知
			if g.luckyBountyHunter.isBountyTarget(t.InstanceID) {
				g.notifyBountyKill(t.InstanceID, playerID, killerName)
			}
			// DAY-315 漁夫野生：Wild 目標擊破通知
			if g.luckyFisherWild.isWildTarget(t.InstanceID) {
				g.luckyFisherWild.onWildKilled(g, t.InstanceID)
			}
			// DAY-305 龍怒蓄積：射擊計數
			if g.luckyDragonWrathV2.isDragonWrathV2Active(playerID) {
				g.luckyDragonWrathV2.addWrathV2(playerID)
			}
			// DAY-305 公會戰：擊破計數
			if g.luckyGuildWar.isGuildWarActive() {
				g.notifyGuildWarKill(playerID, killerName)
			}
			// DAY-307 無限模式：擊破計數
			if g.luckyInfinite.isInfiniteActive(playerID) {
				g.luckyInfinite.notifyInfiniteKill(g, playerID)
			}
			// DAY-307 重生魚：重生目標擊破通知
			if g.luckyRebirth.isRebirthActive() {
				g.luckyRebirth.notifyRebirthKill(g, t.InstanceID)
			}
			// DAY-308 吸血鬼升級：擊破計數（有 session 就通知）
			g.luckyVampireV2.notifyVampireV2Kill(g, p)
			// DAY-308 不死 BOSS：擊破通知
			g.luckyImmortalBoss.notifyImmortalBossKill(g, p)
			// DAY-309 冰鳳凰：凍結期間擊破計數
			if g.luckyIcePhoenix.isIcePhoenixFrozen() {
				g.luckyIcePhoenix.notifyIcePhoenixKill(g, p)
			}
			// DAY-309 龍怒能量：擊破計數
			g.luckyDragonFury.notifyDragonFuryKill(g, p)
			// DAY-309 倍率瀑布：擊破計數
			g.luckyMultCascade.notifyMultCascadeKill(g, p)
			// DAY-309 覺醒 BOSS v2：Power Up 計數
			g.luckyAwakenBossV2.notifyAwakenBossV2Kill(g, p)
			// DAY-310 連擊爆發：擊破計數
			g.luckyComboBurst.onKillDuringComboBurst(g, p)
			// DAY-310 時間炸彈：能量累積
			g.luckyTimeBomb.onKillDuringTimeBomb(g, p)
			// DAY-310 神話覺醒：擊破計數
			g.luckyMythAwaken.onKillDuringMyth(playerID)
			// DAY-312 龍魂融合：魂計數
			g.luckyDragonSoul.onKillDuringDragonSoul(playerID)
			// DAY-312 時空裂縫：擊破計數
			g.luckySpacetimeRift.onKillDuringRift(playerID)
			// DAY-314 多重宇宙：擊破計數
			g.luckyMultiverse.onKill(g, p)
			// DAY-316 時間加速：擊破計數
			if g.luckyTimeAcceleration.isTimeAccelerationActive() {
				g.luckyTimeAcceleration.onKill(playerID)
			}
			// DAY-317 技能連鎖：擊破提升等級
			if g.luckySkillChain.isSkillChainActive(playerID) {
				g.luckySkillChain.onKillDuringSkillChain(g, p)
			}
		}

		// 獎勵通知（單播）
		g.hub.Send(playerID, protocol.MsgReward, protocol.RewardPayload{
			Source:     "target",
			Amount:     reward,
			Multiplier: t.Multiplier,
			NewBalance: p.Coins,
		})

		// 勞動值滿 → 觸發 Bonus
		if laborFull && g.state == StateNormalPlay {
			g.triggerBonus(playerID)
		}
	} else {
		// 未擊破，更新 HP（視覺反饋）
		t.HP = max(1, t.HP-bet.AttackPower/5)
		// BOSS Phase 2：HP 降到 50% 以下時觸發（未擊破時也要檢查）
		if t.Def.Type == data.TypeBoss && !g.bossPhase2 && t.HPPercent() < 0.5 {
			g.bossPhase2 = true
			g.hub.Broadcast(protocol.MsgBossEvent, protocol.BossEventPayload{
				Event:      "phase_change",
				InstanceID: t.InstanceID,
				HP:         t.HP,
				MaxHP:      t.MaxHP,
			})
			log.Printf("[Game] Boss Phase 2 triggered (no kill)! HP: %d/%d", t.HP, t.MaxHP)
		}
		// BOSS Phase 3：HP 降到 20% 以下時觸發（未擊破時也要檢查，絕望模式）
		if t.Def.Type == data.TypeBoss && !g.bossPhase3 && t.HPPercent() < 0.2 {
			g.bossPhase3 = true
			g.hub.Broadcast(protocol.MsgBossEvent, protocol.BossEventPayload{
				Event:      "phase_change",
				Phase:      3,
				InstanceID: t.InstanceID,
				HP:         t.HP,
				MaxHP:      t.MaxHP,
			})
			log.Printf("[Game] Boss Phase 3 triggered (no kill, 絕望模式)! HP: %d/%d", t.HP, t.MaxHP)
		}
		g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
			InstanceID: t.InstanceID,
			HP:         t.HP,
			MaxHP:      t.MaxHP,
			X:          t.X,
			Y:          t.Y,
			IsFleeing:  t.IsFleeing,
		})
	}

	g.hub.Send(playerID, protocol.MsgAttackResult, result)
	g.sendPlayerUpdate(playerID)
}

// ── BOSS 系統 ─────────────────────────────────────────────────

func (g *Game) startBossWarning() {
	g.setState(StateBossWarning)
	g.bossWarningAt = time.Now()
	g.hub.Broadcast(protocol.MsgBossEvent, protocol.BossEventPayload{Event: "warning"})
	log.Printf("[Game] Boss warning!")
}

func (g *Game) spawnBoss() {
	def, ok := data.GetTarget("B001")
	if !ok {
		return
	}
	g.boss = NewTarget(def, GameWidth/2, GameHeight/2)
	g.bossTimer = 60.0
	g.bossPhase2 = false // 重置 Phase 2 狀態
	g.bossPhase3 = false // 重置 Phase 3 狀態
	g.setState(StateBossBattle)

	// 清除部分普通目標（BOSS 期間最多 8 個）
	count := 0
	for id := range g.targets {
		if count >= MaxBossTargets {
			delete(g.targets, id)
		}
		count++
	}

	g.hub.Broadcast(protocol.MsgBossEvent, protocol.BossEventPayload{
		Event:      "spawn",
		InstanceID: g.boss.InstanceID,
		HP:         g.boss.HP,
		MaxHP:      g.boss.MaxHP,
	})
	g.hub.Broadcast(protocol.MsgTargetSpawn, g.targetSpawnPayload(g.boss))
	log.Printf("[Game] Boss spawned!")
}

func (g *Game) handleBossKill(playerID string, p *Player, boss *Target, reward int) {
	// 依剩餘時間計算倍率
	mult := bossRewardMult(g.bossTimer)
	finalReward := int(float64(p.GetBetDef().BetCost) * mult)
	p.AddCoins(finalReward - reward) // 補差額

	g.boss = nil
	g.setState(StateNormalPlay)
	g.nextBossIn = 180 + rand.Float64()*120

	g.hub.Broadcast(protocol.MsgBossEvent, protocol.BossEventPayload{
		Event:      "kill",
		InstanceID: boss.InstanceID,
		Reward:     finalReward,
		Multiplier: mult,
	})
	log.Printf("[Game] Boss killed by %s! Reward: %d (%.0fx)", playerID, finalReward, mult)
}

func bossRewardMult(timeLeft float64) float64 {
	switch {
	case timeLeft > 50:
		return 500
	case timeLeft > 40:
		return 400
	case timeLeft > 30:
		return 300
	case timeLeft > 20:
		return 200
	case timeLeft > 10:
		return 150
	default:
		return 100
	}
}

func (g *Game) triggerBoss() {
	if g.state != StateNormalPlay {
		return
	}
	g.startBossWarning()
}

// ── Bonus 系統 ────────────────────────────────────────────────

func (g *Game) triggerBonus(playerID string) {
	if g.state != StateNormalPlay {
		return
	}
	g.mu.Lock()
	p, ok := g.players[playerID]
	if !ok {
		g.mu.Unlock()
		return
	}
	p.EntryBetCost = p.GetBetDef().BetCost
	p.ResetLabor()
	g.bonusPlayer = playerID
	g.bonusTimer = 15.0
	g.bonusScore = make(map[string]int)
	g.bonusWeedHP = make(map[string]int)
	g.setState(StateBonusGame)
	g.mu.Unlock()

	g.hub.Broadcast(protocol.MsgBonusEvent, protocol.BonusEventPayload{
		Event:    "start",
		TimeLeft: 15.0,
	})
	g.sendPlayerUpdate(playerID)
	log.Printf("[Game] Bonus triggered by %s", playerID)
}

func (g *Game) handleBonusClick(playerID string, req protocol.BonusClickRequest) {
	if g.state != StateBonusGame {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()

	p, ok := g.players[playerID]
	if !ok {
		return
	}

	// 找 Bonus 目標物定義
	var bonusDef *data.BonusTargetDef
	for i := range data.BonusTargets {
		if data.BonusTargets[i].ID == req.TargetID {
			bonusDef = &data.BonusTargets[i]
			break
		}
	}
	if bonusDef == nil {
		return
	}

	// BG002 硬雜草需要連點 2 次
	if bonusDef.Special == "hard" {
		hp := g.bonusWeedHP[req.TargetID]
		if hp == 0 {
			hp = 2
		}
		hp--
		g.bonusWeedHP[req.TargetID] = hp
		if hp > 0 {
			return // 還沒打完
		}
	}

	score := bonusDef.Score
	g.bonusScore[playerID] += score

	g.hub.Broadcast(protocol.MsgBonusEvent, protocol.BonusEventPayload{
		Event:    "click",
		TimeLeft: g.bonusTimer,
		Score:    g.bonusScore[playerID],
	})

	_ = p // 避免 unused
}

func (g *Game) endBonus() {
	p, ok := g.players[g.bonusPlayer]
	if !ok {
		g.setState(StateNormalPlay)
		return
	}

	score := g.bonusScore[g.bonusPlayer]
	mult := 20.0 + float64(score)*0.375
	if mult < 20 {
		mult = 20
	}
	if mult > 50 {
		mult = 50
	}
	reward := int(float64(p.EntryBetCost) * mult)
	p.AddCoins(reward)

	g.setState(StateBonusResult)
	g.hub.Broadcast(protocol.MsgBonusEvent, protocol.BonusEventPayload{
		Event:      "result",
		Score:      score,
		Multiplier: mult,
		Reward:     reward,
	})
	g.sendPlayerUpdate(g.bonusPlayer)
	log.Printf("[Game] Bonus ended. Score: %d, Mult: %.1f, Reward: %d", score, mult, reward)

	// 2 秒後回到正常遊戲
	time.AfterFunc(2*time.Second, func() {
		g.mu.Lock()
		g.setState(StateNormalPlay)
		g.mu.Unlock()
	})
}

// ── AUTO 射擊 ─────────────────────────────────────────────────

func (g *Game) autoFire(p *Player) {
	if g.state != StateNormalPlay && g.state != StateBossBattle {
		return
	}
	// 找最高價值目標
	var best *Target
	bestScore := -1.0
	for _, t := range g.targets {
		score := t.Multiplier * 2.0
		if t.X < 400 {
			score += 20
		}
		if score > bestScore {
			bestScore = score
			best = t
		}
	}
	if g.boss != nil {
		best = g.boss
	}
	if best == nil {
		return
	}
	// 模擬攻擊
	g.handleAttack(p.ID, protocol.AttackRequest{
		TargetID: best.InstanceID,
		ClickX:   best.X,
		ClickY:   best.Y,
	})
}

// ── 輔助函數 ──────────────────────────────────────────────────

func (g *Game) setState(s GameState) {
	g.state = s
	g.hub.Broadcast(protocol.MsgGameState, protocol.GameStatePayload{
		State:     string(s),
		Timestamp: time.Now().UnixMilli(),
	})
}

func (g *Game) handleAutoToggle(playerID string) {
	g.mu.Lock()
	p, ok := g.players[playerID]
	if ok {
		p.IsAuto = !p.IsAuto
	}
	g.mu.Unlock()
	if ok {
		g.sendPlayerUpdate(playerID)
	}
}

func (g *Game) handleBetChange(playerID string, level int) {
	g.mu.Lock()
	p, ok := g.players[playerID]
	if ok {
		if level < 1 {
			level = 1
		}
		if level > 10 {
			level = 10
		}
		p.BetLevel = level
	}
	g.mu.Unlock()
	if ok {
		g.sendPlayerUpdate(playerID)
	}
}

func (g *Game) handleLock(playerID string, targetID string) {
	g.mu.Lock()
	p, ok := g.players[playerID]
	if ok {
		p.LockTargetID = targetID
	}
	g.mu.Unlock()
	if ok {
		g.sendPlayerUpdate(playerID)
	}
}

func (g *Game) sendPlayerUpdate(playerID string) {
	g.mu.RLock()
	p, ok := g.players[playerID]
	g.mu.RUnlock()
	if !ok {
		return
	}
	bet := p.GetBetDef()
	g.hub.Send(playerID, protocol.MsgPlayerUpdate, protocol.PlayerUpdatePayload{
		ID:              p.ID,
		Coins:           p.Coins,
		BetLevel:        p.BetLevel,
		BetCost:         bet.BetCost,
		CharacterID:     p.GetCharacterID(),
		CharacterName:   p.GetCharacterName(),
		LaborValue:      p.LaborValue,
		IsAuto:          p.IsAuto,
		LockTargetID:    p.LockTargetID,
		ProjectileSpeed: bet.ProjectileSpeed,
		FireRate:        bet.FireRate,
		ComboCount:      p.ComboCount,
		ComboMultBonus:  p.GetComboMultBonus(),
	})
}

func (g *Game) GetState() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return string(g.state)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// handleSetDisplayName 設定玩家顯示名稱
func (g *Game) handleSetDisplayName(playerID string, name string) {
	if name == "" || len(name) > 20 {
		return
	}
	g.mu.Lock()
	p, ok := g.players[playerID]
	if ok {
		p.DisplayName = name
	}
	g.mu.Unlock()
	if ok {
		g.sendPlayerUpdate(playerID)
		log.Printf("[Game] Player %s set display name: %s", playerID, name)
	}
}

// ── DAY-306 輔助方法 ──────────────────────────────────────────

// broadcast 廣播 Envelope（供 Lucky handler 使用）
func (g *Game) broadcast(env protocol.Envelope) {
	g.hub.Broadcast(env.Type, env.Payload)
}

// sendAnnounce 廣播公告（供 Lucky handler 使用）
func (g *Game) sendAnnounce(msg, priority, color string) {
	g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  msg,
		Priority: priority,
		Color:    color,
	})
}

// applyAOEDamage 對範圍內所有目標造成百分比傷害，回傳命中數
// cx, cy: 中心座標（若 radius >= 99999 則全場）
// radius: 影響半徑（像素）
// pct: 傷害百分比（0.0-1.0）
func (g *Game) applyAOEDamage(cx, cy, radius, pct float64) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	hitCount := 0
	for _, t := range g.targets {
		if t.Def.Type == "boss" {
			continue // BOSS 不受 AOE 影響
		}
		// 距離判定
		if radius < 99999 {
			dx := t.X - cx
			dy := t.Y - cy
			if dx*dx+dy*dy > radius*radius {
				continue
			}
		}
		// 造成傷害
		damage := int(float64(t.MaxHP) * pct)
		if damage < 1 {
			damage = 1
		}
		t.HP -= damage
		if t.HP < 1 {
			t.HP = 1
		}
		hitCount++

		// 廣播 HP 更新
		g.hub.Broadcast(protocol.MsgTargetUpdate, protocol.TargetUpdatePayload{
			InstanceID: t.InstanceID,
			HP:         t.HP,
			MaxHP:      t.MaxHP,
			X:          t.X,
			Y:          t.Y,
		})
	}
	return hitCount
}

// applyUltimateJudgment — T160 終極審判：全場目標 HP 歸零，每個獎勵 rewardMult 倍
// 返回命中目標數
func (g *Game) applyUltimateJudgment(p *Player, rewardMult float64) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	hitCount := 0
	bet := p.GetBetDef()
	var toDelete []string
	for id, t := range g.targets {
		t.HP = 0
		reward := int(float64(bet.BetCost) * rewardMult)
		p.AddCoins(reward)
		toDelete = append(toDelete, id)
		g.hub.Broadcast(protocol.MsgTargetKill, protocol.TargetKillPayload{
			InstanceID: t.InstanceID,
			DefID:      t.Def.ID,
			Multiplier: t.Multiplier * rewardMult,
			Reward:     reward,
			LaborGain:  0,
			KillerID:   p.ID,
		})
		hitCount++
	}
	for _, id := range toDelete {
		delete(g.targets, id)
	}
	return hitCount
}
