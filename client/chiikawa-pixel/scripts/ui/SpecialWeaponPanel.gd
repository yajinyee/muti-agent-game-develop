## SpecialWeaponPanel.gd — 特殊武器面板（DAY-089，升級 DAY-134，DAY-141）
## 顯示五種特殊武器（炸彈/雷射/冰凍/龍捲風/追蹤飛彈），支援自動充能進度條
## 業界依據：
##   - Fish Road 2026 有 8 tier 武器系統，炸彈/雷射是標配特殊武器
##   - Royal Fishing 2026 Tornado Cannon — 龍捲風掃場，旋轉吸入所有目標
##   - JILI 2026 Auto-Charge — 每次擊破目標自動累積充能，不需要花金幣
##   - thechipotlemenu.com 2026 Automatic Target Locking Weapon — AI 自動追蹤最高倍率目標
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 480  # 六武器版本加寬（DAY-154）
const PANEL_HEIGHT := 90
const BTN_WIDTH    := 72
const BTN_HEIGHT   := 62

# 武器定義（與 Server 端 specialweapon.go 同步）
const WEAPONS = [
	{
		"type": "bomb",
		"name": "炸彈砲",
		"icon": "💣",
		"color": Color(1.0, 0.42, 0.21),
		"cost": 500,
		"max_charges": 3,
		"charge_required": 20,
		"desc": "範圍爆炸\n500/發",
		"purchasable": true
	},
	{
		"type": "laser",
		"name": "雷射砲",
		"icon": "⚡",
		"color": Color(0.0, 1.0, 1.0),
		"cost": 800,
		"max_charges": 3,
		"charge_required": 30,
		"desc": "穿透射擊\n800/發",
		"purchasable": true
	},
	{
		"type": "freeze",
		"name": "冰凍砲",
		"icon": "❄️",
		"color": Color(0.53, 0.81, 0.92),
		"cost": 300,
		"max_charges": 3,
		"charge_required": 15,
		"desc": "全場冰凍\n300/發",
		"purchasable": true
	},
	{
		"type": "tornado",
		"name": "龍捲風",
		"icon": "🌪️",
		"color": Color(0.61, 0.35, 0.71),
		"cost": 0,
		"max_charges": 2,
		"charge_required": 50,
		"desc": "全場掃除\n充能獲得",
		"purchasable": false
	},
	{
		"type": "homing",
		"name": "追蹤彈",
		"icon": "🎯",
		"color": Color(1.0, 0.0, 0.5),
		"cost": 0,
		"max_charges": 3,
		"charge_required": 35,
		"desc": "AI追蹤\n×1.5獎勵",
		"purchasable": false
	},
	{
		"type": "dragon_wrath",
		"name": "龍怒雨",
		"icon": "🐉",
		"color": Color(1.0, 0.27, 0.0),
		"cost": 0,
		"max_charges": 1,
		"charge_required": 60,
		"desc": "流星雨\n全場打擊",
		"purchasable": false
	}
]

# ---- 狀態 ----
var _charges: Dictionary = {"bomb": 0, "laser": 0, "freeze": 0, "tornado": 0, "homing": 0, "dragon_wrath": 0}
var _progress: Dictionary = {"bomb": 0, "laser": 0, "freeze": 0, "tornado": 0, "homing": 0, "dragon_wrath": 0}
var _selected_weapon: String = ""
var _pixel_font: Font = null
var _buttons: Array = []
var _charge_labels: Array = []
var _progress_bars: Array = []  # 充能進度條（DAY-134）

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

	# 四個武器按鈕
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

		# 龍捲風砲：特殊紫色邊框（DAY-134）
		if w["type"] == "tornado":
			var border := ColorRect.new()
			border.position = Vector2(btn_x - 1, 15)
			border.size = Vector2(BTN_WIDTH + 2, BTN_HEIGHT + 2)
			border.color = Color(0.61, 0.35, 0.71, 0.6)
			border.z_index = -1
			bg.add_child(border)

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

		# 充能進度條（底部，DAY-134）
		var prog_bg := ColorRect.new()
		prog_bg.position = Vector2(btn_x, 16 + BTN_HEIGHT - 6)
		prog_bg.size = Vector2(BTN_WIDTH, 5)
		prog_bg.color = Color(0.1, 0.1, 0.2, 0.8)
		bg.add_child(prog_bg)

		var prog_fill := ColorRect.new()
		prog_fill.name = "ProgFill_%s" % w["type"]
		prog_fill.position = Vector2(btn_x, 16 + BTN_HEIGHT - 6)
		prog_fill.size = Vector2(0, 5)
		prog_fill.color = w["color"]
		bg.add_child(prog_fill)
		_progress_bars.append(prog_fill)

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
	if GameManager.has_signal("special_weapon_charged"):
		GameManager.special_weapon_charged.connect(_on_special_weapon_charged)
	if GameManager.has_signal("homing_missile_result"):
		GameManager.homing_missile_result.connect(_on_homing_missile_result)
	if GameManager.has_signal("dragon_wrath_result"):
		GameManager.dragon_wrath_result.connect(_on_dragon_wrath_result)

