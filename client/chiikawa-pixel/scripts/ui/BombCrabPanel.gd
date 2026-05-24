## BombCrabPanel.gd
## ?詨??寥???Ｘ嚗AY-143嚗?
## 璆剔?靘?嚗oyal-fishing.uk 2026?orth 70x, explosive crustacean triggers multiple
## large-scale detonations. Each bomb creates expanding capture zones for massive multi-target eliminations.??
## 憿舐內?詨??寡孛?潛?銝郭?? + ?犖蝯?敶?

extends Control

var pixel_font: Font = null

# ?憿嚗?蝝蜓憿?
const COLOR_BOMB_ORANGE = Color(1.0, 0.45, 0.0)
const COLOR_BOMB_RED    = Color(1.0, 0.2, 0.1)
const COLOR_BOMB_YELLOW = Color(1.0, 0.85, 0.0)
const COLOR_BOMB_WHITE  = Color(1.0, 1.0, 1.0)

## ??????HUD.gd ?澆嚗?
func setup(font: Font) -> void:
	pixel_font = font
	GameManager.bomb_crab_chain.connect(_on_bomb_crab_chain)

## ???詨??寥??閮
func _on_bomb_crab_chain(data: Dictionary) -> void:
	var phase = data.get("phase", "")
	match phase:
		"bomb_start":
			_show_bomb_start(data)
		"explosion":
			_show_explosion_wave(data)
		"result":
			_show_result(data)

## ?詨??寡孛?潭帖撟?bomb_start ?挾嚗?
func _show_bomb_start(data: Dictionary) -> void:
	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# ?璈怠?
	var banner = Control.new()
	banner.name = "BombCrabBanner"
	banner.position = Vector2(0, -60)
	banner.size = Vector2(1280, 52)
	banner.z_index = 86
	canvas_layer.add_child(banner)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.6, 0.1, 0.0, 0.92)
	banner.add_child(bg)

	var lbl = Label.new()
	lbl.text = "? ?詨??寥??嚗?"
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 22)
	lbl.add_theme_color_override("font_color", COLOR_BOMB_YELLOW)
	if is_instance_valid(pixel_font):
		lbl.add_theme_font_override("font", pixel_font)
	banner.add_child(lbl)

	# 皛?
	var tween = banner.create_tween()
	tween.tween_property(banner, "position:y", 0.0, 0.3).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.tween_interval(1.8)
	tween.tween_property(banner, "position:y", -60.0, 0.25).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
	tween.tween_callback(banner.queue_free)

## ?瘜Ｙ??explosion ?挾嚗?
func _show_explosion_wave(data: Dictionary) -> void:
	var wave_index = data.get("wave_index", 0)
	var explode_x = data.get("trigger_x", 640.0)
	var explode_y = data.get("trigger_y", 360.0)

	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# ???嚗?Ｗ?璈???嚗撥摨虫?瘜Ｘ活??嚗?
	var flash_alpha = 0.35 - wave_index * 0.08
	if flash_alpha > 0.05:
		var flash = ColorRect.new()
		flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		flash.color = Color(COLOR_BOMB_ORANGE.r, COLOR_BOMB_ORANGE.g, COLOR_BOMB_ORANGE.b, 0.0)
		flash.z_index = 85
		canvas_layer.add_child(flash)

		var flash_tween = flash.create_tween()
		flash_tween.tween_property(flash, "color:a", flash_alpha, 0.06)
		flash_tween.tween_property(flash, "color:a", 0.0, 0.25)
		flash_tween.tween_callback(flash.queue_free)

	# ???嚗?銝剖?雿蔭嚗?
	_spawn_explosion_circle(explode_x, explode_y, wave_index, canvas_layer)

	# ?蝎?
	_spawn_explosion_particles(explode_x, explode_y, wave_index, canvas_layer)

	# 瘜Ｘ活璅內嚗洵撟暹郭嚗?
	var wave_lbl = Label.new()
	wave_lbl.text = "? WAVE %d" % (wave_index + 1)
	wave_lbl.position = Vector2(explode_x - 60, explode_y - 50)
	wave_lbl.size = Vector2(120, 30)
	wave_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	wave_lbl.add_theme_font_size_override("font_size", 14)
	wave_lbl.add_theme_color_override("font_color", COLOR_BOMB_YELLOW)
	if is_instance_valid(pixel_font):
		wave_lbl.add_theme_font_override("font", pixel_font)
	wave_lbl.z_index = 87
	canvas_layer.add_child(wave_lbl)

	var lbl_tween = wave_lbl.create_tween()
	lbl_tween.tween_property(wave_lbl, "position:y", explode_y - 80, 0.5).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)
	lbl_tween.parallel().tween_property(wave_lbl, "modulate:a", 0.0, 0.5)
	lbl_tween.tween_callback(wave_lbl.queue_free)

