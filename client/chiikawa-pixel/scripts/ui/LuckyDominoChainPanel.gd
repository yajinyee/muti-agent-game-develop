## LuckyDominoChainPanel.gd — T246 幸運骨牌連鎖魚面板
## lucky-panel-agent 負責維護
## Domino Chain Reaction 機制（全新機制，里程碑：全服 ×55.0）
## 骨牌效應：最多 20 個骨牌依序倒下（每個 ×50.0），完美連鎖 → 全服 ×55.0 加成 110 秒
extends Node

const PANEL_COLOR = Color(1.0, 0.55, 0.0)  # 骨牌橙色

var _canvas: CanvasLayer
var _banner: Dictionary
var _indicator: Dictionary
var _settle: Control
var _domino_no: int = 0
var _domino_count: int = 0

func _ready() -> void:
	_canvas = get_tree().get_root().get_node_or_null("Main/HUDLayer")
	if not is_instance_valid(_canvas):
		_canvas = CanvasLayer.new()
		_canvas.layer = 10
		get_tree().get_root().add_child(_canvas)
	_banner    = BaseLuckyPanel.create_banner(_canvas, 100.0, 70)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 160), 75)
	_settle    = BaseLuckyPanel.create_settle_popup(_canvas, 80)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"domino_chain_start":
			_domino_no = 0
			_domino_count = data.get("domino_count", 10)
			var per_domino = data.get("per_domino", 50.0)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🀱💥 DOMINO CHAIN! %s started! %d dominoes falling! Each x%.0f!" % [trigger_name, _domino_count, per_domino],
				PANEL_COLOR, 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 5)
			_show_indicator("DOMINO CHAIN", "0/%d 🀱" % _domino_count)
		"domino_fall":
			_domino_no = data.get("domino_no", _domino_no + 1)
			var chain_mult = data.get("chain_mult", 50.0)
			_show_indicator("DOMINO CHAIN", "%d/%d x%.0f 🀱" % [_domino_no, _domino_count, chain_mult])
		"perfect_chain":
			var perfect_bonus = data.get("perfect_bonus", 0.0)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🀱💥 PERFECT CHAIN! %s: All dominoes fell! Bonus x%.1f!" % [trigger_name, perfect_bonus],
				Color(1.0, 0.8, 0.0), 3.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(1.0, 0.8, 0.0), 7)
		"domino_chain_complete":
			var total = data.get("total_mult", 0.0)
			var is_perfect = data.get("is_perfect", false)
			var bonus = data.get("global_bonus", 55.0)
			var secs = data.get("global_seconds", 110)
			var perfect_str = "🏆 PERFECT! " if is_perfect else ""
			BaseLuckyPanel.show_banner(_banner.panel,
				"🀱💥 DOMINO CHAIN COMPLETE! %s%s: %d dominoes! Total x%.1f! MILESTONE: GLOBAL x%.1f for %ds!" % [perfect_str, trigger_name, _domino_count, total, bonus, secs],
				Color(1.0, 0.7, 0.0), 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(1.0, 0.7, 0.0), 10)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "🀱💥 DOMINO CHAIN!", "size": 22, "color": PANEL_COLOR},
				{"text": "Dominoes: %d" % _domino_count, "size": 16, "color": Color.WHITE},
				{"text": "Per Domino: x50.0+", "size": 16, "color": Color(1.0, 0.9, 0.0)},
				{"text": "Perfect: %s" % ("YES 🏆" if is_perfect else "No"), "size": 16, "color": Color(0.0, 1.0, 0.5) if is_perfect else Color.WHITE},
				{"text": "Total x%.1f" % total, "size": 18, "color": PANEL_COLOR},
				{"text": "🏆 MILESTONE: GLOBAL x%.1f for %ds" % [bonus, secs], "size": 18, "color": Color(1.0, 0.6, 0.0)},
			], 6.0)
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
