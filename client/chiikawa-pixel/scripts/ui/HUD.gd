## HUD.gd — 核心 HUD
## hud-core-agent 負責維護
## DAY-298：Lucky 事件視覺系統改由 LuckyEventSystem 處理
## DAY-336：Lucky 訊號連接全部移入 HUDLuckySignals.gd
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
# DAY-336：HUDLuckySignals 模組（管理所有 148 個 Lucky 訊號）
var _lucky_signals: Node = null
# DAY-342 在線玩家數顯示
var _online_label: Label = null
# DAY-345 每日任務按鈕
var _quest_btn: Button = null
var _daily_quest_panel: Node = null
# DAY-346 每週挑戰按鈕
var _weekly_btn: Button = null
var _weekly_challenge_panel: Node = null

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
	_create_online_label()
	_create_quest_button()
	_create_weekly_button()
	_update_ui()
	# 嘗試自動找 LuckyEventSystem（如果在同一場景樹中）
	call_deferred("_find_lucky_event_system")

	# DAY-336：初始化 HUDLuckySignals 模組（管理所有 148 個 Lucky 訊號）
	_lucky_signals = Node.new()
	_lucky_signals.set_script(load("res://scripts/ui/HUDLuckySignals.gd"))
	add_child(_lucky_signals)
	call_deferred("_init_lucky_signals")

	# DAY-337：Lucky 訊號連接已全部移入 HUDLuckySignals.gd（由 _init_lucky_signals 初始化）

func _process(delta: float) -> void:
	if _boss_active and _boss_time_left > 0:
		_boss_time_left -= delta
		_update_boss_timer()

func _init_lucky_signals() -> void:
	# DAY-336：初始化 HUDLuckySignals 模組
	if is_instance_valid(_lucky_signals):
		_lucky_signals.lucky_event_system = lucky_event_system
		_lucky_signals.connect_all_lucky_signals(self)
		print("[HUD] HUDLuckySignals 初始化完成（DAY-304~319 訊號已接管）")

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

# DAY-340 金幣計數動畫：數字跳動效果
var _displayed_coins: int = 0
var _target_coins: int = 0
var _coin_tween: Tween = null

func _animate_coins_to(new_coins: int) -> void:
	if new_coins == _displayed_coins:
		return
	_target_coins = new_coins
	if is_instance_valid(_coin_tween):
		_coin_tween.kill()
	_coin_tween = create_tween()
	var diff = abs(new_coins - _displayed_coins)
	# 差距越大，動畫越快（最快 0.3s，最慢 0.8s）
	var duration = clamp(0.05 + diff * 0.002, 0.15, 0.8)
	_coin_tween.tween_method(
		func(v: int):
			_displayed_coins = v
			if is_instance_valid(coins_label):
				coins_label.text = "💰 %d" % v
		,
		_displayed_coins,
		new_coins,
		duration
	)
	# 大獎時金幣標籤彈跳
	if new_coins > _displayed_coins + 50:
		_coin_tween.parallel().tween_property(coins_label, "scale", Vector2(1.3, 1.3), 0.1)
		_coin_tween.tween_property(coins_label, "scale", Vector2(1.0, 1.0), 0.1)

func _update_ui() -> void:
	# 金幣：使用動畫計數
	var new_coins = GameManager.get_coins()
	if new_coins != _target_coins:
		_animate_coins_to(new_coins)
	else:
		coins_label.text = "💰 %d" % new_coins
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
	# DAY-342 在線玩家數顯示
	_update_online_display()

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

# ── DAY-297/322 Combo UI ──────────────────────────────────────

func _create_combo_label() -> void:
	# DAY-341 升級：更醒目的 Combo 顯示（帶背景面板）
	var combo_panel = Control.new()
	combo_panel.name = "ComboPanel"
	combo_panel.position = Vector2(800, 2)
	combo_panel.size = Vector2(240, 40)
	combo_panel.z_index = 55
	add_child(combo_panel)
	
	# 半透明背景
	var bg = ColorRect.new()
	bg.size = combo_panel.size
	bg.color = Color(0.0, 0.0, 0.0, 0.6)
	combo_panel.add_child(bg)
	
	_combo_label = Label.new()
	_combo_label.name = "ComboLabel"
	_combo_label.position = Vector2(4, 4)
	_combo_label.size = Vector2(232, 32)
	_combo_label.add_theme_font_size_override("font_size", 18)
	_combo_label.modulate = Color(1.0, 0.85, 0.0)
	combo_panel.add_child(_combo_label)
	combo_panel.visible = false
	# 讓 _combo_label 的 visible 控制整個 panel
	# 用 meta 記錄 panel 引用
	_combo_label.set_meta("panel", combo_panel)
	
	# DAY-322：連接 ComboSystem 訊號（如果存在）
	call_deferred("_connect_combo_system")

