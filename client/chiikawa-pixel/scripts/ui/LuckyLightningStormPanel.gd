## LuckyLightningStormPanel.gd — 幸運閃電風暴魚系統面板（DAY-258）
## 業界原創「閃電風暴+連鎖跳躍+超級閃電+全場電擊」機制
##
## 視覺設計：
##   - 金黃閃電主題（#FFD700 金 + #FFFFFF 白 + #87CEEB 天藍 + #FFF3E0 奶油）
##   - storm_start：金色三次強閃光 + 頂部橫幅 + 「⚡ 閃電風暴！」大字 + 跳躍計數器 + 計時條
##   - storm_broadcast：頂部小橫幅（全服廣播）
##   - storm_jump：金色閃光 + 「⚡ 第N輪 跳躍M個！×1.3」浮動文字 + 跳躍計數器更新
##   - super_lightning：全螢幕三次強閃光 + 「⚡ 超級閃電！×3.0」大字 + 結算彈窗
##   - storm_blast：天藍閃光 + 「⚡ 閃電爆炸！HP-40%」大字
extends CanvasLayer

# 主題顏色
const COLOR_GOLD    = Color("#FFD700")  # 金色（主色）
const COLOR_WHITE   = Color("#FFFFFF")  # 白色（超級閃電）
const COLOR_SKY     = Color("#87CEEB")  # 天藍（爆炸）
const COLOR_CREAM   = Color("#FFF3E0")  # 奶油（副文字）
const COLOR_ORANGE  = Color("#FF8C00")  # 橙色（高能量）
const COLOR_YELLOW  = Color("#FFFF00")  # 亮黃（閃電）

# 計時條
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null

# 跳躍計數器
var _jump_counter: Label = null
var _total_jumps: int = 0
var _super_thresh: int = 5

func _ready() -> void:
	layer = 31  # 幸運閃電風暴魚面板層級（DAY-258）

