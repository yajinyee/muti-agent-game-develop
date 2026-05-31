// Package game — 遊戲主循環與狀態機
// server-core-agent 負責維護
package game

import (
	"encoding/json"
	"fmt"
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

	// DAY-318 新增
	luckyDragonKing     *luckyDragonKingManager
	luckyEternalCycle   *luckyEternalCycleManager
	luckyChaosExplosion *luckyChaosExplosionManager
	luckyDivineRevival  *luckyDivineRevivalManager
	luckyGenesisEpoch   *luckyGenesisEpochManager

	// DAY-319 新增
	luckyEnergyStorm       *luckyEnergyStormManager
	luckyCrystalResonance  *luckyCrystalResonanceManager
	luckyFateJudgment      *luckyFateJudgmentManager
	luckyTimeReversal      *luckyTimeReversalManager
	luckyCosmicSingularity *luckyCosmicSingularityManager

	// DAY-323 新增
	luckyFeverBoost      *luckyFeverBoostManager
	luckyGuildBattle     *luckyGuildBattleManager
	luckyPathFish        *luckyPathFishManager
	luckyChainEel        *luckyChainEelManager
	luckyUltimateMiracle *luckyUltimateMiracleManager

	// DAY-324 新增
	luckyAvalanche        *luckyAvalancheManager
	luckyCrashMultiplier  *luckyCrashMultiplierManager
	luckyMultiplierLadder *luckyMultiplierLadderManager
	luckyIceFishingWheel  *luckyIceFishingWheelManager
	luckyGlobalAvalanche  *luckyGlobalAvalancheManager

	// DAY-325 新增
	luckyFishingNet      *luckyFishingNetManager
	luckyTNTBonus        *luckyTNTBonusManager
	luckyDisturbance     *luckyDisturbanceManager
	luckyPearlMultiplier *luckyPearlMultiplierManager
	luckyRapidRiches     *luckyRapidRichesManager

	// DAY-326 新增
	luckyDiceBonus  *luckyDiceBonusManager
	luckyDualBonus  *luckyDualBonusManager
	luckyCoinRespin *luckyCoinRespinManager

	// DAY-327 新增
	luckyGoldenPot    *luckyGoldenPotManager
	luckyCascadeLock  *luckyCascadeLockManager
	luckyLegendAwaken *luckyLegendAwakenManager
	luckyCrashHarvest *luckyCrashHarvestManager
	luckyCosmicFusion *luckyCosmicFusionManager

	// DAY-328 新增
	luckyMagneticAttraction *luckyMagneticAttractionManager
	luckySuperChain         *luckySuperChainManager
	luckyHolyPillar         *luckyHolyPillarManager
	luckyTimeStop           *luckyTimeStopManager
	luckyCosmicRestart      *luckyCosmicRestartManager

	// DAY-329 新增
	luckyFeverBoostUltimate  *luckyFeverBoostUltimateManager
	luckyRapidRichesUltimate *luckyRapidRichesUltimateManager
	luckyIceFishingMaster    *luckyIceFishingMasterManager
	luckyCosmicMiracle       *luckyCosmicMiracleManager
	luckyGenesisUltimate     *luckyGenesisUltimateManager

	// DAY-331 新增
	luckySharkSpark       *luckySharkSparkManager
	luckyWinterIce        *luckyWinterIceManager
	luckyAtlantisFrenzy   *luckyAtlantisFrenzyManager
	luckyFishingTimeWheel *luckyFishingTimeWheelManager
	luckyUltimateShark    *luckyUltimateSharkManager

	// DAY-332 新增
	luckyWildCollector     *luckyWildCollectorManager
	luckyLightningEelUltra *luckyLightningEelUltraManager
	luckyDominoChain       *luckyDominoChainManager
	luckyImmortalBossUltra *luckyImmortalBossUltraManager
	luckyQuadFusion        *luckyQuadFusionManager

	// DAY-333 新增
	luckyElectricalFrame *luckyElectricalFrameManager
	luckyMagneticRespin  *luckyMagneticRespinManager
	luckyFishermanTrail  *luckyFishermanTrailManager
	luckyGoldenGills     *luckyGoldenGillsManager
	luckyPentaFusion     *luckyPentaFusionManager

	lastTick time.Time

	// DAY-345 每日任務系統
	dailyQuest *DailyQuestSystem

	// DAY-346 每週挑戰系統
	weeklyChallenge *WeeklyChallengeSystem

	// DAY-347 賽季通行證系統
	seasonPass *SeasonPassManager

	// DAY-348 任務幣兌換商店 + 賽季排行榜
	questShop         *QuestShop
	seasonLeaderboard *SeasonLeaderboard

	// DAY-349 成就系統 + 好友排行榜
	achievementSystem *AchievementSystem
	friendSystem      *FriendSystem
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

		// DAY-318 新增
		luckyDragonKing:     newLuckyDragonKingManager(),
		luckyEternalCycle:   newLuckyEternalCycleManager(),
		luckyChaosExplosion: newLuckyChaosExplosionManager(),
		luckyDivineRevival:  newLuckyDivineRevivalManager(),
		luckyGenesisEpoch:   newLuckyGenesisEpochManager(),

		// DAY-319 新增
		luckyEnergyStorm:       newLuckyEnergyStormManager(),
		luckyCrystalResonance:  newLuckyCrystalResonanceManager(),
		luckyFateJudgment:      newLuckyFateJudgmentManager(),
		luckyTimeReversal:      newLuckyTimeReversalManager(),
		luckyCosmicSingularity: newLuckyCosmicSingularityManager(),

		// DAY-323 新增
		luckyFeverBoost:      newLuckyFeverBoostManager(),
		luckyGuildBattle:     newLuckyGuildBattleManager(),
		luckyPathFish:        newLuckyPathFishManager(),
		luckyChainEel:        newLuckyChainEelManager(),
		luckyUltimateMiracle: newLuckyUltimateMiracleManager(),

		// DAY-324 新增
		luckyAvalanche:        newLuckyAvalancheManager(),
		luckyCrashMultiplier:  newLuckyCrashMultiplierManager(),
		luckyMultiplierLadder: newLuckyMultiplierLadderManager(),
		luckyIceFishingWheel:  newLuckyIceFishingWheelManager(),
		luckyGlobalAvalanche:  newLuckyGlobalAvalancheManager(),

		// DAY-325 新增
		luckyFishingNet:      newLuckyFishingNetManager(),
		luckyTNTBonus:        newLuckyTNTBonusManager(),
		luckyDisturbance:     newLuckyDisturbanceManager(),
		luckyPearlMultiplier: newLuckyPearlMultiplierManager(),
		luckyRapidRiches:     newLuckyRapidRichesManager(),

		// DAY-326 新增
		luckyDiceBonus:  newLuckyDiceBonusManager(),
		luckyDualBonus:  newLuckyDualBonusManager(),
		luckyCoinRespin: newLuckyCoinRespinManager(),

		// DAY-327 新增
		luckyGoldenPot:    newLuckyGoldenPotManager(),
		luckyCascadeLock:  newLuckyCascadeLockManager(),
		luckyLegendAwaken: newLuckyLegendAwakenManager(),
		luckyCrashHarvest: newLuckyCrashHarvestManager(),
		luckyCosmicFusion: newLuckyCosmicFusionManager(),

		// DAY-328 新增
		luckyMagneticAttraction: newLuckyMagneticAttractionManager(),
		luckySuperChain:         newLuckySuperChainManager(),
		luckyHolyPillar:         newLuckyHolyPillarManager(),
		luckyTimeStop:           newLuckyTimeStopManager(),
		luckyCosmicRestart:      newLuckyCosmicRestartManager(),

		// DAY-329 新增
		luckyFeverBoostUltimate:  newLuckyFeverBoostUltimateManager(),
		luckyRapidRichesUltimate: newLuckyRapidRichesUltimateManager(),
		luckyIceFishingMaster:    newLuckyIceFishingMasterManager(),
		luckyCosmicMiracle:       newLuckyCosmicMiracleManager(),
		luckyGenesisUltimate:     newLuckyGenesisUltimateManager(),

		// DAY-331 新增
		luckySharkSpark:       newLuckySharkSparkManager(),
		luckyWinterIce:        newLuckyWinterIceManager(),
		luckyAtlantisFrenzy:   newLuckyAtlantisFrenzyManager(),
		luckyFishingTimeWheel: newLuckyFishingTimeWheelManager(),
		luckyUltimateShark:    newLuckyUltimateSharkManager(),

		// DAY-332 新增
		luckyWildCollector:     newLuckyWildCollectorManager(),
		luckyLightningEelUltra: newLuckyLightningEelUltraManager(),
		luckyDominoChain:       newLuckyDominoChainManager(),
		luckyImmortalBossUltra: newLuckyImmortalBossUltraManager(),
		luckyQuadFusion:        newLuckyQuadFusionManager(),

		// DAY-333 新增
		luckyElectricalFrame: newLuckyElectricalFrameManager(),
		luckyMagneticRespin:  newLuckyMagneticRespinManager(),
		luckyFishermanTrail:  newLuckyFishermanTrailManager(),
		luckyGoldenGills:     newLuckyGoldenGillsManager(),
		luckyPentaFusion:     newLuckyPentaFusionManager(),

		// DAY-345 每日任務系統
		dailyQuest: newDailyQuestSystem(),

		// DAY-346 每週挑戰系統
		weeklyChallenge: newWeeklyChallengeSystem(),

		// DAY-347 賽季通行證系統
		seasonPass: NewSeasonPassManager(),

		// DAY-348 任務幣兌換商店 + 賽季排行榜
		questShop:         NewQuestShop(),
		seasonLeaderboard: NewSeasonLeaderboard(time.Now().In(time.FixedZone("UTC+8", 8*60*60)).Format("2006-01")),

		// DAY-349 成就系統 + 好友排行榜
		achievementSystem: newAchievementSystem(),
		friendSystem:      newFriendSystem(),
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
	// DAY-349 好友系統：標記玩家離線
	g.friendSystem.SetOffline(id)
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
	// DAY-345 每日任務系統
	case protocol.MsgDailyQuestRequest:
		g.handleDailyQuestRequest(clientID)
	case protocol.MsgDailyQuestClaim:
		var req protocol.DailyQuestClaimRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			return
		}
		g.handleDailyQuestClaim(clientID, req.QuestID)
	// DAY-346 每週挑戰系統
	case protocol.MsgWeeklyChallengeRequest:
		g.handleWeeklyChallengeRequest(clientID)
	case protocol.MsgWeeklyChallengeClaim:
		var req protocol.WeeklyChallengeClaim
		if err := json.Unmarshal(payload, &req); err != nil {
			return
		}
		g.handleWeeklyChallengeClaim(clientID, req.ChallengeID)
	// DAY-348 任務幣兌換商店 + 賽季排行榜
	case protocol.MsgShopRequest:
		g.handleShopRequest(clientID)
	case protocol.MsgShopPurchase:
		var req protocol.ShopPurchaseRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			return
		}
		g.handleShopPurchase(clientID, req.ItemID)
	case protocol.MsgSeasonLeaderboardRequest:
		g.handleSeasonLeaderboardRequest(clientID)
	// DAY-349 成就系統 + 好友排行榜
	case protocol.MsgAchievementListRequest:
		g.handleAchievementListRequest(clientID)
	case protocol.MsgRoomLeaderboardRequest:
		g.handleRoomLeaderboardRequest(clientID)
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
	g.handleAttackLocked(playerID, req)
}

