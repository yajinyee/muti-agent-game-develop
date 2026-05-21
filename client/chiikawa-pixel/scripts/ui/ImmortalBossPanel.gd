## ImmortalBossPanel.gd — 不死 BOSS 連勝面板（DAY-129）
## 業界依據：JILI Royal Fishing 2026 Immortal Boss
## 顯示不死 BOSS 的出現、命中獎勵、倒數計時和離開動畫
extends Control

# ---- 節點引用 ----
var _banner: PanelContainer       # 頂部橫幅（出現/離開通知）
var _banner_label: Label          # 橫幅文字
var _status_panel: PanelContainer # 右側狀態面板（BOSS 活躍時顯示）
var _boss_icon_label: Label       # BOSS 圖示
var _boss_name_label: Label       # BOSS 名稱
var _timer_label: Label           # 倒數計時
var _hit_count_label: Label       # 命中次數
var _mult_range_label: Label      # 倍率範圍
var _hit_flash: ColorRect         # 命中閃光

# ---- 狀態 ----
var _active: bool = false
var _instance_id: String = ""
var _boss_name: String = ""
var _boss_icon: String = ""
var _boss_color: Color = Color.GOLD
var _remaining_seconds: float = 0.0
var _hit_count: int = 0
var _total_reward: int = 0
var _banner_tween: Tween
var _status_tween: Tween

# ---- 常數 ----
const BANNER_DURATION := 4.0   # 橫幅顯示時間（秒）
const HIT_FLASH_DURATION := 0.3 # 命中閃光時間

func _ready() -> void:
	_build_ui()
	_status_panel.visible = false
	_banner.visible = false
	_hit_flash.visible = false

func _build_ui() -> void:
	# 頂部橫幅（出現/離開通知）
	_banner = PanelContainer.new()
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 56)
	_banner.position = Vector2(0, -60)
	add_child(_banner)

	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.1, 0.05, 0.0, 0.92)
	banner_style.border_color = Color.GOLD
	banner_style.set_border_width_all(2)
	banner_style.set_corner_radius_all(6)
	_banner.add_theme_stylebox_override("panel", banner_style)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_font_size_override("font_size", 18)
	_banner_label.add_theme_color_override("font_color", Color.GOLD)
	_banner.add_child(_banner_label)

	# 右側狀態面板（BOSS 活躍時顯示）
	_status_panel = PanelContainer.new()
	_status_panel.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_status_panel.anchor_left = 1.0
	_status_panel.anchor_right = 1.0
	_status_panel.anchor_top = 0.5
	_status_panel.anchor_bottom = 0.5
	_status_panel.offset_left = -180
	_status_panel.offset_right = -8
	_status_panel.offset_top = -80
	_status_panel.offset_bottom = 80
	add_child(_status_panel)

	var status_style = StyleBoxFlat.new()
	status_style.bg_color = Color(0.08, 0.04, 0.0, 0.88)
	status_style.border_color = Color.GOLD
	status_style.set_border_width_all(2)
	status_style.set_corner_radius_all(8)
	_status_panel.add_theme_stylebox_override("panel", status_style)

	var vbox = VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	vbox.add_theme_constant_override("separation", 4)
	_status_panel.add_child(vbox)

	# BOSS 圖示
	_boss_icon_label = Label.new()
	_boss_icon_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_boss_icon_label.add_theme_font_size_override("font_size", 32)
	vbox.add_child(_boss_icon_label)

	# BOSS 名稱
	_boss_name_label = Label.new()
	_boss_name_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_boss_name_label.add_theme_font_size_override("font_size", 14)
	_boss_name_label.add_theme_color_override("font_color", Color.GOLD)
	vbox.add_child(_boss_name_label)

	# 倍率範圍
	_mult_range_label = Label.new()
	_mult_range_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_mult_range_label.add_theme_font_size_override("font_size", 12)
	_mult_range_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.3))
	vbox.add_child(_mult_range_label)

	# 倒數計時
	_timer_label = Label.new()
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_timer_label.add_theme_font_size_override("font_size", 20)
	_timer_label.add_theme_color_override("font_color", Color.WHITE)
	vbox.add_child(_timer_label)

	# 命中次數
	_hit_count_label = Label.new()
	_hit_count_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_hit_count_label.add_theme_font_size_override("font_size", 12)
	_hit_count_label.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
	vbox.add_child(_hit_count_label)

	# 命中閃光（全螢幕金色閃光）
	_hit_flash = ColorRect.new()
	_hit_flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	_hit_flash.color = Color(1.0, 0.85, 0.0, 0.0)
	_hit_flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_hit_flash)

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
	# 最後 5 秒變紅色閃爍
	if _remaining_seconds <= 5.0:
		_timer_label.add_theme_color_override("font_color", Color.RED)
		_timer_label.modulate.a = 0.5 + 0.5 * sin(Time.get_ticks_msec() * 0.01)
	else:
		_timer_label.add_theme_color_override("font_color", Color.WHITE)
		_timer_label.modulate.a = 1.0

# ---- 外部呼叫 ----

