## LuckyCrashFishPanel.gd — T130 幸運崩潰魚 UI
## lucky-panel-agent 負責維護
## 業界依據：Lucky Fish by AbraCadabra「crash mechanic — multiplier rises until crash, cash out anytime」
## 視覺主題：深紅色 + 上升倍率計數器 + 收割按鈕 + 崩潰爆炸演出
extends CanvasLayer

const LAYER_Z = 30  # CanvasLayer layer 值

# 顏色主題
const COLOR_CRASH    = Color(0.8, 0.1, 0.1)    # 深紅（崩潰主色）
const COLOR_RISING   = Color(1.0, 0.5, 0.1)    # 火橙（倍率上升）
const COLOR_HARVEST  = Color(0.2, 0.9, 0.2)    # 綠色（收割成功）
const COLOR_PERFECT  = Color(1.0, 0.85, 0.0)   # 金色（完美收割）
const COLOR_BG       = Color(0.08, 0.0, 0.0, 0.90)

var _banner: Control = null
var _crash_panel: Control = null
var _result_popup: Control = null
var _flash_overlay: ColorRect = null
var _mult_label: Label = null
var _timer_label: Label = null
var _harvest_btn: Button = null
var _crash_active: bool = false
var _current_mult: float = 1.0
var _my_player_id: String = ""

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_crash_panel()
	GameManager.lucky_crash_fish.connect(_on_lucky_crash_fish)

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(0.8, 0.1, 0.1, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

func _create_crash_panel() -> void:
	# 右上角崩潰倍率顯示
	_crash_panel = Control.new()
	_crash_panel.position = Vector2(960, 10)
	_crash_panel.size = Vector2(300, 120)
	_crash_panel.visible = false
	add_child(_crash_panel)

	var bg = ColorRect.new()
	bg.size = _crash_panel.size
	bg.color = COLOR_BG
	_crash_panel.add_child(bg)

	var title = Label.new()
	title.text = "💥 崩潰倍率進行中"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_CRASH
	_crash_panel.add_child(title)

	# 倍率大字
	_mult_label = Label.new()
	_mult_label.name = "MultLabel"
	_mult_label.text = "×1.0"
	_mult_label.position = Vector2(8, 24)
	_mult_label.size = Vector2(200, 50)
	_mult_label.add_theme_font_size_override("font_size", 36)
	_mult_label.modulate = COLOR_RISING
	_crash_panel.add_child(_mult_label)

	# 計時器
	_timer_label = Label.new()
	_timer_label.name = "TimerLabel"
	_timer_label.text = "??s"
	_timer_label.position = Vector2(220, 30)
	_timer_label.add_theme_font_size_override("font_size", 16)
	_timer_label.modulate = Color(0.8, 0.5, 0.5)
	_crash_panel.add_child(_timer_label)

	# 收割按鈕（只有觸發者才能點）
	_harvest_btn = Button.new()
	_harvest_btn.name = "HarvestBtn"
	_harvest_btn.text = "💰 收割！"
	_harvest_btn.position = Vector2(8, 80)
	_harvest_btn.size = Vector2(284, 32)
	_harvest_btn.add_theme_font_size_override("font_size", 16)
	_harvest_btn.modulate = COLOR_HARVEST
	_harvest_btn.visible = false
	_harvest_btn.pressed.connect(_on_harvest_pressed)
	_crash_panel.add_child(_harvest_btn)

func _on_harvest_pressed() -> void:
	if not _crash_active:
		return
	# 發送收割訊息
	NetworkManager.send("crash_harvest", {})
	_harvest_btn.disabled = true
	_harvest_btn.text = "已收割！"

func _show_start_banner(name: String, crash_in: float) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()
	_banner = Control.new()
	_banner.position = Vector2(0, 100)
	_banner.size = Vector2(1280, 90)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = COLOR_BG
	_banner.add_child(bg)

	# 頂部深紅線
	var top_line = ColorRect.new()
	top_line.size = Vector2(1280, 3)
	top_line.color = COLOR_CRASH
	_banner.add_child(top_line)

	var lbl = Label.new()
	lbl.text = "💥 %s 觸發崩潰倍率！倍率持續上升！" % name
	lbl.position = Vector2(0, 10)
	lbl.size = Vector2(1280, 40)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 28)
	lbl.modulate = COLOR_CRASH
	_banner.add_child(lbl)

	var sub = Label.new()
	sub.text = "在崩潰前收割！收割 ≥×5.0 → 完美收割！全服 ×2.0！"
	sub.position = Vector2(0, 50)
	sub.size = Vector2(1280, 30)
	sub.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	sub.add_theme_font_size_override("font_size", 16)
	sub.modulate = Color(1.0, 0.7, 0.7)
	_banner.add_child(sub)

	# 底部深紅線
	var bot_line = ColorRect.new()
	bot_line.size = Vector2(1280, 3)
	bot_line.position = Vector2(0, 87)
	bot_line.color = COLOR_CRASH
	_banner.add_child(bot_line)

	# 滑入動畫
	_banner.position.y = -90
	var tween = _banner.create_tween()
	tween.tween_property(_banner, "position:y", 100.0, 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.tween_interval(3.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(_banner):
			_banner.queue_free()
			_banner = null
	)

func _show_harvest_popup(name: String, mult: float, reward: int) -> void:
	if is_instance_valid(_result_popup):
		_result_popup.queue_free()

	_do_flash(COLOR_HARVEST, 3)

	_result_popup = Control.new()
	_result_popup.position = Vector2(340, 200)
	_result_popup.size = Vector2(600, 160)
	add_child(_result_popup)

	var bg = ColorRect.new()
	bg.size = _result_popup.size
	bg.color = COLOR_BG
	_result_popup.add_child(bg)

	var title_lbl = Label.new()
	title_lbl.text = "💰 %s 收割！×%.1f" % [name, mult]
	title_lbl.position = Vector2(0, 20)
	title_lbl.size = Vector2(600, 50)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 28)
	title_lbl.modulate = COLOR_HARVEST
	_result_popup.add_child(title_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "獲得 %d 金幣！" % reward
	reward_lbl.position = Vector2(0, 90)
	reward_lbl.size = Vector2(600, 50)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_lbl.add_theme_font_size_override("font_size", 24)
	reward_lbl.modulate = COLOR_PERFECT
	_result_popup.add_child(reward_lbl)

	# 彈出動畫
	_result_popup.scale = Vector2.ZERO
	_result_popup.pivot_offset = Vector2(300, 80)
	var tween = _result_popup.create_tween()
	tween.tween_property(_result_popup, "scale", Vector2(1.0, 1.0), 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
	tween.tween_interval(2.5)
	tween.tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_result_popup):
			_result_popup.queue_free()
			_result_popup = null
	)

