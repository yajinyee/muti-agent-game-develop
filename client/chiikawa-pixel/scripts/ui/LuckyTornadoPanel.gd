## LuckyTornadoPanel.gd — 幸運龍捲風魚系統面板（DAY-248）
## 業界原創「龍捲風吸引+螺旋爆發」機制
##
## 視覺設計：
##   - 青綠龍捲風主題（#1ABC9C + #16A085 + #A3E4D7 + #27AE60）
##   - tornado_start：青綠三次強閃光 + 頂部橫幅 + 「🌪️ 龍捲風！」大字 + 中央龍捲風圓圈 + 計時條
##   - tornado_broadcast：頂部小橫幅（全服廣播）
##   - tornado_spiral：青綠閃光 + 「🌪️ 螺旋吸引！N 個目標」提示（每 2 秒）
##   - tornado_blast：全螢幕三次強閃光 + 「🌪️ 龍捲風爆發！」大字 + 結算彈窗
##   - tornado_end：計時條淡出 + 結束提示
extends CanvasLayer

# 主題顏色
const COLOR_TORNADO  = Color("#1ABC9C")  # 青綠（主題）
const COLOR_DARK     = Color("#16A085")  # 深青綠（強調）
const COLOR_LIGHT    = Color("#A3E4D7")  # 淺青綠（背景）
const COLOR_GREEN    = Color("#27AE60")  # 綠色（爆發）
const COLOR_GOLD     = Color("#F39C12")  # 金色（獎勵）
const COLOR_WHITE    = Color("#FFFFFF")  # 白色

# 計時條
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null

# 龍捲風中央圓圈
var _tornado_circle: Control = null

func _ready() -> void:
	layer = 21  # 幸運龍捲風魚面板層級（DAY-248）

