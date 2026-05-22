## AbyssVortexPanel.gd — 深淵漩渦魚面板（DAY-202）
## 業界依據：Ocean King 2「Vortex Fish — sucks all fish of the same species into a whirlpool.
## Catching a Vortex Fish will suck all fish of the same species in the area into a whirlpool.」
## + SteamDB OceanFest 2026「Abyssal Vortex (Depth 3, persistent whirlpool)」
##
## 視覺設計：
##   - 深淵藍紫主題（#1A0033 + #6600CC + #00CCFF + #FFFFFF）
##   - vortex_start：深藍紫色雙閃光 + 頂部橫幅「🌀 深淵漩渦！」+ 漩渦旋轉動畫（中心）
##   - vortex_pulse：漩渦脈衝閃光 + 擊破浮動文字 + 脈衝計數器
##   - vortex_blast：全螢幕深藍白強閃光 + 爆炸圓圈擴散（深淵爆炸）
##   - vortex_result：右側滑入結算彈窗（脈衝擊破/爆炸擊破/總獎勵）
extends CanvasLayer

var _panel: Control
var _banner: Label
var _pulse_counter: Label
var _vortex_visual: Control  # 漩渦視覺效果容器
var _result_popup: Control

var _vortex_x: float = 640.0
var _vortex_y: float = 360.0
var _pulse_kills: int = 0
var _vortex_tween: Tween = null

func _ready() -> void:
	layer = 43
	_build_ui()

func _build_ui() -> void:
	_panel = Control.new()
	_panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	_panel.visible = false
	_panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_panel)

	# 頂部橫幅
	_banner = Label.new()
	_banner.text = "🌀 深淵漩渦！"
	_banner.add_theme_font_size_override("font_size", 24)
	_banner.add_theme_color_override("font_color", Color("#00CCFF"))
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.position = Vector2(0, 12)
	_banner.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.visible = false
	add_child(_banner)

	# 脈衝計數器
	_pulse_counter = Label.new()
	_pulse_counter.text = "吸入: 0 個"
	_pulse_counter.add_theme_font_size_override("font_size", 18)
	_pulse_counter.add_theme_color_override("font_color", Color("#6600CC"))
	_pulse_counter.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_pulse_counter.position = Vector2(0, 42)
	_pulse_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_pulse_counter.visible = false
	add_child(_pulse_counter)

	# 漩渦視覺效果容器
	_vortex_visual = Control.new()
	_vortex_visual.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_vortex_visual.visible = false
	add_child(_vortex_visual)

	# 結算彈窗
	_result_popup = Control.new()
	_result_popup.visible = false
	_result_popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_popup.size = Vector2(300, 220)
	_result_popup.position = Vector2(-320, -110)
	add_child(_result_popup)

	var popup_bg = ColorRect.new()
	popup_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_bg.color = Color(0.05, 0.0, 0.15, 0.93)
	_result_popup.add_child(popup_bg)

	var popup_label = Label.new()
	popup_label.name = "ResultLabel"
	popup_label.add_theme_font_size_override("font_size", 15)
	popup_label.add_theme_color_override("font_color", Color("#00CCFF"))
	popup_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_popup.add_child(popup_label)

## 處理深淵漩渦訊息
func handle_abyss_vortex(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"vortex_start":
			_on_vortex_start(payload)
		"vortex_pulse":
			_on_vortex_pulse(payload)
		"vortex_blast":
			_on_vortex_blast(payload)
		"vortex_result":
			_on_vortex_result(payload)

## 漩渦開始
func _on_vortex_start(payload: Dictionary) -> void:
	var killer_name = payload.get("killer_name", "玩家")
	_vortex_x = payload.get("vortex_x", 640.0)
	_vortex_y = payload.get("vortex_y", 360.0)
	_pulse_kills = 0

	_panel.visible = true

	# 深藍紫色雙閃光
	_flash_screen(Color("#1A0033"), 0.7)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color("#6600CC"), 0.5)

	# 顯示橫幅
	_banner.text = "🌀 " + killer_name + " 觸發深淵漩渦！"
	_banner.visible = true
	_pulse_counter.text = "吸入: 0 個"
	_pulse_counter.visible = true

	# 啟動漩渦旋轉視覺
	_vortex_visual.visible = true
	_start_vortex_animation()

