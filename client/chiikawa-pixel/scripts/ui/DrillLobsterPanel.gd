## DrillLobsterPanel.gd — 鑽頭龍蝦穿透爆炸面板（DAY-195）
## 業界依據：Royal Fishing JILI「Drill Bit Lobster (80X) — fires a penetrating drill that passes
## through multiple fish before self-detonating, capturing everything in the explosion radius.」
##
## 視覺設計：
##   - 機械橙色主題（#FF6600 + #CC4400）
##   - drill_start：橙色雙閃光 + 頂部橫幅 + 鑽頭方向箭頭
##   - drill_N：鑽頭移動軌跡線 + 命中爆炸圓圈 + 穿透計數器
##   - drill_explode：大型橙紅色爆炸圓圈（300px 半徑）+ 強閃光
##   - drill_result：右側滑入結算彈窗（穿透數/擊破數/獎勵）
extends CanvasLayer

var _panel: Control
var _drill_trail: Control      # 鑽頭軌跡繪製節點
var _counter_label: Label      # 穿透計數器
var _result_popup: Control
var _banner: Label

# 軌跡點列表
var _trail_points: Array = []
var _hit_points: Array = []    # 命中位置（用於繪製爆炸圓圈）

func _ready() -> void:
	layer = 50
	_build_ui()

func _build_ui() -> void:
	_panel = Control.new()
	_panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	_panel.visible = false
	_panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_panel)

	# 頂部橫幅
	_banner = Label.new()
	_banner.text = "🦞 鑽頭龍蝦穿透爆炸！"
	_banner.add_theme_font_size_override("font_size", 22)
	_banner.add_theme_color_override("font_color", Color("#FF6600"))
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.position = Vector2(0, 15)
	_banner.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.visible = false
	add_child(_banner)

	# 穿透計數器
	_counter_label = Label.new()
	_counter_label.text = "穿透: 0 | 擊破: 0"
	_counter_label.add_theme_font_size_override("font_size", 18)
	_counter_label.add_theme_color_override("font_color", Color("#FF6600"))
	_counter_label.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_counter_label.position = Vector2(0, 45)
	_counter_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_counter_label.visible = false
	add_child(_counter_label)

	# 結算彈窗
	_result_popup = Control.new()
	_result_popup.visible = false
	_result_popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_popup.size = Vector2(280, 180)
	_result_popup.position = Vector2(-300, -90)
	add_child(_result_popup)

	var popup_bg = ColorRect.new()
	popup_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_bg.color = Color(0.15, 0.06, 0.0, 0.92)
	_result_popup.add_child(popup_bg)

	var popup_label = Label.new()
	popup_label.name = "ResultLabel"
	popup_label.add_theme_font_size_override("font_size", 17)
	popup_label.add_theme_color_override("font_color", Color("#FF6600"))
	popup_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_popup.add_child(popup_label)

## 處理鑽頭龍蝦訊息
func handle_drill_lobster(payload: Dictionary) -> void:
	var phase = payload.get("phase", "")
	if phase == "drill_start":
		_on_drill_start(payload)
	elif phase.begins_with("drill_") and phase != "drill_explode" and phase != "drill_result":
		_on_drill_step(payload)
	elif phase == "drill_explode":
		_on_drill_explode(payload)
	elif phase == "drill_result":
		_on_drill_result(payload)

func _on_drill_start(payload: Dictionary) -> void:
	_trail_points.clear()
	_hit_points.clear()
	_panel.visible = true
	_banner.visible = true
	_counter_label.visible = true
	_counter_label.text = "穿透: 0 | 擊破: 0"

	var start_x = payload.get("start_x", 400.0)
	var start_y = payload.get("start_y", 300.0)
	_trail_points.append(Vector2(start_x, start_y))

	# 橙色雙閃光
	_flash_screen(Color("#FF6600"), 0.35)
	await get_tree().create_timer(0.4).timeout
	_flash_screen(Color("#FF6600"), 0.25)

func _on_drill_step(payload: Dictionary) -> void:
	var cur_x = payload.get("cur_x", 0.0)
	var cur_y = payload.get("cur_y", 0.0)
	var is_kill = payload.get("is_kill", false)
	var total_kills = payload.get("total_kills", 0)
	var step_index = payload.get("step_index", 0)
	var step_reward = payload.get("step_reward", 0)

	_trail_points.append(Vector2(cur_x, cur_y))

	# 更新計數器
	_counter_label.text = "穿透: %d | 擊破: %d" % [step_index, total_kills]

	if is_kill:
		_hit_points.append(Vector2(cur_x, cur_y))
		# 命中爆炸圓圈
		_spawn_hit_effect(cur_x, cur_y, false)
		# 獎勵浮動文字
		if step_reward > 0:
			_spawn_reward_text(cur_x, cur_y, step_reward)

