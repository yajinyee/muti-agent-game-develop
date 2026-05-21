## DragonWrathPanel.gd
##
## 龍怒蓄力大招面板（DAY-128）
## 業界依據：JILI Royal Fishing 2026 Dragon Wrath — 累積怒氣值釋放全螢幕大招
## 位置：左下角怒氣條 + 大招按鈕
## 設計：
##   - 怒氣條（0-100）：漸層顏色（藍→紫→紅）
##   - 滿怒氣時按鈕閃爍，提示玩家釋放
##   - 大招開始：全螢幕龍形閃光 + 橫幅
##   - 大招結果：顯示擊破數和獎勵

extends Control

# ---- 常數 ----
const MAX_CHARGE := 100
const BAR_WIDTH := 160.0
const BAR_HEIGHT := 16.0
const PANEL_WIDTH := 180.0
const PANEL_HEIGHT := 80.0

# ---- 節點引用 ----
var _bar_bg: ColorRect = null
var _bar_fill: ColorRect = null
var _bar_label: Label = null
var _wrath_btn: Button = null
var _charge_label: Label = null
var _pixel_font: FontFile = null

# ---- 狀態 ----
var _charge: int = 0
var _is_ready: bool = false
var _cooldown: int = 0
var _btn_flash_time: float = 0.0

func setup(font: FontFile) -> void:
	_pixel_font = font
	_build_ui()
	_connect_signals()

func _build_ui() -> void:
	# 面板背景
	var bg = Panel.new()
	bg.position = Vector2(0, 0)
	bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	var bg_style = StyleBoxFlat.new()
	bg_style.bg_color = Color(0.05, 0.05, 0.15, 0.85)
	bg_style.border_color = Color(0.4, 0.2, 0.8, 0.7)
	bg_style.set_border_width_all(1)
	bg_style.corner_radius_top_right = 6
	bg_style.corner_radius_bottom_right = 6
	bg.add_theme_stylebox_override("panel", bg_style)
	add_child(bg)

	# 標題
	var title = Label.new()
	title.position = Vector2(8, 4)
	title.size = Vector2(164, 16)
	title.text = "🐉 龍怒"
	title.add_theme_color_override("font_color", Color(0.8, 0.5, 1.0))
	title.add_theme_font_size_override("font_size", 11)
	if is_instance_valid(_pixel_font):
		title.add_theme_font_override("font", _pixel_font)
	add_child(title)

	# 怒氣條背景
	_bar_bg = ColorRect.new()
	_bar_bg.position = Vector2(8, 22)
	_bar_bg.size = Vector2(BAR_WIDTH, BAR_HEIGHT)
	_bar_bg.color = Color(0.1, 0.1, 0.2)
	add_child(_bar_bg)

	# 怒氣條填充
	_bar_fill = ColorRect.new()
	_bar_fill.position = Vector2(8, 22)
	_bar_fill.size = Vector2(0, BAR_HEIGHT)
	_bar_fill.color = Color(0.3, 0.2, 0.9)
	add_child(_bar_fill)

	# 怒氣值標籤（疊在條上）
	_bar_label = Label.new()
	_bar_label.position = Vector2(8, 22)
	_bar_label.size = Vector2(BAR_WIDTH, BAR_HEIGHT)
	_bar_label.text = "0 / 100"
	_bar_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0, 0.9))
	_bar_label.add_theme_font_size_override("font_size", 10)
	_bar_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_bar_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	if is_instance_valid(_pixel_font):
		_bar_label.add_theme_font_override("font", _pixel_font)
	add_child(_bar_label)

	# 大招按鈕
	_wrath_btn = Button.new()
	_wrath_btn.position = Vector2(8, 44)
	_wrath_btn.size = Vector2(BAR_WIDTH, 28)
	_wrath_btn.text = "🐉 大討伐"
	_wrath_btn.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))
	_wrath_btn.add_theme_font_size_override("font_size", 12)
	_wrath_btn.disabled = true
	if is_instance_valid(_pixel_font):
		_wrath_btn.add_theme_font_override("font", _pixel_font)
	_wrath_btn.pressed.connect(_on_wrath_btn_pressed)
	add_child(_wrath_btn)

func _connect_signals() -> void:
	if GameManager.has_signal("wrath_updated"):
		GameManager.wrath_updated.connect(_on_wrath_updated)
	if GameManager.has_signal("wrath_started"):
		GameManager.wrath_started.connect(_on_wrath_started)
	if GameManager.has_signal("wrath_result"):
		GameManager.wrath_result.connect(_on_wrath_result)

# ---- 事件處理 ----

