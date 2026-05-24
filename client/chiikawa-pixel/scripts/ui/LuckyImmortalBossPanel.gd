## LuckyImmortalBossPanel.gd — 幸運永生 BOSS 魚 UI 面板（DAY-289）
## 業界依據：Royal Fishing Jili「Immortal Boss mechanic — consecutive wins 50X-150X until they leave」
## 視覺設計：永生主題（#8B0000 深紅 + #FF0000 紅 + #FF4500 火橙 + #FFD700 金 + #FFFFFF 白）
##
## Event 類型：
##   - immortal_start：永生 BOSS 降臨（深紅三次強閃光+頂部橫幅+條命指示器）
##   - immortal_kill：永生 BOSS 被擊破（火橙閃光+浮動文字+指示器更新）
##   - immortal_revive：永生 BOSS 復活（紅色閃光+復活提示）
##   - immortal_end：永生終結（全螢幕三次強閃光+「永生終結！全服×3.5！」大字+終結指示器）
##   - immortal_end_expire：永生終結加成結束（淡出）
##   - immortal_timeout：永生 BOSS 超時消散（灰色提示）

extends CanvasLayer

const COLOR_DEEP_RED   := Color(0.55, 0.0,  0.0,  1.0)  # #8B0000 深紅
const COLOR_RED        := Color(1.0,  0.0,  0.0,  1.0)  # #FF0000 紅
const COLOR_FIRE       := Color(1.0,  0.27, 0.0,  1.0)  # #FF4500 火橙
const COLOR_GOLD       := Color(1.0,  0.84, 0.0,  1.0)  # #FFD700 金
const COLOR_WHITE      := Color(1.0,  1.0,  1.0,  1.0)  # 白
const COLOR_GRAY       := Color(0.5,  0.5,  0.5,  1.0)  # 灰

var _lives_left: int = 5
var _kill_count: int = 0
var _current_mult: float = 2.0
var _end_active: bool = false
var _end_remaining: float = 0.0

var _banner: Control = null
var _indicator: Control = null
var _end_indicator: Control = null

func _ready() -> void:
	layer = 62
	_build_ui()

func _build_ui() -> void:
	_banner = _make_banner()
	add_child(_banner)
	_banner.visible = false

	_indicator = _make_indicator()
	add_child(_indicator)
	_indicator.visible = false

	_end_indicator = _make_end_indicator()
	add_child(_end_indicator)
	_end_indicator.visible = false

func _process(delta: float) -> void:
	if _end_active and _end_remaining > 0.0:
		_end_remaining -= delta
		if _end_remaining < 0.0:
			_end_remaining = 0.0
		_update_end_indicator()

func handle_event(data: Dictionary) -> void:
	var event: String = data.get("event", "")
	match event:
		"immortal_start":      _on_immortal_start(data)
		"immortal_kill":       _on_immortal_kill(data)
		"immortal_revive":     _on_immortal_revive(data)
		"immortal_end":        _on_immortal_end(data)
		"immortal_end_expire": _on_immortal_end_expire()
		"immortal_timeout":    _on_immortal_timeout(data)

func _on_immortal_start(data: Dictionary) -> void:
	var player_name: String = data.get("trigger_player_name", "")
	_lives_left = data.get("lives_left", 5)
	_kill_count = 0
	_current_mult = data.get("current_mult", 2.0)

	_flash_screen(COLOR_DEEP_RED, 1, 0.12)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_RED, 1, 0.12)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_FIRE, 1, 0.12)

	_show_banner(
		"⚡💀 永生 BOSS 降臨！",
		"%s 召喚了永生 BOSS！5 條命，打完全服 ×3.5！" % player_name,
		COLOR_DEEP_RED, 3.5
	)

	_indicator.visible = true
	_update_indicator()

func _on_immortal_kill(data: Dictionary) -> void:
	var player_name: String = data.get("trigger_player_name", "")
	var kill_mult: float = data.get("kill_mult", 2.0)
	_kill_count = data.get("kill_count", _kill_count + 1)
	_lives_left = data.get("lives_left", _lives_left - 1)
	_current_mult = data.get("next_mult", _current_mult + 0.5)

	_flash_screen(COLOR_FIRE, 1, 0.08)
	_spawn_float_text("💀 ×%.1f！" % kill_mult, COLOR_GOLD)
	_update_indicator()

func _on_immortal_revive(data: Dictionary) -> void:
	_lives_left = data.get("lives_left", _lives_left)
	_current_mult = data.get("current_mult", _current_mult)

	_flash_screen(COLOR_RED, 1, 0.06)
	_spawn_float_text("💀 復活！×%.1f" % _current_mult, COLOR_RED)
	_update_indicator()

