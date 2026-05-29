extends BaseLuckyPanel
# T231 Lucky Holy Pillar Fish — Holy Pillar Mechanic
# 12 holy pillars descend (each HP -50%), hit >=8 -> Perfect Holy
# Perfect -> Global x46.5 for 93s (surpasses T230 x46.0)

const PANEL_COLOR = Color(0.3, 0.3, 0.0, 0.95)    # Dark gold
const ACCENT_COLOR = Color(1.0, 1.0, 0.0, 1.0)    # Yellow (holy)

var _pillar_label: Label
var _mult_label: Label
var _status_label: Label

func _ready() -> void:
	super._ready()
	_setup_pillar_ui()

func _setup_pillar_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "HOLY PILLAR!"
	subtitle_text = "12 divine pillars descend!"

	_pillar_label = Label.new()
	_pillar_label.text = "Pillars: 0 / 12"
	_pillar_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_pillar_label.add_theme_font_size_override("font_size", 22)
	_pillar_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_pillar_label)

	_mult_label = Label.new()
	_mult_label.text = "Total: x0.0"
	_mult_label.add_theme_color_override("font_color", Color.WHITE)
	_mult_label.add_theme_font_size_override("font_size", 24)
	_mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_mult_label)

	_status_label = Label.new()
	_status_label.text = "Divine judgment descending..."
	_status_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_status_label.add_theme_font_size_override("font_size", 13)
	_status_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_status_label)

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"pillar_start":
			_show_panel()
			_pillar_label.text = "Pillars: 0 / 12"
			_mult_label.text = "Total: x0.0"
			_status_label.text = "Divine judgment descending!"
			_play_start_animation()
		"pillar_hit":
			var idx = event_data.get("pillar_index", 0)
			var total = event_data.get("total_mult", 0.0)
			_pillar_label.text = "Pillars: %d / 12" % idx
			_mult_label.text = "Total: x%.1f" % total
			_status_label.text = "Pillar %d HITS!" % idx
		"pillar_miss":
			var idx = event_data.get("pillar_index", 0)
			_status_label.text = "Pillar %d missed..." % idx
		"pillar_result":
			var hit = event_data.get("hit_pillars", 0)
			var total = event_data.get("total_mult", 0.0)
			var perfect = event_data.get("is_perfect", false)
			_pillar_label.text = "Pillars: %d / 12" % hit
			_mult_label.text = "Total: x%.1f" % total
			if perfect:
				_status_label.text = "PERFECT HOLY PILLAR! Global x46.5 for 93s!"
				_play_perfect_animation()
			else:
				_status_label.text = "Holy Pillar complete! Hit %d/12" % hit
			_schedule_hide(4.0)
