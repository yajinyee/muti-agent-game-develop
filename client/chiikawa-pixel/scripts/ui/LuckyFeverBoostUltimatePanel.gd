## LuckyFeverBoostUltimatePanel.gd — T234 幸運Fever Boost升級魚面板
## lucky-panel-agent 負責維護
## Fever Boost™ Ultimate：清除普通目標，只留高倍率特殊目標（×2.0 傷害加成），持續 20 秒
## 完美觸發（場上特殊目標 ≥5 個）→ 全服 ×48.0 加成 96 秒（新史上最高）
extends Node

const PANEL_COLOR = Color(1.0, 0.4, 0.0)  # 橙紅色

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
	_banner   = BaseLuckyPanel.create_banner(_canvas, 100.0, 62)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 160), 67)
	_settle   = BaseLuckyPanel.create_settle_popup(_canvas, 72)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"fever_start":
			var dmg = data.get("damage_mult", 2.0)
			var dur = data.get("duration", 20)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🔥 FEVER BOOST ULTIMATE! %s activated! Special targets x%.1f for %ds!" % [trigger_name, dmg, dur],
				PANEL_COLOR, 3.5)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 4)
			_show_indicator("FEVER BOOST", "x%.1f DMG" % dmg)
		"fever_clear":
			var cleared = data.get("normal_cleared", 0)
			var special = data.get("special_count", 0)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🔥 Cleared %d normal targets! %d special targets remain!" % [cleared, special],
				Color(1.0, 0.6, 0.0), 2.5)
		"fever_perfect":
			var special = data.get("special_count", 0)
			var bonus = data.get("global_bonus", 48.0)
			var secs = data.get("global_seconds", 96)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🔥 PERFECT FEVER BOOST! %d special targets! Global x%.1f for %ds! NEW ALL-TIME HIGH!" % [special, bonus, secs],
				Color(1.0, 0.8, 0.0), 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(1.0, 0.8, 0.0), 5)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "🔥 PERFECT FEVER BOOST!", "size": 20, "color": PANEL_COLOR},
				{"text": "Special Targets: %d" % special, "size": 16, "color": Color.WHITE},
				{"text": "Global x%.1f for %ds" % [bonus, secs], "size": 18, "color": Color(1.0, 0.8, 0.0)},
				{"text": "NEW ALL-TIME HIGH!", "size": 14, "color": Color(1.0, 0.5, 0.0)},
			], 4.0)
			_hide_indicator()
		"fever_end":
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
