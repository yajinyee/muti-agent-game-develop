## RoulettePanel.gd - DAY-113
## 雙層倍率輪盤 UI：擊殺 BOSS/特殊目標後觸發
## 內圈(8格,1-10x) × 外圈(12格,1-100x) = 最終倍率（最高 1000x）
## 廣播給所有玩家觀看，觸發玩家有互動感
extends Node2D

const PANEL_W := 480
const PANEL_H := 480
const INNER_RADIUS := 80.0
const OUTER_RADIUS := 150.0
const INNER_COUNT := 8
const OUTER_COUNT := 12

var _font: FontFile
var _overlay: ColorRect
var _bg: ColorRect
var _title_label: Label
var _player_label: Label
var _target_label: Label

# 輪盤格子節點
var _inner_rects: Array = []
var _inner_labels: Array = []
var _outer_rects: Array = []
var _outer_labels: Array = []

# 結果顯示
var _inner_result_label: Label
var _outer_result_label: Label
var _mult_label: Label
var _reward_label: Label
var _close_btn: Button

# 狀態
var _inner_data: Array = []
var _outer_data: Array = []
var _is_spinning: bool = false
var _is_self: bool = false
var _session_id: String = ""

# 旋轉動畫狀態
var _inner_spin_tween: Tween = null
var _outer_spin_tween: Tween = null
var _inner_win_idx: int = 0
var _outer_win_idx: int = 0

func setup(font: FontFile) -> void:
	_font = font
	_build_ui()
	_connect_signals()
	hide()

