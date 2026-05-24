## LuckyMirrorDuelPanel.gd — 幸運鏡像對決魚 UI 面板（DAY-270）
## 紫金對決主題：#9B59B6 紫 + #FFD700 金 + #DC143C 紅 + #1E90FF 藍
##
## 事件類型：
##   duel_start      — 對決開始（個人）
##   duel_broadcast  — 全服廣播
##   duel_score      — 積分更新（個人）
##   duel_mirror_reward — 收到鏡像分享獎勵（個人）
##   duel_result     — 對決結算（個人）
##   duel_draw       — 平局（個人）
##   duel_settle_broadcast — 全服廣播結算
##   duel_solo       — 孤獨模式（個人）

extends CanvasLayer

const COLOR_PURPLE   = Color(0.608, 0.349, 0.714)   # #9B59B6
const COLOR_GOLD     = Color(1.0,   0.843, 0.0)     # #FFD700
const COLOR_RED      = Color(0.863, 0.078, 0.235)   # #DC143C
const COLOR_BLUE     = Color(0.118, 0.565, 1.0)     # #1E90FF
const COLOR_WIN      = Color(1.0,   0.843, 0.0)     # 金色（勝利）
const COLOR_LOSE     = Color(0.7,   0.7,   0.7)     # 灰色（失敗）
const COLOR_DRAW     = Color(0.608, 0.349, 0.714)   # 紫色（平局）
const COLOR_SOLO     = Color(0.118, 0.565, 1.0)     # 藍色（孤獨）

var _banner: Control = null
var _score_panel: Control = null
var _timer_bar: ColorRect = null
var _timer_tween: Tween = null
var _my_score_label: Label = null
var _opp_score_label: Label = null
var _duel_duration: float64 = 15.0
var _is_challenger: bool = false
var _opponent_name: String = ""

func _ready() -> void:
	layer = 43  # 比 SpinWheel（42）高一層

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"duel_start":
			_on_duel_start(payload)
		"duel_broadcast":
			_on_duel_broadcast(payload)
		"duel_score":
			_on_duel_score(payload)
		"duel_mirror_reward":
			_on_duel_mirror_reward(payload)
		"duel_result":
			_on_duel_result(payload)
		"duel_draw":
			_on_duel_draw(payload)
		"duel_settle_broadcast":
			_on_duel_settle_broadcast(payload)
		"duel_solo":
			_on_duel_solo(payload)

# ── 對決開始 ──────────────────────────────────────────────────────────────────

func _on_duel_start(payload: Dictionary) -> void:
	_opponent_name = payload.get("opponent_name", "???")
	_is_challenger = payload.get("is_challenger", false)
	_duel_duration = payload.get("duration", 15.0)

	# 三次強閃光
	_flash_screen(COLOR_PURPLE, 3)

	# 頂部橫幅
	_show_banner(
		"🪞 鏡像對決！",
		"vs " + _opponent_name,
		COLOR_PURPLE
	)

	# 積分面板（右側）
	_create_score_panel()

	# 計時條
	_start_timer_bar(_duel_duration, COLOR_PURPLE)

	# 浮動大字
	_spawn_float_text(
		"🪞 對決開始！",
		Vector2(get_viewport().get_visible_rect().size / 2),
		COLOR_PURPLE,
		48
	)

# ── 全服廣播 ──────────────────────────────────────────────────────────────────

func _on_duel_broadcast(payload: Dictionary) -> void:
	var p1: String = payload.get("player_name", "???")
	var p2: String = payload.get("opponent_name", "???")
	_show_mini_banner("🪞 %s vs %s 鏡像對決！" % [p1, p2], COLOR_PURPLE)

# ── 積分更新 ──────────────────────────────────────────────────────────────────

func _on_duel_score(payload: Dictionary) -> void:
	var my_score: int = payload.get("my_score", 0)
	var opp_score: int = payload.get("opponent_score", 0)
	var share: int = payload.get("share_reward", 0)

	_update_score_panel(my_score, opp_score)

	# 輕微閃光
	_flash_screen(COLOR_BLUE, 1, 0.15)

	# 浮動文字
	_spawn_float_text(
		"🪞 +%d 鏡像分享給對手" % share,
		Vector2(get_viewport().get_visible_rect().size.x * 0.5,
				get_viewport().get_visible_rect().size.y * 0.7),
		COLOR_BLUE,
		22
	)

# ── 收到鏡像分享獎勵 ──────────────────────────────────────────────────────────

func _on_duel_mirror_reward(payload: Dictionary) -> void:
	var opp: String = payload.get("opponent_name", "???")
	var share: int = payload.get("share_reward", 0)
	var my_score: int = payload.get("my_score", 0)
	var opp_score: int = payload.get("opponent_score", 0)

	_update_score_panel(my_score, opp_score)

	# 金色閃光
	_flash_screen(COLOR_GOLD, 1, 0.12)

	# 浮動文字
	_spawn_float_text(
		"🪞 %s 分享 +%d！" % [opp, share],
		Vector2(get_viewport().get_visible_rect().size.x * 0.5,
				get_viewport().get_visible_rect().size.y * 0.65),
		COLOR_GOLD,
		24
	)

