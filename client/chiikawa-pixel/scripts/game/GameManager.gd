## GameManager.gd
## 遊戲狀態管理，處理 Server 訊息並更新遊戲
## Autoload 單例

extends Node
class_name GameManagerClass

# 訊號
signal game_state_changed(new_state: String)
signal target_spawned(target_data: Dictionary)
signal target_updated(update_data: Dictionary)
signal target_killed(kill_data: Dictionary)
signal attack_result(result: Dictionary)
signal reward_received(reward: Dictionary)
signal player_updated(player_data: Dictionary)
signal boss_event(event_data: Dictionary)
signal bonus_event(event_data: Dictionary)
signal leaderboard_updated(entries: Array)
signal achievement_unlocked(achievement_data: Dictionary)
signal combo_event(combo_data: Dictionary)  # 連擊事件（DAY-022）
signal mission_updated(missions: Array)     # 任務進度更新（DAY-037）
signal mission_completed(mission_data: Dictionary)  # 任務完成（DAY-037）
signal jackpot_updated(jackpot_data: Dictionary)    # Jackpot 池更新（DAY-048）
signal jackpot_won(win_data: Dictionary)            # Jackpot 中獎（DAY-048）
signal jackpot_animation(anim_data: Dictionary)     # Jackpot 觸發動畫（DAY-095）
signal player_stats_updated(stats_data: Dictionary) # 玩家統計更新（DAY-096）
signal announcement_received(ann_data: Dictionary)  # 全服公告（DAY-097）
signal spectator_joined(spectator_data: Dictionary) # 觀戰者加入通知（DAY-054d）
signal daily_bonus_received(bonus_data: Dictionary) # 每日登入獎勵（DAY-065）
signal spectator_left(spectator_data: Dictionary)  # 觀戰者離開通知（DAY-055）
signal tournament_updated(tournament_data: Dictionary) # 週賽排名更新（DAY-066）
signal daily_tournament_updated(tournament_data: Dictionary) # 每日賽排名更新（DAY-093）
signal title_unlocked(title_data: Dictionary)          # 稱號解鎖通知（DAY-068）
signal skin_updated(skin_data: Dictionary)             # 砲台外觀更新（DAY-071）
signal season_updated(season_data: Dictionary)         # 賽季通行證更新（DAY-072）
signal season_level_up(level_data: Dictionary)         # 賽季等級升級（DAY-072）
signal friend_list_updated(friend_data: Dictionary)    # 好友列表更新（DAY-073）
signal friend_request_received(request_data: Dictionary) # 好友請求通知（DAY-073）
signal friend_updated(update_data: Dictionary)         # 好友狀態更新（DAY-073）
# 好友禮物系統（DAY-101）
signal gift_received(gift_data: Dictionary)            # 收到禮物通知
signal gift_sent(gift_data: Dictionary)                # 送出禮物成功通知
signal gift_status(status_data: Dictionary)            # 今日禮物狀態
signal gift_error(error_data: Dictionary)              # 禮物操作失敗
# 好友挑戰系統（DAY-102）
signal challenge_request(request_data: Dictionary)     # 收到挑戰邀請
signal challenge_updated(update_data: Dictionary)      # 挑戰狀態/分數更新
signal challenge_result(result_data: Dictionary)       # 挑戰結果
signal challenge_error(error_data: Dictionary)         # 挑戰操作失敗
# 私訊系統（DAY-103）
signal dm_received(dm_data: Dictionary)                # 收到私訊
signal dm_sent(dm_data: Dictionary)                    # 發送成功確認
signal dm_error(error_data: Dictionary)                # 發送失敗
signal open_dm_panel(friend_id: String, friend_name: String) # 開啟 DM 面板
# 玩家名片系統（DAY-106）
signal player_card_received(card_data: Dictionary)     # 收到玩家名片資料
# 登入里程碑系統（DAY-107）
signal login_milestone_reached(milestone_data: Dictionary) # 里程碑達成通知
signal login_progress_received(progress_data: Dictionary)  # 登入進度回應
# 超級 Bonus 系統（DAY-108）
signal super_bonus_triggered(bonus_data: Dictionary)       # 超級 Bonus 觸發通知
signal guild_updated(guild_data: Dictionary)           # 公會資訊更新（DAY-074）
signal guild_task_complete(task_data: Dictionary)      # 公會任務完成（DAY-074）
signal guild_message_received(msg_data: Dictionary)    # 公會聊天訊息（DAY-075）
signal guild_war_updated(war_data: Dictionary)         # 公會戰排名更新（DAY-076）
signal guild_war_result(result_data: Dictionary)       # 公會戰結算通知（DAY-076）
signal daily_boss_updated(boss_data: Dictionary)       # 每日 BOSS 狀態更新（DAY-077）
signal daily_boss_defeated(defeat_data: Dictionary)    # 每日 BOSS 擊殺通知（DAY-077）
signal vip_updated(vip_data: Dictionary)               # VIP 狀態更新（DAY-078）
signal vip_level_up(level_data: Dictionary)            # VIP 升級通知（DAY-078）
signal vip_weekly_claimed(claim_data: Dictionary)      # VIP 週獎勵領取通知（DAY-078）
signal event_updated(event_data: Dictionary)           # 限時活動狀態更新（DAY-079）
signal codex_updated(codex_data: Dictionary)           # 圖鑑狀態更新（DAY-081）
signal codex_unlocked(unlock_data: Dictionary)         # 圖鑑條目解鎖通知（DAY-081）
signal codex_complete(complete_data: Dictionary)       # 全圖鑑完成通知（DAY-081）
signal streak_updated(streak_data: Dictionary)         # 連擊狀態更新（DAY-083）
signal streak_reset(reset_data: Dictionary)            # 連擊重置通知（DAY-083）
signal referral_info_received(info_data: Dictionary)   # 推薦碼資訊（DAY-082）
signal referral_success(success_data: Dictionary)      # 推薦碼使用成功（DAY-082）
signal referral_error(error_data: Dictionary)          # 推薦碼使用失敗（DAY-082）
signal wheel_triggered(wheel_data: Dictionary)         # 幸運轉盤觸發（DAY-084）
signal challenge_unlocked(challenge_data: Dictionary)  # 隱藏挑戰解鎖（DAY-085）
signal mission_streak_bonus(streak_data: Dictionary)   # 任務連續完成獎勵（DAY-086）
signal weather_updated(weather_data: Dictionary)       # 天氣狀態更新（DAY-087）
signal chain_explosion(chain_data: Dictionary)         # 連鎖爆炸（DAY-088）
signal chain_target_killed(instance_id: String, multiplier: float) # 連鎖目標擊破（DAY-088）
signal special_weapon_updated(weapon_data: Dictionary) # 特殊武器狀態更新（DAY-089）
signal special_weapon_fired(fire_data: Dictionary)     # 特殊武器發射廣播（DAY-089）
signal mystery_box_updated(box_data: Dictionary)       # 神秘寶箱狀態更新（DAY-090）
signal mystery_box_dropped(drop_data: Dictionary)      # 神秘寶箱掉落通知（DAY-090）
signal mystery_box_opened(open_data: Dictionary)       # 神秘寶箱開箱結果（DAY-090）
# 房間難度系統（DAY-091）
signal room_list_received(room_data: Dictionary)       # 房間列表更新
signal room_switched(switch_data: Dictionary)          # 房間切換成功
signal room_error(error_data: Dictionary)              # 房間操作失敗
# 每日簽到轉盤（DAY-092）
signal daily_spin_state(state_data: Dictionary)        # 每日轉盤狀態
signal daily_spin_result(result_data: Dictionary)      # 每日轉盤結果
# 商店系統（DAY-094）
signal shop_updated(shop_data: Dictionary)             # 商店狀態更新
signal shop_purchased(purchase_data: Dictionary)       # 購買成功通知
signal shop_error(error_data: Dictionary)              # 購買失敗通知

