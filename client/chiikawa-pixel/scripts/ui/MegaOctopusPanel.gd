## MegaOctopusPanel.gd
## 撌典?蝡?頧?Ｘ嚗AY-144嚗?
## 璆剔?靘?嚗ILI Mega Fishing?ega Octopus Wheel ??Defeat that giant octopus and enter
## the bonus wheel round where you have a chance to win massive guaranteed prizes up to 950x.??
## ? T108 敺孛?澆犖頧嚗??潛??蛛?50x-950x嚗??拙振暺??迫

extends Control

var pixel_font: Font = null

# 頧???
var _is_spinning: bool = false
var _spin_duration: float = 3.0
var _spin_elapsed: float = 0.0
var _current_angle: float = 0.0
var _spin_speed: float = 0.0  # 摨?蝘?
var _slots: Array = []
var _overlay: Control = null
var _wheel_container: Control = null
var _stop_btn: Button = null
var _result_shown: bool = false

# 頧憿
const COLOR_PURPLE = Color(0.58, 0.0, 0.83)
const COLOR_GOLD   = Color(1.0, 0.85, 0.0)
const COLOR_WHITE  = Color(1.0, 1.0, 1.0)

## ??????HUD.gd ?澆嚗?
func setup(font: Font) -> void:
	pixel_font = font
	GameManager.mega_octopus_wheel_start.connect(_on_wheel_start)
	GameManager.mega_octopus_wheel_result.connect(_on_wheel_result)

## 頧??嚗erver ?嚗?
func _on_wheel_start(data: Dictionary) -> void:
	if _is_spinning:
		return

	_slots = data.get("slots", [])
	_spin_duration = float(data.get("spin_duration", 3))
	_is_spinning = true
	_spin_elapsed = 0.0
	_result_shown = false
	_spin_speed = 720.0  # ???漲 720摨?蝘?2??蝘?

	_build_wheel_overlay()

