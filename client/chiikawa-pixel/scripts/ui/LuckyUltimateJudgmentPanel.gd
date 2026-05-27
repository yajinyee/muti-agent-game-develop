## LuckyUltimateJudgmentPanel.gd — T160 幸運終極審判魚
## 終極機制：全場目標 HP 歸零（每個獎勵 ×6.0），觸發全服 ×10.0 加成 20 秒
extends BaseLuckyPanel

const PANEL_KEY = "lucky_ultimate_judgment"
const THEME_COLOR = Color(0.8, 0.0, 0.0)  # 深紅色
const ICON = "⚖️"

var _boost_timer: float = 0.0
var _is_boosting: bool = false

func _ready() -> void:
	layer = 55
	_setup_panel(PANEL_KEY, THEME_COLOR, ICON, "終極審判")

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"judgment_start":
			show_banner("⚖️💥 終極審判！全場 HP 歸零！")
			flash_screen(Color(0.8, 0.0, 0.0), 3)
		"judgment_execute":
			var hit = event_data.get("hit_count", 0)
			var mult = event_data.get("reward_mult", 6.0)
			show_banner("⚖️ 審判執行！清場 %d 個！每個獎勵 ×%.1f！" % [hit, mult])
			flash_screen(Color(1.0, 0.5, 0.0), 4)
		"judgment_boost":
			var boost = event_data.get("boost_mult", 10.0)
			var secs = event_data.get("boost_secs", 20)
			var hit = event_data.get("hit_count", 0)
			_boost_timer = float(secs)
			_is_boosting = true
			show_banner("⚖️✨ 終極審判完成！清場 %d 個！全服 ×%.1f 加成 %d 秒！" % [hit, boost, secs])
			show_indicator("終極審判 ×%.1f" % boost, Color(1.0, 0.8, 0.0))
			flash_screen(Color(1.0, 0.8, 0.0), 6)
		"judgment_boost_end":
			_is_boosting = false
			hide_indicator()
			show_banner("⚖️ 終極審判加成結束")

func _process(delta: float) -> void:
	if _is_boosting and _boost_timer > 0:
		_boost_timer -= delta
		update_indicator("終極審判 ×10.0 (%.1fs)" % _boost_timer)
		if _boost_timer <= 0:
			_is_boosting = false
