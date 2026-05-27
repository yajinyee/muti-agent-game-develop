## LuckyGravityFieldPanel.gd — T187 幸運引力場魚 UI
## lucky-panel-agent 負責維護
## DAY-316：引力場系統 — 引力場 15 秒（目標速度 ×0.1），引力爆炸（HP -55%，×9.0），全服 ×17.5 加成 37 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.19, 0.11, 0.57)  # 深紫色（引力）
const PANEL_ICON = "🌀"
const PANEL_TITLE = "引力場"

var _timer_label: Label = null
var _hit_label: Label = null
var _gravity_timer: float = 0.0
var _is_active: bool = false

func _ready() -> void:
	super._ready()
	layer = 82
	_setup_gravity_field_ui()
	GameManager.lucky_gravity_field.connect(_on_lucky_gravity_field)

func _setup_gravity_field_ui() -> void:
	_timer_label = Label.new()
	_timer_label.text = "引力場: 15 秒"
	_timer_label.add_theme_color_override("font_color", Color(0.6, 0.4, 1.0))
	_timer_label.add_theme_font_size_override("font_size", 22)
	_timer_label.position = Vector2(20, 80)
	add_child(_timer_label)

	_hit_label = Label.new()
	_hit_label.text = "目標速度 ×0.1 → 引力爆炸 HP -55%"
	_hit_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_hit_label.add_theme_font_size_override("font_size", 16)
	_hit_label.position = Vector2(20, 115)
	add_child(_hit_label)

func _process(delta: float) -> void:
	if _is_active and _gravity_timer > 0:
		_gravity_timer -= delta
		if is_instance_valid(_timer_label):
			_timer_label.text = "引力場: %.1f 秒" % max(0.0, _gravity_timer)

func _on_lucky_gravity_field(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"gravity_field_start":
			_gravity_timer = 15.0
			_is_active = true
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 啟動！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.3, 0.1, 0.8))
		"gravity_explosion":
			_is_active = false
			var hit_count = data.get("hit_count", 0)
			if is_instance_valid(_hit_label):
				_hit_label.text = "引力爆炸！命中 " + str(hit_count) + " 個！HP -55%！"
		"gravity_perfect":
			var boost_mult = data.get("boost_mult", 17.5)
			var boost_secs = data.get("boost_secs", 37)
			show_settle(PANEL_ICON + " 引力完美！", "命中 " + str(data.get("hit_count", 0)) + " 個！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", PANEL_COLOR)
			flash_screen(Color(0.4, 0.2, 1.0))
			hide_panel()
