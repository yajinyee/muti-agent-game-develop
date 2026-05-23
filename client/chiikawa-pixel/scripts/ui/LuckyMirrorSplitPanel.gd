## LuckyMirrorSplitPanel.gd — 幸運鏡像分裂魚系統面板（DAY-250）
## 業界原創「鏡像分裂+雙重目標」機制
##
## 視覺設計：
##   - 紫色鏡像主題（#8E44AD + #6C3483 + #D7BDE2 + #F5EEF8）
##   - mirror_split_start：紫色三次強閃光 + 頂部橫幅 + 「🪞 鏡像分裂！」大字 + 鏡像線 + 計時條
##   - mirror_split_broadcast：頂部小橫幅（全服廣播）
##   - mirror_split_kill：紫色閃光 + 「🪞 擊破鏡像！×0.6」浮動文字
##   - mirror_split_fade：灰色閃光 + 「🪞 鏡像消融！全服共享」提示 + 結算彈窗
##   - mirror_split_end：計時條淡出
extends CanvasLayer

# 主題顏色
const COLOR_MIRROR  = Color("#8E44AD")  # 紫色（主題）
const COLOR_DARK    = Color("#6C3483")  # 深紫（強調）
const COLOR_LIGHT   = Color("#D7BDE2")  # 淺紫（背景）
const COLOR_FADE    = Color("#7F8C8D")  # 灰色（消融）
const COLOR_GOLD    = Color("#F39C12")  # 金色（獎勵）
const COLOR_WHITE   = Color("#FFFFFF")  # 白色

# 計時條
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null

# 鏡像線
var _mirror_line: ColorRect = null

func _ready() -> void:
	layer = 23  # 幸運鏡像分裂魚面板層級（DAY-250）

## 處理幸運鏡像分裂魚訊息
func handle_lucky_mirror_split(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"mirror_split_start":
			_on_mirror_split_start(payload)
		"mirror_split_broadcast":
			_on_mirror_split_broadcast(payload)
		"mirror_split_kill":
			_on_mirror_split_kill(payload)
		"mirror_split_fade":
			_on_mirror_split_fade(payload)
		"mirror_split_end":
			_on_mirror_split_end(payload)

## mirror_split_start — 鏡像分裂啟動（個人訊息）
func _on_mirror_split_start(payload: Dictionary) -> void:
	var split_count: int = payload.get("split_count", 4)
	var duration_sec: int = payload.get("duration_sec", 15)
	var kill_mult: float = payload.get("kill_mult", 0.6)
	var fade_mult: float = payload.get("fade_mult", 0.3)

	# 紫色三次強閃光
	_flash_screen(COLOR_MIRROR, 0.55, 3)

	# 頂部橫幅
	_show_banner("🪞 鏡像分裂！%d 個目標生成鏡像副本！" % split_count, COLOR_MIRROR, 4.0)

	# 中央大字
	_show_big_text("🪞 鏡像分裂！", COLOR_MIRROR, 52, 2.5)

	# 倍率說明
	_show_sub_text("擊破副本 ×%.1f（個人）  消融 ×%.1f（全服）" % [kill_mult, fade_mult], COLOR_GOLD, 2.5)

	# 場景中央鏡像線
	_show_mirror_line()

	# 右側豎向計時條
	_start_timer_bar(duration_sec)

## mirror_split_broadcast — 全服廣播分裂
func _on_mirror_split_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var split_count: int = payload.get("split_count", 4)
	_show_top_banner("🪞 %s 觸發鏡像分裂！%d 個副本！" % [player_name, split_count], COLOR_MIRROR, 2.5)

## mirror_split_kill — 鏡像副本被擊破
func _on_mirror_split_kill(payload: Dictionary) -> void:
	var reward: int = payload.get("reward", 0)
	var kill_mult: float = payload.get("kill_mult", 0.6)

	# 紫色閃光
	_flash_screen(COLOR_MIRROR, 0.3, 1)

	# 浮動文字
	_show_float_text("🪞 擊破鏡像！×%.1f  +%d" % [kill_mult, reward], COLOR_LIGHT, 1.8)

## mirror_split_fade — 鏡像消融（全服廣播）
func _on_mirror_split_fade(payload: Dictionary) -> void:
	var fade_count: int = payload.get("fade_count", 0)
	var total_reward: int = payload.get("total_reward", 0)
	var fade_mult: float = payload.get("fade_mult", 0.3)

	# 清除鏡像線
	if is_instance_valid(_mirror_line):
		var tween = create_tween()
		tween.tween_property(_mirror_line, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_mirror_line.queue_free)
		_mirror_line = null

	# 停止計時條
	_stop_timer_bar()

	# 灰色閃光（消融感）
	_flash_screen(COLOR_FADE, 0.4, 2)

	# 消融提示
	_show_big_text("🪞 鏡像消融！", COLOR_FADE, 44, 2.0)
	_show_sub_text("%d 個副本消融，×%.1f 全服共享！" % [fade_count, fade_mult], COLOR_GOLD, 2.0)

	# 結算彈窗
	if total_reward > 0:
		_show_result_popup(fade_count, total_reward)

## mirror_split_end — 分裂結束（個人訊息）
func _on_mirror_split_end(payload: Dictionary) -> void:
	var fade_count: int = payload.get("fade_count", 0)
	if fade_count == 0:
		_show_float_text("🪞 鏡像分裂結束（全部擊破）", COLOR_LIGHT, 1.5)

# ─── 鏡像線 ──────────────────────────────────────────────────────────────────

func _show_mirror_line() -> void:
	var vp_size = get_viewport().size
	# 場景中央垂直線（X=500 → 螢幕中央）
	_mirror_line = ColorRect.new()
	_mirror_line.color = Color(COLOR_MIRROR.r, COLOR_MIRROR.g, COLOR_MIRROR.b, 0.5)
	_mirror_line.size = Vector2(3, vp_size.y)
	_mirror_line.position = Vector2(vp_size.x / 2.0 - 1.5, 0)
	add_child(_mirror_line)

	# 閃爍動畫
	var tween = _mirror_line.create_tween().set_loops()
	tween.tween_property(_mirror_line, "modulate:a", 0.3, 0.6)
	tween.tween_property(_mirror_line, "modulate:a", 1.0, 0.6)

# ─── 結算彈窗 ────────────────────────────────────────────────────────────────

func _show_result_popup(fade_count: int, total_reward: int) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(300, 140)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 70)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.06, 0.02, 0.1, 0.92)
	style.border_color = COLOR_MIRROR
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
	title_lbl.text = "🪞 鏡像消融結算"
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", COLOR_MIRROR)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var count_lbl = Label.new()
	count_lbl.text = "消融副本：%d 個" % fade_count
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
	lbl.add_theme_font_size_override("font_size", 14)
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
	var bar_x: float = vp_size.x - 86.0  # 與其他計時條錯開（龍捲風-58，黑洞-72，鏡像分裂-86）
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.06, 0.02, 0.1, 0.7)
	_timer_bar_bg.size = Vector2(bar_w, bar_h)
	_timer_bar_bg.position = Vector2(bar_x, bar_y)
	add_child(_timer_bar_bg)

	_timer_bar = ColorRect.new()
	_timer_bar.color = COLOR_MIRROR
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
