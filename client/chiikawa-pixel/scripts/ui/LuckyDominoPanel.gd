## LuckyDominoPanel.gd — 幸運多米諾魚 UI 面板（DAY-288）
## 業界依據：Fishing Fortune 2026「multiplier cascade system」+ Royal Fishing「chain reactions」
## 視覺設計：多米諾主題（#4A0000 深紅棕 + #8B4513 棕 + #D2691E 橙棕 + #FFD700 金 + #FFFFFF 白）
##
## Event 類型：
##   - domino_start：多米諾觸發（棕色三次強閃光+頂部橫幅+骨牌指示器）
##   - domino_knock：骨牌推倒（橙棕閃光+浮動文字+指示器更新）
##   - domino_next：下一個骨牌提示（輕微閃光+箭頭指示）
##   - domino_perfect：多米諾完美（全螢幕三次強閃光+「完美！全服×2.5！」大字+完美指示器）
##   - domino_perfect_end：完美結束（淡出）
##   - domino_end：超時結算（結算彈窗）

extends CanvasLayer

const COLOR_DEEP_BROWN  := Color(0.29, 0.0,  0.0,  1.0)  # #4A0000 深紅棕
const COLOR_BROWN       := Color(0.55, 0.27, 0.07, 1.0)  # #8B4513 棕
const COLOR_ORANGE_BROWN := Color(0.82, 0.41, 0.12, 1.0) # #D2691E 橙棕
const COLOR_GOLD        := Color(1.0,  0.84, 0.0,  1.0)  # #FFD700 金
const COLOR_WHITE       := Color(1.0,  1.0,  1.0,  1.0)  # 白
const COLOR_GRAY        := Color(0.5,  0.5,  0.5,  1.0)  # 灰

var _knocked_count: int = 0
var _total_count: int = 5
var _accum_mult: float = 1.0
var _perfect_active: bool = false
var _perfect_remaining: float = 0.0

var _banner: Control = null
var _indicator: Control = null
var _perfect_indicator: Control = null

func _ready() -> void:
	layer = 61
	_build_ui()

func _build_ui() -> void:
	_banner = _make_banner()
	add_child(_banner)
	_banner.visible = false

	_indicator = _make_indicator()
	add_child(_indicator)
	_indicator.visible = false

	_perfect_indicator = _make_perfect_indicator()
	add_child(_perfect_indicator)
	_perfect_indicator.visible = false

func _process(delta: float) -> void:
	if _perfect_active and _perfect_remaining > 0.0:
		_perfect_remaining -= delta
		if _perfect_remaining < 0.0:
			_perfect_remaining = 0.0
		_update_perfect_indicator()

func handle_event(data: Dictionary) -> void:
	var event: String = data.get("event", "")
	match event:
		"domino_start":       _on_domino_start(data)
		"domino_knock":       _on_domino_knock(data)
		"domino_next":        _on_domino_next(data)
		"domino_perfect":     _on_domino_perfect(data)
		"domino_perfect_end": _on_domino_perfect_end()
		"domino_end":         _on_domino_end(data)

func _on_domino_start(data: Dictionary) -> void:
	var player_name: String = data.get("player_name", "")
	_knocked_count = 0
	_total_count = data.get("total_count", 5)
	_accum_mult = 1.0

	_flash_screen(COLOR_DEEP_BROWN, 1, 0.12)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_BROWN, 1, 0.12)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_ORANGE_BROWN, 1, 0.12)

	_show_banner(
		"🀱🎯 多米諾連鎖！",
		"%s 觸發多米諾！%d 個骨牌等待連鎖推倒！" % [player_name, _total_count],
		COLOR_ORANGE_BROWN
	)

	_indicator.visible = true
	_update_indicator()

	_spawn_float_text("🀱 多米諾骨牌降臨！", Vector2(512, 300), COLOR_ORANGE_BROWN, 28)

func _on_domino_knock(data: Dictionary) -> void:
	var knocked_idx: int = data.get("knocked_idx", 1)
	var accum_mult: float = data.get("accum_mult", 1.0)

	_knocked_count = knocked_idx
	_accum_mult = accum_mult

	_update_indicator()

	# 顏色隨推倒數：棕→橙棕→金
	var knock_color: Color = COLOR_BROWN
	if knocked_idx >= 4:
		knock_color = COLOR_GOLD
	elif knocked_idx >= 3:
		knock_color = COLOR_ORANGE_BROWN

	_flash_screen(knock_color, 1, 0.08)
	_spawn_float_text("🀱 第%d個骨牌推倒！×%.1f" % [knocked_idx, accum_mult],
		Vector2(512, 320), knock_color, 22)

