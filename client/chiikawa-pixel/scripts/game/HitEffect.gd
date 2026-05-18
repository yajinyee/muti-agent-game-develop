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

# 特效節點腳本（用 _draw 繪製真正的圓形，比 ColorRect 更精確）
const FlashRingScript = preload("res://scripts/effects/FlashRing.gd")
const ShockwaveRingScript = preload("res://scripts/effects/ShockwaveRing.gd")

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

## 大獎全畫面特效（閃白 + 金色粒子雨 + 螢幕扭曲）
func spawn_big_win(pos: Vector2, multiplier: float = 100.0) -> void:
	_ensure_root()
	_spawn_screen_flash(Color(1.0, 0.9, 0.2, 0.6), 0.08, 0.3)
	_spawn_flash_ring(pos, Color.GOLD, 80.0, 0.35)
	_spawn_shockwave(pos, Color.GOLD, 3.0)
	_spawn_particles(pos, Color.GOLD, 20, 120.0, 0.8)
	_spawn_particles(pos, Color.WHITE, 10, 90.0, 0.6)
	# 大獎螢幕扭曲（帶金色色差）
	spawn_screen_shockwave(pos, 0.4)

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

## BOSS 登場特效（全畫面紅色閃爍 + 粒子爆炸 + 震動文字）
## boss_pos: BOSS 實際位置（預設右側進場 x=1100）
func spawn_boss_enter(boss_pos: Vector2 = Vector2(1100, 360)) -> void:
	_ensure_root()
	# 第一波：強烈紅色閃光
	_spawn_screen_flash(Color(0.9, 0.0, 0.0, 0.7), 0.08, 0.25)
	# 第二波：稍弱的橙紅閃光（0.15 秒後）
	get_tree().create_timer(0.15).timeout.connect(func():
		_spawn_screen_flash(Color(0.8, 0.1, 0.0, 0.45), 0.05, 0.3)
	)
	# BOSS 位置爆炸粒子（從 BOSS 實際位置噴射，更有衝擊感）
	_spawn_particles(boss_pos, Color(1.0, 0.2, 0.1), 24, 220.0, 0.65)
	_spawn_particles(boss_pos, Color(1.0, 0.6, 0.0), 14, 160.0, 0.55)
	_spawn_particles(boss_pos, Color.WHITE, 10, 110.0, 0.45)
	# 畫面中心補充粒子（讓整個畫面都有感）
	var center = Vector2(640, 360)
	_spawn_particles(center, Color(1.0, 0.3, 0.1), 10, 150.0, 0.5)
	# BOSS 位置多個衝擊波（由小到大）
	_spawn_shockwave(boss_pos, Color(1.0, 0.2, 0.1), 2.5)
	get_tree().create_timer(0.08).timeout.connect(func():
		_spawn_shockwave(boss_pos, Color(1.0, 0.5, 0.0), 3.5)
	)
	get_tree().create_timer(0.18).timeout.connect(func():
		_spawn_shockwave(boss_pos, Color(0.8, 0.1, 0.0), 2.0)
	)
	# 地面衝擊波（從 BOSS 底部向兩側擴散）
	_spawn_ground_shockwave(Vector2(boss_pos.x, 680.0))
	# 螢幕扭曲衝擊波（真正的螢幕扭曲，帶色差，最震撼的效果）
	spawn_screen_shockwave(boss_pos, 0.55)
	# BOSS 登場文字（畫面中央）
	_spawn_boss_enter_text()

## Bonus 觸發特效（全畫面金色閃爍）
func spawn_bonus_trigger() -> void:
	_ensure_root()
	_spawn_screen_flash(Color(1.0, 0.85, 0.0, 0.5), 0.08, 0.35)

## 角色升級特效（勞動值滿 100 觸發 Bonus 前的慶祝動畫）
## pos: 砲台位置（通常是 Vector2(640, 630)）
## char_id: 角色 ID（決定顏色）
func spawn_level_up(pos: Vector2 = Vector2(640, 630), char_id: String = "chiikawa") -> void:
	_ensure_root()
	var char_color = CHAR_COLORS.get(char_id, Color.WHITE)

	# 1. 全畫面角色色閃光（短暫）
	_spawn_screen_flash(Color(char_color.r, char_color.g, char_color.b, 0.35), 0.06, 0.25)

	# 2. 從砲台位置噴射金色星星粒子（向上噴射）
	_spawn_level_up_stars(pos, char_color)

	# 3. 大閃光環（角色顏色）
	_spawn_flash_ring(pos, char_color, 60.0, 0.4)
	_spawn_flash_ring(pos, Color.GOLD, 40.0, 0.3)

	# 4. 衝擊波
	_spawn_shockwave(pos, char_color, 2.0)

	# 5. 升級文字（畫面中央偏上）
	_spawn_level_up_text(char_id)

