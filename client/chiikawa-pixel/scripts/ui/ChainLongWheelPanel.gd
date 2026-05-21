## ChainLongWheelPanel.gd
## 千龍王強化輪盤面板（DAY-148）
## 業界依據：Royal Fishing JILI 2026「ChainLong King — capture this golden dragon to trigger
## the dual-ring roulette. The ChainLong King itself can award up to 1000X mega wins.」
## 設計：比普通雙環輪盤更震撼的視覺（金龍主題、更大面板、更強閃光）
## 內環（5x/10x/20x/30x/50x）× 外環（2x/3x/5x/7x/10x/20x）= 最高 1000x

extends Control

# ---- 狀態 ----
var _is_active: bool = false
var _spin_duration: float = 4.0
var _elapsed: float = 0.0
var _inner_ring: Array = [5.0, 10.0, 20.0, 30.0, 50.0]
var _outer_ring: Array = [2.0, 3.0, 5.0, 7.0, 10.0, 20.0]
var _inner_angle: float = 0.0
var _outer_angle: float = 0.0
var _inner_speed: float = 420.0  # 度/秒（比普通輪盤快）
var _outer_speed: float = 300.0  # 度/秒
var _base_reward: int = 0
var _target_mult: float = 0.0
var _killer_name: String = ""
var _is_personal: bool = false  # 是否為觸發玩家

# ---- 節點 ----
var _overlay: ColorRect = null
var _panel: Control = null
var _title_label: Label = null
var _subtitle_label: Label = null
var _inner_ring_node: Control = null
var _outer_ring_node: Control = null
var _inner_label: Label = null  # 內環當前值顯示
var _outer_label: Label = null  # 外環當前值顯示
var _stop_btn: Button = null
var _result_panel: Control = null
var _result_label: Label = null
var _bonus_label: Label = null
var _flash_overlay: ColorRect = null
var _dragon_label: Label = null  # 金龍裝飾
var _countdown_label: Label = null
var _observer_label: Label = null  # 旁觀者提示

# ---- 顏色（金龍主題）----
const COLOR_DRAGON_GOLD = Color(1.0, 0.84, 0.0, 1.0)
const COLOR_DRAGON_RED  = Color(0.9, 0.1, 0.1, 1.0)
const COLOR_INNER       = Color(0.9, 0.6, 0.0, 1.0)   # 金色內環
const COLOR_OUTER       = Color(0.8, 0.2, 0.0, 1.0)   # 紅色外環
const COLOR_BG          = Color(0.04, 0.02, 0.08, 0.98)

func _ready() -> void:
	_build_ui()
	set_process(false)
	visible = false
	mouse_filter = Control.MOUSE_FILTER_IGNORE

