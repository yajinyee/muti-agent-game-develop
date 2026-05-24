## LeaderboardPanel.gd
## ??璁?輻蝡?穿?DAY-058嚗?
## 敺?HUD.gd ??嚗create_leaderboard_panel, _create_leaderboard_row,
##                  _show_leaderboard_placeholder, _on_leaderboard_updated, _toggle_leaderboard
extends RefCounted

const MAX_LEADERBOARD_ENTRIES = 5

var _panel: Control = null
var _visible: bool = true
var _toggle_btn: Button = null
var _pixel_font: Font = null

## ????銵??Ｘ
## parent: 閬??亦??嗥?暺?HUD CanvasLayer嚗?
## font: ??摮?嚗??null嚗?
func setup(parent: Node, font: Font) -> void:
	_pixel_font = font
	_create_panel(parent)

## ???Ｘ蝭暺?靘?HUD 摮?嚗?
func get_panel() -> Control:
	return _panel

## ?湔??璁???
func update(entries: Array, my_player_id: String) -> void:
	if not is_instance_valid(_panel):
		return
	var container = _panel.get_node_or_null("EntriesContainer")
	if not container:
		return

	var count = min(entries.size(), MAX_LEADERBOARD_ENTRIES)

	for i in range(MAX_LEADERBOARD_ENTRIES):
		var row = container.get_node_or_null("Row%d" % i)
		if not row:
			continue

		if i >= count:
			row.visible = false
			continue

		row.visible = true
		var entry = entries[i]
		var is_self = entry.get("player_id", "") == my_player_id

		# ?活
		var rank_lbl = row.get_node_or_null("RankLabel")
		if rank_lbl:
			match i:
				0: rank_lbl.text = "??"
				1: rank_lbl.text = "??"
				2: rank_lbl.text = "??"
				_: rank_lbl.text = "#%d" % (i + 1)

		# ?拙振?迂嚗撌梢?鈭殷?
		var name_lbl = row.get_node_or_null("NameLabel")
		if name_lbl:
			var display = entry.get("display_name", "???")
			name_lbl.text = ("?? if is_self else "") + display"
			name_lbl.modulate = Color(1.0, 1.0, 0.4) if is_self else Color.WHITE

		# ?
		var score_lbl = row.get_node_or_null("ScoreLabel")
		if score_lbl:
			var score = entry.get("score", 0)
			score_lbl.text = "?%d" % score
			score_lbl.modulate = Color(1.0, 0.9, 0.3) if is_self else Color(0.9, 0.9, 0.9)

		# ?捏??
		var kill_lbl = row.get_node_or_null("KillLabel")
		if kill_lbl:
			kill_lbl.text = "??d" % entry.get("kill_count", 0)

		# 蝔梯?嚗AY-068嚗?
		var title_lbl = row.get_node_or_null("TitleLabel")
		if title_lbl:
			var title_icon: String = entry.get("title_icon", "")
			var title_name: String = entry.get("title_name", "")
			var title_color: String = entry.get("title_color", "#AAAAAA")
			if title_icon != "" and title_name != "":
				title_lbl.text = title_icon + " " + title_name
				title_lbl.add_theme_color_override("font_color", Color.html(title_color))
			else:
				title_lbl.text = ""

		# ?末????DAY-073嚗??撌梁??拙振?＊蝷?
		var add_friend_btn = row.get_node_or_null("AddFriendBtn")
		if add_friend_btn:
			var entry_player_id: String = entry.get("player_id", "")
			if not is_self and entry_player_id != "":
				add_friend_btn.visible = true
				# 皜????嚗??圈?
				if add_friend_btn.pressed.get_connections().size() > 0:
					for conn in add_friend_btn.pressed.get_connections():
						add_friend_btn.pressed.disconnect(conn["callable"])
				add_friend_btn.pressed.connect(func():
					NetworkManager.send_message({
						"type": "send_friend_request",
						"payload": {"target_id": entry_player_id}
					})
					add_friend_btn.text = "??"
					add_friend_btn.disabled = true
				)
			else:
				add_friend_btn.visible = false

		# ?芸楛???擃漁
		var row_bg = row.get_node_or_null("RowBG")
		if row_bg:
			if is_self:
				row_bg.color = Color(0.15, 0.25, 0.05, 0.8)
			elif i % 2 == 0:
				row_bg.color = Color(0.05, 0.1, 0.25, 0.6)
			else:
				row_bg.color = Color(0.03, 0.07, 0.18, 0.6)

	# ??隤踵擃漲
	var new_height = 30 + count * 40
	if is_instance_valid(_panel):
		var bg = _panel.get_node_or_null("LeaderboardBG")
		if bg:
			bg.size.y = new_height
		_panel.size.y = new_height

## 憿舐內蝑?雿?蝚?
func show_placeholder() -> void:
	if not is_instance_valid(_panel):
		return
	var container = _panel.get_node_or_null("EntriesContainer")
	if not container:
		return

	var row = container.get_node_or_null("Row0")
	if row:
		row.visible = true
		var name_lbl = row.get_node_or_null("NameLabel")
		if name_lbl:
			name_lbl.text = "蝑??拙振..."
			name_lbl.modulate = Color(0.6, 0.6, 0.6)
		var rank_lbl = row.get_node_or_null("RankLabel")
		if rank_lbl:
			rank_lbl.text = ""
		var score_lbl = row.get_node_or_null("ScoreLabel")
		if score_lbl:
			score_lbl.text = ""
		var kill_lbl = row.get_node_or_null("KillLabel")
		if kill_lbl:
			kill_lbl.text = ""

## ??憿舐內/?梯?
func toggle() -> void:
	if not is_instance_valid(_panel):
		return
	_visible = not _visible

	var container = _panel.get_node_or_null("EntriesContainer")
	if container:
		container.visible = _visible

	var bg = _panel.get_node_or_null("LeaderboardBG")
	if bg:
		bg.size.y = 230 if _visible else 28

	if is_instance_valid(_toggle_btn):
		_toggle_btn.text = "?? if _visible else "??

# ---- 蝘??寞? ----

func _create_panel(parent: Node) -> void:
	# 雿蔭嚗OSS 閮???x=900, y=50嚗?摨?80px嚗?隞交?銵?敺?y=140 ??
	var panel = Control.new()
	panel.name = "LeaderboardPanel"
	panel.position = Vector2(900, 140)
	panel.size = Vector2(360, 200)
	panel.z_index = 10
	parent.add_child(panel)
	_panel = panel

	# ?
	var bg = ColorRect.new()
	bg.name = "LeaderboardBG"
	bg.size = Vector2(360, 200)
	bg.color = Color(0.0, 0.05, 0.15, 0.82)
	panel.add_child(bg)

	# 璅???
	var title_bar = ColorRect.new()
	title_bar.size = Vector2(360, 28)
	title_bar.color = Color(0.05, 0.15, 0.4, 0.95)
	panel.add_child(title_bar)

	var title_lbl = Label.new()
	title_lbl.name = "LeaderboardTitle"
	title_lbl.text = "?? ??璁?"
	title_lbl.position = Vector2(10, 4)
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.modulate = Color(1.0, 0.9, 0.3)
	if is_instance_valid(_pixel_font):
		title_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(title_lbl)

	# ????
	var toggle_btn = Button.new()
	toggle_btn.name = "LeaderboardToggle"
	toggle_btn.text = "??"
	toggle_btn.position = Vector2(325, 2)
	toggle_btn.size = Vector2(30, 24)
	toggle_btn.add_theme_font_size_override("font_size", 12)
	toggle_btn.pressed.connect(toggle)
	panel.add_child(toggle_btn)
	_toggle_btn = toggle_btn

	# ??璁??桀捆??
	var entries_container = Control.new()
	entries_container.name = "EntriesContainer"
	entries_container.position = Vector2(0, 30)
	entries_container.size = Vector2(360, 170)
	panel.add_child(entries_container)

	# ?遣蝡?5 ??銵?銵??踹???撱箇?嚗?
	for i in range(MAX_LEADERBOARD_ENTRIES):
		_create_row(entries_container, i)

	# 憿舐內蝑?雿?蝚?
	show_placeholder()

func _create_row(container: Control, index: int) -> void:
	var row = Control.new()
	row.name = "Row%d" % index
	row.position = Vector2(0, index * 40)
	row.size = Vector2(360, 38)
	container.add_child(row)

	# 銵???
	var row_bg = ColorRect.new()
	row_bg.name = "RowBG"
	row_bg.size = Vector2(360, 38)
	if index % 2 == 0:
		row_bg.color = Color(0.05, 0.1, 0.25, 0.6)
	else:
		row_bg.color = Color(0.03, 0.07, 0.18, 0.6)
	row.add_child(row_bg)

	# ?活璅惜
	var rank_lbl = Label.new()
	rank_lbl.name = "RankLabel"
	rank_lbl.position = Vector2(6, 6)
	rank_lbl.size = Vector2(30, 20)
	rank_lbl.add_theme_font_size_override("font_size", 13)
	if is_instance_valid(_pixel_font):
		rank_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(rank_lbl)

	# ?拙振?迂
	var name_lbl = Label.new()
	name_lbl.name = "NameLabel"
	name_lbl.position = Vector2(42, 6)
	name_lbl.size = Vector2(140, 20)
	name_lbl.add_theme_font_size_override("font_size", 12)
	name_lbl.clip_text = true
	if is_instance_valid(_pixel_font):
		name_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(name_lbl)

	# ?
	var score_lbl = Label.new()
	score_lbl.name = "ScoreLabel"
	score_lbl.position = Vector2(188, 6)
	score_lbl.size = Vector2(100, 20)
	score_lbl.add_theme_font_size_override("font_size", 12)
	score_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	if is_instance_valid(_pixel_font):
		score_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(score_lbl)

	# ?捏??
	var kill_lbl = Label.new()
	kill_lbl.name = "KillLabel"
	kill_lbl.position = Vector2(295, 6)
	kill_lbl.size = Vector2(60, 20)
	kill_lbl.add_theme_font_size_override("font_size", 11)
	kill_lbl.modulate = Color(0.7, 0.9, 0.7)
	if is_instance_valid(_pixel_font):
		kill_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(kill_lbl)

	# 蝔梯?璅惜嚗AY-068嚗?憿舐內?典?蝔曹???
	var title_lbl = Label.new()
	title_lbl.name = "TitleLabel"
	title_lbl.position = Vector2(42, 18)
	title_lbl.size = Vector2(200, 14)
	title_lbl.add_theme_font_size_override("font_size", 10)
	title_lbl.modulate = Color(0.7, 0.7, 0.7, 0.9)
	if is_instance_valid(_pixel_font):
		title_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(title_lbl)

	# ?末????DAY-073嚗?
	var add_friend_btn = Button.new()
	add_friend_btn.name = "AddFriendBtn"
	add_friend_btn.text = "??"
	add_friend_btn.position = Vector2(330, 8)
	add_friend_btn.size = Vector2(24, 22)
	add_friend_btn.flat = true
	add_friend_btn.tooltip_text = "?末??"
	add_friend_btn.visible = false  # ?身?梯?嚗pdate() ??璇辣憿舐內
	if is_instance_valid(_pixel_font):
		add_friend_btn.add_theme_font_override("font", _pixel_font)
		add_friend_btn.add_theme_font_size_override("font_size", 11)
	row.add_child(add_friend_btn)

	row.visible = false
