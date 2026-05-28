## LuckySkillChainPanel.gd — T192 幸運技能連鎖魚 UI
## lucky-panel-agent 負責維護
## DAY-317：技能連鎖系統 — Lv.1-10，Lv.10 → 全服 ×20.0 加成 38 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.05, 0.27, 0.60)  # 深藍色（技能連鎖）
const PANEL_ICON = "🔗"
const PANEL_TITLE = "技能連鎖"

var _level_label: Label = null
var _mult_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 87
	_setup_skill_chain_ui()
	GameManager.lucky_skill_chain.connect(_on_lucky_skill_chain)

func _setup_skill_chain_ui() -> void:
	_level_label = Label.new()
	_level_label.text = "技能等級: Lv.1"
	_level_label.add_theme_color_override("font_color", Color(0.4, 0.8, 1.0))
	_level_label.add_theme_font_size_override("font_size", 24)
	_level_label.position = Vector2(20, 80)
	add_child(_level_label)

	_mult_label = Label.new()
	_mult_label.text = "當前倍率: ×2.0"
	_mult_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	_mult_label.add_theme_font_size_override("font_size", 20)
	_mult_label.position = Vector2(20, 110)
	add_child(_mult_label)

	_boost_label = Label.new()
	_boost_label.text = "Lv.10 → 全服 ×20.0 加成 38 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 15)
	_boost_label.position = Vector2(20, 140)
	add_child(_boost_label)

func _on_lucky_skill_chain(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"skill_chain_start":
			var trigger_name = data.get("trigger_name", "玩家")
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 啟動！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.0, 0.3, 0.8))
			if is_instance_valid(_level_label):
				_level_label.text = "技能等級: Lv.1"
			if is_instance_valid(_mult_label):
				_mult_label.text = "當前倍率: ×2.0"
		"skill_chain_level_up":
			var level = data.get("level", 1)
			var level_mult = data.get("level_mult", 2.0)
			if is_instance_valid(_level_label):
				_level_label.text = "技能等級: Lv.%d" % level
				if level >= 10:
					_level_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
				elif level >= 5:
					_level_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.0))
			if is_instance_valid(_mult_label):
				_mult_label.text = "當前倍率: ×%.1f" % level_mult
		"skill_chain_complete":
			var final_level = data.get("final_level", 1)
			var final_mult = data.get("final_mult", 2.0)
			var is_perfect = data.get("is_perfect", false)
			var global_boost_mult = data.get("global_boost_mult", 19.5)
			var global_boost_secs = data.get("global_boost_secs", 35)
			var trigger_name = data.get("trigger_name", "玩家")
			var settle_msg: String
			if is_perfect:
				settle_msg = "完美連鎖！Lv.10！×%.1f！全服 ×%.1f 加成 %d 秒！（新最高）" % [final_mult, global_boost_mult, global_boost_secs]
				flash_screen(Color(1.0, 0.85, 0.0))
			else:
				settle_msg = "Lv.%d（×%.1f）！全服 ×%.1f 加成 %d 秒！" % [final_level, final_mult, global_boost_mult, global_boost_secs]
				flash_screen(Color(0.2, 0.5, 1.0))
			show_settle(PANEL_ICON + " 技能連鎖！", settle_msg, PANEL_COLOR)
			hide_panel()
