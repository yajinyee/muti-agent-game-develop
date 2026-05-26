## LuckyCoopFishPanel.gd — T127 幸運全服合作魚 UI
## lucky-panel-agent 負責維護
## 業界依據：業界原創「全服合作機制 — 所有玩家一起貢獻傷害，達到目標觸發全服大獎」
## 視覺主題：青藍色 + 合作進度條 + 全服廣播
extends CanvasLayer

const LAYER_Z = 27

const COLOR_COOP   = Color(0.0, 0.9, 1.0)    # 青藍
const COLOR_SUCCESS = Color(0.0, 1.0, 0.5)   # 成功綠
const COLOR_BG     = Color(0.0, 0.05, 0.1, 0.88)

var _banner: Control = null
var _progress_panel: Control = null
var _flash_overlay: ColorRect = null

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
	_progress_panel = Control.new()
	_progress_panel.position = Vector2(10, 120)
	_progress_panel.size = Vector2(280, 80)
	_progress_panel.visible = false
	add_child(_progress_panel)

	var bg = ColorRect.new()
	bg.size = _progress_panel.size
	bg.color = COLOR_BG
	_progress_panel.add_child(bg)

	var title = Label.new()
	title.name = "Title"
	title.text = "🤝 全服合作挑戰"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_COOP
	_progress_panel.add_child(title)

	var progress_lbl = Label.new()
	progress_lbl.name = "ProgressLabel"
	progress_lbl.text = "0 / 0 點"
	progress_lbl.position = Vector2(8, 24)
	progress_lbl.add_theme_font_size_override("font_size", 18)
	progress_lbl.modulate = Color.WHITE
	_progress_panel.add_child(progress_lbl)

	# 進度條背景
	var bar_bg = ColorRect.new()
	bar_bg.name = "BarBG"
	bar_bg.size = Vector2(264, 12)
	bar_bg.position = Vector2(8, 48)
	bar_bg.color = Color(0.1, 0.1, 0.2)
	_progress_panel.add_child(bar_bg)

	# 進度條
	var bar = ColorRect.new()
	bar.name = "Bar"
	bar.size = Vector2(0, 12)
	bar.position = Vector2(8, 48)
	bar.color = COLOR_COOP
	_progress_panel.add_child(bar)

	var time_lbl = Label.new()
	time_lbl.name = "TimeLabel"
	time_lbl.text = "20.0s"
	time_lbl.position = Vector2(220, 24)
	time_lbl.add_theme_font_size_override("font_size", 14)
	time_lbl.modulate = Color(0.8, 0.8, 0.8)
	_progress_panel.add_child(time_lbl)

func _update_progress(current: int, target: int, time_left: float) -> void:
	if not is_instance_valid(_progress_panel):
		return
	_progress_panel.visible = true
	var lbl = _progress_panel.get_node_or_null("ProgressLabel")
	if is_instance_valid(lbl):
		lbl.text = "%d / %d 點" % [current, target]
	var bar = _progress_panel.get_node_or_null("Bar")
	var bar_bg = _progress_panel.get_node_or_null("BarBG")
	if is_instance_valid(bar) and is_instance_valid(bar_bg):
		var pct = float(current) / float(max(target, 1))
		bar.size.x = bar_bg.size.x * pct
		if pct >= 0.8:
			bar.color = COLOR_SUCCESS
		elif pct >= 0.5:
			bar.color = Color(0.5, 1.0, 0.5)
		else:
			bar.color = COLOR_COOP
	var time_lbl = _progress_panel.get_node_or_null("TimeLabel")
	if is_instance_valid(time_lbl):
		time_lbl.text = "%.0fs" % max(0, time_left)
		time_lbl.modulate = Color(1.0, 0.3, 0.3) if time_left <= 5 else Color(0.8, 0.8, 0.8)

