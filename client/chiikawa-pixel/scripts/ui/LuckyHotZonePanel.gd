## LuckyHotZonePanel.gd — 幸運熱區魚空間策略面板（DAY-210）
## 業界依據：Ocean King 4 Brand New World 2025
## 「Golden Zone — a glowing area appears on screen, all fish inside
##  receive a 2x multiplier bonus. Zone lasts 8 seconds then explodes.」
##
## 視覺設計：
##   - 橙火熱區主題（#FF6600 + #FF4500 + #FFD700 + #FF8C00）
##   - zone_start：橙色強閃光 + 頂部橫幅 + 在場上畫出熱區圓圈（持續 8 秒）
##   - zone_pulse：熱區圓圈脈衝閃爍 + 「×2 熱區！」浮動文字 + 剩餘時間更新
##   - zone_blast：全螢幕橙紅三次強閃光 + 「🔥 熱區爆炸！」大字 + 結算彈窗
extends CanvasLayer

var _zone_circle: Control = null   # 熱區圓圈視覺
var _zone_timer_bar: Control = null # 熱區計時條
var _zone_x: float = 640.0
var _zone_y: float = 360.0
var _zone_radius: float = 280.0

func _ready() -> void:
	layer = 35  # 幸運熱區面板層級

## 處理幸運熱區魚訊息
func handle_lucky_hot_zone(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"zone_start":
			_on_zone_start(payload)
		"zone_pulse":
			_on_zone_pulse(payload)
		"zone_blast":
			_on_zone_blast(payload)

## 熱區建立 — 橙色強閃光 + 頂部橫幅 + 熱區圓圈
func _on_zone_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "玩家")
	_zone_x = payload.get("zone_x", 640.0)
	_zone_y = payload.get("zone_y", 360.0)
	_zone_radius = payload.get("radius", 280.0)
	var duration: int = payload.get("duration_sec", 8)
	var multiplier: float = payload.get("multiplier", 2.0)

	# 橙色強閃光
	_flash_screen(Color("#FF6600"), 0.7)

	# 頂部橫幅
	var banner = _make_banner(
		"🔥 %s 觸發幸運熱區！熱區內目標 ×%.0f 倍率！快往熱區打！" % [player_name, multiplier],
		Color(0.2, 0.05, 0.0, 0.88),
		Color("#FF6600")
	)
	add_child(banner)
	var tween_b = create_tween()
	tween_b.tween_property(banner, "position:y", 0.0, 0.25)
	tween_b.tween_interval(3.5)
	tween_b.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween_b.tween_callback(banner.queue_free)

	# 建立熱區圓圈視覺
	_create_zone_circle(duration)

	# 底部計時條
	_show_zone_timer_bar(duration)

## 熱區脈衝 — 圓圈閃爍 + 浮動文字
func _on_zone_pulse(payload: Dictionary) -> void:
	var affected: int = payload.get("affected_count", 0)
	var remaining: int = payload.get("remaining_sec", 0)

	# 圓圈脈衝閃爍
	if is_instance_valid(_zone_circle):
		var tween = create_tween()
		tween.tween_property(_zone_circle, "modulate:a", 1.0, 0.08)
		tween.tween_property(_zone_circle, "modulate:a", 0.6, 0.12)

	# 浮動文字（有目標被削血時）
	if affected > 0:
		var float_label = Label.new()
		float_label.text = "🔥 ×2 熱區！-%d 目標" % affected
		float_label.add_theme_font_size_override("font_size", 16)
		float_label.add_theme_color_override("font_color", Color("#FF8C00"))
		float_label.position = Vector2(_zone_x - 80, _zone_y - _zone_radius - 30)
		add_child(float_label)
		var tween2 = create_tween()
		tween2.tween_property(float_label, "position:y", float_label.position.y - 25, 0.6)
		tween2.parallel().tween_property(float_label, "modulate:a", 0.0, 0.6)
		tween2.tween_callback(float_label.queue_free)

## 熱區爆炸 — 全螢幕橙紅三次強閃光 + 大字 + 結算彈窗
func _on_zone_blast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "玩家")
	var killed: int = payload.get("killed_count", 0)
	var reward_per: int = payload.get("reward_per_player", 0)
	var total: int = payload.get("total_reward", 0)

	# 移除熱區圓圈和計時條
	if is_instance_valid(_zone_circle):
		_zone_circle.queue_free()
		_zone_circle = null
	if is_instance_valid(_zone_timer_bar):
		_zone_timer_bar.queue_free()
		_zone_timer_bar = null

	# 全螢幕三次橙紅強閃光
	_flash_screen(Color("#FF4500"), 1.0)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color("#FF6600"), 0.85)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color("#FFD700"), 0.65)

	# 「🔥 熱區爆炸！」大字
	var big_label = Label.new()
	big_label.text = "🔥 熱區爆炸！"
	big_label.add_theme_font_size_override("font_size", 52)
	big_label.add_theme_color_override("font_color", Color("#FF4500"))
	big_label.set_anchors_preset(Control.PRESET_CENTER)
	big_label.position = Vector2(-200, -70)
	add_child(big_label)
	var tween = create_tween()
	tween.tween_property(big_label, "scale", Vector2(1.3, 1.3), 0.25)
	tween.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.2)
	tween.tween_interval(1.5)
	tween.tween_property(big_label, "modulate:a", 0.0, 0.5)
	tween.tween_callback(big_label.queue_free)

	# 結算彈窗
	if killed > 0:
		_show_blast_popup(player_name, killed, reward_per, total)

