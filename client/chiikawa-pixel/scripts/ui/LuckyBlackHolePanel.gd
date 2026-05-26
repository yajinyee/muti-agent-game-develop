## LuckyBlackHolePanel.gd — T133 幸運黑洞魚 UI
## lucky-panel-agent 負責維護
## 業界依據：Godot vortex water shader + 黑洞吸引機制
## 視覺主題：深紫黑色 + 吸引計時條 + 坍縮爆炸演出 + 奇點爆發彈窗
extends CanvasLayer

const LAYER_Z = 33

const COLOR_VOID     = Color(0.3, 0.0, 0.5)    # 深紫（黑洞主色）
const COLOR_PULL     = Color(0.6, 0.0, 1.0)    # 紫色（吸引效果）
const COLOR_COLLAPSE = Color(1.0, 0.3, 1.0)    # 粉紫（坍縮）
const COLOR_SINGULAR = Color(1.0, 0.85, 0.0)   # 金色（奇點爆發）
const COLOR_BG       = Color(0.02, 0.0, 0.05, 0.92)

var _banner: Control = null
var _indicator: Control = null
var _pull_bar: ColorRect = null
var _pull_label: Label = null
var _flash_overlay: ColorRect = null
var _pull_timer: float = 0.0
var _pull_duration: float = 8.0

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_indicator()
	GameManager.lucky_black_hole.connect(_on_lucky_black_hole)

func _process(delta: float) -> void:
	if _pull_timer > 0.0:
		_pull_timer -= delta
		if _pull_timer < 0.0:
			_pull_timer = 0.0
		if is_instance_valid(_pull_bar):
			_pull_bar.size.x = 280.0 * (_pull_timer / _pull_duration)
		if is_instance_valid(_pull_label):
			_pull_label.text = "坍縮倒數：%.1fs" % _pull_timer

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(0.3, 0.0, 0.5, 0.0)
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
	title.text = "🌑 黑洞吸引中"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_PULL
	_indicator.add_child(title)

	var speed_label = Label.new()
	speed_label.text = "目標速度 ×0.2"
	speed_label.position = Vector2(8, 28)
	speed_label.add_theme_font_size_override("font_size", 16)
	speed_label.modulate = COLOR_COLLAPSE
	_indicator.add_child(speed_label)

	_pull_label = Label.new()
	_pull_label.text = "坍縮倒數：8.0s"
	_pull_label.position = Vector2(8, 52)
	_pull_label.add_theme_font_size_override("font_size", 14)
	_pull_label.modulate = Color(0.8, 0.8, 0.8)
	_indicator.add_child(_pull_label)

	var bar_bg = ColorRect.new()
	bar_bg.position = Vector2(8, 78)
	bar_bg.size = Vector2(280, 12)
	bar_bg.color = Color(0.2, 0.2, 0.2)
	_indicator.add_child(bar_bg)

	_pull_bar = ColorRect.new()
	_pull_bar.position = Vector2(8, 78)
	_pull_bar.size = Vector2(280, 12)
	_pull_bar.color = COLOR_PULL
	_indicator.add_child(_pull_bar)

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

func _on_lucky_black_hole(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"black_hole_start":
			_pull_timer = data.get("duration", 8.0)
			_pull_duration = _pull_timer
			_indicator.visible = true
			_flash(COLOR_PULL, 0.5, 0.3)
			_show_banner("🌑 黑洞降臨！目標速度 ×0.2！8 秒後坍縮！", COLOR_PULL)

		"collapse":
			_indicator.visible = false
			_pull_timer = 0.0
			var hit_count = data.get("hit_count", 0)
			_flash(COLOR_COLLAPSE, 0.7, 0.3)
			_flash(COLOR_COLLAPSE, 0.7, 0.3)
			_show_banner("💥 黑洞坍縮！命中 %d 條魚！HP -50%%！" % hit_count, COLOR_COLLAPSE)

		"singularity":
			_flash(COLOR_SINGULAR, 0.8, 0.35)
			_flash(COLOR_SINGULAR, 0.8, 0.35)
			_flash(COLOR_SINGULAR, 0.8, 0.35)
			var hit_count = data.get("hit_count", 0)
			_show_banner("🌑 奇點爆發！吸入 %d 條魚！全服 ×3.0 加成 8 秒！" % hit_count, COLOR_SINGULAR)

		"singularity_end":
			pass

		"black_hole_end":
			_indicator.visible = false
