## LuckyTeleportFishPanel.gd — 幸運傳送魚系統面板（DAY-223）
## 業界原創「傳送混亂」機制
##
## 視覺設計：
##   - 紫色傳送主題（#9B59B6 + #8E44AD + #D7BDE2 + #F5EEF8）
##   - teleport_start：紫色三次強閃光 + 頂部橫幅 + 計時條
##   - teleport_wave：全螢幕紫色閃光 + 「🌀 傳送！」大字 + 浮動位置標記 + 傳送混亂加成提示
##   - teleport_end：紫色淡出
extends CanvasLayer

# 傳送狀態
var _active: bool = false
var _banner: Control = null
var _wave_count: int = 0

# 主題顏色
const COLOR_PRIMARY   = Color("#9B59B6")  # 紫色
const COLOR_DARK      = Color("#8E44AD")  # 深紫
const COLOR_LIGHT     = Color("#D7BDE2")  # 淡紫
const COLOR_PALE      = Color("#F5EEF8")  # 極淡紫
const COLOR_BONUS     = Color("#FFD700")  # 金色（傳送混亂加成）
const COLOR_BG        = Color(0.06, 0.0, 0.1, 0.88)

func _ready() -> void:
	layer = 22  # 幸運傳送魚面板層級

## 處理幸運傳送魚訊息
func handle_lucky_teleport_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"teleport_start":
			_on_teleport_start(payload)
		"teleport_wave":
			_on_teleport_wave(payload)
		"teleport_end":
			_on_teleport_end(payload)

## teleport_start — 傳送漩渦開始
func _on_teleport_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var duration_sec: int = payload.get("duration_sec", 10)
	var max_waves: int = payload.get("max_waves", 4)

	_active = true
	_wave_count = 0

	# 紫色三次強閃光
	_flash_screen(COLOR_PRIMARY, 0.2)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_DARK, 0.18)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_LIGHT, 0.15)

	# 頂部橫幅 + 計時條
	_show_banner(
		"🌀 傳送漩渦！",
		"%s 觸發傳送漩渦！%d 次傳送！傳送後 3 秒內擊破獲得 ×2.5 倍率！" % [player_name, max_waves],
		duration_sec
	)

## teleport_wave — 傳送波次
func _on_teleport_wave(payload: Dictionary) -> void:
	var wave: int = payload.get("wave", 1)
	var max_waves: int = payload.get("max_waves", 4)
	var bonus_sec: int = payload.get("bonus_sec", 3)
	var bonus_mult: float = payload.get("bonus_mult", 2.5)
	var targets: Array = payload.get("targets", [])

	_wave_count = wave

	# 全螢幕紫色閃光
	_flash_screen(COLOR_PRIMARY, 0.15)

	# 「🌀 傳送！」大字
	var vp_size = get_viewport().size
	var big_label = Label.new()
	big_label.text = "🌀 傳送！（%d/%d）" % [wave, max_waves]
	big_label.add_theme_font_size_override("font_size", 44)
	big_label.add_theme_color_override("font_color", COLOR_LIGHT)
	big_label.position = vp_size / 2 - Vector2(130, 28)
	add_child(big_label)

	var tween_label = big_label.create_tween()
	tween_label.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.1)
	tween_label.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_label.tween_interval(0.4)
	tween_label.tween_property(big_label, "modulate:a", 0.0, 0.3)
	tween_label.tween_callback(big_label.queue_free)

	# 傳送混亂加成提示（金色）
	var bonus_label = Label.new()
	bonus_label.text = "⚡ ×%.1f 傳送混亂加成 %d 秒！" % [bonus_mult, bonus_sec]
	bonus_label.add_theme_font_size_override("font_size", 16)
	bonus_label.add_theme_color_override("font_color", COLOR_BONUS)
	bonus_label.position = vp_size / 2 - Vector2(110, -10)
	add_child(bonus_label)

	var tween_bonus = bonus_label.create_tween()
	tween_bonus.tween_property(bonus_label, "position:y", bonus_label.position.y - 20, 0.5)
	tween_bonus.parallel().tween_property(bonus_label, "modulate:a", 0.0, 0.5)
	tween_bonus.tween_callback(bonus_label.queue_free)

	# 傳送目標位置標記（紫色漩渦圓圈）
	for target_info in targets:
		var new_x: float = target_info.get("new_x", 0.0)
		var new_y: float = target_info.get("new_y", 0.0)
		_spawn_teleport_marker(Vector2(new_x, new_y))

