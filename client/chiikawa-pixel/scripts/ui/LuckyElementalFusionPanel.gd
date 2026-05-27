## LuckyElementalFusionPanel.gd — T163 幸運元素融合魚 UI
## lucky-panel-agent 負責維護
## DAY-310：元素融合系統 — 火/冰/雷三元素，全觸發 → 元素爆發全服 ×6.5
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.6, 0.2, 0.9)  # 深紫色
const PANEL_ICON = "⚡"
const PANEL_TITLE = "元素融合"

var _fire_triggered: bool = false
var _ice_triggered: bool = false
var _thunder_triggered: bool = false
var _element_labels: Array = []

func _ready() -> void:
	super._ready()
	layer = 58
	_setup_element_ui()
	GameManager.lucky_elemental_fusion.connect(_on_lucky_elemental_fusion)

func _setup_element_ui() -> void:
	var elements = [
		{"name": "🔥 火元素", "color": Color(1.0, 0.4, 0.1)},
		{"name": "❄️ 冰元素", "color": Color(0.4, 0.8, 1.0)},
		{"name": "⚡ 雷元素", "color": Color(1.0, 0.9, 0.2)},
	]
	for i in range(elements.size()):
		var lbl = Label.new()
		lbl.text = elements[i]["name"] + " ○"
		lbl.add_theme_color_override("font_color", elements[i]["color"])
		lbl.add_theme_font_size_override("font_size", 16)
		lbl.position = Vector2(20, 80 + i * 28)
		add_child(lbl)
		_element_labels.append(lbl)

func _on_lucky_elemental_fusion(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"start":
			_fire_triggered = false
			_ice_triggered = false
			_thunder_triggered = false
			_update_element_display()
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開始！", PANEL_COLOR)
			start_timer(data.get("duration", 25))
			show_panel()
		"element_trigger":
			var element = data.get("element", "")
			match element:
				"fire":
					_fire_triggered = true
					flash_screen(Color(1.0, 0.4, 0.1))
				"ice":
					_ice_triggered = true
					flash_screen(Color(0.4, 0.8, 1.0))
				"thunder":
					_thunder_triggered = true
					flash_screen(Color(1.0, 0.9, 0.2))
			_update_element_display()
		"end":
			var perfect = data.get("perfect", false)
			if perfect:
				show_settle("⚡🔥❄️ 元素爆發！", "全服 ×6.5 加成 14 秒！", PANEL_COLOR)
				flash_screen(Color(0.8, 0.4, 1.0))
			else:
				var count = [_fire_triggered, _ice_triggered, _thunder_triggered].count(true)
				show_settle("元素融合結束", str(count) + "/3 元素觸發", PANEL_COLOR)
			hide_panel()

func _update_element_display() -> void:
	var states = [_fire_triggered, _ice_triggered, _thunder_triggered]
	for i in range(min(_element_labels.size(), states.size())):
		if is_instance_valid(_element_labels[i]):
			var suffix = " ✓" if states[i] else " ○"
			var base_texts = ["🔥 火元素", "❄️ 冰元素", "⚡ 雷元素"]
			_element_labels[i].text = base_texts[i] + suffix
