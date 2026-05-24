## ThunderboltLobsterPanel.gd — 雷霆龍蝦免費射擊面板（DAY-199）
## 業界依據：Royal Fishing JILI「Thunderbolt Lobster — 15 seconds of free play
## followed by automatic shooting from the Thunderbolt Turret.
## Players can earn extra seconds to extend gameplay.」
##
## 視覺設計：
##   - 電黃橙主題（#FFD700 + #FF8C00 + #FFFFFF）
##   - turret_start：全螢幕黃色強閃光 + 頂部橫幅「⚡ 雷霆砲台啟動！」+ 底部時間進度條
##   - turret_shot：目標位置閃光 + 擊破時「⚡ +N 金幣」浮動文字 + 計數器更新
##   - turret_end：結算彈窗（擊破數/總獎勵/延長秒數）+ 4秒後淡出
extends CanvasLayer

var _panel: Control
var _banner: Label
var _kill_counter: Label     # 擊破計數器
var _time_bar: ColorRect     # 底部時間進度條
var _time_bar_bg: ColorRect  # 進度條背景
var _result_popup: Control

var _is_active: bool = false
var _duration: float = 15.0
var _max_duration: float = 30.0
var _elapsed: float = 0.0
var _remaining: float = 15.0
var _kill_count: int = 0

func _ready() -> void:
	layer = 46
	_build_ui()

func _build_ui() -> void:
	_panel = Control.new()
	_panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	_panel.visible = false
	_panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_panel)

	# 頂部橫幅
	_banner = Label.new()
	_banner.text = "⚡ 雷霆砲台啟動！"
	_banner.add_theme_font_size_override("font_size", 24)
	_banner.add_theme_color_override("font_color", Color("#FFD700"))
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.position = Vector2(0, 12)
	_banner.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.visible = false
	add_child(_banner)

	# 擊破計數器
	_kill_counter = Label.new()
	_kill_counter.text = "擊破: 0 個"
	_kill_counter.add_theme_font_size_override("font_size", 18)
	_kill_counter.add_theme_color_override("font_color", Color("#FF8C00"))
	_kill_counter.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_kill_counter.position = Vector2(0, 42)
	_kill_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_kill_counter.visible = false
	add_child(_kill_counter)

	# 底部時間進度條背景
	_time_bar_bg = ColorRect.new()
	_time_bar_bg.color = Color(0.1, 0.1, 0.1, 0.7)
	_time_bar_bg.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	_time_bar_bg.size = Vector2(1280, 14)
	_time_bar_bg.position = Vector2(0, -14)
	_time_bar_bg.visible = false
	add_child(_time_bar_bg)

	# 底部時間進度條
	_time_bar = ColorRect.new()
	_time_bar.color = Color("#FFD700")
	_time_bar.set_anchors_preset(Control.PRESET_BOTTOM_LEFT)
	_time_bar.size = Vector2(1280, 14)
	_time_bar.position = Vector2(0, -14)
	_time_bar.visible = false
	add_child(_time_bar)

	# 結算彈窗
	_result_popup = Control.new()
	_result_popup.visible = false
	_result_popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_popup.size = Vector2(300, 200)
	_result_popup.position = Vector2(-320, -100)
	add_child(_result_popup)

	var popup_bg = ColorRect.new()
	popup_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_bg.color = Color(0.1, 0.08, 0.0, 0.93)
	_result_popup.add_child(popup_bg)

	var popup_label = Label.new()
	popup_label.name = "ResultLabel"
	popup_label.add_theme_font_size_override("font_size", 16)
	popup_label.add_theme_color_override("font_color", Color("#FFD700"))
	popup_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_popup.add_child(popup_label)

