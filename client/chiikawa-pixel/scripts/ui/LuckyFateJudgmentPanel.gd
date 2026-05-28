## LuckyFateJudgmentPanel.gd — T203 幸運命運審判魚 UI
## lucky-panel-agent 負責維護
## DAY-319：命運審判系統 — 隨機 5 個目標各 ×50-×500，全服 ×28.0 加成 56 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(1.0, 0.85, 0.0)  # 金色（命運審判）
const PANEL_ICON = "⚖️"
const PANEL_TITLE = "命運審判"

var _target_label: Label = null
var _mult_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 98
	_setup_fate_judgment_ui()
	GameManager.lucky_fate_judgment.connect(_on_lucky_fate_judgment)

func _setup_fate_judgment_ui() -> void:
	_target_label = Label.new()
	_target_label.text = "命運目標: 0 / 5"
	_target_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_target_label.add_theme_font_size_override("font_size", 22)
	_target_label.position = Vector2(20, 52)
	add_child(_target_label)

	_mult_label = Label.new()
	_mult_label.text = "倍率範圍：×50 ~ ×500"
	_mult_label.add_theme_color_override("font_color", Color(1.0, 0.6, 0.0))
	_mult_label.add_theme_font_size_override("font_size", 18)
	_mult_label.position = Vector2(20, 80)
	add_child(_mult_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×28.0 加成 56 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 108)
	add_child(_boost_label)

func _on_lucky_fate_judgment(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"fate_judgment_start":
			var target_count = data.get("target_count", 0)
			if is_instance_valid(_target_label):
				_target_label.text = "命運目標: " + str(target_count) + " 個"
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 觸發！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(1.0, 0.85, 0.0))
		"fate_judgment_complete":
			var boost_mult = data.get("boost_mult", 28.0)
			var boost_secs = data.get("boost_secs", 56)
			show_settle(PANEL_ICON + " 命運審判完成！",
				"全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！",
				PANEL_COLOR)
			flash_screen(Color(1.0, 0.85, 0.0))
			hide_panel()
