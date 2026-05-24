# RapidRespinPanel.gd ??Rapid Respin 閫貊??Ｘ嚗AY-121嚗?
# 璆剔?靘?嚗eflex Gaming Big Game Fishing Rapid Riches嚗?026-05-14嚗?
# Rapid Respin 閫貊??Ｗ??? + ?璈怠??嚗???＊蝷箏???
extends Control

# ?????憿嚗????甈⊥嚗?
const CHAIN_COLORS = [
	Color(0.3, 0.8, 1.0),   # 蝚?甈∴?憭抵?嚗?.0x嚗?
	Color(0.2, 1.0, 0.4),   # 蝚?甈∴?蝬嚗?.5x嚗?
	Color(1.0, 0.8, 0.0),   # 蝚?甈∴??嚗?.0x嚗?
	Color(1.0, 0.4, 0.0),   # 蝚?甈∴?璈?嚗?.0x嚗?
	Color(1.0, 0.2, 0.8),   # 蝚?甈∴?蝎換嚗?.0x嚗?
]

# ?????璅惜
const CHAIN_LABELS = ["??RAPID RESPIN", "? CHAIN x2", "? CHAIN x3", "?? CHAIN x4", "? MAX CHAIN x5"]

var _banner: Control = null
var _flash_overlay: ColorRect = null

func _ready():
	# 撱箇??刻撟??蝵抬??身?梯?嚗?
	_flash_overlay = ColorRect.new()
	_flash_overlay.size = Vector2(1280, 720)
	_flash_overlay.position = Vector2.ZERO
	_flash_overlay.color = Color(0.3, 0.8, 1.0, 0.0)
	_flash_overlay.z_index = 70
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# ?? GameManager 閮?
	if GameManager.has_signal("rapid_respin"):
		GameManager.rapid_respin.connect(_on_rapid_respin)
	if GameManager.has_signal("rapid_respin_end"):
		GameManager.rapid_respin_end.connect(_on_rapid_respin_end)

func _on_rapid_respin(data: Dictionary) -> void:
	var player_name = data.get("player_name", "?拙振")
	var chain_count = data.get("chain_count", 0)
	var chain_mult = data.get("chain_mult", 1.0)
	var is_chain = data.get("is_chain", false)
	var icon = data.get("icon", "")"
	var player_id = data.get("player_id", "")

	var is_self = (player_id == GameManager.player_data.get("player_id", ""))

	_show_respin_effect(player_name, chain_count, chain_mult, is_chain, icon, is_self)

func _on_rapid_respin_end(data: Dictionary) -> void:
	var player_name = data.get("player_name", "?拙振")
	var total_chain = data.get("total_chain", 1)
	var player_id = data.get("player_id", "")
	var is_self = (player_id == GameManager.player_data.get("player_id", ""))

	if is_self and total_chain >= 2:
		_show_chain_end_banner(total_chain)

