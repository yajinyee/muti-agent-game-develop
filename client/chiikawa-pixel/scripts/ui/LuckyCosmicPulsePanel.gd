## LuckyCosmicPulsePanel.gd — T185 幸運宇宙脈衝魚 UI
## lucky-panel-agent 負責維護
## DAY-315：宇宙脈衝系統 — 全場 HP -45%（每個獎勵 ×12.0），觸發全服 ×16.0 加成 35 秒（新最高）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.29, 0.10, 0.57)  # 深紫色（宇宙）
const PANEL_ICON = "🌌"
const PANEL_TITLE = "宇宙脈衝"

var _hit_count: int = 0
var _hit_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 80
	_setup_cosmic_pulse_ui()
	GameManager.lucky_cosmic_pulse.connect(_on_lucky_cosmic_pulse)

func _setup_cosmic_pulse_ui() -> void:
	_hit_label = Label.new()
	_hit_label.text = "清場: 0 個"
	_hit_label.add_theme_color_override("font_color", Color(0.7, 0.4, 1.0))
	_hit_label.add_theme_font_size_override("font_size", 22)
	_hit_label.position = Vector2(20, 80)
	add_child(_hit_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×16.0 加成 35 秒（新最高）"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 115)
	add_child(_boost_label)

func _on_lucky_cosmic_pulse(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"cosmic_pulse_start":
			_hit_count = 0
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 引動！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.3, 0.0, 0.8))
		"cosmic_pulse_complete":
			_hit_count = data.get("hit_count", 0)
			var boost_mult = data.get("boost_mult", 16.0)
			var boost_secs = data.get("boost_secs", 35)
			if is_instance_valid(_hit_label):
				_hit_label.text = "清場: " + str(_hit_count) + " 個"
			show_settle(PANEL_ICON + " 宇宙脈衝！", "清場 " + str(_hit_count) + " 個！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！（新最高）", PANEL_COLOR)
			flash_screen(Color(0.4, 0.0, 1.0))
			hide_panel()
