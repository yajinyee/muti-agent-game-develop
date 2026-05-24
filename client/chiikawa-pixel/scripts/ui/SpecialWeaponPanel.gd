## SpecialWeaponPanel.gd ???寞?甇血?Ｘ嚗AY-089嚗?蝝?DAY-134嚗AY-141嚗AY-157嚗AY-166嚗?
## 憿舐內銋車?寞?甇血嚗敶??瑕?/?啣?/樴憸?餈質馱憌?/樴?擳/頠???暺?瞍拇蒂嚗??舀?芸???脣漲璇?
## 璆剔?靘?嚗?
##   - Fish Road 2026 ??8 tier 甇血蝟餌絞嚗敶??瑕??舀??畾郎??
##   - Royal Fishing 2026 Tornado Cannon ??樴憸冽??湛????詨??璅?
##   - JILI 2026 Auto-Charge ??瘥活??格??芸?蝝舐??嚗??閬?馳
##   - thechipotlemenu.com 2026 Automatic Target Locking Weapon ??AI ?芸?餈質馱?擃??格?
##   - megafishing.click 2026 Railgun (15x stake) ??蝛輸?湧??賢???蝯扔皜甇血
##   - Ocean King 3 2026 Vortex + Black Hole Fishing 2026 ??暺?瞍拇蒂嚗?亙?璅??嚗AY-166嚗?
extends Node2D

# ---- 撣豢 ----
const PANEL_WIDTH  := 720  # 銋郎?函??砍?撖穿?DAY-166嚗?
const PANEL_HEIGHT := 90
const BTN_WIDTH    := 72
const BTN_HEIGHT   := 62

# 甇血摰儔嚗? Server 蝡?specialweapon.go ?郊嚗?
const WEAPONS = [
	{
		"type": "bomb",
		"name": "?詨???,"
		"icon": "?",
		"color": Color(1.0, 0.42, 0.21),
		"cost": 500,
		"max_charges": 3,
		"charge_required": 20,
		"desc": "蝭??\n500/??,"
		"purchasable": true
	},
	{
		"type": "laser",
		"name": "?瑕???,"
		"icon": "??,"
		"color": Color(0.0, 1.0, 1.0),
		"cost": 800,
		"max_charges": 3,
		"charge_required": 30,
		"desc": "蝛輸??n800/??,"
		"purchasable": true
	},
	{
		"type": "freeze",
		"name": "?啣???,"
		"icon": "??",
		"color": Color(0.53, 0.81, 0.92),
		"cost": 300,
		"max_charges": 3,
		"charge_required": 15,
		"desc": "?典?啣?\n300/??,"
		"purchasable": true
	},
	{
		"type": "tornado",
		"name": "樴憸?,"
		"icon": "?儭?,"
		"color": Color(0.61, 0.35, 0.71),
		"cost": 0,
		"max_charges": 2,
		"charge_required": 50,
		"desc": "?典?\n??脣?",
		"purchasable": false
	},
	{
		"type": "homing",
		"name": "餈質馱敶?,"
		"icon": "?",
		"color": Color(1.0, 0.0, 0.5),
		"cost": 0,
		"max_charges": 3,
		"charge_required": 35,
		"desc": "AI餈質馱\n?1.5?",
		"purchasable": false
	},
	{
		"type": "dragon_wrath",
		"name": "樴",
		"icon": "??",
		"color": Color(1.0, 0.27, 0.0),
		"cost": 0,
		"max_charges": 1,
		"charge_required": 60,
		"desc": "瘚??沔n?典??",
		"purchasable": false
	},
	{
		"type": "torpedo",
		"name": "擳",
		"icon": "??",
		"color": Color(1.0, 0.84, 0.0),
		"cost": -1,
		"max_charges": 2,
		"charge_required": 25,
		"desc": "憭抒??n6x鞎餌",
		"purchasable": false
	},
	{
		"type": "railgun",
		"name": "頠???,"
		"icon": "?",
		"color": Color(0.0, 1.0, 1.0),
		"cost": -1,
		"max_charges": 1,
		"charge_required": 40,
		"desc": "蝛輸?廄n15x鞎餌",
		"purchasable": false
	},
	{
		"type": "black_hole",
		"name": "暺?",
		"icon": "??",
		"color": Color(0.4, 0.0, 0.8),
		"cost": -1,
		"max_charges": 2,
		"charge_required": 45,
		"desc": "?詨?\n10x鞎餌",
		"purchasable": false
	}
]

