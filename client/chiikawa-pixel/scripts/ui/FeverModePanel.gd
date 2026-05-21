## FeverModePanel.gd — 狂熱模式面板（DAY-133）
## 業界依據：Fire Kirin / Ocean King 系列的 Fever Mode
## 顯示狂熱模式觸發進度條、狂熱中的倒數計時、全服廣播橫幅
extends Control

# ---- 節點引用 ----
var _status_panel: PanelContainer   # 左下角狀態面板
var _progress_label: Label          # 觸發進度（3/5）
var _progress_bar_bg: ColorRect
var _progress_bar_fill: ColorRect
var _timer_label: Label             # 狂熱倒數計時
var _mult_label: Label              # 倍率顯示（×1.5）
var _fever_banner: PanelContainer   # 頂部橫幅（全服廣播）
var _banner_label: Label
var _banner_sub_label: Label
var _flash_overlay: ColorRect       # 全螢幕橙紅閃光

# ---- 狀態 ----
var _is_active: bool = false
var _seconds_left: float = 0.0
var _kill_progress: int = 0
var _trigger_kills: int = 5
var _mult_boost: float = 1.5
var _my_player_id: String = ""
var _banner_tween: Tween
var _flash_tween: Tween

const BANNER_DURATION := 4.0
const FLASH_DURATION := 0.3
const PANEL_WIDTH := 130.0
const PANEL_HEIGHT := 90.0

func _ready() -> void:
	_build_ui()
	_status_panel.visible = false
	_fever_banner.visible = false
	_flash_overlay.visible = false

func _process(delta: float) -> void:
	if _is_active and _seconds_left > 0:
		_seconds_left -= delta
		if _seconds_left <= 0:
			_seconds_left = 0
			_is_active = false
			_update_status_ui()
		else:
			_timer_label.text = "🔥 %.1fs" % _seconds_left
			# 最後 5 秒紅色閃爍
			if _seconds_left <= 5.0:
				var blink = int(_seconds_left * 4) % 2 == 0
				_timer_label.add_theme_color_override("font_color",
					Color.RED if blink else Color.WHITE)
			else:
				_timer_label.add_theme_color_override("font_color", Color(1.0, 0.6, 0.1))

func _build_ui() -> void:
	# 全螢幕橙紅閃光
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.3, 0.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅（全服廣播）
	_fever_banner = PanelContainer.new()
	_fever_banner.position = Vector2(0, -80)
	_fever_banner.custom_minimum_size = Vector2(1280, 70)
	add_child(_fever_banner)

	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.15, 0.03, 0.0, 0.92)
	banner_style.border_color = Color(1.0, 0.4, 0.0, 1.0)
	banner_style.set_border_width_all(2)
	_fever_banner.add_theme_stylebox_override("panel", banner_style)

	var banner_vbox = VBoxContainer.new()
	banner_vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	banner_vbox.add_theme_constant_override("separation", 2)
	_fever_banner.add_child(banner_vbox)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.add_theme_font_size_override("font_size", 22)
	_banner_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.1))
	_banner_label.text = "🔥 狂熱模式！"
	banner_vbox.add_child(_banner_label)

	_banner_sub_label = Label.new()
	_banner_sub_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_sub_label.add_theme_font_size_override("font_size", 14)
	_banner_sub_label.add_theme_color_override("font_color", Color(0.9, 0.9, 1.0))
	_banner_sub_label.text = ""
	banner_vbox.add_child(_banner_sub_label)

	# 左下角狀態面板（在連勝面板下方）
	_status_panel = PanelContainer.new()
	_status_panel.position = Vector2(8, 720 - 300)
	_status_panel.custom_minimum_size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	add_child(_status_panel)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.03, 0.0, 0.85)
	style.border_color = Color(1.0, 0.4, 0.0, 0.8)
	style.set_border_width_all(1)
	style.set_corner_radius_all(6)
	_status_panel.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 4)
	_status_panel.add_child(vbox)

	# 進度標籤
	_progress_label = Label.new()
	_progress_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_progress_label.add_theme_font_size_override("font_size", 12)
	_progress_label.add_theme_color_override("font_color", Color(1.0, 0.7, 0.3))
	_progress_label.text = "🔥 0/5"
	vbox.add_child(_progress_label)

	# 進度條
	var bar_container = Control.new()
	bar_container.custom_minimum_size = Vector2(110, 8)
	vbox.add_child(bar_container)

	_progress_bar_bg = ColorRect.new()
	_progress_bar_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	_progress_bar_bg.color = Color(0.2, 0.1, 0.0, 0.8)
	bar_container.add_child(_progress_bar_bg)

	_progress_bar_fill = ColorRect.new()
	_progress_bar_fill.position = Vector2(0, 0)
	_progress_bar_fill.size = Vector2(0, 8)
	_progress_bar_fill.color = Color(1.0, 0.4, 0.0)
	bar_container.add_child(_progress_bar_fill)

	# 倒數計時（狂熱中顯示）
	_timer_label = Label.new()
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_timer_label.add_theme_font_size_override("font_size", 18)
	_timer_label.add_theme_color_override("font_color", Color(1.0, 0.6, 0.1))
	_timer_label.text = ""
	vbox.add_child(_timer_label)

	# 倍率顯示
	_mult_label = Label.new()
	_mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_mult_label.add_theme_font_size_override("font_size", 11)
	_mult_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.3))
	_mult_label.text = ""
	vbox.add_child(_mult_label)

