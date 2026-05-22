## GhostFishPanel.gd — 幽靈魚分身面板（DAY-198）
## 原創設計：幽靈魚出現時生成 2-3 個幻影分身，玩家需找出真身觸發幻影爆炸
##
## 視覺設計：
##   - 幽靈白藍主題（#E0E0FF + #8888FF + #FFFFFF）
##   - ghost_appear：白色閃光 + 頂部橫幅「👻 幽靈魚出現！哪個是真身？」+ 分身計數
##   - phantom_vanish：幻影消散動畫（半透明→消失）+ 「幻影！+N 金幣」浮動文字
##   - real_found：金色強閃光 + 「找到真身！」大字 + 橫幅更新
##   - ghost_explode：白色爆炸圓圈 + 擊破數/獎勵結算彈窗
##   - ghost_escape：幽靈逃跑動畫（淡出）
extends CanvasLayer

var _panel: Control
var _banner: Label
var _clone_counter: Label    # 分身計數器
var _result_popup: Control

func _ready() -> void:
	layer = 47
	_build_ui()

func _build_ui() -> void:
	_panel = Control.new()
	_panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	_panel.visible = false
	_panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_panel)

	# 頂部橫幅
	_banner = Label.new()
	_banner.text = "👻 幽靈魚出現！哪個是真身？"
	_banner.add_theme_font_size_override("font_size", 22)
	_banner.add_theme_color_override("font_color", Color("#8888FF"))
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.position = Vector2(0, 12)
	_banner.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.visible = false
	add_child(_banner)

	# 分身計數器
	_clone_counter = Label.new()
	_clone_counter.text = "場上分身: 0 個"
	_clone_counter.add_theme_font_size_override("font_size", 16)
	_clone_counter.add_theme_color_override("font_color", Color("#E0E0FF"))
	_clone_counter.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_clone_counter.position = Vector2(0, 40)
	_clone_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_clone_counter.visible = false
	add_child(_clone_counter)

	# 結算彈窗
	_result_popup = Control.new()
	_result_popup.visible = false
	_result_popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_popup.size = Vector2(280, 180)
	_result_popup.position = Vector2(-300, -90)
	add_child(_result_popup)

	var popup_bg = ColorRect.new()
	popup_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_bg.color = Color(0.05, 0.05, 0.2, 0.93)
	_result_popup.add_child(popup_bg)

	var popup_label = Label.new()
	popup_label.name = "ResultLabel"
	popup_label.add_theme_font_size_override("font_size", 16)
	popup_label.add_theme_color_override("font_color", Color("#E0E0FF"))
	popup_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_popup.add_child(popup_label)

## 處理幽靈魚訊息
func handle_ghost_fish(payload: Dictionary) -> void:
	var phase = payload.get("phase", "")
	match phase:
		"ghost_appear":
			_on_appear(payload)
		"phantom_vanish":
			_on_phantom_vanish(payload)
		"real_found":
			_on_real_found(payload)
		"ghost_explode":
			_on_explode(payload)
		"ghost_escape":
			_on_escape(payload)

func _on_appear(payload: Dictionary) -> void:
	var clone_count = payload.get("clone_count", 2)
	_panel.visible = true
	_panel.modulate.a = 1.0
	_banner.visible = true
	_clone_counter.visible = true
	_clone_counter.text = "場上分身: %d 個（含真身共 %d 個）" % [clone_count, clone_count + 1]

	# 白色閃光（幽靈感）
	_flash_screen(Color("#FFFFFF"), 0.35)
	await get_tree().create_timer(0.4).timeout
	_flash_screen(Color("#8888FF"), 0.25)

	# 橫幅閃爍（製造懸疑感）
	var tween = create_tween().set_loops(3)
	tween.tween_property(_banner, "modulate:a", 0.3, 0.3)
	tween.tween_property(_banner, "modulate:a", 1.0, 0.3)

func _on_phantom_vanish(payload: Dictionary) -> void:
	var killer_name = payload.get("killer_name", "")
	var reward = payload.get("reward", 0)

	# 幻影消散浮動文字
	var vp_size = get_viewport().get_visible_rect().size
	_spawn_phantom_text(vp_size.x / 2, vp_size.y / 2, reward)

