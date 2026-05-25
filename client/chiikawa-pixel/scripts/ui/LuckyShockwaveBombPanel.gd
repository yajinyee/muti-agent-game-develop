## LuckyShockwaveBombPanel.gd — T112 幸運全場震盪魚面板
## lucky-panel-agent 負責維護
## 震盪主題：橙紅色 + 全場衝擊波
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null

func _ready() -> void:
	layer = 26
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"shockwave_start":
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.42, 0.21), 3)
			BaseLucky.show_banner(_banner, "💥 %s 觸發全場震盪！全場 HP -35%！" % name, Color(1.0, 0.42, 0.21), 2.5)
		"shockwave_hit":
			var hits = data.get("hit_count", 0)
			var reward = data.get("total_reward", 0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "💥 全場震盪"
				if is_instance_valid(value): value.text = "命中 %d 個" % hits
			BaseLucky.spawn_float_text(self, Vector2(640, 300), "💥 命中 %d！" % hits, Color(1.0, 0.42, 0.21), 22)
		"super_shockwave":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.42, 0.21), 3)
			BaseLucky.show_banner(_banner, "💥🌊 超級震盪！%s 全服 ×1.8 加成 6 秒！" % name, Color(1.0, 0.42, 0.21), 3.5)
		"super_end":
			BaseLucky.show_banner(_banner, "💥 超級震盪加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"power_end":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
