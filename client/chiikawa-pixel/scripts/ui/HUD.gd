## HUD.gd — 核心 HUD
## hud-core-agent 負責維護
## DAY-298：Lucky 事件視覺系統改由 LuckyEventSystem 處理
extends CanvasLayer

@onready var coins_label: Label = $TopBar/CoinsLabel
@onready var bet_label: Label = $TopBar/BetLabel
@onready var char_label: Label = $TopBar/CharLabel
@onready var labor_bar: ProgressBar = $TopBar/LaborBar
@onready var labor_label: Label = $TopBar/LaborLabel
@onready var state_label: Label = $TopBar/StateLabel
@onready var auto_btn: Button = $BottomBar/AutoBtn
@onready var lock_btn: Button = $BottomBar/LockBtn
@onready var bet_minus_btn: Button = $BottomBar/BetMinusBtn
@onready var bet_plus_btn: Button = $BottomBar/BetPlusBtn
@onready var boss_btn: Button = $BottomBar/BossBtn
@onready var bonus_btn: Button = $BottomBar/BonusBtn

var _reward_popup: Label = null
var _disconnect_overlay: Control = null
var _boss_timer_panel: Control = null
var _boss_time_left: float = 0.0
var _boss_active: bool = false
var _last_labor: int = 0
# DAY-298：Lucky 事件視覺系統（LuckyEventSystem 節點引用，在 Main.tscn 中設定）
var lucky_event_system: Node = null
# DAY-297 Combo UI
var _combo_label: Label = null
var _last_combo: int = 0

func _ready() -> void:
	GameManager.player_updated.connect(_on_player_updated)
	GameManager.game_state_changed.connect(_on_state_changed)
	GameManager.reward_received.connect(_on_reward_received)
	GameManager.boss_event.connect(_on_boss_event)
	GameManager.bonus_event.connect(_on_bonus_event)
	NetworkManager.connected.connect(_on_reconnected)
	NetworkManager.disconnected.connect(_on_disconnected)

	auto_btn.pressed.connect(func(): NetworkManager.send_auto_toggle())
	lock_btn.pressed.connect(func(): NetworkManager.send_lock(""))
	bet_minus_btn.pressed.connect(func(): NetworkManager.send_bet_change(max(1, GameManager.get_bet_level() - 1)))
	bet_plus_btn.pressed.connect(func(): NetworkManager.send_bet_change(min(10, GameManager.get_bet_level() + 1)))
	boss_btn.pressed.connect(func(): NetworkManager.send_trigger_boss())
	bonus_btn.pressed.connect(func(): NetworkManager.send_trigger_bonus())

	_create_reward_popup()
	_create_disconnect_overlay()
	_create_combo_label()
	_update_ui()
	# 嘗試自動找 LuckyEventSystem（如果在同一場景樹中）
	call_deferred("_find_lucky_event_system")

	# DAY-292 幸運特殊魚訊號連接
	GameManager.lucky_chain_lightning.connect(_on_lucky_chain_lightning)
	GameManager.lucky_crab_torpedo.connect(_on_lucky_crab_torpedo)
	GameManager.lucky_vortex.connect(_on_lucky_vortex)
	GameManager.lucky_golden_dragon.connect(_on_lucky_golden_dragon)
	GameManager.lucky_thunder_lobster.connect(_on_lucky_thunder_lobster)
	GameManager.announce.connect(_on_announce)
	# DAY-293 新增幸運特殊魚訊號連接
	GameManager.lucky_awakened_phoenix.connect(_on_lucky_awakened_phoenix)
	GameManager.lucky_shockwave_bomb.connect(_on_lucky_shockwave_bomb)
	# DAY-294 新增幸運特殊魚訊號連接
	GameManager.lucky_drill_torpedo.connect(_on_lucky_drill_torpedo)
	GameManager.lucky_time_freeze.connect(_on_lucky_time_freeze)
	GameManager.lucky_chain_explosion.connect(_on_lucky_chain_explosion)
	# DAY-295 新增幸運特殊魚訊號連接
	GameManager.lucky_chain_long_king.connect(_on_lucky_chain_long_king)
	GameManager.lucky_dragon_shotgun.connect(_on_lucky_dragon_shotgun)
	GameManager.lucky_rocket_cannon.connect(_on_lucky_rocket_cannon)
	GameManager.lucky_deep_whirlpool.connect(_on_lucky_deep_whirlpool)
	GameManager.lucky_vampire_mult.connect(_on_lucky_vampire_mult)
	# DAY-296 新增幸運特殊魚訊號連接
	GameManager.lucky_mirror_fish.connect(_on_lucky_mirror_fish)
	GameManager.lucky_golden_rain.connect(_on_lucky_golden_rain)
	GameManager.lucky_freeze_bomb.connect(_on_lucky_freeze_bomb)
	GameManager.lucky_thunder_storm.connect(_on_lucky_thunder_storm)
	GameManager.lucky_lucky_wheel.connect(_on_lucky_lucky_wheel)
	# DAY-301 新增幸運特殊魚訊號連接
	GameManager.lucky_jackpot_fish.connect(_on_lucky_jackpot_fish)
	GameManager.lucky_coop_fish.connect(_on_lucky_coop_fish)
	GameManager.lucky_time_warp.connect(_on_lucky_time_warp)
	# DAY-302 新增幸運特殊魚訊號連接
	GameManager.lucky_chain_meteor.connect(_on_lucky_chain_meteor)
	# DAY-303 新增幸運特殊魚訊號連接
	GameManager.lucky_crash_fish.connect(_on_lucky_crash_fish)

func _process(delta: float) -> void:
	if _boss_active and _boss_time_left > 0:
		_boss_time_left -= delta
		_update_boss_timer()

func _find_lucky_event_system() -> void:
	# 在場景樹中尋找 LuckyEventSystem 節點
	var root = get_tree().get_root()
	lucky_event_system = _find_node_by_class(root, "LuckyEventSystem")
	if lucky_event_system == null:
		# 備用：在同層 CanvasLayer 中找
		var parent = get_parent()
		if is_instance_valid(parent):
			lucky_event_system = parent.get_node_or_null("LuckyEventSystem")
	if lucky_event_system != null:
		print("[HUD] LuckyEventSystem found: ", lucky_event_system.get_path())
	else:
		print("[HUD] LuckyEventSystem not found, using fallback banner")

