## LuckyHumpbackWhalePanel.gd — T137 幸運座頭鯨魚 UI
## lucky-panel-agent 負責維護
## 業界依據：Royal Fishing Jili「Humpback Whale 90-150x with 15x base multiplier」
## 視覺主題：深藍色 + 四波聲波指示器 + 完美鯨歌演出
extends CanvasLayer

const LAYER_Z = 37

const COLOR_WHALE    = Color(0.0, 0.4, 0.8)    # 深藍（鯨魚主色）
const COLOR_SONG     = Color(0.0, 0.7, 1.0)    # 亮藍（鯨歌）
const COLOR_PERFECT  = Color(0.0, 0.9, 1.0)    # 青藍（完美鯨歌）
const COLOR_BG       = Color(0.0, 0.02, 0.08, 0.90)

var _banner: Control = null
var _indicator: Control = null
var _wave_labels: Array = []
var _flash_overlay: ColorRect = null

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_indicator()
	GameManager.lucky_humpback_whale.connect(_on_lucky_humpback_whale)

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(0.0, 0.4, 0.8, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

func _create_indicator() -> void:
	_indicator = Control.new()
	_indicator.position = Vector2(960, 10)
	_indicator.size = Vector2(300, 130)
	_indicator.visible = false
	add_child(_indicator)

	var bg = ColorRect.new()
	bg.size = _indicator.size
	bg.color = COLOR_BG
	_indicator.add_child(bg)

	var title = Label.new()
	title.text = "🐋 鯨歌共鳴"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_SONG
	_indicator.add_child(title)

	for i in range(4):
		var lbl = Label.new()
		lbl.text = "第%d波 HP -15%%" % (i + 1)
		lbl.position = Vector2(8, 28 + i * 24)
		lbl.add_theme_font_size_override("font_size", 13)
		lbl.modulate = COLOR_WHALE
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

func _on_lucky_humpback_whale(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"song_start":
			_indicator.visible = true
			_flash(COLOR_SONG, 0.4, 0.2)
			_show_banner("🐋 鯨歌共鳴！四波聲波即將衝擊！", COLOR_SONG)

		"song_wave":
			var wave_num = data.get("wave_num", 1)
			var hit_count = data.get("hit_count", 0)
			_flash(COLOR_SONG, 0.4, 0.2)
			if wave_num - 1 < _wave_labels.size() and is_instance_valid(_wave_labels[wave_num - 1]):
				_wave_labels[wave_num - 1].text = "第%d波 ✓ 命中 %d 條" % [wave_num, hit_count]
				_wave_labels[wave_num - 1].modulate = Color(0.2, 0.9, 0.2)
			_show_banner("🐋 第 %d 波！命中 %d 條魚！" % [wave_num, hit_count], COLOR_SONG)

		"song_perfect":
			_indicator.visible = false
			_flash(COLOR_PERFECT, 0.7, 0.3)
			_flash(COLOR_PERFECT, 0.7, 0.3)
			var total = data.get("total_hit_count", 0)
			_show_banner("🐋 完美鯨歌！四波命中 %d 條！全服 ×3.0 加成 8 秒！" % total, COLOR_PERFECT)

		"song_perfect_end":
			pass

		"song_end":
			_indicator.visible = false
