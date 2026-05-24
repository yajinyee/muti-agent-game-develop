## LuckyFortuneProphecyPanel.gd — 幸運命運預言魚 UI 面板（DAY-274）
## 紫金預言主題：#9B59B6 紫 + #FFD700 金 + #E8DAEF 淡紫白 + #FF8C00 橙
## 業界依據：Lucky Fish by AbraCadabra（2026-05-16）crash mechanic 進化版
##
## 事件類型：
##   prophecy_start              — 預言開始（個人，PredictedMult/ExpiresIn）
##   prophecy_broadcast          — 全服廣播（PlayerName/PredictedMult）
##   prophecy_fulfilled          — 預言成真（ActualMult/ResultMult/BonusReward/TargetName）
##   prophecy_fulfilled_broadcast — 預言成真全服廣播（PlayerName/ActualMult/ResultMult/BonusReward）
##   prophecy_failed             — 預言落空（ActualMult/ResultMult/BonusReward/TargetName）
##   prophecy_expire             — 預言超時（PredictedMult）

extends CanvasLayer

const COLOR_PURPLE  = Color(0.608, 0.349, 0.714)  # #9B59B6 紫
const COLOR_GOLD    = Color(1.0,   0.843, 0.0)    # #FFD700 金
const COLOR_ORANGE  = Color(1.0,   0.549, 0.0)    # #FF8C00 橙
const COLOR_LAVENDER = Color(0.91, 0.855, 0.937)  # #E8DAEF 淡紫白
const COLOR_WHITE   = Color(1.0,   1.0,   1.0)
const COLOR_GREEN   = Color(0.0,   1.0,   0.533)  # #00FF88 翠綠（成真）
const COLOR_RED     = Color(1.0,   0.3,   0.3)    # 落空紅

var _banner: Control = null
var _timer_bar: Control = null
var _predict_indicator: Control = null
var _timer_tween: Tween = null

func _ready() -> void:
	layer = 47  # 比 LuckyResonanceWave（46）高一層

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"prophecy_start":
			_on_prophecy_start(payload)
		"prophecy_broadcast":
			_on_prophecy_broadcast(payload)
		"prophecy_fulfilled":
			_on_prophecy_fulfilled(payload)
		"prophecy_fulfilled_broadcast":
			_on_prophecy_fulfilled_broadcast(payload)
		"prophecy_failed":
			_on_prophecy_failed(payload)
		"prophecy_expire":
			_on_prophecy_expire(payload)

# ── 預言開始 ──────────────────────────────────────────────────────────────────

func _on_prophecy_start(payload: Dictionary) -> void:
	var predicted_mult: float = payload.get("predicted_mult", 2.0)
	var expires_in: int = payload.get("expires_in", 20)

	# 紫色三次強閃光
	_flash_screen(COLOR_PURPLE, 3, 0.45)

	# 頂部橫幅
	_show_banner(
		"🔮 命運預言！",
		"門檻 ×%.1f — 20 秒內能成真嗎？" % predicted_mult,
		COLOR_PURPLE
	)

	# 預言倍率指示器（右上角）
	_show_predict_indicator(predicted_mult)

	# 右側豎向計時條
	_start_timer_bar(expires_in)

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🔮 命運預言！門檻 ×%.1f" % predicted_mult,
		Vector2(vp_size / 2),
		COLOR_PURPLE,
		44
	)

# ── 全服廣播 ──────────────────────────────────────────────────────────────────

func _on_prophecy_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var predicted_mult: float = payload.get("predicted_mult", 2.0)
	_show_mini_banner(
		"🔮 %s 觸發命運預言！門檻 ×%.1f" % [player_name, predicted_mult],
		COLOR_PURPLE
	)

# ── 預言成真 ──────────────────────────────────────────────────────────────────

