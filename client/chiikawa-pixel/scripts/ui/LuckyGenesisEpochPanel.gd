## LuckyGenesisEpochPanel.gd — T200 幸運創世紀元魚 UI
## lucky-panel-agent 負責維護
## DAY-318：創世紀元系統 — 里程碑第 200 個 Lucky 目標，全場 HP 歸零（每個獎勵 ×25.0），全服 ×25.0 加成 50 秒（史上最高）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.0, 0.0, 0.0)   # 純黑色（創世紀元）
const PANEL_BORDER_COLOR = Color(1.0, 1.0, 1.0)  # 純白色邊框（最高階）
const PANEL_ICON = "🌌"
const PANEL_TITLE = "創世紀元"

var _milestone_label: Label = null
var _hit_label: Label = null
var _boost_label: Label = null
var _crown_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 95
	_setup_genesis_epoch_ui()
	GameManager.lucky_genesis_epoch.connect(_on_lucky_genesis_epoch)

func _setup_genesis_epoch_ui() -> void:
	_crown_label = Label.new()
	_crown_label.text = "🎊 里程碑：第 200 個 Lucky 目標 🎊"
	_crown_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	_crown_label.add_theme_font_size_override("font_size", 16)
	_crown_label.position = Vector2(20, 52)
	add_child(_crown_label)

	_milestone_label = Label.new()
	_milestone_label.text = "全場清空 ×25.0（史上最高）"
	_milestone_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.0))
	_milestone_label.add_theme_font_size_override("font_size", 20)
	_milestone_label.position = Vector2(20, 75)
	add_child(_milestone_label)

	_hit_label = Label.new()
	_hit_label.text = "清場: 0 個"
	_hit_label.add_theme_color_override("font_color", Color(0.8, 0.8, 1.0))
	_hit_label.add_theme_font_size_override("font_size", 22)
	_hit_label.position = Vector2(20, 100)
	add_child(_hit_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×25.0 加成 50 秒（史上最高）"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 130)
	add_child(_boost_label)

func _on_lucky_genesis_epoch(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"genesis_epoch_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開啟！", PANEL_COLOR)
			show_panel()
			# 特殊：黑色閃屏 → 白色閃屏 → 金色閃屏（三重閃屏）
			flash_screen(Color(0.0, 0.0, 0.0))
			var tween = create_tween()
			tween.tween_interval(0.1)
			tween.tween_callback(func(): flash_screen(Color(1.0, 1.0, 1.0)))
			tween.tween_interval(0.1)
			tween.tween_callback(func(): flash_screen(Color(1.0, 0.9, 0.0)))
		"genesis_epoch_complete":
			var hit_count = data.get("hit_count", 0)
			var boost_mult = data.get("boost_mult", 25.0)
			var boost_secs = data.get("boost_secs", 50)
			if is_instance_valid(_hit_label):
				_hit_label.text = "清場: " + str(hit_count) + " 個"
			show_settle(PANEL_ICON + " 創世紀元！",
				"清場 " + str(hit_count) + " 個！每個 ×25.0！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！（里程碑第 200 個 Lucky 目標，史上最高）",
				PANEL_COLOR)
			flash_screen(Color(1.0, 1.0, 1.0))
			hide_panel()