func _show_coop_start(name: String, target: int) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()
	_banner = Control.new()
	_banner.position = Vector2(0, 100)
	_banner.size = Vector2(1280, 80)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = COLOR_BG
	_banner.add_child(bg)

	var top_line = ColorRect.new()
	top_line.size = Vector2(1280, 3)
	top_line.color = COLOR_COOP
	_banner.add_child(top_line)

	var lbl = Label.new()
	lbl.text = "🤝 %s 發起全服合作！目標 %d 點！20 秒！" % [name, target]
	lbl.position = Vector2(0, 10)
	lbl.size = Vector2(1280, 40)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 26)
	lbl.modulate = COLOR_COOP
	_banner.add_child(lbl)

	var sub = Label.new()
	sub.text = "每擊破一個目標 +1 點，BOSS +10 點！"
	sub.position = Vector2(0, 50)
	sub.size = Vector2(1280, 25)
	sub.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	sub.add_theme_font_size_override("font_size", 13)
	sub.modulate = Color(0.8, 0.9, 1.0)
	_banner.add_child(sub)

	var bot_line = ColorRect.new()
	bot_line.size = Vector2(1280, 3)
	bot_line.position = Vector2(0, 77)
	bot_line.color = COLOR_COOP
	_banner.add_child(bot_line)

	_banner.position.y = -80
	var tween = _banner.create_tween()
	tween.tween_property(_banner, "position:y", 100.0, 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.tween_interval(3.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(_banner):
			_banner.queue_free()
			_banner = null
	)

func _show_coop_success(name: String, current: int, target: int, boost_mult: float, boost_secs: int) -> void:
	_do_flash(COLOR_SUCCESS, 5)
	ScreenShake.add_trauma(0.7)

	var popup = Control.new()
	popup.position = Vector2(340, 180)
	popup.size = Vector2(600, 220)
	add_child(popup)

	var bg = ColorRect.new()
	bg.size = popup.size
	bg.color = Color(0.0, 0.1, 0.05, 0.92)
	popup.add_child(bg)

	var lbl = Label.new()
	lbl.text = "🤝✨ 全服合作成功！"
	lbl.position = Vector2(0, 15)
	lbl.size = Vector2(600, 50)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 30)
	lbl.modulate = COLOR_SUCCESS
	popup.add_child(lbl)

	var score_lbl = Label.new()
	score_lbl.text = "達成 %d / %d 點！" % [current, target]
	score_lbl.position = Vector2(0, 70)
	score_lbl.size = Vector2(600, 35)
	score_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	score_lbl.add_theme_font_size_override("font_size", 20)
	score_lbl.modulate = COLOR_COOP
	popup.add_child(score_lbl)

	var boost_lbl = Label.new()
	boost_lbl.text = "全服 ×%.0f 加成 %d 秒！" % [boost_mult, boost_secs]
	boost_lbl.position = Vector2(0, 115)
	boost_lbl.size = Vector2(600, 50)
	boost_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	boost_lbl.add_theme_font_size_override("font_size", 26)
	boost_lbl.modulate = Color(1.0, 0.85, 0.0)
	popup.add_child(boost_lbl)

	var hint_lbl = Label.new()
	hint_lbl.text = "依貢獻比例分配獎勵！"
	hint_lbl.position = Vector2(0, 170)
	hint_lbl.size = Vector2(600, 30)
	hint_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	hint_lbl.add_theme_font_size_override("font_size", 13)
	hint_lbl.modulate = Color(0.8, 0.9, 0.8)
	popup.add_child(hint_lbl)

	popup.scale = Vector2.ZERO
	popup.pivot_offset = Vector2(300, 110)
	var tween = popup.create_tween()
	tween.tween_property(popup, "scale", Vector2(1.0, 1.0), 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
	tween.tween_interval(3.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(popup):
			popup.queue_free()
	)

	# 隱藏進度條
	if is_instance_valid(_progress_panel):
		var ptween = _progress_panel.create_tween()
		ptween.tween_property(_progress_panel, "modulate:a", 0.0, 0.5)
		ptween.tween_callback(func():
			if is_instance_valid(_progress_panel):
				_progress_panel.visible = false
				_progress_panel.modulate.a = 1.0
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
			_show_coop_start(name, data.get("target_points", 8))
			_update_progress(0, data.get("target_points", 8), 20.0)
			_do_flash(COLOR_COOP, 3)
			ScreenShake.add_trauma(0.3)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"coop_progress":
			_update_progress(
				data.get("current_points", 0),
				data.get("target_points", 8),
				data.get("time_left", 0.0)
			)
		"coop_success":
			_show_coop_success(
				name,
				data.get("current_points", 0),
				data.get("target_points", 8),
				data.get("boost_mult", 4.0),
				data.get("boost_secs", 8)
			)
		"coop_timeout":
			if is_instance_valid(_progress_panel):
				_progress_panel.visible = false
		"coop_boost_end":
			pass  # 加成結束，靜默處理
