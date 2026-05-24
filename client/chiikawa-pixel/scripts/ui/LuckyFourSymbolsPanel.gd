## LuckyFourSymbolsPanel.gd — 幸運四象大獎魚 UI 面板（DAY-283）
## 四象主題：青龍（綠）/ 白虎（銀）/ 朱雀（紅）/ 玄武（深藍）
## 業界依據：Jackpot Fishing by Jili「四層累進大獎」機制 + 中華文化「四象」主題
##
## 事件類型：
##   symbol_trigger    — 四象大獎觸發（全服，PlayerID/PlayerName/Symbol/SymbolName/Reward/PoolSize）
##   symbol_xuanwu     — 玄武大獎觸發（全服，PlayerName/Reward，最高優先）
##   symbol_pool_update — 大獎池更新（全服，PoolSize）

extends CanvasLayer

const COLOR_QINGLONG = Color(0.0,   0.667, 0.0)    # #00AA00 青龍綠
const COLOR_BAIHU    = Color(0.753, 0.753, 0.753)  # #C0C0C0 白虎銀
const COLOR_ZHUQUE   = Color(1.0,   0.0,   0.0)    # #FF0000 朱雀紅
const COLOR_XUANWU   = Color(0.0,   0.0,   0.502)  # #000080 玄武深藍
const COLOR_GOLD     = Color(1.0,   0.843, 0.0)    # #FFD700 金
const COLOR_WHITE    = Color(1.0,   1.0,   1.0)    # 白

var _banner: Control = null
var _pool_indicator: Control = null
var _pool_label: Label = null

func _ready() -> void:
	layer = 56  # 比 LuckyStarBurst（55）高一層

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"symbol_trigger":
			_on_symbol_trigger(payload)
		"symbol_xuanwu":
			_on_symbol_xuanwu(payload)
		"symbol_pool_update":
			_on_pool_update(payload)

# ── 四象大獎觸發（全服）────────────────────────────────────────────────────────

func _on_symbol_trigger(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var symbol: String = payload.get("symbol", "qinglong")
	var symbol_name: String = payload.get("symbol_name", "青龍")
	var reward: int = payload.get("reward", 0)
	var pool_size: int = payload.get("pool_size", 0)

	# 依四象決定顏色和 emoji
	var tier_color := _get_symbol_color(symbol)
	var tier_emoji := _get_symbol_emoji(symbol)

	# 閃光（依層級強度）
	var flash_times := 1
	if symbol == "zhuque":
		flash_times = 2
	elif symbol == "xuanwu":
		flash_times = 3
	_flash_screen(tier_color, flash_times, 0.55)

	# 頂部橫幅
	_show_banner(
		"%s %s 大獎！" % [tier_emoji, symbol_name],
		"%s 觸發%s大獎！獲得 %d 金幣！大獎池剩餘 %d" % [player_name, symbol_name, reward, pool_size],
		tier_color
	)

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"%s %s 大獎！" % [tier_emoji, symbol_name],
		Vector2(vp_size / 2),
		tier_color,
		44
	)
	_spawn_float_text(
		"+%d 金幣！" % reward,
		Vector2(vp_size.x * 0.5, vp_size.y * 0.6),
		COLOR_GOLD,
		26
	)

	# 結算彈窗（白虎以上才顯示）
	if symbol in ["baihu", "zhuque", "xuanwu"]:
		_show_settle_popup(player_name, symbol_name, tier_emoji, tier_color, reward, pool_size)

	# 5 秒後清除橫幅
	var timer := get_tree().create_timer(5.0)
	timer.timeout.connect(func():
		if is_instance_valid(_banner):
			_banner.queue_free()
			_banner = null
	)

# ── 玄武大獎觸發（全服最高優先）────────────────────────────────────────────────

