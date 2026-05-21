## SpeedRacePanel.gd
## 全服競速獵殺面板（DAY-136）
## 業界依據：soup.io 2025「PvP modes where every shot counts」
## 競速開始時頂部橫幅滑入，倒數計時，第一名擊破時全服廣播

extends Control

# ---- 狀態 ----
var _is_active: bool = false
var _target_instance_id: String = ""
var _target_name: String = ""
var _target_mult: float = 0.0
var _seconds_left: float = 0.0
var _bonus_mult: float = 3.0

# ---- 節點 ----
var _banner: Control = null          # 頂部橫幅
var _banner_label: Label = null      # 橫幅文字
var _timer_label: Label = null       # 倒數計時
var _mult_label: Label = null        # 目標倍率
var _flash_overlay: ColorRect = null # 全螢幕閃光
var _result_popup: Control = null    # 個人結果彈窗
var _result_label: Label = null      # 結果文字

# ---- 顏色 ----
const COLOR_GOLD = Color(1.0, 0.85, 0.0, 1.0)
const COLOR_SILVER = Color(0.8, 0.8, 0.8, 1.0)
const COLOR_BRONZE = Color(0.8, 0.5, 0.2, 1.0)
const COLOR_BG = Color(0.0, 0.0, 0.0, 0.75)
const COLOR_RACE = Color(1.0, 0.7, 0.0, 1.0)  # 競速橙金色

func _ready() -> void:
	_build_ui()
	set_process(false)
	mouse_filter = Control.MOUSE_FILTER_IGNORE

func _build_ui() -> void:
	# 全螢幕閃光（競速開始/結束時）
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.85, 0.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅（競速進行中）
	_banner = Control.new()
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 56)
	_banner.position = Vector2(0, -60)  # 初始在畫面外
	_banner.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_banner)

	var banner_bg = ColorRect.new()
	banner_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	banner_bg.color = Color(0.15, 0.1, 0.0, 0.92)
	banner_bg.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_banner.add_child(banner_bg)

	# 左側競速圖示
	var icon_label = Label.new()
	icon_label.text = "🏆"
	icon_label.position = Vector2(12, 8)
	icon_label.add_theme_font_size_override("font_size", 28)
	_banner.add_child(icon_label)

	# 中央橫幅文字
	_banner_label = Label.new()
	_banner_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", COLOR_RACE)
	_banner_label.add_theme_font_size_override("font_size", 16)
	_banner_label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_banner.add_child(_banner_label)

	# 右側倒數計時
	_timer_label = Label.new()
	_timer_label.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_timer_label.position = Vector2(-80, -14)
	_timer_label.custom_minimum_size = Vector2(70, 28)
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	_timer_label.add_theme_color_override("font_color", Color.WHITE)
	_timer_label.add_theme_font_size_override("font_size", 18)
	_timer_label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_banner.add_child(_timer_label)

	# 個人結果彈窗（右側滑入）
	_result_popup = Control.new()
	_result_popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_popup.position = Vector2(300, -50)  # 初始在畫面外右側
	_result_popup.custom_minimum_size = Vector2(220, 100)
	_result_popup.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_result_popup)

	var result_bg = ColorRect.new()
	result_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	result_bg.color = COLOR_BG
	result_bg.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_result_popup.add_child(result_bg)

	_result_label = Label.new()
	_result_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_label.add_theme_color_override("font_color", COLOR_GOLD)
	_result_label.add_theme_font_size_override("font_size", 16)
	_result_label.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	_result_label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_result_popup.add_child(_result_label)

func _process(delta: float) -> void:
	if not _is_active:
		return
	_seconds_left -= delta
	if _seconds_left <= 0:
		_seconds_left = 0
		set_process(false)
		return
	_update_timer_display()

func _update_timer_display() -> void:
	if not is_instance_valid(_timer_label):
		return
	var secs = int(_seconds_left)
	_timer_label.text = "%ds" % secs
	# 最後 10 秒變紅色閃爍
	if _seconds_left <= 10:
		var blink = fmod(_seconds_left, 0.5) > 0.25
		_timer_label.add_theme_color_override("font_color", Color.RED if blink else Color.WHITE)
	else:
		_timer_label.add_theme_color_override("font_color", Color.WHITE)

# ---- 外部呼叫 ----

