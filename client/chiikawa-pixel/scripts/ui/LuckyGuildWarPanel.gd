## LuckyGuildWarPanel.gd — 幸運公會戰魚系統面板（DAY-257）
## 業界原創「全服分隊競爭→勝隊爆發」機制
##
## 視覺設計：
##   - 紅藍公會戰主題（#DC143C 紅 + #1E90FF 藍 + #FFD700 金 + #9370DB 平局紫）
##   - war_start：紅色三次強閃光 + 頂部橫幅 + 「⚔️ 公會戰！」大字 + 隊伍指示器 + 計時條
##   - war_broadcast：頂部小橫幅（全服廣播）
##   - war_score：右側即時比分面板（每 5 秒更新）
##   - war_result：全螢幕三次強閃光 + 「⚔️ 勝隊！×2.5」大字 + 結算彈窗
##   - war_draw：紫色閃光 + 「⚔️ 平局！×1.8」大字 + 結算彈窗
extends CanvasLayer

# 主題顏色
const COLOR_RED    = Color("#DC143C")  # 紅隊（主色）
const COLOR_BLUE   = Color("#1E90FF")  # 藍隊
const COLOR_GOLD   = Color("#FFD700")  # 金色（勝利）
const COLOR_PURPLE = Color("#9370DB")  # 平局紫
const COLOR_CREAM  = Color("#FFF3E0")  # 奶油（副文字）
const COLOR_WHITE  = Color("#FFFFFF")  # 白色
const COLOR_GRAY   = Color("#808080")  # 灰色

# 計時條
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null

# 比分面板
var _score_panel: PanelContainer = null
var _red_score_lbl: Label = null
var _blue_score_lbl: Label = null

# 隊伍指示器
var _team_indicator: PanelContainer = null

func _ready() -> void:
	layer = 30  # 幸運公會戰魚面板層級（DAY-257）

