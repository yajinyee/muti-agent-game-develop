## CursedPoisonFishPanel.gd — 詛咒毒魚系統面板（DAY-216）
## 業界原創「詛咒反轉」機制
##
## 視覺設計：
##   - 紫色詛咒主題（#9B59B6 + #6C3483 + #E8D5FF + #FF4444）
##   - curse_start：紫色雙閃光 + 頂部橫幅 + 詛咒目標紫色骷髏標記
##   - curse_kill：詛咒消失閃光 + 「☠️ ×2.5 詛咒加成！」浮動文字
##   - curse_escape：紅色警告閃光 + 「⚠️ 詛咒懲罰！×0.5 持續 5 秒」大字
##   - curse_cleanse：白色解咒閃光 + 「✨ 詛咒解除！」大字 + 解咒獎勵
##   - curse_end：標記淡出
extends CanvasLayer

# 詛咒目標視覺節點（instanceID → Control）
var _curse_nodes: Dictionary = {}
var _penalty_overlay: Control = null

func _ready() -> void:
	layer = 29  # 詛咒毒魚面板層級

## 處理詛咒毒魚訊息
func handle_cursed_poison_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"curse_start":
			_on_curse_start(payload)
		"curse_kill":
			_on_curse_kill(payload)
		"curse_escape":
			_on_curse_escape(payload)
		"curse_cleanse":
			_on_curse_cleanse(payload)

## 詛咒標記建立 — 紫色雙閃光 + 頂部橫幅 + 骷髏標記
func _on_curse_start(payload: Dictionary) -> void:
	var cursed_targets: Array = payload.get("cursed_targets", [])
	var curse_mult: float = payload.get("curse_mult", 2.5)
	var player_name: String = payload.get("player_name", "")

	# 紫色雙閃光
	_double_flash(Color("#9B59B6"), 0.50)

	# 頂部橫幅
	var msg: String
	if player_name.is_empty():
		msg = "☠️ 詛咒毒魚出現！%d 個目標被詛咒（×%.1f 倍率）！讓它跑掉受懲罰！" % [cursed_targets.size(), curse_mult]
	else:
		msg = "☠️ %s 觸發詛咒毒魚！%d 個目標被詛咒（×%.1f 倍率）！" % [player_name, cursed_targets.size(), curse_mult]
	var banner = _make_banner(msg, Color(0.06, 0.0, 0.1, 0.88), Color("#E8D5FF"))
	add_child(banner)
	var tw = create_tween()
	tw.tween_interval(4.0)
	tw.tween_callback(func(): if is_instance_valid(banner): banner.queue_free())

	# 建立詛咒骷髏標記
	for target_info in cursed_targets:
		_spawn_curse_marker(target_info)

## 詛咒目標被擊破 — 消失閃光 + 浮動文字
func _on_curse_kill(payload: Dictionary) -> void:
	var instance_id: String = payload.get("instance_id", "")

	# 移除詛咒標記
	if _curse_nodes.has(instance_id):
		var node = _curse_nodes[instance_id]
		if is_instance_valid(node):
			var tw = node.create_tween()
			tw.tween_property(node, "modulate:a", 0.0, 0.2)
			tw.tween_callback(func(): if is_instance_valid(node): node.queue_free())
		_curse_nodes.erase(instance_id)

	# 浮動文字
	_spawn_float_text("☠️ ×2.5 詛咒加成！", Color("#E8D5FF"), 36)

