## LuckyHolyJudgmentPanel.gd — T169 幸運神聖審判魚 UI
## lucky-panel-agent 負責維護
## DAY-312：神聖審判系統 — 25 秒每 5 秒一波神聖光柱（全場 HP -30%），5 波全部命中 ≥5 → 神聖完美全服 ×8.5 加成 18 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.96, 0.50, 0.09)  # 神聖橙金色
const PANEL_ICON = "✨"
const PANEL_TITLE = "神聖審判"

var _current_wave: int = 0
var _total_waves: int = 5
var _wave_labels: Array = []
var _wave_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 64
	_setup_holy_ui()
	GameManager.lucky_holy_judgment.connect(_on_lucky_holy_judgment)

func _setup_holy_ui() -> void:
	# 波次顯示
	_wave_label = Label.new()
	_wave_label.text = "波次: 0/5"
	_wave_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.3))
	_wave_label.add_theme_font_size_override("font_size", 20)
	_wave_label.position = Vector2(20, 80)
	add_child(_wave_label)

	# 5 個波次指示點
	for i in range(5):
		var dot = ColorRect.new()
		dot.size = Vector2(24, 24)
		dot.position = Vector2(20 + i * 36, 110)
		dot.color = Color(0.3, 0.3, 0.3)
		add_child(dot)
		_wave_labels.append(dot)

func _on_lucky_holy_judgment(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"judgment_start":
			_current_wave = 0
			_reset_wave_dots()
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 降臨！", PANEL_COLOR)
			start_timer(data.get("duration", 25))
			show_panel()
			flash_screen(PANEL_COLOR)
		"judgment_wave":
			_current_wave = data.get("wave", 0)
			var hit_count = data.get("hit_count", 0)
			_update_wave_dot(_current_wave - 1, hit_count >= 5)
			if is_instance_valid(_wave_label):
				_wave_label.text = "波次: " + str(_current_wave) + "/5"
		"judgment_perfect":
			var boost_mult = data.get("boost_mult", 8.5)
			var boost_secs = data.get("boost_secs", 18)
			show_settle("✨ 神聖完美！", "5 波全部命中！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", PANEL_COLOR)
			flash_screen(Color(1.0, 0.9, 0.4))
			hide_panel()
		"judgment_end":
			show_settle("神聖審判結束", "完成 " + str(_current_wave) + " 波審判", PANEL_COLOR)
			hide_panel()

func _reset_wave_dots() -> void:
	for dot in _wave_labels:
		if is_instance_valid(dot):
			dot.color = Color(0.3, 0.3, 0.3)

func _update_wave_dot(index: int, success: bool) -> void:
	if index >= 0 and index < _wave_labels.size():
		var dot = _wave_labels[index]
		if is_instance_valid(dot):
			dot.color = Color(0.2, 0.9, 0.2) if success else Color(0.9, 0.2, 0.2)
