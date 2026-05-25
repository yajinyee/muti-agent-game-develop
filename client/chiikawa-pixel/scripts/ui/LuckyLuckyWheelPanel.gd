## LuckyLuckyWheelPanel.gd — T125 幸運大轉盤魚面板
## lucky-panel-agent 負責維護
## 大轉盤主題：彩虹色 + 8格轉盤結果
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null

func _ready() -> void:
	layer = 39
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.5, 0.0), 2)
			BaseLucky.show_banner(_banner, "🎡 %s 觸發幸運大轉盤！8 格隨機獎勵！" % name, Color(1.0, 0.5, 0.0), 2.5)
		"spin":
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "🎡 大轉盤"
				if is_instance_valid(value): value.text = "轉動中..."
		"result":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			var slot = data.get("slot", "")
			var reward = data.get("reward", 0)
			var mult = data.get("multiplier", 1.0)
			var is_jackpot = data.get("is_jackpot", false)
			if is_jackpot:
				BaseLucky.fullscreen_flash(self, Color(1.0, 0.85, 0.0), 5)
				BaseLucky.show_banner(_banner, "🎡🌟 大獎！%s 獲得 %s！" % [name, slot], Color(1.0, 0.85, 0.0), 4.0)
			else:
				BaseLucky.spawn_float_text(self, Vector2(640, 280), "🎡 %s！" % slot, Color(1.0, 0.5, 0.0), 26)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "🎡 大轉盤結算", "size": 18, "color": Color(1.0, 0.5, 0.0)},
				{"text": "結果：%s" % slot, "size": 20, "color": Color.WHITE},
				{"text": "獎勵：+%d" % reward, "size": 22, "color": Color(1.0, 0.85, 0.0)},
			], 3.5)
		"all_hp_half":
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.2, 0.2), 3)
			BaseLucky.show_banner(_banner, "🎡💥 全場 HP -50%！", Color(1.0, 0.2, 0.2), 2.0)
