# WeatherPanel.gd ??憭拇除蝟餌絞?Ｘ嚗AY-087嚗?
# 憿舐內?嗅?憭拇除????牧?閮?
# 憭拇除???＊蝷箏之?敶?
extends Control

# ?Ｘ蝭暺?
var _panel: PanelContainer
var _icon_label: Label
var _name_label: Label
var _desc_label: Label
var _timer_label: Label
var _effect_label: Label

# ?敶?
var _notify_overlay: Control
var _notify_icon: Label
var _notify_name: Label
var _notify_desc: Label

# ?嗅?憭拇除鞈?
var _current_weather: Dictionary = {}
var _remaining_seconds: int = 0

# 憭拇除憿撠?
const WEATHER_COLORS = {
	"clear":    Color(1.0, 0.95, 0.7),   # ??
	"rain":     Color(0.6, 0.8, 1.0),    # 瘛∟?
	"storm":    Color(0.5, 0.5, 0.8),    # 瘛梯?蝝?
	"fog":      Color(0.8, 0.8, 0.8),    # ?啁
	"sunshine": Color(1.0, 0.85, 0.2),   # ??
	"blizzard": Color(0.7, 0.9, 1.0),    # ?啗?
}

const WEATHER_BG_COLORS = {
	"clear":    Color(0.15, 0.15, 0.1, 0.85),
	"rain":     Color(0.05, 0.1, 0.2, 0.85),
	"storm":    Color(0.05, 0.05, 0.15, 0.9),
	"fog":      Color(0.15, 0.15, 0.15, 0.9),
	"sunshine": Color(0.2, 0.15, 0.0, 0.85),
	"blizzard": Color(0.05, 0.1, 0.2, 0.9),
}

func _ready():
	_build_ui()
	# ?? GameManager 閮?
	if GameManager.has_signal("weather_updated"):
		GameManager.weather_updated.connect(_on_weather_updated)

func _build_ui():
	# 銝駁?選??喃?閫?憭拇除?內 + ?迂 + ?嚗?
	_panel = PanelContainer.new()
	_panel.position = Vector2(1180, 8)
	_panel.custom_minimum_size = Vector2(120, 56)
	add_child(_panel)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.1, 0.1, 0.85)
	style.corner_radius_top_left = 6
	style.corner_radius_top_right = 6
	style.corner_radius_bottom_left = 6
	style.corner_radius_bottom_right = 6
	style.border_width_left = 1
	style.border_width_right = 1
	style.border_width_top = 1
	style.border_width_bottom = 1
	style.border_color = Color(0.4, 0.4, 0.4, 0.6)
	_panel.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 2)
	_panel.add_child(vbox)

	# 蝚砌?銵??內 + ?迂
	var row1 = HBoxContainer.new()
	row1.add_theme_constant_override("separation", 4)
	vbox.add_child(row1)

	_icon_label = Label.new()
	_icon_label.text = "?儭?"
	_icon_label.add_theme_font_size_override("font_size", 16)
	row1.add_child(_icon_label)

	_name_label = Label.new()
	_name_label.text = "?游予"
	_name_label.add_theme_font_size_override("font_size", 11)
	_name_label.add_theme_color_override("font_color", Color(1.0, 0.95, 0.7))
	row1.add_child(_name_label)

	# 蝚砌?銵???蝪∟膩
	_effect_label = Label.new()
	_effect_label.text = "甇?虜"
	_effect_label.add_theme_font_size_override("font_size", 9)
	_effect_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	vbox.add_child(_effect_label)

	# 蝚砌?銵??閮?
	_timer_label = Label.new()
	_timer_label.text = "5:00"
	_timer_label.add_theme_font_size_override("font_size", 9)
	_timer_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	vbox.add_child(_timer_label)

	# ?敶?嚗予瘞????憿舐內嚗?
	_build_notify_overlay()

func _build_notify_overlay():
	_notify_overlay = Control.new()
	_notify_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_notify_overlay.visible = false
	_notify_overlay.z_index = 65
	add_child(_notify_overlay)

	# ???
	var bg = ColorRect.new()
	bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0, 0, 0, 0.5)
	_notify_overlay.add_child(bg)

	# 銝剖亢?獢?
	var box = PanelContainer.new()
	box.position = Vector2(440, 260)
	box.custom_minimum_size = Vector2(400, 160)
	_notify_overlay.add_child(box)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.05, 0.15, 0.95)
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	style.border_width_left = 2
	style.border_width_right = 2
	style.border_width_top = 2
	style.border_width_bottom = 2
	style.border_color = Color(0.5, 0.7, 1.0, 0.8)
	box.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	vbox.add_theme_constant_override("separation", 8)
	box.add_child(vbox)

	# 憭拇除霈?璅?
	var title = Label.new()
	title.text = "? 憭拇除霈?"
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_font_size_override("font_size", 12)
	title.add_theme_color_override("font_color", Color(0.7, 0.8, 1.0))
	vbox.add_child(title)

	# 憭拇除?內
	_notify_icon = Label.new()
	_notify_icon.text = "?儭?"
	_notify_icon.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_notify_icon.add_theme_font_size_override("font_size", 36)
	vbox.add_child(_notify_icon)

	# 憭拇除?迂
	_notify_name = Label.new()
	_notify_name.text = "?游予"
	_notify_name.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_notify_name.add_theme_font_size_override("font_size", 18)
	_notify_name.add_theme_color_override("font_color", Color(1.0, 0.95, 0.7))
	vbox.add_child(_notify_name)

	# ??隤芣?
	_notify_desc = Label.new()
	_notify_desc.text = "憸典像瘚芷?嚗迤撣豢?擳?"
	_notify_desc.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_notify_desc.add_theme_font_size_override("font_size", 12)
	_notify_desc.add_theme_color_override("font_color", Color(0.9, 0.9, 0.9))
	vbox.add_child(_notify_desc)

