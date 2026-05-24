## TitlePanel.gd ??蝔梯?憿舐內?Ｘ嚗AY-068嚗?
## 憿舐內?拙振?嗅?蝔梯?嚗圾?蝔梯????粹
extends Node2D

# ?? 撣豢 ??????????????????????????????????????????????????????????????????????
const PANEL_WIDTH  := 200
const PANEL_HEIGHT := 28
const NOTIFY_DURATION := 3.5  # 蝔梯?閫???憿舐內蝘

# ?? 蝭暺?????????????????????????????????????????????????????????????????????
var _font: FontFile
var _bg: ColorRect
var _title_label: Label
var _notify_container: Control  # 蝔梯?閫???摰孵

# ?? ?????????????????????????????????????????????????????????????????????????
var _current_title_id: String = "novice"
var _current_title_name: String = "?唳?閮???"
var _current_title_icon: String = "?"
var _current_title_color: String = "#AAAAAA"
var _notify_queue: Array = []
var _is_showing_notify: bool = false

# ?? ????????????????????????????????????????????????????????????????????????
func setup(font: FontFile) -> void:
	_font = font
	_build_ui()
	_connect_signals()

func _build_ui() -> void:
	# ?
	_bg = ColorRect.new()
	_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_bg.color = Color(0.05, 0.05, 0.15, 0.75)
	add_child(_bg)

	# 蝔梯?璅惜
	_title_label = Label.new()
	_title_label.position = Vector2(6, 4)
	_title_label.size = Vector2(PANEL_WIDTH - 12, PANEL_HEIGHT - 8)
	_title_label.text = _current_title_icon + " " + _current_title_name
	_title_label.add_theme_color_override("font_color", Color.html(_current_title_color))
	if _font:
		_title_label.add_theme_font_override("font", _font)
		_title_label.add_theme_font_size_override("font_size", 13)
	add_child(_title_label)

	# ?摰孵嚗anvasLayer 銝?
	_notify_container = Control.new()
	_notify_container.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(_notify_container)

func _connect_signals() -> void:
	if GameManager.has_signal("title_unlocked"):
		GameManager.title_unlocked.connect(_on_title_unlocked)
	if GameManager.has_signal("player_updated"):
		GameManager.player_updated.connect(_on_player_updated)

# ?? ?湔蝔梯?憿舐內 ??????????????????????????????????????????????????????????????
func update_title(title_id: String, title_name: String, title_icon: String, title_color: String) -> void:
	_current_title_id = title_id
	_current_title_name = title_name
	_current_title_icon = title_icon
	_current_title_color = title_color
	if _title_label:
		_title_label.text = title_icon + " " + title_name
		_title_label.add_theme_color_override("font_color", Color.html(title_color))

# ?? 蝔梯?閫??? ??????????????????????????????????????????????????????????????
func _on_title_unlocked(data: Dictionary) -> void:
	_notify_queue.append(data)
	if not _is_showing_notify:
		_show_next_notify()

func _show_next_notify() -> void:
	if _notify_queue.is_empty():
		_is_showing_notify = false
		return
	_is_showing_notify = true
	var data: Dictionary = _notify_queue.pop_front()
	_spawn_title_notify(data)

func _spawn_title_notify(data: Dictionary) -> void:
	var title_name: String = data.get("title_name", "")"
	var title_icon: String = data.get("title_icon", "??")
	var title_color: String = data.get("title_color", "#FFD700")
	var description: String = data.get("description", "")

	# ?摰孵嚗摰?恍?喃?閫?
	var notify := Control.new()
	notify.set_anchors_and_offsets_preset(Control.PRESET_BOTTOM_RIGHT)
	notify.offset_left = -280
	notify.offset_top = -90
	notify.offset_right = -10
	notify.offset_bottom = -10
	_notify_container.add_child(notify)

	# ?
	var bg := ColorRect.new()
	bg.size = Vector2(270, 80)
	bg.color = Color(0.05, 0.05, 0.15, 0.92)
	notify.add_child(bg)

	# 撌血敶抵??
	var bar := ColorRect.new()
	bar.size = Vector2(4, 80)
	bar.color = Color.html(title_color)
	notify.add_child(bar)

	# 璅?銵????蝔梯?閫??嚗?
	var header := Label.new()
	header.position = Vector2(12, 6)
	header.text = "?? 蝔梯?閫??嚗?"
	header.add_theme_color_override("font_color", Color.html(title_color))
	if _font:
		header.add_theme_font_override("font", _font)
		header.add_theme_font_size_override("font_size", 12)
	notify.add_child(header)

	# 蝔梯??迂
	var name_label := Label.new()
	name_label.position = Vector2(12, 26)
	name_label.text = title_icon + " " + title_name
	name_label.add_theme_color_override("font_color", Color.WHITE)
	if _font:
		name_label.add_theme_font_override("font", _font)
		name_label.add_theme_font_size_override("font_size", 15)
	notify.add_child(name_label)

	# ?膩
	if description != "":
		var desc_label := Label.new()
		desc_label.position = Vector2(12, 52)
		desc_label.text = description
		desc_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
		if _font:
			desc_label.add_theme_font_override("font", _font)
			desc_label.add_theme_font_size_override("font_size", 11)
		notify.add_child(desc_label)

	# 皛?
	notify.modulate.a = 0.0
	notify.position.x += 30
	var tween := notify.create_tween()
	tween.tween_property(notify, "modulate:a", 1.0, 0.25)
	tween.parallel().tween_property(notify, "position:x", notify.position.x - 30, 0.25)
	tween.tween_interval(NOTIFY_DURATION)
	tween.tween_property(notify, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		notify.queue_free()
		_show_next_notify()
	)

# ?? 閮??? ??????????????????????????????????????????????????????????????????
func _on_player_updated(snapshot: Dictionary) -> void:
	var title_id: String = snapshot.get("title_id", "novice")
	var title_name: String = snapshot.get("title_name", "?唳?閮???)"
	var title_icon: String = snapshot.get("title_icon", "?")
	var title_color: String = snapshot.get("title_color", "#AAAAAA")
	update_title(title_id, title_name, title_icon, title_color)
