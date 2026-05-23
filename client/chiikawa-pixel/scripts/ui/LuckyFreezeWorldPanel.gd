## LuckyFreezeWorldPanel.gd — 幸運冰凍世界魚系統面板（DAY-237）
## 業界原創「全場冰凍+冰裂爆發」機制
##
## 視覺設計：
##   - 冰藍主題（#5DADE2 + #2E86C1 + #AED6F1 + #EBF5FB）
##   - freeze_start：冰藍三次強閃光 + 頂部橫幅 + 「❄️ 冰凍世界！」大字 + 計時條 + 速度提示
##   - freeze_crack：全螢幕白色強閃光 + 「❄️ 冰裂爆發！」大字 + HP -50% 提示
##   - freeze_end：冰藍淡出 + 速度恢復提示
extends CanvasLayer

# 冰凍世界狀態
var _active: bool = false
var _duration_sec: int = 8

# 計時條節點
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null

# 主題顏色
const COLOR_PRIMARY  = Color("#5DADE2")  # 冰藍
const COLOR_DARK     = Color("#2E86C1")  # 深藍
const COLOR_PALE     = Color("#AED6F1")  # 淡藍
const COLOR_LIGHT_BG = Color("#EBF5FB")  # 極淡藍
const COLOR_GOLD     = Color("#FFD700")  # 金黃
const COLOR_WHITE    = Color("#FFFFFF")  # 白色

func _ready() -> void:
	layer = 8  # 幸運冰凍世界魚面板層級

## 處理幸運冰凍世界魚訊息
func handle_lucky_freeze_world(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"freeze_start":
			_on_freeze_start(payload)
		"freeze_crack":
			_on_freeze_crack(payload)
		"freeze_end":
			_on_freeze_end(payload)

## freeze_start — 冰凍世界開始
func _on_freeze_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	_duration_sec = payload.get("duration_sec", 8)
	var kill_boost: float = payload.get("kill_boost", 2.0)
	var frozen_count: int = payload.get("frozen_count", 0)
	_active = true

	var vp_size = get_viewport().size

	# 冰藍三次強閃光
	_flash_screen(COLOR_PRIMARY, 0.12)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(COLOR_WHITE, 0.10)
	await get_tree().create_timer(0.07).timeout
	_flash_screen(COLOR_PALE, 0.09)

	# 頂部橫幅
	var banner = Label.new()
	banner.text = "❄️ %s 觸發冰凍世界！%d 個目標速度降低 80%%！" % [player_name, frozen_count]
	banner.add_theme_font_size_override("font_size", 14)
	banner.add_theme_color_override("font_color", COLOR_PALE)
	banner.position = Vector2(vp_size.x / 2 - 160, 6)
	add_child(banner)

	var tween_banner = banner.create_tween()
	tween_banner.tween_interval(float(_duration_sec) - 0.5)
	tween_banner.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween_banner.tween_callback(banner.queue_free)

	# 「❄️ 冰凍世界！」大字
	var big_label = Label.new()
	big_label.text = "❄️ 冰凍世界！"
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
	mult_label.text = "冰凍期間 ×%.1f 倍率加成！" % kill_boost
	mult_label.add_theme_font_size_override("font_size", 13)
	mult_label.add_theme_color_override("font_color", COLOR_GOLD)
	mult_label.position = Vector2(vp_size.x / 2 - 75, vp_size.y / 2 + 28)
	add_child(mult_label)

	var tween_mult = mult_label.create_tween()
	tween_mult.tween_interval(2.0)
	tween_mult.tween_property(mult_label, "modulate:a", 0.0, 0.5)
	tween_mult.tween_callback(mult_label.queue_free)

	# 冰晶粒子效果（四角冰晶）
	_spawn_ice_crystals(vp_size)

	# 底部計時條（冰藍→深藍漸變）
	_spawn_timer_bar(float(_duration_sec))

## freeze_crack — 冰裂爆發
func _on_freeze_crack(payload: Dictionary) -> void:
	var cracked_count: int = payload.get("cracked_count", 0)

	var vp_size = get_viewport().size

	# 清除計時條
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		_timer_bar_bg.queue_free()
		_timer_bar_bg = null

	# 全螢幕白色強閃光（冰裂感）
	_flash_screen(COLOR_WHITE, 0.18)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_PRIMARY, 0.12)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(COLOR_PALE, 0.10)

	# 「❄️ 冰裂爆發！」大字
	var crack_label = Label.new()
	crack_label.text = "❄️ 冰裂爆發！"
	crack_label.add_theme_font_size_override("font_size", 48)
	crack_label.add_theme_color_override("font_color", COLOR_DARK)
	crack_label.position = vp_size / 2 - Vector2(100, 28)
	add_child(crack_label)

	var tween_crack = crack_label.create_tween()
	tween_crack.tween_property(crack_label, "scale", Vector2(1.3, 1.3), 0.12)
	tween_crack.tween_property(crack_label, "scale", Vector2(1.0, 1.0), 0.10)
	tween_crack.tween_interval(0.6)
	tween_crack.tween_property(crack_label, "modulate:a", 0.0, 0.5)
	tween_crack.tween_callback(crack_label.queue_free)

	# HP -50% 提示
	if cracked_count > 0:
		var hp_label = Label.new()
		hp_label.text = "💥 %d 個目標 HP -50%%！" % cracked_count
		hp_label.add_theme_font_size_override("font_size", 16)
		hp_label.add_theme_color_override("font_color", COLOR_GOLD)
		hp_label.position = Vector2(vp_size.x / 2 - 80, vp_size.y / 2 + 30)
		add_child(hp_label)

		var tween_hp = hp_label.create_tween()
		tween_hp.tween_property(hp_label, "position:y", hp_label.position.y - 20, 0.5)
		tween_hp.parallel().tween_property(hp_label, "modulate:a", 0.0, 0.5)
		tween_hp.tween_callback(hp_label.queue_free)

