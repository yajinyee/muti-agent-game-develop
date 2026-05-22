## LuckyTridentPanel.gd — 幸運三叉魚互動三轉盤面板（DAY-211）
## 業界依據：TaDa Gaming TriLuck™ Series 2026
## 「Within the TriLuck™ Series, you can trigger three different feature
##  specifications, ranging from win multipliers, jackpot bonuses,
##  collecting all rewards, and more unique features.」
##
## 視覺設計：
##   - 三叉紫金主題（#9B59B6 + #FFD700 + #E74C3C + #2ECC71）
##   - trident_start：紫色強閃光 + 三個轉盤同時旋轉 + 「點擊停止！」提示
##   - trident_result：三個轉盤依序停止彈跳 + 結算彈窗（三重效果）
##   - trident_mult_end：倍率加成計時條淡出
##   - trident_effect：特效視覺（依類型不同）
##   - trident_broadcast：全服廣播橫幅
extends CanvasLayer

var _trident_panel: Control = null  # 三轉盤主面板
var _mult_bar: Control = null       # 倍率加成計時條
var _is_active: bool = false

func _ready() -> void:
	layer = 34  # 幸運三叉魚面板層級

## 處理幸運三叉魚訊息
func handle_lucky_trident(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"trident_start":
			_on_trident_start(payload)
		"trident_result":
			_on_trident_result(payload)
		"trident_mult_end":
			_on_trident_mult_end()
		"trident_effect":
			_on_trident_effect(payload)
		"trident_effect_end":
			_on_trident_effect_end(payload)
		"trident_broadcast":
			_on_trident_broadcast(payload)

## 三叉儀式開始 — 紫色強閃光 + 三個轉盤
func _on_trident_start(payload: Dictionary) -> void:
	var timeout: int = payload.get("timeout_sec", 12)

	# 紫色強閃光
	_flash_screen(Color("#9B59B6"), 0.75)

	# 建立三轉盤面板
	_create_trident_panel(timeout)
	_is_active = true

## 三叉結算 — 三個轉盤停止 + 結算彈窗
func _on_trident_result(payload: Dictionary) -> void:
	var a_label: String = payload.get("wheel_a_label", "💰 ×10")
	var b_label: String = payload.get("wheel_b_label", "⚡ ×1.5")
	var c_label: String = payload.get("wheel_c_label", "🩸 HP削減")
	var coin: int = payload.get("coin_reward", 0)
	var mult: float = payload.get("mult_boost", 1.5)
	var mult_sec: int = payload.get("mult_sec", 15)
	var effect: String = payload.get("effect", "")
	var effect_desc: String = payload.get("effect_desc", "")
	var is_timeout: bool = payload.get("is_timeout", false)

	# 移除轉盤面板
	if is_instance_valid(_trident_panel):
		_trident_panel.queue_free()
		_trident_panel = null
	_is_active = false

	# 三次閃光（依序）
	_flash_screen(Color("#FFD700"), 0.6)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(Color("#9B59B6"), 0.5)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(Color("#2ECC71"), 0.4)

	# 結算彈窗
	_show_result_popup(a_label, b_label, c_label, coin, mult, mult_sec, effect_desc, is_timeout)

	# 倍率加成計時條
	_show_mult_bar(mult, mult_sec)

## 倍率加成結束
func _on_trident_mult_end() -> void:
	if is_instance_valid(_mult_bar):
		var tween = create_tween()
		tween.tween_property(_mult_bar, "modulate:a", 0.0, 0.4)
		tween.tween_callback(_mult_bar.queue_free)
		_mult_bar = null

## 特效視覺
func _on_trident_effect(payload: Dictionary) -> void:
	var effect: String = payload.get("effect", "")
	var affected: int = payload.get("affected_count", 0)

	match effect:
		"hp_drain":
			# 全場 HP -30% 視覺
			_flash_screen(Color("#E74C3C"), 0.5)
			var label = Label.new()
			label.text = "🩸 全場 HP -30%%！%d 個目標" % affected
			label.add_theme_font_size_override("font_size", 22)
			label.add_theme_color_override("font_color", Color("#E74C3C"))
			label.set_anchors_preset(Control.PRESET_CENTER)
			label.position = Vector2(-180, -40)
			add_child(label)
			var tween = create_tween()
			tween.tween_property(label, "position:y", label.position.y - 30, 0.5)
			tween.parallel().tween_property(label, "modulate:a", 0.0, 0.5)
			tween.tween_callback(label.queue_free)

		"mini_blast":
			# 小型清場視覺
			_flash_screen(Color("#E74C3C"), 0.65)
			var label = Label.new()
			label.text = "💥 小型清場！擊破 %d 個目標！" % affected
			label.add_theme_font_size_override("font_size", 22)
			label.add_theme_color_override("font_color", Color("#FF4500"))
			label.set_anchors_preset(Control.PRESET_CENTER)
			label.position = Vector2(-200, -40)
			add_child(label)
			var tween = create_tween()
			tween.tween_property(label, "position:y", label.position.y - 30, 0.5)
			tween.parallel().tween_property(label, "modulate:a", 0.0, 0.5)
			tween.tween_callback(label.queue_free)

		"free_shot":
			# 免費射擊開始視覺
			_flash_screen(Color("#2ECC71"), 0.6)
			var label = Label.new()
			label.text = "🎯 免費射擊！5 秒！"
			label.add_theme_font_size_override("font_size", 26)
			label.add_theme_color_override("font_color", Color("#2ECC71"))
			label.set_anchors_preset(Control.PRESET_CENTER)
			label.position = Vector2(-160, -40)
			add_child(label)
			var tween = create_tween()
			tween.tween_property(label, "scale", Vector2(1.2, 1.2), 0.2)
			tween.tween_property(label, "scale", Vector2(1.0, 1.0), 0.15)
			tween.tween_interval(1.5)
			tween.tween_property(label, "modulate:a", 0.0, 0.4)
			tween.tween_callback(label.queue_free)

## 特效結束
func _on_trident_effect_end(payload: Dictionary) -> void:
	var effect: String = payload.get("effect", "")
	if effect == "free_shot":
		var label = Label.new()
		label.text = "🎯 免費射擊結束"
		label.add_theme_font_size_override("font_size", 18)
		label.add_theme_color_override("font_color", Color("#95A5A6"))
		label.set_anchors_preset(Control.PRESET_CENTER)
		label.position = Vector2(-120, -30)
		add_child(label)
		var tween = create_tween()
		tween.tween_interval(1.0)
		tween.tween_property(label, "modulate:a", 0.0, 0.4)
		tween.tween_callback(label.queue_free)

## 全服廣播橫幅
func _on_trident_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "玩家")
	var banner = _make_banner(
		"🔱 %s 觸發三叉幸運儀式！" % player_name,
		Color(0.1, 0.05, 0.15, 0.85),
		Color("#9B59B6")
	)
	add_child(banner)
	var tween = create_tween()
	tween.tween_property(banner, "position:y", 0.0, 0.25)
	tween.tween_interval(2.5)
	tween.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(banner.queue_free)

