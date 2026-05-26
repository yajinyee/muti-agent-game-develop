## LuckyTimeWarpPanel.gd — T128 幸運時間扭曲魚 UI
## lucky-panel-agent 負責維護
## 業界依據：業界原創「時間扭曲 — 全場目標移動速度降低 70%，持續 10 秒，傷害 ×2.0」
## 視覺主題：深紫色 + 時鐘 + 扭曲光環 + 計時條
extends CanvasLayer

const LAYER_Z = 28

const COLOR_WARP    = Color(0.55, 0.2, 0.86)  # 紫色
const COLOR_COLLAPSE = Color(0.5, 0.0, 0.8)   # 深紫（崩潰）
const COLOR_DAMAGE  = Color(1.0, 0.5, 0.0)    # 橙色（傷害加成）
const COLOR_BG      = Color(0.04, 0.0, 0.08, 0.88)

var _banner: Control = null
var _warp_indicator: Control = null
var _flash_overlay: ColorRect = null
var _warp_timer: float = 0.0
var _warp_active: bool = false

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_warp_indicator()
	GameManager.lucky_time_warp.connect(_on_lucky_time_warp)

func _process(delta: float) -> void:
	if _warp_active and _warp_timer > 0:
		_warp_timer -= delta
		_update_warp_timer()
		if _warp_timer <= 0:
			_warp_active = false

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(0.55, 0.2, 0.86, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

func _create_warp_indicator() -> void:
	_warp_indicator = Control.new()
	_warp_indicator.position = Vector2(960, 130)
	_warp_indicator.size = Vector2(300, 90)
	_warp_indicator.visible = false
	add_child(_warp_indicator)

	var bg = ColorRect.new()
	bg.size = _warp_indicator.size
	bg.color = COLOR_BG
	_warp_indicator.add_child(bg)

	var title = Label.new()
	title.name = "Title"
	title.text = "⏰ 時間扭曲"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 14)
	title.modulate = COLOR_WARP
	_warp_indicator.add_child(title)

	var dmg_lbl = Label.new()
	dmg_lbl.name = "DmgLabel"
	dmg_lbl.text = "傷害 ×2.0 | 速度 ×0.3"
	dmg_lbl.position = Vector2(8, 24)
	dmg_lbl.add_theme_font_size_override("font_size", 13)
	dmg_lbl.modulate = COLOR_DAMAGE
	_warp_indicator.add_child(dmg_lbl)

	var time_lbl = Label.new()
	time_lbl.name = "TimeLabel"
	time_lbl.text = "10.0s"
	time_lbl.position = Vector2(8, 44)
	time_lbl.add_theme_font_size_override("font_size", 22)
	time_lbl.modulate = Color.WHITE
	_warp_indicator.add_child(time_lbl)

	# 計時條
	var bar_bg = ColorRect.new()
	bar_bg.name = "BarBG"
	bar_bg.size = Vector2(284, 10)
	bar_bg.position = Vector2(8, 72)
	bar_bg.color = Color(0.1, 0.0, 0.15)
	_warp_indicator.add_child(bar_bg)

	var bar = ColorRect.new()
	bar.name = "Bar"
	bar.size = Vector2(284, 10)
	bar.position = Vector2(8, 72)
	bar.color = COLOR_WARP
	_warp_indicator.add_child(bar)

func _update_warp_timer() -> void:
	if not is_instance_valid(_warp_indicator):
		return
	var time_lbl = _warp_indicator.get_node_or_null("TimeLabel")
	if is_instance_valid(time_lbl):
		time_lbl.text = "%.1fs" % max(0, _warp_timer)
		if _warp_timer <= 3:
			time_lbl.modulate = Color(1.0, 0.3, 0.3)
		elif _warp_timer <= 6:
			time_lbl.modulate = Color(1.0, 0.8, 0.2)
		else:
			time_lbl.modulate = Color.WHITE
	var bar = _warp_indicator.get_node_or_null("Bar")
	var bar_bg = _warp_indicator.get_node_or_null("BarBG")
	if is_instance_valid(bar) and is_instance_valid(bar_bg):
		var pct = _warp_timer / 10.0
		bar.size.x = bar_bg.size.x * pct
		if pct <= 0.3:
			bar.color = Color(1.0, 0.3, 0.3)
		elif pct <= 0.6:
			bar.color = Color(0.8, 0.5, 1.0)
		else:
			bar.color = COLOR_WARP

