## LuckyDrillTorpedoPanel.gd — T113 幸運鑽頭魚雷魚面板
## lucky-panel-agent 負責維護
## 鑽頭主題：橙色 + 穿透計數器
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null
var _timer_bar: Dictionary = {}

func _ready() -> void:
	layer = 27
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)
	var ind_panel = _indicator.get("panel")
	if is_instance_valid(ind_panel):
		var bar_container = Control.new()
		bar_container.position = Vector2(8, 50)
		ind_panel.add_child(bar_container)
		_timer_bar = BaseLucky.create_timer_bar(bar_container, 184.0, Color(1.0, 0.55, 0.15))

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.55, 0.15), 2)
			BaseLucky.show_banner(_banner, "🚀 %s 發射鑽頭魚雷！穿透最多 5 個目標！" % name, Color(1.0, 0.55, 0.15), 2.5)
		"penetrate":
			var cnt = data.get("penetrate_cnt", 0)
			var mult = data.get("accum_mult", 1.0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "🚀 鑽頭魚雷"
				if is_instance_valid(value): value.text = "穿透 %d/5 | ×%.1f" % [cnt, mult]
			BaseLucky.update_timer_bar(_timer_bar, float(cnt) / 5.0)
			BaseLucky.spawn_float_text(self, Vector2(640, 300), "穿透！×%.1f" % mult, Color(1.0, 0.55, 0.15), 20)
		"explode":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.4, 0.1), 2)
			BaseLucky.show_banner(_banner, "💥 魚雷終點爆炸！AOE 傷害！", Color(1.0, 0.4, 0.1), 1.5)
		"perfect":
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.85, 0.0), 3)
			BaseLucky.show_banner(_banner, "🚀💥 完美穿透！%s 全服 ×2.2 加成 6 秒！" % name, Color(1.0, 0.85, 0.0), 3.5)
		"perfect_end":
			BaseLucky.show_banner(_banner, "🚀 完美穿透加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"settle":
			var reward = data.get("total_reward", 0)
			var cnt = data.get("penetrate_cnt", 0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "🚀 鑽頭魚雷結算", "size": 18, "color": Color(1.0, 0.55, 0.15)},
				{"text": "穿透：%d 個目標" % cnt, "size": 20, "color": Color.WHITE},
				{"text": "獎勵：+%d" % reward, "size": 22, "color": Color(1.0, 0.85, 0.0)},
			], 3.0)
