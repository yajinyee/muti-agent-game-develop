## RareCatchPanel.gd — 稀有連擊累積倍率面板（DAY-126）
## 業界依據：fishingfortune.app（2026-05-21）確認「multiplier cascade system」
## 連續在 90 秒內擊破稀有目標，倍率從 2x 累積到最高 15x
extends Control

# ---- 常數 ----
const PANEL_WIDTH  := 180
const PANEL_HEIGHT := 80
const FADE_DURATION := 0.3
const RESET_HIDE_DELAY := 2.0

# ---- 節點引用 ----
var _bg: ColorRect
var _left_bar: ColorRect
var _icon_label: Label
var _level_label: Label
var _mult_label: Label
var _timer_label: Label
var _count_dots: Array = []

# ---- 狀態 ----
var _is_active := false
var _current_count := 0
var _seconds_left := 0
var _hide_timer: SceneTreeTimer = null
var _tick_timer: SceneTreeTimer = null

func _ready() -> void:
	_build_ui()
	visible = false
	set_process(false)

func _build_ui() -> void:
	custom_minimum_size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)

	# 背景
	_bg = ColorRect.new()
	_bg.color = Color(0.05, 0.05, 0.15, 0.88)
	_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	add_child(_bg)

	# 左側色條（等級顏色）
	_left_bar = ColorRect.new()
	_left_bar.color = Color(0.0, 0.75, 1.0)
	_left_bar.size = Vector2(4, PANEL_HEIGHT)
	_left_bar.position = Vector2(0, 0)
	add_child(_left_bar)

	# 圖示
	_icon_label = Label.new()
	_icon_label.text = "💎"
	_icon_label.position = Vector2(10, 8)
	_icon_label.add_theme_font_size_override("font_size", 22)
	add_child(_icon_label)

	# 等級名稱
	_level_label = Label.new()
	_level_label.text = "稀有連擊"
	_level_label.position = Vector2(40, 8)
	_level_label.add_theme_font_size_override("font_size", 12)
	_level_label.add_theme_color_override("font_color", Color(0.0, 0.75, 1.0))
	add_child(_level_label)

	# 倍率顯示
	_mult_label = Label.new()
	_mult_label.text = "×2.0"
	_mult_label.position = Vector2(40, 28)
	_mult_label.add_theme_font_size_override("font_size", 20)
	_mult_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.5))
	add_child(_mult_label)

	# 倒數計時
	_timer_label = Label.new()
	_timer_label.text = "90s"
	_timer_label.position = Vector2(PANEL_WIDTH - 36, 8)
	_timer_label.add_theme_font_size_override("font_size", 11)
	_timer_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	add_child(_timer_label)

	# 連擊點（最多 5 個）
	for i in range(5):
		var dot := ColorRect.new()
		dot.size = Vector2(10, 10)
		dot.position = Vector2(40 + i * 14, 56)
		dot.color = Color(0.3, 0.3, 0.3, 0.8)
		add_child(dot)
		_count_dots.append(dot)

# ---- 公開 API ----

## 稀有連擊更新（個人）
func on_rare_catch_update(data: Dictionary) -> void:
	_is_active = true
	_current_count = data.get("count", 1)
	_seconds_left = data.get("seconds_left", 90)
	var mult_boost: float = data.get("mult_boost", 2.0)
	var level_name: String = data.get("level_name", "稀有連擊")
	var icon: String = data.get("icon", "💎")
	var color_str: String = data.get("color", "#00BFFF")
	var is_level_up: bool = data.get("is_level_up", false)

	var color := Color(color_str)

	# 更新 UI
	_left_bar.color = color
	_icon_label.text = icon
	_level_label.text = level_name
	_level_label.add_theme_color_override("font_color", color)

	var mult_str := "%.0f" % mult_boost
	_mult_label.text = "×" + mult_str
	_timer_label.text = str(_seconds_left) + "s"

	# 更新連擊點
	for i in range(5):
		if i < _current_count:
			_count_dots[i].color = color
		else:
			_count_dots[i].color = Color(0.3, 0.3, 0.3, 0.8)

	# 顯示面板
	if not visible:
		visible = true
		modulate.a = 0.0
		var tween := create_tween()
		tween.tween_property(self, "modulate:a", 1.0, FADE_DURATION)

	# 升級時閃爍
	if is_level_up:
		var flash_tween := create_tween()
		flash_tween.tween_property(self, "modulate", Color(2.0, 2.0, 2.0, 1.0), 0.08)
		flash_tween.tween_property(self, "modulate", Color.WHITE, 0.15)

	# 啟動倒數
	_start_countdown()

## 稀有連擊廣播（全服）— 顯示在動態牆，不在此面板處理
func on_rare_catch_broadcast(_data: Dictionary) -> void:
	pass  # 由 ActivityFeedPanel 或 HUD 的 banner 處理

## 稀有連擊重置
func on_rare_catch_reset(_data: Dictionary) -> void:
	_is_active = false
	if _tick_timer != null:
		_tick_timer = null

	# 淡出
	if visible:
		var tween := create_tween()
		tween.tween_interval(RESET_HIDE_DELAY)
		tween.tween_property(self, "modulate:a", 0.0, FADE_DURATION)
		tween.tween_callback(func(): visible = false; modulate.a = 1.0)

# ---- 內部方法 ----

func _start_countdown() -> void:
	if _tick_timer != null:
		_tick_timer = null
	_schedule_tick()

func _schedule_tick() -> void:
	if not _is_active or _seconds_left <= 0:
		return
	_tick_timer = get_tree().create_timer(1.0)
	_tick_timer.timeout.connect(_on_tick)

func _on_tick() -> void:
	if not _is_active:
		return
	_seconds_left -= 1
	if _seconds_left <= 0:
		_timer_label.text = "0s"
		return
	_timer_label.text = str(_seconds_left) + "s"

	# 最後 15 秒變紅色
	if _seconds_left <= 15:
		var blink_color := Color(1.0, 0.3, 0.3) if (_seconds_left % 2 == 0) else Color(0.6, 0.6, 0.6)
		_timer_label.add_theme_color_override("font_color", blink_color)

	_schedule_tick()
