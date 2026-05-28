## LuckyCosmicSingularityPanel.gd — T205 幸運宇宙奇點魚 UI
## lucky-panel-agent 負責維護
## DAY-319：宇宙奇點系統 — 全場 HP 歸零（每個獎勵 ×30.0），全服 ×30.0 加成 60 秒（史上最高）
extends BaseLuckyPanel

const PANEL_COLOR = Color(1.0, 0.0, 1.0)   # 洋紅色（宇宙奇點）
const PANEL_ICON = "🌌"
const PANEL_TITLE = "宇宙奇點"

var _hit_label: Label = null
var _mult_label: Label = null
var _boost_label: Label = null
var _record_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 100
	_setup_cosmic_singularity_ui()
	GameManager.lucky_cosmic_singularity.connect(_on_lucky_cosmic_singularity)

func _setup_cosmic_singularity_ui() -> void:
	_record_label = Label.new()
	_record_label.text = "🏆 史上最高全服倍率 ×30.0"
	_record_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	_record_label.add_theme_font_size_override("font_size", 16)
	_record_label.position = Vector2(20, 52)
	add_child(_record_label)

	_mult_label = Label.new()
	_mult_label.text = "全場清空 ×30.0（史上最高）"
	_mult_label.add_theme_color_override("font_color", Color(1.0, 0.0, 1.0))
	_mult_label.add_theme_font_size_override("font_size", 20)
	_mult_label.position = Vector2(20, 75)
	add_child(_mult_label)

	_hit_label = Label.new()
	_hit_label.text = "清場: 0 個"
	_hit_label.add_theme_color_override("font_color", Color(1.0, 0.5, 1.0))
	_hit_label.add_theme_font_size_override("font_size", 22)
	_hit_label.position = Vector2(20, 100)
	add_child(_hit_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×30.0 加成 60 秒（史上最高）"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 130)
	add_child(_boost_label)

func _on_lucky_cosmic_singularity(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"singularity_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開啟！", PANEL_COLOR)
			show_panel()
			# 三重閃屏：洋紅 → 白 → 洋紅
			flash_screen(Color(1.0, 0.0, 1.0))
			var tween = create_tween()
			tween.tween_interval(0.15)
			tween.tween_callback(func(): flash_screen(Color(1.0, 1.0, 1.0)))
			tween.tween_interval(0.15)
			tween.tween_callback(func(): flash_screen(Color(1.0, 0.0, 1.0)))
		"singularity_complete":
			var hit_count = data.get("hit_count", 0)
			var boost_mult = data.get("boost_mult", 30.0)
			var boost_secs = data.get("boost_secs", 60)
			if is_instance_valid(_hit_label):
				_hit_label.text = "清場: " + str(hit_count) + " 個"
			show_settle(PANEL_ICON + " 宇宙奇點！",
				"清場 " + str(hit_count) + " 個！每個 ×30.0！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！（史上最高）",
				PANEL_COLOR)
			flash_screen(Color(1.0, 1.0, 1.0))
			hide_panel()
