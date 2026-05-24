## LuckyPhoenixRebirthPanel.gd — 幸運鳳凰涅槃魚 UI 面板（DAY-285）
## 業界依據：Royal Fishing Jili「Rainbow Phoenix Power Up」機制
## 視覺設計：鳳凰主題（#FF6B35 火橙 + #FFD700 金 + #FF4500 深橙 + #FFF0E0 淡橙白）
##
## Event 類型：
##   - rebirth_start：鳳凰涅槃開始（火橙三次強閃光+頂部橫幅+涅槃指示器+計時條）
##   - rebirth_kill：涅槃目標被擊破（金色閃光+浮動文字+指示器更新）
##   - rebirth_full：鳳凰完全涅槃（全螢幕三次強閃光+「完全涅槃！全服×3.0！」大字+完全涅槃指示器）
##   - rebirth_full_end：完全涅槃結束（淡出）
##   - rebirth_fade：涅槃消散（灰色閃光+消散文字）

extends CanvasLayer

# ── 顏色常數 ──────────────────────────────────────────────────────────────────
const COLOR_FIRE_ORANGE  := Color(1.0, 0.42, 0.21, 1.0)   # #FF6B35 火橙
const COLOR_GOLD         := Color(1.0, 0.84, 0.0,  1.0)   # #FFD700 金
const COLOR_DEEP_ORANGE  := Color(1.0, 0.27, 0.0,  1.0)   # #FF4500 深橙
const COLOR_LIGHT_ORANGE := Color(1.0, 0.94, 0.88, 1.0)   # #FFF0E0 淡橙白
const COLOR_GRAY         := Color(0.5, 0.5, 0.5, 1.0)     # 消散灰

# ── 狀態 ──────────────────────────────────────────────────────────────────────
var _rebirth_targets: Array = []   # [{instance_id, name, x, y}]
var _killed_count: int = 0
var _total_count: int = 0
var _duration: int = 15
var _timer_remaining: float = 0.0
var _is_active: bool = false
var _full_rebirth_active: bool = false
var _full_rebirth_remaining: float = 0.0

# ── 節點引用 ──────────────────────────────────────────────────────────────────
var _banner: Control = null
var _indicator: Control = null
var _timer_bar: Control = null
var _full_rebirth_indicator: Control = null

func _ready() -> void:
	layer = 58  # 比 LuckyDragonWrathPanel（layer=57）高一層
	_build_ui()

func _build_ui() -> void:
	# 頂部橫幅
	_banner = _make_banner()
	add_child(_banner)
	_banner.visible = false

	# 右上角涅槃指示器
	_indicator = _make_indicator()
	add_child(_indicator)
	_indicator.visible = false

	# 右側豎向計時條
	_timer_bar = _make_timer_bar()
	add_child(_timer_bar)
	_timer_bar.visible = false

	# 完全涅槃指示器（右上角，計時條下方）
	_full_rebirth_indicator = _make_full_rebirth_indicator()
	add_child(_full_rebirth_indicator)
	_full_rebirth_indicator.visible = false

func _process(delta: float) -> void:
	# 更新涅槃計時條
	if _is_active and _timer_remaining > 0.0:
		_timer_remaining -= delta
		if _timer_remaining < 0.0:
			_timer_remaining = 0.0
		_update_timer_bar()

	# 更新完全涅槃倒數
	if _full_rebirth_active and _full_rebirth_remaining > 0.0:
		_full_rebirth_remaining -= delta
		if _full_rebirth_remaining < 0.0:
			_full_rebirth_remaining = 0.0
		_update_full_rebirth_indicator()

# ── 主要事件處理 ──────────────────────────────────────────────────────────────

func handle_event(data: Dictionary) -> void:
	var event: String = data.get("event", "")
	match event:
		"rebirth_start":
			_on_rebirth_start(data)
		"rebirth_kill":
			_on_rebirth_kill(data)
		"rebirth_full":
			_on_rebirth_full(data)
		"rebirth_full_end":
			_on_rebirth_full_end()
		"rebirth_fade":
			_on_rebirth_fade(data)

