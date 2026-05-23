## LuckyZodiacFatePanel.gd — 幸運星座命運魚系統面板（DAY-259）
## 業界原創「星座命運+星座祝福+星座庇護+星座標記」機制
##
## 視覺設計：
##   - 星空紫金主題（#9B59B6 紫 + #FFD700 金 + #87CEEB 天藍 + #FFF3E0 奶油）
##   - zodiac_start：紫色三次強閃光 + 頂部橫幅 + 「✨ 星座命運！」大字 + 星座符號 + 祝福/庇護指示器 + 計時條
##   - zodiac_broadcast：頂部小橫幅（全服廣播）
##   - zodiac_mark_kill：金色閃光 + 「✨ 星座標記擊破！×2.0」浮動文字
##   - zodiac_end：計時條淡出 + 結算提示
extends CanvasLayer

# 主題顏色
const COLOR_PURPLE  = Color("#9B59B6")  # 紫色（主色）
const COLOR_GOLD    = Color("#FFD700")  # 金色（祝福）
const COLOR_SKY     = Color("#87CEEB")  # 天藍（庇護）
const COLOR_CREAM   = Color("#FFF3E0")  # 奶油（副文字）
const COLOR_WHITE   = Color("#FFFFFF")  # 白色
const COLOR_STAR    = Color("#E8DAEF")  # 淡紫（星空）

# 計時條
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null

# 祝福/庇護指示器
var _boost_indicator: Label = null

func _ready() -> void:
	layer = 32  # 幸運星座命運魚面板層級（DAY-259）

