## IceFishingPanel.gd ???圈撟賊?頛芰?Ｘ嚗AY-171嚗?
## 璆剔?靘?嚗ozy Fishing Life嚗?026-05-10嚗inter Wheel ??8 segments x2-x10 multipliers??
## + Ice Fishing Live嚗volution Gaming嚗heel triggers bonus fishing rounds??
## 閬死閮剛?嚗?
##   - wheel_start嚗犖嚗??刻撟????+ 銝剖亢 8 ?潸憚?斗?頧?+ ?迫?? + 5 蝘
##   - wheel_result嚗犖嚗?頛芰蝺拙???蝯???+ ??敶歲? + ?喃?閫閮? 8 蝘?
##   - mult_end嚗犖嚗?瘛∪?閮? + ?喳皛蝯?敶?嚗??湔/??/憿??嚗?
##   - wheel_broadcast嚗??嚗??典?璈怠???鈭箄孛?澆??憚?歹???
##   - ??x嚗??脤???嚗10x嚗蔗?嫣???
extends Node2D

# ---- 撣豢 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0
const WHEEL_SLOTS := 8
const WHEEL_RADIUS := 100.0

# 頛芰?澆?摰儔嚗? Server 蝡臭??湛?
const SLOT_DEFS = [
	{"mult": 2.0, "label": "?2", "color": Color(0.39, 0.71, 0.96)},
	{"mult": 3.0, "label": "?3", "color": Color(0.26, 0.65, 0.96)},
	{"mult": 4.0, "label": "?4", "color": Color(0.13, 0.59, 0.95)},
	{"mult": 5.0, "label": "?5", "color": Color(0.12, 0.53, 0.90)},
	{"mult": 6.0, "label": "?6", "color": Color(0.08, 0.40, 0.75)},
	{"mult": 7.0, "label": "?7", "color": Color(0.05, 0.28, 0.63)},
	{"mult": 8.0, "label": "?8", "color": Color(0.0, 0.74, 0.83)},
	{"mult": 10.0, "label": "?10", "color": Color(0.0, 0.90, 1.0)},
]

# ---- ???----
var _pixel_font: Font = null
var _wheel_node: Node2D = null      # 頛芰蝭暺?
var _stop_btn: Node2D = null        # ?迫??
var _countdown_lbl: Label = null    # ???閮?
var _mult_lbl: Label = null         # ??憿舐內
var _is_spinning: bool = false      # ?臬甇???
var _spin_speed: float = 720.0      # ???漲嚗漲/蝘?
var _target_angle: float = 0.0      # ?格??迫閫漲
var _current_angle: float = 0.0     # ?嗅?閫漲
var _is_stopping: bool = false      # ?臬甇?蝺拙?
var _wheel_result: int = 0          # 頛芰蝯??澆?蝝Ｗ?
var _spin_timer: float = 0.0        # ??閮?
var _spin_duration: float = 5.0     # ??????
var _mult_duration: float = 8.0     # ??????
var _mult_elapsed: float = 0.0      # ??撌脤???
var _is_mult_active: bool = false   # ?臬??瞈瘣颱葉
var _current_mult: float = 1.0      # ?嗅???

# ---- ????----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("ice_fishing_wheel"):
		GameManager.ice_fishing_wheel.connect(_on_ice_fishing_wheel)

# ---- 閮???----
func _process(delta: float) -> void:
	# 頛芰??
	if _is_spinning and not _is_stopping:
		_current_angle += _spin_speed * delta
		if is_instance_valid(_wheel_node):
			_wheel_node.rotation_degrees = _current_angle
		_spin_timer += delta
		if _spin_timer >= _spin_duration:
			# ?芸??迫
			_do_stop_wheel()

	# 頛芰蝺拙?
	elif _is_stopping:
		var diff = _target_angle - _current_angle
		# 蝣箔? diff ?舀迤?賂?????
		while diff < 0:
			diff += 360.0
		if diff < 2.0:
			_current_angle = _target_angle
			_is_spinning = false
			_is_stopping = false
			if is_instance_valid(_wheel_node):
				_wheel_node.rotation_degrees = _current_angle
		else:
			var decel = min(_spin_speed * delta, diff * 0.15)
			_current_angle += decel
			if is_instance_valid(_wheel_node):
				_wheel_node.rotation_degrees = _current_angle

	# ???
	if _is_mult_active:
		_mult_elapsed += delta
		var remaining = _mult_duration - _mult_elapsed
		if remaining < 0.0:
			remaining = 0.0
		if is_instance_valid(_countdown_lbl):
			_countdown_lbl.text = "?%.0f ??%.0f蝘? % [_current_mult, remaining]"

