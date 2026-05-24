## MissionPanel.gd
## 瘥隞餃??Ｘ嚗AY-037嚗? HUD.gd ?? DAY-053嚗?
## 憿舐內隞隞餃??脣漲嚗?湧?????

extends Control

# ??HUD.gd ?典遣蝡?閮剖?
var pixel_font: Font = null

# 閮?嚗遙???嚗蝯?HUD ??撠梢蝟餌絞嚗?
signal mission_completed_notify(mission_data: Dictionary)

var _mission_data: Array = []
var _mission_visible: bool = false
var _mission_reset_at_ms: int = 0

## ??????HUD.gd ?澆嚗?
func setup(font: Font) -> void:
	pixel_font = font
	GameManager.mission_updated.connect(_on_mission_updated)
	GameManager.mission_completed.connect(_on_mission_completed)

## 撱箇?隞餃???嚗opBar ?喳嚗????遙??選?
func create_button(top_bar: Control) -> void:
	if not is_instance_valid(top_bar):
		return
	var btn = Button.new()
	btn.name = "MissionButton"
	btn.text = "?? 隞餃?"
	btn.position = Vector2(750, 4)
	btn.size = Vector2(80, 32)
	btn.add_theme_font_size_override("font_size", 12)
	if is_instance_valid(pixel_font):
		btn.add_theme_font_override("font", pixel_font)
	btn.pressed.connect(_toggle_panel)
	top_bar.add_child(btn)

## ??隞餃??Ｘ憿舐內
func _toggle_panel() -> void:
	_mission_visible = not _mission_visible
	visible = _mission_visible
	if _mission_visible:
		NetworkManager.send("get_missions", {})
		_refresh_mission_list()

## 隞餃??脣漲?湔
func _on_mission_updated(missions: Array) -> void:
	_mission_data = missions
	if _mission_visible:
		_refresh_mission_list()

## ?瑟隞餃??” UI
func _refresh_mission_list() -> void:
	var list = get_node_or_null("MissionList")
	if not is_instance_valid(list):
		return

	# 皜?摰?
	for child in list.get_children():
		child.queue_free()

	# 撱箇?隞餃?璇
	for i in range(_mission_data.size()):
		var m = _mission_data[i]
		_create_mission_row(list, m, i)

	# ?湔?蔭?
	_update_mission_reset_countdown()

