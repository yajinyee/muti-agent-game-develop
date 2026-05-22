## LuckyEggPanel.gd — 幸運彩蛋魚面板（DAY-172）
## 業界依據：JILI Mega Fishing 2026「Giant Prize Fish lets you easily win great prizes,
## with the chance for 5x multipliers」+ Ocean King 2026「Egg Fish drops golden eggs」
## 視覺設計：
##   - egg_start（全服）：頂部小橫幅「有人觸發幸運彩蛋魚！」+ 彩蛋掉落動畫
##   - egg_open（個人）：彩蛋從觸發位置飛出 + 開啟動畫 + 獎勵浮動文字
##     - coins：金色彩蛋 + 金幣雨動畫 + "+XXX 金幣" 浮動文字
##     - mult：粉紅彩蛋 + "×2 加成 5s" 浮動文字 + 右上角倒數計時
##     - weapon：天藍彩蛋 + "武器充能 ×1" 浮動文字 + 武器圖示閃爍
##   - egg_result（個人）：右側滑入結果彈窗（彩蛋數/金幣/倍率/武器）
##   - mult_end（個人）：右上角倒數計時淡出
##   - ≥4個彩蛋：全服廣播橫幅；≥5個：金色雙閃光
extends Node2D

# ---- 常數 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0
const EGG_SIZE := 32.0
const EGG_COLORS = {
	"coins":  Color(1.0, 0.85, 0.0),   # 金色
	"mult":   Color(1.0, 0.41, 0.71),  # 粉紅色
	"weapon": Color(0.0, 0.75, 1.0),   # 天藍色
}
const EGG_ICONS = {
	"coins":  "🪙",
	"mult":   "✨",
	"weapon": "⚡",
}

# ---- 狀態 ----
var _pixel_font: Font = null
var _mult_countdown_lbl: Label = null   # 倍率倒數計時標籤
var _mult_elapsed: float = 0.0          # 倍率已過時間
var _mult_duration: float = 5.0         # 倍率持續時間
var _is_mult_active: bool = false       # 是否倍率激活中
var _mult_stack: int = 0                # 倍率疊加次數（多個彩蛋可疊加）
var _egg_nodes: Array = []              # 彩蛋節點列表
var _result_panel: Node2D = null        # 結果彈窗節點

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("lucky_egg_fish"):
		GameManager.lucky_egg_fish.connect(_on_lucky_egg_fish)

# ---- 計時器 ----
func _process(delta: float) -> void:
	# 倍率倒數計時
	if _is_mult_active:
		_mult_elapsed += delta
		var remaining = _mult_duration - _mult_elapsed
		if remaining <= 0.0:
			_is_mult_active = false
			if is_instance_valid(_mult_countdown_lbl):
				_mult_countdown_lbl.queue_free()
				_mult_countdown_lbl = null
		elif is_instance_valid(_mult_countdown_lbl):
			_mult_countdown_lbl.text = "×2 %.1fs" % remaining

# ---- 訊號處理 ----
func _on_lucky_egg_fish(data: Dictionary) -> void:
	var phase = data.get("phase", "")
	match phase:
		"egg_start":
			_handle_egg_start(data)
		"egg_open":
			_handle_egg_open(data)
		"egg_result":
			_handle_egg_result(data)
		"egg_broadcast":
			_handle_egg_broadcast(data)
		"mult_end":
			_handle_mult_end()

# ---- egg_start：全服廣播橫幅 ----
func _handle_egg_start(data: Dictionary) -> void:
	var player_name = data.get("player_name", "玩家")
	var egg_count = data.get("egg_count", 1)
	_show_broadcast_banner("🥚 %s 觸發幸運彩蛋魚！掉落 %d 個彩蛋！" % [player_name, egg_count])

