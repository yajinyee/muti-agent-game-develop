## LuckyPhantomFishPanel.gd — 幸運幽靈魚系統面板（DAY-245）
## 業界原創「幽靈殘影+死亡後復活攻擊」機制
##
## 視覺設計：
##   - 幽靈紫主題（#8E44AD + #6C3483 + #D7BDE2 + #F5EEF8）
##   - phantom_start：紫色三次強閃光 + 頂部橫幅 + 「👻 幽靈護盾！」大字 + 右側豎向計時條 + 倍率說明
##   - phantom_broadcast：頂部小橫幅（全服廣播）
##   - phantom_ghost_created：幽靈殘影標記（半透明紫色圓圈，閃爍動畫）
##   - phantom_ghost_killed：殘影擊破閃光 + ×1.5 浮動文字
##   - phantom_burst：三次強閃光 + 「👻 幽靈爆發！」大字 + 結算彈窗
##   - phantom_end：計時條淡出 + 護盾結束提示
extends CanvasLayer

# 主題顏色
const COLOR_PHANTOM  = Color("#8E44AD")  # 幽靈紫（主題）
const COLOR_DARK     = Color("#6C3483")  # 深紫（強調）
const COLOR_LIGHT    = Color("#D7BDE2")  # 淺紫（背景）
const COLOR_GHOST    = Color("#BB8FCE")  # 殘影紫（半透明）
const COLOR_GOLD     = Color("#F39C12")  # 金色（獎勵）
const COLOR_WHITE    = Color("#FFFFFF")  # 白色

# 計時條節點
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null
var _duration_sec: int = 12

# 幽靈殘影標記（ghostID → Control）
var _ghost_markers: Dictionary = {}

func _ready() -> void:
	layer = 18  # 幸運幽靈魚面板層級（DAY-245）

## 處理幸運幽靈魚訊息
func handle_lucky_phantom_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"phantom_start":
			_on_phantom_start(payload)
		"phantom_broadcast":
			_on_phantom_broadcast(payload)
		"phantom_ghost_created":
			_on_phantom_ghost_created(payload)
		"phantom_ghost_killed":
			_on_phantom_ghost_killed(payload)
		"phantom_burst":
			_on_phantom_burst(payload)
		"phantom_end":
			_on_phantom_end(payload)

## phantom_start — 幽靈護盾啟動（個人訊息）
func _on_phantom_start(payload: Dictionary) -> void:
	_duration_sec = payload.get("duration_sec", 12)
	var ghost_kill_mult: float = payload.get("ghost_kill_mult", 1.5)
	var burst_mult: float = payload.get("burst_mult", 2.0)

	# 紫色三次強閃光
	_flash_screen(COLOR_PHANTOM, 0.5, 3)

	# 頂部橫幅
	_show_banner("👻 幽靈護盾啟動！擊破目標留下殘影！", COLOR_PHANTOM, 3.5)

	# 中央大字
	_show_big_text("👻 幽靈護盾！", COLOR_PHANTOM, 52, 2.5)

	# 倍率說明
	_show_sub_text("殘影擊破 ×%.1f  幽靈爆發 ×%.1f" % [ghost_kill_mult, burst_mult], COLOR_GOLD, 2.0)

	# 右側豎向計時條
	_start_timer_bar(_duration_sec)

## phantom_broadcast — 全服廣播幽靈護盾啟動
func _on_phantom_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	_show_top_banner("👻 %s 觸發幽靈護盾！" % player_name, COLOR_PHANTOM, 2.5)

## phantom_ghost_created — 幽靈殘影生成
func _on_phantom_ghost_created(payload: Dictionary) -> void:
	var ghost_id: String = payload.get("ghost_id", "")
	var x: float = payload.get("x", 0.0)
	var y: float = payload.get("y", 0.0)
	var ghost_duration_sec: int = payload.get("ghost_duration_sec", 5)

	if ghost_id.is_empty():
		return

	# 建立幽靈殘影標記（半透明紫色圓圈）
	var marker = _create_ghost_marker(x, y, ghost_id)
	_ghost_markers[ghost_id] = marker

	# 5 秒後自動移除標記
	var t = get_tree().create_timer(float(ghost_duration_sec))
	t.timeout.connect(func():
		if _ghost_markers.has(ghost_id):
			var m = _ghost_markers[ghost_id]
			if is_instance_valid(m):
				m.queue_free()
			_ghost_markers.erase(ghost_id)
	)

