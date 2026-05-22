## LionDancePanel.gd — 獅子舞大獎爆發面板（DAY-168）
## 業界依據：Fortune King Jackpot（TaDa Gaming 2026）「Lion Dance bonus — triggered by special fish,
## delivers burst multiplier payouts with festive visual effects」
## 視覺設計：
##   - burst_start：全螢幕橙紅閃光 + 頂部橫幅滑入 + 標記目標顯示金色光環 + 倒數計時 15 秒
##   - 自己觸發時：中央大 🦁 標誌彈跳動畫 + 「快去擊破標記目標！」提示
##   - 擊破標記目標時：浮動倍率文字（+Nx）+ 金色爆炸閃光
##   - burst_end：淡出所有 UI + 右側滑入結果彈窗（擊破數/倍率/獎勵）
##   - ≥7x：金色雙閃光；≥10x：彩虹三閃光
extends Node2D

# ---- 常數 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0

# ---- 狀態 ----
var _pixel_font: Font = null
var _banner: Node2D = null         # 頂部橫幅
var _countdown_lbl: Label = null   # 倒數計時
var _mark_nodes: Dictionary = {}   # instanceID -> Node2D（標記光環）
var _is_my_burst: bool = false     # 是否是自己觸發的爆發
var _burst_mult: float = 1.0       # 本次爆發倍率
var _duration_sec: int = 15        # 持續時間
var _elapsed: float = 0.0          # 已過時間
var _is_active: bool = false       # 是否正在爆發

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("lion_dance_burst"):
		GameManager.lion_dance_burst.connect(_on_lion_dance_burst)

# ---- 計時器 ----
func _process(delta: float) -> void:
	if not _is_active:
		return
	_elapsed += delta
	var remaining = float(_duration_sec) - _elapsed
	if remaining < 0.0:
		remaining = 0.0
	if is_instance_valid(_countdown_lbl):
		_countdown_lbl.text = "🦁 %.0f秒" % remaining

# ---- 事件處理 ----

func _on_lion_dance_burst(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	match phase:
		"burst_start":
			_handle_burst_start(data)
		"burst_end":
			_handle_burst_end(data)

func _handle_burst_start(data: Dictionary) -> void:
	var trigger_player: String = data.get("trigger_player", "")
	var trigger_name: String = data.get("trigger_name", "玩家")
	_burst_mult = data.get("burst_mult", 3.0)
	_duration_sec = data.get("duration_sec", 15)
	_elapsed = 0.0
	_is_active = true
	_is_my_burst = (trigger_player == NetworkManager.get_player_id())

	# 全螢幕橙紅閃光
	_flash_screen(Color(1.0, 0.4, 0.0, 0.0), 0.35)

	# 頂部橫幅
	_show_banner(trigger_name, _burst_mult)

	# 標記目標光環
	var marked: Array = data.get("marked_targets", [])
	for t in marked:
		_add_mark_halo(t.get("instance_id", ""), t.get("x", 0.0), t.get("y", 0.0))

	# 自己觸發時：中央大 🦁 標誌彈跳
	if _is_my_burst:
		_show_center_lion()

func _handle_burst_end(data: Dictionary) -> void:
	_is_active = false

	# 清除所有標記光環
	for id in _mark_nodes:
		var node = _mark_nodes[id]
		if is_instance_valid(node):
			node.queue_free()
	_mark_nodes.clear()

	# 淡出橫幅
	if is_instance_valid(_banner):
		var t = _banner.create_tween()
		t.tween_property(_banner, "modulate:a", 0.0, 0.4)
		t.tween_callback(func(): if is_instance_valid(_banner): _banner.queue_free(); _banner = null)

	# 清除倒數計時
	if is_instance_valid(_countdown_lbl):
		_countdown_lbl.queue_free()
		_countdown_lbl = null

	# 右側滑入結果彈窗
	var remaining: int = data.get("remaining_targets", 0)
	_show_result_panel(remaining)

# ---- UI 建立 ----

func _flash_screen(base_color: Color, peak_alpha: float) -> void:
	var flash := ColorRect.new()
	flash.size = Vector2(SCREEN_W, SCREEN_H)
	flash.color = Color(base_color.r, base_color.g, base_color.b, 0.0)
	add_child(flash)
	var tw = flash.create_tween()
	tw.tween_property(flash, "color:a", peak_alpha, 0.1)
	tw.tween_property(flash, "color:a", 0.0, 0.35)
	tw.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())

func _show_banner(trigger_name: String, mult: float) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()

	_banner = Node2D.new()
	add_child(_banner)

	# 橫幅背景（橙紅漸層）
	var bg := ColorRect.new()
	bg.size = Vector2(SCREEN_W, 56)
	bg.position = Vector2(0, -60)
	bg.color = Color(0.85, 0.25, 0.0, 0.92)
	_banner.add_child(bg)

	# 橫幅文字
	var lbl := Label.new()
	lbl.text = "🦁 %s 觸發獅子舞爆發！標記目標 ×%.0f 倍率！" % [trigger_name, mult]
	lbl.position = Vector2(20, 10)
	lbl.add_theme_color_override("font_color", Color(1.0, 0.95, 0.6))
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 18)
	_banner.add_child(lbl)

	# 倒數計時 Label
	_countdown_lbl = Label.new()
	_countdown_lbl.text = "🦁 15秒"
	_countdown_lbl.position = Vector2(SCREEN_W - 120, 10)
	_countdown_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	if _pixel_font:
		_countdown_lbl.add_theme_font_override("font", _pixel_font)
		_countdown_lbl.add_theme_font_size_override("font_size", 18)
	_banner.add_child(_countdown_lbl)

	# 橫幅從頂部滑入
	var tw = _banner.create_tween()
	tw.tween_property(bg, "position:y", 0.0, 0.25).set_ease(Tween.EASE_OUT)