func _build_ui() -> void:
	# 全螢幕半透明遮罩（比普通輪盤更暗）
	_overlay = ColorRect.new()
	_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_overlay.color = Color(0.0, 0.0, 0.0, 0.88)
	_overlay.mouse_filter = Control.MOUSE_FILTER_STOP
	add_child(_overlay)

	# 全螢幕閃光層
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.84, 0.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 主面板（置中，比普通輪盤更大）
	_panel = Control.new()
	_panel.set_anchors_preset(Control.PRESET_CENTER)
	_panel.custom_minimum_size = Vector2(520, 480)
	_panel.position = Vector2(-260, -240)
	add_child(_panel)

	var panel_bg = ColorRect.new()
	panel_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel_bg.color = COLOR_BG
	_panel.add_child(panel_bg)

	# 金龍邊框（雙層）
	var border_outer = ColorRect.new()
	border_outer.set_anchors_preset(Control.PRESET_FULL_RECT)
	border_outer.color = COLOR_DRAGON_GOLD
	border_outer.custom_minimum_size = Vector2(520, 480)
	_panel.add_child(border_outer)

	var border_inner = ColorRect.new()
	border_inner.color = COLOR_BG
	border_inner.position = Vector2(3, 3)
	border_inner.size = Vector2(514, 474)
	_panel.add_child(border_inner)

	var border_inner2 = ColorRect.new()
	border_inner2.color = COLOR_DRAGON_RED
	border_inner2.position = Vector2(6, 6)
	border_inner2.size = Vector2(508, 468)
	_panel.add_child(border_inner2)

	var border_inner3 = ColorRect.new()
	border_inner3.color = COLOR_BG
	border_inner3.position = Vector2(9, 9)
	border_inner3.size = Vector2(502, 462)
	_panel.add_child(border_inner3)

	# 金龍裝飾
	_dragon_label = Label.new()
	_dragon_label.text = "🐉"
	_dragon_label.position = Vector2(0, 8)
	_dragon_label.size = Vector2(520, 40)
	_dragon_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_dragon_label.add_theme_font_size_override("font_size", 28)
	_panel.add_child(_dragon_label)

	# 標題
	_title_label = Label.new()
	_title_label.text = "千龍王強化輪盤"
	_title_label.position = Vector2(0, 44)
	_title_label.size = Vector2(520, 36)
	_title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_title_label.add_theme_color_override("font_color", COLOR_DRAGON_GOLD)
	_title_label.add_theme_font_size_override("font_size", 24)
	_panel.add_child(_title_label)

	# 副標題（最高倍率提示）
	_subtitle_label = Label.new()
	_subtitle_label.text = "最高 1000x 大獎！"
	_subtitle_label.position = Vector2(0, 80)
	_subtitle_label.size = Vector2(520, 24)
	_subtitle_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_subtitle_label.add_theme_color_override("font_color", COLOR_DRAGON_RED)
	_subtitle_label.add_theme_font_size_override("font_size", 14)
	_panel.add_child(_subtitle_label)

	# 倒數計時
	_countdown_label = Label.new()
	_countdown_label.text = "4.0s"
	_countdown_label.position = Vector2(0, 104)
	_countdown_label.size = Vector2(520, 24)
	_countdown_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_countdown_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_countdown_label.add_theme_font_size_override("font_size", 14)
	_panel.add_child(_countdown_label)

	# 外環（大圓，紅色）
	_outer_ring_node = Control.new()
	_outer_ring_node.position = Vector2(260 - 150, 130)
	_outer_ring_node.custom_minimum_size = Vector2(300, 300)
	_panel.add_child(_outer_ring_node)

	var outer_bg = ColorRect.new()
	outer_bg.color = Color(0.8, 0.2, 0.0, 0.3)
	outer_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	_outer_ring_node.add_child(outer_bg)

	# 外環標籤
	_outer_label = Label.new()
	_outer_label.text = "外環"
	_outer_label.position = Vector2(0, 130)
	_outer_label.size = Vector2(300, 40)
	_outer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_outer_label.add_theme_color_override("font_color", COLOR_OUTER)
	_outer_label.add_theme_font_size_override("font_size", 20)
	_outer_ring_node.add_child(_outer_label)

	# 內環（小圓，金色）
	_inner_ring_node = Control.new()
	_inner_ring_node.position = Vector2(260 - 95, 185)
	_inner_ring_node.custom_minimum_size = Vector2(190, 190)
	_panel.add_child(_inner_ring_node)

	var inner_bg = ColorRect.new()
	inner_bg.color = Color(0.9, 0.6, 0.0, 0.3)
	inner_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	_inner_ring_node.add_child(inner_bg)

	# 內環標籤
	_inner_label = Label.new()
	_inner_label.text = "內環"
	_inner_label.position = Vector2(0, 75)
	_inner_label.size = Vector2(190, 40)
	_inner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_inner_label.add_theme_color_override("font_color", COLOR_INNER)
	_inner_label.add_theme_font_size_override("font_size", 20)
	_inner_ring_node.add_child(_inner_label)

	# 停止按鈕（比普通輪盤更大更醒目）
	_stop_btn = Button.new()
	_stop_btn.text = "⏹ 停止！"
	_stop_btn.position = Vector2(160, 430)
	_stop_btn.size = Vector2(200, 44)
	_stop_btn.add_theme_color_override("font_color", Color(0.0, 0.0, 0.0, 1.0))
	_stop_btn.add_theme_font_size_override("font_size", 18)
	_stop_btn.pressed.connect(_on_stop_pressed)
	_panel.add_child(_stop_btn)

	# 旁觀者提示（非觸發玩家看到的）
	_observer_label = Label.new()
	_observer_label.text = ""
	_observer_label.position = Vector2(0, 430)
	_observer_label.size = Vector2(520, 44)
	_observer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_observer_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
	_observer_label.add_theme_font_size_override("font_size", 14)
	_observer_label.visible = false
	_panel.add_child(_observer_label)

	# 結果面板（停止後顯示）
	_result_panel = Control.new()
	_result_panel.position = Vector2(0, 240)
	_result_panel.size = Vector2(520, 140)
	_result_panel.visible = false
	_panel.add_child(_result_panel)

	_result_label = Label.new()
	_result_label.text = "🎯 內環 10x × 外環 5x = 50x"
	_result_label.position = Vector2(0, 10)
	_result_label.size = Vector2(520, 44)
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_label.add_theme_color_override("font_color", COLOR_DRAGON_GOLD)
	_result_label.add_theme_font_size_override("font_size", 20)
	_result_panel.add_child(_result_label)

	_bonus_label = Label.new()
	_bonus_label.text = "+5000 金幣！"
	_bonus_label.position = Vector2(0, 60)
	_bonus_label.size = Vector2(520, 60)
	_bonus_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_bonus_label.add_theme_color_override("font_color", Color(0.2, 1.0, 0.4, 1.0))
	_bonus_label.add_theme_font_size_override("font_size", 32)
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
	if remaining < 1.5:
		_countdown_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3, 1.0))

	# 更新環的旋轉視覺（金龍主題：金色+紅色脈衝）
	var inner_pulse = 0.6 + 0.4 * sin(_inner_angle * 0.06)
	var outer_pulse = 0.6 + 0.4 * sin(_outer_angle * 0.045)
	if _inner_ring_node:
		_inner_ring_node.modulate = Color(1.0, inner_pulse * 0.8, 0.0, 1.0)
	if _outer_ring_node:
		_outer_ring_node.modulate = Color(1.0, outer_pulse * 0.2, 0.0, 1.0)

	# 動態顯示內外環當前「指向」的值（視覺效果）
	var inner_idx = int(_inner_angle / 72.0) % len(_inner_ring)
	var outer_idx = int(_outer_angle / 60.0) % len(_outer_ring)
	_inner_label.text = "內環 %.0fx" % _inner_ring[inner_idx]
	_outer_label.text = "外環 %.0fx" % _outer_ring[outer_idx]

	# 超時後停止視覺更新
	if _elapsed >= _spin_duration + 0.5:
		_is_active = false
		set_process(false)