## ??????寞?
func _spawn_explosion_circle(cx: float, cy: float, wave: int, canvas_layer: Node) -> void:
	var circle = ColorRect.new()
	var radius = 60.0 + wave * 20.0
	circle.size = Vector2(radius * 2, radius * 2)
	circle.position = Vector2(cx - radius, cy - radius)
	circle.color = Color(COLOR_BOMB_ORANGE.r, COLOR_BOMB_ORANGE.g, COLOR_BOMB_ORANGE.b, 0.7)
	circle.z_index = 84
	canvas_layer.add_child(circle)

	var tween = circle.create_tween()
	tween.tween_property(circle, "size", Vector2(radius * 4, radius * 4), 0.3).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(circle, "position", Vector2(cx - radius * 2, cy - radius * 2), 0.3)
	tween.parallel().tween_property(circle, "color:a", 0.0, 0.3)
	tween.tween_callback(circle.queue_free)

## ???蝎?
func _spawn_explosion_particles(cx: float, cy: float, wave: int, canvas_layer: Node) -> void:
	var rng = RandomNumberGenerator.new()
	rng.randomize()
	var count = 12 + wave * 4

	for i in count:
		var particle = ColorRect.new()
		particle.size = Vector2(rng.randf_range(4, 10), rng.randf_range(4, 10))
		particle.position = Vector2(cx, cy)
		# 憿嚗?/蝝?暺璈?
		var colors = [COLOR_BOMB_ORANGE, COLOR_BOMB_RED, COLOR_BOMB_YELLOW]
		particle.color = colors[rng.randi() % 3]
		particle.z_index = 85
		canvas_layer.add_child(particle)

		var angle = rng.randf_range(0, TAU)
		var speed = rng.randf_range(80, 200 + wave * 40)
		var target_x = cx + cos(angle) * speed
		var target_y = cy + sin(angle) * speed
		var duration = rng.randf_range(0.3, 0.7)

		var tween = particle.create_tween()
		tween.tween_property(particle, "position", Vector2(target_x, target_y), duration).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)
		tween.parallel().tween_property(particle, "modulate:a", 0.0, duration)
		tween.tween_callback(particle.queue_free)

