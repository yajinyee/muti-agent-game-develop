## BubbleLayer.gd
## 海底動態環境層：氣泡 + 光線柱 + 海草搖擺 + 漂浮微粒 + 遠景小魚群
## 在 normal 背景時顯示，提升海底沉浸感

extends Node2D

# ---- 氣泡 ----
var _bubbles: Array = []
var _spawn_timer: float = 0.0
var _spawn_interval: float = 0.35

# ---- 光線柱（從水面射入的光束）----
var _light_beams: Array = []  # {x, width, alpha, speed, phase}

# ---- 海草 ----
var _seaweeds: Array = []  # {x, y, height, segments, color, phase, speed}

# ---- 漂浮微粒（浮游生物/塵埃）----
var _particles: Array = []  # {x, y, radius, alpha, drift_x, drift_y, phase, color}
var _particle_spawn_timer: float = 0.0
var _particle_spawn_interval: float = 0.8

# ---- 遠景小魚群 ----
var _fish_schools: Array = []  # {x, y, dir, speed, fish_count, fish_offsets, color, alpha}
var _fish_spawn_timer: float = 0.0
var _fish_spawn_interval: float = 6.0

# ---- 水面波紋節點 ----
var _water_surface: ColorRect = null
var _water_mat: ShaderMaterial = null

var _active: bool = true
var _time: float = 0.0

const WATER_SURFACE_SHADER = "res://assets/shaders/water_surface.gdshader"

func _ready() -> void:
	# 預生成氣泡
	for i in range(10):
		_spawn_bubble(randf_range(50, 1230), randf_range(80, 680))

	# 初始化光線柱（3-5 條，固定位置，緩慢閃爍）
	for i in range(4):
		_light_beams.append({
			"x": 150.0 + i * 280.0 + randf_range(-60, 60),
			"width": randf_range(40.0, 90.0),
			"alpha": randf_range(0.04, 0.10),
			"phase": randf_range(0.0, TAU),
			"speed": randf_range(0.3, 0.7),
		})

	# 初始化海草（底部，8-12 株）
	for i in range(10):
		var x = 60.0 + i * 120.0 + randf_range(-30, 30)
		var h = randf_range(60.0, 130.0)
		var segs = int(h / 18.0) + 2
		_seaweeds.append({
			"x": x,
			"y": 720.0,  # 底部
			"height": h,
			"segments": segs,
			"color": Color(
				randf_range(0.05, 0.15),
				randf_range(0.35, 0.65),
				randf_range(0.15, 0.35),
				randf_range(0.55, 0.80)
			),
			"phase": randf_range(0.0, TAU),
			"speed": randf_range(0.6, 1.4),
			"amp": randf_range(6.0, 14.0),
		})

	# 預生成漂浮微粒（20 個）
	for i in range(20):
		_spawn_particle(randf_range(0, 1280), randf_range(100, 680))

	# 預生成遠景小魚群（2 群）
	for i in range(2):
		_spawn_fish_school(randf_range(200, 1000), randf_range(150, 550))

	# 建立水面波紋效果
	_setup_water_surface()

func _setup_water_surface() -> void:
	if not ResourceLoader.exists(WATER_SURFACE_SHADER):
		return
	_water_surface = ColorRect.new()
	_water_surface.size = Vector2(1280, 80)
	_water_surface.position = Vector2(0, 0)
	_water_surface.z_index = 3  # 在所有效果上方
	_water_mat = ShaderMaterial.new()
	_water_mat.shader = load(WATER_SURFACE_SHADER)
	_water_surface.material = _water_mat
	call_deferred("add_child", _water_surface)

func set_active(v: bool) -> void:
	_active = v
	if not v:
		_bubbles.clear()
		_particles.clear()
		_fish_schools.clear()
		queue_redraw()
	if is_instance_valid(_water_surface):
		_water_surface.visible = v

