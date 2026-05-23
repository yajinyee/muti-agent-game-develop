## LuckyCloneFishPanel.gd — 幸運分身魚系統面板（DAY-242）
## 業界原創「三方向同時射擊」機制
##
## 視覺設計：
##   - 紫色分身主題（#8E44AD + #9B59B6 + #D7BDE2 + #F5EEF8）
##   - clone_start：紫色雙閃光 + 右側豎向計時條 + 「👥 分身模式！」大字 + 角度說明
##   - clone_broadcast：頂部小橫幅（通知全服有人觸發）
##   - clone_hit：左/右分身子彈命中閃光 + 方向箭頭 + 倍率浮動文字
##   - clone_end：紫色淡出
extends CanvasLayer

# 主題顏色
const COLOR_CLONE    = Color("#8E44AD")  # 紫色（主題）
const COLOR_LIGHT    = Color("#9B59B6")  # 淡紫
const COLOR_PALE     = Color("#F5EEF8")  # 極淡紫
const COLOR_GOLD     = Color("#FFD700")  # 金色（倍率）
const COLOR_LEFT     = Color("#3498DB")  # 藍色（左分身）
const COLOR_RIGHT    = Color("#E74C3C")  # 紅色（右分身）
const COLOR_WHITE    = Color("#FFFFFF")  # 白色

# 計時條節點
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null
var _duration_sec: int = 8

func _ready() -> void:
	layer = 3  # 幸運分身魚面板層級

## 處理幸運分身魚訊息
func handle_lucky_clone_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"clone_start":
			_on_clone_start(payload)
		"clone_broadcast":
			_on_clone_broadcast(payload)
		"clone_hit":
			_on_clone_hit(payload)
		"clone_end":
			_on_clone_end(payload)

## clone_start — 分身模式開始（個人訊息）
func _on_clone_start(payload: Dictionary) -> void:
	_duration_sec = payload.get("duration_sec", 8)
	var angle_deg: float = payload.get("angle_deg", 30.0)
	var kill_chance: float = payload.get("kill_chance", 0.6)
	var kill_mult: float = payload.get("kill_mult", 0.7)
	var vp_size = get_viewport().size

	# 紫色雙閃光
	_flash_screen(COLOR_CLONE, 0.13)
	await get_tree().create_timer(0.09).timeout
	_flash_screen(COLOR_PALE, 0.10)

	# 大字提示
	var big_label = Label.new()
	big_label.text = "👥 分身模式！"
	big_label.add_theme_font_size_override("font_size", 40)
	big_label.add_theme_color_override("font_color", COLOR_CLONE)
	big_label.position = Vector2(vp_size.x / 2 - 100, vp_size.y / 2 - 70)
	add_child(big_label)

	big_label.scale = Vector2(0.7, 0.7)
	var tw_big = big_label.create_tween()
	tw_big.tween_property(big_label, "scale", Vector2(1.1, 1.1), 0.15)
	tw_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tw_big.tween_interval(1.8)
	tw_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tw_big.tween_callback(big_label.queue_free)

	# 說明文字
	var info_label = Label.new()
	info_label.text = "每槍同時發射 3 個方向（±%.0f°）\n分身命中：%.0f%% 擊破 × %.1f 倍率" % [angle_deg, kill_chance * 100, kill_mult]
	info_label.add_theme_font_size_override("font_size", 13)
	info_label.add_theme_color_override("font_color", Color(0.85, 0.75, 0.95))
	info_label.position = Vector2(vp_size.x / 2 - 140, vp_size.y / 2 - 22)
	add_child(info_label)

	var tw_info = info_label.create_tween()
	tw_info.tween_interval(2.5)
	tw_info.tween_property(info_label, "modulate:a", 0.0, 0.4)
	tw_info.tween_callback(info_label.queue_free)

	# 方向示意圖（三條線）
	_draw_direction_hint(vp_size)

	# 右側豎向計時條
	_create_timer_bar(vp_size)

## 繪製方向示意圖
func _draw_direction_hint(vp_size: Vector2) -> void:
	var center_x = vp_size.x / 2
	var center_y = vp_size.y / 2 + 30
	var line_len = 60.0

	# 主方向（向右）
	_draw_arrow_line(center_x, center_y, center_x + line_len, center_y, COLOR_WHITE, 0.8)
	# 左偏 30 度
	var left_x = center_x + line_len * cos(deg_to_rad(30))
	var left_y = center_y - line_len * sin(deg_to_rad(30))
	_draw_arrow_line(center_x, center_y, left_x, left_y, COLOR_LEFT, 0.8)
	# 右偏 30 度
	var right_x = center_x + line_len * cos(deg_to_rad(-30))
	var right_y = center_y - line_len * sin(deg_to_rad(-30))
	_draw_arrow_line(center_x, center_y, right_x, right_y, COLOR_RIGHT, 0.8)

## 繪製箭頭線
func _draw_arrow_line(from_x: float, from_y: float, to_x: float, to_y: float, color: Color, duration: float) -> void:
	var dot_count = 5
	for i in range(dot_count):
		var t = float(i) / float(dot_count - 1)
		var px = from_x + (to_x - from_x) * t
		var py = from_y + (to_y - from_y) * t

		var dot = ColorRect.new()
		dot.color = color
		dot.size = Vector2(5, 5)
		dot.position = Vector2(px - 2.5, py - 2.5)
		add_child(dot)

		var tw = dot.create_tween()
		tw.tween_interval(duration)
		tw.tween_property(dot, "modulate:a", 0.0, 0.3)
		tw.tween_callback(dot.queue_free)

