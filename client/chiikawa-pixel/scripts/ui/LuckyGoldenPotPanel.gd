extends BaseLuckyPanel
# T224 Lucky Golden Pot Fish — Gold Blitz™ Cash Collection
# Games Global "Fishin' Pots of Gold Gold Blitz Ultimate" (2026-05-28)
# 12-grid golden pot, Enhanced Respin (new coin resets 3 spins)
# Copper x5 / Silver x20 / Gold x60 / Platinum x150 / Diamond x200
# Full pot bonus: +x300.0
# Global x43.0 for 88s (new all-time high, surpasses T223 x42.5)

const PANEL_COLOR = Color(0.7, 0.5, 0.0, 0.95)    # Deep gold (Gold Blitz theme)
const ACCENT_COLOR = Color(1.0, 0.85, 0.0, 1.0)   # Bright gold

var _spins_label: Label
var _total_label: Label
var _status_label: Label
var _grid_label: Label

func _ready() -> void:
	super._ready()
	_setup_pot_ui()

func _setup_pot_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "Gold Blitz!"
	subtitle_text = "Cash Collection - Fill the pot!"

	_spins_label = Label.new()
	_spins_label.text = "Spins left: 3"
	_spins_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_spins_label.add_theme_font_size_override("font_size", 22)
	_spins_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_spins_label)

	_total_label = Label.new()
	_total_label.text = "Total: x0.0"
	_total_label.add_theme_color_override("font_color", Color.WHITE)
	_total_label.add_theme_font_size_override("font_size", 20)
	_total_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_total_label)

	_grid_label = Label.new()
	_grid_label.text = "0 / 12 slots filled"
	_grid_label.add_theme_color_override("font_color", Color(0.9, 0.7, 0.2, 1.0))
	_grid_label.add_theme_font_size_override("font_size", 16)
	_grid_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_grid_label)

	_status_label = Label.new()
	_status_label.text = "Copper x5 | Silver x20 | Gold x60 | Platinum x150 | Diamond x200"
	_status_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_status_label.add_theme_font_size_override("font_size", 11)
	_status_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_status_label)

var _filled_slots: int = 0

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"pot_start":
			_filled_slots = 0
			_show_panel()
			_spins_label.text = "Spins left: 3"
			_total_label.text = "Total: x0.0"
			_grid_label.text = "0 / 12 slots filled"
			_status_label.text = "Gold Blitz starting!"
			_play_start_animation()
		"coin_land":
			var coin_name = event_data.get("coin_name", "Copper")
			var coin_mult = event_data.get("coin_mult", 5.0)
			var total = event_data.get("total_mult", 0.0)
			var spins = event_data.get("spins_left", 3)
			_filled_slots += 1
			_spins_label.text = "Spins left: %d (reset!)" % spins
			_total_label.text = "Total: x%.0f" % total
			_grid_label.text = "%d / 12 slots filled" % _filled_slots
			_status_label.text = "%s coin! +x%.0f" % [coin_name, coin_mult]
			_play_coin_animation()
		"spin_tick":
			var spins = event_data.get("spins_left", 0)
			_spins_label.text = "Spins left: %d" % spins
		"full_pot_bonus":
			var bonus = event_data.get("bonus_mult", 300.0)
			var total = event_data.get("total_mult", 0.0)
			_status_label.text = "FULL POT! +x%.0f Gold Blitz bonus!" % bonus
			_total_label.text = "Total: x%.0f" % total
			_grid_label.text = "12 / 12 FULL!"
			_play_perfect_animation()
		"pot_settle":
			var total = event_data.get("total_mult", 0.0)
			var reward = event_data.get("reward", 0)
			var full = event_data.get("full_pot", false)
			title_text = "Gold Blitz Done!"
			if full:
				subtitle_text = "FULL POT! x%.0f | Reward: %d" % [total, reward]
			else:
				subtitle_text = "x%.0f | Reward: %d" % [total, reward]
			_update_title()
			_hide_panel_delayed(3.0)
		"global_boost":
			var g_mult = event_data.get("global_mult", 43.0)
			var dur = event_data.get("duration", 88)
			title_text = "Global x%.1f for %ds!" % [g_mult, dur]
			subtitle_text = "Gold Blitz - NEW ALL-TIME HIGH!"
			_update_title()
			_play_perfect_animation()
			_hide_panel_delayed(6.0)

func _play_start_animation() -> void:
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(1.2, 1.2), 0.15)
	tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.15)

func _play_coin_animation() -> void:
	var tween = create_tween()
	tween.tween_property(_total_label, "modulate", ACCENT_COLOR, 0.05)
	tween.tween_property(_total_label, "modulate", Color.WHITE, 0.15)
