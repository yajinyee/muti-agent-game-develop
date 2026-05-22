## HumpbackWhalePanel.gd — 座頭鯨覺醒面板（DAY-203）
## 業界依據：Royal Fishing JILI「Humpback Whale offers 90-150x with 15x base multiplier.
## Awaken Boss mechanic — triggers wave attack that sweeps the screen.
## The Humpback Whale's signature breach mechanic creates massive splash zones.」
##
## 視覺設計：
##   - 深海藍綠主題（#006994 + #00CED1 + #7FFFD4 + #FFFFFF）
##   - awaken_start：深藍色雙閃光 + 頂部橫幅「🐋 座頭鯨覺醒！」+ 基礎獎勵彈窗 + 波浪計數器
##   - wave_attack：波浪掃場動畫（橫向藍綠波浪）+ 擊破浮動文字 + 波數更新
##   - tidal_wave_start：全螢幕深藍白強閃光 + 「🌊 深海巨浪！」大字
##   - tidal_wave_result：巨浪結果浮動文字
##   - awaken_result：右側滑入結算彈窗（基礎獎勵/波浪擊破/巨浪擊破/總獎勵）
extends CanvasLayer

var _panel: Control
var _banner: Label
var _wave_counter: Label
var _result_popup: Control

var _wave_kills: int = 0
var _tidal_kills: int = 0

func _ready() -> void:
	layer = 42
	_build_ui()

func _build_ui() -> void:
	_panel = Control.new()
	_panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	_panel.visible = false
	_panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_panel)

	# 頂部橫幅
	_banner = Label.new()
	_banner.text = "🐋 座頭鯨覺醒！"
	_banner.add_theme_font_size_override("font_size", 24)
	_banner.add_theme_color_override("font_color", Color("#00CED1"))
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.position = Vector2(0, 12)
	_banner.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.visible = false
	add_child(_banner)

	# 波浪計數器
	_wave_counter = Label.new()
	_wave_counter.text = "波浪: 0 / 3"
	_wave_counter.add_theme_font_size_override("font_size", 18)
	_wave_counter.add_theme_color_override("font_color", Color("#7FFFD4"))
	_wave_counter.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_wave_counter.position = Vector2(0, 42)
	_wave_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_wave_counter.visible = false
	add_child(_wave_counter)

	# 結算彈窗
	_result_popup = Control.new()
	_result_popup.visible = false
	_result_popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_popup.size = Vector2(300, 240)
	_result_popup.position = Vector2(-320, -120)
	add_child(_result_popup)

	var popup_bg = ColorRect.new()
	popup_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_bg.color = Color(0.0, 0.05, 0.15, 0.93)
	_result_popup.add_child(popup_bg)

	var popup_label = Label.new()
	popup_label.name = "ResultLabel"
	popup_label.add_theme_font_size_override("font_size", 15)
	popup_label.add_theme_color_override("font_color", Color("#00CED1"))
	popup_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_popup.add_child(popup_label)

## 處理座頭鯨覺醒訊息
func handle_humpback_whale(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"awaken_start":
			_on_awaken_start(payload)
		"wave_attack":
			_on_wave_attack(payload)
		"tidal_wave_start":
			_on_tidal_wave_start()
		"tidal_wave_result":
			_on_tidal_wave_result(payload)
		"awaken_result":
			_on_awaken_result(payload)

## 覺醒開始
func _on_awaken_start(payload: Dictionary) -> void:
	var killer_name = payload.get("killer_name", "玩家")
	var base_reward = payload.get("base_reward", 0)
	var wave_count = payload.get("wave_count", 3)

	_wave_kills = 0
	_tidal_kills = 0
	_panel.visible = true

	# 深藍色雙閃光
	_flash_screen(Color("#006994"), 0.7)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color("#00CED1"), 0.5)

	# 顯示橫幅
	_banner.text = "🐋 " + killer_name + " 觸發座頭鯨覺醒！"
	_banner.visible = true
	_wave_counter.text = "波浪: 0 / %d" % wave_count
	_wave_counter.visible = true

	# 基礎獎勵彈窗（短暫顯示）
	_spawn_float_text("🐋 基礎獎勵 +%d" % base_reward, Color("#7FFFD4"),
		Vector2(640, 300))

## 波浪攻擊
func _on_wave_attack(payload: Dictionary) -> void:
	var wave_num = payload.get("wave_num", 1)
	var wave_kills = payload.get("wave_kills", 0)
	var wave_reward = payload.get("wave_reward", 0)
	var total_kills = payload.get("total_kills", 0)

	_wave_kills = total_kills
	_wave_counter.text = "波浪: %d / 3" % wave_num

	# 波浪掃場動畫（橫向藍綠波浪）
	_spawn_wave_animation(wave_num)

	# 擊破浮動文字
	if wave_kills > 0:
		_spawn_float_text("🌊 第%d波 ×%d +%d" % [wave_num, wave_kills, wave_reward],
			Color("#00CED1"),
			Vector2(randf_range(300, 900), randf_range(200, 500)))

