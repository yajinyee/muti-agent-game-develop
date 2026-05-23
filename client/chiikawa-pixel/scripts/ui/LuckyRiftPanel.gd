## LuckyRiftPanel.gd — 幸運時空裂縫魚系統面板（DAY-255）
## 業界原創「時空裂縫+傳送吸入+裂縫崩塌」機制
##
## 視覺設計：
##   - 深紫時空主題（#6A0DAD + #9B59B6 + #E8DAEF + #4B0082）
##   - rift_start：深紫三次強閃光 + 頂部橫幅 + 「🌀 時空裂縫！」大字 + 裂縫旋轉圓圈 + 計時條
##   - rift_broadcast：頂部小橫幅（全服廣播）
##   - rift_suck：深紫閃光 + 「🌀 第N次吸入！[目標名] 傳送！×1.6」浮動文字 + 吸入計數器
##   - rift_collapse：全螢幕三次強閃光 + 「🌀 裂縫崩塌！×2.5」大字 + 結算彈窗
extends CanvasLayer

# 主題顏色
const COLOR_RIFT_PURPLE = Color("#6A0DAD")  # 深紫（主色）
const COLOR_VIOLET      = Color("#9B59B6")  # 紫羅蘭（副色）
const COLOR_INDIGO      = Color("#4B0082")  # 靛藍（崩塌）
const COLOR_LAVENDER    = Color("#E8DAEF")  # 薰衣草（副文字）
const COLOR_WHITE       = Color("#FFFFFF")  # 白色

# 計時條
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null

# 裂縫旋轉圓圈
var _rift_circle: Control = null

# 吸入計數器
var _suck_counter: Label = null

func _ready() -> void:
	layer = 28  # 幸運時空裂縫魚面板層級（DAY-255）

