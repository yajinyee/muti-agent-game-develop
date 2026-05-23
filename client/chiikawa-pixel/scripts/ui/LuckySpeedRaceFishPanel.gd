## LuckySpeedRaceFishPanel.gd — 幸運競速賽魚 UI（DAY-265）
## 橙紅競速主題面板
## 主色：#FF6B35 橙紅 + #FFD700 金 + #FF4500 火橙 + #FFF3E0 奶油
extends CanvasLayer

const RACE_COLOR_ORANGE = Color(1.0, 0.42, 0.21, 1.0)   # #FF6B35 橙紅
const RACE_COLOR_GOLD   = Color(1.0, 0.85, 0.0, 1.0)    # #FFD700 金
const RACE_COLOR_FIRE   = Color(1.0, 0.27, 0.0, 1.0)    # #FF4500 火橙
const RACE_COLOR_CREAM  = Color(1.0, 0.95, 0.88, 1.0)   # #FFF3E0 奶油

var _flash_overlay: ColorRect
var _banner: PanelContainer
var _banner_label: Label
var _leaderboard_panel: PanelContainer
var _leaderboard_labels: Array = []
var _timer_bar: ColorRect
var _timer_bar_bg: ColorRect
var _timer_duration: float = 30.0
var _timer_elapsed: float = 0.0
var _timer_active: bool = false

func _ready() -> void:
	layer = 38
	_build_ui()

func _build_ui() -> void:
	var vp_size = get_viewport().size

	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 0.42, 0.21, 0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner = PanelContainer.new()
	_banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 56)
	_banner.modulate.a = 0.0
	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.12, 0.04, 0.0, 0.88)
	banner_style.border_color = RACE_COLOR_ORANGE
	banner_style.set_border_width_all(2)
	banner_style.set_corner_radius_all(6)
	_banner.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_font_size_override("font_size", 18)
	_banner_label.add_theme_color_override("font_color", RACE_COLOR_CREAM)
	_banner.add_child(_banner_label)

	# 排行榜面板（右側，顯示前 3 名）
	_leaderboard_panel = PanelContainer.new()
	_leaderboard_panel.custom_minimum_size = Vector2(200, 120)
	_leaderboard_panel.position = Vector2(vp_size.x - 210, 70)
	_leaderboard_panel.modulate.a = 0.0
	var lb_style = StyleBoxFlat.new()
	lb_style.bg_color = Color(0.1, 0.04, 0.0, 0.88)
	lb_style.border_color = RACE_COLOR_GOLD
	lb_style.set_border_width_all(2)
	lb_style.set_corner_radius_all(6)
	_leaderboard_panel.add_theme_stylebox_override("panel", lb_style)
	add_child(_leaderboard_panel)

	var lb_vbox = VBoxContainer.new()
	lb_vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	_leaderboard_panel.add_child(lb_vbox)

	var lb_title = Label.new()
	lb_title.text = "🏁 競速排行榜"
	lb_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lb_title.add_theme_font_size_override("font_size", 14)
	lb_title.add_theme_color_override("font_color", RACE_COLOR_GOLD)
	lb_vbox.add_child(lb_title)

	for i in range(3):
		var rank_label = Label.new()
		rank_label.text = "—"
		rank_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		rank_label.add_theme_font_size_override("font_size", 13)
		rank_label.add_theme_color_override("font_color", RACE_COLOR_CREAM)
		lb_vbox.add_child(rank_label)
		_leaderboard_labels.append(rank_label)

	# 計時條（右側豎向，x=-268 與其他計時條錯開）
	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.12, 0.04, 0.0, 0.7)
	_timer_bar_bg.size = Vector2(8, 200)
	_timer_bar_bg.position = Vector2(vp_size.x - 268, vp_size.y / 2 - 100)
	_timer_bar_bg.modulate.a = 0.0
	add_child(_timer_bar_bg)

	_timer_bar = ColorRect.new()
	_timer_bar.color = RACE_COLOR_ORANGE
	_timer_bar.size = Vector2(8, 200)
	_timer_bar.position = Vector2(vp_size.x - 268, vp_size.y / 2 - 100)
	_timer_bar.modulate.a = 0.0
	add_child(_timer_bar)

