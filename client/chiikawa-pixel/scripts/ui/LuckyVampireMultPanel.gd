## LuckyVampireMultPanel.gd — T120 幸運吸血鬼魚面板
## lucky-panel-agent 負責維護
## 吸血鬼主題：深紫色 + 吸收計數器
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null

func _ready() -> void:
	layer = 34
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			BaseLucky.fullscreen_flash(self, Color(0.7, 0.2, 0.8), 2)
			BaseLucky.show_banner(_banner, "🧛 %s 觸發吸血鬼！每次擊破吸收倍率！" % name, Color(0.7, 0.2, 0.8), 2.5)
		"absorb":
			var cnt = data.get("absorb_count", 0)
			var mult = data.get("current_mult", 1.0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "🧛 吸血鬼"
				if is_instance_valid(value): value.text = "吸收 %d 次 | ×%.1f" % [cnt, mult]
			BaseLucky.spawn_float_text(self, Vector2(640, 300), "🧛 ×%.1f" % mult, Color(0.7, 0.2, 0.8), 20)
		"mult_mode":
			var mult = data.get("current_mult", 5.0)
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			BaseLucky.fullscreen_flash(self, Color(0.8, 0.0, 0.8), 3)
			BaseLucky.show_banner(_banner, "🧛✨ %s 進入倍率模式！×%.1f！10 秒！" % [name, mult], Color(0.8, 0.0, 0.8), 3.5)
		"mult_end":
			BaseLucky.show_banner(_banner, "🧛 吸血鬼倍率模式結束", Color(0.5, 0.5, 0.5), 1.5)
		"settle":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			var cnt = data.get("absorb_count", 0)
			var mult = data.get("current_mult", 1.0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "🧛 吸血鬼結算", "size": 18, "color": Color(0.7, 0.2, 0.8)},
				{"text": "吸收：%d 次 | ×%.1f" % [cnt, mult], "size": 20, "color": Color.WHITE},
			], 3.0)
