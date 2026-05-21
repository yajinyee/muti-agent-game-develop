## ShockwaveRing.gd
## 衝擊波環 — 用 _draw 繪製真正的圓形衝擊波
## 由 HitEffect.gd 動態建立

extends Node2D

var ring_color: Color = Color.WHITE
var ring_radius: float = 16.0

func _draw() -> void:
	# 主衝擊波環（粗圓弧）
	draw_arc(Vector2.ZERO, ring_radius, 0.0, TAU, 32, Color(ring_color.r, ring_color.g, ring_color.b, 0.65), ring_radius * 0.18)
	# 內側細環（增加層次感）
	draw_arc(Vector2.ZERO, ring_radius * 0.7, 0.0, TAU, 24, Color(ring_color.r, ring_color.g, ring_color.b, 0.35), ring_radius * 0.06)
