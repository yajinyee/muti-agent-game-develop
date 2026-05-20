## WeaponPanel.gd — 武器升級面板（DAY-067）
## 顯示三個武器等級，玩家點擊切換
## 位置：BottomBar 左側
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 200
const PANEL_HEIGHT := 80
const BTN_WIDTH    := 58
const BTN_HEIGHT   := 60

# 武器定義（與 Server 端 data/tables.go 同步）
const WEAPONS = [
	{
		"level": 1,
		"name": "標準砲",
		"icon": "🔫",
		"color": Color(0.9, 0.9, 0.9),
		"extra_cost": 0,
		"power_mod": 1.00,
		"desc": "標準\n無額外費用"
	},
	{
		"level": 2,
		"name": "強化砲",
		"icon": "⚡",
		"color": Color(0.0, 0.9, 1.0),
		"extra_cost": 50,
		"power_mod": 1.25,
		"desc": "+25%\n+50/發"
	},
	{
		"level": 3,
		"name": "超級砲",
		"icon": "🌟",
		"color": Color(1.0, 0.85, 0.0),
		"extra_cost": 150,
		"power_mod": 1.60,
		"desc": "+60%\n+150/發"
	}
]

# ---- 節點引用 ----
var _buttons: Array = []
var _pixel_font: Font = null
var _current_level: int = 1

# ---- 訊號 ----
signal weapon_changed(level: int)

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
	bg.color = Color(0.05, 0.08, 0.18, 0.85)
	add_child(bg)

	# 標題
	var title := Label.new()
	title.position = Vector2(4, 2)
	title.text = "武器"
	title.add_theme_color_override("font_color", Color(0.7, 0.8, 1.0))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 10)
	bg.add_child(title)

	# 三個武器按鈕
	for i in range(3):
		var weapon = WEAPONS[i]
		var btn_x = 4 + i * (BTN_WIDTH + 4)

		var btn_bg := ColorRect.new()
		btn_bg.position = Vector2(btn_x, 14)
		btn_bg.size = Vector2(BTN_WIDTH, BTN_HEIGHT)
		btn_bg.color = Color(0.1, 0.15, 0.3, 0.9)
		bg.add_child(btn_bg)

		# 武器圖示
		var icon_label := Label.new()
		icon_label.position = Vector2(btn_x + 4, 16)
		icon_label.text = weapon["icon"]
		if _pixel_font:
			icon_label.add_theme_font_override("font", _pixel_font)
			icon_label.add_theme_font_size_override("font_size", 18)
		bg.add_child(icon_label)

		# 武器說明
		var desc_label := Label.new()
		desc_label.position = Vector2(btn_x + 2, 38)
		desc_label.text = weapon["desc"]
		desc_label.add_theme_color_override("font_color", weapon["color"])
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
	# 連接 GameManager 的 player_updated 訊號
	if GameManager.has_signal("player_updated"):
		GameManager.player_updated.connect(_on_player_updated)

# ---- 事件處理 ----
func _on_weapon_btn_pressed(level: int) -> void:
	if level == _current_level:
		return
	_current_level = level
	_update_selection()
	# 發送武器升級請求到 Server
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

# ---- UI 更新 ----
func _update_selection() -> void:
	for item in _buttons:
		var is_selected = (item["level"] == _current_level)
		var weapon = WEAPONS[item["level"] - 1]

		if is_selected:
			# 選中：亮色邊框 + 背景高亮
			item["bg"].color = Color(0.15, 0.25, 0.5, 0.95)
			# 加邊框效果（用 StyleBoxFlat）
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
			# 未選中：暗色背景
			item["bg"].color = Color(0.08, 0.12, 0.25, 0.85)
			var style := StyleBoxFlat.new()
			style.bg_color = Color(0.0, 0.0, 0.0, 0.0)
			item["btn"].add_theme_stylebox_override("normal", style)
			item["btn"].add_theme_stylebox_override("hover", style)
			item["btn"].add_theme_stylebox_override("pressed", style)
