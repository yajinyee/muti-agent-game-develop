## GoldenTurtlePanel.gd — 黃金海龜時間停止面板（DAY-159）
## 業界依據：Ocean King 系列「Time Stop」機制
## 擊破黃金海龜後觸發全場時間停止 8 秒，所有目標物暫停移動
## 視覺：金色時鐘停止動畫 + 全螢幕金色光暈 + 倒數計時器
extends Node2D

var _pixel_font: Font = null
var _countdown_label: Label = null
var _overlay: ColorRect = null
var _is_active: bool = false
var _remaining: float = 0.0

func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

func _connect_signals() -> void:
	if GameManager.has_signal("golden_turtle_time_stop"):
		GameManager.golden_turtle_time_stop.connect(_on_golden_turtle_time_stop)

func _process(delta: float) -> void:
	if not _is_active:
		return
	_remaining -= delta
	if _remaining <= 0.0:
		_remaining = 0.0
		_is_active = false
		_hide_countdown()
		return
	if is_instance_valid(_countdown_label):
		_countdown_label.text = "⏱ %.1f" % _remaining

## 處理黃金海龜時間停止事件
func _on_golden_turtle_time_stop(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var duration: float = data.get("duration_secs", 8.0)
	var killer_name: String = data.get("killer_name", "")

	if phase == "time_stop_start":
		_show_time_stop(killer_name, duration)
	elif phase == "time_stop_end":
		_hide_time_stop()

## 顯示時間停止效果
func _show_time_stop(killer_name: String, duration: float) -> void:
	_is_active = true
	_remaining = duration

	# 全螢幕金色半透明遮罩
	_overlay = ColorRect.new()
	_overlay.position = Vector2(-640, -360)  # 相對於面板位置
	_overlay.size = Vector2(1280, 720)
	_overlay.color = Color(1.0, 0.9, 0.0, 0.08)
	add_child(_overlay)

	# 全螢幕金色閃光（短暫）
	var flash := ColorRect.new()
	flash.position = Vector2(-640, -360)
	flash.size = Vector2(1280, 720)
	flash.color = Color(1.0, 0.9, 0.0, 0.35)
	add_child(flash)
	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", 0.0, 0.4)
	flash_tween.tween_callback(func():
		if is_instance_valid(flash): flash.queue_free()
	)

	# 頂部橫幅
	var banner_bg := ColorRect.new()
	banner_bg.name = "TimeBanner"
	banner_bg.position = Vector2(-640, -360)
	banner_bg.size = Vector2(1280, 48)
	banner_bg.color = Color(0.15, 0.12, 0.0, 0.92)
	add_child(banner_bg)

	var banner_lbl := Label.new()
	banner_lbl.name = "TimeBannerLabel"
	banner_lbl.position = Vector2(-640, -356)
	banner_lbl.size = Vector2(1280, 40)
	banner_lbl.text = "🐢 %s 觸發黃金海龜！時間停止 %.0f 秒！" % [killer_name, duration]
	banner_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	banner_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.0))
	banner_lbl.add_theme_font_size_override("font_size", 16)
	if _pixel_font:
		banner_lbl.add_theme_font_override("font", _pixel_font)
	add_child(banner_lbl)

	# 橫幅滑入動畫
	banner_bg.position.y = -400
	banner_lbl.position.y = -396
	var slide_tween = banner_bg.create_tween()
	slide_tween.tween_property(banner_bg, "position:y", -360.0, 0.3)
	var slide_tween2 = banner_lbl.create_tween()
	slide_tween2.tween_property(banner_lbl, "position:y", -356.0, 0.3)

	# 倒數計時器（右上角）
	var countdown_bg := ColorRect.new()
	countdown_bg.name = "CountdownBG"
	countdown_bg.position = Vector2(540, -355)
	countdown_bg.size = Vector2(90, 36)
	countdown_bg.color = Color(0.1, 0.08, 0.0, 0.9)
	add_child(countdown_bg)

	_countdown_label = Label.new()
	_countdown_label.name = "CountdownLabel"
	_countdown_label.position = Vector2(540, -352)
	_countdown_label.size = Vector2(90, 30)
	_countdown_label.text = "⏱ %.1f" % duration
	_countdown_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_countdown_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.0))
	_countdown_label.add_theme_font_size_override("font_size", 14)
	if _pixel_font:
		_countdown_label.add_theme_font_override("font", _pixel_font)
	add_child(_countdown_label)

## 隱藏時間停止效果
func _hide_time_stop() -> void:
	_is_active = false
	_remaining = 0.0
	_hide_countdown()

func _hide_countdown() -> void:
	# 清除所有時間停止 UI
	for child_name in ["TimeBanner", "TimeBannerLabel", "CountdownBG", "CountdownLabel"]:
		var node = get_node_or_null(child_name)
		if is_instance_valid(node):
			var tween = node.create_tween()
			tween.tween_property(node, "modulate:a", 0.0, 0.3)
			tween.tween_callback(func():
				if is_instance_valid(node): node.queue_free()
			)

	if is_instance_valid(_overlay):
		var tween = _overlay.create_tween()
		tween.tween_property(_overlay, "color:a", 0.0, 0.5)
		tween.tween_callback(func():
			if is_instance_valid(_overlay): _overlay.queue_free()
		)
		_overlay = null

	_countdown_label = null
