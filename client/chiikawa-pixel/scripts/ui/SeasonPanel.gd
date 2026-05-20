## SeasonPanel.gd — 賽季通行證面板（DAY-072）
## 顯示賽季積分進度條和 10 個等級獎勵
## 位置：TopBar 下方（可折疊）
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 320
const PANEL_HEIGHT := 200
const BTN_SIZE     := 26

# ---- 節點引用 ----
var _pixel_font: Font = null
var _is_open: bool = false
var _toggle_btn: Button = null
var _panel_bg: ColorRect = null
var _progress_bar: ColorRect = null
var _progress_fill: ColorRect = null
var _points_label: Label = null
var _level_label: Label = null
var _level_buttons: Array = []

# ---- 賽季資料 ----
var _season_data: Dictionary = {
	"season_points": 0,
	"current_level": 0,
	"next_level": 1,
	"points_to_next": 100,
	"progress": 0.0,
	"levels": []
}

# ---- 訊號 ----
signal season_level_claimed(level: int)

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_build_toggle_btn()
	_build_panel()
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

## 建立折疊按鈕（TopBar 上）
func _build_toggle_btn() -> void:
	_toggle_btn = Button.new()
	_toggle_btn.text = "🏆"
	_toggle_btn.size = Vector2(32, 24)
	_toggle_btn.position = Vector2(0, 0)
	_toggle_btn.flat = true
	_toggle_btn.tooltip_text = "賽季通行證"
	if _pixel_font:
		_toggle_btn.add_theme_font_override("font", _pixel_font)
		_toggle_btn.add_theme_font_size_override("font_size", 14)
	add_child(_toggle_btn)

## 建立主面板（預設隱藏）
func _build_panel() -> void:
	_panel_bg = ColorRect.new()
	_panel_bg.position = Vector2(-PANEL_WIDTH + 32, 28)
	_panel_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_panel_bg.color = Color(0.05, 0.03, 0.15, 0.92)
	_panel_bg.visible = false
	add_child(_panel_bg)

	# 標題
	var title := Label.new()
	title.position = Vector2(8, 4)
	title.text = "🏆 賽季通行證"
	title.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 12)
	_panel_bg.add_child(title)

	# 積分標籤
	_points_label = Label.new()
	_points_label.position = Vector2(8, 20)
	_points_label.text = "積分：0"
	_points_label.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
	if _pixel_font:
		_points_label.add_theme_font_override("font", _pixel_font)
		_points_label.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_points_label)

	# 等級標籤
	_level_label = Label.new()
	_level_label.position = Vector2(200, 20)
	_level_label.text = "等級：0/10"
	_level_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _pixel_font:
		_level_label.add_theme_font_override("font", _pixel_font)
		_level_label.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_level_label)

	# 進度條背景
	_progress_bar = ColorRect.new()
	_progress_bar.position = Vector2(8, 34)
	_progress_bar.size = Vector2(PANEL_WIDTH - 16, 10)
	_progress_bar.color = Color(0.15, 0.1, 0.3, 0.9)
	_panel_bg.add_child(_progress_bar)

	# 進度條填充
	_progress_fill = ColorRect.new()
	_progress_fill.position = Vector2(8, 34)
	_progress_fill.size = Vector2(0, 10)
	_progress_fill.color = Color(1.0, 0.85, 0.2)
	_panel_bg.add_child(_progress_fill)

	# 10 個等級按鈕（兩行，每行 5 個）
	_build_level_buttons()

## 建立等級按鈕
func _build_level_buttons() -> void:
	for i in range(10):
		var row = i / 5
		var col = i % 5
		var btn_x = 8 + col * (BTN_SIZE + 4)
		var btn_y = 50 + row * (BTN_SIZE + 24)

		# 按鈕背景
		var btn_bg := ColorRect.new()
		btn_bg.position = Vector2(btn_x, btn_y)
		btn_bg.size = Vector2(BTN_SIZE, BTN_SIZE)
		btn_bg.color = Color(0.1, 0.08, 0.25, 0.9)
		btn_bg.name = "LvlBG_%d" % (i + 1)
		_panel_bg.add_child(btn_bg)

		# 等級圖示
		var icon_label := Label.new()
		icon_label.position = Vector2(btn_x + 2, btn_y + 2)
		icon_label.text = "⭐"
		if _pixel_font:
			icon_label.add_theme_font_override("font", _pixel_font)
			icon_label.add_theme_font_size_override("font_size", 14)
		_panel_bg.add_child(icon_label)

		# 等級數字
		var num_label := Label.new()
		num_label.position = Vector2(btn_x, btn_y + BTN_SIZE + 2)
		num_label.text = "Lv%d" % (i + 1)
		num_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.9))
		if _pixel_font:
			num_label.add_theme_font_override("font", _pixel_font)
			num_label.add_theme_font_size_override("font_size", 8)
		_panel_bg.add_child(num_label)

		# 點擊按鈕
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

