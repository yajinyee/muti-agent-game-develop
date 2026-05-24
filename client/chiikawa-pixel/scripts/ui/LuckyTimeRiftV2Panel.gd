## LuckyTimeRiftV2Panel.gd — 幸運時空裂縫魚 UI 面板（DAY-291）
## 業界依據：Fishing Fortune 2026「Time Freeze mechanic」
## 視覺設計：時空主題（#00E5FF 青藍 + #0D47A1 深藍 + #7B2FBE 紫 + #FFD700 金 + #FFFFFF 白）
##
## Event 類型：
##   - rift_start：時空裂縫開始（青藍三次強閃光+頂部橫幅+凍結計時條+傷害加倍指示器）
##   - rift_end：凍結結束+裂縫爆炸（全螢幕三次強閃光+「裂縫爆炸！HP-30%！」大字）
##   - rift_perfect：時空完美全服加成（全螢幕三次強閃光+「時空完美！全服×2.5！」大字+完美指示器）
##   - rift_perfect_end：時空完美加成結束（淡出）

extends CanvasLayer

const COLOR_CYAN       := Color(0.0,  0.90, 1.0,  1.0)  # #00E5FF 青藍
const COLOR_DEEP_BLUE  := Color(0.05, 0.28, 0.63, 1.0)  # #0D47A1 深藍
const COLOR_PURPLE     := Color(0.48, 0.18, 0.75, 1.0)  # #7B2FBE 紫
const COLOR_GOLD       := Color(1.0,  0.84, 0.0,  1.0)  # #FFD700 金
const COLOR_WHITE      := Color(1.0,  1.0,  1.0,  1.0)  # 白

var _freeze_duration: int = 8
var _freeze_remaining: float = 0.0
var _damage_mult: float = 2.0
var _kill_count: int = 0
var _perfect_active: bool = false
var _perfect_remaining: float = 0.0
var _is_frozen: bool = false

var _banner: Control = null
var _freeze_indicator: Control = null
var _perfect_indicator: Control = null
var _flash_overlay: ColorRect = null

func _ready() -> void:
	layer = 64
	_build_ui()

func _build_ui() -> void:
	# 全螢幕閃光覆蓋層
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(0.0, 0.9, 1.0, 0.0)
	_flash_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	_banner = _make_banner()
	add_child(_banner)
	_banner.visible = false

	_freeze_indicator = _make_freeze_indicator()
	add_child(_freeze_indicator)
	_freeze_indicator.visible = false

	_perfect_indicator = _make_perfect_indicator()
	add_child(_perfect_indicator)
	_perfect_indicator.visible = false

func _process(delta: float) -> void:
	# 凍結計時條倒數
	if _is_frozen and _freeze_remaining > 0.0:
		_freeze_remaining -= delta
		if _freeze_remaining < 0.0:
			_freeze_remaining = 0.0
		_update_freeze_bar()

	# 完美加成倒數
	if _perfect_active and _perfect_remaining > 0.0:
		_perfect_remaining -= delta
		if _perfect_remaining < 0.0:
			_perfect_remaining = 0.0
			_perfect_active = false
		_update_perfect_timer()

