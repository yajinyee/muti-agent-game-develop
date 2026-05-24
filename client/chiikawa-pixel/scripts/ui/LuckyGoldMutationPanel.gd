## LuckyGoldMutationPanel.gd — 幸運黃金突變魚 UI 面板（DAY-281）
## 黃金突變主題：#FFD700 金 + #FFA500 橙金 + #FF8C00 深橙 + #FFFACD 淡金 + #FF6B35 橙紅
## 業界原創「黃金突變+全場感染+突變連鎖」機制
##
## 事件類型：
##   mutation_start  — 黃金突變觸發（全服，PlayerID/PlayerName/TargetIIDs/TargetNames/KillMult/Duration）
##   mutation_kill   — 突變目標被擊破（全服，PlayerName/KilledIID/KillMult）
##   mutation_infect — 感染突變（全服，PlayerName/TargetIIDs/TargetNames/KillMult）
##   mutation_expire — 突變消失（全服）

extends CanvasLayer

const COLOR_GOLD       = Color(1.0,   0.843, 0.0)    # #FFD700 金
const COLOR_ORANGE_GOLD = Color(1.0,  0.647, 0.0)    # #FFA500 橙金
const COLOR_DEEP_ORANGE = Color(1.0,  0.549, 0.0)    # #FF8C00 深橙
const COLOR_LIGHT_GOLD = Color(1.0,   0.980, 0.804)  # #FFFACD 淡金
const COLOR_ORANGE_RED = Color(1.0,   0.420, 0.208)  # #FF6B35 橙紅
const COLOR_WHITE      = Color(1.0,   1.0,   1.0)

var _banner: Control = null
var _mutation_indicator: Control = null
var _mutation_count_label: Label = null

func _ready() -> void:
	layer = 54  # 比 LuckyRareChain（53）高一層

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"mutation_start":
			_on_mutation_start(payload)
		"mutation_kill":
			_on_mutation_kill(payload)
		"mutation_infect":
			_on_mutation_infect(payload)
		"mutation_expire":
			_on_mutation_expire()

# ── 黃金突變觸發（全服）────────────────────────────────────────────────────────

func _on_mutation_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var target_names: Array = payload.get("target_names", [])
	var kill_mult: float = payload.get("kill_mult", 3.0)
	var duration: int = payload.get("duration", 15)
	var count: int = target_names.size()

	# 金色三次強閃光
	_flash_screen(COLOR_GOLD, 3, 0.55)

	# 頂部橫幅
	var names_str := ", ".join(target_names) if target_names.size() > 0 else "???"
	_show_banner(
		"✨ 黃金突變！",
		"%s 觸發黃金突變！%d 個目標變成黃金魚！HP -50%%，擊破得 ×%.1f！%d 秒內快打！" % [player_name, count, kill_mult, duration],
		COLOR_GOLD
	)

	# 突變指示器（右上角）
	_show_mutation_indicator(count, kill_mult, duration)

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"✨ 黃金突變！%d 個目標！" % count,
		Vector2(vp_size / 2),
		COLOR_GOLD,
		40
	)
	_spawn_float_text(
		"HP -50%%，擊破得 ×%.1f！" % kill_mult,
		Vector2(vp_size.x * 0.5, vp_size.y * 0.6),
		COLOR_ORANGE_GOLD,
		22
	)

	# 5 秒後清除橫幅
	var timer := get_tree().create_timer(5.0)
	timer.timeout.connect(func():
		_clear_banner()
	)

# ── 突變目標被擊破（全服）────────────────────────────────────────────────────

func _on_mutation_kill(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var kill_mult: float = payload.get("kill_mult", 3.0)

	# 金色閃光
	_flash_screen(COLOR_GOLD, 1, 0.3)

	# 浮動文字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"✨ %s 擊破黃金魚！×%.1f！" % [player_name, kill_mult],
		Vector2(vp_size.x * 0.5, vp_size.y * 0.45),
		COLOR_GOLD,
		22
	)

# ── 感染突變（全服）──────────────────────────────────────────────────────────

