## DailyBossPanel.gd — 每日特殊 BOSS 挑戰面板（DAY-077）
## 顯示每日 BOSS 的 HP、貢獻排名、倒數計時、獎勵池
## 位置：TopBar 右側（可折疊）
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 280
const PANEL_HEIGHT := 320

# ---- 節點引用 ----
var _pixel_font: Font = null
var _is_open: bool = false
var _toggle_btn: Button = null
var _panel_bg: ColorRect = null
var _content_node: Node2D = null

# ---- 每日 BOSS 資料 ----
var _date_id: String = ""
var _boss_name: String = ""
var _boss_icon: String = "👹"
var _boss_color: String = "#FF4444"
var _description: String = ""
var _max_hp: int = 0
var _current_hp: int = 0
var _hp_percent: float = 1.0
var _status: String = "active"
var _end_at_ms: int = 0
var _reward_pool: int = 0
var _top_contribs: Array = []
var _my_damage: int = 0
var _my_reward: int = 0
var _difficulty_mod: float = 1.0

# ---- 結算資料 ----
var _defeated_data: Dictionary = {}
var _show_defeated: bool = false

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
	_toggle_btn.text = "👹"
	_toggle_btn.size = Vector2(32, 24)
	_toggle_btn.position = Vector2(0, 0)
	_toggle_btn.flat = true
	_toggle_btn.tooltip_text = "每日 BOSS"
	if _pixel_font:
		_toggle_btn.add_theme_font_override("font", _pixel_font)
	add_child(_toggle_btn)
	_toggle_btn.pressed.connect(_on_toggle_pressed)

## 建立面板
func _build_panel() -> void:
	_panel_bg = ColorRect.new()
	_panel_bg.color = Color(0.08, 0.03, 0.03, 0.93)
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
	if GameManager.has_signal("daily_boss_updated"):
		GameManager.daily_boss_updated.connect(_on_daily_boss_updated)
	if GameManager.has_signal("daily_boss_defeated"):
		GameManager.daily_boss_defeated.connect(_on_daily_boss_defeated)

# ---- 事件處理 ----
func _on_toggle_pressed() -> void:
	_is_open = !_is_open
	_panel_bg.visible = _is_open
	_content_node.visible = _is_open
	if _is_open:
		_redraw_panel()
		GameManager.request_daily_boss_status()

func _on_daily_boss_updated(data: Dictionary) -> void:
	_date_id = data.get("date_id", "")
	_boss_name = data.get("boss_name", "???")
	_boss_icon = data.get("boss_icon", "👹")
	_boss_color = data.get("boss_color", "#FF4444")
	_description = data.get("description", "")
	_max_hp = data.get("max_hp", 0)
	_current_hp = data.get("current_hp", 0)
	_hp_percent = data.get("hp_percent", 1.0)
	_status = data.get("status", "active")
	_end_at_ms = data.get("end_at", 0)
	_reward_pool = data.get("reward_pool", 0)
	_top_contribs = data.get("top_contribs", [])
	_my_damage = data.get("my_damage", 0)
	_my_reward = data.get("my_reward", 0)
	_difficulty_mod = data.get("difficulty_mod", 1.0)
	_show_defeated = false
	if _is_open:
		_redraw_panel()
	_update_toggle_badge()

func _on_daily_boss_defeated(data: Dictionary) -> void:
	_defeated_data = data
	_show_defeated = true
	if _is_open:
		_redraw_panel()
	_show_defeat_notification(data)

# ---- 繪製 ----
func _redraw_panel() -> void:
	for child in _content_node.get_children():
		child.queue_free()

	if _show_defeated and not _defeated_data.is_empty():
		_draw_defeated_view()
	else:
		_draw_active_view()

