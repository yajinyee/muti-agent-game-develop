## LuckyLegendDragonPanel.gd — T138 幸運傳說龍魚 UI
## lucky-panel-agent 負責維護
## 業界依據：Royal Fishing Jili「Legend Dragon 120-200x from 20x base multiplier」
## 視覺主題：金橙色 + 噴火計數器 + 傳說龍怒演出
extends CanvasLayer

const LAYER_Z = 38

const COLOR_DRAGON   = Color(1.0, 0.5, 0.0)    # 金橙（龍主色）
const COLOR_FIRE     = Color(1.0, 0.2, 0.0)    # 火紅（噴火）
const COLOR_RAGE     = Color(1.0, 0.85, 0.0)   # 金色（龍怒）
const COLOR_BG       = Color(0.08, 0.03, 0.0, 0.90)

var _banner: Control = null
var _indicator: Control = null
var _breath_label: Label = null
var _flash_overlay: ColorRect = null

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_indicator()
	GameManager.lucky_legend_dragon.connect(_on_lucky_legend_dragon)

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.5, 0.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

func _create_indicator() -> void:
	_indicator = Control.new()
	_indicator.position = Vector2(960, 10)
	_indicator.size = Vector2(300, 90)
	_indicator.visible = false
	add_child(_indicator)

	var bg = ColorRect.new()
	bg.size = _indicator.size
	bg.color = COLOR_BG
	_indicator.add_child(bg)

	var title = Label.new()
	title.text = "🐲 傳說龍降臨"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_DRAGON
	_indicator.add_child(title)

	_breath_label = Label.new()
	_breath_label.text = "噴火：0/4"
	_breath_label.position = Vector2(8, 28)
	_breath_label.add_theme_font_size_override("font_size", 20)
	_breath_label.modulate = COLOR_FIRE
	_indicator.add_child(_breath_label)

	var hint = Label.new()
	hint.text = "每 3 秒噴火一次"
	hint.position = Vector2(8, 58)
	hint.add_theme_font_size_override("font_size", 12)
	hint.modulate = Color(0.7, 0.7, 0.7)
	_indicator.add_child(hint)

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

func _on_lucky_legend_dragon(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"dragon_appear":
			_indicator.visible = true
			if is_instance_valid(_breath_label):
				_breath_label.text = "噴火：0/4"
			_flash(COLOR_DRAGON, 0.5, 0.3)
			_show_banner("🐲 傳說龍降臨！15 秒內 4 次噴火！", COLOR_DRAGON)

		"dragon_breath":
			var breath_num = data.get("breath_num", 1)
			var hit_count = data.get("hit_count", 0)
			var perfect = data.get("perfect_breaths", 0)
			if is_instance_valid(_breath_label):
				_breath_label.text = "噴火：%d/4（完美 %d）" % [breath_num, perfect]
			_flash(COLOR_FIRE, 0.5, 0.2)
			_show_banner("🔥 第 %d 次噴火！命中 %d 條！HP -35%%！" % [breath_num, hit_count], COLOR_FIRE)

		"dragon_rage":
			_indicator.visible = false
			_flash(COLOR_RAGE, 0.8, 0.35)
			_flash(COLOR_RAGE, 0.8, 0.35)
			_flash(COLOR_RAGE, 0.8, 0.35)
			_show_banner("🐲 傳說龍怒！完美噴火 4 次！全服 ×4.0 加成 10 秒！", COLOR_RAGE)

		"dragon_rage_end":
			pass

		"dragon_leave":
			_indicator.visible = false
