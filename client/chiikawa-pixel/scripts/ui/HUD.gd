## HUD.gd — 核心 UI（精簡版，移除所有損壞的 Panel）
## 保留：基本 HUD、BOSS 計時器、獎勵彈窗、斷線提示、側錄按鈕
extends CanvasLayer

const ScreenRecorderScript = preload("res://scripts/ui/ScreenRecorder.gd")
var _screen_recorder = null

# ---- 核心 UI 節點 ----
@onready var coins_label: Label = $TopBar/CoinsLabel
@onready var bet_label: Label = $TopBar/BetLabel
@onready var character_label: Label = $TopBar/CharacterLabel
@onready var labor_bar: ProgressBar = $TopBar/LaborBar
@onready var labor_label: Label = $TopBar/LaborLabel
@onready var auto_button: Button = $BottomBar/AutoButton
@onready var lock_button: Button = $BottomBar/LockButton
@onready var bet_minus_button: Button = $BottomBar/BetMinusButton
@onready var bet_plus_button: Button = $BottomBar/BetPlusButton
@onready var boss_button: Button = $BottomBar/BossButton
@onready var bonus_button: Button = $BottomBar/BonusButton
@onready var reward_popup: Label = $RewardPopup
@onready var state_label: Label = $StateLabel
@onready var warning_overlay: Control = $WarningOverlay
@onready var bonus_overlay: Control = $BonusOverlay

var _reward_popup_base_y: float = 0.0
var _lock_active: bool = false
var _boss_time_left: float = 0.0
var _boss_active: bool = false
var _boss_timer_node: Control = null
var _last_labor_value: int = 0

# Lucky Panel 系統（DAY-289/290）
const LuckyImmortalBossPanelScript = preload("res://scripts/ui/LuckyImmortalBossPanel.gd")
const LuckyWrathChargePanelScript = preload("res://scripts/ui/LuckyWrathChargePanel.gd")
var _lucky_immortal_boss_panel = null
var _lucky_wrath_charge_panel = null

func _ready() -> void:
	GameManager.player_updated.connect(_on_player_updated)
	GameManager.game_state_changed.connect(_on_game_state_changed)
	GameManager.reward_received.connect(_on_reward_received)
	GameManager.boss_event.connect(_on_boss_event)
	GameManager.bonus_event.connect(_on_bonus_event)
	NetworkManager.connected.connect(_on_reconnected)

	auto_button.pressed.connect(_on_auto_pressed)
	lock_button.pressed.connect(_on_lock_pressed)
	bet_minus_button.pressed.connect(_on_bet_minus)
	bet_plus_button.pressed.connect(_on_bet_plus)
	boss_button.pressed.connect(NetworkManager.send_trigger_boss)
	bonus_button.pressed.connect(NetworkManager.send_trigger_bonus)

	for btn in [auto_button, lock_button, bet_minus_button, bet_plus_button, boss_button, bonus_button]:
		if is_instance_valid(btn):
			btn.pressed.connect(func(): AudioManager.play_sfx(AudioManager.SFX.WEED_PULL))

	reward_popup.visible = false
	_reward_popup_base_y = reward_popup.position.y
	_update_ui()
	_create_disconnect_overlay()

	# 初始化 Lucky Immortal Boss Panel（DAY-289）
	_lucky_immortal_boss_panel = LuckyImmortalBossPanelScript.new()
	get_tree().root.add_child(_lucky_immortal_boss_panel)
	if GameManager.has_signal("lucky_immortal_boss"):
		GameManager.lucky_immortal_boss.connect(_on_lucky_immortal_boss)

	# 初始化 Lucky Wrath Charge Panel（DAY-290）
	_lucky_wrath_charge_panel = LuckyWrathChargePanelScript.new()
	get_tree().root.add_child(_lucky_wrath_charge_panel)
	if GameManager.has_signal("lucky_wrath_charge"):
		GameManager.lucky_wrath_charge.connect(_on_lucky_wrath_charge)

	# 初始化側錄系統（DAY-291）
	# CanvasLayer 必須加到 root，不能作為另一個 CanvasLayer 的子節點
	_screen_recorder = ScreenRecorderScript.new()
	get_tree().root.add_child(_screen_recorder)

func _process(delta: float) -> void:
	if _boss_active and _boss_time_left > 0:
		_boss_time_left -= delta
		_update_boss_timer()

