## LuckyAwakenedCrocPanel.gd — T151 幸運覺醒鱷魚 Lucky Panel
## lucky-panel-agent 負責維護
## 業界依據：Jili「Giant Crocodile awakens to hunt fish on the fish farm to accumulate big prizes」
## 覺醒鱷魚主題：深綠色 + 鱷魚獵魚計數器 + 完美覺醒演出
extends CanvasLayer

const THEME_COLOR = Color(0.0, 0.8, 0.3)          # 深綠色
const DARK_COLOR  = Color(0.0, 0.15, 0.05)         # 深綠黑
const GOLD_COLOR  = Color(1.0, 0.85, 0.0)          # 金色

var _banner: Control = null
var _hunt_label: Label = null
var _hunt_count: int = 0

func _ready() -> void:
	layer = 46
	_create_ui()
	GameManager.lucky_awakened_croc.connect(_on_lucky_awakened_croc)

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
	title.text = "🐊 覺醒鱷魚獵魚中！"
	title.add_theme_color_override("font_color", THEME_COLOR)
	title.add_theme_font_size_override("font_size", 20)
	title.position = Vector2(20, 12)
	_banner.add_child(title)

	_hunt_label = Label.new()
	_hunt_label.text = "獵魚：0 條"
	_hunt_label.add_theme_color_override("font_color", GOLD_COLOR)
	_hunt_label.add_theme_font_size_override("font_size", 16)
	_hunt_label.position = Vector2(320, 16)
	_banner.add_child(_hunt_label)

	var hint = Label.new()
	hint.text = "每次獵魚獎勵 ×3.0 | 獵魚 ≥8 → 完美覺醒"
	hint.add_theme_color_override("font_color", Color(0.6, 1.0, 0.6))
	hint.add_theme_font_size_override("font_size", 13)
	hint.position = Vector2(520, 18)
	_banner.add_child(hint)

func _on_lucky_awakened_croc(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"croc_awaken":
			_hunt_count = 0
			_show_banner()
			_flash_screen(THEME_COLOR, 2)
			ScreenShake.add_trauma(0.4)
		"croc_hunt":
			_hunt_count = data.get("hunt_count", _hunt_count)
			if _hunt_label:
				_hunt_label.text = "獵魚：%d 條" % _hunt_count
			var reward = data.get("reward", 0)
			if reward > 0:
				_show_float_text("🐊 +%d" % reward, THEME_COLOR)
		"croc_end":
			_hide_banner()
		"croc_perfect":
			_show_perfect(data.get("hunt_count", 0), data.get("boost_mult", 3.5), data.get("boost_secs", 9))
		"croc_perfect_end":
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

func _show_perfect(hunts: int, mult: float, secs: int) -> void:
	_flash_screen(GOLD_COLOR, 3)
	ScreenShake.add_trauma(0.6)
	var popup = Label.new()
	popup.text = "🐊✨ 完美覺醒！獵魚 %d 條！全服 ×%.1f 加成 %d 秒！" % [hunts, mult, secs]
	popup.add_theme_color_override("font_color", GOLD_COLOR)
	popup.add_theme_font_size_override("font_size", 22)
	popup.position = Vector2(200, 300)
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
