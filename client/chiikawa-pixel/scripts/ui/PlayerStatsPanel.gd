## PlayerStatsPanel.gd
## 玩家個人統計面板（DAY-096）
## 顯示詳細遊戲統計：最佳記錄、命中率、RTP、Jackpot 中獎等

extends Control

var pixel_font: Font = null
var _stats_data: Dictionary = {}

## 初始化面板（由 HUD.gd 呼叫）
func setup(font: Font) -> void:
	pixel_font = font
	_build_panel()
	GameManager.player_stats_updated.connect(_on_stats_updated)

## 建立面板 UI
func _build_panel() -> void:
	# 背景
	var bg = ColorRect.new()
	bg.name = "StatsBG"
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.04, 0.06, 0.12, 0.92)
	add_child(bg)

	# 頂部標題列
	var title_bar = ColorRect.new()
	title_bar.size = Vector2(400, 28)
	title_bar.position = Vector2(0, 0)
	title_bar.color = Color(0.1, 0.2, 0.4, 0.9)
	add_child(title_bar)

	var title_lbl = Label.new()
	title_lbl.text = "📊 個人統計"
	title_lbl.position = Vector2(8, 4)
	title_lbl.size = Vector2(300, 20)
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
	if is_instance_valid(pixel_font):
		title_lbl.add_theme_font_override("font", pixel_font)
	add_child(title_lbl)

	# 關閉按鈕
	var close_btn = Button.new()
	close_btn.text = "✕"
	close_btn.position = Vector2(370, 2)
	close_btn.size = Vector2(28, 24)
	close_btn.add_theme_font_size_override("font_size", 12)
	close_btn.pressed.connect(func(): visible = false)
	add_child(close_btn)

	# 統計內容區域（ScrollContainer）
	var scroll = ScrollContainer.new()
	scroll.name = "StatsScroll"
	scroll.position = Vector2(0, 30)
	scroll.size = Vector2(400, 370)
	add_child(scroll)

	var content = VBoxContainer.new()
	content.name = "StatsContent"
	content.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	scroll.add_child(content)

	# 初始顯示載入中
	var loading_lbl = Label.new()
	loading_lbl.name = "LoadingLabel"
	loading_lbl.text = "載入統計中..."
	loading_lbl.add_theme_font_size_override("font_size", 11)
	loading_lbl.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	if is_instance_valid(pixel_font):
		loading_lbl.add_theme_font_override("font", pixel_font)
	content.add_child(loading_lbl)

	size = Vector2(400, 400)

## 統計更新
func _on_stats_updated(data: Dictionary) -> void:
	_stats_data = data
	_refresh_display()

