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
	// 好友禮物系統（DAY-101）
	MsgSendGift          MessageType = "send_gift"            // 送禮物給好友
	MsgGetGiftStatus     MessageType = "get_gift_status"      // 查詢今日禮物狀態
	// 好友挑戰系統（DAY-102）
	MsgSendChallengeRequest MessageType = "send_challenge_request" // 發起挑戰
	MsgAcceptChallenge      MessageType = "accept_challenge"       // 接受挑戰
	MsgDeclineChallenge     MessageType = "decline_challenge"      // 拒絕挑戰
	// 私訊系統（DAY-103）
	MsgSendDM    MessageType = "send_dm"    // 發送私訊
	// 成就動態牆系統（DAY-112）
	MsgGetActivityFeed   MessageType = "get_activity_feed"   // 查詢最近動態（Client→Server）
	// 雙層倍率輪盤系統（DAY-113）
	MsgSpinRoulette      MessageType = "spin_roulette"       // 玩家手動停止輪盤（Client→Server）
	// Buy Bonus 系統（DAY-114）
	MsgBuyBonus          MessageType = "buy_bonus"           // 購買 Bonus 觸發（Client→Server）
	MsgGetBuyBonusStatus MessageType = "get_buy_bonus_status" // 查詢今日購買狀態（Client→Server）
	// 新手引導系統（DAY-115）
	MsgTutorialAction    MessageType = "tutorial_action"     // 玩家完成引導步驟（Client→Server）
	MsgSkipTutorial      MessageType = "skip_tutorial"       // 跳過引導（Client→Server）
	// Co-op Boss Raid 系統（DAY-115）
	MsgGetRaidStatus     MessageType = "get_raid_status"     // 查詢討伐狀態（Client→Server）
	MsgTriggerRaid       MessageType = "trigger_raid"        // 手動觸發討伐（Prototype 展示用）
	// 碎片收集大獎系統（DAY-116）
	MsgGetFragments      MessageType = "get_fragments"       // 查詢碎片狀態（Client→Server）
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
	MsgLoginMilestone    MessageType = "login_milestone"    // 登入里程碑達成通知（DAY-107）
	MsgLoginProgress     MessageType = "login_progress"     // 登入進度回應（DAY-107）
	MsgSuperBonusReady   MessageType = "super_bonus_ready"  // 超級 Bonus 準備通知（DAY-108）
	MsgTournamentUpdate  MessageType = "tournament_update"  // 週賽排名更新（DAY-066）
	MsgTournamentResult  MessageType = "tournament_result"  // 週賽結算通知（DAY-066）
	MsgGetTournament     MessageType = "get_tournament"     // 查詢週賽/日賽狀態（DAY-093）
	MsgDailyTournamentUpdate MessageType = "daily_tournament_update" // 每日賽排名更新（DAY-093）
	MsgDailyTournamentResult MessageType = "daily_tournament_result" // 每日賽結算通知（DAY-093）
	// 多格式每日賽系統（DAY-111）
	MsgGetMultiFormat        MessageType = "get_multi_format"        // 查詢多格式賽狀態（Client→Server）
	MsgMultiFormatUpdate     MessageType = "multi_format_update"     // 多格式賽排名更新（Server→Client）
	// 商店系統（DAY-094）
	MsgGetShop      MessageType = "get_shop"       // 查詢商店狀態（Client→Server）
	MsgBuyShopItem  MessageType = "buy_shop_item"  // 購買商品（Client→Server）
	MsgShopUpdate   MessageType = "shop_update"    // 商店狀態更新（Server→Client）
	MsgShopPurchased MessageType = "shop_purchased" // 購買成功通知（Server→Client）
	MsgShopError    MessageType = "shop_error"     // 購買失敗通知（Server→Client）
	// 玩家統計系統（DAY-096）
	MsgGetPlayerStats  MessageType = "get_player_stats"  // 查詢個人統計（Client→Server）
	MsgPlayerStatsUpdate MessageType = "player_stats_update" // 個人統計更新（Server→Client）
	// 全服公告系統（DAY-097）
	MsgAnnouncement MessageType = "announcement" // 全服公告廣播（Server→Client）
	// 成就動態牆系統（DAY-112）
	MsgActivityFeedEvent  MessageType = "activity_feed_event"  // 新動態事件廣播（Server→Client）
	MsgActivityFeedHistory MessageType = "activity_feed_history" // 歷史動態回應（Server→Client）
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
	// 任務連續寬限期（DAY-120）
	MsgMissionMercyProtected MessageType = "mission_mercy_protected" // 寬限期保護通知（Server→Client）
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
	// 好友禮物系統（DAY-101）
	MsgGiftReceived      MessageType = "gift_received"       // 收到禮物通知（Server → Client）
	MsgGiftSent          MessageType = "gift_sent"           // 送出禮物成功通知（Server → Client）
	MsgGiftStatus        MessageType = "gift_status"         // 今日禮物狀態（Server → Client）
	MsgGiftError         MessageType = "gift_error"          // 禮物操作失敗（Server → Client）
	// 好友挑戰系統（DAY-102）
	MsgChallengeRequest  MessageType = "challenge_request"   // 收到挑戰通知（Server → Client）
	MsgChallengeUpdate   MessageType = "challenge_update"    // 挑戰狀態/分數更新（Server → Client）
	MsgChallengeResult   MessageType = "challenge_result"    // 挑戰結果通知（Server → Client）
	MsgChallengeError    MessageType = "challenge_error"     // 挑戰操作失敗（Server → Client）
	// 私訊系統（DAY-103）
	MsgDMReceived MessageType = "dm_received" // 收到私訊（Server → Client）
	MsgDMSent     MessageType = "dm_sent"     // 發送成功確認（Server → Client）
	MsgDMError    MessageType = "dm_error"    // 發送失敗（Server → Client）
	// 玩家名片系統（DAY-106）
	MsgGetPlayerCard  MessageType = "get_player_card"  // 查詢玩家名片（Client → Server）
	MsgPlayerCard     MessageType = "player_card"      // 玩家名片資料（Server → Client）
	// 登入里程碑系統（DAY-107）
	MsgGetLoginProgress MessageType = "get_login_progress" // 查詢登入進度（Client → Server）
	// 賽季節日活動系統（DAY-109）
	MsgGetFestival          MessageType = "get_festival"           // 查詢節日狀態（Client → Server）
	MsgClaimFestivalTask    MessageType = "claim_festival_task"    // 領取節日任務獎勵（Client → Server）
	MsgFestivalUpdate       MessageType = "festival_update"        // 節日狀態更新（Server → Client）
	MsgFestivalTaskReady    MessageType = "festival_task_ready"    // 節日任務可領取通知（Server → Client）
	MsgFestivalTaskClaimed  MessageType = "festival_task_claimed"  // 節日任務獎勵領取成功（Server → Client）
	MsgFestivalTitleEarned  MessageType = "festival_title_earned"  // 節日稱號獲得通知（Server → Client）
	MsgFestivalError        MessageType = "festival_error"         // 節日操作失敗（Server → Client）
	// 名人堂系統（DAY-110）
	MsgGetHallOfFame        MessageType = "get_hall_of_fame"       // 查詢名人堂（Client → Server）
	MsgHallOfFameUpdate     MessageType = "hall_of_fame_update"    // 名人堂更新（Server → Client）
	MsgHallOfFameNewRecord  MessageType = "hall_of_fame_new_record" // 新記錄誕生廣播（Server → Client）
	// 智慧推薦系統（DAY-110）
	MsgGetRecommendations   MessageType = "get_recommendations"    // 查詢推薦（Client → Server）
	MsgRecommendations      MessageType = "recommendations"        // 推薦結果（Server → Client）
	// 雙層倍率輪盤系統（DAY-113）
	MsgRouletteStart    MessageType = "roulette_start"    // 輪盤開始（Server→Client，廣播給所有玩家）
	MsgRouletteResult   MessageType = "roulette_result"   // 輪盤結果（Server→Client，廣播給所有玩家）
	// Buy Bonus 系統（DAY-114）
	MsgBuyBonusSuccess  MessageType = "buy_bonus_success"  // 購買成功，Bonus 即將開始（Server→Client）
	MsgBuyBonusError    MessageType = "buy_bonus_error"    // 購買失敗（Server→Client）
	MsgBuyBonusStatus   MessageType = "buy_bonus_status"   // 今日購買狀態（Server→Client）
	// 新手引導系統（DAY-115）
	MsgTutorialStep     MessageType = "tutorial_step"      // 引導步驟（Server→Client）
	// Co-op Boss Raid 系統（DAY-115）
	MsgRaidWarning      MessageType = "raid_warning"       // 討伐警告廣播（Server→Client）
	MsgRaidStart        MessageType = "raid_start"         // 討伐開始廣播（Server→Client）
	MsgRaidUpdate       MessageType = "raid_update"        // 討伐狀態更新廣播（Server→Client，每 3 秒）
	MsgRaidResult       MessageType = "raid_result"        // 討伐結算廣播（Server→Client）
	MsgRaidStatus       MessageType = "raid_status"        // 討伐狀態回應（Server→Client）
	// 碎片收集大獎系統（DAY-116）
	MsgFragmentDrop     MessageType = "fragment_drop"      // 碎片掉落通知（Server→Client）
	MsgFragmentComplete MessageType = "fragment_complete"  // 集齊碎片大獎通知（Server→Client，廣播）
	MsgFragmentStatus   MessageType = "fragment_status"    // 碎片狀態回應（Server→Client）
	// 幸運捕獲系統（DAY-119）
	MsgLuckyCatch       MessageType = "lucky_catch"        // 幸運捕獲觸發廣播（Server→Client，廣播）
	// Rapid Respin 系統（DAY-121）
	MsgRapidRespin      MessageType = "rapid_respin"       // Rapid Respin 觸發廣播（Server→Client，廣播）
	MsgRapidRespinEnd   MessageType = "rapid_respin_end"   // Rapid Respin 連鎖結束通知（Server→Client）
	// 寶藏地圖系統（DAY-122）
	MsgTreasureMapUpdate MessageType = "treasure_map_update" // 寶藏地圖狀態更新（Server→Client）
	MsgTreasureMapLine   MessageType = "treasure_map_line"   // 完成一行/列/對角線通知（Server→Client）
	MsgTreasureMapFull   MessageType = "treasure_map_full"   // 完成整張地圖通知（Server→Client）
	// 寶藏地圖系統（DAY-122）
	MsgGetTreasureMap   MessageType = "get_treasure_map"   // 查詢寶藏地圖狀態（Client→Server）
	// 閃電挑戰系統（DAY-123）
	MsgGetFlashChallenge MessageType = "get_flash_challenge" // 查詢閃電挑戰狀態（Client→Server）
	MsgFlashChallengeStart  MessageType = "flash_challenge_start"  // 閃電挑戰開始廣播（Server→Client）
	MsgFlashChallengeUpdate MessageType = "flash_challenge_update" // 閃電挑戰進度更新（Server→Client）
	MsgFlashChallengeEnd    MessageType = "flash_challenge_end"    // 閃電挑戰結束廣播（Server→Client）
	MsgFlashChallengeReward MessageType = "flash_challenge_reward" // 閃電挑戰獎勵通知（Server→Client，個人）
	// 傳說目標警報系統（DAY-124）
	MsgRareTargetAlert MessageType = "rare_target_alert" // 稀有/傳說目標出現廣播（Server→Client）
	// 個人最佳記錄通知系統（DAY-125）
	MsgPersonalBest MessageType = "personal_best" // 個人最佳記錄通知（Server→Client，個人）
	// 黃金時間系統（DAY-125）
	MsgGetGoldenTime    MessageType = "get_golden_time"    // 查詢黃金時間狀態（Client→Server）
	MsgGoldenTimeStart  MessageType = "golden_time_start"  // 黃金時間開始廣播（Server→Client）
	MsgGoldenTimeEnd    MessageType = "golden_time_end"    // 黃金時間結束廣播（Server→Client）
	MsgGoldenTimeStatus MessageType = "golden_time_status" // 黃金時間狀態回應（Server→Client，個人）
	// 稀有連擊累積倍率系統（DAY-126）
	MsgRareCatchUpdate    MessageType = "rare_catch_update"    // 稀有連擊更新（Server→Client，個人）
	MsgRareCatchBroadcast MessageType = "rare_catch_broadcast" // 稀有連擊廣播（Server→Client，全服）
	MsgRareCatchReset     MessageType = "rare_catch_reset"     // 稀有連擊重置（Server→Client，個人）

	// 天氣湧現事件（DAY-127）
	MsgWeatherSurgeStart MessageType = "weather_surge_start" // 天氣湧現開始（Server→Client，全服廣播）
	MsgWeatherSurgeEnd   MessageType = "weather_surge_end"   // 天氣湧現結束（Server→Client，全服廣播）

	// 龍怒蓄力大招系統（DAY-128）
	MsgWrathUpdate MessageType = "wrath_update" // 怒氣值更新（Server→Client，個人）
	MsgWrathStart  MessageType = "wrath_start"  // 大招開始（Server→Client，全服廣播）
	MsgWrathResult MessageType = "wrath_result" // 大招結果（Server→Client，全服廣播）
	MsgUseWrath    MessageType = "use_wrath"    // 釋放大招（Client→Server）

	// 不死 BOSS 連勝系統（DAY-129）
	MsgImmortalBossSpawn  MessageType = "immortal_boss_spawn"  // 不死 BOSS 出現（Server→Client，全服廣播）
	MsgImmortalBossHit    MessageType = "immortal_boss_hit"    // 命中不死 BOSS（Server→Client，全服廣播）
	MsgImmortalBossLeave  MessageType = "immortal_boss_leave"  // 不死 BOSS 離開（Server→Client，全服廣播）
	MsgImmortalBossStatus MessageType = "immortal_boss_status" // 不死 BOSS 狀態（Server→Client，個人）

	MsgError            MessageType = "error"
	MsgPong             MessageType = "pong"
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
	// DAY-118：加入門檻資訊，讓 Client 能計算進度比例
	MiniThreshold  int `json:"mini_threshold"`  // Mini 觸發門檻
	MinorThreshold int `json:"minor_threshold"` // Minor 觸發門檻
	MajorThreshold int `json:"major_threshold"` // Major 觸發門檻
	GrandThreshold int `json:"grand_threshold"` // Grand 觸發門檻
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