func _on_rebirth_start(data: Dictionary) -> void:
	var player_name: String = data.get("player_name", "")
	var fire_bath_count: int = data.get("fire_bath_count", 0)
	var targets: Array = data.get("rebirth_targets", [])
	_duration = data.get("duration", 15)

	_rebirth_targets = targets
	_killed_count = 0
	_total_count = targets.size()
	_timer_remaining = float(_duration)
	_is_active = true

	# 火橙三次強閃光
	_flash_screen(COLOR_FIRE_ORANGE, 3, 0.12)

	# 顯示頂部橫幅
	_show_banner(
		"🔥🦅 鳳凰涅槃！",
		"%s 召喚鳳凰涅槃！全場火焰洗禮（%d 條），%d 條魚涅槃重生！" % [player_name, fire_bath_count, _total_count],
		COLOR_FIRE_ORANGE
	)

	# 顯示涅槃指示器
	_indicator.visible = true
	_update_indicator()

	# 顯示計時條
	_timer_bar.visible = true
	_update_timer_bar()

	# 浮動文字
	_spawn_float_text("🔥 全場火焰洗禮！%d 條魚受傷！" % fire_bath_count,
		Vector2(512, 300), COLOR_FIRE_ORANGE, 28)
	await get_tree().create_timer(0.8).timeout
	_spawn_float_text("🦅 %d 條涅槃魚出現！擊破得 ×%.0f！" % [_total_count, 4.0],
		Vector2(512, 340), COLOR_GOLD, 24)

func _on_rebirth_kill(data: Dictionary) -> void:
	var player_name: String = data.get("player_name", "")
	var kill_mult: float = data.get("kill_mult", 4.0)

	_killed_count += 1

	# 金色閃光
	_flash_screen(COLOR_GOLD, 1, 0.08)

	# 更新指示器
	_update_indicator()

	# 浮動文字
	_spawn_float_text("🦅 涅槃擊破！×%.0f！（%d/%d）" % [kill_mult, _killed_count, _total_count],
		Vector2(512, 320), COLOR_GOLD, 26)

	# 如果快完成了，加強提示
	if _killed_count == _total_count - 1:
		await get_tree().create_timer(0.3).timeout
		_spawn_float_text("🔥 最後一條！完全涅槃即將觸發！", Vector2(512, 280), COLOR_DEEP_ORANGE, 22)

func _on_rebirth_full(data: Dictionary) -> void:
	var player_name: String = data.get("player_name", "")
	var full_mult: float = data.get("full_rebirth_mult", 3.0)
	var duration: int = data.get("duration", 8)

	_is_active = false
	_full_rebirth_active = true
	_full_rebirth_remaining = float(duration)

	# 隱藏涅槃指示器和計時條
	_indicator.visible = false
	_timer_bar.visible = false

	# 全螢幕三次強閃光（火橙→金→深橙）
	_flash_screen(COLOR_FIRE_ORANGE, 1, 0.15)
	await get_tree().create_timer(0.2).timeout
	_flash_screen(COLOR_GOLD, 1, 0.15)
	await get_tree().create_timer(0.2).timeout
	_flash_screen(COLOR_DEEP_ORANGE, 1, 0.15)

	# 全螢幕大字
	_show_full_rebirth_text(player_name, full_mult, duration)

	# 顯示完全涅槃指示器
	_full_rebirth_indicator.visible = true
	_update_full_rebirth_indicator()

	# 全服廣播橫幅
	_show_banner(
		"🔥🦅🔥 鳳凰完全涅槃！",
		"%s 完全涅槃！全服 ×%.0f 加成 %d 秒！" % [player_name, full_mult, duration],
		COLOR_DEEP_ORANGE
	)

func _on_rebirth_full_end() -> void:
	_full_rebirth_active = false
	_full_rebirth_remaining = 0.0

	# 淡出完全涅槃指示器
	if _full_rebirth_indicator.visible:
		var tween := create_tween()
		tween.tween_property(_full_rebirth_indicator, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func(): _full_rebirth_indicator.visible = false)

	_spawn_float_text("🦅 完全涅槃結束", Vector2(512, 300), COLOR_GRAY, 20)

