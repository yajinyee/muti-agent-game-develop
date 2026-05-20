## JackpotPanel.gd
## Progressive Jackpot 面板（DAY-048，DAY-095 升級四層 + 動畫通知）
## 顯示四個等級的 Jackpot 累積金額，中獎時全畫面慶祝特效

extends Control

# 由 HUD.gd 在建立後設定
var pixel_font: Font = null

var _jackpot_labels: Dictionary = {}  # level -> Label
var _jackpot_history: Array = []      # 最近 5 筆中獎記錄
const MAX_JACKPOT_HISTORY = 5

# 四層等級定義（DAY-095）
const JACKPOT_LEVELS = [
	{"key": "mini",  "label": "MINI",  "color": Color(0.75, 0.75, 0.75), "icon": "🥈", "x": 0},
	{"key": "minor", "label": "MINOR", "color": Color(1.0, 0.85, 0.2),   "icon": "🥇", "x": 160},
	{"key": "major", "label": "MAJOR", "color": Color(1.0, 0.5, 0.1),    "icon": "🔥", "x": 320},
	{"key": "grand", "label": "GRAND", "color": Color(1.0, 0.2, 0.6),    "icon": "👑", "x": 480},
]

## 初始化面板（由 HUD.gd 呼叫）
func setup(font: Font) -> void:
	pixel_font = font
	_build_panel()
	GameManager.jackpot_updated.connect(_on_jackpot_updated)
	GameManager.jackpot_won.connect(_on_jackpot_won)
	GameManager.jackpot_animation.connect(_on_jackpot_animation)

## 建立面板 UI（四層版本）
func _build_panel() -> void:
	# 背景（深色半透明，帶金色邊框）
	var bg = ColorRect.new()
	bg.name = "JackpotBG"
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.05, 0.03, 0.12, 0.85)
	add_child(bg)

	# 金色頂部邊框
	var top_line = ColorRect.new()
	top_line.size = Vector2(640, 2)
	top_line.position = Vector2(0, 0)
	top_line.color = Color(0.90, 0.75, 0.20, 0.80)
	add_child(top_line)

	# 四個 Jackpot 等級（Mini / Minor / Major / Grand）
	for lvl in JACKPOT_LEVELS:
		var container = Control.new()
		container.position = Vector2(lvl["x"], 2)
		container.size = Vector2(160, 32)
		add_child(container)

		# 等級標籤（含圖示）
		var title = Label.new()
		title.text = "%s %s" % [lvl["icon"], lvl["label"]]
		title.position = Vector2(0, 2)
		title.size = Vector2(155, 14)
		title.add_theme_font_size_override("font_size", 9)
		title.add_theme_color_override("font_color", lvl["color"])
		if is_instance_valid(pixel_font):
			title.add_theme_font_override("font", pixel_font)
		container.add_child(title)

		# 金額標籤
		var amount_lbl = Label.new()
		amount_lbl.name = "Amount_" + lvl["key"]
		amount_lbl.text = "---"
		amount_lbl.position = Vector2(0, 16)
		amount_lbl.size = Vector2(155, 16)
		amount_lbl.add_theme_font_size_override("font_size", 12)
		amount_lbl.add_theme_color_override("font_color", Color(1.0, 0.95, 0.7))
		if is_instance_valid(pixel_font):
			amount_lbl.add_theme_font_override("font", pixel_font)
		container.add_child(amount_lbl)
		_jackpot_labels[lvl["key"]] = amount_lbl

	# Jackpot 歷史 ticker（DAY-049）— 顯示最近中獎記錄
	var ticker_bg = ColorRect.new()
	ticker_bg.name = "TickerBG"
	ticker_bg.position = Vector2(0, 36)
	ticker_bg.size = Vector2(640, 18)
	ticker_bg.color = Color(0.02, 0.01, 0.08, 0.75)
	add_child(ticker_bg)

	var ticker_lbl = Label.new()
	ticker_lbl.name = "TickerLabel"
	ticker_lbl.text = "✨ 等待 Jackpot 中獎..."
	ticker_lbl.position = Vector2(8, 38)
	ticker_lbl.size = Vector2(624, 16)
	ticker_lbl.add_theme_font_size_override("font_size", 10)
	ticker_lbl.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 0.7))
	if is_instance_valid(pixel_font):
		ticker_lbl.add_theme_font_override("font", pixel_font)
	add_child(ticker_lbl)

	# 面板高度
	var bg_node = get_node_or_null("JackpotBG")
	if is_instance_valid(bg_node):
		bg_node.size.y = 54
	size.y = 54

## Jackpot 池更新（每 5 秒收到一次，四層版本）
func _on_jackpot_updated(data: Dictionary) -> void:
	for lvl in JACKPOT_LEVELS:
		var lbl = _jackpot_labels.get(lvl["key"])
		if is_instance_valid(lbl):
			var amount = data.get(lvl["key"], 0)
			lbl.text = "🪙%d" % amount
			# Grand 金額大時加閃爍效果
			if lvl["key"] == "grand" and amount > 5000:
				var tween = create_tween()
				tween.tween_property(lbl, "modulate:a", 0.5, 0.2)
				tween.tween_property(lbl, "modulate:a", 1.0, 0.2)

