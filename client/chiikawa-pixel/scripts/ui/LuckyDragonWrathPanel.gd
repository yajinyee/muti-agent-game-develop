## LuckyDragonWrathPanel.gd — 幸運龍怒隕石魚 UI 面板（DAY-284）
## 龍怒主題：火橙（#FF4500）+ 金（#FFD700）+ 深紅（#8B0000）+ 白（#FFFFFF）
## 業界依據：Royal Fishing Jili「Dragon Wrath meteors」機制（2026 最熱門）
##
## 事件類型：
##   wrath_start        — 龍怒隕石開始（全服，PlayerID/PlayerName/MeteorCount）
##   wrath_meteor       — 單顆隕石墜落（全服，PlayerID/MeteorX/MeteorY/HitTargets/AccumMult/MeteorIdx/TotalCount）
##   wrath_end          — 龍怒隕石結算（全服，PlayerID/PlayerName/MeteorCount/TotalHit/AccumMult/Reward/IsPerfect）
##   wrath_perfect      — 龍怒完美全服加成（全服，PlayerName/PerfectMult/Duration）
##   wrath_perfect_end  — 龍怒完美結束（全服）

extends CanvasLayer

const COLOR_FIRE_ORANGE = Color(1.0,   0.271, 0.0)    # #FF4500 火橙
const COLOR_GOLD        = Color(1.0,   0.843, 0.0)    # #FFD700 金
const COLOR_DEEP_RED    = Color(0.545, 0.0,   0.0)    # #8B0000 深紅
const COLOR_WHITE       = Color(1.0,   1.0,   1.0)    # 白
const COLOR_ORANGE      = Color(1.0,   0.647, 0.0)    # #FFA500 橙

var _banner: Control = null
var _meteor_indicator: Control = null
var _meteor_count_label: Label = null
var _accum_mult_label: Label = null
var _perfect_indicator: Control = null
var _perfect_timer_label: Label = null

func _ready() -> void:
	layer = 57  # 比 LuckyFourSymbols（56）高一層

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"wrath_start":
			_on_wrath_start(payload)
		"wrath_meteor":
			_on_wrath_meteor(payload)
		"wrath_end":
			_on_wrath_end(payload)
		"wrath_perfect":
			_on_wrath_perfect(payload)
		"wrath_perfect_end":
			_on_wrath_perfect_end()

# ── 龍怒隕石開始 ────────────────────────────────────────────────────────────────

func _on_wrath_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var meteor_count: int = payload.get("meteor_count", 4)

	# 火橙三次強閃光
	_flash_screen(COLOR_FIRE_ORANGE, 3, 0.6)

	# 頂部橫幅
	_show_banner(
		"🐉🔥 龍怒隕石！",
		"%s 召喚 %d 顆龍怒隕石！" % [player_name, meteor_count],
		COLOR_FIRE_ORANGE
	)

	# 右上角隕石指示器
	_show_meteor_indicator(meteor_count)

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🐉🔥 龍怒隕石！",
		Vector2(vp_size.x * 0.5, vp_size.y * 0.3),
		COLOR_FIRE_ORANGE,
		40
	)
	_spawn_float_text(
		"%d 顆隕石即將墜落！" % meteor_count,
		Vector2(vp_size.x * 0.5, vp_size.y * 0.45),
		COLOR_GOLD,
		22
	)

# ── 單顆隕石墜落 ────────────────────────────────────────────────────────────────

func _on_wrath_meteor(payload: Dictionary) -> void:
	var hit_targets: int = payload.get("hit_targets", 0)
	var accum_mult: float = payload.get("accum_mult", 1.0)
	var meteor_idx: int = payload.get("meteor_idx", 1)
	var total_count: int = payload.get("total_count", 4)
	var meteor_x: float = payload.get("meteor_x", 512.0)
	var meteor_y: float = payload.get("meteor_y", 300.0)

	# 依命中數決定閃光顏色
	var flash_color := COLOR_FIRE_ORANGE if hit_targets > 0 else COLOR_DEEP_RED
	_flash_screen(flash_color, 1, 0.35)

	# 隕石爆炸特效（在墜落位置）
	_spawn_meteor_explosion(Vector2(meteor_x, meteor_y), hit_targets)

	# 浮動文字
	var vp_size := get_viewport().get_visible_rect().size
	if hit_targets > 0:
		_spawn_float_text(
			"💥 命中 %d 個！×%.1f" % [hit_targets, accum_mult],
			Vector2(meteor_x, meteor_y - 30),
			COLOR_GOLD,
			18
		)
	else:
		_spawn_float_text(
			"💨 空砸！",
			Vector2(meteor_x, meteor_y - 30),
			COLOR_DEEP_RED,
			16
		)

	# 更新指示器
	_update_meteor_indicator(meteor_idx, total_count, accum_mult)

