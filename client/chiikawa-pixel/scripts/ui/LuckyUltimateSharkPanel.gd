## LuckyUltimateSharkPanel.gd — T243 幸運終極鯊魚魚面板
## lucky-panel-agent 負責維護
## Ultimate Shark 機制（里程碑：全服 ×53.0 新史上最高）
## 終極鯊魚清場（每個 ×180.0）+ 14 次鯊魚咬合，全服 ×53.0 加成 106 秒
extends Node

const PANEL_COLOR = Color(1.0, 0.27, 0.0)  # 鯊魚橙紅

var _canvas: CanvasLayer
var _banner: Dictionary
var _indicator: Dictionary
var _settle: Control
var _bite_count: int = 0

func _ready() -> void:
	_canvas = get_tree().get_root().get_node_or_null("Main/HUDLayer")
	if not is_instance_valid(_canvas):
		_canvas = CanvasLayer.new()
		_canvas.layer = 10
		get_tree().get_root().add_child(_canvas)
	_banner   = BaseLuckyPanel.create_banner(_canvas, 100.0, 67)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 160), 72)
	_settle   = BaseLuckyPanel.create_settle_popup(_canvas, 77)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"ultimate_shark_start":
			_bite_count = 0
			var per_mult = data.get("per_mult", 180.0)
			var bites = data.get("bite_count", 14)
			var global_target = data.get("global_target", 53.0)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🦈💥 ULTIMATE SHARK! %s triggered! Every target x%.0f! %d bites! MILESTONE: Global x%.1f!" % [trigger_name, per_mult, bites, global_target],
				PANEL_COLOR, 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 7)
			_show_indicator("ULTIMATE SHARK", "Bites: 0/14")
		"shark_bite":
			_bite_count = data.get("bite_no", _bite_count + 1)
			_show_indicator("ULTIMATE SHARK", "Bite %d/14 🦈" % _bite_count)
		"ultimate_shark_milestone":
			var targets = data.get("target_count", 0)
			var per_mult = data.get("per_mult", 180.0)
			var total = data.get("total_mult", 0.0)
			var bonus = data.get("global_bonus", 53.0)
			var secs = data.get("global_seconds", 106)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🦈💥 ULTIMATE SHARK MILESTONE! %d targets x%.0f! Total x%.1f! GLOBAL x%.1f for %ds! NEW RECORD!" % [targets, per_mult, total, bonus, secs],
				Color(1.0, 0.5, 0.0), 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(1.0, 0.5, 0.0), 10)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "🦈💥 ULTIMATE SHARK!", "size": 22, "color": PANEL_COLOR},
				{"text": "Targets Cleared: %d" % targets, "size": 16, "color": Color.WHITE},
				{"text": "Per Target: x%.0f" % per_mult, "size": 16, "color": Color(1.0, 0.9, 0.0)},
				{"text": "Total x%.1f" % total, "size": 18, "color": PANEL_COLOR},
				{"text": "🏆 GLOBAL x%.1f for %ds" % [bonus, secs], "size": 18, "color": Color(1.0, 0.5, 0.0)},
				{"text": "🌟 NEW RECORD: GLOBAL x53.0!", "size": 14, "color": Color(1.0, 1.0, 0.0)},
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
