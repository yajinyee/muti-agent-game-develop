## LuckyDragonKingPanel.gd — 幸運龍王降臨魚系統面板（DAY-254）
## 業界原創「龍王降臨+龍息攻擊+龍王護盾+龍王爆發」機制
##
## 視覺設計：
##   - 深紅龍王主題（#8B0000 + #FF4500 + #FFD700 + #FFF8DC）
##   - dragon_king_start：深紅三次強閃光 + 頂部橫幅 + 「🐉 龍王降臨！」大字 + 護盾指示器 + 計時條
##   - dragon_king_broadcast：頂部小橫幅（全服廣播）
##   - dragon_breath：深紅閃光 + 「🐉 第N次龍息！命中M個！×1.4」浮動文字
##   - dragon_king_burst：全螢幕三次強閃光 + 「🐉 龍王爆發！×3.0 黃金5秒！」大字 + 結算彈窗
##   - dragon_king_burst_broadcast：全服廣播爆發橫幅
extends CanvasLayer

# 主題顏色
const COLOR_DRAGON_RED  = Color("#8B0000")  # 深紅（主色）
const COLOR_FIRE        = Color("#FF4500")  # 火焰橙（爆發）
const COLOR_GOLD        = Color("#FFD700")  # 金色（獎勵）
const COLOR_CREAM       = Color("#FFF8DC")  # 奶油（副文字）
const COLOR_WHITE       = Color("#FFFFFF")  # 白色

# 計時條
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null

# 護盾指示器
var _shield_label: Label = null

# 爆發倍率計時條
var _burst_bar: ColorRect = null
var _burst_tween: Tween = null

func _ready() -> void:
	layer = 27  # 幸運龍王降臨魚面板層級（DAY-254）