# ---- egg_open：個人彩蛋開啟動畫 ----
func _handle_egg_open(data: Dictionary) -> void:
	var egg_result = data.get("egg_result", {})
	var reward_type = egg_result.get("reward_type", "coins")
	var label_text = egg_result.get("label", "")
	var egg_index = data.get("egg_index", 0)
	var trigger_x = data.get("trigger_x", SCREEN_W / 2.0)
	var trigger_y = data.get("trigger_y", SCREEN_H / 2.0)

	# 彩蛋顏色
	var egg_color = EGG_COLORS.get(reward_type, Color.WHITE)
	var egg_icon = EGG_ICONS.get(reward_type, "🥚")

	# 建立彩蛋節點（從觸發位置飛出）
	var egg_node = Node2D.new()
	add_child(egg_node)
	_egg_nodes.append(egg_node)

	# 彩蛋圓形
	var egg_circle = ColorRect.new()
	egg_circle.size = Vector2(EGG_SIZE, EGG_SIZE)
	egg_circle.position = Vector2(-EGG_SIZE / 2.0, -EGG_SIZE / 2.0)
	egg_circle.color = egg_color
	egg_node.add_child(egg_circle)

	# 彩蛋圖示
	var icon_lbl = Label.new()
	icon_lbl.text = egg_icon
	icon_lbl.position = Vector2(-12, -14)
	if _pixel_font:
		icon_lbl.add_theme_font_override("font", _pixel_font)
		icon_lbl.add_theme_font_size_override("font_size", 20)
	egg_node.add_child(icon_lbl)

	# 起始位置（觸發位置）
	egg_node.position = Vector2(trigger_x, trigger_y)

	# 飛出目標位置（分散在畫面中央）
	var spread_x = SCREEN_W / 2.0 + (egg_index - 2) * 80.0
	var spread_y = SCREEN_H / 2.0 - 50.0

	# 飛出動畫
	var tween = egg_node.create_tween()
	tween.tween_property(egg_node, "position",
		Vector2(spread_x, spread_y), 0.3).set_ease(Tween.EASE_OUT)
	tween.tween_interval(0.1)

	# 開啟動畫（縮放爆炸）
	tween.tween_property(egg_node, "scale", Vector2(1.5, 1.5), 0.1)
	tween.tween_property(egg_node, "scale", Vector2(0.8, 0.8), 0.1)
	tween.tween_property(egg_node, "scale", Vector2(1.0, 1.0), 0.1)

	# 顯示獎勵浮動文字
	tween.tween_callback(func():
		_show_reward_float(spread_x, spread_y - 40.0, label_text, egg_color)
	)

	# 特殊處理：倍率加成
	if reward_type == "mult":
		tween.tween_callback(func():
			_activate_mult_display()
		)

	# 淡出彩蛋
	tween.tween_interval(0.5)
	tween.tween_property(egg_node, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(egg_node):
			egg_node.queue_free()
		_egg_nodes.erase(egg_node)
	)

# ---- egg_result：個人結果彈窗 ----
func _handle_egg_result(data: Dictionary) -> void:
	var egg_count = data.get("egg_count", 1)
	var total_coins = data.get("total_coins", 0)
	var mult_count = data.get("mult_count", 0)
	var weapon_count = data.get("weapon_count", 0)

	# 建立結果彈窗（右側滑入）
	if is_instance_valid(_result_panel):
		_result_panel.queue_free()

	_result_panel = Node2D.new()
	add_child(_result_panel)

	# 背景
	var bg = ColorRect.new()
	bg.size = Vector2(220, 140)
	bg.position = Vector2(0, -70)
	bg.color = Color(0.1, 0.1, 0.1, 0.85)
	_result_panel.add_child(bg)

	# 標題
	var title_lbl = Label.new()
	title_lbl.text = "🥚 彩蛋結果"
	title_lbl.position = Vector2(10, -65)
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
		title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_result_panel.add_child(title_lbl)

	# 彩蛋數
	var count_lbl = Label.new()
	count_lbl.text = "彩蛋數：%d 個" % egg_count
	count_lbl.position = Vector2(10, -45)
	if _pixel_font:
		count_lbl.add_theme_font_override("font", _pixel_font)
		count_lbl.add_theme_font_size_override("font_size", 12)
	count_lbl.add_theme_color_override("font_color", Color.WHITE)
	_result_panel.add_child(count_lbl)

	# 金幣獎勵
	if total_coins > 0:
		var coins_lbl = Label.new()
		coins_lbl.text = "🪙 金幣：+%d" % total_coins
		coins_lbl.position = Vector2(10, -25)
		if _pixel_font:
			coins_lbl.add_theme_font_override("font", _pixel_font)
			coins_lbl.add_theme_font_size_override("font_size", 12)
		coins_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
		_result_panel.add_child(coins_lbl)

	# 倍率加成
	if mult_count > 0:
		var mult_lbl = Label.new()
		mult_lbl.text = "✨ 倍率加成：×%d" % mult_count
		mult_lbl.position = Vector2(10, -5)
		if _pixel_font:
			mult_lbl.add_theme_font_override("font", _pixel_font)
			mult_lbl.add_theme_font_size_override("font_size", 12)
		mult_lbl.add_theme_color_override("font_color", Color(1.0, 0.41, 0.71))
		_result_panel.add_child(mult_lbl)

	# 武器充能
	if weapon_count > 0:
		var weapon_lbl = Label.new()
		weapon_lbl.text = "⚡ 武器充能：×%d" % weapon_count
		weapon_lbl.position = Vector2(10, 15)
		if _pixel_font:
			weapon_lbl.add_theme_font_override("font", _pixel_font)
			weapon_lbl.add_theme_font_size_override("font_size", 12)
		weapon_lbl.add_theme_color_override("font_color", Color(0.0, 0.75, 1.0))
		_result_panel.add_child(weapon_lbl)

	# 從右側滑入
	_result_panel.position = Vector2(SCREEN_W + 50, SCREEN_H / 2.0)
	var tween = _result_panel.create_tween()
	tween.tween_property(_result_panel, "position:x", SCREEN_W - 240.0, 0.4).set_ease(Tween.EASE_OUT)

	# ≥5個彩蛋：金色雙閃光
	if egg_count >= 5:
		_flash_screen(Color(1.0, 0.85, 0.0, 0.5))
		tween.tween_interval(0.15)
		tween.tween_callback(func(): _flash_screen(Color(1.0, 0.85, 0.0, 0.5)))

	# 3 秒後淡出
	tween.tween_interval(3.0)
	tween.tween_property(_result_panel, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_result_panel):
			_result_panel.queue_free()
			_result_panel = null
	)

