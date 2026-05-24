## LuckyRareChainPanel.gd — 幸運連鎖稀有魚 UI 面板（DAY-280）
## 連鎖稀有主題：#FF6B35 橙紅 + #FFD700 金 + #FF4500 火橙 + #00BFFF 天藍 + #FFFFFF 白
## 業界原創「稀有連鎖+倍率爬升+時間視窗」機制
##
## 事件類型：
##   chain_start           — 連鎖稀有模式觸發（個人，PlayerID/PlayerName/Duration/WindowSec/MaxLayer）
##   chain_broadcast       — 全服廣播橫幅（PlayerName/Duration）
##   chain_kill            — 連鎖擊破（個人，PlayerID/Layer/Mult/Reward/TotalReward）
##   chain_burst           — 連鎖爆發（個人，PlayerID/PlayerName/Layer/Mult/TotalReward）
##   chain_burst_broadcast — 連鎖爆發全服廣播（PlayerName/Layer/Mult/TotalReward）
##   chain_end             — 模式結束（個人，PlayerID/Layer/TotalReward）

extends CanvasLayer

const COLOR_ORANGE     = Color(1.0,   0.420, 0.208)  # #FF6B35 橙紅
const COLOR_GOLD       = Color(1.0,   0.843, 0.0)    # #FFD700 金
const COLOR_FIRE       = Color(1.0,   0.271, 0.0)    # #FF4500 火橙
const COLOR_SKY        = Color(0.0,   0.749, 1.0)    # #00BFFF 天藍
const COLOR_WHITE      = Color(1.0,   1.0,   1.0)
const COLOR_RED        = Color(1.0,   0.2,   0.2)    # 連鎖中斷警告

# 各層顏色（1-5層）
const LAYER_COLORS = [
	Color(0.0, 0.749, 1.0),    # 第1層：天藍
	Color(0.498, 1.0, 0.0),    # 第2層：草綠
	Color(1.0, 0.843, 0.0),    # 第3層：金
	Color(1.0, 0.420, 0.208),  # 第4層：橙紅
	Color(1.0, 0.271, 0.0),    # 第5層：火橙
]

var _banner: Control = null
var _chain_indicator: Control = null
var _layer_label: Label = null
var _mult_label: Label = null
var _window_bar: Control = null
var _window_bar_fill: ColorRect = null
var _window_timer: float = 0.0
var _window_max: float = 8.0
var _is_active: bool = false

func _ready() -> void:
	layer = 53  # 比 LuckyRainbowBridge（52）高一層

func _process(delta: float) -> void:
	if _is_active and _window_bar_fill != null and is_instance_valid(_window_bar_fill):
		_window_timer -= delta
		if _window_timer < 0.0:
			_window_timer = 0.0
		var ratio := _window_timer / _window_max
		if is_instance_valid(_window_bar_fill):
			_window_bar_fill.size.x = ratio * 120.0
			# 顏色隨時間變化：充足→天藍，緊迫→火橙
			if ratio > 0.5:
				_window_bar_fill.color = COLOR_SKY
			elif ratio > 0.25:
				_window_bar_fill.color = COLOR_ORANGE
			else:
				_window_bar_fill.color = COLOR_FIRE

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"chain_start":
			_on_chain_start(payload)
		"chain_broadcast":
			_on_chain_broadcast(payload)
		"chain_kill":
			_on_chain_kill(payload)
		"chain_burst":
			_on_chain_burst(payload)
		"chain_burst_broadcast":
			_on_chain_burst_broadcast(payload)
		"chain_end":
			_on_chain_end(payload)

# ── 連鎖稀有模式觸發（個人）────────────────────────────────────────────────────

