## DragonTurtlePanel.gd — 龍龜不死 Boss UI 面板（DAY-186）
## 業界依據：Royal Fishing JILI「Immortal Boss mechanic — Golden Toad and Ancient Crocodile
## bosses appear randomly and award consecutive wins ranging from 50X to 150X until they
## leave the screen. This creates extended winning sequences impossible in standard fish games.」
## 顯示龍龜出現橫幅、命中獎勵浮動文字、離開結算彈窗
## Phase: turtle_appear → turtle_hit / my_hit → turtle_leave
extends CanvasLayer

# ---- 常數 ----
const PANEL_COLOR_BG     := Color(0.0, 0.12, 0.06, 0.92)  # 深綠（龜殼感）
const PANEL_COLOR_GREEN  := Color(0.2, 0.9, 0.3, 1.0)     # 亮綠（龜殼高光）
const PANEL_COLOR_GOLD   := Color(1.0, 0.85, 0.0, 1.0)    # 金色（獎勵感）
const PANEL_COLOR_WHITE  := Color(1.0, 1.0, 1.0, 1.0)
const PANEL_COLOR_TEAL   := Color(0.0, 0.75, 0.6, 1.0)    # 青綠（龜殼紋路）

# ---- 節點引用 ----
var _banner_container : Control
var _banner_label     : Label
var _hit_counter      : Label   # 全服命中計數器
var _result_panel     : Control
var _result_label     : Label
var _flash_overlay    : ColorRect

# ---- 狀態 ----
var _is_active        : bool = false
var _total_hits       : int  = 0
var _total_reward     : int  = 0
var _instance_id      : String = ""

func _ready() -> void:
	layer = 59  # 比 PhoenixFishPanel(60) 低一層
	_build_ui()
	hide()

func _build_ui() -> void:
	# 全螢幕閃光 overlay（龜殼綠色）
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(0.2, 0.9, 0.3, 0.0)
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
	banner_style.bg_color = Color(0.0, 0.15, 0.08, 0.92)
	banner_style.corner_radius_top_left = 8
	banner_style.corner_radius_top_right = 8
	banner_style.corner_radius_bottom_left = 8
	banner_style.corner_radius_bottom_right = 8
	banner_style.border_width_left = 2
	banner_style.border_width_right = 2
	banner_style.border_width_top = 2
	banner_style.border_width_bottom = 2
	banner_style.border_color = PANEL_COLOR_GREEN
	_banner_container.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner_container)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	_banner_label.add_theme_font_size_override("font_size", 20)
	_banner_container.add_child(_banner_label)

	# 命中計數器（右上角）
	_hit_counter = Label.new()
	_hit_counter.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	_hit_counter.offset_top = 8
	_hit_counter.offset_right = -8
	_hit_counter.offset_left = -200
	_hit_counter.offset_bottom = 40
	_hit_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	_hit_counter.add_theme_color_override("font_color", PANEL_COLOR_GREEN)
	_hit_counter.add_theme_font_size_override("font_size", 16)
	add_child(_hit_counter)

	# 結果彈窗（右側滑入）
	_result_panel = PanelContainer.new()
	_result_panel.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_panel.offset_right = 0
	_result_panel.offset_left = -280
	_result_panel.offset_top = -80
	_result_panel.offset_bottom = 80
	var result_style := StyleBoxFlat.new()
	result_style.bg_color = Color(0.0, 0.12, 0.06, 0.95)
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

## handle_dragon_turtle 處理龍龜不死 Boss 訊息
func handle_dragon_turtle(payload: Dictionary) -> void:
	var phase : String = payload.get("phase", "")
	match phase:
		"turtle_appear":
			_on_turtle_appear(payload)
		"turtle_hit":
			_on_turtle_hit(payload)
		"my_hit":
			_on_my_hit(payload)
		"turtle_leave":
			_on_turtle_leave(payload)

# ---- 私有處理函數 ----

func _on_turtle_appear(payload: Dictionary) -> void:
	_is_active = true
	_total_hits = 0
	_total_reward = 0
	_instance_id = payload.get("instance_id", "")

	show()
	_result_panel.visible = false
	_banner_container.visible = true
	_hit_counter.visible = true
	_hit_counter.text = "🐢 命中：0 次"

	# 橫幅文字
	_banner_label.text = "🐢 龍龜不死 Boss 出現！命中即得 50-150x 獎勵！"

	# 綠色閃光（龜殼感）
	_flash_green(0.5)

	# 橫幅從頂部滑入
	_banner_container.offset_top = -60
	_banner_container.offset_bottom = -12
	var tween := create_tween()
	tween.tween_property(_banner_container, "offset_top", 8.0, 0.3).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(_banner_container, "offset_bottom", 56.0, 0.3).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

