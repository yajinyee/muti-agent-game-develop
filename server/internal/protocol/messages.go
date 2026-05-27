// Package protocol — WebSocket 訊息協定定義
// server-core-agent + protocol-sync-agent 負責維護
package protocol

// ── Client → Server ──────────────────────────────────────────
const (
	MsgAttack            = "attack"
	MsgLock              = "lock"
	MsgAutoToggle        = "auto_toggle"
	MsgBetChange         = "bet_change"
	MsgBonusClick        = "bonus_click"
	MsgPing              = "ping"
	MsgTriggerBoss       = "trigger_boss"
	MsgTriggerBonus      = "trigger_bonus"
	MsgCollectGoldenCoin = "collect_golden_coin" // T122 黃金雨魚：收集黃金幣
	MsgSetDisplayName    = "set_display_name"    // 設定玩家顯示名稱
)

// ── Server → Client ──────────────────────────────────────────
const (
	MsgGameState    = "game_state"
	MsgTargetSpawn  = "target_spawn"
	MsgTargetUpdate = "target_update"
	MsgTargetKill   = "target_kill"
	MsgAttackResult = "attack_result"
	MsgReward       = "reward"
	MsgBossEvent    = "boss_event"
	MsgBonusEvent   = "bonus_event"
	MsgPlayerUpdate = "player_update"
	MsgError        = "error"
	MsgPong         = "pong"
	MsgAnnounce     = "announce"

	// 幸運特殊魚事件
	MsgLuckyChainLightning  = "lucky_chain_lightning"  // T106 連鎖閃電
	MsgLuckyCrabTorpedo     = "lucky_crab_torpedo"     // T107 螃蟹魚雷
	MsgLuckyVortex          = "lucky_vortex"           // T108 渦旋海葵
	MsgLuckyGoldenDragon    = "lucky_golden_dragon"    // T109 黃金龍魚輪盤
	MsgLuckyThunderLobster  = "lucky_thunder_lobster"  // T110 雷霆龍蝦
	MsgLuckyAwakenedPhoenix = "lucky_awakened_phoenix" // T111 覺醒鳳凰
	MsgLuckyShockwaveBomb   = "lucky_shockwave_bomb"   // T112 全場震盪
	MsgLuckyDrillTorpedo    = "lucky_drill_torpedo"    // T113 鑽頭魚雷
	MsgLuckyTimeFreeze      = "lucky_time_freeze"      // T114 時間凍結
	MsgLuckyChainExplosion  = "lucky_chain_explosion"  // T115 連鎖爆炸

	// DAY-295 新增幸運特殊魚事件
	MsgLuckyChainLongKing   = "lucky_chain_long_king"   // T116 千龍王輪盤
	MsgLuckyDragonShotgun   = "lucky_dragon_shotgun"    // T117 龍力散彈
	MsgLuckyRocketCannon    = "lucky_rocket_cannon"     // T118 火箭砲
	MsgLuckyDeepWhirlpool   = "lucky_deep_whirlpool"    // T119 深海漩渦
	MsgLuckyVampireMult     = "lucky_vampire_mult"      // T120 吸血鬼倍率

	// DAY-296 新增幸運特殊魚事件
	MsgLuckyMirrorFish   = "lucky_mirror_fish"   // T121 鏡像魚
	MsgLuckyGoldenRain   = "lucky_golden_rain"   // T122 黃金雨魚
	MsgLuckyFreezeBomb   = "lucky_freeze_bomb"   // T123 冰凍炸彈魚
	MsgLuckyThunderStorm = "lucky_thunder_storm" // T124 雷暴魚
	MsgLuckyLuckyWheel   = "lucky_lucky_wheel"   // T125 大轉盤魚

	// DAY-301 新增幸運特殊魚事件
	MsgLuckyJackpotFish = "lucky_jackpot_fish" // T126 進階 Jackpot 魚
	MsgLuckyCoopFish    = "lucky_coop_fish"    // T127 全服合作魚
	MsgLuckyTimeWarp    = "lucky_time_warp"    // T128 時間扭曲魚

	// DAY-302 新增幸運特殊魚事件
	MsgLuckyChainMeteor = "lucky_chain_meteor" // T129 連鎖隕石魚

	// DAY-303 新增幸運特殊魚事件
	MsgLuckyCrashFish = "lucky_crash_fish" // T130 崩潰魚（Crash mechanic）

	// DAY-304 新增幸運特殊魚事件
	MsgLuckyElectricEel  = "lucky_electric_eel"  // T131 電鰻魚（持續放電連鎖）
	MsgLuckyAnglerFish   = "lucky_angler_fish"   // T132 巨型安康魚（誘餌+電擊爆炸）
	MsgLuckyBlackHole    = "lucky_black_hole"    // T133 黑洞魚（吸引+坍縮）
	MsgLuckyBountyHunter = "lucky_bounty_hunter" // T134 賞金獵人魚（賞金目標系統）
	MsgLuckyTsunami      = "lucky_tsunami"       // T135 海嘯魚（三波衝擊）

	// DAY-305 新增幸運特殊魚事件
	MsgLuckyDragonWrathV2  = "lucky_dragon_wrath_v2"  // T136 龍怒蓄積魚 v2（升級版）
	MsgLuckyHumpbackWhale  = "lucky_humpback_whale"   // T137 座頭鯨魚（鯨歌共鳴）
	MsgLuckyLegendDragon   = "lucky_legend_dragon"    // T138 傳說龍魚（龍息噴火）
	MsgLuckyGuildWar       = "lucky_guild_war"        // T139 公會戰魚（全服積分）
	MsgLuckyQualityFish    = "lucky_quality_fish"     // T140 品質魚（品質鑑定）

	// DAY-306 新增幸運特殊魚事件
	MsgLuckyTornado      = "lucky_tornado"       // T141 龍捲風魚（橫掃全場）
	MsgLuckyEarthquake   = "lucky_earthquake"    // T142 地震魚（三波同心圓）
	MsgLuckyVolcano      = "lucky_volcano"       // T143 火山魚（熔岩彈雨）
	MsgLuckyCosmicRay    = "lucky_cosmic_ray"    // T144 星際魚（8方向光束）
	MsgLuckyDivineDragon = "lucky_divine_dragon" // T145 神龍魚（神龍降臨）

	// DAY-307 新增幸運特殊魚事件
	MsgLuckyQuantum  = "lucky_quantum"  // T146 量子魚（量子觀測坍縮）
	MsgLuckySupernova = "lucky_supernova" // T147 超新星魚（全場爆炸+倍率加成）
	MsgLuckyInfinite = "lucky_infinite" // T148 無限魚（無限累積倍率）
	MsgLuckyGenesis  = "lucky_genesis"  // T149 創世魚（全場審判）
	MsgLuckyRebirth  = "lucky_rebirth"  // T150 重生魚（死亡目標復活再擊破）

	// DAY-308 新增幸運特殊魚事件
	MsgLuckyAwakenedCroc = "lucky_awakened_croc" // T151 覺醒鱷魚（自動獵魚 20 秒）
	MsgLuckyVampireV2    = "lucky_vampire_v2"    // T152 吸血鬼升級魚（倍率上限 ×10.0）
	MsgLuckySuperAwaken  = "lucky_super_awaken"  // T153 超級覺醒魚（全場審判 + 全服 ×7.0）
	MsgLuckyGiantPrize   = "lucky_giant_prize"   // T154 巨型獎勵魚（5 次隨機大獎）
	MsgLuckyImmortalBoss = "lucky_immortal_boss" // T155 不死 BOSS 魚（5 條命遞增倍率）

	// DAY-309 新增幸運特殊魚事件
	MsgLuckyIcePhoenix       = "lucky_ice_phoenix"       // T156 冰鳳凰魚（冰凍+鳳凰重生）
	MsgLuckyDragonFury       = "lucky_dragon_fury"       // T157 龍怒能量魚（能量累積→全場攻擊）
	MsgLuckyMultCascade      = "lucky_mult_cascade"      // T158 倍率瀑布魚（連續擊破累積倍率）
	MsgLuckyAwakenBossV2     = "lucky_awaken_boss_v2"    // T159 覺醒 BOSS 魚 v2（8 次 Power Up）
	MsgLuckyUltimateJudgment = "lucky_ultimate_judgment" // T160 終極審判魚（全場清空 + 全服 ×10.0）

	// DAY-310 新增幸運特殊魚事件
	MsgLuckyComboBurst      = "lucky_combo_burst"      // T161 連擊爆發魚（連擊累積倍率 ×15.0）
	MsgLuckyTimeBomb        = "lucky_time_bomb"        // T162 時間炸彈魚（30 秒倒數爆炸）
	MsgLuckyElementalFusion = "lucky_elemental_fusion" // T163 元素融合魚（火/冰/雷三元素）
	MsgLuckyTreasureHunter  = "lucky_treasure_hunter"  // T164 寶藏獵人魚（5 個隨機寶藏）
	MsgLuckyMythAwaken      = "lucky_myth_awaken"      // T165 神話覺醒魚（全場倍率 ×3.0，25 秒）

	// DAY-312 新增幸運特殊魚事件
	MsgLuckyStarPortal    = "lucky_star_portal"    // T166 星際門戶魚（傳送 5 個目標到中央）
	MsgLuckyDragonSoul    = "lucky_dragon_soul"    // T167 龍魂融合魚（吸收 50 魂→全場 HP -90%）
	MsgLuckySpacetimeRift = "lucky_spacetime_rift" // T168 時空裂縫魚（每 4 秒瞬間擊破 3 個）
	MsgLuckyHolyJudgment  = "lucky_holy_judgment"  // T169 神聖審判魚（5 波神聖光柱 HP -30%）
	MsgLuckyBigBang       = "lucky_big_bang"       // T170 宇宙大爆炸魚（全場清空 + 全服 ×12.0）

	// DAY-313 新增幸運特殊魚事件（Progressive Jackpot 系列）
	MsgLuckyJackpotPool = "lucky_jackpot_pool" // T171-T175 累積獎池系統（Mini/Minor/Major/Grand）

	// DAY-314 新增幸運特殊魚事件
	MsgLuckyMultiverse  = "lucky_multiverse"   // T176 多重宇宙魚（3 個平行宇宙，全服 ×13.0）
	MsgLuckyTimeLoop    = "lucky_time_loop"    // T177 時間迴圈魚（3 次迴圈，獎勵 ×1.5 遞增，全服 ×10.0）
	MsgLuckyFateWheel   = "lucky_fate_wheel"   // T178 命運之輪魚（3 次旋轉，最高 ×50.0，全服 ×11.0）
	MsgLuckyDivineRealm = "lucky_divine_realm" // T179 神域降臨魚（5 波神域光柱 HP -35%，全服 ×14.0）
	MsgLuckyFinalPower  = "lucky_final_power"  // T180 終焉之力魚（全場清空 ×10.0，全服 ×15.0 新最高）

	// DAY-315 新增幸運特殊魚事件
	MsgLuckyMutation    = "lucky_mutation"     // T181 突變魚（150種突變，最高 ×17.0，全服 ×16.0）
	MsgLuckyArcticStorm = "lucky_arctic_storm" // T182 北極風暴魚（8 波快速連擊，全服 ×16.5）
	MsgLuckyFisherWild  = "lucky_fisher_wild"  // T183 漁夫野生魚（3 個 Wild 目標，全服 ×17.0）
	MsgLuckyRiskLevel   = "lucky_risk_level"   // T184 風險等級魚（5 等級選擇，最高 ×3000，全服 ×17.5）
	MsgLuckyCosmicPulse = "lucky_cosmic_pulse" // T185 宇宙脈衝魚（全場 HP -45%，全服 ×16.0 新最高）

	// DAY-316 新增幸運特殊魚事件
	MsgLuckyMirrorUniverse   = "lucky_mirror_universe"   // T186 鏡像宇宙魚（複製最強 3 個目標，全服 ×17.0）
	MsgLuckyGravityField     = "lucky_gravity_field"     // T187 引力場魚（引力吸引+爆炸 HP -55%，全服 ×17.5）
	MsgLuckyTimeAcceleration = "lucky_time_acceleration" // T188 時間加速魚（目標速度 ×0.15，射擊 ×3.0，全服 ×18.0 新最高）
	MsgLuckyNebulaVortex     = "lucky_nebula_vortex"     // T189 星雲漩渦魚（每秒 HP -8%，持續 20 秒，全服 ×18.5）
	MsgLuckyCosmicJudgment   = "lucky_cosmic_judgment"   // T190 宇宙審判魚（全場清空 ×14.0，全服 ×19.0 新最高）
)
// ── Envelope ─────────────────────────────────────────────────

