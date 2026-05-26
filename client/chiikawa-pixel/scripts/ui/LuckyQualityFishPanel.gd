## LuckyQualityFishPanel.gd — T140 幸運品質魚 UI
## lucky-panel-agent 負責維護
## 業界依據：Fishing Frenzy Chapter 3「Fish Quality tier system that raises the stakes on every cast」
## 視覺主題：彩虹漸層 + 品質等級顯示 + 傳說品質演出
extends CanvasLayer

const LAYER_Z = 40

const COLOR_COMMON    = Color(0.7, 0.7, 0.7)   # 灰色（Common）
const COLOR_RARE      = Color(0.0, 0.5, 1.0)   # 藍色（Rare）
const COLOR_EPIC      = Color(0.6, 0.0, 1.0)   # 紫色（Epic）
const COLOR_LEGENDARY = Color(1.0, 0.85, 0.0)  # 金色（Legendary）
const COLOR_BG        = Color(0.03, 0.03, 0.05, 0.90)

var _banner: Control = null
var _flash_overlay: ColorRect = null

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	GameManager.lucky_quality_fish.connect(_on_lucky_quality_fish)

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.85, 0.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

func _get_tier_color(tier_name: String) -> Color:
	match tier_name:
		"Common": return COLOR_COMMON
		"Rare": return COLOR_RARE
		"Epic": return COLOR_EPIC
		"Legendary": return COLOR_LEGENDARY
	return COLOR_COMMON

func _flash(color: Color, alpha: float = 0.4, duration: float = 0.2) -> void:
	if not is_instance_valid(_flash_overlay):
		return
	_flash_overlay.color = Color(color.r, color.g, color.b, alpha)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)

func _show_banner(text: String, color: Color) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()
	_banner = Control.new()
	_banner.position = Vector2(0, 100)
	_banner.size = Vector2(1280, 70)
	add_child(_banner)
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0, 0, 0, 0.85)
	_banner.add_child(bg)
	var lbl = Label.new()
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 26)
	lbl.modulate = color
	_banner.add_child(lbl)
	var tween = _banner.create_tween()
	tween.tween_interval(2.5)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_banner):
			_banner.queue_free()
			_banner = null
	)

func _on_lucky_quality_fish(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"quality_result":
			var tier_name = data.get("tier_name", "Common")
			var tier_mult = data.get("tier_mult", 2.0)
			var reward = data.get("reward", 0)
			var color = _get_tier_color(tier_name)
			_flash(color, 0.4 + (["Common","Rare","Epic","Legendary"].find(tier_name)) * 0.1, 0.2)
			_show_banner("✨ 品質鑑定！%s 品質！×%.0f！獎勵 %d！" % [tier_name, tier_mult, reward], color)

		"legendary_boost":
			_flash(COLOR_LEGENDARY, 0.8, 0.35)
			_flash(COLOR_LEGENDARY, 0.8, 0.35)
			_flash(COLOR_LEGENDARY, 0.8, 0.35)
			_show_banner("✨ 傳說品質！LEGENDARY！全服 ×5.0 加成 12 秒！", COLOR_LEGENDARY)

		"legendary_end":
			pass
