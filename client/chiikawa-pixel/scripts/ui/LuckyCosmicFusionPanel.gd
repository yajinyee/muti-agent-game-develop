extends BaseLuckyPanel
# T228 Lucky Cosmic Fusion Fish — Ultimate Fusion Mechanic
# 4 Phases: Coin Respin + Cascade Lock + Legend Awaken + Full Field Clear
# Global x45.0 for 90s (NEW ALL-TIME HIGH, surpasses T227 x44.5)

const PANEL_COLOR = Color(0.3, 0.0, 0.4, 0.95)    # Deep cosmic purple
const ACCENT_COLOR = Color(1.0, 0.0, 1.0, 1.0)    # Magenta (cosmic)
const PHASE_COLORS = [
	Color(0.8, 0.6, 0.0, 1.0),   # Phase 1: Gold (Coin Respin)
	Color(0.0, 0.6, 1.0, 1.0),   # Phase 2: Blue (Cascade)
	Color(1.0, 0.3, 0.0, 1.0),   # Phase 3: Orange (Legend)
	Color(1.0, 0.0, 0.5, 1.0),   # Phase 4: Pink (Cosmic Clear)
]

var _phase_label: Label
var _total_label: Label
var _detail_label: Label
var _phase_bar: ColorRect

func _ready() -> void:
	super._ready()
	_setup_fusion_ui()

func _setup_fusion_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "COSMIC FUSION!"
	subtitle_text = "4 phases of ultimate destruction!"

	_phase_label = Label.new()
	_phase_label.text = "Phase: 0 / 4"
	_phase_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_phase_label.add_theme_font_size_override("font_size", 22)
	_phase_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_phase_label)

	_total_label = Label.new()
	_total_label.text = "Total: x0.0"
	_total_label.add_theme_color_override("font_color", Color.WHITE)
	_total_label.add_theme_font_size_override("font_size", 24)
	_total_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_total_label)

	_detail_label = Label.new()
	_detail_label.text = "Coin Respin → Cascade → Legend → Clear"
	_detail_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_detail_label.add_theme_font_size_override("font_size", 12)
	_detail_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_detail_label)

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"fusion_start":
			_show_panel()
			_phase_label.text = "Phase: 0 / 4"
			_total_label.text = "Total: x0.0"
			_detail_label.text = "COSMIC FUSION INITIATED!"
			_play_start_animation()
		"phase_start":
			var phase = event_data.get("phase", 1)
			var name = event_data.get("name", "")
			_phase_label.text = "Phase %d: %s" % [phase, name]
			if phase >= 1 and phase <= 4:
				_phase_label.modulate = PHASE_COLORS[phase - 1]
			_detail_label.text = "Phase %d starting..." % phase
			_play_phase_animation(phase)
		"phase_complete":
			var phase = event_data.get("phase", 1)
			var phase_mult = event_data.get("phase_mult", 0.0)
			var total = event_data.get("total_mult", 0.0)
			_total_label.text = "Total: x%.0f" % total
			_detail_label.text = "Phase %d done! +x%.0f" % [phase, phase_mult]
		"phase1_coin":
			var coin_name = event_data.get("coin_name", "")
			var coin_mult = event_data.get("coin_mult", 0.0)
			var phase_mult = event_data.get("phase_mult", 0.0)
			_detail_label.text = "Coin: %s +x%.0f (Phase: x%.0f)" % [coin_name, coin_mult, phase_mult]
		"phase2_wave":
			var wave = event_data.get("wave", 1)
			var wave_mult = event_data.get("wave_mult", 0.0)
			_detail_label.text = "Wave %d: +x%.1f" % [wave, wave_mult]
		"phase3_awaken":
			var round = event_data.get("round", 1)
			var round_mult = event_data.get("round_mult", 0.0)
			_detail_label.text = "Awaken Round %d: +x%.1f" % [round, round_mult]
		"phase4_clear":
			var cleared = event_data.get("cleared_count", 0)
			var phase_mult = event_data.get("phase_mult", 0.0)
			var total = event_data.get("total_mult", 0.0)
			_total_label.text = "Total: x%.0f" % total
			_detail_label.text = "COSMIC CLEAR! %d targets! +x%.0f" % [cleared, phase_mult]
			_play_perfect_animation()
		"fusion_settle":
			var total = event_data.get("total_mult", 0.0)
			var reward = event_data.get("reward", 0)
			title_text = "COSMIC FUSION COMPLETE!"
			subtitle_text = "x%.0f total! Reward: %d" % [total, reward]
			_update_title()
			_hide_panel_delayed(4.0)
		"global_boost":
			var g_mult = event_data.get("global_mult", 45.0)
			var dur = event_data.get("duration", 90)
			title_text = "Global x%.1f for %ds!" % [g_mult, dur]
			subtitle_text = "COSMIC FUSION - NEW ALL-TIME HIGH!"
			_update_title()
			_play_perfect_animation()
			_hide_panel_delayed(8.0)

func _play_start_animation() -> void:
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(1.3, 1.3), 0.2)
	tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.2)
	# Rainbow flash
	tween.tween_property(self, "modulate", Color(1.0, 0.0, 1.0, 1.0), 0.1)
	tween.tween_property(self, "modulate", Color.WHITE, 0.2)

func _play_phase_animation(phase: int) -> void:
	if phase >= 1 and phase <= 4:
		var tween = create_tween()
		tween.tween_property(self, "modulate", PHASE_COLORS[phase - 1], 0.1)
		tween.tween_property(self, "modulate", Color.WHITE, 0.3)