# ---- ???----
var _charges: Dictionary = {"bomb": 0, "laser": 0, "freeze": 0, "tornado": 0, "homing": 0, "dragon_wrath": 0, "torpedo": 0, "railgun": 0, "black_hole": 0}
var _progress: Dictionary = {"bomb": 0, "laser": 0, "freeze": 0, "tornado": 0, "homing": 0, "dragon_wrath": 0, "torpedo": 0, "railgun": 0, "black_hole": 0}
var _selected_weapon: String = ""
var _pixel_font: Font = null
var _buttons: Array = []
var _charge_labels: Array = []
var _progress_bars: Array = []  # ??脣漲璇?DAY-134嚗?

# ---- 閮? ----
signal weapon_selected(weapon_type: String)

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
	title.text = "?寞?甇血"
	title.add_theme_color_override("font_color", Color(1.0, 0.85, 0.3))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 10)
	bg.add_child(title)

	# ?郎?冽???
	for i in range(WEAPONS.size()):
		var w = WEAPONS[i]
		var btn_x = 4 + i * (BTN_WIDTH + 4)

		# ???
		var btn_bg := ColorRect.new()
		btn_bg.name = "BtnBG_%s" % w["type"]
		btn_bg.position = Vector2(btn_x, 16)
		btn_bg.size = Vector2(BTN_WIDTH, BTN_HEIGHT)
		btn_bg.color = Color(0.1, 0.12, 0.25, 0.9)
		bg.add_child(btn_bg)
		_buttons.append(btn_bg)

		# 樴憸函嚗畾換?脤?獢?DAY-134嚗?
		if w["type"] == "tornado":
			var border := ColorRect.new()
			border.position = Vector2(btn_x - 1, 15)
			border.size = Vector2(BTN_WIDTH + 2, BTN_HEIGHT + 2)
			border.color = Color(0.61, 0.35, 0.71, 0.6)
			border.z_index = -1
			bg.add_child(border)

		# 甇血?內
		var icon_lbl := Label.new()
		icon_lbl.position = Vector2(btn_x + 4, 18)
		icon_lbl.text = w["icon"]
		icon_lbl.add_theme_font_size_override("font_size", 22)
		bg.add_child(icon_lbl)

		# 甇血?迂
		var name_lbl := Label.new()
		name_lbl.position = Vector2(btn_x + 2, 44)
		name_lbl.size = Vector2(BTN_WIDTH - 2, 14)
		name_lbl.text = w["name"]
		name_lbl.add_theme_font_size_override("font_size", 9)
		name_lbl.add_theme_color_override("font_color", w["color"])
		if _pixel_font:
			name_lbl.add_theme_font_override("font", _pixel_font)
		bg.add_child(name_lbl)

		# ??賊?璅惜嚗銝?嚗?
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

		# ??脣漲璇?摨嚗AY-134嚗?
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

		# 暺????
		var area := Area2D.new()
		var col := CollisionShape2D.new()
		var shape := RectangleShape2D.new()
		shape.size = Vector2(BTN_WIDTH, BTN_HEIGHT)
		col.shape = shape
		col.position = Vector2(btn_x + BTN_WIDTH / 2.0, 16 + BTN_HEIGHT / 2.0)
		area.add_child(col)
		add_child(area)

		# ??closure ? weapon type
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
	if GameManager.has_signal("torpedo_result"):
		GameManager.torpedo_result.connect(_on_torpedo_result)
	if GameManager.has_signal("railgun_result"):
		GameManager.railgun_result.connect(_on_railgun_result)
	if GameManager.has_signal("black_hole_result"):
		GameManager.black_hole_result.connect(_on_black_hole_result)