func _connect_combo_system() -> void:
	var tree = get_tree()
	if tree == null:
		return
	var combo_sys = tree.get_root().find_child("ComboSystem", true, false)
	if is_instance_valid(combo_sys):
		if combo_sys.has_signal("combo_updated") and not combo_sys.combo_updated.is_connected(_on_combo_updated):
			combo_sys.combo_updated.connect(_on_combo_updated)
		if combo_sys.has_signal("combo_milestone") and not combo_sys.combo_milestone.is_connected(_on_combo_milestone):
			combo_sys.combo_milestone.connect(_on_combo_milestone)

func _on_combo_updated(count: int) -> void:
	if not is_instance_valid(_combo_label):
		return
	var panel: Control = null
	if _combo_label.has_meta("panel"):
		panel = _combo_label.get_meta("panel")
	if count < 5:
		if is_instance_valid(panel):
			panel.visible = false
		return
	if is_instance_valid(panel):
		panel.visible = true
	_combo_label.text = "🔥 COMBO x%d" % count
	if count >= 50:
		_combo_label.modulate = Color(1.0, 0.0, 0.5)
	elif count >= 20:
		_combo_label.modulate = Color(0.8, 0.0, 1.0)
	elif count >= 10:
		_combo_label.modulate = Color(1.0, 0.3, 0.1)
	else:
		_combo_label.modulate = Color(1.0, 0.85, 0.0)
	_last_combo = count

func _on_combo_milestone(count: int) -> void:
	if not is_instance_valid(_combo_label):
		return
	var tween = _combo_label.create_tween()
	tween.tween_property(_combo_label, "scale", Vector2(1.6, 1.6), 0.08)
	tween.tween_property(_combo_label, "scale", Vector2(1.0, 1.0), 0.12)

func _update_combo_display() -> void:
	# 備用：如果 ComboSystem 不存在，用 GameManager 的 combo 數據
	if not is_instance_valid(_combo_label):
		return
	var combo = 0
	if GameManager.has_method("get_combo_count"):
		combo = GameManager.get_combo_count()
	
	# 取得 panel 引用
	var panel: Control = null
	if _combo_label.has_meta("panel"):
		panel = _combo_label.get_meta("panel")
	
	if combo < 5:
		if is_instance_valid(panel):
			panel.visible = false
		_last_combo = combo
		return
	
	if is_instance_valid(panel):
		panel.visible = true
	
	var bonus = 0.0
	if GameManager.has_method("get_combo_mult_bonus"):
		bonus = GameManager.get_combo_mult_bonus()
	_combo_label.text = "🔥 COMBO x%d (+%.0f%%)" % [combo, bonus * 100]
	if combo >= 30:
		_combo_label.modulate = Color(1.0, 0.3, 0.1)
	elif combo >= 20:
		_combo_label.modulate = Color(1.0, 0.6, 0.1)
	elif combo >= 10:
		_combo_label.modulate = Color(1.0, 0.85, 0.0)
	else:
		_combo_label.modulate = Color(1.0, 1.0, 0.5)
	# DAY-341 里程碑觸發：音效 + 特效
	if combo != _last_combo and combo in [5, 10, 20, 30]:
		var tween = _combo_label.create_tween()
		tween.tween_property(_combo_label, "scale", Vector2(1.6, 1.6), 0.08)
		tween.tween_property(_combo_label, "scale", Vector2(1.0, 1.0), 0.12)
		ScreenShake.add_trauma(0.2 + (combo / 30.0) * 0.3)
		# 播放對應里程碑音效
		AudioManager.play_combo_milestone(combo)
		# 里程碑特效（在 Combo 標籤位置生成爆炸）
		_spawn_combo_milestone_effect(combo)
	_last_combo = combo

