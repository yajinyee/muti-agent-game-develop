## LuckyTimeAccelerationPanel.gd — T188 幸運時間加速魚 UI
## lucky-panel-agent 負責維護
## DAY-316：時間加速系統 — 目標速度 ×0.15，射擊速度 ×3.0，獎勵 ×2.5，擊破 ≥20 → 全服 ×18.0 加成 38 秒（新最高）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.90, 0.40, 0.05)  # 火橙色（時間加速）
const PANEL_ICON = "⚡"
const PANEL_TITLE = "時間加速"

var _timer_label: Label = null
var _kill_label: Label = null
var _accel_timer: float = 0.0
var _kill_count: int = 0
var _is_active: bool = false

func _ready() -> void:
	super._ready()
	layer = 83
	_setup_time_acceleration_ui()
	GameManager.lucky_time_acceleration.connect(_on_lucky_time_acceleration)

func _setup_time_acceleration_ui() -> void:
	_timer_label = Label.new()
	_timer_label.text = "加速: 30 秒"
	_timer_label.add_theme_color_override("font_color", Color(1.0, 0.7, 0.2))
	_timer_label.add_theme_font_size_override("font_size", 22)
	_timer_label.position = Vector2(20, 80)
	add_child(_timer_label)

	_kill_label = Label.new()
	_kill_label.text = "擊破: 0 / 20 個（射擊 ×3.0）"
	_kill_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.5))
	_kill_label.add_theme_font_size_override("font_size", 18)
	_kill_label.position = Vector2(20, 115)
	add_child(_kill_label)

func _process(delta: float) -> void:
	if _is_active and _accel_timer > 0:
		_accel_timer -= delta
		if is_instance_valid(_timer_label):
			_timer_label.text = "加速: %.1f 秒" % max(0.0, _accel_timer)

func _on_lucky_time_acceleration(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"time_acceleration_start":
			_accel_timer = 30.0
			_kill_count = 0
			_is_active = true
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 啟動！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(1.0, 0.5, 0.0))
		"time_acceleration_perfect":
			_is_active = false
			var kills = data.get("kill_count", 0)
			var boost_mult = data.get("boost_mult", 18.0)
			var boost_secs = data.get("boost_secs", 38)
			show_settle(PANEL_ICON + " 時間完美！", "擊破 " + str(kills) + " 個！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！（新最高）", PANEL_COLOR)
			flash_screen(Color(1.0, 0.7, 0.0))
			hide_panel()
		"time_acceleration_end":
			_is_active = false
			var kills = data.get("kill_count", 0)
			if is_instance_valid(_kill_label):
				_kill_label.text = "擊破: " + str(kills) + " 個（需 20 個觸發完美）"
			hide_panel()
