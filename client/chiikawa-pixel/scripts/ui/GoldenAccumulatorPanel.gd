## GoldenAccumulatorPanel.gd — 黃金累積魚系統面板（DAY-214）
## 業界依據：Evolution Ice Fishing Live 2026「random multipliers 2x-10x」
## 業界原創「全服累積爆發」機制
##
## 視覺設計：
##   - 黃金主題（#FFD700 + #FF8C00 + #FFF8DC + #FF4500）
##   - accum_appear：金色閃光 + 頂部橫幅 + 底部累積進度條
##   - accum_progress：進度條更新 + 每 5 點閃光 + 「還差 N 個」提示
##   - early_detonate：金色強閃光 + 「提前引爆！」大字
##   - burst_start/early_burst_start：全螢幕三次金色強閃光 + 「🌟 黃金爆發！」大字 + 倍率計時條
##   - burst_end：計時條淡出
##   - accum_escape：進度條淡出
extends CanvasLayer

var _progress_bar: Control = null   # 底部累積進度條
var _boost_bar: Control = null      # 底部倍率計時條
var _accum_target: int = 20

func _ready() -> void:
	layer = 31  # 黃金累積魚面板層級

## 處理黃金累積魚訊息
func handle_golden_accumulator(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"accum_appear":
			_on_accum_appear(payload)
		"accum_progress":
			_on_accum_progress(payload)
		"early_detonate":
			_on_early_detonate(payload)
		"burst_start", "early_burst_start":
			_on_burst_start(payload)
		"burst_end":
			_on_burst_end()
		"accum_escape":
			_on_accum_escape()

## 黃金累積魚出現 — 金色閃光 + 頂部橫幅 + 底部累積進度條
func _on_accum_appear(payload: Dictionary) -> void:
	_accum_target = payload.get("accum_target", 20)

	# 金色雙閃光
	_double_flash(Color("#FFD700"), 0.55)

	# 頂部橫幅
	var banner = _make_banner(
		"🌟 黃金累積魚出現！全服合力擊破 %d 個目標觸發黃金爆發！" % _accum_target,
		Color(0.1, 0.07, 0.0, 0.88),
		Color("#FFD700")
	)
	add_child(banner)
	var tween_b = create_tween()
	tween_b.tween_property(banner, "position:y", 0.0, 0.25)
	tween_b.tween_interval(4.0)
	tween_b.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween_b.tween_callback(banner.queue_free)

	# 底部累積進度條
	_show_progress_bar(0, _accum_target)

## 累積進度更新 — 進度條更新 + 閃光 + 提示
func _on_accum_progress(payload: Dictionary) -> void:
	var count: int = payload.get("accum_count", 0)
	var target: int = payload.get("accum_target", 20)

	# 更新進度條
	_update_progress_bar(count, target)

	# 每 5 點額外閃光
	if count % 5 == 0:
		_single_flash(Color("#FFD700"), 0.35)

	# 「還差 N 個」浮動提示
	var remaining = target - count
	if remaining > 0:
		var label = Label.new()
		label.text = "還差 %d 個！" % remaining
		label.add_theme_font_size_override("font_size", 20)
		label.add_theme_color_override("font_color", Color("#FFD700"))
		label.position = Vector2(640 - 80, 680)
		label.size = Vector2(160, 30)
		label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		add_child(label)
		var tween = create_tween()
		tween.tween_property(label, "position:y", 650.0, 0.5)
		tween.parallel().tween_property(label, "modulate:a", 0.0, 0.5)
		tween.tween_callback(label.queue_free)

## 提前引爆 — 金色強閃光 + 「提前引爆！」大字
func _on_early_detonate(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "玩家")

	# 移除進度條
	_remove_progress_bar()

	# 金色強閃光
	_triple_flash(Color("#FFD700"), 0.75)

	# 「提前引爆！」大字
	var big_label = Label.new()
	big_label.text = "💥 %s 提前引爆！" % player_name
	big_label.add_theme_font_size_override("font_size", 44)
	big_label.add_theme_color_override("font_color", Color("#FF8C00"))
	big_label.position = Vector2(640 - 220, 280)
	big_label.size = Vector2(440, 60)
	big_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	add_child(big_label)
	var tween_big = create_tween()
	tween_big.tween_property(big_label, "scale", Vector2(1.15, 1.15), 0.1)
	tween_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.1)
	tween_big.tween_interval(1.0)
	tween_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_big.tween_callback(big_label.queue_free)