## 撱箇?頧 overlay
func _build_wheel_overlay() -> void:
	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# ?刻撟蝵?
	_overlay = Control.new()
	_overlay.name = "MegaOctopusOverlay"
	_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_overlay.z_index = 90
	canvas_layer.add_child(_overlay)

	# ??暺?
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.0, 0.0, 0.0, 0.0)
	_overlay.add_child(bg)

	# 瘛∪?
	var bg_tween = bg.create_tween()
	bg_tween.tween_property(bg, "color:a", 0.75, 0.3)

	# 銝剖亢?Ｘ
	_wheel_container = Control.new()
	_wheel_container.name = "WheelContainer"
	_wheel_container.position = Vector2(390, 140)
	_wheel_container.size = Vector2(500, 440)
	_overlay.add_child(_wheel_container)

	# ?Ｘ?
	var panel_bg = ColorRect.new()
	panel_bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	panel_bg.color = Color(0.05, 0.02, 0.12, 0.95)
	_wheel_container.add_child(panel_bg)

	# 蝝怨??
	for border_data in [
		[Vector2(0, 0), Vector2(500, 3)],
		[Vector2(0, 437), Vector2(500, 3)],
		[Vector2(0, 0), Vector2(3, 440)],
		[Vector2(497, 0), Vector2(3, 440)],
	]:
		var border = ColorRect.new()
		border.position = border_data[0]
		border.size = border_data[1]
		border.color = COLOR_PURPLE
		_wheel_container.add_child(border)

	# 璅?
	var title = Label.new()
	title.text = "?? 撌典?蝡?頧嚗?"
	title.position = Vector2(0, 12)
	title.size = Vector2(500, 36)
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_font_size_override("font_size", 22)
	title.add_theme_color_override("font_color", COLOR_PURPLE)
	if is_instance_valid(pixel_font):
		title.add_theme_font_override("font", pixel_font)
	_wheel_container.add_child(title)

	# 頧嚗??澆?敶Ｘ???
	_build_wheel_slots()

	# ??嚗??其?閫耦嚗?
	var pointer = ColorRect.new()
	pointer.size = Vector2(16, 24)
	pointer.position = Vector2(242, 52)
	pointer.color = COLOR_GOLD
	_wheel_container.add_child(pointer)

	# ?閮?璅惜
	var timer_lbl = Label.new()
	timer_lbl.name = "TimerLabel"
	timer_lbl.text = "%.1f" % _spin_duration
	timer_lbl.position = Vector2(0, 340)
	timer_lbl.size = Vector2(500, 30)
	timer_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	timer_lbl.add_theme_font_size_override("font_size", 18)
	timer_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	if is_instance_valid(pixel_font):
		timer_lbl.add_theme_font_override("font", pixel_font)
	_wheel_container.add_child(timer_lbl)

	# ?迫??
	_stop_btn = Button.new()
	_stop_btn.name = "StopBtn"
	_stop_btn.text = "???迫嚗?"
	_stop_btn.position = Vector2(150, 380)
	_stop_btn.size = Vector2(200, 44)
	_stop_btn.add_theme_font_size_override("font_size", 16)
	_stop_btn.add_theme_color_override("font_color", COLOR_GOLD)
	if is_instance_valid(pixel_font):
		_stop_btn.add_theme_font_override("font", pixel_font)
	_stop_btn.pressed.connect(_on_stop_pressed)
	_wheel_container.add_child(_stop_btn)

	# 敶?
	_wheel_container.position.y = 600
	var panel_tween = _wheel_container.create_tween()
	panel_tween.tween_property(_wheel_container, "position:y", 140.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

## 撱箇?頧?澆?嚗??澆?敶Ｘ???
func _build_wheel_slots() -> void:
	if not is_instance_valid(_wheel_container):
		return

	var center = Vector2(250, 200)
	var radius = 130.0
	var slot_count = 8

	for i in slot_count:
		var angle = (float(i) / slot_count) * TAU - PI / 2.0
		var slot_x = center.x + cos(angle) * radius - 40
		var slot_y = center.y + sin(angle) * radius - 20

		var slot_data = {}
		if i < _slots.size():
			slot_data = _slots[i]

		var slot_bg = ColorRect.new()
		slot_bg.name = "Slot_%d" % i
		slot_bg.position = Vector2(slot_x, slot_y)
		slot_bg.size = Vector2(80, 40)
		var slot_color_hex = slot_data.get("color", "#C0C0C0")
		slot_bg.color = Color(slot_color_hex) * 0.3 + Color(0.05, 0.02, 0.12) * 0.7
		_wheel_container.add_child(slot_bg)

		var slot_lbl = Label.new()
		slot_lbl.text = slot_data.get("label", "%dx" % (i * 100 + 50))
		slot_lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		slot_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		slot_lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
		slot_lbl.add_theme_font_size_override("font_size", 11)
		slot_lbl.add_theme_color_override("font_color", Color(slot_color_hex))
		if is_instance_valid(pixel_font):
			slot_lbl.add_theme_font_override("font", pixel_font)
		slot_bg.add_child(slot_lbl)

## 瘥??湔嚗??斗?頧??恬?
func _process(delta: float) -> void:
	if not _is_spinning or not is_instance_valid(_wheel_container):
		return

	_spin_elapsed += delta

	# ?湔?閮?
	var remaining = max(0.0, _spin_duration - _spin_elapsed)
	var timer_lbl = _wheel_container.get_node_or_null("TimerLabel")
	if is_instance_valid(timer_lbl):
		timer_lbl.text = "%.1f 蝘??芸??迫" % remaining
		if remaining < 1.0:
			timer_lbl.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))

	# 頧??嚗摮?銝剖???嚗?
	_current_angle += _spin_speed * delta
	if _current_angle >= 360.0:
		_current_angle -= 360.0

	# ?湔?澆?雿蔭
	var center = Vector2(250, 200)
	var radius = 130.0
	var slot_count = 8
	for i in slot_count:
		var base_angle = (float(i) / slot_count) * TAU - PI / 2.0
		var current_rad = base_angle + deg_to_rad(_current_angle)
		var slot_x = center.x + cos(current_rad) * radius - 40
		var slot_y = center.y + sin(current_rad) * radius - 20

		var slot_node = _wheel_container.get_node_or_null("Slot_%d" % i)
		if is_instance_valid(slot_node):
			slot_node.position = Vector2(slot_x, slot_y)

