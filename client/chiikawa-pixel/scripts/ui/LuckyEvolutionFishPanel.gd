## LuckyEvolutionFishPanel.gd — 幸運進化魚系統面板（DAY-218）
## 業界原創「三段進化」機制
##
## 視覺設計：
##   - 三段進化主題（綠→青→紫→爆炸）
##   - evolution_appear：綠色雙閃光 + 頂部橫幅 + 命中進度條
##   - evolution_hit：小閃光 + 進度條更新 + 「命中 N/9」提示
##   - evolution_stage：對應顏色強閃光 + 「🌟 進化！×N 倍率！」大字 + 進化魚視覺變化
##   - evolution_burst：全螢幕三次強閃光（紫→白）+ 「💥 終極爆發！」52px大字 + 倍率計時條
##   - evolution_burst_end：計時條淡出
##   - evolution_escape：灰色淡出 + 「😔 進化魚逃跑了...」提示
extends CanvasLayer

# 進化狀態
var _evolution_active: bool = false
var _current_stage: int = 0
var _hit_count: int = 0
var _next_hit: int = 3

# 進化 UI 節點
var _evolution_banner: Control = null
var _hit_bar: Control = null
var _boost_bar: Control = null
var _boost_timer: float = 0.0
var _boost_duration: float = 6.0

# 進化階段顏色
const STAGE_COLORS = [
	Color("#00FF88"),  # 進化 1：綠色
	Color("#00CCFF"),  # 進化 2：青色
	Color("#FF00FF"),  # 進化 3：紫色
]

func _ready() -> void:
	layer = 27  # 幸運進化魚面板層級

func _process(delta: float) -> void:
	if _boost_timer > 0.0:
		_boost_timer -= delta
		if _boost_timer <= 0.0:
			_boost_timer = 0.0
			_cleanup_boost_bar()
		else:
			_update_boost_bar()

## 處理幸運進化魚訊息
func handle_lucky_evolution_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"evolution_appear":
			_on_evolution_appear(payload)
		"evolution_hit":
			_on_evolution_hit(payload)
		"evolution_stage":
			_on_evolution_stage(payload)
		"evolution_burst", "evolution_kill_burst":
			_on_evolution_burst(payload)
		"evolution_burst_end":
			_on_evolution_burst_end()
		"evolution_escape":
			_on_evolution_escape(payload)

## 進化魚出現 — 綠色雙閃光 + 頂部橫幅 + 命中進度條
func _on_evolution_appear(payload: Dictionary) -> void:
	_evolution_active = true
	_current_stage = 0
	_hit_count = 0
	_next_hit = payload.get("next_hit", 3)

	# 綠色雙閃光
	_double_flash(Color("#00FF88"), 0.40)

	# 頂部橫幅
	var msg = "🌟 幸運進化魚出現！命中 3 次觸發進化！進化後倍率提升，終極進化後全場爆發！"
	_evolution_banner = _make_banner(msg, Color(0.0, 0.08, 0.04, 0.88), Color("#00FF88"))
	add_child(_evolution_banner)

	# 命中進度條
	_hit_bar = _make_hit_bar()
	add_child(_hit_bar)

## 進化魚被命中 — 小閃光 + 進度條更新
func _on_evolution_hit(payload: Dictionary) -> void:
	_hit_count = payload.get("hit_count", _hit_count + 1)
	_next_hit = payload.get("next_hit", 9)
	var player_name: String = payload.get("player_name", "")

	# 小閃光（依當前階段顏色）
	var flash_color = STAGE_COLORS[min(_current_stage, 2)]
	_single_flash(flash_color, 0.15)

	# 更新進度條
	_update_hit_bar()

	# 浮動文字
	_spawn_float_text("💥 %s 命中！%d/%d" % [player_name, _hit_count, _next_hit], flash_color, 24)

