## LuckyChainLongKingPanel.gd — T116 幸運千龍王輪盤魚面板
## lucky-panel-agent 負責維護
## 千龍王主題：金色 + 雙環輪盤，最高 1000x
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null

func _ready() -> void:
	layer = 30
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.85, 0.0), 2)
			BaseLucky.show_banner(_banner, "🐉👑 %s 觸發千龍王輪盤！最高 1000x！" % name, Color(1.0, 0.85, 0.0), 2.5)
		"spin":
			var inner = data.get("inner_mult", 1.0)
			var outer = data.get("outer_mult", 1.0)
			var final_m = data.get("final_mult", 1.0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "🐉👑 千龍王輪盤"
				if is_instance_valid(value): value.text = "×%.0f × ×%.0f = ×%.0f" % [inner, outer, final_m]
			BaseLucky.spawn_float_text(self, Vector2(640, 280), "🐉 ×%.0f！" % final_m, Color(1.0, 0.85, 0.0), 30)
		"result":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			var reward = data.get("reward", 0)
			var final_m = data.get("final_mult", 1.0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "🐉👑 千龍王輪盤", "size": 18, "color": Color(1.0, 0.85, 0.0)},
				{"text": "最終倍率：×%.0f" % final_m, "size": 26, "color": Color(1.0, 0.85, 0.0)},
				{"text": "獎勵：+%d" % reward, "size": 20, "color": Color.WHITE},
			], 4.0)
		"mega_win":
			var final_m = data.get("final_mult", 1.0)
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.85, 0.0), 5)
			BaseLucky.show_banner(_banner, "🐉👑✨ MEGA WIN！×%.0f！千龍王降臨！" % final_m, Color(1.0, 0.85, 0.0), 4.0)
