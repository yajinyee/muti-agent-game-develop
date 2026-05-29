extends BaseLuckyPanel
# T220 幸運快速暴富魚 — Rapid Riches 機制
# Reflex Gaming「Big Game Fishing Rapid Riches」（2026-05）：5 秒內快速連擊，每次 ×200.0
# 全服 ×41.0 加成 82 秒（新史上最高）

const PANEL_COLOR = Color(1.0, 0.843, 0.0, 0.95)    # 金色（最高階）
const ACCENT_COLOR = Color(1.0, 1.0, 0.0, 1.0)       # 亮黃

var _hit_count: int = 0
var _time_left: float = 5.0
var _is_active: bool = false
var _hit_label: Label
var _timer_label: Label
var _combo_label: Label
var _timer: float = 0.0

func _ready() -> void:
	super._ready()
	_setup_rapid_ui()

func _setup_rapid_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "💰 快速暴富！"
	subtitle_text = "5 秒內快速連擊 ×200.0"

	_timer_label = Label.new()
	_timer_label.text = "剩餘時間：5.0 秒"
	_timer_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_timer_label.add_theme_font_size_override("font_size", 22)
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_timer_label)

	_hit_label = Label.new()
	_hit_label.text = "連擊次數：0 / 10"
	_hit_label.add_theme_color_override("font_color", Color.WHITE)
	_hit_label.add_theme_font_size_override("font_size", 18)
	_hit_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_hit_label)

	_combo_label = Label.new()
	_combo_label.text = "每次獎勵：×200.0"
	_combo_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.0, 1.0))
	_combo_label.add_theme_font_size_override("font_size", 16)
	_combo_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_combo_label)

func _process(delta: float) -> void:
	if not _is_active:
		return
	_timer += delta
	_time_left = max(0.0, 5.0 - _timer)
	_timer_label.text = "剩餘時間：%.1f 秒" % _time_left
	if _time_left <= 0.0:
		_is_active = false

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"rapid_start":
			_hit_count = 0
			_timer = 0.0
			_time_left = 5.0
			_is_active = true
			_show_panel()
			_hit_label.text = "連擊次數：0 / 10"
			_timer_label.text = "剩餘時間：5.0 秒"
			_play_rapid_start_animation()
		"rapid_hit":
			_hit_count = event_data.get("hit_count", 0)
			var reward = event_data.get("reward", 0)
			_hit_label.text = "連擊次數：%d / 10" % _hit_count
			_combo_label.text = "本次獎勵：%d（×200.0）" % reward
			_play_hit_animation()
		"rapid_perfect":
			_is_active = false
			var boost_mult = event_data.get("global_boost_mult", 41.0)
			var boost_secs = event_data.get("global_boost_secs", 82)
			_show_perfect_result(boost_mult, boost_secs)
		"rapid_end":
			_is_active = false
			_hide_panel_delayed(2.0)

func _play_rapid_start_animation() -> void:
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(1.15, 1.15), 0.1)
	tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.1)

func _play_hit_animation() -> void:
	var tween = create_tween()
	tween.tween_property(_hit_label, "modulate", ACCENT_COLOR, 0.05)
	tween.tween_property(_hit_label, "modulate", Color.WHITE, 0.1)

func _show_perfect_result(boost_mult: float, boost_secs: int) -> void:
	title_text = "💰🌟 完美快速暴富！"
	subtitle_text = "全服 ×%.1f 加成 %d 秒！（新史上最高）" % [boost_mult, boost_secs]
	_update_title()
	_play_perfect_animation()
	_hide_panel_delayed(6.0)