# ---- 鈭辣?? ----

func _on_weapon_btn_pressed(wtype: String) -> void:
	var charges = _charges.get(wtype, 0)

	if charges > 0:
		# ???踝??脣??璅芋撘???乩蝙?典?湔郎?剁?
		if wtype == "freeze" or wtype == "tornado" or wtype == "homing" or wtype == "dragon_wrath":
			# ?啣???樴憸函/餈質馱憌?/樴??嚗?乩蝙?剁??函?Ｘ????芸?餈質馱嚗??閬?璅?
			NetworkManager.send_use_special_weapon(wtype, 640.0, 360.0)
			_set_selected("")
		else:
			# ?詨?/?瑕?/擳/頠???暺?嚗脣?豢?璅∪?嚗?敺摰園??璅?蝵?
			if _selected_weapon == wtype:
				_set_selected("")  # ?活暺????豢?
			else:
				_set_selected(wtype)
	else:
		# 瘝??
		var w = _get_weapon_def(wtype)
		if w and w.get("purchasable", false):
			# ?航頃鞎瑞?甇血嚗頃鞎?
			NetworkManager.send_buy_special_weapon(wtype)
		else:
			# 銝鞈潸眺嚗??脤◢/餈質馱憌?嚗?憿舐內??內
			_show_charge_hint(wtype)

func _on_special_weapon_updated(data: Dictionary) -> void:
	_charges["bomb"] = data.get("bomb_charges", 0)
	_charges["laser"] = data.get("laser_charges", 0)
	_charges["freeze"] = data.get("freeze_charges", 0)
	_charges["tornado"] = data.get("tornado_charges", 0)
	_charges["homing"] = data.get("homing_charges", 0)
	_charges["dragon_wrath"] = data.get("dragon_wrath_charges", 0)
	_charges["torpedo"] = data.get("torpedo_charges", 0)
	_charges["railgun"] = data.get("railgun_charges", 0)
	_charges["black_hole"] = data.get("black_hole_charges", 0)
	_progress["bomb"] = data.get("bomb_charge_progress", 0)
	_progress["laser"] = data.get("laser_charge_progress", 0)
	_progress["freeze"] = data.get("freeze_charge_progress", 0)
	_progress["tornado"] = data.get("tornado_charge_progress", 0)
	_progress["homing"] = data.get("homing_charge_progress", 0)
	_progress["dragon_wrath"] = data.get("dragon_wrath_charge_progress", 0)
	_progress["torpedo"] = data.get("torpedo_charge_progress", 0)
	_progress["railgun"] = data.get("railgun_charge_progress", 0)
	_progress["black_hole"] = data.get("black_hole_charge_progress", 0)
	_update_charge_display()

func _on_special_weapon_fired(data: Dictionary) -> void:
	# 皜?豢????
	_set_selected("")

func _on_special_weapon_charged(data: Dictionary) -> void:
	# ?摰??嚗AY-134嚗?
	var wtype = data.get("weapon_type", "")
	var weapon_icon = data.get("weapon_icon", "?")
	var weapon_name = data.get("weapon_name", "")
	var new_charges = data.get("new_charges", 0)

	# ?湔???
	if wtype in _charges:
		_charges[wtype] = new_charges
		_progress[wtype] = 0  # ?蔭?脣漲
		_update_charge_display()

	# 憿舐內?摰??
	_show_charge_complete_effect(wtype, weapon_icon, weapon_name)

# ---- ?摰??寞?嚗AY-134嚗?---