## 進化觸發 — 對應顏色強閃光 + 大字
func _on_evolution_stage(payload: Dictionary) -> void:
	_current_stage = payload.get("stage", 1)
	var stage_name: String = payload.get("stage_name", "進化！")
	var mult_boost: float = payload.get("mult_boost", 1.5)
	var player_name: String = payload.get("player_name", "")

	var stage_color = STAGE_COLORS[min(_current_stage - 1, 2)]

	# 強閃光
	_double_flash(stage_color, 0.50)

	# 大字
	var text = "🌟 %s 觸發%s 倍率 ×%.1f！" % [player_name, stage_name, mult_boost]
	var big_label = _make_big_label(text, stage_color, 40)
	add_child(big_label)
	var tw = big_label.create_tween()
	tw.tween_interval(2.0)
	tw.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tw.tween_callback(func(): if is_instance_valid(big_label): big_label.queue_free())

	# 更新橫幅
	if is_instance_valid(_evolution_banner):
		var label = _evolution_banner.get_node_or_null("Label")
		if label:
			var stage_text = ["", "一段進化", "二段進化", "終極進化"][_current_stage]
			label.text = "🌟 進化魚【%s】倍率 ×%.1f！繼續命中觸發下一段進化！" % [stage_text, mult_boost]
			label.add_theme_color_override("font_color", stage_color)

	# 更新進度條顏色
	_update_hit_bar()

## 終極爆發 — 全螢幕三次強閃光 + 大字 + 倍率計時條
func _on_evolution_burst(payload: Dictionary) -> void:
	_evolution_active = false
	_cleanup_evolution_ui()

	var mult_boost: float = payload.get("mult_boost", 4.0)
	var boost_sec: int = payload.get("boost_sec", 6)
	var affected_count: int = payload.get("affected_count", 0)
	var is_kill_burst: bool = payload.get("event", "") == "evolution_kill_burst"

	# 全螢幕三次強閃光（紫→白）
	_triple_flash_evolution()

	# 大字
	var text: String
	if is_kill_burst:
		text = "💥 進化魚被擊破！提前引爆！×%.1f 倍率加成！" % mult_boost
	else:
		text = "💥 終極爆發！全場 HP -60%%！×%.1f 倍率加成！" % mult_boost
	var big_label = _make_big_label(text, Color("#FF00FF"), 48)
	add_child(big_label)
	var tw = big_label.create_tween()
	tw.tween_interval(2.5)
	tw.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tw.tween_callback(func(): if is_instance_valid(big_label): big_label.queue_free())

	# 副標題
	var sub_text = "全場 %d 個目標 HP -60%% | ×%.1f 倍率加成 %d 秒" % [affected_count, mult_boost, boost_sec]
	var sub_label = _make_big_label(sub_text, Color("#E0B0FF"), 24)
	sub_label.position.y += 60
	add_child(sub_label)
	var tw2 = sub_label.create_tween()
	tw2.tween_interval(2.5)
	tw2.tween_property(sub_label, "modulate:a", 0.0, 0.4)
	tw2.tween_callback(func(): if is_instance_valid(sub_label): sub_label.queue_free())

	# 倍率計時條
	_boost_duration = float(boost_sec)
	_boost_timer = _boost_duration
	_boost_bar = _make_boost_bar(mult_boost, boost_sec)
	add_child(_boost_bar)

## 倍率加成結束
func _on_evolution_burst_end() -> void:
	_boost_timer = 0.0
	_cleanup_boost_bar()

## 進化魚逃跑
func _on_evolution_escape(payload: Dictionary) -> void:
	_evolution_active = false
	_cleanup_evolution_ui()

	var stage: int = payload.get("stage", 0)
	var stage_names = ["未進化", "一段進化", "二段進化", "終極進化"]
	var stage_name = stage_names[min(stage, 3)]

	_single_flash(Color("#888888"), 0.3)
	var label = _make_big_label("😔 進化魚逃跑了...（%s）" % stage_name, Color("#AAAAAA"), 28)
	add_child(label)
	var tw = label.create_tween()
	tw.tween_interval(2.0)
	tw.tween_property(label, "modulate:a", 0.0, 0.5)
	tw.tween_callback(func(): if is_instance_valid(label): label.queue_free())

