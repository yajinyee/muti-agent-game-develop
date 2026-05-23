## LuckyTimeRewindPanel.gd — 幸運時光倒流魚系統面板（DAY-247）
## 業界原創「時光倒流+過去擊破重現」機制
##
## 視覺設計：
##   - 紫色時光主題（#9B59B6 + #8E44AD + #D7BDE2 + #F5EEF8）
##   - rewind_start：紫色三次強閃光 + 頂部橫幅 + 「⏪ 時光倒流！」大字 + 重播目標列表 + HP 恢復提示
##   - rewind_broadcast：頂部小橫幅（全服廣播）
##   - rewind_replay：每個重播目標閃現 + 「⏪ 重播！×1.6」浮動文字（依序間隔 400ms）
##   - rewind_end：三次強閃光 + 「⏪ 時光倒流結束！」大字 + 結算彈窗
extends CanvasLayer

# 主題顏色
const COLOR_PURPLE  = Color("#9B59B6")  # 紫色（主題）
const COLOR_DARK    = Color("#8E44AD")  # 深紫（強調）
const COLOR_LIGHT   = Color("#D7BDE2")  # 淺紫（背景）
const COLOR_GOLD    = Color("#F39C12")  # 金色（獎勵）
const COLOR_GREEN   = Color("#27AE60")  # 綠色（HP 恢復）
const COLOR_WHITE   = Color("#FFFFFF")  # 白色

# 重播計數器
var _replay_idx: int = 0
var _total_count: int = 0

func _ready() -> void:
	layer = 20  # 幸運時光倒流魚面板層級（DAY-247）

## 處理幸運時光倒流魚訊息
func handle_lucky_time_rewind(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"rewind_start":
			_on_rewind_start(payload)
		"rewind_broadcast":
			_on_rewind_broadcast(payload)
		"rewind_replay":
			_on_rewind_replay(payload)
		"rewind_end":
			_on_rewind_end(payload)

## rewind_start — 時光倒流啟動（個人訊息）
func _on_rewind_start(payload: Dictionary) -> void:
	_replay_idx = 0
	_total_count = payload.get("replay_count", 0)
	var replay_mult: float = payload.get("replay_mult", 1.6)
	var restored_count: int = payload.get("restored_count", 0)
	var hp_restore_pct: int = payload.get("hp_restore_pct", 60)
	var replay_names: Array = payload.get("replay_names", [])

	# 紫色三次強閃光
	_flash_screen(COLOR_PURPLE, 0.5, 3)

	# 頂部橫幅
	_show_banner("⏪ 時光倒流！重播 %d 個目標！" % _total_count, COLOR_PURPLE, 4.0)

	# 中央大字
	_show_big_text("⏪ 時光倒流！", COLOR_PURPLE, 52, 2.5)

	# 倍率說明
	_show_sub_text("重播倍率 ×%.1f  HP 恢復 %d%%" % [replay_mult, hp_restore_pct], COLOR_GOLD, 2.0)

	# HP 恢復提示（綠色）
	if restored_count > 0:
		_show_hp_restore_text(restored_count, hp_restore_pct)

	# 重播目標列表（右側小面板）
	if replay_names.size() > 0:
		_show_replay_list(replay_names, replay_mult)

## rewind_broadcast — 全服廣播時光倒流啟動
func _on_rewind_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var replay_count: int = payload.get("replay_count", 0)
	_show_top_banner("⏪ %s 觸發時光倒流！重播 %d 個目標！" % [player_name, replay_count], COLOR_PURPLE, 2.5)

## rewind_replay — 每個重播目標（個人訊息）
func _on_rewind_replay(payload: Dictionary) -> void:
	var target_name: String = payload.get("target_name", "目標")
	var reward: int = payload.get("reward", 0)
	var replay_mult: float = payload.get("replay_mult", 1.6)
	var replay_idx: int = payload.get("replay_idx", 1)
	var total_count: int = payload.get("total_count", 1)

	# 紫色閃光
	_flash_screen(COLOR_PURPLE, 0.3, 1)

	# 浮動獎勵文字（帶序號）
	_show_float_text("⏪ [%d/%d] %s ×%.1f  +%d" % [replay_idx, total_count, target_name, replay_mult, reward], COLOR_GOLD, 1.8)

## rewind_end — 時光倒流結束（個人訊息）
func _on_rewind_end(payload: Dictionary) -> void:
	var replay_count: int = payload.get("replay_count", 0)
	var total_reward: int = payload.get("total_reward", 0)

	if replay_count > 0:
		# 三次強閃光
		_flash_screen(COLOR_DARK, 0.55, 3)

		# 中央大字
		_show_big_text("⏪ 時光倒流結束！", COLOR_PURPLE, 48, 2.0)

		# 結算彈窗
		_show_result_popup(replay_count, total_reward)

# ─── HP 恢復提示 ──────────────────────────────────────────────────────────────

func _show_hp_restore_text(restored_count: int, hp_pct: int) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = "💚 %d 個目標 HP 恢復至 %d%%！" % [restored_count, hp_pct]
	lbl.add_theme_font_size_override("font_size", 16)
	lbl.add_theme_color_override("font_color", COLOR_GREEN)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.size = Vector2(vp_size.x, 36)
	lbl.position = Vector2(0, vp_size.y * 0.35 + 130)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(3.0)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.5)
	tween.tween_callback(lbl.queue_free)

# ─── 重播目標列表 ─────────────────────────────────────────────────────────────

func _show_replay_list(names: Array, replay_mult: float) -> void:
	var vp_size = get_viewport().size
	var panel = PanelContainer.new()
	var panel_h: float = 30.0 + names.size() * 24.0
	panel.size = Vector2(220, panel_h)
	panel.position = Vector2(vp_size.x - 230, vp_size.y * 0.3)
	add_child(panel)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.06, 0.02, 0.1, 0.88)
	style.border_color = COLOR_PURPLE
	style.set_border_width_all(2)
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	panel.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 2)
	panel.add_child(vbox)

	var title_lbl = Label.new()
	title_lbl.text = "⏪ 重播目標 ×%.1f" % replay_mult
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", COLOR_LIGHT)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	for i in range(names.size()):
		var item_lbl = Label.new()
		item_lbl.text = "  %d. %s" % [i + 1, names[i]]
		item_lbl.add_theme_font_size_override("font_size", 12)
		item_lbl.add_theme_color_override("font_color", COLOR_GOLD)
		vbox.add_child(item_lbl)

	# 3 秒後淡出
	var tween = create_tween()
	tween.tween_interval(3.0)
	tween.tween_property(panel, "modulate:a", 0.0, 0.5)
	tween.tween_callback(panel.queue_free)

# ─── 結算彈窗 ────────────────────────────────────────────────────────────────

func _show_result_popup(replay_count: int, total_reward: int) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(300, 150)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 75)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.06, 0.02, 0.1, 0.92)
	style.border_color = COLOR_PURPLE
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
	title_lbl.text = "⏪ 時光倒流結算"
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", COLOR_PURPLE)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var count_lbl = Label.new()
	count_lbl.text = "重播目標：%d 個" % replay_count
	count_lbl.add_theme_font_size_override("font_size", 16)
	count_lbl.add_theme_color_override("font_color", COLOR_LIGHT)
	count_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(count_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "總獎勵：+%d 籌碼" % total_reward
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
	lbl.add_theme_font_size_override("font_size", 18)
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
		randf_range(vp_size.x * 0.25, vp_size.x * 0.65),
		randf_range(vp_size.y * 0.3, vp_size.y * 0.6)
	)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 60, duration)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, duration)
	tween.tween_callback(lbl.queue_free)
