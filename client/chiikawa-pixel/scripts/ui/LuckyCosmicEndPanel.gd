## LuckyCosmicEndPanel.gd — T195 幸運宇宙終焉魚 UI
## lucky-panel-agent 負責維護
## DAY-317：宇宙終焉系統 — 全場 HP 歸零（每個獎勵 ×20.0），觸發全服 ×22.0 加成 45 秒（史上最高）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.05, 0.0, 0.05)   # 近黑色（宇宙終焉）
const PANEL_BORDER_COLOR = Color(1.0, 0.85, 0.0)  # 金色邊框
const PANEL_ICON = "☄️"
const PANEL_TITLE = "宇宙終焉"

var _hit_label: Label = null
var _boost_label: Label = null
var _crown_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 90
	_setup_cosmic_end_ui()
	GameManager.lucky_cosmic_end.connect(_on_lucky_cosmic_end)

func _setup_cosmic_end_ui() -> void:
	_crown_label = Label.new()
	_crown_label.text = "👑 史上最高倍率 ×22.0 👑"
	_crown_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_crown_label.add_theme_font_size_override("font_size", 20)
	_crown_label.position = Vector2(20, 55)
	add_child(_crown_label)

	_hit_label = Label.new()
	_hit_label.text = "清場: 0 個"
	_hit_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.5))
	_hit_label.add_theme_font_size_override("font_size", 22)
	_hit_label.position = Vector2(20, 85)
	add_child(_hit_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×22.0 加成 45 秒（史上最高）"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 120)
	add_child(_boost_label)

func _on_lucky_cosmic_end(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"cosmic_end_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 降臨！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.0, 0.0, 0.0))
			# 特殊：黑色閃屏後金色
			var tween = create_tween()
			tween.tween_interval(0.15)
			tween.tween_callback(func(): flash_screen(Color(1.0, 0.85, 0.0)))
		"cosmic_end_complete":
			var hit_count = data.get("hit_count", 0)
			var boost_mult = data.get("boost_mult", 22.0)
			var boost_secs = data.get("boost_secs", 45)
			if is_instance_valid(_hit_label):
				_hit_label.text = "清場: " + str(hit_count) + " 個"
			show_settle(PANEL_ICON + " 宇宙終焉！",
				"清場 " + str(hit_count) + " 個！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！（史上最高）",
				PANEL_COLOR)
			flash_screen(Color(1.0, 0.85, 0.0))
			hide_panel()
