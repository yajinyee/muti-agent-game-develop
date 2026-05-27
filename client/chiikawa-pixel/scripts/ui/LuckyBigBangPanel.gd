## LuckyBigBangPanel.gd — T170 幸運宇宙大爆炸魚 UI
## lucky-panel-agent 負責維護
## DAY-312：宇宙大爆炸系統 — 全場 HP 歸零（每個獎勵 ×8.0），觸發全服 ×12.0 加成 25 秒（遊戲最高倍率）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.72, 0.07, 0.07)  # 深紅色（爆炸）
const PANEL_ICON = "💥"
const PANEL_TITLE = "宇宙大爆炸"

var _hit_count: int = 0
var _hit_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 65
	_setup_big_bang_ui()
	GameManager.lucky_big_bang.connect(_on_lucky_big_bang)

func _setup_big_bang_ui() -> void:
	# 清場計數
	_hit_label = Label.new()
	_hit_label.text = "清場: 0 個"
	_hit_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))
	_hit_label.add_theme_font_size_override("font_size", 20)
	_hit_label.position = Vector2(20, 80)
	add_child(_hit_label)

	# 全服加成顯示
	_boost_label = Label.new()
	_boost_label.text = "全服 ×12.0 加成 25 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 18)
	_boost_label.position = Vector2(20, 110)
	add_child(_boost_label)

func _on_lucky_big_bang(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"big_bang_start":
			_hit_count = 0
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 引爆！", PANEL_COLOR)
			show_panel()
			# 強烈全螢幕閃光
			flash_screen(Color(1.0, 0.2, 0.2))
		"big_bang_complete":
			_hit_count = data.get("hit_count", 0)
			var boost_mult = data.get("boost_mult", 12.0)
			var boost_secs = data.get("boost_secs", 25)
			if is_instance_valid(_hit_label):
				_hit_label.text = "清場: " + str(_hit_count) + " 個"
			show_settle("💥 宇宙大爆炸！", "清場 " + str(_hit_count) + " 個！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", PANEL_COLOR)
			flash_screen(Color(1.0, 0.5, 0.0))
			hide_panel()
