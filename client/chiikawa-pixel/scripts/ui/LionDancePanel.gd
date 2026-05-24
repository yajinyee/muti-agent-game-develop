## LionDancePanel.gd ?????之???潮?選?DAY-168嚗?
## 璆剔?靘?嚗ortune King Jackpot嚗aDa Gaming 2026嚗ion Dance bonus ??triggered by special fish,
## delivers burst multiplier payouts with festive visual effects??
## 閬死閮剛?嚗?
##   - burst_start嚗?Ｗ?璈??? + ?璈怠?皛 + 璅??格?憿舐內?? + ?閮? 15 蝘?
##   - ?芸楛閫貊??銝剖亢憭??? 璅?敶歲? + ?翰?餅??湔?閮璅???蝷?
##   - ?璅??格???瘚桀?????嚗?Nx嚗? ????
##   - burst_end嚗楚?箸???UI + ?喳皛蝯?敶?嚗??湔/??/?嚗?
##   - ??x嚗??脤???嚗10x嚗蔗?嫣???
## 靽格迤嚗AY-168b嚗??頝??格?蝘餃?嚗? target_updated 閮??湔雿蔭嚗?
extends Node2D

# ---- 撣豢 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0

# ---- ???----
var _pixel_font: Font = null
var _banner: Node2D = null         # ?璈怠?
var _countdown_lbl: Label = null   # ?閮?
var _mark_nodes: Dictionary = {}   # instanceID -> Node2D嚗?閮??堆?
var _is_my_burst: bool = false     # ?臬?航撌梯孛?潛??
var _burst_mult: float = 1.0       # ?祆活???
var _duration_sec: int = 15        # ????
var _elapsed: float = 0.0          # 撌脤???
var _is_active: bool = false       # ?臬甇??

# ---- ????----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("lion_dance_burst"):
		GameManager.lion_dance_burst.connect(_on_lion_dance_burst)
	# 餈質馱?格?蝘餃?嚗?啣??唬?蝵?
	if GameManager.has_signal("target_updated"):
		GameManager.target_updated.connect(_on_target_updated)
	# ?格?鋡急??湔?蝘駁?
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
		_countdown_lbl.text = "?? %.0f蝘? % remaining"

# ---- ?格?雿蔭餈質馱 ----

## ?格?蝘餃???啣??唬?蝵?
func _on_target_updated(data: Dictionary) -> void:
	if not _is_active:
		return
	var instance_id: String = data.get("instance_id", "")
	if not _mark_nodes.has(instance_id):
		return
	var node = _mark_nodes[instance_id]
	if not is_instance_valid(node):
		_mark_nodes.erase(instance_id)
		return
	var x: float = data.get("x", node.position.x)
	var y: float = data.get("y", node.position.y)
	node.position = Vector2(x, y)

## ?格?鋡急??湔?蝘駁?嚗蒂憿舐內???嚗?
func _on_target_killed(data: Dictionary) -> void:
	if not _is_active:
		return
	var instance_id: String = data.get("instance_id", "")
	if not _mark_nodes.has(instance_id):
		return
	var node = _mark_nodes[instance_id]
	_mark_nodes.erase(instance_id)
	if not is_instance_valid(node):
		return
	# ?璅??格?嚗??脩??賊???
	var pos = node.position
	node.queue_free()
	_spawn_mark_kill_effect(pos)

## ?璅??格???閬箸???
func _spawn_mark_kill_effect(pos: Vector2) -> void:
	var lbl := Label.new()
	lbl.text = "?? ?%.0f" % _burst_mult
	lbl.position = pos + Vector2(-30, -20)
	lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 22)
	add_child(lbl)
	var tw = lbl.create_tween()
	tw.tween_property(lbl, "position:y", pos.y - 60, 0.5).set_ease(Tween.EASE_OUT)
	tw.tween_property(lbl, "modulate:a", 0.0, 0.3)
	tw.tween_callback(func(): if is_instance_valid(lbl): lbl.queue_free())

# ---- 鈭辣?? ----

