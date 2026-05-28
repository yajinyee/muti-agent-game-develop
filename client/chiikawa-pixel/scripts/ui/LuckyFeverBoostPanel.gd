## LuckyFeverBoostPanel.gd — T206 幸運 Fever Boost 魚 UI
## lucky-panel-agent 負責維護
## DAY-323：Fever Boost™ 系統 — 30 秒內所有特效機率翻倍，全場倍率 ×2.0，全服 ×31.0 加成 62 秒
## 業界依據：Games Global「Fishin' Pots of Gold」Fever Boost™（2026-05-28 最新）
extends BaseLuckyPanel

const PANEL_COLOR = Color(1.0, 0.4, 0.0)   # 火橙色（Fever Boost）
const PANEL_ICON = "🔥"
const PANEL_TITLE = "Fever Boost"

var _boost_label: Label = null
var _timer_label: Label = null
var _global_label: Label = null
var _fever_timer: float = 0.0
var _fever_active: bool = false

func _ready() -> void:
	super._ready()
	layer = 101
	_setup_fever_boost_ui()
	GameManager.lucky_fever_boost.connect(_on_lucky_fever_boost)

func _setup_fever_boost_ui() -> void:
	_boost_label = Label.new()
	_boost_label.text = "🔥 Fever Boost™ 啟動！"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 0.6, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 22)
	_boost_label.position = Vector2(20, 52)
	add_child(_boost_label)

	_timer_label = Label.new()
	_timer_label.text = "全場倍率 ×2.0 | 30 秒"
	_timer_label.add_theme_color_override("font_color", Color(1.0, 0.8, 0.0))
	_timer_label.add_theme_font_size_override("font_size", 20)
	_timer_label.position = Vector2(20, 80)
	add_child(_timer_label)

	_global_label = Label.new()
	_global_label.text = "全服 ×31.0 加成 62 秒"
	_global_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_global_label.add_theme_font_size_override("font_size", 16)
	_global_label.position = Vector2(20, 110)
	add_child(_global_label)

func _process(delta: float) -> void:
	if _fever_active and _fever_timer > 0.0:
		_fever_timer -= delta
		if is_instance_valid(_timer_label):
			_timer_label.text = "全場倍率 ×2.0 | %.1f 秒" % max(0.0, _fever_timer)
		if _fever_timer <= 0.0:
			_fever_active = false

func _on_lucky_fever_boost(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"fever_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 觸發！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(1.0, 0.4, 0.0))
		"fever_active":
			var boost_secs = data.get("boost_secs", 30)
			_fever_timer = float(boost_secs)
			_fever_active = true
			if is_instance_valid(_boost_label):
				_boost_label.text = "🔥 Fever Boost™ 進行中！"
			if is_instance_valid(_timer_label):
				_timer_label.text = "全場倍率 ×2.0 | %d 秒" % boost_secs
		"fever_complete":
			var global_mult = data.get("global_mult", 31.0)
			var global_secs = data.get("global_secs", 62)
			_fever_active = false
			show_settle(PANEL_ICON + " Fever Boost 完成！",
				"全服 ×%.1f 加成 %d 秒！" % [global_mult, global_secs],
				PANEL_COLOR)
			flash_screen(Color(1.0, 1.0, 0.0))
			hide_panel()
