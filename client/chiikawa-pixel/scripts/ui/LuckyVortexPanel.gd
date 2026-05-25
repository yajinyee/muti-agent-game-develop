## LuckyVortexPanel.gd — T108 幸運渦旋海葵面板
## lucky-panel-agent 負責維護
## 渦旋主題：紫色 + 旋轉特效
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null
var _timer_bar: Dictionary = {}

func _ready() -> void:
	layer = 22
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			BaseLucky.fullscreen_flash(self, Color(0.7, 0.3, 1.0), 2)
			BaseLucky.show_banner(_banner, "🌀 %s 召喚渦旋海葵！全場吸引 5 秒！" % name, Color(0.7, 0.4, 1.0), 2.5)
		"pull":
			var tl = data.get("time_left", 0.0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "🌀 渦旋海葵"
				if is_instance_valid(value): value.text = "剩餘 %.0fs" % tl
		"end":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			BaseLucky.fullscreen_flash(self, Color(0.7, 0.3, 1.0), 3)
			BaseLucky.show_banner(_banner, "🌀 渦旋爆炸！全場 HP -20%！", Color(0.8, 0.5, 1.0), 2.0)
		"settle":
			var reward = data.get("total_reward", 0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "🌀 渦旋海葵結算", "size": 18, "color": Color(0.7, 0.4, 1.0)},
				{"text": "獎勵：+%d" % reward, "size": 22, "color": Color(1.0, 0.85, 0.0)},
			], 3.0)