func _show_respin_effect(player_name: String, chain_count: int, chain_mult: float,
		is_chain: bool, icon: String, is_self: bool) -> void:

	var color_idx = clamp(chain_count, 0, CHAIN_COLORS.size() - 1)
	var color = CHAIN_COLORS[color_idx]
	var label_text = CHAIN_LABELS[color_idx] if chain_count < CHAIN_LABELS.size() else CHAIN_LABELS[-1]

	# ?刻撟?????
	_flash_overlay.color = Color(color.r, color.g, color.b, 0.0)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color:a", 0.35, 0.08)
	flash_tween.tween_property(_flash_overlay, "color:a", 0.0, 0.25)

	# 蝘駁?帖撟?
	if is_instance_valid(_banner):
		_banner.queue_free()

	# 撱箇??璈怠?
	_banner = Control.new()
	_banner.z_index = 71
	add_child(_banner)

	# 璈怠??
	var bg = ColorRect.new()
	bg.size = Vector2(1280, 64)
	bg.position = Vector2(0, -64)  # 敺??典???
	bg.color = Color(color.r * 0.15, color.g * 0.15, color.b * 0.15, 0.95)
	_banner.add_child(bg)

	# ?敶抵??
	var top_bar = ColorRect.new()
	top_bar.size = Vector2(1280, 4)
	top_bar.position = Vector2(0, 0)
	top_bar.color = color
	bg.add_child(top_bar)

	# 摨敶抵??
	var bot_bar = ColorRect.new()
	bot_bar.size = Vector2(1280, 4)
	bot_bar.position = Vector2(0, 60)
	bot_bar.color = color
	bg.add_child(bot_bar)

	# 銝餅?憿????憿?嚗?
	var title_lbl = Label.new()
	title_lbl.text = label_text
	title_lbl.position = Vector2(40, 8)
	title_lbl.add_theme_font_size_override("font_size", 28)
	title_lbl.add_theme_color_override("font_color", color)
	bg.add_child(title_lbl)

	# ??憿舐內
	var mult_lbl = Label.new()
	mult_lbl.text = "?%.1f" % chain_mult
	mult_lbl.position = Vector2(400, 8)
	mult_lbl.add_theme_font_size_override("font_size", 28)
	mult_lbl.add_theme_color_override("font_color", Color.WHITE)
	bg.add_child(mult_lbl)

	# ?拙振?迂
	var name_lbl = Label.new()
	if is_self:
		name_lbl.text = "雿孛?潔?嚗?"
		name_lbl.add_theme_color_override("font_color", Color(1.0, 1.0, 0.3))
	else:
		name_lbl.text = player_name + " 閫貊"
		name_lbl.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	name_lbl.position = Vector2(700, 8)
	name_lbl.add_theme_font_size_override("font_size", 22)
	bg.add_child(name_lbl)

	# ????脣漲?內?剁?撠?暺?
	for i in range(5):
		var dot = ColorRect.new()
		dot.size = Vector2(12, 12)
		dot.position = Vector2(1100 + i * 20, 26)
		if i <= chain_count:
			dot.color = color
		else:
			dot.color = Color(0.3, 0.3, 0.3, 0.8)
		bg.add_child(dot)

	# 皛?
	var slide_tween = create_tween()
	slide_tween.tween_property(bg, "position:y", 0.0, 0.15).set_ease(Tween.EASE_OUT)

	# ?芸?皛嚗? 蝘?嚗?
	var timer = get_tree().create_timer(3.0)
	timer.timeout.connect(func():
		if is_instance_valid(bg):
			var out_tween = create_tween()
			out_tween.tween_property(bg, "position:y", -64.0, 0.2).set_ease(Tween.EASE_IN)
			out_tween.tween_callback(func():
				if is_instance_valid(_banner):
					_banner.queue_free()
					_banner = null
			)
	)

	# ?芸楛閫貊??憭＊蝷粹??脩?摮???
	if is_self:
		_spawn_star_particles(color)

func _show_chain_end_banner(total_chain: int) -> void:
	# ???蝯??＊蝷箇蜇蝯帖撟?
	var end_lbl = Label.new()
	end_lbl.text = "?? Rapid Respin ???蝯?嚗 %d 甈? % total_chain"
	end_lbl.position = Vector2(400, 680)
	end_lbl.add_theme_font_size_override("font_size", 20)
	end_lbl.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	end_lbl.z_index = 71
	add_child(end_lbl)

	var tween = create_tween()
	tween.tween_property(end_lbl, "modulate:a", 0.0, 2.0).set_delay(1.5)
	tween.tween_callback(end_lbl.queue_free)

func _spawn_star_particles(color: Color) -> void:
	# ?函?Ｖ葉憭桃???8 ??敶Ｙ?摮?
	for i in range(8):
		var star = Label.new()
		star.text = "??"
		star.add_theme_font_size_override("font_size", 24)
		star.add_theme_color_override("font_color", color)
		star.z_index = 72

		var angle = i * PI / 4.0
		var start_x = 640.0
		var start_y = 360.0
		star.position = Vector2(start_x, start_y)
		add_child(star)

		var end_x = start_x + cos(angle) * 200.0
		var end_y = start_y + sin(angle) * 200.0

		var tween = create_tween()
		tween.set_parallel(true)
		tween.tween_property(star, "position", Vector2(end_x, end_y), 0.6).set_ease(Tween.EASE_OUT)
		tween.tween_property(star, "modulate:a", 0.0, 0.6).set_delay(0.2)
		tween.tween_callback(star.queue_free).set_delay(0.6)
