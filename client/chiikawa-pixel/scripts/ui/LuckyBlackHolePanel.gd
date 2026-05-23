## LuckyBlackHolePanel.gd — 幸運黑洞魚系統面板（DAY-221）
## 業界原創「重力黑洞」機制
##
## 視覺設計：
##   - 深紫黑洞主題（#8B00FF + #4B0082 + #FF00FF + #E6E6FA）
##   - blackhole_start：紫色三次強閃光 + 頂部橫幅 + 黑洞圓圈（收縮脈衝）+ 底部計時條
##   - blackhole_pulse：黑洞圓圈脈衝閃爍 + 浮動「重力傷害」文字
##   - singularity_blast：全螢幕黑色收縮 + 紫色三次強閃光 + 「🌑 奇點爆炸！」52px大字
##   - singularity_hit：小爆炸圓圈 + 浮動獎勵文字
##   - singularity_result：右側滑入結算彈窗
extends CanvasLayer

# 黑洞狀態
var _blackhole_active: bool = false
var _blackhole_circle: Control = null
var _timer_bar: Control = null
var _banner: Control = null
var _blackhole_x: float = 0.0
var _blackhole_y: float = 0.0
var _blackhole_radius: float = 350.0

# 主題顏色
const COLOR_PRIMARY   = Color("#8B00FF")  # 深紫
const COLOR_DARK      = Color("#4B0082")  # 靛藍
const COLOR_ACCENT    = Color("#FF00FF")  # 洋紅
const COLOR_LIGHT     = Color("#E6E6FA")  # 薰衣草白
const COLOR_BG        = Color(0.05, 0.0, 0.1, 0.88)

func _ready() -> void:
	layer = 24  # 幸運黑洞魚面板層級

## 處理幸運黑洞魚訊息
func handle_lucky_black_hole(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"blackhole_start":
			_on_blackhole_start(payload)
		"blackhole_pulse":
			_on_blackhole_pulse(payload)
		"singularity_blast":
			_on_singularity_blast(payload)
		"singularity_hit":
			_on_singularity_hit(payload)
		"singularity_result":
			_on_singularity_result(payload)

## blackhole_start — 黑洞建立
func _on_blackhole_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var duration_sec: int = payload.get("duration_sec", 10)
	var x: float = payload.get("x", 640.0)
	var y: float = payload.get("y", 360.0)
	var radius: float = payload.get("radius", 350.0)

	_blackhole_active = true
	_blackhole_x = x
	_blackhole_y = y
	_blackhole_radius = radius

	# 紫色三次強閃光
	_flash_screen(COLOR_PRIMARY, 0.2)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_DARK, 0.18)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_ACCENT, 0.15)

	# 建立黑洞圓圈視覺
	_create_blackhole_circle(Vector2(x, y), radius)

	# 頂部橫幅
	_show_banner("🌑 黑洞召喚！", "%s 召喚重力黑洞！10 秒後奇點爆炸！" % player_name, duration_sec)

## blackhole_pulse — 重力脈衝（每秒）
func _on_blackhole_pulse(payload: Dictionary) -> void:
	var affected_count: int = payload.get("affected_count", 0)
	var pulse_num: int = payload.get("pulse_num", 1)

	# 黑洞圓圈脈衝閃爍
	if _blackhole_circle != null and is_instance_valid(_blackhole_circle):
		var tween = _blackhole_circle.create_tween()
		tween.tween_property(_blackhole_circle, "modulate:a", 1.0, 0.1)
		tween.tween_property(_blackhole_circle, "modulate:a", 0.6, 0.2)

	# 浮動文字（重力傷害）
	if affected_count > 0:
		var pos = Vector2(_blackhole_x, _blackhole_y - _blackhole_radius * 0.5)
		_show_float_text("🌑 重力 -%d" % affected_count, COLOR_ACCENT, pos)

