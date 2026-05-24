## AnglerfishPanel.gd
## 撌典?擙?擳?窄蝞梢?選?DAY-145嚗?
## 璆剔?靘?嚗iligames.com 2026?iant Anglerfish can shoot electricity to open treasure chests??
## ? T109 敺孛?潮???餅??喳??圈?餈?撖嗥拳?格?嚗撥?園??窄蝞梁敺?憭???

extends Control

var pixel_font: Font = null

# ?餅?憿嚗??賡瘚蜓憿?
const COLOR_ELECTRIC_BLUE  = Color(0.0, 0.75, 1.0)
const COLOR_ELECTRIC_WHITE = Color(0.8, 0.95, 1.0)
const COLOR_ELECTRIC_CYAN  = Color(0.0, 1.0, 0.9)
const COLOR_GOLD           = Color(1.0, 0.85, 0.0)

## ??????HUD.gd ?澆嚗?
func setup(font: Font) -> void:
	pixel_font = font
	GameManager.anglerfish_shock.connect(_on_anglerfish_shock)

## ??擙?擳????
func _on_anglerfish_shock(data: Dictionary) -> void:
	var phase = data.get("phase", "")
	match phase:
		"shock_start":
			_show_shock_start(data)
		"result":
			_show_result(data)

## ?餅???嚗hock_start ?挾嚗?
func _show_shock_start(data: Dictionary) -> void:
	var trigger_x = data.get("trigger_x", 640.0)
	var trigger_y = data.get("trigger_y", 360.0)
	var chest_ids = data.get("chest_ids", [])

	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# ?璈怠?
	var banner = Control.new()
	banner.name = "AnglerfishBanner"
	banner.position = Vector2(0, -52)
	banner.size = Vector2(1280, 48)
	banner.z_index = 86
	canvas_layer.add_child(banner)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.0, 0.1, 0.3, 0.92)
	banner.add_child(bg)

	var lbl = Label.new()
	if len(chest_ids) > 0:
		lbl.text = "??擙?擳???? %d ?窄蝞梧?" % len(chest_ids)
	else:
		lbl.text = "??擙?擳??"
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 20)
	lbl.add_theme_color_override("font_color", COLOR_ELECTRIC_CYAN)
	if is_instance_valid(pixel_font):
		lbl.add_theme_font_override("font", pixel_font)
	banner.add_child(lbl)

	# 皛?
	var tween = banner.create_tween()
	tween.tween_property(banner, "position:y", 0.0, 0.25).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.tween_interval(1.5)
	tween.tween_property(banner, "position:y", -52.0, 0.2).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
	tween.tween_callback(banner.queue_free)

	# ?餅???嚗??質嚗?
	var flash = ColorRect.new()
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(COLOR_ELECTRIC_BLUE.r, COLOR_ELECTRIC_BLUE.g, COLOR_ELECTRIC_BLUE.b, 0.0)
	flash.z_index = 85
	canvas_layer.add_child(flash)

	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", 0.3, 0.05)
	flash_tween.tween_property(flash, "color:a", 0.0, 0.2)
	flash_tween.tween_callback(flash.queue_free)

	# ?餅?蝎?嚗?閫貊暺?憭???
	_spawn_electric_particles(trigger_x, trigger_y, canvas_layer)

## ???餅?蝎?
func _spawn_electric_particles(cx: float, cy: float, canvas_layer: Node) -> void:
	var rng = RandomNumberGenerator.new()
	rng.randomize()

	for i in 20:
		var particle = ColorRect.new()
		particle.size = Vector2(rng.randf_range(2, 6), rng.randf_range(2, 6))
		particle.position = Vector2(cx, cy)
		var colors = [COLOR_ELECTRIC_BLUE, COLOR_ELECTRIC_WHITE, COLOR_ELECTRIC_CYAN]
		particle.color = colors[rng.randi() % 3]
		particle.z_index = 86
		canvas_layer.add_child(particle)

		var angle = rng.randf_range(0, TAU)
		var speed = rng.randf_range(100, 300)
		var target_x = cx + cos(angle) * speed
		var target_y = cy + sin(angle) * speed
		var duration = rng.randf_range(0.2, 0.5)

		var tween = particle.create_tween()
		tween.tween_property(particle, "position", Vector2(target_x, target_y), duration).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)
		tween.parallel().tween_property(particle, "modulate:a", 0.0, duration)
		tween.tween_callback(particle.queue_free)

