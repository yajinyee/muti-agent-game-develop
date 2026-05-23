## LuckyMagnetFishPanel.gd — 幸運磁力魚系統面板（DAY-232）
## 業界原創「磁力聚集+磁力爆發」機制
##
## 視覺設計：
##   - 藍色磁力主題（#3498DB + #2980B9 + #AED6F1 + #EBF5FB）
##   - magnet_start：藍色雙閃光 + 頂部橫幅 + 「🧲 磁力場啟動！」大字 + 計時條
##   - magnet_pull：藍色閃光 + 「🧲 磁力吸引！」提示 + 移動目標數
##   - magnet_blast：全螢幕三次強閃光 + 「🧲 磁力爆發！」大字 + 結算彈窗
##   - magnet_end：藍色淡出
extends CanvasLayer

# 磁力場狀態
var _active: bool = false
var _player_id: String = ""
var _duration_sec: int = 12
var _pull_count: int = 0

# 計時條節點
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null

# 主題顏色
const COLOR_PRIMARY  = Color("#3498DB")  # 藍色
const COLOR_DARK     = Color("#2980B9")  # 深藍
const COLOR_PALE     = Color("#AED6F1")  # 淡藍
const COLOR_LIGHT_BG = Color("#EBF5FB")  # 極淡藍
const COLOR_GOLD     = Color("#FFD700")  # 金黃
const COLOR_CYAN     = Color("#1ABC9C")  # 青綠（爆發用）

func _ready() -> void:
	layer = 13  # 幸運磁力魚面板層級

## 處理幸運磁力魚訊息
func handle_lucky_magnet_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"magnet_start":
			_on_magnet_start(payload)
		"magnet_pull":
			_on_magnet_pull(payload)
		"magnet_blast":
			_on_magnet_blast(payload)
		"magnet_end":
			_on_magnet_end(payload)

## magnet_start — 磁力場開始
func _on_magnet_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	_player_id = payload.get("player_id", "")
	_duration_sec = payload.get("duration_sec", 12)
	var kill_boost: float = payload.get("kill_boost", 1.8)
	var blast_radius: float = payload.get("blast_radius", 200.0)
	_active = true
	_pull_count = 0

	var vp_size = get_viewport().size

	# 藍色雙閃光
	_flash_screen(COLOR_PRIMARY, 0.16)
	await get_tree().create_timer(0.10).timeout
	_flash_screen(COLOR_DARK, 0.13)

	# 頂部橫幅
	var banner = Label.new()
	banner.text = "🧲 %s 觸發磁力場！所有目標向中央聚集！" % player_name
	banner.add_theme_font_size_override("font_size", 14)
	banner.add_theme_color_override("font_color", COLOR_PALE)
	banner.position = Vector2(vp_size.x / 2 - 160, 6)
	add_child(banner)

	var tween_banner = banner.create_tween()
	tween_banner.tween_interval(float(_duration_sec) - 0.5)
	tween_banner.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween_banner.tween_callback(banner.queue_free)

	# 「🧲 磁力場啟動！」大字
	var big_label = Label.new()
	big_label.text = "🧲 磁力場啟動！"
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
	mult_label.text = "磁力場期間 ×%.1f 倍率加成！" % kill_boost
	mult_label.add_theme_font_size_override("font_size", 13)
	mult_label.add_theme_color_override("font_color", COLOR_GOLD)
	mult_label.position = Vector2(vp_size.x / 2 - 80, vp_size.y / 2 + 28)
	add_child(mult_label)

	var tween_mult = mult_label.create_tween()
	tween_mult.tween_interval(2.0)
	tween_mult.tween_property(mult_label, "modulate:a", 0.0, 0.5)
	tween_mult.tween_callback(mult_label.queue_free)

	# 中央磁力圓圈（視覺提示）
	_spawn_magnet_circle(vp_size, blast_radius)

	# 底部計時條（藍→深藍漸變）
	_spawn_timer_bar(float(_duration_sec))

## magnet_pull — 磁力吸引（每 1.5 秒）
func _on_magnet_pull(payload: Dictionary) -> void:
	var pull_num: int = payload.get("pull_num", 1)
	var moved_count: int = payload.get("moved_count", 0)
	_pull_count = pull_num

	if moved_count == 0:
		return

	var vp_size = get_viewport().size

	# 藍色閃光（輕微）
	_flash_screen(COLOR_PRIMARY, 0.06)

	# 「🧲 磁力吸引！」提示
	var pull_label = Label.new()
	pull_label.text = "🧲 磁力吸引！%d 個目標移動" % moved_count
	pull_label.add_theme_font_size_override("font_size", 13)
	pull_label.add_theme_color_override("font_color", COLOR_PALE)
	pull_label.position = Vector2(vp_size.x / 2 - 70, vp_size.y / 2 - 60)
	add_child(pull_label)

	var tween_p = pull_label.create_tween()
	tween_p.tween_property(pull_label, "position:y", pull_label.position.y - 15, 0.4)
	tween_p.parallel().tween_property(pull_label, "modulate:a", 0.0, 0.4)
	tween_p.tween_callback(pull_label.queue_free)