# ── 龍怒隕石結算 ────────────────────────────────────────────────────────────────

func _on_wrath_end(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var total_hit: int = payload.get("total_hit", 0)
	var accum_mult: float = payload.get("accum_mult", 1.0)
	var reward: int = payload.get("reward", 0)
	var is_perfect: bool = payload.get("is_perfect", false)

	# 清除隕石指示器
	if is_instance_valid(_meteor_indicator):
		_meteor_indicator.queue_free()
		_meteor_indicator = null

	# 清除橫幅
	if is_instance_valid(_banner):
		_banner.queue_free()
		_banner = null

	# 龍怒完美由 wrath_perfect 事件處理，這裡只處理普通結算
	if not is_perfect and accum_mult >= 3.0:
		# 金色三次強閃光
		_flash_screen(COLOR_GOLD, 3, 0.55)

		# 結算彈窗
		_show_settle_popup(player_name, total_hit, accum_mult, reward, false)

# ── 龍怒完美全服加成 ────────────────────────────────────────────────────────────

func _on_wrath_perfect(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var perfect_mult: float = payload.get("perfect_mult", 2.5)
	var duration: int = payload.get("duration", 6)

	# 全螢幕三次強閃光（火橙→金→深紅）
	_flash_screen(COLOR_FIRE_ORANGE, 1, 0.75)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_GOLD, 1, 0.85)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_DEEP_RED, 1, 0.7)

	# 全螢幕大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🐉🔥🐉 龍怒完美！",
		Vector2(vp_size.x * 0.5, vp_size.y * 0.22),
		COLOR_FIRE_ORANGE,
		52
	)
	_spawn_float_text(
		"全服 ×%.1f 加成 %d 秒！" % [perfect_mult, duration],
		Vector2(vp_size.x * 0.5, vp_size.y * 0.38),
		COLOR_GOLD,
		28
	)
	_spawn_float_text(
		"%s 龍怒完美！" % player_name,
		Vector2(vp_size.x * 0.5, vp_size.y * 0.52),
		COLOR_WHITE,
		18
	)

	# 右側完美加成指示器
	_show_perfect_indicator(perfect_mult, duration)

	# 結算彈窗
	_show_settle_popup(player_name, 0, perfect_mult, 0, true)

# ── 龍怒完美結束 ────────────────────────────────────────────────────────────────

func _on_wrath_perfect_end() -> void:
	if is_instance_valid(_perfect_indicator):
		var tween := create_tween()
		tween.tween_property(_perfect_indicator, "modulate:a", 0.0, 0.4)
		tween.tween_callback(_perfect_indicator.queue_free)
		_perfect_indicator = null

# ── 輔助方法 ────────────────────────────────────────────────────────────────────

func _show_meteor_indicator(meteor_count: int) -> void:
	if is_instance_valid(_meteor_indicator):
		_meteor_indicator.queue_free()

	var vp_size := get_viewport().get_visible_rect().size
	var indicator := PanelContainer.new()
	indicator.position = Vector2(vp_size.x - 148, 80)
	indicator.size = Vector2(140, 72)
	indicator.mouse_filter = Control.MOUSE_FILTER_IGNORE

	var style := StyleBoxFlat.new()
	style.bg_color = Color(0, 0, 0, 0.85)
	style.border_color = COLOR_FIRE_ORANGE
	style.set_border_width_all(2)
	style.set_corner_radius_all(5)
	indicator.add_theme_stylebox_override("panel", style)

	var vbox := VBoxContainer.new()
	vbox.mouse_filter = Control.MOUSE_FILTER_IGNORE
	indicator.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "🔥 龍怒隕石"
	title_lbl.add_theme_color_override("font_color", COLOR_FIRE_ORANGE)
	title_lbl.add_theme_font_size_override("font_size", 11)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	vbox.add_child(title_lbl)

	var count_lbl := Label.new()
	count_lbl.name = "MeteorCountLabel"
	count_lbl.text = "0 / %d" % meteor_count
	count_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	count_lbl.add_theme_font_size_override("font_size", 14)
	count_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	count_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	vbox.add_child(count_lbl)

	var mult_lbl := Label.new()
	mult_lbl.name = "AccumMultLabel"
	mult_lbl.text = "×1.0"
	mult_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	mult_lbl.add_theme_font_size_override("font_size", 13)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mult_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	vbox.add_child(mult_lbl)

	add_child(indicator)
	_meteor_indicator = indicator
	_meteor_count_label = count_lbl
	_accum_mult_label = mult_lbl

