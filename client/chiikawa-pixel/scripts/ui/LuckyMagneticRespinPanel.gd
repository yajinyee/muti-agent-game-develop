## LuckyMagneticRespinPanel.gd — T250 幸運磁力連鎖魚面板
## lucky-panel-agent 負責維護
## Golden Gills 磁力連鎖 Respin 機制（Atomic Slot Lab 2026）
## 磁力連鎖 Respin：8 次，每次 ×75.0
## 完美連鎖（≥6次）→ 全服 ×57.0 加成 114 秒
extends Node

const PANEL_COLOR = Color(1.0, 0.84, 0.0)  # 黃金色（磁力連鎖）

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
		"magnetic_respin_start":
			BaseLuckyPanel.show_banner(_banner.panel,
				"🧲✨ MAGNETIC RESPIN! %s activated Golden Gills! 8 respins with x75 multiplier!" % trigger_name,
				PANEL_COLOR, 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 5)
			_show_indicator("MAGNETIC RESPIN", "Respin 0/8 🧲")
		"respin":
			var respin_no = data.get("respin_no", 1)
			var attracted = data.get("attracted", 1)
			_show_indicator("MAGNETIC RESPIN", "Respin %d (%d targets) 🧲" % [respin_no, attracted])
		"magnetic_respin_complete":
			var respin_count = data.get("respin_count", 8)
			var total_targets = data.get("total_targets", 0)
			var bonus = data.get("global_bonus", 57.0)
			var secs = data.get("global_seconds", 114)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🧲✨ MAGNETIC RESPIN COMPLETE! %s: %d respins, %d targets! GLOBAL x%.1f for %ds!" % [trigger_name, respin_count, total_targets, bonus, secs],
				Color(1.0, 0.95, 0.5), 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 6)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "🧲✨ MAGNETIC RESPIN!", "size": 22, "color": PANEL_COLOR},
				{"text": "Golden Gills System (x75 per respin)", "size": 14, "color": Color(1.0, 0.95, 0.5)},
				{"text": "Respins: %d / Targets: %d" % [respin_count, total_targets], "size": 16, "color": Color.WHITE},
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