## 漩渦脈衝
func _on_vortex_pulse(payload: Dictionary) -> void:
	var pulse_kills = payload.get("pulse_kills", 0)
	var pulse_reward = payload.get("pulse_reward", 0)
	var total_kills = payload.get("total_kills", 0)

	_pulse_kills = total_kills
	_pulse_counter.text = "吸入: %d 個" % total_kills

	# 脈衝閃光（藍紫色）
	_flash_screen(Color("#6600CC"), 0.2)

	# 擊破浮動文字
	if pulse_kills > 0:
		_spawn_float_text("🌀 ×%d +%d" % [pulse_kills, pulse_reward],
			Color("#00CCFF"),
			Vector2(_vortex_x + randf_range(-80, 80), _vortex_y - 60))

## 深淵爆炸
func _on_vortex_blast(payload: Dictionary) -> void:
	var blast_kills = payload.get("blast_kills", 0)
	var blast_reward = payload.get("blast_reward", 0)

	# 停止漩渦動畫
	_stop_vortex_animation()

	# 全螢幕深藍白強閃光（爆炸感）
	_flash_screen(Color("#FFFFFF"), 0.9)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(Color("#00CCFF"), 0.7)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(Color("#6600CC"), 0.5)

	# 爆炸圓圈擴散
	_spawn_blast_explosion(Vector2(_vortex_x, _vortex_y), blast_kills)

	# 爆炸浮動文字
	if blast_kills > 0:
		_spawn_float_text("💥 深淵爆炸！×%d +%d" % [blast_kills, blast_reward],
			Color("#FFD700"),
			Vector2(_vortex_x, _vortex_y - 80))

## 最終結算
func _on_vortex_result(payload: Dictionary) -> void:
	var total_kills = payload.get("total_kills", 0)
	var total_reward = payload.get("total_reward", 0)

	# 隱藏橫幅和計數器
	_banner.visible = false
	_pulse_counter.visible = false
	_vortex_visual.visible = false

	# 依擊破數決定閃光
	if total_kills >= 10:
		_flash_screen(Color("#FFD700"), 0.6)
	elif total_kills >= 5:
		_flash_screen(Color("#00CCFF"), 0.4)

	# 顯示結算彈窗
	var result_label = _result_popup.get_node("ResultLabel")
	result_label.text = (
		"🌀 深淵漩渦結算\n\n"
		+ "漩渦吸入: %d 個\n" % _pulse_kills
		+ "深淵爆炸: %d 個\n" % (total_kills - _pulse_kills)
		+ "總擊破: %d 個\n" % total_kills
		+ "總獎勵: %d 金幣" % total_reward
	)
	_result_popup.visible = true
	_result_popup.position = Vector2(get_viewport().get_visible_rect().size.x + 10, -110)

	# 從右側滑入
	var tween = create_tween()
	tween.tween_property(_result_popup, "position:x",
		get_viewport().get_visible_rect().size.x - 320, 0.4)

	# 4 秒後淡出
	await get_tree().create_timer(4.5).timeout
	var fade_tween = create_tween()
	fade_tween.tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	await fade_tween.finished
	_result_popup.visible = false
	_result_popup.modulate.a = 1.0
	_panel.visible = false