## 黃金爆發開始 — 全螢幕三次金色強閃光 + 大字 + 倍率計時條
func _on_burst_start(payload: Dictionary) -> void:
	var boost_mult: float = payload.get("boost_mult", 2.0)
	var boost_sec: int = payload.get("boost_sec", 8)
	var affected: int = payload.get("affected_count", 0)
	var is_early: bool = payload.get("event", "") == "early_burst_start"

	# 移除進度條
	_remove_progress_bar()

	# 全螢幕三次金色強閃光
	_triple_flash(Color("#FFD700"), 0.85)

	# 「🌟 黃金爆發！」大字
	var text = "🌟 黃金爆發！"
	if is_early:
		text = "🌟 提前黃金爆發！"
	var big_label = Label.new()
	big_label.text = text
	big_label.add_theme_font_size_override("font_size", 52)
	big_label.add_theme_color_override("font_color", Color("#FFD700"))
	big_label.position = Vector2(640 - 200, 260)
	big_label.size = Vector2(400, 70)
	big_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	add_child(big_label)
	var tween_big = create_tween()
	tween_big.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.12)
	tween_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.12)
	tween_big.tween_interval(1.5)
	tween_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_big.tween_callback(big_label.queue_free)

	# 副標題（HP -60% + 倍率加成）
	var sub_label = Label.new()
	sub_label.text = "全場 %d 個目標 HP -60%%！×%.0f 倍率加成 %d 秒！" % [affected, boost_mult, boost_sec]
	sub_label.add_theme_font_size_override("font_size", 20)
	sub_label.add_theme_color_override("font_color", Color("#FFF8DC"))
	sub_label.position = Vector2(640 - 240, 330)
	sub_label.size = Vector2(480, 30)
	sub_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	add_child(sub_label)
	var tween_sub = create_tween()
	tween_sub.tween_interval(2.0)
	tween_sub.tween_property(sub_label, "modulate:a", 0.0, 0.4)
	tween_sub.tween_callback(sub_label.queue_free)

	# 底部倍率計時條
	_show_boost_bar(boost_sec)

## 黃金爆發結束 — 計時條淡出
func _on_burst_end() -> void:
	if is_instance_valid(_boost_bar):
		var tween = create_tween()
		tween.tween_property(_boost_bar, "modulate:a", 0.0, 0.3)
		tween.tween_callback(_boost_bar.queue_free)
		_boost_bar = null

## 累積魚逃跑 — 進度條淡出
func _on_accum_escape() -> void:
	_remove_progress_bar()

	var label = Label.new()
	label.text = "🌟 黃金累積魚逃跑了..."
	label.add_theme_font_size_override("font_size", 18)
	label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 0.8))
	label.position = Vector2(640 - 160, 360)
	label.size = Vector2(320, 30)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	add_child(label)
	var tween = create_tween()
	tween.tween_interval(1.5)
	tween.tween_property(label, "modulate:a", 0.0, 0.5)
	tween.tween_callback(label.queue_free)

## 顯示底部累積進度條
func _show_progress_bar(current: int, target: int) -> void:
	if is_instance_valid(_progress_bar):
		_progress_bar.queue_free()

	var bar_container = Control.new()
	bar_container.position = Vector2(0, 700)
	bar_container.size = Vector2(1280, 20)
	_progress_bar = bar_container
	add_child(bar_container)

	# 背景
	var bg = ColorRect.new()
	bg.color = Color(0.1, 0.07, 0.0, 0.7)
	bg.position = Vector2(0, 0)
	bg.size = Vector2(1280, 20)
	bar_container.add_child(bg)

	# 進度條（初始寬度 = current/target * 1280）
	var bar = ColorRect.new()
	bar.color = Color("#FFD700")
	bar.position = Vector2(0, 0)
	var init_width = float(current) / float(target) * 1280.0
	bar.size = Vector2(init_width, 20)
	bar.name = "ProgressFill"
	bar_container.add_child(bar)

	# 標籤
	var label = Label.new()
	label.text = "%d / %d" % [current, target]
	label.add_theme_font_size_override("font_size", 13)
	label.add_theme_color_override("font_color", Color("#FFF8DC"))
	label.position = Vector2(0, 2)
	label.size = Vector2(1280, 16)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.name = "ProgressLabel"
	bar_container.add_child(label)

