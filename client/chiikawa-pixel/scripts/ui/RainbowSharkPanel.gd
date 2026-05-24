## RainbowSharkPanel.gd ??敶抵攳?? UI ?Ｘ嚗AY-180嚗?
## 璆剔?靘?嚗ILI 2026?ainbow Shark ??triggers a rainbow burst that randomly assigns
## 1.5x-3x multiplier bonuses to all targets on screen for 10 seconds??
## 憿舐內敶抵??????璅???璅??閮????潛???
extends CanvasLayer

# ---- 撣豢 ----
const PANEL_COLOR_BG      := Color(0.05, 0.0, 0.1, 0.90)
const PANEL_COLOR_RAINBOW := Color(1.0, 0.4, 0.8, 1.0)   # ?梁?蝝?敶抵??
const PANEL_COLOR_GOLD    := Color(1.0, 0.85, 0.0, 1.0)
const PANEL_COLOR_WHITE   := Color(1.0, 1.0, 1.0, 1.0)
const PANEL_COLOR_CYAN    := Color(0.0, 1.0, 1.0, 1.0)

# ??撠?憿
const MULT_COLORS := {
	1.5: Color(0.5, 1.0, 0.5, 1.0),   # 瘛箇?
	2.0: Color(0.3, 0.8, 1.0, 1.0),   # 憭抵?
	2.5: Color(1.0, 0.8, 0.2, 1.0),   # ??
	3.0: Color(1.0, 0.3, 0.8, 1.0),   # ?梁?蝝?
}

# ---- 蝭暺???----
var _banner_container  : Control
var _banner_label      : Label
var _timer_label       : Label
var _flash_overlay     : ColorRect
var _mult_labels       : Dictionary = {}  # instanceID ??Label嚗?璅惜嚗?

# ---- ???----
var _burst_active      : bool = false
var _burst_remaining   : float = 0.0
var _marked_targets    : Dictionary = {}  # instanceID ??{x, y, burst_mult}
var _rainbow_phase     : float = 0.0     # 敶抵?脩??閮?

func _ready() -> void:
	layer = 66
	_build_ui()
	hide()

func _build_ui() -> void:
	# ?刻撟???overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 0.4, 0.8, 0.0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# ?璈怠?
	_banner_container = PanelContainer.new()
	_banner_container.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner_container.offset_top = 8
	_banner_container.offset_bottom = 56
	_banner_container.offset_left = 60
	_banner_container.offset_right = -60
	var banner_style := StyleBoxFlat.new()
	banner_style.bg_color = Color(0.2, 0.0, 0.3, 0.92)
	banner_style.corner_radius_top_left = 8
	banner_style.corner_radius_top_right = 8
	banner_style.corner_radius_bottom_left = 8
	banner_style.corner_radius_bottom_right = 8
	_banner_container.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner_container)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_RAINBOW)
	_banner_label.add_theme_font_size_override("font_size", 20)
	_banner_container.add_child(_banner_label)

	# ?閮??剁??喃?閫?
	_timer_label = Label.new()
	_timer_label.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	_timer_label.offset_top = 64
	_timer_label.offset_right = -16
	_timer_label.offset_left = -200
	_timer_label.offset_bottom = 96
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	_timer_label.add_theme_color_override("font_color", PANEL_COLOR_RAINBOW)
	_timer_label.add_theme_font_size_override("font_size", 18)
	_timer_label.hide()
	add_child(_timer_label)

func _process(delta: float) -> void:
	if _burst_active:
		_burst_remaining -= delta
		_rainbow_phase += delta * 2.0  # 敶抵?脩???漲
		if _burst_remaining <= 0.0:
			_burst_active = false
			_timer_label.hide()
		else:
			_timer_label.text = "?? 敶抵? %.1f 蝘? % _burst_remaining"
			# 璈怠?敶抵?脩??
			var hue := fmod(_rainbow_phase * 0.1, 1.0)
			_banner_label.add_theme_color_override("font_color", Color.from_hsv(hue, 0.8, 1.0))

## handle_rainbow_shark ????敶抵攳?閮
func handle_rainbow_shark(payload: Dictionary) -> void:
	var phase : String = payload.get("phase", "")
	match phase:
		"burst_start":
			_on_burst_start(payload)
		"burst_end":
			_on_burst_end()

