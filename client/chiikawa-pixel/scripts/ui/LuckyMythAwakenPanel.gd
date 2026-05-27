## LuckyMythAwakenPanel.gd — T165 幸運神話覺醒魚 UI
## lucky-panel-agent 負責維護
## DAY-310：神話覺醒系統 — 全場倍率 ×3.0，25 秒，擊破 ≥15 → 神話完美全服 ×8.0 加成 20 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.9, 0.8, 0.1)  # 神聖金色
const PANEL_ICON = "🌟"
const PANEL_TITLE = "神話覺醒"

var _kill_count: int = 0
var _target_kills: int = 15
var _myth_mult: float = 3.0
var _kill_label: Label = null
var _myth_label: Label = null
var _kill_bar: ProgressBar = null

func _ready() -> void:
	super._ready()
	layer = 60
	_setup_myth_ui()
	GameManager.lucky_myth_awaken.connect(_on_lucky_myth_awaken)

func _setup_myth_ui() -> void:
	# 神話倍率顯示
	_myth_label = Label.new()
	_myth_label.text = "全場倍率 ×3.0"
	_myth_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	_myth_label.add_theme_font_size_override("font_size", 20)
	_myth_label.position = Vector2(20, 80)
	add_child(_myth_label)

	# 擊破計數
	_kill_label = Label.new()
	_kill_label.text = "擊破: 0/15"
	_kill_label.add_theme_color_override("font_color", Color(1.0, 0.8, 0.3))
	_kill_label.add_theme_font_size_override("font_size", 18)
	_kill_label.position = Vector2(20, 110)
	add_child(_kill_label)

	# 擊破進度條
	_kill_bar = ProgressBar.new()
	_kill_bar.min_value = 0
	_kill_bar.max_value = _target_kills
	_kill_bar.value = 0
	_kill_bar.size = Vector2(200, 16)
	_kill_bar.position = Vector2(20, 140)
	add_child(_kill_bar)

func _on_lucky_myth_awaken(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"start":
			_kill_count = 0
			_myth_mult = data.get("myth_mult", 3.0)
			_update_myth_display()
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 降臨！", PANEL_COLOR)
			start_timer(data.get("duration", 25))
			show_panel()
			# 全場金色閃光
			flash_screen(PANEL_COLOR)
		"end":
			var perfect = data.get("perfect", false)
			var final_kills = data.get("final_kills", 0)
			if perfect:
				show_settle("🌟 神話完美！", "擊破 " + str(final_kills) + " 個！全服 ×8.0 加成 20 秒！", PANEL_COLOR)
				flash_screen(Color(1.0, 0.95, 0.5))
			else:
				show_settle("神話覺醒結束", "擊破 " + str(final_kills) + " 個目標", PANEL_COLOR)
			hide_panel()

func _update_myth_display() -> void:
	if is_instance_valid(_kill_label):
		_kill_label.text = "擊破: " + str(_kill_count) + "/" + str(_target_kills)
	if is_instance_valid(_kill_bar):
		_kill_bar.value = min(_kill_count, _target_kills)
