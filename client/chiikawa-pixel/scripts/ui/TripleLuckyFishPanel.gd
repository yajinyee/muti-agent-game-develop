## TripleLuckyFishPanel.gd ??銝?撟賊?擳?選?DAY-190嚗?
## 璆剔???嚗aDa Gaming TriLuck??2026?rigger three different feature specifications simultaneously??
## 銝?????閫貊嚗?撟? + ???? + 甇血?
## 閬死銝駁?嚗??脖??? + 銝?? + 敶抵??

extends Control

const TRIPLE_COLOR := Color(1.0, 0.85, 0.0)   # ?嚗??兢??
const COIN_COLOR   := Color(1.0, 0.9, 0.2)    # 鈭桅?嚗?撟?嚗?
const MULT_COLOR   := Color(0.4, 1.0, 0.4)    # 鈭桃?嚗???嚗?
const WEAPON_COLOR := Color(0.4, 0.8, 1.0)    # 憭抵?嚗郎?典??踝?

var _banner: Control = null
var _mult_bar: Control = null
var _mult_timer: float = 0.0
var _mult_duration: float = 12.0
var _is_my_trigger: bool = false

func _ready() -> void:
	set_process(false)
	if GameManager.has_signal("triple_lucky_fish"):
		GameManager.triple_lucky_fish.connect(_on_triple_lucky_fish)

func _process(delta: float) -> void:
	if _mult_bar == null or not is_instance_valid(_mult_bar):
		set_process(false)
		return
	_mult_timer -= delta
	if _mult_timer <= 0.0:
		_mult_timer = 0.0
		set_process(false)
		_fade_out_mult_bar()
		return
	# ?湔?脣漲璇?
	var pct := _mult_timer / _mult_duration
	var bar_fill = _mult_bar.get_node_or_null("Fill")
	if bar_fill:
		bar_fill.size.x = 280.0 * pct
		# 憿瞍貉?嚗?????
		if pct > 0.5:
			bar_fill.color = MULT_COLOR
		elif pct > 0.25:
			bar_fill.color = Color(1.0, 0.85, 0.0)
		else:
			bar_fill.color = Color(1.0, 0.5, 0.0)

