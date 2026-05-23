## LuckyQuantumFishPanel.gd — 幸運量子魚系統面板（DAY-228）
## 業界原創「量子疊加態」機制
##
## 視覺設計：
##   - 紫色量子主題（#9B59B6 + #D7BDE2 + #FF00FF + #F5EEF8）
##   - quantum_start：紫色三次強閃光 + 頂部橫幅 + 「⚛️ 量子疊加！」大字
##                    + 底部計時條 + 量子態目標菱形標記（紫/灰交替閃爍）
##   - quantum_collapse：閃光 + 坍縮結果（高倍率紫色/低倍率灰色）+ 浮動文字
##   - quantum_blast：全螢幕三次紫色強閃光 + 「⚛️ 量子爆炸！」大字 + 結算彈窗
extends CanvasLayer

# 量子態狀態
var _active: bool = false
var _timer_bar: Control = null
var _quantum_markers: Dictionary = {}  # targetID → marker node

# 主題顏色
const COLOR_PRIMARY   = Color("#9B59B6")  # 紫色
const COLOR_LIGHT     = Color("#D7BDE2")  # 淡紫
const COLOR_MAGENTA   = Color("#FF00FF")  # 洋紅（高倍率）
const COLOR_PALE      = Color("#F5EEF8")  # 極淡紫
const COLOR_GRAY      = Color("#808080")  # 灰色（低倍率）
const COLOR_GOLD      = Color("#FFD700")  # 金黃

func _ready() -> void:
	layer = 17  # 幸運量子魚面板層級

## 處理幸運量子魚訊息
func handle_lucky_quantum_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"quantum_start":
			_on_quantum_start(payload)
		"quantum_collapse":
			_on_quantum_collapse(payload)
		"quantum_blast":
			_on_quantum_blast(payload)

## quantum_start — 量子疊加開始（全服廣播）
func _on_quantum_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var quantum_count: int = payload.get("quantum_count", 0)
	var duration_sec: int = payload.get("duration_sec", 10)
	var high_mult: float = payload.get("high_mult", 3.0)
	var low_mult: float = payload.get("low_mult", 0.8)

	_active = true

	# 紫色三次強閃光
	_flash_screen(COLOR_PRIMARY, 0.18)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_LIGHT, 0.15)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(COLOR_PALE, 0.12)

	var vp_size = get_viewport().size

	# 頂部橫幅
	var banner = Label.new()
	banner.text = "⚛️ %s 觸發量子疊加！%d 個目標進入量子態！" % [player_name, quantum_count]
	banner.add_theme_font_size_override("font_size", 14)
	banner.add_theme_color_override("font_color", COLOR_PRIMARY)
	banner.position = Vector2(vp_size.x / 2 - 180, 6)
	add_child(banner)

	var tween_banner = banner.create_tween()
	tween_banner.tween_interval(4.0)
	tween_banner.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween_banner.tween_callback(banner.queue_free)

	# 「⚛️ 量子疊加！」大字
	var big_label = Label.new()
	big_label.text = "⚛️ 量子疊加！"
	big_label.add_theme_font_size_override("font_size", 52)
	big_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	big_label.position = vp_size / 2 - Vector2(130, 30)
	add_child(big_label)

	var tween_big = big_label.create_tween()
	tween_big.tween_property(big_label, "scale", Vector2(1.25, 1.25), 0.12)
	tween_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_big.tween_interval(0.5)
	tween_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_big.tween_callback(big_label.queue_free)

	# 倍率說明提示
	var hint_label = Label.new()
	hint_label.text = "觀測即坍縮：×%.1f 或 ×%.1f" % [high_mult, low_mult]
	hint_label.add_theme_font_size_override("font_size", 14)
	hint_label.add_theme_color_override("font_color", COLOR_GOLD)
	hint_label.position = Vector2(vp_size.x / 2 - 90, vp_size.y / 2 + 30)
	add_child(hint_label)

	var tween_hint = hint_label.create_tween()
	tween_hint.tween_interval(2.0)
	tween_hint.tween_property(hint_label, "modulate:a", 0.0, 0.5)
	tween_hint.tween_callback(hint_label.queue_free)

	# 底部計時條（紫→洋紅漸變）
	_show_timer_bar(duration_sec)

