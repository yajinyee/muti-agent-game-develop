## LuckyDragonShotgunPanel.gd — T117 幸運龍力散彈魚面板
## lucky-panel-agent 負責維護
## 龍力散彈主題：紫色 + 8方向攻擊
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null

func _ready() -> void:
	layer = 31
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			BaseLucky.fullscreen_flash(self, Color(0.9, 0.4, 1.0), 2)
			BaseLucky.show_banner(_banner, "🐲💥 %s 觸發龍力散彈！8 方向攻擊！" % name, Color(0.9, 0.4, 1.0), 2.5)
		"shotgun_fire":
			var dir = data.get("direction", 0)
			var hits = data.get("total_hits", 0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "🐲 龍力散彈"
				if is_instance_valid(value): value.text = "方向 %d/8 | 命中 %d" % [dir + 1, hits]
			BaseLucky.spawn_float_text(self, Vector2(640, 300), "🐲 方向 %d！" % (dir + 1), Color(0.9, 0.4, 1.0), 18)
		"settle":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			var reward = data.get("total_reward", 0)
			var hits = data.get("total_hits", 0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "🐲 龍力散彈結算", "size": 18, "color": Color(0.9, 0.4, 1.0)},
				{"text": "命中：%d 個目標" % hits, "size": 20, "color": Color.WHITE},
				{"text": "獎勵：+%d" % reward, "size": 22, "color": Color(1.0, 0.85, 0.0)},
			], 3.0)
