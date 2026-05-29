extends BaseLuckyPanel
# T229 Lucky Magnetic Attraction Fish — Magnetic Attraction Mechanic
# Magnetic force pulls all targets to center, each x70.0
# Perfect (>=10 targets) -> Global x45.5 for 91s (surpasses T228 x45.0)

const PANEL_COLOR = Color(0.4, 0.2, 0.0, 0.95)    # Deep orange-brown
const ACCENT_COLOR = Color(1.0, 0.4, 0.0, 1.0)    # Orange (magnetic)

var _hit_label: Label
var _mult_label: Label
var _status_label: Label

func _ready() -> void:
	super._ready()
	_setup_magnetic_ui()

func _setup_magnetic_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "MAGNETIC ATTRACTION!"
	subtitle_text = "All targets pulled to center! x70.0 each!"

	_hit_label = Label.new()
	_hit_label.text = "Targets: 0"
	_hit_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_hit_label.add_theme_font_size_override("font_size", 22)
	_hit_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_hit_label)

	_mult_label = Label.new()
	_mult_label.text = "Total: x0.0"
	_mult_label.add_theme_color_override("font_color", Color.WHITE)
	_mult_label.add_theme_font_size_override("font_size", 24)
	_mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_mult_label)

	_status_label = Label.new()
	_status_label.text = "Magnetic force activating..."
	_status_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_status_label.add_theme_font_size_override("font_size", 13)
	_status_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_status_label)

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"magnetic_start":
			_show_panel()
			_hit_label.text = "Targets: 0"
			_mult_label.text = "Total: x0.0"
			_status_label.text = "Magnetic force activating..."
			_play_start_animation()
		"magnetic_pull":
			var count = event_data.get("target_count", 0)
			var total = event_data.get("total_mult", 0.0)
			_hit_label.text = "Targets: %d" % count
			_mult_label.text = "Total: x%.1f" % total
			_status_label.text = "Pulling all targets to center!"
		"magnetic_result":
			var count = event_data.get("hit_count", 0)
			var total = event_data.get("total_mult", 0.0)
			var perfect = event_data.get("is_perfect", false)
			_hit_label.text = "Targets: %d" % count
			_mult_label.text = "Total: x%.1f" % total
			if perfect:
				_status_label.text = "PERFECT! Global x45.5 for 91s!"
				_play_perfect_animation()
			else:
				_status_label.text = "Magnetic complete!"
			_schedule_hide(4.0)
		"magnetic_bonus":
			var bonus = event_data.get("bonus_mult", 0.0)
			_status_label.text = "Bonus wave! x%.1f!" % bonus