## 啟動漩渦旋轉動畫（同心圓旋轉效果）
func _start_vortex_animation() -> void:
	# 清除舊的視覺元素
	for child in _vortex_visual.get_children():
		child.queue_free()

	# 建立 3 個同心圓（旋轉漩渦感）
	for i in range(3):
		var ring = ColorRect.new()
		var ring_size = 60.0 + i * 50.0
		ring.size = Vector2(ring_size, ring_size)
		ring.position = Vector2(_vortex_x - ring_size / 2, _vortex_y - ring_size / 2)
		var alpha = 0.6 - i * 0.15
		ring.color = Color(0.4, 0.0, 0.8, alpha)
		_vortex_visual.add_child(ring)

		# 旋轉動畫
		var ring_tween = ring.create_tween().set_loops()
		var rot_speed = 1.5 - i * 0.3  # 內圈轉快，外圈轉慢
		ring_tween.tween_property(ring, "rotation", TAU, rot_speed)

	# 中心點（深淵核心）
	var core = ColorRect.new()
	core.size = Vector2(20, 20)
	core.position = Vector2(_vortex_x - 10, _vortex_y - 10)
	core.color = Color("#00CCFF")
	_vortex_visual.add_child(core)

	# 核心脈衝動畫
	var core_tween = core.create_tween().set_loops()
	core_tween.tween_property(core, "modulate:a", 0.2, 0.4)
	core_tween.tween_property(core, "modulate:a", 1.0, 0.4)

## 停止漩渦動畫
func _stop_vortex_animation() -> void:
	for child in _vortex_visual.get_children():
		if is_instance_valid(child):
			child.queue_free()

## 深淵爆炸圓圈動畫
func _spawn_blast_explosion(pos: Vector2, kills: int) -> void:
	# 主爆炸圓圈（深淵藍紫）
	var circle = ColorRect.new()
	var size = 20.0
	circle.size = Vector2(size, size)
	circle.position = pos - Vector2(size / 2, size / 2)
	circle.color = Color("#6600CC")
	add_child(circle)

	var tween = create_tween()
	tween.tween_property(circle, "size", Vector2(600, 600), 0.6)
	tween.parallel().tween_property(circle, "position", pos - Vector2(300, 300), 0.6)
	tween.parallel().tween_property(circle, "modulate:a", 0.0, 0.6)
	tween.tween_callback(circle.queue_free)

	# 外圈（青色）
	var outer = ColorRect.new()
	outer.size = Vector2(40, 40)
	outer.position = pos - Vector2(20, 20)
	outer.color = Color("#00CCFF")
	add_child(outer)

	var outer_tween = create_tween()
	outer_tween.tween_property(outer, "size", Vector2(700, 700), 0.7)
	outer_tween.parallel().tween_property(outer, "position", pos - Vector2(350, 350), 0.7)
	outer_tween.parallel().tween_property(outer, "modulate:a", 0.0, 0.7)
	outer_tween.tween_callback(outer.queue_free)

	# 4 方向深淵射線
	for angle in [0, 90, 180, 270]:
		var line = ColorRect.new()
		line.size = Vector2(6, 50)
		line.color = Color("#00CCFF")
		line.position = pos
		line.rotation_degrees = angle
		add_child(line)
		var lt = create_tween()
		lt.tween_property(line, "position",
			pos + Vector2(cos(deg_to_rad(angle)) * 120, sin(deg_to_rad(angle)) * 120), 0.5)
		lt.parallel().tween_property(line, "modulate:a", 0.0, 0.5)
		lt.tween_callback(line.queue_free)

## 全螢幕閃光
func _flash_screen(color: Color, alpha: float) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, alpha)
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.25)
	tween.tween_callback(flash.queue_free)

## 浮動文字
func _spawn_float_text(text: String, color: Color, pos: Vector2 = Vector2(-1, -1)) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 20)
	label.add_theme_color_override("font_color", color)
	if pos.x < 0:
		label.position = Vector2(randf_range(300, 900), randf_range(200, 500))
	else:
		label.position = pos
	add_child(label)
	var tween = create_tween()
	tween.tween_property(label, "position:y", label.position.y - 70, 1.2)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 1.2)
	tween.tween_callback(label.queue_free)