## quantum_collapse — 量子態坍縮（玩家觀測後）
func _on_quantum_collapse(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var target_id: String = payload.get("target_id", "")
	var collapse_high: bool = payload.get("collapse_high", false)
	var collapse_mult: float = payload.get("collapse_mult", 1.0)
	var collapse_reward: int = payload.get("collapse_reward", 0)

	# 移除量子標記
	if _quantum_markers.has(target_id):
		var marker = _quantum_markers[target_id]
		if is_instance_valid(marker):
			var tween_fade = marker.create_tween()
			tween_fade.tween_property(marker, "modulate:a", 0.0, 0.2)
			tween_fade.tween_callback(marker.queue_free)
		_quantum_markers.erase(target_id)

	var vp_size = get_viewport().size

	if collapse_high:
		# 高倍率坍縮：洋紅閃光 + 大字
		_flash_screen(COLOR_MAGENTA, 0.15)
		var result_label = Label.new()
		result_label.text = "⚛️ 高倍率坍縮！×%.1f" % collapse_mult
		result_label.add_theme_font_size_override("font_size", 28)
		result_label.add_theme_color_override("font_color", COLOR_MAGENTA)
		result_label.position = Vector2(vp_size.x / 2 - 100, vp_size.y / 2 - 20)
		add_child(result_label)

		var tween_r = result_label.create_tween()
		tween_r.tween_property(result_label, "position:y", result_label.position.y - 40, 0.6)
		tween_r.parallel().tween_property(result_label, "modulate:a", 0.0, 0.6)
		tween_r.tween_callback(result_label.queue_free)
	else:
		# 低倍率坍縮：灰色閃光 + 小字
		_flash_screen(COLOR_GRAY, 0.1)
		var result_label = Label.new()
		result_label.text = "低倍率坍縮 ×%.1f" % collapse_mult
		result_label.add_theme_font_size_override("font_size", 18)
		result_label.add_theme_color_override("font_color", COLOR_GRAY)
		result_label.position = Vector2(vp_size.x / 2 - 70, vp_size.y / 2 - 10)
		add_child(result_label)

		var tween_r = result_label.create_tween()
		tween_r.tween_property(result_label, "position:y", result_label.position.y - 25, 0.5)
		tween_r.parallel().tween_property(result_label, "modulate:a", 0.0, 0.5)
		tween_r.tween_callback(result_label.queue_free)

## quantum_blast — 量子爆炸結算
func _on_quantum_blast(payload: Dictionary) -> void:
	var blast_count: int = payload.get("blast_count", 0)
	var total_reward: int = payload.get("total_reward", 0)

	_active = false

	# 清除計時條
	if _timer_bar != null and is_instance_valid(_timer_bar):
		var tween_timer = _timer_bar.create_tween()
		tween_timer.tween_property(_timer_bar, "modulate:a", 0.0, 0.2)
		tween_timer.tween_callback(_timer_bar.queue_free)
		_timer_bar = null

	# 清除所有量子標記
	for tid in _quantum_markers.keys():
		var marker = _quantum_markers[tid]
		if is_instance_valid(marker):
			marker.queue_free()
	_quantum_markers.clear()

	if blast_count == 0:
		return

	# 全螢幕三次紫色強閃光
	_flash_screen(COLOR_MAGENTA, 0.22)
	await get_tree().create_timer(0.14).timeout
	_flash_screen(COLOR_PRIMARY, 0.18)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_LIGHT, 0.15)

	var vp_size = get_viewport().size

	# 「⚛️ 量子爆炸！」大字
	var big_label = Label.new()
	big_label.text = "⚛️ 量子爆炸！"
	big_label.add_theme_font_size_override("font_size", 52)
	big_label.add_theme_color_override("font_color", COLOR_MAGENTA)
	big_label.position = vp_size / 2 - Vector2(130, 30)
	add_child(big_label)

	var tween_big = big_label.create_tween()
	tween_big.tween_property(big_label, "scale", Vector2(1.3, 1.3), 0.12)
	tween_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_big.tween_interval(0.5)
	tween_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_big.tween_callback(big_label.queue_free)

	# 結算彈窗（右側滑入）
	var popup = _create_blast_popup(blast_count, total_reward, vp_size)
	add_child(popup)

	var tween_popup = popup.create_tween()
	tween_popup.tween_property(popup, "position:x", vp_size.x - 220, 0.3)
	tween_popup.tween_interval(2.5)
	tween_popup.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween_popup.tween_callback(popup.queue_free)

## 建立量子爆炸結算彈窗
func _create_blast_popup(blast_count: int, total_reward: int, vp_size: Vector2) -> Control:
	var popup = Control.new()
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2 - 60)

	var bg = ColorRect.new()
	bg.color = Color(0.1, 0.0, 0.15, 0.88)
	bg.size = Vector2(210, 120)
	popup.add_child(bg)

	var border = ColorRect.new()
	border.color = COLOR_MAGENTA
	border.size = Vector2(210, 3)
	border.position = Vector2(0, 0)
	popup.add_child(border)

	var title_lbl = Label.new()
	title_lbl.text = "⚛️ 量子爆炸結算"
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.add_theme_color_override("font_color", COLOR_MAGENTA)
	title_lbl.position = Vector2(10, 10)
	popup.add_child(title_lbl)

	var count_lbl = Label.new()
	count_lbl.text = "擊破目標：%d 個" % blast_count
	count_lbl.add_theme_font_size_override("font_size", 13)
	count_lbl.add_theme_color_override("font_color", COLOR_LIGHT)
	count_lbl.position = Vector2(10, 38)
	popup.add_child(count_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "獎勵：%d" % total_reward
	reward_lbl.add_theme_font_size_override("font_size", 18)
	reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	reward_lbl.position = Vector2(10, 62)
	popup.add_child(reward_lbl)

	return popup

# ---- 輔助函數 ----

## 顯示底部計時條（紫→洋紅漸變）
func _show_timer_bar(duration_sec: int) -> void:
	if _timer_bar != null and is_instance_valid(_timer_bar):
		_timer_bar.queue_free()

	var vp_size = get_viewport().size
	var timer_container = Control.new()
	timer_container.position = Vector2(0, vp_size.y - 6)
	timer_container.size = Vector2(vp_size.x, 6)
	add_child(timer_container)
	_timer_bar = timer_container

	var bar = ColorRect.new()
	bar.color = COLOR_PRIMARY
	bar.size = Vector2(vp_size.x, 6)
	timer_container.add_child(bar)

	var tween = bar.create_tween()
	tween.tween_property(bar, "size:x", 0.0, float(duration_sec))
	tween.parallel().tween_property(bar, "color", COLOR_MAGENTA, float(duration_sec))

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.35)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)