## ?拙振暺??迫??
func _on_stop_pressed() -> void:
	if not _is_spinning:
		return
	if is_instance_valid(_stop_btn):
		_stop_btn.disabled = true
	# ?潮?甇Ｚ??策 Server
	NetworkManager.send_message("mega_octopus_wheel_stop", {})

## 頧蝯?嚗erver ?嚗?
func _on_wheel_result(data: Dictionary) -> void:
	if _result_shown:
		return
	_result_shown = true
	_is_spinning = false

	var result_index = data.get("result_index", 0)
	var multiplier = data.get("multiplier", 50)
	var reward = data.get("reward", 0)
	var slot_label = data.get("slot_label", "%dx" % multiplier)
	var slot_color_hex = data.get("slot_color", "#C0C0C0")

	# ?迫頧?
	if is_instance_valid(_stop_btn):
		_stop_btn.disabled = true

	# 擃漁蝯??澆?
	if is_instance_valid(_wheel_container):
		var slot_node = _wheel_container.get_node_or_null("Slot_%d" % result_index)
		if is_instance_valid(slot_node):
			var highlight_tween = slot_node.create_tween().set_loops(4)
			highlight_tween.tween_property(slot_node, "color", Color(slot_color_hex), 0.15)
			highlight_tween.tween_property(slot_node, "color", Color(slot_color_hex) * 0.3, 0.15)

	# ?剜?單?
	if AudioManager != null:
		AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)

	# ?刻撟?????00x 蝝怨嚗隞??莎?
	var flash_color = COLOR_PURPLE if multiplier >= 300 else COLOR_GOLD
	_spawn_result_flash(flash_color)

	# 憿舐內蝯?
	await get_tree().create_timer(0.8).timeout
	_show_result_popup(slot_label, reward, multiplier, slot_color_hex)

## ?刻撟???
func _spawn_result_flash(color: Color) -> void:
	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return
	var flash = ColorRect.new()
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, 0.0)
	flash.z_index = 91
	canvas_layer.add_child(flash)
	var tween = flash.create_tween()
	tween.tween_property(flash, "color:a", 0.5, 0.1)
	tween.tween_property(flash, "color:a", 0.0, 0.4)
	tween.tween_callback(flash.queue_free)

## 憿舐內蝯?敶?
func _show_result_popup(slot_label: String, reward: int, multiplier: int, color_hex: String) -> void:
	if not is_instance_valid(_wheel_container):
		return

	# 皜頧嚗＊蝷箇???
	for child in _wheel_container.get_children():
		if child.name != "WheelContainer":
			child.queue_free()

	# 蝯?璅?
	var result_title = Label.new()
	result_title.text = "?? %s" % slot_label
	result_title.position = Vector2(0, 80)
	result_title.size = Vector2(500, 80)
	result_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	result_title.add_theme_font_size_override("font_size", 52)
	result_title.add_theme_color_override("font_color", Color(color_hex))
	if is_instance_valid(pixel_font):
		result_title.add_theme_font_override("font", pixel_font)
	_wheel_container.add_child(result_title)

	# ???
	var reward_lbl = Label.new()
	reward_lbl.text = "?? +%d" % reward
	reward_lbl.position = Vector2(0, 180)
	reward_lbl.size = Vector2(500, 60)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_lbl.add_theme_font_size_override("font_size", 36)
	reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	if is_instance_valid(pixel_font):
		reward_lbl.add_theme_font_override("font", pixel_font)
	_wheel_container.add_child(reward_lbl)

	# 敶?
	result_title.modulate.a = 0.0
	result_title.position.y = 120
	var tween = result_title.create_tween()
	tween.tween_property(result_title, "position:y", 80.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(result_title, "modulate:a", 1.0, 0.3)

	reward_lbl.modulate.a = 0.0
	tween.tween_property(reward_lbl, "modulate:a", 1.0, 0.3)

	# 3 蝘???
	tween.tween_interval(3.0)
	tween.tween_callback(_close_overlay)

## ?? overlay
func _close_overlay() -> void:
	_is_spinning = false
	if is_instance_valid(_overlay):
		var tween = _overlay.create_tween()
		tween.tween_property(_overlay, "modulate:a", 0.0, 0.4)
		tween.tween_callback(_overlay.queue_free)
	_overlay = null
	_wheel_container = null
	_stop_btn = null
