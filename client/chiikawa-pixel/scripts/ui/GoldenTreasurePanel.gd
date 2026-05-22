## GoldenTreasurePanel.gd — 黃金寶藏魚 UI 面板（DAY-177）
## 業界依據：Ocean King 3 Plus「Golden Treasure feature」
## 顯示寶藏箱出現、玩家開箱互動、獎勵結果
extends CanvasLayer

# ---- 常數 ----
const PANEL_COLOR_BG     := Color(0.08, 0.06, 0.0, 0.92)
const PANEL_COLOR_GOLD   := Color(1.0, 0.85, 0.0, 1.0)
const PANEL_COLOR_WHITE  := Color(1.0, 1.0, 1.0, 1.0)
const PANEL_COLOR_GREEN  := Color(0.2, 0.9, 0.2, 1.0)
const PANEL_COLOR_BLUE   := Color(0.3, 0.7, 1.0, 1.0)
const CHEST_COUNT        := 3

# ---- 節點引用 ----
var _banner_container  : Control
var _banner_label      : Label
var _chest_container   : Control
var _chest_buttons     : Array = []
var _timer_label       : Label
var _result_panel      : Control
var _result_label      : Label
var _mult_timer_label  : Label
var _flash_overlay     : ColorRect

# ---- 狀態 ----
var _is_my_session     : bool = false
var _timeout_sec       : int = 12
var _elapsed           : float = 0.0
var _session_active    : bool = false
var _mult_active       : bool = false
var _mult_remaining    : float = 0.0
var _opened_chests     : Array = []

func _ready() -> void:
	layer = 67
	_build_ui()
	hide()

func _build_ui() -> void:
	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 0.85, 0.0, 0.0)
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
	banner_style.bg_color = Color(0.4, 0.3, 0.0, 0.92)
	banner_style.corner_radius_top_left = 8
	banner_style.corner_radius_top_right = 8
	banner_style.corner_radius_bottom_left = 8
	banner_style.corner_radius_bottom_right = 8
	_banner_container.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner_container)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	_banner_label.add_theme_font_size_override("font_size", 20)
	_banner_container.add_child(_banner_label)

	# 倒數計時器（右上角）
	_timer_label = Label.new()
	_timer_label.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	_timer_label.offset_top = 64
	_timer_label.offset_right = -16
	_timer_label.offset_left = -160
	_timer_label.offset_bottom = 96
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	_timer_label.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	_timer_label.add_theme_font_size_override("font_size", 18)
	add_child(_timer_label)

	# 寶藏箱容器（中央）
	_chest_container = HBoxContainer.new()
	_chest_container.set_anchors_preset(Control.PRESET_CENTER)
	_chest_container.offset_left = -200
	_chest_container.offset_right = 200
	_chest_container.offset_top = -60
	_chest_container.offset_bottom = 60
	_chest_container.alignment = BoxContainer.ALIGNMENT_CENTER
	_chest_container.add_theme_constant_override("separation", 24)
	add_child(_chest_container)

	# 建立 3 個寶藏箱按鈕
	for i in range(CHEST_COUNT):
		var btn := Button.new()
		btn.text = "💰\n寶藏 %d" % (i + 1)
		btn.custom_minimum_size = Vector2(100, 100)
		btn.add_theme_font_size_override("font_size", 18)
		var btn_style := StyleBoxFlat.new()
		btn_style.bg_color = Color(0.3, 0.22, 0.0, 0.95)
		btn_style.corner_radius_top_left = 12
		btn_style.corner_radius_top_right = 12
		btn_style.corner_radius_bottom_left = 12
		btn_style.corner_radius_bottom_right = 12
		btn_style.border_width_left = 3
		btn_style.border_width_right = 3
		btn_style.border_width_top = 3
		btn_style.border_width_bottom = 3
		btn_style.border_color = PANEL_COLOR_GOLD
		btn.add_theme_stylebox_override("normal", btn_style)
		btn.pressed.connect(_on_chest_pressed.bind(i))
		_chest_container.add_child(btn)
		_chest_buttons.append(btn)

	# 倍率計時器（左上角，倍率激活時顯示）
	_mult_timer_label = Label.new()
	_mult_timer_label.set_anchors_preset(Control.PRESET_TOP_LEFT)
	_mult_timer_label.offset_top = 64
	_mult_timer_label.offset_left = 16
	_mult_timer_label.offset_right = 200
	_mult_timer_label.offset_bottom = 96
	_mult_timer_label.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	_mult_timer_label.add_theme_font_size_override("font_size", 18)
	_mult_timer_label.hide()
	add_child(_mult_timer_label)

	# 結果面板（右側滑入）
	_result_panel = PanelContainer.new()
	_result_panel.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_panel.offset_right = -16
	_result_panel.offset_left = -280
	_result_panel.offset_top = -100
	_result_panel.offset_bottom = 100
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
	result_style.border_color = PANEL_COLOR_GOLD
	_result_panel.add_theme_stylebox_override("panel", result_style)
	_result_panel.modulate.a = 0.0
	add_child(_result_panel)

	_result_label = Label.new()
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_label.add_theme_color_override("font_color", PANEL_COLOR_WHITE)
	_result_label.add_theme_font_size_override("font_size", 16)
	_result_panel.add_child(_result_label)

