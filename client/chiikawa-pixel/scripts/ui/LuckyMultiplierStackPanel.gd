## LuckyMultiplierStackPanel.gd — 幸運倍率疊加魚 UI（DAY-267）
## 翠綠疊加主題面板（Fishing Fortune Multiplier Cascade 風格）
## 主色：#00FF88 翠綠 + #FFD700 金 + #00BFFF 天藍 + #1A1A2E 深藍黑
extends CanvasLayer

const STACK_COLOR_GREEN  = Color(0.0, 1.0, 0.53, 1.0)   # #00FF88 翠綠
const STACK_COLOR_GOLD   = Color(1.0, 0.85, 0.0, 1.0)   # #FFD700 金
const STACK_COLOR_BLUE   = Color(0.0, 0.75, 1.0, 1.0)   # #00BFFF 天藍
const STACK_COLOR_BURST  = Color(1.0, 0.5, 0.0, 1.0)    # #FF8000 爆發橙

var _flash_overlay: ColorRect
var _banner: PanelContainer
var _banner_label: Label
var _stack_counter: Label      # 右上角疊加倍率計數器
var _stack_bar: ColorRect      # 底部疊加進度條
var _stack_bar_bg: ColorRect   # 進度條背景
var _timer_bar: ColorRect      # 右側計時條
var _timer_bar_bg: ColorRect   # 計時條背景
var _session_active: bool = false
var _max_stack: float = 10.0
var _duration: float = 25.0
var _start_time: float = 0.0

func _ready() -> void:
	layer = 40
	_build_ui()

func _build_ui() -> void:
	var vp_size = get_viewport().size

	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(0.0, 1.0, 0.53, 0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner = PanelContainer.new()
	_banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 56)
	_banner.modulate.a = 0.0
	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.0, 0.1, 0.05, 0.88)
	banner_style.border_color = STACK_COLOR_GREEN
	banner_style.set_border_width_all(2)
	banner_style.set_corner_radius_all(6)
	_banner.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_font_size_override("font_size", 18)
	_banner_label.add_theme_color_override("font_color", Color(0.8, 1.0, 0.9))
	_banner.add_child(_banner_label)

	# 疊加倍率計數器（右上角）
	_stack_counter = Label.new()
	_stack_counter.text = "📈 ×1.0"
	_stack_counter.position = Vector2(vp_size.x - 170, 65)
	_stack_counter.add_theme_font_size_override("font_size", 22)
	_stack_counter.add_theme_color_override("font_color", STACK_COLOR_GREEN)
	_stack_counter.modulate.a = 0.0
	add_child(_stack_counter)

	# 底部疊加進度條背景
	_stack_bar_bg = ColorRect.new()
	_stack_bar_bg.color = Color(0.0, 0.1, 0.05, 0.7)
	_stack_bar_bg.size = Vector2(vp_size.x, 12)
	_stack_bar_bg.position = Vector2(0, vp_size.y - 12)
	_stack_bar_bg.modulate.a = 0.0
	add_child(_stack_bar_bg)

	# 底部疊加進度條
	_stack_bar = ColorRect.new()
	_stack_bar.color = STACK_COLOR_GREEN
	_stack_bar.size = Vector2(0, 12)
	_stack_bar.position = Vector2(0, vp_size.y - 12)
	_stack_bar.modulate.a = 0.0
	add_child(_stack_bar)

	# 右側計時條背景（x=-282，與其他計時條錯開）
	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.0, 0.1, 0.05, 0.5)
	_timer_bar_bg.size = Vector2(8, 200)
	_timer_bar_bg.position = Vector2(vp_size.x - 282, vp_size.y / 2 - 100)
	_timer_bar_bg.modulate.a = 0.0
	add_child(_timer_bar_bg)

	# 右側計時條
	_timer_bar = ColorRect.new()
	_timer_bar.color = STACK_COLOR_GREEN
	_timer_bar.size = Vector2(8, 200)
	_timer_bar.position = Vector2(vp_size.x - 282, vp_size.y / 2 - 100)
	_timer_bar.modulate.a = 0.0
	add_child(_timer_bar)

func _process(delta: float) -> void:
	if not _session_active:
		return
	var elapsed = Time.get_ticks_msec() / 1000.0 - _start_time
	var remaining = max(0.0, _duration - elapsed)
	var ratio = remaining / _duration
	_timer_bar.size.y = 200.0 * ratio
	_timer_bar.position.y = (get_viewport().size.y / 2 - 100) + 200.0 * (1.0 - ratio)
	# 剩餘時間少時計時條變金色
	if ratio < 0.3:
		_timer_bar.color = STACK_COLOR_GOLD
	else:
		_timer_bar.color = STACK_COLOR_GREEN

