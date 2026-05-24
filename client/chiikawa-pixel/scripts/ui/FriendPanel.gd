## FriendPanel.gd ??憟賢?蝟餌絞?Ｘ嚗AY-073嚗?
## DAY-101嚗憓旨?抵??頂蝯?+ 憟賢??????
## 憿舐內憟賢??”?末??瘙末????頛旨?抵???
## 雿蔭嚗opBar ?喳嚗??嚗?
extends Node2D

# ---- 撣豢 ----
const PANEL_WIDTH  := 300
const PANEL_HEIGHT := 240

# ---- 蝭暺???----
var _pixel_font: Font = null
var _is_open: bool = false
var _toggle_btn: Button = null
var _panel_bg: ColorRect = null
var _friend_list_container: Node2D = null
var _pending_badge: Label = null
var _gift_status_label: Label = null

# ---- 憟賢?鞈? ----
var _friends: Array = []
var _pending_count: int = 0
var _gift_sent_today: int = 0
var _gift_remaining: int = 3

# ---- 閮? ----
signal friend_request_sent(target_id: String)

# ---- ????----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_build_toggle_btn()
	_build_panel()
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

## 撱箇?????
func _build_toggle_btn() -> void:
	_toggle_btn = Button.new()
	_toggle_btn.text = "?"
	_toggle_btn.size = Vector2(32, 24)
	_toggle_btn.position = Vector2(0, 0)
	_toggle_btn.flat = true
	_toggle_btn.tooltip_text = "憟賢??”"
	if _pixel_font:
		_toggle_btn.add_theme_font_override("font", _pixel_font)
		_toggle_btn.add_theme_font_size_override("font_size", 14)
	add_child(_toggle_btn)

	# 敺???瘙噬蝡?
	_pending_badge = Label.new()
	_pending_badge.position = Vector2(20, -4)
	_pending_badge.text = ""
	_pending_badge.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))
	if _pixel_font:
		_pending_badge.add_theme_font_override("font", _pixel_font)
		_pending_badge.add_theme_font_size_override("font_size", 9)
	_pending_badge.visible = false
	add_child(_pending_badge)

## 撱箇?銝駁??
func _build_panel() -> void:
	_panel_bg = ColorRect.new()
	_panel_bg.position = Vector2(-PANEL_WIDTH + 32, 28)
	_panel_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_panel_bg.color = Color(0.05, 0.03, 0.15, 0.92)
	_panel_bg.visible = false
	add_child(_panel_bg)

	# 璅?
	var title := Label.new()
	title.position = Vector2(8, 4)
	title.text = "? 憟賢??”"
	title.add_theme_color_override("font_color", Color(0.6, 0.9, 1.0))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 12)
	_panel_bg.add_child(title)

	# ?末????
	var add_btn := Button.new()
	add_btn.text = "嚗?憟賢?"
	add_btn.position = Vector2(PANEL_WIDTH - 80, 2)
	add_btn.size = Vector2(72, 20)
	add_btn.flat = false
	if _pixel_font:
		add_btn.add_theme_font_override("font", _pixel_font)
		add_btn.add_theme_font_size_override("font_size", 9)
	add_btn.pressed.connect(_on_add_friend_pressed)
	_panel_bg.add_child(add_btn)

	# 蝳桃???嚗AY-101嚗?
	_gift_status_label = Label.new()
	_gift_status_label.position = Vector2(8, 22)
	_gift_status_label.text = "?? 隞蝳桃嚗擗?3 甈∴?瘥活 500??嚗?"
	_gift_status_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.4))
	if _pixel_font:
		_gift_status_label.add_theme_font_override("font", _pixel_font)
		_gift_status_label.add_theme_font_size_override("font_size", 9)
	_panel_bg.add_child(_gift_status_label)

	# ??蝺?
	var sep := ColorRect.new()
	sep.position = Vector2(4, 36)
	sep.size = Vector2(PANEL_WIDTH - 8, 1)
	sep.color = Color(0.3, 0.3, 0.5, 0.6)
	_panel_bg.add_child(sep)

	# 憟賢??”摰孵
	_friend_list_container = Node2D.new()
	_friend_list_container.position = Vector2(0, 40)
	_panel_bg.add_child(_friend_list_container)

