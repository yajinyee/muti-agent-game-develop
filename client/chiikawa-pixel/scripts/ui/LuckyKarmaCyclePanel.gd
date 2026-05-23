## LuckyKarmaCyclePanel.gd — 幸運命運輪迴魚 UI（DAY-264）
## 紫金命運主題面板
## 主色：#9B59B6 紫 + #FFD700 金 + #E8DAEF 淡紫白
extends CanvasLayer

const KARMA_COLOR_LOW  = Color(0.61, 0.35, 0.71, 1.0)  # #9B59B6 紫（業力 1-4）
const KARMA_COLOR_MID  = Color(1.0, 0.65, 0.0, 1.0)    # #FFA500 橙（業力 5-8）
const KARMA_COLOR_HIGH = Color(1.0, 0.85, 0.0, 1.0)    # #FFD700 金（業力 9-10）

var _flash_overlay: ColorRect
var _banner: PanelContainer
var _banner_label: Label
var _karma_counter: Label
var _karma_bar: ColorRect
var _karma_bar_bg: ColorRect
var _timer_bar: ColorRect
var _timer_bar_bg: ColorRect
var _timer_duration: float = 20.0
var _timer_elapsed: float = 0.0
var _timer_active: bool = false
var _current_karma: int = 0
var _max_karma: int = 10

func _ready() -> void:
	layer = 37
	_build_ui()

func _build_ui() -> void:
	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(0.61, 0.35, 0.71, 0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner = PanelContainer.new()
	_banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 56)
	_banner.modulate.a = 0.0
	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.08, 0.04, 0.12, 0.88)
	banner_style.border_color = Color(0.61, 0.35, 0.71, 1.0)
	banner_style.set_border_width_all(2)
	banner_style.set_corner_radius_all(6)
	_banner.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_font_size_override("font_size", 18)
	_banner_label.add_theme_color_override("font_color", Color(0.91, 0.78, 1.0))
	_banner.add_child(_banner_label)

	# 業力計數器（右上角）
	var vp_size = get_viewport().size
	_karma_counter = Label.new()
	_karma_counter.text = "☯️ 業力 0/10"
	_karma_counter.position = Vector2(vp_size.x - 180, 65)
	_karma_counter.add_theme_font_size_override("font_size", 18)
	_karma_counter.add_theme_color_override("font_color", KARMA_COLOR_LOW)
	_karma_counter.modulate.a = 0.0
	add_child(_karma_counter)

	# 業力進度條（底部）
	_karma_bar_bg = ColorRect.new()
	_karma_bar_bg.color = Color(0.1, 0.05, 0.15, 0.8)
	_karma_bar_bg.size = Vector2(vp_size.x, 12)
	_karma_bar_bg.position = Vector2(0, vp_size.y - 60)
	_karma_bar_bg.modulate.a = 0.0
	add_child(_karma_bar_bg)

	_karma_bar = ColorRect.new()
	_karma_bar.color = KARMA_COLOR_LOW
	_karma_bar.size = Vector2(0, 12)
	_karma_bar.position = Vector2(0, vp_size.y - 60)
	_karma_bar.modulate.a = 0.0
	add_child(_karma_bar)

	# 計時條（右側豎向，x=-254 與其他計時條錯開）
	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.1, 0.05, 0.15, 0.7)
	_timer_bar_bg.size = Vector2(8, 200)
	_timer_bar_bg.position = Vector2(vp_size.x - 254, vp_size.y / 2 - 100)
	_timer_bar_bg.modulate.a = 0.0
	add_child(_timer_bar_bg)

	_timer_bar = ColorRect.new()
	_timer_bar.color = KARMA_COLOR_LOW
	_timer_bar.size = Vector2(8, 200)
	_timer_bar.position = Vector2(vp_size.x - 254, vp_size.y / 2 - 100)
	_timer_bar.modulate.a = 0.0
	add_child(_timer_bar)

