## LuckyLuckTotemPanel.gd — 幸運幸運圖騰魚 UI 面板（DAY-275）
## 翠綠幸運主題：#00FF88 翠綠 + #FFD700 金 + #00BFFF 天藍 + #7FFF00 草綠
## 業界依據：Fish It Luck Totem（2026）「全場幸運加成」機制進化版
##
## 事件類型：
##   totem_start     — 圖騰開始（個人，GlobalMult/PersonalMult/DurationSec）
##   totem_broadcast — 全服廣播（PlayerName/GlobalMult/DurationSec）
##   totem_kill      — 圖騰期間擊破通知（個人，BonusReward）
##   totem_end       — 圖騰結束（全服，PlayerName/TotalKills/TotalReward）

extends CanvasLayer

const COLOR_GREEN   = Color(0.0,   1.0,   0.533)  # #00FF88 翠綠
const COLOR_GOLD    = Color(1.0,   0.843, 0.0)    # #FFD700 金
const COLOR_CYAN    = Color(0.0,   0.749, 1.0)    # #00BFFF 天藍
const COLOR_LIME    = Color(0.498, 1.0,   0.0)    # #7FFF00 草綠
const COLOR_WHITE   = Color(1.0,   1.0,   1.0)

var _banner: Control = null
var _timer_bar: Control = null
var _totem_indicator: Control = null
var _timer_tween: Tween = null

func _ready() -> void:
	layer = 48  # 比 LuckyFortuneProphecy（47）高一層

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"totem_start":
			_on_totem_start(payload)
		"totem_broadcast":
			_on_totem_broadcast(payload)
		"totem_kill":
			_on_totem_kill(payload)
		"totem_end":
			_on_totem_end(payload)

# ── 圖騰開始 ──────────────────────────────────────────────────────────────────

func _on_totem_start(payload: Dictionary) -> void:
	var global_mult: float = payload.get("global_mult", 1.3)
	var personal_mult: float = payload.get("personal_mult", 1.5)
	var duration_sec: int = payload.get("duration_sec", 15)

	# 翠綠三次強閃光
	_flash_screen(COLOR_GREEN, 3, 0.45)

	# 頂部橫幅
	_show_banner(
		"🍀 幸運圖騰！",
		"全服 ×%.1f + 個人 ×%.1f — %d 秒！" % [global_mult, personal_mult, duration_sec],
		COLOR_GREEN
	)

	# 圖騰指示器（右上角）
	_show_totem_indicator(global_mult, personal_mult)

	# 右側豎向計時條
	_start_timer_bar(duration_sec)

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🍀 幸運圖騰！全服 ×%.1f！" % global_mult,
		Vector2(vp_size / 2),
		COLOR_GREEN,
		46
	)

# ── 全服廣播 ──────────────────────────────────────────────────────────────────

func _on_totem_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var global_mult: float = payload.get("global_mult", 1.3)
	var duration_sec: int = payload.get("duration_sec", 15)
	_show_mini_banner(
		"🍀 %s 觸發幸運圖騰！全服 ×%.1f 加成 %d 秒！" % [player_name, global_mult, duration_sec],
		COLOR_GREEN
	)

# ── 圖騰期間擊破通知 ──────────────────────────────────────────────────────────

func _on_totem_kill(payload: Dictionary) -> void:
	var bonus_reward: int = payload.get("bonus_reward", 0)
	if bonus_reward <= 0:
		return

	# 輕微翠綠閃光
	_flash_screen(COLOR_GREEN, 1, 0.15)

	# 浮動文字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🍀 +%d" % bonus_reward,
		Vector2(vp_size.x * 0.5, vp_size.y * 0.45),
		COLOR_GOLD,
		20
	)

# ── 圖騰結束 ──────────────────────────────────────────────────────────────────

