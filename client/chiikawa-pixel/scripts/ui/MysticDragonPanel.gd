## MysticDragonPanel.gd — 神秘龍魚八波攻擊面板（DAY-197）
## 業界依據：Ocean King 3「Mystic Dragon — Catch this fish to get 8 waves and have more
## chances to kill any fish on the screen.」
##
## 視覺設計：
##   - 神秘紫金主題（#8B00FF + #FFD700 + #FF69B4）
##   - dragon_start：紫色三次閃光 + 頂部橫幅 + 波數進度條（8格）
##   - wave_N：龍息掃場動畫（橫向紫色光波）+ 命中目標爆炸 + 波數計數器更新
##   - 第 8 波（龍怒爆發）：全螢幕金色強閃光 + 「🐉 龍怒爆發！」大字 + 全場爆炸
##   - dragon_result：右側滑入結算彈窗（波數/擊破數/全服共享獎勵）
extends CanvasLayer

var _panel: Control
var _banner: Label
var _wave_counter: Label      # 波數計數器
var _progress_bar: HBoxContainer  # 8格波數進度條
var _result_popup: Control
var _wave_cells: Array = []   # 8個進度格

func _ready() -> void:
	layer = 48
	_build_ui()

func _build_ui() -> void:
	_panel = Control.new()
	_panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	_panel.visible = false
	_panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_panel)

	# 頂部橫幅
	_banner = Label.new()
	_banner.text = "🐉 神秘龍魚八波攻擊！"
	_banner.add_theme_font_size_override("font_size", 24)
	_banner.add_theme_color_override("font_color", Color("#8B00FF"))
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.position = Vector2(0, 12)
	_banner.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.visible = false
	add_child(_banner)

	# 波數計數器
	_wave_counter = Label.new()
	_wave_counter.text = "第 0 / 8 波"
	_wave_counter.add_theme_font_size_override("font_size", 18)
	_wave_counter.add_theme_color_override("font_color", Color("#FFD700"))
	_wave_counter.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_wave_counter.position = Vector2(0, 42)
	_wave_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_wave_counter.visible = false
	add_child(_wave_counter)

	# 8格波數進度條（水平排列）
	_progress_bar = HBoxContainer.new()
	_progress_bar.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_progress_bar.position = Vector2(0, 68)
	_progress_bar.size = Vector2(get_viewport().get_visible_rect().size.x, 20)
	_progress_bar.alignment = BoxContainer.ALIGNMENT_CENTER
	_progress_bar.visible = false
	add_child(_progress_bar)

	for i in range(8):
		var cell = ColorRect.new()
		cell.size = Vector2(40, 18)
		cell.color = Color(0.3, 0.0, 0.5, 0.6)  # 未激活：暗紫色
		cell.custom_minimum_size = Vector2(40, 18)
		_progress_bar.add_child(cell)
		_wave_cells.append(cell)

		# 格子間距
		if i < 7:
			var spacer = Control.new()
			spacer.custom_minimum_size = Vector2(4, 0)
			_progress_bar.add_child(spacer)

	# 結算彈窗
	_result_popup = Control.new()
	_result_popup.visible = false
	_result_popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_popup.size = Vector2(300, 200)
	_result_popup.position = Vector2(-320, -100)
	add_child(_result_popup)

	var popup_bg = ColorRect.new()
	popup_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_bg.color = Color(0.08, 0.0, 0.18, 0.93)
	_result_popup.add_child(popup_bg)

	var popup_label = Label.new()
	popup_label.name = "ResultLabel"
	popup_label.add_theme_font_size_override("font_size", 16)
	popup_label.add_theme_color_override("font_color", Color("#FFD700"))
	popup_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_popup.add_child(popup_label)

## 處理神秘龍魚訊息
func handle_mystic_dragon(payload: Dictionary) -> void:
	var phase = payload.get("phase", "")
	match phase:
		"dragon_start":
			_on_dragon_start(payload)
		"dragon_result":
			_on_dragon_result(payload)
		_:
			if phase.begins_with("wave_"):
				_on_wave(payload)

func _on_dragon_start(payload: Dictionary) -> void:
	_panel.visible = true
	_panel.modulate.a = 1.0
	_banner.visible = true
	_wave_counter.visible = true
	_progress_bar.visible = true
	_wave_counter.text = "第 0 / 8 波"

	# 重置進度條
	for cell in _wave_cells:
		cell.color = Color(0.3, 0.0, 0.5, 0.6)

	# 紫色三次閃光
	_flash_screen(Color("#8B00FF"), 0.4)
	await get_tree().create_timer(0.45).timeout
	_flash_screen(Color("#8B00FF"), 0.3)
	await get_tree().create_timer(0.45).timeout
	_flash_screen(Color("#FF69B4"), 0.25)

	# 橫幅滑入動畫
	_banner.position.y = -30
	var tween = create_tween()
	tween.tween_property(_banner, "position:y", 12, 0.3)