## 撱箇??桐?隞餃?璇
func _create_mission_row(container: Control, mission: Dictionary, index: int) -> void:
	var row = Control.new()
	row.position = Vector2(0, index * 52)
	row.size = Vector2(380, 50)
	container.add_child(row)

	var completed = mission.get("completed", false)
	var reward_claimed = mission.get("reward_claimed", false)
	var current = mission.get("current", 0)
	var target = mission.get("target", 1)
	var reward = mission.get("reward", 0)
	var mission_type = mission.get("type", "")
	var is_combo = (mission_type == "combo")

	# ?
	var bg = ColorRect.new()
	bg.size = Vector2(376, 48)
	bg.position = Vector2(2, 1)
	if completed and reward_claimed:
		bg.color = Color(0.05, 0.15, 0.05, 0.7)
	elif completed:
		bg.color = Color(0.05, 0.2, 0.05, 0.85)
	elif is_combo:
		bg.color = Color(0.18, 0.06, 0.02, 0.85)
	else:
		bg.color = Color(0.03, 0.06, 0.18, 0.7)
	row.add_child(bg)

	# combo 隞餃?嚗椰?湔?蝝?璇?
	if is_combo and not completed:
		var side_bar = ColorRect.new()
		side_bar.size = Vector2(3, 46)
		side_bar.position = Vector2(2, 1)
		side_bar.color = Color(1.0, 0.45, 0.1, 0.9)
		row.add_child(side_bar)

	# ?內
	var icon_lbl = Label.new()
	icon_lbl.text = mission.get("icon", "??")
	icon_lbl.position = Vector2(8, 12)
	icon_lbl.add_theme_font_size_override("font_size", 20)
	row.add_child(icon_lbl)

	# combo 隞餃?嚗???內???
	if is_combo and not completed:
		var pulse_tween = row.create_tween().set_loops()
		pulse_tween.tween_property(icon_lbl, "scale", Vector2(1.3, 1.3), 0.4).set_trans(Tween.TRANS_SINE)
		pulse_tween.tween_property(icon_lbl, "scale", Vector2(1.0, 1.0), 0.4).set_trans(Tween.TRANS_SINE)
		var color_tween = row.create_tween().set_loops()
		color_tween.tween_property(icon_lbl, "modulate", Color(1.0, 0.8, 0.2), 0.4)
		color_tween.tween_property(icon_lbl, "modulate", Color(1.0, 0.4, 0.1), 0.4)

	# 隞餃??迂
	var name_lbl = Label.new()
	name_lbl.text = mission.get("name", "")
	name_lbl.position = Vector2(40, 4)
	name_lbl.size = Vector2(200, 20)
	name_lbl.add_theme_font_size_override("font_size", 13)
	if completed:
		name_lbl.modulate = Color(0.5, 1.0, 0.5)
	elif is_combo:
		name_lbl.modulate = Color(1.0, 0.75, 0.3)
	else:
		name_lbl.modulate = Color.WHITE
	if is_instance_valid(pixel_font):
		name_lbl.add_theme_font_override("font", pixel_font)
	row.add_child(name_lbl)

	# ?脣漲??
	var progress_lbl = Label.new()
	progress_lbl.text = "%d / %d" % [current, target]
	progress_lbl.position = Vector2(40, 26)
	progress_lbl.size = Vector2(120, 16)
	progress_lbl.add_theme_font_size_override("font_size", 11)
	progress_lbl.modulate = Color(0.7, 0.9, 0.7) if completed else Color(0.7, 0.7, 0.7)
	if is_instance_valid(pixel_font):
		progress_lbl.add_theme_font_override("font", pixel_font)
	row.add_child(progress_lbl)

	# ?脣漲璇???
	var bar_bg = ColorRect.new()
	bar_bg.size = Vector2(160, 6)
	bar_bg.position = Vector2(40, 42)
	bar_bg.color = Color(0.1, 0.1, 0.1, 0.8)
	row.add_child(bar_bg)

	# ?脣漲璇‵??
	var fill_ratio = float(current) / float(max(target, 1))
	var bar_fill = ColorRect.new()
	bar_fill.size = Vector2(160.0 * fill_ratio, 6)
	bar_fill.position = Vector2(40, 42)
	if completed:
		bar_fill.color = Color(0.3, 1.0, 0.4)
	elif is_combo:
		bar_fill.color = Color(1.0, 0.45, 0.1)
	else:
		bar_fill.color = Color(0.2, 0.6, 1.0)
	row.add_child(bar_fill)

	# ???
	var reward_lbl = Label.new()
	reward_lbl.text = "??%d" % reward
	reward_lbl.position = Vector2(248, 4)
	reward_lbl.size = Vector2(80, 20)
	reward_lbl.add_theme_font_size_override("font_size", 12)
	reward_lbl.modulate = Color(1.0, 0.9, 0.3)
	if is_instance_valid(pixel_font):
		reward_lbl.add_theme_font_override("font", pixel_font)
	row.add_child(reward_lbl)

	# ????嚗????芷???憿舐內嚗?
	if completed and not reward_claimed:
		var claim_btn = Button.new()
		claim_btn.text = "??"
		claim_btn.position = Vector2(300, 12)
		claim_btn.size = Vector2(68, 28)
		claim_btn.add_theme_font_size_override("font_size", 12)
		claim_btn.modulate = Color(0.3, 1.0, 0.4)
		if is_instance_valid(pixel_font):
			claim_btn.add_theme_font_override("font", pixel_font)
		var mission_id = mission.get("id", "")
		claim_btn.pressed.connect(func():
			NetworkManager.send("claim_mission", {"mission_id": mission_id})
			AudioManager.play_sfx(AudioManager.SFX.COIN_DROP)
		)
		row.add_child(claim_btn)
	elif reward_claimed:
		var done_lbl = Label.new()
		done_lbl.text = "??撌脤???"
		done_lbl.position = Vector2(296, 16)
		done_lbl.size = Vector2(76, 20)
		done_lbl.add_theme_font_size_override("font_size", 11)
		done_lbl.modulate = Color(0.5, 0.8, 0.5)
		if is_instance_valid(pixel_font):
			done_lbl.add_theme_font_override("font", pixel_font)
		row.add_child(done_lbl)

