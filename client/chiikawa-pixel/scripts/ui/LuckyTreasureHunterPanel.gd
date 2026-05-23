## LuckyTreasureHunterPanel.gd — 幸運寶藏獵人魚系統面板（DAY-260）
## 業界原創「寶藏地圖碎片+挖掘+寶藏爆發」機制
##
## 視覺設計：
##   - 古銅寶藏主題（#D4A017 古銅金 + #8B4513 深棕 + #FFD700 金 + #FFF8DC 奶油）
##   - treasure_start：古銅三次強閃光 + 頂部橫幅 + 「🗺️ 寶藏獵人！」大字 + 碎片進度條 + 計時條
##   - treasure_broadcast：頂部小橫幅（全服廣播）
##   - treasure_fragment：金色閃光 + 「🗺️ 發現碎片！N/3 ×1.8」浮動文字 + 碎片進度更新
##   - treasure_burst：全螢幕三次強閃光 + 「🏆 寶藏爆發！×5.0」大字 + 結算彈窗
##   - treasure_burst_broadcast：全服廣播橫幅
##   - treasure_timeout：計時條淡出 + 安慰獎提示
extends CanvasLayer

# 主題顏色
const COLOR_BRONZE  = Color("#D4A017")  # 古銅金（主色）
const COLOR_BROWN   = Color("#8B4513")  # 深棕（地圖感）
const COLOR_GOLD    = Color("#FFD700")  # 金色（爆發）
const COLOR_CREAM   = Color("#FFF8DC")  # 奶油（副文字）
const COLOR_WHITE   = Color("#FFFFFF")  # 白色
const COLOR_RED     = Color("#DC143C")  # 紅色（超時）

# 計時條
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null

# 碎片進度條
var _frag_bar: ColorRect = null
var _frag_bar_bg: ColorRect = null
var _frag_label: Label = null
var _frag_count: int = 0
var _frag_target: int = 3

func _ready() -> void:
	layer = 33  # 幸運寶藏獵人魚面板層級（DAY-260）

