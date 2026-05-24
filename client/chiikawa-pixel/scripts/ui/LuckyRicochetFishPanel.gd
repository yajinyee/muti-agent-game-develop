## LuckyRicochetFishPanel.gd ??撟賊???擳頂蝯梢?選?DAY-220嚗?
## 璆剔????敶?敶???
##
## 閬死閮剛?嚗?
##   - 璈??銝駁?嚗?FF8C00 + #FFD700 + #FF4500 + #FFF8DC嚗?
##   - ricochet_start嚗??脤??? + ?璈怠? + 摨閮?璇?
##   - ricochet_bounce嚗?敶?頝∠? + ??? + 瘚桀????
##   - ricochet_end嚗??脫楚??+ ??敶芋撘???蝷?
extends CanvasLayer

# ?????
var _ricochet_active: bool = false
var _my_player_id: String = ""
var _timer_bar: Control = null
var _banner: Control = null

# 銝駁?憿
const COLOR_PRIMARY   = Color("#FF8C00")  # 璈
const COLOR_GOLD      = Color("#FFD700")  # ?
const COLOR_ACCENT    = Color("#FF4500")  # 瘛望?
const COLOR_LIGHT     = Color("#FFF8DC")  # 瘛⊿?
const COLOR_BG        = Color(0.15, 0.08, 0.0, 0.85)

func _ready() -> void:
	layer = 25  # 撟賊???擳?踹惜蝝?
	# ???砍?拙振 ID
	if GameManager.has_method("get_player_id"):
		_my_player_id = GameManager.get_player_id()

## ??撟賊???擳???
func handle_lucky_ricochet_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"ricochet_start":
			_on_ricochet_start(payload)
		"ricochet_bounce":
			_on_ricochet_bounce(payload)
		"ricochet_end":
			_on_ricochet_end(payload)

## ricochet_start ????璅∪???
func _on_ricochet_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var duration_sec: int = payload.get("duration_sec", 8)
	var player_id: String = payload.get("player_id", "")

	_ricochet_active = true

	# 璈????
	_flash_screen(COLOR_PRIMARY, 0.25)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_GOLD, 0.2)

	# ?璈怠?
	_show_banner("? ??璅∪?嚗?, "%s 閫貊??擳?瘥??憭?敶?3 甈? % player_name, duration_sec)

## ricochet_bounce ??摮????賭葉
func _on_ricochet_bounce(payload: Dictionary) -> void:
	var bounce_num: int = payload.get("bounce_num", 1)
	var killed: bool = payload.get("killed", false)
	var reward: int = payload.get("reward", 0)
	var x: float = payload.get("x", 0.0)
	var y: float = payload.get("y", 0.0)

	# ???嚗???甈⊥瘙箏?憭批?嚗?
	var radius: float = 30.0 - float(bounce_num) * 5.0
	_show_bounce_effect(Vector2(x, y), radius, bounce_num)

	# 瘚桀???嚗???湔?憿舐內?嚗?
	if killed and reward > 0:
		var color = COLOR_GOLD if bounce_num == 1 else COLOR_PRIMARY
		_show_float_text("??+%d" % reward, color, Vector2(x, y))

## ricochet_end ????璅∪?蝯?
func _on_ricochet_end(payload: Dictionary) -> void:
	_ricochet_active = false
	_hide_banner()

	# 璈瘛∪?內
	var label = Label.new()
	label.text = "? ??蝯?"
	label.add_theme_font_size_override("font_size", 16)
	label.add_theme_color_override("font_color", COLOR_PRIMARY)
	label.set_anchors_preset(Control.PRESET_CENTER)
	label.position = get_viewport().size / 2 - Vector2(60, 20)
	add_child(label)

	var tween = label.create_tween()
	tween.tween_property(label, "modulate:a", 0.0, 0.8)
	tween.tween_callback(label.queue_free)

# ---- 頛?賣 ----

## 憿舐內?璈怠?
func _show_banner(title: String, subtitle: String, duration_sec: int) -> void:
	_hide_banner()

	var banner = Control.new()
	banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	banner.position = Vector2(0, 8)
	banner.size = Vector2(get_viewport().size.x, 52)
	add_child(banner)

	var bg = ColorRect.new()
	bg.color = COLOR_BG
	bg.size = banner.size
	banner.add_child(bg)

	var title_label = Label.new()
	title_label.text = title
	title_label.add_theme_font_size_override("font_size", 20)
	title_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	title_label.position = Vector2(12, 4)
	banner.add_child(title_label)

	var sub_label = Label.new()
	sub_label.text = subtitle
	sub_label.add_theme_font_size_override("font_size", 13)
	sub_label.add_theme_color_override("font_color", COLOR_LIGHT)
	sub_label.position = Vector2(12, 28)
	banner.add_child(sub_label)

	# 閮?璇?摨嚗???璈撓霈?
	var timer_bar = ColorRect.new()
	timer_bar.name = "TimerBar"
	timer_bar.color = COLOR_PRIMARY
	timer_bar.position = Vector2(0, 48)
	timer_bar.size = Vector2(get_viewport().size.x, 4)
	banner.add_child(timer_bar)

	var tween = banner.create_tween()
	tween.tween_property(timer_bar, "size:x", 0.0, float(duration_sec))
	tween.parallel().tween_property(timer_bar, "color", COLOR_ACCENT, float(duration_sec))

	_banner = banner

## ?梯?璈怠?
func _hide_banner() -> void:
	if _banner != null and is_instance_valid(_banner):
		_banner.queue_free()
	_banner = null

## ?????嚗?????
func _show_bounce_effect(pos: Vector2, radius: float, bounce_num: int) -> void:
	# 憭?嚗??莎?
	var outer = ColorRect.new()
	outer.color = Color(COLOR_PRIMARY.r, COLOR_PRIMARY.g, COLOR_PRIMARY.b, 0.6)
	outer.size = Vector2(radius * 2, radius * 2)
	outer.position = pos - Vector2(radius, radius)
	add_child(outer)

	# ??甈⊥璅?
	var num_label = Label.new()
	num_label.text = "??d" % bounce_num
	num_label.add_theme_font_size_override("font_size", 12)
	num_label.add_theme_color_override("font_color", COLOR_GOLD)
	num_label.position = pos - Vector2(10, 8)
	add_child(num_label)

	# ?湔?
	var tween = outer.create_tween()
	tween.tween_property(outer, "size", Vector2(radius * 4, radius * 4), 0.3)
	tween.parallel().tween_property(outer, "position", pos - Vector2(radius * 2, radius * 2), 0.3)
	tween.parallel().tween_property(outer, "modulate:a", 0.0, 0.3)
	tween.tween_callback(outer.queue_free)

	var tween2 = num_label.create_tween()
	tween2.tween_property(num_label, "modulate:a", 0.0, 0.4)
	tween2.tween_callback(num_label.queue_free)

## ?刻撟?????
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.4)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)

## 瘚桀???
func _show_float_text(text: String, color: Color, pos: Vector2) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 16)
	label.add_theme_color_override("font_color", color)
	label.position = pos - Vector2(30, 15)
	add_child(label)

	var tween = label.create_tween()
	tween.tween_property(label, "position:y", label.position.y - 35, 0.7)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.7)
	tween.tween_callback(label.queue_free)