## 處理幸運閃電風暴魚訊息
func handle_lucky_lightning_storm(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"storm_start":
			_on_storm_start(payload)
		"storm_broadcast":
			_on_storm_broadcast(payload)
		"storm_jump":
			_on_storm_jump(payload)
		"super_lightning":
			_on_super_lightning(payload)
		"storm_blast":
			_on_storm_blast(payload)

## storm_start — 閃電風暴啟動（個人訊息）
func _on_storm_start(payload: Dictionary) -> void:
	var duration_sec: int = payload.get("duration_sec", 12)
	var jump_mult: float = payload.get("jump_mult", 1.3)
	var super_mult: float = payload.get("super_mult", 3.0)
	var super_thresh: int = payload.get("super_thresh", 5)
	_super_thresh = super_thresh
	_total_jumps = 0

	# 金色三次強閃光
	_flash_screen(COLOR_GOLD, 0.65, 3)

	# 頂部橫幅
	_show_banner("⚡ 閃電風暴！每 1.5 秒連鎖跳躍 ×%.1f！累計 %d 跳→超級閃電 ×%.1f！" % [jump_mult, super_thresh, super_mult], COLOR_GOLD, 4.0)

	# 中央大字
	_show_big_text("⚡ 閃電風暴！", COLOR_GOLD, 52, 2.5)
	_show_sub_text("連鎖跳躍 ×%.1f！累計 %d 跳→超級閃電 ×%.1f！" % [jump_mult, super_thresh, super_mult], COLOR_CREAM, 2.5)

	# 跳躍計數器（右上角）
	_show_jump_counter(0, super_thresh)

	# 右側豎向計時條（x=-198 與其他計時條錯開）
	_start_timer_bar(duration_sec, COLOR_GOLD)

## storm_broadcast — 全服廣播閃電風暴
func _on_storm_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var jump_mult: float = payload.get("jump_mult", 1.3)
	var super_mult: float = payload.get("super_mult", 3.0)
	var duration_sec: int = payload.get("duration_sec", 12)
	_show_top_banner("⚡ %s 觸發閃電風暴！連鎖跳躍 ×%.1f！超級閃電 ×%.1f！" % [player_name, jump_mult, super_mult], COLOR_GOLD, 3.0)
	_show_jump_counter(0, _super_thresh)
	_start_timer_bar(duration_sec, COLOR_GOLD)

## storm_jump — 閃電跳躍（全服廣播）
func _on_storm_jump(payload: Dictionary) -> void:
	var round_num: int = payload.get("round", 0)
	var jump_count: int = payload.get("jump_count", 0)
	var total_jumps: int = payload.get("total_jumps", 0)
	var jump_mult: float = payload.get("jump_mult", 1.3)
	var total_reward: int = payload.get("total_reward", 0)
	_total_jumps = total_jumps

	# 金色閃光
	_flash_screen(COLOR_GOLD, 0.12, 1)

	# 更新跳躍計數器
	_update_jump_counter(total_jumps)

	# 浮動文字
	var reward_text = ""
	if total_reward > 0:
		reward_text = " +%d" % total_reward
	_show_float_text("⚡ 第%d輪 跳躍%d個！×%.1f%s" % [round_num, jump_count, jump_mult, reward_text], COLOR_YELLOW, 1.8)

## super_lightning — 超級閃電（全服廣播）
func _on_super_lightning(payload: Dictionary) -> void:
	var super_mult: float = payload.get("super_mult", 3.0)
	var total_reward: int = payload.get("total_reward", 0)
	var total_jumps: int = payload.get("total_jumps", 5)
	var player_name: String = payload.get("player_name", "某玩家")

	# 全螢幕三次強閃光（白色，超級閃電感）
	_flash_screen(COLOR_WHITE, 0.9, 3)

	# 大字
	_show_big_text("⚡ 超級閃電！", COLOR_WHITE, 56, 2.5)
	_show_sub_text("累計 %d 跳！全服獎勵 +%d！×%.1f 大獎！" % [total_jumps, total_reward, super_mult], COLOR_GOLD, 2.5)

	# 結算彈窗
	_show_super_popup(player_name, total_jumps, super_mult, total_reward)

## storm_blast — 閃電爆炸（全服廣播）
func _on_storm_blast(payload: Dictionary) -> void:
	var affected_count: int = payload.get("affected_count", 0)
	var total_jumps: int = payload.get("total_jumps", 0)

	# 停止計時條和計數器
	_stop_timer_bar()
	_clear_jump_counter()

	# 天藍閃光（爆炸感）
	_flash_screen(COLOR_SKY, 0.5, 2)

	# 大字
	_show_big_text("⚡ 閃電爆炸！", COLOR_SKY, 48, 2.5)
	_show_sub_text("%d 個目標 HP-40%%！累計跳躍 %d 次！" % [affected_count, total_jumps], COLOR_CREAM, 2.5)

# ─── 跳躍計數器 ───────────────────────────────────────────────────────────────

func _show_jump_counter(current: int, target: int) -> void:
	_clear_jump_counter()

	var vp_size = get_viewport().size
	_jump_counter = Label.new()
	_jump_counter.text = "⚡ 跳躍 %d/%d" % [current, target]
	_jump_counter.add_theme_font_size_override("font_size", 16)
	_jump_counter.add_theme_color_override("font_color", COLOR_GOLD)
	_jump_counter.position = Vector2(vp_size.x - 130, 55)
	_jump_counter.size = Vector2(120, 30)
	_jump_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	add_child(_jump_counter)

	# 脈衝動畫
	var tween = _jump_counter.create_tween().set_loops()
	tween.tween_property(_jump_counter, "modulate:a", 0.6, 0.5)
	tween.tween_property(_jump_counter, "modulate:a", 1.0, 0.5)

func _update_jump_counter(total: int) -> void:
	if not is_instance_valid(_jump_counter):
		_show_jump_counter(total, _super_thresh)
		return
	_jump_counter.text = "⚡ 跳躍 %d/%d" % [total, _super_thresh]
	# 接近超級閃電時變白色
	if total >= _super_thresh - 1:
		_jump_counter.add_theme_color_override("font_color", COLOR_WHITE)
		var tween = _jump_counter.create_tween()
		tween.tween_property(_jump_counter, "scale", Vector2(1.3, 1.3), 0.1)
		tween.tween_property(_jump_counter, "scale", Vector2(1.0, 1.0), 0.1)

func _clear_jump_counter() -> void:
	if is_instance_valid(_jump_counter):
		_jump_counter.queue_free()
		_jump_counter = null

# ─── 超級閃電結算彈窗 ─────────────────────────────────────────────────────────

func _show_super_popup(player_name: String, total_jumps: int, super_mult: float, total_reward: int) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(340, 180)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 90)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.05, 0.05, 0.93)
	style.border_color = COLOR_WHITE
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
	title_lbl.text = "⚡ 超級閃電！"
	title_lbl.add_theme_font_size_override("font_size", 22)
	title_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var trigger_lbl = Label.new()
	trigger_lbl.text = "觸發者：%s" % player_name
	trigger_lbl.add_theme_font_size_override("font_size", 13)
	trigger_lbl.add_theme_color_override("font_color", COLOR_CREAM)
	trigger_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(trigger_lbl)

	var jumps_lbl = Label.new()
	jumps_lbl.text = "累計跳躍 %d 次！" % total_jumps
	jumps_lbl.add_theme_font_size_override("font_size", 16)
	jumps_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	jumps_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(jumps_lbl)

	var mult_lbl = Label.new()
	mult_lbl.text = "全服 ×%.1f 大獎！" % super_mult
	mult_lbl.add_theme_font_size_override("font_size", 20)
	mult_lbl.add_theme_color_override("font_color", COLOR_YELLOW)
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
	tween.tween_property(popup, "position:x", vp_size.x - 360.0, 0.4).set_ease(Tween.EASE_OUT)
	tween.tween_interval(5.0)
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
	# x=-198 與其他計時條錯開（龍捲風-58，黑洞-72，鏡像分裂-86，量子糾纏-100，武器進化-114，隕石雨-128，龍王-142，時空裂縫-156，全服充能-170，公會戰-184，閃電風暴-198）
	var bar_x: float = vp_size.x - 198.0
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.05, 0.05, 0.0, 0.7)
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
