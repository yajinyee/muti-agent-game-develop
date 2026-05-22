## RainbowLuckyPanel.gd — 彩虹幸運魚面板（DAY-173）
## 業界依據：Fisch Roblox 2026「Rainbow Leviathan — rare rainbow fish that triggers a luck boost event」
## + Fish It 2026「Rainbow Throw — increases luck for rare fish」
## + Ocean King 2026「Rainbow Fish — all players receive a luck boost for 10 seconds」
## 視覺設計：
##   - lucky_start（全服）：全螢幕彩虹閃光 + 頂部橫幅「彩虹幸運時間！擊破機率 +20%！」
##     + 右上角彩虹倒數計時器（10秒）+ 畫面邊緣彩虹光暈
##   - lucky_end（全服）：彩虹光暈淡出 + 倒數計時器淡出
extends Node2D

# ---- 常數 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0
const RAINBOW_COLORS = [
	Color(1.0, 0.0, 0.0, 0.6),   # 紅
	Color(1.0, 0.5, 0.0, 0.6),   # 橙
	Color(1.0, 1.0, 0.0, 0.6),   # 黃
	Color(0.0, 1.0, 0.0, 0.6),   # 綠
	Color(0.0, 0.5, 1.0, 0.6),   # 藍
	Color(0.5, 0.0, 1.0, 0.6),   # 紫
]

# ---- 狀態 ----
var _pixel_font: Font = null
var _countdown_lbl: Label = null    # 右上角倒數計時
var _banner: ColorRect = null       # 頂部橫幅
var _edge_glows: Array = []         # 邊緣彩虹光暈
var _is_active: bool = false        # 是否激活中
var _elapsed: float = 0.0           # 已過時間
var _duration: float = 10.0         # 持續時間
var _color_index: int = 0           # 當前彩虹顏色索引
var _color_timer: float = 0.0       # 顏色切換計時

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("rainbow_lucky_fish"):
		GameManager.rainbow_lucky_fish.connect(_on_rainbow_lucky_fish)

# ---- 計時器 ----
func _process(delta: float) -> void:
	if not _is_active:
		return

	_elapsed += delta
	_color_timer += delta

	# 倒數計時更新
	var remaining = _duration - _elapsed
	if remaining <= 0.0:
		_is_active = false
		_cleanup()
		return

	if is_instance_valid(_countdown_lbl):
		_countdown_lbl.text = "🌈 %.1fs" % remaining

	# 彩虹顏色循環（每 0.3 秒切換一次）
	if _color_timer >= 0.3:
		_color_timer = 0.0
		_color_index = (_color_index + 1) % RAINBOW_COLORS.size()
		_update_edge_glow_color()

# ---- 訊號處理 ----
func _on_rainbow_lucky_fish(data: Dictionary) -> void:
	var phase = data.get("phase", "")
	match phase:
		"lucky_start":
			_handle_lucky_start(data)
		"lucky_end":
			_handle_lucky_end()

# ---- lucky_start：全服彩虹幸運時間開始 ----
func _handle_lucky_start(data: Dictionary) -> void:
	var player_name = data.get("player_name", "玩家")
	var duration_sec = data.get("duration_sec", 10)
	var kill_boost = data.get("kill_boost", 0.20)

	_is_active = true
	_elapsed = 0.0
	_duration = float(duration_sec)
	_color_index = 0
	_color_timer = 0.0

	# 全螢幕彩虹閃光（6色循環）
	_rainbow_flash_sequence()

	# 建立頂部橫幅
	_create_banner(player_name, kill_boost)

	# 建立右上角倒數計時
	_create_countdown()

	# 建立邊緣彩虹光暈
	_create_edge_glows()

# ---- lucky_end：彩虹幸運時間結束 ----
func _handle_lucky_end() -> void:
	_is_active = false
	_cleanup()

# ---- 輔助：彩虹閃光序列 ----
func _rainbow_flash_sequence() -> void:
	for i in range(RAINBOW_COLORS.size()):
		var color = RAINBOW_COLORS[i]
		color.a = 0.4
		var flash = ColorRect.new()
		flash.size = Vector2(SCREEN_W, SCREEN_H)
		flash.position = Vector2(0, 0)
		flash.color = color
		add_child(flash)

		var tween = flash.create_tween()
		tween.tween_interval(float(i) * 0.08)
		tween.tween_property(flash, "modulate:a", 0.0, 0.2)
		tween.tween_callback(func():
			if is_instance_valid(flash):
				flash.queue_free()
		)

