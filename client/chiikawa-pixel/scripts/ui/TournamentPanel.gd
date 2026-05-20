## TournamentPanel.gd — 週賽 + 每日賽面板（DAY-093 升級）
## 顯示本週排名、今日排名、積分、獎勵資訊
## 每 30 秒由 Server 推送更新
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 260
const PANEL_HEIGHT := 340
const PANEL_X      := 10
const PANEL_Y      := 50

# ---- 節點引用 ----
var _panel_bg: ColorRect
var _title_label: Label
var _toggle_btn: Button
var _tab_weekly_btn: Button
var _tab_daily_btn: Button
var _entries_container: VBoxContainer
var _footer_label: Label
var _my_rank_label: Label
var _time_label: Label

# ---- 狀態 ----
var _is_expanded: bool = false
var _active_tab: String = "daily"  # "weekly" or "daily"（預設顯示每日賽）
var _pixel_font: Font = null

# 週賽資料
var _weekly_rankings: Array = []
var _weekly_my_rank: int = 0
var _weekly_my_points: int = 0
var _weekly_seconds_left: int = 0
var _weekly_total_players: int = 0

# 每日賽資料（DAY-093）
var _daily_rankings: Array = []
var _daily_my_rank: int = 0
var _daily_my_points: int = 0
var _daily_seconds_left: int = 0
var _daily_total_players: int = 0

# ---- 初始化 ----
func _ready() -> void:
	# 載入像素字體
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")

	_build_ui()
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

func _build_ui() -> void:
	# 背景面板
	_panel_bg = ColorRect.new()
	_panel_bg.position = Vector2(PANEL_X, PANEL_Y)
	_panel_bg.size = Vector2(PANEL_WIDTH, 36)  # 折疊時只顯示標題列
	_panel_bg.color = Color(0.05, 0.08, 0.18, 0.88)
	add_child(_panel_bg)

	# 標題列
	var title_bar := ColorRect.new()
	title_bar.position = Vector2(0, 0)
	title_bar.size = Vector2(PANEL_WIDTH, 36)
	title_bar.color = Color(0.1, 0.2, 0.5, 0.95)
	_panel_bg.add_child(title_bar)

	# 標題文字
	_title_label = Label.new()
	_title_label.position = Vector2(8, 6)
	_title_label.text = "🏆 錦標賽"
	_title_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	if _pixel_font:
		_title_label.add_theme_font_override("font", _pixel_font)
		_title_label.add_theme_font_size_override("font_size", 14)
	title_bar.add_child(_title_label)

	# 展開/折疊按鈕
	_toggle_btn = Button.new()
	_toggle_btn.position = Vector2(PANEL_WIDTH - 32, 4)
	_toggle_btn.size = Vector2(28, 28)
	_toggle_btn.text = "▼"
	_toggle_btn.flat = true
	_toggle_btn.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
	title_bar.add_child(_toggle_btn)

	# Tab 切換列（每日賽 / 週賽）
	var tab_bar := HBoxContainer.new()
	tab_bar.position = Vector2(0, 36)
	tab_bar.size = Vector2(PANEL_WIDTH, 24)
	_panel_bg.add_child(tab_bar)

	_tab_daily_btn = Button.new()
	_tab_daily_btn.text = "📅 今日賽"
	_tab_daily_btn.custom_minimum_size = Vector2(PANEL_WIDTH / 2, 24)
	_tab_daily_btn.flat = false
	_tab_daily_btn.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	if _pixel_font:
		_tab_daily_btn.add_theme_font_override("font", _pixel_font)
		_tab_daily_btn.add_theme_font_size_override("font_size", 10)
	tab_bar.add_child(_tab_daily_btn)

	_tab_weekly_btn = Button.new()
	_tab_weekly_btn.text = "📆 週賽"
	_tab_weekly_btn.custom_minimum_size = Vector2(PANEL_WIDTH / 2, 24)
	_tab_weekly_btn.flat = false
	_tab_weekly_btn.add_theme_color_override("font_color", Color(0.7, 0.8, 1.0))
	if _pixel_font:
		_tab_weekly_btn.add_theme_font_override("font", _pixel_font)
		_tab_weekly_btn.add_theme_font_size_override("font_size", 10)
	tab_bar.add_child(_tab_weekly_btn)

	# 我的排名（Tab 下方）
	_my_rank_label = Label.new()
	_my_rank_label.position = Vector2(8, 62)
	_my_rank_label.text = "我的排名：未上榜"
	_my_rank_label.add_theme_color_override("font_color", Color(0.7, 0.9, 1.0))
	if _pixel_font:
		_my_rank_label.add_theme_font_override("font", _pixel_font)
		_my_rank_label.add_theme_font_size_override("font_size", 11)
	_panel_bg.add_child(_my_rank_label)

	# 倒數計時
	_time_label = Label.new()
	_time_label.position = Vector2(PANEL_WIDTH - 90, 62)
	_time_label.text = ""
	_time_label.add_theme_color_override("font_color", Color(0.6, 0.8, 0.6))
	if _pixel_font:
		_time_label.add_theme_font_override("font", _pixel_font)
		_time_label.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_time_label)

	# 排名列表容器（展開時顯示）
	_entries_container = VBoxContainer.new()
	_entries_container.position = Vector2(0, 82)
	_entries_container.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT - 100)
	_entries_container.visible = false
	_panel_bg.add_child(_entries_container)

	# 底部說明
	_footer_label = Label.new()
	_footer_label.position = Vector2(8, PANEL_HEIGHT - 20)
	_footer_label.text = "🥇50000  🥈25000  🥉10000"
	_footer_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	if _pixel_font:
		_footer_label.add_theme_font_override("font", _pixel_font)
		_footer_label.add_theme_font_size_override("font_size", 10)
	_footer_label.visible = false
	_panel_bg.add_child(_footer_label)

