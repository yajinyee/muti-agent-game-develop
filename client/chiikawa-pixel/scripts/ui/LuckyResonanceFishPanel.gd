## LuckyResonanceFishPanel.gd — 幸運共鳴魚系統面板（DAY-222）
## 業界原創「全服共鳴」機制
##
## 視覺設計：
##   - 天藍共鳴主題（#00BFFF + #00FFFF + #87CEEB + #E0F8FF）
##   - resonance_start：天藍三次強閃光 + 頂部橫幅 + 底部共鳴進度條
##   - resonance_progress：進度條更新 + 每 5 點閃光 + 「還差 N 槍」提示
##   - resonance_burst：全螢幕三次天藍強閃光 + 「🎵 共鳴爆發！」52px大字 + 倍率計時條
##   - resonance_small_burst：青色閃光 + 「🎵 小型共鳴！」40px大字
##   - resonance_result：右側滑入結算彈窗
##   - resonance_boost_end：計時條淡出
extends CanvasLayer

# 共鳴狀態
var _resonance_active: bool = false
var _resonance_target: int = 30
var _resonance_count: int = 0
var _banner: Control = null
var _progress_bar: Control = null
var _boost_bar: Control = null

# 主題顏色
const COLOR_PRIMARY   = Color("#00BFFF")  # 天藍
const COLOR_CYAN      = Color("#00FFFF")  # 青色
const COLOR_LIGHT     = Color("#87CEEB")  # 淡藍
const COLOR_PALE      = Color("#E0F8FF")  # 極淡藍
const COLOR_BG        = Color(0.0, 0.05, 0.12, 0.88)

func _ready() -> void:
	layer = 23  # 幸運共鳴魚面板層級

## 處理幸運共鳴魚訊息
func handle_lucky_resonance_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"resonance_start":
			_on_resonance_start(payload)
		"resonance_progress":
			_on_resonance_progress(payload)
		"resonance_burst":
			_on_resonance_burst(payload)
		"resonance_small_burst":
			_on_resonance_small_burst(payload)
		"resonance_result":
			_on_resonance_result(payload)
		"resonance_boost_end":
			_on_resonance_boost_end(payload)

## resonance_start — 共鳴模式開始
func _on_resonance_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var duration_sec: int = payload.get("duration_sec", 15)
	var target: int = payload.get("target", 30)

	_resonance_active = true
	_resonance_target = target
	_resonance_count = 0

	# 天藍三次強閃光
	_flash_screen(COLOR_PRIMARY, 0.2)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_CYAN, 0.18)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_LIGHT, 0.15)

	# 頂部橫幅 + 進度條
	_show_banner("🎵 共鳴模式！", "%s 觸發共鳴魚！全服合力射擊 %d 次觸發共鳴爆發！" % [player_name, target], duration_sec)

## resonance_progress — 共鳴能量進度
func _on_resonance_progress(payload: Dictionary) -> void:
	var count: int = payload.get("count", 0)
	var target: int = payload.get("target", 30)

	_resonance_count = count

	# 更新進度條
	if _progress_bar != null and is_instance_valid(_progress_bar):
		var pct: float = float(count) / float(target)
		var tween = _progress_bar.create_tween()
		tween.tween_property(_progress_bar, "size:x", get_viewport().size.x * pct, 0.2)
		# 顏色漸變（天藍→青→白）
		var bar_color = COLOR_PRIMARY.lerp(COLOR_CYAN, pct)
		tween.parallel().tween_property(_progress_bar, "color", bar_color, 0.2)

	# 每 5 點閃光
	_flash_screen(COLOR_PRIMARY, 0.1)

	# 「還差 N 槍」提示
	var remaining: int = target - count
	if remaining > 0:
		var vp_size = get_viewport().size
		var hint = Label.new()
		hint.text = "🎵 還差 %d 槍！" % remaining
		hint.add_theme_font_size_override("font_size", 14)
		hint.add_theme_color_override("font_color", COLOR_CYAN)
		hint.position = Vector2(vp_size.x / 2 - 60, vp_size.y - 80)
		add_child(hint)

		var tween = hint.create_tween()
		tween.tween_property(hint, "modulate:a", 0.0, 0.8)
		tween.tween_callback(hint.queue_free)

