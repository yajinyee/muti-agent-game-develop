## GuildWarPanel.gd — 公會戰面板（DAY-076）
## 顯示本週公會戰排名、積分、倒數計時
## 位置：TopBar 右側（可折疊）
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 280
const PANEL_HEIGHT := 300
const MAX_DISPLAY_GUILDS := 10

# ---- 節點引用 ----
var _pixel_font: Font = null
var _is_open: bool = false
var _toggle_btn: Button = null
var _panel_bg: ColorRect = null
var _content_node: Node2D = null

# ---- 公會戰資料 ----
var _week_id: String = ""
var _status: String = "active"
var _end_at_ms: int = 0
var _rankings: Array = []  # Array of GuildWarScoreEntry
var _my_guild_rank: int = 0
var _my_guild_score: int = 0
var _total_guilds: int = 0

# ---- 結算結果 ----
var _last_result: Dictionary = {}
var _show_result: bool = false

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_build_toggle_btn()
	_build_panel()
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

## 建立折疊按鈕
func _build_toggle_btn() -> void:
	_toggle_btn = Button.new()
	_toggle_btn.text = "🏆"
	_toggle_btn.size = Vector2(32, 24)
	_toggle_btn.position = Vector2(0, 0)
	_toggle_btn.flat = true
	_toggle_btn.tooltip_text = "公會戰"
	if _pixel_font:
		_toggle_btn.add_theme_font_override("font", _pixel_font)
	add_child(_toggle_btn)
	_toggle_btn.pressed.connect(_on_toggle_pressed)

## 建立面板
func _build_panel() -> void:
	_panel_bg = ColorRect.new()
	_panel_bg.color = Color(0.05, 0.05, 0.15, 0.92)
	_panel_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_panel_bg.position = Vector2(0, 28)
	_panel_bg.visible = false
	add_child(_panel_bg)

	_content_node = Node2D.new()
	_content_node.position = Vector2(0, 28)
	_content_node.visible = false
	add_child(_content_node)

## 連接訊號
func _connect_signals() -> void:
	if GameManager.has_signal("guild_war_updated"):
		GameManager.guild_war_updated.connect(_on_guild_war_updated)
	if GameManager.has_signal("guild_war_result"):
		GameManager.guild_war_result.connect(_on_guild_war_result)

# ---- 事件處理 ----
func _on_toggle_pressed() -> void:
	_is_open = !_is_open
	_panel_bg.visible = _is_open
	_content_node.visible = _is_open
	if _is_open:
		_redraw_panel()
		# 請求最新狀態
		if GameManager.has_method("request_guild_war_status"):
			GameManager.request_guild_war_status()

func _on_guild_war_updated(data: Dictionary) -> void:
	_week_id = data.get("week_id", "")
	_status = data.get("status", "active")
	_end_at_ms = data.get("end_at", 0)
	_rankings = data.get("rankings", [])
	_my_guild_rank = data.get("my_guild_rank", 0)
	_my_guild_score = data.get("my_guild_score", 0)
	_total_guilds = data.get("total_guilds", 0)
	_show_result = false
	if _is_open:
		_redraw_panel()
	_update_toggle_badge()

func _on_guild_war_result(data: Dictionary) -> void:
	_last_result = data
	_show_result = true
	if _is_open:
		_redraw_panel()
	# 顯示結算通知
	_show_result_notification(data)

# ---- 繪製 ----
func _redraw_panel() -> void:
	# 清除舊內容
	for child in _content_node.get_children():
		child.queue_free()

	if _show_result and not _last_result.is_empty():
		_draw_result_view()
	else:
		_draw_ranking_view()

