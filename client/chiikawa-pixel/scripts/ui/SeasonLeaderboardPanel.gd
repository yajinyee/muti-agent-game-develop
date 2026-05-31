## SeasonLeaderboardPanel.gd — 賽季排行榜（DAY-348）
## 顯示本賽季 XP 前 20 名，以及玩家自己的排名
extends CanvasLayer

# ── 節點引用 ──────────────────────────────────────────────────
var _panel: PanelContainer
var _title_label: Label
var _season_label: Label
var _list_container: VBoxContainer
var _my_rank_label: Label
var _close_btn: Button
var _refresh_btn: Button

# ── 常數 ──────────────────────────────────────────────────────
const PANEL_WIDTH = 500
const PANEL_HEIGHT = 560
const ROW_HEIGHT = 44

# 排名顏色
const RANK_COLORS = {
	1: Color(1.0, 0.85, 0.1),   # 金
	2: Color(0.8, 0.8, 0.85),   # 銀
	3: Color(0.8, 0.5, 0.2),    # 銅
}

func _ready() -> void:
	layer = 23  # 在賽季通行證上方
	_build_ui()
	_connect_signals()
	visible = false

func _build_ui() -> void:
	# 半透明背景遮罩
	var overlay = ColorRect.new()
	overlay.color = Color(0, 0, 0, 0.65)
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	
	# 主面板
	_panel = PanelContainer.new()
	_panel.set_anchors_preset(Control.PRESET_CENTER)
	_panel.custom_minimum_size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_panel.position = Vector2(-PANEL_WIDTH / 2, -PANEL_HEIGHT / 2)
	add_child(_panel)
	
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.08, 0.18, 0.97)
	style.border_color = Color(0.6, 0.5, 0.9, 1.0)
	style.border_width_left = 2
	style.border_width_right = 2
	style.border_width_top = 2
	style.border_width_bottom = 2
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	_panel.add_theme_stylebox_override("panel", style)
	
	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 8)
	_panel.add_child(vbox)
	
	# ── 標題列 ──────────────────────────────────────────────────
	var header = HBoxContainer.new()
	vbox.add_child(header)
	
	_title_label = Label.new()
	_title_label.text = "🏆 賽季排行榜"
	_title_label.add_theme_font_size_override("font_size", 20)
	_title_label.add_theme_color_override("font_color", Color(0.8, 0.6, 1.0))
	_title_label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	header.add_child(_title_label)
	
	_refresh_btn = Button.new()
	_refresh_btn.text = "🔄"
	_refresh_btn.custom_minimum_size = Vector2(36, 36)
	header.add_child(_refresh_btn)
	
	_close_btn = Button.new()
	_close_btn.text = "✕"
	_close_btn.custom_minimum_size = Vector2(36, 36)
	_close_btn.add_theme_font_size_override("font_size", 18)
	header.add_child(_close_btn)
	
	# ── 賽季資訊 ────────────────────────────────────────────────
	_season_label = Label.new()
	_season_label.text = "賽季：載入中..."
	_season_label.add_theme_font_size_override("font_size", 13)
	_season_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.8))
	vbox.add_child(_season_label)
	
	# 分隔線
	vbox.add_child(HSeparator.new())
	
	# ── 表頭 ────────────────────────────────────────────────────
	var header_row = _create_header_row()
	vbox.add_child(header_row)
	
	# ── 排行榜列表（可滾動）────────────────────────────────────
	var scroll = ScrollContainer.new()
	scroll.custom_minimum_size = Vector2(PANEL_WIDTH - 20, PANEL_HEIGHT - 220)
	scroll.size_flags_vertical = Control.SIZE_EXPAND_FILL
	vbox.add_child(scroll)
	
	_list_container = VBoxContainer.new()
	_list_container.add_theme_constant_override("separation", 3)
	_list_container.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	scroll.add_child(_list_container)
	
	# 分隔線
	vbox.add_child(HSeparator.new())
	
	# ── 我的排名 ────────────────────────────────────────────────
	_my_rank_label = Label.new()
	_my_rank_label.text = "我的排名：載入中..."
	_my_rank_label.add_theme_font_size_override("font_size", 14)
	_my_rank_label.add_theme_color_override("font_color", Color(0.9, 0.9, 0.5))
	_my_rank_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(_my_rank_label)

func _create_header_row() -> Control:
	var row = HBoxContainer.new()
	row.add_theme_constant_override("separation", 0)
	
	var labels = [
		{"text": "排名", "width": 50, "align": HORIZONTAL_ALIGNMENT_CENTER},
		{"text": "玩家", "width": 200, "align": HORIZONTAL_ALIGNMENT_LEFT},
		{"text": "等級", "width": 100, "align": HORIZONTAL_ALIGNMENT_CENTER},
		{"text": "賽季 XP", "width": 120, "align": HORIZONTAL_ALIGNMENT_RIGHT},
	]
	
	for l in labels:
		var label = Label.new()
		label.text = l["text"]
		label.custom_minimum_size = Vector2(l["width"], 0)
		label.add_theme_font_size_override("font_size", 12)
		label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.8))
		label.horizontal_alignment = l["align"]
		row.add_child(label)
	
	return row

