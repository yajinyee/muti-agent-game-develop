## LuckyRerollPanel.gd — 幸運倍率重擲魚 UI 面板（DAY-271）
## 橙金骰子主題：#FF8C00 橙 + #FFD700 金 + #FF4500 火橙 + #FFFFFF 白
##
## 事件類型：
##   reroll_start          — 重擲開始（個人）
##   reroll_broadcast      — 全服廣播
##   reroll_used           — 重擲被使用（個人）
##   reroll_result_broadcast — 全服廣播結果
##   reroll_expire         — session 超時（個人）

extends CanvasLayer

const COLOR_ORANGE  = Color(1.0,   0.549, 0.0)    # #FF8C00
const COLOR_GOLD    = Color(1.0,   0.843, 0.0)    # #FFD700
const COLOR_FIRE    = Color(1.0,   0.271, 0.0)    # #FF4500
const COLOR_WHITE   = Color(1.0,   1.0,   1.0)
const COLOR_LOW     = Color(0.8,   0.6,   0.2)    # 低倍率（棕）
const COLOR_MID     = Color(1.0,   0.7,   0.0)    # 中倍率（橙）
const COLOR_HIGH    = Color(1.0,   0.843, 0.0)    # 高倍率（金）
const COLOR_MAX     = Color(1.0,   0.4,   0.0)    # 最高倍率（火橙）

var _banner: Control = null
var _roll_display: Control = null
var _waiting_indicator: Control = null
var _waiting_tween: Tween = null

func _ready() -> void:
	layer = 44  # 比 MirrorDuel（43）高一層

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"reroll_start":
			_on_reroll_start(payload)
		"reroll_broadcast":
			_on_reroll_broadcast(payload)
		"reroll_used":
			_on_reroll_used(payload)
		"reroll_result_broadcast":
			_on_reroll_result_broadcast(payload)
		"reroll_expire":
			_on_reroll_expire(payload)

# ── 重擲開始 ──────────────────────────────────────────────────────────────────

func _on_reroll_start(payload: Dictionary) -> void:
	var rolls: Array = payload.get("rolls", [])
	var best_mult: float = payload.get("best_mult", 1.0)

	# 三次強閃光（橙色）
	_flash_screen(COLOR_ORANGE, 3)

	# 頂部橫幅
	_show_banner("🎲 倍率重擲！", "最高 ×%.1f 等待下一擊！" % best_mult, COLOR_ORANGE)

	# 顯示三次重擲結果
	_show_roll_display(rolls, best_mult)

	# 等待指示器（右側閃爍）
	_start_waiting_indicator(best_mult)

	# 浮動大字
	_spawn_float_text(
		"🎲 ×%.1f 就緒！" % best_mult,
		Vector2(get_viewport().get_visible_rect().size / 2),
		_mult_color(best_mult),
		46
	)

# ── 全服廣播 ──────────────────────────────────────────────────────────────────

func _on_reroll_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var best_mult: float = payload.get("best_mult", 1.0)
	_show_mini_banner(
		"🎲 %s 重擲 ×%.1f！等待命中！" % [player_name, best_mult],
		_mult_color(best_mult)
	)

# ── 重擲被使用 ────────────────────────────────────────────────────────────────

func _on_reroll_used(payload: Dictionary) -> void:
	var target_name: String = payload.get("target_name", "???")
	var reward: int = payload.get("reward", 0)
	var best_mult: float = payload.get("best_mult", 1.0)
	var rolls: Array = payload.get("rolls", [])

	# 清除等待指示器和橫幅
	_clear_waiting_indicator()
	_clear_banner()
	_clear_roll_display()

	# 閃光強度依倍率
	var flash_times := 3 if best_mult >= 3.0 else 2
	_flash_screen(_mult_color(best_mult), flash_times)

	# 大字
	_spawn_float_text(
		"🎲 ×%.1f 命中！+%d！" % [best_mult, reward],
		Vector2(get_viewport().get_visible_rect().size / 2),
		_mult_color(best_mult),
		50
	)

	# 結算彈窗
	_show_result_popup(rolls, best_mult, target_name, reward)

# ── 全服廣播結果 ──────────────────────────────────────────────────────────────

func _on_reroll_result_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var best_mult: float = payload.get("best_mult", 1.0)
	var reward: int = payload.get("reward", 0)
	_show_mini_banner(
		"🎲 %s 重擲命中！×%.1f +%d！" % [player_name, best_mult, reward],
		COLOR_GOLD
	)

# ── session 超時 ──────────────────────────────────────────────────────────────

func _on_reroll_expire(_payload: Dictionary) -> void:
	_clear_waiting_indicator()
	_clear_banner()
	_clear_roll_display()
	_spawn_float_text(
		"🎲 重擲超時...",
		Vector2(get_viewport().get_visible_rect().size.x * 0.5,
				get_viewport().get_visible_rect().size.y * 0.6),
		Color(0.6, 0.6, 0.6),
		22
	)

# ── 內部 UI 工具 ──────────────────────────────────────────────────────────────

func _mult_color(mult: float) -> Color:
	if mult >= 3.5:
		return COLOR_MAX
	elif mult >= 2.5:
		return COLOR_HIGH
	elif mult >= 1.8:
		return COLOR_MID
	else:
		return COLOR_LOW

