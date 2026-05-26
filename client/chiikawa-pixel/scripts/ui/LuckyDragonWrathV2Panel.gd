## LuckyDragonWrathV2Panel.gd — T136 幸運龍怒蓄積魚 v2 UI
## lucky-panel-agent 負責維護
## 業界依據：Royal Fishing Jili「Dragon Wrath system accumulates with every shot fired」升級版
## 視覺主題：火橙深紅色 + 怒氣條 + 隕石爆發演出 + 完美龍怒彈窗
extends CanvasLayer

const LAYER_Z = 36

const COLOR_WRATH    = Color(1.0, 0.27, 0.0)   # 火橙（怒氣主色）
const COLOR_METEOR   = Color(1.0, 0.6, 0.0)    # 橙色（隕石）
const COLOR_PERFECT  = Color(1.0, 0.85, 0.0)   # 金色（完美龍怒）
const COLOR_BG       = Color(0.1, 0.02, 0.0, 0.90)

var _banner: Control = null
var _indicator: Control = null
var _wrath_bar: ColorRect = null
var _wrath_label: Label = null
var _time_label: Label = null
var _flash_overlay: ColorRect = null
var _wrath_timer: float = 0.0
var _wrath_value: int = 0

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_indicator()
	GameManager.lucky_dragon_wrath_v2.connect(_on_lucky_dragon_wrath_v2)

func _process(delta: float) -> void:
	if _wrath_timer > 0.0:
		_wrath_timer -= delta
		if _wrath_timer < 0.0:
			_wrath_timer = 0.0
		if is_instance_valid(_time_label):
			_time_label.text = "剩餘：%.0fs" % _wrath_timer

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.27, 0.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

func _create_indicator() -> void:
	_indicator = Control.new()
	_indicator.position = Vector2(960, 10)
	_indicator.size = Vector2(300, 120)
	_indicator.visible = false
	add_child(_indicator)

	var bg = ColorRect.new()
	bg.size = _indicator.size
	bg.color = COLOR_BG
	_indicator.add_child(bg)

	var title = Label.new()
	title.text = "🐉 龍怒蓄積中"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_WRATH
	_indicator.add_child(title)

	_wrath_label = Label.new()
	_wrath_label.text = "怒氣：0/30"
	_wrath_label.position = Vector2(8, 28)
	_wrath_label.add_theme_font_size_override("font_size", 18)
	_wrath_label.modulate = COLOR_METEOR
	_indicator.add_child(_wrath_label)

	_time_label = Label.new()
	_time_label.text = "剩餘：30s"
	_time_label.position = Vector2(8, 56)
	_time_label.add_theme_font_size_override("font_size", 14)
	_time_label.modulate = Color(0.8, 0.8, 0.8)
	_indicator.add_child(_time_label)

	var bar_bg = ColorRect.new()
	bar_bg.position = Vector2(8, 82)
	bar_bg.size = Vector2(280, 14)
	bar_bg.color = Color(0.2, 0.2, 0.2)
	_indicator.add_child(bar_bg)

	_wrath_bar = ColorRect.new()
	_wrath_bar.position = Vector2(8, 82)
	_wrath_bar.size = Vector2(0, 14)
	_wrath_bar.color = COLOR_WRATH
	_indicator.add_child(_wrath_bar)

func _flash(color: Color, alpha: float = 0.4, duration: float = 0.2) -> void:
	if not is_instance_valid(_flash_overlay):
		return
	_flash_overlay.color = Color(color.r, color.g, color.b, alpha)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)

func _show_banner(text: String, color: Color) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()
	_banner = Control.new()
	_banner.position = Vector2(0, 100)
	_banner.size = Vector2(1280, 70)
	add_child(_banner)
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0, 0, 0, 0.85)
	_banner.add_child(bg)
	var lbl = Label.new()
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 26)
	lbl.modulate = color
	_banner.add_child(lbl)
	var tween = _banner.create_tween()
	tween.tween_interval(2.5)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_banner):
			_banner.queue_free()
			_banner = null
	)

func _on_lucky_dragon_wrath_v2(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"wrath_start":
			_wrath_timer = data.get("duration", 30.0)
			_wrath_value = 0
			_indicator.visible = true
			if is_instance_valid(_wrath_label):
				_wrath_label.text = "怒氣：0/%d" % data.get("max_wrath", 30)
			if is_instance_valid(_wrath_bar):
				_wrath_bar.size.x = 0
			_flash(COLOR_WRATH, 0.4, 0.2)
			_show_banner("🐉 龍怒蓄積！30 秒內射擊蓄積怒氣！", COLOR_WRATH)

		"wrath_explode":
			_indicator.visible = false
			var wrath = data.get("wrath_value", 0)
			var meteors = data.get("meteor_count", 0)
			_flash(COLOR_METEOR, 0.6, 0.3)
			_flash(COLOR_METEOR, 0.6, 0.3)
			_show_banner("🐉 龍怒爆發！怒氣 %d！%d 顆隕石！" % [wrath, meteors], COLOR_METEOR)

		"wrath_perfect":
			_flash(COLOR_PERFECT, 0.7, 0.3)
			_flash(COLOR_PERFECT, 0.7, 0.3)
			_flash(COLOR_PERFECT, 0.7, 0.3)
			var wrath = data.get("wrath_value", 0)
			_show_banner("🐉 完美龍怒！怒氣 %d！全服 ×3.5 加成 8 秒！" % wrath, COLOR_PERFECT)

		"wrath_perfect_end":
			pass

		"wrath_end":
			_indicator.visible = false
