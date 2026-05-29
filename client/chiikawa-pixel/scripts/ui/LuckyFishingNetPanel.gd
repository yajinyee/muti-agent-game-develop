extends BaseLuckyPanel
# T216 幸運漁網魚 — Fishing Net 機制
# BGaming Fishing Club 2（2026-04）：撒網捕獲全場所有目標，每個獎勵 ×60.0
# 全服 ×38.5 加成 77 秒

const PANEL_COLOR = Color(0.118, 0.565, 1.0, 0.92)  # 深海藍
const ACCENT_COLOR = Color(0.0, 0.8, 1.0, 1.0)       # 亮藍

var _caught_count: int = 0
var _net_mult: float = 60.0
var _net_label: Label
var _caught_label: Label
var _progress_bar: ProgressBar

func _ready() -> void:
	super._ready()
	_setup_fishing_net_ui()

func _setup_fishing_net_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "🎣 漁網撒出！"
	subtitle_text = "全場目標被捕獲 ×60.0"

	_net_label = Label.new()
	_net_label.text = "漁網倍率：×60.0"
	_net_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_net_label.add_theme_font_size_override("font_size", 18)
	_net_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_net_label)

	_caught_label = Label.new()
	_caught_label.text = "捕獲目標：0 個"
	_caught_label.add_theme_color_override("font_color", Color.WHITE)
	_caught_label.add_theme_font_size_override("font_size", 16)
	_caught_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_caught_label)

	_progress_bar = ProgressBar.new()
	_progress_bar.min_value = 0
	_progress_bar.max_value = 10
	_progress_bar.value = 0
	_progress_bar.custom_minimum_size = Vector2(200, 20)
	content_container.add_child(_progress_bar)

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"net_cast":
			_net_mult = event_data.get("net_mult", 60.0)
			_show_panel()
			_net_label.text = "漁網倍率：×%.1f" % _net_mult
			_play_cast_animation()
		"net_haul":
			_caught_count = event_data.get("caught_count", 0)
			_caught_label.text = "捕獲目標：%d 個" % _caught_count
			_progress_bar.value = min(_caught_count, 10)
		"net_perfect":
			var boost_mult = event_data.get("global_boost_mult", 38.5)
			var boost_secs = event_data.get("global_boost_secs", 77)
			_show_perfect_result(boost_mult, boost_secs)
		"net_end":
			_hide_panel_delayed(2.0)

func _play_cast_animation() -> void:
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(1.1, 1.1), 0.15)
	tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.15)

func _show_perfect_result(boost_mult: float, boost_secs: int) -> void:
	title_text = "🎣✨ 完美漁網！"
	subtitle_text = "全服 ×%.1f 加成 %d 秒！" % [boost_mult, boost_secs]
	_update_title()
	_play_perfect_animation()
	_hide_panel_delayed(5.0)
