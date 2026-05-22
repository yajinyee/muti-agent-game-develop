## RouletteCrabPanel.gd — 黃金輪盤螃蟹面板（DAY-167）
## 業界依據：King of Treasures Plus 2026「Roulette Crab — triggers Golden Roulette bonus game,
## player hits SHOOT to stop wheel, wins the amount listed where it stops.」
## 視覺設計：
##   - roulette_crab_start：全螢幕金色閃光 + 中央輪盤 UI（8格旋轉）+ 「點擊停止！」提示
##   - 旋轉動畫：輪盤持續旋轉，倒數計時 4 秒
##   - roulette_crab_result：輪盤停止 + 指針指向結果格 + 結果彈窗（倍率+獎勵）
##   - 自己觸發時：停止按鈕可點擊；旁觀者：只看動畫
##   - ≥100x：金色雙閃光；≥150x：彩虹三閃光
extends Node2D

# ---- 常數 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0
const WHEEL_SLOTS = [10.0, 20.0, 30.0, 50.0, 80.0, 100.0, 150.0, 200.0]
const WHEEL_COLORS = [
	Color(0.8, 0.8, 0.2),   # 10x 黃色
	Color(0.9, 0.6, 0.1),   # 20x 橙黃
	Color(1.0, 0.5, 0.0),   # 30x 橙色
	Color(1.0, 0.3, 0.0),   # 50x 橙紅
	Color(0.9, 0.1, 0.1),   # 80x 紅色
	Color(0.8, 0.0, 0.8),   # 100x 紫色
	Color(0.2, 0.5, 1.0),   # 150x 藍色
	Color(0.0, 0.9, 0.5),   # 200x 翠綠
]

# ---- 狀態 ----
var _pixel_font: Font = null
var _wheel_node: Node2D = null    # 輪盤節點
var _spin_tween: Tween = null     # 旋轉 tween
var _stop_btn: ColorRect = null   # 停止按鈕
var _countdown_lbl: Label = null  # 倒數計時
var _result_panel: Node2D = null  # 結果面板
var _is_my_wheel: bool = false    # 是否是自己的輪盤
var _spin_secs: float = 4.0       # 旋轉持續秒數
var _elapsed: float = 0.0         # 已過時間
var _is_spinning: bool = false    # 是否正在旋轉

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("roulette_crab_start"):
		GameManager.roulette_crab_start.connect(_on_roulette_crab_start)
	if GameManager.has_signal("roulette_crab_result"):
		GameManager.roulette_crab_result.connect(_on_roulette_crab_result)

# ---- 計時器 ----
func _process(delta: float) -> void:
	if not _is_spinning:
		return
	_elapsed += delta
	var remaining = _spin_secs - _elapsed
	if remaining < 0:
		remaining = 0.0
	if is_instance_valid(_countdown_lbl):
		_countdown_lbl.text = "%.1f" % remaining

# ---- 事件處理 ----

func _on_roulette_crab_start(data: Dictionary) -> void:
	var player_id: String = data.get("player_id", "")
	var player_name: String = data.get("player_name", "玩家")
	_is_my_wheel = (player_id == NetworkManager.get_player_id())
	_spin_secs = data.get("spin_secs", 4.0)
	_elapsed = 0.0
	_is_spinning = true

	# 全螢幕金色閃光
	var flash := ColorRect.new()
	flash.size = Vector2(SCREEN_W, SCREEN_H)
	flash.color = Color(1.0, 0.85, 0.0, 0.0)
	add_child(flash)
	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", 0.4, 0.12)
	flash_tween.tween_property(flash, "color:a", 0.0, 0.3)
	flash_tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())

	# 建立輪盤 UI
	_build_wheel_ui(player_name)

