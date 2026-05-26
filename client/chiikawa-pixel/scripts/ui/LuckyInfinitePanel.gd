## LuckyInfinitePanel.gd — T148 幸運無限魚 Lucky Panel
## lucky-panel-agent 負責維護
## 業界依據：Infinite multiplier accumulation — 每次擊破倍率無限累積
## 無限主題：紫色 + 倍率計數器 + 無限符號
extends CanvasLayer

const THEME_COLOR = Color(0.6, 0.2, 1.0)        # 紫色
const DARK_COLOR  = Color(0.15, 0.05, 0.25)      # 深紫
const GOLD_COLOR  = Color(1.0, 0.85, 0.0)        # 金色

var _banner: Control = null
var _mult_label: Label = null
var _kill_label: Label = null

func _ready() -> void:
	layer = 43
	_create_ui()
	GameManager.lucky_infinite.connect(_on_lucky_infinite)

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
	title.text = "♾️ 無限模式！"
	title.add_theme_color_override("font_color", THEME_COLOR)
	title.add_theme_font_size_override("font_size", 20)
	title.position = Vector2(20, 12)
	_banner.add_child(title)

	_mult_label = Label.new()
	_mult_label.text = "×1.0"
	_mult_label.add_theme_color_override("font_color", GOLD_COLOR)
	_mult_label.add_theme_font_size_override("font_size", 22)
	_mult_label.position = Vector2(280, 10)
	_banner.add_child(_mult_label)

	_kill_label = Label.new()
	_kill_label.text = "0 次"
	_kill_label.add_theme_color_override("font_color", Color(0.8, 0.6, 1.0))
	_kill_label.add_theme_font_size_override("font_size", 14)
	_kill_label.position = Vector2(380, 18)
	_banner.add_child(_kill_label)

func _on_lucky_infinite(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"infinite_start":
			_show_banner()
			_flash_screen(THEME_COLOR, 2)
		"infinite_kill":
			var mult = data.get("accum_mult", 1.0)
			var kills = data.get("kill_count", 0)
			if _mult_label:
				_mult_label.text = "×%.1f" % mult
				# 倍率越高顏色越亮
				if mult >= 20.0:
					_mult_label.modulate = Color(1.0, 0.3, 0.1)
				elif mult >= 10.0:
					_mult_label.modulate = Color(1.0, 0.6, 0.1)
				else:
					_mult_label.modulate = GOLD_COLOR
			if _kill_label:
				_kill_label.text = "%d 次" % kills
		"infinite_end":
			_hide_banner()
		"infinite_perfect":
			_show_perfect(data.get("accum_mult", 1.0), data.get("boost_mult", 6.0), data.get("boost_sec", 15))
		"infinite_perfect_end":
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

func _show_perfect(accum_mult: float, boost_mult: float, secs: int) -> void:
	_flash_screen(GOLD_COLOR, 3)
	var popup = Label.new()
	popup.text = "♾️ 無限完美！累積 ×%.1f！全服 ×%.0f 加成 %d 秒！" % [accum_mult, boost_mult, secs]
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
