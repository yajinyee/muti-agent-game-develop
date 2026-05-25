## LuckyTimeFreezePanel.gd — T114 幸運時間凍結魚面板
## lucky-panel-agent 負責維護
## 冰凍主題：冰藍色 + 凍結計時器
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null
var _timer_bar: Dictionary = {}

func _ready() -> void:
	layer = 28
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)
	var ind_panel = _indicator.get("panel")
	if is_instance_valid(ind_panel):
		var bar_container = Control.new()
		bar_container.position = Vector2(8, 50)
		ind_panel.add_child(bar_container)
		_timer_bar = BaseLucky.create_timer_bar(bar_container, 184.0, Color(0.0, 0.9, 1.0))

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"freeze_start":
			BaseLucky.fullscreen_flash(self, Color(0.0, 0.9, 1.0), 2)
			BaseLucky.show_banner(_banner, "❄️ %s 觸發時間凍結！全場凍結 8 秒！傷害 ×1.8！" % name, Color(0.0, 0.9, 1.0), 2.5)
		"freeze_tick":
			var tl = data.get("time_left", 0.0)
			var kills = data.get("kill_count", 0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "❄️ 時間凍結"
				if is_instance_valid(value): value.text = "%.0fs | 擊破 %d 條" % [tl, kills]
			BaseLucky.update_timer_bar(_timer_bar, tl / 8.0)
		"freeze_end":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			BaseLucky.fullscreen_flash(self, Color(0.0, 0.9, 1.0), 3)
			BaseLucky.show_banner(_banner, "❄️💥 冰裂爆炸！全場 HP -25%！", Color(0.6, 0.9, 1.0), 2.0)
		"perfect_freeze":
			var kills = data.get("kill_count", 0)
			BaseLucky.fullscreen_flash(self, Color(0.0, 0.9, 1.0), 3)
			BaseLucky.show_banner(_banner, "❄️✨ 完美凍結！%s 擊破 %d 條！全服 ×2.0 加成 5 秒！" % [name, kills], Color(0.0, 0.9, 1.0), 3.5)
		"perfect_end":
			BaseLucky.show_banner(_banner, "❄️ 完美凍結加成結束", Color(0.7, 0.7, 0.7), 1.5)
