## LuckyRocketCannonPanel.gd — T118 幸運火箭砲魚面板
## lucky-panel-agent 負責維護
## 火箭砲主題：橙紅色 + 3枚火箭序列
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null

func _ready() -> void:
	layer = 32
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.5, 0.2), 2)
			BaseLucky.show_banner(_banner, "🚀💥 %s 召喚火箭砲！3 枚火箭！" % name, Color(1.0, 0.5, 0.2), 2.5)
		"rocket_launch":
			var no = data.get("rocket_no", 1)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "🚀 火箭砲"
				if is_instance_valid(value): value.text = "第 %d/3 枚發射！" % no
			BaseLucky.spawn_float_text(self, Vector2(640, 300), "🚀 第 %d 枚！" % no, Color(1.0, 0.5, 0.2), 22)
		"rocket_explode":
			var no = data.get("rocket_no", 1)
			var hits = data.get("hit_targets", [])
			BaseLucky.spawn_float_text(self, Vector2(640, 280), "💥 爆炸！命中 %d！" % hits.size(), Color(1.0, 0.4, 0.1), 24)
		"settle":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			var reward = data.get("total_reward", 0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "🚀 火箭砲結算", "size": 18, "color": Color(1.0, 0.5, 0.2)},
				{"text": "3 枚火箭全部爆炸！", "size": 18, "color": Color.WHITE},
				{"text": "獎勵：+%d" % reward, "size": 22, "color": Color(1.0, 0.85, 0.0)},
			], 3.0)