# ---- 銝餉?鈭辣?? ----
func _on_ice_fishing_wheel(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var player_id: String = data.get("player_id", "")
	var player_name: String = data.get("player_name", "")
	var multiplier: float = data.get("multiplier", 1.0)
	var label: String = data.get("label", "??")
	var wheel_result: int = data.get("wheel_result", 0)
	var duration_sec: int = data.get("duration_sec", 8)
	var kill_count: int = data.get("kill_count", 0)
	var total_bonus: int = data.get("total_bonus", 0)

	# ?斗?臬?航撌?
	var my_id: String = ""
	if GameManager.has_method("get_player_id"):
		my_id = GameManager.get_player_id()
	var is_mine: bool = (player_id == my_id)

	match phase:
		"wheel_start":
			if is_mine:
				_start_wheel(wheel_result, multiplier, data.get("spin_sec", 5))
		"wheel_broadcast":
			if not is_mine:
				_show_broadcast_banner(player_name)
		"wheel_result":
			if is_mine:
				_show_wheel_result(wheel_result, multiplier, label, duration_sec)
		"wheel_result_broadcast":
			if not is_mine:
				_show_result_broadcast(player_name, multiplier, label)
		"mult_end":
			if is_mine:
				_end_mult(multiplier, kill_count, total_bonus)

# ---- ????頛芰 ----
func _start_wheel(wheel_result: int, multiplier: float, spin_sec: int) -> void:
	_wheel_result = wheel_result
	_current_mult = multiplier
	_is_spinning = true
	_is_stopping = false
	_spin_timer = 0.0
	_spin_duration = float(spin_sec)
	_current_angle = 0.0

	# 閮??格??迫閫漲嚗?蝯??潭????剁?
	var slot_angle = 360.0 / WHEEL_SLOTS
	_target_angle = 360.0 * 5 + (360.0 - wheel_result * slot_angle - slot_angle / 2.0)

	# ?刻撟????
	_flash_screen(Color(0.0, 0.8, 1.0, 0.5), 0.35)

	# 撱箇?頛芰
	_create_wheel()

	# 撱箇??迫??
	_create_stop_button()

# ---- 撱箇?頛芰 ----
func _create_wheel() -> void:
	if is_instance_valid(_wheel_node):
		_wheel_node.queue_free()

	_wheel_node = Node2D.new()
	_wheel_node.position = Vector2(SCREEN_W / 2.0, SCREEN_H / 2.0)
	add_child(_wheel_node)

	# 頛芰???
	var bg = ColorRect.new()
	bg.size = Vector2(WHEEL_RADIUS * 2 + 20, WHEEL_RADIUS * 2 + 20)
	bg.position = Vector2(-(WHEEL_RADIUS + 10), -(WHEEL_RADIUS + 10))
	bg.color = Color(0.0, 0.1, 0.3, 0.9)
	_wheel_node.add_child(bg)

	# 8 ?摮?
	var slot_angle = 360.0 / WHEEL_SLOTS
	for i in range(WHEEL_SLOTS):
		var slot_def = SLOT_DEFS[i]
		var angle_rad = deg_to_rad(i * slot_angle)

		# ?澆??
		var slot_bg = ColorRect.new()
		slot_bg.size = Vector2(44, 44)
		slot_bg.position = Vector2(
			cos(angle_rad) * WHEEL_RADIUS - 22,
			sin(angle_rad) * WHEEL_RADIUS - 22
		)
		slot_bg.color = slot_def["color"]
		_wheel_node.add_child(slot_bg)

		# ?澆???
		var slot_lbl = Label.new()
		slot_lbl.text = slot_def["label"]
		slot_lbl.position = Vector2(
			cos(angle_rad) * WHEEL_RADIUS - 16,
			sin(angle_rad) * WHEEL_RADIUS - 10
		)
		slot_lbl.add_theme_color_override("font_color", Color.WHITE)
		if _pixel_font:
			slot_lbl.add_theme_font_override("font", _pixel_font)
		slot_lbl.add_theme_font_size_override("font_size", 14)
		_wheel_node.add_child(slot_lbl)

	# 銝剖???
	var center = ColorRect.new()
	center.size = Vector2(30, 30)
	center.position = Vector2(-15, -15)
	center.color = Color(0.0, 0.5, 0.8, 1.0)
	_wheel_node.add_child(center)

	# ??嚗??剁?
	var pointer = Label.new()
	pointer.text = "??"
	pointer.position = Vector2(-8, -(WHEEL_RADIUS + 30))
	pointer.add_theme_color_override("font_color", Color(1.0, 0.9, 0.0))
	if _pixel_font:
		pointer.add_theme_font_override("font", _pixel_font)
	pointer.add_theme_font_size_override("font_size", 20)
	# ??銝??憚?斗?頧???嗥?暺?
	add_child(pointer)
	pointer.position = Vector2(SCREEN_W / 2.0 - 8, SCREEN_H / 2.0 - WHEEL_RADIUS - 30)

# ---- 撱箇??迫?? ----
func _create_stop_button() -> void:
	if is_instance_valid(_stop_btn):
		_stop_btn.queue_free()

	_stop_btn = Node2D.new()
	_stop_btn.position = Vector2(SCREEN_W / 2.0, SCREEN_H / 2.0 + WHEEL_RADIUS + 50)
	add_child(_stop_btn)

	var btn_bg = ColorRect.new()
	btn_bg.size = Vector2(120, 40)
	btn_bg.position = Vector2(-60, -20)
	btn_bg.color = Color(0.0, 0.5, 0.8, 0.9)
	_stop_btn.add_child(btn_bg)

	var btn_lbl = Label.new()
	btn_lbl.text = "?? ?迫"
	btn_lbl.position = Vector2(-30, -12)
	btn_lbl.add_theme_color_override("font_color", Color.WHITE)
	if _pixel_font:
		btn_lbl.add_theme_font_override("font", _pixel_font)
	btn_lbl.add_theme_font_size_override("font_size", 16)
	_stop_btn.add_child(btn_lbl)

	# ?????
	var tween = _stop_btn.create_tween().set_loops()
	tween.tween_property(_stop_btn, "modulate:a", 0.6, 0.4)
	tween.tween_property(_stop_btn, "modulate:a", 1.0, 0.4)

# ---- ???迫頛芰 ----
func _do_stop_wheel() -> void:
	if not _is_spinning or _is_stopping:
		return
	_is_stopping = true
	_spin_timer = _spin_duration  # ?脫迫?活閫貊

	# 皜??迫??
	if is_instance_valid(_stop_btn):
		_stop_btn.queue_free()

	# ?潮?甇Ｚ??舐策 Server
	if GameManager.has_method("send_message"):
		GameManager.send_message("ice_fishing_wheel_stop", {})

# ---- 憿舐內頛芰蝯? ----
func _show_wheel_result(wheel_result: int, multiplier: float, label: String, duration_sec: int) -> void:
	_is_spinning = false
	_is_stopping = false
	_is_mult_active = true
	_mult_elapsed = 0.0
	_mult_duration = float(duration_sec)
	_current_mult = multiplier

	# 皜?頛芰
	if is_instance_valid(_wheel_node):
		var tween = create_tween()
		tween.tween_interval(0.5)
		tween.tween_property(_wheel_node, "modulate:a", 0.0, 0.3)
		tween.tween_callback(func(): if is_instance_valid(_wheel_node): _wheel_node.queue_free())

	# ??敶歲?
	var mult_node = Node2D.new()
	mult_node.position = Vector2(SCREEN_W / 2.0, SCREEN_H / 2.0)
	add_child(mult_node)

	var mult_lbl = Label.new()
	mult_lbl.text = label
	mult_lbl.position = Vector2(-30, -20)
	mult_lbl.add_theme_color_override("font_color", Color(0.0, 0.9, 1.0))
	if _pixel_font:
		mult_lbl.add_theme_font_override("font", _pixel_font)
	mult_lbl.add_theme_font_size_override("font_size", 48)
	mult_node.add_child(mult_lbl)

	var tween2 = create_tween()
	tween2.tween_property(mult_node, "scale", Vector2(1.5, 1.5), 0.2).set_ease(Tween.EASE_OUT)
	tween2.tween_property(mult_node, "scale", Vector2(1.0, 1.0), 0.15)
	tween2.tween_interval(0.8)
	tween2.tween_property(mult_node, "modulate:a", 0.0, 0.3)
	tween2.tween_callback(func(): if is_instance_valid(mult_node): mult_node.queue_free())

	# 撱箇??閮?
	_create_mult_countdown()

	# 擃???
	if multiplier >= 10.0:
		_flash_screen(Color(0.0, 1.0, 1.0, 0.6), 0.15)
		await get_tree().create_timer(0.2).timeout
		_flash_screen(Color(0.5, 0.0, 1.0, 0.5), 0.15)
		await get_tree().create_timer(0.2).timeout
		_flash_screen(Color(0.0, 1.0, 0.5, 0.4), 0.15)
	elif multiplier >= 7.0:
		_flash_screen(Color(0.0, 0.8, 1.0, 0.5), 0.15)
		await get_tree().create_timer(0.2).timeout
		_flash_screen(Color(0.0, 0.8, 1.0, 0.35), 0.15)
	else:
		_flash_screen(Color(0.0, 0.7, 1.0, 0.4), 0.2)

# ---- 撱箇????閮? ----
func _create_mult_countdown() -> void:
	if is_instance_valid(_countdown_lbl):
		_countdown_lbl.queue_free()

	_countdown_lbl = Label.new()
	_countdown_lbl.text = "?%.0f ??%.0f蝘? % [_current_mult, _mult_duration]"
	_countdown_lbl.position = Vector2(SCREEN_W - 150, 60)
	_countdown_lbl.add_theme_color_override("font_color", Color(0.0, 0.9, 1.0))
	if _pixel_font:
		_countdown_lbl.add_theme_font_override("font", _pixel_font)
	_countdown_lbl.add_theme_font_size_override("font_size", 18)
	add_child(_countdown_lbl)

# ---- ?冽?撱?璈怠? ----
func _show_broadcast_banner(player_name: String) -> void:
	var banner = Node2D.new()
	banner.position = Vector2(SCREEN_W / 2.0, -40)
	add_child(banner)

	var bg = ColorRect.new()
	bg.size = Vector2(400, 36)
	bg.position = Vector2(-200, -18)
	bg.color = Color(0.0, 0.2, 0.5, 0.8)
	banner.add_child(bg)

	var lbl = Label.new()
	lbl.text = "? %s 閫貊?圈撟賊?頛芰嚗? % player_name"
	lbl.position = Vector2(-185, -12)
	lbl.add_theme_color_override("font_color", Color(0.5, 0.9, 1.0))
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
	lbl.add_theme_font_size_override("font_size", 14)
	banner.add_child(lbl)

	var tween = create_tween()
	tween.tween_property(banner, "position:y", 28.0, 0.25).set_ease(Tween.EASE_OUT)
	tween.tween_interval(2.0)
	tween.tween_property(banner, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func(): if is_instance_valid(banner): banner.queue_free())

# ---- ?冽?蝯?撱? ----
func _show_result_broadcast(player_name: String, multiplier: float, label: String) -> void:
	var banner = Node2D.new()
	banner.position = Vector2(SCREEN_W / 2.0, -40)
	add_child(banner)

	var bg = ColorRect.new()
	bg.size = Vector2(440, 36)
	bg.position = Vector2(-220, -18)
	bg.color = Color(0.0, 0.15, 0.4, 0.85)
	banner.add_child(bg)

	var lbl = Label.new()
	lbl.text = "? %s ?圈頛芰 %s嚗???8 蝘?" % [player_name, label]
	lbl.position = Vector2(-210, -12)
	lbl.add_theme_color_override("font_color", Color(0.0, 0.9, 1.0))
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
	lbl.add_theme_font_size_override("font_size", 14)
	banner.add_child(lbl)

	var tween = create_tween()
	tween.tween_property(banner, "position:y", 28.0, 0.25).set_ease(Tween.EASE_OUT)
	tween.tween_interval(2.5)
	tween.tween_property(banner, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func(): if is_instance_valid(banner): banner.queue_free())

# ---- ??蝯? ----
func _end_mult(multiplier: float, kill_count: int, total_bonus: int) -> void:
	_is_mult_active = false

	# 皜??閮?
	if is_instance_valid(_countdown_lbl):
		var tween = create_tween()
		tween.tween_property(_countdown_lbl, "modulate:a", 0.0, 0.3)
		tween.tween_callback(func(): if is_instance_valid(_countdown_lbl): _countdown_lbl.queue_free())

	# 憿舐內蝯?敶?
	await get_tree().create_timer(0.3).timeout
	_show_result_popup(multiplier, kill_count, total_bonus)

# ---- 蝯?敶? ----
func _show_result_popup(multiplier: float, kill_count: int, total_bonus: int) -> void:
	var popup = Node2D.new()
	popup.position = Vector2(SCREEN_W + 200, SCREEN_H / 2.0 - 80)
	add_child(popup)

	var bg = ColorRect.new()
	bg.size = Vector2(260, 160)
	bg.position = Vector2(-130, -80)
	bg.color = Color(0.0, 0.1, 0.3, 0.92)
	popup.add_child(bg)

	var border = ColorRect.new()
	border.size = Vector2(264, 164)
	border.position = Vector2(-132, -82)
	border.color = Color(0.0, 0.7, 1.0, 0.8)
	popup.add_child(border)
	popup.move_child(border, 0)

	var title_lbl = Label.new()
	title_lbl.text = "? ?圈頛芰蝯?嚗?"
	title_lbl.position = Vector2(-120, -70)
	title_lbl.add_theme_color_override("font_color", Color(0.0, 0.9, 1.0))
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
	title_lbl.add_theme_font_size_override("font_size", 16)
	popup.add_child(title_lbl)

	var mult_lbl = Label.new()
	mult_lbl.text = "頛芰??嚗?.0f" % multiplier
	mult_lbl.position = Vector2(-110, -30)
	mult_lbl.add_theme_color_override("font_color", Color(0.8, 0.95, 1.0))
	if _pixel_font:
		mult_lbl.add_theme_font_override("font", _pixel_font)
	mult_lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(mult_lbl)

	var kill_lbl = Label.new()
	kill_lbl.text = "??格?嚗?d ?? % kill_count"
	kill_lbl.position = Vector2(-110, 0)
	kill_lbl.add_theme_color_override("font_color", Color(0.8, 0.95, 1.0))
	if _pixel_font:
		kill_lbl.add_theme_font_override("font", _pixel_font)
	kill_lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(kill_lbl)

	var bonus_lbl = Label.new()
	bonus_lbl.text = "憿??嚗?d" % total_bonus
	bonus_lbl.position = Vector2(-110, 30)
	bonus_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	if _pixel_font:
		bonus_lbl.add_theme_font_override("font", _pixel_font)
	bonus_lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(bonus_lbl)

	var tween = create_tween()
	tween.tween_property(popup, "position:x", SCREEN_W - 160.0, 0.35).set_ease(Tween.EASE_OUT)
	tween.tween_interval(2.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func(): if is_instance_valid(popup): popup.queue_free())

# ---- ?刻撟???----
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.size = Vector2(SCREEN_W, SCREEN_H)
	flash.color = color
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())
