## AwakenBossPanel.gd — 覺醒 BOSS 面板（DAY-130）
## 業界依據：JILI Royal Fishing 2026 Awaken Boss
## 顯示覺醒 BOSS 的出現、命中進度、Power Up 爆發和離開動畫
## 特色：Power Up 進度條 + 爆發閃光效果
extends Control

# ---- 節點引用 ----
var _banner: PanelContainer
var _banner_label: Label
var _status_panel: PanelContainer
var _boss_icon_label: Label
var _boss_name_label: Label
var _timer_label: Label
var _hit_count_label: Label
var _mult_range_label: Label
var _powerup_bar_bg: ColorRect
var _powerup_bar_fill: ColorRect
var _powerup_label: Label
var _powerup_flash: ColorRect  # Power Up 爆發閃光

# ---- 狀態 ----
var _active: bool = false
var _instance_id: String = ""
var _boss_name: String = ""
var _boss_icon: String = ""
var _boss_color: Color = Color(1.0, 0.27, 0.0)
var _remaining_seconds: float = 0.0
var _hit_count: int = 0
var _powerup_progress: float = 0.0
var _powerup_threshold: int = 5
var _banner_tween: Tween
var _status_tween: Tween

const BANNER_DURATION := 4.0
const POWERUP_FLASH_DURATION := 0.5

func _ready() -> void:
	_build_ui()
	_status_panel.visible = false
	_banner.visible = false
	_powerup_flash.visible = false

func _build_ui() -> void:
	# 頂部橫幅
	_banner = PanelContainer.new()
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 60)
	_banner.position = Vector2(0, -65)
	add_child(_banner)

	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.15, 0.03, 0.0, 0.93)
	banner_style.border_color = Color(1.0, 0.27, 0.0)
	banner_style.set_border_width_all(2)
	banner_style.set_corner_radius_all(6)
	_banner.add_theme_stylebox_override("panel", banner_style)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_font_size_override("font_size", 17)
	_banner_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.0))
	_banner.add_child(_banner_label)

	# 右側狀態面板（比不死 BOSS 面板稍高，因為有 Power Up 進度條）
	_status_panel = PanelContainer.new()
	_status_panel.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_status_panel.anchor_left = 1.0
	_status_panel.anchor_right = 1.0
	_status_panel.anchor_top = 0.5
	_status_panel.anchor_bottom = 0.5
	_status_panel.offset_left = -185
	_status_panel.offset_right = -8
	_status_panel.offset_top = -100
	_status_panel.offset_bottom = 100
	add_child(_status_panel)

	var status_style = StyleBoxFlat.new()
	status_style.bg_color = Color(0.1, 0.02, 0.0, 0.9)
	status_style.border_color = Color(1.0, 0.27, 0.0)
	status_style.set_border_width_all(2)
	status_style.set_corner_radius_all(8)
	_status_panel.add_theme_stylebox_override("panel", status_style)

	var vbox = VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	vbox.add_theme_constant_override("separation", 4)
	_status_panel.add_child(vbox)

	_boss_icon_label = Label.new()
	_boss_icon_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_boss_icon_label.add_theme_font_size_override("font_size", 30)
	vbox.add_child(_boss_icon_label)

	_boss_name_label = Label.new()
	_boss_name_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_boss_name_label.add_theme_font_size_override("font_size", 14)
	_boss_name_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.0))
	vbox.add_child(_boss_name_label)

	_mult_range_label = Label.new()
	_mult_range_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_mult_range_label.add_theme_font_size_override("font_size", 11)
	_mult_range_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.3))
	vbox.add_child(_mult_range_label)

	_timer_label = Label.new()
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_timer_label.add_theme_font_size_override("font_size", 20)
	_timer_label.add_theme_color_override("font_color", Color.WHITE)
	vbox.add_child(_timer_label)

	_hit_count_label = Label.new()
	_hit_count_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_hit_count_label.add_theme_font_size_override("font_size", 11)
	_hit_count_label.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
	vbox.add_child(_hit_count_label)

	# Power Up 進度條
	_powerup_label = Label.new()
	_powerup_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_powerup_label.add_theme_font_size_override("font_size", 11)
	_powerup_label.add_theme_color_override("font_color", Color(1.0, 0.7, 0.0))
	_powerup_label.text = "⚡ Power Up"
	vbox.add_child(_powerup_label)

	var bar_container = Control.new()
	bar_container.custom_minimum_size = Vector2(140, 14)
	vbox.add_child(bar_container)

	_powerup_bar_bg = ColorRect.new()
	_powerup_bar_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	_powerup_bar_bg.color = Color(0.2, 0.1, 0.0, 0.8)
	bar_container.add_child(_powerup_bar_bg)

	_powerup_bar_fill = ColorRect.new()
	_powerup_bar_fill.anchor_top = 0.0
	_powerup_bar_fill.anchor_bottom = 1.0
	_powerup_bar_fill.anchor_left = 0.0
	_powerup_bar_fill.anchor_right = 0.0
	_powerup_bar_fill.color = Color(1.0, 0.5, 0.0)
	bar_container.add_child(_powerup_bar_fill)

	# Power Up 爆發閃光（全螢幕橙紅色）
	_powerup_flash = ColorRect.new()
	_powerup_flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	_powerup_flash.color = Color(1.0, 0.3, 0.0, 0.0)
	_powerup_flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_powerup_flash)

