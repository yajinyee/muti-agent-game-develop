## LuckyMeteorShowerPanel.gd — 幸運星際隕石魚系統面板（DAY-253）
## 業界原創「隕石雨+隨機轟炸+隕石連擊+最終隕石」機制
##
## 視覺設計：
##   - 深紅隕石主題（#C0392B + #E74C3C + #F39C12 + #FFF3E0）
##   - meteor_start：深紅三次強閃光 + 頂部橫幅 + 「☄️ 星際隕石雨！」大字 + 計時條
##   - meteor_bomb：深紅閃光 + 「☄️ 第N輪 轟炸命中！×1.3」浮動文字 + 隕石落點特效
##   - meteor_combo：全螢幕三次強閃光 + 「☄️ 隕石連擊！×3.0」大字 + 連擊計數器
##   - meteor_final：全螢幕三次強閃光 + 「☄️ 最終隕石！×2.0」大字 + 結算彈窗
##   - meteor_end：計時條淡出 + 結束提示
extends CanvasLayer

# 主題顏色
const COLOR_DARK_RED  = Color("#C0392B")  # 深紅（主色）
const COLOR_RED       = Color("#E74C3C")  # 紅色（連擊/最終）
const COLOR_ORANGE    = Color("#E67E22")  # 橙色（轟炸）
const COLOR_GOLD      = Color("#F39C12")  # 金色（獎勵）
const COLOR_LIGHT     = Color("#FAD7A0")  # 淺橙（副文字）
const COLOR_WHITE     = Color("#FFFFFF")  # 白色

# 計時條
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null

# 連擊計數器
var _combo_label: Label = null
var _combo_count: int = 0

func _ready() -> void:
	layer = 26  # 幸運星際隕石魚面板層級（DAY-253）

