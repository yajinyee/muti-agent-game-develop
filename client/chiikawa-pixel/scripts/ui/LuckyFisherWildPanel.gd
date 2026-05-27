## LuckyFisherWildPanel.gd — T183 幸運漁夫野生魚 UI
## lucky-panel-agent 負責維護
## DAY-315：漁夫野生系統 — 標記 3 個 Wild 目標（擊破獎勵 ×5.0），全部擊破 → 全服 ×17.0 加成 35 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.09, 0.40, 0.75)  # 深藍色（漁夫）
const PANEL_ICON = "🎣"
const PANEL_TITLE = "漁夫野生"

var _wild_label: Label = null
var _timer_label: Label = null
var _timer: float = 0.0
var _active: bool = false

func _ready() -> void:
	super._ready()
	layer = 78
	_setup_fisher_wild_ui()
	GameManager.lucky_fisher_wild.connect(_on_lucky_fisher_wild)

func _setup_fisher_wild_ui() -> void:
	_wild_label = Label.new()
	_wild_label.text = "Wild 目標：0 / 3"
	_wild_label.add_theme_color_override("font_color", Color(0.3, 0.7, 1.0))
	_wild_label.add_theme_font_size_override("font_size", 22)
	_wild_label.position = Vector2(20, 80)
	add_child(_wild_label)

	_timer_label = Label.new()
	_timer_label.text = "剩餘：30 秒"
	_timer_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_timer_label.add_theme_font_size_override("font_size", 18)
	_timer_label.position = Vector2(20, 110)
	add_child(_timer_label)

func _process(delta: float) -> void:
	if _active and _timer > 0:
		_timer -= delta
		if is_instance_valid(_timer_label):
			_timer_label.text = "剩餘：" + str(int(_timer)) + " 秒"
		if _timer <= 0:
			_active = false

func _on_lucky_fisher_wild(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"fisher_wild_start":
			var wild_count = data.get("wild_count", 3)
			_timer = 30.0
			_active = true
			if is_instance_valid(_wild_label):
				_wild_label.text = "Wild 目標：0 / " + str(wild_count)
			show_banner(PANEL_ICON + " " + PANEL_TITLE + "！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.1, 0.4, 0.8))
		"wild_killed":
			var killed = data.get("killed", 0)
			var total = data.get("total", 3)
			if is_instance_valid(_wild_label):
				_wild_label.text = "Wild 目標：" + str(killed) + " / " + str(total)
		"fisher_wild_perfect":
			_active = false
			var killed = data.get("killed", 3)
			var boost_mult = data.get("boost_mult", 17.0)
			var boost_secs = data.get("boost_secs", 35)
			show_settle(PANEL_ICON + " 完美漁夫！", "擊破 " + str(killed) + " 個 Wild！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", PANEL_COLOR)
			flash_screen(Color(0.1, 0.5, 1.0))
			hide_panel()
		"fisher_wild_timeout":
			_active = false
			var killed = data.get("killed", 0)
			var total = data.get("total", 3)
			show_settle(PANEL_ICON + " 時間到！", "擊破 " + str(killed) + " / " + str(total) + " 個 Wild 目標", PANEL_COLOR)
			hide_panel()
