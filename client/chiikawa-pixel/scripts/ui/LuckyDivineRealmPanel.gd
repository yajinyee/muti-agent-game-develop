## LuckyDivineRealmPanel.gd — T179 幸運神域降臨魚 UI
## lucky-panel-agent 負責維護
## DAY-314：神域降臨系統 — 5 波神域光柱（每波 HP -35%），5 波全部命中 ≥6 → 全服 ×14.0 加成 30 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.98, 0.66, 0.15)  # 神聖橙金色
const PANEL_ICON = "✨"
const PANEL_TITLE = "神域降臨"

var _wave_dots: Array = []
var _wave_label: Label = null
var _boost_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 74
	_setup_divine_realm_ui()
	GameManager.lucky_divine_realm.connect(_on_lucky_divine_realm)

func _setup_divine_realm_ui() -> void:
	# 5 個波次指示點
	for i in range(5):
		var dot = ColorRect.new()
		dot.size = Vector2(24, 24)
		dot.position = Vector2(20 + i * 32, 80)
		dot.color = Color(0.3, 0.3, 0.3)
		add_child(dot)
		_wave_dots.append(dot)

	# 波次標籤
	_wave_label = Label.new()
	_wave_label.text = "波次: 0/5"
	_wave_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_wave_label.add_theme_font_size_override("font_size", 20)
	_wave_label.position = Vector2(20, 115)
	add_child(_wave_label)

	# 全服加成顯示
	_boost_label = Label.new()
	_boost_label.text = "5 波完美 → 全服 ×14.0 加成 30 秒"
	_boost_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	_boost_label.add_theme_font_size_override("font_size", 14)
	_boost_label.position = Vector2(20, 148)
	add_child(_boost_label)

func _on_lucky_divine_realm(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"divine_realm_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 降臨！", PANEL_COLOR)
			show_panel()
			for dot in _wave_dots:
				if is_instance_valid(dot):
					dot.color = Color(0.3, 0.3, 0.3)
			if is_instance_valid(_wave_label):
				_wave_label.text = "波次: 0/5"
		"divine_wave":
			var wave_no = data.get("wave_no", 1)
			var hits = data.get("hits", 0)
			if is_instance_valid(_wave_label):
				_wave_label.text = "波次: %d/5 (命中 %d)" % [wave_no, hits]
			var idx = wave_no - 1
			if idx < _wave_dots.size() and is_instance_valid(_wave_dots[idx]):
				_wave_dots[idx].color = Color(0.3, 1.0, 0.3) if hits >= 6 else Color(1.0, 0.5, 0.0)
			show_banner(PANEL_ICON + " 第 %d 波！命中 %d 個！" % [wave_no, hits], PANEL_COLOR, 1.0)
		"divine_realm_perfect":
			var boost_mult = data.get("boost_mult", 14.0)
			var boost_secs = data.get("boost_secs", 30)
			show_settle(PANEL_ICON + " 神域完美！", "全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", Color(1.0, 0.85, 0.0))
			flash_screen(Color(1.0, 0.9, 0.0))
			hide_panel()
		"divine_realm_end":
			hide_panel()
