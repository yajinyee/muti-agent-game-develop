## LuckyJackpotFishPanel.gd — T126 幸運進階 Jackpot 魚 UI
## lucky-panel-agent 負責維護
## 業界依據：Jackpot Fishing Jili「Progressive Jackpot — Grand/Major/Minor/Mini 四層獎池」
## 視覺主題：金色 + 四層獎池顯示 + 王冠 + 閃光
extends CanvasLayer

const LAYER_Z = 26  # CanvasLayer layer 值

# 顏色主題
const COLOR_GRAND  = Color(1.0, 0.85, 0.0)   # 金色（Grand）
const COLOR_MAJOR  = Color(1.0, 0.55, 0.0)   # 橙色（Major）
const COLOR_MINOR  = Color(0.8, 0.8, 0.9)    # 銀色（Minor）
const COLOR_MINI   = Color(0.7, 0.4, 0.2)    # 銅色（Mini）
const COLOR_BG     = Color(0.05, 0.03, 0.0, 0.88)

var _banner: Control = null
var _pool_panel: Control = null
var _result_popup: Control = null
var _flash_overlay: ColorRect = null

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_pool_panel()
	GameManager.lucky_jackpot_fish.connect(_on_lucky_jackpot_fish)

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.85, 0.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

func _create_pool_panel() -> void:
	# 右上角四層獎池顯示
	_pool_panel = Control.new()
	_pool_panel.position = Vector2(960, 10)
	_pool_panel.size = Vector2(300, 110)
	_pool_panel.visible = false
	add_child(_pool_panel)

	var bg = ColorRect.new()
	bg.size = _pool_panel.size
	bg.color = COLOR_BG
	_pool_panel.add_child(bg)

	var title = Label.new()
	title.text = "🏆 JACKPOT POOL"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_GRAND
	_pool_panel.add_child(title)

	# 四層獎池標籤
	var tier_names = ["GRAND", "MAJOR", "MINOR", "MINI"]
	var tier_colors = [COLOR_GRAND, COLOR_MAJOR, COLOR_MINOR, COLOR_MINI]
	for i in 4:
		var lbl = Label.new()
		lbl.name = "Tier%d" % i
		lbl.text = "%s: ---" % tier_names[i]
		lbl.position = Vector2(8, 24 + i * 20)
		lbl.add_theme_font_size_override("font_size", 14)
		lbl.modulate = tier_colors[i]
		_pool_panel.add_child(lbl)

func _update_pool_display(mini: int, minor: int, major: int, grand: int) -> void:
	if not is_instance_valid(_pool_panel):
		return
	_pool_panel.visible = true
	var pools = [grand, major, minor, mini]
	for i in 4:
		var lbl = _pool_panel.get_node_or_null("Tier%d" % i)
		if is_instance_valid(lbl):
			lbl.text = ["GRAND", "MAJOR", "MINOR", "MINI"][i] + ": %d" % pools[i]

func _show_trigger_banner(name: String, mini: int, minor: int, major: int, grand: int) -> void:
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

	# 頂部金色線
	var top_line = ColorRect.new()
	top_line.size = Vector2(1280, 3)
	top_line.color = COLOR_GRAND
	_banner.add_child(top_line)

	var lbl = Label.new()
	lbl.text = "🏆 %s 觸發進階 Jackpot！" % name
	lbl.position = Vector2(0, 12)
	lbl.size = Vector2(1280, 40)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 28)
	lbl.modulate = COLOR_GRAND
	_banner.add_child(lbl)

	var sub = Label.new()
	sub.text = "Mini:%d  Minor:%d  Major:%d  Grand:%d" % [mini, minor, major, grand]
	sub.position = Vector2(0, 50)
	sub.size = Vector2(1280, 30)
	sub.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	sub.add_theme_font_size_override("font_size", 14)
	sub.modulate = Color(0.9, 0.9, 0.9)
	_banner.add_child(sub)

	# 底部金色線
	var bot_line = ColorRect.new()
	bot_line.size = Vector2(1280, 3)
	bot_line.position = Vector2(0, 87)
	bot_line.color = COLOR_GRAND
	_banner.add_child(bot_line)

	# 滑入動畫
	_banner.position.y = -90
	var tween = _banner.create_tween()
	tween.tween_property(_banner, "position:y", 100.0, 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.tween_interval(2.5)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(_banner):
			_banner.queue_free()
			_banner = null
	)

