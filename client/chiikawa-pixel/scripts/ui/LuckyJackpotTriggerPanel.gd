## LuckyJackpotTriggerPanel.gd — T175 幸運 Jackpot Trigger 魚 UI
## progressive-jackpot-agent 負責維護
## DAY-313：Jackpot Trigger 系統 — 擊破後隨機觸發四層之一
## 機率分佈：Mini 60% / Minor 30% / Major 8% / Grand 2%
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.9, 0.7, 0.1)  # 金黃色（觸發器）
const PANEL_ICON = "🎰✨"
const PANEL_TITLE = "Jackpot Trigger"

# 四層顏色
const TIER_COLORS = {
	"mini":  Color(0.2, 0.8, 0.2),
	"minor": Color(0.13, 0.59, 0.95),
	"major": Color(1.0, 0.6, 0.0),
	"grand": Color(1.0, 0.85, 0.0),
}

# 四層機率顯示
var _prob_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 70
	_setup_trigger_ui()
	GameManager.lucky_jackpot_pool.connect(_on_lucky_jackpot_pool)

func _setup_trigger_ui() -> void:
	_prob_label = Label.new()
	_prob_label.text = "Mini 60% | Minor 30% | Major 8% | Grand 2%"
	_prob_label.add_theme_color_override("font_color", Color(0.9, 0.9, 0.9))
	_prob_label.add_theme_font_size_override("font_size", 14)
	_prob_label.position = Vector2(10, 80)
	add_child(_prob_label)

func _on_lucky_jackpot_pool(data: Dictionary) -> void:
	var event = data.get("event", "")
	if event != "jackpot_win":
		return

	var tier = data.get("tier", "")
	var tier_name = data.get("tier_name", "Jackpot")
	var reward = data.get("reward", 0)
	var player_name = data.get("player_name", "玩家")
	var target_id = data.get("target_id", "")

	# 只有 T175 觸發的才顯示此 Panel
	if target_id != "T175":
		return

	var color = TIER_COLORS.get(tier, PANEL_COLOR)
	show_banner(PANEL_ICON + " " + tier_name + "！" + player_name + " 獲得 " + str(reward) + " 金幣！", color)
	show_panel()
	flash_screen(color)
	show_settle(PANEL_ICON + " " + tier_name + "！", player_name + " 獲得 " + str(reward) + " 金幣！", color)
	hide_panel()
