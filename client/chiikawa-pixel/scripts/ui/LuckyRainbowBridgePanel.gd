## LuckyRainbowBridgePanel.gd — 幸運彩虹橋魚 UI 面板（DAY-279）
## 彩虹橋主題：#FF69B4 粉紅 + #FFD700 金 + #00BFFF 天藍 + #7FFF00 草綠 + #FF8C00 橙
## 業界原創「彩虹橋連接+跨目標連鎖傷害+彩虹爆發」機制
##
## 事件類型：
##   bridge_start     — 彩虹橋觸發（全服，PlayerID/PlayerName/TargetIIDs/TargetNames/Duration）
##   bridge_chain     — 連鎖傷害（全服，PlayerName/KilledIID/OtherIIDs/KilledCount/TotalCount）
##   bridge_burst     — 彩虹爆發（全服，PlayerName/BurstMult/BurstSeconds）
##   bridge_burst_end — 彩虹爆發結束（全服）
##   bridge_fade      — 彩虹消散（全服，PlayerName/KilledCount/TotalCount/FadedCount）

extends CanvasLayer

const COLOR_PINK       = Color(1.0,   0.412, 0.706)  # #FF69B4 粉紅
const COLOR_GOLD       = Color(1.0,   0.843, 0.0)    # #FFD700 金
const COLOR_SKY        = Color(0.0,   0.749, 1.0)    # #00BFFF 天藍
const COLOR_GREEN      = Color(0.498, 1.0,   0.0)    # #7FFF00 草綠
const COLOR_ORANGE     = Color(1.0,   0.549, 0.0)    # #FF8C00 橙
const COLOR_WHITE      = Color(1.0,   1.0,   1.0)

# 彩虹顏色循環（用於爆發動畫）
const RAINBOW_COLORS = [
	Color(1.0, 0.0, 0.0),    # 紅
	Color(1.0, 0.5, 0.0),    # 橙
	Color(1.0, 1.0, 0.0),    # 黃
	Color(0.0, 1.0, 0.0),    # 綠
	Color(0.0, 0.5, 1.0),    # 藍
	Color(0.5, 0.0, 1.0),    # 紫
]

var _banner: Control = null
var _bridge_indicator: Control = null
var _bridge_count_label: Label = null
var _burst_indicator: Control = null
var _burst_tween: Tween = null

func _ready() -> void:
	layer = 52  # 比 LuckyTimeRift（51）高一層

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"bridge_start":
			_on_bridge_start(payload)
		"bridge_chain":
			_on_bridge_chain(payload)
		"bridge_burst":
			_on_bridge_burst(payload)
		"bridge_burst_end":
			_on_bridge_burst_end()
		"bridge_fade":
			_on_bridge_fade(payload)

# ── 彩虹橋觸發（全服）────────────────────────────────────────────────────────

func _on_bridge_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var target_names: Array = payload.get("target_names", [])
	var duration: int = payload.get("duration", 12)
	var total_count: int = target_names.size()

	# 彩虹三次強閃光
	_flash_rainbow(3)

	# 頂部橫幅
	var names_str := ", ".join(target_names) if target_names.size() > 0 else "???"
	_show_banner(
		"🌈 彩虹橋！",
		"%s 連接了 %d 個目標！打一個，其他 HP -40%%！%d 秒內全打完觸發彩虹爆發！" % [player_name, total_count, duration],
		COLOR_PINK
	)

	# 彩虹橋指示器（右上角）
	_show_bridge_indicator(total_count, 0, duration)

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🌈 彩虹橋！%d 個目標連接！" % total_count,
		Vector2(vp_size / 2),
		COLOR_PINK,
		38
	)
	_spawn_float_text(
		"打一個，其他 HP -40%%！",
		Vector2(vp_size.x * 0.5, vp_size.y * 0.6),
		COLOR_SKY,
		22
	)

	# 5 秒後清除橫幅
	var timer := get_tree().create_timer(5.0)
	timer.timeout.connect(func():
		_clear_banner()
	)

# ── 連鎖傷害（全服）──────────────────────────────────────────────────────────