func _on_domino_next(data: Dictionary) -> void:
	var next_idx: int = data.get("next_idx", 1)
	_spawn_float_text("➡ 第%d個骨牌！快去打！" % next_idx,
		Vector2(512, 280), COLOR_ORANGE_BROWN, 18)

func _on_domino_perfect(data: Dictionary) -> void:
	var player_name: String = data.get("player_name", "")
	var perfect_mult: float = data.get("perfect_mult", 2.5)
	var duration: int = data.get("duration", 7)
	var accum_mult: float = data.get("accum_mult", 7.5)

	_perfect_active = true
	_perfect_remaining = float(duration)

	_flash_screen(COLOR_DEEP_BROWN, 1, 0.15)
	await get_tree().create_timer(0.2).timeout
	_flash_screen(COLOR_ORANGE_BROWN, 1, 0.15)
	await get_tree().create_timer(0.2).timeout
	_flash_screen(COLOR_GOLD, 1, 0.15)

	_show_full_perfect_text(player_name, accum_mult, perfect_mult, duration)

	_perfect_indicator.visible = true
	_update_perfect_indicator()

	_show_banner(
		"🀱🎯🀱 多米諾完美！",
		"%s 多米諾完美！累積 ×%.1f！全服 ×%.1f 加成 %d 秒！" % [player_name, accum_mult, perfect_mult, duration],
		COLOR_GOLD
	)

	_indicator.visible = false

func _on_domino_perfect_end() -> void:
	_perfect_active = false
	_perfect_remaining = 0.0
	if _perfect_indicator.visible:
		var tween := create_tween()
		tween.tween_property(_perfect_indicator, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func(): _perfect_indicator.visible = false)
	_spawn_float_text("🀱 多米諾完美結束", Vector2(512, 300), COLOR_GRAY, 20)

func _on_domino_end(data: Dictionary) -> void:
	var knocked_idx: int = data.get("knocked_idx", 0)
	var total_count: int = data.get("total_count", 5)
	var accum_mult: float = data.get("accum_mult", 1.0)

	_indicator.visible = false

	if knocked_idx >= 2:
		_spawn_float_text("🀱 多米諾結算！推倒 %d/%d 個！×%.1f" % [knocked_idx, total_count, accum_mult],
			Vector2(512, 300), COLOR_ORANGE_BROWN, 22)

# ── UI 建構 ───────────────────────────────────────────────────────────────────

func _make_banner() -> Control:
	var panel := PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_WIDE)
	panel.position = Vector2(0, -80)
	panel.size = Vector2(1024, 72)
	var style := StyleBoxFlat.new()
	style.bg_color = Color(0.0, 0.0, 0.0, 0.85)
	style.border_color = COLOR_ORANGE_BROWN
	style.set_border_width_all(2)
	style.set_corner_radius_all(6)
	panel.add_theme_stylebox_override("panel", style)
	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)
	var title_lbl := Label.new()
	title_lbl.name = "TitleLabel"
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", COLOR_ORANGE_BROWN)
	vbox.add_child(title_lbl)
	var msg_lbl := Label.new()
	msg_lbl.name = "MsgLabel"
	msg_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	msg_lbl.add_theme_font_size_override("font_size", 14)
	msg_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	msg_lbl.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	vbox.add_child(msg_lbl)
	return panel

func _make_indicator() -> Control:
	var panel := PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	panel.position = Vector2(-200, 80)
	panel.size = Vector2(180, 70)
	var style := StyleBoxFlat.new()
	style.bg_color = Color(0.0, 0.0, 0.0, 0.8)
	style.border_color = COLOR_ORANGE_BROWN
	style.set_border_width_all(2)
	style.set_corner_radius_all(6)
	panel.add_theme_stylebox_override("panel", style)
	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)
	var title_lbl := Label.new()
	title_lbl.name = "IndicatorTitle"
	title_lbl.text = "🀱 多米諾骨牌"
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", COLOR_ORANGE_BROWN)
	vbox.add_child(title_lbl)
	var count_lbl := Label.new()
	count_lbl.name = "CountLabel"
	count_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	count_lbl.add_theme_font_size_override("font_size", 18)
	count_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	vbox.add_child(count_lbl)
	var mult_lbl := Label.new()
	mult_lbl.name = "MultLabel"
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mult_lbl.add_theme_font_size_override("font_size", 12)
	mult_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	vbox.add_child(mult_lbl)
	return panel