# ---- 事件處理 ----

func _on_weapon_btn_pressed(wtype: String) -> void:
	var charges = _charges.get(wtype, 0)

	if charges > 0:
		# 有充能：進入「選擇目標」模式（或直接使用全場武器）
		if wtype == "freeze" or wtype == "tornado" or wtype == "homing" or wtype == "dragon_wrath":
			# 冰凍砲/龍捲風砲/追蹤飛彈/龍怒流星雨：直接使用（全畫面效果或自動追蹤，不需要選擇目標）
			NetworkManager.send_use_special_weapon(wtype, 640.0, 360.0)
			_set_selected("")
		else:
			# 炸彈/雷射：進入選擇模式，等待玩家點擊目標位置
			if _selected_weapon == wtype:
				_set_selected("")  # 再次點擊取消選擇
			else:
				_set_selected(wtype)
	else:
		# 沒有充能
		var w = _get_weapon_def(wtype)
		if w and w.get("purchasable", false):
			# 可購買的武器：購買
			NetworkManager.send_buy_special_weapon(wtype)
		else:
			# 不可購買（龍捲風/追蹤飛彈）：顯示充能提示
			_show_charge_hint(wtype)

func _on_special_weapon_updated(data: Dictionary) -> void:
	_charges["bomb"] = data.get("bomb_charges", 0)
	_charges["laser"] = data.get("laser_charges", 0)
	_charges["freeze"] = data.get("freeze_charges", 0)
	_charges["tornado"] = data.get("tornado_charges", 0)
	_charges["homing"] = data.get("homing_charges", 0)
	_charges["dragon_wrath"] = data.get("dragon_wrath_charges", 0)
	_progress["bomb"] = data.get("bomb_charge_progress", 0)
	_progress["laser"] = data.get("laser_charge_progress", 0)
	_progress["freeze"] = data.get("freeze_charge_progress", 0)
	_progress["tornado"] = data.get("tornado_charge_progress", 0)
	_progress["homing"] = data.get("homing_charge_progress", 0)
	_progress["dragon_wrath"] = data.get("dragon_wrath_charge_progress", 0)
	_update_charge_display()

func _on_special_weapon_fired(data: Dictionary) -> void:
	# 清除選擇狀態
	_set_selected("")

func _on_special_weapon_charged(data: Dictionary) -> void:
	# 充能完成通知（DAY-134）
	var wtype = data.get("weapon_type", "")
	var weapon_icon = data.get("weapon_icon", "🔫")
	var weapon_name = data.get("weapon_name", "")
	var new_charges = data.get("new_charges", 0)

	# 更新充能數
	if wtype in _charges:
		_charges[wtype] = new_charges
		_progress[wtype] = 0  # 重置進度
		_update_charge_display()

	# 顯示充能完成動畫
	_show_charge_complete_effect(wtype, weapon_icon, weapon_name)

# ---- 充能完成特效（DAY-134）----

func _show_charge_complete_effect(wtype: String, icon: String, name: String) -> void:
	# 找到對應按鈕，播放閃爍動畫
	var btn_idx = _get_weapon_index(wtype)
	if btn_idx < 0 or btn_idx >= _buttons.size():
		return

	var btn_bg = _buttons[btn_idx]
	if not is_instance_valid(btn_bg):
		return

	# 閃爍動畫：金色閃光
	var tween = btn_bg.create_tween()
	tween.tween_property(btn_bg, "color", Color(1.0, 0.9, 0.2, 1.0), 0.1)
	tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.1)
	tween.tween_property(btn_bg, "color", Color(1.0, 0.9, 0.2, 1.0), 0.1)
	tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.1)

	# 縮放彈跳
	var scale_tween = btn_bg.create_tween()
	scale_tween.tween_property(btn_bg, "scale", Vector2(1.15, 1.15), 0.12)
	scale_tween.tween_property(btn_bg, "scale", Vector2(1.0, 1.0), 0.12)

	# 顯示充能完成提示（頂部橫幅）
	_show_charge_banner(icon, name)