func _on_real_found(payload: Dictionary) -> void:
	var killer_name = payload.get("killer_name", "")

	# 金色強閃光
	_flash_screen(Color("#FFD700"), 0.5)

	# 橫幅更新
	_banner.text = "👻💥 找到真身！幻影爆炸！"
	_banner.add_theme_color_override("font_color", Color("#FFD700"))

	# 「找到真身！」大字
	var vp_size = get_viewport().get_visible_rect().size
	var label = Label.new()
	label.text = "👻 找到真身！"
	label.add_theme_font_size_override("font_size", 30)
	label.add_theme_color_override("font_color", Color("#FFD700"))
	label.set_anchors_preset(Control.PRESET_CENTER)
	label.position = Vector2(-80, -30)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "scale", Vector2(1.2, 1.2), 0.15)
	tween.tween_property(label, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(1.0)
	tween.tween_property(label, "modulate:a", 0.0, 0.4)
	await tween.finished
	label.queue_free()

func _on_explode(payload: Dictionary) -> void:
	var explode_kills = payload.get("explode_kills", 0)
	var explode_reward = payload.get("explode_reward", 0)
	var killer_name = payload.get("killer_name", "")

	# 白色爆炸圓圈
	_spawn_explode_effect()

	# 顯示結算彈窗
	_result_popup.visible = true
	_result_popup.modulate.a = 0.0
	var vp_size = get_viewport().get_visible_rect().size
	_result_popup.position.x = vp_size.x

	var label = _result_popup.get_node("ResultLabel")
	label.text = "👻 幽靈爆炸！\n幻影擊破: %d 個\n獎勵: %d 金幣" % [explode_kills, explode_reward]

	var tween = create_tween()
	tween.tween_property(_result_popup, "modulate:a", 1.0, 0.3)
	tween.parallel().tween_property(_result_popup, "position:x", vp_size.x - 300, 0.3)

	# 4 秒後淡出
	await get_tree().create_timer(4.0).timeout
	_fade_out()

func _on_escape(payload: Dictionary) -> void:
	# 幽靈魚逃跑，淡出 UI
	_banner.text = "👻 幽靈魚逃跑了..."
	_banner.add_theme_color_override("font_color", Color("#888888"))
	await get_tree().create_timer(1.5).timeout
	_fade_out()

## 幻影消散浮動文字
func _spawn_phantom_text(x: float, y: float, reward: int) -> void:
	var label = Label.new()
	label.text = "幻影！+%d 金幣" % reward
	label.add_theme_font_size_override("font_size", 18)
	label.add_theme_color_override("font_color", Color("#8888FF"))
	label.position = Vector2(x - 60, y - 20)
	label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_panel.add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "position:y", y - 60, 0.7)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.7)
	await tween.finished
	label.queue_free()

## 爆炸效果
func _spawn_explode_effect() -> void:
	var vp_size = get_viewport().get_visible_rect().size
	var effect = Control.new()
	effect.position = Vector2(vp_size.x / 2, vp_size.y / 2)
	effect.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_panel.add_child(effect)

	var circle = ColorRect.new()
	circle.size = Vector2(200, 200)
	circle.position = Vector2(-100, -100)
	circle.color = Color(0.9, 0.9, 1.0, 0.5)
	circle.mouse_filter = Control.MOUSE_FILTER_IGNORE
	effect.add_child(circle)

	var tween = create_tween()
	tween.tween_property(circle, "scale", Vector2(3.0, 3.0), 0.5)
	tween.parallel().tween_property(circle, "modulate:a", 0.0, 0.5)
	await tween.finished
	effect.queue_free()

func _fade_out() -> void:
	var tween = create_tween()
	tween.tween_property(_panel, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_banner, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_clone_counter, "modulate:a", 0.0, 0.5)
	await tween.finished
	_panel.visible = false
	_panel.modulate.a = 1.0
	_result_popup.visible = false
	_banner.visible = false
	_banner.modulate.a = 1.0
	_banner.text = "👻 幽靈魚出現！哪個是真身？"
	_banner.add_theme_color_override("font_color", Color("#8888FF"))
	_clone_counter.visible = false
	_clone_counter.modulate.a = 1.0

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