func _draw_ranking_view() -> void:
	var y := 8.0

	# 標題
	var title_lbl := Label.new()
	title_lbl.text = "🏆 公會戰 — " + _week_id
	title_lbl.position = Vector2(8, y)
	title_lbl.size = Vector2(PANEL_WIDTH - 16, 20)
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
		title_lbl.add_theme_font_size_override("font_size", 12)
	title_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	_content_node.add_child(title_lbl)
	y += 22

	# 倒數計時
	var countdown_lbl := Label.new()
	countdown_lbl.name = "CountdownLabel"
	countdown_lbl.text = _get_countdown_text()
	countdown_lbl.position = Vector2(8, y)
	countdown_lbl.size = Vector2(PANEL_WIDTH - 16, 16)
	if _pixel_font:
		countdown_lbl.add_theme_font_override("font", _pixel_font)
		countdown_lbl.add_theme_font_size_override("font_size", 10)
	countdown_lbl.add_theme_color_override("font_color", Color(0.7, 0.9, 1.0))
	_content_node.add_child(countdown_lbl)
	y += 18

	# 分隔線
	var sep := ColorRect.new()
	sep.color = Color(0.3, 0.3, 0.5, 0.8)
	sep.size = Vector2(PANEL_WIDTH - 16, 1)
	sep.position = Vector2(8, y)
	_content_node.add_child(sep)
	y += 6

	# 我的排名摘要
	if _my_guild_rank > 0:
		var my_lbl := Label.new()
		my_lbl.text = "我的公會：第 %d 名 (%d 分)" % [_my_guild_rank, _my_guild_score]
		my_lbl.position = Vector2(8, y)
		my_lbl.size = Vector2(PANEL_WIDTH - 16, 16)
		if _pixel_font:
			my_lbl.add_theme_font_override("font", _pixel_font)
			my_lbl.add_theme_font_size_override("font_size", 10)
		my_lbl.add_theme_color_override("font_color", Color(0.4, 1.0, 0.6))
		_content_node.add_child(my_lbl)
		y += 18

	# 排名列表
	var display_count := min(_rankings.size(), MAX_DISPLAY_GUILDS)
	for i in range(display_count):
		var entry = _rankings[i]
		var rank: int = entry.get("rank", i + 1)
		var guild_name: String = entry.get("guild_name", "???")
		var guild_icon: String = entry.get("guild_icon", "⚔️")
		var score: int = entry.get("score", 0)
		var is_my_guild: bool = entry.get("is_my_guild", false)

		# 排名行背景（我的公會高亮）
		if is_my_guild:
			var row_bg := ColorRect.new()
			row_bg.color = Color(0.2, 0.4, 0.2, 0.5)
			row_bg.size = Vector2(PANEL_WIDTH - 16, 20)
			row_bg.position = Vector2(8, y - 2)
			_content_node.add_child(row_bg)

		# 排名圖示
		var rank_icon := "🥇" if rank == 1 else ("🥈" if rank == 2 else ("🥉" if rank == 3 else str(rank) + "."))
		var row_lbl := Label.new()
		row_lbl.text = "%s %s%s  %d分" % [rank_icon, guild_icon, guild_name, score]
		row_lbl.position = Vector2(12, y)
		row_lbl.size = Vector2(PANEL_WIDTH - 24, 18)
		if _pixel_font:
			row_lbl.add_theme_font_override("font", _pixel_font)
			row_lbl.add_theme_font_size_override("font_size", 10)
		var row_color := Color(1.0, 0.9, 0.3) if rank <= 3 else (Color(0.4, 1.0, 0.6) if is_my_guild else Color(0.85, 0.85, 0.85))
		row_lbl.add_theme_color_override("font_color", row_color)
		_content_node.add_child(row_lbl)
		y += 22

	# 總參與公會數
	if _total_guilds > 0:
		var total_lbl := Label.new()
		total_lbl.text = "共 %d 個公會參戰" % _total_guilds
		total_lbl.position = Vector2(8, y + 4)
		total_lbl.size = Vector2(PANEL_WIDTH - 16, 14)
		if _pixel_font:
			total_lbl.add_theme_font_override("font", _pixel_font)
			total_lbl.add_theme_font_size_override("font_size", 9)
		total_lbl.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
		_content_node.add_child(total_lbl)