func _process(delta: float) -> void:
	if not _active:
		return
	_remaining_seconds -= delta
	if _remaining_seconds < 0:
		_remaining_seconds = 0
	_update_timer_display()

func _update_timer_display() -> void:
	if not is_instance_valid(_timer_label):
		return
	var secs = int(_remaining_seconds)
	_timer_label.text = "%02d" % secs
	if _remaining_seconds <= 5.0:
		_timer_label.add_theme_color_override("font_color", Color.RED)
		_timer_label.modulate.a = 0.5 + 0.5 * sin(Time.get_ticks_msec() * 0.01)
	else:
		_timer_label.add_theme_color_override("font_color", Color.WHITE)
		_timer_label.modulate.a = 1.0

func _update_powerup_bar() -> void:
	if not is_instance_valid(_powerup_bar_fill):
		return
	_powerup_bar_fill.anchor_right = _powerup_progress
	# 接近滿時顏色變亮
	if _powerup_progress >= 0.8:
		_powerup_bar_fill.color = Color(1.0, 0.8, 0.0)
	else:
		_powerup_bar_fill.color = Color(1.0, 0.5, 0.0)

# ---- 外部呼叫 ----

func on_awaken_boss_spawn(data: Dictionary) -> void:
	_instance_id = data.get("instance_id", "")
	_boss_name = data.get("boss_name", "覺醒 BOSS")
	_boss_icon = data.get("boss_icon", "🐉")
	_remaining_seconds = data.get("duration_seconds", 30.0)
	_hit_count = 0
	_powerup_progress = 0.0
	_powerup_threshold = data.get("powerup_threshold", 5)
	_active = true

	var color_str = data.get("boss_color", "#FF4500")
	_boss_color = Color(color_str)

	var min_mult = data.get("min_mult", 90.0)
	var max_mult = data.get("max_mult", 200.0)
	var pu_min = data.get("powerup_min_mult", 6.0)
	var pu_max = data.get("powerup_max_mult", 10.0)

	_boss_icon_label.text = _boss_icon
	_boss_name_label.text = _boss_name
	_mult_range_label.text = "%.0fx~%.0fx (⚡×%.0f-%.0f)" % [min_mult, max_mult, pu_min, pu_max]
	_hit_count_label.text = "命中 0 次"
	_powerup_label.text = "⚡ Power Up (0/%d)" % _powerup_threshold
	_update_powerup_bar()

	# 更新邊框顏色
	var status_style = StyleBoxFlat.new()
	status_style.bg_color = Color(0.1, 0.02, 0.0, 0.9)
	status_style.border_color = _boss_color
	status_style.set_border_width_all(2)
	status_style.set_corner_radius_all(8)
	_status_panel.add_theme_stylebox_override("panel", status_style)

	_status_panel.visible = true
	_status_panel.modulate.a = 0.0
	if _status_tween:
		_status_tween.kill()
	_status_tween = create_tween()
	_status_tween.tween_property(_status_panel, "modulate:a", 1.0, 0.4)

	_show_banner(data.get("message", "%s %s 覺醒！" % [_boss_icon, _boss_name]), _boss_color, true)
	_flash_powerup(Color(1.0, 0.3, 0.0, 0.4))

