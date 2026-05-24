## LuckyGoldenHurricanePanel.gd — 幸運黃金颶風魚 UI 面板（DAY-276）
## 黃金颶風主題：#FFD700 金 + #FFA500 橙 + #FF8C00 深橙 + #FFFACD 淡金
## 業界依據：Royal Fishing Jili 2026「AOE 旋風掃場」機制進化版
##
## 事件類型：
##   hurricane_start     — 颶風開始（全服，PlayerID/PlayerName/DurationSec/TargetCount）
##   hurricane_sweep     — 颶風掃過目標（全服，InstanceID/DefID/HPDamage/AccumMult）
##   hurricane_end       — 颶風結算（全服，PlayerID/PlayerName/SweptCount/FinalMult/TotalReward）
##   hurricane_broadcast — 全服廣播橫幅（PlayerName/SweptCount/FinalMult）

extends CanvasLayer

const COLOR_GOLD      = Color(1.0,   0.843, 0.0)    # #FFD700 金
const COLOR_ORANGE    = Color(1.0,   0.647, 0.0)    # #FFA500 橙
const COLOR_DEEP_ORG  = Color(1.0,   0.549, 0.0)    # #FF8C00 深橙
const COLOR_LIGHT_GLD = Color(1.0,   0.980, 0.804)  # #FFFACD 淡金
const COLOR_WHITE     = Color(1.0,   1.0,   1.0)
const COLOR_RED       = Color(1.0,   0.2,   0.2)

var _banner: Control = null
var _timer_bar: Control = null
var _mult_indicator: Control = null
var _mult_label: Label = null
var _timer_tween: Tween = null
var _sweep_count: int = 0

func _ready() -> void:
	layer = 49  # 比 LuckyLuckTotem（48）高一層

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"hurricane_start":
			_on_hurricane_start(payload)
		"hurricane_sweep":
			_on_hurricane_sweep(payload)
		"hurricane_end":
			_on_hurricane_end(payload)
		"hurricane_broadcast":
			_on_hurricane_broadcast(payload)

# ── 颶風開始 ──────────────────────────────────────────────────────────────────

func _on_hurricane_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var duration_sec: int = payload.get("duration_sec", 6)
	var target_count: int = payload.get("target_count", 0)
	_sweep_count = 0

	# 金色三次強閃光
	_flash_screen(COLOR_GOLD, 3, 0.5)

	# 頂部橫幅
	_show_banner(
		"🌪️ 黃金颶風！",
		"%s 觸發！場上 %d 個目標，螺旋掃場 %d 秒！" % [player_name, target_count, duration_sec],
		COLOR_GOLD
	)

	# 倍率指示器（右上角）
	_show_mult_indicator(1.0)

	# 右側豎向計時條
	_start_timer_bar(duration_sec)

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🌪️ 黃金颶風！螺旋掃場！",
		Vector2(vp_size / 2),
		COLOR_GOLD,
		44
	)

# ── 颶風掃過目標 ──────────────────────────────────────────────────────────────

func _on_hurricane_sweep(payload: Dictionary) -> void:
	var hp_damage: int = payload.get("hp_damage", 0)
	var accum_mult: float = payload.get("accum_mult", 1.0)
	_sweep_count += 1

	# 更新倍率指示器
	_update_mult_indicator(accum_mult)

	# 橙色閃光（每次掃過）
	_flash_screen(COLOR_ORANGE, 1, 0.18)

	# 浮動文字（掃過目標）
	var vp_size := get_viewport().get_visible_rect().size
	var pos := Vector2(
		vp_size.x * (0.3 + randf() * 0.4),
		vp_size.y * (0.3 + randf() * 0.4)
	)
	_spawn_float_text(
		"🌪️ HP -%d  ×%.1f" % [hp_damage, accum_mult],
		pos,
		COLOR_ORANGE,
		18
	)

# ── 颶風結算 ──────────────────────────────────────────────────────────────────

