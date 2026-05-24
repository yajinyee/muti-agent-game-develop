## WeatherSurgePanel.gd
##
## 憭拇除皝抒鈭辣?Ｘ嚗AY-127嚗?
## 憭拇除???孛?潛??璅黎皝改??冽?撱?璈怠?
## 雿蔭嚗??其葉憭格??伐?z_index=77嚗暺????Ｘ 76 銋?嚗?
## 閮剛?嚗?
##   - 璈怠?敺??冽??伐?憿舐內皝抒?迂??蝷箝閮?
##   - 皝抒???喃?閫＊蝷箏?????蝷箏嚗????馳擳????嚗?
##   - 皝抒蝯??楚??

extends Control

# ---- 撣豢 ----
const BANNER_HEIGHT := 64.0
const SLIDE_DURATION := 0.4
const INDICATOR_SIZE := Vector2(140, 48)

# ---- 蝭暺???----
var _banner: Panel = null
var _banner_icon: Label = null
var _banner_title: Label = null
var _banner_desc: Label = null
var _banner_timer: Label = null
var _indicator: Panel = null
var _indicator_rare: Label = null
var _indicator_gold: Label = null
var _pixel_font: FontFile = null

# ---- ???----
var _is_active: bool = false
var _end_time: float = 0.0
var _surge_name: String = ""
var _rare_bonus: float = 0.0
var _gold_bonus: float = 0.0
var _banner_color: Color = Color(0.29, 0.56, 0.85)

func setup(font: FontFile) -> void:
	_pixel_font = font
	_build_ui()
	_connect_signals()

func _build_ui() -> void:
	# ?璈怠?
	_banner = Panel.new()
	_banner.name = "WeatherSurgeBanner"
	_banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	_banner.size = Vector2(1280, BANNER_HEIGHT)
	_banner.position = Vector2(0, -BANNER_HEIGHT)  # ???函?Ｗ?
	_banner.visible = false
	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.1, 0.2, 0.4, 0.92)
	banner_style.border_color = _banner_color
	banner_style.set_border_width_all(2)
	banner_style.corner_radius_bottom_left = 6
	banner_style.corner_radius_bottom_right = 6
	_banner.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner)

	# 璈怠??內
	_banner_icon = Label.new()
	_banner_icon.position = Vector2(20, 10)
	_banner_icon.size = Vector2(44, 44)
	_banner_icon.add_theme_font_size_override("font_size", 32)
	_banner_icon.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.add_child(_banner_icon)

	# 璈怠?璅?
	_banner_title = Label.new()
	_banner_title.position = Vector2(72, 6)
	_banner_title.size = Vector2(600, 28)
	_banner_title.add_theme_color_override("font_color", Color(0.4, 0.85, 1.0))
	_banner_title.add_theme_font_size_override("font_size", 18)
	if is_instance_valid(_pixel_font):
		_banner_title.add_theme_font_override("font", _pixel_font)
	_banner.add_child(_banner_title)

	# 璈怠??膩
	_banner_desc = Label.new()
	_banner_desc.position = Vector2(72, 34)
	_banner_desc.size = Vector2(700, 22)
	_banner_desc.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
	_banner_desc.add_theme_font_size_override("font_size", 12)
	if is_instance_valid(_pixel_font):
		_banner_desc.add_theme_font_override("font", _pixel_font)
	_banner.add_child(_banner_desc)

	# ?閮?嚗?湛?
	_banner_timer = Label.new()
	_banner_timer.position = Vector2(1100, 18)
	_banner_timer.size = Vector2(160, 28)
	_banner_timer.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	_banner_timer.add_theme_font_size_override("font_size", 16)
	_banner_timer.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	if is_instance_valid(_pixel_font):
		_banner_timer.add_theme_font_override("font", _pixel_font)
	_banner.add_child(_banner_timer)

	# ?喃?閫???蝷箏
	_indicator = Panel.new()
	_indicator.name = "WeatherSurgeIndicator"
	_indicator.position = Vector2(1280 - INDICATOR_SIZE.x - 8, 720 - INDICATOR_SIZE.y - 8)
	_indicator.size = INDICATOR_SIZE
	_indicator.visible = false
	var ind_style = StyleBoxFlat.new()
	ind_style.bg_color = Color(0.05, 0.1, 0.2, 0.85)
	ind_style.border_color = Color(0.3, 0.6, 0.9, 0.8)
	ind_style.set_border_width_all(1)
	ind_style.corner_radius_top_left = 4
	ind_style.corner_radius_top_right = 4
	ind_style.corner_radius_bottom_left = 4
	ind_style.corner_radius_bottom_right = 4
	_indicator.add_theme_stylebox_override("panel", ind_style)
	add_child(_indicator)

	# ?內?冽?憿?
	var ind_title = Label.new()
	ind_title.position = Vector2(4, 2)
	ind_title.size = Vector2(132, 16)
	ind_title.add_theme_color_override("font_color", Color(0.4, 0.85, 1.0))
	ind_title.add_theme_font_size_override("font_size", 10)
	ind_title.text = "?? 皝抒銝?"
	if is_instance_valid(_pixel_font):
		ind_title.add_theme_font_override("font", _pixel_font)
	_indicator.add_child(ind_title)

	# 蝔????
	_indicator_rare = Label.new()
	_indicator_rare.position = Vector2(4, 18)
	_indicator_rare.size = Vector2(132, 14)
	_indicator_rare.add_theme_color_override("font_color", Color(0.6, 0.9, 1.0))
	_indicator_rare.add_theme_font_size_override("font_size", 10)
	if is_instance_valid(_pixel_font):
		_indicator_rare.add_theme_font_override("font", _pixel_font)
	_indicator.add_child(_indicator_rare)

	# ?馳擳???
	_indicator_gold = Label.new()
	_indicator_gold.position = Vector2(4, 32)
	_indicator_gold.size = Vector2(132, 14)
	_indicator_gold.add_theme_color_override("font_color", Color(1.0, 0.85, 0.3))
	_indicator_gold.add_theme_font_size_override("font_size", 10)
	if is_instance_valid(_pixel_font):
		_indicator_gold.add_theme_font_override("font", _pixel_font)
	_indicator.add_child(_indicator_gold)

