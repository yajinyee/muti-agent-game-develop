## LuckyDragonKingPanel.gd — T196 幸運龍王輪盤魚 UI
## lucky-panel-agent 負責維護
## DAY-318：龍王輪盤系統 — 雙環輪盤，內環 × 外環 = 最高 ×25.0，全服 ×23.0 加成 46 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.15, 0.05, 0.0)   # 深橙色（龍王）
const PANEL_BORDER_COLOR = Color(1.0, 0.7, 0.0)  # 金橙色邊框
const PANEL_ICON = "🐉"
const PANEL_TITLE = "龍王輪盤"

var _inner_label: Label = null
var _outer_label: Label = null
var _total_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 91
	_setup_dragon_king_ui()
	GameManager.lucky_dragon_king.connect(_on_lucky_dragon_king)

func _setup_dragon_king_ui() -> void:
	_inner_label = Label.new()
	_inner_label.text = "內環: ×?"
	_inner_label.add_theme_color_override("font_color", Color(1.0, 0.8, 0.2))
	_inner_label.add_theme_font_size_override("font_size", 20)
	_inner_label.position = Vector2(20, 55)
	add_child(_inner_label)

	_outer_label = Label.new()
	_outer_label.text = "外環: ×?"
	_outer_label.add_theme_color_override("font_color", Color(1.0, 0.6, 0.1))
	_outer_label.add_theme_font_size_override("font_size", 20)
	_outer_label.position = Vector2(20, 80)
	add_child(_outer_label)

	_total_label = Label.new()
	_total_label.text = "合計: ×?"
	_total_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.0))
	_total_label.add_theme_font_size_override("font_size", 24)
	_total_label.position = Vector2(20, 108)
	add_child(_total_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×23.0 加成 46 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 138)
	add_child(_boost_label)

func _on_lucky_dragon_king(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"dragon_king_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 啟動！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(1.0, 0.7, 0.0))
		"dragon_king_complete":
			var inner_mult = data.get("inner_mult", 1.0)
			var outer_mult = data.get("outer_mult", 1.0)
			var total_mult = data.get("total_mult", 1.0)
			var boost_mult = data.get("boost_mult", 23.0)
			var boost_secs = data.get("boost_secs", 46)
			if is_instance_valid(_inner_label):
				_inner_label.text = "內環: ×" + str(int(inner_mult))
			if is_instance_valid(_outer_label):
				_outer_label.text = "外環: ×" + str(int(outer_mult))
			if is_instance_valid(_total_label):
				_total_label.text = "合計: ×" + str(snapped(total_mult, 0.1))
			show_settle(PANEL_ICON + " 龍王輪盤！",
				"內環 ×" + str(int(inner_mult)) + " × 外環 ×" + str(int(outer_mult)) + " = ×" + str(snapped(total_mult, 0.1)) + "！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！",
				PANEL_COLOR)
			flash_screen(Color(1.0, 0.85, 0.0))
			hide_panel()
