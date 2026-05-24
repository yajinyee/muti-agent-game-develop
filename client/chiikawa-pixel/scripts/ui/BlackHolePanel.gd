## BlackHolePanel.gd ??暺?瞍拇蒂甇血閬死???Ｘ嚗AY-166嚗?
## 璆剔?靘?嚗?
##   - Ocean King 3 2026 Vortex 璈 ???曄蔭敺撘?璅?銝剖?蝘餃?嚗?蝯??豢???
##   - Black Hole Fishing 2026嚗team嚗??券?瘣擳??詨??拇?嚗?026 撟湔??啗隅??
## 閬死閮剛?嚗?
##   - black_hole_place嚗?曄蔭雿蔭憿舐內蝝怨瞍拇蒂? + ?冽?璈怠??XX ?曄蔭鈭?瘣???
##   - black_hole_suck嚗憬皜行憭?+ ?詨閮?剁??迤?典??N ?璅?..??
##   - result嚗?Ｗ?蝝怨??? + ?喳皛蝯?敶?嚗?交/????嚗?
##   - ?芸楛?曄蔭??銝剖亢憭??? 璅?敶歲?
extends Node2D

# ---- 撣豢 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0

# ---- ???----
var _pixel_font: Font = null
var _vortex_node: Node2D = null  # 瞍拇蒂閬死蝭暺?
var _banner_node: Node2D = null  # ?璈怠?
var _result_panel: Node2D = null # 蝯??Ｘ

# ---- ????----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("black_hole_result"):
		GameManager.black_hole_result.connect(_on_black_hole_result)

# ---- 鈭辣?? ----

func _on_black_hole_result(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var shooter_id: String = data.get("shooter_id", "")
	var shooter_name: String = data.get("shooter_name", "?拙振")
	var cx: float = data.get("center_x", SCREEN_W / 2.0)
	var cy: float = data.get("center_y", SCREEN_H / 2.0)
	var sucked_count: int = data.get("sucked_count", 0)
	var total_reward: int = data.get("total_reward", 0)
	var cost: int = data.get("cost", 0)
	var is_self: bool = (shooter_id == NetworkManager.get_player_id())

	match phase:
		"black_hole_place":
			_show_vortex(cx, cy, shooter_name, is_self)
		"black_hole_suck":
			_show_suck_effect(cx, cy, sucked_count)
		"result":
			_show_result(cx, cy, sucked_count, data.get("hit_targets", []), total_reward, cost, is_self)

# ---- 暺??曄蔭閬死 ----

func _show_vortex(cx: float, cy: float, shooter_name: String, is_self: bool) -> void:
	# 皜??瞍拇蒂
	if is_instance_valid(_vortex_node):
		_vortex_node.queue_free()

	_vortex_node = Node2D.new()
	_vortex_node.position = Vector2(cx, cy)
	add_child(_vortex_node)

	# 瞍拇蒂憭?嚗換?脣??堆?
	var outer_ring := ColorRect.new()
	outer_ring.size = Vector2(80, 80)
	outer_ring.position = Vector2(-40, -40)
	outer_ring.color = Color(0.4, 0.0, 0.8, 0.6)
	_vortex_node.add_child(outer_ring)

	# 瞍拇蒂銝剖?嚗楛蝝怨嚗?
	var inner := ColorRect.new()
	inner.size = Vector2(30, 30)
	inner.position = Vector2(-15, -15)
	inner.color = Color(0.15, 0.0, 0.35, 0.9)
	_vortex_node.add_child(inner)

	# 瞍拇蒂?內
	var vortex_lbl := Label.new()
	vortex_lbl.text = "??"
	vortex_lbl.position = Vector2(-18, -22)
	vortex_lbl.add_theme_font_size_override("font_size", 36)
	_vortex_node.add_child(vortex_lbl)

	# ???
	var spin_tween = _vortex_node.create_tween().set_loops()
	spin_tween.tween_property(_vortex_node, "rotation_degrees", 360.0, 1.5)

	# ?璈怠?
	_show_banner("?? %s ?曄蔭鈭?瘣憬皜佗?" % shooter_name, Color(0.6, 0.2, 1.0))

	# ?芸楛?曄蔭??銝剖亢憭扳?隤?
	if is_self:
		_show_self_vortex_effect()

func _show_self_vortex_effect() -> void:
	# ?刻撟換?脤???
	var flash := ColorRect.new()
	flash.size = Vector2(SCREEN_W, SCREEN_H)
	flash.position = Vector2(-position.x, -position.y)
	flash.color = Color(0.3, 0.0, 0.6, 0.0)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "color:a", 0.35, 0.15)
	tween.tween_property(flash, "color:a", 0.0, 0.3)
	tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())

	# 銝剖亢憭??? 璅?敶歲
	var big_lbl := Label.new()
	big_lbl.text = "?? 暺?瞍拇蒂嚗?"
	big_lbl.position = Vector2(SCREEN_W / 2.0 - 100, SCREEN_H / 2.0 - 30)
	big_lbl.size = Vector2(200, 60)
	big_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	big_lbl.add_theme_color_override("font_color", Color(0.8, 0.4, 1.0))
	big_lbl.add_theme_font_size_override("font_size", 28)
	if _pixel_font:
		big_lbl.add_theme_font_override("font", _pixel_font)
	add_child(big_lbl)

	var bounce_tween = big_lbl.create_tween()
	bounce_tween.tween_property(big_lbl, "scale", Vector2(1.3, 1.3), 0.15)
	bounce_tween.tween_property(big_lbl, "scale", Vector2(1.0, 1.0), 0.15)
	bounce_tween.tween_interval(1.0)
	bounce_tween.tween_property(big_lbl, "modulate:a", 0.0, 0.4)
	bounce_tween.tween_callback(func(): if is_instance_valid(big_lbl): big_lbl.queue_free())

