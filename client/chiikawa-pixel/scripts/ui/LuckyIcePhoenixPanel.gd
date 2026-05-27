## LuckyIcePhoenixPanel.gd — T156 幸運冰鳳凰魚
## 業界依據：Royal Fishing「Ice Phoenix 180-300x」
## 冰凍全場 10 秒（傷害 ×1.5），鳳凰重生爆炸（HP -60%），命中 ≥8 → 完美鳳凰全服 ×5.5 加成 12 秒
extends BaseLuckyPanel

const PANEL_KEY = "lucky_ice_phoenix"
const THEME_COLOR = Color(0.0, 0.8, 1.0)  # 冰藍色
const ICON = "🧊🔥"

var _freeze_timer: float = 0.0
var _freeze_duration: float = 10.0
var _kill_count: int = 0
var _is_frozen: bool = false

func _ready() -> void:
	layer = 51
	_setup_panel(PANEL_KEY, THEME_COLOR, ICON, "冰鳳凰")

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"freeze_start":
			_freeze_duration = event_data.get("duration", 10.0)
			_freeze_timer = _freeze_duration
			_is_frozen = true
			_kill_count = 0
			show_banner("🧊 冰鳳凰降臨！全場凍結 %d 秒！傷害 ×1.5！" % int(_freeze_duration))
			show_indicator("凍結中", THEME_COLOR)
			start_timer(_freeze_duration)
			flash_screen(THEME_COLOR, 2)
		"phoenix_kill":
			_kill_count = event_data.get("kill_count", 0)
			update_indicator("凍結中 擊破:%d" % _kill_count)
		"phoenix_rebirth":
			_is_frozen = false
			var hit = event_data.get("hit_count", 0)
			show_banner("🔥 鳳凰重生！爆炸命中 %d 個！" % hit)
			flash_screen(Color(1.0, 0.5, 0.0), 3)
		"phoenix_perfect":
			var boost = event_data.get("boost_mult", 5.5)
			var secs = event_data.get("boost_secs", 12)
			show_banner("🔥✨ 完美鳳凰！全服 ×%.1f 加成 %d 秒！" % [boost, secs])
			show_indicator("完美鳳凰 ×%.1f" % boost, Color(1.0, 0.8, 0.0))
			flash_screen(Color(1.0, 0.8, 0.0), 4)
		"phoenix_perfect_end":
			hide_indicator()

func _process(delta: float) -> void:
	if _is_frozen and _freeze_timer > 0:
		_freeze_timer -= delta
		update_indicator("凍結 %.1fs 擊破:%d" % [_freeze_timer, _kill_count])
		if _freeze_timer <= 0:
			_is_frozen = false
