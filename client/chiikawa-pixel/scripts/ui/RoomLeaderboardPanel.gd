## RoomLeaderboardPanel.gd — 同場玩家排行榜（DAY-349）
## 顯示同場在線玩家的賽季 XP 排名、擊破數、最高倍率
## social-ui-agent 負責維護
extends CanvasLayer

# ── 節點引用 ──────────────────────────────────────────────────
var _panel: PanelContainer
var _title_label: Label
var _online_label: Label
var _list_container: VBoxContainer
var _my_rank_label: Label
var _close_btn: Button
var _refresh_btn: Button

# ── 常數 ──────────────────────────────────────────────────────
const PANEL_WIDTH = 480
const PANEL_HEIGHT = 520
const ROW_HEIGHT = 48

# 排名顏色
const RANK_COLORS = {
	1: Color(1.0, 0.85, 0.1),   # 金
	2: Color(0.8, 0.8, 0.85),   # 銀
	3: Color(0.8, 0.5, 0.2),    # 銅
}

const RANK_BADGES = {
	1: "🥇",
	2: "🥈",
	3: "🥉",
}

func _ready() -> void:
	layer = 25  # 在成就面板上方
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
	style.bg_color = Color(0.04, 0.08, 0.16, 0.97)
	style.border_color = Color(0.2, 0.8, 1.0, 1.0)
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
	_title_label.text = "👥 同場排行榜"
	_title_label.add_theme_font_size_override("font_size", 20)
	_title_label.add_theme_color_override("font_color", Color(0.2, 0.9, 1.0))
	_title_label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	header.add_child(_title_label)

	_refresh_btn = Button.new()
	_refresh_btn.text = "🔄"
	_refresh_btn.custom_minimum_size = Vector2(36, 36)
	header.add_child(_refresh_btn)

	_close_btn = Button.new()
	_close_btn.text = "✕"
	_close_btn.custom_minimum_size = Vector2(36, 36)
	header.add_child(_close_btn)

	# 在線人數
	_online_label = Label.new()
	_online_label.text = "👤 在線：0 人"
	_online_label.add_theme_font_size_override("font_size", 13)
	_online_label.add_theme_color_override("font_color", Color(0.5, 0.9, 0.5))
	vbox.add_child(_online_label)

	# 欄位標題
	var col_header = HBoxContainer.new()
	col_header.add_theme_constant_override("separation", 4)
	vbox.add_child(col_header)

	var rank_h = Label.new()
	rank_h.text = "排名"
	rank_h.custom_minimum_size = Vector2(50, 0)
	rank_h.add_theme_font_size_override("font_size", 12)
	rank_h.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	col_header.add_child(rank_h)

	var name_h = Label.new()
	name_h.text = "玩家"
	name_h.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	name_h.add_theme_font_size_override("font_size", 12)
	name_h.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	col_header.add_child(name_h)

	var xp_h = Label.new()
	xp_h.text = "賽季XP"
	xp_h.custom_minimum_size = Vector2(70, 0)
	xp_h.add_theme_font_size_override("font_size", 12)
	xp_h.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	col_header.add_child(xp_h)

	var kills_h = Label.new()
	kills_h.text = "擊破"
	kills_h.custom_minimum_size = Vector2(55, 0)
	kills_h.add_theme_font_size_override("font_size", 12)
	kills_h.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	col_header.add_child(kills_h)

	# 分隔線
	var sep = HSeparator.new()
	vbox.add_child(sep)

	# 排行榜列表（可捲動）
	var scroll = ScrollContainer.new()
	scroll.size_flags_vertical = Control.SIZE_EXPAND_FILL
	scroll.custom_minimum_size = Vector2(0, PANEL_HEIGHT - 200)
	vbox.add_child(scroll)

	_list_container = VBoxContainer.new()
	_list_container.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	_list_container.add_theme_constant_override("separation", 2)
	scroll.add_child(_list_container)

	# 我的排名（底部固定）
	var sep2 = HSeparator.new()
	vbox.add_child(sep2)

	_my_rank_label = Label.new()
	_my_rank_label.text = "我的排名：-"
	_my_rank_label.add_theme_font_size_override("font_size", 14)
	_my_rank_label.add_theme_color_override("font_color", Color(0.9, 0.9, 0.5))
	_my_rank_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(_my_rank_label)

