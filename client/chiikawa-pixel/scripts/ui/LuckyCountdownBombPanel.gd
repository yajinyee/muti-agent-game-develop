## LuckyCountdownBombPanel.gd — 幸運倒數炸彈魚 UI（DAY-268）
## 紅橙炸彈主題面板（全服合力充能倒數爆炸）
## 主色：#FF4500 火橙 + #FFD700 金 + #FF0000 紅 + #1A1A2E 深藍黑
extends CanvasLayer

const BOMB_COLOR_FIRE   = Color(1.0, 0.27, 0.0, 1.0)   # #FF4500 火橙
const BOMB_COLOR_GOLD   = Color(1.0, 0.85, 0.0, 1.0)   # #FFD700 金
const BOMB_COLOR_RED    = Color(1.0, 0.0, 0.0, 1.0)    # #FF0000 紅
const BOMB_COLOR_BURST  = Color(1.0, 1.0, 0.0, 1.0)    # #FFFF00 爆炸黃

var _flash_overlay: ColorRect
var _banner: PanelContainer
var _banner_label: Label
var _charge_counter: Label    # 右上角充能計數器
var _charge_bar: ColorRect    # 底部充能進度條
var _charge_bar_bg: ColorRect # 進度條背景
var _countdown_label: Label   # 倒數計時顯示
var _session_active: bool = false
var _countdown: float = 10.0
var _start_time: float = 0.0
var _max_charge: int = 10

func _ready() -> void:
	layer = 41
	_build_ui()

func _build_ui() -> void:
	var vp_size = get_viewport().size

	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 0.27, 0.0, 0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner = PanelContainer.new()
	_banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 56)
	_banner.modulate.a = 0.0
	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.15, 0.04, 0.0, 0.88)
	banner_style.border_color = BOMB_COLOR_FIRE
	banner_style.set_border_width_all(2)
	banner_style.set_corner_radius_all(6)
	_banner.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_font_size_override("font_size", 18)
	_banner_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.8))
	_banner.add_child(_banner_label)

	# 充能計數器（右上角）
	_charge_counter = Label.new()
	_charge_counter.text = "💣 充能 0/10"
	_charge_counter.position = Vector2(vp_size.x - 180, 65)
	_charge_counter.add_theme_font_size_override("font_size", 20)
	_charge_counter.add_theme_color_override("font_color", BOMB_COLOR_FIRE)
	_charge_counter.modulate.a = 0.0
	add_child(_charge_counter)

	# 倒數計時顯示（右上角，充能計數器下方）
	_countdown_label = Label.new()
	_countdown_label.text = "⏱ 10.0s"
	_countdown_label.position = Vector2(vp_size.x - 160, 92)
	_countdown_label.add_theme_font_size_override("font_size", 16)
	_countdown_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.8))
	_countdown_label.modulate.a = 0.0
	add_child(_countdown_label)

	# 底部充能進度條背景
	_charge_bar_bg = ColorRect.new()
	_charge_bar_bg.color = Color(0.15, 0.04, 0.0, 0.7)
	_charge_bar_bg.size = Vector2(vp_size.x, 12)
	_charge_bar_bg.position = Vector2(0, vp_size.y - 24)
	_charge_bar_bg.modulate.a = 0.0
	add_child(_charge_bar_bg)

	# 底部充能進度條
	_charge_bar = ColorRect.new()
	_charge_bar.color = BOMB_COLOR_FIRE
	_charge_bar.size = Vector2(0, 12)
	_charge_bar.position = Vector2(0, vp_size.y - 24)
	_charge_bar.modulate.a = 0.0
	add_child(_charge_bar)

func _process(_delta: float) -> void:
	if not _session_active:
		return
	var elapsed = Time.get_ticks_msec() / 1000.0 - _start_time
	var remaining = max(0.0, _countdown - elapsed)
	_countdown_label.text = "⏱ %.1fs" % remaining
	# 剩餘時間少時倒數變紅色並閃爍
	if remaining < 3.0:
		_countdown_label.add_theme_color_override("font_color", BOMB_COLOR_RED)
	elif remaining < 5.0:
		_countdown_label.add_theme_color_override("font_color", BOMB_COLOR_FIRE)
	else:
		_countdown_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.8))

## 處理來自 GameManager 的倒數炸彈事件
func handle_event(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"bomb_start":
			_on_bomb_start(payload)
		"bomb_charge":
			_on_bomb_charge(payload)
		"bomb_explode":
			_on_bomb_explode(payload)

## 倒數炸彈觸發（全服）
func _on_bomb_start(payload: Dictionary) -> void:
	_countdown = payload.get("countdown", 10.0)
	_max_charge = payload.get("max_charge", 10)
	_start_time = Time.get_ticks_msec() / 1000.0
	_session_active = true

	var player_name = payload.get("player_name", "玩家")
	var burst_mult = payload.get("burst_mult", 3.0)

	# 火橙三次強閃光
	_flash_triple(Color(1.0, 0.27, 0.0, 0.45))

	# 顯示橫幅
	_banner_label.text = "💣 %s 觸發倒數炸彈！10 秒倒數！全服充能越多爆炸越強！滿充能 ×%.1f！" % [player_name, burst_mult]
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)
	tween.tween_interval(5.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.5)

	# 顯示充能計數器
	_charge_counter.text = "💣 充能 0/%d" % _max_charge
	_charge_counter.add_theme_color_override("font_color", BOMB_COLOR_FIRE)
	var counter_tween = create_tween()
	counter_tween.tween_property(_charge_counter, "modulate:a", 1.0, 0.3)
	counter_tween.parallel().tween_property(_countdown_label, "modulate:a", 1.0, 0.3)

	# 顯示進度條
	var bar_tween = create_tween()
	bar_tween.tween_property(_charge_bar_bg, "modulate:a", 1.0, 0.3)
	bar_tween.parallel().tween_property(_charge_bar, "modulate:a", 1.0, 0.3)

	# 大字提示
	_show_big_text("💣 倒數炸彈！", BOMB_COLOR_FIRE)

