## LuckyFreezeBombPanel.gd — T123 幸運冰凍炸彈魚面板
## lucky-panel-agent 負責維護
## 冰凍炸彈主題：冰藍色 + 凍結倒數
extends CanvasLayer

const BaseLucky = preload("res://scripts/ui/BaseLuckyPanel.gd")

var _banner: Control = null
var _indicator: Dictionary = {}
var _settle: Control = null
var _timer_bar: Dictionary = {}

func _ready() -> void:
	layer = 37
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
		"trigger":
			BaseLucky.fullscreen_flash(self, Color(0.0, 0.9, 1.0), 2)
			BaseLucky.show_banner(_banner, "🧊💣 %s 投擲冰凍炸彈！凍結 3 秒後爆炸！" % name, Color(0.0, 0.9, 1.0), 2.5)
		"freeze":
			var tl = data.get("time_left", 0.0)
			var ind_panel = _indicator.get("panel")
			var title = _indicator.get("title")
			var value = _indicator.get("value")
			if is_instance_valid(ind_panel):
				ind_panel.visible = true
				if is_instance_valid(title): title.text = "🧊 冰凍炸彈"
				if is_instance_valid(value): value.text = "凍結 %.0fs 後爆炸！" % tl
			BaseLucky.update_timer_bar(_timer_bar, tl / 3.0)
		"explode":
			var ind_panel = _indicator.get("panel")
			if is_instance_valid(ind_panel): ind_panel.visible = false
			BaseLucky.fullscreen_flash(self, Color(0.0, 0.9, 1.0), 3)
			BaseLucky.show_banner(_banner, "💣💥 冰凍爆炸！HP -60%！", Color(0.0, 0.9, 1.0), 2.0)
		"ice_burst":
			var hits = data.get("hit_count", 0)
			BaseLucky.fullscreen_flash(self, Color(1.0, 0.85, 0.0), 3)
			BaseLucky.show_banner(_banner, "🧊✨ 冰爆完美！%s 命中 %d 個！全服 ×2.2 加成 5 秒！" % [name, hits], Color(1.0, 0.85, 0.0), 3.5)
		"ice_burst_end":
			BaseLucky.show_banner(_banner, "🧊 冰爆完美加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"settle":
			var reward = data.get("total_reward", 0)
			BaseLucky.show_settle_popup(_settle, [
				{"text": "🧊 冰凍炸彈結算", "size": 18, "color": Color(0.0, 0.9, 1.0)},
				{"text": "獎勵：+%d" % reward, "size": 22, "color": Color(1.0, 0.85, 0.0)},
			], 3.0)