func _on_chain_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var duration: int = payload.get("duration", 20)
	var window_sec: int = payload.get("window_sec", 8)
	var max_layer: int = payload.get("max_layer", 5)

	_is_active = true
	_window_max = float(window_sec)
	_window_timer = float(window_sec)

	# 橙紅三次強閃光
	_flash_screen(COLOR_ORANGE, 3, 0.5)

	# 頂部橫幅
	_show_banner(
		"🔗 連鎖稀有模式！",
		"擊破稀有魚（×15+）連鎖！8 秒內連鎖，倍率最高 ×10.0！持續 %d 秒！" % duration,
		COLOR_ORANGE
	)

	# 連鎖指示器（右上角）
	_show_chain_indicator(0, max_layer, window_sec)

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🔗 連鎖稀有模式！",
		Vector2(vp_size / 2),
		COLOR_ORANGE,
		38
	)
	_spawn_float_text(
		"擊破稀有魚（×15+）連鎖倍率！",
		Vector2(vp_size.x * 0.5, vp_size.y * 0.6),
		COLOR_SKY,
		20
	)

	# 5 秒後清除橫幅
	var timer := get_tree().create_timer(5.0)
	timer.timeout.connect(func():
		_clear_banner()
	)

# ── 全服廣播橫幅 ──────────────────────────────────────────────────────────────

func _on_chain_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var duration: int = payload.get("duration", 20)
	_show_mini_banner(
		"🔗 %s 觸發連鎖稀有模式！%d 秒內連鎖擊破稀有魚，倍率最高 ×10.0！" % [player_name, duration],
		COLOR_ORANGE
	)

# ── 連鎖擊破（個人）──────────────────────────────────────────────────────────

func _on_chain_kill(payload: Dictionary) -> void:
	var layer: int = payload.get("layer", 1)
	var mult: float = payload.get("mult", 1.5)
	var reward: int = payload.get("reward", 0)

	# 重置連鎖視窗計時
	_window_timer = _window_max

	# 更新指示器
	_update_chain_layer(layer, mult)

	# 閃光（顏色依層數）
	var flash_color: Color = LAYER_COLORS[min(layer-1, LAYER_COLORS.size()-1)]
	_flash_screen(flash_color, 1, 0.3)

	# 浮動文字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🔗 第%d層！×%.1f！+%d！" % [layer, mult, reward],
		Vector2(vp_size.x * 0.5, vp_size.y * 0.45),
		flash_color,
		28
	)

# ── 連鎖爆發（個人）──────────────────────────────────────────────────────────

func _on_chain_burst(payload: Dictionary) -> void:
	var layer: int = payload.get("layer", 5)
	var mult: float = payload.get("mult", 10.0)
	var total_reward: int = payload.get("total_reward", 0)

	_is_active = false
	_clear_chain_indicator()
	_clear_banner()

	# 火橙三次強閃光
	_flash_screen(COLOR_FIRE, 3, 0.6)

	# 頂部橫幅
	_show_banner(
		"🔗 連鎖爆發！",
		"達成 %d 層連鎖！×%.1f 大獎！總獎勵 +%d！" % [layer, mult, total_reward],
		COLOR_GOLD
	)

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🔗 連鎖爆發！×%.1f！" % mult,
		Vector2(vp_size / 2),
		COLOR_GOLD,
		44
	)
	_spawn_float_text(
		"總獎勵 +%d！" % total_reward,
		Vector2(vp_size.x * 0.5, vp_size.y * 0.6),
		COLOR_FIRE,
		28
	)

	# 結算彈窗
	_show_burst_popup(layer, mult, total_reward)

	# 5 秒後清除橫幅
	var timer := get_tree().create_timer(5.0)
	timer.timeout.connect(func():
		_clear_banner()
	)

# ── 連鎖爆發全服廣播 ──────────────────────────────────────────────────────────

func _on_chain_burst_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var layer: int = payload.get("layer", 5)
	var mult: float = payload.get("mult", 10.0)
	var total_reward: int = payload.get("total_reward", 0)
	_show_mini_banner(
		"🔗 %s 達成 %d 層連鎖爆發！×%.1f 大獎！總獎勵 +%d！" % [player_name, layer, mult, total_reward],
		COLOR_GOLD
	)

# ── 模式結束（個人）──────────────────────────────────────────────────────────

