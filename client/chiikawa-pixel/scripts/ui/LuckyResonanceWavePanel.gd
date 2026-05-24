## LuckyResonanceWavePanel.gd — 幸運共鳴波魚 UI 面板（DAY-273）
## 天藍波紋主題：#00BFFF 天藍 + #1E90FF 道奇藍 + #FFD700 金 + #00FF88 翠綠
## 業界依據：Royal Fishing / Jili 2026「連鎖閃電+群體攻擊」趨勢進化版
##
## 事件類型：
##   wave_start          — 共鳴波開始（個人）
##   wave_broadcast      — 全服廣播
##   wave_layer          — 每層波結果（Layer/Radius/AffectedCount/ExplodeCount/Mult/Reward）
##   wave_burst          — 共鳴爆發（TotalExplode/TotalReward/BurstMult/BurstDurSec）
##   wave_burst_end      — 共鳴爆發結束
##   wave_result         — 未達爆發門檻的結算

extends CanvasLayer

const COLOR_WAVE    = Color(0.0,   0.749, 1.0)    # #00BFFF 天藍
const COLOR_DEEP    = Color(0.118, 0.565, 1.0)    # #1E90FF 道奇藍
const COLOR_GOLD    = Color(1.0,   0.843, 0.0)    # #FFD700 金
const COLOR_GREEN   = Color(0.0,   1.0,   0.533)  # #00FF88 翠綠
const COLOR_WHITE   = Color(1.0,   1.0,   1.0)
const COLOR_BURST   = Color(0.0,   1.0,   0.533)  # 爆發時翠綠

var _banner: Control = null
var _layer_counter: Control = null
var _burst_indicator: Control = null
var _burst_tween: Tween = null

func _ready() -> void:
	layer = 46  # 比 LuckyQualityMutation（45）高一層

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"wave_start":
			_on_wave_start(payload)
		"wave_broadcast":
			_on_wave_broadcast(payload)
		"wave_layer":
			_on_wave_layer(payload)
		"wave_burst":
			_on_wave_burst(payload)
		"wave_burst_end":
			_on_wave_burst_end(payload)
		"wave_result":
			_on_wave_result(payload)

# ── 共鳴波開始 ────────────────────────────────────────────────────────────────

func _on_wave_start(payload: Dictionary) -> void:
	# 天藍色閃光
	_flash_screen(COLOR_WAVE, 2)

	# 頂部橫幅
	_show_banner("🌊 共鳴波！", "3 層同心圓擴散中...", COLOR_WAVE)

	# 層數計數器（右上角）
	_show_layer_counter(0)

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🌊 共鳴波觸發！",
		Vector2(vp_size / 2),
		COLOR_WAVE,
		42
	)

# ── 全服廣播 ──────────────────────────────────────────────────────────────────

func _on_wave_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	_show_mini_banner(
		"🌊 %s 觸發共鳴波！3 層擴散！" % player_name,
		COLOR_WAVE
	)

# ── 每層波結果 ────────────────────────────────────────────────────────────────

func _on_wave_layer(payload: Dictionary) -> void:
	var layer_num: int = payload.get("layer", 1)
	var explode_count: int = payload.get("explode_count", 0)
	var affected_count: int = payload.get("affected_count", 0)
	var mult: float = payload.get("mult", 1.0)
	var reward: int = payload.get("reward", 0)
	var radius: float = payload.get("radius", 150.0)

	# 更新層數計數器
	_update_layer_counter(layer_num)

	# 輕微閃光（依引爆數）
	if explode_count > 0:
		_flash_screen(COLOR_WAVE, 1, 0.25)

	# 浮動文字（顯示引爆數）
	var vp_size := get_viewport().get_visible_rect().size
	if explode_count > 0:
		_spawn_float_text(
			"🌊 第%d層：%d 個引爆！×%.1f +%d" % [layer_num, explode_count, mult, reward],
			Vector2(vp_size.x * 0.5, vp_size.y * 0.4 + layer_num * 40),
			COLOR_WAVE,
			22
		)
	else:
		_spawn_float_text(
			"🌊 第%d層：波及 %d 個目標" % [layer_num, affected_count],
			Vector2(vp_size.x * 0.5, vp_size.y * 0.4 + layer_num * 40),
			Color(0.5, 0.7, 0.9),
			18
		)

# ── 共鳴爆發 ──────────────────────────────────────────────────────────────────

func _on_wave_burst(payload: Dictionary) -> void:
	var total_explode: int = payload.get("total_explode", 0)
	var total_reward: int = payload.get("total_reward", 0)
	var burst_mult: float = payload.get("burst_mult", 1.5)
	var burst_dur_sec: int = payload.get("burst_dur_sec", 8)

	# 清除層數計數器和橫幅
	_clear_layer_counter()
	_clear_banner()

	# 全螢幕三次強閃光（翠綠）
	_flash_screen(COLOR_BURST, 3, 0.5)

	# 大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🌊 共鳴爆發！×%.1f 全服加成！" % burst_mult,
		Vector2(vp_size / 2),
		COLOR_BURST,
		50
	)

	# 爆發指示器（右側常駐）
	_start_burst_indicator(burst_mult, burst_dur_sec)

	# 結算彈窗
	_show_burst_popup(total_explode, total_reward, burst_mult, burst_dur_sec)

# ── 共鳴爆發結束 ──────────────────────────────────────────────────────────────

