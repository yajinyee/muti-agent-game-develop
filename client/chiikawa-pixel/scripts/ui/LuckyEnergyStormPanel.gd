## LuckyEnergyStormPanel.gd — T201 幸運能量風暴魚 UI
## lucky-panel-agent 負責維護
## DAY-319：能量風暴系統 — 5 波連鎖電擊（每波全場 HP -30%），全服 ×26.0 加成 52 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.0, 1.0, 1.0)   # 青藍色（能量風暴）
const PANEL_ICON = "⚡"
const PANEL_TITLE = "能量風暴"

var _wave_label: Label = null
var _hit_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 96
	_setup_energy_storm_ui()
	GameManager.lucky_energy_storm.connect(_on_lucky_energy_storm)

func _setup_energy_storm_ui() -> void:
	_wave_label = Label.new()
	_wave_label.text = "波次: 0 / 5"
	_wave_label.add_theme_color_override("font_color", Color(0.0, 1.0, 1.0))
	_wave_label.add_theme_font_size_override("font_size", 22)
	_wave_label.position = Vector2(20, 52)
	add_child(_wave_label)

	_hit_label = Label.new()
	_hit_label.text = "總命中: 0"
	_hit_label.add_theme_color_override("font_color", Color(0.8, 1.0, 1.0))
	_hit_label.add_theme_font_size_override("font_size", 20)
	_hit_label.position = Vector2(20, 80)
	add_child(_hit_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×26.0 加成 52 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 108)
	add_child(_boost_label)

func _on_lucky_energy_storm(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"energy_storm_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 觸發！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.0, 1.0, 1.0))
		"storm_wave":
			var wave = data.get("wave", 0)
			var hit_count = data.get("hit_count", 0)
			if is_instance_valid(_wave_label):
				_wave_label.text = "波次: " + str(wave) + " / 5"
			if is_instance_valid(_hit_label):
				_hit_label.text = "本波命中: " + str(hit_count)
		"storm_complete":
			var total_hit = data.get("total_hit", 0)
			var boost_mult = data.get("boost_mult", 26.0)
			var boost_secs = data.get("boost_secs", 52)
			show_settle(PANEL_ICON + " 能量風暴完成！",
				"總命中 " + str(total_hit) + " 個！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！",
				PANEL_COLOR)
			flash_screen(Color(0.0, 1.0, 1.0))
			hide_panel()