func _build_ui() -> void:
	# 全螢幕半透明遮罩
	_overlay = ColorRect.new()
	_overlay.color = Color(0.0, 0.0, 0.0, 0.82)
	_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_overlay.size = Vector2(1280, 720)
	_overlay.position = Vector2(-640, -180)
	add_child(_overlay)

	# 主面板背景（深藍色）
	_bg = ColorRect.new()
	_bg.color = Color(0.04, 0.04, 0.18, 0.98)
	_bg.size = Vector2(PANEL_W, PANEL_H)
	_bg.position = Vector2(-PANEL_W / 2, -PANEL_H / 2)
	add_child(_bg)

	# 金色頂部邊框
	var border_top = ColorRect.new()
	border_top.color = Color(1.0, 0.85, 0.1, 1.0)
	border_top.size = Vector2(PANEL_W, 4)
	border_top.position = Vector2(-PANEL_W / 2, -PANEL_H / 2)
	add_child(border_top)

	# 標題
	_title_label = Label.new()
	_title_label.text = "🎡 雙層倍率輪盤"
	_title_label.position = Vector2(-PANEL_W / 2 + 12, -PANEL_H / 2 + 10)
	_title_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
	if _font:
		_title_label.add_theme_font_override("font", _font)
		_title_label.add_theme_font_size_override("font_size", 20)
	add_child(_title_label)

	# 玩家名稱
	_player_label = Label.new()
	_player_label.text = ""
	_player_label.position = Vector2(-PANEL_W / 2 + 12, -PANEL_H / 2 + 36)
	_player_label.add_theme_color_override("font_color", Color(0.7, 0.9, 1.0))
	if _font:
		_player_label.add_theme_font_override("font", _font)
		_player_label.add_theme_font_size_override("font_size", 13)
	add_child(_player_label)

	# 目標名稱
	_target_label = Label.new()
	_target_label.text = ""
	_target_label.position = Vector2(-PANEL_W / 2 + 12, -PANEL_H / 2 + 54)
	_target_label.add_theme_color_override("font_color", Color(1.0, 0.6, 0.2))
	if _font:
		_target_label.add_theme_font_override("font", _font)
		_target_label.add_theme_font_size_override("font_size", 12)
	add_child(_target_label)

	# 外圈格子（12格）
	for i in range(OUTER_COUNT):
		var angle := (i * TAU / OUTER_COUNT) - PI / 2
		var cx := OUTER_RADIUS * cos(angle)
		var cy := OUTER_RADIUS * sin(angle)

		var rect = ColorRect.new()
		rect.size = Vector2(52, 26)
		rect.position = Vector2(cx - 26, cy - 13)
		rect.color = Color(0.1, 0.1, 0.3, 0.9)
		add_child(rect)
		_outer_rects.append(rect)

		var lbl = Label.new()
		lbl.text = "?"
		lbl.position = Vector2(cx - 22, cy - 9)
		lbl.add_theme_color_override("font_color", Color.WHITE)
		if _font:
			lbl.add_theme_font_override("font", _font)
			lbl.add_theme_font_size_override("font_size", 12)
		add_child(lbl)
		_outer_labels.append(lbl)

	# 外圈標籤
	var outer_ring_label = Label.new()
	outer_ring_label.text = "外圈"
	outer_ring_label.position = Vector2(-16, -OUTER_RADIUS - 20)
	outer_ring_label.add_theme_color_override("font_color", Color(0.6, 0.8, 1.0))
	if _font:
		outer_ring_label.add_theme_font_override("font", _font)
		outer_ring_label.add_theme_font_size_override("font_size", 11)
	add_child(outer_ring_label)

	# 內圈格子（8格）
	for i in range(INNER_COUNT):
		var angle := (i * TAU / INNER_COUNT) - PI / 2
		var cx := INNER_RADIUS * cos(angle)
		var cy := INNER_RADIUS * sin(angle)

		var rect = ColorRect.new()
		rect.size = Vector2(44, 22)
		rect.position = Vector2(cx - 22, cy - 11)
		rect.color = Color(0.15, 0.15, 0.35, 0.9)
		add_child(rect)
		_inner_rects.append(rect)

		var lbl = Label.new()
		lbl.text = "?"
		lbl.position = Vector2(cx - 18, cy - 8)
		lbl.add_theme_color_override("font_color", Color.WHITE)
		if _font:
			lbl.add_theme_font_override("font", _font)
			lbl.add_theme_font_size_override("font_size", 12)
		add_child(lbl)
		_inner_labels.append(lbl)

	# 內圈標籤
	var inner_ring_label = Label.new()
	inner_ring_label.text = "內圈"
	inner_ring_label.position = Vector2(-16, -INNER_RADIUS - 16)
	inner_ring_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.5))
	if _font:
		inner_ring_label.add_theme_font_override("font", _font)
		inner_ring_label.add_theme_font_size_override("font_size", 11)
	add_child(inner_ring_label)

	# 中心圓（裝飾）
	var center_bg = ColorRect.new()
	center_bg.size = Vector2(50, 50)
	center_bg.position = Vector2(-25, -25)
	center_bg.color = Color(0.08, 0.08, 0.25, 1.0)
	add_child(center_bg)

	var center_label = Label.new()
	center_label.text = "×"
	center_label.position = Vector2(-10, -14)
	center_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
	if _font:
		center_label.add_theme_font_override("font", _font)
		center_label.add_theme_font_size_override("font_size", 24)
	add_child(center_label)

	# 結果區域（底部）
	var result_bg = ColorRect.new()
	result_bg.color = Color(0.06, 0.06, 0.2, 0.95)
	result_bg.size = Vector2(PANEL_W - 20, 100)
	result_bg.position = Vector2(-PANEL_W / 2 + 10, PANEL_H / 2 - 115)
	add_child(result_bg)

	_inner_result_label = Label.new()
	_inner_result_label.text = "內圈：旋轉中..."
	_inner_result_label.position = Vector2(-PANEL_W / 2 + 20, PANEL_H / 2 - 110)
	_inner_result_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.5))
	if _font:
		_inner_result_label.add_theme_font_override("font", _font)
		_inner_result_label.add_theme_font_size_override("font_size", 14)
	add_child(_inner_result_label)

	_outer_result_label = Label.new()
	_outer_result_label.text = "外圈：旋轉中..."
	_outer_result_label.position = Vector2(-PANEL_W / 2 + 20, PANEL_H / 2 - 90)
	_outer_result_label.add_theme_color_override("font_color", Color(0.6, 0.8, 1.0))
	if _font:
		_outer_result_label.add_theme_font_override("font", _font)
		_outer_result_label.add_theme_font_size_override("font_size", 14)
	add_child(_outer_result_label)

	_mult_label = Label.new()
	_mult_label.text = ""
	_mult_label.position = Vector2(-PANEL_W / 2 + 20, PANEL_H / 2 - 68)
	_mult_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
	if _font:
		_mult_label.add_theme_font_override("font", _font)
		_mult_label.add_theme_font_size_override("font_size", 22)
	add_child(_mult_label)

	_reward_label = Label.new()
	_reward_label.text = ""
	_reward_label.position = Vector2(-PANEL_W / 2 + 20, PANEL_H / 2 - 44)
	_reward_label.add_theme_color_override("font_color", Color(0.5, 1.0, 0.5))
	if _font:
		_reward_label.add_theme_font_override("font", _font)
		_reward_label.add_theme_font_size_override("font_size", 16)
	add_child(_reward_label)

	# 關閉按鈕
	_close_btn = Button.new()
	_close_btn.text = "太棒了！繼續"
	_close_btn.size = Vector2(160, 36)
	_close_btn.position = Vector2(-80, PANEL_H / 2 - 44)
	_close_btn.pressed.connect(_on_close)
	if _font:
		_close_btn.add_theme_font_override("font", _font)
		_close_btn.add_theme_font_size_override("font_size", 14)
	_close_btn.hide()
	add_child(_close_btn)