## DAY-341 Combo 里程碑視覺特效
func _spawn_combo_milestone_effect(combo: int) -> void:
	# 在 Combo 標籤位置生成爆炸粒子
	var pos = Vector2(920, 24)  # Combo 標籤中心位置
	var color: Color
	var particle_count: int
	match combo:
		5:
			color = Color(1.0, 1.0, 0.5)
			particle_count = 6
		10:
			color = Color(1.0, 0.85, 0.0)
			particle_count = 10
		20:
			color = Color(1.0, 0.5, 0.0)
			particle_count = 16
		30:
			color = Color(1.0, 0.2, 0.1)
			particle_count = 24
		_:
			return
	
	# 生成粒子
	for i in particle_count:
		var dot = ColorRect.new()
		var size = randf_range(4.0, 8.0)
		dot.size = Vector2(size, size)
		dot.color = color
		dot.position = pos - dot.size / 2
		dot.z_index = 65
		add_child(dot)
		var angle = (float(i) / float(particle_count)) * TAU + randf_range(-0.3, 0.3)
		var dist = randf_range(30.0, 80.0)
		var target = pos + Vector2(cos(angle), sin(angle)) * dist
		var tween = dot.create_tween()
		tween.tween_property(dot, "position", target - dot.size / 2, 0.4).set_ease(Tween.EASE_OUT)
		tween.parallel().tween_property(dot, "modulate:a", 0.0, 0.4)
		tween.tween_callback(func(): if is_instance_valid(dot): dot.queue_free())
	
	# 里程碑文字（30連擊才顯示大文字）
	if combo >= 20:
		var milestone_label = Label.new()
		milestone_label.z_index = 70
		match combo:
			20:
				milestone_label.text = "🔥 COMBO x20!"
				milestone_label.modulate = Color(1.0, 0.5, 0.0)
			30:
				milestone_label.text = "💥 MAX COMBO x30!"
				milestone_label.modulate = Color(1.0, 0.2, 0.1)
		milestone_label.add_theme_font_size_override("font_size", 28)
		milestone_label.position = Vector2(480, 200)
		add_child(milestone_label)
		var tween2 = milestone_label.create_tween()
		tween2.tween_property(milestone_label, "scale", Vector2(1.5, 1.5), 0.1)
		tween2.tween_property(milestone_label, "scale", Vector2(1.0, 1.0), 0.1)
		tween2.tween_interval(0.8)
		tween2.tween_property(milestone_label, "modulate:a", 0.0, 0.4)
		tween2.tween_callback(func(): if is_instance_valid(milestone_label): milestone_label.queue_free())


# ── DAY-337 重構完成 ────────────────────────────────────────
# 所有 Lucky 函數已移入 HUDLuckySignals.gd
# HUD.gd 只保留核心 HUD 功能

# ── DAY-342 在線玩家數顯示 ────────────────────────────────────

func _create_online_label() -> void:
	_online_label = Label.new()
	_online_label.name = "OnlineLabel"
	_online_label.position = Vector2(1050, 6)
	_online_label.size = Vector2(180, 30)
	_online_label.add_theme_font_size_override("font_size", 14)
	_online_label.modulate = Color(0.6, 1.0, 0.6)
	_online_label.text = "👥 1 在線"
	_online_label.z_index = 55
	add_child(_online_label)

func _update_online_display() -> void:
	if not is_instance_valid(_online_label):
		return
	var count = 1
	if GameManager.has_method("get_online_count"):
		count = GameManager.get_online_count()
	if count <= 1:
		_online_label.text = "👤 1 在線"
		_online_label.modulate = Color(0.7, 0.7, 0.7)
	elif count <= 3:
		_online_label.text = "👥 %d 在線" % count
		_online_label.modulate = Color(0.6, 1.0, 0.6)
	else:
		_online_label.text = "👥 %d 在線 🔥" % count
		_online_label.modulate = Color(1.0, 0.85, 0.0)

# ── DAY-345 每日任務按鈕 ──────────────────────────────────────

func _create_quest_button() -> void:
	_quest_btn = Button.new()
	_quest_btn.name = "QuestBtn"
	_quest_btn.text = "🎯"
	_quest_btn.position = Vector2(1050, 40)
	_quest_btn.size = Vector2(40, 30)
	_quest_btn.add_theme_font_size_override("font_size", 16)
	_quest_btn.z_index = 55
	_quest_btn.tooltip_text = "每日任務"
	add_child(_quest_btn)

	# 建立每日任務面板
	_daily_quest_panel = load("res://scripts/ui/DailyQuestPanel.gd").new()
	_daily_quest_panel.name = "DailyQuestPanel"
	add_child(_daily_quest_panel)

	# 連接按鈕
	_quest_btn.pressed.connect(func():
		if is_instance_valid(_daily_quest_panel) and _daily_quest_panel.has_method("_toggle_panel"):
			_daily_quest_panel._toggle_panel()
	)

# ── DAY-346 每週挑戰按鈕 ──────────────────────────────────────

func _create_weekly_button() -> void:
	_weekly_btn = Button.new()
	_weekly_btn.name = "WeeklyBtn"
	_weekly_btn.text = "🏆"
	_weekly_btn.position = Vector2(1000, 40)
	_weekly_btn.size = Vector2(40, 30)
	_weekly_btn.add_theme_font_size_override("font_size", 16)
	_weekly_btn.z_index = 55
	_weekly_btn.tooltip_text = "每週挑戰"
	add_child(_weekly_btn)

	# 建立每週挑戰面板
	_weekly_challenge_panel = load("res://scripts/ui/WeeklyChallengePanel.gd").new()
	_weekly_challenge_panel.name = "WeeklyChallengePanel"
	add_child(_weekly_challenge_panel)

	# 連接按鈕
	_weekly_btn.pressed.connect(func():
		if is_instance_valid(_weekly_challenge_panel) and _weekly_challenge_panel.has_method("_toggle_panel"):
			_weekly_challenge_panel._toggle_panel()
	)
