## LuckyMirrorFishPanel.gd — 幸運鏡像魚系統面板（DAY-215）
## 業界原創「鏡像複製」機制
##
## 視覺設計：
##   - 青色鏡像主題（#00FFFF + #00CCFF + #88FFFF + #0088FF）
##   - mirror_start：青色雙閃光 + 頂部橫幅 + 鏡像分身標記（菱形輪廓）
##   - mirror_kill：鏡像分身消失閃光 + 「×1.5 鏡像加成！」浮動文字
##   - mirror_blast：全螢幕青色三次強閃光 + 「🪞 鏡像爆炸！」大字
##   - mirror_result：右側滑入結算彈窗
extends CanvasLayer

# 鏡像分身視覺節點（mirrorID → Control）
var _mirror_nodes: Dictionary = {}

func _ready() -> void:
	layer = 30  # 幸運鏡像魚面板層級

## 處理幸運鏡像魚訊息
func handle_lucky_mirror_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"mirror_start":
			_on_mirror_start(payload)
		"mirror_kill":
			_on_mirror_kill(payload)
		"mirror_blast":
			_on_mirror_blast(payload)
		"mirror_result":
			_on_mirror_result(payload)

## 鏡像複製開始 — 青色雙閃光 + 頂部橫幅 + 鏡像分身標記
func _on_mirror_start(payload: Dictionary) -> void:
	var mirrors: Array = payload.get("mirrors", [])
	var mult_boost: float = payload.get("mult_boost", 1.5)
	var player_name: String = payload.get("player_name", "")

	# 青色雙閃光
	_double_flash(Color("#00FFFF"), 0.50)

	# 頂部橫幅
	var msg = "🪞 %s 觸發幸運鏡像魚！%d 個鏡像分身出現（×%.1f 倍率）！" % [player_name, mirrors.size(), mult_boost]
	var banner = _make_banner(msg, Color(0.0, 0.08, 0.12, 0.88), Color("#00FFFF"))
	add_child(banner)
	var tw = create_tween()
	tw.tween_interval(4.0)
	tw.tween_callback(func(): if is_instance_valid(banner): banner.queue_free())

	# 建立鏡像分身標記
	for mirror_info in mirrors:
		_spawn_mirror_marker(mirror_info)

## 鏡像分身被擊破 — 消失閃光 + 浮動文字
func _on_mirror_kill(payload: Dictionary) -> void:
	var mirror_id: String = payload.get("mirror_id", "")

	# 移除鏡像標記
	if _mirror_nodes.has(mirror_id):
		var node = _mirror_nodes[mirror_id]
		if is_instance_valid(node):
			# 消失閃光
			var tw = node.create_tween()
			tw.tween_property(node, "modulate:a", 0.0, 0.25)
			tw.tween_callback(func(): if is_instance_valid(node): node.queue_free())
		_mirror_nodes.erase(mirror_id)

	# 浮動文字「×1.5 鏡像加成！」
	_spawn_float_text("🪞 ×1.5 鏡像加成！", Color("#00FFFF"), 36)

## 鏡像爆炸 — 全螢幕青色三次強閃光 + 大字
func _on_mirror_blast(payload: Dictionary) -> void:
	var blast_count: int = payload.get("blast_count", 0)

	# 清除所有鏡像標記
	_clear_all_mirrors()

	# 全螢幕三次青色強閃光
	_triple_flash(Color("#00FFFF"), 0.75)

	# 大字
	var label = _make_big_label("🪞 鏡像爆炸！", Color("#00FFFF"), 52)
	add_child(label)
	var sub = _make_big_label("爆炸 %d 個分身！" % blast_count, Color("#88FFFF"), 32)
	sub.position.y += 70
	add_child(sub)

	var tw = create_tween()
	tw.tween_interval(2.5)
	tw.tween_callback(func():
		if is_instance_valid(label): label.queue_free()
		if is_instance_valid(sub): sub.queue_free()
	)

## 鏡像結算 — 右側滑入結算彈窗
func _on_mirror_result(payload: Dictionary) -> void:
	var killed_count: int = payload.get("killed_count", 0)
	var blast_count: int = payload.get("blast_count", 0)
	var total_reward: int = payload.get("total_reward", 0)

	var panel = _make_result_panel(killed_count, blast_count, total_reward)
	add_child(panel)

	# 右側滑入
	var vp = get_viewport()
	var vp_size = vp.get_visible_rect().size if vp else Vector2(1280, 720)
	panel.position = Vector2(vp_size.x + 10, vp_size.y * 0.5 - 80)
	var tw = create_tween()
	tw.tween_property(panel, "position:x", vp_size.x - 310.0, 0.35).set_ease(Tween.EASE_OUT)
	tw.tween_interval(3.5)
	tw.tween_property(panel, "position:x", vp_size.x + 10.0, 0.3)
	tw.tween_callback(func(): if is_instance_valid(panel): panel.queue_free())

