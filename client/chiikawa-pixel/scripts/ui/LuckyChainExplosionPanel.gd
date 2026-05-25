## LuckyChainExplosionPanel.gd — T115 幸運連鎖爆炸魚面板
## lucky-panel-agent 負責維護
## 連鎖爆炸主題：火紅色 + 連鎖計數器
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null

func _ready() -> void:
	layer = 29
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"chain_start":
			BaseLucky.fullscreen_flash(self, Color(0.9, 0.2, 0.15), 2)
			BaseLucky.show_banner(_banner, "💥 %s 觸發連鎖爆炸！12 秒連鎖模式！" % name, Color(0.9, 0.2, 0.15), 2.5)
		"chain_explode":
			var cnt = data.get("chain_count", 0)
			var mult = data.get("accum_mult", 1.0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "💥 連鎖爆炸"
				if is_instance_valid(value): value.text = "×%.1f | %d 次" % [mult, cnt]
			BaseLucky.spawn_float_text(self, Vector2(640, 300), "💥 ×%.1f" % mult, Color(0.9, 0.2, 0.15), 20)
		"chain_burst":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.85, 0.0), 3)
			BaseLucky.show_banner(_banner, "💥🔥 連鎖爆發！%s 全服 ×2.5 加成 6 秒！" % name, Color(1.0, 0.85, 0.0), 3.5)
		"burst_end":
			BaseLucky.show_banner(_banner, "💥 連鎖爆發加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"chain_end":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			var cnt = data.get("chain_count", 0)
			var reward = data.get("total_reward", 0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "💥 連鎖爆炸結算", "size": 18, "color": Color(0.9, 0.2, 0.15)},
				{"text": "連鎖：%d 次" % cnt, "size": 20, "color": Color.WHITE},
				{"text": "獎勵：+%d" % reward, "size": 22, "color": Color(1.0, 0.85, 0.0)},
			], 3.0)
