## DMPanel.gd ???拙振蝘??Ｘ嚗AY-103嚗?
## 憟賢??隞乩??貊??閮??Ｙ?閮銝?敺?＊蝷?
## ?游???FriendPanel嚗??末??蝔梢???閰望?嚗?
extends Node2D

# ---- 撣豢 ----
const PANEL_WIDTH  := 320
const PANEL_HEIGHT := 280
const MAX_DISPLAY_MESSAGES := 20

# ---- 蝭暺???----
var _pixel_font: Font = null
var _panel_bg: ColorRect = null
var _messages_container: Node2D = null
var _input_field: LineEdit = null
var _send_btn: Button = null
var _title_label: Label = null
var _unread_badge: Label = null
var _toggle_btn: Button = null

# ---- ???----
var _is_open: bool = false
var _current_friend_id: String = ""
var _current_friend_name: String = ""
var _messages: Array = []  # {from_id, from_name, content, sent_at, is_offline}
var _unread_count: int = 0

# ---- 閮? ----
signal dm_opened(friend_id: String)

# ---- ????----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_build_toggle_btn()
	_build_panel()
	_connect_signals()
	visible = false

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

## 撱箇?????嚗＊蝷箏 TopBar嚗?
func _build_toggle_btn() -> void:
	_toggle_btn = Button.new()
	_toggle_btn.text = "?"
	_toggle_btn.size = Vector2(32, 24)
	_toggle_btn.position = Vector2(0, 0)
	_toggle_btn.flat = true
	_toggle_btn.tooltip_text = "蝘?"
	if _pixel_font:
		_toggle_btn.add_theme_font_override("font", _pixel_font)
		_toggle_btn.add_theme_font_size_override("font_size", 14)
	add_child(_toggle_btn)

	# ?芾?敺賜?
	_unread_badge = Label.new()
	_unread_badge.position = Vector2(20, -4)
	_unread_badge.text = ""
	_unread_badge.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))
	if _pixel_font:
		_unread_badge.add_theme_font_override("font", _pixel_font)
		_unread_badge.add_theme_font_size_override("font_size", 9)
	_unread_badge.visible = false
	add_child(_unread_badge)

## 撱箇?銝駁??
func _build_panel() -> void:
	_panel_bg = ColorRect.new()
	_panel_bg.position = Vector2(-PANEL_WIDTH + 32, 28)
	_panel_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_panel_bg.color = Color(0.04, 0.02, 0.14, 0.95)
	_panel_bg.visible = false
	add_child(_panel_bg)

	# 璅?
	_title_label = Label.new()
	_title_label.position = Vector2(8, 4)
	_title_label.text = "? 蝘?"
	_title_label.add_theme_color_override("font_color", Color(0.6, 0.9, 1.0))
	if _pixel_font:
		_title_label.add_theme_font_override("font", _pixel_font)
		_title_label.add_theme_font_size_override("font_size", 12)
	_panel_bg.add_child(_title_label)

	# ????
	var close_btn := Button.new()
	close_btn.text = "??"
	close_btn.position = Vector2(PANEL_WIDTH - 24, 2)
	close_btn.size = Vector2(20, 20)
	close_btn.flat = true
	close_btn.add_theme_color_override("font_color", Color(0.8, 0.4, 0.4))
	if _pixel_font:
		close_btn.add_theme_font_override("font", _pixel_font)
		close_btn.add_theme_font_size_override("font_size", 10)
	close_btn.pressed.connect(func(): _close_panel())
	_panel_bg.add_child(close_btn)

	# ??蝺?
	var sep := ColorRect.new()
	sep.position = Vector2(4, 22)
	sep.size = Vector2(PANEL_WIDTH - 8, 1)
	sep.color = Color(0.3, 0.3, 0.5, 0.6)
	_panel_bg.add_child(sep)

	# 閮摰孵嚗????
	_messages_container = Node2D.new()
	_messages_container.position = Vector2(0, 26)
	_panel_bg.add_child(_messages_container)

	# 頛詨獢?
	_input_field = LineEdit.new()
	_input_field.position = Vector2(4, PANEL_HEIGHT - 30)
	_input_field.size = Vector2(PANEL_WIDTH - 60, 24)
	_input_field.placeholder_text = "頛詨閮..."
	_input_field.max_length = 200
	if _pixel_font:
		_input_field.add_theme_font_override("font", _pixel_font)
		_input_field.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_input_field)

	# ?潮???
	_send_btn = Button.new()
	_send_btn.text = "?潮?"
	_send_btn.position = Vector2(PANEL_WIDTH - 54, PANEL_HEIGHT - 30)
	_send_btn.size = Vector2(50, 24)
	if _pixel_font:
		_send_btn.add_theme_font_override("font", _pixel_font)
		_send_btn.add_theme_font_size_override("font_size", 9)
	_send_btn.pressed.connect(_on_send_pressed)
	_panel_bg.add_child(_send_btn)

	# 頛詨獢?Enter ?潮?
	_input_field.text_submitted.connect(func(_t): _on_send_pressed())