func _flash_screen(color: Color, times: int, alpha: float = 0.45) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var flash := ColorRect.new()
	flash.size = vp_size
	flash.color = Color(color.r, color.g, color.b, 0.0)
	add_child(flash)

	var tween := create_tween()
	for i in range(times):
		tween.tween_property(flash, "color:a", alpha, 0.07)
		tween.tween_property(flash, "color:a", 0.0, 0.1)
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
	panel.modulate = Color(color.r, color.g, color.b, 0.92)
	add_child(panel)
	_banner = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var lbl1 := Label.new()
	lbl1.text = title
	lbl1.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl1.add_theme_font_size_override("font_size", 28)
	lbl1.add_theme_color_override("font_color", Color.WHITE)
	vbox.add_child(lbl1)

	var lbl2 := Label.new()
	lbl2.text = subtitle
	lbl2.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl2.add_theme_font_size_override("font_size", 18)
	lbl2.add_theme_color_override("font_color", Color(1, 1, 0.8))
	vbox.add_child(lbl2)

	var tween := create_tween()
	tween.tween_property(panel, "position:y", 0.0, 0.25)

func _show_mini_banner(text: String, color: Color) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var lbl := Label.new()
	lbl.text = text
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 18)
	lbl.add_theme_color_override("font_color", color)
	lbl.position = Vector2(0, 8)
	lbl.size = Vector2(vp_size.x, 28)
	add_child(lbl)

	var tween := create_tween()
	tween.tween_interval(3.0)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)

func _clear_banner() -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()
	_banner = null

func _show_roll_display(rolls: Array, best_mult: float) -> void:
	_clear_roll_display()
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x - 160, 80)
	panel.size = Vector2(150, 110)
	panel.modulate = Color(0.15, 0.1, 0.0, 0.92)
	add_child(panel)
	_roll_display = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "🎲 重擲結果"
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	vbox.add_child(title_lbl)

	for i in range(rolls.size()):
		var mult: float = rolls[i]
		var lbl := Label.new()
		var is_best := (mult == best_mult and mult > 1.0)
		lbl.text = "第%d擲：×%.1f%s" % [i+1, mult, " ★" if is_best else ""]
		lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		lbl.add_theme_font_size_override("font_size", 16)
		lbl.add_theme_color_override("font_color",
			_mult_color(mult) if mult > 1.0 else Color(0.5, 0.5, 0.5))
		vbox.add_child(lbl)

	var best_lbl := Label.new()
	best_lbl.text = "最高：×%.1f" % best_mult
	best_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	best_lbl.add_theme_font_size_override("font_size", 18)
	best_lbl.add_theme_color_override("font_color", _mult_color(best_mult))
	vbox.add_child(best_lbl)

func _clear_roll_display() -> void:
	if is_instance_valid(_roll_display):
		_roll_display.queue_free()
	_roll_display = null

func _start_waiting_indicator(best_mult: float) -> void:
	_clear_waiting_indicator()
	var vp_size := get_viewport().get_visible_rect().size

	var lbl := Label.new()
	lbl.text = "🎲 ×%.1f" % best_mult
	lbl.position = Vector2(vp_size.x - 100, vp_size.y * 0.5 - 20)
	lbl.size = Vector2(90, 40)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 22)
	lbl.add_theme_color_override("font_color", _mult_color(best_mult))
	add_child(lbl)
	_waiting_indicator = lbl

	# 脈衝動畫
	_waiting_tween = create_tween().set_loops()
	_waiting_tween.tween_property(lbl, "modulate:a", 0.3, 0.5)
	_waiting_tween.tween_property(lbl, "modulate:a", 1.0, 0.5)

func _clear_waiting_indicator() -> void:
	if _waiting_tween != null:
		_waiting_tween.kill()
		_waiting_tween = null
	if is_instance_valid(_waiting_indicator):
		_waiting_indicator.queue_free()
	_waiting_indicator = null

func _show_result_popup(rolls: Array, best_mult: float, target_name: String, reward: int) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var popup := PanelContainer.new()
	popup.size = Vector2(240, 150)
	popup.position = Vector2(vp_size.x + 10, vp_size.y * 0.3)
	popup.modulate = Color(0.1, 0.07, 0.0, 0.95)
	add_child(popup)

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	popup.add_child(vbox)

	var lbl_title := Label.new()
	lbl_title.text = "🎲 重擲命中！"
	lbl_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_title.add_theme_font_size_override("font_size", 20)
	lbl_title.add_theme_color_override("font_color", _mult_color(best_mult))
	vbox.add_child(lbl_title)

	var lbl_rolls := Label.new()
	var rolls_str := " / ".join(rolls.map(func(m): return "×%.1f" % m))
	lbl_rolls.text = rolls_str
	lbl_rolls.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_rolls.add_theme_font_size_override("font_size", 14)
	lbl_rolls.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	vbox.add_child(lbl_rolls)

	var lbl_best := Label.new()
	lbl_best.text = "最高 ×%.1f" % best_mult
	lbl_best.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_best.add_theme_font_size_override("font_size", 22)
	lbl_best.add_theme_color_override("font_color", _mult_color(best_mult))
	vbox.add_child(lbl_best)

	var lbl_reward := Label.new()
	lbl_reward.text = "%s +%d" % [target_name, reward]
	lbl_reward.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_reward.add_theme_font_size_override("font_size", 16)
	lbl_reward.add_theme_color_override("font_color", COLOR_GOLD)
	vbox.add_child(lbl_reward)

	var tween := create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 250.0, 0.3)
	tween.tween_interval(4.0)
	tween.tween_property(popup, "position:x", vp_size.x + 10.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(popup):
			popup.queue_free()
	)

func _spawn_float_text(text: String, pos: Vector2, color: Color, font_size: int = 28) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.position = pos - Vector2(200, 30)
	lbl.size = Vector2(400, 60)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", font_size)
	lbl.add_theme_color_override("font_color", color)
	add_child(lbl)

	var tween := create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 60.0, 1.2)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 1.2)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)