## 升級星星粒子（從砲台向上噴射，帶重力下落）
func _spawn_level_up_stars(pos: Vector2, color: Color) -> void:
	if not is_instance_valid(_scene_root):
		return

	var star_count = 16
	for i in star_count:
		var star = ColorRect.new()
		var sz = randf_range(4.0, 9.0)
		star.size = Vector2(sz, sz)
		star.position = pos + Vector2(randf_range(-20, 20), randf_range(-10, 10))
		star.color = Color.GOLD if randf() > 0.4 else color
		star.z_index = 12
		_scene_root.add_child(star)

		# 向上噴射（帶隨機角度，偏向上方）
		var angle = randf_range(-PI * 0.85, -PI * 0.15)  # 上方扇形
		var speed = randf_range(180.0, 380.0)
		var target = pos + Vector2(cos(angle) * speed, sin(angle) * speed)
		# 加重力（向下偏移）
		target.y += randf_range(80.0, 160.0)

		var duration = randf_range(0.5, 0.9)
		var tween = star.create_tween()
		tween.tween_property(star, "position", target, duration).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_QUAD)
		tween.parallel().tween_property(star, "modulate:a", 0.0, duration * 0.7)
		tween.tween_callback(func():
			if is_instance_valid(star):
				star.queue_free()
		)

## 升級文字動畫（畫面中央偏上）
func _spawn_level_up_text(char_id: String) -> void:
	if not is_instance_valid(_scene_root):
		return

	var canvas = CanvasLayer.new()
	canvas.layer = 97
	_scene_root.add_child(canvas)

	# 角色對應的升級文字
	var texts = {
		"chiikawa": "BONUS READY!",
		"hachiware": "BONUS READY!",
		"usagi": "BONUS READY!",
	}
	var char_color = CHAR_COLORS.get(char_id, Color.WHITE)

	# 陰影
	var shadow = Label.new()
	shadow.text = texts.get(char_id, "BONUS READY!")
	shadow.position = Vector2(322, 202)
	shadow.add_theme_font_size_override("font_size", 52)
	shadow.modulate = Color(0.0, 0.0, 0.0, 0.7)
	shadow.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	shadow.size = Vector2(640, 70)
	canvas.add_child(shadow)

	# 主文字
	var label = Label.new()
	label.text = texts.get(char_id, "BONUS READY!")
	label.position = Vector2(320, 200)
	label.add_theme_font_size_override("font_size", 52)
	label.modulate = Color(char_color.r * 1.2, char_color.g * 1.2, char_color.b * 1.2, 1.0)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.size = Vector2(640, 70)
	canvas.add_child(label)

	# 副標題（金色星星）
	var sub = Label.new()
	sub.text = "★ ★ ★"
	sub.position = Vector2(320, 260)
	sub.add_theme_font_size_override("font_size", 28)
	sub.modulate = Color(1.0, 0.9, 0.2, 0.0)
	sub.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	sub.size = Vector2(640, 40)
	canvas.add_child(sub)

	# 動畫：彈入 → 停留 → 淡出
	label.scale = Vector2(0.2, 0.2)
	shadow.scale = Vector2(0.2, 0.2)

	var tween = canvas.create_tween()
	# 彈入（0.25 秒，BACK 彈性）
	tween.tween_property(label, "scale", Vector2(1.05, 1.05), 0.25).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.parallel().tween_property(shadow, "scale", Vector2(1.05, 1.05), 0.25).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.parallel().tween_property(sub, "modulate:a", 1.0, 0.25)
	# 回彈
	tween.tween_property(label, "scale", Vector2(1.0, 1.0), 0.08)
	tween.parallel().tween_property(shadow, "scale", Vector2(1.0, 1.0), 0.08)
	# 停留（0.7 秒）
	tween.tween_interval(0.7)
	# 淡出（0.3 秒）
	tween.tween_property(label, "modulate:a", 0.0, 0.3)
	tween.parallel().tween_property(shadow, "modulate:a", 0.0, 0.3)
	tween.parallel().tween_property(sub, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(canvas):
			canvas.queue_free()
	)

