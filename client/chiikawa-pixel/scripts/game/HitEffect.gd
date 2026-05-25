## HitEffect.gd — 命中特效、擊破粒子、獎勵跳字
## hit-effect-agent 負責維護
## DAY-298：強化粒子特效、加入高倍率爆炸演出、改善視覺衝擊力
extends Node

var _scene_root: Node = null

func _ready() -> void:
	call_deferred("_find_scene_root")

func _find_scene_root() -> void:
	# 自動尋找 Main 場景根節點
	var tree = get_tree()
	if tree == null:
		return
	var root = tree.get_root()
	if root == null:
		return
	# 找第一個 Node2D 子節點作為場景根
	for child in root.get_children():
		if child is Node2D:
			_scene_root = child
			return
	_scene_root = root

func set_scene_root(root: Node) -> void:
	_scene_root = root

## 命中特效（閃光環 + 小粒子）
func spawn_hit(pos: Vector2, char_id: String) -> void:
	if not is_instance_valid(_scene_root):
		return
	var color = _char_color(char_id)
	_spawn_flash_ring(pos, color, 0.12)
	# 小粒子（3個）
	_spawn_mini_particles(pos, color, 3)

## 擊破特效（爆炸粒子 + 閃光）
func spawn_kill(pos: Vector2, multiplier: float) -> void:
	if not is_instance_valid(_scene_root):
		return
	var color = _mult_color(multiplier)

	if multiplier >= 50:
		# 高倍率：大爆炸演出
		_spawn_flash_ring(pos, color, 0.35)
		_spawn_flash_ring(pos, Color.WHITE, 0.2)
		_spawn_particles(pos, color, 12)
		_spawn_ring_burst(pos, color)
	elif multiplier >= 20:
		# 中高倍率：中等爆炸
		_spawn_flash_ring(pos, color, 0.28)
		_spawn_particles(pos, color, 8)
		_spawn_ring_burst(pos, color)
	elif multiplier >= 10:
		# 中倍率：標準爆炸
		_spawn_flash_ring(pos, color, 0.22)
		_spawn_particles(pos, color, 6)
	else:
		# 低倍率：小爆炸
		_spawn_flash_ring(pos, color, 0.15)
		_spawn_particles(pos, color, 4)

## 大獎演出（≥20x）
func spawn_big_win(pos: Vector2, multiplier: float) -> void:
	if not is_instance_valid(_scene_root):
		return
	_spawn_big_win_text(pos, multiplier)
	if multiplier >= 100:
		ScreenShake.add_trauma(0.9)
		_spawn_star_burst(pos, _mult_color(multiplier))
	elif multiplier >= 50:
		ScreenShake.add_trauma(0.7)
		_spawn_star_burst(pos, _mult_color(multiplier))
	else:
		ScreenShake.add_trauma(0.5)

## 獎勵跳字
func spawn_reward_text(pos: Vector2, amount: int, multiplier: float) -> void:
	if not is_instance_valid(_scene_root):
		return
	var label = Label.new()
	var icon = "💰"
	var font_size = 18
	if multiplier >= 100:
		icon = "🌟"
		label.modulate = Color(1.0, 0.3, 0.1)
		font_size = 26
	elif multiplier >= 50:
		icon = "💫"
		label.modulate = Color(1.0, 0.5, 0.0)
		font_size = 22
	elif multiplier >= 20:
		icon = "⭐"
		label.modulate = Color(1.0, 0.85, 0.0)
		font_size = 20
	elif multiplier >= 10:
		icon = "✨"
		label.modulate = Color(1.0, 1.0, 0.4)
	label.text = "%s +%d" % [icon, amount]
	label.position = pos + Vector2(-30, -20)
	label.add_theme_font_size_override("font_size", font_size)
	label.z_index = 50
	_scene_root.add_child(label)

	# 高倍率跳字有彈跳效果
	var tween = label.create_tween()
	if multiplier >= 20:
		tween.tween_property(label, "scale", Vector2(1.3, 1.3), 0.08)
		tween.tween_property(label, "scale", Vector2(1.0, 1.0), 0.06)
	tween.tween_property(label, "position:y", label.position.y - 70, 0.8)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.8)
	tween.tween_callback(func(): if is_instance_valid(label): label.queue_free())

## Hit Stop（打擊感）
func hit_stop(duration: float = 0.04) -> void:
	Engine.time_scale = 0.0
	get_tree().create_timer(duration, true, false, true).timeout.connect(
		func(): Engine.time_scale = 1.0
	)

## 連鎖特效（連鎖閃電等 Lucky 系統用）
func spawn_chain_effect(from_pos: Vector2, to_pos: Vector2, color: Color) -> void:
	if not is_instance_valid(_scene_root):
		return
	# 閃電線段（用多個小方塊模擬）
	var steps = 8
	var prev = from_pos
	for i in range(1, steps + 1):
		var t = float(i) / float(steps)
		var mid = from_pos.lerp(to_pos, t)
		# 隨機偏移模擬閃電
		if i < steps:
			mid += Vector2(randf_range(-15, 15), randf_range(-15, 15))
		var seg = ColorRect.new()
		seg.size = Vector2(4, 4)
		seg.position = mid - seg.size / 2
		seg.color = color
		seg.z_index = 45
		_scene_root.add_child(seg)
		var tween = seg.create_tween()
		tween.tween_interval(0.02 * i)
		tween.tween_property(seg, "modulate:a", 0.0, 0.15)
		tween.tween_callback(func(): if is_instance_valid(seg): seg.queue_free())
		prev = mid

# ── 內部輔助 ──────────────────────────────────────────────────

