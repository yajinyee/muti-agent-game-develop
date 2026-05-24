## PlayerJourneyPanel.gd ???拙振???銵冽嚗AY-108嚗?
## 璆剔?靘?嚗verymatrix.com嚗?026嚗Ⅱ隤???拙振????2026 撟?iGaming ?詨?蝡嗥??
## ?嚗?
##   - ?葉憿舐內??脣漲嚗遙???餃/鞈賢迤/VIP/頧/?行?鞈踝?
##   - 霈摰嗡??潛??啜?憭抵???暻潦?
##   - 暺????桀?湔頝唾?撠??Ｘ
##   - ?喃?閫??? ????嚗_index=85
extends CanvasLayer

# ---- 蝭暺???----
var _panel: PanelContainer
var _scroll: ScrollContainer
var _content: VBoxContainer
var _toggle_btn: Button

# ---- ??翰??----
var _mission_data: Array = []
var _login_streak: int = 0
var _login_next_milestone_days: int = 0
var _login_days_to_next: int = 0
var _season_level: int = 0
var _season_points: int = 0
var _season_points_to_next: int = 0
var _season_progress: float = 0.0
var _vip_level: int = 0
var _vip_name: String = ""
var _daily_spin_available: bool = false
var _tournament_rank: int = 0
var _tournament_points: int = 0
var _daily_tournament_rank: int = 0
var _daily_tournament_points: int = 0

# ---- ?憛?暺??剁??冽?湔嚗?---
var _mission_section: Control = null
var _login_section: Control = null
var _season_section: Control = null
var _vip_section: Control = null
var _spin_section: Control = null
var _tournament_section: Control = null

func _ready() -> void:
	layer = 85
	_build_toggle_button()
	_build_panel()
	_connect_signals()

# ---- UI 撱箸? ----

func _build_toggle_button() -> void:
	_toggle_btn = Button.new()
	_toggle_btn.text = "??"
	_toggle_btn.custom_minimum_size = Vector2(40, 40)
	_toggle_btn.position = Vector2(1230, 4)
	_toggle_btn.add_theme_font_size_override("font_size", 18)
	_toggle_btn.pressed.connect(_on_toggle_pressed)
	add_child(_toggle_btn)