func _on_bridge_chain(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var killed_count: int = payload.get("killed_count", 0)
	var total_count: int = payload.get("total_count", 3)
	var other_count: int = payload.get("other_iids", []).size()

	# 更新指示器
	_update_bridge_count(killed_count, total_count)

	# 輕微彩虹閃光
	_flash_screen(RAINBOW_COLORS[killed_count % RAINBOW_COLORS.size()], 1, 0.25)

	# 浮動文字
	var vp_size := get_viewport().get_visible_rect().size
	if other_count > 0:
		_spawn_float_text(
			"🌈 連鎖！其他 %d 個目標 HP -40%%！" % other_count,
			Vector2(vp_size.x * 0.5, vp_size.y * 0.45),
			COLOR_SKY,
			20
		)
	_spawn_float_text(
		"擊破 %d/%d" % [killed_count, total_count],
		Vector2(vp_size.x * 0.5, vp_size.y * 0.55),
		COLOR_GOLD,
		24
	)

# ── 彩虹爆發（全服）──────────────────────────────────────────────────────────

func _on_bridge_burst(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var burst_mult: float = payload.get("burst_mult", 2.0)
	var burst_seconds: int = payload.get("burst_seconds", 6)

	# 清除橋指示器
	_clear_bridge_indicator()

	# 全螢幕彩虹三次強閃光
	_flash_rainbow(3)

	# 頂部橫幅
	_show_banner(
		"🌈 彩虹爆發！",
		"%s 擊破全部目標！全服 ×%.1f 加成 %d 秒！" % [player_name, burst_mult, burst_seconds],
		COLOR_GOLD
	)

	# 彩虹爆發指示器（右上角）
	_show_burst_indicator(burst_mult, burst_seconds)

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🌈 彩虹爆發！全服 ×%.1f！" % burst_mult,
		Vector2(vp_size / 2),
		COLOR_GOLD,
		44
	)
	_spawn_float_text(
		"持續 %d 秒！" % burst_seconds,
		Vector2(vp_size.x * 0.5, vp_size.y * 0.6),
		COLOR_PINK,
		28
	)

	# 結算彈窗
	_show_burst_popup(player_name, burst_mult, burst_seconds)

	# 橫幅 6 秒後清除
	var timer := get_tree().create_timer(float(burst_seconds))
	timer.timeout.connect(func():
		_clear_banner()
	)

# ── 彩虹爆發結束 ──────────────────────────────────────────────────────────────

func _on_bridge_burst_end() -> void:
	_clear_burst_indicator()
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🌈 彩虹爆發結束",
		Vector2(vp_size.x * 0.5, vp_size.y * 0.5),
		COLOR_PINK,
		20
	)

# ── 彩虹消散（全服）──────────────────────────────────────────────────────────

func _on_bridge_fade(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var killed_count: int = payload.get("killed_count", 0)
	var total_count: int = payload.get("total_count", 3)
	var faded_count: int = payload.get("faded_count", 0)

	# 清除指示器
	_clear_bridge_indicator()
	_clear_banner()

	# 灰色閃光（消散感）
	_flash_screen(Color(0.6, 0.6, 0.6), 1, 0.2)

	# 浮動文字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🌈 彩虹消散... 擊破 %d/%d" % [killed_count, total_count],
		Vector2(vp_size / 2),
		COLOR_ORANGE,
		28
	)
	if faded_count > 0:
		_spawn_float_text(
			"剩餘 %d 個目標 HP -60%%（安慰獎）" % faded_count,
			Vector2(vp_size.x * 0.5, vp_size.y * 0.6),
			COLOR_SKY,
			18
		)

# ── 內部 UI 工具 ──────────────────────────────────────────────────────────────

func _flash_rainbow(times: int) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var flash := ColorRect.new()
	flash.size = vp_size
	flash.color = Color(1.0, 0.0, 0.0, 0.0)
	add_child(flash)

	var tween := create_tween()
	for i in range(times):
		var c: Variant = RAINBOW_COLORS[i % RAINBOW_COLORS.size()]
		tween.tween_property(flash, "color", Color(c.r, c.g, c.b, 0.45), 0.07)
		tween.tween_property(flash, "color:a", 0.0, 0.10)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)

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
	panel.modulate = Color(0.05, 0.02, 0.08, 0.92)
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

