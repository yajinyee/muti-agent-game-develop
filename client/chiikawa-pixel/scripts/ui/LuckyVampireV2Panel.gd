## LuckyVampireV2Panel.gd — T152 幸運吸血鬼升級魚 Lucky Panel
## lucky-panel-agent 負責維護
## 業界依據：Jili「Vampire multiplier increases the more you fight, up to X5」升級版（最高 ×10.0）
## 吸血鬼升級主題：深紫色 + 吸收計數器 + 倍率顯示 + 完美吸血演出
extends CanvasLayer

const THEME_COLOR = Color(0.6, 0.0, 0.8)          # 深紫色
const DARK_COLOR  = Color(0.1, 0.0, 0.15)          # 深紫黑
const GOLD_COLOR  = Color(1.0, 0.85, 0.0)          # 金色

var _banner: Control = null
var _absorb_label: Label = null
var _mult_label: Label = null
var _absorb_count: int = 0

func _ready() -> void:
	layer = 47
	_create_ui()
	GameManager.lucky_vampire_v2.connect(_on_lucky_vampire_v2)

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
	title.text = "🧛 吸血鬼升級模式！"
	title.add_theme_color_override("font_color", THEME_COLOR)
	title.add_theme_font_size_override("font_size", 20)
	title.position = Vector2(20, 12)
	_banner.add_child(title)

	_absorb_label = Label.new()
	_absorb_label.text = "吸收：0 次"
	_absorb_label.add_theme_color_override("font_color", GOLD_COLOR)
	_absorb_label.add_theme_font_size_override("font_size", 16)
	_absorb_label.position = Vector2(320, 16)
	_banner.add_child(_absorb_label)

	_mult_label = Label.new()
	_mult_label.text = "×1.0"
	_mult_label.add_theme_color_override("font_color", Color(0.9, 0.5, 1.0))
	_mult_label.add_theme_font_size_override("font_size", 20)
	_mult_label.position = Vector2(520, 12)
	_banner.add_child(_mult_label)

	var hint = Label.new()
	hint.text = "每次擊破 +1.5x | 最高 ×10.0 | 吸收 ≥10 → 完美"
	hint.add_theme_color_override("font_color", Color(0.8, 0.6, 1.0))
	hint.add_theme_font_size_override("font_size", 13)
	hint.position = Vector2(640, 18)
	_banner.add_child(hint)

func _on_lucky_vampire_v2(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"vampire_v2_start":
			_absorb_count = 0
			_show_banner()
			_flash_screen(THEME_COLOR, 2)
			ScreenShake.add_trauma(0.3)
		"absorb_v2":
			_absorb_count = data.get("absorb_count", _absorb_count)
			var mult = data.get("current_mult", 1.0)
			if _absorb_label:
				_absorb_label.text = "吸收：%d 次" % _absorb_count
			if _mult_label:
				_mult_label.text = "×%.1f" % mult
		"mult_mode_v2":
			var mult = data.get("current_mult", 1.0)
			_flash_screen(THEME_COLOR, 2)
			_show_float_text("🧛 倍率模式！×%.1f！" % mult, THEME_COLOR)
		"vampire_v2_end":
			_hide_banner()
		"vampire_v2_perfect":
			_show_perfect(data.get("absorb_count", 0), data.get("final_mult", 1.0), data.get("boost_mult", 4.0), data.get("boost_secs", 10))
		"vampire_v2_perfect_end":
			_hide_banner()

func _show_banner() -> void:
	if _banner:
		_banner.visible = true
		_banner.modulate.a = 0.0
		var tween = create_tween()
		tween.tween_property(_banner, "modulate:a", 1.0, 0.3)

func _hide_banner() -> void:
	if _banner:
		var tween = create_tween()
		tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
		tween.tween_callback(func(): _banner.visible = false)

func _show_perfect(absorbs: int, final_mult: float, boost_mult: float, secs: int) -> void:
	_flash_screen(GOLD_COLOR, 3)
	ScreenShake.add_trauma(0.6)
	var popup = Label.new()
	popup.text = "🧛✨ 完美吸血！吸收 %d 次！×%.1f！全服 ×%.1f 加成 %d 秒！" % [absorbs, final_mult, boost_mult, secs]
	popup.add_theme_color_override("font_color", GOLD_COLOR)
	popup.add_theme_font_size_override("font_size", 20)
	popup.position = Vector2(150, 300)
	add_child(popup)
	var tween = create_tween()
	tween.tween_property(popup, "position:y", 260.0, 0.5)
	tween.tween_interval(2.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)

func _show_float_text(text: String, color: Color) -> void:
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 18)
	lbl.position = Vector2(randf_range(400, 800), 400)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 60, 0.6)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.6)
	tween.tween_callback(lbl.queue_free)

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
