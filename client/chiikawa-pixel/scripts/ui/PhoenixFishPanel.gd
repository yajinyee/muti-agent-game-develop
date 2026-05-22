## PhoenixFishPanel.gd — 鳳凰魚涅槃重生 UI 面板（DAY-185）
## 業界依據：Ocean King 3 Plus「Phoenix Fish — when defeated, the Phoenix Fish triggers a
## rebirth explosion that deals massive damage to all fish on screen, with the Phoenix
## rising from the ashes to grant a 30-second luck boost」
## 顯示涅槃爆炸（全場火焰）、鳳凰重生橫幅（+30% 倒數計時）、結果彈窗
extends CanvasLayer

# ---- 常數 ----
const PANEL_COLOR_BG      := Color(0.10, 0.02, 0.0, 0.92)
const PANEL_COLOR_FIRE    := Color(1.0, 0.35, 0.0, 1.0)   # 火焰橙紅
const PANEL_COLOR_GOLD    := Color(1.0, 0.85, 0.0, 1.0)   # 金色（重生感）
const PANEL_COLOR_WHITE   := Color(1.0, 1.0, 1.0, 1.0)
const PANEL_COLOR_REBIRTH := Color(1.0, 0.6, 0.1, 1.0)    # 重生橙

# ---- 節點引用 ----
var _banner_container  : Control
var _banner_label      : Label
var _rebirth_bar       : Control   # 重生加成倒數計時條
var _rebirth_fill      : ColorRect
var _rebirth_label     : Label     # 「+30% 重生加成」文字
var _result_panel      : Control
var _result_label      : Label
var _flash_overlay     : ColorRect

# ---- 狀態 ----
var _rebirth_active    : bool = false
var _rebirth_remaining : float = 0.0
var _rebirth_total     : float = 30.0

func _ready() -> void:
	layer = 60
	_build_ui()
	hide()

func _build_ui() -> void:
	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 0.35, 0.0, 0.0)
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
	banner_style.bg_color = Color(0.18, 0.04, 0.0, 0.92)
	banner_style.corner_radius_top_left = 8
	banner_style.corner_radius_top_right = 8
	banner_style.corner_radius_bottom_left = 8
	banner_style.corner_radius_bottom_right = 8
	banner_style.border_width_left = 2
	banner_style.border_width_right = 2
	banner_style.border_width_top = 2
	banner_style.border_width_bottom = 2
	banner_style.border_color = PANEL_COLOR_FIRE
	_banner_container.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner_container)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	_banner_label.add_theme_font_size_override("font_size", 20)
	_banner_container.add_child(_banner_label)

	# 重生加成進度條（底部，顯示剩餘時間）
	_rebirth_bar = Control.new()
	_rebirth_bar.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	_rebirth_bar.offset_bottom = 0
	_rebirth_bar.offset_top = -20
	_rebirth_bar.offset_left = 0
	_rebirth_bar.offset_right = 0
	add_child(_rebirth_bar)

	var bar_bg := ColorRect.new()
	bar_bg.color = Color(0.1, 0.05, 0.0, 0.85)
	bar_bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_rebirth_bar.add_child(bar_bg)

	_rebirth_fill = ColorRect.new()
	_rebirth_fill.color = PANEL_COLOR_FIRE
	_rebirth_fill.set_anchors_preset(Control.PRESET_LEFT_WIDE)
	_rebirth_fill.offset_right = 0
	_rebirth_bar.add_child(_rebirth_fill)

	_rebirth_label = Label.new()
	_rebirth_label.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_rebirth_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_rebirth_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_rebirth_label.add_theme_color_override("font_color", PANEL_COLOR_WHITE)
	_rebirth_label.add_theme_font_size_override("font_size", 13)
	_rebirth_label.text = "🔥 鳳凰重生 +30%"
	_rebirth_bar.add_child(_rebirth_label)

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
	result_style.border_color = PANEL_COLOR_FIRE
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
	if _rebirth_active:
		_rebirth_remaining -= delta
		if _rebirth_remaining <= 0.0:
			_rebirth_active = false
			_rebirth_bar.hide()
		else:
			# 更新進度條
			var pct := _rebirth_remaining / _rebirth_total
			_rebirth_fill.offset_right = _rebirth_bar.size.x * pct - _rebirth_bar.size.x
			# 顏色從金→橙→紅（時間越少越緊迫）
			if pct > 0.6:
				_rebirth_fill.color = PANEL_COLOR_GOLD
			elif pct > 0.3:
				_rebirth_fill.color = PANEL_COLOR_REBIRTH
			else:
				_rebirth_fill.color = PANEL_COLOR_FIRE
			# 更新倒數文字
			_rebirth_label.text = "🔥 鳳凰重生 +30%%  剩餘 %ds" % int(_rebirth_remaining)

