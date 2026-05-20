## FriendPanel.gd — 好友系統面板（DAY-073）
## 顯示好友列表、好友請求、好友積分比較
## 位置：TopBar 右側（可折疊）
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 280
const PANEL_HEIGHT := 220

# ---- 節點引用 ----
var _pixel_font: Font = null
var _is_open: bool = false
var _toggle_btn: Button = null
var _panel_bg: ColorRect = null
var _friend_list_container: Node2D = null
var _pending_badge: Label = null

# ---- 好友資料 ----
var _friends: Array = []
var _pending_count: int = 0

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

	# 好友列表容器
	_friend_list_container = Node2D.new()
	_friend_list_container.position = Vector2(0, 24)
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

func _on_toggle_pressed() -> void:
	_is_open = !_is_open
	_panel_bg.visible = _is_open
	if _is_open:
		# 開啟時請求最新好友列表
		NetworkManager.send_message({"type": "get_friend_list", "payload": {}})

func _on_add_friend_pressed() -> void:
	# 顯示輸入框（簡單實作：用 OS.get_clipboard 或固定 ID）
	_show_add_friend_dialog()

func _on_friend_list_updated(data: Dictionary) -> void:
	_friends = data.get("friends", [])
	_pending_count = data.get("pending_count", 0)
	_refresh_ui()

func _on_friend_request_received(data: Dictionary) -> void:
	var from_name = data.get("display_name", data.get("from_id", "?"))
	# 顯示好友請求通知
	_show_friend_request_notification(data.get("from_id", ""), from_name)
	# 更新待處理徽章
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

## 更新 UI
func _refresh_ui() -> void:
	# 清除舊的好友列表
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

	# 顯示好友列表（最多 6 個）
	var max_show = min(_friends.size(), 6)
	for i in range(max_show):
		var friend = _friends[i]
		_build_friend_row(i, friend)

## 建立好友行
func _build_friend_row(index: int, friend: Dictionary) -> void:
	var row_y = index * 30
	var is_online = friend.get("is_online", false)
	var display_name = friend.get("display_name", "?")
	var season_level = friend.get("season_level", 0)
	var coins = friend.get("coins", 0)
	var friend_id = friend.get("player_id", "")

	# 在線狀態指示
	var status_dot := ColorRect.new()
	status_dot.position = Vector2(8, row_y + 10)
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
	level_label.position = Vector2(130, row_y + 4)
	level_label.text = "Lv%d" % season_level
	level_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _pixel_font:
		level_label.add_theme_font_override("font", _pixel_font)
		level_label.add_theme_font_size_override("font_size", 9)
	_friend_list_container.add_child(level_label)

	# 金幣
	var coins_label := Label.new()
	coins_label.position = Vector2(165, row_y + 4)
	coins_label.text = "🪙%d" % coins
	coins_label.add_theme_color_override("font_color", Color(0.9, 0.8, 0.4))
	if _pixel_font:
		coins_label.add_theme_font_override("font", _pixel_font)
		coins_label.add_theme_font_size_override("font_size", 9)
	_friend_list_container.add_child(coins_label)

	# 移除按鈕
	var remove_btn := Button.new()
	remove_btn.text = "✕"
	remove_btn.position = Vector2(PANEL_WIDTH - 28, row_y + 4)
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

## 顯示加好友對話框（簡單版：輸入玩家 ID）
func _show_add_friend_dialog() -> void:
	# 建立簡單輸入提示
	var hint := Label.new()
	hint.text = "請在聊天輸入對方的玩家 ID\n（功能開發中）"
	hint.position = Vector2(8, -40)
	hint.add_theme_color_override("font_color", Color(0.8, 0.9, 1.0))
	if _pixel_font:
		hint.add_theme_font_override("font", _pixel_font)
		hint.add_theme_font_size_override("font_size", 10)
	add_child(hint)

	var tween = create_tween()
	tween.tween_interval(2.0)
	tween.tween_property(hint, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(hint):
			hint.queue_free()
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

	# 接受/拒絕按鈕
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
	notify.position = Vector2(-100, -30)
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