## 憿舐內蝯?敶?嚗esult ?挾嚗?
func _show_result(data: Dictionary) -> void:
	var total_reward = data.get("total_reward", 0)
	var killer_id = data.get("killer_id", "")
	var killer_name = data.get("killer_name", "")
	var killed_targets = data.get("killed_targets", [])
	var is_self = killer_id == NetworkManager.get_player_id()

	if total_reward <= 0:
		return

	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# 蝯?敶?嚗?湔??伐?
	var panel = Control.new()
	panel.name = "BombCrabResult"
	panel.position = Vector2(1280, 120)
	panel.size = Vector2(280, 160)
	panel.z_index = 88
	canvas_layer.add_child(panel)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.12, 0.04, 0.02, 0.92)
	panel.add_child(bg)

	# ??嚗?蝝嚗?
	var border_color = COLOR_BOMB_ORANGE if is_self else Color(0.8, 0.4, 0.1)
	for side in [Vector2(0, 0), Vector2(278, 0), Vector2(0, 158), Vector2(278, 158)]:
		var corner = ColorRect.new()
		corner.size = Vector2(2, 2)
		corner.position = side
		corner.color = border_color
		panel.add_child(corner)

	# 璅?
	var title_lbl = Label.new()
	title_lbl.text = "? ?詨??寧??賂?" if is_self else "? %s ?敶嚗? % killer_name"
	title_lbl.position = Vector2(8, 8)
	title_lbl.size = Vector2(264, 24)
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", COLOR_BOMB_YELLOW)
	if is_instance_valid(pixel_font):
		title_lbl.add_theme_font_override("font", pixel_font)
	panel.add_child(title_lbl)

	# ???
	var kill_lbl = Label.new()
	kill_lbl.text = "??葆?嚗?d ?璅? % len(killed_targets)"
	kill_lbl.position = Vector2(8, 36)
	kill_lbl.size = Vector2(264, 20)
	kill_lbl.add_theme_font_size_override("font_size", 11)
	kill_lbl.add_theme_color_override("font_color", COLOR_BOMB_WHITE)
	if is_instance_valid(pixel_font):
		kill_lbl.add_theme_font_override("font", pixel_font)
	panel.add_child(kill_lbl)

	# 瘜Ｘ活??
	var wave_counts = [0, 0, 0]
	for entry in killed_targets:
		var wi = entry.get("wave_index", 0)
		if wi < 3:
			wave_counts[wi] += 1
	var wave_lbl = Label.new()
	wave_lbl.text = "Wave1:%d  Wave2:%d  Wave3:%d" % [wave_counts[0], wave_counts[1], wave_counts[2]]
	wave_lbl.position = Vector2(8, 58)
	wave_lbl.size = Vector2(264, 18)
	wave_lbl.add_theme_font_size_override("font_size", 10)
	wave_lbl.add_theme_color_override("font_color", Color(0.9, 0.7, 0.5))
	if is_instance_valid(pixel_font):
		wave_lbl.add_theme_font_override("font", pixel_font)
	panel.add_child(wave_lbl)

	# ?
	var reward_lbl = Label.new()
	reward_lbl.text = "?? +%d" % total_reward
	reward_lbl.position = Vector2(8, 82)
	reward_lbl.size = Vector2(264, 32)
	reward_lbl.add_theme_font_size_override("font_size", 22)
	reward_lbl.add_theme_color_override("font_color", COLOR_BOMB_YELLOW)
	if is_instance_valid(pixel_font):
		reward_lbl.add_theme_font_override("font", pixel_font)
	panel.add_child(reward_lbl)

	# ?芸楛閫貊?????
	if is_self:
		var self_lbl = Label.new()
		self_lbl.text = "??雿孛?潔??詨??對?"
		self_lbl.position = Vector2(8, 118)
		self_lbl.size = Vector2(264, 18)
		self_lbl.add_theme_font_size_override("font_size", 10)
		self_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
		if is_instance_valid(pixel_font):
			self_lbl.add_theme_font_override("font", pixel_font)
		panel.add_child(self_lbl)

		# ???
		var flash = ColorRect.new()
		flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		flash.color = Color(1.0, 0.85, 0.0, 0.0)
		flash.z_index = 1
		panel.add_child(flash)
		var flash_tween = flash.create_tween()
		flash_tween.tween_property(flash, "color:a", 0.4, 0.1)
		flash_tween.tween_property(flash, "color:a", 0.0, 0.3)
		flash_tween.tween_callback(flash.queue_free)

	# 皛?
	var tween = panel.create_tween()
	tween.tween_property(panel, "position:x", 990.0, 0.35).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.tween_interval(3.5)
	tween.tween_property(panel, "position:x", 1280.0, 0.3).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
	tween.tween_callback(panel.queue_free)

	# ?剜?單?
	if is_self and AudioManager != null:
		AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
