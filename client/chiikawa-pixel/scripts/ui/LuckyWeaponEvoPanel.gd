## LuckyWeaponEvoPanel.gd — 幸運武器進化魚系統面板（DAY-252）
## 業界原創「武器進化+穿透+武器爆發」機制
##
## 視覺設計：
##   - 橙紅武器主題（#E67E22 + #E74C3C + #FAD7A0 + #FFF3E0）
##   - weapon_evo_start：橙色三次強閃光 + 頂部橫幅 + 「⚔️ 武器進化！等級2」大字 + 武器等級指示器 + 計時條
##   - weapon_evo_broadcast：頂部小橫幅（全服廣播）
##   - weapon_evo_upgrade：紅色強閃光 + 「⚔️ 武器升級！等級3 穿透！」大字
##   - weapon_evo_pierce：橙色閃光 + 「⚔️ 穿透命中！×0.8」浮動文字
##   - weapon_evo_burst：全螢幕三次強閃光 + 「⚔️ 武器爆發！3連射！」大字 + 結算彈窗
##   - weapon_evo_end：計時條淡出 + 結算彈窗
extends CanvasLayer

# 主題顏色
const COLOR_ORANGE  = Color("#E67E22")  # 橙色（等級2）
const COLOR_RED     = Color("#E74C3C")  # 紅色（等級3）
const COLOR_LIGHT   = Color("#FAD7A0")  # 淺橙（背景）
const COLOR_GOLD    = Color("#F39C12")  # 金色（獎勵）
const COLOR_WHITE   = Color("#FFFFFF")  # 白色

# 計時條
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null

# 武器等級指示器
var _level_indicator: Label = null

# 當前等級
var _current_level: int = 0

func _ready() -> void:
	layer = 25  # 幸運武器進化魚面板層級（DAY-252）