## 處理幸運公會戰魚訊息
func handle_lucky_guild_war(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"war_start":
			_on_war_start(payload)
		"war_broadcast":
			_on_war_broadcast(payload)
		"war_score":
			_on_war_score(payload)
		"war_result":
			_on_war_result(payload)
		"war_draw":
			_on_war_draw(payload)

## war_start — 公會戰啟動（個人訊息）
func _on_war_start(payload: Dictionary) -> void:
	var duration_sec: int = payload.get("duration_sec", 30)
	var player_team: int = payload.get("player_team", 0)  # 0=紅隊, 1=藍隊
	var win_mult: float = payload.get("win_mult", 2.5)
	var lose_mult: float = payload.get("lose_mult", 1.2)

	var my_team_color = COLOR_RED if player_team == 0 else COLOR_BLUE
	var my_team_name = "🔴 紅隊" if player_team == 0 else "🔵 藍隊"

	# 紅色三次強閃光
	_flash_screen(my_team_color, 0.65, 3)

	# 頂部橫幅
	_show_banner("⚔️ 公會戰開始！你是 %s！勝隊 ×%.1f 倍率加成！" % [my_team_name, win_mult], my_team_color, 4.0)

	# 中央大字
	_show_big_text("⚔️ 公會戰！", my_team_color, 52, 2.5)
	_show_sub_text("你是 %s！30秒競爭擊破數！勝隊 ×%.1f，敗隊 ×%.1f！" % [my_team_name, win_mult, lose_mult], COLOR_CREAM, 2.5)

	# 隊伍指示器（右上角）
	_show_team_indicator(player_team)

	# 比分面板（右側）
	_show_score_panel(0, 0)

	# 右側豎向計時條（x=-184 與其他計時條錯開）
	_start_timer_bar(duration_sec, my_team_color)

## war_broadcast — 全服廣播公會戰
func _on_war_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var win_mult: float = payload.get("win_mult", 2.5)
	var duration_sec: int = payload.get("duration_sec", 30)
	_show_top_banner("⚔️ %s 發起公會戰！🔴紅隊 vs 🔵藍隊，30秒競爭！勝隊 ×%.1f！" % [player_name, win_mult], COLOR_RED, 3.0)
	# 廣播時也顯示比分面板
	_show_score_panel(0, 0)
	_start_timer_bar(duration_sec, COLOR_RED)

## war_score — 即時比分廣播（每 5 秒）
func _on_war_score(payload: Dictionary) -> void:
	var red_score: int = payload.get("red_score", 0)
	var blue_score: int = payload.get("blue_score", 0)
	_update_score_panel(red_score, blue_score)

## war_result — 公會戰結算（有勝負）
func _on_war_result(payload: Dictionary) -> void:
	var red_score: int = payload.get("red_score", 0)
	var blue_score: int = payload.get("blue_score", 0)
	var win_team: int = payload.get("win_team", 0)  # 0=紅隊, 1=藍隊
	var win_mult: float = payload.get("win_mult", 2.5)
	var lose_mult: float = payload.get("lose_mult", 1.2)
	var boost_sec: int = payload.get("boost_sec", 5)
	var player_name: String = payload.get("player_name", "某玩家")

	# 停止計時條和比分面板
	_stop_timer_bar()
	_clear_score_panel()
	_clear_team_indicator()

	var win_color = COLOR_RED if win_team == 0 else COLOR_BLUE
	var win_name = "🔴 紅隊" if win_team == 0 else "🔵 藍隊"
	var win_score = red_score if win_team == 0 else blue_score
	var lose_score = blue_score if win_team == 0 else red_score

	# 全螢幕三次強閃光
	_flash_screen(win_color, 0.75, 3)

	# 大字
	_show_big_text("⚔️ %s 勝利！" % win_name, win_color, 52, 2.5)
	_show_sub_text("%d vs %d！勝隊 ×%.1f 倍率加成 %ds！" % [win_score, lose_score, win_mult, boost_sec], COLOR_CREAM, 2.5)

	# 結算彈窗
	_show_result_popup(player_name, red_score, blue_score, win_team, win_mult, lose_mult, boost_sec, false)

## war_draw — 平局結算
func _on_war_draw(payload: Dictionary) -> void:
	var red_score: int = payload.get("red_score", 0)
	var blue_score: int = payload.get("blue_score", 0)
	var win_mult: float = payload.get("win_mult", 1.8)
	var boost_sec: int = payload.get("boost_sec", 5)
	var player_name: String = payload.get("player_name", "某玩家")

	# 停止計時條和比分面板
	_stop_timer_bar()
	_clear_score_panel()
	_clear_team_indicator()

	# 紫色閃光（平局）
	_flash_screen(COLOR_PURPLE, 0.5, 2)

	# 大字
	_show_big_text("⚔️ 勢均力敵！平局！", COLOR_PURPLE, 48, 2.5)
	_show_sub_text("%d vs %d！雙隊 ×%.1f 倍率加成 %ds！" % [red_score, blue_score, win_mult, boost_sec], COLOR_CREAM, 2.5)

	# 結算彈窗
	_show_result_popup(player_name, red_score, blue_score, -1, win_mult, win_mult, boost_sec, true)

# ─── 比分面板 ─────────────────────────────────────────────────────────────────

func _show_score_panel(red: int, blue: int) -> void:
	_clear_score_panel()

	var vp_size = get_viewport().size
	_score_panel = PanelContainer.new()
	_score_panel.size = Vector2(140, 100)
	_score_panel.position = Vector2(vp_size.x - 155, vp_size.y * 0.35)
	add_child(_score_panel)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.05, 0.1, 0.88)
	style.border_color = COLOR_GOLD
	style.set_border_width_all(2)
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	_score_panel.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 4)
	_score_panel.add_child(vbox)

	var title_lbl = Label.new()
	title_lbl.text = "⚔️ 公會戰"
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	_red_score_lbl = Label.new()
	_red_score_lbl.text = "🔴 紅隊：%d" % red
	_red_score_lbl.add_theme_font_size_override("font_size", 16)
	_red_score_lbl.add_theme_color_override("font_color", COLOR_RED)
	_red_score_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(_red_score_lbl)

	_blue_score_lbl = Label.new()
	_blue_score_lbl.text = "🔵 藍隊：%d" % blue
	_blue_score_lbl.add_theme_font_size_override("font_size", 16)
	_blue_score_lbl.add_theme_color_override("font_color", COLOR_BLUE)
	_blue_score_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(_blue_score_lbl)

