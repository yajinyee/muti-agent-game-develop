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
	// 特殊武器系統（DAY-089，升級 DAY-134，DAY-154）
	MsgSpecialWeaponUpdate  MessageType = "special_weapon_update"  // 特殊武器狀態更新
	MsgSpecialWeaponFired   MessageType = "special_weapon_fired"   // 特殊武器發射廣播（所有玩家可見）
	MsgSpecialWeaponCharged MessageType = "special_weapon_charged" // 自動充能完成通知（DAY-134）
	MsgHomingMissileResult  MessageType = "homing_missile_result"  // 追蹤飛彈命中結果（DAY-141）
	MsgDragonWrathCharge    MessageType = "dragon_wrath_charge"    // 龍怒怒氣值更新（DAY-154）
	MsgDragonWrathResult    MessageType = "dragon_wrath_result"    // 龍怒流星雨結果（DAY-154）
	MsgTorpedoResult        MessageType = "torpedo_result"         // 魚雷爆炸結果（DAY-155）
	MsgRailgunResult        MessageType = "railgun_result"         // 軌道炮穿透結果（DAY-157）
	MsgBlackHoleResult      MessageType = "black_hole_result"      // 黑洞漩渦爆炸結果（DAY-166）
	MsgRouletteCrabStart    MessageType = "roulette_crab_start"    // 黃金輪盤螃蟹開始（DAY-167）
	MsgRouletteCrabStop     MessageType = "roulette_crab_stop"     // 黃金輪盤螃蟹停止（DAY-167）
	// 冰釣幸運輪盤系統（DAY-171）
	MsgIceFishingWheelStop MessageType = "ice_fishing_wheel_stop" // 冰釣輪盤停止（Client→Server，DAY-171）
	MsgRouletteCrabResult   MessageType = "roulette_crab_result"   // 黃金輪盤螃蟹結果（DAY-167）
	MsgRouletteCrabStatus   MessageType = "roulette_crab_status"   // 黃金輪盤螃蟹冷卻狀態（DAY-167）
	MsgGoldenTurtleTimeStop MessageType = "golden_turtle_time_stop" // 黃金海龜時間停止（DAY-159）
	MsgLuckyStarFish        MessageType = "lucky_star_fish"         // 幸運星魚全場倍率翻倍（DAY-160）
	MsgGoldenSharkBerserk   MessageType = "golden_shark_berserk"    // 黃金鯊魚全服狂暴模式（DAY-161）
	MsgMoneyFishReward      MessageType = "money_fish_reward"       // 金幣魚王即時獎勵（DAY-162）
	MsgCaptainFishRace      MessageType = "captain_fish_race"       // 船長魚全服競速模式（DAY-163）
	MsgAbyssWhale           MessageType = "abyss_whale"             // 深淵巨鯨全服 Boss 挑戰（DAY-164）
	MsgRoyalChainLightning  MessageType = "royal_chain_lightning"  // 皇家閃電鰻持續連鎖電擊（DAY-156）
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

	// 覺醒 BOSS 系統（DAY-130）
	MsgAwakenBossSpawn   MessageType = "awaken_boss_spawn"   // 覺醒 BOSS 出現（Server→Client，全服廣播）
	MsgAwakenBossHit     MessageType = "awaken_boss_hit"     // 命中覺醒 BOSS（Server→Client，全服廣播）
	MsgAwakenBossPowerUp MessageType = "awaken_boss_powerup" // Power Up 觸發（Server→Client，全服廣播）
	MsgAwakenBossLeave   MessageType = "awaken_boss_leave"   // 覺醒 BOSS 離開（Server→Client，全服廣播）
	MsgAwakenBossStatus  MessageType = "awaken_boss_status"  // 覺醒 BOSS 狀態（Server→Client，個人）

	// 連勝獎勵系統（DAY-131）
	MsgWinStreakUpdate    MessageType = "win_streak_update"    // 連勝更新（Server→Client，個人）
	MsgWinStreakMilestone MessageType = "win_streak_milestone" // 里程碑達成（Server→Client，個人/全服）
	MsgWinStreakReset     MessageType = "win_streak_reset"     // 連勝重置（Server→Client，個人）

	// 閃電鰻連鎖攻擊系統（DAY-132）
	MsgLightningEelChain  MessageType = "lightning_eel_chain"  // 連鎖攻擊結果廣播（Server→Client，全服）
	MsgLightningEelStatus MessageType = "lightning_eel_status" // 閃電鰻冷卻狀態（Server→Client，個人）

	// 狂熱模式系統（DAY-133）
	MsgFeverModeStart  MessageType = "fever_mode_start"  // 狂熱模式開始（Server→Client，全服廣播）
	MsgFeverModeEnd    MessageType = "fever_mode_end"    // 狂熱模式結束（Server→Client，個人）
	MsgFeverModeStatus MessageType = "fever_mode_status" // 狂熱模式狀態（Server→Client，個人）

	// 失敗補償系統（DAY-135）
	MsgUnluckyBonus       MessageType = "unlucky_bonus"        // 失敗補償觸發（Server→Client，個人）
	MsgUnluckyBonusStatus MessageType = "unlucky_bonus_status" // 失敗補償狀態（Server→Client，個人）

	// 競速獵殺系統（DAY-136）
	MsgSpeedRaceStart  MessageType = "speed_race_start"  // 競速開始廣播（Server→Client，全服）
	MsgSpeedRaceEnd    MessageType = "speed_race_end"    // 競速結束廣播（Server→Client，全服）
	MsgSpeedRaceCancel MessageType = "speed_race_cancel" // 競速取消廣播（Server→Client，全服）
	MsgSpeedRaceResult MessageType = "speed_race_result" // 競速個人結果（Server→Client，個人）

	// 全服目標懸賞系統（DAY-137）
	MsgPostBounty   MessageType = "post_bounty"    // 玩家下懸賞（Client→Server）
	MsgGetBounties  MessageType = "get_bounties"   // 查詢懸賞列表（Client→Server）
	MsgBountyPosted MessageType = "bounty_posted"  // 懸賞發布廣播（Server→Client，全服）
	MsgBountyClaimed MessageType = "bounty_claimed" // 懸賞領取通知（Server→Client，個人）
	MsgBountyKilled  MessageType = "bounty_killed"  // 懸賞目標擊破廣播（Server→Client，全服）
	MsgBountyExpired MessageType = "bounty_expired" // 懸賞過期通知（Server→Client，個人+全服）
	MsgBountyList    MessageType = "bounty_list"    // 懸賞列表回應（Server→Client，個人）
	MsgBountyError   MessageType = "bounty_error"   // 懸賞操作失敗（Server→Client，個人）

	// 全服倍率風暴系統（DAY-138）
	MsgMultStormStart MessageType = "mult_storm_start" // 風暴開始廣播（Server→Client，全服）
	MsgMultStormEnd   MessageType = "mult_storm_end"   // 風暴結束廣播（Server→Client，全服）

	// 雙環輪盤系統（DAY-139）
	MsgDualRouletteStart  MessageType = "dual_roulette_start"  // 輪盤開始（Server→Client，個人）
	MsgDualRouletteStop   MessageType = "dual_roulette_stop"   // 玩家停止輪盤（Client→Server）
	MsgDualRouletteResult MessageType = "dual_roulette_result" // 輪盤結果（Server→Client，個人）
	MsgDualRouletteStatus MessageType = "dual_roulette_status" // 輪盤狀態（Server→Client，個人）

	// 全服 Mega Catch 事件系統（DAY-140）
	MsgMegaCatchStart  MessageType = "mega_catch_start"  // 事件開始廣播（Server→Client，全服）
	MsgMegaCatchEnd    MessageType = "mega_catch_end"    // 事件結束廣播（Server→Client，全服）
	MsgMegaCatchStatus MessageType = "mega_catch_status" // 事件狀態（Server→Client，個人）

	// 鑽頭龍蝦連帶效果系統（DAY-142）
	MsgDrillLobsterChain MessageType = "drill_lobster_chain" // 鑽頭連帶效果廣播（Server→Client，全服）

	// 炸彈蟹連環爆炸系統（DAY-143）
	MsgBombCrabChain MessageType = "bomb_crab_chain" // 炸彈蟹連環爆炸廣播（Server→Client，全服）

	// 巨型章魚轉盤系統（DAY-144）
	MsgMegaOctopusWheelStart  MessageType = "mega_octopus_wheel_start"  // 轉盤開始（Server→Client，個人）
	MsgMegaOctopusWheelStop   MessageType = "mega_octopus_wheel_stop"   // 玩家停止轉盤（Client→Server）
	MsgMegaOctopusWheelResult MessageType = "mega_octopus_wheel_result" // 轉盤結果（Server→Client，個人）

	// 巨型鮟鱇魚電擊寶箱系統（DAY-145）
	MsgAnglerfishShock MessageType = "anglerfish_shock" // 電擊開寶箱廣播（Server→Client，全服）

	// 巨型鹹水鱷魚獵魚累積系統（DAY-146）
	MsgCrocodileHunt MessageType = "crocodile_hunt" // 鱷魚獵魚廣播（Server→Client，全服）

	// 夢幻巨型獎勵魚系統（DAY-147）
	MsgGiantPrizeFish MessageType = "giant_prize_fish" // 夢幻獎勵模式廣播（Server→Client，全服）

	// 千龍王強化輪盤系統（DAY-148）
	MsgChainLongWheelStart  MessageType = "chainlong_wheel_start"  // 輪盤開始（Server→Client，個人）
	MsgChainLongWheelStop   MessageType = "chainlong_wheel_stop"   // 玩家停止輪盤（Client→Server）
	MsgChainLongWheelResult MessageType = "chainlong_wheel_result" // 輪盤結果（Server→Client，個人+全服）
	MsgChainLongWheelStatus MessageType = "chainlong_wheel_status" // 冷卻狀態（Server→Client，個人）

	// 黃金水母全場電擊系統（DAY-149）
	MsgGoldenJellyfishShock MessageType = "golden_jellyfish_shock" // 全場電擊廣播（Server→Client，全服）

	// 雷霆龍蝦免費射擊系統（DAY-150）
	MsgThunderboltLobsterActivate MessageType = "thunderbolt_lobster_activate" // 免費射擊模式開始（Server→Client，個人+全服）
	MsgThunderboltLobsterShot     MessageType = "thunderbolt_lobster_shot"     // 自動射擊一次（Server→Client，全服）
	MsgThunderboltLobsterEnd      MessageType = "thunderbolt_lobster_end"      // 免費射擊模式結束（Server→Client，個人+全服）

	// 彩虹鳳凰 Power Up 系統（DAY-151）
	MsgRainbowPhoenixActivate MessageType = "rainbow_phoenix_activate" // Power Up 開始（Server→Client，個人+全服）
	MsgRainbowPhoenixEnd      MessageType = "rainbow_phoenix_end"      // Power Up 結束（Server→Client，個人+全服）
	MsgRainbowPhoenixStatus   MessageType = "rainbow_phoenix_status"   // Power Up 狀態（Server→Client，個人）

	// 吸血鬼成長倍率系統（DAY-152）
	MsgVampireGrow      MessageType = "vampire_grow"       // 吸血鬼倍率成長（Server→Client，全服）
	MsgVampireBloodMoon MessageType = "vampire_blood_moon" // 血月模式觸發（Server→Client，全服）
	MsgVampireKilled    MessageType = "vampire_killed"     // 吸血鬼被擊破（Server→Client，全服）

	// 水晶龍收集大獎系統（DAY-153）
	MsgCrystalDragonDrop   MessageType = "crystal_dragon_drop"   // 水晶掉落（Server→Client，全服）
	MsgCrystalDragonUpdate MessageType = "crystal_dragon_update" // 水晶進度更新（Server→Client，全服）
	MsgCrystalDragonReward MessageType = "crystal_dragon_reward" // 地獄龍大獎（Server→Client，全服）
	MsgCrystalDragonStatus MessageType = "crystal_dragon_status" // 水晶狀態（Server→Client，個人，登入時）

	// 獅子舞大獎爆發系統（DAY-168）
	MsgLionDanceBurst MessageType = "lion_dance_burst" // 獅子舞爆發廣播（Server→Client，全服）

	// 漩渦魚群吸引系統（DAY-169）
	MsgVortexFish MessageType = "vortex_fish" // 漩渦魚群廣播（Server→Client，全服）

	// 冰凍炸彈魚系統（DAY-170）
	MsgFreezeBomb MessageType = "freeze_bomb" // 冰凍炸彈魚廣播（Server→Client，全服）

	// 冰釣幸運輪盤系統（DAY-171）
	MsgIceFishingWheel MessageType = "ice_fishing_wheel" // 冰釣幸運輪盤廣播（Server→Client）

	// 幸運彩蛋魚系統（DAY-172）
	MsgLuckyEggFish MessageType = "lucky_egg_fish" // 幸運彩蛋魚廣播（Server→Client）

	// 彩虹幸運魚系統（DAY-173）
	MsgRainbowLuckyFish MessageType = "rainbow_lucky_fish" // 彩虹幸運魚廣播（Server→Client，全服）

	// 海葵觸手攻擊系統（DAY-174）
	MsgSeaAnemone MessageType = "sea_anemone" // 海葵觸手攻擊廣播（Server→Client，全服）

	// 幸運骰子魚系統（DAY-175）
	MsgLuckyDiceFish MessageType = "lucky_dice_fish" // 幸運骰子魚廣播（Server→Client）

	// 火焰風暴魚系統（DAY-176）
	MsgFireStormFish MessageType = "fire_storm_fish" // 火焰風暴魚廣播（Server→Client，全服）

	// 黃金寶藏魚系統（DAY-177）
	MsgGoldenTreasureFish    MessageType = "golden_treasure_fish"      // 黃金寶藏魚廣播（Server→Client）
	MsgGoldenTreasureOpen    MessageType = "golden_treasure_open"      // 玩家開箱請求（Client→Server）

	// 美人魚治癒系統（DAY-178）
	MsgMermaidHealing MessageType = "mermaid_healing" // 美人魚治癒廣播（Server→Client）

	// 幸運草魚系統（DAY-179）
	MsgLuckyCloverFish MessageType = "lucky_clover_fish" // 幸運草魚廣播（Server→Client，全服）

	// 彩虹鯊魚爆發系統（DAY-180）
	MsgRainbowSharkBurst MessageType = "rainbow_shark_burst" // 彩虹鯊魚爆發廣播（Server→Client，全服）

	// 雷霆鯊魚連鎖閃電系統（DAY-181）
	MsgThunderSharkChain MessageType = "thunder_shark_chain" // 雷霆鯊魚連鎖閃電廣播（Server→Client，全服）

	// 吸血鬼魚累積倍率系統（DAY-182）
	MsgVampireFish MessageType = "vampire_fish" // 吸血鬼魚倍率廣播（Server→Client）

	// 閃電魚自動連鎖系統（DAY-183）
	MsgLightningAutoChain MessageType = "lightning_auto_chain" // 閃電魚自動連鎖廣播（Server→Client，全服）

	// 隕石魚隕石雨系統（DAY-184）
	MsgMeteorFish MessageType = "meteor_fish" // 隕石魚隕石雨廣播（Server→Client，全服）

	// 鳳凰魚涅槃重生系統（DAY-185）
	MsgPhoenixFish MessageType = "phoenix_fish" // 鳳凰魚涅槃重生廣播（Server→Client，全服）

	// 龍龜不死 Boss 系統（DAY-186）
	MsgDragonTurtle MessageType = "dragon_turtle" // 龍龜不死 Boss 廣播（Server→Client，全服）

	// 連鎖爆炸魚系統（DAY-187）
	MsgChainBomb MessageType = "chain_bomb" // 連鎖爆炸魚廣播（Server→Client，全服）

	MsgCrocodileHunter MessageType = "crocodile_hunter" // 巨型鱷魚獵食廣播（Server→Client，全服，DAY-188）

	MsgTimeBombFish MessageType = "time_bomb_fish" // 時間炸彈魚廣播（Server→Client，全服，DAY-189）

	// 三重幸運魚系統（DAY-190）
	MsgTripleLuckyFish MessageType = "triple_lucky_fish" // 三重幸運魚廣播（Server→Client）

	// 魚群驚嚇連帶系統（DAY-191）
	MsgSchoolPanic MessageType = "school_panic" // 魚群驚嚇廣播（Server→Client，全服）

	// 搖滾骷髏演唱會系統（DAY-192）
	MsgRockSkeletonConcert MessageType = "rock_skeleton_concert" // 搖滾骷髏演唱會廣播（Server→Client，全服）

	// 電流水母電流網路系統（DAY-193）
	MsgElectricJellyfish MessageType = "electric_jellyfish" // 電流水母電流網路廣播（Server→Client，全服）

	// 長龍王雙環輪盤系統（DAY-194）
	MsgChainLongKing     MessageType = "chainlong_king"      // 長龍王雙環輪盤廣播（Server→Client，個人+全服）
	MsgChainLongKingStop MessageType = "chainlong_king_stop" // 玩家停止輪盤（Client→Server）

	// 鑽頭龍蝦穿透爆炸系統（DAY-195）
	MsgDrillLobster MessageType = "drill_lobster" // 鑽頭龍蝦穿透爆炸廣播（Server→Client，全服）

	// 巨型鮟鱇魚電擊寶箱系統（DAY-196）
	MsgAnglerfishElectric MessageType = "anglerfish_electric" // 巨型鮟鱇魚電擊廣播（Server→Client，全服）

	// 神秘龍魚八波攻擊系統（DAY-197）
	MsgMysticDragon MessageType = "mystic_dragon" // 神秘龍魚八波龍息攻擊廣播（Server→Client，全服）

	// 幽靈魚分身系統（DAY-198）
	MsgGhostFish MessageType = "ghost_fish" // 幽靈魚分身廣播（Server→Client，全服）

	// 雷霆龍蝦免費射擊系統（DAY-199）
	MsgThunderboltLobster MessageType = "thunderbolt_lobster" // 雷霆龍蝦免費射擊廣播（Server→Client，全服）

	// 冰鳳凰覺醒 BOSS 系統（DAY-200）
	MsgIcePhoenix MessageType = "ice_phoenix" // 冰鳳凰覺醒廣播（Server→Client，全服）

	// 連環炸彈蟹系統（DAY-201）
	MsgSerialBombCrab MessageType = "serial_bomb_crab" // 連環炸彈蟹廣播（Server→Client，全服）

	// 深淵漩渦魚系統（DAY-202）
	MsgAbyssVortex MessageType = "abyss_vortex" // 深淵漩渦廣播（Server→Client，全服）

	// 座頭鯨覺醒系統（DAY-203）
	MsgHumpbackWhale MessageType = "humpback_whale" // 座頭鯨覺醒廣播（Server→Client，全服）

	// 自由旋轉魚免費射擊系統（DAY-204）
	MsgFreeSpinFish MessageType = "free_spin_fish" // 自由旋轉魚廣播（Server→Client，個人+全服）

	// 獎池龍 Jackpot 抽獎系統（DAY-205）
	MsgJackpotDragon MessageType = "jackpot_dragon" // 獎池龍廣播（Server→Client，全服）

	MsgError MessageType = "error"
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
	PlayerID       string             `json:"player_id"`
	BombCharges    int                `json:"bomb_charges"`
	LaserCharges   int                `json:"laser_charges"`
	FreezeCharges  int                `json:"freeze_charges"`
	TornadoCharges int                `json:"tornado_charges"`  // DAY-134
	HomingCharges  int                `json:"homing_charges"`   // DAY-141
	DragonWrathCharges int            `json:"dragon_wrath_charges"` // DAY-154
	TorpedoCharges     int            `json:"torpedo_charges"`      // DAY-155
	RailgunCharges     int            `json:"railgun_charges"`      // DAY-157
	BlackHoleCharges   int            `json:"black_hole_charges"`   // DAY-166
	NewBalance     int                `json:"new_balance"`      // 購買後的新餘額（0=使用操作）
	Definitions    []SpecialWeaponDef `json:"definitions"`      // 武器定義（首次發送時填入）
	// 充能進度（DAY-134）
	BombChargeProgress    int `json:"bomb_charge_progress"`
	LaserChargeProgress   int `json:"laser_charge_progress"`
	FreezeChargeProgress  int `json:"freeze_charge_progress"`
	TornadoChargeProgress int `json:"tornado_charge_progress"`
	HomingChargeProgress  int `json:"homing_charge_progress"`       // DAY-141
	DragonWrathChargeProgress int `json:"dragon_wrath_charge_progress"` // DAY-154
	TorpedoChargeProgress     int `json:"torpedo_charge_progress"`      // DAY-155
	RailgunChargeProgress     int `json:"railgun_charge_progress"`      // DAY-157
	BlackHoleChargeProgress   int `json:"black_hole_charge_progress"`   // DAY-166
}