# 遊戲狀態
var current_state: String = "normal_play"
var player_data: Dictionary = {}
var targets: Dictionary = {}  # instance_id -> target_data

# 角色顏色對應（規格書 5章）
const CHARACTER_COLORS = {
	"chiikawa": Color(1.0, 0.6, 0.8),   # 粉紅
	"hachiware": Color(0.4, 0.6, 1.0),  # 藍色
	"usagi": Color(1.0, 0.9, 0.2),      # 黃色
}

# 角色名稱
const CHARACTER_NAMES = {
	"chiikawa": "Chiikawa",
	"hachiware": "Hachiware",
	"usagi": "Usagi",
}

func _ready() -> void:
	# 連接 NetworkManager 訊號
	NetworkManager.message_received.connect(_on_message_received)
	NetworkManager.connected.connect(_on_connected)
	NetworkManager.disconnected.connect(_on_disconnected)
	# 啟動資產預載入（背景執行，不阻塞遊戲）
	call_deferred("_start_preloading")

func _start_preloading() -> void:
	if LoadingManager != null:
		LoadingManager.preload_all()

func _on_connected() -> void:
	print("[GameManager] Connected to server")

func _on_disconnected() -> void:
	print("[GameManager] Disconnected from server")

## 處理 Server 訊息
func _on_message_received(type: String, payload: Dictionary) -> void:
	match type:
		"game_state":
			_handle_game_state(payload)
		"target_spawn":
			_handle_target_spawn(payload)
		"target_update":
			_handle_target_update(payload)
		"target_kill":
			_handle_target_kill(payload)
		"attack_result":
			_handle_attack_result(payload)
		"reward":
			_handle_reward(payload)
		"player_update":
			_handle_player_update(payload)
		"boss_event":
			_handle_boss_event(payload)
		"bonus_event":
			_handle_bonus_event(payload)
		"leaderboard":
			_handle_leaderboard(payload)
		"achievement":
			_handle_achievement(payload)
		"combo_event":
			_handle_combo_event(payload)
		"mission_update":
			_handle_mission_update(payload)
		"mission_complete":
			_handle_mission_complete(payload)
		"jackpot_update":
			_handle_jackpot_update(payload)
		"jackpot_win":
			_handle_jackpot_win(payload)
		"jackpot_animation":
			_handle_jackpot_animation(payload)
		"player_stats_update":
			_handle_player_stats_update(payload)
		"announcement":
			_handle_announcement(payload)
		"daily_bonus":
			_handle_daily_bonus(payload)
		"spectator_join":
			_handle_spectator_join(payload)
		"spectator_leave":
			_handle_spectator_leave(payload)
		"tournament_update":
			_handle_tournament_update(payload)
		# 每日賽排名更新（DAY-093）
		"daily_tournament_update":
			_handle_daily_tournament_update(payload)
		"title_unlocked":
			_handle_title_unlocked(payload)
		"skin_update":
			_handle_skin_update(payload)
		"season_update":
			_handle_season_update(payload)
		"season_level_up":
			_handle_season_level_up(payload)
		"friend_list":
			_handle_friend_list(payload)
		"friend_request":
			_handle_friend_request(payload)
		"friend_update":
			_handle_friend_update(payload)
		# 好友禮物系統（DAY-101）
		"gift_received":
			_handle_gift_received(payload)
		"gift_sent":
			_handle_gift_sent(payload)
		"gift_status":
			_handle_gift_status(payload)
		"gift_error":
			_handle_gift_error(payload)
		# 好友挑戰系統（DAY-102）
		"challenge_request":
			_handle_challenge_request(payload)
		"challenge_update":
			_handle_challenge_updated(payload)
		"challenge_result":
			_handle_challenge_result(payload)
		"challenge_error":
			_handle_challenge_error(payload)
		# 私訊系統（DAY-103）
		"dm_received":
			_handle_dm_received(payload)
		"dm_sent":
			_handle_dm_sent(payload)
		"dm_error":
			_handle_dm_error(payload)
		# 玩家名片系統（DAY-106）
		"player_card":
			emit_signal("player_card_received", payload)
		# 登入里程碑系統（DAY-107）
		"login_milestone":
			_handle_login_milestone(payload)
		"login_progress":
			_handle_login_progress(payload)
		# 超級 Bonus 系統（DAY-108）
		"super_bonus_ready":
			_handle_super_bonus_ready(payload)
		"guild_update":
			_handle_guild_update(payload)
		"guild_list":
			_handle_guild_list(payload)
		"guild_task_complete":
			_handle_guild_task_complete(payload)
		"guild_error":
			_handle_guild_error(payload)
		"guild_message":
			_handle_guild_message(payload)
		"guild_war_update":
			_handle_guild_war_update(payload)
		"guild_war_result":
			_handle_guild_war_result(payload)
		"daily_boss_update":
			_handle_daily_boss_update(payload)
		"daily_boss_defeated":
			_handle_daily_boss_defeated(payload)
		"vip_update":
			_handle_vip_update(payload)
		"vip_level_up":
			_handle_vip_level_up(payload)
		"vip_weekly_claimed":
			_handle_vip_weekly_claimed(payload)
		"event_update":
			_handle_event_update(payload)
		"codex_update":
			_handle_codex_update(payload)
		"codex_unlock":
			_handle_codex_unlock(payload)
		"codex_complete":
			_handle_codex_complete(payload)
		"referral_info":
			_handle_referral_info(payload)
		"referral_success":
			_handle_referral_success(payload)
		"referral_error":
			_handle_referral_error(payload)
		"wheel_trigger":
			_handle_wheel_trigger(payload)
		"challenge_unlocked":
			_handle_challenge_unlocked(payload)
		"mission_streak_bonus":
			_handle_mission_streak_bonus(payload)
		"weather_update":
			_handle_weather_update(payload)
		"chain_explosion":
			_handle_chain_explosion(payload)
		"special_weapon_update":
			_handle_special_weapon_update(payload)
		"special_weapon_fired":
			_handle_special_weapon_fired(payload)
		"mystery_box_drop":
			_handle_mystery_box_drop(payload)
		"mystery_box_update":
			_handle_mystery_box_update(payload)
		"mystery_box_opened":
			_handle_mystery_box_opened(payload)
		"streak_update":
			_handle_streak_update(payload)
		"streak_reset":
			_handle_streak_reset(payload)
		# 房間難度系統（DAY-091）
		"room_list":
			_handle_room_list(payload)
		"room_switched":
			_handle_room_switched(payload)
		"room_error":
			_handle_room_error(payload)
		# 每日簽到轉盤（DAY-092）
		"daily_spin_state":
			_handle_daily_spin_state(payload)
		"daily_spin_result":
			_handle_daily_spin_result(payload)
		# 商店系統（DAY-094）
		"shop_update":
			_handle_shop_update(payload)
		"shop_purchased":
			_handle_shop_purchased(payload)
		"shop_error":
			_handle_shop_error(payload)
		"error":
			_handle_error(payload)
		"pong":
			pass  # Ping/Pong 心跳
		# 賽季節日活動系統（DAY-109）
		"festival_update":
			_handle_festival_update(payload)
		"festival_task_ready":
			_handle_festival_task_ready(payload)
		"festival_task_claimed":
			_handle_festival_task_claimed(payload)
		"festival_title_earned":
			_handle_festival_title_earned(payload)
		"festival_error":
			_handle_festival_error(payload)

