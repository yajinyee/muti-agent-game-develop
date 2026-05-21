## VampirePanel.gd
## 吸血鬼成長倍率面板（DAY-152）
## 業界依據：jiligames.com 2026「The explicit multiplier of vampires increases the more you fight,
## and there is a chance that you can enter the multiplier mode, up to X5.」
## 設計：深紅血月主題；vampire_grow 顯示倍率成長動畫；
## vampire_blood_moon 全螢幕血紅閃光+橫幅；vampire_killed 右側滑入結果彈窗

extends Control

# ---- 節點 ----
var _blood_moon_banner: Control = null
var _blood_moon_label: Label = null
var _result_panel: Control = null
var _result_title: Label = null
var _result_stats: Label = null
var _result_reward: Label = null
var _flash_overlay: ColorRect = null

# ---- 浮動倍率標籤池（顯示成長動畫）----
var _float_labels: Array = []
const MAX_FLOAT_LABELS = 5

# ---- 顏色（深紅血月主題）----
const COLOR_BLOOD   = Color(0.8, 0.0, 0.0, 1.0)
const COLOR_DARK    = Color(0.5, 0.0, 0.0, 1.0)
const COLOR_CRIMSON = Color(1.0, 0.1, 0.1, 1.0)
const COLOR_BG      = Color(0.06, 0.0, 0.0, 0.95)

# ---- 狀態 ----
var _blood_moon_active: bool = false

func _ready() -> void:
	_build_ui()
	visible = false
	mouse_filter = Control.MOUSE_FILTER_IGNORE