func _process(delta: float) -> void:
	# 更新開箱倒數計時
	if _session_active and _is_my_session:
		_elapsed += delta
		var remaining := max(0.0, _timeout_sec - _elapsed)
		_timer_label.text = "⏱ %.1f 秒" % remaining
		if remaining <= 3.0:
			_timer_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3, 1.0))

	# 更新倍率倒數計時
	if _mult_active:
		_mult_remaining -= delta
		if _mult_remaining <= 0.0:
			_mult_active = false
			_mult_timer_label.hide()
		else:
			_mult_timer_label.text = "×3 倍率 %.1f 秒" % _mult_remaining

## handle_golden_treasure — 處理黃金寶藏魚訊息
func handle_golden_treasure(payload: Dictionary) -> void:
	var phase : String = payload.get("phase", "")
	match phase:
		"treasure_start":
			_on_treasure_start(payload)
		"treasure_broadcast":
			_on_treasure_broadcast(payload)
		"treasure_open":
			_on_treasure_open(payload)
		"treasure_mult_start":
			_on_mult_start(payload)
		"treasure_mult_end":
			_on_mult_end()
		"treasure_auto_open":
			_on_auto_open(payload)
		"treasure_end":
			_on_treasure_end()
		"treasure_weapon_charge":
			_on_weapon_charge(payload)

## _on_treasure_start — 寶藏箱出現（個人）
func _on_treasure_start(payload: Dictionary) -> void:
	_is_my_session = true
	_timeout_sec = payload.get("timeout_sec", 12)	_elapsed = 0.0
	_session_active = true
	_opened_chests.clear()

	# 重置箱子按鈕
	for i in range(CHEST_COUNT):
		_chest_buttons[i].text = "💰\n寶藏 %d" % (i + 1)
		_chest_buttons[i].disabled = false
		_chest_buttons[i].modulate = Color.WHITE

	show()

	# 全螢幕金色閃光
	_flash_screen(Color(1.0, 0.85, 0.0, 0.5), 0.4)

	# 橫幅
	_banner_label.text = "💰 黃金寶藏！點擊開啟寶藏箱！"
	_banner_container.modulate.a = 0.0
	_banner_container.show()
	var tween := create_tween()
	tween.tween_property(_banner_container, "modulate:a", 1.0, 0.3)

	# 寶藏箱彈入動畫
	_chest_container.show()
	for i in range(CHEST_COUNT):
		var btn = _chest_buttons[i]
		btn.scale = Vector2(0.3, 0.3)
		var btn_tween := create_tween()
		btn_tween.tween_interval(i * 0.15)
		btn_tween.tween_property(btn, "scale", Vector2(1.0, 1.0), 0.3).set_trans(Tween.TRANS_BACK)

	# 倒數計時
	_timer_label.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	_timer_label.show()

## _on_treasure_broadcast — 全服廣播（其他玩家看到）
func _on_treasure_broadcast(payload: Dictionary) -> void:
	if _is_my_session:
		return  # 自己的 session 已由 treasure_start 處理
	var player_name : String = payload.get("player_name", "")
	_show_broadcast_banner("💰 %s 觸發了黃金寶藏！" % player_name)