## 處理 Server 事件
func handle_event(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"rift_start":
			_on_rift_start(payload)
		"rift_end":
			_on_rift_end(payload)
		"rift_perfect":
			_on_rift_perfect(payload)
		"rift_perfect_end":
			_on_rift_perfect_end()

## rift_start：時空裂縫開始
func _on_rift_start(payload: Dictionary) -> void:
	_freeze_duration = payload.get("freeze_duration", 8)
	_freeze_remaining = float(_freeze_duration)
	_damage_mult = payload.get("damage_mult", 2.0)
	_kill_count = 0
	_is_frozen = true

	var trigger_name: String = payload.get("trigger_player_name", "玩家")

	# 三次青藍強閃光
	_flash_screen(COLOR_CYAN, 3)

	# 顯示頂部橫幅
	_show_banner("⏸️ 時空裂縫！傷害 ×%.1f" % _damage_mult, COLOR_CYAN)

	# 顯示凍結指示器
	_freeze_indicator.visible = true
	_update_freeze_bar()

	# 浮動文字
	_spawn_float_text("⏸️ %s 觸發時空裂縫！" % trigger_name, COLOR_CYAN, Vector2(640, 300))

## rift_end：凍結結束+裂縫爆炸
func _on_rift_end(payload: Dictionary) -> void:
	_is_frozen = false
	_freeze_remaining = 0.0
	_freeze_indicator.visible = false
	_banner.visible = false

	var kill_count: int = payload.get("kill_count", 0)
	var explosion_dmg: float = payload.get("explosion_dmg", 0.3)

	# 全螢幕三次強閃光（白色爆炸感）
	_flash_screen(COLOR_WHITE, 3)

	# 大字顯示
	_spawn_big_text("💥 裂縫爆炸！HP -%d%%！" % int(explosion_dmg * 100), COLOR_CYAN)

	# 浮動文字（擊破數）
	if kill_count > 0:
		_spawn_float_text("凍結期間擊破 %d 個目標！" % kill_count, COLOR_GOLD, Vector2(640, 400))

## rift_perfect：時空完美全服加成
func _on_rift_perfect(payload: Dictionary) -> void:
	var kill_count: int = payload.get("kill_count", 0)
	var perfect_mult: float = payload.get("perfect_mult", 2.5)
	var perfect_duration: int = payload.get("perfect_duration", 6)

	_perfect_active = true
	_perfect_remaining = float(perfect_duration)

	# 全螢幕三次強閃光（金色）
	_flash_screen(COLOR_GOLD, 3)

	# 大字顯示
	_spawn_big_text("⏸️ 時空完美！全服 ×%.1f！" % perfect_mult, COLOR_GOLD)

	# 顯示完美指示器
	_perfect_indicator.visible = true
	_update_perfect_timer()

	# 浮動文字
	_spawn_float_text("擊破 %d 個目標！時空完美！" % kill_count, COLOR_GOLD, Vector2(640, 350))

## rift_perfect_end：時空完美加成結束
func _on_rift_perfect_end() -> void:
	_perfect_active = false
	_perfect_remaining = 0.0

	# 淡出完美指示器
	var tween := create_tween()
	tween.tween_property(_perfect_indicator, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func(): _perfect_indicator.visible = false; _perfect_indicator.modulate.a = 1.0)

# ─── UI 建構 ───────────────────────────────────────────────────────────────

func _make_banner() -> Control:
	var c := Control.new()
	c.set_anchors_preset(Control.PRESET_TOP_WIDE)
	c.custom_minimum_size = Vector2(0, 56)

	var bg := ColorRect.new()
	bg.color = Color(0.0, 0.9, 1.0, 0.85)
	bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	c.add_child(bg)

	var lbl := Label.new()
	lbl.name = "Label"
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.set_anchors_preset(Control.PRESET_FULL_RECT)
	lbl.add_theme_color_override("font_color", COLOR_WHITE)
	lbl.add_theme_font_size_override("font_size", 22)
	c.add_child(lbl)
	return c

func _show_banner(text: String, color: Color) -> void:
	var bg := _banner.get_child(0) as ColorRect
	bg.color = Color(color.r, color.g, color.b, 0.85)
	var lbl := _banner.get_node("Label") as Label
	lbl.text = text
	_banner.visible = true
	# 5 秒後淡出
	var tween := create_tween()
	tween.tween_interval(5.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func(): _banner.visible = false; _banner.modulate.a = 1.0)

func _make_freeze_indicator() -> Control:
	var c := Control.new()
	c.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	c.position = Vector2(-220, 80)
	c.custom_minimum_size = Vector2(200, 70)

	var bg := ColorRect.new()
	bg.color = Color(0.0, 0.1, 0.2, 0.85)
	bg.size = Vector2(200, 70)
	c.add_child(bg)

	var title_lbl := Label.new()
	title_lbl.name = "TitleLabel"
	title_lbl.text = "⏸️ 時空裂縫"
	title_lbl.position = Vector2(8, 4)
	title_lbl.add_theme_color_override("font_color", COLOR_CYAN)
	title_lbl.add_theme_font_size_override("font_size", 13)
	c.add_child(title_lbl)

	var dmg_lbl := Label.new()
	dmg_lbl.name = "DmgLabel"
	dmg_lbl.text = "傷害 ×2.0"
	dmg_lbl.position = Vector2(8, 22)
	dmg_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	dmg_lbl.add_theme_font_size_override("font_size", 12)
	c.add_child(dmg_lbl)

	# 凍結計時條背景
	var bar_bg := ColorRect.new()
	bar_bg.name = "BarBg"
	bar_bg.color = Color(0.1, 0.1, 0.2, 1.0)
	bar_bg.position = Vector2(8, 42)
	bar_bg.size = Vector2(184, 14)
	c.add_child(bar_bg)

	# 凍結計時條前景
	var bar_fg := ColorRect.new()
	bar_fg.name = "BarFg"
	bar_fg.color = COLOR_CYAN
	bar_fg.position = Vector2(8, 42)
	bar_fg.size = Vector2(184, 14)
	c.add_child(bar_fg)

	# 計時文字
	var time_lbl := Label.new()
	time_lbl.name = "TimeLabel"
	time_lbl.text = "8.0s"
	time_lbl.position = Vector2(8, 56)
	time_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	time_lbl.add_theme_font_size_override("font_size", 11)
	c.add_child(time_lbl)

	return c

func _update_freeze_bar() -> void:
	if not is_instance_valid(_freeze_indicator):
		return
	var bar_fg := _freeze_indicator.get_node_or_null("BarFg") as ColorRect
	var time_lbl := _freeze_indicator.get_node_or_null("TimeLabel") as Label
	if bar_fg:
		var pct := clamp(_freeze_remaining / float(_freeze_duration), 0.0, 1.0)
		bar_fg.size.x = 184.0 * pct
		# 顏色隨時間：青藍→紫→金
		if pct > 0.5:
			bar_fg.color = COLOR_CYAN
		elif pct > 0.25:
			bar_fg.color = COLOR_PURPLE
		else:
			bar_fg.color = COLOR_GOLD
	if time_lbl:
		time_lbl.text = "%.1fs" % _freeze_remaining

func _make_perfect_indicator() -> Control:
	var c := Control.new()
	c.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	c.position = Vector2(-220, 160)
	c.custom_minimum_size = Vector2(200, 60)

	var bg := ColorRect.new()
	bg.color = Color(0.2, 0.15, 0.0, 0.9)
	bg.size = Vector2(200, 60)
	c.add_child(bg)

	var title_lbl := Label.new()
	title_lbl.name = "TitleLabel"
	title_lbl.text = "⏸️ 時空完美！"
	title_lbl.position = Vector2(8, 4)
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	title_lbl.add_theme_font_size_override("font_size", 14)
	c.add_child(title_lbl)

	var mult_lbl := Label.new()
	mult_lbl.name = "MultLabel"
	mult_lbl.text = "全服 ×2.5"
	mult_lbl.position = Vector2(8, 24)
	mult_lbl.add_theme_color_override("font_color", COLOR_CYAN)
	mult_lbl.add_theme_font_size_override("font_size", 13)
	c.add_child(mult_lbl)

	var timer_lbl := Label.new()
	timer_lbl.name = "TimerLabel"
	timer_lbl.text = "6.0s"
	timer_lbl.position = Vector2(8, 42)
	timer_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	timer_lbl.add_theme_font_size_override("font_size", 12)
	c.add_child(timer_lbl)

	return c

func _update_perfect_timer() -> void:
	if not is_instance_valid(_perfect_indicator):
		return
	var timer_lbl := _perfect_indicator.get_node_or_null("TimerLabel") as Label
	if timer_lbl:
		timer_lbl.text = "%.1fs" % _perfect_remaining
	# 脈衝動畫
	var scale_val := 1.0 + 0.05 * sin(Time.get_ticks_msec() * 0.005)
	_perfect_indicator.scale = Vector2(scale_val, scale_val)

# ─── 特效工具 ──────────────────────────────────────────────────────────────

func _flash_screen(color: Color, times: int) -> void:
	if not is_instance_valid(_flash_overlay):
		return
	var tween := create_tween()
	for i in range(times):
		tween.tween_property(_flash_overlay, "color:a", 0.45, 0.06)
		tween.tween_property(_flash_overlay, "color:a", 0.0, 0.10)
	_flash_overlay.color = Color(color.r, color.g, color.b, 0.0)

func _spawn_float_text(text: String, color: Color, pos: Vector2) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 18)
	lbl.position = pos
	add_child(lbl)
	var tween := create_tween()
	tween.tween_property(lbl, "position:y", pos.y - 60, 1.2)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 1.2)
	tween.tween_callback(func(): if is_instance_valid(lbl): lbl.queue_free())

func _spawn_big_text(text: String, color: Color) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 36)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.set_anchors_preset(Control.PRESET_CENTER)
	lbl.position = Vector2(640 - 300, 280)
	lbl.custom_minimum_size = Vector2(600, 60)
	add_child(lbl)
	# 縮放動畫
	lbl.scale = Vector2(0.5, 0.5)
	var tween := create_tween()
	tween.tween_property(lbl, "scale", Vector2(1.1, 1.1), 0.15)
	tween.tween_property(lbl, "scale", Vector2(1.0, 1.0), 0.08)
	tween.tween_interval(1.5)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func(): if is_instance_valid(lbl): lbl.queue_free())