func _on_lion_dance_burst(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	match phase:
		"burst_start":
			_handle_burst_start(data)
		"burst_end":
			_handle_burst_end(data)

func _handle_burst_start(data: Dictionary) -> void:
	var trigger_player: String = data.get("trigger_player", "")
	var trigger_name: String = data.get("trigger_name", "?拙振")
	_burst_mult = data.get("burst_mult", 3.0)
	_duration_sec = data.get("duration_sec", 15)
	_elapsed = 0.0
	_is_active = true
	_is_my_burst = (trigger_player == NetworkManager.get_player_id())

	# ?刻撟?蝝???
	_flash_screen(Color(1.0, 0.4, 0.0, 0.0), 0.35)

	# ?璈怠?
	_show_banner(trigger_name, _burst_mult)

	# 璅??格??
	var marked: Array = data.get("marked_targets", [])
	for t in marked:
		_add_mark_halo(t.get("instance_id", ""), t.get("x", 0.0), t.get("y", 0.0))

	# ?芸楛閫貊??銝剖亢憭??? 璅?敶歲
	if _is_my_burst:
		_show_center_lion()

func _handle_burst_end(data: Dictionary) -> void:
	_is_active = false

	# 皜???閮???
	for id in _mark_nodes:
		var node = _mark_nodes[id]
		if is_instance_valid(node):
			node.queue_free()
	_mark_nodes.clear()

	# 瘛∪璈怠?
	if is_instance_valid(_banner):
		var t = _banner.create_tween()
		t.tween_property(_banner, "modulate:a", 0.0, 0.4)
		t.tween_callback(func(): if is_instance_valid(_banner): _banner.queue_free(); _banner = null)

	# 皜?閮?
	if is_instance_valid(_countdown_lbl):
		_countdown_lbl.queue_free()
		_countdown_lbl = null

	# ?喳皛蝯?敶?
	var remaining: int = data.get("remaining_targets", 0)
	_show_result_panel(remaining)

# ---- UI 撱箇? ----

func _flash_screen(base_color: Color, peak_alpha: float) -> void:
	var flash := ColorRect.new()
	flash.size = Vector2(SCREEN_W, SCREEN_H)
	flash.color = Color(base_color.r, base_color.g, base_color.b, 0.0)
	add_child(flash)
	var tw = flash.create_tween()
	tw.tween_property(flash, "color:a", peak_alpha, 0.1)
	tw.tween_property(flash, "color:a", 0.0, 0.35)
	tw.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())

func _show_banner(trigger_name: String, mult: float) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()

	_banner = Node2D.new()
	add_child(_banner)

	# 璈怠??嚗?蝝撓撅歹?
	var bg := ColorRect.new()
	bg.size = Vector2(SCREEN_W, 56)
	bg.position = Vector2(0, -60)
	bg.color = Color(0.85, 0.25, 0.0, 0.92)
	_banner.add_child(bg)

	# 璈怠???
	var lbl := Label.new()
	lbl.text = "?? %s 閫貊?????潘?璅??格? ?%.0f ??嚗? % [trigger_name, mult]"
	lbl.position = Vector2(20, 10)
	lbl.add_theme_color_override("font_color", Color(1.0, 0.95, 0.6))
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 18)
	_banner.add_child(lbl)

	# ?閮? Label
	_countdown_lbl = Label.new()
	_countdown_lbl.text = "?? 15蝘?"
	_countdown_lbl.position = Vector2(SCREEN_W - 120, 10)
	_countdown_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	if _pixel_font:
		_countdown_lbl.add_theme_font_override("font", _pixel_font)
		_countdown_lbl.add_theme_font_size_override("font_size", 18)
	_banner.add_child(_countdown_lbl)

	# 璈怠?敺??冽???
	var tw = _banner.create_tween()
	tw.tween_property(bg, "position:y", 0.0, 0.25).set_ease(Tween.EASE_OUT)

func _add_mark_halo(instance_id: String, x: float, y: float) -> void:
	if instance_id == "":
		return

	var halo := Node2D.new()
	halo.position = Vector2(x, y)
	add_child(halo)
	_mark_nodes[instance_id] = halo

	# ??嚗?銵??恬?
	var ring := ColorRect.new()
	ring.size = Vector2(64, 64)
	ring.position = Vector2(-32, -32)
	ring.color = Color(1.0, 0.85, 0.0, 0.0)
	halo.add_child(ring)

	# ????嚗?摰 ring 蝭暺?蝭暺?斗??芸??迫嚗?
	var tw = ring.create_tween().set_loops()
	tw.tween_property(ring, "color:a", 0.7, 0.4)
	tw.tween_property(ring, "color:a", 0.2, 0.4)

	# ??璅惜
	var lbl := Label.new()
	lbl.text = "?%.0f" % _burst_mult
	lbl.position = Vector2(-20, -48)
	lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.0))
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 14)
	halo.add_child(lbl)

