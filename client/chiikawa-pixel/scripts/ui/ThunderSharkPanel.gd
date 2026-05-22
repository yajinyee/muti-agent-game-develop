## ThunderSharkPanel.gd — 雷霆鯊魚連鎖閃電 UI 面板（DAY-181）
## 業界依據：JILI Jackpot Fishing「Thunder Shark brings unique abilities —
## chain lightning that jumps between nearby fish, with no distance limit」
## 顯示連鎖閃電開始、每跳閃電動畫、跳數計數器、結果彈窗
extends CanvasLayer

# ---- 常數 ----
const PANEL_COLOR_BG      := Color(0.0, 0.05, 0.15, 0.90)
const PANEL_COLOR_YELLOW  := Color(1.0, 0.9, 0.0, 1.0)    # 閃電黃
const PANEL_COLOR_CYAN    := Color(0.0, 1.0, 1.0, 1.0)    # 電藍
const PANEL_COLOR_WHITE   := Color(1.0, 1.0, 1.0, 1.0)

# ---- 節點引用 ----
var _banner_container  : Control
var _banner_label      : Label
var _jump_counter      : Label
var _result_panel      : Control
var _result_label      : Label
var _flash_overlay     : ColorRect

# ---- 狀態 ----
var _jump_count        : int = 0
var _kill_count        : int = 0

func _ready() -> void:
	layer = 64
	_build_ui()
	hide()

func _build_ui() -> void:
	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 0.9, 0.0, 0.0)
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
	banner_style.bg_color = Color(0.0, 0.1, 0.3, 0.92)
	banner_style.corner_radius_top_left = 8
	banner_style.corner_radius_top_right = 8
	banner_style.corner_radius_bottom_left = 8
	banner_style.corner_radius_bottom_right = 8
	_banner_container.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner_container)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_YELLOW)
	_banner_label.add_theme_font_size_override("font_size", 20)
	_banner_container.add_child(_banner_label)

	# 跳數計數器（右上角）
	_jump_counter = Label.new()
	_jump_counter.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	_jump_counter.offset_top = 64
	_jump_counter.offset_right = -16
	_jump_counter.offset_left = -200
	_jump_counter.offset_bottom = 96
	_jump_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	_jump_counter.add_theme_color_override("font_color", PANEL_COLOR_CYAN)
	_jump_counter.add_theme_font_size_override("font_size", 18)
	_jump_counter.hide()
	add_child(_jump_counter)

	# 結果彈窗（右側滑入）
	_result_panel = PanelContainer.new()
	_result_panel.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_panel.offset_right = -16
	_result_panel.offset_left = -280
	_result_panel.offset_top = -80
	_result_panel.offset_bottom = 80
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
	result_style.border_color = PANEL_COLOR_YELLOW
	_result_panel.add_theme_stylebox_override("panel", result_style)
	_result_panel.modulate.a = 0.0
	add_child(_result_panel)

	_result_label = Label.new()
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_label.add_theme_color_override("font_color", PANEL_COLOR_WHITE)
	_result_label.add_theme_font_size_override("font_size", 18)
	_result_panel.add_child(_result_label)

## handle_thunder_shark — 處理雷霆鯊魚訊息
func handle_thunder_shark(payload: Dictionary) -> void:
	var phase : String = payload.get("phase", "")
	if phase == "chain_start":
		_on_chain_start(payload)
	elif phase.begins_with("jump_"):
		_on_jump(payload)
	elif phase == "result":
		_on_result(payload)

## _on_chain_start — 連鎖閃電開始
func _on_chain_start(payload: Dictionary) -> void:
	_jump_count = 0
	_kill_count = 0
	var killer_name : String = payload.get("killer_name", "")

	show()

	# 全螢幕黃色閃光
	_flash_screen(Color(1.0, 0.9, 0.0, 0.5), 0.3)

	# 橫幅
	_banner_label.text = "⚡ %s 觸發雷霆鯊魚！全場連鎖閃電！" % killer_name
	_banner_container.modulate.a = 0.0
	_banner_container.show()
	var tween := create_tween()
	tween.tween_property(_banner_container, "modulate:a", 1.0, 0.3)

	# 跳數計數器
	_jump_counter.text = "⚡ 0 跳 / 0 擊破"
	_jump_counter.show()

## _on_jump — 每跳閃電
func _on_jump(payload: Dictionary) -> void:
	_jump_count = payload.get("jump_num", _jump_count + 1)
	var jump_x : float = payload.get("jump_x", 0.0)
	var jump_y : float = payload.get("jump_y", 0.0)

	# 更新計數器
	_jump_counter.text = "⚡ %d 跳 / %d 擊破" % [_jump_count, _kill_count]

	# 閃電閃光（小）
	_flash_screen(Color(1.0, 1.0, 0.0, 0.2), 0.1)

	# 在目標位置顯示閃電符號
	_spawn_lightning_at(Vector2(jump_x, jump_y))

## _on_result — 連鎖結果
func _on_result(payload: Dictionary) -> void:
	var total_jumps : int = payload.get("total_jumps", 0)
	var total_kills : int = payload.get("total_kills", 0)
	var total_reward : int = payload.get("total_reward", 0)
	var killer_name : String = payload.get("killer_name", "")

	_jump_count = total_jumps
	_kill_count = total_kills
	_jump_counter.text = "⚡ %d 跳 / %d 擊破" % [total_jumps, total_kills]

	# 結果彈窗
	_result_label.text = "⚡ 雷霆連鎖結果\n跳數：%d\n擊破：%d\n獎勵：+%d" % [total_jumps, total_kills, total_reward]
	_result_panel.offset_left = 0
	_result_panel.modulate.a = 1.0
	var tween := create_tween()
	tween.tween_property(_result_panel, "offset_left", -280.0, 0.4).set_trans(Tween.TRANS_BACK)

	# 高跳數特效
	if total_jumps >= 15:
		_flash_screen(Color(1.0, 1.0, 0.0, 0.6), 0.4)
	elif total_jumps >= 10:
		_flash_screen(Color(1.0, 0.9, 0.0, 0.4), 0.3)

	# 3 秒後淡出
	await get_tree().create_timer(3.5).timeout
	var fade := create_tween()
	fade.tween_property(_result_panel, "modulate:a", 0.0, 0.5)
	fade.parallel().tween_property(_banner_container, "modulate:a", 0.0, 0.5)
	await fade.finished
	_jump_counter.hide()
	_banner_container.hide()
	hide()

## _spawn_lightning_at — 在指定位置顯示閃電符號
func _spawn_lightning_at(pos: Vector2) -> void:
	var lbl := Label.new()
	lbl.text = "⚡"
	lbl.add_theme_font_size_override("font_size", 24)
	lbl.add_theme_color_override("font_color", PANEL_COLOR_YELLOW)
	lbl.position = pos - Vector2(12, 12)
	lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(lbl)
	var tween := create_tween()
	tween.tween_property(lbl, "scale", Vector2(1.5, 1.5), 0.1)
	tween.tween_property(lbl, "scale", Vector2(1.0, 1.0), 0.05)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.3)
	await tween.finished
	if is_instance_valid(lbl):
		lbl.queue_free()

## _flash_screen — 全螢幕閃光
func _flash_screen(color: Color, duration: float) -> void:
	_flash_overlay.color = color
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)
