## MeteorFishPanel.gd — 隕石魚隕石雨 UI 面板（DAY-184）
## 業界依據：Royal Fishing JILI「Dragon Wrath — unleash a massive meteorite attack across
## the centre screen, simultaneously targeting multiple fish」
## 顯示隕石雨開始、每顆隕石落點動畫（☄️ 符號+爆炸圓圈）、隕石計數器、結果彈窗
extends CanvasLayer

# ---- 常數 ----
const PANEL_COLOR_BG      := Color(0.08, 0.04, 0.0, 0.92)
const PANEL_COLOR_ORANGE  := Color(1.0, 0.45, 0.0, 1.0)   # 橙紅（隕石感）
const PANEL_COLOR_FIRE    := Color(1.0, 0.7, 0.1, 1.0)    # 火焰黃
const PANEL_COLOR_WHITE   := Color(1.0, 1.0, 1.0, 1.0)
const PANEL_COLOR_BOSS    := Color(1.0, 0.2, 0.2, 1.0)    # 紅色（命中 BOSS）

# ---- 節點引用 ----
var _banner_container  : Control
var _banner_label      : Label
var _meteor_counter    : Label
var _result_panel      : Control
var _result_label      : Label
var _flash_overlay     : ColorRect

# ---- 狀態 ----
var _meteor_count      : int = 0
var _kill_count        : int = 0
var _total_meteors     : int = 0

func _ready() -> void:
	layer = 61
	_build_ui()
	hide()

func _build_ui() -> void:
	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 0.45, 0.0, 0.0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner_container = PanelContainer.new()
	_banner_container.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner_container.offset_top = 8
	_banner_container.offset_bottom = 56
	_banner_container.offset_left = 80
	_banner_container.offset_right = -80
	var banner_style := StyleBoxFlat.new()
	banner_style.bg_color = Color(0.15, 0.05, 0.0, 0.92)
	banner_style.corner_radius_top_left = 8
	banner_style.corner_radius_top_right = 8
	banner_style.corner_radius_bottom_left = 8
	banner_style.corner_radius_bottom_right = 8
	banner_style.border_width_left = 2
	banner_style.border_width_right = 2
	banner_style.border_width_top = 2
	banner_style.border_width_bottom = 2
	banner_style.border_color = PANEL_COLOR_ORANGE
	_banner_container.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner_container)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_FIRE)
	_banner_label.add_theme_font_size_override("font_size", 20)
	_banner_container.add_child(_banner_label)

	# 隕石計數器（右上角）
	_meteor_counter = Label.new()
	_meteor_counter.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	_meteor_counter.offset_top = 64
	_meteor_counter.offset_right = -16
	_meteor_counter.offset_left = -240
	_meteor_counter.offset_bottom = 96
	_meteor_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	_meteor_counter.add_theme_color_override("font_color", PANEL_COLOR_ORANGE)
	_meteor_counter.add_theme_font_size_override("font_size", 18)
	_meteor_counter.hide()
	add_child(_meteor_counter)

	# 結果彈窗（右側滑入）
	_result_panel = PanelContainer.new()
	_result_panel.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_panel.offset_right = -16
	_result_panel.offset_left = -300
	_result_panel.offset_top = -90
	_result_panel.offset_bottom = 90
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
	result_style.border_color = PANEL_COLOR_ORANGE
	_result_panel.add_theme_stylebox_override("panel", result_style)
	_result_panel.modulate.a = 0.0
	add_child(_result_panel)

	_result_label = Label.new()
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_label.add_theme_color_override("font_color", PANEL_COLOR_WHITE)
	_result_label.add_theme_font_size_override("font_size", 18)
	_result_panel.add_child(_result_label)

## handle_meteor_fish — 處理隕石魚隕石雨訊息
func handle_meteor_fish(payload: Dictionary) -> void:
	var phase : String = payload.get("phase", "")
	if phase == "meteor_start":
		_on_meteor_start(payload)
	elif phase.begins_with("meteor_") and phase != "meteor_result":
		_on_meteor_hit(payload)
	elif phase == "meteor_result":
		_on_meteor_result(payload)

## _on_meteor_start — 隕石雨開始
func _on_meteor_start(payload: Dictionary) -> void:
	_total_meteors = payload.get("meteor_count", 5)
	_meteor_count = 0
	_kill_count = 0
	var killer_name : String = payload.get("killer_name", "")

	show()

	# 全螢幕橙紅閃光（三次，製造「天降神兵」的衝擊感）
	_flash_screen(Color(1.0, 0.45, 0.0, 0.6), 0.15)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(Color(1.0, 0.45, 0.0, 0.4), 0.15)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(Color(1.0, 0.45, 0.0, 0.25), 0.2)

	# 橫幅
	_banner_label.text = "☄️ %s 觸發隕石魚！%d 顆隕石從天而降！" % [killer_name, _total_meteors]
	_banner_container.modulate.a = 0.0
	_banner_container.show()
	var tween := create_tween()
	tween.tween_property(_banner_container, "modulate:a", 1.0, 0.3)

	# 隕石計數器
	_meteor_counter.text = "☄️ 0 / %d 顆 | 0 擊破" % _total_meteors
	_meteor_counter.show()