func _build_panel() -> void:
	_panel = PanelContainer.new()
	_panel.custom_minimum_size = Vector2(360, 560)
	_panel.position = Vector2(910, 50)
	_panel.visible = false
	add_child(_panel)

	var vbox := VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 0)
	_panel.add_child(vbox)

	# 璅???
	var header := HBoxContainer.new()
	header.custom_minimum_size = Vector2(0, 44)
	var header_bg := ColorRect.new()
	header_bg.color = Color(0.05, 0.1, 0.25, 1.0)
	header_bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	header.add_child(header_bg)

	var title_lbl := Label.new()
	title_lbl.text = "  ?? 隞??"
	title_lbl.add_theme_font_size_override("font_size", 16)
	title_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	title_lbl.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	title_lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	header.add_child(title_lbl)

	var close_btn := Button.new()
	close_btn.text = "??"
	close_btn.custom_minimum_size = Vector2(40, 40)
	close_btn.add_theme_font_size_override("font_size", 14)
	close_btn.pressed.connect(func(): _panel.visible = false)
	header.add_child(close_btn)
	vbox.add_child(header)

	# ?脣????
	_scroll = ScrollContainer.new()
	_scroll.custom_minimum_size = Vector2(0, 500)
	_scroll.horizontal_scroll_mode = ScrollContainer.SCROLL_MODE_DISABLED
	vbox.add_child(_scroll)

	_content = VBoxContainer.new()
	_content.add_theme_constant_override("separation", 4)
	_content.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	_scroll.add_child(_content)

	# 撱箇???憛?
	_mission_section = _build_section("?? 隞隞餃?", Color(0.2, 0.6, 1.0))
	_login_section = _build_section("? ?餃???", Color(1.0, 0.5, 0.1))
	_season_section = _build_section("?? 鞈賢迤??霅?, Color(0.8, 0.3, 1.0))"
	_vip_section = _build_section("?? VIP 蝑?", Color(1.0, 0.85, 0.0))
	_spin_section = _build_section("? 瘥頧", Color(0.3, 1.0, 0.5))
	_tournament_section = _build_section("?? ?行?鞈?, Color(1.0, 0.3, 0.3))"

## _build_section 撱箇?銝?????憛?
func _build_section(title: String, color: Color) -> Control:
	var container := VBoxContainer.new()
	container.add_theme_constant_override("separation", 2)
	_content.add_child(container)

	# ?憛?憿?
	var header := HBoxContainer.new()
	header.custom_minimum_size = Vector2(0, 36)
	var hbg := ColorRect.new()
	hbg.color = Color(color.r * 0.15, color.g * 0.15, color.b * 0.15, 0.9)
	hbg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	header.add_child(hbg)

	var left_bar := ColorRect.new()
	left_bar.custom_minimum_size = Vector2(4, 0)
	left_bar.color = color
	header.add_child(left_bar)

	var title_lbl := Label.new()
	title_lbl.text = "  " + title
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.add_theme_color_override("font_color", color)
	title_lbl.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	title_lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	header.add_child(title_lbl)
	container.add_child(header)

	# ?批捆???
	var body := VBoxContainer.new()
	body.name = "Body"
	body.add_theme_constant_override("separation", 4)
	var body_bg := ColorRect.new()
	body_bg.color = Color(0.05, 0.05, 0.1, 0.7)
	body_bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	body.add_child(body_bg)
	container.add_child(body)

	return container

# ---- 閮??? ----

func _connect_signals() -> void:
	GameManager.mission_updated.connect(_on_mission_updated)
	GameManager.login_progress_received.connect(_on_login_progress)
	GameManager.season_updated.connect(_on_season_updated)
	GameManager.vip_updated.connect(_on_vip_updated)
	GameManager.daily_spin_state.connect(_on_daily_spin_state)
	GameManager.tournament_updated.connect(_on_tournament_updated)
	GameManager.daily_tournament_updated.connect(_on_daily_tournament_updated)

# ---- 鞈??湔 ----

func _on_mission_updated(missions: Array) -> void:
	_mission_data = missions
	if _panel.visible:
		_refresh_mission_section()

func _on_login_progress(data: Dictionary) -> void:
	_login_streak = data.get("current_streak", 0)
	_login_next_milestone_days = data.get("next_milestone_days", 0)
	_login_days_to_next = data.get("days_to_next", 0)
	if _panel.visible:
		_refresh_login_section()

func _on_season_updated(data: Dictionary) -> void:
	_season_level = data.get("current_level", 0)
	_season_points = data.get("season_points", 0)
	_season_points_to_next = data.get("points_to_next", 0)
	_season_progress = data.get("progress", 0.0)
	if _panel.visible:
		_refresh_season_section()

func _on_vip_updated(data: Dictionary) -> void:
	_vip_level = data.get("level", 0)
	_vip_name = data.get("name", "")
	if _panel.visible:
		_refresh_vip_section()

func _on_daily_spin_state(data: Dictionary) -> void:
	_daily_spin_available = not data.get("spun_today", false)
	if _panel.visible:
		_refresh_spin_section()

func _on_tournament_updated(data: Dictionary) -> void:
	_tournament_rank = data.get("player_rank", 0)
	_tournament_points = data.get("player_points", 0)
	if _panel.visible:
		_refresh_tournament_section()

func _on_daily_tournament_updated(data: Dictionary) -> void:
	_daily_tournament_rank = data.get("player_rank", 0)
	_daily_tournament_points = data.get("player_points", 0)
	if _panel.visible:
		_refresh_tournament_section()

# ---- ?Ｘ?? ----

func _on_toggle_pressed() -> void:
	_panel.visible = not _panel.visible
	if _panel.visible:
		_refresh_all()
		# 隢???啗???
		GameManager.request_login_progress()

## _refresh_all ??渡????憛?
func _refresh_all() -> void:
	_refresh_mission_section()
	_refresh_login_section()
	_refresh_season_section()
	_refresh_vip_section()
	_refresh_spin_section()
	_refresh_tournament_section()

# ---- ?憛??----

func _refresh_mission_section() -> void:
	var body := _mission_section.get_node_or_null("Body")
	if not is_instance_valid(body):
		return
	_clear_body(body)

	if _mission_data.is_empty():
		_add_row(body, "  ?怎隞餃?鞈?", Color(0.6, 0.6, 0.6))
		return

	var completed := 0
	var total := _mission_data.size()
	for m in _mission_data:
		if m.get("completed", false):
			completed += 1

	# 摰??脣漲
	var prog_row := HBoxContainer.new()
	prog_row.add_theme_constant_override("separation", 8)
	var prog_lbl := Label.new()
	prog_lbl.text = "  摰? %d/%d" % [completed, total]
	prog_lbl.add_theme_font_size_override("font_size", 13)
	prog_lbl.add_theme_color_override("font_color",
		Color(0.3, 1.0, 0.4) if completed == total else Color(0.8, 0.8, 0.8))
	prog_row.add_child(prog_lbl)
	body.add_child(prog_row)

	# ?遙????
	for m in _mission_data:
		var name: String = m.get("name", "")
		var current: int = m.get("current", 0)
		var target: int = m.get("target", 1)
		var done: bool = m.get("completed", false)
		var claimed: bool = m.get("reward_claimed", false)
		var reward: int = m.get("reward", 0)
		var icon: String = m.get("icon", "??")

		var row := HBoxContainer.new()
		row.add_theme_constant_override("separation", 6)

		var status_lbl := Label.new()
		status_lbl.text = "  " + ("?? if claimed else ("??" if done else icon))"
		status_lbl.add_theme_font_size_override("font_size", 13)
		row.add_child(status_lbl)

		var info_vbox := VBoxContainer.new()
		info_vbox.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		var name_lbl := Label.new()
		name_lbl.text = name
		name_lbl.add_theme_font_size_override("font_size", 12)
		name_lbl.add_theme_color_override("font_color",
			Color(0.5, 0.5, 0.5) if claimed else Color(0.9, 0.9, 0.9))
		info_vbox.add_child(name_lbl)

		if not done:
			var prog_bar := ProgressBar.new()
			prog_bar.custom_minimum_size = Vector2(0, 8)
			prog_bar.min_value = 0
			prog_bar.max_value = target
			prog_bar.value = current
			prog_bar.show_percentage = false
			info_vbox.add_child(prog_bar)
		row.add_child(info_vbox)

		var reward_lbl := Label.new()
		reward_lbl.text = "+%d  " % reward
		reward_lbl.add_theme_font_size_override("font_size", 12)
		reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
		row.add_child(reward_lbl)

		body.add_child(row)

func _refresh_login_section() -> void:
	var body := _login_section.get_node_or_null("Body")
	if not is_instance_valid(body):
		return
	_clear_body(body)

	# ???憭拇
	var streak_row := HBoxContainer.new()
	var streak_lbl := Label.new()
	streak_lbl.text = "  ? ??? %d 憭? % _login_streak"
	streak_lbl.add_theme_font_size_override("font_size", 14)
	streak_lbl.add_theme_color_override("font_color", Color(1.0, 0.6, 0.2))
	streak_row.add_child(streak_lbl)
	body.add_child(streak_row)

	# 銝???蝔?
	if _login_next_milestone_days > 0:
		var next_row := HBoxContainer.new()
		var next_lbl := Label.new()
		next_lbl.text = "  ??頝蝚?%d 憭拚?蝔??榆 %d 憭? % [_login_next_milestone_days, _login_days_to_next]"
		next_lbl.add_theme_font_size_override("font_size", 12)
		next_lbl.add_theme_color_override("font_color", Color(0.7, 0.9, 1.0))
		next_row.add_child(next_lbl)
		body.add_child(next_row)

		# ?脣漲璇?
		var prog := ProgressBar.new()
		prog.custom_minimum_size = Vector2(0, 10)
		prog.min_value = 0
		prog.max_value = _login_next_milestone_days
		prog.value = _login_streak
		prog.show_percentage = false
		var prog_container := HBoxContainer.new()
		prog_container.add_theme_constant_override("separation", 0)
		var spacer := Control.new()
		spacer.custom_minimum_size = Vector2(8, 0)
		prog_container.add_child(spacer)
		prog.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		prog_container.add_child(prog)
		var spacer2 := Control.new()
		spacer2.custom_minimum_size = Vector2(8, 0)
		prog_container.add_child(spacer2)
		body.add_child(prog_container)
	else:
		_add_row(body, "  ?? 撌脤?????蝔?嚗?, Color(1.0, 0.85, 0.0))"

func _refresh_season_section() -> void:
	var body := _season_section.get_node_or_null("Body")
	if not is_instance_valid(body):
		return
	_clear_body(body)

	var level_row := HBoxContainer.new()
	var level_lbl := Label.new()
	level_lbl.text = "  LV%d  %d pts" % [_season_level, _season_points]
	level_lbl.add_theme_font_size_override("font_size", 14)
	level_lbl.add_theme_color_override("font_color", Color(0.9, 0.6, 1.0))
	level_row.add_child(level_lbl)
	body.add_child(level_row)

	if _season_points_to_next > 0:
		var prog := ProgressBar.new()
		prog.custom_minimum_size = Vector2(0, 10)
		prog.min_value = 0
		prog.max_value = 1.0
		prog.value = _season_progress
		prog.show_percentage = false
		var prog_container := HBoxContainer.new()
		var spacer := Control.new()
		spacer.custom_minimum_size = Vector2(8, 0)
		prog_container.add_child(spacer)
		prog.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		prog_container.add_child(prog)
		var spacer2 := Control.new()
		spacer2.custom_minimum_size = Vector2(8, 0)
		prog_container.add_child(spacer2)
		body.add_child(prog_container)

		var pts_lbl := Label.new()
		pts_lbl.text = "  頝銝?蝝?撌?%d pts" % _season_points_to_next
		pts_lbl.add_theme_font_size_override("font_size", 12)
		pts_lbl.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
		body.add_child(pts_lbl)

func _refresh_vip_section() -> void:
	var body := _vip_section.get_node_or_null("Body")
	if not is_instance_valid(body):
		return
	_clear_body(body)

	var vip_row := HBoxContainer.new()
	var vip_lbl := Label.new()
	var vip_text := "  ?? VIP %d" % _vip_level
	if _vip_name != "":
		vip_text += "  (%s)" % _vip_name
	vip_lbl.text = vip_text
	vip_lbl.add_theme_font_size_override("font_size", 14)
	vip_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	vip_row.add_child(vip_lbl)
	body.add_child(vip_row)

func _refresh_spin_section() -> void:
	var body := _spin_section.get_node_or_null("Body")
	if not is_instance_valid(body):
		return
	_clear_body(body)

	var spin_row := HBoxContainer.new()
	var spin_lbl := Label.new()
	if _daily_spin_available:
		spin_lbl.text = "  ??隞頧?舐嚗?????"
		spin_lbl.add_theme_color_override("font_color", Color(0.3, 1.0, 0.4))
	else:
		spin_lbl.text = "  ??隞撌脰?嚗?憭拙?靘?"
		spin_lbl.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	spin_lbl.add_theme_font_size_override("font_size", 13)
	spin_row.add_child(spin_lbl)
	body.add_child(spin_row)

func _refresh_tournament_section() -> void:
	var body := _tournament_section.get_node_or_null("Body")
	if not is_instance_valid(body):
		return
	_clear_body(body)

	# ?梯魚
	var week_row := HBoxContainer.new()
	var week_lbl := Label.new()
	if _tournament_rank > 0:
		week_lbl.text = "  ?? ?梯魚 #%d  %d pts" % [_tournament_rank, _tournament_points]
		week_lbl.add_theme_color_override("font_color",
			Color(1.0, 0.85, 0.0) if _tournament_rank <= 3 else Color(0.8, 0.8, 0.8))
	else:
		week_lbl.text = "  ?? ?梯魚 ?芯?璁?"
		week_lbl.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	week_lbl.add_theme_font_size_override("font_size", 13)
	week_row.add_child(week_lbl)
	body.add_child(week_row)

	# ?亥魚
	var day_row := HBoxContainer.new()
	var day_lbl := Label.new()
	if _daily_tournament_rank > 0:
		day_lbl.text = "  ?? ?亥魚 #%d  %d pts" % [_daily_tournament_rank, _daily_tournament_points]
		day_lbl.add_theme_color_override("font_color",
			Color(1.0, 0.85, 0.0) if _daily_tournament_rank <= 3 else Color(0.8, 0.8, 0.8))
	else:
		day_lbl.text = "  ?? ?亥魚 ?芯?璁?"
		day_lbl.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	day_lbl.add_theme_font_size_override("font_size", 13)
	day_row.add_child(day_lbl)
	body.add_child(day_row)

# ---- 撌亙?賣 ----

func _clear_body(body: Control) -> void:
	for child in body.get_children():
		if child is ColorRect:
			continue  # 靽??
		child.queue_free()

func _add_row(parent: Control, text: String, color: Color) -> void:
	var row := HBoxContainer.new()
	var lbl := Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 13)
	lbl.add_theme_color_override("font_color", color)
	row.add_child(lbl)
	parent.add_child(row)
