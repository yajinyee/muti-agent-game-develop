## LuckyChainMeteorPanel.gd — T129 幸運連鎖隕石魚 UI
## lucky-panel-agent 負責維護
## 業界依據：Royal Fishing Jili「Dragon Wrath meteors — accumulate wrath, unleash meteorite attack」
## 視覺主題：火橙色 + 隕石計數器 + AOE 半徑指示器 + 完美隕石雨演出
extends CanvasLayer

const LAYER_Z = 29  # CanvasLayer layer 值

# 顏色主題
const COLOR_METEOR   = Color(1.0, 0.5, 0.1)    # 火橙（隕石主色）
const COLOR_PERFECT  = Color(1.0, 0.85, 0.0)   # 金色（完美隕石雨）
const COLOR_MISS     = Color(0.5, 0.5, 0.5)    # 灰色（空揮）
const COLOR_BG       = Color(0.1, 0.03, 0.0, 0.88)

var _banner: Control = null
var _meteor_panel: Control = null
var _result_popup: Control = null
var _flash_overlay: ColorRect = null
var _meteor_dots: Array = []  # 5 個隕石指示點
var _radius_label: Label = null
var _hit_label: Label = null

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_meteor_panel()
	GameManager.lucky_chain_meteor.connect(_on_lucky_chain_meteor)

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.5, 0.1, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

func _create_meteor_panel() -> void:
	# 右上角隕石進度顯示
	_meteor_panel = Control.new()
	_meteor_panel.position = Vector2(960, 10)
	_meteor_panel.size = Vector2(300, 100)
	_meteor_panel.visible = false
	add_child(_meteor_panel)

	var bg = ColorRect.new()
	bg.size = _meteor_panel.size
	bg.color = COLOR_BG
	_meteor_panel.add_child(bg)

	var title = Label.new()
	title.text = "☄️ 連鎖隕石雨"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_METEOR
	_meteor_panel.add_child(title)

	# 5 個隕石指示點（橫排）
	for i in 5:
		var dot = ColorRect.new()
		dot.name = "Dot%d" % i
		dot.position = Vector2(8 + i * 56, 26)
		dot.size = Vector2(48, 24)
		dot.color = Color(0.2, 0.1, 0.0)  # 未觸發：暗色
		_meteor_panel.add_child(dot)

		var dot_lbl = Label.new()
		dot_lbl.name = "DotLabel%d" % i
		dot_lbl.text = "%d" % (i + 1)
		dot_lbl.position = Vector2(8 + i * 56, 26)
		dot_lbl.size = Vector2(48, 24)
		dot_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		dot_lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
		dot_lbl.add_theme_font_size_override("font_size", 12)
		dot_lbl.modulate = Color(0.5, 0.5, 0.5)
		_meteor_panel.add_child(dot_lbl)

		_meteor_dots.append(dot)

	# AOE 半徑標籤
	_radius_label = Label.new()
	_radius_label.name = "RadiusLabel"
	_radius_label.text = "AOE r=150"
	_radius_label.position = Vector2(8, 56)
	_radius_label.add_theme_font_size_override("font_size", 13)
	_radius_label.modulate = Color(1.0, 0.7, 0.4)
	_meteor_panel.add_child(_radius_label)

	# 命中計數
	_hit_label = Label.new()
	_hit_label.name = "HitLabel"
	_hit_label.text = "命中 0 個"
	_hit_label.position = Vector2(160, 56)
	_hit_label.add_theme_font_size_override("font_size", 13)
	_hit_label.modulate = Color(0.9, 0.6, 0.3)
	_meteor_panel.add_child(_hit_label)

func _update_meteor_dot(idx: int, hit: bool) -> void:
	if idx < 0 or idx >= _meteor_dots.size():
		return
	var dot = _meteor_dots[idx]
	if not is_instance_valid(dot):
		return
	if hit:
		dot.color = COLOR_METEOR
		var dot_lbl = _meteor_panel.get_node_or_null("DotLabel%d" % idx)
		if is_instance_valid(dot_lbl):
			dot_lbl.modulate = Color(1.0, 1.0, 1.0)
		# 閃光動畫
		var tween = dot.create_tween()
		tween.tween_property(dot, "modulate", Color(2.0, 1.5, 0.5), 0.1)
		tween.tween_property(dot, "modulate", Color(1.0, 1.0, 1.0), 0.2)
	else:
		dot.color = Color(0.3, 0.3, 0.3)  # 空揮：灰色
		var dot_lbl = _meteor_panel.get_node_or_null("DotLabel%d" % idx)
		if is_instance_valid(dot_lbl):
			dot_lbl.modulate = Color(0.5, 0.5, 0.5)

func _reset_meteor_dots() -> void:
	for i in _meteor_dots.size():
		var dot = _meteor_dots[i]
		if is_instance_valid(dot):
			dot.color = Color(0.2, 0.1, 0.0)
		var dot_lbl = _meteor_panel.get_node_or_null("DotLabel%d" % i)
		if is_instance_valid(dot_lbl):
			dot_lbl.modulate = Color(0.5, 0.5, 0.5)