## 憿舐內蝯?敶?嚗esult ?挾嚗?
func _show_result(data: Dictionary) -> void:
	var total_reward = data.get("total_reward", 0)
	var killer_id = data.get("killer_id", "")
	var killer_name = data.get("killer_name", "")
	var opened_chests = data.get("opened_chests", [])
	var is_self = killer_id == NetworkManager.get_player_id()

	if total_reward <= 0:
		return

	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# 蝯?敶?嚗?湔??伐?
	var panel = Control.new()
	panel.name = "AnglerfishResult"
	panel.position = Vector2(1280, 100)
	panel.size = Vector2(270, 150)
	panel.z_index = 87
	canvas_layer.add_child(panel)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.0, 0.05, 0.15, 0.92)
	panel.add_child(bg)

	# ??嚗?嚗?
	var border_color = COLOR_ELECTRIC_CYAN if is_self else COLOR_ELECTRIC_BLUE
	var top_border = ColorRect.new()
	top_border.size = Vector2(270, 2)
	top_border.position = Vector2(0, 0)
	top_border.color = border_color
	panel.add_child(top_border)

	# 璅?
	var title_lbl = Label.new()
	title_lbl.text = "???餅??窄蝞梧?" if is_self else "??%s ???" % killer_name
	title_lbl.position = Vector2(8, 8)
	title_lbl.size = Vector2(254, 22)
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", COLOR_ELECTRIC_CYAN)
	if is_instance_valid(pixel_font):
		title_lbl.add_theme_font_override("font", pixel_font)
	panel.add_child(title_lbl)

	# ??撖嗥拳??
	var chest_lbl = Label.new()
	chest_lbl.text = "??撖嗥拳嚗?d ?? % len(opened_chests)"
	chest_lbl.position = Vector2(8, 34)
	chest_lbl.size = Vector2(254, 18)
	chest_lbl.add_theme_font_size_override("font_size", 11)
	chest_lbl.add_theme_color_override("font_color", COLOR_ELECTRIC_WHITE)
	if is_instance_valid(pixel_font):
		chest_lbl.add_theme_font_override("font", pixel_font)
	panel.add_child(chest_lbl)

	# 撖嗥拳???”嚗?憭＊蝷?3 ??
	var display_count = min(len(opened_chests), 3)
	for i in display_count:
		var chest = opened_chests[i]
		var chest_item = Label.new()
		chest_item.text = "  ? ?%.0f ??+%d" % [chest.get("multiplier", 0), chest.get("reward", 0)]
		chest_item.position = Vector2(8, 54 + i * 16)
		chest_item.size = Vector2(254, 16)
		chest_item.add_theme_font_size_override("font_size", 10)
		chest_item.add_theme_color_override("font_color", COLOR_GOLD)
		if is_instance_valid(pixel_font):
			chest_item.add_theme_font_override("font", pixel_font)
		panel.add_child(chest_item)

	# 蝮賜???
	var reward_lbl = Label.new()
	reward_lbl.text = "?? +%d" % total_reward
	reward_lbl.position = Vector2(8, 108)
	reward_lbl.size = Vector2(254, 30)
	reward_lbl.add_theme_font_size_override("font_size", 20)
	reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	if is_instance_valid(pixel_font):
		reward_lbl.add_theme_font_override("font", pixel_font)
	panel.add_child(reward_lbl)

	# ?芸楛閫貊???餉???
	if is_self:
		var flash = ColorRect.new()
		flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		flash.color = Color(COLOR_ELECTRIC_BLUE.r, COLOR_ELECTRIC_BLUE.g, COLOR_ELECTRIC_BLUE.b, 0.0)
		flash.z_index = 1
		panel.add_child(flash)
		var flash_tween = flash.create_tween()
		flash_tween.tween_property(flash, "color:a", 0.5, 0.08)
		flash_tween.tween_property(flash, "color:a", 0.0, 0.3)
		flash_tween.tween_callback(flash.queue_free)

	# 皛?
	var tween = panel.create_tween()
	tween.tween_property(panel, "position:x", 1000.0, 0.3).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.tween_interval(3.5)
	tween.tween_property(panel, "position:x", 1280.0, 0.25).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
	tween.tween_callback(panel.queue_free)