func _find_node_by_class(node: Node, class_name_str: String) -> Node:
	if node.get_script() != null:
		var path = node.get_script().resource_path
		if path.ends_with(class_name_str + ".gd"):
			return node
	for child in node.get_children():
		var result = _find_node_by_class(child, class_name_str)
		if result != null:
			return result
	return null

func _on_player_updated(_data: Dictionary) -> void:
	_update_ui()

func _update_ui() -> void:
	coins_label.text = "💰 %d" % GameManager.get_coins()
	var lv = GameManager.get_bet_level()
	var cost = GameManager.get_bet_cost()
	bet_label.text = "BET LV%d (%d)" % [lv, cost]
	char_label.text = GameManager.get_character_name()
	char_label.modulate = GameManager.get_character_color()

	var labor = GameManager.get_labor_value()
	labor_bar.value = labor
	if labor >= 80:
		labor_label.text = "⚡%d/100" % labor
		labor_label.modulate = Color(1.0, 0.9, 0.2)
	else:
		labor_label.text = "%d/100" % labor
		labor_label.modulate = Color.WHITE
	if labor >= 100 and _last_labor < 100:
		ScreenShake.add_trauma(0.3)
	_last_labor = labor

	if GameManager.is_auto():
		auto_btn.text = "AUTO ON"
		auto_btn.modulate = Color(0.3, 1.0, 0.3)
	else:
		auto_btn.text = "AUTO"
		auto_btn.modulate = Color.WHITE

	var lock_id = GameManager.get_lock_target_id()
	if lock_id != "":
		lock_btn.text = "🔒 LOCK"
		lock_btn.modulate = Color(1.0, 0.8, 0.2)
	else:
		lock_btn.text = "🔓 LOCK"
		lock_btn.modulate = Color(0.7, 0.7, 0.7)

	# Combo 顯示
	_update_combo_display()

func _on_state_changed(new_state: String) -> void:
	state_label.text = new_state.to_upper().replace("_", " ")
	match new_state:
		"boss_battle":
			AudioManager.play_bgm(AudioManager.BGM.BOSS_ENTER)
		"boss_result":
			AudioManager.play_bgm(AudioManager.BGM.MAIN_GAME)
		"bonus_game":
			AudioManager.play_bgm(AudioManager.BGM.BONUS_GAME)
		"bonus_result", "normal_play":
			AudioManager.play_bgm(AudioManager.BGM.MAIN_GAME)

func _on_reward_received(reward: Dictionary) -> void:
	var amount = reward.get("amount", 0)
	var mult = reward.get("multiplier", 1.0)
	if amount <= 0:
		return
	_show_reward_popup(amount, mult)

func _show_reward_popup(amount: int, mult: float) -> void:
	if not is_instance_valid(_reward_popup):
		return
	var icon = "💰"
	if mult >= 100:
		icon = "🌟"
		_reward_popup.modulate = Color(1.0, 0.3, 0.1)
	elif mult >= 20:
		icon = "⭐"
		_reward_popup.modulate = Color(1.0, 0.85, 0.0)
	else:
		_reward_popup.modulate = Color.WHITE
	_reward_popup.text = "%s +%d  x%.0f" % [icon, amount, mult]
	_reward_popup.visible = true
	_reward_popup.modulate.a = 1.0
	_reward_popup.position.y = 350
	var tween = create_tween()
	tween.tween_property(_reward_popup, "position:y", 280.0, 0.7)
	tween.parallel().tween_property(_reward_popup, "modulate:a", 0.0, 0.7)
	tween.tween_callback(func(): if is_instance_valid(_reward_popup): _reward_popup.visible = false)

func _on_boss_event(event_data: Dictionary) -> void:
	match event_data.get("event", ""):
		"warning":
			AudioManager.play_sfx(AudioManager.SFX.BOSS_WARNING)
		"spawn":
			_start_boss_timer()
		"phase_change":
			AudioManager.play_bgm(AudioManager.BGM.BOSS_RAGE)
		"kill", "timeout":
			_stop_boss_timer()
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)

func _on_bonus_event(event_data: Dictionary) -> void:
	match event_data.get("event", ""):
		"start":
			AudioManager.play_sfx(AudioManager.SFX.BONUS_READY)

# ── BOSS 計時器 ───────────────────────────────────────────────

func _start_boss_timer() -> void:
	_boss_time_left = 60.0
	_boss_active = true
	if is_instance_valid(_boss_timer_panel):
		_boss_timer_panel.queue_free()

	var panel = Control.new()
	panel.name = "BossTimerPanel"
	panel.position = Vector2(900, 50)
	panel.size = Vector2(340, 75)
	add_child(panel)
	_boss_timer_panel = panel

	var bg = ColorRect.new()
	bg.size = panel.size
	bg.color = Color(0.1, 0.0, 0.0, 0.85)
	panel.add_child(bg)

	var title = Label.new()
	title.text = "⚔ BOSS BATTLE"
	title.position = Vector2(10, 5)
	title.add_theme_font_size_override("font_size", 15)
	title.modulate = Color(1.0, 0.3, 0.3)
	panel.add_child(title)

	var timer_lbl = Label.new()
	timer_lbl.name = "TimerLabel"
	timer_lbl.text = "60.0s"
	timer_lbl.position = Vector2(10, 26)
	timer_lbl.add_theme_font_size_override("font_size", 26)
	timer_lbl.modulate = Color(1.0, 0.9, 0.2)
	panel.add_child(timer_lbl)

	var mult_lbl = Label.new()
	mult_lbl.name = "MultLabel"
	mult_lbl.text = "500x"
	mult_lbl.position = Vector2(190, 26)
	mult_lbl.add_theme_font_size_override("font_size", 26)
	mult_lbl.modulate = Color(1.0, 0.5, 0.0)
	panel.add_child(mult_lbl)

	var hint = Label.new()
	hint.text = "Kill faster = higher reward!"
	hint.position = Vector2(10, 56)
	hint.add_theme_font_size_override("font_size", 11)
	hint.modulate = Color(0.8, 0.8, 0.8)
	panel.add_child(hint)

