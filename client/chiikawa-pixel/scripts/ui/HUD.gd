## HUD.gd
## 主遊戲 UI（規格書 11章）

extends CanvasLayer

const PixelTheme = preload("res://scripts/ui/PixelTheme.gd")

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

# BOSS 計時器（規格書 28.3：顯示剩餘時間與對應倍率）
var _boss_time_left: float = 0.0
var _boss_active: bool = false
var _boss_timer_node: Control = null

# 像素字體（規格書美術要求）
var _pixel_font: Font = null
const PIXEL_FONT_PATH = "res://assets/fonts/pixel8.fnt"

func _ready() -> void:
	# 套用像素風格 Theme（讓所有按鈕和 UI 元素有一致的像素風格）
	var pixel_theme = PixelTheme.create()
	# 套用到 TopBar 和 BottomBar 的所有子節點
	var top_bar = get_node_or_null("TopBar")
	var bottom_bar = get_node_or_null("BottomBar")
	if is_instance_valid(top_bar):
		top_bar.theme = pixel_theme
		# TopBar 背景（深海藍半透明）
		var top_bg = ColorRect.new()
		top_bg.name = "PixelBG"
		top_bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		top_bg.color = Color(0.03, 0.06, 0.18, 0.88)
		top_bg.z_index = -1
		top_bar.add_child(top_bg)
		top_bar.move_child(top_bg, 0)
		# TopBar 底部邊框線（金色）
		var top_line = ColorRect.new()
		top_line.name = "BottomLine"
		top_line.size = Vector2(1280, 2)
		top_line.position = Vector2(0, 38)
		top_line.color = Color(0.90, 0.75, 0.20, 0.60)
		top_bar.add_child(top_line)
	if is_instance_valid(bottom_bar):
		bottom_bar.theme = pixel_theme
		# BottomBar 背景（深海藍半透明）
		var bot_bg = ColorRect.new()
		bot_bg.name = "PixelBG"
		bot_bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		bot_bg.color = Color(0.03, 0.06, 0.18, 0.88)
		bot_bg.z_index = -1
		bottom_bar.add_child(bot_bg)
		bottom_bar.move_child(bot_bg, 0)
		# BottomBar 頂部邊框線（金色）
		var bot_line = ColorRect.new()
		bot_line.name = "TopLine"
		bot_line.size = Vector2(1280, 2)
		bot_line.position = Vector2(0, 0)
		bot_line.color = Color(0.90, 0.75, 0.20, 0.60)
		bottom_bar.add_child(bot_line)

	# 載入像素字體
	if ResourceLoader.exists(PIXEL_FONT_PATH):
		_pixel_font = load(PIXEL_FONT_PATH)
		_apply_pixel_font()

	GameManager.player_updated.connect(_on_player_updated)
	GameManager.game_state_changed.connect(_on_game_state_changed)
	GameManager.reward_received.connect(_on_reward_received)
	GameManager.boss_event.connect(_on_boss_event)
	GameManager.bonus_event.connect(_on_bonus_event)
	GameManager.leaderboard_updated.connect(_on_leaderboard_updated)
	GameManager.achievement_unlocked.connect(_on_achievement_unlocked)

	# 斷線/重連提示
	NetworkManager.disconnected.connect(_on_disconnected)
	NetworkManager.connected.connect(_on_reconnected)

	auto_button.pressed.connect(_on_auto_pressed)
	lock_button.pressed.connect(_on_lock_pressed)
	bet_minus_button.pressed.connect(_on_bet_minus)
	bet_plus_button.pressed.connect(_on_bet_plus)
	boss_button.pressed.connect(NetworkManager.send_trigger_boss)
	bonus_button.pressed.connect(NetworkManager.send_trigger_bonus)

	# UI 按鈕點擊音效（規格書 audio-map.json：ui.click = weed_pull.wav）
	for btn in [auto_button, lock_button, bet_minus_button, bet_plus_button, boss_button, bonus_button]:
		if is_instance_valid(btn):
			btn.pressed.connect(func(): AudioManager.play_sfx(AudioManager.SFX.WEED_PULL))

	reward_popup.visible = false
	_reward_popup_base_y = reward_popup.position.y

	# WarningLabel 像素風格（大字、紅色、陰影）
	var warning_label = get_node_or_null("WarningOverlay/WarningLabel")
	if is_instance_valid(warning_label):
		warning_label.add_theme_font_size_override("font_size", 72)
		warning_label.add_theme_color_override("font_color", Color(1.0, 0.15, 0.15))
		warning_label.add_theme_color_override("font_shadow_color", Color(0.5, 0.0, 0.0, 0.8))
		warning_label.add_theme_constant_override("shadow_offset_x", 3)
		warning_label.add_theme_constant_override("shadow_offset_y", 3)
		if is_instance_valid(_pixel_font):
			warning_label.add_theme_font_override("font", _pixel_font)

	# StateLabel 像素風格（右上角狀態顯示）
	var state_lbl = get_node_or_null("StateLabel")
	if is_instance_valid(state_lbl):
		state_lbl.add_theme_font_size_override("font_size", 11)
		state_lbl.add_theme_color_override("font_color", Color(0.6, 0.8, 1.0, 0.7))

	_update_ui()
	_create_disconnect_overlay()
	_create_leaderboard_panel()
	_create_achievement_queue()
	_create_lobby_overlay()  # 大廳 UI（DAY-020）

