## ElectricJellyfishPanel.gd ???餅?瘞湔??餅?蝬脰楝?Ｘ嚗AY-193嚗?
## 璆剔?靘?嚗ing of Ocean 2026?lectric jellyfish chains current between adjacent targets,
## paying multipliers from every link in the chain??
## 閬死銝駁?嚗??脤瘚?+ ?餅???蝺???+ 蝬脰楝?閬死

extends Control

const ELECTRIC_COLOR := Color(0.0, 1.0, 1.0)   # ?嚗瘚?嚗?
const KILL_COLOR     := Color(1.0, 1.0, 0.0)   # 暺嚗??湛?
const MISS_COLOR     := Color(0.3, 0.8, 1.0)   # 瘛∟?嚗?嚗?

var _banner: Control = null
var _link_counter: Label = null
var _total_links: int = 0
var _total_kills: int = 0

func _ready() -> void:
	if GameManager.has_signal("electric_jellyfish"):
		GameManager.electric_jellyfish.connect(_on_electric_jellyfish)

func _on_electric_jellyfish(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	if phase == "network_start":
		_show_network_start(data)
	elif phase == "network_result":
		_show_network_result(data)
	elif phase.begins_with("link_"):
		_show_link(data)

# ?? network_start ??????????????????????????????????????????????????????????????
func _show_network_start(data: Dictionary) -> void:
	var killer_name: String = data.get("killer_name", "")
	var link_count: int = data.get("link_count", 0)
	_total_links = link_count
	_total_kills = 0

	# ???嚗甈∴?
	_flash_screen(ELECTRIC_COLOR, 0.5)
	var t1 = get_tree().create_timer(0.15)
	t1.timeout.connect(func(): _flash_screen(ELECTRIC_COLOR, 0.35))

	# ?璈怠?
	_show_banner("?♀?%s ?瘚偌瘥?撱箇? %d 璇瘚?嚗? % [killer_name, link_count], ELECTRIC_COLOR)"

	# ??閮??
	_show_link_counter(link_count)

# ?? link_N ?餅??? ?????????????????????????????????????????????????????????????
func _show_link(data: Dictionary) -> void:
	var xa: float = data.get("x_a", 0.0)
	var ya: float = data.get("y_a", 0.0)
	var xb: float = data.get("x_b", 0.0)
	var yb: float = data.get("y_b", 0.0)
	var is_kill: bool = data.get("is_kill", false)
	var reward: int = data.get("reward", 0)
	var link_index: int = data.get("link_index", 0)

	var line_color := KILL_COLOR if is_kill else MISS_COLOR

	# 蝜芾ˊ?餅???蝺?
	_draw_electric_line(xa, ya, xb, yb, line_color)

	# ??＊蝷箇??菜筑??摮?
	if is_kill and reward > 0:
		var mid_x := (xa + xb) / 2.0
		var mid_y := (ya + yb) / 2.0
		_show_floating_reward("+%d" % reward, mid_x, mid_y, KILL_COLOR)
		_total_kills += 1

	# ?湔閮??
	if _link_counter and is_instance_valid(_link_counter):
		_link_counter.text = "??%d/%d ?? | ? %d" % [link_index, _total_links, _total_kills]

# ?? network_result ?餅?蝬脰楝蝯? ?????????????????????????????????????????????????
func _show_network_result(data: Dictionary) -> void:
	var total_kills: int = data.get("total_kills", 0)
	var total_reward: int = data.get("total_reward", 0)
	var link_count: int = data.get("link_count", 0)

	# 瘛∪璈怠?
	_hide_banner()

	# ?喳皛蝯?敶?
	_show_result_popup(link_count, total_kills, total_reward)

	# 瘛∪閮??
	if _link_counter and is_instance_valid(_link_counter):
		var tween = create_tween()
		tween.tween_property(_link_counter, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_link_counter.queue_free)
		_link_counter = null

	# 憭折?????憭???
	if link_count >= 12:
		_flash_screen(ELECTRIC_COLOR, 0.6)
	elif total_kills >= 5:
		_flash_screen(KILL_COLOR, 0.45)

# ?? 頛?賣 ????????????????????????????????????????????????????????????????????

func _flash_screen(color: Color, alpha: float) -> void:
	var overlay = ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, alpha)
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tween = create_tween()
	tween.tween_property(overlay, "modulate:a", 0.0, 0.3)
	tween.tween_callback(overlay.queue_free)

func _show_banner(text: String, color: Color) -> void:
	if _banner != null and is_instance_valid(_banner):
		_banner.queue_free()

	_banner = Control.new()
	_banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 50)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.color = Color(0.0, 0.1, 0.15, 0.9)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_banner.add_child(bg)

	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", 19)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.add_child(label)

	_banner.position.y = -50
	var tween = create_tween()
	tween.tween_property(_banner, "position:y", 0.0, 0.25).set_trans(Tween.TRANS_BACK)