func _on_roulette_crab_result(data: Dictionary) -> void:
	_is_spinning = false
	var slot_index: int = data.get("slot_index", 0)
	var wheel_result: float = data.get("wheel_result", 10.0)
	var bonus_reward: int = data.get("bonus_reward", 0)
	var player_name: String = data.get("player_name", "玩家")
	var is_self: bool = (data.get("player_id", "") == NetworkManager.get_player_id())

	# 停止輪盤旋轉，指針指向結果格
	_stop_wheel_at_slot(slot_index)

	# 顯示結果
	await get_tree().create_timer(0.8).timeout
	_show_result(wheel_result, bonus_reward, player_name, is_self)

# ---- 輪盤 UI ----

func _build_wheel_ui(player_name: String) -> void:
	if is_instance_valid(_wheel_node):
		_wheel_node.queue_free()

	_wheel_node = Node2D.new()
	_wheel_node.position = Vector2(SCREEN_W / 2.0, SCREEN_H / 2.0)
	add_child(_wheel_node)

	# 輪盤背景（深色圓形）
	var bg := ColorRect.new()
	bg.size = Vector2(280, 280)
	bg.position = Vector2(-140, -140)
	bg.color = Color(0.05, 0.04, 0.02, 0.95)
	_wheel_node.add_child(bg)

	# 金色邊框
	var border := ColorRect.new()
	border.size = Vector2(284, 284)
	border.position = Vector2(-142, -142)
	border.color = Color(1.0, 0.85, 0.0, 0.8)
	border.z_index = -1
	_wheel_node.add_child(border)

	# 8 個格子（圍繞中心排列）
	var slot_size = Vector2(60, 30)
	for i in range(8):
		var angle = (float(i) / 8.0) * TAU - PI / 2.0
		var radius = 100.0
		var slot_x = cos(angle) * radius - slot_size.x / 2.0
		var slot_y = sin(angle) * radius - slot_size.y / 2.0

		var slot_bg := ColorRect.new()
		slot_bg.size = slot_size
		slot_bg.position = Vector2(slot_x, slot_y)
		slot_bg.color = WHEEL_COLORS[i]
		_wheel_node.add_child(slot_bg)

		var slot_lbl := Label.new()
		slot_lbl.text = "×%.0f" % WHEEL_SLOTS[i]
		slot_lbl.size = slot_size
		slot_lbl.position = Vector2(slot_x, slot_y)
		slot_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		slot_lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
		slot_lbl.add_theme_color_override("font_color", Color.WHITE)
		slot_lbl.add_theme_font_size_override("font_size", 11)
		if _pixel_font:
			slot_lbl.add_theme_font_override("font", _pixel_font)
		_wheel_node.add_child(slot_lbl)

	# 中心圖示
	var center_lbl := Label.new()
	center_lbl.text = "🦀"
	center_lbl.position = Vector2(-18, -22)
	center_lbl.add_theme_font_size_override("font_size", 36)
	_wheel_node.add_child(center_lbl)

	# 指針（頂部）
	var pointer := ColorRect.new()
	pointer.size = Vector2(8, 20)
	pointer.position = Vector2(-4, -150)
	pointer.color = Color(1.0, 0.2, 0.2)
	_wheel_node.add_child(pointer)

	# 旋轉動畫（持續旋轉）
	_spin_tween = _wheel_node.create_tween().set_loops()
	_spin_tween.tween_property(_wheel_node, "rotation_degrees", 360.0, 1.2)

	# 標題
	var title_lbl := Label.new()
	title_lbl.text = "🦀 %s 的黃金輪盤！" % player_name
	title_lbl.position = Vector2(-140, -175)
	title_lbl.size = Vector2(280, 24)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	title_lbl.add_theme_font_size_override("font_size", 13)
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
	_wheel_node.add_child(title_lbl)

	# 倒數計時
	_countdown_lbl = Label.new()
	_countdown_lbl.text = "%.1f" % _spin_secs
	_countdown_lbl.position = Vector2(-30, 150)
	_countdown_lbl.size = Vector2(60, 24)
	_countdown_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_countdown_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	_countdown_lbl.add_theme_font_size_override("font_size", 16)
	if _pixel_font:
		_countdown_lbl.add_theme_font_override("font", _pixel_font)
	_wheel_node.add_child(_countdown_lbl)

	# 停止按鈕（只有自己才能點擊）
	if _is_my_wheel:
		_stop_btn = ColorRect.new()
		_stop_btn.size = Vector2(100, 32)
		_stop_btn.position = Vector2(-50, 180)
		_stop_btn.color = Color(0.8, 0.1, 0.1, 0.95)
		_wheel_node.add_child(_stop_btn)

		var stop_lbl := Label.new()
		stop_lbl.text = "🛑 停止！"
		stop_lbl.size = Vector2(100, 32)
		stop_lbl.position = Vector2(-50, 180)
		stop_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		stop_lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
		stop_lbl.add_theme_color_override("font_color", Color.WHITE)
		stop_lbl.add_theme_font_size_override("font_size", 13)
		if _pixel_font:
			stop_lbl.add_theme_font_override("font", _pixel_font)
		_wheel_node.add_child(stop_lbl)

		# 停止按鈕點擊
		var area := Area2D.new()
		var col := CollisionShape2D.new()
		var shape := RectangleShape2D.new()
		shape.size = Vector2(100, 32)
		col.shape = shape
		col.position = Vector2(SCREEN_W / 2.0, SCREEN_H / 2.0 + 196)
		area.add_child(col)
		add_child(area)
		area.input_event.connect(func(_viewport, event, _shape_idx):
			if event is InputEventMouseButton and event.pressed and event.button_index == MOUSE_BUTTON_LEFT:
				if _is_spinning:
					NetworkManager.send_message("roulette_crab_stop", {})
					_is_spinning = false
		)
	else:
		# 旁觀者提示
		var watch_lbl := Label.new()
		watch_lbl.text = "旁觀中..."
		watch_lbl.position = Vector2(-40, 180)
		watch_lbl.size = Vector2(80, 24)
		watch_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		watch_lbl.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
		watch_lbl.add_theme_font_size_override("font_size", 11)
		if _pixel_font:
			watch_lbl.add_theme_font_override("font", _pixel_font)
		_wheel_node.add_child(watch_lbl)