## 處理來自 GameManager 的倍率疊加事件
func handle_event(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"stack_start":
			_on_stack_start(payload)
		"stack_broadcast":
			_on_stack_broadcast(payload)
		"stack_update":
			_on_stack_update(payload)
		"stack_burst":
			_on_stack_burst(payload)
		"stack_burst_broadcast":
			_on_stack_burst_broadcast(payload)
		"stack_settle":
			_on_stack_settle(payload)

## 倍率疊加觸發（個人）
func _on_stack_start(payload: Dictionary) -> void:
	_max_stack = payload.get("max_stack", 10.0)
	_duration = payload.get("duration", 25.0)
	_start_time = Time.get_ticks_msec() / 1000.0
	_session_active = true

	# 翠綠三次強閃光
	_flash_triple(Color(0.0, 1.0, 0.53, 0.45))

	# 顯示橫幅
	_banner_label.text = "📈 倍率疊加！每次擊破 +0.3x，最高 ×%.0f！" % _max_stack
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)
	tween.tween_interval(4.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.5)

	# 顯示疊加計數器
	_stack_counter.text = "📈 ×1.0"
	_stack_counter.add_theme_color_override("font_color", STACK_COLOR_GREEN)
	var counter_tween = create_tween()
	counter_tween.tween_property(_stack_counter, "modulate:a", 1.0, 0.3)

	# 顯示進度條
	var bar_tween = create_tween()
	bar_tween.tween_property(_stack_bar_bg, "modulate:a", 1.0, 0.3)
	bar_tween.parallel().tween_property(_stack_bar, "modulate:a", 1.0, 0.3)

	# 顯示計時條
	var timer_tween = create_tween()
	timer_tween.tween_property(_timer_bar_bg, "modulate:a", 1.0, 0.3)
	timer_tween.parallel().tween_property(_timer_bar, "modulate:a", 1.0, 0.3)

	# 大字提示
	_show_big_text("📈 倍率疊加！", STACK_COLOR_GREEN)

## 全服廣播
func _on_stack_broadcast(payload: Dictionary) -> void:
	var trigger = payload.get("player_name", "玩家")
	_banner_label.text = "📈 %s 觸發倍率疊加！每次擊破 +0.3x！" % trigger
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)
	tween.tween_interval(3.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.3)

## 每次擊破疊加更新（個人）
func _on_stack_update(payload: Dictionary) -> void:
	var current_stack = payload.get("current_stack", 1.0)
	var kill_count = payload.get("kill_count", 0)
	var target_name = payload.get("target_name", "")
	var reward = payload.get("reward", 0)

	# 更新疊加計數器
	var stack_ratio = current_stack / _max_stack
	var counter_color: Color
	if stack_ratio >= 0.8:
		counter_color = STACK_COLOR_GOLD
	elif stack_ratio >= 0.5:
		counter_color = STACK_COLOR_BLUE
	else:
		counter_color = STACK_COLOR_GREEN
	_stack_counter.text = "📈 ×%.1f" % current_stack
	_stack_counter.add_theme_color_override("font_color", counter_color)

	# 更新進度條
	var vp_size = get_viewport().size
	var bar_width = vp_size.x * stack_ratio
	var bar_tween = create_tween()
	bar_tween.tween_property(_stack_bar, "size:x", bar_width, 0.15)
	# 接近滿時進度條變金色
	if stack_ratio >= 0.8:
		bar_tween.parallel().tween_property(_stack_bar, "color", STACK_COLOR_GOLD, 0.15)
	else:
		bar_tween.parallel().tween_property(_stack_bar, "color", STACK_COLOR_GREEN, 0.15)

	# 計數器脈衝
	var pulse = create_tween()
	pulse.tween_property(_stack_counter, "scale", Vector2(1.25, 1.25), 0.07)
	pulse.tween_property(_stack_counter, "scale", Vector2(1.0, 1.0), 0.1)

	# 翠綠閃光
	_flash_overlay.color = Color(0.0, 1.0, 0.53, 0)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color:a", 0.12, 0.04)
	flash_tween.tween_property(_flash_overlay, "color:a", 0.0, 0.1)

	# 浮動文字
	if target_name != "":
		_show_float_text("📈 %s ×%.1f +%d" % [target_name, current_stack, reward], counter_color)

	_ = kill_count

## 倍率爆發（達到 10.0x，個人）
func _on_stack_burst(payload: Dictionary) -> void:
	var total_stack = payload.get("total_stack", 10.0)
	var burst_reward = payload.get("burst_reward", 0)
	var total_reward = payload.get("total_reward", 0)
	var player_name = payload.get("player_name", "玩家")

	_session_active = false

	# 全螢幕三次強閃光（金色爆發）
	_flash_triple(Color(1.0, 0.85, 0.0, 0.7))

	# 大字提示
	_show_big_text("📈 倍率爆發！×%.0f → ×%.0f！" % [total_stack, 20.0], STACK_COLOR_GOLD)

	# 結算彈窗
	_show_burst_popup(player_name, total_stack, burst_reward, total_reward)

	# 隱藏計時條和進度條
	_hide_session_ui()