func _process(delta: float) -> void:
	if not _timer_active:
		return
	_timer_elapsed += delta
	var pct = 1.0 - (_timer_elapsed / _timer_duration)
	pct = clamp(pct, 0.0, 1.0)
	_timer_bar.size.y = 200.0 * pct

	# 計時條顏色：剩餘時間少時變金色
	if pct < 0.3:
		_timer_bar.color = RACE_COLOR_GOLD
	elif pct < 0.6:
		_timer_bar.color = RACE_COLOR_FIRE
	else:
		_timer_bar.color = RACE_COLOR_ORANGE

	if _timer_elapsed >= _timer_duration:
		_timer_active = false

## 處理來自 GameManager 的競速賽事件
func handle_event(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"race_start":
			_on_race_start(payload)
		"race_broadcast":
			_on_race_broadcast(payload)
		"race_leaderboard":
			_on_race_leaderboard(payload)
		"race_result":
			_on_race_result(payload)

## 競速賽啟動（個人）
func _on_race_start(payload: Dictionary) -> void:
	_timer_duration = float(payload.get("duration_sec", 30))
	_timer_elapsed = 0.0
	_timer_active = true

	var rank1 = payload.get("rank1_mult", 4.0)

	# 橙紅三次強閃光
	_flash_triple(Color(1.0, 0.42, 0.21, 0.5))

	# 顯示橫幅
	_banner_label.text = "🏁 競速賽開始！30秒內擊破最多目標！第1名 ×%.1f！" % rank1
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)
	tween.tween_interval(_timer_duration)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.5)

	# 顯示排行榜和計時條
	_leaderboard_panel.modulate.a = 1.0
	_timer_bar_bg.modulate.a = 1.0
	_timer_bar.modulate.a = 1.0

	# 大字提示
	_show_big_text("🏁 競速賽！", RACE_COLOR_ORANGE)

## 全服廣播
func _on_race_broadcast(payload: Dictionary) -> void:
	var trigger = payload.get("player_name", "玩家")
	var rank1 = payload.get("rank1_mult", 4.0)
	_banner_label.text = "🏁 %s 發起競速賽！第1名 ×%.1f！" % [trigger, rank1]
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)
	tween.tween_interval(3.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.3)

## 即時排行榜廣播（每 5 秒）
func _on_race_leaderboard(payload: Dictionary) -> void:
	var leaderboard = payload.get("leaderboard", [])
	var rank_emojis = ["🥇", "🥈", "🥉"]
	var rank_colors = [RACE_COLOR_GOLD, Color(0.75, 0.75, 0.75), Color(0.8, 0.5, 0.2)]

	for i in range(3):
		if i < leaderboard.size():
			var entry = leaderboard[i]
			var name = entry.get("player_name", "—")
			var score = entry.get("score", 0)
			_leaderboard_labels[i].text = "%s %s: %d 擊" % [rank_emojis[i], name, score]
			_leaderboard_labels[i].add_theme_color_override("font_color", rank_colors[i])
		else:
			_leaderboard_labels[i].text = "—"
			_leaderboard_labels[i].add_theme_color_override("font_color", RACE_COLOR_CREAM)

	# 排行榜更新時輕微閃光
	var tween = create_tween()
	tween.tween_property(_leaderboard_panel, "modulate", Color(1.3, 1.3, 1.0, 1.0), 0.1)
	tween.tween_property(_leaderboard_panel, "modulate", Color(1.0, 1.0, 1.0, 1.0), 0.15)