func _show_charge_complete_effect(wtype: String, icon: String, name: String) -> void:
	# ?曉撠???嚗?暸?????
	var btn_idx = _get_weapon_index(wtype)
	if btn_idx < 0 or btn_idx >= _buttons.size():
		return

	var btn_bg = _buttons[btn_idx]
	if not is_instance_valid(btn_bg):
		return

	# ???嚗??脤???
	var tween = btn_bg.create_tween()
	tween.tween_property(btn_bg, "color", Color(1.0, 0.9, 0.2, 1.0), 0.1)
	tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.1)
	tween.tween_property(btn_bg, "color", Color(1.0, 0.9, 0.2, 1.0), 0.1)
	tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.1)

	# 蝮格敶歲
	var scale_tween = btn_bg.create_tween()
	scale_tween.tween_property(btn_bg, "scale", Vector2(1.15, 1.15), 0.12)
	scale_tween.tween_property(btn_bg, "scale", Vector2(1.0, 1.0), 0.12)

	# 憿舐內?摰??內嚗??冽帖撟?
	_show_charge_banner(icon, name)

func _show_charge_banner(icon: String, name: String) -> void:
	# ?券?蹂??寥＊蝷箇?急?蝷?
	var banner := Label.new()
	banner.text = "%s %s ?摰?嚗? % [icon, name]"
	banner.position = Vector2(0, -22)
	banner.size = Vector2(PANEL_WIDTH, 20)
	banner.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	banner.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	banner.add_theme_font_size_override("font_size", 11)
	if _pixel_font:
		banner.add_theme_font_override("font", _pixel_font)

	# ?
	var banner_bg := ColorRect.new()
	banner_bg.position = Vector2(0, -22)
	banner_bg.size = Vector2(PANEL_WIDTH, 20)
	banner_bg.color = Color(0.1, 0.08, 0.02, 0.9)
	add_child(banner_bg)
	add_child(banner)

	# 2 蝘?瘛∪
	var tween = banner.create_tween()
	tween.tween_interval(1.5)
	tween.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(banner): banner.queue_free()
		if is_instance_valid(banner_bg): banner_bg.queue_free()
	)

func _show_charge_hint(wtype: String) -> void:
	# 憿舐內??內嚗??脤◢?脖??航頃鞎瑟?嚗?
	var w = _get_weapon_def(wtype)
	if not w:
		return
	var required = w.get("charge_required", 50)
	var current = _progress.get(wtype, 0)

	var hint := Label.new()
	hint.text = "??格??嚗?%d/%d)" % [current, required]
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

# ---- ?豢?璅∪? ----

func _set_selected(wtype: String) -> void:
	_selected_weapon = wtype
	_update_button_highlight()
	emit_signal("weapon_selected", wtype)

func get_selected_weapon() -> String:
	return _selected_weapon

func clear_selection() -> void:
	_set_selected("")

# ---- UI ?湔 ----

func _update_charge_display() -> void:
	for i in range(WEAPONS.size()):
		var w = WEAPONS[i]
		var wtype = w["type"]
		var charges = _charges.get(wtype, 0)
		var progress = _progress.get(wtype, 0)
		var required = w.get("charge_required", 20)

		# ?湔??賊?
		if i < _charge_labels.size():
			var lbl = _charge_labels[i]
			if is_instance_valid(lbl):
				lbl.text = str(charges)
				if charges > 0:
					lbl.add_theme_color_override("font_color", Color(1.0, 1.0, 0.3))
				else:
					lbl.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))

		# ?湔???憿
		if i < _buttons.size():
			var btn_bg = _buttons[i]
			if is_instance_valid(btn_bg):
				if charges > 0:
					btn_bg.color = Color(0.1, 0.15, 0.3, 0.95)
				else:
					btn_bg.color = Color(0.08, 0.08, 0.15, 0.85)

		# ?湔??脣漲璇?DAY-134嚗?
		if i < _progress_bars.size():
			var prog_fill = _progress_bars[i]
			if is_instance_valid(prog_fill):
				var ratio = float(progress) / float(required) if required > 0 else 0.0
				ratio = clampf(ratio, 0.0, 1.0)
				var tween = prog_fill.create_tween()
				tween.tween_property(prog_fill, "size:x", BTN_WIDTH * ratio, 0.15)
				# ?亥??遛????
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
					# ?訾葉???鈭桅?獢???
					btn_bg.color = Color(0.2, 0.3, 0.6, 1.0)
					var tween = btn_bg.create_tween()
					tween.tween_property(btn_bg, "scale", Vector2(1.05, 1.05), 0.1)
					tween.tween_property(btn_bg, "scale", Vector2(1.0, 1.0), 0.1)
				else:
					var charges = _charges.get(w["type"], 0)
					btn_bg.color = Color(0.1, 0.15, 0.3, 0.95) if charges > 0 else Color(0.08, 0.08, 0.15, 0.85)

