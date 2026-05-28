## LuckyTimeReversalPanel.gd — T204 幸運時間逆流魚 UI
## lucky-panel-agent 負責維護
## DAY-319：時間逆流系統 — 最近死亡的 10 個目標全部復活（HP 100%，獎勵 ×5.0），全服 ×29.0 加成 58 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.5, 0.5, 1.0)   # 深藍紫色（時間逆流）
const PANEL_ICON = "⏪"
const PANEL_TITLE = "時間逆流"

var _revive_label: Label = null
var _mult_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 99
	_setup_time_reversal_ui()
	GameManager.lucky_time_reversal.connect(_on_lucky_time_reversal)

func _setup_time_reversal_ui() -> void:
	_revive_label = Label.new()
	_revive_label.text = "復活目標: 0 個"
	_revive_label.add_theme_color_override("font_color", Color(0.5, 0.5, 1.0))
	_revive_label.add_theme_font_size_override("font_size", 22)
	_revive_label.position = Vector2(20, 52)
	add_child(_revive_label)

	_mult_label = Label.new()
	_mult_label.text = "復活目標獎勵 ×5.0"
	_mult_label.add_theme_color_override("font_color", Color(0.7, 0.7, 1.0))
	_mult_label.add_theme_font_size_override("font_size", 18)
	_mult_label.position = Vector2(20, 80)
	add_child(_mult_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×29.0 加成 58 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 108)
	add_child(_boost_label)

func _on_lucky_time_reversal(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"time_reversal_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 觸發！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.5, 0.5, 1.0))
		"time_reversal_revived":
			var revive_count = data.get("revive_count", 0)
			if is_instance_valid(_revive_label):
				_revive_label.text = "復活目標: " + str(revive_count) + " 個"
		"time_reversal_complete":
			var boost_mult = data.get("boost_mult", 29.0)
			var boost_secs = data.get("boost_secs", 58)
			show_settle(PANEL_ICON + " 時間逆流完成！",
				"復活 10 個目標！獎勵 ×5.0！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！",
				PANEL_COLOR)
			flash_screen(Color(0.5, 0.5, 1.0))
			hide_panel()