## 建立三轉盤主面板
func _create_trident_panel(timeout: int) -> void:
	if is_instance_valid(_trident_panel):
		_trident_panel.queue_free()

	var panel = PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_CENTER)
	panel.position = Vector2(-320, -180)
	panel.size = Vector2(640, 360)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.02, 0.1, 0.95)
	style.border_color = Color("#9B59B6")
	style.set_border_width_all(3)
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	panel.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.add_child(vbox)

	# 標題
	var title = Label.new()
	title.text = "🔱 三叉幸運儀式"
	title.add_theme_font_size_override("font_size", 24)
	title.add_theme_color_override("font_color", Color("#9B59B6"))
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title)

	var sub = Label.new()
	sub.text = "點擊停止轉盤！（%d 秒後自動停止）" % timeout
	sub.add_theme_font_size_override("font_size", 14)
	sub.add_theme_color_override("font_color", Color("#BDC3C7"))
	sub.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(sub)

	# 三個轉盤橫排
	var hbox = HBoxContainer.new()
	hbox.set_h_size_flags(Control.SIZE_EXPAND_FILL)
	hbox.alignment = BoxContainer.ALIGNMENT_CENTER
	vbox.add_child(hbox)

	var wheel_colors = [Color("#FFD700"), Color("#9B59B6"), Color("#2ECC71")]
	var wheel_labels = ["💰 金幣", "⚡ 倍率", "✨ 特效"]
	var wheel_items = [
		["×10", "×20", "×30", "×50", "×100"],
		["×1.5", "×2.0", "×2.5", "×3.0", "×5.0"],
		["HP削減", "免費射擊", "全服廣播", "小型清場"],
	]

	for i in range(3):
		var wheel_container = _make_wheel(wheel_labels[i], wheel_colors[i], wheel_items[i])
		hbox.add_child(wheel_container)

	# 停止按鈕
	var stop_btn = Button.new()
	stop_btn.text = "🛑 停止轉盤"
	stop_btn.add_theme_font_size_override("font_size", 18)
	stop_btn.pressed.connect(_on_stop_pressed)
	vbox.add_child(stop_btn)

	add_child(panel)
	_trident_panel = panel

	# 淡入動畫
	panel.modulate.a = 0.0
	var tween = create_tween()
	tween.tween_property(panel, "modulate:a", 1.0, 0.3)

