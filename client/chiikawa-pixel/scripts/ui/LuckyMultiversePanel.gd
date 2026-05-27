## LuckyMultiversePanel.gd — T176 幸運多重宇宙魚 UI
## lucky-panel-agent 負責維護
## DAY-314：多重宇宙系統 — 開啟 3 個平行宇宙，每個宇宙擊破 5 個目標，全部完成 → 全服 ×13.0 加成 28 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.49, 0.11, 0.64)  # 深紫色（宇宙）
const PANEL_ICON = "🌌"
const PANEL_TITLE = "多重宇宙"

var _universe_labels: Array = []
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 71
	_setup_multiverse_ui()
	GameManager.lucky_multiverse.connect(_on_lucky_multiverse)

func _setup_multiverse_ui() -> void:
	# 三個宇宙進度指示器
	for i in range(3):
		var lbl = Label.new()
		lbl.text = "宇宙 %d: 0/5" % (i + 1)
		lbl.add_theme_color_override("font_color", Color(0.8, 0.5, 1.0))
		lbl.add_theme_font_size_override("font_size", 18)
		lbl.position = Vector2(20, 80 + i * 28)
		add_child(lbl)
		_universe_labels.append(lbl)

	# 全服加成顯示
	_boost_label = Label.new()
	_boost_label.text = "全服 ×13.0 加成 28 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 170)
	add_child(_boost_label)

func _on_lucky_multiverse(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"multiverse_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開啟！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.5, 0.1, 0.8))
			for lbl in _universe_labels:
				if is_instance_valid(lbl):
					lbl.text = "宇宙 %d: 0/5" % (_universe_labels.find(lbl) + 1)
		"universe_progress":
			var universes = data.get("universes", [0, 0, 0])
			var target = data.get("target", 5)
			for i in range(min(3, universes.size())):
				if i < _universe_labels.size() and is_instance_valid(_universe_labels[i]):
					var cnt = universes[i]
					_universe_labels[i].text = "宇宙 %d: %d/%d" % [i + 1, cnt, target]
					_universe_labels[i].modulate = Color(0.3, 1.0, 0.3) if cnt >= target else Color(0.8, 0.5, 1.0)
		"multiverse_perfect":
			var boost_mult = data.get("boost_mult", 13.0)
			var boost_secs = data.get("boost_secs", 28)
			show_settle(PANEL_ICON + " 多重宇宙完美！", "全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", Color(0.8, 0.3, 1.0))
			flash_screen(Color(0.6, 0.0, 1.0))
			hide_panel()
		"multiverse_end":
			hide_panel()
