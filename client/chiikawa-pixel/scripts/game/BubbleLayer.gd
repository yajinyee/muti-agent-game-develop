## BubbleLayer.gd
## 海底氣泡動畫層，在 normal 背景時顯示浮動氣泡

extends Node2D

var _bubbles: Array = []
var _spawn_timer: float = 0.0
var _spawn_interval: float = 0.35
var _active: bool = true

func _ready() -> void:
	# 預生成氣泡，讓畫面一開始就有
	for i in range(10):
		_spawn_bubble(randf_range(50, 1230), randf_range(80, 680))

func set_active(v: bool) -> void:
	_active = v
	if not v:
		_bubbles.clear()
		queue_redraw()

func _process(delta: float) -> void:
	if not _active:
		return

	# 生成新氣泡
	_spawn_timer += delta
	if _spawn_timer >= _spawn_interval:
		_spawn_timer = 0.0
		_spawn_bubble(randf_range(50, 1230), 730.0)

	# 更新氣泡位置和透明度
	var time = Time.get_ticks_msec() / 1000.0
	var i = _bubbles.size() - 1
	while i >= 0:
		var b = _bubbles[i]
		b["y"] -= b["speed"] * delta
		# 左右搖擺（正弦波）
		b["x"] += sin(time * b["wobble_speed"] + b["wobble_phase"]) * b["wobble_amp"] * delta
		# 接近頂部時淡出
		if b["y"] < 80:
			b["alpha"] -= delta * 1.8
		if b["alpha"] <= 0.0 or b["y"] < -10:
			_bubbles.remove_at(i)
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

func _draw() -> void:
	for b in _bubbles:
		var pos = Vector2(b["x"], b["y"])
		var r: float = b["radius"]
		var a: float = b["alpha"]
		# 氣泡外圈（淡藍白色）
		draw_arc(pos, r, 0.0, TAU, 16, Color(0.75, 0.92, 1.0, a * 0.85), 1.2)
		# 高光（左上角小白點，讓氣泡有立體感）
		draw_circle(pos + Vector2(-r * 0.3, -r * 0.35), r * 0.22, Color(1.0, 1.0, 1.0, a * 0.55))
