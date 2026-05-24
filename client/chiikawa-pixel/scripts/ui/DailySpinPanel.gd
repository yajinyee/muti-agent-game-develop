## DailySpinPanel.gd
## 瘥蝪賢頧?Ｘ嚗AY-092嚗?
## 瘥予?餃?臭誑?祥頧?甈∴???? 7 憭拙頧?蝝???
## 璆剔?靘?嚗Gaming 2026 ??梢???璈嚗??仿?????????35%+

extends Control

var _pixel_font: Font = null
var _slots: Array = []
var _is_super: bool = false
var _can_spin: bool = false
var _login_streak: int = 0
var _next_spin_at: int = 0

# UI 蝭暺?
var _panel_bg: ColorRect = null
var _title_label: Label = null
var _streak_label: Label = null
var _spin_btn: Button = null
var _close_btn: Button = null
var _slot_nodes: Array = []
var _result_label: Label = null
var _countdown_label: Label = null

# ????
var _is_spinning: bool = false
var _spin_angle: float = 0.0
var _spin_speed: float = 0.0
var _target_slot: int = -1
var _spin_timer: float = 0.0
var _result_timer: float = 0.0

# 頧??
const WHEEL_RADIUS = 120.0
const WHEEL_CENTER = Vector2(640, 320)

func setup(font: Font) -> void:
	_pixel_font = font
	_build_ui()
	_connect_signals()
	visible = false

func _build_ui() -> void:
	# ???
	var overlay = ColorRect.new()
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.color = Color(0, 0, 0, 0.70)
	overlay.z_index = -1
	add_child(overlay)

	# 銝駁??
	_panel_bg = ColorRect.new()
	_panel_bg.size = Vector2(560, 500)
	_panel_bg.position = Vector2(360, 60)
	_panel_bg.color = Color(0.04, 0.07, 0.18, 0.97)
	add_child(_panel_bg)

	# ??
	var border = ColorRect.new()
	border.size = Vector2(564, 504)
	border.position = Vector2(358, 58)
	border.color = Color(0.90, 0.75, 0.20, 0.80)
	border.z_index = -1
	add_child(border)

	# 璅?
	_title_label = Label.new()
	_title_label.text = "? 瘥蝪賢頧"
	_title_label.position = Vector2(360, 70)
	_title_label.size = Vector2(560, 36)
	_title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_title_label.add_theme_color_override("font_color", Color(1.0, 0.90, 0.20))
	_title_label.add_theme_font_size_override("font_size", 22)
	if is_instance_valid(_pixel_font):
		_title_label.add_theme_font_override("font", _pixel_font)
	add_child(_title_label)

	# ????餃憭拇
	_streak_label = Label.new()
	_streak_label.text = "????餃嚗? 憭?"
	_streak_label.position = Vector2(360, 108)
	_streak_label.size = Vector2(560, 24)
	_streak_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_streak_label.add_theme_color_override("font_color", Color(0.7, 0.9, 1.0))
	_streak_label.add_theme_font_size_override("font_size", 14)
	if is_instance_valid(_pixel_font):
		_streak_label.add_theme_font_override("font", _pixel_font)
	add_child(_streak_label)

	# 頧?澆?嚗????耦??嚗?
	for i in range(8):
		var slot_node = _create_slot_node(i)
		add_child(slot_node)
		_slot_nodes.append(slot_node)

	# 銝剖???
	var pointer = Label.new()
	pointer.text = "??"
	pointer.position = Vector2(WHEEL_CENTER.x - 10, WHEEL_CENTER.y - WHEEL_RADIUS - 30)
	pointer.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))
	pointer.add_theme_font_size_override("font_size", 24)
	if is_instance_valid(_pixel_font):
		pointer.add_theme_font_override("font", _pixel_font)
	add_child(pointer)

	# 頧???
	_spin_btn = Button.new()
	_spin_btn.text = "? 頧?嚗?"
	_spin_btn.position = Vector2(560, 430)
	_spin_btn.size = Vector2(160, 40)
	_spin_btn.add_theme_color_override("font_color", Color(1.0, 0.90, 0.20))
	_spin_btn.add_theme_font_size_override("font_size", 16)
	if is_instance_valid(_pixel_font):
		_spin_btn.add_theme_font_override("font", _pixel_font)
	add_child(_spin_btn)

	# ????
	_close_btn = Button.new()
	_close_btn.text = "??"
	_close_btn.position = Vector2(890, 68)
	_close_btn.size = Vector2(30, 30)
	_close_btn.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	_close_btn.add_theme_font_size_override("font_size", 14)
	if is_instance_valid(_pixel_font):
		_close_btn.add_theme_font_override("font", _pixel_font)
	add_child(_close_btn)

	# 蝯?璅惜
	_result_label = Label.new()
	_result_label.position = Vector2(360, 390)
	_result_label.size = Vector2(560, 36)
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_label.add_theme_font_size_override("font_size", 18)
	_result_label.visible = false
	if is_instance_valid(_pixel_font):
		_result_label.add_theme_font_override("font", _pixel_font)
	add_child(_result_label)

	# ?閮?璅惜
	_countdown_label = Label.new()
	_countdown_label.position = Vector2(360, 480)
	_countdown_label.size = Vector2(560, 24)
	_countdown_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_countdown_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.7))
	_countdown_label.add_theme_font_size_override("font_size", 12)
	_countdown_label.visible = false
	if is_instance_valid(_pixel_font):
		_countdown_label.add_theme_font_override("font", _pixel_font)
	add_child(_countdown_label)

