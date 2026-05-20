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
	// VIP 等級系統（DAY-078）
	MsgGetVIPStatus      MessageType = "get_vip_status"       // 查詢 VIP 狀態
	MsgClaimVIPWeekly    MessageType = "claim_vip_weekly"     // 領取 VIP 週獎勵
	// 限時活動系統（DAY-079）
	MsgGetEventStatus    MessageType = "get_event_status"     // 查詢限時活動狀態
	// 魚類圖鑑系統（DAY-081）
	MsgGetCodex          MessageType = "get_codex"            // 查詢圖鑑狀態
	// 推薦碼系統（DAY-082）
	MsgGetReferralInfo   MessageType = "get_referral_info"    // 查詢推薦碼資訊
	MsgUseReferralCode   MessageType = "use_referral_code"    // 使用推薦碼
	// 特殊武器系統（DAY-089）
	MsgBuySpecialWeapon  MessageType = "buy_special_weapon"   // 購買特殊武器
	MsgUseSpecialWeapon  MessageType = "use_special_weapon"   // 使用特殊武器
	MsgGetSpecialWeapons MessageType = "get_special_weapons"  // 查詢特殊武器狀態
	// 神秘寶箱系統（DAY-090）
	MsgOpenMysteryBox    MessageType = "open_mystery_box"     // 開箱請求
	MsgGetMysteryBoxes   MessageType = "get_mystery_boxes"    // 查詢持有寶箱
	// 房間難度系統（DAY-091）
	MsgGetRoomList       MessageType = "get_room_list"        // 查詢房間列表
	MsgSwitchRoom        MessageType = "switch_room"          // 切換房間
	// 每日簽到轉盤（DAY-092）
	MsgGetDailySpin      MessageType = "get_daily_spin"       // 查詢每日轉盤狀態
	MsgDailySpin         MessageType = "daily_spin"           // 執行每日轉盤
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
	MsgJackpotAnimation MessageType = "jackpot_animation" // Jackpot 觸發動畫通知（DAY-095）
	MsgDailyBonus        MessageType = "daily_bonus"        // 每日登入獎勵（DAY-065）
	MsgTournamentUpdate  MessageType = "tournament_update"  // 週賽排名更新（DAY-066）
	MsgTournamentResult  MessageType = "tournament_result"  // 週賽結算通知（DAY-066）
	MsgGetTournament     MessageType = "get_tournament"     // 查詢週賽/日賽狀態（DAY-093）
	MsgDailyTournamentUpdate MessageType = "daily_tournament_update" // 每日賽排名更新（DAY-093）
	MsgDailyTournamentResult MessageType = "daily_tournament_result" // 每日賽結算通知（DAY-093）
	// 商店系統（DAY-094）
	MsgGetShop      MessageType = "get_shop"       // 查詢商店狀態（Client→Server）
	MsgBuyShopItem  MessageType = "buy_shop_item"  // 購買商品（Client→Server）
	MsgShopUpdate   MessageType = "shop_update"    // 商店狀態更新（Server→Client）
	MsgShopPurchased MessageType = "shop_purchased" // 購買成功通知（Server→Client）
	MsgShopError    MessageType = "shop_error"     // 購買失敗通知（Server→Client）
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
	// VIP 等級系統（DAY-078）
	MsgVIPUpdate         MessageType = "vip_update"          // VIP 狀態更新
	MsgVIPLevelUp        MessageType = "vip_level_up"        // VIP 升級通知
	MsgVIPWeeklyClaimed  MessageType = "vip_weekly_claimed"  // VIP 週獎勵領取通知
	// 限時活動系統（DAY-079）
	MsgEventUpdate       MessageType = "event_update"        // 限時活動狀態更新
	// 魚類圖鑑系統（DAY-081）
	MsgCodexUpdate       MessageType = "codex_update"        // 圖鑑狀態更新
	MsgCodexUnlock       MessageType = "codex_unlock"        // 圖鑑條目解鎖通知
	MsgCodexComplete     MessageType = "codex_complete"      // 全圖鑑完成通知
	// 推薦碼系統（DAY-082）
	MsgReferralInfo      MessageType = "referral_info"       // 推薦碼資訊（Server → Client）
	MsgReferralSuccess   MessageType = "referral_success"    // 推薦碼使用成功通知
	MsgReferralError     MessageType = "referral_error"      // 推薦碼使用失敗通知
	// 連擊系統（DAY-083）
	MsgStreakUpdate      MessageType = "streak_update"       // 連擊狀態更新
	MsgStreakReset       MessageType = "streak_reset"        // 連擊重置通知
	// 幸運轉盤系統（DAY-084）
	MsgWheelTrigger     MessageType = "wheel_trigger"       // 轉盤觸發通知
	// 隱藏挑戰系統（DAY-085）
	MsgChallengeUnlocked MessageType = "challenge_unlocked" // 挑戰解鎖通知
	// 每日任務連續完成獎勵（DAY-086）
	MsgMissionStreakBonus MessageType = "mission_streak_bonus" // 連續完成獎勵通知
	// 天氣系統（DAY-087）
	MsgWeatherUpdate MessageType = "weather_update" // 天氣狀態更新
	// 連鎖爆炸系統（DAY-088）
	MsgChainExplosion MessageType = "chain_explosion" // 連鎖爆炸通知
	// 特殊武器系統（DAY-089）
	MsgSpecialWeaponUpdate MessageType = "special_weapon_update" // 特殊武器狀態更新
	MsgSpecialWeaponFired  MessageType = "special_weapon_fired"  // 特殊武器發射廣播（所有玩家可見）
	// 神秘寶箱系統（DAY-090）
	MsgMysteryBoxDrop    MessageType = "mystery_box_drop"    // 寶箱掉落通知（擊破目標後）
	MsgMysteryBoxUpdate  MessageType = "mystery_box_update"  // 持有寶箱狀態更新
	MsgMysteryBoxOpened  MessageType = "mystery_box_opened"  // 開箱結果通知
	// 房間難度系統（DAY-091）
	MsgRoomList          MessageType = "room_list"           // 房間列表（Server → Client）
	MsgRoomSwitched      MessageType = "room_switched"       // 房間切換成功通知
	MsgRoomError         MessageType = "room_error"          // 房間操作失敗通知
	// 每日簽到轉盤（DAY-092）
	MsgDailySpinState    MessageType = "daily_spin_state"    // 每日轉盤狀態（Server → Client）
	MsgDailySpinResult   MessageType = "daily_spin_result"   // 每日轉盤結果（Server → Client）
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

