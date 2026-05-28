## LuckyMultiplierLadderPanel.gd — T213 幸運倍率梯魚 UI
## lucky-panel-agent 負責維護
## DAY-324：倍率梯系統 — 每次擊破提升梯度（Lv.1-10），Lv.10 → 全服 ×37.0 加成 74 秒
## 業界依據：Relax Gaming「Cod of Thunder Dream Drop」Multiplier Ladder（2026）
extends BaseLuckyPanel

const PANEL_COLOR = Color(1.0, 0.84, 0.0)   # 金色（倍率梯）
const PANEL_ICON = "🪜"
const PANEL_TITLE = "倍率梯"

var _level_label: Label = null
var _mult_label: Label = null
var _kills_label: Label = null
var _boost_label: Label = null
var _timer_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 108
	_setup_ladder_ui()
	GameManager.lucky_multiplier_ladder.connect(_on_lucky_multiplier_ladder)

func _setup_ladder_ui() -> void:
	_level_label = Label.new()
	_level_label.text = "等級: Lv.0 / 10"
	_level_label.add_theme_color_override("font_color", Color(1.0, 0.84, 0.0))
	_level_label.add_theme_font_size_override("font_size", 24)
	_level_label.position = Vector2(20, 52)
	add_child(_level_label)

	_mult_label = Label.new()
	_mult_label.text = "本級倍率: ×0.0"
	_mult_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	_mult_label.add_theme_font_size_override("font_size", 20)
	_mult_label.position = Vector2(20, 80)
	add_child(_mult_label)

	_kills_label = Label.new()
	_kills_label.text = "擊破: 0 個"
	_kills_label.add_theme_color_override("font_color", Color(0.9, 0.9, 0.9))
	_kills_label.add_theme_font_size_override("font_size", 16)
	_kills_label.position = Vector2(20, 106)
	add_child(_kills_label)

	_boost_label = Label.new()
	_boost_label.text = "Lv.10 → 全服 ×37.0 加成 74 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 14)
	_boost_label.position = Vector2(20, 130)
	add_child(_boost_label)

func _on_lucky_multiplier_ladder(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"ladder_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開始！", PANEL_COLOR)
			show_panel()
			flash_screen(PANEL_COLOR)
		"level_up":
			var level = data.get("level", 0)
			var level_mult = data.get("level_mult", 0.0)
			var kills = data.get("kills", 0)
			if is_instance_valid(_level_label):
				_level_label.text = "等級: Lv." + str(level) + " / 10"
			if is_instance_valid(_mult_label):
				_mult_label.text = "本級倍率: ×" + str(snapped(level_mult, 0.1))
			if is_instance_valid(_kills_label):
				_kills_label.text = "擊破: " + str(kills) + " 個"
			# 等級越高顏色越亮
			if level >= 8 and is_instance_valid(_level_label):
				_level_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.0))
		"ladder_max":
			var global_mult = data.get("global_mult", 37.0)
			var global_secs = data.get("global_secs", 74)
			show_settle(PANEL_ICON + " 倍率梯頂端！",
				"Lv.10 達成！全服 ×" + str(global_mult) + " 加成 " + str(global_secs) + " 秒！",
				PANEL_COLOR)
			flash_screen(PANEL_COLOR)
			hide_panel()
		"ladder_timeout":
			var kills = data.get("kills", 0)
			show_settle(PANEL_ICON + " 倍率梯超時",
				"最終等級 Lv." + str(kills),
				Color(0.7, 0.7, 0.7))
			hide_panel()