## 競速賽結算（全服廣播）
func _on_race_result(payload: Dictionary) -> void:
	var winner_name = payload.get("winner_name", "玩家")
	var winner_score = payload.get("winner_score", 0)
	var rank1_mult = payload.get("rank1_mult", 4.0)
	var boost_sec = payload.get("boost_sec", 5)
	var leaderboard = payload.get("leaderboard", [])
	var total = payload.get("total_players", 0)

	_timer_active = false

	# 全螢幕三次強閃光（金色）
	_flash_triple(Color(1.0, 0.85, 0.0, 0.7))

	# 大字顯示冠軍
	_show_big_text("🏆 %s 奪冠！×%.1f" % [winner_name, rank1_mult], RACE_COLOR_GOLD)

	# 結算彈窗
	_show_result_popup(winner_name, winner_score, rank1_mult, boost_sec, leaderboard, total)

	# 隱藏排行榜和計時條
	var tween = create_tween()
	tween.tween_interval(0.5)
	tween.tween_property(_leaderboard_panel, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_timer_bar_bg, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_timer_bar, "modulate:a", 0.0, 0.5)

## 工具函數：三次強閃光
func _flash_triple(color: Color) -> void:
	var tween = create_tween()
	for i in range(3):
		tween.tween_property(_flash_overlay, "color:a", color.a, 0.06)
		tween.tween_property(_flash_overlay, "color:a", 0.0, 0.1)
	_flash_overlay.color = Color(color.r, color.g, color.b, 0)

## 工具函數：顯示大字
func _show_big_text(text: String, color: Color) -> void:
	var label = Label.new()
	label.text = text
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.add_theme_font_size_override("font_size", 30)
	label.add_theme_color_override("font_color", color)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.position.y -= 80
	label.modulate.a = 0.0
	add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "modulate:a", 1.0, 0.15)
	tween.tween_property(label, "position:y", label.position.y - 40, 1.2)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.4).set_delay(0.8)
	tween.tween_callback(label.queue_free)

## 工具函數：顯示結算彈窗
func _show_result_popup(winner: String, score: int, mult: float, boost_sec: int,
		leaderboard: Array, total: int) -> void:
	var popup = PanelContainer.new()
	popup.custom_minimum_size = Vector2(300, 160)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.08, 0.03, 0.0, 0.92)
	style.border_color = RACE_COLOR_GOLD
	style.set_border_width_all(2)
	style.set_corner_radius_all(8)
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	popup.add_child(vbox)

	var title_label = Label.new()
	title_label.text = "🏆 競速賽結算！"
	title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_label.add_theme_font_size_override("font_size", 18)
	title_label.add_theme_color_override("font_color", RACE_COLOR_GOLD)
	vbox.add_child(title_label)

	var winner_label = Label.new()
	winner_label.text = "🥇 冠軍：%s（%d 擊）" % [winner, score]
	winner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	winner_label.add_theme_font_size_override("font_size", 15)
	winner_label.add_theme_color_override("font_color", RACE_COLOR_GOLD)
	vbox.add_child(winner_label)

	var rank_emojis = ["🥇", "🥈", "🥉"]
	for i in range(min(leaderboard.size(), 3)):
		var entry = leaderboard[i]
		var entry_mult = entry.get("mult", 0.0)
		var entry_label = Label.new()
		entry_label.text = "%s %s ×%.1f" % [rank_emojis[i], entry.get("player_name", "—"), entry_mult]
		entry_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		entry_label.add_theme_font_size_override("font_size", 13)
		entry_label.add_theme_color_override("font_color", RACE_COLOR_CREAM)
		vbox.add_child(entry_label)

	var boost_label = Label.new()
	boost_label.text = "×%.1f 倍率加成 %ds！（共 %d 人）" % [mult, boost_sec, total]
	boost_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	boost_label.add_theme_font_size_override("font_size", 13)
	boost_label.add_theme_color_override("font_color", RACE_COLOR_ORANGE)
	vbox.add_child(boost_label)

	var vp_size = get_viewport().size
	popup.position = Vector2(vp_size.x, vp_size.y / 2 - 80)
	popup.modulate.a = 0.0
	add_child(popup)

	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 320, 0.3)
	tween.parallel().tween_property(popup, "modulate:a", 1.0, 0.3)
	tween.tween_interval(4.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.4)
	tween.tween_callback(popup.queue_free)
