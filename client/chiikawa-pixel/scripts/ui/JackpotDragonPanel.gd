## JackpotDragonPanel.gd — 獎池龍 Jackpot 抽獎面板（DAY-205）
## 業界依據：JILI Jackpot Fishing
## 「special targets like the Jackpot Fish and Jackpot Dragon offering chances at substantial prizes.
##  With the potential for high payouts up to 1000 times the bet.」
##
## 視覺設計：
##   - 金龍主題（#FFD700 + #FF6B35 + #FF0080 + #FFFFFF）
##   - dragon_draw（Mini/Minor）：金色閃光 + 橫幅 + 等級彈窗 + 獎勵文字
##   - dragon_draw（Major）：橙紅色雙閃光 + 大字 + 結算彈窗
##   - dragon_draw（Grand）：全螢幕粉紅金色爆炸 + 「GRAND JACKPOT！」大字 + 結算彈窗
extends CanvasLayer

func _ready() -> void:
	layer = 60  # 高於其他面板，Jackpot 最重要

## 處理獎池龍訊息
func handle_jackpot_dragon(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	if event == "dragon_draw":
		_on_dragon_draw(payload)

## 獎池龍抽獎結果
func _on_dragon_draw(payload: Dictionary) -> void:
	var player_name = payload.get("player_name", "玩家")
	var level = payload.get("level", "mini")
	var level_name = payload.get("level_name", "MINI")
	var level_color = payload.get("level_color", "#C0C0C0")
	var level_icon = payload.get("level_icon", "🥈")
	var amount = payload.get("amount", 0)
	var is_grand = payload.get("is_grand", false)
	var is_major = payload.get("is_major", false)

	match level:
		"grand":
			_show_grand_jackpot(player_name, level_name, level_icon, amount)
		"major":
			_show_major_jackpot(player_name, level_name, level_icon, amount)
		"minor":
			_show_minor_jackpot(player_name, level_name, level_icon, amount)
		_:  # mini
			_show_mini_jackpot(player_name, level_name, level_icon, amount)

## Grand Jackpot — 全螢幕爆炸
func _show_grand_jackpot(player_name: String, level_name: String, level_icon: String, amount: int) -> void:
	# 全螢幕粉紅色三次強閃光
	_flash_screen(Color("#FF0080"), 0.9)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color("#FFD700"), 0.8)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color("#FFFFFF"), 0.7)

	# 「GRAND JACKPOT！」大字
	var big_label = Label.new()
	big_label.text = "%s GRAND JACKPOT！" % level_icon
	big_label.add_theme_font_size_override("font_size", 52)
	big_label.add_theme_color_override("font_color", Color("#FF0080"))
	big_label.set_anchors_preset(Control.PRESET_CENTER)
	big_label.position = Vector2(-300, -60)
	add_child(big_label)

	var tween = create_tween()
	tween.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.3)
	tween.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.2)
	tween.tween_interval(1.5)
	tween.tween_property(big_label, "modulate:a", 0.0, 0.8)
	tween.tween_callback(big_label.queue_free)

	# 結算彈窗
	_show_result_popup(player_name, level_name, level_icon, amount, Color("#FF0080"))

## Major Jackpot — 橙紅色雙閃光
func _show_major_jackpot(player_name: String, level_name: String, level_icon: String, amount: int) -> void:
	_flash_screen(Color("#FF6B35"), 0.7)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color("#FFD700"), 0.5)

	# 「MAJOR JACKPOT！」大字
	var big_label = Label.new()
	big_label.text = "%s MAJOR JACKPOT！" % level_icon
	big_label.add_theme_font_size_override("font_size", 40)
	big_label.add_theme_color_override("font_color", Color("#FF6B35"))
	big_label.set_anchors_preset(Control.PRESET_CENTER)
	big_label.position = Vector2(-250, -50)
	add_child(big_label)

	var tween = create_tween()
	tween.tween_property(big_label, "scale", Vector2(1.1, 1.1), 0.25)
	tween.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.15)
	tween.tween_interval(1.2)
	tween.tween_property(big_label, "modulate:a", 0.0, 0.6)
	tween.tween_callback(big_label.queue_free)

	_show_result_popup(player_name, level_name, level_icon, amount, Color("#FF6B35"))

## Minor Jackpot — 金色閃光
func _show_minor_jackpot(player_name: String, level_name: String, level_icon: String, amount: int) -> void:
	_flash_screen(Color("#FFD700"), 0.5)

	_spawn_float_text("%s %s JACKPOT +%d" % [level_icon, level_name, amount],
		Color("#FFD700"), Vector2(640, 280))

	_show_result_popup(player_name, level_name, level_icon, amount, Color("#FFD700"))

## Mini Jackpot — 銀色小提示
func _show_mini_jackpot(player_name: String, level_name: String, level_icon: String, amount: int) -> void:
	_spawn_float_text("%s %s +%d" % [level_icon, level_name, amount],
		Color("#C0C0C0"), Vector2(randf_range(400, 800), randf_range(200, 450)))

## 結算彈窗
func _show_result_popup(player_name: String, level_name: String, level_icon: String, amount: int, color: Color) -> void:
	var popup = Control.new()
	popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	popup.size = Vector2(280, 160)
	popup.position = Vector2(get_viewport().get_visible_rect().size.x + 10, -80)
	add_child(popup)

	var bg = ColorRect.new()
	bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.05, 0.02, 0.08, 0.93)
	popup.add_child(bg)

	var label = Label.new()
	label.text = (
		"%s 獎池龍抽獎\n\n"
		% level_icon
		+ "%s: %s\n" % [player_name, level_name]
		+ "獎勵: %d 金幣" % amount
	)
	label.add_theme_font_size_override("font_size", 15)
	label.add_theme_color_override("font_color", color)
	label.set_anchors_preset(Control.PRESET_FULL_RECT)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	popup.add_child(label)

	# 從右側滑入
	var tween = create_tween()
	tween.tween_property(popup, "position:x",
		get_viewport().get_visible_rect().size.x - 300, 0.4)
	tween.tween_interval(3.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)

## 全螢幕閃光
func _flash_screen(color: Color, alpha: float) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, alpha)
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.3)
	tween.tween_callback(flash.queue_free)

## 浮動文字
func _spawn_float_text(text: String, color: Color, pos: Vector2) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 20)
	label.add_theme_color_override("font_color", color)
	label.position = pos
	add_child(label)
	var tween = create_tween()
	tween.tween_property(label, "position:y", label.position.y - 70, 1.2)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 1.2)
	tween.tween_callback(label.queue_free)