func _update_boss_timer() -> void:
	if not is_instance_valid(_boss_timer_panel):
		return
	var tl = _boss_timer_panel.get_node_or_null("TimerLabel")
	var ml = _boss_timer_panel.get_node_or_null("MultLabel")
	if is_instance_valid(tl):
		tl.text = "%.1fs" % max(0, _boss_time_left)
		tl.modulate = Color(1.0, 0.3, 0.3) if _boss_time_left <= 10 else Color(1.0, 0.9, 0.2)
	if is_instance_valid(ml):
		var m = 100
		if _boss_time_left > 50: m = 500
		elif _boss_time_left > 40: m = 400
		elif _boss_time_left > 30: m = 300
		elif _boss_time_left > 20: m = 200
		elif _boss_time_left > 10: m = 150
		ml.text = "%dx" % m

func _stop_boss_timer() -> void:
	_boss_active = false
	if is_instance_valid(_boss_timer_panel):
		var tween = create_tween()
		tween.tween_property(_boss_timer_panel, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func():
			if is_instance_valid(_boss_timer_panel):
				_boss_timer_panel.queue_free()
				_boss_timer_panel = null
		)

# ── 獎勵彈窗 ─────────────────────────────────────────────────

func _create_reward_popup() -> void:
	_reward_popup = Label.new()
	_reward_popup.visible = false
	_reward_popup.position = Vector2(540, 350)
	_reward_popup.add_theme_font_size_override("font_size", 24)
	_reward_popup.z_index = 50
	add_child(_reward_popup)

# ── 斷線提示 ─────────────────────────────────────────────────

func _create_disconnect_overlay() -> void:
	_disconnect_overlay = Control.new()
	_disconnect_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_disconnect_overlay.visible = false
	_disconnect_overlay.z_index = 100
	add_child(_disconnect_overlay)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0, 0, 0, 0.7)
	_disconnect_overlay.add_child(bg)

	var msg = Label.new()
	msg.text = "📡 DISCONNECTED\nReconnecting..."
	msg.position = Vector2(500, 330)
	msg.add_theme_font_size_override("font_size", 22)
	msg.modulate = Color(1.0, 0.4, 0.4)
	_disconnect_overlay.add_child(msg)

func _on_disconnected() -> void:
	if is_instance_valid(_disconnect_overlay):
		_disconnect_overlay.visible = true

func _on_reconnected() -> void:
	if is_instance_valid(_disconnect_overlay):
		_disconnect_overlay.visible = false

# ── DAY-298 幸運特殊魚 UI（改由 LuckyEventSystem 處理）────────

func _show_lucky_banner(text: String, color: Color, duration: float = 2.5) -> void:
	if is_instance_valid(lucky_event_system):
		lucky_event_system.show_banner(text, color, duration)
	else:
		# 備用：直接在 HUD 上顯示簡單橫幅
		_show_fallback_banner(text, color, duration)

func _show_lucky_event(lucky_key: String, msg: String, duration: float = 2.5) -> void:
	if is_instance_valid(lucky_event_system):
		lucky_event_system.show_lucky_banner(lucky_key, msg, duration)
	else:
		_show_fallback_banner(msg, Color.WHITE, duration)

func _update_lucky_indicator(title: String, value: String, bar_pct: float = -1.0, color: Color = Color(1.0, 0.85, 0.0)) -> void:
	if is_instance_valid(lucky_event_system):
		lucky_event_system.update_indicator(title, value, bar_pct, color)

func _hide_lucky_indicator() -> void:
	if is_instance_valid(lucky_event_system):
		lucky_event_system.hide_indicator()

func _show_lucky_settle(lines: Array, duration: float = 3.5) -> void:
	if is_instance_valid(lucky_event_system):
		lucky_event_system.show_settle(lines, duration)

# 備用橫幅（LuckyEventSystem 不可用時）
var _fallback_banner: Control = null
func _show_fallback_banner(text: String, color: Color, duration: float = 2.5) -> void:
	if not is_instance_valid(_fallback_banner):
		_fallback_banner = Control.new()
		_fallback_banner.position = Vector2(0, 120)
		_fallback_banner.size = Vector2(1280, 80)
		_fallback_banner.z_index = 60
		add_child(_fallback_banner)
		var bg = ColorRect.new()
		bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		bg.color = Color(0, 0, 0, 0.75)
		_fallback_banner.add_child(bg)
		var lbl = Label.new()
		lbl.name = "BannerLabel"
		lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
		lbl.add_theme_font_size_override("font_size", 28)
		_fallback_banner.add_child(lbl)
	var lbl = _fallback_banner.get_node_or_null("BannerLabel")
	if is_instance_valid(lbl):
		lbl.text = text
		lbl.modulate = color
	_fallback_banner.visible = true
	_fallback_banner.modulate.a = 1.0
	var tween = create_tween()
	tween.tween_interval(duration - 0.5)
	tween.tween_property(_fallback_banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_fallback_banner):
			_fallback_banner.visible = false
	)

# ── DAY-297 Combo UI ──────────────────────────────────────────

func _create_combo_label() -> void:
	_combo_label = Label.new()
	_combo_label.name = "ComboLabel"
	_combo_label.position = Vector2(830, 8)
	_combo_label.size = Vector2(160, 30)
	_combo_label.add_theme_font_size_override("font_size", 16)
	_combo_label.modulate = Color(1.0, 0.85, 0.0)
	_combo_label.visible = false
	add_child(_combo_label)

