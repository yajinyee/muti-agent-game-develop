## FortuneCoinFishPanel.gd — 幸運金幣魚即時獎勵面板（DAY-209）
## 業界依據：Galaxsys King of Ocean 2026
## 「Money Fish trigger instant payouts.」
##
## 視覺設計：
##   - 金幣黃色主題（#FFD700 + #FF8C00 + #FFF8DC + #FFFFFF）
##   - coin_burst：金色強閃光 + 「💰 ×N」大字彈跳 + 右側滑入獎勵彈窗
##   - coin_broadcast：頂部小橫幅（讓其他玩家看到）
##   - golden_burst_start：全螢幕金色爆炸 + 「💥 黃金爆發！」大字 + 底部計時條
##   - golden_burst_end：計時條淡出
extends CanvasLayer

var _burst_bar: Control = null

func _ready() -> void:
	layer = 36  # 幸運金幣魚面板層級

## 處理幸運金幣魚訊息
func handle_fortune_coin_fish(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"coin_burst":
			_on_coin_burst(payload)
		"coin_broadcast":
			_on_coin_broadcast(payload)
		"golden_burst_start":
			_on_golden_burst_start(payload)
		"golden_burst_end":
			_on_golden_burst_end()

## 個人金幣爆發 — 金色強閃光 + 大字彈跳 + 獎勵彈窗
func _on_coin_burst(payload: Dictionary) -> void:
	var multiplier: int = payload.get("multiplier", 5)
	var reward: int = payload.get("reward", 0)
	var label: String = payload.get("label", "💰 ×5")

	# 金色強閃光（倍率越高越強）
	var flash_alpha = 0.5 + float(multiplier) / 100.0
	_flash_screen(Color("#FFD700"), clampf(flash_alpha, 0.5, 0.95))

	# 「💰 ×N」大字彈跳（中央）
	var font_size = 36 + multiplier  # 倍率越高字越大
	font_size = clampi(font_size, 36, 72)
	var big_label = Label.new()
	big_label.text = label
	big_label.add_theme_font_size_override("font_size", font_size)
	big_label.add_theme_color_override("font_color", Color("#FFD700"))
	big_label.set_anchors_preset(Control.PRESET_CENTER)
	big_label.position = Vector2(-80, -50)
	add_child(big_label)
	var tween = create_tween()
	tween.tween_property(big_label, "scale", Vector2(1.4, 1.4), 0.2)
	tween.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.15)
	tween.tween_interval(0.8)
	tween.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween.tween_callback(big_label.queue_free)

	# 右側滑入獎勵彈窗
	_show_coin_popup(multiplier, reward, label)

## 全服廣播橫幅（讓其他玩家看到）
func _on_coin_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "玩家")
	var label: String = payload.get("label", "💰 ×5")

	var banner = _make_small_banner(
		"💰 %s 觸發幸運金幣魚！%s" % [player_name, label],
		Color("#8B6914"),
		Color("#FFD700")
	)
	add_child(banner)
	var tween = create_tween()
	tween.tween_property(banner, "position:y", 60.0, 0.25)
	tween.tween_interval(2.5)
	tween.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(banner.queue_free)

## 黃金爆發開始 — 全螢幕金色爆炸 + 大字 + 底部計時條
func _on_golden_burst_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "玩家")
	var affected: int = payload.get("affected_count", 0)
	var burst_sec: int = payload.get("burst_sec", 5)

	# 全螢幕三次金色強閃光
	_flash_screen(Color("#FFD700"), 1.0)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color("#FF8C00"), 0.8)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color("#FFFFFF"), 0.6)

	# 「💥 黃金爆發！」大字
	var big_label = Label.new()
	big_label.text = "💥 黃金爆發！"
	big_label.add_theme_font_size_override("font_size", 52)
	big_label.add_theme_color_override("font_color", Color("#FFD700"))
	big_label.set_anchors_preset(Control.PRESET_CENTER)
	big_label.position = Vector2(-200, -70)
	add_child(big_label)
	var tween = create_tween()
	tween.tween_property(big_label, "scale", Vector2(1.3, 1.3), 0.25)
	tween.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.2)
	tween.tween_interval(1.5)
	tween.tween_property(big_label, "modulate:a", 0.0, 0.5)
	tween.tween_callback(big_label.queue_free)

	# 副標題
	var sub_label = Label.new()
	sub_label.text = "%s 觸發！%d 個目標 HP -80%%！快打！" % [player_name, affected]
	sub_label.add_theme_font_size_override("font_size", 20)
	sub_label.add_theme_color_override("font_color", Color("#FF8C00"))
	sub_label.set_anchors_preset(Control.PRESET_CENTER)
	sub_label.position = Vector2(-240, -10)
	add_child(sub_label)
	var tween2 = create_tween()
	tween2.tween_interval(2.0)
	tween2.tween_property(sub_label, "modulate:a", 0.0, 0.5)
	tween2.tween_callback(sub_label.queue_free)

	# 底部計時條
	_show_burst_bar(burst_sec)

