## FreezeBombPanel.gd ???啣??詨?擳?選?DAY-170嚗?
## 璆剔?靘?嚗ing of Ocean 2026?he freezing blast pauses an entire school for a few seconds ??
## useful when a high-tier creature is escaping the frame.??
## 閬死閮剛?嚗?
##   - freeze_start嚗?Ｗ??啗??? + ?璈怠?皛 + ?寞??格?霈? + ?閮? 6 蝘?
##   - ?芸楛閫貊??銝剖亢憭??? 璅?敶歲? + ?畾璅歇?啣?嚗翰?餅??湛???蝷?
##   - ?啣???嚗畾璅＊蝷箏?嗅?????莎?+ ?閮?
##   - freeze_end嚗?嗥?鋆???+ 瘛∪???UI
##   - ????????????敶抵銝???
extends Node2D

# ---- 撣豢 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0

# ---- ???----
var _pixel_font: Font = null
var _banner: Node2D = null          # ?璈怠?
var _countdown_lbl: Label = null    # ?閮?
var _frozen_nodes: Dictionary = {}  # instanceID -> Node2D嚗?嗅???
var _is_my_freeze: bool = false     # ?臬?航撌梯孛?潛??啣?
var _duration_sec: int = 6          # ?啣?????
var _elapsed: float = 0.0           # 撌脤???
var _is_active: bool = false        # ?臬甇??啣?銝?
var _frozen_count: int = 0          # 鋡怠???格???

# ---- ????----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("freeze_bomb"):
		GameManager.freeze_bomb.connect(_on_freeze_bomb)
	# 餈質馱?格?蝘餃?嚗?啣?嗅???蝵?
	if GameManager.has_signal("target_updated"):
		GameManager.target_updated.connect(_on_target_updated)
	# ?格?鋡急??湔?蝘駁?唳??
	if GameManager.has_signal("target_killed"):
		GameManager.target_killed.connect(_on_target_killed)

# ---- 閮???----
func _process(delta: float) -> void:
	if not _is_active:
		return
	_elapsed += delta
	var remaining = float(_duration_sec) - _elapsed
	if remaining < 0.0:
		remaining = 0.0
	if is_instance_valid(_countdown_lbl):
		_countdown_lbl.text = "?? %.0f蝘? % remaining"

# ---- ?格?雿蔭餈質馱 ----
func _on_target_updated(data: Dictionary) -> void:
	if not _is_active:
		return
	var instance_id: String = data.get("instance_id", "")
	if not _frozen_nodes.has(instance_id):
		return
	var node = _frozen_nodes[instance_id]
	if not is_instance_valid(node):
		_frozen_nodes.erase(instance_id)
		return
	var x: float = data.get("x", node.position.x)
	var y: float = data.get("y", node.position.y)
	node.position = Vector2(x, y)

func _on_target_killed(data: Dictionary) -> void:
	if not _is_active:
		return
	var instance_id: String = data.get("instance_id", "")
	if not _frozen_nodes.has(instance_id):
		return
	var node = _frozen_nodes[instance_id]
	_frozen_nodes.erase(instance_id)
	# ?唳蝣??
	if is_instance_valid(node):
		var tween = create_tween()
		tween.tween_property(node, "scale", Vector2(1.5, 1.5), 0.1)
		tween.tween_property(node, "modulate:a", 0.0, 0.15)
		tween.tween_callback(func(): if is_instance_valid(node): node.queue_free())

