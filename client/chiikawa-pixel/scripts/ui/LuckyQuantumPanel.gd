## LuckyQuantumPanel.gd — T146 幸運量子魚 Lucky Panel
## lucky-panel-agent 負責維護
## 業界依據：量子力學 Crash mechanic 升級版 — 量子疊加態觀測機制
## 量子主題：青藍色 + 量子粒子特效 + 觀測計數器
extends CanvasLayer

const THEME_COLOR = Color(0.0, 0.9, 1.0)       # 青藍色
const DARK_COLOR  = Color(0.0, 0.2, 0.3)        # 深青藍
const GOLD_COLOR  = Color(1.0, 0.85, 0.0)       # 金色

var _banner: Control = null
var _observe_label: Label = null
var _observed_count: int = 0

func _ready() -> void:
	layer = 41
	_create_ui()
	GameManager.lucky_quantum.connect(_on_lucky_quantum)

func _create_ui() -> void:
	_banner = Control.new()
	_banner.visible = false
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 56)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.color = Color(DARK_COLOR.r, DARK_COLOR.g, DARK_COLOR.b, 0.88)
	bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	_banner.add_child(bg)

	var top_line = ColorRect.new()
	top_line.color = THEME_COLOR
	top_line.size = Vector2(1280, 3)
	top_line.position = Vector2(0, 0)
	_banner.add_child(top_line)

	var title = Label.new()
	title.text = "⚛️ 量子觀測！"
	title.add_theme_color_override("font_color", THEME_COLOR)
	title.add_theme_font_size_override("font_size", 20)
	title.position = Vector2(20, 12)
	_banner.add_child(title)

	_observe_label = Label.new()
	_observe_label.text = "觀測：0 個"
	_observe_label.add_theme_color_override("font_color", GOLD_COLOR)
	_observe_label.add_theme_font_size_override("font_size", 16)
	_observe_label.position = Vector2(300, 16)
	_banner.add_child(_observe_label)

func _on_lucky_quantum(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"quantum_observe":
			_observed_count = 0
			_show_banner()
			_flash_screen(THEME_COLOR, 2)
		"quantum_result":
			_observed_count = data.get("observed_count", _observed_count)
			if _observe_label:
				_observe_label.text = "觀測：%d 個" % _observed_count
		"quantum_collapse":
			_show_perfect(data.get("boost_mult", 5.5), data.get("boost_sec", 12))
		"quantum_collapse_end":
			_hide_banner()

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
	popup.text = "⚛️ 量子坍縮！全服 ×%.1f 加成 %d 秒！" % [mult, secs]
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
