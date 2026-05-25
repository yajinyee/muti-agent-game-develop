## HitEffect.gd — 命中特效、擊破粒子、獎勵跳字
## hit-effect-agent 負責維護
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

## 命中特效（閃光環）
func spawn_hit(pos: Vector2, char_id: String) -> void:
	if not is_instance_valid(_scene_root):
		return
	var color = _char_color(char_id)
	_spawn_flash_ring(pos, color, 0.15)

## 擊破特效（爆炸粒子）
func spawn_kill(pos: Vector2, multiplier: float) -> void:
	if not is_instance_valid(_scene_root):
		return
	var color = _mult_color(multiplier)
	_spawn_flash_ring(pos, color, 0.25)
	_spawn_particles(pos, color)

## 大獎演出（≥20x）
func spawn_big_win(pos: Vector2, multiplier: float) -> void:
	if not is_instance_valid(_scene_root):
		return
	_spawn_big_win_text(pos, multiplier)
	ScreenShake.add_trauma(0.7)

## 獎勵跳字
func spawn_reward_text(pos: Vector2, amount: int, multiplier: float) -> void:
	if not is_instance_valid(_scene_root):
		return
	var label = Label.new()
	var icon = "💰"
	if multiplier >= 100:
		icon = "🌟"
		label.modulate = Color(1.0, 0.3, 0.1)
	elif multiplier >= 20:
		icon = "⭐"
		label.modulate = Color(1.0, 0.85, 0.0)
	elif multiplier >= 10:
		icon = "✨"
		label.modulate = Color(1.0, 1.0, 0.4)
	label.text = "%s +%d" % [icon, amount]
	label.position = pos + Vector2(-30, -20)
	label.add_theme_font_size_override("font_size", 18)
	label.z_index = 50
	_scene_root.add_child(label)
	var tween = label.create_tween()
	tween.tween_property(label, "position:y", label.position.y - 60, 0.7)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.7)
	tween.tween_callback(func(): if is_instance_valid(label): label.queue_free())

## Hit Stop（打擊感）
func hit_stop(duration: float = 0.04) -> void:
	Engine.time_scale = 0.0
	get_tree().create_timer(duration, true, false, true).timeout.connect(
		func(): Engine.time_scale = 1.0
	)

# ── 內部輔助 ──────────────────────────────────────────────────

func _spawn_flash_ring(pos: Vector2, color: Color, duration: float) -> void:
	var ring = ColorRect.new()
	ring.size = Vector2(48, 48)
	ring.position = pos - ring.size / 2
	ring.color = Color(color.r, color.g, color.b, 0.8)
	ring.z_index = 40
	_scene_root.add_child(ring)
	var tween = ring.create_tween()
	tween.tween_property(ring, "scale", Vector2(2.0, 2.0), duration)
	tween.parallel().tween_property(ring, "modulate:a", 0.0, duration)
	tween.tween_callback(func(): if is_instance_valid(ring): ring.queue_free())

func _spawn_particles(pos: Vector2, color: Color) -> void:
	for i in 6:
		var dot = ColorRect.new()
		dot.size = Vector2(6, 6)
		dot.color = color
		dot.position = pos
		dot.z_index = 40
		_scene_root.add_child(dot)
		var angle = i * PI / 3.0
		var target = pos + Vector2(cos(angle), sin(angle)) * 40
		var tween = dot.create_tween()
		tween.tween_property(dot, "position", target, 0.3)
		tween.parallel().tween_property(dot, "modulate:a", 0.0, 0.3)
		tween.tween_callback(func(): if is_instance_valid(dot): dot.queue_free())

func _spawn_big_win_text(pos: Vector2, multiplier: float) -> void:
	var label = Label.new()
	label.text = "✨ x%.0f BIG WIN!" % multiplier
	label.position = Vector2(640 - 100, 300)
	label.add_theme_font_size_override("font_size", 36)
	label.modulate = Color(1.0, 0.85, 0.0)
	label.z_index = 60
	_scene_root.add_child(label)
	var tween = label.create_tween()
	tween.tween_property(label, "scale", Vector2(1.3, 1.3), 0.15)
	tween.tween_property(label, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(1.0)
	tween.tween_property(label, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func(): if is_instance_valid(label): label.queue_free())

func _char_color(char_id: String) -> Color:
	match char_id:
		"hachiware": return Color(0.4, 0.6, 1.0)
		"usagi": return Color(1.0, 0.9, 0.2)
		_: return Color(1.0, 0.6, 0.8)

func _mult_color(mult: float) -> Color:
	if mult >= 50: return Color(1.0, 0.3, 0.1)
	if mult >= 20: return Color(1.0, 0.85, 0.0)
	if mult >= 10: return Color(1.0, 1.0, 0.4)
	return Color(0.8, 0.8, 0.8)
