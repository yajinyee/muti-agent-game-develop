## LuckyCloverPanel.gd — 幸運草魚 UI 面板（DAY-179）
## 業界依據：Ocean King 3 Plus「Lucky Shamrock Leprechaun Boss」
## 顯示幸運草爆發、全服加成、幸運草金幣發放
extends CanvasLayer

# ---- 常數 ----
const PANEL_COLOR_BG     := Color(0.0, 0.1, 0.05, 0.90)
const PANEL_COLOR_GREEN  := Color(0.0, 1.0, 0.5, 1.0)    # 春綠色（幸運草感）
const PANEL_COLOR_GOLD   := Color(1.0, 0.85, 0.0, 1.0)
const PANEL_COLOR_WHITE  := Color(1.0, 1.0, 1.0, 1.0)

# ---- 節點引用 ----
var _banner_container  : Control
var _banner_label      : Label
var _boost_timer_label : Label
var _gift_popup        : Control
var _gift_label        : Label
var _flash_overlay     : ColorRect

# ---- 狀態 ----
var _boost_active      : bool = false
var _boost_remaining   : float = 0.0

func _ready() -> void:
	layer = 65
	_build_ui()
	hide()

func _build_ui() -> void:
	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(0.0, 1.0, 0.5, 0.0)
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
	banner_style.bg_color = Color(0.0, 0.35, 0.15, 0.92)
	banner_style.corner_radius_top_left = 8
	banner_style.corner_radius_top_right = 8
	banner_style.corner_radius_bottom_left = 8
	banner_style.corner_radius_bottom_right = 8
	_banner_container.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner_container)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_GREEN)
	_banner_label.add_theme_font_size_override("font_size", 20)
	_banner_container.add_child(_banner_label)

	# 加成計時器（右上角）
	_boost_timer_label = Label.new()
	_boost_timer_label.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	_boost_timer_label.offset_top = 64
	_boost_timer_label.offset_right = -16
	_boost_timer_label.offset_left = -200
	_boost_timer_label.offset_bottom = 96
	_boost_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	_boost_timer_label.add_theme_color_override("font_color", PANEL_COLOR_GREEN)
	_boost_timer_label.add_theme_font_size_override("font_size", 18)
	_boost_timer_label.hide()
	add_child(_boost_timer_label)

	# 幸運草金幣彈窗（中央）
	_gift_popup = PanelContainer.new()
	_gift_popup.set_anchors_preset(Control.PRESET_CENTER)
	_gift_popup.offset_left = -150
	_gift_popup.offset_right = 150
	_gift_popup.offset_top = -70
	_gift_popup.offset_bottom = 70
	var popup_style := StyleBoxFlat.new()
	popup_style.bg_color = PANEL_COLOR_BG
	popup_style.corner_radius_top_left = 16
	popup_style.corner_radius_top_right = 16
	popup_style.corner_radius_bottom_left = 16
	popup_style.corner_radius_bottom_right = 16
	popup_style.border_width_left = 2
	popup_style.border_width_right = 2
	popup_style.border_width_top = 2
	popup_style.border_width_bottom = 2
	popup_style.border_color = PANEL_COLOR_GREEN
	_gift_popup.add_theme_stylebox_override("panel", popup_style)
	_gift_popup.modulate.a = 0.0
	add_child(_gift_popup)

	_gift_label = Label.new()
	_gift_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_gift_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_gift_label.add_theme_color_override("font_color", PANEL_COLOR_WHITE)
	_gift_label.add_theme_font_size_override("font_size", 22)
	_gift_popup.add_child(_gift_label)

func _process(delta: float) -> void:
	if _boost_active:
		_boost_remaining -= delta
		if _boost_remaining <= 0.0:
			_boost_active = false
			_boost_timer_label.hide()
		else:
			_boost_timer_label.text = "🍀 +50%% %.1f 秒" % _boost_remaining

## handle_lucky_clover — 處理幸運草魚訊息
func handle_lucky_clover(payload: Dictionary) -> void:
	var phase : String = payload.get("phase", "")
	match phase:
		"clover_start":
			_on_clover_start(payload)
		"clover_gift":
			_on_clover_gift(payload)
		"clover_end":
			_on_clover_end()

## _on_clover_start — 幸運草爆發開始（全服）
func _on_clover_start(payload: Dictionary) -> void:
	_boost_active = true
	_boost_remaining = float(payload.get("boost_duration_sec", 10))
	var player_name : String = payload.get("player_name", "")

	show()

	# 全螢幕綠色閃光
	_flash_screen(Color(0.0, 1.0, 0.5, 0.4), 0.4)

	# 橫幅
	_banner_label.text = "🍀 %s 觸發幸運草爆發！全服 +50%% 加成！" % player_name
	_banner_container.modulate.a = 0.0
	_banner_container.show()
	var tween := create_tween()
	tween.tween_property(_banner_container, "modulate:a", 1.0, 0.3)

	# 加成計時器
	_boost_timer_label.show()

## _on_clover_gift — 幸運草金幣發放（個人）
func _on_clover_gift(payload: Dictionary) -> void:
	var gift_amount : int = payload.get("gift_amount", 0)
	var gift_mult : int = payload.get("gift_mult", 10)

	# 彈窗彈跳
	_gift_label.text = "🍀 幸運草金幣！\n+%d 金幣\n（×%d 倍）" % [gift_amount, gift_mult]
	_gift_popup.scale = Vector2(0.5, 0.5)
	_gift_popup.modulate.a = 1.0
	var tween := create_tween()
	tween.tween_property(_gift_popup, "scale", Vector2(1.0, 1.0), 0.3).set_trans(Tween.TRANS_BACK)

	# 浮動文字
	_spawn_float_text("🍀 +%d" % gift_amount, PANEL_COLOR_GOLD)

	# 3 秒後淡出彈窗
	await get_tree().create_timer(3.0).timeout
	var fade := create_tween()
	fade.tween_property(_gift_popup, "modulate:a", 0.0, 0.4)

## _on_clover_end — 幸運草爆發結束（全服）
func _on_clover_end() -> void:
	_boost_active = false
	_boost_timer_label.hide()
	_banner_container.hide()
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
	lbl.add_theme_font_size_override("font_size", 22)
	lbl.set_anchors_preset(Control.PRESET_CENTER)
	lbl.offset_left = -80
	lbl.offset_right = 80
	lbl.offset_top = -20
	lbl.offset_bottom = 20
	add_child(lbl)
	var tween := create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 55, 0.9)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.9)
	await tween.finished
	lbl.queue_free()
