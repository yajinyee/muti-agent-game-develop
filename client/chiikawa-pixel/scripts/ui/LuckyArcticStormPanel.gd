## LuckyArcticStormPanel.gd — T182 幸運北極風暴魚 UI
## lucky-panel-agent 負責維護
## DAY-315：北極風暴系統 — 8 波快速冰雪攻擊（每 0.3 秒），全部命中 → 全服 ×16.5 加成 33 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.01, 0.53, 0.82)  # 冰藍色（北極）
const PANEL_ICON = "❄️"
const PANEL_TITLE = "北極風暴"

var _wave_label: Label = null
var _hit_label: Label = null
var _perfect_waves: int = 0

func _ready() -> void:
	super._ready()
	layer = 77
	_setup_arctic_storm_ui()
	GameManager.lucky_arctic_storm.connect(_on_lucky_arctic_storm)

func _setup_arctic_storm_ui() -> void:
	_wave_label = Label.new()
	_wave_label.text = "波次：0 / 8"
	_wave_label.add_theme_color_override("font_color", Color(0.5, 0.9, 1.0))
	_wave_label.add_theme_font_size_override("font_size", 22)
	_wave_label.position = Vector2(20, 80)
	add_child(_wave_label)

	_hit_label = Label.new()
	_hit_label.text = "命中：0 個"
	_hit_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	_hit_label.add_theme_font_size_override("font_size", 18)
	_hit_label.position = Vector2(20, 110)
	add_child(_hit_label)

func _on_lucky_arctic_storm(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"arctic_storm_start":
			_perfect_waves = 0
			if is_instance_valid(_wave_label):
				_wave_label.text = "波次：0 / 8"
			if is_instance_valid(_hit_label):
				_hit_label.text = "命中：0 個"
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 來襲！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.0, 0.5, 1.0))
		"arctic_wave":
			var wave = data.get("wave", 0)
			var total_hit = data.get("total_hit", 0)
			if is_instance_valid(_wave_label):
				_wave_label.text = "波次：" + str(wave) + " / 8"
			if is_instance_valid(_hit_label):
				_hit_label.text = "命中：" + str(total_hit) + " 個"
		"arctic_storm_perfect":
			var total_hit = data.get("total_hit", 0)
			var boost_mult = data.get("boost_mult", 16.5)
			var boost_secs = data.get("boost_secs", 33)
			show_settle(PANEL_ICON + " 完美北極風暴！", "8 波全中！命中 " + str(total_hit) + " 個！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", PANEL_COLOR)
			flash_screen(Color(0.0, 0.7, 1.0))
			hide_panel()
		"arctic_storm_complete":
			var total_hit = data.get("total_hit", 0)
			var perfect_waves = data.get("perfect_waves", 0)
			show_settle(PANEL_ICON + " 北極風暴完成！", "命中 " + str(total_hit) + " 個，完美波次 " + str(perfect_waves) + " / 8", PANEL_COLOR)
			hide_panel()