func _on_wave(payload: Dictionary) -> void:
	var wave_index = payload.get("wave_index", 0)
	var wave_kills = payload.get("wave_kills", 0)
	var wave_reward = payload.get("wave_reward", 0)
	var total_kills = payload.get("total_kills", 0)
	var is_final_wave = payload.get("is_final_wave", false)

	# 更新波數計數器
	_wave_counter.text = "第 %d / 8 波" % wave_index

	# 更新進度條（激活對應格子）
	if wave_index >= 1 and wave_index <= 8:
		var cell = _wave_cells[wave_index - 1]
		if is_final_wave:
			cell.color = Color("#FFD700")  # 第 8 波：金色
		else:
			cell.color = Color("#8B00FF")  # 普通波：紫色

	if is_final_wave:
		# 第 8 波（龍怒爆發）：全螢幕金色強閃光
		_flash_screen(Color("#FFD700"), 0.6)
		await get_tree().create_timer(0.1).timeout
		_flash_screen(Color("#FF6600"), 0.4)

		# 「龍怒爆發！」大字
		_spawn_final_wave_text(wave_kills, wave_reward)

		# 橫幅更新
		_banner.text = "🐉💥 龍怒爆發！"
		_banner.add_theme_color_override("font_color", Color("#FFD700"))
	else:
		# 普通波：龍息掃場動畫
		_spawn_wave_sweep(wave_index)

		if wave_kills > 0:
			# 命中浮動文字
			var vp_size = get_viewport().get_visible_rect().size
			_spawn_reward_text(vp_size.x / 2, vp_size.y / 2 - 30, wave_reward)

func _on_dragon_result(payload: Dictionary) -> void:
	var total_kills = payload.get("total_kills", 0)
	var total_reward = payload.get("total_reward", 0)
	var killer_name = payload.get("killer_name", "")

	# 顯示結算彈窗（右側滑入）
	_result_popup.visible = true
	_result_popup.modulate.a = 0.0
	var vp_size = get_viewport().get_visible_rect().size
	_result_popup.position.x = vp_size.x

	var label = _result_popup.get_node("ResultLabel")
	label.text = "🐉 神秘龍魚\n八波攻擊完成\n擊破: %d 個目標\n全服共享: %d 金幣" % [total_kills, total_reward]

	var tween = create_tween()
	tween.tween_property(_result_popup, "modulate:a", 1.0, 0.3)
	tween.parallel().tween_property(_result_popup, "position:x", vp_size.x - 320, 0.3)

	# 依擊破數決定特效
	if total_kills >= 15:
		await get_tree().create_timer(0.3).timeout
		_flash_screen(Color("#FFD700"), 0.5)
	elif total_kills >= 8:
		await get_tree().create_timer(0.3).timeout
		_flash_screen(Color("#8B00FF"), 0.35)

	# 5 秒後淡出
	await get_tree().create_timer(5.0).timeout
	_fade_out()

## 龍息掃場動畫（橫向紫色光波）
func _spawn_wave_sweep(wave_index: int) -> void:
	var sweep = ColorRect.new()
	var vp_size = get_viewport().get_visible_rect().size
	sweep.size = Vector2(vp_size.x, 8)
	sweep.position = Vector2(0, vp_size.y * 0.3 + wave_index * 30)
	sweep.color = Color(0.55, 0.0, 1.0, 0.5)
	sweep.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_panel.add_child(sweep)

	var tween = create_tween()
	tween.tween_property(sweep, "scale:x", 0.0, 0.6)
	tween.parallel().tween_property(sweep, "modulate:a", 0.0, 0.6)
	await tween.finished
	sweep.queue_free()

## 第 8 波龍怒爆發大字
func _spawn_final_wave_text(kills: int, reward: int) -> void:
	var vp_size = get_viewport().get_visible_rect().size
	var label = Label.new()
	label.text = "🐉 龍怒爆發！\n擊破 %d 個！\n+%d 金幣" % [kills, reward]
	label.add_theme_font_size_override("font_size", 32)
	label.add_theme_color_override("font_color", Color("#FFD700"))
	label.set_anchors_preset(Control.PRESET_CENTER)
	label.position = Vector2(-120, -60)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(label)

	# 彈跳動畫
	var tween = create_tween()
	tween.tween_property(label, "scale", Vector2(1.2, 1.2), 0.2)
	tween.tween_property(label, "scale", Vector2(1.0, 1.0), 0.15)
	tween.tween_interval(1.5)
	tween.tween_property(label, "modulate:a", 0.0, 0.5)
	await tween.finished
	label.queue_free()

## 獎勵浮動文字
func _spawn_reward_text(x: float, y: float, reward: int) -> void:
	var label = Label.new()
	label.text = "+%d" % reward
	label.add_theme_font_size_override("font_size", 20)
	label.add_theme_color_override("font_color", Color("#8B00FF"))
	label.position = Vector2(x - 25, y)
	label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_panel.add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "position:y", y - 50, 0.7)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.7)
	await tween.finished
	label.queue_free()

func _fade_out() -> void:
	var tween = create_tween()
	tween.tween_property(_panel, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_banner, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_wave_counter, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_progress_bar, "modulate:a", 0.0, 0.5)
	await tween.finished
	_panel.visible = false
	_panel.modulate.a = 1.0
	_result_popup.visible = false
	_banner.visible = false
	_banner.modulate.a = 1.0
	_banner.text = "🐉 神秘龍魚八波攻擊！"
	_banner.add_theme_color_override("font_color", Color("#8B00FF"))
	_wave_counter.visible = false
	_wave_counter.modulate.a = 1.0
	_progress_bar.visible = false
	_progress_bar.modulate.a = 1.0

func _flash_screen(color: Color, intensity: float) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, intensity)
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.28)
	await tween.finished
	flash.queue_free()
