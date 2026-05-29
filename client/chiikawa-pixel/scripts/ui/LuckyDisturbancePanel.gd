extends BaseLuckyPanel
# T218 幸運擾動魚 — Disturbance System
# Fisch Roblox（2026-01）：活躍度越高倍率越高，最高 ×50.0
# 全服 ×39.5 加成 79 秒

const PANEL_COLOR = Color(0.0, 0.808, 0.820, 0.92)  # 深青色
const ACCENT_COLOR = Color(0.0, 1.0, 0.9, 1.0)       # 亮青

var _disturbance: int = 0
var _disturb_mult: float = 5.0
var _disturbed_count: int = 0
var _disturbance_label: Label
var _mult_label: Label
var _disturb_bar: ProgressBar

func _ready() -> void:
	super._ready()
	_setup_disturbance_ui()

func _setup_disturbance_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "🌊 擾動爆發！"
	subtitle_text = "活躍度越高倍率越高"

	_disturbance_label = Label.new()
	_disturbance_label.text = "擾動值：0 / 30"
	_disturbance_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_disturbance_label.add_theme_font_size_override("font_size", 18)
	_disturbance_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_disturbance_label)

	_disturb_bar = ProgressBar.new()
	_disturb_bar.min_value = 0
	_disturb_bar.max_value = 30
	_disturb_bar.value = 0
	_disturb_bar.custom_minimum_size = Vector2(200, 20)
	content_container.add_child(_disturb_bar)

	_mult_label = Label.new()
	_mult_label.text = "擾動倍率：×5.0"
	_mult_label.add_theme_color_override("font_color", Color.WHITE)
	_mult_label.add_theme_font_size_override("font_size", 16)
	_mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_mult_label)

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"disturbance_start":
			_disturbance = event_data.get("disturbance", 1)
			_disturb_mult = event_data.get("disturb_mult", 5.0)
			_show_panel()
			_disturbance_label.text = "擾動值：%d / 30" % _disturbance
			_disturb_bar.value = _disturbance
			_mult_label.text = "擾動倍率：×%.1f" % _disturb_mult
			_play_disturbance_animation()
		"disturbance_hit":
			_disturbed_count += 1
			_mult_label.text = "影響目標：%d 個" % _disturbed_count
		"disturbance_perfect":
			var boost_mult = event_data.get("global_boost_mult", 39.5)
			var boost_secs = event_data.get("global_boost_secs", 79)
			_show_perfect_result(boost_mult, boost_secs)
		"disturbance_end":
			_hide_panel_delayed(2.0)

func _play_disturbance_animation() -> void:
	var tween = create_tween().set_loops(3)
	tween.tween_property(_disturbance_label, "modulate", ACCENT_COLOR, 0.2)
	tween.tween_property(_disturbance_label, "modulate", Color.WHITE, 0.2)

func _show_perfect_result(boost_mult: float, boost_secs: int) -> void:
	title_text = "🌊✨ 完美擾動！"
	subtitle_text = "全服 ×%.1f 加成 %d 秒！" % [boost_mult, boost_secs]
	_update_title()
	_play_perfect_animation()
	_hide_panel_delayed(5.0)
