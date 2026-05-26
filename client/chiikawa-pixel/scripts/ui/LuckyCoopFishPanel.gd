## LuckyCoopFishPanel.gd — T127 幸運全服合作魚 UI
## lucky-panel-agent 負責維護
## 業界依據：Fishing Fortune 2026「Cooperative Challenge — all players contribute damage within time limit」
## 視覺主題：青藍色 + 合作進度條 + 全服廣播 + 握手圖案
extends CanvasLayer

const LAYER_Z = 27  # CanvasLayer layer 值

# 顏色主題
const COLOR_COOP    = Color(0.0, 0.9, 1.0)   # 青藍（合作主色）
const COLOR_SUCCESS = Color(0.2, 1.0, 0.5)   # 綠色（成功）
const COLOR_TIMEOUT = Color(0.6, 0.6, 0.6)   # 灰色（超時）
const COLOR_BG      = Color(0.0, 0.05, 0.1, 0.88)

var _banner: Control = null
var _progress_panel: Control = null
var _result_popup: Control = null
var _flash_overlay: ColorRect = null
var _progress_bar: ColorRect = null
var _progress_label: Label = null
var _timer_label: Label = null

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_progress_panel()
	GameManager.lucky_coop_fish.connect(_on_lucky_coop_fish)

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(0.0, 0.9, 1.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

func _create_progress_panel() -> void:
	# 右上角合作進度顯示
	_progress_panel = Control.new()
	_progress_panel.position = Vector2(960, 10)
	_progress_panel.size = Vector2(300, 90)
	_progress_panel.visible = false
	add_child(_progress_panel)

	var bg = ColorRect.new()
	bg.size = _progress_panel.size
	bg.color = COLOR_BG
	_progress_panel.add_child(bg)

	var title = Label.new()
	title.text = "🤝 全服合作挑戰"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_COOP
	_progress_panel.add_child(title)

	# 進度條背景
	var bar_bg = ColorRect.new()
	bar_bg.position = Vector2(8, 26)
	bar_bg.size = Vector2(284, 18)
	bar_bg.color = Color(0.1, 0.1, 0.1)
	_progress_panel.add_child(bar_bg)

	# 進度條填充
	_progress_bar = ColorRect.new()
	_progress_bar.position = Vector2(8, 26)
	_progress_bar.size = Vector2(0, 18)
	_progress_bar.color = COLOR_COOP
	_progress_panel.add_child(_progress_bar)

	# 進度文字
	_progress_label = Label.new()
	_progress_label.name = "ProgressLabel"
	_progress_label.text = "0/8 點"
	_progress_label.position = Vector2(8, 48)
	_progress_label.add_theme_font_size_override("font_size", 14)
	_progress_label.modulate = COLOR_COOP
	_progress_panel.add_child(_progress_label)

	# 計時器
	_timer_label = Label.new()
	_timer_label.name = "TimerLabel"
	_timer_label.text = "20s"
	_timer_label.position = Vector2(220, 48)
	_timer_label.add_theme_font_size_override("font_size", 14)
	_timer_label.modulate = Color(1.0, 0.9, 0.5)
	_progress_panel.add_child(_timer_label)

func _update_progress(current: int, target: int, time_left: float) -> void:
	if not is_instance_valid(_progress_panel):
		return
	_progress_panel.visible = true
	var pct = float(current) / float(max(target, 1))
	if is_instance_valid(_progress_bar):
		_progress_bar.size.x = 284.0 * pct
		# 顏色隨進度變化
		if pct >= 0.8:
			_progress_bar.color = COLOR_SUCCESS
		elif pct >= 0.5:
			_progress_bar.color = Color(0.5, 1.0, 0.7)
		else:
			_progress_bar.color = COLOR_COOP
	if is_instance_valid(_progress_label):
		_progress_label.text = "%d/%d 點" % [current, target]
	if is_instance_valid(_timer_label):
		_timer_label.text = "%.0fs" % time_left
		if time_left <= 5.0:
			_timer_label.modulate = Color(1.0, 0.3, 0.3)
		else:
			_timer_label.modulate = Color(1.0, 0.9, 0.5)

func _show_start_banner(name: String, target: int) -> void:
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

	# 頂部青藍線
	var top_line = ColorRect.new()
	top_line.size = Vector2(1280, 3)
	top_line.color = COLOR_COOP
	_banner.add_child(top_line)

	var lbl = Label.new()
	lbl.text = "🤝 %s 發起全服合作挑戰！" % name
	lbl.position = Vector2(0, 10)
	lbl.size = Vector2(1280, 40)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 28)
	lbl.modulate = COLOR_COOP
	_banner.add_child(lbl)

	var sub = Label.new()
	sub.text = "目標：%d 點！20 秒內全服合力達成！" % target
	sub.position = Vector2(0, 50)
	sub.size = Vector2(1280, 30)
	sub.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	sub.add_theme_font_size_override("font_size", 16)
	sub.modulate = Color(0.8, 0.95, 1.0)
	_banner.add_child(sub)

	# 底部青藍線
	var bot_line = ColorRect.new()
	bot_line.size = Vector2(1280, 3)
	bot_line.position = Vector2(0, 87)
	bot_line.color = COLOR_COOP
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