func _connect_signals() -> void:
	if is_instance_valid(_close_btn):
		_close_btn.pressed.connect(_on_close)
	if is_instance_valid(_refresh_btn):
		_refresh_btn.pressed.connect(_on_refresh)
	# 連接 GameManager 訊號
	if GameManager.has_signal("room_leaderboard"):
		GameManager.room_leaderboard.connect(_on_room_leaderboard)

func _on_close() -> void:
	visible = false

func _on_refresh() -> void:
	GameManager.request_room_leaderboard()

func show_panel() -> void:
	visible = true
	_on_refresh()

# ── 排行榜更新 ───────────────────────────────────────────────

func _on_room_leaderboard(data: Dictionary) -> void:
	var entries = data.get("entries", [])
	var my_rank = data.get("my_rank", -1)
	var online_count = data.get("online_count", 0)

	_online_label.text = "👥 在線：%d 人" % online_count

	if my_rank > 0:
		_my_rank_label.text = "我的排名：第 %d 名" % my_rank
	else:
		_my_rank_label.text = "我的排名：未上榜"

	# 清空舊列表
	for child in _list_container.get_children():
		child.queue_free()

	if entries.is_empty():
		var empty_lbl = Label.new()
		empty_lbl.text = "目前沒有其他玩家在線"
		empty_lbl.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))
		empty_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		_list_container.add_child(empty_lbl)
		return

	var my_player_id = GameManager.get_player_id()

	for entry in entries:
		_add_entry_row(entry, my_player_id)

func _add_entry_row(entry: Dictionary, my_player_id: String) -> void:
	var rank = entry.get("rank", 0)
	var player_id = entry.get("player_id", "")
	var display_name = entry.get("display_name", "玩家")
	var season_xp = entry.get("season_xp", 0)
	var total_kills = entry.get("total_kills", 0)
	var best_mult = entry.get("best_mult", 0.0)
	var is_me = (player_id == my_player_id)

	var row = HBoxContainer.new()
	row.custom_minimum_size = Vector2(0, ROW_HEIGHT)
	row.add_theme_constant_override("separation", 4)
	_list_container.add_child(row)

	# 背景（自己的行高亮）
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	if is_me:
		bg.color = Color(0.1, 0.2, 0.35, 0.9)
	elif rank <= 3:
		bg.color = Color(0.08, 0.12, 0.2, 0.8)
	else:
		bg.color = Color(0.04, 0.06, 0.1, 0.6)
	row.add_child(bg)

	# 自己的行加邊框
	if is_me:
		var border = ColorRect.new()
		border.size = Vector2(3, ROW_HEIGHT)
		border.color = Color(0.2, 0.8, 1.0)
		row.add_child(border)

	# 排名
	var rank_lbl = Label.new()
	rank_lbl.custom_minimum_size = Vector2(50, 0)
	rank_lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	if rank <= 3:
		rank_lbl.text = RANK_BADGES.get(rank, str(rank))
		rank_lbl.add_theme_font_size_override("font_size", 20)
	else:
		rank_lbl.text = "#%d" % rank
		rank_lbl.add_theme_font_size_override("font_size", 14)
		rank_lbl.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	row.add_child(rank_lbl)

	# 玩家名稱
	var name_lbl = Label.new()
	name_lbl.text = display_name + (" ← 我" if is_me else "")
	name_lbl.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	name_lbl.add_theme_font_size_override("font_size", 14)
	if rank <= 3:
		name_lbl.add_theme_color_override("font_color", RANK_COLORS.get(rank, Color.WHITE))
	elif is_me:
		name_lbl.add_theme_color_override("font_color", Color(0.4, 0.9, 1.0))
	else:
		name_lbl.add_theme_color_override("font_color", Color(0.85, 0.85, 0.85))
	name_lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	row.add_child(name_lbl)

	# 賽季 XP
	var xp_lbl = Label.new()
	xp_lbl.text = "%d XP" % season_xp
	xp_lbl.custom_minimum_size = Vector2(70, 0)
	xp_lbl.add_theme_font_size_override("font_size", 13)
	xp_lbl.add_theme_color_override("font_color", Color(0.8, 0.6, 1.0))
	xp_lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	row.add_child(xp_lbl)

	# 擊破數
	var kills_lbl = Label.new()
	kills_lbl.text = "%d" % total_kills
	kills_lbl.custom_minimum_size = Vector2(55, 0)
	kills_lbl.add_theme_font_size_override("font_size", 13)
	kills_lbl.add_theme_color_override("font_color", Color(0.9, 0.7, 0.3))
	kills_lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	row.add_child(kills_lbl)
