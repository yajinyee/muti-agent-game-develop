## LuckyTimeLoopPanel.gd — T177 幸運時間迴圈魚 UI
## lucky-panel-agent 負責維護
## DAY-314：時間迴圈系統 — 3 次迴圈（每次 15 秒，獎勵 ×1.5 遞增），全部完成 → 全服 ×10.0 加成 22 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.08, 0.39, 0.75)  # 深藍色（時間）
const PANEL_ICON = "⏰"
const PANEL_TITLE = "時間迴圈"

var _loop_label: Label = null
var _mult_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 72
	_setup_time_loop_ui()
	GameManager.lucky_time_loop.connect(_on_lucky_time_loop)

func _setup_time_loop_ui() -> void:
	# 迴圈計數
	_loop_label = Label.new()
	_loop_label.text = "迴圈: 0/3"
	_loop_label.add_theme_color_override("font_color", Color(0.4, 0.7, 1.0))
	_loop_label.add_theme_font_size_override("font_size", 22)
	_loop_label.position = Vector2(20, 80)
	add_child(_loop_label)

	# 當前倍率
	_mult_label = Label.new()
	_mult_label.text = "獎勵倍率: ×1.5"
	_mult_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_mult_label.add_theme_font_size_override("font_size", 20)
	_mult_label.position = Vector2(20, 110)
	add_child(_mult_label)

	# 全服加成顯示
	_boost_label = Label.new()
	_boost_label.text = "完成全部 → 全服 ×10.0 加成 22 秒"
	_boost_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	_boost_label.add_theme_font_size_override("font_size", 14)
	_boost_label.position = Vector2(20, 145)
	add_child(_boost_label)

func _on_lucky_time_loop(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"time_loop_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開始！", PANEL_COLOR)
			show_panel()
			if is_instance_valid(_loop_label):
				_loop_label.text = "迴圈: 0/3"
			if is_instance_valid(_mult_label):
				_mult_label.text = "獎勵倍率: ×1.5"
		"loop_reset":
			var loop_no = data.get("loop_no", 1)
			var loop_mult = data.get("loop_mult", 1.5)
			if is_instance_valid(_loop_label):
				_loop_label.text = "迴圈: %d/3" % loop_no
			if is_instance_valid(_mult_label):
				_mult_label.text = "獎勵倍率: ×%.1f" % loop_mult
			show_banner(PANEL_ICON + " 第 %d 次迴圈！×%.1f！" % [loop_no, loop_mult], PANEL_COLOR, 1.5)
		"time_loop_perfect":
			var boost_mult = data.get("boost_mult", 10.0)
			var boost_secs = data.get("boost_secs", 22)
			show_settle(PANEL_ICON + " 時間迴圈完美！", "全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", Color(0.2, 0.6, 1.0))
			flash_screen(Color(0.0, 0.5, 1.0))
			hide_panel()
