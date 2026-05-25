## LuckyCrabTorpedoPanel.gd — T107 幸運螃蟹魚雷面板
## lucky-panel-agent 負責維護
## 螃蟹主題：橙紅色 + 爆炸特效
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null

func _ready() -> void:
	layer = 21
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.4, 0.1), 2)
			BaseLucky.show_banner(_banner, "🦀 %s 發射螃蟹魚雷！3 次 AOE 爆炸！" % name, Color(1.0, 0.5, 0.2), 2.5)
		"explosion":
			var no = data.get("explosion_no", 1)
			var hits = data.get("hit_count", 0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "🦀 螃蟹魚雷"
				if is_instance_valid(value): value.text = "第 %d/3 枚 | 命中 %d" % [no, hits]
			BaseLucky.spawn_float_text(self, Vector2(640, 300), "💥 爆炸 %d/3！" % no, Color(1.0, 0.5, 0.2), 24)
		"settle":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			var reward = data.get("total_reward", 0)
			var hits = data.get("total_hits", 0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "🦀 螃蟹魚雷結算", "size": 18, "color": Color(1.0, 0.5, 0.2)},
				{"text": "命中：%d 個目標" % hits, "size": 20, "color": Color.WHITE},
				{"text": "獎勵：+%d" % reward, "size": 22, "color": Color(1.0, 0.85, 0.0)},
			], 3.0)
