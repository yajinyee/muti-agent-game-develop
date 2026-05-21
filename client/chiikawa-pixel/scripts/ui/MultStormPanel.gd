## MultStormPanel.gd
## 全服倍率風暴面板（DAY-138）
## 業界依據：findingdulcinea.com 2026「admin events multiply luck by thousands」
## 風暴觸發時全螢幕特效 + 頂部橫幅 + 倒數計時，製造「全場瘋狂」高峰體驗

extends Control

# ---- 狀態 ----
var _is_active: bool = false
var _seconds_left: float = 0.0
var _mult_boost: float = 1.0
var _tier_color: Color = Color.WHITE

# ---- 節點 ----
var _banner: Control = null
var _banner_label: Label = null
var _timer_label: Label = null
var _mult_label: Label = null
var _flash_overlay: ColorRect = null
var _particle_overlay: Control = null  # 粒子效果層

# ---- 顏色 ----
const COLOR_LIGHTNING = Color(1.0, 0.88, 0.2, 1.0)  # 閃電黃
const COLOR_TSUNAMI   = Color(0.29, 0.56, 0.85, 1.0) # 海嘯藍
const COLOR_RAINBOW   = Color(1.0, 0.41, 0.71, 1.0)  # 彩虹粉

func _ready() -> void:
	_build_ui()
	set_process(false)
	mouse_filter = Control.MOUSE_FILTER_IGNORE

func _build_ui() -> void:
	# 全螢幕閃光層
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.88, 0.2, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅（風暴進行中）
	_banner = Control.new()
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 52)
	_banner.position = Vector2(0, -56)  # 初始在畫面外
	_banner.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_banner)

	var banner_bg = ColorRect.new()
	banner_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	banner_bg.color = Color(0.05, 0.05, 0.1, 0.95)
	banner_bg.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_banner.add_child(banner_bg)

	# 左側倍率顯示
	_mult_label = Label.new()
	_mult_label.position = Vector2(10, 6)
	_mult_label.custom_minimum_size = Vector2(80, 40)
	_mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_mult_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_mult_label.add_theme_font_size_override("font_size", 22)
	_mult_label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_banner.add_child(_mult_label)

	# 中央橫幅文字
	_banner_label = Label.new()
	_banner_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_font_size_override("font_size", 15)
	_banner_label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_banner.add_child(_banner_label)

	# 右側倒數計時
	_timer_label = Label.new()
	_timer_label.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_timer_label.position = Vector2(-75, -14)
	_timer_label.custom_minimum_size = Vector2(65, 28)
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	_timer_label.add_theme_color_override("font_color", Color.WHITE)
	_timer_label.add_theme_font_size_override("font_size", 18)
	_timer_label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_banner.add_child(_timer_label)

func _process(delta: float) -> void:
	if not _is_active:
		return
	_seconds_left -= delta
	if _seconds_left <= 0:
		_seconds_left = 0
		set_process(false)
		return
	_update_timer()

func _update_timer() -> void:
	if not is_instance_valid(_timer_label):
		return
	var secs = int(_seconds_left)
	_timer_label.text = "%ds" % secs
	# 最後 5 秒紅色閃爍
	if _seconds_left <= 5:
		var blink = fmod(_seconds_left, 0.4) > 0.2
		_timer_label.add_theme_color_override("font_color", Color.RED if blink else Color.WHITE)
	else:
		_timer_label.add_theme_color_override("font_color", Color.WHITE)

# ---- 外部呼叫 ----

## 風暴開始（全服廣播）
func on_storm_start(data: Dictionary) -> void:
	var tier_name: String = data.get("tier_name", "⚡ 倍率風暴")
	var tier_icon: String = data.get("tier_icon", "⚡")
	var tier_color_hex: String = data.get("tier_color", "#FFE066")
	_mult_boost = data.get("mult_boost", 2.0)
	_seconds_left = data.get("seconds_left", 20.0)
	_is_active = true

	# 解析顏色
	_tier_color = Color.html(tier_color_hex) if tier_color_hex.begins_with("#") else COLOR_LIGHTNING

	# 更新 UI
	if is_instance_valid(_mult_label):
		_mult_label.text = "×%.0f" % _mult_boost
		_mult_label.add_theme_color_override("font_color", _tier_color)

	if is_instance_valid(_banner_label):
		_banner_label.text = "%s 全場倍率 ×%.0f！" % [tier_name, _mult_boost]
		_banner_label.add_theme_color_override("font_color", _tier_color)

	# 橫幅滑入
	_show_banner()

	# 全螢幕閃光（依等級強度）
	var flash_alpha = 0.3 + (_mult_boost - 2.0) * 0.1  # 2x=0.3, 3x=0.4, 5x=0.6
	_flash_screen(Color(_tier_color.r, _tier_color.g, _tier_color.b, flash_alpha), 0.6)

	# 彩虹風暴：雙閃光
	if _mult_boost >= 5.0:
		var tween = create_tween()
		tween.tween_interval(0.4)
		tween.tween_callback(func():
			_flash_screen(Color(1.0, 0.41, 0.71, 0.5), 0.5)
		)

	# 開始倒數
	set_process(true)

## 風暴結束（全服廣播）
func on_storm_end(_data: Dictionary) -> void:
	_is_active = false
	set_process(false)

	if is_instance_valid(_banner_label):
		_banner_label.text = "倍率風暴結束，回歸正常！"
		_banner_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7, 1.0))

	if is_instance_valid(_mult_label):
		_mult_label.text = ""

	# 2 秒後橫幅滑出
	var tween = create_tween()
	tween.tween_interval(2.0)
	tween.tween_callback(_hide_banner)

# ---- 私有方法 ----

func _show_banner() -> void:
	if not is_instance_valid(_banner):
		return
	_banner.position = Vector2(0, -56)
	var tween = create_tween()
	tween.tween_property(_banner, "position", Vector2(0, 0), 0.35).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

func _hide_banner() -> void:
	if not is_instance_valid(_banner):
		return
	var tween = create_tween()
	tween.tween_property(_banner, "position", Vector2(0, -56), 0.25).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)

func _flash_screen(color: Color, duration: float) -> void:
	if not is_instance_valid(_flash_overlay):
		return
	_flash_overlay.color = color
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)
