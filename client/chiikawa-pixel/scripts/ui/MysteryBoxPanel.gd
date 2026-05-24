## MysteryBoxPanel.gd ??蟡?撖嗥拳?Ｘ嚗AY-090嚗?
## 憿舐內??撖嗥拳?賊?嚗???蝞梧?憿舐內?拳?????
## 璆剔?靘?嚗erdbot.com 2026-05-02 蝣箄??ystery rewards? 2026 撟湔??梢???璈
extends Node2D

# ---- 撣豢 ----
const PANEL_WIDTH  := 200
const PANEL_HEIGHT := 80

# 蝔?漲憿摰儔
const RARITY_COLORS = {
	"common":    Color(0.63, 0.63, 0.63),
	"rare":      Color(0.25, 0.41, 0.88),
	"epic":      Color(0.61, 0.35, 0.71),
	"legendary": Color(1.0, 0.84, 0.0),
}

const RARITY_NAMES = {
	"common":    "?桅?,"
	"rare":      "蝔??,"
	"epic":      "?脰帘",
	"legendary": "?唾牧",
}

# ---- ???----
var _inventory: Dictionary = {}  # rarity -> count
var _pixel_font: Font = null
var _box_labels: Dictionary = {}  # rarity -> Label
var _box_buttons: Dictionary = {} # rarity -> ColorRect

# ---- ????----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_build_ui()
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

func _build_ui() -> void:
	# ?
	var bg := ColorRect.new()
	bg.position = Vector2(0, 0)
	bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	bg.color = Color(0.05, 0.08, 0.18, 0.88)
	add_child(bg)

	# 璅?
	var title := Label.new()
	title.position = Vector2(4, 2)
	title.text = "蟡?撖嗥拳"
	title.add_theme_color_override("font_color", Color(1.0, 0.84, 0.0))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 10)
	bg.add_child(title)

	# ???漲??
	var rarities = ["common", "rare", "epic", "legendary"]
	var icons = ["?", "??", "?", "??"]
	for i in range(rarities.size()):
		var rarity = rarities[i]
		var btn_x = 4 + i * 48

		# ???
		var btn_bg := ColorRect.new()
		btn_bg.name = "BtnBG_%s" % rarity
		btn_bg.position = Vector2(btn_x, 16)
		btn_bg.size = Vector2(44, 58)
		btn_bg.color = Color(0.08, 0.08, 0.15, 0.85)
		bg.add_child(btn_bg)
		_box_buttons[rarity] = btn_bg

		# ?內
		var icon_lbl := Label.new()
		icon_lbl.position = Vector2(btn_x + 4, 18)
		icon_lbl.text = icons[i]
		icon_lbl.add_theme_font_size_override("font_size", 20)
		bg.add_child(icon_lbl)

		# 蝔?漲?迂
		var name_lbl := Label.new()
		name_lbl.position = Vector2(btn_x + 2, 42)
		name_lbl.size = Vector2(42, 12)
		name_lbl.text = RARITY_NAMES.get(rarity, rarity)
		name_lbl.add_theme_font_size_override("font_size", 8)
		name_lbl.add_theme_color_override("font_color", RARITY_COLORS.get(rarity, Color.WHITE))
		if _pixel_font:
			name_lbl.add_theme_font_override("font", _pixel_font)
		bg.add_child(name_lbl)

		# ?賊?璅惜
		var count_lbl := Label.new()
		count_lbl.name = "Count_%s" % rarity
		count_lbl.position = Vector2(btn_x + 28, 18)
		count_lbl.size = Vector2(16, 14)
		count_lbl.text = "0"
		count_lbl.add_theme_font_size_override("font_size", 10)
		count_lbl.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))
		count_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
		if _pixel_font:
			count_lbl.add_theme_font_override("font", _pixel_font)
		bg.add_child(count_lbl)
		_box_labels[rarity] = count_lbl

		# 暺????
		var area := Area2D.new()
		var col := CollisionShape2D.new()
		var shape := RectangleShape2D.new()
		shape.size = Vector2(44, 58)
		col.shape = shape
		col.position = Vector2(btn_x + 22, 16 + 29)
		area.add_child(col)
		add_child(area)

		var r = rarity
		area.input_event.connect(func(_viewport, event, _shape_idx):
			if event is InputEventMouseButton and event.pressed and event.button_index == MOUSE_BUTTON_LEFT:
				_on_box_btn_pressed(r)
		)