# ---- 銝餉?鈭辣?? ----
func _on_freeze_bomb(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var trigger_id: String = data.get("trigger_id", "")
	var trigger_name: String = data.get("trigger_name", "")"
	var freeze_x: float = data.get("freeze_x", SCREEN_W / 2.0)
	var freeze_y: float = data.get("freeze_y", SCREEN_H / 2.0)
	var frozen_count: int = data.get("frozen_count", 0)
	var duration_sec: int = data.get("duration_sec", 6)

	match phase:
		"freeze_start":
			var frozen_targets = data.get("frozen_targets", [])
			_start_freeze(trigger_id, trigger_name, freeze_x, freeze_y, frozen_count, duration_sec, frozen_targets)
		"freeze_end":
			_end_freeze(frozen_count)

# ---- ?啣??? ----
func _start_freeze(trigger_id: String, trigger_name: String, fx: float, fy: float,
		frozen_count: int, duration_sec: int, frozen_targets: Array) -> void:
	_is_active = true
	_duration_sec = duration_sec
	_elapsed = 0.0
	_frozen_count = frozen_count

	# ?斗?臬?航撌梯孛??
	var my_id: String = ""
	if GameManager.has_method("get_player_id"):
		my_id = GameManager.get_player_id()
	_is_my_freeze = (trigger_id == my_id)

	# ?刻撟????
	_flash_screen(Color(0.0, 0.8, 1.0, 0.55), 0.4)

	# 撱箇??璈怠?
	_create_banner(trigger_name, frozen_count)

	# 撱箇??閮?
	_create_countdown()

	# ?箸??◤?啣??璅遣蝡?嗅???
	for entry in frozen_targets:
		var instance_id: String = entry.get("instance_id", "")
		var ex: float = entry.get("x", fx)
		var ey: float = entry.get("y", fy)
		if instance_id != "":
			_create_ice_halo(instance_id, ex, ey)

	# ?芸楛閫貊??銝剖亢憭??? 璅?敶歲
	if _is_my_freeze:
		_show_my_trigger_anim()

	# 憭?????
	if frozen_count >= 5:
		await get_tree().create_timer(0.5).timeout
		_flash_screen(Color(0.5, 0.9, 1.0, 0.4), 0.2)
		await get_tree().create_timer(0.25).timeout
		_flash_screen(Color(0.8, 0.95, 1.0, 0.3), 0.2)
	elif frozen_count >= 3:
		await get_tree().create_timer(0.5).timeout
		_flash_screen(Color(0.3, 0.8, 1.0, 0.35), 0.2)

# ---- 撱箇??璈怠? ----
func _create_banner(trigger_name: String, frozen_count: int) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()

	_banner = Node2D.new()
	_banner.position = Vector2(SCREEN_W / 2.0, -60)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.size = Vector2(580, 52)
	bg.position = Vector2(-290, -26)
	bg.color = Color(0.0, 0.15, 0.4, 0.88)
	_banner.add_child(bg)

	var lbl = Label.new()
	lbl.text = "?? %s 閫貊?啣??詨?擳?%d ?畾璅◤?啣?嚗? % [trigger_name, frozen_count]"
	lbl.position = Vector2(-275, -18)
	lbl.add_theme_color_override("font_color", Color(0.6, 0.95, 1.0))
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
	lbl.add_theme_font_size_override("font_size", 16)
	_banner.add_child(lbl)

	var tween = create_tween()
	tween.tween_property(_banner, "position:y", 36.0, 0.3).set_ease(Tween.EASE_OUT)

# ---- 撱箇??閮? ----
func _create_countdown() -> void:
	if is_instance_valid(_countdown_lbl):
		_countdown_lbl.queue_free()

	_countdown_lbl = Label.new()
	_countdown_lbl.text = "?? %d蝘? % _duration_sec"
	_countdown_lbl.position = Vector2(SCREEN_W - 120, 60)
	_countdown_lbl.add_theme_color_override("font_color", Color(0.6, 0.95, 1.0))
	if _pixel_font:
		_countdown_lbl.add_theme_font_override("font", _pixel_font)
	_countdown_lbl.add_theme_font_size_override("font_size", 20)
	add_child(_countdown_lbl)

# ---- 撱箇??唳?? ----
func _create_ice_halo(instance_id: String, x: float, y: float) -> void:
	var halo = Node2D.new()
	halo.position = Vector2(x, y)
	add_child(halo)

	# ?唳憭?
	var ring = ColorRect.new()
	ring.size = Vector2(56, 56)
	ring.position = Vector2(-28, -28)
	ring.color = Color(0.4, 0.85, 1.0, 0.35)
	halo.add_child(ring)

	# ?唳?內
	var icon_lbl = Label.new()
	icon_lbl.text = "??"
	icon_lbl.position = Vector2(-10, -10)
	if _pixel_font:
		icon_lbl.add_theme_font_override("font", _pixel_font)
	icon_lbl.add_theme_font_size_override("font_size", 18)
	halo.add_child(icon_lbl)

	# ???
	var tween = halo.create_tween().set_loops()
	tween.tween_property(halo, "modulate:a", 0.5, 0.6)
	tween.tween_property(halo, "modulate:a", 1.0, 0.6)

	_frozen_nodes[instance_id] = halo

# ---- ?芸楛閫貊? ----
func _show_my_trigger_anim() -> void:
	var anim_node = Node2D.new()
	anim_node.position = Vector2(SCREEN_W / 2.0, SCREEN_H / 2.0)
	add_child(anim_node)

	var lbl = Label.new()
	lbl.text = "??"
	lbl.position = Vector2(-24, -24)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
	lbl.add_theme_font_size_override("font_size", 48)
	anim_node.add_child(lbl)

	var sub_lbl = Label.new()
	sub_lbl.text = "?寞??格?撌脣??敹怠?嚗?"
	sub_lbl.position = Vector2(-80, 30)
	sub_lbl.add_theme_color_override("font_color", Color(0.6, 0.95, 1.0))
	if _pixel_font:
		sub_lbl.add_theme_font_override("font", _pixel_font)
	sub_lbl.add_theme_font_size_override("font_size", 14)
	anim_node.add_child(sub_lbl)

	var tween = create_tween()
	tween.tween_property(anim_node, "scale", Vector2(1.4, 1.4), 0.15).set_ease(Tween.EASE_OUT)
	tween.tween_property(anim_node, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(1.2)
	tween.tween_property(anim_node, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func(): if is_instance_valid(anim_node): anim_node.queue_free())

# ---- ?啣?蝯? ----
func _end_freeze(frozen_count: int) -> void:
	_is_active = false

	# 皜????嗅???蝣??嚗?
	for instance_id in _frozen_nodes.keys():
		var node = _frozen_nodes[instance_id]
		if is_instance_valid(node):
			var tween = create_tween()
			tween.tween_property(node, "scale", Vector2(1.3, 1.3), 0.1)
			tween.tween_property(node, "modulate:a", 0.0, 0.2)
			tween.tween_callback(func(): if is_instance_valid(node): node.queue_free())
	_frozen_nodes.clear()

	# 皜?璈怠?
	if is_instance_valid(_banner):
		var tween2 = create_tween()
		tween2.tween_property(_banner, "modulate:a", 0.0, 0.3)
		tween2.tween_callback(func(): if is_instance_valid(_banner): _banner.queue_free())

	# 皜??閮?
	if is_instance_valid(_countdown_lbl):
		var tween3 = create_tween()
		tween3.tween_property(_countdown_lbl, "modulate:a", 0.0, 0.3)
		tween3.tween_callback(func(): if is_instance_valid(_countdown_lbl): _countdown_lbl.queue_free())

	# ?唳蝣???
	_flash_screen(Color(0.7, 0.95, 1.0, 0.3), 0.25)

# ---- ?刻撟???----
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.size = Vector2(SCREEN_W, SCREEN_H)
	flash.color = color
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())
