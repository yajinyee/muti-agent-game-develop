## LuckyChainMeteorPanel.gd — T129 幸運連鎖隕石魚面板
## lucky-panel-agent 負責維護
extends CanvasLayer

const BaseLuckyPanelScript = preload("res://scripts/ui/BaseLuckyPanel.gd")
var _base: Node

func _ready() -> void:
	layer = 29
	_base = BaseLuckyPanelScript.new()
	_base.name = "BaseLuckyPanel"
	add_child(_base)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"meteor_start":
			_base.show_banner("☄️ 連鎖隕石雨！", Color(0.8, 0.3, 0.1))
			_base.show_indicator("隕石 0/5", Color(1.0, 0.5, 0.2))
			_base.flash_screen(Color(0.8, 0.3, 0.1), 3)
		"meteor_hit":
			var idx = data.get("meteor_index", 1)
			var radius = data.get("aoe_radius", 150)
			_base.update_indicator("隕石 %d/5 r=%.0f" % [idx, radius], Color(1.0, 0.6, 0.2))
			_base.flash_screen(Color(1.0, 0.4, 0.1), 1)
		"meteor_miss":
			var idx = data.get("meteor_index", 1)
			_base.update_indicator("隕石 %d/5 空揮！" % idx, Color(0.6, 0.6, 0.6))
		"meteor_perfect":
			_base.show_banner("完美隕石雨！全服×2.5！", Color(1.0, 0.85, 0.0))
			_base.flash_screen(Color(1.0, 0.85, 0.0), 3)
			_base.show_indicator("×2.5 加成中", Color(1.0, 0.85, 0.0))
		"meteor_perfect_end":
			_base.hide_indicator()
			_base.hide_banner()