func _show_warp_start(name: String, duration: float, speed_mult: float, damage_mult: float) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()
	_banner = Control.new()
	_banner.position = Vector2(0, 100)
	_banner.size = Vector2(1280, 80)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = COLOR_BG
	_banner.add_child(bg)

	var top_line = ColorRect.new()
	top_line.size = Vector2(1280, 3)
	top_line.color = COLOR_WARP
	_banner.add_child(top_line)

	var lbl = Label.new()
	lbl.text = "⏰ %s 觸發時間扭曲！全場慢速 %.0f 秒！傷害 ×%.0f！" % [name, duration, damage_mult]
	lbl.position = Vector2(0, 10)
	lbl.size = Vector2(1280, 40)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 26)
	lbl.modulate = COLOR_WARP
	_banner.add_child(lbl)

	var sub = Label.new()
	sub.text = "目標移動速度 ×%.1f！趁現在瘋狂射擊！" % speed_mult
	sub.position = Vector2(0, 50)
	sub.size = Vector2(1280, 25)
	sub.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	sub.add_theme_font_size_override("font_size", 13)
	sub.modulate = Color(0.8, 0.7, 1.0)
	_banner.add_child(sub)

	var bot_line = ColorRect.new()
	bot_line.size = Vector2(1280, 3)
	bot_line.position = Vector2(0, 77)
	bot_line.color = COLOR_WARP
	_banner.add_child(bot_line)

	_banner.position.y = -80
	var tween = _banner.create_tween()
	tween.tween_property(_banner, "position:y", 100.0, 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.tween_interval(3.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(_banner):
			_banner.queue_free()
			_banner = null
	)

	# 啟動計時器
	_warp_timer = duration
	_warp_active = true
	if is_instance_valid(_warp_indicator):
		_warp_indicator.visible = true
		_warp_indicator.modulate.a = 1.0

func _show_warp_end(kill_count: int) -> void:
	_warp_active = false
	if is_instance_valid(_warp_indicator):
		var tween = _warp_indicator.create_tween()
		tween.tween_property(_warp_indicator, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func():
			if is_instance_valid(_warp_indicator):
				_warp_indicator.visible = false
				_warp_indicator.modulate.a = 1.0
		)

	# 扭曲結束爆炸提示
	var end_lbl = Label.new()
	end_lbl.text = "⏰💥 時間扭曲結束！全場 HP -20%！擊破 %d 條！" % kill_count
	end_lbl.position = Vector2(0, 300)
	end_lbl.size = Vector2(1280, 40)
	end_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	end_lbl.add_theme_font_size_override("font_size", 22)
	end_lbl.modulate = COLOR_WARP
	add_child(end_lbl)
	var tween = end_lbl.create_tween()
	tween.tween_property(end_lbl, "position:y", 260.0, 0.5)
	tween.tween_interval(1.5)
	tween.tween_property(end_lbl, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(end_lbl):
			end_lbl.queue_free()
	)
	ScreenShake.add_trauma(0.5)

func _show_time_collapse(name: String, kill_count: int, boost_mult: float, boost_secs: int) -> void:
	_do_flash(COLOR_COLLAPSE, 5)
	ScreenShake.add_trauma(0.7)

	var popup = Control.new()
	popup.position = Vector2(340, 200)
	popup.size = Vector2(600, 180)
	add_child(popup)

	var bg = ColorRect.new()
	bg.size = popup.size
	bg.color = Color(0.08, 0.0, 0.12, 0.92)
	popup.add_child(bg)

	var lbl = Label.new()
	lbl.text = "⏰💥 時間崩潰！"
	lbl.position = Vector2(0, 15)
	lbl.size = Vector2(600, 50)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 32)
	lbl.modulate = COLOR_COLLAPSE
	popup.add_child(lbl)

	var detail_lbl = Label.new()
	detail_lbl.text = "%s 扭曲期間擊破 %d 條！" % [name, kill_count]
	detail_lbl.position = Vector2(0, 70)
	detail_lbl.size = Vector2(600, 35)
	detail_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	detail_lbl.add_theme_font_size_override("font_size", 18)
	detail_lbl.modulate = Color(0.8, 0.6, 1.0)
	popup.add_child(detail_lbl)

	var boost_lbl = Label.new()
	boost_lbl.text = "全服 ×%.0f 加成 %d 秒！" % [boost_mult, boost_secs]
	boost_lbl.position = Vector2(0, 115)
	boost_lbl.size = Vector2(600, 50)
	boost_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	boost_lbl.add_theme_font_size_override("font_size", 26)
	boost_lbl.modulate = Color(1.0, 0.85, 0.0)
	popup.add_child(boost_lbl)

	popup.scale = Vector2.ZERO
	popup.pivot_offset = Vector2(300, 90)
	var tween = popup.create_tween()
	tween.tween_property(popup, "scale", Vector2(1.0, 1.0), 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
	tween.tween_interval(3.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(popup):
			popup.queue_free()
	)

func _do_flash(color: Color, count: int) -> void:
	if not is_instance_valid(_flash_overlay):
		return
	var tween = _flash_overlay.create_tween()
	for i in count:
		tween.tween_property(_flash_overlay, "color", Color(color.r, color.g, color.b, 0.45), 0.08)
		tween.tween_property(_flash_overlay, "color", Color(color.r, color.g, color.b, 0.0), 0.12)

func _on_lucky_time_warp(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"warp_start":
			_show_warp_start(
				name,
				data.get("duration", 10.0),
				data.get("speed_mult", 0.3),
				data.get("damage_mult", 2.0)
			)
			_do_flash(COLOR_WARP, 3)
			ScreenShake.add_trauma(0.35)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"warp_end":
			_show_warp_end(data.get("kill_count", 0))
		"time_collapse":
			_show_time_collapse(
				name,
				data.get("kill_count", 0),
				data.get("boost_mult", 2.5),
				data.get("boost_secs", 6)
			)
		"collapse_end":
			pass  # 靜默處理
