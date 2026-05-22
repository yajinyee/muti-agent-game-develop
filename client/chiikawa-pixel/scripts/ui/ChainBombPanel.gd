## ChainBombPanel.gd — 連鎖爆炸魚 UI 面板（DAY-187）
## 業界依據：Royal Fishing「chain reaction mechanic — players can trigger multiple explosions
## to capture additional fish within a blast radius」
## 顯示連鎖爆炸開始橫幅、每層爆炸圓圈動畫、結果彈窗
## Phase: chain_start → chain_explode(×N) → chain_result
extends CanvasLayer

# ---- 常數 ----
const PANEL_COLOR_BG     := Color(0.15, 0.02, 0.0, 0.92)  # 深紅（爆炸感）
const PANEL_COLOR_RED    := Color(1.0, 0.2, 0.1, 1.0)     # 亮紅（爆炸）
const PANEL_COLOR_ORANGE := Color(1.0, 0.55, 0.0, 1.0)    # 橙色（3層連鎖）
const PANEL_COLOR_GOLD   := Color(1.0, 0.85, 0.0, 1.0)    # 金色（5層連鎖）
const PANEL_COLOR_WHITE  := Color(1.0, 1.0, 1.0, 1.0)

# ---- 節點引用 ----
var _banner_container : Control
var _banner_label     : Label
var _chain_counter    : Label   # 連鎖層數計數器
var _result_panel     : Control
var _result_label     : Label
var _flash_overlay    : ColorRect

# ---- 狀態 ----
var _chain_depth      : int = 0
var _total_kills      : int = 0
var _total_reward     : int = 0

func _ready() -> void:
	layer = 58  # 比 DragonTurtlePanel(59) 低一層
	_build_ui()
	hide()

func _build_ui() -> void:
	# 全螢幕閃光 overlay（爆炸紅色）
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 0.2, 0.1, 0.0)
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
	banner_style.bg_color = Color(0.18, 0.02, 0.0, 0.92)
	banner_style.corner_radius_top_left = 8
	banner_style.corner_radius_top_right = 8
	banner_style.corner_radius_bottom_left = 8
	banner_style.corner_radius_bottom_right = 8
	banner_style.border_width_left = 2
	banner_style.border_width_right = 2
	banner_style.border_width_top = 2
	banner_style.border_width_bottom = 2
	banner_style.border_color = PANEL_COLOR_RED
	_banner_container.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner_container)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	_banner_label.add_theme_font_size_override("font_size", 20)
	_banner_container.add_child(_banner_label)

	# 連鎖層數計數器（右上角）
	_chain_counter = Label.new()
	_chain_counter.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	_chain_counter.offset_top = 8
	_chain_counter.offset_right = -8
	_chain_counter.offset_left = -200
	_chain_counter.offset_bottom = 40
	_chain_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	_chain_counter.add_theme_color_override("font_color", PANEL_COLOR_ORANGE)
	_chain_counter.add_theme_font_size_override("font_size", 16)
	add_child(_chain_counter)

	# 結果彈窗（右側滑入）
	_result_panel = PanelContainer.new()
	_result_panel.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_panel.offset_right = 0
	_result_panel.offset_left = -300
	_result_panel.offset_top = -90
	_result_panel.offset_bottom = 90
	var result_style := StyleBoxFlat.new()
	result_style.bg_color = Color(0.15, 0.02, 0.0, 0.95)
	result_style.corner_radius_top_left = 10
	result_style.corner_radius_top_right = 10
	result_style.corner_radius_bottom_left = 10
	result_style.corner_radius_bottom_right = 10
	result_style.border_width_left = 2
	result_style.border_width_right = 2
	result_style.border_width_top = 2
	result_style.border_width_bottom = 2
	result_style.border_color = PANEL_COLOR_GOLD
	_result_panel.add_theme_stylebox_override("panel", result_style)
	_result_panel.visible = false
	add_child(_result_panel)

	_result_label = Label.new()
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_label.add_theme_color_override("font_color", PANEL_COLOR_WHITE)
	_result_label.add_theme_font_size_override("font_size", 15)
	_result_panel.add_child(_result_label)

# ---- 公開 API ----

## handle_chain_bomb 處理連鎖爆炸魚訊息
func handle_chain_bomb(payload: Dictionary) -> void:
	var phase : String = payload.get("phase", "")
	match phase:
		"chain_start":
			_on_chain_start(payload)
		"chain_explode":
			_on_chain_explode(payload)
		"chain_result":
			_on_chain_result(payload)

# ---- 私有處理函數 ----

func _on_chain_start(payload: Dictionary) -> void:
	_chain_depth = 0
	_total_kills = 0
	_total_reward = 0

	show()
	_result_panel.visible = false
	_banner_container.visible = true
	_chain_counter.visible = true
	_chain_counter.text = "💥 連鎖：第 0 層"

	var killer_name : String = payload.get("killer_name", "")
	_banner_label.text = "💥 %s 觸發連鎖爆炸！" % killer_name

	# 紅色閃光（爆炸感）
	_flash_red(0.6)

	# 橫幅從頂部滑入
	_banner_container.offset_top = -60
	_banner_container.offset_bottom = -12
	var tween := create_tween()
	tween.tween_property(_banner_container, "offset_top", 8.0, 0.25).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(_banner_container, "offset_bottom", 56.0, 0.25).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

