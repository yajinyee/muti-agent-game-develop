## LuckyQualityMutationPanel.gd — 幸運品質突變魚 UI 面板（DAY-272）
## 品質主題：灰/藍/紫/橙/彩虹 五色品質系統
## 業界依據：Fishing Frenzy Chapter 3 Quality Roll + Fisch Mutation 機制
##
## 事件類型：
##   mutation_start          — 品質突變開始（個人）
##   mutation_broadcast      — 全服廣播（Legendary+）
##   mutation_used           — 品質突變被使用（個人）
##   mutation_result_broadcast — 全服廣播結果（Mythic）
##   mutation_expire         — session 超時（個人）

extends CanvasLayer

# 品質顏色常數
const COLOR_NORMAL    = Color(0.67, 0.67, 0.67)  # #AAAAAA 灰
const COLOR_RARE      = Color(0.29, 0.565, 0.851) # #4A90D9 藍
const COLOR_EPIC      = Color(0.608, 0.349, 0.714) # #9B59B6 紫
const COLOR_LEGENDARY = Color(1.0,  0.549, 0.0)   # #FF8C00 橙
const COLOR_MYTHIC    = Color(1.0,  0.412, 0.706)  # #FF69B4 彩虹粉（動態）
const COLOR_GOLD      = Color(1.0,  0.843, 0.0)   # #FFD700 金

var _banner: Control = null
var _waiting_indicator: Control = null
var _waiting_tween: Tween = null
var _rainbow_tween: Tween = null

func _ready() -> void:
	layer = 45  # 比 LuckyReroll（44）高一層

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"mutation_start":
			_on_mutation_start(payload)
		"mutation_broadcast":
			_on_mutation_broadcast(payload)
		"mutation_used":
			_on_mutation_used(payload)
		"mutation_result_broadcast":
			_on_mutation_result_broadcast(payload)
		"mutation_expire":
			_on_mutation_expire(payload)

# ── 品質突變開始 ──────────────────────────────────────────────────────────────

func _on_mutation_start(payload: Dictionary) -> void:
	var quality: String = payload.get("quality", "normal")
	var quality_mult: float = payload.get("quality_mult", 1.0)
	var quality_color: String = payload.get("quality_color", "#AAAAAA")
	var quality_emoji: String = payload.get("quality_emoji", "⬜")
	var color := _parse_color(quality_color)

	# 閃光次數依品質
	var flash_times := _flash_times_for_quality(quality)
	_flash_screen(color, flash_times)

	# 頂部橫幅
	var quality_name_zh := _quality_name_zh(quality)
	_show_banner(
		"%s %s 突變！" % [quality_emoji, quality_name_zh],
		"×%.1f 等待下一擊！" % quality_mult,
		color
	)

	# 等待指示器（右側）
	_start_waiting_indicator(quality_emoji, quality_mult, color, quality == "mythic")

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"%s ×%.1f 突變就緒！" % [quality_emoji, quality_mult],
		Vector2(vp_size / 2),
		color,
		44
	)

# ── 全服廣播（Legendary+）────────────────────────────────────────────────────

func _on_mutation_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var quality: String = payload.get("quality", "legendary")
	var quality_mult: float = payload.get("quality_mult", 6.0)
	var quality_emoji: String = payload.get("quality_emoji", "🟠")
	var quality_color: String = payload.get("quality_color", "#FF8C00")
	var color := _parse_color(quality_color)
	_show_mini_banner(
		"%s %s 觸發 %s 突變！×%.1f！" % [quality_emoji, player_name, _quality_name_zh(quality), quality_mult],
		color
	)

# ── 品質突變被使用 ────────────────────────────────────────────────────────────

func _on_mutation_used(payload: Dictionary) -> void:
	var target_name: String = payload.get("target_name", "???")
	var reward: int = payload.get("reward", 0)
	var quality: String = payload.get("quality", "normal")
	var quality_mult: float = payload.get("quality_mult", 1.0)
	var quality_emoji: String = payload.get("quality_emoji", "⬜")
	var quality_color: String = payload.get("quality_color", "#AAAAAA")
	var color := _parse_color(quality_color)

	# 清除等待指示器和橫幅
	_clear_waiting_indicator()
	_clear_banner()

	# 閃光
	var flash_times := _flash_times_for_quality(quality)
	_flash_screen(color, flash_times)

	# 大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"%s ×%.1f 命中！+%d！" % [quality_emoji, quality_mult, reward],
		Vector2(vp_size / 2),
		color,
		48
	)

	# 結算彈窗
	_show_result_popup(quality, quality_mult, quality_emoji, color, target_name, reward)

# ── 全服廣播結果（Mythic）────────────────────────────────────────────────────

