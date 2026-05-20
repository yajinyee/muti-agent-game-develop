## FriendPanel.gd — 好友系統面板（DAY-073）
## DAY-101：新增禮物贈送系統 + 好友持久化支援
## 顯示好友列表、好友請求、好友積分比較、禮物贈送
## 位置：TopBar 右側（可折疊）
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 300
const PANEL_HEIGHT := 240

# ---- 節點引用 ----
var _pixel_font: Font = null
var _is_open: bool = false
var _toggle_btn: Button = null
var _panel_bg: ColorRect = null
var _friend_list_container: Node2D = null
var _pending_badge: Label = null
var _gift_status_label: Label = null

# ---- 好友資料 ----
var _friends: Array = []
var _pending_count: int = 0
var _gift_sent_today: int = 0
var _gift_remaining: int = 3

# ---- 訊號 ----
signal friend_request_sent(target_id: String)

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_build_toggle_btn()
	_build_panel()
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

## 建立折疊按鈕
func _build_toggle_btn() -> void:
	_toggle_btn = Button.new()
	_toggle_btn.text = "👥"
	_toggle_btn.size = Vector2(32, 24)
	_toggle_btn.position = Vector2(0, 0)
	_toggle_btn.flat = true
	_toggle_btn.tooltip_text = "好友列表"
	if _pixel_font:
		_toggle_btn.add_theme_font_override("font", _pixel_font)
		_toggle_btn.add_theme_font_size_override("font_size", 14)
	add_child(_toggle_btn)

	# 待處理請求徽章
	_pending_badge = Label.new()
	_pending_badge.position = Vector2(20, -4)
	_pending_badge.text = ""
	_pending_badge.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))
	if _pixel_font:
		_pending_badge.add_theme_font_override("font", _pixel_font)
		_pending_badge.add_theme_font_size_override("font_size", 9)
	_pending_badge.visible = false
	add_child(_pending_badge)

## 建立主面板
func _build_panel() -> void:
	_panel_bg = ColorRect.new()
	_panel_bg.position = Vector2(-PANEL_WIDTH + 32, 28)
	_panel_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_panel_bg.color = Color(0.05, 0.03, 0.15, 0.92)
	_panel_bg.visible = false
	add_child(_panel_bg)

	# 標題
	var title := Label.new()
	title.position = Vector2(8, 4)
	title.text = "👥 好友列表"
	title.add_theme_color_override("font_color", Color(0.6, 0.9, 1.0))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 12)
	_panel_bg.add_child(title)

	# 加好友按鈕
	var add_btn := Button.new()
	add_btn.text = "＋加好友"
	add_btn.position = Vector2(PANEL_WIDTH - 80, 2)
	add_btn.size = Vector2(72, 20)
	add_btn.flat = false
	if _pixel_font:
		add_btn.add_theme_font_override("font", _pixel_font)
		add_btn.add_theme_font_size_override("font_size", 9)
	add_btn.pressed.connect(_on_add_friend_pressed)
	_panel_bg.add_child(add_btn)

	# 禮物狀態列（DAY-101）
	_gift_status_label = Label.new()
	_gift_status_label.position = Vector2(8, 22)
	_gift_status_label.text = "🎁 今日禮物：剩餘 3 次（每次 500🪙）"
	_gift_status_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.4))
	if _pixel_font:
		_gift_status_label.add_theme_font_override("font", _pixel_font)
		_gift_status_label.add_theme_font_size_override("font_size", 9)
	_panel_bg.add_child(_gift_status_label)

	# 分隔線
	var sep := ColorRect.new()
	sep.position = Vector2(4, 36)
	sep.size = Vector2(PANEL_WIDTH - 8, 1)
	sep.color = Color(0.3, 0.3, 0.5, 0.6)
	_panel_bg.add_child(sep)

	# 好友列表容器
	_friend_list_container = Node2D.new()
	_friend_list_container.position = Vector2(0, 40)
	_panel_bg.add_child(_friend_list_container)