## 深海巨浪開始
func _on_tidal_wave_start() -> void:
	# 全螢幕深藍白強閃光
	_flash_screen(Color("#FFFFFF"), 0.9)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(Color("#00CED1"), 0.7)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(Color("#006994"), 0.5)

	# 「深海巨浪！」大字
	var big_label = Label.new()
	big_label.text = "🌊 深海巨浪！"
	big_label.add_theme_font_size_override("font_size", 48)
	big_label.add_theme_color_override("font_color", Color("#7FFFD4"))
	big_label.set_anchors_preset(Control.PRESET_CENTER)
	big_label.position = Vector2(-200, -40)
	add_child(big_label)

	var tween = create_tween()
	tween.tween_property(big_label, "scale", Vector2(1.3, 1.3), 0.3)
	tween.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.2)
	tween.tween_interval(1.0)
	tween.tween_property(big_label, "modulate:a", 0.0, 0.5)
	tween.tween_callback(big_label.queue_free)

	# 更新橫幅
	_banner.text = "🌊 深海巨浪！全場清掃！"
	_banner.add_theme_color_override("font_color", Color("#7FFFD4"))

## 深海巨浪結果
func _on_tidal_wave_result(payload: Dictionary) -> void:
	var tidal_kills = payload.get("tidal_kills", 0)
	var tidal_reward = payload.get("tidal_reward", 0)
	_tidal_kills = tidal_kills

	if tidal_kills > 0:
		_spawn_float_text("🌊💥 巨浪擊破 %d 個！+%d" % [tidal_kills, tidal_reward],
			Color("#FFD700"),
			Vector2(640, 250))

## 最終結算
func _on_awaken_result(payload: Dictionary) -> void:
	var base_reward = payload.get("base_reward", 0)
	var total_kills = payload.get("total_kills", 0)
	var total_reward = payload.get("total_reward", 0)
	var has_tidal = payload.get("has_tidal", false)

	# 隱藏橫幅和計數器
	_banner.visible = false
	_wave_counter.visible = false

	# 依總獎勵決定閃光
	if total_reward >= 100:
		_flash_screen(Color("#FFD700"), 0.6)
	elif total_reward >= 50:
		_flash_screen(Color("#00CED1"), 0.4)

	# 顯示結算彈窗
	var result_label = _result_popup.get_node("ResultLabel")
	var tidal_line = ""
	if has_tidal:
		tidal_line = "\n🌊 深海巨浪: %d 個" % _tidal_kills
	result_label.text = (
		"🐋 座頭鯨覺醒結算\n\n"
		+ "基礎獎勵: %d\n" % base_reward
		+ "波浪擊破: %d 個\n" % (_wave_kills - _tidal_kills)
		+ tidal_line
		+ "\n總擊破: %d 個\n" % total_kills
		+ "總獎勵: %d 金幣" % total_reward
	)
	_result_popup.visible = true
	_result_popup.position = Vector2(get_viewport().get_visible_rect().size.x + 10, -120)

	# 從右側滑入
	var tween = create_tween()
	tween.tween_property(_result_popup, "position:x",
		get_viewport().get_visible_rect().size.x - 320, 0.4)

	# 4.5 秒後淡出
	await get_tree().create_timer(4.5).timeout
	var fade_tween = create_tween()
	fade_tween.tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	await fade_tween.finished
	_result_popup.visible = false
	_result_popup.modulate.a = 1.0
	_panel.visible = false

## 波浪掃場動畫（橫向藍綠波浪）
func _spawn_wave_animation(wave_num: int) -> void:
	var wave_y = 150.0 + (wave_num - 1) * 150.0  # 每波在不同高度
	var wave = ColorRect.new()
	wave.size = Vector2(0, 30)
	wave.position = Vector2(0, wave_y)
	wave.color = Color(0.0, 0.8, 0.82, 0.6)  # 青色半透明
	add_child(wave)

	var tween = create_tween()
	tween.tween_property(wave, "size:x", 1280.0, 0.6)
	tween.parallel().tween_property(wave, "modulate:a", 0.0, 0.8)
	tween.tween_callback(wave.queue_free)

## 全螢幕閃光
func _flash_screen(color: Color, alpha: float) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, alpha)
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.25)
	tween.tween_callback(flash.queue_free)

## 浮動文字
func _spawn_float_text(text: String, color: Color, pos: Vector2 = Vector2(-1, -1)) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 20)
	label.add_theme_color_override("font_color", color)
	if pos.x < 0:
		label.position = Vector2(randf_range(300, 900), randf_range(200, 500))
	else:
		label.position = pos
	add_child(label)
	var tween = create_tween()
	tween.tween_property(label, "position:y", label.position.y - 70, 1.2)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 1.2)
	tween.tween_callback(label.queue_free)
