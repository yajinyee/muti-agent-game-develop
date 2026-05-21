## GoldenTimePanel.gd — 黃金時間系統面板（DAY-125）
## 業界依據：Fire Kirin / Ocean King 系列的 Golden Time 機制
## 全場目標物倍率暫時提升，製造「全場瘋狂」的高峰體驗
## 業界研究顯示 Golden Time 讓短期參與度提升 40%+
extends Control

# ---- 常數 ----
const PANEL_WIDTH  := 400
const PANEL_HEIGHT := 120
const BANNER_DURATION := 4.0   # 橫幅顯示 4 秒後縮回
const FLASH_INTERVAL  := 0.25  # 閃爍間隔

# 等級顏色
const TIER_COLORS = {
	0: Color(0.75, 0.75, 0.75, 1.0),  # Silver：銀色
	1: Color(1.00, 0.85, 0.00, 1.0),  # Gold：金色
	2: Color(1.00, 0.40, 0.80, 1.0),  # Rainbow：粉紅
}
const TIER_BG_COLORS = {
	0: Color(0.10, 0.10, 0.15, 0.92),
	1: Color(0.15, 0.10, 0.00, 0.92),
	2: Color(0.15, 0.00, 0.15, 0.92),
}

# ---- 節點引用 ----
var _bg: ColorRect
var _top_bar: ColorRect
var _icon_label: Label
var _tier_label: Label
var _mult_label: Label
var _timer_label: Label
var _desc_label: Label
var _flash_overlay: ColorRect

# ---- 狀態 ----
var _is_active := false
var _current_tier := 1
var _seconds_left := 0
var _mult_boost := 2.0
var _flash_count := 0
var _flash_timer: SceneTreeTimer = null
var _banner_timer: SceneTreeTimer = null
var _tick_timer: SceneTreeTimer = null

func _ready() -> void:
	_build_ui()
	visible = false
	set_process(false)

func _build_ui() -> void:
	custom_minimum_size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)

	# 背景
	_bg = ColorRect.new()
	_bg.color = TIER_BG_COLORS[1]
	_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	add_child(_bg)

	# 頂部色條
	_top_bar = ColorRect.new()
	_top_bar.color = TIER_COLORS[1]
	_top_bar.size = Vector2(PANEL_WIDTH, 5)
	_top_bar.position = Vector2(0, 0)
	add_child(_top_bar)

	# 底部色條
	var bot_bar := ColorRect.new()
	bot_bar.color = TIER_COLORS[1]
	bot_bar.size = Vector2(PANEL_WIDTH, 5)
	bot_bar.position = Vector2(0, PANEL_HEIGHT - 5)
	add_child(bot_bar)

	# 圖示
	_icon_label = Label.new()
	_icon_label.text = "✨"
	_icon_label.position = Vector2(12, 14)
	_icon_label.add_theme_font_size_override("font_size", 36)
	add_child(_icon_label)

	# 等級名稱
	_tier_label = Label.new()
	_tier_label.text = "✨ 黃金時間"
	_tier_label.position = Vector2(60, 12)
	_tier_label.add_theme_font_size_override("font_size", 20)
	_tier_label.add_theme_color_override("font_color", TIER_COLORS[1])
	add_child(_tier_label)

	# 倍率顯示
	_mult_label = Label.new()
	_mult_label.text = "全場獎勵 ×2.0"
	_mult_label.position = Vector2(60, 40)
	_mult_label.add_theme_font_size_override("font_size", 16)
	_mult_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.6))
	add_child(_mult_label)

	# 說明文字
	_desc_label = Label.new()
	_desc_label.text = "把握機會！擊破目標獲得更多獎勵！"
	_desc_label.position = Vector2(60, 64)
	_desc_label.add_theme_font_size_override("font_size", 12)
	_desc_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	add_child(_desc_label)

	# 倒數計時
	_timer_label = Label.new()
	_timer_label.text = "45s"
	_timer_label.position = Vector2(PANEL_WIDTH - 60, 14)
	_timer_label.add_theme_font_size_override("font_size", 28)
	_timer_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	add_child(_timer_label)

	# 閃光覆蓋層（全螢幕閃光用）
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 0.85, 0.0, 0.0)
	_flash_overlay.size = Vector2(1280, 720)
	_flash_overlay.position = Vector2(-640, -360)  # 相對於面板中心
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