# ─── 內部工具函數 ───────────────────────────────────────────────────────────

## 建立鏡像分身標記（菱形輪廓 + 閃爍）
func _spawn_mirror_marker(mirror_info: Dictionary) -> void:
	var mirror_id: String = mirror_info.get("mirror_id", "")
	var x: float = mirror_info.get("x", 640.0)
	var y: float = mirror_info.get("y", 360.0)
	var mirror_mult: float = mirror_info.get("mirror_mult", 1.5)

	var container = Control.new()
	container.position = Vector2(x - 32, y - 32)
	container.size = Vector2(64, 64)

	# 菱形輪廓（用 4 個細長 ColorRect 組成）
	var colors = [Color("#00FFFF"), Color("#88FFFF")]
	for i in range(4):
		var line = ColorRect.new()
		line.color = colors[i % 2]
		match i:
			0: line.position = Vector2(28, 0);  line.size = Vector2(8, 32)   # 上
			1: line.position = Vector2(28, 32); line.size = Vector2(8, 32)   # 下
			2: line.position = Vector2(0, 28);  line.size = Vector2(32, 8)   # 左
			3: line.position = Vector2(32, 28); line.size = Vector2(32, 8)   # 右
		container.add_child(line)

	# 倍率標籤
	var mult_label = Label.new()
	mult_label.text = "×%.1f" % mirror_mult
	mult_label.add_theme_color_override("font_color", Color("#00FFFF"))
	mult_label.add_theme_font_size_override("font_size", 14)
	mult_label.position = Vector2(8, 22)
	container.add_child(mult_label)

	# 閃爍動畫
	var tw = container.create_tween().set_loops()
	tw.tween_property(container, "modulate:a", 0.4, 0.5)
	tw.tween_property(container, "modulate:a", 1.0, 0.5)

	add_child(container)
	_mirror_nodes[mirror_id] = container

## 清除所有鏡像標記
func _clear_all_mirrors() -> void:
	for mirror_id in _mirror_nodes.keys():
		var node = _mirror_nodes[mirror_id]
		if is_instance_valid(node):
			node.queue_free()
	_mirror_nodes.clear()

## 建立結算彈窗
func _make_result_panel(killed_count: int, blast_count: int, total_reward: int) -> Control:
	var panel = ColorRect.new()
	panel.color = Color(0.0, 0.08, 0.12, 0.92)
	panel.size = Vector2(300, 160)

	var border = ColorRect.new()
	border.color = Color("#00FFFF")
	border.size = Vector2(300, 4)
	panel.add_child(border)

	var title = Label.new()
	title.text = "🪞 鏡像結算"
	title.add_theme_color_override("font_color", Color("#00FFFF"))
	title.add_theme_font_size_override("font_size", 22)
	title.position = Vector2(12, 14)
	panel.add_child(title)

	var lines = [
		"玩家擊破：%d 個分身" % killed_count,
		"鏡像爆炸：%d 個分身" % blast_count,
		"全服獎勵：+%d 金幣" % total_reward,
	]
	for i in range(lines.size()):
		var lbl = Label.new()
		lbl.text = lines[i]
		lbl.add_theme_color_override("font_color", Color("#88FFFF"))
		lbl.add_theme_font_size_override("font_size", 16)
		lbl.position = Vector2(12, 50 + i * 28)
		panel.add_child(lbl)

	return panel

## 青色雙閃光
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
	label.add_theme_font_size_override("font_size", 18)
	label.position = Vector2(16, 12)
	label.size = Vector2(vp_size.x - 32, 32)
	banner.add_child(label)
	return banner

## 建立大字標籤（置中）
func _make_big_label(text: String, color: Color, font_size: int) -> Label:
	var vp = get_viewport()
	var vp_size = vp.get_visible_rect().size if vp else Vector2(1280, 720)
	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", font_size)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.size = Vector2(vp_size.x, font_size + 16)
	label.position = Vector2(0, vp_size.y * 0.38)
	return label

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
