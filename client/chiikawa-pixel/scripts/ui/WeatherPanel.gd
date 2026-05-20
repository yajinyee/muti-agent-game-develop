# WeatherPanel.gd — 天氣系統面板（DAY-087）
# 顯示當前天氣狀態、效果說明、倒數計時
# 天氣切換時顯示大通知彈窗
extends Control

# 面板節點
var _panel: PanelContainer
var _icon_label: Label
var _name_label: Label
var _desc_label: Label
var _timer_label: Label
var _effect_label: Label

# 通知彈窗
var _notify_overlay: Control
var _notify_icon: Label
var _notify_name: Label
var _notify_desc: Label

# 當前天氣資料
var _current_weather: Dictionary = {}
var _remaining_seconds: int = 0

# 天氣顏色對應
const WEATHER_COLORS = {
	"clear":    Color(1.0, 0.95, 0.7),   # 暖黃
	"rain":     Color(0.6, 0.8, 1.0),    # 淡藍
	"storm":    Color(0.5, 0.5, 0.8),    # 深藍紫
	"fog":      Color(0.8, 0.8, 0.8),    # 灰白
	"sunshine": Color(1.0, 0.85, 0.2),   # 金黃
	"blizzard": Color(0.7, 0.9, 1.0),    # 冰藍
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
	# 連接 GameManager 訊號
	if GameManager.has_signal("weather_updated"):
		GameManager.weather_updated.connect(_on_weather_updated)

func _build_ui():
	# 主面板（右上角，天氣圖示 + 名稱 + 倒數）
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

	# 第一行：圖示 + 名稱
	var row1 = HBoxContainer.new()
	row1.add_theme_constant_override("separation", 4)
	vbox.add_child(row1)

	_icon_label = Label.new()
	_icon_label.text = "☀️"
	_icon_label.add_theme_font_size_override("font_size", 16)
	row1.add_child(_icon_label)

	_name_label = Label.new()
	_name_label.text = "晴天"
	_name_label.add_theme_font_size_override("font_size", 11)
	_name_label.add_theme_color_override("font_color", Color(1.0, 0.95, 0.7))
	row1.add_child(_name_label)

	# 第二行：效果簡述
	_effect_label = Label.new()
	_effect_label.text = "正常"
	_effect_label.add_theme_font_size_override("font_size", 9)
	_effect_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	vbox.add_child(_effect_label)

	# 第三行：倒數計時
	_timer_label = Label.new()
	_timer_label.text = "5:00"
	_timer_label.add_theme_font_size_override("font_size", 9)
	_timer_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	vbox.add_child(_timer_label)

	# 通知彈窗（天氣切換時顯示）
	_build_notify_overlay()

func _build_notify_overlay():
	_notify_overlay = Control.new()
	_notify_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_notify_overlay.visible = false
	_notify_overlay.z_index = 65
	add_child(_notify_overlay)

	# 半透明背景
	var bg = ColorRect.new()
	bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0, 0, 0, 0.5)
	_notify_overlay.add_child(bg)

	# 中央通知框
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

	# 天氣變化標題
	var title = Label.new()
	title.text = "🌤 天氣變化"
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_font_size_override("font_size", 12)
	title.add_theme_color_override("font_color", Color(0.7, 0.8, 1.0))
	vbox.add_child(title)

	# 天氣圖示
	_notify_icon = Label.new()
	_notify_icon.text = "☀️"
	_notify_icon.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_notify_icon.add_theme_font_size_override("font_size", 36)
	vbox.add_child(_notify_icon)

	# 天氣名稱
	_notify_name = Label.new()
	_notify_name.text = "晴天"
	_notify_name.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_notify_name.add_theme_font_size_override("font_size", 18)
	_notify_name.add_theme_color_override("font_color", Color(1.0, 0.95, 0.7))
	vbox.add_child(_notify_name)

	# 效果說明
	_notify_desc = Label.new()
	_notify_desc.text = "風平浪靜，正常捕魚"
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
	# 快到期時變紅
	if secs < 30:
		_timer_label.add_theme_color_override("font_color", Color(1.0, 0.4, 0.4))
	else:
		_timer_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))

func _on_weather_updated(data: Dictionary):
	var is_new = data.get("is_new", false)
	_current_weather = data
	_remaining_seconds = float(data.get("remaining_seconds", 300))

	var wtype = data.get("type", "clear")
	var icon = data.get("icon", "☀️")
	var name = data.get("name", "晴天")
	var desc = data.get("description", "")
	var reward_mult = data.get("reward_mult", 1.0)
	var speed_mult = data.get("speed_mult", 1.0)

	# 更新主面板
	_icon_label.text = icon
	_name_label.text = name

	# 效果簡述
	var effects = []
	if reward_mult > 1.0:
		effects.append("獎勵×%.1f" % reward_mult)
	if speed_mult > 1.0:
		effects.append("速度×%.1f" % speed_mult)
	if data.get("rare_chance_bonus", 0.0) > 0:
		effects.append("稀有+%d%%" % int(data.get("rare_chance_bonus", 0.0) * 100))
	if data.get("gold_fish_bonus", 0.0) > 0:
		effects.append("金幣魚+%d%%" % int(data.get("gold_fish_bonus", 0.0) * 100))
	if data.get("boss_chance_bonus", 0.0) > 0:
		effects.append("BOSS+%d%%" % int(data.get("boss_chance_bonus", 0.0) * 100))
	if data.get("fog_effect", false):
		effects.append("濃霧視野")

	if effects.size() > 0:
		_effect_label.text = " ".join(effects)
	else:
		_effect_label.text = "正常"

	# 更新顏色
	var color = WEATHER_COLORS.get(wtype, Color.WHITE)
	_name_label.add_theme_color_override("font_color", color)

	# 更新面板背景
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

	# 天氣切換時顯示通知彈窗
	if is_new:
		_show_weather_notify(icon, name, desc, color)

func _show_weather_notify(icon: String, name: String, desc: String, color: Color):
	_notify_icon.text = icon
	_notify_name.text = name
	_notify_name.add_theme_color_override("font_color", color)
	_notify_desc.text = desc

	# 更新通知框邊框顏色
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

	# 淡入 → 停留 → 淡出
	var tween = create_tween()
	tween.tween_property(_notify_overlay, "modulate:a", 1.0, 0.3)
	tween.tween_interval(2.5)
	tween.tween_property(_notify_overlay, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func(): _notify_overlay.visible = false)