## 套用像素字體到所有 Label
func _apply_pixel_font() -> void:
	if not is_instance_valid(_pixel_font):
		return
	var labels = [coins_label, bet_label, character_label, labor_label, reward_popup, state_label]
	for label in labels:
		if is_instance_valid(label):
			label.add_theme_font_override("font", _pixel_font)
	# 按鈕字體
	var buttons = [auto_button, lock_button, bet_minus_button, bet_plus_button, boss_button, bonus_button]
	for btn in buttons:
		if is_instance_valid(btn):
			btn.add_theme_font_override("font", _pixel_font)

var _last_labor_value: int = 0  # 追蹤上次勞動值，偵測升滿觸發

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

	# 偵測勞動值剛達到 100（觸發升級特效）
	if labor >= 100 and _last_labor_value < 100:
		var char_id = GameManager.player_data.get("character_id", "chiikawa")
		HitEffect.spawn_level_up(Vector2(640, 630), char_id)
		ScreenShake.add_trauma(0.3)
	_last_labor_value = labor

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
			_show_boss_incoming_preview()  # BOSS 血條預覽動畫
		"spawn":
			AudioManager.play_bgm(AudioManager.BGM.BOSS_ENTER)
			_hide_boss_incoming_preview()  # 隱藏預覽，顯示正式計時器
			_start_boss_timer()
		"phase_change":
			AudioManager.play_bgm(AudioManager.BGM.BOSS_RAGE)
		"kill", "timeout":
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			AudioManager.play_bgm(AudioManager.BGM.MAIN_GAME)
			_stop_boss_timer()
			# BOSS 擊殺慶祝特效
			if event_data.get("event", "") == "kill":
				HitEffect.spawn_big_win(Vector2(640, 360), 100.0)
				ScreenShake.add_trauma(0.7)

# ── BOSS 計時器 UI（規格書 28.3）──────────────────────────────

func _start_boss_timer() -> void:
	_boss_time_left = 60.0
	_boss_active = true

	# 建立 BOSS 計時器面板
	if is_instance_valid(_boss_timer_node):
		_boss_timer_node.queue_free()

	var panel = Control.new()
	panel.name = "BossTimerPanel"
	panel.set_anchors_and_offsets_preset(Control.PRESET_TOP_RIGHT)
	panel.position = Vector2(900, 50)
	panel.size = Vector2(360, 80)
	add_child(panel)
	_boss_timer_node = panel

	# 背景
	var bg = ColorRect.new()
	bg.size = Vector2(360, 80)
	bg.color = Color(0.1, 0.0, 0.0, 0.85)
	panel.add_child(bg)

	# BOSS 標題
	var title = Label.new()
	title.name = "BossTitle"
	title.text = "⚔ BOSS BATTLE"
	title.position = Vector2(10, 5)
	title.add_theme_font_size_override("font_size", 16)
	title.modulate = Color(1.0, 0.3, 0.3)
	if is_instance_valid(_pixel_font):
		title.add_theme_font_override("font", _pixel_font)
	panel.add_child(title)

	# 剩餘時間
	var timer_lbl = Label.new()
	timer_lbl.name = "BossTimeLabel"
	timer_lbl.text = "60.0s"
	timer_lbl.position = Vector2(10, 28)
	timer_lbl.add_theme_font_size_override("font_size", 28)
	timer_lbl.modulate = Color(1.0, 0.9, 0.2)
	if is_instance_valid(_pixel_font):
		timer_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(timer_lbl)

	# 倍率提示
	var mult_lbl = Label.new()
	mult_lbl.name = "BossMultLabel"
	mult_lbl.text = "500x"
	mult_lbl.position = Vector2(200, 28)
	mult_lbl.add_theme_font_size_override("font_size", 28)
	mult_lbl.modulate = Color(1.0, 0.5, 0.0)
	if is_instance_valid(_pixel_font):
		mult_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(mult_lbl)

	# 倍率說明
	var hint_lbl = Label.new()
	hint_lbl.name = "BossHintLabel"
	hint_lbl.text = "Kill faster = higher reward!"
	hint_lbl.position = Vector2(10, 60)
	hint_lbl.add_theme_font_size_override("font_size", 12)
	hint_lbl.modulate = Color(0.8, 0.8, 0.8)
	if is_instance_valid(_pixel_font):
		hint_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(hint_lbl)

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

# ── BOSS 進場預覽 UI（警告階段顯示 BOSS 血條從 0 填滿）──────────

var _boss_preview_node: Control = null