func _hide_banner() -> void:
	if _banner != null and is_instance_valid(_banner):
		var tween = create_tween()
		tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
		tween.tween_callback(_banner.queue_free)
		_banner = null

func _show_link_counter(link_count: int) -> void:
	if _link_counter != null and is_instance_valid(_link_counter):
		_link_counter.queue_free()

	_link_counter = Label.new()
	_link_counter.text = "??0/%d ?? | ? 0" % link_count
	_link_counter.add_theme_color_override("font_color", ELECTRIC_COLOR)
	_link_counter.add_theme_font_size_override("font_size", 16)
	_link_counter.set_anchors_preset(Control.PRESET_BOTTOM_RIGHT)
	_link_counter.offset_left = -260
	_link_counter.offset_top = -42
	_link_counter.offset_right = -10
	_link_counter.offset_bottom = -10
	add_child(_link_counter)

func _draw_electric_line(x1: float, y1: float, x2: float, y2: float, color: Color) -> void:
	# ?典??? ColorRect 璅⊥?餅?蝺??賊??嚗?
	var steps := 8
	var prev_x := x1
	var prev_y := y1
	var dx := (x2 - x1) / steps
	var dy := (y2 - y1) / steps

	for i in range(1, steps + 1):
		var nx := x1 + dx * i + randf_range(-6, 6)  # ?賊??宏
		var ny := y1 + dy * i + randf_range(-6, 6)

		# 蝜芾ˊ蝺挾嚗蝝圈 ColorRect 餈撮嚗?
		var seg_len := sqrt((nx - prev_x) * (nx - prev_x) + (ny - prev_y) * (ny - prev_y))
		if seg_len < 1.0:
			prev_x = nx
			prev_y = ny
			continue

		var seg = ColorRect.new()
		seg.color = Color(color.r, color.g, color.b, 0.8)
		seg.size = Vector2(seg_len, 3)
		seg.position = Vector2(prev_x, prev_y - 1.5)

		# ??蝺挾
		var angle := atan2(ny - prev_y, nx - prev_x)
		seg.rotation = angle
		seg.pivot_offset = Vector2(0, 1.5)

		add_child(seg)

		# 瘛∪?
		var tween = create_tween()
		tween.tween_property(seg, "modulate:a", 0.0, 0.6)
		tween.tween_callback(seg.queue_free)

		prev_x = nx
		prev_y = ny

	# ?典蝡舫＊蝷粹瘚?暺?撠?暺?
	_spawn_node_dot(x1, y1, color)
	_spawn_node_dot(x2, y2, color)

func _spawn_node_dot(x: float, y: float, color: Color) -> void:
	var dot = ColorRect.new()
	dot.color = Color(color.r, color.g, color.b, 0.9)
	dot.size = Vector2(8, 8)
	dot.position = Vector2(x - 4, y - 4)
	add_child(dot)
	var tween = create_tween()
	tween.tween_property(dot, "modulate:a", 0.0, 0.5)
	tween.tween_callback(dot.queue_free)

func _show_floating_reward(text: String, x: float, y: float, color: Color) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", 15)
	label.position = Vector2(x - 20, y - 10)
	add_child(label)

	var tween = create_tween()
	tween.set_parallel(true)
	tween.tween_property(label, "position:y", y - 45, 0.7)
	tween.tween_property(label, "modulate:a", 0.0, 0.7)
	tween.chain().tween_callback(label.queue_free)

func _show_result_popup(link_count: int, total_kills: int, total_reward: int) -> void:
	var popup = Control.new()
	popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	popup.custom_minimum_size = Vector2(210, 120)
	popup.offset_left = -220
	popup.offset_top = -60
	popup.offset_right = -10
	popup.offset_bottom = 60
	add_child(popup)

	var bg = ColorRect.new()
	bg.color = Color(0.0, 0.08, 0.12, 0.92)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	popup.add_child(bg)

	var vbox = VBoxContainer.new()
	vbox.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	vbox.add_theme_constant_override("separation", 6)
	popup.add_child(vbox)

	_add_label(vbox, "?♀??餅?蝬脰楝蝯?嚗?, ELECTRIC_COLOR, 16)"
	_add_label(vbox, "?餅???嚗?d 璇? % link_count, Color.WHITE, 14)"
	_add_label(vbox, "??格?嚗?d ?? % total_kills, KILL_COLOR, 14)"
	_add_label(vbox, "蝮賜??蛛?%d ?馳" % total_reward, ELECTRIC_COLOR, 14)

	popup.position.x += 230
	var tween = create_tween()
	tween.tween_property(popup, "position:x", popup.position.x - 230, 0.3).set_trans(Tween.TRANS_BACK)

	var timer = get_tree().create_timer(4.5)
	timer.timeout.connect(func():
		if is_instance_valid(popup):
			var t2 = create_tween()
			t2.tween_property(popup, "modulate:a", 0.0, 0.5)
			t2.tween_callback(popup.queue_free)
	)

func _add_label(parent: Control, text: String, color: Color, size: int) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", size)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	parent.add_child(label)
