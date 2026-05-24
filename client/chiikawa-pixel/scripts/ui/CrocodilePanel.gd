## CrocodilePanel.gd
## 撌典?暽寞偌敼琿??菟?蝝舐??Ｘ嚗AY-146嚗?
## 璆剔?靘?嚗iligames.com 2026?iant crocodiles awaken to hunt fish on the fish farm to accumulate big prizes!??
## ? T110 敺孛?潮捧擳擳芋撘?8 蝘?芸??菜捏?桅璅?蝝舐??

extends Control

var pixel_font: Font = null

# 敼琿?銝駁?憿嚗楛蝬嚗?
const COLOR_CROC_GREEN  = Color(0.13, 0.55, 0.13)
const COLOR_CROC_DARK   = Color(0.05, 0.25, 0.05)
const COLOR_CROC_YELLOW = Color(0.8, 0.9, 0.1)
const COLOR_GOLD        = Color(1.0, 0.85, 0.0)
const COLOR_WHITE       = Color(1.0, 1.0, 1.0)

# ?菟??脣漲餈質馱
var _hunt_count: int = 0
var _max_hunts: int = 6
var _total_reward: int = 0
var _banner: Control = null
var _progress_bar: ColorRect = null
var _hunt_label: Label = null

## ??????HUD.gd ?澆嚗?
func setup(font: Font) -> void:
	pixel_font = font
	GameManager.crocodile_hunt.connect(_on_crocodile_hunt)

## ??敼琿??菟?閮
func _on_crocodile_hunt(data: Dictionary) -> void:
	var phase = data.get("phase", "")
	match phase:
		"awaken":
			_show_awaken(data)
		"hunt":
			_show_hunt(data)
		"result":
			_show_result(data)

## 敼琿?閬粹?嚗waken ?挾嚗?
func _show_awaken(data: Dictionary) -> void:
	_hunt_count = 0
	_max_hunts = data.get("max_hunts", 6)
	_total_reward = 0

	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# ?璈怠?嚗楛蝬嚗?
	_banner = Control.new()
	_banner.name = "CrocodileBanner"
	_banner.position = Vector2(0, -60)
	_banner.size = Vector2(1280, 56)
	_banner.z_index = 87
	canvas_layer.add_child(_banner)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.05, 0.2, 0.05, 0.92)
	_banner.add_child(bg)

	# 璅?
	var title_lbl = Label.new()
	title_lbl.text = "?? 撌典?敼琿?閬粹?嚗?憪擳?"
	title_lbl.position = Vector2(0, 4)
	title_lbl.size = Vector2(1280, 28)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", COLOR_CROC_YELLOW)
	if is_instance_valid(pixel_font):
		title_lbl.add_theme_font_override("font", pixel_font)
	_banner.add_child(title_lbl)

	# ?菟??脣漲璅惜
	_hunt_label = Label.new()
	_hunt_label.name = "HuntLabel"
	_hunt_label.text = "?菜捏?脣漲嚗? / %d" % _max_hunts
	_hunt_label.position = Vector2(0, 32)
	_hunt_label.size = Vector2(1280, 20)
	_hunt_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_hunt_label.add_theme_font_size_override("font_size", 12)
	_hunt_label.add_theme_color_override("font_color", COLOR_WHITE)
	if is_instance_valid(pixel_font):
		_hunt_label.add_theme_font_override("font", pixel_font)
	_banner.add_child(_hunt_label)

	# ?脣漲璇???
	var bar_bg = ColorRect.new()
	bar_bg.position = Vector2(440, 50)
	bar_bg.size = Vector2(400, 6)
	bar_bg.color = Color(0.1, 0.1, 0.1, 0.8)
	_banner.add_child(bar_bg)

	# ?脣漲璇‵??
	_progress_bar = ColorRect.new()
	_progress_bar.name = "ProgressBar"
	_progress_bar.position = Vector2(440, 50)
	_progress_bar.size = Vector2(0, 6)
	_progress_bar.color = COLOR_CROC_GREEN
	_banner.add_child(_progress_bar)

	# 皛?
	var tween = _banner.create_tween()
	tween.tween_property(_banner, "position:y", 0.0, 0.3).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

	# ?刻撟??脤???
	var flash = ColorRect.new()
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(COLOR_CROC_GREEN.r, COLOR_CROC_GREEN.g, COLOR_CROC_GREEN.b, 0.0)
	flash.z_index = 86
	canvas_layer.add_child(flash)
	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", 0.25, 0.08)
	flash_tween.tween_property(flash, "color:a", 0.0, 0.3)
	flash_tween.tween_callback(flash.queue_free)

