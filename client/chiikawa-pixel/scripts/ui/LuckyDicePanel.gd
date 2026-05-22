## LuckyDicePanel.gd — 幸運骰子魚面板（DAY-175）
## 業界依據：Ocean King 3 Plus「Fast Bomb — randomly triggered bonus」
## + 捕魚機業界「Dice Roll bonus — roll dice to determine reward multiplier」
## 視覺設計：
##   - dice_start（個人）：全螢幕金色閃光 + 中央兩顆骰子滾動動畫（2秒）
##   - dice_broadcast（全服）：頂部小橫幅「有人觸發幸運骰子！」
##   - dice_result（個人）：骰子緩停顯示點數 + 結果彈窗（點數/獎勵/標籤）
##     - 點數7：金色光暈；點數12：橙紅雙閃光；點數2：紫色閃光
##   - dice_jackpot（全服）：全服廣播橫幅「XXX 擲出大六！」
extends Node2D

# ---- 常數 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0
const DICE_SIZE := 64.0
const DICE_FACES = [
	[],                                    # 0（不用）
	[[0.5, 0.5]],                          # 1：中心
	[[0.25, 0.25], [0.75, 0.75]],          # 2：對角
	[[0.25, 0.25], [0.5, 0.5], [0.75, 0.75]], # 3：對角+中心
	[[0.25, 0.25], [0.75, 0.25], [0.25, 0.75], [0.75, 0.75]], # 4：四角
	[[0.25, 0.25], [0.75, 0.25], [0.5, 0.5], [0.25, 0.75], [0.75, 0.75]], # 5：四角+中心
	[[0.25, 0.2], [0.75, 0.2], [0.25, 0.5], [0.75, 0.5], [0.25, 0.8], [0.75, 0.8]], # 6：兩列三行
]

# ---- 狀態 ----
var _pixel_font: Font = null
var _dice_container: Node2D = null  # 骰子容器
var _die1_node: Node2D = null       # 骰子1節點
var _die2_node: Node2D = null       # 骰子2節點
var _is_rolling: bool = false       # 是否正在滾動
var _roll_elapsed: float = 0.0      # 滾動已過時間
var _roll_duration: float = 2.0     # 滾動持續時間
var _roll_face1: int = 1            # 骰子1當前顯示面
var _roll_face2: int = 1            # 骰子2當前顯示面
var _face_timer: float = 0.0        # 換面計時
var _result_panel: Node2D = null    # 結果彈窗

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("lucky_dice_fish"):
		GameManager.lucky_dice_fish.connect(_on_lucky_dice_fish)

# ---- 計時器 ----
func _process(delta: float) -> void:
	if not _is_rolling:
		return

	_roll_elapsed += delta
	_face_timer += delta

	# 每 0.1 秒換一次骰子面（滾動效果）
	if _face_timer >= 0.1:
		_face_timer = 0.0
		_roll_face1 = randi() % 6 + 1
		_roll_face2 = randi() % 6 + 1
		_update_dice_display(_roll_face1, _roll_face2)

	if _roll_elapsed >= _roll_duration:
		_is_rolling = false

# ---- 訊號處理 ----
func _on_lucky_dice_fish(data: Dictionary) -> void:
	var phase = data.get("phase", "")
	match phase:
		"dice_start":
			_handle_dice_start(data)
		"dice_broadcast":
			_handle_dice_broadcast(data)
		"dice_result":
			_handle_dice_result(data)
		"dice_jackpot":
			_handle_dice_jackpot(data)

