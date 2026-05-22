## LightningAutoChainPanel.gd — 閃電魚自動連鎖 UI 面板（DAY-183）
## 業界依據：Ocean King 3 Monster Awaken「Lightning Fish — Catching a Lightning Fish will
## trigger a Lightning Chain. Lightning Chain will continue to catch fish automatically
## until time runs out.」
## 顯示自動連鎖開始、每次自動攻擊閃電動畫、攻擊計數器、結果彈窗
extends CanvasLayer

# ---- 常數 ----
const PANEL_COLOR_BG      := Color(0.0, 0.05, 0.1, 0.90)
const PANEL_COLOR_YELLOW  := Color(1.0, 1.0, 0.0, 1.0)    # 亮黃（閃電感）
const PANEL_COLOR_CYAN    := Color(0.0, 1.0, 1.0, 1.0)    # 電藍
const PANEL_COLOR_WHITE   := Color(1.0, 1.0, 1.0, 1.0)

# ---- 節點引用 ----
var _banner_container  : Control
var _banner_label      : Label
var _attack_counter    : Label
var _timer_bar         : Control   # 進度條（顯示剩餘時間）
var _timer_fill        : ColorRect
var _result_panel      : Control
var _result_label      : Label
var _flash_overlay     : ColorRect

# ---- 狀態 ----
var _chain_active      : bool = false
var _time_remaining    : float = 0.0
var _total_duration    : float = 8.0
var _attack_count      : int = 0
var _kill_count        : int = 0

func _ready() -> void:
	layer = 62
	_build_ui()
	hide()

func _build_ui() -> void:
	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 1.0, 0.0, 0.0)
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
	banner_style.bg_color = Color(0.0, 0.1, 0.25, 0.92)
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

	# 攻擊計數器（右上角）
	_attack_counter = Label.new()
	_attack_counter.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	_attack_counter.offset_top = 64
	_attack_counter.offset_right = -16
	_attack_counter.offset_left = -220
	_attack_counter.offset_bottom = 96
	_attack_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	_attack_counter.add_theme_color_override("font_color", PANEL_COLOR_CYAN)
	_attack_counter.add_theme_font_size_override("font_size", 18)
	_attack_counter.hide()
	add_child(_attack_counter)

	# 時間進度條（底部）
	_timer_bar = Control.new()
	_timer_bar.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	_timer_bar.offset_bottom = 0
	_timer_bar.offset_top = -12
	_timer_bar.offset_left = 0
	_timer_bar.offset_right = 0
	add_child(_timer_bar)

	var bar_bg := ColorRect.new()
	bar_bg.color = Color(0.1, 0.1, 0.1, 0.8)
	bar_bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_timer_bar.add_child(bar_bg)

	_timer_fill = ColorRect.new()
	_timer_fill.color = PANEL_COLOR_YELLOW
	_timer_fill.set_anchors_preset(Control.PRESET_LEFT_WIDE)
	_timer_fill.offset_right = 0
	_timer_bar.add_child(_timer_fill)

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

func _process(delta: float) -> void:
	if _chain_active:
		_time_remaining -= delta
		if _time_remaining <= 0.0:
			_chain_active = false
			_timer_bar.hide()
		else:
			# 更新進度條
			var pct := _time_remaining / _total_duration
			_timer_fill.offset_right = _timer_bar.size.x * pct - _timer_bar.size.x
			# 顏色從黃→橙→紅
			if pct > 0.5:
				_timer_fill.color = PANEL_COLOR_YELLOW
			elif pct > 0.25:
				_timer_fill.color = Color(1.0, 0.6, 0.0, 1.0)
			else:
				_timer_fill.color = Color(1.0, 0.2, 0.0, 1.0)

## handle_lightning_auto_chain — 處理閃電魚自動連鎖訊息
func handle_lightning_auto_chain(payload: Dictionary) -> void:
	var phase : String = payload.get("phase", "")
	if phase == "chain_start":
		_on_chain_start(payload)
	elif phase.begins_with("auto_"):
		_on_auto_attack(payload)
	elif phase == "result":
		_on_result(payload)