## Jackpot 觸發動畫通知（DAY-095）— 廣播給所有玩家
func _on_jackpot_animation(data: Dictionary) -> void:
	var level = data.get("level", "mini")
	var level_name = data.get("level_name", level.to_upper())
	var level_color_hex = data.get("level_color", "#FFFFFF")
	var amount = data.get("amount", 0)
	var winner_name = data.get("winner_name", "")
	var is_grand = data.get("is_grand", false)
	var is_major = data.get("is_major", false)

	# 解析顏色
	var level_color = Color.WHITE
	if level_color_hex.begins_with("#"):
		level_color = Color(level_color_hex)

	# 依等級觸發不同強度的動畫
	if is_grand:
		_show_grand_jackpot_animation(level_name, amount, winner_name, level_color)
	elif is_major:
		_show_major_jackpot_animation(level_name, amount, winner_name, level_color)
	else:
		_show_mini_jackpot_animation(level_name, amount, winner_name, level_color)

## Grand Jackpot 全畫面動畫（最強特效）
func _show_grand_jackpot_animation(level_name: String, amount: int, winner_name: String, color: Color) -> void:
	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# 全畫面金色閃光
	var flash = ColorRect.new()
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(1.0, 0.9, 0.2, 0.0)
	flash.z_index = 195
	canvas_layer.add_child(flash)

	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", 0.6, 0.1)
	flash_tween.tween_property(flash, "color:a", 0.0, 0.3)
	flash_tween.tween_callback(flash.queue_free)

	# 螢幕震動
	if ScreenShake != null:
		ScreenShake.add_trauma(0.9)

	# 三波金幣雨
	for i in 3:
		var timer = get_tree().create_timer(i * 0.35)
		timer.timeout.connect(func():
			_spawn_jackpot_coin_rain(color, 25)
		)

	# 大獎特效
	if HitEffect != null:
		HitEffect.spawn_big_win(Vector2(640, 360), 100.0)

## Major Jackpot 半畫面動畫
func _show_major_jackpot_animation(level_name: String, amount: int, winner_name: String, color: Color) -> void:
	if ScreenShake != null:
		ScreenShake.add_trauma(0.6)
	_spawn_jackpot_coin_rain(color, 16)
	var timer = get_tree().create_timer(0.4)
	timer.timeout.connect(func():
		_spawn_jackpot_coin_rain(color, 12)
	)
	if HitEffect != null:
		HitEffect.spawn_big_win(Vector2(640, 360), 50.0)

## Mini/Minor Jackpot 小動畫
func _show_mini_jackpot_animation(level_name: String, amount: int, winner_name: String, color: Color) -> void:
	_spawn_jackpot_coin_rain(color, 8)

## Jackpot 中獎！顯示慶祝面板
func _on_jackpot_won(data: Dictionary) -> void:
	var level = data.get("level", "mini")
	var amount = data.get("amount", 0)
	var winner_name = data.get("winner_name", "")
	var is_self = data.get("winner_id", "") == NetworkManager.get_player_id()

	# 播放大獎音效
	if AudioManager != null:
		AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)

	# 全畫面慶祝 overlay
	_show_jackpot_celebration(level, amount, winner_name, is_self)

