## LuckyImmortalBossUltraPanel.gd — T247 幸運不死BOSS升級魚面板
## lucky-panel-agent 負責維護
## Immortal Boss Ultra 機制（Royal Fishing Immortal Boss 升級版）
## 不死 BOSS 連續獎勵：5 次復活（×100→×150→×200→×250→×300），全服 ×55.5 加成 111 秒
extends Node

const PANEL_COLOR = Color(0.55, 0.0, 0.0)  # 深紅色

var _canvas: CanvasLayer
var _banner: Dictionary
var _indicator: Dictionary
var _settle: Control
var _revive_no: int = 0

func _ready() -> void:
	_canvas = get_tree().get_root().get_node_or_null("Main/HUDLayer")
	if not is_instance_valid(_canvas):
		_canvas = CanvasLayer.new()
		_canvas.layer = 10
		get_tree().get_root().add_child(_canvas)
	_banner    = BaseLuckyPanel.create_banner(_canvas, 100.0, 71)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 160), 76)
	_settle    = BaseLuckyPanel.create_settle_popup(_canvas, 81)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"immortal_boss_ultra_start":
			_revive_no = 0
			var max_mult = data.get("max_mult", 300.0)
			BaseLuckyPanel.show_banner(_banner.panel,
				"💀👑 IMMORTAL BOSS ULTRA! %s summoned! 5 revivals! Up to x%.0f!" % [trigger_name, max_mult],
				PANEL_COLOR, 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 6)
			_show_indicator("IMMORTAL BOSS", "Revive 0/5 💀")
		"boss_revive":
			_revive_no = data.get("revive_no", _revive_no + 1)
			var revive_mult = data.get("revive_mult", 100.0)
			var hp_pct = data.get("hp_percent", 50.0)
			_show_indicator("IMMORTAL BOSS", "Revive %d/5 x%.0f 💀" % [_revive_no, revive_mult])
			BaseLuckyPanel.show_banner(_banner.panel,
				"💀 BOSS REVIVED! Revive %d/5: x%.0f! HP: %.0f%%" % [_revive_no, revive_mult, hp_pct],
				Color(0.8, 0.0, 0.0), 2.0)
		"immortal_boss_ultra_complete":
			var total = data.get("total_mult", 0.0)
			var bonus = data.get("global_bonus", 55.5)
			var secs = data.get("global_seconds", 111)
			BaseLuckyPanel.show_banner(_banner.panel,
				"💀👑 IMMORTAL BOSS ULTRA COMPLETE! %s: 5 revivals! Total x%.1f! GLOBAL x%.1f for %ds!" % [trigger_name, total, bonus, secs],
				Color(1.0, 0.2, 0.2), 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, Color(1.0, 0.2, 0.2), 9)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "💀👑 IMMORTAL BOSS ULTRA!", "size": 22, "color": PANEL_COLOR},
				{"text": "5 Revivals: x100→x150→x200→x250→x300", "size": 14, "color": Color.WHITE},
				{"text": "Total x%.1f" % total, "size": 18, "color": Color(1.0, 0.3, 0.3)},
				{"text": "🌍 GLOBAL x%.1f for %ds" % [bonus, secs], "size": 18, "color": Color(1.0, 0.2, 0.2)},
			], 5.0)
			_hide_indicator()

func _show_indicator(title: String, value: String) -> void:
	if _indicator.is_empty(): return
	var t = _indicator.get("title")
	var v = _indicator.get("value")
	var p = _indicator.get("panel")
	if is_instance_valid(t): t.text = title
	if is_instance_valid(v): v.text = value
	if is_instance_valid(p): p.visible = true

func _hide_indicator() -> void:
	var p = _indicator.get("panel")
	if is_instance_valid(p): p.visible = false
