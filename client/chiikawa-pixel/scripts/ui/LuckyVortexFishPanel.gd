## LuckyVortexFishPanel.gd — 幸運漩渦魚系統面板（DAY-234）
## 業界原創「漩渦旋轉+高倍率加成」機制
##
## 視覺設計：
##   - 青綠漩渦主題（#16A085 + #1ABC9C + #A3E4D7 + #E8F8F5）
##   - vortex_start：青綠三次強閃光 + 頂部橫幅 + 「🌀 漩渦啟動！」大字 + 漩渦圓圈 + 計時條
##   - vortex_rotate：青綠閃光 + 「🌀 旋轉！」提示 + 旋轉目標數
##   - vortex_blast：全螢幕三次強閃光 + 「🌀 漩渦爆發！」大字 + 結算彈窗
##   - vortex_end：青綠淡出
extends CanvasLayer

# 漩渦狀態
var _active: bool = false
var _duration_sec: int = 8
var _rotate_count: int = 0

# 計時條節點
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null

# 漩渦圓圈節點
var _vortex_ring: ColorRect = null

# 主題顏色
const COLOR_PRIMARY  = Color("#16A085")  # 青綠
const COLOR_BRIGHT   = Color("#1ABC9C")  # 亮青綠
const COLOR_PALE     = Color("#A3E4D7")  # 淡青綠
const COLOR_LIGHT_BG = Color("#E8F8F5")  # 極淡青綠
const COLOR_GOLD     = Color("#FFD700")  # 金黃

func _ready() -> void:
	layer = 11  # 幸運漩渦魚面板層級

## 處理幸運漩渦魚訊息
func handle_lucky_vortex_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"vortex_start":
			_on_vortex_start(payload)
		"vortex_rotate":
			_on_vortex_rotate(payload)
		"vortex_blast":
			_on_vortex_blast(payload)
		"vortex_end":
			_on_vortex_end(payload)

## vortex_start — 漩渦模式開始
func _on_vortex_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	_duration_sec = payload.get("duration_sec", 8)
	var kill_boost: float = payload.get("kill_boost", 2.2)
	var radius: float = payload.get("radius", 300.0)
	_active = true
	_rotate_count = 0

	var vp_size = get_viewport().size

	# 青綠三次強閃光
	_flash_screen(COLOR_PRIMARY, 0.12)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(COLOR_BRIGHT, 0.10)
	await get_tree().create_timer(0.07).timeout
	_flash_screen(COLOR_PALE, 0.09)

	# 頂部橫幅
	var banner = Label.new()
	banner.text = "🌀 %s 觸發漩渦！目標開始旋轉！" % player_name
	banner.add_theme_font_size_override("font_size", 14)
	banner.add_theme_color_override("font_color", COLOR_PALE)
	banner.position = Vector2(vp_size.x / 2 - 140, 6)
	add_child(banner)

	var tween_banner = banner.create_tween()
	tween_banner.tween_interval(float(_duration_sec) - 0.5)
	tween_banner.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween_banner.tween_callback(banner.queue_free)

	# 「🌀 漩渦啟動！」大字
	var big_label = Label.new()
	big_label.text = "🌀 漩渦啟動！"
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
	mult_label.text = "漩渦期間 ×%.1f 倍率加成！" % kill_boost
	mult_label.add_theme_font_size_override("font_size", 13)
	mult_label.add_theme_color_override("font_color", COLOR_GOLD)
	mult_label.position = Vector2(vp_size.x / 2 - 75, vp_size.y / 2 + 28)
	add_child(mult_label)

	var tween_mult = mult_label.create_tween()
	tween_mult.tween_interval(2.0)
	tween_mult.tween_property(mult_label, "modulate:a", 0.0, 0.5)
	tween_mult.tween_callback(mult_label.queue_free)

	# 漩渦圓圈（視覺提示）
	_spawn_vortex_ring(vp_size, radius)

	# 底部計時條（青綠漸變）
	_spawn_timer_bar(float(_duration_sec))

## vortex_rotate — 漩渦旋轉（每 2 秒）
func _on_vortex_rotate(payload: Dictionary) -> void:
	var rotate_num: int = payload.get("rotate_num", 1)
	var rotated_count: int = payload.get("rotated_count", 0)
	_rotate_count = rotate_num

	if rotated_count == 0:
		return

	var vp_size = get_viewport().size

	# 青綠閃光（輕微）
	_flash_screen(COLOR_PRIMARY, 0.06)

	# 「🌀 旋轉！」提示
	var rotate_label = Label.new()
	rotate_label.text = "🌀 旋轉！%d 個目標" % rotated_count
	rotate_label.add_theme_font_size_override("font_size", 13)
	rotate_label.add_theme_color_override("font_color", COLOR_PALE)
	rotate_label.position = Vector2(vp_size.x / 2 - 55, vp_size.y / 2 - 60)
	add_child(rotate_label)

	var tween_r = rotate_label.create_tween()
	tween_r.tween_property(rotate_label, "position:y", rotate_label.position.y - 15, 0.4)
	tween_r.parallel().tween_property(rotate_label, "modulate:a", 0.0, 0.4)
	tween_r.tween_callback(rotate_label.queue_free)

	# 漩渦圓圈旋轉動畫（視覺反饋）
	if is_instance_valid(_vortex_ring):
		var tween_ring = _vortex_ring.create_tween()
		tween_ring.tween_property(_vortex_ring, "modulate:a", 0.6, 0.15)
		tween_ring.tween_property(_vortex_ring, "modulate:a", 0.2, 0.15)