func _show_center_lion() -> void:
	var lbl := Label.new()
	lbl.text = "??"
	lbl.position = Vector2(SCREEN_W / 2 - 40, SCREEN_H / 2 - 60)
	lbl.add_theme_color_override("font_color", Color(1.0, 0.6, 0.0))
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 72)
	add_child(lbl)

	# 敶歲?
	var tw = lbl.create_tween()
	tw.tween_property(lbl, "position:y", SCREEN_H / 2 - 100, 0.15).set_ease(Tween.EASE_OUT)
	tw.tween_property(lbl, "position:y", SCREEN_H / 2 - 60, 0.12).set_ease(Tween.EASE_IN)
	tw.tween_property(lbl, "position:y", SCREEN_H / 2 - 80, 0.1).set_ease(Tween.EASE_OUT)
	tw.tween_property(lbl, "position:y", SCREEN_H / 2 - 60, 0.1).set_ease(Tween.EASE_IN)
	tw.tween_property(lbl, "modulate:a", 0.0, 0.5)
	tw.tween_callback(func(): if is_instance_valid(lbl): lbl.queue_free())

	# ?舀?憿?
	var sub := Label.new()
	sub.text = "敹怠?璅??格?嚗?"
	sub.position = Vector2(SCREEN_W / 2 - 100, SCREEN_H / 2 + 20)
	sub.add_theme_color_override("font_color", Color(1.0, 0.85, 0.3))
	if _pixel_font:
		sub.add_theme_font_override("font", _pixel_font)
		sub.add_theme_font_size_override("font_size", 20)
	add_child(sub)
	var sub_tw = sub.create_tween()
	sub_tw.tween_interval(2.0)
	sub_tw.tween_property(sub, "modulate:a", 0.0, 0.5)
	sub_tw.tween_callback(func(): if is_instance_valid(sub): sub.queue_free())

func _show_result_panel(remaining: int) -> void:
	var panel := Node2D.new()
	panel.position = Vector2(SCREEN_W + 10, SCREEN_H / 2 - 80)
	add_child(panel)

	# ?Ｘ?
	var bg := ColorRect.new()
	bg.size = Vector2(280, 160)
	bg.color = Color(0.15, 0.08, 0.0, 0.95)
	panel.add_child(bg)

	# ??
	var border := ColorRect.new()
	border.size = Vector2(280, 4)
	border.color = Color(1.0, 0.6, 0.0, 1.0)
	panel.add_child(border)

	# 璅?
	var title := Label.new()
	title.text = "?? ?????潛???"
	title.position = Vector2(10, 12)
	title.add_theme_color_override("font_color", Color(1.0, 0.85, 0.3))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 16)
	panel.add_child(title)

	# ??
	var mult_lbl := Label.new()
	mult_lbl.text = "???嚗?.0f" % _burst_mult
	mult_lbl.position = Vector2(10, 45)
	mult_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.5))
	if _pixel_font:
		mult_lbl.add_theme_font_override("font", _pixel_font)
		mult_lbl.add_theme_font_size_override("font_size", 14)
	panel.add_child(mult_lbl)

	# ?拚??芣???
	var remain_lbl := Label.new()
	remain_lbl.text = "?芣??湔?閮?%d ?? % remaining"
	remain_lbl.position = Vector2(10, 75)
	remain_lbl.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	if _pixel_font:
		remain_lbl.add_theme_font_override("font", _pixel_font)
		remain_lbl.add_theme_font_size_override("font_size", 14)
	panel.add_child(remain_lbl)

	# 敺?湔???
	var tw = panel.create_tween()
	tw.tween_property(panel, "position:x", SCREEN_W - 300, 0.3).set_ease(Tween.EASE_OUT)
	tw.tween_interval(3.0)
	tw.tween_property(panel, "modulate:a", 0.0, 0.5)
	tw.tween_callback(func(): if is_instance_valid(panel): panel.queue_free())

	# ??x ????
	if _burst_mult >= 7.0:
		_flash_screen(Color(1.0, 0.7, 0.0, 0.0), 0.5)
		await get_tree().create_timer(0.2).timeout
		_flash_screen(Color(1.0, 0.7, 0.0, 0.0), 0.4)

	# ??0x 敶抵銝???
	if _burst_mult >= 10.0:
		await get_tree().create_timer(0.4).timeout
		_flash_screen(Color(0.5, 0.0, 1.0, 0.0), 0.35)


# ---- 撣豢 ----
