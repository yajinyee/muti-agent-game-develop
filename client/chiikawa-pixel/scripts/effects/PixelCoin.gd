## PixelCoin.gd
## 像素金幣節點（用 _draw 繪製帶高光的金幣）
## 由 TargetManager 的 T105 金幣雨使用

extends Node2D

func _draw() -> void:
	# 金幣主體（金色圓形）
	draw_circle(Vector2.ZERO, 6.0, Color(1.0, 0.82, 0.0))
	# 金幣邊框（深金色）
	draw_arc(Vector2.ZERO, 6.0, 0.0, TAU, 16, Color(0.75, 0.55, 0.0), 1.5)
	# 高光（左上角白色小圓）
	draw_circle(Vector2(-2.0, -2.0), 1.8, Color(1.0, 1.0, 0.9, 0.8))
	# 中心 ¥ 符號（用小矩形模擬）
	draw_rect(Rect2(-1.0, -2.5, 2.0, 5.0), Color(0.75, 0.55, 0.0, 0.7))
	draw_rect(Rect2(-3.0, -0.5, 6.0, 1.0), Color(0.75, 0.55, 0.0, 0.7))