func _connect_signals() -> void:
	if GameManager.has_signal("weather_surge_started"):
		GameManager.weather_surge_started.connect(_on_weather_surge_started)
	if GameManager.has_signal("weather_surge_ended"):
		GameManager.weather_surge_ended.connect(_on_weather_surge_ended)

# ---- 鈭辣?? ----

func _on_weather_surge_started(data: Dictionary) -> void:
	_surge_name = data.get("surge_name", "憭拇除皝抒")
	var icon = data.get("surge_icon", "??")
	var message = data.get("surge_message", "蝔?璅之?鳩?橘?")
	var duration = data.get("duration", 30)
	_rare_bonus = data.get("rare_bonus", 0.0)
	_gold_bonus = data.get("gold_bonus", 0.0)
	var color_hex = data.get("color", "#4A90D9")
	_banner_color = Color(color_hex)
	_is_active = true
	_end_time = Time.get_ticks_msec() / 1000.0 + duration

	# ?湔璈怠?璅??憿
	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.1, 0.2, 0.4, 0.92)
	banner_style.border_color = _banner_color
	banner_style.set_border_width_all(2)
	banner_style.corner_radius_bottom_left = 6
	banner_style.corner_radius_bottom_right = 6
	_banner.add_theme_stylebox_override("panel", banner_style)

	# ?湔?批捆
	_banner_icon.text = icon
	_banner_title.text = _surge_name + "嚗?"
	_banner_desc.text = message

	# ?湔?內??
	_indicator_rare.text = "? 蝔??+" + str(int(_rare_bonus * 100)) + "%"
	_indicator_gold.text = "?? ?馳擳?+" + str(int(_gold_bonus * 100)) + "%"

	# 憿舐內銝行???
	_banner.visible = true
	_indicator.visible = true
	_banner.position = Vector2(0, -BANNER_HEIGHT)
	var tween = create_tween()
	tween.tween_property(_banner, "position", Vector2(0, 0), SLIDE_DURATION).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)

	# ????
	_do_flash_effect()

func _on_weather_surge_ended(data: Dictionary) -> void:
	_is_active = false
	_banner_timer.text = "蝯?"
	_banner_timer.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))

	# 皛銝阡??
	var tween = create_tween()
	tween.tween_property(_banner, "position", Vector2(0, -BANNER_HEIGHT), SLIDE_DURATION).set_ease(Tween.EASE_IN)
	tween.tween_callback(func():
		_banner.visible = false
	)
	# ?內?冽楚??
	var tween2 = create_tween()
	tween2.tween_property(_indicator, "modulate:a", 0.0, 0.5)
	tween2.tween_callback(func():
		_indicator.visible = false
		_indicator.modulate.a = 1.0
	)

func _do_flash_effect() -> void:
	# ?刻撟?恍???憭抵??莎?
	var flash = ColorRect.new()
	flash.color = Color(_banner_color.r, _banner_color.g, _banner_color.b, 0.25)
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.z_index = 200
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)

# ---- 瘥??湔?閮? ----

func _process(delta: float) -> void:
	if not _is_active or not is_instance_valid(_banner_timer):
		return
	var now = Time.get_ticks_msec() / 1000.0
	var remaining = _end_time - now
	if remaining <= 0:
		_banner_timer.text = "00:00"
		return
	var secs = int(remaining)
	_banner_timer.text = "%02d:%02d" % [secs / 60, secs % 60]
	# ?敺?10 蝘??脤???
	if remaining <= 10:
		var blink = fmod(now * 2.0, 1.0) > 0.5
		_banner_timer.add_theme_color_override("font_color",
			Color(1.0, 0.3, 0.3) if blink else Color(1.0, 0.6, 0.6))
	else:
		_banner_timer.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
