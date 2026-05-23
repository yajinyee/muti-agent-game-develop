## LuckyMirrorWorldPanel.gd — 幸運鏡面世界魚系統面板（DAY-236）
## 業界原創「全場鏡像反轉+鏡面崩潰」機制
##
## 視覺設計：
##   - 紫色鏡面主題（#8E44AD + #6C3483 + #D7BDE2 + #F5EEF8）
##   - mirror_start：紫色三次強閃光 + 頂部橫幅 + 「🪞 鏡面世界！」大字 + 計時條 + 倍率提示
##   - mirror_collapse：全螢幕白色強閃光 + 「🪞 鏡面崩潰！」大字 + HP -35% 提示
##   - mirror_end：紫色淡出
extends CanvasLayer

# 鏡面世界狀態
var _active: bool = false
var _duration_sec: int = 10

# 計時條節點
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null

# 主題顏色
const COLOR_PRIMARY  = Color("#8E44AD")  # 紫色
const COLOR_DARK     = Color("#6C3483")  # 深紫
const COLOR_PALE     = Color("#D7BDE2")  # 淡紫
const COLOR_LIGHT_BG = Color("#F5EEF8")  # 極淡紫
const COLOR_GOLD     = Color("#FFD700")  # 金黃
const COLOR_ORANGE   = Color("#E67E22")  # 橙色

func _ready() -> void:
	layer = 9  # 幸運鏡面世界魚面板層級

## 處理幸運鏡面世界魚訊息
func handle_lucky_mirror_world(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"mirror_start":
			_on_mirror_start(payload)
		"mirror_collapse":
			_on_mirror_collapse(payload)
		"mirror_end":
			_on_mirror_end(payload)

## mirror_start — 鏡面世界開始
func _on_mirror_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	_duration_sec = payload.get("duration_sec", 10)
	var kill_boost: float = payload.get("kill_boost", 2.3)
	_active = true

	var vp_size = get_viewport().size

	# 紫色三次強閃光
	_flash_screen(COLOR_PRIMARY, 0.12)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(Color.WHITE, 0.10)
	await get_tree().create_timer(0.07).timeout
	_flash_screen(COLOR_PALE, 0.09)

	# 頂部橫幅
	var banner = Label.new()
	banner.text = "🪞 %s 觸發鏡面世界！目標左右翻轉！" % player_name
	banner.add_theme_font_size_override("font_size", 14)
	banner.add_theme_color_override("font_color", COLOR_PALE)
	banner.position = Vector2(vp_size.x / 2 - 140, 6)
	add_child(banner)

	var tween_banner = banner.create_tween()
	tween_banner.tween_interval(float(_duration_sec) - 0.5)
	tween_banner.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween_banner.tween_callback(banner.queue_free)

	# 「🪞 鏡面世界！」大字
	var big_label = Label.new()
	big_label.text = "🪞 鏡面世界！"
	big_label.add_theme_font_size_override("font_size", 48)
	big_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	big_label.position = vp_size / 2 - Vector2(100, 28)
	add_child(big_label)

	var tween_big = big_label.create_tween()
	tween_big.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.10)
	tween_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_big.tween_interval(0.5)
	tween_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_big.tween_callback(big_label.queue_free)

	# 倍率提示
	var mult_label = Label.new()
	mult_label.text = "鏡面期間 ×%.1f 倍率加成！" % kill_boost
	mult_label.add_theme_font_size_override("font_size", 13)
	mult_label.add_theme_color_override("font_color", COLOR_GOLD)
	mult_label.position = Vector2(vp_size.x / 2 - 75, vp_size.y / 2 + 28)
	add_child(mult_label)

	var tween_mult = mult_label.create_tween()
	tween_mult.tween_interval(2.0)
	tween_mult.tween_property(mult_label, "modulate:a", 0.0, 0.5)
	tween_mult.tween_callback(mult_label.queue_free)

	# 鏡面線（場景中央垂直線）
	_spawn_mirror_line(vp_size)

	# 底部計時條（紫→深紫漸變）
	_spawn_timer_bar(float(_duration_sec))