func _show_success_popup(boost: float, secs: int) -> void:
	if is_instance_valid(_result_popup):
		_result_popup.queue_free()

	_do_flash(COLOR_SUCCESS, 5)

	_result_popup = Control.new()
	_result_popup.position = Vector2(340, 200)
	_result_popup.size = Vector2(600, 180)
	add_child(_result_popup)

	var bg = ColorRect.new()
	bg.size = _result_popup.size
	bg.color = COLOR_BG
	_result_popup.add_child(bg)

	var border = ColorRect.new()
	border.size = _result_popup.size
	border.color = Color(0.0, 0.9, 1.0, 0.25)
	_result_popup.add_child(border)

	var title_lbl = Label.new()
	title_lbl.text = "🤝✨ 全服合作成功！"
	title_lbl.position = Vector2(0, 20)
	title_lbl.size = Vector2(600, 50)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 30)
	title_lbl.modulate = COLOR_SUCCESS
	_result_popup.add_child(title_lbl)

	var boost_lbl = Label.new()
	boost_lbl.text = "全服 ×%.0f 加成 %d 秒！" % [boost, secs]
	boost_lbl.position = Vector2(0, 80)
	boost_lbl.size = Vector2(600, 50)
	boost_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	boost_lbl.add_theme_font_size_override("font_size", 26)
	boost_lbl.modulate = COLOR_COOP
	_result_popup.add_child(boost_lbl)

	# 彈出動畫
	_result_popup.scale = Vector2.ZERO
	_result_popup.pivot_offset = Vector2(300, 90)
	var tween = _result_popup.create_tween()
	tween.tween_property(_result_popup, "scale", Vector2(1.0, 1.0), 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
	tween.tween_interval(3.0)
	tween.tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_result_popup):
			_result_popup.queue_free()
			_result_popup = null
	)

	# 隱藏進度面板
	if is_instance_valid(_progress_panel):
		_progress_panel.visible = false

func _show_boost_indicator(boost: float, secs: int) -> void:
	# 全服加成持續指示器
	var boost_panel = Control.new()
	boost_panel.position = Vector2(0, 200)
	boost_panel.size = Vector2(1280, 50)
	add_child(boost_panel)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.0, 0.2, 0.1, 0.85)
	boost_panel.add_child(bg)

	var lbl = Label.new()
	lbl.text = "🤝 全服合作加成 ×%.0f 進行中！" % boost
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 18)
	lbl.modulate = COLOR_SUCCESS
	boost_panel.add_child(lbl)

	# 脈動動畫
	var tween = boost_panel.create_tween().set_loops(secs * 2)
	tween.tween_property(boost_panel, "modulate:a", 0.6, 0.25)
	tween.tween_property(boost_panel, "modulate:a", 1.0, 0.25)
	boost_panel.create_tween().tween_interval(float(secs)).tween_callback(func():
		if is_instance_valid(boost_panel):
			boost_panel.queue_free()
	)

func _do_flash(color: Color, count: int) -> void:
	if not is_instance_valid(_flash_overlay):
		return
	var tween = _flash_overlay.create_tween()
	for i in count:
		tween.tween_property(_flash_overlay, "color", Color(color.r, color.g, color.b, 0.45), 0.08)
		tween.tween_property(_flash_overlay, "color", Color(color.r, color.g, color.b, 0.0), 0.12)

func _on_lucky_coop_fish(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"coop_start":
			var target = data.get("target_points", 8)
			_show_start_banner(name, target)
			_update_progress(0, target, 20.0)
			_do_flash(COLOR_COOP, 3)
			ScreenShake.add_trauma(0.3)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"coop_progress":
			var current = data.get("current_points", 0)
			var target = data.get("target_points", 8)
			var tl = data.get("time_left", 0.0)
			_update_progress(current, target, tl)
		"coop_success":
			var boost = data.get("boost_mult", 4.0)
			var secs = data.get("boost_secs", 8)
			_show_success_popup(boost, secs)
			_show_boost_indicator(boost, secs)
			_do_flash(COLOR_SUCCESS, 5)
			ScreenShake.add_trauma(0.7)
		"coop_timeout":
			if is_instance_valid(_progress_panel):
				_progress_panel.visible = false
		"coop_boost_end":
			pass  # HUD 已處理橫幅
