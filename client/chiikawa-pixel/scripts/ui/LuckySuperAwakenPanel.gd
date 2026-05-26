## LuckySuperAwakenPanel.gd — T153 幸運超級覺醒魚 Lucky Panel
## lucky-panel-agent 負責維護
## 業界依據：Jili「Super Awakening Performance, bonus up to 3000x」
## 超級覺醒主題：火橙色 + 全場審判演出 + 全服 ×7.0 加成
extends CanvasLayer

const THEME_COLOR = Color(1.0, 0.4, 0.0)          # 火橙色
const DARK_COLOR  = Color(0.2, 0.05, 0.0)          # 深棕黑
const GOLD_COLOR  = Color(1.0, 0.85, 0.0)          # 金色

var _banner: Control = null
var _boost_label: Label = null

func _ready() -> void:
	layer = 48
	_create_ui()
	GameManager.lucky_super_awaken.connect(_on_lucky_super_awaken)

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
	title.text = "⚡ 超級覺醒！全服 ×7.0！"
	title.add_theme_color_override("font_color", THEME_COLOR)
	title.add_theme_font_size_override("font_size", 20)
	title.position = Vector2(20, 12)
	_banner.add_child(title)

	_boost_label = Label.new()
	_boost_label.text = "全服加成 ×7.0"
	_boost_label.add_theme_color_override("font_color", GOLD_COLOR)
	_boost_label.add_theme_font_size_override("font_size", 18)
	_boost_label.position = Vector2(480, 14)
	_banner.add_child(_boost_label)

func _on_lucky_super_awaken(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"super_awaken_start":
			_flash_screen(THEME_COLOR, 3)
			ScreenShake.add_trauma(0.7)
			_show_float_text("⚡ 超級覺醒！全場審判！", THEME_COLOR)
		"super_awaken_result":
			var hits = data.get("hit_count", 0)
			var reward = data.get("total_reward", 0)
			_show_float_text("⚡ 審判 %d 條！+%d！" % [hits, reward], GOLD_COLOR)
		"super_awaken_boost":
			_show_banner()
			_flash_screen(GOLD_COLOR, 3)
			ScreenShake.add_trauma(0.8)
		"super_awaken_boost_end":
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

func _show_float_text(text: String, color: Color) -> void:
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 24)
	lbl.position = Vector2(300, 300)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "scale", Vector2(1.3, 1.3), 0.1)
	tween.tween_property(lbl, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_property(lbl, "position:y", lbl.position.y - 80, 0.8)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.8)
	tween.tween_callback(lbl.queue_free)

func _flash_screen(color: Color, times: int) -> void:
	var overlay = ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, 0.0)
	overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tween = create_tween()
	for i in range(times):
		tween.tween_property(overlay, "color:a", 0.45, 0.08)
		tween.tween_property(overlay, "color:a", 0.0, 0.12)
	tween.tween_callback(overlay.queue_free)
