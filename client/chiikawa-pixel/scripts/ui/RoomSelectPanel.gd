## RoomSelectPanel.gd
## ?輸???漲?豢??Ｘ嚗AY-091嚗?
## 4 ?摨佗???/銝剔?/擃?/VIP
## 璆剔?靘?嚗cean King 蝟餃?憭摨行? 2026 撟湔?擳?璅?

extends Control

# ?輸?鞈?
var _rooms: Array = []
var _current_room: String = "beginner"
var _pixel_font: Font = null

# UI 蝭暺?
var _panel_bg: ColorRect = null
var _title_label: Label = null
var _room_buttons: Array = []
var _close_btn: Button = null
var _room_container: VBoxContainer = null

# ???
var _notify_label: Label = null
var _notify_timer: float = 0.0

# ??漲憿撠?
const DIFF_COLORS = {
	"beginner":     Color(0.30, 0.75, 0.30),  # 蝬
	"intermediate": Color(0.13, 0.59, 0.95),  # ?
	"advanced":     Color(1.00, 0.60, 0.00),  # 璈
	"vip":          Color(0.61, 0.15, 0.69),  # 蝝怨
}

func setup(font: Font) -> void:
	_pixel_font = font
	_build_ui()
	_connect_signals()
	# ?身?梯?
	visible = false

func _build_ui() -> void:
	# ????桃蔗
	var overlay = ColorRect.new()
	overlay.name = "Overlay"
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.color = Color(0, 0, 0, 0.65)
	overlay.z_index = -1
	add_child(overlay)

	# 銝駁??
	_panel_bg = ColorRect.new()
	_panel_bg.name = "PanelBG"
	_panel_bg.size = Vector2(520, 480)
	_panel_bg.position = Vector2(380, 80)
	_panel_bg.color = Color(0.05, 0.08, 0.20, 0.97)
	add_child(_panel_bg)

	# ?Ｘ??
	var border = ColorRect.new()
	border.name = "Border"
	border.size = Vector2(524, 484)
	border.position = Vector2(378, 78)
	border.color = Color(0.90, 0.75, 0.20, 0.80)
	border.z_index = -1
	add_child(border)

	# 璅?
	_title_label = Label.new()
	_title_label.name = "TitleLabel"
	_title_label.text = "?? ?豢??輸?"
	_title_label.position = Vector2(380, 90)
	_title_label.size = Vector2(520, 36)
	_title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_title_label.add_theme_color_override("font_color", Color(1.0, 0.90, 0.20))
	_title_label.add_theme_font_size_override("font_size", 20)
	if is_instance_valid(_pixel_font):
		_title_label.add_theme_font_override("font", _pixel_font)
	add_child(_title_label)

	# ?舀?憿?
	var subtitle = Label.new()
	subtitle.text = "銝???漲?????萄???Jackpot 蝝舐??漲"
	subtitle.position = Vector2(380, 122)
	subtitle.size = Vector2(520, 24)
	subtitle.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	subtitle.add_theme_color_override("font_color", Color(0.7, 0.7, 0.8))
	subtitle.add_theme_font_size_override("font_size", 12)
	if is_instance_valid(_pixel_font):
		subtitle.add_theme_font_override("font", _pixel_font)
	add_child(subtitle)

	# ?輸???摰孵
	_room_container = VBoxContainer.new()
	_room_container.name = "RoomContainer"
	_room_container.position = Vector2(400, 155)
	_room_container.size = Vector2(480, 280)
	_room_container.add_theme_constant_override("separation", 8)
	add_child(_room_container)

	# ????
	_close_btn = Button.new()
	_close_btn.name = "CloseBtn"
	_close_btn.text = "????"
	_close_btn.position = Vector2(560, 520)
	_close_btn.size = Vector2(120, 32)
	_close_btn.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	_close_btn.add_theme_font_size_override("font_size", 14)
	if is_instance_valid(_pixel_font):
		_close_btn.add_theme_font_override("font", _pixel_font)
	add_child(_close_btn)

	# ???璅惜
	_notify_label = Label.new()
	_notify_label.name = "NotifyLabel"
	_notify_label.position = Vector2(380, 450)
	_notify_label.size = Vector2(520, 36)
	_notify_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_notify_label.add_theme_font_size_override("font_size", 16)
	_notify_label.visible = false
	if is_instance_valid(_pixel_font):
		_notify_label.add_theme_font_override("font", _pixel_font)
	add_child(_notify_label)

