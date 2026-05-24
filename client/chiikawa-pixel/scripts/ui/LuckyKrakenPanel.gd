## LuckyKrakenPanel.gd — 幸運深海克拉肯魚 UI 面板（DAY-286）
## 業界依據：Kraken Unleashed「Kraken Reel + 多段觸手攻擊」機制
## 視覺設計：深海主題（#0A1E50 深藍 + #00B4DC 青藍 + #6432B4 紫 + #FFFFFF 白）
##
## Event 類型：
##   - kraken_start：克拉肯召喚（深藍三次強閃光+頂部橫幅+觸手計數器）
##   - kraken_tentacle：觸手攻擊（青藍閃光+浮動文字+計數器更新）
##   - kraken_end：克拉肯結算（結算彈窗）
##   - kraken_fury：克拉肯狂怒（全螢幕三次強閃光+「狂怒！全服×2.8！」大字+狂怒指示器）
##   - kraken_fury_end：狂怒結束（淡出）

extends CanvasLayer

const COLOR_DEEP_BLUE  := Color(0.04, 0.12, 0.31, 1.0)  # #0A1E50 深藍
const COLOR_CYAN_BLUE  := Color(0.0,  0.71, 0.86, 1.0)  # #00B4DC 青藍
const COLOR_PURPLE     := Color(0.39, 0.20, 0.71, 1.0)  # #6432B4 紫
const COLOR_LIGHT_CYAN := Color(0.59, 0.90, 1.0,  1.0)  # 淡青
const COLOR_GRAY       := Color(0.5,  0.5,  0.5,  1.0)  # 灰

var _hit_count: int = 0
var _total_count: int = 8
var _accum_mult: float = 1.0
var _fury_active: bool = false
var _fury_remaining: float = 0.0

var _banner: Control = null
var _indicator: Control = null
var _fury_indicator: Control = null

func _ready() -> void:
	layer = 59
	_build_ui()

func _build_ui() -> void:
	_banner = _make_banner()
	add_child(_banner)
	_banner.visible = false

	_indicator = _make_indicator()
	add_child(_indicator)
	_indicator.visible = false

	_fury_indicator = _make_fury_indicator()
	add_child(_fury_indicator)
	_fury_indicator.visible = false

func _process(delta: float) -> void:
	if _fury_active and _fury_remaining > 0.0:
		_fury_remaining -= delta
		if _fury_remaining < 0.0:
			_fury_remaining = 0.0
		_update_fury_indicator()

func handle_event(data: Dictionary) -> void:
	var event: String = data.get("event", "")
	match event:
		"kraken_start":   _on_kraken_start(data)
		"kraken_tentacle": _on_kraken_tentacle(data)
		"kraken_end":     _on_kraken_end(data)
		"kraken_fury":    _on_kraken_fury(data)
		"kraken_fury_end": _on_kraken_fury_end()

func _on_kraken_start(data: Dictionary) -> void:
	var player_name: String = data.get("player_name", "")
	_hit_count = 0
	_total_count = data.get("tentacle_count", 8)
	_accum_mult = 1.0

	_flash_screen(COLOR_DEEP_BLUE, 3, 0.12)

	_show_banner(
		"🦑🌊 深海克拉肯！",
		"%s 召喚深海克拉肯！%d 條觸手即將攻擊！" % [player_name, _total_count],
		COLOR_CYAN_BLUE
	)

	_indicator.visible = true
	_update_indicator()

	_spawn_float_text("🦑 深海克拉肯降臨！", Vector2(512, 300), COLOR_CYAN_BLUE, 28)

func _on_kraken_tentacle(data: Dictionary) -> void:
	var hit_targets: int = data.get("hit_targets", 0)
	var accum_mult: float = data.get("accum_mult", 1.0)
	var tentacle_idx: int = data.get("tentacle_idx", 1)

	_accum_mult = accum_mult
	if hit_targets > 0:
		_hit_count += 1

	_update_indicator()

	if hit_targets > 0:
		_flash_screen(COLOR_CYAN_BLUE, 1, 0.06)
		_spawn_float_text("🦑 觸手 %d 命中 %d 條！×%.1f" % [tentacle_idx, hit_targets, accum_mult],
			Vector2(512, 320), COLOR_CYAN_BLUE, 22)
	else:
		_spawn_float_text("💨 觸手 %d 空揮..." % tentacle_idx,
			Vector2(512, 320), COLOR_GRAY, 18)

func _on_kraken_end(data: Dictionary) -> void:
	var accum_mult: float = data.get("accum_mult", 1.0)
	var reward: int = data.get("reward", 0)
	var total_hit: int = data.get("total_hit", 0)
	var is_fury: bool = data.get("is_fury", false)

	_indicator.visible = false

	if not is_fury and accum_mult >= 3.0:
		_spawn_float_text("🦑 克拉肯結算！×%.1f 命中 %d 條！+%d 金幣！" % [accum_mult, total_hit, reward],
			Vector2(512, 300), COLOR_CYAN_BLUE, 24)