func _handle_game_state(payload: Dictionary) -> void:
	var new_state = payload.get("state", "")
	if new_state != current_state:
		current_state = new_state
		print("[GameManager] State: ", current_state)
		emit_signal("game_state_changed", current_state)

func _handle_target_spawn(payload: Dictionary) -> void:
	var instance_id = payload.get("instance_id", "")
	targets[instance_id] = payload
	emit_signal("target_spawned", payload)

func _handle_target_update(payload: Dictionary) -> void:
	var instance_id = payload.get("instance_id", "")
	if targets.has(instance_id):
		targets[instance_id].merge(payload, true)
	emit_signal("target_updated", payload)

func _handle_target_kill(payload: Dictionary) -> void:
	var instance_id = payload.get("instance_id", "")
	targets.erase(instance_id)
	emit_signal("target_killed", payload)

func _handle_attack_result(payload: Dictionary) -> void:
	emit_signal("attack_result", payload)

func _handle_reward(payload: Dictionary) -> void:
	# 更新本地金幣顯示
	if payload.has("new_balance"):
		player_data["coins"] = payload["new_balance"]
	emit_signal("reward_received", payload)

func _handle_player_update(payload: Dictionary) -> void:
	player_data = payload
	emit_signal("player_updated", payload)

func _handle_boss_event(payload: Dictionary) -> void:
	emit_signal("boss_event", payload)
	# BGM 切換：BOSS Phase 2 時切換到 boss_rage
	var event = payload.get("event", "")
	if event == "phase_change" and payload.get("phase", 1) == 2:
		if AudioManager != null:
			AudioManager.play_bgm(AudioManager.BGM.BOSS_RAGE)
	elif event == "kill":
		# BOSS 擊敗：短暫靜音，等 boss_result 狀態切換回主 BGM
		if AudioManager != null:
			AudioManager.stop_bgm_briefly()