func _create_slot_node(index: int) -> Control:
	var angle = (float(index) / 8.0) * TAU - PI / 2.0
	var pos = WHEEL_CENTER + Vector2(cos(angle), sin(angle)) * WHEEL_RADIUS

	var container = Control.new()
	container.name = "Slot%d" % index
	container.position = pos - Vector2(32, 32)
	container.size = Vector2(64, 64)

	# ???
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.08, 0.12, 0.28, 0.90)
	container.add_child(bg)

	# ?內
	var icon_label = Label.new()
	icon_label.name = "Icon"
	icon_label.text = "?"
	icon_label.position = Vector2(0, 8)
	icon_label.size = Vector2(64, 28)
	icon_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	icon_label.add_theme_font_size_override("font_size", 20)
	if is_instance_valid(_pixel_font):
		icon_label.add_theme_font_override("font", _pixel_font)
	container.add_child(icon_label)

	# 璅惜
	var text_label = Label.new()
	text_label.name = "Text"
	text_label.text = ""
	text_label.position = Vector2(0, 38)
	text_label.size = Vector2(64, 20)
	text_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	text_label.add_theme_font_size_override("font_size", 9)
	text_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.9))
	if is_instance_valid(_pixel_font):
		text_label.add_theme_font_override("font", _pixel_font)
	container.add_child(text_label)

	return container

func _connect_signals() -> void:
	if is_instance_valid(_spin_btn):
		_spin_btn.pressed.connect(_on_spin_pressed)
	if is_instance_valid(_close_btn):
		_close_btn.pressed.connect(_on_close_pressed)

	if GameManager.has_signal("daily_spin_state"):
		GameManager.daily_spin_state.connect(_on_daily_spin_state)
	if GameManager.has_signal("daily_spin_result"):
		GameManager.daily_spin_result.connect(_on_daily_spin_result)

func show_panel() -> void:
	visible = true
	if NetworkManager.has_method("send_get_daily_spin"):
		NetworkManager.send_get_daily_spin()

func _on_close_pressed() -> void:
	visible = false

func _on_spin_pressed() -> void:
	if not _can_spin or _is_spinning:
		return
	if NetworkManager.has_method("send_daily_spin"):
		NetworkManager.send_daily_spin()

func _on_daily_spin_state(data: Dictionary) -> void:
	_can_spin = data.get("can_spin", false)
	_is_super = data.get("is_super", false)
	_login_streak = data.get("login_streak", 0)
	_next_spin_at = data.get("next_spin_at", 0)

	# ?湔?澆?
	var slots_key = "super_slots" if _is_super else "normal_slots"
	_slots = data.get(slots_key, [])
	_update_slot_display()
	_update_ui_state()

func _update_slot_display() -> void:
	for i in range(min(_slots.size(), _slot_nodes.size())):
		var slot_data = _slots[i]
		var node = _slot_nodes[i]
		if not is_instance_valid(node):
			continue

		var icon_label = node.get_node_or_null("Icon")
		var text_label = node.get_node_or_null("Text")
		var bg = node.get_node_or_null("ColorRect")

		if is_instance_valid(icon_label):
			icon_label.text = slot_data.get("icon", "?")

		if is_instance_valid(text_label):
			var label = slot_data.get("label", "")
			# 蝮桃璅惜
			if len(label) > 8:
				label = label.substr(0, 8)
			text_label.text = label

		# 頞?頧?澆????脤?獢?
		if is_instance_valid(bg) and slot_data.get("is_super", false):
			bg.color = Color(0.15, 0.12, 0.05, 0.95)

func _update_ui_state() -> void:
	# ?湔璅?嚗?蝝??歹?
	if _is_super:
		_title_label.text = "?? 頞?頧嚗?? %d 憭抬?" % _login_streak
		_title_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.20))
	else:
		_title_label.text = "? 瘥蝪賢頧"
		_title_label.add_theme_color_override("font_color", Color(1.0, 0.90, 0.20))

	# ?湔???憭拇
	var streak_text = "????餃嚗?d 憭? % _login_streak"
	if _login_streak >= 7:
		streak_text += " ?? 頞?頧閫??嚗?"
	elif _login_streak >= 3:
		streak_text += " 嚗? %d 憭抵圾??蝝??歹?" % (7 - _login_streak)
	_streak_label.text = streak_text

	# ?湔?????
	if is_instance_valid(_spin_btn):
		if _can_spin:
			_spin_btn.text = "? 頧?嚗??祥嚗?"
			_spin_btn.disabled = false
			_spin_btn.add_theme_color_override("font_color", Color(1.0, 0.90, 0.20))
		else:
			_spin_btn.text = "???予??"
			_spin_btn.disabled = true
			_spin_btn.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))

	# ?閮?
	if not _can_spin and _next_spin_at > 0:
		_countdown_label.visible = true
		_update_countdown()
	else:
		_countdown_label.visible = false

