extends BaseLuckyPanel
# T230 Lucky Super Chain Fish — Super Chain Mechanic
# 5 chain reactions, each x80.0
# Perfect (>=5 chains) -> Global x46.0 for 92s (surpasses T229 x45.5)

const PANEL_COLOR = Color(0.0, 0.3, 0.4, 0.95)    # Deep cyan
const ACCENT_COLOR = Color(0.0, 1.0, 1.0, 1.0)    # Cyan (chain)

var _chain_label: Label
var _mult_label: Label
var _status_label: Label

func _ready() -> void:
	super._ready()
	_setup_chain_ui()

func _setup_chain_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "SUPER CHAIN!"
	subtitle_text = "5 chain reactions! Each x80.0!"

	_chain_label = Label.new()
	_chain_label.text = "Chain: 0 / 5"
	_chain_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_chain_label.add_theme_font_size_override("font_size", 22)
	_chain_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_chain_label)

	_mult_label = Label.new()
	_mult_label.text = "Total: x0.0"
	_mult_label.add_theme_color_override("font_color", Color.WHITE)
	_mult_label.add_theme_font_size_override("font_size", 24)
	_mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_mult_label)

	_status_label = Label.new()
	_status_label.text = "Chain reaction starting..."
	_status_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_status_label.add_theme_font_size_override("font_size", 13)
	_status_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_status_label)

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"chain_start":
			_show_panel()
			_chain_label.text = "Chain: 0 / 5"
			_mult_label.text = "Total: x0.0"
			_status_label.text = "Chain reaction starting!"
			_play_start_animation()
		"chain_hit":
			var idx = event_data.get("chain_index", 0)
			var total = event_data.get("total_mult", 0.0)
			var is_bonus = event_data.get("is_bonus", false)
			_chain_label.text = "Chain: %d / 5" % idx
			_mult_label.text = "Total: x%.1f" % total
			if is_bonus:
				_status_label.text = "BONUS CHAIN %d! x80.0!" % idx
			else:
				_status_label.text = "Chain %d! x80.0!" % idx
		"chain_result":
			var count = event_data.get("chain_count", 0)
			var total = event_data.get("total_mult", 0.0)
			var perfect = event_data.get("is_perfect", false)
			_chain_label.text = "Chain: %d / 5" % count
			_mult_label.text = "Total: x%.1f" % total
			if perfect:
				_status_label.text = "SUPER CHAIN BURST! Global x46.0 for 92s!"
				_play_perfect_animation()
			else:
				_status_label.text = "Chain complete!"
			_schedule_hide(4.0)