## ??閮?
func _connect_signals() -> void:
	_toggle_btn.pressed.connect(_on_toggle_pressed)

	if GameManager.has_signal("dm_received"):
		GameManager.dm_received.connect(_on_dm_received)
	if GameManager.has_signal("dm_sent"):
		GameManager.dm_sent.connect(_on_dm_sent)
	if GameManager.has_signal("dm_error"):
		GameManager.dm_error.connect(_on_dm_error)

## ???摰末??撠店
func open_conversation(friend_id: String, friend_name: String) -> void:
	_current_friend_id = friend_id
	_current_friend_name = friend_name
	_title_label.text = "? ??%s ??閰? % friend_name"
	_is_open = true
	visible = true
	_panel_bg.visible = true
	_refresh_messages()
	emit_signal("dm_opened", friend_id)

func _close_panel() -> void:
	_is_open = false
	_panel_bg.visible = false
	visible = false

func _on_toggle_pressed() -> void:
	if _is_open:
		_close_panel()
	else:
		# 憿舐內憟賢??豢??內
		_show_friend_select_hint()

func _show_friend_select_hint() -> void:
	var hint := Label.new()
	hint.text = "隢?憟賢??”暺?憟賢??迂??撠店"
	hint.position = Vector2(-200, -30)
	hint.add_theme_color_override("font_color", Color(0.8, 0.8, 0.6))
	if _pixel_font:
		hint.add_theme_font_override("font", _pixel_font)
		hint.add_theme_font_size_override("font_size", 10)
	add_child(hint)
	var tween = create_tween()
	tween.tween_interval(2.0)
	tween.tween_property(hint, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(hint): hint.queue_free()
	)

func _on_send_pressed() -> void:
	var content = _input_field.text.strip_edges()
	if content.is_empty() or _current_friend_id.is_empty():
		return

	NetworkManager.send_message({
		"type": "send_dm",
		"payload": {
			"to_id": _current_friend_id,
			"content": content
		}
	})
	_input_field.text = ""

## ?嗅蝘?
func _on_dm_received(data: Dictionary) -> void:
	var from_id = data.get("from_id", "")
	var from_name = data.get("from_name", "憟賢?")
	var content = data.get("content", "")
	var is_offline = data.get("is_offline", false)

	# ?閮閮?
	_messages.append({
		"from_id": from_id,
		"from_name": from_name,
		"content": content,
		"sent_at": data.get("sent_at", 0),
		"is_mine": false,
		"is_offline": is_offline
	})

	# 憒??Ｘ??銝?嗅?撠店嚗?圈＊蝷?
	if _is_open and from_id == _current_friend_id:
		_refresh_messages()
	else:
		# 憓??芾?閮
		_unread_count += 1
		_update_unread_badge()
		# 憿舐內瘚桀??
		var prefix = "? [?Ｙ?]" if is_offline else "?"
		_show_notification("%s %s嚗?s" % [prefix, from_name, content.substr(0, 30)],
			Color(0.4, 0.8, 1.0))

