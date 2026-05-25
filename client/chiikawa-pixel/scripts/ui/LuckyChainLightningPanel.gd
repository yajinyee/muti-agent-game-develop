## LuckyChainLightningPanel.gd — T106 幸運連鎖閃電魚面板
## lucky-panel-agent 負責維護
## 閃電主題：青藍色 + 電弧特效
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null
var _timer_bar: Dictionary = {}

func _ready() -> void:
	layer = 20
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			BaseLucky.fullscreen_flash(self, Color(0.0, 0.9, 1.0), 2)
			BaseLucky.show_banner(_banner, "⚡ %s 觸發連鎖閃電！" % name, Color(0.0, 0.9, 1.0), 2.5)
		"chain_hit":
			var chain = data.get("chain_count", 0)
			var mult = data.get("multiplier", 1.0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "⚡ 連鎖閃電"
				if is_instance_valid(value): value.text = "連鎖 %d | ×%.1f" % [chain, mult]
			BaseLucky.spawn_float_text(self, Vector2(640, 300), "⚡ ×%.1f" % mult, Color(0.0, 0.9, 1.0), 22)
		"settle":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			var reward = data.get("total_reward", 0)
			var mult = data.get("multiplier", 1.0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "⚡ 連鎖閃電結算", "size": 18, "color": Color(0.0, 0.9, 1.0)},
				{"text": "倍率：×%.1f" % mult, "size": 22, "color": Color(1.0, 0.85, 0.0)},
				{"text": "獎勵：+%d" % reward, "size": 20, "color": Color.WHITE},
			], 3.0)
