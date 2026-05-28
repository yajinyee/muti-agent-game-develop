## LuckyAvalanchePanel.gd — T211 幸運雪崩魚 UI
## lucky-panel-agent 負責維護
## DAY-324：雪崩連鎖系統 — 8 波連鎖消除，每波倍率 +5.0，全服 ×36.0 加成 72 秒
## 業界依據：Relax Gaming「Cod of Thunder Dream Drop」Multiplier Ladder（2026）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.0, 0.75, 1.0)   # 冰藍色（雪崩）
const PANEL_ICON = "❄️"
const PANEL_TITLE = "雪崩連鎖"

var _wave_label: Label = null
var _mult_label: Label = null
var _hit_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 106
	_setup_avalanche_ui()
	GameManager.lucky_avalanche.connect(_on_lucky_avalanche)

func _setup_avalanche_ui() -> void:
	_wave_label = Label.new()
	_wave_label.text = "波次: 0 / 8"
	_wave_label.add_theme_color_override("font_color", Color(0.0, 0.9, 1.0))
	_wave_label.add_theme_font_size_override("font_size", 22)
	_wave_label.position = Vector2(20, 52)
	add_child(_wave_label)

	_mult_label = Label.new()
	_mult_label.text = "本波倍率: ×0.0"
	_mult_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	_mult_label.add_theme_font_size_override("font_size", 20)
	_mult_label.position = Vector2(20, 78)
	add_child(_mult_label)

	_hit_label = Label.new()
	_hit_label.text = "命中波次: 0"
	_hit_label.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
	_hit_label.add_theme_font_size_override("font_size", 18)
	_hit_label.position = Vector2(20, 104)
	add_child(_hit_label)

	_boost_label = Label.new()
	_boost_label.text = "完美雪崩 → 全服 ×36.0 加成 72 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 14)
	_boost_label.position = Vector2(20, 130)
	add_child(_boost_label)

func _on_lucky_avalanche(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"avalanche_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開始！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.0, 0.75, 1.0))
		"wave_hit":
			var wave = data.get("wave", 0)
			var wave_mult = data.get("wave_mult", 0.0)
			var hit_count = data.get("hit_count", 0)
			if is_instance_valid(_wave_label):
				_wave_label.text = "波次: " + str(wave) + " / 8"
			if is_instance_valid(_mult_label):
				_mult_label.text = "本波倍率: ×" + str(wave_mult)
			if is_instance_valid(_hit_label):
				_hit_label.text = "命中波次: " + str(hit_count)
		"avalanche_perfect":
			var hit_count = data.get("hit_count", 0)
			var global_mult = data.get("global_mult", 36.0)
			var global_secs = data.get("global_secs", 72)
			show_settle(PANEL_ICON + " 完美雪崩！",
				str(hit_count) + " 波命中！全服 ×" + str(global_mult) + " 加成 " + str(global_secs) + " 秒！",
				PANEL_COLOR)
			flash_screen(Color(0.0, 0.75, 1.0))
			hide_panel()
		"avalanche_end":
			hide_panel()