func _on_symbol_xuanwu(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var reward: int = payload.get("reward", 0)

	# 全螢幕三次強閃光（深藍→金→深藍）
	_flash_screen(COLOR_XUANWU, 1, 0.7)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_GOLD, 1, 0.8)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_XUANWU, 1, 0.7)

	# 全螢幕大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🐢🐢🐢 玄武大獎！",
		Vector2(vp_size.x * 0.5, vp_size.y * 0.25),
		COLOR_XUANWU,
		56
	)
	_spawn_float_text(
		"%s 獲得 %d 金幣！" % [player_name, reward],
		Vector2(vp_size.x * 0.5, vp_size.y * 0.42),
		COLOR_GOLD,
		28
	)
	_spawn_float_text(
		"大獎池已重置！",
		Vector2(vp_size.x * 0.5, vp_size.y * 0.55),
		COLOR_WHITE,
		18
	)

# ── 大獎池更新（全服）──────────────────────────────────────────────────────────

func _on_pool_update(payload: Dictionary) -> void:
	var pool_size: int = payload.get("pool_size", 0)
	_update_pool_indicator(pool_size)

# ── 輔助方法 ────────────────────────────────────────────────────────────────────

func _get_symbol_color(symbol: String) -> Color:
	match symbol:
		"qinglong": return COLOR_QINGLONG
		"baihu":    return COLOR_BAIHU
		"zhuque":   return COLOR_ZHUQUE
		"xuanwu":   return COLOR_XUANWU
	return COLOR_GOLD

func _get_symbol_emoji(symbol: String) -> String:
	match symbol:
		"qinglong": return "🐉"
		"baihu":    return "🐯"
		"zhuque":   return "🦅"
		"xuanwu":   return "🐢"
	return "🏆"

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

func _update_pool_indicator(pool_size: int) -> void:
	if not is_instance_valid(_pool_indicator):
		var vp_size := get_viewport().get_visible_rect().size
		var indicator := PanelContainer.new()
		indicator.position = Vector2(vp_size.x - 148, 240)
		indicator.size = Vector2(140, 56)
		indicator.mouse_filter = Control.MOUSE_FILTER_IGNORE

		var style := StyleBoxFlat.new()
		style.bg_color = Color(0, 0, 0, 0.82)
		style.border_color = COLOR_GOLD
		style.set_border_width_all(2)
		style.set_corner_radius_all(5)
		indicator.add_theme_stylebox_override("panel", style)

		var vbox := VBoxContainer.new()
		vbox.mouse_filter = Control.MOUSE_FILTER_IGNORE
		indicator.add_child(vbox)

		var title_lbl := Label.new()
		title_lbl.text = "🏆 四象大獎池"
		title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
		title_lbl.add_theme_font_size_override("font_size", 11)
		title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		title_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
		vbox.add_child(title_lbl)

		var pool_lbl := Label.new()
		pool_lbl.name = "PoolLabel"
		pool_lbl.add_theme_color_override("font_color", COLOR_WHITE)
		pool_lbl.add_theme_font_size_override("font_size", 14)
		pool_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		pool_lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
		vbox.add_child(pool_lbl)

		add_child(indicator)
		_pool_indicator = indicator
		_pool_label = pool_lbl

	if is_instance_valid(_pool_label):
		_pool_label.text = "%d" % pool_size

func _show_settle_popup(player_name: String, symbol_name: String, emoji: String, color: Color, reward: int, pool_size: int) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var popup := PanelContainer.new()
	popup.position = Vector2(vp_size.x + 10, vp_size.y * 0.3)
	popup.size = Vector2(220, 130)
	popup.mouse_filter = Control.MOUSE_FILTER_IGNORE

	var style := StyleBoxFlat.new()
	style.bg_color = Color(0, 0, 0, 0.92)
	style.border_color = color
	style.set_border_width_all(3)
	style.set_corner_radius_all(8)
	popup.add_theme_stylebox_override("panel", style)

	var vbox := VBoxContainer.new()
	vbox.mouse_filter = Control.MOUSE_FILTER_IGNORE
	popup.add_child(vbox)

	var lines := [
		["%s %s 大獎！" % [emoji, symbol_name], color, 16],
		[player_name, COLOR_WHITE, 12],
		["獎勵：+%d 金幣" % reward, COLOR_GOLD, 14],
		["大獎池剩餘：%d" % pool_size, COLOR_WHITE, 11],
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
	tween.tween_interval(4.0)
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
