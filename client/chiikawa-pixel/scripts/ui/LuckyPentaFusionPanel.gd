## LuckyPentaFusionPanel.gd — T253 幸運五重終極魚面板
## lucky-panel-agent 負責維護
## Penta Fusion Ultimate 機制（里程碑：全服 ×58.5 新史上最高）
## 五重機制融合：電擊框架 + 磁力連鎖 + 漁夫路徑 + 黃金鰓 Jackpot + Quad Fusion
## 全服 ×58.5 加成 117 秒（新史上最高，超越 T248 的 ×56.0）
extends Node

const PANEL_COLOR = Color(1.0, 0.41, 0.71)  # 熱粉紅色（五重融合）

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
	_banner    = BaseLuckyPanel.create_banner(_canvas, 100.0, 74)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 180), 79)
	_settle    = BaseLuckyPanel.create_settle_popup(_canvas, 84)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"penta_fusion_start":
			_phase = 0
			BaseLuckyPanel.show_banner(_banner.panel,
				"💎🌸 PENTA FUSION ULTIMATE! %s activated 5-Phase Fusion! Frame+Magnetic+Trail+Gills+Quad! MILESTONE: GLOBAL x58.5!" % trigger_name,
				PANEL_COLOR, 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 10)
			_show_indicator("PENTA FUSION", "Phase 0/5 💎")
		"phase1_frame":
			_phase = 1
			var hit_no = data.get("hit_no", 1)
			var global_mult = data.get("global_mult", 2.0)
			_show_indicator("PHASE 1: FRAME", "Hit %d x%.0f ⚡" % [hit_no, global_mult])
		"phase1_complete":
			BaseLuckyPanel.show_banner(_banner.panel,
				"💎 PHASE 1 COMPLETE! Electrical Frame done! Moving to Phase 2: Magnetic Respin!",
				Color(0.0, 1.0, 1.0), 2.5)
		"phase2_respin":
			_phase = 2
			var respin_no = data.get("respin_no", 1)
			_show_indicator("PHASE 2: MAGNETIC", "Respin %d x75 🧲" % respin_no)
		"phase2_complete":
			BaseLuckyPanel.show_banner(_banner.panel,
				"💎 PHASE 2 COMPLETE! Magnetic Respin done! Moving to Phase 3: Fisherman Trail!",
				Color(1.0, 0.84, 0.0), 2.5)
		"phase3_node":
			_phase = 3
			var node_no = data.get("node_no", 1)
			var mult = data.get("mult", 10.0)
			_show_indicator("PHASE 3: TRAIL", "Node %d x%.0f 🎣" % [node_no, mult])
		"phase3_complete":
			BaseLuckyPanel.show_banner(_banner.panel,
				"💎 PHASE 3 COMPLETE! Fisherman Trail done! Moving to Phase 4: Golden Gills!",
				Color(1.0, 0.55, 0.0), 2.5)
		"phase4_grand_jackpot":
			_phase = 4
			var mult = data.get("mult", 2000.0)
			_show_indicator("PHASE 4: JACKPOT", "GRAND x%.0f 💰" % mult)
			BaseLuckyPanel.show_banner(_banner.panel,
				"💰 PHASE 4: GRAND JACKPOT! x%.0f!" % mult,
				Color(1.0, 0.84, 0.0), 2.5)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(1.0, 0.84, 0.0), 5)
		"phase4_complete":
			BaseLuckyPanel.show_banner(_banner.panel,
				"💎 PHASE 4 COMPLETE! Golden Gills done! Moving to Phase 5: Quad Fusion!",
				Color(1.0, 0.84, 0.0), 2.5)
		"phase5_clear":
			_phase = 5
			var clear_no = data.get("clear_no", 1)
			var mult = data.get("mult", 200.0)
			_show_indicator("PHASE 5: QUAD", "Clear %d x%.0f 🌌" % [clear_no, mult])
		"penta_fusion_milestone":
			var total = data.get("total_mult", 0.0)
			var bonus = data.get("global_bonus", 58.5)
			var secs = data.get("global_seconds", 117)
			BaseLuckyPanel.show_banner(_banner.panel,
				"💎🌸 PENTA FUSION MILESTONE! %s: 5 phases complete! Total x%.1f! MILESTONE: GLOBAL x%.1f for %ds! NEW RECORD!" % [trigger_name, total, bonus, secs],
				Color(1.0, 0.7, 0.9), 6.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(1.0, 0.7, 0.9), 15)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "💎🌸 PENTA FUSION ULTIMATE!", "size": 22, "color": PANEL_COLOR},
				{"text": "Phase 1: Electrical Frame (x1→x1024)", "size": 12, "color": Color(0.0, 1.0, 1.0)},
				{"text": "Phase 2: Magnetic Respin (8×x75)", "size": 12, "color": Color(1.0, 0.84, 0.0)},
				{"text": "Phase 3: Fisherman Trail (10 nodes)", "size": 12, "color": Color(1.0, 0.55, 0.0)},
				{"text": "Phase 4: Golden Gills Grand Jackpot", "size": 12, "color": Color(1.0, 0.84, 0.0)},
				{"text": "Phase 5: Quad Fusion (5×x200)", "size": 12, "color": Color(1.0, 0.0, 1.0)},
				{"text": "Total x%.1f" % total, "size": 18, "color": PANEL_COLOR},
				{"text": "🏆 MILESTONE: GLOBAL x%.1f for %ds" % [bonus, secs], "size": 18, "color": Color(1.0, 0.7, 0.9)},
				{"text": "🌟 NEW RECORD: GLOBAL x58.5!", "size": 14, "color": Color(1.0, 1.0, 0.0)},
			], 8.0)
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