## singularity_blast — 奇點爆炸開始
func _on_singularity_blast(_payload: Dictionary) -> void:
	_blackhole_active = false
	_hide_banner()

	# 移除黑洞圓圈
	if _blackhole_circle != null and is_instance_valid(_blackhole_circle):
		var tween = _blackhole_circle.create_tween()
		tween.tween_property(_blackhole_circle, "scale", Vector2(0.1, 0.1), 0.3)
		tween.parallel().tween_property(_blackhole_circle, "modulate:a", 0.0, 0.3)
		tween.tween_callback(_blackhole_circle.queue_free)
		_blackhole_circle = null

	# 全螢幕黑色收縮效果
	var black_overlay = ColorRect.new()
	black_overlay.color = Color(0, 0, 0, 0.0)
	black_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(black_overlay)

	var tween_black = black_overlay.create_tween()
	tween_black.tween_property(black_overlay, "color:a", 0.7, 0.2)
	tween_black.tween_property(black_overlay, "color:a", 0.0, 0.15)
	tween_black.tween_callback(black_overlay.queue_free)

	await get_tree().create_timer(0.1).timeout

	# 紫色三次強閃光
	_flash_screen(COLOR_PRIMARY, 0.15)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color.WHITE, 0.12)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(COLOR_ACCENT, 0.15)

	# 「🌑 奇點爆炸！」大字
	var big_label = Label.new()
	big_label.text = "🌑 奇點爆炸！"
	big_label.add_theme_font_size_override("font_size", 52)
	big_label.add_theme_color_override("font_color", COLOR_ACCENT)
	big_label.set_anchors_preset(Control.PRESET_CENTER)
	var vp_size = get_viewport().size
	big_label.position = vp_size / 2 - Vector2(160, 30)
	add_child(big_label)

	var tween_label = big_label.create_tween()
	tween_label.tween_property(big_label, "scale", Vector2(1.3, 1.3), 0.15)
	tween_label.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.1)
	tween_label.tween_interval(0.5)
	tween_label.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_label.tween_callback(big_label.queue_free)

## singularity_hit — 單個目標被奇點擊破
func _on_singularity_hit(payload: Dictionary) -> void:
	var reward: int = payload.get("reward", 0)
	var x: float = payload.get("x", 0.0)
	var y: float = payload.get("y", 0.0)

	# 小爆炸圓圈（紫色）
	_show_explosion_circle(Vector2(x, y), 25.0)

	# 浮動獎勵文字
	if reward > 0:
		_show_float_text("+%d" % reward, COLOR_LIGHT, Vector2(x, y))

## singularity_result — 奇點爆炸結算
func _on_singularity_result(payload: Dictionary) -> void:
	var killed_count: int = payload.get("killed_count", 0)
	var total_reward: int = payload.get("total_reward", 0)

	if killed_count <= 0:
		return

	# 右側滑入結算彈窗
	_show_result_popup(killed_count, total_reward)

# ---- 輔助函數 ----

## 建立黑洞圓圈視覺（三層同心圓 + 脈衝動畫）
func _create_blackhole_circle(center: Vector2, radius: float) -> void:
	if _blackhole_circle != null and is_instance_valid(_blackhole_circle):
		_blackhole_circle.queue_free()

	var container = Control.new()
	container.position = center - Vector2(radius, radius)
	container.size = Vector2(radius * 2, radius * 2)
	add_child(container)
	_blackhole_circle = container

	# 外圈（深紫，半透明）
	var outer = ColorRect.new()
	outer.color = Color(COLOR_PRIMARY.r, COLOR_PRIMARY.g, COLOR_PRIMARY.b, 0.15)
	outer.size = Vector2(radius * 2, radius * 2)
	container.add_child(outer)

	# 中圈（靛藍，更不透明）
	var mid_r = radius * 0.65
	var mid = ColorRect.new()
	mid.color = Color(COLOR_DARK.r, COLOR_DARK.g, COLOR_DARK.b, 0.25)
	mid.size = Vector2(mid_r * 2, mid_r * 2)
	mid.position = Vector2(radius - mid_r, radius - mid_r)
	container.add_child(mid)

	# 內圈（黑色核心）
	var inner_r = radius * 0.25
	var inner = ColorRect.new()
	inner.color = Color(0.0, 0.0, 0.0, 0.85)
	inner.size = Vector2(inner_r * 2, inner_r * 2)
	inner.position = Vector2(radius - inner_r, radius - inner_r)
	container.add_child(inner)

	# 脈衝動畫（sin 波動透明度）
	var tween = container.create_tween().set_loops()
	tween.tween_property(container, "modulate:a", 0.5, 0.8)
	tween.tween_property(container, "modulate:a", 1.0, 0.8)

