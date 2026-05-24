## JackpotPanel.gd
## Progressive Jackpot ?Ｘ嚗AY-048嚗AY-095 ???惜 + ??嚗AY-118 Meter ?脣漲璇?
## 憿舐內??蝝? Jackpot 蝝舐??? + ?澆??脣漲璇?銝剔???恍?嗥??寞?

extends Control

# ??HUD.gd ?典遣蝡?閮剖?
var pixel_font: Font = null

var _jackpot_labels: Dictionary = {}  # level -> Label
var _jackpot_meters: Dictionary = {}  # level -> ColorRect嚗脣漲璇‵??
var _jackpot_meter_mats: Dictionary = {} # level -> ShaderMaterial
var _jackpot_history: Array = []      # ?餈?5 蝑葉????
const MAX_JACKPOT_HISTORY = 5

# Jackpot ?瑼鳴?撠? Server 蝡?jackpot.go嚗?
const JACKPOT_THRESHOLDS = {
	"mini":  300,
	"minor": 1000,
	"major": 3000,
	"grand": 15000,
}

# ?惜蝑?摰儔嚗AY-095嚗?
const JACKPOT_LEVELS = [
	{"key": "mini",  "label": "MINI",  "color": Color(0.75, 0.75, 0.75), "icon": "??", "x": 0},
	{"key": "minor", "label": "MINOR", "color": Color(1.0, 0.85, 0.2),   "icon": "??", "x": 160},
	{"key": "major", "label": "MAJOR", "color": Color(1.0, 0.5, 0.1),    "icon": "?", "x": 320},
	{"key": "grand", "label": "GRAND", "color": Color(1.0, 0.2, 0.6),    "icon": "??", "x": 480},
]

const METER_SHADER_PATH = "res://assets/shaders/jackpot_meter.gdshader"

## ????選???HUD.gd ?澆嚗?
func setup(font: Font) -> void:
	pixel_font = font
	_build_panel()
	GameManager.jackpot_updated.connect(_on_jackpot_updated)
	GameManager.jackpot_won.connect(_on_jackpot_won)
	GameManager.jackpot_animation.connect(_on_jackpot_animation)

## 撱箇??Ｘ UI嚗?撅斤??穿?
func _build_panel() -> void:
	# ?嚗楛?脣???嚗葆???嚗?
	var bg = ColorRect.new()
	bg.name = "JackpotBG"
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.05, 0.03, 0.12, 0.85)
	add_child(bg)

	# ????
	var top_line = ColorRect.new()
	top_line.size = Vector2(640, 2)
	top_line.position = Vector2(0, 0)
	top_line.color = Color(0.90, 0.75, 0.20, 0.80)
	add_child(top_line)

	# ??Jackpot 蝑?嚗ini / Minor / Major / Grand嚗?
	for lvl in JACKPOT_LEVELS:
		var container = Control.new()
		container.position = Vector2(lvl["x"], 2)
		container.size = Vector2(160, 44)  # DAY-118嚗?摨血? 32 ??44嚗??脣漲璇?
		add_child(container)

		# 蝑?璅惜嚗?內嚗?
		var title = Label.new()
		title.text = "%s %s" % [lvl["icon"], lvl["label"]]
		title.position = Vector2(0, 2)
		title.size = Vector2(155, 14)
		title.add_theme_font_size_override("font_size", 9)
		title.add_theme_color_override("font_color", lvl["color"])
		if is_instance_valid(pixel_font):
			title.add_theme_font_override("font", pixel_font)
		container.add_child(title)

		# ??璅惜
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

		# ?? DAY-118嚗ackpot Meter ?脣漲璇???
		# ?脣漲璇???
		var meter_bg = ColorRect.new()
		meter_bg.name = "MeterBG_" + lvl["key"]
		meter_bg.position = Vector2(2, 34)
		meter_bg.size = Vector2(154, 8)
		meter_bg.color = Color(0.05, 0.05, 0.1, 0.9)
		container.add_child(meter_bg)

		# ?脣漲璇‵??撣?Shader嚗?
		var meter_fill = ColorRect.new()
		meter_fill.name = "MeterFill_" + lvl["key"]
		meter_fill.position = Vector2(2, 34)
		meter_fill.size = Vector2(0, 8)  # ??撖砍漲 0嚗 _on_jackpot_updated ?湔
		meter_fill.color = lvl["color"]
		container.add_child(meter_fill)
		_jackpot_meters[lvl["key"]] = meter_fill

		# 憟 Jackpot Meter Shader嚗????剁?
		if ResourceLoader.exists(METER_SHADER_PATH):
			var mat = ShaderMaterial.new()
			mat.shader = load(METER_SHADER_PATH)
			mat.set_shader_parameter("bar_color", lvl["color"])
			mat.set_shader_parameter("fill_ratio", 0.0)
			mat.set_shader_parameter("glow_intensity", 0.8)
			mat.set_shader_parameter("time_offset", lvl["x"] * 0.01)  # ??蝝???蝘颱???
			meter_fill.material = mat
			_jackpot_meter_mats[lvl["key"]] = mat

	# Jackpot 甇瑕 ticker嚗AY-049嚗?憿舐內?餈葉????
	var ticker_bg = ColorRect.new()
	ticker_bg.name = "TickerBG"
	ticker_bg.position = Vector2(0, 48)  # DAY-118嚗?蝵桐?蝘?
	ticker_bg.size = Vector2(640, 18)
	ticker_bg.color = Color(0.02, 0.01, 0.08, 0.75)
	add_child(ticker_bg)

	var ticker_lbl = Label.new()
	ticker_lbl.name = "TickerLabel"
	ticker_lbl.text = "??蝑? Jackpot 銝剔?..."
	ticker_lbl.position = Vector2(8, 50)  # DAY-118嚗?蝵桐?蝘?
	ticker_lbl.size = Vector2(624, 16)
	ticker_lbl.add_theme_font_size_override("font_size", 10)
	ticker_lbl.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 0.7))
	if is_instance_valid(pixel_font):
		ticker_lbl.add_theme_font_override("font", pixel_font)
	add_child(ticker_lbl)

	# ?Ｘ擃漲嚗AY-118嚗? 54 ??66嚗?
	var bg_node = get_node_or_null("JackpotBG")
	if is_instance_valid(bg_node):
		bg_node.size.y = 66
	size.y = 66

