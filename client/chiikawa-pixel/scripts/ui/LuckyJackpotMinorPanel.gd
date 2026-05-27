## LuckyJackpotMinorPanel.gd — T172 幸運 Minor Jackpot 魚 UI
## progressive-jackpot-agent 負責維護
## DAY-313：Minor Jackpot 系統 — 擊破後直接觸發 Minor Jackpot（200x 起跳累積獎池）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.13, 0.59, 0.95)  # 藍色（Minor）
const PANEL_ICON = "🎰💙"
const PANEL_TITLE = "Minor Jackpot"

var _pool_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 67
	_setup_minor_jackpot_ui()
	GameManager.lucky_jackpot_pool.connect(_on_lucky_jackpot_pool)

func _setup_minor_jackpot_ui() -> void:
	_pool_label = Label.new()
	_pool_label.text = "Minor Pool: 200x"
	_pool_label.add_theme_color_override("font_color", Color(0.13, 0.59, 0.95))
	_pool_label.add_theme_font_size_override("font_size", 20)
	_pool_label.position = Vector2(20, 80)
	add_child(_pool_label)

func _on_lucky_jackpot_pool(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"pool_update":
			var minor = data.get("minor", 200.0)
			if is_instance_valid(_pool_label):
				_pool_label.text = "Minor Pool: %.0fx" % minor
		"jackpot_win":
			var tier = data.get("tier", "")
			if tier != "minor":
				return
			var reward = data.get("reward", 0)
			var player_name = data.get("player_name", "玩家")
			show_banner(PANEL_ICON + " Minor Jackpot！" + player_name + " 獲得 " + str(reward) + " 金幣！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.13, 0.59, 0.95))
			show_settle("🎰💙 Minor Jackpot！", player_name + " 獲得 " + str(reward) + " 金幣！", PANEL_COLOR)
			hide_panel()