func _on_totem_end(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var total_kills: int = payload.get("total_kills", 0)
	var total_reward: int = payload.get("total_reward", 0)

	# 清除計時條、橫幅、指示器
	_clear_timer_bar()
	_clear_banner()
	_clear_totem_indicator()

	# 翠綠閃光（結束）
	_flash_screen(COLOR_GREEN, 1, 0.25)

	# 全服廣播橫幅
	_show_mini_banner(
		"🍀 幸運圖騰結束！全服共擊破 %d 個目標，總獎勵 +%d！" % [total_kills, total_reward],
		COLOR_GOLD
	)

	# 結算彈窗
	_show_end_popup(player_name, total_kills, total_reward)

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
	panel.modulate = Color(0.0, 0.15, 0.1, 0.92)
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
	tween.tween_interval(3.0)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)

func _show_totem_indicator(global_mult: float, personal_mult: float) -> void:
	_clear_totem_indicator()
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x - 130, 80)
	panel.size = Vector2(120, 60)
	panel.modulate = Color(0.0, 0.15, 0.1, 0.9)
	add_child(panel)
	_totem_indicator = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var lbl1 := Label.new()
	lbl1.text = "🍀 全服 ×%.1f" % global_mult
	lbl1.add_theme_color_override("font_color", COLOR_GREEN)
	lbl1.add_theme_font_size_override("font_size", 14)
	lbl1.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(lbl1)

	var lbl2 := Label.new()
	lbl2.text = "個人 ×%.1f" % personal_mult
	lbl2.add_theme_color_override("font_color", COLOR_GOLD)
	lbl2.add_theme_font_size_override("font_size", 14)
	lbl2.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(lbl2)

	# 脈衝動畫
	var tween := panel.create_tween().set_loops()
	tween.tween_property(panel, "modulate:a", 0.7, 0.5)
	tween.tween_property(panel, "modulate:a", 1.0, 0.5)

func _clear_totem_indicator() -> void:
	if is_instance_valid(_totem_indicator):
		_totem_indicator.queue_free()
	_totem_indicator = null

func _start_timer_bar(duration_sec: int) -> void:
	_clear_timer_bar()
	var vp_size := get_viewport().get_visible_rect().size

	var bar_bg := ColorRect.new()
	bar_bg.position = Vector2(vp_size.x - 338, 80)  # x=-338 與其他計時條錯開
	bar_bg.size = Vector2(8, 200)
	bar_bg.color = Color(0.0, 0.2, 0.1, 0.7)
	add_child(bar_bg)
	_timer_bar = bar_bg

	var bar_fill := ColorRect.new()
	bar_fill.position = Vector2(0, 0)
	bar_fill.size = Vector2(8, 200)
	bar_fill.color = COLOR_GREEN
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

func _show_end_popup(player_name: String, total_kills: int, total_reward: int) -> void:
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x + 10, vp_size.y * 0.3)
	panel.size = Vector2(260, 140)
	panel.modulate = Color(0.0, 0.12, 0.08, 0.95)
	add_child(panel)

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "🍀 幸運圖騰結算"
	title_lbl.add_theme_color_override("font_color", COLOR_GREEN)
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var trigger_lbl := Label.new()
	trigger_lbl.text = "觸發者：%s" % player_name
	trigger_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	trigger_lbl.add_theme_font_size_override("font_size", 13)
	trigger_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(trigger_lbl)

	var kills_lbl := Label.new()
	kills_lbl.text = "全服擊破：%d 個目標" % total_kills
	kills_lbl.add_theme_color_override("font_color", COLOR_CYAN)
	kills_lbl.add_theme_font_size_override("font_size", 14)
	kills_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(kills_lbl)

	var reward_lbl := Label.new()
	reward_lbl.text = "總額外獎勵：+%d" % total_reward
	reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	reward_lbl.add_theme_font_size_override("font_size", 16)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	# 右側滑入動畫
	var tween := create_tween()
	tween.tween_property(panel, "position:x", vp_size.x - 270.0, 0.3).set_ease(Tween.EASE_OUT)
	tween.tween_interval(4.0)
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
