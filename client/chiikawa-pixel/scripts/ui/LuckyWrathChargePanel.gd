## LuckyWrathChargePanel.gd — 幸運怒氣蓄積魚 UI 面板（DAY-290）
## 業界依據：Royal Fishing Jili「Dragon Wrath system accumulates with every shot fired」
## 視覺設計：怒氣主題（#8B0000 深紅 + #FF4500 火橙 + #FF6B35 橙紅 + #FFD700 金 + #FFFFFF 白）
##
## Event 類型：
##   - wrath_start：怒氣蓄積開始（深紅三次強閃光+頂部橫幅+怒氣條指示器）
##   - wrath_charge：怒氣值更新（橙紅閃光+怒氣條更新）
##   - wrath_explode：怒氣爆發（全螢幕三次強閃光+「怒氣爆發！X顆隕石！」大字）
##   - wrath_meteor：單顆隕石（閃光+浮動文字）
##   - wrath_settle：爆發結算（結算彈窗）
##   - wrath_perfect：完美怒氣全服加成（全螢幕三次強閃光+「完美怒氣！全服×2.8！」大字+終結指示器）
##   - wrath_perfect_end：完美怒氣加成結束（淡出）

extends CanvasLayer

const COLOR_DEEP_RED   := Color(0.55, 0.0,  0.0,  1.0)  # #8B0000 深紅
const COLOR_FIRE       := Color(1.0,  0.27, 0.0,  1.0)  # #FF4500 火橙
const COLOR_ORANGE_RED := Color(1.0,  0.42, 0.21, 1.0)  # #FF6B35 橙紅
const COLOR_GOLD       := Color(1.0,  0.84, 0.0,  1.0)  # #FFD700 金
const COLOR_WHITE      := Color(1.0,  1.0,  1.0,  1.0)  # 白

var _wrath_value: int = 0
var _max_wrath: int = 20
var _timeout_sec: int = 25
var _perfect_active: bool = false
var _perfect_remaining: float = 0.0

var _banner: Control = null
var _wrath_indicator: Control = null
var _perfect_indicator: Control = null

func _ready() -> void:
	layer = 63
	_build_ui()

func _build_ui() -> void:
	_banner = _make_banner()
	add_child(_banner)
	_banner.visible = false

	_wrath_indicator = _make_wrath_indicator()
	add_child(_wrath_indicator)
	_wrath_indicator.visible = false

	_perfect_indicator = _make_perfect_indicator()
	add_child(_perfect_indicator)
	_perfect_indicator.visible = false

func _process(delta: float) -> void:
	if _perfect_active and _perfect_remaining > 0.0:
		_perfect_remaining -= delta
		if _perfect_remaining < 0.0:
			_perfect_remaining = 0.0
		_update_perfect_indicator()

func handle_event(data: Dictionary) -> void:
	var event: String = data.get("event", "")
	match event:
		"wrath_start":      _on_wrath_start(data)
		"wrath_charge":     _on_wrath_charge(data)
		"wrath_explode":    _on_wrath_explode(data)
		"wrath_meteor":     _on_wrath_meteor(data)
		"wrath_settle":     _on_wrath_settle(data)
		"wrath_perfect":    _on_wrath_perfect(data)
		"wrath_perfect_end": _on_wrath_perfect_end()

func _on_wrath_start(data: Dictionary) -> void:
	var player_name: String = data.get("player_name", "")
	_wrath_value = 0
	_max_wrath = data.get("max_wrath", 20)
	_timeout_sec = data.get("timeout_sec", 25)

	_flash_screen(COLOR_DEEP_RED, 1, 0.12)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_FIRE, 1, 0.12)

	_show_banner(
		"🔥💥 怒氣蓄積！",
		"%s 開始蓄積怒氣！打越多魚，隕石越多！最多 %d 顆！" % [player_name, _max_wrath],
		COLOR_DEEP_RED, 3.5
	)

	_wrath_indicator.visible = true
	_update_wrath_indicator()

func _on_wrath_charge(data: Dictionary) -> void:
	_wrath_value = data.get("wrath_value", _wrath_value + 1)
	_flash_screen(COLOR_ORANGE_RED, 1, 0.04)
	_update_wrath_indicator()