func _connect_signals() -> void:
	if GameManager.has_signal("mystery_box_updated"):
		GameManager.mystery_box_updated.connect(_on_mystery_box_updated)
	if GameManager.has_signal("mystery_box_dropped"):
		GameManager.mystery_box_dropped.connect(_on_mystery_box_dropped)
	if GameManager.has_signal("mystery_box_opened"):
		GameManager.mystery_box_opened.connect(_on_mystery_box_opened)

# ---- 鈭辣?? ----

func _on_box_btn_pressed(rarity: String) -> void:
	var count = _inventory.get(rarity, 0)
	if count <= 0:
		return
	# ?潮?蝞梯?瘙?
	NetworkManager.send_open_mystery_box(rarity)

func _on_mystery_box_updated(data: Dictionary) -> void:
	_inventory.clear()
	var inv = data.get("inventory", [])
	for entry in inv:
		_inventory[entry.get("rarity", "")] = entry.get("count", 0)
	_update_display()

func _on_mystery_box_dropped(data: Dictionary) -> void:
	# 憿舐內撖嗥拳??嚗?格??拐?蝵殷?
	var rarity = data.get("rarity", "common")
	var drop_x = data.get("drop_x", 640.0)
	var drop_y = data.get("drop_y", 360.0)
	_show_drop_animation(rarity, drop_x, drop_y)

func _on_mystery_box_opened(data: Dictionary) -> void:
	# 憿舐內?拳蝯?
	_show_open_result(data)

# ---- UI ?湔 ----

func _update_display() -> void:
	for rarity in _box_labels:
		var lbl = _box_labels[rarity]
		if not is_instance_valid(lbl):
			continue
		var count = _inventory.get(rarity, 0)
		lbl.text = str(count)
		if count > 0:
			lbl.add_theme_color_override("font_color", RARITY_COLORS.get(rarity, Color.WHITE))
			# ???鈭株絲
			var btn_bg = _box_buttons.get(rarity)
			if is_instance_valid(btn_bg):
				btn_bg.color = Color(0.12, 0.15, 0.3, 0.95)
		else:
			lbl.add_theme_color_override("font_color", Color(0.4, 0.4, 0.4))
			var btn_bg = _box_buttons.get(rarity)
			if is_instance_valid(btn_bg):
				btn_bg.color = Color(0.08, 0.08, 0.15, 0.85)

# ---- ? ----

func _show_drop_animation(rarity: String, drop_x: float, drop_y: float) -> void:
	# ?函璅雿蔭憿舐內撖嗥拳??
	var canvas = get_viewport().get_canvas_item()
	var lbl = Label.new()
	lbl.text = _get_rarity_icon(rarity)
	lbl.position = Vector2(drop_x - 16, drop_y - 16)
	lbl.add_theme_font_size_override("font_size", 28)
	lbl.modulate = RARITY_COLORS.get(rarity, Color.WHITE)
	get_parent().add_child(lbl)

	# 銝? + 瘛∪?
	var tween = lbl.create_tween()
	tween.set_parallel(true)
	tween.tween_property(lbl, "position:y", drop_y - 80, 1.0).set_ease(Tween.EASE_OUT)
	tween.tween_property(lbl, "modulate:a", 0.0, 1.0).set_delay(0.5)
	tween.tween_callback(lbl.queue_free).set_delay(1.0)

