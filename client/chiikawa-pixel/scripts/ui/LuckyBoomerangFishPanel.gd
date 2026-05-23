## LuckyBoomerangFishPanel.gd — 幸運迴旋鏢魚系統面板（DAY-231）
## 業界原創「迴旋鏢來回穿透」機制
##
## 視覺設計：
##   - 橙棕迴旋鏢主題（#E67E22 + #D35400 + #FAD7A0 + #FFF3E0）
##   - boomerang_start：橙色雙閃光 + 頂部橫幅 + 「🪃 迴旋鏢模式！」大字 + 計時條
##   - boomerang_hit：命中閃光 + 折返次數標記 + 浮動獎勵文字
##   - boomerang_end：橙色淡出
extends CanvasLayer

# 迴旋鏢狀態
var _active: bool = false
var _player_id: String = ""
var _duration_sec: int = 10

# 計時條節點
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null

# 主題顏色
const COLOR_PRIMARY  = Color("#E67E22")  # 橙色
const COLOR_DARK     = Color("#D35400")  # 深橙
const COLOR_PALE     = Color("#FAD7A0")  # 淡橙
const COLOR_LIGHT_BG = Color("#FFF3E0")  # 極淡橙
const COLOR_GOLD     = Color("#FFD700")  # 金黃

func _ready() -> void:
	layer = 14  # 幸運迴旋鏢魚面板層級

## 處理幸運迴旋鏢魚訊息
func handle_lucky_boomerang_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"boomerang_start":
			_on_boomerang_start(payload)
		"boomerang_hit":
			_on_boomerang_hit(payload)
		"boomerang_end":
			_on_boomerang_end(payload)

## boomerang_start — 迴旋鏢模式開始
func _on_boomerang_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	_player_id = payload.get("player_id", "")
	_duration_sec = payload.get("duration_sec", 10)
	var max_bounces: int = payload.get("max_bounces", 3)
	var kill_mult: float = payload.get("kill_mult", 0.65)
	_active = true

	var vp_size = get_viewport().size

	# 橙色雙閃光
	_flash_screen(COLOR_PRIMARY, 0.16)
	await get_tree().create_timer(0.10).timeout
	_flash_screen(COLOR_DARK, 0.13)

	# 頂部橫幅
	var banner = Label.new()
	banner.text = "🪃 %s 觸發迴旋鏢模式！最多折返 %d 次！" % [player_name, max_bounces]
	banner.add_theme_font_size_override("font_size", 14)
	banner.add_theme_color_override("font_color", COLOR_PALE)
	banner.position = Vector2(vp_size.x / 2 - 160, 6)
	add_child(banner)

	var tween_banner = banner.create_tween()
	tween_banner.tween_interval(float(_duration_sec) - 0.5)
	tween_banner.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween_banner.tween_callback(banner.queue_free)

	# 「🪃 迴旋鏢模式！」大字
	var big_label = Label.new()
	big_label.text = "🪃 迴旋鏢模式！"
	big_label.add_theme_font_size_override("font_size", 48)
	big_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	big_label.position = vp_size / 2 - Vector2(110, 28)
	add_child(big_label)

	var tween_big = big_label.create_tween()
	tween_big.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.10)
	tween_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_big.tween_interval(0.5)
	tween_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_big.tween_callback(big_label.queue_free)

	# 倍率提示
	var mult_label = Label.new()
	mult_label.text = "每次命中 ×%.2f 倍率！" % kill_mult
	mult_label.add_theme_font_size_override("font_size", 13)
	mult_label.add_theme_color_override("font_color", COLOR_GOLD)
	mult_label.position = Vector2(vp_size.x / 2 - 70, vp_size.y / 2 + 28)
	add_child(mult_label)

	var tween_mult = mult_label.create_tween()
	tween_mult.tween_interval(2.0)
	tween_mult.tween_property(mult_label, "modulate:a", 0.0, 0.5)
	tween_mult.tween_callback(mult_label.queue_free)

	# 底部計時條
	_spawn_timer_bar(float(_duration_sec))