## 顯示 BOSS 進場預覽（警告階段 3 秒）
## 血條從 0 緩慢填滿到 100%，增加期待感
func _show_boss_incoming_preview() -> void:
	if is_instance_valid(_boss_preview_node):
		_boss_preview_node.queue_free()

	var panel = Control.new()
	panel.name = "BossIncomingPreview"
	panel.position = Vector2(320, 280)  # 畫面中央偏下
	panel.size = Vector2(640, 120)
	panel.z_index = 90
	panel.modulate.a = 0.0
	add_child(panel)
	_boss_preview_node = panel

	# 背景
	var bg = ColorRect.new()
	bg.size = Vector2(640, 120)
	bg.color = Color(0.05, 0.0, 0.0, 0.88)
	panel.add_child(bg)

	# 頂部紅色邊條（閃爍）
	var top_bar = ColorRect.new()
	top_bar.name = "TopBar"
	top_bar.size = Vector2(640, 4)
	top_bar.color = Color(1.0, 0.1, 0.1, 1.0)
	panel.add_child(top_bar)

	# BOSS 名稱
	var name_lbl = Label.new()
	name_lbl.name = "BossNameLabel"
	name_lbl.text = "那個孩子"
	name_lbl.position = Vector2(20, 12)
	name_lbl.add_theme_font_size_override("font_size", 22)
	name_lbl.modulate = Color(1.0, 0.3, 0.3)
	if is_instance_valid(_pixel_font):
		name_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(name_lbl)

	# BOSS 副標題
	var sub_lbl = Label.new()
	sub_lbl.text = "BOSS  HP: 3000"
	sub_lbl.position = Vector2(20, 40)
	sub_lbl.add_theme_font_size_override("font_size", 13)
	sub_lbl.modulate = Color(0.8, 0.8, 0.8)
	if is_instance_valid(_pixel_font):
		sub_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(sub_lbl)

	# HP 條背景
	var hp_bg = ColorRect.new()
	hp_bg.size = Vector2(600, 20)
	hp_bg.position = Vector2(20, 65)
	hp_bg.color = Color(0.15, 0.0, 0.0, 1.0)
	panel.add_child(hp_bg)

	# HP 條（從 0 填滿）
	var hp_bar = ColorRect.new()
	hp_bar.name = "BossHPBar"
	hp_bar.size = Vector2(0, 20)  # 初始寬度 0
	hp_bar.position = Vector2(20, 65)
	hp_bar.color = Color(0.9, 0.1, 0.1, 1.0)
	panel.add_child(hp_bar)

	# HP 條高光（頂部亮線）
	var hp_shine = ColorRect.new()
	hp_shine.name = "BossHPShine"
	hp_shine.size = Vector2(0, 4)
	hp_shine.position = Vector2(20, 65)
	hp_shine.color = Color(1.0, 0.5, 0.5, 0.6)
	panel.add_child(hp_shine)

	# 倍率提示
	var mult_lbl = Label.new()
	mult_lbl.text = "MAX 500x"
	mult_lbl.position = Vector2(490, 12)
	mult_lbl.add_theme_font_size_override("font_size", 18)
	mult_lbl.modulate = Color(1.0, 0.6, 0.0)
	if is_instance_valid(_pixel_font):
		mult_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(mult_lbl)

	# 倒數文字
	var countdown_lbl = Label.new()
	countdown_lbl.name = "CountdownLabel"
	countdown_lbl.text = "3"
	countdown_lbl.position = Vector2(295, 88)
	countdown_lbl.add_theme_font_size_override("font_size", 20)
	countdown_lbl.modulate = Color(1.0, 0.9, 0.2)
	countdown_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	countdown_lbl.size = Vector2(50, 28)
	if is_instance_valid(_pixel_font):
		countdown_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(countdown_lbl)

	# 動畫序列
	var tween = panel.create_tween()

	# 1. 淡入（0.2 秒）
	tween.tween_property(panel, "modulate:a", 1.0, 0.2)

	# 2. HP 條從 0 填滿（2.5 秒，模擬 BOSS 充能）
	tween.parallel().tween_property(hp_bar, "size:x", 600.0, 2.5).set_ease(Tween.EASE_IN).set_trans(Tween.TRANS_QUAD)
	tween.parallel().tween_property(hp_shine, "size:x", 600.0, 2.5).set_ease(Tween.EASE_IN).set_trans(Tween.TRANS_QUAD)

	# 3. 倒數 3→2→1
	tween.tween_callback(func():
		if is_instance_valid(countdown_lbl):
			countdown_lbl.text = "2"
			var t2 = countdown_lbl.create_tween()
			t2.tween_property(countdown_lbl, "scale", Vector2(1.5, 1.5), 0.1)
			t2.tween_property(countdown_lbl, "scale", Vector2(1.0, 1.0), 0.1)
	)
	tween.tween_interval(0.8)
	tween.tween_callback(func():
		if is_instance_valid(countdown_lbl):
			countdown_lbl.text = "1"
			countdown_lbl.modulate = Color(1.0, 0.4, 0.4)
			var t3 = countdown_lbl.create_tween()
			t3.tween_property(countdown_lbl, "scale", Vector2(1.8, 1.8), 0.1)
			t3.tween_property(countdown_lbl, "scale", Vector2(1.0, 1.0), 0.1)
	)
	tween.tween_interval(0.5)

	# 4. HP 條閃爍（充滿後閃爍 3 次）
	for _i in 3:
		tween.tween_property(hp_bar, "modulate", Color(2.0, 0.5, 0.5, 1.0), 0.06)
		tween.tween_property(hp_bar, "modulate", Color.WHITE, 0.06)

	# 頂部邊條閃爍動畫（獨立 tween，持續整個警告期間）
	var bar_tween = top_bar.create_tween().set_loops()
	bar_tween.tween_property(top_bar, "modulate:a", 0.3, 0.2)
	bar_tween.tween_property(top_bar, "modulate:a", 1.0, 0.2)

## 隱藏 BOSS 進場預覽（BOSS 正式出現時）
func _hide_boss_incoming_preview() -> void:
	if not is_instance_valid(_boss_preview_node):
		return
	var tween = create_tween()
	tween.tween_property(_boss_preview_node, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(_boss_preview_node):
			_boss_preview_node.queue_free()
			_boss_preview_node = null
	)