## 處理雷霆龍蝦訊息
func handle_thunderbolt_lobster(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"turret_start":
			_on_turret_start(payload)
		"turret_shot":
			_on_turret_shot(payload)
		"turret_end":
			_on_turret_end(payload)

## 砲台啟動
func _on_turret_start(payload: Dictionary) -> void:
	_is_active = true
	_duration = payload.get("duration", 15.0)
	_max_duration = payload.get("max_duration", 30.0)
	_remaining = _duration
	_kill_count = 0

	var killer_name = payload.get("killer_name", "玩家")

	# 全螢幕黃色強閃光（3次）
	_flash_screen(Color("#FFD700"), 0.7)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(Color("#FF8C00"), 0.5)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(Color("#FFD700"), 0.4)

	# 顯示橫幅
	_banner.text = "⚡ " + killer_name + " 觸發雷霆砲台！"
	_banner.visible = true
	_kill_counter.text = "擊破: 0 個"
	_kill_counter.visible = true

	# 顯示進度條
	_time_bar_bg.visible = true
	_time_bar.visible = true
	_time_bar.size.x = 1280.0

	# 啟動進度條更新
	_panel.visible = true

## 自動射擊一次
func _on_turret_shot(payload: Dictionary) -> void:
	if not _is_active:
		return

	var killed = payload.get("killed", false)
	var reward = payload.get("reward", 0)
	var kill_count = payload.get("kill_count", 0)
	var remaining = payload.get("remaining", _remaining)

	# 更新剩餘時間
	_remaining = remaining

	# 更新進度條（顏色依剩餘時間）
	var ratio = _remaining / _max_duration
	_time_bar.size.x = 1280.0 * ratio
	if ratio > 0.5:
		_time_bar.color = Color("#FFD700")  # 黃色
	elif ratio > 0.25:
		_time_bar.color = Color("#FF8C00")  # 橙色
	else:
		_time_bar.color = Color("#FF4500")  # 紅橙色

	if killed:
		_kill_count = kill_count
		_kill_counter.text = "擊破: " + str(kill_count) + " 個"

		# 小閃光
		_flash_screen(Color("#FFD700"), 0.25)

		# 浮動獎勵文字
		if reward > 0:
			_spawn_float_text("⚡ +" + str(reward) + " 金幣", Color("#FFD700"))

		# 每 5 次擊破額外閃光
		if kill_count % 5 == 0 and kill_count > 0:
			_flash_screen(Color("#FF8C00"), 0.45)

## 砲台結束
func _on_turret_end(payload: Dictionary) -> void:
	_is_active = false
	var kill_count = payload.get("kill_count", 0)
	var total_reward = payload.get("total_reward", 0)
	var extend_sec = payload.get("extend_sec", 0.0)

	# 隱藏橫幅和計數器
	_banner.visible = false
	_kill_counter.visible = false
	_time_bar.visible = false
	_time_bar_bg.visible = false

	# 顯示結算彈窗
	var result_label = _result_popup.get_node("ResultLabel")
	var extend_text = ""
	if extend_sec > 0:
		extend_text = "\n⏱ 延長 " + str(snapped(extend_sec, 0.1)) + " 秒"
	result_label.text = (
		"⚡ 雷霆砲台結束！\n\n"
		+ "擊破: " + str(kill_count) + " 個\n"
		+ "總獎勵: " + str(total_reward) + " 金幣"
		+ extend_text
	)
	_result_popup.visible = true

	# 依擊破數決定閃光
	if kill_count >= 20:
		_flash_screen(Color("#FFD700"), 0.8)
	elif kill_count >= 10:
		_flash_screen(Color("#FF8C00"), 0.6)
	elif kill_count >= 5:
		_flash_screen(Color("#FFD700"), 0.4)

	# 4 秒後淡出
	await get_tree().create_timer(4.0).timeout
	var tween = create_tween()
	tween.tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	await tween.finished
	_result_popup.visible = false
	_result_popup.modulate.a = 1.0
	_panel.visible = false

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
func _spawn_float_text(text: String, color: Color) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 20)
	label.add_theme_color_override("font_color", color)
	label.position = Vector2(randf_range(300, 900), randf_range(200, 500))
	add_child(label)
	var tween = create_tween()
	tween.tween_property(label, "position:y", label.position.y - 60, 1.0)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 1.0)
	tween.tween_callback(label.queue_free)