## 處理幸運星座命運魚訊息
func handle_lucky_zodiac_fate(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"zodiac_start":
			_on_zodiac_start(payload)
		"zodiac_broadcast":
			_on_zodiac_broadcast(payload)
		"zodiac_mark_kill":
			_on_zodiac_mark_kill(payload)
		"zodiac_end":
			_on_zodiac_end(payload)

## zodiac_start — 星座命運啟動（個人訊息）
func _on_zodiac_start(payload: Dictionary) -> void:
	var zodiac: String = payload.get("zodiac", "未知星座")
	var zodiac_emoji: String = payload.get("zodiac_emoji", "✨")
	var zodiac_color_str: String = payload.get("zodiac_color", "#9B59B6")
	var boost_type: String = payload.get("boost_type", "shield")
	var boost_mult: float = payload.get("boost_mult", 1.5)
	var boost_duration: int = payload.get("boost_duration", 5)
	var mark_count: int = payload.get("mark_count", 3)
	var mark_mult: float = payload.get("mark_mult", 2.0)
	var mark_duration: int = payload.get("mark_duration", 15)

	var zodiac_color = Color(zodiac_color_str)
	var is_blessed = boost_type == "bless"

	# 三次強閃光（對應星座顏色）
	_flash_screen(zodiac_color, 0.65, 3)

	# 頂部橫幅
	var boost_text = "×%.1f 星座祝福 %ds！" % [boost_mult, boost_duration] if is_blessed else "×%.1f 星座庇護 %ds！" % [boost_mult, boost_duration]
	_show_banner("%s 今日星座：%s！%s%d 個目標被星座標記 ×%.1f！" % [zodiac_emoji, zodiac, boost_text, mark_count, mark_mult], zodiac_color, 4.5)

	# 中央大字
	var big_text = "%s 今日星座：%s！" % [zodiac_emoji, zodiac]
	_show_big_text(big_text, zodiac_color, 48, 3.0)

	# 副文字
	var sub_text = boost_text + " %d 個目標被星座標記 ×%.1f！" % [mark_count, mark_mult]
	_show_sub_text(sub_text, COLOR_CREAM, 3.0)

	# 祝福/庇護指示器（右上角）
	var indicator_color = COLOR_GOLD if is_blessed else COLOR_SKY
	var indicator_text = "%s ×%.1f 星座祝福！" % [zodiac_emoji, boost_mult] if is_blessed else "%s ×%.1f 星座庇護！" % [zodiac_emoji, boost_mult]
	_show_boost_indicator(indicator_text, indicator_color, float(boost_duration))

	# 右側豎向計時條（x=-212 與其他計時條錯開）
	_start_timer_bar(mark_duration, zodiac_color)

## zodiac_broadcast — 全服廣播星座命運
func _on_zodiac_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var zodiac: String = payload.get("zodiac", "未知星座")
	var zodiac_emoji: String = payload.get("zodiac_emoji", "✨")
	var zodiac_color_str: String = payload.get("zodiac_color", "#9B59B6")
	var blessed_count: int = payload.get("blessed_count", 0)
	var mark_count: int = payload.get("mark_count", 3)
	var mark_mult: float = payload.get("mark_mult", 2.0)

	var zodiac_color = Color(zodiac_color_str)
	_show_top_banner("%s %s 觸發星座命運！今日星座：%s！%d 人獲得祝福！%d 個目標被標記 ×%.1f！" % [zodiac_emoji, player_name, zodiac, blessed_count, mark_count, mark_mult], zodiac_color, 3.5)
	_start_timer_bar(15, zodiac_color)

## zodiac_mark_kill — 星座標記目標被擊破（全服廣播）
func _on_zodiac_mark_kill(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var target_name: String = payload.get("target_name", "目標")
	var mark_mult: float = payload.get("mark_mult", 2.0)
	var total_reward: int = payload.get("total_reward", 0)

	# 金色閃光
	_flash_screen(COLOR_GOLD, 0.18, 1)

	# 浮動文字
	var reward_text = ""
	if total_reward > 0:
		reward_text = " 全服+%d" % total_reward
	_show_float_text("✨ %s 擊破星座標記！%s ×%.1f%s" % [player_name, target_name, mark_mult, reward_text], COLOR_GOLD, 2.0)

## zodiac_end — 星座標記結束（全服廣播）
func _on_zodiac_end(payload: Dictionary) -> void:
	var zodiac: String = payload.get("zodiac", "未知星座")
	var zodiac_emoji: String = payload.get("zodiac_emoji", "✨")

	# 停止計時條和指示器
	_stop_timer_bar()
	_clear_boost_indicator()

	# 淡出提示
	_show_float_text("%s %s 星座命運結束！" % [zodiac_emoji, zodiac], COLOR_PURPLE, 2.0)

# ─── 祝福/庇護指示器 ──────────────────────────────────────────────────────────

func _show_boost_indicator(text: String, color: Color, duration: float) -> void:
	_clear_boost_indicator()

	var vp_size = get_viewport().size
	_boost_indicator = Label.new()
	_boost_indicator.text = text
	_boost_indicator.add_theme_font_size_override("font_size", 15)
	_boost_indicator.add_theme_color_override("font_color", color)
	_boost_indicator.position = Vector2(vp_size.x - 160, 90)
	_boost_indicator.size = Vector2(150, 30)
	_boost_indicator.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	add_child(_boost_indicator)

	# 脈衝動畫
	var tween = _boost_indicator.create_tween().set_loops()
	tween.tween_property(_boost_indicator, "modulate:a", 0.5, 0.6)
	tween.tween_property(_boost_indicator, "modulate:a", 1.0, 0.6)

	# 持續時間後自動清除
	var timer = get_tree().create_timer(duration)
	timer.timeout.connect(func():
		if is_instance_valid(_boost_indicator):
			var fade = create_tween()
			fade.tween_property(_boost_indicator, "modulate:a", 0.0, 0.5)
			fade.tween_callback(func():
				if is_instance_valid(_boost_indicator):
					_boost_indicator.queue_free()
					_boost_indicator = null
			)
	)

func _clear_boost_indicator() -> void:
	if is_instance_valid(_boost_indicator):
		_boost_indicator.queue_free()
		_boost_indicator = null

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
	lbl.add_theme_font_size_override("font_size", 15)
	lbl.add_theme_color_override("font_color", Color(0.05, 0.05, 0.05))
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
	lbl.add_theme_color_override("font_color", Color(0.05, 0.05, 0.05))
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
		randf_range(vp_size.x * 0.25, vp_size.x * 0.65),
		randf_range(vp_size.y * 0.25, vp_size.y * 0.55)
	)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 60, duration)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, duration)
	tween.tween_callback(lbl.queue_free)

func _start_timer_bar(duration_sec: int, color: Color) -> void:
	_stop_timer_bar()

	var vp_size = get_viewport().size
	var bar_h: float = vp_size.y * 0.5
	var bar_w: float = 12.0
	# x=-212 與其他計時條錯開（龍捲風-58，黑洞-72，鏡像分裂-86，量子糾纏-100，武器進化-114，隕石雨-128，龍王-142，時空裂縫-156，全服充能-170，公會戰-184，閃電風暴-198，星座命運-212）
	var bar_x: float = vp_size.x - 212.0
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.05, 0.0, 0.05, 0.7)
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