## phantom_ghost_killed — 幽靈殘影被擊破
func _on_phantom_ghost_killed(payload: Dictionary) -> void:
	var ghost_id: String = payload.get("ghost_id", "")
	var reward: int = payload.get("reward", 0)
	var kill_mult: float = payload.get("kill_mult", 1.5)

	# 移除殘影標記
	if _ghost_markers.has(ghost_id):
		var marker = _ghost_markers[ghost_id]
		if is_instance_valid(marker):
			# 擊破閃光
			var tween = create_tween()
			tween.tween_property(marker, "modulate:a", 0.0, 0.3)
			tween.tween_callback(marker.queue_free)
		_ghost_markers.erase(ghost_id)

	# 浮動獎勵文字
	_show_float_text("👻 ×%.1f  +%d" % [kill_mult, reward], COLOR_GHOST, 1.5)

## phantom_burst — 幽靈爆發
func _on_phantom_burst(payload: Dictionary) -> void:
	var ghost_count: int = payload.get("ghost_count", 0)
	var total_reward: int = payload.get("total_reward", 0)
	var burst_mult: float = payload.get("burst_mult", 2.0)

	# 清除所有殘影標記
	for ghost_id in _ghost_markers.keys():
		var marker = _ghost_markers[ghost_id]
		if is_instance_valid(marker):
			marker.queue_free()
	_ghost_markers.clear()

	if ghost_count > 0:
		# 三次強閃光
		_flash_screen(COLOR_DARK, 0.6, 3)

		# 中央大字
		_show_big_text("👻 幽靈爆發！", COLOR_PHANTOM, 56, 2.5)

		# 結算彈窗
		_show_result_popup(ghost_count, total_reward, burst_mult)
	else:
		# 無殘影時的提示
		_show_big_text("👻 護盾結束", COLOR_LIGHT, 40, 1.5)

## phantom_end — 幽靈護盾結束
func _on_phantom_end(_payload: Dictionary) -> void:
	# 停止計時條
	if is_instance_valid(_timer_tween):
		_timer_tween.kill()
	if is_instance_valid(_timer_bar):
		var tween = create_tween()
		tween.tween_property(_timer_bar, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_timer_bar.queue_free)
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		var tween2 = create_tween()
		tween2.tween_property(_timer_bar_bg, "modulate:a", 0.0, 0.5)
		tween2.tween_callback(_timer_bar_bg.queue_free)
		_timer_bar_bg = null

# ─── 幽靈殘影標記 ───────────────────────────────────────────────────────────

func _create_ghost_marker(x: float, y: float, ghost_id: String) -> Control:
	var vp_size = get_viewport().size
	var container = Control.new()
	container.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(container)

	# 半透明紫色圓圈（用 ColorRect 模擬）
	var circle = ColorRect.new()
	circle.color = Color(COLOR_GHOST.r, COLOR_GHOST.g, COLOR_GHOST.b, 0.5)
	circle.size = Vector2(48, 48)
	circle.position = Vector2(x - 24, y - 24)
	container.add_child(circle)

	# 幽靈標籤
	var label = Label.new()
	label.text = "👻"
	label.add_theme_font_size_override("font_size", 24)
	label.position = Vector2(x - 16, y - 20)
	container.add_child(label)

	# 閃爍動畫
	var tween = create_tween()
	tween.set_loops()
	tween.tween_property(circle, "modulate:a", 0.2, 0.5)
	tween.tween_property(circle, "modulate:a", 0.7, 0.5)

	return container

# ─── 結算彈窗 ────────────────────────────────────────────────────────────────

