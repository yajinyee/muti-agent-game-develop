## LuckyGiantPrizePanel.gd — T154 幸運巨型獎勵魚 Lucky Panel
## lucky-panel-agent 負責維護
## 業界依據：Jili「Giant Prize Fish lets you easily win great prizes, with the chance for 5x multipliers」
## 巨型獎勵主題：金色 + 5 次大獎演出 + 完美大獎全服加成
extends CanvasLayer

const THEME_COLOR = Color(1.0, 0.85, 0.0)          # 金色
const DARK_COLOR  = Color(0.15, 0.12, 0.0)          # 深金黑
const WHITE_COLOR = Color(1.0, 1.0, 1.0)            # 白色

var _banner: Control = null
var _prize_label: Label = null
var _prize_no: int = 0

func _ready() -> void:
	layer = 49
	_create_ui()
	GameManager.lucky_giant_prize.connect(_on_lucky_giant_prize)

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
	title.text = "🎁 巨型獎勵魚！5 次大獎！"
	title.add_theme_color_override("font_color", THEME_COLOR)
	title.add_theme_font_size_override("font_size", 20)
	title.position = Vector2(20, 12)
	_banner.add_child(title)

	_prize_label = Label.new()
	_prize_label.text = "第 0/5 次"
	_prize_label.add_theme_color_override("font_color", WHITE_COLOR)
	_prize_label.add_theme_font_size_override("font_size", 16)
	_prize_label.position = Vector2(380, 16)
	_banner.add_child(_prize_label)

func _on_lucky_giant_prize(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"giant_prize_start":
			_prize_no = 0
			_show_banner()
			_flash_screen(THEME_COLOR, 2)
			ScreenShake.add_trauma(0.4)
		"prize_drop":
			_prize_no = data.get("prize_no", _prize_no)
			var mult = data.get("prize_mult", 1.0)
			var reward = data.get("reward", 0)
			if _prize_label:
				_prize_label.text = "第 %d/5 次" % _prize_no
			_show_prize_drop(mult, reward)
		"giant_prize_end":
			_hide_banner()
		"giant_prize_perfect":
			_show_perfect(data.get("total_reward", 0), data.get("boost_mult", 4.5), data.get("boost_secs", 10))
		"giant_prize_perfect_end":
			_hide_banner()

func _show_banner() -> void:
	if _banner:
		_banner.visible = true
		_banner.modulate.a = 0.0
		var tween = create_tween()
		tween.tween_property(_banner, "modulate:a", 1.0, 0.3)

func _hide_banner() -> void:
	if _banner:
		var tween = create_tween()
		tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
		tween.tween_callback(func(): _banner.visible = false)

func _show_prize_drop(mult: float, reward: int) -> void:
	var lbl = Label.new()
	lbl.text = "🎁 ×%.0f  +%d" % [mult, reward]
	var color = THEME_COLOR
	if mult >= 30:
		color = Color(1.0, 0.3, 0.1)
		ScreenShake.add_trauma(0.3)
	elif mult >= 20:
		color = Color(1.0, 0.6, 0.0)
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 22)
	lbl.position = Vector2(randf_range(300, 700), 350)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "scale", Vector2(1.3, 1.3), 0.1)
	tween.tween_property(lbl, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_property(lbl, "position:y", lbl.position.y - 70, 0.7)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.7)
	tween.tween_callback(lbl.queue_free)

func _show_perfect(total: int, mult: float, secs: int) -> void:
	_flash_screen(THEME_COLOR, 3)
	ScreenShake.add_trauma(0.6)
	var popup = Label.new()
	popup.text = "🎁✨ 完美大獎！獲得 %d！全服 ×%.1f 加成 %d 秒！" % [total, mult, secs]
	popup.add_theme_color_override("font_color", THEME_COLOR)
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