## 像素化過場效果（場景切換時使用）
## 先像素化（0→1，duration_in 秒），再還原（1→0，duration_out 秒）
## callback 在最像素化時執行（適合在此切換場景內容）
func pixelate_transition(duration_in: float = 0.2, duration_out: float = 0.3, callback: Callable = Callable()) -> void:
	_ensure_root()
	
	var shader_path = "res://assets/shaders/pixelate_transition.gdshader"
	if not ResourceLoader.exists(shader_path):
		# 備用：直接執行 callback
		if callback.is_valid():
			callback.call()
		return
	
	# 建立全畫面覆蓋層
	var canvas = CanvasLayer.new()
	canvas.layer = 99
	_scene_root.add_child(canvas)
	
	var rect = ColorRect.new()
	rect.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	rect.color = Color.WHITE  # 顏色不重要，shader 會覆蓋
	canvas.add_child(rect)
	
	# 套用像素化 shader
	var mat = ShaderMaterial.new()
	mat.shader = load(shader_path)
	mat.set_shader_parameter("pixelate_amount", 0.0)
	mat.set_shader_parameter("screen_size", Vector2(1280.0, 720.0))
	rect.material = mat
	
	# 像素化進入
	var tween = rect.create_tween()
	tween.tween_method(func(v): mat.set_shader_parameter("pixelate_amount", v), 0.0, 1.0, duration_in)
	tween.tween_callback(func():
		# 最像素化時執行 callback（切換場景內容）
		if callback.is_valid():
			callback.call()
	)
	# 像素化還原
	tween.tween_method(func(v): mat.set_shader_parameter("pixelate_amount", v), 1.0, 0.0, duration_out)
	tween.tween_callback(func():
		if is_instance_valid(canvas):
			canvas.queue_free()
	)

# ── 內部實作 ──────────────────────────────────────────

func _ensure_root() -> void:
	if not is_instance_valid(_scene_root):
		_scene_root = get_tree().current_scene

## 閃光環（真正的圓形擴散，用 _draw 繪製）
func _spawn_flash_ring(pos: Vector2, color: Color, radius: float, duration: float) -> void:
	if not is_instance_valid(_scene_root):
		return

	# 用 FlashRingScript 繪製真正的圓形（比 ColorRect 更精確）
	var ring = FlashRingScript.new()
	ring.position = pos
	ring.z_index = 10
	ring.ring_color = color
	ring.ring_radius = radius
	_scene_root.add_child(ring)

	var tween = ring.create_tween()
	tween.tween_property(ring, "scale", Vector2(2.2, 2.2), duration * 0.6).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(ring, "modulate:a", 0.0, duration)
	tween.tween_callback(func():
		if is_instance_valid(ring):
			ring.queue_free()
	)

## 衝擊波（向外擴散的環）
func _spawn_shockwave(pos: Vector2, color: Color, scale_factor: float) -> void:
	if not is_instance_valid(_scene_root):
		return

	# 用 ShockwaveRingScript 繪製真正的圓形衝擊波環
	var wave = ShockwaveRingScript.new()
	wave.position = pos
	wave.z_index = 9
	wave.ring_color = color
	wave.ring_radius = 12.0 * scale_factor
	_scene_root.add_child(wave)

	var tween = wave.create_tween()
	tween.tween_property(wave, "scale", Vector2(4.5 * scale_factor, 4.5 * scale_factor), 0.28).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(wave, "modulate:a", 0.0, 0.28)
	tween.tween_callback(func():
		if is_instance_valid(wave):
			wave.queue_free()
	)

## 地面衝擊波（BOSS 登場時從底部向兩側擴散的橫向波）
func _spawn_ground_shockwave(pos: Vector2) -> void:
	if not is_instance_valid(_scene_root):
		return

	# 左右各一條橫向衝擊線
	for direction in [-1, 1]:
		var line = ColorRect.new()
		line.size = Vector2(4, 40)
		line.position = pos + Vector2(-2, -20)
		line.color = Color(1.0, 0.3, 0.1, 0.8)
		line.z_index = 9
		_scene_root.add_child(line)

		var target_x = pos.x + direction * 700.0
		var tween = line.create_tween()
		tween.tween_property(line, "position:x", target_x, 0.35).set_ease(Tween.EASE_OUT)
		tween.parallel().tween_property(line, "modulate:a", 0.0, 0.35)
		tween.parallel().tween_property(line, "size:y", 8.0, 0.35)
		tween.tween_callback(func():
			if is_instance_valid(line):
				line.queue_free()
		)