func _get_boss_multiplier_text(time_left: float) -> String:
	if time_left <= 10:
		return "100x"
	elif time_left <= 20:
		return "150x"
	elif time_left <= 30:
		return "200x"
	elif time_left <= 40:
		return "300x"
	elif time_left <= 50:
		return "400x"
	else:
		return "500x"

func _get_boss_multiplier_color(time_left: float) -> Color:
	if time_left <= 10:
		return Color(0.6, 0.6, 0.6)  # 灰（最低倍率）
	elif time_left <= 20:
		return Color(0.4, 0.8, 1.0)  # 藍
	elif time_left <= 30:
		return Color(0.4, 1.0, 0.4)  # 綠
	elif time_left <= 40:
		return Color(1.0, 0.9, 0.2)  # 黃
	elif time_left <= 50:
		return Color(1.0, 0.5, 0.0)  # 橙
	else:
		return Color(1.0, 0.2, 0.2)  # 紅（最高倍率）

func _process(delta: float) -> void:
	if not _boss_active or not is_instance_valid(_boss_timer_node):
		return

	_boss_time_left = max(0.0, _boss_time_left - delta)

	var timer_lbl = _boss_timer_node.get_node_or_null("BossTimeLabel")
	var mult_lbl = _boss_timer_node.get_node_or_null("BossMultLabel")

	if timer_lbl:
		timer_lbl.text = "%.1fs" % _boss_time_left
		# 最後 10 秒閃爍警告
		if _boss_time_left <= 10.0:
			var flash = int(_boss_time_left * 4) % 2 == 0
			timer_lbl.modulate = Color.RED if flash else Color.WHITE
		else:
			timer_lbl.modulate = Color(1.0, 0.9, 0.2)

	if mult_lbl:
		var mult_text = _get_boss_multiplier_text(_boss_time_left)
		var mult_color = _get_boss_multiplier_color(_boss_time_left)
		mult_lbl.text = mult_text
		mult_lbl.modulate = mult_color

		# 倍率變化時放大提示
		if mult_text != mult_lbl.get_meta("last_mult", ""):
			mult_lbl.set_meta("last_mult", mult_text)
			var tween = create_tween()
			tween.tween_property(mult_lbl, "scale", Vector2(1.4, 1.4), 0.1)
			tween.tween_property(mult_lbl, "scale", Vector2(1.0, 1.0), 0.1)

	# FPS 顯示（DEBUG 模式）
	if OS.is_debug_build():
		_update_fps_display()

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

# ── 斷線/重連 UI ──────────────────────────────────────────────

var _disconnect_overlay: Control = null
var _reconnect_dots_timer: float = 0.0
var _reconnect_dots: int = 0
var _is_disconnected: bool = false

func _create_disconnect_overlay() -> void:
	var overlay = Control.new()
	overlay.name = "DisconnectOverlay"
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.visible = false
	overlay.z_index = 100
	add_child(overlay)
	_disconnect_overlay = overlay

	# 半透明黑色背景
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.0, 0.0, 0.0, 0.75)
	overlay.add_child(bg)

	# 斷線圖示
	var icon_label = Label.new()
	icon_label.name = "DisconnectIcon"
	icon_label.text = "📡"
	icon_label.position = Vector2(580, 290)
	icon_label.add_theme_font_size_override("font_size", 48)
	overlay.add_child(icon_label)

	# 斷線文字
	var msg_label = Label.new()
	msg_label.name = "DisconnectMsg"
	msg_label.text = "連線中斷"
	msg_label.position = Vector2(540, 355)
	msg_label.add_theme_font_size_override("font_size", 24)
	msg_label.modulate = Color(1.0, 0.4, 0.4)
	if is_instance_valid(_pixel_font):
		msg_label.add_theme_font_override("font", _pixel_font)
	overlay.add_child(msg_label)

	# 重連中文字（帶動態點點）
	var reconnect_label = Label.new()
	reconnect_label.name = "ReconnectLabel"
	reconnect_label.text = "重新連線中..."
	reconnect_label.position = Vector2(520, 390)
	reconnect_label.add_theme_font_size_override("font_size", 18)
	reconnect_label.modulate = Color(0.8, 0.8, 0.8)
	if is_instance_valid(_pixel_font):
		reconnect_label.add_theme_font_override("font", _pixel_font)
	overlay.add_child(reconnect_label)

func _on_disconnected() -> void:
	_is_disconnected = true
	if is_instance_valid(_disconnect_overlay):
		_disconnect_overlay.visible = true
		# 閃爍動畫
		var tween = create_tween().set_loops()
		tween.tween_property(_disconnect_overlay, "modulate:a", 0.7, 0.5)
		tween.tween_property(_disconnect_overlay, "modulate:a", 1.0, 0.5)