## 處理幸運寶藏獵人魚訊息
func handle_lucky_treasure_hunter(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"treasure_start":
			_on_treasure_start(payload)
		"treasure_broadcast":
			_on_treasure_broadcast(payload)
		"treasure_fragment":
			_on_treasure_fragment(payload)
		"treasure_burst":
			_on_treasure_burst(payload)
		"treasure_burst_broadcast":
			_on_treasure_burst_broadcast(payload)
		"treasure_timeout":
			_on_treasure_timeout(payload)

## treasure_start — 寶藏獵人啟動（個人訊息）
func _on_treasure_start(payload: Dictionary) -> void:
	var duration_sec: int = payload.get("duration_sec", 20)
	var frag_target: int = payload.get("frag_target", 3)
	var frag_mult: float = payload.get("frag_mult", 1.8)
	var burst_mult: float = payload.get("burst_mult", 5.0)
	var find_chance: float = payload.get("find_chance", 0.3)
	_frag_target = frag_target
	_frag_count = 0

	# 古銅三次強閃光
	_flash_screen(COLOR_BRONZE, 0.6, 3)

	# 頂部橫幅
	_show_banner("🗺️ 寶藏獵人！每次擊破 %.0f%% 機率發現碎片 ×%.1f！集齊 %d 個→×%.1f 大獎！" % [find_chance * 100, frag_mult, frag_target, burst_mult], COLOR_BRONZE, 4.5)

	# 中央大字
	_show_big_text("🗺️ 寶藏獵人！", COLOR_BRONZE, 50, 3.0)
	_show_sub_text("集齊 %d 個碎片→×%.1f 大獎！每次擊破 %.0f%% 機率發現！" % [frag_target, burst_mult, find_chance * 100], COLOR_CREAM, 3.0)

	# 碎片進度條（底部）
	_show_frag_bar(0, frag_target)

	# 右側豎向計時條（x=-226 與其他計時條錯開）
	_start_timer_bar(duration_sec, COLOR_BRONZE)

## treasure_broadcast — 全服廣播寶藏獵人
func _on_treasure_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var frag_target: int = payload.get("frag_target", 3)
	var burst_mult: float = payload.get("burst_mult", 5.0)
	_show_top_banner("🗺️ %s 觸發寶藏獵人！集齊 %d 個碎片→×%.1f 大獎！" % [player_name, frag_target, burst_mult], COLOR_BRONZE, 3.5)

## treasure_fragment — 發現碎片（個人）
func _on_treasure_fragment(payload: Dictionary) -> void:
	var fragments: int = payload.get("fragments", 0)
	var frag_target: int = payload.get("frag_target", 3)
	var frag_mult: float = payload.get("frag_mult", 1.8)
	var reward: int = payload.get("reward", 0)
	var target_name: String = payload.get("target_name", "目標")
	_frag_count = fragments

	# 金色閃光
	_flash_screen(COLOR_GOLD, 0.22, 1)

	# 更新碎片進度條
	_update_frag_bar(fragments, frag_target)

	# 浮動文字
	var reward_text = ""
	if reward > 0:
		reward_text = " +%d" % reward
	_show_float_text("🗺️ 發現碎片！%d/%d ×%.1f%s" % [fragments, frag_target, frag_mult, reward_text], COLOR_GOLD, 2.0)

## treasure_burst — 寶藏爆發（個人）
func _on_treasure_burst(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var burst_mult: float = payload.get("burst_mult", 5.0)
	var reward: int = payload.get("reward", 0)
	var frag_target: int = payload.get("frag_target", 3)

	# 停止計時條和進度條
	_stop_timer_bar()
	_clear_frag_bar()

	# 全螢幕三次強閃光（金色）
	_flash_screen(COLOR_GOLD, 0.9, 3)

	# 大字
	_show_big_text("🏆 寶藏爆發！", COLOR_GOLD, 54, 3.0)
	_show_sub_text("集齊 %d 個碎片！×%.1f 大獎！+%d！" % [frag_target, burst_mult, reward], COLOR_CREAM, 3.0)

	# 結算彈窗
	_show_burst_popup(player_name, frag_target, burst_mult, reward)

## treasure_burst_broadcast — 寶藏爆發全服廣播
func _on_treasure_burst_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var burst_mult: float = payload.get("burst_mult", 5.0)
	var reward: int = payload.get("reward", 0)
	_show_top_banner("🏆 %s 集齊寶藏！×%.1f 大獎！+%d！" % [player_name, burst_mult, reward], COLOR_GOLD, 4.0)

## treasure_timeout — 寶藏消失（個人）
func _on_treasure_timeout(payload: Dictionary) -> void:
	var fragments: int = payload.get("fragments", 0)
	var frag_target: int = payload.get("frag_target", 3)
	var reward: int = payload.get("reward", 0)

	# 停止計時條和進度條
	_stop_timer_bar()
	_clear_frag_bar()

	if fragments > 0:
		_flash_screen(COLOR_BROWN, 0.3, 1)
		_show_float_text("🗺️ 寶藏消失！收集了 %d/%d 個碎片，安慰獎 +%d" % [fragments, frag_target, reward], COLOR_CREAM, 2.5)
	else:
		_show_float_text("🗺️ 寶藏消失！下次再試！", COLOR_BROWN, 2.0)

# ─── 碎片進度條 ───────────────────────────────────────────────────────────────

func _show_frag_bar(current: int, target: int) -> void:
	_clear_frag_bar()

	var vp_size = get_viewport().size
	var bar_w: float = vp_size.x * 0.4
	var bar_h: float = 18.0
	var bar_x: float = vp_size.x * 0.3
	var bar_y: float = vp_size.y - 60.0

	_frag_bar_bg = ColorRect.new()
	_frag_bar_bg.color = Color(0.1, 0.05, 0.0, 0.8)
	_frag_bar_bg.size = Vector2(bar_w, bar_h)
	_frag_bar_bg.position = Vector2(bar_x, bar_y)
	add_child(_frag_bar_bg)

	_frag_bar = ColorRect.new()
	_frag_bar.color = COLOR_BRONZE
	var fill_w = bar_w * float(current) / float(target) if target > 0 else 0.0
	_frag_bar.size = Vector2(fill_w, bar_h)
	_frag_bar.position = Vector2(bar_x, bar_y)
	add_child(_frag_bar)

	_frag_label = Label.new()
	_frag_label.text = "🗺️ 碎片 %d/%d" % [current, target]
	_frag_label.add_theme_font_size_override("font_size", 13)
	_frag_label.add_theme_color_override("font_color", COLOR_CREAM)
	_frag_label.position = Vector2(bar_x, bar_y - 22)
	_frag_label.size = Vector2(bar_w, 20)
	_frag_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	add_child(_frag_label)

func _update_frag_bar(current: int, target: int) -> void:
	if not is_instance_valid(_frag_bar):
		_show_frag_bar(current, target)
		return

	var vp_size = get_viewport().size
	var bar_w: float = vp_size.x * 0.4
	var fill_w = bar_w * float(current) / float(target) if target > 0 else 0.0

	var tween = create_tween()
	tween.tween_property(_frag_bar, "size:x", fill_w, 0.2).set_ease(Tween.EASE_OUT)

	if is_instance_valid(_frag_label):
		_frag_label.text = "🗺️ 碎片 %d/%d" % [current, target]
		# 接近集齊時變金色
		if current >= target - 1:
			_frag_label.add_theme_color_override("font_color", COLOR_GOLD)
			_frag_bar.color = COLOR_GOLD

func _clear_frag_bar() -> void:
	if is_instance_valid(_frag_bar):
		_frag_bar.queue_free()
		_frag_bar = null
	if is_instance_valid(_frag_bar_bg):
		_frag_bar_bg.queue_free()
		_frag_bar_bg = null
	if is_instance_valid(_frag_label):
		_frag_label.queue_free()
		_frag_label = null

# ─── 寶藏爆發結算彈窗 ─────────────────────────────────────────────────────────

func _show_burst_popup(player_name: String, frag_count: int, burst_mult: float, reward: int) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(340, 190)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 95)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.08, 0.04, 0.0, 0.93)
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
	title_lbl.text = "🏆 寶藏爆發！"
	title_lbl.add_theme_font_size_override("font_size", 24)
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var hunter_lbl = Label.new()
	hunter_lbl.text = "寶藏獵人：%s" % player_name
	hunter_lbl.add_theme_font_size_override("font_size", 13)
	hunter_lbl.add_theme_color_override("font_color", COLOR_CREAM)
	hunter_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(hunter_lbl)

	var frag_lbl = Label.new()
	frag_lbl.text = "集齊 %d 個碎片！" % frag_count
	frag_lbl.add_theme_font_size_override("font_size", 16)
	frag_lbl.add_theme_color_override("font_color", COLOR_BRONZE)
	frag_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(frag_lbl)

	var mult_lbl = Label.new()
	mult_lbl.text = "×%.1f 大獎！" % burst_mult
	mult_lbl.add_theme_font_size_override("font_size", 22)
	mult_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "獎勵：+%d" % reward
	reward_lbl.add_theme_font_size_override("font_size", 18)
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
	lbl.add_theme_font_size_override("font_size", 14)
	lbl.add_theme_color_override("font_color", Color(0.05, 0.02, 0.0))
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
	lbl.add_theme_color_override("font_color", Color(0.05, 0.02, 0.0))
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
	# x=-226 與其他計時條錯開（...閃電風暴-198，星座命運-212，寶藏獵人-226）
	var bar_x: float = vp_size.x - 226.0
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.08, 0.04, 0.0, 0.7)
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