func _on_rebirth_fade(data: Dictionary) -> void:
	var remaining: int = data.get("remaining", 0)

	_is_active = false
	_timer_remaining = 0.0

	# 隱藏指示器和計時條
	_indicator.visible = false
	_timer_bar.visible = false

	# 灰色閃光
	_flash_screen(COLOR_GRAY, 1, 0.08)

	if remaining > 0:
		_spawn_float_text("🦅 涅槃消散... 尚有 %d 條未擊破" % remaining,
			Vector2(512, 300), COLOR_GRAY, 22)
	else:
		_spawn_float_text("🦅 涅槃完成！", Vector2(512, 300), COLOR_GOLD, 22)

# ── UI 建構輔助 ───────────────────────────────────────────────────────────────

func _make_banner() -> Control:
	var panel := PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_WIDE)
	panel.position = Vector2(0, -80)
	panel.size = Vector2(1024, 72)

	var style := StyleBoxFlat.new()
	style.bg_color = Color(0.0, 0.0, 0.0, 0.85)
	style.border_color = COLOR_FIRE_ORANGE
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
	title_lbl.add_theme_color_override("font_color", COLOR_FIRE_ORANGE)
	vbox.add_child(title_lbl)

	var msg_lbl := Label.new()
	msg_lbl.name = "MsgLabel"
	msg_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	msg_lbl.add_theme_font_size_override("font_size", 14)
	msg_lbl.add_theme_color_override("font_color", COLOR_LIGHT_ORANGE)
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
	style.border_color = COLOR_FIRE_ORANGE
	style.set_border_width_all(2)
	style.set_corner_radius_all(6)
	panel.add_theme_stylebox_override("panel", style)

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.name = "IndicatorTitle"
	title_lbl.text = "🦅 涅槃重生"
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", COLOR_FIRE_ORANGE)
	vbox.add_child(title_lbl)

	var count_lbl := Label.new()
	count_lbl.name = "CountLabel"
	count_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	count_lbl.add_theme_font_size_override("font_size", 18)
	count_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	vbox.add_child(count_lbl)

	var mult_lbl := Label.new()
	mult_lbl.name = "MultLabel"
	mult_lbl.text = "擊破 ×4.0"
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mult_lbl.add_theme_font_size_override("font_size", 12)
	mult_lbl.add_theme_color_override("font_color", COLOR_LIGHT_ORANGE)
	vbox.add_child(mult_lbl)

	return panel

func _make_timer_bar() -> Control:
	var container := Control.new()
	container.set_anchors_preset(Control.PRESET_RIGHT_WIDE)
	container.position = Vector2(-28, 160)
	container.size = Vector2(16, 200)

	var bg := ColorRect.new()
	bg.name = "BG"
	bg.color = Color(0.1, 0.1, 0.1, 0.8)
	bg.size = Vector2(16, 200)
	container.add_child(bg)

	var fill := ColorRect.new()
	fill.name = "Fill"
	fill.color = COLOR_FIRE_ORANGE
	fill.size = Vector2(16, 200)
	fill.position = Vector2(0, 0)
	container.add_child(fill)

	return container

func _make_full_rebirth_indicator() -> Control:
	var panel := PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	panel.position = Vector2(-200, 160)
	panel.size = Vector2(180, 60)

	var style := StyleBoxFlat.new()
	style.bg_color = Color(0.2, 0.0, 0.0, 0.9)
	style.border_color = COLOR_DEEP_ORANGE
	style.set_border_width_all(2)
	style.set_corner_radius_all(6)
	panel.add_theme_stylebox_override("panel", style)

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.name = "FullTitle"
	title_lbl.text = "🔥 完全涅槃 ×3.0"
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", COLOR_DEEP_ORANGE)
	vbox.add_child(title_lbl)

	var timer_lbl := Label.new()
	timer_lbl.name = "FullTimerLabel"
	timer_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	timer_lbl.add_theme_font_size_override("font_size", 18)
	timer_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	vbox.add_child(timer_lbl)

	return panel