func _connect_signals() -> void:
	if GameManager.has_signal("roulette_started"):
		GameManager.roulette_started.connect(_on_roulette_started)
	if GameManager.has_signal("roulette_result"):
		GameManager.roulette_result.connect(_on_roulette_result)

func _on_roulette_started(data: Dictionary) -> void:
	if _is_spinning:
		return

	_session_id = data.get("session_id", "")
	_is_self = data.get("is_self", false)
	var player_name: String = data.get("player_name", "")
	var target_name: String = data.get("target_name", "")
	_inner_data = data.get("inner_segments", [])
	_outer_data = data.get("outer_segments", [])

	# 更新格子顯示
	_update_segments()

	# 更新標籤
	if _is_self:
		_player_label.text = "🎯 你觸發了雙層輪盤！"
		_player_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.3))
	else:
		_player_label.text = "👀 %s 觸發了雙層輪盤" % player_name
		_player_label.add_theme_color_override("font_color", Color(0.7, 0.9, 1.0))

	_target_label.text = "擊殺：%s" % target_name
	_inner_result_label.text = "內圈：旋轉中..."
	_outer_result_label.text = "外圈：旋轉中..."
	_mult_label.text = ""
	_reward_label.text = ""
	_close_btn.hide()
	_is_spinning = true

	show()

func _on_roulette_result(data: Dictionary) -> void:
	if data.get("session_id", "") != _session_id:
		return

	_inner_win_idx = data.get("inner", {}).get("segment_index", 0)
	_outer_win_idx = data.get("outer", {}).get("segment_index", 0)
	var inner_mult: float = data.get("inner", {}).get("multiplier", 1.0)
	var outer_mult: float = data.get("outer", {}).get("multiplier", 1.0)
	var final_mult: float = data.get("final_mult", 1.0)
	var base_reward: int = data.get("base_reward", 0)
	var final_reward: int = data.get("final_reward", 0)
	var is_jackpot: bool = data.get("is_jackpot", false)
	var is_mega_win: bool = data.get("is_mega_win", false)

	# 播放旋轉動畫，3秒後顯示結果
	_play_dual_spin_animation(inner_mult, outer_mult, final_mult, base_reward, final_reward, is_jackpot, is_mega_win)

func _update_segments() -> void:
	for i in range(min(_inner_data.size(), INNER_COUNT)):
		var seg = _inner_data[i]
		if i < _inner_labels.size():
			_inner_labels[i].text = seg.get("label", "?")
		if i < _inner_rects.size():
			var color_hex: String = seg.get("color", "#4CAF50")
			_inner_rects[i].color = Color(color_hex) * 0.6

	for i in range(min(_outer_data.size(), OUTER_COUNT)):
		var seg = _outer_data[i]
		if i < _outer_labels.size():
			_outer_labels[i].text = seg.get("label", "?")
		if i < _outer_rects.size():
			var color_hex: String = seg.get("color", "#4CAF50")
			_outer_rects[i].color = Color(color_hex) * 0.6

func _play_dual_spin_animation(inner_mult: float, outer_mult: float, final_mult: float,
		base_reward: int, final_reward: int, is_jackpot: bool, is_mega_win: bool) -> void:
	# 內圈旋轉（快速）
	var inner_steps := 24 + _inner_win_idx
	var outer_steps := 32 + _outer_win_idx  # 外圈轉更多圈

	var tween = create_tween()

	# 同時旋轉內外圈
	for step in range(max(inner_steps, outer_steps)):
		var inner_idx := step % INNER_COUNT
		var outer_idx := step % OUTER_COUNT
		var duration := 0.04
		if step > max(inner_steps, outer_steps) - 10:
			duration = 0.04 + (max(inner_steps, outer_steps) - step) * 0.05

		tween.tween_callback(func():
			_highlight_inner(inner_idx if step < inner_steps else _inner_win_idx)
			_highlight_outer(outer_idx if step < outer_steps else _outer_win_idx)
		)
		tween.tween_interval(duration)

	# 最終停止
	tween.tween_callback(func():
		_highlight_inner(_inner_win_idx)
		_highlight_outer(_outer_win_idx)
		_show_result(inner_mult, outer_mult, final_mult, base_reward, final_reward, is_jackpot, is_mega_win)
	)

