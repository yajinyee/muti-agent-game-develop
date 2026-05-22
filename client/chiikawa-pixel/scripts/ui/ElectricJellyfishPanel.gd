## ElectricJellyfishPanel.gd — 電流水母電流網路面板（DAY-193）
## 業界依據：King of Ocean 2026「electric jellyfish chains current between adjacent targets,
## paying multipliers from every link in the chain」
## 視覺主題：青色電流 + 電流連接線動畫 + 網路拓撲視覺

extends Control

const ELECTRIC_COLOR := Color(0.0, 1.0, 1.0)   # 青色（電流感）
const KILL_COLOR     := Color(1.0, 1.0, 0.0)   # 黃色（擊破）
const MISS_COLOR     := Color(0.3, 0.8, 1.0)   # 淡藍（未擊破）

var _banner: Control = null
var _link_counter: Label = null
var _total_links: int = 0
var _total_kills: int = 0

func _ready() -> void:
	if GameManager.has_signal("electric_jellyfish"):
		GameManager.electric_jellyfish.connect(_on_electric_jellyfish)

func _on_electric_jellyfish(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	if phase == "network_start":
		_show_network_start(data)
	elif phase == "network_result":
		_show_network_result(data)
	elif phase.begins_with("link_"):
		_show_link(data)

# ── network_start ──────────────────────────────────────────────────────────────
func _show_network_start(data: Dictionary) -> void:
	var killer_name: String = data.get("killer_name", "")
	var link_count: int = data.get("link_count", 0)
	_total_links = link_count
	_total_kills = 0

	# 青色閃光（兩次）
	_flash_screen(ELECTRIC_COLOR, 0.5)
	var t1 = get_tree().create_timer(0.15)
	t1.timeout.connect(func(): _flash_screen(ELECTRIC_COLOR, 0.35))

	# 頂部橫幅
	_show_banner("⚡🪼 %s 的電流水母！建立 %d 條電流連接！" % [killer_name, link_count], ELECTRIC_COLOR)

	# 連接計數器
	_show_link_counter(link_count)

# ── link_N 電流連接 ─────────────────────────────────────────────────────────────
func _show_link(data: Dictionary) -> void:
	var xa: float = data.get("x_a", 0.0)
	var ya: float = data.get("y_a", 0.0)
	var xb: float = data.get("x_b", 0.0)
	var yb: float = data.get("y_b", 0.0)
	var is_kill: bool = data.get("is_kill", false)
	var reward: int = data.get("reward", 0)
	var link_index: int = data.get("link_index", 0)

	var line_color := KILL_COLOR if is_kill else MISS_COLOR

	# 繪製電流連接線
	_draw_electric_line(xa, ya, xb, yb, line_color)

	# 擊破時顯示獎勵浮動文字
	if is_kill and reward > 0:
		var mid_x := (xa + xb) / 2.0
		var mid_y := (ya + yb) / 2.0
		_show_floating_reward("+%d" % reward, mid_x, mid_y, KILL_COLOR)
		_total_kills += 1

	# 更新計數器
	if _link_counter and is_instance_valid(_link_counter):
		_link_counter.text = "⚡ %d/%d 連接 | 擊破 %d" % [link_index, _total_links, _total_kills]

# ── network_result 電流網路結果 ─────────────────────────────────────────────────
func _show_network_result(data: Dictionary) -> void:
	var total_kills: int = data.get("total_kills", 0)
	var total_reward: int = data.get("total_reward", 0)
	var link_count: int = data.get("link_count", 0)

	# 淡出橫幅
	_hide_banner()

	# 右側滑入結算彈窗
	_show_result_popup(link_count, total_kills, total_reward)

	# 淡出計數器
	if _link_counter and is_instance_valid(_link_counter):
		var tween = create_tween()
		tween.tween_property(_link_counter, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_link_counter.queue_free)
		_link_counter = null

	# 大量連接時額外閃光
	if link_count >= 12:
		_flash_screen(ELECTRIC_COLOR, 0.6)
	elif total_kills >= 5:
		_flash_screen(KILL_COLOR, 0.45)

# ── 輔助函數 ────────────────────────────────────────────────────────────────────

func _flash_screen(color: Color, alpha: float) -> void:
	var overlay = ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, alpha)
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tween = create_tween()
	tween.tween_property(overlay, "modulate:a", 0.0, 0.3)
	tween.tween_callback(overlay.queue_free)

func _show_banner(text: String, color: Color) -> void:
	if _banner != null and is_instance_valid(_banner):
		_banner.queue_free()

	_banner = Control.new()
	_banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 50)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.color = Color(0.0, 0.1, 0.15, 0.9)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_banner.add_child(bg)

	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", 19)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.add_child(label)

	_banner.position.y = -50
	var tween = create_tween()
	tween.tween_property(_banner, "position:y", 0.0, 0.25).set_trans(Tween.TRANS_BACK)