func _on_player_updated(_data: Dictionary) -> void:
	_update_ui()

func _update_ui() -> void:
	coins_label.text = "💰 %d" % GameManager.get_coins()
	var lv = GameManager.get_bet_level()
	var cost = GameManager.get_bet_cost()
	bet_label.text = "BET LV%d (%d)" % [lv, cost]
	character_label.text = GameManager.get_character_name()
	character_label.modulate = GameManager.get_character_color()
	var labor = GameManager.get_labor_value()
	labor_bar.value = labor
	if labor >= 80:
		labor_label.text = "⚡%d/100" % labor
		labor_label.modulate = Color(1.0, 0.9, 0.2)
	else:
		labor_label.text = "%d/100" % labor
		labor_label.modulate = Color.WHITE
	if labor >= 100 and _last_labor_value < 100:
		ScreenShake.add_trauma(0.3)
	_last_labor_value = labor
	if GameManager.is_auto():
		auto_button.modulate = Color(0.3, 1.0, 0.3)
		auto_button.text = "AUTO ON"
	else:
		auto_button.modulate = Color.WHITE
		auto_button.text = "AUTO"
	var lock_id = GameManager.get_lock_target_id()
	if lock_id != "":
		lock_button.modulate = Color(1.0, 0.8, 0.2)
		lock_button.text = "🔒 LOCK"
		_lock_active = true
	else:
		lock_button.modulate = Color(0.7, 0.7, 0.7)
		lock_button.text = "🔓 LOCK"
		_lock_active = false

func _on_game_state_changed(new_state: String) -> void:
	state_label.text = new_state.to_upper().replace("_", " ")
	match new_state:
		"boss_warning":
			_show_boss_warning()
		"boss_battle":
			warning_overlay.visible = false
		"bonus_game":
			bonus_overlay.visible = true
		"bonus_result", "normal_play", "boss_result":
			bonus_overlay.visible = false
			warning_overlay.visible = false

func _show_boss_warning() -> void:
	warning_overlay.visible = true
	var tween = create_tween().set_loops(8)
	tween.tween_property(warning_overlay, "modulate:a", 0.1, 0.18)
	tween.tween_property(warning_overlay, "modulate:a", 1.0, 0.18)

func _on_reward_received(reward: Dictionary) -> void:
	var amount = reward.get("amount", 0)
	var multiplier = reward.get("multiplier", 1.0)
	if amount <= 0:
		return
	_show_reward_popup(amount, multiplier)

func _show_reward_popup(amount: int, multiplier: float) -> void:
	var icon = "💰"
	if multiplier >= 100:
		icon = "🌟"
		reward_popup.modulate = Color(1.0, 0.3, 0.1, 1.0)
	elif multiplier >= 20:
		icon = "⭐"
		reward_popup.modulate = Color(1.0, 0.85, 0.0, 1.0)
	elif multiplier >= 10:
		icon = "✨"
		reward_popup.modulate = Color(1.0, 1.0, 0.4, 1.0)
	else:
		icon = "💰"
		reward_popup.modulate = Color(1.0, 1.0, 1.0, 1.0)
	reward_popup.text = "%s +%d  x%.0f" % [icon, amount, multiplier]
	reward_popup.position.y = _reward_popup_base_y
	reward_popup.visible = true
	var tween = create_tween()
	tween.tween_property(reward_popup, "position:y", _reward_popup_base_y - 70, 0.7)
	tween.parallel().tween_property(reward_popup, "modulate:a", 0.0, 0.7)
	tween.tween_callback(func():
		reward_popup.visible = false
		reward_popup.position.y = _reward_popup_base_y
	)

func _on_boss_event(event_data: Dictionary) -> void:
	match event_data.get("event", ""):
		"warning":
			AudioManager.stop_bgm_briefly()
			AudioManager.play_sfx(AudioManager.SFX.BOSS_WARNING)
		"spawn":
			AudioManager.play_bgm(AudioManager.BGM.BOSS_ENTER)
			_start_boss_timer()
		"phase_change":
			AudioManager.play_bgm(AudioManager.BGM.BOSS_RAGE)
		"kill", "timeout":
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			AudioManager.play_bgm(AudioManager.BGM.MAIN_GAME)
			_stop_boss_timer()
			if event_data.get("event", "") == "kill":
				HitEffect.spawn_big_win(Vector2(640, 360), 100.0)
				ScreenShake.add_trauma(0.7)