## 敼琿??菜捏嚗unt ?挾嚗?
func _show_hunt(data: Dictionary) -> void:
	_hunt_count += 1
	var hunt_reward = data.get("hunt_reward", 0)
	_total_reward += hunt_reward

	# ?湔?脣漲璇?
	if is_instance_valid(_progress_bar) and _max_hunts > 0:
		var ratio = float(_hunt_count) / float(_max_hunts)
		var tween = _progress_bar.create_tween()
		tween.tween_property(_progress_bar, "size:x", 400.0 * ratio, 0.2).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)

	# ?湔?脣漲璅惜
	if is_instance_valid(_hunt_label):
		_hunt_label.text = "?菜捏?脣漲嚗?d / %d  蝝舐?嚗?d" % [_hunt_count, _max_hunts, _total_reward]

	# ?菜捏瘚桀???
	var canvas_layer = get_parent()
	if is_instance_valid(canvas_layer) and hunt_reward > 0:
		var float_lbl = Label.new()
		float_lbl.text = "+%d" % hunt_reward
		float_lbl.position = Vector2(
			randf_range(400, 880),
			randf_range(200, 500)
		)
		float_lbl.size = Vector2(100, 24)
		float_lbl.add_theme_font_size_override("font_size", 14)
		float_lbl.add_theme_color_override("font_color", COLOR_CROC_YELLOW)
		if is_instance_valid(pixel_font):
			float_lbl.add_theme_font_override("font", pixel_font)
		float_lbl.z_index = 88
		canvas_layer.add_child(float_lbl)

		var tween = float_lbl.create_tween()
		tween.tween_property(float_lbl, "position:y", float_lbl.position.y - 40, 0.6).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)
		tween.parallel().tween_property(float_lbl, "modulate:a", 0.0, 0.6)
		tween.tween_callback(float_lbl.queue_free)

## 憿舐內蝯?嚗esult ?挾嚗?
func _show_result(data: Dictionary) -> void:
	var total_reward = data.get("total_reward", 0)
	var killer_id = data.get("killer_id", "")
	var killer_name = data.get("killer_name", "")
	var hunted_targets = data.get("hunted_targets", [])
	var is_self = killer_id == NetworkManager.get_player_id()

	# ??璈怠?
	if is_instance_valid(_banner):
		var tween = _banner.create_tween()
		tween.tween_property(_banner, "position:y", -60.0, 0.3).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
		tween.tween_callback(_banner.queue_free)
	_banner = null

	if total_reward <= 0:
		return

	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# 蝯?敶?嚗?湔??伐?
	var panel = Control.new()
	panel.name = "CrocodileResult"
	panel.position = Vector2(1280, 80)
	panel.size = Vector2(280, 160)
	panel.z_index = 88
	canvas_layer.add_child(panel)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.03, 0.12, 0.03, 0.92)
	panel.add_child(bg)

	# 蝬???
	var top_border = ColorRect.new()
	top_border.size = Vector2(280, 2)
	top_border.color = COLOR_CROC_GREEN
	panel.add_child(top_border)

	# 璅?
	var title_lbl = Label.new()
	title_lbl.text = "?? 敼琿??菟?摰?嚗? if is_self else "?? %s ?捧擳?" % killer_name"
	title_lbl.position = Vector2(8, 8)
	title_lbl.size = Vector2(264, 22)
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", COLOR_CROC_YELLOW)
	if is_instance_valid(pixel_font):
		title_lbl.add_theme_font_override("font", pixel_font)
	panel.add_child(title_lbl)

	# ?菜捏??
	var hunt_lbl = Label.new()
	hunt_lbl.text = "?菜捏?格?嚗?d ?? % len(hunted_targets)"
	hunt_lbl.position = Vector2(8, 34)
	hunt_lbl.size = Vector2(264, 18)
	hunt_lbl.add_theme_font_size_override("font_size", 11)
	hunt_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	if is_instance_valid(pixel_font):
		hunt_lbl.add_theme_font_override("font", pixel_font)
	panel.add_child(hunt_lbl)

	# ?菜捏?”嚗?憭＊蝷?3 ??
	var display_count = min(len(hunted_targets), 3)
	for i in display_count:
		var entry = hunted_targets[i]
		var item_lbl = Label.new()
		item_lbl.text = "  ?? ?%.0f ??+%d" % [entry.get("multiplier", 0), entry.get("reward", 0)]
		item_lbl.position = Vector2(8, 54 + i * 16)
		item_lbl.size = Vector2(264, 16)
		item_lbl.add_theme_font_size_override("font_size", 10)
		item_lbl.add_theme_color_override("font_color", COLOR_CROC_YELLOW)
		if is_instance_valid(pixel_font):
			item_lbl.add_theme_font_override("font", pixel_font)
		panel.add_child(item_lbl)

	# 蝮賜???
	var reward_lbl = Label.new()
	reward_lbl.text = "?? +%d" % total_reward
	reward_lbl.position = Vector2(8, 112)
	reward_lbl.size = Vector2(264, 32)
	reward_lbl.add_theme_font_size_override("font_size", 22)
	reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	if is_instance_valid(pixel_font):
		reward_lbl.add_theme_font_override("font", pixel_font)
	panel.add_child(reward_lbl)

	# ?芸楛閫貊??蝬??
	if is_self:
		var flash = ColorRect.new()
		flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		flash.color = Color(COLOR_CROC_GREEN.r, COLOR_CROC_GREEN.g, COLOR_CROC_GREEN.b, 0.0)
		flash.z_index = 1
		panel.add_child(flash)
		var flash_tween = flash.create_tween()
		flash_tween.tween_property(flash, "color:a", 0.4, 0.1)
		flash_tween.tween_property(flash, "color:a", 0.0, 0.3)
		flash_tween.tween_callback(flash.queue_free)

		if AudioManager != null:
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)

	# 皛?
	var tween = panel.create_tween()
	tween.tween_property(panel, "position:x", 990.0, 0.35).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.tween_interval(3.5)
	tween.tween_property(panel, "position:x", 1280.0, 0.3).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
	tween.tween_callback(panel.queue_free)
