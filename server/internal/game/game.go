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
	"digital-twin/server/internal/game/activityfeed"
	"digital-twin/server/internal/game/achievement"
	"digital-twin/server/internal/game/awakenboss"
	"digital-twin/server/internal/game/chain"
	"digital-twin/server/internal/game/challenge"
	"digital-twin/server/internal/game/combat"
	"digital-twin/server/internal/game/dailyboss"
	"digital-twin/server/internal/game/dm"
	"digital-twin/server/internal/game/fragment"
	"digital-twin/server/internal/game/flashchallenge"
	"digital-twin/server/internal/game/goldentime"
	"digital-twin/server/internal/game/immortalboss"
	"digital-twin/server/internal/game/mysterybox"
	"digital-twin/server/internal/game/rarecatch"
	"digital-twin/server/internal/game/respin"
	"digital-twin/server/internal/game/treasuremap"
	"digital-twin/server/internal/game/specialweapon"
	"digital-twin/server/internal/game/dailyspin"
	"digital-twin/server/internal/game/event"
	"digital-twin/server/internal/game/friend"
	"digital-twin/server/internal/game/friendchallenge"
	"digital-twin/server/internal/game/guild"
	"digital-twin/server/internal/game/guildwar"
	"digital-twin/server/internal/game/jackpot"
	"digital-twin/server/internal/game/mission"
	"digital-twin/server/internal/game/season"
	"digital-twin/server/internal/game/shop"
	"digital-twin/server/internal/game/announce"
	"digital-twin/server/internal/game/state"
	"digital-twin/server/internal/game/target"
	"digital-twin/server/internal/game/tournament"
	"digital-twin/server/internal/game/vip"
	"digital-twin/server/internal/game/referral"
	"digital-twin/server/internal/game/weather"
	"digital-twin/server/internal/game/wheel"
	"digital-twin/server/internal/game/winstreak"
	"digital-twin/server/internal/game/lightningeel"
	"digital-twin/server/internal/game/fevermode"
	"digital-twin/server/internal/anticheat"
	"digital-twin/server/internal/game/festival"
	"digital-twin/server/internal/game/halloffame"
	"digital-twin/server/internal/game/roulette"
	raidboss "digital-twin/server/internal/game/raidBoss"
	"digital-twin/server/internal/game/unlucky"
	"digital-twin/server/internal/game/speedrace"
	"digital-twin/server/internal/game/bounty"
	"digital-twin/server/internal/game/multstorm"
	"digital-twin/server/internal/game/dualroulette"
	"digital-twin/server/internal/game/megacatch"
	"digital-twin/server/internal/game/megaoctopus"
	"digital-twin/server/internal/game/chainlongwheel"
	"digital-twin/server/internal/game/crystaldragon"
	"digital-twin/server/internal/game/roulettecrab"
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
	dailyTournamentMgr *tournament.DailyTournament // 每日賽管理器（DAY-093）
	multiFormatMgr *tournament.MultiFormatTournament // 多格式每日賽管理器（DAY-111）
	Season      *season.Manager     // 賽季通行證管理器（DAY-072）
	Friends     *friend.Manager     // 好友系統管理器（DAY-073）
	FriendChallenge *friendchallenge.Manager // 好友挑戰管理器（DAY-102）
	DM          *dm.Manager         // 私訊管理器（DAY-103）
	Guild       *guild.Manager      // 公會系統管理器（DAY-074）
	GuildWar    *guildwar.Manager   // 公會戰管理器（DAY-076）
	DailyBoss   *dailyboss.Manager  // 每日 BOSS 挑戰管理器（DAY-077）
	VIP         *vip.Manager        // VIP 等級管理器（DAY-078）
	Event       *event.Manager      // 限時活動管理器（DAY-079）
	Referral    *referral.Manager   // 推薦碼管理器（DAY-082）
	Wheel       *wheel.Manager      // 幸運轉盤管理器（DAY-084）
	Challenge   *challenge.Manager  // 隱藏挑戰管理器（DAY-085）
	Weather     *weather.Manager    // 天氣系統管理器（DAY-087）
	Chain       *chain.Manager      // 連鎖爆炸管理器（DAY-088）
	SpecialWeapon *specialweapon.Manager // 特殊武器管理器（DAY-089）
	MysteryBox    *mysterybox.Manager    // 神秘寶箱管理器（DAY-090）
	DailySpin     *dailyspin.Manager     // 每日簽到轉盤管理器（DAY-092）
	Shop          *shop.Manager          // 商店管理器（DAY-094）
	Announce      *announce.Manager      // 全服公告管理器（DAY-097）
	AntiCheat     *anticheat.Manager     // 異常行為偵測管理器（DAY-105）
	Festival      *festival.Manager      // 賽季節日活動管理器（DAY-109）
	HallOfFame    *halloffame.Manager    // 全服名人堂管理器（DAY-110）
	ActivityFeed  *activityfeed.Manager  // 成就動態牆管理器（DAY-112）
	Roulette      *roulette.Manager      // 雙層倍率輪盤管理器（DAY-113）
	RaidBoss      *raidboss.Manager      // Co-op Boss Raid 管理器（DAY-115）
	Fragment      *fragment.Manager      // 碎片收集大獎管理器（DAY-116）
	RespinMgr     *respin.Manager        // Rapid Respin 管理器（DAY-121）
	TreasureMap   *treasuremap.Manager   // 寶藏地圖管理器（DAY-122）
	FlashChallenge *flashchallenge.Manager // 閃電挑戰管理器（DAY-123）
	GoldenTime     *goldentime.Manager     // 黃金時間管理器（DAY-125）
	RareCatch      *rarecatch.Manager      // 稀有連擊累積倍率管理器（DAY-126）
	ImmortalBoss   *immortalboss.Manager   // 不死 BOSS 連勝管理器（DAY-129）
	AwakenBoss     *awakenboss.Manager     // 覺醒 BOSS 系統管理器（DAY-130）
	WinStreak      *winstreak.Manager      // 連勝獎勵系統管理器（DAY-131）
	LightningEel   *lightningeel.Manager   // 閃電鰻連鎖攻擊管理器（DAY-132）
	FeverMode      *fevermode.Manager      // 狂熱模式管理器（DAY-133）
	UnluckyBonus   *unlucky.Manager        // 失敗補償系統管理器（DAY-135）
	SpeedRace      *speedrace.Manager      // 競速獵殺系統管理器（DAY-136）
	Bounty         *bounty.Manager         // 全服目標懸賞系統管理器（DAY-137）
	MultStorm      *multstorm.Manager      // 全服倍率風暴系統管理器（DAY-138）
	DualRoulette   *dualroulette.Manager   // 雙環輪盤系統管理器（DAY-139）
	MegaCatch      *megacatch.Manager      // 全服 Mega Catch 事件系統管理器（DAY-140）
	MegaOctopus    *megaoctopus.Manager    // 巨型章魚轉盤系統管理器（DAY-144）
	GiantPrizeFish *giantPrizeFishManager  // 夢幻巨型獎勵魚系統管理器（DAY-147）
	ChainLongWheel *chainlongwheel.Manager // 千龍王強化輪盤系統管理器（DAY-148）
	ThunderboltLobster *thunderboltLobsterManager // 雷霆龍蝦免費射擊系統管理器（DAY-150）
	RainbowPhoenix     *rainbowPhoenixManager     // 彩虹鳳凰 Power Up 系統管理器（DAY-151）
	GoldenTurtle       *goldenTurtleManager       // 黃金海龜時間停止系統管理器（DAY-159）
	LuckyStarFish      *luckyStarFishManager      // 幸運星魚全場倍率翻倍管理器（DAY-160）
	GoldenShark        *goldenSharkManager        // 黃金鯊魚全服狂暴模式管理器（DAY-161）
	CaptainFish        *captainFishManager        // 船長魚全服競速模式管理器（DAY-163）
	AbyssWhale         *abyssWhaleManager         // 深淵巨鯨全服 Boss 挑戰管理器（DAY-164）
	RouletteCrab       *roulettecrab.Manager      // 黃金輪盤螃蟹系統管理器（DAY-167）
	CrystalDragon      *crystaldragon.Manager     // 水晶龍收集大獎系統管理器（DAY-153）
	LionDance          *lionDanceManager          // 獅子舞大獎爆發系統管理器（DAY-168）
	VortexFish         *vortexFishManager         // 漩渦魚群吸引系統管理器（DAY-169）
	FreezeBomb         *freezeBombManager         // 冰凍炸彈魚系統管理器（DAY-170）
	IceFishing         *iceFishingManager         // 冰釣幸運輪盤系統管理器（DAY-171）
	LuckyEgg           *luckyEggManager           // 幸運彩蛋魚系統管理器（DAY-172）
	RainbowLucky       *rainbowLuckyManager       // 彩虹幸運魚系統管理器（DAY-173）
	// DAY-174 海葵觸手攻擊系統：無需管理器（stateless，每次擊破獨立觸發）
	LuckyDice          *luckyDiceManager          // 幸運骰子魚系統管理器（DAY-175）
	FireStorm          *fireStormManager          // 火焰風暴魚系統管理器（DAY-176）
	GoldenTreasure     *goldenTreasureManager     // 黃金寶藏魚系統管理器（DAY-177）
	Mermaid            *mermaidManager            // 美人魚治癒系統管理器（DAY-178）
	LuckyClover        *luckyCloverManager        // 幸運草魚系統管理器（DAY-179）
	RainbowShark       *rainbowSharkManager       // 彩虹鯊魚爆發系統管理器（DAY-180）
	// DAY-181 雷霆鯊魚連鎖閃電系統：無需管理器（stateless，每次擊破獨立觸發）
	VampireFish        *vampireFishManager        // 吸血鬼魚累積倍率系統管理器（DAY-182）
	LightningAutoChain *lightningAutoChainManager // 閃電魚自動連鎖系統管理器（DAY-183）
	MeteorFish         *meteorFishManager         // 隕石魚隕石雨系統管理器（DAY-184）
	PhoenixFish        *phoenixFishManager        // 鳳凰魚涅槃重生系統管理器（DAY-185）
	DragonTurtle       *dragonTurtleManager       // 龍龜不死 Boss 系統管理器（DAY-186）
	CrocodileHunter    *crocodileHunterManager    // 巨型鱷魚獵食系統管理器（DAY-188）
	TimeBomb           *timeBombManager           // 時間炸彈魚系統管理器（DAY-189）
	TripleLucky        *tripleLuckyFishManager    // 三重幸運魚系統管理器（DAY-190）
	SchoolPanic        *schoolPanicManager        // 魚群驚嚇連帶系統管理器（DAY-191）
	RockSkeleton       *rockSkeletonManager       // 搖滾骷髏演唱會系統管理器（DAY-192）
	ChainLongKing      *chainLongKingManager      // 長龍王雙環輪盤系統管理器（DAY-194）
	DrillLobster       *drillLobsterManager       // 鑽頭龍蝦穿透爆炸系統管理器（DAY-195）
	AnglerfishElectric *anglerfishManager         // 巨型鮟鱇魚電擊寶箱系統管理器（DAY-196）
	MysticDragon       *mysticDragonManager       // 神秘龍魚八波攻擊系統管理器（DAY-197）
	GhostFish          *ghostFishManager          // 幽靈魚分身系統管理器（DAY-198）
	ThunderboltLobsterV2 *tbLobsterV2Manager       // 雷霆龍蝦免費射擊系統管理器（DAY-199）
	IcePhoenix         *icePhoenixManager          // 冰鳳凰覺醒 BOSS 系統管理器（DAY-200）
	SerialBombCrab     *serialBombCrabManager      // 連環炸彈蟹系統管理器（DAY-201）
	AbyssVortex        *abyssVortexManager         // 深淵漩渦魚系統管理器（DAY-202）
	HumpbackWhale      *humpbackWhaleManager       // 座頭鯨覺醒系統管理器（DAY-203）
	FreeSpinFish       *freeSpinFishManager        // 自由旋轉魚免費射擊系統管理器（DAY-204）
	JackpotDragon      *jackpotDragonManager       // 獎池龍 Jackpot 抽獎系統管理器（DAY-205）
	CometFish          *cometFishManager           // 彗星魚連鎖爆炸系統管理器（DAY-206）
	GoldenWaveFish     *goldenWaveFishManager      // 黃金波浪魚全場倍率衝擊系統管理器（DAY-207）
	DragonKing         *dragonKingManager          // 深海龍王全服合力蓄力系統管理器（DAY-208）
	FortuneCoinFish    *fortuneCoinFishManager     // 幸運金幣魚即時獎勵系統管理器（DAY-209）
	LuckyHotZone       *luckyHotZoneManager        // 幸運熱區魚空間策略系統管理器（DAY-210）
	LuckyTrident       *luckyTridentManager        // 幸運三叉魚互動三轉盤系統管理器（DAY-211）
	TimeFreeze         *timeFreezeManager          // 時間凍結魚系統管理器（DAY-212）
	RainbowPrism       *rainbowPrismManager        // 彩虹稜鏡魚系統管理器（DAY-213）
	GoldenAccumulator  *goldenAccumulatorManager   // 黃金累積魚系統管理器（DAY-214）
	LuckyMirrorFish    *luckyMirrorFishManager     // 幸運鏡像魚系統管理器（DAY-215）
	CursedPoisonFish   *cursedPoisonFishManager    // 詛咒毒魚系統管理器（DAY-216）
	LuckyAuctionFish   *luckyAuctionFishManager    // 幸運拍賣魚系統管理器（DAY-217）
	LuckyEvolutionFish *luckyEvolutionFishManager  // 幸運進化魚系統管理器（DAY-218）
	LuckyInfectionFish *luckyInfectionFishManager  // 幸運連鎖感染魚系統管理器（DAY-219）
	LuckyRicochetFish  *luckyRicochetFishManager   // 幸運反彈魚系統管理器（DAY-220）
	LuckyBlackHole     *luckyBlackHoleManager      // 幸運黑洞魚系統管理器（DAY-221）
	LuckyResonanceFish *luckyResonanceFishManager  // 幸運共鳴魚系統管理器（DAY-222）
	LuckyTeleportFish  *luckyTeleportFishManager   // 幸運傳送魚系統管理器（DAY-223）
	LuckySplitFish     *luckySplitFishManager      // 幸運分裂魚系統管理器（DAY-224）
	LuckyChargeFish    *luckyChargeFishManager     // 幸運充能魚系統管理器（DAY-225）
	LuckyChainBomb     *luckyChainBombManager      // 幸運鏈鎖爆炸魚系統管理器（DAY-226）
	LuckyMirrorTime    *luckyMirrorTimeManager     // 幸運鏡像時空魚系統管理器（DAY-227）
	LuckyQuantumFish   *luckyQuantumFishManager    // 幸運量子魚系統管理器（DAY-228）
	LuckyParasiteFish  *luckyParasiteFishManager   // 幸運寄生魚系統管理器（DAY-229）
	LuckyStormFish     *luckyStormFishManager      // 幸運風暴魚系統管理器（DAY-230）
	LuckyBoomerangFish *luckyBoomerangFishManager  // 幸運迴旋鏢魚系統管理器（DAY-231）
	LuckyMagnetFish    *luckyMagnetFishManager     // 幸運磁力魚系統管理器（DAY-232）
	LuckyEchoFish      *luckyEchoFishManager       // 幸運回聲魚系統管理器（DAY-233）
	LuckyVortexFish    *luckyVortexFishManager     // 幸運漩渦魚系統管理器（DAY-234）
	LuckyTimeBombFish  *luckyTimeBombFishManager   // 幸運時間炸彈魚系統管理器（DAY-235）
	LuckyMirrorWorld   *luckyMirrorWorldManager    // 幸運鏡面世界魚系統管理器（DAY-236）
	LuckyFreezeWorld   *luckyFreezeWorldManager    // 幸運冰凍世界魚系統管理器（DAY-237）
	LuckyGravityFlip   *luckyGravityFlipManager    // 幸運重力反轉魚系統管理器（DAY-238）
	LuckySynergyBurst  *luckySynergyBurstManager   // 幸運共鳴爆發魚系統管理器（DAY-239）
	LuckyBetFish       *luckyBetFishManager         // 幸運賭注魚系統管理器（DAY-240）
	LuckyChainReaction *luckyChainReactionManager   // 幸運連鎖反應魚系統管理器（DAY-241）
	LuckyCloneFish     *luckyCloneFishManager        // 幸運分身魚系統管理器（DAY-242）
	LuckyProphecyFish  *luckyProphecyFishManager     // 幸運預言魚系統管理器（DAY-243）
	LuckyFlagFish      *luckyFlagFishManager          // 幸運奪旗魚系統管理器（DAY-244）
	LuckyPhantomFish   *luckyPhantomFishManager        // 幸運幽靈魚系統管理器（DAY-245）
	LuckyCrystalBallFish *luckyCrystalBallFishManager   // 幸運水晶球魚系統管理器（DAY-246）
	LuckyTimeRewind      *luckyTimeRewindManager          // 幸運時光倒流魚系統管理器（DAY-247）
	LuckyTornado         *luckyTornadoManager              // 幸運龍捲風魚系統管理器（DAY-248）
	LuckyBlackHoleExplosion *luckyBlackHoleExplosionManager  // 幸運黑洞爆炸魚系統管理器（DAY-249）
	LuckyMirrorSplit        *luckyMirrorSplitManager          // 幸運鏡像分裂魚系統管理器（DAY-250）
	LuckyQuantumEntangle    *luckyQuantumEntangleManager      // 幸運量子糾纏魚系統管理器（DAY-251）

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
	lastDailyTournamentAt time.Time // 每日賽排名廣播計時（每 30 秒一次，DAY-093）
	lastMultiFormatAt     time.Time // 多格式賽排名廣播計時（每 30 秒一次，DAY-111）
	lastGuildWarAt     time.Time  // 公會戰廣播計時（每 60 秒一次，DAY-076）
	lastDailyBossAt    time.Time  // 每日 BOSS 廣播計時（每 30 秒一次，DAY-077）
	lastEventAt        time.Time  // 限時活動廣播計時（每 30 秒一次，DAY-079）
	lastWeatherAt      time.Time  // 天氣廣播計時（每 30 秒一次，DAY-087）
	lastAutoSaveAt     time.Time  // 玩家資料定期自動儲存計時（每 60 秒，DAY-099）
	lastRaidTickAt     time.Time  // Co-op Raid 狀態廣播計時（每 3 秒，DAY-115）
	lastGoldenTimeTickAt time.Time // 黃金時間 tick 計時（每秒，DAY-125）
	lastRareCatchTickAt  time.Time // 稀有連擊過期檢查計時（每 5 秒，DAY-126）
	lastMegaCatchTickAt  time.Time // Mega Catch 隨機觸發計時（每 60 秒，DAY-140）

	// 補償機制
	lastHighRewardAt time.Time
	bonusSpecialBonus float64

	// 天氣湧現事件（DAY-127）
	weatherSurgeActive    bool      // 是否正在湧現
	weatherSurgeEndAt     time.Time // 湧現結束時間
	weatherSurgeRareBonus float64   // 稀有目標加成
	weatherSurgeGoldBonus float64   // 金幣魚加成
	weatherSurgeName      string    // 湧現名稱（用於廣播）

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
		dailyTournamentMgr: tournament.NewDaily(),
		multiFormatMgr:     tournament.NewMultiFormat(),
		Season:             season.New(),
		Friends:            friend.New(),
		FriendChallenge:    friendchallenge.New(),
		DM:                 dm.New(),
		Guild:              guild.New(),
		GuildWar:           guildwar.New(),
		DailyBoss:          dailyboss.New(),
		VIP:                vip.New(),
		Event:              event.New(30 * time.Minute),
		Referral:           referral.NewManager(),
		Wheel:              wheel.NewManager(),
		Challenge:          challenge.NewManager(),
		Weather:            weather.New(),
		Chain:              chain.NewDefault(),
		SpecialWeapon:      specialweapon.New(),
		MysteryBox:         mysterybox.New(),
		DailySpin:          dailyspin.NewManager(),
		Shop:               shop.New(),
		Announce:           announce.NewManager(),
		AntiCheat:          anticheat.New(),
		Festival:           festival.New(),
		HallOfFame:         halloffame.New(),
		ActivityFeed:       activityfeed.New(),
		Roulette:           roulette.NewManager(),
		RaidBoss:           raidboss.New(),
		Fragment:           fragment.New(),
		RespinMgr:          respin.New(),
		TreasureMap:        treasuremap.New(),
		FlashChallenge:     flashchallenge.New(),
		GoldenTime:         goldentime.New(),
		RareCatch:          rarecatch.New(),
		ImmortalBoss:       immortalboss.New(),
		AwakenBoss:         awakenboss.New(),
		WinStreak:          winstreak.New(),
		LightningEel:       lightningeel.New(),
		FeverMode:          fevermode.New(),
		UnluckyBonus:       unlucky.NewDefault(),
		SpeedRace:          speedrace.NewDefault(),
		Bounty:             bounty.NewDefault(),
		MultStorm:          multstorm.NewDefault(),
		DualRoulette:       dualroulette.NewDefault(),
		MegaCatch:          megacatch.NewDefault(),
		MegaOctopus:        megaoctopus.NewManager(),
		GiantPrizeFish:     newGiantPrizeFishManager(),
		ChainLongWheel:     chainlongwheel.New(),
		ThunderboltLobster: newThunderboltLobsterManager(),
		RainbowPhoenix:     newRainbowPhoenixManager(),
		GoldenTurtle:       newGoldenTurtleManager(),
		LuckyStarFish:      newLuckyStarFishManager(),
		GoldenShark:        newGoldenSharkManager(),
		CaptainFish:        newCaptainFishManager(),
		AbyssWhale:         newAbyssWhaleManager(),
		RouletteCrab:       roulettecrab.New(),
		CrystalDragon:      crystaldragon.New(),
		LionDance:          newLionDanceManager(),
		VortexFish:         newVortexFishManager(),
		FreezeBomb:         newFreezeBombManager(),
		IceFishing:         newIceFishingManager(),
		LuckyEgg:           newLuckyEggManager(),
		RainbowLucky:       newRainbowLuckyManager(),
		LuckyDice:          newLuckyDiceManager(),
		FireStorm:          newFireStormManager(),
		GoldenTreasure:     newGoldenTreasureManager(),
		Mermaid:            newMermaidManager(),
		LuckyClover:        newLuckyCloverManager(),
		RainbowShark:       newRainbowSharkManager(),
		VampireFish:        newVampireFishManager(),
		LightningAutoChain: newLightningAutoChainManager(),
		MeteorFish:         newMeteorFishManager(),
		PhoenixFish:        newPhoenixFishManager(),
		DragonTurtle:       newDragonTurtleManager(),
		CrocodileHunter:    newCrocodileHunterManager(),
		TimeBomb:           newTimeBombManager(),
		TripleLucky:        newTripleLuckyFishManager(),
		SchoolPanic:        newSchoolPanicManager(),
		RockSkeleton:       newRockSkeletonManager(),
		ChainLongKing:      newChainLongKingManager(),
		DrillLobster:       newDrillLobsterManager(),
		AnglerfishElectric: newAnglerfishManager(),
		MysticDragon:       newMysticDragonManager(),
		GhostFish:          newGhostFishManager(),
		ThunderboltLobsterV2: newTBLobsterV2Manager(),
		IcePhoenix:         newIcePhoenixManager(),
		SerialBombCrab:     newSerialBombCrabManager(),
		AbyssVortex:        newAbyssVortexManager(),
		HumpbackWhale:      newHumpbackWhaleManager(),
		FreeSpinFish:       newFreeSpinFishManager(),
		JackpotDragon:      newJackpotDragonManager(),
		CometFish:          newCometFishManager(),
		GoldenWaveFish:     newGoldenWaveFishManager(),
		DragonKing:         newDragonKingManager(),
		FortuneCoinFish:    newFortuneCoinFishManager(),
		LuckyHotZone:       newLuckyHotZoneManager(),
		LuckyTrident:       newLuckyTridentManager(),
		TimeFreeze:         newTimeFreezeManager(),
		RainbowPrism:       newRainbowPrismManager(),
		GoldenAccumulator:  newGoldenAccumulatorManager(),
		LuckyMirrorFish:    newLuckyMirrorFishManager(),
		CursedPoisonFish:   newCursedPoisonFishManager(),
		LuckyAuctionFish:   newLuckyAuctionFishManager(),
		LuckyEvolutionFish: newLuckyEvolutionFishManager(),
		LuckyInfectionFish: newLuckyInfectionFishManager(),
		LuckyRicochetFish:  newLuckyRicochetFishManager(),
		LuckyBlackHole:     newLuckyBlackHoleManager(),
		LuckyResonanceFish: newLuckyResonanceFishManager(),
		LuckyTeleportFish:  newLuckyTeleportFishManager(),
		LuckySplitFish:     newLuckySplitFishManager(),
		LuckyChargeFish:    newLuckyChargeFishManager(),
		LuckyChainBomb:     newLuckyChainBombManager(),
		LuckyMirrorTime:    newLuckyMirrorTimeManager(),
		LuckyQuantumFish:   newLuckyQuantumFishManager(),
		LuckyParasiteFish:  newLuckyParasiteFishManager(),
		LuckyStormFish:     newLuckyStormFishManager(),
		LuckyBoomerangFish: newLuckyBoomerangFishManager(),
		LuckyMagnetFish:    newLuckyMagnetFishManager(),
		LuckyEchoFish:      newLuckyEchoFishManager(),
		LuckyVortexFish:    newLuckyVortexFishManager(),
		LuckyTimeBombFish:  newLuckyTimeBombFishManager(),
		LuckyMirrorWorld:   newLuckyMirrorWorldManager(),
		LuckyFreezeWorld:   newLuckyFreezeWorldManager(),
		LuckyGravityFlip:   newLuckyGravityFlipManager(),
		LuckySynergyBurst:  newLuckySynergyBurstManager(),
		LuckyBetFish:       newLuckyBetFishManager(),
		LuckyChainReaction: newLuckyChainReactionManager(),
		LuckyCloneFish:     newLuckyCloneFishManager(),
		LuckyProphecyFish:  newLuckyProphecyFishManager(),
		LuckyFlagFish:      newLuckyFlagFishManager(),
		LuckyPhantomFish:   newLuckyPhantomFishManager(),
		LuckyCrystalBallFish: newLuckyCrystalBallFishManager(),
		LuckyTimeRewind:    newLuckyTimeRewindManager(),
		LuckyTornado:       newLuckyTornadoManager(),
		LuckyBlackHoleExplosion: newLuckyBlackHoleExplosionManager(),
		LuckyMirrorSplit:        newLuckyMirrorSplitManager(),
		LuckyQuantumEntangle:    newLuckyQuantumEntangleManager(),
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
	// 啟動好友挑戰結算計時器（DAY-102）
	g.startChallengeTicker()
	// 啟動巨型章魚轉盤超時計時器（DAY-144）
	g.startMegaOctopusTicker()
	go g.gameLoop()
}

