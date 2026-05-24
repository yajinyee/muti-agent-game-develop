## LuckySpinWheelPanel.gd — 幸運輪盤魚 UI（DAY-269）
## 粉金輪盤主題面板（個人幸運輪盤）
## 主色：#FF69B4 粉紅 + #FFD700 金 + #FF4500 火橙 + #1A1A2E 深藍黑
extends CanvasLayer

const WHEEL_COLOR_PINK   = Color(1.0, 0.41, 0.71, 1.0)  # #FF69B4 粉紅
const WHEEL_COLOR_GOLD   = Color(1.0, 0.85, 0.0, 1.0)   # #FFD700 金
const WHEEL_COLOR_FIRE   = Color(1.0, 0.27, 0.0, 1.0)   # #FF4500 火橙
const WHEEL_COLOR_GREEN  = Color(0.0, 1.0, 0.53, 1.0)   # #00FF88 翠綠
const WHEEL_COLOR_RED    = Color(1.0, 0.2, 0.2, 1.0)    # 懲罰紅

# 扇區顏色對應（×0.5/×2.0/×3.0/×5.0/×8.0）
const SECTOR_COLORS = [
	Color(1.0, 0.2, 0.2, 1.0),   # ×0.5 紅（懲罰）
	Color(0.0, 0.75, 1.0, 1.0),  # ×2.0 天藍
	Color(0.0, 1.0, 0.53, 1.0),  # ×3.0 翠綠
	Color(1.0, 0.85, 0.0, 1.0),  # ×5.0 金
	Color(1.0, 0.27, 0.0, 1.0),  # ×8.0 火橙
]

var _flash_overlay: ColorRect
var _banner: PanelContainer
var _banner_label: Label
var _mult_indicator: Label    # 右上角倍率指示器
var _timer_bar: ColorRect     # 右側計時條
var _timer_bar_bg: ColorRect  # 計時條背景
var _session_active: bool = false
var _current_mult: float = 1.0
var _duration: float = 8.0
var _start_time: float = 0.0

func _ready() -> void:
	layer = 42
	_build_ui()

func _build_ui() -> void:
	var vp_size = get_viewport().size

	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 0.41, 0.71, 0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner = PanelContainer.new()
	_banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 56)
	_banner.modulate.a = 0.0
	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.1, 0.02, 0.06, 0.88)
	banner_style.border_color = WHEEL_COLOR_PINK
	banner_style.set_border_width_all(2)
	banner_style.set_corner_radius_all(6)
	_banner.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_font_size_override("font_size", 18)
	_banner_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.95))
	_banner.add_child(_banner_label)

	# 倍率指示器（右上角）
	_mult_indicator = Label.new()
	_mult_indicator.text = "🎡 ×?"
	_mult_indicator.position = Vector2(vp_size.x - 150, 65)
	_mult_indicator.add_theme_font_size_override("font_size", 22)
	_mult_indicator.add_theme_color_override("font_color", WHEEL_COLOR_PINK)
	_mult_indicator.modulate.a = 0.0
	add_child(_mult_indicator)

	# 右側計時條背景（x=-296，與其他計時條錯開）
	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.1, 0.02, 0.06, 0.5)
	_timer_bar_bg.size = Vector2(8, 200)
	_timer_bar_bg.position = Vector2(vp_size.x - 296, vp_size.y / 2 - 100)
	_timer_bar_bg.modulate.a = 0.0
	add_child(_timer_bar_bg)

	# 右側計時條
	_timer_bar = ColorRect.new()
	_timer_bar.color = WHEEL_COLOR_PINK
	_timer_bar.size = Vector2(8, 200)
	_timer_bar.position = Vector2(vp_size.x - 296, vp_size.y / 2 - 100)
	_timer_bar.modulate.a = 0.0
	add_child(_timer_bar)

