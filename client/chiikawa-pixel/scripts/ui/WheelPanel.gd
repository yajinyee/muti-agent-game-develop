## WheelPanel.gd - DAY-084
## 幸運轉盤 UI：擊殺特殊目標後觸發，顯示旋轉動畫和獎勵結果
extends Node2D

const PANEL_W := 320
const PANEL_H := 360
const SLOT_COUNT := 8

var _font: FontFile
var _bg: ColorRect
var _overlay: ColorRect
var _title_label: Label
var _target_label: Label
var _wheel_container: Control
var _slot_rects: Array = []
var _slot_labels: Array = []
var _pointer: ColorRect
var _result_label: Label
var _reward_label: Label
var _close_btn: Button

var _slots_data: Array = []
var _win_index: int = 0
var _base_reward: int = 0
var _final_reward: int = 0
var _is_spinning: bool = false
var _spin_tween: Tween = null

# 轉盤旋轉狀態
var _current_angle: float = 0.0
var _target_angle: float = 0.0

func setup(font: FontFile) -> void:
	_font = font
	_build_ui()
	_connect_signals()
	hide()

func _build_ui() -> void:
	# 全螢幕半透明遮罩
	_overlay = ColorRect.new()
	_overlay.color = Color(0.0, 0.0, 0.0, 0.75)
	_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_overlay.size = Vector2(1280, 720)
	_overlay.position = Vector2(-640, -180)
	add_child(_overlay)

	# 主面板背景
	_bg = ColorRect.new()
	_bg.color = Color(0.05, 0.05, 0.15, 0.97)
	_bg.size = Vector2(PANEL_W, PANEL_H)
	_bg.position = Vector2(-PANEL_W / 2, -PANEL_H / 2)
	add_child(_bg)

	# 金色邊框
	var border = ColorRect.new()
	border.color = Color(1.0, 0.85, 0.1, 1.0)
	border.size = Vector2(PANEL_W, 3)
	border.position = Vector2(-PANEL_W / 2, -PANEL_H / 2)
	add_child(border)
	var border_b = ColorRect.new()
	border_b.color = Color(1.0, 0.85, 0.1, 1.0)
	border_b.size = Vector2(PANEL_W, 3)
	border_b.position = Vector2(-PANEL_W / 2, PANEL_H / 2 - 3)
	add_child(border_b)

	# 標題
	_title_label = Label.new()
	_title_label.text = "🎰 幸運轉盤"
	_title_label.position = Vector2(-PANEL_W / 2 + 10, -PANEL_H / 2 + 8)
	_title_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
	if _font:
		_title_label.add_theme_font_override("font", _font)
		_title_label.add_theme_font_size_override("font_size", 18)
	add_child(_title_label)

	# 目標物名稱
	_target_label = Label.new()
	_target_label.text = ""
	_target_label.position = Vector2(-PANEL_W / 2 + 10, -PANEL_H / 2 + 32)
	_target_label.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
	if _font:
		_target_label.add_theme_font_override("font", _font)
		_target_label.add_theme_font_size_override("font_size", 12)
	add_child(_target_label)

	# 轉盤容器（8個格子排成圓形）
	_wheel_container = Control.new()
	_wheel_container.size = Vector2(200, 200)
	_wheel_container.position = Vector2(-100, -80)
	add_child(_wheel_container)

	# 建立 8 個格子
	var radius := 80.0
	for i in range(SLOT_COUNT):
		var angle := (i * TAU / SLOT_COUNT) - PI / 2
		var cx := radius * cos(angle)
		var cy := radius * sin(angle)

		var slot_bg = ColorRect.new()
		slot_bg.size = Vector2(56, 28)
		slot_bg.position = Vector2(100 + cx - 28, 100 + cy - 14)
		slot_bg.color = Color(0.15, 0.15, 0.25, 0.9)
		_wheel_container.add_child(slot_bg)
		_slot_rects.append(slot_bg)

		var slot_lbl = Label.new()
		slot_lbl.text = "?"
		slot_lbl.position = Vector2(100 + cx - 24, 100 + cy - 10)
		slot_lbl.add_theme_color_override("font_color", Color.WHITE)
		if _font:
			slot_lbl.add_theme_font_override("font", _font)
			slot_lbl.add_theme_font_size_override("font_size", 14)
		_wheel_container.add_child(slot_lbl)
		_slot_labels.append(slot_lbl)

	# 中心指針（三角形用 ColorRect 模擬）
	_pointer = ColorRect.new()
	_pointer.color = Color(1.0, 0.3, 0.3, 1.0)
	_pointer.size = Vector2(12, 20)
	_pointer.position = Vector2(94, 72)
	_wheel_container.add_child(_pointer)

	# 結果標籤
	_result_label = Label.new()
	_result_label.text = ""
	_result_label.position = Vector2(-PANEL_W / 2 + 10, 100)
	_result_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
	if _font:
		_result_label.add_theme_font_override("font", _font)
		_result_label.add_theme_font_size_override("font_size", 22)
	add_child(_result_label)

	# 獎勵標籤
	_reward_label = Label.new()
	_reward_label.text = ""
	_reward_label.position = Vector2(-PANEL_W / 2 + 10, 130)
	_reward_label.add_theme_color_override("font_color", Color(0.5, 1.0, 0.5))
	if _font:
		_reward_label.add_theme_font_override("font", _font)
		_reward_label.add_theme_font_size_override("font_size", 16)
	add_child(_reward_label)

	# 關閉按鈕
	_close_btn = Button.new()
	_close_btn.text = "繼續"
	_close_btn.size = Vector2(120, 36)
	_close_btn.position = Vector2(-60, 150)
	_close_btn.pressed.connect(_on_close)
	if _font:
		_close_btn.add_theme_font_override("font", _font)
		_close_btn.add_theme_font_size_override("font_size", 14)
	_close_btn.hide()
	add_child(_close_btn)

