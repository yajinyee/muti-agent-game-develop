## MegaCatchPanel.gd
## 全服 Mega Catch 事件面板（DAY-140）
## 業界依據：Ocean King 系列「Mega Catch」— 全場高倍率目標湧現 + 獎勵翻倍
## BOSS 擊殺後 60% 機率觸發，或每分鐘 5% 機率隨機觸發

extends Control

# ---- 狀態 ----
var _is_active: bool = false
var _seconds_left: float = 0.0
var _total_duration: float = 12.0
var _reward_boost: float = 1.5
var _tier_color: Color = Color(0.4, 0.8, 1.0, 1.0)

# ---- 節點 ----
var _banner: Control = null
var _banner_label: Label = null
var _timer_label: Label = null
var _boost_label: Label = null
var _flash_overlay: ColorRect = null
var _progress_bar: Control = null
var _progress_fill: ColorRect = null

# ---- 顏色 ----
const COLOR_MEGA    = Color(0.4, 0.8, 1.0, 1.0)   # 大豐收藍
const COLOR_SUPER   = Color(1.0, 0.84, 0.0, 1.0)  # 超級豐收金
const COLOR_LEGEND  = Color(1.0, 0.41, 0.71, 1.0) # 傳說豐收粉

func _ready() -> void:
	_build_ui()
	set_process(false)
	visible = false
	mouse_filter = Control.MOUSE_FILTER_IGNORE

func _build_ui() -> void:
	# 全螢幕閃光層
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(0.4, 0.8, 1.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner = Control.new()
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 56)
	_banner.position = Vector2(0, -60)
	_banner.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_banner)

	var banner_bg = ColorRect.new()
	banner_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	banner_bg.color = Color(0.02, 0.08, 0.18, 0.96)
	_banner.add_child(banner_bg)

	# 主標題
	_banner_label = Label.new()
	_banner_label.text = "🎣 大豐收！全場獎勵 ×1.5！"
	_banner_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", COLOR_MEGA)
	_banner_label.add_theme_font_size_override("font_size", 20)
	_banner.add_child(_banner_label)

	# 倒數計時（右側）
	_timer_label = Label.new()
	_timer_label.text = "12.0s"
	_timer_label.position = Vector2(-80, 0)
	_timer_label.size = Vector2(76, 56)
	_timer_label.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	_timer_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_timer_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_timer_label.add_theme_font_size_override("font_size", 14)
	_banner.add_child(_timer_label)

	# 進度條（橫幅底部）
	_progress_bar = Control.new()
	_progress_bar.position = Vector2(0, 50)
	_progress_bar.size = Vector2(1280, 6)
	_banner.add_child(_progress_bar)

	var bar_bg = ColorRect.new()
	bar_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	bar_bg.color = Color(0.1, 0.1, 0.2, 1.0)
	_progress_bar.add_child(bar_bg)

	_progress_fill = ColorRect.new()
	_progress_fill.color = COLOR_MEGA
	_progress_fill.position = Vector2(0, 0)
	_progress_fill.size = Vector2(1280, 6)
	_progress_bar.add_child(_progress_fill)

func _process(delta: float) -> void:
	if not _is_active:
		return

	_seconds_left -= delta
	if _seconds_left <= 0:
		_seconds_left = 0
		_is_active = false
		set_process(false)
		return

	# 更新倒數計時
	_timer_label.text = "%.1fs" % _seconds_left
	if _seconds_left < 3.0:
		_timer_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3, 1.0))

	# 更新進度條
	var pct = _seconds_left / _total_duration
	if _progress_fill:
		_progress_fill.size = Vector2(1280.0 * pct, 6)

# ---- 公開 API ----

## on_mega_catch_start 收到 mega_catch_start 時呼叫
func on_mega_catch_start(data: Dictionary) -> void:
	var tier_name: String = data.get("tier_name", "🎣 大豐收")
	var tier_icon: String = data.get("tier_icon", "🎣")
	var reward_boost: float = data.get("reward_boost", 1.5)
	var duration: float = data.get("duration", 12.0)
	var seconds_left: float = data.get("seconds_left", duration)

	_reward_boost = reward_boost
	_total_duration = duration
	_seconds_left = seconds_left
	_is_active = true

	# 決定顏色
	if reward_boost >= 3.0:
		_tier_color = COLOR_LEGEND
	elif reward_boost >= 2.0:
		_tier_color = COLOR_SUPER
	else:
		_tier_color = COLOR_MEGA

	# 更新 UI
	_banner_label.text = "%s %s！全場獎勵 ×%.0f！" % [tier_icon, tier_name, reward_boost]
	_banner_label.add_theme_color_override("font_color", _tier_color)
	_progress_fill.color = _tier_color
	_timer_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))

	# 顯示橫幅（從頂部滑入）
	visible = true
	_banner.position = Vector2(0, -60)
	var tween = create_tween()
	tween.tween_property(_banner, "position", Vector2(0, 0), 0.3).set_ease(Tween.EASE_OUT)

	# 全螢幕閃光
	_flash_overlay.color = Color(_tier_color.r, _tier_color.g, _tier_color.b, 0.45)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color",
		Color(_tier_color.r, _tier_color.g, _tier_color.b, 0.0), 0.5)

	set_process(true)

## on_mega_catch_end 收到 mega_catch_end 時呼叫
func on_mega_catch_end(_data: Dictionary) -> void:
	_is_active = false
	set_process(false)

	# 橫幅滑出
	var tween = create_tween()
	tween.tween_property(_banner, "position", Vector2(0, -60), 0.3).set_ease(Tween.EASE_IN)
	tween.tween_callback(func(): visible = false)
