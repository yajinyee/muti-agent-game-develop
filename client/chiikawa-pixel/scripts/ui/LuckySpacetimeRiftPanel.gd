## LuckySpacetimeRiftPanel.gd — T168 幸運時空裂縫魚 UI
## lucky-panel-agent 負責維護
## DAY-312：時空裂縫系統 — 20 秒每 4 秒瞬間擊破 3 個目標（獎勵 ×4.0），擊破 ≥12 → 時空完美全服 ×7.5 加成 16 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.08, 0.27, 0.75)  # 深藍色（時空）
const PANEL_ICON = "⏳"
const PANEL_TITLE = "時空裂縫"

var _kill_count: int = 0
var _target_kills: int = 12
var _current_wave: int = 0
var _kill_label: Label = null
var _wave_label: Label = null
var _kill_bar: ProgressBar = null

func _ready() -> void:
	super._ready()
	layer = 63
	_setup_rift_ui()
	GameManager.lucky_spacetime_rift.connect(_on_lucky_spacetime_rift)

func _setup_rift_ui() -> void:
	# 波次顯示
	_wave_label = Label.new()
	_wave_label.text = "波次: 0/5"
	_wave_label.add_theme_color_override("font_color", Color(0.4, 0.7, 1.0))
	_wave_label.add_theme_font_size_override("font_size", 18)
	_wave_label.position = Vector2(20, 80)
	add_child(_wave_label)

	# 擊破計數
	_kill_label = Label.new()
	_kill_label.text = "擊破: 0/12"
	_kill_label.add_theme_color_override("font_color", Color(0.6, 0.85, 1.0))
	_kill_label.add_theme_font_size_override("font_size", 18)
	_kill_label.position = Vector2(20, 108)
	add_child(_kill_label)

	# 擊破進度條
	_kill_bar = ProgressBar.new()
	_kill_bar.min_value = 0
	_kill_bar.max_value = _target_kills
	_kill_bar.value = 0
	_kill_bar.size = Vector2(200, 16)
	_kill_bar.position = Vector2(20, 136)
	add_child(_kill_bar)

func _on_lucky_spacetime_rift(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"rift_open":
			_kill_count = 0
			_current_wave = 0
			_update_rift_display()
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 撕裂！", PANEL_COLOR)
			start_timer(data.get("duration", 20))
			show_panel()
			flash_screen(PANEL_COLOR)
		"rift_wave":
			_current_wave = data.get("wave", 0)
			_kill_count += data.get("kill_count", 0)
			_update_rift_display()
		"rift_perfect":
			var boost_mult = data.get("boost_mult", 7.5)
			var boost_secs = data.get("boost_secs", 16)
			show_settle("⏳ 時空完美！", "擊破 " + str(_kill_count) + " 個！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", PANEL_COLOR)
			flash_screen(Color(0.05, 0.28, 0.9))
			hide_panel()
		"rift_end":
			show_settle("時空裂縫關閉", "擊破 " + str(_kill_count) + " 個目標", PANEL_COLOR)
			hide_panel()

func _update_rift_display() -> void:
	if is_instance_valid(_wave_label):
		_wave_label.text = "波次: " + str(_current_wave) + "/5"
	if is_instance_valid(_kill_label):
		_kill_label.text = "擊破: " + str(_kill_count) + "/" + str(_target_kills)
	if is_instance_valid(_kill_bar):
		_kill_bar.value = min(_kill_count, _target_kills)
