extends BaseLuckyPanel
# T221 Lucky Dice Bonus Fish — Dice Bonus mechanic
# BGaming "Shark & Spark Hold & Win" (2026-05-25): roll dice 3 times
# 1-3 dots x50.0, 4-5 dots x150.0, 6 dots x300.0
# Global x41.5 for 83s

const PANEL_COLOR = Color(0.9, 0.3, 0.1, 0.95)    # Orange-red (dice theme)
const ACCENT_COLOR = Color(1.0, 0.8, 0.0, 1.0)    # Gold

var _roll_index: int = 0
var _total_mult: float = 0.0
var _roll_label: Label
var _total_label: Label
var _dice_label: Label

func _ready() -> void:
	super._ready()
	_setup_dice_ui()

func _setup_dice_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "Dice Bonus!"
	subtitle_text = "Roll 3 dice - max x300.0 each!"

	_dice_label = Label.new()
	_dice_label.text = "Roll: -"
	_dice_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_dice_label.add_theme_font_size_override("font_size", 28)
	_dice_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_dice_label)

	_roll_label = Label.new()
	_roll_label.text = "Roll 0 / 3"
	_roll_label.add_theme_color_override("font_color", Color.WHITE)
	_roll_label.add_theme_font_size_override("font_size", 18)
	_roll_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_roll_label)

	_total_label = Label.new()
	_total_label.text = "Total: x0.0"
	_total_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.0, 1.0))
	_total_label.add_theme_font_size_override("font_size", 20)
	_total_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_total_label)

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"dice_start":
			_roll_index = 0
			_total_mult = 0.0
			_show_panel()
			_roll_label.text = "Roll 0 / 3"
			_total_label.text = "Total: x0.0"
			_dice_label.text = "Rolling..."
			_play_start_animation()
		"dice_roll":
			_roll_index = event_data.get("roll_index", 0)
			var dice_val = event_data.get("dice_value", 1)
			var roll_mult = event_data.get("roll_mult", 50.0)
			_total_mult = event_data.get("total_mult", 0.0)
			_dice_label.text = "Dice: %d  (+x%.0f)" % [dice_val, roll_mult]
			_roll_label.text = "Roll %d / 3" % _roll_index
			_total_label.text = "Total: x%.0f" % _total_mult
			_play_roll_animation()
		"dice_settle":
			var total = event_data.get("total_mult", 0.0)
			var reward = event_data.get("reward", 0)
			title_text = "Dice Settled!"
			subtitle_text = "Total x%.0f | Reward: %d" % [total, reward]
			_update_title()
			_hide_panel_delayed(3.0)
		"global_boost":
			var mult = event_data.get("global_mult", 41.5)
			var dur = event_data.get("duration", 83)
			title_text = "Global x%.1f for %ds!" % [mult, dur]
			subtitle_text = "Dice Bonus - New record!"
			_update_title()
			_play_perfect_animation()
			_hide_panel_delayed(5.0)

func _play_start_animation() -> void:
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(1.2, 1.2), 0.15)
	tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.15)

func _play_roll_animation() -> void:
	var tween = create_tween()
	tween.tween_property(_dice_label, "modulate", ACCENT_COLOR, 0.05)
	tween.tween_property(_dice_label, "modulate", Color.WHITE, 0.15)
