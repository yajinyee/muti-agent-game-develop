## AnnouncementPanel.gd
## ?冽??砍?蝟餌絞?Ｘ嚗AY-097嚗?
## 憿舐內?之鈭辣?嚗ackpot 銝剔??之?OSS ?捏?摰嗅??亦?

extends Control

var pixel_font: Font = null

# ?砍?雿?嚗???＊蝷箏云憭?
var _queue: Array = []
var _is_showing: bool = false
const MAX_QUEUE = 5

# ?砍?憿舐內雿蔭嚗銝?嚗?
const PANEL_WIDTH = 380
const PANEL_HEIGHT = 60
const PANEL_X = 900  # 1280 - 380
const PANEL_Y_START = 680  # 720 - 40

## ??????HUD.gd ?澆嚗?
func setup(font: Font) -> void:
	pixel_font = font
	GameManager.announcement_received.connect(_on_announcement_received)

## ?嗅?砍?
func _on_announcement_received(data: Dictionary) -> void:
	var priority = data.get("priority", 1)

	# 雿???拙振?/?ａ?嚗?函摰嗅??＊蝷?
	if priority <= 1:
		# 蝪∪?嚗??芸??砍??湔?雿?雿?頞? 2 ??
		var low_count = 0
		for item in _queue:
			if item.get("priority", 1) <= 1:
				low_count += 1
		if low_count >= 2:
			return

	# 雿?皛蹂?撠曹?璉??芸?
	if _queue.size() >= MAX_QUEUE:
		# ?暹?雿??蝘駁
		var min_priority = 999
		var min_idx = -1
		for i in _queue.size():
			if _queue[i].get("priority", 1) < min_priority:
				min_priority = _queue[i].get("priority", 1)
				min_idx = i
		if min_idx >= 0 and min_priority < priority:
			_queue.remove_at(min_idx)
		else:
			return

	_queue.append(data)

	if not _is_showing:
		_show_next()

## 憿舐內銝????
func _show_next() -> void:
	if _queue.is_empty():
		_is_showing = false
		return

	_is_showing = true
	var data = _queue.pop_front()
	_display_announcement(data)

## 憿舐內?砍?
func _display_announcement(data: Dictionary) -> void:
	var priority = data.get("priority", 1)
	var title = data.get("title", "?砍?")
	var message = data.get("message", "")
	var icon = data.get("icon", "?")
	var color_hex = data.get("color", "#FFFFFF")
	var duration = data.get("duration", 3000)
	var event_type = data.get("event_type", "")

	# 閫??憿
	var ann_color = Color.WHITE
	if color_hex.begins_with("#"):
		ann_color = Color(color_hex)

	# 撱箇??砍?蝭暺?
	var ann_node = Control.new()
	ann_node.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	ann_node.position = Vector2(PANEL_X + PANEL_WIDTH, PANEL_Y_START)  # 敺?湔???
	ann_node.z_index = 80
	add_child(ann_node)

	# ?嚗??芸?蝝矽?湧?摨血?憿嚗?
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	var bg_alpha = 0.75 + (priority - 1) * 0.05
	bg.color = Color(0.05, 0.05, 0.15, bg_alpha)
	ann_node.add_child(bg)

	# 撌血敶抵??嚗??芸?蝝?
	var side_bar = ColorRect.new()
	side_bar.size = Vector2(4, PANEL_HEIGHT)
	side_bar.position = Vector2(0, 0)
	side_bar.color = ann_color
	ann_node.add_child(side_bar)

	# ?內
	var icon_lbl = Label.new()
	icon_lbl.text = icon
	icon_lbl.position = Vector2(10, 8)
	icon_lbl.size = Vector2(32, 44)
	icon_lbl.add_theme_font_size_override("font_size", 28)
	ann_node.add_child(icon_lbl)

	# 璅?
	var title_lbl = Label.new()
	title_lbl.text = title
	title_lbl.position = Vector2(48, 6)
	title_lbl.size = Vector2(PANEL_WIDTH - 56, 22)
	title_lbl.add_theme_font_size_override("font_size", 12)
	title_lbl.add_theme_color_override("font_color", ann_color)
	if is_instance_valid(pixel_font):
		title_lbl.add_theme_font_override("font", pixel_font)
	ann_node.add_child(title_lbl)

	# 閮
	var msg_lbl = Label.new()
	msg_lbl.text = message
	msg_lbl.position = Vector2(48, 30)
	msg_lbl.size = Vector2(PANEL_WIDTH - 56, 24)
	msg_lbl.add_theme_font_size_override("font_size", 10)
	msg_lbl.add_theme_color_override("font_color", Color(0.85, 0.85, 0.85))
	if is_instance_valid(pixel_font):
		msg_lbl.add_theme_font_override("font", pixel_font)
	ann_node.add_child(msg_lbl)

	# Grand Jackpot ?寞???嚗??脤???獢?
	if event_type == "grand_jackpot":
		var glow = ColorRect.new()
		glow.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		glow.color = Color(1.0, 0.85, 0.0, 0.0)
		ann_node.add_child(glow)
		var glow_tween = glow.create_tween().set_loops(3)
		glow_tween.tween_property(glow, "color:a", 0.3, 0.2)
		glow_tween.tween_property(glow, "color:a", 0.0, 0.2)

	# ?嚗??喳皛
	var duration_sec = duration / 1000.0
	var tween = ann_node.create_tween()

	# 皛
	tween.tween_property(ann_node, "position:x", float(PANEL_X), 0.25).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

	# ??
	tween.tween_interval(duration_sec - 0.5)

	# 瘛∪
	tween.tween_property(ann_node, "modulate:a", 0.0, 0.25)
	tween.tween_callback(ann_node.queue_free)
	tween.tween_callback(_on_announcement_done)

## ?砍?憿舐內摰嚗＊蝷箔?銝??
func _on_announcement_done() -> void:
	# ?剜??敺＊蝷箔?銝??
	var timer = get_tree().create_timer(0.15)
	timer.timeout.connect(_show_next)