func _spawn_flash_ring(pos: Vector2, color: Color, duration: float) -> void:
	var ring = ColorRect.new()
	ring.size = Vector2(48, 48)
	ring.position = pos - ring.size / 2
	ring.color = Color(color.r, color.g, color.b, 0.85)
	ring.z_index = 40
	_scene_root.add_child(ring)
	var tween = ring.create_tween()
	tween.tween_property(ring, "scale", Vector2(2.5, 2.5), duration)
	tween.parallel().tween_property(ring, "modulate:a", 0.0, duration)
	tween.tween_callback(func(): if is_instance_valid(ring): ring.queue_free())

func _spawn_mini_particles(pos: Vector2, color: Color, count: int) -> void:
	for i in count:
		var dot = ColorRect.new()
		dot.size = Vector2(4, 4)
		dot.color = color
		dot.position = pos
		dot.z_index = 40
		_scene_root.add_child(dot)
		var angle = (float(i) / float(count)) * TAU + randf_range(-0.3, 0.3)
		var dist = randf_range(15, 25)
		var target = pos + Vector2(cos(angle), sin(angle)) * dist
		var tween = dot.create_tween()
		tween.tween_property(dot, "position", target, 0.2)
		tween.parallel().tween_property(dot, "modulate:a", 0.0, 0.2)
		tween.tween_callback(func(): if is_instance_valid(dot): dot.queue_free())

func _spawn_particles(pos: Vector2, color: Color, count: int = 6) -> void:
	for i in count:
		var dot = ColorRect.new()
		var size = randf_range(4, 8)
		dot.size = Vector2(size, size)
		dot.color = color
		dot.position = pos
		dot.z_index = 40
		_scene_root.add_child(dot)
		var angle = (float(i) / float(count)) * TAU + randf_range(-0.2, 0.2)
		var dist = randf_range(30, 60)
		var target = pos + Vector2(cos(angle), sin(angle)) * dist
		var tween = dot.create_tween()
		tween.tween_property(dot, "position", target, randf_range(0.25, 0.4))
		tween.parallel().tween_property(dot, "modulate:a", 0.0, 0.35)
		tween.tween_callback(func(): if is_instance_valid(dot): dot.queue_free())

func _spawn_ring_burst(pos: Vector2, color: Color) -> void:
	# 擴散環（高倍率擊破的視覺衝擊）
	for ring_i in 2:
		var ring = ColorRect.new()
		var start_size = 20.0 + ring_i * 10.0
		ring.size = Vector2(start_size, start_size)
		ring.position = pos - ring.size / 2
		ring.color = Color(color.r, color.g, color.b, 0.6 - ring_i * 0.2)
		ring.z_index = 38
		_scene_root.add_child(ring)
		var end_scale = 3.5 - ring_i * 0.5
		var tween = ring.create_tween()
		tween.tween_interval(ring_i * 0.06)
		tween.tween_property(ring, "scale", Vector2(end_scale, end_scale), 0.3)
		tween.parallel().tween_property(ring, "modulate:a", 0.0, 0.3)
		tween.tween_callback(func(): if is_instance_valid(ring): ring.queue_free())

func _spawn_star_burst(pos: Vector2, color: Color) -> void:
	# 星形爆炸（超高倍率用）
	var star_count = 8
	for i in star_count:
		var star = Label.new()
		star.text = "★"
		star.position = pos + Vector2(-8, -8)
		star.add_theme_font_size_override("font_size", 20)
		star.modulate = color
		star.z_index = 55
		_scene_root.add_child(star)
		var angle = (float(i) / float(star_count)) * TAU
		var dist = randf_range(50, 90)
		var target = pos + Vector2(cos(angle), sin(angle)) * dist
		var tween = star.create_tween()
		tween.tween_property(star, "position", target, 0.5)
		tween.parallel().tween_property(star, "scale", Vector2(1.5, 1.5), 0.2)
		tween.parallel().tween_property(star, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func(): if is_instance_valid(star): star.queue_free())

func _spawn_big_win_text(pos: Vector2, multiplier: float) -> void:
	var label = Label.new()
	var font_size = 36
	if multiplier >= 100:
		label.text = "🌟 x%.0f MEGA WIN!" % multiplier
		label.modulate = Color(1.0, 0.3, 0.1)
		font_size = 42
	elif multiplier >= 50:
		label.text = "💫 x%.0f BIG WIN!" % multiplier
		label.modulate = Color(1.0, 0.6, 0.0)
		font_size = 38
	else:
		label.text = "✨ x%.0f BIG WIN!" % multiplier
		label.modulate = Color(1.0, 0.85, 0.0)
	label.position = Vector2(640 - 120, 290)
	label.add_theme_font_size_override("font_size", font_size)
	label.z_index = 60
	_scene_root.add_child(label)
	var tween = label.create_tween()
	tween.tween_property(label, "scale", Vector2(1.4, 1.4), 0.12)
	tween.tween_property(label, "scale", Vector2(1.0, 1.0), 0.08)
	tween.tween_interval(0.9)
	tween.tween_property(label, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func(): if is_instance_valid(label): label.queue_free())

func _char_color(char_id: String) -> Color:
	match char_id:
		"hachiware": return Color(0.4, 0.6, 1.0)
		"usagi": return Color(1.0, 0.9, 0.2)
		_: return Color(1.0, 0.6, 0.8)

func _mult_color(mult: float) -> Color:
	if mult >= 100: return Color(1.0, 0.3, 0.1)
	if mult >= 50: return Color(1.0, 0.5, 0.0)
	if mult >= 20: return Color(1.0, 0.85, 0.0)
	if mult >= 10: return Color(1.0, 1.0, 0.4)
	return Color(0.8, 0.8, 0.8)
