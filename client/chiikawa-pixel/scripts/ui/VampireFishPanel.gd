## VampireFishPanel.gd — 吸血鬼魚累積倍率 UI 面板（DAY-182）
## 業界依據：JILI 2026「The explicit multiplier of vampires increases the more you fight,
## and there is a chance that you can enter the multiplier mode, up to X5」
## 顯示吸血鬼模式激活、倍率累積動畫、倒數計時、最終結果
extends CanvasLayer

# ---- 常數 ----
const PANEL_COLOR_BG      := Color(0.1, 0.0, 0.05, 0.92)
const PANEL_COLOR_RED     := Color(0.8, 0.0, 0.1, 1.0)    # 深紅（吸血鬼感）
const PANEL_COLOR_CRIMSON := Color(1.0, 0.1, 0.2, 1.0)    # 緋紅
const PANEL_COLOR_GOLD    := Color(1.0, 0.85, 0.0, 1.0)
const PANEL_COLOR_WHITE   := Color(1.0, 1.0, 1.0, 1.0)

# ---- 節點引用 ----
var _banner_container  : Control
var _banner_label      : Label
var _mult_display      : Label    # 大倍率顯示（中央）
var _timer_label       : Label
var _result_panel      : Control
var _result_label      : Label
var _flash_overlay     : ColorRect

# ---- 狀態 ----
var _vampire_active    : bool = false
var _time_remaining    : float = 0.0
var _current_mult      : float = 1.0
var _kill_count        : int = 0
var _is_my_mode        : bool = false  # 是否是自己的吸血鬼模式

func _ready() -> void:
	layer = 63
	_build_ui()
	hide()

func _build_ui() -> void:
	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(0.8, 0.0, 0.1, 0.0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner_container = PanelContainer.new()
	_banner_container.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner_container.offset_top = 8
	_banner_container.offset_bottom = 56
	_banner_container.offset_left = 80
	_banner_container.offset_right = -80
	var banner_style := StyleBoxFlat.new()
	banner_style.bg_color = Color(0.3, 0.0, 0.05, 0.92)
	banner_style.corner_radius_top_left = 8
	banner_style.corner_radius_top_right = 8
	banner_style.corner_radius_bottom_left = 8
	banner_style.corner_radius_bottom_right = 8
	_banner_container.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner_container)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_CRIMSON)
	_banner_label.add_theme_font_size_override("font_size", 20)
	_banner_container.add_child(_banner_label)

	# 大倍率顯示（中央，只有自己的模式才顯示）
	_mult_display = Label.new()
	_mult_display.set_anchors_preset(Control.PRESET_CENTER)
	_mult_display.offset_left = -100
	_mult_display.offset_right = 100
	_mult_display.offset_top = -40
	_mult_display.offset_bottom = 40
	_mult_display.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_mult_display.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_mult_display.add_theme_color_override("font_color", PANEL_COLOR_CRIMSON)
	_mult_display.add_theme_font_size_override("font_size", 48)
	_mult_display.hide()
	add_child(_mult_display)

	# 倒數計時器（右上角）
	_timer_label = Label.new()
	_timer_label.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	_timer_label.offset_top = 64
	_timer_label.offset_right = -16
	_timer_label.offset_left = -200
	_timer_label.offset_bottom = 96
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	_timer_label.add_theme_color_override("font_color", PANEL_COLOR_RED)
	_timer_label.add_theme_font_size_override("font_size", 18)
	_timer_label.hide()
	add_child(_timer_label)

	# 結果彈窗（右側滑入）
	_result_panel = PanelContainer.new()
	_result_panel.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_panel.offset_right = -16
	_result_panel.offset_left = -280
	_result_panel.offset_top = -80
	_result_panel.offset_bottom = 80
	var result_style := StyleBoxFlat.new()
	result_style.bg_color = PANEL_COLOR_BG
	result_style.corner_radius_top_left = 12
	result_style.corner_radius_top_right = 12
	result_style.corner_radius_bottom_left = 12
	result_style.corner_radius_bottom_right = 12
	result_style.border_width_left = 2
	result_style.border_width_right = 2
	result_style.border_width_top = 2
	result_style.border_width_bottom = 2
	result_style.border_color = PANEL_COLOR_CRIMSON
	_result_panel.add_theme_stylebox_override("panel", result_style)
	_result_panel.modulate.a = 0.0
	add_child(_result_panel)

	_result_label = Label.new()
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_label.add_theme_color_override("font_color", PANEL_COLOR_WHITE)
	_result_label.add_theme_font_size_override("font_size", 18)
	_result_panel.add_child(_result_label)

func _process(delta: float) -> void:
	if _vampire_active and _is_my_mode:
		_time_remaining -= delta
		if _time_remaining <= 0.0:
			_vampire_active = false
			_timer_label.hide()
			_mult_display.hide()
		else:
			_timer_label.text = "🧛 %.1f 秒" % _time_remaining
			# 倍率越高，顏色越紅
			var intensity := (_current_mult - 1.0) / (5.0 - 1.0)
			var color := Color(0.8 + intensity * 0.2, 0.1 - intensity * 0.1, 0.2 - intensity * 0.1, 1.0)
			_mult_display.add_theme_color_override("font_color", color)