// ---- Progressive Jackpot（DAY-048，DAY-095 升級四層）----

// JackpotUpdatePayload Jackpot 池更新廣播（每 5 秒）
type JackpotUpdatePayload struct {
	Mini  int `json:"mini"`  // Mini Jackpot 當前金額
	Minor int `json:"minor"` // Minor Jackpot 當前金額（DAY-095）
	Major int `json:"major"` // Major Jackpot 當前金額
	Grand int `json:"grand"` // Grand Jackpot 當前金額
}

// JackpotWinPayload Jackpot 中獎通知
type JackpotWinPayload struct {
	Level      string `json:"level"`       // "mini" / "minor" / "major" / "grand"
	LevelName  string `json:"level_name"`  // "MINI" / "MINOR" / "MAJOR" / "GRAND"
	LevelColor string `json:"level_color"` // 顯示顏色（DAY-095）
	LevelIcon  string `json:"level_icon"`  // 圖示（DAY-095）
	Amount     int    `json:"amount"`      // 中獎金額
	WinnerID   string `json:"winner_id"`   // 中獎玩家 ID
	WinnerName string `json:"winner_name"` // 中獎玩家顯示名稱
	NewBalance int    `json:"new_balance"` // 中獎後餘額（只對中獎玩家有意義）
	IsGrand    bool   `json:"is_grand"`    // 是否是 Grand Jackpot（觸發全畫面動畫）
}