func _on_hurricane_end(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var swept_count: int = payload.get("swept_count", 0)
	var final_mult: float = payload.get("final_mult", 1.0)
	var total_reward: int = payload.get("total_reward", 0)

	# 清除計時條、橫幅、指示器
	_clear_timer_bar()
	_clear_banner()
	_clear_mult_indicator()

	# 根據最終倍率決定閃光強度
	if final_mult >= 5.0:
		# 金色三次強閃光（高倍率）
		_flash_screen(COLOR_GOLD, 3, 0.55)
		var vp_size := get_viewport().get_visible_rect().size
		_spawn_float_text(
			"🌪️ 颶風結算！×%.1f 累積倍率！" % final_mult,
			Vector2(vp_size / 2),
			COLOR_GOLD,
			42
		)
	else:
		# 橙色閃光（普通）
		_flash_screen(COLOR_ORANGE, 2, 0.3)

	# 結算彈窗
	_show_end_popup(player_name, swept_count, final_mult, total_reward)

# ── 全服廣播橫幅 ──────────────────────────────────────────────────────────────

func _on_hurricane_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var swept_count: int = payload.get("swept_count", 0)
	var final_mult: float = payload.get("final_mult", 1.0)

	var color := COLOR_GOLD if final_mult >= 5.0 else COLOR_ORANGE
	_show_mini_banner(
		"🌪️ %s 的黃金颶風結算！掃過 %d 個目標，累積倍率 ×%.1f！" % [player_name, swept_count, final_mult],
		color
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
	panel.modulate = Color(0.15, 0.1, 0.0, 0.92)
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
	sub_lbl.add_theme_font_size_override("font_size", 14)
	sub_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(sub_lbl)

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
	tween.tween_interval(3.5)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)

func _show_mult_indicator(mult: float) -> void:
	_clear_mult_indicator()
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x - 130, 80)
	panel.size = Vector2(120, 50)
	panel.modulate = Color(0.15, 0.1, 0.0, 0.9)
	add_child(panel)
	_mult_indicator = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	_mult_label = Label.new()
	_mult_label.text = "🌪️ ×%.1f" % mult
	_mult_label.add_theme_color_override("font_color", COLOR_GOLD)
	_mult_label.add_theme_font_size_override("font_size", 18)
	_mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(_mult_label)

	var sub_lbl := Label.new()
	sub_lbl.text = "累積倍率"
	sub_lbl.add_theme_color_override("font_color", COLOR_LIGHT_GLD)
	sub_lbl.add_theme_font_size_override("font_size", 11)
	sub_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(sub_lbl)

	# 脈衝動畫
	var tween := panel.create_tween().set_loops()
	tween.tween_property(panel, "modulate:a", 0.75, 0.4)
	tween.tween_property(panel, "modulate:a", 1.0, 0.4)

func _update_mult_indicator(mult: float) -> void:
	if not is_instance_valid(_mult_label):
		return
	_mult_label.text = "🌪️ ×%.1f" % mult
	# 顏色隨倍率變化：橙→深橙→金
	if mult >= 5.0:
		_mult_label.add_theme_color_override("font_color", COLOR_GOLD)
	elif mult >= 3.0:
		_mult_label.add_theme_color_override("font_color", COLOR_DEEP_ORG)
	else:
		_mult_label.add_theme_color_override("font_color", COLOR_ORANGE)

func _clear_mult_indicator() -> void:
	if is_instance_valid(_mult_indicator):
		_mult_indicator.queue_free()
	_mult_indicator = null
	_mult_label = null

func _start_timer_bar(duration_sec: int) -> void:
	_clear_timer_bar()
	var vp_size := get_viewport().get_visible_rect().size

	var bar_bg := ColorRect.new()
	bar_bg.position = Vector2(vp_size.x - 352, 80)  # x=-352 與其他計時條錯開
	bar_bg.size = Vector2(8, 200)
	bar_bg.color = Color(0.2, 0.1, 0.0, 0.7)
	add_child(bar_bg)
	_timer_bar = bar_bg

	var bar_fill := ColorRect.new()
	bar_fill.position = Vector2(0, 0)
	bar_fill.size = Vector2(8, 200)
	bar_fill.color = COLOR_GOLD
	bar_bg.add_child(bar_fill)

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

func _show_end_popup(player_name: String, swept_count: int, final_mult: float, total_reward: int) -> void:
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x + 10, vp_size.y * 0.3)
	panel.size = Vector2(270, 150)
	panel.modulate = Color(0.12, 0.08, 0.0, 0.95)
	add_child(panel)

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "🌪️ 黃金颶風結算"
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var trigger_lbl := Label.new()
	trigger_lbl.text = "觸發者：%s" % player_name
	trigger_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	trigger_lbl.add_theme_font_size_override("font_size", 13)
	trigger_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(trigger_lbl)

	var swept_lbl := Label.new()
	swept_lbl.text = "掃過目標：%d 個" % swept_count
	swept_lbl.add_theme_color_override("font_color", COLOR_ORANGE)
	swept_lbl.add_theme_font_size_override("font_size", 14)
	swept_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(swept_lbl)

	var mult_lbl := Label.new()
	mult_lbl.text = "累積倍率：×%.1f" % final_mult
	var mult_color := COLOR_GOLD if final_mult >= 5.0 else COLOR_DEEP_ORG
	mult_lbl.add_theme_color_override("font_color", mult_color)
	mult_lbl.add_theme_font_size_override("font_size", 18)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	var reward_lbl := Label.new()
	reward_lbl.text = "額外獎勵：+%d" % total_reward
	reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	reward_lbl.add_theme_font_size_override("font_size", 15)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	# 右側滑入動畫
	var tween := create_tween()
	tween.tween_property(panel, "position:x", vp_size.x - 280.0, 0.3).set_ease(Tween.EASE_OUT)
	tween.tween_interval(4.5)
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