func _process(delta: float) -> void:
	if not _timer_active:
		return
	_timer_elapsed += delta
	var pct = 1.0 - (_timer_elapsed / _timer_duration)
	pct = clamp(pct, 0.0, 1.0)
	_timer_bar.size.y = 200.0 * pct

	# 計時條顏色隨業力變化
	var karma_pct = float(_current_karma) / float(_max_karma)
	if karma_pct >= 0.9:
		_timer_bar.color = KARMA_COLOR_HIGH
	elif karma_pct >= 0.5:
		_timer_bar.color = KARMA_COLOR_MID
	else:
		_timer_bar.color = KARMA_COLOR_LOW

	if _timer_elapsed >= _timer_duration:
		_timer_active = false

## 處理來自 GameManager 的命運輪迴事件
func handle_event(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"karma_start":
			_on_karma_start(payload)
		"karma_broadcast":
			_on_karma_broadcast(payload)
		"karma_update":
			_on_karma_update(payload)
		"karma_burst":
			_on_karma_burst(payload)
		"karma_burst_broadcast":
			_on_karma_burst_broadcast(payload)
		"karma_settle":
			_on_karma_settle(payload)
		"karma_expire":
			_on_karma_expire()

## 觸發命運輪迴（個人）
func _on_karma_start(payload: Dictionary) -> void:
	_timer_duration = float(payload.get("duration", 20))
	_max_karma = payload.get("max_karma", 10)
	_current_karma = 0
	_timer_elapsed = 0.0
	_timer_active = true

	# 紫色三次強閃光
	_flash_triple(Color(0.61, 0.35, 0.71, 0.5))

	# 顯示橫幅
	_banner_label.text = "☯️ 命運輪迴！累積 10 業力獲得 ×15.0 大獎！"
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)
	tween.tween_interval(_timer_duration)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.5)

	# 顯示業力計數器和進度條
	_karma_counter.modulate.a = 1.0
	_karma_bar_bg.modulate.a = 1.0
	_karma_bar.modulate.a = 1.0
	_timer_bar_bg.modulate.a = 1.0
	_timer_bar.modulate.a = 1.0
	_update_karma_display()

	# 大字提示
	_show_big_text("☯️ 命運輪迴！", KARMA_COLOR_LOW)

## 全服廣播
func _on_karma_broadcast(payload: Dictionary) -> void:
	var trigger = payload.get("trigger_name", "玩家")
	_banner_label.text = "☯️ %s 觸發命運輪迴！" % trigger
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)
	tween.tween_interval(3.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.3)

## 業力更新（個人）
func _on_karma_update(payload: Dictionary) -> void:
	_current_karma = payload.get("karma", 0)
	_update_karma_display()

	# 業力顏色閃光
	var karma_color = _get_karma_color()
	_flash_overlay.color = Color(karma_color.r, karma_color.g, karma_color.b, 0.2)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, 0.15)

	# 業力接近滿時脈衝動畫
	if _current_karma >= 8:
		var pulse = create_tween().set_loops(2)
		pulse.tween_property(_karma_counter, "scale", Vector2(1.2, 1.2), 0.1)
		pulse.tween_property(_karma_counter, "scale", Vector2(1.0, 1.0), 0.1)

## 命運爆發（個人）
func _on_karma_burst(payload: Dictionary) -> void:
	var karma = payload.get("karma", 10)
	var mult = payload.get("mult", 15.0)
	var reward = payload.get("reward", 0)

	_timer_active = false

	# 全螢幕三次強閃光（金色）
	_flash_triple(Color(1.0, 0.85, 0.0, 0.7))

	# 大字顯示
	_show_big_text("☯️ 命運爆發！×%.1f" % mult, KARMA_COLOR_HIGH)

	# 結算彈窗
	_show_result_popup("☯️ 命運爆發！", karma, mult, reward, KARMA_COLOR_HIGH)

	# 隱藏計數器和進度條
	_hide_karma_ui()

