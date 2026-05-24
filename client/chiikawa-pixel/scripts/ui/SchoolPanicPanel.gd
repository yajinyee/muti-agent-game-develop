## SchoolPanicPanel.gd ??擳黎撽???葆?Ｘ嚗AY-191嚗?
## 璆剔???嚗cean King 3 Plus?chool of Fish ??when one fish is caught, others scatter in panic??
## 閬死銝駁?嚗??脤???+ 擳黎???? + ?閮?

extends Control

const PANIC_COLOR  := Color(1.0, 0.55, 0.0)   # 璈嚗???
const WARN_COLOR   := Color(1.0, 0.8, 0.0)    # 暺嚗郎??

var _banner: Control = null
var _timer_label: Label = null
var _panic_timer: float = 0.0
var _is_active: bool = false

func _ready() -> void:
	set_process(false)
	if GameManager.has_signal("school_panic"):
		GameManager.school_panic.connect(_on_school_panic)

func _process(delta: float) -> void:
	if not _is_active:
		set_process(false)
		return
	_panic_timer -= delta
	if _panic_timer <= 0.0:
		_panic_timer = 0.0
		_is_active = false
		set_process(false)
		_hide_all()
		return
	# ?湔?閮?
	if _timer_label and is_instance_valid(_timer_label):
		_timer_label.text = "?? 撽?銝?%.1f 蝘? % _panic_timer"
		# ?敺?3 蝘?蝝??
		if _panic_timer <= 3.0:
			_timer_label.add_theme_color_override("font_color", Color.RED)
		else:
			_timer_label.add_theme_color_override("font_color", PANIC_COLOR)

func _on_school_panic(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	match phase:
		"panic_start":
			_show_panic_start(data)
		"panic_end":
			_hide_all()

func _show_panic_start(data: Dictionary) -> void:
	var target_count: int = data.get("target_count", 0)
	var duration: float = data.get("duration", 8.0)
	var killer_name: String = data.get("killer_name", "")

	_panic_timer = duration
	_is_active = true

	# 璈??嚗甈∴?
	_flash_screen(PANIC_COLOR, 0.55)
	var t1 = get_tree().create_timer(0.2)
	t1.timeout.connect(func(): _flash_screen(PANIC_COLOR, 0.4))

	# ?璈怠?
	_show_banner("?? 擳黎撽?嚗?d 璇? HP 皜?嚗? % target_count, PANIC_COLOR)"

	# 銝剖亢憭批??內嚗撌梯孛?潭?嚗?
	var my_id: String = GameManager.my_player_id if GameManager.has_method("get_my_player_id") else ""
	var killer_id: String = data.get("killer_id", "")
	if killer_id != "" and killer_id == my_id:
		_show_center_popup("?? 擳黎撽?閫貊嚗n敹恍??游蝷?嚗?, PANIC_COLOR)"

	# 摨?閮?璅惜
	_show_timer_label(duration)

	set_process(true)

func _flash_screen(color: Color, alpha: float) -> void:
	var overlay = ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, alpha)
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tween = create_tween()
	tween.tween_property(overlay, "modulate:a", 0.0, 0.3)
	tween.tween_callback(overlay.queue_free)

func _show_banner(text: String, color: Color) -> void:
	if _banner != null and is_instance_valid(_banner):
		_banner.queue_free()

	_banner = Control.new()
	_banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 52)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.color = Color(color.r * 0.25, color.g * 0.25, color.b * 0.1, 0.9)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_banner.add_child(bg)

	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", 20)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.add_child(label)

	# 皛?
	_banner.position.y = -52
	var tween = create_tween()
	tween.tween_property(_banner, "position:y", 0.0, 0.25).set_trans(Tween.TRANS_BACK)

func _show_center_popup(text: String, color: Color) -> void:
	var popup = Label.new()
	popup.text = text
	popup.add_theme_color_override("font_color", color)
	popup.add_theme_font_size_override("font_size", 28)
	popup.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	popup.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	add_child(popup)

	# 敶歲?
	popup.scale = Vector2(0.5, 0.5)
	popup.pivot_offset = popup.size / 2.0
	var tween = create_tween()
	tween.tween_property(popup, "scale", Vector2(1.0, 1.0), 0.3).set_trans(Tween.TRANS_BACK)

	# 2 蝘?瘛∪
	var timer = get_tree().create_timer(2.0)
	timer.timeout.connect(func():
		if is_instance_valid(popup):
			var t2 = create_tween()
			t2.tween_property(popup, "modulate:a", 0.0, 0.4)
			t2.tween_callback(popup.queue_free)
	)

func _show_timer_label(duration: float) -> void:
	if _timer_label != null and is_instance_valid(_timer_label):
		_timer_label.queue_free()

	_timer_label = Label.new()
	_timer_label.text = "?? 撽?銝?%.1f 蝘? % duration"
	_timer_label.add_theme_color_override("font_color", PANIC_COLOR)
	_timer_label.add_theme_font_size_override("font_size", 16)
	_timer_label.set_anchors_preset(Control.PRESET_BOTTOM_RIGHT)
	_timer_label.offset_left = -200
	_timer_label.offset_top = -40
	_timer_label.offset_right = -10
	_timer_label.offset_bottom = -10
	add_child(_timer_label)

func _hide_all() -> void:
	if _banner != null and is_instance_valid(_banner):
		var tween = create_tween()
		tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
		tween.tween_callback(_banner.queue_free)
		_banner = null
	if _timer_label != null and is_instance_valid(_timer_label):
		var tween2 = create_tween()
		tween2.tween_property(_timer_label, "modulate:a", 0.0, 0.4)
		tween2.tween_callback(_timer_label.queue_free)
		_timer_label = null