## ?潮??Ⅱ隤?
func _on_dm_sent(data: Dictionary) -> void:
	var remaining = data.get("remaining", 0)
	# ??芸楛???臬閮?嚗? input_field ??嚗?撌脫?蝛綽??隞亙? payload ??
	# ?ㄐ?芣?啗??賊＊蝷?
	if remaining <= 10:
		_show_notification("? 隞?拚? %d ???? % remaining, Color(0.8, 0.8, 0.6))"

## ?潮仃??
func _on_dm_error(data: Dictionary) -> void:
	var msg = data.get("message", "?潮仃??)"
	_show_notification("??%s" % msg, Color(1.0, 0.4, 0.4))

## ?瑟閮憿舐內
func _refresh_messages() -> void:
	for child in _messages_container.get_children():
		child.queue_free()

	# ?芷＊蝷箄??嗅?憟賢???閰?
	var filtered = []
	for m in _messages:
		if m.get("from_id") == _current_friend_id or m.get("is_mine", false):
			filtered.append(m)

	# ?憭＊蝷箸???MAX_DISPLAY_MESSAGES ??
	var start = max(0, filtered.size() - MAX_DISPLAY_MESSAGES)
	var y_offset = 0
	for i in range(start, filtered.size()):
		var m = filtered[i]
		y_offset = _build_message_bubble(y_offset, m)

	# 皜?芾?
	_unread_count = 0
	_update_unread_badge()

## 撱箇?閮瘞?部
func _build_message_bubble(y: int, msg_data: Dictionary) -> int:
	var is_mine = msg_data.get("is_mine", false)
	var from_name = msg_data.get("from_name", "?")
	var content = msg_data.get("content", "")
	var is_offline = msg_data.get("is_offline", false)

	# ?迂璅惜
	var name_label := Label.new()
	name_label.position = Vector2(8 if not is_mine else PANEL_WIDTH - 120, y)
	name_label.text = ("雿? if is_mine else from_name) + ("嚗蝺?" if is_offline else "")"
	name_label.add_theme_color_override("font_color",
		Color(0.6, 0.9, 1.0) if is_mine else Color(1.0, 0.85, 0.5))
	if _pixel_font:
		name_label.add_theme_font_override("font", _pixel_font)
		name_label.add_theme_font_size_override("font_size", 8)
	_messages_container.add_child(name_label)
	y += 12

	# 閮瘞?部?
	var bubble_bg := ColorRect.new()
	var bubble_w = min(int(len(content) * 6 + 16), PANEL_WIDTH - 24)
	bubble_bg.position = Vector2(8 if not is_mine else PANEL_WIDTH - bubble_w - 8, y)
	bubble_bg.size = Vector2(bubble_w, 20)
	bubble_bg.color = Color(0.15, 0.1, 0.35) if is_mine else Color(0.1, 0.15, 0.3)
	_messages_container.add_child(bubble_bg)

	# 閮??
	var content_label := Label.new()
	content_label.position = Vector2(4, 2)
	content_label.text = content
	content_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	if _pixel_font:
		content_label.add_theme_font_override("font", _pixel_font)
		content_label.add_theme_font_size_override("font_size", 9)
	bubble_bg.add_child(content_label)

	return y + 26

## ?湔?芾?敺賜?
func _update_unread_badge() -> void:
	if not is_instance_valid(_unread_badge):
		return
	if _unread_count > 0:
		_unread_badge.text = str(_unread_count)
		_unread_badge.visible = true
	else:
		_unread_badge.visible = false

## 憿舐內瘚桀??
func _show_notification(text: String, color: Color) -> void:
	var notify := Label.new()
	notify.text = text
	notify.position = Vector2(-200, -30)
	notify.add_theme_color_override("font_color", color)
	if _pixel_font:
		notify.add_theme_font_override("font", _pixel_font)
		notify.add_theme_font_size_override("font_size", 10)
	add_child(notify)

	var tween = create_tween()
	tween.tween_interval(3.0)
	tween.tween_property(notify, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(notify): notify.queue_free()
	)