## 處理幸運時空裂縫魚訊息
func handle_lucky_rift(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"rift_start":
			_on_rift_start(payload)
		"rift_broadcast":
			_on_rift_broadcast(payload)
		"rift_suck":
			_on_rift_suck(payload)
		"rift_collapse":
			_on_rift_collapse(payload)

## rift_start — 裂縫啟動（個人訊息）
func _on_rift_start(payload: Dictionary) -> void:
	var duration_sec: int = payload.get("duration_sec", 18)
	var suck_mult: float = payload.get("suck_mult", 1.6)
	var collapse_mult: float = payload.get("collapse_mult", 2.5)

	# 深紫三次強閃光
	_flash_screen(COLOR_RIFT_PURPLE, 0.6, 3)

	# 頂部橫幅
	_show_banner("🌀 時空裂縫！每 3 秒吸入最近目標傳送！×%.1f 倍率！全服共享！" % suck_mult, COLOR_RIFT_PURPLE, 4.0)

	# 中央大字
	_show_big_text("🌀 時空裂縫！", COLOR_RIFT_PURPLE, 52, 2.5)
	_show_sub_text("吸入 ×%.1f  崩塌 ×%.1f（全服 AOE）  最多吸入 5 個目標！" % [suck_mult, collapse_mult], COLOR_LAVENDER, 2.5)

	# 裂縫旋轉圓圈（場景中央）
	_show_rift_circle()

	# 吸入計數器
	_show_suck_counter(0, 5)

	# 右側豎向計時條（x=-156 與其他計時條錯開）
	_start_timer_bar(duration_sec, COLOR_RIFT_PURPLE)

## rift_broadcast — 全服廣播裂縫
func _on_rift_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var suck_mult: float = payload.get("suck_mult", 1.6)
	_show_top_banner("🌀 %s 開啟時空裂縫！吸入 ×%.1f 全服共享！" % [player_name, suck_mult], COLOR_RIFT_PURPLE, 2.5)

## rift_suck — 裂縫吸入（全服廣播）
func _on_rift_suck(payload: Dictionary) -> void:
	var suck_num: int = payload.get("suck_num", 1)
	var suck_count: int = payload.get("suck_count", 1)
	var max_suck: int = payload.get("max_suck", 5)
	var target_name: String = payload.get("target_name", "目標")
	var mult: float = payload.get("mult", 1.6)
	var total_reward: int = payload.get("total_reward", 0)

	# 深紫閃光（輕微）
	_flash_screen(COLOR_RIFT_PURPLE, 0.15, 1)

	# 浮動文字
	_show_float_text("🌀 第%d次吸入！%s 傳送！×%.1f  全服+%d" % [suck_num, target_name, mult, total_reward], COLOR_VIOLET, 1.8)

	# 更新吸入計數器
	_update_suck_counter(suck_count, max_suck)

## rift_collapse — 裂縫崩塌（全服廣播）
func _on_rift_collapse(payload: Dictionary) -> void:
	var drain_count: int = payload.get("drain_count", 0)
	var collapse_mult: float = payload.get("collapse_mult", 2.5)
	var total_reward: int = payload.get("total_reward", 0)
	var suck_count: int = payload.get("suck_count", 0)

	# 停止計時條和裂縫圓圈
	_stop_timer_bar()
	_clear_rift_circle()
	_clear_suck_counter()

	# 全螢幕三次強閃光（靛藍，崩塌感）
	_flash_screen(COLOR_INDIGO, 0.75, 3)

	# 大字
	_show_big_text("🌀 裂縫崩塌！", COLOR_INDIGO, 52, 2.5)
	_show_sub_text("%d 個目標 HP -50%%！全服 AOE ×%.1f！共吸入 %d 個目標！" % [drain_count, collapse_mult, suck_count], COLOR_LAVENDER, 2.5)

	# 結算彈窗
	_show_collapse_popup(drain_count, collapse_mult, total_reward, suck_count)

# ─── 裂縫旋轉圓圈 ────────────────────────────────────────────────────────────

func _show_rift_circle() -> void:
	var vp_size = get_viewport().size
	_rift_circle = Control.new()
	_rift_circle.size = Vector2(120, 120)
	_rift_circle.position = Vector2(vp_size.x / 2.0 - 60, vp_size.y / 2.0 - 60)
	add_child(_rift_circle)

	# 多層同心圓（用 ColorRect 模擬）
	for i in range(3):
		var ring = ColorRect.new()
		var ring_size = 120.0 - i * 30.0
		ring.size = Vector2(ring_size, ring_size)
		ring.position = Vector2((120.0 - ring_size) / 2.0, (120.0 - ring_size) / 2.0)
		ring.color = Color(COLOR_RIFT_PURPLE.r, COLOR_RIFT_PURPLE.g, COLOR_RIFT_PURPLE.b, 0.3 - i * 0.08)
		_rift_circle.add_child(ring)

	# 中心 🌀 標記
	var center_lbl = Label.new()
	center_lbl.text = "🌀"
	center_lbl.add_theme_font_size_override("font_size", 36)
	center_lbl.position = Vector2(30, 30)
	_rift_circle.add_child(center_lbl)

	# 持續旋轉動畫
	var tween = _rift_circle.create_tween().set_loops()
	tween.tween_property(_rift_circle, "rotation", TAU, 3.0).set_ease(Tween.EASE_IN_OUT)

	# 脈衝縮放動畫
	var pulse_tween = _rift_circle.create_tween().set_loops()
	pulse_tween.tween_property(_rift_circle, "scale", Vector2(1.12, 1.12), 0.6)
	pulse_tween.tween_property(_rift_circle, "scale", Vector2(1.0, 1.0), 0.6)

func _clear_rift_circle() -> void:
	if is_instance_valid(_rift_circle):
		var tween = create_tween()
		tween.tween_property(_rift_circle, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_rift_circle.queue_free)
		_rift_circle = null

# ─── 吸入計數器 ──────────────────────────────────────────────────────────────

func _show_suck_counter(current: int, max_count: int) -> void:
	var vp_size = get_viewport().size
	_suck_counter = Label.new()
	_suck_counter.text = "🌀 吸入 %d/%d" % [current, max_count]
	_suck_counter.add_theme_font_size_override("font_size", 18)
	_suck_counter.add_theme_color_override("font_color", COLOR_VIOLET)
	_suck_counter.position = Vector2(vp_size.x - 140, vp_size.y * 0.25 - 30)
	add_child(_suck_counter)

	# 脈衝動畫
	var tween = _suck_counter.create_tween().set_loops()
	tween.tween_property(_suck_counter, "modulate:a", 0.6, 0.4)
	tween.tween_property(_suck_counter, "modulate:a", 1.0, 0.4)

func _update_suck_counter(current: int, max_count: int) -> void:
	if is_instance_valid(_suck_counter):
		_suck_counter.text = "🌀 吸入 %d/%d" % [current, max_count]
		# 閃爍提示更新
		var tween = _suck_counter.create_tween()
		tween.tween_property(_suck_counter, "scale", Vector2(1.2, 1.2), 0.1)
		tween.tween_property(_suck_counter, "scale", Vector2(1.0, 1.0), 0.1)

func _clear_suck_counter() -> void:
	if is_instance_valid(_suck_counter):
		var tween = create_tween()
		tween.tween_property(_suck_counter, "modulate:a", 0.0, 0.4)
		tween.tween_callback(_suck_counter.queue_free)
		_suck_counter = null

# ─── 裂縫崩塌結算彈窗 ────────────────────────────────────────────────────────

func _show_collapse_popup(drain_count: int, collapse_mult: float, total_reward: int, suck_count: int) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(320, 180)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 90)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.04, 0.0, 0.08, 0.93)
	style.border_color = COLOR_INDIGO
	style.set_border_width_all(3)
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 8)
	popup.add_child(vbox)

	var title_lbl = Label.new()
	title_lbl.text = "🌀 裂縫崩塌結算"
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", COLOR_INDIGO)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var suck_lbl = Label.new()
	suck_lbl.text = "共吸入傳送 %d 個目標" % suck_count
	suck_lbl.add_theme_font_size_override("font_size", 14)
	suck_lbl.add_theme_color_override("font_color", COLOR_LAVENDER)
	suck_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(suck_lbl)

	var drain_lbl = Label.new()
	drain_lbl.text = "%d 個目標 HP -50%%" % drain_count
	drain_lbl.add_theme_font_size_override("font_size", 15)
	drain_lbl.add_theme_color_override("font_color", COLOR_LAVENDER)
	drain_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(drain_lbl)

	var mult_lbl = Label.new()
	mult_lbl.text = "全服 AOE ×%.1f" % collapse_mult
	mult_lbl.add_theme_font_size_override("font_size", 18)
	mult_lbl.add_theme_color_override("font_color", COLOR_VIOLET)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "全服獎勵：+%d" % total_reward
	reward_lbl.add_theme_font_size_override("font_size", 16)
	reward_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 340.0, 0.4).set_ease(Tween.EASE_OUT)
	tween.tween_interval(4.5)
	tween.tween_property(popup, "position:x", vp_size.x + 10.0, 0.4).set_ease(Tween.EASE_IN)
	tween.tween_callback(popup.queue_free)