## boomerang_hit — 迴旋鏢命中目標
func _on_boomerang_hit(payload: Dictionary) -> void:
	var bounce_num: int = payload.get("bounce_num", 1)
	var killed: bool = payload.get("killed", false)
	var reward: int = payload.get("reward", 0)
	var x: float = payload.get("x", 500.0)
	var y: float = payload.get("y", 300.0)

	var vp_size = get_viewport().size

	# 命中閃光（橙色）
	_flash_screen(COLOR_PRIMARY, 0.08)

	# 折返次數標記（在命中位置附近）
	var bounce_label = Label.new()
	bounce_label.text = "🪃 第%d折返" % bounce_num
	bounce_label.add_theme_font_size_override("font_size", 13)
	bounce_label.add_theme_color_override("font_color", COLOR_PALE)
	# 轉換到螢幕座標（假設遊戲座標和螢幕座標比例）
	var screen_x = x * vp_size.x / 1000.0
	var screen_y = y * vp_size.y / 600.0
	bounce_label.position = Vector2(screen_x - 30, screen_y - 20)
	add_child(bounce_label)

	var tween_b = bounce_label.create_tween()
	tween_b.tween_property(bounce_label, "position:y", bounce_label.position.y - 18, 0.4)
	tween_b.parallel().tween_property(bounce_label, "modulate:a", 0.0, 0.4)
	tween_b.tween_callback(bounce_label.queue_free)

	# 擊破獎勵浮動文字（金色）
	if killed and reward > 0:
		var reward_label = Label.new()
		reward_label.text = "+%d" % reward
		reward_label.add_theme_font_size_override("font_size", 16)
		reward_label.add_theme_color_override("font_color", COLOR_GOLD)
		reward_label.position = Vector2(screen_x - 15, screen_y - 35)
		add_child(reward_label)

		var tween_r = reward_label.create_tween()
		tween_r.tween_property(reward_label, "position:y", reward_label.position.y - 25, 0.5)
		tween_r.parallel().tween_property(reward_label, "modulate:a", 0.0, 0.5)
		tween_r.tween_callback(reward_label.queue_free)

## boomerang_end — 迴旋鏢模式結束
func _on_boomerang_end(_payload: Dictionary) -> void:
	_active = false

	# 清除計時條
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		_timer_bar_bg.queue_free()
		_timer_bar_bg = null

	var vp_size = get_viewport().size

	# 橙色淡出提示
	var end_label = Label.new()
	end_label.text = "🪃 迴旋鏢模式結束"
	end_label.add_theme_font_size_override("font_size", 14)
	end_label.add_theme_color_override("font_color", COLOR_DARK)
	end_label.position = Vector2(vp_size.x / 2 - 70, vp_size.y - 40)
	add_child(end_label)

	var tween_end = end_label.create_tween()
	tween_end.tween_interval(0.5)
	tween_end.tween_property(end_label, "modulate:a", 0.0, 0.5)
	tween_end.tween_callback(end_label.queue_free)

## 建立底部計時條
func _spawn_timer_bar(duration: float) -> void:
	var vp_size = get_viewport().size

	# 背景條
	var bg = ColorRect.new()
	bg.color = Color(0.1, 0.1, 0.1, 0.6)
	bg.size = Vector2(vp_size.x, 8)
	bg.position = Vector2(0, vp_size.y - 8)
	add_child(bg)
	_timer_bar_bg = bg

	# 計時條（橙→深橙漸變）
	var bar = ColorRect.new()
	bar.color = COLOR_PRIMARY
	bar.size = Vector2(vp_size.x, 8)
	bar.position = Vector2(0, vp_size.y - 8)
	add_child(bar)
	_timer_bar = bar

	# 計時條縮短動畫
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
	flash.color = Color(color.r, color.g, color.b, 0.28)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)