func _update_score_panel(red: int, blue: int) -> void:
	if not is_instance_valid(_score_panel):
		_show_score_panel(red, blue)
		return
	if is_instance_valid(_red_score_lbl):
		_red_score_lbl.text = "🔴 紅隊：%d" % red
		# 領先時放大
		if red > blue:
			var tween = _red_score_lbl.create_tween()
			tween.tween_property(_red_score_lbl, "scale", Vector2(1.2, 1.2), 0.1)
			tween.tween_property(_red_score_lbl, "scale", Vector2(1.0, 1.0), 0.1)
	if is_instance_valid(_blue_score_lbl):
		_blue_score_lbl.text = "🔵 藍隊：%d" % blue
		if blue > red:
			var tween = _blue_score_lbl.create_tween()
			tween.tween_property(_blue_score_lbl, "scale", Vector2(1.2, 1.2), 0.1)
			tween.tween_property(_blue_score_lbl, "scale", Vector2(1.0, 1.0), 0.1)

func _clear_score_panel() -> void:
	if is_instance_valid(_score_panel):
		_score_panel.queue_free()
		_score_panel = null
	_red_score_lbl = null
	_blue_score_lbl = null

# ─── 隊伍指示器 ───────────────────────────────────────────────────────────────

func _show_team_indicator(player_team: int) -> void:
	_clear_team_indicator()

	var vp_size = get_viewport().size
	_team_indicator = PanelContainer.new()
	_team_indicator.size = Vector2(120, 40)
	_team_indicator.position = Vector2(vp_size.x - 130, 10)
	add_child(_team_indicator)

	var style = StyleBoxFlat.new()
	style.bg_color = COLOR_RED if player_team == 0 else COLOR_BLUE
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	_team_indicator.add_theme_stylebox_override("panel", style)

	var lbl = Label.new()
	lbl.text = "🔴 紅隊" if player_team == 0 else "🔵 藍隊"
	lbl.add_theme_font_size_override("font_size", 14)
	lbl.add_theme_color_override("font_color", COLOR_WHITE)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_team_indicator.add_child(lbl)

	# 脈衝動畫
	var tween = _team_indicator.create_tween().set_loops()
	tween.tween_property(_team_indicator, "modulate:a", 0.7, 0.6)
	tween.tween_property(_team_indicator, "modulate:a", 1.0, 0.6)

func _clear_team_indicator() -> void:
	if is_instance_valid(_team_indicator):
		_team_indicator.queue_free()
		_team_indicator = null

# ─── 結算彈窗 ─────────────────────────────────────────────────────────────────

func _show_result_popup(trigger_name: String, red_score: int, blue_score: int,
		win_team: int, win_mult: float, lose_mult: float, boost_sec: int, is_draw: bool) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(360, 220)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 110)
	add_child(popup)

	var border_color = COLOR_PURPLE if is_draw else (COLOR_RED if win_team == 0 else COLOR_BLUE)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.05, 0.1, 0.93)
	style.border_color = border_color
	style.set_border_width_all(3)
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 8)
	popup.add_child(vbox)

	# 標題
	var title_text = "⚔️ 公會戰平局！" if is_draw else ("⚔️ 🔴紅隊勝利！" if win_team == 0 else "⚔️ 🔵藍隊勝利！")
	var title_color = COLOR_PURPLE if is_draw else (COLOR_RED if win_team == 0 else COLOR_BLUE)
	var title_lbl = Label.new()
	title_lbl.text = title_text
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", title_color)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	# 發起者
	var trigger_lbl = Label.new()
	trigger_lbl.text = "戰爭發起者：%s" % trigger_name
	trigger_lbl.add_theme_font_size_override("font_size", 12)
	trigger_lbl.add_theme_color_override("font_color", COLOR_CREAM)
	trigger_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(trigger_lbl)

	# 比分
	var score_lbl = Label.new()
	score_lbl.text = "🔴 紅隊 %d  vs  🔵 藍隊 %d" % [red_score, blue_score]
	score_lbl.add_theme_font_size_override("font_size", 16)
	score_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	score_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(score_lbl)

	# 勝隊倍率
	if is_draw:
		var mult_lbl = Label.new()
		mult_lbl.text = "雙隊 ×%.1f 倍率加成 %ds！" % [win_mult, boost_sec]
		mult_lbl.add_theme_font_size_override("font_size", 18)
		mult_lbl.add_theme_color_override("font_color", COLOR_PURPLE)
		mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		vbox.add_child(mult_lbl)
	else:
		var win_lbl = Label.new()
		win_lbl.text = "勝隊 ×%.1f 倍率加成 %ds！" % [win_mult, boost_sec]
		win_lbl.add_theme_font_size_override("font_size", 18)
		win_lbl.add_theme_color_override("font_color", COLOR_GOLD)
		win_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		vbox.add_child(win_lbl)

		var lose_lbl = Label.new()
		lose_lbl.text = "敗隊 ×%.1f 安慰獎" % lose_mult
		lose_lbl.add_theme_font_size_override("font_size", 14)
		lose_lbl.add_theme_color_override("font_color", COLOR_GRAY)
		lose_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		vbox.add_child(lose_lbl)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 380.0, 0.4).set_ease(Tween.EASE_OUT)
	tween.tween_interval(6.0)
	tween.tween_property(popup, "position:x", vp_size.x + 10.0, 0.4).set_ease(Tween.EASE_IN)
	tween.tween_callback(popup.queue_free)

