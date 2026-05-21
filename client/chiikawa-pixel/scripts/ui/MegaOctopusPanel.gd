## MegaOctopusPanel.gd
## 巨型章魚轉盤面板（DAY-144）
## 業界依據：JILI Mega Fishing「Mega Octopus Wheel – Defeat that giant octopus and enter
## the bonus wheel round where you have a chance to win massive guaranteed prizes up to 950x.」
## 擊破 T108 後觸發個人轉盤，8格獎勵（50x-950x），玩家點擊停止

extends Control

var pixel_font: Font = null

# 轉盤狀態
var _is_spinning: bool = false
var _spin_duration: float = 3.0
var _spin_elapsed: float = 0.0
var _current_angle: float = 0.0
var _spin_speed: float = 0.0  # 度/秒
var _slots: Array = []
var _overlay: Control = null
var _wheel_container: Control = null
var _stop_btn: Button = null
var _result_shown: bool = false

# 轉盤顏色
const COLOR_PURPLE = Color(0.58, 0.0, 0.83)
const COLOR_GOLD   = Color(1.0, 0.85, 0.0)
const COLOR_WHITE  = Color(1.0, 1.0, 1.0)

## 初始化（由 HUD.gd 呼叫）
func setup(font: Font) -> void:
	pixel_font = font
	GameManager.mega_octopus_wheel_start.connect(_on_wheel_start)
	GameManager.mega_octopus_wheel_result.connect(_on_wheel_result)

## 轉盤開始（Server 通知）
func _on_wheel_start(data: Dictionary) -> void:
	if _is_spinning:
		return

	_slots = data.get("slots", [])
	_spin_duration = float(data.get("spin_duration", 3))
	_is_spinning = true
	_spin_elapsed = 0.0
	_result_shown = false
	_spin_speed = 720.0  # 初始速度 720度/秒（2圈/秒）

	_build_wheel_overlay()

