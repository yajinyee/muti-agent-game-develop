## LuckyChainExplosionPanel.gd — 幸運連鎖爆炸魚 UI（DAY-266）
## 火焰爆炸主題面板
## 主色：#FF4500 火橙 + #FF6B35 橙紅 + #FFD700 金 + #FF0000 紅
extends CanvasLayer

const EXPLOSION_COLOR_FIRE   = Color(1.0, 0.27, 0.0, 1.0)   # #FF4500 火橙
const EXPLOSION_COLOR_ORANGE = Color(1.0, 0.42, 0.21, 1.0)  # #FF6B35 橙紅
const EXPLOSION_COLOR_GOLD   = Color(1.0, 0.85, 0.0, 1.0)   # #FFD700 金
const EXPLOSION_COLOR_RED    = Color(1.0, 0.0, 0.0, 1.0)    # #FF0000 紅

var _flash_overlay: ColorRect
var _banner: PanelContainer
var _banner_label: Label
var _layer_counter: Label

func _ready() -> void:
	layer = 39
	_build_ui()

func _build_ui() -> void:
	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 0.27, 0.0, 0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner = PanelContainer.new()
	_banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 56)
	_banner.modulate.a = 0.0
	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.15, 0.04, 0.0, 0.88)
	banner_style.border_color = EXPLOSION_COLOR_FIRE
	banner_style.set_border_width_all(2)
	banner_style.set_corner_radius_all(6)
	_banner.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_font_size_override("font_size", 18)
	_banner_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.8))
	_banner.add_child(_banner_label)

	# 連鎖層數計數器（右上角）
	var vp_size = get_viewport().size
	_layer_counter = Label.new()
	_layer_counter.text = "💥 第 1 層"
	_layer_counter.position = Vector2(vp_size.x - 160, 65)
	_layer_counter.add_theme_font_size_override("font_size", 20)
	_layer_counter.add_theme_color_override("font_color", EXPLOSION_COLOR_FIRE)
	_layer_counter.modulate.a = 0.0
	add_child(_layer_counter)

## 處理來自 GameManager 的連鎖爆炸事件
func handle_event(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"explosion_start":
			_on_explosion_start(payload)
		"explosion_broadcast":
			_on_explosion_broadcast(payload)
		"explosion_layer":
			_on_explosion_layer(payload)
		"explosion_result":
			_on_explosion_result(payload)

## 連鎖爆炸觸發（個人）
func _on_explosion_start(payload: Dictionary) -> void:
	var layer1_mult = payload.get("layer1_mult", 2.0)

	# 火橙三次強閃光
	_flash_triple(Color(1.0, 0.27, 0.0, 0.5))

	# 顯示橫幅
	_banner_label.text = "💥 連鎖爆炸！最多 3 層連鎖！第 1 層 ×%.1f！" % layer1_mult
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)
	tween.tween_interval(4.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.5)

	# 顯示層數計數器
	_layer_counter.text = "💥 第 1 層"
	_layer_counter.add_theme_color_override("font_color", EXPLOSION_COLOR_FIRE)
	_layer_counter.modulate.a = 1.0

	# 大字提示
	_show_big_text("💥 連鎖爆炸！", EXPLOSION_COLOR_FIRE)

## 全服廣播
func _on_explosion_broadcast(payload: Dictionary) -> void:
	var trigger = payload.get("player_name", "玩家")
	var mult = payload.get("layer1_mult", 2.0)
	_banner_label.text = "💥 %s 觸發連鎖爆炸！×%.1f！" % [trigger, mult]
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)
	tween.tween_interval(3.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.3)

## 每層爆炸廣播
func _on_explosion_layer(payload: Dictionary) -> void:
	var layer_num = payload.get("layer", 1)
	var mult = payload.get("mult", 1.0)
	var reward = payload.get("reward", 0)
	var target_name = payload.get("target_name", "")
	var x = payload.get("x", 0.0)
	var y = payload.get("y", 0.0)
	var explode_count = payload.get("explode_count", 0)
	var affected_count = payload.get("affected_count", 0)

	# 更新層數計數器
	var layer_colors = [EXPLOSION_COLOR_FIRE, EXPLOSION_COLOR_ORANGE, EXPLOSION_COLOR_GOLD]
	var layer_color = layer_colors[min(layer_num - 1, 2)]
	_layer_counter.text = "💥 第 %d 層" % layer_num
	_layer_counter.add_theme_color_override("font_color", layer_color)

	# 爆炸閃光（層數越高越強）
	var flash_alpha = 0.3 + (layer_num - 1) * 0.1
	_flash_overlay.color = Color(layer_color.r, layer_color.g, layer_color.b, 0)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color:a", flash_alpha, 0.05)
	tween.tween_property(_flash_overlay, "color:a", 0.0, 0.12)

	# 爆炸圓圈特效（在爆炸位置）
	if x > 0 or y > 0:
		_show_explosion_circle(x, y, layer_num, layer_color)

	# 浮動文字
	if layer_num == 1 and target_name != "":
		_show_float_text("💥 %s 引爆！×%.1f +%d" % [target_name, mult, reward], layer_color)
	elif layer_num == 2 and explode_count > 0:
		_show_float_text("💥 第 2 層！%d 個引爆！×%.1f +%d" % [explode_count, mult, reward], layer_color)
	elif layer_num == 3 and explode_count > 0:
		_show_float_text("💥 第 3 層！%d 個引爆！×%.1f +%d" % [explode_count, mult, reward], layer_color)

	# 層數計數器脈衝
	var pulse = create_tween()
	pulse.tween_property(_layer_counter, "scale", Vector2(1.3, 1.3), 0.08)
	pulse.tween_property(_layer_counter, "scale", Vector2(1.0, 1.0), 0.12)

	_ = affected_count

