## GoldenSharkPanel.gd — 黃金鯊魚全服狂暴模式面板（DAY-161）
## 擊破黃金鯊魚後觸發全服狂暴模式 12 秒，所有獎勵 ×1.5
## 視覺：橙紅鯊魚爆炸 + 全螢幕橙色光暈 + 倒數計時器 + ×1.5 標誌
## 與幸運星魚不同：全服共享，任何玩家擊破都讓全服受益
extends Node2D

var _pixel_font: Font = null
var _countdown_label: Label = null
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
	if GameManager.has_signal("golden_shark_berserk"):
		GameManager.golden_shark_berserk.connect(_on_golden_shark_berserk)

func _process(delta: float) -> void:
	if not _is_active:
		return
	_remaining -= delta
	if _remaining <= 0.0:
		_remaining = 0.0
		_is_active = false
		_hide_berserk()
		return
	if is_instance_valid(_countdown_label):
		_countdown_label.text = "🦈 ×1.5 %.1f秒" % _remaining

## 處理黃金鯊魚狂暴模式事件
func _on_golden_shark_berserk(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var duration: float = data.get("duration_secs", 12.0)
	var killer_name: String = data.get("killer_name", "")
	var killer_id: String = data.get("killer_id", "")

	if phase == "berserk_start":
		_show_berserk(killer_name, killer_id, duration)
	elif phase == "berserk_end":
		_hide_berserk()

## 顯示狂暴模式效果
func _show_berserk(killer_name: String, killer_id: String, duration: float) -> void:
	_is_active = true
	_remaining = duration
	_my_player_id = NetworkManager.get_player_id() if NetworkManager.has_method("get_player_id") else ""

	var is_me = (killer_id == _my_player_id)

	# 全螢幕橙紅半透明遮罩（狂暴感）
	_overlay = ColorRect.new()
	_overlay.position = Vector2(-640, -360)
	_overlay.size = Vector2(1280, 720)
	_overlay.color = Color(1.0, 0.4, 0.0, 0.07)
	add_child(_overlay)

	# 全螢幕橙色閃光
	var flash := ColorRect.new()
	flash.position = Vector2(-640, -360)
	flash.size = Vector2(1280, 720)
	flash.color = Color(1.0, 0.5, 0.0, 0.55 if is_me else 0.35)
	add_child(flash)
	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", 0.0, 0.5)
	flash_tween.tween_callback(func():
		if is_instance_valid(flash): flash.queue_free()
	)

	# 頂部橫幅（橙紅主題）
	var banner_bg := ColorRect.new()
	banner_bg.name = "SharkBanner"
	banner_bg.position = Vector2(-640, -360)
	banner_bg.size = Vector2(1280, 52)
	banner_bg.color = Color(0.18, 0.06, 0.0, 0.92)
	add_child(banner_bg)

	var banner_text: String
	if is_me:
		banner_text = "🦈 你擊破黃金鯊魚！全服狂暴模式！所有獎勵 ×1.5 持續 %.0f 秒！" % duration
	else:
		banner_text = "🦈 %s 擊破黃金鯊魚！全服狂暴模式！所有獎勵 ×1.5 持續 %.0f 秒！" % [killer_name, duration]

	var banner_lbl := Label.new()
	banner_lbl.name = "SharkBannerLabel"
	banner_lbl.position = Vector2(-640, -354)
	banner_lbl.size = Vector2(1280, 44)
	banner_lbl.text = banner_text
	banner_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	banner_lbl.add_theme_color_override("font_color", Color(1.0, 0.7, 0.1))
	banner_lbl.add_theme_font_size_override("font_size", 16)
	if _pixel_font:
		banner_lbl.add_theme_font_override("font", _pixel_font)
	add_child(banner_lbl)

	# 橫幅滑入動畫
	banner_bg.position.y = -410
	banner_lbl.position.y = -404
	var slide_tween = banner_bg.create_tween()
	slide_tween.tween_property(banner_bg, "position:y", -360.0, 0.3)
	var slide_tween2 = banner_lbl.create_tween()
	slide_tween2.tween_property(banner_lbl, "position:y", -354.0, 0.3)

	# 倒數計時器（右上角）
	var countdown_bg := ColorRect.new()
	countdown_bg.name = "SharkCountdownBG"
	countdown_bg.position = Vector2(510, -355)
	countdown_bg.size = Vector2(130, 36)
	countdown_bg.color = Color(0.15, 0.05, 0.0, 0.9)
	add_child(countdown_bg)

	_countdown_label = Label.new()
	_countdown_label.name = "SharkCountdown"
	_countdown_label.position = Vector2(510, -352)
	_countdown_label.size = Vector2(130, 30)
	_countdown_label.text = "🦈 ×1.5 %.1f秒" % duration
	_countdown_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_countdown_label.add_theme_color_override("font_color", Color(1.0, 0.7, 0.1))
	_countdown_label.add_theme_font_size_override("font_size", 13)
	if _pixel_font:
		_countdown_label.add_theme_font_override("font", _pixel_font)
	add_child(_countdown_label)

	# 自己觸發時：中央大 ×1.5 標誌（彈跳動畫）
	if is_me:
		var mult_bg := ColorRect.new()
		mult_bg.name = "BerserkMultBG"
		mult_bg.position = Vector2(-90, -55)
		mult_bg.size = Vector2(180, 90)
		mult_bg.color = Color(0.18, 0.06, 0.0, 0.95)
		add_child(mult_bg)

		var mult_lbl := Label.new()
		mult_lbl.name = "BerserkMultLabel"
		mult_lbl.position = Vector2(-90, -53)
		mult_lbl.size = Vector2(180, 86)
		mult_lbl.text = "×1.5"
		mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		mult_lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
		mult_lbl.add_theme_color_override("font_color", Color(1.0, 0.7, 0.1))
		mult_lbl.add_theme_font_size_override("font_size", 44)
		if _pixel_font:
			mult_lbl.add_theme_font_override("font", _pixel_font)
		add_child(mult_lbl)

		# 彈跳動畫
		mult_lbl.scale = Vector2(0.5, 0.5)
		var scale_tween = mult_lbl.create_tween()
		scale_tween.tween_property(mult_lbl, "scale", Vector2(1.2, 1.2), 0.2)
		scale_tween.tween_property(mult_lbl, "scale", Vector2(1.0, 1.0), 0.1)

		# 2 秒後淡出
		var fade_tween = mult_lbl.create_tween()
		fade_tween.tween_interval(2.0)
		fade_tween.tween_property(mult_lbl, "modulate:a", 0.0, 0.5)
		fade_tween.tween_callback(func():
			if is_instance_valid(mult_lbl): mult_lbl.queue_free()
			if is_instance_valid(mult_bg): mult_bg.queue_free()
		)

## 隱藏狂暴模式效果
func _hide_berserk() -> void:
	_is_active = false
	_remaining = 0.0

	for child_name in ["SharkBanner", "SharkBannerLabel", "SharkCountdownBG", "SharkCountdown", "BerserkMultBG", "BerserkMultLabel"]:
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