## 命運爆發全服廣播
func _on_karma_burst_broadcast(payload: Dictionary) -> void:
	var player_name = payload.get("player_name", "玩家")
	var mult = payload.get("mult", 15.0)
	var reward = payload.get("reward", 0)

	_banner_label.text = "☯️ %s 業力滿溢！命運爆發！×%.1f +%d！" % [player_name, mult, reward]
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)
	tween.tween_interval(4.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.3)

## 業力結算（個人）
func _on_karma_settle(payload: Dictionary) -> void:
	var karma = payload.get("karma", 0)
	var mult = payload.get("mult", 0.0)
	var reward = payload.get("reward", 0)

	_timer_active = false
	_show_big_text("☯️ 業力結算 %d/10 ×%.1f" % [karma, mult], KARMA_COLOR_MID)
	_show_result_popup("☯️ 業力結算", karma, mult, reward, KARMA_COLOR_MID)
	_hide_karma_ui()

## 業力消散（個人）
func _on_karma_expire() -> void:
	_timer_active = false
	_hide_karma_ui()

## 工具函數：更新業力顯示
func _update_karma_display() -> void:
	var karma_color = _get_karma_color()
	_karma_counter.text = "☯️ 業力 %d/%d" % [_current_karma, _max_karma]
	_karma_counter.add_theme_color_override("font_color", karma_color)

	# 更新進度條
	var vp_size = get_viewport().size
	var pct = float(_current_karma) / float(_max_karma)
	_karma_bar.size.x = vp_size.x * pct
	_karma_bar.color = karma_color

## 工具函數：取得業力顏色
func _get_karma_color() -> Color:
	var pct = float(_current_karma) / float(_max_karma)
	if pct >= 0.9:
		return KARMA_COLOR_HIGH
	elif pct >= 0.5:
		return KARMA_COLOR_MID
	return KARMA_COLOR_LOW

## 工具函數：隱藏業力 UI
func _hide_karma_ui() -> void:
	var tween = create_tween()
	tween.tween_property(_karma_counter, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_karma_bar_bg, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_karma_bar, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_timer_bar_bg, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_timer_bar, "modulate:a", 0.0, 0.5)

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
	label.add_theme_font_size_override("font_size", 30)
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

## 工具函數：顯示結算彈窗
func _show_result_popup(title: String, karma: int, mult: float, reward: int, color: Color) -> void:
	var popup = PanelContainer.new()
	popup.custom_minimum_size = Vector2(280, 130)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.02, 0.1, 0.92)
	style.border_color = color
	style.set_border_width_all(2)
	style.set_corner_radius_all(8)
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	popup.add_child(vbox)

	var title_label = Label.new()
	title_label.text = title
	title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_label.add_theme_font_size_override("font_size", 18)
	title_label.add_theme_color_override("font_color", color)
	vbox.add_child(title_label)

	var karma_label = Label.new()
	karma_label.text = "業力：%d/%d" % [karma, _max_karma]
	karma_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	karma_label.add_theme_font_size_override("font_size", 15)
	karma_label.add_theme_color_override("font_color", Color(0.91, 0.78, 1.0))
	vbox.add_child(karma_label)

	var mult_label = Label.new()
	mult_label.text = "倍率：×%.1f" % mult
	mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mult_label.add_theme_font_size_override("font_size", 16)
	mult_label.add_theme_color_override("font_color", KARMA_COLOR_HIGH)
	vbox.add_child(mult_label)

	var reward_label = Label.new()
	reward_label.text = "獎勵：+%d 🪙" % reward
	reward_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_label.add_theme_font_size_override("font_size", 16)
	reward_label.add_theme_color_override("font_color", Color(0.0, 1.0, 0.53))
	vbox.add_child(reward_label)

	var vp_size = get_viewport().size
	popup.position = Vector2(vp_size.x, vp_size.y / 2 - 65)
	popup.modulate.a = 0.0
	add_child(popup)

	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 300, 0.3)
	tween.parallel().tween_property(popup, "modulate:a", 1.0, 0.3)
	tween.tween_interval(3.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.4)
	tween.tween_callback(popup.queue_free)