// JackpotAnimationPayload Jackpot 觸發動畫通知（DAY-095）
// 廣播給所有玩家，觸發對應等級的動畫效果
type JackpotAnimationPayload struct {
	Level      string `json:"level"`       // "mini" / "minor" / "major" / "grand"
	LevelName  string `json:"level_name"`  // 顯示名稱
	LevelColor string `json:"level_color"` // 顯示顏色
	LevelIcon  string `json:"level_icon"`  // 圖示
	Amount     int    `json:"amount"`      // 中獎金額
	WinnerName string `json:"winner_name"` // 中獎玩家名稱
	IsGrand    bool   `json:"is_grand"`    // Grand 觸發全畫面特效
	IsMajor    bool   `json:"is_major"`    // Major 觸發半畫面特效
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

// ---- 每日賽系統（DAY-093）----

// DailyTournamentUpdatePayload 每日賽排名更新廣播（每 30 秒）
type DailyTournamentUpdatePayload struct {
	DayStart     int64                 `json:"day_start"`     // Unix ms
	DayEnd       int64                 `json:"day_end"`       // Unix ms
	SecondsLeft  int64                 `json:"seconds_left"`  // 距離結束秒數
	Rankings     []TournamentRankEntry `json:"rankings"`      // 前 10 名
	TotalPlayers int                   `json:"total_players"` // 今日參賽人數
	PlayerRank   int                   `json:"player_rank"`   // 接收者的排名（0=未上榜）
	PlayerPoints int                   `json:"player_points"` // 接收者的積分
}

// DailyTournamentResultPayload 每日賽結算通知（每日結束時廣播）
type DailyTournamentResultPayload struct {
	Date       string                `json:"date"`        // "2026-05-20"
	Rankings   []TournamentRankEntry `json:"rankings"`
	Prize      int                   `json:"prize"`       // 接收者獲得的獎勵（0=未獲獎）
	PrizeLabel string                `json:"prize_label"` // 獎勵標籤
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

// ---- VIP 等級系統（DAY-078）----

// VIPUpdatePayload VIP 狀態更新（Server → Client）
// 在玩家加入、消費、升級時發送
type VIPUpdatePayload struct {
	PlayerID       string  `json:"player_id"`
	TotalSpend     int     `json:"total_spend"`
	VIPLevel       int     `json:"vip_level"`
	TierName       string  `json:"tier_name"`
	TierIcon       string  `json:"tier_icon"`
	TierColor      string  `json:"tier_color"`
	CashbackRate   float64 `json:"cashback_rate"`
	DailyBonusMult float64 `json:"daily_bonus_mult"`
	WeeklyBonus    int     `json:"weekly_bonus"`
	NextLevel      int     `json:"next_level"`
	SpendToNext    int     `json:"spend_to_next"`
	Progress       float64 `json:"progress"`
	CanClaimWeekly bool    `json:"can_claim_weekly"`
}

// VIPLevelUpPayload VIP 升級通知（Server → Client）
type VIPLevelUpPayload struct {
	PlayerID    string `json:"player_id"`
	NewLevel    int    `json:"new_level"`
	TierName    string `json:"tier_name"`
	TierIcon    string `json:"tier_icon"`
	TierColor   string `json:"tier_color"`
	TitleID     string `json:"title_id"`
	TitleName   string `json:"title_name"`
	WeeklyBonus int    `json:"weekly_bonus"`
}

// VIPWeeklyClaimedPayload VIP 週獎勵領取通知（Server → Client）
type VIPWeeklyClaimedPayload struct {
	PlayerID   string `json:"player_id"`
	VIPLevel   int    `json:"vip_level"`
	TierName   string `json:"tier_name"`
	Coins      int    `json:"coins"`
	NewBalance int    `json:"new_balance"`
}

// ---- 限時活動系統（DAY-079）----

// EventUpdatePayload 限時活動狀態更新（Server → Client）
// 在玩家加入、活動切換時發送
type EventUpdatePayload struct {
	Type          string  `json:"type"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Icon          string  `json:"icon"`
	Color         string  `json:"color"`
	IsActive      bool    `json:"is_active"`
	EndAt         int64   `json:"end_at"`    // Unix ms
	TimeLeft      float64 `json:"time_left"` // 秒
	RewardMult    float64 `json:"reward_mult"`
	SpawnMult     float64 `json:"spawn_mult"`
	KillChanceAdd float64 `json:"kill_chance_add"`
}

// ---- 魚類圖鑑系統（DAY-081）----

// CodexEntryPayload 圖鑑條目（用於廣播）
type CodexEntryPayload struct {
	TargetID      string  `json:"target_id"`
	TargetName    string  `json:"target_name"`
	Rarity        string  `json:"rarity"`         // "common"/"rare"/"epic"/"legendary"
	Unlocked      bool    `json:"unlocked"`
	UnlockedAt    int64   `json:"unlocked_at"`    // Unix ms，0=未解鎖
	KillCount     int     `json:"kill_count"`
	MaxMultiplier float64 `json:"max_multiplier"`
}

// CodexUpdatePayload 圖鑑狀態更新（Server → Client）
// 玩家加入時發送完整圖鑑，解鎖時也發送
type CodexUpdatePayload struct {
	Entries       []CodexEntryPayload `json:"entries"`
	UnlockedCount int                 `json:"unlocked_count"`
	TotalCount    int                 `json:"total_count"`
	IsComplete    bool                `json:"is_complete"`
}

// CodexUnlockPayload 圖鑑條目解鎖通知（Server → Client）
// 首次擊破某種目標物時發送
type CodexUnlockPayload struct {
	TargetID      string  `json:"target_id"`
	TargetName    string  `json:"target_name"`
	Rarity        string  `json:"rarity"`
	Reward        int     `json:"reward"`        // 解鎖獎勵金幣
	NewBalance    int     `json:"new_balance"`
	UnlockedCount int     `json:"unlocked_count"`
	TotalCount    int     `json:"total_count"`
}

// CodexCompletePayload 全圖鑑完成通知（Server → Client）
type CodexCompletePayload struct {
	Reward     int    `json:"reward"`      // 完成獎勵金幣
	NewBalance int    `json:"new_balance"`
	TitleID    string `json:"title_id"`    // 解鎖的稱號 ID
	TitleName  string `json:"title_name"`
}

// ---- 推薦碼系統（DAY-082）----

// UseReferralCodePayload 使用推薦碼請求（Client → Server）
type UseReferralCodePayload struct {
	Code string `json:"code"` // 推薦碼（6位英數字）
}

// ReferralInfoPayload 推薦碼資訊（Server → Client）
type ReferralInfoPayload struct {
	MyCode        string `json:"my_code"`        // 我的推薦碼
	UsedCode      string `json:"used_code"`      // 我使用的推薦碼（空=未使用）
	ReferredBy    string `json:"referred_by"`    // 推薦我的玩家 ID
	ReferralCount int    `json:"referral_count"` // 我成功推薦的人數
	TotalReward   int    `json:"total_reward"`   // 累計推薦獎勵
	ReferrerReward int   `json:"referrer_reward"` // 每次推薦獎勵金幣
	RefereeReward  int   `json:"referee_reward"`  // 被推薦人獎勵金幣
	MaxReferrals   int   `json:"max_referrals"`   // 最多推薦人數
}

// ReferralSuccessPayload 推薦碼使用成功通知（Server → Client）
type ReferralSuccessPayload struct {
	Code        string `json:"code"`
	ReferrerID  string `json:"referrer_id"`
	Reward      int    `json:"reward"`      // 被推薦人獲得的獎勵
	NewBalance  int    `json:"new_balance"`
	Message     string `json:"message"`
}

// ReferralErrorPayload 推薦碼使用失敗通知（Server → Client）
type ReferralErrorPayload struct {
	Code    string `json:"code"`
	Reason  string `json:"reason"`
}

// ---- 連擊系統（DAY-082）----

// StreakUpdatePayload 連擊狀態更新（Server → Client）
// 每次擊破目標後發送
type StreakUpdatePayload struct {
	Current    int     `json:"current"`     // 當前連擊數
	MultBonus  float64 `json:"mult_bonus"`  // 獎勵倍率加成
	LevelName  string  `json:"level_name"`  // 等級名稱
	LevelColor string  `json:"level_color"` // 顯示顏色（hex）
	IsNewLevel bool    `json:"is_new_level"` // 是否剛升到新等級
	MaxStreak  int     `json:"max_streak"`  // 本局最高連擊
}

// StreakResetPayload 連擊重置通知（Server → Client）
// 超時未擊破時發送
type StreakResetPayload struct {
	FinalStreak int `json:"final_streak"` // 重置前的連擊數
	MaxStreak   int `json:"max_streak"`   // 本局最高連擊
}

// ---- 幸運轉盤系統（DAY-084）----

// WheelSlotPayload 轉盤格子定義（用於廣播）
type WheelSlotPayload struct {
	Multiplier float64 `json:"multiplier"`
	Label      string  `json:"label"`
	Color      string  `json:"color"`
}

// WheelTriggerPayload 轉盤觸發通知（Server → Client）
// 擊殺特殊目標後觸發，Client 播放轉盤動畫後顯示結果
type WheelTriggerPayload struct {
	PlayerID    string             `json:"player_id"`
	TargetID    string             `json:"target_id"`    // 觸發的目標物 ID
	TargetName  string             `json:"target_name"`  // 目標物名稱
	Slots       []WheelSlotPayload `json:"slots"`        // 所有格子定義
	WinIndex    int                `json:"win_index"`    // 中獎格子索引
	Multiplier  float64            `json:"multiplier"`   // 中獎倍率
	BaseReward  int                `json:"base_reward"`  // 基礎獎勵
	FinalReward int                `json:"final_reward"` // 最終獎勵
	NewBalance  int                `json:"new_balance"`  // 新金幣餘額
}

// ---- 每日任務連續完成獎勵（DAY-086）----

// MissionStreakBonusPayload 連續完成獎勵通知（Server → Client）
type MissionStreakBonusPayload struct {
	Streak     int    `json:"streak"`      // 當前連續天數
	MaxStreak  int    `json:"max_streak"`  // 歷史最高連續天數
	Reward     int    `json:"reward"`      // 本次獎勵金幣
	Label      string `json:"label"`       // 獎勵標籤（如「連續 7 天 🏆」）
	NewBalance int    `json:"new_balance"` // 領取後餘額
}

// ---- 天氣系統（DAY-087）----

// WeatherUpdatePayload 天氣狀態更新（Server → Client）
// 玩家加入時發送當前天氣，天氣切換時廣播給所有玩家
type WeatherUpdatePayload struct {
	Type             string  `json:"type"`              // 天氣類型
	Name             string  `json:"name"`              // 顯示名稱
	Icon             string  `json:"icon"`              // 圖示（emoji）
	Description      string  `json:"description"`       // 效果說明
	RemainingSeconds int     `json:"remaining_seconds"` // 剩餘秒數
	SpawnRateMult    float64 `json:"spawn_rate_mult"`   // 目標生成倍率
	RewardMult       float64 `json:"reward_mult"`       // 獎勵倍率
	SpeedMult        float64 `json:"speed_mult"`        // 目標移動速度倍率
	RareChanceBonus  float64 `json:"rare_chance_bonus"` // 稀有目標出現機率加成
	GoldFishBonus    float64 `json:"gold_fish_bonus"`   // 金幣魚出現機率加成
	BossChanceBonus  float64 `json:"boss_chance_bonus"` // BOSS 出現機率加成
	FogEffect        bool    `json:"fog_effect"`        // 是否有濃霧效果
	IsNew            bool    `json:"is_new"`            // 是否剛切換（用於 Client 端顯示通知）
}

// ---- 連鎖爆炸系統（DAY-088）----

// ChainKillEntry 連鎖擊破的目標條目
type ChainKillEntry struct {
	InstanceID string  `json:"instance_id"`
	DefID      string  `json:"def_id"`
	Multiplier float64 `json:"multiplier"`
	Reward     int     `json:"reward"`
}

// ChainExplosionPayload 連鎖爆炸通知（Server → Client）
// 擊破目標後觸發連鎖，周圍目標同時爆炸
type ChainExplosionPayload struct {
	TriggerID   string           `json:"trigger_id"`   // 觸發連鎖的目標 ID
	Level       int              `json:"level"`        // 連鎖等級（1-4）
	LevelName   string           `json:"level_name"`   // 等級名稱
	LevelColor  string           `json:"level_color"`  // 顯示顏色（hex）
	Chains      []ChainKillEntry `json:"chains"`       // 被連鎖擊破的目標列表
	TotalReward int              `json:"total_reward"` // 連鎖總獎勵
	BonusMult   float64          `json:"bonus_mult"`   // 連鎖獎勵倍率加成
	PlayerID    string           `json:"player_id"`    // 觸發玩家 ID
}

// ---- 特殊武器系統（DAY-089）----

// BuySpecialWeaponPayload 購買特殊武器請求（Client → Server）
type BuySpecialWeaponPayload struct {
	WeaponType string `json:"weapon_type"` // "bomb" / "laser" / "freeze"
}

// UseSpecialWeaponPayload 使用特殊武器請求（Client → Server）
type UseSpecialWeaponPayload struct {
	WeaponType string  `json:"weapon_type"` // "bomb" / "laser" / "freeze"
	ClickX     float64 `json:"click_x"`     // 點擊位置 X（炸彈/雷射用）
	ClickY     float64 `json:"click_y"`     // 點擊位置 Y（炸彈/雷射用）
}

// SpecialWeaponDef 特殊武器定義（用於廣播）
type SpecialWeaponDef struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Cost        int    `json:"cost"`
	MaxCharges  int    `json:"max_charges"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}

// SpecialWeaponUpdatePayload 特殊武器狀態更新（Server → Client）
// 玩家加入、購買、使用後發送
type SpecialWeaponUpdatePayload struct {
	PlayerID      string             `json:"player_id"`
	BombCharges   int                `json:"bomb_charges"`
	LaserCharges  int                `json:"laser_charges"`
	FreezeCharges int                `json:"freeze_charges"`
	NewBalance    int                `json:"new_balance"`    // 購買後的新餘額（0=使用操作）
	Definitions   []SpecialWeaponDef `json:"definitions"`    // 武器定義（首次發送時填入）
}

// SpecialWeaponHitEntry 特殊武器命中的目標條目
type SpecialWeaponHitEntry struct {
	InstanceID string  `json:"instance_id"`
	DefID      string  `json:"def_id"`
	Multiplier float64 `json:"multiplier"`
	Reward     int     `json:"reward"`
	Killed     bool    `json:"killed"` // 是否被擊破
}

// SpecialWeaponFiredPayload 特殊武器發射廣播（Server → Client，廣播給所有玩家）
// 讓所有玩家看到特殊武器效果
type SpecialWeaponFiredPayload struct {
	PlayerID    string                  `json:"player_id"`
	WeaponType  string                  `json:"weapon_type"`
	ClickX      float64                 `json:"click_x"`
	ClickY      float64                 `json:"click_y"`
	HitTargets  []SpecialWeaponHitEntry `json:"hit_targets"`
	TotalReward int                     `json:"total_reward"`
	NewBalance  int                     `json:"new_balance"` // 只有發射者有值
	FreezeTime  float64                 `json:"freeze_time"` // 冰凍持續秒數（freeze 武器用）
}

// ---- 神秘寶箱系統（DAY-090）----

// OpenMysteryBoxPayload 開箱請求（Client → Server）
type OpenMysteryBoxPayload struct {
	Rarity string `json:"rarity"` // "common" / "rare" / "epic" / "legendary"
}

// MysteryBoxInventoryEntry 持有寶箱條目
type MysteryBoxInventoryEntry struct {
	Rarity    string `json:"rarity"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Color     string `json:"color"`
	GlowColor string `json:"glow_color"`
	Count     int    `json:"count"`
}

// MysteryBoxUpdatePayload 持有寶箱狀態更新（Server → Client）
type MysteryBoxUpdatePayload struct {
	PlayerID  string                     `json:"player_id"`
	Inventory []MysteryBoxInventoryEntry `json:"inventory"`
	Total     int                        `json:"total"` // 總持有數量
}

// MysteryBoxDropPayload 寶箱掉落通知（Server → Client）
// 擊破目標後掉落寶箱時廣播
type MysteryBoxDropPayload struct {
	PlayerID  string `json:"player_id"`
	Rarity    string `json:"rarity"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Color     string `json:"color"`
	GlowColor string `json:"glow_color"`
	DropX     float64 `json:"drop_x"` // 掉落位置（目標物位置）
	DropY     float64 `json:"drop_y"`
}

// MysteryBoxRewardPayload 開箱獎勵條目
type MysteryBoxRewardPayload struct {
	Type   string `json:"type"`
	Amount int    `json:"amount"`
	Label  string `json:"label"`
	Icon   string `json:"icon"`
	Color  string `json:"color"`
}

// MysteryBoxOpenedPayload 開箱結果通知（Server → Client）
type MysteryBoxOpenedPayload struct {
	PlayerID        string                  `json:"player_id"`
	Rarity          string                  `json:"rarity"`
	BoxName         string                  `json:"box_name"`
	BoxIcon         string                  `json:"box_icon"`
	Reward          MysteryBoxRewardPayload `json:"reward"`
	NewBalance      int                     `json:"new_balance"`
	PendingMultMult float64                 `json:"pending_mult"` // 待使用倍率（0=無）
	RemainingBoxes  int                     `json:"remaining_boxes"` // 該稀有度剩餘數量
}

// ---- 房間難度系統（DAY-091）----

// RoomDifficultyInfo 房間難度資訊（用於 UI 顯示）
type RoomDifficultyInfo struct {
	ID              string  `json:"id"`               // "beginner"/"intermediate"/"advanced"/"vip"
	Name            string  `json:"name"`             // 顯示名稱
	Icon            string  `json:"icon"`             // emoji 圖示
	Color           string  `json:"color"`            // 主題色
	MinBetCost      int     `json:"min_bet_cost"`     // 最低 bet 金幣
	MaxBetCost      int     `json:"max_bet_cost"`     // 最高 bet 金幣
	MaxPlayers      int     `json:"max_players"`      // 最大玩家數
	PlayerCount     int     `json:"player_count"`     // 當前玩家數
	RewardMult      float64 `json:"reward_mult"`      // 獎勵倍率
	JackpotMult     float64 `json:"jackpot_mult"`     // Jackpot 倍率
	EntryFee        int     `json:"entry_fee"`        // 進場費（0=免費）
	Description     string  `json:"description"`      // 房間描述
	IsAvailable     bool    `json:"is_available"`     // 是否可進入（未滿且有足夠金幣）
	IsCurrent       bool    `json:"is_current"`       // 是否是當前房間
}

// RoomListPayload 房間列表（Server → Client）
type RoomListPayload struct {
	Rooms       []RoomDifficultyInfo `json:"rooms"`
	CurrentRoom string               `json:"current_room"` // 當前所在房間 ID
}

// SwitchRoomPayload 切換房間請求（Client → Server）
type SwitchRoomPayload struct {
	RoomID string `json:"room_id"` // 目標房間 ID（"beginner"/"intermediate"/"advanced"/"vip"）
}

// RoomSwitchedPayload 房間切換成功通知（Server → Client）
type RoomSwitchedPayload struct {
	RoomID      string  `json:"room_id"`
	RoomName    string  `json:"room_name"`
	RoomIcon    string  `json:"room_icon"`
	RoomColor   string  `json:"room_color"`
	RewardMult  float64 `json:"reward_mult"`
	JackpotMult float64 `json:"jackpot_mult"`
	EntryFee    int     `json:"entry_fee"`    // 已扣除的進場費
	NewBalance  int     `json:"new_balance"`  // 扣除進場費後的餘額
}

// RoomErrorPayload 房間操作失敗通知（Server → Client）
type RoomErrorPayload struct {
	Code    string `json:"code"`    // "room_full"/"insufficient_coins"/"room_not_found"
	Message string `json:"message"`
}

// ---- 每日簽到轉盤（DAY-092）----

// DailySpinSlotPayload 轉盤格子資訊
type DailySpinSlotPayload struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Amount  int    `json:"amount"`
	Label   string `json:"label"`
	Icon    string `json:"icon"`
	Color   string `json:"color"`
	IsSuper bool   `json:"is_super"`
}

// DailySpinStatePayload 每日轉盤狀態（Server → Client）
type DailySpinStatePayload struct {
	CanSpin     bool                   `json:"can_spin"`
	IsSuper     bool                   `json:"is_super"`
	LoginStreak int                    `json:"login_streak"`
	TotalSpins  int                    `json:"total_spins"`
	NextSpinAt  int64                  `json:"next_spin_at"` // Unix ms
	NormalSlots []DailySpinSlotPayload `json:"normal_slots"`
	SuperSlots  []DailySpinSlotPayload `json:"super_slots"`
}

// DailySpinResultPayload 每日轉盤結果（Server → Client）
type DailySpinResultPayload struct {
	SlotIndex   int                  `json:"slot_index"`
	Slot        DailySpinSlotPayload `json:"slot"`
	IsSuper     bool                 `json:"is_super"`
	LoginStreak int                  `json:"login_streak"`
	NextSpinAt  int64                `json:"next_spin_at"`
	NewBalance  int                  `json:"new_balance"`
	// 特殊獎勵欄位
	SeasonPoints  int     `json:"season_points"`   // 賽季積分（0=無）
	MultBonus     float64 `json:"mult_bonus"`      // 倍率加成（0=無）
	MysteryBoxRarity string `json:"mystery_box_rarity"` // 寶箱稀有度（""=無）
}

// ---- 商店系統（DAY-094）----

// ShopItemReward 商品獎勵（用於 Payload）
type ShopItemReward struct {
	Coins        int     `json:"coins"`
	BombCharge   int     `json:"bomb_charge"`
	LaserCharge  int     `json:"laser_charge"`
	FreezeCharge int     `json:"freeze_charge"`
	AttackMult   float64 `json:"attack_mult"`
	SeasonPoints int     `json:"season_points"`
}

// ShopItem 商品資訊（用於 Payload）
type ShopItem struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Type           string         `json:"type"`
	Price          int            `json:"price"`
	OrigPrice      int            `json:"orig_price"`
	Reward         ShopItemReward `json:"reward"`
	Stock          int            `json:"stock"`
	LimitPerDay    int            `json:"limit_per_day"`
	IsFlashSale    bool           `json:"is_flash_sale"`
	FlashEndAt     int64          `json:"flash_end_at"`
	PurchasedToday int            `json:"purchased_today"` // 今日已購買次數
}

// ShopUpdatePayload 商店狀態更新（Server → Client）
type ShopUpdatePayload struct {
	Items          []ShopItem `json:"items"`
	FlashSaleEndAt int64      `json:"flash_sale_end_at"` // Unix ms
	SecondsLeft    int64      `json:"seconds_left"`      // 特賣剩餘秒數
}

// BuyShopItemPayload 購買商品請求（Client → Server）
type BuyShopItemPayload struct {
	ItemID string `json:"item_id"`
}

// ShopPurchasedPayload 購買成功通知（Server → Client）
type ShopPurchasedPayload struct {
	ItemID     string         `json:"item_id"`
	ItemName   string         `json:"item_name"`
	Price      int            `json:"price"`
	NewBalance int            `json:"new_balance"`
	Reward     ShopItemReward `json:"reward"`
}

// ShopErrorPayload 購買失敗通知（Server → Client）
type ShopErrorPayload struct {
	ItemID string `json:"item_id"`
	Reason string `json:"reason"` // "item_not_found" / "insufficient_coins" / "daily_limit_reached" / "out_of_stock"
}