## 螢幕扭曲衝擊波（BOSS 登場 / 大獎用，真正的螢幕扭曲效果）
## 比 ColorRect 衝擊波視覺效果強很多，帶色差（chromatic aberration）
func spawn_screen_shockwave(world_pos: Vector2, duration: float = 0.5) -> void:
	_ensure_root()
	var shader_path = "res://assets/shaders/shockwave_distortion.gdshader"
	if not ResourceLoader.exists(shader_path):
		return

	# 建立全畫面 CanvasLayer + ColorRect
	var canvas = CanvasLayer.new()
	canvas.layer = 95
	_scene_root.add_child(canvas)

	var rect = ColorRect.new()
	rect.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	canvas.add_child(rect)

	var mat = ShaderMaterial.new()
	mat.shader = load(shader_path)

	# 將世界座標轉換為螢幕 UV（0-1 範圍）
	var viewport_size = Vector2(1280.0, 720.0)
	var normalized_pos = world_pos / viewport_size
	mat.set_shader_parameter("center", normalized_pos)
	mat.set_shader_parameter("radius", 0.0)
	mat.set_shader_parameter("strength", 0.06)
	rect.material = mat

	# 動畫：radius 從 0 擴散到 0.8（衝擊波向外擴散）
	var tween = rect.create_tween()
	tween.tween_method(func(v): mat.set_shader_parameter("radius", v), 0.0, 0.8, duration)
	tween.parallel().tween_method(func(v): mat.set_shader_parameter("strength", v), 0.06, 0.0, duration)
	tween.tween_callback(func():
		if is_instance_valid(canvas):
			canvas.queue_free()
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

## BOSS 登場文字動畫（畫面中央大字）
func _spawn_boss_enter_text() -> void:
	if not is_instance_valid(_scene_root):
		return

	var canvas = CanvasLayer.new()
	canvas.layer = 98
	_scene_root.add_child(canvas)

	# 外框陰影文字（黑色，稍微偏移）
	var shadow = Label.new()
	shadow.text = "⚔ BOSS ⚔"
	shadow.position = Vector2(642, 282)  # 偏移 2px
	shadow.add_theme_font_size_override("font_size", 64)
	shadow.modulate = Color(0.0, 0.0, 0.0, 0.8)
	shadow.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	shadow.size = Vector2(640, 80)
	shadow.position.x = 320
	canvas.add_child(shadow)

	# 主文字（紅色）
	var label = Label.new()
	label.text = "⚔ BOSS ⚔"
	label.position = Vector2(320, 280)
	label.add_theme_font_size_override("font_size", 64)
	label.modulate = Color(1.0, 0.15, 0.15, 1.0)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.size = Vector2(640, 80)
	canvas.add_child(label)

	# 副標題
	var sub = Label.new()
	sub.text = "INCOMING!"
	sub.position = Vector2(320, 355)
	sub.add_theme_font_size_override("font_size", 28)
	sub.modulate = Color(1.0, 0.6, 0.0, 1.0)
	sub.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	sub.size = Vector2(640, 40)
	canvas.add_child(sub)

	# 動畫：從小放大 → 停留 → 縮小消失
	label.scale = Vector2(0.3, 0.3)
	shadow.scale = Vector2(0.3, 0.3)
	sub.modulate.a = 0.0

	var tween = canvas.create_tween()
	# 放大衝擊（0.2 秒）
	tween.tween_property(label, "scale", Vector2(1.1, 1.1), 0.2).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.parallel().tween_property(shadow, "scale", Vector2(1.1, 1.1), 0.2).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.parallel().tween_property(sub, "modulate:a", 1.0, 0.2)
	# 回彈（0.1 秒）
	tween.tween_property(label, "scale", Vector2(1.0, 1.0), 0.1)
	tween.parallel().tween_property(shadow, "scale", Vector2(1.0, 1.0), 0.1)
	# 停留（0.8 秒）
	tween.tween_interval(0.8)
	# 淡出（0.3 秒）
	tween.tween_property(label, "modulate:a", 0.0, 0.3)
	tween.parallel().tween_property(shadow, "modulate:a", 0.0, 0.3)
	tween.parallel().tween_property(sub, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(canvas):
			canvas.queue_free()
	)
