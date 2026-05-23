## LuckyElementFusionPanel.gd — 幸運元素融合魚 UI（DAY-263）
## 火/水/風三元素主題面板
## 元素顏色：火=紅橙(#FF4500)、水=天藍(#00BFFF)、風=翠綠(#00FF88)
extends CanvasLayer

const ELEMENT_COLORS = {
	"fire":  Color(1.0, 0.27, 0.0, 1.0),   # #FF4500 火紅
	"water": Color(0.0, 0.75, 1.0, 1.0),   # #00BFFF 天藍
	"wind":  Color(0.0, 1.0, 0.53, 1.0),   # #00FF88 翠綠
}
const ELEMENT_ICONS = {
	"fire":  "🔥",
	"water": "💧",
	"wind":  "🌪️",
}
const ELEMENT_NAMES = {
	"fire":  "火元素",
	"water": "水元素",
	"wind":  "風元素",
}

var _flash_overlay: ColorRect
var _banner: PanelContainer
var _banner_label: Label
var _collect_label: Label
var _collect_timer: float = 0.0
var _fragment_markers: Dictionary = {}  # instanceID -> Control

func _ready() -> void:
	layer = 36
	_build_ui()

func _build_ui() -> void:
	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1, 0.42, 0.21, 0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner = PanelContainer.new()
	_banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 56)
	_banner.modulate.a = 0.0
	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.1, 0.05, 0.0, 0.88)
	banner_style.border_color = Color(1.0, 0.42, 0.21, 1.0)
	banner_style.set_border_width_all(2)
	banner_style.set_corner_radius_all(6)
	_banner.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_font_size_override("font_size", 20)
	_banner_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_banner.add_child(_banner_label)

	# 收集提示浮動文字
	_collect_label = Label.new()
	_collect_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_collect_label.add_theme_font_size_override("font_size", 22)
	_collect_label.add_theme_color_override("font_color", Color(1, 1, 1))
	_collect_label.position = Vector2(get_viewport().size.x / 2 - 120, get_viewport().size.y / 2 - 60)
	_collect_label.modulate.a = 0.0
	add_child(_collect_label)

func _process(delta: float) -> void:
	if _collect_timer > 0:
		_collect_timer -= delta
		if _collect_timer <= 0:
			_collect_label.modulate.a = 0.0

## 處理來自 GameManager 的元素融合事件
func handle_event(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"fusion_start":
			_on_fusion_start(payload)
		"fusion_broadcast":
			_on_fusion_broadcast(payload)
		"fragment_collect":
			_on_fragment_collect(payload)
		"fusion_burst":
			_on_fusion_burst(payload)
		"fusion_partial":
			_on_fusion_partial(payload)
		"fusion_single":
			_on_fusion_single(payload)
		"fusion_burst_broadcast":
			_on_fusion_burst_broadcast(payload)
		"fusion_expire":
			_on_fusion_expire()

## 觸發元素融合（個人）
func _on_fusion_start(payload: Dictionary) -> void:
	var trigger = payload.get("trigger_name", "玩家")
	var count = payload.get("fragment_count", 3)
	var duration = payload.get("duration", 25)

	# 橙色三次強閃光
	_flash_triple(Color(1.0, 0.42, 0.21, 0.5))

	# 顯示橫幅
	_banner_label.text = "🔥💧🌪️ 元素融合！集齊三種元素獲得 ×6.0！"
	_show_banner(Color(1.0, 0.42, 0.21, 1.0), duration)

## 全服廣播元素碎片位置
func _on_fusion_broadcast(payload: Dictionary) -> void:
	var trigger = payload.get("trigger_name", "玩家")
	var fragments = payload.get("fragments", [])

	# 顯示小橫幅
	_banner_label.text = "🔥💧🌪️ %s 觸發元素融合！找到帶元素的目標！" % trigger
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)
	tween.tween_interval(3.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.3)

## 收集元素碎片（個人）
func _on_fragment_collect(payload: Dictionary) -> void:
	var element = payload.get("element", "fire")
	var count = payload.get("collected_count", 1)
	var reward = payload.get("fragment_reward", 0)
	var elem_icon = ELEMENT_ICONS.get(element, "✨")
	var elem_name = ELEMENT_NAMES.get(element, "元素")
	var elem_color = ELEMENT_COLORS.get(element, Color.WHITE)

	# 元素顏色閃光
	_flash_overlay.color = Color(elem_color.r, elem_color.g, elem_color.b, 0.35)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, 0.25)

	# 浮動文字
	_collect_label.text = "%s 收集 %s！%d/3 ×0.5 +%d" % [elem_icon, elem_name, count, reward]
	_collect_label.add_theme_color_override("font_color", elem_color)
	_collect_label.modulate.a = 1.0
	_collect_timer = 1.8