func _on_chain_end(payload: Dictionary) -> void:
	var layer: int = payload.get("layer", 0)
	var total_reward: int = payload.get("total_reward", 0)

	_is_active = false
	_clear_chain_indicator()
	_clear_banner()

	var vp_size := get_viewport().get_visible_rect().size
	if layer > 0:
		_spawn_float_text(
			"🔗 連鎖結束 第%d層 總獎勵 +%d" % [layer, total_reward],
			Vector2(vp_size.x * 0.5, vp_size.y * 0.5),
			COLOR_ORANGE,
			20
		)
	else:
		_spawn_float_text(
			"🔗 連鎖稀有模式結束",
			Vector2(vp_size.x * 0.5, vp_size.y * 0.5),
			COLOR_SKY,
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
	panel.modulate = Color(0.08, 0.03, 0.01, 0.92)
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

func _show_chain_indicator(layer: int, max_layer: int, window_sec: int) -> void:
	_clear_chain_indicator()
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x - 155, 80)
	panel.size = Vector2(145, 110)
	panel.modulate = Color(0.08, 0.03, 0.01, 0.92)
	add_child(panel)
	_chain_indicator = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "🔗 連鎖稀有"
	title_lbl.add_theme_color_override("font_color", COLOR_ORANGE)
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	_layer_label = Label.new()
	_layer_label.text = "第 %d/%d 層" % [layer, max_layer]
	_layer_label.add_theme_color_override("font_color", COLOR_SKY)
	_layer_label.add_theme_font_size_override("font_size", 16)
	_layer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(_layer_label)

	_mult_label = Label.new()
	_mult_label.text = "×1.5"
	_mult_label.add_theme_color_override("font_color", COLOR_GOLD)
	_mult_label.add_theme_font_size_override("font_size", 20)
	_mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(_mult_label)

	# 連鎖視窗進度條
	var bar_bg := ColorRect.new()
	bar_bg.size = Vector2(120, 8)
	bar_bg.color = Color(0.2, 0.2, 0.2, 0.8)
	vbox.add_child(bar_bg)
	_window_bar = bar_bg

	_window_bar_fill = ColorRect.new()
	_window_bar_fill.size = Vector2(120, 8)
	_window_bar_fill.color = COLOR_SKY
	_window_bar_fill.position = Vector2(0, 0)
	bar_bg.add_child(_window_bar_fill)

	# 脈衝動畫
	var tween := panel.create_tween().set_loops()
	tween.tween_property(panel, "modulate:a", 0.75, 0.4)
	tween.tween_property(panel, "modulate:a", 1.0, 0.4)

func _update_chain_layer(layer: int, mult: float) -> void:
	if is_instance_valid(_layer_label):
		_layer_label.text = "第 %d/5 層" % layer
		var c: Color = LAYER_COLORS[min(layer-1, LAYER_COLORS.size()-1)]
		_layer_label.add_theme_color_override("font_color", c)
	if is_instance_valid(_mult_label):
		_mult_label.text = "×%.1f" % mult
		var c: Color = LAYER_COLORS[min(layer-1, LAYER_COLORS.size()-1)]
		_mult_label.add_theme_color_override("font_color", c)

func _clear_chain_indicator() -> void:
	if is_instance_valid(_chain_indicator):
		_chain_indicator.queue_free()
	_chain_indicator = null
	_layer_label = null
	_mult_label = null
	_window_bar = null
	_window_bar_fill = null

func _show_burst_popup(layer: int, mult: float, total_reward: int) -> void:
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x + 10, vp_size.y * 0.3)
	panel.size = Vector2(270, 140)
	panel.modulate = Color(0.08, 0.03, 0.01, 0.95)
	add_child(panel)

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "🔗 連鎖爆發！"
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	title_lbl.add_theme_font_size_override("font_size", 22)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var layer_lbl := Label.new()
	layer_lbl.text = "達成 %d 層連鎖！" % layer
	layer_lbl.add_theme_color_override("font_color", COLOR_FIRE)
	layer_lbl.add_theme_font_size_override("font_size", 16)
	layer_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(layer_lbl)

	var mult_lbl := Label.new()
	mult_lbl.text = "連鎖倍率：×%.1f" % mult
	mult_lbl.add_theme_color_override("font_color", COLOR_ORANGE)
	mult_lbl.add_theme_font_size_override("font_size", 20)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	var reward_lbl := Label.new()
	reward_lbl.text = "總獎勵：+%d 籌碼！" % total_reward
	reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	reward_lbl.add_theme_font_size_override("font_size", 18)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

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
