## LuckySynergyBurstPanel.gd — 幸運共鳴爆發魚系統面板（DAY-239）
## 業界原創「多效疊加共鳴爆發」機制
##
## 視覺設計：
##   - 粉紅共鳴主題（#FF6B9D + #C39BD3 + #F8C8DC + #FFF0F8）
##   - synergy_full：粉紅三次強閃光 + 頂部橫幅 + 「✨ 共鳴爆發！」大字 + 效果列表 + 計時條
##   - synergy_small：紫色閃光 + 「✨ 小型共鳴！」大字 + 效果名稱
##   - synergy_base：橙色閃光 + 「✨ 基礎爆發！」大字 + HP -30% 提示
##   - synergy_end：粉紅淡出
extends CanvasLayer

# 計時條節點
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null

# 主題顏色
const COLOR_FULL     = Color("#FF6B9D")  # 粉紅（共鳴爆發）
const COLOR_SMALL    = Color("#C39BD3")  # 紫色（小型共鳴）
const COLOR_BASE     = Color("#F39C12")  # 橙色（基礎爆發）
const COLOR_PALE     = Color("#F8C8DC")  # 淡粉
const COLOR_GOLD     = Color("#FFD700")  # 金黃
const COLOR_WHITE    = Color("#FFFFFF")  # 白色

func _ready() -> void:
	layer = 6  # 幸運共鳴爆發魚面板層級