func _on_reconnected() -> void:
	_is_disconnected = false
	if is_instance_valid(_disconnect_overlay):
		# 顯示「已重新連線」然後淡出
		var msg = _disconnect_overlay.get_node_or_null("DisconnectMsg")
		if msg:
			msg.text = "已重新連線 ✓"
			msg.modulate = Color(0.3, 1.0, 0.3)
		var reconnect = _disconnect_overlay.get_node_or_null("ReconnectLabel")
		if reconnect:
			reconnect.visible = false

		var tween = create_tween()
		tween.tween_interval(1.0)
		tween.tween_property(_disconnect_overlay, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func():
			if is_instance_valid(_disconnect_overlay):
				_disconnect_overlay.visible = false
				_disconnect_overlay.modulate.a = 1.0
				# 重置文字
				var m = _disconnect_overlay.get_node_or_null("DisconnectMsg")
				if m:
					m.text = "連線中斷"
					m.modulate = Color(1.0, 0.4, 0.4)
				var r = _disconnect_overlay.get_node_or_null("ReconnectLabel")
				if r:
					r.visible = true
		)

# ── 排行榜 UI ──────────────────────────────────────────────────

var _leaderboard_panel: Control = null
var _leaderboard_visible: bool = true
var _leaderboard_toggle_btn: Button = null
const MAX_LEADERBOARD_ENTRIES = 5

func _create_leaderboard_panel() -> void:
	# 排行榜容器（右上角，BOSS 計時器下方，避免重疊）
	# BOSS 計時器在 x=900, y=50，高度 80px → 排行榜從 y=140 開始
	var panel = Control.new()
	panel.name = "LeaderboardPanel"
	panel.position = Vector2(900, 140)
	panel.size = Vector2(360, 200)
	panel.z_index = 10
	add_child(panel)
	_leaderboard_panel = panel

	# 背景
	var bg = ColorRect.new()
	bg.name = "LeaderboardBG"
	bg.size = Vector2(360, 200)
	bg.color = Color(0.0, 0.05, 0.15, 0.82)
	panel.add_child(bg)

	# 標題列
	var title_bar = ColorRect.new()
	title_bar.size = Vector2(360, 28)
	title_bar.color = Color(0.05, 0.15, 0.4, 0.95)
	panel.add_child(title_bar)

	var title_lbl = Label.new()
	title_lbl.name = "LeaderboardTitle"
	title_lbl.text = "🏆 排行榜"
	title_lbl.position = Vector2(10, 4)
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.modulate = Color(1.0, 0.9, 0.3)
	if is_instance_valid(_pixel_font):
		title_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(title_lbl)

	# 折疊按鈕
	var toggle_btn = Button.new()
	toggle_btn.name = "LeaderboardToggle"
	toggle_btn.text = "▲"
	toggle_btn.position = Vector2(325, 2)
	toggle_btn.size = Vector2(30, 24)
	toggle_btn.add_theme_font_size_override("font_size", 12)
	toggle_btn.pressed.connect(_toggle_leaderboard)
	panel.add_child(toggle_btn)
	_leaderboard_toggle_btn = toggle_btn

	# 排行榜條目容器
	var entries_container = Control.new()
	entries_container.name = "EntriesContainer"
	entries_container.position = Vector2(0, 30)
	entries_container.size = Vector2(360, 170)
	panel.add_child(entries_container)

	# 預建 5 個條目（動態更新文字）
	for i in range(MAX_LEADERBOARD_ENTRIES):
		_create_leaderboard_row(entries_container, i)

	# 初始顯示「等待玩家...」
	_show_leaderboard_placeholder()

func _create_leaderboard_row(container: Control, index: int) -> void:
	var row = Control.new()
	row.name = "Row%d" % index
	row.position = Vector2(0, index * 32)
	row.size = Vector2(360, 30)
	container.add_child(row)

	# 交替背景色
	var row_bg = ColorRect.new()
	row_bg.name = "RowBG"
	row_bg.size = Vector2(360, 30)
	if index % 2 == 0:
		row_bg.color = Color(0.05, 0.1, 0.25, 0.6)
	else:
		row_bg.color = Color(0.03, 0.07, 0.18, 0.6)
	row.add_child(row_bg)

	# 名次標籤
	var rank_lbl = Label.new()
	rank_lbl.name = "RankLabel"
	rank_lbl.position = Vector2(6, 6)
	rank_lbl.size = Vector2(30, 20)
	rank_lbl.add_theme_font_size_override("font_size", 13)
	if is_instance_valid(_pixel_font):
		rank_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(rank_lbl)

	# 玩家名稱
	var name_lbl = Label.new()
	name_lbl.name = "NameLabel"
	name_lbl.position = Vector2(42, 6)
	name_lbl.size = Vector2(140, 20)
	name_lbl.add_theme_font_size_override("font_size", 12)
	name_lbl.clip_text = true
	if is_instance_valid(_pixel_font):
		name_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(name_lbl)

	# 分數
	var score_lbl = Label.new()
	score_lbl.name = "ScoreLabel"
	score_lbl.position = Vector2(188, 6)
	score_lbl.size = Vector2(100, 20)
	score_lbl.add_theme_font_size_override("font_size", 12)
	score_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	if is_instance_valid(_pixel_font):
		score_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(score_lbl)

	# 擊破數
	var kill_lbl = Label.new()
	kill_lbl.name = "KillLabel"
	kill_lbl.position = Vector2(295, 6)
	kill_lbl.size = Vector2(60, 20)
	kill_lbl.add_theme_font_size_override("font_size", 11)
	kill_lbl.modulate = Color(0.7, 0.9, 0.7)
	if is_instance_valid(_pixel_font):
		kill_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(kill_lbl)

	row.visible = false

