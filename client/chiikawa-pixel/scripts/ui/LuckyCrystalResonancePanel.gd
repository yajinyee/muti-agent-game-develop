## LuckyCrystalResonancePanel.gd — T202 幸運水晶共鳴魚 UI
## lucky-panel-agent 負責維護
## DAY-319：水晶共鳴系統 — 全場共鳴爆炸（每個獎勵 ×30.0），全服 ×27.0 加成 54 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.88, 0.88, 1.0)  # 水晶白藍色
const PANEL_ICON = "💎"
const PANEL_TITLE = "水晶共鳴"

var _hit_label: Label = null
var _mult_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 97
	_setup_crystal_resonance_ui()
	GameManager.lucky_crystal_resonance.connect(_on_lucky_crystal_resonance)

func _setup_crystal_resonance_ui() -> void:
	_mult_label = Label.new()
	_mult_label.text = "全場共鳴爆炸 ×30.0"
	_mult_label.add_theme_color_override("font_color", Color(0.88, 0.88, 1.0))
	_mult_label.add_theme_font_size_override("font_size", 22)
	_mult_label.position = Vector2(20, 52)
	add_child(_mult_label)

	_hit_label = Label.new()
	_hit_label.text = "清場: 0 個"
	_hit_label.add_theme_color_override("font_color", Color(0.8, 0.8, 1.0))
	_hit_label.add_theme_font_size_override("font_size", 20)
	_hit_label.position = Vector2(20, 80)
	add_child(_hit_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×27.0 加成 54 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 108)
	add_child(_boost_label)

func _on_lucky_crystal_resonance(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"crystal_resonance_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 觸發！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.88, 0.88, 1.0))
		"crystal_resonance_complete":
			var hit_count = data.get("hit_count", 0)
			var boost_mult = data.get("boost_mult", 27.0)
			var boost_secs = data.get("boost_secs", 54)
			if is_instance_valid(_hit_label):
				_hit_label.text = "清場: " + str(hit_count) + " 個"
			show_settle(PANEL_ICON + " 水晶共鳴完成！",
				"清場 " + str(hit_count) + " 個！每個 ×30.0！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！",
				PANEL_COLOR)
			flash_screen(Color(0.88, 0.88, 1.0))
			hide_panel()