func _on_immortal_end(data: Dictionary) -> void:
	var player_name: String = data.get("trigger_player_name", "")
	var global_mult: float = data.get("global_mult", 3.5)
	var duration: int = data.get("global_duration_sec", 8)
	_kill_count = data.get("kill_count", _kill_count)
	_end_active = true
	_end_remaining = float(duration)

	_indicator.visible = false

	# 全螢幕三次強閃光
	_flash_screen(COLOR_DEEP_RED, 1, 0.15)
	await get_tree().create_timer(0.18).timeout
	_flash_screen(COLOR_RED, 1, 0.15)
	await get_tree().create_timer(0.18).timeout
	_flash_screen(COLOR_GOLD, 1, 0.15)

	_show_banner(
		"💀 永生終結！",
		"%s 擊敗永生 BOSS %d 次！全服 ×%.1f 加成 %d 秒！" % [player_name, _kill_count, global_mult, duration],
		COLOR_DEEP_RED, 4.0
	)

	_end_indicator.visible = true
	_update_end_indicator()

func _on_immortal_end_expire() -> void:
	_end_active = false
	var tween := create_tween()
	tween.tween_property(_end_indicator, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func(): _end_indicator.visible = false)

func _on_immortal_timeout(data: Dictionary) -> void:
	var kill_count: int = data.get("kill_count", 0)
	var lives_left: int = data.get("lives_left", 0)
	_indicator.visible = false
	_spawn_float_text("⏰ 消散（擊破 %d 次）" % kill_count, COLOR_GRAY)

# ---- UI 建構 ----

func _make_banner() -> Control:
	var c := ColorRect.new()
	c.color = Color(0, 0, 0, 0.75)
	c.size = Vector2(700, 72)
	c.position = Vector2(50, 8)
	var lbl := Label.new()
	lbl.name = "Title"
	lbl.text = "⚡💀 永生 BOSS！"
	lbl.add_theme_color_override("font_color", COLOR_GOLD)
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

func _make_indicator() -> Control:
	var c := ColorRect.new()
	c.color = Color(0.1, 0.0, 0.0, 0.85)
	c.size = Vector2(130, 70)
	c.position = Vector2(660, 90)
	var title := Label.new()
	title.text = "💀 永生 BOSS"
	title.add_theme_color_override("font_color", COLOR_GOLD)
	title.add_theme_font_size_override("font_size", 11)
	title.position = Vector2(6, 4)
	c.add_child(title)
	var lives_lbl := Label.new()
	lives_lbl.name = "Lives"
	lives_lbl.text = "條命：5"
	lives_lbl.add_theme_color_override("font_color", COLOR_RED)
	lives_lbl.add_theme_font_size_override("font_size", 13)
	lives_lbl.position = Vector2(6, 22)
	c.add_child(lives_lbl)
	var mult_lbl := Label.new()
	mult_lbl.name = "Mult"
	mult_lbl.text = "下次：×2.0"
	mult_lbl.add_theme_color_override("font_color", COLOR_FIRE)
	mult_lbl.add_theme_font_size_override("font_size", 13)
	mult_lbl.position = Vector2(6, 42)
	c.add_child(mult_lbl)
	return c

func _make_end_indicator() -> Control:
	var c := ColorRect.new()
	c.color = Color(0.1, 0.0, 0.0, 0.9)
	c.size = Vector2(130, 70)
	c.position = Vector2(660, 90)
	var title := Label.new()
	title.text = "💀 永生終結"
	title.add_theme_color_override("font_color", COLOR_GOLD)
	title.add_theme_font_size_override("font_size", 11)
	title.position = Vector2(6, 4)
	c.add_child(title)
	var mult_lbl := Label.new()
	mult_lbl.name = "Mult"
	mult_lbl.text = "全服 ×3.5"
	mult_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	mult_lbl.add_theme_font_size_override("font_size", 15)
	mult_lbl.position = Vector2(6, 22)
	c.add_child(mult_lbl)
	var timer_lbl := Label.new()
	timer_lbl.name = "Timer"
	timer_lbl.text = "8.0s"
	timer_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	timer_lbl.add_theme_font_size_override("font_size", 13)
	timer_lbl.position = Vector2(6, 46)
	c.add_child(timer_lbl)
	return c

# ---- 更新 UI ----

func _update_indicator() -> void:
	if not is_instance_valid(_indicator):
		return
	var lives_lbl := _indicator.get_node_or_null("Lives")
	if lives_lbl:
		lives_lbl.text = "條命：%d" % _lives_left
		# 顏色隨條命數變化
		if _lives_left >= 4:
			lives_lbl.add_theme_color_override("font_color", COLOR_RED)
		elif _lives_left >= 2:
			lives_lbl.add_theme_color_override("font_color", COLOR_FIRE)
		else:
			lives_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	var mult_lbl := _indicator.get_node_or_null("Mult")
	if mult_lbl:
		mult_lbl.text = "下次：×%.1f" % _current_mult

func _update_end_indicator() -> void:
	if not is_instance_valid(_end_indicator):
		return
	var timer_lbl := _end_indicator.get_node_or_null("Timer")
	if timer_lbl:
		timer_lbl.text = "%.1fs" % _end_remaining
		# 脈衝動畫
		var alpha := 0.7 + 0.3 * sin(Time.get_ticks_msec() * 0.005)
		_end_indicator.modulate.a = alpha

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