## 建立熱區圓圈視覺（持續 duration 秒）
func _create_zone_circle(duration: int) -> void:
	if is_instance_valid(_zone_circle):
		_zone_circle.queue_free()

	# 用 Node2D 繪製圓圈（在 CanvasLayer 上）
	var circle_node = _ZoneCircle.new()
	circle_node.zone_x = _zone_x
	circle_node.zone_y = _zone_y
	circle_node.zone_radius = _zone_radius
	add_child(circle_node)
	_zone_circle = circle_node

	# 圓圈淡入
	circle_node.modulate.a = 0.0
	var tween = create_tween()
	tween.tween_property(circle_node, "modulate:a", 0.75, 0.4)

## 底部熱區計時條
func _show_zone_timer_bar(duration: int) -> void:
	if is_instance_valid(_zone_timer_bar):
		_zone_timer_bar.queue_free()

	var bar_container = Control.new()
	bar_container.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	bar_container.position = Vector2(0, -30)
	bar_container.size = Vector2(1280, 30)
	add_child(bar_container)
	_zone_timer_bar = bar_container

	var bg = ColorRect.new()
	bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.15, 0.05, 0.0, 0.85)
	bar_container.add_child(bg)

	var bar = ColorRect.new()
	bar.name = "ZoneBar"
	bar.set_anchors_preset(Control.PRESET_FULL_RECT)
	bar.color = Color("#FF6600")
	bar_container.add_child(bar)

	var label = Label.new()
	label.text = "🔥 幸運熱區！熱區內目標 ×2 倍率！"
	label.add_theme_font_size_override("font_size", 14)
	label.add_theme_color_override("font_color", Color("#FFFFFF"))
	label.set_anchors_preset(Control.PRESET_FULL_RECT)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	bar_container.add_child(label)

	var tween = create_tween()
	tween.tween_property(bar, "size:x", 0.0, float(duration))
	tween.parallel().tween_method(
		func(t: float):
			if is_instance_valid(bar):
				if t > 0.5:
					bar.color = Color("#FF6600")
				elif t > 0.25:
					bar.color = Color("#FF4500")
				else:
					bar.color = Color("#CC2200"),
		1.0, 0.0, float(duration)
	)

## 爆炸結算彈窗
func _show_blast_popup(player_name: String, killed: int, reward_per: int, total: int) -> void:
	var popup = PanelContainer.new()
	popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	popup.position = Vector2(1400, -100)
	popup.size = Vector2(260, 160)

	var border_color = Color("#FF6600")
	if killed >= 10:
		border_color = Color("#FF4500")

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.12, 0.04, 0.0, 0.92)
	style.border_color = border_color
	style.set_border_width_all(3)
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup.add_child(vbox)

	var title_lbl = Label.new()
	title_lbl.text = "🔥 熱區爆炸結算"
	title_lbl.add_theme_font_size_override("font_size", 16)
	title_lbl.add_theme_color_override("font_color", Color("#FF6600"))
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var killed_lbl = Label.new()
	killed_lbl.text = "擊破 %d 個目標" % killed
	killed_lbl.add_theme_font_size_override("font_size", 14)
	killed_lbl.add_theme_color_override("font_color", Color("#FFD700"))
	killed_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(killed_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "每人 +%d 金幣" % reward_per
	reward_lbl.add_theme_font_size_override("font_size", 18)
	reward_lbl.add_theme_color_override("font_color", Color("#FFFFFF"))
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	var total_lbl = Label.new()
	total_lbl.text = "總獎勵 %d" % total
	total_lbl.add_theme_font_size_override("font_size", 12)
	total_lbl.add_theme_color_override("font_color", Color("#FF8C00"))
	total_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(total_lbl)

	add_child(popup)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", 1010.0, 0.3)
	tween.tween_interval(3.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)

## 頂部橫幅
func _make_banner(text: String, bg_color: Color, text_color: Color) -> PanelContainer:
	var panel = PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_WIDE)
	panel.position = Vector2(0, -44)
	panel.size = Vector2(1280, 44)

	var style = StyleBoxFlat.new()
	style.bg_color = bg_color
	panel.add_theme_stylebox_override("panel", style)

	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 17)
	label.add_theme_color_override("font_color", text_color)
	label.set_anchors_preset(Control.PRESET_FULL_RECT)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	panel.add_child(label)
	return panel

## 全螢幕閃光
func _flash_screen(color: Color, alpha: float) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, alpha)
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.22)
	tween.tween_callback(flash.queue_free)

## 熱區圓圈繪製節點（內部類別）
class _ZoneCircle extends Node2D:
	var zone_x: float = 640.0
	var zone_y: float = 360.0
	var zone_radius: float = 280.0
	var _time: float = 0.0

	func _process(delta: float) -> void:
		_time += delta
		queue_redraw()

	func _draw() -> void:
		var pos = Vector2(zone_x, zone_y)
		var pulse = 0.6 + 0.4 * sin(_time * 3.0)  # 脈衝效果

		# 外圈（橙色，虛線感）
		draw_arc(pos, zone_radius, 0, TAU, 64,
			Color(1.0, 0.4, 0.0, 0.8 * pulse), 3.0)

		# 中圈（橙黃色）
		draw_arc(pos, zone_radius - 8, 0, TAU, 64,
			Color(1.0, 0.6, 0.0, 0.5 * pulse), 2.0)

		# 內圈（金色）
		draw_arc(pos, zone_radius - 16, 0, TAU, 64,
			Color(1.0, 0.85, 0.0, 0.35 * pulse), 1.5)

		# 填充（半透明橙色）
		draw_circle(pos, zone_radius,
			Color(1.0, 0.4, 0.0, 0.08 * pulse))

		# 「×2」中心文字（用 draw_string 模擬）
		# 注意：Godot 4 的 draw_string 需要 Font 資源，這裡用圓點代替
		draw_circle(pos, 6.0, Color(1.0, 0.6, 0.0, 0.9 * pulse))