func _show_charge_banner(icon: String, name: String) -> void:
	# 在面板上方顯示短暫提示
	var banner := Label.new()
	banner.text = "%s %s 充能完成！" % [icon, name]
	banner.position = Vector2(0, -22)
	banner.size = Vector2(PANEL_WIDTH, 20)
	banner.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	banner.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	banner.add_theme_font_size_override("font_size", 11)
	if _pixel_font:
		banner.add_theme_font_override("font", _pixel_font)

	# 背景
	var banner_bg := ColorRect.new()
	banner_bg.position = Vector2(0, -22)
	banner_bg.size = Vector2(PANEL_WIDTH, 20)
	banner_bg.color = Color(0.1, 0.08, 0.02, 0.9)
	add_child(banner_bg)
	add_child(banner)

	# 2 秒後淡出
	var tween = banner.create_tween()
	tween.tween_interval(1.5)
	tween.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(banner): banner.queue_free()
		if is_instance_valid(banner_bg): banner_bg.queue_free()
	)

func _show_charge_hint(wtype: String) -> void:
	# 顯示充能提示（龍捲風砲不可購買時）
	var w = _get_weapon_def(wtype)
	if not w:
		return
	var required = w.get("charge_required", 50)
	var current = _progress.get(wtype, 0)

	var hint := Label.new()
	hint.text = "擊破目標充能！(%d/%d)" % [current, required]
	hint.position = Vector2(0, -22)
	hint.size = Vector2(PANEL_WIDTH, 20)
	hint.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	hint.add_theme_color_override("font_color", Color(0.8, 0.6, 1.0))
	hint.add_theme_font_size_override("font_size", 10)
	if _pixel_font:
		hint.add_theme_font_override("font", _pixel_font)

	var hint_bg := ColorRect.new()
	hint_bg.position = Vector2(0, -22)
	hint_bg.size = Vector2(PANEL_WIDTH, 20)
	hint_bg.color = Color(0.08, 0.05, 0.15, 0.9)
	add_child(hint_bg)
	add_child(hint)

	var tween = hint.create_tween()
	tween.tween_interval(1.8)
	tween.tween_property(hint, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(hint): hint.queue_free()
		if is_instance_valid(hint_bg): hint_bg.queue_free()
	)

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
		var progress = _progress.get(wtype, 0)
		var required = w.get("charge_required", 20)

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

		# 更新充能進度條（DAY-134）
		if i < _progress_bars.size():
			var prog_fill = _progress_bars[i]
			if is_instance_valid(prog_fill):
				var ratio = float(progress) / float(required) if required > 0 else 0.0
				ratio = clampf(ratio, 0.0, 1.0)
				var tween = prog_fill.create_tween()
				tween.tween_property(prog_fill, "size:x", BTN_WIDTH * ratio, 0.15)
				# 接近充滿時閃爍
				if ratio > 0.8:
					prog_fill.modulate = Color(1.5, 1.5, 0.5)
				else:
					prog_fill.modulate = Color(1.0, 1.0, 1.0)

func _update_button_highlight() -> void:
	for i in range(WEAPONS.size()):
		var w = WEAPONS[i]
		if i < _buttons.size():
			var btn_bg = _buttons[i]
			if is_instance_valid(btn_bg):
				if w["type"] == _selected_weapon:
					# 選中狀態：亮邊框效果
					btn_bg.color = Color(0.2, 0.3, 0.6, 1.0)
					var tween = btn_bg.create_tween()
					tween.tween_property(btn_bg, "scale", Vector2(1.05, 1.05), 0.1)
					tween.tween_property(btn_bg, "scale", Vector2(1.0, 1.0), 0.1)
				else:
					var charges = _charges.get(w["type"], 0)
					btn_bg.color = Color(0.1, 0.15, 0.3, 0.95) if charges > 0 else Color(0.08, 0.08, 0.15, 0.85)

# ---- 輔助函數 ----

func _get_weapon_def(wtype: String) -> Dictionary:
	for w in WEAPONS:
		if w["type"] == wtype:
			return w
	return {}

func _get_weapon_index(wtype: String) -> int:
	for i in range(WEAPONS.size()):
		if WEAPONS[i]["type"] == wtype:
			return i
	return -1

