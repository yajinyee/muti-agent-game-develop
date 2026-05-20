## LeaderboardPanel.gd
## 排行榜面板獨立腳本（DAY-058）
## 從 HUD.gd 拆分：_create_leaderboard_panel, _create_leaderboard_row,
##                  _show_leaderboard_placeholder, _on_leaderboard_updated, _toggle_leaderboard
extends RefCounted

const MAX_LEADERBOARD_ENTRIES = 5

var _panel: Control = null
var _visible: bool = true
var _toggle_btn: Button = null
var _pixel_font: Font = null

## 初始化排行榜面板
## parent: 要加入的父節點（HUD CanvasLayer）
## font: 像素字體（可為 null）
func setup(parent: Node, font: Font) -> void:
	_pixel_font = font
	_create_panel(parent)

## 取得面板節點（供 HUD 存取）
func get_panel() -> Control:
	return _panel

## 更新排行榜資料
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

		# 名次
		var rank_lbl = row.get_node_or_null("RankLabel")
		if rank_lbl:
			match i:
				0: rank_lbl.text = "🥇"
				1: rank_lbl.text = "🥈"
				2: rank_lbl.text = "🥉"
				_: rank_lbl.text = "#%d" % (i + 1)

		# 玩家名稱（自己高亮）
		var name_lbl = row.get_node_or_null("NameLabel")
		if name_lbl:
			var display = entry.get("display_name", "???")
			name_lbl.text = ("▶" if is_self else "") + display
			name_lbl.modulate = Color(1.0, 1.0, 0.4) if is_self else Color.WHITE

		# 分數
		var score_lbl = row.get_node_or_null("ScoreLabel")
		if score_lbl:
			var score = entry.get("score", 0)
			score_lbl.text = "💰%d" % score
			score_lbl.modulate = Color(1.0, 0.9, 0.3) if is_self else Color(0.9, 0.9, 0.9)

		# 擊殺數
		var kill_lbl = row.get_node_or_null("KillLabel")
		if kill_lbl:
			kill_lbl.text = "⚔%d" % entry.get("kill_count", 0)

		# 稱號（DAY-068）
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

		# 自己的行背景高亮
		var row_bg = row.get_node_or_null("RowBG")
		if row_bg:
			if is_self:
				row_bg.color = Color(0.15, 0.25, 0.05, 0.8)
			elif i % 2 == 0:
				row_bg.color = Color(0.05, 0.1, 0.25, 0.6)
			else:
				row_bg.color = Color(0.03, 0.07, 0.18, 0.6)

	# 動態調整高度
	var new_height = 30 + count * 40
	if is_instance_valid(_panel):
		var bg = _panel.get_node_or_null("LeaderboardBG")
		if bg:
			bg.size.y = new_height
		_panel.size.y = new_height

## 顯示等待佔位符
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
			name_lbl.text = "等待玩家..."
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

## 切換顯示/隱藏
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
		_toggle_btn.text = "▲" if _visible else "▼"

# ---- 私有方法 ----

func _create_panel(parent: Node) -> void:
	# 位置：BOSS 計時器 x=900, y=50，高度 80px，所以排行榜從 y=140 開始
	var panel = Control.new()
	panel.name = "LeaderboardPanel"
	panel.position = Vector2(900, 140)
	panel.size = Vector2(360, 200)
	panel.z_index = 10
	parent.add_child(panel)
	_panel = panel

	# 背景
	var bg = ColorRect.new()
	bg.name = "LeaderboardBG"
	bg.size = Vector2(360, 200)
	bg.color = Color(0.0, 0.05, 0.15, 0.82)
	panel.add_child(bg)

	# 標題列
	var title_bar = ColorRect.new()
	title_bar.size = Vector2(360, 28)
	title_bar.color = Color(0.05, 0.15, 0.4, 0.95)
	panel.add_child(title_bar)

	var title_lbl = Label.new()
	title_lbl.name = "LeaderboardTitle"
	title_lbl.text = "🏆 排行榜"
	title_lbl.position = Vector2(10, 4)
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.modulate = Color(1.0, 0.9, 0.3)
	if is_instance_valid(_pixel_font):
		title_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(title_lbl)

	# 折疊按鈕
	var toggle_btn = Button.new()
	toggle_btn.name = "LeaderboardToggle"
	toggle_btn.text = "▲"
	toggle_btn.position = Vector2(325, 2)
	toggle_btn.size = Vector2(30, 24)
	toggle_btn.add_theme_font_size_override("font_size", 12)
	toggle_btn.pressed.connect(toggle)
	panel.add_child(toggle_btn)
	_toggle_btn = toggle_btn

	# 排行榜條目容器
	var entries_container = Control.new()
	entries_container.name = "EntriesContainer"
	entries_container.position = Vector2(0, 30)
	entries_container.size = Vector2(360, 170)
	panel.add_child(entries_container)

	# 預建立 5 個排行榜行（避免動態建立）
	for i in range(MAX_LEADERBOARD_ENTRIES):
		_create_row(entries_container, i)

	# 顯示等待佔位符
	show_placeholder()

func _create_row(container: Control, index: int) -> void:
	var row = Control.new()
	row.name = "Row%d" % index
	row.position = Vector2(0, index * 40)
	row.size = Vector2(360, 38)
	container.add_child(row)

	# 行背景
	var row_bg = ColorRect.new()
	row_bg.name = "RowBG"
	row_bg.size = Vector2(360, 38)
	if index % 2 == 0:
		row_bg.color = Color(0.05, 0.1, 0.25, 0.6)
	else:
		row_bg.color = Color(0.03, 0.07, 0.18, 0.6)
	row.add_child(row_bg)

	# 名次標籤
	var rank_lbl = Label.new()
	rank_lbl.name = "RankLabel"
	rank_lbl.position = Vector2(6, 6)
	rank_lbl.size = Vector2(30, 20)
	rank_lbl.add_theme_font_size_override("font_size", 13)
	if is_instance_valid(_pixel_font):
		rank_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(rank_lbl)

	# 玩家名稱
	var name_lbl = Label.new()
	name_lbl.name = "NameLabel"
	name_lbl.position = Vector2(42, 6)
	name_lbl.size = Vector2(140, 20)
	name_lbl.add_theme_font_size_override("font_size", 12)
	name_lbl.clip_text = true
	if is_instance_valid(_pixel_font):
		name_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(name_lbl)

	# 分數
	var score_lbl = Label.new()
	score_lbl.name = "ScoreLabel"
	score_lbl.position = Vector2(188, 6)
	score_lbl.size = Vector2(100, 20)
	score_lbl.add_theme_font_size_override("font_size", 12)
	score_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	if is_instance_valid(_pixel_font):
		score_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(score_lbl)

	# 擊殺數
	var kill_lbl = Label.new()
	kill_lbl.name = "KillLabel"
	kill_lbl.position = Vector2(295, 6)
	kill_lbl.size = Vector2(60, 20)
	kill_lbl.add_theme_font_size_override("font_size", 11)
	kill_lbl.modulate = Color(0.7, 0.9, 0.7)
	if is_instance_valid(_pixel_font):
		kill_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(kill_lbl)

	# 稱號標籤（DAY-068）— 顯示在名稱下方
	var title_lbl = Label.new()
	title_lbl.name = "TitleLabel"
	title_lbl.position = Vector2(42, 18)
	title_lbl.size = Vector2(200, 14)
	title_lbl.add_theme_font_size_override("font_size", 10)
	title_lbl.modulate = Color(0.7, 0.7, 0.7, 0.9)
	if is_instance_valid(_pixel_font):
		title_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(title_lbl)

	row.visible = false