func _show_crash_popup(name: String, mult: float) -> void:
	if is_instance_valid(_result_popup):
		_result_popup.queue_free()

	_do_flash(COLOR_CRASH, 5)

	_result_popup = Control.new()
	_result_popup.position = Vector2(340, 200)
	_result_popup.size = Vector2(600, 140)
	add_child(_result_popup)

	var bg = ColorRect.new()
	bg.size = _result_popup.size
	bg.color = COLOR_BG
	_result_popup.add_child(bg)

	var title_lbl = Label.new()
	title_lbl.text = "💥 崩潰！×%.1f 歸零！" % mult
	title_lbl.position = Vector2(0, 20)
	title_lbl.size = Vector2(600, 50)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 30)
	title_lbl.modulate = COLOR_CRASH
	_result_popup.add_child(title_lbl)

	var sub_lbl = Label.new()
	sub_lbl.text = "%s 沒有及時收割！" % name
	sub_lbl.position = Vector2(0, 80)
	sub_lbl.size = Vector2(600, 40)
	sub_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	sub_lbl.add_theme_font_size_override("font_size", 18)
	sub_lbl.modulate = Color(0.8, 0.5, 0.5)
	_result_popup.add_child(sub_lbl)

	# 彈出動畫
	_result_popup.scale = Vector2.ZERO
	_result_popup.pivot_offset = Vector2(300, 70)
	var tween = _result_popup.create_tween()
	tween.tween_property(_result_popup, "scale", Vector2(1.0, 1.0), 0.2).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.tween_interval(2.0)
	tween.tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_result_popup):
			_result_popup.queue_free()
			_result_popup = null
	)