// handleAttackLocked — 攻擊處理（不加鎖版本，供 tick 內部呼叫）
// DAY-338 修復死鎖：autoFire 在 tick 的 g.mu.Lock() 保護下執行，不能再次加鎖
func (g *Game) handleAttackLocked(playerID string, req protocol.AttackRequest) {
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
		// Disturbance 系統：記錄最近擊破（T218 擾動魚）
		p.AddRecentKill()
		effectiveMult := t.Multiplier * (1.0 + comboBonus)

		// DAY-345 每日任務：連擊達成計數（goroutine 避免死鎖）
		go g.notifyQuestProgress(playerID, "combo", p.ComboCount)

		// DAY-346 每週挑戰：連擊達成計數（goroutine 避免死鎖）
		comboVal := p.ComboCount
		go g.notifyWeeklyChallengeProgress(playerID, "combo", false, 0.0, comboVal)

		// DAY-347 賽季通行證：連擊里程碑 XP
		if comboVal == 5 {
			go g.addSeasonXP(playerID, XPPerCombo5, "combo")
		} else if comboVal == 10 {
			go g.addSeasonXP(playerID, XPPerCombo10, "combo")
		}

		// DAY-349 成就系統：連擊觸發
		go g.notifyAchievements(playerID, g.achievementSystem.OnCombo(playerID, comboVal))

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
		// DAY-318 新增全服加成
		dragonKingBoost := g.luckyDragonKing.getDragonKingMult()
		if dragonKingBoost > 1.0 {
			effectiveMult *= dragonKingBoost
		}
		eternalCycleBoost := g.luckyEternalCycle.getEternalCycleMult()
		if eternalCycleBoost > 1.0 {
			effectiveMult *= eternalCycleBoost
		}
		chaosExplosionBoost := g.luckyChaosExplosion.getChaosExplosionMult()
		if chaosExplosionBoost > 1.0 {
			effectiveMult *= chaosExplosionBoost
		}
		divineRevivalBoost := g.luckyDivineRevival.getDivineRevivalMult()
		if divineRevivalBoost > 1.0 {
			effectiveMult *= divineRevivalBoost
		}
		genesisEpochBoost := g.luckyGenesisEpoch.getGenesisEpochMult()
		if genesisEpochBoost > 1.0 {
			effectiveMult *= genesisEpochBoost
		}
		// DAY-319 新增全服加成
		energyStormBoost := g.luckyEnergyStorm.getEnergyStormMult()
		if energyStormBoost > 1.0 {
			effectiveMult *= energyStormBoost
		}
		crystalResonanceBoost := g.luckyCrystalResonance.getCrystalResonanceMult()
		if crystalResonanceBoost > 1.0 {
			effectiveMult *= crystalResonanceBoost
		}
		fateJudgmentBoost := g.luckyFateJudgment.getFateJudgmentMult()
		if fateJudgmentBoost > 1.0 {
			effectiveMult *= fateJudgmentBoost
		}
		// 命運審判：命運目標倍率加成
		fateTargetMult := g.luckyFateJudgment.getFateTargetMult(t.InstanceID)
		if fateTargetMult > 1.0 {
			effectiveMult *= fateTargetMult
			g.luckyFateJudgment.onFateTargetKilled(t.InstanceID)
		}
		timeReversalBoost := g.luckyTimeReversal.getTimeReversalMult()
		if timeReversalBoost > 1.0 {
			effectiveMult *= timeReversalBoost
		}
		// 時間逆流：逆流目標倍率加成
		reversalTargetMult := g.luckyTimeReversal.getReversalTargetMult(t.InstanceID)
		if reversalTargetMult > 1.0 {
			effectiveMult *= reversalTargetMult
			g.luckyTimeReversal.onReversalTargetKilled(t.InstanceID)
		}
		cosmicSingularityBoost := g.luckyCosmicSingularity.getCosmicSingularityMult()
		if cosmicSingularityBoost > 1.0 {
			effectiveMult *= cosmicSingularityBoost
		}
		// DAY-323 新增全服加成
		feverBoostMult := g.luckyFeverBoost.getFeverBoostMult()
		if feverBoostMult > 1.0 {
			effectiveMult *= feverBoostMult
		}
		guildBattleMult := g.luckyGuildBattle.getGuildBattleMult()
		if guildBattleMult > 1.0 {
			effectiveMult *= guildBattleMult
		}
		pathFishMult := g.luckyPathFish.getPathFishMult()
		if pathFishMult > 1.0 {
			effectiveMult *= pathFishMult
		}
		// 路徑魚：擊破時推進路徑
		pathKillMult := g.luckyPathFish.onKillDuringPath(g, p)
		if pathKillMult > 1.0 {
			effectiveMult *= pathKillMult
		}
		chainEelMult := g.luckyChainEel.getChainEelMult()
		if chainEelMult > 1.0 {
			effectiveMult *= chainEelMult
		}
		ultimateMiracleMult := g.luckyUltimateMiracle.getUltimateMiracleMult()
		if ultimateMiracleMult > 1.0 {
			effectiveMult *= ultimateMiracleMult
		}
		// DAY-328 新增全服加成
		magneticAttractionMult := g.luckyMagneticAttraction.getMagneticAttractionMult()
		if magneticAttractionMult > 1.0 {
			effectiveMult *= magneticAttractionMult
		}
		superChainMult := g.luckySuperChain.getSuperChainMult()
		if superChainMult > 1.0 {
			effectiveMult *= superChainMult
		}
		holyPillarMult := g.luckyHolyPillar.getHolyPillarMult()
		if holyPillarMult > 1.0 {
			effectiveMult *= holyPillarMult
		}
		timeStopMult := g.luckyTimeStop.getTimeStopMult()
		if timeStopMult > 1.0 {
			effectiveMult *= timeStopMult
		}
		cosmicRestartMult := g.luckyCosmicRestart.getCosmicRestartMult()
		if cosmicRestartMult > 1.0 {
			effectiveMult *= cosmicRestartMult
		}
		// DAY-329 新增倍率計算
		feverBoostUltimateMult := g.luckyFeverBoostUltimate.getFeverBoostUltimateMult()
		if feverBoostUltimateMult > 1.0 {
			effectiveMult *= feverBoostUltimateMult
		}
		rapidRichesUltimateMult := g.luckyRapidRichesUltimate.getRapidRichesUltimateMult()
		if rapidRichesUltimateMult > 1.0 {
			effectiveMult *= rapidRichesUltimateMult
		}
		iceFishingMasterMult := g.luckyIceFishingMaster.getIceFishingMasterMult()
		if iceFishingMasterMult > 1.0 {
			effectiveMult *= iceFishingMasterMult
		}
		cosmicMiracleMult := g.luckyCosmicMiracle.getCosmicMiracleMult()
		if cosmicMiracleMult > 1.0 {
			effectiveMult *= cosmicMiracleMult
		}
		genesisUltimateMult := g.luckyGenesisUltimate.getGenesisUltimateMult()
		if genesisUltimateMult > 1.0 {
			effectiveMult *= genesisUltimateMult
		}
		// DAY-331 新增倍率加成
		sharkSparkMult := g.luckySharkSpark.getSharkSparkMult()
		if sharkSparkMult > 1.0 {
			effectiveMult *= sharkSparkMult
		}
		winterIceMult := g.luckyWinterIce.getWinterIceMult()
		if winterIceMult > 1.0 {
			effectiveMult *= winterIceMult
		}
		atlantisFrenzyMult := g.luckyAtlantisFrenzy.getAtlantisFrenzyMult()
		if atlantisFrenzyMult > 1.0 {
			effectiveMult *= atlantisFrenzyMult
		}
		fishingTimeWheelMult := g.luckyFishingTimeWheel.getFishingTimeWheelMult()
		if fishingTimeWheelMult > 1.0 {
			effectiveMult *= fishingTimeWheelMult
		}
		ultimateSharkMult := g.luckyUltimateShark.getUltimateSharkMult()
		if ultimateSharkMult > 1.0 {
			effectiveMult *= ultimateSharkMult
		}
		// DAY-332 新增倍率加成
		wildCollectorMult := g.luckyWildCollector.getWildCollectorMult()
		if wildCollectorMult > 1.0 {
			effectiveMult *= wildCollectorMult
		}
		lightningEelUltraMult := g.luckyLightningEelUltra.getLightningEelUltraMult()
		if lightningEelUltraMult > 1.0 {
			effectiveMult *= lightningEelUltraMult
		}
		dominoChainMult := g.luckyDominoChain.getDominoChainMult()
		if dominoChainMult > 1.0 {
			effectiveMult *= dominoChainMult
		}
		immortalBossUltraMult := g.luckyImmortalBossUltra.getImmortalBossUltraMult()
		if immortalBossUltraMult > 1.0 {
			effectiveMult *= immortalBossUltraMult
		}
		quadFusionMult := g.luckyQuadFusion.getQuadFusionMult()
		if quadFusionMult > 1.0 {
			effectiveMult *= quadFusionMult
		}
		// DAY-333 新增倍率加成
		electricalFrameMult := g.luckyElectricalFrame.getElectricalFrameMult()
		if electricalFrameMult > 1.0 {
			effectiveMult *= electricalFrameMult
		}
		magneticRespinMult := g.luckyMagneticRespin.getMagneticRespinMult()
		if magneticRespinMult > 1.0 {
			effectiveMult *= magneticRespinMult
		}
		fishermanTrailMult := g.luckyFishermanTrail.getFishermanTrailMult()
		if fishermanTrailMult > 1.0 {
			effectiveMult *= fishermanTrailMult
		}
		goldenGillsMult := g.luckyGoldenGills.getGoldenGillsMult()
		if goldenGillsMult > 1.0 {
			effectiveMult *= goldenGillsMult
		}
		pentaFusionMult := g.luckyPentaFusion.getPentaFusionMult()
		if pentaFusionMult > 1.0 {
			effectiveMult *= pentaFusionMult
		}
		// 公會戰：擊破計數
		g.luckyGuildBattle.onKillDuringBattle(playerID)
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
			// DAY-318 新增
			case isLuckyDragonKingFish(t.Def.ID):
				g.luckyDragonKing.tryLuckyDragonKingFish(g, p)
			case isLuckyEternalCycleFish(t.Def.ID):
				g.luckyEternalCycle.tryLuckyEternalCycleFish(g, p)
			case isLuckyChaosExplosionFish(t.Def.ID):
				g.luckyChaosExplosion.tryLuckyChaosExplosionFish(g, p)
			case isLuckyDivineRevivalFish(t.Def.ID):
				g.luckyDivineRevival.tryLuckyDivineRevivalFish(g, p)
			case isLuckyGenesisEpochFish(t.Def.ID):
				g.luckyGenesisEpoch.tryLuckyGenesisEpochFish(g, p)
			// DAY-319 新增
			case isLuckyEnergyStormFish(t.Def.ID):
				g.luckyEnergyStorm.tryLuckyEnergyStormFish(g, p)
			case isLuckyCrystalResonanceFish(t.Def.ID):
				g.luckyCrystalResonance.tryLuckyCrystalResonanceFish(g, p)
			case isLuckyFateJudgmentFish(t.Def.ID):
				g.luckyFateJudgment.tryLuckyFateJudgmentFish(g, p)
			case isLuckyTimeReversalFish(t.Def.ID):
				g.luckyTimeReversal.tryLuckyTimeReversalFish(g, p)
			case isLuckyCosmicSingularityFish(t.Def.ID):
				g.luckyCosmicSingularity.tryLuckyCosmicSingularityFish(g, p)
			// DAY-323 新增
			case isLuckyFeverBoostFish(t.Def.ID):
				g.luckyFeverBoost.tryLuckyFeverBoostFish(g, p)
			case isLuckyGuildBattleFish(t.Def.ID):
				g.luckyGuildBattle.tryLuckyGuildBattleFish(g, p)
			case isLuckyPathFish(t.Def.ID):
				g.luckyPathFish.tryLuckyPathFish(g, p)
			case isLuckyChainEelFish(t.Def.ID):
				g.luckyChainEel.tryLuckyChainEelFish(g, p)
			case isLuckyUltimateMiracleFish(t.Def.ID):
				g.luckyUltimateMiracle.tryLuckyUltimateMiracleFish(g, p)
			// DAY-324 新增
			case isLuckyAvalancheFish(t.Def.ID):
				g.luckyAvalanche.tryLuckyAvalancheFish(g, p)
			case isLuckyCrashMultiplierFish(t.Def.ID):
				g.luckyCrashMultiplier.tryLuckyCrashMultiplierFish(g, p)
			case isLuckyMultiplierLadderFish(t.Def.ID):
				g.luckyMultiplierLadder.tryLuckyMultiplierLadderFish(g, p)
			case isLuckyIceFishingWheelFish(t.Def.ID):
				g.luckyIceFishingWheel.tryLuckyIceFishingWheelFish(g, p)
			case isLuckyGlobalAvalancheFish(t.Def.ID):
				g.luckyGlobalAvalanche.tryLuckyGlobalAvalancheFish(g, p)
			// DAY-325 新增
			case isLuckyFishingNetFish(t.Def.ID):
				g.luckyFishingNet.tryLuckyFishingNetFish(g, p)
			case isLuckyTNTBonusFish(t.Def.ID):
				g.luckyTNTBonus.tryLuckyTNTBonusFish(g, p)
			case isLuckyDisturbanceFish(t.Def.ID):
				g.luckyDisturbance.tryLuckyDisturbanceFish(g, p)
			case isLuckyPearlMultiplierFish(t.Def.ID):
				g.luckyPearlMultiplier.tryLuckyPearlMultiplierFish(g, p)
			case isLuckyRapidRichesFish(t.Def.ID):
				g.luckyRapidRiches.tryLuckyRapidRichesFish(g, p)
			// DAY-326 新增
			case isLuckyDiceBonusFish(t.Def.ID):
				g.luckyDiceBonus.tryLuckyDiceBonusFish(g, p)
			case isLuckyDualBonusFish(t.Def.ID):
				g.luckyDualBonus.tryLuckyDualBonusFish(g, p)
			case isLuckyCoinRespinFish(t.Def.ID):
				g.luckyCoinRespin.tryLuckyCoinRespinFish(g, p)
			// DAY-327 新增
			case isLuckyGoldenPotFish(t.Def.ID):
				g.luckyGoldenPot.tryLuckyGoldenPotFish(g, p)
			case isLuckyCascadeLockFish(t.Def.ID):
				g.luckyCascadeLock.tryLuckyCascadeLockFish(g, p)
			case isLuckyLegendAwakenFish(t.Def.ID):
				g.luckyLegendAwaken.tryLuckyLegendAwakenFish(g, p)
			case isLuckyCrashHarvestFish(t.Def.ID):
				g.luckyCrashHarvest.tryLuckyCrashHarvestFish(g, p)
			case isLuckyCosmicFusionFish(t.Def.ID):
				g.luckyCosmicFusion.tryLuckyCosmicFusionFish(g, p)
			// DAY-328 新增
			case isLuckyMagneticAttractionFish(t.Def.ID):
				g.luckyMagneticAttraction.tryLuckyMagneticAttractionFish(g, p)
			case isLuckySuperChainFish(t.Def.ID):
				g.luckySuperChain.tryLuckySuperChainFish(g, p)
			case isLuckyHolyPillarFish(t.Def.ID):
				g.luckyHolyPillar.tryLuckyHolyPillarFish(g, p)
			case isLuckyTimeStopFish(t.Def.ID):
				g.luckyTimeStop.tryLuckyTimeStopFish(g, p)
			case isLuckyCosmicRestartFish(t.Def.ID):
				g.luckyCosmicRestart.tryLuckyCosmicRestartFish(g, p)
			// DAY-329 新增
			case isLuckyFeverBoostUltimateFish(t.Def.ID):
				g.luckyFeverBoostUltimate.tryLuckyFeverBoostUltimateFish(g, p)
			case isLuckyRapidRichesUltimateFish(t.Def.ID):
				g.luckyRapidRichesUltimate.tryLuckyRapidRichesUltimateFish(g, p)
			case isLuckyIceFishingMasterFish(t.Def.ID):
				g.luckyIceFishingMaster.tryLuckyIceFishingMasterFish(g, p)
			case isLuckyCosmicMiracleFish(t.Def.ID):
				g.luckyCosmicMiracle.tryLuckyCosmicMiracleFish(g, p)
			case isLuckyGenesisUltimateFish(t.Def.ID):
				g.luckyGenesisUltimate.tryLuckyGenesisUltimateFish(g, p)
			// DAY-331 新增
			case isLuckySharkSparkFish(t.Def.ID):
				g.luckySharkSpark.tryLuckySharkSparkFish(g, p)
			case isLuckyWinterIceFish(t.Def.ID):
				g.luckyWinterIce.tryLuckyWinterIceFish(g, p)
			case isLuckyAtlantisFrenzyFish(t.Def.ID):
				g.luckyAtlantisFrenzy.tryLuckyAtlantisFrenzyFish(g, p)
			case isLuckyFishingTimeWheelFish(t.Def.ID):
				g.luckyFishingTimeWheel.tryLuckyFishingTimeWheelFish(g, p)
			case isLuckyUltimateSharkFish(t.Def.ID):
				g.luckyUltimateShark.tryLuckyUltimateSharkFish(g, p)
			// DAY-332 新增
			case isLuckyWildCollectorFish(t.Def.ID):
				g.luckyWildCollector.tryLuckyWildCollectorFish(g, p)
			case isLuckyLightningEelUltraFish(t.Def.ID):
				g.luckyLightningEelUltra.tryLuckyLightningEelUltraFish(g, p)
			case isLuckyDominoChainFish(t.Def.ID):
				g.luckyDominoChain.tryLuckyDominoChainFish(g, p)
			case isLuckyImmortalBossUltraFish(t.Def.ID):
				g.luckyImmortalBossUltra.tryLuckyImmortalBossUltraFish(g, p)
			case isLuckyQuadFusionFish(t.Def.ID):
				g.luckyQuadFusion.tryLuckyQuadFusionFish(g, p)
			// DAY-333 新增
			case isLuckyElectricalFrameFish(t.Def.ID):
				g.luckyElectricalFrame.tryLuckyElectricalFrameFish(g, p)
			case isLuckyMagneticRespinFish(t.Def.ID):
				g.luckyMagneticRespin.tryLuckyMagneticRespinFish(g, p)
			case isLuckyFishermanTrailFish(t.Def.ID):
				g.luckyFishermanTrail.tryLuckyFishermanTrailFish(g, p)
			case isLuckyGoldenGillsFish(t.Def.ID):
				g.luckyGoldenGills.tryLuckyGoldenGillsFish(g, p)
			case isLuckyPentaFusionFish(t.Def.ID):
				g.luckyPentaFusion.tryLuckyPentaFusionFish(g, p)
			}
			if g.luckyChainExplosion.isChainExplosionActive(playerID) {
				g.notifyChainExplosionKill(playerID, killerName, t.X, t.Y)
			}
			// 凍結期間擊破計數
			if g.luckyTimeFreeze.isTimeFreezeActive() {
				g.luckyTimeFreeze.notifyFreezeKill(playerID)
			}
			// T232 時間停止凍結期間擊破計數
			if g.luckyTimeStop.freezeActive {
				g.luckyTimeStop.notifyFreezeKill(playerID)
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
			// DAY-324 倍率梯：擊破通知
			g.luckyMultiplierLadder.onKill(g, playerID)
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

		// DAY-342 高倍率擊破全服公告（≥50x）
		if t.Multiplier >= 50 {
			playerName := p.GetDisplayName()
			var announceMsg string
			var announceColor string
			if t.Multiplier >= 1000 {
				announceMsg = "🌟 " + playerName + " 擊破 " + t.Def.Name + "！獲得 " + fmt.Sprintf("%.0fx", t.Multiplier) + " 超級大獎！"
				announceColor = "#FF00FF"
			} else if t.Multiplier >= 100 {
				announceMsg = "💫 " + playerName + " 擊破 " + t.Def.Name + "！獲得 " + fmt.Sprintf("%.0fx", t.Multiplier) + " 大獎！"
				announceColor = "#FF4500"
			} else {
				announceMsg = "✨ " + playerName + " 擊破 " + t.Def.Name + "！獲得 " + fmt.Sprintf("%.0fx", t.Multiplier) + " 獎勵！"
				announceColor = "#FFD700"
			}
			g.hub.Broadcast(protocol.MsgAnnounce, protocol.AnnouncePayload{
				Message:  announceMsg,
				Priority: "high",
				Color:    announceColor,
			})
		}

		// 勞動值滿 → 觸發 Bonus
		if laborFull && g.state == StateNormalPlay {
			// DAY-338 修復死鎖：在 tick 的鎖保護下，使用不加鎖版本
			g.triggerBonusLocked(playerID)
		}

		// DAY-345 每日任務：擊破目標計數（在鎖外呼叫，避免死鎖）
		go g.notifyQuestProgress(playerID, "kill", 1)

		// DAY-346 每週挑戰：擊破目標計數（goroutine 避免死鎖）
		isLucky := len(t.Def.ID) >= 4 && t.Def.ID[:1] == "T" && t.Def.ID[1:] >= "106" // T106+ 為 Lucky 目標物
		killMult := t.Multiplier
		go g.notifyWeeklyChallengeProgress(playerID, "kill", isLucky, killMult, 1)

		// DAY-347 賽季通行證：擊破目標 XP
		xpSource := "kill"
		xpGain := XPPerKill
		if t.Def.Type == data.TypeBoss {
			xpGain = XPPerBossKill
			xpSource = "boss"
		}
		go g.addSeasonXP(playerID, xpGain, xpSource)

		// DAY-349 成就系統：擊破觸發
		isFirstKill := g.achievementSystem.GetProgress(playerID).KillCount == 0
		go g.notifyAchievements(playerID, g.achievementSystem.OnKill(playerID, effectiveMult, isFirstKill))

		// DAY-349 好友系統：更新玩家資訊
		go g.updateFriendEntry(playerID)
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
	// DAY-339 多人投射物顯示：廣播攻擊給其他玩家（不包含自己）
	targetX := req.ClickX
	targetY := req.ClickY
	if t != nil {
		targetX = t.X
		targetY = t.Y
	}
	g.hub.BroadcastExcept(playerID, protocol.MsgOtherPlayerAttack, protocol.OtherPlayerAttackPayload{
		PlayerID:    playerID,
		CharacterID: p.GetCharacterID(),
		TargetX:     targetX,
		TargetY:     targetY,
		IsHit:       t != nil,
	})
	// DAY-338 修復死鎖：在 tick 的鎖保護下，使用不加鎖版本
	g.sendPlayerUpdateLocked(playerID)
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

	// DAY-349 成就系統：BOSS 擊破觸發
	go g.notifyAchievements(playerID, g.achievementSystem.OnBossKill(playerID))
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

// triggerBonusLocked — 不加鎖版本，供 tick 內部呼叫（已持有 g.mu.Lock()）
// DAY-338 修復死鎖：handleAttackLocked 在 tick 的鎖保護下，不能再次加鎖
func (g *Game) triggerBonusLocked(playerID string) {
	if g.state != StateNormalPlay {
		return
	}
	p, ok := g.players[playerID]
	if !ok {
		return
	}
	p.EntryBetCost = p.GetBetDef().BetCost
	p.ResetLabor()
	g.bonusPlayer = playerID
	g.bonusTimer = 15.0
	g.bonusScore = make(map[string]int)
	g.bonusWeedHP = make(map[string]int)
	g.setState(StateBonusGame)

	g.hub.Broadcast(protocol.MsgBonusEvent, protocol.BonusEventPayload{
		Event:    "start",
		TimeLeft: 15.0,
	})
	g.sendPlayerUpdateLocked(playerID)
	log.Printf("[Game] Bonus triggered by %s (locked)", playerID)

	// DAY-345 每日任務：Bonus 觸發計數（goroutine 避免死鎖）
	go g.notifyQuestProgress(playerID, "bonus", 1)

	// DAY-346 每週挑戰：Bonus 觸發計數（goroutine 避免死鎖）
	go g.notifyWeeklyChallengeProgress(playerID, "bonus", false, 0.0, 1)

	// DAY-347 賽季通行證：完成 Bonus XP
	go g.addSeasonXP(playerID, XPPerBonus, "bonus")

	// DAY-349 成就系統：Bonus 完成觸發
	go g.notifyAchievements(playerID, g.achievementSystem.OnBonusComplete(playerID))
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
	// DAY-338 修復死鎖：autoFire 在 tick 的鎖保護下執行，直接呼叫 handleAttackLocked
	g.handleAttackLocked(p.ID, protocol.AttackRequest{
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
	g.sendPlayerUpdateWithPlayer(playerID, p)
}

// sendPlayerUpdateLocked — 不加鎖版本，供 tick 內部呼叫（已持有 g.mu.Lock()）
// DAY-338 修復死鎖：handleAttackLocked 在 tick 的鎖保護下，不能再次加鎖
func (g *Game) sendPlayerUpdateLocked(playerID string) {
	p, ok := g.players[playerID]
	if !ok {
		return
	}
	g.sendPlayerUpdateWithPlayer(playerID, p)
}

// sendPlayerUpdateWithPlayer — 實際發送邏輯（不加鎖）
func (g *Game) sendPlayerUpdateWithPlayer(playerID string, p *Player) {
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
		// DAY-342 在線玩家數
		OnlineCount:     len(g.players),
		// DAY-344 玩家顯示名稱
		DisplayName:     p.GetDisplayName(),
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

// ── DAY-345 每日任務系統 Handler ─────────────────────────────

// handleDailyQuestRequest 處理玩家請求任務狀態
func (g *Game) handleDailyQuestRequest(playerID string) {
	status := g.dailyQuest.GetPlayerStatus(playerID)
	payload := protocol.DailyQuestUpdatePayload{
		QuestCoins: status.QuestCoins,
		ResetAt:    status.ResetAt,
	}
	for _, q := range status.Quests {
		payload.Quests = append(payload.Quests, protocol.QuestStatusPayload{
			ID:          q.ID,
			Name:        q.Name,
			Description: q.Description,
			Target:      q.Target,
			Progress:    q.Progress,
			Completed:   q.Completed,
			Claimed:     q.Claimed,
			Reward:      q.Reward,
		})
	}
	g.hub.Send(playerID, protocol.MsgDailyQuestUpdate, payload)
}

// handleDailyQuestClaim 處理玩家領取任務獎勵
func (g *Game) handleDailyQuestClaim(playerID string, questID string) {
	coins := g.dailyQuest.ClaimReward(playerID, questID)
	if coins < 0 {
		g.hub.Send(playerID, protocol.MsgError, map[string]string{
			"message": "任務未完成或已領取",
		})
		return
	}
	log.Printf("[DailyQuest] Player %s claimed quest %s, earned %d coins", playerID, questID, coins)
	// 發送更新後的任務狀態
	g.handleDailyQuestRequest(playerID)
	// 公告
	g.hub.Send(playerID, protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  fmt.Sprintf("🎯 任務完成！獲得 %d 任務幣", coins),
		Priority: "normal",
		Color:    "#FFD700",
	})
}

// notifyQuestProgress 通知任務進度（擊破目標時呼叫）
func (g *Game) notifyQuestProgress(playerID string, eventType string, value int) {
	var completed bool
	var questName string
	var reward int

	switch eventType {
	case "kill":
		completed, questName, reward = g.dailyQuest.OnKillTarget(playerID)
	case "combo":
		completed, questName, reward = g.dailyQuest.OnComboReach(playerID, value)
	case "bonus":
		completed, questName, reward = g.dailyQuest.OnTriggerBonus(playerID)
	}

	if completed {
		// 通知任務完成
		g.hub.Send(playerID, protocol.MsgDailyQuestComplete, protocol.DailyQuestCompletePayload{
			QuestID:   questName,
			QuestName: questName,
			Reward:    reward,
			Message:   fmt.Sprintf("🎯 任務完成：%s！點擊領取 %d 任務幣", questName, reward),
		})
		log.Printf("[DailyQuest] Player %s completed quest: %s (reward: %d)", playerID, questName, reward)
		// 同步更新任務狀態
		g.handleDailyQuestRequest(playerID)

		// DAY-349 成就系統：每日任務完成觸發
		go g.notifyAchievements(playerID, g.achievementSystem.OnQuestComplete(playerID))
	}
}

// ── DAY-349 成就系統 Handler ──────────────────────────────────

// notifyAchievements 通知玩家新解鎖的成就，並發放獎勵
func (g *Game) notifyAchievements(playerID string, defs []*AchievementDef) {
	if len(defs) == 0 {
		return
	}
	g.mu.Lock()
	p, ok := g.players[playerID]
	if !ok {
		g.mu.Unlock()
		return
	}
	totalReward := 0
	for _, def := range defs {
		totalReward += def.Reward
	}
	if totalReward > 0 {
		p.AddCoins(totalReward)
	}
	g.mu.Unlock()

	for _, def := range defs {
		g.hub.Send(playerID, protocol.MsgAchievementUnlock, protocol.AchievementUnlockPayload{
			ID:          string(def.ID),
			Name:        def.Name,
			Description: def.Description,
			Icon:        def.Icon,
			Rarity:      def.Rarity,
			Reward:      def.Reward,
		})
		log.Printf("[Achievement] Player %s unlocked: %s (%s) +%d coins", playerID, def.Name, def.Rarity, def.Reward)
	}
}

// handleAchievementListRequest 處理玩家請求成就列表
func (g *Game) handleAchievementListRequest(playerID string) {
	unlocked := g.achievementSystem.GetAllUnlocked(playerID)
	unlockedMap := make(map[AchievementType]bool)
	for _, def := range unlocked {
		unlockedMap[def.ID] = true
	}

	var entries []*protocol.AchievementEntryPayload
	for _, def := range achievementDefs {
		entries = append(entries, &protocol.AchievementEntryPayload{
			ID:          string(def.ID),
			Name:        def.Name,
			Description: def.Description,
			Icon:        def.Icon,
			Rarity:      def.Rarity,
			Reward:      def.Reward,
			Unlocked:    unlockedMap[def.ID],
		})
	}

	g.hub.Send(playerID, protocol.MsgAchievementList, protocol.AchievementListPayload{
		Achievements:  entries,
		TotalCount:    len(achievementDefs),
		UnlockedCount: len(unlocked),
	})
}

// handleRoomLeaderboardRequest 處理玩家請求同場排行榜
func (g *Game) handleRoomLeaderboardRequest(playerID string) {
	entries := g.friendSystem.GetRoomLeaderboard()
	myRank := g.friendSystem.GetPlayerRank(playerID)

	var payload []*protocol.RoomLeaderboardEntryPayload
	for i, e := range entries {
		payload = append(payload, &protocol.RoomLeaderboardEntryPayload{
			Rank:        i + 1,
			PlayerID:    e.PlayerID,
			DisplayName: e.DisplayName,
			SeasonXP:    e.SeasonXP,
			TotalKills:  e.TotalKills,
			BestMult:    e.BestMult,
			IsOnline:    e.IsOnline,
		})
	}

	g.hub.Send(playerID, protocol.MsgRoomLeaderboard, protocol.RoomLeaderboardPayload{
		Entries:     payload,
		MyRank:      myRank,
		OnlineCount: g.friendSystem.GetOnlineCount(),
	})
}

// updateFriendEntry 更新好友系統中的玩家資訊
func (g *Game) updateFriendEntry(playerID string) {
	g.mu.RLock()
	p, ok := g.players[playerID]
	if !ok {
		g.mu.RUnlock()
		return
	}
	displayName := p.GetDisplayName()
	coins := p.Coins
	g.mu.RUnlock()

	// 取得賽季 XP
	seasonXP := 0
	if g.seasonPass != nil {
		state := g.seasonPass.GetOrCreateState(playerID)
		seasonXP = state.CurrentXP
	}

	// 取得累積擊破數（從成就進度）
	prog := g.achievementSystem.GetProgress(playerID)
	totalKills := prog.KillCount

	// 取得最高倍率（從成就解鎖記錄推算）
	bestMult := 0.0
	if g.achievementSystem.isUnlocked(playerID, AchievTypeMult1000) {
		bestMult = 1000.0
	} else if g.achievementSystem.isUnlocked(playerID, AchievTypeMult500) {
		bestMult = 500.0
	} else if g.achievementSystem.isUnlocked(playerID, AchievTypeMult100) {
		bestMult = 100.0
	} else if g.achievementSystem.isUnlocked(playerID, AchievTypeMult50) {
		bestMult = 50.0
	}

	g.friendSystem.UpdatePlayer(playerID, displayName, seasonXP, totalKills, bestMult)

	// 金幣成就觸發
	go g.notifyAchievements(playerID, g.achievementSystem.OnCoinsUpdate(playerID, coins))
}

// ── DAY-346 每週挑戰系統 Handler ─────────────────────────────

// handleWeeklyChallengeRequest 處理玩家請求每週挑戰狀態
func (g *Game) handleWeeklyChallengeRequest(playerID string) {
	status := g.weeklyChallenge.GetPlayerStatus(playerID)
	payload := protocol.WeeklyChallengeUpdatePayload{
		WeeklyCoins: status.WeeklyCoins,
		WeekKey:     status.WeekKey,
		ResetAt:     status.ResetAt,
	}
	for _, c := range status.Challenges {
		payload.Challenges = append(payload.Challenges, protocol.ChallengeStatusPayload{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			Target:      c.Target,
			Progress:    c.Progress,
			Completed:   c.Completed,
			Claimed:     c.Claimed,
			Reward:      c.Reward,
			Tier:        c.Tier,
		})
	}
	g.hub.Send(playerID, protocol.MsgWeeklyChallengeUpdate, payload)
}

// handleWeeklyChallengeClaim 處理玩家領取每週挑戰獎勵
func (g *Game) handleWeeklyChallengeClaim(playerID string, challengeID string) {
	coins := g.weeklyChallenge.ClaimReward(playerID, challengeID)
	if coins < 0 {
		g.hub.Send(playerID, protocol.MsgError, map[string]string{
			"message": "挑戰未完成或已領取",
		})
		return
	}
	log.Printf("[WeeklyChallenge] Player %s claimed challenge %s, earned %d coins", playerID, challengeID, coins)
	// 發送更新後的挑戰狀態
	g.handleWeeklyChallengeRequest(playerID)
	// 公告
	g.hub.Send(playerID, protocol.MsgAnnounce, protocol.AnnouncePayload{
		Message:  fmt.Sprintf("🏆 週間挑戰完成！獲得 %d 任務幣", coins),
		Priority: "high",
		Color:    "#FF8C00",
	})
}

// notifyWeeklyChallengeProgress 通知每週挑戰進度
// eventType: "kill"（isLucky, mult）/ "combo"（value）/ "bonus"
func (g *Game) notifyWeeklyChallengeProgress(playerID string, eventType string, isLucky bool, mult float64, value int) {
	var completedList []WeeklyChallengeComplete

	switch eventType {
	case "kill":
		completedList = g.weeklyChallenge.OnKillTarget(playerID, isLucky, mult)
	case "combo":
		completedList = g.weeklyChallenge.OnComboReach(playerID, value)
	case "bonus":
		completedList = g.weeklyChallenge.OnTriggerBonus(playerID)
	}

	for _, c := range completedList {
		tierEmoji := "🥉"
		if c.Tier == 2 {
			tierEmoji = "🥈"
		} else if c.Tier == 3 {
			tierEmoji = "🥇"
		}
		g.hub.Send(playerID, protocol.MsgWeeklyChallengeComplete, protocol.WeeklyChallengeCompletePayload{
			ChallengeID:   c.ChallengeID,
			ChallengeName: c.ChallengeName,
			Reward:        c.Reward,
			Tier:          c.Tier,
			Message:       fmt.Sprintf("%s 週間挑戰完成：%s！點擊領取 %d 任務幣", tierEmoji, c.ChallengeName, c.Reward),
		})
		log.Printf("[WeeklyChallenge] Player %s completed challenge: %s (tier: %d, reward: %d)", playerID, c.ChallengeName, c.Tier, c.Reward)
		// 同步更新挑戰狀態
		g.handleWeeklyChallengeRequest(playerID)
	}
}

// ── DAY-347 賽季通行證系統 Handler ───────────────────────────

// addSeasonXP 增加賽季 XP 並通知玩家
func (g *Game) addSeasonXP(playerID string, xp int, source string) {
	levelUp, newLevel, newXP := g.seasonPass.AddXP(playerID, xp)

	// 計算下一等級所需 XP
	nextLevelXP := -1
	if newLevel < len(g.seasonPass.tiers) {
		nextLevelXP = g.seasonPass.tiers[newLevel].RequiredXP
	}

	// 計算剩餘天數
	daysLeft := int(time.Until(g.seasonPass.GetSeasonEnd()).Hours() / 24)

	// 發送賽季狀態更新
	g.hub.Send(playerID, protocol.MsgSeasonPassUpdate, protocol.SeasonPassUpdatePayload{
		CurrentXP:    newXP,
		CurrentLevel: newLevel,
		NextLevelXP:  nextLevelXP,
		IsPremium:    false,
		SeasonID:     g.seasonPass.GetSeasonID(),
		DaysLeft:     daysLeft,
		XPGained:     xp,
		XPSource:     source,
	})

	// 升級通知
	if levelUp && newLevel >= 1 && newLevel <= len(g.seasonPass.tiers) {
		tier := g.seasonPass.tiers[newLevel-1]
		g.hub.Send(playerID, protocol.MsgSeasonPassLevelUp, protocol.SeasonPassLevelUpPayload{
			NewLevel:      newLevel,
			LevelName:     tier.Name,
			BadgeName:     tier.BadgeName,
			FreeReward:    tier.FreeReward,
			PremiumReward: tier.PremiumReward,
			IsPremium:     false,
		})
		log.Printf("[SeasonPass] Player %s leveled up to %d (%s)", playerID, newLevel, tier.Name)
	}
}

// ── DAY-348 任務幣兌換商店 + 賽季排行榜 Handler ──────────────

// handleShopRequest 處理玩家請求商店資訊
func (g *Game) handleShopRequest(playerID string) {
	// 取得玩家任務幣（每日 + 每週）
	questCoins := g.dailyQuest.GetQuestCoins(playerID) + g.weeklyChallenge.GetWeeklyCoins(playerID)

	// 轉換商品列表
	items := g.questShop.GetItems()
	itemPayloads := make([]protocol.ShopItemPayload, len(items))
	for i, item := range items {
		itemPayloads[i] = protocol.ShopItemPayload{
			ID:          item.ID,
			Type:        string(item.Type),
			Name:        item.Name,
			Description: item.Description,
			Cost:        item.Cost,
			Value:       item.Value,
			Icon:        item.Icon,
		}
	}

	// 取得有效效果
	effects := g.questShop.GetActiveEffectsSummary(playerID)

	g.hub.Send(playerID, protocol.MsgShopItems, protocol.ShopItemsPayload{
		Items:      itemPayloads,
		QuestCoins: questCoins,
		Effects:    effects,
	})
}

// handleShopPurchase 處理玩家購買道具
func (g *Game) handleShopPurchase(playerID string, itemID string) {
	// 取得玩家任務幣（每日 + 每週合計）
	dailyCoins := g.dailyQuest.GetQuestCoins(playerID)
	weeklyCoins := g.weeklyChallenge.GetWeeklyCoins(playerID)
	totalCoins := dailyCoins + weeklyCoins

	// 嘗試購買
	success, errMsg, item := g.questShop.Purchase(playerID, itemID, totalCoins)
	if !success {
		g.hub.Send(playerID, protocol.MsgShopPurchaseResult, protocol.ShopPurchaseResultPayload{
			Success: false,
			ItemID:  itemID,
			Message: errMsg,
			QuestCoins: totalCoins,
		})
		return
	}

	// 扣除任務幣（優先扣每日，不足再扣每週）
	cost := item.Cost
	if dailyCoins >= cost {
		g.dailyQuest.SpendQuestCoins(playerID, cost)
	} else {
		g.dailyQuest.SpendQuestCoins(playerID, dailyCoins)
		g.weeklyChallenge.SpendQuestCoins(playerID, cost-dailyCoins)
	}

	// 計算剩餘任務幣
	remainingCoins := g.dailyQuest.GetQuestCoins(playerID) + g.weeklyChallenge.GetWeeklyCoins(playerID)

	// 立即生效的獎勵（CoinBonus 類型）
	coinReward := 0
	if item.Type == ShopItemCoinBonus {
		// 立即給予金幣
		g.mu.Lock()
		p, ok := g.players[playerID]
		if ok {
			p.AddCoins(item.Value)
			coinReward = item.Value
			g.sendPlayerUpdateWithPlayer(playerID, p)
		}
		g.mu.Unlock()
		// 標記為已使用
		g.questShop.ConsumeEffect(playerID, ShopItemCoinBonus)
	}

	g.hub.Send(playerID, protocol.MsgShopPurchaseResult, protocol.ShopPurchaseResultPayload{
		Success:    true,
		ItemID:     item.ID,
		ItemName:   item.Name,
		Cost:       item.Cost,
		Message:    fmt.Sprintf("成功購買「%s」！", item.Name),
		QuestCoins: remainingCoins,
		CoinReward: coinReward,
	})

	// 發送效果更新
	effects := g.questShop.GetActiveEffectsSummary(playerID)
	g.hub.Send(playerID, protocol.MsgShopEffectUpdate, protocol.ShopEffectUpdatePayload{
		Effects: effects,
	})

	log.Printf("[Shop] Player %s purchased %s (cost: %d, remaining coins: %d)", playerID, item.Name, item.Cost, remainingCoins)
}

// handleSeasonLeaderboardRequest 處理玩家請求賽季排行榜
func (g *Game) handleSeasonLeaderboardRequest(playerID string) {
	// 先更新當前玩家的排行榜資料
	g.mu.RLock()
	p, ok := g.players[playerID]
	if ok {
		state := g.seasonPass.GetSnapshot(playerID)
		xp, _ := state["current_xp"].(int)
		level, _ := state["current_level"].(int)
		levelName := ""
		badge := ""
		if level >= 1 && level <= len(g.seasonPass.tiers) {
			tier := g.seasonPass.tiers[level-1]
			levelName = tier.Name
			badge = tier.BadgeName
		}
		displayName := p.GetDisplayName()
		g.mu.RUnlock()
		g.seasonLeaderboard.UpdatePlayer(playerID, displayName, xp, level, levelName, badge)
	} else {
		g.mu.RUnlock()
	}

	// 取得排行榜快照
	snapshot := g.seasonLeaderboard.GetSnapshot(playerID)

	// 轉換為 Payload
	top20Raw, _ := snapshot["top20"].([]LeaderboardEntry)
	top20 := make([]protocol.LeaderboardEntryPayload, len(top20Raw))
	for i, e := range top20Raw {
		top20[i] = protocol.LeaderboardEntryPayload{
			Rank:        e.Rank,
			PlayerID:    e.PlayerID,
			DisplayName: e.DisplayName,
			SeasonXP:    e.SeasonXP,
			Level:       e.Level,
			LevelName:   e.LevelName,
			Badge:       e.Badge,
		}
	}

	myRank, _ := snapshot["my_rank"].(int)
	payload := protocol.SeasonLeaderboardPayload{
		SeasonID:    g.seasonLeaderboard.seasonID,
		Top20:       top20,
		MyRank:      myRank,
		LastUpdated: time.Now().UnixMilli(),
	}

	if myEntry, ok := snapshot["my_entry"].(*LeaderboardEntry); ok && myEntry != nil {
		e := protocol.LeaderboardEntryPayload{
			Rank:        myEntry.Rank,
			PlayerID:    myEntry.PlayerID,
			DisplayName: myEntry.DisplayName,
			SeasonXP:    myEntry.SeasonXP,
			Level:       myEntry.Level,
			LevelName:   myEntry.LevelName,
			Badge:       myEntry.Badge,
		}
		payload.MyEntry = &e
	}

	g.hub.Send(playerID, protocol.MsgSeasonLeaderboard, payload)
}