## 顯示頂部橫幅
func _show_banner(title: String, subtitle: String, duration_sec: int) -> void:
	_hide_banner()

	var banner = Control.new()
	banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	banner.position = Vector2(0, 8)
	banner.size = Vector2(get_viewport().size.x, 52)
	add_child(banner)

	var bg = ColorRect.new()
	bg.color = COLOR_BG
	bg.size = banner.size
	banner.add_child(bg)

	var title_label = Label.new()
	title_label.text = title
	title_label.add_theme_font_size_override("font_size", 20)
	title_label.add_theme_color_override("font_color", COLOR_ACCENT)
	title_label.position = Vector2(12, 4)
	banner.add_child(title_label)

	var sub_label = Label.new()
	sub_label.text = subtitle
	sub_label.add_theme_font_size_override("font_size", 13)
	sub_label.add_theme_color_override("font_color", COLOR_LIGHT)
	sub_label.position = Vector2(12, 28)
	banner.add_child(sub_label)

	# 計時條（底部，深紫→洋紅漸變）
	var timer_bar = ColorRect.new()
	timer_bar.name = "TimerBar"
	timer_bar.color = COLOR_PRIMARY
	timer_bar.position = Vector2(0, 48)
	timer_bar.size = Vector2(get_viewport().size.x, 4)
	banner.add_child(timer_bar)

	var tween = banner.create_tween()
	tween.tween_property(timer_bar, "size:x", 0.0, float(duration_sec))
	tween.parallel().tween_property(timer_bar, "color", COLOR_ACCENT, float(duration_sec))

	_banner = banner

## 隱藏橫幅
func _hide_banner() -> void:
	if _banner != null and is_instance_valid(_banner):
		_banner.queue_free()
	_banner = null

## 爆炸圓圈效果
func _show_explosion_circle(pos: Vector2, radius: float) -> void:
	var circle = ColorRect.new()
	circle.color = Color(COLOR_PRIMARY.r, COLOR_PRIMARY.g, COLOR_PRIMARY.b, 0.7)
	circle.size = Vector2(radius * 2, radius * 2)
	circle.position = pos - Vector2(radius, radius)
	add_child(circle)

	var tween = circle.create_tween()
	tween.tween_property(circle, "size", Vector2(radius * 4, radius * 4), 0.35)
	tween.parallel().tween_property(circle, "position", pos - Vector2(radius * 2, radius * 2), 0.35)
	tween.parallel().tween_property(circle, "modulate:a", 0.0, 0.35)
	tween.tween_callback(circle.queue_free)

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.45)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)

## 浮動文字
func _show_float_text(text: String, color: Color, pos: Vector2) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 16)
	label.add_theme_color_override("font_color", color)
	label.position = pos - Vector2(30, 15)
	add_child(label)

	var tween = label.create_tween()
	tween.tween_property(label, "position:y", label.position.y - 40, 0.8)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.8)
	tween.tween_callback(label.queue_free)

## 結算彈窗（右側滑入）
func _show_result_popup(killed_count: int, total_reward: int) -> void:
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
	border.color = COLOR_ACCENT
	border.size = Vector2(popup.size.x, 3)
	border.position = Vector2(0, 0)
	popup.add_child(border)

	var title_label = Label.new()
	title_label.text = "🌑 奇點爆炸結算"
	title_label.add_theme_font_size_override("font_size", 14)
	title_label.add_theme_color_override("font_color", COLOR_ACCENT)
	title_label.position = Vector2(8, 8)
	popup.add_child(title_label)

	var killed_label = Label.new()
	killed_label.text = "消滅目標：%d 個" % killed_count
	killed_label.add_theme_font_size_override("font_size", 13)
	killed_label.add_theme_color_override("font_color", COLOR_LIGHT)
	killed_label.position = Vector2(8, 32)
	popup.add_child(killed_label)

	var reward_label = Label.new()
	reward_label.text = "全服獎勵：%d 金幣" % total_reward
	reward_label.add_theme_font_size_override("font_size", 13)
	reward_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	reward_label.position = Vector2(8, 56)
	popup.add_child(reward_label)

	# 右側滑入動畫
	var tween = popup.create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 230.0, 0.3)
	tween.tween_interval(2.5)
	tween.tween_property(popup, "position:x", vp_size.x + 10.0, 0.3)
	tween.tween_callback(popup.queue_free)
