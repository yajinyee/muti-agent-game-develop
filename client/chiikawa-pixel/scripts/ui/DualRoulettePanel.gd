## DualRoulettePanel.gd
## 雙環輪盤系統面板（DAY-139）
## 業界依據：Royal Fishing JILI 2026 ChainLong King Dual-Ring Roulette
## 擊破高倍率目標後觸發，內外圈相乘最高 150x，製造「技巧感」
## 玩家點擊停止按鈕，內外圈停止，顯示最終倍率和獎勵

extends Control

# ---- 狀態 ----
var _is_active: bool = false
var _spin_duration: float = 3.0
var _elapsed: float = 0.0
var _inner_ring: Array = [2.0, 3.0, 5.0, 8.0, 10.0]
var _outer_ring: Array = [2.0, 3.0, 5.0, 7.0, 10.0, 15.0]
var _inner_angle: float = 0.0
var _outer_angle: float = 0.0
var _inner_speed: float = 360.0  # 度/秒（內環較快）
var _outer_speed: float = 240.0  # 度/秒（外環較慢）
var _base_reward: int = 0
var _target_mult: float = 0.0

# ---- 節點 ----
var _overlay: ColorRect = null
var _panel: Control = null
var _title_label: Label = null
var _inner_ring_node: Control = null
var _outer_ring_node: Control = null
var _inner_pointer: Control = null
var _outer_pointer: Control = null
var _stop_btn: Button = null
var _result_panel: Control = null
var _result_label: Label = null
var _bonus_label: Label = null
var _flash_overlay: ColorRect = null
var _countdown_label: Label = null

# ---- 顏色 ----
const COLOR_GOLD   = Color(1.0, 0.84, 0.0, 1.0)
const COLOR_INNER  = Color(0.2, 0.6, 1.0, 1.0)   # 藍色內環
const COLOR_OUTER  = Color(1.0, 0.5, 0.1, 1.0)   # 橙色外環
const COLOR_BG     = Color(0.05, 0.05, 0.15, 0.97)

func _ready() -> void:
	_build_ui()
	set_process(false)
	visible = false
	mouse_filter = Control.MOUSE_FILTER_IGNORE

func _build_ui() -> void:
	# 全螢幕半透明遮罩
	_overlay = ColorRect.new()
	_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_overlay.color = Color(0.0, 0.0, 0.0, 0.75)
	_overlay.mouse_filter = Control.MOUSE_FILTER_STOP
	add_child(_overlay)

	# 全螢幕閃光層
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.84, 0.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 主面板（置中）
	_panel = Control.new()
	_panel.set_anchors_preset(Control.PRESET_CENTER)
	_panel.custom_minimum_size = Vector2(480, 420)
	_panel.position = Vector2(-240, -210)
	add_child(_panel)

	var panel_bg = ColorRect.new()
	panel_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel_bg.color = COLOR_BG
	_panel.add_child(panel_bg)

	# 金色邊框
	var border = ColorRect.new()
	border.set_anchors_preset(Control.PRESET_FULL_RECT)
	border.color = Color(1.0, 0.84, 0.0, 0.6)
	border.custom_minimum_size = Vector2(480, 420)
	_panel.add_child(border)

	var inner_bg = ColorRect.new()
	inner_bg.color = COLOR_BG
	inner_bg.position = Vector2(3, 3)
	inner_bg.size = Vector2(474, 414)
	_panel.add_child(inner_bg)

	# 標題
	_title_label = Label.new()
	_title_label.text = "🎡 雙環輪盤"
	_title_label.position = Vector2(0, 12)
	_title_label.size = Vector2(480, 36)
	_title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_title_label.add_theme_color_override("font_color", COLOR_GOLD)
	_title_label.add_theme_font_size_override("font_size", 22)
	_panel.add_child(_title_label)

	# 倒數計時
	_countdown_label = Label.new()
	_countdown_label.text = "3.0s"
	_countdown_label.position = Vector2(0, 48)
	_countdown_label.size = Vector2(480, 24)
	_countdown_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_countdown_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_countdown_label.add_theme_font_size_override("font_size", 14)
	_panel.add_child(_countdown_label)

	# 外環（大圓，橙色）
	_outer_ring_node = Control.new()
	_outer_ring_node.position = Vector2(240 - 140, 80)
	_outer_ring_node.custom_minimum_size = Vector2(280, 280)
	_panel.add_child(_outer_ring_node)

	# 內環（小圓，藍色）
	_inner_ring_node = Control.new()
	_inner_ring_node.position = Vector2(240 - 90, 130)
	_inner_ring_node.custom_minimum_size = Vector2(180, 180)
	_panel.add_child(_inner_ring_node)

	# 停止按鈕
	_stop_btn = Button.new()
	_stop_btn.text = "⏹ 停止！"
	_stop_btn.position = Vector2(140, 370)
	_stop_btn.size = Vector2(200, 44)
	_stop_btn.add_theme_color_override("font_color", Color(0.0, 0.0, 0.0, 1.0))
	_stop_btn.add_theme_font_size_override("font_size", 18)
	_stop_btn.pressed.connect(_on_stop_pressed)
	_panel.add_child(_stop_btn)

	# 結果面板（停止後顯示）
	_result_panel = Control.new()
	_result_panel.position = Vector2(0, 200)
	_result_panel.size = Vector2(480, 120)
	_result_panel.visible = false
	_panel.add_child(_result_panel)

	_result_label = Label.new()
	_result_label.text = "🎯 內環 5x × 外環 10x = 50x"
	_result_label.position = Vector2(0, 10)
	_result_label.size = Vector2(480, 40)
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_label.add_theme_color_override("font_color", COLOR_GOLD)
	_result_label.add_theme_font_size_override("font_size", 18)
	_result_panel.add_child(_result_label)

	_bonus_label = Label.new()
	_bonus_label.text = "+5000 金幣！"
	_bonus_label.position = Vector2(0, 55)
	_bonus_label.size = Vector2(480, 48)
	_bonus_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_bonus_label.add_theme_color_override("font_color", Color(0.2, 1.0, 0.4, 1.0))
	_bonus_label.add_theme_font_size_override("font_size", 28)
	_result_panel.add_child(_bonus_label)