# ---- 輔助：建立頂部橫幅 ----
func _create_banner(player_name: String, kill_boost: float) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()

	_banner = ColorRect.new()
	_banner.size = Vector2(SCREEN_W, 40)
	_banner.position = Vector2(0, -40)
	_banner.color = Color(0.1, 0.0, 0.2, 0.9)
	add_child(_banner)

	var lbl = Label.new()
	lbl.text = "🌈 %s 觸發彩虹幸運魚！擊破機率 +%.0f%%！" % [player_name, kill_boost * 100]
	lbl.position = Vector2(10, 8)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 14)
	lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_banner.add_child(lbl)

	# 從頂部滑入
	_banner.position = Vector2(0, 0)
	var tween = _banner.create_tween()
	tween.tween_interval(_duration - 0.5)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_banner):
			_banner.queue_free()
			_banner = null
	)

# ---- 輔助：建立右上角倒數計時 ----
func _create_countdown() -> void:
	if is_instance_valid(_countdown_lbl):
		_countdown_lbl.queue_free()

	_countdown_lbl = Label.new()
	_countdown_lbl.text = "🌈 %.1fs" % _duration
	_countdown_lbl.position = Vector2(SCREEN_W - 110, 50)
	if _pixel_font:
		_countdown_lbl.add_theme_font_override("font", _pixel_font)
		_countdown_lbl.add_theme_font_size_override("font_size", 16)
	_countdown_lbl.add_theme_color_override("font_color", Color(1.0, 0.41, 0.71))
	add_child(_countdown_lbl)

	# 彈跳動畫
	var tween = _countdown_lbl.create_tween()
	tween.tween_property(_countdown_lbl, "scale", Vector2(1.3, 1.3), 0.1)
	tween.tween_property(_countdown_lbl, "scale", Vector2(1.0, 1.0), 0.15)

# ---- 輔助：建立邊緣彩虹光暈 ----
func _create_edge_glows() -> void:
	_cleanup_edge_glows()

	# 四條邊緣光暈（上下左右）
	var edges = [
		{"size": Vector2(SCREEN_W, 8), "pos": Vector2(0, 0)},
		{"size": Vector2(SCREEN_W, 8), "pos": Vector2(0, SCREEN_H - 8)},
		{"size": Vector2(8, SCREEN_H), "pos": Vector2(0, 0)},
		{"size": Vector2(8, SCREEN_H), "pos": Vector2(SCREEN_W - 8, 0)},
	]

	for edge_def in edges:
		var glow = ColorRect.new()
		glow.size = edge_def["size"]
		glow.position = edge_def["pos"]
		glow.color = RAINBOW_COLORS[0]
		add_child(glow)
		_edge_glows.append(glow)

# ---- 輔助：更新邊緣光暈顏色 ----
func _update_edge_glow_color() -> void:
	var color = RAINBOW_COLORS[_color_index]
	for glow in _edge_glows:
		if is_instance_valid(glow):
			glow.color = color

# ---- 輔助：清理邊緣光暈 ----
func _cleanup_edge_glows() -> void:
	for glow in _edge_glows:
		if is_instance_valid(glow):
			glow.queue_free()
	_edge_glows.clear()

# ---- 輔助：清理所有 UI ----
func _cleanup() -> void:
	if is_instance_valid(_countdown_lbl):
		var tween = _countdown_lbl.create_tween()
		tween.tween_property(_countdown_lbl, "modulate:a", 0.0, 0.3)
		tween.tween_callback(func():
			if is_instance_valid(_countdown_lbl):
				_countdown_lbl.queue_free()
				_countdown_lbl = null
		)

	if is_instance_valid(_banner):
		var tween = _banner.create_tween()
		tween.tween_property(_banner, "modulate:a", 0.0, 0.3)
		tween.tween_callback(func():
			if is_instance_valid(_banner):
				_banner.queue_free()
				_banner = null
		)

	# 邊緣光暈淡出
	for glow in _edge_glows:
		if is_instance_valid(glow):
			var tween = glow.create_tween()
			tween.tween_property(glow, "modulate:a", 0.0, 0.5)
			tween.tween_callback(func():
				if is_instance_valid(glow):
					glow.queue_free()
			)
	_edge_glows.clear()
