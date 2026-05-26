## GameManager.gd — 遊戲狀態機與訊號分發
## game-state-agent 負責維護
extends Node

# ── 訊號 ─────────────────────────────────────────────────────
signal player_updated(data: Dictionary)
signal game_state_changed(new_state: String)
signal reward_received(reward: Dictionary)
signal attack_result(result: Dictionary)
signal target_spawned(data: Dictionary)
signal target_updated(data: Dictionary)
signal target_killed(data: Dictionary)
signal boss_event(event_data: Dictionary)
signal bonus_event(event_data: Dictionary)
signal announce(data: Dictionary)
# DAY-292 幸運特殊魚訊號
signal lucky_chain_lightning(data: Dictionary)
signal lucky_crab_torpedo(data: Dictionary)
signal lucky_vortex(data: Dictionary)
signal lucky_golden_dragon(data: Dictionary)
signal lucky_thunder_lobster(data: Dictionary)
# DAY-293 新增幸運特殊魚訊號
signal lucky_awakened_phoenix(data: Dictionary)
signal lucky_shockwave_bomb(data: Dictionary)
# DAY-294 新增幸運特殊魚訊號
signal lucky_drill_torpedo(data: Dictionary)
signal lucky_time_freeze(data: Dictionary)
signal lucky_chain_explosion(data: Dictionary)
# DAY-295 新增幸運特殊魚訊號
signal lucky_chain_long_king(data: Dictionary)
signal lucky_dragon_shotgun(data: Dictionary)
signal lucky_rocket_cannon(data: Dictionary)
signal lucky_deep_whirlpool(data: Dictionary)
signal lucky_vampire_mult(data: Dictionary)
# DAY-296 新增幸運特殊魚訊號
signal lucky_mirror_fish(data: Dictionary)
signal lucky_golden_rain(data: Dictionary)
signal lucky_freeze_bomb(data: Dictionary)
signal lucky_thunder_storm(data: Dictionary)
signal lucky_lucky_wheel(data: Dictionary)
# DAY-301 新增幸運特殊魚訊號
signal lucky_jackpot_fish(data: Dictionary)
signal lucky_coop_fish(data: Dictionary)
signal lucky_time_warp(data: Dictionary)
# DAY-302 新增幸運特殊魚訊號
signal lucky_chain_meteor(data: Dictionary)

# ── 玩家資料快取 ──────────────────────────────────────────────
var player_data: Dictionary = {}
var current_state: String = "normal_play"

func _ready() -> void:
	NetworkManager.message_received.connect(_on_message)
	NetworkManager.connected.connect(_on_connected)
	NetworkManager.disconnected.connect(_on_disconnected)

func _on_connected() -> void:
	print("[GameManager] Connected to server")

func _on_disconnected() -> void:
	print("[GameManager] Disconnected from server")

func _on_message(type: String, payload: Dictionary) -> void:
	match type:
		"game_state":
			current_state = payload.get("state", "normal_play")
			emit_signal("game_state_changed", current_state)
		"player_update":
			player_data = payload
			emit_signal("player_updated", payload)
		"target_spawn":
			emit_signal("target_spawned", payload)
		"target_update":
			emit_signal("target_updated", payload)
		"target_kill":
			emit_signal("target_killed", payload)
		"attack_result":
			emit_signal("attack_result", payload)
		"reward":
			emit_signal("reward_received", payload)
		"boss_event":
			emit_signal("boss_event", payload)
		"bonus_event":
			emit_signal("bonus_event", payload)
		"announce":
			emit_signal("announce", payload)
		# DAY-292 幸運特殊魚事件
		"lucky_chain_lightning":
			emit_signal("lucky_chain_lightning", payload)
		"lucky_crab_torpedo":
			emit_signal("lucky_crab_torpedo", payload)
		"lucky_vortex":
			emit_signal("lucky_vortex", payload)
		"lucky_golden_dragon":
			emit_signal("lucky_golden_dragon", payload)
		"lucky_thunder_lobster":
			emit_signal("lucky_thunder_lobster", payload)
		# DAY-293 新增幸運特殊魚事件
		"lucky_awakened_phoenix":
			emit_signal("lucky_awakened_phoenix", payload)
		"lucky_shockwave_bomb":
			emit_signal("lucky_shockwave_bomb", payload)
		# DAY-294 新增幸運特殊魚事件
		"lucky_drill_torpedo":
			emit_signal("lucky_drill_torpedo", payload)
		"lucky_time_freeze":
			emit_signal("lucky_time_freeze", payload)
		"lucky_chain_explosion":
			emit_signal("lucky_chain_explosion", payload)
		# DAY-295 新增幸運特殊魚事件
		"lucky_chain_long_king":
			emit_signal("lucky_chain_long_king", payload)
		"lucky_dragon_shotgun":
			emit_signal("lucky_dragon_shotgun", payload)
		"lucky_rocket_cannon":
			emit_signal("lucky_rocket_cannon", payload)
		"lucky_deep_whirlpool":
			emit_signal("lucky_deep_whirlpool", payload)
		"lucky_vampire_mult":
			emit_signal("lucky_vampire_mult", payload)
		# DAY-296 新增幸運特殊魚事件
		"lucky_mirror_fish":
			emit_signal("lucky_mirror_fish", payload)
		"lucky_golden_rain":
			emit_signal("lucky_golden_rain", payload)
		"lucky_freeze_bomb":
			emit_signal("lucky_freeze_bomb", payload)
		"lucky_thunder_storm":
			emit_signal("lucky_thunder_storm", payload)
		"lucky_lucky_wheel":
			emit_signal("lucky_lucky_wheel", payload)
		# DAY-301 新增幸運特殊魚事件
		"lucky_jackpot_fish":
			emit_signal("lucky_jackpot_fish", payload)
		"lucky_coop_fish":
			emit_signal("lucky_coop_fish", payload)
		"lucky_time_warp":
			emit_signal("lucky_time_warp", payload)
		# DAY-302 新增幸運特殊魚事件
		"lucky_chain_meteor":
			emit_signal("lucky_chain_meteor", payload)
		"pong":
			pass
		"error":
			push_error("[GameManager] Server error: " + str(payload))

# ── 玩家資料存取 ──────────────────────────────────────────────
func get_coins() -> int:
	return player_data.get("coins", 0)

func get_bet_level() -> int:
	return player_data.get("bet_level", 1)

func get_bet_cost() -> int:
	return player_data.get("bet_cost", 1)

func get_character_id() -> String:
	return player_data.get("character_id", "chiikawa")

func get_character_name() -> String:
	return player_data.get("character_name", "Chiikawa")

func get_character_color() -> Color:
	match get_character_id():
		"hachiware": return Color(0.4, 0.6, 1.0)
		"usagi": return Color(1.0, 0.9, 0.2)
		_: return Color(1.0, 0.6, 0.8)

func get_labor_value() -> int:
	return player_data.get("labor_value", 0)

func is_auto() -> bool:
	return player_data.get("is_auto", false)

func get_lock_target_id() -> String:
	return player_data.get("lock_target_id", "")

func get_fire_rate() -> float:
	return player_data.get("fire_rate", 2.0)

func get_projectile_speed() -> float:
	return player_data.get("projectile_speed", 700.0)

func get_combo_count() -> int:
	return player_data.get("combo_count", 0)

func get_combo_mult_bonus() -> float:
	return player_data.get("combo_mult_bonus", 0.0)