## 連接訊號
func _connect_signals() -> void:
	_toggle_btn.pressed.connect(_on_toggle_pressed)

	if GameManager.has_signal("friend_list_updated"):
		GameManager.friend_list_updated.connect(_on_friend_list_updated)
	if GameManager.has_signal("friend_request_received"):
		GameManager.friend_request_received.connect(_on_friend_request_received)
	if GameManager.has_signal("friend_updated"):
		GameManager.friend_updated.connect(_on_friend_updated)
	# 禮物系統訊號（DAY-101）
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
		# 開啟時請求最新好友列表 + 禮物狀態
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
			_show_notification("👥 %s 上線了！" % friend_name, Color(0.4, 0.9, 0.4))
		"offline":
			_show_notification("👥 %s 下線了" % friend_name, Color(0.6, 0.6, 0.6))
		"accepted":
			_show_notification("👥 %s 接受了你的好友請求！" % friend_name, Color(0.4, 0.8, 1.0))
			NetworkManager.send_message({"type": "get_friend_list", "payload": {}})
		"removed":
			_show_notification("👥 %s 移除了你的好友" % friend_name, Color(1.0, 0.5, 0.5))
			NetworkManager.send_message({"type": "get_friend_list", "payload": {}})

# ---- 禮物系統 handler（DAY-101）----

func _on_gift_received(data: Dictionary) -> void:
	var from_name = data.get("display_name", "好友")
	var amount = data.get("amount", 500)
	var new_balance = data.get("new_balance", 0)
	_show_notification("🎁 %s 送你 %d🪙！（餘額：%d）" % [from_name, amount, new_balance],
		Color(1.0, 0.85, 0.2))

func _on_gift_sent(data: Dictionary) -> void:
	var to_name = data.get("display_name", "好友")
	var amount = data.get("amount", 500)
	_gift_sent_today = data.get("sent_today", _gift_sent_today)
	_gift_remaining = data.get("remaining", _gift_remaining)
	_update_gift_status_label()
	_show_notification("🎁 已送 %d🪙 給 %s！（今日剩餘 %d 次）" % [amount, to_name, _gift_remaining],
		Color(0.4, 1.0, 0.6))
	# 重新整理好友列表（更新禮物按鈕狀態）
	_refresh_ui()

func _on_gift_status(data: Dictionary) -> void:
	_gift_sent_today = data.get("sent_today", 0)
	_gift_remaining = data.get("remaining", 3)
	_update_gift_status_label()

func _on_gift_error(data: Dictionary) -> void:
	var msg = data.get("message", "禮物發送失敗")
	_show_notification("❌ %s" % msg, Color(1.0, 0.4, 0.4))

func _update_gift_status_label() -> void:
	if not is_instance_valid(_gift_status_label):
		return
	if _gift_remaining > 0:
		_gift_status_label.text = "🎁 今日禮物：剩餘 %d 次（每次 500🪙）" % _gift_remaining
		_gift_status_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.4))
	else:
		_gift_status_label.text = "🎁 今日禮物已送完（明日重置）"
		_gift_status_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))

## 更新 UI
func _refresh_ui() -> void:
	for child in _friend_list_container.get_children():
		child.queue_free()

	_update_pending_badge()

	if _friends.is_empty():
		var empty_label := Label.new()
		empty_label.position = Vector2(8, 4)
		empty_label.text = "還沒有好友，快去加好友吧！"
		empty_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.8))
		if _pixel_font:
			empty_label.add_theme_font_override("font", _pixel_font)
			empty_label.add_theme_font_size_override("font_size", 10)
		_friend_list_container.add_child(empty_label)
		return

	# 顯示好友列表（最多 5 個，留空間給禮物按鈕）
	var max_show = min(_friends.size(), 5)
	for i in range(max_show):
		var friend_data = _friends[i]
		_build_friend_row(i, friend_data)

