## LuckyServerChargePanel.gd — 幸運全服充能魚系統面板（DAY-256）
## 業界原創「全服共同充能→全服大爆發」機制
##
## 視覺設計：
##   - 橙金充能主題（#FF8C00 + #FFD700 + #FF4500 + #FFF3E0）
##   - charge_start：橙色三次強閃光 + 頂部橫幅 + 「⚡ 全服充能！」大字 + 充能進度條 + 計時條
##   - charge_broadcast：頂部小橫幅（全服廣播）
##   - charge_progress：橙色閃光 + 充能進度條更新 + 「⚡ N/20」浮動文字
##   - charge_burst：全螢幕三次強閃光 + 「⚡ 全服大爆發！×2.0」大字 + 結算彈窗
##   - charge_fail：灰色閃光 + 「⚡ 充能失敗！安慰獎」提示
extends CanvasLayer

# 主題顏色
const COLOR_CHARGE_ORANGE = Color("#FF8C00")  # 充能橙（主色）
const COLOR_GOLD          = Color("#FFD700")  # 金色（大爆發）
const COLOR_FIRE          = Color("#FF4500")  # 火焰橙（爆發）
const COLOR_CREAM         = Color("#FFF3E0")  # 奶油（副文字）
const COLOR_GRAY          = Color("#808080")  # 灰色（失敗）
const COLOR_WHITE         = Color("#FFFFFF")  # 白色

# 計時條
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null

# 充能進度條
var _charge_bar: ColorRect = null
var _charge_bar_bg: ColorRect = null
var _charge_label: Label = null
var _charge_target: int = 20

func _ready() -> void:
	layer = 29  # 幸運全服充能魚面板層級（DAY-256）

