// Package protocol — WebSocket 訊息協定定義
// server-core-agent + protocol-sync-agent 負責維護
package protocol

// ── Client → Server ──────────────────────────────────────────
const (
	MsgAttack       = "attack"
	MsgLock         = "lock"
	MsgAutoToggle   = "auto_toggle"
	MsgBetChange    = "bet_change"
	MsgBonusClick   = "bonus_click"
	MsgPing         = "ping"
	MsgTriggerBoss  = "trigger_boss"
	MsgTriggerBonus = "trigger_bonus"
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