# ---- 頛?賣 ----

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

## 餈質馱憌??賭葉蝯?嚗AY-141嚗?
func _on_homing_missile_result(data: Dictionary) -> void:
	var killed: bool = data.get("killed", false)
	var multiplier: float = data.get("multiplier", 0.0)
	var final_reward: int = data.get("final_reward", 0)
	var message: String = data.get("message", "")

	if not killed or final_reward <= 0:
		return

	# 憿舐內餈質馱憌??賭葉蝯?嚗?蝝敶?嚗?
	var result_lbl := Label.new()
	result_lbl.text = "? ?%.0f ??+%d" % [multiplier, final_reward]
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

	# 銝筑瘛∪?
	var tween = result_lbl.create_tween()
	tween.tween_property(result_lbl, "position:y", result_lbl.position.y - 20, 1.0)
	tween.parallel().tween_property(result_lbl, "modulate:a", 0.0, 1.0)
	tween.tween_callback(func():
		if is_instance_valid(result_lbl): result_lbl.queue_free()
		if is_instance_valid(result_bg): result_bg.queue_free()
	)

	# 餈質馱憌?????嚗?蝝嚗?
	var homing_idx = _get_weapon_index("homing")
	if homing_idx >= 0 and homing_idx < _buttons.size():
		var btn_bg = _buttons[homing_idx]
		if is_instance_valid(btn_bg):
			var flash_tween = btn_bg.create_tween()
			flash_tween.tween_property(btn_bg, "color", Color(0.5, 0.0, 0.25, 1.0), 0.08)
			flash_tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.08)
			flash_tween.tween_property(btn_bg, "color", Color(0.5, 0.0, 0.25, 1.0), 0.08)
			flash_tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.08)

## 樴??蝯?嚗AY-154嚗?
func _on_dragon_wrath_result(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var total_reward: int = data.get("total_reward", 0)
	var killer_id: String = data.get("killer_id", "")

	# ?芾??撌梯孛?潛?蝯?
	if killer_id != NetworkManager.get_player_id():
		return

	if phase == "result" and total_reward > 0:
		# 憿舐內樴??蝯?嚗?蝝敶?嚗?
		var result_lbl := Label.new()
		result_lbl.text = "?? 瘚???+%d" % total_reward
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

		# 銝筑瘛∪?
		var tween = result_lbl.create_tween()
		tween.tween_property(result_lbl, "position:y", result_lbl.position.y - 24, 1.2)
		tween.parallel().tween_property(result_lbl, "modulate:a", 0.0, 1.2)
		tween.tween_callback(func():
			if is_instance_valid(result_lbl): result_lbl.queue_free()
			if is_instance_valid(result_bg): result_bg.queue_free()
		)

		# 樴?????璈??莎?
		var dw_idx = _get_weapon_index("dragon_wrath")
		if dw_idx >= 0 and dw_idx < _buttons.size():
			var btn_bg = _buttons[dw_idx]
			if is_instance_valid(btn_bg):
				var flash_tween = btn_bg.create_tween()
				flash_tween.tween_property(btn_bg, "color", Color(0.5, 0.13, 0.0, 1.0), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.5, 0.13, 0.0, 1.0), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.08)