func _update_meteor_indicator(meteor_idx: int, total_count: int, accum_mult: float) -> void:
	if is_instance_valid(_meteor_count_label):
		_meteor_count_label.text = "%d / %d" % [meteor_idx, total_count]
	if is_instance_valid(_accum_mult_label):
		_accum_mult_label.text = "×%.1f" % accum_mult
		# 顏色隨倍率變化
		var mult_color := COLOR_GOLD
		if accum_mult >= 5.0:
			mult_color = COLOR_FIRE_ORANGE
		elif accum_mult >= 3.0:
			mult_color = COLOR_ORANGE
		_accum_mult_label.add_theme_color_override("font_color", mult_color)

func _show_perfect_indicator(perfect_mult: float, duration: int) -> void:
	if is_instance_valid(_perfect_indicator):
		_perfect_indicator.queue_free()

	var vp_size := get_viewport().get_visible_rect().size
	var indicator := PanelContainer.new()
	indicator.position = Vector2(vp_size.x - 148, 160)
	indicator.size = Vector2(140, 72)
	indicator.mouse_filter = Control.MOUSE_FILTER_IGNORE

	var style := StyleBoxFlat.new()
	style.bg_color = Color(0, 0, 0, 0.9)
	style.border_color = COLOR_FIRE_ORANGE
	style.set_border_width_all(3)
	style.set_corner_radius_all(6)
	indicator.add_theme_stylebox_override("panel", style)

	var vbox := VBoxContainer.new()
	vbox.mouse_filter = Control.MOUSE_FILTER_IGNORE
	indicator.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "🐉 龍怒完美！"
	title_lbl.add_theme_color_override("font_color", COLOR_FIRE_ORANGE)
	title_lbl.add_theme_font_size_override("font_size", 11)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	vbox.add_child(title_lbl)

	var mult_lbl := Label.new()
	mult_lbl.text = "全服 ×%.1f" % perfect_mult
	mult_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	mult_lbl.add_theme_font_size_override("font_size", 14)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mult_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	vbox.add_child(mult_lbl)

	var timer_lbl := Label.new()
	timer_lbl.name = "PerfectTimerLabel"
	timer_lbl.text = "%d 秒" % duration
	timer_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	timer_lbl.add_theme_font_size_override("font_size", 12)
	timer_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	timer_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	vbox.add_child(timer_lbl)

	add_child(indicator)
	_perfect_indicator = indicator
	_perfect_timer_label = timer_lbl

	# 脈衝動畫
	var tween := create_tween().set_loops()
	tween.tween_property(indicator, "modulate:a", 0.7, 0.4)
	tween.tween_property(indicator, "modulate:a", 1.0, 0.4)

	# 倒數計時
	var remaining := duration
	var countdown_timer := get_tree().create_timer(1.0)
	countdown_timer.timeout.connect(func():
		remaining -= 1
		if is_instance_valid(_perfect_timer_label):
			_perfect_timer_label.text = "%d 秒" % max(0, remaining)
	)

func _spawn_meteor_explosion(pos: Vector2, hit_count: int) -> void:
	# 爆炸圓圈
	var circle := ColorRect.new()
	var size := 80.0 if hit_count > 0 else 40.0
	circle.size = Vector2(size, size)
	circle.position = pos - Vector2(size * 0.5, size * 0.5)
	circle.color = Color(1.0, 0.271, 0.0, 0.7) if hit_count > 0 else Color(0.545, 0.0, 0.0, 0.5)
	circle.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(circle)

	var tween := create_tween()
	tween.tween_property(circle, "scale", Vector2(2.0, 2.0), 0.3)
	tween.parallel().tween_property(circle, "modulate:a", 0.0, 0.3)
	tween.tween_callback(circle.queue_free)