## 更新累積進度條
func _update_progress_bar(current: int, target: int) -> void:
	if not is_instance_valid(_progress_bar):
		return

	var bar = _progress_bar.get_node_or_null("ProgressFill")
	var label = _progress_bar.get_node_or_null("ProgressLabel")

	if is_instance_valid(bar):
		var new_width = float(current) / float(target) * 1280.0
		var tween = create_tween()
		tween.tween_property(bar, "size:x", new_width, 0.2)
		# 顏色漸變（金→橙→紅橙）
		var t_ratio = float(current) / float(target)
		var new_color = Color("#FFD700").lerp(Color("#FF4500"), t_ratio)
		tween.parallel().tween_property(bar, "color", new_color, 0.2)

	if is_instance_valid(label):
		label.text = "%d / %d" % [current, target]

## 移除進度條
func _remove_progress_bar() -> void:
	if is_instance_valid(_progress_bar):
		var tween = create_tween()
		tween.tween_property(_progress_bar, "modulate:a", 0.0, 0.2)
		tween.tween_callback(_progress_bar.queue_free)
		_progress_bar = null

## 顯示底部倍率計時條
func _show_boost_bar(duration: int) -> void:
	if is_instance_valid(_boost_bar):
		_boost_bar.queue_free()

	var bar_container = Control.new()
	bar_container.position = Vector2(0, 700)
	bar_container.size = Vector2(1280, 20)
	_boost_bar = bar_container
	add_child(bar_container)

	var bg = ColorRect.new()
	bg.color = Color(0.1, 0.07, 0.0, 0.7)
	bg.position = Vector2(0, 0)
	bg.size = Vector2(1280, 20)
	bar_container.add_child(bg)

	var bar = ColorRect.new()
	bar.color = Color("#FFD700")
	bar.position = Vector2(0, 0)
	bar.size = Vector2(1280, 20)
	bar_container.add_child(bar)

	var label = Label.new()
	label.text = "🌟 ×2 黃金加成中..."
	label.add_theme_font_size_override("font_size", 13)
	label.add_theme_color_override("font_color", Color("#FFF8DC"))
	label.position = Vector2(0, 2)
	label.size = Vector2(1280, 16)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	bar_container.add_child(label)

	# 計時條縮短動畫（顏色漸變：金→橙→紅橙）
	var tween = create_tween()
	tween.tween_property(bar, "size:x", 0.0, float(duration))
	tween.parallel().tween_method(
		func(t: float) -> void:
			var c = Color("#FFD700").lerp(Color("#FF4500"), t)
			bar.color = c,
		0.0, 1.0, float(duration)
	)

## 三次強閃光效果
func _triple_flash(color: Color, alpha: float) -> void:
	for i in range(3):
		var flash = ColorRect.new()
		flash.color = Color(color.r, color.g, color.b, 0.0)
		flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		add_child(flash)
		var delay = i * 0.18
		var tween = create_tween()
		tween.tween_interval(delay)
		tween.tween_property(flash, "color:a", alpha, 0.08)
		tween.tween_property(flash, "color:a", 0.0, 0.14)
		tween.tween_callback(flash.queue_free)

## 雙閃光效果
func _double_flash(color: Color, alpha: float) -> void:
	for i in range(2):
		var flash = ColorRect.new()
		flash.color = Color(color.r, color.g, color.b, 0.0)
		flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		add_child(flash)
		var delay = i * 0.18
		var tween = create_tween()
		tween.tween_interval(delay)
		tween.tween_property(flash, "color:a", alpha, 0.08)
		tween.tween_property(flash, "color:a", 0.0, 0.14)
		tween.tween_callback(flash.queue_free)

## 單次閃光效果
func _single_flash(color: Color, alpha: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.0)
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "color:a", alpha, 0.06)
	tween.tween_property(flash, "color:a", 0.0, 0.12)
	tween.tween_callback(flash.queue_free)

## 建立頂部橫幅
func _make_banner(text: String, bg_color: Color, border_color: Color) -> PanelContainer:
	var panel = PanelContainer.new()
	panel.position = Vector2(0, -60)
	panel.size = Vector2(1280, 52)

	var style = StyleBoxFlat.new()
	style.bg_color = bg_color
	style.border_color = border_color
	style.border_width_bottom = 2
	panel.add_theme_stylebox_override("panel", style)

	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 18)
	label.add_theme_color_override("font_color", Color("#FFFFFF"))
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	panel.add_child(label)

	return panel