## 連接訊號
func _connect_signals() -> void:
	_toggle_btn.pressed.connect(_on_toggle_pressed)
	for item in _level_buttons:
		item["btn"].pressed.connect(_on_level_btn_pressed.bind(item["level"]))

	# 連接 GameManager 訊號
	if GameManager.has_signal("season_updated"):
		GameManager.season_updated.connect(_on_season_updated)
	if GameManager.has_signal("season_level_up"):
		GameManager.season_level_up.connect(_on_season_level_up)

func _on_toggle_pressed() -> void:
	_is_open = !_is_open
	_panel_bg.visible = _is_open

func _on_level_btn_pressed(level: int) -> void:
	# 找到等級資料
	var levels = _season_data.get("levels", [])
	for lvl in levels:
		if lvl.get("level") == level:
			if lvl.get("unlocked", false) and not lvl.get("claimed", false):
				# 發送領取請求
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
	# 顯示升級通知
	_show_level_up_notification(data)

## 更新 UI
func _refresh_ui() -> void:
	var points = _season_data.get("season_points", 0)
	var current_level = _season_data.get("current_level", 0)
	var progress = _season_data.get("progress", 0.0)
	var points_to_next = _season_data.get("points_to_next", 100)

	if is_instance_valid(_points_label):
		_points_label.text = "積分：%d（距下級：%d）" % [points, points_to_next]

	if is_instance_valid(_level_label):
		_level_label.text = "等級：%d/10" % current_level

	# 更新進度條
	if is_instance_valid(_progress_fill):
		var max_width = PANEL_WIDTH - 16
		_progress_fill.size.x = max_width * progress

	# 更新等級按鈕狀態
	var levels = _season_data.get("levels", [])
	for item in _level_buttons:
		var level = item["level"]
		var bg = item["bg"]
		var icon = item["icon"]
		if not is_instance_valid(bg):
			continue

		# 找到對應等級資料
		var lvl_data = {}
		for l in levels:
			if l.get("level") == level:
				lvl_data = l
				break

		var claimed = lvl_data.get("claimed", false)
		var unlocked = lvl_data.get("unlocked", false)
		var special_type = lvl_data.get("special_type", "")

		if claimed:
			# 已領取：綠色
			bg.color = Color(0.05, 0.25, 0.05, 0.9)
			if is_instance_valid(icon):
				icon.text = "✅"
		elif unlocked:
			# 可領取：金色閃爍
			bg.color = Color(0.3, 0.25, 0.05, 0.95)
			if is_instance_valid(icon):
				if special_type == "skin":
					icon.text = "🎨"
				elif special_type == "title":
					icon.text = "👑"
				else:
					icon.text = "💰"
		else:
			# 未解鎖：灰色
			bg.color = Color(0.1, 0.08, 0.25, 0.9)
			if is_instance_valid(icon):
				icon.text = "🔒"

## 顯示升級通知
func _show_level_up_notification(data: Dictionary) -> void:
	var level = data.get("level", 0)
	var coin_reward = data.get("coin_reward", 0)
	var special_type = data.get("special_type", "")
	var special_name = data.get("special_name", "")

	var text = "🏆 賽季等級 %d！\n+%d 金幣" % [level, coin_reward]
	if special_type == "skin":
		text += "\n🎨 解鎖：%s" % special_name
	elif special_type == "title":
		text += "\n👑 解鎖：%s" % special_name

	# 建立通知標籤
	var notify := Label.new()
	notify.text = text
	notify.position = Vector2(-100, -60)
	notify.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _pixel_font:
		notify.add_theme_font_override("font", _pixel_font)
		notify.add_theme_font_size_override("font_size", 11)
	add_child(notify)

	# 動畫：彈入 → 停留 → 淡出
	var tween = create_tween()
	tween.tween_property(notify, "scale", Vector2(1.1, 1.1), 0.15)
	tween.tween_property(notify, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(2.0)
	tween.tween_property(notify, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(notify):
			notify.queue_free()
	)