func _flash_screen(color: Color, times: int, alpha: float) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var flash := ColorRect.new()
	flash.size = vp_size
	flash.color = Color(color.r, color.g, color.b, alpha)
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)

	var tween := create_tween()
	for i in range(times):
		tween.tween_property(flash, "modulate:a", 1.0, 0.0)
		tween.tween_property(flash, "modulate:a", 0.0, 0.14)
		if i < times - 1:
			tween.tween_interval(0.07)
	tween.tween_callback(flash.queue_free)

func _show_banner(title: String, message: String, color: Color) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()

	var vp_size := get_viewport().get_visible_rect().size
	var banner := PanelContainer.new()
	banner.position = Vector2(vp_size.x * 0.1, 8)
	banner.size = Vector2(vp_size.x * 0.8, 56)
	banner.mouse_filter = Control.MOUSE_FILTER_IGNORE

	var style := StyleBoxFlat.new()
	style.bg_color = Color(0, 0, 0, 0.85)
	style.border_color = color
	style.set_border_width_all(3)
	style.set_corner_radius_all(6)
	banner.add_theme_stylebox_override("panel", style)

	var vbox := VBoxContainer.new()
	vbox.mouse_filter = Control.MOUSE_FILTER_IGNORE
	banner.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = title
	title_lbl.add_theme_color_override("font_color", color)
	title_lbl.add_theme_font_size_override("font_size", 15)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	vbox.add_child(title_lbl)

	var msg_lbl := Label.new()
	msg_lbl.text = message
	msg_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	msg_lbl.add_theme_font_size_override("font_size", 11)
	msg_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	msg_lbl.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	msg_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	vbox.add_child(msg_lbl)

	add_child(banner)
	_banner = banner

	var tween := create_tween()
	tween.tween_property(banner, "modulate:a", 1.0, 0.25)

	var timer := get_tree().create_timer(5.0)
	timer.timeout.connect(func():
		if is_instance_valid(_banner) and _banner == banner:
			var fade := create_tween()
			fade.tween_property(banner, "modulate:a", 0.0, 0.3)
			fade.tween_callback(banner.queue_free)
			_banner = null
	)

func _show_settle_popup(player_name: String, total_hit: int, accum_mult: float, reward: int, is_perfect: bool) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var popup := PanelContainer.new()
	popup.position = Vector2(vp_size.x + 10, vp_size.y * 0.3)
	popup.size = Vector2(220, 140)
	popup.mouse_filter = Control.MOUSE_FILTER_IGNORE

	var border_color := COLOR_FIRE_ORANGE if is_perfect else COLOR_GOLD
	var style := StyleBoxFlat.new()
	style.bg_color = Color(0, 0, 0, 0.92)
	style.border_color = border_color
	style.set_border_width_all(3)
	style.set_corner_radius_all(8)
	popup.add_theme_stylebox_override("panel", style)

	var vbox := VBoxContainer.new()
	vbox.mouse_filter = Control.MOUSE_FILTER_IGNORE
	popup.add_child(vbox)

	var lines: Array
	if is_perfect:
		lines = [
			["🐉🔥 龍怒完美！", COLOR_FIRE_ORANGE, 16],
			[player_name, COLOR_WHITE, 12],
			["全服 ×%.1f 加成！" % accum_mult, COLOR_GOLD, 14],
			["6 秒全服加成！", COLOR_ORANGE, 12],
		]
	else:
		lines = [
			["🔥 龍怒結算", COLOR_FIRE_ORANGE, 16],
			[player_name, COLOR_WHITE, 12],
			["命中 %d 個目標" % total_hit, COLOR_ORANGE, 12],
			["累積 ×%.1f" % accum_mult, COLOR_GOLD, 14],
			["獎勵：+%d 金幣" % reward, COLOR_GOLD, 13],
		]

	for line in lines:
		var lbl := Label.new()
		lbl.text = line[0]
		lbl.add_theme_color_override("font_color", line[1])
		lbl.add_theme_font_size_override("font_size", line[2])
		lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
		vbox.add_child(lbl)

	add_child(popup)

	var tween := create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 230, 0.35)
	tween.tween_interval(4.5)
	tween.tween_property(popup, "position:x", vp_size.x + 10, 0.3)
	tween.tween_callback(popup.queue_free)

func _spawn_float_text(text: String, pos: Vector2, color: Color, font_size: int = 20) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", font_size)
	lbl.position = pos - Vector2(250, 0)
	lbl.size = Vector2(500, 70)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(lbl)

	var tween := create_tween()
	tween.tween_property(lbl, "position:y", pos.y - 70, 1.4)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 1.4)
	tween.tween_callback(lbl.queue_free)