# ---- ?詨?? ----

func _show_suck_effect(cx: float, cy: float, sucked_count: int) -> void:
	if not is_instance_valid(_vortex_node):
		return

	# 瞍拇蒂?游之?
	var expand_tween = _vortex_node.create_tween()
	expand_tween.tween_property(_vortex_node, "scale", Vector2(1.5, 1.5), 0.3)
	expand_tween.tween_property(_vortex_node, "scale", Vector2(1.2, 1.2), 0.2)

	# ?詨閮??
	var suck_lbl := Label.new()
	suck_lbl.text = "甇??詨 %d ?璅?.." % sucked_count
	suck_lbl.position = Vector2(cx - 80, cy + 50)
	suck_lbl.size = Vector2(160, 20)
	suck_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	suck_lbl.add_theme_color_override("font_color", Color(0.8, 0.5, 1.0))
	suck_lbl.add_theme_font_size_override("font_size", 11)
	if _pixel_font:
		suck_lbl.add_theme_font_override("font", _pixel_font)
	add_child(suck_lbl)

	var tween = suck_lbl.create_tween()
	tween.tween_interval(1.2)
	tween.tween_property(suck_lbl, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func(): if is_instance_valid(suck_lbl): suck_lbl.queue_free())

# ---- ?蝯? ----