func _handle_bonus_event(payload: Dictionary) -> void:
	emit_signal("bonus_event", payload)

func _handle_leaderboard(payload: Dictionary) -> void:
	var entries = payload.get("entries", [])
	emit_signal("leaderboard_updated", entries)

func _handle_achievement(payload: Dictionary) -> void:
	print("[GameManager] Achievement unlocked: ", payload.get("name", ""))
	emit_signal("achievement_unlocked", payload)

func _handle_combo_event(payload: Dictionary) -> void:
	## 連擊事件（DAY-022）
	var combo_count = payload.get("combo_count", 1)
	print("[GameManager] COMBO x%d!" % combo_count)
	emit_signal("combo_event", payload)

## 任務進度更新（DAY-037）
func _handle_mission_update(payload: Dictionary) -> void:
	var missions = payload.get("missions", [])
	var reset_at = payload.get("reset_at", 0)
	emit_signal("mission_updated", missions)
	# 傳遞重置時間給 HUD（DAY-038）
	if reset_at > 0:
		var hud = get_node_or_null("/root/Main/HUD")
		if is_instance_valid(hud) and hud.has_method("set_mission_reset_at"):
			hud.set_mission_reset_at(reset_at)

## 任務完成通知（DAY-037）
func _handle_mission_complete(payload: Dictionary) -> void:
	print("[GameManager] Mission completed: ", payload.get("name", ""))
	emit_signal("mission_completed", payload)

func _handle_error(payload: Dictionary) -> void:
	push_warning("[GameManager] Server error: " + str(payload))

## Jackpot 池更新（DAY-048）
func _handle_jackpot_update(payload: Dictionary) -> void:
	emit_signal("jackpot_updated", payload)

