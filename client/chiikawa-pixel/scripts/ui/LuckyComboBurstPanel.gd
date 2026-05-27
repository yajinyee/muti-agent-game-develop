## LuckyComboBurstPanel.gd — T161 幸運連擊爆發魚 UI
## lucky-panel-agent 負責維護
## DAY-310：連擊爆發系統 — 連擊累積倍率最高 ×15.0，Combo ≥10 → 完美連擊全服 ×5.5
extends BaseLuckyPanel

const PANEL_COLOR = Color(1.0, 0.4, 0.1)  # 火橙色
const PANEL_ICON = "🔥"
const PANEL_TITLE = "連擊爆發"

var _combo_count: int = 0
var _current_mult: float = 1.0
var _max_combo: int = 10
var _max_mult: float = 15.0
var _combo_label: Label = null
var _mult_label: Label = null
var _combo_bar: ProgressBar = null

func _ready() -> void:
	super._ready()
	layer = 56
	_setup_combo_ui()
	GameManager.lucky_combo_burst.connect(_on_lucky_combo_burst)

func _setup_combo_ui() -> void:
	# 連擊計數器
	_combo_label = Label.new()
	_combo_label.text = "COMBO: 0"
	_combo_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	_combo_label.add_theme_font_size_override("font_size", 20)
	_combo_label.position = Vector2(20, 80)
	add_child(_combo_label)

	# 倍率顯示
	_mult_label = Label.new()
	_mult_label.text = "×1.0"
	_mult_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.1))
	_mult_label.add_theme_font_size_override("font_size", 28)
	_mult_label.position = Vector2(20, 110)
	add_child(_mult_label)

	# 連擊進度條
	_combo_bar = ProgressBar.new()
	_combo_bar.min_value = 0
	_combo_bar.max_value = _max_combo
	_combo_bar.value = 0
	_combo_bar.size = Vector2(200, 16)
	_combo_bar.position = Vector2(20, 145)
	add_child(_combo_bar)

func _on_lucky_combo_burst(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"start":
			_combo_count = 0
			_current_mult = 1.0
			_update_combo_display()
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開始！", PANEL_COLOR)
			start_timer(data.get("duration", 20))
			show_panel()
		"combo_update":
			_combo_count = data.get("combo", 0)
			_current_mult = data.get("mult", 1.0)
			_update_combo_display()
			# 連擊特效
			if _combo_count % 5 == 0:
				flash_screen(PANEL_COLOR)
		"end":
			var perfect = data.get("perfect", false)
			if perfect:
				show_settle("🔥 完美連擊！", "全服 ×5.5 加成 12 秒！", PANEL_COLOR)
				flash_screen(Color(1.0, 0.6, 0.1))
			else:
				show_settle("連擊結束", "Combo: " + str(_combo_count), PANEL_COLOR)
			hide_panel()

func _update_combo_display() -> void:
	if is_instance_valid(_combo_label):
		_combo_label.text = "COMBO: " + str(_combo_count)
	if is_instance_valid(_mult_label):
		_mult_label.text = "×" + ("%.1f" % _current_mult)
	if is_instance_valid(_combo_bar):
		_combo_bar.value = min(_combo_count, _max_combo)
