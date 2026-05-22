## DragonKingPanel.gd — 深海龍王全服合力蓄力面板（DAY-208）
## 業界依據：Royal Fishing JILI「Dragon Wrath — accumulate wrath value through shooting,
##  then unleash devastating meteor strikes across the entire screen.」
##
## 視覺設計：
##   - 深紅龍焰主題（#FF4500 + #FF0000 + #FF6600 + #FFD700 + #1A0000）
##   - charge_start：深紅三次閃光 + 頂部橫幅 + 底部蓄力進度條（0/20）
##   - charge_progress：進度條更新 + 每 5 點額外閃光 + 「還差 N 槍！」提示
##   - meteor_rain_start：全螢幕深紅爆炸 + 「🐉 龍怒隕石雨！」大字
##   - small_meteor_start：橙色閃光 + 「🐉 小型龍怒！」大字
##   - meteor_hit：隕石落點爆炸圓圈 + 浮動獎勵文字
##   - meteor_rain_result / small_meteor_result：右側滑入結算彈窗
extends CanvasLayer

var _charge_bar: Control = null
var _charge_label: Label = null
var _charge_current: int = 0
var _charge_target: int = 20

func _ready() -> void:
	layer = 37  # 深海龍王面板層級

## 處理深海龍王訊息
func handle_dragon_king(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"charge_start":
			_on_charge_start(payload)
		"charge_progress":
			_on_charge_progress(payload)
		"meteor_rain_start":
			_on_meteor_start(payload, true)
		"small_meteor_start":
			_on_meteor_start(payload, false)
		"meteor_hit":
			_on_meteor_hit(payload)
		"meteor_rain_result":
			_on_meteor_result(payload, true)
		"small_meteor_result":
			_on_meteor_result(payload, false)

## 蓄力開始 — 深紅三次閃光 + 頂部橫幅 + 底部蓄力進度條
func _on_charge_start(payload: Dictionary) -> void:
	_charge_target = payload.get("charge_target", 20)
	var charge_sec = payload.get("charge_sec", 12)
	_charge_current = 0

	# 深紅三次閃光
	_flash_screen(Color("#FF4500"), 0.8)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(Color("#FF0000"), 0.6)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(Color("#FF6600"), 0.4)

	# 頂部橫幅
	var banner = _make_banner(
		"🐉 深海龍王覺醒！全服合力射擊 %d 次觸發龍怒隕石雨！（%d 秒）" % [_charge_target, charge_sec],
		Color("#8B0000"),
		Color("#FFD700")
	)
	add_child(banner)
	var tween = create_tween()
	tween.tween_property(banner, "position:y", 10.0, 0.3)
	tween.tween_interval(float(charge_sec) - 0.5)
	tween.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(banner.queue_free)

	# 底部蓄力進度條
	_show_charge_bar(charge_sec)

## 蓄力進度更新 — 進度條更新 + 每 5 點額外閃光
func _on_charge_progress(payload: Dictionary) -> void:
	_charge_current = payload.get("current", 0)
	_charge_target = payload.get("charge_target", 20)

	# 更新進度條
	if is_instance_valid(_charge_bar):
		var bar = _charge_bar.get_node_or_null("ChargeBar")
		if is_instance_valid(bar):
			var pct = float(_charge_current) / float(_charge_target)
			bar.size.x = 1280.0 * pct
			# 顏色漸變：深紅→橙→金
			if pct < 0.4:
				bar.color = Color("#FF4500")
			elif pct < 0.7:
				bar.color = Color("#FF6600")
			else:
				bar.color = Color("#FFD700")

	# 更新文字
	if is_instance_valid(_charge_label):
		var remaining = _charge_target - _charge_current
		if remaining > 0:
			_charge_label.text = "🐉 龍怒蓄力：%d / %d  （還差 %d 槍！）" % [_charge_current, _charge_target, remaining]
		else:
			_charge_label.text = "🐉 龍怒已滿！隕石雨即將降臨！"

	# 每 5 點額外閃光
	if _charge_current > 0 and _charge_current % 5 == 0:
		_flash_screen(Color("#FF4500"), 0.3)

	# 達到目標時強閃光
	if _charge_current >= _charge_target:
		_flash_screen(Color("#FF0000"), 0.7)
		await get_tree().create_timer(0.08).timeout
		_flash_screen(Color("#FFD700"), 0.5)

## 隕石雨開始 — 全螢幕爆炸 + 大字
func _on_meteor_start(payload: Dictionary, is_full: bool) -> void:
	# 清除蓄力進度條
	if is_instance_valid(_charge_bar):
		var tween = create_tween()
		tween.tween_property(_charge_bar, "modulate:a", 0.0, 0.3)
		tween.tween_callback(_charge_bar.queue_free)
		_charge_bar = null
		_charge_label = null

	if is_full:
		# 龍怒隕石雨：全螢幕三次深紅強閃光
		_flash_screen(Color("#FF0000"), 1.0)
		await get_tree().create_timer(0.1).timeout
		_flash_screen(Color("#FF4500"), 0.8)
		await get_tree().create_timer(0.1).timeout
		_flash_screen(Color("#FFD700"), 0.6)

		# 「🐉 龍怒隕石雨！」大字
		var big_label = Label.new()
		big_label.text = "🐉 龍怒隕石雨！"
		big_label.add_theme_font_size_override("font_size", 52)
		big_label.add_theme_color_override("font_color", Color("#FF4500"))
		big_label.set_anchors_preset(Control.PRESET_CENTER)
		big_label.position = Vector2(-220, -70)
		add_child(big_label)
		var tween = create_tween()
		tween.tween_property(big_label, "scale", Vector2(1.3, 1.3), 0.3)
		tween.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.2)
		tween.tween_interval(1.5)
		tween.tween_property(big_label, "modulate:a", 0.0, 0.6)
		tween.tween_callback(big_label.queue_free)
	else:
		# 小型龍怒：橙色閃光
		_flash_screen(Color("#FF6600"), 0.6)

		# 「🐉 小型龍怒！」大字
		var big_label = Label.new()
		big_label.text = "🐉 小型龍怒！"
		big_label.add_theme_font_size_override("font_size", 40)
		big_label.add_theme_color_override("font_color", Color("#FF6600"))
		big_label.set_anchors_preset(Control.PRESET_CENTER)
		big_label.position = Vector2(-180, -55)
		add_child(big_label)
		var tween = create_tween()
		tween.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.25)
		tween.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.2)
		tween.tween_interval(1.2)
		tween.tween_property(big_label, "modulate:a", 0.0, 0.5)
		tween.tween_callback(big_label.queue_free)

