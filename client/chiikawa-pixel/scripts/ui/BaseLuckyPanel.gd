## BaseLuckyPanel.gd — 幸運特殊魚面板基礎類別
## lucky-panel-agent 負責維護
## 所有 LuckyXxxPanel 的共用基礎，減少重複程式碼
## 使用方式：extends BaseLuckyPanel（但 Godot 4 inner class 不支援 extends，改用組合模式）
## 本腳本作為 HUD 的輔助工具，提供標準化的 Lucky Panel 建立和動畫方法
extends Node

# ── 常數 ─────────────────────────────────────────────────────
const BANNER_HEIGHT = 80.0
const INDICATOR_SIZE = Vector2(200, 60)
const SETTLE_PANEL_SIZE = Vector2(320, 160)

# ── 標準顏色主題 ──────────────────────────────────────────────
const THEME_GOLD    = Color(1.0, 0.85, 0.0)
const THEME_RED     = Color(1.0, 0.2, 0.2)
const THEME_BLUE    = Color(0.0, 0.6, 1.0)
const THEME_ICE     = Color(0.0, 0.9, 1.0)
const THEME_PURPLE  = Color(0.7, 0.2, 0.9)
const THEME_ORANGE  = Color(1.0, 0.5, 0.1)
const THEME_GREEN   = Color(0.2, 0.9, 0.2)
const THEME_PINK    = Color(1.0, 0.42, 0.71)

# ── 建立標準橫幅 ─────────────────────────────────────────────
## 在 canvas_layer 上建立一個全寬橫幅
## 回傳 {panel, label} 字典
static func create_banner(canvas_layer: CanvasLayer, y_pos: float = 120.0, z_index: int = 60) -> Dictionary:
	var panel = Control.new()
	panel.position = Vector2(0, y_pos)
	panel.size = Vector2(1280, BANNER_HEIGHT)
	panel.visible = false
	panel.z_index = z_index
	canvas_layer.add_child(panel)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0, 0, 0, 0.8)
	panel.add_child(bg)

	var lbl = Label.new()
	lbl.name = "BannerLabel"
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 28)
	lbl.modulate = THEME_GOLD
	panel.add_child(lbl)

	return {"panel": panel, "label": lbl}

# ── 顯示橫幅動畫 ─────────────────────────────────────────────
static func show_banner(panel: Control, text: String, color: Color, duration: float = 2.5) -> void:
	if not is_instance_valid(panel):
		return
	var lbl = panel.get_node_or_null("BannerLabel")
	if is_instance_valid(lbl):
		lbl.text = text
		lbl.modulate = color
	panel.visible = true
	panel.modulate.a = 1.0
	var tween = panel.create_tween()
	tween.tween_interval(max(0.1, duration - 0.5))
	tween.tween_property(panel, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(panel):
			panel.visible = false
	)

# ── 建立右上角指示器 ─────────────────────────────────────────
## 顯示進度/計時/倍率等資訊
static func create_indicator(canvas_layer: CanvasLayer, pos: Vector2, z_index: int = 65) -> Dictionary:
	var panel = Control.new()
	panel.position = pos
	panel.size = INDICATOR_SIZE
	panel.visible = false
	panel.z_index = z_index
	canvas_layer.add_child(panel)

	var bg = ColorRect.new()
	bg.size = INDICATOR_SIZE
	bg.color = Color(0, 0, 0, 0.75)
	panel.add_child(bg)

	var title_lbl = Label.new()
	title_lbl.name = "TitleLabel"
	title_lbl.position = Vector2(8, 4)
	title_lbl.size = Vector2(184, 24)
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.modulate = THEME_GOLD
	panel.add_child(title_lbl)

	var value_lbl = Label.new()
	value_lbl.name = "ValueLabel"
	value_lbl.position = Vector2(8, 28)
	value_lbl.size = Vector2(184, 28)
	value_lbl.add_theme_font_size_override("font_size", 22)
	value_lbl.modulate = Color.WHITE
	panel.add_child(value_lbl)

	return {"panel": panel, "title": title_lbl, "value": value_lbl}