## mirror_collapse — 鏡面崩潰
func _on_mirror_collapse(payload: Dictionary) -> void:
	var collapsed_count: int = payload.get("collapsed_count", 0)

	var vp_size = get_viewport().size

	# 清除計時條
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		_timer_bar_bg.queue_free()
		_timer_bar_bg = null

	# 全螢幕白色強閃光（鏡面破碎感）
	_flash_screen(Color.WHITE, 0.15)
	await get_tree().create_timer(0.10).timeout
	_flash_screen(COLOR_PRIMARY, 0.12)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(COLOR_PALE, 0.10)

	# 「🪞 鏡面崩潰！」大字
	var collapse_label = Label.new()
	collapse_label.text = "🪞 鏡面崩潰！"
	collapse_label.add_theme_font_size_override("font_size", 48)
	collapse_label.add_theme_color_override("font_color", COLOR_DARK)
	collapse_label.position = vp_size / 2 - Vector2(100, 28)
	add_child(collapse_label)

	var tween_collapse = collapse_label.create_tween()
	tween_collapse.tween_property(collapse_label, "scale", Vector2(1.3, 1.3), 0.12)
	tween_collapse.tween_property(collapse_label, "scale", Vector2(1.0, 1.0), 0.10)
	tween_collapse.tween_interval(0.6)
	tween_collapse.tween_property(collapse_label, "modulate:a", 0.0, 0.5)
	tween_collapse.tween_callback(collapse_label.queue_free)

	# HP -35% 提示
	if collapsed_count > 0:
		var hp_label = Label.new()
		hp_label.text = "💥 %d 個目標 HP -35%%！" % collapsed_count
		hp_label.add_theme_font_size_override("font_size", 16)
		hp_label.add_theme_color_override("font_color", COLOR_ORANGE)
		hp_label.position = Vector2(vp_size.x / 2 - 90, vp_size.y / 2 + 30)
		add_child(hp_label)

		var tween_hp = hp_label.create_tween()
		tween_hp.tween_property(hp_label, "position:y", hp_label.position.y - 20, 0.5)
		tween_hp.parallel().tween_property(hp_label, "modulate:a", 0.0, 0.5)
		tween_hp.tween_callback(hp_label.queue_free)

## mirror_end — 鏡面世界結束
func _on_mirror_end(_payload: Dictionary) -> void:
	_active = false

	# 清除計時條（若崩潰前已清除則跳過）
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		_timer_bar_bg.queue_free()
		_timer_bar_bg = null

	var vp_size = get_viewport().size

	# 紫色淡出提示
	var end_label = Label.new()
	end_label.text = "🪞 鏡面世界結束"
	end_label.add_theme_font_size_override("font_size", 14)
	end_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	end_label.position = Vector2(vp_size.x / 2 - 55, vp_size.y - 40)
	add_child(end_label)

	var tween_end = end_label.create_tween()
	tween_end.tween_interval(0.5)
	tween_end.tween_property(end_label, "modulate:a", 0.0, 0.5)
	tween_end.tween_callback(end_label.queue_free)

## 建立鏡面線（場景中央垂直線，視覺提示）
func _spawn_mirror_line(vp_size: Vector2) -> void:
	var center_x = vp_size.x / 2.0

	var line = ColorRect.new()
	line.color = Color(COLOR_PRIMARY.r, COLOR_PRIMARY.g, COLOR_PRIMARY.b, 0.25)
	line.size = Vector2(3, vp_size.y)
	line.position = Vector2(center_x - 1.5, 0)
	add_child(line)

	# 閃爍動畫後淡出
	var tween_line = line.create_tween()
	tween_line.tween_property(line, "modulate:a", 0.6, 0.3)
	tween_line.tween_property(line, "modulate:a", 0.15, 0.3)
	tween_line.tween_property(line, "modulate:a", 0.6, 0.3)
	tween_line.tween_property(line, "modulate:a", 0.15, 0.3)
	tween_line.tween_interval(float(_duration_sec) - 1.5)
	tween_line.tween_property(line, "modulate:a", 0.0, 0.5)
	tween_line.tween_callback(line.queue_free)

## 建立底部計時條（紫→深紫漸變）
func _spawn_timer_bar(duration: float) -> void:
	var vp_size = get_viewport().size

	var bg = ColorRect.new()
	bg.color = Color(0.1, 0.1, 0.1, 0.6)
	bg.size = Vector2(vp_size.x, 8)
	bg.position = Vector2(0, vp_size.y - 8)
	add_child(bg)
	_timer_bar_bg = bg

	var bar = ColorRect.new()
	bar.color = COLOR_PRIMARY
	bar.size = Vector2(vp_size.x, 8)
	bar.position = Vector2(0, vp_size.y - 8)
	add_child(bar)
	_timer_bar = bar

	var tween_bar = bar.create_tween()
	tween_bar.tween_property(bar, "size:x", 0.0, duration)
	tween_bar.tween_callback(func():
		if is_instance_valid(bar):
			bar.queue_free()
		if is_instance_valid(bg):
			bg.queue_free()
	)

# ---- 輔助函數 ----

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.26)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)