## 隞餃?摰??
func _on_mission_completed(mission_data: Dictionary) -> void:
	# ?瑟隞餃??Ｘ
	if _mission_visible:
		NetworkManager.send("get_missions", {})
	# ? HUD 憿舐內?停?
	mission_completed_notify.emit(mission_data)

## 閮剖?隞餃??蔭??嚗 GameManager ?澆嚗?
func set_mission_reset_at(reset_at_ms: int) -> void:
	_mission_reset_at_ms = reset_at_ms
	_update_mission_reset_countdown()

## ?湔隞餃??蔭?憿舐內
func _update_mission_reset_countdown() -> void:
	var countdown_lbl = get_node_or_null("ResetCountdown")
	if not is_instance_valid(countdown_lbl):
		countdown_lbl = Label.new()
		countdown_lbl.name = "ResetCountdown"
		countdown_lbl.position = Vector2(8, 272)
		countdown_lbl.size = Vector2(364, 20)
		countdown_lbl.add_theme_font_size_override("font_size", 11)
		countdown_lbl.modulate = Color(0.6, 0.7, 0.6)
		if is_instance_valid(pixel_font):
			countdown_lbl.add_theme_font_override("font", pixel_font)
		add_child(countdown_lbl)

	if _mission_reset_at_ms > 0:
		var now_ms = int(Time.get_unix_time_from_system() * 1000)
		var diff_sec = int((_mission_reset_at_ms - now_ms) / 1000)
		if diff_sec > 0:
			var hours = diff_sec / 3600
			var mins = (diff_sec % 3600) / 60
			countdown_lbl.text = "?? ?蔭?嚗?dh %02dm嚗TC+8 00:00嚗? % [hours, mins]"
		else:
			countdown_lbl.text = "?? 隞餃??喳??蔭..."
	else:
		countdown_lbl.text = "?? ?蔭??嚗???00:00嚗TC+8嚗?"

## 撱箇??Ｘ UI嚗 HUD.gd ?澆 setup 敺?銵?
func _ready() -> void:
	name = "MissionPanel"
	position = Vector2(640, 50)
	size = Vector2(380, 300)
	z_index = 80
	visible = false
	_build_panel_ui()

func _build_panel_ui() -> void:
	# ?
	var bg = ColorRect.new()
	bg.size = Vector2(380, 300)
	bg.color = Color(0.02, 0.05, 0.15, 0.92)
	add_child(bg)

	# ???
	var top_line = ColorRect.new()
	top_line.size = Vector2(380, 3)
	top_line.color = Color(0.9, 0.75, 0.2, 0.8)
	add_child(top_line)

	# 璅?
	var title = Label.new()
	title.name = "MissionTitle"
	title.text = "?? 隞隞餃?"
	title.position = Vector2(12, 8)
	title.add_theme_font_size_override("font_size", 16)
	title.modulate = Color(1.0, 0.9, 0.3)
	if is_instance_valid(pixel_font):
		title.add_theme_font_override("font", pixel_font)
	add_child(title)

	# ????
	var close_btn = Button.new()
	close_btn.text = "??"
	close_btn.position = Vector2(348, 4)
	close_btn.size = Vector2(28, 24)
	close_btn.add_theme_font_size_override("font_size", 12)
	close_btn.pressed.connect(func():
		_mission_visible = false
		visible = false
	)
	add_child(close_btn)

	# 隞餃??”摰孵
	var list = Control.new()
	list.name = "MissionList"
	list.position = Vector2(0, 36)
	list.size = Vector2(380, 264)
	add_child(list)

	# ??憿舐內???乩葉...??
	var loading = Label.new()
	loading.name = "LoadingLabel"
	loading.text = "頛隞餃?銝?.."
	loading.position = Vector2(120, 100)
	loading.add_theme_font_size_override("font_size", 14)
	loading.modulate = Color(0.6, 0.6, 0.6)
	if is_instance_valid(pixel_font):
		loading.add_theme_font_override("font", pixel_font)
	list.add_child(loading)
