## LuckyStarFishPanel.gd ??撟賊????典??蝧餃?選?DAY-160嚗?
## ?撟賊???敺孛?澆?游?蝧餃?10 蝘???????2
## 閬死嚗??脫?????+ ?刻撟??脣???+ ?閮???+ ?2 璅?
extends Node2D

var _pixel_font: Font = null
var _countdown_label: Label = null
var _mult_label: Label = null
var _overlay: ColorRect = null
var _is_active: bool = false
var _remaining: float = 0.0
var _my_player_id: String = ""

func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

func _connect_signals() -> void:
	if GameManager.has_signal("lucky_star_fish"):
		GameManager.lucky_star_fish.connect(_on_lucky_star_fish)

func _process(delta: float) -> void:
	if not _is_active:
		return
	_remaining -= delta
	if _remaining <= 0.0:
		_remaining = 0.0
		_is_active = false
		_hide_lucky()
		return
	if is_instance_valid(_countdown_label):
		_countdown_label.text = "潃??2 %.1f蝘? % _remaining"

## ??撟賊???鈭辣
func _on_lucky_star_fish(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var duration: float = data.get("duration_secs", 10.0)
	var killer_name: String = data.get("killer_name", "")
	var killer_id: String = data.get("killer_id", "")

	if phase == "lucky_start":
		_show_lucky(killer_name, killer_id, duration)
	elif phase == "lucky_end":
		_hide_lucky()

## 憿舐內??蝧餃???
func _show_lucky(killer_name: String, killer_id: String, duration: float) -> void:
	_is_active = true
	_remaining = duration
	_my_player_id = NetworkManager.get_player_id() if NetworkManager.has_method("get_player_id") else ""

	var is_me = (killer_id == _my_player_id)

	# ?刻撟??脣????桃蔗嚗撌梯孛?潭??游撥??
	_overlay = ColorRect.new()
	_overlay.position = Vector2(-640, -360)
	_overlay.size = Vector2(1280, 720)
	_overlay.color = Color(1.0, 0.9, 0.0, 0.06 if not is_me else 0.12)
	add_child(_overlay)

	# ?刻撟??脤???
	var flash := ColorRect.new()
	flash.position = Vector2(-640, -360)
	flash.size = Vector2(1280, 720)
	flash.color = Color(1.0, 0.9, 0.0, 0.5 if is_me else 0.3)
	add_child(flash)
	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", 0.0, 0.5)
	flash_tween.tween_callback(func():
		if is_instance_valid(flash): flash.queue_free()
	)

	# ?璈怠?
	var banner_bg := ColorRect.new()
	banner_bg.name = "LuckyBanner"
	banner_bg.position = Vector2(-640, -360)
	banner_bg.size = Vector2(1280, 52)
	banner_bg.color = Color(0.15, 0.12, 0.0, 0.92)
	add_child(banner_bg)

	var banner_text = "潃?%s ?撟賊???嚗?渡????2 ?? %.0f 蝘?" % [killer_name, duration]
	if is_me:
		banner_text = "潃?雿??游兢??擳?雿?? ?2 ?? %.0f 蝘?" % duration

	var banner_lbl := Label.new()
	banner_lbl.name = "LuckyBannerLabel"
	banner_lbl.position = Vector2(-640, -354)
	banner_lbl.size = Vector2(1280, 44)
	banner_lbl.text = banner_text
	banner_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	banner_lbl.add_theme_color_override("font_color", Color(1.0, 0.95, 0.2))
	banner_lbl.add_theme_font_size_override("font_size", 16)
	if _pixel_font:
		banner_lbl.add_theme_font_override("font", _pixel_font)
	add_child(banner_lbl)

	# 璈怠?皛?
	banner_bg.position.y = -410
	banner_lbl.position.y = -404
	var slide_tween = banner_bg.create_tween()
	slide_tween.tween_property(banner_bg, "position:y", -360.0, 0.3)
	var slide_tween2 = banner_lbl.create_tween()
	slide_tween2.tween_property(banner_lbl, "position:y", -354.0, 0.3)

	# ?閮??剁??喃?閫?
	var countdown_bg := ColorRect.new()
	countdown_bg.name = "LuckyCountdownBG"
	countdown_bg.position = Vector2(510, -355)
	countdown_bg.size = Vector2(120, 36)
	countdown_bg.color = Color(0.12, 0.1, 0.0, 0.9)
	add_child(countdown_bg)

	_countdown_label = Label.new()
	_countdown_label.name = "LuckyCountdown"
	_countdown_label.position = Vector2(510, -352)
	_countdown_label.size = Vector2(120, 30)
	_countdown_label.text = "潃??2 %.1f蝘? % duration"
	_countdown_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_countdown_label.add_theme_color_override("font_color", Color(1.0, 0.95, 0.2))
	_countdown_label.add_theme_font_size_override("font_size", 13)
	if _pixel_font:
		_countdown_label.add_theme_font_override("font", _pixel_font)
	add_child(_countdown_label)

	# ?芸楛閫貊??銝剖亢憭??2 璅?
	if is_me:
		var mult_bg := ColorRect.new()
		mult_bg.name = "MultBG"
		mult_bg.position = Vector2(-80, -50)
		mult_bg.size = Vector2(160, 80)
		mult_bg.color = Color(0.15, 0.12, 0.0, 0.95)
		add_child(mult_bg)

		_mult_label = Label.new()
		_mult_label.name = "MultLabel"
		_mult_label.position = Vector2(-80, -48)
		_mult_label.size = Vector2(160, 76)
		_mult_label.text = "?2"
		_mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		_mult_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
		_mult_label.add_theme_color_override("font_color", Color(1.0, 0.95, 0.2))
		_mult_label.add_theme_font_size_override("font_size", 48)
		if _pixel_font:
			_mult_label.add_theme_font_override("font", _pixel_font)
		add_child(_mult_label)

		# ?2 璅?蝮格敶歲?
		_mult_label.scale = Vector2(0.5, 0.5)
		var scale_tween = _mult_label.create_tween()
		scale_tween.tween_property(_mult_label, "scale", Vector2(1.2, 1.2), 0.2)
		scale_tween.tween_property(_mult_label, "scale", Vector2(1.0, 1.0), 0.1)

		# 2 蝘?瘛∪ ?2 璅?
		var fade_tween = _mult_label.create_tween()
		fade_tween.tween_interval(2.0)
		fade_tween.tween_property(_mult_label, "modulate:a", 0.0, 0.5)
		fade_tween.tween_callback(func():
			if is_instance_valid(_mult_label): _mult_label.queue_free()
			if is_instance_valid(mult_bg): mult_bg.queue_free()
			_mult_label = null
		)

## ?梯???蝧餃???
func _hide_lucky() -> void:
	_is_active = false
	_remaining = 0.0

	for child_name in ["LuckyBanner", "LuckyBannerLabel", "LuckyCountdownBG", "LuckyCountdown", "MultBG", "MultLabel"]:
		var node = get_node_or_null(child_name)
		if is_instance_valid(node):
			var tween = node.create_tween()
			tween.tween_property(node, "modulate:a", 0.0, 0.3)
			tween.tween_callback(func():
				if is_instance_valid(node): node.queue_free()
			)

	if is_instance_valid(_overlay):
		var tween = _overlay.create_tween()
		tween.tween_property(_overlay, "color:a", 0.0, 0.5)
		tween.tween_callback(func():
			if is_instance_valid(_overlay): _overlay.queue_free()
		)
		_overlay = null

	_countdown_label = null
	_mult_label = null