# ── UI 更新輔助 ───────────────────────────────────────────────────────────────

func _show_banner(title: String, msg: String, color: Color) -> void:
	var title_lbl: Label = _banner.get_node("VBoxContainer/TitleLabel")
	var msg_lbl: Label = _banner.get_node("VBoxContainer/MsgLabel")
	title_lbl.text = title
	title_lbl.add_theme_color_override("font_color", color)
	msg_lbl.text = msg
	_banner.visible = true
	_banner.modulate.a = 1.0

	# 滑入動畫
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
	count_lbl.text = "%d / %d 擊破" % [_killed_count, _total_count]

	# 顏色隨進度變化
	var progress: float = float(_killed_count) / float(max(_total_count, 1))
	if progress >= 0.67:
		count_lbl.add_theme_color_override("font_color", COLOR_DEEP_ORANGE)
	elif progress >= 0.33:
		count_lbl.add_theme_color_override("font_color", COLOR_FIRE_ORANGE)
	else:
		count_lbl.add_theme_color_override("font_color", COLOR_GOLD)

	# 脈衝動畫
	var tween := create_tween()
	tween.tween_property(_indicator, "scale", Vector2(1.05, 1.05), 0.1)
	tween.tween_property(_indicator, "scale", Vector2(1.0, 1.0), 0.1)

func _update_timer_bar() -> void:
	if not _timer_bar.visible:
		return
	var fill: ColorRect = _timer_bar.get_node("Fill")
	var pct: float = _timer_remaining / float(max(_duration, 1))
	fill.size.y = 200.0 * pct
	fill.position.y = 200.0 * (1.0 - pct)

	# 顏色隨剩餘時間變化
	if pct > 0.5:
		fill.color = COLOR_FIRE_ORANGE
	elif pct > 0.25:
		fill.color = Color(1.0, 0.5, 0.0, 1.0)  # 橙
	else:
		fill.color = COLOR_DEEP_ORANGE

func _update_full_rebirth_indicator() -> void:
	if not _full_rebirth_indicator.visible:
		return
	var timer_lbl: Label = _full_rebirth_indicator.get_node("VBoxContainer/FullTimerLabel")
	timer_lbl.text = "%.1f 秒" % _full_rebirth_remaining

	# 脈衝動畫
	var tween := create_tween()
	tween.tween_property(_full_rebirth_indicator, "modulate:a", 0.7, 0.4)
	tween.tween_property(_full_rebirth_indicator, "modulate:a", 1.0, 0.4)

func _show_full_rebirth_text(player_name: String, mult: float, duration: int) -> void:
	# 全螢幕大字
	var lbl := Label.new()
	lbl.text = "🔥🦅🔥 鳳凰完全涅槃！\n全服 ×%.0f 加成 %d 秒！" % [mult, duration]
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.set_anchors_preset(Control.PRESET_CENTER)
	lbl.position = Vector2(512 - 300, 300 - 50)
	lbl.size = Vector2(600, 100)
	lbl.add_theme_font_size_override("font_size", 32)
	lbl.add_theme_color_override("font_color", COLOR_GOLD)
	add_child(lbl)

	# 縮放動畫
	lbl.scale = Vector2(0.5, 0.5)
	var tween := create_tween()
	tween.tween_property(lbl, "scale", Vector2(1.1, 1.1), 0.3).set_ease(Tween.EASE_OUT)
	tween.tween_property(lbl, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(2.5)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.5)
	tween.tween_callback(lbl.queue_free)

# ── 通用特效 ──────────────────────────────────────────────────────────────────

func _flash_screen(color: Color, times: int, duration: float) -> void:
	var overlay := ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, 0.0)
	overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)

	var tween := create_tween()
	for i in range(times):
		tween.tween_property(overlay, "color:a", 0.45, duration * 0.4)
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
