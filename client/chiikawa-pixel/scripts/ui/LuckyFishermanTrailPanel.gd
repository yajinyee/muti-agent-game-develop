## LuckyFishermanTrailPanel.gd — T251 幸運漁夫路徑魚面板
## lucky-panel-agent 負責維護
## Bigger Bites 進階路徑機制（Reflex Gaming + Bragg 2026）
## Fishermen 符號收集 + 路徑升級（10 個節點，最高 ×500）
## 路徑完成（≥8節點）→ 全服 ×57.5 加成 115 秒
extends Node

const PANEL_COLOR = Color(1.0, 0.55, 0.0)  # 橙色（漁夫路徑）

var _canvas: CanvasLayer
var _banner: Dictionary
var _indicator: Dictionary
var _settle: Control

func _ready() -> void:
	_canvas = get_tree().get_root().get_node_or_null("Main/HUDLayer")
	if not is_instance_valid(_canvas):
		_canvas = CanvasLayer.new()
		_canvas.layer = 10
		get_tree().get_root().add_child(_canvas)
	_banner    = BaseLuckyPanel.create_banner(_canvas, 100.0, 73)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 170), 78)
	_settle    = BaseLuckyPanel.create_settle_popup(_canvas, 83)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"fisherman_trail_start":
			BaseLuckyPanel.show_banner(_banner.panel,
				"🎣🛤️ FISHERMAN TRAIL! %s activated Bigger Bites! 10 trail nodes, max x500!" % trigger_name,
				PANEL_COLOR, 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 5)
			_show_indicator("FISHERMAN TRAIL", "Node 0/10 🎣")
		"trail_node":
			var node_no = data.get("node_no", 1)
			var mult = data.get("mult", 10.0)
			var upgrade = data.get("upgrade", "Fish Upgrade")
			_show_indicator("FISHERMAN TRAIL", "Node %d x%.0f 🎣" % [node_no, mult])
			BaseLuckyPanel.spawn_float_text(
				_canvas, Vector2(640, 300),
				"%s! x%.0f" % [upgrade, mult],
				PANEL_COLOR, 18)
		"fisherman_trail_complete":
			var nodes_reached = data.get("nodes_reached", 10)
			var bonus = data.get("global_bonus", 57.5)
			var secs = data.get("global_seconds", 115)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🎣🛤️ FISHERMAN TRAIL COMPLETE! %s: %d nodes! GLOBAL x%.1f for %ds!" % [trigger_name, nodes_reached, bonus, secs],
				Color(1.0, 0.75, 0.3), 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 6)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "🎣🛤️ FISHERMAN TRAIL!", "size": 22, "color": PANEL_COLOR},
				{"text": "Bigger Bites System (max x500)", "size": 14, "color": Color(1.0, 0.75, 0.3)},
				{"text": "Nodes Reached: %d / 10" % nodes_reached, "size": 16, "color": Color.WHITE},
				{"text": "🌐 GLOBAL x%.1f for %ds" % [bonus, secs], "size": 18, "color": PANEL_COLOR},
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
