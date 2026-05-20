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
	MsgUpgradeWeapon   MessageType = "upgrade_weapon"    // 武器升級（DAY-067）
	MsgSetTitle        MessageType = "set_title"          // 設定顯示稱號（DAY-068）
	MsgBuySkin         MessageType = "buy_skin"           // 購買砲台外觀（DAY-071）
	MsgEquipSkin       MessageType = "equip_skin"         // 裝備砲台外觀（DAY-071）
	MsgClaimSeasonLevel MessageType = "claim_season_level" // 領取賽季等級獎勵（DAY-072）
	MsgSendFriendRequest  MessageType = "send_friend_request"  // 發送好友請求（DAY-073）
	MsgAcceptFriendRequest MessageType = "accept_friend_request" // 接受好友請求（DAY-073）
	MsgRejectFriendRequest MessageType = "reject_friend_request" // 拒絕好友請求（DAY-073）
	MsgRemoveFriend       MessageType = "remove_friend"         // 移除好友（DAY-073）
	MsgGetFriendList      MessageType = "get_friend_list"       // 查詢好友列表（DAY-073）
	// 公會系統（DAY-074）
	MsgCreateGuild   MessageType = "create_guild"   // 建立公會
	MsgJoinGuild     MessageType = "join_guild"     // 加入公會
	MsgLeaveGuild    MessageType = "leave_guild"    // 退出公會
	MsgKickGuildMember MessageType = "kick_guild_member" // 踢出成員
	MsgPromoteGuildMember MessageType = "promote_guild_member" // 升職成員
	MsgGetGuildInfo  MessageType = "get_guild_info" // 查詢公會資訊
	MsgGetGuildList  MessageType = "get_guild_list" // 查詢公會列表
	MsgGuildChat     MessageType = "guild_chat"     // 公會聊天（DAY-075）
	// 公會戰系統（DAY-076）
	MsgGetGuildWarStatus MessageType = "get_guild_war_status" // 查詢公會戰狀態
	// 每日 BOSS 挑戰（DAY-077）
	MsgGetDailyBoss      MessageType = "get_daily_boss"       // 查詢每日 BOSS 狀態
	MsgDailyBossAttack   MessageType = "daily_boss_attack"    // 對每日 BOSS 攻擊
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
	MsgSkinUpdate        MessageType = "skin_update"        // 砲台外觀更新（DAY-071）
	MsgSeasonUpdate      MessageType = "season_update"      // 賽季通行證更新（DAY-072）
	MsgSeasonLevelUp     MessageType = "season_level_up"    // 賽季等級升級通知（DAY-072）
	MsgFriendList        MessageType = "friend_list"        // 好友列表（DAY-073）
	MsgFriendRequest     MessageType = "friend_request"     // 好友請求通知（DAY-073）
	MsgFriendUpdate      MessageType = "friend_update"      // 好友狀態更新（DAY-073）
	// 公會系統（DAY-074）
	MsgGuildUpdate       MessageType = "guild_update"       // 公會資訊更新
	MsgGuildList         MessageType = "guild_list"         // 公會列表
	MsgGuildTaskComplete MessageType = "guild_task_complete" // 公會任務完成通知
	MsgGuildError        MessageType = "guild_error"        // 公會操作錯誤
	MsgGuildMessage      MessageType = "guild_message"      // 公會聊天訊息（DAY-075）
	// 公會戰系統（DAY-076）
	MsgGuildWarUpdate    MessageType = "guild_war_update"    // 公會戰排名更新
	MsgGuildWarResult    MessageType = "guild_war_result"    // 公會戰結算通知
	// 每日 BOSS 挑戰（DAY-077）
	MsgDailyBossUpdate   MessageType = "daily_boss_update"   // 每日 BOSS 狀態更新
	MsgDailyBossDefeated MessageType = "daily_boss_defeated" // 每日 BOSS 擊殺通知
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
	// 品質等級（DAY-070）
	Quality      string `json:"quality"`       // "normal"/"rare"/"epic"/"legendary"
	QualityColor string `json:"quality_color"` // 光暈顏色（hex，空字串=無光暈）
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
	Quality     string  `json:"quality"`      // 品質等級（DAY-070）
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
	// 稱號（DAY-068）
	TitleID    string `json:"title_id"`
	TitleName  string `json:"title_name"`
	TitleIcon  string `json:"title_icon"`
	TitleColor string `json:"title_color"`
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