## resonance_burst — 共鳴爆發（完整）
func _on_resonance_burst(payload: Dictionary) -> void:
	_resonance_active = false
	_hide_banner()

	var total_shots: int = payload.get("total_shots", 0)
	var boost_mult: float = payload.get("boost_mult", 1.8)
	var boost_sec: int = payload.get("boost_sec", 6)

	# 全螢幕三次天藍強閃光
	_flash_screen(COLOR_PRIMARY, 0.15)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color.WHITE, 0.12)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(COLOR_CYAN, 0.15)

	# 「🎵 共鳴爆發！」大字
	var big_label = Label.new()
	big_label.text = "🎵 共鳴爆發！"
	big_label.add_theme_font_size_override("font_size", 52)
	big_label.add_theme_color_override("font_color", COLOR_CYAN)
	big_label.set_anchors_preset(Control.PRESET_CENTER)
	var vp_size = get_viewport().size
	big_label.position = vp_size / 2 - Vector2(150, 30)
	add_child(big_label)

	var tween_label = big_label.create_tween()
	tween_label.tween_property(big_label, "scale", Vector2(1.3, 1.3), 0.15)
	tween_label.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.1)
	tween_label.tween_interval(0.5)
	tween_label.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_label.tween_callback(big_label.queue_free)

	# 副標題（全服合力 N 槍）
	var sub_label = Label.new()
	sub_label.text = "全服合力 %d 槍！×%.1f 倍率加成 %d 秒！" % [total_shots, boost_mult, boost_sec]
	sub_label.add_theme_font_size_override("font_size", 18)
	sub_label.add_theme_color_override("font_color", COLOR_LIGHT)
	sub_label.position = vp_size / 2 - Vector2(140, -10)
	add_child(sub_label)

	var tween_sub = sub_label.create_tween()
	tween_sub.tween_interval(0.8)
	tween_sub.tween_property(sub_label, "modulate:a", 0.0, 0.5)
	tween_sub.tween_callback(sub_label.queue_free)

	# 倍率計時條（底部）
	_show_boost_bar(boost_mult, boost_sec)

## resonance_small_burst — 小型共鳴
func _on_resonance_small_burst(payload: Dictionary) -> void:
	_resonance_active = false
	_hide_banner()

	var total_shots: int = payload.get("total_shots", 0)
	var boost_mult: float = payload.get("boost_mult", 1.3)
	var boost_sec: int = payload.get("boost_sec", 3)

	# 青色閃光
	_flash_screen(COLOR_LIGHT, 0.2)

	# 「🎵 小型共鳴！」大字
	var big_label = Label.new()
	big_label.text = "🎵 小型共鳴！"
	big_label.add_theme_font_size_override("font_size", 40)
	big_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	big_label.set_anchors_preset(Control.PRESET_CENTER)
	var vp_size = get_viewport().size
	big_label.position = vp_size / 2 - Vector2(120, 25)
	add_child(big_label)

	var tween_label = big_label.create_tween()
	tween_label.tween_interval(0.6)
	tween_label.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_label.tween_callback(big_label.queue_free)

	# 倍率計時條
	_show_boost_bar(boost_mult, boost_sec)

## resonance_result — 爆發結算
func _on_resonance_result(payload: Dictionary) -> void:
	var affected_count: int = payload.get("affected_count", 0)
	var reward_pool: int = payload.get("reward_pool", 0)
	var total_shots: int = payload.get("total_shots", 0)

	if affected_count <= 0 and reward_pool <= 0:
		return

	# 右側滑入結算彈窗
	_show_result_popup(affected_count, reward_pool, total_shots)

## resonance_boost_end — 倍率加成結束
func _on_resonance_boost_end(_payload: Dictionary) -> void:
	if _boost_bar != null and is_instance_valid(_boost_bar):
		var tween = _boost_bar.create_tween()
		tween.tween_property(_boost_bar, "modulate:a", 0.0, 0.4)
		tween.tween_callback(_boost_bar.queue_free)
		_boost_bar = null

# ---- 輔助函數 ----

## 顯示頂部橫幅 + 進度條
func _show_banner(title: String, subtitle: String, duration_sec: int) -> void:
	_hide_banner()

	var banner = Control.new()
	banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	banner.position = Vector2(0, 8)
	banner.size = Vector2(get_viewport().size.x, 60)
	add_child(banner)

	var bg = ColorRect.new()
	bg.color = COLOR_BG
	bg.size = banner.size
	banner.add_child(bg)

	var title_label = Label.new()
	title_label.text = title
	title_label.add_theme_font_size_override("font_size", 20)
	title_label.add_theme_color_override("font_color", COLOR_CYAN)
	title_label.position = Vector2(12, 4)
	banner.add_child(title_label)

	var sub_label = Label.new()
	sub_label.text = subtitle
	sub_label.add_theme_font_size_override("font_size", 12)
	sub_label.add_theme_color_override("font_color", COLOR_PALE)
	sub_label.position = Vector2(12, 28)
	banner.add_child(sub_label)

	# 共鳴進度條（底部，天藍→青漸變）
	var progress_bg = ColorRect.new()
	progress_bg.color = Color(0.1, 0.1, 0.2, 0.8)
	progress_bg.position = Vector2(0, 52)
	progress_bg.size = Vector2(get_viewport().size.x, 8)
	banner.add_child(progress_bg)

	var progress_bar = ColorRect.new()
	progress_bar.name = "ProgressBar"
	progress_bar.color = COLOR_PRIMARY
	progress_bar.position = Vector2(0, 52)
	progress_bar.size = Vector2(0, 8)
	banner.add_child(progress_bar)
	_progress_bar = progress_bar

	# 計時條（最底部，細線）
	var timer_bar = ColorRect.new()
	timer_bar.color = Color(COLOR_LIGHT.r, COLOR_LIGHT.g, COLOR_LIGHT.b, 0.5)
	timer_bar.position = Vector2(0, 56)
	timer_bar.size = Vector2(get_viewport().size.x, 4)
	banner.add_child(timer_bar)

	var tween = banner.create_tween()
	tween.tween_property(timer_bar, "size:x", 0.0, float(duration_sec))

	_banner = banner

