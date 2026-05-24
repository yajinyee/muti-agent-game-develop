## LuckyFlagFishPanel.gd ??撟賊?憟芣?擳頂蝯梢?選?DAY-244嚗?
## 璆剔?????奎?准???
##
## 閬死閮剛?嚗?
##   - 蝝?嗆?銝駁?嚗?E74C3C + #FFFFFF + #F39C12 + #2ECC71嚗?
##   - flag_start嚗??脖?甈∪撥?? + ?璈怠? + ????冽??嗆?嚗之摮?+ ?格?璅? + 閮?璇?+ ??隤芣?
##   - flag_rank_update嚗?湔???選??單??湔嚗? 3 蝘?
##   - flag_captured嚗??脖?甈∪撥?? + ???憟芣???嚗?.0?之摮?+ 蝯?敶?
##   - flag_timeout嚗??脤???+ ????嗆?蝯?嚗?蝷?+ ??蝯?
##   - flag_escaped嚗?脫?蝷?
##   - flag_auto_blast嚗??脤???+ ????芸??嚗?蝷?
extends CanvasLayer

# 銝駁?憿
const COLOR_FLAG    = Color("#E74C3C")  # 蝝嚗蜓憿?
const COLOR_GOLD    = Color("#F39C12")  # ?嚗洵 1 ??
const COLOR_SILVER  = Color("#BDC3C7")  # ??莎?蝚?2 ??
const COLOR_BRONZE  = Color("#E67E22")  # ?嚗洵 3 ??
const COLOR_GREEN   = Color("#2ECC71")  # 蝬嚗???
const COLOR_FAIL    = Color("#7F8C8D")  # ?啗嚗仃??
const COLOR_WHITE   = Color("#FFFFFF")  # ?質

# 閮?璇?暺?
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null
var _duration_sec: int = 15

# ???Ｘ蝭暺?
var _rank_panel: Control = null
var _rank_labels: Array = []

# ?格?璅?蝭暺?
var _target_marker: Control = null

func _ready() -> void:
	layer = 1  # 撟賊?憟芣?擳?踹惜蝝?

## ??撟賊?憟芣?擳???
func handle_lucky_flag_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"flag_start":
			_on_flag_start(payload)
		"flag_rank_update":
			_on_flag_rank_update(payload)
		"flag_captured":
			_on_flag_captured(payload)
		"flag_timeout":
			_on_flag_timeout(payload)
		"flag_escaped":
			_on_flag_escaped(payload)
		"flag_auto_blast":
			_on_flag_auto_blast(payload)