## on_immortal_boss_spawn — 不死 BOSS 出現
func on_immortal_boss_spawn(data: Dictionary) -> void:
	_instance_id = data.get("instance_id", "")
	_boss_name = data.get("boss_name", "不死 BOSS")
	_boss_icon = data.get("boss_icon", "👾")
	_remaining_seconds = data.get("duration_seconds", 25.0)
	_hit_count = 0
	_total_reward = 0
	_active = true

	var color_str = data.get("boss_color", "#FFD700")
	_boss_color = Color(color_str)

	var min_mult = data.get("min_mult", 50.0)
	var max_mult = data.get("max_mult", 150.0)

	# 更新狀態面板
	_boss_icon_label.text = _boss_icon
	_boss_name_label.text = _boss_name
	_mult_range_label.text = "%.0fx ~ %.0fx" % [min_mult, max_mult]
	_hit_count_label.text = "命中 0 次"

	# 更新邊框顏色
	var status_style = StyleBoxFlat.new()
	status_style.bg_color = Color(0.08, 0.04, 0.0, 0.88)
	status_style.border_color = _boss_color
	status_style.set_border_width_all(2)
	status_style.set_corner_radius_all(8)
	_status_panel.add_theme_stylebox_override("panel", status_style)

	# 顯示狀態面板（滑入）
	_status_panel.visible = true
	_status_panel.modulate.a = 0.0
	if _status_tween:
		_status_tween.kill()
	_status_tween = create_tween()
	_status_tween.tween_property(_status_panel, "modulate:a", 1.0, 0.4)

	# 顯示橫幅
	var msg = data.get("message", "%s %s 出現了！" % [_boss_icon, _boss_name])
	_show_banner(msg, _boss_color, true)

	# 全螢幕金色閃光
	_flash_screen(_boss_color)

## on_immortal_boss_hit — 命中不死 BOSS
func on_immortal_boss_hit(data: Dictionary) -> void:
	if data.get("instance_id", "") != _instance_id:
		return

	_hit_count = data.get("hit_count", _hit_count)
	_total_reward = data.get("total_reward", _total_reward)
	var mult = data.get("multiplier", 0.0)
	var reward = data.get("reward", 0)
	var is_high_mult = data.get("is_high_mult", false)
	var player_name = data.get("player_name", "")

	# 更新命中次數
	_hit_count_label.text = "命中 %d 次" % _hit_count

	# 命中閃光
	var flash_color = Color(1.0, 0.85, 0.0, 0.35) if is_high_mult else Color(1.0, 0.85, 0.0, 0.15)
	_flash_screen(flash_color)

	# 高倍率時顯示橫幅
	if is_high_mult:
		_show_banner("🌟 %s 命中 %.0fx！+%d 金幣" % [player_name, mult, reward], Color.GOLD, false)

## on_immortal_boss_leave — 不死 BOSS 離開
func on_immortal_boss_leave(data: Dictionary) -> void:
	if data.get("instance_id", "") != _instance_id:
		return

	_active = false
	var msg = data.get("message", "%s %s 離開了！" % [_boss_icon, _boss_name])
	_show_banner(msg, Color(0.7, 0.7, 0.7), false)

	# 狀態面板淡出
	if _status_tween:
		_status_tween.kill()
	_status_tween = create_tween()
	_status_tween.tween_property(_status_panel, "modulate:a", 0.0, 0.6)
	_status_tween.tween_callback(func(): _status_panel.visible = false)

## on_immortal_boss_status — 登入時恢復狀態
func on_immortal_boss_status(data: Dictionary) -> void:
	if not data.get("active", false):
		return
	_instance_id = data.get("instance_id", "")
	_boss_name = data.get("boss_name", "不死 BOSS")
	_boss_icon = data.get("boss_icon", "👾")
	_remaining_seconds = data.get("remaining_seconds", 0.0)
	_hit_count = data.get("hit_count", 0)
	_total_reward = data.get("total_reward", 0)
	_active = _remaining_seconds > 0

	if _active:
		var min_mult = data.get("min_mult", 50.0)
		var max_mult = data.get("max_mult", 150.0)
		_boss_icon_label.text = _boss_icon
		_boss_name_label.text = _boss_name
		_mult_range_label.text = "%.0fx ~ %.0fx" % [min_mult, max_mult]
		_hit_count_label.text = "命中 %d 次" % _hit_count
		_status_panel.visible = true
		_status_panel.modulate.a = 1.0

# ---- 私有方法 ----

func _show_banner(text: String, color: Color, is_spawn: bool) -> void:
	_banner_label.text = text
	_banner_label.add_theme_color_override("font_color", color)

	# 更新橫幅邊框顏色
	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.1, 0.05, 0.0, 0.92)
	banner_style.border_color = color
	banner_style.set_border_width_all(2)
	banner_style.set_corner_radius_all(6)
	_banner.add_theme_stylebox_override("panel", banner_style)

	_banner.visible = true
	if _banner_tween:
		_banner_tween.kill()
	_banner_tween = create_tween()

	# 滑入
	_banner.position.y = -60
	_banner_tween.tween_property(_banner, "position:y", 0.0, 0.3).set_ease(Tween.EASE_OUT)
	# 停留
	_banner_tween.tween_interval(BANNER_DURATION)
	# 滑出
	_banner_tween.tween_property(_banner, "position:y", -60.0, 0.3).set_ease(Tween.EASE_IN)
	_banner_tween.tween_callback(func(): _banner.visible = false)

	# 出現時額外閃爍
	if is_spawn:
		_banner_tween.parallel().tween_property(_banner, "modulate:a", 0.3, 0.1)
		_banner_tween.parallel().tween_property(_banner, "modulate:a", 1.0, 0.1)

func _flash_screen(color: Color) -> void:
	_hit_flash.color = color
	_hit_flash.visible = true
	var tween = create_tween()
	tween.tween_property(_hit_flash, "color:a", 0.0, HIT_FLASH_DURATION)
	tween.tween_callback(func(): _hit_flash.visible = false)