## 競速開始（全服廣播）
func on_speed_race_start(data: Dictionary) -> void:
	_target_instance_id = data.get("target_instance_id", "")
	_target_name = data.get("target_name", "目標")
	_target_mult = data.get("target_mult", 0.0)
	_seconds_left = data.get("seconds_left", 30.0)
	_bonus_mult = data.get("bonus_mult", 3.0)
	_is_active = true

	# 更新橫幅文字
	if is_instance_valid(_banner_label):
		_banner_label.text = "🏆 競速獵殺！搶先擊破【%s】(×%.0f) 獲得 %.0fx 獎勵！" % [_target_name, _target_mult, _bonus_mult]

	# 橫幅滑入
	_show_banner()
	# 全螢幕金色閃光
	_flash_screen(Color(1.0, 0.85, 0.0, 0.35), 0.3)
	# 開始倒數
	set_process(true)

## 競速結束（第一名擊破）
func on_speed_race_end(data: Dictionary) -> void:
	_is_active = false
	set_process(false)

	var winner_name = data.get("winner_name", "玩家")
	var target_name = data.get("target_name", _target_name)
	var bonus_mult = data.get("bonus_mult", _bonus_mult)

	# 更新橫幅文字
	if is_instance_valid(_banner_label):
		_banner_label.text = "🥇 %s 搶先擊破【%s】！獲得 %.1fx 獎勵！" % [winner_name, target_name, bonus_mult]
		_banner_label.add_theme_color_override("font_color", COLOR_GOLD)

	# 全螢幕金色閃光（更強烈）
	_flash_screen(Color(1.0, 0.85, 0.0, 0.5), 0.5)

	# 3 秒後橫幅滑出
	var tween = create_tween()
	tween.tween_interval(3.0)
	tween.tween_callback(_hide_banner)

## 競速取消（超時或目標消失）
func on_speed_race_cancel(data: Dictionary) -> void:
	_is_active = false
	set_process(false)

	var target_name = data.get("target_name", _target_name)
	if is_instance_valid(_banner_label):
		_banner_label.text = "⏰ 競速超時！【%s】無人搶先擊破。" % target_name
		_banner_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7, 1.0))

	# 2 秒後橫幅滑出
	var tween = create_tween()
	tween.tween_interval(2.0)
	tween.tween_callback(_hide_banner)

## 個人競速結果（名次 + 獎勵倍率）
func on_speed_race_result(data: Dictionary) -> void:
	var rank = data.get("rank", 0)
	var bonus_mult = data.get("bonus_mult", 1.0)
	var rank_icon = data.get("rank_icon", "")
	var message = data.get("message", "")

	# 依名次設定顏色
	var color = Color.WHITE
	match rank:
		1: color = COLOR_GOLD
		2: color = COLOR_SILVER
		3: color = COLOR_BRONZE

	if is_instance_valid(_result_label):
		_result_label.text = message
		_result_label.add_theme_color_override("font_color", color)

	# 第一名：更強烈的閃光
	if rank == 1:
		_flash_screen(Color(1.0, 0.9, 0.0, 0.6), 0.6)

	# 結果彈窗從右側滑入
	_show_result_popup()

	# 4 秒後滑出
	var tween = create_tween()
	tween.tween_interval(4.0)
	tween.tween_callback(_hide_result_popup)

# ---- 私有方法 ----

func _show_banner() -> void:
	if not is_instance_valid(_banner):
		return
	_banner.position = Vector2(0, -60)
	var tween = create_tween()
	tween.tween_property(_banner, "position", Vector2(0, 0), 0.3).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

func _hide_banner() -> void:
	if not is_instance_valid(_banner):
		return
	var tween = create_tween()
	tween.tween_property(_banner, "position", Vector2(0, -60), 0.25).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)

func _show_result_popup() -> void:
	if not is_instance_valid(_result_popup):
		return
	_result_popup.position = Vector2(300, -50)
	var tween = create_tween()
	tween.tween_property(_result_popup, "position", Vector2(0, -50), 0.3).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

func _hide_result_popup() -> void:
	if not is_instance_valid(_result_popup):
		return
	var tween = create_tween()
	tween.tween_property(_result_popup, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func(): 
		if is_instance_valid(_result_popup):
			_result_popup.modulate.a = 1.0
			_result_popup.position = Vector2(300, -50)
	)

func _flash_screen(color: Color, duration: float) -> void:
	if not is_instance_valid(_flash_overlay):
		return
	_flash_overlay.color = color
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)