func _connect_signals() -> void:
	if is_instance_valid(_close_btn):
		_close_btn.pressed.connect(_on_close_pressed)

	# ?? GameManager 閮?
	if GameManager.has_signal("room_list_received"):
		GameManager.room_list_received.connect(_on_room_list_received)
	if GameManager.has_signal("room_switched"):
		GameManager.room_switched.connect(_on_room_switched)
	if GameManager.has_signal("room_error"):
		GameManager.room_error.connect(_on_room_error)

func show_panel() -> void:
	visible = true
	# 隢???唳??銵?
	if NetworkManager.has_method("send_get_room_list"):
		NetworkManager.send_get_room_list()

func _on_close_pressed() -> void:
	visible = false

func _on_room_list_received(data: Dictionary) -> void:
	_rooms = data.get("rooms", [])
	_current_room = data.get("current_room", "beginner")
	_rebuild_room_buttons()

func _rebuild_room_buttons() -> void:
	# 皜????
	for btn in _room_buttons:
		if is_instance_valid(btn):
			btn.queue_free()
	_room_buttons.clear()

	# 皜摰孵摮?暺?
	for child in _room_container.get_children():
		child.queue_free()

	# 撱箇?瘥????銵?
	for room_data in _rooms:
		var row = _build_room_row(room_data)
		_room_container.add_child(row)
		_room_buttons.append(row)

func _build_room_row(room_data: Dictionary) -> Control:
	var diff_id = room_data.get("id", "beginner")
	var name_str = room_data.get("name", "")
	var icon = room_data.get("icon", "??")
	var color_hex = room_data.get("color", "#4CAF50")
	var min_bet = room_data.get("min_bet_cost", 0)
	var max_bet = room_data.get("max_bet_cost", 0)
	var player_count = room_data.get("player_count", 0)
	var max_players = room_data.get("max_players", 16)
	var reward_mult = room_data.get("reward_mult", 1.0)
	var jackpot_mult = room_data.get("jackpot_mult", 1.0)
	var entry_fee = room_data.get("entry_fee", 0)
	var is_available = room_data.get("is_available", true)
	var is_current = room_data.get("is_current", false)
	var description = room_data.get("description", "")

	var diff_color = DIFF_COLORS.get(diff_id, Color.WHITE)

	# 銵捆??
	var row = Control.new()
	row.custom_minimum_size = Vector2(480, 64)
	row.size = Vector2(480, 64)

	# ?
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	if is_current:
		bg.color = Color(diff_color.r, diff_color.g, diff_color.b, 0.30)
	elif not is_available:
		bg.color = Color(0.1, 0.1, 0.1, 0.5)
	else:
		bg.color = Color(0.08, 0.10, 0.22, 0.80)
	row.add_child(bg)

	# 撌血??嚗摨阡??莎?
	var left_bar = ColorRect.new()
	left_bar.size = Vector2(4, 60)
	left_bar.position = Vector2(0, 0)
	left_bar.color = diff_color if is_available else Color(0.3, 0.3, 0.3)
	row.add_child(left_bar)

	# ?內 + ?迂
	var icon_label = Label.new()
	icon_label.text = icon + " " + name_str
	icon_label.position = Vector2(12, 6)
	icon_label.size = Vector2(160, 24)
	icon_label.add_theme_color_override("font_color", diff_color if is_available else Color(0.5, 0.5, 0.5))
	icon_label.add_theme_font_size_override("font_size", 16)
	if is_instance_valid(_pixel_font):
		icon_label.add_theme_font_override("font", _pixel_font)
	row.add_child(icon_label)

	# ?膩
	var desc_label = Label.new()
	desc_label.text = description
	desc_label.position = Vector2(12, 32)
	desc_label.size = Vector2(200, 20)
	desc_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.7))
	desc_label.add_theme_font_size_override("font_size", 10)
	if is_instance_valid(_pixel_font):
		desc_label.add_theme_font_override("font", _pixel_font)
	row.add_child(desc_label)

	# ??鞈?
	var mult_label = Label.new()
	var mult_text = "? ?%.1f  Jackpot ?%.1f" % [reward_mult, jackpot_mult]
	mult_label.text = mult_text
	mult_label.position = Vector2(220, 6)
	mult_label.size = Vector2(160, 24)
	mult_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.20) if is_available else Color(0.4, 0.4, 0.4))
	mult_label.add_theme_font_size_override("font_size", 12)
	if is_instance_valid(_pixel_font):
		mult_label.add_theme_font_override("font", _pixel_font)
	row.add_child(mult_label)

	# Bet 蝭?
	var bet_label = Label.new()
	bet_label.text = "Bet: %d ~ %d" % [min_bet, max_bet]
	bet_label.position = Vector2(220, 32)
	bet_label.size = Vector2(160, 20)
	bet_label.add_theme_color_override("font_color", Color(0.6, 0.8, 0.6) if is_available else Color(0.4, 0.4, 0.4))
	bet_label.add_theme_font_size_override("font_size", 11)
	if is_instance_valid(_pixel_font):
		bet_label.add_theme_font_override("font", _pixel_font)
	row.add_child(bet_label)

	# 鈭箸
	var count_label = Label.new()
	count_label.text = "? %d/%d" % [player_count, max_players]
	count_label.position = Vector2(390, 6)
	count_label.size = Vector2(80, 24)
	count_label.add_theme_color_override("font_color", Color(0.7, 0.9, 1.0) if is_available else Color(0.4, 0.4, 0.4))
	count_label.add_theme_font_size_override("font_size", 12)
	if is_instance_valid(_pixel_font):
		count_label.add_theme_font_override("font", _pixel_font)
	row.add_child(count_label)

	# ?脣鞎?
	if entry_fee > 0:
		var fee_label = Label.new()
		fee_label.text = "? %d" % entry_fee
		fee_label.position = Vector2(390, 32)
		fee_label.size = Vector2(80, 20)
		fee_label.add_theme_color_override("font_color", Color(1.0, 0.7, 0.2) if is_available else Color(0.5, 0.4, 0.2))
		fee_label.add_theme_font_size_override("font_size", 11)
		if is_instance_valid(_pixel_font):
			fee_label.add_theme_font_override("font", _pixel_font)
		row.add_child(fee_label)

	# ?脣/?嗅???
	var enter_btn = Button.new()
	enter_btn.position = Vector2(390, 32) if entry_fee == 0 else Vector2(390, 32)
	enter_btn.size = Vector2(80, 24)
	if is_current:
		enter_btn.text = "???嗅?"
		enter_btn.position = Vector2(390, 18)
		enter_btn.add_theme_color_override("font_color", Color(0.3, 1.0, 0.3))
		enter_btn.disabled = true
	elif not is_available:
		enter_btn.text = "撌脫遛"
		enter_btn.position = Vector2(390, 18)
		enter_btn.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))
		enter_btn.disabled = true
	else:
		enter_btn.text = "?脣 ??"
		enter_btn.position = Vector2(390, 18)
		enter_btn.add_theme_color_override("font_color", diff_color)
		enter_btn.pressed.connect(func(): _on_enter_room(diff_id))
	enter_btn.add_theme_font_size_override("font_size", 12)
	if is_instance_valid(_pixel_font):
		enter_btn.add_theme_font_override("font", _pixel_font)
	row.add_child(enter_btn)

	# ?嗅??輸?璅?
	if is_current:
		var current_badge = Label.new()
		current_badge.text = "? ?桀?"
		current_badge.position = Vector2(12, 32)
		current_badge.size = Vector2(80, 20)
		current_badge.add_theme_color_override("font_color", Color(0.3, 1.0, 0.3))
		current_badge.add_theme_font_size_override("font_size", 11)
		if is_instance_valid(_pixel_font):
			current_badge.add_theme_font_override("font", _pixel_font)
		row.add_child(current_badge)

	return row