func _on_wrath_updated(data: Dictionary) -> void:
	_charge = data.get("charge", 0)
	_is_ready = data.get("is_ready", false)
	_cooldown = data.get("cooldown", 0)
	_update_display()

func _on_wrath_btn_pressed() -> void:
	if not _is_ready:
		return
	if NetworkManager.has_method("send_use_wrath"):
		NetworkManager.send_use_wrath()

func _on_wrath_started(data: Dictionary) -> void:
	var player_name = data.get("player_name", "玩家")
	var msg = data.get("message", player_name + " 釋放了大討伐！")
	_show_wrath_effect(msg)

func _on_wrath_result(data: Dictionary) -> void:
	var killed = data.get("killed_count", 0)
	var reward = data.get("total_reward", 0)
	var is_me = data.get("player_id", "") == GameManager.get_meta("player_id", "")
	if is_me:
		_show_result_popup(killed, reward)

# ---- 視覺更新 ----

func _update_display() -> void:
	if not is_instance_valid(_bar_fill):
		return

	# 更新怒氣條寬度
	var ratio = float(_charge) / float(MAX_CHARGE)
	var tween = create_tween()
	tween.tween_property(_bar_fill, "size:x", BAR_WIDTH * ratio, 0.15)

	# 更新怒氣條顏色（藍→紫→紅）
	var bar_color: Color
	if ratio < 0.5:
		bar_color = Color(0.3, 0.2, 0.9).lerp(Color(0.7, 0.2, 0.9), ratio * 2.0)
	else:
		bar_color = Color(0.7, 0.2, 0.9).lerp(Color(1.0, 0.2, 0.2), (ratio - 0.5) * 2.0)
	_bar_fill.color = bar_color

	# 更新標籤
	if is_instance_valid(_bar_label):
		_bar_label.text = str(_charge) + " / " + str(MAX_CHARGE)

	# 更新按鈕狀態
	if is_instance_valid(_wrath_btn):
		if _is_ready:
			_wrath_btn.disabled = false
			_wrath_btn.text = "🐉 大討伐！"
			_wrath_btn.add_theme_color_override("font_color", Color(1.0, 0.8, 0.2))
		elif _cooldown > 0:
			_wrath_btn.disabled = true
			_wrath_btn.text = "冷卻 %ds" % _cooldown
			_wrath_btn.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))
		else:
			_wrath_btn.disabled = true
			_wrath_btn.text = "🐉 大討伐"
			_wrath_btn.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))

func _show_wrath_effect(msg: String) -> void:
	# 全螢幕龍形閃光（紫紅色）
	var flash = ColorRect.new()
	flash.color = Color(0.6, 0.1, 0.9, 0.35)
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.z_index = 200
	get_tree().root.add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.6)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)

	# 頂部橫幅
	var banner = Label.new()
	banner.text = "🐉 " + msg
	banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	banner.size.y = 40
	banner.add_theme_color_override("font_color", Color(1.0, 0.7, 0.2))
	banner.add_theme_font_size_override("font_size", 20)
	banner.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	banner.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	banner.z_index = 201
	banner.modulate.a = 0.0
	get_tree().root.add_child(banner)
	var tween2 = create_tween()
	tween2.tween_property(banner, "modulate:a", 1.0, 0.2)
	tween2.tween_interval(2.5)
	tween2.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween2.tween_callback(func():
		if is_instance_valid(banner):
			banner.queue_free()
	)

func _show_result_popup(killed: int, reward: int) -> void:
	var popup = Label.new()
	popup.text = "🐉 大討伐！擊破 %d 個目標，獲得 %d 金幣！" % [killed, reward]
	popup.position = Vector2(640 - 200, 360)
	popup.size = Vector2(400, 40)
	popup.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	popup.add_theme_font_size_override("font_size", 16)
	popup.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup.z_index = 202
	popup.modulate.a = 0.0
	get_tree().root.add_child(popup)
	var tween = create_tween()
	tween.tween_property(popup, "modulate:a", 1.0, 0.2)
	tween.tween_property(popup, "position:y", 300.0, 1.0).set_ease(Tween.EASE_OUT)
	tween.tween_interval(1.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(popup):
			popup.queue_free()
	)

# ---- 每幀更新（按鈕閃爍）----

func _process(delta: float) -> void:
	if not _is_ready or not is_instance_valid(_wrath_btn):
		return
	_btn_flash_time += delta
	var blink = fmod(_btn_flash_time * 2.0, 1.0) > 0.5
	_wrath_btn.add_theme_color_override("font_color",
		Color(1.0, 0.9, 0.2) if blink else Color(1.0, 0.6, 0.1))