## handle_vampire_fish — 處理吸血鬼魚訊息
func handle_vampire_fish(payload: Dictionary, my_player_id: String) -> void:
	var phase : String = payload.get("phase", "")
	var player_id : String = payload.get("player_id", "")
	_is_my_mode = (player_id == my_player_id)

	match phase:
		"vampire_start":
			_on_vampire_start(payload)
		"vampire_broadcast":
			_on_vampire_broadcast(payload)
		"mult_update":
			if _is_my_mode:
				_on_mult_update(payload)
		"vampire_end":
			if _is_my_mode:
				_on_vampire_end(payload)

## _on_vampire_start — 吸血鬼模式激活（個人）
func _on_vampire_start(payload: Dictionary) -> void:
	_vampire_active = true
	_time_remaining = float(payload.get("duration_sec", 15))
	_current_mult = payload.get("current_mult", 1.0)
	_kill_count = 0

	show()

	# 全螢幕深紅閃光
	_flash_screen(Color(0.8, 0.0, 0.1, 0.5), 0.4)

	# 橫幅
	_banner_label.text = "🧛 吸血鬼模式激活！越打越強！最高 ×5.0！"
	_banner_container.modulate.a = 0.0
	_banner_container.show()
	var tween := create_tween()
	tween.tween_property(_banner_container, "modulate:a", 1.0, 0.3)

	# 大倍率顯示
	_mult_display.text = "×%.1f" % _current_mult
	_mult_display.scale = Vector2(0.5, 0.5)
	_mult_display.show()
	var mult_tween := create_tween()
	mult_tween.tween_property(_mult_display, "scale", Vector2(1.0, 1.0), 0.4).set_trans(Tween.TRANS_BACK)

	# 倒數計時器
	_timer_label.show()

## _on_vampire_broadcast — 全服廣播（其他玩家看到）
func _on_vampire_broadcast(payload: Dictionary) -> void:
	if _is_my_mode:
		return  # 自己的模式已在 vampire_start 處理
	var player_name : String = payload.get("player_name", "")
	show()
	_banner_label.text = "🧛 %s 進入吸血鬼模式！" % player_name
	_banner_container.modulate.a = 0.0
	_banner_container.show()
	var tween := create_tween()
	tween.tween_property(_banner_container, "modulate:a", 1.0, 0.3)
	# 3 秒後淡出橫幅
	await get_tree().create_timer(3.0).timeout
	var fade := create_tween()
	fade.tween_property(_banner_container, "modulate:a", 0.0, 0.5)
	await fade.finished
	if not _vampire_active:
		hide()

## _on_mult_update — 倍率更新（個人）
func _on_mult_update(payload: Dictionary) -> void:
	_current_mult = payload.get("current_mult", _current_mult)
	_kill_count += 1

	# 更新大倍率顯示
	_mult_display.text = "×%.1f" % _current_mult

	# 倍率提升彈跳動畫
	var tween := create_tween()
	tween.tween_property(_mult_display, "scale", Vector2(1.3, 1.3), 0.08)
	tween.tween_property(_mult_display, "scale", Vector2(1.0, 1.0), 0.12)

	# 達到 3.0x 以上時閃光
	if _current_mult >= 3.0:
		_flash_screen(Color(0.8, 0.0, 0.1, 0.2), 0.15)

	# 達到 5.0x 時特殊效果
	if _current_mult >= 5.0:
		_flash_screen(Color(1.0, 0.0, 0.0, 0.5), 0.3)
		_mult_display.add_theme_color_override("font_color", PANEL_COLOR_GOLD)

## _on_vampire_end — 吸血鬼模式結束（個人）
func _on_vampire_end(payload: Dictionary) -> void:
	_vampire_active = false
	var final_mult : float = payload.get("current_mult", _current_mult)
	var kill_count : int = payload.get("kill_count", _kill_count)

	_timer_label.hide()
	_mult_display.hide()

	# 結果彈窗
	_result_label.text = "🧛 吸血鬼模式結束\n最終倍率：×%.1f\n擊破：%d 個" % [final_mult, kill_count]
	_result_panel.offset_left = 0
	_result_panel.modulate.a = 1.0
	var tween := create_tween()
	tween.tween_property(_result_panel, "offset_left", -280.0, 0.4).set_trans(Tween.TRANS_BACK)

	# 高倍率特效
	if final_mult >= 4.0:
		_flash_screen(Color(1.0, 0.0, 0.0, 0.5), 0.4)
	elif final_mult >= 3.0:
		_flash_screen(Color(0.8, 0.0, 0.1, 0.3), 0.3)

	# 3 秒後淡出
	await get_tree().create_timer(3.5).timeout
	var fade := create_tween()
	fade.tween_property(_result_panel, "modulate:a", 0.0, 0.5)
	fade.parallel().tween_property(_banner_container, "modulate:a", 0.0, 0.5)
	await fade.finished
	_banner_container.hide()
	hide()

## _flash_screen — 全螢幕閃光
func _flash_screen(color: Color, duration: float) -> void:
	_flash_overlay.color = color
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)