func _show_start_banner(name: String) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()
	_banner = Control.new()
	_banner.position = Vector2(0, 100)
	_banner.size = Vector2(1280, 90)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = COLOR_BG
	_banner.add_child(bg)

	# 頂部火橙線
	var top_line = ColorRect.new()
	top_line.size = Vector2(1280, 3)
	top_line.color = COLOR_METEOR
	_banner.add_child(top_line)

	var lbl = Label.new()
	lbl.text = "☄️ %s 觸發連鎖隕石雨！5 顆隕石依序落下！" % name
	lbl.position = Vector2(0, 10)
	lbl.size = Vector2(1280, 40)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 28)
	lbl.modulate = COLOR_METEOR
	_banner.add_child(lbl)

	var sub = Label.new()
	sub.text = "每顆命中觸發連鎖擴大！5 顆全命中 → 完美隕石雨！全服 ×2.5！"
	sub.position = Vector2(0, 50)
	sub.size = Vector2(1280, 30)
	sub.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	sub.add_theme_font_size_override("font_size", 16)
	sub.modulate = Color(1.0, 0.8, 0.5)
	_banner.add_child(sub)

	# 底部火橙線
	var bot_line = ColorRect.new()
	bot_line.size = Vector2(1280, 3)
	bot_line.position = Vector2(0, 87)
	bot_line.color = COLOR_METEOR
	_banner.add_child(bot_line)

	# 滑入動畫
	_banner.position.y = -90
	var tween = _banner.create_tween()
	tween.tween_property(_banner, "position:y", 100.0, 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.tween_interval(3.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(_banner):
			_banner.queue_free()
			_banner = null
	)

func _show_perfect_popup(name: String, boost: float, secs: int) -> void:
	if is_instance_valid(_result_popup):
		_result_popup.queue_free()

	_do_flash(COLOR_PERFECT, 5)

	_result_popup = Control.new()
	_result_popup.position = Vector2(340, 200)
	_result_popup.size = Vector2(600, 180)
	add_child(_result_popup)

	var bg = ColorRect.new()
	bg.size = _result_popup.size
	bg.color = COLOR_BG
	_result_popup.add_child(bg)

	var border = ColorRect.new()
	border.size = _result_popup.size
	border.color = Color(1.0, 0.85, 0.0, 0.25)
	_result_popup.add_child(border)

	var title_lbl = Label.new()
	title_lbl.text = "☄️✨ 完美隕石雨！"
	title_lbl.position = Vector2(0, 15)
	title_lbl.size = Vector2(600, 50)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 32)
	title_lbl.modulate = COLOR_PERFECT
	_result_popup.add_child(title_lbl)

	var detail_lbl = Label.new()
	detail_lbl.text = "%s 5 顆全命中！" % name
	detail_lbl.position = Vector2(0, 70)
	detail_lbl.size = Vector2(600, 30)
	detail_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	detail_lbl.add_theme_font_size_override("font_size", 18)
	detail_lbl.modulate = Color(0.9, 0.9, 0.9)
	_result_popup.add_child(detail_lbl)

	var boost_lbl = Label.new()
	boost_lbl.text = "全服 ×%.0f 加成 %d 秒！" % [boost, secs]
	boost_lbl.position = Vector2(0, 110)
	boost_lbl.size = Vector2(600, 50)
	boost_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	boost_lbl.add_theme_font_size_override("font_size", 24)
	boost_lbl.modulate = COLOR_METEOR
	_result_popup.add_child(boost_lbl)

	# 彈出動畫
	_result_popup.scale = Vector2.ZERO
	_result_popup.pivot_offset = Vector2(300, 90)
	var tween = _result_popup.create_tween()
	tween.tween_property(_result_popup, "scale", Vector2(1.0, 1.0), 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
	tween.tween_interval(3.5)
	tween.tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_result_popup):
			_result_popup.queue_free()
			_result_popup = null
	)

	# 隱藏隕石面板
	if is_instance_valid(_meteor_panel):
		var tween2 = _meteor_panel.create_tween()
		tween2.tween_interval(3.5)
		tween2.tween_property(_meteor_panel, "modulate:a", 0.0, 0.5)
		tween2.tween_callback(func():
			if is_instance_valid(_meteor_panel):
				_meteor_panel.visible = false
				_meteor_panel.modulate.a = 1.0
		)

func _do_flash(color: Color, count: int) -> void:
	if not is_instance_valid(_flash_overlay):
		return
	var tween = _flash_overlay.create_tween()
	for i in count:
		tween.tween_property(_flash_overlay, "color", Color(color.r, color.g, color.b, 0.5), 0.08)
		tween.tween_property(_flash_overlay, "color", Color(color.r, color.g, color.b, 0.0), 0.12)

func _on_lucky_chain_meteor(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("player_name", "玩家")
	match event:
		"meteor_start":
			_show_start_banner(name)
			_reset_meteor_dots()
			if is_instance_valid(_meteor_panel):
				_meteor_panel.visible = true
			if is_instance_valid(_radius_label):
				_radius_label.text = "AOE r=150"
			if is_instance_valid(_hit_label):
				_hit_label.text = "命中 0 個"
			_do_flash(COLOR_METEOR, 3)
			ScreenShake.add_trauma(0.5)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"meteor_hit":
			var idx = data.get("meteor_index", 1) - 1  # 轉為 0-based
			var radius = data.get("aoe_radius", 150.0)
			var hits = data.get("hit_count", 0)
			_update_meteor_dot(idx, true)
			if is_instance_valid(_radius_label):
				_radius_label.text = "AOE r=%.0f" % radius
			if is_instance_valid(_hit_label):
				_hit_label.text = "命中 %d 個" % hits
			_do_flash(COLOR_METEOR, 1)
			ScreenShake.add_trauma(0.35)
		"meteor_miss":
			var idx = data.get("meteor_index", 1) - 1
			_update_meteor_dot(idx, false)
		"meteor_perfect":
			var boost = data.get("boost_mult", 2.5)
			var secs = data.get("boost_secs", 7)
			_show_perfect_popup(name, boost, secs)
			ScreenShake.add_trauma(0.8)
		"meteor_perfect_end":
			pass  # HUD 已處理橫幅
