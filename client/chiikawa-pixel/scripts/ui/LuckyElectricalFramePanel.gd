## LuckyElectricalFramePanel.gd — T249 幸運電擊框架魚面板
## lucky-panel-agent 負責維護
## Catfish Hunters 電擊框架機制（Nolimit City 2026）
## 每次命中全局倍率翻倍（×1→×1024），最多 10 次
## 完美連鎖（≥8次）→ 全服 ×56.5 加成 113 秒
extends Node

const PANEL_COLOR = Color(0.0, 1.0, 1.0)  # 青色（電擊框架）

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
		"electrical_frame_start":
			BaseLuckyPanel.show_banner(_banner.panel,
				"⚡🎣 ELECTRICAL FRAME! %s activated Catfish Hunters! Global multiplier doubles each hit! Max x1024!" % trigger_name,
				PANEL_COLOR, 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 5)
			_show_indicator("ELECTRICAL FRAME", "Hit 0/10 ⚡")
		"frame_hit":
			var hit_no = data.get("hit_no", 1)
			var global_mult = data.get("global_mult", 2.0)
			_show_indicator("ELECTRICAL FRAME", "Hit %d x%.0f ⚡" % [hit_no, global_mult])
		"electrical_frame_complete":
			var hit_count = data.get("hit_count", 10)
			var final_mult = data.get("final_mult", 1024.0)
			var bonus = data.get("global_bonus", 56.5)
			var secs = data.get("global_seconds", 113)
			BaseLuckyPanel.show_banner(_banner.panel,
				"⚡🎣 ELECTRICAL FRAME COMPLETE! %s: %d hits, final x%.0f! GLOBAL x%.1f for %ds!" % [trigger_name, hit_count, final_mult, bonus, secs],
				Color(0.5, 1.0, 1.0), 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 6)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "⚡🎣 ELECTRICAL FRAME!", "size": 22, "color": PANEL_COLOR},
				{"text": "Catfish Hunters System", "size": 14, "color": Color(0.7, 1.0, 1.0)},
				{"text": "Hits: %d / Final Mult: x%.0f" % [hit_count, final_mult], "size": 16, "color": Color.WHITE},
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
