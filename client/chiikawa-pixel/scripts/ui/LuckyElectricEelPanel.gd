## LuckyElectricEelPanel.gd — T131 幸運電鰻魚 UI
## lucky-panel-agent 負責維護
## 業界依據：Royal Fishing Jili「60x lightning eel chain reactions」升級版
## 視覺主題：電黃色 + 閃電計數器 + 連鎖加速指示器 + 超級放電演出
extends CanvasLayer

const LAYER_Z = 31

const COLOR_ELECTRIC = Color(1.0, 0.95, 0.0)   # 電黃（主色）
const COLOR_CHAIN    = Color(0.0, 0.9, 1.0)    # 青藍（連鎖加速）
const COLOR_SUPER    = Color(1.0, 0.85, 0.0)   # 金色（超級放電）
const COLOR_BG       = Color(0.05, 0.05, 0.0, 0.90)

var _banner: Control = null
var _indicator: Control = null
var _shock_label: Label = null
var _time_label: Label = null
var _flash_overlay: ColorRect = null

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_indicator()
	GameManager.lucky_electric_eel.connect(_on_lucky_electric_eel)

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.95, 0.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

func _create_indicator() -> void:
	_indicator = Control.new()
	_indicator.position = Vector2(960, 10)
	_indicator.size = Vector2(300, 100)
	_indicator.visible = false
	add_child(_indicator)

	var bg = ColorRect.new()
	bg.size = _indicator.size
	bg.color = COLOR_BG
	_indicator.add_child(bg)

	var title = Label.new()
	title.text = "⚡ 電鰻放電中"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_ELECTRIC
	_indicator.add_child(title)

	_shock_label = Label.new()
	_shock_label.text = "電擊次數：0"
	_shock_label.position = Vector2(8, 28)
	_shock_label.add_theme_font_size_override("font_size", 18)
	_shock_label.modulate = COLOR_ELECTRIC
	_indicator.add_child(_shock_label)

	_time_label = Label.new()
	_time_label.text = "剩餘：12.0s"
	_time_label.position = Vector2(8, 56)
	_time_label.add_theme_font_size_override("font_size", 14)
	_time_label.modulate = Color(0.8, 0.8, 0.8)
	_indicator.add_child(_time_label)

func _flash(color: Color, alpha: float = 0.35, duration: float = 0.15) -> void:
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
	tween.tween_interval(2.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_banner):
			_banner.queue_free()
			_banner = null
	)

func _on_lucky_electric_eel(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"eel_start":
			_indicator.visible = true
			if is_instance_valid(_shock_label):
				_shock_label.text = "電擊次數：0"
			if is_instance_valid(_time_label):
				_time_label.text = "剩餘：%.1fs" % data.get("duration", 12.0)
			_flash(COLOR_ELECTRIC, 0.4, 0.2)
			_show_banner("⚡ 電鰻放電！持續 12 秒連鎖電擊！", COLOR_ELECTRIC)

		"eel_shock":
			var shock_count = data.get("shock_count", 0)
			var hit_count = data.get("hit_count", 0)
			var time_left = data.get("time_left", 0.0)
			if is_instance_valid(_shock_label):
				_shock_label.text = "電擊次數：%d（命中 %d）" % [shock_count, hit_count]
			if is_instance_valid(_time_label):
				_time_label.text = "剩餘：%.1fs" % time_left
			_flash(COLOR_ELECTRIC, 0.25, 0.1)

		"eel_end":
			_indicator.visible = false

		"eel_super":
			_indicator.visible = false
			_flash(COLOR_SUPER, 0.7, 0.3)
			_flash(COLOR_SUPER, 0.7, 0.3)
			_flash(COLOR_SUPER, 0.7, 0.3)
			var shock_count = data.get("shock_count", 0)
			_show_banner("⚡ 超級放電！電擊 %d 次！全服 ×2.5 加成 7 秒！" % shock_count, COLOR_SUPER)

		"eel_super_end":
			pass