func _process(delta: float) -> void:
	if not _active:
		return

	_time += delta

	# 生成新氣泡
	_spawn_timer += delta
	if _spawn_timer >= _spawn_interval:
		_spawn_timer = 0.0
		_spawn_bubble(randf_range(50, 1230), 730.0)

	# 生成新漂浮微粒
	_particle_spawn_timer += delta
	if _particle_spawn_timer >= _particle_spawn_interval:
		_particle_spawn_timer = 0.0
		_spawn_particle(randf_range(0, 1280), randf_range(200, 680))

	# 生成新魚群
	_fish_spawn_timer += delta
	if _fish_spawn_timer >= _fish_spawn_interval:
		_fish_spawn_timer = 0.0
		# 從左或右邊緣生成
		var from_left = randf() > 0.5
		var x = -80.0 if from_left else 1360.0
		_spawn_fish_school(x, randf_range(150, 500))

	# 更新氣泡
	var i = _bubbles.size() - 1
	while i >= 0:
		var b = _bubbles[i]
		b["y"] -= b["speed"] * delta
		b["x"] += sin(_time * b["wobble_speed"] + b["wobble_phase"]) * b["wobble_amp"] * delta
		if b["y"] < 80:
			b["alpha"] -= delta * 1.8
		if b["alpha"] <= 0.0 or b["y"] < -10:
			# 氣泡到達水面時播放破裂音（低機率，避免太吵）
			if b["y"] < 80 and randf() < 0.25 and AudioManager != null:
				AudioManager.play_sfx(AudioManager.SFX.BUBBLE_POP)
			_bubbles.remove_at(i)
		i -= 1

	# 更新漂浮微粒
	i = _particles.size() - 1
	while i >= 0:
		var p = _particles[i]
		# 緩慢漂移（受水流影響）
		p["x"] += p["drift_x"] * delta + sin(_time * 0.3 + p["phase"]) * 4.0 * delta
		p["y"] += p["drift_y"] * delta + cos(_time * 0.2 + p["phase"] * 1.3) * 2.0 * delta
		# 邊界消失
		if p["x"] < -20 or p["x"] > 1300 or p["y"] < 50 or p["y"] > 740:
			_particles.remove_at(i)
		i -= 1

	# 更新魚群
	i = _fish_schools.size() - 1
	while i >= 0:
		var school = _fish_schools[i]
		school["x"] += school["dir"] * school["speed"] * delta
		# 輕微上下漂移
		school["y"] += sin(_time * 0.4 + school["phase"]) * 8.0 * delta
		# 超出畫面就移除
		if school["x"] < -200 or school["x"] > 1480:
			_fish_schools.remove_at(i)
		i -= 1

	queue_redraw()

func _spawn_bubble(x: float, y: float) -> void:
	_bubbles.append({
		"x": x,
		"y": y,
		"radius": randf_range(3.0, 9.0),
		"speed": randf_range(35.0, 90.0),
		"alpha": randf_range(0.25, 0.65),
		"wobble_phase": randf_range(0.0, TAU),
		"wobble_speed": randf_range(1.5, 3.5),
		"wobble_amp": randf_range(8.0, 18.0),
	})

func _spawn_particle(x: float, y: float) -> void:
	# 浮游生物：微小發光點（藍白色/綠色）
	var is_bioluminescent = randf() < 0.3  # 30% 機率是發光的
	var color: Color
	if is_bioluminescent:
		# 發光浮游生物（藍綠色）
		color = Color(
			randf_range(0.1, 0.4),
			randf_range(0.7, 1.0),
			randf_range(0.6, 1.0),
			randf_range(0.4, 0.8)
		)
	else:
		# 普通塵埃（白色半透明）
		color = Color(0.85, 0.92, 1.0, randf_range(0.1, 0.3))

	_particles.append({
		"x": x,
		"y": y,
		"radius": randf_range(1.0, 3.5) if is_bioluminescent else randf_range(0.8, 2.0),
		"alpha": color.a,
		"drift_x": randf_range(-8.0, 8.0),
		"drift_y": randf_range(-3.0, 3.0),
		"phase": randf_range(0.0, TAU),
		"color": color,
		"glow": is_bioluminescent,
	})

