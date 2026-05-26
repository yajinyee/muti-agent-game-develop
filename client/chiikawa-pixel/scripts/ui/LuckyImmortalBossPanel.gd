## LuckyImmortalBossPanel.gd — T155 幸運不死 BOSS 魚 Lucky Panel
## lucky-panel-agent 負責維護
## 業界依據：Royal Fishing「Immortal Boss mechanic — consecutive wins 50X-150X until they leave the screen」
## 不死 BOSS 主題：深紅色 + 生命條 + 倍率遞增顯示 + 完美不死演出
extends CanvasLayer

const THEME_COLOR = Color(0.9, 0.1, 0.1)           # 深紅色
const DARK_COLOR  = Color(0.15, 0.0, 0.0)           # 深紅黑
const GOLD_COLOR  = Color(1.0, 0.85, 0.0)           # 金色

var _banner: Control = null
var _lives_label: Label = null
var _mult_label: Label = null
var _lives_left: int = 5

func _ready() -> void:
	layer = 50
	_create_ui()
	GameManager.lucky_immortal_boss.connect(_on_lucky_immortal_boss)

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
	title.text = "💀 不死 BOSS 降臨！"
	title.add_theme_color_override("font_color", THEME_COLOR)
	title.add_theme_font_size_override("font_size", 20)
	title.position = Vector2(20, 12)
	_banner.add_child(title)

	_lives_label = Label.new()
	_lives_label.text = "生命：❤️❤️❤️❤️❤️"
	_lives_label.add_theme_color_override("font_color", THEME_COLOR)
	_lives_label.add_theme_font_size_override("font_size", 16)
	_lives_label.position = Vector2(320, 16)
	_banner.add_child(_lives_label)

	_mult_label = Label.new()
	_mult_label.text = "×2.0"
	_mult_label.add_theme_color_override("font_color", GOLD_COLOR)
	_mult_label.add_theme_font_size_override("font_size", 20)
	_mult_label.position = Vector2(620, 12)
	_banner.add_child(_mult_label)

	var hint = Label.new()
	hint.text = "每次擊破倍率 +0.5x | 耗盡 5 條命 → 完美不死"
	hint.add_theme_color_override("font_color", Color(1.0, 0.6, 0.6))
	hint.add_theme_font_size_override("font_size", 13)
	hint.position = Vector2(720, 18)
	_banner.add_child(hint)

func _on_lucky_immortal_boss(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"immortal_spawn":
			_lives_left = 5
			_show_banner()
			_flash_screen(THEME_COLOR, 3)
			ScreenShake.add_trauma(0.6)
			_update_lives_display()
		"immortal_kill":
			_lives_left = data.get("lives_left", _lives_left)
			var mult = data.get("current_mult", 2.0)
			var reward = data.get("reward", 0)
			_update_lives_display()
			if _mult_label:
				_mult_label.text = "×%.1f" % mult
			_show_float_text("💀 +%d  ×%.1f" % [reward, mult], THEME_COLOR)
			ScreenShake.add_trauma(0.3)
		"immortal_timeout":
			_hide_banner()
		"immortal_perfect":
			_show_perfect(data.get("total_reward", 0), data.get("boost_mult", 5.0), data.get("boost_secs", 12))
		"immortal_perfect_end":
			_hide_banner()

func _update_lives_display() -> void:
	if not _lives_label:
		return
	var hearts = ""
	for i in range(_lives_left):
		hearts += "❤️"
	for i in range(5 - _lives_left):
		hearts += "🖤"
	_lives_label.text = "生命：" + hearts

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

func _show_perfect(total: int, mult: float, secs: int) -> void:
	_flash_screen(GOLD_COLOR, 3)
	ScreenShake.add_trauma(0.8)
	var popup = Label.new()
	popup.text = "💀✨ 完美不死！耗盡 5 條命！獲得 %d！全服 ×%.1f 加成 %d 秒！" % [total, mult, secs]
	popup.add_theme_color_override("font_color", GOLD_COLOR)
	popup.add_theme_font_size_override("font_size", 20)
	popup.position = Vector2(150, 300)
	add_child(popup)
	var tween = create_tween()
	tween.tween_property(popup, "position:y", 260.0, 0.5)
	tween.tween_interval(2.5)
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