func _show_leaderboard_placeholder() -> void:
	if not is_instance_valid(_leaderboard_panel):
		return
	var container = _leaderboard_panel.get_node_or_null("EntriesContainer")
	if not container:
		return

	# 顯示第一行作為佔位符
	var row = container.get_node_or_null("Row0")
	if row:
		row.visible = true
		var name_lbl = row.get_node_or_null("NameLabel")
		if name_lbl:
			name_lbl.text = "等待玩家加入..."
			name_lbl.modulate = Color(0.6, 0.6, 0.6)
		var rank_lbl = row.get_node_or_null("RankLabel")
		if rank_lbl:
			rank_lbl.text = ""
		var score_lbl = row.get_node_or_null("ScoreLabel")
		if score_lbl:
			score_lbl.text = ""
		var kill_lbl = row.get_node_or_null("KillLabel")
		if kill_lbl:
			kill_lbl.text = ""

func _on_leaderboard_updated(entries: Array) -> void:
	if not is_instance_valid(_leaderboard_panel):
		return
	var container = _leaderboard_panel.get_node_or_null("EntriesContainer")
	if not container:
		return

	var my_player_id = GameManager.get_player_id()
	var count = min(entries.size(), MAX_LEADERBOARD_ENTRIES)

	# 更新各條目
	for i in range(MAX_LEADERBOARD_ENTRIES):
		var row = container.get_node_or_null("Row%d" % i)
		if not row:
			continue

		if i >= count:
			row.visible = false
			continue

		row.visible = true
		var entry = entries[i]
		var is_self = entry.get("player_id", "") == my_player_id

		# 名次
		var rank_lbl = row.get_node_or_null("RankLabel")
		if rank_lbl:
			match i:
				0: rank_lbl.text = "🥇"
				1: rank_lbl.text = "🥈"
				2: rank_lbl.text = "🥉"
				_: rank_lbl.text = "#%d" % (i + 1)

		# 玩家名稱（自己高亮）
		var name_lbl = row.get_node_or_null("NameLabel")
		if name_lbl:
			var display = entry.get("display_name", "???")
			name_lbl.text = ("▶ " if is_self else "") + display
			name_lbl.modulate = Color(1.0, 1.0, 0.4) if is_self else Color.WHITE

		# 分數
		var score_lbl = row.get_node_or_null("ScoreLabel")
		if score_lbl:
			var score = entry.get("score", 0)
			score_lbl.text = "🪙%d" % score
			score_lbl.modulate = Color(1.0, 0.9, 0.3) if is_self else Color(0.9, 0.9, 0.9)

		# 擊破數
		var kill_lbl = row.get_node_or_null("KillLabel")
		if kill_lbl:
			kill_lbl.text = "×%d" % entry.get("kill_count", 0)

		# 自己的行加邊框高亮
		var row_bg = row.get_node_or_null("RowBG")
		if row_bg:
			if is_self:
				row_bg.color = Color(0.15, 0.25, 0.05, 0.8)
			elif i % 2 == 0:
				row_bg.color = Color(0.05, 0.1, 0.25, 0.6)
			else:
				row_bg.color = Color(0.03, 0.07, 0.18, 0.6)

	# 更新面板高度
	var new_height = 30 + count * 32
	if is_instance_valid(_leaderboard_panel):
		var bg = _leaderboard_panel.get_node_or_null("LeaderboardBG")
		if bg:
			bg.size.y = new_height
		_leaderboard_panel.size.y = new_height

func _toggle_leaderboard() -> void:
	if not is_instance_valid(_leaderboard_panel):
		return
	_leaderboard_visible = not _leaderboard_visible

	var container = _leaderboard_panel.get_node_or_null("EntriesContainer")
	if container:
		container.visible = _leaderboard_visible

	var bg = _leaderboard_panel.get_node_or_null("LeaderboardBG")
	if bg:
		bg.size.y = 200 if _leaderboard_visible else 28

	if is_instance_valid(_leaderboard_toggle_btn):
		_leaderboard_toggle_btn.text = "▲" if _leaderboard_visible else "▼"

# ── FPS 顯示（DEBUG 模式）──────────────────────────────────────

var _fps_label: Label = null
var _fps_update_timer: float = 0.0
var _perf_panel: Control = null  # 完整效能監控面板

func _update_fps_display() -> void:
	_fps_update_timer += get_process_delta_time()
	if _fps_update_timer < 0.5:
		return
	_fps_update_timer = 0.0

	# 首次建立完整效能面板
	if _perf_panel == null:
		_create_perf_panel()

	if not is_instance_valid(_perf_panel):
		return

	var fps = Engine.get_frames_per_second()
	var mem_mb = PerformanceMonitor.snapshot_memory_mb
	var draw_calls = PerformanceMonitor.snapshot_draw_calls
	var nodes = PerformanceMonitor.snapshot_nodes
	var quality = PerformanceMonitor.current_quality
	var quality_str = ["HIGH", "MED", "LOW"][quality]

	# FPS 行
	var fps_lbl = _perf_panel.get_node_or_null("FPSLine")
	if fps_lbl:
		fps_lbl.text = "FPS: %d  [%s]" % [fps, quality_str]
		if fps < 30:
			fps_lbl.modulate = Color(1.0, 0.3, 0.3, 0.95)
		elif fps < 50:
			fps_lbl.modulate = Color(1.0, 0.8, 0.2, 0.9)
		else:
			fps_lbl.modulate = Color(0.4, 1.0, 0.5, 0.85)

	# 記憶體行
	var mem_lbl = _perf_panel.get_node_or_null("MemLine")
	if mem_lbl:
		mem_lbl.text = "MEM: %.1f MB" % mem_mb
		if mem_mb > 200.0:
			mem_lbl.modulate = Color(1.0, 0.5, 0.2, 0.9)
		else:
			mem_lbl.modulate = Color(0.7, 0.9, 1.0, 0.8)

	# Draw Calls 行
	var dc_lbl = _perf_panel.get_node_or_null("DCLine")
	if dc_lbl:
		dc_lbl.text = "DC: %d  Nodes: %d" % [draw_calls, nodes]
		if draw_calls > 500:
			dc_lbl.modulate = Color(1.0, 0.6, 0.2, 0.9)
		else:
			dc_lbl.modulate = Color(0.7, 0.9, 1.0, 0.8)

