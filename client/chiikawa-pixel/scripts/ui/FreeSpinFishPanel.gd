## FreeSpinFishPanel.gd — 自由旋轉魚免費射擊面板（DAY-204）
## 業界依據：Galaxsys King of Ocean 2026
## 「Free Spin Fish, Captain Fish, and Money Fish trigger bonus rounds,
##  extra multipliers, and instant payouts.」
##
## 視覺設計：
##   - 青色旋轉主題（#00CED1 + #7FFFD4 + #FFD700 + #FFFFFF）
##   - free_spin_start：青色雙閃光 + 頂部橫幅「🌀 免費射擊！」+ 底部時間進度條
##   - free_spin_shot：小閃光 + 浮動獎勵文字 + 計數器更新 + 進度條更新
##   - free_spin_end：結算彈窗（擊破數/總獎勵/延長秒數）+ 4秒後淡出
##   - free_spin_broadcast：全服廣播橫幅（其他玩家看到）
extends CanvasLayer

var _panel: Control
var _banner: Label
var _kill_counter: Label
var _time_bar: ColorRect
var _time_bar_bg: ColorRect
var _result_popup: Control

var _max_duration: float = 20.0
var _kill_count: int = 0
var _is_active: bool = false

func _ready() -> void:
	layer = 41  # 低於其他面板，不遮擋重要 UI
	_build_ui()

func _build_ui() -> void:
	_panel = Control.new()
	_panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	_panel.visible = false
	_panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_panel)

	# 頂部橫幅
	_banner = Label.new()
	_banner.text = "🌀 免費射擊！"
	_banner.add_theme_font_size_override("font_size", 22)
	_banner.add_theme_color_override("font_color", Color("#00CED1"))
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.position = Vector2(0, 10)
	_banner.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.visible = false
	add_child(_banner)

	# 擊破計數器
	_kill_counter = Label.new()
	_kill_counter.text = "擊破: 0"
	_kill_counter.add_theme_font_size_override("font_size", 16)
	_kill_counter.add_theme_color_override("font_color", Color("#7FFFD4"))
	_kill_counter.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_kill_counter.position = Vector2(0, 38)
	_kill_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_kill_counter.visible = false
	add_child(_kill_counter)

	# 底部時間進度條背景
	_time_bar_bg = ColorRect.new()
	_time_bar_bg.size = Vector2(1280, 12)
	_time_bar_bg.position = Vector2(0, 708)
	_time_bar_bg.color = Color(0.1, 0.1, 0.1, 0.8)
	_time_bar_bg.visible = false
	add_child(_time_bar_bg)

	# 底部時間進度條
	_time_bar = ColorRect.new()
	_time_bar.size = Vector2(1280, 12)
	_time_bar.position = Vector2(0, 708)
	_time_bar.color = Color("#00CED1")
	_time_bar.visible = false
	add_child(_time_bar)

	# 結算彈窗
	_result_popup = Control.new()
	_result_popup.visible = false
	_result_popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_popup.size = Vector2(280, 200)
	_result_popup.position = Vector2(-300, -100)
	add_child(_result_popup)

	var popup_bg = ColorRect.new()
	popup_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_bg.color = Color(0.0, 0.08, 0.12, 0.92)
	_result_popup.add_child(popup_bg)

	var popup_label = Label.new()
	popup_label.name = "ResultLabel"
	popup_label.add_theme_font_size_override("font_size", 14)
	popup_label.add_theme_color_override("font_color", Color("#00CED1"))
	popup_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_popup.add_child(popup_label)

## 處理自由旋轉魚訊息
func handle_free_spin_fish(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"free_spin_start":
			_on_free_spin_start(payload)
		"free_spin_shot":
			_on_free_spin_shot(payload)
		"free_spin_end":
			_on_free_spin_end(payload)
		"free_spin_broadcast":
			_on_free_spin_broadcast(payload)

## 免費射擊開始（個人）
func _on_free_spin_start(payload: Dictionary) -> void:
	var duration = payload.get("duration", 10.0)
	_max_duration = payload.get("max_duration", 20.0)
	_kill_count = 0
	_is_active = true
	_panel.visible = true

	# 青色雙閃光
	_flash_screen(Color("#00CED1"), 0.7)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color("#7FFFD4"), 0.5)

	# 顯示橫幅和計數器
	_banner.text = "🌀 免費射擊！不扣費！"
	_banner.visible = true
	_kill_counter.text = "擊破: 0"
	_kill_counter.visible = true

	# 顯示時間進度條
	_time_bar_bg.visible = true
	_time_bar.visible = true
	_time_bar.size.x = 1280.0

	# 啟動進度條動畫
	_animate_time_bar(duration)

