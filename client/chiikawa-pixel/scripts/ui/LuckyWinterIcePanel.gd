## LuckyWinterIcePanel.gd — T240 幸運冬季冰釣魚面板
## lucky-panel-agent 負責維護
## Winter Ice Fishing 機制（業界依據：BGaming Winter Fishing Club 2026-01）
## 冰下魚群 + 53格輪盤，3 次旋轉（最高 ×500），全服 ×51.5 加成 103 秒
extends Node

const PANEL_COLOR = Color(0.53, 0.81, 0.98)  # 冰藍色

var _canvas: CanvasLayer
var _banner: Dictionary
var _indicator: Dictionary
var _settle: Control
var _spin_count: int = 0
var _best_spin: float = 0.0

func _ready() -> void:
	_canvas = get_tree().get_root().get_node_or_null("Main/HUDLayer")
	if not is_instance_valid(_canvas):
		_canvas = CanvasLayer.new()
		_canvas.layer = 10
		get_tree().get_root().add_child(_canvas)
	_banner   = BaseLuckyPanel.create_banner(_canvas, 100.0, 64)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 160), 69)
	_settle   = BaseLuckyPanel.create_settle_popup(_canvas, 74)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"winter_ice_start":
			_spin_count = 0
			_best_spin = 0.0
			var spins = data.get("spins", 3)
			var global_target = data.get("global_target", 51.5)
			BaseLuckyPanel.show_banner(_banner.panel,
				"❄️🎣 WINTER ICE FISHING! %s triggered! 53-segment wheel! %d spins! Max x500! Global x%.1f!" % [trigger_name, spins, global_target],
				PANEL_COLOR, 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 5)
			_show_indicator("WINTER ICE", "Spin 0/3")
		"wheel_spin":
			_spin_count = data.get("spin_no", _spin_count + 1)
			var segment = data.get("segment", "Leaf")
			var spin_mult = data.get("spin_mult", 1.0)
			if spin_mult > _best_spin:
				_best_spin = spin_mult
			_show_indicator("WINTER ICE", "Spin %d: %s x%.1f" % [_spin_count, segment, spin_mult])
		"winter_ice_complete":
			var best = data.get("best_spin", _best_spin)
			var total = data.get("total_mult", 0.0)
			var bonus = data.get("global_bonus", 51.5)
			var secs = data.get("global_seconds", 103)
			BaseLuckyPanel.show_banner(_banner.panel,
				"❄️🎣 WINTER ICE COMPLETE! Best x%.1f! Total x%.1f! Global x%.1f for %ds!" % [best, total, bonus, secs],
				Color(0.0, 0.9, 1.0), 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(0.0, 0.9, 1.0), 7)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "❄️🎣 WINTER ICE FISHING!", "size": 22, "color": PANEL_COLOR},
				{"text": "Best Spin: x%.1f" % best, "size": 16, "color": Color(1.0, 0.9, 0.0)},
				{"text": "Total x%.1f" % total, "size": 18, "color": PANEL_COLOR},
				{"text": "🏆 GLOBAL x%.1f for %ds" % [bonus, secs], "size": 18, "color": Color(0.0, 0.9, 1.0)},
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
