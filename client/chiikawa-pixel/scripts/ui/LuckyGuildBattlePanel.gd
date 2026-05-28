## LuckyGuildBattlePanel.gd — T207 幸運公會戰魚 UI
## lucky-panel-agent 負責維護
## DAY-323：公會戰系統 — 全服公會戰 45 秒，擊破最多目標的玩家獲得 ×35.0，全服 ×32.0 加成 64 秒
## 業界依據：Fishing Frenzy Chapter 3「Guild Wars + Boss Fish」（2026-05-27）
extends BaseLuckyPanel

const PANEL_COLOR = Color(1.0, 0.85, 0.0)  # 金色（公會戰）
const PANEL_ICON = "⚔️"
const PANEL_TITLE = "公會戰"

var _battle_label: Label = null
var _timer_label: Label = null
var _kill_label: Label = null
var _global_label: Label = null
var _battle_timer: float = 0.0
var _battle_active: bool = false
var _my_kills: int = 0

func _ready() -> void:
	super._ready()
	layer = 102
	_setup_guild_battle_ui()
	GameManager.lucky_guild_battle.connect(_on_lucky_guild_battle)

func _setup_guild_battle_ui() -> void:
	_battle_label = Label.new()
	_battle_label.text = "⚔️ 公會戰開始！"
	_battle_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_battle_label.add_theme_font_size_override("font_size", 22)
	_battle_label.position = Vector2(20, 52)
	add_child(_battle_label)

	_timer_label = Label.new()
	_timer_label.text = "剩餘時間：45 秒"
	_timer_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	_timer_label.add_theme_font_size_override("font_size", 20)
	_timer_label.position = Vector2(20, 80)
	add_child(_timer_label)

	_kill_label = Label.new()
	_kill_label.text = "我的擊破：0"
	_kill_label.add_theme_color_override("font_color", Color(0.0, 1.0, 0.5))
	_kill_label.add_theme_font_size_override("font_size", 20)
	_kill_label.position = Vector2(20, 108)
	add_child(_kill_label)

	_global_label = Label.new()
	_global_label.text = "勝者獲得 ×35.0 | 全服 ×32.0 加成 64 秒"
	_global_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_global_label.add_theme_font_size_override("font_size", 14)
	_global_label.position = Vector2(20, 136)
	add_child(_global_label)

func _process(delta: float) -> void:
	if _battle_active and _battle_timer > 0.0:
		_battle_timer -= delta
		if is_instance_valid(_timer_label):
			_timer_label.text = "剩餘時間：%.1f 秒" % max(0.0, _battle_timer)
		if _battle_timer <= 0.0:
			_battle_active = false

func _on_lucky_guild_battle(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"battle_start":
			var battle_secs = data.get("battle_secs", 45)
			_battle_timer = float(battle_secs)
			_battle_active = true
			_my_kills = 0
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開始！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(1.0, 0.85, 0.0))
		"battle_complete":
			var winner_id = data.get("winner_id", "")
			var winner_kills = data.get("winner_kills", 0)
			var winner_mult = data.get("winner_mult", 35.0)
			var global_mult = data.get("global_mult", 32.0)
			var global_secs = data.get("global_secs", 64)
			_battle_active = false
			show_settle(PANEL_ICON + " 公會戰結束！",
				"勝者 %s（%d 擊破）×%.1f！全服 ×%.1f 加成 %d 秒！" % [winner_id, winner_kills, winner_mult, global_mult, global_secs],
				PANEL_COLOR)
			flash_screen(Color(1.0, 1.0, 0.0))
			hide_panel()
