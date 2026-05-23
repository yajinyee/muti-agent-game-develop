## LuckyTimeCapsulePanel.gd — 幸運時間膠囊魚系統面板（DAY-261）
## 業界原創「時間膠囊+預存獎勵+追加存入+膠囊開啟」機制
##
## 視覺設計：
##   - 深藍時間主題（#4A90D9 天藍 + #1A3A5C 深藍 + #FFD700 金 + #E8F4FD 淡藍白）
##   - capsule_start：天藍三次強閃光 + 頂部橫幅 + 「⏳ 時間膠囊！」大字 + 追加存入計數器 + 計時條
##   - capsule_broadcast：頂部小橫幅（全服廣播）
##   - capsule_deposit：天藍閃光 + 「⏳ 追加存入！N/5 ×0.5」浮動文字 + 計數器更新
##   - capsule_open：全螢幕三次強閃光 + 「⏳ 膠囊開啟！」大字 + 結算彈窗
##   - capsule_open_broadcast：全服廣播橫幅
extends CanvasLayer

# 主題顏色
const COLOR_BLUE    = Color("#4A90D9")  # 天藍（主色）
const COLOR_DARK    = Color("#1A3A5C")  # 深藍（背景感）
const COLOR_GOLD    = Color("#FFD700")  # 金色（開啟）
const COLOR_LIGHT   = Color("#E8F4FD")  # 淡藍白（副文字）
const COLOR_WHITE   = Color("#FFFFFF")  # 白色
const COLOR_CYAN    = Color("#00BFFF")  # 青藍（追加存入）

# 計時條
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null

# 追加存入計數器
var _deposit_label: Label = null
var _deposit_count: int = 0
var _max_deposits: int = 5

func _ready() -> void:
	layer = 34  # 幸運時間膠囊魚面板層級（DAY-261）