func _start_boss_timer() -> void:
	_boss_time_left = 60.0
	_boss_active = true
	if is_instance_valid(_boss_timer_node):
		_boss_timer_node.queue_free()
	var panel = Control.new()
	panel.name = "BossTimerPanel"
	panel.position = Vector2(900, 50)
	panel.size = Vector2(360, 80)
	add_child(panel)
	_boss_timer_node = panel
	var bg = ColorRect.new()
	bg.size = Vector2(360, 80)
	bg.color = Color(0.1, 0.0, 0.0, 0.85)
	panel.add_child(bg)
	var title = Label.new()
	title.text = "⚔ BOSS BATTLE"
	title.position = Vector2(10, 5)
	title.add_theme_font_size_override("font_size", 16)
	title.modulate = Color(1.0, 0.3, 0.3)
	panel.add_child(title)
	var timer_lbl = Label.new()
	timer_lbl.name = "BossTimeLabel"
	timer_lbl.text = "60.0s"
	timer_lbl.position = Vector2(10, 28)
	timer_lbl.add_theme_font_size_override("font_size", 28)
	timer_lbl.modulate = Color(1.0, 0.9, 0.2)
	panel.add_child(timer_lbl)
	var mult_lbl = Label.new()
	mult_lbl.name = "BossMultLabel"
	mult_lbl.text = "500x"
	mult_lbl.position = Vector2(200, 28)
	mult_lbl.add_theme_font_size_override("font_size", 28)
	mult_lbl.modulate = Color(1.0, 0.5, 0.0)
	panel.add_child(mult_lbl)
	var hint_lbl = Label.new()
	hint_lbl.text = "Kill faster = higher reward!"
	hint_lbl.position = Vector2(10, 60)
	hint_lbl.add_theme_font_size_override("font_size", 12)
	hint_lbl.modulate = Color(0.8, 0.8, 0.8)
	panel.add_child(hint_lbl)

func _update_boss_timer() -> void:
	if not is_instance_valid(_boss_timer_node):
		return
	var timer_lbl = _boss_timer_node.get_node_or_null("BossTimeLabel")
	var mult_lbl = _boss_timer_node.get_node_or_null("BossMultLabel")
	if is_instance_valid(timer_lbl):
		timer_lbl.text = "%.1fs" % max(0, _boss_time_left)
		if _boss_time_left <= 10:
			timer_lbl.modulate = Color(1.0, 0.3, 0.3)
		elif _boss_time_left <= 20:
			timer_lbl.modulate = Color(1.0, 0.7, 0.2)
		else:
			timer_lbl.modulate = Color(1.0, 0.9, 0.2)
	if is_instance_valid(mult_lbl):
		var mult = 100
		if _boss_time_left > 50: mult = 500
		elif _boss_time_left > 40: mult = 400
		elif _boss_time_left > 30: mult = 300
		elif _boss_time_left > 20: mult = 200
		elif _boss_time_left > 10: mult = 150
		mult_lbl.text = "%dx" % mult

func _stop_boss_timer() -> void:
	_boss_active = false
	if is_instance_valid(_boss_timer_node):
		var tween = create_tween()
		tween.tween_property(_boss_timer_node, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func():
			if is_instance_valid(_boss_timer_node):
				_boss_timer_node.queue_free()
				_boss_timer_node = null
		)

func _on_bonus_event(event_data: Dictionary) -> void:
	match event_data.get("event", ""):
		"start":
			bonus_overlay.visible = true
		"end", "result":
			bonus_overlay.visible = false

func _on_auto_pressed() -> void:
	NetworkManager.send_auto_toggle()

func _on_lock_pressed() -> void:
	NetworkManager.send_lock("")

func _on_bet_minus() -> void:
	NetworkManager.send_bet_change(max(1, GameManager.get_bet_level() - 1))

func _on_bet_plus() -> void:
	NetworkManager.send_bet_change(min(10, GameManager.get_bet_level() + 1))

# ---- 斷線提示 ----
var _disconnect_overlay: Control = null
var _is_disconnected: bool = false