func _on_turtle_hit(payload: Dictionary) -> void:
	if not _is_active:
		return

	_total_hits = payload.get("total_hits", _total_hits)
	_total_reward = payload.get("total_reward", _total_reward)

	# 更新命中計數器
	_hit_counter.text = "🐢 命中：%d 次" % _total_hits

	# 命中玩家名稱 + 獎勵浮動文字
	var hitter_name : String = payload.get("hitter_name", "")
	var hit_mult    : int    = payload.get("hit_mult", 0)
	if hitter_name != "" and hit_mult > 0:
		_spawn_hit_float_text(hitter_name, hit_mult)

	# 小閃光
	_flash_green(0.2)

func _on_my_hit(payload: Dictionary) -> void:
	var hit_reward : int = payload.get("hit_reward", 0)
	var hit_mult   : int = payload.get("hit_mult", 0)

	# 個人命中：金色大字彈跳
	_spawn_my_hit_effect(hit_reward, hit_mult)

	# 高倍率（≥120x）額外閃光
	if hit_mult >= 120:
		_flash_gold(0.4)
	elif hit_mult >= 100:
		_flash_green(0.3)

func _on_turtle_leave(payload: Dictionary) -> void:
	_is_active = false
	_total_hits   = payload.get("total_hits", _total_hits)
	_total_reward = payload.get("total_reward", _total_reward)

	# 隱藏橫幅和計數器
	_banner_container.visible = false
	_hit_counter.visible = false

	# 顯示結果彈窗（右側滑入）
	_result_panel.visible = true
	_result_label.text = "🐢 龍龜不死 Boss 離開\n\n全服命中：%d 次\n總獎勵：%d 金幣\n\n感謝所有玩家的參與！" % [_total_hits, _total_reward]

	_result_panel.offset_right = 300  # 從右側外面開始
	_result_panel.offset_left = 20
	var tween := create_tween()
	tween.tween_property(_result_panel, "offset_right", 0.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(_result_panel, "offset_left", -280.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

	# 金色閃光（結算感）
	_flash_gold(0.6)

	# 5 秒後淡出
	await get_tree().create_timer(5.0).timeout
	var fade_tween := create_tween()
	fade_tween.tween_property(_result_panel, "modulate:a", 0.0, 0.5)
	await fade_tween.finished
	hide()
	_result_panel.modulate.a = 1.0

# ---- 視覺效果 ----

## _flash_green 綠色閃光（龜殼感）
func _flash_green(intensity: float) -> void:
	_flash_overlay.color = Color(0.2, 0.9, 0.3, intensity * 0.35)
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, 0.25)

## _flash_gold 金色閃光（獎勵感）
func _flash_gold(intensity: float) -> void:
	_flash_overlay.color = Color(1.0, 0.85, 0.0, intensity * 0.4)
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, 0.3)

## _spawn_hit_float_text 生成命中浮動文字（全服廣播用）
func _spawn_hit_float_text(hitter_name: String, hit_mult: int) -> void:
	var label := Label.new()
	label.text = "%s ×%d" % [hitter_name, hit_mult]
	label.add_theme_color_override("font_color", PANEL_COLOR_TEAL)
	label.add_theme_font_size_override("font_size", 14)
	label.set_anchors_preset(Control.PRESET_CENTER)
	# 隨機位置（避免重疊）
	label.offset_left = randf_range(-200, 200)
	label.offset_top = randf_range(-60, 60)
	add_child(label)

	var tween := create_tween()
	tween.tween_property(label, "offset_top", label.offset_top - 40, 1.2)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 1.2)
	await tween.finished
	if is_instance_valid(label):
		label.queue_free()

## _spawn_my_hit_effect 個人命中效果（金色大字彈跳）
func _spawn_my_hit_effect(hit_reward: int, hit_mult: int) -> void:
	var label := Label.new()
	label.text = "+%d 金幣 (×%d)" % [hit_reward, hit_mult]
	label.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	label.add_theme_font_size_override("font_size", 22)
	label.set_anchors_preset(Control.PRESET_CENTER)
	label.offset_left = -120
	label.offset_top = -20
	add_child(label)

	# 彈跳動畫
	var tween := create_tween()
	tween.tween_property(label, "scale", Vector2(1.3, 1.3), 0.1).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.tween_property(label, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_property(label, "offset_top", label.offset_top - 60, 1.0)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 1.0)
	await tween.finished
	if is_instance_valid(label):
		label.queue_free()