func _on_mutation_infect(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var target_names: Array = payload.get("target_names", [])
	var kill_mult: float = payload.get("kill_mult", 2.0)

	# 橙金閃光
	_flash_screen(COLOR_ORANGE_GOLD, 1, 0.25)

	# 浮動文字
	var vp_size := get_viewport().get_visible_rect().size
	var name_str: String = target_names[0] if target_names.size() > 0 else "???"
	_show_mini_banner(
		"✨ 感染！%s 讓 %s 也突變了！擊破得 ×%.1f！" % [player_name, name_str, kill_mult],
		COLOR_ORANGE_GOLD
	)

# ── 突變消失（全服）──────────────────────────────────────────────────────────

func _on_mutation_expire() -> void:
	_clear_mutation_indicator()
	_clear_banner()

	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"✨ 黃金突變消失",
		Vector2(vp_size.x * 0.5, vp_size.y * 0.5),
		COLOR_DEEP_ORANGE,
		18
	)

# ── 內部 UI 工具 ──────────────────────────────────────────────────────────────

func _flash_screen(color: Color, times: int, alpha: float = 0.35) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var flash := ColorRect.new()
	flash.size = vp_size
	flash.color = Color(color.r, color.g, color.b, 0.0)
	add_child(flash)

	var tween := create_tween()
	for i in range(times):
		tween.tween_property(flash, "color:a", alpha, 0.07)
		tween.tween_property(flash, "color:a", 0.0, 0.10)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)

func _show_banner(title: String, subtitle: String, color: Color) -> void:
	_clear_banner()
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(0, -80)
	panel.size = Vector2(vp_size.x, 72)
	panel.modulate = Color(0.08, 0.06, 0.0, 0.92)
	add_child(panel)
	_banner = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = title
	title_lbl.add_theme_color_override("font_color", color)
	title_lbl.add_theme_font_size_override("font_size", 22)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var sub_lbl := Label.new()
	sub_lbl.text = subtitle
	sub_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	sub_lbl.add_theme_font_size_override("font_size", 13)
	sub_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(sub_lbl)

	var tween := create_tween()
	tween.tween_property(panel, "position:y", 0.0, 0.22).set_ease(Tween.EASE_OUT)

func _clear_banner() -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()
	_banner = null

func _show_mini_banner(text: String, color: Color) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var lbl := Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 16)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.position = Vector2(0, 4)
	lbl.size = Vector2(vp_size.x, 28)
	add_child(lbl)

	var tween := create_tween()
	tween.tween_interval(3.5)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)

func _show_mutation_indicator(count: int, kill_mult: float, duration: int) -> void:
	_clear_mutation_indicator()
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x - 150, 80)
	panel.size = Vector2(140, 90)
	panel.modulate = Color(0.08, 0.06, 0.0, 0.92)
	add_child(panel)
	_mutation_indicator = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "✨ 黃金突變"
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	_mutation_count_label = Label.new()
	_mutation_count_label.text = "%d 個黃金魚" % count
	_mutation_count_label.add_theme_color_override("font_color", COLOR_ORANGE_GOLD)
	_mutation_count_label.add_theme_font_size_override("font_size", 18)
	_mutation_count_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(_mutation_count_label)

	var mult_lbl := Label.new()
	mult_lbl.text = "擊破 ×%.1f" % kill_mult
	mult_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	mult_lbl.add_theme_font_size_override("font_size", 16)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	var dur_lbl := Label.new()
	dur_lbl.text = "⏱ %d 秒" % duration
	dur_lbl.add_theme_color_override("font_color", COLOR_LIGHT_GOLD)
	dur_lbl.add_theme_font_size_override("font_size", 12)
	dur_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(dur_lbl)

	# 金色脈衝動畫
	var tween := panel.create_tween().set_loops()
	tween.tween_property(panel, "modulate:a", 0.75, 0.5)
	tween.tween_property(panel, "modulate:a", 1.0, 0.5)

func _clear_mutation_indicator() -> void:
	if is_instance_valid(_mutation_indicator):
		_mutation_indicator.queue_free()
	_mutation_indicator = null
	_mutation_count_label = null

func _spawn_float_text(text: String, pos: Vector2, color: Color, font_size: int = 28) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", font_size)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.position = pos - Vector2(300, font_size * 0.5)
	lbl.size = Vector2(600, font_size * 2)
	add_child(lbl)

	var tween := create_tween()
	tween.tween_property(lbl, "position:y", pos.y - 60, 0.8).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.8).set_delay(0.5)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)