func _create_perf_panel() -> void:
	var panel = Control.new()
	panel.name = "PerfPanel"
	panel.position = Vector2(8, 670)
	panel.size = Vector2(220, 56)
	panel.z_index = 200
	add_child(panel)
	_perf_panel = panel

	# 半透明背景
	var bg = ColorRect.new()
	bg.size = Vector2(220, 56)
	bg.color = Color(0.0, 0.0, 0.0, 0.55)
	panel.add_child(bg)

	# 左側綠色邊條（DEBUG 標識）
	var side = ColorRect.new()
	side.size = Vector2(3, 56)
	side.color = Color(0.2, 1.0, 0.4, 0.8)
	panel.add_child(side)

	# FPS 行
	var fps_lbl = Label.new()
	fps_lbl.name = "FPSLine"
	fps_lbl.position = Vector2(8, 4)
	fps_lbl.size = Vector2(210, 16)
	fps_lbl.add_theme_font_size_override("font_size", 12)
	fps_lbl.text = "FPS: --  [HIGH]"
	fps_lbl.modulate = Color(0.4, 1.0, 0.5, 0.85)
	if is_instance_valid(_pixel_font):
		fps_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(fps_lbl)

	# 記憶體行
	var mem_lbl = Label.new()
	mem_lbl.name = "MemLine"
	mem_lbl.position = Vector2(8, 22)
	mem_lbl.size = Vector2(210, 16)
	mem_lbl.add_theme_font_size_override("font_size", 12)
	mem_lbl.text = "MEM: -- MB"
	mem_lbl.modulate = Color(0.7, 0.9, 1.0, 0.8)
	if is_instance_valid(_pixel_font):
		mem_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(mem_lbl)

	# Draw Calls 行
	var dc_lbl = Label.new()
	dc_lbl.name = "DCLine"
	dc_lbl.position = Vector2(8, 40)
	dc_lbl.size = Vector2(210, 16)
	dc_lbl.add_theme_font_size_override("font_size", 12)
	dc_lbl.text = "DC: --  Nodes: --"
	dc_lbl.modulate = Color(0.7, 0.9, 1.0, 0.8)
	if is_instance_valid(_pixel_font):
		dc_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(dc_lbl)

# ── 成就通知系統 ──────────────────────────────────────────────

var _achievement_queue: Array = []   # 待顯示的成就佇列
var _achievement_showing: bool = false
var _achievement_panel: Control = null

func _create_achievement_queue() -> void:
	# 成就通知面板（右下角，初始隱藏）
	var panel = Control.new()
	panel.name = "AchievementPanel"
	panel.position = Vector2(900, 650)  # 右下角
	panel.size = Vector2(360, 80)
	panel.z_index = 50
	panel.visible = false
	add_child(panel)
	_achievement_panel = panel

	# 背景（深色半透明，金色邊框感）
	var bg = ColorRect.new()
	bg.name = "AchBG"
	bg.size = Vector2(360, 80)
	bg.color = Color(0.08, 0.06, 0.02, 0.92)
	panel.add_child(bg)

	# 金色頂部邊條
	var top_bar = ColorRect.new()
	top_bar.size = Vector2(360, 4)
	top_bar.color = Color(1.0, 0.85, 0.1, 1.0)
	panel.add_child(top_bar)

	# 成就圖示（大 emoji）
	var icon_lbl = Label.new()
	icon_lbl.name = "AchIcon"
	icon_lbl.text = "🏆"
	icon_lbl.position = Vector2(8, 18)
	icon_lbl.add_theme_font_size_override("font_size", 36)
	panel.add_child(icon_lbl)

	# 「成就解鎖！」標題
	var title_lbl = Label.new()
	title_lbl.name = "AchTitle"
	title_lbl.text = "成就解鎖！"
	title_lbl.position = Vector2(58, 8)
	title_lbl.add_theme_font_size_override("font_size", 11)
	title_lbl.modulate = Color(1.0, 0.85, 0.1)
	if is_instance_valid(_pixel_font):
		title_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(title_lbl)

	# 成就名稱
	var name_lbl = Label.new()
	name_lbl.name = "AchName"
	name_lbl.text = ""
	name_lbl.position = Vector2(58, 26)
	name_lbl.add_theme_font_size_override("font_size", 16)
	name_lbl.modulate = Color.WHITE
	if is_instance_valid(_pixel_font):
		name_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(name_lbl)

	# 成就描述
	var desc_lbl = Label.new()
	desc_lbl.name = "AchDesc"
	desc_lbl.text = ""
	desc_lbl.position = Vector2(58, 50)
	desc_lbl.add_theme_font_size_override("font_size", 11)
	desc_lbl.modulate = Color(0.8, 0.8, 0.8)
	if is_instance_valid(_pixel_font):
		desc_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(desc_lbl)

