## SeasonPanel.gd ??鞈賢迤??霅?選?DAY-072嚗?
## 憿舐內鞈賢迤蝛??脣漲璇? 10 ??蝝???
## 雿蔭嚗opBar 銝嚗??嚗?
extends Node2D

# ---- 撣豢 ----
const PANEL_WIDTH  := 320
const PANEL_HEIGHT := 200
const BTN_SIZE     := 26

# ---- 蝭暺???----
var _pixel_font: Font = null
var _is_open: bool = false
var _toggle_btn: Button = null
var _panel_bg: ColorRect = null
var _progress_bar: ColorRect = null
var _progress_fill: ColorRect = null
var _points_label: Label = null
var _level_label: Label = null
var _level_buttons: Array = []

# ---- 鞈賢迤鞈? ----
var _season_data: Dictionary = {
	"season_points": 0,
	"current_level": 0,
	"next_level": 1,
	"points_to_next": 100,
	"progress": 0.0,
	"levels": []
}

# ---- 閮? ----
signal season_level_claimed(level: int)

# ---- ????----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_build_toggle_btn()
	_build_panel()
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

## 撱箇?????嚗opBar 銝?
func _build_toggle_btn() -> void:
	_toggle_btn = Button.new()
	_toggle_btn.text = "??"
	_toggle_btn.size = Vector2(32, 24)
	_toggle_btn.position = Vector2(0, 0)
	_toggle_btn.flat = true
	_toggle_btn.tooltip_text = "鞈賢迤??霅?"
	if _pixel_font:
		_toggle_btn.add_theme_font_override("font", _pixel_font)
		_toggle_btn.add_theme_font_size_override("font_size", 14)
	add_child(_toggle_btn)

## 撱箇?銝駁?選??身?梯?嚗?
func _build_panel() -> void:
	_panel_bg = ColorRect.new()
	_panel_bg.position = Vector2(-PANEL_WIDTH + 32, 28)
	_panel_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_panel_bg.color = Color(0.05, 0.03, 0.15, 0.92)
	_panel_bg.visible = false
	add_child(_panel_bg)

	# 璅?
	var title := Label.new()
	title.position = Vector2(8, 4)
	title.text = "?? 鞈賢迤??霅?"
	title.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 12)
	_panel_bg.add_child(title)

	# 蝛?璅惜
	_points_label = Label.new()
	_points_label.position = Vector2(8, 20)
	_points_label.text = "蝛?嚗?"
	_points_label.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
	if _pixel_font:
		_points_label.add_theme_font_override("font", _pixel_font)
		_points_label.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_points_label)

	# 蝑?璅惜
	_level_label = Label.new()
	_level_label.position = Vector2(200, 20)
	_level_label.text = "蝑?嚗?/10"
	_level_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _pixel_font:
		_level_label.add_theme_font_override("font", _pixel_font)
		_level_label.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_level_label)

	# ?脣漲璇???
	_progress_bar = ColorRect.new()
	_progress_bar.position = Vector2(8, 34)
	_progress_bar.size = Vector2(PANEL_WIDTH - 16, 10)
	_progress_bar.color = Color(0.15, 0.1, 0.3, 0.9)
	_panel_bg.add_child(_progress_bar)

	# ?脣漲璇‵??
	_progress_fill = ColorRect.new()
	_progress_fill.position = Vector2(8, 34)
	_progress_fill.size = Vector2(0, 10)
	_progress_fill.color = Color(1.0, 0.85, 0.2)
	_panel_bg.add_child(_progress_fill)

	# 10 ??蝝????抵?嚗?銵?5 ??
	_build_level_buttons()

## 撱箇?蝑???
func _build_level_buttons() -> void:
	for i in range(10):
		var row = i / 5
		var col = i % 5
		var btn_x = 8 + col * (BTN_SIZE + 4)
		var btn_y = 50 + row * (BTN_SIZE + 24)

		# ???
		var btn_bg := ColorRect.new()
		btn_bg.position = Vector2(btn_x, btn_y)
		btn_bg.size = Vector2(BTN_SIZE, BTN_SIZE)
		btn_bg.color = Color(0.1, 0.08, 0.25, 0.9)
		btn_bg.name = "LvlBG_%d" % (i + 1)
		_panel_bg.add_child(btn_bg)

		# 蝑??內
		var icon_label := Label.new()
		icon_label.position = Vector2(btn_x + 2, btn_y + 2)
		icon_label.text = "潃?"
		if _pixel_font:
			icon_label.add_theme_font_override("font", _pixel_font)
			icon_label.add_theme_font_size_override("font_size", 14)
		_panel_bg.add_child(icon_label)

		# 蝑??詨?
		var num_label := Label.new()
		num_label.position = Vector2(btn_x, btn_y + BTN_SIZE + 2)
		num_label.text = "Lv%d" % (i + 1)
		num_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.9))
		if _pixel_font:
			num_label.add_theme_font_override("font", _pixel_font)
			num_label.add_theme_font_size_override("font_size", 8)
		_panel_bg.add_child(num_label)

		# 暺???
		var btn := Button.new()
		btn.position = Vector2(btn_x, btn_y)
		btn.size = Vector2(BTN_SIZE, BTN_SIZE)
		btn.flat = true
		btn.text = ""
		btn.set_meta("level", i + 1)
		_panel_bg.add_child(btn)

		_level_buttons.append({
			"btn": btn,
			"bg": btn_bg,
			"icon": icon_label,
			"num": num_label,
			"level": i + 1
		})

