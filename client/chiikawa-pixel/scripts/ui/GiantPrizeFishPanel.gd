## GiantPrizeFishPanel.gd ??憭Ｗ劂撌典??擳?選?DAY-147嚗?
## 璆剔?靘?嚗iligames.com 2026?he dreamy Giant Prize Fish lets you easily win great prizes,
## with the chance for 5x multipliers??
## 閫貊?拙振??10 蝘????渡????5嚗?誨?剖丐撟餅芋撘?憪?蝯?
## 閬死閮剛?嚗?蝝丐撟颱蜓憿??璈怠?皛嚗閮?嚗? ??憿舐內嚗?Ｗ?蝎???嚗???蝒?
extends Control

# ---- 撣豢 ----
const PANEL_COLOR_BG    := Color(0.95, 0.4, 0.7, 0.92)   # 蝎?銝駁??
const PANEL_COLOR_GOLD  := Color(1.0, 0.85, 0.0, 1.0)    # ???
const PANEL_COLOR_WHITE := Color(1.0, 1.0, 1.0, 1.0)
const PANEL_COLOR_PINK  := Color(1.0, 0.4, 0.7, 1.0)     # 蝎???
const PANEL_COLOR_DREAM := Color(0.98, 0.75, 0.9, 1.0)   # 憭Ｗ劂瘛∠?

# ---- ???----
var _pixel_font: Font = null
var _is_active: bool = false
var _is_my_session: bool = false   # ?臬?航撌梯孛?潛?
var _duration: float = 10.0
var _elapsed: float = 0.0
var _mult_bonus: float = 5.0
var _killer_name: String = ""

# ---- ??蝭暺?----
var _banner: Control = null
var _timer_label: Label = null
var _progress_bar: ColorRect = null

## setup ????HUD.gd ?澆嚗? GameManager 閮?
func setup(font: Font) -> void:
	_pixel_font = font
	GameManager.giant_prize_fish.connect(_on_giant_prize_fish)

## _on_giant_prize_fish ???? Server 撱??丐撟餌??菟?鈭辣
func _on_giant_prize_fish(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var killer_id: String = data.get("killer_id", "")
	_killer_name = data.get("killer_name", "?拙振")
	_mult_bonus = data.get("mult_bonus", 5.0)
	_duration = float(data.get("duration", 10))
	_is_my_session = (killer_id == NetworkManager.get_player_id())

	match phase:
		"activate":
			_on_activate(data)
		"end":
			_on_end(data)

func _process(delta: float) -> void:
	if not _is_active:
		return

	_elapsed += delta
	var remaining = max(0.0, _duration - _elapsed)
	var pct = remaining / _duration

	# ?湔?閮?
	if is_instance_valid(_timer_label):
		_timer_label.text = "??%.1fs" % remaining

	# ?湔?脣漲璇?
	if is_instance_valid(_progress_bar):
		_progress_bar.size.x = 180.0 * pct

	# ?敺?3 蝘?蝝??
	if is_instance_valid(_timer_label) and remaining <= 3.0:
		var blink = sin(_elapsed * 8.0) > 0.0
		_timer_label.add_theme_color_override("font_color",
			Color(1.0, 0.2, 0.2) if blink else PANEL_COLOR_PINK)

	# ???堆??皜?
	if _elapsed >= _duration + 0.5:
		_is_active = false

func _on_activate(data: Dictionary) -> void:
	_is_active = true
	_elapsed = 0.0

	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# ?璈怠?嚗?蝝丐撟颱蜓憿?
	_banner = Control.new()
	_banner.name = "GiantPrizeFishBanner"
	_banner.position = Vector2(0, -60)
	_banner.size = Vector2(1280, 56)
	_banner.z_index = 88
	canvas_layer.add_child(_banner)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.6, 0.1, 0.35, 0.92)
	_banner.add_child(bg)

	# 璅?
	var title_lbl = Label.new()
	if _is_my_session:
		title_lbl.text = "??憭Ｗ劂?璅∪?嚗?.0f ??嚗? % _mult_bonus"
	else:
		title_lbl.text = "??%s 閫貊憭Ｗ劂?嚗?.0f ??嚗? % [_killer_name, _mult_bonus]"
	title_lbl.position = Vector2(0, 4)
	title_lbl.size = Vector2(1280, 28)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	if is_instance_valid(_pixel_font):
		title_lbl.add_theme_font_override("font", _pixel_font)
	_banner.add_child(title_lbl)

	# ?閮?璅惜
	_timer_label = Label.new()
	_timer_label.name = "TimerLabel"
	_timer_label.text = "??%.1fs" % _duration
	_timer_label.position = Vector2(0, 32)
	_timer_label.size = Vector2(1280, 20)
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_timer_label.add_theme_font_size_override("font_size", 12)
	_timer_label.add_theme_color_override("font_color", PANEL_COLOR_PINK)
	if is_instance_valid(_pixel_font):
		_timer_label.add_theme_font_override("font", _pixel_font)
	_banner.add_child(_timer_label)

	# ?脣漲璇???
	var bar_bg = ColorRect.new()
	bar_bg.position = Vector2(550, 50)
	bar_bg.size = Vector2(180, 6)
	bar_bg.color = Color(0.1, 0.1, 0.1, 0.8)
	_banner.add_child(bar_bg)

	# ?脣漲璇‵??
	_progress_bar = ColorRect.new()
	_progress_bar.name = "ProgressBar"
	_progress_bar.position = Vector2(550, 50)
	_progress_bar.size = Vector2(180, 6)
	_progress_bar.color = PANEL_COLOR_PINK
	_banner.add_child(_progress_bar)

	# 璈怠?皛?
	var tween = _banner.create_tween()
	tween.tween_property(_banner, "position:y", 0.0, 0.3).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

	# ?刻撟?蝝????芸楛閫貊?撘瑞?嚗?
	var flash_alpha = 0.4 if _is_my_session else 0.2
	var flash = ColorRect.new()
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(1.0, 0.4, 0.7, 0.0)
	flash.z_index = 87
	canvas_layer.add_child(flash)
	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", flash_alpha, 0.1)
	flash_tween.tween_property(flash, "color:a", 0.0, 0.4)
	flash_tween.tween_callback(flash.queue_free)

	# ?芸楛閫貊????憭Ｗ劂??蝎?
	if _is_my_session:
		_spawn_dream_particles(canvas_layer)

