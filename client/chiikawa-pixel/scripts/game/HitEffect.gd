## HitEffect.gd
## 命中特效系統（Autoload 單例）
## 提供：命中閃光、擊殺爆炸、大獎全畫面、Hit Stop（時間凍結）
##
## 使用方式：
##   HitEffect.spawn_hit(pos, char_id)          # 普通命中
##   HitEffect.spawn_kill(pos, multiplier)       # 擊殺爆炸
##   HitEffect.spawn_big_win(pos, multiplier)    # 大獎特效
##   HitEffect.hit_stop(0.05)                    # 短暫時間凍結

extends Node

# 角色顏色
const CHAR_COLORS = {
	"chiikawa": Color(1.0, 0.6, 0.8),
	"hachiware": Color(0.4, 0.6, 1.0),
	"usagi": Color(1.0, 0.9, 0.2),
}

# 倍率對應顏色
const MULTIPLIER_COLORS = {
	100: Color(1.0, 0.2, 0.2),   # 紅（超高倍）
	20:  Color(1.0, 0.8, 0.0),   # 金（高倍）
	5:   Color(0.4, 1.0, 0.4),   # 綠（中倍）
	0:   Color(1.0, 1.0, 1.0),   # 白（低倍）
}

var _scene_root: Node = null
var _hit_stop_active: bool = false

func _ready() -> void:
	# 等待場景樹就緒
	call_deferred("_find_scene_root")

func _find_scene_root() -> void:
	_scene_root = get_tree().current_scene

# ── 公開 API ──────────────────────────────────────────

## 普通命中特效（小閃光 + 少量粒子）
func spawn_hit(pos: Vector2, char_id: String = "chiikawa") -> void:
	_ensure_root()
	var color = CHAR_COLORS.get(char_id, Color.WHITE)
	_spawn_flash_ring(pos, color, 18.0, 0.12)
	_spawn_particles(pos, color, 4, 25.0, 0.25)

## 擊殺爆炸特效（大閃光 + 多粒子 + 衝擊波）
func spawn_kill(pos: Vector2, multiplier: float = 1.0) -> void:
	_ensure_root()
	var color = _get_multiplier_color(multiplier)
	var scale_factor = clamp(1.0 + multiplier * 0.01, 1.0, 2.5)

	_spawn_flash_ring(pos, color, 32.0 * scale_factor, 0.20)
	_spawn_shockwave(pos, color, scale_factor)
	_spawn_particles(pos, color, 8, 50.0 * scale_factor, 0.45)

	# 高倍率加強特效
	if multiplier >= 20:
		_spawn_flash_ring(pos, Color.WHITE, 48.0, 0.15)
		_spawn_particles(pos, Color.GOLD, 6, 70.0, 0.5)

## 大獎全畫面特效（閃白 + 金色粒子雨）
func spawn_big_win(pos: Vector2, multiplier: float = 100.0) -> void:
	_ensure_root()
	_spawn_screen_flash(Color(1.0, 0.9, 0.2, 0.6), 0.08, 0.3)
	_spawn_flash_ring(pos, Color.GOLD, 80.0, 0.35)
	_spawn_shockwave(pos, Color.GOLD, 3.0)
	_spawn_particles(pos, Color.GOLD, 20, 120.0, 0.8)
	_spawn_particles(pos, Color.WHITE, 10, 90.0, 0.6)

## Hit Stop — 短暫凍結時間（增加打擊感）
## duration: 凍結秒數（建議 0.03~0.08）
func hit_stop(duration: float = 0.05) -> void:
	if _hit_stop_active:
		return
	_hit_stop_active = true
	Engine.time_scale = 0.0
	await get_tree().create_timer(duration, true, false, true).timeout
	Engine.time_scale = 1.0
	_hit_stop_active = false

## BOSS 登場特效（全畫面紅色閃爍）
func spawn_boss_enter() -> void:
	_ensure_root()
	_spawn_screen_flash(Color(0.8, 0.0, 0.0, 0.5), 0.1, 0.4)
	_spawn_screen_flash(Color(0.8, 0.0, 0.0, 0.3), 0.05, 0.2)

## Bonus 觸發特效（全畫面金色閃爍）
func spawn_bonus_trigger() -> void:
	_ensure_root()
	_spawn_screen_flash(Color(1.0, 0.85, 0.0, 0.5), 0.08, 0.35)

# ── 內部實作 ──────────────────────────────────────────