## ??閮?
func _connect_signals() -> void:
	_toggle_btn.pressed.connect(_on_toggle_pressed)

	if GameManager.has_signal("friend_list_updated"):
		GameManager.friend_list_updated.connect(_on_friend_list_updated)
	if GameManager.has_signal("friend_request_received"):
		GameManager.friend_request_received.connect(_on_friend_request_received)
	if GameManager.has_signal("friend_updated"):
		GameManager.friend_updated.connect(_on_friend_updated)
	# 蝳桃蝟餌絞閮?嚗AY-101嚗?
	if GameManager.has_signal("gift_received"):
		GameManager.gift_received.connect(_on_gift_received)
	if GameManager.has_signal("gift_sent"):
		GameManager.gift_sent.connect(_on_gift_sent)
	if GameManager.has_signal("gift_status"):
		GameManager.gift_status.connect(_on_gift_status)
	if GameManager.has_signal("gift_error"):
		GameManager.gift_error.connect(_on_gift_error)

func _on_toggle_pressed() -> void:
	_is_open = !_is_open
	_panel_bg.visible = _is_open
	if _is_open:
		# ????瘙??啣末??銵?+ 蝳桃???
		NetworkManager.send_message({"type": "get_friend_list", "payload": {}})
		NetworkManager.send_message({"type": "get_gift_status", "payload": {}})

func _on_add_friend_pressed() -> void:
	_show_add_friend_dialog()

func _on_friend_list_updated(data: Dictionary) -> void:
	_friends = data.get("friends", [])
	_pending_count = data.get("pending_count", 0)
	_refresh_ui()

func _on_friend_request_received(data: Dictionary) -> void:
	var from_name = data.get("display_name", data.get("from_id", "?"))
	_show_friend_request_notification(data.get("from_id", ""), from_name)
	_pending_count += 1
	_update_pending_badge()

func _on_friend_updated(data: Dictionary) -> void:
	var event = data.get("event", "")
	var friend_name = data.get("display_name", "")
	match event:
		"online":
			_show_notification("? %s 銝?鈭?" % friend_name, Color(0.4, 0.9, 0.4))
		"offline":
			_show_notification("? %s 銝?鈭? % friend_name, Color(0.6, 0.6, 0.6))"
		"accepted":
			_show_notification("? %s ?亙?鈭??末??瘙?" % friend_name, Color(0.4, 0.8, 1.0))
			NetworkManager.send_message({"type": "get_friend_list", "payload": {}})
		"removed":
			_show_notification("? %s 蝘駁鈭??末?? % friend_name, Color(1.0, 0.5, 0.5))"
			NetworkManager.send_message({"type": "get_friend_list", "payload": {}})

# ---- 蝳桃蝟餌絞 handler嚗AY-101嚗?---

func _on_gift_received(data: Dictionary) -> void:
	var from_name = data.get("display_name", "憟賢?")
	var amount = data.get("amount", 500)
	var new_balance = data.get("new_balance", 0)
	_show_notification("?? %s ?? %d??嚗?擗?嚗?d嚗? % [from_name, amount, new_balance],"
		Color(1.0, 0.85, 0.2))

func _on_gift_sent(data: Dictionary) -> void:
	var to_name = data.get("display_name", "憟賢?")
	var amount = data.get("amount", 500)
	_gift_sent_today = data.get("sent_today", _gift_sent_today)
	_gift_remaining = data.get("remaining", _gift_remaining)
	_update_gift_status_label()
	_show_notification("?? 撌脤?%d?? 蝯?%s嚗?隞?拚? %d 甈∴?" % [amount, to_name, _gift_remaining],
		Color(0.4, 1.0, 0.6))
	# ??渡?憟賢??”嚗?啁旨?拇?????
	_refresh_ui()

func _on_gift_status(data: Dictionary) -> void:
	_gift_sent_today = data.get("sent_today", 0)
	_gift_remaining = data.get("remaining", 3)
	_update_gift_status_label()

func _on_gift_error(data: Dictionary) -> void:
	var msg = data.get("message", "蝳桃?潮仃??)"
	_show_notification("??%s" % msg, Color(1.0, 0.4, 0.4))

func _update_gift_status_label() -> void:
	if not is_instance_valid(_gift_status_label):
		return
	if _gift_remaining > 0:
		_gift_status_label.text = "?? 隞蝳桃嚗擗?%d 甈∴?瘥活 500??嚗? % _gift_remaining"
		_gift_status_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.4))
	else:
		_gift_status_label.text = "?? 隞蝳桃撌脤?嚗??仿?蝵殷?"
		_gift_status_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))

