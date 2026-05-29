## LuckyCosmicMiraclePanel.gd — T237 幸運宇宙奇蹟魚面板
## lucky-panel-agent 負責維護
## Cosmic Miracle：全場 HP 歸零（每個目標 ×120.0）+ 8 道宇宙光柱
## 命中 ≥12 個 → 完美奇蹟，全服 ×49.5 加成 99 秒（超越 T236 的 ×49.0）
extends Node

const PANEL_COLOR = Color(0.58, 0.0, 0.83)  # 深紫色

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
		"miracle_start":
			_pillar_count = 0
			var per_mult = data.get("per_mult", 120.0)
			var pillars = data.get("pillar_count", 8)
			BaseLuckyPanel.show_banner(_banner.panel,
				"✨ COSMIC MIRACLE! %s summoned %d cosmic pillars! Every target x%.0f!" % [trigger_name, pillars, per_mult],
				PANEL_COLOR, 3.5)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 4)
			_show_indicator("COSMIC MIRACLE", "x%.0f/target" % per_mult)
		"miracle_pillar":
			_pillar_count = data.get("pillar_no", _pillar_count + 1)
			_show_indicator("COSMIC MIRACLE", "Pillar %d/8" % _pillar_count)
		"miracle_perfect":
			var targets = data.get("target_count", 0)
			var total = data.get("total_mult", 0.0)
			var bonus = data.get("global_bonus", 49.5)
			var secs = data.get("global_seconds", 99)
			BaseLuckyPanel.show_banner(_banner.panel,
				"✨ PERFECT COSMIC MIRACLE! %d targets cleared! Total x%.1f! Global x%.1f for %ds!" % [targets, total, bonus, secs],
				Color(0.8, 0.0, 1.0), 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(0.8, 0.0, 1.0), 5)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "✨ PERFECT COSMIC MIRACLE!", "size": 20, "color": PANEL_COLOR},
				{"text": "Targets Cleared: %d" % targets, "size": 16, "color": Color.WHITE},
				{"text": "Total x%.1f" % total, "size": 18, "color": PANEL_COLOR},
				{"text": "Global x%.1f for %ds" % [bonus, secs], "size": 16, "color": Color(1.0, 0.8, 0.0)},
			], 4.0)
			_hide_indicator()
		"miracle_end":
			var targets = data.get("target_count", 0)
			var total = data.get("total_mult", 0.0)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "✨ Cosmic Miracle End", "size": 18, "color": PANEL_COLOR},
				{"text": "Targets: %d  Total: x%.1f" % [targets, total], "size": 16, "color": Color.WHITE},
			], 2.5)
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
