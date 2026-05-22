## GoldenWaveFishPanel.gd — 黃金波浪魚全場倍率衝擊面板（DAY-207）
## 業界依據：Ocean King 4 Brand New World
## 「Golden Wave Fish — triggers a golden tidal wave that sweeps across the entire screen,
##  temporarily boosting all multipliers by 2x for 8 seconds.」
##
## 視覺設計：
##   - 黃金海浪主題（#FFD700 + #FF8C00 + #FFF8DC + #FFFFFF）
##   - wave_start：金色強閃光 + 頂部橫幅 + 波浪預告
##   - wave_column：黃金波浪柱（從左到右依序出現，每 150ms 一列）
##   - boost_start：全螢幕金色爆炸 + 「×2 黃金加成！」大字 + 底部計時進度條
##   - boost_end：進度條淡出 + 結算彈窗
extends CanvasLayer

var _boost_bar: Control = null
var _boost_tween: Tween = null

func _ready() -> void:
	layer = 38  # 黃金波浪面板層級

## 處理黃金波浪魚訊息
func handle_golden_wave_fish(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"wave_start":
			_on_wave_start(payload)
		"wave_column":
			_on_wave_column(payload)
		"boost_start":
			_on_boost_start(payload)
		"boost_end":
			_on_boost_end()

## 波浪開始 — 金色強閃光 + 頂部橫幅
func _on_wave_start(payload: Dictionary) -> void:
	var columns = payload.get("columns", 8)
	var boost_mult = payload.get("boost_mult", 2.0)
	var boost_sec = payload.get("boost_sec", 8)

	# 金色強閃光
	_flash_screen(Color("#FFD700"), 0.7)

	# 頂部橫幅
	var banner = _make_banner(
		"🌊 黃金波浪！掃場 %d 列 → ×%.0f 倍率加成 %d 秒！" % [columns, boost_mult, boost_sec],
		Color("#FF8C00"),
		Color("#1A0A00")
	)
	add_child(banner)
	var tween = create_tween()
	tween.tween_property(banner, "position:y", 10.0, 0.3)
	tween.tween_interval(3.0)
	tween.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(banner.queue_free)

## 波浪掃過一列 — 黃金波浪柱動畫
func _on_wave_column(payload: Dictionary) -> void:
	var col_x = payload.get("col_x", 640.0)
	var kills = payload.get("kill_count", 0)
	var reward = payload.get("reward", 0)
	var col_index = payload.get("col_index", 0)

	# 黃金波浪柱（從底部升起）
	var wave_col = ColorRect.new()
	wave_col.color = Color(1.0, 0.85, 0.0, 0.6)
	wave_col.size = Vector2(160, 0)
	wave_col.position = Vector2(col_x - 80, 720)
	add_child(wave_col)

	var tween = create_tween()
	tween.tween_property(wave_col, "size:y", 720.0, 0.12)
	tween.parallel().tween_property(wave_col, "position:y", 0.0, 0.12)
	tween.tween_property(wave_col, "modulate:a", 0.0, 0.25)
	tween.tween_callback(wave_col.queue_free)

	# 擊破浮動文字
	if kills > 0:
		var float_label = Label.new()
		float_label.text = "💰 ×%d  +%d" % [kills, reward]
		float_label.add_theme_font_size_override("font_size", 16)
		float_label.add_theme_color_override("font_color", Color("#FFD700"))
		float_label.position = Vector2(col_x - 30, 300)
		add_child(float_label)
		var tween2 = create_tween()
		tween2.tween_property(float_label, "position:y", 240.0, 0.7)
		tween2.parallel().tween_property(float_label, "modulate:a", 0.0, 0.7)
		tween2.tween_callback(float_label.queue_free)

	# 最後一列（col_index == 7）額外強閃光
	if col_index >= 7:
		_flash_screen(Color("#FFD700"), 0.5)

## 黃金加成開始 — 全螢幕金色爆炸 + 大字 + 底部計時進度條
func _on_boost_start(payload: Dictionary) -> void:
	var boost_mult = payload.get("boost_mult", 2.0)
	var boost_sec = payload.get("boost_sec", 8)
	var total_kills = payload.get("total_kills", 0)
	var total_reward = payload.get("total_reward", 0)

	# 全螢幕三次金色強閃光
	_flash_screen(Color("#FFD700"), 0.9)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color("#FF8C00"), 0.7)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color("#FFFFFF"), 0.6)

	# 「×2 黃金加成！」大字
	var big_label = Label.new()
	big_label.text = "✨ ×%.0f 黃金加成！" % boost_mult
	big_label.add_theme_font_size_override("font_size", 48)
	big_label.add_theme_color_override("font_color", Color("#FFD700"))
	big_label.set_anchors_preset(Control.PRESET_CENTER)
	big_label.position = Vector2(-200, -60)
	add_child(big_label)
	var tween = create_tween()
	tween.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.25)
	tween.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.2)
	tween.tween_interval(1.5)
	tween.tween_property(big_label, "modulate:a", 0.0, 0.6)
	tween.tween_callback(big_label.queue_free)

	# 底部計時進度條
	_show_boost_bar(boost_sec)

	# 結算彈窗（右側滑入）
	if total_kills > 0:
		await get_tree().create_timer(0.5).timeout
		_show_wave_result(total_kills, total_reward, boost_mult, boost_sec)

