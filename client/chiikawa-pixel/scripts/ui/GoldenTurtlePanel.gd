## GoldenTurtlePanel.gd ??暺?瘚琿????迫?Ｘ嚗AY-159嚗?
## 璆剔?靘?嚗cean King 蝟餃??ime Stop????
## ?暺?瘚琿?敺孛?澆?湔???甇?8 蝘???璅?怠?蝘餃?
## 閬死嚗??脫???甇Ｗ???+ ?刻撟??脣???+ ?閮???
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
		_countdown_label.text = "??%.1f" % _remaining

## ??暺?瘚琿????迫鈭辣
func _on_golden_turtle_time_stop(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var duration: float = data.get("duration_secs", 8.0)
	var killer_name: String = data.get("killer_name", "")

	if phase == "time_stop_start":
		_show_time_stop(killer_name, duration)
	elif phase == "time_stop_end":
		_hide_time_stop()

## 憿舐內???迫??
func _show_time_stop(killer_name: String, duration: float) -> void:
	_is_active = true
	_remaining = duration

	# ?刻撟??脣????桃蔗
	_overlay = ColorRect.new()
	_overlay.position = Vector2(-640, -360)  # ?詨??潮?蹂?蝵?
	_overlay.size = Vector2(1280, 720)
	_overlay.color = Color(1.0, 0.9, 0.0, 0.08)
	add_child(_overlay)

	# ?刻撟??脤????剜嚗?
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

	# ?璈怠?
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
	banner_lbl.text = "? %s 閫貊暺?瘚琿?嚗???甇?%.0f 蝘?" % [killer_name, duration]
	banner_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	banner_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.0))
	banner_lbl.add_theme_font_size_override("font_size", 16)
	if _pixel_font:
		banner_lbl.add_theme_font_override("font", _pixel_font)
	add_child(banner_lbl)

	# 璈怠?皛?
	banner_bg.position.y = -400
	banner_lbl.position.y = -396
	var slide_tween = banner_bg.create_tween()
	slide_tween.tween_property(banner_bg, "position:y", -360.0, 0.3)
	var slide_tween2 = banner_lbl.create_tween()
	slide_tween2.tween_property(banner_lbl, "position:y", -356.0, 0.3)

	# ?閮??剁??喃?閫?
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
	_countdown_label.text = "??%.1f" % duration
	_countdown_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_countdown_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.0))
	_countdown_label.add_theme_font_size_override("font_size", 14)
	if _pixel_font:
		_countdown_label.add_theme_font_override("font", _pixel_font)
	add_child(_countdown_label)

## ?梯????迫??
func _hide_time_stop() -> void:
	_is_active = false
	_remaining = 0.0
	_hide_countdown()

func _hide_countdown() -> void:
	# 皜?????甇?UI
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
