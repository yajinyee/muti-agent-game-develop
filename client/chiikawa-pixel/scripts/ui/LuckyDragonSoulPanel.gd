## LuckyDragonSoulPanel.gd — T167 幸運龍魂融合魚 UI
## lucky-panel-agent 負責維護
## DAY-312：龍魂融合系統 — 30 秒吸收龍魂（最高 50 魂），50 魂 → 龍魂爆發全場 HP -90%，全服 ×9.0 加成 18 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.83, 0.18, 0.18)  # 龍紅色
const PANEL_ICON = "🐉"
const PANEL_TITLE = "龍魂融合"

var _soul_count: int = 0
var _max_souls: int = 50
var _soul_label: Label = null
var _soul_bar: ProgressBar = null

func _ready() -> void:
	super._ready()
	layer = 62
	_setup_dragon_soul_ui()
	GameManager.lucky_dragon_soul.connect(_on_lucky_dragon_soul)

func _setup_dragon_soul_ui() -> void:
	# 龍魂計數
	_soul_label = Label.new()
	_soul_label.text = "龍魂: 0/50"
	_soul_label.add_theme_color_override("font_color", Color(1.0, 0.4, 0.1))
	_soul_label.add_theme_font_size_override("font_size", 20)
	_soul_label.position = Vector2(20, 80)
	add_child(_soul_label)

	# 龍魂進度條
	_soul_bar = ProgressBar.new()
	_soul_bar.min_value = 0
	_soul_bar.max_value = _max_souls
	_soul_bar.value = 0
	_soul_bar.size = Vector2(200, 16)
	_soul_bar.position = Vector2(20, 110)
	add_child(_soul_bar)

func _on_lucky_dragon_soul(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"soul_fusion_start":
			_soul_count = 0
			_max_souls = data.get("max_souls", 50)
			_update_soul_display()
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開始！", PANEL_COLOR)
			start_timer(data.get("duration", 30))
			show_panel()
			flash_screen(PANEL_COLOR)
		"soul_burst":
			_soul_count = data.get("soul_count", 0)
			_update_soul_display()
		"soul_perfect":
			var boost_mult = data.get("boost_mult", 9.0)
			var boost_secs = data.get("boost_secs", 18)
			show_settle("🐉 龍魂完美！", "集滿 50 魂！全場 HP -90%！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", PANEL_COLOR)
			flash_screen(Color(1.0, 0.44, 0.07))
			hide_panel()
		"soul_end":
			show_settle("龍魂融合結束", "吸收 " + str(_soul_count) + " 魂", PANEL_COLOR)
			hide_panel()

func _update_soul_display() -> void:
	if is_instance_valid(_soul_label):
		_soul_label.text = "龍魂: " + str(_soul_count) + "/" + str(_max_souls)
	if is_instance_valid(_soul_bar):
		_soul_bar.value = min(_soul_count, _max_souls)