## _on_meteor_hit — 每顆隕石落點
func _on_meteor_hit(payload: Dictionary) -> void:
	_meteor_count = payload.get("meteor_num", _meteor_count + 1)
	var target_x : float = payload.get("target_x", 0.0)
	var target_y : float = payload.get("target_y", 0.0)
	var is_boss  : bool  = payload.get("is_boss", false)

	# 更新計數器
	_meteor_counter.text = "☄️ %d / %d 顆 | %d 擊破" % [_meteor_count, _total_meteors, _kill_count]

	# 小閃光（BOSS 命中用紅色）
	if is_boss:
		_flash_screen(Color(1.0, 0.2, 0.2, 0.25), 0.12)
	else:
		_flash_screen(Color(1.0, 0.45, 0.0, 0.15), 0.08)

	# 在目標位置顯示隕石符號 + 爆炸圓圈
	_spawn_meteor_at(Vector2(target_x, target_y), is_boss)

## _on_meteor_result — 隕石雨結果
func _on_meteor_result(payload: Dictionary) -> void:
	var total_meteors : int = payload.get("meteor_count", 0)
	var total_kills   : int = payload.get("total_kills", 0)
	var total_reward  : int = payload.get("total_reward", 0)

	_meteor_count = total_meteors
	_kill_count = total_kills
	_meteor_counter.text = "☄️ %d / %d 顆 | %d 擊破" % [total_meteors, total_meteors, total_kills]

	# 結果彈窗
	_result_label.text = "☄️ 隕石雨結束\n隕石：%d 顆\n擊破：%d 個\n獎勵：+%d" % [total_meteors, total_kills, total_reward]
	_result_panel.offset_left = 0
	_result_panel.modulate.a = 1.0
	var tween := create_tween()
	tween.tween_property(_result_panel, "offset_left", -300.0, 0.4).set_trans(Tween.TRANS_BACK)

	# 高擊破數特效
	if total_kills >= 7:
		# 金色強閃光（≥7 擊破）
		_flash_screen(Color(1.0, 0.9, 0.0, 0.7), 0.2)
		await get_tree().create_timer(0.2).timeout
		_flash_screen(Color(1.0, 0.9, 0.0, 0.5), 0.3)
	elif total_kills >= 4:
		# 橙色雙閃光（≥4 擊破）
		_flash_screen(Color(1.0, 0.45, 0.0, 0.5), 0.2)
		await get_tree().create_timer(0.2).timeout
		_flash_screen(Color(1.0, 0.45, 0.0, 0.3), 0.25)

	# 3.5 秒後淡出
	await get_tree().create_timer(3.5).timeout
	var fade := create_tween()
	fade.tween_property(_result_panel, "modulate:a", 0.0, 0.5)
	fade.parallel().tween_property(_banner_container, "modulate:a", 0.0, 0.5)
	await fade.finished
	_meteor_counter.hide()
	_banner_container.hide()
	hide()

## _spawn_meteor_at — 在指定位置顯示隕石符號 + 爆炸圓圈
func _spawn_meteor_at(pos: Vector2, is_boss: bool) -> void:
	# 隕石符號（從上方飛入）
	var lbl := Label.new()
	lbl.text = "☄️"
	lbl.add_theme_font_size_override("font_size", 28)
	lbl.position = Vector2(pos.x - 14, pos.y - 80)  # 從上方 80px 開始
	lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(lbl)

	# 隕石飛入動畫
	var tween := create_tween()
	tween.tween_property(lbl, "position:y", pos.y - 14, 0.18).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
	tween.tween_property(lbl, "scale", Vector2(1.5, 1.5), 0.06)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.2)
	await tween.finished
	if is_instance_valid(lbl):
		lbl.queue_free()

	# 爆炸圓圈（BOSS 命中用紅色，普通用橙色）
	var circle_color := PANEL_COLOR_BOSS if is_boss else PANEL_COLOR_ORANGE
	_spawn_explosion_circle(pos, circle_color)

## _spawn_explosion_circle — 爆炸圓圈動畫
func _spawn_explosion_circle(pos: Vector2, color: Color) -> void:
	# 用多個 ColorRect 模擬爆炸圓圈（4個方向的短線）
	for i in range(4):
		var bar := ColorRect.new()
		bar.color = color
		bar.size = Vector2(4, 20)
		var angle := i * 90.0
		var rad := deg_to_rad(angle)
		var offset := Vector2(cos(rad), sin(rad)) * 12
		bar.position = pos + offset - Vector2(2, 10)
		bar.rotation = rad
		bar.mouse_filter = Control.MOUSE_FILTER_IGNORE
		add_child(bar)
		var tween := create_tween()
		tween.tween_property(bar, "position", pos + offset * 2.5 - Vector2(2, 10), 0.2)
		tween.parallel().tween_property(bar, "modulate:a", 0.0, 0.2)
		await tween.finished
		if is_instance_valid(bar):
			bar.queue_free()

## _flash_screen — 全螢幕閃光
func _flash_screen(color: Color, duration: float) -> void:
	_flash_overlay.color = color
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)
