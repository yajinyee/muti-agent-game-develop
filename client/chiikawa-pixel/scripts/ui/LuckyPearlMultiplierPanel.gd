extends BaseLuckyPanel
# T219 幸運珍珠倍率魚 — Pearl Multiplier 機制
# BGaming Shark & Spark Hold & Win（2026-05）：場上每個目標都有珍珠倍率（×1-×100）
# 全服 ×40.0 加成 80 秒（新里程碑：全服 ×40.0）

const PANEL_COLOR = Color(1.0, 0.843, 0.0, 0.92)    # 金色
const ACCENT_COLOR = Color(1.0, 1.0, 0.5, 1.0)       # 亮金

var _pearl_count: int = 0
var _collected: int = 0
var _pearl_label: Label
var _collected_label: Label
var _milestone_label: Label

func _ready() -> void:
	super._ready()
	_setup_pearl_ui()

func _setup_pearl_ui() -> void:
	panel_color = PANEL_COLOR
	title_text = "🦪 珍珠降臨！"
	subtitle_text = "每個目標獲得珍珠倍率（×1-×100）"

	_pearl_label = Label.new()
	_pearl_label.text = "珍珠目標：0 個"
	_pearl_label.add_theme_color_override("font_color", ACCENT_COLOR)
	_pearl_label.add_theme_font_size_override("font_size", 18)
	_pearl_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_pearl_label)

	_collected_label = Label.new()
	_collected_label.text = "已收集：0 個"
	_collected_label.add_theme_color_override("font_color", Color.WHITE)
	_collected_label.add_theme_font_size_override("font_size", 16)
	_collected_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_collected_label)

	_milestone_label = Label.new()
	_milestone_label.text = "🏆 目標：全服 ×40.0"
	_milestone_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.0, 1.0))
	_milestone_label.add_theme_font_size_override("font_size", 14)
	_milestone_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	content_container.add_child(_milestone_label)

func handle_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"pearl_assign":
			_pearl_count = event_data.get("pearl_count", 0)
			_show_panel()
			_pearl_label.text = "珍珠目標：%d 個" % _pearl_count
			_play_pearl_animation()
		"pearl_on_target":
			var pearl_mult = event_data.get("pearl_mult", 1.0)
			# 高倍率珍珠特別顯示
			if pearl_mult >= 50.0:
				_show_high_pearl(pearl_mult)
		"pearl_perfect":
			var boost_mult = event_data.get("global_boost_mult", 40.0)
			var boost_secs = event_data.get("global_boost_secs", 80)
			_show_perfect_result(boost_mult, boost_secs)
		"pearl_end":
			_hide_panel_delayed(2.0)

func _show_high_pearl(mult: float) -> void:
	var tween = create_tween()
	_milestone_label.text = "💎 高倍珍珠 ×%.0f！" % mult
	tween.tween_property(_milestone_label, "modulate", Color(1.0, 0.8, 0.0, 1.0), 0.1)
	tween.tween_property(_milestone_label, "modulate", Color.WHITE, 0.3)

func _play_pearl_animation() -> void:
	var tween = create_tween().set_loops(4)
	tween.tween_property(self, "modulate", Color(1.0, 1.0, 0.8, 1.0), 0.2)
	tween.tween_property(self, "modulate", Color.WHITE, 0.2)

func _show_perfect_result(boost_mult: float, boost_secs: int) -> void:
	title_text = "🦪🌟 珍珠完美收集！"
	subtitle_text = "全服 ×%.1f 加成 %d 秒！（新里程碑）" % [boost_mult, boost_secs]
	_update_title()
	_play_perfect_animation()
	_hide_panel_delayed(6.0)
