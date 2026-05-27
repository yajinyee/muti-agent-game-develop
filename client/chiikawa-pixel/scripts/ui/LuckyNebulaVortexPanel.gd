## LuckyNebulaVortexPanel.gd — T189 幸運星雲漩渦魚 UI
## lucky-panel-agent 負責維護
## DAY-316：星雲漩渦系統 — 每秒全場 HP -8%，持續 20 秒，累積命中 ≥20 → 全服 ×18.5 加成 39 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.29, 0.10, 0.40)  # 深紫色（星雲）
const PANEL_ICON = "🌌"
const PANEL_TITLE = "星雲漩渦"

var _wave_label: Label = null
var _hit_label: Label = null
var _current_wave: int = 0
var _total_hits: int = 0

func _ready() -> void:
	super._ready()
	layer = 84
	_setup_nebula_vortex_ui()
	GameManager.lucky_nebula_vortex.connect(_on_lucky_nebula_vortex)

func _setup_nebula_vortex_ui() -> void:
	_wave_label = Label.new()
	_wave_label.text = "波次: 0 / 20"
	_wave_label.add_theme_color_override("font_color", Color(0.7, 0.4, 1.0))
	_wave_label.add_theme_font_size_override("font_size", 22)
	_wave_label.position = Vector2(20, 80)
	add_child(_wave_label)

	_hit_label = Label.new()
	_hit_label.text = "累積命中: 0 次（每秒 HP -8%）"
	_hit_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_hit_label.add_theme_font_size_override("font_size", 16)
	_hit_label.position = Vector2(20, 115)
	add_child(_hit_label)

func _on_lucky_nebula_vortex(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"nebula_vortex_start":
			_current_wave = 0
			_total_hits = 0
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 召喚！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.4, 0.0, 0.8))
		"nebula_wave":
			_current_wave = data.get("wave", 0)
			_total_hits += data.get("hit_count", 0)
			if is_instance_valid(_wave_label):
				_wave_label.text = "波次: " + str(_current_wave) + " / 20"
			if is_instance_valid(_hit_label):
				_hit_label.text = "累積命中: " + str(_total_hits) + " 次"
		"nebula_perfect":
			var boost_mult = data.get("boost_mult", 18.5)
			var boost_secs = data.get("boost_secs", 39)
			show_settle(PANEL_ICON + " 星雲完美！", "累積命中 " + str(data.get("total_hits", 0)) + " 次！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", PANEL_COLOR)
			flash_screen(Color(0.5, 0.1, 1.0))
			hide_panel()
		"nebula_vortex_end":
			hide_panel()
