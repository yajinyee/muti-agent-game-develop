## LuckyCosmicPulsePanel.gd — 幸運宇宙脈衝魚 UI 面板（DAY-287）
## 業界依據：TaDa Gaming 2026「Cosmic」主題 + Fishing Fortune 2026「pulse wave mechanics」
## 視覺設計：宇宙主題（#0D0025 深紫黑 + #7B2FBE 紫 + #C77DFF 淡紫 + #E0AAFF 淡紫白 + #FFD700 金）
##
## Event 類型：
##   - pulse_start：宇宙脈衝開始（紫色三次強閃光+頂部橫幅+脈衝指示器）
##   - pulse_wave：脈衝波結果（顏色隨層數閃光+浮動文字+指示器更新）
##   - pulse_end：宇宙脈衝結算（結算彈窗）
##   - pulse_resonance：宇宙共振（全螢幕三次強閃光+「宇宙共振！全服×2.2！」大字+共振指示器）
##   - pulse_resonance_end：共振結束（淡出）

extends CanvasLayer

const COLOR_DEEP_SPACE  := Color(0.05, 0.0,  0.15, 1.0)  # #0D0025 深紫黑
const COLOR_PURPLE      := Color(0.48, 0.18, 0.75, 1.0)  # #7B2FBE 紫
const COLOR_LIGHT_PURPLE := Color(0.78, 0.49, 1.0,  1.0) # #C77DFF 淡紫
const COLOR_PALE_PURPLE := Color(0.88, 0.67, 1.0,  1.0)  # #E0AAFF 淡紫白
const COLOR_GOLD        := Color(1.0,  0.84, 0.0,  1.0)  # #FFD700 金
const COLOR_GRAY        := Color(0.5,  0.5,  0.5,  1.0)  # 灰

var _wave_idx: int = 0
var _total_waves: int = 3
var _accum_mult: float = 1.0
var _total_hit: int = 0
var _resonance_active: bool = false
var _resonance_remaining: float = 0.0

var _banner: Control = null
var _indicator: Control = null
var _resonance_indicator: Control = null

func _ready() -> void:
	layer = 60
	_build_ui()

func _build_ui() -> void:
	_banner = _make_banner()
	add_child(_banner)
	_banner.visible = false

	_indicator = _make_indicator()
	add_child(_indicator)
	_indicator.visible = false

	_resonance_indicator = _make_resonance_indicator()
	add_child(_resonance_indicator)
	_resonance_indicator.visible = false

func _process(delta: float) -> void:
	if _resonance_active and _resonance_remaining > 0.0:
		_resonance_remaining -= delta
		if _resonance_remaining < 0.0:
			_resonance_remaining = 0.0
		_update_resonance_indicator()

func handle_event(data: Dictionary) -> void:
	var event: String = data.get("event", "")
	match event:
		"pulse_start":          _on_pulse_start(data)
		"pulse_wave":           _on_pulse_wave(data)
		"pulse_end":            _on_pulse_end(data)
		"pulse_resonance":      _on_pulse_resonance(data)
		"pulse_resonance_end":  _on_pulse_resonance_end()

func _on_pulse_start(data: Dictionary) -> void:
	var player_name: String = data.get("player_name", "")
	_wave_idx = 0
	_total_waves = data.get("wave_count", 3)
	_accum_mult = 1.0
	_total_hit = 0

	_flash_screen(COLOR_DEEP_SPACE, 1, 0.12)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_PURPLE, 1, 0.12)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_LIGHT_PURPLE, 1, 0.12)

	_show_banner(
		"🌌✨ 宇宙脈衝！",
		"%s 觸發宇宙脈衝！%d 層脈衝波即將擴散！" % [player_name, _total_waves],
		COLOR_LIGHT_PURPLE
	)

	_indicator.visible = true
	_update_indicator()

	_spawn_float_text("🌌 宇宙脈衝降臨！", Vector2(512, 300), COLOR_LIGHT_PURPLE, 28)

func _on_pulse_wave(data: Dictionary) -> void:
	var wave_idx: int = data.get("wave_idx", 1)
	var hit_targets: int = data.get("hit_targets", 0)
	var accum_mult: float = data.get("accum_mult", 1.0)
	var radius: float = data.get("radius", 200.0)

	_wave_idx = wave_idx
	_accum_mult = accum_mult
	_total_hit += hit_targets

	_update_indicator()

	# 顏色隨層數變化：第1層淡紫→第2層紫→第3層金
	var wave_color: Color = COLOR_PALE_PURPLE
	if wave_idx == 2:
		wave_color = COLOR_LIGHT_PURPLE
	elif wave_idx == 3:
		wave_color = COLOR_GOLD

	_flash_screen(wave_color, 1, 0.08)

	if hit_targets > 0:
		_spawn_float_text("🌌 第%d層脈衝！命中 %d 條！×%.1f" % [wave_idx, hit_targets, accum_mult],
			Vector2(512, 320), wave_color, 22)
	else:
		_spawn_float_text("💨 第%d層脈衝... 未命中" % wave_idx,
			Vector2(512, 320), COLOR_GRAY, 18)

func _on_pulse_end(data: Dictionary) -> void:
	var accum_mult: float = data.get("accum_mult", 1.0)
	var reward: int = data.get("reward", 0)
	var total_hit: int = data.get("total_hit", 0)
	var is_resonance: bool = data.get("is_resonance", false)

	_indicator.visible = false

	if not is_resonance and accum_mult >= 3.0:
		_spawn_float_text("🌌 脈衝結算！×%.1f 命中 %d 條！+%d 金幣！" % [accum_mult, total_hit, reward],
			Vector2(512, 300), COLOR_LIGHT_PURPLE, 24)