func _draw_active_view() -> void:
	var y := 8.0

	# 標題
	var title_lbl := Label.new()
	title_lbl.text = "每日 BOSS：%s%s" % [_boss_icon, _boss_name]
	title_lbl.position = Vector2(8, y)
	title_lbl.size = Vector2(PANEL_WIDTH - 16, 20)
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
		title_lbl.add_theme_font_size_override("font_size", 12)
	title_lbl.add_theme_color_override("font_color", Color(1.0, 0.4, 0.3))
	_content_node.add_child(title_lbl)
	y += 22

	# 描述
	if _description != "":
		var desc_lbl := Label.new()
		desc_lbl.text = _description
		desc_lbl.position = Vector2(8, y)
		desc_lbl.size = Vector2(PANEL_WIDTH - 16, 28)
		desc_lbl.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
		if _pixel_font:
			desc_lbl.add_theme_font_override("font", _pixel_font)
			desc_lbl.add_theme_font_size_override("font_size", 9)
		desc_lbl.add_theme_color_override("font_color", Color(0.75, 0.75, 0.75))
		_content_node.add_child(desc_lbl)
		y += 30

	# HP 條
	var hp_bg := ColorRect.new()
	hp_bg.color = Color(0.3, 0.1, 0.1, 0.8)
	hp_bg.size = Vector2(PANEL_WIDTH - 16, 16)
	hp_bg.position = Vector2(8, y)
	_content_node.add_child(hp_bg)

	var hp_fill := ColorRect.new()
	var fill_width := (PANEL_WIDTH - 16) * _hp_percent
	hp_fill.color = Color(0.9, 0.2, 0.2) if _hp_percent > 0.5 else (Color(1.0, 0.6, 0.1) if _hp_percent > 0.25 else Color(1.0, 0.2, 0.2))
	hp_fill.size = Vector2(fill_width, 16)
	hp_fill.position = Vector2(8, y)
	_content_node.add_child(hp_fill)

	var hp_lbl := Label.new()
	hp_lbl.name = "HPLabel"
	hp_lbl.text = "%d / %d" % [_current_hp, _max_hp]
	hp_lbl.position = Vector2(8, y)
	hp_lbl.size = Vector2(PANEL_WIDTH - 16, 16)
	hp_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	if _pixel_font:
		hp_lbl.add_theme_font_override("font", _pixel_font)
		hp_lbl.add_theme_font_size_override("font_size", 9)
	hp_lbl.add_theme_color_override("font_color", Color.WHITE)
	_content_node.add_child(hp_lbl)
	y += 20

	# 獎勵池 + 倒數
	var info_lbl := Label.new()
	info_lbl.text = "🪙 獎勵池：%d  |  %s" % [_reward_pool, _get_countdown_text()]
	info_lbl.position = Vector2(8, y)
	info_lbl.size = Vector2(PANEL_WIDTH - 16, 16)
	if _pixel_font:
		info_lbl.add_theme_font_override("font", _pixel_font)
		info_lbl.add_theme_font_size_override("font_size", 9)
	info_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	_content_node.add_child(info_lbl)
	y += 18

	# 我的貢獻
	if _my_damage > 0:
		var my_lbl := Label.new()
		my_lbl.text = "我的傷害：%d" % _my_damage
		my_lbl.position = Vector2(8, y)
		my_lbl.size = Vector2(PANEL_WIDTH - 16, 14)
		if _pixel_font:
			my_lbl.add_theme_font_override("font", _pixel_font)
			my_lbl.add_theme_font_size_override("font_size", 9)
		my_lbl.add_theme_color_override("font_color", Color(0.4, 1.0, 0.6))
		_content_node.add_child(my_lbl)
		y += 16

	# 分隔線
	var sep := ColorRect.new()
	sep.color = Color(0.4, 0.2, 0.2, 0.8)
	sep.size = Vector2(PANEL_WIDTH - 16, 1)
	sep.position = Vector2(8, y)
	_content_node.add_child(sep)
	y += 6

	# 貢獻排名標題
	var rank_title := Label.new()
	rank_title.text = "貢獻排名"
	rank_title.position = Vector2(8, y)
	rank_title.size = Vector2(PANEL_WIDTH - 16, 14)
	if _pixel_font:
		rank_title.add_theme_font_override("font", _pixel_font)
		rank_title.add_theme_font_size_override("font_size", 10)
	rank_title.add_theme_color_override("font_color", Color(0.9, 0.7, 0.3))
	_content_node.add_child(rank_title)
	y += 16

	# 貢獻者列表
	for i in range(min(_top_contribs.size(), 5)):
		var entry = _top_contribs[i]
		var rank: int = entry.get("rank", i + 1)
		var name: String = entry.get("display_name", "???")
		var damage: int = entry.get("damage", 0)
		var is_me: bool = entry.get("is_me", false)

		if is_me:
			var row_bg := ColorRect.new()
			row_bg.color = Color(0.2, 0.4, 0.2, 0.4)
			row_bg.size = Vector2(PANEL_WIDTH - 16, 18)
			row_bg.position = Vector2(8, y - 2)
			_content_node.add_child(row_bg)

		var rank_icon := "🥇" if rank == 1 else ("🥈" if rank == 2 else ("🥉" if rank == 3 else str(rank) + "."))
		var row_lbl := Label.new()
		row_lbl.text = "%s %s  %d傷" % [rank_icon, name, damage]
		row_lbl.position = Vector2(12, y)
		row_lbl.size = Vector2(PANEL_WIDTH - 24, 16)
		if _pixel_font:
			row_lbl.add_theme_font_override("font", _pixel_font)
			row_lbl.add_theme_font_size_override("font_size", 9)
		var row_color := Color(1.0, 0.9, 0.3) if rank <= 3 else (Color(0.4, 1.0, 0.6) if is_me else Color(0.8, 0.8, 0.8))
		row_lbl.add_theme_color_override("font_color", row_color)
		_content_node.add_child(row_lbl)
		y += 18

	# 難度修正提示
	if _difficulty_mod < 1.0:
		var diff_lbl := Label.new()
		diff_lbl.text = "⬇ 難度已降低 %d%%" % [int((1.0 - _difficulty_mod) * 100)]
		diff_lbl.position = Vector2(8, y + 4)
		diff_lbl.size = Vector2(PANEL_WIDTH - 16, 14)
		if _pixel_font:
			diff_lbl.add_theme_font_override("font", _pixel_font)
			diff_lbl.add_theme_font_size_override("font_size", 9)
		diff_lbl.add_theme_color_override("font_color", Color(0.6, 0.8, 1.0))
		_content_node.add_child(diff_lbl)