# ---- 公開 API ----

## show_wheel 顯示千龍王輪盤（收到 chainlong_wheel_start 時呼叫）
func show_wheel(killer_name: String, target_mult: float, base_reward: int,
		spin_duration: float, inner_ring: Array, outer_ring: Array, is_personal: bool) -> void:
	_killer_name = killer_name
	_target_mult = target_mult
	_base_reward = base_reward
	_spin_duration = spin_duration
	_inner_ring = inner_ring
	_outer_ring = outer_ring
	_is_personal = is_personal
	_elapsed = 0.0
	_inner_angle = 0.0
	_outer_angle = 0.0
	_is_active = true

	# 更新標題
	_title_label.text = "千龍王強化輪盤（%.0fx）" % target_mult
	_countdown_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))

	# 隱藏結果面板
	_result_panel.visible = false

	# 根據是否為觸發玩家決定顯示
	if is_personal:
		_stop_btn.visible = true
		_stop_btn.disabled = false
		_observer_label.visible = false
	else:
		_stop_btn.visible = false
		_observer_label.visible = true
		_observer_label.text = "🐉 %s 正在旋轉千龍王輪盤..." % killer_name

	# 顯示面板（從下方滑入，帶金龍閃光）
	visible = true
	_panel.position = Vector2(-260, 400)
	var tween = create_tween()
	tween.tween_property(_panel, "position", Vector2(-260, -240), 0.4).set_ease(Tween.EASE_OUT)

	# 全螢幕金龍閃光（比普通輪盤更強烈）
	_flash_overlay.color = Color(1.0, 0.5, 0.0, 0.7)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color", Color(1.0, 0.84, 0.0, 0.3), 0.2)
	flash_tween.tween_property(_flash_overlay, "color", Color(1.0, 0.84, 0.0, 0.0), 0.4)

	set_process(true)

