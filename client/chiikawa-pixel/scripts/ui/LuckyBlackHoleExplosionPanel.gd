## LuckyBlackHoleExplosionPanel.gd — 幸運黑洞爆炸魚系統面板（DAY-249）
## 業界原創「黑洞吸收+能量爆炸」機制
##
## 視覺設計：
##   - 深黑宇宙主題（#2C3E50 + #1A252F + #85929E + #F39C12）
##   - blackhole_start：深黑三次強閃光 + 頂部橫幅 + 「🕳️ 黑洞生成！」大字 + 中央黑洞圓圈 + 能量計數器 + 計時條
##   - blackhole_broadcast：頂部小橫幅（全服廣播）
##   - blackhole_absorb：吸收閃光 + 「🕳️ 吸收 [N/6] 目標名」浮動文字 + 能量計數器更新
##   - blackhole_explosion：全螢幕三次強閃光 + 「🕳️ 黑洞爆炸！」大字 + 結算彈窗
##   - blackhole_end：計時條淡出 + 結束提示
extends CanvasLayer

# 主題顏色
const COLOR_BLACK_HOLE = Color("#2C3E50")  # 深黑（主題）
const COLOR_DARK       = Color("#1A252F")  # 更深黑（強調）
const COLOR_GRAY       = Color("#85929E")  # 灰色（次要）
const COLOR_GOLD       = Color("#F39C12")  # 金色（獎勵/能量）
const COLOR_ORANGE     = Color("#E67E22")  # 橙色（爆炸）
const COLOR_WHITE      = Color("#FFFFFF")  # 白色

# 計時條
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null

# 黑洞圓圈
var _black_hole_circle: Control = null

# 能量計數器
var _energy_label: Label = null
var _energy_count: int = 0
var _max_absorb: int = 6

func _ready() -> void:
	layer = 22  # 幸運黑洞爆炸魚面板層級（DAY-249）