# ---- dice_start：骰子開始滾動 ----
func _handle_dice_start(data: Dictionary) -> void:
	var roll_ms = data.get("roll_ms", 2000)
	_roll_duration = float(roll_ms) / 1000.0

	# 全螢幕金色閃光
	_flash_screen(Color(1.0, 0.85, 0.0, 0.4))

	# 建立骰子容器
	if is_instance_valid(_dice_container):
		_dice_container.queue_free()

	_dice_container = Node2D.new()
	_dice_container.position = Vector2(SCREEN_W / 2.0, SCREEN_H / 2.0)
	add_child(_dice_container)

	# 建立兩顆骰子
	_die1_node = _create_die_node(-DICE_SIZE - 10, 0)
	_die2_node = _create_die_node(10, 0)
	_dice_container.add_child(_die1_node)
	_dice_container.add_child(_die2_node)

	# 開始滾動
	_is_rolling = true
	_roll_elapsed = 0.0
	_face_timer = 0.0

	# 彈跳動畫
	var tween = _dice_container.create_tween()
	tween.tween_property(_dice_container, "scale", Vector2(1.3, 1.3), 0.15)
	tween.tween_property(_dice_container, "scale", Vector2(1.0, 1.0), 0.15)

# ---- dice_broadcast：全服廣播橫幅 ----
func _handle_dice_broadcast(data: Dictionary) -> void:
	var player_name = data.get("player_name", "玩家")
	_show_broadcast_banner("🎲 %s 觸發幸運骰子！" % player_name)

# ---- dice_result：骰子結果 ----
func _handle_dice_result(data: Dictionary) -> void:
	var die1 = data.get("die1", 1)
	var die2 = data.get("die2", 1)
	var sum = data.get("sum", 2)
	var reward = data.get("reward", 0)
	var label = data.get("label", "")

	# 停止滾動，顯示最終點數
	_is_rolling = false
	_update_dice_display(die1, die2)

	# 根據點數決定特效
	if sum == 12:
		_flash_screen(Color(1.0, 0.27, 0.0, 0.5))
		var tween = create_tween()
		tween.tween_interval(0.15)
		tween.tween_callback(func(): _flash_screen(Color(1.0, 0.27, 0.0, 0.5)))
	elif sum == 2:
		_flash_screen(Color(0.58, 0.0, 0.83, 0.5))
	elif sum == 7:
		_flash_screen(Color(1.0, 0.85, 0.0, 0.4))

	# 骰子緩停動畫（縮放）
	if is_instance_valid(_dice_container):
		var tween = _dice_container.create_tween()
		tween.tween_property(_dice_container, "scale", Vector2(1.2, 1.2), 0.1)
		tween.tween_property(_dice_container, "scale", Vector2(1.0, 1.0), 0.15)

	# 建立結果彈窗
	_show_result_panel(die1, die2, sum, reward, label)

	# 3 秒後清理骰子
	if is_instance_valid(_dice_container):
		var tween = _dice_container.create_tween()
		tween.tween_interval(3.0)
		tween.tween_property(_dice_container, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func():
			if is_instance_valid(_dice_container):
				_dice_container.queue_free()
				_dice_container = null
		)

# ---- dice_jackpot：大六全服廣播 ----
func _handle_dice_jackpot(data: Dictionary) -> void:
	var player_name = data.get("player_name", "玩家")
	var reward = data.get("reward", 0)
	_show_broadcast_banner("🎲🎲 %s 擲出大六！獲得 %d 金幣！" % [player_name, reward])

# ---- 輔助：建立骰子節點 ----
func _create_die_node(offset_x: float, offset_y: float) -> Node2D:
	var die = Node2D.new()
	die.position = Vector2(offset_x, offset_y)

	# 骰子背景（白色圓角矩形）
	var bg = ColorRect.new()
	bg.size = Vector2(DICE_SIZE, DICE_SIZE)
	bg.position = Vector2(0, -DICE_SIZE / 2.0)
	bg.color = Color(0.95, 0.95, 0.95)
	die.add_child(bg)

	return die

# ---- 輔助：更新骰子顯示 ----
func _update_dice_display(face1: int, face2: int) -> void:
	if is_instance_valid(_die1_node):
		_draw_die_face(_die1_node, face1)
	if is_instance_valid(_die2_node):
		_draw_die_face(_die2_node, face2)