## 黃金爆發結束
func _on_golden_burst_end() -> void:
	if is_instance_valid(_burst_bar):
		var tween = create_tween()
		tween.tween_property(_burst_bar, "modulate:a", 0.0, 0.4)
		tween.tween_callback(_burst_bar.queue_free)
		_burst_bar = null

## 個人獎勵彈窗
func _show_coin_popup(multiplier: int, reward: int, label: String) -> void:
	var popup = PanelContainer.new()
	popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	popup.position = Vector2(1400, -80)
	popup.size = Vector2(240, 130)

	var border_color = Color("#FFD700")
	if multiplier >= 50:
		border_color = Color("#FF8C00")
	elif multiplier >= 20:
		border_color = Color("#FFD700")

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.08, 0.0, 0.92)
	style.border_color = border_color
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
	title_lbl.text = "💰 幸運金幣魚！"
	title_lbl.add_theme_font_size_override("font_size", 16)
	title_lbl.add_theme_color_override("font_color", Color("#FFD700"))
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var mult_lbl = Label.new()
	mult_lbl.text = label
	mult_lbl.add_theme_font_size_override("font_size", 22)
	mult_lbl.add_theme_color_override("font_color", border_color)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "+%d 金幣" % reward
	reward_lbl.add_theme_font_size_override("font_size", 14)
	reward_lbl.add_theme_color_override("font_color", Color("#FFFFFF"))
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	add_child(popup)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", 1030.0, 0.3)
	tween.tween_interval(2.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)

## 底部黃金爆發計時條
func _show_burst_bar(burst_sec: int) -> void:
	if is_instance_valid(_burst_bar):
		_burst_bar.queue_free()

	var bar_container = Control.new()
	bar_container.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	bar_container.position = Vector2(0, -30)
	bar_container.size = Vector2(1280, 30)
	add_child(bar_container)
	_burst_bar = bar_container

	var bg = ColorRect.new()
	bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.1, 0.08, 0.0, 0.85)
	bar_container.add_child(bg)

	var bar = ColorRect.new()
	bar.name = "BurstBar"
	bar.set_anchors_preset(Control.PRESET_FULL_RECT)
	bar.color = Color("#FFD700")
	bar_container.add_child(bar)

	var label = Label.new()
	label.text = "💥 黃金爆發！全場 HP -80%%！快打！"
	label.add_theme_font_size_override("font_size", 14)
	label.add_theme_color_override("font_color", Color("#1A0A00"))
	label.set_anchors_preset(Control.PRESET_FULL_RECT)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	bar_container.add_child(label)

	var tween = create_tween()
	tween.tween_property(bar, "size:x", 0.0, float(burst_sec))
	tween.parallel().tween_method(
		func(t: float):
			if is_instance_valid(bar):
				if t > 0.5:
					bar.color = Color("#FFD700")
				elif t > 0.25:
					bar.color = Color("#FF8C00")
				else:
					bar.color = Color("#FF4500"),
		1.0, 0.0, float(burst_sec)
	)

## 小橫幅（全服廣播用）
func _make_small_banner(text: String, bg_color: Color, text_color: Color) -> PanelContainer:
	var panel = PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_WIDE)
	panel.position = Vector2(0, -40)
	panel.size = Vector2(1280, 36)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(bg_color.r, bg_color.g, bg_color.b, 0.80)
	panel.add_theme_stylebox_override("panel", style)

	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 16)
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
	tween.tween_property(flash, "modulate:a", 0.0, 0.22)
	tween.tween_callback(flash.queue_free)