## Jackpot 中獎通知（DAY-048）
func _handle_jackpot_win(payload: Dictionary) -> void:
	var level = payload.get("level", "mini")
	var amount = payload.get("amount", 0)
	var winner_name = payload.get("winner_name", "")
	print("[GameManager] JACKPOT WIN! Level=%s Amount=%d Winner=%s" % [level, amount, winner_name])
	emit_signal("jackpot_won", payload)

## Jackpot 觸發動畫通知（DAY-095）
func _handle_jackpot_animation(payload: Dictionary) -> void:
	var level = payload.get("level", "mini")
	var amount = payload.get("amount", 0)
	print("[GameManager] JACKPOT ANIMATION! Level=%s Amount=%d" % [level, amount])
	emit_signal("jackpot_animation", payload)

## 玩家統計更新（DAY-096）
func _handle_player_stats_update(payload: Dictionary) -> void:
	emit_signal("player_stats_updated", payload)

## 全服公告（DAY-097）
func _handle_announcement(payload: Dictionary) -> void:
	emit_signal("announcement_received", payload)

## 觀戰者加入通知（DAY-054d）
func _handle_spectator_join(payload: Dictionary) -> void:
	var count = payload.get("spectator_count", 1)
	print("[GameManager] Spectator joined! Total spectators: %d" % count)
	emit_signal("spectator_joined", payload)

## 每日登入獎勵（DAY-065）
func _handle_daily_bonus(payload: Dictionary) -> void:
	var streak = payload.get("streak", 1)
	var reward = payload.get("reward", 0)
	var is_new = payload.get("is_new_streak", false)
	if is_new:
		print("[GameManager] Daily bonus! Streak=%d Reward=%d" % [streak, reward])
		emit_signal("daily_bonus_received", payload)

## 觀戰者離開通知（DAY-055）
func _handle_spectator_leave(payload: Dictionary) -> void:
	var count = payload.get("spectator_count", 0)
	print("[GameManager] Spectator left! Remaining spectators: %d" % count)
	emit_signal("spectator_left", payload)

## 週賽排名更新（DAY-066）
func _handle_tournament_update(payload: Dictionary) -> void:
	var rank = payload.get("player_rank", 0)
	var points = payload.get("player_points", 0)
	if rank > 0:
		print("[GameManager] Tournament rank=%d points=%d" % [rank, points])
	emit_signal("tournament_updated", payload)

## 每日賽排名更新（DAY-093）
func _handle_daily_tournament_update(payload: Dictionary) -> void:
	var rank = payload.get("player_rank", 0)
	var points = payload.get("player_points", 0)
	if rank > 0:
		print("[GameManager] Daily Tournament rank=%d points=%d" % [rank, points])
	emit_signal("daily_tournament_updated", payload)

## 商店狀態更新（DAY-094）
func _handle_shop_update(payload: Dictionary) -> void:
	emit_signal("shop_updated", payload)

## 購買成功通知（DAY-094）
func _handle_shop_purchased(payload: Dictionary) -> void:
	var item_name = payload.get("item_name", "")
	print("[GameManager] Shop purchased: %s" % item_name)
	emit_signal("shop_purchased", payload)

## 購買失敗通知（DAY-094）
func _handle_shop_error(payload: Dictionary) -> void:
	var reason = payload.get("reason", "")
	print("[GameManager] Shop error: %s" % reason)
	emit_signal("shop_error", payload)

## 處理稱號解鎖通知（DAY-068）
func _handle_title_unlocked(payload: Dictionary) -> void:
	emit_signal("title_unlocked", payload)
	# 播放稱號解鎖音效（用 big_win 音效）
	if AudioManager != null:
		AudioManager.play_sfx("big_win")

## 處理砲台外觀更新（DAY-071）
func _handle_skin_update(payload: Dictionary) -> void:
	emit_signal("skin_updated", payload)

## 處理賽季通行證更新（DAY-072）
func _handle_season_update(payload: Dictionary) -> void:
	emit_signal("season_updated", payload)

## 處理賽季等級升級（DAY-072）
func _handle_season_level_up(payload: Dictionary) -> void:
	emit_signal("season_level_up", payload)

## 處理好友列表更新（DAY-073）
func _handle_friend_list(payload: Dictionary) -> void:
	emit_signal("friend_list_updated", payload)

## 處理好友請求通知（DAY-073）
func _handle_friend_request(payload: Dictionary) -> void:
	emit_signal("friend_request_received", payload)

## 處理好友狀態更新（DAY-073）
func _handle_friend_update(payload: Dictionary) -> void:
	emit_signal("friend_updated", payload)

## 處理禮物收到通知（DAY-101）
func _handle_gift_received(payload: Dictionary) -> void:
	emit_signal("gift_received", payload)

