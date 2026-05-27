## LuckyJackpotMiniPanel.gd — T171 幸運 Mini Jackpot 魚 UI
## progressive-jackpot-agent 負責維護
## DAY-313：Mini Jackpot 系統 — 擊破後直接觸發 Mini Jackpot（50x 起跳累積獎池）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.2, 0.8, 0.2)  # 綠色（Mini）
const PANEL_ICON = "🎰"
const PANEL_TITLE = "Mini Jackpot"

var _pool_label: Label = null
var _win_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 66
	_setup_mini_jackpot_ui()
	GameManager.lucky_jackpot_pool.connect(_on_lucky_jackpot_pool)

func _setup_mini_jackpot_ui() -> void:
	_pool_label = Label.new()
	_pool_label.text = "Mini Pool: 50x"
	_pool_label.add_theme_color_override("font_color", Color(0.2, 1.0, 0.2))
	_pool_label.add_theme_font_size_override("font_size", 20)
	_pool_label.position = Vector2(20, 80)
	add_child(_pool_label)

	_win_label = Label.new()
	_win_label.text = ""
	_win_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_win_label.add_theme_font_size_override("font_size", 24)
	_win_label.position = Vector2(20, 110)
	add_child(_win_label)

func _on_lucky_jackpot_pool(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"pool_update":
			var mini = data.get("mini", 50.0)
			if is_instance_valid(_pool_label):
				_pool_label.text = "Mini Pool: %.0fx" % mini
		"jackpot_win":
			var tier = data.get("tier", "")
			if tier != "mini":
				return
			var reward = data.get("reward", 0)
			var player_name = data.get("player_name", "玩家")
			show_banner(PANEL_ICON + " Mini Jackpot！" + player_name + " 獲得 " + str(reward) + " 金幣！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.2, 1.0, 0.2))
			if is_instance_valid(_win_label):
				_win_label.text = "🏆 +" + str(reward) + " 金幣！"
			show_settle("🎰 Mini Jackpot！", player_name + " 獲得 " + str(reward) + " 金幣！", PANEL_COLOR)
			hide_panel()