## ??閮?
func _connect_signals() -> void:
	_toggle_btn.pressed.connect(_on_toggle_pressed)
	for item in _level_buttons:
		item["btn"].pressed.connect(_on_level_btn_pressed.bind(item["level"]))

	# ?? GameManager 閮?
	if GameManager.has_signal("season_updated"):
		GameManager.season_updated.connect(_on_season_updated)
	if GameManager.has_signal("season_level_up"):
		GameManager.season_level_up.connect(_on_season_level_up)

func _on_toggle_pressed() -> void:
	_is_open = !_is_open
	_panel_bg.visible = _is_open

func _on_level_btn_pressed(level: int) -> void:
	# ?曉蝑?鞈?
	var levels = _season_data.get("levels", [])
	for lvl in levels:
		if lvl.get("level") == level:
			if lvl.get("unlocked", false) and not lvl.get("claimed", false):
				# ?潮???瘙?
				NetworkManager.send_message({
					"type": "claim_season_level",
					"payload": {"level": level}
				})
				emit_signal("season_level_claimed", level)
			return

func _on_season_updated(data: Dictionary) -> void:
	_season_data = data
	_refresh_ui()

func _on_season_level_up(data: Dictionary) -> void:
	# 憿舐內???
	_show_level_up_notification(data)

## ?湔 UI
func _refresh_ui() -> void:
	var points = _season_data.get("season_points", 0)
	var current_level = _season_data.get("current_level", 0)
	var progress = _season_data.get("progress", 0.0)
	var points_to_next = _season_data.get("points_to_next", 100)

	if is_instance_valid(_points_label):
		_points_label.text = "蝛?嚗?d嚗?銝?嚗?d嚗? % [points, points_to_next]"

	if is_instance_valid(_level_label):
		_level_label.text = "蝑?嚗?d/10" % current_level

	# ?湔?脣漲璇?
	if is_instance_valid(_progress_fill):
		var max_width = PANEL_WIDTH - 16
		_progress_fill.size.x = max_width * progress

	# ?湔蝑??????
	var levels = _season_data.get("levels", [])
	for item in _level_buttons:
		var level = item["level"]
		var bg = item["bg"]
		var icon = item["icon"]
		if not is_instance_valid(bg):
			continue

		# ?曉撠?蝑?鞈?
		var lvl_data = {}
		for l in levels:
			if l.get("level") == level:
				lvl_data = l
				break

		var claimed = lvl_data.get("claimed", false)
		var unlocked = lvl_data.get("unlocked", false)
		var special_type = lvl_data.get("special_type", "")

		if claimed:
			# 撌脤???蝬
			bg.color = Color(0.05, 0.25, 0.05, 0.9)
			if is_instance_valid(icon):
				icon.text = "??"
		elif unlocked:
			# ?舫??????
			bg.color = Color(0.3, 0.25, 0.05, 0.95)
			if is_instance_valid(icon):
				if special_type == "skin":
					icon.text = "?"
				elif special_type == "title":
					icon.text = "??"
				else:
					icon.text = "?"
		else:
			# ?芾圾???啗
			bg.color = Color(0.1, 0.08, 0.25, 0.9)
			if is_instance_valid(icon):
				icon.text = "??"

## 憿舐內???
func _show_level_up_notification(data: Dictionary) -> void:
	var level = data.get("level", 0)
	var coin_reward = data.get("coin_reward", 0)
	var special_type = data.get("special_type", "")
	var special_name = data.get("special_name", "")

	var text = "?? 鞈賢迤蝑? %d嚗n+%d ?馳" % [level, coin_reward]
	if special_type == "skin":
		text += "\n? 閫??嚗?s" % special_name
	elif special_type == "title":
		text += "\n?? 閫??嚗?s" % special_name

	# 撱箇??璅惜
	var notify := Label.new()
	notify.text = text
	notify.position = Vector2(-100, -60)
	notify.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _pixel_font:
		notify.add_theme_font_override("font", _pixel_font)
		notify.add_theme_font_size_override("font_size", 11)
	add_child(notify)

	# ?嚗??????? ??瘛∪
	var tween = create_tween()
	tween.tween_property(notify, "scale", Vector2(1.1, 1.1), 0.15)
	tween.tween_property(notify, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(2.0)
	tween.tween_property(notify, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(notify):
			notify.queue_free()
	)