## _on_burst_start ??敶抵???嚗??
func _on_burst_start(payload: Dictionary) -> void:
	_burst_active = true
	_burst_remaining = float(payload.get("duration_sec", 10))
	_rainbow_phase = 0.0
	var trigger_name : String = payload.get("trigger_name", "")
	var marked_targets : Array = payload.get("marked_targets", [])

	# ?脣?璅??格?
	_marked_targets.clear()
	for mk in marked_targets:
		_marked_targets[mk.get("instance_id", "")] = mk

	show()

	# ?刻撟蔗?寥???銝活嚗?
	_flash_rainbow_triple()

	# 璈怠?
	_banner_label.text = "?? %s 閫貊敶抵攳??嚗?d ?璅敺蔗?孵???" % [trigger_name, marked_targets.size()]
	_banner_container.modulate.a = 0.0
	_banner_container.show()
	var tween := create_tween()
	tween.tween_property(_banner_container, "modulate:a", 1.0, 0.3)

	# ?閮???
	_timer_label.show()

	# ?典銝＊蝷箸??璅???璅惜
	await get_tree().create_timer(0.3).timeout
	_spawn_mult_labels(marked_targets)

## _on_burst_end ??敶抵?蝯?嚗??
func _on_burst_end() -> void:
	_burst_active = false
	_burst_remaining = 0.0
	_timer_label.hide()
	_banner_container.hide()

	# 皜???璅惜
	_clear_mult_labels()

	# 瘛∪
	var tween := create_tween()
	tween.tween_property(self, "modulate:a", 0.0, 0.5)
	await tween.finished
	modulate.a = 1.0
	hide()

## _spawn_mult_labels ???典銝＊蝷箸??璅???璅惜
func _spawn_mult_labels(marked_targets: Array) -> void:
	_clear_mult_labels()
	for mk in marked_targets:
		var instance_id : String = mk.get("instance_id", "")
		var x : float = mk.get("x", 0.0)
		var y : float = mk.get("y", 0.0)
		var burst_mult : float = mk.get("burst_mult", 1.5)

		var lbl := Label.new()
		lbl.text = "?%.1f" % burst_mult
		lbl.add_theme_font_size_override("font_size", 14)

		# ?寞????豢?憿
		var color := PANEL_COLOR_WHITE
		if burst_mult >= 3.0:
			color = MULT_COLORS[3.0]
		elif burst_mult >= 2.5:
			color = MULT_COLORS[2.5]
		elif burst_mult >= 2.0:
			color = MULT_COLORS[2.0]
		else:
			color = MULT_COLORS[1.5]
		lbl.add_theme_color_override("font_color", color)

		# 雿蔭嚗璅??對?
		lbl.position = Vector2(x - 20, y - 30)
		lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE

		# 敶歲?
		lbl.scale = Vector2(0.3, 0.3)
		add_child(lbl)
		var tween := create_tween()
		tween.tween_property(lbl, "scale", Vector2(1.0, 1.0), 0.25).set_trans(Tween.TRANS_BACK)

		# ????嚗????游翰嚗?
		var blink_speed := 0.8 if burst_mult < 2.5 else 0.4
		_start_blink(lbl, blink_speed)

		_mult_labels[instance_id] = lbl

## _start_blink ??霈?蝐斗?蝥???
func _start_blink(lbl: Label, speed: float) -> void:
	var tween := lbl.create_tween().set_loops()
	tween.tween_property(lbl, "modulate:a", 0.4, speed)
	tween.tween_property(lbl, "modulate:a", 1.0, speed)

## _clear_mult_labels ??皜???璅惜
func _clear_mult_labels() -> void:
	for lbl in _mult_labels.values():
		if is_instance_valid(lbl):
			lbl.queue_free()
	_mult_labels.clear()

## remove_mark ??蝘駁撌脫??渡璅???璅惜嚗 GameManager ?澆嚗?
func remove_mark(instance_id: String) -> void:
	if _mult_labels.has(instance_id):
		var lbl : Label = _mult_labels[instance_id]
		if is_instance_valid(lbl):
			# ????豢?憭勗???
			var tween := create_tween()
			tween.tween_property(lbl, "scale", Vector2(2.0, 2.0), 0.15)
			tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.15)
			await tween.finished
			if is_instance_valid(lbl):
				lbl.queue_free()
		_mult_labels.erase(instance_id)

## _flash_rainbow_triple ??敶抵銝活??
func _flash_rainbow_triple() -> void:
	var colors := [
		Color(1.0, 0.3, 0.3, 0.35),  # 蝝?
		Color(0.3, 1.0, 0.3, 0.35),  # 蝬?
		Color(0.3, 0.3, 1.0, 0.35),  # ??
	]
	for c in colors:
		_flash_overlay.color = c
		var tween := create_tween()
		tween.tween_property(_flash_overlay, "color:a", 0.0, 0.2)
		await get_tree().create_timer(0.25).timeout

## _spawn_float_text ??瘚桀???
func _spawn_float_text(text: String, color: Color, pos: Vector2) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 20)
	lbl.position = pos
	lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(lbl)
	var tween := create_tween()
	tween.tween_property(lbl, "position:y", pos.y - 60, 0.9)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.9)
	await tween.finished
	if is_instance_valid(lbl):
		lbl.queue_free()
