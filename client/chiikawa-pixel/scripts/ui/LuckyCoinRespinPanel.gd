extends BaseLuckyPanel
# T223 Lucky Coin Respin Fish — Coin Respin mechanic
# BGaming "Shark & Spark Hold & Win" (2026-05-28): Hold & Win style
# 9-grid board, coins land and lock, each new coin resets 3 spins
# Bronze x10 / Silver x30 / Gold x80 / Diamond x200
# Full board bonus: +x500.0
# Global x42.5 for 86s (new all-time high)

const PANEL_COLOR = Color(0.8, 0.6, 0.0, 0.95)    # Dark gold (Hold & Win theme)
const ACCENT_COLOR = Color(1.0, 0.9, 0.2, 1.0)    # Bright gold

var _grid_labels: Array = []
var _spins_label: Label
var _total_label: Label
var _status_label: Label

func _ready() -> void:
	super._ready()
	_setup_respin_ui()

func _setup_respin_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "Coin Respin!"
	subtitle_text = "Hold & Win - Fill the board!"

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

	_status_label = Label.new()
	_status_label.text = "Bronze x10 | Silver x30 | Gold x80 | Diamond x200"
	_status_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_status_label.add_theme_font_size_override("font_size", 12)
	_status_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_status_label)

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"respin_start":
			_show_panel()
			_spins_label.text = "Spins left: 3"
			_total_label.text = "Total: x0.0"
			_status_label.text = "Hold & Win starting!"
			_play_start_animation()
		"coin_land":
			var coin_name = event_data.get("coin_name", "Bronze")
			var coin_mult = event_data.get("coin_mult", 10.0)
			var total = event_data.get("total_mult", 0.0)
			var spins = event_data.get("spins_left", 3)
			_spins_label.text = "Spins left: %d (reset!)" % spins
			_total_label.text = "Total: x%.0f" % total
			_status_label.text = "%s coin landed! +x%.0f" % [coin_name, coin_mult]
			_play_coin_animation()
		"spin_tick":
			var spins = event_data.get("spins_left", 0)
			_spins_label.text = "Spins left: %d" % spins
		"full_board_bonus":
			var bonus = event_data.get("bonus_mult", 500.0)
			var total = event_data.get("total_mult", 0.0)
			_status_label.text = "FULL BOARD! +x%.0f bonus!" % bonus
			_total_label.text = "Total: x%.0f" % total
			_play_perfect_animation()
		"respin_settle":
			var total = event_data.get("total_mult", 0.0)
			var reward = event_data.get("reward", 0)
			var full = event_data.get("full_board", false)
			title_text = "Coin Respin Done!"
			if full:
				subtitle_text = "FULL BOARD! x%.0f | Reward: %d" % [total, reward]
			else:
				subtitle_text = "x%.0f | Reward: %d" % [total, reward]
			_update_title()
			_hide_panel_delayed(3.0)
		"global_boost":
			var g_mult = event_data.get("global_mult", 42.5)
			var dur = event_data.get("duration", 86)
			title_text = "Global x%.1f for %ds!" % [g_mult, dur]
			subtitle_text = "Coin Respin - NEW ALL-TIME HIGH!"
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
