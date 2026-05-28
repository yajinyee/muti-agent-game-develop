## LuckySpacetimeFoldPanel.gd — T194 幸運時空折疊魚 UI
## lucky-panel-agent 負責維護
## DAY-317：時空折疊系統 — 20 秒目標倍率 ×3.0，射擊速度 ×2.0，觸發全服 ×21.0 加成 42 秒（新最高）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.42, 0.10, 0.60)  # 深紫色（時空折疊）
const PANEL_ICON = "🌀"
const PANEL_TITLE = "時空折疊"

var _timer_label: Label = null
var _effect_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 89
	_setup_spacetime_fold_ui()
	GameManager.lucky_spacetime_fold.connect(_on_lucky_spacetime_fold)

func _setup_spacetime_fold_ui() -> void:
	_timer_label = Label.new()
	_timer_label.text = "折疊時間: 20 秒"
	_timer_label.add_theme_color_override("font_color", Color(0.8, 0.5, 1.0))
	_timer_label.add_theme_font_size_override("font_size", 22)
	_timer_label.position = Vector2(20, 80)
	add_child(_timer_label)

	_effect_label = Label.new()
	_effect_label.text = "目標倍率 ×3.0 | 射擊速度 ×2.0"
	_effect_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	_effect_label.add_theme_font_size_override("font_size", 18)
	_effect_label.position = Vector2(20, 110)
	add_child(_effect_label)

	_boost_label = Label.new()
	_boost_label.text = "結束後全服 ×21.0 加成 42 秒（新最高）"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 15)
	_boost_label.position = Vector2(20, 140)
	add_child(_boost_label)

func _on_lucky_spacetime_fold(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"spacetime_fold_start":
			var trigger_name = data.get("trigger_name", "玩家")
			var target_mult = data.get("target_mult", 3.0)
			var fire_rate_mult = data.get("fire_rate_mult", 2.0)
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 啟動！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.5, 0.0, 0.8))
			if is_instance_valid(_effect_label):
				_effect_label.text = "目標倍率 ×%.1f | 射擊速度 ×%.1f" % [target_mult, fire_rate_mult]
		"spacetime_fold_complete":
			var boost_mult = data.get("boost_mult", 21.0)
			var boost_secs = data.get("boost_secs", 42)
			var trigger_name = data.get("trigger_name", "玩家")
			show_settle(PANEL_ICON + " 時空折疊！",
				"折疊結束！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！（新最高）",
				PANEL_COLOR)
			flash_screen(Color(0.7, 0.2, 1.0))
			hide_panel()