func _on_drill_explode(payload: Dictionary) -> void:
	var cur_x = payload.get("cur_x", 400.0)
	var cur_y = payload.get("cur_y", 300.0)
	var explode_kills = payload.get("explode_kills", 0)
	var explode_reward = payload.get("explode_reward", 0)

	# 大型爆炸圓圈（300px 半徑）
	_spawn_hit_effect(cur_x, cur_y, true)

	# 橙紅色強閃光
	_flash_screen(Color("#FF4500"), 0.5)
	await get_tree().create_timer(0.55).timeout
	_flash_screen(Color("#FF6600"), 0.35)

	if explode_kills > 0:
		_spawn_reward_text(cur_x, cur_y - 40, explode_reward)

func _on_drill_result(payload: Dictionary) -> void:
	var penetrate_count = payload.get("penetrate_count", 0)
	var total_kills = payload.get("total_kills", 0)
	var total_reward = payload.get("total_reward", 0)
	var killer_name = payload.get("killer_name", "")

	# 顯示結算彈窗（右側滑入）
	_result_popup.visible = true
	_result_popup.modulate.a = 0.0
	_result_popup.position.x = get_viewport().get_visible_rect().size.x

	var label = _result_popup.get_node("ResultLabel")
	label.text = "🦞 鑽頭龍蝦爆炸\n穿透: %d 個\n擊破: %d 個\n獎勵: %d 金幣" % [penetrate_count, total_kills, total_reward]

	var tween = create_tween()
	tween.tween_property(_result_popup, "modulate:a", 1.0, 0.3)
	tween.parallel().tween_property(_result_popup, "position:x",
		get_viewport().get_visible_rect().size.x - 300, 0.3)

	# 依擊破數決定特效
	if total_kills >= 7:
		_flash_screen(Color("#FFD700"), 0.4)
	elif total_kills >= 4:
		_flash_screen(Color("#FF6600"), 0.3)

	# 4 秒後淡出
	await get_tree().create_timer(4.0).timeout
	_fade_out()

func _spawn_hit_effect(x: float, y: float, is_big: bool) -> void:
	var effect = Control.new()
	effect.position = Vector2(x, y)
	effect.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_panel.add_child(effect)

	var radius = 300.0 if is_big else 60.0
	var color = Color("#FF4500") if is_big else Color("#FF6600")

	# 用 ColorRect 模擬爆炸圓圈（簡化版）
	var circle = ColorRect.new()
	circle.size = Vector2(radius * 2, radius * 2)
	circle.position = Vector2(-radius, -radius)
	circle.color = Color(color.r, color.g, color.b, 0.5)
	circle.mouse_filter = Control.MOUSE_FILTER_IGNORE
	effect.add_child(circle)

	var tween = create_tween()
	tween.tween_property(circle, "scale", Vector2(1.5, 1.5), 0.4)
	tween.parallel().tween_property(circle, "modulate:a", 0.0, 0.4)
	await tween.finished
	effect.queue_free()

func _spawn_reward_text(x: float, y: float, reward: int) -> void:
	var label = Label.new()
	label.text = "+%d" % reward
	label.add_theme_font_size_override("font_size", 20)
	label.add_theme_color_override("font_color", Color("#FF6600"))
	label.position = Vector2(x - 30, y - 20)
	label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_panel.add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "position:y", y - 60, 0.8)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.8)
	await tween.finished
	label.queue_free()

func _fade_out() -> void:
	var tween = create_tween()
	tween.tween_property(_panel, "modulate:a", 0.0, 0.4)
	tween.parallel().tween_property(_result_popup, "modulate:a", 0.0, 0.4)
	await tween.finished
	_panel.visible = false
	_panel.modulate.a = 1.0
	_result_popup.visible = false
	_banner.visible = false
	_counter_label.visible = false
	_trail_points.clear()
	_hit_points.clear()

func _flash_screen(color: Color, intensity: float) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, intensity)
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.25)
	await tween.finished
	flash.queue_free()