# ── 對決結算 ──────────────────────────────────────────────────────────────────

func _on_duel_result(payload: Dictionary) -> void:
	var is_winner: bool = payload.get("is_winner", false)
	var my_score: int = payload.get("my_score", 0)
	var opp_score: int = payload.get("opponent_score", 0)
	var mult: float = payload.get("result_mult", 1.0)
	var dur: float = payload.get("duration", 5.0)

	# 清除計時條和積分面板
	_clear_timer_bar()
	_clear_score_panel()
	_clear_banner()

	if is_winner:
		# 勝利：全螢幕三次強閃光 + 金色大字
		_flash_screen(COLOR_GOLD, 3)
		_spawn_float_text(
			"🏆 對決勝利！×%.1f！" % mult,
			Vector2(get_viewport().get_visible_rect().size / 2),
			COLOR_GOLD,
			52
		)
		_show_result_popup(true, my_score, opp_score, mult, dur)
	else:
		# 失敗：紅色閃光 + 灰色文字
		_flash_screen(COLOR_RED, 1, 0.2)
		_spawn_float_text(
			"💔 對決失敗 ×%.1f" % mult,
			Vector2(get_viewport().get_visible_rect().size / 2),
			COLOR_LOSE,
			40
		)
		_show_result_popup(false, my_score, opp_score, mult, dur)

# ── 平局 ──────────────────────────────────────────────────────────────────────

func _on_duel_draw(payload: Dictionary) -> void:
	var my_score: int = payload.get("my_score", 0)
	var opp_score: int = payload.get("opponent_score", 0)
	var mult: float = payload.get("result_mult", 1.5)
	var dur: float = payload.get("duration", 5.0)

	_clear_timer_bar()
	_clear_score_panel()
	_clear_banner()

	# 紫色閃光 + 平局大字
	_flash_screen(COLOR_PURPLE, 2)
	_spawn_float_text(
		"🤝 平局！×%.1f！" % mult,
		Vector2(get_viewport().get_visible_rect().size / 2),
		COLOR_DRAW,
		46
	)
	_show_result_popup(true, my_score, opp_score, mult, dur)

# ── 全服廣播結算 ──────────────────────────────────────────────────────────────

func _on_duel_settle_broadcast(payload: Dictionary) -> void:
	var p1: String = payload.get("player_name", "???")
	var p2: String = payload.get("opponent_name", "???")
	var winner: String = payload.get("winner_name", "")
	var win_mult: float = payload.get("win_mult", 1.5)
	var s1: int = payload.get("score1", 0)
	var s2: int = payload.get("score2", 0)

	var msg: String
	if winner == "":
		msg = "🪞 %s vs %s 平局！（%d:%d）" % [p1, p2, s1, s2]
	else:
		msg = "🏆 %s 對決勝利！（%d:%d）×%.1f！" % [winner, s1, s2, win_mult]

	_show_mini_banner(msg, COLOR_GOLD if winner != "" else COLOR_DRAW)

# ── 孤獨模式 ──────────────────────────────────────────────────────────────────

func _on_duel_solo(payload: Dictionary) -> void:
	var mult: float = payload.get("solo_mult", 1.5)
	var dur: float = payload.get("duration", 10.0)

	# 藍色閃光
	_flash_screen(COLOR_SOLO, 2)

	_spawn_float_text(
		"🌊 孤獨模式 ×%.1f（%ds）" % [mult, int(dur)],
		Vector2(get_viewport().get_visible_rect().size / 2),
		COLOR_SOLO,
		40
	)

	# 計時條（藍色）
	_start_timer_bar(dur, COLOR_SOLO)

	# 超時後清除
	var t := get_tree().create_timer(dur)
	t.timeout.connect(func(): _clear_timer_bar())

# ── 內部 UI 工具 ──────────────────────────────────────────────────────────────

func _flash_screen(color: Color, times: int, alpha: float = 0.45) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var flash := ColorRect.new()
	flash.size = vp_size
	flash.color = Color(color.r, color.g, color.b, 0.0)
	add_child(flash)

	var tween := create_tween()
	for i in range(times):
		tween.tween_property(flash, "color:a", alpha, 0.07)
		tween.tween_property(flash, "color:a", 0.0, 0.1)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)

func _show_banner(title: String, subtitle: String, color: Color) -> void:
	_clear_banner()
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(0, -80)
	panel.size = Vector2(vp_size.x, 72)
	panel.modulate = Color(color.r, color.g, color.b, 0.92)
	add_child(panel)
	_banner = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var lbl1 := Label.new()
	lbl1.text = title
	lbl1.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl1.add_theme_font_size_override("font_size", 28)
	lbl1.add_theme_color_override("font_color", Color.WHITE)
	vbox.add_child(lbl1)

	var lbl2 := Label.new()
	lbl2.text = subtitle
	lbl2.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl2.add_theme_font_size_override("font_size", 18)
	lbl2.add_theme_color_override("font_color", Color(1, 1, 0.8))
	vbox.add_child(lbl2)

	# 滑入動畫
	var tween := create_tween()
	tween.tween_property(panel, "position:y", 0.0, 0.25)