# ─── 內部 UI 工具函數 ───────────────────────────────────────────────────────

func _cleanup_evolution_ui() -> void:
	if is_instance_valid(_evolution_banner):
		_evolution_banner.queue_free()
		_evolution_banner = null
	if is_instance_valid(_hit_bar):
		_hit_bar.queue_free()
		_hit_bar = null

func _cleanup_boost_bar() -> void:
	if is_instance_valid(_boost_bar):
		var tw = _boost_bar.create_tween()
		tw.tween_property(_boost_bar, "modulate:a", 0.0, 0.3)
		tw.tween_callback(func(): if is_instance_valid(_boost_bar): _boost_bar.queue_free())
		_boost_bar = null

func _make_hit_bar() -> Control:
	var container = Control.new()
	container.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	container.position.y = -52
	container.size = Vector2(get_viewport().get_visible_rect().size.x, 24)

	var bg = ColorRect.new()
	bg.color = Color(0.05, 0.05, 0.05, 0.75)
	bg.size = Vector2(container.size.x, 20)
	bg.position = Vector2(0, 2)
	container.add_child(bg)

	var bar = ColorRect.new()
	bar.name = "Bar"
	bar.color = Color("#00FF88")
	bar.size = Vector2(0.0, 14)
	bar.position = Vector2((container.size.x - 480.0) / 2.0, 5)
	container.add_child(bar)

	var label = Label.new()
	label.name = "HitLabel"
	label.text = "命中 0/3"
	label.add_theme_color_override("font_color", Color("#00FF88"))
	label.add_theme_font_size_override("font_size", 14)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.set_anchors_preset(Control.PRESET_FULL_RECT)
	container.add_child(label)

	return container

func _update_hit_bar() -> void:
	if not is_instance_valid(_hit_bar):
		return
	var bar = _hit_bar.get_node_or_null("Bar")
	var label = _hit_bar.get_node_or_null("HitLabel")
	if not bar or not label:
		return

	var max_hit = 9  # 終極進化需要 9 次命中
	var pct = float(_hit_count) / float(max_hit)
	bar.size.x = 480.0 * pct

	# 顏色依進化階段
	var stage_color = STAGE_COLORS[min(_current_stage, 2)]
	bar.color = stage_color
	label.add_theme_color_override("font_color", stage_color)
	label.text = "命中 %d/%d（第 %d 段進化）" % [_hit_count, _next_hit, _current_stage + 1]

func _make_boost_bar(mult_boost: float, boost_sec: int) -> Control:
	var container = Control.new()
	container.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	container.position.y = -28
	container.size = Vector2(get_viewport().get_visible_rect().size.x, 20)

	var bg = ColorRect.new()
	bg.color = Color(0.1, 0.0, 0.1, 0.75)
	bg.size = Vector2(container.size.x, 16)
	bg.position = Vector2(0, 2)
	container.add_child(bg)

	var bar = ColorRect.new()
	bar.name = "Bar"
	bar.color = Color("#FF00FF")
	bar.size = Vector2(480.0, 12)
	bar.position = Vector2((container.size.x - 480.0) / 2.0, 4)
	container.add_child(bar)

	var label = Label.new()
	label.name = "BoostLabel"
	label.text = "🌟 ×%.1f 倍率加成 %ds" % [mult_boost, boost_sec]
	label.add_theme_color_override("font_color", Color("#FF00FF"))
	label.add_theme_font_size_override("font_size", 14)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.set_anchors_preset(Control.PRESET_FULL_RECT)
	container.add_child(label)

	return container

