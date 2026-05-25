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
