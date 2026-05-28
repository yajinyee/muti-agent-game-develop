## LuckyCrashMultiplierPanel.gd — T212 幸運崩潰倍率魚 UI
## lucky-panel-agent 負責維護
## DAY-324：崩潰倍率系統 — 倍率持續上升直到崩潰，完美收割（≥40.0）→ 全服 ×36.5 加成 73 秒
## 業界依據：cardsrealm.com「Hybrid Crash Game」趨勢（2026-05）
extends BaseLuckyPanel

const PANEL_COLOR = Color(1.0, 0.27, 0.0)   # 火橙紅色（崩潰）
const PANEL_ICON = "💥"
const PANEL_TITLE = "崩潰倍率"

var _mult_label: Label = null
var _status_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 107
	_setup_crash_ui()
	GameManager.lucky_crash_multiplier.connect(_on_lucky_crash_multiplier)

func _setup_crash_ui() -> void:
	_mult_label = Label.new()
	_mult_label.text = "當前倍率: ×1.0"
	_mult_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.0))
	_mult_label.add_theme_font_size_override("font_size", 28)
	_mult_label.position = Vector2(20, 52)
	add_child(_mult_label)

	_status_label = Label.new()
	_status_label.text = "倍率上升中... 隨時可崩潰！"
	_status_label.add_theme_color_override("font_color", Color(1.0, 0.8, 0.0))
	_status_label.add_theme_font_size_override("font_size", 16)
	_status_label.position = Vector2(20, 88)
	add_child(_status_label)

	_boost_label = Label.new()
	_boost_label.text = "完美收割 ≥×40.0 → 全服 ×36.5 加成 73 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 14)
	_boost_label.position = Vector2(20, 115)
	add_child(_boost_label)

func _on_lucky_crash_multiplier(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"crash_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開始！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(1.0, 0.27, 0.0))
		"mult_update":
			var current_mult = data.get("current_mult", 1.0)
			if is_instance_valid(_mult_label):
				_mult_label.text = "當前倍率: ×" + str(snapped(current_mult, 0.1))
			# 倍率越高顏色越紅
			if current_mult >= 40.0:
				if is_instance_valid(_mult_label):
					_mult_label.add_theme_color_override("font_color", Color(1.0, 0.0, 0.0))
			elif current_mult >= 20.0:
				if is_instance_valid(_mult_label):
					_mult_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.0))
		"crashed":
			var final_mult = data.get("final_mult", 1.0)
			show_settle(PANEL_ICON + " 崩潰！",
				"倍率在 ×" + str(snapped(final_mult, 0.1)) + " 時崩潰！",
				Color(0.8, 0.0, 0.0))
			hide_panel()
		"crash_perfect":
			var final_mult = data.get("final_mult", 1.0)
			var global_mult = data.get("global_mult", 36.5)
			var global_secs = data.get("global_secs", 73)
			show_settle(PANEL_ICON + " 完美收割！",
				"倍率 ×" + str(snapped(final_mult, 0.1)) + "！全服 ×" + str(global_mult) + " 加成 " + str(global_secs) + " 秒！",
				PANEL_COLOR)
			flash_screen(PANEL_COLOR)
			hide_panel()
		"crash_end":
			hide_panel()
