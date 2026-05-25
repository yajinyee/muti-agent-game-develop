## LuckyThunderStormPanel.gd — T124 幸運雷暴魚面板
## lucky-panel-agent 負責維護
## 雷暴主題：黃色 + 閃電計數器
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null
var _timer_bar: Dictionary = {}

func _ready() -> void:
	layer = 38
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)
	var ind_panel = _indicator.get("panel")
	if is_instance_valid(ind_panel):
		var bar_container = Control.new()
		bar_container.position = Vector2(8, 50)
		ind_panel.add_child(bar_container)
		_timer_bar = BaseLucky.create_timer_bar(bar_container, 184.0, Color(1.0, 0.9, 0.0))

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.9, 0.0), 2)
			BaseLucky.show_banner(_banner, "⛈️ %s 觸發雷暴！10 秒 6-7 道閃電！" % name, Color(1.0, 0.9, 0.0), 2.5)
		"lightning":
			var no = data.get("lightning_no", 1)
			var hits = data.get("hit_count", 0)
			var tl = data.get("time_left", 0.0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "⛈️ 雷暴"
				if is_instance_valid(value): value.text = "閃電 %d | %.0fs" % [no, tl]
			BaseLucky.update_timer_bar(_timer_bar, tl / 10.0)
			BaseLucky.spawn_float_text(self, Vector2(640, 300), "⚡ 閃電 %d！命中 %d！" % [no, hits], Color(1.0, 0.9, 0.0), 20)
		"thunder_perfect":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.9, 0.0), 3)
			BaseLucky.show_banner(_banner, "⛈️✨ 雷暴完美！%s 全服 ×2.3 加成 6 秒！" % name, Color(1.0, 0.9, 0.0), 3.5)
		"thunder_perfect_end":
			BaseLucky.show_banner(_banner, "⛈️ 雷暴完美加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"settle":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			var reward = data.get("total_reward", 0)
			var hits = data.get("total_hits", 0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "⛈️ 雷暴結算", "size": 18, "color": Color(1.0, 0.9, 0.0)},
				{"text": "命中：%d 個目標" % hits, "size": 20, "color": Color.WHITE},
				{"text": "獎勵：+%d" % reward, "size": 22, "color": Color(1.0, 0.85, 0.0)},
			], 3.0)
