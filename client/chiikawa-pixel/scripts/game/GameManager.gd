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
# DAY-303 新增幸運特殊魚訊號
signal lucky_crash_fish(data: Dictionary)
# DAY-304 新增幸運特殊魚訊號
signal lucky_electric_eel(data: Dictionary)
signal lucky_angler_fish(data: Dictionary)
signal lucky_black_hole(data: Dictionary)
signal lucky_bounty_hunter(data: Dictionary)
signal lucky_tsunami(data: Dictionary)
# DAY-305 新增幸運特殊魚訊號
signal lucky_dragon_wrath_v2(data: Dictionary)
signal lucky_humpback_whale(data: Dictionary)
signal lucky_legend_dragon(data: Dictionary)
signal lucky_guild_war(data: Dictionary)
signal lucky_quality_fish(data: Dictionary)
# DAY-306 新增幸運特殊魚訊號
signal lucky_tornado(data: Dictionary)
signal lucky_earthquake(data: Dictionary)
signal lucky_volcano(data: Dictionary)
signal lucky_cosmic_ray(data: Dictionary)
signal lucky_divine_dragon(data: Dictionary)
# DAY-307 新增幸運特殊魚訊號
signal lucky_quantum(data: Dictionary)
signal lucky_supernova(data: Dictionary)
signal lucky_infinite(data: Dictionary)
signal lucky_genesis(data: Dictionary)
signal lucky_rebirth(data: Dictionary)
# DAY-308 新增幸運特殊魚訊號
signal lucky_awakened_croc(data: Dictionary)
signal lucky_vampire_v2(data: Dictionary)
signal lucky_super_awaken(data: Dictionary)
signal lucky_giant_prize(data: Dictionary)
signal lucky_immortal_boss(data: Dictionary)
# DAY-309 新增幸運特殊魚訊號
signal lucky_ice_phoenix(data: Dictionary)
signal lucky_dragon_fury(data: Dictionary)
signal lucky_mult_cascade(data: Dictionary)
signal lucky_awaken_boss_v2(data: Dictionary)
signal lucky_ultimate_judgment(data: Dictionary)
# DAY-310 新增幸運特殊魚訊號
signal lucky_combo_burst(data: Dictionary)
signal lucky_time_bomb(data: Dictionary)
signal lucky_elemental_fusion(data: Dictionary)
signal lucky_treasure_hunter(data: Dictionary)
signal lucky_myth_awaken(data: Dictionary)
# DAY-312 新增幸運特殊魚訊號
signal lucky_star_portal(data: Dictionary)
signal lucky_dragon_soul(data: Dictionary)
signal lucky_spacetime_rift(data: Dictionary)
signal lucky_holy_judgment(data: Dictionary)
signal lucky_big_bang(data: Dictionary)
# DAY-313 新增幸運特殊魚訊號（Progressive Jackpot 系列）
signal lucky_jackpot_pool(data: Dictionary)
# DAY-314 新增幸運特殊魚訊號
signal lucky_multiverse(data: Dictionary)
signal lucky_time_loop(data: Dictionary)
signal lucky_fate_wheel(data: Dictionary)
signal lucky_divine_realm(data: Dictionary)
signal lucky_final_power(data: Dictionary)
# DAY-315 新增幸運特殊魚訊號
signal lucky_mutation(data: Dictionary)
signal lucky_arctic_storm(data: Dictionary)
signal lucky_fisher_wild(data: Dictionary)
signal lucky_risk_level(data: Dictionary)
signal lucky_cosmic_pulse(data: Dictionary)
# DAY-316 新增幸運特殊魚訊號
signal lucky_mirror_universe(data: Dictionary)
signal lucky_gravity_field(data: Dictionary)
signal lucky_time_acceleration(data: Dictionary)
signal lucky_nebula_vortex(data: Dictionary)
signal lucky_cosmic_judgment(data: Dictionary)
# DAY-317 新增幸運特殊魚訊號
signal lucky_pvp_battle(data: Dictionary)
signal lucky_skill_chain(data: Dictionary)
signal lucky_global_explosion(data: Dictionary)
signal lucky_spacetime_fold(data: Dictionary)
signal lucky_cosmic_end(data: Dictionary)
# DAY-318 新增幸運特殊魚訊號
signal lucky_dragon_king(data: Dictionary)
signal lucky_eternal_cycle(data: Dictionary)
signal lucky_chaos_explosion(data: Dictionary)
signal lucky_divine_revival(data: Dictionary)
signal lucky_genesis_epoch(data: Dictionary)
# DAY-319 新增幸運特殊魚訊號
signal lucky_energy_storm(data: Dictionary)
signal lucky_crystal_resonance(data: Dictionary)
signal lucky_fate_judgment(data: Dictionary)
signal lucky_time_reversal(data: Dictionary)
signal lucky_cosmic_singularity(data: Dictionary)
# DAY-323 新增幸運特殊魚訊號
signal lucky_fever_boost(data: Dictionary)
signal lucky_guild_battle(data: Dictionary)
signal lucky_path_fish(data: Dictionary)
signal lucky_chain_eel(data: Dictionary)
signal lucky_ultimate_miracle(data: Dictionary)
# DAY-324 新增
signal lucky_avalanche(data: Dictionary)
signal lucky_crash_multiplier(data: Dictionary)
signal lucky_multiplier_ladder(data: Dictionary)
signal lucky_ice_fishing_wheel(data: Dictionary)
signal lucky_global_avalanche(data: Dictionary)

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
		# DAY-303 新增幸運特殊魚事件
		"lucky_crash_fish":
			emit_signal("lucky_crash_fish", payload)
		# DAY-304 新增幸運特殊魚事件
		"lucky_electric_eel":
			emit_signal("lucky_electric_eel", payload)
		"lucky_angler_fish":
			emit_signal("lucky_angler_fish", payload)
		"lucky_black_hole":
			emit_signal("lucky_black_hole", payload)
		"lucky_bounty_hunter":
			emit_signal("lucky_bounty_hunter", payload)
		"lucky_tsunami":
			emit_signal("lucky_tsunami", payload)
		# DAY-305 新增幸運特殊魚事件
		"lucky_dragon_wrath_v2":
			emit_signal("lucky_dragon_wrath_v2", payload)
		"lucky_humpback_whale":
			emit_signal("lucky_humpback_whale", payload)
		"lucky_legend_dragon":
			emit_signal("lucky_legend_dragon", payload)
		"lucky_guild_war":
			emit_signal("lucky_guild_war", payload)
		"lucky_quality_fish":
			emit_signal("lucky_quality_fish", payload)
		# DAY-306 新增幸運特殊魚事件
		"lucky_tornado":
			emit_signal("lucky_tornado", payload)
		"lucky_earthquake":
			emit_signal("lucky_earthquake", payload)
		"lucky_volcano":
			emit_signal("lucky_volcano", payload)
		"lucky_cosmic_ray":
			emit_signal("lucky_cosmic_ray", payload)
		"lucky_divine_dragon":
			emit_signal("lucky_divine_dragon", payload)
		# DAY-307 新增幸運特殊魚事件
		"lucky_quantum":
			emit_signal("lucky_quantum", payload)
		"lucky_supernova":
			emit_signal("lucky_supernova", payload)
		"lucky_infinite":
			emit_signal("lucky_infinite", payload)
		"lucky_genesis":
			emit_signal("lucky_genesis", payload)
		"lucky_rebirth":
			emit_signal("lucky_rebirth", payload)
		# DAY-308 新增幸運特殊魚事件
		"lucky_awakened_croc":
			emit_signal("lucky_awakened_croc", payload)
		"lucky_vampire_v2":
			emit_signal("lucky_vampire_v2", payload)
		"lucky_super_awaken":
			emit_signal("lucky_super_awaken", payload)
		"lucky_giant_prize":
			emit_signal("lucky_giant_prize", payload)
		"lucky_immortal_boss":
			emit_signal("lucky_immortal_boss", payload)
		# DAY-309 新增幸運特殊魚事件
		"lucky_ice_phoenix":
			emit_signal("lucky_ice_phoenix", payload)
		"lucky_dragon_fury":
			emit_signal("lucky_dragon_fury", payload)
		"lucky_mult_cascade":
			emit_signal("lucky_mult_cascade", payload)
		"lucky_awaken_boss_v2":
			emit_signal("lucky_awaken_boss_v2", payload)
		"lucky_ultimate_judgment":
			emit_signal("lucky_ultimate_judgment", payload)
		# DAY-310 新增
		"lucky_combo_burst":
			emit_signal("lucky_combo_burst", payload)
		"lucky_time_bomb":
			emit_signal("lucky_time_bomb", payload)
		"lucky_elemental_fusion":
			emit_signal("lucky_elemental_fusion", payload)
		"lucky_treasure_hunter":
			emit_signal("lucky_treasure_hunter", payload)
		"lucky_myth_awaken":
			emit_signal("lucky_myth_awaken", payload)
		# DAY-312 新增
		"lucky_star_portal":
			emit_signal("lucky_star_portal", payload)
		"lucky_dragon_soul":
			emit_signal("lucky_dragon_soul", payload)
		"lucky_spacetime_rift":
			emit_signal("lucky_spacetime_rift", payload)
		"lucky_holy_judgment":
			emit_signal("lucky_holy_judgment", payload)
		"lucky_big_bang":
			emit_signal("lucky_big_bang", payload)
		# DAY-313 新增
		"lucky_jackpot_pool":
			emit_signal("lucky_jackpot_pool", payload)
		# DAY-314 新增
		"lucky_multiverse":
			emit_signal("lucky_multiverse", payload)
		"lucky_time_loop":
			emit_signal("lucky_time_loop", payload)
		"lucky_fate_wheel":
			emit_signal("lucky_fate_wheel", payload)
		"lucky_divine_realm":
			emit_signal("lucky_divine_realm", payload)
		"lucky_final_power":
			emit_signal("lucky_final_power", payload)
		# DAY-315 新增
		"lucky_mutation":
			emit_signal("lucky_mutation", payload)
		"lucky_arctic_storm":
			emit_signal("lucky_arctic_storm", payload)
		"lucky_fisher_wild":
			emit_signal("lucky_fisher_wild", payload)
		"lucky_risk_level":
			emit_signal("lucky_risk_level", payload)
		"lucky_cosmic_pulse":
			emit_signal("lucky_cosmic_pulse", payload)
		# DAY-316 新增
		"lucky_mirror_universe":
			emit_signal("lucky_mirror_universe", payload)
		"lucky_gravity_field":
			emit_signal("lucky_gravity_field", payload)
		"lucky_time_acceleration":
			emit_signal("lucky_time_acceleration", payload)
		"lucky_nebula_vortex":
			emit_signal("lucky_nebula_vortex", payload)
		"lucky_cosmic_judgment":
			emit_signal("lucky_cosmic_judgment", payload)
		# DAY-317 新增
		"lucky_pvp_battle":
			emit_signal("lucky_pvp_battle", payload)
		"lucky_skill_chain":
			emit_signal("lucky_skill_chain", payload)
		"lucky_global_explosion":
			emit_signal("lucky_global_explosion", payload)
		"lucky_spacetime_fold":
			emit_signal("lucky_spacetime_fold", payload)
		"lucky_cosmic_end":
			emit_signal("lucky_cosmic_end", payload)
		# DAY-318 新增
		"lucky_dragon_king":
			emit_signal("lucky_dragon_king", payload)
		"lucky_eternal_cycle":
			emit_signal("lucky_eternal_cycle", payload)
		"lucky_chaos_explosion":
			emit_signal("lucky_chaos_explosion", payload)
		"lucky_divine_revival":
			emit_signal("lucky_divine_revival", payload)
		"lucky_genesis_epoch":
			emit_signal("lucky_genesis_epoch", payload)
		# DAY-319 新增
		"lucky_energy_storm":
			emit_signal("lucky_energy_storm", payload)
		"lucky_crystal_resonance":
			emit_signal("lucky_crystal_resonance", payload)
		"lucky_fate_judgment":
			emit_signal("lucky_fate_judgment", payload)
		"lucky_time_reversal":
			emit_signal("lucky_time_reversal", payload)
		"lucky_cosmic_singularity":
			emit_signal("lucky_cosmic_singularity", payload)
		# DAY-323 新增
		"lucky_fever_boost":
			emit_signal("lucky_fever_boost", payload)
		"lucky_guild_battle":
			emit_signal("lucky_guild_battle", payload)
		"lucky_path_fish":
			emit_signal("lucky_path_fish", payload)
		"lucky_chain_eel":
			emit_signal("lucky_chain_eel", payload)
		"lucky_ultimate_miracle":
			emit_signal("lucky_ultimate_miracle", payload)
		# DAY-324 新增
		"lucky_avalanche":
			emit_signal("lucky_avalanche", payload)
		"lucky_crash_multiplier":
			emit_signal("lucky_crash_multiplier", payload)
		"lucky_multiplier_ladder":
			emit_signal("lucky_multiplier_ladder", payload)
		"lucky_ice_fishing_wheel":
			emit_signal("lucky_ice_fishing_wheel", payload)
		"lucky_global_avalanche":
			emit_signal("lucky_global_avalanche", payload)
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

func get_player_id() -> String:
	return NetworkManager.get_player_id()