func _show_result_popup(ghost_count: int, total_reward: int, burst_mult: float) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(320, 160)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 80)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.08, 0.02, 0.12, 0.92)
	style.border_color = COLOR_PHANTOM
	style.set_border_width_all(3)
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 8)
	popup.add_child(vbox)

	var title_lbl = Label.new()
	title_lbl.text = "👻 幽靈爆發！"
	title_lbl.add_theme_font_size_override("font_size", 22)
	title_lbl.add_theme_color_override("font_color", COLOR_PHANTOM)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var count_lbl = Label.new()
	count_lbl.text = "殘影數量：%d 個" % ghost_count
	count_lbl.add_theme_font_size_override("font_size", 16)
	count_lbl.add_theme_color_override("font_color", COLOR_LIGHT)
	count_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(count_lbl)

	var mult_lbl = Label.new()
	mult_lbl.text = "爆發倍率：×%.1f" % burst_mult
	mult_lbl.add_theme_font_size_override("font_size", 16)
	mult_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "總獎勵：+%d 籌碼" % total_reward
	reward_lbl.add_theme_font_size_override("font_size", 18)
	reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 340.0, 0.4).set_ease(Tween.EASE_OUT)
	tween.tween_interval(3.0)
	tween.tween_property(popup, "position:x", vp_size.x + 10.0, 0.4).set_ease(Tween.EASE_IN)
	tween.tween_callback(popup.queue_free)

# ─── 通用 UI 工具 ─────────────────────────────────────────────────────────────

func _flash_screen(color: Color, alpha: float, times: int) -> void:
	var vp_size = get_viewport().size
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.0)
	flash.size = vp_size
	add_child(flash)
	var tween = create_tween()
	for i in range(times):
		tween.tween_property(flash, "color:a", alpha, 0.08)
		tween.tween_property(flash, "color:a", 0.0, 0.12)
	tween.tween_callback(flash.queue_free)

func _show_banner(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var banner = PanelContainer.new()
	banner.size = Vector2(vp_size.x, 48)
	banner.position = Vector2(0, 0)
	add_child(banner)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(color.r, color.g, color.b, 0.88)
	banner.add_theme_stylebox_override("panel", style)
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 18)
	lbl.add_theme_color_override("font_color", COLOR_WHITE)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	banner.add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(banner.queue_free)

func _show_top_banner(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var banner = PanelContainer.new()
	banner.size = Vector2(vp_size.x * 0.7, 36)
	banner.position = Vector2(vp_size.x * 0.15, 52)
	add_child(banner)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(color.r, color.g, color.b, 0.82)
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	banner.add_theme_stylebox_override("panel", style)
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 14)
	lbl.add_theme_color_override("font_color", COLOR_WHITE)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	banner.add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(banner.queue_free)

func _show_big_text(text: String, color: Color, font_size: int, duration: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", font_size)
	lbl.add_theme_color_override("font_color", color)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.size = Vector2(vp_size.x, 80)
	lbl.position = Vector2(0, vp_size.y * 0.35)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "scale", Vector2(1.15, 1.15), 0.15).set_ease(Tween.EASE_OUT)
	tween.tween_property(lbl, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(duration - 0.5)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.5)
	tween.tween_callback(lbl.queue_free)

func _show_sub_text(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 18)
	lbl.add_theme_color_override("font_color", color)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.size = Vector2(vp_size.x, 40)
	lbl.position = Vector2(0, vp_size.y * 0.35 + 80)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.4)
	tween.tween_callback(lbl.queue_free)

func _show_float_text(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 20)
	lbl.add_theme_color_override("font_color", color)
	lbl.position = Vector2(
		randf_range(vp_size.x * 0.3, vp_size.x * 0.7),
		randf_range(vp_size.y * 0.3, vp_size.y * 0.6)
	)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 60, duration)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, duration)
	tween.tween_callback(lbl.queue_free)

func _start_timer_bar(duration_sec: int) -> void:
	var vp_size = get_viewport().size
	var bar_h: float = vp_size.y * 0.5
	var bar_w: float = 12.0
	var bar_x: float = vp_size.x - 28.0
	var bar_y: float = vp_size.y * 0.25

	# 背景條
	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.1, 0.0, 0.15, 0.7)
	_timer_bar_bg.size = Vector2(bar_w, bar_h)
	_timer_bar_bg.position = Vector2(bar_x, bar_y)
	add_child(_timer_bar_bg)

	# 計時條（從上往下縮短）
	_timer_bar = ColorRect.new()
	_timer_bar.color = COLOR_PHANTOM
	_timer_bar.size = Vector2(bar_w, bar_h)
	_timer_bar.position = Vector2(bar_x, bar_y)
	add_child(_timer_bar)

	# 動畫：從滿到空
	_timer_tween = create_tween()
	_timer_tween.tween_property(_timer_bar, "size:y", 0.0, float(duration_sec)).set_ease(Tween.EASE_IN_OUT)