## 建立單個轉盤
func _make_wheel(label_text: String, color: Color, items: Array) -> VBoxContainer:
	var vbox = VBoxContainer.new()
	vbox.custom_minimum_size = Vector2(180, 200)

	var title = Label.new()
	title.text = label_text
	title.add_theme_font_size_override("font_size", 16)
	title.add_theme_color_override("font_color", color)
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title)

	# 轉盤顯示框（旋轉的項目）
	var frame = PanelContainer.new()
	frame.custom_minimum_size = Vector2(160, 120)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.08, 0.04, 0.12, 0.9)
	style.border_color = color
	style.set_border_width_all(2)
	style.corner_radius_top_left = 6
	style.corner_radius_top_right = 6
	style.corner_radius_bottom_left = 6
	style.corner_radius_bottom_right = 6
	frame.add_theme_stylebox_override("panel", style)

	var item_label = Label.new()
	item_label.text = items[0]
	item_label.add_theme_font_size_override("font_size", 22)
	item_label.add_theme_color_override("font_color", color)
	item_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	item_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	item_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	frame.add_child(item_label)
	vbox.add_child(frame)

	# 旋轉動畫（模擬轉盤效果）
	var idx = 0
	var timer = get_tree().create_timer(0.0)
	var spin_func = func():
		pass
	spin_func = func():
		if not is_instance_valid(item_label):
			return
		idx = (idx + 1) % items.size()
		item_label.text = items[idx]
		var t = create_tween()
		t.tween_property(item_label, "modulate:a", 0.3, 0.05)
		t.tween_property(item_label, "modulate:a", 1.0, 0.05)
		get_tree().create_timer(0.12).timeout.connect(spin_func, CONNECT_ONE_SHOT)

	get_tree().create_timer(0.12).timeout.connect(spin_func, CONNECT_ONE_SHOT)

	return vbox

## 停止按鈕點擊
func _on_stop_pressed() -> void:
	# 發送停止訊號給 Server
	if NetworkManager.has_method("send_message"):
		NetworkManager.send_message("lucky_trident_stop", {"player_id": ""})

## 結算彈窗
func _show_result_popup(a: String, b: String, c: String, coin: int, mult: float, mult_sec: int, effect_desc: String, is_timeout: bool) -> void:
	var popup = PanelContainer.new()
	popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	popup.position = Vector2(1400, -130)
	popup.size = Vector2(280, 260)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.02, 0.1, 0.95)
	style.border_color = Color("#9B59B6")
	style.set_border_width_all(3)
	style.corner_radius_top_left = 10
	style.corner_radius_top_right = 10
	style.corner_radius_bottom_left = 10
	style.corner_radius_bottom_right = 10
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup.add_child(vbox)

	var title = Label.new()
	title.text = "🔱 三叉幸運結算" + (" (超時)" if is_timeout else "")
	title.add_theme_font_size_override("font_size", 16)
	title.add_theme_color_override("font_color", Color("#9B59B6"))
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title)

	for item in [[a, Color("#FFD700")], [b, Color("#9B59B6")], [c, Color("#2ECC71")]]:
		var lbl = Label.new()
		lbl.text = item[0]
		lbl.add_theme_font_size_override("font_size", 20)
		lbl.add_theme_color_override("font_color", item[1])
		lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		vbox.add_child(lbl)

	var coin_lbl = Label.new()
	coin_lbl.text = "+%d 金幣" % coin
	coin_lbl.add_theme_font_size_override("font_size", 16)
	coin_lbl.add_theme_color_override("font_color", Color("#FFFFFF"))
	coin_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(coin_lbl)

	var effect_lbl = Label.new()
	effect_lbl.text = effect_desc
	effect_lbl.add_theme_font_size_override("font_size", 13)
	effect_lbl.add_theme_color_override("font_color", Color("#BDC3C7"))
	effect_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(effect_lbl)

	add_child(popup)

	var tween = create_tween()
	tween.tween_property(popup, "position:x", 1010.0, 0.3)
	tween.tween_interval(4.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)

## 倍率加成計時條
func _show_mult_bar(mult: float, mult_sec: int) -> void:
	if is_instance_valid(_mult_bar):
		_mult_bar.queue_free()

	var bar_container = Control.new()
	bar_container.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	bar_container.position = Vector2(0, -30)
	bar_container.size = Vector2(1280, 30)
	add_child(bar_container)
	_mult_bar = bar_container

	var bg = ColorRect.new()
	bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.05, 0.02, 0.1, 0.85)
	bar_container.add_child(bg)

	var bar = ColorRect.new()
	bar.set_anchors_preset(Control.PRESET_FULL_RECT)
	bar.color = Color("#9B59B6")
	bar_container.add_child(bar)

	var label = Label.new()
	label.text = "⚡ 三叉倍率加成 ×%.1f 進行中！" % mult
	label.add_theme_font_size_override("font_size", 14)
	label.add_theme_color_override("font_color", Color("#FFFFFF"))
	label.set_anchors_preset(Control.PRESET_FULL_RECT)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	bar_container.add_child(label)

	var tween = create_tween()
	tween.tween_property(bar, "size:x", 0.0, float(mult_sec))
	tween.parallel().tween_method(
		func(t: float):
			if is_instance_valid(bar):
				if t > 0.5:
					bar.color = Color("#9B59B6")
				elif t > 0.25:
					bar.color = Color("#8E44AD")
				else:
					bar.color = Color("#6C3483"),
		1.0, 0.0, float(mult_sec)
	)

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