## _on_chain_start — 自動連鎖開始
func _on_chain_start(payload: Dictionary) -> void:
	_chain_active = true
	_total_duration = float(payload.get("duration_sec", 8))
	_time_remaining = _total_duration
	_attack_count = 0
	_kill_count = 0
	var killer_name : String = payload.get("killer_name", "")

	show()

	# 全螢幕黃色閃光
	_flash_screen(Color(1.0, 1.0, 0.0, 0.5), 0.3)

	# 橫幅
	_banner_label.text = "⚡ %s 觸發閃電魚！自動連鎖 %d 秒！" % [killer_name, int(_total_duration)]
	_banner_container.modulate.a = 0.0
	_banner_container.show()
	var tween := create_tween()
	tween.tween_property(_banner_container, "modulate:a", 1.0, 0.3)

	# 攻擊計數器
	_attack_counter.text = "⚡ 0 攻擊 / 0 擊破"
	_attack_counter.show()

	# 時間進度條
	_timer_bar.show()

## _on_auto_attack — 每次自動攻擊
func _on_auto_attack(payload: Dictionary) -> void:
	_attack_count = payload.get("attack_num", _attack_count + 1)
	var target_x : float = payload.get("target_x", 0.0)
	var target_y : float = payload.get("target_y", 0.0)

	# 更新計數器
	_attack_counter.text = "⚡ %d 攻擊 / %d 擊破" % [_attack_count, _kill_count]

	# 閃電閃光（小）
	_flash_screen(Color(1.0, 1.0, 0.0, 0.15), 0.08)

	# 在目標位置顯示閃電符號
	_spawn_lightning_at(Vector2(target_x, target_y))

## _on_result — 自動連鎖結果
func _on_result(payload: Dictionary) -> void:
	var total_attacks : int = payload.get("total_attacks", 0)
	var total_kills : int = payload.get("total_kills", 0)
	var total_reward : int = payload.get("total_reward", 0)

	_chain_active = false
	_attack_count = total_attacks
	_kill_count = total_kills
	_attack_counter.text = "⚡ %d 攻擊 / %d 擊破" % [total_attacks, total_kills]
	_timer_bar.hide()

	# 結果彈窗
	_result_label.text = "⚡ 閃電自動連鎖結束\n攻擊：%d 次\n擊破：%d 個\n獎勵：+%d" % [total_attacks, total_kills, total_reward]
	_result_panel.offset_left = 0
	_result_panel.modulate.a = 1.0
	var tween := create_tween()
	tween.tween_property(_result_panel, "offset_left", -280.0, 0.4).set_trans(Tween.TRANS_BACK)

	# 高擊破數特效
	if total_kills >= 10:
		_flash_screen(Color(1.0, 1.0, 0.0, 0.6), 0.4)
	elif total_kills >= 6:
		_flash_screen(Color(1.0, 0.9, 0.0, 0.3), 0.3)

	# 3 秒後淡出
	await get_tree().create_timer(3.5).timeout
	var fade := create_tween()
	fade.tween_property(_result_panel, "modulate:a", 0.0, 0.5)
	fade.parallel().tween_property(_banner_container, "modulate:a", 0.0, 0.5)
	await fade.finished
	_attack_counter.hide()
	_banner_container.hide()
	hide()

## _spawn_lightning_at — 在指定位置顯示閃電符號
func _spawn_lightning_at(pos: Vector2) -> void:
	var lbl := Label.new()
	lbl.text = "⚡"
	lbl.add_theme_font_size_override("font_size", 20)
	lbl.add_theme_color_override("font_color", PANEL_COLOR_YELLOW)
	lbl.position = pos - Vector2(10, 10)
	lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(lbl)
	var tween := create_tween()
	tween.tween_property(lbl, "scale", Vector2(1.4, 1.4), 0.08)
	tween.tween_property(lbl, "scale", Vector2(1.0, 1.0), 0.04)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.25)
	await tween.finished
	if is_instance_valid(lbl):
		lbl.queue_free()

## _flash_screen — 全螢幕閃光
func _flash_screen(color: Color, duration: float) -> void:
	_flash_overlay.color = color
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)