func _show_open_result(data: Dictionary) -> void:
	# 撱箇??拳蝯?敶?
	var canvas = CanvasLayer.new()
	canvas.layer = 80
	get_parent().add_child(canvas)

	var rarity = data.get("rarity", "common")
	var box_name = data.get("box_name", "撖嗥拳")
	var box_icon = data.get("box_icon", "?")
	var reward = data.get("reward", {})
	var reward_label = reward.get("label", "")
	var reward_icon = reward.get("icon", "??")
	var rarity_color = RARITY_COLORS.get(rarity, Color.WHITE)

	# ???
	var overlay = ColorRect.new()
	overlay.color = Color(0, 0, 0, 0.5)
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.mouse_filter = Control.MOUSE_FILTER_STOP
	canvas.add_child(overlay)

	# 銝駁??
	var panel = Control.new()
	panel.position = Vector2(440, 260)
	panel.size = Vector2(400, 220)
	canvas.add_child(panel)

	# ?
	var bg = ColorRect.new()
	bg.size = Vector2(400, 220)
	bg.color = Color(0.05, 0.05, 0.15, 0.97)
	panel.add_child(bg)

	# ?敶抵璇?
	var top_bar = ColorRect.new()
	top_bar.size = Vector2(400, 4)
	top_bar.color = rarity_color
	panel.add_child(top_bar)

	# 撖嗥拳?內嚗之嚗?
	var box_lbl = Label.new()
	box_lbl.text = box_icon
	box_lbl.position = Vector2(160, 20)
	box_lbl.add_theme_font_size_override("font_size", 48)
	panel.add_child(box_lbl)

	# 撖嗥拳?迂
	var name_lbl = Label.new()
	name_lbl.text = box_name
	name_lbl.position = Vector2(0, 80)
	name_lbl.size = Vector2(400, 24)
	name_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	name_lbl.add_theme_font_size_override("font_size", 16)
	name_lbl.modulate = rarity_color
	if is_instance_valid(_pixel_font):
		name_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(name_lbl)

	# ?嚗之摮?
	var reward_lbl = Label.new()
	reward_lbl.text = "%s %s" % [reward_icon, reward_label]
	reward_lbl.position = Vector2(0, 110)
	reward_lbl.size = Vector2(400, 36)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_lbl.add_theme_font_size_override("font_size", 26)
	reward_lbl.modulate = Color(1.0, 0.95, 0.4)
	if is_instance_valid(_pixel_font):
		reward_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(reward_lbl)

	# ?????內嚗???嚗?
	var pending_mult = data.get("pending_mult", 0.0)
	if pending_mult > 1.0:
		var mult_lbl = Label.new()
		mult_lbl.text = "??銝活?餅? ?%.1f嚗?0蝘??嚗? % pending_mult"
		mult_lbl.position = Vector2(0, 152)
		mult_lbl.size = Vector2(400, 20)
		mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		mult_lbl.add_theme_font_size_override("font_size", 12)
		mult_lbl.modulate = Color(1.0, 0.85, 0.0)
		if is_instance_valid(_pixel_font):
			mult_lbl.add_theme_font_override("font", _pixel_font)
		panel.add_child(mult_lbl)

	# 蝣箄???
	var btn = Button.new()
	btn.text = "憭芣?鈭?"
	btn.position = Vector2(150, 178)
	btn.size = Vector2(100, 32)
	btn.add_theme_font_size_override("font_size", 14)
	if is_instance_valid(_pixel_font):
		btn.add_theme_font_override("font", _pixel_font)
	panel.add_child(btn)

	# 敶?
	panel.scale = Vector2(0.5, 0.5)
	panel.modulate.a = 0.0
	var tween = panel.create_tween()
	tween.set_parallel(true)
	tween.tween_property(panel, "scale", Vector2(1.0, 1.0), 0.3).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.tween_property(panel, "modulate:a", 1.0, 0.2)

	# ?唾牧撖嗥拳嚗?恍???
	if rarity == "legendary":
		var flash = ColorRect.new()
		flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		flash.color = Color(1.0, 0.84, 0.0, 0.0)
		canvas.add_child(flash)
		var flash_tween = flash.create_tween()
		flash_tween.tween_property(flash, "color:a", 0.4, 0.15)
		flash_tween.tween_property(flash, "color:a", 0.0, 0.3)
		flash_tween.tween_callback(flash.queue_free)

	# ??
	btn.pressed.connect(func():
		var close_tween = panel.create_tween()
		close_tween.set_parallel(true)
		close_tween.tween_property(panel, "scale", Vector2(0.8, 0.8), 0.15)
		close_tween.tween_property(panel, "modulate:a", 0.0, 0.15)
		close_tween.tween_callback(canvas.queue_free).set_delay(0.15)
	)

	# 3 蝘??芸???
	get_tree().create_timer(3.0).timeout.connect(func():
		if is_instance_valid(canvas):
			canvas.queue_free()
	)

func _get_rarity_icon(rarity: String) -> String:
	match rarity:
		"common": return "?"
		"rare": return "??"
		"epic": return "?"
		"legendary": return "??"
	return "?"
