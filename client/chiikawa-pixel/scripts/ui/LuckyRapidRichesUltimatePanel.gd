## LuckyRapidRichesUltimatePanel.gd — T235 幸運快速暴富升級魚面板
## lucky-panel-agent 負責維護
## Rapid Riches Ultimate：3 秒極速連擊視窗（每次擊破 ×300.0），連擊 ≥10 次 → 完美暴富
## 完美觸發 → 全服 ×48.5 加成 97 秒（超越 T234 的 ×48.0）
extends Node

const PANEL_COLOR = Color(1.0, 0.85, 0.0)  # 金黃色

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
		"rapid_start":
			var per_mult = data.get("per_mult", 300.0)
			var window = data.get("window_secs", 3)
			BaseLuckyPanel.show_banner(_banner.panel,
				"⚡ RAPID RICHES ULTIMATE! %s activated! %ds combo window! Each kill x%.0f!" % [trigger_name, window, per_mult],
				PANEL_COLOR, 3.5)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 4)
			_show_indicator("RAPID RICHES", "x%.0f/kill" % per_mult)
		"rapid_perfect":
			var combo = data.get("combo_hits", 0)
			var total = data.get("total_mult", 0.0)
			var bonus = data.get("global_bonus", 48.5)
			var secs = data.get("global_seconds", 97)
			BaseLuckyPanel.show_banner(_banner.panel,
				"⚡ PERFECT RAPID RICHES! %s hit %d combos! Total x%.1f! Global x%.1f for %ds!" % [trigger_name, combo, total, bonus, secs],
				Color(1.0, 1.0, 0.0), 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(1.0, 1.0, 0.0), 5)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "⚡ PERFECT RAPID RICHES!", "size": 20, "color": PANEL_COLOR},
				{"text": "Combo Hits: %d" % combo, "size": 16, "color": Color.WHITE},
				{"text": "Total x%.1f" % total, "size": 18, "color": PANEL_COLOR},
				{"text": "Global x%.1f for %ds" % [bonus, secs], "size": 16, "color": Color(1.0, 0.8, 0.0)},
			], 4.0)
			_hide_indicator()
		"rapid_end":
			var combo = data.get("combo_hits", 0)
			var total = data.get("total_mult", 0.0)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "⚡ Rapid Riches End", "size": 18, "color": PANEL_COLOR},
				{"text": "Combo: %d  Total: x%.1f" % [combo, total], "size": 16, "color": Color.WHITE},
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