## 黃金加成結束
func _on_boost_end() -> void:
	if is_instance_valid(_boost_bar):
		var tween = create_tween()
		tween.tween_property(_boost_bar, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_boost_bar.queue_free)
		_boost_bar = null

## 底部計時進度條
func _show_boost_bar(boost_sec: int) -> void:
	# 清除舊的進度條
	if is_instance_valid(_boost_bar):
		_boost_bar.queue_free()

	var bar_container = Control.new()
	bar_container.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	bar_container.position = Vector2(0, -30)
	bar_container.size = Vector2(1280, 30)
	add_child(bar_container)
	_boost_bar = bar_container

	# 背景
	var bg = ColorRect.new()
	bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.1, 0.08, 0.0, 0.85)
	bar_container.add_child(bg)

	# 進度條
	var bar = ColorRect.new()
	bar.name = "BoostBar"
	bar.set_anchors_preset(Control.PRESET_FULL_RECT)
	bar.color = Color("#FFD700")
	bar_container.add_child(bar)

	# 文字
	var label = Label.new()
	label.text = "✨ ×2 黃金加成進行中..."
	label.add_theme_font_size_override("font_size", 14)
	label.add_theme_color_override("font_color", Color("#1A0A00"))
	label.set_anchors_preset(Control.PRESET_FULL_RECT)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	bar_container.add_child(label)

	# 進度條縮短動畫（顏色漸變：金→橙→紅橙）
	_boost_tween = create_tween()
	_boost_tween.tween_property(bar, "size:x", 0.0, float(boost_sec))
	_boost_tween.parallel().tween_method(
		func(t: float):
			if is_instance_valid(bar):
				if t > 0.6:
					bar.color = Color("#FFD700")
				elif t > 0.3:
					bar.color = Color("#FF8C00")
				else:
					bar.color = Color("#FF4500"),
		1.0, 0.0, float(boost_sec)
	)

## 波浪結算彈窗
func _show_wave_result(kills: int, reward: int, boost_mult: float, boost_sec: int) -> void:
	var popup = PanelContainer.new()
	popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	popup.position = Vector2(1400, -90)
	popup.size = Vector2(270, 170)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.08, 0.0, 0.92)
	style.border_color = Color("#FFD700")
	style.set_border_width_all(3)
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup.add_child(vbox)

	var title_lbl = Label.new()
	title_lbl.text = "🌊 黃金波浪結算"
	title_lbl.add_theme_font_size_override("font_size", 16)
	title_lbl.add_theme_color_override("font_color", Color("#FFD700"))
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var kills_lbl = Label.new()
	kills_lbl.text = "波浪擊破：%d 個" % kills
	kills_lbl.add_theme_font_size_override("font_size", 14)
	kills_lbl.add_theme_color_override("font_color", Color("#FF8C00"))
	kills_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(kills_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "全服獎勵：%d 金幣" % reward
	reward_lbl.add_theme_font_size_override("font_size", 14)
	reward_lbl.add_theme_color_override("font_color", Color("#FFFFFF"))
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	var boost_lbl = Label.new()
	boost_lbl.text = "✨ ×%.0f 倍率加成 %d 秒！" % [boost_mult, boost_sec]
	boost_lbl.add_theme_font_size_override("font_size", 14)
	boost_lbl.add_theme_color_override("font_color", Color("#FFD700"))
	boost_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(boost_lbl)

	add_child(popup)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", 1000.0, 0.35)
	tween.tween_interval(4.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.6)
	tween.tween_callback(popup.queue_free)

## 橫幅工廠
func _make_banner(text: String, bg_color: Color, text_color: Color) -> PanelContainer:
	var panel = PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_WIDE)
	panel.position = Vector2(0, -60)
	panel.size = Vector2(1280, 50)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(bg_color.r, bg_color.g, bg_color.b, 0.88)
	panel.add_theme_stylebox_override("panel", style)

	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 20)
	label.add_theme_color_override("font_color", text_color)
	label.set_anchors_preset(Control.PRESET_FULL_RECT)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	panel.add_child(label)
	return panel

## 全螢幕閃光
func _flash_screen(color: Color, alpha: float) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, alpha)
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.25)
	tween.tween_callback(flash.queue_free)