## 處理禮物送出成功通知（DAY-101）
func _handle_gift_sent(payload: Dictionary) -> void:
	emit_signal("gift_sent", payload)

## 處理禮物狀態（DAY-101）
func _handle_gift_status(payload: Dictionary) -> void:
	emit_signal("gift_status", payload)

## 處理禮物錯誤（DAY-101）
func _handle_gift_error(payload: Dictionary) -> void:
	emit_signal("gift_error", payload)

## 處理挑戰邀請（DAY-102）
func _handle_challenge_request(payload: Dictionary) -> void:
	emit_signal("challenge_request", payload)

## 處理挑戰狀態更新（DAY-102）
func _handle_challenge_updated(payload: Dictionary) -> void:
	emit_signal("challenge_updated", payload)

## 處理挑戰結果（DAY-102）
func _handle_challenge_result(payload: Dictionary) -> void:
	emit_signal("challenge_result", payload)

## 處理挑戰錯誤（DAY-102）
func _handle_challenge_error(payload: Dictionary) -> void:
	emit_signal("challenge_error", payload)

## 處理收到私訊（DAY-103）
func _handle_dm_received(payload: Dictionary) -> void:
	emit_signal("dm_received", payload)

## 處理私訊發送確認（DAY-103）
func _handle_dm_sent(payload: Dictionary) -> void:
	emit_signal("dm_sent", payload)

## 處理私訊發送失敗（DAY-103）
func _handle_dm_error(payload: Dictionary) -> void:
	emit_signal("dm_error", payload)

## 處理公會資訊更新（DAY-074）
func _handle_guild_update(payload: Dictionary) -> void:
	emit_signal("guild_updated", payload)

## 處理公會列表（DAY-074）
func _handle_guild_list(payload: Dictionary) -> void:
	# 如果沒有公會，自動加入第一個公會（簡化版）
	var guilds: Array = payload.get("guilds", [])
	if guilds.size() > 0:
		var first_guild: Dictionary = guilds[0]
		var guild_id: String = first_guild.get("guild_id", "")
		if guild_id != "":
			send_message("join_guild", {"guild_id": guild_id})

## 處理公會任務完成（DAY-074）
func _handle_guild_task_complete(payload: Dictionary) -> void:
	emit_signal("guild_task_complete", payload)

## 處理公會錯誤（DAY-074）
func _handle_guild_error(payload: Dictionary) -> void:
	var msg: String = payload.get("message", "公會操作失敗")
	print("[GameManager] Guild error: ", msg)

## 處理公會聊天訊息（DAY-075）
func _handle_guild_message(payload: Dictionary) -> void:
	emit_signal("guild_message_received", payload)

## 處理公會戰排名更新（DAY-076）
func _handle_guild_war_update(payload: Dictionary) -> void:
	emit_signal("guild_war_updated", payload)

## 處理公會戰結算通知（DAY-076）
func _handle_guild_war_result(payload: Dictionary) -> void:
	emit_signal("guild_war_result", payload)

## 請求公會戰狀態（DAY-076）
func request_guild_war_status() -> void:
	NetworkManager.send_message("get_guild_war_status", {})

## 處理每日 BOSS 狀態更新（DAY-077）
func _handle_daily_boss_update(payload: Dictionary) -> void:
	emit_signal("daily_boss_updated", payload)

## 處理每日 BOSS 擊殺通知（DAY-077）
func _handle_daily_boss_defeated(payload: Dictionary) -> void:
	emit_signal("daily_boss_defeated", payload)

## 請求每日 BOSS 狀態（DAY-077）
func request_daily_boss_status() -> void:
	NetworkManager.send_message("get_daily_boss", {})

## 取得顯示名稱
func get_display_name() -> String:
	return player_data.get("display_name", "玩家")

## 取得目前角色顏色
func get_character_color() -> Color:
	var char_id = player_data.get("character_id", "chiikawa")
	return CHARACTER_COLORS.get(char_id, Color.WHITE)

## 取得目前角色名稱
func get_character_name() -> String:
	var char_id = player_data.get("character_id", "chiikawa")
	return CHARACTER_NAMES.get(char_id, "吉伊卡哇")

## 取得目前金幣
func get_coins() -> int:
	return player_data.get("coins", 0)

## 取得勞動值
func get_labor_value() -> int:
	return player_data.get("labor_value", 0)

## 取得投注等級
func get_bet_level() -> int:
	return player_data.get("bet_level", 1)

## 取得投注消耗
func get_bet_cost() -> int:
	return player_data.get("bet_cost", 1)

## 是否自動攻擊
func is_auto() -> bool:
	return player_data.get("is_auto", false)

## 取得鎖定目標 ID
func get_lock_target_id() -> String:
	return player_data.get("lock_target_id", "")

