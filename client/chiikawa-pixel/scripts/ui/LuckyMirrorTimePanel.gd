## LuckyMirrorTimePanel.gd — 幸運鏡像時空魚系統面板（DAY-227）
## 業界原創「時間倒流」機制
##
## 視覺設計：
##   - 天藍時空主題（#00BFFF + #87CEEB + #E0F8FF + #FF6B35）
##   - time_rewind_start：全螢幕天藍三次強閃光 + 頂部橫幅 + 「⏪ 時間倒流！」大字
##                        + 底部計時條 + 「×2.0 倍率加成」提示
##   - time_collapse：全螢幕橙色閃光 + 「💥 時間崩潰！」大字 + HP -40% 提示
extends CanvasLayer

# 時間倒流狀態
var _active: bool = false
var _timer_bar: Control = null
var _boost_label: Control = null

# 主題顏色
const COLOR_PRIMARY   = Color("#00BFFF")  # 天藍
const COLOR_LIGHT     = Color("#87CEEB")  # 淡藍
const COLOR_PALE      = Color("#E0F8FF")  # 極淡藍
const COLOR_COLLAPSE  = Color("#FF6B35")  # 橙色（崩潰）
const COLOR_GOLD      = Color("#FFD700")  # 金黃

func _ready() -> void:
	layer = 18  # 幸運鏡像時空魚面板層級

## 處理幸運鏡像時空魚訊息
func handle_lucky_mirror_time(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"time_rewind_start":
			_on_time_rewind_start(payload)
		"time_collapse":
			_on_time_collapse(payload)

## time_rewind_start — 時間倒流開始（全服廣播）
func _on_time_rewind_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var boost_mult: float = payload.get("boost_mult", 2.0)
	var duration_sec: int = payload.get("duration_sec", 8)
	var rewind_count: int = payload.get("rewind_count", 0)

	_active = true

	# 全螢幕天藍三次強閃光
	_flash_screen(COLOR_PRIMARY, 0.18)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_LIGHT, 0.15)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(COLOR_PALE, 0.12)

	var vp_size = get_viewport().size

	# 頂部橫幅
	var banner = Label.new()
	banner.text = "⏪ %s 觸發時間倒流！%d 個目標 HP 回滿！" % [player_name, rewind_count]
	banner.add_theme_font_size_override("font_size", 14)
	banner.add_theme_color_override("font_color", COLOR_PRIMARY)
	banner.position = Vector2(vp_size.x / 2 - 170, 6)
	add_child(banner)

	var tween_banner = banner.create_tween()
	tween_banner.tween_interval(3.5)
	tween_banner.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween_banner.tween_callback(banner.queue_free)

	# 「⏪ 時間倒流！」大字
	var big_label = Label.new()
	big_label.text = "⏪ 時間倒流！"
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

	# 「×2.0 倍率加成」提示（持續顯示）
	if _boost_label != null and is_instance_valid(_boost_label):
		_boost_label.queue_free()

	var boost_label = Label.new()
	boost_label.text = "⏪ ×%.1f 倍率加成中！" % boost_mult
	boost_label.add_theme_font_size_override("font_size", 16)
	boost_label.add_theme_color_override("font_color", COLOR_GOLD)
	boost_label.position = Vector2(vp_size.x / 2 - 80, vp_size.y - 50)
	add_child(boost_label)
	_boost_label = boost_label

	# 閃爍動畫
	var tween_boost = boost_label.create_tween().set_loops()
	tween_boost.tween_property(boost_label, "modulate", Color(1.5, 1.5, 0.5, 1.0), 0.5)
	tween_boost.tween_property(boost_label, "modulate", Color.WHITE, 0.5)

	# 底部計時條（天藍→橙色漸變）
	_show_timer_bar(duration_sec)

## time_collapse — 時間崩潰
func _on_time_collapse(payload: Dictionary) -> void:
	var collapse_count: int = payload.get("collapse_count", 0)

	_active = false

	# 清除倍率提示
	if _boost_label != null and is_instance_valid(_boost_label):
		var tween_fade = _boost_label.create_tween()
		tween_fade.tween_property(_boost_label, "modulate:a", 0.0, 0.3)
		tween_fade.tween_callback(_boost_label.queue_free)
		_boost_label = null

	# 清除計時條
	if _timer_bar != null and is_instance_valid(_timer_bar):
		var tween_timer = _timer_bar.create_tween()
		tween_timer.tween_property(_timer_bar, "modulate:a", 0.0, 0.2)
		tween_timer.tween_callback(_timer_bar.queue_free)
		_timer_bar = null

	# 全螢幕橙色閃光
	_flash_screen(COLOR_COLLAPSE, 0.22)

	var vp_size = get_viewport().size

	# 「💥 時間崩潰！」大字
	var big_label = Label.new()
	big_label.text = "💥 時間崩潰！"
	big_label.add_theme_font_size_override("font_size", 48)
	big_label.add_theme_color_override("font_color", COLOR_COLLAPSE)
	big_label.position = vp_size / 2 - Vector2(120, 28)
	add_child(big_label)

	var tween_big = big_label.create_tween()
	tween_big.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.1)
	tween_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_big.tween_interval(0.4)
	tween_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_big.tween_callback(big_label.queue_free)

	# HP -40% 提示
	var hp_label = Label.new()
	hp_label.text = "⚡ %d 個目標 HP -40%%！" % collapse_count
	hp_label.add_theme_font_size_override("font_size", 18)
	hp_label.add_theme_color_override("font_color", COLOR_GOLD)
	hp_label.position = vp_size / 2 - Vector2(100, -20)
	add_child(hp_label)

	var tween_hp = hp_label.create_tween()
	tween_hp.tween_property(hp_label, "position:y", hp_label.position.y - 30, 0.6)
	tween_hp.parallel().tween_property(hp_label, "modulate:a", 0.0, 0.6)
	tween_hp.tween_callback(hp_label.queue_free)

# ---- 輔助函數 ----

## 顯示底部計時條（天藍→橙色漸變）
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
	tween.parallel().tween_property(bar, "color", COLOR_COLLAPSE, float(duration_sec))

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.35)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)
