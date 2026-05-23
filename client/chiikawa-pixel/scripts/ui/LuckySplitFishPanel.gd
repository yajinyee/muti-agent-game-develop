## LuckySplitFishPanel.gd — 幸運分裂魚系統面板（DAY-224）
## 業界原創「一魚分三」機制
##
## 視覺設計：
##   - 橙紅分裂主題（#FF6B35 + #FF4500 + #FFD700 + #FFF0E6）
##   - split_start：橙紅三次強閃光 + 頂部橫幅 + 分裂碎片標記（三角形散佈）
##   - split_blast：全螢幕橙紅三次強閃光 + 「💥 二次爆炸！」大字 + 結算彈窗
##   - split_end：橙色淡出
extends CanvasLayer

# 分裂狀態
var _active: bool = false
var _banner: Control = null

# 主題顏色
const COLOR_PRIMARY   = Color("#FF6B35")  # 橙紅
const COLOR_DARK      = Color("#FF4500")  # 深橙紅
const COLOR_GOLD      = Color("#FFD700")  # 金色
const COLOR_PALE      = Color("#FFF0E6")  # 極淡橙
const COLOR_BG        = Color(0.12, 0.04, 0.0, 0.88)

func _ready() -> void:
	layer = 21  # 幸運分裂魚面板層級

## 處理幸運分裂魚訊息
func handle_lucky_split_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"split_start":
			_on_split_start(payload)
		"split_blast":
			_on_split_blast(payload)
		"split_end":
			_on_split_end(payload)

## split_start — 分裂爆炸開始
func _on_split_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var duration_sec: int = payload.get("duration_sec", 8)
	var kill_mult: float = payload.get("kill_mult", 1.8)
	var fragments: Array = payload.get("fragments", [])
	var origin_x: float = payload.get("origin_x", 0.0)
	var origin_y: float = payload.get("origin_y", 0.0)

	_active = true

	# 橙紅三次強閃光
	_flash_screen(COLOR_PRIMARY, 0.18)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(COLOR_DARK, 0.15)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(COLOR_GOLD, 0.12)

	# 頂部橫幅 + 計時條
	_show_banner(
		"💥 分裂爆炸！",
		"%s 觸發分裂！一魚分三！×%.1f 倍率加成！" % [player_name, kill_mult],
		duration_sec
	)

	# 分裂爆炸中心特效
	_spawn_split_burst(Vector2(origin_x, origin_y))

	# 分裂碎片標記（三角形散佈）
	for frag in fragments:
		var fx: float = frag.get("x", 0.0)
		var fy: float = frag.get("y", 0.0)
		_spawn_fragment_marker(Vector2(fx, fy), kill_mult)

## split_blast — 二次爆炸結算
func _on_split_blast(payload: Dictionary) -> void:
	_active = false
	_hide_banner()

	var blast_count: int = payload.get("blast_count", 0)
	var total_reward: int = payload.get("total_reward", 0)

	if blast_count <= 0:
		return

	# 全螢幕橙紅三次強閃光
	_flash_screen(COLOR_DARK, 0.15)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(COLOR_PRIMARY, 0.12)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(COLOR_GOLD, 0.1)

	# 「💥 二次爆炸！」大字
	var vp_size = get_viewport().size
	var big_label = Label.new()
	big_label.text = "💥 二次爆炸！×%d" % blast_count
	big_label.add_theme_font_size_override("font_size", 48)
	big_label.add_theme_color_override("font_color", COLOR_GOLD)
	big_label.position = vp_size / 2 - Vector2(140, 28)
	add_child(big_label)

	var tween_label = big_label.create_tween()
	tween_label.tween_property(big_label, "scale", Vector2(1.25, 1.25), 0.12)
	tween_label.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_label.tween_interval(0.5)
	tween_label.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_label.tween_callback(big_label.queue_free)

	# 結算彈窗（右側滑入）
	if total_reward > 0:
		_show_result_popup(blast_count, total_reward)

## split_end — 分裂結束（所有碎片都被玩家打掉）
func _on_split_end(_payload: Dictionary) -> void:
	_active = false
	_hide_banner()

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
	title_label.add_theme_color_override("font_color", COLOR_GOLD)
	title_label.position = Vector2(12, 4)
	banner.add_child(title_label)

	var sub_label = Label.new()
	sub_label.text = subtitle
	sub_label.add_theme_font_size_override("font_size", 11)
	sub_label.add_theme_color_override("font_color", COLOR_PALE)
	sub_label.position = Vector2(12, 28)
	banner.add_child(sub_label)

	# 計時條（底部，橙紅→深橙漸變）
	var timer_bg = ColorRect.new()
	timer_bg.color = Color(0.15, 0.05, 0.0, 0.8)
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