## 處理幸運全服充能魚訊息
func handle_lucky_server_charge(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"charge_start":
			_on_charge_start(payload)
		"charge_broadcast":
			_on_charge_broadcast(payload)
		"charge_progress":
			_on_charge_progress(payload)
		"charge_burst":
			_on_charge_burst(payload)
		"charge_fail":
			_on_charge_fail(payload)

## charge_start — 充能啟動（個人訊息）
func _on_charge_start(payload: Dictionary) -> void:
	var duration_sec: int = payload.get("duration_sec", 30)
	var charge_target: int = payload.get("charge_target", 20)
	var burst_mult: float = payload.get("burst_mult", 2.0)
	_charge_target = charge_target

	# 橙色三次強閃光
	_flash_screen(COLOR_CHARGE_ORANGE, 0.6, 3)

	# 頂部橫幅
	_show_banner("⚡ 全服充能！全服一起打魚累積 %d 次→大爆發！×%.1f 全服共享！" % [charge_target, burst_mult], COLOR_CHARGE_ORANGE, 4.0)

	# 中央大字
	_show_big_text("⚡ 全服充能！", COLOR_CHARGE_ORANGE, 52, 2.5)
	_show_sub_text("全服一起打魚！累積 %d 次→×%.1f 全場大爆發！" % [charge_target, burst_mult], COLOR_CREAM, 2.5)

	# 充能進度條（底部）
	_show_charge_bar(0, charge_target)

	# 右側豎向計時條（x=-170 與其他計時條錯開）
	_start_timer_bar(duration_sec, COLOR_CHARGE_ORANGE)

## charge_broadcast — 全服廣播充能
func _on_charge_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var charge_target: int = payload.get("charge_target", 20)
	var burst_mult: float = payload.get("burst_mult", 2.0)
	_charge_target = charge_target
	_show_top_banner("⚡ %s 開啟全服充能！全服一起打魚累積 %d 次→×%.1f 大爆發！" % [player_name, charge_target, burst_mult], COLOR_CHARGE_ORANGE, 3.0)
	# 廣播時也顯示進度條（從 0 開始）
	_show_charge_bar(0, charge_target)
	# 計時條
	var duration_sec: int = payload.get("duration_sec", 30)
	_start_timer_bar(duration_sec, COLOR_CHARGE_ORANGE)

## charge_progress — 充能進度更新（全服廣播）
func _on_charge_progress(payload: Dictionary) -> void:
	var charge_count: int = payload.get("charge_count", 0)
	var charge_target: int = payload.get("charge_target", 20)

	# 輕微橙色閃光
	_flash_screen(COLOR_CHARGE_ORANGE, 0.08, 1)

	# 更新進度條
	_update_charge_bar(charge_count, charge_target)

	# 每 5 次顯示浮動文字
	if charge_count % 5 == 0 or charge_count >= charge_target - 3:
		var remaining = charge_target - charge_count
		if remaining > 0:
			_show_float_text("⚡ %d/%d  還差 %d 次！" % [charge_count, charge_target, remaining], COLOR_GOLD, 1.5)
		else:
			_show_float_text("⚡ %d/%d  即將爆發！" % [charge_count, charge_target], COLOR_FIRE, 1.5)

## charge_burst — 全服大爆發（全服廣播）
func _on_charge_burst(payload: Dictionary) -> void:
	var kill_count: int = payload.get("kill_count", 0)
	var burst_mult: float = payload.get("burst_mult", 2.0)
	var total_reward: int = payload.get("total_reward", 0)
	var player_name: String = payload.get("player_name", "某玩家")

	# 停止計時條和進度條
	_stop_timer_bar()
	_clear_charge_bar()

	# 全螢幕三次強閃光（金色，大爆發感）
	_flash_screen(COLOR_GOLD, 0.8, 3)

	# 大字
	_show_big_text("⚡ 全服大爆發！", COLOR_GOLD, 52, 2.5)
	_show_sub_text("擊破 %d 個目標！全服獎勵 +%d！×%.1f 全服共享！" % [kill_count, total_reward, burst_mult], COLOR_CREAM, 2.5)

	# 結算彈窗
	_show_burst_popup(player_name, kill_count, burst_mult, total_reward)

## charge_fail — 充能失敗（全服廣播）
func _on_charge_fail(payload: Dictionary) -> void:
	var charge_count: int = payload.get("charge_count", 0)
	var charge_target: int = payload.get("charge_target", 20)
	var total_reward: int = payload.get("total_reward", 0)

	# 停止計時條和進度條
	_stop_timer_bar()
	_clear_charge_bar()

	# 灰色閃光
	_flash_screen(COLOR_GRAY, 0.3, 1)

	# 提示文字
	if total_reward > 0:
		_show_big_text("⚡ 充能失敗！", COLOR_GRAY, 44, 2.0)
		_show_sub_text("累積 %d/%d 次，安慰獎 +%d！下次加油！" % [charge_count, charge_target, total_reward], COLOR_CREAM, 2.0)
	else:
		_show_big_text("⚡ 充能失敗！", COLOR_GRAY, 44, 2.0)
		_show_sub_text("累積 %d/%d 次，下次加油！" % [charge_count, charge_target], COLOR_CREAM, 2.0)

# ─── 充能進度條 ──────────────────────────────────────────────────────────────

func _show_charge_bar(current: int, target: int) -> void:
	# 先清除舊的
	_clear_charge_bar()

	var vp_size = get_viewport().size
	var bar_w: float = vp_size.x * 0.6
	var bar_h: float = 24.0
	var bar_x: float = vp_size.x * 0.2
	var bar_y: float = vp_size.y - 60.0

	# 背景
	_charge_bar_bg = ColorRect.new()
	_charge_bar_bg.color = Color(0.1, 0.05, 0.0, 0.85)
	_charge_bar_bg.size = Vector2(bar_w, bar_h)
	_charge_bar_bg.position = Vector2(bar_x, bar_y)
	add_child(_charge_bar_bg)

	# 進度條
	var fill_w = bar_w * float(current) / float(target) if target > 0 else 0.0
	_charge_bar = ColorRect.new()
	_charge_bar.color = COLOR_CHARGE_ORANGE
	_charge_bar.size = Vector2(fill_w, bar_h)
	_charge_bar.position = Vector2(bar_x, bar_y)
	add_child(_charge_bar)

	# 標籤
	_charge_label = Label.new()
	_charge_label.text = "⚡ 全服充能 %d/%d" % [current, target]
	_charge_label.add_theme_font_size_override("font_size", 14)
	_charge_label.add_theme_color_override("font_color", COLOR_WHITE)
	_charge_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_charge_label.size = Vector2(bar_w, bar_h)
	_charge_label.position = Vector2(bar_x, bar_y)
	add_child(_charge_label)

func _update_charge_bar(current: int, target: int) -> void:
	if not is_instance_valid(_charge_bar):
		_show_charge_bar(current, target)
		return

	var vp_size = get_viewport().size
	var bar_w: float = vp_size.x * 0.6
	var fill_w = bar_w * float(current) / float(target) if target > 0 else 0.0

	# 動畫更新進度條寬度
	var tween = _charge_bar.create_tween()
	tween.tween_property(_charge_bar, "size:x", fill_w, 0.2).set_ease(Tween.EASE_OUT)

	# 更新標籤
	if is_instance_valid(_charge_label):
		_charge_label.text = "⚡ 全服充能 %d/%d" % [current, target]

	# 接近目標時變金色
	if float(current) / float(target) >= 0.8:
		_charge_bar.color = COLOR_GOLD
		if is_instance_valid(_charge_label):
			_charge_label.add_theme_color_override("font_color", COLOR_GOLD)

func _clear_charge_bar() -> void:
	if is_instance_valid(_charge_bar):
		_charge_bar.queue_free()
		_charge_bar = null
	if is_instance_valid(_charge_bar_bg):
		_charge_bar_bg.queue_free()
		_charge_bar_bg = null
	if is_instance_valid(_charge_label):
		_charge_label.queue_free()
		_charge_label = null

# ─── 大爆發結算彈窗 ──────────────────────────────────────────────────────────

func _show_burst_popup(player_name: String, kill_count: int, burst_mult: float, total_reward: int) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(340, 180)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 90)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.06, 0.0, 0.93)
	style.border_color = COLOR_GOLD
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
	title_lbl.text = "⚡ 全服大爆發結算"
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var pioneer_lbl = Label.new()
	pioneer_lbl.text = "充能先鋒：%s" % player_name
	pioneer_lbl.add_theme_font_size_override("font_size", 13)
	pioneer_lbl.add_theme_color_override("font_color", COLOR_CREAM)
	pioneer_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(pioneer_lbl)

	var kill_lbl = Label.new()
	kill_lbl.text = "擊破全場 %d 個目標！" % kill_count
	kill_lbl.add_theme_font_size_override("font_size", 16)
	kill_lbl.add_theme_color_override("font_color", COLOR_FIRE)
	kill_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(kill_lbl)

	var mult_lbl = Label.new()
	mult_lbl.text = "全服 ×%.1f 共享獎勵" % burst_mult
	mult_lbl.add_theme_font_size_override("font_size", 18)
	mult_lbl.add_theme_color_override("font_color", COLOR_GOLD)
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
	# 先清除舊的
	_stop_timer_bar()

	var vp_size = get_viewport().size
	var bar_h: float = vp_size.y * 0.5
	var bar_w: float = 12.0
	# x=-170 與其他計時條錯開（龍捲風-58，黑洞-72，鏡像分裂-86，量子糾纏-100，武器進化-114，隕石雨-128，龍王-142，時空裂縫-156，全服充能-170）
	var bar_x: float = vp_size.x - 170.0
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.1, 0.05, 0.0, 0.7)
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