## 處理幸運星際隕石魚訊息
func handle_lucky_meteor_shower(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"meteor_start":
			_on_meteor_start(payload)
		"meteor_bomb":
			_on_meteor_bomb(payload)
		"meteor_combo":
			_on_meteor_combo(payload)
		"meteor_final":
			_on_meteor_final(payload)
		"meteor_end":
			_on_meteor_end(payload)

## meteor_start — 隕石雨開始（全服廣播）
func _on_meteor_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var duration_sec: int = payload.get("duration_sec", 8)
	var bomb_mult: float = payload.get("bomb_mult", 1.3)
	var final_mult: float = payload.get("final_mult", 2.0)
	_combo_count = 0

	# 深紅三次強閃光
	_flash_screen(COLOR_DARK_RED, 0.55, 3)

	# 頂部橫幅
	_show_banner("☄️ %s 觸發星際隕石雨！每秒轟炸 2 個目標！×%.1f 倍率！全服共享！" % [player_name, bomb_mult], COLOR_DARK_RED, 4.0)

	# 中央大字
	_show_big_text("☄️ 星際隕石雨！", COLOR_DARK_RED, 52, 2.5)
	_show_sub_text("轟炸倍率 ×%.1f  最終隕石 ×%.1f  連擊 ×3.0  全服共享！" % [bomb_mult, final_mult], COLOR_GOLD, 2.5)

	# 右側豎向計時條（x=-128 與其他計時條錯開）
	_start_timer_bar(duration_sec, COLOR_DARK_RED)

## meteor_bomb — 隕石轟炸命中（全服廣播）
func _on_meteor_bomb(payload: Dictionary) -> void:
	var round_num: int = payload.get("round", 1)
	var target_name: String = payload.get("target_name", "目標")
	var mult: float = payload.get("mult", 1.3)
	var reward: int = payload.get("reward", 0)
	var x: float = payload.get("x", 0.0)
	var y: float = payload.get("y", 0.0)

	# 深紅閃光（輕微）
	_flash_screen(COLOR_DARK_RED, 0.2, 1)

	# 隕石落點特效（在目標位置）
	if x > 0 and y > 0:
		_show_meteor_impact(x, y)

	# 浮動文字
	_show_float_text_at("☄️ 第%d輪 %s ×%.1f  +%d" % [round_num, target_name, mult, reward], COLOR_ORANGE, 1.8, x, y)

## meteor_combo — 隕石連擊（全服廣播）
func _on_meteor_combo(payload: Dictionary) -> void:
	var target_name: String = payload.get("target_name", "目標")
	var mult: float = payload.get("mult", 3.0)
	var reward: int = payload.get("reward", 0)
	var x: float = payload.get("x", 0.0)
	var y: float = payload.get("y", 0.0)
	_combo_count += 1

	# 全螢幕三次強閃光（紅色，連擊感）
	_flash_screen(COLOR_RED, 0.65, 3)

	# 大字
	_show_big_text("☄️ 隕石連擊！", COLOR_RED, 52, 2.5)
	_show_sub_text("%s 被連續命中 3 次！×%.1f  全服+%d！" % [target_name, mult, reward], COLOR_GOLD, 2.5)

	# 連擊計數器更新
	_update_combo_counter(_combo_count)

	# 隕石落點特效（更大）
	if x > 0 and y > 0:
		_show_meteor_impact_big(x, y)

## meteor_final — 最終隕石（全服廣播）
func _on_meteor_final(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var target_name: String = payload.get("target_name", "目標")
	var mult: float = payload.get("mult", 2.0)
	var reward: int = payload.get("reward", 0)
	var x: float = payload.get("x", 0.0)
	var y: float = payload.get("y", 0.0)

	# 停止計時條
	_stop_timer_bar()

	# 全螢幕三次強閃光（深紅，最終感）
	_flash_screen(COLOR_DARK_RED, 0.7, 3)

	# 大字
	_show_big_text("☄️ 最終隕石！", COLOR_DARK_RED, 52, 2.5)
	_show_sub_text("命中 %s！×%.1f  全服+%d！" % [target_name, mult, reward], COLOR_GOLD, 2.5)

	# 最終隕石落點特效（最大）
	if x > 0 and y > 0:
		_show_meteor_impact_final(x, y)

	# 結算彈窗
	_show_final_popup(player_name, target_name, mult, reward)

## meteor_end — 隕石雨結束（全服廣播，無最終隕石時）
func _on_meteor_end(payload: Dictionary) -> void:
	_stop_timer_bar()
	_clear_combo_counter()
	_combo_count = 0
	_show_float_text("☄️ 隕石雨結束", COLOR_LIGHT, 1.5)

# ─── 隕石落點特效 ─────────────────────────────────────────────────────────────

func _show_meteor_impact(x: float, y: float) -> void:
	var vp_size = get_viewport().size
	# 轉換遊戲座標到螢幕座標（假設遊戲寬度 1000px）
	var screen_x = x / 1000.0 * vp_size.x
	var screen_y = y / 600.0 * vp_size.y

	# 隕石落點圓圈（橙色）
	var impact = ColorRect.new()
	impact.color = Color(COLOR_ORANGE.r, COLOR_ORANGE.g, COLOR_ORANGE.b, 0.7)
	impact.size = Vector2(40, 40)
	impact.position = Vector2(screen_x - 20, screen_y - 20)
	add_child(impact)

	var tween = create_tween()
	tween.tween_property(impact, "scale", Vector2(2.0, 2.0), 0.3).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(impact, "modulate:a", 0.0, 0.3)
	tween.tween_callback(impact.queue_free)

func _show_meteor_impact_big(x: float, y: float) -> void:
	var vp_size = get_viewport().size
	var screen_x = x / 1000.0 * vp_size.x
	var screen_y = y / 600.0 * vp_size.y

	# 連擊落點（紅色，更大）
	var impact = ColorRect.new()
	impact.color = Color(COLOR_RED.r, COLOR_RED.g, COLOR_RED.b, 0.8)
	impact.size = Vector2(64, 64)
	impact.position = Vector2(screen_x - 32, screen_y - 32)
	add_child(impact)

	var tween = create_tween()
	tween.tween_property(impact, "scale", Vector2(2.5, 2.5), 0.4).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(impact, "modulate:a", 0.0, 0.4)
	tween.tween_callback(impact.queue_free)

func _show_meteor_impact_final(x: float, y: float) -> void:
	var vp_size = get_viewport().size
	var screen_x = x / 1000.0 * vp_size.x
	var screen_y = y / 600.0 * vp_size.y

	# 最終隕石落點（深紅，最大）
	var impact = ColorRect.new()
	impact.color = Color(COLOR_DARK_RED.r, COLOR_DARK_RED.g, COLOR_DARK_RED.b, 0.9)
	impact.size = Vector2(96, 96)
	impact.position = Vector2(screen_x - 48, screen_y - 48)
	add_child(impact)

	var tween = create_tween()
	tween.tween_property(impact, "scale", Vector2(3.0, 3.0), 0.5).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(impact, "modulate:a", 0.0, 0.5)
	tween.tween_callback(impact.queue_free)

# ─── 連擊計數器 ───────────────────────────────────────────────────────────────

func _update_combo_counter(count: int) -> void:
	if not is_instance_valid(_combo_label):
		var vp_size = get_viewport().size
		_combo_label = Label.new()
		_combo_label.add_theme_font_size_override("font_size", 20)
		_combo_label.add_theme_color_override("font_color", COLOR_RED)
		_combo_label.position = Vector2(vp_size.x - 130, vp_size.y * 0.25 - 30)
		add_child(_combo_label)
		# 脈衝動畫
		var tween = _combo_label.create_tween().set_loops()
		tween.tween_property(_combo_label, "scale", Vector2(1.2, 1.2), 0.3)
		tween.tween_property(_combo_label, "scale", Vector2(1.0, 1.0), 0.3)

	_combo_label.text = "☄️ 連擊 ×%d" % count

func _clear_combo_counter() -> void:
	if is_instance_valid(_combo_label):
		var tween = create_tween()
		tween.tween_property(_combo_label, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_combo_label.queue_free)
		_combo_label = null

# ─── 最終隕石結算彈窗 ────────────────────────────────────────────────────────

func _show_final_popup(player_name: String, target_name: String, mult: float, reward: int) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(320, 160)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 80)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.02, 0.02, 0.93)
	style.border_color = COLOR_DARK_RED
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
	title_lbl.text = "☄️ 最終隕石結算"
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", COLOR_DARK_RED)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var target_lbl = Label.new()
	target_lbl.text = "命中目標：%s" % target_name
	target_lbl.add_theme_font_size_override("font_size", 15)
	target_lbl.add_theme_color_override("font_color", COLOR_LIGHT)
	target_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(target_lbl)

	var mult_lbl = Label.new()
	mult_lbl.text = "倍率：×%.1f  觸發者：%s" % [mult, player_name]
	mult_lbl.add_theme_font_size_override("font_size", 14)
	mult_lbl.add_theme_color_override("font_color", COLOR_ORANGE)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "全服共享：+%d 籌碼" % reward
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

func _show_float_text_at(text: String, color: Color, duration: float, x: float, y: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 16)
	lbl.add_theme_color_override("font_color", color)
	# 轉換遊戲座標到螢幕座標
	var screen_x = x / 1000.0 * vp_size.x if x > 0 else randf_range(vp_size.x * 0.3, vp_size.x * 0.7)
	var screen_y = y / 600.0 * vp_size.y if y > 0 else randf_range(vp_size.y * 0.3, vp_size.y * 0.6)
	lbl.position = Vector2(screen_x, screen_y)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "position:y", screen_y - 60, duration)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, duration)
	tween.tween_callback(lbl.queue_free)

func _start_timer_bar(duration_sec: int, color: Color) -> void:
	var vp_size = get_viewport().size
	var bar_h: float = vp_size.y * 0.5
	var bar_w: float = 12.0
	# x=-128 與其他計時條錯開（龍捲風-58，黑洞-72，鏡像分裂-86，量子糾纏-100，武器進化-114，隕石雨-128）
	var bar_x: float = vp_size.x - 128.0
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.1, 0.02, 0.02, 0.7)
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