func _on_prophecy_fulfilled(payload: Dictionary) -> void:
	var actual_mult: float = payload.get("actual_mult", 2.0)
	var result_mult: float = payload.get("result_mult", 3.0)
	var bonus_reward: int = payload.get("bonus_reward", 0)
	var target_name: String = payload.get("target_name", "???")

	# 清除計時條和橫幅
	_clear_timer_bar()
	_clear_banner()
	_clear_predict_indicator()

	# 金色三次強閃光（成真！）
	_flash_screen(COLOR_GOLD, 3, 0.55)

	# 大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🔮 預言成真！×%.1f 命中！" % actual_mult,
		Vector2(vp_size / 2),
		COLOR_GOLD,
		52
	)

	# 結算彈窗
	_show_result_popup(true, actual_mult, result_mult, bonus_reward, target_name)

# ── 預言成真全服廣播 ──────────────────────────────────────────────────────────

func _on_prophecy_fulfilled_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var actual_mult: float = payload.get("actual_mult", 2.0)
	var bonus_reward: int = payload.get("bonus_reward", 0)
	_show_mini_banner(
		"🔮 %s 預言成真！×%.1f 命中！+%d！" % [player_name, actual_mult, bonus_reward],
		COLOR_GOLD
	)

# ── 預言落空 ──────────────────────────────────────────────────────────────────

func _on_prophecy_failed(payload: Dictionary) -> void:
	var actual_mult: float = payload.get("actual_mult", 1.0)
	var result_mult: float = payload.get("result_mult", 1.2)
	var bonus_reward: int = payload.get("bonus_reward", 0)
	var target_name: String = payload.get("target_name", "???")

	# 清除計時條和橫幅
	_clear_timer_bar()
	_clear_banner()
	_clear_predict_indicator()

	# 橙色閃光（落空安慰）
	_flash_screen(COLOR_ORANGE, 1, 0.3)

	# 浮動文字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🔮 預言落空... ×%.1f 安慰獎 +%d" % [result_mult, bonus_reward],
		Vector2(vp_size / 2),
		COLOR_ORANGE,
		32
	)

	# 結算彈窗
	_show_result_popup(false, actual_mult, result_mult, bonus_reward, target_name)

# ── 預言超時 ──────────────────────────────────────────────────────────────────

func _on_prophecy_expire(payload: Dictionary) -> void:
	var predicted_mult: float = payload.get("predicted_mult", 2.0)

	_clear_timer_bar()
	_clear_banner()
	_clear_predict_indicator()

	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🔮 預言超時... 門檻 ×%.1f 未達成" % predicted_mult,
		Vector2(vp_size.x * 0.5, vp_size.y * 0.5),
		Color(0.5, 0.4, 0.6),
		22
	)

# ── 內部 UI 工具 ──────────────────────────────────────────────────────────────

func _flash_screen(color: Color, times: int, alpha: float = 0.4) -> void:
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
	panel.modulate = Color(0.1, 0.05, 0.15, 0.92)
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
	sub_lbl.add_theme_color_override("font_color", COLOR_LAVENDER)
	sub_lbl.add_theme_font_size_override("font_size", 14)
	sub_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(sub_lbl)

	# 滑入動畫
	var tween := create_tween()
	tween.tween_property(panel, "position:y", 0.0, 0.25).set_ease(Tween.EASE_OUT)

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
	tween.tween_interval(2.5)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)

func _show_predict_indicator(predicted_mult: float) -> void:
	_clear_predict_indicator()
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x - 130, 80)
	panel.size = Vector2(120, 52)
	panel.modulate = Color(0.1, 0.05, 0.2, 0.9)
	add_child(panel)
	_predict_indicator = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var label1 := Label.new()
	label1.text = "🔮 預言門檻"
	label1.add_theme_color_override("font_color", COLOR_LAVENDER)
	label1.add_theme_font_size_override("font_size", 12)
	label1.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(label1)

	var label2 := Label.new()
	label2.text = "×%.1f" % predicted_mult
	# 顏色依門檻高低：低=紫，中=橙，高=金
	var c := COLOR_PURPLE
	if predicted_mult >= 6.0:
		c = COLOR_GOLD
	elif predicted_mult >= 4.0:
		c = COLOR_ORANGE
	label2.add_theme_color_override("font_color", c)
	label2.add_theme_font_size_override("font_size", 22)
	label2.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(label2)

	# 脈衝動畫
	var tween := panel.create_tween().set_loops()
	tween.tween_property(panel, "modulate:a", 0.7, 0.6)
	tween.tween_property(panel, "modulate:a", 1.0, 0.6)