## 分裂爆炸中心特效（放射狀爆炸圓圈）
func _spawn_split_burst(pos: Vector2) -> void:
	for i in range(3):
		var burst = ColorRect.new()
		var size = 20.0 + i * 15.0
		burst.color = Color(COLOR_PRIMARY.r, COLOR_PRIMARY.g, COLOR_PRIMARY.b, 0.7 - i * 0.2)
		burst.size = Vector2(size, size)
		burst.position = pos - Vector2(size / 2, size / 2)
		add_child(burst)

		var tween = burst.create_tween()
		tween.tween_property(burst, "scale", Vector2(3.0, 3.0), 0.3)
		tween.parallel().tween_property(burst, "modulate:a", 0.0, 0.3)
		tween.tween_callback(burst.queue_free)

## 分裂碎片標記（橙色菱形 + 倍率標籤）
func _spawn_fragment_marker(pos: Vector2, kill_mult: float) -> void:
	var marker = Control.new()
	marker.position = pos - Vector2(18, 18)
	marker.size = Vector2(36, 36)
	add_child(marker)

	# 菱形輪廓（4個 ColorRect 組成）
	var corners = [
		Vector2(14, 0), Vector2(28, 14), Vector2(14, 28), Vector2(0, 14)
	]
	for j in range(4):
		var dot = ColorRect.new()
		dot.color = COLOR_GOLD
		dot.size = Vector2(4, 4)
		dot.position = corners[j]
		marker.add_child(dot)

	# 倍率標籤
	var mult_label = Label.new()
	mult_label.text = "×%.1f" % kill_mult
	mult_label.add_theme_font_size_override("font_size", 11)
	mult_label.add_theme_color_override("font_color", COLOR_GOLD)
	mult_label.position = Vector2(-8, -16)
	marker.add_child(mult_label)

	# 閃爍動畫
	var tween = marker.create_tween().set_loops(4)
	tween.tween_property(marker, "modulate:a", 0.4, 0.3)
	tween.tween_property(marker, "modulate:a", 1.0, 0.3)

	# 8 秒後淡出
	var fade_timer = get_tree().create_timer(7.5)
	fade_timer.timeout.connect(func():
		if is_instance_valid(marker):
			var fade = marker.create_tween()
			fade.tween_property(marker, "modulate:a", 0.0, 0.5)
			fade.tween_callback(marker.queue_free)
	)

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.35)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)

## 結算彈窗（右側滑入）
func _show_result_popup(blast_count: int, total_reward: int) -> void:
	var vp_size = get_viewport().size
	var popup = Control.new()
	popup.size = Vector2(220, 90)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2 - 45)
	add_child(popup)

	var bg = ColorRect.new()
	bg.color = COLOR_BG
	bg.size = popup.size
	popup.add_child(bg)

	var border = ColorRect.new()
	border.color = COLOR_GOLD
	border.size = Vector2(popup.size.x, 3)
	popup.add_child(border)

	var title_label = Label.new()
	title_label.text = "💥 分裂二次爆炸"
	title_label.add_theme_font_size_override("font_size", 14)
	title_label.add_theme_color_override("font_color", COLOR_GOLD)
	title_label.position = Vector2(8, 8)
	popup.add_child(title_label)

	var blast_label = Label.new()
	blast_label.text = "爆炸碎片：%d 個" % blast_count
	blast_label.add_theme_font_size_override("font_size", 12)
	blast_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	blast_label.position = Vector2(8, 32)
	popup.add_child(blast_label)

	var reward_label = Label.new()
	reward_label.text = "獎勵：%d 金幣" % total_reward
	reward_label.add_theme_font_size_override("font_size", 12)
	reward_label.add_theme_color_override("font_color", COLOR_PALE)
	reward_label.position = Vector2(8, 56)
	popup.add_child(reward_label)

	# 右側滑入動畫
	var tween = popup.create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 230.0, 0.3)
	tween.tween_interval(3.0)
	tween.tween_property(popup, "position:x", vp_size.x + 10.0, 0.3)
	tween.tween_callback(popup.queue_free)
