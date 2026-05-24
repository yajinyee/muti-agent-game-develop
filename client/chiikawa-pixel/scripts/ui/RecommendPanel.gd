## RecommendPanel.gd ???箸?刻?Ｘ嚗AY-110嚗?
## ?寞??拙振銵璅∪?嚗＊蝷箏犖???脣遣霅?
extends CanvasLayer

var _panel: PanelContainer
var _recs_container: VBoxContainer
var _loading_label: Label

func _ready():
	layer = 86
	_build_ui()
	hide()

func _build_ui():
	# ???
	var bg = ColorRect.new()
	bg.color = Color(0, 0, 0, 0.5)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(bg)
	bg.gui_input.connect(func(e): if e is InputEventMouseButton and e.pressed: hide())

	# 銝駁?選??喳皛嚗?
	_panel = PanelContainer.new()
	_panel.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_panel.custom_minimum_size = Vector2(320, 400)
	_panel.offset_left = -330
	_panel.offset_top = -200
	_panel.offset_right = -10
	_panel.offset_bottom = 200
	add_child(_panel)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 10)
	_panel.add_child(vbox)

	# 璅???
	var title_row = HBoxContainer.new()
	vbox.add_child(title_row)

	var title_lbl = Label.new()
	title_lbl.text = "? ?箸撱箄降"
	title_lbl.add_theme_font_size_override("font_size", 18)
	title_lbl.add_theme_color_override("font_color", Color(0.3, 0.9, 1.0))
	title_lbl.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	title_row.add_child(title_lbl)

	var close_btn = Button.new()
	close_btn.text = "??"
	close_btn.custom_minimum_size = Vector2(28, 28)
	close_btn.pressed.connect(hide)
	title_row.add_child(close_btn)

	var subtitle = Label.new()
	subtitle.text = "?寞?雿??蝧嚗雿?頨急??遣霅?"
	subtitle.add_theme_font_size_override("font_size", 11)
	subtitle.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	subtitle.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	vbox.add_child(subtitle)

	var sep = HSeparator.new()
	vbox.add_child(sep)

	# 頛銝剜?蝷?
	_loading_label = Label.new()
	_loading_label.text = "????銝?.."
	_loading_label.add_theme_font_size_override("font_size", 14)
	_loading_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	_loading_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(_loading_label)

	# ?刻?”
	_recs_container = VBoxContainer.new()
	_recs_container.add_theme_constant_override("separation", 8)
	_recs_container.size_flags_vertical = Control.SIZE_EXPAND_FILL
	vbox.add_child(_recs_container)

	# 摨??
	var refresh_btn = Button.new()
	refresh_btn.text = "?? ???"
	refresh_btn.pressed.connect(_on_refresh_pressed)
	vbox.add_child(refresh_btn)

func show_panel():
	show()
	_loading_label.show()
	_clear_recs()
	GameManager.request_recommendations()

func update_recommendations(data: Dictionary):
	_loading_label.hide()
	_clear_recs()

	var recs = data.get("recommendations", [])
	if recs.is_empty():
		var empty_lbl = Label.new()
		empty_lbl.text = "???桀?瘝??孵撱箄降\n蝜潛?靽?嚗?"
		empty_lbl.add_theme_font_size_override("font_size", 13)
		empty_lbl.add_theme_color_override("font_color", Color(0.5, 1.0, 0.5))
		empty_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		empty_lbl.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
		_recs_container.add_child(empty_lbl)
		return

	for rec in recs:
		var card = _create_rec_card(rec)
		_recs_container.add_child(card)

func _create_rec_card(rec: Dictionary) -> PanelContainer:
	var panel = PanelContainer.new()

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.2, 0.3, 0.9)
	style.border_width_left = 3
	style.border_color = _get_priority_color(rec.get("priority", 3))
	style.corner_radius_top_left = 6
	style.corner_radius_top_right = 6
	style.corner_radius_bottom_left = 6
	style.corner_radius_bottom_right = 6
	panel.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 4)
	panel.add_child(vbox)

	# 璅?銵?
	var title_row = HBoxContainer.new()
	vbox.add_child(title_row)

	var icon_lbl = Label.new()
	icon_lbl.text = rec.get("icon", "?")
	icon_lbl.add_theme_font_size_override("font_size", 18)
	title_row.add_child(icon_lbl)

	var title_lbl = Label.new()
	title_lbl.text = rec.get("title", "")
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	title_lbl.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	title_row.add_child(title_lbl)

	# 靽∪?摨?
	var confidence = rec.get("confidence", 0.5)
	var conf_lbl = Label.new()
	conf_lbl.text = "%.0f%%" % (confidence * 100)
	conf_lbl.add_theme_font_size_override("font_size", 11)
	conf_lbl.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	title_row.add_child(conf_lbl)

	# ?膩
	var desc_lbl = Label.new()
	desc_lbl.text = rec.get("description", "")
	desc_lbl.add_theme_font_size_override("font_size", 11)
	desc_lbl.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	desc_lbl.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	vbox.add_child(desc_lbl)

	# 憒??遣霅唳?瘜函?蝝?憿舐內敹恍?????
	var target_bet = rec.get("target_bet_lv", 0)
	if target_bet > 0:
		var btn = Button.new()
		btn.text = "????LV%d" % target_bet
		btn.custom_minimum_size = Vector2(0, 28)
		btn.pressed.connect(func(): _on_switch_bet_pressed(target_bet))
		vbox.add_child(btn)

	return panel

func _get_priority_color(priority: int) -> Color:
	match priority:
		1: return Color(1.0, 0.3, 0.3)  # 蝝 = 擃??
		2: return Color(1.0, 0.7, 0.2)  # 璈 = 銝剖??
		_: return Color(0.3, 0.7, 1.0)  # ? = 雿??

func _clear_recs():
	for child in _recs_container.get_children():
		child.queue_free()

func _on_refresh_pressed():
	_loading_label.show()
	_clear_recs()
	GameManager.request_recommendations()

func _on_switch_bet_pressed(bet_level: int):
	GameManager.send_bet_change(bet_level)
	hide()