func _spawn_fish_school(x: float, y: float) -> void:
	# 遠景小魚群（5-9 條小魚，半透明，統一方向游動）
	var dir = 1.0 if x < 640 else -1.0
	var fish_count = randi_range(5, 9)
	var offsets: Array = []
	for j in range(fish_count):
		offsets.append(Vector2(
			randf_range(-60, 60),
			randf_range(-25, 25)
		))
	# 魚的顏色（深海魚：藍灰色/銀色）
	var hue = randf_range(0.55, 0.70)  # 藍色系
	var fish_color = Color.from_hsv(hue, randf_range(0.2, 0.5), randf_range(0.5, 0.8), 0.35)

	_fish_schools.append({
		"x": x,
		"y": y,
		"dir": dir,
		"speed": randf_range(40.0, 80.0),
		"fish_count": fish_count,
		"fish_offsets": offsets,
		"color": fish_color,
		"alpha": randf_range(0.2, 0.4),
		"size": randf_range(6.0, 12.0),  # 魚的大小
		"phase": randf_range(0.0, TAU),
	})

func _draw() -> void:
	if not _active:
		return

	# 1. 光線柱（最底層，半透明白色梯形從頂部射入）
	_draw_light_beams()

	# 2. 遠景小魚群（在海草下方，遠景感）
	_draw_fish_schools()

	# 3. 海草（底部，在氣泡下方）
	_draw_seaweeds()

	# 4. 漂浮微粒（浮游生物/塵埃）
	_draw_particles()

	# 5. 氣泡（最上層）
	_draw_bubbles()

func _draw_light_beams() -> void:
	for beam in _light_beams:
		# 光線柱亮度隨時間緩慢脈動
		var pulse = sin(_time * beam["speed"] + beam["phase"]) * 0.5 + 0.5
		var alpha = beam["alpha"] * (0.6 + pulse * 0.4)

		var x: float = beam["x"]
		var w_top: float = beam["width"] * 0.4   # 頂部較窄
		var w_bot: float = beam["width"] * 1.2   # 底部較寬（光線擴散）
		var top_y: float = 0.0
		var bot_y: float = 720.0

		# 用多邊形畫梯形光柱（頂部窄、底部寬）
		var pts = PackedVector2Array([
			Vector2(x - w_top, top_y),
			Vector2(x + w_top, top_y),
			Vector2(x + w_bot, bot_y),
			Vector2(x - w_bot, bot_y),
		])
		# 漸層效果：頂部較亮，底部較暗
		var colors = PackedColorArray([
			Color(0.85, 0.95, 1.0, alpha * 1.2),
			Color(0.85, 0.95, 1.0, alpha * 1.2),
			Color(0.7, 0.88, 1.0, alpha * 0.2),
			Color(0.7, 0.88, 1.0, alpha * 0.2),
		])
		draw_polygon(pts, colors)

func _draw_seaweeds() -> void:
	for sw in _seaweeds:
		var base_x: float = sw["x"]
		var base_y: float = sw["y"]
		var h: float = sw["height"]
		var segs: int = sw["segments"]
		var color: Color = sw["color"]
		var amp: float = sw["amp"]

		# 從底部往上畫分段曲線（每段用 draw_line）
		var prev_pos = Vector2(base_x, base_y)
		for s in range(1, segs + 1):
			var t = float(s) / float(segs)
			# 搖擺：越靠頂部搖擺越大
			var sway = sin(_time * sw["speed"] + sw["phase"] + t * 2.0) * amp * t
			var seg_x = base_x + sway
			var seg_y = base_y - h * t
			var cur_pos = Vector2(seg_x, seg_y)

			# 線條粗細：底部粗、頂部細
			var width = lerp(3.5, 1.0, t)
			# 顏色：底部深、頂部亮
			var seg_color = Color(
				color.r * (0.6 + t * 0.4),
				color.g * (0.7 + t * 0.3),
				color.b * (0.6 + t * 0.4),
				color.a * (0.8 + t * 0.2)
			)
			draw_line(prev_pos, cur_pos, seg_color, width)
			prev_pos = cur_pos

		# 頂部小葉片（橢圓）
		var tip_sway = sin(_time * sw["speed"] + sw["phase"] + 2.0) * amp
		var tip_pos = Vector2(base_x + tip_sway, base_y - h)
		draw_circle(tip_pos, 4.0, Color(color.r * 1.2, color.g * 1.3, color.b * 1.0, color.a * 0.9))