## 隱藏橫幅
func _hide_banner() -> void:
	if _banner != null and is_instance_valid(_banner):
		_banner.queue_free()
	_banner = null
	_progress_bar = null

## 顯示倍率計時條（底部）
func _show_boost_bar(boost_mult: float, boost_sec: int) -> void:
	if _boost_bar != null and is_instance_valid(_boost_bar):
		_boost_bar.queue_free()

	var vp_size = get_viewport().size
	var bar_container = Control.new()
	bar_container.position = Vector2(0, vp_size.y - 20)
	bar_container.size = Vector2(vp_size.x, 20)
	add_child(bar_container)
	_boost_bar = bar_container

	var bar_bg = ColorRect.new()
	bar_bg.color = Color(0.0, 0.1, 0.2, 0.8)
	bar_bg.size = bar_container.size
	bar_container.add_child(bar_bg)

	var bar = ColorRect.new()
	bar.color = COLOR_CYAN
	bar.size = Vector2(vp_size.x, 20)
	bar_container.add_child(bar)

	var label = Label.new()
	label.text = "🎵 ×%.1f 共鳴加成" % boost_mult
	label.add_theme_font_size_override("font_size", 13)
	label.add_theme_color_override("font_color", Color.BLACK)
	label.position = Vector2(vp_size.x / 2 - 60, 2)
	bar_container.add_child(label)

	var tween = bar.create_tween()
	tween.tween_property(bar, "size:x", 0.0, float(boost_sec))
	tween.parallel().tween_property(bar, "color", COLOR_PRIMARY, float(boost_sec))

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.4)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)

## 結算彈窗（右側滑入）
func _show_result_popup(affected_count: int, reward_pool: int, total_shots: int) -> void:
	var vp_size = get_viewport().size
	var popup = Control.new()
	popup.size = Vector2(230, 100)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2 - 50)
	add_child(popup)

	var bg = ColorRect.new()
	bg.color = COLOR_BG
	bg.size = popup.size
	popup.add_child(bg)

	var border = ColorRect.new()
	border.color = COLOR_CYAN
	border.size = Vector2(popup.size.x, 3)
	popup.add_child(border)

	var title_label = Label.new()
	title_label.text = "🎵 共鳴爆發結算"
	title_label.add_theme_font_size_override("font_size", 14)
	title_label.add_theme_color_override("font_color", COLOR_CYAN)
	title_label.position = Vector2(8, 8)
	popup.add_child(title_label)

	var shots_label = Label.new()
	shots_label.text = "全服合力：%d 槍" % total_shots
	shots_label.add_theme_font_size_override("font_size", 12)
	shots_label.add_theme_color_override("font_color", COLOR_LIGHT)
	shots_label.position = Vector2(8, 32)
	popup.add_child(shots_label)

	var affected_label = Label.new()
	affected_label.text = "HP 削減：%d 個目標" % affected_count
	affected_label.add_theme_font_size_override("font_size", 12)
	affected_label.add_theme_color_override("font_color", COLOR_PALE)
	affected_label.position = Vector2(8, 54)
	popup.add_child(affected_label)

	var reward_label = Label.new()
	reward_label.text = "獎勵池：%d 金幣（按貢獻分配）" % reward_pool
	reward_label.add_theme_font_size_override("font_size", 11)
	reward_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	reward_label.position = Vector2(8, 76)
	popup.add_child(reward_label)

	# 右側滑入動畫
	var tween = popup.create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 240.0, 0.3)
	tween.tween_interval(3.0)
	tween.tween_property(popup, "position:x", vp_size.x + 10.0, 0.3)
	tween.tween_callback(popup.queue_free)