## 建立右側豎向計時條
func _create_timer_bar(vp_size: Vector2) -> void:
	if is_instance_valid(_timer_bar_bg):
		_timer_bar_bg.queue_free()
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()

	var bar_x = vp_size.x - 18
	var bar_h = 120.0

	# 背景
	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.15, 0.05, 0.2, 0.7)
	_timer_bar_bg.size = Vector2(10, bar_h)
	_timer_bar_bg.position = Vector2(bar_x, vp_size.y / 2 - bar_h / 2)
	add_child(_timer_bar_bg)

	# 計時條
	_timer_bar = ColorRect.new()
	_timer_bar.color = COLOR_CLONE
	_timer_bar.size = Vector2(10, bar_h)
	_timer_bar.position = Vector2(bar_x, vp_size.y / 2 - bar_h / 2)
	add_child(_timer_bar)

	# 計時條動畫（從上往下縮短）
	if is_instance_valid(_timer_tween):
		_timer_tween.kill()
	_timer_tween = _timer_bar.create_tween()
	_timer_tween.tween_property(_timer_bar, "size:y", 0.0, float(_duration_sec))
	_timer_tween.tween_callback(func():
		if is_instance_valid(_timer_bar_bg):
			_timer_bar_bg.queue_free()
		if is_instance_valid(_timer_bar):
			_timer_bar.queue_free()
	)

## clone_broadcast — 通知全服有人觸發分身模式
func _on_clone_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var duration_sec: int = payload.get("duration_sec", 8)
	var vp_size = get_viewport().size

	# 頂部小橫幅
	var banner = Label.new()
	banner.text = "👥 %s 進入分身模式！%d 秒" % [player_name, duration_sec]
	banner.add_theme_font_size_override("font_size", 13)
	banner.add_theme_color_override("font_color", COLOR_PALE)
	banner.position = Vector2(vp_size.x / 2 - 110, 6)
	add_child(banner)

	var tw = banner.create_tween()
	tw.tween_interval(3.0)
	tw.tween_property(banner, "modulate:a", 0.0, 0.4)
	tw.tween_callback(banner.queue_free)

## clone_hit — 分身子彈命中
func _on_clone_hit(payload: Dictionary) -> void:
	var side: String = payload.get("side", "left")
	var killed: bool = payload.get("killed", false)
	var reward: int = payload.get("reward", 0)
	var hit_x: float = payload.get("x", 0.0)
	var hit_y: float = payload.get("y", 0.0)

	# 依左右選顏色
	var hit_color = COLOR_LEFT if side == "left" else COLOR_RIGHT
	var side_icon = "◀" if side == "left" else "▶"

	# 命中閃光（小範圍）
	var flash = ColorRect.new()
	flash.color = Color(hit_color.r, hit_color.g, hit_color.b, 0.35)
	flash.size = Vector2(60, 60)
	flash.position = Vector2(hit_x - 30, hit_y - 30)
	add_child(flash)

	var tw_flash = flash.create_tween()
	tw_flash.tween_property(flash, "modulate:a", 0.0, 0.15)
	tw_flash.tween_callback(flash.queue_free)

	if killed and reward > 0:
		# 倍率浮動文字
		var mult_label = Label.new()
		mult_label.text = "%s ×0.7 +%d" % [side_icon, reward]
		mult_label.add_theme_font_size_override("font_size", 16)
		mult_label.add_theme_color_override("font_color", hit_color)
		mult_label.position = Vector2(hit_x - 30, hit_y - 30)
		add_child(mult_label)

		var tw_mult = mult_label.create_tween()
		tw_mult.tween_property(mult_label, "position:y", hit_y - 65, 0.5)
		tw_mult.parallel().tween_property(mult_label, "modulate:a", 0.0, 0.5)
		tw_mult.tween_callback(mult_label.queue_free)

## clone_end — 分身模式結束
func _on_clone_end(_payload: Dictionary) -> void:
	# 清除計時條
	if is_instance_valid(_timer_tween):
		_timer_tween.kill()
	if is_instance_valid(_timer_bar_bg):
		var tw = _timer_bar_bg.create_tween()
		tw.tween_property(_timer_bar_bg, "modulate:a", 0.0, 0.3)
		tw.tween_callback(_timer_bar_bg.queue_free)
	if is_instance_valid(_timer_bar):
		var tw2 = _timer_bar.create_tween()
		tw2.tween_property(_timer_bar, "modulate:a", 0.0, 0.3)
		tw2.tween_callback(_timer_bar.queue_free)

	# 結束提示
	var vp_size = get_viewport().size
	var end_label = Label.new()
	end_label.text = "👥 分身模式結束"
	end_label.add_theme_font_size_override("font_size", 14)
	end_label.add_theme_color_override("font_color", Color(0.6, 0.4, 0.7))
	end_label.position = Vector2(vp_size.x / 2 - 70, vp_size.y - 60)
	add_child(end_label)

	var tw = end_label.create_tween()
	tw.tween_interval(1.5)
	tw.tween_property(end_label, "modulate:a", 0.0, 0.4)
	tw.tween_callback(end_label.queue_free)

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.40)
	flash.size = vp_size
	flash.position = Vector2.ZERO
	add_child(flash)

	var tw = flash.create_tween()
	tw.tween_property(flash, "modulate:a", 0.0, duration)
	tw.tween_callback(flash.queue_free)