func _process(delta: float) -> void:
	if not _is_active:
		return

	_elapsed += delta
	_inner_angle += _inner_speed * delta
	_outer_angle += _outer_speed * delta

	# 更新倒數計時
	var remaining = max(0.0, _spin_duration - _elapsed)
	_countdown_label.text = "%.1fs" % remaining
	if remaining < 1.0:
		_countdown_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3, 1.0))

	# 更新環的旋轉視覺（用 modulate 閃爍模擬旋轉感）
	var inner_pulse = 0.7 + 0.3 * sin(_inner_angle * 0.05)
	var outer_pulse = 0.7 + 0.3 * sin(_outer_angle * 0.04)
	if _inner_ring_node:
		_inner_ring_node.modulate = Color(inner_pulse, inner_pulse, 1.0, 1.0)
	if _outer_ring_node:
		_outer_ring_node.modulate = Color(1.0, outer_pulse, outer_pulse * 0.5, 1.0)

	# 超時自動停止（Server 會發送結果，這裡只是視覺提示）
	if _elapsed >= _spin_duration + 0.5:
		_is_active = false
		set_process(false)

# ---- 公開 API ----

## show_roulette 顯示輪盤（收到 dual_roulette_start 時呼叫）
func show_roulette(target_mult: float, base_reward: int, spin_duration: float,
		inner_ring: Array, outer_ring: Array) -> void:
	_target_mult = target_mult
	_base_reward = base_reward
	_spin_duration = spin_duration
	_inner_ring = inner_ring
	_outer_ring = outer_ring
	_elapsed = 0.0
	_inner_angle = 0.0
	_outer_angle = 0.0
	_is_active = true

	# 更新標題
	_title_label.text = "🎡 雙環輪盤（觸發：%.0fx）" % target_mult
	_countdown_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))

	# 隱藏結果面板，顯示停止按鈕
	_result_panel.visible = false
	_stop_btn.visible = true
	_stop_btn.disabled = false

	# 顯示面板（從下方滑入）
	visible = true
	_panel.position = Vector2(-240, 300)
	var tween = create_tween()
	tween.tween_property(_panel, "position", Vector2(-240, -210), 0.35).set_ease(Tween.EASE_OUT)

	# 全螢幕金色閃光
	_flash_overlay.color = Color(1.0, 0.84, 0.0, 0.5)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color", Color(1.0, 0.84, 0.0, 0.0), 0.4)

	set_process(true)

## show_result 顯示結果（收到 dual_roulette_result 時呼叫）
func show_result(inner_result: float, outer_result: float, combined: float,
		bonus_reward: int, _new_balance: int) -> void:
	_is_active = false
	set_process(false)

	# 停止按鈕禁用
	_stop_btn.disabled = true
	_stop_btn.visible = false

	# 顯示結果
	_result_label.text = "🎯 內環 %.0fx × 外環 %.0fx = %.0fx" % [inner_result, outer_result, combined]
	_bonus_label.text = "+%d 金幣！" % bonus_reward
	_result_panel.visible = true

	# 根據倍率決定顏色
	if combined >= 100.0:
		_result_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.8, 1.0))  # 粉紅（超高）
		_bonus_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.8, 1.0))
		_do_mega_flash()
	elif combined >= 50.0:
		_result_label.add_theme_color_override("font_color", COLOR_GOLD)
		_bonus_label.add_theme_color_override("font_color", COLOR_GOLD)
		_do_gold_flash()
	else:
		_result_label.add_theme_color_override("font_color", Color(0.2, 1.0, 0.4, 1.0))
		_bonus_label.add_theme_color_override("font_color", Color(0.2, 1.0, 0.4, 1.0))

	# 3 秒後自動關閉
	var close_timer = get_tree().create_timer(3.0)
	close_timer.timeout.connect(_close_panel)

## hide_panel 隱藏面板
func hide_panel() -> void:
	_is_active = false
	set_process(false)
	visible = false

# ---- 私有方法 ----

func _on_stop_pressed() -> void:
	if not _is_active:
		return
	_stop_btn.disabled = true
	# 通知 Server 停止
	var gm = get_node_or_null("/root/GameManager")
	if gm and gm.has_method("send_dual_roulette_stop"):
		gm.send_dual_roulette_stop()

func _close_panel() -> void:
	var tween = create_tween()
	tween.tween_property(_panel, "position", Vector2(-240, 300), 0.3).set_ease(Tween.EASE_IN)
	tween.tween_callback(func(): visible = false)

func _do_gold_flash() -> void:
	_flash_overlay.color = Color(1.0, 0.84, 0.0, 0.6)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.84, 0.0, 0.0), 0.5)

func _do_mega_flash() -> void:
	# 雙閃光（超高倍率）
	_flash_overlay.color = Color(1.0, 0.3, 0.8, 0.7)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.3, 0.8, 0.0), 0.3)
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.84, 0.0, 0.6), 0.1)
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.84, 0.0, 0.0), 0.4)