## flag_start ???嗆???嚗?誨?哨?
func _on_flag_start(payload: Dictionary) -> void:
	_duration_sec = payload.get("duration_sec", 15)
	var player_name: String = payload.get("player_name", "Player")
	var winner_mult: float = payload.get("winner_mult", 4.0)
	var second_mult: float = payload.get("second_mult", 2.0)
	var third_mult: float = payload.get("third_mult", 1.5)
	var target_x: float = payload.get("x", 0.0)
	var target_y: float = payload.get("y", 0.0)
	var vp_size = get_viewport().size

	# 蝝銝活撘琿???
	_flash_screen(COLOR_FLAG, 0.55, 3)

	# ?璈怠?
	_show_banner("? ?冽??嗆?嚗???撟璅敞蝛???", COLOR_FLAG, 4.0)

	# 銝剖亢憭批?
	_show_big_text("? ?冽??嗆?嚗?, COLOR_FLAG, 52, 2.5)"

	# ??隤芣?
	_show_sub_text("蝚????%.1f  蝚????%.1f  蝚????%.1f" % [winner_mult, second_mult, third_mult], COLOR_GOLD, 2.0)

	# 閫貊??蝷?
	_show_top_banner("? %s 閫貊?嗆?嚗?d 蝘蝛??擃敺??%.1f ??嚗? % [player_name, _duration_sec, winner_mult], COLOR_FLAG, 3.5)"

	# ?喳鞊?閮?璇?
	_create_timer_bar(_duration_sec)

	# ?格?璅?
	_create_target_marker(target_x, target_y)

	# 撱箇????Ｘ
	_create_rank_panel()

## flag_rank_update ???單????湔嚗?誨?哨?
func _on_flag_rank_update(payload: Dictionary) -> void:
	var rank_list: Array = payload.get("rank_list", [])
	var remaining_sec: int = payload.get("remaining_sec", 0)
	_update_rank_panel(rank_list, remaining_sec)

## flag_captured ?????格?鋡急??湛??冽?撱?嚗?
func _on_flag_captured(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "?摰?)"
	var killer_mult: float = payload.get("killer_mult", 4.0)
	var killer_rank: int = payload.get("killer_rank", 1)
	var reward: int = payload.get("reward", 0)
	var rank_list: Array = payload.get("rank_list", [])

	# 皜閮?璇??格?璅?
	_clear_timer_bar()
	_clear_target_marker()
	_clear_rank_panel()

	# ?寞????豢?憿
	var color = COLOR_GOLD
	if killer_rank == 2:
		color = COLOR_SILVER
	elif killer_rank == 3:
		color = COLOR_BRONZE

	# ??
	_flash_screen(color, 0.6, 3)

	# 銝剖亢憭批?
	var rank_text = ["蝚???, "蝚???, "蝚???].get(killer_rank - 1, "")"
	_show_big_text("? %s 憟芣?嚗?s ?%.1f" % [player_name, rank_text, killer_mult], color, 3.0)

	# 蝯?敶?
	_show_result_popup(rank_list, reward)

## flag_timeout ???嗆????堆??冽?撱?嚗?
func _on_flag_timeout(payload: Dictionary) -> void:
	var rank_list: Array = payload.get("rank_list", [])

	# 皜閮?璇??格?璅?
	_clear_timer_bar()
	_clear_target_marker()
	_clear_rank_panel()

	# 璈??
	_flash_screen(COLOR_BRONZE, 0.4, 2)

	# ?內??
	_show_big_text("? ?嗆?蝯?嚗?, COLOR_BRONZE, 2.5)"

	# 蝯?敶?
	_show_result_popup(rank_list, 0)

## flag_escaped ?????格?瘨仃嚗?誨?哨?
func _on_flag_escaped(payload: Dictionary) -> void:
	_clear_timer_bar()
	_clear_target_marker()
	_clear_rank_panel()
	_show_big_text("? ???格???鈭?", COLOR_FAIL, 2.0)

## flag_auto_blast ?????格??芸??嚗?誨?哨?
func _on_flag_auto_blast(payload: Dictionary) -> void:
	var reward: int = payload.get("reward", 0)

	_clear_timer_bar()
	_clear_target_marker()
	_clear_rank_panel()

	# 璈??
	_flash_screen(COLOR_BRONZE, 0.5, 2)

	# ?內??
	_show_big_text("? ???芸??嚗??+%d ?馳嚗? % reward, COLOR_BRONZE, 2.5)"

# ??? ?折 UI 撌亙?賣 ???????????????????????????????????????????????????????

## ?刻撟???
func _flash_screen(color: Color, alpha: float, times: int) -> void:
	var vp_size = get_viewport().size
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, alpha)
	flash.size = vp_size
	flash.position = Vector2.ZERO
	add_child(flash)

	var tween = create_tween()
	for i in range(times):
		tween.tween_property(flash, "modulate:a", 1.0, 0.05)
		tween.tween_property(flash, "modulate:a", 0.0, 0.1)
	tween.tween_callback(flash.queue_free)

## ?璈怠?
func _show_banner(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var banner = ColorRect.new()
	banner.color = Color(color.r, color.g, color.b, 0.85)
	banner.size = Vector2(vp_size.x, 48)
	banner.position = Vector2(0, 0)
	add_child(banner)

	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 20)
	label.add_theme_color_override("font_color", Color.WHITE)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	label.size = banner.size
	label.position = Vector2.ZERO
	banner.add_child(label)

	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(banner.queue_free)

## ?撠帖撟?
func _show_top_banner(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var banner = ColorRect.new()
	banner.color = Color(color.r, color.g, color.b, 0.75)
	banner.size = Vector2(vp_size.x * 0.7, 36)
	banner.position = Vector2(vp_size.x * 0.15, 56)
	add_child(banner)

	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 16)
	label.add_theme_color_override("font_color", Color.WHITE)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	label.size = banner.size
	label.position = Vector2.ZERO
	banner.add_child(label)

	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(banner.queue_free)

## 銝剖亢憭批?
func _show_big_text(text: String, color: Color, font_size: int, duration: float) -> void:
	var vp_size = get_viewport().size
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", font_size)
	label.add_theme_color_override("font_color", color)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.size = Vector2(vp_size.x, 80)
	label.position = Vector2(0, vp_size.y * 0.35)
	add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "position:y", vp_size.y * 0.30, 0.3)
	tween.tween_interval(duration - 0.8)
	tween.tween_property(label, "modulate:a", 0.0, 0.5)
	tween.tween_callback(label.queue_free)

## ?舀?憿?摮?
func _show_sub_text(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 20)
	label.add_theme_color_override("font_color", color)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.size = Vector2(vp_size.x, 40)
	label.position = Vector2(0, vp_size.y * 0.45)
	add_child(label)

	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(label, "modulate:a", 0.0, 0.4)
	tween.tween_callback(label.queue_free)

## 撱箇??喳鞊?閮?璇?
func _create_timer_bar(duration: int) -> void:
	_clear_timer_bar()
	var vp_size = get_viewport().size

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.3, 0.05, 0.05, 0.6)
	_timer_bar_bg.size = Vector2(16, 200)
	_timer_bar_bg.position = Vector2(vp_size.x - 28, vp_size.y * 0.3)
	add_child(_timer_bar_bg)

	_timer_bar = ColorRect.new()
	_timer_bar.color = COLOR_FLAG
	_timer_bar.size = Vector2(16, 200)
	_timer_bar.position = Vector2(vp_size.x - 28, vp_size.y * 0.3)
	add_child(_timer_bar)

	_timer_tween = create_tween()
	_timer_tween.tween_property(_timer_bar, "size:y", 0.0, float(duration))

## 皜閮?璇?
func _clear_timer_bar() -> void:
	if is_instance_valid(_timer_tween):
		_timer_tween.kill()
		_timer_tween = null
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		_timer_bar_bg.queue_free()
		_timer_bar_bg = null

## 撱箇??格?璅?
func _create_target_marker(target_x: float, target_y: float) -> void:
	_clear_target_marker()
	if target_x <= 0 and target_y <= 0:
		return

	_target_marker = Control.new()
	_target_marker.position = Vector2(target_x - 20, target_y - 55)
	_target_marker.size = Vector2(40, 40)
	add_child(_target_marker)

	var label = Label.new()
	label.text = "?"
	label.add_theme_font_size_override("font_size", 30)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.size = Vector2(40, 40)
	label.position = Vector2.ZERO
	_target_marker.add_child(label)

	# 銝?瘚桀??
	var tween = _target_marker.create_tween().set_loops()
	tween.tween_property(label, "position:y", -5.0, 0.5)
	tween.tween_property(label, "position:y", 0.0, 0.5)

## 皜?格?璅?
func _clear_target_marker() -> void:
	if is_instance_valid(_target_marker):
		_target_marker.queue_free()
		_target_marker = null

## 撱箇????Ｘ嚗?湛?
func _create_rank_panel() -> void:
	_clear_rank_panel()
	var vp_size = get_viewport().size

	_rank_panel = Control.new()
	_rank_panel.position = Vector2(vp_size.x - 180, vp_size.y * 0.15)
	_rank_panel.size = Vector2(170, 200)
	add_child(_rank_panel)

	# ?
	var bg = ColorRect.new()
	bg.color = Color(0.1, 0.0, 0.0, 0.75)
	bg.size = Vector2(170, 200)
	bg.position = Vector2.ZERO
	_rank_panel.add_child(bg)

	# 璅?
	var title = Label.new()
	title.text = "? ?嗆???"
	title.add_theme_font_size_override("font_size", 16)
	title.add_theme_color_override("font_color", COLOR_FLAG)
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.size = Vector2(170, 28)
	title.position = Vector2(0, 4)
	_rank_panel.add_child(title)

	# ?????憭?5 ??
	_rank_labels = []
	for i in range(5):
		var row = Label.new()
		row.text = ""
		row.add_theme_font_size_override("font_size", 14)
		row.add_theme_color_override("font_color", COLOR_WHITE)
		row.size = Vector2(170, 28)
		row.position = Vector2(4, 32 + i * 30)
		_rank_panel.add_child(row)
		_rank_labels.append(row)

## ?湔???Ｘ
func _update_rank_panel(rank_list: Array, remaining_sec: int) -> void:
	if not is_instance_valid(_rank_panel):
		return

	for i in range(_rank_labels.size()):
		var label = _rank_labels[i]
		if not is_instance_valid(label):
			continue
		if i < rank_list.size():
			var entry = rank_list[i]
			var rank = entry.get("rank", i + 1)
			var name_str = entry.get("player_name", "?")
			var score = entry.get("score", 0)
			var color = COLOR_WHITE
			if rank == 1:
				color = COLOR_GOLD
			elif rank == 2:
				color = COLOR_SILVER
			elif rank == 3:
				color = COLOR_BRONZE
			label.text = "#%d %s (%d)" % [rank, name_str, score]
			label.add_theme_color_override("font_color", color)
		else:
			label.text = ""

## 皜???Ｘ
func _clear_rank_panel() -> void:
	if is_instance_valid(_rank_panel):
		_rank_panel.queue_free()
		_rank_panel = null
	_rank_labels = []

## 蝯?敶?
func _show_result_popup(rank_list: Array, reward: int) -> void:
	var vp_size = get_viewport().size
	var popup = ColorRect.new()
	popup.color = Color(0.1, 0.0, 0.0, 0.92)
	popup.size = Vector2(300, 200)
	popup.position = Vector2(vp_size.x, vp_size.y * 0.35)
	add_child(popup)

	var title_label = Label.new()
	title_label.text = "? ?嗆?蝯?"
	title_label.add_theme_font_size_override("font_size", 22)
	title_label.add_theme_color_override("font_color", COLOR_FLAG)
	title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_label.size = Vector2(300, 36)
	title_label.position = Vector2(0, 8)
	popup.add_child(title_label)

	# ???”
	var rank_colors = [COLOR_GOLD, COLOR_SILVER, COLOR_BRONZE]
	for i in range(min(rank_list.size(), 3)):
		var entry = rank_list[i]
		var rank = entry.get("rank", i + 1)
		var name_str = entry.get("player_name", "?")
		var score = entry.get("score", 0)
		var mult = entry.get("mult", 1.0)
		var row = Label.new()
		row.text = "#%d %s  %d?? ?%.1f" % [rank, name_str, score, mult]
		row.add_theme_font_size_override("font_size", 16)
		row.add_theme_color_override("font_color", rank_colors[i] if i < rank_colors.size() else COLOR_WHITE)
		row.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		row.size = Vector2(300, 28)
		row.position = Vector2(0, 48 + i * 32)
		popup.add_child(row)

	# 敺?湔???
	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 320, 0.3)
	tween.tween_interval(3.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)