## 處理幸運武器進化魚訊息
func handle_lucky_weapon_evo(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"weapon_evo_start":
			_on_weapon_evo_start(payload)
		"weapon_evo_broadcast":
			_on_weapon_evo_broadcast(payload)
		"weapon_evo_upgrade":
			_on_weapon_evo_upgrade(payload)
		"weapon_evo_pierce":
			_on_weapon_evo_pierce(payload)
		"weapon_evo_burst":
			_on_weapon_evo_burst(payload)
		"weapon_evo_end":
			_on_weapon_evo_end(payload)

## weapon_evo_start — 武器進化啟動（個人訊息）
func _on_weapon_evo_start(payload: Dictionary) -> void:
	var level: int = payload.get("level", 2)
	var duration_sec: int = payload.get("duration_sec", 12)
	var mult: float = payload.get("mult", 1.5)
	var hit_bonus: float = payload.get("hit_bonus", 0.3)
	_current_level = level

	# 橙色三次強閃光
	_flash_screen(COLOR_ORANGE, 0.5, 3)

	# 頂部橫幅
	_show_banner("⚔️ 武器進化！等級 %d！命中率+%d%%，倍率 ×%.1f！" % [level, int(hit_bonus * 100), mult], COLOR_ORANGE, 4.0)

	# 中央大字
	_show_big_text("⚔️ 武器進化！", COLOR_ORANGE, 52, 2.5)
	_show_sub_text("等級 %d  命中率+%d%%  倍率 ×%.1f  再擊破 T210 升等！" % [level, int(hit_bonus * 100), mult], COLOR_GOLD, 2.5)

	# 武器等級指示器
	_show_level_indicator(level)

	# 右側豎向計時條
	_start_timer_bar(duration_sec, COLOR_ORANGE)

## weapon_evo_broadcast — 全服廣播進化
func _on_weapon_evo_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var level: int = payload.get("level", 2)
	var mult: float = payload.get("mult", 1.5)
	var color = COLOR_ORANGE if level < 3 else COLOR_RED
	_show_top_banner("⚔️ %s 武器進化到等級 %d！倍率 ×%.1f！" % [player_name, level, mult], color, 2.5)

## weapon_evo_upgrade — 武器升級到等級 3
func _on_weapon_evo_upgrade(payload: Dictionary) -> void:
	var level: int = payload.get("level", 3)
	var mult: float = payload.get("mult", 2.5)
	_current_level = level

	# 紅色強閃光（升級感）
	_flash_screen(COLOR_RED, 0.6, 3)

	# 大字
	_show_big_text("⚔️ 武器升級！等級 3！", COLOR_RED, 52, 2.5)
	_show_sub_text("穿透效果啟動！倍率 ×%.1f！" % mult, COLOR_GOLD, 2.5)

	# 更新等級指示器
	_update_level_indicator(level)

## weapon_evo_pierce — 穿透命中
func _on_weapon_evo_pierce(payload: Dictionary) -> void:
	var mult: float = payload.get("mult", 0.8)
	var reward: int = payload.get("reward", 0)

	# 橙色閃光
	_flash_screen(COLOR_ORANGE, 0.25, 1)

	# 浮動文字
	_show_float_text("⚔️ 穿透命中！×%.1f  +%d" % [mult, reward], COLOR_LIGHT, 1.5)

## weapon_evo_burst — 武器爆發（全服廣播）
func _on_weapon_evo_burst(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var burst_count: int = payload.get("burst_count", 3)
	var total_reward: int = payload.get("total_reward", 0)
	var burst_mult: float = payload.get("burst_mult", 1.2)

	# 全螢幕三次強閃光（紅色，爆發感）
	_flash_screen(COLOR_RED, 0.65, 3)

	# 大字
	_show_big_text("⚔️ 武器爆發！", COLOR_RED, 52, 2.5)
	_show_sub_text("%s %d 連射！×%.1f  全服+%d！" % [player_name, burst_count, burst_mult, total_reward], COLOR_GOLD, 2.5)

## weapon_evo_end — 進化結束（個人訊息）
func _on_weapon_evo_end(payload: Dictionary) -> void:
	var burst_count: int = payload.get("burst_count", 0)
	var total_reward: int = payload.get("total_reward", 0)
	var burst_mult: float = payload.get("burst_mult", 1.2)

	# 停止計時條
	_stop_timer_bar()
	_clear_level_indicator()
	_current_level = 0

	if burst_count > 0:
		# 全螢幕閃光
		_flash_screen(COLOR_RED, 0.55, 2)
		_show_big_text("⚔️ 武器爆發！%d 連射！" % burst_count, COLOR_RED, 48, 2.5)
		# 結算彈窗
		_show_burst_popup(burst_count, total_reward, burst_mult)
	else:
		_show_float_text("⚔️ 武器進化結束", COLOR_LIGHT, 1.5)

# ─── 武器等級指示器 ──────────────────────────────────────────────────────────

func _show_level_indicator(level: int) -> void:
	var vp_size = get_viewport().size
	_level_indicator = Label.new()
	_level_indicator.text = "⚔️ LV%d" % level
	_level_indicator.add_theme_font_size_override("font_size", 22)
	var color = COLOR_ORANGE if level < 3 else COLOR_RED
	_level_indicator.add_theme_color_override("font_color", color)
	_level_indicator.position = Vector2(vp_size.x - 130, vp_size.y * 0.25 - 30)
	add_child(_level_indicator)

	# 脈衝動畫
	var tween = _level_indicator.create_tween().set_loops()
	tween.tween_property(_level_indicator, "scale", Vector2(1.2, 1.2), 0.4)
	tween.tween_property(_level_indicator, "scale", Vector2(1.0, 1.0), 0.4)

func _update_level_indicator(level: int) -> void:
	if is_instance_valid(_level_indicator):
		_level_indicator.text = "⚔️ LV%d" % level
		var color = COLOR_ORANGE if level < 3 else COLOR_RED
		_level_indicator.add_theme_color_override("font_color", color)

func _clear_level_indicator() -> void:
	if is_instance_valid(_level_indicator):
		var tween = create_tween()
		tween.tween_property(_level_indicator, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_level_indicator.queue_free)
		_level_indicator = null

# ─── 武器爆發結算彈窗 ────────────────────────────────────────────────────────

func _show_burst_popup(burst_count: int, total_reward: int, burst_mult: float) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(300, 150)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 75)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.12, 0.03, 0.02, 0.93)
	style.border_color = COLOR_RED
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
	title_lbl.text = "⚔️ 武器爆發結算"
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", COLOR_RED)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var burst_lbl = Label.new()
	burst_lbl.text = "%d 連射  ×%.1f 倍率" % [burst_count, burst_mult]
	burst_lbl.add_theme_font_size_override("font_size", 16)
	burst_lbl.add_theme_color_override("font_color", COLOR_LIGHT)
	burst_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(burst_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "個人獎勵：+%d 籌碼" % total_reward
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
	lbl.add_theme_font_size_override("font_size", 17)
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
	# x=-114 與其他計時條錯開（龍捲風-58，黑洞-72，鏡像分裂-86，量子糾纏-100，武器進化-114）
	var bar_x: float = vp_size.x - 114.0
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.12, 0.03, 0.02, 0.7)
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