func _on_pulse_resonance(data: Dictionary) -> void:
	var player_name: String = data.get("player_name", "")
	var res_mult: float = data.get("res_mult", 2.2)
	var duration: int = data.get("duration", 6)
	var total_hit: int = data.get("total_hit", 0)

	_resonance_active = true
	_resonance_remaining = float(duration)

	_flash_screen(COLOR_DEEP_SPACE, 1, 0.15)
	await get_tree().create_timer(0.2).timeout
	_flash_screen(COLOR_PURPLE, 1, 0.15)
	await get_tree().create_timer(0.2).timeout
	_flash_screen(COLOR_GOLD, 1, 0.15)

	_show_full_resonance_text(player_name, res_mult, duration, total_hit)

	_resonance_indicator.visible = true
	_update_resonance_indicator()

	_show_banner(
		"🌌🌟🌌 宇宙共振！",
		"%s 宇宙共振！命中 %d 條！全服 ×%.1f 加成 %d 秒！" % [player_name, total_hit, res_mult, duration],
		COLOR_GOLD
	)

func _on_pulse_resonance_end() -> void:
	_resonance_active = false
	_resonance_remaining = 0.0
	if _resonance_indicator.visible:
		var tween := create_tween()
		tween.tween_property(_resonance_indicator, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func(): _resonance_indicator.visible = false)
	_spawn_float_text("🌌 宇宙共振結束", Vector2(512, 300), COLOR_GRAY, 20)

# ── UI 建構 ───────────────────────────────────────────────────────────────────

func _make_banner() -> Control:
	var panel := PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_WIDE)
	panel.position = Vector2(0, -80)
	panel.size = Vector2(1024, 72)
	var style := StyleBoxFlat.new()
	style.bg_color = Color(0.0, 0.0, 0.0, 0.85)
	style.border_color = COLOR_LIGHT_PURPLE
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
	title_lbl.add_theme_color_override("font_color", COLOR_LIGHT_PURPLE)
	vbox.add_child(title_lbl)
	var msg_lbl := Label.new()
	msg_lbl.name = "MsgLabel"
	msg_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	msg_lbl.add_theme_font_size_override("font_size", 14)
	msg_lbl.add_theme_color_override("font_color", COLOR_PALE_PURPLE)
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
	style.border_color = COLOR_LIGHT_PURPLE
	style.set_border_width_all(2)
	style.set_corner_radius_all(6)
	panel.add_theme_stylebox_override("panel", style)
	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)
	var title_lbl := Label.new()
	title_lbl.name = "IndicatorTitle"
	title_lbl.text = "🌌 脈衝波"
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", COLOR_LIGHT_PURPLE)
	vbox.add_child(title_lbl)
	var wave_lbl := Label.new()
	wave_lbl.name = "WaveLabel"
	wave_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	wave_lbl.add_theme_font_size_override("font_size", 18)
	wave_lbl.add_theme_color_override("font_color", COLOR_PALE_PURPLE)
	vbox.add_child(wave_lbl)
	var mult_lbl := Label.new()
	mult_lbl.name = "MultLabel"
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mult_lbl.add_theme_font_size_override("font_size", 12)
	mult_lbl.add_theme_color_override("font_color", COLOR_PALE_PURPLE)
	vbox.add_child(mult_lbl)
	return panel

func _make_resonance_indicator() -> Control:
	var panel := PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	panel.position = Vector2(-200, 160)
	panel.size = Vector2(180, 60)
	var style := StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.0, 0.1, 0.9)
	style.border_color = COLOR_GOLD
	style.set_border_width_all(2)
	style.set_corner_radius_all(6)
	panel.add_theme_stylebox_override("panel", style)
	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)
	var title_lbl := Label.new()
	title_lbl.name = "ResTitle"
	title_lbl.text = "🌟 宇宙共振 ×2.2"
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	vbox.add_child(title_lbl)
	var timer_lbl := Label.new()
	timer_lbl.name = "ResTimerLabel"
	timer_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	timer_lbl.add_theme_font_size_override("font_size", 18)
	timer_lbl.add_theme_color_override("font_color", COLOR_LIGHT_PURPLE)
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
	var wave_lbl: Label = _indicator.get_node("VBoxContainer/WaveLabel")
	var mult_lbl: Label = _indicator.get_node("VBoxContainer/MultLabel")
	wave_lbl.text = "第 %d / %d 層" % [_wave_idx, _total_waves]
	mult_lbl.text = "累積 ×%.1f（命中 %d）" % [_accum_mult, _total_hit]
	# 顏色隨層數：淡紫→紫→金
	if _wave_idx >= 3:
		wave_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	elif _wave_idx >= 2:
		wave_lbl.add_theme_color_override("font_color", COLOR_LIGHT_PURPLE)
	else:
		wave_lbl.add_theme_color_override("font_color", COLOR_PALE_PURPLE)

func _update_resonance_indicator() -> void:
	if not _resonance_indicator.visible:
		return
	var timer_lbl: Label = _resonance_indicator.get_node("VBoxContainer/ResTimerLabel")
	timer_lbl.text = "%.1f 秒" % _resonance_remaining
	var tween := create_tween()
	tween.tween_property(_resonance_indicator, "modulate:a", 0.7, 0.35)
	tween.tween_property(_resonance_indicator, "modulate:a", 1.0, 0.35)

func _show_full_resonance_text(player_name: String, mult: float, duration: int, total_hit: int) -> void:
	var lbl := Label.new()
	lbl.text = "🌌🌟🌌 宇宙共振！\n命中 %d 條！全服 ×%.1f 加成 %d 秒！" % [total_hit, mult, duration]
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
