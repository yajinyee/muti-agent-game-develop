## LuckyPathFishPanel.gd — T208 幸運路徑魚 UI
## lucky-panel-agent 負責維護
## DAY-323：路徑系統 — 路徑越遠倍率越高（最高 ×20,000），全服 ×33.0 加成 66 秒
## 業界依據：Fish Road「路徑越遠倍率越高，最高 20,000x」（fishroad.eu）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.0, 1.0, 1.0)   # 青藍色（路徑）
const PANEL_ICON = "🛤️"
const PANEL_TITLE = "路徑魚"

var _step_label: Label = null
var _mult_label: Label = null
var _timer_label: Label = null
var _global_label: Label = null
var _path_timer: float = 0.0
var _path_active: bool = false
var _current_step: int = 0
var _current_mult: float = 1.5

func _ready() -> void:
	super._ready()
	layer = 103
	_setup_path_fish_ui()
	GameManager.lucky_path_fish.connect(_on_lucky_path_fish)

func _setup_path_fish_ui() -> void:
	_step_label = Label.new()
	_step_label.text = "🛤️ 路徑：第 1 步"
	_step_label.add_theme_color_override("font_color", Color(0.0, 1.0, 1.0))
	_step_label.add_theme_font_size_override("font_size", 22)
	_step_label.position = Vector2(20, 52)
	add_child(_step_label)

	_mult_label = Label.new()
	_mult_label.text = "當前倍率：×1.5"
	_mult_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_mult_label.add_theme_font_size_override("font_size", 24)
	_mult_label.position = Vector2(20, 80)
	add_child(_mult_label)

	_timer_label = Label.new()
	_timer_label.text = "剩餘時間：40 秒"
	_timer_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	_timer_label.add_theme_font_size_override("font_size", 18)
	_timer_label.position = Vector2(20, 110)
	add_child(_timer_label)

	_global_label = Label.new()
	_global_label.text = "最高 ×20,000 | 全服 ×33.0 加成 66 秒"
	_global_label.add_theme_color_override("font_color", Color(0.5, 1.0, 1.0))
	_global_label.add_theme_font_size_override("font_size", 14)
	_global_label.position = Vector2(20, 136)
	add_child(_global_label)

func _process(delta: float) -> void:
	if _path_active and _path_timer > 0.0:
		_path_timer -= delta
		if is_instance_valid(_timer_label):
			_timer_label.text = "剩餘時間：%.1f 秒" % max(0.0, _path_timer)
		if _path_timer <= 0.0:
			_path_active = false

func _on_lucky_path_fish(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"path_start":
			var path_secs = data.get("path_secs", 40)
			_path_timer = float(path_secs)
			_path_active = true
			_current_step = 0
			_current_mult = 1.5
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開啟！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.0, 1.0, 1.0))
		"path_advance":
			var step = data.get("step", 0)
			var mult = data.get("mult", 1.5)
			_current_step = step
			_current_mult = mult
			if is_instance_valid(_step_label):
				_step_label.text = "🛤️ 路徑：第 %d 步" % (step + 1)
			if is_instance_valid(_mult_label):
				_mult_label.text = "當前倍率：×%.1f" % mult
			# 高倍率時閃光
			if mult >= 1000.0:
				flash_screen(Color(1.0, 1.0, 0.0))
			elif mult >= 100.0:
				flash_screen(Color(0.0, 1.0, 1.0))
		"path_complete":
			var final_step = data.get("final_step", 0)
			var final_mult = data.get("final_mult", 1.5)
			var global_mult = data.get("global_mult", 33.0)
			var global_secs = data.get("global_secs", 66)
			_path_active = false
			show_settle(PANEL_ICON + " 路徑完成！",
				"第 %d 步（×%.1f）！全服 ×%.1f 加成 %d 秒！" % [final_step + 1, final_mult, global_mult, global_secs],
				PANEL_COLOR)
			flash_screen(Color(1.0, 1.0, 1.0))
			hide_panel()
