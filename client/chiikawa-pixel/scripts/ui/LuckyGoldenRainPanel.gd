## LuckyGoldenRainPanel.gd — T122 幸運黃金雨魚面板
## lucky-panel-agent 負責維護
## 黃金雨主題：金色 + 金幣收集計數器
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null

func _ready() -> void:
	layer = 36
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.85, 0.0), 2)
			BaseLucky.show_banner(_banner, "🌧️💰 %s 觸發黃金雨！收集金幣！" % name, Color(1.0, 0.85, 0.0), 2.5)
		"coin_spawn":
			var count = data.get("coin_count", 0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "💰 黃金雨"
				if is_instance_valid(value): value.text = "金幣 %d 個！" % count
		"coin_collect":
			var collected = data.get("collected", 0)
			var total = data.get("total_coins", 0)
			var ind_panel = _indicator.get("panel")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(value): value.text = "收集 %d/%d" % [collected, total]
			BaseLucky.spawn_float_text(self, Vector2(640, 300), "💰 +1", Color(1.0, 0.85, 0.0), 18)
		"golden_harvest":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.85, 0.0), 3)
			BaseLucky.show_banner(_banner, "💰✨ 黃金豐收！%s 全服 ×2.0 加成 6 秒！" % name, Color(1.0, 0.85, 0.0), 3.5)
		"harvest_end":
			BaseLucky.show_banner(_banner, "💰 黃金豐收加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"settle":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			var collected = data.get("collected", 0)
			var reward = data.get("total_reward", 0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "💰 黃金雨結算", "size": 18, "color": Color(1.0, 0.85, 0.0)},
				{"text": "收集：%d 個金幣" % collected, "size": 20, "color": Color.WHITE},
				{"text": "獎勵：+%d" % reward, "size": 22, "color": Color(1.0, 0.85, 0.0)},
			], 3.0)
