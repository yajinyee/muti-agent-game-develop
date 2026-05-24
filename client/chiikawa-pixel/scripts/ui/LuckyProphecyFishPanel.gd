## LuckyProphecyFishPanel.gd ??撟賊???擳頂蝯梢?選?DAY-243嚗?
## 璆剔????閮???格?????
##
## 閬死閮剛?嚗?
##   - 蝝恍???銝駁?嚗?9B59B6 + #F39C12 + #D7BDE2 + #FFF9E6嚗?
##   - prophecy_start嚗換?脖?甈∪撥?? + ?璈怠? + ??????嚗之摮?+ ?格?璅? + 閮?璇?
##   - prophecy_broadcast嚗??典?璈怠?嚗?冽??犖閫貊??嚗?
##   - prophecy_fulfilled嚗??脖?甈∪撥?? + ???????嚗?.5?之摮?+ 蝯?敶?
##   - prophecy_broadcast_fulfilled嚗??典?璈怠?嚗?冽?????嚗?
##   - prophecy_transfer嚗??脤???+ ?????頧宏嚗?蝷?+ ?啁璅?閮?
##   - prophecy_broadcast_transfer嚗??典?璈怠?嚗?冽???頧宏嚗?
##   - prophecy_fail嚗?脤???+ ?????憭望?嚗P-20%??蝷?
extends CanvasLayer

# 銝駁?憿
const COLOR_PROPHECY = Color("#9B59B6")  # 蝝怨嚗蜓憿?
const COLOR_GOLD     = Color("#F39C12")  # ?嚗???
const COLOR_ORANGE   = Color("#E67E22")  # 璈嚗?蝘鳴?
const COLOR_PALE     = Color("#FFF9E6")  # 璆菜楚暺?
const COLOR_FAIL     = Color("#7F8C8D")  # ?啗嚗仃??
const COLOR_WHITE    = Color("#FFFFFF")  # ?質
const COLOR_MARK     = Color("#E74C3C")  # 蝝嚗璅?閮?

# 閮?璇?暺?
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null
var _duration_sec: int = 12

# ?格?璅?蝭暺?憿舐內?券?閮?格?銝嚗?
var _target_marker: Control = null
var _current_target_id: String = ""

func _ready() -> void:
	layer = 2  # 撟賊???擳?踹惜蝝?

## ??撟賊???擳???
func handle_lucky_prophecy_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"prophecy_start":
			_on_prophecy_start(payload)
		"prophecy_broadcast":
			_on_prophecy_broadcast(payload)
		"prophecy_fulfilled":
			_on_prophecy_fulfilled(payload)
		"prophecy_broadcast_fulfilled":
			_on_prophecy_broadcast_fulfilled(payload)
		"prophecy_transfer":
			_on_prophecy_transfer(payload)
		"prophecy_broadcast_transfer":
			_on_prophecy_broadcast_transfer(payload)
		"prophecy_fail":
			_on_prophecy_fail(payload)

## prophecy_start ??????嚗犖閮嚗?
func _on_prophecy_start(payload: Dictionary) -> void:
	_duration_sec = payload.get("duration_sec", 12)
	var kill_mult: float = payload.get("kill_mult", 3.5)
	_current_target_id = payload.get("target_id", "")
	var target_x: float = payload.get("x", 0.0)
	var target_y: float = payload.get("y", 0.0)
	var vp_size = get_viewport().size

	# 蝝怨銝活撘琿???
	_flash_screen(COLOR_PROPHECY, 0.5, 3)

	# ?璈怠?
		_show_banner("Prophecy Activated!", COLOR_PROPHECY, 3.5)

	# 銝剖亢憭批?
		_show_big_text("Prophecy!", COLOR_PROPHECY, 52, 2.5)

	# ??隤芣?
		_show_sub_text("Kill mult: x%.1f" % kill_mult, COLOR_GOLD, 2.0)

	# ?喳鞊?閮?璇?
	_create_timer_bar(_duration_sec)

	# ?格?璅?嚗?格?雿蔭憿舐內??璅?嚗?
	_create_target_marker(target_x, target_y)

