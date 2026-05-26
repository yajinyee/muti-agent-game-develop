## LuckySupernovaPanel.gd — T147 幸運超新星魚 Lucky Panel
## lucky-panel-agent 負責維護
## 業界依據：Supernova explosion mechanic — 全場爆炸+倍率加成
## 超新星主題：火橙色 + 爆炸特效 + 命中計數器
extends CanvasLayer

const THEME_COLOR = Color(1.0, 0.4, 0.1)        # 火橙色
const DARK_COLOR  = Color(0.3, 0.1, 0.0)         # 深棕
const GOLD_COLOR  = Color(1.0, 0.85, 0.0)        # 金色

var _banner: Control = null
var _hit_label: Label = null
var _mult_label: Label = null

func _ready() -> void:
	layer = 42
	_create_ui()
	GameManager.lucky_supernova.connect(_on_lucky_supernova)

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
	title.text = "💥 超新星爆炸！"
	title.add_theme_color_override("font_color", THEME_COLOR)
	title.add_theme_font_size_override("font_size", 20)
	title.position = Vector2(20, 12)
	_banner.add_child(title)

	_hit_label = Label.new()
	_hit_label.text = "命中：0 個"
	_hit_label.add_theme_color_override("font_color", GOLD_COLOR)
	_hit_label.add_theme_font_size_override("font_size", 16)
	_hit_label.position = Vector2(300, 16)
	_banner.add_child(_hit_label)

	_mult_label = Label.new()
	_mult_label.text = ""
	_mult_label.add_theme_color_override("font_color", Color(1.0, 0.6, 0.2))
	_mult_label.add_theme_font_size_override("font_size", 14)
	_mult_label.position = Vector2(500, 18)
	_banner.add_child(_mult_label)

func _on_lucky_supernova(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"supernova_explode":
			_show_banner()
			_flash_screen(THEME_COLOR, 3)
			ScreenShake.add_trauma(1.0)
		"supernova_boost":
			var hits = data.get("hit_count", 0)
			if _hit_label:
				_hit_label.text = "命中：%d 個" % hits
			if _mult_label:
				_mult_label.text = "5秒 ×3.0 加成！"
		"supernova_end":
			_hide_banner()
		"supernova_perfect":
			_show_perfect(data.get("boost_mult", 5.5), data.get("boost_sec", 12))
		"supernova_perfect_end":
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

func _show_perfect(mult: float, secs: int) -> void:
	_flash_screen(GOLD_COLOR, 3)
	var popup = Label.new()
	popup.text = "💥 超新星完美！全服 ×%.1f 加成 %d 秒！" % [mult, secs]
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
		tween.tween_property(overlay, "color:a", 0.4, 0.08)
		tween.tween_property(overlay, "color:a", 0.0, 0.12)
	tween.tween_callback(overlay.queue_free)
