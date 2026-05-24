# ChainExplosionPanel.gd ????????Ｘ嚗AY-088嚗?
# 憿舐內????蝑??璅?蜇?
# ??????恍銝剖亢憿舐內憭抒?
extends Control

# ???蝑?憿
const CHAIN_COLORS = {
	1: Color(1.0, 1.0, 1.0),    # 撠??嚗??
	2: Color(0.0, 0.75, 1.0),   # 銝剝??嚗予??
	3: Color(1.0, 0.84, 0.0),   # 憭折??嚗???
	4: Color(1.0, 0.27, 0.0),   # 頞????嚗?蝝?
}

const CHAIN_SIZES = {
	1: 24,
	2: 28,
	3: 34,
	4: 42,
}

func _ready():
	# ?? GameManager 閮?
	if GameManager.has_signal("chain_explosion"):
		GameManager.chain_explosion.connect(_on_chain_explosion)

func _on_chain_explosion(data: Dictionary) -> void:
	var level = data.get("level", 1)
	var level_name = data.get("level_name", "")"
	var chains = data.get("chains", [])
	var total_reward = data.get("total_reward", 0)
	var bonus_mult = data.get("bonus_mult", 1.0)
	var player_id = data.get("player_id", "")

	# ?芸?閫貊?拙振憿舐內閰喟敦?嚗隞摰園＊蝷箇陛??
	var is_self = (player_id == GameManager.player_data.get("player_id", ""))

	_show_chain_notify(level, level_name, len(chains), total_reward, bonus_mult, is_self)

	# 撠??◤????璅?曄??貊??
	for entry in chains:
		var instance_id = entry.get("instance_id", "")
		var mult = entry.get("multiplier", 1.0)
		# ? TargetManager ?剜?????寞?
		if GameManager.has_signal("chain_target_killed"):
			GameManager.emit_signal("chain_target_killed", instance_id, mult)

func _show_chain_notify(level: int, level_name: String, count: int, reward: int, bonus_mult: float, is_self: bool) -> void:
	var color = CHAIN_COLORS.get(level, Color.WHITE)
	var font_size = CHAIN_SIZES.get(level, 24)

	# 撱箇??蝭暺?
	var notify = Control.new()
	notify.z_index = 70
	add_child(notify)

	# ???蝑?璅惜嚗?Ｖ葉憭桀?銝?
	var name_lbl = Label.new()
	name_lbl.text = level_name
	name_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	name_lbl.add_theme_font_size_override("font_size", font_size)
	name_lbl.add_theme_color_override("font_color", color)
	name_lbl.add_theme_color_override("font_shadow_color", Color(0, 0, 0, 0.8))
	name_lbl.add_theme_constant_override("shadow_offset_x", 2)
	name_lbl.add_theme_constant_override("shadow_offset_y", 2)
	name_lbl.position = Vector2(540, 280)
	name_lbl.size = Vector2(200, 50)
	notify.add_child(name_lbl)

	# ?格??賊? + ?嚗撠撌梢＊蝷綽?
	if is_self and count > 0:
		var detail_lbl = Label.new()
		detail_lbl.text = "?%d ?格?  +%d" % [count, reward]
		detail_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		detail_lbl.add_theme_font_size_override("font_size", int(font_size * 0.6))
		detail_lbl.add_theme_color_override("font_color", color * 0.9)
		detail_lbl.position = Vector2(520, 280 + font_size + 4)
		detail_lbl.size = Vector2(240, 30)
		notify.add_child(detail_lbl)

		# ????嚗葉???隞乩??＊蝷綽?
		if bonus_mult > 1.0:
			var mult_lbl = Label.new()
			mult_lbl.text = "????? ?%.1f" % bonus_mult
			mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
			mult_lbl.add_theme_font_size_override("font_size", int(font_size * 0.5))
			mult_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
			mult_lbl.position = Vector2(520, 280 + font_size + 36)
			mult_lbl.size = Vector2(240, 24)
			notify.add_child(mult_lbl)

	# ?嚗葬?曉??????? ??銝宏瘛∪
	notify.scale = Vector2(0.3, 0.3)
	notify.modulate.a = 0.0

	var tween = notify.create_tween()
	tween.tween_property(notify, "scale", Vector2(1.0, 1.0), 0.15).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(notify, "modulate:a", 1.0, 0.15)
	tween.tween_interval(0.8 if level < 3 else 1.2)
	tween.tween_property(notify, "position:y", notify.position.y - 60, 0.4)
	tween.parallel().tween_property(notify, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func(): if is_instance_valid(notify): notify.queue_free())

	# 頞????嚗?憭?恍??
	if level >= 4:
		_spawn_mega_flash(color)

func _spawn_mega_flash(color: Color) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, 0.3)
	flash.z_index = 68
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())
