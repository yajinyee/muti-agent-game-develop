## TripleLuckyFishPanel.gd — 三重幸運魚面板（DAY-190）
## 業界靈感：TaDa Gaming TriLuck™ 2026「trigger three different feature specifications simultaneously」
## 三重效果同時觸發：金幣雨 + 倍率加成 + 武器充能
## 視覺主題：金色三葉草 + 三重光環 + 彩虹閃光

extends Control

const TRIPLE_COLOR := Color(1.0, 0.85, 0.0)   # 金色（三重幸運）
const COIN_COLOR   := Color(1.0, 0.9, 0.2)    # 亮金（金幣雨）
const MULT_COLOR   := Color(0.4, 1.0, 0.4)    # 亮綠（倍率加成）
const WEAPON_COLOR := Color(0.4, 0.8, 1.0)    # 天藍（武器充能）

var _banner: Control = null
var _mult_bar: Control = null
var _mult_timer: float = 0.0
var _mult_duration: float = 12.0
var _is_my_trigger: bool = false

func _ready() -> void:
	set_process(false)
	if GameManager.has_signal("triple_lucky_fish"):
		GameManager.triple_lucky_fish.connect(_on_triple_lucky_fish)

func _process(delta: float) -> void:
	if _mult_bar == null or not is_instance_valid(_mult_bar):
		set_process(false)
		return
	_mult_timer -= delta
	if _mult_timer <= 0.0:
		_mult_timer = 0.0
		set_process(false)
		_fade_out_mult_bar()
		return
	# 更新進度條
	var pct := _mult_timer / _mult_duration
	var bar_fill = _mult_bar.get_node_or_null("Fill")
	if bar_fill:
		bar_fill.size.x = 280.0 * pct
		# 顏色漸變：綠→黃→橙
		if pct > 0.5:
			bar_fill.color = MULT_COLOR
		elif pct > 0.25:
			bar_fill.color = Color(1.0, 0.85, 0.0)
		else:
			bar_fill.color = Color(1.0, 0.5, 0.0)

func _on_triple_lucky_fish(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	match phase:
		"triple_start":
			_is_my_trigger = true
			_show_triple_start(data)
		"triple_broadcast":
			if not _is_my_trigger:
				_show_broadcast_banner(data)
		"mult_end":
			_is_my_trigger = false
			_fade_out_mult_bar()

func _show_triple_start(data: Dictionary) -> void:
	var coin_reward: int = data.get("coin_reward", 0)
	var coin_mult: float = data.get("coin_mult", 0.0)
	var mult_duration: float = data.get("mult_duration", 12.0)
	var weapon_charged: String = data.get("weapon_charged", "龍怒")

	_mult_duration = mult_duration

	# 全螢幕三重彩虹閃光（三次，間隔 150ms）
	_triple_rainbow_flash()

	# 頂部橫幅滑入
	_show_banner("🍀 三重幸運觸發！", TRIPLE_COLOR)

	# 三個效果同時顯示（錯開 200ms）
	var timer1 = get_tree().create_timer(0.3)
	timer1.timeout.connect(func():
		_show_effect_popup("💰 金幣雨", "+%d 金幣（×%.0f）" % [coin_reward, coin_mult], COIN_COLOR, Vector2(200, 280))
	)
	var timer2 = get_tree().create_timer(0.5)
	timer2.timeout.connect(func():
		_show_effect_popup("✨ 倍率加成", "+50%% 持續 %.0f 秒" % mult_duration, MULT_COLOR, Vector2(640, 280))
	)
	var timer3 = get_tree().create_timer(0.7)
	timer3.timeout.connect(func():
		_show_effect_popup("⚡ 武器充能", "%s 充能 1 發" % weapon_charged, WEAPON_COLOR, Vector2(1080, 280))
	)

	# 底部倍率加成進度條
	var timer4 = get_tree().create_timer(0.9)
	timer4.timeout.connect(func():
		_show_mult_bar(mult_duration)
	)

func _show_broadcast_banner(data: Dictionary) -> void:
	var player_name: String = data.get("player_name", "")
	var coin_reward: int = data.get("coin_reward", 0)
	_show_banner("🍀 %s 觸發三重幸運！+%d 金幣" % [player_name, coin_reward], TRIPLE_COLOR)

func _triple_rainbow_flash() -> void:
	# 三次彩虹閃光（金→綠→藍）
	var colors := [TRIPLE_COLOR, MULT_COLOR, WEAPON_COLOR]
	for i in range(3):
		var t = get_tree().create_timer(float(i) * 0.15)
		var c = colors[i]
		t.timeout.connect(func():
			_flash_screen(c, 0.5)
		)

func _flash_screen(color: Color, alpha: float) -> void:
	var overlay = ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, alpha)
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tween = create_tween()
	tween.tween_property(overlay, "modulate:a", 0.0, 0.25)
	tween.tween_callback(overlay.queue_free)

