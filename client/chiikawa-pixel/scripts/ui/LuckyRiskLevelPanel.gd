## LuckyRiskLevelPanel.gd — T184 幸運風險等級魚 UI
## lucky-panel-agent 負責維護
## DAY-315：風險等級系統 — 5 等級選擇（低 ×5.0 / 中 ×20.0 / 高 ×100.0 / 極高 ×500.0 / 最高 ×3000.0）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.90, 0.32, 0.0)  # 火橙色（風險）
const PANEL_ICON = "🎰"
const PANEL_TITLE = "風險等級"

var _level_label: Label = null
var _mult_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 79
	_setup_risk_level_ui()
	GameManager.lucky_risk_level.connect(_on_lucky_risk_level)

func _setup_risk_level_ui() -> void:
	_level_label = Label.new()
	_level_label.text = "等級：選擇中..."
	_level_label.add_theme_color_override("font_color", Color(1.0, 0.6, 0.2))
	_level_label.add_theme_font_size_override("font_size", 20)
	_level_label.position = Vector2(20, 80)
	add_child(_level_label)

	_mult_label = Label.new()
	_mult_label.text = "最高 ×3000"
	_mult_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_mult_label.add_theme_font_size_override("font_size", 18)
	_mult_label.position = Vector2(20, 110)
	add_child(_mult_label)

func _on_lucky_risk_level(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"risk_level_start":
			if is_instance_valid(_level_label):
				_level_label.text = "等級：選擇中..."
			if is_instance_valid(_mult_label):
				_mult_label.text = "最高 ×3000"
			show_banner(PANEL_ICON + " " + PANEL_TITLE + "！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.9, 0.3, 0.0))
		"risk_level_jackpot":
			var risk_name = data.get("risk_name", "最高風險")
			var risk_mult = data.get("risk_mult", 3000.0)
			var boost_mult = data.get("boost_mult", 17.5)
			var boost_secs = data.get("boost_secs", 36)
			if is_instance_valid(_level_label):
				_level_label.text = "等級：" + risk_name
			if is_instance_valid(_mult_label):
				_mult_label.text = "×" + str(risk_mult) + " 大獎！"
			show_settle(PANEL_ICON + " 最高風險！", "抽中 ×" + str(risk_mult) + "！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", PANEL_COLOR)
			flash_screen(Color(1.0, 0.85, 0.0))
			hide_panel()
		"risk_level_result":
			var risk_name = data.get("risk_name", "低風險")
			var risk_mult = data.get("risk_mult", 5.0)
			if is_instance_valid(_level_label):
				_level_label.text = "等級：" + risk_name
			if is_instance_valid(_mult_label):
				_mult_label.text = "×" + str(risk_mult)
			show_settle(PANEL_ICON + " 風險結果！", "【" + risk_name + "】×" + str(risk_mult) + "！", PANEL_COLOR)
			hide_panel()