func _on_enter_room(diff_id: String) -> void:
	if NetworkManager.has_method("send_switch_room"):
		NetworkManager.send_switch_room(diff_id)

func _on_room_switched(data: Dictionary) -> void:
	var room_name = data.get("room_name", "")
	var room_icon = data.get("room_icon", "??")
	var reward_mult = data.get("reward_mult", 1.0)
	var entry_fee = data.get("entry_fee", 0)

	# 憿舐內?????
	var msg = "%s %s 撌脤脣嚗????%.1f" % [room_icon, room_name, reward_mult]
	if entry_fee > 0:
		msg += "嚗脣鞎?%d嚗? % entry_fee"
	_show_notify(msg, Color(0.3, 1.0, 0.3))

	# ?隢??輸??”嚗?啁????
	if NetworkManager.has_method("send_get_room_list"):
		NetworkManager.send_get_room_list()

func _on_room_error(data: Dictionary) -> void:
	var message = data.get("message", "??憭望?")
	_show_notify("??" + message, Color(1.0, 0.3, 0.3))

func _show_notify(text: String, color: Color) -> void:
	if not is_instance_valid(_notify_label):
		return
	_notify_label.text = text
	_notify_label.add_theme_color_override("font_color", color)
	_notify_label.visible = true
	_notify_timer = 3.0

func _process(delta: float) -> void:
	if _notify_timer > 0:
		_notify_timer -= delta
		if _notify_timer <= 0:
			if is_instance_valid(_notify_label):
				_notify_label.visible = false
