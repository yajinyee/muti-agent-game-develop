## LuckyFishingTimeWheelPanel.gd — T242 幸運釣魚時間魚面板
## lucky-panel-agent 負責維護
## Fishing Time Wheel 機制（業界依據：BGaming Fishing Time 2026-04）
## 命運輪盤 + 倍率疊加，5 次旋轉（最高 ×10000），全服 ×52.5 加成 105 秒
extends Node

const PANEL_COLOR = Color(1.0, 0.75, 0.0)  # 金橙色

var _canvas: CanvasLayer
var _banner: Dictionary
var _indicator: Dictionary
var _settle: Control
var _spin_count: int = 0
var _best_spin: float = 0.0
var _total_so_far: float = 0.0

func _ready() -> void:
	_canvas = get_tree().get_root().get_node_or_null("Main/HUDLayer")
	if not is_instance_valid(_canvas):
		_canvas = CanvasLayer.new()
		_canvas.layer = 10
		get_tree().get_root().add_child(_canvas)
	_banner   = BaseLuckyPanel.create_banner(_canvas, 100.0, 66)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 160), 71)
	_settle   = BaseLuckyPanel.create_settle_popup(_canvas, 76)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"fishing_time_start":
			_spin_count = 0
			_best_spin = 0.0
			_total_so_far = 0.0
			var spins = data.get("spins", 5)
			var max_mult = data.get("max_mult", 10000.0)
			var global_target = data.get("global_target", 52.5)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🎣🎡 FISHING TIME WHEEL! %s triggered! %d spins! Max x%.0f! Global x%.1f!" % [trigger_name, spins, max_mult, global_target],
				PANEL_COLOR, 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 6)
			_show_indicator("FISHING TIME", "Spin 0/5")
		"wheel_spin":
			_spin_count = data.get("spin_no", _spin_count + 1)
			var spin_mult = data.get("spin_mult", 0.0)
			_total_so_far = data.get("total_so_far", _total_so_far + spin_mult)
			if spin_mult > _best_spin:
				_best_spin = spin_mult
			_show_indicator("FISHING TIME", "Spin %d: x%.0f (Total: x%.0f)" % [_spin_count, spin_mult, _total_so_far])
		"fishing_time_complete":
			var best = data.get("best_spin", _best_spin)
			var total = data.get("total_mult", _total_so_far)
			var bonus = data.get("global_bonus", 52.5)
			var secs = data.get("global_seconds", 105)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🎣🎡 FISHING TIME COMPLETE! Best x%.0f! Total x%.1f! Global x%.1f for %ds!" % [best, total, bonus, secs],
				Color(1.0, 1.0, 0.0), 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(1.0, 1.0, 0.0), 8)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "🎣🎡 FISHING TIME WHEEL!", "size": 22, "color": PANEL_COLOR},
				{"text": "Best Spin: x%.0f" % best, "size": 16, "color": Color(1.0, 0.9, 0.0)},
				{"text": "Total x%.1f" % total, "size": 18, "color": PANEL_COLOR},
				{"text": "🏆 GLOBAL x%.1f for %ds" % [bonus, secs], "size": 18, "color": Color(1.0, 1.0, 0.0)},
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