## 重新整理顯示
func _refresh_display() -> void:
	var scroll = get_node_or_null("StatsScroll")
	if not is_instance_valid(scroll):
		return
	var content = scroll.get_node_or_null("StatsContent")
	if not is_instance_valid(content):
		return

	# 清除舊內容
	for child in content.get_children():
		child.queue_free()

	var d = _stats_data

	# 計算遊戲時長顯示
	var play_sec = d.get("total_play_time_sec", 0)
	var play_str = _format_time(play_sec)

	# 計算命中率
	var hit_rate = d.get("hit_rate", 0.0)
	var hit_rate_str = "%.1f%%" % (hit_rate * 100.0)

	# 計算 RTP
	var rtp = d.get("rtp", 0.0)
	var rtp_str = "%.1f%%" % (rtp * 100.0)
	var rtp_color = Color(0.4, 1.0, 0.4) if rtp >= 0.9 else (Color(1.0, 0.8, 0.2) if rtp >= 0.7 else Color(1.0, 0.4, 0.4))

	# 各區塊
	_add_section(content, "🎮 遊戲概況")
	_add_stat_row(content, "總場次", "%d 場" % d.get("total_sessions", 0))
	_add_stat_row(content, "總遊戲時間", play_str)
	_add_stat_row(content, "總射擊次數", "%d 次" % d.get("total_shots", 0))
	_add_stat_row(content, "總擊破次數", "%d 次" % d.get("total_kills", 0))
	_add_stat_row(content, "命中率", hit_rate_str)

	_add_section(content, "💰 金幣統計")
	_add_stat_row(content, "總投注", "🪙%d" % d.get("total_bet", 0))
	_add_stat_row(content, "總獲得", "🪙%d" % d.get("total_reward", 0))
	_add_stat_row_colored(content, "實際 RTP", rtp_str, rtp_color)
	_add_stat_row(content, "歷史最高金幣", "🪙%d" % d.get("max_coins", 0))

	_add_section(content, "🏆 最佳記錄")
	_add_stat_row(content, "最高單次倍率", "%.1fx" % d.get("best_multiplier", 0.0))
	_add_stat_row(content, "最高連擊數", "%d 連擊" % d.get("best_streak", 0))
	_add_stat_row(content, "單場最高得分", "🪙%d" % d.get("best_session_score", 0))
	_add_stat_row(content, "單次 Bonus 最高", "🪙%d" % d.get("best_bonus_reward", 0))

	_add_section(content, "⚡ 特殊事件")
	_add_stat_row(content, "觸發 Bonus", "%d 次" % d.get("total_bonuses", 0))
	_add_stat_row(content, "擊殺 BOSS", "%d 次" % d.get("total_boss_kills", 0))

	_add_section(content, "🎰 Jackpot 統計")
	var jackpot_wins = d.get("jackpot_wins", 0)
	_add_stat_row(content, "總中獎次數", "%d 次" % jackpot_wins)
	if jackpot_wins > 0:
		_add_stat_row(content, "🥈 Mini", "%d 次" % d.get("jackpot_mini_wins", 0))
		_add_stat_row(content, "🥇 Minor", "%d 次" % d.get("jackpot_minor_wins", 0))
		_add_stat_row(content, "🔥 Major", "%d 次" % d.get("jackpot_major_wins", 0))
		_add_stat_row(content, "👑 Grand", "%d 次" % d.get("jackpot_grand_wins", 0))
		_add_stat_row(content, "總 Jackpot 獲得", "🪙%d" % d.get("total_jackpot_payout", 0))

## 加入區塊標題
func _add_section(parent: Control, title: String) -> void:
	var sep = ColorRect.new()
	sep.size = Vector2(380, 1)
	sep.color = Color(0.3, 0.4, 0.6, 0.5)
	sep.custom_minimum_size = Vector2(0, 1)
	parent.add_child(sep)

	var lbl = Label.new()
	lbl.text = title
	lbl.add_theme_font_size_override("font_size", 11)
	lbl.add_theme_color_override("font_color", Color(0.6, 0.8, 1.0))
	if is_instance_valid(pixel_font):
		lbl.add_theme_font_override("font", pixel_font)
	parent.add_child(lbl)

## 加入統計行
func _add_stat_row(parent: Control, label: String, value: String) -> void:
	_add_stat_row_colored(parent, label, value, Color(1.0, 0.95, 0.7))

## 加入帶顏色的統計行
func _add_stat_row_colored(parent: Control, label: String, value: String, value_color: Color) -> void:
	var row = HBoxContainer.new()
	row.custom_minimum_size = Vector2(380, 18)
	parent.add_child(row)

	var label_lbl = Label.new()
	label_lbl.text = label
	label_lbl.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	label_lbl.add_theme_font_size_override("font_size", 10)
	label_lbl.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	if is_instance_valid(pixel_font):
		label_lbl.add_theme_font_override("font", pixel_font)
	row.add_child(label_lbl)

	var value_lbl = Label.new()
	value_lbl.text = value
	value_lbl.add_theme_font_size_override("font_size", 10)
	value_lbl.add_theme_color_override("font_color", value_color)
	if is_instance_valid(pixel_font):
		value_lbl.add_theme_font_override("font", pixel_font)
	row.add_child(value_lbl)

## 格式化時間（秒 → 時:分:秒）
func _format_time(seconds: int) -> String:
	if seconds < 60:
		return "%d 秒" % seconds
	elif seconds < 3600:
		return "%d 分 %d 秒" % [seconds / 60, seconds % 60]
	else:
		var h = seconds / 3600
		var m = (seconds % 3600) / 60
		return "%d 時 %d 分" % [h, m]

## 顯示面板並請求最新統計
func show_panel() -> void:
	visible = true
	NetworkManager.send_get_player_stats()
