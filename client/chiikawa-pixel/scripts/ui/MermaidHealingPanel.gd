## MermaidHealingPanel.gd — 美人魚治癒 UI 面板（DAY-178）
## 業界依據：Ocean King 3 Plus「The Mermaid feature」
## 顯示美人魚治癒、幸運加成效果
extends CanvasLayer

# ---- 常數 ----
const PANEL_COLOR_BG     := Color(0.0, 0.1, 0.15, 0.90)
const PANEL_COLOR_TEAL   := Color(0.0, 0.8, 0.85, 1.0)   # 深青色（美人魚感）
const PANEL_COLOR_GOLD   := Color(1.0, 0.85, 0.0, 1.0)
const PANEL_COLOR_WHITE  := Color(1.0, 1.0, 1.0, 1.0)
const PANEL_COLOR_GREEN  := Color(0.2, 1.0, 0.5, 1.0)

# ---- 節點引用 ----
var _banner_container  : Control
var _banner_label      : Label
var _heal_popup        : Control
var _heal_label        : Label
var _luck_timer_label  : Label
var _flash_overlay     : ColorRect

# ---- 狀態 ----
var _is_my_heal        : bool = false
var _luck_active       : bool = false
var _luck_remaining    : float = 0.0

func _ready() -> void:
	layer = 66
	_build_ui()
	hide()

func _build_ui() -> void:
	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(0.0, 0.8, 0.85, 0.0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner_container = PanelContainer.new()
	_banner_container.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner_container.offset_top = 8
	_banner_container.offset_bottom = 56
	_banner_container.offset_left = 80
	_banner_container.offset_right = -80
	var banner_style := StyleBoxFlat.new()
	banner_style.bg_color = Color(0.0, 0.3, 0.35, 0.92)
	banner_style.corner_radius_top_left = 8
	banner_style.corner_radius_top_right = 8
	banner_style.corner_radius_bottom_left = 8
	banner_style.corner_radius_bottom_right = 8
	_banner_container.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner_container)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_TEAL)
	_banner_label.add_theme_font_size_override("font_size", 20)
	_banner_container.add_child(_banner_label)

	# 治癒彈窗（中央）
	_heal_popup = PanelContainer.new()
	_heal_popup.set_anchors_preset(Control.PRESET_CENTER)
	_heal_popup.offset_left = -140
	_heal_popup.offset_right = 140
	_heal_popup.offset_top = -60
	_heal_popup.offset_bottom = 60
	var popup_style := StyleBoxFlat.new()
	popup_style.bg_color = PANEL_COLOR_BG
	popup_style.corner_radius_top_left = 16
	popup_style.corner_radius_top_right = 16
	popup_style.corner_radius_bottom_left = 16
	popup_style.corner_radius_bottom_right = 16
	popup_style.border_width_left = 2
	popup_style.border_width_right = 2
	popup_style.border_width_top = 2
	popup_style.border_width_bottom = 2
	popup_style.border_color = PANEL_COLOR_TEAL
	_heal_popup.add_theme_stylebox_override("panel", popup_style)
	_heal_popup.modulate.a = 0.0
	add_child(_heal_popup)

	_heal_label = Label.new()
	_heal_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_heal_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_heal_label.add_theme_color_override("font_color", PANEL_COLOR_WHITE)
	_heal_label.add_theme_font_size_override("font_size", 22)
	_heal_popup.add_child(_heal_label)

	# 幸運加成計時器（左上角）
	_luck_timer_label = Label.new()
	_luck_timer_label.set_anchors_preset(Control.PRESET_TOP_LEFT)
	_luck_timer_label.offset_top = 64
	_luck_timer_label.offset_left = 16
	_luck_timer_label.offset_right = 220
	_luck_timer_label.offset_bottom = 96
	_luck_timer_label.add_theme_color_override("font_color", PANEL_COLOR_TEAL)
	_luck_timer_label.add_theme_font_size_override("font_size", 18)
	_luck_timer_label.hide()
	add_child(_luck_timer_label)

func _process(delta: float) -> void:
	if _luck_active:
		_luck_remaining -= delta
		if _luck_remaining <= 0.0:
			_luck_active = false
			_luck_timer_label.hide()
		else:
			_luck_timer_label.text = "🧜 幸運 +20%% %.1f 秒" % _luck_remaining