# ---- 輔助：繪製骰子面 ----
func _draw_die_face(die_node: Node2D, face: int) -> void:
	# 清除舊的點
	for child in die_node.get_children():
		if child.name.begins_with("dot_"):
			child.queue_free()

	if face < 1 or face > 6:
		return

	var dots = DICE_FACES[face]
	for i in range(dots.size()):
		var dot_pos = dots[i]
		var dot = ColorRect.new()
		dot.name = "dot_%d" % i
		dot.size = Vector2(10, 10)
		dot.position = Vector2(
			dot_pos[0] * DICE_SIZE - 5,
			dot_pos[1] * DICE_SIZE - DICE_SIZE / 2.0 - 5
		)
		dot.color = Color(0.1, 0.1, 0.1)
		die_node.add_child(dot)

# ---- 輔助：顯示結果彈窗 ----
func _show_result_panel(die1: int, die2: int, sum: int, reward: int, label: String) -> void:
	if is_instance_valid(_result_panel):
		_result_panel.queue_free()

	_result_panel = Node2D.new()
	add_child(_result_panel)

	# 背景顏色依點數
	var bg_color = Color(0.1, 0.1, 0.1, 0.9)
	if sum == 12:
		bg_color = Color(0.3, 0.1, 0.0, 0.9)
	elif sum == 2:
		bg_color = Color(0.15, 0.0, 0.2, 0.9)

	var bg = ColorRect.new()
	bg.size = Vector2(220, 110)
	bg.position = Vector2(0, -55)
	bg.color = bg_color
	_result_panel.add_child(bg)

	# 標籤
	var label_lbl = Label.new()
	label_lbl.text = label
	label_lbl.position = Vector2(10, -50)
	if _pixel_font:
		label_lbl.add_theme_font_override("font", _pixel_font)
		label_lbl.add_theme_font_size_override("font_size", 14)
	label_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_result_panel.add_child(label_lbl)

	# 點數顯示
	var sum_lbl = Label.new()
	sum_lbl.text = "🎲 %d + %d = %d" % [die1, die2, sum]
	sum_lbl.position = Vector2(10, -28)
	if _pixel_font:
		sum_lbl.add_theme_font_override("font", _pixel_font)
		sum_lbl.add_theme_font_size_override("font_size", 13)
	sum_lbl.add_theme_color_override("font_color", Color.WHITE)
	_result_panel.add_child(sum_lbl)

	# 獎勵
	var reward_lbl = Label.new()
	reward_lbl.text = "🪙 獎勵：+%d 金幣" % reward
	reward_lbl.position = Vector2(10, -6)
	if _pixel_font:
		reward_lbl.add_theme_font_override("font", _pixel_font)
		reward_lbl.add_theme_font_size_override("font_size", 13)
	reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_result_panel.add_child(reward_lbl)

	# 從右側滑入
	_result_panel.position = Vector2(SCREEN_W + 50, SCREEN_H / 2.0)
	var tween = _result_panel.create_tween()
	tween.tween_property(_result_panel, "position:x", SCREEN_W - 240.0, 0.4).set_ease(Tween.EASE_OUT)

	# 3 秒後淡出
	tween.tween_interval(3.0)
	tween.tween_property(_result_panel, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_result_panel):
			_result_panel.queue_free()
			_result_panel = null
	)

# ---- 輔助：全服廣播橫幅 ----
func _show_broadcast_banner(text: String) -> void:
	var banner = ColorRect.new()
	banner.size = Vector2(SCREEN_W, 34)
	banner.position = Vector2(0, 0)
	banner.color = Color(0.1, 0.08, 0.0, 0.9)
	add_child(banner)

	var lbl = Label.new()
	lbl.text = text
	lbl.position = Vector2(10, 6)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 13)
	lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	banner.add_child(lbl)

	var tween = banner.create_tween()
	tween.tween_interval(2.5)
	tween.tween_property(banner, "modulate:a", 0.0, 0.3)
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
