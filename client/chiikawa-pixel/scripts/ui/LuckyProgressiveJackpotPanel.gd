## LuckyProgressiveJackpotPanel.gd — 幸運累積大獎池魚系統面板（DAY-262）
## 業界原創「全服累積大獎池+貢獻比例分配+大獎池爆發」機制
##
## 視覺設計：
##   - 金色大獎主題（#FFD700 金 + #FF8C00 橙金 + #1A1A2E 深藍黑 + #FFFFFF 白）
##   - jackpot_update：右上角大獎池金額顯示（持續更新，脈衝動畫）
##   - jackpot_burst：全螢幕三次強閃光 + 「💰 大獎池爆發！」大字 + 個人獎勵 + 結算彈窗
##   - jackpot_burst_broadcast：全服廣播橫幅（含觸發者/最高獎勵者）
extends CanvasLayer

# 主題顏色
const COLOR_GOLD    = Color("#FFD700")  # 金色（主色）
const COLOR_ORANGE  = Color("#FF8C00")  # 橙金（大獎池）
const COLOR_DARK    = Color("#1A1A2E")  # 深藍黑（背景）
const COLOR_WHITE   = Color("#FFFFFF")  # 白色
const COLOR_GREEN   = Color("#00FF88")  # 綠色（個人獎勵）
const COLOR_CREAM   = Color("#FFF8DC")  # 奶油（副文字）

# 大獎池顯示
var _pool_label: Label = null
var _pool_bg: PanelContainer = null
var _pool_tween: Tween = null
var _current_pool: int = 0

func _ready() -> void:
	layer = 35  # 幸運累積大獎池魚面板層級（DAY-262）
	# 初始化大獎池顯示
	_init_pool_display()

## 初始化大獎池顯示（右上角常駐）
func _init_pool_display() -> void:
	var vp_size = get_viewport().size

	_pool_bg = PanelContainer.new()
	_pool_bg.size = Vector2(160, 40)
	_pool_bg.position = Vector2(vp_size.x - 170.0, 10.0)
	add_child(_pool_bg)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.05, 0.1, 0.75)
	style.border_color = COLOR_GOLD
	style.set_border_width_all(2)
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	_pool_bg.add_theme_stylebox_override("panel", style)

	_pool_label = Label.new()
	_pool_label.text = "💰 大獎池：---"
	_pool_label.add_theme_font_size_override("font_size", 13)
	_pool_label.add_theme_color_override("font_color", COLOR_GOLD)
	_pool_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_pool_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_pool_label.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_pool_bg.add_child(_pool_label)

