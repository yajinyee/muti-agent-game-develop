## WeaponPanel.gd ??甇血???Ｘ嚗AY-067嚗?
## 憿舐內銝郎?函?蝝??拙振暺???
## 雿蔭嚗ottomBar 撌血
extends Node2D

# ---- 撣豢 ----
const PANEL_WIDTH  := 200
const PANEL_HEIGHT := 80
const BTN_WIDTH    := 58
const BTN_HEIGHT   := 60

# 甇血摰儔嚗? Server 蝡?data/tables.go ?郊嚗?
const WEAPONS = [
	{
		"level": 1,
		"name": "璅???,"
		"icon": "?",
		"color": Color(0.9, 0.9, 0.9),
		"extra_cost": 0,
		"power_mod": 1.00,
		"desc": "璅?\n?⊿?憭祥??"
	},
	{
		"level": 2,
		"name": "撘瑕???,"
		"icon": "??,"
		"color": Color(0.0, 0.9, 1.0),
		"extra_cost": 50,
		"power_mod": 1.25,
		"desc": "+25%\n+50/??"
	},
	{
		"level": 3,
		"name": "頞???,"
		"icon": "??",
		"color": Color(1.0, 0.85, 0.0),
		"extra_cost": 150,
		"power_mod": 1.60,
		"desc": "+60%\n+150/??"
	}
]

# ---- 蝭暺???----
var _buttons: Array = []
var _pixel_font: Font = null
var _current_level: int = 1

# ---- 閮? ----
signal weapon_changed(level: int)

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
	bg.color = Color(0.05, 0.08, 0.18, 0.85)
	add_child(bg)

	# 璅?
	var title := Label.new()
	title.position = Vector2(4, 2)
	title.text = "甇血"
	title.add_theme_color_override("font_color", Color(0.7, 0.8, 1.0))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 10)
	bg.add_child(title)

	# 銝郎?冽???
	for i in range(3):
		var weapon = WEAPONS[i]
		var btn_x = 4 + i * (BTN_WIDTH + 4)

		var btn_bg := ColorRect.new()
		btn_bg.position = Vector2(btn_x, 14)
		btn_bg.size = Vector2(BTN_WIDTH, BTN_HEIGHT)
		btn_bg.color = Color(0.1, 0.15, 0.3, 0.9)
		bg.add_child(btn_bg)

		# 甇血?內
		var icon_label := Label.new()
		icon_label.position = Vector2(btn_x + 4, 16)
		icon_label.text = weapon["icon"]
		if _pixel_font:
			icon_label.add_theme_font_override("font", _pixel_font)
			icon_label.add_theme_font_size_override("font_size", 18)
		bg.add_child(icon_label)

		# 甇血隤芣?
		var desc_label := Label.new()
		desc_label.position = Vector2(btn_x + 2, 38)
		desc_label.text = weapon["desc"]
		desc_label.add_theme_color_override("font_color", weapon["color"])
		if _pixel_font:
			desc_label.add_theme_font_override("font", _pixel_font)
			desc_label.add_theme_font_size_override("font_size", 9)
		bg.add_child(desc_label)

		# 暺????Button嚗?
		var btn := Button.new()
		btn.position = Vector2(btn_x, 14)
		btn.size = Vector2(BTN_WIDTH, BTN_HEIGHT)
		btn.flat = true
		btn.text = ""
		btn.set_meta("weapon_level", weapon["level"])
		bg.add_child(btn)

		_buttons.append({
			"btn": btn,
			"bg": btn_bg,
			"icon": icon_label,
			"desc": desc_label,
			"level": weapon["level"]
		})

	_update_selection()

func _connect_signals() -> void:
	for item in _buttons:
		item["btn"].pressed.connect(_on_weapon_btn_pressed.bind(item["level"]))
	# ?? GameManager ??player_updated 閮?
	if GameManager.has_signal("player_updated"):
		GameManager.player_updated.connect(_on_player_updated)

# ---- 鈭辣?? ----
func _on_weapon_btn_pressed(level: int) -> void:
	if level == _current_level:
		return
	_current_level = level
	_update_selection()
	# ?潮郎?典?蝝?瘙 Server
	NetworkManager.send_message({
		"type": "upgrade_weapon",
		"payload": {"weapon_level": level}
	})
	emit_signal("weapon_changed", level)

func _on_player_updated(player_data: Dictionary) -> void:
	var weapon_level = player_data.get("weapon_level", 1)
	if weapon_level != _current_level:
		_current_level = weapon_level
		_update_selection()

# ---- UI ?湔 ----
func _update_selection() -> void:
	for item in _buttons:
		var is_selected = (item["level"] == _current_level)
		var weapon = WEAPONS[item["level"] - 1]

		if is_selected:
			# ?訾葉嚗漁?脤?獢?+ ?擃漁
			item["bg"].color = Color(0.15, 0.25, 0.5, 0.95)
			# ??獢?????StyleBoxFlat嚗?
			var style := StyleBoxFlat.new()
			style.bg_color = Color(0.15, 0.25, 0.5, 0.95)
			style.border_width_left = 2
			style.border_width_right = 2
			style.border_width_top = 2
			style.border_width_bottom = 2
			style.border_color = weapon["color"]
			item["btn"].add_theme_stylebox_override("normal", style)
			item["btn"].add_theme_stylebox_override("hover", style)
			item["btn"].add_theme_stylebox_override("pressed", style)
		else:
			# ?芷銝哨???
			item["bg"].color = Color(0.08, 0.12, 0.25, 0.85)
			var style := StyleBoxFlat.new()
			style.bg_color = Color(0.0, 0.0, 0.0, 0.0)
			item["btn"].add_theme_stylebox_override("normal", style)
			item["btn"].add_theme_stylebox_override("hover", style)
			item["btn"].add_theme_stylebox_override("pressed", style)