# ---- egg_broadcast：全服廣播橫幅 ----
func _handle_egg_broadcast(data: Dictionary) -> void:
	var player_name = data.get("player_name", "玩家")
	var egg_count = data.get("egg_count", 1)
	var total_coins = data.get("total_coins", 0)
	var mult_count = data.get("mult_count", 0)

	var msg = "🥚 %s 幸運彩蛋魚掉落 %d 個彩蛋！" % [player_name, egg_count]
	if mult_count > 0:
		msg += " %d 次倍率加成！" % mult_count
	if total_coins > 0:
		msg += " +%d 金幣！" % total_coins
	_show_broadcast_banner(msg)

# ---- mult_end：倍率結束 ----
func _handle_mult_end() -> void:
	_is_mult_active = false
	_mult_stack = 0
	if is_instance_valid(_mult_countdown_lbl):
		var tween = _mult_countdown_lbl.create_tween()
		tween.tween_property(_mult_countdown_lbl, "modulate:a", 0.0, 0.3)
		tween.tween_callback(func():
			if is_instance_valid(_mult_countdown_lbl):
				_mult_countdown_lbl.queue_free()
				_mult_countdown_lbl = null
		)

# ---- 輔助：激活倍率顯示 ----
func _activate_mult_display() -> void:
	_is_mult_active = true
	_mult_elapsed = 0.0
	_mult_duration = 5.0
	_mult_stack += 1

	# 建立右上角倒數計時標籤
	if not is_instance_valid(_mult_countdown_lbl):
		_mult_countdown_lbl = Label.new()
		add_child(_mult_countdown_lbl)
		_mult_countdown_lbl.position = Vector2(SCREEN_W - 120, 60)
		if _pixel_font:
			_mult_countdown_lbl.add_theme_font_override("font", _pixel_font)
			_mult_countdown_lbl.add_theme_font_size_override("font_size", 16)
		_mult_countdown_lbl.add_theme_color_override("font_color", Color(1.0, 0.41, 0.71))

	_mult_countdown_lbl.text = "×2 5.0s"

	# 彈跳動畫
	var tween = _mult_countdown_lbl.create_tween()
	tween.tween_property(_mult_countdown_lbl, "scale", Vector2(1.4, 1.4), 0.1)
	tween.tween_property(_mult_countdown_lbl, "scale", Vector2(1.0, 1.0), 0.15)

# ---- 輔助：獎勵浮動文字 ----
func _show_reward_float(x: float, y: float, text: String, color: Color) -> void:
	var lbl = Label.new()
	lbl.text = text
	lbl.position = Vector2(x - 40, y)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 14)
	lbl.add_theme_color_override("font_color", color)
	add_child(lbl)

	var tween = lbl.create_tween()
	tween.tween_property(lbl, "position:y", y - 40.0, 0.8).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.8)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)

# ---- 輔助：全服廣播橫幅 ----
func _show_broadcast_banner(text: String) -> void:
	var banner = ColorRect.new()
	banner.size = Vector2(SCREEN_W, 36)
	banner.position = Vector2(0, -36)
	banner.color = Color(0.1, 0.1, 0.1, 0.85)
	add_child(banner)

	var lbl = Label.new()
	lbl.text = text
	lbl.position = Vector2(10, 6)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 14)
	lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	banner.add_child(lbl)

	# 從頂部滑入
	banner.position = Vector2(0, 0)
	var tween = banner.create_tween()
	tween.tween_interval(2.5)
	tween.tween_property(banner, "position:y", -36.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(banner):
			banner.queue_free()
	)

# ---- 輔助：全螢幕閃光 ----
func _flash_screen(color: Color) -> void:
	var flash = ColorRect.new()
	flash.size = Vector2(SCREEN_W, SCREEN_H)
	flash.position = Vector2(0, 0)
	flash.color = color
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.25)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)
