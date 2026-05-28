## LuckyEternalCyclePanel.gd — T197 幸運永恆循環魚 UI
## lucky-panel-agent 負責維護
## DAY-318：永恆循環系統 — 10 波遞增獎勵（×1.0 → ×10.0），全服 ×23.5 加成 47 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.0, 0.05, 0.15)   # 深藍色（永恆）
const PANEL_BORDER_COLOR = Color(0.4, 0.6, 1.0)  # 藍色邊框
const PANEL_ICON = "♾️"
const PANEL_TITLE = "永恆循環"

var _wave_label: Label = null
var _mult_label: Label = null
var _total_label: Label = null
var _boost_label: Label = null
var _current_wave: int = 0

func _ready() -> void:
	super._ready()
	layer = 92
	_setup_eternal_cycle_ui()
	GameManager.lucky_eternal_cycle.connect(_on_lucky_eternal_cycle)

func _setup_eternal_cycle_ui() -> void:
	_wave_label = Label.new()
	_wave_label.text = "波次: 0 / 10"
	_wave_label.add_theme_color_override("font_color", Color(0.6, 0.8, 1.0))
	_wave_label.add_theme_font_size_override("font_size", 20)
	_wave_label.position = Vector2(20, 55)
	add_child(_wave_label)

	_mult_label = Label.new()
	_mult_label.text = "本波倍率: ×?"
	_mult_label.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
	_mult_label.add_theme_font_size_override("font_size", 22)
	_mult_label.position = Vector2(20, 80)
	add_child(_mult_label)

	_total_label = Label.new()
	_total_label.text = "累積獎勵: 0"
	_total_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.0))
	_total_label.add_theme_font_size_override("font_size", 18)
	_total_label.position = Vector2(20, 108)
	add_child(_total_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×23.5 加成 47 秒"
	_boost_label.add_theme_color_override("font_color", Color(0.6, 0.8, 1.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 135)
	add_child(_boost_label)

func _on_lucky_eternal_cycle(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"eternal_cycle_start":
			_current_wave = 0
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 啟動！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.4, 0.6, 1.0))
		"eternal_cycle_wave":
			_current_wave = data.get("wave", 0)
			var wave_mult = data.get("wave_mult", 1.0)
			if is_instance_valid(_wave_label):
				_wave_label.text = "波次: " + str(_current_wave) + " / 10"
			if is_instance_valid(_mult_label):
				_mult_label.text = "本波倍率: ×" + str(int(wave_mult))
		"eternal_cycle_complete":
			var total_reward = data.get("total_reward", 0)
			var boost_mult = data.get("boost_mult", 23.5)
			var boost_secs = data.get("boost_secs", 47)
			if is_instance_valid(_total_label):
				_total_label.text = "累積獎勵: " + str(total_reward)
			show_settle(PANEL_ICON + " 永恆循環完成！",
				"10 波全部完成！總獎勵 " + str(total_reward) + "！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！",
				PANEL_COLOR)
			flash_screen(Color(0.4, 0.6, 1.0))
			hide_panel()
