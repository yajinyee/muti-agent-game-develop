## SkinPanel.gd ???脣憭??Ｘ嚗AY-071嚗?
## 憿舐內?航頃鞎?鋆???啣?閫
## 雿蔭嚗ottomBar ?喳嚗eaponPanel ??嚗?
extends Node2D

# ---- 撣豢 ----
const PANEL_WIDTH  := 240
const PANEL_HEIGHT := 80
const BTN_WIDTH    := 54
const BTN_HEIGHT   := 60

# 憭?摰儔嚗? Server 蝡?ws/protocol.go ?郊嚗?
const SKINS = [
	{
		"id": "default",
		"name": "璅?",
		"icon": "?",
		"price": 0,
		"cannon_color": Color(0.9, 0.9, 0.9),
		"bullet_color": Color(1.0, 1.0, 0.8),
		"desc": "?祥\n?身"
	},
	{
		"id": "golden",
		"name": "暺?",
		"icon": "??,"
		"price": 5000,
		"cannon_color": Color(1.0, 0.843, 0.0),
		"bullet_color": Color(1.0, 0.647, 0.0),
		"desc": "5000\n暺?"
	},
	{
		"id": "sakura",
		"name": "瑹餉",
		"icon": "?",
		"price": 8000,
		"cannon_color": Color(1.0, 0.714, 0.773),
		"bullet_color": Color(1.0, 0.412, 0.706),
		"desc": "8000\n??"
	},
	{
		"id": "rainbow",
		"name": "敶抵",
		"icon": "??",
		"price": 20000,
		"cannon_color": Color(1.0, 0.412, 0.706),
		"bullet_color": Color(0.0, 1.0, 1.0),
		"desc": "20000\n?唾牧"
	}
]

# ---- 蝭暺???----
var _buttons: Array = []
var _pixel_font: Font = null
var _equipped_skin: String = "default"
var _owned_skins: Array = ["default"]
var _player_coins: int = 0

# ---- 閮? ----
signal skin_buy_requested(skin_id: String)
signal skin_equip_requested(skin_id: String)

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
	bg.color = Color(0.08, 0.05, 0.18, 0.85)
	add_child(bg)

	# 璅?
	var title := Label.new()
	title.position = Vector2(4, 2)
	title.text = "憭?"
	title.add_theme_color_override("font_color", Color(1.0, 0.8, 1.0))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 10)
	bg.add_child(title)

	# ??閫??
	for i in range(4):
		var skin = SKINS[i]
		var btn_x = 4 + i * (BTN_WIDTH + 2)

		var btn_bg := ColorRect.new()
		btn_bg.position = Vector2(btn_x, 14)
		btn_bg.size = Vector2(BTN_WIDTH, BTN_HEIGHT)
		btn_bg.color = Color(0.1, 0.08, 0.25, 0.9)
		btn_bg.name = "BtnBG_" + skin["id"]
		bg.add_child(btn_bg)

		# 憭??內
		var icon_label := Label.new()
		icon_label.position = Vector2(btn_x + 4, 16)
		icon_label.text = skin["icon"]
		if _pixel_font:
			icon_label.add_theme_font_override("font", _pixel_font)
			icon_label.add_theme_font_size_override("font_size", 18)
		bg.add_child(icon_label)

		# 憭?隤芣?
		var desc_label := Label.new()
		desc_label.position = Vector2(btn_x + 2, 38)
		desc_label.text = skin["desc"]
		desc_label.add_theme_color_override("font_color", skin["cannon_color"])
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
		btn.set_meta("skin_id", skin["id"])
		btn.set_meta("skin_price", skin["price"])
		bg.add_child(btn)

		_buttons.append({
			"btn": btn,
			"bg": btn_bg,
			"skin_id": skin["id"],
			"price": skin["price"]
		})

func _connect_signals() -> void:
	for item in _buttons:
		item["btn"].pressed.connect(_on_skin_btn_pressed.bind(item["skin_id"], item["price"]))

	# ?? GameManager 閮?
	if GameManager.has_signal("skin_updated"):
		GameManager.skin_updated.connect(_on_skin_updated)
	if GameManager.has_signal("player_updated"):
		GameManager.player_updated.connect(_on_player_updated)

func _on_skin_btn_pressed(skin_id: String, price: int) -> void:
	# 撌脰???銝?隞颱?鈭?
	if skin_id == _equipped_skin:
		return

	# 撌脫????湔鋆?
	if skin_id in _owned_skins:
		emit_signal("skin_equip_requested", skin_id)
		NetworkManager.send_message({
			"type": "equip_skin",
			"payload": {"skin_id": skin_id}
		})
		return

	# ?芣???蝣箄?鞈潸眺嚗?撟?雲憭??潮?
	if _player_coins >= price:
		emit_signal("skin_buy_requested", skin_id)
		NetworkManager.send_message({
			"type": "buy_skin",
			"payload": {"skin_id": skin_id}
		})
	else:
		# ?馳銝雲嚗＊蝷箸?蝷?
		_show_insufficient_coins(price)

func _on_skin_updated(data: Dictionary) -> void:
	_equipped_skin = data.get("equipped_skin", "default")
	var owned = data.get("owned_skins", ["default"])
	_owned_skins = owned
	_refresh_ui()

func _on_player_updated(data: Dictionary) -> void:
	_player_coins = data.get("coins", 0)
	_refresh_ui()

## ?湔 UI ???擃漁撌脰????啣??芣???
func _refresh_ui() -> void:
	for item in _buttons:
		var skin_id = item["skin_id"]
		var bg = item["bg"]
		if not is_instance_valid(bg):
			continue

		if skin_id == _equipped_skin:
			# 撌脰??????擃漁
			bg.color = Color(0.3, 0.25, 0.05, 0.95)
		elif skin_id in _owned_skins:
			# 撌脫??鋆?嚗???
			bg.color = Color(0.05, 0.15, 0.35, 0.9)
		else:
			# ?芣????啗
			bg.color = Color(0.1, 0.08, 0.25, 0.9)

## ?馳銝雲?內
func _show_insufficient_coins(price: int) -> void:
	# 撱箇??冽??內璅惜
	var hint := Label.new()
	hint.text = "?馳銝雲嚗?閬?%d" % price
	hint.position = Vector2(0, -20)
	hint.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))
	if _pixel_font:
		hint.add_theme_font_override("font", _pixel_font)
		hint.add_theme_font_size_override("font_size", 10)
	add_child(hint)

	# 1.5 蝘?瘨仃
	var tween = create_tween()
	tween.tween_property(hint, "modulate:a", 0.0, 1.5)
	tween.tween_callback(hint.queue_free)

## ???嗅?鋆???閫摰儔
func get_equipped_skin_def() -> Dictionary:
	for skin in SKINS:
		if skin["id"] == _equipped_skin:
			return skin
	return SKINS[0]  # ?身餈? default