## 建立好友行（DAY-101：加入禮物按鈕）
func _build_friend_row(index: int, friend_data: Dictionary) -> void:
	var row_y = index * 38
	var is_online = friend_data.get("is_online", false)
	var display_name = friend_data.get("display_name", "?")
	var season_level = friend_data.get("season_level", 0)
	var coins = friend_data.get("coins", 0)
	var friend_id = friend_data.get("player_id", "")

	# 行背景（hover 效果）
	var row_bg := ColorRect.new()
	row_bg.position = Vector2(4, row_y)
	row_bg.size = Vector2(PANEL_WIDTH - 8, 34)
	row_bg.color = Color(0.1, 0.08, 0.25, 0.5) if index % 2 == 0 else Color(0.08, 0.06, 0.2, 0.3)
	_friend_list_container.add_child(row_bg)

	# 在線狀態指示
	var status_dot := ColorRect.new()
	status_dot.position = Vector2(8, row_y + 13)
	status_dot.size = Vector2(8, 8)
	status_dot.color = Color(0.3, 1.0, 0.3) if is_online else Color(0.5, 0.5, 0.5)
	_friend_list_container.add_child(status_dot)

	# 名稱
	var name_label := Label.new()
	name_label.position = Vector2(20, row_y + 4)
	name_label.text = display_name
	name_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0) if is_online else Color(0.7, 0.7, 0.7))
	if _pixel_font:
		name_label.add_theme_font_override("font", _pixel_font)
		name_label.add_theme_font_size_override("font_size", 10)
	_friend_list_container.add_child(name_label)

	# 賽季等級
	var level_label := Label.new()
	level_label.position = Vector2(20, row_y + 18)
	level_label.text = "Lv%d  🪙%d" % [season_level, coins]
	level_label.add_theme_color_override("font_color", Color(0.8, 0.75, 0.5))
	if _pixel_font:
		level_label.add_theme_font_override("font", _pixel_font)
		level_label.add_theme_font_size_override("font_size", 8)
	_friend_list_container.add_child(level_label)

	# 禮物按鈕（DAY-101）
	var gift_btn := Button.new()
	var can_gift = _gift_remaining > 0
	gift_btn.text = "🎁" if can_gift else "✗"
	gift_btn.position = Vector2(PANEL_WIDTH - 56, row_y + 7)
	gift_btn.size = Vector2(24, 20)
	gift_btn.flat = true
	gift_btn.disabled = not can_gift
	gift_btn.tooltip_text = "送 500🪙 禮物" if can_gift else "今日禮物已送完"
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

	# 移除按鈕
	var remove_btn := Button.new()
	remove_btn.text = "✕"
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

## 更新待處理徽章
func _update_pending_badge() -> void:
	if not is_instance_valid(_pending_badge):
		return
	if _pending_count > 0:
		_pending_badge.text = str(_pending_count)
		_pending_badge.visible = true
	else:
		_pending_badge.visible = false

## 顯示加好友對話框（輸入玩家 ID）
func _show_add_friend_dialog() -> void:
	var dialog_bg := ColorRect.new()
	dialog_bg.position = Vector2(-PANEL_WIDTH + 32, 28 + PANEL_HEIGHT + 4)
	dialog_bg.size = Vector2(PANEL_WIDTH, 50)
	dialog_bg.color = Color(0.08, 0.05, 0.2, 0.95)
	dialog_bg.name = "AddFriendDialog"
	add_child(dialog_bg)

	var hint := Label.new()
	hint.position = Vector2(4, 4)
	hint.text = "輸入玩家 ID（前8碼）："
	hint.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
	if _pixel_font:
		hint.add_theme_font_override("font", _pixel_font)
		hint.add_theme_font_size_override("font_size", 9)
	dialog_bg.add_child(hint)

	var line_edit := LineEdit.new()
	line_edit.position = Vector2(4, 18)
	line_edit.size = Vector2(PANEL_WIDTH - 60, 22)
	line_edit.placeholder_text = "玩家 ID..."
	line_edit.max_length = 36
	if _pixel_font:
		line_edit.add_theme_font_override("font", _pixel_font)
		line_edit.add_theme_font_size_override("font_size", 10)
	dialog_bg.add_child(line_edit)

	var confirm_btn := Button.new()
	confirm_btn.position = Vector2(PANEL_WIDTH - 54, 18)
	confirm_btn.size = Vector2(50, 22)
	confirm_btn.text = "發送"
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
			_show_notification("好友請求已發送！", Color(0.4, 0.9, 0.4))
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

## 顯示好友請求通知
func _show_friend_request_notification(from_id: String, from_name: String) -> void:
	var notify := Label.new()
	notify.text = "👥 %s 想加你為好友！" % from_name
	notify.position = Vector2(-120, -50)
	notify.add_theme_color_override("font_color", Color(0.4, 0.8, 1.0))
	if _pixel_font:
		notify.add_theme_font_override("font", _pixel_font)
		notify.add_theme_font_size_override("font_size", 10)
	add_child(notify)

	var accept_btn := Button.new()
	accept_btn.text = "✓ 接受"
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

## 顯示通知
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
