## LuckyRebirthPanel.gd — T150 幸運重生魚 Lucky Panel
## lucky-panel-agent 負責維護
## 業界依據：Phoenix rebirth mechanic — 死亡目標復活再擊破，雙重獎勵
## 重生主題：火橙色 + 鳳凰羽毛特效 + 重生計數器
extends CanvasLayer

const THEME_COLOR = Color(1.0, 0.4, 0.1)         # 火橙色
const DARK_COLOR  = Color(0.25, 0.08, 0.0)        # 深棕
const GOLD_COLOR  = Color(1.0, 0.85, 0.0)         # 金色

var _banner: Control = null
var _rebirth_label: Label = null
var _rebirth_count: int = 0

func _ready() -> void:
	layer = 45
	_create_ui()
	GameManager.lucky_rebirth.connect(_on_lucky_rebirth)

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
	title.text = "🔥 重生之力！"
	title.add_theme_color_override("font_color", THEME_COLOR)
	title.add_theme_font_size_override("font_size", 20)
	title.position = Vector2(20, 12)
	_banner.add_child(title)

	_rebirth_label = Label.new()
	_rebirth_label.text = "重生擊破：0 個"
	_rebirth_label.add_theme_color_override("font_color", GOLD_COLOR)
	_rebirth_label.add_theme_font_size_override("font_size", 16)
	_rebirth_label.position = Vector2(280, 16)
	_banner.add_child(_rebirth_label)

	var hint = Label.new()
	hint.text = "復活目標獎勵 ×3.0"
	hint.add_theme_color_override("font_color", Color(1.0, 0.7, 0.4))
	hint.add_theme_font_size_override("font_size", 13)
	hint.position = Vector2(520, 18)
	_banner.add_child(hint)

func _on_lucky_rebirth(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"rebirth_start":
			_rebirth_count = 0
			_show_banner()
			_flash_screen(THEME_COLOR, 2)
		"rebirth_kill":
			_rebirth_count = data.get("rebirth_kills", _rebirth_count)
			if _rebirth_label:
				_rebirth_label.text = "重生擊破：%d 個" % _rebirth_count
		"rebirth_end":
			_hide_banner()
		"rebirth_perfect":
			_show_perfect(data.get("rebirth_kills", 0), data.get("boost_mult", 6.5), data.get("boost_sec", 15))
		"rebirth_perfect_end":
			_hide_banner()

func _show_banner() -> void:
	if _banner:
		_banner.visible = true
		var tween = create_tween()
		_banner.modulate.a = 0.0
		tween.tween_property(_banner, "modulate:a", 1.0, 0.3)

func _hide_banner() -> void:
	if _banner:
		var tween = create_tween()
		tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
		tween.tween_callback(func(): _banner.visible = false)

func _show_perfect(kills: int, mult: float, secs: int) -> void:
	_flash_screen(GOLD_COLOR, 3)
	var popup = Label.new()
	popup.text = "🔥 完美重生！重生擊破 %d 個！全服 ×%.1f 加成 %d 秒！" % [kills, mult, secs]
	popup.add_theme_color_override("font_color", GOLD_COLOR)
	popup.add_theme_font_size_override("font_size", 22)
	popup.position = Vector2(200, 300)
	add_child(popup)
	var tween = create_tween()
	tween.tween_property(popup, "position:y", 260.0, 0.5)
	tween.tween_interval(2.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)

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