## 顯示 Jackpot 慶祝畫面（四層版本）
func _show_jackpot_celebration(level: String, amount: int, winner_name: String, is_self: bool) -> void:
	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	var overlay = Control.new()
	overlay.name = "JackpotCelebration"
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.z_index = 200
	canvas_layer.add_child(overlay)

	# 半透明黑色背景
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.0, 0.0, 0.0, 0.0)
	overlay.add_child(bg)

	# 等級顏色（四層）
	var level_color = Color.WHITE
	var level_icon = "✨"
	for lvl in JACKPOT_LEVELS:
		if lvl["key"] == level:
			level_color = lvl["color"]
			level_icon = lvl["icon"]
			break
	var level_name = level.to_upper()

	# 主標題
	var title = Label.new()
	title.text = "%s %s JACKPOT %s" % [level_icon, level_name, level_icon]
	title.position = Vector2(0, 200)
	title.size = Vector2(1280, 80)
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_font_size_override("font_size", 52)
	title.add_theme_color_override("font_color", level_color)
	title.add_theme_color_override("font_shadow_color", Color(0.0, 0.0, 0.0, 0.9))
	title.add_theme_constant_override("shadow_offset_x", 4)
	title.add_theme_constant_override("shadow_offset_y", 4)
	if is_instance_valid(pixel_font):
		title.add_theme_font_override("font", pixel_font)
	overlay.add_child(title)

	# 金額
	var amount_lbl = Label.new()
	amount_lbl.text = "🪙 %d" % amount
	amount_lbl.position = Vector2(0, 290)
	amount_lbl.size = Vector2(1280, 60)
	amount_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	amount_lbl.add_theme_font_size_override("font_size", 44)
	amount_lbl.add_theme_color_override("font_color", Color(1.0, 0.95, 0.5))
	if is_instance_valid(pixel_font):
		amount_lbl.add_theme_font_override("font", pixel_font)
	overlay.add_child(amount_lbl)

	# 中獎者名稱
	var winner_text = ("🎉 YOU WIN!" if is_self else "🎉 %s WINS!" % winner_name)
	var winner_lbl = Label.new()
	winner_lbl.text = winner_text
	winner_lbl.position = Vector2(0, 360)
	winner_lbl.size = Vector2(1280, 40)
	winner_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	winner_lbl.add_theme_font_size_override("font_size", 28)
	winner_lbl.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	if is_instance_valid(pixel_font):
		winner_lbl.add_theme_font_override("font", pixel_font)
	overlay.add_child(winner_lbl)

	# 動畫：背景淡入 → 標題彈入 → 停留 → 淡出
	var tween = create_tween()
	tween.tween_property(bg, "color", Color(0.0, 0.0, 0.0, 0.75), 0.3)
	title.position.y = 400
	title.modulate.a = 0.0
	tween.tween_property(title, "position:y", 200.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(title, "modulate:a", 1.0, 0.3)
	amount_lbl.modulate.a = 0.0
	tween.tween_property(amount_lbl, "modulate:a", 1.0, 0.3)
	winner_lbl.modulate.a = 0.0
	tween.tween_property(winner_lbl, "modulate:a", 1.0, 0.3)
	tween.tween_interval(3.0)
	tween.tween_property(overlay, "modulate:a", 0.0, 0.5)
	tween.tween_callback(overlay.queue_free)

	# 記錄到 Jackpot 歷史 ticker
	_add_jackpot_history_entry(level, amount, winner_name, is_self)

## 生成 Jackpot 金幣雨特效
func _spawn_jackpot_coin_rain(color: Color, count: int) -> void:
	var rng = RandomNumberGenerator.new()
	rng.randomize()
	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return
	for i in count:
		var coin = ColorRect.new()
		coin.size = Vector2(8, 8)
		coin.color = color
		coin.position = Vector2(rng.randf_range(100, 1180), -20)
		coin.z_index = 190
		canvas_layer.add_child(coin)

		var target_y = rng.randf_range(200, 700)
		var target_x = coin.position.x + rng.randf_range(-80, 80)
		var duration = rng.randf_range(0.6, 1.2)

		var tween = coin.create_tween()
		tween.tween_property(coin, "position", Vector2(target_x, target_y), duration).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
		tween.parallel().tween_property(coin, "rotation", rng.randf_range(-PI, PI), duration)
		tween.tween_property(coin, "modulate:a", 0.0, 0.3)
		tween.tween_callback(coin.queue_free)

## 加入一筆 Jackpot 中獎記錄並更新 ticker
func _add_jackpot_history_entry(level: String, amount: int, winner_name: String, is_self: bool) -> void:
	var level_icons = {"mini": "🥈", "minor": "🥇", "major": "🔥", "grand": "👑"}
	var icon = level_icons.get(level, "✨")
	var name_display = "YOU" if is_self else winner_name
	var entry_text = "%s %s: %s 🪙%d" % [icon, level.to_upper(), name_display, amount]

	_jackpot_history.insert(0, entry_text)
	if _jackpot_history.size() > MAX_JACKPOT_HISTORY:
		_jackpot_history.resize(MAX_JACKPOT_HISTORY)

	var ticker_lbl = get_node_or_null("TickerLabel")
	if not is_instance_valid(ticker_lbl):
		return

	ticker_lbl.text = entry_text
	if is_self:
		ticker_lbl.add_theme_color_override("font_color", Color(1.0, 0.95, 0.3, 1.0))
	else:
		var level_colors = {
			"mini":  Color(0.75, 0.75, 0.75),
			"minor": Color(1.0, 0.85, 0.2),
			"major": Color(1.0, 0.5, 0.1),
			"grand": Color(1.0, 0.2, 0.6)
		}
		ticker_lbl.add_theme_color_override("font_color", level_colors.get(level, Color.WHITE))

	# 閃爍動畫
	var tween = create_tween()
	tween.tween_property(ticker_lbl, "modulate:a", 0.3, 0.1)
	tween.tween_property(ticker_lbl, "modulate:a", 1.0, 0.1)
	tween.tween_property(ticker_lbl, "modulate:a", 0.3, 0.1)
	tween.tween_property(ticker_lbl, "modulate:a", 1.0, 0.1)

	# 5 秒後切換到下一筆
	if _jackpot_history.size() > 1:
		var timer = get_tree().create_timer(5.0)
		timer.timeout.connect(func():
			if is_instance_valid(ticker_lbl) and _jackpot_history.size() > 1:
				var cur_idx = ticker_lbl.get_meta("ticker_idx", 0)
				var next_idx = (cur_idx + 1) % _jackpot_history.size()
				ticker_lbl.set_meta("ticker_idx", next_idx)
				ticker_lbl.text = _jackpot_history[next_idx]
				ticker_lbl.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7, 0.6))
		)