func _make_perfect_indicator() -> Control:
	var panel := PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	panel.position = Vector2(-200, 160)
	panel.size = Vector2(180, 60)
	var style := StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.05, 0.0, 0.9)
	style.border_color = COLOR_GOLD
	style.set_border_width_all(2)
	style.set_corner_radius_all(6)
	panel.add_theme_stylebox_override("panel", style)
	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)
	var title_lbl := Label.new()
	title_lbl.name = "PerfectTitle"
	title_lbl.text = "🎯 多米諾完美 ×2.5"
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	vbox.add_child(title_lbl)
	var timer_lbl := Label.new()
	timer_lbl.name = "PerfectTimerLabel"
	timer_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	timer_lbl.add_theme_font_size_override("font_size", 18)
	timer_lbl.add_theme_color_override("font_color", COLOR_ORANGE_BROWN)
	vbox.add_child(timer_lbl)
	return panel

# ── UI 更新 ───────────────────────────────────────────────────────────────────

func _show_banner(title: String, msg: String, color: Color) -> void:
	var title_lbl: Label = _banner.get_node("VBoxContainer/TitleLabel")
	var msg_lbl: Label = _banner.get_node("VBoxContainer/MsgLabel")
	title_lbl.text = title
	title_lbl.add_theme_color_override("font_color", color)
	msg_lbl.text = msg
	_banner.visible = true
	_banner.modulate.a = 1.0
	_banner.position.y = -80
	var tween := create_tween()
	tween.tween_property(_banner, "position:y", 0.0, 0.3).set_ease(Tween.EASE_OUT)
	tween.tween_interval(4.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func(): _banner.visible = false)

func _update_indicator() -> void:
	if not _indicator.visible:
		return
	var count_lbl: Label = _indicator.get_node("VBoxContainer/CountLabel")
	var mult_lbl: Label = _indicator.get_node("VBoxContainer/MultLabel")
	count_lbl.text = "%d / %d 推倒" % [_knocked_count, _total_count]
	mult_lbl.text = "累積 ×%.1f" % _accum_mult
	if _knocked_count >= 4:
		count_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	elif _knocked_count >= 3:
		count_lbl.add_theme_color_override("font_color", COLOR_ORANGE_BROWN)
	else:
		count_lbl.add_theme_color_override("font_color", COLOR_WHITE)

func _update_perfect_indicator() -> void:
	if not _perfect_indicator.visible:
		return
	var timer_lbl: Label = _perfect_indicator.get_node("VBoxContainer/PerfectTimerLabel")
	timer_lbl.text = "%.1f 秒" % _perfect_remaining
	var tween := create_tween()
	tween.tween_property(_perfect_indicator, "modulate:a", 0.7, 0.35)
	tween.tween_property(_perfect_indicator, "modulate:a", 1.0, 0.35)

func _show_full_perfect_text(player_name: String, accum_mult: float, perfect_mult: float, duration: int) -> void:
	var lbl := Label.new()
	lbl.text = "🀱🎯🀱 多米諾完美！\n累積 ×%.1f！全服 ×%.1f 加成 %d 秒！" % [accum_mult, perfect_mult, duration]
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.position = Vector2(512 - 300, 300 - 50)
	lbl.size = Vector2(600, 100)
	lbl.add_theme_font_size_override("font_size", 30)
	lbl.add_theme_color_override("font_color", COLOR_GOLD)
	add_child(lbl)
	lbl.scale = Vector2(0.5, 0.5)
	var tween := create_tween()
	tween.tween_property(lbl, "scale", Vector2(1.1, 1.1), 0.3).set_ease(Tween.EASE_OUT)
	tween.tween_property(lbl, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(2.5)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.5)
	tween.tween_callback(lbl.queue_free)

func _flash_screen(color: Color, times: int, duration: float) -> void:
	var overlay := ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, 0.0)
	overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tween := create_tween()
	for i in range(times):
		tween.tween_property(overlay, "color:a", 0.4, duration * 0.4)
		tween.tween_property(overlay, "color:a", 0.0, duration * 0.6)
	tween.tween_callback(overlay.queue_free)

func _spawn_float_text(text: String, pos: Vector2, color: Color, font_size: int = 22) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.position = pos
	lbl.add_theme_font_size_override("font_size", font_size)
	lbl.add_theme_color_override("font_color", color)
	add_child(lbl)
	var tween := create_tween()
	tween.tween_property(lbl, "position:y", pos.y - 60, 1.2).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 1.2)
	tween.tween_callback(lbl.queue_free)
