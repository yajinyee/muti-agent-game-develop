## SpecialWeaponPanel.gd — 特殊武器面板（DAY-089）
## 顯示三種特殊武器（炸彈/雷射/冰凍），玩家點擊購買或使用
## 業界依據：Fish Road 2026 有 8 tier 武器系統，炸彈/雷射是標配特殊武器
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 240
const PANEL_HEIGHT := 80
const BTN_WIDTH    := 70
const BTN_HEIGHT   := 60

# 武器定義（與 Server 端 specialweapon.go 同步）
const WEAPONS = [
	{
		"type": "bomb",
		"name": "炸彈砲",
		"icon": "💣",
		"color": Color(1.0, 0.42, 0.21),
		"cost": 500,
		"max_charges": 3,
		"desc": "範圍爆炸\n500/發"
	},
	{
		"type": "laser",
		"name": "雷射砲",
		"icon": "⚡",
		"color": Color(0.0, 1.0, 1.0),
		"cost": 800,
		"max_charges": 3,
		"desc": "穿透射擊\n800/發"
	},
	{
		"type": "freeze",
		"name": "冰凍砲",
		"icon": "❄️",
		"color": Color(0.53, 0.81, 0.92),
		"cost": 300,
		"max_charges": 3,
		"desc": "全場冰凍\n300/發"
	}
]

# ---- 狀態 ----
var _charges: Dictionary = {"bomb": 0, "laser": 0, "freeze": 0}
var _selected_weapon: String = ""  # 當前選中的武器（等待點擊目標）
var _pixel_font: Font = null
var _buttons: Array = []
var _charge_labels: Array = []

# ---- 訊號 ----
signal weapon_selected(weapon_type: String)

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
	bg.color = Color(0.05, 0.08, 0.18, 0.88)
	add_child(bg)

	# 標題
	var title := Label.new()
	title.position = Vector2(4, 2)
	title.text = "特殊武器"
	title.add_theme_color_override("font_color", Color(1.0, 0.85, 0.3))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 10)
	bg.add_child(title)

	# 三個武器按鈕
	for i in range(WEAPONS.size()):
		var w = WEAPONS[i]
		var btn_x = 4 + i * (BTN_WIDTH + 4)

		# 按鈕背景
		var btn_bg := ColorRect.new()
		btn_bg.name = "BtnBG_%s" % w["type"]
		btn_bg.position = Vector2(btn_x, 16)
		btn_bg.size = Vector2(BTN_WIDTH, BTN_HEIGHT)
		btn_bg.color = Color(0.1, 0.12, 0.25, 0.9)
		bg.add_child(btn_bg)
		_buttons.append(btn_bg)

		# 武器圖示
		var icon_lbl := Label.new()
		icon_lbl.position = Vector2(btn_x + 4, 18)
		icon_lbl.text = w["icon"]
		icon_lbl.add_theme_font_size_override("font_size", 22)
		bg.add_child(icon_lbl)

		# 武器名稱
		var name_lbl := Label.new()
		name_lbl.position = Vector2(btn_x + 2, 44)
		name_lbl.size = Vector2(BTN_WIDTH - 2, 14)
		name_lbl.text = w["name"]
		name_lbl.add_theme_font_size_override("font_size", 9)
		name_lbl.add_theme_color_override("font_color", w["color"])
		if _pixel_font:
			name_lbl.add_theme_font_override("font", _pixel_font)
		bg.add_child(name_lbl)

		# 充能數量標籤（右上角）
		var charge_lbl := Label.new()
		charge_lbl.name = "Charge_%s" % w["type"]
		charge_lbl.position = Vector2(btn_x + BTN_WIDTH - 18, 18)
		charge_lbl.size = Vector2(18, 14)
		charge_lbl.text = "0"
		charge_lbl.add_theme_font_size_override("font_size", 11)
		charge_lbl.add_theme_color_override("font_color", Color(1.0, 1.0, 0.3))
		charge_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
		if _pixel_font:
			charge_lbl.add_theme_font_override("font", _pixel_font)
		bg.add_child(charge_lbl)
		_charge_labels.append(charge_lbl)

		# 點擊區域
		var area := Area2D.new()
		var col := CollisionShape2D.new()
		var shape := RectangleShape2D.new()
		shape.size = Vector2(BTN_WIDTH, BTN_HEIGHT)
		col.shape = shape
		col.position = Vector2(btn_x + BTN_WIDTH / 2.0, 16 + BTN_HEIGHT / 2.0)
		area.add_child(col)
		add_child(area)

		# 用 closure 捕獲 weapon type
		var wtype = w["type"]
		area.input_event.connect(func(_viewport, event, _shape_idx):
			if event is InputEventMouseButton and event.pressed and event.button_index == MOUSE_BUTTON_LEFT:
				_on_weapon_btn_pressed(wtype)
		)