func _build_ui() -> void:
	# 全螢幕閃光層
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(0.8, 0.0, 0.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 血月橫幅（初始隱藏）
	_blood_moon_banner = Control.new()
	_blood_moon_banner.position = Vector2(0, -80)
	_blood_moon_banner.size = Vector2(1280, 72)
	_blood_moon_banner.visible = false
	add_child(_blood_moon_banner)

	var banner_bg = ColorRect.new()
	banner_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	banner_bg.color = Color(0.12, 0.0, 0.0, 0.95)
	_blood_moon_banner.add_child(banner_bg)

	var banner_border = ColorRect.new()
	banner_border.color = COLOR_BLOOD
	banner_border.position = Vector2(0, 68)
	banner_border.size = Vector2(1280, 4)
	_blood_moon_banner.add_child(banner_border)

	_blood_moon_label = Label.new()
	_blood_moon_label.text = "🩸 吸血鬼進入血月模式！×5 倍率！"
	_blood_moon_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	_blood_moon_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_blood_moon_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_blood_moon_label.add_theme_color_override("font_color", COLOR_CRIMSON)
	_blood_moon_label.add_theme_font_size_override("font_size", 26)
	_blood_moon_banner.add_child(_blood_moon_label)

	# 右側結果彈窗
	_result_panel = Control.new()
	_result_panel.position = Vector2(1280, 120)
	_result_panel.size = Vector2(320, 240)
	_result_panel.visible = false
	add_child(_result_panel)

	var result_bg = ColorRect.new()
	result_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	result_bg.color = COLOR_BG
	_result_panel.add_child(result_bg)

	var result_border = ColorRect.new()
	result_border.set_anchors_preset(Control.PRESET_FULL_RECT)
	result_border.color = COLOR_BLOOD
	result_border.custom_minimum_size = Vector2(320, 240)
	_result_panel.add_child(result_border)

	var result_inner = ColorRect.new()
	result_inner.color = COLOR_BG
	result_inner.position = Vector2(2, 2)
	result_inner.size = Vector2(316, 236)
	_result_panel.add_child(result_inner)

	_result_title = Label.new()
	_result_title.text = "🦇 吸血鬼擊破！"
	_result_title.position = Vector2(0, 10)
	_result_title.size = Vector2(320, 36)
	_result_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_title.add_theme_color_override("font_color", COLOR_CRIMSON)
	_result_title.add_theme_font_size_override("font_size", 18)
	_result_panel.add_child(_result_title)

	_result_stats = Label.new()
	_result_stats.text = ""
	_result_stats.position = Vector2(12, 50)
	_result_stats.size = Vector2(296, 130)
	_result_stats.add_theme_color_override("font_color", Color(0.9, 0.9, 0.9, 1.0))
	_result_stats.add_theme_font_size_override("font_size", 14)
	_result_stats.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	_result_panel.add_child(_result_stats)

	_result_reward = Label.new()
	_result_reward.text = "+0 金幣"
	_result_reward.position = Vector2(0, 190)
	_result_reward.size = Vector2(320, 40)
	_result_reward.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_reward.add_theme_color_override("font_color", Color(0.2, 1.0, 0.4, 1.0))
	_result_reward.add_theme_font_size_override("font_size", 22)
	_result_panel.add_child(_result_reward)

# ---- 公開 API ----

## handle_grow 處理吸血鬼倍率成長（收到 vampire_grow 時呼叫）
func handle_grow(data: Dictionary) -> void:
	var mult_bonus: float = data.get("mult_bonus", 1.0)
	var phase_name: String = data.get("phase_name", "")
	var phase_changed: bool = data.get("phase_changed", false)
	var new_mult: float = data.get("new_mult", 1.0)
	var hit_count: int = data.get("hit_count", 0)

	# 顯示浮動倍率文字
	_spawn_float_label("🦇 ×%.0f" % mult_bonus, _get_phase_color(mult_bonus))

	# 階段變化時小閃光
	if phase_changed:
		visible = true
		_flash_overlay.color = Color(0.6, 0.0, 0.0, 0.3)
		var tween = create_tween()
		tween.tween_property(_flash_overlay, "color", Color(0.6, 0.0, 0.0, 0.0), 0.2)

## handle_blood_moon 處理血月模式觸發（收到 vampire_blood_moon 時呼叫）
func handle_blood_moon(data: Dictionary) -> void:
	_blood_moon_active = true
	var mult_bonus: float = data.get("mult_bonus", 5.0)

	_blood_moon_label.text = "🩸 吸血鬼進入血月模式！×%.0f 倍率！誰能擊破它？" % mult_bonus

	# 顯示血月橫幅
	visible = true
	_blood_moon_banner.visible = true
	_blood_moon_banner.position = Vector2(0, -80)
	var tween = create_tween()
	tween.tween_property(_blood_moon_banner, "position", Vector2(0, 0), 0.3).set_ease(Tween.EASE_OUT)

	# 全螢幕血紅閃光
	_flash_overlay.color = Color(0.8, 0.0, 0.0, 0.6)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color", Color(0.8, 0.0, 0.0, 0.0), 0.5)

	# 5 秒後橫幅滑出（但保持 visible=true 直到吸血鬼被擊破）
	var timer = get_tree().create_timer(5.0)
	timer.timeout.connect(func():
		var out_tween = create_tween()
		out_tween.tween_property(_blood_moon_banner, "position", Vector2(0, -80), 0.3).set_ease(Tween.EASE_IN)
	)

## handle_killed 處理吸血鬼被擊破（收到 vampire_killed 時呼叫）
func handle_killed(data: Dictionary) -> void:
	_blood_moon_active = false
	var killer_name: String = data.get("killer_name", "")
	var hit_count: int = data.get("hit_count", 0)
	var mult_bonus: float = data.get("mult_bonus", 1.0)
	var final_mult: float = data.get("final_mult", 1.0)
	var final_reward: int = data.get("final_reward", 0)
	var phase_name: String = data.get("phase_name", "沉睡")

	# 隱藏血月橫幅
	if _blood_moon_banner.visible:
		var tween = create_tween()
		tween.tween_property(_blood_moon_banner, "position", Vector2(0, -80), 0.3).set_ease(Tween.EASE_IN)
		tween.tween_callback(func(): _blood_moon_banner.visible = false)

	# 建立統計文字
	var stats_text = ""
	stats_text += "🦇 擊破者：%s\n" % killer_name
	stats_text += "💥 命中次數：%d 次\n" % hit_count
	stats_text += "🩸 最終階段：%s\n" % phase_name
	stats_text += "📈 倍率加成：×%.0f\n" % mult_bonus
	stats_text += "🎯 最終倍率：×%.1f\n" % final_mult

	_result_stats.text = stats_text
	_result_reward.text = "+%d 金幣" % final_reward

	# 血月模式擊破：血紅色獎勵文字
	if mult_bonus >= 5.0:
		_result_reward.add_theme_color_override("font_color", COLOR_CRIMSON)
		_do_blood_flash()
	elif mult_bonus >= 3.5:
		_result_reward.add_theme_color_override("font_color", COLOR_BLOOD)
	else:
		_result_reward.add_theme_color_override("font_color", Color(0.2, 1.0, 0.4, 1.0))

	# 顯示結果面板（從右側滑入）
	visible = true
	_result_panel.visible = true
	_result_panel.position = Vector2(1280, 120)
	var tween = create_tween()
	tween.tween_property(_result_panel, "position", Vector2(950, 120), 0.35).set_ease(Tween.EASE_OUT)

	# 3.5 秒後關閉
	var close_timer = get_tree().create_timer(3.5)
	close_timer.timeout.connect(_close_panel)

func _close_panel() -> void:
	var tween = create_tween()
	tween.tween_property(_result_panel, "position", Vector2(1280, 120), 0.3).set_ease(Tween.EASE_IN)
	tween.tween_callback(func():
		_result_panel.visible = false
		if not _blood_moon_active:
			visible = false
	)

func _spawn_float_label(text: String, color: Color) -> void:
	# 建立浮動文字（在畫面中央偏上）
	var label = Label.new()
	label.text = text
	label.position = Vector2(540 + randf_range(-60, 60), 300 + randf_range(-30, 30))
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", 20)
	add_child(label)
	visible = true

	var tween = create_tween()
	tween.tween_property(label, "position", label.position + Vector2(0, -60), 0.8)
	tween.parallel().tween_property(label, "modulate", Color(1, 1, 1, 0), 0.8)
	tween.tween_callback(label.queue_free)

func _get_phase_color(mult_bonus: float) -> Color:
	if mult_bonus >= 5.0:
		return COLOR_CRIMSON
	elif mult_bonus >= 3.5:
		return COLOR_BLOOD
	elif mult_bonus >= 2.0:
		return COLOR_DARK
	else:
		return Color(0.7, 0.7, 0.7, 1.0)

func _do_blood_flash() -> void:
	_flash_overlay.color = Color(0.8, 0.0, 0.0, 0.6)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color", Color(0.8, 0.0, 0.0, 0.0), 0.25)
	tween.tween_property(_flash_overlay, "color", Color(0.6, 0.0, 0.0, 0.4), 0.15)
	tween.tween_property(_flash_overlay, "color", Color(0.6, 0.0, 0.0, 0.0), 0.3)
