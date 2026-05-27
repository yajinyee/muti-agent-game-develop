## LuckyAwakenBossV2Panel.gd — T159 幸運覺醒 BOSS 魚 v2
## 業界依據：Royal Fishing「Awaken Boss Power Up 6x-10x」升級版
## 覺醒後 8 次 Power Up（每次 8x-15x 隨機），全部命中 → 完美覺醒全服 ×7.0 加成 15 秒
extends BaseLuckyPanel

const PANEL_KEY = "lucky_awaken_boss_v2"
const THEME_COLOR = Color(1.0, 0.6, 0.0)  # 金橙色
const ICON = "⚡"

var _shots_left: int = 8
var _hit_count: int = 0
var _awaken_timer: float = 0.0
var _is_active: bool = false

func _ready() -> void:
	layer = 54
	_setup_panel(PANEL_KEY, THEME_COLOR, ICON, "覺醒BOSS v2")

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"awaken_start":
			_shots_left = event_data.get("shots", 8)
			_hit_count = 0
			_awaken_timer = event_data.get("duration", 25.0)
			_is_active = true
			show_banner("⚡ 覺醒 BOSS！8 次 Power Up（8x-15x）！25 秒！")
			show_indicator("Power Up: 0/8", THEME_COLOR)
			start_timer(_awaken_timer)
			flash_screen(THEME_COLOR, 3)
		"power_up":
			_shots_left = event_data.get("shots_left", 0)
			_hit_count = event_data.get("hit_count", 0)
			var mult = event_data.get("power_up_mult", 8.0)
			var reward = event_data.get("reward", 0)
			update_indicator("Power Up: %d/8 (×%.1f)" % [_hit_count, mult])
			show_floating_text("⚡ ×%.1f！+%d" % [mult, reward], Color(1.0, 0.8, 0.0))
		"awaken_perfect":
			_is_active = false
			var boost = event_data.get("boost_mult", 7.0)
			var secs = event_data.get("boost_secs", 15)
			show_banner("⚡✨ 完美覺醒！8 次全命中！全服 ×%.1f 加成 %d 秒！" % [boost, secs])
			show_indicator("完美覺醒 ×%.1f" % boost, Color(1.0, 0.8, 0.0))
			flash_screen(Color(1.0, 0.8, 0.0), 5)
		"awaken_perfect_end":
			hide_indicator()
		"awaken_timeout":
			_is_active = false
			var hit = event_data.get("hit_count", 0)
			show_banner("⚡ 覺醒結束！命中 %d/8 次" % hit)
			hide_indicator()

func _process(delta: float) -> void:
	if _is_active and _awaken_timer > 0:
		_awaken_timer -= delta
		update_indicator("Power Up: %d/8 (%.1fs)" % [_hit_count, _awaken_timer])
