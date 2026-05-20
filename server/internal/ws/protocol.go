// Package ws 定義 WebSocket 通訊協定
package ws

// MessageType 訊息類型
type MessageType string

// Client → Server
const (
	MsgAttack     MessageType = "attack"
	MsgLock       MessageType = "lock"
	MsgAutoToggle MessageType = "auto_toggle"
	MsgBetChange  MessageType = "bet_change"
	MsgBonusClick MessageType = "bonus_click"
	MsgPing       MessageType = "ping"
	// Prototype 展示用
	MsgTriggerBoss     MessageType = "trigger_boss"
	MsgTriggerBonus    MessageType = "trigger_bonus"
	MsgSetDisplayName  MessageType = "set_display_name"  // 設定顯示名稱（DAY-021）
	MsgClaimMission    MessageType = "claim_mission"     // 領取任務獎勵（DAY-037）
	MsgGetMissions     MessageType = "get_missions"      // 查詢任務列表（DAY-037）
	MsgClientPerf      MessageType = "client_perf"       // Client 端效能數據上報（DAY-045）
)

// Server → Client
const (
	MsgGameState    MessageType = "game_state"
	MsgTargetSpawn  MessageType = "target_spawn"
	MsgTargetUpdate MessageType = "target_update"
	MsgTargetKill   MessageType = "target_kill"
	MsgAttackResult MessageType = "attack_result"
	MsgReward       MessageType = "reward"
	MsgBossEvent    MessageType = "boss_event"
	MsgBonusEvent   MessageType = "bonus_event"
	MsgPlayerUpdate MessageType = "player_update"
	MsgLeaderboard  MessageType = "leaderboard"
	MsgAchievement  MessageType = "achievement"
	MsgComboEvent   MessageType = "combo_event"      // 連擊事件（DAY-022）
	MsgSpectatorJoin  MessageType = "spectator_join"  // 觀戰者加入通知（DAY-023）
	MsgSpectatorLeave MessageType = "spectator_leave" // 觀戰者離開通知（DAY-055）
	MsgMissionUpdate MessageType = "mission_update"  // 任務進度更新（DAY-037）
	MsgMissionComplete MessageType = "mission_complete" // 任務完成通知（DAY-037）
	MsgJackpotUpdate   MessageType = "jackpot_update"   // Jackpot 池更新（DAY-048）
	MsgJackpotWin      MessageType = "jackpot_win"      // Jackpot 中獎通知（DAY-048）
	MsgDailyBonus        MessageType = "daily_bonus"        // 每日登入獎勵（DAY-065）
	MsgTournamentUpdate  MessageType = "tournament_update"  // 週賽排名更新（DAY-066）
	MsgTournamentResult  MessageType = "tournament_result"  // 週賽結算通知（DAY-066）
	MsgError        MessageType = "error"
	MsgPong         MessageType = "pong"
)

// Message 通用訊息結構
type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

// ---- Client → Server Payloads ----

// AttackPayload 攻擊請求
type AttackPayload struct {
	TargetID string  `json:"target_id"` // 空字串=自由攻擊
	ClickX   float64 `json:"click_x"`
	ClickY   float64 `json:"click_y"`
}

// LockPayload 鎖定目標
type LockPayload struct {
	TargetID string `json:"target_id"` // 空字串=解除鎖定
}

// BetChangePayload 切換投注
type BetChangePayload struct {
	BetLevel int `json:"bet_level"`
}

// BonusClickPayload Bonus 點擊
type BonusClickPayload struct {
	TargetID string  `json:"target_id"`
	ClickX   float64 `json:"click_x"`
	ClickY   float64 `json:"click_y"`
}

// ---- Server → Client Payloads ----

// GameStatePayload 遊戲狀態
type GameStatePayload struct {
	State     string `json:"state"`
	Timestamp int64  `json:"timestamp"`
}

// TargetSpawnPayload 目標生成
type TargetSpawnPayload struct {
	InstanceID string  `json:"instance_id"`
	DefID      string  `json:"def_id"`
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	HP         int     `json:"hp"`
	MaxHP      int     `json:"max_hp"`
	Speed      float64 `json:"speed"`
	Lifetime   float64 `json:"lifetime"`
	Behavior   string  `json:"behavior"`
	Multiplier float64 `json:"multiplier"` // 目標倍率（Client 顯示用）
}