## 擳?蝯?嚗AY-155嚗?
func _on_torpedo_result(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var total_reward: int = data.get("total_reward", 0)
	var shooter_id: String = data.get("shooter_id", "")
	var cost: int = data.get("cost", 0)

	# ?芾??撌梯孛?潛?蝯?
	if shooter_id != NetworkManager.get_player_id():
		return

	if phase == "result" and total_reward > 0:
		# 憿舐內擳蝯?嚗??脣?蝒?
		var net_reward = total_reward - cost
		var result_lbl := Label.new()
		if net_reward > 0:
			result_lbl.text = "?? +%d (鞎?d)" % [total_reward, cost]
		else:
			result_lbl.text = "?? %d (鞎?d)" % [total_reward, cost]
		result_lbl.position = Vector2(PANEL_WIDTH / 2.0 - 70, -40)
		result_lbl.size = Vector2(140, 20)
		result_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		var color = Color(1.0, 0.84, 0.0) if net_reward > 0 else Color(0.8, 0.4, 0.0)
		result_lbl.add_theme_color_override("font_color", color)
		result_lbl.add_theme_font_size_override("font_size", 12)
		if _pixel_font:
			result_lbl.add_theme_font_override("font", _pixel_font)

		var result_bg := ColorRect.new()
		result_bg.position = Vector2(PANEL_WIDTH / 2.0 - 72, -42)
		result_bg.size = Vector2(144, 24)
		result_bg.color = Color(0.1, 0.08, 0.0, 0.92)
		add_child(result_bg)
		add_child(result_lbl)

		# 銝筑瘛∪?
		var tween = result_lbl.create_tween()
		tween.tween_property(result_lbl, "position:y", result_lbl.position.y - 24, 1.2)
		tween.parallel().tween_property(result_lbl, "modulate:a", 0.0, 1.2)
		tween.tween_callback(func():
			if is_instance_valid(result_lbl): result_lbl.queue_free()
			if is_instance_valid(result_bg): result_bg.queue_free()
		)

		# 擳????嚗??莎?
		var torpedo_idx = _get_weapon_index("torpedo")
		if torpedo_idx >= 0 and torpedo_idx < _buttons.size():
			var btn_bg = _buttons[torpedo_idx]
			if is_instance_valid(btn_bg):
				var flash_tween = btn_bg.create_tween()
				flash_tween.tween_property(btn_bg, "color", Color(0.5, 0.42, 0.0, 1.0), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.5, 0.42, 0.0, 1.0), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.08)

## 頠??桃忽????DAY-157嚗?
func _on_railgun_result(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var total_reward: int = data.get("total_reward", 0)
	var shooter_id: String = data.get("shooter_id", "")
	var cost: int = data.get("cost", 0)

	# ?芾??撌梯孛?潛?蝯?
	if shooter_id != NetworkManager.get_player_id():
		return

	if phase == "result" and total_reward > 0:
		# 憿舐內頠??桃????敶?嚗?
		var net_reward = total_reward - cost
		var result_lbl := Label.new()
		if net_reward > 0:
			result_lbl.text = "? +%d (鞎?d)" % [total_reward, cost]
		else:
			result_lbl.text = "? %d (鞎?d)" % [total_reward, cost]
		result_lbl.position = Vector2(PANEL_WIDTH / 2.0 - 70, -40)
		result_lbl.size = Vector2(140, 20)
		result_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		var color = Color(0.0, 1.0, 1.0) if net_reward > 0 else Color(0.0, 0.6, 0.6)
		result_lbl.add_theme_color_override("font_color", color)
		result_lbl.add_theme_font_size_override("font_size", 12)
		if _pixel_font:
			result_lbl.add_theme_font_override("font", _pixel_font)

		var result_bg := ColorRect.new()
		result_bg.position = Vector2(PANEL_WIDTH / 2.0 - 72, -42)
		result_bg.size = Vector2(144, 24)
		result_bg.color = Color(0.0, 0.08, 0.1, 0.92)
		add_child(result_bg)
		add_child(result_lbl)

		# 銝筑瘛∪?
		var tween = result_lbl.create_tween()
		tween.tween_property(result_lbl, "position:y", result_lbl.position.y - 24, 1.2)
		tween.parallel().tween_property(result_lbl, "modulate:a", 0.0, 1.2)
		tween.tween_callback(func():
			if is_instance_valid(result_lbl): result_lbl.queue_free()
			if is_instance_valid(result_bg): result_bg.queue_free()
		)

		# 頠??格??????嚗?
		var railgun_idx = _get_weapon_index("railgun")
		if railgun_idx >= 0 and railgun_idx < _buttons.size():
			var btn_bg = _buttons[railgun_idx]
			if is_instance_valid(btn_bg):
				var flash_tween = btn_bg.create_tween()
				flash_tween.tween_property(btn_bg, "color", Color(0.0, 0.5, 0.5, 1.0), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.0, 0.5, 0.5, 1.0), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.08)

