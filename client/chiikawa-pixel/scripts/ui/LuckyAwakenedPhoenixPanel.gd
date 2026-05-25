## LuckyAwakenedPhoenixPanel.gd — T111 幸運覺醒鳳凰魚面板
## lucky-panel-agent 負責維護
## 鳳凰主題：火橙色 + Power Up 計數器
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null
var _timer_bar: Dictionary = {}

func _ready() -> void:
	layer = 25
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)
	var ind_panel = _indicator.get("panel")
	if is_instance_valid(ind_panel):
		var bar_container = Control.new()
		bar_container.position = Vector2(8, 50)
		ind_panel.add_child(bar_container)
		_timer_bar = BaseLucky.create_timer_bar(bar_container, 184.0, Color(1.0, 0.6, 0.2))

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"awaken_start":
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.6, 0.2), 2)
			BaseLucky.show_banner(_banner, "🔥 %s 觸發覺醒鳳凰！下 5 次攻擊 Power Up！" % name, Color(1.0, 0.6, 0.2), 2.5)
		"power_up":
			var mult = data.get("power_up_mult", 6.0)
			var shots = data.get("shots_left", 0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "🔥 覺醒鳳凰"
				if is_instance_valid(value): value.text = "×%.0f | 剩 %d 次" % [mult, shots]
			BaseLucky.update_timer_bar(_timer_bar, float(shots) / 5.0)
			BaseLucky.spawn_float_text(self, Vector2(640, 300), "🔥 Power Up ×%.0f！" % mult, Color(1.0, 0.6, 0.2), 22)
		"perfect_awaken":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.85, 0.0), 3)
			BaseLucky.show_banner(_banner, "🔥✨ 完美覺醒！%s 全服 ×2.0 加成 8 秒！" % name, Color(1.0, 0.85, 0.0), 3.5)
		"perfect_end":
			BaseLucky.show_banner(_banner, "🔥 完美覺醒加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"awaken_end":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			var reward = data.get("total_reward", 0)
			var hits = data.get("hit_count", 0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "🔥 覺醒鳳凰結算", "size": 18, "color": Color(1.0, 0.6, 0.2)},
				{"text": "命中：%d 次" % hits, "size": 20, "color": Color.WHITE},
				{"text": "獎勵：+%d" % reward, "size": 22, "color": Color(1.0, 0.85, 0.0)},
			], 3.0)