func _update_countdown() -> void:
	if not is_instance_valid(_countdown_label):
		return
	var now_ms = Time.get_unix_time_from_system() * 1000
	var diff_ms = _next_spin_at - now_ms
	if diff_ms <= 0:
		_countdown_label.text = "?臭誑頧?嚗?"
		return
	var diff_sec = int(diff_ms / 1000)
	var hours = diff_sec / 3600
	var minutes = (diff_sec % 3600) / 60
	var seconds = diff_sec % 60
	_countdown_label.text = "銝活頧嚗?02d:%02d:%02d" % [hours, minutes, seconds]

func _on_daily_spin_result(data: Dictionary) -> void:
	var slot_index = data.get("slot_index", 0)
	var slot = data.get("slot", {})
	var is_super = data.get("is_super", false)
	var new_balance = data.get("new_balance", 0)

	# ?剜頧?
	_play_spin_animation(slot_index, slot, is_super, new_balance)

	# ?湔???
	_can_spin = false
	_update_ui_state()

func _play_spin_animation(target_slot: int, slot: Dictionary, is_super: bool, new_balance: int) -> void:
	_is_spinning = true
	_target_slot = target_slot
	_spin_timer = 2.5  # 頧??? 2.5 蝘?_process ?其?閮? progress嚗?
	_spin_angle = 0.0  # ?蔭??閫漲

	# 2.5 蝘??迫?銝阡＊蝷箇???
	var tween = create_tween()
	tween.tween_delay(2.5)
	tween.tween_callback(func():
		if not is_instance_valid(self):
			return
		_is_spinning = false
		_spin_timer = 0.0
		# ?迫敺??澆?甇訾??圈?甇Ｖ?蝵?
		_reset_slot_positions()
		_show_result(slot, is_super, new_balance)
		_highlight_slot(target_slot)
	)

func _reset_slot_positions() -> void:
	# ?蝯?敺??澆?甇訾??圈?甇Ｖ?蝵?
	for i in range(_slot_nodes.size()):
		var node = _slot_nodes[i]
		if not is_instance_valid(node):
			continue
		var angle_deg = (float(i) / 8.0) * 360.0 - 90.0
		var angle_rad = deg_to_rad(angle_deg)
		var pos = WHEEL_CENTER + Vector2(cos(angle_rad), sin(angle_rad)) * WHEEL_RADIUS
		node.position = pos - Vector2(32, 32)

func _highlight_slot(index: int) -> void:
	if index < 0 or index >= _slot_nodes.size():
		return
	var node = _slot_nodes[index]
	if not is_instance_valid(node):
		return
	var bg = node.get_node_or_null("ColorRect")
	if is_instance_valid(bg):
		var tween = create_tween().set_loops(3)
		tween.tween_property(bg, "color", Color(0.8, 0.7, 0.1, 0.95), 0.15)
		tween.tween_property(bg, "color", Color(0.08, 0.12, 0.28, 0.90), 0.15)

func _show_result(slot: Dictionary, is_super: bool, new_balance: int) -> void:
	if not is_instance_valid(_result_label):
		return
	var icon = slot.get("icon", "??")
	var label = slot.get("label", "?")
	var result_text = "%s ?脣?嚗?s嚗? % [icon, label]"
	if is_super:
		result_text = "?? 頞??嚗? + result_text"
	_result_label.text = result_text
	_result_label.add_theme_color_override("font_color",
		Color(1.0, 0.85, 0.20) if is_super else Color(0.3, 1.0, 0.3))
	_result_label.visible = true
	_result_timer = 4.0

func _process(delta: float) -> void:
	if not visible:
		return

	# 頧???嚗?閬箸???
	if _is_spinning and _spin_timer > 0:
		_spin_timer -= delta
		var elapsed = 2.5 - max(_spin_timer, 0.0)
		var progress = clamp(elapsed / 2.5, 0.0, 1.0)
		# ?翰敺?楨??ease out嚗?
		var speed = lerp(720.0, 30.0, ease(progress, 2.5))
		_spin_angle += speed * delta
		for i in range(_slot_nodes.size()):
			var node = _slot_nodes[i]
			if not is_instance_valid(node):
				continue
			var base_angle = (float(i) / 8.0) * 360.0 - 90.0
			var current_angle_deg = base_angle + _spin_angle
			var current_angle_rad = deg_to_rad(current_angle_deg)
			var pos = WHEEL_CENTER + Vector2(cos(current_angle_rad), sin(current_angle_rad)) * WHEEL_RADIUS
			node.position = pos - Vector2(32, 32)

	# 蝯?憿舐內閮?
	if _result_timer > 0:
		_result_timer -= delta
		if _result_timer <= 0 and is_instance_valid(_result_label):
			_result_label.visible = false

	# ?閮??湔嚗?蝘?
	if not _can_spin and _next_spin_at > 0:
		_update_countdown()
