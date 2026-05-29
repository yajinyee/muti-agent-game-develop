extends BaseLuckyPanel
# T225 Lucky Cascade Lock Fish — Cascading Wins + Locked Multipliers
# BGaming "Shark & Spark Hold & Win" (2026-05-28)
# 8 waves of cascading wins, Pearl symbols x2-x10 bonus
# Perfect 8 waves: +x50.0 bonus
# Global x43.5 for 87s (surpasses T224 x43.0)

const PANEL_COLOR = Color(0.0, 0.4, 0.8, 0.95)    # Deep blue (Cascade theme)
const ACCENT_COLOR = Color(0.0, 0.9, 1.0, 1.0)    # Cyan

var _wave_label: Label
var _total_label: Label
var _locked_label: Label
var _status_label: Label

func _ready() -> void:
	super._ready()
	_setup_cascade_ui()

func _setup_cascade_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "Cascade Lock!"
	subtitle_text = "8 waves of cascading wins!"

	_wave_label = Label.new()
	_wave_label.text = "Wave: 0 / 8"
	_wave_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_wave_label.add_theme_font_size_override("font_size", 22)
	_wave_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_wave_label)

	_total_label = Label.new()
	_total_label.text = "Total: x0.0"
	_total_label.add_theme_color_override("font_color", Color.WHITE)
	_total_label.add_theme_font_size_override("font_size", 20)
	_total_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_total_label)

	_locked_label = Label.new()
	_locked_label.text = "Locked: 0 multipliers"
	_locked_label.add_theme_color_override("font_color", Color(0.9, 0.7, 0.2, 1.0))
	_locked_label.add_theme_font_size_override("font_size", 14)
	_locked_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_locked_label)

	_status_label = Label.new()
	_status_label.text = "Pearl symbols give x2-x10 bonus!"
	_status_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_status_label.add_theme_font_size_override("font_size", 12)
	_status_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_status_label)

var _locked_count: int = 0

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"cascade_start":
			_locked_count = 0
			_show_panel()
			_wave_label.text = "Wave: 0 / 8"
			_total_label.text = "Total: x0.0"
			_locked_label.text = "Locked: 0 multipliers"
			_status_label.text = "Cascade starting!"
			_play_start_animation()
		"wave_hit":
			var wave = event_data.get("wave", 1)
			var wave_mult = event_data.get("wave_mult", 0.0)
			var total = event_data.get("total_mult", 0.0)
			var has_pearl = event_data.get("has_pearl", false)
			var pearl_mult = event_data.get("pearl_mult", 1.0)
			_locked_count += 1
			_wave_label.text = "Wave: %d / 8" % wave
			_total_label.text = "Total: x%.0f" % total
			_locked_label.text = "Locked: %d multipliers" % _locked_count
			if has_pearl:
				_status_label.text = "Wave %d: x%.1f (Pearl x%.0f!)" % [wave, wave_mult, pearl_mult]
			else:
				_status_label.text = "Wave %d: x%.1f locked!" % [wave, wave_mult]
			_play_wave_animation()
		"wave_miss":
			var wave = event_data.get("wave", 1)
			_wave_label.text = "Wave: %d / 8 (miss)" % wave
			_status_label.text = "Wave %d missed..." % wave
		"perfect_cascade":
			var bonus = event_data.get("bonus_mult", 50.0)
			var total = event_data.get("total_mult", 0.0)
			_status_label.text = "PERFECT CASCADE! +x%.0f bonus!" % bonus
			_total_label.text = "Total: x%.0f" % total
			_play_perfect_animation()
		"cascade_settle":
			var total = event_data.get("total_mult", 0.0)
			var reward = event_data.get("reward", 0)
			var perfect = event_data.get("perfect_cascade", false)
			title_text = "Cascade Done!"
			if perfect:
				subtitle_text = "PERFECT! x%.0f | Reward: %d" % [total, reward]
			else:
				subtitle_text = "x%.0f | Reward: %d" % [total, reward]
			_update_title()
			_hide_panel_delayed(3.0)
		"global_boost":
			var g_mult = event_data.get("global_mult", 43.5)
			var dur = event_data.get("duration", 87)
			title_text = "Global x%.1f for %ds!" % [g_mult, dur]
			subtitle_text = "Cascade Lock - All-time high!"
			_update_title()
			_play_perfect_animation()
			_hide_panel_delayed(6.0)

func _play_start_animation() -> void:
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(1.15, 1.15), 0.12)
	tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.12)

func _play_wave_animation() -> void:
	var tween = create_tween()
	tween.tween_property(_wave_label, "modulate", ACCENT_COLOR, 0.05)
	tween.tween_property(_wave_label, "modulate", Color.WHITE, 0.2)