// Envelope は WebSocket の送受信フォーマット
type Envelope struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// ── Payloads: Server → Client ────────────────────────────────

type GameStatePayload struct {
	State     string `json:"state"`
	Timestamp int64  `json:"timestamp"`
}

type TargetSpawnPayload struct {
	InstanceID string  `json:"instance_id"`
	DefID      string  `json:"def_id"`
	Name       string  `json:"name"`
	Type       string  `json:"type"` // basic | special | boss
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	HP         int     `json:"hp"`
	MaxHP      int     `json:"max_hp"`
	Speed      float64 `json:"speed"`
	Lifetime   float64 `json:"lifetime"`
	Behavior   string  `json:"behavior"` // linear | sink | flee | fast
	Multiplier float64 `json:"multiplier"`
}

type TargetUpdatePayload struct {
	InstanceID string  `json:"instance_id"`
	HP         int     `json:"hp"`
	MaxHP      int     `json:"max_hp"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	IsFleeing  bool    `json:"is_fleeing,omitempty"`
}

type TargetKillPayload struct {
	InstanceID string  `json:"instance_id"`
	DefID      string  `json:"def_id"`
	Multiplier float64 `json:"multiplier"`
	Reward     int     `json:"reward"`
	LaborGain  int     `json:"labor_gain"`
	KillerID   string  `json:"killer_id"`
}

type AttackResultPayload struct {
	TargetID    string  `json:"target_id"`
	IsHit       bool    `json:"is_hit"`
	IsKill      bool    `json:"is_kill"`
	Damage      int     `json:"damage"`
	Reward      int     `json:"reward"`
	LaborGain   int     `json:"labor_gain"`
	CharacterID string  `json:"character_id"`
	Multiplier  float64 `json:"multiplier"`
	// DAY-311 加入命中位置（供 Client 特效使用）
	PosX float64 `json:"pos_x"`
	PosY float64 `json:"pos_y"`
}

type RewardPayload struct {
	Source     string  `json:"source"`
	Amount     int     `json:"amount"`
	Multiplier float64 `json:"multiplier"`
	NewBalance int     `json:"new_balance"`
}

type BossEventPayload struct {
	Event      string  `json:"event"` // warning | spawn | phase_change | kill | timeout
	InstanceID string  `json:"instance_id"`
	Phase      int     `json:"phase"`
	HP         int     `json:"hp"`
	MaxHP      int     `json:"max_hp"`
	Reward     int     `json:"reward,omitempty"`
	Multiplier float64 `json:"multiplier,omitempty"`
}

type BonusEventPayload struct {
	Event      string  `json:"event"` // start | tick | end | result
	TimeLeft   float64 `json:"time_left"`
	Score      int     `json:"score"`
	Multiplier float64 `json:"multiplier"`
	Reward     int     `json:"reward,omitempty"`
}

type PlayerUpdatePayload struct {
	ID             string  `json:"id"`
	Coins          int     `json:"coins"`
	BetLevel       int     `json:"bet_level"`
	BetCost        int     `json:"bet_cost"`
	CharacterID    string  `json:"character_id"`
	CharacterName  string  `json:"character_name"`
	LaborValue     int     `json:"labor_value"`
	IsAuto         bool    `json:"is_auto"`
	LockTargetID   string  `json:"lock_target_id"`
	ProjectileSpeed float64 `json:"projectile_speed"`
	FireRate       float64 `json:"fire_rate"`
	// DAY-297 Combo 系統
	ComboCount     int     `json:"combo_count"`
	ComboMultBonus float64 `json:"combo_mult_bonus"`
}

type AnnouncePayload struct {
	Message  string `json:"message"`
	Priority string `json:"priority"` // low | normal | high | critical
	Color    string `json:"color"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ── Payloads: Client → Server ────────────────────────────────

type AttackRequest struct {
	TargetID string  `json:"target_id"`
	ClickX   float64 `json:"click_x"`
	ClickY   float64 `json:"click_y"`
}

type LockRequest struct {
	TargetID string `json:"target_id"`
}

type BetChangeRequest struct {
	BetLevel int `json:"bet_level"`
}

type BonusClickRequest struct {
	TargetID string  `json:"target_id"`
	ClickX   float64 `json:"click_x"`
	ClickY   float64 `json:"click_y"`
}

// ── Lucky 特殊魚 Payloads ─────────────────────────────────────

// LuckyChainLightningPayload T106 連鎖閃電
type LuckyChainLightningPayload struct {
	Event        string   `json:"event"`         // "trigger" | "chain_hit" | "settle"
	TriggerID    string   `json:"trigger_id"`    // 觸發玩家 ID
	TriggerName  string   `json:"trigger_name"`  // 觸發玩家名稱
	HitTargets   []string `json:"hit_targets"`   // 被閃電命中的目標 instance_id 列表
	ChainCount   int      `json:"chain_count"`   // 連鎖次數
	TotalReward  int      `json:"total_reward"`  // 總獎勵
	Multiplier   float64  `json:"multiplier"`    // 最終倍率
}

// LuckyCrabTorpedoPayload T107 螃蟹魚雷
type LuckyCrabTorpedoPayload struct {
	Event       string    `json:"event"`        // "trigger" | "explosion" | "settle"
	TriggerID   string    `json:"trigger_id"`
	TriggerName string    `json:"trigger_name"`
	ExplosionX  float64   `json:"explosion_x"`  // 爆炸中心 X
	ExplosionY  float64   `json:"explosion_y"`  // 爆炸中心 Y
	HitTargets  []string  `json:"hit_targets"`  // 被爆炸命中的目標
	ExplosionNo int       `json:"explosion_no"` // 第幾次爆炸（1-3）
	TotalReward int       `json:"total_reward"`
}

// LuckyVortexPayload T108 渦旋海葵
type LuckyVortexPayload struct {
	Event       string  `json:"event"`        // "trigger" | "pull" | "damage" | "end"
	TriggerID   string  `json:"trigger_id"`
	TriggerName string  `json:"trigger_name"`
	TimeLeft    float64 `json:"time_left"`    // 渦旋剩餘時間
	HitCount    int     `json:"hit_count"`    // 本次傷害命中數
	TotalReward int     `json:"total_reward"`
}

// LuckyGoldenDragonPayload T109 黃金龍魚輪盤
type LuckyGoldenDragonPayload struct {
	Event       string  `json:"event"`        // "trigger" | "spin" | "result"
	TriggerID   string  `json:"trigger_id"`
	TriggerName string  `json:"trigger_name"`
	InnerMult   float64 `json:"inner_mult"`   // 內環倍率
	OuterMult   float64 `json:"outer_mult"`   // 外環倍率
	FinalMult   float64 `json:"final_mult"`   // 最終倍率（內×外）
	Reward      int     `json:"reward"`
}

// LuckyThunderLobsterPayload T110 雷霆龍蝦
type LuckyThunderLobsterPayload struct {
	Event       string  `json:"event"`        // "trigger" | "auto_fire" | "end"
	TriggerID   string  `json:"trigger_id"`
	TriggerName string  `json:"trigger_name"`
	TimeLeft    float64 `json:"time_left"`    // 自動射擊剩餘時間
	KillCount   int     `json:"kill_count"`   // 已擊破數
	TotalReward int     `json:"total_reward"`
}

// LuckyAwakenedPhoenixPayload T111 覺醒鳳凰
type LuckyAwakenedPhoenixPayload struct {
	Event       string  `json:"event"`        // "awaken_start" | "power_up" | "perfect_awaken" | "perfect_end" | "awaken_end"
	TriggerID   string  `json:"trigger_id"`
	TriggerName string  `json:"trigger_name"`
	PowerUpMult float64 `json:"power_up_mult,omitempty"` // 本次 Power Up 倍率（6-10x）
	ShotsLeft   int     `json:"shots_left,omitempty"`    // 剩餘 Power Up 次數
	HitCount    int     `json:"hit_count,omitempty"`     // 命中次數
	TotalReward int     `json:"total_reward,omitempty"`
}

// LuckyShockwaveBombPayload T112 全場震盪
type LuckyShockwaveBombPayload struct {
	Event       string `json:"event"`        // "shockwave_start" | "shockwave_hit" | "super_shockwave" | "super_end" | "power_end"
	TriggerID   string `json:"trigger_id"`
	TriggerName string `json:"trigger_name"`
	HitCount    int    `json:"hit_count,omitempty"`    // 震盪命中目標數
	TotalReward int    `json:"total_reward,omitempty"`
}

// LuckyDrillTorpedoPayload T113 鑽頭魚雷
type LuckyDrillTorpedoPayload struct {
	Event        string   `json:"event"`          // "trigger" | "penetrate" | "explode" | "perfect" | "perfect_end"
	TriggerID    string   `json:"trigger_id"`
	TriggerName  string   `json:"trigger_name"`
	HitTargets   []string `json:"hit_targets,omitempty"`
	PenetrateCnt int      `json:"penetrate_cnt,omitempty"`
	ExplodeX     float64  `json:"explode_x,omitempty"`
	ExplodeY     float64  `json:"explode_y,omitempty"`
	AccumMult    float64  `json:"accum_mult,omitempty"`
	TotalReward  int      `json:"total_reward,omitempty"`
}

// LuckyTimeFreezePayload T114 時間凍結
type LuckyTimeFreezePayload struct {
	Event       string  `json:"event"`        // "freeze_start" | "freeze_end" | "perfect_freeze" | "perfect_end"
	TriggerID   string  `json:"trigger_id"`
	TriggerName string  `json:"trigger_name"`
	Duration    float64 `json:"duration,omitempty"`
	KillCount   int     `json:"kill_count,omitempty"`
}

// LuckyChainExplosionPayload T115 連鎖爆炸
type LuckyChainExplosionPayload struct {
	Event       string   `json:"event"`        // "chain_start" | "chain_explode" | "chain_burst" | "burst_end" | "chain_end"
	TriggerID   string   `json:"trigger_id"`
	TriggerName string   `json:"trigger_name"`
	Duration    float64  `json:"duration,omitempty"`
	ExplodeX    float64  `json:"explode_x,omitempty"`
	ExplodeY    float64  `json:"explode_y,omitempty"`
	HitTargets  []string `json:"hit_targets,omitempty"`
	ChainCount  int      `json:"chain_count,omitempty"`
	AccumMult   float64  `json:"accum_mult,omitempty"`
	TotalReward int      `json:"total_reward,omitempty"`
}

// LuckyChainLongKingPayload T116 千龍王輪盤
type LuckyChainLongKingPayload struct {
	Event       string  `json:"event"`        // "trigger" | "spin" | "result" | "mega_win"
	TriggerID   string  `json:"trigger_id"`
	TriggerName string  `json:"trigger_name"`
	InnerMult   float64 `json:"inner_mult,omitempty"`
	OuterMult   float64 `json:"outer_mult,omitempty"`
	FinalMult   float64 `json:"final_mult,omitempty"`
	IsMegaWin   bool    `json:"is_mega_win,omitempty"` // 是否觸發 1000x Mega Win
	Reward      int     `json:"reward,omitempty"`
}

// LuckyDragonShotgunPayload T117 龍力散彈
type LuckyDragonShotgunPayload struct {
	Event       string   `json:"event"`        // "trigger" | "shotgun_fire" | "settle"
	TriggerID   string   `json:"trigger_id"`
	TriggerName string   `json:"trigger_name"`
	Direction   int      `json:"direction,omitempty"`   // 方向 0-7（8方向）
	HitTargets  []string `json:"hit_targets,omitempty"` // 命中目標
	TotalHits   int      `json:"total_hits,omitempty"`  // 總命中數
	TotalReward int      `json:"total_reward,omitempty"`
}

// LuckyRocketCannonPayload T118 火箭砲
type LuckyRocketCannonPayload struct {
	Event       string   `json:"event"`        // "trigger" | "rocket_launch" | "rocket_explode" | "settle"
	TriggerID   string   `json:"trigger_id"`
	TriggerName string   `json:"trigger_name"`
	RocketNo    int      `json:"rocket_no,omitempty"`    // 第幾枚火箭（1-3）
	ExplodeX    float64  `json:"explode_x,omitempty"`
	ExplodeY    float64  `json:"explode_y,omitempty"`
	HitTargets  []string `json:"hit_targets,omitempty"`
	TotalReward int      `json:"total_reward,omitempty"`
}

// LuckyDeepWhirlpoolPayload T119 深海漩渦
type LuckyDeepWhirlpoolPayload struct {
	Event       string  `json:"event"`        // "trigger" | "whirlpool_damage" | "settle"
	TriggerID   string  `json:"trigger_id"`
	TriggerName string  `json:"trigger_name"`
	HitCount    int     `json:"hit_count,omitempty"`
	TotalReward int     `json:"total_reward,omitempty"`
}

// LuckyVampireMultPayload T120 吸血鬼倍率
type LuckyVampireMultPayload struct {
	Event       string  `json:"event"`        // "trigger" | "absorb" | "mult_mode" | "mult_end" | "settle"
	TriggerID   string  `json:"trigger_id"`
	TriggerName string  `json:"trigger_name"`
	AbsorbCount int     `json:"absorb_count,omitempty"` // 已吸收次數
	CurrentMult float64 `json:"current_mult,omitempty"` // 當前倍率
	TimeLeft    float64 `json:"time_left,omitempty"`    // 倍率模式剩餘時間
	TotalReward int     `json:"total_reward,omitempty"`
}

// ── DAY-296 新增 Lucky 特殊魚 Payloads ───────────────────────

// LuckyMirrorFishPayload T121 鏡像魚
type LuckyMirrorFishPayload struct {
	Event       string `json:"event"`                  // "trigger" | "mirror_hit" | "perfect_mirror" | "perfect_end" | "settle" | "timeout"
	TriggerID   string `json:"trigger_id"`
	TriggerName string `json:"trigger_name"`
	ShotsLeft   int    `json:"shots_left,omitempty"`   // 剩餘鏡像次數
	HitCount    int    `json:"hit_count,omitempty"`    // 命中次數
	TotalReward int    `json:"total_reward,omitempty"`
}

// GoldenCoinInfo 黃金幣位置資訊
type GoldenCoinInfo struct {
	CoinID int     `json:"coin_id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

// LuckyGoldenRainPayload T122 黃金雨魚
type LuckyGoldenRainPayload struct {
	Event          string           `json:"event"`                    // "trigger" | "coin_collect" | "golden_harvest" | "harvest_end" | "settle"
	TriggerID      string           `json:"trigger_id"`
	TriggerName    string           `json:"trigger_name"`
	TotalCoins     int              `json:"total_coins,omitempty"`    // 生成的黃金幣總數
	CoinPositions  []GoldenCoinInfo `json:"coin_positions,omitempty"` // 黃金幣位置
	CollectorID    string           `json:"collector_id,omitempty"`   // 收集者 ID
	CoinID         int              `json:"coin_id,omitempty"`        // 被收集的幣 ID
	CollectedCoins int              `json:"collected_coins,omitempty"` // 已收集數
	TotalReward    int              `json:"total_reward,omitempty"`
}

// LuckyFreezeBombPayload T123 冰凍炸彈魚
type LuckyFreezeBombPayload struct {
	Event         string   `json:"event"`                    // "freeze_start" | "bomb_explode" | "perfect_freeze" | "perfect_end"
	TriggerID     string   `json:"trigger_id"`
	TriggerName   string   `json:"trigger_name"`
	BombX         float64  `json:"bomb_x,omitempty"`
	BombY         float64  `json:"bomb_y,omitempty"`
	FreezeRadius  float64  `json:"freeze_radius,omitempty"`
	FrozenTargets []string `json:"frozen_targets,omitempty"` // 被凍結的目標 instance_id
	Duration      float64  `json:"duration,omitempty"`       // 凍結持續時間
	HitCount      int      `json:"hit_count,omitempty"`      // 爆炸命中數
	TotalReward   int      `json:"total_reward,omitempty"`
}

// LuckyThunderStormPayload T124 雷暴魚
type LuckyThunderStormPayload struct {
	Event          string   `json:"event"`                    // "storm_start" | "lightning_strike" | "perfect_storm" | "perfect_end" | "storm_end"
	TriggerID      string   `json:"trigger_id"`
	TriggerName    string   `json:"trigger_name"`
	LightningCount int      `json:"lightning_count,omitempty"` // 閃電總數
	Duration       float64  `json:"duration,omitempty"`
	StrikeX        float64  `json:"strike_x,omitempty"`
	StrikeY        float64  `json:"strike_y,omitempty"`
	HitTargets     []string `json:"hit_targets,omitempty"`
	StrikeNo       int      `json:"strike_no,omitempty"`       // 第幾道閃電
	HitStrikes     int      `json:"hit_strikes,omitempty"`     // 命中的閃電數
	AccumMult      float64  `json:"accum_mult,omitempty"`
	TotalReward    int      `json:"total_reward,omitempty"`
}

// LuckyLuckyWheelPayload T125 大轉盤魚
type LuckyLuckyWheelPayload struct {
	Event       string  `json:"event"`                  // "trigger" | "spin_result"
	TriggerID   string  `json:"trigger_id"`
	TriggerName string  `json:"trigger_name"`
	PoolSize    int     `json:"pool_size,omitempty"`    // 大獎池大小
	SlotName    string  `json:"slot_name,omitempty"`    // 轉到的格子名稱
	SlotType    string  `json:"slot_type,omitempty"`    // "mult" | "aoe" | "jackpot"
	SlotMult    float64 `json:"slot_mult,omitempty"`    // 倍率
	Reward      int     `json:"reward,omitempty"`
}

// CollectGoldenCoinRequest T122 黃金雨魚：收集黃金幣請求
type CollectGoldenCoinRequest struct {
	CoinID int `json:"coin_id"`
}

// SetDisplayNameRequest 設定玩家顯示名稱請求
type SetDisplayNameRequest struct {
	Name string `json:"name"`
}

// ── DAY-301 新增 Lucky 特殊魚 Payloads ───────────────────────

// LuckyJackpotFishPayload T126 進階 Jackpot 魚
type LuckyJackpotFishPayload struct {
	Event       string  `json:"event"`                  // "trigger" | "jackpot_result" | "grand_boost" | "grand_boost_end"
	TriggerID   string  `json:"trigger_id"`
	TriggerName string  `json:"trigger_name"`
	MiniPool    int     `json:"mini_pool,omitempty"`
	MinorPool   int     `json:"minor_pool,omitempty"`
	MajorPool   int     `json:"major_pool,omitempty"`
	GrandPool   int     `json:"grand_pool,omitempty"`
	TierName    string  `json:"tier_name,omitempty"`    // "Mini" | "Minor" | "Major" | "Grand"
	TierIdx     int     `json:"tier_idx,omitempty"`     // 0-3
	Reward      int     `json:"reward,omitempty"`
	BoostMult   float64 `json:"boost_mult,omitempty"`
	BoostSecs   int     `json:"boost_secs,omitempty"`
}

// LuckyCoopFishPayload T127 全服合作魚
type LuckyCoopFishPayload struct {
	Event         string  `json:"event"`                   // "coop_start" | "coop_progress" | "coop_success" | "coop_timeout" | "coop_boost_end"
	TriggerID     string  `json:"trigger_id"`
	TriggerName   string  `json:"trigger_name"`
	TargetPoints  int     `json:"target_points,omitempty"`
	CurrentPoints int     `json:"current_points,omitempty"`
	TimeLeft      float64 `json:"time_left,omitempty"`
	BoostMult     float64 `json:"boost_mult,omitempty"`
	BoostSecs     int     `json:"boost_secs,omitempty"`
}

// LuckyTimeWarpPayload T128 時間扭曲魚
type LuckyTimeWarpPayload struct {
	Event       string  `json:"event"`                  // "warp_start" | "warp_end" | "time_collapse" | "collapse_end"
	TriggerID   string  `json:"trigger_id"`
	TriggerName string  `json:"trigger_name"`
	Duration    float64 `json:"duration,omitempty"`
	SpeedMult   float64 `json:"speed_mult,omitempty"`   // 目標移動速度倍率（0.3）
	DamageMult  float64 `json:"damage_mult,omitempty"`  // 傷害倍率（2.0）
	KillCount   int     `json:"kill_count,omitempty"`
	BoostMult   float64 `json:"boost_mult,omitempty"`
	BoostSecs   int     `json:"boost_secs,omitempty"`
}

// ── DAY-302 新增 Lucky 特殊魚 Payloads ───────────────────────

// LuckyChainMeteorPayload T129 連鎖隕石魚事件
type LuckyChainMeteorPayload struct {
	Event       string  `json:"event"`        // meteor_start / meteor_hit / meteor_miss / meteor_perfect / meteor_perfect_end
	PlayerID    string  `json:"player_id"`
	PlayerName  string  `json:"player_name"`
	MeteorIndex int     `json:"meteor_index"` // 第幾顆（1-5）
	AOERadius   float64 `json:"aoe_radius"`   // 當前 AOE 半徑
	HitCount    int     `json:"hit_count"`    // 命中目標數
	PerfectMult float64 `json:"perfect_mult"` // 完美加成倍率
	ExpiresAt   int64   `json:"expires_at"`   // 完美加成到期時間
}

// ── DAY-303 新增 Lucky 特殊魚 Payloads ───────────────────────

// LuckyCrashFishPayload T130 崩潰魚（Crash mechanic）
type LuckyCrashFishPayload struct {
	Event       string  `json:"event"`        // crash_start / mult_rise / harvest / crash / perfect_harvest / perfect_end
	PlayerID    string  `json:"player_id"`
	PlayerName  string  `json:"player_name"`
	CurrentMult float64 `json:"current_mult"` // 當前倍率
	CrashIn     float64 `json:"crash_in"`     // 距離崩潰的秒數（僅 crash_start）
	TimeLeft    float64 `json:"time_left"`    // 剩餘時間
	Reward      int     `json:"reward"`       // 收割獎勵
	BoostMult   float64 `json:"boost_mult"`   // 完美收割全服加成倍率
	BoostSecs   int     `json:"boost_secs"`   // 完美收割加成秒數
}

// ── DAY-304 新增 Lucky 特殊魚 Payloads ───────────────────────

// LuckyElectricEelPayload T131 電鰻魚（持續放電連鎖）
type LuckyElectricEelPayload struct {
	Event      string  `json:"event"`       // eel_start / eel_shock / eel_end / eel_super / eel_super_end
	PlayerID   string  `json:"player_id"`
	PlayerName string  `json:"player_name"`
	Duration   float64 `json:"duration,omitempty"`   // 持續時間（秒）
	HitCount   int     `json:"hit_count,omitempty"`  // 本次電擊命中數
	ShockCount int     `json:"shock_count,omitempty"` // 累積電擊次數
	TimeLeft   float64 `json:"time_left,omitempty"`  // 剩餘時間
	BoostMult  float64 `json:"boost_mult,omitempty"` // 超級放電加成倍率
	BoostSec   int     `json:"boost_sec,omitempty"`  // 加成秒數
}

// LuckyAnglerFishPayload T132 巨型安康魚（誘餌+電擊爆炸）
type LuckyAnglerFishPayload struct {
	Event      string  `json:"event"`        // lure_start / explosion / perfect / perfect_end / lure_end
	PlayerID   string  `json:"player_id"`
	PlayerName string  `json:"player_name"`
	LureSec    int     `json:"lure_sec,omitempty"`    // 誘餌持續秒數
	DamageMult float64 `json:"damage_mult,omitempty"` // 誘餌期間傷害倍率
	HitCount   int     `json:"hit_count,omitempty"`   // 電擊命中數
	BoostMult  float64 `json:"boost_mult,omitempty"`  // 完美誘捕加成倍率
	BoostSec   int     `json:"boost_sec,omitempty"`   // 加成秒數
}

// LuckyBlackHolePayload T133 黑洞魚（吸引+坍縮）
type LuckyBlackHolePayload struct {
	Event      string  `json:"event"`        // black_hole_start / collapse / singularity / singularity_end / black_hole_end
	PlayerID   string  `json:"player_id"`
	PlayerName string  `json:"player_name"`
	Duration   float64 `json:"duration,omitempty"`   // 黑洞持續秒數
	HitCount   int     `json:"hit_count,omitempty"`  // 坍縮命中數
	BoostMult  float64 `json:"boost_mult,omitempty"` // 奇點爆發加成倍率
	BoostSec   int     `json:"boost_sec,omitempty"`  // 加成秒數
}

// LuckyBountyHunterPayload T134 賞金獵人魚（賞金目標系統）
type LuckyBountyHunterPayload struct {
	Event         string   `json:"event"`                    // bounty_start / bounty_kill / bounty_perfect / bounty_perfect_end / bounty_timeout
	PlayerID      string   `json:"player_id"`
	PlayerName    string   `json:"player_name"`
	BountyTargets []string `json:"bounty_targets,omitempty"` // 賞金目標 instance_id 列表
	TotalBounty   int      `json:"total_bounty,omitempty"`   // 賞金目標總數
	Duration      float64  `json:"duration,omitempty"`       // 任務持續秒數
	KillerID      string   `json:"killer_id,omitempty"`      // 擊破者 ID
	KillerName    string   `json:"killer_name,omitempty"`    // 擊破者名稱
	KillCount     int      `json:"kill_count,omitempty"`     // 已擊破賞金目標數
	BoostMult     float64  `json:"boost_mult,omitempty"`     // 完美賞金加成倍率
	BoostSec      int      `json:"boost_sec,omitempty"`      // 加成秒數
}

// LuckyTsunamiPayload T135 海嘯魚（三波衝擊）
type LuckyTsunamiPayload struct {
	Event         string  `json:"event"`                     // tsunami_warning / wave_hit / tsunami_perfect / tsunami_perfect_end / tsunami_end
	PlayerID      string  `json:"player_id"`
	PlayerName    string  `json:"player_name"`
	WaveCount     int     `json:"wave_count,omitempty"`      // 總波數
	WaveNum       int     `json:"wave_num,omitempty"`        // 當前波次（1-3）
	HitCount      int     `json:"hit_count,omitempty"`       // 本波命中數
	DamagePct     float64 `json:"damage_pct,omitempty"`      // 本波傷害百分比
	TotalHitCount int     `json:"total_hit_count,omitempty"` // 三波累積命中數
	BoostMult     float64 `json:"boost_mult,omitempty"`      // 完美海嘯加成倍率
	BoostSec      int     `json:"boost_sec,omitempty"`       // 加成秒數
}

// ── DAY-305 新增 Lucky 特殊魚 Payloads ───────────────────────

// LuckyDragonWrathV2Payload T136 龍怒蓄積魚 v2
type LuckyDragonWrathV2Payload struct {
	Event       string  `json:"event"`        // wrath_start / wrath_explode / wrath_perfect / wrath_perfect_end / wrath_end
	PlayerID    string  `json:"player_id"`
	PlayerName  string  `json:"player_name"`
	Duration    float64 `json:"duration,omitempty"`
	MaxWrath    int     `json:"max_wrath,omitempty"`
	WrathValue  int     `json:"wrath_value,omitempty"`
	MeteorCount int     `json:"meteor_count,omitempty"`
	HitCount    int     `json:"hit_count,omitempty"`
	BoostMult   float64 `json:"boost_mult,omitempty"`
	BoostSec    int     `json:"boost_sec,omitempty"`
}

// LuckyHumpbackWhalePayload T137 座頭鯨魚（鯨歌共鳴）
type LuckyHumpbackWhalePayload struct {
	Event         string  `json:"event"`                     // song_start / song_wave / song_perfect / song_perfect_end / song_end
	PlayerID      string  `json:"player_id"`
	PlayerName    string  `json:"player_name"`
	WaveCount     int     `json:"wave_count,omitempty"`
	WaveNum       int     `json:"wave_num,omitempty"`
	HitCount      int     `json:"hit_count,omitempty"`
	DamagePct     float64 `json:"damage_pct,omitempty"`
	TotalHitCount int     `json:"total_hit_count,omitempty"`
	BoostMult     float64 `json:"boost_mult,omitempty"`
	BoostSec      int     `json:"boost_sec,omitempty"`
}

// LuckyLegendDragonPayload T138 傳說龍魚（龍息噴火）
type LuckyLegendDragonPayload struct {
	Event          string  `json:"event"`                      // dragon_appear / dragon_breath / dragon_rage / dragon_rage_end / dragon_leave
	PlayerID       string  `json:"player_id"`
	PlayerName     string  `json:"player_name"`
	Duration       float64 `json:"duration,omitempty"`
	BreathNum      int     `json:"breath_num,omitempty"`
	HitCount       int     `json:"hit_count,omitempty"`
	PerfectBreaths int     `json:"perfect_breaths,omitempty"`
	BoostMult      float64 `json:"boost_mult,omitempty"`
	BoostSec       int     `json:"boost_sec,omitempty"`
}

// LuckyGuildWarPayload T139 公會戰魚（全服積分）
type LuckyGuildWarPayload struct {
	Event         string  `json:"event"`                    // war_start / war_progress / war_victory / war_victory_end / war_timeout
	PlayerID      string  `json:"player_id"`
	PlayerName    string  `json:"player_name"`
	TargetPoints  int     `json:"target_points,omitempty"`
	CurrentPoints int     `json:"current_points,omitempty"`
	Duration      float64 `json:"duration,omitempty"`
	KillerID      string  `json:"killer_id,omitempty"`
	KillerName    string  `json:"killer_name,omitempty"`
	BoostMult     float64 `json:"boost_mult,omitempty"`
	BoostSec      int     `json:"boost_sec,omitempty"`
}

// LuckyQualityFishPayload T140 品質魚（品質鑑定）
type LuckyQualityFishPayload struct {
	Event      string  `json:"event"`        // quality_result / legendary_boost / legendary_end
	PlayerID   string  `json:"player_id"`
	PlayerName string  `json:"player_name"`
	TierName   string  `json:"tier_name,omitempty"`   // Common / Rare / Epic / Legendary
	TierMult   float64 `json:"tier_mult,omitempty"`
	Reward     int     `json:"reward,omitempty"`
	BoostMult  float64 `json:"boost_mult,omitempty"`
	BoostSec   int     `json:"boost_sec,omitempty"`
}

// ── DAY-306 新增 Lucky 特殊魚 Payloads ───────────────────────

// LuckyTornadoPayload T141 龍捲風魚（橫掃全場）
type LuckyTornadoPayload struct {
	Event      string  `json:"event"`                  // tornado_start / tornado_sweep / tornado_end / tornado_perfect / tornado_perfect_end
	PlayerID   string  `json:"player_id"`
	PlayerName string  `json:"player_name"`
	Duration   float64 `json:"duration,omitempty"`
	WaveNum    int     `json:"wave_num,omitempty"`
	HitCount   int     `json:"hit_count,omitempty"`
	KillCount  int     `json:"kill_count,omitempty"`
	BoostMult  float64 `json:"boost_mult,omitempty"`
	BoostSec   int     `json:"boost_sec,omitempty"`
}

// LuckyEarthquakePayload T142 地震魚（三波同心圓）
type LuckyEarthquakePayload struct {
	Event         string  `json:"event"`                     // quake_warning / quake_wave / quake_end / quake_perfect / quake_perfect_end
	PlayerID      string  `json:"player_id"`
	PlayerName    string  `json:"player_name"`
	WaveCount     int     `json:"wave_count,omitempty"`
	WaveNum       int     `json:"wave_num,omitempty"`
	DamagePct     float64 `json:"damage_pct,omitempty"`
	HitCount      int     `json:"hit_count,omitempty"`
	TotalHitCount int     `json:"total_hit_count,omitempty"`
	BoostMult     float64 `json:"boost_mult,omitempty"`
	BoostSec      int     `json:"boost_sec,omitempty"`
}

// LuckyVolcanoPayload T143 火山魚（熔岩彈雨）
type LuckyVolcanoPayload struct {
	Event      string  `json:"event"`                  // volcano_erupt / lava_bomb / volcano_end / volcano_perfect / volcano_perfect_end
	PlayerID   string  `json:"player_id"`
	PlayerName string  `json:"player_name"`
	BombCount  int     `json:"bomb_count,omitempty"`
	BombNum    int     `json:"bomb_num,omitempty"`
	BombX      float64 `json:"bomb_x,omitempty"`
	BombY      float64 `json:"bomb_y,omitempty"`
	HitCount   int     `json:"hit_count,omitempty"`
	HitBombs   int     `json:"hit_bombs,omitempty"`
	BoostMult  float64 `json:"boost_mult,omitempty"`
	BoostSec   int     `json:"boost_sec,omitempty"`
}

// LuckyCosmicRayPayload T144 星際魚（8方向光束）
type LuckyCosmicRayPayload struct {
	Event         string  `json:"event"`                     // cosmic_start / cosmic_ray / cosmic_end / cosmic_perfect / cosmic_perfect_end
	PlayerID      string  `json:"player_id"`
	PlayerName    string  `json:"player_name"`
	RayCount      int     `json:"ray_count,omitempty"`
	Direction     int     `json:"direction,omitempty"`       // 0-7（8方向）
	HitCount      int     `json:"hit_count,omitempty"`
	TotalHitCount int     `json:"total_hit_count,omitempty"`
	BoostMult     float64 `json:"boost_mult,omitempty"`
	BoostSec      int     `json:"boost_sec,omitempty"`
}

// LuckyDivineDragonPayload T145 神龍魚（神龍降臨）
type LuckyDivineDragonPayload struct {
	Event        string  `json:"event"`                   // dragon_descend / dragon_claw / dragon_leave / dragon_perfect / dragon_perfect_end
	PlayerID     string  `json:"player_id"`
	PlayerName   string  `json:"player_name"`
	Duration     float64 `json:"duration,omitempty"`
	ClawCount    int     `json:"claw_count,omitempty"`
	ClawNum      int     `json:"claw_num,omitempty"`
	HitCount     int     `json:"hit_count,omitempty"`
	PerfectClaws int     `json:"perfect_claws,omitempty"`
	BoostMult    float64 `json:"boost_mult,omitempty"`
	BoostSec     int     `json:"boost_sec,omitempty"`
}

// ── DAY-307 新增 Lucky 特殊魚 Payloads ───────────────────────

// LuckyQuantumPayload T146 量子魚（量子觀測坍縮）
type LuckyQuantumPayload struct {
	Event         string  `json:"event"`                      // quantum_observe / quantum_result / quantum_collapse / quantum_collapse_end
	PlayerID      string  `json:"player_id"`
	PlayerName    string  `json:"player_name"`
	ObservedCount int     `json:"observed_count,omitempty"`   // 被觀測到的目標數
	BoostMult     float64 `json:"boost_mult,omitempty"`
	BoostSec      int     `json:"boost_sec,omitempty"`
}

// LuckySupernovaPayload T147 超新星魚（全場爆炸+倍率加成）
type LuckySupernovaPayload struct {
	Event      string  `json:"event"`                   // supernova_explode / supernova_boost / supernova_end / supernova_perfect / supernova_perfect_end
	PlayerID   string  `json:"player_id"`
	PlayerName string  `json:"player_name"`
	HitCount   int     `json:"hit_count,omitempty"`
	MultBoost  float64 `json:"mult_boost,omitempty"`    // 5秒倍率加成
	BoostSec   int     `json:"boost_sec,omitempty"`
	BoostMult  float64 `json:"boost_mult,omitempty"`
}

// LuckyInfinitePayload T148 無限魚（無限累積倍率）
type LuckyInfinitePayload struct {
	Event      string  `json:"event"`                   // infinite_start / infinite_kill / infinite_end / infinite_perfect / infinite_perfect_end
	PlayerID   string  `json:"player_id"`
	PlayerName string  `json:"player_name"`
	Duration   float64 `json:"duration,omitempty"`
	AccumMult  float64 `json:"accum_mult,omitempty"`    // 累積倍率
	KillCount  int     `json:"kill_count,omitempty"`
	BoostMult  float64 `json:"boost_mult,omitempty"`
	BoostSec   int     `json:"boost_sec,omitempty"`
}

// LuckyGenesisPayload T149 創世魚（全場審判）
type LuckyGenesisPayload struct {
	Event       string  `json:"event"`                    // genesis_descend / genesis_judgment / genesis_blessing / genesis_blessing_end
	PlayerID    string  `json:"player_id"`
	PlayerName  string  `json:"player_name"`
	KillCount   int     `json:"kill_count,omitempty"`
	TotalReward int     `json:"total_reward,omitempty"`
	MultBoost   float64 `json:"mult_boost,omitempty"`     // 每個目標的倍率加成
	BoostMult   float64 `json:"boost_mult,omitempty"`
	BoostSec    int     `json:"boost_sec,omitempty"`
}

// LuckyRebirthPayload T150 重生魚（死亡目標復活再擊破）
type LuckyRebirthPayload struct {
	Event        string  `json:"event"`                    // rebirth_start / rebirth_kill / rebirth_end / rebirth_perfect / rebirth_perfect_end
	PlayerID     string  `json:"player_id"`
	PlayerName   string  `json:"player_name"`
	Duration     float64 `json:"duration,omitempty"`
	RebirthKills int     `json:"rebirth_kills,omitempty"`  // 重生後被擊破的目標數
	BoostMult    float64 `json:"boost_mult,omitempty"`
	BoostSec     int     `json:"boost_sec,omitempty"`
}
