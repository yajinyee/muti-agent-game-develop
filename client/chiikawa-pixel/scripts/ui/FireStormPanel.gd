## FireStormPanel.gd ???怎憸冽擳?UI ?Ｘ嚗AY-176嚗?
## 璆剔?靘?嚗cean King 3 Plus?ire Storm feature??
## 憿舐內?怎憸冽閫貊????蝔???
extends CanvasLayer

# ---- 撣豢 ----
const PANEL_COLOR_BG    := Color(0.15, 0.05, 0.0, 0.88)
const PANEL_COLOR_FIRE  := Color(1.0, 0.35, 0.0, 1.0)   # 璈??莎??怎??
const PANEL_COLOR_GOLD  := Color(1.0, 0.85, 0.0, 1.0)
const PANEL_COLOR_WHITE := Color(1.0, 1.0, 1.0, 1.0)
const PANEL_COLOR_RED   := Color(1.0, 0.2, 0.2, 1.0)

# ---- 蝭暺???----
var _banner_container : Control
var _banner_label     : Label
var _result_panel     : Control
var _result_label     : Label
var _burn_counter     : Label
var _flash_overlay    : ColorRect

# ---- ???----
var _burn_count       : int = 0
var _total_reward     : int = 0
var _target_count     : int = 0
var _is_my_trigger    : bool = false

func _ready() -> void:
	layer = 68
	_build_ui()
	hide()

func _build_ui() -> void:
	# ?刻撟???overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 0.35, 0.0, 0.0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# ?璈怠?
	_banner_container = PanelContainer.new()
	_banner_container.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner_container.offset_top = 8
	_banner_container.offset_bottom = 56
	_banner_container.offset_left = 80
	_banner_container.offset_right = -80
	var banner_style := StyleBoxFlat.new()
	banner_style.bg_color = Color(0.6, 0.1, 0.0, 0.92)
	banner_style.corner_radius_top_left = 8
	banner_style.corner_radius_top_right = 8
	banner_style.corner_radius_bottom_left = 8
	banner_style.corner_radius_bottom_right = 8
	_banner_container.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner_container)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_FIRE)
	_banner_label.add_theme_font_size_override("font_size", 20)
	_banner_container.add_child(_banner_label)

	# ??閮?剁??喃?閫?
	_burn_counter = Label.new()
	_burn_counter.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	_burn_counter.offset_top = 64
	_burn_counter.offset_right = -16
	_burn_counter.offset_left = -200
	_burn_counter.offset_bottom = 96
	_burn_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	_burn_counter.add_theme_color_override("font_color", PANEL_COLOR_FIRE)
	_burn_counter.add_theme_font_size_override("font_size", 18)
	add_child(_burn_counter)

	# 蝯??Ｘ嚗?湔??伐?
	_result_panel = PanelContainer.new()
	_result_panel.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_panel.offset_right = -16
	_result_panel.offset_left = -280
	_result_panel.offset_top = -100
	_result_panel.offset_bottom = 100
	var result_style := StyleBoxFlat.new()
	result_style.bg_color = PANEL_COLOR_BG
	result_style.corner_radius_top_left = 12
	result_style.corner_radius_top_right = 12
	result_style.corner_radius_bottom_left = 12
	result_style.corner_radius_bottom_right = 12
	result_style.border_width_left = 2
	result_style.border_width_right = 2
	result_style.border_width_top = 2
	result_style.border_width_bottom = 2
	result_style.border_color = PANEL_COLOR_FIRE
	_result_panel.add_theme_stylebox_override("panel", result_style)
	_result_panel.modulate.a = 0.0
	add_child(_result_panel)

	_result_label = Label.new()
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_label.add_theme_color_override("font_color", PANEL_COLOR_WHITE)
	_result_label.add_theme_font_size_override("font_size", 16)
	_result_panel.add_child(_result_label)

## handle_fire_storm ?????怎憸冽擳???
func handle_fire_storm(payload: Dictionary) -> void:
	var phase : String = payload.get("phase", "")
	match phase:
		"fire_start":
			_on_fire_start(payload)
		"fire_burn":
			_on_fire_burn(payload)
		"fire_end":
			_on_fire_end(payload)