## 建立轉盤 overlay
func _build_wheel_overlay() -> void:
	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# 全螢幕遮罩
	_overlay = Control.new()
	_overlay.name = "MegaOctopusOverlay"
	_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_overlay.z_index = 90
	canvas_layer.add_child(_overlay)

	# 半透明黑色背景
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.0, 0.0, 0.0, 0.0)
	_overlay.add_child(bg)

	# 淡入背景
	var bg_tween = bg.create_tween()
	bg_tween.tween_property(bg, "color:a", 0.75, 0.3)

	# 中央面板
	_wheel_container = Control.new()
	_wheel_container.name = "WheelContainer"
	_wheel_container.position = Vector2(390, 140)
	_wheel_container.size = Vector2(500, 440)
	_overlay.add_child(_wheel_container)

	# 面板背景
	var panel_bg = ColorRect.new()
	panel_bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	panel_bg.color = Color(0.05, 0.02, 0.12, 0.95)
	_wheel_container.add_child(panel_bg)

	# 紫色邊框
	for border_data in [
		[Vector2(0, 0), Vector2(500, 3)],
		[Vector2(0, 437), Vector2(500, 3)],
		[Vector2(0, 0), Vector2(3, 440)],
		[Vector2(497, 0), Vector2(3, 440)],
	]:
		var border = ColorRect.new()
		border.position = border_data[0]
		border.size = border_data[1]
		border.color = COLOR_PURPLE
		_wheel_container.add_child(border)

	# 標題
	var title = Label.new()
	title.text = "🐙 巨型章魚轉盤！"
	title.position = Vector2(0, 12)
	title.size = Vector2(500, 36)
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_font_size_override("font_size", 22)
	title.add_theme_color_override("font_color", COLOR_PURPLE)
	if is_instance_valid(pixel_font):
		title.add_theme_font_override("font", pixel_font)
	_wheel_container.add_child(title)

	# 轉盤（8格圓形排列）
	_build_wheel_slots()

	# 指針（頂部三角形）
	var pointer = ColorRect.new()
	pointer.size = Vector2(16, 24)
	pointer.position = Vector2(242, 52)
	pointer.color = COLOR_GOLD
	_wheel_container.add_child(pointer)

	# 倒數計時標籤
	var timer_lbl = Label.new()
	timer_lbl.name = "TimerLabel"
	timer_lbl.text = "%.1f" % _spin_duration
	timer_lbl.position = Vector2(0, 340)
	timer_lbl.size = Vector2(500, 30)
	timer_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	timer_lbl.add_theme_font_size_override("font_size", 18)
	timer_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	if is_instance_valid(pixel_font):
		timer_lbl.add_theme_font_override("font", pixel_font)
	_wheel_container.add_child(timer_lbl)

	# 停止按鈕
	_stop_btn = Button.new()
	_stop_btn.name = "StopBtn"
	_stop_btn.text = "⏹ 停止！"
	_stop_btn.position = Vector2(150, 380)
	_stop_btn.size = Vector2(200, 44)
	_stop_btn.add_theme_font_size_override("font_size", 16)
	_stop_btn.add_theme_color_override("font_color", COLOR_GOLD)
	if is_instance_valid(pixel_font):
		_stop_btn.add_theme_font_override("font", pixel_font)
	_stop_btn.pressed.connect(_on_stop_pressed)
	_wheel_container.add_child(_stop_btn)

	# 彈入動畫
	_wheel_container.position.y = 600
	var panel_tween = _wheel_container.create_tween()
	panel_tween.tween_property(_wheel_container, "position:y", 140.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

## 建立轉盤格子（8格圓形排列）
func _build_wheel_slots() -> void:
	if not is_instance_valid(_wheel_container):
		return

	var center = Vector2(250, 200)
	var radius = 130.0
	var slot_count = 8

	for i in slot_count:
		var angle = (float(i) / slot_count) * TAU - PI / 2.0
		var slot_x = center.x + cos(angle) * radius - 40
		var slot_y = center.y + sin(angle) * radius - 20

		var slot_data = {}
		if i < _slots.size():
			slot_data = _slots[i]

		var slot_bg = ColorRect.new()
		slot_bg.name = "Slot_%d" % i
		slot_bg.position = Vector2(slot_x, slot_y)
		slot_bg.size = Vector2(80, 40)
		var slot_color_hex = slot_data.get("color", "#C0C0C0")
		slot_bg.color = Color(slot_color_hex) * 0.3 + Color(0.05, 0.02, 0.12) * 0.7
		_wheel_container.add_child(slot_bg)

		var slot_lbl = Label.new()
		slot_lbl.text = slot_data.get("label", "%dx" % (i * 100 + 50))
		slot_lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		slot_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		slot_lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
		slot_lbl.add_theme_font_size_override("font_size", 11)
		slot_lbl.add_theme_color_override("font_color", Color(slot_color_hex))
		if is_instance_valid(pixel_font):
			slot_lbl.add_theme_font_override("font", pixel_font)
		slot_bg.add_child(slot_lbl)

## 每幀更新（轉盤旋轉動畫）
func _process(delta: float) -> void:
	if not _is_spinning or not is_instance_valid(_wheel_container):
		return

	_spin_elapsed += delta

	# 更新倒數計時
	var remaining = max(0.0, _spin_duration - _spin_elapsed)
	var timer_lbl = _wheel_container.get_node_or_null("TimerLabel")
	if is_instance_valid(timer_lbl):
		timer_lbl.text = "%.1f 秒後自動停止" % remaining
		if remaining < 1.0:
			timer_lbl.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))

	# 轉盤旋轉（格子繞中心旋轉）
	_current_angle += _spin_speed * delta
	if _current_angle >= 360.0:
		_current_angle -= 360.0

	# 更新格子位置
	var center = Vector2(250, 200)
	var radius = 130.0
	var slot_count = 8
	for i in slot_count:
		var base_angle = (float(i) / slot_count) * TAU - PI / 2.0
		var current_rad = base_angle + deg_to_rad(_current_angle)
		var slot_x = center.x + cos(current_rad) * radius - 40
		var slot_y = center.y + sin(current_rad) * radius - 20

		var slot_node = _wheel_container.get_node_or_null("Slot_%d" % i)
		if is_instance_valid(slot_node):
			slot_node.position = Vector2(slot_x, slot_y)