## vortex_blast — 漩渦爆發
func _on_vortex_blast(payload: Dictionary) -> void:
	var killed_count: int = payload.get("killed_count", 0)
	var total_reward: int = payload.get("total_reward", 0)

	var vp_size = get_viewport().size

	# 清除計時條
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		_timer_bar_bg.queue_free()
		_timer_bar_bg = null
	# 清除漩渦圓圈
	if is_instance_valid(_vortex_ring):
		_vortex_ring.queue_free()
		_vortex_ring = null

	# 全螢幕三次強閃光（青綠→白→亮青綠）
	_flash_screen(COLOR_PRIMARY, 0.12)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(Color.WHITE, 0.10)
	await get_tree().create_timer(0.07).timeout
	_flash_screen(COLOR_BRIGHT, 0.12)

	# 「🌀 漩渦爆發！」大字
	var blast_label = Label.new()
	blast_label.text = "🌀 漩渦爆發！"
	blast_label.add_theme_font_size_override("font_size", 52)
	blast_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	blast_label.position = vp_size / 2 - Vector2(110, 32)
	add_child(blast_label)

	var tween_blast = blast_label.create_tween()
	tween_blast.tween_property(blast_label, "scale", Vector2(1.3, 1.3), 0.12)
	tween_blast.tween_property(blast_label, "scale", Vector2(1.0, 1.0), 0.10)
	tween_blast.tween_interval(0.8)
	tween_blast.tween_property(blast_label, "modulate:a", 0.0, 0.5)
	tween_blast.tween_callback(blast_label.queue_free)

	# 結算彈窗（右側滑入）
	if killed_count > 0:
		_spawn_blast_result_panel(vp_size, killed_count, total_reward)

## vortex_end — 漩渦模式結束
func _on_vortex_end(_payload: Dictionary) -> void:
	_active = false
	_rotate_count = 0

	# 清除計時條
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		_timer_bar_bg.queue_free()
		_timer_bar_bg = null
	# 清除漩渦圓圈
	if is_instance_valid(_vortex_ring):
		_vortex_ring.queue_free()
		_vortex_ring = null

	var vp_size = get_viewport().size

	# 青綠淡出提示
	var end_label = Label.new()
	end_label.text = "🌀 漩渦結束"
	end_label.add_theme_font_size_override("font_size", 14)
	end_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	end_label.position = Vector2(vp_size.x / 2 - 45, vp_size.y - 40)
	add_child(end_label)

	var tween_end = end_label.create_tween()
	tween_end.tween_interval(0.5)
	tween_end.tween_property(end_label, "modulate:a", 0.0, 0.5)
	tween_end.tween_callback(end_label.queue_free)

## 建立漩渦圓圈（視覺提示）
func _spawn_vortex_ring(vp_size: Vector2, radius: float) -> void:
	var scale_x = vp_size.x / 1000.0
	var scale_y = vp_size.y / 600.0
	var center_x = 500.0 * scale_x
	var center_y = 300.0 * scale_y
	var radius_px = radius * (scale_x + scale_y) / 2.0

	# 外圈（青綠半透明）
	var ring = ColorRect.new()
	ring.color = Color(COLOR_PRIMARY.r, COLOR_PRIMARY.g, COLOR_PRIMARY.b, 0.12)
	ring.size = Vector2(radius_px * 2, radius_px * 2)
	ring.position = Vector2(center_x - radius_px, center_y - radius_px)
	add_child(ring)
	_vortex_ring = ring

	# 旋轉脈衝動畫
	var tween_ring = ring.create_tween().set_loops(4)
	tween_ring.tween_property(ring, "modulate:a", 0.4, 0.5)
	tween_ring.tween_property(ring, "modulate:a", 0.1, 0.5)

## 建立漩渦爆發結算彈窗
func _spawn_blast_result_panel(vp_size: Vector2, killed_count: int, total_reward: int) -> void:
	var panel = ColorRect.new()
	panel.color = Color(0.05, 0.25, 0.20, 0.88)
	panel.size = Vector2(200, 90)
	panel.position = Vector2(vp_size.x + 10, vp_size.y / 2 - 45)
	add_child(panel)

	var title = Label.new()
	title.text = "🌀 漩渦爆發結算"
	title.add_theme_font_size_override("font_size", 13)
	title.add_theme_color_override("font_color", COLOR_PALE)
	title.position = Vector2(10, 8)
	panel.add_child(title)

	var kill_label = Label.new()
	kill_label.text = "擊破：%d 個目標" % killed_count
	kill_label.add_theme_font_size_override("font_size", 12)
	kill_label.add_theme_color_override("font_color", Color.WHITE)
	kill_label.position = Vector2(10, 30)
	panel.add_child(kill_label)

	var reward_label = Label.new()
	reward_label.text = "全服共享：+%d" % total_reward
	reward_label.add_theme_font_size_override("font_size", 12)
	reward_label.add_theme_color_override("font_color", COLOR_GOLD)
	reward_label.position = Vector2(10, 52)
	panel.add_child(reward_label)

	var tween_panel = panel.create_tween()
	tween_panel.tween_property(panel, "position:x", vp_size.x - 210, 0.3)
	tween_panel.tween_interval(2.5)
	tween_panel.tween_property(panel, "modulate:a", 0.0, 0.4)
	tween_panel.tween_callback(panel.queue_free)

## 建立底部計時條
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
