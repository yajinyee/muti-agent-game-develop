## FreezeBombPanel.gd — 冰凍炸彈魚面板（DAY-170）
## 業界依據：King of Ocean 2026「The freezing blast pauses an entire school for a few seconds —
## useful when a high-tier creature is escaping the frame.」
## 視覺設計：
##   - freeze_start：全螢幕冰藍閃光 + 頂部橫幅滑入 + 特殊目標變冰藍色 + 倒數計時 6 秒
##   - 自己觸發時：中央大 ❄️ 標誌彈跳動畫 + 「特殊目標已冰凍！快去擊破！」提示
##   - 冰凍期間：特殊目標顯示冰晶光暈（藍白色）+ 倒數計時
##   - freeze_end：冰晶碎裂動畫 + 淡出所有 UI
##   - ≥3個：雙閃光；≥5個：彩虹三閃光
extends Node2D

# ---- 常數 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0

# ---- 狀態 ----
var _pixel_font: Font = null
var _banner: Node2D = null          # 頂部橫幅
var _countdown_lbl: Label = null    # 倒數計時
var _frozen_nodes: Dictionary = {}  # instanceID -> Node2D（冰晶光暈）
var _is_my_freeze: bool = false     # 是否是自己觸發的冰凍
var _duration_sec: int = 6          # 冰凍持續時間
var _elapsed: float = 0.0           # 已過時間
var _is_active: bool = false        # 是否正在冰凍中
var _frozen_count: int = 0          # 被冰凍的目標數

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("freeze_bomb"):
		GameManager.freeze_bomb.connect(_on_freeze_bomb)
	# 追蹤目標移動，更新冰晶光暈位置
	if GameManager.has_signal("target_updated"):
		GameManager.target_updated.connect(_on_target_updated)
	# 目標被擊破時移除冰晶光暈
	if GameManager.has_signal("target_killed"):
		GameManager.target_killed.connect(_on_target_killed)

# ---- 計時器 ----
func _process(delta: float) -> void:
	if not _is_active:
		return
	_elapsed += delta
	var remaining = float(_duration_sec) - _elapsed
	if remaining < 0.0:
		remaining = 0.0
	if is_instance_valid(_countdown_lbl):
		_countdown_lbl.text = "❄️ %.0f秒" % remaining

# ---- 目標位置追蹤 ----
func _on_target_updated(data: Dictionary) -> void:
	if not _is_active:
		return
	var instance_id: String = data.get("instance_id", "")
	if not _frozen_nodes.has(instance_id):
		return
	var node = _frozen_nodes[instance_id]
	if not is_instance_valid(node):
		_frozen_nodes.erase(instance_id)
		return
	var x: float = data.get("x", node.position.x)
	var y: float = data.get("y", node.position.y)
	node.position = Vector2(x, y)

func _on_target_killed(data: Dictionary) -> void:
	if not _is_active:
		return
	var instance_id: String = data.get("instance_id", "")
	if not _frozen_nodes.has(instance_id):
		return
	var node = _frozen_nodes[instance_id]
	_frozen_nodes.erase(instance_id)
	# 冰晶碎裂動畫
	if is_instance_valid(node):
		var tween = create_tween()
		tween.tween_property(node, "scale", Vector2(1.5, 1.5), 0.1)
		tween.tween_property(node, "modulate:a", 0.0, 0.15)
		tween.tween_callback(func(): if is_instance_valid(node): node.queue_free())