func _draw_defeated_view() -> void:
	var y := 8.0
	var boss_name: String = _defeated_data.get("boss_name", "???")
	var boss_icon: String = _defeated_data.get("boss_icon", "👹")
	var killer_name: String = _defeated_data.get("killer_name", "???")
	var my_reward: int = _defeated_data.get("my_reward", 0)
	var rankings = _defeated_data.get("rankings", [])

	# 標題
	var title_lbl := Label.new()
	title_lbl.text = "%s%s 已被擊敗！" % [boss_icon, boss_name]
	title_lbl.position = Vector2(8, y)
	title_lbl.size = Vector2(PANEL_WIDTH - 16, 20)
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
		title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	_content_node.add_child(title_lbl)
	y += 24

	# 擊殺者
	var killer_lbl := Label.new()
	killer_lbl.text = "最後一擊：%s" % killer_name
	killer_lbl.position = Vector2(8, y)
	killer_lbl.size = Vector2(PANEL_WIDTH - 16, 16)
	if _pixel_font:
		killer_lbl.add_theme_font_override("font", _pixel_font)
		killer_lbl.add_theme_font_size_override("font_size", 10)
	killer_lbl.add_theme_color_override("font_color", Color(1.0, 0.5, 0.3))
	_content_node.add_child(killer_lbl)
	y += 20

	# 我的獎勵
	if my_reward > 0:
		var reward_lbl := Label.new()
		reward_lbl.text = "你獲得：🪙 %d 金幣！" % my_reward
		reward_lbl.position = Vector2(8, y)
		reward_lbl.size = Vector2(PANEL_WIDTH - 16, 18)
		if _pixel_font:
			reward_lbl.add_theme_font_override("font", _pixel_font)
			reward_lbl.add_theme_font_size_override("font_size", 12)
		reward_lbl.add_theme_color_override("font_color", Color(0.4, 1.0, 0.6))
		_content_node.add_child(reward_lbl)
		y += 22

	# 前三名
	var sep := ColorRect.new()
	sep.color = Color(0.4, 0.2, 0.1, 0.8)
	sep.size = Vector2(PANEL_WIDTH - 16, 1)
	sep.position = Vector2(8, y)
	_content_node.add_child(sep)
	y += 8

	for i in range(min(rankings.size(), 5)):
		var entry = rankings[i]
		var rank: int = entry.get("rank", i + 1)
		var name: String = entry.get("display_name", "???")
		var damage: int = entry.get("damage", 0)
		var reward: int = entry.get("reward", 0)
		var is_me: bool = entry.get("is_me", false)

		var rank_icon := "🥇" if rank == 1 else ("🥈" if rank == 2 else ("🥉" if rank == 3 else str(rank) + "."))
		var row_lbl := Label.new()
		row_lbl.text = "%s %s  %d傷  +%d🪙" % [rank_icon, name, damage, reward]
		row_lbl.position = Vector2(12, y)
		row_lbl.size = Vector2(PANEL_WIDTH - 24, 16)
		if _pixel_font:
			row_lbl.add_theme_font_override("font", _pixel_font)
			row_lbl.add_theme_font_size_override("font_size", 9)
		var row_color := Color(1.0, 0.9, 0.3) if rank <= 3 else (Color(0.4, 1.0, 0.6) if is_me else Color(0.8, 0.8, 0.8))
		row_lbl.add_theme_color_override("font_color", row_color)
		_content_node.add_child(row_lbl)
		y += 18

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
		_show_defeated = false
		_redraw_panel()
	)

