## LuckyThunderLobsterPanel.gd — T110 幸運雷霆龍蝦面板
## lucky-panel-agent 負責維護
## 雷霆主題：橙色 + 自動射擊計時器
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null
var _timer_bar: Dictionary = {}

func _ready() -> void:
	layer = 24
	_banner = BaseLucky.create_banner(self, 120.0, 60).panel
	_indicator = BaseLucky.create_indicator(self, Vector2(1060, 100), 65)
	_settle = BaseLucky.create_settle_popup(self, 70)
	# 建立計時條
	var ind_panel = _indicator.get("panel")
	if is_instance_valid(ind_panel):
		var bar_container = Control.new()
		bar_container.position = Vector2(8, 50)
		ind_panel.add_child(bar_container)
		_timer_bar = BaseLucky.create_timer_bar(bar_container, 184.0, Color(1.0, 0.5, 0.2))

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.5, 0.2), 2)
			BaseLucky.show_banner(_banner, "🦞⚡ %s 觸發雷霆龍蝦！15 秒免費射擊！" % name, Color(1.0, 0.5, 0.2), 2.5)
		"auto_fire":
			var tl = data.get("time_left", 0.0)
			var kills = data.get("kill_count", 0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "🦞⚡ 雷霆模式"
				if is_instance_valid(value): value.text = "%.0fs | 擊破 %d 條" % [tl, kills]
			BaseLucky.update_timer_bar(_timer_bar, tl / 15.0)
		"end":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			var reward = data.get("total_reward", 0)
			var kills = data.get("kill_count", 0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "🦞 雷霆龍蝦結算", "size": 18, "color": Color(1.0, 0.5, 0.2)},
				{"text": "擊破：%d 條" % kills, "size": 20, "color": Color.WHITE},
				{"text": "獎勵：+%d" % reward, "size": 22, "color": Color(1.0, 0.85, 0.0)},
			], 3.0)