func _process(delta: float):
	if _remaining_seconds > 0:
		_remaining_seconds -= delta
		if _remaining_seconds < 0:
			_remaining_seconds = 0
		_update_timer_display()

func _update_timer_display():
	var secs = int(_remaining_seconds)
	var mins = secs / 60
	var s = secs % 60
	_timer_label.text = "%d:%02d" % [mins, s]
	# 敹怠??霈?
	if secs < 30:
		_timer_label.add_theme_color_override("font_color", Color(1.0, 0.4, 0.4))
	else:
		_timer_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))

func _on_weather_updated(data: Dictionary):
	var is_new = data.get("is_new", false)
	_current_weather = data
	_remaining_seconds = float(data.get("remaining_seconds", 300))

	var wtype = data.get("type", "clear")
	var icon = data.get("icon", "")"
	var name = data.get("name", "?游予")
	var desc = data.get("description", "")
	var reward_mult = data.get("reward_mult", 1.0)
	var speed_mult = data.get("speed_mult", 1.0)

	# ?湔銝駁??
	_icon_label.text = icon
	_name_label.text = name

	# ??蝪∟膩
	var effects = []
	if reward_mult > 1.0:
		effects.append("??%.1f" % reward_mult)
	if speed_mult > 1.0:
		effects.append("?漲?%.1f" % speed_mult)
	if data.get("rare_chance_bonus", 0.0) > 0:
		effects.append("蝔??%d%%" % int(data.get("rare_chance_bonus", 0.0) * 100))
	if data.get("gold_fish_bonus", 0.0) > 0:
		effects.append("?馳擳?%d%%" % int(data.get("gold_fish_bonus", 0.0) * 100))
	if data.get("boss_chance_bonus", 0.0) > 0:
		effects.append("BOSS+%d%%" % int(data.get("boss_chance_bonus", 0.0) * 100))
	if data.get("fog_effect", false):
		effects.append("瞈閬?")

	if effects.size() > 0:
		_effect_label.text = " ".join(effects)
	else:
		_effect_label.text = "甇?虜"

	# ?湔憿
	var color = WEATHER_COLORS.get(wtype, Color.WHITE)
	_name_label.add_theme_color_override("font_color", color)

	# ?湔?Ｘ?
	var bg_color = WEATHER_BG_COLORS.get(wtype, Color(0.1, 0.1, 0.1, 0.85))
	var style = StyleBoxFlat.new()
	style.bg_color = bg_color
	style.corner_radius_top_left = 6
	style.corner_radius_top_right = 6
	style.corner_radius_bottom_left = 6
	style.corner_radius_bottom_right = 6
	style.border_width_left = 1
	style.border_width_right = 1
	style.border_width_top = 1
	style.border_width_bottom = 1
	style.border_color = color * 0.6
	_panel.add_theme_stylebox_override("panel", style)

	# 憭拇除???＊蝷粹敶?
	if is_new:
		_show_weather_notify(icon, name, desc, color)

func _show_weather_notify(icon: String, name: String, desc: String, color: Color):
	_notify_icon.text = icon
	_notify_name.text = name
	_notify_name.add_theme_color_override("font_color", color)
	_notify_desc.text = desc

	# ?湔?獢?獢???
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.05, 0.15, 0.95)
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	style.border_width_left = 2
	style.border_width_right = 2
	style.border_width_top = 2
	style.border_width_bottom = 2
	style.border_color = color * 0.8
	if _notify_overlay.get_child_count() > 1:
		var box = _notify_overlay.get_child(1)
		if box is PanelContainer:
			box.add_theme_stylebox_override("panel", style)

	_notify_overlay.visible = true
	_notify_overlay.modulate.a = 0.0

	# 瘛∪ ???? ??瘛∪
	var tween = create_tween()
	tween.tween_property(_notify_overlay, "modulate:a", 1.0, 0.3)
	tween.tween_interval(2.5)
	tween.tween_property(_notify_overlay, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func(): _notify_overlay.visible = false)