func _ensure_root() -> void:
	if not is_instance_valid(_scene_root):
		_scene_root = get_tree().current_scene

## 閃光環（圓形擴散）
func _spawn_flash_ring(pos: Vector2, color: Color, radius: float, duration: float) -> void:
	if not is_instance_valid(_scene_root):
		return

	# 用多個 ColorRect 模擬圓形（像素風格）
	var ring = Node2D.new()
	ring.position = pos
	ring.z_index = 10

	# 中心亮點
	var center = ColorRect.new()
	center.size = Vector2(radius * 0.5, radius * 0.5)
	center.position = -center.size / 2
	center.color = Color(color.r, color.g, color.b, 0.9)
	ring.add_child(center)

	# 外環（4 個方向的矩形）
	for i in 4:
		var bar = ColorRect.new()
		bar.size = Vector2(radius * 0.3, radius * 0.15)
		var angle = i * PI / 2.0
		bar.position = Vector2(cos(angle) * radius * 0.4, sin(angle) * radius * 0.4) - bar.size / 2
		bar.color = Color(color.r, color.g, color.b, 0.7)
		ring.add_child(bar)

	_scene_root.add_child(ring)

	var tween = ring.create_tween()
	tween.tween_property(ring, "scale", Vector2(2.2, 2.2), duration * 0.6)
	tween.parallel().tween_property(ring, "modulate:a", 0.0, duration)
	tween.tween_callback(func():
		if is_instance_valid(ring):
			ring.queue_free()
	)

## 衝擊波（向外擴散的環）
func _spawn_shockwave(pos: Vector2, color: Color, scale_factor: float) -> void:
	if not is_instance_valid(_scene_root):
		return

	var wave = ColorRect.new()
	var size = 20.0 * scale_factor
	wave.size = Vector2(size, size)
	wave.position = pos - Vector2(size, size) / 2
	wave.color = Color(color.r, color.g, color.b, 0.6)
	wave.z_index = 9
	_scene_root.add_child(wave)

	var tween = wave.create_tween()
	tween.tween_property(wave, "scale", Vector2(4.0 * scale_factor, 4.0 * scale_factor), 0.25)
	tween.parallel().tween_property(wave, "modulate:a", 0.0, 0.25)
	tween.tween_callback(func():
		if is_instance_valid(wave):
			wave.queue_free()
	)

## 粒子噴射
func _spawn_particles(pos: Vector2, color: Color, count: int, spread: float, duration: float) -> void:
	if not is_instance_valid(_scene_root):
		return

	for i in count:
		var p = ColorRect.new()
		var psize = randf_range(3.0, 7.0)
		p.size = Vector2(psize, psize)
		p.position = pos + Vector2(randf_range(-8, 8), randf_range(-8, 8))
		p.color = color
		p.z_index = 11
		_scene_root.add_child(p)

		# 隨機方向噴射
		var angle = randf() * TAU
		var dist = randf_range(spread * 0.4, spread)
		var target = pos + Vector2(cos(angle) * dist, sin(angle) * dist - spread * 0.3)

		var tween = p.create_tween()
		tween.tween_property(p, "position", target, duration)
		tween.parallel().tween_property(p, "modulate:a", 0.0, duration * 0.8)
		tween.tween_callback(func():
			if is_instance_valid(p):
				p.queue_free()
		)

## 全畫面閃光（CanvasLayer 上的 ColorRect）
func _spawn_screen_flash(color: Color, hold: float, fade: float) -> void:
	if not is_instance_valid(_scene_root):
		return

	var canvas = CanvasLayer.new()
	canvas.layer = 100  # 最上層
	_scene_root.add_child(canvas)

	var rect = ColorRect.new()
	rect.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	rect.color = color
	canvas.add_child(rect)

	var tween = rect.create_tween()
	tween.tween_interval(hold)
	tween.tween_property(rect, "modulate:a", 0.0, fade)
	tween.tween_callback(func():
		if is_instance_valid(canvas):
			canvas.queue_free()
	)

## 依倍率取得顏色
func _get_multiplier_color(multiplier: float) -> Color:
	if multiplier >= 100:
		return MULTIPLIER_COLORS[100]
	elif multiplier >= 20:
		return MULTIPLIER_COLORS[20]
	elif multiplier >= 5:
		return MULTIPLIER_COLORS[5]
	else:
		return MULTIPLIER_COLORS[0]