func _show_banner(text: String, color: Color) -> void:
	if _banner != null and is_instance_valid(_banner):
		_banner.queue_free()

	_banner = Control.new()
	_banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 60)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.color = Color(color.r * 0.3, color.g * 0.3, color.b * 0.3, 0.92)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_banner.add_child(bg)

	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", 22)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.add_child(label)

	# 滑入動畫
	_banner.position.y = -60
	var tween = create_tween()
	tween.tween_property(_banner, "position:y", 0.0, 0.3).set_trans(Tween.TRANS_BACK)
	# 4 秒後淡出
	var timer = get_tree().create_timer(4.0)
	timer.timeout.connect(func():
		if is_instance_valid(_banner):
			var t2 = create_tween()
			t2.tween_property(_banner, "modulate:a", 0.0, 0.4)
			t2.tween_callback(_banner.queue_free)
	)

func _show_effect_popup(title: String, desc: String, color: Color, pos: Vector2) -> void:
	var popup = Control.new()
	popup.position = pos - Vector2(80, 60)
	popup.custom_minimum_size = Vector2(160, 120)
	add_child(popup)

	var bg = ColorRect.new()
	bg.color = Color(color.r * 0.2, color.g * 0.2, color.b * 0.2, 0.9)
	bg.size = Vector2(160, 120)
	bg.position = Vector2.ZERO
	popup.add_child(bg)

	# 邊框
	var border = ColorRect.new()
	border.color = color
	border.size = Vector2(160, 4)
	border.position = Vector2(0, 0)
	popup.add_child(border)

	var title_label = Label.new()
	title_label.text = title
	title_label.add_theme_color_override("font_color", color)
	title_label.add_theme_font_size_override("font_size", 18)
	title_label.position = Vector2(8, 12)
	title_label.size = Vector2(144, 30)
	title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup.add_child(title_label)

	var desc_label = Label.new()
	desc_label.text = desc
	desc_label.add_theme_color_override("font_color", Color.WHITE)
	desc_label.add_theme_font_size_override("font_size", 14)
	desc_label.position = Vector2(8, 50)
	desc_label.size = Vector2(144, 60)
	desc_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	desc_label.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	popup.add_child(desc_label)

	# 彈跳動畫
	popup.scale = Vector2(0.5, 0.5)
	popup.pivot_offset = Vector2(80, 60)
	var tween = create_tween()
	tween.tween_property(popup, "scale", Vector2(1.0, 1.0), 0.3).set_trans(Tween.TRANS_BACK)

	# 3 秒後淡出
	var timer = get_tree().create_timer(3.0)
	timer.timeout.connect(func():
		if is_instance_valid(popup):
			var t2 = create_tween()
			t2.tween_property(popup, "modulate:a", 0.0, 0.4)
			t2.tween_callback(popup.queue_free)
	)

func _show_mult_bar(duration: float) -> void:
	if _mult_bar != null and is_instance_valid(_mult_bar):
		_mult_bar.queue_free()

	_mult_bar = Control.new()
	_mult_bar.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	_mult_bar.offset_top = -40
	_mult_bar.offset_bottom = 0
	add_child(_mult_bar)

	# 背景
	var bg = ColorRect.new()
	bg.color = Color(0.1, 0.1, 0.1, 0.8)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_mult_bar.add_child(bg)

	# 進度條填充
	var fill = ColorRect.new()
	fill.name = "Fill"
	fill.color = MULT_COLOR
	fill.size = Vector2(280.0, 20)
	fill.position = Vector2(get_viewport_rect().size.x / 2.0 - 140.0, 10)
	_mult_bar.add_child(fill)

	# 標籤
	var label = Label.new()
	label.text = "🍀 三重幸運 +50%% 倍率加成"
	label.add_theme_color_override("font_color", MULT_COLOR)
	label.add_theme_font_size_override("font_size", 14)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_mult_bar.add_child(label)

	_mult_timer = duration
	set_process(true)

func _fade_out_mult_bar() -> void:
	if _mult_bar != null and is_instance_valid(_mult_bar):
		var tween = create_tween()
		tween.tween_property(_mult_bar, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_mult_bar.queue_free)
		_mult_bar = null
