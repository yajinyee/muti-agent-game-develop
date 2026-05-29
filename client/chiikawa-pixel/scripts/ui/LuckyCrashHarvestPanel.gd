extends BaseLuckyPanel
# T227 Lucky Crash Harvest Fish — Crash Harvest Mechanic
# Lucky Fish AbraCadabra (2026-05): Crash mechanic
# Multiplier rises continuously, cash out before it crashes!
# Perfect harvest (>=50x, no crash): Global x44.5 for 89s
# Max multiplier: x1000

const PANEL_COLOR = Color(0.5, 0.1, 0.0, 0.95)    # Dark red-orange (Crash theme)
const ACCENT_COLOR = Color(1.0, 0.3, 0.0, 1.0)    # Bright red-orange
const SAFE_COLOR = Color(0.2, 0.9, 0.2, 1.0)      # Green (safe zone)
const DANGER_COLOR = Color(1.0, 0.2, 0.2, 1.0)    # Red (danger zone)

var _mult_label: Label
var _status_label: Label
var _risk_label: Label

func _ready() -> void:
	super._ready()
	_setup_crash_ui()

func _setup_crash_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "Crash Harvest!"
	subtitle_text = "Cash out before it crashes!"

	_mult_label = Label.new()
	_mult_label.text = "x1.0"
	_mult_label.add_theme_color_override("font_color", SAFE_COLOR)
	_mult_label.add_theme_font_size_override("font_size", 36)
	_mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_mult_label)

	_risk_label = Label.new()
	_risk_label.text = "Risk: LOW"
	_risk_label.add_theme_color_override("font_color", SAFE_COLOR)
	_risk_label.add_theme_font_size_override("font_size", 16)
	_risk_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_risk_label)

	_status_label = Label.new()
	_status_label.text = "Multiplier rising... Perfect harvest >= x50!"
	_status_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_status_label.add_theme_font_size_override("font_size", 12)
	_status_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_status_label)

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"crash_start":
			_show_panel()
			_mult_label.text = "x1.0"
			_mult_label.modulate = SAFE_COLOR
			_risk_label.text = "Risk: LOW"
			_risk_label.modulate = SAFE_COLOR
			_status_label.text = "Multiplier rising... Perfect >= x50!"
			_play_start_animation()
		"mult_tick":
			var mult = event_data.get("current_mult", 1.0)
			_mult_label.text = "x%.1f" % mult
			# Color changes with risk
			if mult >= 50.0:
				_mult_label.modulate = DANGER_COLOR
				_risk_label.text = "Risk: EXTREME"
				_risk_label.modulate = DANGER_COLOR
			elif mult >= 20.0:
				_mult_label.modulate = Color(1.0, 0.5, 0.0, 1.0)
				_risk_label.text = "Risk: HIGH"
				_risk_label.modulate = Color(1.0, 0.5, 0.0, 1.0)
			elif mult >= 10.0:
				_mult_label.modulate = Color(1.0, 0.8, 0.0, 1.0)
				_risk_label.text = "Risk: MEDIUM"
				_risk_label.modulate = Color(1.0, 0.8, 0.0, 1.0)
			else:
				_mult_label.modulate = SAFE_COLOR
				_risk_label.text = "Risk: LOW"
				_risk_label.modulate = SAFE_COLOR
		"crashed":
			var crash_mult = event_data.get("crash_mult", 0.0)
			var harvest = event_data.get("harvest_mult", 1.0)
			_mult_label.text = "CRASHED!"
			_mult_label.modulate = DANGER_COLOR
			_status_label.text = "Crashed at x%.1f! Consolation: x%.1f" % [crash_mult, harvest]
			_play_crash_animation()
		"harvested":
			var harvest = event_data.get("harvest_mult", 0.0)
			var perfect = event_data.get("perfect", false)
			_mult_label.text = "x%.1f" % harvest
			if perfect:
				_status_label.text = "PERFECT HARVEST! x%.1f!" % harvest
				_play_perfect_animation()
			else:
				_status_label.text = "Harvested: x%.1f" % harvest
		"crash_settle":
			var harvest = event_data.get("harvest_mult", 0.0)
			var reward = event_data.get("reward", 0)
			var crashed = event_data.get("crashed", false)
			var perfect = event_data.get("perfect_harvest", false)
			title_text = "Crash Harvest Done!"
			if perfect:
				subtitle_text = "PERFECT! x%.1f | Reward: %d" % [harvest, reward]
			elif crashed:
				subtitle_text = "Crashed! x%.1f | Reward: %d" % [harvest, reward]
			else:
				subtitle_text = "x%.1f | Reward: %d" % [harvest, reward]
			_update_title()
			_hide_panel_delayed(3.0)
		"global_boost":
			var g_mult = event_data.get("global_mult", 44.5)
			var dur = event_data.get("duration", 89)
			title_text = "Global x%.1f for %ds!" % [g_mult, dur]
			subtitle_text = "Perfect Harvest - All-time high!"
			_update_title()
			_play_perfect_animation()
			_hide_panel_delayed(6.0)

func _play_start_animation() -> void:
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(1.15, 1.15), 0.12)
	tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.12)

func _play_crash_animation() -> void:
	var tween = create_tween()
	for i in range(3):
		tween.tween_property(self, "modulate", Color(1.0, 0.2, 0.2, 1.0), 0.05)
		tween.tween_property(self, "modulate", Color.WHITE, 0.1)