# ── 建立計時條 ───────────────────────────────────────────────
static func create_timer_bar(parent: Control, width: float = 180.0, color: Color = THEME_GOLD) -> Dictionary:
	var bg = ColorRect.new()
	bg.name = "TimerBarBG"
	bg.size = Vector2(width, 8)
	bg.color = Color(0.2, 0.2, 0.2, 0.8)
	parent.add_child(bg)

	var bar = ColorRect.new()
	bar.name = "TimerBar"
	bar.size = Vector2(width, 8)
	bar.color = color
	parent.add_child(bar)

	return {"bg": bg, "bar": bar, "max_width": width}

# ── 更新計時條 ───────────────────────────────────────────────
static func update_timer_bar(bar_dict: Dictionary, pct: float) -> void:
	var bar = bar_dict.get("bar")
	if not is_instance_valid(bar):
		return
	var max_w = bar_dict.get("max_width", 180.0)
	bar.size.x = max_w * clamp(pct, 0.0, 1.0)
	# 顏色隨進度：金→橙→紅
	if pct > 0.6:
		bar.color = THEME_GOLD
	elif pct > 0.3:
		bar.color = THEME_ORANGE
	else:
		bar.color = THEME_RED

# ── 建立結算彈窗 ─────────────────────────────────────────────
static func create_settle_popup(canvas_layer: CanvasLayer, z_index: int = 70) -> Control:
	var panel = Control.new()
	panel.position = Vector2(1280, 300)  # 從右側滑入
	panel.size = SETTLE_PANEL_SIZE
	panel.visible = false
	panel.z_index = z_index
	canvas_layer.add_child(panel)

	var bg = ColorRect.new()
	bg.size = SETTLE_PANEL_SIZE
	bg.color = Color(0.05, 0.05, 0.1, 0.92)
	panel.add_child(bg)

	var border = ColorRect.new()
	border.size = Vector2(SETTLE_PANEL_SIZE.x, 3)
	border.color = THEME_GOLD
	panel.add_child(border)

	return panel

# ── 顯示結算彈窗（從右側滑入）────────────────────────────────
static func show_settle_popup(panel: Control, lines: Array, duration: float = 3.0) -> void:
	if not is_instance_valid(panel):
		return
	# 清除舊內容（保留 bg 和 border）
	for child in panel.get_children():
		if child is Label:
			child.queue_free()

	# 加入文字行
	for i in range(lines.size()):
		var line = lines[i]
		var lbl = Label.new()
		lbl.text = line.get("text", "")
		lbl.position = Vector2(12, 20 + i * 32)
		lbl.add_theme_font_size_override("font_size", line.get("size", 16))
		lbl.modulate = line.get("color", Color.WHITE)
		panel.add_child(lbl)

	panel.visible = true
	panel.position.x = 1280
	var tween = panel.create_tween()
	tween.tween_property(panel, "position:x", 960.0, 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.tween_interval(duration - 0.8)
	tween.tween_property(panel, "position:x", 1280.0, 0.3).set_ease(Tween.EASE_IN)
	tween.tween_callback(func():
		if is_instance_valid(panel):
			panel.visible = false
	)

# ── 全螢幕閃光 ───────────────────────────────────────────────
static func fullscreen_flash(canvas_layer: CanvasLayer, color: Color, times: int = 3) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.color = color
	flash.color.a = 0.0
	flash.z_index = 90
	canvas_layer.add_child(flash)

	var tween = flash.create_tween()
	for i in range(times):
		tween.tween_property(flash, "modulate:a", 0.7, 0.06)
		tween.tween_property(flash, "modulate:a", 0.0, 0.1)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)

# ── 脈動動畫 ─────────────────────────────────────────────────
static func start_pulse(node: Node, min_alpha: float = 0.3, max_alpha: float = 1.0, period: float = 0.5) -> Tween:
	if not is_instance_valid(node):
		return null
	var tween = node.create_tween().set_loops()
	tween.tween_property(node, "modulate:a", min_alpha, period / 2)
	tween.tween_property(node, "modulate:a", max_alpha, period / 2)
	return tween

# ── 浮動文字 ─────────────────────────────────────────────────
static func spawn_float_text(parent: Node, pos: Vector2, text: String, color: Color, font_size: int = 20) -> void:
	var lbl = Label.new()
	lbl.text = text
	lbl.position = pos
	lbl.add_theme_font_size_override("font_size", font_size)
	lbl.modulate = color
	lbl.z_index = 80
	parent.add_child(lbl)

	var tween = lbl.create_tween()
	tween.tween_property(lbl, "position:y", pos.y - 50, 0.8)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.8)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)