// SpecialWeaponChargedPayload 自動充能完成通知（Server → Client，DAY-134）
// 擊破目標累積充能滿時發送
type SpecialWeaponChargedPayload struct {
	PlayerID   string `json:"player_id"`
	WeaponType string `json:"weapon_type"` // "bomb" / "laser" / "freeze" / "tornado"
	WeaponName string `json:"weapon_name"`
	WeaponIcon string `json:"weapon_icon"`
	NewCharges int    `json:"new_charges"` // 充能後的發數
	Message    string `json:"message"`
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

// ---- 覺醒 BOSS 系統 Payloads（DAY-130）----

// AwakenBossSpawnPayload 覺醒 BOSS 出現廣播（Server → Client，全服）（DAY-130）
type AwakenBossSpawnPayload struct {
	InstanceID       string  `json:"instance_id"`         // BOSS 實例 ID
	BossType         string  `json:"boss_type"`           // BOSS 類型
	BossName         string  `json:"boss_name"`           // BOSS 名稱
	BossIcon         string  `json:"boss_icon"`           // BOSS 圖示
	BossColor        string  `json:"boss_color"`          // 顯示顏色
	MinMult          float64 `json:"min_mult"`            // 基礎最小倍率
	MaxMult          float64 `json:"max_mult"`            // 基礎最大倍率
	PowerUpMinMult   float64 `json:"powerup_min_mult"`    // Power Up 最小加成
	PowerUpMaxMult   float64 `json:"powerup_max_mult"`    // Power Up 最大加成
	PowerUpThreshold int     `json:"powerup_threshold"`   // 觸發 Power Up 所需命中次數
	DurationSeconds  float64 `json:"duration_seconds"`    // 在場時間
	Message          string  `json:"message"`             // 公告訊息
}

// AwakenBossHitPayload 命中覺醒 BOSS 廣播（Server → Client，全服）（DAY-130）
type AwakenBossHitPayload struct {
	InstanceID      string  `json:"instance_id"`       // BOSS 實例 ID
	PlayerID        string  `json:"player_id"`         // 命中玩家 ID
	PlayerName      string  `json:"player_name"`       // 命中玩家名稱
	Multiplier      float64 `json:"multiplier"`        // 本次倍率
	Reward          int     `json:"reward"`            // 本次獎勵
	NewBalance      int     `json:"new_balance"`       // 玩家新餘額（只發給命中者）
	HitCount        int     `json:"hit_count"`         // 累計命中次數
	PowerUpProgress float64 `json:"powerup_progress"`  // Power Up 進度（0.0-1.0）
	TotalReward     int     `json:"total_reward"`      // 累計總獎勵
}

// AwakenBossPowerUpPayload Power Up 觸發廣播（Server → Client，全服）（DAY-130）
type AwakenBossPowerUpPayload struct {
	InstanceID   string  `json:"instance_id"`    // BOSS 實例 ID
	BossName     string  `json:"boss_name"`      // BOSS 名稱
	BossIcon     string  `json:"boss_icon"`      // BOSS 圖示
	PlayerID     string  `json:"player_id"`      // 觸發玩家 ID
	PlayerName   string  `json:"player_name"`    // 觸發玩家名稱
	Multiplier   float64 `json:"multiplier"`     // Power Up 倍率
	Reward       int     `json:"reward"`         // Power Up 獎勵
	NewBalance   int     `json:"new_balance"`    // 玩家新餘額
	PowerUpCount int     `json:"powerup_count"`  // 第幾次 Power Up
	Message      string  `json:"message"`        // 廣播訊息
}

// AwakenBossLeavePayload 覺醒 BOSS 離開廣播（Server → Client，全服）（DAY-130）
type AwakenBossLeavePayload struct {
	InstanceID   string `json:"instance_id"`   // BOSS 實例 ID
	BossName     string `json:"boss_name"`     // BOSS 名稱
	BossIcon     string `json:"boss_icon"`     // BOSS 圖示
	HitCount     int    `json:"hit_count"`     // 總命中次數
	PowerUpCount int    `json:"powerup_count"` // 總 Power Up 次數
	TotalReward  int    `json:"total_reward"`  // 總獎勵
	Message      string `json:"message"`      // 離開訊息
}

// AwakenBossStatusPayload 覺醒 BOSS 狀態（Server → Client，個人）（DAY-130）
type AwakenBossStatusPayload struct {
	Active           bool    `json:"active"`             // 是否有活躍 BOSS
	InstanceID       string  `json:"instance_id"`        // BOSS 實例 ID
	BossType         string  `json:"boss_type"`          // BOSS 類型
	BossName         string  `json:"boss_name"`          // BOSS 名稱
	BossIcon         string  `json:"boss_icon"`          // BOSS 圖示
	BossColor        string  `json:"boss_color"`         // 顯示顏色
	MinMult          float64 `json:"min_mult"`           // 基礎最小倍率
	MaxMult          float64 `json:"max_mult"`           // 基礎最大倍率
	PowerUpMinMult   float64 `json:"powerup_min_mult"`   // Power Up 最小加成
	PowerUpMaxMult   float64 `json:"powerup_max_mult"`   // Power Up 最大加成
	PowerUpThreshold int     `json:"powerup_threshold"`  // 觸發 Power Up 所需命中次數
	HitCount         int     `json:"hit_count"`          // 累計命中次數
	PowerUpCount     int     `json:"powerup_count"`      // Power Up 次數
	PowerUpProgress  float64 `json:"powerup_progress"`   // Power Up 進度
	TotalReward      int     `json:"total_reward"`       // 累計總獎勵
	RemainingSeconds float64 `json:"remaining_seconds"`  // 剩餘秒數
}

// ---- 連勝獎勵系統 Payloads（DAY-131）----

// WinStreakUpdatePayload 連勝更新（Server → Client，個人）（DAY-131）
type WinStreakUpdatePayload struct {
	Current           int     `json:"current"`             // 當前連勝次數
	MaxStreak         int     `json:"max_streak"`          // 本 session 最高連勝
	NextMilestone     int     `json:"next_milestone"`      // 下一個里程碑次數
	NextMilestoneName string  `json:"next_milestone_name"` // 下一個里程碑名稱
	ProgressToNext    float64 `json:"progress_to_next"`    // 到下一個里程碑的進度
	SecondsToExpiry   float64 `json:"seconds_to_expiry"`   // 超時倒數（秒）
}

// WinStreakMilestonePayload 里程碑達成（Server → Client）（DAY-131）
type WinStreakMilestonePayload struct {
	PlayerID    string  `json:"player_id"`    // 玩家 ID
	PlayerName  string  `json:"player_name"`  // 玩家名稱
	Streak      int     `json:"streak"`       // 達成的連勝次數
	Level       int     `json:"level"`        // 里程碑等級
	LevelName   string  `json:"level_name"`   // 里程碑名稱
	Icon        string  `json:"icon"`         // 圖示
	Color       string  `json:"color"`        // 顏色
	BonusReward int     `json:"bonus_reward"` // 額外獎勵
	NewBalance  int     `json:"new_balance"`  // 新餘額
	Broadcast   bool    `json:"broadcast"`    // 是否全服廣播
}

// WinStreakResetPayload 連勝重置（Server → Client，個人）（DAY-131）
type WinStreakResetPayload struct {
	FinalStreak int `json:"final_streak"` // 最終連勝次數
	MaxStreak   int `json:"max_streak"`   // 本 session 最高連勝
}

// ---- 閃電鰻連鎖攻擊系統 Payloads（DAY-132）----

// LightningEelJumpEntry 單次跳躍結果
type LightningEelJumpEntry struct {
	TargetInstanceID string  `json:"target_instance_id"` // 被跳躍的目標 instance ID
	TargetDefID      string  `json:"target_def_id"`      // 目標定義 ID
	TargetName       string  `json:"target_name"`        // 目標名稱
	Killed           bool    `json:"killed"`             // 是否擊破
	Multiplier       float64 `json:"multiplier"`         // 目標原始倍率
	Reward           int64   `json:"reward"`             // 實際獎勵（Killed 才有）
	JumpIndex        int     `json:"jump_index"`         // 第幾次跳躍（1-based）
}

// LightningEelChainPayload 連鎖攻擊結果廣播（Server → Client，全服）（DAY-132）
type LightningEelChainPayload struct {
	PlayerID        string                  `json:"player_id"`         // 觸發玩家 ID
	PlayerName      string                  `json:"player_name"`       // 觸發玩家名稱
	TriggerTargetID string                  `json:"trigger_target_id"` // 觸發連鎖的閃電鰻 instance ID
	Jumps           []LightningEelJumpEntry `json:"jumps"`             // 所有跳躍結果
	TotalKills      int                     `json:"total_kills"`       // 總擊破數
	TotalReward     int64                   `json:"total_reward"`      // 總獎勵
	NewBalance      int64                   `json:"new_balance"`       // 觸發玩家新餘額
}

// LightningEelStatusPayload 閃電鰻冷卻狀態（Server → Client，個人）（DAY-132）
type LightningEelStatusPayload struct {
	PlayerID     string `json:"player_id"`      // 玩家 ID
	CooldownLeft int    `json:"cooldown_left"`  // 冷卻剩餘秒數（0 = 可觸發）
	MaxJumps     int    `json:"max_jumps"`      // 最大跳躍次數
	JumpRange    float64 `json:"jump_range"`   // 跳躍範圍（Client 用於視覺）
}

// ---- 狂熱模式系統 Payloads（DAY-133）----

// FeverModeStartPayload 狂熱模式開始廣播（Server → Client，全服）（DAY-133）
type FeverModeStartPayload struct {
	PlayerID    string  `json:"player_id"`    // 觸發玩家 ID
	PlayerName  string  `json:"player_name"`  // 觸發玩家名稱
	SecondsLeft int     `json:"seconds_left"` // 剩餘秒數
	MultBoost   float64 `json:"mult_boost"`   // 倍率加成（1.5）
	IsSelf      bool    `json:"is_self"`      // 是否為自己觸發（Client 端填充）
}

// FeverModeEndPayload 狂熱模式結束（Server → Client，個人）（DAY-133）
type FeverModeEndPayload struct {
	PlayerID     string `json:"player_id"`     // 玩家 ID
	TotalFevered int    `json:"total_fevered"` // 本 session 觸發次數
	CooldownLeft int    `json:"cooldown_left"` // 冷卻剩餘秒數
}

// FeverModeStatusPayload 狂熱模式狀態（Server → Client，個人）（DAY-133）
type FeverModeStatusPayload struct {
	PlayerID     string  `json:"player_id"`     // 玩家 ID
	IsActive     bool    `json:"is_active"`     // 是否正在狂熱中
	SecondsLeft  int     `json:"seconds_left"`  // 剩餘秒數（0 = 未觸發）
	CooldownLeft int     `json:"cooldown_left"` // 冷卻剩餘秒數（0 = 可觸發）
	MultBoost    float64 `json:"mult_boost"`    // 當前倍率加成
	KillProgress int     `json:"kill_progress"` // 觸發進度（0-5）
	TriggerKills int     `json:"trigger_kills"` // 觸發所需擊破數（5）
	TotalFevered int     `json:"total_fevered"` // 本 session 觸發次數
}

// ---- 失敗補償系統（DAY-135）----

// UnluckyBonusPayload 失敗補償觸發通知（Server → Client，個人）
type UnluckyBonusPayload struct {
	PlayerID    string `json:"player_id"`
	BonusAmount int    `json:"bonus_amount"` // 補償金額
	NewBalance  int    `json:"new_balance"`  // 補償後餘額
	Message     string `json:"message"`      // 顯示訊息
}

// UnluckyBonusStatusPayload 失敗補償狀態（Server → Client，個人）
type UnluckyBonusStatusPayload struct {
	PlayerID     string `json:"player_id"`
	ShotCount    int    `json:"shot_count"`    // 已追蹤射擊次數
	TrackingMax  int    `json:"tracking_max"`  // 最大追蹤次數
	NetLoss      int    `json:"net_loss"`       // 淨虧損
	RatioPercent int    `json:"ratio_percent"` // 花費/回報百分比
	CooldownLeft int    `json:"cooldown_left"` // 冷卻剩餘秒數
	BonusCount   int    `json:"bonus_count"`   // 累計補償次數
}

// ---- 競速獵殺系統 Payloads（DAY-136）----

// SpeedRaceStartPayload 競速獵殺開始廣播（Server → Client，全服）
type SpeedRaceStartPayload struct {
	TargetInstanceID string  `json:"target_instance_id"` // 競速目標 instanceID
	TargetDefID      string  `json:"target_def_id"`      // 目標定義 ID
	TargetName       string  `json:"target_name"`        // 目標名稱
	TargetMult       float64 `json:"target_mult"`        // 目標倍率
	SecondsLeft      float64 `json:"seconds_left"`       // 競速剩餘秒數
	BonusMult        float64 `json:"bonus_mult"`         // 第一名獎勵倍率（3.0）
	Message          string  `json:"message"`            // 顯示訊息
}

// SpeedRaceEndPayload 競速獵殺結束廣播（Server → Client，全服）
type SpeedRaceEndPayload struct {
	WinnerID    string  `json:"winner_id"`    // 第一名玩家 ID
	WinnerName  string  `json:"winner_name"`  // 第一名玩家名稱
	TargetName  string  `json:"target_name"`  // 競速目標名稱
	TargetMult  float64 `json:"target_mult"`  // 競速目標倍率
	BonusMult   float64 `json:"bonus_mult"`   // 第一名獎勵倍率
	Message     string  `json:"message"`      // 顯示訊息
}

// SpeedRaceCancelPayload 競速獵殺取消廣播（Server → Client，全服）
type SpeedRaceCancelPayload struct {
	TargetInstanceID string `json:"target_instance_id"` // 競速目標 instanceID
	TargetName       string `json:"target_name"`        // 目標名稱
	Message          string `json:"message"`            // 顯示訊息
}

// SpeedRaceResultPayload 競速個人結果（Server → Client，個人）
type SpeedRaceResultPayload struct {
	PlayerID    string  `json:"player_id"`    // 玩家 ID
	DisplayName string  `json:"display_name"` // 玩家名稱
	Rank        int     `json:"rank"`         // 名次（1/2/3）
	BonusMult   float64 `json:"bonus_mult"`   // 獎勵倍率
	RankIcon    string  `json:"rank_icon"`    // 名次圖示（🥇/🥈/🥉）
	Message     string  `json:"message"`      // 顯示訊息
}

// ---- 全服目標懸賞系統 Payloads（DAY-137）----

// PostBountyPayload 玩家下懸賞請求（Client → Server）
type PostBountyPayload struct {
	TargetInstanceID string `json:"target_instance_id"` // 目標 instanceID
	Amount           int    `json:"amount"`             // 懸賞金額
}

// BountySnap 懸賞快照（用於列表）
type BountySnap struct {
	BountyID         string  `json:"bounty_id"`
	TargetInstanceID string  `json:"target_instance_id"`
	TargetDefID      string  `json:"target_def_id"`
	TargetName       string  `json:"target_name"`
	TargetMult       float64 `json:"target_mult"`
	PosterID         string  `json:"poster_id"`
	PosterName       string  `json:"poster_name"`
	Amount           int     `json:"amount"`
	SecondsLeft      float64 `json:"seconds_left"`
}

// BountyPostedPayload 懸賞發布廣播（Server → Client，全服）
type BountyPostedPayload struct {
	BountyID         string  `json:"bounty_id"`
	TargetInstanceID string  `json:"target_instance_id"`
	TargetDefID      string  `json:"target_def_id"`
	TargetName       string  `json:"target_name"`
	TargetMult       float64 `json:"target_mult"`
	PosterID         string  `json:"poster_id"`
	PosterName       string  `json:"poster_name"`
	Amount           int     `json:"amount"`
	SecondsLeft      float64 `json:"seconds_left"`
	Message          string  `json:"message"`
}

// BountyClaimedPayload 懸賞領取通知（Server → Client，個人）
type BountyClaimedPayload struct {
	KillerID    string `json:"killer_id"`
	KillerName  string `json:"killer_name"`
	TotalAmount int    `json:"total_amount"`
	BountyCount int    `json:"bounty_count"`
	NewBalance  int    `json:"new_balance"`
	Message     string `json:"message"`
}

// BountyKilledPayload 懸賞目標擊破廣播（Server → Client，全服）
type BountyKilledPayload struct {
	KillerID    string `json:"killer_id"`
	KillerName  string `json:"killer_name"`
	TargetName  string `json:"target_name"`
	TotalAmount int    `json:"total_amount"`
	BountyCount int    `json:"bounty_count"`
	Message     string `json:"message"`
}

// BountyExpiredPayload 懸賞過期通知（Server → Client）
type BountyExpiredPayload struct {
	BountyID   string `json:"bounty_id"`
	TargetName string `json:"target_name"`
	Amount     int    `json:"amount"`
	Message    string `json:"message"`
}

// BountyListPayload 懸賞列表回應（Server → Client，個人）
type BountyListPayload struct {
	Bounties     []BountySnap `json:"bounties"`
	CooldownLeft int          `json:"cooldown_left"` // 玩家下懸賞冷卻剩餘秒數
}

// BountyErrorPayload 懸賞操作失敗（Server → Client，個人）
type BountyErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ---- 全服倍率風暴系統 Payloads（DAY-138）----

// MultStormStartPayload 風暴開始廣播（Server → Client，全服）
type MultStormStartPayload struct {
	TierName    string  `json:"tier_name"`    // 風暴等級名稱
	TierIcon    string  `json:"tier_icon"`    // 圖示
	TierColor   string  `json:"tier_color"`   // 顏色（hex）
	MultBoost   float64 `json:"mult_boost"`   // 倍率加成（2.0/3.0/5.0）
	SecondsLeft float64 `json:"seconds_left"` // 持續秒數
	Message     string  `json:"message"`      // 顯示訊息
}

// MultStormEndPayload 風暴結束廣播（Server → Client，全服）
type MultStormEndPayload struct {
	Message string `json:"message"` // 顯示訊息
}

// ---- 雙環輪盤系統 Payloads（DAY-139）----

// DualRouletteStartPayload 輪盤開始（Server → Client，個人）
type DualRouletteStartPayload struct {
	PlayerID     string    `json:"player_id"`     // 觸發玩家 ID
	TargetMult   float64   `json:"target_mult"`   // 觸發目標倍率
	BaseReward   int       `json:"base_reward"`   // 基礎獎勵（輪盤加成的基礎）
	SpinDuration float64   `json:"spin_duration"` // 旋轉持續秒數
	InnerRing    []float64 `json:"inner_ring"`    // 內環倍率選項
	OuterRing    []float64 `json:"outer_ring"`    // 外環倍率選項
}

// DualRouletteResultPayload 輪盤結果（Server → Client，個人）
type DualRouletteResultPayload struct {
	PlayerID    string  `json:"player_id"`    // 玩家 ID
	InnerResult float64 `json:"inner_result"` // 內環停止結果
	OuterResult float64 `json:"outer_result"` // 外環停止結果
	Combined    float64 `json:"combined"`     // 最終倍率 = 內 × 外
	BonusReward int     `json:"bonus_reward"` // 額外獎勵金幣
	NewBalance  int     `json:"new_balance"`  // 新餘額
}

// DualRouletteStatusPayload 輪盤狀態（Server → Client，個人）
type DualRouletteStatusPayload struct {
	CooldownLeft int `json:"cooldown_left"` // 冷卻剩餘秒數（0 = 可觸發）
}

// ---- 全服 Mega Catch 事件系統 Payloads（DAY-140）----

// MegaCatchStartPayload 事件開始廣播（Server → Client，全服）
type MegaCatchStartPayload struct {
	TierName    string  `json:"tier_name"`    // 等級名稱
	TierIcon    string  `json:"tier_icon"`    // 圖示
	TierColor   string  `json:"tier_color"`   // 顏色
	RewardBoost float64 `json:"reward_boost"` // 獎勵倍率加成
	SpawnBoost  float64 `json:"spawn_boost"`  // 稀有目標生成加成
	Duration    float64 `json:"duration"`     // 持續秒數
	SecondsLeft float64 `json:"seconds_left"` // 剩餘秒數
}

// MegaCatchEndPayload 事件結束廣播（Server → Client，全服）
type MegaCatchEndPayload struct {
	Message string `json:"message"` // 顯示訊息
}

// MegaCatchStatusPayload 事件狀態（Server → Client，個人）
type MegaCatchStatusPayload struct {
	IsActive    bool    `json:"is_active"`    // 是否活躍
	TierName    string  `json:"tier_name"`    // 等級名稱
	TierIcon    string  `json:"tier_icon"`    // 圖示
	TierColor   string  `json:"tier_color"`   // 顏色
	RewardBoost float64 `json:"reward_boost"` // 獎勵倍率加成
	SpawnBoost  float64 `json:"spawn_boost"`  // 稀有目標生成加成
	SecondsLeft float64 `json:"seconds_left"` // 剩餘秒數
}

// HomingMissileResultPayload 追蹤飛彈命中結果（Server → Client，DAY-141）
// 追蹤飛彈自動鎖定倍率最高的目標，100% 命中，獎勵 ×1.5
type HomingMissileResultPayload struct {
	PlayerID   string  `json:"player_id"`   // 使用者 ID
	TargetID   string  `json:"target_id"`   // 命中的目標 InstanceID
	DefID      string  `json:"def_id"`      // 目標定義 ID
	Multiplier float64 `json:"multiplier"`  // 目標倍率
	BaseReward int     `json:"base_reward"` // 基礎獎勵
	FinalReward int    `json:"final_reward"` // 最終獎勵（×1.5）
	NewBalance int     `json:"new_balance"` // 命中後餘額
	Killed     bool    `json:"killed"`      // 是否擊破（100% 命中，但不一定擊破）
	Message    string  `json:"message"`     // 顯示訊息
}

// DragonWrathChargePayload 龍怒怒氣值更新（Server → Client，DAY-154）
// 每次射擊後發送，讓 Client 更新怒氣條
type DragonWrathChargePayload struct {
	PlayerID    string `json:"player_id"`    // 玩家 ID
	Progress    int    `json:"progress"`     // 當前怒氣進度（0-60）
	Required    int    `json:"required"`     // 充滿所需（60）
	Charges     int    `json:"charges"`      // 當前持有發數
	MaxCharges  int    `json:"max_charges"`  // 最大持有發數（1）
	JustCharged bool   `json:"just_charged"` // 是否剛充滿一發
}

// DragonWrathMeteorEntry 龍怒流星雨命中的目標條目（DAY-154）
type DragonWrathMeteorEntry struct {
	InstanceID  string  `json:"instance_id"`  // 目標 InstanceID
	DefID       string  `json:"def_id"`       // 目標定義 ID
	Multiplier  float64 `json:"multiplier"`   // 目標倍率
	Reward      int     `json:"reward"`       // 獎勵
	Killed      bool    `json:"killed"`       // 是否擊破
	IsImmortal  bool    `json:"is_immortal"`  // 是否是不死 BOSS
	IsBoss      bool    `json:"is_boss"`      // 是否是 BOSS
}

// DragonWrathResultPayload 龍怒流星雨結果廣播（Server → Client，DAY-154）
// Phase: "wrath_start" → "meteor_N" → "result"
type DragonWrathResultPayload struct {
	KillerID      string                   `json:"killer_id"`      // 觸發玩家 ID
	KillerName    string                   `json:"killer_name"`    // 觸發玩家名稱
	Phase         string                   `json:"phase"`          // 當前階段
	MeteorIndex   int                      `json:"meteor_index"`   // 流星序號（meteor_N 時）
	MeteorX       float64                  `json:"meteor_x"`       // 流星落點 X
	MeteorY       float64                  `json:"meteor_y"`       // 流星落點 Y
	HitTargets    []DragonWrathMeteorEntry `json:"hit_targets"`    // 命中目標（result 時）
	TotalReward   int                      `json:"total_reward"`   // 總獎勵（result 時）
	NewBalance    int                      `json:"new_balance"`    // 結果後餘額（result 時，僅觸發者）
	ImmortalHits  int                      `json:"immortal_hits"`  // 命中不死 BOSS 次數
	ImmortalReward int                     `json:"immortal_reward"` // 不死 BOSS 獎勵
}

// TorpedoKillEntry 魚雷命中的目標條目（DAY-155）
type TorpedoKillEntry struct {
	InstanceID string  `json:"instance_id"` // 目標 InstanceID
	DefID      string  `json:"def_id"`      // 目標定義 ID
	Multiplier float64 `json:"multiplier"`  // 目標倍率
	Reward     int     `json:"reward"`      // 獎勵金幣
	Killed     bool    `json:"killed"`      // 是否擊破
}

// TorpedoResultPayload 魚雷爆炸結果廣播（Server → Client，DAY-155）
// Phase: "torpedo_launch" → "explosion" → "result"
type TorpedoResultPayload struct {
	ShooterID   string             `json:"shooter_id"`   // 射擊玩家 ID
	ShooterName string             `json:"shooter_name"` // 射擊玩家名稱
	Phase       string             `json:"phase"`        // 當前階段
	TargetX     float64            `json:"target_x"`     // 爆炸中心 X
	TargetY     float64            `json:"target_y"`     // 爆炸中心 Y
	HitTargets  []TorpedoKillEntry `json:"hit_targets"`  // 命中目標（result 時）
	TotalReward int                `json:"total_reward"` // 總獎勵（result 時）
	NewBalance  int                `json:"new_balance"`  // 結果後餘額（result 時，僅射擊者）
	Cost        int                `json:"cost"`         // 魚雷費用（6x betLevel）
}

// RailgunKillEntry 軌道炮穿透命中的目標條目（DAY-157）
type RailgunKillEntry struct {
	InstanceID string  `json:"instance_id"` // 目標 InstanceID
	DefID      string  `json:"def_id"`      // 目標定義 ID
	X          float64 `json:"x"`           // 目標 X 座標（用於排序穿透順序）
	Multiplier float64 `json:"multiplier"`  // 目標倍率
	Reward     int     `json:"reward"`      // 獎勵金幣
	Killed     bool    `json:"killed"`      // 是否擊破
}

// RailgunResultPayload 軌道炮穿透結果廣播（Server → Client，DAY-157）
// Phase: "railgun_charge" → "railgun_fire" → "result"
type RailgunResultPayload struct {
	ShooterID   string             `json:"shooter_id"`   // 射擊玩家 ID
	ShooterName string             `json:"shooter_name"` // 射擊玩家名稱
	Phase       string             `json:"phase"`        // 當前階段
	TargetY     float64            `json:"target_y"`     // 光束 Y 座標
	HitTargets  []RailgunKillEntry `json:"hit_targets"`  // 命中目標（result 時）
	TotalReward int                `json:"total_reward"` // 總獎勵（result 時）
	NewBalance  int                `json:"new_balance"`  // 結果後餘額（result 時，僅射擊者）
	Cost        int                `json:"cost"`         // 軌道炮費用（15x betLevel）
}

// ---- 黑洞漩渦（DAY-166）----

// BlackHoleKillEntry 黑洞擊破的目標條目（DAY-166）
type BlackHoleKillEntry struct {
	InstanceID string  `json:"instance_id"` // 目標 InstanceID
	DefID      string  `json:"def_id"`      // 目標 DefID
	Multiplier float64 `json:"multiplier"`  // 目標倍率
	Reward     int     `json:"reward"`      // 獎勵金幣
}

// BlackHoleResultPayload 黑洞漩渦爆炸結果廣播（Server → Client，DAY-166）
// Phase: "black_hole_place" → "black_hole_suck" → "black_hole_explode" → "result"
type BlackHoleResultPayload struct {
	ShooterID   string               `json:"shooter_id"`   // 放置玩家 ID
	ShooterName string               `json:"shooter_name"` // 放置玩家名稱
	Phase       string               `json:"phase"`        // 當前階段
	CenterX     float64              `json:"center_x"`     // 黑洞中心 X
	CenterY     float64              `json:"center_y"`     // 黑洞中心 Y
	Radius      float64              `json:"radius"`       // 吸引半徑（300px）
	SuckedCount int                  `json:"sucked_count"` // 被吸入的目標數（suck 階段）
	HitTargets  []BlackHoleKillEntry `json:"hit_targets"`  // 命中目標（result 時）
	TotalReward int                  `json:"total_reward"` // 總獎勵（result 時）
	NewBalance  int                  `json:"new_balance"`  // 結果後餘額（result 時，僅放置者）
	Cost        int                  `json:"cost"`         // 黑洞費用（10x betLevel）
}

// ---- 黃金輪盤螃蟹（DAY-167）----

// RouletteCrabStartPayload 黃金輪盤螃蟹開始廣播（Server → Client，DAY-167）
// 擊破 T125 後觸發，廣播給所有玩家（旁觀者看到輪盤旋轉）
type RouletteCrabStartPayload struct {
	PlayerID    string  `json:"player_id"`    // 觸發玩家 ID
	PlayerName  string  `json:"player_name"`  // 觸發玩家名稱
	TargetMult  float64 `json:"target_mult"`  // 螃蟹本身的倍率
	BaseReward  int     `json:"base_reward"`  // 螃蟹擊破的基礎獎勵
	SpinSecs    float64 `json:"spin_secs"`    // 旋轉持續秒數（4 秒）
	WheelSlots  []float64 `json:"wheel_slots"` // 輪盤格子（8格）
}

// RouletteCrabResultPayload 黃金輪盤螃蟹結果廣播（Server → Client，DAY-167）
// 玩家停止或超時後廣播
type RouletteCrabResultPayload struct {
	PlayerID    string  `json:"player_id"`    // 觸發玩家 ID
	PlayerName  string  `json:"player_name"`  // 觸發玩家名稱
	WheelResult float64 `json:"wheel_result"` // 輪盤結果倍率（10x-200x）
	SlotIndex   int     `json:"slot_index"`   // 輪盤格子索引（0-7，用於 Client 動畫定位）
	BaseReward  int     `json:"base_reward"`  // 螃蟹擊破的基礎獎勵
	BonusReward int     `json:"bonus_reward"` // 輪盤額外獎勵（基礎獎勵 × 輪盤倍率）
	NewBalance  int     `json:"new_balance"`  // 結果後餘額（僅觸發玩家有值）
	IsAutoStop  bool    `json:"is_auto_stop"` // 是否超時自動停止
}

// RouletteCrabStatusPayload 黃金輪盤螃蟹冷卻狀態（Server → Client，DAY-167）
// 玩家登入時發送
type RouletteCrabStatusPayload struct {
	PlayerID     string `json:"player_id"`
	CooldownLeft int    `json:"cooldown_left"` // 冷卻剩餘秒數（0=可觸發）
}

// GoldenTurtleTimeStopPayload 黃金海龜時間停止廣播（Server → Client，DAY-159）
// Phase: "time_stop_start" → "time_stop_end"
type GoldenTurtleTimeStopPayload struct {
	TriggerID    string  `json:"trigger_id"`    // 觸發的 T119 InstanceID
	TriggerX     float64 `json:"trigger_x"`     // 觸發位置 X
	TriggerY     float64 `json:"trigger_y"`     // 觸發位置 Y
	KillerID     string  `json:"killer_id"`     // 擊破玩家 ID
	KillerName   string  `json:"killer_name"`   // 擊破玩家名稱
	Phase        string  `json:"phase"`         // 當前階段
	DurationSecs float64 `json:"duration_secs"` // 停止時間（秒）
}

// LuckyStarFishPayload 幸運星魚全場倍率翻倍廣播（Server → Client，DAY-160）
// Phase: "lucky_start" → "lucky_end"
type LuckyStarFishPayload struct {
	TriggerID    string  `json:"trigger_id"`    // 觸發的 T120 InstanceID
	TriggerX     float64 `json:"trigger_x"`     // 觸發位置 X
	TriggerY     float64 `json:"trigger_y"`     // 觸發位置 Y
	KillerID     string  `json:"killer_id"`     // 擊破玩家 ID
	KillerName   string  `json:"killer_name"`   // 擊破玩家名稱
	Phase        string  `json:"phase"`         // 當前階段
	DurationSecs float64 `json:"duration_secs"` // 倍率翻倍時間（秒）
	MultBonus    float64 `json:"mult_bonus"`    // 倍率加成（2.0 = 翻倍）
}

// GoldenSharkBerserkPayload 黃金鯊魚全服狂暴模式廣播（Server → Client，DAY-161）
// Phase: "berserk_start" → "berserk_end"
type GoldenSharkBerserkPayload struct {
	TriggerID    string  `json:"trigger_id"`    // 觸發的 T121 InstanceID
	TriggerX     float64 `json:"trigger_x"`     // 觸發位置 X
	TriggerY     float64 `json:"trigger_y"`     // 觸發位置 Y
	KillerID     string  `json:"killer_id"`     // 擊破玩家 ID
	KillerName   string  `json:"killer_name"`   // 擊破玩家名稱
	Phase        string  `json:"phase"`         // 當前階段：berserk_start / berserk_end
	DurationSecs float64 `json:"duration_secs"` // 狂暴模式持續時間（秒）
	MultBonus    float64 `json:"mult_bonus"`    // 倍率加成（1.5 = 全場 ×1.5）
}

// MoneyFishRewardPayload 金幣魚王即時獎勵廣播（Server → Client，DAY-162）
type MoneyFishRewardPayload struct {
	TriggerID     string  `json:"trigger_id"`     // 觸發的 T122 InstanceID
	TriggerX      float64 `json:"trigger_x"`      // 觸發位置 X
	TriggerY      float64 `json:"trigger_y"`      // 觸發位置 Y
	KillerID      string  `json:"killer_id"`      // 擊破玩家 ID
	KillerName    string  `json:"killer_name"`    // 擊破玩家名稱
	InstantReward int     `json:"instant_reward"` // 即時獎勵金幣數
	MultUsed      int     `json:"mult_used"`      // 使用的倍率（20-50）
	BetLevel      int     `json:"bet_level"`      // 玩家當前 betLevel
}

// CaptainRaceEntry 船長魚競速排名條目（DAY-163）
type CaptainRaceEntry struct {
	Rank        int    `json:"rank"`         // 排名（1-based）
	PlayerID    string `json:"player_id"`    // 玩家 ID
	PlayerName  string `json:"player_name"`  // 玩家名稱
	KillCount   int    `json:"kill_count"`   // 擊破數
	TotalReward int    `json:"total_reward"` // 累積獎勵
}

// CaptainFishRacePayload 船長魚全服競速模式廣播（Server → Client，DAY-163）
// Phase: "race_start" → "race_update" → "race_end" → "race_reward"（個人）
type CaptainFishRacePayload struct {
	TriggerID     string             `json:"trigger_id"`     // 觸發的 T123 InstanceID
	TriggerX      float64            `json:"trigger_x"`      // 觸發位置 X
	TriggerY      float64            `json:"trigger_y"`      // 觸發位置 Y
	KillerID      string             `json:"killer_id"`      // 擊破玩家 ID
	KillerName    string             `json:"killer_name"`    // 擊破玩家名稱
	Phase         string             `json:"phase"`          // 當前階段
	DurationSecs  float64            `json:"duration_secs"`  // 競速持續時間（秒）
	RemainingTime float64            `json:"remaining_time"` // 剩餘時間（秒）
	Entries       []CaptainRaceEntry `json:"entries"`        // 排名列表
	MyRank        int                `json:"my_rank"`        // 我的排名（race_reward 時）
	MyBonus       int                `json:"my_bonus"`       // 我的獎勵（race_reward 時）
	MyKillCount   int                `json:"my_kill_count"`  // 我的擊破數（race_reward 時）
}

// RoyalChainLightningEntry 皇家閃電鰻連鎖電擊的目標條目（DAY-156）
type RoyalChainLightningEntry struct {
	InstanceID string  `json:"instance_id"` // 目標 InstanceID
	DefID      string  `json:"def_id"`      // 目標定義 ID
	Multiplier float64 `json:"multiplier"`  // 目標倍率
	Reward     int     `json:"reward"`      // 獎勵金幣
	Killed     bool    `json:"killed"`      // 是否擊破
	JumpIndex  int     `json:"jump_index"`  // 第幾跳（1-15）
	FromX      float64 `json:"from_x"`      // 電擊起點 X（上一個目標位置）
	FromY      float64 `json:"from_y"`      // 電擊起點 Y
	ToX        float64 `json:"to_x"`        // 電擊終點 X（本目標位置）
	ToY        float64 `json:"to_y"`        // 電擊終點 Y
}

// RoyalChainLightningPayload 皇家閃電鰻持續連鎖電擊廣播（Server → Client，DAY-156）
// Phase: "chain_start" → "jump_N"（每跳一次）→ "result"
type RoyalChainLightningPayload struct {
	TriggerID     string                     `json:"trigger_id"`     // 觸發的 T118 InstanceID
	TriggerX      float64                    `json:"trigger_x"`      // 觸發位置 X
	TriggerY      float64                    `json:"trigger_y"`      // 觸發位置 Y
	KillerID      string                     `json:"killer_id"`      // 擊破玩家 ID
	KillerName    string                     `json:"killer_name"`    // 擊破玩家名稱
	Phase         string                     `json:"phase"`          // 當前階段
	JumpIndex     int                        `json:"jump_index"`     // 當前跳數（jump_N 時）
	JumpEntry     *RoyalChainLightningEntry  `json:"jump_entry"`     // 本跳命中目標（jump_N 時）
	AllEntries    []RoyalChainLightningEntry `json:"all_entries"`    // 所有命中目標（result 時）
	TotalReward   int                        `json:"total_reward"`   // 總獎勵（result 時）
	TotalJumps    int                        `json:"total_jumps"`    // 總跳數（result 時）
}

// DrillKillEntry 鑽頭連帶擊破的目標條目（DAY-142）
type DrillKillEntry struct {
	InstanceID string  `json:"instance_id"`
	DefID      string  `json:"def_id"`
	Multiplier float64 `json:"multiplier"`
	Reward     int     `json:"reward"`
	Phase      string  `json:"phase"` // "penetrate" or "explosion"
}

// DrillLobsterChainPayload 鑽頭龍蝦連帶效果廣播（Server → Client，DAY-142）
// Phase: "drill_start" → "explosion" → "result"
type DrillLobsterChainPayload struct {
	TriggerID     string          `json:"trigger_id"`    // 觸發的 T106 InstanceID
	TriggerX      float64         `json:"trigger_x"`     // 觸發位置 X
	TriggerY      float64         `json:"trigger_y"`     // 觸發位置 Y
	Phase         string          `json:"phase"`         // 當前階段
	PenetrateIDs  []string        `json:"penetrate_ids"` // 穿透路徑上的目標 ID（drill_start 時）
	ExplodeIDs    []string        `json:"explode_ids"`   // 爆炸範圍內的目標 ID（explosion 時）
	KilledTargets []DrillKillEntry `json:"killed_targets"` // 所有被擊破的目標（result 時）
	TotalReward   int             `json:"total_reward"`  // 總獎勵（result 時）
	KillerID      string          `json:"killer_id"`     // 觸發玩家 ID
	KillerName    string          `json:"killer_name"`   // 觸發玩家名稱
}

// ---- 炸彈蟹連環爆炸系統（DAY-143）----

// BombCrabKillEntry 炸彈蟹連帶擊破記錄
type BombCrabKillEntry struct {
	InstanceID string  `json:"instance_id"` // 被擊破的目標 InstanceID
	DefID      string  `json:"def_id"`      // 目標定義 ID
	Multiplier float64 `json:"multiplier"`  // 目標倍率
	Reward     int     `json:"reward"`      // 獎勵金幣
	WaveIndex  int     `json:"wave_index"`  // 第幾波爆炸（0/1/2）
}

// BombCrabChainPayload 炸彈蟹連環爆炸廣播（Server → Client，DAY-143）
// Phase: "bomb_start" → "explosion"（×3波）→ "result"
type BombCrabChainPayload struct {
	TriggerID     string              `json:"trigger_id"`     // 觸發的 T107 InstanceID
	TriggerX      float64             `json:"trigger_x"`      // 觸發位置 X
	TriggerY      float64             `json:"trigger_y"`      // 觸發位置 Y
	Phase         string              `json:"phase"`          // 當前階段
	WaveIndex     int                 `json:"wave_index"`     // 當前波次（0/1/2）
	TotalWaves    int                 `json:"total_waves"`    // 總波數（3）
	ExplodeIDs    []string            `json:"explode_ids"`    // 本波爆炸範圍內的目標 ID
	KilledTargets []BombCrabKillEntry `json:"killed_targets"` // 所有被擊破的目標（result 時）
	TotalReward   int                 `json:"total_reward"`   // 總獎勵（result 時）
	KillerID      string              `json:"killer_id"`      // 觸發玩家 ID
	KillerName    string              `json:"killer_name"`    // 觸發玩家名稱
}

// ---- 巨型章魚轉盤系統（DAY-144）----

// OctopusWheelSlotPayload 巨型章魚轉盤格子資訊
type OctopusWheelSlotPayload struct {
	Index      int    `json:"index"`      // 格子索引（0-7）
	Multiplier int    `json:"multiplier"` // 倍率
	Color      string `json:"color"`      // 顯示顏色（hex）
	Label      string `json:"label"`      // 顯示文字
}

// MegaOctopusWheelStartPayload 轉盤開始通知（Server → Client，個人）
type MegaOctopusWheelStartPayload struct {
	TriggerID    string                    `json:"trigger_id"`    // 觸發的 T108 InstanceID
	SpinDuration int                       `json:"spin_duration"` // 旋轉時間（秒）
	Slots        []OctopusWheelSlotPayload `json:"slots"`         // 轉盤格子定義
}

// MegaOctopusWheelResultPayload 轉盤結果（Server → Client，個人）
type MegaOctopusWheelResultPayload struct {
	ResultIndex int    `json:"result_index"` // 結果格子索引
	Multiplier  int    `json:"multiplier"`   // 獲得倍率
	Reward      int    `json:"reward"`       // 獲得金幣
	NewBalance  int    `json:"new_balance"`  // 新餘額
	SlotLabel   string `json:"slot_label"`   // 格子文字（如 "950x 👑"）
	SlotColor   string `json:"slot_color"`   // 格子顏色
}

// ---- 巨型鮟鱇魚電擊寶箱系統（DAY-145）----

// AnglerfishChestEntry 電擊開啟的寶箱記錄
type AnglerfishChestEntry struct {
	InstanceID string  `json:"instance_id"` // 寶箱 InstanceID
	Multiplier float64 `json:"multiplier"`  // 寶箱倍率
	Reward     int     `json:"reward"`      // 獎勵金幣
	X          float64 `json:"x"`           // 寶箱位置 X
	Y          float64 `json:"y"`           // 寶箱位置 Y
}

// AnglerfishShockPayload 鮟鱇魚電擊廣播（Server → Client，DAY-145）
// Phase: "shock_start" → "result"
type AnglerfishShockPayload struct {
	TriggerID    string                 `json:"trigger_id"`    // 觸發的 T109 InstanceID
	TriggerX     float64                `json:"trigger_x"`     // 觸發位置 X
	TriggerY     float64                `json:"trigger_y"`     // 觸發位置 Y
	Phase        string                 `json:"phase"`         // 當前階段
	ChestIDs     []string               `json:"chest_ids"`     // 電擊範圍內的寶箱 ID
	OpenedChests []AnglerfishChestEntry `json:"opened_chests"` // 開啟的寶箱（result 時）
	TotalReward  int                    `json:"total_reward"`  // 總獎勵（result 時）
	KillerID     string                 `json:"killer_id"`     // 觸發玩家 ID
	KillerName   string                 `json:"killer_name"`   // 觸發玩家名稱
}

// ---- 巨型鹹水鱷魚獵魚累積系統（DAY-146）----

// CrocodileHuntEntry 鱷魚獵殺記錄
type CrocodileHuntEntry struct {
	InstanceID string  `json:"instance_id"` // 被獵殺的目標 InstanceID
	DefID      string  `json:"def_id"`      // 目標定義 ID
	Multiplier float64 `json:"multiplier"`  // 目標倍率
	Reward     int     `json:"reward"`      // 獎勵金幣
	HuntIndex  int     `json:"hunt_index"`  // 第幾次獵殺（0-5）
}

// CrocodileHuntPayload 鱷魚獵魚廣播（Server → Client，DAY-146）
// Phase: "awaken" → "hunt"（×N次）→ "result"
type CrocodileHuntPayload struct {
	TriggerID     string               `json:"trigger_id"`     // 觸發的 T110 InstanceID
	TriggerX      float64              `json:"trigger_x"`      // 觸發位置 X
	TriggerY      float64              `json:"trigger_y"`      // 觸發位置 Y
	Phase         string               `json:"phase"`          // 當前階段
	HuntDuration  int                  `json:"hunt_duration"`  // 獵魚持續時間（秒）
	MaxHunts      int                  `json:"max_hunts"`      // 最多獵殺數
	HuntIndex     int                  `json:"hunt_index"`     // 當前獵殺次數（hunt 階段）
	HuntedID      string               `json:"hunted_id"`      // 本次獵殺的目標 ID（hunt 階段）
	HuntReward    int                  `json:"hunt_reward"`    // 本次獵殺獎勵（hunt 階段）
	HuntedTargets []CrocodileHuntEntry `json:"hunted_targets"` // 所有獵殺記錄（result 時）
	TotalReward   int                  `json:"total_reward"`   // 總獎勵（result 時）
	KillerID      string               `json:"killer_id"`      // 觸發玩家 ID
	KillerName    string               `json:"killer_name"`    // 觸發玩家名稱
}

// GiantPrizeFishPayload 夢幻巨型獎勵魚廣播（Server → Client，DAY-147）
// Phase: "activate" → "end"
// 觸發玩家在 10 秒內所有擊破獎勵 ×5
type GiantPrizeFishPayload struct {
	TriggerID    string  `json:"trigger_id"`    // 觸發的 T111 InstanceID
	TriggerX     float64 `json:"trigger_x"`     // 觸發位置 X
	TriggerY     float64 `json:"trigger_y"`     // 觸發位置 Y
	Phase        string  `json:"phase"`         // "activate" 或 "end"
	MultBonus    float64 `json:"mult_bonus"`    // 倍率加成（5.0）
	Duration     int     `json:"duration"`      // 持續時間（秒，10）
	KillerID     string  `json:"killer_id"`     // 觸發玩家 ID
	KillerName   string  `json:"killer_name"`   // 觸發玩家名稱
	TotalReward  int     `json:"total_reward"`  // 夢幻模式期間總獎勵（end 時）
	KillCount    int     `json:"kill_count"`    // 夢幻模式期間擊破數（end 時）
}

// ChainLongWheelStartPayload 千龍王輪盤開始（Server → Client，個人，DAY-148）
// 觸發玩家看到全螢幕輪盤，其他玩家看到廣播通知
type ChainLongWheelStartPayload struct {
	InstanceID  string  `json:"instance_id"`  // 千龍王 InstanceID
	KillerID    string  `json:"killer_id"`    // 觸發玩家 ID
	KillerName  string  `json:"killer_name"`  // 觸發玩家名稱
	TargetMult  float64 `json:"target_mult"`  // 千龍王本身的倍率
	BaseReward  int     `json:"base_reward"`  // 千龍王擊破的基礎獎勵
	InnerSlots  []float64 `json:"inner_slots"` // 內環選項（5x/10x/20x/30x/50x）
	OuterSlots  []float64 `json:"outer_slots"` // 外環選項（2x/3x/5x/7x/10x/20x）
	SpinSecs    float64 `json:"spin_secs"`    // 旋轉持續秒數
	Message     string  `json:"message"`      // 廣播訊息
}

// ChainLongWheelResultPayload 千龍王輪盤結果（Server → Client，DAY-148）
// 個人：完整結果；全服：只顯示玩家名稱和最終倍率
type ChainLongWheelResultPayload struct {
	KillerID    string  `json:"killer_id"`    // 觸發玩家 ID
	KillerName  string  `json:"killer_name"`  // 觸發玩家名稱
	TargetMult  float64 `json:"target_mult"`  // 千龍王本身的倍率
	InnerResult float64 `json:"inner_result"` // 內環結果
	OuterResult float64 `json:"outer_result"` // 外環結果
	Combined    float64 `json:"combined"`     // 最終倍率（內 × 外）
	BaseReward  int     `json:"base_reward"`  // 基礎獎勵
	BonusReward int     `json:"bonus_reward"` // 額外獎勵（基礎 × 最終倍率）
	NewBalance  int     `json:"new_balance"`  // 新餘額（個人）
	IsMegaWin   bool    `json:"is_mega_win"`  // 是否大獎（≥200x）
	IsPersonal  bool    `json:"is_personal"`  // 是否為觸發玩家（個人通知）
	Message     string  `json:"message"`      // 廣播訊息
}

// ChainLongWheelStatusPayload 千龍王輪盤冷卻狀態（Server → Client，個人，DAY-148）
type ChainLongWheelStatusPayload struct {
	CooldownLeft int  `json:"cooldown_left"` // 冷卻剩餘秒數（0=可觸發）
	HasActive    bool `json:"has_active"`    // 是否有活躍 session
}

// GoldenJellyfishShockEntry 黃金水母電擊單個目標的記錄
type GoldenJellyfishShockEntry struct {
	TargetInstanceID string  `json:"target_instance_id"`
	TargetDefID      string  `json:"target_def_id"`
	TargetName       string  `json:"target_name"`
	Killed           bool    `json:"killed"`
	Multiplier       float64 `json:"multiplier"`
	Reward           int     `json:"reward"`
	ShockIndex       int     `json:"shock_index"` // 第幾個被電擊（0-based）
}

// GoldenJellyfishShockPayload 黃金水母全場電擊廣播（Server → Client，全服，DAY-149）
// Phase: "shock_start" → "shock_N"（逐一電擊）→ "result"
type GoldenJellyfishShockPayload struct {
	TriggerID    string                       `json:"trigger_id"`    // 觸發的 T113 InstanceID
	TriggerX     float64                      `json:"trigger_x"`     // 觸發位置 X
	TriggerY     float64                      `json:"trigger_y"`     // 觸發位置 Y
	Phase        string                       `json:"phase"`         // "shock_start" / "shock" / "result"
	KillerID     string                       `json:"killer_id"`     // 觸發玩家 ID
	KillerName   string                       `json:"killer_name"`   // 觸發玩家名稱
	Targets      []GoldenJellyfishShockEntry  `json:"targets"`       // 電擊目標列表
	TotalKills   int                          `json:"total_kills"`   // 總擊破數（result 時）
	TotalReward  int                          `json:"total_reward"`  // 總獎勵（result 時）
	NewBalance   int                          `json:"new_balance"`   // 新餘額（result 時，個人）
	Message      string                       `json:"message"`       // 廣播訊息
}

// ThunderboltLobsterActivatePayload 雷霆龍蝦免費射擊模式開始（Server → Client，個人+全服，DAY-150）
type ThunderboltLobsterActivatePayload struct {
	TriggerID    string `json:"trigger_id"`    // 觸發的 T114 InstanceID
	TriggerX     float64 `json:"trigger_x"`   // 觸發位置 X
	TriggerY     float64 `json:"trigger_y"`   // 觸發位置 Y
	KillerID     string `json:"killer_id"`     // 觸發玩家 ID
	KillerName   string `json:"killer_name"`   // 觸發玩家名稱
	Duration     int    `json:"duration"`      // 免費射擊持續秒數（15）
	ShotInterval int    `json:"shot_interval"` // 自動射擊間隔毫秒（500）
	Message      string `json:"message"`       // 廣播訊息
}

// ThunderboltLobsterShotPayload 雷霆龍蝦自動射擊一次（Server → Client，全服，DAY-150）
type ThunderboltLobsterShotPayload struct {
	KillerID     string  `json:"killer_id"`     // 觸發玩家 ID
	KillerName   string  `json:"killer_name"`   // 觸發玩家名稱
	TargetID     string  `json:"target_id"`     // 被射擊的目標 InstanceID
	TargetDefID  string  `json:"target_def_id"` // 目標 DefID
	TargetName   string  `json:"target_name"`   // 目標名稱
	TargetX      float64 `json:"target_x"`      // 目標位置 X
	TargetY      float64 `json:"target_y"`      // 目標位置 Y
	IsKill       bool    `json:"is_kill"`        // 是否擊破
	Reward       int     `json:"reward"`         // 獎勵（擊破時）
	Multiplier   float64 `json:"multiplier"`     // 目標倍率
	ShotIndex    int     `json:"shot_index"`     // 第幾次射擊（0-based）
	ShotsLeft    int     `json:"shots_left"`     // 剩餘射擊次數
}

// ThunderboltLobsterEndPayload 雷霆龍蝦免費射擊模式結束（Server → Client，個人+全服，DAY-150）
type ThunderboltLobsterEndPayload struct {
	KillerID     string `json:"killer_id"`     // 觸發玩家 ID
	KillerName   string `json:"killer_name"`   // 觸發玩家名稱
	TotalShots   int    `json:"total_shots"`   // 總射擊次數
	TotalKills   int    `json:"total_kills"`   // 總擊破數
	TotalReward  int    `json:"total_reward"`  // 總獎勵
	NewBalance   int    `json:"new_balance"`   // 新餘額（個人）
	Message      string `json:"message"`       // 廣播訊息
}

// RainbowPhoenixActivatePayload 彩虹鳳凰 Power Up 開始（Server → Client，個人+全服，DAY-151）
type RainbowPhoenixActivatePayload struct {
	TriggerID   string  `json:"trigger_id"`    // 觸發的 T115 InstanceID
	TriggerX    float64 `json:"trigger_x"`     // 觸發位置 X
	TriggerY    float64 `json:"trigger_y"`     // 觸發位置 Y
	KillerID    string  `json:"killer_id"`     // 觸發玩家 ID
	KillerName  string  `json:"killer_name"`   // 觸發玩家名稱
	PowerUpMult float64 `json:"power_up_mult"` // Power Up 倍率（6x-10x）
	Duration    int     `json:"duration"`      // 持續秒數（8）
	Message     string  `json:"message"`       // 廣播訊息
}

// RainbowPhoenixEndPayload 彩虹鳳凰 Power Up 結束（Server → Client，個人+全服，DAY-151）
type RainbowPhoenixEndPayload struct {
	KillerID    string  `json:"killer_id"`     // 觸發玩家 ID
	KillerName  string  `json:"killer_name"`   // 觸發玩家名稱
	PowerUpMult float64 `json:"power_up_mult"` // Power Up 倍率
	TotalKills  int     `json:"total_kills"`   // Power Up 期間擊破數
	TotalReward int     `json:"total_reward"`  // Power Up 期間總獎勵
	NewBalance  int     `json:"new_balance"`   // 新餘額（個人）
	Message     string  `json:"message"`       // 廣播訊息
}

// RainbowPhoenixStatusPayload 彩虹鳳凰 Power Up 狀態（Server → Client，個人，DAY-151）
type RainbowPhoenixStatusPayload struct {
	IsActive    bool    `json:"is_active"`     // 是否在 Power Up 模式中
	PowerUpMult float64 `json:"power_up_mult"` // 當前 Power Up 倍率（0 = 無）
}

// VampireGrowPayload 吸血鬼倍率成長（Server → Client，全服，DAY-152）
type VampireGrowPayload struct {
	InstanceID   string  `json:"instance_id"`    // 吸血鬼 InstanceID
	HitCount     int     `json:"hit_count"`      // 當前命中次數
	MultBonus    float64 `json:"mult_bonus"`     // 當前倍率加成（1.0/2.0/3.5/5.0）
	NewMult      float64 `json:"new_mult"`       // 當前實際倍率
	PhaseName    string  `json:"phase_name"`     // 階段名稱（沉睡/覺醒/狂暴/血月）
	PhaseChanged bool    `json:"phase_changed"`  // 是否剛進入新階段
}

// VampireBloodMoonPayload 血月模式觸發（Server → Client，全服，DAY-152）
type VampireBloodMoonPayload struct {
	InstanceID string  `json:"instance_id"` // 吸血鬼 InstanceID
	HitCount   int     `json:"hit_count"`   // 當前命中次數
	MultBonus  float64 `json:"mult_bonus"`  // 血月倍率加成（5.0）
	NewMult    float64 `json:"new_mult"`    // 當前實際倍率
	Message    string  `json:"message"`     // 廣播訊息
}

// VampireKilledPayload 吸血鬼被擊破（Server → Client，全服，DAY-152）
type VampireKilledPayload struct {
	InstanceID  string  `json:"instance_id"`   // 吸血鬼 InstanceID
	KillerID    string  `json:"killer_id"`     // 擊破玩家 ID
	KillerName  string  `json:"killer_name"`   // 擊破玩家名稱
	HitCount    int     `json:"hit_count"`     // 最終命中次數
	MultBonus   float64 `json:"mult_bonus"`    // 最終倍率加成
	FinalMult   float64 `json:"final_mult"`    // 最終實際倍率
	FinalReward int     `json:"final_reward"`  // 最終獎勵
	PhaseName   string  `json:"phase_name"`    // 最終階段名稱
	Message     string  `json:"message"`       // 廣播訊息
}

// ---- 水晶龍收集大獎系統（DAY-153）----

// CrystalDragonDropPayload 水晶掉落（Server → Client，全服）（DAY-153）
type CrystalDragonDropPayload struct {
	KillerID      string  `json:"killer_id"`      // 擊破玩家 ID
	KillerName    string  `json:"killer_name"`    // 擊破玩家名稱
	CrystalsGain  int     `json:"crystals_gain"`  // 本次掉落水晶數量
	TotalCrystals int     `json:"total_crystals"` // 全服水晶總數
	Goal          int     `json:"goal"`           // 目標水晶數量
	Progress      float64 `json:"progress"`       // 進度（0.0-1.0）
	Message       string  `json:"message"`        // 廣播訊息
}

// CrystalDragonUpdatePayload 水晶進度更新（Server → Client，全服）（DAY-153）
type CrystalDragonUpdatePayload struct {
	TotalCrystals int     `json:"total_crystals"` // 全服水晶總數
	Goal          int     `json:"goal"`           // 目標水晶數量
	Progress      float64 `json:"progress"`       // 進度（0.0-1.0）
}

// CrystalContributorEntry 水晶貢獻者（DAY-153）
type CrystalContributorEntry struct {
	PlayerID   string `json:"player_id"`   // 玩家 ID
	PlayerName string `json:"player_name"` // 玩家名稱
	Crystals   int    `json:"crystals"`    // 貢獻水晶數量
	Reward     int    `json:"reward"`      // 獲得獎勵
}

// CrystalDragonRewardPayload 地獄龍大獎（Server → Client，全服）（DAY-153）
type CrystalDragonRewardPayload struct {
	Contributors []CrystalContributorEntry `json:"contributors"` // 貢獻者列表
	TotalReward  int                       `json:"total_reward"` // 總獎勵（所有貢獻者合計）
	Message      string                    `json:"message"`      // 廣播訊息
}

// CrystalDragonStatusPayload 水晶狀態（Server → Client，個人，登入時）（DAY-153）
type CrystalDragonStatusPayload struct {
	TotalCrystals int     `json:"total_crystals"` // 全服水晶總數
	Goal          int     `json:"goal"`           // 目標水晶數量
	Progress      float64 `json:"progress"`       // 進度（0.0-1.0）
	CooldownSecs  int     `json:"cooldown_secs"`  // 冷卻剩餘秒數
}

// ---- 深淵巨鯨全服 Boss 挑戰系統（DAY-164）----

// AbyssWhaleEntry 深淵巨鯨貢獻者記錄
type AbyssWhaleEntry struct {
	Rank       int     `json:"rank"`        // 排名（1-based）
	PlayerID   string  `json:"player_id"`   // 玩家 ID
	PlayerName string  `json:"player_name"` // 玩家名稱
	Damage     int     `json:"damage"`      // 累積傷害
	Ratio      float64 `json:"ratio"`       // 貢獻比例（0.0-1.0）
	Bonus      int     `json:"bonus"`       // 獲得獎勵
}

// AbyssWhalePayload 深淵巨鯨全服 Boss 挑戰廣播（Server → Client，DAY-164）
// Phase: "whale_spawn" → "whale_hp_update"（多次）→ "whale_killed" → "whale_reward"（個人）
type AbyssWhalePayload struct {
	Phase       string            `json:"phase"`        // 當前階段
	InstanceID  string            `json:"instance_id"`  // 深淵巨鯨 InstanceID
	X           float64           `json:"x"`            // 位置 X
	Y           float64           `json:"y"`            // 位置 Y
	TotalHP     int               `json:"total_hp"`     // 總 HP（500）
	CurrentHP   int               `json:"current_hp"`   // 當前 HP
	HPPercent   float64           `json:"hp_percent"`   // HP 百分比（0.0-1.0）
	AttackerID  string            `json:"attacker_id"`  // 最後攻擊者 ID（hp_update 時）
	KillerID    string            `json:"killer_id"`    // 擊破玩家 ID（killed/reward 時）
	KillerName  string            `json:"killer_name"`  // 擊破玩家名稱
	Entries     []AbyssWhaleEntry `json:"entries"`      // 貢獻者列表（killed/reward 時）
	TotalDamage int               `json:"total_damage"` // 總傷害（killed 時）
	MyRank      int               `json:"my_rank"`      // 我的排名（reward 時）
	MyBonus     int               `json:"my_bonus"`     // 我的獎勵（reward 時）
	MyDamage    int               `json:"my_damage"`    // 我的傷害（reward 時）
	MyRatio     float64           `json:"my_ratio"`     // 我的貢獻比例（reward 時）
}

// ---- 獅子舞大獎爆發系統（DAY-168）----

// LionDanceMarkedTarget 獅子舞標記目標
type LionDanceMarkedTarget struct {
	InstanceID string  `json:"instance_id"` // 被標記的目標 InstanceID
	X          float64 `json:"x"`           // 目標位置 X
	Y          float64 `json:"y"`           // 目標位置 Y
}

// LionDanceBurstPayload 獅子舞爆發廣播（Server → Client，DAY-168）
// Phase: "burst_start" → "burst_end"
type LionDanceBurstPayload struct {
	Phase            string                  `json:"phase"`             // 當前階段
	TriggerPlayer    string                  `json:"trigger_player"`    // 觸發玩家 ID
	TriggerName      string                  `json:"trigger_name"`      // 觸發玩家名稱
	BurstMult        float64                 `json:"burst_mult"`        // 爆發倍率（3-10x）
	MarkedTargets    []LionDanceMarkedTarget `json:"marked_targets"`    // 標記目標列表（burst_start 時）
	DurationSec      int                     `json:"duration_sec"`      // 持續時間（秒）
	RemainingTargets int                     `json:"remaining_targets"` // 剩餘未擊破的標記目標數（burst_end 時）
}

// ---- 漩渦魚群吸引系統（DAY-169）----

// VortexKillEntry 漩渦吸入的目標記錄
type VortexKillEntry struct {
	InstanceID string  `json:"instance_id"` // 被吸入的目標 InstanceID
	DefID      string  `json:"def_id"`      // 目標定義 ID
	Multiplier float64 `json:"multiplier"`  // 目標倍率
	Reward     int     `json:"reward"`      // 獲得獎勵
	X          float64 `json:"x"`           // 目標位置 X
	Y          float64 `json:"y"`           // 目標位置 Y
}

// VortexFishPayload 漩渦魚群廣播（Server → Client，DAY-169）
// Phase: "vortex_start" → "vortex_suck"（多次）→ "vortex_end"
type VortexFishPayload struct {
	Phase        string           `json:"phase"`         // 當前階段
	TriggerID    string           `json:"trigger_id"`    // 觸發玩家 ID
	TriggerName  string           `json:"trigger_name"`  // 觸發玩家名稱
	VortexX      float64          `json:"vortex_x"`      // 漩渦中心 X（漩渦魚位置）
	VortexY      float64          `json:"vortex_y"`      // 漩渦中心 Y
	GroupName    string           `json:"group_name"`    // 目標群組名稱（例如「基礎目標群」）
	TargetCount  int              `json:"target_count"`  // 預計吸入目標數
	SuckIndex    int              `json:"suck_index"`    // 當前吸入的目標索引（vortex_suck 時）
	SuckEntry    *VortexKillEntry `json:"suck_entry"`    // 當前被吸入的目標（vortex_suck 時）
	KilledCount  int              `json:"killed_count"`  // 實際擊破數（vortex_end 時）
	TotalReward  int              `json:"total_reward"`  // 總獎勵（vortex_end 時）
	KilledEntries []VortexKillEntry `json:"killed_entries"` // 所有被擊破的目標（vortex_end 時）
}

// ---- 冰凍炸彈魚系統（DAY-170）----

// FreezeBombEntry 被冰凍的特殊目標記錄
type FreezeBombEntry struct {
	InstanceID string  `json:"instance_id"` // 被冰凍的目標 InstanceID
	DefID      string  `json:"def_id"`      // 目標定義 ID
	X          float64 `json:"x"`           // 目標位置 X
	Y          float64 `json:"y"`           // 目標位置 Y
}

// FreezeBombPayload 冰凍炸彈魚廣播（Server → Client，DAY-170）
// Phase: "freeze_start" → "freeze_end"
type FreezeBombPayload struct {
	Phase         string            `json:"phase"`          // 當前階段
	TriggerID     string            `json:"trigger_id"`     // 觸發玩家 ID
	TriggerName   string            `json:"trigger_name"`   // 觸發玩家名稱
	FreezeX       float64           `json:"freeze_x"`       // 冰凍炸彈位置 X
	FreezeY       float64           `json:"freeze_y"`       // 冰凍炸彈位置 Y
	FrozenCount   int               `json:"frozen_count"`   // 被冰凍的目標數
	DurationSec   int               `json:"duration_sec"`   // 冰凍持續時間（秒）
	FrozenTargets []FreezeBombEntry `json:"frozen_targets"` // 被冰凍的目標列表（freeze_start 時）
}

// ---- 冰釣幸運輪盤系統（DAY-171）----

// IceFishingWheelPayload 冰釣幸運輪盤廣播（Server → Client，DAY-171）
// Phase: "wheel_start"（個人）→ "wheel_broadcast"（全服）→ "wheel_result"（個人）
//        → "wheel_result_broadcast"（全服，≥5x）→ "mult_end"（個人）
type IceFishingWheelPayload struct {
	Phase       string  `json:"phase"`        // 當前階段
	PlayerID    string  `json:"player_id"`    // 觸發玩家 ID
	PlayerName  string  `json:"player_name"`  // 觸發玩家名稱
	WheelResult int     `json:"wheel_result"` // 輪盤結果格子索引（wheel_result 時）
	Multiplier  float64 `json:"multiplier"`   // 輪盤倍率（2x-10x）
	Label       string  `json:"label"`        // 格子顯示文字（×2-×10）
	Color       string  `json:"color"`        // 格子顏色
	SpinSec     int     `json:"spin_sec"`     // 旋轉時間（秒，wheel_start 時）
	DurationSec int     `json:"duration_sec"` // 倍率持續時間（秒，wheel_result 時）
	KillCount   int     `json:"kill_count"`   // 倍率期間擊破數（mult_end 時）
	TotalBonus  int     `json:"total_bonus"`  // 倍率期間額外獎勵（mult_end 時）
}

// MsgIceFishingWheelStop Client → Server：玩家手動停止輪盤
// （使用 MsgUseSpecialWeapon 的 action 欄位，或獨立訊息）

// ---- 幸運彩蛋魚系統（DAY-172）----

// LuckyEggResult 單個彩蛋開啟結果
type LuckyEggResult struct {
	EggIndex    int     `json:"egg_index"`    // 彩蛋索引（0-4）
	RewardType  string  `json:"reward_type"`  // 獎勵類型：coins/mult/weapon
	CoinsReward int     `json:"coins_reward"` // 金幣獎勵（coins 類型）
	MultBoost   float64 `json:"mult_boost"`   // 倍率加成（mult 類型，2.0）
	DurationSec int     `json:"duration_sec"` // 倍率持續時間（mult 類型，5秒）
	Label       string  `json:"label"`        // 顯示文字
	Color       string  `json:"color"`        // 顯示顏色
}

// LuckyEggFishPayload 幸運彩蛋魚廣播（Server → Client，DAY-172）
// Phase: "egg_start"（全服）→ "egg_open"（個人，每個彩蛋）→ "egg_result"（個人）
//        → "egg_broadcast"（全服，≥4個）→ "mult_end"（個人，倍率結束）
type LuckyEggFishPayload struct {
	Phase       string          `json:"phase"`        // 當前階段
	PlayerID    string          `json:"player_id"`    // 觸發玩家 ID
	PlayerName  string          `json:"player_name"`  // 觸發玩家名稱
	EggCount    int             `json:"egg_count"`    // 彩蛋總數
	EggIndex    int             `json:"egg_index"`    // 當前彩蛋索引（egg_open 時）
	EggResult   LuckyEggResult  `json:"egg_result"`   // 當前彩蛋結果（egg_open 時）
	EggResults  []LuckyEggResult `json:"egg_results"` // 所有彩蛋結果（egg_result 時）
	TotalCoins  int             `json:"total_coins"`  // 金幣彩蛋總獎勵
	MultCount   int             `json:"mult_count"`   // 倍率彩蛋數量
	WeaponCount int             `json:"weapon_count"` // 武器充能彩蛋數量
	TriggerX    float64         `json:"trigger_x"`    // 觸發位置 X（egg_start 時）
	TriggerY    float64         `json:"trigger_y"`    // 觸發位置 Y（egg_start 時）
}

// ---- 彩虹幸運魚系統（DAY-173）----

// RainbowLuckyFishPayload 彩虹幸運魚廣播（Server → Client，DAY-173）
// Phase: "lucky_start"（全服）→ "lucky_end"（全服）
type RainbowLuckyFishPayload struct {
	Phase       string  `json:"phase"`        // 當前階段：lucky_start/lucky_end
	PlayerName  string  `json:"player_name"`  // 觸發玩家名稱
	DurationSec int     `json:"duration_sec"` // 持續時間（秒，lucky_start 時）
	KillBoost   float64 `json:"kill_boost"`   // 擊破機率加成（0.20 = +20%，lucky_start 時）
	TriggerX    float64 `json:"trigger_x"`    // 觸發位置 X（lucky_start 時）
	TriggerY    float64 `json:"trigger_y"`    // 觸發位置 Y（lucky_start 時）
}

// ---- 海葵觸手攻擊系統（DAY-174）----

// SeaAnemoneHitEntry 觸手命中記錄
type SeaAnemoneHitEntry struct {
	InstanceID string  `json:"instance_id"` // 被命中目標 ID
	DefID      string  `json:"def_id"`      // 目標定義 ID
	X          float64 `json:"x"`           // 目標位置 X
	Y          float64 `json:"y"`           // 目標位置 Y
	Multiplier float64 `json:"multiplier"`  // 目標倍率
	Direction  int     `json:"direction"`   // 觸手方向（0-7）
	Angle      float64 `json:"angle"`       // 觸手角度（度）
	IsKill     bool    `json:"is_kill"`     // 是否擊破
	Reward     int     `json:"reward"`      // 獎勵金幣
}

// SeaAnemonePayload 海葵觸手攻擊廣播（Server → Client，DAY-174）
// Phase: "tentacle_start"（全服）→ "tentacle_hit"/"tentacle_miss"（全服，每個方向）
//        → "tentacle_result"（全服）
type SeaAnemonePayload struct {
	Phase       string               `json:"phase"`        // 當前階段
	TriggerID   string               `json:"trigger_id"`   // 觸發目標 ID
	TriggerX    float64              `json:"trigger_x"`    // 觸發位置 X
	TriggerY    float64              `json:"trigger_y"`    // 觸發位置 Y
	KillerID    string               `json:"killer_id"`    // 觸發玩家 ID
	KillerName  string               `json:"killer_name"`  // 觸發玩家名稱
	Directions  int                  `json:"directions"`   // 觸手方向數（tentacle_start 時）
	Direction   int                  `json:"direction"`    // 當前觸手方向（tentacle_hit/miss 時）
	Angle       float64              `json:"angle"`        // 當前觸手角度（度）
	HitID       string               `json:"hit_id"`       // 命中目標 ID（tentacle_hit 時）
	HitX        float64              `json:"hit_x"`        // 命中目標位置 X
	HitY        float64              `json:"hit_y"`        // 命中目標位置 Y
	IsKill      bool                 `json:"is_kill"`      // 是否擊破（tentacle_hit 時）
	Reward      int                  `json:"reward"`       // 獎勵金幣（tentacle_hit 時）
	Multiplier  float64              `json:"multiplier"`   // 目標倍率（tentacle_hit 時）
	HitEntries  []SeaAnemoneHitEntry `json:"hit_entries"`  // 所有命中記錄（tentacle_result 時）
	KillCount   int                  `json:"kill_count"`   // 擊破數（tentacle_result 時）
	TotalReward int                  `json:"total_reward"` // 總獎勵（tentacle_result 時）
}

// ---- 幸運骰子魚系統（DAY-175）----

// LuckyDiceFishPayload 幸運骰子魚廣播（Server → Client，DAY-175）
// Phase: "dice_start"（個人）→ "dice_broadcast"（全服）→ "dice_result"（個人）
//        → "dice_jackpot"（全服，點數12時）
type LuckyDiceFishPayload struct {
	Phase      string  `json:"phase"`       // 當前階段
	PlayerID   string  `json:"player_id"`   // 觸發玩家 ID
	PlayerName string  `json:"player_name"` // 觸發玩家名稱
	Die1       int     `json:"die1"`        // 骰子1點數（1-6，dice_result 時）
	Die2       int     `json:"die2"`        // 骰子2點數（1-6，dice_result 時）
	Sum        int     `json:"sum"`         // 點數之和（2-12，dice_result 時）
	Reward     int     `json:"reward"`      // 獎勵金幣（dice_result 時）
	Label      string  `json:"label"`       // 結果標籤（dice_result 時）
	NewBalance int     `json:"new_balance"` // 新餘額（dice_result 時）
	RollMs     int     `json:"roll_ms"`     // 骰子滾動時間（ms，dice_start 時）
	TriggerX   float64 `json:"trigger_x"`   // 觸發位置 X（dice_start 時）
	TriggerY   float64 `json:"trigger_y"`   // 觸發位置 Y（dice_start 時）
}

// ---- 火焰風暴魚系統（DAY-176）----

// FireStormFishPayload 火焰風暴魚廣播（Server → Client，DAY-176）
// Phase: "fire_start"（全服）→ "fire_burn"（全服，每個目標）→ "fire_end"（全服）
type FireStormFishPayload struct {
	Phase       string   `json:"phase"`        // 當前階段
	PlayerID    string   `json:"player_id"`    // 觸發玩家 ID
	PlayerName  string   `json:"player_name"`  // 觸發玩家名稱
	TargetCount int      `json:"target_count"` // 標記目標數（fire_start 時）
	TargetIDs   []string `json:"target_ids"`   // 標記目標 ID 列表（fire_start 時）
	DurationSec int      `json:"duration_sec"` // 持續時間（fire_start 時）
	TargetID    string   `json:"target_id"`    // 當前燃燒目標 ID（fire_burn 時）
	Reward      int      `json:"reward"`       // 燃燒獎勵（fire_burn 時）
	Skipped     bool     `json:"skipped"`      // 目標已消失（fire_burn 時）
	BurnedCount int      `json:"burned_count"` // 成功燃燒數（fire_end 時）
	TotalReward int      `json:"total_reward"` // 總獎勵（fire_end 時）
}

// ---- 黃金寶藏魚系統（DAY-177）----

// GoldenTreasureFishPayload 黃金寶藏魚廣播（Server → Client，DAY-177）
// Phase: "treasure_start"（個人）→ "treasure_broadcast"（全服）
//        → "treasure_open"（個人，玩家開箱）→ "treasure_mult_start"（個人）
//        → "treasure_mult_end"（個人）→ "treasure_auto_open"（個人，超時自動開）
//        → "treasure_end"（個人）→ "treasure_weapon_charge"（個人）
type GoldenTreasureFishPayload struct {
	Phase           string  `json:"phase"`             // 當前階段
	PlayerID        string  `json:"player_id"`         // 觸發玩家 ID
	PlayerName      string  `json:"player_name"`       // 觸發玩家名稱
	ChestCount      int     `json:"chest_count"`       // 寶藏箱數量（treasure_start 時）
	TimeoutSec      int     `json:"timeout_sec"`       // 超時時間（treasure_start 時）
	ChestID         int     `json:"chest_id"`          // 箱子編號（treasure_open 時）
	RewardType      string  `json:"reward_type"`       // 獎勵類型（coins/mult/weapon）
	Reward          int     `json:"reward"`            // 金幣獎勵（treasure_open 時）
	MultActivated   bool    `json:"mult_activated"`    // 倍率是否激活（treasure_open 時）
	MultBonus       float64 `json:"mult_bonus"`        // 倍率加成值（treasure_mult_start 時）
	MultDurationSec int     `json:"mult_duration_sec"` // 倍率持續時間（treasure_mult_start 時）
	WeaponCharge    int     `json:"weapon_charge"`     // 武器充能量（treasure_weapon_charge 時）
	IsAuto          bool    `json:"is_auto"`           // 是否自動開啟（treasure_auto_open 時）
}

// GoldenTreasureOpenPayload 玩家開箱請求（Client → Server，DAY-177）
type GoldenTreasureOpenPayload struct {
	ChestID int `json:"chest_id"` // 箱子編號（0-2）
}

// ---- 美人魚治癒系統（DAY-178）----

// MermaidHealingPayload 美人魚治癒廣播（Server → Client，DAY-178）
// Phase: "heal_start"（個人）→ "heal_broadcast"（全服）
//        → "luck_start"（全服）→ "luck_end"（全服）
type MermaidHealingPayload struct {
	Phase                string  `json:"phase"`                  // 當前階段
	PlayerID             string  `json:"player_id"`              // 觸發玩家 ID
	PlayerName           string  `json:"player_name"`            // 觸發玩家名稱
	HealAmount           int     `json:"heal_amount"`            // 治癒金幣數（heal_start 時）
	NewBalance           int     `json:"new_balance"`            // 新餘額（heal_start 時）
	LuckBoostPercent     float64 `json:"luck_boost_percent"`     // 幸運加成比例（luck_start 時）
	LuckBoostDurationSec int     `json:"luck_boost_duration_sec"` // 幸運加成持續時間（luck_start 時）
}

// ---- 幸運草魚系統（DAY-179）----

// LuckyCloverFishPayload 幸運草魚廣播（Server → Client，DAY-179）
// Phase: "clover_start"（全服）→ "clover_gift"（個人）→ "clover_end"（全服）
type LuckyCloverFishPayload struct {
	Phase            string  `json:"phase"`              // 當前階段
	PlayerID         string  `json:"player_id"`          // 觸發玩家 ID
	PlayerName       string  `json:"player_name"`        // 觸發玩家名稱
	BoostPercent     float64 `json:"boost_percent"`      // 加成比例（clover_start 時）
	BoostDurationSec int     `json:"boost_duration_sec"` // 加成持續時間（clover_start 時）
	GiftAmount       int     `json:"gift_amount"`        // 幸運草金幣數（clover_gift 時）
	GiftMult         int     `json:"gift_mult"`          // 幸運草金幣倍率（clover_gift 時）
	NewBalance       int     `json:"new_balance"`        // 新餘額（clover_gift 時）
}

// ---- 彩虹鯊魚爆發系統（DAY-180）----

// RainbowSharkMarkedTargetPayload 彩虹爆發標記目標（用於 payload）
type RainbowSharkMarkedTargetPayload struct {
	InstanceID string  `json:"instance_id"`
	DefID      string  `json:"def_id"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	BurstMult  float64 `json:"burst_mult"` // 1.5/2.0/2.5/3.0
}

// RainbowSharkBurstPayload 彩虹鯊魚爆發廣播（Server → Client，DAY-180）
// Phase: "burst_start"（全服）→ "burst_end"（全服）
type RainbowSharkBurstPayload struct {
	Phase           string                              `json:"phase"`             // 當前階段
	TriggerPlayerID string                              `json:"trigger_player_id"` // 觸發玩家 ID
	TriggerName     string                              `json:"trigger_name"`      // 觸發玩家名稱
	MarkedTargets   []RainbowSharkMarkedTargetPayload   `json:"marked_targets"`    // 標記目標列表
	DurationSec     int                                 `json:"duration_sec"`      // 持續時間（秒）
}

// ---- 雷霆鯊魚連鎖閃電系統（DAY-181）----

// ThunderSharkChainPayload 雷霆鯊魚連鎖閃電廣播（Server → Client，DAY-181）
// Phase: "chain_start" → "jump_N"（每跳）→ "result"
type ThunderSharkChainPayload struct {
	Phase       string  `json:"phase"`        // 當前階段
	TriggerID   string  `json:"trigger_id"`   // 觸發目標 ID
	TriggerX    float64 `json:"trigger_x"`    // 觸發位置 X
	TriggerY    float64 `json:"trigger_y"`    // 觸發位置 Y
	JumpTarget  string  `json:"jump_target"`  // 跳躍目標 ID（jump_N 時）
	JumpX       float64 `json:"jump_x"`       // 跳躍目標位置 X
	JumpY       float64 `json:"jump_y"`       // 跳躍目標位置 Y
	JumpNum     int     `json:"jump_num"`     // 當前跳數（jump_N 時）
	TotalJumps  int     `json:"total_jumps"`  // 總跳數（result 時）
	TotalKills  int     `json:"total_kills"`  // 總擊破數（result 時）
	TotalReward int     `json:"total_reward"` // 總獎勵（result 時）
	KillerID    string  `json:"killer_id"`    // 觸發玩家 ID
	KillerName  string  `json:"killer_name"`  // 觸發玩家名稱
}

// ---- 吸血鬼魚累積倍率系統（DAY-182）----

// VampireFishPayload 吸血鬼魚倍率廣播（Server → Client，DAY-182）
// Phase: "vampire_start"（個人）→ "vampire_broadcast"（全服）→ "mult_update"（個人）→ "vampire_end"（個人）
type VampireFishPayload struct {
	Phase       string  `json:"phase"`        // 當前階段
	PlayerID    string  `json:"player_id"`    // 玩家 ID
	PlayerName  string  `json:"player_name"`  // 玩家名稱
	CurrentMult float64 `json:"current_mult"` // 當前倍率
	MaxMult     float64 `json:"max_mult"`     // 最高倍率（vampire_start 時）
	DurationSec int     `json:"duration_sec"` // 持續時間（vampire_start 時）
	KillCount   int     `json:"kill_count"`   // 擊破數（vampire_end 時）
}

// ---- 閃電魚自動連鎖系統（DAY-183）----

// LightningAutoChainPayload 閃電魚自動連鎖廣播（Server → Client，DAY-183）
// Phase: "chain_start" → "auto_N"（每次自動攻擊）→ "result"
type LightningAutoChainPayload struct {
	Phase        string  `json:"phase"`         // 當前階段
	TriggerID    string  `json:"trigger_id"`    // 觸發目標 ID
	TriggerX     float64 `json:"trigger_x"`     // 觸發位置 X
	TriggerY     float64 `json:"trigger_y"`     // 觸發位置 Y
	TargetID     string  `json:"target_id"`     // 攻擊目標 ID（auto_N 時）
	TargetX      float64 `json:"target_x"`      // 攻擊目標位置 X
	TargetY      float64 `json:"target_y"`      // 攻擊目標位置 Y
	AttackNum    int     `json:"attack_num"`    // 當前攻擊次數（auto_N 時）
	TotalAttacks int     `json:"total_attacks"` // 總攻擊次數（result 時）
	TotalKills   int     `json:"total_kills"`   // 總擊破數（result 時）
	TotalReward  int     `json:"total_reward"`  // 總獎勵（result 時）
	KillerID     string  `json:"killer_id"`     // 觸發玩家 ID
	KillerName   string  `json:"killer_name"`   // 觸發玩家名稱
	DurationSec  int     `json:"duration_sec"`  // 持續時間（chain_start 時）
}

// ---- 隕石魚隕石雨系統（DAY-184）----

// MeteorFishPayload 隕石魚隕石雨廣播（Server → Client，DAY-184）
// Phase: "meteor_start" → "meteor_N"（每顆隕石落點）→ "meteor_result"
type MeteorFishPayload struct {
	Phase       string  `json:"phase"`        // 當前階段
	TriggerID   string  `json:"trigger_id"`   // 觸發目標 ID
	TriggerX    float64 `json:"trigger_x"`    // 觸發位置 X
	TriggerY    float64 `json:"trigger_y"`    // 觸發位置 Y
	TargetID    string  `json:"target_id"`    // 命中目標 ID（meteor_N 時）
	TargetX     float64 `json:"target_x"`     // 命中目標位置 X
	TargetY     float64 `json:"target_y"`     // 命中目標位置 Y
	MeteorNum   int     `json:"meteor_num"`   // 當前隕石編號（meteor_N 時）
	MeteorCount int     `json:"meteor_count"` // 總隕石數（meteor_start/result 時）
	TotalKills  int     `json:"total_kills"`  // 總擊破數（result 時）
	TotalReward int     `json:"total_reward"` // 總獎勵（result 時）
	KillerID    string  `json:"killer_id"`    // 觸發玩家 ID
	KillerName  string  `json:"killer_name"`  // 觸發玩家名稱
	IsBoss      bool    `json:"is_boss"`      // 是否命中 BOSS（meteor_N 時）
}

// ---- 鳳凰魚涅槃重生系統（DAY-185）----

// PhoenixFishPayload 鳳凰魚涅槃重生廣播（Server → Client，DAY-185）
// Phase: "phoenix_explode" → "phoenix_rebirth" → "rebirth_end"
type PhoenixFishPayload struct {
	Phase       string  `json:"phase"`        // 當前階段
	TriggerID   string  `json:"trigger_id"`   // 觸發目標 ID
	TriggerX    float64 `json:"trigger_x"`    // 觸發位置 X
	TriggerY    float64 `json:"trigger_y"`    // 觸發位置 Y
	TotalKills  int     `json:"total_kills"`  // 總擊破數（rebirth 時）
	TotalReward int     `json:"total_reward"` // 總獎勵（rebirth 時）
	KillerID    string  `json:"killer_id"`    // 觸發玩家 ID
	KillerName  string  `json:"killer_name"`  // 觸發玩家名稱
	BoostPct    int     `json:"boost_pct"`    // 重生加成百分比（rebirth 時，30）
	BoostSec    int     `json:"boost_sec"`    // 重生加成持續秒數（rebirth 時，30）
}

// ---- 龍龜不死 Boss 系統（DAY-186）----

// DragonTurtlePayload 龍龜不死 Boss 廣播（Server → Client，DAY-186）
// Phase: "turtle_appear" → "turtle_hit" / "my_hit" → "turtle_leave"
type DragonTurtlePayload struct {
	Phase       string  `json:"phase"`                  // 當前階段
	InstanceID  string  `json:"instance_id"`            // 龍龜 InstanceID
	X           float64 `json:"x,omitempty"`            // 出現位置 X（turtle_appear）
	Y           float64 `json:"y,omitempty"`            // 出現位置 Y（turtle_appear）
	HitterID    string  `json:"hitter_id,omitempty"`    // 命中玩家 ID（turtle_hit）
	HitterName  string  `json:"hitter_name,omitempty"`  // 命中玩家名稱（turtle_hit）
	HitReward   int     `json:"hit_reward,omitempty"`   // 本次命中獎勵（turtle_hit/my_hit）
	HitMult     int     `json:"hit_mult,omitempty"`     // 本次命中倍率（turtle_hit/my_hit）
	TotalHits   int     `json:"total_hits,omitempty"`   // 全服總命中數
	TotalReward int     `json:"total_reward,omitempty"` // 全服總獎勵（turtle_leave）
}

// ---- 連鎖爆炸魚系統（DAY-187）----

// ChainBombPayload 連鎖爆炸魚廣播（Server → Client，DAY-187）
// Phase: "chain_start" → "chain_explode"(×N) → "chain_result"
type ChainBombPayload struct {
	Phase       string  `json:"phase"`                  // 當前階段
	TriggerID   string  `json:"trigger_id,omitempty"`   // 觸發目標 InstanceID
	TriggerX    float64 `json:"trigger_x,omitempty"`    // 觸發位置 X（chain_start/chain_explode）
	TriggerY    float64 `json:"trigger_y,omitempty"`    // 觸發位置 Y（chain_start/chain_explode）
	KillerID    string  `json:"killer_id,omitempty"`    // 觸發玩家 ID
	KillerName  string  `json:"killer_name,omitempty"`  // 觸發玩家名稱
	ChainDepth  int     `json:"chain_depth,omitempty"`  // 當前連鎖層數（chain_explode/chain_result）
	KillCount   int     `json:"kill_count,omitempty"`   // 本層擊破數（chain_explode）
	Reward      int     `json:"reward,omitempty"`       // 本層獎勵（chain_explode）
	TotalKills  int     `json:"total_kills,omitempty"`  // 總擊破數（chain_result）
	TotalReward int     `json:"total_reward,omitempty"` // 總獎勵（chain_result）
}

// ---- 巨型鱷魚獵食系統（DAY-188）----

// CrocodileHunterPayload 巨型鱷魚獵食廣播（Server → Client，DAY-188）
// Phase: "croc_appear" → "croc_hunt"(×N) / "croc_miss" → "croc_killed" / "croc_leave"
type CrocodileHunterPayload struct {
	Phase       string  `json:"phase"`                  // 當前階段
	InstanceID  string  `json:"instance_id,omitempty"`  // 鱷魚 InstanceID
	KillerID    string  `json:"killer_id,omitempty"`    // 擊破玩家 ID（croc_killed）
	KillerName  string  `json:"killer_name,omitempty"`  // 擊破玩家名稱（croc_killed）
	HuntIndex   int     `json:"hunt_index,omitempty"`   // 第幾次獵食（croc_hunt）
	HuntCount   int     `json:"hunt_count,omitempty"`   // 總獵食次數（croc_killed/croc_leave）
	MaxHunts    int     `json:"max_hunts,omitempty"`    // 最大獵食次數（croc_appear）
	TargetName  string  `json:"target_name,omitempty"`  // 被獵食目標名稱（croc_hunt/croc_miss）
	TargetMult  float64 `json:"target_mult,omitempty"`  // 被獵食目標倍率（croc_hunt）
	TargetX     float64 `json:"target_x,omitempty"`     // 被獵食目標位置 X（croc_hunt/croc_miss）
	TargetY     float64 `json:"target_y,omitempty"`     // 被獵食目標位置 Y（croc_hunt/croc_miss）
	HuntReward  int     `json:"hunt_reward,omitempty"`  // 本次獵食獎勵（croc_hunt）
	TotalPool   int     `json:"total_pool,omitempty"`   // 累積獎池（croc_hunt/croc_killed/croc_leave）
	PoolBonus   int     `json:"pool_bonus,omitempty"`   // 獎池加成（croc_killed）
	BaseReward  int     `json:"base_reward,omitempty"`  // 基礎獎勵（croc_killed）
	TotalReward int     `json:"total_reward,omitempty"` // 總獎勵（croc_killed）
	NewBalance  int     `json:"new_balance,omitempty"`  // 玩家新餘額（croc_killed）
	Message     string  `json:"message,omitempty"`      // 廣播訊息
}

// ---- 時間炸彈魚系統（DAY-189）----

// TimeBombFishPayload 時間炸彈魚廣播（Server → Client，DAY-189）
// Phase: "bomb_appear" → "bomb_tick"(×N) → "bomb_defused" / "bomb_explode" → "bomb_result" / "defuse_end"
type TimeBombFishPayload struct {
	Phase         string `json:"phase"`                    // 當前階段
	InstanceID    string `json:"instance_id,omitempty"`    // 炸彈魚 InstanceID
	Countdown     int    `json:"countdown,omitempty"`      // 剩餘秒數（bomb_appear/bomb_tick）
	KillerID      string `json:"killer_id,omitempty"`      // 拆彈玩家 ID（bomb_defused）
	KillerName    string `json:"killer_name,omitempty"`    // 拆彈玩家名稱（bomb_defused）
	BaseReward    int    `json:"base_reward,omitempty"`    // 基礎擊破獎勵（bomb_defused）
	NewBalance    int    `json:"new_balance,omitempty"`    // 玩家新餘額（bomb_defused）
	BonusPct      int    `json:"bonus_pct,omitempty"`      // 拆彈加成百分比（bomb_defused）
	BonusDuration int    `json:"bonus_duration,omitempty"` // 拆彈加成持續秒數（bomb_defused）
	KillCount     int    `json:"kill_count,omitempty"`     // 爆炸擊破數（bomb_result）
	TotalReward   int    `json:"total_reward,omitempty"`   // 爆炸總獎勵（bomb_result）
	Message       string `json:"message,omitempty"`        // 廣播訊息
}

// TripleLuckyFishPayload 三重幸運魚廣播（Server → Client，DAY-190）
// Phase: "triple_start"（個人詳細）/ "triple_broadcast"（全服廣播）/ "mult_end"（倍率結束）
type TripleLuckyFishPayload struct {
	Phase         string  `json:"phase"`                    // 當前階段
	PlayerID      string  `json:"player_id,omitempty"`      // 觸發玩家 ID
	PlayerName    string  `json:"player_name,omitempty"`    // 觸發玩家名稱
	CoinReward    int     `json:"coin_reward,omitempty"`    // 金幣雨獎勵
	CoinMult      float64 `json:"coin_mult,omitempty"`      // 金幣雨倍率（20-50x）
	MultBonus     float64 `json:"mult_bonus,omitempty"`     // 倍率加成（0.5 = +50%）
	MultDuration  float64 `json:"mult_duration,omitempty"`  // 倍率加成持續秒數
	MultEndUnix   int64   `json:"mult_end_unix,omitempty"`  // 倍率加成結束時間（Unix timestamp）
	WeaponCharged string  `json:"weapon_charged,omitempty"` // 充能的武器名稱（龍怒/魚雷/軌道炮）
	NewBalance    int     `json:"new_balance,omitempty"`    // 玩家新餘額（triple_start）
	Message       string  `json:"message,omitempty"`        // 廣播訊息
}

// SchoolPanicPayload 魚群驚嚇廣播（Server → Client，DAY-191）
// Phase: "panic_start" → "panic_end"
type SchoolPanicPayload struct {
	Phase        string   `json:"phase"`                     // 當前階段
	TriggerID    string   `json:"trigger_id,omitempty"`      // 觸發的魚群領袖 InstanceID
	KillerID     string   `json:"killer_id,omitempty"`       // 觸發玩家 ID
	KillerName   string   `json:"killer_name,omitempty"`     // 觸發玩家名稱
	PanicTargets []string `json:"panic_targets,omitempty"`   // 受驚嚇的目標 InstanceID 列表
	TargetCount  int      `json:"target_count,omitempty"`    // 受驚嚇目標數量
	Duration     float64  `json:"duration,omitempty"`        // 驚嚇持續秒數
	Message      string   `json:"message,omitempty"`         // 廣播訊息
}

// RockSkeletonConcertPayload 搖滾骷髏演唱會廣播（Server → Client，DAY-192）
// Phase: "concert_start" → "note_N" → "beat_result" → "awakening" → "encore_start"/"concert_end" → "encore_end"
type RockSkeletonConcertPayload struct {
	Phase           string    `json:"phase"`                        // 當前階段
	TriggerID       string    `json:"trigger_id,omitempty"`         // 觸發的搖滾骷髏魚 InstanceID
	TriggerX        float64   `json:"trigger_x,omitempty"`          // 觸發位置 X
	TriggerY        float64   `json:"trigger_y,omitempty"`          // 觸發位置 Y
	KillerID        string    `json:"killer_id,omitempty"`          // 觸發玩家 ID
	KillerName      string    `json:"killer_name,omitempty"`        // 觸發玩家名稱
	DurationSec     int       `json:"duration_sec,omitempty"`       // 演唱會持續秒數
	Beat            int       `json:"beat,omitempty"`               // 當前拍數（1-15）
	NoteTargetIDs   []string  `json:"note_target_ids,omitempty"`    // 音符炸彈目標 InstanceID 列表
	NoteTargetXs    []float64 `json:"note_target_xs,omitempty"`     // 音符炸彈目標 X 座標列表
	NoteTargetYs    []float64 `json:"note_target_ys,omitempty"`     // 音符炸彈目標 Y 座標列表
	BeatKills       int       `json:"beat_kills,omitempty"`         // 本拍擊破數
	BeatReward      int       `json:"beat_reward,omitempty"`        // 本拍獎勵
	TotalKills      int       `json:"total_kills,omitempty"`        // 累計擊破數
	TotalReward     int       `json:"total_reward,omitempty"`       // 累計獎勵
	IsAwakening     bool      `json:"is_awakening,omitempty"`       // 是否在超級覺醒狀態
	AwakenedTargets []string  `json:"awakened_targets,omitempty"`   // 超級覺醒影響的目標 InstanceID 列表
	AwakenedCount   int       `json:"awakened_count,omitempty"`     // 超級覺醒影響目標數量
	EncoreDuration  int       `json:"encore_duration,omitempty"`    // 安可加成持續秒數
	EncoreBonus     float64   `json:"encore_bonus,omitempty"`       // 安可加成比例（0.30 = +30%）
	Message         string    `json:"message,omitempty"`            // 廣播訊息
}

// ElectricLinkResult 電流連接結果（DAY-193）
type ElectricLinkResult struct {
	IDA        string  `json:"id_a"`                    // 目標 A InstanceID
	IDB        string  `json:"id_b"`                    // 目標 B InstanceID
	XA         float64 `json:"x_a"`                     // 目標 A X 座標
	YA         float64 `json:"y_a"`                     // 目標 A Y 座標
	XB         float64 `json:"x_b"`                     // 目標 B X 座標
	YB         float64 `json:"y_b"`                     // 目標 B Y 座標
	IsKill     bool    `json:"is_kill"`                 // 是否擊破
	KilledID   string  `json:"killed_id,omitempty"`     // 被擊破的目標 InstanceID
	KilledMult float64 `json:"killed_mult,omitempty"`   // 被擊破目標的倍率
	Reward     int     `json:"reward,omitempty"`        // 獎勵金幣
}

// ElectricJellyfishPayload 電流水母電流網路廣播（Server → Client，DAY-193）
// Phase: "network_start" → "link_N" → "network_result"
type ElectricJellyfishPayload struct {
	Phase       string               `json:"phase"`                    // 當前階段
	TriggerID   string               `json:"trigger_id,omitempty"`     // 觸發的電流水母 InstanceID
	TriggerX    float64              `json:"trigger_x,omitempty"`      // 觸發位置 X
	TriggerY    float64              `json:"trigger_y,omitempty"`      // 觸發位置 Y
	KillerID    string               `json:"killer_id,omitempty"`      // 觸發玩家 ID
	KillerName  string               `json:"killer_name,omitempty"`    // 觸發玩家名稱
	LinkCount   int                  `json:"link_count,omitempty"`     // 電流連接總數
	LinkIndex   int                  `json:"link_index,omitempty"`     // 當前連接序號（1-based）
	IDA         string               `json:"id_a,omitempty"`           // 目標 A InstanceID
	IDB         string               `json:"id_b,omitempty"`           // 目標 B InstanceID
	XA          float64              `json:"x_a,omitempty"`            // 目標 A X 座標
	YA          float64              `json:"y_a,omitempty"`            // 目標 A Y 座標
	XB          float64              `json:"x_b,omitempty"`            // 目標 B X 座標
	YB          float64              `json:"y_b,omitempty"`            // 目標 B Y 座標
	IsKill      bool                 `json:"is_kill,omitempty"`        // 是否擊破
	KilledID    string               `json:"killed_id,omitempty"`      // 被擊破的目標 InstanceID
	KilledMult  float64              `json:"killed_mult,omitempty"`    // 被擊破目標的倍率
	Reward      int                  `json:"reward,omitempty"`         // 本次連接獎勵
	TotalKills  int                  `json:"total_kills,omitempty"`    // 累計擊破數
	TotalReward int                  `json:"total_reward,omitempty"`   // 累計獎勵
	Links       []ElectricLinkResult `json:"links,omitempty"`          // 所有連接結果（result 階段）
}

// ChainLongKingPayload 長龍王雙環輪盤廣播（Server → Client，DAY-194）
// Phase: "roulette_start" → "inner_stop" → "outer_stop" → "result"
// 千倍大獎: "mega_win" → "mega_broadcast"
// 全服廣播: "broadcast"（≥100x）
type ChainLongKingPayload struct {
	Phase       string `json:"phase"`                    // 當前階段
	InstanceID  string `json:"instance_id,omitempty"`    // 觸發的長龍王 InstanceID
	PlayerName  string `json:"player_name,omitempty"`    // 觸發玩家名稱（全服廣播用）
	InnerRing   []int  `json:"inner_ring,omitempty"`     // 內環倍率定義（roulette_start 時發送）
	OuterRing   []int  `json:"outer_ring,omitempty"`     // 外環乘數定義（roulette_start 時發送）
	InnerResult int    `json:"inner_result,omitempty"`   // 內環停止結果
	OuterResult int    `json:"outer_result,omitempty"`   // 外環停止結果
	TotalMult   int    `json:"total_mult,omitempty"`     // 最終倍率（內環 × 外環）
	Reward      int    `json:"reward,omitempty"`         // 獎勵金幣
	IsBigWin    bool   `json:"is_big_win,omitempty"`     // 是否大獎（≥100x）
	IsMega      bool   `json:"is_mega,omitempty"`        // 是否千倍大獎
	IsTimeout   bool   `json:"is_timeout,omitempty"`     // 是否超時自動停止
}

// ChainLongKingStopPayload 玩家停止輪盤（Client → Server，DAY-194）
type ChainLongKingStopPayload struct {
	InstanceID string `json:"instance_id"` // 輪盤 InstanceID（防止重複停止）
}

// DrillLobsterPayload 鑽頭龍蝦穿透爆炸廣播（Server → Client，DAY-195）
// Phase: "drill_start" → "drill_1"..."drill_5" → "drill_explode" → "drill_result"
type DrillLobsterPayload struct {
	Phase          string  `json:"phase"`                     // 當前階段
	TriggerID      string  `json:"trigger_id,omitempty"`      // 觸發的鑽頭龍蝦 InstanceID
	KillerID       string  `json:"killer_id,omitempty"`       // 觸發玩家 ID
	KillerName     string  `json:"killer_name,omitempty"`     // 觸發玩家名稱
	StartX         float64 `json:"start_x,omitempty"`         // 鑽頭出發 X
	StartY         float64 `json:"start_y,omitempty"`         // 鑽頭出發 Y
	DirX           float64 `json:"dir_x,omitempty"`           // 鑽頭方向 X（單位向量）
	DirY           float64 `json:"dir_y,omitempty"`           // 鑽頭方向 Y（單位向量）
	StepIndex      int     `json:"step_index,omitempty"`      // 當前步驟序號（1-5）
	CurX           float64 `json:"cur_x,omitempty"`           // 當前位置 X
	CurY           float64 `json:"cur_y,omitempty"`           // 當前位置 Y
	IsKill         bool    `json:"is_kill,omitempty"`         // 本步是否擊破目標
	KilledID       string  `json:"killed_id,omitempty"`       // 被擊破目標 InstanceID
	KilledName     string  `json:"killed_name,omitempty"`     // 被擊破目標名稱
	KilledMult     float64 `json:"killed_mult,omitempty"`     // 被擊破目標倍率
	StepReward     int     `json:"step_reward,omitempty"`     // 本步獎勵
	TotalKills     int     `json:"total_kills,omitempty"`     // 累計擊破數
	ExplodeRadius  float64 `json:"explode_radius,omitempty"`  // 爆炸半徑
	ExplodeKills   int     `json:"explode_kills,omitempty"`   // 爆炸擊破數
	ExplodeReward  int     `json:"explode_reward,omitempty"`  // 爆炸獎勵
	PenetrateCount int     `json:"penetrate_count,omitempty"` // 穿透命中目標數
	TotalReward    int     `json:"total_reward,omitempty"`    // 總獎勵
}

// AnglerfishElectricPayload 巨型鮟鱇魚電擊寶箱廣播（Server → Client，DAY-196）
// Phase: "anglerfish_appear" → "zap_1"..."zap_8" / "super_zap_start"/"super_zap_N"/"super_zap_result" → "anglerfish_killed"/"anglerfish_leave"
type AnglerfishElectricPayload struct {
	Phase        string  `json:"phase"`                   // 當前階段
	InstanceID   string  `json:"instance_id,omitempty"`   // 鮟鱇魚 InstanceID
	ZapIndex     int     `json:"zap_index,omitempty"`     // 電擊序號（1-8）
	TargetID     string  `json:"target_id,omitempty"`     // 被電擊目標 InstanceID
	TargetDefID  string  `json:"target_def_id,omitempty"` // 被電擊目標定義 ID
	TargetX      float64 `json:"target_x,omitempty"`      // 被電擊目標位置 X
	TargetY      float64 `json:"target_y,omitempty"`      // 被電擊目標位置 Y
	IsKill       bool    `json:"is_kill,omitempty"`       // 是否擊破目標
	IsTreasure   bool    `json:"is_treasure,omitempty"`   // 是否為寶箱開箱
	TreasureMult float64 `json:"treasure_mult,omitempty"` // 寶箱開箱倍率（3-5x）
	ZapReward    int     `json:"zap_reward,omitempty"`    // 本次電擊獎勵
	IsSuperZap   bool    `json:"is_super_zap,omitempty"`  // 是否為超級電擊
	IsEmpty      bool    `json:"is_empty,omitempty"`      // 場上無目標（空電擊）
	TargetCount  int     `json:"target_count,omitempty"`  // 超級電擊目標數
	SuperKills   int     `json:"super_kills,omitempty"`   // 超級電擊擊破數
	SuperReward  int     `json:"super_reward,omitempty"`  // 超級電擊總獎勵
	KillerID     string  `json:"killer_id,omitempty"`     // 擊破鮟鱇魚的玩家 ID
	KillerName   string  `json:"killer_name,omitempty"`   // 擊破鮟鱇魚的玩家名稱
	ZapCount     int     `json:"zap_count,omitempty"`     // 累計電擊次數
	TotalPool    int     `json:"total_pool,omitempty"`    // 累積電擊獎池
	PoolBonus    int     `json:"pool_bonus,omitempty"`    // 獎池加成（40%）
	BaseReward   int     `json:"base_reward,omitempty"`   // 基礎倍率獎勵
	TotalReward  int     `json:"total_reward,omitempty"`  // 總獎勵
}

// MysticDragonPayload 神秘龍魚八波龍息攻擊廣播（Server → Client，DAY-197）
// Phase: "dragon_start" → "wave_1"..."wave_8" → "dragon_result"
type MysticDragonPayload struct {
	Phase       string `json:"phase"`                  // 當前階段
	TriggerID   string `json:"trigger_id,omitempty"`   // 觸發的神秘龍魚 InstanceID
	KillerID    string `json:"killer_id,omitempty"`    // 觸發玩家 ID
	KillerName  string `json:"killer_name,omitempty"`  // 觸發玩家名稱
	TotalWaves  int    `json:"total_waves,omitempty"`  // 總波數（dragon_start 時）
	WaveIndex   int    `json:"wave_index,omitempty"`   // 當前波數（1-8）
	WaveKills   int    `json:"wave_kills,omitempty"`   // 本波擊破數
	WaveReward  int    `json:"wave_reward,omitempty"`  // 本波獎勵
	TotalKills  int    `json:"total_kills,omitempty"`  // 累計擊破數
	TotalReward int    `json:"total_reward,omitempty"` // 累計獎勵
	IsFinalWave bool   `json:"is_final_wave,omitempty"` // 是否為第 8 波（龍怒爆發）
}

// GhostFishPayload 幽靈魚分身廣播（Server → Client，DAY-198）
// Phase: "ghost_appear" → "phantom_vanish"（幻影被擊破）/ "real_found"（真身被找到）→ "ghost_explode" → "ghost_escape"（逃跑）
type GhostFishPayload struct {
	Phase         string   `json:"phase"`                    // 當前階段
	RealID        string   `json:"real_id,omitempty"`        // 真身 InstanceID
	CloneIDs      []string `json:"clone_ids,omitempty"`      // 幻影分身 InstanceID 列表
	CloneID       string   `json:"clone_id,omitempty"`       // 被擊破的幻影分身 ID（phantom_vanish 時）
	CloneCount    int      `json:"clone_count,omitempty"`    // 幻影分身數量（ghost_appear 時）
	KillerID      string   `json:"killer_id,omitempty"`      // 擊破玩家 ID
	KillerName    string   `json:"killer_name,omitempty"`    // 擊破玩家名稱
	Reward        int      `json:"reward,omitempty"`         // 安慰獎（phantom_vanish 時）
	ExplodeKills  int      `json:"explode_kills,omitempty"`  // 幻影爆炸擊破數（ghost_explode 時）
	ExplodeReward int      `json:"explode_reward,omitempty"` // 幻影爆炸獎勵（ghost_explode 時）
}

// ThunderboltLobsterPayload 雷霆龍蝦免費射擊廣播（Server → Client，DAY-199）
// Event: "turret_start" → "turret_shot"（每次自動射擊）→ "turret_end"
type ThunderboltLobsterPayload struct {
	Event       string  `json:"event"`                  // 事件類型
	KillerName  string  `json:"killer_name,omitempty"`  // 觸發者名稱（turret_start 時）
	Duration    float64 `json:"duration,omitempty"`     // 基礎持續時間（turret_start 時）
	MaxDuration float64 `json:"max_duration,omitempty"` // 最大持續時間（turret_start 時）
	TargetID    string  `json:"target_id,omitempty"`    // 射擊目標 ID（turret_shot 時）
	Killed      bool    `json:"killed,omitempty"`       // 是否擊破（turret_shot 時）
	Reward      int64   `json:"reward,omitempty"`       // 本次獎勵（turret_shot 時）
	KillCount   int     `json:"kill_count,omitempty"`   // 累計擊破數
	TotalReward int64   `json:"total_reward,omitempty"` // 累計總獎勵
	Remaining   float64 `json:"remaining,omitempty"`    // 剩餘時間（秒）
	ExtendSec   float64 `json:"extend_sec,omitempty"`   // 已延長秒數（turret_end 時）
}

// IcePhoenixPowerUpResult 冰鳳凰 Power Up 單次攻擊結果
type IcePhoenixPowerUpResult struct {
	TargetID string  `json:"target_id"`
	Mult     float64 `json:"mult"`
	Killed   bool    `json:"killed"`
	Reward   int     `json:"reward"`
}

// IcePhoenixPayload 冰鳳凰覺醒廣播（Server → Client，DAY-200）
// Event: "awaken_start" → "power_up_shot"（每次 Power Up）→ "frost_burst_start"（可選）→ "frost_burst_result"（可選）→ "awaken_result"
type IcePhoenixPayload struct {
	Event         string                   `json:"event"`
	KillerName    string                   `json:"killer_name,omitempty"`
	BaseReward    int                      `json:"base_reward,omitempty"`
	ShotIndex     int                      `json:"shot_index,omitempty"`
	TotalShots    int                      `json:"total_shots,omitempty"`
	PowerUpResult IcePhoenixPowerUpResult  `json:"power_up_result,omitempty"`
	PowerUpKills  int                      `json:"power_up_kills,omitempty"`
	PowerUpReward int                      `json:"power_up_reward,omitempty"`
	FrostKills    int                      `json:"frost_kills,omitempty"`
	FrostReward   int                      `json:"frost_reward,omitempty"`
	TotalReward   int                      `json:"total_reward,omitempty"`
	HasFrost      bool                     `json:"has_frost,omitempty"`
}

// SerialBombCrabPayload 連環炸彈蟹廣播（Server → Client，DAY-201）
// Event: "bomb_start" → "bomb_explode"（每顆炸彈）× N → "bomb_result"
type SerialBombCrabPayload struct {
	Event       string  `json:"event"`
	KillerName  string  `json:"killer_name,omitempty"`
	BombCount   int     `json:"bomb_count,omitempty"`
	BombIndex   int     `json:"bomb_index,omitempty"`
	KillX       float64 `json:"kill_x,omitempty"`
	KillY       float64 `json:"kill_y,omitempty"`
	BombX       float64 `json:"bomb_x,omitempty"`
	BombY       float64 `json:"bomb_y,omitempty"`
	BombKills   int     `json:"bomb_kills,omitempty"`
	BombReward  int     `json:"bomb_reward,omitempty"`
	TotalKills  int     `json:"total_kills,omitempty"`
	TotalReward int     `json:"total_reward,omitempty"`
}

// AbyssVortexPayload 深淵漩渦魚廣播（Server → Client，DAY-202）
// Event: "vortex_start" → "vortex_pulse"（每次脈衝）× 10 → "vortex_blast" → "vortex_result"
type AbyssVortexPayload struct {
	Event       string  `json:"event"`
	KillerName  string  `json:"killer_name,omitempty"`
	VortexX     float64 `json:"vortex_x,omitempty"`
	VortexY     float64 `json:"vortex_y,omitempty"`
	Duration    int     `json:"duration,omitempty"`
	PulseNum    int     `json:"pulse_num,omitempty"`
	PulseKills  int     `json:"pulse_kills,omitempty"`
	PulseReward int     `json:"pulse_reward,omitempty"`
	BlastKills  int     `json:"blast_kills,omitempty"`
	BlastReward int     `json:"blast_reward,omitempty"`
	TotalKills  int     `json:"total_kills,omitempty"`
	TotalReward int     `json:"total_reward,omitempty"`
}

// HumpbackWhalePayload 座頭鯨覺醒廣播（Server → Client，DAY-203）
// Event: "awaken_start" → "wave_attack"（每波）× 3 → ["tidal_wave_start" → "tidal_wave_result"] → "awaken_result"
type HumpbackWhalePayload struct {
	Event       string `json:"event"`
	KillerName  string `json:"killer_name,omitempty"`
	BaseReward  int    `json:"base_reward,omitempty"`
	WaveCount   int    `json:"wave_count,omitempty"`
	WaveNum     int    `json:"wave_num,omitempty"`
	WaveKills   int    `json:"wave_kills,omitempty"`
	WaveReward  int    `json:"wave_reward,omitempty"`
	TidalKills  int    `json:"tidal_kills,omitempty"`
	TidalReward int    `json:"tidal_reward,omitempty"`
	TotalKills  int    `json:"total_kills,omitempty"`
	TotalReward int    `json:"total_reward,omitempty"`
	HasTidal    bool   `json:"has_tidal,omitempty"`
}

// FreeSpinFishPayload 自由旋轉魚免費射擊廣播（Server → Client，DAY-204）
// Event: "free_spin_start" → "free_spin_shot"（每次射擊）× N → "free_spin_end"
// "free_spin_broadcast" 廣播給全服（讓其他玩家知道有人觸發了免費射擊）
type FreeSpinFishPayload struct {
	Event       string  `json:"event"`
	PlayerID    string  `json:"player_id,omitempty"`
	PlayerName  string  `json:"player_name,omitempty"`
	Duration    float64 `json:"duration,omitempty"`
	MaxDuration float64 `json:"max_duration,omitempty"`
	TargetID    string  `json:"target_id,omitempty"`
	TargetX     float64 `json:"target_x,omitempty"`
	TargetY     float64 `json:"target_y,omitempty"`
	Killed      bool    `json:"killed,omitempty"`
	Reward      int     `json:"reward,omitempty"`
	KillCount   int     `json:"kill_count,omitempty"`
	Remaining   float64 `json:"remaining,omitempty"`
	TotalReward int     `json:"total_reward,omitempty"`
	ExtendSec   float64 `json:"extend_sec,omitempty"`
}

// JackpotDragonPayload 獎池龍 Jackpot 抽獎廣播（Server → Client，DAY-205）
// Event: "dragon_draw" — 擊破獎池龍後觸發抽獎，廣播給全服
type JackpotDragonPayload struct {
	Event      string `json:"event"`
	PlayerID   string `json:"player_id,omitempty"`
	PlayerName string `json:"player_name,omitempty"`
	Level      string `json:"level"`       // "mini" / "minor" / "major" / "grand"
	LevelName  string `json:"level_name"`  // "MINI" / "MINOR" / "MAJOR" / "GRAND"
	LevelColor string `json:"level_color"` // 顏色代碼
	LevelIcon  string `json:"level_icon"`  // 圖示
	Amount     int    `json:"amount"`      // 獎勵金額
	IsGrand    bool   `json:"is_grand,omitempty"`
	IsMajor    bool   `json:"is_major,omitempty"`
}
