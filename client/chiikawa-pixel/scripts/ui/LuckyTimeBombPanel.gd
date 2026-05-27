## LuckyTimeBombPanel.gd — T162 幸運時間炸彈魚 UI
## lucky-panel-agent 負責維護
## DAY-310：時間炸彈系統 — 30 秒倒數，能量 ≥20 → 完美爆炸全服 ×6.0
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.9, 0.2, 0.1)  # 深紅色
const PANEL_ICON = "💣"
const PANEL_TITLE = "時間炸彈"

var _energy: int = 0
var _max_energy: int = 30
var _energy_label: Label = null
var _energy_bar: ProgressBar = null
var _countdown_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 57
	_setup_bomb_ui()
	GameManager.lucky_time_bomb.connect(_on_lucky_time_bomb)

func _setup_bomb_ui() -> void:
	# 炸彈能量顯示
	_energy_label = Label.new()
	_energy_label.text = "能量: 0/30"
	_energy_label.add_theme_color_override("font_color", Color(1.0, 0.8, 0.1))
	_energy_label.add_theme_font_size_override("font_size", 18)
	_energy_label.position = Vector2(20, 80)
	add_child(_energy_label)

	# 能量條
	_energy_bar = ProgressBar.new()
	_energy_bar.min_value = 0
	_energy_bar.max_value = _max_energy
	_energy_bar.value = 0
	_energy_bar.size = Vector2(200, 20)
	_energy_bar.position = Vector2(20, 110)
	add_child(_energy_bar)

	# 倒數標籤
	_countdown_label = Label.new()
	_countdown_label.text = "💣 30s"
	_countdown_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.1))
	_countdown_label.add_theme_font_size_override("font_size", 24)
	_countdown_label.position = Vector2(20, 140)
	add_child(_countdown_label)

func _on_lucky_time_bomb(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"start":
			_energy = 0
			_update_bomb_display()
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 啟動！", PANEL_COLOR)
			start_timer(data.get("duration", 30))
			show_panel()
		"energy_update":
			_energy = data.get("energy", 0)
			_update_bomb_display()
			# 能量高時閃爍警告
			if _energy >= 20:
				flash_screen(Color(1.0, 0.4, 0.1))
		"explode":
			var perfect = data.get("perfect", false)
			var final_energy = data.get("final_energy", 0)
			var dmg_pct = data.get("dmg_pct", 0.0)
			if perfect:
				show_settle("💣 完美爆炸！", "能量 " + str(final_energy) + "！全服 ×6.0！", PANEL_COLOR)
				flash_screen(Color(1.0, 0.2, 0.0))
			else:
				show_settle("爆炸！", "傷害 " + ("%.0f" % (dmg_pct * 100)) + "%", PANEL_COLOR)
			hide_panel()

func _update_bomb_display() -> void:
	if is_instance_valid(_energy_label):
		_energy_label.text = "能量: " + str(_energy) + "/" + str(_max_energy)
	if is_instance_valid(_energy_bar):
		_energy_bar.value = _energy
		# 能量高時變紅
		if _energy >= 20:
			_energy_bar.modulate = Color(1.0, 0.3, 0.1)
		else:
			_energy_bar.modulate = Color(1.0, 0.8, 0.1)
