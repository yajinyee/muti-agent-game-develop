## HUD.gd
## 主遊戲 UI（規格書 11章）

extends CanvasLayer

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

func _ready() -> void:
	GameManager.player_updated.connect(_on_player_updated)
	GameManager.game_state_changed.connect(_on_game_state_changed)
	GameManager.reward_received.connect(_on_reward_received)
	GameManager.boss_event.connect(_on_boss_event)
	GameManager.bonus_event.connect(_on_bonus_event)

	auto_button.pressed.connect(_on_auto_pressed)
	lock_button.pressed.connect(_on_lock_pressed)
	bet_minus_button.pressed.connect(_on_bet_minus)
	bet_plus_button.pressed.connect(_on_bet_plus)
	boss_button.pressed.connect(NetworkManager.send_trigger_boss)
	bonus_button.pressed.connect(NetworkManager.send_trigger_bonus)

	reward_popup.visible = false
	_reward_popup_base_y = reward_popup.position.y
	_update_ui()

func _on_player_updated(_data: Dictionary) -> void:
	_update_ui()

func _update_ui() -> void:
	coins_label.text = "🪙 %d" % GameManager.get_coins()

	var lv = GameManager.get_bet_level()
	var cost = GameManager.get_bet_cost()
	bet_label.text = "BET LV%d  (%d/shot)" % [lv, cost]

	character_label.text = "▶ %s" % GameManager.get_character_name()
	character_label.modulate = GameManager.get_character_color()

	var labor = GameManager.get_labor_value()
	labor_bar.value = labor
	# 勞動值接近滿時變色
	if labor >= 80:
		labor_label.text = "⚡ %d/100" % labor
		labor_label.modulate = Color(1.0, 0.9, 0.2)
	else:
		labor_label.text = "💪 %d/100" % labor
		labor_label.modulate = Color.WHITE

	# Auto 按鈕
	if GameManager.is_auto():
		auto_button.modulate = Color(0.3, 1.0, 0.3)
		auto_button.text = "AUTO ON"
	else:
		auto_button.modulate = Color.WHITE
		auto_button.text = "AUTO"

	# Lock 按鈕
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
	# 依倍率決定顯示內容
	var icon = "🪙"
	if multiplier >= 100:
		icon = "💰"
		reward_popup.modulate = Color(1.0, 0.3, 0.1, 1.0)
	elif multiplier >= 20:
		icon = "💰"
		reward_popup.modulate = Color(1.0, 0.85, 0.0, 1.0)
	elif multiplier >= 10:
		icon = "🪙"
		reward_popup.modulate = Color(1.0, 1.0, 0.4, 1.0)
	else:
		icon = "🪙"
		reward_popup.modulate = Color(1.0, 1.0, 1.0, 1.0)

	reward_popup.text = "%s +%d  ×%.0f" % [icon, amount, multiplier]
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
		"phase_change":
			AudioManager.play_bgm(AudioManager.BGM.BOSS_RAGE)
		"kill", "timeout":
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			AudioManager.play_bgm(AudioManager.BGM.MAIN_GAME)

func _on_bonus_event(event_data: Dictionary) -> void:
	match event_data.get("event", ""):
		"ready":
			AudioManager.play_sfx(AudioManager.SFX.BONUS_READY)
		"start":
			AudioManager.play_bgm(AudioManager.BGM.BONUS_GAME)
		"tick":
			var tl = event_data.get("time_left", 0.0)
			var timer_lbl = bonus_overlay.get_node_or_null("TimerLabel")
			if timer_lbl:
				timer_lbl.text = "%.1f" % tl
		"end":
			_show_reward_popup(event_data.get("reward", 0), event_data.get("multiplier", 50.0))
			AudioManager.play_bgm(AudioManager.BGM.MAIN_GAME)

func _on_auto_pressed() -> void:
	NetworkManager.send_auto_toggle()

func _on_lock_pressed() -> void:
	# 解除鎖定
	NetworkManager.send_lock("")

func _on_bet_minus() -> void:
	NetworkManager.send_bet_change(max(1, GameManager.get_bet_level() - 1))

func _on_bet_plus() -> void:
	NetworkManager.send_bet_change(min(10, GameManager.get_bet_level() + 1))
