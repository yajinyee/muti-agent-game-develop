## SerialBombCrabPanel.gd — 連環炸彈蟹面板（DAY-201）
## 業界依據：Royal Fishing JILI「Serial Bomb Crab (70x) — triggers large-scale multiple
## explosions across screen, capturing fish within each explosion range.
## Each bomb creates expanding capture zones for massive multi-target eliminations.」
##
## 視覺設計：
##   - 橙紅爆炸主題（#FF4500 + #FF8C00 + #FFD700）
##   - bomb_start：橙紅色閃光 + 頂部橫幅「💣 連環炸彈蟹！N 顆炸彈！」+ 炸彈計數器
##   - bomb_explode：爆炸圓圈擴散（橙→紅→黑）+ 4方向短線 + 擊破數浮動文字
##   - bomb_result：右側滑入結算彈窗（炸彈數/擊破數/總獎勵）
extends CanvasLayer

var _panel: Control
var _banner: Label
var _bomb_counter: Label
var _result_popup: Control

var _total_bombs: int = 0
var _current_bomb: int = 0

func _ready() -> void:
	layer = 44
	_build_ui()

func _build_ui() -> void:
	_panel = Control.new()
	_panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	_panel.visible = false
	_panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_panel)

	# 頂部橫幅
	_banner = Label.new()
	_banner.text = "💣 連環炸彈蟹！"
	_banner.add_theme_font_size_override("font_size", 24)
	_banner.add_theme_color_override("font_color", Color("#FF4500"))
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.position = Vector2(0, 12)
	_banner.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.visible = false
	add_child(_banner)

	# 炸彈計數器
	_bomb_counter = Label.new()
	_bomb_counter.text = "炸彈: 0 / 0"
	_bomb_counter.add_theme_font_size_override("font_size", 18)
	_bomb_counter.add_theme_color_override("font_color", Color("#FF8C00"))
	_bomb_counter.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_bomb_counter.position = Vector2(0, 42)
	_bomb_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_bomb_counter.visible = false
	add_child(_bomb_counter)

	# 結算彈窗
	_result_popup = Control.new()
	_result_popup.visible = false
	_result_popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_popup.size = Vector2(300, 200)
	_result_popup.position = Vector2(-320, -100)
	add_child(_result_popup)

	var popup_bg = ColorRect.new()
	popup_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_bg.color = Color(0.12, 0.04, 0.0, 0.93)
	_result_popup.add_child(popup_bg)

	var popup_label = Label.new()
	popup_label.name = "ResultLabel"
	popup_label.add_theme_font_size_override("font_size", 16)
	popup_label.add_theme_color_override("font_color", Color("#FF8C00"))
	popup_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_popup.add_child(popup_label)

## 處理連環炸彈蟹訊息
func handle_serial_bomb_crab(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"bomb_start":
			_on_bomb_start(payload)
		"bomb_explode":
			_on_bomb_explode(payload)
		"bomb_result":
			_on_bomb_result(payload)

## 連環爆炸開始
func _on_bomb_start(payload: Dictionary) -> void:
	var killer_name = payload.get("killer_name", "玩家")
	var bomb_count = payload.get("bomb_count", 3)

	_total_bombs = bomb_count
	_current_bomb = 0
	_panel.visible = true

	# 橙紅色閃光
	_flash_screen(Color("#FF4500"), 0.6)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(Color("#FF8C00"), 0.4)

	# 顯示橫幅
	_banner.text = "💣 " + killer_name + " 觸發連環炸彈蟹！" + str(bomb_count) + " 顆炸彈！"
	_banner.visible = true
	_bomb_counter.text = "炸彈: 0 / " + str(bomb_count)
	_bomb_counter.visible = true

## 單顆炸彈爆炸
func _on_bomb_explode(payload: Dictionary) -> void:
	var bomb_index = payload.get("bomb_index", 1)
	var bomb_count = payload.get("bomb_count", 3)
	var bomb_x = payload.get("bomb_x", 640.0)
	var bomb_y = payload.get("bomb_y", 360.0)
	var bomb_kills = payload.get("bomb_kills", 0)
	var bomb_reward = payload.get("bomb_reward", 0)

	_current_bomb = bomb_index
	_bomb_counter.text = "炸彈: %d / %d" % [bomb_index, bomb_count]

	# 爆炸圓圈動畫（在炸彈位置）
	_spawn_explosion(Vector2(bomb_x, bomb_y), bomb_kills)

	# 小閃光
	_flash_screen(Color("#FF4500"), 0.3)

	# 浮動文字
	if bomb_kills > 0:
		_spawn_float_text("💣 ×%d +%d" % [bomb_kills, bomb_reward],
			Color("#FFD700"), Vector2(bomb_x, bomb_y - 40))

## 最終結算
func _on_bomb_result(payload: Dictionary) -> void:
	var total_kills = payload.get("total_kills", 0)
	var total_reward = payload.get("total_reward", 0)

	# 隱藏橫幅和計數器
	_banner.visible = false
	_bomb_counter.visible = false

	# 依擊破數決定閃光
	if total_kills >= 10:
		_flash_screen(Color("#FFD700"), 0.7)
	elif total_kills >= 5:
		_flash_screen(Color("#FF8C00"), 0.5)

	# 顯示結算彈窗
	var result_label = _result_popup.get_node("ResultLabel")
	result_label.text = (
		"💣 連環炸彈蟹結算\n\n"
		+ "炸彈數: " + str(_total_bombs) + " 顆\n"
		+ "擊破: " + str(total_kills) + " 個\n"
		+ "總獎勵: " + str(total_reward) + " 金幣"
	)
	_result_popup.visible = true

	# 4 秒後淡出
	await get_tree().create_timer(4.0).timeout
	var tween = create_tween()
	tween.tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	await tween.finished
	_result_popup.visible = false
	_result_popup.modulate.a = 1.0
	_panel.visible = false

## 爆炸圓圈動畫
func _spawn_explosion(pos: Vector2, kills: int) -> void:
	# 主爆炸圓圈
	var circle = ColorRect.new()
	var size = 20.0
	circle.size = Vector2(size, size)
	circle.position = pos - Vector2(size / 2, size / 2)
	var color = Color("#FF4500") if kills == 0 else Color("#FFD700")
	circle.color = color
	add_child(circle)

	var tween = create_tween()
	tween.tween_property(circle, "size", Vector2(500, 500), 0.5)
	tween.parallel().tween_property(circle, "position", pos - Vector2(250, 250), 0.5)
	tween.parallel().tween_property(circle, "modulate:a", 0.0, 0.5)
	tween.tween_callback(circle.queue_free)

	# 4 方向短線
	for angle in [0, 90, 180, 270]:
		var line = ColorRect.new()
		line.size = Vector2(8, 40)
		line.color = Color("#FF8C00")
		line.position = pos
		line.rotation_degrees = angle
		add_child(line)
		var lt = create_tween()
		lt.tween_property(line, "position", pos + Vector2(cos(deg_to_rad(angle)) * 80, sin(deg_to_rad(angle)) * 80), 0.4)
		lt.parallel().tween_property(line, "modulate:a", 0.0, 0.4)
		lt.tween_callback(line.queue_free)

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
	tween.tween_property(label, "position:y", label.position.y - 60, 1.0)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 1.0)
	tween.tween_callback(label.queue_free)