// ---- 武器升級系統（DAY-067）----

// UpgradeWeaponPayload 武器升級請求（Client → Server）
type UpgradeWeaponPayload struct {
	WeaponLevel int `json:"weapon_level"` // 目標武器等級（1/2/3）
}

// ---- 稱號系統（DAY-068）----

// MsgTitleUnlocked 稱號解鎖通知（Server → Client）
// 當玩家解鎖新稱號時廣播
const MsgTitleUnlocked MessageType = "title_unlocked"

// TitleUnlockedPayload 稱號解鎖通知
type TitleUnlockedPayload struct {
	TitleID    string `json:"title_id"`
	TitleName  string `json:"title_name"`
	TitleIcon  string `json:"title_icon"`
	TitleColor string `json:"title_color"`
	Description string `json:"description"`
}

// SetTitlePayload 設定顯示稱號請求
type SetTitlePayload struct {
	TitleID string `json:"title_id"`
}

// ---- 砲台外觀系統（DAY-071）----

// SkinDef 砲台外觀定義
type SkinDef struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`       // 0 = 免費
	CannonColor string `json:"cannon_color"` // 砲台顏色（hex）
	BulletColor string `json:"bullet_color"` // 投射物顏色（hex）
	GlowColor   string `json:"glow_color"`   // 光暈顏色（hex，空=無光暈）
	Icon        string `json:"icon"`
}

// AvailableSkins 可用外觀列表
var AvailableSkins = []SkinDef{
	{
		ID: "default", Name: "標準砲台", Description: "預設外觀",
		Price: 0, CannonColor: "", BulletColor: "", GlowColor: "", Icon: "🔫",
	},
	{
		ID: "golden", Name: "黃金砲台", Description: "閃耀的黃金外觀，彰顯財富",
		Price: 5000, CannonColor: "#FFD700", BulletColor: "#FFA500", GlowColor: "#FFD700", Icon: "✨",
	},
	{
		ID: "rainbow", Name: "彩虹砲台", Description: "七彩光芒，傳說等級外觀",
		Price: 20000, CannonColor: "#FF69B4", BulletColor: "#00FFFF", GlowColor: "#FF00FF", Icon: "🌈",
	},
	{
		ID: "sakura", Name: "櫻花砲台", Description: "吉伊卡哇限定，粉嫩可愛",
		Price: 8000, CannonColor: "#FFB7C5", BulletColor: "#FF69B4", GlowColor: "#FFB7C5", Icon: "🌸",
	},
}

// BuySkinPayload 購買外觀請求（Client → Server）
type BuySkinPayload struct {
	SkinID string `json:"skin_id"`
}

// EquipSkinPayload 裝備外觀請求（Client → Server）
type EquipSkinPayload struct {
	SkinID string `json:"skin_id"`
}

// SkinUpdatePayload 外觀更新通知（Server → Client）
type SkinUpdatePayload struct {
	PlayerID    string   `json:"player_id"`
	EquippedSkin string  `json:"equipped_skin"`
	OwnedSkins  []string `json:"owned_skins"`
	NewBalance  int      `json:"new_balance"` // 購買後的新餘額（0=裝備操作）
}

// ---- 賽季通行證系統（DAY-072）----

// ClaimSeasonLevelPayload 領取賽季等級獎勵請求（Client → Server）
type ClaimSeasonLevelPayload struct {
	Level int `json:"level"` // 要領取的等級（1-10）
}

// SeasonUpdatePayload 賽季通行證更新廣播（Server → Client）
// 在玩家加入、積分增加、領取獎勵時發送
type SeasonUpdatePayload struct {
	PlayerID     string              `json:"player_id"`
	SeasonPoints int                 `json:"season_points"`
	CurrentLevel int                 `json:"current_level"`
	NextLevel    int                 `json:"next_level"`    // 0 = 已滿級
	PointsToNext int                 `json:"points_to_next"`
	Progress     float64             `json:"progress"`      // 0.0-1.0
	Levels       []SeasonLevelStatus `json:"levels"`
}

// SeasonLevelStatus 賽季等級狀態（用於 SeasonUpdatePayload）
type SeasonLevelStatus struct {
	Level        int    `json:"level"`
	PointsNeeded int    `json:"points_needed"`
	CoinReward   int    `json:"coin_reward"`
	SpecialType  string `json:"special_type"`  // "" / "skin" / "title"
	SpecialID    string `json:"special_id"`
	SpecialName  string `json:"special_name"`
	Icon         string `json:"icon"`
	Claimed      bool   `json:"claimed"`
	Unlocked     bool   `json:"unlocked"`
}

// SeasonLevelUpPayload 賽季等級升級通知（Server → Client）
// 當玩家成功領取等級獎勵時發送
type SeasonLevelUpPayload struct {
	PlayerID    string `json:"player_id"`
	Level       int    `json:"level"`
	CoinReward  int    `json:"coin_reward"`
	NewBalance  int    `json:"new_balance"`
	SpecialType string `json:"special_type"` // "" / "skin" / "title"
	SpecialID   string `json:"special_id"`
	SpecialName string `json:"special_name"`
}

// ---- 好友系統（DAY-073）----

// SendFriendRequestPayload 發送好友請求（Client → Server）
type SendFriendRequestPayload struct {
	TargetID string `json:"target_id"` // 目標玩家 ID
}

// AcceptFriendRequestPayload 接受好友請求（Client → Server）
type AcceptFriendRequestPayload struct {
	FromID string `json:"from_id"` // 發送請求的玩家 ID
}

// RejectFriendRequestPayload 拒絕好友請求（Client → Server）
type RejectFriendRequestPayload struct {
	FromID string `json:"from_id"`
}

// RemoveFriendPayload 移除好友（Client → Server）
type RemoveFriendPayload struct {
	FriendID string `json:"friend_id"`
}

// FriendInfoPayload 好友資訊（用於列表顯示）
type FriendInfoPayload struct {
	PlayerID     string `json:"player_id"`
	DisplayName  string `json:"display_name"`
	IsOnline     bool   `json:"is_online"`
	Coins        int    `json:"coins"`
	KillCount    int    `json:"kill_count"`
	TitleName    string `json:"title_name"`
	TitleIcon    string `json:"title_icon"`
	SeasonLevel  int    `json:"season_level"`
	SeasonPoints int    `json:"season_points"`
}

// FriendListPayload 好友列表（Server → Client）
type FriendListPayload struct {
	Friends        []FriendInfoPayload `json:"friends"`
	PendingCount   int                 `json:"pending_count"` // 待處理請求數
}

// FriendRequestPayload 好友請求通知（Server → Client）
// 當有人發送好友請求時通知目標玩家
type FriendRequestPayload struct {
	FromID      string `json:"from_id"`
	DisplayName string `json:"display_name"`
}

// FriendUpdatePayload 好友狀態更新（Server → Client）
// 當好友上線/下線/積分變化時廣播
type FriendUpdatePayload struct {
	FriendID    string `json:"friend_id"`
	DisplayName string `json:"display_name"`
	IsOnline    bool   `json:"is_online"`
	Event       string `json:"event"` // "online" / "offline" / "accepted" / "removed"
}

// ---- 公會系統（DAY-074）----

// CreateGuildPayload 建立公會請求（Client → Server）
type CreateGuildPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// JoinGuildPayload 加入公會請求（Client → Server）
type JoinGuildPayload struct {
	GuildID string `json:"guild_id"`
}

// KickGuildMemberPayload 踢出成員請求（Client → Server）
type KickGuildMemberPayload struct {
	TargetID string `json:"target_id"`
}

// PromoteGuildMemberPayload 升職成員請求（Client → Server）
type PromoteGuildMemberPayload struct {
	TargetID string `json:"target_id"`
}

// GuildMemberInfo 公會成員資訊（用於廣播）
type GuildMemberInfo struct {
	PlayerID     string `json:"player_id"`
	DisplayName  string `json:"display_name"`
	Role         string `json:"role"`         // "leader" / "officer" / "member"
	IsOnline     bool   `json:"is_online"`
	Contribution int    `json:"contribution"`
}

// GuildTaskInfo 公會任務資訊（用於廣播）
type GuildTaskInfo struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Target      int    `json:"target"`
	Current     int    `json:"current"`
	Reward      int    `json:"reward"`
	Completed   bool   `json:"completed"`
	ResetAt     int64  `json:"reset_at"` // Unix ms
}

// GuildUpdatePayload 公會資訊更新（Server → Client）
// 在玩家加入公會、任務進度更新、成員變動時發送
type GuildUpdatePayload struct {
	GuildID     string            `json:"guild_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Icon        string            `json:"icon"`
	Level       int               `json:"level"`
	Exp         int               `json:"exp"`
	Members     []GuildMemberInfo `json:"members"`
	Tasks       []GuildTaskInfo   `json:"tasks"`
	TotalKills  int               `json:"total_kills"`
	TotalCoins  int               `json:"total_coins"`
	// 接收者的角色（方便 Client 判斷操作權限）
	MyRole      string `json:"my_role"` // "leader" / "officer" / "member" / ""（不在公會）
}

