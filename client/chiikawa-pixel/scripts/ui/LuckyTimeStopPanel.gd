extends BaseLuckyPanel
# T232 Lucky Time Stop Fish — Time Stop Mechanic
# Freeze all targets 15s (damage x5.0), freeze end HP -70%
# Kill >=15 during freeze -> Perfect, Global x47.0 for 94s

const PANEL_COLOR = Color(0.0, 0.2, 0.4, 0.95)    # Deep blue
const ACCENT_COLOR = Color(0.0, 0.8, 1.0, 1.0)    # Cyan-blue (ice)

var _timer_label: Label
var _kill_label: Label
var _status_label: Label

func _ready() -> void:
	super._ready()
	_setup_timestop_ui()

func _setup_timestop_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "TIME STOP!"
	subtitle_text = "All time frozen! 15s x5.0 damage!"

	_timer_label = Label.new()
	_timer_label.text = "Time: 15s"
	_timer_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_timer_label.add_theme_font_size_override("font_size", 26)
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_timer_label)

	_kill_label = Label.new()
	_kill_label.text = "Kills: 0"
	_kill_label.add_theme_color_override("font_color", Color.WHITE)
	_kill_label.add_theme_font_size_override("font_size", 22)
	_kill_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_kill_label)

	_status_label = Label.new()
	_status_label.text = "Time frozen! Kill 15+ for PERFECT!"
	_status_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_status_label.add_theme_font_size_override("font_size", 13)
	_status_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_status_label)

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"freeze_start":
			_show_panel()
			_timer_label.text = "Time: 15s"
			_kill_label.text = "Kills: 0"
			_status_label.text = "Time frozen! Kill 15+ for PERFECT!"
			_play_start_animation()
		"freeze_tick":
			var time_left = event_data.get("time_left", 0)
			var kills = event_data.get("kill_count", 0)
			_timer_label.text = "Time: %ds" % time_left
			_kill_label.text = "Kills: %d" % kills
			if kills >= 15:
				_status_label.text = "PERFECT ACHIEVED! Keep going!"
			else:
				_status_label.text = "Need %d more kills!" % (15 - kills)
		"freeze_end":
			var kills = event_data.get("kill_count", 0)
			var perfect = event_data.get("is_perfect", false)
			_timer_label.text = "FREEZE END!"
			_kill_label.text = "Kills: %d" % kills
			if perfect:
				_status_label.text = "PERFECT TIME STOP! Global x47.0 for 94s!"
				_play_perfect_animation()
			else:
				_status_label.text = "Time Stop complete! Killed %d" % kills
			_schedule_hide(4.0)
