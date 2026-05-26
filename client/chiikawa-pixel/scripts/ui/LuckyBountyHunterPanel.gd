## LuckyBountyHunterPanel.gd — T134 幸運賞金獵人魚 UI
## lucky-panel-agent 負責維護
## 業界依據：Fishing Frenzy Chapter 3「Guild Wars + Boss Fish + quality values」
## 視覺主題：火橙色 + 賞金目標計數器 + 完美賞金演出
extends CanvasLayer

const LAYER_Z = 34

const COLOR_BOUNTY   = Color(1.0, 0.5, 0.0)    # 火橙（賞金主色）
const COLOR_KILL     = Color(1.0, 0.85, 0.0)   # 金色（擊破賞金）
const COLOR_PERFECT  = Color(1.0, 0.3, 0.1)    # 火紅（完美賞金）
const COLOR_BG       = Color(0.1, 0.04, 0.0, 0.90)

var _banner: Control = null
var _indicator: Control = null
var _kill_label: Label = null
var _time_label: Label = null
var _flash_overlay: ColorRect = null
var _bounty_timer: float = 0.0

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_indicator()
	GameManager.lucky_bounty_hunter.connect(_on_lucky_bounty_hunter)

func _process(delta: float) -> void:
	if _bounty_timer > 0.0:
		_bounty_timer -= delta
		if _bounty_timer < 0.0:
			_bounty_timer = 0.0
		if is_instance_valid(_time_label):
			_time_label.text = "剩餘：%.0fs" % _bounty_timer

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.5, 0.0, 0.0)
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
	title.text = "🎯 賞金獵人任務"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_BOUNTY
	_indicator.add_child(title)

	_kill_label = Label.new()
	_kill_label.text = "賞金目標：0/3"
	_kill_label.position = Vector2(8, 28)
	_kill_label.add_theme_font_size_override("font_size", 18)
	_kill_label.modulate = COLOR_KILL
	_indicator.add_child(_kill_label)

	_time_label = Label.new()
	_time_label.text = "剩餘：30s"
	_time_label.position = Vector2(8, 56)
	_time_label.add_theme_font_size_override("font_size", 14)
	_time_label.modulate = Color(0.8, 0.8, 0.8)
	_indicator.add_child(_time_label)

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

func _on_lucky_bounty_hunter(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"bounty_start":
			_bounty_timer = data.get("duration", 30.0)
			_indicator.visible = true
			if is_instance_valid(_kill_label):
				_kill_label.text = "賞金目標：0/%d" % data.get("total_bounty", 3)
			_flash(COLOR_BOUNTY, 0.4, 0.2)
			_show_banner("🎯 賞金任務！標記 %d 個賞金目標！30 秒內全部獵殺！" % data.get("total_bounty", 3), COLOR_BOUNTY)

		"bounty_kill":
			var kill_count = data.get("kill_count", 0)
			var total = data.get("total_bounty", 3)
			if is_instance_valid(_kill_label):
				_kill_label.text = "賞金目標：%d/%d" % [kill_count, total]
			_flash(COLOR_KILL, 0.35, 0.15)

		"bounty_perfect":
			_indicator.visible = false
			_bounty_timer = 0.0
			_flash(COLOR_PERFECT, 0.7, 0.3)
			_flash(COLOR_PERFECT, 0.7, 0.3)
			_flash(COLOR_PERFECT, 0.7, 0.3)
			var kill_count = data.get("kill_count", 0)
			_show_banner("🎯 完美賞金！獵殺全部 %d 個目標！全服 ×3.5 加成 8 秒！" % kill_count, COLOR_PERFECT)

		"bounty_perfect_end":
			pass

		"bounty_timeout":
			_indicator.visible = false
			_bounty_timer = 0.0