func _highlight_inner(idx: int) -> void:
	for i in range(_inner_rects.size()):
		if i < _inner_data.size():
			var color_hex: String = _inner_data[i].get("color", "#4CAF50")
			if i == idx:
				_inner_rects[i].color = Color.WHITE
				if i < _inner_labels.size():
					_inner_labels[i].add_theme_color_override("font_color", Color(0.1, 0.1, 0.1))
			else:
				_inner_rects[i].color = Color(color_hex) * 0.4
				if i < _inner_labels.size():
					_inner_labels[i].add_theme_color_override("font_color", Color.WHITE)

func _highlight_outer(idx: int) -> void:
	for i in range(_outer_rects.size()):
		if i < _outer_data.size():
			var color_hex: String = _outer_data[i].get("color", "#4CAF50")
			if i == idx:
				_outer_rects[i].color = Color.WHITE
				if i < _outer_labels.size():
					_outer_labels[i].add_theme_color_override("font_color", Color(0.1, 0.1, 0.1))
			else:
				_outer_rects[i].color = Color(color_hex) * 0.4
				if i < _outer_labels.size():
					_outer_labels[i].add_theme_color_override("font_color", Color.WHITE)

func _show_result(inner_mult: float, outer_mult: float, final_mult: float,
		base_reward: int, final_reward: int, is_jackpot: bool, is_mega_win: bool) -> void:
	_is_spinning = false

	# 保持中獎格子高亮
	if _inner_win_idx < _inner_rects.size() and _inner_data.size() > _inner_win_idx:
		var color_hex: String = _inner_data[_inner_win_idx].get("color", "#FFD700")
		_inner_rects[_inner_win_idx].color = Color(color_hex)
	if _outer_win_idx < _outer_rects.size() and _outer_data.size() > _outer_win_idx:
		var color_hex: String = _outer_data[_outer_win_idx].get("color", "#FFD700")
		_outer_rects[_outer_win_idx].color = Color(color_hex)

	_inner_result_label.text = "內圈：%.0fx" % inner_mult
	_outer_result_label.text = "外圈：%.0fx" % outer_mult

	# 最終倍率顯示（依等級變色）
	if is_jackpot:
		_mult_label.text = "🌟 最終倍率：%.0fx（傳說！）" % final_mult
		_mult_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
		_play_jackpot_effect()
	elif is_mega_win:
		_mult_label.text = "💥 最終倍率：%.0fx（超大獎！）" % final_mult
		_mult_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.1))
		_play_megawin_effect()
	else:
		_mult_label.text = "✨ 最終倍率：%.0fx" % final_mult
		_mult_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))

	if _is_self:
		_reward_label.text = "🪙 獎勵：+%d 金幣（基礎 %d）" % [final_reward, base_reward]
	else:
		_reward_label.text = "🪙 獎勵：%d 金幣" % final_reward

	# 縮放彈入動畫
	scale = Vector2(0.8, 0.8)
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(1.1, 1.1), 0.15)
	tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.1)

	_close_btn.show()

func _play_jackpot_effect() -> void:
	# 全畫面金色閃光（傳說等級）
	var flash = ColorRect.new()
	flash.color = Color(1.0, 0.85, 0.1, 0.0)
	flash.size = Vector2(1280, 720)
	flash.position = Vector2(-640, -180)
	flash.z_index = 200
	add_child(flash)

	var tween = create_tween()
	tween.tween_property(flash, "color:a", 0.6, 0.15)
	tween.tween_property(flash, "color:a", 0.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)

func _play_megawin_effect() -> void:
	# 橙色閃光（超大獎）
	var flash = ColorRect.new()
	flash.color = Color(1.0, 0.5, 0.1, 0.0)
	flash.size = Vector2(1280, 720)
	flash.position = Vector2(-640, -180)
	flash.z_index = 200
	add_child(flash)

	var tween = create_tween()
	tween.tween_property(flash, "color:a", 0.4, 0.12)
	tween.tween_property(flash, "color:a", 0.0, 0.25)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)

func _on_close() -> void:
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(0.0, 0.0), 0.15)
	tween.tween_callback(func():
		scale = Vector2(1.0, 1.0)
		hide()
	)
