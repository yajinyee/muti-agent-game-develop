## RockSkeletonConcertPanel.gd ???遝撉琿?瞍??選?DAY-192嚗?
## 璆剔?靘?嚗ILI 2026?ock Skeleton Concert ??Rock Skeleton and Super Awakening Performance, up to 3,000x??
## 閬死銝駁?嚗?蝝?遝 + ?喟泵?詨?? + 頞?閬粹?擃蔭 + 摰???脣漲璇?

extends Control

const ROCK_COLOR     := Color(1.0, 0.0, 1.0)    # 瘣??莎??遝??
const BEAT_COLOR     := Color(1.0, 0.4, 0.0)    # 璈嚗蝚衣敶?
const AWAKEN_COLOR   := Color(1.0, 0.8, 0.0)    # ?嚗?蝝死??
const ENCORE_COLOR   := Color(0.0, 1.0, 0.8)    # ???莎?摰嚗?

var _banner: Control = null
var _beat_counter: Label = null
var _encore_bar: Control = null
var _encore_timer: float = 0.0
var _is_encore_active: bool = false
var _total_kills: int = 0
var _is_awakening: bool = false

func _ready() -> void:
	set_process(false)
	if GameManager.has_signal("rock_skeleton_concert"):
		GameManager.rock_skeleton_concert.connect(_on_rock_skeleton_concert)

func _process(delta: float) -> void:
	if _is_encore_active:
		_encore_timer -= delta
		if _encore_timer <= 0.0:
			_encore_timer = 0.0
			_is_encore_active = false
			set_process(false)
			_hide_encore_bar()
		else:
			_update_encore_bar()