## magnet_blast — 磁力爆發
func _on_magnet_blast(payload: Dictionary) -> void:
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

	# 全螢幕三次強閃光（藍→白→青）
	_flash_screen(COLOR_PRIMARY, 0.12)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(Color.WHITE, 0.10)
	await get_tree().create_timer(0.07).timeout
	_flash_screen(COLOR_CYAN, 0.12)

	# 「🧲 磁力爆發！」大字
	var blast_label = Label.new()
	blast_label.text = "🧲 磁力爆發！"
	blast_label.add_theme_font_size_override("font_size", 52)
	blast_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	blast_label.position = vp_size / 2 - Vector2(120, 32)
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

## magnet_end — 磁力場結束
func _on_magnet_end(_payload: Dictionary) -> void:
	_active = false
	_pull_count = 0

	# 清除計時條
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		_timer_bar_bg.queue_free()
		_timer_bar_bg = null

	var vp_size = get_viewport().size

	# 藍色淡出提示
	var end_label = Label.new()
	end_label.text = "🧲 磁力場結束"
	end_label.add_theme_font_size_override("font_size", 14)
	end_label.add_theme_color_override("font_color", COLOR_DARK)
	end_label.position = Vector2(vp_size.x / 2 - 55, vp_size.y - 40)
	add_child(end_label)

	var tween_end = end_label.create_tween()
	tween_end.tween_interval(0.5)
	tween_end.tween_property(end_label, "modulate:a", 0.0, 0.5)
	tween_end.tween_callback(end_label.queue_free)

## 建立中央磁力圓圈（視覺提示）
func _spawn_magnet_circle(vp_size: Vector2, blast_radius: float) -> void:
	# 磁力圓圈（藍色輪廓，中央位置）
	# 使用 ColorRect 模擬圓圈（Godot 2D 沒有直接的圓形 UI 節點）
	var scale_x = vp_size.x / 1000.0
	var scale_y = vp_size.y / 600.0
	var center_x = 500.0 * scale_x
	var center_y = 300.0 * scale_y
	var radius_px = blast_radius * (scale_x + scale_y) / 2.0

	# 外圈（藍色半透明）
	var outer_ring = ColorRect.new()
	outer_ring.color = Color(COLOR_PRIMARY.r, COLOR_PRIMARY.g, COLOR_PRIMARY.b, 0.15)
	outer_ring.size = Vector2(radius_px * 2, radius_px * 2)
	outer_ring.position = Vector2(center_x - radius_px, center_y - radius_px)
	add_child(outer_ring)

	# 圓圈脈衝動畫
	var tween_ring = outer_ring.create_tween().set_loops(6)
	tween_ring.tween_property(outer_ring, "modulate:a", 0.5, 0.6)
	tween_ring.tween_property(outer_ring, "modulate:a", 0.1, 0.6)

	# 12 秒後自動清除
	get_tree().create_timer(12.0).timeout.connect(func():
		if is_instance_valid(outer_ring):
			outer_ring.queue_free()
	)

## 建立磁力爆發結算彈窗
func _spawn_blast_result_panel(vp_size: Vector2, killed_count: int, total_reward: int) -> void:
	var panel = ColorRect.new()
	panel.color = Color(0.05, 0.15, 0.35, 0.88)
	panel.size = Vector2(200, 90)
	panel.position = Vector2(vp_size.x + 10, vp_size.y / 2 - 45)
	add_child(panel)

	# 標題
	var title = Label.new()
	title.text = "🧲 磁力爆發結算"
	title.add_theme_font_size_override("font_size", 13)
	title.add_theme_color_override("font_color", COLOR_PALE)
	title.position = Vector2(10, 8)
	panel.add_child(title)

	# 擊破數
	var kill_label = Label.new()
	kill_label.text = "擊破：%d 個目標" % killed_count
	kill_label.add_theme_font_size_override("font_size", 12)
	kill_label.add_theme_color_override("font_color", Color.WHITE)
	kill_label.position = Vector2(10, 30)
	panel.add_child(kill_label)

	# 全服共享獎勵
	var reward_label = Label.new()
	reward_label.text = "全服共享：+%d" % total_reward
	reward_label.add_theme_font_size_override("font_size", 12)
	reward_label.add_theme_color_override("font_color", COLOR_GOLD)
	reward_label.position = Vector2(10, 52)
	panel.add_child(reward_label)

	# 右側滑入動畫
	var tween_panel = panel.create_tween()
	tween_panel.tween_property(panel, "position:x", vp_size.x - 210, 0.3)
	tween_panel.tween_interval(2.5)
	tween_panel.tween_property(panel, "modulate:a", 0.0, 0.4)
	tween_panel.tween_callback(panel.queue_free)

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

	# 計時條（藍→深藍漸變）
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