# ─── 通用 UI 工具 ─────────────────────────────────────────────────────────────

func _flash_screen(color: Color, alpha: float, times: int) -> void:
	var vp_size = get_viewport().size
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.0)
	flash.size = vp_size
	add_child(flash)
	var tween = create_tween()
	for i in range(times):
		tween.tween_property(flash, "color:a", alpha, 0.08)
		tween.tween_property(flash, "color:a", 0.0, 0.12)
	tween.tween_callback(flash.queue_free)

func _show_banner(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var banner = PanelContainer.new()
	banner.size = Vector2(vp_size.x, 48)
	banner.position = Vector2(0, 0)
	add_child(banner)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(color.r, color.g, color.b, 0.88)
	banner.add_theme_stylebox_override("panel", style)
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 16)
	lbl.add_theme_color_override("font_color", COLOR_WHITE)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	banner.add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(banner.queue_free)

func _show_top_banner(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var banner = PanelContainer.new()
	banner.size = Vector2(vp_size.x * 0.7, 36)
	banner.position = Vector2(vp_size.x * 0.15, 52)
	add_child(banner)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(color.r, color.g, color.b, 0.82)
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	banner.add_theme_stylebox_override("panel", style)
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 13)
	lbl.add_theme_color_override("font_color", COLOR_WHITE)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	banner.add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(banner.queue_free)

func _show_big_text(text: String, color: Color, font_size: int, duration: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", font_size)
	lbl.add_theme_color_override("font_color", color)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.size = Vector2(vp_size.x, 80)
	lbl.position = Vector2(0, vp_size.y * 0.35)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "scale", Vector2(1.15, 1.15), 0.15).set_ease(Tween.EASE_OUT)
	tween.tween_property(lbl, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(duration - 0.5)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.5)
	tween.tween_callback(lbl.queue_free)

func _show_sub_text(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 13)
	lbl.add_theme_color_override("font_color", color)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.size = Vector2(vp_size.x, 40)
	lbl.position = Vector2(0, vp_size.y * 0.35 + 80)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.4)
	tween.tween_callback(lbl.queue_free)

func _start_timer_bar(duration_sec: int, color: Color) -> void:
	_stop_timer_bar()

	var vp_size = get_viewport().size
	var bar_h: float = vp_size.y * 0.5
	var bar_w: float = 12.0
	# x=-184 與其他計時條錯開（龍捲風-58，黑洞-72，鏡像分裂-86，量子糾纏-100，武器進化-114，隕石雨-128，龍王-142，時空裂縫-156，全服充能-170，公會戰-184）
	var bar_x: float = vp_size.x - 184.0
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.05, 0.05, 0.1, 0.7)
	_timer_bar_bg.size = Vector2(bar_w, bar_h)
	_timer_bar_bg.position = Vector2(bar_x, bar_y)
	add_child(_timer_bar_bg)

	_timer_bar = ColorRect.new()
	_timer_bar.color = color
	_timer_bar.size = Vector2(bar_w, bar_h)
	_timer_bar.position = Vector2(bar_x, bar_y)
	add_child(_timer_bar)

	_timer_tween = create_tween()
	_timer_tween.tween_property(_timer_bar, "size:y", 0.0, float(duration_sec)).set_ease(Tween.EASE_IN_OUT)

func _stop_timer_bar() -> void:
	if is_instance_valid(_timer_tween):
		_timer_tween.kill()
	if is_instance_valid(_timer_bar):
		var tween = create_tween()
		tween.tween_property(_timer_bar, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_timer_bar.queue_free)
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		var tween2 = create_tween()
		tween2.tween_property(_timer_bar_bg, "modulate:a", 0.0, 0.5)
		tween2.tween_callback(_timer_bar_bg.queue_free)
		_timer_bar_bg = null
