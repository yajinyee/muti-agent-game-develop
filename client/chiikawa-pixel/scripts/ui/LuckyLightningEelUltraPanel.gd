## LuckyLightningEelUltraPanel.gd — T245 幸運閃電鰻升級魚面板
## lucky-panel-agent 負責維護
## Lightning Eel Ultra 機制（Royal Fishing 60x Lightning Eel 升級版）
## 8 條鰻魚依序觸發（每條 ×90.0），連鎖跳躍 3 次，全服 ×54.5 加成 109 秒
extends Node

const PANEL_COLOR = Color(0.0, 1.0, 1.0)  # 青色閃電

var _canvas: CanvasLayer
var _banner: Dictionary
var _indicator: Dictionary
var _settle: Control
var _eel_no: int = 0
var _jump_no: int = 0

func _ready() -> void:
	_canvas = get_tree().get_root().get_node_or_null("Main/HUDLayer")
	if not is_instance_valid(_canvas):
		_canvas = CanvasLayer.new()
		_canvas.layer = 10
		get_tree().get_root().add_child(_canvas)
	_banner    = BaseLuckyPanel.create_banner(_canvas, 100.0, 69)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 160), 74)
	_settle    = BaseLuckyPanel.create_settle_popup(_canvas, 79)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"lightning_eel_ultra_start":
			_eel_no = 0
			_jump_no = 0
			var eel_count = data.get("eel_count", 8)
			var per_eel = data.get("per_eel", 90.0)
			BaseLuckyPanel.show_banner(_banner.panel,
				"⚡🐍 LIGHTNING EEL ULTRA! %s unleashed %d eels! Each x%.0f! Chain jumps x3!" % [trigger_name, eel_count, per_eel],
				PANEL_COLOR, 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 6)
			_show_indicator("LIGHTNING EEL", "Eel 0/8 ⚡")
		"eel_chain_jump":
			_eel_no = data.get("eel_no", _eel_no)
			_jump_no = data.get("jump_no", _jump_no)
			var jump_mult = data.get("jump_mult", 36.0)
			_show_indicator("LIGHTNING EEL", "Eel %d Jump %d x%.0f ⚡" % [_eel_no, _jump_no, jump_mult])
		"lightning_eel_ultra_complete":
			var total = data.get("total_mult", 0.0)
			var perfect = data.get("perfect_bonus", 0.0)
			var bonus = data.get("global_bonus", 54.5)
			var secs = data.get("global_seconds", 109)
			BaseLuckyPanel.show_banner(_banner.panel,
				"⚡🐍 LIGHTNING EEL ULTRA COMPLETE! %s: 8 eels! Total x%.1f! GLOBAL x%.1f for %ds!" % [trigger_name, total, bonus, secs],
				Color(0.0, 0.9, 1.0), 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(0.0, 0.9, 1.0), 9)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "⚡🐍 LIGHTNING EEL ULTRA!", "size": 22, "color": PANEL_COLOR},
				{"text": "8 Eels × 3 Jumps", "size": 16, "color": Color.WHITE},
				{"text": "Per Eel: x90.0", "size": 16, "color": Color(1.0, 1.0, 0.0)},
				{"text": "Perfect Bonus: x%.1f" % perfect, "size": 16, "color": Color(0.0, 1.0, 0.8)},
				{"text": "Total x%.1f" % total, "size": 18, "color": PANEL_COLOR},
				{"text": "🌍 GLOBAL x%.1f for %ds" % [bonus, secs], "size": 18, "color": Color(0.0, 0.8, 1.0)},
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