func _connect_signals() -> void:
	_close_btn.pressed.connect(_on_close)
	_refresh_btn.pressed.connect(_on_refresh)
	GameManager.season_leaderboard_received.connect(_on_leaderboard_received)

func show_panel() -> void:
	visible = true
	_on_refresh()

func _on_close() -> void:
	visible = false

func _on_refresh() -> void:
	GameManager.send_message("season_leaderboard_request", {})
	_season_label.text = "載入中..."

# ── 排行榜資料渲染 ────────────────────────────────────────────

func _on_leaderboard_received(data: Dictionary) -> void:
	var season_id = data.get("season_id", "")
	var top20 = data.get("top20", [])
	var my_rank = data.get("my_rank", -1)
	var my_entry = data.get("my_entry", null)
	
	_season_label.text = "賽季：%s  |  共 %d 名玩家" % [season_id, top20.size()]
	
	# 清空列表
	for child in _list_container.get_children():
		child.queue_free()
	
	if top20.size() == 0:
		var empty_label = Label.new()
		empty_label.text = "本賽季尚無排名資料"
		empty_label.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))
		empty_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		_list_container.add_child(empty_label)
	else:
		for entry in top20:
			_list_container.add_child(_create_entry_row(entry))
	
	# 更新我的排名
	if my_rank > 0 and my_entry != null:
		var badge = my_entry.get("badge", "")
		var level_name = my_entry.get("level_name", "")
		var xp = my_entry.get("season_xp", 0)
		_my_rank_label.text = "我的排名：第 %d 名  %s %s  XP: %d" % [my_rank, badge, level_name, xp]
	elif my_rank == -1:
		_my_rank_label.text = "我的排名：尚未上榜（繼續遊玩獲得 XP！）"
	else:
		_my_rank_label.text = "我的排名：第 %d 名" % my_rank

func _create_entry_row(entry: Dictionary) -> Control:
	var rank = entry.get("rank", 0)
	var display_name = entry.get("display_name", "玩家")
	var level = entry.get("level", 1)
	var level_name = entry.get("level_name", "")
	var badge = entry.get("badge", "")
	var xp = entry.get("season_xp", 0)
	var player_id = entry.get("player_id", "")
	
	# 判斷是否是自己
	var is_me = player_id == GameManager.get_player_id()
	
	var panel = PanelContainer.new()
	var style = StyleBoxFlat.new()
	if is_me:
		style.bg_color = Color(0.2, 0.15, 0.35, 0.9)
		style.border_color = Color(0.8, 0.6, 1.0, 0.8)
		style.border_width_left = 1
		style.border_width_right = 1
		style.border_width_top = 1
		style.border_width_bottom = 1
	else:
		style.bg_color = Color(0.1, 0.08, 0.18, 0.7)
	style.corner_radius_top_left = 4
	style.corner_radius_top_right = 4
	style.corner_radius_bottom_left = 4
	style.corner_radius_bottom_right = 4
	panel.add_theme_stylebox_override("panel", style)
	panel.custom_minimum_size = Vector2(0, ROW_HEIGHT)
	
	var row = HBoxContainer.new()
	row.add_theme_constant_override("separation", 0)
	panel.add_child(row)
	
	# 排名
	var rank_label = Label.new()
	var rank_text = str(rank)
	if rank == 1:
		rank_text = "🥇"
	elif rank == 2:
		rank_text = "🥈"
	elif rank == 3:
		rank_text = "🥉"
	rank_label.text = rank_text
	rank_label.custom_minimum_size = Vector2(50, 0)
	rank_label.add_theme_font_size_override("font_size", 16 if rank <= 3 else 14)
	var rank_color = RANK_COLORS.get(rank, Color(0.8, 0.8, 0.8))
	rank_label.add_theme_color_override("font_color", rank_color)
	rank_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	rank_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	row.add_child(rank_label)
	
	# 玩家名稱
	var name_label = Label.new()
	name_label.text = display_name + (" ← 我" if is_me else "")
	name_label.custom_minimum_size = Vector2(200, 0)
	name_label.add_theme_font_size_override("font_size", 13)
	name_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.7) if is_me else Color(0.9, 0.9, 0.9))
	name_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	row.add_child(name_label)
	
	# 等級
	var level_label = Label.new()
	level_label.text = "%s Lv.%d" % [badge, level]
	level_label.custom_minimum_size = Vector2(100, 0)
	level_label.add_theme_font_size_override("font_size", 12)
	level_label.add_theme_color_override("font_color", Color(0.7, 0.9, 0.7))
	level_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	level_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	row.add_child(level_label)
	
	# XP
	var xp_label = Label.new()
	xp_label.text = "%d XP" % xp
	xp_label.custom_minimum_size = Vector2(120, 0)
	xp_label.add_theme_font_size_override("font_size", 13)
	xp_label.add_theme_color_override("font_color", Color(0.9, 0.8, 0.3))
	xp_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	xp_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	row.add_child(xp_label)
	
	return panel