func _show_bridge_indicator(total: int, killed: int, duration: int) -> void:
	_clear_bridge_indicator()
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x - 150, 80)
	panel.size = Vector2(140, 90)
	panel.modulate = Color(0.05, 0.02, 0.08, 0.92)
	add_child(panel)
	_bridge_indicator = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "🌈 彩虹橋"
	title_lbl.add_theme_color_override("font_color", COLOR_PINK)
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	_bridge_count_label = Label.new()
	_bridge_count_label.text = "擊破 %d/%d" % [killed, total]
	_bridge_count_label.add_theme_color_override("font_color", COLOR_GOLD)
	_bridge_count_label.add_theme_font_size_override("font_size", 18)
	_bridge_count_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(_bridge_count_label)

	var dur_lbl := Label.new()
	dur_lbl.text = "⏱ %d 秒" % duration
	dur_lbl.add_theme_color_override("font_color", COLOR_SKY)
	dur_lbl.add_theme_font_size_override("font_size", 13)
	dur_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(dur_lbl)

	# 彩虹脈衝動畫
	var tween := panel.create_tween().set_loops()
	tween.tween_property(panel, "modulate:a", 0.7, 0.5)
	tween.tween_property(panel, "modulate:a", 1.0, 0.5)

func _update_bridge_count(killed: int, total: int) -> void:
	if is_instance_valid(_bridge_count_label):
		_bridge_count_label.text = "擊破 %d/%d" % [killed, total]
		# 顏色隨進度變化
		var progress := float(killed) / float(total) if total > 0 else 0.0
		if progress >= 0.67:
			_bridge_count_label.add_theme_color_override("font_color", COLOR_GOLD)
		elif progress >= 0.33:
			_bridge_count_label.add_theme_color_override("font_color", COLOR_ORANGE)
		else:
			_bridge_count_label.add_theme_color_override("font_color", COLOR_SKY)

func _clear_bridge_indicator() -> void:
	if is_instance_valid(_bridge_indicator):
		_bridge_indicator.queue_free()
	_bridge_indicator = null
	_bridge_count_label = null

func _show_burst_indicator(burst_mult: float, burst_seconds: int) -> void:
	_clear_burst_indicator()
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x - 150, 80)
	panel.size = Vector2(140, 80)
	panel.modulate = Color(0.05, 0.02, 0.08, 0.95)
	add_child(panel)
	_burst_indicator = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "🌈 彩虹爆發！"
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var mult_lbl := Label.new()
	mult_lbl.text = "全服 ×%.1f" % burst_mult
	mult_lbl.add_theme_color_override("font_color", COLOR_PINK)
	mult_lbl.add_theme_font_size_override("font_size", 20)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	# 彩虹循環動畫
	_burst_tween = panel.create_tween().set_loops()
	for c in RAINBOW_COLORS:
		_burst_tween.tween_property(mult_lbl, "modulate", c, 0.3)

func _clear_burst_indicator() -> void:
	if _burst_tween != null:
		_burst_tween.kill()
		_burst_tween = null
	if is_instance_valid(_burst_indicator):
		_burst_indicator.queue_free()
	_burst_indicator = null

func _show_burst_popup(player_name: String, burst_mult: float, burst_seconds: int) -> void:
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x + 10, vp_size.y * 0.3)
	panel.size = Vector2(270, 140)
	panel.modulate = Color(0.05, 0.02, 0.08, 0.95)
	add_child(panel)

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "🌈 彩虹爆發！"
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	title_lbl.add_theme_font_size_override("font_size", 22)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var trigger_lbl := Label.new()
	trigger_lbl.text = "觸發者：%s" % player_name
	trigger_lbl.add_theme_color_override("font_color", COLOR_PINK)
	trigger_lbl.add_theme_font_size_override("font_size", 14)
	trigger_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(trigger_lbl)

	var mult_lbl := Label.new()
	mult_lbl.text = "全服 ×%.1f 加成" % burst_mult
	mult_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	mult_lbl.add_theme_font_size_override("font_size", 20)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	var dur_lbl := Label.new()
	dur_lbl.text = "持續 %d 秒！" % burst_seconds
	dur_lbl.add_theme_color_override("font_color", COLOR_SKY)
	dur_lbl.add_theme_font_size_override("font_size", 16)
	dur_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(dur_lbl)

	# 右側滑入動畫
	var tween := create_tween()
	tween.tween_property(panel, "position:x", vp_size.x - 280.0, 0.3).set_ease(Tween.EASE_OUT)
	tween.tween_interval(5.0)
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