## 取得玩家 ID（用於排行榜標記自己）
func get_player_id() -> String:
	return player_data.get("id", "")

## 處理 VIP 狀態更新（DAY-078）
func _handle_vip_update(payload: Dictionary) -> void:
	emit_signal("vip_updated", payload)

## 處理 VIP 升級通知（DAY-078）
func _handle_vip_level_up(payload: Dictionary) -> void:
	emit_signal("vip_level_up", payload)

## 處理 VIP 週獎勵領取通知（DAY-078）
func _handle_vip_weekly_claimed(payload: Dictionary) -> void:
	emit_signal("vip_weekly_claimed", payload)

## 請求 VIP 狀態（DAY-078）
func request_vip_status() -> void:
	NetworkManager.send_message("get_vip_status", {})

## 領取 VIP 週獎勵（DAY-078）
func claim_vip_weekly() -> void:
	NetworkManager.send_message("claim_vip_weekly", {})

## 處理限時活動狀態更新（DAY-079）
func _handle_event_update(payload: Dictionary) -> void:
	emit_signal("event_updated", payload)

## 請求限時活動狀態（DAY-079）
func request_event_status() -> void:
	NetworkManager.send_message("get_event_status", {})

## 處理圖鑑狀態更新（DAY-081）
func _handle_codex_update(payload: Dictionary) -> void:
	emit_signal("codex_updated", payload)

## 處理圖鑑條目解鎖通知（DAY-081）
func _handle_codex_unlock(payload: Dictionary) -> void:
	emit_signal("codex_unlocked", payload)

## 處理全圖鑑完成通知（DAY-081）
func _handle_codex_complete(payload: Dictionary) -> void:
	emit_signal("codex_complete", payload)

## 請求圖鑑狀態（DAY-081）
func request_codex() -> void:
	NetworkManager.send_message("get_codex", {})

## 處理推薦碼資訊（DAY-082）
func _handle_referral_info(payload: Dictionary) -> void:
	emit_signal("referral_info_received", payload)

## 處理推薦碼使用成功（DAY-082）
func _handle_referral_success(payload: Dictionary) -> void:
	emit_signal("referral_success", payload)

## 處理推薦碼使用失敗（DAY-082）
func _handle_referral_error(payload: Dictionary) -> void:
	emit_signal("referral_error", payload)

## 請求推薦碼資訊（DAY-082）
func request_referral_info() -> void:
	NetworkManager.send_message("get_referral_info", {})

## 使用推薦碼（DAY-082）
func use_referral_code(code: String) -> void:
	NetworkManager.send_message("use_referral_code", {"code": code})

## 處理連擊狀態更新（DAY-083）
func _handle_streak_update(payload: Dictionary) -> void:
	emit_signal("streak_updated", payload)

## 處理連擊重置通知（DAY-083）
func _handle_streak_reset(payload: Dictionary) -> void:
	emit_signal("streak_reset", payload)

## 處理幸運轉盤觸發（DAY-084）
func _handle_wheel_trigger(payload: Dictionary) -> void:
	emit_signal("wheel_triggered", payload)

## 處理隱藏挑戰解鎖（DAY-085）
func _handle_challenge_unlocked(payload: Dictionary) -> void:
	emit_signal("challenge_unlocked", payload)

## 處理任務連續完成獎勵（DAY-086）
func _handle_mission_streak_bonus(payload: Dictionary) -> void:
	emit_signal("mission_streak_bonus", payload)

## 處理天氣狀態更新（DAY-087）
func _handle_weather_update(payload: Dictionary) -> void:
	emit_signal("weather_updated", payload)

## 處理連鎖爆炸（DAY-088）
func _handle_chain_explosion(payload: Dictionary) -> void:
	emit_signal("chain_explosion", payload)

## 處理特殊武器狀態更新（DAY-089）
func _handle_special_weapon_update(payload: Dictionary) -> void:
	emit_signal("special_weapon_updated", payload)

## 處理特殊武器發射廣播（DAY-089）
func _handle_special_weapon_fired(payload: Dictionary) -> void:
	emit_signal("special_weapon_fired", payload)

## 處理神秘寶箱掉落通知（DAY-090）
func _handle_mystery_box_drop(payload: Dictionary) -> void:
	emit_signal("mystery_box_dropped", payload)

## 處理神秘寶箱狀態更新（DAY-090）
func _handle_mystery_box_update(payload: Dictionary) -> void:
	emit_signal("mystery_box_updated", payload)

## 處理神秘寶箱開箱結果（DAY-090）
func _handle_mystery_box_opened(payload: Dictionary) -> void:
	emit_signal("mystery_box_opened", payload)

## 處理房間列表（DAY-091）
func _handle_room_list(payload: Dictionary) -> void:
	emit_signal("room_list_received", payload)

