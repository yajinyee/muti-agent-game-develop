## LuckyGenesisUltimatePanel.gd — T238 幸運創世終極魚面板
## lucky-panel-agent 負責維護
## Genesis Ultimate（里程碑：全服 ×50.0）
## 全場清空（每個目標 ×150.0）+ 12 道創世光柱，全服 ×50.0 加成 100 秒
## 史上第一個全服 ×50.0 里程碑！
extends Node

const PANEL_COLOR = Color(1.0, 0.84, 0.0)  # 純金色

var _canvas: CanvasLayer
var _banner: Dictionary
var _indicator: Dictionary
var _settle: Control
var _pillar_count: int = 0

func _ready() -> void:
	_canvas = get_tree().get_root().get_node_or_null("Main/HUDLayer")
	if not is_instance_valid(_canvas):
		_canvas = CanvasLayer.new()
		_canvas.layer = 10
		get_tree().get_root().add_child(_canvas)
	_banner   = BaseLuckyPanel.create_banner(_canvas, 100.0, 62)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 160), 67)
	_settle   = BaseLuckyPanel.create_settle_popup(_canvas, 72)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"genesis_start":
			_pillar_count = 0
			var per_mult = data.get("per_mult", 150.0)
			var pillars = data.get("pillar_count", 12)
			var global_target = data.get("global_target", 50.0)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🌟 GENESIS ULTIMATE! %s triggered ULTIMATE CREATION! %d pillars! x%.0f/target! MILESTONE: Global x%.1f!" % [trigger_name, pillars, per_mult, global_target],
				PANEL_COLOR, 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 6)
			_show_indicator("GENESIS ULTIMATE", "x%.0f/target" % per_mult)
		"genesis_pillar":
			_pillar_count = data.get("pillar_no", _pillar_count + 1)
			_show_indicator("GENESIS ULTIMATE", "Pillar %d/12" % _pillar_count)
		"genesis_milestone":
			var targets = data.get("target_count", 0)
			var total = data.get("total_mult", 0.0)
			var bonus = data.get("global_bonus", 50.0)
			var secs = data.get("global_seconds", 100)
			var milestone = data.get("milestone", "")
			BaseLuckyPanel.show_banner(_banner.panel,
				"🌟 GENESIS ULTIMATE MILESTONE! %d targets! Total x%.1f! GLOBAL x%.1f for %ds! HISTORY MADE!" % [targets, total, bonus, secs],
				Color(1.0, 1.0, 0.0), 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(1.0, 1.0, 0.0), 8)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "🌟 GENESIS ULTIMATE!", "size": 22, "color": PANEL_COLOR},
				{"text": "Targets Cleared: %d" % targets, "size": 16, "color": Color.WHITE},
				{"text": "Total x%.1f" % total, "size": 18, "color": PANEL_COLOR},
				{"text": "🏆 GLOBAL x%.1f for %ds" % [bonus, secs], "size": 18, "color": Color(1.0, 1.0, 0.0)},
				{"text": "MILESTONE: FIRST EVER x50.0!", "size": 14, "color": Color(1.0, 0.5, 0.0)},
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