// ---- 多格式每日賽系統（DAY-111）----

// MultiFormatRankEntry 多格式賽排名單筆記錄
type MultiFormatRankEntry struct {
	Rank        int     `json:"rank"`
	PlayerID    string  `json:"player_id"`
	DisplayName string  `json:"display_name"`
	Score       float64 `json:"score"`
	ScoreLabel  string  `json:"score_label"`  // 格式化後的分數（如 "50x"、"5000"）
	Prize       int     `json:"prize"`
	PrizeLabel  string  `json:"prize_label"`
	IsSelf      bool    `json:"is_self"`
}

// MultiFormatUpdatePayload 多格式賽排名更新廣播（每 30 秒）
type MultiFormatUpdatePayload struct {
	DayStart      int64                  `json:"day_start"`       // Unix ms
	DayEnd        int64                  `json:"day_end"`         // Unix ms
	SecondsLeft   int64                  `json:"seconds_left"`    // 距離結束秒數
	TodayFormat   string                 `json:"today_format"`    // "score"/"multiplier"/"reward"/"bet"
	FormatName    string                 `json:"format_name"`     // "積分賽"/"最高倍率賽"等
	FormatIcon    string                 `json:"format_icon"`     // "⭐"/"⚡"/"💰"/"🎯"
	FormatUnit    string                 `json:"format_unit"`     // "分"/"x"/"金幣"
	FormatDesc    string                 `json:"format_desc"`     // 格式說明
	Rankings      []MultiFormatRankEntry `json:"rankings"`        // 前 10 名
	TotalPlayers  int                    `json:"total_players"`   // 今日參賽人數
	PlayerRank    int                    `json:"player_rank"`     // 接收者的排名（0=未上榜）
	PlayerScore   float64                `json:"player_score"`    // 接收者的分數
	NextFormat    string                 `json:"next_format"`     // 明日格式
	NextFormatName string                `json:"next_format_name"` // 明日格式名稱
	NextFormatIcon string                `json:"next_format_icon"` // 明日格式圖示
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

// ---- 好友禮物系統（DAY-101）----

// SendGiftPayload 送禮物請求（Client → Server）
type SendGiftPayload struct {
	FriendID string `json:"friend_id"`
}

// GetGiftStatusPayload 查詢禮物狀態請求（Client → Server）
type GetGiftStatusPayload struct{}

// GiftReceivedPayload 收到禮物通知（Server → Client）
type GiftReceivedPayload struct {
	FromID      string `json:"from_id"`
	DisplayName string `json:"display_name"`
	Amount      int    `json:"amount"`
	NewBalance  int    `json:"new_balance"`
}

// GiftSentPayload 送出禮物成功通知（Server → Client）
type GiftSentPayload struct {
	ToID        string `json:"to_id"`
	DisplayName string `json:"display_name"`
	Amount      int    `json:"amount"`
	SentToday   int    `json:"sent_today"`   // 今日已送次數
	Remaining   int    `json:"remaining"`    // 今日剩餘次數
}

// GiftStatusPayload 今日禮物狀態（Server → Client）
type GiftStatusPayload struct {
	SentToday int `json:"sent_today"`
	Remaining int `json:"remaining"`
	MaxDaily  int `json:"max_daily"`
	Amount    int `json:"amount"` // 每次禮物金幣數
}

// GiftErrorPayload 禮物操作失敗（Server → Client）
type GiftErrorPayload struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

// ---- 好友挑戰系統（DAY-102）----

// SendChallengeRequestPayload 發起挑戰請求（Client → Server）
type SendChallengeRequestPayload struct {
	FriendID string `json:"friend_id"`
}

// AcceptChallengePayload 接受挑戰請求（Client → Server）
type AcceptChallengePayload struct {
	ChallengeID string `json:"challenge_id"`
}

// DeclineChallengePayload 拒絕挑戰請求（Client → Server）
type DeclineChallengePayload struct {
	ChallengeID string `json:"challenge_id"`
}

// ChallengeRequestPayload 收到挑戰通知（Server → Client）
type ChallengeRequestPayload struct {
	ChallengeID    string `json:"challenge_id"`
	ChallengerID   string `json:"challenger_id"`
	ChallengerName string `json:"challenger_name"`
	Stake          int    `json:"stake"`
	ExpiresInSec   int    `json:"expires_in_sec"`
}

// ChallengeUpdatePayload 挑戰狀態/分數更新（Server → Client）
type ChallengeUpdatePayload struct {
	ChallengeID   string `json:"challenge_id"`
	Status        string `json:"status"`
	OpponentID    string `json:"opponent_id"`
	OpponentName  string `json:"opponent_name,omitempty"`
	Stake         int    `json:"stake,omitempty"`
	MyScore       int    `json:"my_score"`
	OpponentScore int    `json:"opponent_score"`
	TimeRemaining int    `json:"time_remaining"` // 剩餘秒數
}

// ChallengeResultPayload 挑戰結果通知（Server → Client）
type ChallengeResultPayload struct {
	ChallengeID   string `json:"challenge_id"`
	IsWinner      bool   `json:"is_winner"`
	IsDraw        bool   `json:"is_draw"`
	WinnerName    string `json:"winner_name"`
	MyScore       int    `json:"my_score"`
	OpponentScore int    `json:"opponent_score"`
	OpponentID    string `json:"opponent_id"`
	OpponentName  string `json:"opponent_name"`
	Prize         int    `json:"prize"` // 獲得的金幣（勝者=2x賭注，平局=退回賭注）
}

// ChallengeErrorPayload 挑戰操作失敗（Server → Client）
type ChallengeErrorPayload struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

// ---- 私訊系統（DAY-103）----

// SendDMPayload 發送私訊請求（Client → Server）
type SendDMPayload struct {
	ToID    string `json:"to_id"`
	Content string `json:"content"`
}

// DMReceivedPayload 收到私訊通知（Server → Client）
type DMReceivedPayload struct {
	MessageID string `json:"message_id"`
	FromID    string `json:"from_id"`
	FromName  string `json:"from_name"`
	Content   string `json:"content"`
	SentAt    int64  `json:"sent_at"` // Unix milliseconds
	IsOffline bool   `json:"is_offline,omitempty"` // 是否為離線訊息
}

// DMSentPayload 發送成功確認（Server → Client）
type DMSentPayload struct {
	MessageID string `json:"message_id"`
	ToID      string `json:"to_id"`
	SentToday int    `json:"sent_today"`
	Remaining int    `json:"remaining"`
}

// DMErrorPayload 發送失敗（Server → Client）
type DMErrorPayload struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
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
	MercyUsed  bool   `json:"mercy_used"`  // 是否使用了寬限期（DAY-120）
	MercyLeft  int    `json:"mercy_left"`  // 本週剩餘寬限次數（DAY-120）
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

// ---- 玩家統計系統（DAY-096）----

// PlayerStatsPayload 玩家個人統計（Server → Client）
type PlayerStatsPayload struct {
	TotalSessions      int     `json:"total_sessions"`
	TotalPlayTimeSec   int64   `json:"total_play_time_sec"`
	TotalShots         int     `json:"total_shots"`
	TotalKills         int     `json:"total_kills"`
	TotalBet           int     `json:"total_bet"`
	TotalReward        int     `json:"total_reward"`
	TotalBonuses       int     `json:"total_bonuses"`
	TotalBossKills     int     `json:"total_boss_kills"`
	BestMultiplier     float64 `json:"best_multiplier"`
	BestStreak         int     `json:"best_streak"`
	BestSessionScore   int     `json:"best_session_score"`
	BestBonusReward    int     `json:"best_bonus_reward"`
	MaxCoins           int     `json:"max_coins"`
	JackpotWins        int     `json:"jackpot_wins"`
	JackpotMiniWins    int     `json:"jackpot_mini_wins"`
	JackpotMinorWins   int     `json:"jackpot_minor_wins"`
	JackpotMajorWins   int     `json:"jackpot_major_wins"`
	JackpotGrandWins   int     `json:"jackpot_grand_wins"`
	TotalJackpotPayout int     `json:"total_jackpot_payout"`
	HitRate            float64 `json:"hit_rate"`
	RTP                float64 `json:"rtp"`
	FirstPlayAtMs      int64   `json:"first_play_at_ms"`
	LastPlayAtMs       int64   `json:"last_play_at_ms"`
}

// ---- 全服公告系統（DAY-097）----

// AnnouncementPayload 全服公告廣播（Server → Client）
type AnnouncementPayload struct {
	ID         string `json:"id"`
	EventType  string `json:"event_type"`  // "jackpot_win" / "big_win" / "boss_kill" 等
	Priority   int    `json:"priority"`    // 1=低 2=普通 3=高 4=最高
	Title      string `json:"title"`       // 公告標題
	Message    string `json:"message"`     // 公告內容
	PlayerName string `json:"player_name"` // 相關玩家名稱
	Amount     int    `json:"amount"`      // 相關金額
	Icon       string `json:"icon"`        // 顯示圖示
	Color      string `json:"color"`       // 顯示顏色（hex）
	Duration   int    `json:"duration"`    // 顯示時長（毫秒）
	CreatedAtMs int64 `json:"created_at_ms"`
}

// ---- 玩家名片系統（DAY-106）----

// GetPlayerCardPayload 查詢玩家名片請求（Client → Server）
type GetPlayerCardPayload struct {
	TargetPlayerID string `json:"target_player_id"`
}

// PlayerCardPayload 玩家名片資料（Server → Client）
type PlayerCardPayload struct {
	PlayerID         string  `json:"player_id"`
	DisplayName      string  `json:"display_name"`
	TitleName        string  `json:"title_name"`
	TitleIcon        string  `json:"title_icon"`
	TitleColor       string  `json:"title_color"`
	VIPLevel         int     `json:"vip_level"`
	VIPName          string  `json:"vip_name"`
	GuildName        string  `json:"guild_name"`
	GuildRole        string  `json:"guild_role"`
	KillCount        int     `json:"kill_count"`
	MaxCoins         int     `json:"max_coins"`
	BestStreak       int     `json:"best_streak"`
	BestMult         float64 `json:"best_mult"`
	JackpotWins      int     `json:"jackpot_wins"`
	AchievementCount int     `json:"achievement_count"`
	LoginStreak      int     `json:"login_streak"`
	RTP              float64 `json:"rtp"`
	IsOnline         bool    `json:"is_online"`
}

// ---- 登入里程碑系統（DAY-107）----

// MilestoneRewardPayload 里程碑獎勵項目（用於通知 Client）
type MilestoneRewardPayload struct {
	Type   string `json:"type"`   // "coins" / "mystery_box" / "title"
	Amount int    `json:"amount"` // 金幣數量 / 寶箱數量
	Rarity string `json:"rarity"` // 寶箱稀有度（type=mystery_box 時）
	TitleID string `json:"title_id"` // 稱號 ID（type=title 時）
}

// LoginMilestonePayload 登入里程碑達成通知（Server → Client）
type LoginMilestonePayload struct {
	Days        int                      `json:"days"`        // 達到的里程碑天數
	Name        string                   `json:"name"`        // 里程碑名稱
	Description string                   `json:"description"` // 描述
	Icon        string                   `json:"icon"`        // 圖示 emoji
	Color       string                   `json:"color"`       // 顏色（hex）
	Rewards     []MilestoneRewardPayload `json:"rewards"`     // 獎勵列表
	CoinsGained int                      `json:"coins_gained"` // 本次獲得金幣總計
	NewBalance  int                      `json:"new_balance"`  // 領取後餘額
}

// MilestoneInfoPayload 里程碑資訊（用於 LoginProgressPayload）
type MilestoneInfoPayload struct {
	Days        int                      `json:"days"`
	Name        string                   `json:"name"`
	Icon        string                   `json:"icon"`
	Color       string                   `json:"color"`
	Rewards     []MilestoneRewardPayload `json:"rewards"`
	IsReached   bool                     `json:"is_reached"` // 是否已達到
}

// LoginProgressPayload 登入進度回應（Server → Client）
type LoginProgressPayload struct {
	CurrentStreak  int                    `json:"current_streak"`  // 當前連續天數
	MaxStreak      int                    `json:"max_streak"`      // 歷史最高
	NextMilestoneDays int                 `json:"next_milestone_days"` // 下一個里程碑天數（0=已全達成）
	DaysToNext     int                    `json:"days_to_next"`    // 距離下一個里程碑還差幾天
	Milestones     []MilestoneInfoPayload `json:"milestones"`      // 所有里程碑狀態
}

// ---- 超級 Bonus 系統（DAY-108）----

// SuperBonusReadyPayload 超級 Bonus 準備通知（Server → Client）
// 當玩家連續觸發 3 次 Bonus 後，第 4 次觸發超級 Bonus
type SuperBonusReadyPayload struct {
	ComboCount  int    `json:"combo_count"`  // 當前連續次數
	MultBonus   float64 `json:"mult_bonus"`  // 倍率加成（1.5x / 2.0x / 3.0x）
	Label       string `json:"label"`        // 顯示文字（"SUPER BONUS!" / "MEGA BONUS!" / "ULTRA BONUS!"）
	Color       string `json:"color"`        // 顏色（hex）
}

// ---- 賽季節日活動系統（DAY-109）----

// ClaimFestivalTaskPayload 領取節日任務獎勵請求（Client → Server）
type ClaimFestivalTaskPayload struct {
	TaskID string `json:"task_id"`
}

// FestivalTaskReadyPayload 節日任務可領取通知（Server → Client）
type FestivalTaskReadyPayload struct {
	TaskID string `json:"task_id"`
}

// FestivalTaskClaimedPayload 節日任務獎勵領取成功（Server → Client）
type FestivalTaskClaimedPayload struct {
	TaskID      string `json:"task_id"`
	RewardCoins int    `json:"reward_coins"`
}

// FestivalTitleEarnedPayload 節日稱號獲得通知（Server → Client）
type FestivalTitleEarnedPayload struct {
	TitleID    string `json:"title_id"`
	TitleName  string `json:"title_name"`
	TitleColor string `json:"title_color"`
}

// FestivalErrorPayload 節日操作失敗（Server → Client）
type FestivalErrorPayload struct {
	Message string `json:"message"`
}

// ---- 名人堂系統（DAY-110）----

// HallEntryPayload 名人堂條目（Server → Client）
type HallEntryPayload struct {
	PlayerID    string  `json:"player_id"`
	DisplayName string  `json:"display_name"`
	RecordType  string  `json:"record_type"`
	RecordLabel string  `json:"record_label"` // 中文標籤
	RecordIcon  string  `json:"record_icon"`  // 圖示
	Value       float64 `json:"value"`
	Description string  `json:"description"`
	AchievedAt  int64   `json:"achieved_at_ms"`
	BetLevel    int     `json:"bet_level"`
	CharacterID int     `json:"character_id"`
}

// HallOfFameUpdatePayload 名人堂更新（Server → Client）
type HallOfFameUpdatePayload struct {
	Records   []HallEntryPayload `json:"records"`
	UpdatedAt int64              `json:"updated_at_ms"`
}

// HallOfFameNewRecordPayload 新記錄誕生廣播（Server → Client）
type HallOfFameNewRecordPayload struct {
	Entry       HallEntryPayload `json:"entry"`
	OldHolder   string           `json:"old_holder"`   // 舊記錄持有者名稱（空=首次記錄）
	OldValue    float64          `json:"old_value"`    // 舊記錄數值（0=首次記錄）
	IsFirstTime bool             `json:"is_first_time"` // 是否是首次建立此記錄
}

// ---- 智慧推薦系統（DAY-110）----

// RecommendationPayload 單條推薦（Server → Client）
type RecommendationPayload struct {
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	Priority    int     `json:"priority"`
	TargetBetLv int     `json:"target_bet_lv,omitempty"`
	Confidence  float64 `json:"confidence"`
}

// RecommendationsPayload 推薦結果（Server → Client）
type RecommendationsPayload struct {
	Recommendations []RecommendationPayload `json:"recommendations"`
	GeneratedAt     int64                   `json:"generated_at_ms"`
}

// ---- 成就動態牆系統（DAY-112）----

// ActivityFeedEventPayload 動態牆事件廣播（Server → Client）
// 當有重要事件發生時廣播給所有玩家
type ActivityFeedEventPayload struct {
	ID          string `json:"id"`           // 唯一事件 ID
	EventType   string `json:"event_type"`   // "achievement"/"title"/"jackpot"/"boss_kill"/"mega_win"/"streak_record"/"hall_of_fame"/"season_level"/"milestone"
	PlayerID    string `json:"player_id"`    // 玩家 ID
	DisplayName string `json:"display_name"` // 玩家顯示名稱
	Icon        string `json:"icon"`         // 事件圖示（emoji）
	Title       string `json:"title"`        // 事件標題（如「Player1 解鎖成就」）
	Detail      string `json:"detail"`       // 事件詳情（如「討伐傳說」）
	Rarity      string `json:"rarity"`       // "common"/"uncommon"/"rare"/"epic"/"legendary"
	Timestamp   int64  `json:"timestamp"`    // Unix ms
}

// ActivityFeedHistoryPayload 歷史動態回應（Server → Client）
// 玩家上線時或主動查詢時發送最近 10 條
type ActivityFeedHistoryPayload struct {
	Events []ActivityFeedEventPayload `json:"events"`
	Total  int                        `json:"total"` // 總事件數（最多 50）
}

// ---- 雙層倍率輪盤系統（DAY-113）----

// RouletteSegmentPayload 輪盤格子定義（供 Client 顯示用）
type RouletteSegmentPayload struct {
	Multiplier float64 `json:"multiplier"`
	Label      string  `json:"label"`
	Color      string  `json:"color"`
}

// RouletteSpinPayload 單次旋轉結果
type RouletteSpinPayload struct {
	SegmentIndex int     `json:"segment_index"` // 停在哪個格子（Client 動畫用）
	Multiplier   float64 `json:"multiplier"`
	Label        string  `json:"label"`
	Color        string  `json:"color"`
}

// RouletteStartPayload 輪盤開始廣播（Server → Client）
// 廣播給所有玩家，讓大家都能看到輪盤動畫
type RouletteStartPayload struct {
	SessionID      string                   `json:"session_id"`
	PlayerID       string                   `json:"player_id"`
	PlayerName     string                   `json:"player_name"`
	TargetDefID    string                   `json:"target_def_id"`
	TargetName     string                   `json:"target_name"`
	BaseReward     int                      `json:"base_reward"`
	InnerSegments  []RouletteSegmentPayload `json:"inner_segments"`  // 內圈格子定義
	OuterSegments  []RouletteSegmentPayload `json:"outer_segments"`  // 外圈格子定義
	SpinDurationMs int                      `json:"spin_duration_ms"` // 旋轉動畫時長（ms）
}

// RouletteResultPayload 輪盤結果廣播（Server → Client）
// 廣播給所有玩家，顯示最終結果
type RouletteResultPayload struct {
	SessionID   string              `json:"session_id"`
	PlayerID    string              `json:"player_id"`
	PlayerName  string              `json:"player_name"`
	Inner       RouletteSpinPayload `json:"inner"`        // 內圈結果
	Outer       RouletteSpinPayload `json:"outer"`        // 外圈結果
	FinalMult   float64             `json:"final_mult"`   // 最終倍率（內圈 × 外圈）
	BaseReward  int                 `json:"base_reward"`
	FinalReward int                 `json:"final_reward"`
	NewBalance  int                 `json:"new_balance"`  // 中獎後餘額（只對觸發玩家有意義）
	IsJackpot   bool                `json:"is_jackpot"`   // ≥500x，觸發全畫面特效
	IsMegaWin   bool                `json:"is_mega_win"`  // ≥100x，觸發大獎特效
	IsSelf      bool                `json:"is_self"`      // 是否為自己觸發（Client 端標記）
}

// ---- Buy Bonus 系統（DAY-114）----

// BuyBonusPayload 購買 Bonus 請求（Client → Server）
type BuyBonusPayload struct {
	BonusType string `json:"bonus_type"` // "standard"（標準）或 "tnt"（TNT，更貴更強）
}

// BuyBonusSuccessPayload 購買成功通知（Server → Client）
type BuyBonusSuccessPayload struct {
	BonusType   string `json:"bonus_type"`   // "standard" / "tnt"
	Cost        int    `json:"cost"`         // 實際扣除金幣
	NewBalance  int    `json:"new_balance"`  // 購買後餘額
	DailyLeft   int    `json:"daily_left"`   // 今日剩餘購買次數
	MultBonus   float64 `json:"mult_bonus"`  // 本次 Bonus 的倍率加成（TNT=1.5x）
}

// BuyBonusErrorPayload 購買失敗通知（Server → Client）
type BuyBonusErrorPayload struct {
	Reason  string `json:"reason"`  // "insufficient_coins" / "daily_limit" / "game_busy"
	Message string `json:"message"` // 顯示給玩家的訊息
	Cost    int    `json:"cost"`    // 需要的金幣數
	Balance int    `json:"balance"` // 當前餘額
}

// BuyBonusStatusPayload 今日購買狀態（Server → Client）
type BuyBonusStatusPayload struct {
	DailyLimit    int  `json:"daily_limit"`    // 每日上限（3次）
	DailyUsed     int  `json:"daily_used"`     // 今日已購買次數
	DailyLeft     int  `json:"daily_left"`     // 今日剩餘次數
	StandardCost  int  `json:"standard_cost"`  // 標準 Bonus 費用
	TNTCost       int  `json:"tnt_cost"`       // TNT Bonus 費用
	CanBuy        bool `json:"can_buy"`        // 是否可以購買（未達上限且遊戲狀態正常）
}

// ---- 新手引導系統（DAY-115）----

// TutorialStepPayload 引導步驟（Server → Client）
type TutorialStepPayload struct {
	Step      int     `json:"step"`       // 當前步驟（1-4）
	TotalStep int     `json:"total_step"` // 總步驟數（3）
	Title     string  `json:"title"`      // 步驟標題
	Desc      string  `json:"desc"`       // 步驟說明
	Highlight string  `json:"highlight"`  // 高亮區域（"game_area"/"bet_buttons"/"labor_bar"/"none"）
	Action    string  `json:"action"`     // 期望動作（"shoot"/"bet_change"/"bonus"/"complete"）
	ArrowX    float64 `json:"arrow_x"`    // 引導箭頭 X 座標
	ArrowY    float64 `json:"arrow_y"`    // 引導箭頭 Y 座標
}

// TutorialActionPayload 玩家完成引導步驟（Client → Server）
type TutorialActionPayload struct {
	Action string `json:"action"` // "shoot_done"/"bet_done"/"bonus_done"/"skip"
}

// ---- Co-op Boss Raid 系統（DAY-115）----

// RaidContributorPayload 貢獻者資料
type RaidContributorPayload struct {
	PlayerID    string `json:"player_id"`
	DisplayName string `json:"display_name"`
	Damage      int    `json:"damage"`
	Reward      int    `json:"reward"`   // 結算後才有值
	Rank        int    `json:"rank"`
}

// RaidWarningPayload 討伐警告廣播（Server → Client）
type RaidWarningPayload struct {
	RaidID    string `json:"raid_id"`
	BossName  string `json:"boss_name"`
	MaxHP     int    `json:"max_hp"`
	RewardPool int   `json:"reward_pool"`
	StartsIn  int    `json:"starts_in"` // 幾秒後開始（30）
}

// RaidStartPayload 討伐開始廣播（Server → Client）
type RaidStartPayload struct {
	RaidID     string `json:"raid_id"`
	BossName   string `json:"boss_name"`
	HP         int    `json:"hp"`
	MaxHP      int    `json:"max_hp"`
	RewardPool int    `json:"reward_pool"`
	Duration   int    `json:"duration"` // 討伐持續秒數（300）
}

// RaidUpdatePayload 討伐狀態更新廣播（Server → Client，每 3 秒）
type RaidUpdatePayload struct {
	RaidID       string                   `json:"raid_id"`
	HP           int                      `json:"hp"`
	MaxHP        int                      `json:"max_hp"`
	TimeLeft     float64                  `json:"time_left"`
	Contributors []*RaidContributorPayload `json:"contributors"`
}

// RaidResultPayload 討伐結算廣播（Server → Client）
type RaidResultPayload struct {
	RaidID       string                   `json:"raid_id"`
	BossName     string                   `json:"boss_name"`
	Defeated     bool                     `json:"defeated"`     // true=擊殺，false=超時
	RewardPool   int                      `json:"reward_pool"`
	Contributors []*RaidContributorPayload `json:"contributors"` // 含個人獎勵
	MyReward     int                      `json:"my_reward"`    // 本玩家獎勵（0=未參與）
	MyRank       int                      `json:"my_rank"`      // 本玩家排名（0=未參與）
}

// RaidStatusPayload 討伐狀態回應（Server → Client）
type RaidStatusPayload struct {
	State        string                   `json:"state"`         // "idle"/"warning"/"active"/"result"
	RaidID       string                   `json:"raid_id"`
	BossName     string                   `json:"boss_name"`
	HP           int                      `json:"hp"`
	MaxHP        int                      `json:"max_hp"`
	RewardPool   int                      `json:"reward_pool"`
	TimeLeft     float64                  `json:"time_left"`
	Contributors []*RaidContributorPayload `json:"contributors"`
	CanTrigger   bool                     `json:"can_trigger"` // 今日是否可觸發
}

// ---- 碎片收集大獎系統（DAY-116）----

// FragmentDropPayload 碎片掉落通知（Server → Client）
// 擊破目標後有機率掉落碎片，飛向收集槽
type FragmentDropPayload struct {
	FragmentType string `json:"fragment_type"` // "bronze"/"silver"/"gold"
	Label        string `json:"label"`         // 顯示名稱（如「銅碎片」）
	Color        string `json:"color"`         // 顯示顏色（hex）
	NewCount     int    `json:"new_count"`     // 目前持有數量
	Required     int    `json:"required"`      // 集齊需要數量
	DropX        float64 `json:"drop_x"`      // 掉落位置 X（動畫起點）
	DropY        float64 `json:"drop_y"`      // 掉落位置 Y（動畫起點）
}

// FragmentCompletePayload 集齊碎片大獎通知（Server → Client，廣播給所有玩家）
type FragmentCompletePayload struct {
	PlayerID     string `json:"player_id"`
	DisplayName  string `json:"display_name"`
	FragmentType string `json:"fragment_type"` // "bronze"/"silver"/"gold"
	Label        string `json:"label"`         // 大獎名稱（如「金碎片大獎」）
	Color        string `json:"color"`         // 顯示顏色
	Reward       int    `json:"reward"`        // 獲得金幣
	NewBalance   int    `json:"new_balance"`   // 中獎後餘額（只對觸發玩家有意義）
	IsSelf       bool   `json:"is_self"`       // 是否為自己觸發
}

// FragmentStatusPayload 碎片狀態回應（Server → Client）
type FragmentStatusPayload struct {
	Bronze   int `json:"bronze"`    // 目前銅碎片數量
	Silver   int `json:"silver"`    // 目前銀碎片數量
	Gold     int `json:"gold"`      // 目前金碎片數量
	Required int `json:"required"`  // 集齊需要數量（固定 5）
}

// LuckyCatchPayload 幸運捕獲廣播（Server → Client，廣播）（DAY-119）
type LuckyCatchPayload struct {
	PlayerID    string  `json:"player_id"`    // 觸發玩家 ID
	PlayerName  string  `json:"player_name"`  // 觸發玩家名稱
	TargetDefID string  `json:"target_def_id"` // 被捕獲的目標物定義 ID
	TargetName  string  `json:"target_name"`  // 目標物名稱
	Multiplier  float64 `json:"multiplier"`   // 目標物倍率
	BonusMult   float64 `json:"bonus_mult"`   // 幸運加成倍率（2.0-5.0x）
	Reward      int     `json:"reward"`       // 最終獎勵金幣
	TriggerType string  `json:"trigger_type"` // 觸發類型：streak/weather/festival
	Icon        string  `json:"icon"`         // 顯示圖示
}

// MissionMercyProtectedPayload 任務連續寬限期保護通知（Server → Client）（DAY-120）
type MissionMercyProtectedPayload struct {
	Streak    int    `json:"streak"`     // 被保護的連續天數
	MercyLeft int    `json:"mercy_left"` // 本週剩餘寬限次數（使用後）
	Message   string `json:"message"`    // 顯示訊息
}

// RapidRespinPayload Rapid Respin 觸發廣播（Server → Client，廣播）（DAY-121）
type RapidRespinPayload struct {
	PlayerID    string  `json:"player_id"`    // 觸發玩家 ID
	PlayerName  string  `json:"player_name"`  // 觸發玩家名稱
	ChainCount  int     `json:"chain_count"`  // 當前連鎖次數（0=第一次，1=第二次...）
	ChainMult   float64 `json:"chain_mult"`   // 當前連鎖倍率（1.0/1.5/2.0/3.0/5.0）
	IsChain     bool    `json:"is_chain"`     // 是否為連鎖觸發
	MaxChain    int     `json:"max_chain"`    // 最大連鎖次數（5）
	Icon        string  `json:"icon"`         // 顯示圖示（⚡🔄）
}

// RapidRespinEndPayload Rapid Respin 連鎖結束通知（Server → Client）（DAY-121）
type RapidRespinEndPayload struct {
	PlayerID   string `json:"player_id"`   // 觸發玩家 ID
	PlayerName string `json:"player_name"` // 觸發玩家名稱
	TotalChain int    `json:"total_chain"` // 總連鎖次數
}

// ---- 寶藏地圖系統（DAY-122）----

// TreasureMapCellPayload 地圖格子狀態
type TreasureMapCellPayload struct {
	Row    int    `json:"row"`
	Col    int    `json:"col"`
	DefID  string `json:"def_id"`  // 對應目標物 ID
	Name   string `json:"name"`    // 顯示名稱
	Icon   string `json:"icon"`    // 顯示圖示
	Filled bool   `json:"filled"`  // 是否已填滿
}

// TreasureMapUpdatePayload 寶藏地圖狀態更新（Server → Client）（DAY-122）
type TreasureMapUpdatePayload struct {
	Cells       []TreasureMapCellPayload `json:"cells"`        // 所有格子狀態
	FilledCount int                      `json:"filled_count"` // 已填滿格子數
	LinesCount  int                      `json:"lines_count"`  // 已完成行/列/對角線數
	FullDone    bool                     `json:"full_done"`    // 是否已完成整張地圖
	Date        string                   `json:"date"`         // 今日日期（YYYY-MM-DD）
}

// TreasureMapLinePayload 完成一行/列/對角線通知（Server → Client）（DAY-122）
type TreasureMapLinePayload struct {
	LineType   string `json:"line_type"`  // "row0"/"row1"/"row2"/"col0"/"col1"/"col2"/"diag0"/"diag1"
	Reward     int    `json:"reward"`     // 獎勵金幣
	NewBalance int    `json:"new_balance"` // 新餘額
	Message    string `json:"message"`    // 顯示訊息
}

// TreasureMapFullPayload 完成整張地圖通知（Server → Client）（DAY-122）
type TreasureMapFullPayload struct {
	Reward     int    `json:"reward"`      // 傳說寶藏獎勵金幣
	NewBalance int    `json:"new_balance"` // 新餘額
	Message    string `json:"message"`     // 顯示訊息
}

// ---- 閃電挑戰系統（DAY-123）----

// FlashChallengePlayerSnap 玩家進度快照
type FlashChallengePlayerSnap struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	Progress   int    `json:"progress"`
	Completed  bool   `json:"completed"`
}