// TargetUpdatePayload 目標狀態更新
type TargetUpdatePayload struct {
	InstanceID string  `json:"instance_id"`
	HP         int     `json:"hp"`
	MaxHP      int     `json:"max_hp"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Phase      int     `json:"phase"`      // BOSS 用
	IsFleeing  bool    `json:"is_fleeing"` // 寶箱怪用
}

// TargetKillPayload 目標擊破
type TargetKillPayload struct {
	InstanceID  string  `json:"instance_id"`
	DefID       string  `json:"def_id"`
	Multiplier  float64 `json:"multiplier"`
	Reward      int     `json:"reward"`
	LaborGain   int     `json:"labor_gain"`
	KillerID    string  `json:"killer_id"`
}

// AttackResultPayload 攻擊結果
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

// RewardPayload 獎勵發放
type RewardPayload struct {
	Source     string  `json:"source"` // "target", "boss", "bonus"
	Amount     int     `json:"amount"`
	Multiplier float64 `json:"multiplier"`
	NewBalance int     `json:"new_balance"`
}

// BossEventPayload BOSS 事件
type BossEventPayload struct {
	Event      string  `json:"event"` // "warning", "spawn", "phase_change", "kill"
	InstanceID string  `json:"instance_id"`
	Phase      int     `json:"phase"`
	HP         int     `json:"hp"`
	MaxHP      int     `json:"max_hp"`
	Reward     int     `json:"reward"`
	Multiplier float64 `json:"multiplier"`
}

// BonusEventPayload Bonus 事件
type BonusEventPayload struct {
	Event      string  `json:"event"` // "ready", "start", "tick", "end"
	TimeLeft   float64 `json:"time_left"`
	Score      int     `json:"score"`
	Multiplier float64 `json:"multiplier"`
	Reward     int     `json:"reward"`
}

// ErrorPayload 錯誤訊息
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// LeaderboardEntry 排行榜單筆記錄
type LeaderboardEntry struct {
	Rank        int    `json:"rank"`
	PlayerID    string `json:"player_id"`
	DisplayName string `json:"display_name"`
	Score       int    `json:"score"`       // 本局累積獎勵
	MaxCoins    int    `json:"max_coins"`   // 歷史最高金幣
	KillCount   int    `json:"kill_count"`  // 本局擊破數
	IsSelf      bool   `json:"is_self"`     // 是否為自己（Client 端標記用）
}

// LeaderboardPayload 排行榜廣播
type LeaderboardPayload struct {
	Entries   []LeaderboardEntry `json:"entries"`
	Timestamp int64              `json:"timestamp"`
}

// AchievementPayload 成就解鎖通知
type AchievementPayload struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	UnlockedAt  int64  `json:"unlocked_at"` // Unix milliseconds
}

// ComboEventPayload 連擊事件（DAY-022）
type ComboEventPayload struct {
	ComboCount  int     `json:"combo_count"`  // 當前連擊數（2+）
	LaborBonus  float64 `json:"labor_bonus"`  // 勞動值加成係數（0.1/0.2/0.3）
	PlayerID    string  `json:"player_id"`    // 觸發連擊的玩家
}

// ---- 任務系統 Payloads（DAY-037）----

// MissionPayload 單一任務狀態
type MissionPayload struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Target      int    `json:"target"`
	Current     int    `json:"current"`
	Completed   bool   `json:"completed"`
	RewardClaimed bool `json:"reward_claimed"`
	Reward      int    `json:"reward"`
}

// MissionUpdatePayload 任務進度更新廣播
type MissionUpdatePayload struct {
	PlayerID      string           `json:"player_id"`
	Missions      []MissionPayload `json:"missions"`
	ResetAt       int64            `json:"reset_at"`        // Unix ms，下次重置時間（UTC+8 00:00）
	ResetTimezone string           `json:"reset_timezone"`  // 重置時區說明（"UTC+8"）
}

// MissionCompletePayload 任務完成通知
type MissionCompletePayload struct {
	MissionID   string `json:"mission_id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Reward      int    `json:"reward"`
}

// ClaimMissionPayload 領取任務獎勵請求
type ClaimMissionPayload struct {
	MissionID string `json:"mission_id"`
}

