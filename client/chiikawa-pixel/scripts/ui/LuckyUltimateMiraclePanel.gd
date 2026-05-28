## LuckyUltimateMiraclePanel.gd — T210 幸運終極奇蹟魚 UI
## lucky-panel-agent 負責維護
## DAY-323：終極奇蹟系統 — 全場 HP 歸零（每個獎勵 ×50.0），全服 ×35.0 加成 70 秒（新史上最高）
## 業界依據：終極奇蹟機制 + 2026 最高倍率設計（16888x 吉祥數字）
extends BaseLuckyPanel

const PANEL_COLOR = Color(1.0, 1.0, 1.0)   # 純白色（終極奇蹟）
const PANEL_ICON = "🌟"
const PANEL_TITLE = "終極奇蹟"

var _miracle_label: Label = null
var _mult_label: Label = null
var _hit_label: Label = null
var _boost_label: Label = null
var _record_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 105
	_setup_ultimate_miracle_ui()
	GameManager.lucky_ultimate_miracle.connect(_on_lucky_ultimate_miracle)

func _setup_ultimate_miracle_ui() -> void:
	_record_label = Label.new()
	_record_label.text = "🏆 新史上最高全服倍率 ×35.0"
	_record_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_record_label.add_theme_font_size_override("font_size", 16)
	_record_label.position = Vector2(20, 52)
	add_child(_record_label)

	_mult_label = Label.new()
	_mult_label.text = "全場清空 ×50.0（史上最高單次）"
	_mult_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	_mult_label.add_theme_font_size_override("font_size", 20)
	_mult_label.position = Vector2(20, 75)
	add_child(_mult_label)

	_hit_label = Label.new()
	_hit_label.text = "清場: 0 個"
	_hit_label.add_theme_color_override("font_color", Color(0.8, 0.8, 1.0))
	_hit_label.add_theme_font_size_override("font_size", 22)
	_hit_label.position = Vector2(20, 100)
	add_child(_hit_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×35.0 加成 70 秒（新史上最高）"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 130)
	add_child(_boost_label)

func _on_lucky_ultimate_miracle(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"miracle_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開啟！", PANEL_COLOR)
			show_panel()
			# 四重閃屏：白 → 金 → 白 → 金
			flash_screen(Color(1.0, 1.0, 1.0))
			var tween = create_tween()
			tween.tween_interval(0.12)
			tween.tween_callback(func(): flash_screen(Color(1.0, 0.85, 0.0)))
			tween.tween_interval(0.12)
			tween.tween_callback(func(): flash_screen(Color(1.0, 1.0, 1.0)))
			tween.tween_interval(0.12)
			tween.tween_callback(func(): flash_screen(Color(1.0, 0.85, 0.0)))
		"miracle_complete":
			var hit_count = data.get("hit_count", 0)
			var boost_mult = data.get("boost_mult", 35.0)
			var boost_secs = data.get("boost_secs", 70)
			if is_instance_valid(_hit_label):
				_hit_label.text = "清場: " + str(hit_count) + " 個"
			show_settle(PANEL_ICON + " 終極奇蹟！",
				"清場 " + str(hit_count) + " 個！每個 ×50.0！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！（新史上最高）",
				PANEL_COLOR)
			# 最強閃屏：白色全螢幕
			flash_screen(Color(1.0, 1.0, 1.0))
			var tween2 = create_tween()
			tween2.tween_interval(0.1)
			tween2.tween_callback(func(): flash_screen(Color(1.0, 1.0, 1.0)))
			tween2.tween_interval(0.1)
			tween2.tween_callback(func(): flash_screen(Color(1.0, 0.85, 0.0)))
			hide_panel()
