## LuckyGenesisPanel.gd — T149 幸運創世魚 Lucky Panel
## lucky-panel-agent 負責維護
## 業界依據：Ultimate boss mechanic — 召喚創世神，全場目標直接擊破
## 創世主題：金色 + 神聖光芒 + 審判計數器
extends CanvasLayer

const THEME_COLOR = Color(1.0, 0.85, 0.0)        # 金色
const DARK_COLOR  = Color(0.2, 0.15, 0.0)         # 深金棕
const WHITE_COLOR = Color(1.0, 1.0, 1.0)          # 白色

var _banner: Control = null
var _kill_label: Label = null
var _reward_label: Label = null

func _ready() -> void:
	layer = 44
	_create_ui()
	GameManager.lucky_genesis.connect(_on_lucky_genesis)

func _create_ui() -> void:
	_banner = Control.new()
	_banner.visible = false
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 56)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.color = Color(DARK_COLOR.r, DARK_COLOR.g, DARK_COLOR.b, 0.92)
	bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	_banner.add_child(bg)

	var top_line = ColorRect.new()
	top_line.color = THEME_COLOR
	top_line.size = Vector2(1280, 4)
	top_line.position = Vector2(0, 0)
	_banner.add_child(top_line)

	var title = Label.new()
	title.text = "🌟 創世神降臨！"
	title.add_theme_color_override("font_color", THEME_COLOR)
	title.add_theme_font_size_override("font_size", 22)
	title.position = Vector2(20, 10)
	_banner.add_child(title)

	_kill_label = Label.new()
	_kill_label.text = "審判：0 個"
	_kill_label.add_theme_color_override("font_color", WHITE_COLOR)
	_kill_label.add_theme_font_size_override("font_size", 16)
	_kill_label.position = Vector2(320, 16)
	_banner.add_child(_kill_label)

	_reward_label = Label.new()
	_reward_label.text = "每個 ×5.0"
	_reward_label.add_theme_color_override("font_color", THEME_COLOR)
	_reward_label.add_theme_font_size_override("font_size", 14)
	_reward_label.position = Vector2(500, 18)
	_banner.add_child(_reward_label)

func _on_lucky_genesis(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"genesis_descend":
			_show_banner()
			_flash_screen(THEME_COLOR, 3)
			ScreenShake.add_trauma(0.8)
		"genesis_judgment":
			var kills = data.get("kill_count", 0)
			if _kill_label:
				_kill_label.text = "審判：%d 個" % kills
			_flash_screen(WHITE_COLOR, 2)
			ScreenShake.add_trauma(1.0)
		"genesis_blessing":
			_show_blessing(data.get("kill_count", 0), data.get("boost_mult", 6.0), data.get("boost_sec", 15))
		"genesis_blessing_end":
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

func _show_blessing(kills: int, mult: float, secs: int) -> void:
	_flash_screen(THEME_COLOR, 4)
	var popup = Label.new()
	popup.text = "🌟 創世祝福！審判 %d 個！全服 ×%.0f 加成 %d 秒！" % [kills, mult, secs]
	popup.add_theme_color_override("font_color", THEME_COLOR)
	popup.add_theme_font_size_override("font_size", 24)
	popup.position = Vector2(150, 300)
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
		tween.tween_property(overlay, "color:a", 0.45, 0.08)
		tween.tween_property(overlay, "color:a", 0.0, 0.12)
	tween.tween_callback(overlay.queue_free)