## Jackpot 瘙?堆?瘥?5 蝘?唬?甈∴??惜?嚗?
func _on_jackpot_updated(data: Dictionary) -> void:
	for lvl in JACKPOT_LEVELS:
		var lbl = _jackpot_labels.get(lvl["key"])
		if is_instance_valid(lbl):
			var amount = data.get(lvl["key"], 0)
			lbl.text = "??%d" % amount
			# Grand ??憭扳???????
			if lvl["key"] == "grand" and amount > 5000:
				var tween = create_tween()
				tween.tween_property(lbl, "modulate:a", 0.5, 0.2)
				tween.tween_property(lbl, "modulate:a", 1.0, 0.2)

		# ?? DAY-118嚗??Jackpot Meter ?脣漲璇???
		var meter = _jackpot_meters.get(lvl["key"])
		if is_instance_valid(meter):
			var amount = data.get(lvl["key"], 0)
			# ?芸?雿輻 Server ?喃???瑼鳴???砍撣豢
			var threshold_key = lvl["key"] + "_threshold"
			var threshold = data.get(threshold_key, JACKPOT_THRESHOLDS.get(lvl["key"], 1000))
			var ratio = clamp(float(amount) / float(threshold), 0.0, 1.0)
			var target_width = 154.0 * ratio

			# 撟單???湔?脣漲璇祝摨?
			var tween = create_tween()
			tween.tween_property(meter, "size:x", target_width, 0.4).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)

			# ?湔 Shader ?
			var mat = _jackpot_meter_mats.get(lvl["key"])
			if is_instance_valid(mat):
				mat.set_shader_parameter("fill_ratio", ratio)
				# ?亥?閫貊??>80%嚗?撘瑞??
				var glow = 0.8 + ratio * 1.5
				mat.set_shader_parameter("glow_intensity", glow)

			# ?亥?閫貊??>90%嚗???璅惜??
			if ratio >= 0.9 and is_instance_valid(lbl):
				var flash_tween = create_tween().set_loops(2)
				flash_tween.tween_property(lbl, "modulate", Color(2.0, 2.0, 2.0), 0.12)
				flash_tween.tween_property(lbl, "modulate", Color.WHITE, 0.12)

## Jackpot 閫貊??嚗AY-095嚗?撱?蝯行??摰?
func _on_jackpot_animation(data: Dictionary) -> void:
	var level = data.get("level", "mini")
	var level_name = data.get("level_name", level.to_upper())
	var level_color_hex = data.get("level_color", "#FFFFFF")
	var amount = data.get("amount", 0)
	var winner_name = data.get("winner_name", "")
	var is_grand = data.get("is_grand", false)
	var is_major = data.get("is_major", false)

	# 閫??憿
	var level_color = Color.WHITE
	if level_color_hex.begins_with("#"):
		level_color = Color(level_color_hex)

	# 靘?蝝孛?潔??撥摨衣??
	if is_grand:
		_show_grand_jackpot_animation(level_name, amount, winner_name, level_color)
	elif is_major:
		_show_major_jackpot_animation(level_name, amount, winner_name, level_color)
	else:
		_show_mini_jackpot_animation(level_name, amount, winner_name, level_color)