# ---- 公開 API ----

func set_player_id(pid: String) -> void:
	_my_player_id = pid

## 狂熱模式開始（全服廣播）
func on_fever_start(data: Dictionary) -> void:
	var player_id: String = data.get("player_id", "")
	var player_name: String = data.get("player_name", "???")
	var seconds_left: int = data.get("seconds_left", 15)
	var mult_boost: float = data.get("mult_boost", 1.5)
	var is_self: bool = (player_id == _my_player_id)

	# 全螢幕橙紅閃光
	_play_flash(is_self)

	# 頂部橫幅
	_banner_label.text = "🔥 %s 進入狂熱模式！" % player_name
	_banner_sub_label.text = "獎勵 ×%.1f！持續 %d 秒！" % [mult_boost, seconds_left]
	if is_self:
		_banner_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.1))
	else:
		_banner_label.add_theme_color_override("font_color", Color(0.9, 0.9, 1.0))
	_show_banner()

	# 自己觸發時更新狀態面板
	if is_self:
		_is_active = true
		_seconds_left = float(seconds_left)
		_mult_boost = mult_boost
		_update_status_ui()

## 狂熱模式結束（個人）
func on_fever_end(data: Dictionary) -> void:
	_is_active = false
	_seconds_left = 0.0
	_update_status_ui()

## 狂熱模式狀態更新（個人）
func on_fever_status(data: Dictionary) -> void:
	var is_active: bool = data.get("is_active", false)
	var seconds_left: int = data.get("seconds_left", 0)
	var kill_progress: int = data.get("kill_progress", 0)
	var trigger_kills: int = data.get("trigger_kills", 5)
	var mult_boost: float = data.get("mult_boost", 1.5)

	_is_active = is_active
	_seconds_left = float(seconds_left)
	_kill_progress = kill_progress
	_trigger_kills = trigger_kills
	_mult_boost = mult_boost
	_update_status_ui()

# ---- 私有方法 ----

func _update_status_ui() -> void:
	if _is_active:
		_status_panel.visible = true
		_progress_label.text = "🔥 狂熱中！"
		_progress_bar_fill.size.x = 110.0  # 滿格
		_progress_bar_fill.color = Color(1.0, 0.3, 0.0)
		_timer_label.text = "🔥 %.1fs" % _seconds_left
		_mult_label.text = "×%.1f 獎勵！" % _mult_boost
	else:
		if _kill_progress > 0:
			_status_panel.visible = true
			_progress_label.text = "🔥 %d/%d" % [_kill_progress, _trigger_kills]
			var ratio = float(_kill_progress) / float(_trigger_kills)
			_progress_bar_fill.size.x = 110.0 * ratio
			_progress_bar_fill.color = Color(1.0, 0.4 + ratio * 0.3, 0.0)
			_timer_label.text = ""
			_mult_label.text = ""
		else:
			_status_panel.visible = false
			_timer_label.text = ""
			_mult_label.text = ""

func _play_flash(is_self: bool) -> void:
	if _flash_tween:
		_flash_tween.kill()
	_flash_overlay.visible = true
	var alpha = 0.5 if is_self else 0.2
	_flash_overlay.color = Color(1.0, 0.3, 0.0, alpha)
	_flash_tween = create_tween()
	_flash_tween.tween_property(_flash_overlay, "color:a", 0.0, FLASH_DURATION)
	_flash_tween.tween_callback(func(): _flash_overlay.visible = false)

func _show_banner() -> void:
	if _banner_tween:
		_banner_tween.kill()
	_fever_banner.visible = true
	_fever_banner.position.y = -80
	_banner_tween = create_tween()
	_banner_tween.tween_property(_fever_banner, "position:y", 0.0, 0.3).set_ease(Tween.EASE_OUT)
	_banner_tween.tween_interval(BANNER_DURATION)
	_banner_tween.tween_property(_fever_banner, "position:y", -80.0, 0.3).set_ease(Tween.EASE_IN)
	_banner_tween.tween_callback(func(): _fever_banner.visible = false)
