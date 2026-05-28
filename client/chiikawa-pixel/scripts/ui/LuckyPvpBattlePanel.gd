## LuckyPvpBattlePanel.gd — T191 幸運 PvP 競技魚 UI
## lucky-panel-agent 負責維護
## DAY-317：PvP 競技系統 — 全服競技 30 秒，勝者 ×20.0，全服 ×19.5 加成 40 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.72, 0.11, 0.11)  # 深紅色（PvP 競技）
const PANEL_ICON = "⚔️"
const PANEL_TITLE = "PvP 競技"

var _timer_label: Label = null
var _winner_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 86
	_setup_pvp_battle_ui()
	GameManager.lucky_pvp_battle.connect(_on_lucky_pvp_battle)

func _setup_pvp_battle_ui() -> void:
	_timer_label = Label.new()
	_timer_label.text = "競技時間: 30 秒"
	_timer_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	_timer_label.add_theme_font_size_override("font_size", 22)
	_timer_label.position = Vector2(20, 80)
	add_child(_timer_label)

	_winner_label = Label.new()
	_winner_label.text = "擊破最多目標者獲得 ×20.0"
	_winner_label.add_theme_color_override("font_color", Color(1.0, 0.4, 0.4))
	_winner_label.add_theme_font_size_override("font_size", 18)
	_winner_label.position = Vector2(20, 110)
	add_child(_winner_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×19.5 加成 40 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 140)
	add_child(_boost_label)

func _on_lucky_pvp_battle(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"pvp_battle_start":
			var trigger_name = data.get("trigger_name", "玩家")
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開始！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.8, 0.0, 0.0))
			if is_instance_valid(_timer_label):
				_timer_label.text = "競技時間: 30 秒"
			if is_instance_valid(_winner_label):
				_winner_label.text = trigger_name + " 發起競技！"
		"pvp_battle_progress":
			var time_left = data.get("time_left", 30)
			if is_instance_valid(_timer_label):
				_timer_label.text = "競技時間: %d 秒" % time_left
		"pvp_battle_complete":
			var winner_name = data.get("winner_name", "玩家")
			var winner_boost_mult = data.get("winner_boost_mult", 20.0)
			var winner_boost_secs = data.get("winner_boost_secs", 35)
			var global_boost_mult = data.get("global_boost_mult", 19.5)
			var global_boost_secs = data.get("global_boost_secs", 40)
			if is_instance_valid(_winner_label):
				_winner_label.text = "勝者: " + winner_name
			if is_instance_valid(_boost_label):
				_boost_label.text = "全服 ×" + str(global_boost_mult) + " 加成 " + str(global_boost_secs) + " 秒"
			show_settle(PANEL_ICON + " PvP 競技！",
				winner_name + " 獲勝！個人 ×" + str(winner_boost_mult) + " 加成 " + str(winner_boost_secs) + " 秒！全服 ×" + str(global_boost_mult) + " 加成 " + str(global_boost_secs) + " 秒！",
				PANEL_COLOR)
			flash_screen(Color(1.0, 0.1, 0.1))
			hide_panel()