func _on_achievement_unlocked(achievement_data: Dictionary) -> void:
	_achievement_queue.append(achievement_data)
	if not _achievement_showing:
		_show_next_achievement()

func _show_next_achievement() -> void:
	if _achievement_queue.is_empty() or not is_instance_valid(_achievement_panel):
		_achievement_showing = false
		return

	_achievement_showing = true
	var data = _achievement_queue.pop_front()

	# 更新面板內容
	var icon_lbl = _achievement_panel.get_node_or_null("AchIcon")
	var name_lbl = _achievement_panel.get_node_or_null("AchName")
	var desc_lbl = _achievement_panel.get_node_or_null("AchDesc")

	if icon_lbl:
		icon_lbl.text = data.get("icon", "🏆")
	if name_lbl:
		name_lbl.text = data.get("name", "")
	if desc_lbl:
		desc_lbl.text = data.get("description", "")

	# 播放音效（用 bonus_ready 音效）
	AudioManager.play_sfx(AudioManager.SFX.BONUS_READY)

	# 動畫：從右側滑入 → 停留 3 秒 → 滑出
	_achievement_panel.modulate.a = 1.0
	_achievement_panel.position.x = 1300.0  # 畫面外右側
	_achievement_panel.visible = true

	var tween = create_tween()
	# 滑入（0.35 秒）
	tween.tween_property(_achievement_panel, "position:x", 900.0, 0.35).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	# 停留 3 秒
	tween.tween_interval(3.0)
	# 滑出（0.3 秒）
	tween.tween_property(_achievement_panel, "position:x", 1300.0, 0.3).set_ease(Tween.EASE_IN)
	tween.tween_callback(func():
		if is_instance_valid(_achievement_panel):
			_achievement_panel.visible = false
		# 顯示下一個成就（若佇列不空）
		_show_next_achievement()
	)

# ── 大廳 UI（DAY-020）──────────────────────────────────────────

var _lobby_overlay: Control = null
var _lobby_manager: Control = null

## 建立大廳 overlay（初始隱藏，可由「切換房間」按鈕呼叫）
func _create_lobby_overlay() -> void:
	# 建立全螢幕 overlay 容器
	var overlay = Control.new()
	overlay.name = "LobbyOverlay"
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.visible = false
	overlay.z_index = 150
	add_child(overlay)
	_lobby_overlay = overlay

	# 建立 LobbyManager UI
	var lobby_script = load("res://scripts/ui/LobbyManager.gd")
	if lobby_script:
		_lobby_manager = lobby_script.new()
		_lobby_manager.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		overlay.add_child(_lobby_manager)
		# 連接房間選擇訊號
		_lobby_manager.room_selected.connect(_on_lobby_room_selected)

	# 在 TopBar 加入「切換房間」按鈕
	var top_bar = get_node_or_null("TopBar")
	if is_instance_valid(top_bar):
		var switch_btn = Button.new()
		switch_btn.name = "SwitchRoomBtn"
		switch_btn.text = "🏠"
		switch_btn.position = Vector2(1240, 6)
		switch_btn.size = Vector2(32, 28)
		switch_btn.add_theme_font_size_override("font_size", 16)
		switch_btn.pressed.connect(_show_lobby)
		switch_btn.tooltip_text = "切換房間"
		if is_instance_valid(_pixel_font):
			switch_btn.add_theme_font_override("font", _pixel_font)
		top_bar.add_child(switch_btn)

## 顯示大廳
func _show_lobby() -> void:
	if not is_instance_valid(_lobby_overlay):
		return
	_lobby_overlay.visible = true
	_lobby_overlay.modulate.a = 0.0
	var tween = create_tween()
	tween.tween_property(_lobby_overlay, "modulate:a", 1.0, 0.3)
	# 觸發房間列表刷新
	if is_instance_valid(_lobby_manager) and _lobby_manager.has_method("show_lobby"):
		_lobby_manager.show_lobby()

## 大廳選擇房間後的回調
func _on_lobby_room_selected(room_id: String) -> void:
	print("[HUD] Room selected: ", room_id)
	# 淡出大廳
	var tween = create_tween()
	tween.tween_property(_lobby_overlay, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(_lobby_overlay):
			_lobby_overlay.visible = false
	)
	# 顯示切換房間提示
	var notify_lbl = Label.new()
	notify_lbl.text = "切換到 %s..." % room_id
	notify_lbl.position = Vector2(440, 360)
	notify_lbl.size = Vector2(400, 40)
	notify_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	notify_lbl.add_theme_font_size_override("font_size", 18)
	notify_lbl.modulate = Color(0.4, 1.0, 0.5)
	if is_instance_valid(_pixel_font):
		notify_lbl.add_theme_font_override("font", _pixel_font)
	add_child(notify_lbl)
	var t2 = create_tween()
	t2.tween_interval(1.5)
	t2.tween_property(notify_lbl, "modulate:a", 0.0, 0.5)
	t2.tween_callback(func():
		if is_instance_valid(notify_lbl):
			notify_lbl.queue_free()
	)