## Grand Jackpot ?函?Ｗ??恬??撘瑞??
func _show_grand_jackpot_animation(level_name: String, amount: int, winner_name: String, color: Color) -> void:
	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# ?函?ａ??脤???
	var flash = ColorRect.new()
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(1.0, 0.9, 0.2, 0.0)
	flash.z_index = 195
	canvas_layer.add_child(flash)

	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", 0.6, 0.1)
	flash_tween.tween_property(flash, "color:a", 0.0, 0.3)
	flash_tween.tween_callback(flash.queue_free)

	# ?Ｗ???
	if ScreenShake != null:
		ScreenShake.add_trauma(0.9)

	# 銝郭?馳??
	for i in 3:
		var timer = get_tree().create_timer(i * 0.35)
		timer.timeout.connect(func():
			_spawn_jackpot_coin_rain(color, 25)
		)

	# 憭抒??寞?
	if HitEffect != null:
		HitEffect.spawn_big_win(Vector2(640, 360), 100.0)

## Major Jackpot ??Ｗ???
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

## Mini/Minor Jackpot 撠???
func _show_mini_jackpot_animation(level_name: String, amount: int, winner_name: String, color: Color) -> void:
	_spawn_jackpot_coin_rain(color, 8)

## Jackpot 銝剔?嚗＊蝷箸蟡??
func _on_jackpot_won(data: Dictionary) -> void:
	var level = data.get("level", "mini")
	var amount = data.get("amount", 0)
	var winner_name = data.get("winner_name", "")
	var is_self = data.get("winner_id", "") == NetworkManager.get_player_id()

	# ?剜憭抒??單?
	if AudioManager != null:
		AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)

	# ?函?Ｘ蟡?overlay
	_show_jackpot_celebration(level, amount, winner_name, is_self)

## 憿舐內 Jackpot ?嗥??恍嚗?撅斤??穿?
func _show_jackpot_celebration(level: String, amount: int, winner_name: String, is_self: bool) -> void:
	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	var overlay = Control.new()
	overlay.name = "JackpotCelebration"
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.z_index = 200
	canvas_layer.add_child(overlay)

	# ??暺?
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.0, 0.0, 0.0, 0.0)
	overlay.add_child(bg)

	# 蝑?憿嚗?撅歹?
	var level_color = Color.WHITE
	var level_icon = "??"
	for lvl in JACKPOT_LEVELS:
		if lvl["key"] == level:
			level_color = lvl["color"]
			level_icon = lvl["icon"]
			break
	var level_name = level.to_upper()

	# 銝餅?憿?
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

	# ??
	var amount_lbl = Label.new()
	amount_lbl.text = "?? %d" % amount
	amount_lbl.position = Vector2(0, 290)
	amount_lbl.size = Vector2(1280, 60)
	amount_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	amount_lbl.add_theme_font_size_override("font_size", 44)
	amount_lbl.add_theme_color_override("font_color", Color(1.0, 0.95, 0.5))
	if is_instance_valid(pixel_font):
		amount_lbl.add_theme_font_override("font", pixel_font)
	overlay.add_child(amount_lbl)

	# 銝剔???蝔?
	var winner_text = ("?? YOU WIN!" if is_self else "?? %s WINS!" % winner_name)
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

	# ?嚗??舀楚????璅?敶 ???? ??瘛∪
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

	# 閮???Jackpot 甇瑕 ticker
	_add_jackpot_history_entry(level, amount, winner_name, is_self)

## ?? Jackpot ?馳?函??
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

## ?銝蝑?Jackpot 銝剔?閮?銝行??ticker
func _add_jackpot_history_entry(level: String, amount: int, winner_name: String, is_self: bool) -> void:
	var level_icons = {"mini": "??", "minor": "??", "major": "?", "grand": "??"}
	var icon = level_icons.get(level, "??)"
	var name_display = "YOU" if is_self else winner_name
	var entry_text = "%s %s: %s ??%d" % [icon, level.to_upper(), name_display, amount]

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

	# ???
	var tween = create_tween()
	tween.tween_property(ticker_lbl, "modulate:a", 0.3, 0.1)
	tween.tween_property(ticker_lbl, "modulate:a", 1.0, 0.1)
	tween.tween_property(ticker_lbl, "modulate:a", 0.3, 0.1)
	tween.tween_property(ticker_lbl, "modulate:a", 1.0, 0.1)

	# 5 蝘????唬?銝蝑?
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