## ?湔 UI
func _refresh_ui() -> void:
	for child in _friend_list_container.get_children():
		child.queue_free()

	_update_pending_badge()

	if _friends.is_empty():
		var empty_label := Label.new()
		empty_label.position = Vector2(8, 4)
		empty_label.text = "???末??敹怠?末?嚗?"
		empty_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.8))
		if _pixel_font:
			empty_label.add_theme_font_override("font", _pixel_font)
			empty_label.add_theme_font_size_override("font_size", 10)
		_friend_list_container.add_child(empty_label)
		return

	# 憿舐內憟賢??”嚗?憭?5 ???征?策蝳桃??嚗?
	var max_show = min(_friends.size(), 5)
	for i in range(max_show):
		var friend_data = _friends[i]
		_build_friend_row(i, friend_data)

## 撱箇?憟賢?銵?DAY-101嚗??亦旨?拇???
func _build_friend_row(index: int, friend_data: Dictionary) -> void:
	var row_y = index * 38
	var is_online = friend_data.get("is_online", false)
	var display_name = friend_data.get("display_name", "?")
	var season_level = friend_data.get("season_level", 0)
	var coins = friend_data.get("coins", 0)
	var friend_id = friend_data.get("player_id", "")

	# 銵??荔?hover ??嚗?
	var row_bg := ColorRect.new()
	row_bg.position = Vector2(4, row_y)
	row_bg.size = Vector2(PANEL_WIDTH - 8, 34)
	row_bg.color = Color(0.1, 0.08, 0.25, 0.5) if index % 2 == 0 else Color(0.08, 0.06, 0.2, 0.3)
	_friend_list_container.add_child(row_bg)

	# ?函????蝷?
	var status_dot := ColorRect.new()
	status_dot.position = Vector2(8, row_y + 13)
	status_dot.size = Vector2(8, 8)
	status_dot.color = Color(0.3, 1.0, 0.3) if is_online else Color(0.5, 0.5, 0.5)
	_friend_list_container.add_child(status_dot)

	# ?迂
	var name_label := Label.new()
	name_label.position = Vector2(20, row_y + 4)
	name_label.text = display_name
	name_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0) if is_online else Color(0.7, 0.7, 0.7))
	if _pixel_font:
		name_label.add_theme_font_override("font", _pixel_font)
		name_label.add_theme_font_size_override("font_size", 10)
	_friend_list_container.add_child(name_label)

	# 鞈賢迤蝑?
	var level_label := Label.new()
	level_label.position = Vector2(20, row_y + 18)
	level_label.text = "Lv%d  ??%d" % [season_level, coins]
	level_label.add_theme_color_override("font_color", Color(0.8, 0.75, 0.5))
	if _pixel_font:
		level_label.add_theme_font_override("font", _pixel_font)
		level_label.add_theme_font_size_override("font_size", 8)
	_friend_list_container.add_child(level_label)

	# 蝳桃??嚗AY-101嚗?
	var gift_btn := Button.new()
	var can_gift = _gift_remaining > 0
	gift_btn.text = "??" if can_gift else "??"
	gift_btn.position = Vector2(PANEL_WIDTH - 56, row_y + 7)
	gift_btn.size = Vector2(24, 20)
	gift_btn.flat = true
	gift_btn.disabled = not can_gift
	gift_btn.tooltip_text = "??500?? 蝳桃" if can_gift else "隞蝳桃撌脤?"
	if _pixel_font:
		gift_btn.add_theme_font_override("font", _pixel_font)
		gift_btn.add_theme_font_size_override("font_size", 11)
	if can_gift:
		gift_btn.pressed.connect(func():
			NetworkManager.send_message({
				"type": "send_gift",
				"payload": {"friend_id": friend_id}
			})
		)
	_friend_list_container.add_child(gift_btn)

	# ???嚗AY-102嚗?
	var challenge_btn := Button.new()
	challenge_btn.text = "??"
	challenge_btn.position = Vector2(PANEL_WIDTH - 80, row_y + 7)
	challenge_btn.size = Vector2(22, 20)
	challenge_btn.flat = true
	challenge_btn.tooltip_text = "?潸絲 1v1 ?嚗陪瘜?1000??嚗?"
	if _pixel_font:
		challenge_btn.add_theme_font_override("font", _pixel_font)
		challenge_btn.add_theme_font_size_override("font_size", 11)
	challenge_btn.pressed.connect(func():
		NetworkManager.send_message({
			"type": "send_challenge_request",
			"payload": {"friend_id": friend_id}
		})
	)
	_friend_list_container.add_child(challenge_btn)

	# ?唾??舀???DAY-103嚗?
	var dm_btn := Button.new()
	dm_btn.text = "?"
	dm_btn.position = Vector2(PANEL_WIDTH - 104, row_y + 7)
	dm_btn.size = Vector2(22, 20)
	dm_btn.flat = true
	dm_btn.tooltip_text = "?喟?閮策 %s" % display_name
	if _pixel_font:
		dm_btn.add_theme_font_override("font", _pixel_font)
		dm_btn.add_theme_font_size_override("font_size", 11)
	dm_btn.pressed.connect(func():
		# ? HUD ?? DM ?Ｘ
		if GameManager.has_signal("open_dm_panel"):
			GameManager.emit_signal("open_dm_panel", friend_id, display_name)
	)
	_friend_list_container.add_child(dm_btn)

	# 蝘駁??
	var remove_btn := Button.new()
	remove_btn.text = "??"
	remove_btn.position = Vector2(PANEL_WIDTH - 28, row_y + 7)
	remove_btn.size = Vector2(20, 20)
	remove_btn.flat = true
	remove_btn.add_theme_color_override("font_color", Color(1.0, 0.4, 0.4))
	if _pixel_font:
		remove_btn.add_theme_font_override("font", _pixel_font)
		remove_btn.add_theme_font_size_override("font_size", 9)
	remove_btn.pressed.connect(func():
		NetworkManager.send_message({
			"type": "remove_friend",
			"payload": {"friend_id": friend_id}
		})
	)
	_friend_list_container.add_child(remove_btn)