## prophecy_broadcast ???冽?撱??犖閫貊??
func _on_prophecy_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "?摰?)"
	var duration_sec: int = payload.get("duration_sec", 12)
	_show_top_banner("? %s 閫貊??嚗?d 蝘餈質馱???格?嚗? % [player_name, duration_sec], COLOR_PROPHECY, 3.0)"

## prophecy_fulfilled ??????嚗犖閮嚗?
func _on_prophecy_fulfilled(payload: Dictionary) -> void:
	var kill_mult: float = payload.get("kill_mult", 3.5)
	var reward: int = payload.get("reward", 0)

	# 皜閮?璇??格?璅?
	_clear_timer_bar()
	_clear_target_marker()

	# ?銝活撘琿???
	_flash_screen(COLOR_GOLD, 0.6, 3)

	# 銝剖亢憭批?
	_show_big_text("? ????嚗?.1f" % kill_mult, COLOR_GOLD, 3.0)

	# 蝯?敶?
	_show_reward_popup(reward, kill_mult)

## prophecy_broadcast_fulfilled ???冽?撱?????
func _on_prophecy_broadcast_fulfilled(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "?摰?)"
	var kill_mult: float = payload.get("kill_mult", 3.5)
	var reward: int = payload.get("reward", 0)
	_show_top_banner("? %s ????嚗?.1f ??嚗敺?%d ?馳嚗? % [player_name, kill_mult, reward], COLOR_GOLD, 3.5)"

## prophecy_transfer ????頧宏嚗犖閮嚗?
func _on_prophecy_transfer(payload: Dictionary) -> void:
	var new_target_id: String = payload.get("target_id", "")
	var new_x: float = payload.get("x", 0.0)
	var new_y: float = payload.get("y", 0.0)
	var transfer_count: int = payload.get("transfer_count", 1)
	var kill_mult: float = payload.get("kill_mult", 3.5)

	_current_target_id = new_target_id

	# 璈??
	_flash_screen(COLOR_ORANGE, 0.3, 1)

	# ?內??
	_show_big_text("? ??頧宏嚗?蝚?%d 甈∴?" % transfer_count, COLOR_ORANGE, 2.0)

	# ?湔?格?璅?雿蔭
	_clear_target_marker()
	_create_target_marker(new_x, new_y)

## prophecy_broadcast_transfer ???冽?撱???頧宏
func _on_prophecy_broadcast_transfer(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "?摰?)"
	var transfer_count: int = payload.get("transfer_count", 1)
	_show_top_banner("? %s ??閮頧宏嚗?蝚?%d 甈∴?" % [player_name, transfer_count], COLOR_ORANGE, 2.5)

## prophecy_fail ????憭望?嚗?誨?哨?
func _on_prophecy_fail(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "?摰?)"
	var affected_count: int = payload.get("affected_count", 0)

	# 皜閮?璇??格?璅?
	_clear_timer_bar()
	_clear_target_marker()

	# ?啗??
	_flash_screen(COLOR_FAIL, 0.4, 2)

	# ?內??
	_show_big_text("? ??憭望?嚗P -20%%", COLOR_FAIL, 2.5)
	_show_sub_text("%s ??閮?芾??嚗?d ?璅??啣摰喉?" % [player_name, affected_count], COLOR_FAIL, 2.0)

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

## ?璈怠?嚗?蝥＊蝷綽?
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

## ?撠帖撟??冽?撱??剁?
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
func _show_big_text(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 48)
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
	label.add_theme_font_size_override("font_size", 22)
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

	# ?璇?
	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.2, 0.1, 0.3, 0.6)
	_timer_bar_bg.size = Vector2(16, 200)
	_timer_bar_bg.position = Vector2(vp_size.x - 28, vp_size.y * 0.3)
	add_child(_timer_bar_bg)

	# 閮?璇?蝝徉??撓霈?敺?敺銝葬?哨?
	_timer_bar = ColorRect.new()
	_timer_bar.color = COLOR_PROPHECY
	_timer_bar.size = Vector2(16, 200)
	_timer_bar.position = Vector2(vp_size.x - 28, vp_size.y * 0.3)
	add_child(_timer_bar)

	# 閮?璇??恬?敺?200 蝮桀 0嚗?
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

## 撱箇??格?璅?嚗?格?雿蔭憿舐內????? 璅?嚗?
func _create_target_marker(target_x: float, target_y: float) -> void:
	_clear_target_marker()
	if target_x <= 0 and target_y <= 0:
		return

	_target_marker = Control.new()
	_target_marker.position = Vector2(target_x - 20, target_y - 50)
	_target_marker.size = Vector2(40, 40)
	add_child(_target_marker)

	# ????閮?摮?
	var label = Label.new()
	label.text = "?"
	label.add_theme_font_size_override("font_size", 28)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.size = Vector2(40, 40)
	label.position = Vector2.ZERO
	_target_marker.add_child(label)

	# ???
	var tween = _target_marker.create_tween().set_loops()
	tween.tween_property(label, "modulate:a", 0.3, 0.4)
	tween.tween_property(label, "modulate:a", 1.0, 0.4)

## 皜?格?璅?
func _clear_target_marker() -> void:
	if is_instance_valid(_target_marker):
		_target_marker.queue_free()
		_target_marker = null
	_current_target_id = ""

## 蝯?敶?
func _show_reward_popup(reward: int, kill_mult: float) -> void:
	var vp_size = get_viewport().size
	var popup = ColorRect.new()
	popup.color = Color(0.1, 0.05, 0.2, 0.92)
	popup.size = Vector2(280, 120)
	popup.position = Vector2(vp_size.x, vp_size.y * 0.4)  # 敺?湔???
	add_child(popup)

	var title_label = Label.new()
	title_label.text = "? ????嚗?"
	title_label.add_theme_font_size_override("font_size", 22)
	title_label.add_theme_color_override("font_color", COLOR_GOLD)
	title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_label.size = Vector2(280, 40)
	title_label.position = Vector2(0, 8)
	popup.add_child(title_label)

	var mult_label = Label.new()
	mult_label.text = "?%.1f ????" % kill_mult
	mult_label.add_theme_font_size_override("font_size", 18)
	mult_label.add_theme_color_override("font_color", COLOR_PROPHECY)
	mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mult_label.size = Vector2(280, 32)
	mult_label.position = Vector2(0, 44)
	popup.add_child(mult_label)

	var reward_label = Label.new()
	reward_label.text = "+%d ?馳" % reward
	reward_label.add_theme_font_size_override("font_size", 20)
	reward_label.add_theme_color_override("font_color", COLOR_GOLD)
	reward_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_label.size = Vector2(280, 32)
	reward_label.position = Vector2(0, 76)
	popup.add_child(reward_label)

	# 敺?湔???
	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 300, 0.3)
	tween.tween_interval(2.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)