func _create_disconnect_overlay() -> void:
	var overlay = Control.new()
	overlay.name = "DisconnectOverlay"
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.visible = false
	overlay.z_index = 100
	add_child(overlay)
	_disconnect_overlay = overlay
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.0, 0.0, 0.0, 0.7)
	overlay.add_child(bg)
	var msg = Label.new()
	msg.name = "DisconnectMsg"
	msg.text = "DISCONNECTED"
	msg.position = Vector2(540, 340)
	msg.add_theme_font_size_override("font_size", 24)
	msg.modulate = Color(1.0, 0.3, 0.3)
	overlay.add_child(msg)
	var reconnect = Label.new()
	reconnect.name = "ReconnectLabel"
	reconnect.text = "Reconnecting..."
	reconnect.position = Vector2(560, 380)
	reconnect.add_theme_font_size_override("font_size", 14)
	reconnect.modulate = Color(0.8, 0.8, 0.8)
	overlay.add_child(reconnect)

func _on_disconnected() -> void:
	_is_disconnected = true
	if is_instance_valid(_disconnect_overlay):
		_disconnect_overlay.visible = true

func _on_reconnected() -> void:
	_is_disconnected = false
	if is_instance_valid(_disconnect_overlay):
		var tween = create_tween()
		tween.tween_interval(1.0)
		tween.tween_property(_disconnect_overlay, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func():
			if is_instance_valid(_disconnect_overlay):
				_disconnect_overlay.visible = false
				_disconnect_overlay.modulate.a = 1.0
		)

# ---- 成就通知（簡化版）----
func _on_achievement_unlocked(achievement_data: Dictionary) -> void:
	var name_text = achievement_data.get("name", "Achievement")
	var popup = Label.new()
	popup.text = "🏆 " + name_text
	popup.position = Vector2(900, 600)
	popup.add_theme_font_size_override("font_size", 14)
	popup.modulate = Color(1.0, 0.85, 0.1)
	add_child(popup)
	var tween = create_tween()
	tween.tween_property(popup, "position:y", 560.0, 0.5)
	tween.tween_interval(1.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(popup):
			popup.queue_free()
	)

# ---- 連擊事件（簡化版）----
func _on_combo_event(combo_data: Dictionary) -> void:
	var combo = combo_data.get("current", 0)
	if combo < 3:
		return
	var popup = Label.new()
	popup.text = "🔥 x%d COMBO!" % combo
	popup.position = Vector2(580, 200)
	popup.add_theme_font_size_override("font_size", 20)
	popup.modulate = Color(1.0, 0.5, 0.1)
	add_child(popup)
	var tween = create_tween()
	tween.tween_property(popup, "scale", Vector2(1.3, 1.3), 0.1)
	tween.tween_property(popup, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(0.8)
	tween.tween_property(popup, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(popup):
			popup.queue_free()
	)

# ---- 排行榜（簡化版）----
func _on_leaderboard_updated(_entries: Array) -> void:
	pass  # 簡化版不顯示排行榜

# ---- 觀戰者（簡化版）----
func _on_spectator_joined(_data: Dictionary) -> void:
	pass

func _on_spectator_left(_data: Dictionary) -> void:
	pass

# ---- 每日獎勵（簡化版）----
func _on_daily_bonus_received(bonus_data: Dictionary) -> void:
	var amount = bonus_data.get("amount", 0)
	if amount <= 0:
		return
	var popup = Label.new()
	popup.text = "🎁 Daily Bonus: +%d" % amount
	popup.position = Vector2(540, 300)
	popup.add_theme_font_size_override("font_size", 18)
	popup.modulate = Color(0.3, 1.0, 0.5)
	add_child(popup)
	var tween = create_tween()
	tween.tween_property(popup, "position:y", 260.0, 0.5)
	tween.tween_interval(2.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(popup):
			popup.queue_free()
	)

# ---- 幸運永生 BOSS 魚（DAY-289）----
func _on_lucky_immortal_boss(data: Dictionary) -> void:
	if is_instance_valid(_lucky_immortal_boss_panel):
		_lucky_immortal_boss_panel.handle_event(data)

# ---- 幸運怒氣蓄積魚（DAY-290）----
func _on_lucky_wrath_charge(data: Dictionary) -> void:
	if is_instance_valid(_lucky_wrath_charge_panel):
		_lucky_wrath_charge_panel.handle_event(data)
