## LuckyStarPortalPanel.gd — T166 幸運星際門戶魚 UI
## lucky-panel-agent 負責維護
## DAY-312：星際門戶系統 — 傳送 5 個目標到中央（HP -50%），完美門戶全服 ×5.5 加成 12 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.49, 0.11, 0.64)  # 深紫色（星際）
const PANEL_ICON = "🌌"
const PANEL_TITLE = "星際門戶"

var _teleport_count: int = 0
var _target_teleports: int = 5
var _teleport_label: Label = null
var _teleport_bar: ProgressBar = null

func _ready() -> void:
	super._ready()
	layer = 61
	_setup_portal_ui()
	GameManager.lucky_star_portal.connect(_on_lucky_star_portal)

func _setup_portal_ui() -> void:
	# 傳送計數
	_teleport_label = Label.new()
	_teleport_label.text = "傳送: 0/5"
	_teleport_label.add_theme_color_override("font_color", Color(0.8, 0.4, 1.0))
	_teleport_label.add_theme_font_size_override("font_size", 20)
	_teleport_label.position = Vector2(20, 80)
	add_child(_teleport_label)

	# 傳送進度條
	_teleport_bar = ProgressBar.new()
	_teleport_bar.min_value = 0
	_teleport_bar.max_value = _target_teleports
	_teleport_bar.value = 0
	_teleport_bar.size = Vector2(200, 16)
	_teleport_bar.position = Vector2(20, 110)
	add_child(_teleport_bar)

func _on_lucky_star_portal(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"portal_open":
			_teleport_count = 0
			_update_portal_display()
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開啟！", PANEL_COLOR)
			show_panel()
			flash_screen(PANEL_COLOR)
		"portal_teleport":
			_teleport_count = data.get("teleport_count", 0)
			_update_portal_display()
		"portal_perfect":
			var boost_mult = data.get("boost_mult", 5.5)
			var boost_secs = data.get("boost_secs", 12)
			show_settle("🌌 完美門戶！", "傳送 " + str(_teleport_count) + " 個！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", PANEL_COLOR)
			flash_screen(Color(0.8, 0.4, 1.0))
			hide_panel()
		"portal_end":
			show_settle("星際門戶關閉", "傳送 " + str(_teleport_count) + " 個目標", PANEL_COLOR)
			hide_panel()

func _update_portal_display() -> void:
	if is_instance_valid(_teleport_label):
		_teleport_label.text = "傳送: " + str(_teleport_count) + "/" + str(_target_teleports)
	if is_instance_valid(_teleport_bar):
		_teleport_bar.value = min(_teleport_count, _target_teleports)