func _show_jackpot_result(name: String, tier_name: String, tier_idx: int, reward: int) -> void:
	if is_instance_valid(_result_popup):
		_result_popup.queue_free()

	# 依層級選顏色
	var tier_colors = [COLOR_MINI, COLOR_MINOR, COLOR_MAJOR, COLOR_GRAND]
	var tier_color = tier_colors[clamp(tier_idx, 0, 3)]
	var is_grand = tier_idx == 3

	# 全螢幕閃光（Grand 更強）
	var flash_count = 5 if is_grand else 3
	_do_flash(tier_color, flash_count)

	_result_popup = Control.new()
	_result_popup.position = Vector2(340, 200)
	_result_popup.size = Vector2(600, 200)
	add_child(_result_popup)

	var bg = ColorRect.new()
	bg.size = _result_popup.size
	bg.color = COLOR_BG
	_result_popup.add_child(bg)

	var border = ColorRect.new()
	border.size = _result_popup.size
	border.color = Color(tier_color.r, tier_color.g, tier_color.b, 0.3)
	_result_popup.add_child(border)

	var tier_lbl = Label.new()
	tier_lbl.text = "🏆 %s JACKPOT！" % tier_name.to_upper()
	tier_lbl.position = Vector2(0, 20)
	tier_lbl.size = Vector2(600, 50)
	tier_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	tier_lbl.add_theme_font_size_override("font_size", 32 if is_grand else 26)
	tier_lbl.modulate = tier_color
	_result_popup.add_child(tier_lbl)

	var name_lbl = Label.new()
	name_lbl.text = name
	name_lbl.position = Vector2(0, 70)
	name_lbl.size = Vector2(600, 30)
	name_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	name_lbl.add_theme_font_size_override("font_size", 16)
	name_lbl.modulate = Color(0.9, 0.9, 0.9)
	_result_popup.add_child(name_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "獲得 %d 金幣！" % reward
	reward_lbl.position = Vector2(0, 110)
	reward_lbl.size = Vector2(600, 50)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_lbl.add_theme_font_size_override("font_size", 28)
	reward_lbl.modulate = COLOR_GRAND
	_result_popup.add_child(reward_lbl)

	# 彈出動畫
	_result_popup.scale = Vector2.ZERO
	_result_popup.pivot_offset = Vector2(300, 100)
	var tween = _result_popup.create_tween()
	tween.tween_property(_result_popup, "scale", Vector2(1.0, 1.0), 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
	tween.tween_interval(3.0 if is_grand else 2.0)
	tween.tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_result_popup):
			_result_popup.queue_free()
			_result_popup = null
	)

func _show_grand_boost(name: String, mult: float, secs: int) -> void:
	# Grand Jackpot 全服加成橫幅
	var boost_banner = Control.new()
	boost_banner.position = Vector2(0, 200)
	boost_banner.size = Vector2(1280, 60)
	add_child(boost_banner)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.3, 0.2, 0.0, 0.9)
	boost_banner.add_child(bg)

	var lbl = Label.new()
	lbl.text = "🏆✨ GRAND JACKPOT！%s 全服 ×%.0f 加成 %d 秒！" % [name, mult, secs]
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 22)
	lbl.modulate = COLOR_GRAND
	boost_banner.add_child(lbl)

	# 脈動動畫
	var tween = boost_banner.create_tween().set_loops(5)
	tween.tween_property(boost_banner, "modulate:a", 0.5, 0.5)
	tween.tween_property(boost_banner, "modulate:a", 1.0, 0.5)
	boost_banner.create_tween().tween_interval(float(secs)).tween_callback(func():
		if is_instance_valid(boost_banner):
			boost_banner.queue_free()
	)

func _do_flash(color: Color, count: int) -> void:
	if not is_instance_valid(_flash_overlay):
		return
	var tween = _flash_overlay.create_tween()
	for i in count:
		tween.tween_property(_flash_overlay, "color", Color(color.r, color.g, color.b, 0.5), 0.08)
		tween.tween_property(_flash_overlay, "color", Color(color.r, color.g, color.b, 0.0), 0.12)

func _on_lucky_jackpot_fish(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_trigger_banner(name,
				data.get("mini_pool", 1000),
				data.get("minor_pool", 5000),
				data.get("major_pool", 20000),
				data.get("grand_pool", 50000)
			)
			_update_pool_display(
				data.get("mini_pool", 1000),
				data.get("minor_pool", 5000),
				data.get("major_pool", 20000),
				data.get("grand_pool", 50000)
			)
			_do_flash(COLOR_GRAND, 3)
			ScreenShake.add_trauma(0.4)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"jackpot_result":
			_show_jackpot_result(
				name,
				data.get("tier_name", "Mini"),
				data.get("tier_idx", 0),
				data.get("reward", 0)
			)
			_update_pool_display(
				data.get("mini_pool", 1000),
				data.get("minor_pool", 5000),
				data.get("major_pool", 20000),
				data.get("grand_pool", 50000)
			)
		"grand_boost":
			_show_grand_boost(name, data.get("boost_mult", 3.0), data.get("boost_secs", 10))
			ScreenShake.add_trauma(0.8)
		"grand_boost_end":
			if is_instance_valid(_pool_panel):
				var tween = _pool_panel.create_tween()
				tween.tween_property(_pool_panel, "modulate:a", 0.0, 0.5)
				tween.tween_callback(func():
					if is_instance_valid(_pool_panel):
						_pool_panel.visible = false
						_pool_panel.modulate.a = 1.0
				)