func _draw_particles() -> void:
	for p in _particles:
		var pos = Vector2(p["x"], p["y"])
		var r: float = p["radius"]
		var color: Color = p["color"]
		# 閃爍效果（發光微粒）
		var flicker = 1.0
		if p["glow"]:
			flicker = 0.6 + sin(_time * 3.0 + p["phase"]) * 0.4
			# 發光暈圈
			draw_circle(pos, r * 2.5, Color(color.r, color.g, color.b, color.a * 0.15 * flicker))
			draw_circle(pos, r * 1.5, Color(color.r, color.g, color.b, color.a * 0.3 * flicker))
		# 主體
		draw_circle(pos, r, Color(color.r, color.g, color.b, color.a * flicker))

func _draw_fish_schools() -> void:
	for school in _fish_schools:
		var base_x: float = school["x"]
		var base_y: float = school["y"]
		var dir: float = school["dir"]
		var color: Color = school["color"]
		var sz: float = school["size"]

		for j in range(school["fish_count"]):
			var offset: Vector2 = school["fish_offsets"][j]
			# 每條魚有輕微的上下擺動（模擬游泳）
			var swim_y = sin(_time * 1.8 + school["phase"] + j * 0.7) * 4.0
			var fx = base_x + offset.x
			var fy = base_y + offset.y + swim_y

			# 畫小魚（橢圓身體 + 三角形尾巴）
			_draw_small_fish(Vector2(fx, fy), sz, dir, color)

func _draw_small_fish(pos: Vector2, size: float, dir: float, color: Color) -> void:
	# 身體（橢圓，用多邊形近似）
	var body_pts = PackedVector2Array()
	var body_steps = 8
	for k in range(body_steps):
		var angle = TAU * k / body_steps
		var bx = pos.x + cos(angle) * size * dir
		var by = pos.y + sin(angle) * size * 0.45
		body_pts.append(Vector2(bx, by))
	draw_polygon(body_pts, PackedColorArray([color] * body_steps))

	# 尾巴（三角形）
	var tail_x = pos.x - size * dir * 0.8
	var tail_pts = PackedVector2Array([
		Vector2(tail_x, pos.y),
		Vector2(tail_x - size * dir * 0.7, pos.y - size * 0.5),
		Vector2(tail_x - size * dir * 0.7, pos.y + size * 0.5),
	])
	draw_polygon(tail_pts, PackedColorArray([
		Color(color.r, color.g, color.b, color.a * 0.8),
		Color(color.r, color.g, color.b, color.a * 0.5),
		Color(color.r, color.g, color.b, color.a * 0.5),
	]))

func _draw_bubbles() -> void:
	for b in _bubbles:
		var pos = Vector2(b["x"], b["y"])
		var r: float = b["radius"]
		var a: float = b["alpha"]
		# 氣泡外圈（淡藍白色）
		draw_arc(pos, r, 0.0, TAU, 16, Color(0.75, 0.92, 1.0, a * 0.85), 1.2)
		# 高光（左上角小白點，讓氣泡有立體感）
		draw_circle(pos + Vector2(-r * 0.3, -r * 0.35), r * 0.22, Color(1.0, 1.0, 1.0, a * 0.55))
		# 大氣泡加內部反光（r > 6 才加）
		if r > 6.0:
			draw_arc(pos, r * 0.6, PI * 0.8, PI * 1.6, 10, Color(0.85, 0.95, 1.0, a * 0.25), 0.8)
