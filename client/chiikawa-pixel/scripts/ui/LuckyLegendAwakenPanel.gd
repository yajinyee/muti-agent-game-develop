extends BaseLuckyPanel
# T226 Lucky Legend Awaken Fish — Legend Dragon Awaken Upgrade
# Royal Fishing Jili (2026): Humpback Whale 90-150x / Legend Dragon 120-200x
# 8 consecutive rewards, each increasing
# Global x44.0 for 88s (surpasses T225 x43.5)

const PANEL_COLOR = Color(0.6, 0.1, 0.0, 0.95)    # Deep crimson (Legend theme)
const ACCENT_COLOR = Color(1.0, 0.4, 0.0, 1.0)    # Fire orange

var _round_label: Label
var _total_label: Label
var _mode_label: Label
var _reward_label: Label

func _ready() -> void:
	super._ready()
	_setup_awaken_ui()

func _setup_awaken_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "Legend Awaken!"
	subtitle_text = "8 consecutive rewards!"

	_mode_label = Label.new()
	_mode_label.text = "Mode: ?"
	_mode_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_mode_label.add_theme_font_size_override("font_size", 20)
	_mode_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_mode_label)

	_round_label = Label.new()
	_round_label.text = "Round: 0 / 8"
	_round_label.add_theme_color_override("font_color", Color.WHITE)
	_round_label.add_theme_font_size_override("font_size", 18)
	_round_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_round_label)

	_total_label = Label.new()
	_total_label.text = "Total: x0.0"
	_total_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0, 1.0))
	_total_label.add_theme_font_size_override("font_size", 22)
	_total_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_total_label)

	_reward_label = Label.new()
	_reward_label.text = "Humpback: 90-150x | Legend: 120-200x"
	_reward_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_reward_label.add_theme_font_size_override("font_size", 11)
	_reward_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_reward_label)

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"awaken_start":
			_show_panel()
			var mode = event_data.get("mode_name", "Legend Dragon")
			var base = event_data.get("base_mult", 20.0)
			_mode_label.text = "Mode: %s (x%.0f base)" % [mode, base]
			_round_label.text = "Round: 0 / 8"
			_total_label.text = "Total: x0.0"
			_reward_label.text = "Awakening..."
			_play_start_animation()
		"awaken_reward":
			var round = event_data.get("round", 1)
			var round_mult = event_data.get("round_mult", 0.0)
			var total = event_data.get("total_mult", 0.0)
			var mode = event_data.get("mode_name", "Legend Dragon")
			_round_label.text = "Round: %d / 8" % round
			_total_label.text = "Total: x%.0f" % total
			_reward_label.text = "%s Round %d: +x%.1f" % [mode, round, round_mult]
			_play_reward_animation()
		"awaken_settle":
			var total = event_data.get("total_mult", 0.0)
			var mode = event_data.get("mode_name", "Legend Dragon")
			title_text = "Legend Awaken Done!"
			subtitle_text = "%s x%.0f total!" % [mode, total]
			_update_title()
			_hide_panel_delayed(3.0)
		"global_boost":
			var g_mult = event_data.get("global_mult", 44.0)
			var dur = event_data.get("duration", 88)
			var mode = event_data.get("mode_name", "Legend Dragon")
			title_text = "Global x%.1f for %ds!" % [g_mult, dur]
			subtitle_text = "%s Awaken - All-time high!" % mode
			_update_title()
			_play_perfect_animation()
			_hide_panel_delayed(6.0)

func _play_start_animation() -> void:
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(1.2, 1.2), 0.15)
	tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.15)

func _play_reward_animation() -> void:
	var tween = create_tween()
	tween.tween_property(_total_label, "modulate", ACCENT_COLOR, 0.06)
	tween.tween_property(_total_label, "modulate", Color.WHITE, 0.2)