func _on_mutation_result_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var quality_mult: float = payload.get("quality_mult", 10.0)
	var reward: int = payload.get("reward", 0)
	_show_mini_banner(
		"🌈 %s 神話突變命中！×%.1f +%d！" % [player_name, quality_mult, reward],
		COLOR_MYTHIC
	)

# ── session 超時 ──────────────────────────────────────────────────────────────

func _on_mutation_expire(_payload: Dictionary) -> void:
	_clear_waiting_indicator()
	_clear_banner()
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"✨ 突變超時...",
		Vector2(vp_size.x * 0.5, vp_size.y * 0.6),
		Color(0.6, 0.6, 0.6),
		22
	)

# ── 內部 UI 工具 ──────────────────────────────────────────────────────────────

func _quality_name_zh(quality: String) -> String:
	match quality:
		"normal":    return "普通"
		"rare":      return "稀有"
		"epic":      return "史詩"
		"legendary": return "傳說"
		"mythic":    return "神話"
	return "普通"

func _flash_times_for_quality(quality: String) -> int:
	match quality:
		"mythic":    return 3
		"legendary": return 3
		"epic":      return 2
		"rare":      return 1
		_:           return 1

func _parse_color(hex: String) -> Color:
	# 解析 #RRGGBB 格式
	if hex.begins_with("#") and hex.length() == 7:
		var r := hex.substr(1, 2).hex_to_int() / 255.0
		var g := hex.substr(3, 2).hex_to_int() / 255.0
		var b := hex.substr(5, 2).hex_to_int() / 255.0
		return Color(r, g, b)
	return Color.WHITE

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
	panel.modulate = Color(color.r * 0.3, color.g * 0.3, color.b * 0.3, 0.92)
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

func _start_waiting_indicator(emoji: String, mult: float, color: Color, is_mythic: bool) -> void:
	_clear_waiting_indicator()
	var vp_size := get_viewport().get_visible_rect().size

	var lbl := Label.new()
	lbl.text = "%s ×%.1f" % [emoji, mult]
	lbl.position = Vector2(vp_size.x - 110, vp_size.y * 0.5 - 20)
	lbl.size = Vector2(100, 40)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 22)
	lbl.add_theme_color_override("font_color", color)
	add_child(lbl)
	_waiting_indicator = lbl

	if is_mythic:
		# Mythic：彩虹色循環動畫
		_rainbow_tween = create_tween().set_loops()
		_rainbow_tween.tween_method(func(h: float):
			if is_instance_valid(lbl):
				lbl.add_theme_color_override("font_color", Color.from_hsv(h, 0.9, 1.0))
		, 0.0, 1.0, 1.5)
	else:
		# 普通：脈衝透明度動畫
		_waiting_tween = create_tween().set_loops()
		_waiting_tween.tween_property(lbl, "modulate:a", 0.3, 0.5)
		_waiting_tween.tween_property(lbl, "modulate:a", 1.0, 0.5)

func _clear_waiting_indicator() -> void:
	if _waiting_tween != null:
		_waiting_tween.kill()
		_waiting_tween = null
	if _rainbow_tween != null:
		_rainbow_tween.kill()
		_rainbow_tween = null
	if is_instance_valid(_waiting_indicator):
		_waiting_indicator.queue_free()
	_waiting_indicator = null

func _show_result_popup(quality: String, quality_mult: float, emoji: String, color: Color, target_name: String, reward: int) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var popup := PanelContainer.new()
	popup.size = Vector2(240, 140)
	popup.position = Vector2(vp_size.x + 10, vp_size.y * 0.3)
	popup.modulate = Color(color.r * 0.15, color.g * 0.15, color.b * 0.15, 0.95)
	add_child(popup)

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	popup.add_child(vbox)

	var lbl_title := Label.new()
	lbl_title.text = "%s %s 突變命中！" % [emoji, _quality_name_zh(quality)]
	lbl_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_title.add_theme_font_size_override("font_size", 18)
	lbl_title.add_theme_color_override("font_color", color)
	vbox.add_child(lbl_title)

	var lbl_mult := Label.new()
	lbl_mult.text = "×%.1f 倍率" % quality_mult
	lbl_mult.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_mult.add_theme_font_size_override("font_size", 26)
	lbl_mult.add_theme_color_override("font_color", color)
	vbox.add_child(lbl_mult)

	var lbl_target := Label.new()
	lbl_target.text = "目標：%s" % target_name
	lbl_target.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_target.add_theme_font_size_override("font_size", 14)
	lbl_target.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	vbox.add_child(lbl_target)

	var lbl_reward := Label.new()
	lbl_reward.text = "獎勵 +%d" % reward
	lbl_reward.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl_reward.add_theme_font_size_override("font_size", 20)
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