## 單顆隕石命中 — 爆炸圓圈 + 浮動獎勵文字
func _on_meteor_hit(payload: Dictionary) -> void:
	var mx = payload.get("meteor_x", 640.0)
	var my = payload.get("meteor_y", 360.0)
	var kills = payload.get("kill_count", 0)
	var reward = payload.get("reward", 0)
	var is_full = payload.get("is_full", false)

	# 隕石落點閃光
	var flash_color = Color("#FF4500") if is_full else Color("#FF6600")
	_flash_at(mx, my, flash_color, 0.5)

	# 爆炸圓圈擴散（雙圓圈）
	var radius = 350.0 if is_full else 250.0
	_spawn_explosion_ring(mx, my, radius, flash_color)

	# 浮動獎勵文字
	if kills > 0:
		var float_label = Label.new()
		float_label.text = "💥 ×%d  +%d" % [kills, reward]
		float_label.add_theme_font_size_override("font_size", 18)
		float_label.add_theme_color_override("font_color", Color("#FFD700"))
		float_label.position = Vector2(mx - 40, my - 20)
		add_child(float_label)
		var tween = create_tween()
		tween.tween_property(float_label, "position:y", my - 80.0, 0.8)
		tween.parallel().tween_property(float_label, "modulate:a", 0.0, 0.8)
		tween.tween_callback(float_label.queue_free)

