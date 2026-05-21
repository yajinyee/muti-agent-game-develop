## CaptainFishPanel.gd — 船長魚全服競速模式面板（DAY-163）
## 擊破船長魚後觸發全服競速 30 秒，擊破最多目標獲得大獎
## 視覺：藍色海軍主題 + 倒數計時 + 即時排行榜 + 結果彈窗
## 業界依據：King of Ocean 2026「Captain Fish trigger bonus rounds」
extends Node2D

var _pixel_font: Font = null
var _countdown_label: Label = null
var _leaderboard_bg: ColorRect = null
var _leaderboard_labels: Array = []
var _is_active: bool = false
var _remaining: float = 0.0
var _my_player_id: String = ""

func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

func _connect_signals() -> void:
	if GameManager.has_signal("captain_fish_race"):
		GameManager.captain_fish_race.connect(_on_captain_fish_race)

func _process(delta: float) -> void:
	if not _is_active:
		return
	_remaining -= delta
	if _remaining <= 0.0:
		_remaining = 0.0
		return
	if is_instance_valid(_countdown_label):
		_countdown_label.text = "⚓ 競速 %.0f秒" % _remaining

## 處理船長魚競速事件
func _on_captain_fish_race(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var duration: float = data.get("duration_secs", 30.0)
	var remaining: float = data.get("remaining_time", 30.0)
	var killer_name: String = data.get("killer_name", "")
	var killer_id: String = data.get("killer_id", "")
	var entries: Array = data.get("entries", [])

	_my_player_id = NetworkManager.get_player_id() if NetworkManager.has_method("get_player_id") else ""

	match phase:
		"race_start":
			_show_race_start(killer_name, killer_id, duration)
		"race_update":
			_update_leaderboard(entries, remaining)
		"race_end":
			_show_race_end(entries)
		"race_reward":
			_show_my_reward(data)

## 顯示競速開始
func _show_race_start(killer_name: String, killer_id: String, duration: float) -> void:
	_is_active = true
	_remaining = duration
	var is_me = (killer_id == _my_player_id)

	# 全螢幕藍色閃光
	var flash := ColorRect.new()
	flash.position = Vector2(-640, -360)
	flash.size = Vector2(1280, 720)
	flash.color = Color(0.1, 0.3, 0.9, 0.4 if is_me else 0.25)
	add_child(flash)
	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", 0.0, 0.5)
	flash_tween.tween_callback(func():
		if is_instance_valid(flash): flash.queue_free()
	)

	# 頂部橫幅（藍色海軍主題）
	var banner_bg := ColorRect.new()
	banner_bg.name = "CaptainBanner"
	banner_bg.position = Vector2(-640, -360)
	banner_bg.size = Vector2(1280, 52)
	banner_bg.color = Color(0.0, 0.08, 0.25, 0.92)
	add_child(banner_bg)

	var banner_text: String
	if is_me:
		banner_text = "⚓ 你擊破船長魚！全服競速開始！30 秒內擊破最多目標獲得大獎！"
	else:
		banner_text = "⚓ %s 擊破船長魚！全服競速開始！30 秒內擊破最多目標獲得大獎！" % killer_name

	var banner_lbl := Label.new()
	banner_lbl.name = "CaptainBannerLabel"
	banner_lbl.position = Vector2(-640, -354)
	banner_lbl.size = Vector2(1280, 44)
	banner_lbl.text = banner_text
	banner_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	banner_lbl.add_theme_color_override("font_color", Color(0.5, 0.85, 1.0))
	banner_lbl.add_theme_font_size_override("font_size", 16)
	if _pixel_font:
		banner_lbl.add_theme_font_override("font", _pixel_font)
	add_child(banner_lbl)

	# 橫幅滑入
	banner_bg.position.y = -410
	banner_lbl.position.y = -404
	var slide = banner_bg.create_tween()
	slide.tween_property(banner_bg, "position:y", -360.0, 0.3)
	var slide2 = banner_lbl.create_tween()
	slide2.tween_property(banner_lbl, "position:y", -354.0, 0.3)

	# 倒數計時器（右上角）
	var countdown_bg := ColorRect.new()
	countdown_bg.name = "CaptainCountdownBG"
	countdown_bg.position = Vector2(510, -355)
	countdown_bg.size = Vector2(120, 36)
	countdown_bg.color = Color(0.0, 0.06, 0.2, 0.9)
	add_child(countdown_bg)

	_countdown_label = Label.new()
	_countdown_label.name = "CaptainCountdown"
	_countdown_label.position = Vector2(510, -352)
	_countdown_label.size = Vector2(120, 30)
	_countdown_label.text = "⚓ 競速 %.0f秒" % duration
	_countdown_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_countdown_label.add_theme_color_override("font_color", Color(0.5, 0.85, 1.0))
	_countdown_label.add_theme_font_size_override("font_size", 13)
	if _pixel_font:
		_countdown_label.add_theme_font_override("font", _pixel_font)
	add_child(_countdown_label)

	# 左側排行榜面板
	_leaderboard_bg = ColorRect.new()
	_leaderboard_bg.name = "CaptainLeaderboard"
	_leaderboard_bg.position = Vector2(-630, -200)
	_leaderboard_bg.size = Vector2(180, 160)
	_leaderboard_bg.color = Color(0.0, 0.06, 0.2, 0.88)
	add_child(_leaderboard_bg)

	var lb_title := Label.new()
	lb_title.name = "CaptainLBTitle"
	lb_title.position = Vector2(-630, -198)
	lb_title.size = Vector2(180, 22)
	lb_title.text = "🏆 競速排行"
	lb_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lb_title.add_theme_color_override("font_color", Color(0.5, 0.85, 1.0))
	lb_title.add_theme_font_size_override("font_size", 13)
	if _pixel_font:
		lb_title.add_theme_font_override("font", _pixel_font)
	add_child(lb_title)

## 更新排行榜
func _update_leaderboard(entries: Array, remaining: float) -> void:
	_remaining = remaining
	# 清除舊排行榜條目
	for lbl in _leaderboard_labels:
		if is_instance_valid(lbl): lbl.queue_free()
	_leaderboard_labels.clear()

	# 顯示前 5 名
	for i in range(min(entries.size(), 5)):
		var entry = entries[i]
		var rank_icon = ["🥇", "🥈", "🥉", "4.", "5."][i]
		var name_str: String = entry.get("player_name", "?")
		var kills: int = entry.get("kill_count", 0)
		var is_me = (entry.get("player_id", "") == _my_player_id)

		var lbl := Label.new()
		lbl.text = "%s %s: %d" % [rank_icon, name_str.substr(0, 6), kills]
		lbl.position = Vector2(-628, -174 + i * 22)
		lbl.size = Vector2(176, 20)
		lbl.add_theme_color_override("font_color", Color(1.0, 0.95, 0.3) if is_me else Color(0.8, 0.9, 1.0))
		lbl.add_theme_font_size_override("font_size", 11)
		if _pixel_font:
			lbl.add_theme_font_override("font", _pixel_font)
		add_child(lbl)
		_leaderboard_labels.append(lbl)

## 顯示競速結束
func _show_race_end(entries: Array) -> void:
	_is_active = false
	_remaining = 0.0

	# 隱藏倒數計時器
	for node_name in ["CaptainCountdownBG", "CaptainCountdown"]:
		var node = get_node_or_null(node_name)
		if is_instance_valid(node):
			node.queue_free()
	_countdown_label = null

	# 更新排行榜顯示最終結果
	_update_leaderboard(entries, 0.0)

	# 3 秒後淡出所有 UI
	var timer = get_tree().create_timer(3.0)
	timer.timeout.connect(func():
		_hide_all()
	)

## 顯示個人獎勵
func _show_my_reward(data: Dictionary) -> void:
	var my_rank: int = data.get("my_rank", 0)
	var my_bonus: int = data.get("my_bonus", 0)
	var my_kills: int = data.get("my_kill_count", 0)

	if my_bonus <= 0:
		return

	var rank_text = ["🥇 第一名", "🥈 第二名", "🥉 第三名"].get(my_rank - 1, "第%d名" % my_rank)

	# 中央獎勵彈窗
	var popup_bg := ColorRect.new()
	popup_bg.name = "CaptainRewardPopup"
	popup_bg.position = Vector2(-120, -70)
	popup_bg.size = Vector2(240, 120)
	popup_bg.color = Color(0.0, 0.06, 0.2, 0.96)
	add_child(popup_bg)

	var title_lbl := Label.new()
	title_lbl.text = "⚓ 競速結束！%s" % rank_text
	title_lbl.position = Vector2(-120, -68)
	title_lbl.size = Vector2(240, 26)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_color_override("font_color", Color(0.5, 0.85, 1.0))
	title_lbl.add_theme_font_size_override("font_size", 14)
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
	add_child(title_lbl)

	var kills_lbl := Label.new()
	kills_lbl.text = "擊破 %d 個目標" % my_kills
	kills_lbl.position = Vector2(-120, -40)
	kills_lbl.size = Vector2(240, 22)
	kills_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	kills_lbl.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
	kills_lbl.add_theme_font_size_override("font_size", 13)
	if _pixel_font:
		kills_lbl.add_theme_font_override("font", _pixel_font)
	add_child(kills_lbl)

	var bonus_lbl := Label.new()
	bonus_lbl.text = "+%d 金幣獎勵！" % my_bonus
	bonus_lbl.position = Vector2(-120, -16)
	bonus_lbl.size = Vector2(240, 30)
	bonus_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	bonus_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.1))
	bonus_lbl.add_theme_font_size_override("font_size", 20)
	if _pixel_font:
		bonus_lbl.add_theme_font_override("font", _pixel_font)
	add_child(bonus_lbl)

	# 彈跳動畫
	popup_bg.scale = Vector2(0.8, 0.8)
	var scale_tween = popup_bg.create_tween()
	scale_tween.tween_property(popup_bg, "scale", Vector2(1.05, 1.05), 0.15)
	scale_tween.tween_property(popup_bg, "scale", Vector2(1.0, 1.0), 0.08)

	# 4 秒後淡出
	for node in [popup_bg, title_lbl, kills_lbl, bonus_lbl]:
		var fade = node.create_tween()
		fade.tween_interval(4.0)
		fade.tween_property(node, "modulate:a", 0.0, 0.5)
		fade.tween_callback(func():
			if is_instance_valid(node): node.queue_free()
		)

## 隱藏所有 UI
func _hide_all() -> void:
	_is_active = false
	for child_name in ["CaptainBanner", "CaptainBannerLabel", "CaptainCountdownBG",
						"CaptainCountdown", "CaptainLeaderboard", "CaptainLBTitle"]:
		var node = get_node_or_null(child_name)
		if is_instance_valid(node):
			var tween = node.create_tween()
			tween.tween_property(node, "modulate:a", 0.0, 0.4)
			tween.tween_callback(func():
				if is_instance_valid(node): node.queue_free()
			)
	for lbl in _leaderboard_labels:
		if is_instance_valid(lbl):
			var tween = lbl.create_tween()
			tween.tween_property(lbl, "modulate:a", 0.0, 0.4)
			tween.tween_callback(func():
				if is_instance_valid(lbl): lbl.queue_free()
			)
	_leaderboard_labels.clear()
	_leaderboard_bg = null
	_countdown_label = null
