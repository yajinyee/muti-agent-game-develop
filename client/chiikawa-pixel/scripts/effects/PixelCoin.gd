## PixelCoin.gd — 像素金幣特效
## hit-effect-agent 負責維護
## 用於金幣雨、黃金雨等特效的可重用金幣節點
extends Node2D

const COIN_COLORS = [
	Color(1.0, 0.85, 0.0),   # 金色
	Color(1.0, 0.9, 0.3),    # 亮金
	Color(0.9, 0.7, 0.0),    # 深金
]

var _rect: ColorRect = null
var _label: Label = null
var _tween: Tween = null

func _ready() -> void:
	_rect = ColorRect.new()
	_rect.size = Vector2(14, 14)
	_rect.position = Vector2(-7, -7)
	_rect.color = COIN_COLORS[randi() % COIN_COLORS.size()]
	add_child(_rect)

	# 金幣符號
	_label = Label.new()
	_label.text = "¥"
	_label.position = Vector2(-6, -8)
	_label.add_theme_font_size_override("font_size", 10)
	_label.modulate = Color(0.6, 0.4, 0.0)
	add_child(_label)

## 播放金幣飛向目標的動畫
func fly_to(target: Vector2, duration: float = 0.5, delay: float = 0.0) -> void:
	if _tween != null and _tween.is_valid():
		_tween.kill()
	scale = Vector2.ZERO
	_tween = create_tween()
	if delay > 0:
		_tween.tween_interval(delay)
	# 彈出
	_tween.tween_property(self, "scale", Vector2(1.2, 1.2), 0.1).set_ease(Tween.EASE_OUT)
	_tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.05)
	# 飛行
	_tween.tween_property(self, "position", target, duration).set_ease(Tween.EASE_IN).set_trans(Tween.TRANS_QUAD)
	_tween.parallel().tween_property(self, "scale", Vector2(0.3, 0.3), duration)
	_tween.tween_callback(func(): queue_free())

## 播放金幣掉落動畫
func drop_from(start: Vector2, end: Vector2, duration: float = 0.6) -> void:
	position = start
	if _tween != null and _tween.is_valid():
		_tween.kill()
	scale = Vector2.ZERO
	_tween = create_tween()
	_tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.1).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	_tween.tween_property(self, "position", end, duration * 0.6).set_ease(Tween.EASE_OUT)
	_tween.tween_property(self, "position:y", end.y + 30, duration * 0.4).set_ease(Tween.EASE_IN)
	_tween.parallel().tween_property(self, "modulate:a", 0.0, duration * 0.4)
	_tween.tween_callback(func(): queue_free())

## 靜態工廠方法：在場景中生成金幣特效
static func spawn_coin_rain(parent: Node, center: Vector2, count: int = 8) -> void:
	for i in count:
		var coin = PixelCoin.new()
		coin.position = center
		coin.z_index = 45
		parent.add_child(coin)
		var angle = randf_range(-PI * 0.8, -PI * 0.2)  # 向上散射
		var dist = randf_range(50, 120)
		var end = center + Vector2(cos(angle), sin(angle)) * dist
		coin.drop_from(center, end, randf_range(0.4, 0.7))

static func spawn_coin_collect(parent: Node, from: Vector2, to: Vector2, count: int = 5) -> void:
	for i in count:
		var coin = PixelCoin.new()
		var offset = Vector2(randf_range(-30, 30), randf_range(-30, 30))
		coin.position = from + offset
		coin.z_index = 45
		parent.add_child(coin)
		coin.fly_to(to, randf_range(0.3, 0.6), float(i) * 0.05)
