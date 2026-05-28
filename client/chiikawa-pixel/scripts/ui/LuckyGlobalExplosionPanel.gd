## LuckyGlobalExplosionPanel.gd — T193 幸運全服大爆炸魚 UI
## lucky-panel-agent 負責維護
## DAY-317：全服大爆炸系統 — 全場 HP 歸零（每個獎勵 ×15.0），觸發全服 ×20.5 加成 40 秒（新最高）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.80, 0.05, 0.05)  # 深紅色（全服大爆炸）
const PANEL_ICON = "💥"
const PANEL_TITLE = "全服大爆炸"

var _hit_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 88
	_setup_global_explosion_ui()
	GameManager.lucky_global_explosion.connect(_on_lucky_global_explosion)

func _setup_global_explosion_ui() -> void:
	_hit_label = Label.new()
	_hit_label.text = "清場: 0 個"
	_hit_label.add_theme_color_override("font_color", Color(1.0, 0.4, 0.4))
	_hit_label.add_theme_font_size_override("font_size", 22)
	_hit_label.position = Vector2(20, 80)
	add_child(_hit_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×20.5 加成 40 秒（新最高）"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 115)
	add_child(_boost_label)

func _on_lucky_global_explosion(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"global_explosion_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 引爆！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(1.0, 0.0, 0.0))
		"global_explosion_complete":
			var hit_count = data.get("hit_count", 0)
			var boost_mult = data.get("boost_mult", 20.5)
			var boost_secs = data.get("boost_secs", 40)
			if is_instance_valid(_hit_label):
				_hit_label.text = "清場: " + str(hit_count) + " 個"
			show_settle(PANEL_ICON + " 全服大爆炸！",
				"清場 " + str(hit_count) + " 個！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！（新最高）",
				PANEL_COLOR)
			flash_screen(Color(1.0, 0.2, 0.0))
			hide_panel()