// FlashChallengeStartPayload 閃電挑戰開始廣播（Server → Client）（DAY-123）
type FlashChallengeStartPayload struct {
	Type        string                      `json:"type"`
	Title       string                      `json:"title"`
	Description string                      `json:"description"`
	Icon        string                      `json:"icon"`
	Color       string                      `json:"color"`
	Target      int                         `json:"target"`
	TargetDefID string                      `json:"target_def_id"`
	Duration    int                         `json:"duration"`
	TimeLeft    int                         `json:"time_left"`
	BaseReward  int                         `json:"base_reward"`
	BonusReward int                         `json:"bonus_reward"`
	TopPlayers  []FlashChallengePlayerSnap  `json:"top_players"`
}

// FlashChallengeUpdatePayload 閃電挑戰進度更新（Server → Client）（DAY-123）
type FlashChallengeUpdatePayload struct {
	PlayerID   string                     `json:"player_id"`
	PlayerName string                     `json:"player_name"`
	Progress   int                        `json:"progress"`
	Target     int                        `json:"target"`
	Completed  bool                       `json:"completed"`
	TimeLeft   int                        `json:"time_left"`
	TopPlayers []FlashChallengePlayerSnap `json:"top_players"`
}

// FlashChallengeEndPayload 閃電挑戰結束廣播（Server → Client）（DAY-123）
type FlashChallengeEndPayload struct {
	Success    bool                       `json:"success"`    // 是否有人完成
	Title      string                     `json:"title"`
	Icon       string                     `json:"icon"`
	TopPlayers []FlashChallengePlayerSnap `json:"top_players"`
	Message    string                     `json:"message"`
}

