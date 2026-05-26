## LuckyTsunamiPanel.gd — T135 幸運海嘯魚 UI
## lucky-panel-agent 負責維護
## 業界依據：Fishing Fortune 2026「multiplier cascade system」
## 視覺主題：深藍色 + 三波衝擊指示器 + 完美海嘯演出
extends CanvasLayer

const LAYER_Z = 35

const COLOR_WAVE1    = Color(0.0, 0.5, 1.0)    # 藍色（第一波）
const COLOR_WAVE2    = Color(0.0, 0.7, 1.0)    # 亮藍（第二波）
const COLOR_WAVE3    = Color(0.0, 0.9, 1.0)    # 青藍（第三波）
const COLOR_PERFECT  = Color(0.0, 0.85, 1.0)   # 海藍（完美海嘯）
const COLOR_BG       = Color(0.0, 0.03, 0.1, 0.90)

var _banner: Control = null
var _indicator: Control = null
var _wave_labels: Array = []
var _flash_overlay: ColorRect = null

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_indicator()
	GameManager.lucky_tsunami.connect(_on_lucky_tsunami)

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(0.0, 0.5, 1.0, 0.0)
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
	title.text = "🌊 海嘯預警"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_WAVE3
	_indicator.add_child(title)

	var wave_colors = [COLOR_WAVE1, COLOR_WAVE2, COLOR_WAVE3]
	var wave_texts = ["第一波 HP -20%", "第二波 HP -30%", "第三波 HP -40%"]
	for i in range(3):
		var lbl = Label.new()
		lbl.text = wave_texts[i]
		lbl.position = Vector2(8, 28 + i * 26)
		lbl.add_theme_font_size_override("font_size", 14)
		lbl.modulate = wave_colors[i]
		_indicator.add_child(lbl)
		_wave_labels.append(lbl)

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

func _on_lucky_tsunami(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"tsunami_warning":
			_indicator.visible = true
			# 重置波次顯示
			var wave_colors = [COLOR_WAVE1, COLOR_WAVE2, COLOR_WAVE3]
			var wave_texts = ["第一波 HP -20%", "第二波 HP -30%", "第三波 HP -40%"]
			for i in range(min(_wave_labels.size(), 3)):
				if is_instance_valid(_wave_labels[i]):
					_wave_labels[i].text = wave_texts[i]
					_wave_labels[i].modulate = wave_colors[i]
			_flash(COLOR_WAVE1, 0.4, 0.2)
			_show_banner("🌊 海嘯預警！三波衝擊即將來臨！", COLOR_WAVE3)

		"wave_hit":
			var wave_num = data.get("wave_num", 1)
			var hit_count = data.get("hit_count", 0)
			var dmg_pct = data.get("damage_pct", 0.2)
			var wave_colors = [COLOR_WAVE1, COLOR_WAVE2, COLOR_WAVE3]
			var color = wave_colors[min(wave_num - 1, 2)]
			_flash(color, 0.5 + (wave_num - 1) * 0.1, 0.25)
			if wave_num - 1 < _wave_labels.size() and is_instance_valid(_wave_labels[wave_num - 1]):
				_wave_labels[wave_num - 1].text = "第%d波 ✓ 命中 %d 條" % [wave_num, hit_count]
				_wave_labels[wave_num - 1].modulate = Color(0.2, 0.9, 0.2)
			_show_banner("🌊 第 %d 波！命中 %d 條魚！HP -%d%%！" % [wave_num, hit_count, int(dmg_pct * 100)], color)

		"tsunami_perfect":
			_indicator.visible = false
			_flash(COLOR_PERFECT, 0.7, 0.3)
			_flash(COLOR_PERFECT, 0.7, 0.3)
			_flash(COLOR_PERFECT, 0.7, 0.3)
			var total = data.get("total_hit_count", 0)
			_show_banner("🌊 完美海嘯！三波命中 %d 條魚！全服 ×3.2 加成 8 秒！" % total, COLOR_PERFECT)

		"tsunami_perfect_end":
			pass

		"tsunami_end":
			_indicator.visible = false