func on_awaken_boss_hit(data: Dictionary) -> void:
	if data.get("instance_id", "") != _instance_id:
		return

	_hit_count = data.get("hit_count", _hit_count)
	_powerup_progress = data.get("powerup_progress", _powerup_progress)

	var hits_since = int(_powerup_progress * _powerup_threshold)
	_hit_count_label.text = "命中 %d 次" % _hit_count
	_powerup_label.text = "⚡ Power Up (%d/%d)" % [hits_since, _powerup_threshold]
	_update_powerup_bar()

func on_awaken_boss_powerup(data: Dictionary) -> void:
	if data.get("instance_id", "") != _instance_id:
		return

	var mult = data.get("multiplier", 0.0)
	var player_name = data.get("player_name", "")
	var powerup_count = data.get("powerup_count", 1)

	# 重置進度條
	_powerup_progress = 0.0
	_powerup_label.text = "⚡ Power Up (0/%d)" % _powerup_threshold
	_update_powerup_bar()

	# Power Up 爆發效果
	_flash_powerup(Color(1.0, 0.5, 0.0, 0.6))

	# 顯示橫幅
	_show_banner("⚡ %s Power Up #%d！%.0fx！" % [player_name, powerup_count, mult], Color(1.0, 0.7, 0.0), false)

func on_awaken_boss_leave(data: Dictionary) -> void:
	if data.get("instance_id", "") != _instance_id:
		return

	_active = false
	var msg = data.get("message", "%s %s 離去！" % [_boss_icon, _boss_name])
	_show_banner(msg, Color(0.6, 0.6, 0.6), false)

	if _status_tween:
		_status_tween.kill()
	_status_tween = create_tween()
	_status_tween.tween_property(_status_panel, "modulate:a", 0.0, 0.6)
	_status_tween.tween_callback(func(): _status_panel.visible = false)

func on_awaken_boss_status(data: Dictionary) -> void:
	if not data.get("active", false):
		return
	_instance_id = data.get("instance_id", "")
	_boss_name = data.get("boss_name", "覺醒 BOSS")
	_boss_icon = data.get("boss_icon", "🐉")
	_remaining_seconds = data.get("remaining_seconds", 0.0)
	_hit_count = data.get("hit_count", 0)
	_powerup_progress = data.get("powerup_progress", 0.0)
	_powerup_threshold = data.get("powerup_threshold", 5)
	_active = _remaining_seconds > 0

	if _active:
		var min_mult = data.get("min_mult", 90.0)
		var max_mult = data.get("max_mult", 200.0)
		_boss_icon_label.text = _boss_icon
		_boss_name_label.text = _boss_name
		_mult_range_label.text = "%.0fx ~ %.0fx" % [min_mult, max_mult]
		_hit_count_label.text = "命中 %d 次" % _hit_count
		var hits_since = int(_powerup_progress * _powerup_threshold)
		_powerup_label.text = "⚡ Power Up (%d/%d)" % [hits_since, _powerup_threshold]
		_update_powerup_bar()
		_status_panel.visible = true
		_status_panel.modulate.a = 1.0

# ---- 私有方法 ----

func _show_banner(text: String, color: Color, is_spawn: bool) -> void:
	_banner_label.text = text
	_banner_label.add_theme_color_override("font_color", color)

	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.15, 0.03, 0.0, 0.93)
	banner_style.border_color = color
	banner_style.set_border_width_all(2)
	banner_style.set_corner_radius_all(6)
	_banner.add_theme_stylebox_override("panel", banner_style)

	_banner.visible = true
	if _banner_tween:
		_banner_tween.kill()
	_banner_tween = create_tween()
	_banner.position.y = -65
	_banner_tween.tween_property(_banner, "position:y", 0.0, 0.3).set_ease(Tween.EASE_OUT)
	_banner_tween.tween_interval(BANNER_DURATION)
	_banner_tween.tween_property(_banner, "position:y", -65.0, 0.3).set_ease(Tween.EASE_IN)
	_banner_tween.tween_callback(func(): _banner.visible = false)

	if is_spawn:
		_banner_tween.parallel().tween_property(_banner, "modulate:a", 0.3, 0.1)
		_banner_tween.parallel().tween_property(_banner, "modulate:a", 1.0, 0.1)

func _flash_powerup(color: Color) -> void:
	_powerup_flash.color = color
	_powerup_flash.visible = true
	var tween = create_tween()
	tween.tween_property(_powerup_flash, "color:a", 0.0, POWERUP_FLASH_DURATION)
	tween.tween_callback(func(): _powerup_flash.visible = false)