## 充能更新（全服）
func _on_bomb_charge(payload: Dictionary) -> void:
	var charge_count = payload.get("charge_count", 0)
	var max_charge = payload.get("max_charge", 10)
	var charger_name = payload.get("charger_name", "玩家")

	# 更新充能計數器
	var charge_ratio = float(charge_count) / float(max_charge)
	var counter_color: Color
	if charge_ratio >= 0.8:
		counter_color = BOMB_COLOR_GOLD
	elif charge_ratio >= 0.5:
		counter_color = BOMB_COLOR_FIRE
	else:
		counter_color = Color(1.0, 0.5, 0.2)
	_charge_counter.text = "💣 充能 %d/%d" % [charge_count, max_charge]
	_charge_counter.add_theme_color_override("font_color", counter_color)

	# 更新進度條
	var vp_size = get_viewport().size
	var bar_width = vp_size.x * charge_ratio
	var bar_tween = create_tween()
	bar_tween.tween_property(_charge_bar, "size:x", bar_width, 0.12)
	if charge_ratio >= 0.8:
		bar_tween.parallel().tween_property(_charge_bar, "color", BOMB_COLOR_GOLD, 0.12)
	else:
		bar_tween.parallel().tween_property(_charge_bar, "color", BOMB_COLOR_FIRE, 0.12)

	# 計數器脈衝
	var pulse = create_tween()
	pulse.tween_property(_charge_counter, "scale", Vector2(1.2, 1.2), 0.06)
	pulse.tween_property(_charge_counter, "scale", Vector2(1.0, 1.0), 0.09)

	# 輕微閃光
	_flash_overlay.color = Color(1.0, 0.27, 0.0, 0)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color:a", 0.1, 0.04)
	flash_tween.tween_property(_flash_overlay, "color:a", 0.0, 0.08)

	# 浮動文字
	_show_float_text("💣 %s 充能！%d/%d" % [charger_name, charge_count, max_charge], counter_color)

## 炸彈爆炸（全服）
func _on_bomb_explode(payload: Dictionary) -> void:
	var charge_count = payload.get("charge_count", 0)
	var mult = payload.get("mult", 1.0)
	var total_reward = payload.get("total_reward", 0)
	var is_burst = payload.get("is_burst", false)

	_session_active = false

	if is_burst:
		# 滿充能爆炸：全螢幕三次強閃光（金色）
		_flash_triple(Color(1.0, 0.85, 0.0, 0.75))
		_show_big_text("💣 滿充能爆炸！×%.1f！" % mult, BOMB_COLOR_GOLD)
	else:
		# 普通爆炸：全螢幕三次強閃光（火橙）
		_flash_triple(Color(1.0, 0.27, 0.0, 0.6))
		_show_big_text("💣 炸彈爆炸！充能 %d 次 ×%.1f！" % [charge_count, mult], BOMB_COLOR_FIRE)

	# 結算彈窗
	_show_explode_popup(charge_count, mult, total_reward, is_burst)

	# 隱藏計時和進度條
	_hide_session_ui()

## 工具函數：隱藏 session UI
func _hide_session_ui() -> void:
	var tween = create_tween()
	tween.tween_property(_charge_counter, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_countdown_label, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_charge_bar, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_charge_bar_bg, "modulate:a", 0.0, 0.5)

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

## 工具函數：顯示爆炸結算彈窗
func _show_explode_popup(charge_count: int, mult: float, total_reward: int, is_burst: bool) -> void:
	var popup = PanelContainer.new()
	popup.custom_minimum_size = Vector2(280, 130)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.03, 0.0, 0.92)
	var border_color = BOMB_COLOR_GOLD if is_burst else BOMB_COLOR_FIRE
	style.border_color = border_color
	style.set_border_width_all(2)
	style.set_corner_radius_all(8)
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	popup.add_child(vbox)

	var title_label = Label.new()
	title_label.text = "💣 滿充能爆炸！" if is_burst else "💣 炸彈爆炸！"
	title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_label.add_theme_font_size_override("font_size", 18)
	title_label.add_theme_color_override("font_color", border_color)
	vbox.add_child(title_label)

	var charge_label = Label.new()
	charge_label.text = "充能次數：%d / 10" % charge_count
	charge_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	charge_label.add_theme_font_size_override("font_size", 14)
	charge_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.8))
	vbox.add_child(charge_label)

	var mult_label = Label.new()
	mult_label.text = "爆炸倍率：×%.1f" % mult
	mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mult_label.add_theme_font_size_override("font_size", 15)
	mult_label.add_theme_color_override("font_color", BOMB_COLOR_FIRE)
	vbox.add_child(mult_label)

	var reward_label = Label.new()
	reward_label.text = "全服 AOE：+%d 🪙" % total_reward
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
	tween.tween_interval(4.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.4)
	tween.tween_callback(popup.queue_free)
