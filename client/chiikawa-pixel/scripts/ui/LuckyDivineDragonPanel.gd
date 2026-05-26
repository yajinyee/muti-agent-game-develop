## LuckyDivineDragonPanel.gd — T145 幸運神龍魚 Lucky Panel
## lucky-panel-agent 負責維護
## 業界依據：Royal Fishing Jili「Divine Dragon descends from heavens」
## 神龍主題：金色 + 爪擊計數器 + 完美指示器
extends CanvasLayer

const THEME_COLOR = Color(1.0, 0.85, 0.0)        # 金色
const DARK_COLOR  = Color(0.35, 0.28, 0.0)        # 深金
const RED_COLOR   = Color(0.9, 0.1, 0.1)          # 紅色

var _banner: Control = null
var _claw_dots: Array = []
var _claw_label: Label = null
var _perfect_claws: int = 0

func _ready() -> void:
	layer = 40
	_create_ui()
	GameManager.lucky_divine_dragon.connect(_on_lucky_divine_dragon)

func _create_ui() -> void:
	_banner = Control.new()
	_banner.visible = false
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 64)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.color = Color(DARK_COLOR.r, DARK_COLOR.g, DARK_COLOR.b, 0.92)
	bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	_banner.add_child(bg)

	var top_line = ColorRect.new()
	top_line.color = THEME_COLOR
	top_line.size = Vector2(1280, 4)
	_banner.add_child(top_line)

	var title = Label.new()
	title.text = "🐉 神龍降臨！"
	title.add_theme_color_override("font_color", THEME_COLOR)
	title.add_theme_font_size_override("font_size", 22)
	title.position = Vector2(20, 14)
	_banner.add_child(title)

	# 5爪擊指示點
	for i in range(5):
		var dot = ColorRect.new()
		dot.size = Vector2(22, 22)
		dot.position = Vector2(220 + i * 30, 18)
		dot.color = Color(0.4, 0.4, 0.4)
		_banner.add_child(dot)
		_claw_dots.append(dot)

	_claw_label = Label.new()
	_claw_label.text = "完美爪擊：0/5"
	_claw_label.add_theme_color_override("font_color", THEME_COLOR)
	_claw_label.add_theme_font_size_override("font_size", 16)
	_claw_label.position = Vector2(380, 18)
	_banner.add_child(_claw_label)

func _on_lucky_divine_dragon(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"dragon_descend":
			_perfect_claws = 0
			_show_banner()
			_flash_screen(THEME_COLOR, 3)
		"dragon_claw":
			var claw_num = data.get("claw_num", 1) - 1
			_perfect_claws = data.get("perfect_claws", _perfect_claws)
			if claw_num < _claw_dots.size():
				var hit = data.get("hit_count", 0)
				_claw_dots[claw_num].color = THEME_COLOR if hit >= 5 else RED_COLOR
			if _claw_label:
				_claw_label.text = "完美爪擊：%d/5" % _perfect_claws
		"dragon_leave":
			_hide_banner()
		"dragon_perfect":
			_show_perfect(data.get("boost_mult", 5.0), data.get("boost_sec", 12))
		"dragon_perfect_end":
			_hide_banner()

func _show_banner() -> void:
	if _banner:
		_banner.visible = true
		for dot in _claw_dots:
			dot.color = Color(0.4, 0.4, 0.4)
		var tween = create_tween()
		_banner.modulate.a = 0.0
		tween.tween_property(_banner, "modulate:a", 1.0, 0.3)

func _hide_banner() -> void:
	if _banner:
		var tween = create_tween()
		tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
		tween.tween_callback(func(): _banner.visible = false)

func _show_perfect(mult: float, secs: int) -> void:
	_flash_screen(THEME_COLOR, 4)
	var popup = Label.new()
	popup.text = "🐉 神龍完美！全服 ×%.1f 加成 %d 秒！" % [mult, secs]
	popup.add_theme_color_override("font_color", THEME_COLOR)
	popup.add_theme_font_size_override("font_size", 24)
	popup.position = Vector2(180, 300)
	add_child(popup)
	var tween = create_tween()
	tween.tween_property(popup, "position:y", 250.0, 0.5)
	tween.tween_interval(2.5)
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