func _connect_signals() -> void:
	if GameManager.has_signal("wheel_triggered"):
		GameManager.wheel_triggered.connect(_on_wheel_triggered)

func _on_wheel_triggered(data: Dictionary) -> void:
	if _is_spinning:
		return

	_slots_data = data.get("slots", [])
	_win_index = data.get("win_index", 0)
	_base_reward = data.get("base_reward", 0)
	_final_reward = data.get("final_reward", 0)
	var target_name: String = data.get("target_name", "")
	var multiplier: float = data.get("multiplier", 1.0)

	# 更新格子顯示
	for i in range(min(_slots_data.size(), SLOT_COUNT)):
		var slot = _slots_data[i]
		if i < _slot_labels.size():
			_slot_labels[i].text = slot.get("label", "?")
		if i < _slot_rects.size():
			var color_hex: String = slot.get("color", "#4CAF50")
			_slot_rects[i].color = Color(color_hex)

	_target_label.text = "擊殺：%s" % target_name
	_result_label.text = ""
	_reward_label.text = ""
	_close_btn.hide()
	_is_spinning = true

	# 顯示面板
	show()
	_play_spin_animation(multiplier)

func _play_spin_animation(multiplier: float) -> void:
	# 高亮格子依序閃爍，模擬轉盤旋轉
	var spin_steps := 24 + _win_index  # 至少轉 3 圈
	var step_duration := 0.05
	var current_step := 0

	_spin_tween = create_tween()

	for step in range(spin_steps):
		var highlight_idx := step % SLOT_COUNT
		var duration := step_duration
		# 後半段減速
		if step > spin_steps - 8:
			duration = step_duration + (spin_steps - step) * 0.04

		_spin_tween.tween_callback(func():
			_highlight_slot(highlight_idx)
		)
		_spin_tween.tween_interval(duration)

	# 最後停在中獎格子
	_spin_tween.tween_callback(func():
		_highlight_slot(_win_index)
		_show_result(multiplier)
	)

func _highlight_slot(idx: int) -> void:
	for i in range(_slot_rects.size()):
		if i < _slot_rects.size():
			var base_color: Color = Color(0.15, 0.15, 0.25, 0.9)
			if _slots_data.size() > i:
				var color_hex: String = _slots_data[i].get("color", "#4CAF50")
				base_color = Color(color_hex)
			if i == idx:
				_slot_rects[i].color = Color.WHITE
				if i < _slot_labels.size():
					_slot_labels[i].add_theme_color_override("font_color", Color(0.1, 0.1, 0.1))
			else:
				_slot_rects[i].color = base_color * 0.5
				if i < _slot_labels.size():
					_slot_labels[i].add_theme_color_override("font_color", Color.WHITE)

func _show_result(multiplier: float) -> void:
	_is_spinning = false

	# 中獎格子保持高亮
	if _win_index < _slot_rects.size() and _slots_data.size() > _win_index:
		var color_hex: String = _slots_data[_win_index].get("color", "#FFD700")
		_slot_rects[_win_index].color = Color(color_hex)

	_result_label.text = "🎉 %.0fx 大獎！" % multiplier
	_reward_label.text = "+%d 金幣（基礎 %d）" % [_final_reward, _base_reward]

	# 放大動畫
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(1.15, 1.15), 0.12)
	tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.1)

	# 顯示關閉按鈕
	_close_btn.show()

func _on_close() -> void:
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(0.0, 0.0), 0.15)
	tween.tween_callback(func():
		scale = Vector2(1.0, 1.0)
		hide()
	)
