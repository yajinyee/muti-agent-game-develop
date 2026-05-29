extends BaseLuckyPanel
# T222 Lucky Dual Bonus Fish — Dual Bonus mechanic
# BGaming "Fishing Club 2" (2026-04): choose Bonus A (Coin Collect) or Bonus B (Risk Wheel)
# Bonus A: 5 coins x80.0, global x41.8 for 84s
# Bonus B: Risk Wheel max x500.0, global x42.0 for 85s

const PANEL_COLOR = Color(0.2, 0.6, 1.0, 0.95)    # Blue (dual choice theme)
const ACCENT_COLOR = Color(0.0, 1.0, 0.8, 1.0)    # Cyan

var _status_label: Label
var _result_label: Label
var _choice_timer: float = 0.0
var _choice_active: bool = false

func _ready() -> void:
	super._ready()
	_setup_dual_ui()

func _setup_dual_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "Dual Bonus!"
	subtitle_text = "Choose your bonus!"

	_status_label = Label.new()
	_status_label.text = "Waiting for choice..."
	_status_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_status_label.add_theme_font_size_override("font_size", 20)
	_status_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_status_label)

	_result_label = Label.new()
	_result_label.text = ""
	_result_label.add_theme_color_override("font_color", Color.WHITE)
	_result_label.add_theme_font_size_override("font_size", 16)
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_result_label)

func _process(delta: float) -> void:
	if not _choice_active:
		return
	_choice_timer -= delta
	if _choice_timer > 0:
		_status_label.text = "Choose in %.0fs..." % _choice_timer
	else:
		_choice_active = false

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"choice_start":
			_choice_timer = float(event_data.get("timeout", 10))
			_choice_active = true
			_show_panel()
			_status_label.text = "Choose in 10s..."
			_result_label.text = "A: Coin Collect (5x x80)\nB: Risk Wheel (max x500)"
			_play_start_animation()
		"coin_collect":
			var idx = event_data.get("coin_index", 0)
			var mult = event_data.get("coin_mult", 80.0)
			var total = event_data.get("total_mult", 0.0)
			_status_label.text = "Coin %d: x%.0f" % [idx, mult]
			_result_label.text = "Total: x%.0f" % total
		"bonus_a_settle":
			var total = event_data.get("total_mult", 0.0)
			var reward = event_data.get("reward", 0)
			_choice_active = false
			title_text = "Coin Collect!"
			subtitle_text = "x%.0f | Reward: %d" % [total, reward]
			_update_title()
			_hide_panel_delayed(3.0)
		"bonus_b_settle":
			var mult = event_data.get("wheel_mult", 0.0)
			var reward = event_data.get("reward", 0)
			_choice_active = false
			title_text = "Risk Wheel!"
			subtitle_text = "x%.0f | Reward: %d" % [mult, reward]
			_update_title()
			_hide_panel_delayed(3.0)
		"global_boost":
			var g_mult = event_data.get("global_mult", 42.0)
			var dur = event_data.get("duration", 85)
			var bonus_type = event_data.get("bonus_type", "A")
			title_text = "Global x%.1f for %ds!" % [g_mult, dur]
			subtitle_text = "Dual Bonus %s - New record!" % bonus_type
			_update_title()
			_play_perfect_animation()
			_hide_panel_delayed(5.0)

func _play_start_animation() -> void:
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(1.15, 1.15), 0.12)
	tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.12)