## 玩家點擊停止按鈕
func _on_stop_pressed() -> void:
	if not _is_spinning:
		return
	if is_instance_valid(_stop_btn):
		_stop_btn.disabled = true
	# 發送停止訊號給 Server
	NetworkManager.send_message("mega_octopus_wheel_stop", {})

## 轉盤結果（Server 通知）
func _on_wheel_result(data: Dictionary) -> void:
	if _result_shown:
		return
	_result_shown = true
	_is_spinning = false

	var result_index = data.get("result_index", 0)
	var multiplier = data.get("multiplier", 50)
	var reward = data.get("reward", 0)
	var slot_label = data.get("slot_label", "%dx" % multiplier)
	var slot_color_hex = data.get("slot_color", "#C0C0C0")

	# 停止轉盤動畫
	if is_instance_valid(_stop_btn):
		_stop_btn.disabled = true

	# 高亮結果格子
	if is_instance_valid(_wheel_container):
		var slot_node = _wheel_container.get_node_or_null("Slot_%d" % result_index)
		if is_instance_valid(slot_node):
			var highlight_tween = slot_node.create_tween().set_loops(4)
			highlight_tween.tween_property(slot_node, "color", Color(slot_color_hex), 0.15)
			highlight_tween.tween_property(slot_node, "color", Color(slot_color_hex) * 0.3, 0.15)

	# 播放音效
	if AudioManager != null:
		AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)

	# 全螢幕閃光（≥300x 紫色，其他金色）
	var flash_color = COLOR_PURPLE if multiplier >= 300 else COLOR_GOLD
	_spawn_result_flash(flash_color)

	# 顯示結果
	await get_tree().create_timer(0.8).timeout
	_show_result_popup(slot_label, reward, multiplier, slot_color_hex)

## 全螢幕閃光
func _spawn_result_flash(color: Color) -> void:
	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return
	var flash = ColorRect.new()
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, 0.0)
	flash.z_index = 91
	canvas_layer.add_child(flash)
	var tween = flash.create_tween()
	tween.tween_property(flash, "color:a", 0.5, 0.1)
	tween.tween_property(flash, "color:a", 0.0, 0.4)
	tween.tween_callback(flash.queue_free)

## 顯示結果彈窗
func _show_result_popup(slot_label: String, reward: int, multiplier: int, color_hex: String) -> void:
	if not is_instance_valid(_wheel_container):
		return

	# 清除轉盤，顯示結果
	for child in _wheel_container.get_children():
		if child.name != "WheelContainer":
			child.queue_free()

	# 結果標題
	var result_title = Label.new()
	result_title.text = "🎉 %s" % slot_label
	result_title.position = Vector2(0, 80)
	result_title.size = Vector2(500, 80)
	result_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	result_title.add_theme_font_size_override("font_size", 52)
	result_title.add_theme_color_override("font_color", Color(color_hex))
	if is_instance_valid(pixel_font):
		result_title.add_theme_font_override("font", pixel_font)
	_wheel_container.add_child(result_title)

	# 獎勵金額
	var reward_lbl = Label.new()
	reward_lbl.text = "🪙 +%d" % reward
	reward_lbl.position = Vector2(0, 180)
	reward_lbl.size = Vector2(500, 60)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_lbl.add_theme_font_size_override("font_size", 36)
	reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	if is_instance_valid(pixel_font):
		reward_lbl.add_theme_font_override("font", pixel_font)
	_wheel_container.add_child(reward_lbl)

	# 彈入動畫
	result_title.modulate.a = 0.0
	result_title.position.y = 120
	var tween = result_title.create_tween()
	tween.tween_property(result_title, "position:y", 80.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(result_title, "modulate:a", 1.0, 0.3)

	reward_lbl.modulate.a = 0.0
	tween.tween_property(reward_lbl, "modulate:a", 1.0, 0.3)

	# 3 秒後關閉
	tween.tween_interval(3.0)
	tween.tween_callback(_close_overlay)

## 關閉 overlay
func _close_overlay() -> void:
	_is_spinning = false
	if is_instance_valid(_overlay):
		var tween = _overlay.create_tween()
		tween.tween_property(_overlay, "modulate:a", 0.0, 0.4)
		tween.tween_callback(_overlay.queue_free)
	_overlay = null
	_wheel_container = null
	_stop_btn = null