func _show_perfect_popup(name: String, mult: float, boost: float, secs: int) -> void:
	if is_instance_valid(_result_popup):
		_result_popup.queue_free()

	_do_flash(COLOR_PERFECT, 5)

	_result_popup = Control.new()
	_result_popup.position = Vector2(340, 180)
	_result_popup.size = Vector2(600, 200)
	add_child(_result_popup)

	var bg = ColorRect.new()
	bg.size = _result_popup.size
	bg.color = COLOR_BG
	_result_popup.add_child(bg)

	var title_lbl = Label.new()
	title_lbl.text = "💰✨ 完美收割！"
	title_lbl.position = Vector2(0, 15)
	title_lbl.size = Vector2(600, 50)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 32)
	title_lbl.modulate = COLOR_PERFECT
	_result_popup.add_child(title_lbl)

	var detail_lbl = Label.new()
	detail_lbl.text = "%s 在 ×%.1f 時收割！" % [name, mult]
	detail_lbl.position = Vector2(0, 70)
	detail_lbl.size = Vector2(600, 30)
	detail_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	detail_lbl.add_theme_font_size_override("font_size", 18)
	detail_lbl.modulate = Color(0.9, 0.9, 0.9)
	_result_popup.add_child(detail_lbl)

	var boost_lbl = Label.new()
	boost_lbl.text = "全服 ×%.0f 加成 %d 秒！" % [boost, secs]
	boost_lbl.position = Vector2(0, 120)
	boost_lbl.size = Vector2(600, 50)
	boost_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	boost_lbl.add_theme_font_size_override("font_size", 24)
	boost_lbl.modulate = COLOR_HARVEST
	_result_popup.add_child(boost_lbl)

	# 彈出動畫
	_result_popup.scale = Vector2.ZERO
	_result_popup.pivot_offset = Vector2(300, 100)
	var tween = _result_popup.create_tween()
	tween.tween_property(_result_popup, "scale", Vector2(1.0, 1.0), 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
	tween.tween_interval(3.5)
	tween.tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_result_popup):
			_result_popup.queue_free()
			_result_popup = null
	)

func _do_flash(color: Color, count: int) -> void:
	if not is_instance_valid(_flash_overlay):
		return
	var tween = _flash_overlay.create_tween()
	for i in count:
		tween.tween_property(_flash_overlay, "color", Color(color.r, color.g, color.b, 0.5), 0.08)
		tween.tween_property(_flash_overlay, "color", Color(color.r, color.g, color.b, 0.0), 0.12)

func _on_lucky_crash_fish(data: Dictionary) -> void:
	var event = data.get("event", "")
	var player_id = data.get("player_id", "")
	var name = data.get("player_name", "玩家")
	var mult = data.get("current_mult", 1.0)
	var is_me = (player_id == GameManager.get_player_id())

	match event:
		"crash_start":
			_crash_active = true
			_current_mult = 1.0
			_my_player_id = player_id
			_show_start_banner(name, data.get("crash_in", 10.0))
			if is_instance_valid(_crash_panel):
				_crash_panel.visible = true
			if is_instance_valid(_mult_label):
				_mult_label.text = "×1.0"
				_mult_label.modulate = COLOR_RISING
			if is_instance_valid(_harvest_btn):
				_harvest_btn.visible = is_me
				_harvest_btn.disabled = false
				_harvest_btn.text = "💰 收割！"
			_do_flash(COLOR_CRASH, 3)
			ScreenShake.add_trauma(0.4)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"mult_rise":
			_current_mult = mult
			if is_instance_valid(_mult_label):
				_mult_label.text = "×%.1f" % mult
				# 顏色隨倍率變化
				if mult >= 7.0:
					_mult_label.modulate = COLOR_PERFECT
				elif mult >= 5.0:
					_mult_label.modulate = Color(1.0, 0.85, 0.0)
				elif mult >= 3.0:
					_mult_label.modulate = Color(1.0, 0.6, 0.2)
				else:
					_mult_label.modulate = COLOR_RISING
			if is_instance_valid(_timer_label):
				var tl = data.get("time_left", 0.0)
				_timer_label.text = "%.0fs" % tl
		"harvest":
			_crash_active = false
			if is_instance_valid(_crash_panel):
				_crash_panel.visible = false
			if is_instance_valid(_harvest_btn):
				_harvest_btn.visible = false
			var reward = data.get("reward", 0)
			_show_harvest_popup(name, mult, reward)
			ScreenShake.add_trauma(0.3)
		"crash":
			_crash_active = false
			if is_instance_valid(_crash_panel):
				_crash_panel.visible = false
			if is_instance_valid(_harvest_btn):
				_harvest_btn.visible = false
			_show_crash_popup(name, mult)
			_do_flash(COLOR_CRASH, 4)
			ScreenShake.add_trauma(0.6)
		"perfect_harvest":
			var boost = data.get("boost_mult", 2.0)
			var secs = data.get("boost_secs", 5)
			_show_perfect_popup(name, mult, boost, secs)
			ScreenShake.add_trauma(0.7)
		"perfect_end":
			pass  # HUD 已處理橫幅
