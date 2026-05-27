## LuckyFinalPowerPanel.gd — T180 幸運終焉之力魚 UI
## lucky-panel-agent 負責維護
## DAY-314：終焉之力系統 — 全場 HP 歸零（每個獎勵 ×10.0），觸發全服 ×15.0 加成 30 秒（新最高倍率）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.72, 0.07, 0.07)  # 深紅色（終焉）
const PANEL_ICON = "💀"
const PANEL_TITLE = "終焉之力"

var _hit_count: int = 0
var _hit_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 75
	_setup_final_power_ui()
	GameManager.lucky_final_power.connect(_on_lucky_final_power)

func _setup_final_power_ui() -> void:
	# 清場計數
	_hit_label = Label.new()
	_hit_label.text = "清場: 0 個"
	_hit_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))
	_hit_label.add_theme_font_size_override("font_size", 22)
	_hit_label.position = Vector2(20, 80)
	add_child(_hit_label)

	# 全服加成顯示（超越 T170）
	_boost_label = Label.new()
	_boost_label.text = "全服 ×15.0 加成 30 秒（新最高）"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 115)
	add_child(_boost_label)

func _on_lucky_final_power(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"final_power_start":
			_hit_count = 0
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 引動！", PANEL_COLOR)
			show_panel()
			# 強烈全螢幕閃光（比 T170 更強）
			flash_screen(Color(1.0, 0.0, 0.0))
		"final_power_complete":
			_hit_count = data.get("hit_count", 0)
			var boost_mult = data.get("boost_mult", 15.0)
			var boost_secs = data.get("boost_secs", 30)
			if is_instance_valid(_hit_label):
				_hit_label.text = "清場: " + str(_hit_count) + " 個"
			show_settle(PANEL_ICON + " 終焉之力！", "清場 " + str(_hit_count) + " 個！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", PANEL_COLOR)
			flash_screen(Color(0.8, 0.0, 0.0))
			hide_panel()
