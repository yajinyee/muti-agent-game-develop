extends BaseLuckyPanel
# T233 Lucky Cosmic Restart Fish — Cosmic Restart Mechanic
# Full field clear (each target x100.0), Global x47.5 for 95s (NEW ALL-TIME HIGH)

const PANEL_COLOR = Color(0.2, 0.0, 0.3, 0.95)    # Deep cosmic purple
const ACCENT_COLOR = Color(1.0, 0.0, 1.0, 1.0)    # Magenta (cosmic)

var _target_label: Label
var _mult_label: Label
var _status_label: Label

func _ready() -> void:
	super._ready()
	_setup_restart_ui()

func _setup_restart_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "COSMIC RESTART!"
	subtitle_text = "Universe restarted! Every target x100.0!"

	_target_label = Label.new()
	_target_label.text = "Targets: 0"
	_target_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_target_label.add_theme_font_size_override("font_size", 22)
	_target_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_target_label)

	_mult_label = Label.new()
	_mult_label.text = "Total: x0.0"
	_mult_label.add_theme_color_override("font_color", Color.WHITE)
	_mult_label.add_theme_font_size_override("font_size", 26)
	_mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_mult_label)

	_status_label = Label.new()
	_status_label.text = "Cosmic restart initiating..."
	_status_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_status_label.add_theme_font_size_override("font_size", 13)
	_status_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_status_label)

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"restart_start":
			_show_panel()
			_target_label.text = "Targets: 0"
			_mult_label.text = "Total: x0.0"
			_status_label.text = "COSMIC RESTART INITIATED!"
			_play_start_animation()
		"restart_result":
			var count = event_data.get("target_count", 0)
			var total = event_data.get("total_mult", 0.0)
			var global_bonus = event_data.get("global_bonus", 47.5)
			var global_secs = event_data.get("global_seconds", 95)
			_target_label.text = "Targets: %d" % count
			_mult_label.text = "Total: x%.1f" % total
			_status_label.text = "UNIVERSE RESTARTED! Global x%.1f for %ds! NEW ALL-TIME HIGH!" % [global_bonus, global_secs]
			_play_perfect_animation()
			_schedule_hide(5.0)