// ---- Client 端效能數據（DAY-045）----

// ClientPerfPayload Client 端效能數據上報
// Client 每 30 秒發送一次，Server 記錄並暴露到 /metrics
type ClientPerfPayload struct {
	FPS        float64 `json:"fps"`         // 當前平均 FPS
	MemoryMB   float64 `json:"memory_mb"`   // 靜態記憶體使用（MB）
	DrawCalls  int     `json:"draw_calls"`  // 每幀 Draw Call 數
	NodeCount  int     `json:"node_count"`  // 場景節點數
	PingMs     int     `json:"ping_ms"`     // Client 端測量的 ping 延遲（ms）
	Quality    string  `json:"quality"`     // 效能等級（HIGH/MEDIUM/LOW）
	Timestamp  int64   `json:"timestamp"`   // Unix ms
}

// ---- Progressive Jackpot（DAY-048）----

// JackpotUpdatePayload Jackpot 池更新廣播（每 5 秒）
type JackpotUpdatePayload struct {
	Mini    int                  `json:"mini"`    // Mini Jackpot 當前金額
	Major   int                  `json:"major"`   // Major Jackpot 當前金額
	Grand   int                  `json:"grand"`   // Grand Jackpot 當前金額
}

// JackpotWinPayload Jackpot 中獎通知
type JackpotWinPayload struct {
	Level       string `json:"level"`        // "mini" / "major" / "grand"
	Amount      int    `json:"amount"`       // 中獎金額
	WinnerID    string `json:"winner_id"`    // 中獎玩家 ID
	WinnerName  string `json:"winner_name"`  // 中獎玩家顯示名稱
	NewBalance  int    `json:"new_balance"`  // 中獎後餘額（只對中獎玩家有意義）
}

// ---- 每日登入獎勵（DAY-065）----

// DailyBonusPayload 每日登入獎勵通知
// Server 在玩家加入時檢查是否是新的一天，計算連續天數，發放獎勵
type DailyBonusPayload struct {
	Streak      int    `json:"streak"`        // 連續登入天數（1=首次/重置，2=連續2天...）
	Reward      int    `json:"reward"`        // 本次獎勵金幣數
	NewBalance  int    `json:"new_balance"`   // 領取後餘額
	IsNewStreak bool   `json:"is_new_streak"` // 是否是今天第一次登入（false=今天已領過）
	MaxStreak   int    `json:"max_streak"`    // 最高連續天數（用於顯示里程碑）
}

// ---- 週賽系統（DAY-066）----

// TournamentRankEntry 週賽排名單筆記錄
type TournamentRankEntry struct {
	Rank        int    `json:"rank"`
	PlayerID    string `json:"player_id"`
	DisplayName string `json:"display_name"`
	Points      int    `json:"points"`
	Prize       int    `json:"prize"`        // 獎勵金幣（前三名才有）
	PrizeLabel  string `json:"prize_label"`  // 獎勵標籤（"🥇 週賽冠軍" 等）
	IsSelf      bool   `json:"is_self"`      // 是否為自己（Client 端標記用）
}

// TournamentUpdatePayload 週賽排名更新廣播（每 30 秒）
type TournamentUpdatePayload struct {
	WeekStart    int64                 `json:"week_start"`    // Unix ms
	WeekEnd      int64                 `json:"week_end"`      // Unix ms
	SecondsLeft  int64                 `json:"seconds_left"`  // 距離結束秒數
	Rankings     []TournamentRankEntry `json:"rankings"`      // 前 10 名
	TotalPlayers int                   `json:"total_players"` // 本週參賽人數
	PlayerRank   int                   `json:"player_rank"`   // 接收者的排名（0=未上榜）
	PlayerPoints int                   `json:"player_points"` // 接收者的積分
}

// TournamentResultPayload 週賽結算通知（週結束時廣播）
type TournamentResultPayload struct {
	WeekStart int64                 `json:"week_start"`
	WeekEnd   int64                 `json:"week_end"`
	Rankings  []TournamentRankEntry `json:"rankings"`
	Prize     int                   `json:"prize"`       // 接收者獲得的獎勵（0=未獲獎）
	PrizeLabel string               `json:"prize_label"` // 獎勵標籤
}
