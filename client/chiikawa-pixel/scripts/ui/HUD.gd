## HUD.gd — 核心 HUD
## hud-core-agent 負責維護
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
	_update_ui()

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