func _on_rock_skeleton_concert(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	match phase:
		"concert_start":
			_show_concert_start(data)
		"awakening":
			_show_awakening(data)
		"encore_start":
			_show_encore_start(data)
		"concert_end":
			_show_concert_end(data)
		"encore_end":
			_hide_encore_bar()
		"beat_result":
			_update_beat_counter(data)
		_:
			# note_N ??憿舐內?喟泵?詨??
			if phase.begins_with("note_"):
				_show_note_bomb(data)

# ?? concert_start ??????????????????????????????????????????????????????????????
func _show_concert_start(data: Dictionary) -> void:
	var killer_name: String = data.get("killer_name", "")
	var duration: int = data.get("duration_sec", 15)
	_total_kills = 0
	_is_awakening = false

	# 銝活瘣??脤????遝?嚗?
	_flash_screen(ROCK_COLOR, 0.6)
	var t1 = get_tree().create_timer(0.15)
	t1.timeout.connect(func(): _flash_screen(ROCK_COLOR, 0.45))
	var t2 = get_tree().create_timer(0.30)
	t2.timeout.connect(func(): _flash_screen(ROCK_COLOR, 0.3))

	# ?璈怠?
	_show_banner("??? %s ??皛暸疝擃??望?嚗?蝥?%d 蝘?" % [killer_name, duration], ROCK_COLOR)

	# ?閮??
	_show_beat_counter()

# ?? note_N ?喟泵?詨? ?????????????????????????????????????????????????????????????
func _show_note_bomb(data: Dictionary) -> void:
	var beat: int = data.get("beat", 0)
	var target_xs: Array = data.get("note_target_xs", [])
	var target_ys: Array = data.get("note_target_ys", [])
	var is_awakening: bool = data.get("is_awakening", false)

	var note_color := AWAKEN_COLOR if is_awakening else BEAT_COLOR

	# ?冽??璅?蝵桅＊蝷粹蝚衣泵??
	for i in range(min(target_xs.size(), target_ys.size())):
		var tx: float = target_xs[i]
		var ty: float = target_ys[i]
		_spawn_note_at(tx, ty, note_color, beat)

# ?? beat_result ?祆?蝯? ????????????????????????????????????????????????????????
func _update_beat_counter(data: Dictionary) -> void:
	_total_kills = data.get("total_kills", _total_kills)
	var beat_kills: int = data.get("beat_kills", 0)
	var beat_reward: int = data.get("beat_reward", 0)
	var is_awakening: bool = data.get("is_awakening", false)

	if _beat_counter and is_instance_valid(_beat_counter):
		var icon := "??" if is_awakening else "?"
		_beat_counter.text = "%s ? %d ?? % [icon, _total_kills]"
		var color := AWAKEN_COLOR if is_awakening else ROCK_COLOR
		_beat_counter.add_theme_color_override("font_color", color)

	# ?祆????湔?憿舐內瘚桀????
	if beat_kills > 0 and beat_reward > 0:
		_show_floating_reward("+%d ?馳 ?%d" % [beat_reward, beat_kills], BEAT_COLOR)

# ?? awakening 頞?閬粹?擃蔭 ??????????????????????????????????????????????????????
func _show_awakening(data: Dictionary) -> void:
	var awakened_count: int = data.get("awakened_count", 0)
	_is_awakening = true

	# ?撘琿???銝活嚗?
	_flash_screen(AWAKEN_COLOR, 0.7)
	var t1 = get_tree().create_timer(0.12)
	t1.timeout.connect(func(): _flash_screen(AWAKEN_COLOR, 0.55))
	var t2 = get_tree().create_timer(0.24)
	t2.timeout.connect(func(): _flash_screen(AWAKEN_COLOR, 0.4))

	# ?湔璈怠?
	_show_banner("???? 頞?閬粹?嚗?d ?璅?HP ?? 70%%嚗? % awakened_count, AWAKEN_COLOR)"

	# 銝剖亢憭批??內
	_show_center_popup("??頞?閬粹?嚗n?典?格? HP ?? 70%%嚗?, AWAKEN_COLOR)"

	# ?湔閮?券???
	if _beat_counter and is_instance_valid(_beat_counter):
		_beat_counter.add_theme_color_override("font_color", AWAKEN_COLOR)

# ?? encore_start 摰?? ???????????????????????????????????????????????????????
func _show_encore_start(data: Dictionary) -> void:
	var total_kills: int = data.get("total_kills", 0)
	var total_reward: int = data.get("total_reward", 0)
	var encore_duration: int = data.get("encore_duration", 10)
	var encore_bonus: float = data.get("encore_bonus", 0.30)

	_encore_timer = float(encore_duration)
	_is_encore_active = true
	set_process(true)

	# ???脣撥??嚗??荔?嚗?
	_flash_screen(ENCORE_COLOR, 0.65)
	var t1 = get_tree().create_timer(0.15)
	t1.timeout.connect(func(): _flash_screen(ENCORE_COLOR, 0.5))

	# ?湔璈怠?
	_show_banner("? 摰嚗???%d ???冽? +%d%% ?? %d 蝘?" % [
		total_kills, int(encore_bonus * 100), encore_duration
	], ENCORE_COLOR)

	# ?喳皛蝯?敶?
	_show_result_popup(total_kills, total_reward, true, encore_bonus, encore_duration)

	# 摨摰???脣漲璇?
	_show_encore_bar(encore_duration)

# ?? concert_end 瞍?????∪??荔???????????????????????????????????????????????
func _show_concert_end(data: Dictionary) -> void:
	var total_kills: int = data.get("total_kills", 0)
	var total_reward: int = data.get("total_reward", 0)

	# 瘛∪璈怠?
	_hide_banner()

	# ?喳皛蝯?敶?
	_show_result_popup(total_kills, total_reward, false, 0.0, 0)

	# 瘛∪閮??
	if _beat_counter and is_instance_valid(_beat_counter):
		var tween = create_tween()
		tween.tween_property(_beat_counter, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_beat_counter.queue_free)
		_beat_counter = null

# ?? 頛?賣 ????????????????????????????????????????????????????????????????????

func _flash_screen(color: Color, alpha: float) -> void:
	var overlay = ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, alpha)
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tween = create_tween()
	tween.tween_property(overlay, "modulate:a", 0.0, 0.35)
	tween.tween_callback(overlay.queue_free)

func _show_banner(text: String, color: Color) -> void:
	if _banner != null and is_instance_valid(_banner):
		_banner.queue_free()

	_banner = Control.new()
	_banner.set_anchors_and_offsets_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 52)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.color = Color(color.r * 0.2, color.g * 0.1, color.b * 0.2, 0.92)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_banner.add_child(bg)

	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", 19)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.add_child(label)

	# 皛?
	_banner.position.y = -52
	var tween = create_tween()
	tween.tween_property(_banner, "position:y", 0.0, 0.25).set_trans(Tween.TRANS_BACK)

func _hide_banner() -> void:
	if _banner != null and is_instance_valid(_banner):
		var tween = create_tween()
		tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
		tween.tween_callback(_banner.queue_free)
		_banner = null

func _show_beat_counter() -> void:
	if _beat_counter != null and is_instance_valid(_beat_counter):
		_beat_counter.queue_free()

	_beat_counter = Label.new()
	_beat_counter.text = "? ? 0 ??"
	_beat_counter.add_theme_color_override("font_color", ROCK_COLOR)
	_beat_counter.add_theme_font_size_override("font_size", 18)
	_beat_counter.set_anchors_preset(Control.PRESET_BOTTOM_RIGHT)
	_beat_counter.offset_left = -220
	_beat_counter.offset_top = -44
	_beat_counter.offset_right = -10
	_beat_counter.offset_bottom = -10
	add_child(_beat_counter)

func _spawn_note_at(x: float, y: float, color: Color, beat: int) -> void:
	# ?喟泵蝚西?嚗 ???恬?靘??訾漱?選?
	var note_label = Label.new()
	note_label.text = "?? if beat % 2 == 0 else "??
	note_label.add_theme_color_override("font_color", color)
	note_label.add_theme_font_size_override("font_size", 28)
	note_label.position = Vector2(x - 14, y - 14)
	add_child(note_label)

	# ??憌? + 瘛∪
	var tween = create_tween()
	tween.set_parallel(true)
	tween.tween_property(note_label, "position:y", y - 60, 0.5).set_trans(Tween.TRANS_QUAD)
	tween.tween_property(note_label, "modulate:a", 0.0, 0.5)
	tween.chain().tween_callback(note_label.queue_free)

	# ???
	var circle = ColorRect.new()
	circle.color = Color(color.r, color.g, color.b, 0.35)
	circle.size = Vector2(20, 20)
	circle.position = Vector2(x - 10, y - 10)
	add_child(circle)
	var tween2 = create_tween()
	tween2.set_parallel(true)
	tween2.tween_property(circle, "size", Vector2(80, 80), 0.3)
	tween2.tween_property(circle, "position", Vector2(x - 40, y - 40), 0.3)
	tween2.tween_property(circle, "modulate:a", 0.0, 0.3)
	tween2.chain().tween_callback(circle.queue_free)

func _show_center_popup(text: String, color: Color) -> void:
	var popup = Label.new()
	popup.text = text
	popup.add_theme_color_override("font_color", color)
	popup.add_theme_font_size_override("font_size", 30)
	popup.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	popup.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	add_child(popup)

	popup.scale = Vector2(0.5, 0.5)
	popup.pivot_offset = popup.size / 2.0
	var tween = create_tween()
	tween.tween_property(popup, "scale", Vector2(1.0, 1.0), 0.3).set_trans(Tween.TRANS_BACK)

	var timer = get_tree().create_timer(2.5)
	timer.timeout.connect(func():
		if is_instance_valid(popup):
			var t2 = create_tween()
			t2.tween_property(popup, "modulate:a", 0.0, 0.4)
			t2.tween_callback(popup.queue_free)
	)

func _show_floating_reward(text: String, color: Color) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", 16)
	# ?冽?雿蔭嚗葉憭桀??喉?
	label.position = Vector2(600 + randf_range(-80, 80), 300 + randf_range(-40, 40))
	add_child(label)

	var tween = create_tween()
	tween.set_parallel(true)
	tween.tween_property(label, "position:y", label.position.y - 50, 0.8)
	tween.tween_property(label, "modulate:a", 0.0, 0.8)
	tween.chain().tween_callback(label.queue_free)

func _show_result_popup(total_kills: int, total_reward: int, has_encore: bool, encore_bonus: float, encore_duration: int) -> void:
	var popup = Control.new()
	popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	popup.custom_minimum_size = Vector2(220, 140)
	popup.offset_left = -230
	popup.offset_top = -70
	popup.offset_right = -10
	popup.offset_bottom = 70
	add_child(popup)

	var bg = ColorRect.new()
	bg.color = Color(0.1, 0.0, 0.15, 0.92)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	popup.add_child(bg)

	var vbox = VBoxContainer.new()
	vbox.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	vbox.add_theme_constant_override("separation", 6)
	popup.add_child(vbox)

	var title_color := ENCORE_COLOR if has_encore else ROCK_COLOR
	var title_text := "? 摰嚗??望?蝯?嚗? if has_encore else "? 瞍????""
	_add_result_label(vbox, title_text, title_color, 16)
	_add_result_label(vbox, "??格?嚗?d ?? % total_kills, Color.WHITE, 14)"
	_add_result_label(vbox, "蝮賜??蛛?%d ?馳" % total_reward, BEAT_COLOR, 14)
	if has_encore:
		_add_result_label(vbox, "摰??嚗?%d%% ? %d 蝘? % [int(encore_bonus * 100), encore_duration], ENCORE_COLOR, 14)"

	# 敺?湔???
	popup.position.x += 240
	var tween = create_tween()
	tween.tween_property(popup, "position:x", popup.position.x - 240, 0.3).set_trans(Tween.TRANS_BACK)

	# 5 蝘?瘛∪
	var timer = get_tree().create_timer(5.0)
	timer.timeout.connect(func():
		if is_instance_valid(popup):
			var t2 = create_tween()
			t2.tween_property(popup, "modulate:a", 0.0, 0.5)
			t2.tween_callback(popup.queue_free)
	)

func _add_result_label(parent: Control, text: String, color: Color, size: int) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", size)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	parent.add_child(label)

func _show_encore_bar(duration: int) -> void:
	if _encore_bar != null and is_instance_valid(_encore_bar):
		_encore_bar.queue_free()

	_encore_bar = Control.new()
	_encore_bar.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	_encore_bar.custom_minimum_size = Vector2(0, 28)
	_encore_bar.offset_top = -28
	_encore_bar.offset_bottom = 0
	add_child(_encore_bar)

	var bg = ColorRect.new()
	bg.color = Color(0.0, 0.15, 0.15, 0.85)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_encore_bar.add_child(bg)

	var fill = ColorRect.new()
	fill.name = "Fill"
	fill.color = ENCORE_COLOR
	fill.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_encore_bar.add_child(fill)

	var label = Label.new()
	label.name = "Label"
	label.text = "? 摰?? +30%% ?拚? %d 蝘? % duration"
	label.add_theme_color_override("font_color", Color.WHITE)
	label.add_theme_font_size_override("font_size", 14)
	label.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_encore_bar.add_child(label)

func _update_encore_bar() -> void:
	if _encore_bar == null or not is_instance_valid(_encore_bar):
		return
	var fill = _encore_bar.get_node_or_null("Fill")
	var label = _encore_bar.get_node_or_null("Label")
	if fill:
		fill.anchor_right = _encore_timer / float(RockSkeletonEncoreDurationSec if false else 10)
	if label:
		label.text = "? 摰?? +30%% ?拚? %.1f 蝘? % _encore_timer"

func _hide_encore_bar() -> void:
	if _encore_bar != null and is_instance_valid(_encore_bar):
		var tween = create_tween()
		tween.tween_property(_encore_bar, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_encore_bar.queue_free)
		_encore_bar = null
	# 瘛∪璈怠????詨
	_hide_banner()
	if _beat_counter and is_instance_valid(_beat_counter):
		var tween2 = create_tween()
		tween2.tween_property(_beat_counter, "modulate:a", 0.0, 0.5)
		tween2.tween_callback(_beat_counter.queue_free)
		_beat_counter = null