func _process(_delta: float) -> void:
	if not _session_active:
		return
	var elapsed = Time.get_ticks_msec() / 1000.0 - _start_time
	var remaining = max(0.0, _duration - elapsed)
	var ratio = remaining / _duration
	_timer_bar.size.y = 200.0 * ratio
	_timer_bar.position.y = (get_viewport().size.y / 2 - 100) + 200.0 * (1.0 - ratio)
	if ratio < 0.3:
		_timer_bar.color = WHEEL_COLOR_GOLD
	else:
		_timer_bar.color = _get_mult_color(_current_mult)

## 處理來自 GameManager 的輪盤事件
func handle_event(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"spin_start":
			_on_spin_start(payload)
		"spin_broadcast":
			_on_spin_broadcast(payload)
		"spin_result":
			_on_spin_result(payload)
		"spin_boost_kill":
			_on_spin_boost_kill(payload)
		"spin_expire":
			_on_spin_expire(payload)

## 輪盤開始旋轉（個人）
func _on_spin_start(payload: Dictionary) -> void:
	var player_name = payload.get("player_name", "玩家")

	# 粉紅三次強閃光
	_flash_triple(Color(1.0, 0.41, 0.71, 0.4))

	# 顯示橫幅（旋轉中）
	_banner_label.text = "🎡 %s 觸發幸運輪盤！轉動中..." % player_name
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)

	# 顯示倍率指示器（旋轉動畫）
	_mult_indicator.text = "🎡 旋轉中..."
	_mult_indicator.add_theme_color_override("font_color", WHEEL_COLOR_PINK)
	var counter_tween = create_tween()
	counter_tween.tween_property(_mult_indicator, "modulate:a", 1.0, 0.3)

	# 大字提示
	_show_big_text("🎡 幸運輪盤！", WHEEL_COLOR_PINK)

## 全服廣播
func _on_spin_broadcast(payload: Dictionary) -> void:
	var player_name = payload.get("player_name", "玩家")
	var result_mult = payload.get("result_mult", 1.0)
	_banner_label.text = "🎡 %s 幸運輪盤結果：×%.1f！" % [player_name, result_mult]
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.1)
	tween.tween_interval(3.5)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.3)

## 輪盤結果（個人）
func _on_spin_result(payload: Dictionary) -> void:
	var result_mult = payload.get("result_mult", 1.0)
	var sector_index = payload.get("sector_index", 0)
	_duration = payload.get("duration", 8.0)
	_current_mult = result_mult
	_start_time = Time.get_ticks_msec() / 1000.0
	_session_active = true

	var mult_color = _get_mult_color(result_mult)

	# 更新倍率指示器
	_mult_indicator.text = "🎡 ×%.1f" % result_mult
	_mult_indicator.add_theme_color_override("font_color", mult_color)

	# 計數器脈衝
	var pulse = create_tween()
	pulse.tween_property(_mult_indicator, "scale", Vector2(1.4, 1.4), 0.1)
	pulse.tween_property(_mult_indicator, "scale", Vector2(1.0, 1.0), 0.15)

	# 顯示計時條
	var timer_tween = create_tween()
	timer_tween.tween_property(_timer_bar_bg, "modulate:a", 1.0, 0.3)
	timer_tween.parallel().tween_property(_timer_bar, "modulate:a", 1.0, 0.3)
	_timer_bar.color = mult_color

	# 根據倍率決定閃光強度
	if result_mult >= 8.0:
		_flash_triple(Color(1.0, 0.27, 0.0, 0.7))
		_show_big_text("🎡 ×8.0 大獎！8 秒黃金時間！", WHEEL_COLOR_FIRE)
	elif result_mult >= 5.0:
		_flash_triple(Color(1.0, 0.85, 0.0, 0.6))
		_show_big_text("🎡 ×%.1f！8 秒加成！" % result_mult, WHEEL_COLOR_GOLD)
	elif result_mult <= 0.5:
		_flash_triple(Color(1.0, 0.2, 0.2, 0.5))
		_show_big_text("🎡 ×0.5 懲罰！小心！", WHEEL_COLOR_RED)
	else:
		_flash_triple(Color(1.0, 0.41, 0.71, 0.4))
		_show_big_text("🎡 ×%.1f！8 秒加成！" % result_mult, mult_color)

	# 更新橫幅
	_banner_label.text = "🎡 輪盤結果：×%.1f！8 秒加成！" % result_mult
	var banner_tween = create_tween()
	banner_tween.tween_interval(4.0)
	banner_tween.tween_property(_banner, "modulate:a", 0.0, 0.5)

	# (unused variable suppressed)