## 處理幸運龍王降臨魚訊息
func handle_lucky_dragon_king(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"dragon_king_start":
			_on_dragon_king_start(payload)
		"dragon_king_broadcast":
			_on_dragon_king_broadcast(payload)
		"dragon_breath":
			_on_dragon_breath(payload)
		"dragon_king_burst":
			_on_dragon_king_burst(payload)
		"dragon_king_burst_broadcast":
			_on_dragon_king_burst_broadcast(payload)

## dragon_king_start — 龍王降臨啟動（個人訊息）
func _on_dragon_king_start(payload: Dictionary) -> void:
	var duration_sec: int = payload.get("duration_sec", 15)
	var breath_mult: float = payload.get("breath_mult", 1.4)
	var burst_mult: float = payload.get("burst_mult", 3.0)
	var has_shield: bool = payload.get("has_shield", false)

	# 深紅三次強閃光
	_flash_screen(COLOR_DRAGON_RED, 0.6, 3)

	# 頂部橫幅
	_show_banner("🐉 龍王降臨！每 2 秒龍息攻擊 3 個目標！×%.1f 倍率！全服共享！" % breath_mult, COLOR_DRAGON_RED, 4.0)

	# 中央大字
	_show_big_text("🐉 龍王降臨！", COLOR_DRAGON_RED, 52, 2.5)
	_show_sub_text("龍息 ×%.1f  爆發 ×%.1f（5秒）  龍王護盾已啟動！" % [breath_mult, burst_mult], COLOR_GOLD, 2.5)

	# 護盾指示器
	if has_shield:
		_show_shield_indicator()

	# 右側豎向計時條（x=-142 與其他計時條錯開）
	_start_timer_bar(duration_sec, COLOR_DRAGON_RED)

## dragon_king_broadcast — 全服廣播降臨
func _on_dragon_king_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var breath_mult: float = payload.get("breath_mult", 1.4)
	_show_top_banner("🐉 %s 召喚龍王降臨！龍息 ×%.1f 全服共享！" % [player_name, breath_mult], COLOR_DRAGON_RED, 2.5)

## dragon_breath — 龍息攻擊（全服廣播）
func _on_dragon_breath(payload: Dictionary) -> void:
	var breath_num: int = payload.get("breath_num", 1)
	var hit_count: int = payload.get("hit_count", 0)
	var mult: float = payload.get("mult", 1.4)
	var total_reward: int = payload.get("total_reward", 0)

	# 深紅閃光（輕微）
	_flash_screen(COLOR_DRAGON_RED, 0.18, 1)

	# 浮動文字
	_show_float_text("🐉 第%d次龍息！命中%d個！×%.1f  全服+%d" % [breath_num, hit_count, mult, total_reward], COLOR_FIRE, 1.8)

## dragon_king_burst — 龍王爆發（個人訊息）
func _on_dragon_king_burst(payload: Dictionary) -> void:
	var drain_count: int = payload.get("drain_count", 0)
	var burst_mult: float = payload.get("burst_mult", 3.0)
	var burst_sec: int = payload.get("burst_sec", 5)

	# 停止降臨計時條
	_stop_timer_bar()
	_clear_shield_indicator()

	# 全螢幕三次強閃光（火焰橙，爆發感）
	_flash_screen(COLOR_FIRE, 0.7, 3)

	# 大字
	_show_big_text("🐉 龍王爆發！", COLOR_FIRE, 52, 2.5)
	_show_sub_text("%d 個目標 HP -60%%！×%.1f 倍率加成 %d 秒！" % [drain_count, burst_mult, burst_sec], COLOR_GOLD, 2.5)

	# 爆發倍率計時條（金色，短暫）
	_start_burst_bar(burst_sec)

	# 結算彈窗
	_show_burst_popup(drain_count, burst_mult, burst_sec)

## dragon_king_burst_broadcast — 龍王爆發全服廣播
func _on_dragon_king_burst_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var drain_count: int = payload.get("drain_count", 0)
	_show_top_banner("🐉 %s 龍王爆發！%d 個目標 HP -60%%！" % [player_name, drain_count], COLOR_FIRE, 2.5)

# ─── 護盾指示器 ───────────────────────────────────────────────────────────────

func _show_shield_indicator() -> void:
	var vp_size = get_viewport().size
	_shield_label = Label.new()
	_shield_label.text = "🛡️ 龍王護盾"
	_shield_label.add_theme_font_size_override("font_size", 20)
	_shield_label.add_theme_color_override("font_color", COLOR_GOLD)
	_shield_label.position = Vector2(vp_size.x - 130, vp_size.y * 0.25 - 30)
	add_child(_shield_label)

	# 脈衝動畫
	var tween = _shield_label.create_tween().set_loops()
	tween.tween_property(_shield_label, "modulate:a", 0.5, 0.5)
	tween.tween_property(_shield_label, "modulate:a", 1.0, 0.5)

func _clear_shield_indicator() -> void:
	if is_instance_valid(_shield_label):
		var tween = create_tween()
		tween.tween_property(_shield_label, "modulate:a", 0.0, 0.4)
		tween.tween_callback(_shield_label.queue_free)
		_shield_label = null

# ─── 爆發倍率計時條 ───────────────────────────────────────────────────────────

func _start_burst_bar(duration_sec: int) -> void:
	var vp_size = get_viewport().size
	var bar_h: float = vp_size.y * 0.15
	var bar_w: float = 16.0
	var bar_x: float = vp_size.x - 142.0
	var bar_y: float = vp_size.y * 0.25

	_burst_bar = ColorRect.new()
	_burst_bar.color = COLOR_GOLD
	_burst_bar.size = Vector2(bar_w, bar_h)
	_burst_bar.position = Vector2(bar_x, bar_y)
	add_child(_burst_bar)

	_burst_tween = create_tween()
	_burst_tween.tween_property(_burst_bar, "size:y", 0.0, float(duration_sec)).set_ease(Tween.EASE_IN_OUT)
	_burst_tween.tween_callback(_burst_bar.queue_free)

# ─── 龍王爆發結算彈窗 ────────────────────────────────────────────────────────

func _show_burst_popup(drain_count: int, burst_mult: float, burst_sec: int) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(320, 160)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 80)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.08, 0.0, 0.0, 0.93)
	style.border_color = COLOR_FIRE
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
	title_lbl.text = "🐉 龍王爆發結算"
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", COLOR_FIRE)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var drain_lbl = Label.new()
	drain_lbl.text = "%d 個目標 HP -60%%" % drain_count
	drain_lbl.add_theme_font_size_override("font_size", 15)
	drain_lbl.add_theme_color_override("font_color", COLOR_CREAM)
	drain_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(drain_lbl)

	var mult_lbl = Label.new()
	mult_lbl.text = "個人倍率加成：×%.1f" % burst_mult
	mult_lbl.add_theme_font_size_override("font_size", 18)
	mult_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	var time_lbl = Label.new()
	time_lbl.text = "持續 %d 秒！趕快打！" % burst_sec
	time_lbl.add_theme_font_size_override("font_size", 14)
	time_lbl.add_theme_color_override("font_color", COLOR_FIRE)
	time_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(time_lbl)

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
	# x=-142 與其他計時條錯開（龍捲風-58，黑洞-72，鏡像分裂-86，量子糾纏-100，武器進化-114，隕石雨-128，龍王-142）
	var bar_x: float = vp_size.x - 142.0
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.08, 0.0, 0.0, 0.7)
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