## 處理幸運累積大獎池魚訊息
func handle_lucky_progressive_jackpot(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"jackpot_update":
			_on_jackpot_update(payload)
		"jackpot_burst":
			_on_jackpot_burst(payload)
		"jackpot_burst_broadcast":
			_on_jackpot_burst_broadcast(payload)

## jackpot_update — 大獎池定期更新
func _on_jackpot_update(payload: Dictionary) -> void:
	var pool: int = payload.get("pool", 0)
	_update_pool_display(pool)

## jackpot_burst — 大獎池爆發（個人）
func _on_jackpot_burst(payload: Dictionary) -> void:
	var trigger_name: String = payload.get("trigger_name", "某玩家")
	var pool: int = payload.get("pool", 0)
	var kills: int = payload.get("kills", 0)
	var total_kills: int = payload.get("total_kills", 0)
	var pct: float = payload.get("pct", 0.0)
	var reward: int = payload.get("reward", 0)

	# 更新大獎池顯示（爆發後重置）
	_update_pool_display(100)

	# 全螢幕三次強閃光（金色）
	_flash_screen(COLOR_GOLD, 0.9, 3)

	# 大字
	_show_big_text("💰 大獎池爆發！", COLOR_GOLD, 52, 3.0)
	_show_sub_text("由 %s 觸發！你貢獻了 %d/%d 次（%.1f%%）！" % [trigger_name, kills, total_kills, pct * 100], COLOR_CREAM, 3.0)

	# 個人獎勵浮動文字
	_show_float_text("💰 大獎池分配！+%d" % reward, COLOR_GREEN, 2.5)

	# 結算彈窗
	_show_burst_popup(trigger_name, pool, kills, total_kills, pct, reward)

## jackpot_burst_broadcast — 大獎池爆發全服廣播
func _on_jackpot_burst_broadcast(payload: Dictionary) -> void:
	var trigger_name: String = payload.get("trigger_name", "某玩家")
	var pool: int = payload.get("pool", 0)
	var top_name: String = payload.get("top_name", "")
	var top_reward: int = payload.get("top_reward", 0)
	var player_count: int = payload.get("player_count", 0)
	_show_top_banner("💰 %s 觸發大獎池爆發！池=%d，%d 位玩家分配！最高：%s +%d！" % [trigger_name, pool, player_count, top_name, top_reward], COLOR_GOLD, 5.0)

# ─── 大獎池顯示 ───────────────────────────────────────────────────────────────

func _update_pool_display(pool: int) -> void:
	_current_pool = pool
	if not is_instance_valid(_pool_label):
		return
	_pool_label.text = "💰 大獎池：%d" % pool

	# 大獎池越大，顏色越亮
	if pool >= 5000:
		_pool_label.add_theme_color_override("font_color", COLOR_WHITE)
		# 快速脈衝
		if is_instance_valid(_pool_tween):
			_pool_tween.kill()
		_pool_tween = _pool_label.create_tween().set_loops()
		_pool_tween.tween_property(_pool_label, "modulate:a", 0.4, 0.3)
		_pool_tween.tween_property(_pool_label, "modulate:a", 1.0, 0.3)
	elif pool >= 2000:
		_pool_label.add_theme_color_override("font_color", COLOR_ORANGE)
	else:
		_pool_label.add_theme_color_override("font_color", COLOR_GOLD)

# ─── 大獎池爆發結算彈窗 ───────────────────────────────────────────────────────

func _show_burst_popup(trigger_name: String, pool: int, kills: int, total_kills: int, pct: float, reward: int) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(360, 230)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 115)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.04, 0.0, 0.95)
	style.border_color = COLOR_GOLD
	style.set_border_width_all(3)
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 7)
	popup.add_child(vbox)

	var title_lbl = Label.new()
	title_lbl.text = "💰 大獎池爆發！"
	title_lbl.add_theme_font_size_override("font_size", 26)
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var trigger_lbl = Label.new()
	trigger_lbl.text = "觸發者：%s" % trigger_name
	trigger_lbl.add_theme_font_size_override("font_size", 13)
	trigger_lbl.add_theme_color_override("font_color", COLOR_CREAM)
	trigger_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(trigger_lbl)

	var pool_lbl = Label.new()
	pool_lbl.text = "大獎池總額：%d" % pool
	pool_lbl.add_theme_font_size_override("font_size", 16)
	pool_lbl.add_theme_color_override("font_color", COLOR_ORANGE)
	pool_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(pool_lbl)

	var contrib_lbl = Label.new()
	contrib_lbl.text = "你的貢獻：%d/%d 次（%.1f%%）" % [kills, total_kills, pct * 100]
	contrib_lbl.add_theme_font_size_override("font_size", 14)
	contrib_lbl.add_theme_color_override("font_color", COLOR_CREAM)
	contrib_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(contrib_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "你的分配：+%d" % reward
	reward_lbl.add_theme_font_size_override("font_size", 22)
	reward_lbl.add_theme_color_override("font_color", COLOR_GREEN)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	var tip_lbl = Label.new()
	tip_lbl.text = "打得越多，分得越多！"
	tip_lbl.add_theme_font_size_override("font_size", 12)
	tip_lbl.add_theme_color_override("font_color", COLOR_CREAM)
	tip_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(tip_lbl)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 380.0, 0.4).set_ease(Tween.EASE_OUT)
	tween.tween_interval(6.0)
	tween.tween_property(popup, "position:x", vp_size.x + 10.0, 0.4).set_ease(Tween.EASE_IN)
	tween.tween_callback(popup.queue_free)

# ─── 通用 UI 工具 ─────────────────────────────────────────────────────────────

func _flash_screen(color: Color, alpha: float, times: int) -> void:
	var vp_size = get_viewport().size
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.0)
	flash.size = vp_size
	add_child(flash)
	var tween = create_tween()
	for i in range(times):
		tween.tween_property(flash, "color:a", alpha, 0.08)
		tween.tween_property(flash, "color:a", 0.0, 0.12)
	tween.tween_callback(flash.queue_free)

func _show_top_banner(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var banner = PanelContainer.new()
	banner.size = Vector2(vp_size.x * 0.75, 38)
	banner.position = Vector2(vp_size.x * 0.125, 52)
	add_child(banner)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(color.r, color.g, color.b, 0.88)
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	banner.add_theme_stylebox_override("panel", style)
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 13)
	lbl.add_theme_color_override("font_color", Color(0.05, 0.04, 0.0))
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	banner.add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(banner.queue_free)

func _show_big_text(text: String, color: Color, font_size: int, duration: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", font_size)
	lbl.add_theme_color_override("font_color", color)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.size = Vector2(vp_size.x, 80)
	lbl.position = Vector2(0, vp_size.y * 0.35)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "scale", Vector2(1.2, 1.2), 0.15).set_ease(Tween.EASE_OUT)
	tween.tween_property(lbl, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(duration - 0.5)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.5)
	tween.tween_callback(lbl.queue_free)

func _show_sub_text(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 13)
	lbl.add_theme_color_override("font_color", color)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.size = Vector2(vp_size.x, 40)
	lbl.position = Vector2(0, vp_size.y * 0.35 + 80)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.4)
	tween.tween_callback(lbl.queue_free)

func _show_float_text(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 18)
	lbl.add_theme_color_override("font_color", color)
	lbl.position = Vector2(
		randf_range(vp_size.x * 0.25, vp_size.x * 0.65),
		randf_range(vp_size.y * 0.25, vp_size.y * 0.55)
	)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 70, duration)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, duration)
	tween.tween_callback(lbl.queue_free)
