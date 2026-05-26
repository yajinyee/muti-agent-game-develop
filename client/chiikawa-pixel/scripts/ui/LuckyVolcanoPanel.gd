## LuckyVolcanoPanel.gd — T143 幸運火山魚 Lucky Panel
## lucky-panel-agent 負責維護
## 業界依據：Jili Games 2026「Volcano eruption」
## 火山主題：火紅色 + 熔岩彈計數器 + 命中指示
extends CanvasLayer

const THEME_COLOR = Color(1.0, 0.27, 0.0)       # 火紅色
const DARK_COLOR  = Color(0.35, 0.08, 0.0)       # 深紅
const GOLD_COLOR  = Color(1.0, 0.85, 0.0)        # 金色

var _banner: Control = null
var _bomb_label: Label = null
var _hit_label: Label = null
var _bomb_count: int = 0
var _hit_bombs: int = 0

func _ready() -> void:
	layer = 38
	_create_ui()
	GameManager.lucky_volcano.connect(_on_lucky_volcano)

func _create_ui() -> void:
	_banner = Control.new()
	_banner.visible = false
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 60)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.color = Color(DARK_COLOR.r, DARK_COLOR.g, DARK_COLOR.b, 0.88)
	bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	_banner.add_child(bg)

	var top_line = ColorRect.new()
	top_line.color = THEME_COLOR
	top_line.size = Vector2(1280, 3)
	_banner.add_child(top_line)

	var title = Label.new()
	title.text = "🌋 火山爆發！"
	title.add_theme_color_override("font_color", THEME_COLOR)
	title.add_theme_font_size_override("font_size", 20)
	title.position = Vector2(20, 12)
	_banner.add_child(title)

	_bomb_label = Label.new()
	_bomb_label.text = "熔岩彈：0/10"
	_bomb_label.add_theme_color_override("font_color", THEME_COLOR)
	_bomb_label.add_theme_font_size_override("font_size", 16)
	_bomb_label.position = Vector2(200, 16)
	_banner.add_child(_bomb_label)

	_hit_label = Label.new()
	_hit_label.text = "命中：0"
	_hit_label.add_theme_color_override("font_color", GOLD_COLOR)
	_hit_label.add_theme_font_size_override("font_size", 16)
	_hit_label.position = Vector2(380, 16)
	_banner.add_child(_hit_label)

func _on_lucky_volcano(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"volcano_erupt":
			_bomb_count = 0
			_hit_bombs = 0
			_show_banner()
			_flash_screen(THEME_COLOR, 3)
		"lava_bomb":
			_bomb_count = data.get("bomb_num", _bomb_count)
			_hit_bombs = data.get("hit_bombs", _hit_bombs)
			if _bomb_label:
				_bomb_label.text = "熔岩彈：%d/10" % _bomb_count
			if _hit_label:
				_hit_label.text = "命中：%d" % _hit_bombs
			# 在落點顯示爆炸特效
			_show_bomb_effect(data.get("bomb_x", 0.0), data.get("bomb_y", 0.0))
		"volcano_end":
			_hide_banner()
		"volcano_perfect":
			_show_perfect(data.get("boost_mult", 4.2), data.get("boost_sec", 10))
		"volcano_perfect_end":
			_hide_banner()

func _show_bomb_effect(bx: float, by: float) -> void:
	if bx <= 0 and by <= 0:
		return
	var dot = ColorRect.new()
	dot.size = Vector2(24, 24)
	dot.position = Vector2(bx - 12, by - 12)
	dot.color = THEME_COLOR
	add_child(dot)
	var tween = create_tween()
	tween.tween_property(dot, "scale", Vector2(2.5, 2.5), 0.15)
	tween.tween_property(dot, "modulate:a", 0.0, 0.2)
	tween.tween_callback(dot.queue_free)

func _show_banner() -> void:
	if _banner:
		_banner.visible = true
		var tween = create_tween()
		_banner.modulate.a = 0.0
		tween.tween_property(_banner, "modulate:a", 1.0, 0.3)

func _hide_banner() -> void:
	if _banner:
		var tween = create_tween()
		tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
		tween.tween_callback(func(): _banner.visible = false)

func _show_perfect(mult: float, secs: int) -> void:
	_flash_screen(GOLD_COLOR, 3)
	var popup = Label.new()
	popup.text = "🌋 完美火山！全服 ×%.1f 加成 %d 秒！" % [mult, secs]
	popup.add_theme_color_override("font_color", GOLD_COLOR)
	popup.add_theme_font_size_override("font_size", 22)
	popup.position = Vector2(200, 300)
	add_child(popup)
	var tween = create_tween()
	tween.tween_property(popup, "position:y", 260.0, 0.5)
	tween.tween_interval(2.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)

func _flash_screen(color: Color, times: int) -> void:
	var overlay = ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, 0.0)
	overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tween = create_tween()
	for i in range(times):
		tween.tween_property(overlay, "color:a", 0.35, 0.08)
		tween.tween_property(overlay, "color:a", 0.0, 0.12)
	tween.tween_callback(overlay.queue_free)