## 暺?瞍拇蒂?蝯?嚗AY-166嚗?
func _on_black_hole_result(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var total_reward: int = data.get("total_reward", 0)
	var shooter_id: String = data.get("shooter_id", "")
	var cost: int = data.get("cost", 0)
	var sucked_count: int = data.get("sucked_count", 0)

	# ?芾??撌梯孛?潛?蝯?
	if shooter_id != NetworkManager.get_player_id():
		return

	if phase == "result" and total_reward > 0:
		# 憿舐內暺?蝯?嚗換?脣?蝒?
		var net_reward = total_reward - cost
		var result_lbl := Label.new()
		if net_reward > 0:
			result_lbl.text = "?? ??d??+%d" % [sucked_count, total_reward]
		else:
			result_lbl.text = "?? ??d??%d" % [sucked_count, total_reward]
		result_lbl.position = Vector2(PANEL_WIDTH / 2.0 - 70, -40)
		result_lbl.size = Vector2(140, 20)
		result_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		var color = Color(0.6, 0.2, 1.0) if net_reward > 0 else Color(0.4, 0.1, 0.6)
		result_lbl.add_theme_color_override("font_color", color)
		result_lbl.add_theme_font_size_override("font_size", 12)
		if _pixel_font:
			result_lbl.add_theme_font_override("font", _pixel_font)

		var result_bg := ColorRect.new()
		result_bg.position = Vector2(PANEL_WIDTH / 2.0 - 72, -42)
		result_bg.size = Vector2(144, 24)
		result_bg.color = Color(0.06, 0.0, 0.12, 0.92)
		add_child(result_bg)
		add_child(result_lbl)

		# 銝筑瘛∪?
		var tween = result_lbl.create_tween()
		tween.tween_property(result_lbl, "position:y", result_lbl.position.y - 24, 1.2)
		tween.parallel().tween_property(result_lbl, "modulate:a", 0.0, 1.2)
		tween.tween_callback(func():
			if is_instance_valid(result_lbl): result_lbl.queue_free()
			if is_instance_valid(result_bg): result_bg.queue_free()
		)

		# 暺?????嚗換?莎?
		var bh_idx = _get_weapon_index("black_hole")
		if bh_idx >= 0 and bh_idx < _buttons.size():
			var btn_bg = _buttons[bh_idx]
			if is_instance_valid(btn_bg):
				var flash_tween = btn_bg.create_tween()
				flash_tween.tween_property(btn_bg, "color", Color(0.3, 0.0, 0.5, 1.0), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.3, 0.0, 0.5, 1.0), 0.08)
				flash_tween.tween_property(btn_bg, "color", Color(0.1, 0.15, 0.3, 0.95), 0.08)