## 追蹤飛彈命中結果（DAY-141）
func _on_homing_missile_result(data: Dictionary) -> void:
	var killed: bool = data.get("killed", false)
	var multiplier: float = data.get("multiplier", 0.0)
	var final_reward: int = data.get("final_reward", 0)
	var message: String = data.get("message", "")

	if not killed or final_reward <= 0:
		return

	# 顯示追蹤飛彈命中結果（粉紅色彈窗）
	var result_lbl := Label.new()
	result_lbl.text = "🎯 ×%.0f → +%d" % [multiplier, final_reward]
	result_lbl.position = Vector2(PANEL_WIDTH / 2.0 - 60, -40)
	result_lbl.size = Vector2(120, 20)
	result_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	result_lbl.add_theme_color_override("font_color", Color(1.0, 0.0, 0.5))
	result_lbl.add_theme_font_size_override("font_size", 13)
	if _pixel_font:
		result_lbl.add_theme_font_override("font", _pixel_font)

	var result_bg := ColorRect.new()
	result_bg.position = Vector2(PANEL_WIDTH / 2.0 - 62, -42)
	result_bg.size = Vector2(124, 24)
	result_bg.color = Color(0.1, 0.0, 0.08, 0.92)
	add_child(result_bg)
	add_child(result_lbl)

	# 上浮淡出動畫
	var tween = result_lbl.create_tween()
	tween.tween_property(result_lbl, "position:y", result_lbl.position.y - 20, 1.0)
	tween.parallel().tween_property(result_lbl, "modulate:a", 0.0, 1.0)
	tween.tween_callback(func():
		if is_instance_valid(result_lbl): result_lbl.queue_free()
		if is_instance_valid(result_bg): result_bg.queue_free()
	)

	# 追蹤飛彈按鈕閃爍（粉紅色）
	var homing_idx = _get_weapon_index("homing")
	if homing_idx >= 0 and homing_idx < _buttons.size():
		var btn_bg = _buttons[homing_idx]
		if is_instance_valid(btn_bg):
			var flash_tween = btn_bg.create_tween()
			flash_tween.tween_property(btn_bg, "color", Color(0.5, 0.0, 0.25, 1.0), 0.08)
			flash_tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.08)
			flash_tween.tween_property(btn_bg, "color", Color(0.5, 0.0, 0.25, 1.0), 0.08)
			flash_tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.08)

## 龍怒流星雨結果（DAY-154）
func _on_dragon_wrath_result(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var total_reward: int = data.get("total_reward", 0)
	var killer_id: String = data.get("killer_id", "")

	# 只處理自己觸發的結果
	if killer_id != NetworkManager.get_player_id():
		return

	if phase == "result" and total_reward > 0:
		# 顯示龍怒流星雨結果（橙紅色彈窗）
		var result_lbl := Label.new()
		result_lbl.text = "🐉 流星雨 +%d" % total_reward
		result_lbl.position = Vector2(PANEL_WIDTH / 2.0 - 70, -40)
		result_lbl.size = Vector2(140, 20)
		result_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		result_lbl.add_theme_color_override("font_color", Color(1.0, 0.4, 0.0))
		result_lbl.add_theme_font_size_override("font_size", 13)
		if _pixel_font:
			result_lbl.add_theme_font_override("font", _pixel_font)

		var result_bg := ColorRect.new()
		result_bg.position = Vector2(PANEL_WIDTH / 2.0 - 72, -42)
		result_bg.size = Vector2(144, 24)
		result_bg.color = Color(0.15, 0.05, 0.0, 0.92)
		add_child(result_bg)
		add_child(result_lbl)

		# 上浮淡出動畫
		var tween = result_lbl.create_tween()
		tween.tween_property(result_lbl, "position:y", result_lbl.position.y - 24, 1.2)
		tween.parallel().tween_property(result_lbl, "modulate:a", 0.0, 1.2)
		tween.tween_callback(func():
			if is_instance_valid(result_lbl): result_lbl.queue_free()
			if is_instance_valid(result_bg): result_bg.queue_free()
		)

		# 龍怒按鈕閃爍（橙紅色）
		var dw_idx = _get_weapon_index("dragon_wrath")
		if dw_idx >= 0 and dw_idx < _buttons.size():
			var btn_bg = _buttons[dw_idx]
			if is_instance_valid(btn_bg):
				var flash_tween = btn_bg.create_tween()
				flash_tween.tween_property(btn_bg, "color", Color(0.5, 0.13, 0.0, 1.0), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.5, 0.13, 0.0, 1.0), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.08)