// Stop 停止遊戲（graceful shutdown，儲存所有玩家資料，DAY-099）
func (g *Game) Stop() {
	g.saveAllPlayersOnShutdown()
	close(g.stopCh)
}

// AddPlayer 加入玩家（從 Store 恢復狀態，若無則建立新玩家）
func (g *Game) AddPlayer(playerID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.Players[playerID]; !exists {
		p := player.NewPlayer(playerID, g.initialCoins)

		// 從 Store 恢復玩家完整狀態（DAY-098）
		var savedState *store.PlayerState
		if g.store != nil {
			g.restoreFullPlayerState(p)
			// 恢復好友關係（DAY-101）
			g.restoreFriendState(p.ID)
			// 取得基礎 savedState 供 checkAndSendDailyBonus 使用
			if saved, err := g.store.LoadPlayer(playerID); err == nil && saved != nil {
				savedState = saved
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
				// 發放離線期間收到的禮物（DAY-101）
				g.deliverPendingGifts(pp)
				// 發送離線期間收到的私訊（DAY-103）
				g.deliverPendingDMs(pp)
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
				// 初始化隱藏挑戰（DAY-085）
				g.Challenge.InitPlayer(playerID)
				// 發送天氣狀態（DAY-087）
				g.sendWeatherUpdate(playerID, false)
				// 發送特殊武器狀態（DAY-089）
				g.sendSpecialWeaponUpdate(pp, true)
				// 發送神秘寶箱狀態（DAY-090）
				g.sendMysteryBoxUpdate(pp)
				// 發送每日轉盤狀態（DAY-092）
				g.handleGetDailySpin(pp.ID)
				// 發送週賽/每日賽狀態（DAY-093）
				g.sendTournamentUpdate(pp.ID)
				g.sendDailyTournamentUpdate(pp.ID)
				// 發送多格式賽狀態（DAY-111）
				g.sendMultiFormatUpdate(pp.ID)
				// 發送成就動態牆歷史（DAY-112）
				g.sendActivityFeedHistory(pp.ID)
				// 發送商店狀態（DAY-094）
				g.sendShopUpdate(pp)
				// 啟動統計 Session（DAY-096）
				if pp.Stats != nil {
					pp.Stats.StartSession()
				}
				// 發送個人統計（DAY-096）
				g.sendPlayerStats(pp)
				// 全服公告：玩家加入（DAY-097）
				g.announcePlayerJoin(pp.DisplayName)
				// 異常偵測：建立玩家記錄（DAY-105）
				g.AntiCheat.EnsureRecord(playerID, pp.DisplayName)
				// 發送登入進度（DAY-107）
				g.handleGetLoginProgress(pp)
				// 發送節日狀態（DAY-109）
				g.sendFestivalState(pp)
				// 發送碎片狀態（DAY-116）
				g.sendFragmentStatus(pp)
				// 發送寶藏地圖狀態（DAY-122）
				g.sendTreasureMapUpdate(pp)
				// 發送黃金時間狀態（DAY-125）
				g.handleGetGoldenTime(pp)
				// 發送龍怒蓄力大招狀態（DAY-128）
				g.sendWrathStatus(pp)
				// 發送不死 BOSS 狀態（DAY-129）
				g.sendImmortalBossStatus(pp)
				// 發送覺醒 BOSS 狀態（DAY-130）
				g.sendAwakenBossStatus(pp)
				// 發送閃電鰻冷卻狀態（DAY-132）
				g.sendLightningEelStatus(pp)
				// 發送狂熱模式狀態（DAY-133）
				g.sendFeverModeStatus(pp)
				// 發送失敗補償狀態（DAY-135）
				g.sendUnluckyBonusStatus(pp)
				// 發送競速獵殺狀態（DAY-136）
				g.sendSpeedRaceStatus(pp)
				// 發送懸賞列表（DAY-137）
				g.sendBountyStatus(pp)
				// 發送倍率風暴狀態（DAY-138）
				g.sendMultStormStatus(pp.ID)
				// 發送雙環輪盤狀態（DAY-139）
				g.sendDualRouletteStatus(pp.ID)
				// 發送 Mega Catch 狀態（DAY-140）
				g.sendMegaCatchStatus(pp.ID)
				// 發送巨型章魚轉盤狀態（DAY-144）
				g.sendMegaOctopusStatus(pp)
				// 發送千龍王輪盤狀態（DAY-148）
				g.sendChainLongWheelStatus(pp)
				// 發送彩虹鳳凰 Power Up 狀態（DAY-151）
				g.sendRainbowPhoenixStatus(pp)
				// 發送水晶龍收集進度狀態（DAY-153）
				g.sendCrystalDragonStatus(pp)
				// 發送輪盤螃蟹冷卻狀態（DAY-167）
				g.sendRouletteCrabStatus(pp)
				// 任務連續寬限期檢查（DAY-120）
				go g.checkMissionStreakMercy(pp)
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

	// 結束統計 Session（DAY-096）
	if p != nil && p.Stats != nil {
		p.Stats.RecordSessionScore(p.SessionScore)
		p.Stats.EndSession()
	}
	// 全服公告：玩家離開（DAY-097）
	if p != nil {
		g.announcePlayerLeave(p.DisplayName)
	}

	// 儲存玩家完整狀態到 Store（DAY-098）
	if g.store != nil && p != nil {
		// 結束統計 Session 後再儲存（確保 TotalPlayTime 正確）
		g.saveFullPlayerState(p)
		// 儲存好友關係（DAY-101）
		g.saveFriendState(p.ID)
	}

	log.Printf("[Game] Player %s left game %s", playerID, g.ID)

	// 通知好友下線（DAY-073）
	if p != nil {
		go g.notifyFriendsOffline(playerID, p.DisplayName)
		// 強制結算進行中的挑戰（DAY-102）
		if c := g.FriendChallenge.ForceFinish(playerID); c != nil {
			go g.settleChallengeResult(c)
		}
		// 更新公會在線狀態（DAY-074）
		g.Guild.SetOnlineStatus(playerID, false)
		// 清理挑戰 session（DAY-085）
		g.Challenge.RemovePlayer(playerID)
		// 清理特殊武器狀態（DAY-089）
		g.SpecialWeapon.RemovePlayer(playerID)
		// 清理神秘寶箱狀態（DAY-090）
		g.MysteryBox.RemovePlayer(playerID)
		// 清理異常偵測記錄（DAY-105）
		g.AntiCheat.RemoveRecord(playerID)
		// 清理節日進度（DAY-109）
		g.Festival.RemovePlayer(playerID)
		// 清理輪盤 session（DAY-113）
		g.Roulette.CancelSession(playerID)
		// 清理碎片狀態（DAY-116）
		g.Fragment.RemovePlayer(playerID)
		// 清理 Rapid Respin session（DAY-121）
		g.RespinMgr.RemovePlayer(playerID)
		// 清理寶藏地圖狀態（DAY-122）
		g.TreasureMap.RemovePlayer(playerID)
		// 清理稀有連擊 session（DAY-126）
		g.RareCatch.RemovePlayer(playerID)
		// 清理連勝記錄（DAY-131）
		if g.WinStreak != nil {
			g.WinStreak.RemovePlayer(playerID)
		}
		// 清理閃電鰻冷卻記錄（DAY-132）
		if g.LightningEel != nil {
			g.LightningEel.RemovePlayer(playerID)
		}
		// 清理狂熱模式記錄（DAY-133）
		if g.FeverMode != nil {
			g.FeverMode.RemovePlayer(playerID)
		}
		// 清理失敗補償記錄（DAY-135）
		if g.UnluckyBonus != nil {
			g.UnluckyBonus.RemovePlayer(playerID)
		}
		// 清理雙環輪盤 session（DAY-139）
		if g.DualRoulette != nil {
			g.DualRoulette.RemovePlayer(playerID)
		}
		// 清理巨型章魚轉盤 session（DAY-144）
		if g.MegaOctopus != nil {
			g.MegaOctopus.RemovePlayer(playerID)
		}
		// 清理千龍王輪盤 session（DAY-148）
		if g.ChainLongWheel != nil {
			g.ChainLongWheel.RemovePlayer(playerID)
		}
		// 清理雷霆龍蝦免費射擊 session（DAY-150）
		if g.ThunderboltLobster != nil {
			g.ThunderboltLobster.RemovePlayer(playerID)
		}
		// 清理彩虹鳳凰 Power Up session（DAY-151）
		if g.RainbowPhoenix != nil {
			g.RainbowPhoenix.RemovePlayer(playerID)
		}
		// FlashChallenge 不需要清理（進度保留到挑戰結束）
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
	case ws.MsgTriggerRaid:
		// Co-op Boss Raid 手動觸發（Prototype 展示用，DAY-115）
		go g.triggerRaid()
	case ws.MsgGetRaidStatus:
		// 查詢討伐狀態（DAY-115）
		g.handleGetRaidStatus(p)
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
	// 特殊武器系統（DAY-089）
	case ws.MsgBuySpecialWeapon:
		g.handleBuySpecialWeapon(p, msg)
	case ws.MsgUseSpecialWeapon:
		g.handleUseSpecialWeapon(p, msg)
	case ws.MsgGetSpecialWeapons:
		g.handleGetSpecialWeapons(p)
	// 神秘寶箱系統（DAY-090）
	case ws.MsgOpenMysteryBox:
		g.handleOpenMysteryBox(p, msg)
	case ws.MsgGetMysteryBoxes:
		g.handleGetMysteryBoxes(p)
	// 房間難度系統（DAY-091）
	case ws.MsgGetRoomList:
		g.handleGetRoomList(clientID)
	case ws.MsgSwitchRoom:
		var payload ws.SwitchRoomPayload
		if err := remarshal(msg.Payload, &payload); err == nil {
			g.handleSwitchRoom(clientID, payload)
		}
	// 每日簽到轉盤（DAY-092）
	case ws.MsgGetDailySpin:
		g.handleGetDailySpin(clientID)
	case ws.MsgDailySpin:
		g.handleDailySpin(clientID)
	// 週賽/每日賽查詢（DAY-093）
	case ws.MsgGetTournament:
		g.handleGetTournament(p)
	// 商店系統（DAY-094）
	case ws.MsgGetShop:
		g.handleGetShop(p)
	case ws.MsgBuyShopItem:
		g.handleBuyShopItem(p, msg)
	// 玩家統計系統（DAY-096）
	case ws.MsgGetPlayerStats:
		g.handleGetPlayerStats(p)
	// 好友禮物系統（DAY-101）
	case ws.MsgSendGift:
		g.handleSendGift(p, msg)
	case ws.MsgGetGiftStatus:
		g.handleGetGiftStatus(p)
	// 好友挑戰系統（DAY-102）
	case ws.MsgSendChallengeRequest:
		g.handleSendChallengeRequest(p, msg)
	case ws.MsgAcceptChallenge:
		g.handleAcceptChallenge(p, msg)
	case ws.MsgDeclineChallenge:
		g.handleDeclineChallenge(p, msg)
	// 私訊系統（DAY-103）
	case ws.MsgSendDM:
		g.handleSendDM(p, msg)
	// 玩家名片系統（DAY-106）
	case ws.MsgGetPlayerCard:
		g.handleGetPlayerCard(p, msg)
	// 登入里程碑系統（DAY-107）
	case ws.MsgGetLoginProgress:
		g.handleGetLoginProgress(p)
	// 賽季節日活動系統（DAY-109）
	case ws.MsgGetFestival:
		g.handleGetFestival(p)
	case ws.MsgClaimFestivalTask:
		var payload ws.ClaimFestivalTaskPayload
		if err := remarshal(msg.Payload, &payload); err == nil {
			g.handleClaimFestivalTask(p, payload.TaskID)
		}
	// 名人堂系統（DAY-110）
	case ws.MsgGetHallOfFame:
		g.handleGetHallOfFame(p)
	// 智慧推薦系統（DAY-110）
	case ws.MsgGetRecommendations:
		g.handleGetRecommendations(p)
	// 成就動態牆系統（DAY-112）
	case ws.MsgGetActivityFeed:
		g.handleGetActivityFeed(p)
	// Buy Bonus 系統（DAY-114）
	case ws.MsgBuyBonus:
		g.handleBuyBonus(p, msg)
	case ws.MsgGetBuyBonusStatus:
		g.handleGetBuyBonusStatus(p)
	// 新手引導系統（DAY-115）
	case ws.MsgTutorialAction, ws.MsgSkipTutorial:
		g.handleTutorialAction(p, msg)
	// 碎片收集大獎系統（DAY-116）
	case ws.MsgGetFragments:
		g.handleGetFragments(p)
	// 寶藏地圖系統（DAY-122）
	case ws.MsgGetTreasureMap:
		g.handleGetTreasureMap(p)
	// 閃電挑戰系統（DAY-123）
	case ws.MsgGetFlashChallenge:
		g.handleGetFlashChallenge(p)
	// 黃金時間系統（DAY-125）
	case ws.MsgGetGoldenTime:
		g.handleGetGoldenTime(p)
	// 龍怒蓄力大招系統（DAY-128）
	case ws.MsgUseWrath:
		g.handleUseWrath(clientID)
	// 全服目標懸賞系統（DAY-137）
	case ws.MsgPostBounty:
		g.handlePostBounty(p, msg)
	case ws.MsgGetBounties:
		g.sendBountyStatus(p)
	// 雙環輪盤系統（DAY-139）
	case ws.MsgDualRouletteStop:
		g.handleDualRouletteStop(clientID)
	// 巨型章魚轉盤系統（DAY-144）
	case ws.MsgMegaOctopusWheelStop:
		g.mu.RLock()
		p, ok := g.Players[clientID]
		g.mu.RUnlock()
		if ok {
			g.handleMegaOctopusStop(p)
		}
	// 千龍王強化輪盤系統（DAY-148）
	case ws.MsgChainLongWheelStop:
		g.handleChainLongWheelStop(clientID)
	// 黃金輪盤螃蟹系統（DAY-167）
	case ws.MsgRouletteCrabStop:
		g.handleRouletteCrabWheelStop(p)
	// 冰釣幸運輪盤系統（DAY-171）
	case ws.MsgIceFishingWheelStop:
		go g.handleIceFishingWheelStop(p)
	// 黃金寶藏魚系統（DAY-177）
	case ws.MsgGoldenTreasureOpen:
		var payload ws.GoldenTreasureOpenPayload
		if err := remarshal(msg.Payload, &payload); err == nil {
			go g.handleGoldenTreasureOpen(p, payload.ChestID)
		}
	// 長龍王雙環輪盤系統（DAY-194）
	case ws.MsgChainLongKingStop:
		var payload ws.ChainLongKingStopPayload
		if err := remarshal(msg.Payload, &payload); err == nil {
			g.handleChainLongKingStop(p, payload)
		}
	// 幸運三叉魚互動三轉盤系統（DAY-211）
	case ws.MsgLuckyTridentStop:
		g.handleLuckyTridentStop(p, msg)
	// 幸運拍賣魚競標出價（DAY-217）
	case ws.MsgLuckyAuctionBid:
		g.handleLuckyAuctionBid(p)
	// 幸運賭注魚選擇（DAY-240）
	case ws.MsgLuckyBetChoice:
		g.handleLuckyBetChoice(p, msg)
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
	// 玩家統計：記錄射擊（DAY-096）
	g.notifyStatsShot(p, betCost)
	// 多格式每日賽：記錄投注（DAY-111，投注競賽格式用）
	go g.multiFormatMgr.RecordShot(p.ID, p.DisplayName, betCost)
	// 龍怒蓄力大招：每次射擊累積怒氣（DAY-128）
	go g.notifyWrathShot(p)
	// 龍怒流星雨武器：每次射擊累積充能（DAY-154）
	go g.notifyDragonWrathShot(p)
	// 深海龍王：蓄力模式中每次射擊累積龍怒值（DAY-208）
	go g.notifyDragonKingShot()
	// 幸運共鳴魚：共鳴模式中每次射擊累積共鳴能量（DAY-222）
	go g.notifyResonanceShot(p)
	// 幸運充能魚：充能模式中每次射擊累積充能值（DAY-225）
	go g.notifyChargeShot(p)
	// 不死 BOSS：每次射擊嘗試命中（DAY-129）
	go g.tryImmortalBossHit(p)
	// 覺醒 BOSS：每次射擊嘗試命中（DAY-130）
	go g.tryAwakenBossHit(p)
	// 異常偵測：記錄攻擊（DAY-105）
	if alert := g.AntiCheat.RecordAttack(p.ID, betCost); alert != nil {
		log.Printf("[AntiCheat] Alert triggered for player %s: %s", p.ID, alert.Message)
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
		EventKillAdd:   g.getEventKillChanceAdd() + g.GetRainbowLuckyBoost() + g.getTimeFreezeHitBonus(), // 限時活動+彩虹幸運魚+時間凍結擊破率加成（DAY-079/DAY-173/DAY-212）
	}

	result := combat.ProcessAttack(req, t)

	// Progressive Jackpot 貢獻（DAY-048）：每次攻擊抽取 0.5% 進入 Jackpot 池
	// 套用房間難度 Jackpot 倍率（DAY-091）：高難度房間累積更快
	jackpotBetCost := int(float64(betCost) * g.getRoomJackpotMult(p))
	if jackpotBetCost < betCost {
		jackpotBetCost = betCost
	}
	// 套用節日 Jackpot 倍率（DAY-109）：節日期間 Jackpot 累積更快
	festivalJackpotMult := g.getFestivalJackpotMult()
	if festivalJackpotMult > 1.0 {
		jackpotBetCost = int(float64(jackpotBetCost) * festivalJackpotMult)
	}
	if jackpotWin := g.jackpotMgr.Contribute(jackpotBetCost, p.ID); jackpotWin != nil {
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
		// 失敗補償：記錄射擊（有擊破，reward > 0）（DAY-135）
		go g.notifyUnluckyShot(p, betCost, result.Reward)
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
	} else if result.IsHit && t != nil && isVampire(t.DefID) {
		// T116 吸血鬼：命中後倍率成長（DAY-152）
		go g.notifyVampireHit(p, t)
	}

	// 深淵巨鯨：命中時記錄傷害（DAY-164）
	if result.IsHit && t != nil && isAbyssWhale(t.DefID) {
		go g.notifyAbyssWhaleHit(p, t.InstanceID, betCost)
	}

	// 時間凍結魚：凍結期間命中目標時記錄（DAY-212）
	if result.IsHit && t != nil && !result.IsKill {
		g.notifyTimeFreezeTargetHit(t)
	}

	// 幸運進化魚：命中進化魚時累積命中次數（DAY-218）
	if result.IsHit && t != nil && !result.IsKill && isLuckyEvolutionFish(t.DefID) {
		go g.notifyLuckyEvolutionFishHit(t.InstanceID, p)
	}

	// 幸運反彈魚：命中目標後觸發反彈（DAY-220）
	if result.IsHit && t != nil && g.isRicochetActive(p.ID) {
		go g.doRicochetBounce(p, t.X, t.Y, t.InstanceID, 0)
	}

	// 幸運迴旋鏢魚：命中目標後觸發折返（DAY-231）
	// 計算射擊方向（從砲台位置到目標位置）
	if result.IsHit && t != nil && g.isBoomerangActive(p.ID) {
		// 預設射擊方向為向右（-1,0 表示從右往左折返）
		dirX := -1.0
		dirY := 0.0
		if t.X > 0 {
			// 從砲台（左側）到目標的方向
			dirX = 1.0
		}
		go g.doBoomerangBounce(p, t.X, t.Y, dirX, dirY, t.InstanceID, 0)
	}

	// 幸運分身魚：射擊時同時發射兩個偏移方向的分身子彈（DAY-242）
	if t != nil && g.isCloneModeActive(p.ID) {
		// 砲台在左側（x≈0），目標在右側
		go g.doCloneShots(p, 0, t.Y, t.X, t.Y, t.InstanceID)
	}

	// 幸運奪旗魚：命中旗幟目標時累積搶旗積分（DAY-244）
	if result.IsHit && t != nil && !result.IsKill && g.isFlagTarget(t.InstanceID) {
		g.recordFlagHit(p.ID, p.DisplayName)
	}

	// 龍龜不死 Boss：命中時直接給獎勵（DAY-186）
	if result.IsHit && t != nil && isDragonTurtle(t.DefID) {
		go g.notifyDragonTurtleHit(p, t.InstanceID)
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

	// 失敗補償：未擊破時記錄（spend=betCost, reward=0）（DAY-135）
	if !result.IsKill {
		go g.notifyUnluckyShot(p, betCost, 0)
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
	// 套用天氣獎勵倍率（DAY-087）
	weatherRewardMult := g.Weather.GetRewardMult()
	if weatherRewardMult > 1.0 {
		finalReward = int(float64(finalReward) * weatherRewardMult)
	}
	// 套用連擊倍率（DAY-083）
	streakMult := g.notifyStreakKill(p)
	if streakMult > 1.0 {
		finalReward = int(float64(finalReward) * streakMult)
	}
	// 套用房間難度獎勵倍率（DAY-091）
	roomRewardMult := g.getRoomRewardMult(p)
	if roomRewardMult > 1.0 {
		finalReward = int(float64(finalReward) * roomRewardMult)
	}
	// 套用節日獎勵倍率（DAY-109）
	festivalRewardMult := g.getFestivalRewardMult()
	if festivalRewardMult > 1.0 {
		finalReward = int(float64(finalReward) * festivalRewardMult)
	}
	// 套用黃金時間倍率（DAY-125）
	goldenTimeMult := g.GoldenTime.GetMultBoost()
	if goldenTimeMult > 1.0 {
		finalReward = int(float64(finalReward) * goldenTimeMult)
	}
	// 套用稀有連擊倍率（DAY-126）
	rareCatchMult := g.notifyRareCatchKill(p, t.DefID)
	if rareCatchMult > 1.0 {
		finalReward = int(float64(finalReward) * rareCatchMult)
	}
	// 套用狂熱模式倍率（DAY-133）
	feverMult := g.notifyFeverModeKill(p)
	if feverMult > 1.0 {
		finalReward = int(float64(finalReward) * feverMult)
	}
	// 套用競速獵殺倍率（DAY-136）
	speedRaceMult := g.notifySpeedRaceKill(p, t.InstanceID)
	if speedRaceMult > 1.0 {
		finalReward = int(float64(finalReward) * speedRaceMult)
	}
	// 套用倍率風暴加成（DAY-138）
	stormMult := g.getMultStormBoost()
	if stormMult > 1.0 {
		finalReward = int(float64(finalReward) * stormMult)
	}
	// 套用 Mega Catch 獎勵倍率（DAY-140）
	megaCatchMult := g.getMegaCatchRewardBoost()
	if megaCatchMult > 1.0 {
		finalReward = int(float64(finalReward) * megaCatchMult)
	}
	// 套用夢幻巨型獎勵魚倍率（DAY-147）
	giantPrizeMult := g.getGiantPrizeFishMult(p.ID)
	if giantPrizeMult > 1.0 {
		finalReward = int(float64(finalReward) * giantPrizeMult)
		// 記錄夢幻模式期間的擊破
		go g.recordGiantPrizeFishKill(p.ID, finalReward)
	}
	// 套用彩虹鳳凰 Power Up 倍率（DAY-151）
	rainbowPhoenixMult := g.getRainbowPhoenixMult(p.ID)
	if rainbowPhoenixMult > 1.0 {
		finalReward = int(float64(finalReward) * rainbowPhoenixMult)
		// 記錄 Power Up 期間的擊破
		go g.notifyRainbowPhoenixKill(p, finalReward)
	}
	// 套用幸運星魚倍率翻倍（DAY-160）
	luckyStarMult := g.getLuckyStarMult(p.ID)
	if luckyStarMult > 1.0 {
		finalReward = int(float64(finalReward) * luckyStarMult)
	}
	// 套用黃金鯊魚全服狂暴倍率（DAY-161）
	goldenSharkMult := g.getGoldenSharkMult()
	if goldenSharkMult > 1.0 {
		finalReward = int(float64(finalReward) * goldenSharkMult)
	}
	// 套用獅子舞爆發倍率（DAY-168）
	lionDanceMult := g.getLionDanceMult(p.ID, t.InstanceID)
	if lionDanceMult > 1.0 {
		finalReward = int(float64(finalReward) * lionDanceMult)
	}
	// 套用冰釣幸運輪盤倍率（DAY-171）
	iceFishingMult := g.getIceFishingMult(p.ID)
	if iceFishingMult > 1.0 {
		bonusReward := int(float64(finalReward) * (iceFishingMult - 1.0))
		finalReward = int(float64(finalReward) * iceFishingMult)
		go g.recordIceFishingKill(p.ID, bonusReward)
	}
	// 套用幸運彩蛋魚倍率加成（DAY-172）
	luckyEggMult := g.getLuckyEggMult(p.ID)
	if luckyEggMult > 1.0 {
		finalReward = int(float64(finalReward) * luckyEggMult)
	}
	// 套用黃金寶藏魚倍率加成（DAY-177）
	goldenTreasureMult := g.getGoldenTreasureMult(p.ID)
	if goldenTreasureMult > 1.0 {
		finalReward = int(float64(finalReward) * goldenTreasureMult)
	}
	// 套用美人魚幸運加成（DAY-178，+20% 加法）
	mermaidBoost := g.getMermaidLuckBoost()
	if mermaidBoost > 0.0 {
		finalReward = finalReward + int(float64(finalReward)*mermaidBoost)
	}
	// 套用幸運草魚加成（DAY-179，+50% 加法，全服共享）
	cloverBoost := g.getLuckyCloverBoost()
	if cloverBoost > 0.0 {
		finalReward = finalReward + int(float64(finalReward)*cloverBoost)
	}
	// 套用鳳凰魚重生加成（DAY-185，+30% 加法，全服共享，30 秒）
	phoenixBoost := g.getPhoenixRebirthBoost()
	if phoenixBoost > 0.0 {
		finalReward = finalReward + int(float64(finalReward)*phoenixBoost)
	}
	// 套用時間炸彈魚拆彈加成（DAY-189，+25% 加法，全服共享，15 秒）
	timeBombBoost := g.getTimeBombDefuseBoost()
	if timeBombBoost > 0.0 {
		finalReward = finalReward + int(float64(finalReward)*timeBombBoost)
	}
	// 套用三重幸運魚倍率加成（DAY-190，+50% 加法，個人，12 秒）
	tripleLuckyBonus := g.getTripleLuckyMultBonus(p.ID)
	if tripleLuckyBonus > 0.0 {
		finalReward = finalReward + int(float64(finalReward)*tripleLuckyBonus)
	}
	// 套用搖滾骷髏演唱會安可加成（DAY-192，+30% 加法，全服共享，10 秒）
	rockSkeletonEncoreBoost := g.getRockSkeletonEncoreBoost()
	if rockSkeletonEncoreBoost > 0.0 {
		finalReward = finalReward + int(float64(finalReward)*rockSkeletonEncoreBoost)
	}
	// 套用黃金波浪魚加成（DAY-207，×2.0 乘法，全服共享，8 秒）
	goldenWaveBoost := g.getGoldenWaveBoost()
	if goldenWaveBoost > 1.0 {
		finalReward = int(float64(finalReward) * goldenWaveBoost)
	}
	// 套用幸運熱區魚加成（DAY-210，×2.0 乘法，空間限定，8 秒）
	hotZoneMult := g.getLuckyHotZoneMultiplier(t.X, t.Y)
	if hotZoneMult > 1.0 {
		finalReward = int(float64(finalReward) * hotZoneMult)
	}
	// 套用幸運三叉魚倍率加成（DAY-211，×1.5-5.0 乘法，個人，15 秒）
	tridentMult := g.getLuckyTridentMultBoost(p.ID)
	if tridentMult > 1.0 {
		finalReward = int(float64(finalReward) * tridentMult)
	}
	// 套用彩虹稜鏡魚倍率加成（DAY-213，×1.5-5.0 乘法，目標染色，10 秒）
	prismMult := g.getRainbowPrismMultiplier(t.InstanceID)
	if prismMult > 1.0 {
		finalReward = int(float64(finalReward) * prismMult)
		// 移除已擊破的染色標記
		g.removeRainbowPrismColor(t.InstanceID)
	}
	// 套用黃金累積魚爆發倍率加成（DAY-214，×2.0 乘法，全服共享，8 秒）
	goldenAccumBoost := g.getGoldenAccumulatorBoost()
	if goldenAccumBoost > 1.0 {
		finalReward = int(float64(finalReward) * goldenAccumBoost)
	}
	// 套用幸運鏡像魚分身倍率加成（DAY-215，×1.5 乘法，擊破鏡像分身時）
	mirrorMult := g.getLuckyMirrorMultiplier(t.InstanceID)
	if mirrorMult > 1.0 {
		finalReward = int(float64(finalReward) * mirrorMult)
	}
	// 套用詛咒毒魚倍率加成（DAY-216，×2.5 乘法，擊破詛咒目標時）
	cursedMult := g.getCursedPoisonKillMult(t.InstanceID)
	if cursedMult > 1.0 {
		finalReward = int(float64(finalReward) * cursedMult)
		g.removeCursedEntry(t.InstanceID)
	}
	// 套用詛咒懲罰倍率（DAY-216，×0.5 乘法，詛咒目標逃跑後持續 5 秒）
	penaltyMult := g.getCursedPoisonPenaltyMult(p.ID)
	if penaltyMult < 1.0 {
		finalReward = int(float64(finalReward) * penaltyMult)
	}
	// 套用幸運拍賣魚大獎控制權倍率（DAY-217，×0.85 乘法，控制權持有者）
	auctionControlMult := g.getLuckyAuctionControlMult(p.ID)
	if auctionControlMult < 1.0 {
		finalReward = int(float64(finalReward) * auctionControlMult)
	}
	// 套用幸運進化魚終極爆發倍率（DAY-218，×4.0 乘法，全服共享，6 秒）
	evolutionBurstBoost := g.getLuckyEvolutionBurstBoost()
	if evolutionBurstBoost > 1.0 {
		finalReward = int(float64(finalReward) * evolutionBurstBoost)
	}
	// 套用幸運進化魚本身倍率加成（DAY-218，×1.5/2.5/4.0 乘法，依進化階段）
	evolutionKillMult := g.getLuckyEvolutionKillMult(t.InstanceID)
	if evolutionKillMult > 1.0 {
		finalReward = int(float64(finalReward) * evolutionKillMult)
	}
	// 套用幸運連鎖感染魚感染目標倍率加成（DAY-219，×2.0 乘法，感染目標被擊破時）
	infectionMult := g.getLuckyInfectionKillMult(t.InstanceID)
	if infectionMult > 1.0 {
		finalReward = int(float64(finalReward) * infectionMult)
	}
	// 套用幸運黑洞魚空間倍率加成（DAY-221，×2.0 乘法，黑洞範圍內目標被擊破時）
	blackHoleMult := g.getLuckyBlackHoleMultiplier(t.X, t.Y)
	if blackHoleMult > 1.0 {
		finalReward = int(float64(finalReward) * blackHoleMult)
	}
	// 套用幸運共鳴魚倍率加成（DAY-222，×1.3-1.8 乘法，全服共享，共鳴爆發期間）
	resonanceBoost := g.getLuckyResonanceBoost()
	if resonanceBoost > 1.0 {
		finalReward = int(float64(finalReward) * resonanceBoost)
	}
	// 套用幸運傳送魚傳送混亂倍率加成（DAY-223，×2.5 乘法，個人，傳送後 3 秒內）
	teleportBonus := g.getLuckyTeleportBonus(p)
	if teleportBonus > 1.0 {
		finalReward = int(float64(finalReward) * teleportBonus)
	}
	// 套用幸運分裂魚碎片倍率加成（DAY-224，×1.8 乘法，個人，擊破分裂碎片時）
	splitFragMult := g.getLuckySplitFragmentMult(t.InstanceID)
	if splitFragMult > 1.0 {
		finalReward = int(float64(finalReward) * splitFragMult)
		g.removeLuckySplitFragment(t.InstanceID)
	}
	// 套用幸運充能魚充能爆發倍率（DAY-225，×5.0 乘法，個人，一次性）
	chargeInstanceID := g.getChargeInstanceID(p.ID)
	chargeBurst := g.getLuckyChargeBurst(p)
	if chargeBurst > 1.0 {
		finalReward = int(float64(finalReward) * chargeBurst)
		go g.notifyLuckyChargeBurstUsed(p, chargeInstanceID, finalReward)
	}
	// 套用幸運鏈鎖爆炸魚引爆標記倍率（DAY-226，×1.5 乘法，個人，引爆標記目標）
	chainBombMult := g.getLuckyChainBombMult(t.InstanceID)
	if chainBombMult > 1.0 {
		finalReward = int(float64(finalReward) * chainBombMult)
		go g.notifyChainBombKill(p, t)
	}
	// 套用幸運鏡像時空魚時間倒流倍率（DAY-227，×2.0 乘法，全服，時間倒流期間）
	mirrorTimeBoost := g.getLuckyMirrorTimeBoost()
	if mirrorTimeBoost > 1.0 {
		finalReward = int(float64(finalReward) * mirrorTimeBoost)
	}
	// 套用幸運量子魚坍縮倍率（DAY-228，×3.0 或 ×0.8 乘法，個人，量子態目標）
	if g.isQuantumTarget(t.InstanceID) {
		quantumMult, isQuantum := g.getLuckyQuantumCollapseMult(t.InstanceID)
		if isQuantum {
			finalReward = int(float64(finalReward) * quantumMult)
			go g.notifyQuantumKill(p, t.InstanceID, quantumMult, finalReward)
			g.removeQuantumEntry(t.InstanceID)
		}
	}
	// 套用幸運寄生魚擊破倍率（DAY-229，×2.2 乘法，個人，寄生目標）
	if g.isParasiteTarget(t.InstanceID) {
		parasiteMult := g.getLuckyParasiteKillMult(t.InstanceID)
		if parasiteMult > 1.0 {
			finalReward = int(float64(finalReward) * parasiteMult)
			go g.notifyParasiteKill(p, t.InstanceID, finalReward)
		}
	}
	// 套用幸運風暴魚擊破倍率（DAY-230，×2.5 乘法，空間限定，風暴範圍內）
	luckyStormMult := g.getLuckyStormMultiplier(t.InstanceID)
	if luckyStormMult > 1.0 {
		finalReward = int(float64(finalReward) * luckyStormMult)
	}
	// 套用幸運磁力魚磁力場倍率加成（DAY-232，×1.8 乘法，全服，磁力場期間）
	magnetBoost := g.getLuckyMagnetBoost()
	if magnetBoost > 1.0 {
		finalReward = int(float64(finalReward) * magnetBoost)
	}
	// 套用幸運回聲魚分身倍率加成（DAY-233，×1.5/2.0/2.5 乘法，個人，回聲分身擊破）
	echoMult, isEchoKill, _ := g.getLuckyEchoKillMult(t.InstanceID)
	if isEchoKill && echoMult > 1.0 {
		finalReward = int(float64(finalReward) * echoMult)
	}
	// 套用幸運漩渦魚漩渦模式倍率加成（DAY-234，×2.2 乘法，全服，漩渦模式期間）
	vortexBoost := g.getLuckyVortexBoost()
	if vortexBoost > 1.0 {
		finalReward = int(float64(finalReward) * vortexBoost)
	}
	// 套用幸運時間炸彈魚提前引爆倍率加成（DAY-235，×2.0 乘法，個人，炸彈目標）
	timeBombMult := g.getLuckyTimeBombEarlyMult(t.InstanceID)
	if timeBombMult > 1.0 {
		finalReward = int(float64(finalReward) * timeBombMult)
	}
	// 套用幸運鏡面世界魚鏡面世界倍率加成（DAY-236，×2.3 乘法，全服，鏡面世界期間）
	mirrorWorldBoost := g.getLuckyMirrorWorldBoost()
	if mirrorWorldBoost > 1.0 {
		finalReward = int(float64(finalReward) * mirrorWorldBoost)
	}
	// 套用幸運冰凍世界魚冰凍世界倍率加成（DAY-237，×2.0 乘法，全服，冰凍世界期間）
	freezeWorldBoost := g.getLuckyFreezeWorldBoost()
	if freezeWorldBoost > 1.0 {
		finalReward = int(float64(finalReward) * freezeWorldBoost)
	}
	// 套用幸運重力反轉魚重力反轉倍率加成（DAY-238，×2.1 乘法，全服，重力反轉期間）
	gravityFlipBoost := g.getLuckyGravityFlipBoost()
	if gravityFlipBoost > 1.0 {
		finalReward = int(float64(finalReward) * gravityFlipBoost)
	}
	// 套用幸運共鳴爆發魚共鳴爆發額外倍率（DAY-239，×1.3-1.8 乘法，個人，共鳴期間）
	synergyMult := g.getLuckySynergyBurstMult(p.ID)
	if synergyMult > 1.0 {
		finalReward = int(float64(finalReward) * synergyMult)
	}
	// 套用幸運賭注魚賭注倍率（DAY-240，×0.0-8.0 乘法，個人，一次性消耗）
	// 失敗時 mult=0.0（歸零），成功時 mult=2.0/4.0/8.0
	betMult := g.getLuckyBetFishMult(p.ID)
	if betMult != 1.0 {
		finalReward = int(float64(finalReward) * betMult)
		if finalReward < 0 {
			finalReward = 0
		}
	}
	// 套用幸運連鎖反應魚起點倍率（DAY-241，×1.4 乘法，個人，擊破連鎖起點時）
	chainReactionStarterMult := g.getLuckyChainReactionStarterMult(t.InstanceID)
	if chainReactionStarterMult > 1.0 {
		finalReward = int(float64(finalReward) * chainReactionStarterMult)
	}
	// 套用幸運預言魚倍率（DAY-243，×3.5 乘法，個人，擊破預言目標時）
	prophecyMult := g.getLuckyProphecyMult(p.ID, t.InstanceID)
	if prophecyMult > 1.0 {
		finalReward = int(float64(finalReward) * prophecyMult)
	}
	// 套用幸運奪旗魚倍率（DAY-244，×4.0/2.0/1.5 乘法，個人，擊破旗幟目標時）
	flagMult := g.getFlagWinnerMult(p.ID, t.InstanceID)
	if flagMult > 1.0 {
		finalReward = int(float64(finalReward) * flagMult)
	}
	// 套用幸運龍捲風魚龍捲風期間倍率加成（DAY-248，×2.2 乘法，全服，龍捲風期間）
	tornadoBoost := g.LuckyTornado.getLuckyTornadoBoost()
	if tornadoBoost > 1.0 {
		finalReward = int(float64(finalReward) * tornadoBoost)
	}
	// 套用彩虹鯊魚爆發倍率（DAY-180，乘法，全服共享，每個目標倍率不同）
	rainbowSharkMult := g.getRainbowSharkMult(t.InstanceID)
	if rainbowSharkMult > 1.0 {
		finalReward = int(float64(finalReward) * rainbowSharkMult)
		// 移除已擊破的彩虹標記
		g.removeRainbowSharkMark(t.InstanceID)
	}
	// 套用吸血鬼魚累積倍率（DAY-182，乘法，個人，越打越高）
	// 注意：先記錄擊破（累積倍率），再套用新倍率
	vampireMult := g.recordVampireKill(p)
	if vampireMult > 1.0 {
		finalReward = int(float64(finalReward) * vampireMult)
	}
	// 船長魚競速：記錄擊破（DAY-163）
	if g.IsCaptainRaceActive() {
		go g.recordCaptainRaceKill(p, finalReward)
	}
	// 雙環輪盤：擊破高倍率目標後嘗試觸發（DAY-139）
	go g.tryDualRoulette(p, float64(t.Multiplier), finalReward)
	// 懸賞領取：擊破懸賞目標獲得額外金幣（DAY-137）
	go g.notifyBountyKill(p, t.InstanceID)
	// 龍怒蓄力大招：擊破目標累積怒氣（DAY-128）
	go g.notifyWrathKill(p, t.Multiplier)
	// 連勝獎勵：記錄擊破（DAY-131）
	go g.notifyWinStreakKill(p)
	rewardUnlocks := p.AddReward(finalReward)
	killUnlocks := p.AddKill()

	// 異常偵測：記錄獎勵和金幣（DAY-105）
	go func() {
		if alert := g.AntiCheat.RecordReward(p.ID, finalReward); alert != nil {
			log.Printf("[AntiCheat] RTP Alert for player %s: %s", p.ID, alert.Message)
		}
		snap := p.Snapshot()
		if alert := g.AntiCheat.RecordCoins(p.ID, int64(snap.Coins)); alert != nil {
			log.Printf("[AntiCheat] Coin Spike Alert for player %s: %s", p.ID, alert.Message)
		}
	}()

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
	// 每日賽積分：擊破目標（DAY-093）
	g.notifyDailyTournamentKill(p, result.Multiplier)
	// 多格式每日賽：記錄擊破（DAY-111）
	go g.multiFormatMgr.RecordKill(p.ID, p.DisplayName, result.Multiplier, finalReward, data.GetBetDef(p.BetLevel).BetCost)
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
	// 幸運轉盤：擊破特殊目標後觸發（DAY-084）
	g.notifyWheelKill(p, t.DefID, finalReward)
	// 雙層倍率輪盤：擊破 BOSS/特殊目標後觸發（DAY-113）
	go g.notifyRouletteKill(p, t.DefID, finalReward)
	// Co-op Boss Raid：記錄傷害貢獻（DAY-115）
	go g.notifyRaidKill(p, finalReward)
	// 隱藏挑戰：記錄擊破事件（DAY-085）
	g.notifyChallengeKill(p, t.DefID, result.Multiplier, finalReward)
	// 連鎖爆炸：擊破後嘗試觸發連鎖（DAY-088）
	go g.notifyChainKill(p, t.InstanceID, t.X, t.Y, result.Multiplier, t.DefID)
	// 神秘寶箱：擊破後嘗試掉落（DAY-090）
	isBoss := t.DefID == "B001"
	go g.notifyMysteryBoxKill(p, t.X, t.Y, isBoss)
	// 碎片收集：擊破後嘗試掉落（DAY-116）
	go g.notifyFragmentKill(p, t.DefID, t.X, t.Y, isBoss)
	// 玩家統計：記錄擊破（DAY-096）
	g.notifyStatsKill(p, result.Multiplier, finalReward)
	// 全服公告：大獎（DAY-097，≥20x 才公告）
	if result.Multiplier >= 20 {
		g.announceBigWin(p.DisplayName, result.Multiplier, finalReward)
	}
	// 動態牆：超大獎（DAY-112，≥50x 才廣播）
	if result.Multiplier >= 50 {
		go g.notifyFeedMegaWin(p, result.Multiplier, finalReward)
	}
	// 好友挑戰：更新分數（DAY-102）
	g.notifyChallengeKillScore(p, finalReward)
	// 節日任務：記錄擊破（DAY-109）
	go g.notifyFestivalKill(p, t.DefID)
	// 名人堂：倍率和金幣記錄（DAY-110）
	go g.notifyHallOfFameKill(p, float64(t.Multiplier), finalReward)
	// 幸運捕獲：天氣加成期間觸發（DAY-119）
	if g.Weather.GetRewardMult() > 1.0 {
		go g.tryLuckyCatch(p, "weather")
	}
	// 幸運捕獲：節日期間觸發（DAY-119）
	if g.getFestivalRewardMult() > 1.0 {
		go g.tryLuckyCatch(p, "festival")
	}
	// Rapid Respin：擊破後嘗試觸發（DAY-121）
	go g.notifyRespinKill(p, finalReward)
	// 寶藏地圖：記錄擊破（DAY-122）
	go g.notifyTreasureMapKill(p, t.DefID)
	// 閃電挑戰：記錄擊破（DAY-123）
	streak := 0
	if p.Streak != nil {
		streak = p.Streak.GetSnapshot().Current
	}
	go g.notifyFlashChallengeKill(p, t.DefID, float64(t.Multiplier), streak)
	// 閃電鰻連鎖攻擊：擊破 T103 時觸發（DAY-132）
	if isLightningEel(t.DefID) {
		go g.tryLightningEelChain(p, t.InstanceID, t.X, t.Y)
	}
	// 鑽頭龍蝦連帶效果：擊破 T106 時觸發（DAY-142）
	if isDrillLobster(t.DefID) {
		go g.tryDrillLobsterChain(p, t.InstanceID, t.X, t.Y)
	}
	// 炸彈蟹連環爆炸：擊破 T107 時觸發（DAY-143）
	if isBombCrab(t.DefID) {
		go g.tryBombCrabChain(p, t.InstanceID, t.X, t.Y)
	}
	// 巨型章魚轉盤：擊破 T108 時觸發（DAY-144）
	if isMegaOctopus(t.DefID) {
		go g.tryMegaOctopusWheel(p, t.InstanceID)
	}
	// 巨型鮟鱇魚電擊寶箱：擊破 T109 時觸發（DAY-145）
	if isAnglerfish(t.DefID) {
		go g.tryAnglerfishShock(p, t.InstanceID, t.X, t.Y)
	}
	// 巨型鹹水鱷魚獵魚：擊破 T110 時觸發（DAY-146）
	if isCrocodile(t.DefID) {
		go g.tryCrocodileHunt(p, t.InstanceID, t.X, t.Y)
	}
	// 夢幻巨型獎勵魚：擊破 T111 時觸發夢幻獎勵模式（DAY-147）
	if isGiantPrizeFish(t.DefID) {
		go g.tryGiantPrizeFish(p, t.InstanceID, t.X, t.Y)
	}
	// 千龍王強化輪盤：擊破 T112 時觸發（DAY-148）
	if isChainLongKing(t.DefID) {
		go g.tryChainLongWheel(p, t.InstanceID, finalReward, t.Multiplier)
	}
	// 黃金水母全場電擊：擊破 T113 時觸發（DAY-149）
	if isGoldenJellyfish(t.DefID) {
		go g.tryGoldenJellyfishShock(p, t.InstanceID, t.X, t.Y)
	}
	// 雷霆龍蝦免費射擊：擊破 T114 時觸發（DAY-150）
	if isThunderboltLobster(t.DefID) {
		go g.tryThunderboltLobster(p, t.InstanceID, t.X, t.Y)
	}
	// 彩虹鳳凰 Power Up：擊破 T115 時觸發（DAY-151）
	if isRainbowPhoenix(t.DefID) {
		go g.tryRainbowPhoenix(p, t.InstanceID, t.X, t.Y)
	}
	// 吸血鬼成長倍率：擊破 T116 時廣播最終結果（DAY-152）
	if isVampire(t.DefID) {
		go g.notifyVampireKill(p, t, finalReward)
	}
	// 水晶龍收集大獎：擊破 T117 時掉落水晶（DAY-153）
	if isCrystalDragon(t.DefID) {
		go g.notifyCrystalDragonKill(p, t)
	}
	// 皇家閃電鰻持續連鎖電擊：擊破 T118 時觸發（DAY-156）
	if isRoyalChainLightning(t.DefID) {
		go g.notifyRoyalChainLightningKill(p, t.InstanceID, t.X, t.Y)
	}
	// 黃金海龜時間停止：擊破 T119 時觸發（DAY-159）
	if isGoldenTurtle(t.DefID) {
		go g.tryGoldenTurtleTimeStop(p, t.InstanceID, t.X, t.Y)
	}
	// 幸運星魚全場倍率翻倍：擊破 T120 時觸發（DAY-160）
	if isLuckyStarFish(t.DefID) {
		go g.tryLuckyStarFish(p, t.InstanceID, t.X, t.Y)
	}
	// 黃金鯊魚全服狂暴模式：擊破 T121 時觸發（DAY-161）
	if isGoldenShark(t.DefID) {
		go g.tryGoldenSharkBerserk(p, t.InstanceID, t.X, t.Y)
	}
	// 金幣魚王即時獎勵：擊破 T122 時觸發（DAY-162）
	if isMoneyFish(t.DefID) {
		go g.notifyMoneyFishKill(p, t.InstanceID, t.X, t.Y)
	}
	// 船長魚全服競速：擊破 T123 時觸發（DAY-163）
	if isCaptainFish(t.DefID) {
		go g.tryCaptainFishRace(p, t.InstanceID, t.X, t.Y)
	}
	// 深淵巨鯨全服 Boss 挑戰：擊破 T124 時觸發（DAY-164）
	if isAbyssWhale(t.DefID) {
		go g.notifyAbyssWhaleKill(p, t.InstanceID, t.X, t.Y)
	}
	// 黃金輪盤螃蟹：擊破 T125 時觸發（DAY-167）
	if isRouletteCrab(t.DefID) {
		go g.tryRouletteCrabWheel(p, t.Multiplier, finalReward)
	}
	// 獅子舞大獎爆發：擊破 T126 時觸發（DAY-168）
	if isLionDance(t.DefID) {
		go g.tryLionDanceBurst(p, t.InstanceID, t.X, t.Y)
	}
	// 漩渦魚群吸引：擊破 T127 時觸發（DAY-169）
	if isVortexFish(t.DefID) {
		go g.tryVortexFishSuck(p, t.InstanceID, t.X, t.Y)
	}
	// 冰凍炸彈魚：擊破 T128 時觸發（DAY-170）
	if isFreezeBomb(t.DefID) {
		go g.tryFreezeBomb(p, t.InstanceID, t.X, t.Y)
	}
	// 冰釣幸運輪盤：擊破 T129 時觸發（DAY-171）
	if isIceFish(t.DefID) {
		go g.tryIceFishingWheel(p, t.InstanceID, t.X, t.Y)
	}
	// 幸運彩蛋魚：擊破 T130 時觸發（DAY-172）
	if isLuckyEggFish(t.DefID) {
		go g.tryLuckyEggFish(p, t.InstanceID, t.X, t.Y)
	}
	// 彩虹幸運魚：擊破 T131 時觸發（DAY-173）
	if isRainbowLuckyFish(t.DefID) {
		go g.tryRainbowLuckyFish(p.DisplayName, t.X, t.Y)
	}
	// 海葵觸手攻擊：擊破 T132 時觸發（DAY-174）
	if isSeaAnemone(t.DefID) {
		go g.trySeaAnemone(p, t.InstanceID, t.X, t.Y)
	}
	// 幸運骰子魚：擊破 T133 時觸發（DAY-175）
	if isLuckyDiceFish(t.DefID) {
		go g.tryLuckyDiceFish(p, t.X, t.Y)
	}
	// 火焰風暴魚：擊破 T134 時觸發（DAY-176）
	if isFireStormFish(t.DefID) {
		go g.tryFireStormFish(p)
	}
	// 黃金寶藏魚：擊破 T135 時觸發（DAY-177）
	if isGoldenTreasureFish(t.DefID) {
		go g.tryGoldenTreasureFish(p)
	}
	// 美人魚：擊破 T136 時觸發（DAY-178）
	if isMermaid(t.DefID) {
		go g.tryMermaidHealing(p)
	}
	// 幸運草魚：擊破 T137 時觸發（DAY-179）
	if isLuckyCloverFish(t.DefID) {
		go g.tryLuckyCloverFish(p)
	}
	// 彩虹鯊魚：擊破 T138 時觸發（DAY-180）
	if isRainbowShark(t.DefID) {
		go g.tryRainbowSharkBurst(p, t.InstanceID)
	}
	// 雷霆鯊魚：擊破 T139 時觸發（DAY-181）
	if isThunderShark(t.DefID) {
		go g.tryThunderSharkChain(p, t.InstanceID, t.X, t.Y)
	}
	// 吸血鬼魚：擊破 T140 時觸發（DAY-182）
	if isVampireFish(t.DefID) {
		go g.tryVampireFish(p)
	}
	// 閃電魚：擊破 T141 時觸發（DAY-183）
	if isLightningAutoFish(t.DefID) {
		go g.tryLightningAutoChain(p, t.InstanceID, t.X, t.Y)
	}
	// 隕石魚：擊破 T142 時觸發（DAY-184）
	if isMeteorFish(t.DefID) {
		go g.tryMeteorFishShower(p, t.InstanceID, t.X, t.Y)
	}
	// 鳳凰魚：擊破 T143 時觸發（DAY-185）
	if isPhoenixFish(t.DefID) {
		go g.tryPhoenixFishRebirth(p, t.InstanceID, t.X, t.Y)
	}
	// 連鎖爆炸魚：擊破 T145 時觸發（DAY-187）
	if isChainBomb(t.DefID) {
		go g.tryChainBombExplosion(p, t.InstanceID, t.X, t.Y)
	}
	// 巨型鱷魚：擊破 T146 時觸發獎池結算（DAY-188）
	if isCrocodileHunter(t.DefID) {
		go g.notifyCrocodileHunterKill(p, t.InstanceID, t.Multiplier)
	}
	// 時間炸彈魚：擊破 T147 時觸發拆彈成功（DAY-189）
	if isTimeBombFish(t.DefID) {
		go g.notifyTimeBombDefuse(p, t.InstanceID, t.Multiplier)
	}
	// 三重幸運魚：擊破 T148 時觸發三重效果（DAY-190）
	if isTripleLuckyFish(t.DefID) {
		go g.tryTripleLuckyFish(p, t.InstanceID)
	}
	// 魚群領袖：擊破 T149 時觸發魚群驚嚇（DAY-191）
	if isSchoolLeader(t.DefID) {
		go g.trySchoolPanic(p, t.InstanceID)
	}
	// 搖滾骷髏魚：擊破 T150 時觸發演唱會（DAY-192）
	if isRockSkeleton(t.DefID) {
		go g.tryRockSkeletonConcert(p, t.InstanceID, t.X, t.Y)
	}
	// 電流水母：擊破 T151 時觸發電流網路（DAY-193）
	if isElectricJellyfish(t.DefID) {
		go g.tryElectricJellyfishNetwork(p, t.InstanceID, t.X, t.Y)
	}
	// 長龍王：擊破 T152 時觸發雙環輪盤（DAY-194）
	if isChainLongKingFish(t.DefID) {
		go g.tryChainLongKingRoulette(p, t.InstanceID)
	}
	// 鑽頭龍蝦：擊破 T153 時觸發穿透爆炸（DAY-195）
	if isDrillBitLobster(t.DefID) {
		go g.tryDrillLobsterPenetrate(p, t.InstanceID, t.X, t.Y)
	}
	// 巨型鮟鱇魚：擊破 T154 時觸發獎池結算（DAY-196）
	if isAnglerfishElectric(t.DefID) {
		go g.notifyAnglerfishKill(p, t.InstanceID, t.Multiplier)
	}
	// 神秘龍魚：擊破 T155 時觸發八波龍息攻擊（DAY-197）
	if isMysticDragon(t.DefID) {
		go g.tryMysticDragonWaves(p, t.InstanceID)
	}
	// 幽靈魚真身：擊破 T156 時觸發幻影爆炸（DAY-198）
	if isGhostFish(t.DefID) {
		go g.notifyGhostFishRealKill(p, t.InstanceID, t.Multiplier)
	}
	// 幽靈魚幻影分身：擊破 T156C 時給安慰獎（DAY-198）
	if isGhostFishClone(t.DefID) {
		go g.notifyGhostFishCloneKill(p, t.InstanceID)
	}
	// 雷霆龍蝦：擊破 T157 時觸發免費射擊模式（DAY-199）
	if isThunderboltLobsterV2(t.DefID) {
		go g.tryThunderboltLobsterFreePlay(p, t.Multiplier)
	}
	// 冰鳳凰：擊破 T158 時觸發覺醒（DAY-200）
	if isIcePhoenix(t.DefID) {
		go g.tryIcePhoenixAwaken(p, t.Multiplier)
	}
	// 連環炸彈蟹：擊破 T159 時觸發連環爆炸（DAY-201）
	if isSerialBombCrab(t.DefID) {
		go g.trySerialBombCrabExplosion(p, t.X, t.Y)
	}
	// 深淵漩渦魚：擊破 T160 時觸發深淵漩渦（DAY-202）
	if isAbyssVortexFish(t.DefID) {
		go g.tryAbyssVortexPull(p, t.X, t.Y)
	}
	// 座頭鯨：擊破 T161 時觸發鯨歌覺醒（DAY-203）
	if isHumpbackWhale(t.DefID) {
		go g.tryHumpbackWhaleAwaken(p, t.Multiplier)
	}
	// 自由旋轉魚：擊破 T162 時觸發個人免費射擊模式（DAY-204）
	if isFreeSpinFish(t.DefID) {
		go g.tryFreeSpinFishMode(p, t.Multiplier)
	}
	// 獎池龍：擊破 T163 時觸發獎池抽獎（DAY-205）
	if isJackpotDragon(t.DefID) {
		go g.tryJackpotDragonDraw(p)
	}
	// 彗星魚：擊破 T164 時提前引爆超新星（DAY-206）
	if isCometFish(t.DefID) {
		go g.notifyCometFishKill(p, t)
	}
	// 黃金波浪魚：擊破 T165 時觸發黃金波浪（DAY-207）
	if isGoldenWaveFish(t.DefID) {
		go g.tryGoldenWaveFish(t)
	}
	// 深海龍王：擊破 T166 時觸發全服合力蓄力模式（DAY-208）
	if isDragonKingFish(t.DefID) {
		go g.tryDragonKingCharge(t)
	}
	// 幸運金幣魚：擊破 T167 時觸發即時金幣爆發（DAY-209）
	if isFortuneCoinFish(t.DefID) {
		go g.tryFortuneCoinFish(p, t)
	}
	// 幸運熱區魚：擊破 T168 時在場上建立幸運熱區（DAY-210）
	if isLuckyHotZoneFish(t.DefID) {
		go g.tryLuckyHotZone(p, t)
	}
	// 幸運三叉魚：擊破 T169 時觸發三叉幸運儀式（DAY-211）
	if isLuckyTridentFish(t.DefID) {
		go g.tryLuckyTrident(p, t)
	}
	// 時間凍結魚：擊破 T170 時觸發時間凍結（DAY-212）
	if isTimeFreezeFish(t.DefID) {
		go g.tryTimeFreezeFish(p)
	}
	// 彩虹稜鏡魚：擊破 T171 時觸發稜鏡折射（DAY-213）
	if isRainbowPrismFish(t.DefID) {
		go g.tryRainbowPrismFish(p)
	}
	// 黃金累積魚：擊破 T172 時提前引爆（DAY-214）
	if isGoldenAccumulatorFish(t.DefID) {
		go g.notifyGoldenAccumulatorFishKill(p)
	}
	// 黃金累積魚：任何目標被擊破時累積計數（DAY-214）
	go g.notifyGoldenAccumulatorKill()
	// 幸運鏡像魚：擊破 T173 時觸發鏡像複製（DAY-215）
	if isLuckyMirrorFish(t.DefID) {
		go g.tryLuckyMirrorFish(p)
	}
	// 幸運鏡像魚：擊破鏡像分身時套用倍率加成（DAY-215）
	if g.isLuckyMirrorEntry(t.InstanceID) {
		go g.removeLuckyMirrorEntry(t.InstanceID)
	}
	// 詛咒毒魚：擊破 T174 本身時解除所有詛咒（DAY-216）
	if isCursedPoisonFish(t.DefID) {
		go g.tryCleanseAllCurses(p)
	}
	// 幸運拍賣魚：擊破 T175 本身時提前結束競標（DAY-217）
	if isLuckyAuctionFish(t.DefID) {
		go g.notifyLuckyAuctionFishKill(p)
	}
	// 幸運進化魚：擊破 T176 本身時立即觸發終極爆發（DAY-218）
	if isLuckyEvolutionFish(t.DefID) {
		go g.notifyLuckyEvolutionFishKill(p, t.InstanceID)
	}
	// 幸運連鎖感染魚：擊破 T177 本身時觸發感染（DAY-219）
	if isLuckyInfectionFish(t.DefID) {
		go g.tryLuckyInfectionFish(p)
	}
	// 幸運連鎖感染魚：擊破感染目標時移除感染標記（DAY-219）
	if g.isInfectedTarget(t.InstanceID) {
		go g.notifyLuckyInfectionFishKill(t.InstanceID)
	}
	// 幸運反彈魚：擊破 T178 本身時觸發反彈模式（DAY-220）
	if isLuckyRicochetFish(t.DefID) {
		go g.tryLuckyRicochetFish(p)
	}
	// 幸運黑洞魚：擊破 T179 本身時觸發黑洞（DAY-221）
	if isLuckyBlackHoleFish(t.DefID) {
		go g.tryLuckyBlackHoleFish(p)
	}
	// 幸運黑洞魚：黑洞範圍內目標被擊破時，倍率加成已在 finalMult 計算（DAY-221）
	// 幸運共鳴魚：擊破 T180 本身時觸發共鳴模式（DAY-222）
	if isLuckyResonanceFish(t.DefID) {
		go g.tryLuckyResonanceFish(p)
	}
	// 幸運傳送魚：擊破 T181 本身時觸發傳送漩渦（DAY-223）
	if isLuckyTeleportFish(t.DefID) {
		go g.tryLuckyTeleportFish(p)
	}
	// 幸運分裂魚：擊破 T182 本身時觸發分裂爆炸（DAY-224）
	if isLuckySplitFish(t.DefID) {
		go g.tryLuckySplitFish(p, t.X, t.Y, t.HP)
	}
	// 幸運充能魚：擊破 T183 本身時觸發充能模式（DAY-225）
	if isLuckyChargeFish(t.DefID) {
		go g.tryLuckyChargeFish(p)
	}
	// 幸運鏈鎖爆炸魚：擊破 T184 本身時觸發引爆標記（DAY-226）
	if isLuckyChainBombFish(t.DefID) {
		go g.tryLuckyChainBombFish(p)
	}
	// 幸運鏡像時空魚：擊破 T185 本身時觸發時間倒流（DAY-227）
	if isLuckyMirrorTimeFish(t.DefID) {
		go g.tryLuckyMirrorTimeFish(p)
	}
	// 幸運量子魚：擊破 T186 本身時觸發量子疊加（DAY-228）
	if isLuckyQuantumFish(t.DefID) {
		go g.tryLuckyQuantumFish(p)
	}
	// 幸運寄生魚：擊破 T187 本身時觸發寄生釋放（DAY-229）
	if isLuckyParasiteFish(t.DefID) {
		go g.tryLuckyParasiteFish(p)
	}
	// 幸運風暴魚：擊破 T188 本身時觸發風暴（DAY-230）
	if isLuckyStormFish(t.DefID) {
		go g.tryLuckyStormFish(p)
	}
	// 幸運迴旋鏢魚：擊破 T189 本身時觸發迴旋鏢模式（DAY-231）
	if isLuckyBoomerangFish(t.DefID) {
		go g.tryLuckyBoomerangFish(p)
	}
	// 幸運磁力魚：擊破 T190 本身時觸發磁力場（DAY-232）
	if isLuckyMagnetFish(t.DefID) {
		go g.tryLuckyMagnetFish(p)
	}
	// 幸運回聲魚：擊破 T191 本身時觸發回聲模式（DAY-233）
	if isLuckyEchoFish(t.DefID) {
		go g.tryLuckyEchoFish(p)
	}
	// 幸運漩渦魚：擊破 T192 本身時觸發漩渦模式（DAY-234）
	if isLuckyVortexFish(t.DefID) {
		go g.tryLuckyVortexFish(p)
	}
	// 幸運時間炸彈魚：擊破 T193 本身時觸發時間炸彈（DAY-235）
	if isLuckyTimeBombFish(t.DefID) {
		go g.tryLuckyTimeBombFish(p)
	}
	// 幸運時間炸彈魚：炸彈目標被玩家提前引爆（DAY-235）
	if g.isTimeBombTarget(t.InstanceID) {
		g.notifyTimeBombKill(p, t.InstanceID, t.X, t.Y)
	}
	// 幸運鏡面世界魚：擊破 T194 本身時觸發鏡面世界（DAY-236）
	if isLuckyMirrorWorldFish(t.DefID) {
		go g.tryLuckyMirrorWorldFish(p)
	}
	// 幸運冰凍世界魚：擊破 T195 本身時觸發冰凍世界（DAY-237）
	if isLuckyFreezeWorldFish(t.DefID) {
		go g.tryLuckyFreezeWorldFish(p)
	}
	// 幸運重力反轉魚：擊破 T196 本身時觸發重力反轉（DAY-238）
	if isLuckyGravityFlipFish(t.DefID) {
		go g.tryLuckyGravityFlipFish(p)
	}
	// 幸運共鳴爆發魚：擊破 T197 本身時觸發共鳴爆發（DAY-239）
	if isLuckySynergyBurstFish(t.DefID) {
		go g.tryLuckySynergyBurstFish(p)
	}
	// 幸運賭注魚：擊破 T198 本身時觸發賭注選擇（DAY-240）
	if isLuckyBetFish(t.DefID) {
		go g.tryLuckyBetFish(p)
	}
	// 幸運連鎖反應魚：擊破 T199 本身時觸發連鎖起點標記（DAY-241）
	if isLuckyChainReactionFish(t.DefID) {
		go g.tryLuckyChainReactionFish(p)
	}
	// 幸運連鎖反應：玩家擊破連鎖起點時觸發連鎖（DAY-241）
	if g.isChainReactionStarter(t.InstanceID) {
		g.removeChainReactionStarter(t.InstanceID)
		go g.notifyChainReactionKill(p, t)
	}
	// 幸運分身魚：擊破 T200 本身時觸發分身模式（DAY-242）
	if isLuckyCloneFish(t.DefID) {
		go g.tryLuckyCloneFish(p)
	}
	// 幸運預言魚：擊破 T201 本身時觸發預言（DAY-243）
	if isLuckyProphecyFish(t.DefID) {
		go g.tryLuckyProphecyFish(p)
	}
	// 幸運預言魚：玩家擊破預言目標時觸發預言成真（DAY-243）
	if g.isProphecyTarget(p.ID, t.InstanceID) {
		go g.notifyProphecyKill(p, t.InstanceID, finalReward)
	}
	// 幸運奪旗魚：擊破 T202 本身時觸發搶旗（DAY-244）
	if isLuckyFlagFish(t.DefID) {
		go g.tryLuckyFlagFish(p)
	}
	// 幸運奪旗魚：玩家擊破旗幟目標時觸發奪旗結算（DAY-244）
	if g.isFlagTarget(t.InstanceID) {
		go g.notifyFlagTargetKill(p, t.InstanceID, finalReward)
	}
	// 幸運幽靈魚：擊破 T203 本身時觸發幽靈護盾（DAY-245）
	if isLuckyPhantomFish(t.DefID) {
		go g.tryLuckyPhantomFish(p)
	}
	// 幸運幽靈魚：玩家在幽靈護盾期間擊破任何目標時，留下幽靈殘影（DAY-245）
	if !isLuckyPhantomFish(t.DefID) && g.LuckyPhantomFish.isPhantomShieldActive(p.ID) {
		go g.createPhantomGhost(p, t.DefID, t.X, t.Y)
	}
	// 幸運水晶球魚：擊破 T204 本身時觸發水晶預言（DAY-246）
	if isLuckyCrystalBallFish(t.DefID) {
		go g.tryLuckyCrystalBallFish(p)
	}
	// 幸運水晶球魚：玩家擊破水晶預言目標時觸發必中獎勵（DAY-246）
	if isCrystal, crystalOwner := g.LuckyCrystalBallFish.isCrystalBallTarget(t.InstanceID); isCrystal {
		if crystalOwner == p.ID {
			go g.notifyCrystalBallKill(p, t.InstanceID, finalReward)
		}
	}
	// 幸運時光倒流魚：擊破 T205 本身時觸發時光倒流（DAY-247）
	if isLuckyTimeRewindFish(t.DefID) {
		go g.tryLuckyTimeRewindFish(p)
	}
	// 幸運時光倒流魚：記錄所有擊破歷史（DAY-247，供時光倒流重播使用）
	if !isLuckyTimeRewindFish(t.DefID) {
		name := t.DefID
		if def, ok := data.Targets[t.DefID]; ok {
			name = def.Name
		}
		g.LuckyTimeRewind.recordKillHistory(p.ID, t.InstanceID, t.DefID, name, int(t.Multiplier))
	}
	// 幸運龍捲風魚：擊破 T206 本身時觸發龍捲風（DAY-248）
	if isLuckyTornadoFish(t.DefID) {
		go g.tryLuckyTornadoFish(p)
	}
	// 幸運黑洞爆炸魚：擊破 T207 本身時觸發黑洞（DAY-249）
	if isLuckyBlackHoleExplosionFish(t.DefID) {
		go g.tryLuckyBlackHoleExplosionFish(p)
	}
	// 幸運鏡像分裂魚：擊破 T208 本身時觸發鏡像分裂（DAY-250）
	if isLuckyMirrorSplitFish(t.DefID) {
		go g.tryLuckyMirrorSplitFish(p)
	}
	// 幸運鏡像分裂魚：擊破鏡像副本時給予個人獎勵（DAY-250）
	if isMirror, mirrorEntry := g.LuckyMirrorSplit.isMirrorSplitTarget(t.InstanceID); isMirror && mirrorEntry != nil {
		go g.notifyMirrorSplitKill(p, t.InstanceID, mirrorEntry)
	}
	// 幸運量子糾纏魚：擊破 T209 本身時觸發量子糾纏（DAY-251）
	if isLuckyQuantumEntangleFish(t.DefID) {
		go g.tryLuckyQuantumEntangleFish(p)
	}
	// 幸運量子糾纏魚：擊破糾纏目標時觸發同步爆炸或量子共鳴（DAY-251）
	if isEntangled, sessionID := g.LuckyQuantumEntangle.isQuantumEntangleTarget(t.InstanceID); isEntangled && sessionID != "" {
		go g.notifyQuantumEntangleKill(p, t.InstanceID, sessionID)
	}
	// 幸運回聲魚：玩家在回聲模式中擊破任何目標時，觸發回聲分身（DAY-233）
	if !isLuckyEchoFish(t.DefID) && g.isEchoModeActive(p.ID) {
		go g.notifyEchoKill(p, t, t.InstanceID)
	}
	// 幸運回聲魚：擊破回聲分身時，觸發下一層回聲（DAY-233）
	if isEcho, echoEntry := g.isEchoTarget(t.InstanceID); isEcho {
		g.removeEchoEntry(t.InstanceID)
		go g.notifyEchoTargetKill(p, t, t.InstanceID, echoEntry.layer)
	}
	// S-Rank 傳說目標召喚深淵巨鯨：擊破傳說品質目標後 15% 機率觸發（DAY-165）
	if t.Quality == target.QualityLegendary && !isAbyssWhale(t.DefID) {
		go g.tryLegendarySummonWhale(p, t.X, t.Y)
	}
	// 特殊武器自動充能：每次擊破累積充能進度（DAY-134）
	go g.notifySpecialWeaponCharge(p, t.Multiplier)
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
			// 競速獵殺：目標消失時取消競速（DAY-136）
			go g.cancelSpeedRaceIfTarget(id)
			// 懸賞：目標消失時取消懸賞並退款（DAY-137）
			go g.cancelBountyForTarget(id)
			// 龍龜不死 Boss：離開畫面時廣播結算（DAY-186）
			if isDragonTurtle(t.DefID) {
				go g.notifyDragonTurtleLeave(id)
			}
			// 巨型鱷魚：離開畫面時廣播結算（DAY-188）
			if isCrocodileHunter(t.DefID) {
				go g.onCrocodileHunterLeave(id, "timeout")
			}
			// 巨型鮟鱇魚：離開畫面時廣播結算（DAY-196）
			if isAnglerfishElectric(t.DefID) {
				go g.onAnglerfishLeave(id, "timeout")
			}
			// 幽靈魚真身：離開畫面時移除所有幻影分身（DAY-198）
			if isGhostFish(t.DefID) {
				go g.onGhostFishLeave(id)
			}
			// 黃金累積魚：離開畫面時通知累積系統（DAY-214）
			if isGoldenAccumulatorFish(t.DefID) {
				go g.notifyGoldenAccumulatorLeave(id)
			}
			// 詛咒毒魚：詛咒目標離開畫面時觸發懲罰（DAY-216）
			if g.isCursedTarget(id) {
				go g.notifyCursedTargetEscape(id)
			}
			// 幸運進化魚：進化魚離開畫面時清除狀態（DAY-218）
			if isLuckyEvolutionFish(t.DefID) {
				go g.notifyLuckyEvolutionFishLeave(id)
			}
			// 幸運回聲魚：回聲分身消失時清除記錄（DAY-233）
			go g.notifyEchoTargetExpire(id)
			// 幸運預言魚：預言目標消失時嘗試轉移預言（DAY-243）
			go g.notifyProphecyTargetGone(id)
			// 幸運奪旗魚：旗幟目標消失時觸發自動爆炸（DAY-244）
			go g.notifyFlagTargetGone(id)
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
	// 每日賽排名廣播（每 30 秒，DAY-093）
	shouldBroadcastDailyTournament := now.Sub(g.lastDailyTournamentAt) >= 30*time.Second
	if shouldBroadcastDailyTournament {
		g.lastDailyTournamentAt = now
	}
	// 多格式賽排名廣播（每 30 秒，DAY-111）
	shouldBroadcastMultiFormat := now.Sub(g.lastMultiFormatAt) >= 30*time.Second
	if shouldBroadcastMultiFormat {
		g.lastMultiFormatAt = now
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
	// 天氣廣播（每 30 秒，DAY-087）
	shouldBroadcastWeather := now.Sub(g.lastWeatherAt) >= 30*time.Second
	if shouldBroadcastWeather {
		g.lastWeatherAt = now
	}
	// 玩家資料定期自動儲存（每 60 秒，DAY-099）
	shouldAutoSave := now.Sub(g.lastAutoSaveAt) >= 60*time.Second
	if shouldAutoSave {
		g.lastAutoSaveAt = now
	}
	// Co-op Boss Raid tick（每 3 秒，DAY-115）
	shouldTickRaid := now.Sub(g.lastRaidTickAt) >= 3*time.Second
	if shouldTickRaid {
		g.lastRaidTickAt = now
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
	// 每日賽排名廣播（每 30 秒，DAY-093）
	if shouldBroadcastDailyTournament {
		go g.broadcastDailyTournament()
	}
	// 多格式賽排名廣播（每 30 秒，DAY-111）
	if shouldBroadcastMultiFormat {
		go g.broadcastMultiFormat()
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
	// 天氣 Tick + 廣播（每 30 秒，DAY-087）
	if shouldBroadcastWeather {
		go g.tickAndBroadcastWeather()
	}
	// 玩家資料定期自動儲存（每 60 秒，DAY-099）
	if shouldAutoSave {
		go g.autoSaveAllPlayers()
	}
	// Co-op Boss Raid tick（每 3 秒，DAY-115）
	if shouldTickRaid {
		go g.tickRaidUpdate()
	}
	// 閃電挑戰 tick（每次 update，DAY-123）
	go g.tickFlashChallenge()
	// 黃金時間 tick（每次 update，DAY-125）
	go g.tickGoldenTime()
	// 稀有連擊過期檢查（每 5 秒，DAY-126）
	if now.Sub(g.lastRareCatchTickAt) >= 5*time.Second {
		g.lastRareCatchTickAt = now
		go g.tickRareCatchExpiry()
	}	// 天氣湧現事件過期檢查（每次 update，DAY-127）
	g.tickWeatherSurge()
	// 不死 BOSS 過期檢查（每次 update，DAY-129）
	go g.tickImmortalBoss()
	// 覺醒 BOSS 過期檢查（每次 update，DAY-130）
	go g.tickAwakenBoss()
	// 連勝超時檢查（每次 update，DAY-131）
	go g.tickWinStreakExpiry()
	// 狂熱模式過期檢查（每次 update，DAY-133）
	go g.tickFeverModeExpiry()
	// 競速獵殺超時檢查（每次 update，DAY-136）
	go g.tickSpeedRace()
	// 懸賞過期檢查（每次 update，DAY-137）
	go g.tickBountyExpiry()
	// 倍率風暴 tick（每秒，DAY-138）
	if now.Sub(g.lastRareCatchTickAt) >= 1*time.Second {
		go g.tickMultStorm()
	}
	// 雙環輪盤超時自動停止（每次 update，DAY-139）
	go g.tickDualRoulette()
	// Mega Catch 過期檢查（每次 update，DAY-140）
	go g.tickMegaCatch()
	// Mega Catch 每分鐘隨機觸發（DAY-140）
	if now.Sub(g.lastMegaCatchTickAt) >= 60*time.Second {
		g.lastMegaCatchTickAt = now
		go g.tryMegaCatchRandom()
	}
	// 千龍王輪盤超時自動停止（每次 update，DAY-148）
	go g.tickChainLongWheel()
	// 水晶龍水晶衰減檢查（每次 update，DAY-153）
	go g.tickCrystalDragonDecay()
	// 輪盤螃蟹超時檢查（每次 update，DAY-167）
	go g.tickRouletteCrabWheel()
	// Rapid Respin session 過期檢查（每次 update，DAY-121）
	g.checkRespinSessionExpiry()
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

	// 整合天氣加成（DAY-127）
	rareBonus := 0.0
	goldFishBonus := 0.0
	if g.Weather != nil {
		rareBonus = g.Weather.GetRareChanceBonus()
		goldFishBonus = g.Weather.GetGoldFishBonus()
	}
	// 天氣湧現事件加成（DAY-127）
	if g.weatherSurgeActive {
		rareBonus += g.weatherSurgeRareBonus
		goldFishBonus += g.weatherSurgeGoldBonus
	}
	// Mega Catch 生成加成（DAY-140）
	rareBonus += g.getMegaCatchSpawnBoost()
	def := g.SpawnSys.PickTargetDef(betLevel, bonusSpecial, rareBonus, goldFishBonus)
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

	// 傳說/史詩品質目標出現公告（DAY-124）
	if t.Quality == target.QualityLegendary {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgRareTargetAlert,
			Payload: ws.RareTargetAlertPayload{
				InstanceID: instanceID,
				DefID:      def.ID,
				Name:       def.Name,
				Quality:    string(t.Quality),
				Multiplier: int(t.Multiplier),
				Icon:       "⭐",
				Message:    "傳說目標出現！" + def.Name + "（×" + fmt.Sprintf("%.0f", t.Multiplier) + "）",
				Color:      "#FFD700",
			},
		})
	} else if t.Quality == target.QualityEpic {
		g.Hub.Broadcast(&ws.Message{
			Type: ws.MsgRareTargetAlert,
			Payload: ws.RareTargetAlertPayload{
				InstanceID: instanceID,
				DefID:      def.ID,
				Name:       def.Name,
				Quality:    string(t.Quality),
				Multiplier: int(t.Multiplier),
				Icon:       "💜",
				Message:    "史詩目標出現！" + def.Name + "（×" + fmt.Sprintf("%.0f", t.Multiplier) + "）",
				Color:      "#9B59B6",
			},
		})
	}

	// 不死 BOSS：每次生成目標時嘗試觸發（DAY-129）
	go g.trySpawnImmortalBoss()
	// 覺醒 BOSS：每次生成目標時嘗試觸發（DAY-130）
	go g.trySpawnAwakenBoss()
	// 競速獵殺：高倍率目標生成時嘗試觸發（DAY-136）
	if t.Multiplier >= 10.0 {
		go g.tryStartSpeedRace(instanceID, def.ID, def.Name, t.Multiplier)
	}
	// 深淵巨鯨：T124 生成時通知全服（DAY-164）
	if isAbyssWhale(def.ID) {
		go g.notifyAbyssWhaleSpawn(instanceID, x, y)
	}
	// 龍龜不死 Boss：T144 生成時通知全服（DAY-186）
	if isDragonTurtle(def.ID) {
		go g.notifyDragonTurtleSpawn(instanceID, x, y)
	}
	// 巨型鱷魚：T146 生成時觸發獵食模式（DAY-188）
	if isCrocodileHunter(def.ID) {
		go g.tryCrocodileHunterSpawn(instanceID)
	}
	// 時間炸彈魚：T147 生成時觸發倒數（DAY-189）
	if isTimeBombFish(def.ID) {
		go g.tryTimeBombFishSpawn(instanceID)
	}
	// 巨型鮟鱇魚：T154 生成時觸發電擊模式（DAY-196）
	if isAnglerfishElectric(def.ID) {
		go g.notifyAnglerfishSpawn(instanceID)
	}
	// 幽靈魚：T156 生成時觸發幻影分身（DAY-198）
	if isGhostFish(def.ID) {
		go g.notifyGhostFishSpawn(instanceID, x, y, t.Multiplier)
	}
	// 彗星魚：T164 生成時觸發軌跡飛行（DAY-206）
	if isCometFish(def.ID) {
		go g.notifyCometFishSpawn(t)
	}
	// 黃金累積魚：T172 生成時啟動累積系統（DAY-214）
	if isGoldenAccumulatorFish(def.ID) {
		go g.notifyGoldenAccumulatorSpawn(t)
	}
	// 詛咒毒魚：T174 生成時自動詛咒場上目標（DAY-216）
	// 注意：詛咒毒魚生成時立即詛咒，不需要玩家擊破才觸發
	if isCursedPoisonFish(def.ID) {
		// 用一個虛擬玩家觸發詛咒（取第一個在線玩家，或直接呼叫內部函數）
		go g.spawnCursedPoisonMarks()
	}
	// 幸運拍賣魚：T175 生成時開啟全服競標（DAY-217）
	if isLuckyAuctionFish(def.ID) {
		go g.notifyLuckyAuctionFishSpawn(t)
	}
	// 幸運進化魚：T176 生成時初始化進化狀態（DAY-218）
	if isLuckyEvolutionFish(def.ID) {
		go g.notifyLuckyEvolutionFishSpawn(t)
	}
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
	// VIP（DAY-106）
	VIPLevel    int                    `json:"vip_level"`
	VIPName     string                 `json:"vip_name"`
	// 公會（DAY-106）
	GuildID     string                 `json:"guild_id"`
	GuildName   string                 `json:"guild_name"`
	GuildRole   string                 `json:"guild_role"`
	// 統計亮點（DAY-106）
	BestStreak  int                    `json:"best_streak"`
	BestMult    float64                `json:"best_mult"`
	JackpotWins int                    `json:"jackpot_wins"`
	TotalBet    int                    `json:"total_bet"`
	TotalReward int                    `json:"total_reward"`
	RTP         float64                `json:"rtp"`
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

	// VIP 資訊（DAY-106）
	vipSnap := g.VIP.GetSnapshot(playerID)
	vipLevel := vipSnap.VIPLevel
	vipName := vipSnap.TierName

	// 公會資訊（DAY-106）
	guildID := g.Guild.GetPlayerGuildID(playerID)
	guildName := ""
	guildRole := ""
	if guildID != "" {
		if gd := g.Guild.GetGuild(guildID); gd != nil {
			guildName = gd.Name
			if member, ok := gd.Members[playerID]; ok {
				guildRole = string(member.Role)
			}
		}
	}

	// 統計亮點（DAY-106）
	bestStreak := 0
	bestMult := 0.0
	jackpotWins := 0
	totalBet := 0
	totalReward := 0
	rtp := 0.0
	if p.Stats != nil {
		statsSnap := p.Stats.Snapshot()
		bestStreak = statsSnap.BestStreak
		bestMult = statsSnap.BestMultiplier
		jackpotWins = statsSnap.JackpotWins
		totalBet = statsSnap.TotalBet
		totalReward = statsSnap.TotalReward
		rtp = statsSnap.RTP
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
		VIPLevel:       vipLevel,
		VIPName:        vipName,
		GuildID:        guildID,
		GuildName:      guildName,
		GuildRole:      guildRole,
		BestStreak:     bestStreak,
		BestMult:       bestMult,
		JackpotWins:    jackpotWins,
		TotalBet:       totalBet,
		TotalReward:    totalReward,
		RTP:            rtp,
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
	// 動態牆廣播（DAY-112）
	g.mu.RLock()
	p, pok := g.Players[playerID]
	g.mu.RUnlock()
	if pok {
		g.notifyFeedAchievement(p, u.Name, u.Icon, u.Type)
	}
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
				// 動態牆：稱號獲得（DAY-112）
				g.notifyFeedTitle(p, titleDef.Name, titleDef.Icon, titleDef.Priority)
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