# ---- 公開 API ----

## 黃金時間開始
func show_golden_time_start(data: Dictionary) -> void:
	_is_active = true
	_current_tier = data.get("tier", 1)
	_seconds_left = data.get("seconds_left", data.get("duration", 45))
	_mult_boost = data.get("mult_boost", 2.0)

	var color := TIER_COLORS.get(_current_tier, TIER_COLORS[1])
	var bg_color := TIER_BG_COLORS.get(_current_tier, TIER_BG_COLORS[1])

	# 更新 UI
	_bg.color = bg_color
	_top_bar.color = color
	_icon_label.text = data.get("icon", "✨")
	_tier_label.text = data.get("tier_name", "✨ 黃金時間")
	_tier_label.add_theme_color_override("font_color", color)

	var mult_str := "%.1f" % _mult_boost
	if _mult_boost == float(int(_mult_boost)):
		mult_str = "%d" % int(_mult_boost)
	_mult_label.text = "全場獎勵 ×" + mult_str
	_timer_label.text = str(_seconds_left) + "s"

	# 定位：頂部中央滑入
	position = Vector2(640 - PANEL_WIDTH / 2, -PANEL_HEIGHT)
	visible = true

	# 滑入動畫
	var tween := create_tween()
	tween.tween_property(self, "position:y", 8, 0.4).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)

	# 全螢幕閃光
	_do_screen_flash(color, _current_tier)

	# 啟動倒數計時
	_start_countdown()

	# 自動隱藏（黃金時間結束後）
	if _banner_timer != null:
		_banner_timer = null
	_banner_timer = get_tree().create_timer(float(_seconds_left) + 1.0)
	_banner_timer.timeout.connect(_on_banner_timeout)

## 黃金時間結束
func show_golden_time_end(_data: Dictionary) -> void:
	_is_active = false
	if _tick_timer != null:
		_tick_timer = null

	# 更新顯示
	_timer_label.text = "結束"
	_timer_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	_mult_label.text = "黃金時間結束"

	# 1.5 秒後滑出
	var tween := create_tween()
	tween.tween_interval(1.5)
	tween.tween_property(self, "position:y", -PANEL_HEIGHT, 0.35).set_ease(Tween.EASE_IN)
	tween.tween_callback(func(): visible = false)

## 更新狀態（玩家加入時同步）
func update_status(data: Dictionary) -> void:
	if not data.get("is_active", false):
		visible = false
		return
	show_golden_time_start(data)

# ---- 內部方法 ----

func _do_screen_flash(color: Color, tier: int) -> void:
	# 全螢幕閃光
	_flash_overlay.color = Color(color.r, color.g, color.b, 0.0)
	var flash_tween := create_tween()
	flash_tween.tween_property(_flash_overlay, "color:a", 0.35, 0.08)
	flash_tween.tween_property(_flash_overlay, "color:a", 0.0, 0.25)

	# 彩虹等級：額外閃爍
	if tier == 2:
		flash_tween.tween_property(_flash_overlay, "color:a", 0.25, 0.08)
		flash_tween.tween_property(_flash_overlay, "color:a", 0.0, 0.25)

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

	# 最後 10 秒變紅色閃爍
	if _seconds_left <= 10:
		var blink_color := Color(1.0, 0.3, 0.3) if (_seconds_left % 2 == 0) else Color(1.0, 0.9, 0.3)
		_timer_label.add_theme_color_override("font_color", blink_color)

	_schedule_tick()

func _on_banner_timeout() -> void:
	if not _is_active:
		return
	# 黃金時間自然結束（如果 Server 沒有發 end 訊息）
	var tween := create_tween()
	tween.tween_property(self, "position:y", -PANEL_HEIGHT, 0.35).set_ease(Tween.EASE_IN)
	tween.tween_callback(func(): visible = false)
	_is_active = false
