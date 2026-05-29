## LuckySharkSparkPanel.gd — T239 幸運鯊魚閃電魚面板
## lucky-panel-agent 負責維護
## Shark & Spark 機制（業界依據：BGaming Shark & Spark Hold & Win 2026-05-30 最新）
## 鯊魚閃電 + 珍珠倍率組合，閃電連鎖 6 條（每條 ×80.0），全服 ×51.0 加成 102 秒（里程碑）
extends Node

const PANEL_COLOR = Color(0.0, 0.75, 1.0)  # 深海藍

var _canvas: CanvasLayer
var _banner: Dictionary
var _indicator: Dictionary
var _settle: Control
var _chain_count: int = 0

func _ready() -> void:
	_canvas = get_tree().get_root().get_node_or_null("Main/HUDLayer")
	if not is_instance_valid(_canvas):
		_canvas = CanvasLayer.new()
		_canvas.layer = 10
		get_tree().get_root().add_child(_canvas)
	_banner   = BaseLuckyPanel.create_banner(_canvas, 100.0, 63)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 160), 68)
	_settle   = BaseLuckyPanel.create_settle_popup(_canvas, 73)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"shark_spark_start":
			_chain_count = 0
			var chains = data.get("chain_count", 6)
			var per_chain = data.get("per_chain", 80.0)
			var global_target = data.get("global_target", 51.0)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🦈⚡ SHARK & SPARK! %s triggered! %d chains x%.0f! Pearl Multipliers! MILESTONE: Global x%.1f!" % [trigger_name, chains, per_chain, global_target],
				PANEL_COLOR, 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 6)
			_show_indicator("SHARK & SPARK", "Chains: 0/6")
		"pearl_assigned":
			var count = data.get("pearl_count", 0)
			_show_indicator("SHARK & SPARK", "Pearls: %d" % count)
		"chain_strike":
			_chain_count = data.get("chain_no", _chain_count + 1)
			var chain_mult = data.get("chain_mult", 80.0)
			_show_indicator("SHARK & SPARK", "Chain %d/6 x%.0f" % [_chain_count, chain_mult])
		"shark_spark_complete":
			var chains = data.get("chain_count", 6)
			var total_chain = data.get("total_chain", 0.0)
			var total_pearl = data.get("total_pearl", 0.0)
			var total = data.get("total_mult", 0.0)
			var bonus = data.get("global_bonus", 51.0)
			var secs = data.get("global_seconds", 102)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🦈⚡ SHARK & SPARK COMPLETE! %d chains + Pearl x%.1f! Total x%.1f! MILESTONE: Global x%.1f for %ds!" % [chains, total_pearl, total, bonus, secs],
				Color(0.0, 1.0, 1.0), 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(0.0, 1.0, 1.0), 8)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "🦈⚡ SHARK & SPARK!", "size": 22, "color": PANEL_COLOR},
				{"text": "Chains: %d x%.0f" % [chains, total_chain / max(chains, 1)], "size": 16, "color": Color.WHITE},
				{"text": "Pearl Bonus: x%.1f" % total_pearl, "size": 16, "color": Color(1.0, 0.9, 0.0)},
				{"text": "Total x%.1f" % total, "size": 18, "color": PANEL_COLOR},
				{"text": "🏆 GLOBAL x%.1f for %ds" % [bonus, secs], "size": 18, "color": Color(0.0, 1.0, 1.0)},
				{"text": "MILESTONE: GLOBAL x51.0!", "size": 14, "color": Color(1.0, 0.5, 0.0)},
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
