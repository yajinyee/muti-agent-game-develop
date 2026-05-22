## CometFishPanel.gd — 彗星魚連鎖爆炸面板（DAY-206）
## 業界依據：Ocean King 3 Plus「Comet Fish — streaks across the screen leaving a trail of explosions」
##
## 視覺設計：
##   - 彗星橙白主題（#FF8C00 + #FFD700 + #FFFFFF + #FF4500）
##   - comet_appear：橙色閃光 + 頂部橫幅 + 彗星軌跡預覽線
##   - trail_blast：爆炸圓圈擴散 + 橙色閃光 + 擊破浮動文字
##   - early_supernova：金色強閃光 + 「提前引爆！」大字
##   - supernova：全螢幕橙紅三次強閃光 + 「☄️ 超新星爆炸！」大字 + 結算彈窗
extends CanvasLayer

func _ready() -> void:
	layer = 39  # 彗星魚面板層級

## 處理彗星魚訊息
func handle_comet_fish(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"comet_appear":
			_on_comet_appear(payload)
		"trail_blast":
			_on_trail_blast(payload)
		"early_supernova":
			_on_early_supernova(payload)
		"supernova":
			_on_supernova(payload)

## 彗星出現 — 橙色閃光 + 頂部橫幅
func _on_comet_appear(payload: Dictionary) -> void:
	var trail_count = payload.get("trail_count", 7)

	# 橙色閃光
	_flash_screen(Color("#FF8C00"), 0.5)

	# 頂部橫幅
	var banner = _make_banner(
		"☄️ 彗星魚出現！沿途 %d 次爆炸！" % trail_count,
		Color("#FF8C00"),
		Color("#1A0A00")
	)
	add_child(banner)
	var tween = create_tween()
	tween.tween_property(banner, "position:y", 10.0, 0.3)
	tween.tween_interval(2.5)
	tween.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(banner.queue_free)

## 軌跡爆炸 — 爆炸圓圈 + 浮動文字
func _on_trail_blast(payload: Dictionary) -> void:
	var x = payload.get("x", 640.0)
	var y = payload.get("y", 360.0)
	var kills = payload.get("kill_count", 0)
	var reward = payload.get("reward", 0)
	var blast_idx = payload.get("blast_index", 1)

	# 爆炸圓圈（橙色，擴散）
	var circle = ColorRect.new()
	circle.color = Color("#FF8C00")
	circle.size = Vector2(20, 20)
	circle.position = Vector2(x - 10, y - 10)
	circle.modulate.a = 0.7
	add_child(circle)
	var tween = create_tween()
	tween.tween_property(circle, "size", Vector2(400, 400), 0.4)
	tween.parallel().tween_property(circle, "position", Vector2(x - 200, y - 200), 0.4)
	tween.parallel().tween_property(circle, "modulate:a", 0.0, 0.4)
	tween.tween_callback(circle.queue_free)

	# 擊破浮動文字
	if kills > 0:
		var float_label = Label.new()
		float_label.text = "💥 ×%d  +%d" % [kills, reward]
		float_label.add_theme_font_size_override("font_size", 18)
		float_label.add_theme_color_override("font_color", Color("#FFD700"))
		float_label.position = Vector2(x - 40, y - 20)
		add_child(float_label)
		var tween2 = create_tween()
		tween2.tween_property(float_label, "position:y", y - 60, 0.8)
		tween2.parallel().tween_property(float_label, "modulate:a", 0.0, 0.8)
		tween2.tween_callback(float_label.queue_free)

	# 第 7 個爆炸點（最後一個）額外閃光
	if blast_idx >= 7:
		_flash_screen(Color("#FFD700"), 0.3)

## 提前引爆 — 金色強閃光 + 大字
func _on_early_supernova(payload: Dictionary) -> void:
	var player_name = payload.get("player_name", "玩家")

	# 金色強閃光
	_flash_screen(Color("#FFD700"), 0.7)

	# 「提前引爆！」大字
	var big_label = Label.new()
	big_label.text = "⚡ %s 提前引爆！" % player_name
	big_label.add_theme_font_size_override("font_size", 36)
	big_label.add_theme_color_override("font_color", Color("#FFD700"))
	big_label.set_anchors_preset(Control.PRESET_CENTER)
	big_label.position = Vector2(-200, -40)
	add_child(big_label)
	var tween = create_tween()
	tween.tween_property(big_label, "scale", Vector2(1.15, 1.15), 0.2)
	tween.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.15)
	tween.tween_interval(0.8)
	tween.tween_property(big_label, "modulate:a", 0.0, 0.5)
	tween.tween_callback(big_label.queue_free)

## 超新星爆炸 — 全螢幕三次強閃光 + 大字 + 結算彈窗
func _on_supernova(payload: Dictionary) -> void:
	var kills = payload.get("kill_count", 0)
	var reward = payload.get("reward", 0)
	var is_early = payload.get("is_early", false)
	var player_name = payload.get("player_name", "")

	# 全螢幕三次強閃光（橙→金→白）
	_flash_screen(Color("#FF4500"), 0.85)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(Color("#FFD700"), 0.75)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(Color("#FFFFFF"), 0.65)

	# 「☄️ 超新星爆炸！」大字
	var title_text = "☄️ 超新星爆炸！"
	if is_early and player_name != "":
		title_text = "☄️ %s 引爆超新星！" % player_name
	var big_label = Label.new()
	big_label.text = title_text
	big_label.add_theme_font_size_override("font_size", 44)
	big_label.add_theme_color_override("font_color", Color("#FF8C00"))
	big_label.set_anchors_preset(Control.PRESET_CENTER)
	big_label.position = Vector2(-280, -80)
	add_child(big_label)
	var tween = create_tween()
	tween.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.25)
	tween.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.2)
	tween.tween_interval(1.2)
	tween.tween_property(big_label, "modulate:a", 0.0, 0.6)
	tween.tween_callback(big_label.queue_free)

	# 結算彈窗（右側滑入）
	if kills > 0:
		await get_tree().create_timer(0.4).timeout
		_show_result_popup(kills, reward, is_early)

## 結算彈窗
func _show_result_popup(kills: int, reward: int, is_early: bool) -> void:
	var popup = PanelContainer.new()
	popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	popup.position = Vector2(1400, -80)
	popup.size = Vector2(260, 160)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.05, 0.0, 0.92)
	style.border_color = Color("#FF8C00")
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
	title_lbl.text = "☄️ 超新星結算"
	title_lbl.add_theme_font_size_override("font_size", 16)
	title_lbl.add_theme_color_override("font_color", Color("#FF8C00"))
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var kills_lbl = Label.new()
	kills_lbl.text = "擊破目標：%d 個" % kills
	kills_lbl.add_theme_font_size_override("font_size", 14)
	kills_lbl.add_theme_color_override("font_color", Color("#FFD700"))
	kills_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(kills_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "全服獎勵：%d 金幣" % reward
	reward_lbl.add_theme_font_size_override("font_size", 14)
	reward_lbl.add_theme_color_override("font_color", Color("#FFFFFF"))
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	if is_early:
		var early_lbl = Label.new()
		early_lbl.text = "⚡ 提前引爆加成！"
		early_lbl.add_theme_font_size_override("font_size", 12)
		early_lbl.add_theme_color_override("font_color", Color("#00FFFF"))
		early_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		vbox.add_child(early_lbl)

	add_child(popup)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", 1000.0, 0.35)
	tween.tween_interval(3.0)
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