func _connect_signals() -> void:
	_toggle_btn.pressed.connect(_on_toggle_pressed)
	_tab_daily_btn.pressed.connect(_on_tab_daily_pressed)
	_tab_weekly_btn.pressed.connect(_on_tab_weekly_pressed)
	# 連接 GameManager 的週賽更新訊號
	if GameManager.has_signal("tournament_updated"):
		GameManager.tournament_updated.connect(_on_tournament_updated)
	# 連接 GameManager 的每日賽更新訊號（DAY-093）
	if GameManager.has_signal("daily_tournament_updated"):
		GameManager.daily_tournament_updated.connect(_on_daily_tournament_updated)

# ---- 訊號處理 ----
func _on_toggle_pressed() -> void:
	_is_expanded = !_is_expanded
	_toggle_btn.text = "▲" if _is_expanded else "▼"
	_entries_container.visible = _is_expanded
	_footer_label.visible = _is_expanded

	if _is_expanded:
		_panel_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
		_rebuild_entries()
	else:
		_panel_bg.size = Vector2(PANEL_WIDTH, 82)

func _on_tab_daily_pressed() -> void:
	_active_tab = "daily"
	_tab_daily_btn.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	_tab_weekly_btn.add_theme_color_override("font_color", Color(0.7, 0.8, 1.0))
	_update_my_rank_label()
	_update_time_label()
	if _is_expanded:
		_rebuild_entries()

func _on_tab_weekly_pressed() -> void:
	_active_tab = "weekly"
	_tab_weekly_btn.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	_tab_daily_btn.add_theme_color_override("font_color", Color(0.7, 0.8, 1.0))
	_update_my_rank_label()
	_update_time_label()
	if _is_expanded:
		_rebuild_entries()

func _on_tournament_updated(data: Dictionary) -> void:
	_weekly_rankings = data.get("rankings", [])
	_weekly_my_rank = data.get("player_rank", 0)
	_weekly_my_points = data.get("player_points", 0)
	_weekly_seconds_left = data.get("seconds_left", 0)
	_weekly_total_players = data.get("total_players", 0)

	if _active_tab == "weekly":
		_update_my_rank_label()
		_update_time_label()
		if _is_expanded:
			_rebuild_entries()

func _on_daily_tournament_updated(data: Dictionary) -> void:
	_daily_rankings = data.get("rankings", [])
	_daily_my_rank = data.get("player_rank", 0)
	_daily_my_points = data.get("player_points", 0)
	_daily_seconds_left = data.get("seconds_left", 0)
	_daily_total_players = data.get("total_players", 0)

	if _active_tab == "daily":
		_update_my_rank_label()
		_update_time_label()
		if _is_expanded:
			_rebuild_entries()

# ---- UI 更新 ----
func _update_my_rank_label() -> void:
	var my_rank := _daily_my_rank if _active_tab == "daily" else _weekly_my_rank
	var my_points := _daily_my_points if _active_tab == "daily" else _weekly_my_points

	if my_rank > 0:
		var rank_icon := _get_rank_icon(my_rank)
		_my_rank_label.text = "%s #%d  %d分" % [rank_icon, my_rank, my_points]
		if my_rank <= 3:
			_my_rank_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
		else:
			_my_rank_label.add_theme_color_override("font_color", Color(0.7, 0.9, 1.0))
	else:
		_my_rank_label.text = "我的積分：%d分" % my_points if my_points > 0 else "我的排名：未上榜"
		_my_rank_label.add_theme_color_override("font_color", Color(0.6, 0.7, 0.8))

func _update_time_label() -> void:
	var seconds_left := _daily_seconds_left if _active_tab == "daily" else _weekly_seconds_left
	if seconds_left <= 0:
		_time_label.text = "結算中..."
		return
	var days := seconds_left / 86400
	var hours := (seconds_left % 86400) / 3600
	var mins := (seconds_left % 3600) / 60
	if days > 0:
		_time_label.text = "%dd%dh" % [days, hours]
	elif hours > 0:
		_time_label.text = "%dh%dm" % [hours, mins]
	else:
		_time_label.text = "%dm" % mins

