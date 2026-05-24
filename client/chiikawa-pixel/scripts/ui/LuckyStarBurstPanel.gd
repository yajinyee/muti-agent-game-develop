## LuckyStarBurstPanel.gd — 幸運星爆魚 UI 面板（DAY-282）
## 星爆主題：#FFD700 金 + #00BFFF 天藍 + #FF69B4 粉紅 + #7FFF00 草綠 + #FFFFFF 白
## 業界原創「星爆連鎖+全場星雨+倍率爆炸」機制
##
## 事件類型：
##   burst_start     — 星爆觸發（全服，PlayerID/PlayerName/BurstCount/AccumMult/Duration）
##   burst_explode   — 星爆點爆炸（全服，BurstIndex/HitCount/AccumMult）
##   burst_resonance — 星爆共鳴觸發（全服，GlobalMult/GlobalDuration）
##   burst_end       — 星爆結算（全服，PlayerName/TotalBursts/TotalHits/FinalMult/TotalReward）

extends CanvasLayer

const COLOR_GOLD       = Color(1.0,   0.843, 0.0)    # #FFD700 金
const COLOR_SKY_BLUE   = Color(0.0,   0.749, 1.0)    # #00BFFF 天藍
const COLOR_PINK       = Color(1.0,   0.412, 0.706)  # #FF69B4 粉紅
const COLOR_LIME       = Color(0.498, 1.0,   0.0)    # #7FFF00 草綠
const COLOR_WHITE      = Color(1.0,   1.0,   1.0)
const COLOR_ORANGE     = Color(1.0,   0.647, 0.0)    # #FFA500 橙

var _banner: Control = null
var _burst_indicator: Control = null
var _burst_count_label: Label = null
var _accum_mult_label: Label = null
var _resonance_indicator: Control = null
var _resonance_timer_label: Label = null
var _resonance_tween: Tween = null

func _ready() -> void:
	layer = 55  # 比 LuckyGoldMutation（54）高一層

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"burst_start":
			_on_burst_start(payload)
		"burst_explode":
			_on_burst_explode(payload)
		"burst_resonance":
			_on_burst_resonance(payload)
		"burst_end":
			_on_burst_end(payload)

# ── 星爆觸發（全服）────────────────────────────────────────────────────────────

func _on_burst_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var burst_count: int = payload.get("burst_count", 5)
	var duration: int = payload.get("duration", 3)

	# 天藍+金色三次強閃光
	_flash_screen(COLOR_SKY_BLUE, 2, 0.45)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_GOLD, 1, 0.55)

	# 頂部橫幅
	_show_banner(
		"⭐ 星爆！",
		"%s 觸發星爆！%d 個星爆點即將爆炸！全場 HP -35%%！" % [player_name, burst_count],
		COLOR_GOLD
	)

	# 星爆指示器（右上角）
	_show_burst_indicator(burst_count, 1.0)

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"⭐ 星爆！%d 個星爆點！" % burst_count,
		Vector2(vp_size / 2),
		COLOR_GOLD,
		40
	)
	_spawn_float_text(
		"全場 HP -35%%！",
		Vector2(vp_size.x * 0.5, vp_size.y * 0.6),
		COLOR_SKY_BLUE,
		22
	)

	# 5 秒後清除橫幅
	var timer := get_tree().create_timer(5.0)
	timer.timeout.connect(func():
		if is_instance_valid(_banner):
			_banner.queue_free()
			_banner = null
	)

# ── 星爆點爆炸（全服）──────────────────────────────────────────────────────────

func _on_burst_explode(payload: Dictionary) -> void:
	var burst_index: int = payload.get("burst_index", 1)
	var hit_count: int = payload.get("hit_count", 0)
	var accum_mult: float = payload.get("accum_mult", 1.0)

	# 爆炸閃光（顏色隨累積倍率變化）
	var flash_color := COLOR_SKY_BLUE
	if accum_mult >= 4.0:
		flash_color = COLOR_GOLD
	elif accum_mult >= 2.5:
		flash_color = COLOR_ORANGE
	elif accum_mult >= 1.8:
		flash_color = COLOR_PINK
	_flash_screen(flash_color, 1, 0.3)

	# 更新指示器
	if is_instance_valid(_accum_mult_label):
		_accum_mult_label.text = "×%.2f" % accum_mult
		_accum_mult_label.modulate = flash_color

	# 浮動文字
	var vp_size := get_viewport().get_visible_rect().size
	var rand_x := randf_range(vp_size.x * 0.3, vp_size.x * 0.7)
	var rand_y := randf_range(vp_size.y * 0.3, vp_size.y * 0.7)
	_spawn_float_text(
		"⭐ 爆炸 #%d！命中 %d 個！×%.2f" % [burst_index, hit_count, accum_mult],
		Vector2(rand_x, rand_y),
		flash_color,
		18
	)

