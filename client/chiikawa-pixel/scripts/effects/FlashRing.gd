## FlashRing.gd
## 命中閃光環 — 用 _draw 繪製真正的圓形（比 ColorRect 更精確）
## 由 HitEffect.gd 動態建立

extends Node2D

var ring_color: Color = Color.WHITE
var ring_radius: float = 24.0

func _draw() -> void:
	# 中心實心圓（高亮核心）
	draw_circle(Vector2.ZERO, ring_radius * 0.28, Color(ring_color.r, ring_color.g, ring_color.b, 0.95))
	# 外環（空心圓弧，像素風格用多段線模擬）
	draw_arc(Vector2.ZERO, ring_radius * 0.55, 0.0, TAU, 24, Color(ring_color.r, ring_color.g, ring_color.b, 0.75), ring_radius * 0.12)
	# 外圈光暈（更大、更透明）
	draw_arc(Vector2.ZERO, ring_radius * 0.85, 0.0, TAU, 20, Color(ring_color.r, ring_color.g, ring_color.b, 0.35), ring_radius * 0.08)
	# 4 個方向的光芒射線
	for i in 4:
		var angle = i * PI / 2.0 + PI / 4.0
		var inner = Vector2(cos(angle), sin(angle)) * ring_radius * 0.35
		var outer = Vector2(cos(angle), sin(angle)) * ring_radius * 0.9
		draw_line(inner, outer, Color(ring_color.r, ring_color.g, ring_color.b, 0.6), ring_radius * 0.06)