## _on_treasure_open — 開箱結果
func _on_treasure_open(payload: Dictionary) -> void:
	var chest_id : int = payload.get("chest_id", 0)
	var reward_type : String = payload.get("reward_type", "coins")
	var reward : int = payload.get("reward", 0)

	_opened_chests.append(chest_id)

	# 更新箱子按鈕顯示
	if chest_id < _chest_buttons.size():
		var btn = _chest_buttons[chest_id]
		btn.disabled = true
		match reward_type:
			"coins":
				btn.text = "🪙\n+%d" % reward
				btn.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
				_spawn_float_text("+%d 金幣" % reward, PANEL_COLOR_GOLD)
			"mult":
				btn.text = "⚡\n×3 倍率"
				btn.add_theme_color_override("font_color", Color(0.5, 1.0, 0.5, 1.0))
				_spawn_float_text("×3 倍率激活！", Color(0.5, 1.0, 0.5, 1.0))
			"weapon":
				btn.text = "🔫\n武器充能"
				btn.add_theme_color_override("font_color", PANEL_COLOR_BLUE)
				_spawn_float_text("武器充能！", PANEL_COLOR_BLUE)

	# 小閃光
	_flash_screen(Color(1.0, 0.85, 0.0, 0.2), 0.2)

## _on_mult_start — 倍率加成激活
func _on_mult_start(payload: Dictionary) -> void:
	_mult_active = true
	_mult_remaining = float(payload.get("mult_duration_sec", 8))
	_mult_timer_label.show()
	_mult_timer_label.text = "×3 倍率 %.1f 秒" % _mult_remaining

## _on_mult_end — 倍率加成結束
func _on_mult_end() -> void:
	_mult_active = false
	_mult_timer_label.hide()

## _on_auto_open — 超時自動開啟
func _on_auto_open(payload: Dictionary) -> void:
	var chest_id : int = payload.get("chest_id", 0)
	var reward : int = payload.get("reward", 0)
	if chest_id < _chest_buttons.size():
		var btn = _chest_buttons[chest_id]
		btn.disabled = true
		btn.text = "🪙\n+%d\n(自動)" % reward
		btn.modulate = Color(0.7, 0.7, 0.7, 1.0)

## _on_treasure_end — 寶藏結束
func _on_treasure_end() -> void:
	_session_active = false
	_is_my_session = false
	_timer_label.hide()
	_chest_container.hide()

	# 結果面板
	var total_opened := _opened_chests.size()
	_result_label.text = "💰 寶藏結束\n開啟：%d / %d 個箱子" % [total_opened, CHEST_COUNT]
	var tween := create_tween()
	tween.tween_property(_result_panel, "modulate:a", 1.0, 0.3)

	await get_tree().create_timer(3.0).timeout
	_fade_out()

## _on_weapon_charge — 武器充能
func _on_weapon_charge(_payload: Dictionary) -> void:
	_spawn_float_text("🔫 武器充能！", PANEL_COLOR_BLUE)

## _on_chest_pressed — 玩家點擊寶藏箱
func _on_chest_pressed(chest_id: int) -> void:
	if not _is_my_session or not _session_active:
		return
	if chest_id in _opened_chests:
		return
	# 發送開箱請求（優先透過 GameManager）
	if GameManager.has_method("send_golden_treasure_open"):
		GameManager.send_golden_treasure_open(chest_id)
	elif NetworkManager.has_method("send"):
		NetworkManager.send("golden_treasure_open", {"chest_id": chest_id})

## _show_broadcast_banner — 顯示全服廣播橫幅（短暫）
func _show_broadcast_banner(text: String) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	lbl.add_theme_font_size_override("font_size", 16)
	lbl.set_anchors_preset(Control.PRESET_TOP_WIDE)
	lbl.offset_top = 60
	lbl.offset_bottom = 90
	lbl.modulate.a = 0.0
	add_child(lbl)
	show()
	var tween := create_tween()
	tween.tween_property(lbl, "modulate:a", 1.0, 0.3)
	await get_tree().create_timer(2.5).timeout
	var fade := create_tween()
	fade.tween_property(lbl, "modulate:a", 0.0, 0.4)
	await fade.finished
	lbl.queue_free()
	if not _is_my_session:
		hide()

## _flash_screen — 全螢幕閃光
func _flash_screen(color: Color, duration: float) -> void:
	_flash_overlay.color = color
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)

## _spawn_float_text — 浮動文字
func _spawn_float_text(text: String, color: Color) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 20)
	lbl.set_anchors_preset(Control.PRESET_CENTER)
	lbl.offset_left = -80
	lbl.offset_right = 80
	lbl.offset_top = -20
	lbl.offset_bottom = 20
	add_child(lbl)
	var tween := create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 50, 0.9)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.9)
	await tween.finished
	lbl.queue_free()

## _fade_out — 淡出所有 UI
func _fade_out() -> void:
	var tween := create_tween()
	tween.tween_property(self, "modulate:a", 0.0, 0.4)
	await tween.finished
	modulate.a = 1.0
	_banner_container.hide()
	_result_panel.modulate.a = 0.0
	hide()
