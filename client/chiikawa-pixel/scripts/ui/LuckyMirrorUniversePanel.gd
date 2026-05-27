## LuckyMirrorUniversePanel.gd — T186 幸運鏡像宇宙魚 UI
## lucky-panel-agent 負責維護
## DAY-316：鏡像宇宙系統 — 複製場上最強 3 個目標（HP 50%，獎勵 ×2.0），全服 ×17.0 加成 36 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.10, 0.14, 0.55)  # 深藍色（鏡像宇宙）
const PANEL_ICON = "🪞"
const PANEL_TITLE = "鏡像宇宙"

var _mirror_count: int = 0
var _mirror_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 81
	_setup_mirror_universe_ui()
	GameManager.lucky_mirror_universe.connect(_on_lucky_mirror_universe)

func _setup_mirror_universe_ui() -> void:
	_mirror_label = Label.new()
	_mirror_label.text = "鏡像目標: 0 個"
	_mirror_label.add_theme_color_override("font_color", Color(0.5, 0.7, 1.0))
	_mirror_label.add_theme_font_size_override("font_size", 22)
	_mirror_label.position = Vector2(20, 80)
	add_child(_mirror_label)

	_boost_label = Label.new()
	_boost_label.text = "全部擊破 → 全服 ×17.0 加成 36 秒"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 16)
	_boost_label.position = Vector2(20, 115)
	add_child(_boost_label)

func _on_lucky_mirror_universe(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"mirror_universe_start":
			_mirror_count = 0
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開啟！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.1, 0.2, 0.8))
		"mirror_universe_active":
			_mirror_count = data.get("mirror_count", 0)
			if is_instance_valid(_mirror_label):
				_mirror_label.text = "鏡像目標: " + str(_mirror_count) + " 個（HP 50%，×2.0）"
		"mirror_perfect":
			var boost_mult = data.get("boost_mult", 17.0)
			var boost_secs = data.get("boost_secs", 36)
			show_settle(PANEL_ICON + " 鏡像完美！", "全部擊破！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", PANEL_COLOR)
			flash_screen(Color(0.2, 0.4, 1.0))
			hide_panel()
		"mirror_universe_end":
			hide_panel()
