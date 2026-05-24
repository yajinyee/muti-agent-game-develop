## VortexFishPanel.gd ??瞍拇蒂擳黎?詨??Ｘ嚗AY-169嚗?
## 璆剔?靘?嚗cean King嚗oogle Play 2026嚗ortex Fish ??catching a Vortex Fish will suck
## all fish of the same species in the area into a whirlpool, capturing them all at once.??
## 閬死閮剛?嚗?
##   - vortex_start嚗?Ｗ?瘛梯??? + ?璈怠?皛 + 瞍拇蒂???嚗葉憭殷?+ ?格??賊??內
##   - ?芸楛閫貊??銝剖亢憭??? 璅?敶歲? + ?憬皜血撘葉嚗?蝷?
##   - vortex_suck嚗??璅?嚗璅??憬皜虫葉敹?? + ?詨閮??+ 撠???
##   - vortex_end嚗?Ｗ????? + ?喳皛蝯?敶?嚗?交/?嚗?
##   - ????????????敶抵銝???
extends Node2D

# ---- 撣豢 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0

# ---- ???----
var _pixel_font: Font = null
var _banner: Node2D = null          # ?璈怠?
var _vortex_center: Node2D = null   # 瞍拇蒂銝剖??蝭暺?
var _suck_counter_lbl: Label = null # ?詨閮??
var _is_my_vortex: bool = false     # ?臬?航撌梯孛?潛?瞍拇蒂
var _vortex_x: float = SCREEN_W / 2.0
var _vortex_y: float = SCREEN_H / 2.0
var _target_count: int = 0          # ???詨?格???
var _sucked_count: int = 0          # 撌脣?交
var _is_active: bool = false        # ?臬甇?瞍拇蒂銝?
var _vortex_rotation: float = 0.0   # 瞍拇蒂??閫漲

# ---- ????----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("vortex_fish"):
		GameManager.vortex_fish.connect(_on_vortex_fish)

# ---- 瞍拇蒂??? ----
func _process(delta: float) -> void:
	if not _is_active:
		return
	_vortex_rotation += delta * 180.0  # 瘥??? 180 摨?
	if is_instance_valid(_vortex_center):
		_vortex_center.rotation_degrees = _vortex_rotation