func _update_boost_bar() -> void:
	if not is_instance_valid(_boost_bar):
		return
	var bar = _boost_bar.get_node_or_null("Bar")
	var label = _boost_bar.get_node_or_null("BoostLabel")
	if not bar:
		return
	var pct = _boost_timer / _boost_duration
	bar.size.x = 480.0 * pct
	# 顏色漸變：紫→藍紫→深藍
	if pct > 0.6:
		bar.color = Color("#FF00FF")
	elif pct > 0.3:
		bar.color = Color("#8800FF")
	else:
		bar.color = Color("#4400CC")
	if label:
		label.text = "🌟 ×4.0 倍率加成 %.1fs" % _boost_timer

func _make_banner(text: String, bg_color: Color, text_color: Color) -> Control:
	var panel = PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_WIDE)
	panel.position.y = 8
	panel.size.x = 620
	panel.position.x = (get_viewport().get_visible_rect().size.x - 620) / 2.0

	var style = StyleBoxFlat.new()
	style.bg_color = bg_color
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	panel.add_theme_stylebox_override("panel", style)

	var label = Label.new()
	label.name = "Label"
	label.text = text
	label.add_theme_color_override("font_color", text_color)
	label.add_theme_font_size_override("font_size", 17)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	panel.add_child(label)
	return panel

func _make_big_label(text: String, color: Color, font_size: int) -> Label:
	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", font_size)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.set_anchors_preset(Control.PRESET_CENTER)
	label.position.y = get_viewport().get_visible_rect().size.y * 0.35
	label.position.x = -310
	label.size.x = 620
	return label

func _spawn_float_text(text: String, color: Color, font_size: int) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", font_size)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	var vp_size = get_viewport().get_visible_rect().size
	label.position = Vector2(vp_size.x / 2.0 - 200, vp_size.y * 0.50)
	label.size.x = 400
	add_child(label)
	var tw = label.create_tween()
	tw.tween_property(label, "position:y", label.position.y - 50, 1.0)
	tw.parallel().tween_property(label, "modulate:a", 0.0, 1.0)
	tw.tween_callback(func(): if is_instance_valid(label): label.queue_free())

func _double_flash(color: Color, duration: float) -> void:
	var overlay = ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, 0.0)
	overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tw = overlay.create_tween()
	tw.tween_property(overlay, "color:a", 0.55, duration * 0.3)
	tw.tween_property(overlay, "color:a", 0.0, duration * 0.3)
	tw.tween_property(overlay, "color:a", 0.45, duration * 0.2)
	tw.tween_property(overlay, "color:a", 0.0, duration * 0.2)
	tw.tween_callback(func(): if is_instance_valid(overlay): overlay.queue_free())

func _single_flash(color: Color, duration: float) -> void:
	var overlay = ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, 0.0)
	overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tw = overlay.create_tween()
	tw.tween_property(overlay, "color:a", 0.35, duration * 0.4)
	tw.tween_property(overlay, "color:a", 0.0, duration * 0.6)
	tw.tween_callback(func(): if is_instance_valid(overlay): overlay.queue_free())

func _triple_flash_evolution() -> void:
	# 紫→白→紫 三次強閃光
	var overlay = ColorRect.new()
	overlay.color = Color(1.0, 0.0, 1.0, 0.0)
	overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tw = overlay.create_tween()
	tw.tween_property(overlay, "color:a", 0.75, 0.10)
	tw.tween_property(overlay, "color:a", 0.0, 0.12)
	tw.tween_property(overlay, "color", Color(1.0, 1.0, 1.0, 0.0), 0.01)
	tw.tween_property(overlay, "color:a", 0.65, 0.10)
	tw.tween_property(overlay, "color:a", 0.0, 0.12)
	tw.tween_property(overlay, "color", Color(1.0, 0.0, 1.0, 0.0), 0.01)
	tw.tween_property(overlay, "color:a", 0.55, 0.10)
	tw.tween_property(overlay, "color:a", 0.0, 0.15)
	tw.tween_callback(func(): if is_instance_valid(overlay): overlay.queue_free())