func _rebuild_entries() -> void:
	for child in _entries_container.get_children():
		child.queue_free()

	var rankings := _daily_rankings if _active_tab == "daily" else _weekly_rankings

	if rankings.is_empty():
		var empty_label := Label.new()
		empty_label.text = "  尚無參賽者"
		empty_label.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))
		if _pixel_font:
			empty_label.add_theme_font_override("font", _pixel_font)
			empty_label.add_theme_font_size_override("font_size", 11)
		_entries_container.add_child(empty_label)

		# 更新底部獎勵說明
		if _active_tab == "daily":
			_footer_label.text = "🥇5000  🥈2000  🥉1000"
		else:
			_footer_label.text = "🥇50000  🥈25000  🥉10000"
		return

	for entry in rankings:
		_add_entry_row(entry)

	# 更新底部獎勵說明
	if _active_tab == "daily":
		_footer_label.text = "🥇5000  🥈2000  🥉1000"
	else:
		_footer_label.text = "🥇50000  🥈25000  🥉10000"

func _add_entry_row(entry: Dictionary) -> void:
	var rank: int = entry.get("rank", 0)
	var display_name: String = entry.get("display_name", "???")
	var points: int = entry.get("points", 0)
	var prize: int = entry.get("prize", 0)
	var is_self: bool = entry.get("is_self", false)

	var row := HBoxContainer.new()
	row.custom_minimum_size = Vector2(PANEL_WIDTH - 8, 22)

	# 背景（自己高亮）
	var row_bg := ColorRect.new()
	row_bg.color = Color(0.2, 0.4, 0.8, 0.3) if is_self else Color(0.0, 0.0, 0.0, 0.0)
	row_bg.size = Vector2(PANEL_WIDTH - 8, 22)
	row.add_child(row_bg)

	# 排名圖示
	var rank_label := Label.new()
	rank_label.text = _get_rank_icon(rank)
	rank_label.custom_minimum_size = Vector2(28, 22)
	rank_label.add_theme_color_override("font_color", _get_rank_color(rank))
	if _pixel_font:
		rank_label.add_theme_font_override("font", _pixel_font)
		rank_label.add_theme_font_size_override("font_size", 12)
	row.add_child(rank_label)

	# 玩家名稱
	var name_label := Label.new()
	# 截斷長名稱
	var short_name := display_name if display_name.length() <= 8 else display_name.substr(0, 7) + "…"
	name_label.text = short_name
	name_label.custom_minimum_size = Vector2(100, 22)
	name_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.8) if is_self else Color(0.9, 0.9, 0.9))
	if _pixel_font:
		name_label.add_theme_font_override("font", _pixel_font)
		name_label.add_theme_font_size_override("font_size", 11)
	row.add_child(name_label)

	# 積分
	var pts_label := Label.new()
	pts_label.text = "%d" % points
	pts_label.custom_minimum_size = Vector2(50, 22)
	pts_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	pts_label.add_theme_color_override("font_color", Color(0.6, 1.0, 0.6))
	if _pixel_font:
		pts_label.add_theme_font_override("font", _pixel_font)
		pts_label.add_theme_font_size_override("font_size", 11)
	row.add_child(pts_label)

	# 獎勵（前三名顯示）
	if prize > 0:
		var prize_label := Label.new()
		prize_label.text = " 💰"
		prize_label.custom_minimum_size = Vector2(24, 22)
		prize_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
		if _pixel_font:
			prize_label.add_theme_font_override("font", _pixel_font)
			prize_label.add_theme_font_size_override("font_size", 11)
		row.add_child(prize_label)

	_entries_container.add_child(row)

# ---- 工具函數 ----
func _get_rank_icon(rank: int) -> String:
	match rank:
		1: return "🥇"
		2: return "🥈"
		3: return "🥉"
		_: return "#%d" % rank

func _get_rank_color(rank: int) -> Color:
	match rank:
		1: return Color(1.0, 0.85, 0.0)   # 金色
		2: return Color(0.8, 0.8, 0.85)   # 銀色
		3: return Color(0.8, 0.5, 0.2)    # 銅色
		_: return Color(0.7, 0.7, 0.7)    # 灰色

# ---- 每幀更新倒數 ----
func _process(delta: float) -> void:
	if _active_tab == "daily" and _daily_seconds_left > 0:
		_daily_seconds_left -= int(delta)
		if _daily_seconds_left < 0:
			_daily_seconds_left = 0
		_update_time_label()
	elif _active_tab == "weekly" and _weekly_seconds_left > 0:
		_weekly_seconds_left -= int(delta)
		if _weekly_seconds_left < 0:
			_weekly_seconds_left = 0
		_update_time_label()