## handle_phoenix_fish — 處理鳳凰魚涅槃重生訊息
func handle_phoenix_fish(payload: Dictionary) -> void:
	var phase : String = payload.get("phase", "")
	if phase == "phoenix_explode":
		_on_phoenix_explode(payload)
	elif phase == "phoenix_rebirth":
		_on_phoenix_rebirth(payload)
	elif phase == "rebirth_end":
		_on_rebirth_end(payload)

## _on_phoenix_explode — 涅槃爆炸
func _on_phoenix_explode(payload: Dictionary) -> void:
	var killer_name : String = payload.get("killer_name", "")

	show()

	# 全螢幕火焰爆炸（四次漸強閃光，製造「全場燃燒」感）
	_flash_screen(Color(1.0, 0.35, 0.0, 0.4), 0.1)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color(1.0, 0.35, 0.0, 0.6), 0.12)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(Color(1.0, 0.5, 0.0, 0.75), 0.15)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(Color(1.0, 0.7, 0.0, 0.5), 0.3)

	# 橫幅（爆炸中）
	_banner_label.text = "🔥 %s 的鳳凰魚涅槃爆炸！全場燃燒！" % killer_name
	_banner_container.modulate.a = 0.0
	_banner_container.show()
	var tween := create_tween()
	tween.tween_property(_banner_container, "modulate:a", 1.0, 0.25)

	# 在全場生成火焰粒子（8個隨機位置）
	for i in range(8):
		var rx := randf_range(100, 1180)
		var ry := randf_range(80, 620)
		_spawn_fire_at(Vector2(rx, ry))
		await get_tree().create_timer(0.04).timeout

## _on_phoenix_rebirth — 鳳凰重生（爆炸結束，加成開始）
func _on_phoenix_rebirth(payload: Dictionary) -> void:
	var total_kills  : int = payload.get("total_kills", 0)
	var total_reward : int = payload.get("total_reward", 0)
	var boost_sec    : int = payload.get("boost_sec", 30)
	var killer_name  : String = payload.get("killer_name", "")

	# 激活重生加成計時
	_rebirth_active = true
	_rebirth_total = float(boost_sec)
	_rebirth_remaining = float(boost_sec)
	_rebirth_bar.show()

	# 橫幅更新（重生中）
	_banner_label.text = "🔥 鳳凰重生！全服 +30%% 加成 %d 秒！" % boost_sec
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_GOLD)

	# 金色閃光（重生感）
	_flash_screen(Color(1.0, 0.85, 0.0, 0.5), 0.4)

	# 結果彈窗
	_result_label.text = "🔥 鳳凰涅槃\n擊破：%d 個\n獎勵：+%d\n全服 +30%% 加成！" % [total_kills, total_reward]
	_result_panel.offset_left = 0
	_result_panel.modulate.a = 1.0
	var tween := create_tween()
	tween.tween_property(_result_panel, "offset_left", -300.0, 0.4).set_trans(Tween.TRANS_BACK)

	# 高擊破數特效
	if total_kills >= 8:
		await get_tree().create_timer(0.4).timeout
		_flash_screen(Color(1.0, 0.85, 0.0, 0.7), 0.3)
		await get_tree().create_timer(0.3).timeout
		_flash_screen(Color(1.0, 0.85, 0.0, 0.5), 0.4)
	elif total_kills >= 5:
		await get_tree().create_timer(0.4).timeout
		_flash_screen(Color(1.0, 0.5, 0.0, 0.5), 0.3)

	# 3 秒後結果彈窗淡出（重生條繼續顯示）
	await get_tree().create_timer(3.5).timeout
	var fade := create_tween()
	fade.tween_property(_result_panel, "modulate:a", 0.0, 0.5)

## _on_rebirth_end — 重生加成結束
func _on_rebirth_end(_payload: Dictionary) -> void:
	_rebirth_active = false
	_rebirth_remaining = 0.0

	# 橫幅更新（重生結束）
	_banner_label.text = "🔥 鳳凰重生加成結束"
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_FIRE)

	# 淡出所有 UI
	await get_tree().create_timer(1.5).timeout
	var fade := create_tween()
	fade.tween_property(_banner_container, "modulate:a", 0.0, 0.5)
	fade.parallel().tween_property(_rebirth_bar, "modulate:a", 0.0, 0.5)
	await fade.finished
	_banner_container.hide()
	_rebirth_bar.hide()
	_rebirth_bar.modulate.a = 1.0
	hide()

## _spawn_fire_at — 在指定位置生成火焰符號
func _spawn_fire_at(pos: Vector2) -> void:
	var lbl := Label.new()
	lbl.text = "🔥"
	lbl.add_theme_font_size_override("font_size", 32)
	lbl.position = pos - Vector2(16, 16)
	lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(lbl)
	var tween := create_tween()
	tween.tween_property(lbl, "scale", Vector2(1.6, 1.6), 0.12)
	tween.tween_property(lbl, "position:y", pos.y - 40, 0.3)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.3)
	await tween.finished
	if is_instance_valid(lbl):
		lbl.queue_free()

## _flash_screen — 全螢幕閃光
func _flash_screen(color: Color, duration: float) -> void:
	_flash_overlay.color = color
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)