# ---- 銝餉?鈭辣?? ----
func _on_vortex_fish(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var trigger_id: String = data.get("trigger_id", "")
	var trigger_name: String = data.get("trigger_name", "")"
	var vortex_x: float = data.get("vortex_x", SCREEN_W / 2.0)
	var vortex_y: float = data.get("vortex_y", SCREEN_H / 2.0)
	var group_name: String = data.get("group_name", "?箇??格?蝢?)"
	var target_count: int = data.get("target_count", 0)

	match phase:
		"vortex_start":
			_start_vortex(trigger_id, trigger_name, vortex_x, vortex_y, group_name, target_count)
		"vortex_suck":
			var suck_entry = data.get("suck_entry", null)
			var suck_index: int = data.get("suck_index", 0)
			_on_suck(suck_entry, suck_index)
		"vortex_end":
			var killed_count: int = data.get("killed_count", 0)
			var total_reward: int = data.get("total_reward", 0)
			_end_vortex(killed_count, total_reward)

# ---- 瞍拇蒂?? ----
func _start_vortex(trigger_id: String, trigger_name: String, vx: float, vy: float, group_name: String, target_count: int) -> void:
	_is_active = true
	_vortex_x = vx
	_vortex_y = vy
	_target_count = target_count
	_sucked_count = 0
	_vortex_rotation = 0.0

	# ?斗?臬?航撌梯孛??
	var my_id: String = ""
	if GameManager.has_method("get_player_id"):
		my_id = GameManager.get_player_id()
	_is_my_vortex = (trigger_id == my_id)

	# ?刻撟楛????
	_flash_screen(Color(0.0, 0.5, 1.0, 0.6), 0.4)

	# 撱箇?瞍拇蒂銝剖??
	_create_vortex_center(vx, vy)

	# 撱箇??璈怠?
	_create_banner(trigger_name, group_name, target_count)

	# 撱箇??詨閮??
	_create_suck_counter()

	# ?芸楛閫貊??銝剖亢憭??? 璅?敶歲
	if _is_my_vortex:
		_show_my_trigger_anim()

# ---- 撱箇?瞍拇蒂銝剖?? ----
func _create_vortex_center(vx: float, vy: float) -> void:
	if is_instance_valid(_vortex_center):
		_vortex_center.queue_free()

	_vortex_center = Node2D.new()
	_vortex_center.position = Vector2(vx, vy)
	add_child(_vortex_center)

	# 瞍拇蒂憭?嚗楛?嚗?
	for i in range(3):
		var ring = ColorRect.new()
		var r := 40.0 + i * 20.0
		ring.size = Vector2(r * 2, r * 2)
		ring.position = Vector2(-r, -r)
		ring.color = Color(0.0, 0.3 + i * 0.2, 1.0, 0.3 - i * 0.08)
		_vortex_center.add_child(ring)

	# 瞍拇蒂銝剖???
	var center_dot = ColorRect.new()
	center_dot.size = Vector2(20, 20)
	center_dot.position = Vector2(-10, -10)
	center_dot.color = Color(0.0, 0.8, 1.0, 0.9)
	_vortex_center.add_child(center_dot)

	# 瞍拇蒂?內
	var icon_lbl = Label.new()
	icon_lbl.text = "??"
	icon_lbl.position = Vector2(-16, -16)
	if _pixel_font:
		icon_lbl.add_theme_font_override("font", _pixel_font)
	icon_lbl.add_theme_font_size_override("font_size", 28)
	_vortex_center.add_child(icon_lbl)

# ---- 撱箇??璈怠? ----
func _create_banner(trigger_name: String, group_name: String, target_count: int) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()

	_banner = Node2D.new()
	_banner.position = Vector2(SCREEN_W / 2.0, -60)
	add_child(_banner)

	# 璈怠??
	var bg = ColorRect.new()
	bg.size = Vector2(600, 52)
	bg.position = Vector2(-300, -26)
	bg.color = Color(0.0, 0.2, 0.6, 0.88)
	_banner.add_child(bg)

	# 璈怠???
	var lbl = Label.new()
	lbl.text = "?? %s 閫貊瞍拇蒂擳??詨 %d ??s嚗? % [trigger_name, target_count, group_name]"
	lbl.position = Vector2(-290, -18)
	lbl.add_theme_color_override("font_color", Color(0.4, 0.9, 1.0))
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
	lbl.add_theme_font_size_override("font_size", 16)
	_banner.add_child(lbl)

	# 璈怠?皛?
	var tween = create_tween()
	tween.tween_property(_banner, "position:y", 36.0, 0.3).set_ease(Tween.EASE_OUT)

# ---- 撱箇??詨閮??----
func _create_suck_counter() -> void:
	if is_instance_valid(_suck_counter_lbl):
		_suck_counter_lbl.queue_free()

	_suck_counter_lbl = Label.new()
	_suck_counter_lbl.text = "?詨嚗? / %d" % _target_count
	_suck_counter_lbl.position = Vector2(SCREEN_W / 2.0 - 80, SCREEN_H - 80)
	_suck_counter_lbl.add_theme_color_override("font_color", Color(0.4, 0.9, 1.0))
	if _pixel_font:
		_suck_counter_lbl.add_theme_font_override("font", _pixel_font)
	_suck_counter_lbl.add_theme_font_size_override("font_size", 18)
	add_child(_suck_counter_lbl)

# ---- ?芸楛閫貊? ----
func _show_my_trigger_anim() -> void:
	var anim_node = Node2D.new()
	anim_node.position = Vector2(SCREEN_W / 2.0, SCREEN_H / 2.0)
	add_child(anim_node)

	var lbl = Label.new()
	lbl.text = "??"
	lbl.position = Vector2(-24, -24)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
	lbl.add_theme_font_size_override("font_size", 48)
	anim_node.add_child(lbl)

	var sub_lbl = Label.new()
	sub_lbl.text = "瞍拇蒂?詨?銝哨?"
	sub_lbl.position = Vector2(-60, 30)
	sub_lbl.add_theme_color_override("font_color", Color(0.4, 0.9, 1.0))
	if _pixel_font:
		sub_lbl.add_theme_font_override("font", _pixel_font)
	sub_lbl.add_theme_font_size_override("font_size", 16)
	anim_node.add_child(sub_lbl)

	# 敶歲?
	var tween = create_tween()
	tween.tween_property(anim_node, "scale", Vector2(1.4, 1.4), 0.15).set_ease(Tween.EASE_OUT)
	tween.tween_property(anim_node, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(1.0)
	tween.tween_property(anim_node, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func(): if is_instance_valid(anim_node): anim_node.queue_free())

# ---- ?格?鋡怠??----
func _on_suck(suck_entry, suck_index: int) -> void:
	if suck_entry == null:
		return

	_sucked_count += 1

	# ?湔閮??
	if is_instance_valid(_suck_counter_lbl):
		_suck_counter_lbl.text = "?詨嚗?d / %d" % [_sucked_count, _target_count]

	# ?格?憌?瞍拇蒂銝剖?????
	var entry_x: float = suck_entry.get("x", _vortex_x)
	var entry_y: float = suck_entry.get("y", _vortex_y)
	var reward: int = suck_entry.get("reward", 0)

	_spawn_suck_particle(entry_x, entry_y, reward)

	# 撠???
	_flash_screen(Color(0.0, 0.6, 1.0, 0.15), 0.1)

# ---- ???詨蝎?? ----
func _spawn_suck_particle(from_x: float, from_y: float, reward: int) -> void:
	var particle = Node2D.new()
	particle.position = Vector2(from_x, from_y)
	add_child(particle)

	# 蝎??內
	var lbl = Label.new()
	lbl.text = "?"
	lbl.position = Vector2(-8, -8)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
	lbl.add_theme_font_size_override("font_size", 16)
	particle.add_child(lbl)

	# ???
	if reward > 0:
		var reward_lbl = Label.new()
		reward_lbl.text = "+%d" % reward
		reward_lbl.position = Vector2(10, -8)
		reward_lbl.add_theme_color_override("font_color", Color(0.4, 1.0, 0.8))
		if _pixel_font:
			reward_lbl.add_theme_font_override("font", _pixel_font)
		reward_lbl.add_theme_font_size_override("font_size", 12)
		particle.add_child(reward_lbl)

	# 憌?瞍拇蒂銝剖?
	var tween = create_tween()
	tween.tween_property(particle, "position", Vector2(_vortex_x, _vortex_y), 0.3).set_ease(Tween.EASE_IN)
	tween.tween_property(particle, "modulate:a", 0.0, 0.1)
	tween.tween_callback(func(): if is_instance_valid(particle): particle.queue_free())

# ---- 瞍拇蒂蝯? ----
func _end_vortex(killed_count: int, total_reward: int) -> void:
	_is_active = false

	# 皜?瞍拇蒂銝剖?
	if is_instance_valid(_vortex_center):
		var tween = create_tween()
		tween.tween_property(_vortex_center, "scale", Vector2(2.0, 2.0), 0.2)
		tween.tween_property(_vortex_center, "modulate:a", 0.0, 0.2)
		tween.tween_callback(func(): if is_instance_valid(_vortex_center): _vortex_center.queue_free())

	# 皜?璈怠?
	if is_instance_valid(_banner):
		var tween2 = create_tween()
		tween2.tween_property(_banner, "modulate:a", 0.0, 0.3)
		tween2.tween_callback(func(): if is_instance_valid(_banner): _banner.queue_free())

	# 皜?閮??
	if is_instance_valid(_suck_counter_lbl):
		var tween3 = create_tween()
		tween3.tween_property(_suck_counter_lbl, "modulate:a", 0.0, 0.3)
		tween3.tween_callback(func(): if is_instance_valid(_suck_counter_lbl): _suck_counter_lbl.queue_free())

	# ???
	if killed_count >= 8:
		# 敶抵銝???
		_flash_screen(Color(0.0, 0.8, 1.0, 0.7), 0.15)
		await get_tree().create_timer(0.2).timeout
		_flash_screen(Color(0.4, 0.0, 1.0, 0.6), 0.15)
		await get_tree().create_timer(0.2).timeout
		_flash_screen(Color(0.0, 1.0, 0.5, 0.5), 0.15)
	elif killed_count >= 5:
		# ????
		_flash_screen(Color(0.0, 0.6, 1.0, 0.6), 0.15)
		await get_tree().create_timer(0.2).timeout
		_flash_screen(Color(0.0, 0.6, 1.0, 0.4), 0.15)
	else:
		_flash_screen(Color(0.0, 0.5, 1.0, 0.5), 0.2)

	# 憿舐內蝯?敶?
	await get_tree().create_timer(0.3).timeout
	_show_result_popup(killed_count, total_reward)

# ---- 蝯?敶? ----
func _show_result_popup(killed_count: int, total_reward: int) -> void:
	var popup = Node2D.new()
	popup.position = Vector2(SCREEN_W + 200, SCREEN_H / 2.0 - 80)
	add_child(popup)

	# 敶??
	var bg = ColorRect.new()
	bg.size = Vector2(260, 160)
	bg.position = Vector2(-130, -80)
	bg.color = Color(0.0, 0.1, 0.3, 0.92)
	popup.add_child(bg)

	# ??
	var border = ColorRect.new()
	border.size = Vector2(264, 164)
	border.position = Vector2(-132, -82)
	border.color = Color(0.0, 0.6, 1.0, 0.8)
	popup.add_child(border)
	popup.move_child(border, 0)

	# 璅?
	var title_lbl = Label.new()
	title_lbl.text = "?? 瞍拇蒂擳之鞊嚗?"
	title_lbl.position = Vector2(-120, -70)
	title_lbl.add_theme_color_override("font_color", Color(0.4, 0.9, 1.0))
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
	title_lbl.add_theme_font_size_override("font_size", 16)
	popup.add_child(title_lbl)

	# ?詨??
	var killed_lbl = Label.new()
	killed_lbl.text = "?詨?格?嚗?d ?? % killed_count"
	killed_lbl.position = Vector2(-110, -30)
	killed_lbl.add_theme_color_override("font_color", Color(0.8, 0.95, 1.0))
	if _pixel_font:
		killed_lbl.add_theme_font_override("font", _pixel_font)
	killed_lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(killed_lbl)

	# ?
	var reward_lbl = Label.new()
	reward_lbl.text = "?脣??嚗?d" % total_reward
	reward_lbl.position = Vector2(-110, 0)
	reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	if _pixel_font:
		reward_lbl.add_theme_font_override("font", _pixel_font)
	reward_lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(reward_lbl)

	# 閰?
	var comment_lbl = Label.new()
	if killed_count >= 8:
		comment_lbl.text = "?? ?唾牧瞍拇蒂嚗?"
		comment_lbl.add_theme_color_override("font_color", Color(0.0, 1.0, 1.0))
	elif killed_count >= 5:
		comment_lbl.text = "? 憭扯??塚?"
		comment_lbl.add_theme_color_override("font_color", Color(0.4, 0.9, 1.0))
	else:
		comment_lbl.text = "?? 瞍拇蒂摰?"
		comment_lbl.add_theme_color_override("font_color", Color(0.7, 0.85, 1.0))
	comment_lbl.position = Vector2(-110, 35)
	if _pixel_font:
		comment_lbl.add_theme_font_override("font", _pixel_font)
	comment_lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(comment_lbl)

	# 皛?
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