func _on_kraken_fury(data: Dictionary) -> void:
	var player_name: String = data.get("player_name", "")
	var fury_mult: float = data.get("fury_mult", 2.8)
	var duration: int = data.get("duration", 7)

	_fury_active = true
	_fury_remaining = float(duration)

	_flash_screen(COLOR_DEEP_BLUE, 1, 0.15)
	await get_tree().create_timer(0.2).timeout
	_flash_screen(COLOR_CYAN_BLUE, 1, 0.15)
	await get_tree().create_timer(0.2).timeout
	_flash_screen(COLOR_PURPLE, 1, 0.15)

	_show_full_fury_text(player_name, fury_mult, duration)

	_fury_indicator.visible = true
	_update_fury_indicator()

	_show_banner(
		"🦑🌊🦑 克拉肯狂怒！",
		"%s 克拉肯狂怒！全服 ×%.1f 加成 %d 秒！" % [player_name, fury_mult, duration],
		COLOR_PURPLE
	)

func _on_kraken_fury_end() -> void:
	_fury_active = false
	_fury_remaining = 0.0
	if _fury_indicator.visible:
		var tween := create_tween()
		tween.tween_property(_fury_indicator, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func(): _fury_indicator.visible = false)
	_spawn_float_text("🦑 克拉肯狂怒結束", Vector2(512, 300), COLOR_GRAY, 20)

# ── UI 建構 ───────────────────────────────────────────────────────────────────

func _make_banner() -> Control:
	var panel := PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_WIDE)
	panel.position = Vector2(0, -80)
	panel.size = Vector2(1024, 72)
	var style := StyleBoxFlat.new()
	style.bg_color = Color(0.0, 0.0, 0.0, 0.85)
	style.border_color = COLOR_CYAN_BLUE
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
	title_lbl.add_theme_color_override("font_color", COLOR_CYAN_BLUE)
	vbox.add_child(title_lbl)
	var msg_lbl := Label.new()
	msg_lbl.name = "MsgLabel"
	msg_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	msg_lbl.add_theme_font_size_override("font_size", 14)
	msg_lbl.add_theme_color_override("font_color", COLOR_LIGHT_CYAN)
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
	style.border_color = COLOR_CYAN_BLUE
	style.set_border_width_all(2)
	style.set_corner_radius_all(6)
	panel.add_theme_stylebox_override("panel", style)
	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)
	var title_lbl := Label.new()
	title_lbl.name = "IndicatorTitle"
	title_lbl.text = "🦑 觸手攻擊"
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", COLOR_CYAN_BLUE)
	vbox.add_child(title_lbl)
	var count_lbl := Label.new()
	count_lbl.name = "CountLabel"
	count_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	count_lbl.add_theme_font_size_override("font_size", 18)
	count_lbl.add_theme_color_override("font_color", COLOR_LIGHT_CYAN)
	vbox.add_child(count_lbl)
	var mult_lbl := Label.new()
	mult_lbl.name = "MultLabel"
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mult_lbl.add_theme_font_size_override("font_size", 12)
	mult_lbl.add_theme_color_override("font_color", COLOR_LIGHT_CYAN)
	vbox.add_child(mult_lbl)
	return panel

func _make_fury_indicator() -> Control:
	var panel := PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	panel.position = Vector2(-200, 160)
	panel.size = Vector2(180, 60)
	var style := StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.0, 0.2, 0.9)
	style.border_color = COLOR_PURPLE
	style.set_border_width_all(2)
	style.set_corner_radius_all(6)
	panel.add_theme_stylebox_override("panel", style)
	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)
	var title_lbl := Label.new()
	title_lbl.name = "FuryTitle"
	title_lbl.text = "🌊 克拉肯狂怒 ×2.8"
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", COLOR_PURPLE)
	vbox.add_child(title_lbl)
	var timer_lbl := Label.new()
	timer_lbl.name = "FuryTimerLabel"
	timer_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	timer_lbl.add_theme_font_size_override("font_size", 18)
	timer_lbl.add_theme_color_override("font_color", COLOR_CYAN_BLUE)
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
	count_lbl.text = "%d / %d 命中" % [_hit_count, _total_count]
	mult_lbl.text = "累積 ×%.1f" % _accum_mult
	var progress: float = float(_hit_count) / float(max(_total_count, 1))
	if progress >= 0.75:
		count_lbl.add_theme_color_override("font_color", COLOR_PURPLE)
	elif progress >= 0.5:
		count_lbl.add_theme_color_override("font_color", COLOR_CYAN_BLUE)
	else:
		count_lbl.add_theme_color_override("font_color", COLOR_LIGHT_CYAN)

func _update_fury_indicator() -> void:
	if not _fury_indicator.visible:
		return
	var timer_lbl: Label = _fury_indicator.get_node("VBoxContainer/FuryTimerLabel")
	timer_lbl.text = "%.1f 秒" % _fury_remaining
	var tween := create_tween()
	tween.tween_property(_fury_indicator, "modulate:a", 0.7, 0.35)
	tween.tween_property(_fury_indicator, "modulate:a", 1.0, 0.35)

func _show_full_fury_text(player_name: String, mult: float, duration: int) -> void:
	var lbl := Label.new()
	lbl.text = "🦑🌊🦑 克拉肯狂怒！\n全服 ×%.1f 加成 %d 秒！" % [mult, duration]
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.position = Vector2(512 - 300, 300 - 50)
	lbl.size = Vector2(600, 100)
	lbl.add_theme_font_size_override("font_size", 32)
	lbl.add_theme_color_override("font_color", COLOR_CYAN_BLUE)
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