func _on_wave_burst_end(_payload: Dictionary) -> void:
	_clear_burst_indicator()
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🌊 共鳴爆發結束",
		Vector2(vp_size.x * 0.5, vp_size.y * 0.6),
		Color(0.5, 0.7, 0.9),
		20
	)

# ── 未達爆發門檻的結算 ────────────────────────────────────────────────────────

func _on_wave_result(payload: Dictionary) -> void:
	var total_explode: int = payload.get("total_explode", 0)
	var total_reward: int = payload.get("total_reward", 0)

	_clear_layer_counter()
	_clear_banner()

	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🌊 共鳴波結束：%d 個引爆 +%d" % [total_explode, total_reward],
		Vector2(vp_size.x * 0.5, vp_size.y * 0.5),
		COLOR_WAVE,
		24
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
	panel.modulate = Color(0.0, 0.2, 0.3, 0.92)
	add_child(panel)
	_banner = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var lbl1 := Label.new()
	lbl1.text = title
	lbl1.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl1.add_theme_font_size_override("font_size", 28)
	lbl1.add_theme_color_override("font_color", color)
	vbox.add_child(lbl1)

	var lbl2 := Label.new()
	lbl2.text = subtitle
	lbl2.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl2.add_theme_font_size_override("font_size", 18)
	lbl2.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
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

func _show_layer_counter(current_layer: int) -> void:
	_clear_layer_counter()
	var vp_size := get_viewport().get_visible_rect().size

	var lbl := Label.new()
	lbl.text = "🌊 第 %d/3 層" % current_layer
	lbl.position = Vector2(vp_size.x - 120, 80)
	lbl.size = Vector2(110, 36)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 20)
	lbl.add_theme_color_override("font_color", COLOR_WAVE)
	add_child(lbl)
	_layer_counter = lbl

func _update_layer_counter(current_layer: int) -> void:
	if is_instance_valid(_layer_counter):
		_layer_counter.text = "🌊 第 %d/3 層" % current_layer
		# 顏色隨層數加深
		match current_layer:
			1: _layer_counter.add_theme_color_override("font_color", COLOR_WAVE)
			2: _layer_counter.add_theme_color_override("font_color", COLOR_DEEP)
			3: _layer_counter.add_theme_color_override("font_color", COLOR_GOLD)

func _clear_layer_counter() -> void:
	if is_instance_valid(_layer_counter):
		_layer_counter.queue_free()
	_layer_counter = null

func _start_burst_indicator(burst_mult: float, dur_sec: int) -> void:
	_clear_burst_indicator()
	var vp_size := get_viewport().get_visible_rect().size

	var lbl := Label.new()
	lbl.text = "🌊 ×%.1f 全服！" % burst_mult
	lbl.position = Vector2(vp_size.x - 130, vp_size.y * 0.5 - 20)
	lbl.size = Vector2(120, 40)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 22)
	lbl.add_theme_color_override("font_color", COLOR_BURST)
	add_child(lbl)
	_burst_indicator = lbl

	# 脈衝動畫
	_burst_tween = create_tween().set_loops()
	_burst_tween.tween_property(lbl, "modulate:a", 0.4, 0.6)
	_burst_tween.tween_property(lbl, "modulate:a", 1.0, 0.6)

	# 自動消失
	var timer := create_tween()
	timer.tween_interval(float(dur_sec))
	timer.tween_callback(func():
		_clear_burst_indicator()
	)

func _clear_burst_indicator() -> void:
	if _burst_tween != null:
		_burst_tween.kill()
		_burst_tween = null
	if is_instance_valid(_burst_indicator):
		_burst_indicator.queue_free()
	_burst_indicator = null

func _show_burst_popup(total_explode: int, total_reward: int, burst_mult: float, dur_sec: int) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var popup := PanelContainer.new()
	popup.size = Vector2(240, 140)
	popup.position = Vector2(vp_size.x + 10, vp_size.y * 0.3)
	popup.modulate = Color(0.0, 0.15, 0.25, 0.95)
	add_child(popup)

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	popup.add_child(vbox)

	var lbl_title := Label.new()
	lbl_title.text = "🌊 共鳴爆發！"
	lbl_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_title.add_theme_font_size_override("font_size", 22)
	lbl_title.add_theme_color_override("font_color", COLOR_BURST)
	vbox.add_child(lbl_title)

	var lbl_explode := Label.new()
	lbl_explode.text = "引爆 %d 個目標" % total_explode
	lbl_explode.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_explode.add_theme_font_size_override("font_size", 16)
	lbl_explode.add_theme_color_override("font_color", COLOR_WAVE)
	vbox.add_child(lbl_explode)

	var lbl_burst := Label.new()
	lbl_burst.text = "全服 ×%.1f 加成 %ds" % [burst_mult, dur_sec]
	lbl_burst.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_burst.add_theme_font_size_override("font_size", 20)
	lbl_burst.add_theme_color_override("font_color", COLOR_GOLD)
	vbox.add_child(lbl_burst)

	var lbl_reward := Label.new()
	lbl_reward.text = "總獎勵 +%d" % total_reward
	lbl_reward.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_reward.add_theme_font_size_override("font_size", 18)
	lbl_reward.add_theme_color_override("font_color", COLOR_GOLD)
	vbox.add_child(lbl_reward)

	var tween := create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 250.0, 0.3)
	tween.tween_interval(5.0)
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