## 處理幸運時間膠囊魚訊息
func handle_lucky_time_capsule(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"capsule_start":
			_on_capsule_start(payload)
		"capsule_broadcast":
			_on_capsule_broadcast(payload)
		"capsule_deposit":
			_on_capsule_deposit(payload)
		"capsule_open":
			_on_capsule_open(payload)
		"capsule_open_broadcast":
			_on_capsule_open_broadcast(payload)

## capsule_start — 時間膠囊封存啟動（個人訊息）
func _on_capsule_start(payload: Dictionary) -> void:
	var duration_sec: int = payload.get("duration_sec", 15)
	var seal_target: String = payload.get("seal_target", "神秘目標")
	var seal_mult: float = payload.get("seal_mult", 2.5)
	var seal_reward: int = payload.get("seal_reward", 0)
	var max_deposits: int = payload.get("max_deposits", 5)
	var deposit_mult: float = payload.get("deposit_mult", 0.5)
	_deposit_count = 0
	_max_deposits = max_deposits

	# 天藍三次強閃光
	_flash_screen(COLOR_BLUE, 0.6, 3)

	# 頂部橫幅
	_show_banner("⏳ 時間膠囊！封存 %s ×%.1f！每次擊破追加存入 ×%.1f（最多 %d 次）！15秒後開啟！" % [seal_target, seal_mult, deposit_mult, max_deposits], COLOR_BLUE, 4.5)

	# 中央大字
	_show_big_text("⏳ 時間膠囊！", COLOR_BLUE, 50, 3.0)
	_show_sub_text("封存 %s ×%.1f！封存獎勵 +%d！" % [seal_target, seal_mult, seal_reward], COLOR_LIGHT, 3.0)

	# 追加存入計數器（右上角）
	_show_deposit_counter(0, max_deposits)

	# 右側豎向計時條（x=-240 與其他計時條錯開）
	_start_timer_bar(duration_sec, COLOR_BLUE)

## capsule_broadcast — 全服廣播時間膠囊
func _on_capsule_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var seal_target: String = payload.get("seal_target", "神秘目標")
	var seal_mult: float = payload.get("seal_mult", 2.5)
	_show_top_banner("⏳ %s 觸發時間膠囊！封存 %s ×%.1f！15秒後開啟！" % [player_name, seal_target, seal_mult], COLOR_BLUE, 3.5)

## capsule_deposit — 追加存入（個人）
func _on_capsule_deposit(payload: Dictionary) -> void:
	var deposit_count: int = payload.get("deposit_count", 0)
	var max_deposits: int = payload.get("max_deposits", 5)
	var deposit_mult: float = payload.get("deposit_mult", 0.5)
	var reward: int = payload.get("reward", 0)
	var target_name: String = payload.get("target_name", "目標")
	_deposit_count = deposit_count

	# 天藍閃光
	_flash_screen(COLOR_CYAN, 0.18, 1)

	# 更新計數器
	_update_deposit_counter(deposit_count, max_deposits)

	# 浮動文字
	var reward_text = ""
	if reward > 0:
		reward_text = " +%d" % reward
	_show_float_text("⏳ 追加存入！%d/%d ×%.1f%s" % [deposit_count, max_deposits, deposit_mult, reward_text], COLOR_CYAN, 2.0)

## capsule_open — 膠囊開啟（個人）
func _on_capsule_open(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var seal_target: String = payload.get("seal_target", "神秘目標")
	var seal_mult: float = payload.get("seal_mult", 2.5)
	var seal_reward: int = payload.get("seal_reward", 0)
	var deposit_count: int = payload.get("deposit_count", 0)
	var max_deposits: int = payload.get("max_deposits", 5)
	var total_reward: int = payload.get("total_reward", 0)

	# 停止計時條和計數器
	_stop_timer_bar()
	_clear_deposit_counter()

	# 全螢幕三次強閃光（金色）
	_flash_screen(COLOR_GOLD, 0.9, 3)

	# 大字
	_show_big_text("⏳ 膠囊開啟！", COLOR_GOLD, 52, 3.0)
	_show_sub_text("封存 %s ×%.1f + 追加 %d 次！總獎勵 +%d！" % [seal_target, seal_mult, deposit_count, total_reward], COLOR_LIGHT, 3.0)

	# 結算彈窗
	_show_open_popup(player_name, seal_target, seal_mult, seal_reward, deposit_count, max_deposits, total_reward)

## capsule_open_broadcast — 膠囊開啟全服廣播
func _on_capsule_open_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var seal_target: String = payload.get("seal_target", "神秘目標")
	var deposit_count: int = payload.get("deposit_count", 0)
	var total_reward: int = payload.get("total_reward", 0)
	_show_top_banner("⏳ %s 開啟時間膠囊！%s + 追加 %d 次！總獎勵 +%d！" % [player_name, seal_target, deposit_count, total_reward], COLOR_GOLD, 4.0)

# ─── 追加存入計數器 ────────────────────────────────────────────────────────────

func _show_deposit_counter(current: int, max_d: int) -> void:
	_clear_deposit_counter()

	var vp_size = get_viewport().size
	_deposit_label = Label.new()
	_deposit_label.text = "⏳ 存入 %d/%d" % [current, max_d]
	_deposit_label.add_theme_font_size_override("font_size", 15)
	_deposit_label.add_theme_color_override("font_color", COLOR_CYAN)
	_deposit_label.position = Vector2(vp_size.x - 160.0, 80.0)
	_deposit_label.size = Vector2(150, 30)
	_deposit_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	add_child(_deposit_label)

	# 脈衝動畫
	var tween = _deposit_label.create_tween().set_loops()
	tween.tween_property(_deposit_label, "modulate:a", 0.5, 0.7)
	tween.tween_property(_deposit_label, "modulate:a", 1.0, 0.7)

func _update_deposit_counter(current: int, max_d: int) -> void:
	if not is_instance_valid(_deposit_label):
		_show_deposit_counter(current, max_d)
		return
	_deposit_label.text = "⏳ 存入 %d/%d" % [current, max_d]
	# 接近滿時變金色
	if current >= max_d - 1:
		_deposit_label.add_theme_color_override("font_color", COLOR_GOLD)
	# 縮放脈衝
	var tween = create_tween()
	tween.tween_property(_deposit_label, "scale", Vector2(1.2, 1.2), 0.1)
	tween.tween_property(_deposit_label, "scale", Vector2(1.0, 1.0), 0.1)

func _clear_deposit_counter() -> void:
	if is_instance_valid(_deposit_label):
		_deposit_label.queue_free()
		_deposit_label = null

# ─── 膠囊開啟結算彈窗 ─────────────────────────────────────────────────────────

func _show_open_popup(player_name: String, seal_target: String, seal_mult: float, seal_reward: int, deposit_count: int, max_deposits: int, total_reward: int) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(360, 220)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 110)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.04, 0.1, 0.2, 0.93)
	style.border_color = COLOR_GOLD
	style.set_border_width_all(3)
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 7)
	popup.add_child(vbox)

	var title_lbl = Label.new()
	title_lbl.text = "⏳ 膠囊開啟！"
	title_lbl.add_theme_font_size_override("font_size", 26)
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var owner_lbl = Label.new()
	owner_lbl.text = "膠囊主人：%s" % player_name
	owner_lbl.add_theme_font_size_override("font_size", 13)
	owner_lbl.add_theme_color_override("font_color", COLOR_LIGHT)
	owner_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(owner_lbl)

	var seal_lbl = Label.new()
	seal_lbl.text = "封存：%s ×%.1f  +%d" % [seal_target, seal_mult, seal_reward]
	seal_lbl.add_theme_font_size_override("font_size", 14)
	seal_lbl.add_theme_color_override("font_color", COLOR_BLUE)
	seal_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(seal_lbl)

	var dep_lbl = Label.new()
	dep_lbl.text = "追加存入：%d/%d 次" % [deposit_count, max_deposits]
	dep_lbl.add_theme_font_size_override("font_size", 14)
	dep_lbl.add_theme_color_override("font_color", COLOR_CYAN)
	dep_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(dep_lbl)

	var total_lbl = Label.new()
	total_lbl.text = "總獎勵：+%d" % total_reward
	total_lbl.add_theme_font_size_override("font_size", 22)
	total_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	total_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(total_lbl)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 380.0, 0.4).set_ease(Tween.EASE_OUT)
	tween.tween_interval(5.5)
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
	lbl.add_theme_font_size_override("font_size", 13)
	lbl.add_theme_color_override("font_color", Color(0.02, 0.05, 0.1))
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
	lbl.add_theme_color_override("font_color", Color(0.02, 0.05, 0.1))
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
	# x=-240 與其他計時條錯開（...寶藏獵人-226，時間膠囊-240）
	var bar_x: float = vp_size.x - 240.0
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.02, 0.06, 0.12, 0.7)
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