## 處理幸運共鳴爆發魚訊息
func handle_lucky_synergy_burst(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"synergy_full":
			_on_synergy_full(payload)
		"synergy_small":
			_on_synergy_small(payload)
		"synergy_base":
			_on_synergy_base(payload)
		"synergy_end":
			_on_synergy_end(payload)

## synergy_full — 共鳴爆發（≥2 個效果）
func _on_synergy_full(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var effect_count: int = payload.get("effect_count", 0)
	var effect_names: Array = payload.get("effect_names", [])
	var extra_mult: float = payload.get("extra_mult", 1.5)
	var duration_sec: int = payload.get("duration_sec", 6)

	var vp_size = get_viewport().size

	# 粉紅三次強閃光
	_flash_screen(COLOR_FULL, 0.14)
	await get_tree().create_timer(0.09).timeout
	_flash_screen(COLOR_WHITE, 0.11)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(COLOR_PALE, 0.10)

	# 頂部橫幅
	var banner = Label.new()
	banner.text = "✨ %s 觸發共鳴爆發！%d 個效果疊加！" % [player_name, effect_count]
	banner.add_theme_font_size_override("font_size", 14)
	banner.add_theme_color_override("font_color", COLOR_PALE)
	banner.position = Vector2(vp_size.x / 2 - 150, 6)
	add_child(banner)

	var tween_banner = banner.create_tween()
	tween_banner.tween_interval(float(duration_sec) - 0.5)
	tween_banner.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween_banner.tween_callback(banner.queue_free)

	# 「✨ 共鳴爆發！」大字
	var big_label = Label.new()
	big_label.text = "✨ 共鳴爆發！"
	big_label.add_theme_font_size_override("font_size", 52)
	big_label.add_theme_color_override("font_color", COLOR_FULL)
	big_label.position = vp_size / 2 - Vector2(100, 32)
	add_child(big_label)

	var tween_big = big_label.create_tween()
	tween_big.tween_property(big_label, "scale", Vector2(1.25, 1.25), 0.12)
	tween_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.10)
	tween_big.tween_interval(0.5)
	tween_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_big.tween_callback(big_label.queue_free)

	# 效果列表（顯示疊加的效果名稱）
	var names_str = " + ".join(effect_names)
	var effects_label = Label.new()
	effects_label.text = "【%s】×%.1f" % [names_str, extra_mult]
	effects_label.add_theme_font_size_override("font_size", 13)
	effects_label.add_theme_color_override("font_color", COLOR_GOLD)
	effects_label.position = Vector2(vp_size.x / 2 - 100, vp_size.y / 2 + 30)
	add_child(effects_label)

	var tween_effects = effects_label.create_tween()
	tween_effects.tween_interval(2.5)
	tween_effects.tween_property(effects_label, "modulate:a", 0.0, 0.5)
	tween_effects.tween_callback(effects_label.queue_free)

	# 底部計時條（粉紅→深粉漸變）
	_spawn_timer_bar(float(duration_sec), COLOR_FULL)

## synergy_small — 小型共鳴（1 個效果）
func _on_synergy_small(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var effect_names: Array = payload.get("effect_names", [])
	var extra_mult: float = payload.get("extra_mult", 1.3)
	var duration_sec: int = payload.get("duration_sec", 4)

	var vp_size = get_viewport().size

	# 紫色閃光
	_flash_screen(COLOR_SMALL, 0.12)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(COLOR_WHITE, 0.09)

	# 「✨ 小型共鳴！」大字
	var small_label = Label.new()
	small_label.text = "✨ 小型共鳴！"
	small_label.add_theme_font_size_override("font_size", 40)
	small_label.add_theme_color_override("font_color", COLOR_SMALL)
	small_label.position = vp_size / 2 - Vector2(80, 24)
	add_child(small_label)

	var tween_small = small_label.create_tween()
	tween_small.tween_property(small_label, "scale", Vector2(1.15, 1.15), 0.10)
	tween_small.tween_property(small_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_small.tween_interval(0.4)
	tween_small.tween_property(small_label, "modulate:a", 0.0, 0.4)
	tween_small.tween_callback(small_label.queue_free)

	# 效果名稱
	var effect_name = effect_names[0] if effect_names.size() > 0 else ""
	var info_label = Label.new()
	info_label.text = "%s 效果 ×%.1f（%s 觸發）" % [effect_name, extra_mult, player_name]
	info_label.add_theme_font_size_override("font_size", 12)
	info_label.add_theme_color_override("font_color", COLOR_GOLD)
	info_label.position = Vector2(vp_size.x / 2 - 80, vp_size.y / 2 + 26)
	add_child(info_label)

	var tween_info = info_label.create_tween()
	tween_info.tween_interval(1.5)
	tween_info.tween_property(info_label, "modulate:a", 0.0, 0.4)
	tween_info.tween_callback(info_label.queue_free)

	# 底部計時條（紫色）
	_spawn_timer_bar(float(duration_sec), COLOR_SMALL)

## synergy_base — 基礎爆發（0 個效果）
func _on_synergy_base(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var extra_mult: float = payload.get("extra_mult", 1.8)
	var duration_sec: int = payload.get("duration_sec", 5)
	var damaged_count: int = payload.get("damaged_count", 0)

	var vp_size = get_viewport().size

	# 橙色閃光
	_flash_screen(COLOR_BASE, 0.13)
	await get_tree().create_timer(0.09).timeout
	_flash_screen(COLOR_WHITE, 0.10)

	# 「✨ 基礎爆發！」大字
	var base_label = Label.new()
	base_label.text = "✨ 基礎爆發！"
	base_label.add_theme_font_size_override("font_size", 40)
	base_label.add_theme_color_override("font_color", COLOR_BASE)
	base_label.position = vp_size / 2 - Vector2(80, 24)
	add_child(base_label)

	var tween_base = base_label.create_tween()
	tween_base.tween_property(base_label, "scale", Vector2(1.15, 1.15), 0.10)
	tween_base.tween_property(base_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_base.tween_interval(0.4)
	tween_base.tween_property(base_label, "modulate:a", 0.0, 0.4)
	tween_base.tween_callback(base_label.queue_free)

	# HP -30% 提示
	if damaged_count > 0:
		var hp_label = Label.new()
		hp_label.text = "💥 %d 個目標 HP -30%%！個人 ×%.1f（%s）" % [damaged_count, extra_mult, player_name]
		hp_label.add_theme_font_size_override("font_size", 12)
		hp_label.add_theme_color_override("font_color", COLOR_GOLD)
		hp_label.position = Vector2(vp_size.x / 2 - 100, vp_size.y / 2 + 26)
		add_child(hp_label)

		var tween_hp = hp_label.create_tween()
		tween_hp.tween_property(hp_label, "position:y", hp_label.position.y - 18, 0.5)
		tween_hp.parallel().tween_property(hp_label, "modulate:a", 0.0, 0.5)
		tween_hp.tween_callback(hp_label.queue_free)

	# 底部計時條（橙色）
	_spawn_timer_bar(float(duration_sec), COLOR_BASE)

## synergy_end — 共鳴結束
func _on_synergy_end(_payload: Dictionary) -> void:
	# 清除計時條
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		_timer_bar_bg.queue_free()
		_timer_bar_bg = null

## 建立底部計時條
func _spawn_timer_bar(duration: float, color: Color) -> void:
	var vp_size = get_viewport().size

	# 清除舊的計時條
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
	if is_instance_valid(_timer_bar_bg):
		_timer_bar_bg.queue_free()

	var bg = ColorRect.new()
	bg.color = Color(0.1, 0.1, 0.1, 0.6)
	bg.size = Vector2(vp_size.x, 8)
	bg.position = Vector2(0, vp_size.y - 16)  # 稍高一點，避免與其他計時條重疊
	add_child(bg)
	_timer_bar_bg = bg

	var bar = ColorRect.new()
	bar.color = color
	bar.size = Vector2(vp_size.x, 8)
	bar.position = Vector2(0, vp_size.y - 16)
	add_child(bar)
	_timer_bar = bar

	var tween_bar = bar.create_tween()
	tween_bar.tween_property(bar, "size:x", 0.0, duration)
	tween_bar.tween_callback(func():
		if is_instance_valid(bar):
			bar.queue_free()
		if is_instance_valid(bg):
			bg.queue_free()
	)

# ---- 輔助函數 ----

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.28)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)