## freeze_end — 冰凍世界結束
func _on_freeze_end(_payload: Dictionary) -> void:
	_active = false

	# 清除計時條（若崩潰前已清除則跳過）
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		_timer_bar_bg.queue_free()
		_timer_bar_bg = null

	var vp_size = get_viewport().size

	# 速度恢復提示
	var end_label = Label.new()
	end_label.text = "❄️ 冰凍結束，速度恢復"
	end_label.add_theme_font_size_override("font_size", 14)
	end_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	end_label.position = Vector2(vp_size.x / 2 - 65, vp_size.y - 40)
	add_child(end_label)

	var tween_end = end_label.create_tween()
	tween_end.tween_interval(0.5)
	tween_end.tween_property(end_label, "modulate:a", 0.0, 0.5)
	tween_end.tween_callback(end_label.queue_free)

## 建立冰晶粒子效果（四角）
func _spawn_ice_crystals(vp_size: Vector2) -> void:
	var corners = [
		Vector2(20, 20),
		Vector2(vp_size.x - 40, 20),
		Vector2(20, vp_size.y - 40),
		Vector2(vp_size.x - 40, vp_size.y - 40),
	]
	for corner in corners:
		var crystal = Label.new()
		crystal.text = "❄"
		crystal.add_theme_font_size_override("font_size", 24)
		crystal.add_theme_color_override("font_color", COLOR_PALE)
		crystal.position = corner
		add_child(crystal)

		var tween_c = crystal.create_tween()
		tween_c.tween_property(crystal, "modulate:a", 0.8, 0.3)
		tween_c.tween_property(crystal, "modulate:a", 0.2, 0.3)
		tween_c.tween_property(crystal, "modulate:a", 0.8, 0.3)
		tween_c.tween_property(crystal, "modulate:a", 0.2, 0.3)
		tween_c.tween_interval(float(_duration_sec) - 1.5)
		tween_c.tween_property(crystal, "modulate:a", 0.0, 0.5)
		tween_c.tween_callback(crystal.queue_free)

## 建立底部計時條（冰藍→深藍漸變）
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