## 全融合爆發（個人）
func _on_fusion_burst(payload: Dictionary) -> void:
	var mult = payload.get("mult", 6.0)
	var reward = payload.get("reward", 0)

	# 全螢幕三次強閃光（橙紅）
	_flash_triple(Color(1.0, 0.42, 0.21, 0.7))

	# 大字顯示
	_show_big_text("🔥💧🌪️ 元素全融合！×%.0f" % mult, Color(1.0, 0.85, 0.0))

	# 結算彈窗
	_show_result_popup("🔥💧🌪️ 元素全融合！", mult, reward, Color(1.0, 0.42, 0.21))

## 部分融合（個人）
func _on_fusion_partial(payload: Dictionary) -> void:
	var mult = payload.get("mult", 2.5)
	var reward = payload.get("reward", 0)
	var count = payload.get("collected_count", 2)

	_flash_triple(Color(0.8, 0.6, 0.0, 0.5))
	_show_big_text("⚡ 部分融合！%d/3 ×%.1f" % [count, mult], Color(1.0, 0.8, 0.0))
	_show_result_popup("⚡ 部分融合！", mult, reward, Color(1.0, 0.8, 0.0))

## 元素殘留（個人）
func _on_fusion_single(payload: Dictionary) -> void:
	var mult = payload.get("mult", 1.3)
	var reward = payload.get("reward", 0)

	_flash_triple(Color(0.5, 0.5, 0.5, 0.4))
	_show_big_text("✨ 元素殘留！×%.1f" % mult, Color(0.8, 0.8, 0.8))
	_show_result_popup("✨ 元素殘留！", mult, reward, Color(0.8, 0.8, 0.8))

## 全融合全服廣播
func _on_fusion_burst_broadcast(payload: Dictionary) -> void:
	var player_name = payload.get("player_name", "玩家")
	var mult = payload.get("mult", 6.0)
	var reward = payload.get("reward", 0)

	_banner_label.text = "🔥💧🌪️ %s 集齊三種元素！×%.0f 大獎 +%d！" % [player_name, mult, reward]
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)
	tween.tween_interval(4.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.3)

## 元素碎片消失
func _on_fusion_expire() -> void:
	_banner_label.text = "🔥💧🌪️ 元素碎片消失了..."
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 0.6, 0.2)
	tween.tween_interval(2.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.5)

## 工具函數：三次強閃光
func _flash_triple(color: Color) -> void:
	var tween = create_tween()
	for i in range(3):
		tween.tween_property(_flash_overlay, "color:a", color.a, 0.06)
		tween.tween_property(_flash_overlay, "color:a", 0.0, 0.1)
	_flash_overlay.color = Color(color.r, color.g, color.b, 0)

## 工具函數：顯示橫幅（帶計時條）
func _show_banner(border_color: Color, duration: float) -> void:
	var style = _banner.get_theme_stylebox("panel") as StyleBoxFlat
	if style:
		style.border_color = border_color
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)
	tween.tween_interval(duration)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.5)

## 工具函數：顯示大字
func _show_big_text(text: String, color: Color) -> void:
	var label = Label.new()
	label.text = text
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.add_theme_font_size_override("font_size", 32)
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

## 工具函數：顯示結算彈窗
func _show_result_popup(title: String, mult: float, reward: int, color: Color) -> void:
	var popup = PanelContainer.new()
	popup.custom_minimum_size = Vector2(280, 120)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.05, 0.1, 0.92)
	style.border_color = color
	style.set_border_width_all(2)
	style.set_corner_radius_all(8)
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	popup.add_child(vbox)

	var title_label = Label.new()
	title_label.text = title
	title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_label.add_theme_font_size_override("font_size", 18)
	title_label.add_theme_color_override("font_color", color)
	vbox.add_child(title_label)

	var mult_label = Label.new()
	mult_label.text = "倍率：×%.1f" % mult
	mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mult_label.add_theme_font_size_override("font_size", 16)
	mult_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	vbox.add_child(mult_label)

	var reward_label = Label.new()
	reward_label.text = "獎勵：+%d 🪙" % reward
	reward_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_label.add_theme_font_size_override("font_size", 16)
	reward_label.add_theme_color_override("font_color", Color(0.0, 1.0, 0.53))
	vbox.add_child(reward_label)

	# 右側滑入
	var vp_size = get_viewport().size
	popup.position = Vector2(vp_size.x, vp_size.y / 2 - 60)
	popup.modulate.a = 0.0
	add_child(popup)

	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 300, 0.3)
	tween.parallel().tween_property(popup, "modulate:a", 1.0, 0.3)
	tween.tween_interval(3.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.4)
	tween.tween_callback(popup.queue_free)