## 加成期間擊破（個人）
func _on_spin_boost_kill(payload: Dictionary) -> void:
	var target_name = payload.get("target_name", "")
	var reward = payload.get("reward", 0)
	var current_mult = payload.get("current_mult", 1.0)
	var mult_color = _get_mult_color(current_mult)

	# 輕微閃光
	_flash_overlay.color = Color(mult_color.r, mult_color.g, mult_color.b, 0)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color:a", 0.1, 0.04)
	flash_tween.tween_property(_flash_overlay, "color:a", 0.0, 0.08)

	# 浮動文字
	if target_name != "":
		_show_float_text("🎡 %s ×%.1f +%d" % [target_name, current_mult, reward], mult_color)

## 加成結束（個人）
func _on_spin_expire(payload: Dictionary) -> void:
	var total_reward = payload.get("total_reward", 0)
	var kill_count = payload.get("kill_count", 0)
	var current_mult = payload.get("current_mult", 1.0)

	_session_active = false

	# 結算提示
	_show_big_text("🎡 輪盤加成結束！×%.1f 擊破 %d 個 +%d" % [current_mult, kill_count, total_reward],
		_get_mult_color(current_mult))

	# 隱藏 UI
	var tween = create_tween()
	tween.tween_property(_mult_indicator, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_timer_bar, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_timer_bar_bg, "modulate:a", 0.0, 0.5)

## 工具函數：根據倍率取得顏色
func _get_mult_color(mult: float) -> Color:
	if mult <= 0.5:
		return WHEEL_COLOR_RED
	elif mult <= 2.0:
		return Color(0.0, 0.75, 1.0, 1.0)  # 天藍
	elif mult <= 3.0:
		return WHEEL_COLOR_GREEN
	elif mult <= 5.0:
		return WHEEL_COLOR_GOLD
	else:
		return WHEEL_COLOR_FIRE

## 工具函數：三次強閃光
func _flash_triple(color: Color) -> void:
	var tween = create_tween()
	for i in range(3):
		tween.tween_property(_flash_overlay, "color:a", color.a, 0.06)
		tween.tween_property(_flash_overlay, "color:a", 0.0, 0.1)
	_flash_overlay.color = Color(color.r, color.g, color.b, 0)

## 工具函數：顯示大字
func _show_big_text(text: String, color: Color) -> void:
	var label = Label.new()
	label.text = text
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.add_theme_font_size_override("font_size", 28)
	label.add_theme_color_override("font_color", color)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.position.y -= 80
	label.modulate.a = 0.0
	add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "modulate:a", 1.0, 0.15)
	tween.tween_property(label, "position:y", label.position.y - 40, 1.2)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.4).set_delay(0.8)
	tween.tween_callback(label.queue_free)

## 工具函數：浮動文字
func _show_float_text(text: String, color: Color) -> void:
	var label = Label.new()
	label.text = text
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.add_theme_font_size_override("font_size", 18)
	label.add_theme_color_override("font_color", color)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.position.y -= 40
	label.modulate.a = 0.0
	add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "modulate:a", 1.0, 0.1)
	tween.tween_property(label, "position:y", label.position.y - 25, 0.7)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.3).set_delay(0.4)
	tween.tween_callback(label.queue_free)
