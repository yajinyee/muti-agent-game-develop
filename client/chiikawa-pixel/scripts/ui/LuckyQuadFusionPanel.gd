## LuckyQuadFusionPanel.gd — T248 幸運四重終極融合魚面板
## lucky-panel-agent 負責維護
## Quad Fusion Ultimate 機制（里程碑：全服 ×56.0 新史上最高）
## 四重機制融合：Wild Collector + Lightning Eel + Domino Chain + Immortal Boss
## 全服 ×56.0 加成 112 秒（新史上最高，超越 T246 的 ×55.0）
extends Node

const PANEL_COLOR = Color(1.0, 0.0, 1.0)  # 洋紅色（四重融合）

var _canvas: CanvasLayer
var _banner: Dictionary
var _indicator: Dictionary
var _settle: Control
var _phase: int = 0

func _ready() -> void:
	_canvas = get_tree().get_root().get_node_or_null("Main/HUDLayer")
	if not is_instance_valid(_canvas):
		_canvas = CanvasLayer.new()
		_canvas.layer = 10
		get_tree().get_root().add_child(_canvas)
	_banner    = BaseLuckyPanel.create_banner(_canvas, 100.0, 72)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 160), 77)
	_settle    = BaseLuckyPanel.create_settle_popup(_canvas, 82)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"quad_fusion_start":
			_phase = 0
			BaseLuckyPanel.show_banner(_banner.panel,
				"💎🌌 QUAD FUSION ULTIMATE! %s activated 4-Phase Fusion! Wild+Eel+Domino+Boss! MILESTONE: GLOBAL x56.0!" % trigger_name,
				PANEL_COLOR, 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 8)
			_show_indicator("QUAD FUSION", "Phase 0/4 💎")
		"phase1_spin":
			_phase = 1
			var mult = data.get("mult", 2.0)
			var spin_no = data.get("spin_no", 1)
			_show_indicator("PHASE 1: WILD", "Spin %d x%.0f 🃏" % [spin_no, mult])
		"phase1_complete":
			BaseLuckyPanel.show_banner(_banner.panel,
				"💎 PHASE 1 COMPLETE! Wild Collector done! Moving to Phase 2: Lightning Eel!",
				Color(1.0, 0.85, 0.0), 2.5)
		"phase2_eel":
			_phase = 2
			var eel_no = data.get("eel_no", 1)
			var mult = data.get("mult", 90.0)
			_show_indicator("PHASE 2: EEL", "Eel %d x%.0f ⚡" % [eel_no, mult])
		"phase2_complete":
			BaseLuckyPanel.show_banner(_banner.panel,
				"💎 PHASE 2 COMPLETE! Lightning Eel done! Moving to Phase 3: Domino Chain!",
				Color(0.0, 1.0, 1.0), 2.5)
		"phase3_domino":
			_phase = 3
			var domino_no = data.get("domino_no", 1)
			var mult = data.get("mult", 50.0)
			_show_indicator("PHASE 3: DOMINO", "Domino %d x%.0f 🀱" % [domino_no, mult])
		"phase3_complete":
			BaseLuckyPanel.show_banner(_banner.panel,
				"💎 PHASE 3 COMPLETE! Domino Chain done! Moving to Phase 4: Immortal Boss!",
				Color(1.0, 0.55, 0.0), 2.5)
		"phase4_boss":
			_phase = 4
			var revive_no = data.get("revive_no", 1)
			var mult = data.get("mult", 100.0)
			_show_indicator("PHASE 4: BOSS", "Revive %d x%.0f 💀" % [revive_no, mult])
		"quad_fusion_milestone":
			var total = data.get("total_mult", 0.0)
			var bonus = data.get("global_bonus", 56.0)
			var secs = data.get("global_seconds", 112)
			BaseLuckyPanel.show_banner(_banner.panel,
				"💎🌌 QUAD FUSION MILESTONE! %s: 4 phases complete! Total x%.1f! MILESTONE: GLOBAL x%.1f for %ds! NEW RECORD!" % [trigger_name, total, bonus, secs],
				Color(1.0, 0.5, 1.0), 6.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(1.0, 0.5, 1.0), 12)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "💎🌌 QUAD FUSION ULTIMATE!", "size": 22, "color": PANEL_COLOR},
				{"text": "Phase 1: Wild Collector (x2→x3→x10)", "size": 13, "color": Color(1.0, 0.85, 0.0)},
				{"text": "Phase 2: Lightning Eel (8×x90)", "size": 13, "color": Color(0.0, 1.0, 1.0)},
				{"text": "Phase 3: Domino Chain (20×x50+)", "size": 13, "color": Color(1.0, 0.55, 0.0)},
				{"text": "Phase 4: Immortal Boss (x100→x300)", "size": 13, "color": Color(1.0, 0.2, 0.2)},
				{"text": "Total x%.1f" % total, "size": 18, "color": PANEL_COLOR},
				{"text": "🏆 MILESTONE: GLOBAL x%.1f for %ds" % [bonus, secs], "size": 18, "color": Color(1.0, 0.5, 1.0)},
				{"text": "🌟 NEW RECORD: GLOBAL x56.0!", "size": 14, "color": Color(1.0, 1.0, 0.0)},
			], 7.0)
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