func _add_mark_halo(instance_id: String, x: float, y: float) -> void:
	if instance_id == "":
		return

	var halo := Node2D.new()
	halo.position = Vector2(x, y)
	add_child(halo)
	_mark_nodes[instance_id] = halo

	# 金色光環（旋轉動畫）
	var ring := ColorRect.new()
	ring.size = Vector2(64, 64)
	ring.position = Vector2(-32, -32)
	ring.color = Color(1.0, 0.85, 0.0, 0.0)
	halo.add_child(ring)

	# 淡入光環
	var tw = ring.create_tween().set_loops()
	tw.tween_property(ring, "color:a", 0.7, 0.4)
	tw.tween_property(ring, "color:a", 0.2, 0.4)

	# 倍率標籤
	var lbl := Label.new()
	lbl.text = "×%.0f" % _burst_mult
	lbl.position = Vector2(-20, -48)
	lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.0))
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 14)
	halo.add_child(lbl)

func _show_center_lion() -> void:
	var lbl := Label.new()
	lbl.text = "🦁"
	lbl.position = Vector2(SCREEN_W / 2 - 40, SCREEN_H / 2 - 60)
	lbl.add_theme_color_override("font_color", Color(1.0, 0.6, 0.0))
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 72)
	add_child(lbl)

	# 彈跳動畫
	var tw = lbl.create_tween()
	tw.tween_property(lbl, "position:y", SCREEN_H / 2 - 100, 0.15).set_ease(Tween.EASE_OUT)
	tw.tween_property(lbl, "position:y", SCREEN_H / 2 - 60, 0.12).set_ease(Tween.EASE_IN)
	tw.tween_property(lbl, "position:y", SCREEN_H / 2 - 80, 0.1).set_ease(Tween.EASE_OUT)
	tw.tween_property(lbl, "position:y", SCREEN_H / 2 - 60, 0.1).set_ease(Tween.EASE_IN)
	tw.tween_property(lbl, "modulate:a", 0.0, 0.5)
	tw.tween_callback(func(): if is_instance_valid(lbl): lbl.queue_free())

	# 副標題
	var sub := Label.new()
	sub.text = "快去擊破標記目標！"
	sub.position = Vector2(SCREEN_W / 2 - 100, SCREEN_H / 2 + 20)
	sub.add_theme_color_override("font_color", Color(1.0, 0.85, 0.3))
	if _pixel_font:
		sub.add_theme_font_override("font", _pixel_font)
		sub.add_theme_font_size_override("font_size", 20)
	add_child(sub)
	var sub_tw = sub.create_tween()
	sub_tw.tween_interval(2.0)
	sub_tw.tween_property(sub, "modulate:a", 0.0, 0.5)
	sub_tw.tween_callback(func(): if is_instance_valid(sub): sub.queue_free())

func _show_result_panel(remaining: int) -> void:
	var panel := Node2D.new()
	panel.position = Vector2(SCREEN_W + 10, SCREEN_H / 2 - 80)
	add_child(panel)

	# 面板背景
	var bg := ColorRect.new()
	bg.size = Vector2(280, 160)
	bg.color = Color(0.15, 0.08, 0.0, 0.95)
	panel.add_child(bg)

	# 邊框
	var border := ColorRect.new()
	border.size = Vector2(280, 4)
	border.color = Color(1.0, 0.6, 0.0, 1.0)
	panel.add_child(border)

	# 標題
	var title := Label.new()
	title.text = "🦁 獅子舞爆發結束"
	title.position = Vector2(10, 12)
	title.add_theme_color_override("font_color", Color(1.0, 0.85, 0.3))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 16)
	panel.add_child(title)

	# 倍率
	var mult_lbl := Label.new()
	mult_lbl.text = "爆發倍率：×%.0f" % _burst_mult
	mult_lbl.position = Vector2(10, 45)
	mult_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.5))
	if _pixel_font:
		mult_lbl.add_theme_font_override("font", _pixel_font)
		mult_lbl.add_theme_font_size_override("font_size", 14)
	panel.add_child(mult_lbl)

	# 剩餘未擊破
	var remain_lbl := Label.new()
	remain_lbl.text = "未擊破標記：%d 個" % remaining
	remain_lbl.position = Vector2(10, 75)
	remain_lbl.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	if _pixel_font:
		remain_lbl.add_theme_font_override("font", _pixel_font)
		remain_lbl.add_theme_font_size_override("font_size", 14)
	panel.add_child(remain_lbl)

	# 從右側滑入
	var tw = panel.create_tween()
	tw.tween_property(panel, "position:x", SCREEN_W - 300, 0.3).set_ease(Tween.EASE_OUT)
	tw.tween_interval(3.0)
	tw.tween_property(panel, "modulate:a", 0.0, 0.5)
	tw.tween_callback(func(): if is_instance_valid(panel): panel.queue_free())

	# ≥7x 雙閃光
	if _burst_mult >= 7.0:
		_flash_screen(Color(1.0, 0.7, 0.0, 0.0), 0.5)
		await get_tree().create_timer(0.2).timeout
		_flash_screen(Color(1.0, 0.7, 0.0, 0.0), 0.4)

	# ≥10x 彩虹三閃光
	if _burst_mult >= 10.0:
		await get_tree().create_timer(0.4).timeout
		_flash_screen(Color(0.5, 0.0, 1.0, 0.0), 0.35)