func _show_mini_banner(text: String, color: Color) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var lbl := Label.new()
	lbl.text = text
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 18)
	lbl.add_theme_color_override("font_color", color)
	lbl.position = Vector2(0, 8)
	lbl.size = Vector2(vp_size.x, 28)
	add_child(lbl)

	var tween := create_tween()
	tween.tween_interval(3.0)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)

func _clear_banner() -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()
	_banner = null

func _create_score_panel() -> void:
	_clear_score_panel()
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x - 160, 80)
	panel.size = Vector2(150, 90)
	panel.modulate = Color(0.608, 0.349, 0.714, 0.9)
	add_child(panel)
	_score_panel = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "🪞 對決積分"
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	vbox.add_child(title_lbl)

	_my_score_label = Label.new()
	_my_score_label.text = "我：0"
	_my_score_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_my_score_label.add_theme_font_size_override("font_size", 20)
	_my_score_label.add_theme_color_override("font_color", Color.WHITE)
	vbox.add_child(_my_score_label)

	_opp_score_label = Label.new()
	_opp_score_label.text = "對手：0"
	_opp_score_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_opp_score_label.add_theme_font_size_override("font_size", 20)
	_opp_score_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	vbox.add_child(_opp_score_label)

func _update_score_panel(my_score: int, opp_score: int) -> void:
	if is_instance_valid(_my_score_label):
		_my_score_label.text = "我：%d" % my_score
		_my_score_label.add_theme_color_override("font_color",
			COLOR_GOLD if my_score >= opp_score else Color.WHITE)
	if is_instance_valid(_opp_score_label):
		_opp_score_label.text = "對手：%d" % opp_score
		_opp_score_label.add_theme_color_override("font_color",
			COLOR_RED if opp_score > my_score else Color(0.8, 0.8, 0.8))

func _clear_score_panel() -> void:
	if is_instance_valid(_score_panel):
		_score_panel.queue_free()
	_score_panel = null
	_my_score_label = null
	_opp_score_label = null

func _start_timer_bar(duration: float, color: Color) -> void:
	_clear_timer_bar()
	var vp_size := get_viewport().get_visible_rect().size

	var bar := ColorRect.new()
	bar.size = Vector2(8, vp_size.y * 0.6)
	bar.position = Vector2(vp_size.x - 310, vp_size.y * 0.2)  # x=-310 與其他計時條錯開
	bar.color = color
	add_child(bar)
	_timer_bar = bar

	_timer_tween = create_tween()
	_timer_tween.tween_property(bar, "size:y", 0.0, duration)
	_timer_tween.tween_callback(func():
		if is_instance_valid(bar):
			bar.queue_free()
		_timer_bar = null
	)

func _clear_timer_bar() -> void:
	if _timer_tween != null:
		_timer_tween.kill()
		_timer_tween = null
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
	_timer_bar = null

func _show_result_popup(is_win: bool, my_score: int, opp_score: int, mult: float, dur: float) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var popup := PanelContainer.new()
	popup.size = Vector2(260, 140)
	popup.position = Vector2(vp_size.x + 10, vp_size.y * 0.35)
	popup.modulate = Color(0.1, 0.1, 0.15, 0.95)
	add_child(popup)

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	popup.add_child(vbox)

	var result_color := COLOR_WIN if is_win else COLOR_LOSE
	var result_text := "🏆 勝利！" if is_win else "💔 失敗"

	var lbl_result := Label.new()
	lbl_result.text = result_text
	lbl_result.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_result.add_theme_font_size_override("font_size", 24)
	lbl_result.add_theme_color_override("font_color", result_color)
	vbox.add_child(lbl_result)

	var lbl_score := Label.new()
	lbl_score.text = "比分：%d : %d" % [my_score, opp_score]
	lbl_score.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_score.add_theme_font_size_override("font_size", 18)
	lbl_score.add_theme_color_override("font_color", Color.WHITE)
	vbox.add_child(lbl_score)

	var lbl_mult := Label.new()
	lbl_mult.text = "×%.1f 加成 %ds" % [mult, int(dur)]
	lbl_mult.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_mult.add_theme_font_size_override("font_size", 20)
	lbl_mult.add_theme_color_override("font_color", result_color)
	vbox.add_child(lbl_mult)

	# 右側滑入動畫
	var tween := create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 270.0, 0.3)
	tween.tween_interval(3.5)
	tween.tween_property(popup, "position:x", vp_size.x + 10.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(popup):
			popup.queue_free()
	)

func _spawn_float_text(text: String, pos: Vector2, color: Color, font_size: int = 28) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.position = pos - Vector2(200, 30)
	lbl.size = Vector2(400, 60)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", font_size)
	lbl.add_theme_color_override("font_color", color)
	add_child(lbl)

	var tween := create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 60.0, 1.2)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 1.2)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)