## 爆發全服廣播
func _on_stack_burst_broadcast(payload: Dictionary) -> void:
	var player_name = payload.get("player_name", "玩家")
	var total_stack = payload.get("total_stack", 10.0)
	var burst_reward = payload.get("burst_reward", 0)
	_banner_label.text = "📈 %s 倍率疊加達到 ×%.0f！爆發獎勵 +%d！" % [player_name, total_stack, burst_reward]
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)
	tween.tween_interval(4.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.4)

## 超時結算（個人）
func _on_stack_settle(payload: Dictionary) -> void:
	var final_stack = payload.get("final_stack", 1.0)
	var kill_count = payload.get("kill_count", 0)
	var total_reward = payload.get("total_reward", 0)

	_session_active = false

	# 天藍閃光（結算）
	_flash_overlay.color = Color(0.0, 0.75, 1.0, 0)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color:a", 0.3, 0.1)
	flash_tween.tween_property(_flash_overlay, "color:a", 0.0, 0.2)

	# 大字提示
	_show_big_text("📈 疊加結算 ×%.1f / %d 次 +%d" % [final_stack, kill_count, total_reward], STACK_COLOR_BLUE)

	# 隱藏計時條和進度條
	_hide_session_ui()

## 工具函數：隱藏 session UI
func _hide_session_ui() -> void:
	var tween = create_tween()
	tween.tween_property(_stack_counter, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_stack_bar, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_stack_bar_bg, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_timer_bar, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_timer_bar_bg, "modulate:a", 0.0, 0.5)

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

## 工具函數：浮動文字
func _show_float_text(text: String, color: Color) -> void:
	var label = Label.new()
	label.text = text
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.add_theme_font_size_override("font_size", 20)
	label.add_theme_color_override("font_color", color)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.position.y -= 40
	label.modulate.a = 0.0
	add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "modulate:a", 1.0, 0.1)
	tween.tween_property(label, "position:y", label.position.y - 30, 0.8)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.3).set_delay(0.5)
	tween.tween_callback(label.queue_free)

## 工具函數：顯示爆發結算彈窗
func _show_burst_popup(player_name: String, total_stack: float, burst_reward: int, total_reward: int) -> void:
	var popup = PanelContainer.new()
	popup.custom_minimum_size = Vector2(300, 150)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.08, 0.02, 0.92)
	style.border_color = STACK_COLOR_GOLD
	style.set_border_width_all(2)
	style.set_corner_radius_all(8)
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	popup.add_child(vbox)

	var title_label = Label.new()
	title_label.text = "📈 倍率爆發！"
	title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_label.add_theme_font_size_override("font_size", 20)
	title_label.add_theme_color_override("font_color", STACK_COLOR_GOLD)
	vbox.add_child(title_label)

	var player_label = Label.new()
	player_label.text = "玩家：%s" % player_name
	player_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	player_label.add_theme_font_size_override("font_size", 14)
	player_label.add_theme_color_override("font_color", Color(0.8, 1.0, 0.9))
	vbox.add_child(player_label)

	var stack_label = Label.new()
	stack_label.text = "最終疊加：×%.0f → 爆發 ×20.0" % total_stack
	stack_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	stack_label.add_theme_font_size_override("font_size", 15)
	stack_label.add_theme_color_override("font_color", STACK_COLOR_BURST)
	vbox.add_child(stack_label)

	var burst_label = Label.new()
	burst_label.text = "爆發獎勵：+%d 🪙" % burst_reward
	burst_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	burst_label.add_theme_font_size_override("font_size", 16)
	burst_label.add_theme_color_override("font_color", STACK_COLOR_GOLD)
	vbox.add_child(burst_label)

	var total_label = Label.new()
	total_label.text = "總獎勵：+%d 🪙" % total_reward
	total_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	total_label.add_theme_font_size_override("font_size", 14)
	total_label.add_theme_color_override("font_color", Color(0.0, 1.0, 0.53))
	vbox.add_child(total_label)

	var vp_size = get_viewport().size
	popup.position = Vector2(vp_size.x, vp_size.y / 2 - 75)
	popup.modulate.a = 0.0
	add_child(popup)

	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 320, 0.3)
	tween.parallel().tween_property(popup, "modulate:a", 1.0, 0.3)
	tween.tween_interval(5.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.4)
	tween.tween_callback(popup.queue_free)