func _on_chain_explode(payload: Dictionary) -> void:
	_chain_depth = payload.get("chain_depth", _chain_depth)
	var kill_count : int = payload.get("kill_count", 0)
	var reward     : int = payload.get("reward", 0)
	var x          : float = payload.get("trigger_x", 640.0)
	var y          : float = payload.get("trigger_y", 360.0)

	# 更新計數器
	_chain_counter.text = "💥 連鎖：第 %d 層" % _chain_depth

	# 依層數改變顏色
	var chain_color := PANEL_COLOR_RED
	if _chain_depth >= 3:
		chain_color = PANEL_COLOR_ORANGE
		_chain_counter.add_theme_color_override("font_color", PANEL_COLOR_ORANGE)
	if _chain_depth >= 5:
		chain_color = PANEL_COLOR_GOLD
		_chain_counter.add_theme_color_override("font_color", PANEL_COLOR_GOLD)

	# 爆炸圓圈動畫（在爆炸位置）
	_spawn_explosion_circle(x, y, chain_color)

	# 閃光（層數越高越強）
	var intensity := 0.2 + _chain_depth * 0.08
	_flash_red(minf(intensity, 0.6))

	# 獎勵浮動文字
	if reward > 0:
		_spawn_reward_float(x, y, reward, chain_color)

func _on_chain_result(payload: Dictionary) -> void:
	_chain_depth  = payload.get("chain_depth", _chain_depth)
	_total_kills  = payload.get("total_kills", 0)
	_total_reward = payload.get("total_reward", 0)
	var killer_name : String = payload.get("killer_name", "")

	# 隱藏橫幅和計數器
	_banner_container.visible = false
	_chain_counter.visible = false

	# 決定結果圖示
	var icon := "💥"
	if _chain_depth >= 3:
		icon = "🔥💥"
	if _chain_depth >= 5:
		icon = "🌟💥💥"

	# 顯示結果彈窗
	_result_panel.visible = true
	_result_label.text = "%s 連鎖爆炸結算\n\n觸發者：%s\n連鎖層數：%d 層\n擊破目標：%d 個\n總獎勵：%d 金幣" % [
		icon, killer_name, _chain_depth, _total_kills, _total_reward
	]

	# 依層數決定結果顏色
	var result_color := PANEL_COLOR_RED
	if _chain_depth >= 3:
		result_color = PANEL_COLOR_ORANGE
	if _chain_depth >= 5:
		result_color = PANEL_COLOR_GOLD
	_result_label.add_theme_color_override("font_color", result_color)

	# 右側滑入
	_result_panel.offset_right = 300
	_result_panel.offset_left = 20
	var tween := create_tween()
	tween.tween_property(_result_panel, "offset_right", 0.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(_result_panel, "offset_left", -300.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

	# 最終閃光
	if _chain_depth >= 5:
		_flash_gold(0.7)
	elif _chain_depth >= 3:
		_flash_orange(0.5)
	else:
		_flash_red(0.4)

	# 4 秒後淡出
	await get_tree().create_timer(4.0).timeout
	var fade_tween := create_tween()
	fade_tween.tween_property(_result_panel, "modulate:a", 0.0, 0.5)
	await fade_tween.finished
	hide()
	_result_panel.modulate.a = 1.0

# ---- 視覺效果 ----

func _flash_red(intensity: float) -> void:
	_flash_overlay.color = Color(1.0, 0.2, 0.1, intensity * 0.4)
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, 0.2)

func _flash_orange(intensity: float) -> void:
	_flash_overlay.color = Color(1.0, 0.55, 0.0, intensity * 0.4)
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, 0.25)

func _flash_gold(intensity: float) -> void:
	_flash_overlay.color = Color(1.0, 0.85, 0.0, intensity * 0.4)
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, 0.3)

## _spawn_explosion_circle 在指定位置生成爆炸圓圈動畫
func _spawn_explosion_circle(x: float, y: float, color: Color) -> void:
	# 外圈擴散
	var circle := ColorRect.new()
	circle.color = Color(color.r, color.g, color.b, 0.6)
	circle.size = Vector2(20, 20)
	circle.position = Vector2(x - 10, y - 10)
	add_child(circle)

	var tween := create_tween()
	tween.tween_property(circle, "size", Vector2(400, 400), 0.4).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(circle, "position", Vector2(x - 200, y - 200), 0.4)
	tween.parallel().tween_property(circle, "color:a", 0.0, 0.4)
	await tween.finished
	if is_instance_valid(circle):
		circle.queue_free()

	# 4 方向短線（爆炸感）
	for i in range(4):
		var line := ColorRect.new()
		line.color = color
		line.size = Vector2(4, 30)
		var angle := i * 90.0
		var rad := deg_to_rad(angle)
		line.position = Vector2(x + cos(rad) * 20, y + sin(rad) * 20)
		add_child(line)

		var line_tween := create_tween()
		line_tween.tween_property(line, "position", Vector2(x + cos(rad) * 80, y + sin(rad) * 80), 0.3)
		line_tween.parallel().tween_property(line, "modulate:a", 0.0, 0.3)
		await line_tween.finished
		if is_instance_valid(line):
			line.queue_free()

## _spawn_reward_float 獎勵浮動文字
func _spawn_reward_float(x: float, y: float, reward: int, color: Color) -> void:
	var label := Label.new()
	label.text = "+%d" % reward
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", 18)
	label.position = Vector2(x - 30, y - 20)
	add_child(label)

	var tween := create_tween()
	tween.tween_property(label, "position:y", y - 70, 1.0)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 1.0)
	await tween.finished
	if is_instance_valid(label):
		label.queue_free()
