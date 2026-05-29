## LuckyWildCollectorPanel.gd — T244 幸運野生收集魚面板
## lucky-panel-agent 負責維護
## Wild Collector 機制（BGaming Big Boat Big Catch 升級版）
## Wild 收集：每 4 個 Wild → 額外旋轉（×2→×3→×10），全服 ×54.0 加成 108 秒
extends Node

const PANEL_COLOR = Color(1.0, 0.85, 0.0)  # 黃金色

var _canvas: CanvasLayer
var _banner: Dictionary
var _indicator: Dictionary
var _settle: Control
var _stage: int = 0
var _spin_count: int = 0

func _ready() -> void:
	_canvas = get_tree().get_root().get_node_or_null("Main/HUDLayer")
	if not is_instance_valid(_canvas):
		_canvas = CanvasLayer.new()
		_canvas.layer = 10
		get_tree().get_root().add_child(_canvas)
	_banner    = BaseLuckyPanel.create_banner(_canvas, 100.0, 68)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 160), 73)
	_settle    = BaseLuckyPanel.create_settle_popup(_canvas, 78)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"wild_collector_start":
			_stage = 0
			_spin_count = 0
			BaseLuckyPanel.show_banner(_banner.panel,
				"🃏✨ WILD COLLECTOR! %s activated! Collecting Wilds for bonus spins x2→x3→x10!" % trigger_name,
				PANEL_COLOR, 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 5)
			_show_indicator("WILD COLLECTOR", "Stage 1/3 🃏")
		"wild_collected":
			_stage = data.get("stage", _stage)
			var spin_mult = data.get("spin_mult", 2.0)
			var bonus_spins = data.get("bonus_spins", 3)
			_show_indicator("WILD COLLECTOR", "Stage %d: x%.0f (%d spins)" % [_stage, spin_mult, bonus_spins])
		"bonus_spin":
			_spin_count += 1
			var spin_mult = data.get("spin_mult", 2.0)
			_show_indicator("WILD COLLECTOR", "Spin %d x%.0f 🎰" % [_spin_count, spin_mult])
		"wild_collector_complete":
			var total = data.get("total_mult", 0.0)
			var bonus = data.get("global_bonus", 54.0)
			var secs = data.get("global_seconds", 108)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🃏✨ WILD COLLECTOR COMPLETE! %s: 10 spins! Total x%.1f! GLOBAL x%.1f for %ds!" % [trigger_name, total, bonus, secs],
				Color(1.0, 0.9, 0.0), 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(1.0, 0.9, 0.0), 8)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "🃏✨ WILD COLLECTOR!", "size": 22, "color": PANEL_COLOR},
				{"text": "Bonus Spins: 10", "size": 16, "color": Color.WHITE},
				{"text": "Stages: x2 → x3 → x10", "size": 16, "color": Color(1.0, 0.9, 0.0)},
				{"text": "Total x%.1f" % total, "size": 18, "color": PANEL_COLOR},
				{"text": "🌍 GLOBAL x%.1f for %ds" % [bonus, secs], "size": 18, "color": Color(1.0, 0.8, 0.0)},
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