# ─── 通用 UI 工具 ─────────────────────────────────────────────────────────────

func _flash_screen(color: Color, alpha: float, times: int) -> void:
	var vp_size = get_viewport().size
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.0)
	flash.size = vp_size
	add_child(flash)
	var tween = create_tween()
	for i in range(times):
		tween.tween_property(flash, "color:a", alpha, 0.08)
		tween.tween_property(flash, "color:a", 0.0, 0.12)
	tween.tween_callback(flash.queue_free)

func _show_banner(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var banner = PanelContainer.new()
	banner.size = Vector2(vp_size.x, 48)
	banner.position = Vector2(0, 0)
	add_child(banner)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(color.r, color.g, color.b, 0.88)
	banner.add_theme_stylebox_override("panel", style)
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 16)
	lbl.add_theme_color_override("font_color", COLOR_WHITE)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	banner.add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(banner.queue_free)

func _show_top_banner(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var banner = PanelContainer.new()
	banner.size = Vector2(vp_size.x * 0.7, 36)
	banner.position = Vector2(vp_size.x * 0.15, 52)
	add_child(banner)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(color.r, color.g, color.b, 0.82)
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	banner.add_theme_stylebox_override("panel", style)
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 13)
	lbl.add_theme_color_override("font_color", COLOR_WHITE)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	banner.add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(banner.queue_free)

func _show_big_text(text: String, color: Color, font_size: int, duration: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", font_size)
	lbl.add_theme_color_override("font_color", color)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.size = Vector2(vp_size.x, 80)
	lbl.position = Vector2(0, vp_size.y * 0.35)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "scale", Vector2(1.15, 1.15), 0.15).set_ease(Tween.EASE_OUT)
	tween.tween_property(lbl, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(duration - 0.5)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.5)
	tween.tween_callback(lbl.queue_free)

func _show_sub_text(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 13)
	lbl.add_theme_color_override("font_color", color)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.size = Vector2(vp_size.x, 40)
	lbl.position = Vector2(0, vp_size.y * 0.35 + 80)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.4)
	tween.tween_callback(lbl.queue_free)

func _show_float_text(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 17)
	lbl.add_theme_color_override("font_color", color)
	lbl.position = Vector2(
		randf_range(vp_size.x * 0.3, vp_size.x * 0.7),
		randf_range(vp_size.y * 0.3, vp_size.y * 0.6)
	)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 55, duration)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, duration)
	tween.tween_callback(lbl.queue_free)

func _start_timer_bar(duration_sec: int, color: Color) -> void:
	var vp_size = get_viewport().size
	var bar_h: float = vp_size.y * 0.5
	var bar_w: float = 12.0
	# x=-156 與其他計時條錯開（龍捲風-58，黑洞-72，鏡像分裂-86，量子糾纏-100，武器進化-114，隕石雨-128，龍王-142，時空裂縫-156）
	var bar_x: float = vp_size.x - 156.0
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.04, 0.0, 0.08, 0.7)
	_timer_bar_bg.size = Vector2(bar_w, bar_h)
	_timer_bar_bg.position = Vector2(bar_x, bar_y)
	add_child(_timer_bar_bg)

	_timer_bar = ColorRect.new()
	_timer_bar.color = color
	_timer_bar.size = Vector2(bar_w, bar_h)
	_timer_bar.position = Vector2(bar_x, bar_y)
	add_child(_timer_bar)

	_timer_tween = create_tween()
	_timer_tween.tween_property(_timer_bar, "size:y", 0.0, float(duration_sec)).set_ease(Tween.EASE_IN_OUT)

func _stop_timer_bar() -> void:
	if is_instance_valid(_timer_tween):
		_timer_tween.kill()
	if is_instance_valid(_timer_bar):
		var tween = create_tween()
		tween.tween_property(_timer_bar, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_timer_bar.queue_free)
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		var tween2 = create_tween()
		tween2.tween_property(_timer_bar_bg, "modulate:a", 0.0, 0.5)
		tween2.tween_callback(_timer_bar_bg.queue_free)
		_timer_bar_bg = null