func _draw_result_view() -> void:
	var y := 8.0
	var rankings = _last_result.get("rankings", [])
	var my_rank: int = _last_result.get("my_rank", 0)
	var my_reward: int = _last_result.get("my_reward", 0)

	# 標題
	var title_lbl := Label.new()
	title_lbl.text = "🏆 公會戰結算！"
	title_lbl.position = Vector2(8, y)
	title_lbl.size = Vector2(PANEL_WIDTH - 16, 20)
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
		title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	_content_node.add_child(title_lbl)
	y += 24

	# 我的結果
	if my_rank > 0:
		var my_lbl := Label.new()
		my_lbl.text = "你的公會：第 %d 名" % my_rank
		my_lbl.position = Vector2(8, y)
		my_lbl.size = Vector2(PANEL_WIDTH - 16, 18)
		if _pixel_font:
			my_lbl.add_theme_font_override("font", _pixel_font)
			my_lbl.add_theme_font_size_override("font_size", 12)
		my_lbl.add_theme_color_override("font_color", Color(0.4, 1.0, 0.6))
		_content_node.add_child(my_lbl)
		y += 20

		if my_reward > 0:
			var reward_lbl := Label.new()
			reward_lbl.text = "獲得獎勵：🪙 %d" % my_reward
			reward_lbl.position = Vector2(8, y)
			reward_lbl.size = Vector2(PANEL_WIDTH - 16, 18)
			if _pixel_font:
				reward_lbl.add_theme_font_override("font", _pixel_font)
				reward_lbl.add_theme_font_size_override("font_size", 12)
			reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
			_content_node.add_child(reward_lbl)
			y += 22

	# 前三名
	var sep := ColorRect.new()
	sep.color = Color(0.3, 0.3, 0.5, 0.8)
	sep.size = Vector2(PANEL_WIDTH - 16, 1)
	sep.position = Vector2(8, y)
	_content_node.add_child(sep)
	y += 8

	for i in range(min(rankings.size(), 3)):
		var entry = rankings[i]
		var rank: int = entry.get("rank", i + 1)
		var guild_name: String = entry.get("guild_name", "???")
		var guild_icon: String = entry.get("guild_icon", "⚔️")
		var score: int = entry.get("score", 0)
		var reward: int = entry.get("reward", 0)

		var rank_icon := "🥇" if rank == 1 else ("🥈" if rank == 2 else "🥉")
		var row_lbl := Label.new()
		row_lbl.text = "%s %s%s  %d分  +%d🪙" % [rank_icon, guild_icon, guild_name, score, reward]
		row_lbl.position = Vector2(12, y)
		row_lbl.size = Vector2(PANEL_WIDTH - 24, 18)
		if _pixel_font:
			row_lbl.add_theme_font_override("font", _pixel_font)
			row_lbl.add_theme_font_size_override("font_size", 10)
		row_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
		_content_node.add_child(row_lbl)
		y += 20

	# 關閉按鈕
	var close_btn := Button.new()
	close_btn.text = "關閉"
	close_btn.size = Vector2(80, 24)
	close_btn.position = Vector2((PANEL_WIDTH - 80) / 2, y + 8)
	if _pixel_font:
		close_btn.add_theme_font_override("font", _pixel_font)
		close_btn.add_theme_font_size_override("font_size", 10)
	_content_node.add_child(close_btn)
	close_btn.pressed.connect(func():
		_show_result = false
		_redraw_panel()
	)

# ---- 工具函數 ----
func _get_countdown_text() -> String:
	if _end_at_ms <= 0:
		return "等待下一場..."
	var now_ms := Time.get_ticks_msec() + int(Time.get_unix_time_from_system() * 1000) - Time.get_ticks_msec()
	var now_unix_ms := int(Time.get_unix_time_from_system() * 1000)
	var remaining_ms := _end_at_ms - now_unix_ms
	if remaining_ms <= 0:
		return "結算中..."
	var remaining_sec := remaining_ms / 1000
	var days := remaining_sec / 86400
	var hours := (remaining_sec % 86400) / 3600
	var minutes := (remaining_sec % 3600) / 60
	if days > 0:
		return "剩餘：%d天 %d時 %d分" % [days, hours, minutes]
	elif hours > 0:
		return "剩餘：%d時 %d分" % [hours, minutes]
	else:
		return "剩餘：%d分鐘" % minutes

func _update_toggle_badge() -> void:
	if _my_guild_rank > 0 and _my_guild_rank <= 3:
		_toggle_btn.text = "🏆" + str(_my_guild_rank)
	else:
		_toggle_btn.text = "🏆"

func _show_result_notification(data: Dictionary) -> void:
	var my_rank: int = data.get("my_rank", 0)
	var my_reward: int = data.get("my_reward", 0)
	if my_rank > 0 and my_reward > 0:
		# 透過 HUD 的成就通知系統顯示
		if get_parent() and get_parent().has_method("show_achievement_notify"):
			var msg := "🏆 公會戰結算！第 %d 名，獲得 %d 金幣！" % [my_rank, my_reward]
			get_parent().show_achievement_notify(msg, "gold")

# ---- 每幀更新倒數 ----
func _process(_delta: float) -> void:
	if _is_open and not _show_result:
		var countdown_node = _content_node.find_child("CountdownLabel", false, false)
		if countdown_node:
			countdown_node.text = _get_countdown_text()
