## LuckyChaosExplosionPanel.gd — T198 幸運混沌爆炸魚 UI
## lucky-panel-agent 負責維護
## DAY-318：混沌爆炸系統 — 隨機 3-8 個目標同時爆炸，倍率疊加最高 ×30.0，全服 ×24.0 加成 48 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.15, 0.0, 0.0)   # 深紅色（混沌）
const PANEL_BORDER_COLOR = Color(1.0, 0.3, 0.0)  # 火橙色邊框
const PANEL_ICON = "💥"
const PANEL_TITLE = "混沌爆炸"

var _explode_label: Label = null
var _mult_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 93
	_setup_chaos_explosion_ui()
	GameManager.lucky_chaos_explosion.connect(_on_lucky_chaos_explosion)

func _setup_chaos_explosion_ui() -> void:
	_explode_label = Label.new()
	_explode_label.text = "爆炸目標: ?"
	_explode_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.2))
	_explode_label.add_theme_font_size_override("font_size", 22)
	_explode_label.position = Vector2(20, 55)
	add_child(_explode_label)

	_mult_label = Label.new()
	_mult_label.text = "倍率疊加: ×?"
	_mult_label.add_theme_color_override("font_color", Color(1.0, 0.8, 0.0))
	_mult_label.add_theme_font_size_override("font_size", 24)
	_mult_label.position = Vector2(20, 82)
	add_child(_mult_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×24.0 加成 48 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 0.6, 0.2))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 115)
	add_child(_boost_label)

func _on_lucky_chaos_explosion(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"chaos_explosion_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 引爆！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(1.0, 0.3, 0.0))
		"chaos_explosion_complete":
			var explode_count = data.get("explode_count", 0)
			var total_mult = data.get("total_mult", 1.0)
			var boost_mult = data.get("boost_mult", 24.0)
			var boost_secs = data.get("boost_secs", 48)
			if is_instance_valid(_explode_label):
				_explode_label.text = "爆炸目標: " + str(explode_count) + " 個"
			if is_instance_valid(_mult_label):
				_mult_label.text = "倍率疊加: ×" + str(snapped(total_mult, 0.1))
			show_settle(PANEL_ICON + " 混沌爆炸！",
				"爆炸 " + str(explode_count) + " 個目標！倍率疊加 ×" + str(snapped(total_mult, 0.1)) + "！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！",
				PANEL_COLOR)
			flash_screen(Color(1.0, 0.5, 0.0))
			hide_panel()