func _update_combo_display() -> void:
	if not is_instance_valid(_combo_label):
		return
	var combo = GameManager.get_combo_count()
	var bonus = GameManager.get_combo_mult_bonus()
	if combo < 5:
		_combo_label.visible = false
		return
	_combo_label.visible = true
	_combo_label.text = "🔥 COMBO x%d (+%.0f%%)" % [combo, bonus * 100]
	if combo >= 30:
		_combo_label.modulate = Color(1.0, 0.3, 0.1)
	elif combo >= 20:
		_combo_label.modulate = Color(1.0, 0.6, 0.1)
	elif combo >= 10:
		_combo_label.modulate = Color(1.0, 0.85, 0.0)
	else:
		_combo_label.modulate = Color(1.0, 1.0, 0.5)
	# 新 Combo 等級時放大動畫
	if combo != _last_combo and combo in [5, 10, 20, 30]:
		var tween = _combo_label.create_tween()
		tween.tween_property(_combo_label, "scale", Vector2(1.4, 1.4), 0.1)
		tween.tween_property(_combo_label, "scale", Vector2(1.0, 1.0), 0.1)
		ScreenShake.add_trauma(0.2)
	_last_combo = combo

func _on_lucky_chain_lightning(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_event("chain_lightning", "⚡ %s 觸發連鎖閃電！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"chain_hit":
			var chain = data.get("chain_count", 0)
			var mult = data.get("multiplier", 1.0)
			_show_lucky_banner("⚡ 連鎖 %d！×%.1f" % [chain, mult], Color(0.0, 0.9, 1.0), 1.0)
		"settle":
			var reward = data.get("total_reward", 0)
			if reward > 0:
				_show_reward_popup(reward, data.get("multiplier", 1.0))

func _on_lucky_crab_torpedo(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_event("crab_torpedo", "🦀 %s 發射螃蟹魚雷！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.35)
		"explosion":
			var no = data.get("explosion_no", 1)
			_show_lucky_banner("💥 魚雷爆炸 %d/3！" % no, Color(1.0, 0.6, 0.2), 0.8)
			ScreenShake.add_trauma(0.5)
		"settle":
			var reward = data.get("total_reward", 0)
			if reward > 0:
				_show_reward_popup(reward, 3.0)

func _on_lucky_vortex(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_event("vortex", "🌀 %s 召喚渦旋海葵！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"pull":
			var tl = data.get("time_left", 0.0)
			_update_lucky_indicator("🌀 渦旋海葵", "%.0fs" % tl, tl / 8.0, Color(0.7, 0.4, 1.0))
		"end":
			_hide_lucky_indicator()
			_show_lucky_banner("🌀 渦旋爆炸！全場 HP -20%！", Color(0.8, 0.5, 1.0))
			ScreenShake.add_trauma(0.6)

func _on_lucky_golden_dragon(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_event("golden_dragon", "🐉 %s 觸發黃金龍魚輪盤！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"spin":
			var inner = data.get("inner_mult", 1.0)
			var outer = data.get("outer_mult", 1.0)
			var final_m = data.get("final_mult", 1.0)
			_show_lucky_banner("🐉 內環 ×%.0f × 外環 ×%.0f = ×%.0f！" % [inner, outer, final_m], Color(1.0, 0.85, 0.0), 3.0)
		"result":
			var reward = data.get("reward", 0)
			var final_m = data.get("final_mult", 1.0)
			if reward > 0:
				_show_reward_popup(reward, final_m)
			if final_m >= 100:
				ScreenShake.add_trauma(0.8)

func _on_lucky_thunder_lobster(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_event("thunder_lobster", "🦞⚡ %s 觸發雷霆龍蝦！15 秒免費射擊！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"auto_fire":
			var tl = data.get("time_left", 0.0)
			var kills = data.get("kill_count", 0)
			_update_lucky_indicator("🦞⚡ 雷霆模式", "%.0fs | %d 條" % [tl, kills], tl / 15.0, Color(1.0, 0.5, 0.2))
		"end":
			_hide_lucky_indicator()
			var reward = data.get("total_reward", 0)
			var kills = data.get("kill_count", 0)
			_show_lucky_banner("🦞 雷霆結束！擊破 %d 條，獎勵 %d！" % [kills, reward], Color(1.0, 0.7, 0.3))
			if reward > 0:
				_show_reward_popup(reward, float(kills))

func _on_announce(data: Dictionary) -> void:
	var msg = data.get("message", "")
	var priority = data.get("priority", "normal")
	var color_str = data.get("color", "#FFFFFF")
	var color = Color.WHITE
	# 解析 hex 顏色
	if color_str.begins_with("#") and color_str.length() == 7:
		var r = color_str.substr(1, 2).hex_to_int() / 255.0
		var g = color_str.substr(3, 2).hex_to_int() / 255.0
		var b = color_str.substr(5, 2).hex_to_int() / 255.0
		color = Color(r, g, b)

	var duration = 2.0
	match priority:
		"high": duration = 3.0
		"critical": duration = 4.0

	_show_lucky_banner(msg, color, duration)

func _process_announce_queue() -> void:
	pass  # DAY-298：佇列邏輯已移至 LuckyEventSystem

# ── DAY-293 新增幸運特殊魚事件處理 ───────────────────────────

func _on_lucky_awakened_phoenix(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"awaken_start":
			_show_lucky_event("awakened_phoenix", "🔥 %s 觸發覺醒鳳凰！下 5 次攻擊 Power Up！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.35)
		"power_up":
			var mult = data.get("power_up_mult", 6.0)
			var shots = data.get("shots_left", 0)
			_update_lucky_indicator("🔥 覺醒鳳凰", "×%.0f | 剩 %d 次" % [mult, shots], float(shots) / 5.0, Color(1.0, 0.6, 0.2))
		"perfect_awaken":
			_hide_lucky_indicator()
			_show_lucky_banner("🔥✨ 完美覺醒！%s 全服 ×2.0 加成 8 秒！" % name, Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.6)
		"perfect_end":
			_show_lucky_banner("🔥 完美覺醒加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"awaken_end":
			_hide_lucky_indicator()
			var reward = data.get("total_reward", 0)
			var hits = data.get("hit_count", 0)
			_show_lucky_banner("🔥 覺醒結束！命中 %d 次，獎勵 %d！" % [hits, reward], Color(1.0, 0.7, 0.3))
			if reward > 0:
				_show_reward_popup(reward, float(hits))

func _on_lucky_shockwave_bomb(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"shockwave_start":
			_show_lucky_event("shockwave_bomb", "💥 %s 觸發全場震盪！全場 HP -35%！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.7)
		"shockwave_hit":
			var hits = data.get("hit_count", 0)
			var reward = data.get("total_reward", 0)
			_show_lucky_banner("💥 震盪命中 %d 個目標！獎勵 %d！" % [hits, reward], Color(1.0, 0.5, 0.2))
			if reward > 0:
				_show_reward_popup(reward, float(hits) * 0.5)
		"super_shockwave":
			_show_lucky_banner("💥🌊 超級震盪！%s 全服 ×1.8 加成 6 秒！" % name, Color(1.0, 0.42, 0.21), 3.5)
			ScreenShake.add_trauma(0.8)
		"super_end":
			_show_lucky_banner("💥 超級震盪加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"power_end":
			_show_lucky_banner("💥 %s 的震盪強化結束" % name, Color(0.6, 0.6, 0.6), 1.5)

# ── DAY-294 新增幸運特殊魚事件處理 ───────────────────────────

func _on_lucky_drill_torpedo(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_event("drill_torpedo", "🚀 %s 發射鑽頭魚雷！穿透最多 5 個目標！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"penetrate":
			var cnt = data.get("penetrate_cnt", 0)
			var mult = data.get("accum_mult", 1.0)
			_update_lucky_indicator("🚀 鑽頭魚雷", "穿透 %d | ×%.1f" % [cnt, mult], float(cnt) / 5.0, Color(1.0, 0.55, 0.15))
		"explode":
			_hide_lucky_indicator()
			_show_lucky_banner("💥 魚雷終點爆炸！AOE 傷害！", Color(1.0, 0.4, 0.1))
			ScreenShake.add_trauma(0.55)
		"perfect":
			_show_lucky_banner("🚀💥 完美穿透！%s 全服 ×2.2 加成 6 秒！" % name, Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.65)
		"perfect_end":
			_show_lucky_banner("🚀 完美穿透加成結束", Color(0.7, 0.7, 0.7), 1.5)

func _on_lucky_time_freeze(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"freeze_start":
			_show_lucky_event("time_freeze", "❄️ %s 觸發時間凍結！全場凍結 8 秒！傷害 ×1.8！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.3)
		"freeze_end":
			_hide_lucky_indicator()
			_show_lucky_banner("❄️💥 冰裂爆炸！全場 HP -25%！", Color(0.6, 0.9, 1.0))
			ScreenShake.add_trauma(0.5)
		"perfect_freeze":
			var kills = data.get("kill_count", 0)
			_show_lucky_banner("❄️✨ 完美凍結！%s 擊破 %d 條！全服 ×2.0 加成 5 秒！" % [name, kills], Color(0.0, 0.9, 1.0), 3.5)
			ScreenShake.add_trauma(0.6)
		"perfect_end":
			_show_lucky_banner("❄️ 完美凍結加成結束", Color(0.7, 0.7, 0.7), 1.5)

func _on_lucky_chain_explosion(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"chain_start":
			_show_lucky_event("chain_explosion", "💥 %s 觸發連鎖爆炸！12 秒連鎖模式！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"chain_explode":
			var cnt = data.get("chain_count", 0)
			var mult = data.get("accum_mult", 1.0)
			_update_lucky_indicator("💥 連鎖爆炸", "×%.1f | %d 次" % [mult, cnt], -1.0, Color(0.9, 0.2, 0.15))
			ScreenShake.add_trauma(0.25)
		"chain_burst":
			_hide_lucky_indicator()
			_show_lucky_banner("💥🔥 連鎖爆發！%s 全服 ×2.5 加成 6 秒！" % name, Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.7)
		"burst_end":
			_show_lucky_banner("💥 連鎖爆發加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"chain_end":
			_hide_lucky_indicator()
			var cnt = data.get("chain_count", 0)
			var reward = data.get("total_reward", 0)
			_show_lucky_banner("💥 連鎖結束！%d 次連鎖，獎勵 %d！" % [cnt, reward], Color(1.0, 0.6, 0.3))
			if reward > 0:
				_show_reward_popup(reward, float(cnt) * 0.5)

# ── DAY-295 新增幸運特殊魚事件處理 ───────────────────────────

func _on_lucky_chain_long_king(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_event("chain_long_king", "🐉👑 %s 觸發千龍王輪盤！最高 1000x！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.5)
		"spin":
			var inner = data.get("inner_mult", 1.0)
			var outer = data.get("outer_mult", 1.0)
			var final_m = data.get("final_mult", 1.0)
			_show_lucky_banner("🐉 內環 ×%.0f × 外環 ×%.0f = ×%.0f！" % [inner, outer, final_m], Color(1.0, 0.85, 0.0), 3.0)
		"result":
			var reward = data.get("reward", 0)
			var final_m = data.get("final_mult", 1.0)
			if reward > 0:
				_show_reward_popup(reward, final_m)
			if final_m >= 200:
				ScreenShake.add_trauma(0.7)
		"mega_win":
			var final_m = data.get("final_mult", 1.0)
			_show_lucky_banner("🐉👑✨ MEGA WIN！×%.0f！千龍王降臨！" % final_m, Color(1.0, 0.85, 0.0), 4.0)
			ScreenShake.add_trauma(1.0)

func _on_lucky_dragon_shotgun(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_event("dragon_shotgun", "🐲💥 %s 觸發龍力散彈！8 方向攻擊！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.45)
		"shotgun_fire":
			var dir = data.get("direction", 0)
			var hits = data.get("total_hits", 0)
			_show_lucky_banner("🐲 方向 %d 命中！總計 %d 個！" % [dir + 1, hits], Color(0.9, 0.4, 1.0), 0.6)
			ScreenShake.add_trauma(0.2)
		"settle":
			var reward = data.get("total_reward", 0)
			var hits = data.get("total_hits", 0)
			_show_lucky_banner("🐲 散彈結束！命中 %d 個，獎勵 %d！" % [hits, reward], Color(0.8, 0.5, 1.0))
			if reward > 0:
				_show_reward_popup(reward, float(hits) * 0.4)

func _on_lucky_rocket_cannon(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_event("rocket_cannon", "🚀💥 %s 召喚火箭砲！3 枚火箭！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"rocket_launch":
			var no = data.get("rocket_no", 1)
			_show_lucky_banner("🚀 第 %d 枚火箭發射！" % no, Color(1.0, 0.5, 0.2), 0.6)
		"rocket_explode":
			var no = data.get("rocket_no", 1)
			var hits = data.get("hit_targets", [])
			_show_lucky_banner("💥 火箭 %d 爆炸！命中 %d 個！" % [no, hits.size()], Color(1.0, 0.4, 0.1), 0.8)
			ScreenShake.add_trauma(0.5)
		"settle":
			var reward = data.get("total_reward", 0)
			_show_lucky_banner("🚀 火箭砲結束！獎勵 %d！" % reward, Color(1.0, 0.6, 0.3))
			if reward > 0:
				_show_reward_popup(reward, 3.0)

func _on_lucky_deep_whirlpool(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_event("deep_whirlpool", "🌊🌀 %s 觸發深海漩渦！全場 HP -50%！6 秒！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.5)
		"whirlpool_damage":
			var hits = data.get("hit_count", 0)
			_show_lucky_banner("🌀 漩渦傷害！命中 %d 個！" % hits, Color(0.2, 0.7, 1.0), 0.7)
		"settle":
			var reward = data.get("total_reward", 0)
			_show_lucky_banner("🌊 深海漩渦結束！獎勵 %d！" % reward, Color(0.0, 0.8, 1.0))
			if reward > 0:
				_show_reward_popup(reward, 5.0)

func _on_lucky_vampire_mult(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_event("vampire_mult", "🧛 %s 觸發吸血鬼！每次擊破吸收倍率！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.3)
		"absorb":
			var cnt = data.get("absorb_count", 0)
			var mult = data.get("current_mult", 1.0)
			_show_lucky_banner("🧛 吸收 %d 次！當前 ×%.1f" % [cnt, mult], Color(0.7, 0.2, 0.8), 0.7)
		"mult_mode":
			var mult = data.get("current_mult", 5.0)
			_show_lucky_banner("🧛✨ %s 進入倍率模式！×%.1f！10 秒！" % [name, mult], Color(0.8, 0.0, 0.8), 3.5)
			ScreenShake.add_trauma(0.6)
		"mult_end":
			_show_lucky_banner("🧛 吸血鬼倍率模式結束", Color(0.5, 0.5, 0.5), 1.5)
		"settle":
			var cnt = data.get("absorb_count", 0)
			var mult = data.get("current_mult", 1.0)
			_show_lucky_banner("🧛 吸血結束！吸收 %d 次，最終 ×%.1f！" % [cnt, mult], Color(0.6, 0.1, 0.7))

# ── DAY-296 新增幸運特殊魚事件處理 ───────────────────────────

func _on_lucky_mirror_fish(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			var shots = data.get("shots_left", 3)
			_show_lucky_event("mirror_fish", "🪞 %s 觸發鏡像魚！下 %d 次攻擊自動複製！" % [name, shots])
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.3)
		"mirror_hit":
			var shots = data.get("shots_left", 0)
			var hits = data.get("hit_count", 0)
			_show_lucky_banner("🪞 鏡像命中！已命中 %d 次，剩餘 %d 次" % [hits, shots], Color(0.88, 0.67, 1.0), 0.8)
		"perfect_mirror":
			_show_lucky_banner("🪞✨ 完美鏡像！%s 全服 ×1.8 加成 5 秒！" % name, Color(0.88, 0.67, 1.0), 3.5)
			ScreenShake.add_trauma(0.5)
		"perfect_end":
			_show_lucky_banner("🪞 完美鏡像加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"settle":
			var hits = data.get("hit_count", 0)
			var reward = data.get("total_reward", 0)
			_show_lucky_banner("🪞 鏡像結束！命中 %d 次，獎勵 %d！" % [hits, reward], Color(0.8, 0.6, 1.0))
			if reward > 0:
				_show_reward_popup(reward, float(hits))
		"timeout":
			_show_lucky_banner("🪞 鏡像時間到！", Color(0.6, 0.6, 0.6), 1.5)

func _on_lucky_golden_rain(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			var total = data.get("total_coins", 10)
			_show_lucky_event("golden_rain", "🌧️💰 %s 觸發黃金雨！%d 個黃金幣！快去收集！" % [name, total])
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.35)
			# 在畫面上生成可點擊的黃金幣
			var positions = data.get("coin_positions", [])
			_spawn_golden_coins(positions)
		"coin_collect":
			var collected = data.get("collected_coins", 0)
			_show_lucky_banner("💰 收集 %d 個黃金幣！" % collected, Color(1.0, 0.85, 0.0), 0.6)
		"golden_harvest":
			var collected = data.get("collected_coins", 0)
			var reward = data.get("total_reward", 0)
			_show_lucky_banner("💰✨ 黃金豐收！%s 收集 %d 個！全服 ×2.0 加成 6 秒！" % [name, collected], Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.6)
			if reward > 0:
				_show_reward_popup(reward, 2.0)
		"harvest_end":
			_show_lucky_banner("💰 黃金豐收加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"settle":
			var collected = data.get("collected_coins", 0)
			var reward = data.get("total_reward", 0)
			_show_lucky_banner("🌧️ 黃金雨結束！收集 %d 個，獎勵 %d！" % [collected, reward], Color(1.0, 0.85, 0.0))
			if reward > 0:
				_show_reward_popup(reward, float(collected) * 0.3)
			_clear_golden_coins()

var _golden_coins: Array = []

func _spawn_golden_coins(positions: Array) -> void:
	_clear_golden_coins()
	for coin_data in positions:
		var coin_id = coin_data.get("coin_id", 0)
		var x = coin_data.get("x", 0.0)
		var y = coin_data.get("y", 0.0)
		var btn = Button.new()
		btn.text = "💰"
		btn.position = Vector2(x - 20, y - 20)
		btn.size = Vector2(40, 40)
		btn.z_index = 55
		btn.add_theme_font_size_override("font_size", 20)
		btn.set_meta("coin_id", coin_id)
		btn.pressed.connect(func():
			if is_instance_valid(btn):
				NetworkManager.send_collect_golden_coin(btn.get_meta("coin_id"))
				var tween = btn.create_tween()
				tween.tween_property(btn, "scale", Vector2(2.0, 2.0), 0.1)
				tween.parallel().tween_property(btn, "modulate:a", 0.0, 0.2)
				tween.tween_callback(func(): if is_instance_valid(btn): btn.queue_free())
		)
		add_child(btn)
		_golden_coins.append(btn)
		# 進場動畫
		btn.scale = Vector2.ZERO
		var tween = btn.create_tween()
		tween.tween_property(btn, "scale", Vector2(1.0, 1.0), 0.2).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)

func _clear_golden_coins() -> void:
	for coin in _golden_coins:
		if is_instance_valid(coin):
			coin.queue_free()
	_golden_coins.clear()

func _on_lucky_freeze_bomb(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"freeze_start":
			var frozen = data.get("frozen_targets", [])
			_show_lucky_event("freeze_bomb", "❄️💣 %s 投擲冰凍炸彈！%d 個目標凍結！3 秒後爆炸！" % [name, frozen.size()])
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.3)
		"bomb_explode":
			var hits = data.get("hit_count", 0)
			_show_lucky_banner("❄️💥 冰凍炸彈爆炸！命中 %d 個！HP -60%！" % hits, Color(0.4, 0.9, 1.0))
			ScreenShake.add_trauma(0.65)
		"perfect_freeze":
			var hits = data.get("hit_count", 0)
			_show_lucky_banner("❄️💥✨ 冰爆完美！%s 命中 %d 個！全服 ×2.2 加成 5 秒！" % [name, hits], Color(0.0, 0.9, 1.0), 3.5)
			ScreenShake.add_trauma(0.7)
		"perfect_end":
			_show_lucky_banner("❄️ 冰爆完美加成結束", Color(0.7, 0.7, 0.7), 1.5)

func _on_lucky_thunder_storm(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"storm_start":
			var count = data.get("lightning_count", 6)
			_show_lucky_event("thunder_storm", "⛈️ %s 召喚雷暴！%d 道閃電！10 秒！" % [name, count])
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"lightning_strike":
			var no = data.get("strike_no", 1)
			var hits = data.get("hit_targets", [])
			var mult = data.get("accum_mult", 1.0)
			_show_lucky_banner("⚡ 第 %d 道閃電！命中 %d 個！累積 ×%.1f" % [no, hits.size(), mult], Color(1.0, 0.9, 0.2), 0.7)
			if hits.size() > 0:
				ScreenShake.add_trauma(0.2)
		"perfect_storm":
			var strikes = data.get("hit_strikes", 6)
			_show_lucky_banner("⛈️✨ 雷暴完美！%s %d 道全命中！全服 ×2.3 加成 6 秒！" % [name, strikes], Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.7)
		"perfect_end":
			_show_lucky_banner("⛈️ 雷暴完美加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"storm_end":
			var strikes = data.get("hit_strikes", 0)
			var mult = data.get("accum_mult", 1.0)
			var reward = data.get("total_reward", 0)
			_show_lucky_banner("⛈️ 雷暴結束！%d 道命中，累積 ×%.1f，獎勵 %d！" % [strikes, mult, reward], Color(1.0, 0.85, 0.0))
			if reward > 0:
				_show_reward_popup(reward, mult)

func _on_lucky_lucky_wheel(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			var pool = data.get("pool_size", 20000)
			_show_lucky_event("lucky_wheel", "🎡 %s 觸發幸運大轉盤！大獎池 %d！" % [name, pool])
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.35)
		"spin_result":
			var slot_name = data.get("slot_name", "×2")
			var slot_type = data.get("slot_type", "mult")
			var reward = data.get("reward", 0)
			match slot_type:
				"jackpot":
					_show_lucky_banner("🎡🏆 %s 中大獎！%s！獎勵 %d！" % [name, slot_name, reward], Color(1.0, 0.85, 0.0), 4.0)
					ScreenShake.add_trauma(0.8)
					if reward > 0:
						_show_reward_popup(reward, 100.0)
				"aoe":
					_show_lucky_banner("🎡💥 %s 轉到 %s！全場 HP -50%！" % [name, slot_name], Color(1.0, 0.42, 0.71))
					ScreenShake.add_trauma(0.5)
				"mult":
					var mult = data.get("slot_mult", 2.0)
					_show_lucky_banner("🎡 %s 轉到 %s！獎勵 %d！" % [name, slot_name, reward], Color(1.0, 0.42, 0.71))
					if reward > 0:
						_show_reward_popup(reward, mult)

# ── DAY-301 新增幸運特殊魚事件處理 ───────────────────────────

func _on_lucky_jackpot_fish(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_event("jackpot_fish", "🏆 %s 觸發進階 Jackpot！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"jackpot_result":
			var tier = data.get("tier_name", "Mini")
			var reward = data.get("reward", 0)
			var tier_idx = data.get("tier_idx", 0)
			var colors = [Color(0.7, 0.4, 0.2), Color(0.8, 0.8, 0.9), Color(1.0, 0.55, 0.0), Color(1.0, 0.85, 0.0)]
			_show_lucky_banner("🏆 %s 中 %s Jackpot！獲得 %d！" % [name, tier, reward], colors[clamp(tier_idx, 0, 3)], 3.0)
			if reward > 0:
				_show_reward_popup(reward, float(tier_idx + 1) * 10.0)
			if tier_idx == 3:
				ScreenShake.add_trauma(0.8)
		"grand_boost":
			var mult = data.get("boost_mult", 3.0)
			var secs = data.get("boost_secs", 10)
			_show_lucky_banner("🏆✨ GRAND JACKPOT！%s 全服 ×%.0f 加成 %d 秒！" % [name, mult, secs], Color(1.0, 0.85, 0.0), 4.0)
			ScreenShake.add_trauma(0.9)
		"grand_boost_end":
			_show_lucky_banner("🏆 Grand Jackpot 加成結束", Color(0.7, 0.7, 0.7), 1.5)

func _on_lucky_coop_fish(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"coop_start":
			var target = data.get("target_points", 8)
			_show_lucky_event("coop_fish", "🤝 %s 發起全服合作！目標 %d 點！20 秒！" % [name, target])
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.3)
		"coop_progress":
			var current = data.get("current_points", 0)
			var target = data.get("target_points", 8)
			var tl = data.get("time_left", 0.0)
			_update_lucky_indicator("🤝 全服合作", "%d/%d 點" % [current, target], float(current) / float(max(target, 1)), Color(0.0, 0.9, 1.0))
		"coop_success":
			_hide_lucky_indicator()
			var boost = data.get("boost_mult", 4.0)
			var secs = data.get("boost_secs", 8)
			_show_lucky_banner("🤝✨ 全服合作成功！全服 ×%.0f 加成 %d 秒！" % [boost, secs], Color(0.0, 1.0, 0.5), 3.5)
			ScreenShake.add_trauma(0.7)
		"coop_timeout":
			_hide_lucky_indicator()
			var current = data.get("current_points", 0)
			var target = data.get("target_points", 8)
			_show_lucky_banner("🤝 合作挑戰時間到！達成 %d/%d 點" % [current, target], Color(0.6, 0.6, 0.6), 2.0)
		"coop_boost_end":
			_show_lucky_banner("🤝 全服合作加成結束", Color(0.6, 0.6, 0.6), 1.5)

func _on_lucky_time_warp(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"warp_start":
			var duration = data.get("duration", 10.0)
			var dmg = data.get("damage_mult", 2.0)
			_show_lucky_event("time_warp", "⏰ %s 觸發時間扭曲！全場慢速 %.0f 秒！傷害 ×%.0f！" % [name, duration, dmg])
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.35)
		"warp_end":
			_hide_lucky_indicator()
			var kills = data.get("kill_count", 0)
			_show_lucky_banner("⏰💥 時間扭曲結束！全場 HP -20%！擊破 %d 條！" % kills, Color(0.55, 0.2, 0.86))
			ScreenShake.add_trauma(0.5)
		"time_collapse":
			var kills = data.get("kill_count", 0)
			var boost = data.get("boost_mult", 2.5)
			var secs = data.get("boost_secs", 6)
			_show_lucky_banner("⏰💥 時間崩潰！%s 擊破 %d 條！全服 ×%.0f 加成 %d 秒！" % [name, kills, boost, secs], Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.7)
		"collapse_end":
			_show_lucky_banner("⏰ 時間崩潰加成結束", Color(0.6, 0.6, 0.6), 1.5)

# ── DAY-302 新增幸運特殊魚事件處理 ───────────────────────────

func _on_lucky_chain_meteor(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("player_name", "玩家")
	match event:
		"meteor_start":
			_show_lucky_event("chain_meteor", "☄️ %s 觸發連鎖隕石雨！5 顆隕石依序落下！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.5)
		"meteor_hit":
			var idx = data.get("meteor_index", 1)
			var radius = data.get("aoe_radius", 150.0)
			var hits = data.get("hit_count", 0)
			_update_lucky_indicator("☄️ 連鎖隕石", "第 %d/5 顆 r=%.0f 命中 %d" % [idx, radius, hits], float(idx) / 5.0, Color(1.0, 0.5, 0.2))
			ScreenShake.add_trauma(0.35)
		"meteor_miss":
			var idx = data.get("meteor_index", 1)
			_update_lucky_indicator("☄️ 連鎖隕石", "第 %d/5 顆 空揮！" % idx, float(idx) / 5.0, Color(0.6, 0.6, 0.6))
		"meteor_perfect":
			_hide_lucky_indicator()
			_show_lucky_banner("☄️✨ 完美隕石雨！%s 全服 ×2.5 加成 7 秒！" % name, Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.8)
		"meteor_perfect_end":
			_show_lucky_banner("☄️ 完美隕石雨加成結束", Color(0.7, 0.7, 0.7), 1.5)

# ── DAY-303 新增幸運特殊魚事件處理 ───────────────────────────

func _on_lucky_crash_fish(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("player_name", "玩家")
	var mult = data.get("current_mult", 1.0)
	match event:
		"crash_start":
			_show_lucky_event("crash_fish", "💥 %s 觸發崩潰倍率！倍率持續上升！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"mult_rise":
			_update_lucky_indicator("💥 崩潰倍率", "×%.1f" % mult, -1.0, Color(0.8, 0.1, 0.1))
		"harvest":
			_hide_lucky_indicator()
			var reward = data.get("reward", 0)
			_show_lucky_banner("💰 %s 收割！×%.1f！獎勵 %d！" % [name, mult, reward], Color(0.2, 0.9, 0.2), 2.5)
			if reward > 0:
				_show_reward_popup(reward, mult)
		"crash":
			_hide_lucky_indicator()
			_show_lucky_banner("💥 崩潰！%s 的 ×%.1f 歸零！" % [name, mult], Color(0.8, 0.1, 0.1), 2.5)
			ScreenShake.add_trauma(0.6)
		"perfect_harvest":
			var boost = data.get("boost_mult", 2.0)
			var secs = data.get("boost_secs", 5)
			_show_lucky_banner("💰✨ 完美收割！%s ×%.1f！全服 ×%.0f 加成 %d 秒！" % [name, mult, boost, secs], Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.7)
		"perfect_end":
			_show_lucky_banner("💰 完美收割加成結束", Color(0.7, 0.7, 0.7), 1.5)
