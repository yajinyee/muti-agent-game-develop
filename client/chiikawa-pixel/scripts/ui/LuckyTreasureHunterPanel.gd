## LuckyTreasureHunterPanel.gd — T164 幸運寶藏獵人魚 UI
## lucky-panel-agent 負責維護
## DAY-310：寶藏獵人系統 — 5 個隨機寶藏，全部擊破 → 完美寶藏全服 ×7.0
extends BaseLuckyPanel

const PANEL_COLOR = Color(1.0, 0.75, 0.1)  # 金色
const PANEL_ICON = "💎"
const PANEL_TITLE = "寶藏獵人"

var _found: int = 0
var _total: int = 5
var _treasure_ids: Dictionary = {}  # instanceID -> mult
var _found_label: Label = null
var _treasure_dots: Array = []

func _ready() -> void:
	super._ready()
	layer = 59
	_setup_treasure_ui()
	GameManager.lucky_treasure_hunter.connect(_on_lucky_treasure_hunter)

func _setup_treasure_ui() -> void:
	_found_label = Label.new()
	_found_label.text = "寶藏: 0/5"
	_found_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	_found_label.add_theme_font_size_override("font_size", 20)
	_found_label.position = Vector2(20, 80)
	add_child(_found_label)

	# 5 個寶藏指示點
	for i in range(5):
		var dot = Label.new()
		dot.text = "💎"
		dot.add_theme_font_size_override("font_size", 18)
		dot.position = Vector2(20 + i * 40, 115)
		dot.modulate = Color(0.5, 0.5, 0.5)  # 灰色 = 未找到
		add_child(dot)
		_treasure_dots.append(dot)

func _on_lucky_treasure_hunter(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"start":
			_found = 0
			_treasure_ids.clear()
			var treasures = data.get("treasures", [])
			for t in treasures:
				_treasure_ids[t.get("instance_id", "")] = t.get("mult", 10.0)
			_update_treasure_display()
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開始！", PANEL_COLOR)
			start_timer(data.get("duration", 30))
			show_panel()
		"treasure_found":
			_found = data.get("found", 0)
			_update_treasure_display()
			flash_screen(PANEL_COLOR)
			if data.get("perfect", false):
				show_settle("💎 完美寶藏！", "全服 ×7.0 加成 15 秒！", PANEL_COLOR)
				hide_panel()

func _update_treasure_display() -> void:
	if is_instance_valid(_found_label):
		_found_label.text = "寶藏: " + str(_found) + "/" + str(_total)
	for i in range(min(_treasure_dots.size(), _total)):
		if is_instance_valid(_treasure_dots[i]):
			if i < _found:
				_treasure_dots[i].modulate = Color(1.0, 0.9, 0.2)  # 金色 = 已找到
			else:
				_treasure_dots[i].modulate = Color(0.5, 0.5, 0.5)  # 灰色 = 未找到
