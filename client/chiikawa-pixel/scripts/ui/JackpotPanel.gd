## JackpotPanel.gd
## Progressive Jackpot 面板（DAY-048，從 HUD.gd 拆分 DAY-053）
## 顯示三個等級的 Jackpot 累積金額，中獎時全畫面慶祝特效

extends Control

# 由 HUD.gd 在建立後設定
var pixel_font: Font = null

var _jackpot_labels: Dictionary = {}  # level -> Label
var _jackpot_history: Array = []      # 最近 5 筆中獎記錄
const MAX_JACKPOT_HISTORY = 5

## 初始化面板（由 HUD.gd 呼叫）
func setup(font: Font) -> void:
	pixel_font = font
	_build_panel()
	GameManager.jackpot_updated.connect(_on_jackpot_updated)
	GameManager.jackpot_won.connect(_on_jackpot_won)

## 建立面板 UI
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

	# 三個 Jackpot 等級（Mini / Major / Grand）
	var levels = [
		{"key": "mini",  "label": "MINI",  "color": Color(0.6, 0.9, 1.0), "x": 20},
		{"key": "major", "label": "MAJOR", "color": Color(1.0, 0.8, 0.2), "x": 220},
		{"key": "grand", "label": "GRAND", "color": Color(1.0, 0.3, 0.3), "x": 420},
	]

	for lvl in levels:
		var container = Control.new()
		container.position = Vector2(lvl["x"], 2)
		container.size = Vector2(200, 32)
		add_child(container)

		# 等級標籤
		var title = Label.new()
		title.text = lvl["label"]
		title.position = Vector2(0, 2)
		title.size = Vector2(80, 14)
		title.add_theme_font_size_override("font_size", 10)
		title.add_theme_color_override("font_color", lvl["color"])
		if is_instance_valid(pixel_font):
			title.add_theme_font_override("font", pixel_font)
		container.add_child(title)

		# 金額標籤
		var amount_lbl = Label.new()
		amount_lbl.name = "Amount_" + lvl["key"]
		amount_lbl.text = "---"
		amount_lbl.position = Vector2(0, 16)
		amount_lbl.size = Vector2(180, 16)
		amount_lbl.add_theme_font_size_override("font_size", 13)
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

	# 面板高度擴展（加入 ticker 後從 36 增加到 54）
	var bg_node = get_node_or_null("JackpotBG")
	if is_instance_valid(bg_node):
		bg_node.size.y = 54
	size.y = 54

## Jackpot 池更新（每 5 秒收到一次）
func _on_jackpot_updated(data: Dictionary) -> void:
	var levels = ["mini", "major", "grand"]
	for lvl in levels:
		var lbl = _jackpot_labels.get(lvl)
		if is_instance_valid(lbl):
			var amount = data.get(lvl, 0)
			lbl.text = "🪙%d" % amount
			# 脈動動畫（金額越大越明顯）
			if amount > 0:
				var tween = create_tween()
				tween.tween_property(lbl, "modulate:a", 0.6, 0.15)
				tween.tween_property(lbl, "modulate:a", 1.0, 0.15)

## Jackpot 中獎！全畫面慶祝特效
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

## 顯示 Jackpot 慶祝畫面
func _show_jackpot_celebration(level: String, amount: int, winner_name: String, is_self: bool) -> void:
	# 建立全畫面 overlay（掛在 CanvasLayer 的父節點上）
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

	# 等級顏色
	var level_colors = {
		"mini":  Color(0.6, 0.9, 1.0),
		"major": Color(1.0, 0.8, 0.2),
		"grand": Color(1.0, 0.3, 0.3),
	}
	var level_color = level_colors.get(level, Color.WHITE)
	var level_name = level.to_upper()

	# 主標題
	var title = Label.new()
	title.text = "✨ %s JACKPOT ✨" % level_name
	title.position = Vector2(0, 200)
	title.size = Vector2(1280, 80)
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_font_size_override("font_size", 56)
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

	# 螢幕震動（Grand 最強）
	if ScreenShake != null:
		var trauma = {"mini": 0.4, "major": 0.6, "grand": 0.9}.get(level, 0.4)
		ScreenShake.add_trauma(trauma)

	# 觸發 HitEffect 大獎特效（依等級強度不同）
	if HitEffect != null:
		match level:
			"grand":
				HitEffect.spawn_big_win(Vector2(640, 360), 100.0)
				for i in 3:
					var delay_t = get_tree().create_timer(i * 0.4)
					delay_t.timeout.connect(func():
						_spawn_jackpot_coin_rain(level_color, 20)
					)
			"major":
				HitEffect.spawn_big_win(Vector2(640, 360), 50.0)
				_spawn_jackpot_coin_rain(level_color, 14)
				var delay_t = get_tree().create_timer(0.5)
				delay_t.timeout.connect(func():
					_spawn_jackpot_coin_rain(level_color, 10)
				)
			"mini":
				_spawn_jackpot_coin_rain(level_color, 8)

	# 記錄到 Jackpot 歷史 ticker（DAY-049）
	_add_jackpot_history_entry(level, amount, winner_name, is_self)

## 生成 Jackpot 金幣雨特效（DAY-049）
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
	var level_icons = {"mini": "💙", "major": "💛", "grand": "❤️"}
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
		var level_colors = {"mini": Color(0.6, 0.9, 1.0), "major": Color(1.0, 0.8, 0.2), "grand": Color(1.0, 0.4, 0.4)}
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