func _hide_banner() -> void:
	if _banner != null and is_instance_valid(_banner):
		var tween = create_tween()
		tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
		tween.tween_callback(_banner.queue_free)
		_banner = null

func _show_link_counter(link_count: int) -> void:
	if _link_counter != null and is_instance_valid(_link_counter):
		_link_counter.queue_free()

	_link_counter = Label.new()
	_link_counter.text = "⚡ 0/%d 連接 | 擊破 0" % link_count
	_link_counter.add_theme_color_override("font_color", ELECTRIC_COLOR)
	_link_counter.add_theme_font_size_override("font_size", 16)
	_link_counter.set_anchors_preset(Control.PRESET_BOTTOM_RIGHT)
	_link_counter.offset_left = -260
	_link_counter.offset_top = -42
	_link_counter.offset_right = -10
	_link_counter.offset_bottom = -10
	add_child(_link_counter)

func _draw_electric_line(x1: float, y1: float, x2: float, y2: float, color: Color) -> void:
	# 用多個小 ColorRect 模擬電流線（鋸齒狀）
	var steps := 8
	var prev_x := x1
	var prev_y := y1
	var dx := (x2 - x1) / steps
	var dy := (y2 - y1) / steps

	for i in range(1, steps + 1):
		var nx := x1 + dx * i + randf_range(-6, 6)  # 鋸齒偏移
		var ny := y1 + dy * i + randf_range(-6, 6)

		# 繪製線段（用細長 ColorRect 近似）
		var seg_len := sqrt((nx - prev_x) * (nx - prev_x) + (ny - prev_y) * (ny - prev_y))
		if seg_len < 1.0:
			prev_x = nx
			prev_y = ny
			continue

		var seg = ColorRect.new()
		seg.color = Color(color.r, color.g, color.b, 0.8)
		seg.size = Vector2(seg_len, 3)
		seg.position = Vector2(prev_x, prev_y - 1.5)

		# 旋轉線段
		var angle := atan2(ny - prev_y, nx - prev_x)
		seg.rotation = angle
		seg.pivot_offset = Vector2(0, 1.5)

		add_child(seg)

		# 淡出動畫
		var tween = create_tween()
		tween.tween_property(seg, "modulate:a", 0.0, 0.6)
		tween.tween_callback(seg.queue_free)

		prev_x = nx
		prev_y = ny

	# 在兩端顯示電流節點（小圓點）
	_spawn_node_dot(x1, y1, color)
	_spawn_node_dot(x2, y2, color)

func _spawn_node_dot(x: float, y: float, color: Color) -> void:
	var dot = ColorRect.new()
	dot.color = Color(color.r, color.g, color.b, 0.9)
	dot.size = Vector2(8, 8)
	dot.position = Vector2(x - 4, y - 4)
	add_child(dot)
	var tween = create_tween()
	tween.tween_property(dot, "modulate:a", 0.0, 0.5)
	tween.tween_callback(dot.queue_free)

func _show_floating_reward(text: String, x: float, y: float, color: Color) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", 15)
	label.position = Vector2(x - 20, y - 10)
	add_child(label)

	var tween = create_tween()
	tween.set_parallel(true)
	tween.tween_property(label, "position:y", y - 45, 0.7)
	tween.tween_property(label, "modulate:a", 0.0, 0.7)
	tween.chain().tween_callback(label.queue_free)

func _show_result_popup(link_count: int, total_kills: int, total_reward: int) -> void:
	var popup = Control.new()
	popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	popup.custom_minimum_size = Vector2(210, 120)
	popup.offset_left = -220
	popup.offset_top = -60
	popup.offset_right = -10
	popup.offset_bottom = 60
	add_child(popup)

	var bg = ColorRect.new()
	bg.color = Color(0.0, 0.08, 0.12, 0.92)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	popup.add_child(bg)

	var vbox = VBoxContainer.new()
	vbox.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	vbox.add_theme_constant_override("separation", 6)
	popup.add_child(vbox)

	_add_label(vbox, "⚡🪼 電流網路結束！", ELECTRIC_COLOR, 16)
	_add_label(vbox, "電流連接：%d 條" % link_count, Color.WHITE, 14)
	_add_label(vbox, "擊破目標：%d 個" % total_kills, KILL_COLOR, 14)
	_add_label(vbox, "總獎勵：%d 金幣" % total_reward, ELECTRIC_COLOR, 14)

	popup.position.x += 230
	var tween = create_tween()
	tween.tween_property(popup, "position:x", popup.position.x - 230, 0.3).set_trans(Tween.TRANS_BACK)

	var timer = get_tree().create_timer(4.5)
	timer.timeout.connect(func():
		if is_instance_valid(popup):
			var t2 = create_tween()
			t2.tween_property(popup, "modulate:a", 0.0, 0.5)
			t2.tween_callback(popup.queue_free)
	)

func _add_label(parent: Control, text: String, color: Color, size: int) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", size)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	parent.add_child(label)