## teleport_end — 傳送漩渦結束
func _on_teleport_end(_payload: Dictionary) -> void:
	_active = false
	_hide_banner()

	# 紫色淡出
	_flash_screen(COLOR_DARK, 0.3)

# ---- 輔助函數 ----

## 顯示頂部橫幅 + 計時條
func _show_banner(title: String, subtitle: String, duration_sec: int) -> void:
	_hide_banner()

	var banner = Control.new()
	banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	banner.position = Vector2(0, 8)
	banner.size = Vector2(get_viewport().size.x, 56)
	add_child(banner)

	var bg = ColorRect.new()
	bg.color = COLOR_BG
	bg.size = banner.size
	banner.add_child(bg)

	var title_label = Label.new()
	title_label.text = title
	title_label.add_theme_font_size_override("font_size", 20)
	title_label.add_theme_color_override("font_color", COLOR_LIGHT)
	title_label.position = Vector2(12, 4)
	banner.add_child(title_label)

	var sub_label = Label.new()
	sub_label.text = subtitle
	sub_label.add_theme_font_size_override("font_size", 11)
	sub_label.add_theme_color_override("font_color", COLOR_PALE)
	sub_label.position = Vector2(12, 28)
	banner.add_child(sub_label)

	# 計時條（底部，紫→深紫漸變）
	var timer_bg = ColorRect.new()
	timer_bg.color = Color(0.1, 0.0, 0.15, 0.8)
	timer_bg.position = Vector2(0, 50)
	timer_bg.size = Vector2(get_viewport().size.x, 6)
	banner.add_child(timer_bg)

	var timer_bar = ColorRect.new()
	timer_bar.color = COLOR_PRIMARY
	timer_bar.position = Vector2(0, 50)
	timer_bar.size = Vector2(get_viewport().size.x, 6)
	banner.add_child(timer_bar)

	var tween = banner.create_tween()
	tween.tween_property(timer_bar, "size:x", 0.0, float(duration_sec))
	tween.parallel().tween_property(timer_bar, "color", COLOR_DARK, float(duration_sec))

	_banner = banner

## 隱藏橫幅
func _hide_banner() -> void:
	if _banner != null and is_instance_valid(_banner):
		_banner.queue_free()
	_banner = null

## 傳送位置標記（紫色漩渦圓圈）
func _spawn_teleport_marker(pos: Vector2) -> void:
	var marker = Control.new()
	marker.position = pos - Vector2(16, 16)
	marker.size = Vector2(32, 32)
	add_child(marker)

	# 外圈（紫色）
	var outer = ColorRect.new()
	outer.color = Color(COLOR_PRIMARY.r, COLOR_PRIMARY.g, COLOR_PRIMARY.b, 0.6)
	outer.size = Vector2(32, 32)
	marker.add_child(outer)

	# 內圈（淡紫）
	var inner = ColorRect.new()
	inner.color = Color(COLOR_LIGHT.r, COLOR_LIGHT.g, COLOR_LIGHT.b, 0.4)
	inner.size = Vector2(20, 20)
	inner.position = Vector2(6, 6)
	marker.add_child(inner)

	# 漩渦動畫（旋轉 + 縮放 + 淡出）
	var tween = marker.create_tween()
	tween.tween_property(marker, "rotation_degrees", 360.0, 0.4)
	tween.parallel().tween_property(marker, "scale", Vector2(1.5, 1.5), 0.2)
	tween.tween_property(marker, "modulate:a", 0.0, 0.3)
	tween.tween_callback(marker.queue_free)

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.35)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)