## ?湔敺??噬蝡?
func _update_pending_badge() -> void:
	if not is_instance_valid(_pending_badge):
		return
	if _pending_count > 0:
		_pending_badge.text = str(_pending_count)
		_pending_badge.visible = true
	else:
		_pending_badge.visible = false

## 憿舐內?末??閰望?嚗撓?亦摰?ID嚗?
func _show_add_friend_dialog() -> void:
	var dialog_bg := ColorRect.new()
	dialog_bg.position = Vector2(-PANEL_WIDTH + 32, 28 + PANEL_HEIGHT + 4)
	dialog_bg.size = Vector2(PANEL_WIDTH, 50)
	dialog_bg.color = Color(0.08, 0.05, 0.2, 0.95)
	dialog_bg.name = "AddFriendDialog"
	add_child(dialog_bg)

	var hint := Label.new()
	hint.position = Vector2(4, 4)
	hint.text = "頛詨?拙振 ID嚗?8蝣潘?嚗?"
	hint.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
	if _pixel_font:
		hint.add_theme_font_override("font", _pixel_font)
		hint.add_theme_font_size_override("font_size", 9)
	dialog_bg.add_child(hint)

	var line_edit := LineEdit.new()
	line_edit.position = Vector2(4, 18)
	line_edit.size = Vector2(PANEL_WIDTH - 60, 22)
	line_edit.placeholder_text = "?拙振 ID..."
	line_edit.max_length = 36
	if _pixel_font:
		line_edit.add_theme_font_override("font", _pixel_font)
		line_edit.add_theme_font_size_override("font_size", 10)
	dialog_bg.add_child(line_edit)

	var confirm_btn := Button.new()
	confirm_btn.position = Vector2(PANEL_WIDTH - 54, 18)
	confirm_btn.size = Vector2(50, 22)
	confirm_btn.text = "?潮?"
	if _pixel_font:
		confirm_btn.add_theme_font_override("font", _pixel_font)
		confirm_btn.add_theme_font_size_override("font_size", 9)
	dialog_bg.add_child(confirm_btn)

	var send_fn = func():
		var target_id = line_edit.text.strip_edges()
		if target_id.length() >= 4:
			NetworkManager.send_message({
				"type": "send_friend_request",
				"payload": {"target_id": target_id}
			})
			emit_signal("friend_request_sent", target_id)
			_show_notification("憟賢?隢?撌脩??", Color(0.4, 0.9, 0.4))
		if is_instance_valid(dialog_bg):
			dialog_bg.queue_free()

	confirm_btn.pressed.connect(send_fn)
	line_edit.text_submitted.connect(func(_t): send_fn.call())

	var tween = create_tween()
	tween.tween_interval(5.0)
	tween.tween_callback(func():
		if is_instance_valid(dialog_bg):
			dialog_bg.queue_free()
	)

