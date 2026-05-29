## LuckyAtlantisFrenzyPanel.gd — T241 幸運大西洋狂潮魚面板
## lucky-panel-agent 負責維護
## Big Atlantis Frenzy 機制（業界依據：BGaming Big Atlantis Frenzy 2025-2026）
## 亞特蘭提斯爆炸 + 連鎖消除，7 波 Fish 符號（×5-×500），全服 ×52.0 加成 104 秒
extends Node

const PANEL_COLOR = Color(0.12, 0.56, 1.0)  # 亞特蘭提斯藍

var _canvas: CanvasLayer
var _banner: Dictionary
var _indicator: Dictionary
var _settle: Control
var _wave_count: int = 0

func _ready() -> void:
	_canvas = get_tree().get_root().get_node_or_null("Main/HUDLayer")
	if not is_instance_valid(_canvas):
		_canvas = CanvasLayer.new()
		_canvas.layer = 10
		get_tree().get_root().add_child(_canvas)
	_banner   = BaseLuckyPanel.create_banner(_canvas, 100.0, 65)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 160), 70)
	_settle   = BaseLuckyPanel.create_settle_popup(_canvas, 75)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"atlantis_start":
			_wave_count = 0
			var waves = data.get("waves", 7)
			var global_target = data.get("global_target", 52.0)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🌊🏛️ BIG ATLANTIS FRENZY! %s triggered! %d cascade waves! Fish x5-x500! Global x%.1f!" % [trigger_name, waves, global_target],
				PANEL_COLOR, 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 6)
			_show_indicator("ATLANTIS FRENZY", "Wave 0/7")
		"cascade_wave":
			_wave_count = data.get("wave_no", _wave_count + 1)
			var fish_count = data.get("fish_count", 0)
			var wave_mult = data.get("wave_mult", 0.0)
			_show_indicator("ATLANTIS FRENZY", "Wave %d: %d fish x%.1f" % [_wave_count, fish_count, wave_mult])
		"atlantis_complete":
			var waves = data.get("wave_count", _wave_count)
			var total = data.get("total_mult", 0.0)
			var bonus = data.get("global_bonus", 52.0)
			var secs = data.get("global_seconds", 104)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🌊🏛️ ATLANTIS FRENZY COMPLETE! %d waves! Total x%.1f! Global x%.1f for %ds!" % [waves, total, bonus, secs],
				Color(0.0, 0.8, 1.0), 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(0.0, 0.8, 1.0), 8)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "🌊🏛️ BIG ATLANTIS FRENZY!", "size": 22, "color": PANEL_COLOR},
				{"text": "Cascade Waves: %d" % waves, "size": 16, "color": Color.WHITE},
				{"text": "Total x%.1f" % total, "size": 18, "color": PANEL_COLOR},
				{"text": "🏆 GLOBAL x%.1f for %ds" % [bonus, secs], "size": 18, "color": Color(0.0, 0.8, 1.0)},
			], 5.0)
			_hide_indicator()

func _show_indicator(title: String, value: String) -> void:
	if _indicator.is_empty(): return
	var t = _indicator.get("title")
	var v = _indicator.get("value")
	var p = _indicator.get("panel")
	if is_instance_valid(t): t.text = title
	if is_instance_valid(v): v.text = value
	if is_instance_valid(p): p.visible = true

func _hide_indicator() -> void:
	var p = _indicator.get("panel")
	if is_instance_valid(p): p.visible = false
