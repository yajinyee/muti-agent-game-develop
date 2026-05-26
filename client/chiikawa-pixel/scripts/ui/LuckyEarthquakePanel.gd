## LuckyEarthquakePanel.gd — T142 幸運地震魚 Lucky Panel
## lucky-panel-agent 負責維護
## 業界依據：Fishing Fortune 2026「Earthquake shockwave」
## 地震主題：棕橙色 + 三波指示器 + 命中計數
extends CanvasLayer

const THEME_COLOR = Color(0.8, 0.4, 0.1)       # 棕橙色
const DARK_COLOR  = Color(0.3, 0.15, 0.05)      # 深棕
const GOLD_COLOR  = Color(1.0, 0.85, 0.0)       # 金色

var _banner: Control = null
var _wave_dots: Array = []
var _hit_label: Label = null
var _total_hit: int = 0

func _ready() -> void:
	layer = 37
	_create_ui()
	GameManager.lucky_earthquake.connect(_on_lucky_earthquake)

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
	title.text = "🌍 地震波！"
	title.add_theme_color_override("font_color", THEME_COLOR)
	title.add_theme_font_size_override("font_size", 20)
	title.position = Vector2(20, 12)
	_banner.add_child(title)

	# 三波指示點
	for i in range(3):
		var dot = ColorRect.new()
		dot.size = Vector2(20, 20)
		dot.position = Vector2(200 + i * 30, 18)
		dot.color = Color(0.4, 0.4, 0.4)
		_banner.add_child(dot)
		_wave_dots.append(dot)

	_hit_label = Label.new()
	_hit_label.text = "命中：0"
	_hit_label.add_theme_color_override("font_color", GOLD_COLOR)
	_hit_label.add_theme_font_size_override("font_size", 16)
	_hit_label.position = Vector2(320, 16)
	_banner.add_child(_hit_label)

func _on_lucky_earthquake(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"quake_warning":
			_total_hit = 0
			_show_banner()
			_flash_screen(THEME_COLOR, 2)
		"quake_wave":
			var wave_num = data.get("wave_num", 1) - 1
			_total_hit = data.get("total_hit_count", _total_hit)
			if wave_num < _wave_dots.size():
				_wave_dots[wave_num].color = THEME_COLOR
			if _hit_label:
				_hit_label.text = "命中：%d" % _total_hit
		"quake_end":
			_hide_banner()
		"quake_perfect":
			_show_perfect(data.get("boost_mult", 4.0), data.get("boost_sec", 9))
		"quake_perfect_end":
			_hide_banner()

func _show_banner() -> void:
	if _banner:
		_banner.visible = true
		for dot in _wave_dots:
			dot.color = Color(0.4, 0.4, 0.4)
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
	popup.text = "🌍 完美地震！全服 ×%.1f 加成 %d 秒！" % [mult, secs]
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