// GuildListEntry 公會列表單筆記錄（用於搜尋）
type GuildListEntry struct {
	GuildID     string `json:"guild_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Level       int    `json:"level"`
	MemberCount int    `json:"member_count"`
	OnlineCount int    `json:"online_count"`
}

// GuildListPayload 公會列表（Server → Client）
type GuildListPayload struct {
	Guilds []GuildListEntry `json:"guilds"`
}

// GuildTaskCompletePayload 公會任務完成通知（Server → Client）
type GuildTaskCompletePayload struct {
	GuildID    string `json:"guild_id"`
	GuildName  string `json:"guild_name"`
	TaskID     string `json:"task_id"`
	TaskName   string `json:"task_name"`
	TaskIcon   string `json:"task_icon"`
	Reward     int    `json:"reward"`     // 每人獎勵
	NewBalance int    `json:"new_balance"` // 領取後餘額
}

// GuildErrorPayload 公會操作錯誤（Server → Client）
type GuildErrorPayload struct {
	Operation string `json:"operation"` // 操作類型
	Message   string `json:"message"`   // 錯誤訊息
}

// ---- 公會聊天室（DAY-075）----

// GuildChatPayload 公會聊天訊息（Client → Server）
type GuildChatPayload struct {
	Message string `json:"message"` // 聊天內容（最多 100 字）
}

// GuildMessagePayload 公會聊天廣播（Server → Client）
type GuildMessagePayload struct {
	GuildID     string `json:"guild_id"`
	PlayerID    string `json:"player_id"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`    // "leader" / "officer" / "member"
	Message     string `json:"message"`
	Timestamp   int64  `json:"timestamp"` // Unix ms
}

// ---- 公會戰系統（DAY-076）----

// GuildWarScoreEntry 公會戰積分條目（用於排名顯示）
type GuildWarScoreEntry struct {
	Rank        int    `json:"rank"`
	GuildID     string `json:"guild_id"`
	GuildName   string `json:"guild_name"`
	GuildIcon   string `json:"guild_icon"`
	MemberCount int    `json:"member_count"`
	Score       int    `json:"score"`
	KillScore   int    `json:"kill_score"`
	BossScore   int    `json:"boss_score"`
	BonusScore  int    `json:"bonus_score"`
	IsMyGuild   bool   `json:"is_my_guild"` // 是否為玩家所在公會
}

// GuildWarUpdatePayload 公會戰排名更新（Server → Client）
type GuildWarUpdatePayload struct {
	WeekID      string               `json:"week_id"`
	Status      string               `json:"status"`       // "active" / "settling" / "idle"
	EndAt       int64                `json:"end_at"`       // Unix timestamp（結束時間）
	Rankings    []GuildWarScoreEntry `json:"rankings"`
	MyGuildRank int                  `json:"my_guild_rank"` // 0 = 不在公會
	MyGuildScore int                 `json:"my_guild_score"`
	TotalGuilds int                  `json:"total_guilds"`
}

// GuildWarResultEntry 公會戰結算條目
type GuildWarResultEntry struct {
	Rank        int    `json:"rank"`
	GuildID     string `json:"guild_id"`
	GuildName   string `json:"guild_name"`
	GuildIcon   string `json:"guild_icon"`
	Score       int    `json:"score"`
	Reward      int    `json:"reward"` // 每人獎勵金幣
}

// GuildWarResultPayload 公會戰結算通知（Server → Client）
type GuildWarResultPayload struct {
	WeekID      string                `json:"week_id"`
	Rankings    []GuildWarResultEntry `json:"rankings"`
	MyRank      int                   `json:"my_rank"`    // 0 = 不在公會
	MyReward    int                   `json:"my_reward"`  // 本次獲得獎勵
	NextWarAt   int64                 `json:"next_war_at"` // 下週開始時間
}

// ---- 每日 BOSS 挑戰（DAY-077）----

// DailyBossAttackPayload 對每日 BOSS 攻擊（Client → Server）
type DailyBossAttackPayload struct {
	Damage int `json:"damage"` // 造成的傷害（由 Server 驗證）
}

// DailyBossContributorEntry 貢獻者條目
type DailyBossContributorEntry struct {
	Rank        int    `json:"rank"`
	PlayerID    string `json:"player_id"`
	DisplayName string `json:"display_name"`
	Damage      int    `json:"damage"`
	Reward      int    `json:"reward"` // 結算後填入
	IsMe        bool   `json:"is_me"`
}

// DailyBossUpdatePayload 每日 BOSS 狀態更新（Server → Client）
type DailyBossUpdatePayload struct {
	DateID        string                      `json:"date_id"`
	BossID        string                      `json:"boss_id"`
	BossName      string                      `json:"boss_name"`
	BossIcon      string                      `json:"boss_icon"`
	BossColor     string                      `json:"boss_color"`
	Description   string                      `json:"description"`
	MaxHP         int                         `json:"max_hp"`
	CurrentHP     int                         `json:"current_hp"`
	HPPercent     float64                     `json:"hp_percent"`
	Status        string                      `json:"status"`        // "active" / "defeated" / "expired"
	EndAt         int64                       `json:"end_at"`        // Unix ms
	RewardPool    int                         `json:"reward_pool"`
	TopContribs   []DailyBossContributorEntry `json:"top_contribs"`  // 前 5 名
	MyDamage      int                         `json:"my_damage"`
	MyReward      int                         `json:"my_reward"`     // 結算後填入
	DifficultyMod float64                     `json:"difficulty_mod"`
}

// DailyBossDefeatedPayload 每日 BOSS 擊殺通知（Server → Client）
type DailyBossDefeatedPayload struct {
	DateID      string                      `json:"date_id"`
	BossName    string                      `json:"boss_name"`
	BossIcon    string                      `json:"boss_icon"`
	KillerID    string                      `json:"killer_id"`
	KillerName  string                      `json:"killer_name"`
	Rankings    []DailyBossContributorEntry `json:"rankings"`
	MyReward    int                         `json:"my_reward"`
	TotalDamage int                         `json:"total_damage"`
}