## 最終結算廣播
func _on_explosion_result(payload: Dictionary) -> void:
	var total_layers = payload.get("total_layers", 1)
	var total_explode = payload.get("total_explode", 1)
	var total_reward = payload.get("total_reward", 0)
	var player_name = payload.get("player_name", "玩家")

	# 3 層連鎖時全螢幕強閃光
	if total_layers >= 3:
		_flash_triple(Color(1.0, 0.85, 0.0, 0.7))
		_show_big_text("💥 3 層連鎖！%d 個引爆！" % total_explode, EXPLOSION_COLOR_GOLD)
	elif total_layers >= 2:
		_flash_triple(Color(1.0, 0.42, 0.21, 0.5))
		_show_big_text("💥 2 層連鎖！%d 個引爆！" % total_explode, EXPLOSION_COLOR_ORANGE)

	# 結算彈窗
	_show_result_popup(player_name, total_layers, total_explode, total_reward)

	# 隱藏層數計數器
	var tween = create_tween()
	tween.tween_interval(0.5)
	tween.tween_property(_layer_counter, "modulate:a", 0.0, 0.5)

## 工具函數：爆炸圓圈特效
func _show_explosion_circle(cx: float, cy: float, layer_num: int, color: Color) -> void:
	var circle = ColorRect.new()
	var base_size = 40.0 + layer_num * 20.0
	circle.color = Color(color.r, color.g, color.b, 0.5)
	circle.size = Vector2(base_size, base_size)
	circle.position = Vector2(cx - base_size / 2, cy - base_size / 2)
	add_child(circle)

	var tween = create_tween()
	tween.tween_property(circle, "size", Vector2(base_size * 3, base_size * 3), 0.3)
	tween.parallel().tween_property(circle, "position",
		Vector2(cx - base_size * 1.5, cy - base_size * 1.5), 0.3)
	tween.parallel().tween_property(circle, "color:a", 0.0, 0.3)
	tween.tween_callback(circle.queue_free)

## 工具函數：三次強閃光
func _flash_triple(color: Color) -> void:
	var tween = create_tween()
	for i in range(3):
		tween.tween_property(_flash_overlay, "color:a", color.a, 0.06)
		tween.tween_property(_flash_overlay, "color:a", 0.0, 0.1)
	_flash_overlay.color = Color(color.r, color.g, color.b, 0)

## 工具函數：顯示大字
func _show_big_text(text: String, color: Color) -> void:
	var label = Label.new()
	label.text = text
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.add_theme_font_size_override("font_size", 30)
	label.add_theme_color_override("font_color", color)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.position.y -= 80
	label.modulate.a = 0.0
	add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "modulate:a", 1.0, 0.15)
	tween.tween_property(label, "position:y", label.position.y - 40, 1.2)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.4).set_delay(0.8)
	tween.tween_callback(label.queue_free)

## 工具函數：浮動文字
func _show_float_text(text: String, color: Color) -> void:
	var label = Label.new()
	label.text = text
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.add_theme_font_size_override("font_size", 20)
	label.add_theme_color_override("font_color", color)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.position.y -= 40
	label.modulate.a = 0.0
	add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "modulate:a", 1.0, 0.1)
	tween.tween_property(label, "position:y", label.position.y - 30, 0.8)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.3).set_delay(0.5)
	tween.tween_callback(label.queue_free)

## 工具函數：顯示結算彈窗
func _show_result_popup(player_name: String, total_layers: int, total_explode: int, total_reward: int) -> void:
	var popup = PanelContainer.new()
	popup.custom_minimum_size = Vector2(280, 130)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.03, 0.0, 0.92)
	var border_color = EXPLOSION_COLOR_FIRE if total_layers < 3 else EXPLOSION_COLOR_GOLD
	style.border_color = border_color
	style.set_border_width_all(2)
	style.set_corner_radius_all(8)
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	popup.add_child(vbox)

	var title_label = Label.new()
	title_label.text = "💥 連鎖爆炸結算！"
	title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_label.add_theme_font_size_override("font_size", 18)
	title_label.add_theme_color_override("font_color", border_color)
	vbox.add_child(title_label)

	var trigger_label = Label.new()
	trigger_label.text = "觸發者：%s" % player_name
	trigger_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	trigger_label.add_theme_font_size_override("font_size", 14)
	trigger_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.8))
	vbox.add_child(trigger_label)

	var layers_label = Label.new()
	layers_label.text = "連鎖層數：%d 層 / 共 %d 個引爆" % [total_layers, total_explode]
	layers_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	layers_label.add_theme_font_size_override("font_size", 15)
	layers_label.add_theme_color_override("font_color", EXPLOSION_COLOR_ORANGE)
	vbox.add_child(layers_label)

	var reward_label = Label.new()
	reward_label.text = "總獎勵：+%d 🪙" % total_reward
	reward_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_label.add_theme_font_size_override("font_size", 16)
	reward_label.add_theme_color_override("font_color", Color(0.0, 1.0, 0.53))
	vbox.add_child(reward_label)

	var vp_size = get_viewport().size
	popup.position = Vector2(vp_size.x, vp_size.y / 2 - 65)
	popup.modulate.a = 0.0
	add_child(popup)

	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 300, 0.3)
	tween.parallel().tween_property(popup, "modulate:a", 1.0, 0.3)
	tween.tween_interval(4.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.4)
	tween.tween_callback(popup.queue_free)
