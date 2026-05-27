## LuckyFateWheelPanel.gd — T178 幸運命運之輪魚 UI
## lucky-panel-agent 負責維護
## DAY-314：命運之輪系統 — 3 次旋轉（最高 ×50.0），連續 3 次 ≥20x → 全服 ×11.0 加成 24 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.97, 0.50, 0.09)  # 火橙色（命運）
const PANEL_ICON = "🎡"
const PANEL_TITLE = "命運之輪"

var _spin_labels: Array = []
var _total_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 73
	_setup_fate_wheel_ui()
	GameManager.lucky_fate_wheel.connect(_on_lucky_fate_wheel)

func _setup_fate_wheel_ui() -> void:
	# 三次旋轉結果
	for i in range(3):
		var lbl = Label.new()
		lbl.text = "第 %d 轉: -" % (i + 1)
		lbl.add_theme_color_override("font_color", Color(1.0, 0.7, 0.3))
		lbl.add_theme_font_size_override("font_size", 20)
		lbl.position = Vector2(20, 80 + i * 30)
		add_child(lbl)
		_spin_labels.append(lbl)

	# 總獎勵
	_total_label = Label.new()
	_total_label.text = "3 次 ≥20x → 全服 ×11.0"
	_total_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_total_label.add_theme_font_size_override("font_size", 16)
	_total_label.position = Vector2(20, 175)
	add_child(_total_label)

func _on_lucky_fate_wheel(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"fate_wheel_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 旋轉！", PANEL_COLOR)
			show_panel()
			for lbl in _spin_labels:
				if is_instance_valid(lbl):
					lbl.text = "第 %d 轉: -" % (_spin_labels.find(lbl) + 1)
		"wheel_spin":
			var spin_no = data.get("spin_no", 1)
			var mult = data.get("mult", 1.0)
			var reward = data.get("reward", 0)
			var idx = spin_no - 1
			if idx < _spin_labels.size() and is_instance_valid(_spin_labels[idx]):
				_spin_labels[idx].text = "第 %d 轉: ×%.0f (+%d)" % [spin_no, mult, reward]
				_spin_labels[idx].modulate = Color(1.0, 0.85, 0.0) if mult >= 20 else Color(1.0, 0.7, 0.3)
			show_banner(PANEL_ICON + " 第 %d 轉 ×%.0f！" % [spin_no, mult], PANEL_COLOR, 1.2)
		"fate_wheel_perfect":
			var boost_mult = data.get("boost_mult", 11.0)
			var boost_secs = data.get("boost_secs", 24)
			show_settle(PANEL_ICON + " 命運完美！", "全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", Color(1.0, 0.6, 0.0))
			flash_screen(Color(1.0, 0.5, 0.0))
			hide_panel()
		"fate_wheel_end":
			hide_panel()
