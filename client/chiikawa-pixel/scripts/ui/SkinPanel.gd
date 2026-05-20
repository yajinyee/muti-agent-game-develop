## SkinPanel.gd — 砲台外觀面板（DAY-071）
## 顯示可購買/裝備的砲台外觀
## 位置：BottomBar 右側（WeaponPanel 旁邊）
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 240
const PANEL_HEIGHT := 80
const BTN_WIDTH    := 54
const BTN_HEIGHT   := 60

# 外觀定義（與 Server 端 ws/protocol.go 同步）
const SKINS = [
	{
		"id": "default",
		"name": "標準",
		"icon": "🔫",
		"price": 0,
		"cannon_color": Color(0.9, 0.9, 0.9),
		"bullet_color": Color(1.0, 1.0, 0.8),
		"desc": "免費\n預設"
	},
	{
		"id": "golden",
		"name": "黃金",
		"icon": "✨",
		"price": 5000,
		"cannon_color": Color(1.0, 0.843, 0.0),
		"bullet_color": Color(1.0, 0.647, 0.0),
		"desc": "5000\n黃金"
	},
	{
		"id": "sakura",
		"name": "櫻花",
		"icon": "🌸",
		"price": 8000,
		"cannon_color": Color(1.0, 0.714, 0.773),
		"bullet_color": Color(1.0, 0.412, 0.706),
		"desc": "8000\n限定"
	},
	{
		"id": "rainbow",
		"name": "彩虹",
		"icon": "🌈",
		"price": 20000,
		"cannon_color": Color(1.0, 0.412, 0.706),
		"bullet_color": Color(0.0, 1.0, 1.0),
		"desc": "20000\n傳說"
	}
]

# ---- 節點引用 ----
var _buttons: Array = []
var _pixel_font: Font = null
var _equipped_skin: String = "default"
var _owned_skins: Array = ["default"]
var _player_coins: int = 0

# ---- 訊號 ----
signal skin_buy_requested(skin_id: String)
signal skin_equip_requested(skin_id: String)

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_build_ui()
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

func _build_ui() -> void:
	# 背景
	var bg := ColorRect.new()
	bg.position = Vector2(0, 0)
	bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	bg.color = Color(0.08, 0.05, 0.18, 0.85)
	add_child(bg)

	# 標題
	var title := Label.new()
	title.position = Vector2(4, 2)
	title.text = "外觀"
	title.add_theme_color_override("font_color", Color(1.0, 0.8, 1.0))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 10)
	bg.add_child(title)

	# 四個外觀按鈕
	for i in range(4):
		var skin = SKINS[i]
		var btn_x = 4 + i * (BTN_WIDTH + 2)

		var btn_bg := ColorRect.new()
		btn_bg.position = Vector2(btn_x, 14)
		btn_bg.size = Vector2(BTN_WIDTH, BTN_HEIGHT)
		btn_bg.color = Color(0.1, 0.08, 0.25, 0.9)
		btn_bg.name = "BtnBG_" + skin["id"]
		bg.add_child(btn_bg)

		# 外觀圖示
		var icon_label := Label.new()
		icon_label.position = Vector2(btn_x + 4, 16)
		icon_label.text = skin["icon"]
		if _pixel_font:
			icon_label.add_theme_font_override("font", _pixel_font)
			icon_label.add_theme_font_size_override("font_size", 18)
		bg.add_child(icon_label)

		# 外觀說明
		var desc_label := Label.new()
		desc_label.position = Vector2(btn_x + 2, 38)
		desc_label.text = skin["desc"]
		desc_label.add_theme_color_override("font_color", skin["cannon_color"])
		if _pixel_font:
			desc_label.add_theme_font_override("font", _pixel_font)
			desc_label.add_theme_font_size_override("font_size", 9)
		bg.add_child(desc_label)

		# 點擊區域（Button）
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

	# 連接 GameManager 訊號
	if GameManager.has_signal("skin_updated"):
		GameManager.skin_updated.connect(_on_skin_updated)
	if GameManager.has_signal("player_updated"):
		GameManager.player_updated.connect(_on_player_updated)

func _on_skin_btn_pressed(skin_id: String, price: int) -> void:
	# 已裝備：不做任何事
	if skin_id == _equipped_skin:
		return

	# 已擁有：直接裝備
	if skin_id in _owned_skins:
		emit_signal("skin_equip_requested", skin_id)
		NetworkManager.send_message({
			"type": "equip_skin",
			"payload": {"skin_id": skin_id}
		})
		return

	# 未擁有：確認購買（金幣足夠才發送）
	if _player_coins >= price:
		emit_signal("skin_buy_requested", skin_id)
		NetworkManager.send_message({
			"type": "buy_skin",
			"payload": {"skin_id": skin_id}
		})
	else:
		# 金幣不足：顯示提示
		_show_insufficient_coins(price)

func _on_skin_updated(data: Dictionary) -> void:
	_equipped_skin = data.get("equipped_skin", "default")
	var owned = data.get("owned_skins", ["default"])
	_owned_skins = owned
	_refresh_ui()

func _on_player_updated(data: Dictionary) -> void:
	_player_coins = data.get("coins", 0)
	_refresh_ui()

## 更新 UI 狀態（高亮已裝備，灰化未擁有）
func _refresh_ui() -> void:
	for item in _buttons:
		var skin_id = item["skin_id"]
		var bg = item["bg"]
		if not is_instance_valid(bg):
			continue

		if skin_id == _equipped_skin:
			# 已裝備：金色邊框高亮
			bg.color = Color(0.3, 0.25, 0.05, 0.95)
		elif skin_id in _owned_skins:
			# 已擁有未裝備：藍色
			bg.color = Color(0.05, 0.15, 0.35, 0.9)
		else:
			# 未擁有：灰色
			bg.color = Color(0.1, 0.08, 0.25, 0.9)

## 金幣不足提示
func _show_insufficient_coins(price: int) -> void:
	# 建立臨時提示標籤
	var hint := Label.new()
	hint.text = "金幣不足！需要 %d" % price
	hint.position = Vector2(0, -20)
	hint.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))
	if _pixel_font:
		hint.add_theme_font_override("font", _pixel_font)
		hint.add_theme_font_size_override("font_size", 10)
	add_child(hint)

	# 1.5 秒後消失
	var tween = create_tween()
	tween.tween_property(hint, "modulate:a", 0.0, 1.5)
	tween.tween_callback(hint.queue_free)

## 取得當前裝備的外觀定義
func get_equipped_skin_def() -> Dictionary:
	for skin in SKINS:
		if skin["id"] == _equipped_skin:
			return skin
	return SKINS[0]  # 預設返回 default