## 憿舐內憟賢?隢??
func _show_friend_request_notification(from_id: String, from_name: String) -> void:
	var notify := Label.new()
	notify.text = "? %s ?喳?雿憟賢?嚗? % from_name"
	notify.position = Vector2(-120, -50)
	notify.add_theme_color_override("font_color", Color(0.4, 0.8, 1.0))
	if _pixel_font:
		notify.add_theme_font_override("font", _pixel_font)
		notify.add_theme_font_size_override("font_size", 10)
	add_child(notify)

	var accept_btn := Button.new()
	accept_btn.text = "???亙?"
	accept_btn.position = Vector2(-120, -30)
	accept_btn.size = Vector2(60, 20)
	if _pixel_font:
		accept_btn.add_theme_font_override("font", _pixel_font)
		accept_btn.add_theme_font_size_override("font_size", 9)
	accept_btn.pressed.connect(func():
		NetworkManager.send_message({
			"type": "accept_friend_request",
			"payload": {"from_id": from_id}
		})
		if is_instance_valid(notify): notify.queue_free()
		if is_instance_valid(accept_btn): accept_btn.queue_free()
	)
	add_child(accept_btn)

	var tween = create_tween()
	tween.tween_interval(8.0)
	tween.tween_callback(func():
		if is_instance_valid(notify): notify.queue_free()
		if is_instance_valid(accept_btn): accept_btn.queue_free()
	)

## 憿舐內?
func _show_notification(text: String, color: Color) -> void:
	var notify := Label.new()
	notify.text = text
	notify.position = Vector2(-120, -30)
	notify.add_theme_color_override("font_color", color)
	if _pixel_font:
		notify.add_theme_font_override("font", _pixel_font)
		notify.add_theme_font_size_override("font_size", 10)
	add_child(notify)

	var tween = create_tween()
	tween.tween_interval(2.5)
	tween.tween_property(notify, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(notify):
			notify.queue_free()
	)