# ---- 停止輪盤 ----

func _stop_wheel_at_slot(slot_index: int) -> void:
	if not is_instance_valid(_wheel_node):
		return

	# 停止旋轉 tween
	if is_instance_valid(_spin_tween):
		_spin_tween.kill()

	# 計算目標角度（讓指針指向對應格子）
	# 格子 i 的角度 = (i / 8) * 360 - 90 度（從頂部開始）
	# 指針在頂部（-90度），所以目標旋轉 = -(slot_angle)
	var slot_angle = (float(slot_index) / 8.0) * 360.0
	var target_rotation = -slot_angle

	# 緩停動畫（0.6 秒）
	var stop_tween = _wheel_node.create_tween()
	stop_tween.tween_property(_wheel_node, "rotation_degrees", target_rotation, 0.6)

	# 停止按鈕淡出
	if is_instance_valid(_stop_btn):
		var fade_tween = _stop_btn.create_tween()
		fade_tween.tween_property(_stop_btn, "modulate:a", 0.0, 0.3)

# ---- 結果顯示 ----

func _show_result(wheel_result: float, bonus_reward: int, player_name: String, is_self: bool) -> void:
	# 結果格子閃爍
	if is_instance_valid(_wheel_node):
		var flash_tween = _wheel_node.create_tween()
		flash_tween.tween_property(_wheel_node, "modulate", Color(2.0, 2.0, 0.5), 0.1)
		flash_tween.tween_property(_wheel_node, "modulate", Color(1.0, 1.0, 1.0), 0.1)
		flash_tween.tween_property(_wheel_node, "modulate", Color(2.0, 2.0, 0.5), 0.1)
		flash_tween.tween_property(_wheel_node, "modulate", Color(1.0, 1.0, 1.0), 0.1)

	# 結果彈窗
	if is_instance_valid(_result_panel):
		_result_panel.queue_free()

	_result_panel = Node2D.new()
	_result_panel.position = Vector2(SCREEN_W / 2.0 - 100, SCREEN_H / 2.0 + 160)
	add_child(_result_panel)

	var panel_bg := ColorRect.new()
	panel_bg.size = Vector2(200, 80)
	panel_bg.color = Color(0.05, 0.04, 0.02, 0.95)
	_result_panel.add_child(panel_bg)

	var border := ColorRect.new()
	border.size = Vector2(202, 82)
	border.position = Vector2(-1, -1)
	border.color = Color(1.0, 0.85, 0.0, 0.8)
	border.z_index = -1
	_result_panel.add_child(border)

	var result_lbl := Label.new()
	result_lbl.text = "🦀 ×%.0f！" % wheel_result
	result_lbl.position = Vector2(4, 6)
	result_lbl.size = Vector2(192, 28)
	result_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	result_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	result_lbl.add_theme_font_size_override("font_size", 20)
	if _pixel_font:
		result_lbl.add_theme_font_override("font", _pixel_font)
	_result_panel.add_child(result_lbl)

	var reward_lbl := Label.new()
	if is_self:
		reward_lbl.text = "+%d 金幣！" % bonus_reward
		reward_lbl.add_theme_color_override("font_color", Color(0.3, 1.0, 0.5))
	else:
		reward_lbl.text = "%s 獲得 %d 金幣" % [player_name, bonus_reward]
		reward_lbl.add_theme_color_override("font_color", Color(0.9, 0.8, 0.3))
	reward_lbl.position = Vector2(4, 40)
	reward_lbl.size = Vector2(192, 24)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_lbl.add_theme_font_size_override("font_size", 13)
	if _pixel_font:
		reward_lbl.add_theme_font_override("font", _pixel_font)
	_result_panel.add_child(reward_lbl)

	# 高倍率特效
	if wheel_result >= 150.0:
		_show_rainbow_flash()
	elif wheel_result >= 100.0:
		_show_double_flash()

	# 3 秒後清除所有 UI
	await get_tree().create_timer(3.0).timeout
	if is_instance_valid(_wheel_node):
		var fade_tween = _wheel_node.create_tween()
		fade_tween.tween_property(_wheel_node, "modulate:a", 0.0, 0.5)
		fade_tween.tween_callback(func(): if is_instance_valid(_wheel_node): _wheel_node.queue_free())
	if is_instance_valid(_result_panel):
		var fade_tween2 = _result_panel.create_tween()
		fade_tween2.tween_property(_result_panel, "modulate:a", 0.0, 0.5)
		fade_tween2.tween_callback(func(): if is_instance_valid(_result_panel): _result_panel.queue_free())

# ---- 特效 ----

func _show_double_flash() -> void:
	for i in range(2):
		var flash := ColorRect.new()
		flash.size = Vector2(SCREEN_W, SCREEN_H)
		flash.color = Color(1.0, 0.85, 0.0, 0.0)
		add_child(flash)
		var tween = flash.create_tween()
		tween.tween_interval(float(i) * 0.2)
		tween.tween_property(flash, "color:a", 0.35, 0.1)
		tween.tween_property(flash, "color:a", 0.0, 0.2)
		tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())

func _show_rainbow_flash() -> void:
	var colors = [Color(1.0, 0.0, 0.0), Color(0.0, 1.0, 0.0), Color(0.0, 0.0, 1.0)]
	for i in range(3):
		var flash := ColorRect.new()
		flash.size = Vector2(SCREEN_W, SCREEN_H)
		flash.color = Color(colors[i].r, colors[i].g, colors[i].b, 0.0)
		add_child(flash)
		var tween = flash.create_tween()
		tween.tween_interval(float(i) * 0.15)
		tween.tween_property(flash, "color:a", 0.3, 0.08)
		tween.tween_property(flash, "color:a", 0.0, 0.15)
		tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())
