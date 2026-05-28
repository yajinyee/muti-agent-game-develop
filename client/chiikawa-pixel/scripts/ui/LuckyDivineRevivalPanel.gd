## LuckyDivineRevivalPanel.gd — T199 幸運神聖復活魚 UI
## lucky-panel-agent 負責維護
## DAY-318：神聖復活系統 — 最近死亡的 5 個目標全部復活（HP 80%，獎勵 ×4.0），全服 ×24.5 加成 49 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.1, 0.1, 0.0)   # 深金色（神聖）
const PANEL_BORDER_COLOR = Color(1.0, 0.95, 0.4)  # 亮金色邊框
const PANEL_ICON = "✨"
const PANEL_TITLE = "神聖復活"

var _revive_label: Label = null
var _reward_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 94
	_setup_divine_revival_ui()
	GameManager.lucky_divine_revival.connect(_on_lucky_divine_revival)

func _setup_divine_revival_ui() -> void:
	_revive_label = Label.new()
	_revive_label.text = "復活目標: ?"
	_revive_label.add_theme_color_override("font_color", Color(1.0, 0.95, 0.4))
	_revive_label.add_theme_font_size_override("font_size", 22)
	_revive_label.position = Vector2(20, 55)
	add_child(_revive_label)

	_reward_label = Label.new()
	_reward_label.text = "強化獎勵: ×4.0"
	_reward_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	_reward_label.add_theme_font_size_override("font_size", 20)
	_reward_label.position = Vector2(20, 82)
	add_child(_reward_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×24.5 加成 49 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 0.95, 0.4))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 112)
	add_child(_boost_label)

func _on_lucky_divine_revival(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"divine_revival_start":
			var revive_count = data.get("revive_count", 0)
			if is_instance_valid(_revive_label):
				_revive_label.text = "復活目標: " + str(revive_count) + " 個"
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 降臨！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(1.0, 0.95, 0.4))
		"divine_revival_complete":
			var revived_count = data.get("revived_ids", []).size()
			var boost_mult = data.get("boost_mult", 24.5)
			var boost_secs = data.get("boost_secs", 49)
			show_settle(PANEL_ICON + " 神聖復活完成！",
				"復活 " + str(revived_count) + " 個強化目標！獎勵 ×4.0！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！",
				PANEL_COLOR)
			flash_screen(Color(1.0, 0.95, 0.4))
			hide_panel()