# ---- 主要事件處理 ----
func _on_freeze_bomb(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var trigger_id: String = data.get("trigger_id", "")
	var trigger_name: String = data.get("trigger_name", "冰凍炸彈魚")
	var freeze_x: float = data.get("freeze_x", SCREEN_W / 2.0)
	var freeze_y: float = data.get("freeze_y", SCREEN_H / 2.0)
	var frozen_count: int = data.get("frozen_count", 0)
	var duration_sec: int = data.get("duration_sec", 6)

	match phase:
		"freeze_start":
			var frozen_targets = data.get("frozen_targets", [])
			_start_freeze(trigger_id, trigger_name, freeze_x, freeze_y, frozen_count, duration_sec, frozen_targets)
		"freeze_end":
			_end_freeze(frozen_count)

# ---- 冰凍開始 ----
func _start_freeze(trigger_id: String, trigger_name: String, fx: float, fy: float,
		frozen_count: int, duration_sec: int, frozen_targets: Array) -> void:
	_is_active = true
	_duration_sec = duration_sec
	_elapsed = 0.0
	_frozen_count = frozen_count

	# 判斷是否是自己觸發
	var my_id: String = ""
	if GameManager.has_method("get_player_id"):
		my_id = GameManager.get_player_id()
	_is_my_freeze = (trigger_id == my_id)

	# 全螢幕冰藍閃光
	_flash_screen(Color(0.0, 0.8, 1.0, 0.55), 0.4)

	# 建立頂部橫幅
	_create_banner(trigger_name, frozen_count)

	# 建立倒數計時
	_create_countdown()

	# 為每個被冰凍的目標建立冰晶光暈
	for entry in frozen_targets:
		var instance_id: String = entry.get("instance_id", "")
		var ex: float = entry.get("x", fx)
		var ey: float = entry.get("y", fy)
		if instance_id != "":
			_create_ice_halo(instance_id, ex, ey)

	# 自己觸發時：中央大 ❄️ 標誌彈跳
	if _is_my_freeze:
		_show_my_trigger_anim()

	# 多閃光效果
	if frozen_count >= 5:
		await get_tree().create_timer(0.5).timeout
		_flash_screen(Color(0.5, 0.9, 1.0, 0.4), 0.2)
		await get_tree().create_timer(0.25).timeout
		_flash_screen(Color(0.8, 0.95, 1.0, 0.3), 0.2)
	elif frozen_count >= 3:
		await get_tree().create_timer(0.5).timeout
		_flash_screen(Color(0.3, 0.8, 1.0, 0.35), 0.2)

# ---- 建立頂部橫幅 ----
func _create_banner(trigger_name: String, frozen_count: int) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()

	_banner = Node2D.new()
	_banner.position = Vector2(SCREEN_W / 2.0, -60)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.size = Vector2(580, 52)
	bg.position = Vector2(-290, -26)
	bg.color = Color(0.0, 0.15, 0.4, 0.88)
	_banner.add_child(bg)

	var lbl = Label.new()
	lbl.text = "❄️ %s 觸發冰凍炸彈魚！%d 個特殊目標被冰凍！" % [trigger_name, frozen_count]
	lbl.position = Vector2(-275, -18)
	lbl.add_theme_color_override("font_color", Color(0.6, 0.95, 1.0))
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
	lbl.add_theme_font_size_override("font_size", 16)
	_banner.add_child(lbl)

	var tween = create_tween()
	tween.tween_property(_banner, "position:y", 36.0, 0.3).set_ease(Tween.EASE_OUT)

# ---- 建立倒數計時 ----
func _create_countdown() -> void:
	if is_instance_valid(_countdown_lbl):
		_countdown_lbl.queue_free()

	_countdown_lbl = Label.new()
	_countdown_lbl.text = "❄️ %d秒" % _duration_sec
	_countdown_lbl.position = Vector2(SCREEN_W - 120, 60)
	_countdown_lbl.add_theme_color_override("font_color", Color(0.6, 0.95, 1.0))
	if _pixel_font:
		_countdown_lbl.add_theme_font_override("font", _pixel_font)
	_countdown_lbl.add_theme_font_size_override("font_size", 20)
	add_child(_countdown_lbl)

# ---- 建立冰晶光暈 ----
func _create_ice_halo(instance_id: String, x: float, y: float) -> void:
	var halo = Node2D.new()
	halo.position = Vector2(x, y)
	add_child(halo)

	# 冰晶外圈
	var ring = ColorRect.new()
	ring.size = Vector2(56, 56)
	ring.position = Vector2(-28, -28)
	ring.color = Color(0.4, 0.85, 1.0, 0.35)
	halo.add_child(ring)

	# 冰晶圖示
	var icon_lbl = Label.new()
	icon_lbl.text = "❄️"
	icon_lbl.position = Vector2(-10, -10)
	if _pixel_font:
		icon_lbl.add_theme_font_override("font", _pixel_font)
	icon_lbl.add_theme_font_size_override("font_size", 18)
	halo.add_child(icon_lbl)

	# 閃爍動畫
	var tween = halo.create_tween().set_loops()
	tween.tween_property(halo, "modulate:a", 0.5, 0.6)
	tween.tween_property(halo, "modulate:a", 1.0, 0.6)

	_frozen_nodes[instance_id] = halo

# ---- 自己觸發動畫 ----
func _show_my_trigger_anim() -> void:
	var anim_node = Node2D.new()
	anim_node.position = Vector2(SCREEN_W / 2.0, SCREEN_H / 2.0)
	add_child(anim_node)

	var lbl = Label.new()
	lbl.text = "❄️"
	lbl.position = Vector2(-24, -24)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
	lbl.add_theme_font_size_override("font_size", 48)
	anim_node.add_child(lbl)

	var sub_lbl = Label.new()
	sub_lbl.text = "特殊目標已冰凍！快去擊破！"
	sub_lbl.position = Vector2(-80, 30)
	sub_lbl.add_theme_color_override("font_color", Color(0.6, 0.95, 1.0))
	if _pixel_font:
		sub_lbl.add_theme_font_override("font", _pixel_font)
	sub_lbl.add_theme_font_size_override("font_size", 14)
	anim_node.add_child(sub_lbl)

	var tween = create_tween()
	tween.tween_property(anim_node, "scale", Vector2(1.4, 1.4), 0.15).set_ease(Tween.EASE_OUT)
	tween.tween_property(anim_node, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(1.2)
	tween.tween_property(anim_node, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func(): if is_instance_valid(anim_node): anim_node.queue_free())

# ---- 冰凍結束 ----
func _end_freeze(frozen_count: int) -> void:
	_is_active = false

	# 清理所有冰晶光暈（碎裂動畫）
	for instance_id in _frozen_nodes.keys():
		var node = _frozen_nodes[instance_id]
		if is_instance_valid(node):
			var tween = create_tween()
			tween.tween_property(node, "scale", Vector2(1.3, 1.3), 0.1)
			tween.tween_property(node, "modulate:a", 0.0, 0.2)
			tween.tween_callback(func(): if is_instance_valid(node): node.queue_free())
	_frozen_nodes.clear()

	# 清理橫幅
	if is_instance_valid(_banner):
		var tween2 = create_tween()
		tween2.tween_property(_banner, "modulate:a", 0.0, 0.3)
		tween2.tween_callback(func(): if is_instance_valid(_banner): _banner.queue_free())

	# 清理倒數計時
	if is_instance_valid(_countdown_lbl):
		var tween3 = create_tween()
		tween3.tween_property(_countdown_lbl, "modulate:a", 0.0, 0.3)
		tween3.tween_callback(func(): if is_instance_valid(_countdown_lbl): _countdown_lbl.queue_free())

	# 冰晶碎裂閃光
	_flash_screen(Color(0.7, 0.95, 1.0, 0.3), 0.25)

# ---- 全螢幕閃光 ----
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.size = Vector2(SCREEN_W, SCREEN_H)
	flash.color = color
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())