# ---- 工具函數 ----
func _get_countdown_text() -> String:
	if _end_at_ms <= 0:
		return "等待重置..."
	var now_unix_ms := int(Time.get_unix_time_from_system() * 1000)
	var remaining_ms := _end_at_ms - now_unix_ms
	if remaining_ms <= 0:
		return "重置中..."
	var remaining_sec := remaining_ms / 1000
	var hours := remaining_sec / 3600
	var minutes := (remaining_sec % 3600) / 60
	if hours > 0:
		return "%d時%d分" % [hours, minutes]
	else:
		return "%d分鐘" % minutes

func _update_toggle_badge() -> void:
	if _status == "defeated":
		_toggle_btn.text = "👹✓"
	elif _my_damage > 0:
		_toggle_btn.text = "👹⚔"
	else:
		_toggle_btn.text = "👹"

func _show_defeat_notification(data: Dictionary) -> void:
	var boss_name: String = data.get("boss_name", "???")
	var boss_icon: String = data.get("boss_icon", "👹")
	var my_reward: int = data.get("my_reward", 0)
	if my_reward > 0:
		if get_parent() and get_parent().has_method("show_achievement_notify"):
			var msg := "%s%s 已擊敗！你獲得 %d 金幣！" % [boss_icon, boss_name, my_reward]
			get_parent().show_achievement_notify(msg, "gold")
	else:
		if get_parent() and get_parent().has_method("show_achievement_notify"):
			var msg := "%s%s 已被擊敗！" % [boss_icon, boss_name]
			get_parent().show_achievement_notify(msg, "normal")

# ---- 每幀更新 ----
func _process(_delta: float) -> void:
	if _is_open and not _show_defeated:
		var hp_lbl = _content_node.find_child("HPLabel", false, false)
		if hp_lbl:
			hp_lbl.text = "%d / %d" % [_current_hp, _max_hp]
