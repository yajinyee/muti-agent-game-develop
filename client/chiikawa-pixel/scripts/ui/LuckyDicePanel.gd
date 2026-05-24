## LuckyDicePanel.gd ??撟賊?撉啣?擳?選?DAY-175嚗?
## 璆剔?靘?嚗cean King 3 Plus?ast Bomb ??randomly triggered bonus??
## + ??璈平?ice Roll bonus ??roll dice to determine reward multiplier??
## 閬死閮剛?嚗?
##   - dice_start嚗犖嚗??刻撟??脤???+ 銝剖亢?拚?撉啣?皛曉??嚗?蝘?
##   - dice_broadcast嚗??嚗??典?璈怠???鈭箄孛?澆兢?狐摮???
##   - dice_result嚗犖嚗?撉啣?蝺拙?憿舐內暺 + 蝯?敶?嚗????/璅惜嚗?
##     - 暺7嚗??脣???暺12嚗?蝝???嚗???嚗換?脤???
##   - dice_jackpot嚗??嚗?誨?剜帖撟XX ?脣憭批嚗?
extends Node2D

# ---- 撣豢 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0
const DICE_SIZE := 64.0
const DICE_FACES = [
	[],                                    # 0嚗??剁?
	[[0.5, 0.5]],                          # 1嚗葉敹?
	[[0.25, 0.25], [0.75, 0.75]],          # 2嚗?閫?
	[[0.25, 0.25], [0.5, 0.5], [0.75, 0.75]], # 3嚗?閫?銝剖?
	[[0.25, 0.25], [0.75, 0.25], [0.25, 0.75], [0.75, 0.75]], # 4嚗?閫?
	[[0.25, 0.25], [0.75, 0.25], [0.5, 0.5], [0.25, 0.75], [0.75, 0.75]], # 5嚗?閫?銝剖?
	[[0.25, 0.2], [0.75, 0.2], [0.25, 0.5], [0.75, 0.5], [0.25, 0.8], [0.75, 0.8]], # 6嚗??銵?
]

# ---- ???----
var _pixel_font: Font = null
var _dice_container: Node2D = null  # 撉啣?摰孵
var _die1_node: Node2D = null       # 撉啣?1蝭暺?
var _die2_node: Node2D = null       # 撉啣?2蝭暺?
var _is_rolling: bool = false       # ?臬甇?皛曉?
var _roll_elapsed: float = 0.0      # 皛曉?撌脤???
var _roll_duration: float = 2.0     # 皛曉?????
var _roll_face1: int = 1            # 撉啣?1?嗅?憿舐內??
var _roll_face2: int = 1            # 撉啣?2?嗅?憿舐內??
var _face_timer: float = 0.0        # ?閮?
var _result_panel: Node2D = null    # 蝯?敶?

# ---- ????----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("lucky_dice_fish"):
		GameManager.lucky_dice_fish.connect(_on_lucky_dice_fish)

# ---- 閮???----
func _process(delta: float) -> void:
	if not _is_rolling:
		return

	_roll_elapsed += delta
	_face_timer += delta

	# 瘥?0.1 蝘?銝甈⊿狐摮嚗遝????
	if _face_timer >= 0.1:
		_face_timer = 0.0
		_roll_face1 = randi() % 6 + 1
		_roll_face2 = randi() % 6 + 1
		_update_dice_display(_roll_face1, _roll_face2)

	if _roll_elapsed >= _roll_duration:
		_is_rolling = false

# ---- 閮??? ----
func _on_lucky_dice_fish(data: Dictionary) -> void:
	var phase = data.get("phase", "")
	match phase:
		"dice_start":
			_handle_dice_start(data)
		"dice_broadcast":
			_handle_dice_broadcast(data)
		"dice_result":
			_handle_dice_result(data)
		"dice_jackpot":
			_handle_dice_jackpot(data)

