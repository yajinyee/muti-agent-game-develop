## LuckyDragonFuryPanel.gd — T157 幸運龍怒能量魚
## 業界依據：Royal Fishing「Dragon Fury energy accumulation → full-screen attack」
## 能量累積 15 秒（每次擊破 +10），滿 100 → 龍怒全場（HP -80%），命中 ≥10 → 完美龍怒全服 ×6.0 加成 13 秒
extends BaseLuckyPanel

const PANEL_KEY = "lucky_dragon_fury"
const THEME_COLOR = Color(1.0, 0.3, 0.0)  # 火橙色
const ICON = "🐉"

var _energy: int = 0
var _energy_timer: float = 0.0
var _is_active: bool = false

func _ready() -> void:
	layer = 52
	_setup_panel(PANEL_KEY, THEME_COLOR, ICON, "龍怒能量")

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"energy_start":
			_energy = 0
			_energy_timer = event_data.get("duration", 15.0)
			_is_active = true
			show_banner("🐉 龍怒蓄能！15 秒內累積 100 能量！")
			show_indicator("能量: 0/100", THEME_COLOR)
			start_timer(_energy_timer)
			flash_screen(THEME_COLOR, 2)
		"energy_gain":
			_energy = event_data.get("energy", 0)
			update_indicator("能量: %d/100" % _energy)
		"fury_unleash":
			_is_active = false
			show_banner("🐉💥 龍怒爆發！全場 HP -80%%！")
			flash_screen(Color(1.0, 0.1, 0.0), 4)
		"fury_hit":
			var hit = event_data.get("hit_count", 0)
			show_banner("🐉 龍怒命中 %d 個！" % hit)
		"fury_perfect":
			var boost = event_data.get("boost_mult", 6.0)
			var secs = event_data.get("boost_secs", 13)
			show_banner("🐉✨ 完美龍怒！全服 ×%.1f 加成 %d 秒！" % [boost, secs])
			show_indicator("完美龍怒 ×%.1f" % boost, Color(1.0, 0.8, 0.0))
			flash_screen(Color(1.0, 0.8, 0.0), 5)
		"fury_perfect_end":
			hide_indicator()
		"energy_timeout":
			_is_active = false
			var energy = event_data.get("energy", 0)
			show_banner("🐉 能量不足（%d/100），龍怒未觸發" % energy)
			hide_indicator()

func _process(delta: float) -> void:
	if _is_active and _energy_timer > 0:
		_energy_timer -= delta
		update_indicator("能量: %d/100 (%.1fs)" % [_energy, _energy_timer])