## 處理幸運龍捲風魚訊息
func handle_lucky_tornado(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"tornado_start":
			_on_tornado_start(payload)
		"tornado_broadcast":
			_on_tornado_broadcast(payload)
		"tornado_spiral":
			_on_tornado_spiral(payload)
		"tornado_blast":
			_on_tornado_blast(payload)
		"tornado_end":
			_on_tornado_end(payload)

## tornado_start — 龍捲風啟動（個人訊息）
func _on_tornado_start(payload: Dictionary) -> void:
	var duration_sec: int = payload.get("duration_sec", 12)
	var kill_mult: float = payload.get("kill_mult", 2.2)
	var blast_mult: float = payload.get("blast_mult", 1.5)
	var center_x: float = payload.get("center_x", 500.0)
	var center_y: float = payload.get("center_y", 300.0)

	# 青綠三次強閃光
	_flash_screen(COLOR_TORNADO, 0.5, 3)

	# 頂部橫幅
	_show_banner("🌪️ 龍捲風！螺旋吸引所有目標！", COLOR_TORNADO, 4.0)

	# 中央大字
	_show_big_text("🌪️ 龍捲風！", COLOR_TORNADO, 52, 2.5)

	# 倍率說明
	_show_sub_text("擊破 ×%.1f  爆發 ×%.1f（全服共享）" % [kill_mult, blast_mult], COLOR_GOLD, 2.0)

	# 中央龍捲風圓圈（旋轉動畫）
	_show_tornado_circle(center_x, center_y)

	# 右側豎向計時條
	_start_timer_bar(duration_sec)

## tornado_broadcast — 全服廣播龍捲風啟動
func _on_tornado_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	_show_top_banner("🌪️ %s 觸發龍捲風！" % player_name, COLOR_TORNADO, 2.5)

## tornado_spiral — 螺旋吸引（每 2 秒）
func _on_tornado_spiral(payload: Dictionary) -> void:
	var pull_count: int = payload.get("pull_count", 1)
	var moved_count: int = payload.get("moved_count", 0)

	# 青綠閃光
	_flash_screen(COLOR_TORNADO, 0.25, 1)

	# 提示文字
	_show_float_text("🌪️ 螺旋吸引 #%d！%d 個目標" % [pull_count, moved_count], COLOR_LIGHT, 1.5)

	# 龍捲風圓圈旋轉加速
	if is_instance_valid(_tornado_circle):
		var tween = create_tween()
		tween.tween_property(_tornado_circle, "rotation", _tornado_circle.rotation + PI * 0.5, 0.3)

## tornado_blast — 龍捲風爆發
func _on_tornado_blast(payload: Dictionary) -> void:
	var blast_count: int = payload.get("blast_count", 0)
	var total_reward: int = payload.get("total_reward", 0)
	var blast_mult: float = payload.get("blast_mult", 1.5)

	# 清除龍捲風圓圈
	if is_instance_valid(_tornado_circle):
		var tween = create_tween()
		tween.tween_property(_tornado_circle, "scale", Vector2(3.0, 3.0), 0.3)
		tween.parallel().tween_property(_tornado_circle, "modulate:a", 0.0, 0.3)
		tween.tween_callback(_tornado_circle.queue_free)
		_tornado_circle = null

	# 停止計時條
	_stop_timer_bar()

	# 全螢幕三次強閃光
	_flash_screen(COLOR_GREEN, 0.6, 3)

	# 中央大字
	_show_big_text("🌪️ 龍捲風爆發！", COLOR_GREEN, 52, 2.5)

	# 爆發說明
	_show_sub_text("捲走 %d 個目標！×%.1f 全服共享！" % [blast_count, blast_mult], COLOR_GOLD, 2.0)

	# 結算彈窗
	if blast_count > 0:
		_show_result_popup(blast_count, total_reward)

## tornado_end — 龍捲風結束（個人訊息）
func _on_tornado_end(payload: Dictionary) -> void:
	var blast_count: int = payload.get("blast_count", 0)
	var total_reward: int = payload.get("total_reward", 0)
	# 個人結算已在 tornado_blast 顯示，這裡只做清理
	if blast_count == 0:
		_show_float_text("🌪️ 龍捲風結束", COLOR_LIGHT, 1.5)

# ─── 龍捲風圓圈 ──────────────────────────────────────────────────────────────

func _show_tornado_circle(center_x: float, center_y: float) -> void:
	var vp_size = get_viewport().size
	# 將遊戲座標轉換為螢幕座標（假設遊戲畫面 1000x600 → 螢幕）
	var scale_x = vp_size.x / 1000.0
	var scale_y = vp_size.y / 600.0
	var screen_x = center_x * scale_x
	var screen_y = center_y * scale_y

	_tornado_circle = Control.new()
	_tornado_circle.position = Vector2(screen_x, screen_y)
	add_child(_tornado_circle)

	# 多層同心圓（模擬龍捲風）
	for i in range(4):
		var radius = 40.0 + i * 30.0
		var circle = ColorRect.new()
		circle.size = Vector2(radius * 2, radius * 2)
		circle.position = Vector2(-radius, -radius)
		circle.color = Color(COLOR_TORNADO.r, COLOR_TORNADO.g, COLOR_TORNADO.b, 0.15 - i * 0.03)
		_tornado_circle.add_child(circle)

	# 中心點
	var center_dot = ColorRect.new()
	center_dot.size = Vector2(16, 16)
	center_dot.position = Vector2(-8, -8)
	center_dot.color = COLOR_TORNADO
	_tornado_circle.add_child(center_dot)

	# 持續旋轉動畫
	var tween = _tornado_circle.create_tween().set_loops()
	tween.tween_property(_tornado_circle, "rotation", TAU, 2.0).set_ease(Tween.EASE_IN_OUT)

	# 脈衝縮放
	var pulse_tween = _tornado_circle.create_tween().set_loops()
	pulse_tween.tween_property(_tornado_circle, "scale", Vector2(1.15, 1.15), 0.8)
	pulse_tween.tween_property(_tornado_circle, "scale", Vector2(1.0, 1.0), 0.8)

# ─── 結算彈窗 ────────────────────────────────────────────────────────────────

func _show_result_popup(blast_count: int, total_reward: int) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(300, 140)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 70)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.0, 0.08, 0.06, 0.92)
	style.border_color = COLOR_TORNADO
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
	title_lbl.text = "🌪️ 龍捲風爆發結算"
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", COLOR_TORNADO)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var count_lbl = Label.new()
	count_lbl.text = "捲走目標：%d 個" % blast_count
	count_lbl.add_theme_font_size_override("font_size", 16)
	count_lbl.add_theme_color_override("font_color", COLOR_LIGHT)
	count_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(count_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "全服獎勵：+%d 籌碼" % total_reward
	reward_lbl.add_theme_font_size_override("font_size", 18)
	reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 320.0, 0.4).set_ease(Tween.EASE_OUT)
	tween.tween_interval(3.5)
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
	lbl.add_theme_font_size_override("font_size", 18)
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
	lbl.add_theme_font_size_override("font_size", 14)
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
	lbl.add_theme_font_size_override("font_size", 16)
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
	lbl.add_theme_font_size_override("font_size", 18)
	lbl.add_theme_color_override("font_color", color)
	lbl.position = Vector2(
		randf_range(vp_size.x * 0.3, vp_size.x * 0.7),
		randf_range(vp_size.y * 0.3, vp_size.y * 0.6)
	)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 60, duration)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, duration)
	tween.tween_callback(lbl.queue_free)

func _start_timer_bar(duration_sec: int) -> void:
	var vp_size = get_viewport().size
	var bar_h: float = vp_size.y * 0.5
	var bar_w: float = 12.0
	var bar_x: float = vp_size.x - 58.0  # 與其他計時條錯開
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.0, 0.08, 0.06, 0.7)
	_timer_bar_bg.size = Vector2(bar_w, bar_h)
	_timer_bar_bg.position = Vector2(bar_x, bar_y)
	add_child(_timer_bar_bg)

	_timer_bar = ColorRect.new()
	_timer_bar.color = COLOR_TORNADO
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