func _on_end(data: Dictionary) -> void:
	_is_active = false
	var total_reward: int = data.get("total_reward", 0)
	var kill_count: int = data.get("kill_count", 0)

	# 璈怠?皛
	if is_instance_valid(_banner):
		var tween = _banner.create_tween()
		tween.tween_property(_banner, "position:y", -60.0, 0.3).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
		tween.tween_callback(_banner.queue_free)
	_banner = null
	_timer_label = null
	_progress_bar = null

	# ?芣??芸楛閫貊銝???＊蝷箇???蝒?
	if not _is_my_session or total_reward <= 0:
		return

	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# 蝯?敶?嚗?湔??伐?
	var panel = Control.new()
	panel.name = "GiantPrizeFishResult"
	panel.position = Vector2(1280, 80)
	panel.size = Vector2(280, 140)
	panel.z_index = 89
	canvas_layer.add_child(panel)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.1, 0.03, 0.08, 0.95)
	panel.add_child(bg)

	# 蝎?撌血??
	var left_border = ColorRect.new()
	left_border.size = Vector2(4, 140)
	left_border.color = PANEL_COLOR_PINK
	panel.add_child(left_border)

	# 璅?
	var title_lbl = Label.new()
	title_lbl.text = "??憭Ｗ劂璅∪?蝯?嚗?"
	title_lbl.position = Vector2(12, 8)
	title_lbl.size = Vector2(260, 22)
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	if is_instance_valid(_pixel_font):
		title_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(title_lbl)

	# ???
	var kill_lbl = Label.new()
	kill_lbl.text = "??格?嚗?d ?? % kill_count"
	kill_lbl.position = Vector2(12, 34)
	kill_lbl.size = Vector2(260, 18)
	kill_lbl.add_theme_font_size_override("font_size", 12)
	kill_lbl.add_theme_color_override("font_color", PANEL_COLOR_WHITE)
	if is_instance_valid(_pixel_font):
		kill_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(kill_lbl)

	# ??隤芣?
	var mult_lbl = Label.new()
	mult_lbl.text = "?%.0f ??撌脣??? % _mult_bonus"
	mult_lbl.position = Vector2(12, 56)
	mult_lbl.size = Vector2(260, 18)
	mult_lbl.add_theme_font_size_override("font_size", 12)
	mult_lbl.add_theme_color_override("font_color", PANEL_COLOR_PINK)
	if is_instance_valid(_pixel_font):
		mult_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(mult_lbl)

	# 蝮賜???
	var reward_lbl = Label.new()
	reward_lbl.text = "?? +%d" % total_reward
	reward_lbl.position = Vector2(12, 82)
	reward_lbl.size = Vector2(260, 40)
	reward_lbl.add_theme_font_size_override("font_size", 26)
	reward_lbl.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	if is_instance_valid(_pixel_font):
		reward_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(reward_lbl)

	# 蝎???
	var flash = ColorRect.new()
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(1.0, 0.4, 0.7, 0.0)
	flash.z_index = 1
	panel.add_child(flash)
	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", 0.5, 0.1)
	flash_tween.tween_property(flash, "color:a", 0.0, 0.4)
	flash_tween.tween_callback(flash.queue_free)

	# 皛?
	var tween = panel.create_tween()
	tween.tween_property(panel, "position:x", 990.0, 0.35).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.tween_interval(3.5)
	tween.tween_property(panel, "position:x", 1280.0, 0.3).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
	tween.tween_callback(panel.queue_free)

## _spawn_dream_particles ????憭Ｗ劂??蝎?嚗筑??摮芋?穿?
func _spawn_dream_particles(canvas_layer: Node) -> void:
	var stars = ["??, "潃?, "?", "??", "??]"
	for i in range(5):
		var star_label = Label.new()
		star_label.text = stars[i % stars.size()]
		star_label.add_theme_font_size_override("font_size", 24 + randi() % 16)
		star_label.position = Vector2(
			randf_range(100, 1180),
			randf_range(100, 620)
		)
		star_label.z_index = 90
		canvas_layer.add_child(star_label)

		var star_tween = star_label.create_tween()
		star_tween.tween_property(star_label, "position:y", star_label.position.y - 80, 1.2)
		star_tween.parallel().tween_property(star_label, "modulate:a", 0.0, 1.2)
		star_tween.tween_callback(star_label.queue_free)

# ---- 撣豢 ----