func _show_result(cx: float, cy: float, sucked_count: int, hit_targets: Array, total_reward: int, cost: int, is_self: bool) -> void:
	# 皜瞍拇蒂
	if is_instance_valid(_vortex_node):
		var fade_tween = _vortex_node.create_tween()
		fade_tween.tween_property(_vortex_node, "modulate:a", 0.0, 0.3)
		fade_tween.tween_callback(func(): if is_instance_valid(_vortex_node): _vortex_node.queue_free())

	# ?刻撟??賊???蝝怨嚗?
	var flash := ColorRect.new()
	flash.size = Vector2(SCREEN_W, SCREEN_H)
	flash.position = Vector2(-position.x, -position.y)
	flash.color = Color(0.4, 0.0, 0.8, 0.0)
	add_child(flash)

	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", 0.5, 0.1)
	flash_tween.tween_property(flash, "color:a", 0.0, 0.4)
	flash_tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())

	# ???嚗暺?雿蔭嚗?
	var explosion_lbl := Label.new()
	explosion_lbl.text = "?"
	explosion_lbl.position = Vector2(cx - 24, cy - 24)
	explosion_lbl.add_theme_font_size_override("font_size", 48)
	add_child(explosion_lbl)

	var exp_tween = explosion_lbl.create_tween()
	exp_tween.tween_property(explosion_lbl, "scale", Vector2(2.0, 2.0), 0.3)
	exp_tween.parallel().tween_property(explosion_lbl, "modulate:a", 0.0, 0.3)
	exp_tween.tween_callback(func(): if is_instance_valid(explosion_lbl): explosion_lbl.queue_free())

	# 皜璈怠?
	if is_instance_valid(_banner_node):
		var banner_tween = _banner_node.create_tween()
		banner_tween.tween_interval(0.5)
		banner_tween.tween_property(_banner_node, "position:x", SCREEN_W + 10, 0.4)
		banner_tween.tween_callback(func(): if is_instance_valid(_banner_node): _banner_node.queue_free())

	if total_reward <= 0:
		return

	# ?喳皛蝯??Ｘ
	var net_reward = total_reward - cost
	var kill_count = hit_targets.size()

	if is_instance_valid(_result_panel):
		_result_panel.queue_free()

	_result_panel = Node2D.new()
	_result_panel.position = Vector2(SCREEN_W + 10, SCREEN_H / 2.0 - 70)
	add_child(_result_panel)

	var panel_bg := ColorRect.new()
	panel_bg.size = Vector2(200, 140)
	panel_bg.color = Color(0.06, 0.0, 0.12, 0.95)
	_result_panel.add_child(panel_bg)

	# ??
	var border := ColorRect.new()
	border.size = Vector2(202, 142)
	border.position = Vector2(-1, -1)
	border.color = Color(0.5, 0.1, 0.9, 0.8)
	border.z_index = -1
	_result_panel.add_child(border)

	# 璅?
	var title_lbl := Label.new()
	title_lbl.text = "?? 暺?瞍拇蒂?嚗?"
	title_lbl.position = Vector2(4, 6)
	title_lbl.size = Vector2(192, 20)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_color_override("font_color", Color(0.8, 0.4, 1.0))
	title_lbl.add_theme_font_size_override("font_size", 12)
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
	_result_panel.add_child(title_lbl)

	# ?詨??
	var suck_lbl := Label.new()
	suck_lbl.text = "?詨?格?嚗?d ?? % sucked_count"
	suck_lbl.position = Vector2(8, 32)
	suck_lbl.size = Vector2(184, 18)
	suck_lbl.add_theme_color_override("font_color", Color(0.7, 0.5, 1.0))
	suck_lbl.add_theme_font_size_override("font_size", 11)
	if _pixel_font:
		suck_lbl.add_theme_font_override("font", _pixel_font)
	_result_panel.add_child(suck_lbl)

	# ???
	var kill_lbl := Label.new()
	kill_lbl.text = "??格?嚗?d ?? % kill_count"
	kill_lbl.position = Vector2(8, 52)
	kill_lbl.size = Vector2(184, 18)
	kill_lbl.add_theme_color_override("font_color", Color(0.9, 0.6, 1.0))
	kill_lbl.add_theme_font_size_override("font_size", 11)
	if _pixel_font:
		kill_lbl.add_theme_font_override("font", _pixel_font)
	_result_panel.add_child(kill_lbl)

	# 鞎餌
	var cost_lbl := Label.new()
	cost_lbl.text = "鞎餌嚗?%d" % cost
	cost_lbl.position = Vector2(8, 72)
	cost_lbl.size = Vector2(184, 18)
	cost_lbl.add_theme_color_override("font_color", Color(0.8, 0.4, 0.4))
	cost_lbl.add_theme_font_size_override("font_size", 11)
	if _pixel_font:
		cost_lbl.add_theme_font_override("font", _pixel_font)
	_result_panel.add_child(cost_lbl)

	# ?
	var reward_lbl := Label.new()
	reward_lbl.text = "?嚗?%d" % total_reward
	reward_lbl.position = Vector2(8, 92)
	reward_lbl.size = Vector2(184, 18)
	reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	reward_lbl.add_theme_font_size_override("font_size", 11)
	if _pixel_font:
		reward_lbl.add_theme_font_override("font", _pixel_font)
	_result_panel.add_child(reward_lbl)

	# 瘛冽??
	var net_lbl := Label.new()
	if net_reward > 0:
		net_lbl.text = "瘛冽??+%d ?? % net_reward"
		net_lbl.add_theme_color_override("font_color", Color(0.3, 1.0, 0.5))
	else:
		net_lbl.text = "瘛冽??%d" % net_reward
		net_lbl.add_theme_color_override("font_color", Color(0.8, 0.4, 0.4))
	net_lbl.position = Vector2(8, 112)
	net_lbl.size = Vector2(184, 18)
	net_lbl.add_theme_font_size_override("font_size", 12)
	if _pixel_font:
		net_lbl.add_theme_font_override("font", _pixel_font)
	_result_panel.add_child(net_lbl)

	# 皛?
	var slide_tween = _result_panel.create_tween()
	slide_tween.tween_property(_result_panel, "position:x", SCREEN_W - 210.0, 0.4)

	# 擃??湔嚗???
	if kill_count >= 5:
		_show_double_flash()

	# 4 蝘?皛
	slide_tween.tween_interval(4.0)
	slide_tween.tween_property(_result_panel, "position:x", SCREEN_W + 10.0, 0.4)
	slide_tween.tween_callback(func(): if is_instance_valid(_result_panel): _result_panel.queue_free())

# ---- ?璈怠? ----

func _show_banner(text: String, color: Color) -> void:
	if is_instance_valid(_banner_node):
		_banner_node.queue_free()

	_banner_node = Node2D.new()
	_banner_node.position = Vector2(-SCREEN_W, 10)
	add_child(_banner_node)

	var bg := ColorRect.new()
	bg.size = Vector2(SCREEN_W, 28)
	bg.color = Color(0.06, 0.0, 0.12, 0.92)
	_banner_node.add_child(bg)

	var lbl := Label.new()
	lbl.text = text
	lbl.size = Vector2(SCREEN_W, 28)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 14)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
	_banner_node.add_child(lbl)

	# 皛?
	var tween = _banner_node.create_tween()
	tween.tween_property(_banner_node, "position:x", 0.0, 0.35)

# ---- ????擃??湔嚗?---

func _show_double_flash() -> void:
	for i in range(2):
		var flash := ColorRect.new()
		flash.size = Vector2(SCREEN_W, SCREEN_H)
		flash.position = Vector2(-position.x, -position.y)
		flash.color = Color(0.5, 0.1, 1.0, 0.0)
		add_child(flash)

		var tween = flash.create_tween()
		tween.tween_interval(float(i) * 0.25)
		tween.tween_property(flash, "color:a", 0.4, 0.1)
		tween.tween_property(flash, "color:a", 0.0, 0.2)
		tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())