func _on_triple_lucky_fish(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	match phase:
		"triple_start":
			_is_my_trigger = true
			_show_triple_start(data)
		"triple_broadcast":
			if not _is_my_trigger:
				_show_broadcast_banner(data)
		"mult_end":
			_is_my_trigger = false
			_fade_out_mult_bar()

func _show_triple_start(data: Dictionary) -> void:
	var coin_reward: int = data.get("coin_reward", 0)
	var coin_mult: float = data.get("coin_mult", 0.0)
	var mult_duration: float = data.get("mult_duration", 12.0)
	var weapon_charged: String = data.get("weapon_charged", "")"

	_mult_duration = mult_duration

	# ?刻撟??蔗?寥???銝活嚗???150ms嚗?
	_triple_rainbow_flash()

	# ?璈怠?皛
	_show_banner("?? 銝?撟賊?閫貊嚗?, TRIPLE_COLOR)"

	# 銝????＊蝷綽??舫? 200ms嚗?
	var timer1 = get_tree().create_timer(0.3)
	timer1.timeout.connect(func():
		_show_effect_popup("? ?馳??, "+%d ?馳嚗?.0f嚗? % [coin_reward, coin_mult], COIN_COLOR, Vector2(200, 280))
	)
	var timer2 = get_tree().create_timer(0.5)
	timer2.timeout.connect(func():
		_show_effect_popup("??????", "+50%% ?? %.0f 蝘? % mult_duration, MULT_COLOR, Vector2(640, 280))"
	)
	var timer3 = get_tree().create_timer(0.7)
	timer3.timeout.connect(func():
		_show_effect_popup("??甇血?", "%s ? 1 ?? % weapon_charged, WEAPON_COLOR, Vector2(1080, 280))"
	)

	# 摨?????脣漲璇?
	var timer4 = get_tree().create_timer(0.9)
	timer4.timeout.connect(func():
		_show_mult_bar(mult_duration)
	)

func _show_broadcast_banner(data: Dictionary) -> void:
	var player_name: String = data.get("player_name", "")
	var coin_reward: int = data.get("coin_reward", 0)
	_show_banner("?? %s 閫貊銝?撟賊?嚗?%d ?馳" % [player_name, coin_reward], TRIPLE_COLOR)

func _triple_rainbow_flash() -> void:
	# 銝活敶抵??嚗?????嚗?
	var colors := [TRIPLE_COLOR, MULT_COLOR, WEAPON_COLOR]
	for i in range(3):
		var t = get_tree().create_timer(float(i) * 0.15)
		var c = colors[i]
		t.timeout.connect(func():
			_flash_screen(c, 0.5)
		)

func _flash_screen(color: Color, alpha: float) -> void:
	var overlay = ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, alpha)
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tween = create_tween()
	tween.tween_property(overlay, "modulate:a", 0.0, 0.25)
	tween.tween_callback(overlay.queue_free)

func _show_banner(text: String, color: Color) -> void:
	if _banner != null and is_instance_valid(_banner):
		_banner.queue_free()

	_banner = Control.new()
	_banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 60)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.color = Color(color.r * 0.3, color.g * 0.3, color.b * 0.3, 0.92)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_banner.add_child(bg)

	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", 22)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.add_child(label)

	# 皛?
	_banner.position.y = -60
	var tween = create_tween()
	tween.tween_property(_banner, "position:y", 0.0, 0.3).set_trans(Tween.TRANS_BACK)
	# 4 蝘?瘛∪
	var timer = get_tree().create_timer(4.0)
	timer.timeout.connect(func():
		if is_instance_valid(_banner):
			var t2 = create_tween()
			t2.tween_property(_banner, "modulate:a", 0.0, 0.4)
			t2.tween_callback(_banner.queue_free)
	)

func _show_effect_popup(title: String, desc: String, color: Color, pos: Vector2) -> void:
	var popup = Control.new()
	popup.position = pos - Vector2(80, 60)
	popup.custom_minimum_size = Vector2(160, 120)
	add_child(popup)

	var bg = ColorRect.new()
	bg.color = Color(color.r * 0.2, color.g * 0.2, color.b * 0.2, 0.9)
	bg.size = Vector2(160, 120)
	bg.position = Vector2.ZERO
	popup.add_child(bg)

	# ??
	var border = ColorRect.new()
	border.color = color
	border.size = Vector2(160, 4)
	border.position = Vector2(0, 0)
	popup.add_child(border)

	var title_label = Label.new()
	title_label.text = title
	title_label.add_theme_color_override("font_color", color)
	title_label.add_theme_font_size_override("font_size", 18)
	title_label.position = Vector2(8, 12)
	title_label.size = Vector2(144, 30)
	title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup.add_child(title_label)

	var desc_label = Label.new()
	desc_label.text = desc
	desc_label.add_theme_color_override("font_color", Color.WHITE)
	desc_label.add_theme_font_size_override("font_size", 14)
	desc_label.position = Vector2(8, 50)
	desc_label.size = Vector2(144, 60)
	desc_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	desc_label.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	popup.add_child(desc_label)

	# 敶歲?
	popup.scale = Vector2(0.5, 0.5)
	popup.pivot_offset = Vector2(80, 60)
	var tween = create_tween()
	tween.tween_property(popup, "scale", Vector2(1.0, 1.0), 0.3).set_trans(Tween.TRANS_BACK)

	# 3 蝘?瘛∪
	var timer = get_tree().create_timer(3.0)
	timer.timeout.connect(func():
		if is_instance_valid(popup):
			var t2 = create_tween()
			t2.tween_property(popup, "modulate:a", 0.0, 0.4)
			t2.tween_callback(popup.queue_free)
	)

func _show_mult_bar(duration: float) -> void:
	if _mult_bar != null and is_instance_valid(_mult_bar):
		_mult_bar.queue_free()

	_mult_bar = Control.new()
	_mult_bar.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	_mult_bar.offset_top = -40
	_mult_bar.offset_bottom = 0
	add_child(_mult_bar)

	# ?
	var bg = ColorRect.new()
	bg.color = Color(0.1, 0.1, 0.1, 0.8)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_mult_bar.add_child(bg)

	# ?脣漲璇‵??
	var fill = ColorRect.new()
	fill.name = "Fill"
	fill.color = MULT_COLOR
	fill.size = Vector2(280.0, 20)
	fill.position = Vector2(get_viewport_rect().size.x / 2.0 - 140.0, 10)
	_mult_bar.add_child(fill)

	# 璅惜
	var label = Label.new()
	label.text = "?? 銝?撟賊? +50%% ????"
	label.add_theme_color_override("font_color", MULT_COLOR)
	label.add_theme_font_size_override("font_size", 14)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_mult_bar.add_child(label)

	_mult_timer = duration
	set_process(true)

func _fade_out_mult_bar() -> void:
	if _mult_bar != null and is_instance_valid(_mult_bar):
		var tween = create_tween()
		tween.tween_property(_mult_bar, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_mult_bar.queue_free)
		_mult_bar = null