# ---- dice_start嚗狐摮?憪遝??----
func _handle_dice_start(data: Dictionary) -> void:
	var roll_ms = data.get("roll_ms", 2000)
	_roll_duration = float(roll_ms) / 1000.0

	# ?刻撟??脤???
	_flash_screen(Color(1.0, 0.85, 0.0, 0.4))

	# 撱箇?撉啣?摰孵
	if is_instance_valid(_dice_container):
		_dice_container.queue_free()

	_dice_container = Node2D.new()
	_dice_container.position = Vector2(SCREEN_W / 2.0, SCREEN_H / 2.0)
	add_child(_dice_container)

	# 撱箇??拚?撉啣?
	_die1_node = _create_die_node(-DICE_SIZE - 10, 0)
	_die2_node = _create_die_node(10, 0)
	_dice_container.add_child(_die1_node)
	_dice_container.add_child(_die2_node)

	# ??皛曉?
	_is_rolling = true
	_roll_elapsed = 0.0
	_face_timer = 0.0

	# 敶歲?
	var tween = _dice_container.create_tween()
	tween.tween_property(_dice_container, "scale", Vector2(1.3, 1.3), 0.15)
	tween.tween_property(_dice_container, "scale", Vector2(1.0, 1.0), 0.15)

# ---- dice_broadcast嚗?誨?剜帖撟?----
func _handle_dice_broadcast(data: Dictionary) -> void:
	var player_name = data.get("player_name", "?拙振")
	_show_broadcast_banner("? %s 閫貊撟賊?撉啣?嚗? % player_name)"

# ---- dice_result嚗狐摮???----
func _handle_dice_result(data: Dictionary) -> void:
	var die1 = data.get("die1", 1)
	var die2 = data.get("die2", 1)
	var sum = data.get("sum", 2)
	var reward = data.get("reward", 0)
	var label = data.get("label", "")

	# ?迫皛曉?嚗＊蝷箸?蝯???
	_is_rolling = false
	_update_dice_display(die1, die2)

	# ?寞?暺瘙箏??寞?
	if sum == 12:
		_flash_screen(Color(1.0, 0.27, 0.0, 0.5))
		var tween = create_tween()
		tween.tween_interval(0.15)
		tween.tween_callback(func(): _flash_screen(Color(1.0, 0.27, 0.0, 0.5)))
	elif sum == 2:
		_flash_screen(Color(0.58, 0.0, 0.83, 0.5))
	elif sum == 7:
		_flash_screen(Color(1.0, 0.85, 0.0, 0.4))

	# 撉啣?蝺拙??嚗葬?橘?
	if is_instance_valid(_dice_container):
		var tween = _dice_container.create_tween()
		tween.tween_property(_dice_container, "scale", Vector2(1.2, 1.2), 0.1)
		tween.tween_property(_dice_container, "scale", Vector2(1.0, 1.0), 0.15)

	# 撱箇?蝯?敶?
	_show_result_panel(die1, die2, sum, reward, label)

	# 3 蝘?皜?撉啣?
	if is_instance_valid(_dice_container):
		var tween = _dice_container.create_tween()
		tween.tween_interval(3.0)
		tween.tween_property(_dice_container, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func():
			if is_instance_valid(_dice_container):
				_dice_container.queue_free()
				_dice_container = null
		)

# ---- dice_jackpot嚗之?剖?誨??----
func _handle_dice_jackpot(data: Dictionary) -> void:
	var player_name = data.get("player_name", "?拙振")
	var reward = data.get("reward", 0)
	_show_broadcast_banner("?? %s ?脣憭批嚗敺?%d ?馳嚗? % [player_name, reward])"

# ---- 頛嚗遣蝡狐摮?暺?----
func _create_die_node(offset_x: float, offset_y: float) -> Node2D:
	var die = Node2D.new()
	die.position = Vector2(offset_x, offset_y)

	# 撉啣??嚗?脣?閫敶ｇ?
	var bg = ColorRect.new()
	bg.size = Vector2(DICE_SIZE, DICE_SIZE)
	bg.position = Vector2(0, -DICE_SIZE / 2.0)
	bg.color = Color(0.95, 0.95, 0.95)
	die.add_child(bg)

	return die

# ---- 頛嚗?圈狐摮＊蝷?----
func _update_dice_display(face1: int, face2: int) -> void:
	if is_instance_valid(_die1_node):
		_draw_die_face(_die1_node, face1)
	if is_instance_valid(_die2_node):
		_draw_die_face(_die2_node, face2)

# ---- 頛嚗鼓鋆賡狐摮 ----
func _draw_die_face(die_node: Node2D, face: int) -> void:
	# 皜??暺?
	for child in die_node.get_children():
		if child.name.begins_with("dot_"):
			child.queue_free()

	if face < 1 or face > 6:
		return

	var dots = DICE_FACES[face]
	for i in range(dots.size()):
		var dot_pos = dots[i]
		var dot = ColorRect.new()
		dot.name = "dot_%d" % i
		dot.size = Vector2(10, 10)
		dot.position = Vector2(
			dot_pos[0] * DICE_SIZE - 5,
			dot_pos[1] * DICE_SIZE - DICE_SIZE / 2.0 - 5
		)
		dot.color = Color(0.1, 0.1, 0.1)
		die_node.add_child(dot)

# ---- 頛嚗＊蝷箇???蝒?----
func _show_result_panel(die1: int, die2: int, sum: int, reward: int, label: String) -> void:
	if is_instance_valid(_result_panel):
		_result_panel.queue_free()

	_result_panel = Node2D.new()
	add_child(_result_panel)

	# ?憿靘???
	var bg_color = Color(0.1, 0.1, 0.1, 0.9)
	if sum == 12:
		bg_color = Color(0.3, 0.1, 0.0, 0.9)
	elif sum == 2:
		bg_color = Color(0.15, 0.0, 0.2, 0.9)

	var bg = ColorRect.new()
	bg.size = Vector2(220, 110)
	bg.position = Vector2(0, -55)
	bg.color = bg_color
	_result_panel.add_child(bg)

	# 璅惜
	var label_lbl = Label.new()
	label_lbl.text = label
	label_lbl.position = Vector2(10, -50)
	if _pixel_font:
		label_lbl.add_theme_font_override("font", _pixel_font)
		label_lbl.add_theme_font_size_override("font_size", 14)
	label_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_result_panel.add_child(label_lbl)

	# 暺憿舐內
	var sum_lbl = Label.new()
	sum_lbl.text = "? %d + %d = %d" % [die1, die2, sum]
	sum_lbl.position = Vector2(10, -28)
	if _pixel_font:
		sum_lbl.add_theme_font_override("font", _pixel_font)
		sum_lbl.add_theme_font_size_override("font_size", 13)
	sum_lbl.add_theme_color_override("font_color", Color.WHITE)
	_result_panel.add_child(sum_lbl)

	# ?
	var reward_lbl = Label.new()
	reward_lbl.text = "?? ?嚗?%d ?馳" % reward
	reward_lbl.position = Vector2(10, -6)
	if _pixel_font:
		reward_lbl.add_theme_font_override("font", _pixel_font)
		reward_lbl.add_theme_font_size_override("font_size", 13)
	reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_result_panel.add_child(reward_lbl)

	# 敺?湔???
	_result_panel.position = Vector2(SCREEN_W + 50, SCREEN_H / 2.0)
	var tween = _result_panel.create_tween()
	tween.tween_property(_result_panel, "position:x", SCREEN_W - 240.0, 0.4).set_ease(Tween.EASE_OUT)

	# 3 蝘?瘛∪
	tween.tween_interval(3.0)
	tween.tween_property(_result_panel, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_result_panel):
			_result_panel.queue_free()
			_result_panel = null
	)

# ---- 頛嚗?誨?剜帖撟?----
func _show_broadcast_banner(text: String) -> void:
	var banner = ColorRect.new()
	banner.size = Vector2(SCREEN_W, 34)
	banner.position = Vector2(0, 0)
	banner.color = Color(0.1, 0.08, 0.0, 0.9)
	add_child(banner)

	var lbl = Label.new()
	lbl.text = text
	lbl.position = Vector2(10, 6)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 13)
	lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	banner.add_child(lbl)

	var tween = banner.create_tween()
	tween.tween_interval(2.5)
	tween.tween_property(banner, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(banner):
			banner.queue_free()
	)

# ---- 頛嚗?Ｗ??? ----
func _flash_screen(color: Color) -> void:
	var flash = ColorRect.new()
	flash.size = Vector2(SCREEN_W, SCREEN_H)
	flash.position = Vector2(0, 0)
	flash.color = color
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.25)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)