// FlashChallengeRewardPayload 閃電挑戰獎勵通知（Server → Client，個人）（DAY-123）
type FlashChallengeRewardPayload struct {
	PlayerID   string `json:"player_id"`
	Progress   int    `json:"progress"`
	Target     int    `json:"target"`
	Completed  bool   `json:"completed"`
	Reward     int    `json:"reward"`
	NewBalance int    `json:"new_balance"`
	Message    string `json:"message"`
}

// FlashChallengeStatusPayload 閃電挑戰狀態（Server → Client）（DAY-123）
type FlashChallengeStatusPayload struct {
	Active      bool                       `json:"active"`
	Type        string                     `json:"type"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	Icon        string                     `json:"icon"`
	Color       string                     `json:"color"`
	Target      int                        `json:"target"`
	TargetDefID string                     `json:"target_def_id"`
	Duration    int                        `json:"duration"`
	TimeLeft    int                        `json:"time_left"`
	BaseReward  int                        `json:"base_reward"`
	BonusReward int                        `json:"bonus_reward"`
	MyProgress  int                        `json:"my_progress"`
	MyCompleted bool                       `json:"my_completed"`
	TopPlayers  []FlashChallengePlayerSnap `json:"top_players"`
}

// ---- 傳說目標警報系統（DAY-124）----

// RareTargetAlertPayload 稀有/傳說目標出現廣播（Server → Client）（DAY-124）
type RareTargetAlertPayload struct {
	InstanceID string `json:"instance_id"` // 目標實例 ID
	DefID      string `json:"def_id"`      // 目標定義 ID
	Name       string `json:"name"`        // 目標名稱
	Quality    string `json:"quality"`     // 品質等級（epic/legendary）
	Multiplier int    `json:"multiplier"`  // 倍率
	Icon       string `json:"icon"`        // 圖示（⭐/💜）
	Message    string `json:"message"`     // 顯示訊息
	Color      string `json:"color"`       // 顏色（hex）
}

// ---- 個人最佳記錄通知系統（DAY-125）----

// PersonalBestPayload 個人最佳記錄通知（Server → Client，個人）（DAY-125）
type PersonalBestPayload struct {
	RecordType  string  `json:"record_type"`  // 記錄類型：multiplier/streak/reward/coins
	OldValue    float64 `json:"old_value"`    // 舊記錄
	NewValue    float64 `json:"new_value"`    // 新記錄
	Label       string  `json:"label"`        // 顯示文字（如「最高倍率」）
	Icon        string  `json:"icon"`         // 圖示
	Message     string  `json:"message"`      // 完整訊息
}

// ---- 黃金時間系統（DAY-125）----

// GoldenTimeStartPayload 黃金時間開始廣播（Server → Client）（DAY-125）
type GoldenTimeStartPayload struct {
	Tier        int     `json:"tier"`         // 等級（0=Silver, 1=Gold, 2=Rainbow）
	TierName    string  `json:"tier_name"`    // 等級名稱（如「✨ 黃金時間」）
	MultBoost   float64 `json:"mult_boost"`   // 倍率加成（1.5/2.0/3.0）
	Duration    int     `json:"duration"`     // 持續秒數
	SecondsLeft int     `json:"seconds_left"` // 剩餘秒數（開始時等於 Duration）
	Icon        string  `json:"icon"`         // 圖示（⚡/✨/🌈）
	Color       string  `json:"color"`        // 主色（hex）
	BgColor     string  `json:"bg_color"`     // 背景色（hex）
	TriggerType string  `json:"trigger_type"` // 觸發類型（boss_kill/random/flash_combo/raid_victory）
	Message     string  `json:"message"`      // 廣播訊息
}

// GoldenTimeEndPayload 黃金時間結束廣播（Server → Client）（DAY-125）
type GoldenTimeEndPayload struct {
	Tier     int    `json:"tier"`      // 等級
	TierName string `json:"tier_name"` // 等級名稱
	Message  string `json:"message"`   // 結束訊息
}

// GoldenTimeStatusPayload 黃金時間狀態回應（Server → Client，個人）（DAY-125）
type GoldenTimeStatusPayload struct {
	IsActive    bool    `json:"is_active"`    // 是否進行中
	Tier        int     `json:"tier"`         // 等級
	TierName    string  `json:"tier_name"`    // 等級名稱
	MultBoost   float64 `json:"mult_boost"`   // 倍率加成
	SecondsLeft int     `json:"seconds_left"` // 剩餘秒數
	Icon        string  `json:"icon"`         // 圖示
	Color       string  `json:"color"`        // 主色
	BgColor     string  `json:"bg_color"`     // 背景色
	TriggerType string  `json:"trigger_type"` // 觸發類型
}

// ---- 稀有連擊累積倍率系統（DAY-126）----

// RareCatchUpdatePayload 稀有連擊更新（Server → Client，個人）（DAY-126）
type RareCatchUpdatePayload struct {
	Count       int     `json:"count"`        // 當前連擊數
	MultBoost   float64 `json:"mult_boost"`   // 當前倍率加成
	LevelName   string  `json:"level_name"`   // 等級名稱
	Icon        string  `json:"icon"`         // 圖示
	Color       string  `json:"color"`        // 顏色（hex）
	SecondsLeft int     `json:"seconds_left"` // 距離超時的剩餘秒數
	IsLevelUp   bool    `json:"is_level_up"`  // 是否升級
}

// RareCatchBroadcastPayload 稀有連擊廣播（Server → Client，全服）（DAY-126）
// 達到 ×5.0 以上時廣播，讓全場玩家知道有人在連擊稀有目標
type RareCatchBroadcastPayload struct {
	PlayerID   string  `json:"player_id"`   // 玩家 ID
	PlayerName string  `json:"player_name"` // 玩家名稱
	Count      int     `json:"count"`       // 連擊數
	MultBoost  float64 `json:"mult_boost"`  // 倍率加成
	LevelName  string  `json:"level_name"`  // 等級名稱
	Icon       string  `json:"icon"`        // 圖示
	Color      string  `json:"color"`       // 顏色
	Message    string  `json:"message"`     // 廣播訊息
}

// RareCatchResetPayload 稀有連擊重置（Server → Client，個人）（DAY-126）
type RareCatchResetPayload struct {
	FinalCount int    `json:"final_count"` // 最終連擊數
	Message    string `json:"message"`     // 重置訊息
}

// WeatherSurgeStartPayload 天氣湧現開始（Server → Client，全服廣播）（DAY-127）
type WeatherSurgeStartPayload struct {
	SurgeName    string  `json:"surge_name"`    // 湧現名稱（如「暴風湧現」）
	SurgeIcon    string  `json:"surge_icon"`    // 圖示（emoji）
	SurgeMessage string  `json:"surge_message"` // 廣播訊息
	Duration     int     `json:"duration"`      // 持續秒數
	RareBonus    float64 `json:"rare_bonus"`    // 稀有目標加成（0.0-1.0）
	GoldBonus    float64 `json:"gold_bonus"`    // 金幣魚加成（0.0-1.0）
	Color        string  `json:"color"`         // 橫幅顏色（hex）
}

// WeatherSurgeEndPayload 天氣湧現結束（Server → Client，全服廣播）（DAY-127）
type WeatherSurgeEndPayload struct {
	SurgeName string `json:"surge_name"` // 湧現名稱
	Message   string `json:"message"`    // 結束訊息
}

// WrathUpdatePayload 怒氣值更新（Server → Client，個人）（DAY-128）
type WrathUpdatePayload struct {
	Charge    int  `json:"charge"`     // 當前怒氣值（0-100）
	MaxCharge int  `json:"max_charge"` // 最大怒氣值（100）
	IsReady   bool `json:"is_ready"`   // 是否可以釋放大招
	Cooldown  int  `json:"cooldown"`   // 冷卻剩餘秒數（0 = 可用）
}

// WrathStartPayload 大招開始廣播（Server → Client，全服）（DAY-128）
type WrathStartPayload struct {
	PlayerID   string `json:"player_id"`   // 釋放大招的玩家 ID
	PlayerName string `json:"player_name"` // 玩家名稱
	Icon       string `json:"icon"`        // 圖示
	Message    string `json:"message"`     // 廣播訊息
}

// WrathKillEntry 大招擊破的目標（DAY-128）
type WrathKillEntry struct {
	InstanceID string  `json:"instance_id"` // 目標實例 ID
	DefID      string  `json:"def_id"`      // 目標定義 ID
	Reward     int     `json:"reward"`      // 獎勵金幣
	Multiplier float64 `json:"multiplier"`  // 目標倍率
}

// WrathResultPayload 大招結果（Server → Client，全服廣播）（DAY-128）
type WrathResultPayload struct {
	PlayerID    string          `json:"player_id"`    // 玩家 ID
	PlayerName  string          `json:"player_name"`  // 玩家名稱
	KilledCount int             `json:"killed_count"` // 擊破目標數
	TotalReward int             `json:"total_reward"` // 總獎勵
	NewBalance  int             `json:"new_balance"`  // 新金幣餘額
	Targets     []WrathKillEntry `json:"targets"`     // 擊破的目標列表
}

// ---- 不死 BOSS 連勝系統 Payloads（DAY-129）----

// ImmortalBossSpawnPayload 不死 BOSS 出現廣播（Server → Client，全服）（DAY-129）
type ImmortalBossSpawnPayload struct {
	InstanceID       string  `json:"instance_id"`        // BOSS 實例 ID
	BossType         string  `json:"boss_type"`          // BOSS 類型（golden_toad/ancient_crocodile）
	BossName         string  `json:"boss_name"`          // BOSS 名稱
	BossIcon         string  `json:"boss_icon"`          // BOSS 圖示
	BossColor        string  `json:"boss_color"`         // 顯示顏色
	MinMult          float64 `json:"min_mult"`           // 最小倍率
	MaxMult          float64 `json:"max_mult"`           // 最大倍率
	DurationSeconds  float64 `json:"duration_seconds"`   // 在場時間（秒）
	Message          string  `json:"message"`            // 公告訊息
}

// ImmortalBossHitPayload 命中不死 BOSS 廣播（Server → Client，全服）（DAY-129）
type ImmortalBossHitPayload struct {
	InstanceID  string  `json:"instance_id"`  // BOSS 實例 ID
	PlayerID    string  `json:"player_id"`    // 命中玩家 ID
	PlayerName  string  `json:"player_name"`  // 命中玩家名稱
	Multiplier  float64 `json:"multiplier"`   // 本次倍率
	Reward      int     `json:"reward"`       // 本次獎勵
	NewBalance  int     `json:"new_balance"`  // 玩家新餘額（只發給命中者）
	HitCount    int     `json:"hit_count"`    // 累計命中次數
	TotalReward int     `json:"total_reward"` // 累計總獎勵
	IsHighMult  bool    `json:"is_high_mult"` // 是否高倍率（≥100x）
}

// ImmortalBossLeavePayload 不死 BOSS 離開廣播（Server → Client，全服）（DAY-129）
type ImmortalBossLeavePayload struct {
	InstanceID  string `json:"instance_id"`  // BOSS 實例 ID
	BossName    string `json:"boss_name"`    // BOSS 名稱
	BossIcon    string `json:"boss_icon"`    // BOSS 圖示
	HitCount    int    `json:"hit_count"`    // 總命中次數
	TotalReward int    `json:"total_reward"` // 總獎勵
	Message     string `json:"message"`     // 離開訊息
}

// ImmortalBossStatusPayload 不死 BOSS 狀態（Server → Client，個人）（DAY-129）
type ImmortalBossStatusPayload struct {
	Active           bool    `json:"active"`            // 是否有活躍 BOSS
	InstanceID       string  `json:"instance_id"`       // BOSS 實例 ID
	BossType         string  `json:"boss_type"`         // BOSS 類型
	BossName         string  `json:"boss_name"`         // BOSS 名稱
	BossIcon         string  `json:"boss_icon"`         // BOSS 圖示
	BossColor        string  `json:"boss_color"`        // 顯示顏色
	MinMult          float64 `json:"min_mult"`          // 最小倍率
	MaxMult          float64 `json:"max_mult"`          // 最大倍率
	HitCount         int     `json:"hit_count"`         // 累計命中次數
	TotalReward      int     `json:"total_reward"`      // 累計總獎勵
	RemainingSeconds float64 `json:"remaining_seconds"` // 剩餘秒數
}
