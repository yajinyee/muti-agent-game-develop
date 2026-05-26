## LuckyAnglerFishPanel.gd — T132 幸運巨型安康魚 UI
## lucky-panel-agent 負責維護
## 業界依據：Jili「Giant Anglerfish shoots electricity to open treasure chests」
## 視覺主題：深海青藍色 + 誘餌計時條 + 電擊爆炸演出 + 完美誘捕彈窗
extends CanvasLayer

const LAYER_Z = 32

const COLOR_LURE     = Color(0.0, 0.9, 1.0)    # 青藍（誘餌主色）
const COLOR_EXPLODE  = Color(1.0, 0.9, 0.0)    # 金黃（電擊爆炸）
const COLOR_PERFECT  = Color(0.0, 1.0, 0.5)    # 翠綠（完美誘捕）
const COLOR_BG       = Color(0.0, 0.05, 0.1, 0.90)

var _banner: Control = null
var _indicator: Control = null
var _lure_bar: ColorRect = null
var _lure_label: Label = null
var _flash_overlay: ColorRect = null
var _lure_timer: float = 0.0
var _lure_duration: float = 5.0

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_indicator()
	GameManager.lucky_angler_fish.connect(_on_lucky_angler_fish)

func _process(delta: float) -> void:
	if _lure_timer > 0.0:
		_lure_timer -= delta
		if _lure_timer < 0.0:
			_lure_timer = 0.0
		if is_instance_valid(_lure_bar):
			_lure_bar.size.x = 280.0 * (_lure_timer / _lure_duration)
		if is_instance_valid(_lure_label):
			_lure_label.text = "誘餌倒數：%.1fs" % _lure_timer

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(0.0, 0.9, 1.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

func _create_indicator() -> void:
	_indicator = Control.new()
	_indicator.position = Vector2(960, 10)
	_indicator.size = Vector2(300, 110)
	_indicator.visible = false
	add_child(_indicator)

	var bg = ColorRect.new()
	bg.size = _indicator.size
	bg.color = COLOR_BG
	_indicator.add_child(bg)

	var title = Label.new()
	title.text = "🎣 安康魚誘餌中"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_LURE
	_indicator.add_child(title)

	var dmg_label = Label.new()
	dmg_label.text = "傷害 ×1.8"
	dmg_label.position = Vector2(8, 28)
	dmg_label.add_theme_font_size_override("font_size", 16)
	dmg_label.modulate = COLOR_EXPLODE
	_indicator.add_child(dmg_label)

	_lure_label = Label.new()
	_lure_label.text = "誘餌倒數：5.0s"
	_lure_label.position = Vector2(8, 52)
	_lure_label.add_theme_font_size_override("font_size", 14)
	_lure_label.modulate = Color(0.8, 0.8, 0.8)
	_indicator.add_child(_lure_label)

	# 誘餌計時條
	var bar_bg = ColorRect.new()
	bar_bg.position = Vector2(8, 78)
	bar_bg.size = Vector2(280, 12)
	bar_bg.color = Color(0.2, 0.2, 0.2)
	_indicator.add_child(bar_bg)

	_lure_bar = ColorRect.new()
	_lure_bar.position = Vector2(8, 78)
	_lure_bar.size = Vector2(280, 12)
	_lure_bar.color = COLOR_LURE
	_indicator.add_child(_lure_bar)

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

func _on_lucky_angler_fish(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"lure_start":
			_lure_timer = data.get("lure_sec", 5.0)
			_lure_duration = _lure_timer
			_indicator.visible = true
			_flash(COLOR_LURE, 0.4, 0.2)
			_show_banner("🎣 深海誘餌！傷害 ×1.8！5 秒後電擊爆炸！", COLOR_LURE)

		"explosion":
			_indicator.visible = false
			_lure_timer = 0.0
			var hit_count = data.get("hit_count", 0)
			_flash(COLOR_EXPLODE, 0.6, 0.25)
			_show_banner("⚡ 電擊爆炸！命中 %d 條魚！HP -30%%！" % hit_count, COLOR_EXPLODE)

		"perfect":
			_flash(COLOR_PERFECT, 0.7, 0.3)
			_flash(COLOR_PERFECT, 0.7, 0.3)
			var hit_count = data.get("hit_count", 0)
			_show_banner("🎣 完美誘捕！命中 %d 條！全服 ×2.8 加成 7 秒！" % hit_count, COLOR_PERFECT)

		"perfect_end":
			pass

		"lure_end":
			_indicator.visible = false
