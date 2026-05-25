## LuckyDeepWhirlpoolPanel.gd — T119 幸運深海漩渦魚面板
## lucky-panel-agent 負責維護
## 深海漩渦主題：青藍色 + 6秒持續傷害
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null
var _timer_bar: Dictionary = {}

func _ready() -> void:
	layer = 33
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)
	var ind_panel = _indicator.get("panel")
	if is_instance_valid(ind_panel):
		var bar_container = Control.new()
		bar_container.position = Vector2(8, 50)
		ind_panel.add_child(bar_container)
		_timer_bar = BaseLucky.create_timer_bar(bar_container, 184.0, Color(0.2, 0.7, 1.0))

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			BaseLucky.fullscreen_flash(self, Color(0.2, 0.7, 1.0), 2)
			BaseLucky.show_banner(_banner, "🌊🌀 %s 觸發深海漩渦！全場 HP -50%！6 秒！" % name, Color(0.2, 0.7, 1.0), 2.5)
		"whirlpool_damage":
			var hits = data.get("hit_count", 0)
			var tl = data.get("time_left", 0.0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "🌀 深海漩渦"
				if is_instance_valid(value): value.text = "%.0fs | 命中 %d" % [tl, hits]
			BaseLucky.update_timer_bar(_timer_bar, tl / 6.0)
			BaseLucky.spawn_float_text(self, Vector2(640, 300), "🌀 命中 %d！" % hits, Color(0.2, 0.7, 1.0), 20)
		"settle":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			var reward = data.get("total_reward", 0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "🌊 深海漩渦結算", "size": 18, "color": Color(0.2, 0.7, 1.0)},
				{"text": "全場 HP -50% 完成！", "size": 18, "color": Color.WHITE},
				{"text": "獎勵：+%d" % reward, "size": 22, "color": Color(1.0, 0.85, 0.0)},
			], 3.0)