# ── 星爆共鳴觸發（全服）────────────────────────────────────────────────────────

func _on_burst_resonance(payload: Dictionary) -> void:
	var global_mult: float = payload.get("global_mult", 2.0)
	var global_dur: int = payload.get("global_duration", 5)

	# 全螢幕三次強閃光（彩色）
	_flash_screen(COLOR_PINK, 1, 0.5)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_GOLD, 1, 0.6)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_LIME, 1, 0.5)

	# 大字提示
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"⭐✨ 星爆共鳴！",
		Vector2(vp_size.x * 0.5, vp_size.y * 0.35),
		COLOR_GOLD,
		52
	)
	_spawn_float_text(
		"全服 ×%.1f 加成 %d 秒！" % [global_mult, global_dur],
		Vector2(vp_size.x * 0.5, vp_size.y * 0.5),
		COLOR_SKY_BLUE,
		28
	)

	# 共鳴指示器（右側）
	_show_resonance_indicator(global_mult, global_dur)

# ── 星爆結算（全服）────────────────────────────────────────────────────────────

func _on_burst_end(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var total_bursts: int = payload.get("total_bursts", 0)
	var total_hits: int = payload.get("total_hits", 0)
	var final_mult: float = payload.get("final_mult", 1.0)
	var total_reward: int = payload.get("total_reward", 0)

	# 清除指示器
	if is_instance_valid(_burst_indicator):
		_burst_indicator.queue_free()
		_burst_indicator = null

	# 高倍率才顯示結算彈窗
	if final_mult >= 3.0:
		_flash_screen(COLOR_GOLD, 2, 0.5)
		_show_settle_popup(player_name, total_bursts, total_hits, final_mult, total_reward)

	# 廣播橫幅
	_show_banner(
		"⭐ 星爆結算！",
		"%s 星爆結算！%d 次爆炸，命中 %d 個目標，最終倍率 ×%.1f！" % [player_name, total_bursts, total_hits, final_mult],
		COLOR_GOLD if final_mult >= 3.0 else COLOR_SKY_BLUE
	)

	var timer := get_tree().create_timer(4.0)
	timer.timeout.connect(func():
		if is_instance_valid(_banner):
			_banner.queue_free()
			_banner = null
	)

# ── 輔助方法 ────────────────────────────────────────────────────────────────────

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
		tween.tween_property(flash, "modulate:a", 0.0, 0.12)
		if i < times - 1:
			tween.tween_interval(0.06)
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
	style.bg_color = Color(0, 0, 0, 0.82)
	style.border_color = color
	style.set_border_width_all(2)
	style.set_corner_radius_all(6)
	banner.add_theme_stylebox_override("panel", style)

	var vbox := VBoxContainer.new()
	vbox.mouse_filter = Control.MOUSE_FILTER_IGNORE
	banner.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = title
	title_lbl.add_theme_color_override("font_color", color)
	title_lbl.add_theme_font_size_override("font_size", 14)
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

	# 滑入動畫
	banner.modulate.a = 0.0
	var tween := create_tween()
	tween.tween_property(banner, "modulate:a", 1.0, 0.25)

func _show_burst_indicator(burst_count: int, accum_mult: float) -> void:
	if is_instance_valid(_burst_indicator):
		_burst_indicator.queue_free()

	var vp_size := get_viewport().get_visible_rect().size
	var indicator := PanelContainer.new()
	indicator.position = Vector2(vp_size.x - 148, 80)
	indicator.size = Vector2(140, 72)
	indicator.mouse_filter = Control.MOUSE_FILTER_IGNORE

	var style := StyleBoxFlat.new()
	style.bg_color = Color(0, 0, 0, 0.85)
	style.border_color = COLOR_GOLD
	style.set_border_width_all(2)
	style.set_corner_radius_all(5)
	indicator.add_theme_stylebox_override("panel", style)

	var vbox := VBoxContainer.new()
	vbox.mouse_filter = Control.MOUSE_FILTER_IGNORE
	indicator.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "⭐ 星爆"
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	vbox.add_child(title_lbl)

	var count_lbl := Label.new()
	count_lbl.text = "%d 個星爆點" % burst_count
	count_lbl.add_theme_color_override("font_color", COLOR_SKY_BLUE)
	count_lbl.add_theme_font_size_override("font_size", 11)
	count_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	count_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	vbox.add_child(count_lbl)

	var mult_lbl := Label.new()
	mult_lbl.text = "×%.2f" % accum_mult
	mult_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	mult_lbl.add_theme_font_size_override("font_size", 16)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mult_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	vbox.add_child(mult_lbl)

	add_child(indicator)
	_burst_indicator = indicator
	_accum_mult_label = mult_lbl

	# 脈衝動畫
	var tween := indicator.create_tween().set_loops()
	tween.tween_property(indicator, "modulate:a", 0.7, 0.5)
	tween.tween_property(indicator, "modulate:a", 1.0, 0.5)

func _show_resonance_indicator(global_mult: float, global_dur: int) -> void:
	if is_instance_valid(_resonance_indicator):
		_resonance_indicator.queue_free()
	if _resonance_tween != null:
		_resonance_tween.kill()

	var vp_size := get_viewport().get_visible_rect().size
	var indicator := PanelContainer.new()
	indicator.position = Vector2(vp_size.x - 148, 160)
	indicator.size = Vector2(140, 64)
	indicator.mouse_filter = Control.MOUSE_FILTER_IGNORE

	var style := StyleBoxFlat.new()
	style.bg_color = Color(0, 0, 0, 0.88)
	style.border_color = COLOR_PINK
	style.set_border_width_all(2)
	style.set_corner_radius_all(5)
	indicator.add_theme_stylebox_override("panel", style)

	var vbox := VBoxContainer.new()
	vbox.mouse_filter = Control.MOUSE_FILTER_IGNORE
	indicator.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "⭐✨ 共鳴！"
	title_lbl.add_theme_color_override("font_color", COLOR_PINK)
	title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	vbox.add_child(title_lbl)

	var mult_lbl := Label.new()
	mult_lbl.text = "全服 ×%.1f" % global_mult
	mult_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	mult_lbl.add_theme_font_size_override("font_size", 14)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mult_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	vbox.add_child(mult_lbl)

	var timer_lbl := Label.new()
	timer_lbl.text = "%d 秒" % global_dur
	timer_lbl.add_theme_color_override("font_color", COLOR_SKY_BLUE)
	timer_lbl.add_theme_font_size_override("font_size", 11)
	timer_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	timer_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	vbox.add_child(timer_lbl)

	add_child(indicator)
	_resonance_indicator = indicator
	_resonance_timer_label = timer_lbl

	# 彩虹循環動畫
	_resonance_tween = indicator.create_tween().set_loops()
	_resonance_tween.tween_property(indicator, "modulate", Color(1.0, 0.8, 0.0), 0.4)
	_resonance_tween.tween_property(indicator, "modulate", Color(0.0, 0.8, 1.0), 0.4)
	_resonance_tween.tween_property(indicator, "modulate", Color(1.0, 0.4, 0.7), 0.4)
	_resonance_tween.tween_property(indicator, "modulate", Color(0.5, 1.0, 0.0), 0.4)

	# 倒數計時
	var elapsed := 0.0
	while elapsed < float(global_dur):
		await get_tree().create_timer(1.0).timeout
		elapsed += 1.0
		if is_instance_valid(_resonance_timer_label):
			var remaining: int = max(0, global_dur - int(elapsed))
			_resonance_timer_label.text = "%d 秒" % remaining

	if is_instance_valid(_resonance_indicator):
		_resonance_indicator.queue_free()
		_resonance_indicator = null
	if _resonance_tween != null:
		_resonance_tween.kill()
		_resonance_tween = null

func _show_settle_popup(player_name: String, total_bursts: int, total_hits: int, final_mult: float, total_reward: int) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var popup := PanelContainer.new()
	popup.position = Vector2(vp_size.x + 10, vp_size.y * 0.3)
	popup.size = Vector2(220, 140)
	popup.mouse_filter = Control.MOUSE_FILTER_IGNORE

	var style := StyleBoxFlat.new()
	style.bg_color = Color(0, 0, 0, 0.92)
	style.border_color = COLOR_GOLD
	style.set_border_width_all(3)
	style.set_corner_radius_all(8)
	popup.add_theme_stylebox_override("panel", style)

	var vbox := VBoxContainer.new()
	vbox.mouse_filter = Control.MOUSE_FILTER_IGNORE
	popup.add_child(vbox)

	var lines := [
		["⭐ 星爆結算！", COLOR_GOLD, 16],
		[player_name, COLOR_WHITE, 12],
		["爆炸次數：%d" % total_bursts, COLOR_SKY_BLUE, 12],
		["命中目標：%d" % total_hits, COLOR_SKY_BLUE, 12],
		["最終倍率：×%.2f" % final_mult, COLOR_GOLD, 14],
		["獎勵：+%d" % total_reward, COLOR_LIME, 14],
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

	# 右側滑入動畫
	var tween := create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 230, 0.35)
	tween.tween_interval(3.5)
	tween.tween_property(popup, "position:x", vp_size.x + 10, 0.3)
	tween.tween_callback(popup.queue_free)

func _spawn_float_text(text: String, pos: Vector2, color: Color, font_size: int = 20) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", font_size)
	lbl.position = pos - Vector2(200, 0)
	lbl.size = Vector2(400, 60)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(lbl)

	var tween := create_tween()
	tween.tween_property(lbl, "position:y", pos.y - 60, 1.2)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 1.2)
	tween.tween_callback(lbl.queue_free)