## handle_mermaid_healing — 處理美人魚治癒訊息
func handle_mermaid_healing(payload: Dictionary) -> void:
	var phase : String = payload.get("phase", "")
	match phase:
		"heal_start":
			_on_heal_start(payload)
		"heal_broadcast":
			_on_heal_broadcast(payload)
		"luck_start":
			_on_luck_start(payload)
		"luck_end":
			_on_luck_end()

## _on_heal_start — 治癒開始（個人）
func _on_heal_start(payload: Dictionary) -> void:
	_is_my_heal = true
	var heal_amount : int = payload.get("heal_amount", 0)

	show()

	# 全螢幕青色閃光
	_flash_screen(Color(0.0, 0.8, 0.85, 0.45), 0.5)

	# 橫幅
	_banner_label.text = "🧜 美人魚降臨！恢復 +%d 金幣！" % heal_amount
	_banner_container.modulate.a = 0.0
	_banner_container.show()
	var tween := create_tween()
	tween.tween_property(_banner_container, "modulate:a", 1.0, 0.3)

	# 治癒彈窗彈跳
	_heal_label.text = "🧜‍♀️\n+%d 金幣\n治癒！" % heal_amount
	_heal_popup.scale = Vector2(0.5, 0.5)
	_heal_popup.modulate.a = 1.0
	var popup_tween := create_tween()
	popup_tween.tween_property(_heal_popup, "scale", Vector2(1.0, 1.0), 0.3).set_trans(Tween.TRANS_BACK)

	# 浮動治癒文字
	_spawn_float_text("+%d 💚" % heal_amount, PANEL_COLOR_GREEN)

	# 2 秒後淡出彈窗
	await get_tree().create_timer(2.0).timeout
	var fade := create_tween()
	fade.tween_property(_heal_popup, "modulate:a", 0.0, 0.4)

## _on_heal_broadcast — 全服廣播（其他玩家看到）
func _on_heal_broadcast(payload: Dictionary) -> void:
	if _is_my_heal:
		return
	var player_name : String = payload.get("player_name", "")
	var heal_amount : int = payload.get("heal_amount", 0)
	_show_broadcast_banner("🧜 %s 遇見美人魚！恢復了 %d 金幣！" % [player_name, heal_amount])

## _on_luck_start — 幸運加成開始（全服）
func _on_luck_start(payload: Dictionary) -> void:
	_luck_active = true
	_luck_remaining = float(payload.get("luck_boost_duration_sec", 20))
	_luck_timer_label.show()
	show()

## _on_luck_end — 幸運加成結束（全服）
func _on_luck_end() -> void:
	_luck_active = false
	_luck_timer_label.hide()
	_is_my_heal = false
	_banner_container.hide()
	hide()

## _show_broadcast_banner — 顯示全服廣播橫幅
func _show_broadcast_banner(text: String) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_color_override("font_color", PANEL_COLOR_TEAL)
	lbl.add_theme_font_size_override("font_size", 16)
	lbl.set_anchors_preset(Control.PRESET_TOP_WIDE)
	lbl.offset_top = 60
	lbl.offset_bottom = 90
	lbl.modulate.a = 0.0
	add_child(lbl)
	show()
	var tween := create_tween()
	tween.tween_property(lbl, "modulate:a", 1.0, 0.3)
	await get_tree().create_timer(2.5).timeout
	var fade := create_tween()
	fade.tween_property(lbl, "modulate:a", 0.0, 0.4)
	await fade.finished
	lbl.queue_free()

## _flash_screen — 全螢幕閃光
func _flash_screen(color: Color, duration: float) -> void:
	_flash_overlay.color = color
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)

## _spawn_float_text — 浮動文字
func _spawn_float_text(text: String, color: Color) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 24)
	lbl.set_anchors_preset(Control.PRESET_CENTER)
	lbl.offset_left = -80
	lbl.offset_right = 80
	lbl.offset_top = -20
	lbl.offset_bottom = 20
	add_child(lbl)
	var tween := create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 60, 1.0)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 1.0)
	await tween.finished
	lbl.queue_free()