func _on_wrath_explode(data: Dictionary) -> void:
	var meteor_count: int = data.get("meteor_count", 1)
	var is_perfect: bool = data.get("is_perfect", false)

	_wrath_indicator.visible = false

	_flash_screen(COLOR_DEEP_RED, 1, 0.15)
	await get_tree().create_timer(0.18).timeout
	_flash_screen(COLOR_FIRE, 1, 0.15)
	await get_tree().create_timer(0.18).timeout
	_flash_screen(COLOR_GOLD, 1, 0.15)

	var title_text := "🔥💥 怒氣爆發！"
	if is_perfect:
		title_text = "🔥💥 完美怒氣爆發！"
	_show_banner(
		title_text,
		"%d 顆隕石即將落下！" % meteor_count,
		COLOR_FIRE, 2.5
	)

func _on_wrath_meteor(data: Dictionary) -> void:
	var hit_count: int = data.get("hit_count", 0)
	if hit_count > 0:
		_flash_screen(COLOR_FIRE, 1, 0.05)
		_spawn_float_text("💥 命中 %d！" % hit_count, COLOR_GOLD)

func _on_wrath_settle(data: Dictionary) -> void:
	var wrath_value: int = data.get("wrath_value", 0)
	var meteor_count: int = data.get("meteor_count", 0)
	var total_hit: int = data.get("total_hit", 0)
	var is_perfect: bool = data.get("is_perfect", false)

	if total_hit >= 5 or is_perfect:
		_show_banner(
			"🔥 怒氣結算",
			"怒氣 %d/%d，%d 顆隕石，命中 %d 個目標！" % [wrath_value, _max_wrath, meteor_count, total_hit],
			COLOR_GOLD, 3.0
		)

func _on_wrath_perfect(data: Dictionary) -> void:
	var player_name: String = data.get("player_name", "")
	var global_mult: float = data.get("global_mult", 2.8)
	var duration: int = data.get("global_duration_sec", 7)
	_perfect_active = true
	_perfect_remaining = float(duration)

	_flash_screen(COLOR_DEEP_RED, 1, 0.15)
	await get_tree().create_timer(0.18).timeout
	_flash_screen(COLOR_FIRE, 1, 0.15)
	await get_tree().create_timer(0.18).timeout
	_flash_screen(COLOR_GOLD, 1, 0.15)

	_show_banner(
		"🔥 完美怒氣！",
		"%s 完美怒氣！全服 ×%.1f 加成 %d 秒！" % [player_name, global_mult, duration],
		COLOR_FIRE, 4.0
	)

	_perfect_indicator.visible = true
	_update_perfect_indicator()

func _on_wrath_perfect_end() -> void:
	_perfect_active = false
	var tween := create_tween()
	tween.tween_property(_perfect_indicator, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func(): _perfect_indicator.visible = false)

# ---- UI 建構 ----

func _make_banner() -> Control:
	var c := ColorRect.new()
	c.color = Color(0, 0, 0, 0.75)
	c.size = Vector2(700, 72)
	c.position = Vector2(50, 8)
	var lbl := Label.new()
	lbl.name = "Title"
	lbl.text = "🔥💥 怒氣蓄積！"
	lbl.add_theme_color_override("font_color", COLOR_FIRE)
	lbl.add_theme_font_size_override("font_size", 20)
	lbl.position = Vector2(12, 4)
	c.add_child(lbl)
	var sub := Label.new()
	sub.name = "Sub"
	sub.text = ""
	sub.add_theme_color_override("font_color", COLOR_WHITE)
	sub.add_theme_font_size_override("font_size", 13)
	sub.position = Vector2(12, 30)
	sub.size = Vector2(676, 36)
	sub.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	c.add_child(sub)
	return c

func _make_wrath_indicator() -> Control:
	var c := ColorRect.new()
	c.color = Color(0.1, 0.0, 0.0, 0.85)
	c.size = Vector2(130, 80)
	c.position = Vector2(660, 90)
	var title := Label.new()
	title.text = "🔥 怒氣蓄積"
	title.add_theme_color_override("font_color", COLOR_FIRE)
	title.add_theme_font_size_override("font_size", 11)
	title.position = Vector2(6, 4)
	c.add_child(title)
	var wrath_lbl := Label.new()
	wrath_lbl.name = "Wrath"
	wrath_lbl.text = "0/20"
	wrath_lbl.add_theme_color_override("font_color", COLOR_ORANGE_RED)
	wrath_lbl.add_theme_font_size_override("font_size", 18)
	wrath_lbl.position = Vector2(6, 22)
	c.add_child(wrath_lbl)
	var bar := ColorRect.new()
	bar.name = "WrathBar"
	bar.color = COLOR_FIRE
	bar.size = Vector2(0, 8)
	bar.position = Vector2(6, 50)
	c.add_child(bar)
	var bar_bg := ColorRect.new()
	bar_bg.color = Color(0.3, 0.1, 0.0, 1.0)
	bar_bg.size = Vector2(118, 8)
	bar_bg.position = Vector2(6, 50)
	c.add_child(bar_bg)
	c.move_child(bar, c.get_child_count() - 1)
	var hint := Label.new()
	hint.text = "打魚蓄積怒氣！"
	hint.add_theme_color_override("font_color", COLOR_WHITE)
	hint.add_theme_font_size_override("font_size", 10)
	hint.position = Vector2(6, 62)
	c.add_child(hint)
	return c