## 免費射擊一次（個人）
func _on_free_spin_shot(payload: Dictionary) -> void:
	if not _is_active:
		return

	var killed = payload.get("killed", false)
	var reward = payload.get("reward", 0)
	var kill_count = payload.get("kill_count", 0)
	var remaining = payload.get("remaining", 0.0)
	var target_x = payload.get("target_x", 640.0)
	var target_y = payload.get("target_y", 360.0)

	_kill_count = kill_count
	_kill_counter.text = "擊破: %d" % kill_count

	# 更新進度條
	var ratio = clampf(remaining / _max_duration, 0.0, 1.0)
	_time_bar.size.x = 1280.0 * ratio
	# 顏色漸變：青→黃→橙
	if ratio > 0.5:
		_time_bar.color = Color("#00CED1")
	elif ratio > 0.25:
		_time_bar.color = Color("#FFD700")
	else:
		_time_bar.color = Color("#FF8C00")

	if killed and reward > 0:
		# 小閃光
		_flash_screen(Color("#7FFFD4"), 0.25)
		# 浮動獎勵文字
		_spawn_float_text("🌀 +%d" % reward, Color("#FFD700"),
			Vector2(target_x, target_y))

		# 每 5 次擊破額外閃光
		if kill_count % 5 == 0 and kill_count > 0:
			_flash_screen(Color("#00CED1"), 0.4)
			_spawn_float_text("🌀 %d 連擊！" % kill_count, Color("#00CED1"),
				Vector2(640, 300))

## 免費射擊結束（個人）
func _on_free_spin_end(payload: Dictionary) -> void:
	_is_active = false
	var kill_count = payload.get("kill_count", 0)
	var total_reward = payload.get("total_reward", 0)
	var extend_sec = payload.get("extend_sec", 0.0)

	# 隱藏橫幅和計數器
	_banner.visible = false
	_kill_counter.visible = false
	_time_bar.visible = false
	_time_bar_bg.visible = false

	# 依擊破數決定閃光
	if kill_count >= 10:
		_flash_screen(Color("#FFD700"), 0.6)
	elif kill_count >= 5:
		_flash_screen(Color("#00CED1"), 0.4)

	# 顯示結算彈窗
	var result_label = _result_popup.get_node("ResultLabel")
	var extend_line = ""
	if extend_sec > 0:
		extend_line = "\n⏱ 延長: %.1f 秒" % extend_sec
	result_label.text = (
		"🌀 免費射擊結算\n\n"
		+ "擊破: %d 個\n" % kill_count
		+ extend_line
		+ "\n總獎勵: %d 金幣" % total_reward
	)
	_result_popup.visible = true
	_result_popup.modulate.a = 1.0
	_result_popup.position = Vector2(get_viewport().get_visible_rect().size.x + 10, -100)

	# 從右側滑入
	var tween = create_tween()
	tween.tween_property(_result_popup, "position:x",
		get_viewport().get_visible_rect().size.x - 300, 0.4)

	# 4 秒後淡出
	await get_tree().create_timer(4.0).timeout
	var fade_tween = create_tween()
	fade_tween.tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	await fade_tween.finished
	_result_popup.visible = false
	_panel.visible = false

## 全服廣播（其他玩家看到）
func _on_free_spin_broadcast(payload: Dictionary) -> void:
	var player_name = payload.get("player_name", "玩家")
	# 短暫顯示廣播橫幅
	var broadcast_label = Label.new()
	broadcast_label.text = "🌀 %s 觸發免費射擊！" % player_name
	broadcast_label.add_theme_font_size_override("font_size", 16)
	broadcast_label.add_theme_color_override("font_color", Color("#00CED1"))
	broadcast_label.set_anchors_preset(Control.PRESET_TOP_WIDE)
	broadcast_label.position = Vector2(0, 65)
	broadcast_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	add_child(broadcast_label)

	var tween = create_tween()
	tween.tween_interval(2.5)
	tween.tween_property(broadcast_label, "modulate:a", 0.0, 0.5)
	tween.tween_callback(broadcast_label.queue_free)

## 進度條動畫（隨時間縮短）
func _animate_time_bar(duration: float) -> void:
	var tween = create_tween()
	tween.tween_property(_time_bar, "size:x", 0.0, duration)

## 全螢幕閃光
func _flash_screen(color: Color, alpha: float) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, alpha)
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.25)
	tween.tween_callback(flash.queue_free)

## 浮動文字
func _spawn_float_text(text: String, color: Color, pos: Vector2 = Vector2(-1, -1)) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 18)
	label.add_theme_color_override("font_color", color)
	if pos.x < 0:
		label.position = Vector2(randf_range(300, 900), randf_range(200, 500))
	else:
		label.position = pos
	add_child(label)
	var tween = create_tween()
	tween.tween_property(label, "position:y", label.position.y - 60, 1.0)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 1.0)
	tween.tween_callback(label.queue_free)