## _on_fire_start ???怎憸冽??
func _on_fire_start(payload: Dictionary) -> void:
	_target_count = payload.get("target_count", 0)
	_burn_count = 0
	_total_reward = 0
	var player_name : String = payload.get("player_name", "")
	var my_id : String = ""
	if GameManager.has_method("get_player_id"):
		my_id = GameManager.get_player_id()
	_is_my_trigger = (payload.get("player_id", "") == my_id and my_id != "")

	show()

	# ?刻撟?蝝???
	_flash_screen(Color(1.0, 0.35, 0.0, 0.5), 0.3)

	# 璈怠?皛
	_banner_label.text = "? %s 閫貊?怎憸冽嚗?d ?璅??葉嚗? % [player_name, _target_count]"
	_banner_container.modulate.a = 0.0
	_banner_container.show()
	var tween := create_tween()
	tween.tween_property(_banner_container, "modulate:a", 1.0, 0.3)

	# 閮??
	_burn_counter.text = "? 0 / %d" % _target_count
	_burn_counter.show()

	# ?芸楛閫貊??銝剖亢憭抒?唳?隤?
	if _is_my_trigger:
		_show_center_label("? ?怎憸冽嚗n敹怠?璅??格?嚗?, PANEL_COLOR_FIRE)"

## _on_fire_burn ???桀璅???
func _on_fire_burn(payload: Dictionary) -> void:
	var skipped : bool = payload.get("skipped", false)
	if skipped:
		return

	_burn_count += 1
	var reward : int = payload.get("reward", 0)
	_total_reward += reward

	# ?湔閮??
	_burn_counter.text = "? %d / %d" % [_burn_count, _target_count]

	# 撠???
	_flash_screen(Color(1.0, 0.5, 0.0, 0.25), 0.15)

	# 瘚桀????
	if reward > 0 and _is_my_trigger:
		_spawn_float_text("+%d" % reward, PANEL_COLOR_GOLD)

## _on_fire_end ???怎憸冽蝯?
func _on_fire_end(payload: Dictionary) -> void:
	var burned_count : int = payload.get("burned_count", 0)
	var total_reward : int = payload.get("total_reward", 0)

	# 蝯??Ｘ皛
	_result_label.text = "? ?怎憸冽蝯?\n??嚗?d ?璅n?嚗?%d ?馳" % [burned_count, total_reward]
	_result_panel.offset_right = 300
	_result_panel.offset_left = 20
	var tween := create_tween()
	tween.tween_property(_result_panel, "modulate:a", 1.0, 0.3)

	# 憭扯?璅∠?堆?????
	if burned_count >= 6:
		_double_flash(Color(1.0, 0.35, 0.0, 0.6))

	# 3 蝘?瘛∪
	await get_tree().create_timer(3.0).timeout
	_fade_out()

## _flash_screen ???刻撟???
func _flash_screen(color: Color, duration: float) -> void:
	_flash_overlay.color = color
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)

## _double_flash ??????
func _double_flash(color: Color) -> void:
	_flash_screen(color, 0.2)
	await get_tree().create_timer(0.25).timeout
	_flash_screen(color, 0.2)

## _show_center_label ??銝剖亢憭扳?隤?頝?
func _show_center_label(text: String, color: Color) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 28)
	lbl.set_anchors_preset(Control.PRESET_CENTER)
	lbl.offset_left = -200
	lbl.offset_right = 200
	lbl.offset_top = -60
	lbl.offset_bottom = 60
	lbl.scale = Vector2(0.5, 0.5)
	add_child(lbl)
	var tween := create_tween()
	tween.tween_property(lbl, "scale", Vector2(1.0, 1.0), 0.25).set_trans(Tween.TRANS_BACK)
	await get_tree().create_timer(2.0).timeout
	var fade := create_tween()
	fade.tween_property(lbl, "modulate:a", 0.0, 0.3)
	await fade.finished
	lbl.queue_free()

## _spawn_float_text ??瘚桀????
func _spawn_float_text(text: String, color: Color) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 20)
	lbl.set_anchors_preset(Control.PRESET_CENTER)
	lbl.offset_left = -60 + randi() % 120 - 60
	lbl.offset_top = -20
	lbl.offset_right = lbl.offset_left + 120
	lbl.offset_bottom = 20
	add_child(lbl)
	var tween := create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 40, 0.8)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.8)
	await tween.finished
	lbl.queue_free()

## _fade_out ??瘛∪???UI
func _fade_out() -> void:
	var tween := create_tween()
	tween.tween_property(self, "modulate:a", 0.0, 0.4)
	await tween.finished
	modulate.a = 1.0
	_banner_container.hide()
	_burn_counter.hide()
	_result_panel.modulate.a = 0.0
	hide()
