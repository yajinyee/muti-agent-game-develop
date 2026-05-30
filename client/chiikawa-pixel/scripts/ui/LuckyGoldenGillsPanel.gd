## LuckyGoldenGillsPanel.gd — T252 幸運黃金鰓魚面板
## lucky-panel-agent 負責維護
## Golden Gills Jackpot Respin 機制（Atomic Slot Lab 2026）
## 磁力連鎖 + 4 層 Jackpot（Mini/Minor/Major/Grand）
## Grand Jackpot → 全服 ×58.0 加成 116 秒
extends Node

const PANEL_COLOR = Color(1.0, 0.84, 0.0)  # 黃金色（黃金鰓）
const JACKPOT_COLORS = {
	"Mini":  Color(0.8, 0.8, 0.8),
	"Minor": Color(0.0, 0.8, 1.0),
	"Major": Color(0.8, 0.0, 1.0),
	"Grand": Color(1.0, 0.84, 0.0),
}

var _canvas: CanvasLayer
var _banner: Dictionary
var _indicator: Dictionary
var _settle: Control

func _ready() -> void:
	_canvas = get_tree().get_root().get_node_or_null("Main/HUDLayer")
	if not is_instance_valid(_canvas):
		_canvas = CanvasLayer.new()
		_canvas.layer = 10
		get_tree().get_root().add_child(_canvas)
	_banner    = BaseLuckyPanel.create_banner(_canvas, 100.0, 73)
	_indicator = BaseLuckyPanel.create_indicator(_canvas, Vector2(1060, 170), 78)
	_settle    = BaseLuckyPanel.create_settle_popup(_canvas, 83)

func handle_event(data: Dictionary) -> void:
	var event = data.get("event", "")
	var trigger_name = data.get("trigger_name", "???")
	match event:
		"golden_gills_start":
			BaseLuckyPanel.show_banner(_banner.panel,
				"🐟💰 GOLDEN GILLS! %s activated Jackpot Respin! 4-tier Jackpot: Mini/Minor/Major/Grand!" % trigger_name,
				PANEL_COLOR, 4.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 5)
			_show_indicator("GOLDEN GILLS", "Respin 0/6 💰")
		"jackpot_triggered":
			var tier = data.get("jackpot_tier", "Mini")
			var color = JACKPOT_COLORS.get(tier, PANEL_COLOR)
			BaseLuckyPanel.show_banner(_banner.panel,
				"💰 %s JACKPOT TRIGGERED! %s wins the %s Jackpot!" % [tier.to_upper(), trigger_name, tier],
				color, 3.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, color, 4)
		"respin":
			var respin_no = data.get("respin_no", 1)
			_show_indicator("GOLDEN GILLS", "Respin %d/6 💰" % respin_no)
		"golden_gills_complete":
			var jackpots = data.get("triggered_jackpots", [])
			var bonus = data.get("global_bonus", 58.0)
			var secs = data.get("global_seconds", 116)
			BaseLuckyPanel.show_banner(_banner.panel,
				"🐟💰 GOLDEN GILLS COMPLETE! %s: Jackpots=%s! GLOBAL x%.1f for %ds!" % [trigger_name, str(jackpots), bonus, secs],
				Color(1.0, 0.95, 0.5), 5.0)
			BaseLuckyPanel.fullscreen_flash(_canvas, PANEL_COLOR, 7)
			BaseLuckyPanel.show_settle_popup(_settle, [
				{"text": "🐟💰 GOLDEN GILLS!", "size": 22, "color": PANEL_COLOR},
				{"text": "Golden Gills Jackpot Respin", "size": 14, "color": Color(1.0, 0.95, 0.5)},
				{"text": "Jackpots: %s" % str(jackpots), "size": 14, "color": Color.WHITE},
				{"text": "🌐 GLOBAL x%.1f for %ds" % [bonus, secs], "size": 18, "color": PANEL_COLOR},
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