## show_result 顯示結果（收到 chainlong_wheel_result 時呼叫）
func show_result(inner_result: float, outer_result: float, combined: float,
		bonus_reward: int, _new_balance: int, is_mega_win: bool, is_personal: bool) -> void:
	_is_active = false
	set_process(false)

	# 停止按鈕禁用
	_stop_btn.disabled = true
	_stop_btn.visible = false
	_observer_label.visible = false

	# 顯示結果
	_result_label.text = "🎯 內環 %.0fx × 外環 %.0fx = %.0fx" % [inner_result, outer_result, combined]
	if is_personal:
		_bonus_label.text = "+%d 金幣！" % bonus_reward
	else:
		_bonus_label.text = "%.0fx 大獎！" % combined
	_result_panel.visible = true

	# 根據倍率決定顏色和閃光強度
	if combined >= 500.0:
		# 傳說級（500x+）：彩虹閃光
		_result_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.8, 1.0))
		_bonus_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.8, 1.0))
		_do_legendary_flash()
	elif combined >= 200.0 or is_mega_win:
		# 大獎（200x+）：金龍閃光
		_result_label.add_theme_color_override("font_color", COLOR_DRAGON_GOLD)
		_bonus_label.add_theme_color_override("font_color", COLOR_DRAGON_GOLD)
		_do_mega_flash()
	elif combined >= 50.0:
		# 中獎（50x+）：金色閃光
		_result_label.add_theme_color_override("font_color", COLOR_DRAGON_GOLD)
		_bonus_label.add_theme_color_override("font_color", Color(0.2, 1.0, 0.4, 1.0))
		_do_gold_flash()
	else:
		_result_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
		_bonus_label.add_theme_color_override("font_color", Color(0.2, 1.0, 0.4, 1.0))

	# 非觸發玩家：3 秒後自動關閉
	# 觸發玩家：4 秒後自動關閉（讓他多看一下）
	var close_delay = 4.0 if is_personal else 3.0
	var close_timer = get_tree().create_timer(close_delay)
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
	if gm and gm.has_method("send_chainlong_wheel_stop"):
		gm.send_chainlong_wheel_stop()

func _close_panel() -> void:
	var tween = create_tween()
	tween.tween_property(_panel, "position", Vector2(-260, 400), 0.35).set_ease(Tween.EASE_IN)
	tween.tween_callback(func(): visible = false)

func _do_gold_flash() -> void:
	_flash_overlay.color = Color(1.0, 0.84, 0.0, 0.5)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.84, 0.0, 0.0), 0.5)

func _do_mega_flash() -> void:
	# 金龍雙閃光（大獎）
	_flash_overlay.color = Color(1.0, 0.5, 0.0, 0.8)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.5, 0.0, 0.0), 0.3)
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.84, 0.0, 0.6), 0.1)
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.84, 0.0, 0.0), 0.5)

func _do_legendary_flash() -> void:
	# 傳說三閃光（500x+）
	_flash_overlay.color = Color(1.0, 0.3, 0.8, 0.9)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.3, 0.8, 0.0), 0.2)
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.84, 0.0, 0.7), 0.1)
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.84, 0.0, 0.0), 0.2)
	tween.tween_property(_flash_overlay, "color", Color(0.5, 0.3, 1.0, 0.6), 0.1)
	tween.tween_property(_flash_overlay, "color", Color(0.5, 0.3, 1.0, 0.0), 0.4)