func _clear_predict_indicator() -> void:
	if is_instance_valid(_predict_indicator):
		_predict_indicator.queue_free()
	_predict_indicator = null

func _start_timer_bar(duration_sec: int) -> void:
	_clear_timer_bar()
	var vp_size := get_viewport().get_visible_rect().size

	var bar_bg := ColorRect.new()
	bar_bg.position = Vector2(vp_size.x - 324, 80)  # x=-324 與其他計時條錯開
	bar_bg.size = Vector2(8, 200)
	bar_bg.color = Color(0.2, 0.1, 0.3, 0.7)
	add_child(bar_bg)
	_timer_bar = bar_bg

	var bar_fill := ColorRect.new()
	bar_fill.position = Vector2(0, 0)
	bar_fill.size = Vector2(8, 200)
	bar_fill.color = COLOR_PURPLE
	bar_bg.add_child(bar_fill)

	# 計時條縮短動畫
	if _timer_tween != null and _timer_tween.is_valid():
		_timer_tween.kill()
	_timer_tween = create_tween()
	_timer_tween.tween_property(bar_fill, "size:y", 0.0, float(duration_sec))
	_timer_tween.tween_callback(func():
		_clear_timer_bar()
	)

func _clear_timer_bar() -> void:
	if _timer_tween != null and _timer_tween.is_valid():
		_timer_tween.kill()
	_timer_tween = null
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
	_timer_bar = null

func _show_result_popup(fulfilled: bool, actual_mult: float, result_mult: float, bonus_reward: int, target_name: String) -> void:
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x + 10, vp_size.y * 0.3)
	panel.size = Vector2(260, 160)
	if fulfilled:
		panel.modulate = Color(0.15, 0.12, 0.0, 0.95)
	else:
		panel.modulate = Color(0.1, 0.05, 0.15, 0.95)
	add_child(panel)

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	if fulfilled:
		title_lbl.text = "🔮 預言成真！"
		title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	else:
		title_lbl.text = "🔮 預言落空"
		title_lbl.add_theme_color_override("font_color", COLOR_ORANGE)
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var target_lbl := Label.new()
	target_lbl.text = "目標：%s" % target_name
	target_lbl.add_theme_color_override("font_color", COLOR_LAVENDER)
	target_lbl.add_theme_font_size_override("font_size", 13)
	target_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(target_lbl)

	var mult_lbl := Label.new()
	mult_lbl.text = "實際倍率：×%.1f" % actual_mult
	mult_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	mult_lbl.add_theme_font_size_override("font_size", 14)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	var result_lbl := Label.new()
	result_lbl.text = "加成：×%.1f" % result_mult
	if fulfilled:
		result_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	else:
		result_lbl.add_theme_color_override("font_color", COLOR_ORANGE)
	result_lbl.add_theme_font_size_override("font_size", 16)
	result_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(result_lbl)

	var reward_lbl := Label.new()
	reward_lbl.text = "額外獎勵：+%d" % bonus_reward
	reward_lbl.add_theme_color_override("font_color", COLOR_GREEN)
	reward_lbl.add_theme_font_size_override("font_size", 15)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	# 右側滑入動畫
	var tween := create_tween()
	tween.tween_property(panel, "position:x", vp_size.x - 270.0, 0.3).set_ease(Tween.EASE_OUT)
	tween.tween_interval(3.5)
	tween.tween_property(panel, "position:x", vp_size.x + 10.0, 0.3).set_ease(Tween.EASE_IN)
	tween.tween_callback(func():
		if is_instance_valid(panel):
			panel.queue_free()
	)

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