func _make_perfect_indicator() -> Control:
	var c := ColorRect.new()
	c.color = Color(0.1, 0.0, 0.0, 0.9)
	c.size = Vector2(130, 70)
	c.position = Vector2(660, 90)
	var title := Label.new()
	title.text = "🔥 完美怒氣"
	title.add_theme_color_override("font_color", COLOR_GOLD)
	title.add_theme_font_size_override("font_size", 11)
	title.position = Vector2(6, 4)
	c.add_child(title)
	var mult_lbl := Label.new()
	mult_lbl.name = "Mult"
	mult_lbl.text = "全服 ×2.8"
	mult_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	mult_lbl.add_theme_font_size_override("font_size", 15)
	mult_lbl.position = Vector2(6, 22)
	c.add_child(mult_lbl)
	var timer_lbl := Label.new()
	timer_lbl.name = "Timer"
	timer_lbl.text = "7.0s"
	timer_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	timer_lbl.add_theme_font_size_override("font_size", 13)
	timer_lbl.position = Vector2(6, 46)
	c.add_child(timer_lbl)
	return c

# ---- 更新 UI ----

func _update_wrath_indicator() -> void:
	if not is_instance_valid(_wrath_indicator):
		return
	var wrath_lbl := _wrath_indicator.get_node_or_null("Wrath")
	if wrath_lbl:
		wrath_lbl.text = "%d/%d" % [_wrath_value, _max_wrath]
		var pct := float(_wrath_value) / float(_max_wrath)
		if pct >= 0.75:
			wrath_lbl.add_theme_color_override("font_color", COLOR_GOLD)
		elif pct >= 0.5:
			wrath_lbl.add_theme_color_override("font_color", COLOR_FIRE)
		else:
			wrath_lbl.add_theme_color_override("font_color", COLOR_ORANGE_RED)
	var bar := _wrath_indicator.get_node_or_null("WrathBar")
	if bar:
		var pct := float(_wrath_value) / float(_max_wrath)
		bar.size.x = 118.0 * pct

func _update_perfect_indicator() -> void:
	if not is_instance_valid(_perfect_indicator):
		return
	var timer_lbl := _perfect_indicator.get_node_or_null("Timer")
	if timer_lbl:
		timer_lbl.text = "%.1fs" % _perfect_remaining
		var alpha := 0.7 + 0.3 * sin(Time.get_ticks_msec() * 0.005)
		_perfect_indicator.modulate.a = alpha

# ---- 特效工具 ----

func _show_banner(title: String, sub: String, color: Color, duration: float) -> void:
	if not is_instance_valid(_banner):
		return
	var title_lbl := _banner.get_node_or_null("Title")
	if title_lbl:
		title_lbl.text = title
		title_lbl.add_theme_color_override("font_color", color)
	var sub_lbl := _banner.get_node_or_null("Sub")
	if sub_lbl:
		sub_lbl.text = sub
	_banner.visible = true
	_banner.modulate.a = 1.0
	await get_tree().create_timer(duration).timeout
	if is_instance_valid(_banner):
		var tween := create_tween()
		tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
		tween.tween_callback(func(): if is_instance_valid(_banner): _banner.visible = false)

func _flash_screen(color: Color, times: int, duration: float) -> void:
	var flash := ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.35)
	flash.size = Vector2(800, 600)
	flash.position = Vector2.ZERO
	add_child(flash)
	for i in range(times):
		flash.visible = true
		await get_tree().create_timer(duration).timeout
		flash.visible = false
		if i < times - 1:
			await get_tree().create_timer(0.05).timeout
	if is_instance_valid(flash):
		flash.queue_free()

func _spawn_float_text(text: String, color: Color) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 22)
	lbl.position = Vector2(300 + randf_range(-80, 80), 280)
	add_child(lbl)
	var tween := create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 80, 1.2)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 1.2)
	tween.tween_callback(func(): if is_instance_valid(lbl): lbl.queue_free())
