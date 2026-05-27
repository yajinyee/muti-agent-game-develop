## LuckyJackpotMajorPanel.gd — T173 幸運 Major Jackpot 魚 UI
## progressive-jackpot-agent 負責維護
## DAY-313：Major Jackpot 系統 — 擊破後直接觸發 Major Jackpot（1000x 起跳累積獎池）
extends BaseLuckyPanel

const PANEL_COLOR = Color(1.0, 0.6, 0.0)  # 橙色（Major）
const PANEL_ICON = "🎰🔥"
const PANEL_TITLE = "Major Jackpot"

var _pool_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 68
	_setup_major_jackpot_ui()
	GameManager.lucky_jackpot_pool.connect(_on_lucky_jackpot_pool)

func _setup_major_jackpot_ui() -> void:
	_pool_label = Label.new()
	_pool_label.text = "Major Pool: 1000x"
	_pool_label.add_theme_color_override("font_color", Color(1.0, 0.6, 0.0))
	_pool_label.add_theme_font_size_override("font_size", 20)
	_pool_label.position = Vector2(20, 80)
	add_child(_pool_label)

func _on_lucky_jackpot_pool(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"pool_update":
			var major = data.get("major", 1000.0)
			if is_instance_valid(_pool_label):
				_pool_label.text = "Major Pool: %.0fx" % major
		"jackpot_win":
			var tier = data.get("tier", "")
			if tier != "major":
				return
			var reward = data.get("reward", 0)
			var player_name = data.get("player_name", "玩家")
			show_banner(PANEL_ICON + " MAJOR JACKPOT！" + player_name + " 獲得 " + str(reward) + " 金幣！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(1.0, 0.6, 0.0))
			ScreenShake.add_trauma(0.6)
			show_settle("🎰🔥 Major Jackpot！", player_name + " 獲得 " + str(reward) + " 金幣！", PANEL_COLOR)
			hide_panel()