## 詛咒目標逃跑 — 紅色警告閃光 + 懲罰提示
func _on_curse_escape(payload: Dictionary) -> void:
	var penalty_sec: int = payload.get("penalty_sec", 5)
	var instance_id: String = payload.get("instance_id", "")

	# 移除詛咒標記
	if _curse_nodes.has(instance_id):
		var node = _curse_nodes[instance_id]
		if is_instance_valid(node):
			node.queue_free()
		_curse_nodes.erase(instance_id)

	# 紅色警告閃光
	_double_flash(Color("#FF4444"), 0.65)

	# 懲罰大字
	var vp = get_viewport()
	var vp_size = vp.get_visible_rect().size if vp else Vector2(1280, 720)
	var label = Label.new()
	label.text = "⚠️ 詛咒懲罰！×0.5 持續 %d 秒！" % penalty_sec
	label.add_theme_color_override("font_color", Color("#FF4444"))
	label.add_theme_font_size_override("font_size", 40)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.size = Vector2(vp_size.x, 56)
	label.position = Vector2(0, vp_size.y * 0.42)
	add_child(label)

	# 懲罰計時條（底部紅色）
	_penalty_overlay = _make_penalty_bar(penalty_sec)
	add_child(_penalty_overlay)

	var tw = create_tween()
	tw.tween_interval(2.5)
	tw.tween_callback(func(): if is_instance_valid(label): label.queue_free())

## 解除詛咒 — 白色解咒閃光 + 大字 + 解咒獎勵
func _on_curse_cleanse(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var cleanse_reward: int = payload.get("cleanse_reward", 0)
	var killed_count: int = payload.get("killed_count", 0)

	# 清除所有詛咒標記
	_clear_all_curses()

	# 清除懲罰條
	if is_instance_valid(_penalty_overlay):
		_penalty_overlay.queue_free()
		_penalty_overlay = null

	# 白色解咒閃光
	_triple_flash(Color("#FFFFFF"), 0.60)

	# 解咒大字
	var vp = get_viewport()
	var vp_size = vp.get_visible_rect().size if vp else Vector2(1280, 720)
	var label = Label.new()
	label.text = "✨ %s 解除詛咒！" % player_name
	label.add_theme_color_override("font_color", Color("#E8D5FF"))
	label.add_theme_font_size_override("font_size", 44)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.size = Vector2(vp_size.x, 60)
	label.position = Vector2(0, vp_size.y * 0.38)
	add_child(label)

	var sub = Label.new()
	sub.text = "解除 %d 個詛咒！獲得 +%d 金幣！" % [killed_count, cleanse_reward]
	sub.add_theme_color_override("font_color", Color("#FFD700"))
	sub.add_theme_font_size_override("font_size", 28)
	sub.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	sub.size = Vector2(vp_size.x, 40)
	sub.position = Vector2(0, vp_size.y * 0.38 + 68)
	add_child(sub)

	var tw = create_tween()
	tw.tween_interval(3.0)
	tw.tween_callback(func():
		if is_instance_valid(label): label.queue_free()
		if is_instance_valid(sub): sub.queue_free()
	)

# ─── 內部工具函數 ───────────────────────────────────────────────────────────

## 建立詛咒骷髏標記（紫色骷髏 + 閃爍）
func _spawn_curse_marker(target_info: Dictionary) -> void:
	# 詛咒標記沒有位置資訊（Server 不廣播目標位置），
	# 改為在畫面角落顯示詛咒計數器
	var instance_id: String = target_info.get("instance_id", "")
	var curse_mult: float = target_info.get("curse_mult", 2.5)

	# 建立小型詛咒標記（右側垂直排列）
	var vp = get_viewport()
	var vp_size = vp.get_visible_rect().size if vp else Vector2(1280, 720)
	var idx = _curse_nodes.size()

	var marker = ColorRect.new()
	marker.color = Color(0.1, 0.0, 0.15, 0.85)
	marker.size = Vector2(120, 36)
	marker.position = Vector2(vp_size.x - 130, 80 + idx * 44)

	var lbl = Label.new()
	lbl.text = "☠️ ×%.1f" % curse_mult
	lbl.add_theme_color_override("font_color", Color("#E8D5FF"))
	lbl.add_theme_font_size_override("font_size", 18)
	lbl.position = Vector2(8, 8)
	marker.add_child(lbl)

	# 閃爍動畫
	var tw = marker.create_tween().set_loops()
	tw.tween_property(marker, "modulate:a", 0.5, 0.6)
	tw.tween_property(marker, "modulate:a", 1.0, 0.6)

	add_child(marker)
	_curse_nodes[instance_id] = marker

## 清除所有詛咒標記
func _clear_all_curses() -> void:
	for id in _curse_nodes.keys():
		var node = _curse_nodes[id]
		if is_instance_valid(node):
			node.queue_free()
	_curse_nodes.clear()

## 建立懲罰計時條（底部紅色）
func _make_penalty_bar(duration_sec: int) -> Control:
	var vp = get_viewport()
	var vp_size = vp.get_visible_rect().size if vp else Vector2(1280, 720)

	var container = Control.new()
	container.size = Vector2(vp_size.x, 12)
	container.position = Vector2(0, vp_size.y - 12)

	var bg = ColorRect.new()
	bg.color = Color(0.3, 0.0, 0.0, 0.8)
	bg.size = Vector2(vp_size.x, 12)
	container.add_child(bg)

	var bar = ColorRect.new()
	bar.color = Color("#FF4444")
	bar.size = Vector2(vp_size.x, 12)
	container.add_child(bar)

	# 計時條縮短動畫
	var tw = container.create_tween()
	tw.tween_property(bar, "size:x", 0.0, float(duration_sec))
	tw.tween_callback(func(): if is_instance_valid(container): container.queue_free())

	return container

## 紫色雙閃光
func _double_flash(color: Color, alpha: float) -> void:
	var vp = get_viewport()
	var vp_size = vp.get_visible_rect().size if vp else Vector2(1280, 720)
	for i in range(2):
		var flash = ColorRect.new()
		flash.color = Color(color.r, color.g, color.b, alpha)
		flash.size = vp_size
		add_child(flash)
		var tw = flash.create_tween()
		tw.tween_interval(float(i) * 0.18)
		tw.tween_property(flash, "modulate:a", 0.0, 0.22)
		tw.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())