## 處理房間切換成功（DAY-091）
func _handle_room_switched(payload: Dictionary) -> void:
	emit_signal("room_switched", payload)

## 處理房間操作失敗（DAY-091）
func _handle_room_error(payload: Dictionary) -> void:
	emit_signal("room_error", payload)

## 處理每日轉盤狀態（DAY-092）
func _handle_daily_spin_state(payload: Dictionary) -> void:
	emit_signal("daily_spin_state", payload)

## 處理每日轉盤結果（DAY-092）
func _handle_daily_spin_result(payload: Dictionary) -> void:
	emit_signal("daily_spin_result", payload)

## 處理登入里程碑達成通知（DAY-107）
func _handle_login_milestone(payload: Dictionary) -> void:	var days: int = payload.get("days", 0)
	var name: String = payload.get("name", "")
	print("[GameManager] Login milestone reached! Day=%d Name=%s" % [days, name])
	emit_signal("login_milestone_reached", payload)
	# 播放大獎音效
	if AudioManager != null:
		AudioManager.play_sfx("big_win")

## 處理登入進度回應（DAY-107）
func _handle_login_progress(payload: Dictionary) -> void:
	var streak: int = payload.get("current_streak", 0)
	print("[GameManager] Login progress: streak=%d" % streak)
	emit_signal("login_progress_received", payload)

## 請求登入進度（DAY-107）
func request_login_progress() -> void:
	NetworkManager.send_message("get_login_progress", {})

## 發送訊息（通用）
func send_message(type: String, payload: Dictionary) -> void:
	NetworkManager.send_message(type, payload)

## 處理超級 Bonus 通知（DAY-108）
func _handle_super_bonus_ready(payload: Dictionary) -> void:
	var label: String = payload.get("label", "SUPER BONUS!")
	var mult: float = payload.get("mult_bonus", 1.5)
	var combo: int = payload.get("combo_count", 3)
	print("[GameManager] %s combo=%d mult=%.1fx" % [label, combo, mult])
	emit_signal("super_bonus_triggered", payload)
	# 播放大獎音效
	if AudioManager != null:
		AudioManager.play_sfx("big_win")

## 發送訊息（通用）
func send_message(type: String, payload: Dictionary) -> void:
	NetworkManager.send_message(type, payload)

# ---- 賽季節日活動系統（DAY-109）----
signal festival_updated(festival_data: Dictionary)      # 節日狀態更新
signal festival_task_ready_signal(task_id: String)      # 節日任務可領取
signal festival_task_claimed_signal(task_data: Dictionary) # 節日任務獎勵領取成功
signal festival_title_earned_signal(title_data: Dictionary) # 節日稱號獲得
signal festival_error_signal(error_data: Dictionary)    # 節日操作失敗

## 處理節日狀態更新（DAY-109）
func _handle_festival_update(payload: Dictionary) -> void:
	var festival_type: String = payload.get("type", "none")
	var is_active: bool = payload.get("is_active", false)
	if is_active:
		print("[GameManager] Festival active: %s" % festival_type)
	emit_signal("festival_updated", payload)

## 處理節日任務可領取通知（DAY-109）
func _handle_festival_task_ready(payload: Dictionary) -> void:
	var task_id: String = payload.get("task_id", "")
	print("[GameManager] Festival task ready: %s" % task_id)
	emit_signal("festival_task_ready_signal", task_id)
	# 播放任務完成音效
	if AudioManager != null:
		AudioManager.play_sfx("bonus_ready")

## 處理節日任務獎勵領取成功（DAY-109）
func _handle_festival_task_claimed(payload: Dictionary) -> void:
	var task_id: String = payload.get("task_id", "")
	var coins: int = payload.get("reward_coins", 0)
	print("[GameManager] Festival task claimed: %s reward=%d" % [task_id, coins])
	emit_signal("festival_task_claimed_signal", payload)

## 處理節日稱號獲得通知（DAY-109）
func _handle_festival_title_earned(payload: Dictionary) -> void:
	var title_name: String = payload.get("title_name", "")
	print("[GameManager] Festival title earned: %s" % title_name)
	emit_signal("festival_title_earned_signal", payload)
	# 播放大獎音效
	if AudioManager != null:
		AudioManager.play_sfx("big_win")

## 處理節日操作失敗（DAY-109）
func _handle_festival_error(payload: Dictionary) -> void:
	emit_signal("festival_error_signal", payload)

## 請求節日狀態（DAY-109）
func request_festival() -> void:
	NetworkManager.send_message("get_festival", {})

## 發送領取節日任務獎勵請求（DAY-109）
func send_claim_festival_task(task_id: String) -> void:
	NetworkManager.send_message("claim_festival_task", {"task_id": task_id})