## 處理幸運黑洞爆炸魚訊息
func handle_lucky_black_hole_explosion(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"blackhole_start":
			_on_blackhole_start(payload)
		"blackhole_broadcast":
			_on_blackhole_broadcast(payload)
		"blackhole_absorb":
			_on_blackhole_absorb(payload)
		"blackhole_explosion":
			_on_blackhole_explosion(payload)
		"blackhole_end":
			_on_blackhole_end(payload)

## blackhole_start — 黑洞生成（個人訊息）
func _on_blackhole_start(payload: Dictionary) -> void:
	var duration_sec: int = payload.get("duration_sec", 10)
	var max_absorb: int = payload.get("max_absorb", 6)
	var absorb_mult: float = payload.get("absorb_mult", 1.2)
	var center_x: float = payload.get("center_x", 500.0)
	var center_y: float = payload.get("center_y", 300.0)

	_max_absorb = max_absorb
	_energy_count = 0

	# 深黑三次強閃光
	_flash_screen(COLOR_BLACK_HOLE, 0.6, 3)

	# 頂部橫幅
	_show_banner("🕳️ 黑洞生成！吸收目標累積能量！", COLOR_BLACK_HOLE, 4.0)

	# 中央大字
	_show_big_text("🕳️ 黑洞！", COLOR_GRAY, 52, 2.5)

	# 倍率說明
	_show_sub_text("吸收 ×%.1f（個人）  爆炸能量×目標數×0.8（全服）" % absorb_mult, COLOR_GOLD, 2.5)

	# 中央黑洞圓圈
	_show_black_hole_circle(center_x, center_y)

	# 右側能量計數器
	_show_energy_counter()

	# 右側豎向計時條
	_start_timer_bar(duration_sec)

## blackhole_broadcast — 全服廣播黑洞生成
func _on_blackhole_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	_show_top_banner("🕳️ %s 觸發黑洞！" % player_name, COLOR_BLACK_HOLE, 2.5)

## blackhole_absorb — 吸收目標（全服廣播）
func _on_blackhole_absorb(payload: Dictionary) -> void:
	var absorb_count: int = payload.get("absorb_count", 1)
	var max_absorb: int = payload.get("max_absorb", 6)
	var target_name: String = payload.get("target_name", "目標")
	var reward: int = payload.get("reward", 0)

	_energy_count = absorb_count

	# 深黑閃光
	_flash_screen(COLOR_BLACK_HOLE, 0.3, 1)

	# 吸收浮動文字
	_show_float_text("🕳️ [%d/%d] 吸收 %s！" % [absorb_count, max_absorb, target_name], COLOR_GRAY, 1.8)

	# 個人獎勵提示
	if reward > 0:
		_show_float_text("+%d 籌碼" % reward, COLOR_GOLD, 1.5)

	# 更新能量計數器
	_update_energy_counter(absorb_count, max_absorb)

## blackhole_explosion — 黑洞爆炸（全服廣播）
func _on_blackhole_explosion(payload: Dictionary) -> void:
	var energy: int = payload.get("energy", 0)
	var target_count: int = payload.get("target_count", 0)
	var blast_mult: float = payload.get("blast_mult", 1.0)
	var total_reward: int = payload.get("total_reward", 0)

	# 清除黑洞圓圈（爆炸擴散）
	if is_instance_valid(_black_hole_circle):
		var tween = create_tween()
		tween.tween_property(_black_hole_circle, "scale", Vector2(4.0, 4.0), 0.4)
		tween.parallel().tween_property(_black_hole_circle, "modulate:a", 0.0, 0.4)
		tween.tween_callback(_black_hole_circle.queue_free)
		_black_hole_circle = null

	# 清除能量計數器
	if is_instance_valid(_energy_label):
		var tween2 = create_tween()
		tween2.tween_property(_energy_label, "modulate:a", 0.0, 0.3)
		tween2.tween_callback(_energy_label.queue_free)
		_energy_label = null

	# 停止計時條
	_stop_timer_bar()

	# 全螢幕三次強閃光（橙色爆炸感）
	_flash_screen(COLOR_ORANGE, 0.7, 3)

	# 中央大字
	_show_big_text("🕳️ 黑洞爆炸！", COLOR_ORANGE, 52, 2.5)

	# 爆炸說明
	_show_sub_text("能量 %d × 目標 %d × 0.8 = ×%.1f！" % [energy, target_count, blast_mult], COLOR_GOLD, 2.5)

	# 結算彈窗
	if total_reward > 0:
		_show_result_popup(energy, blast_mult, total_reward)

## blackhole_end — 黑洞結束（個人訊息）
func _on_blackhole_end(payload: Dictionary) -> void:
	var energy: int = payload.get("energy", 0)
	if energy == 0:
		_show_float_text("🕳️ 黑洞結束（未吸收目標）", COLOR_GRAY, 1.5)

# ─── 黑洞圓圈 ────────────────────────────────────────────────────────────────

func _show_black_hole_circle(center_x: float, center_y: float) -> void:
	var vp_size = get_viewport().size
	var scale_x = vp_size.x / 1000.0
	var scale_y = vp_size.y / 600.0
	var screen_x = center_x * scale_x
	var screen_y = center_y * scale_y

	_black_hole_circle = Control.new()
	_black_hole_circle.position = Vector2(screen_x, screen_y)
	add_child(_black_hole_circle)

	# 多層同心圓（模擬黑洞）
	for i in range(5):
		var radius = 20.0 + i * 25.0
		var circle = ColorRect.new()
		circle.size = Vector2(radius * 2, radius * 2)
		circle.position = Vector2(-radius, -radius)
		# 越外層越透明，越深色
		var alpha = 0.5 - i * 0.08
		circle.color = Color(COLOR_BLACK_HOLE.r, COLOR_BLACK_HOLE.g, COLOR_BLACK_HOLE.b, alpha)
		_black_hole_circle.add_child(circle)

	# 中心黑點
	var center_dot = ColorRect.new()
	center_dot.size = Vector2(12, 12)
	center_dot.position = Vector2(-6, -6)
	center_dot.color = Color(0.0, 0.0, 0.0, 1.0)
	_black_hole_circle.add_child(center_dot)

	# 持續旋轉動畫（逆時針，模擬黑洞吸引）
	var tween = _black_hole_circle.create_tween().set_loops()
	tween.tween_property(_black_hole_circle, "rotation", -TAU, 3.0).set_ease(Tween.EASE_IN_OUT)

	# 脈衝縮放（黑洞呼吸感）
	var pulse_tween = _black_hole_circle.create_tween().set_loops()
	pulse_tween.tween_property(_black_hole_circle, "scale", Vector2(1.2, 1.2), 1.0)
	pulse_tween.tween_property(_black_hole_circle, "scale", Vector2(0.9, 0.9), 1.0)

# ─── 能量計數器 ──────────────────────────────────────────────────────────────

func _show_energy_counter() -> void:
	var vp_size = get_viewport().size
	_energy_label = Label.new()
	_energy_label.text = "⚡ 0/%d" % _max_absorb
	_energy_label.add_theme_font_size_override("font_size", 20)
	_energy_label.add_theme_color_override("font_color", COLOR_GOLD)
	_energy_label.position = Vector2(vp_size.x - 100, vp_size.y * 0.15)
	add_child(_energy_label)

func _update_energy_counter(count: int, max_count: int) -> void:
	if not is_instance_valid(_energy_label):
		return
	_energy_label.text = "⚡ %d/%d" % [count, max_count]
	# 閃爍動畫
	var tween = create_tween()
	tween.tween_property(_energy_label, "modulate", Color(COLOR_ORANGE, 1.0), 0.1)
	tween.tween_property(_energy_label, "modulate", Color(1, 1, 1, 1), 0.1)

# ─── 結算彈窗 ────────────────────────────────────────────────────────────────

func _show_result_popup(energy: int, blast_mult: float, total_reward: int) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(320, 160)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 80)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.05, 0.08, 0.92)
	style.border_color = COLOR_ORANGE
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
	title_lbl.text = "🕳️ 黑洞爆炸結算"
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", COLOR_ORANGE)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var energy_lbl = Label.new()
	energy_lbl.text = "能量：%d 個目標" % energy
	energy_lbl.add_theme_font_size_override("font_size", 16)
	energy_lbl.add_theme_color_override("font_color", COLOR_GRAY)
	energy_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(energy_lbl)

	var mult_lbl = Label.new()
	mult_lbl.text = "爆炸倍率：×%.1f" % blast_mult
	mult_lbl.add_theme_font_size_override("font_size", 18)
	mult_lbl.add_theme_color_override("font_color", COLOR_ORANGE)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "全服獎勵：+%d 籌碼" % total_reward
	reward_lbl.add_theme_font_size_override("font_size", 18)
	reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 340.0, 0.4).set_ease(Tween.EASE_OUT)
	tween.tween_interval(4.0)
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
	var bar_x: float = vp_size.x - 72.0  # 與其他計時條錯開（龍捲風在 -58，黑洞在 -72）
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.05, 0.05, 0.08, 0.7)
	_timer_bar_bg.size = Vector2(bar_w, bar_h)
	_timer_bar_bg.position = Vector2(bar_x, bar_y)
	add_child(_timer_bar_bg)

	_timer_bar = ColorRect.new()
	_timer_bar.color = COLOR_BLACK_HOLE
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