## 全螢幕三次強閃光
func _triple_flash(color: Color, alpha: float) -> void:
	var vp = get_viewport()
	var vp_size = vp.get_visible_rect().size if vp else Vector2(1280, 720)
	for i in range(3):
		var flash = ColorRect.new()
		flash.color = Color(color.r, color.g, color.b, alpha)
		flash.size = vp_size
		add_child(flash)
		var tw = flash.create_tween()
		tw.tween_interval(float(i) * 0.22)
		tw.tween_property(flash, "modulate:a", 0.0, 0.28)
		tw.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())

## 建立頂部橫幅
func _make_banner(text: String, bg_color: Color, text_color: Color) -> Control:
	var vp = get_viewport()
	var vp_size = vp.get_visible_rect().size if vp else Vector2(1280, 720)
	var banner = ColorRect.new()
	banner.color = bg_color
	banner.size = Vector2(vp_size.x, 52)
	banner.position = Vector2(0, 8)

	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", text_color)
	label.add_theme_font_size_override("font_size", 17)
	label.position = Vector2(16, 12)
	label.size = Vector2(vp_size.x - 32, 32)
	banner.add_child(label)
	return banner

## 浮動文字（置中偏上）
func _spawn_float_text(text: String, color: Color, font_size: int) -> void:
	var vp = get_viewport()
	var vp_size = vp.get_visible_rect().size if vp else Vector2(1280, 720)
	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", font_size)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.size = Vector2(vp_size.x, font_size + 16)
	label.position = Vector2(0, vp_size.y * 0.3)
	add_child(label)
	var tw = create_tween()
	tw.tween_property(label, "position:y", label.position.y - 60.0, 1.2)
	tw.parallel().tween_property(label, "modulate:a", 0.0, 1.2)
	tw.tween_callback(func(): if is_instance_valid(label): label.queue_free())