## 隕石雨結算 — 右側滑入結算彈窗
func _on_meteor_result(payload: Dictionary, is_full: bool) -> void:
	var total_kills = payload.get("total_kills", 0)
	var total_reward = payload.get("total_reward", 0)

	var title = "🐉 龍怒隕石雨結算" if is_full else "🐉 小型龍怒結算"
	var border_color = Color("#FF4500") if is_full else Color("#FF6600")

	var popup = PanelContainer.new()
	popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	popup.position = Vector2(1400, -90)
	popup.size = Vector2(270, 150)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.0, 0.0, 0.92)
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
	title_lbl.text = title
	title_lbl.add_theme_font_size_override("font_size", 16)
	title_lbl.add_theme_color_override("font_color", border_color)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var kills_lbl = Label.new()
	kills_lbl.text = "擊破目標：%d 個" % total_kills
	kills_lbl.add_theme_font_size_override("font_size", 14)
	kills_lbl.add_theme_color_override("font_color", Color("#FF8C00"))
	kills_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(kills_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "全服獎勵：%d 金幣" % total_reward
	reward_lbl.add_theme_font_size_override("font_size", 14)
	reward_lbl.add_theme_color_override("font_color", Color("#FFFFFF"))
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	add_child(popup)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", 1000.0, 0.35)
	tween.tween_interval(4.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.6)
	tween.tween_callback(popup.queue_free)

## 底部蓄力進度條
func _show_charge_bar(charge_sec: int) -> void:
	if is_instance_valid(_charge_bar):
		_charge_bar.queue_free()

	var bar_container = Control.new()
	bar_container.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	bar_container.position = Vector2(0, -30)
	bar_container.size = Vector2(1280, 30)
	add_child(bar_container)
	_charge_bar = bar_container

	# 背景
	var bg = ColorRect.new()
	bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.1, 0.0, 0.0, 0.85)
	bar_container.add_child(bg)

	# 進度條（初始為 0）
	var bar = ColorRect.new()
	bar.name = "ChargeBar"
	bar.size = Vector2(0, 30)
	bar.color = Color("#FF4500")
	bar_container.add_child(bar)

	# 文字
	var label = Label.new()
	label.name = "ChargeLabel"
	label.text = "🐉 龍怒蓄力：0 / %d  （還差 %d 槍！）" % [_charge_target, _charge_target]
	label.add_theme_font_size_override("font_size", 14)
	label.add_theme_color_override("font_color", Color("#FFD700"))
	label.set_anchors_preset(Control.PRESET_FULL_RECT)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	bar_container.add_child(label)
	_charge_label = label

	# 超時後自動淡出
	var tween = create_tween()
	tween.tween_interval(float(charge_sec))
	tween.tween_property(bar_container, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(bar_container):
			bar_container.queue_free()
		_charge_bar = null
		_charge_label = null
	)

## 爆炸圓圈擴散
func _spawn_explosion_ring(cx: float, cy: float, max_radius: float, color: Color) -> void:
	# 外圈
	var ring = ColorRect.new()
	ring.size = Vector2(4, 4)
	ring.position = Vector2(cx - 2, cy - 2)
	ring.color = Color(color.r, color.g, color.b, 0.8)
	add_child(ring)
	var tween = create_tween()
	tween.tween_property(ring, "size", Vector2(max_radius * 2, max_radius * 2), 0.5)
	tween.parallel().tween_property(ring, "position", Vector2(cx - max_radius, cy - max_radius), 0.5)
	tween.parallel().tween_property(ring, "modulate:a", 0.0, 0.5)
	tween.tween_callback(ring.queue_free)

	# 內圈（較小，較快）
	var inner = ColorRect.new()
	inner.size = Vector2(4, 4)
	inner.position = Vector2(cx - 2, cy - 2)
	inner.color = Color(1.0, 0.9, 0.0, 0.6)
	add_child(inner)
	var tween2 = create_tween()
	tween2.tween_property(inner, "size", Vector2(max_radius, max_radius), 0.35)
	tween2.parallel().tween_property(inner, "position", Vector2(cx - max_radius * 0.5, cy - max_radius * 0.5), 0.35)
	tween2.parallel().tween_property(inner, "modulate:a", 0.0, 0.35)
	tween2.tween_callback(inner.queue_free)

## 局部閃光（隕石落點）
func _flash_at(x: float, y: float, color: Color, alpha: float) -> void:
	var flash = ColorRect.new()
	flash.size = Vector2(200, 200)
	flash.position = Vector2(x - 100, y - 100)
	flash.color = Color(color.r, color.g, color.b, alpha)
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.3)
	tween.tween_callback(flash.queue_free)

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

## 橫幅工廠
func _make_banner(text: String, bg_color: Color, text_color: Color) -> PanelContainer:
	var panel = PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_WIDE)
	panel.position = Vector2(0, -60)
	panel.size = Vector2(1280, 50)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(bg_color.r, bg_color.g, bg_color.b, 0.88)
	panel.add_theme_stylebox_override("panel", style)

	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 18)
	label.add_theme_color_override("font_color", text_color)
	label.set_anchors_preset(Control.PRESET_FULL_RECT)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	panel.add_child(label)
	return panel
