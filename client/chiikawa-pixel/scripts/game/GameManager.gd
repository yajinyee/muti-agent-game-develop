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
		"error":
			_handle_error(payload)
		"pong":
			pass  # Ping/Pong 心跳

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

func _handle_error(payload: Dictionary) -> void:
	push_warning("[GameManager] Server error: " + str(payload))

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
