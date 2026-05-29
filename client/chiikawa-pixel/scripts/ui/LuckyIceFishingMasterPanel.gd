## LuckyIceFishingMasterPanel.gd — T236 幸運冰釣大師魚面板
## lucky-panel-agent 負責維護
## Ice Fishing Master：5 次旋轉（每次最高 ×8000），最高單次 ≥3000 → 完美冰釣
## 完美觸發 → 全服 ×49.0 加成 98 秒（超越 T235 的 ×48.5）
extends Node

const PANEL_COLOR = Color(0.0, 0.75, 1.0)  # 冰藍色

var _canvas: CanvasLayer
var _banner: Dictionary
var _indicator: Dictionary
var _settle: Control
var _spin_count: int = 0

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
		"ice_start":
			_spin_count = 0
			var max_mult = data.get("max_mult", 8000)
			BaseLuckyPanel.show_banner(_banner.panel,
				"❄️ ICE FISHING MASTER! %s activated! 5 spins, max x%d!" % [trigger_name, max_mult],
				PANEL_COLOR, 3.5)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 3)
			_show_indicator("ICE FISHING", "Spin 0/5")
		"ice_spin":
			_spin_count = data.get("spin_no", _spin_count + 1)
			var spin_mult = data.get("spin_mult", 0.0)
			var total = data.get("total_reward", 0)
			_show_indicator("ICE FISHING", "Spin %d: x%.0f" % [_spin_count, spin_mult])
			if spin_mult >= 1000:
				BaseLuckyPanel.show_banner(_banner.panel,
					"❄️ BIG SPIN! x%.0f!" % spin_mult,
					Color(0.0, 1.0, 1.0), 2.0)
		"ice_perfect":
			var max_spin = data.get("max_spin", 0.0)
			var total = data.get("total_reward", 0)
			var bonus = data.get("global_bonus", 49.0)
			var secs = data.get("global_seconds", 98)
			BaseLuckyPanel.show_banner(_banner.panel,
				"❄️ PERFECT ICE FISHING! Max spin x%.0f! Global x%.1f for %ds!" % [max_spin, bonus, secs],
				Color(0.0, 1.0, 1.0), 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(0.0, 1.0, 1.0), 5)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "❄️ PERFECT ICE FISHING!", "size": 20, "color": PANEL_COLOR},
				{"text": "Max Spin: x%.0f" % max_spin, "size": 18, "color": Color.WHITE},
				{"text": "Total Reward: %d" % total, "size": 16, "color": PANEL_COLOR},
				{"text": "Global x%.1f for %ds" % [bonus, secs], "size": 16, "color": Color(1.0, 0.8, 0.0)},
			], 4.0)
			_hide_indicator()
		"ice_end":
			var max_spin = data.get("max_spin", 0.0)
			var total = data.get("total_reward", 0)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "❄️ Ice Fishing End", "size": 18, "color": PANEL_COLOR},
				{"text": "Max: x%.0f  Total: %d" % [max_spin, total], "size": 16, "color": Color.WHITE},
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