func _connect_signals() -> void:
	if GameManager.has_signal("special_weapon_updated"):
		GameManager.special_weapon_updated.connect(_on_special_weapon_updated)
	if GameManager.has_signal("special_weapon_fired"):
		GameManager.special_weapon_fired.connect(_on_special_weapon_fired)

# ---- 事件處理 ----

func _on_weapon_btn_pressed(wtype: String) -> void:
	var charges = _charges.get(wtype, 0)

	if charges > 0:
		# 有充能：進入「選擇目標」模式（或直接使用冰凍）
		if wtype == "freeze":
			# 冰凍砲：直接使用（全畫面效果，不需要選擇目標）
			NetworkManager.send_use_special_weapon(wtype, 640.0, 360.0)
			_set_selected("")
		else:
			# 炸彈/雷射：進入選擇模式，等待玩家點擊目標位置
			if _selected_weapon == wtype:
				_set_selected("")  # 再次點擊取消選擇
			else:
				_set_selected(wtype)
	else:
		# 沒有充能：購買
		NetworkManager.send_buy_special_weapon(wtype)

func _on_special_weapon_updated(data: Dictionary) -> void:
	_charges["bomb"] = data.get("bomb_charges", 0)
	_charges["laser"] = data.get("laser_charges", 0)
	_charges["freeze"] = data.get("freeze_charges", 0)
	_update_charge_display()

func _on_special_weapon_fired(data: Dictionary) -> void:
	# 清除選擇狀態
	_set_selected("")
	# 顯示命中效果（由 TargetManager 處理）

# ---- 選擇模式 ----

func _set_selected(wtype: String) -> void:
	_selected_weapon = wtype
	_update_button_highlight()
	emit_signal("weapon_selected", wtype)

func get_selected_weapon() -> String:
	return _selected_weapon

func clear_selection() -> void:
	_set_selected("")

# ---- UI 更新 ----

func _update_charge_display() -> void:
	for i in range(WEAPONS.size()):
		var w = WEAPONS[i]
		var wtype = w["type"]
		var charges = _charges.get(wtype, 0)

		# 更新充能數量
		if i < _charge_labels.size():
			var lbl = _charge_labels[i]
			if is_instance_valid(lbl):
				lbl.text = str(charges)
				if charges > 0:
					lbl.add_theme_color_override("font_color", Color(1.0, 1.0, 0.3))
				else:
					lbl.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))

		# 更新按鈕背景顏色
		if i < _buttons.size():
			var btn_bg = _buttons[i]
			if is_instance_valid(btn_bg):
				if charges > 0:
					btn_bg.color = Color(0.1, 0.15, 0.3, 0.95)
				else:
					btn_bg.color = Color(0.08, 0.08, 0.15, 0.85)

func _update_button_highlight() -> void:
	for i in range(WEAPONS.size()):
		var w = WEAPONS[i]
		if i < _buttons.size():
			var btn_bg = _buttons[i]
			if is_instance_valid(btn_bg):
				if w["type"] == _selected_weapon:
					# 選中狀態：亮邊框效果（用 modulate 模擬）
					btn_bg.color = Color(0.2, 0.3, 0.6, 1.0)
					# 縮放動畫
					var tween = btn_bg.create_tween()
					tween.tween_property(btn_bg, "scale", Vector2(1.05, 1.05), 0.1)
					tween.tween_property(btn_bg, "scale", Vector2(1.0, 1.0), 0.1)
				else:
					var charges = _charges.get(w["type"], 0)
					btn_bg.color = Color(0.1, 0.15, 0.3, 0.95) if charges > 0 else Color(0.08, 0.08, 0.15, 0.85)
