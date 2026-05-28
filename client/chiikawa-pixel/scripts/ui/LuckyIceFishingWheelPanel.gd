## LuckyIceFishingWheelPanel.gd — T214 幸運冰釣輪盤魚 UI
## lucky-panel-agent 負責維護
## DAY-324：冰釣輪盤系統 — 3 次旋轉（最高 ×5000），最高單次 ≥2000 → 全服 ×37.5 加成 75 秒
## 業界依據：Evolution Gaming「Ice Fishing」最高 5000x（2026）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.0, 0.81, 0.82)   # 冰青色（冰釣）
const PANEL_ICON = "❄️"
const PANEL_TITLE = "冰釣輪盤"

var _spin_label: Label = null
var _result_label: Label = null
var _total_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 109
	_setup_wheel_ui()
	GameManager.lucky_ice_fishing_wheel.connect(_on_lucky_ice_fishing_wheel)

func _setup_wheel_ui() -> void:
	_spin_label = Label.new()
	_spin_label.text = "旋轉: 0 / 3"
	_spin_label.add_theme_color_override("font_color", Color(0.0, 0.9, 1.0))
	_spin_label.add_theme_font_size_override("font_size", 22)
	_spin_label.position = Vector2(20, 52)
	add_child(_spin_label)

	_result_label = Label.new()
	_result_label.text = "本次: ---"
	_result_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	_result_label.add_theme_font_size_override("font_size", 24)
	_result_label.position = Vector2(20, 78)
	add_child(_result_label)

	_total_label = Label.new()
	_total_label.text = "累計: ×0"
	_total_label.add_theme_color_override("font_color", Color(1.0, 0.84, 0.0))
	_total_label.add_theme_font_size_override("font_size", 18)
	_total_label.position = Vector2(20, 108)
	add_child(_total_label)

	_boost_label = Label.new()
	_boost_label.text = "最高 ×5000！≥×2000 → 全服 ×37.5 加成 75 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 13)
	_boost_label.position = Vector2(20, 132)
	add_child(_boost_label)

func _on_lucky_ice_fishing_wheel(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"wheel_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開始！", PANEL_COLOR)
			show_panel()
			flash_screen(PANEL_COLOR)
		"spin_result":
			var spin = data.get("spin", 0)
			var mult = data.get("mult", 0.0)
			var label = data.get("label", "")
			var total_mult = data.get("total_mult", 0.0)
			if is_instance_valid(_spin_label):
				_spin_label.text = "旋轉: " + str(spin) + " / 3"
			if is_instance_valid(_result_label):
				_result_label.text = "本次: " + label
			if is_instance_valid(_total_label):
				_total_label.text = "累計: ×" + str(int(total_mult))
			# 高倍率閃光
			if mult >= 2000:
				flash_screen(Color(1.0, 1.0, 1.0))
			elif mult >= 1000:
				flash_screen(PANEL_COLOR)
		"wheel_jackpot":
			var max_mult = data.get("max_mult", 0.0)
			var global_mult = data.get("global_mult", 37.5)
			var global_secs = data.get("global_secs", 75)
			show_settle(PANEL_ICON + " 冰釣大獎！",
				"最高 ×" + str(int(max_mult)) + "！全服 ×" + str(global_mult) + " 加成 " + str(global_secs) + " 秒！",
				PANEL_COLOR)
			flash_screen(Color(1.0, 1.0, 1.0))
			hide_panel()
		"wheel_end":
			var max_mult = data.get("max_mult", 0.0)
			show_settle(PANEL_ICON + " 冰釣結束",
				"最高 ×" + str(int(max_mult)),
				PANEL_COLOR)
			hide_panel()
