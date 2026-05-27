## LuckyMultCascadePanel.gd — T158 幸運倍率瀑布魚
## 業界依據：Fishing Fortune「Multiplier Cascade 2x→500x」
## 30 秒倍率瀑布（每次擊破 +0.5x，最高 ×20.0），達到 ×15.0 → 完美瀑布全服 ×6.5 加成 14 秒
extends BaseLuckyPanel

const PANEL_KEY = "lucky_mult_cascade"
const THEME_COLOR = Color(0.1, 0.4, 1.0)  # 深藍色
const ICON = "🌊"

var _current_mult: float = 1.0
var _cascade_timer: float = 0.0
var _is_active: bool = false

func _ready() -> void:
	layer = 53
	_setup_panel(PANEL_KEY, THEME_COLOR, ICON, "倍率瀑布")

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"cascade_start":
			_current_mult = 1.0
			_cascade_timer = event_data.get("duration", 30.0)
			_is_active = true
			show_banner("🌊 倍率瀑布！30 秒每次擊破 +0.5x！最高 ×20.0！")
			show_indicator("倍率: ×1.0", THEME_COLOR)
			start_timer(_cascade_timer)
			flash_screen(THEME_COLOR, 2)
		"cascade_rise":
			_current_mult = event_data.get("current_mult", 1.0)
			var kill_count = event_data.get("kill_count", 0)
			update_indicator("倍率: ×%.1f (擊破:%d)" % [_current_mult, kill_count])
		"cascade_perfect":
			var peak = event_data.get("peak_mult", 15.0)
			var boost = event_data.get("boost_mult", 6.5)
			var secs = event_data.get("boost_secs", 14)
			_is_active = false
			show_banner("🌊✨ 完美瀑布！倍率達 ×%.1f！全服 ×%.1f 加成 %d 秒！" % [peak, boost, secs])
			show_indicator("完美瀑布 ×%.1f" % boost, Color(1.0, 0.8, 0.0))
			flash_screen(Color(1.0, 0.8, 0.0), 4)
		"cascade_perfect_end":
			hide_indicator()
		"cascade_end":
			_is_active = false
			var final_mult = event_data.get("final_mult", 1.0)
			show_banner("🌊 瀑布結束！最終倍率 ×%.1f" % final_mult)
			hide_indicator()

func _process(delta: float) -> void:
	if _is_active and _cascade_timer > 0:
		_cascade_timer -= delta
		update_indicator("倍率: ×%.1f (%.1fs)" % [_current_mult, _cascade_timer])
