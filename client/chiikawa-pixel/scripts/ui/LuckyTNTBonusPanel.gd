extends BaseLuckyPanel
# T217 幸運 TNT 爆炸魚 — TNT Bonus 機制
# BGaming Fishing Club 2（2026-04）：水下大爆炸，全場 HP -80%，每個 ×100.0
# 全服 ×39.0 加成 78 秒

const PANEL_COLOR = Color(1.0, 0.271, 0.0, 0.92)   # 火橙紅
const ACCENT_COLOR = Color(1.0, 0.6, 0.0, 1.0)      # 橙黃

var _countdown: int = 3
var _blasted_count: int = 0
var _tnt_mult: float = 100.0
var _countdown_label: Label
var _blast_label: Label
var _tnt_icon_label: Label

func _ready() -> void:
	super._ready()
	_setup_tnt_ui()

func _setup_tnt_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "💣 TNT 引爆！"
	subtitle_text = "水下大爆炸 ×100.0"

	_tnt_icon_label = Label.new()
	_tnt_icon_label.text = "💣"
	_tnt_icon_label.add_theme_font_size_override("font_size", 48)
	_tnt_icon_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_tnt_icon_label)

	_countdown_label = Label.new()
	_countdown_label.text = "引爆倒數：3"
	_countdown_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_countdown_label.add_theme_font_size_override("font_size", 22)
	_countdown_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_countdown_label)

	_blast_label = Label.new()
	_blast_label.text = "炸毀目標：0 個"
	_blast_label.add_theme_color_override("font_color", Color.WHITE)
	_blast_label.add_theme_font_size_override("font_size", 16)
	_blast_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_blast_label)

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"tnt_countdown":
			_tnt_mult = event_data.get("tnt_mult", 100.0)
			_show_panel()
			_countdown_label.text = "引爆倒數：3"
			_play_countdown_animation()
		"tnt_tick":
			_countdown = event_data.get("countdown", 3)
			_countdown_label.text = "引爆倒數：%d" % _countdown
			_play_tick_animation()
		"tnt_explode":
			_blasted_count = event_data.get("blasted_count", 0)
			_blast_label.text = "炸毀目標：%d 個" % _blasted_count
			_countdown_label.text = "💥 爆炸！"
			_play_explode_animation()
		"tnt_perfect":
			var boost_mult = event_data.get("global_boost_mult", 39.0)
			var boost_secs = event_data.get("global_boost_secs", 78)
			_show_perfect_result(boost_mult, boost_secs)
		"tnt_end":
			_hide_panel_delayed(2.0)

func _play_countdown_animation() -> void:
	var tween = create_tween().set_loops(3)
	tween.tween_property(_tnt_icon_label, "scale", Vector2(1.2, 1.2), 0.3)
	tween.tween_property(_tnt_icon_label, "scale", Vector2(1.0, 1.0), 0.3)

func _play_tick_animation() -> void:
	var tween = create_tween()
	tween.tween_property(_countdown_label, "modulate", Color(1.0, 0.3, 0.0, 1.0), 0.1)
	tween.tween_property(_countdown_label, "modulate", Color.WHITE, 0.1)

func _play_explode_animation() -> void:
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(1.3, 1.3), 0.1)
	tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.2)

func _show_perfect_result(boost_mult: float, boost_secs: int) -> void:
	title_text = "💣✨ TNT 完美爆炸！"
	subtitle_text = "全服 ×%.1f 加成 %d 秒！" % [boost_mult, boost_secs]
	_update_title()
	_play_perfect_animation()
	_hide_panel_delayed(5.0)
