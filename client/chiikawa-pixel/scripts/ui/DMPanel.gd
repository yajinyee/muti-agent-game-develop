## DMPanel.gd — 玩家私訊面板（DAY-103）
## 好友間可以互相發送私訊，離線訊息上線後自動顯示
## 整合到 FriendPanel（點擊好友名稱開啟對話框）
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 320
const PANEL_HEIGHT := 280
const MAX_DISPLAY_MESSAGES := 20

# ---- 節點引用 ----
var _pixel_font: Font = null
var _panel_bg: ColorRect = null
var _messages_container: Node2D = null
var _input_field: LineEdit = null
var _send_btn: Button = null
var _title_label: Label = null
var _unread_badge: Label = null
var _toggle_btn: Button = null

# ---- 狀態 ----
var _is_open: bool = false
var _current_friend_id: String = ""
var _current_friend_name: String = ""
var _messages: Array = []  # {from_id, from_name, content, sent_at, is_offline}
var _unread_count: int = 0

# ---- 訊號 ----
signal dm_opened(friend_id: String)

# ---- 初始化 ----
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

## 建立折疊按鈕（顯示在 TopBar）
func _build_toggle_btn() -> void:
	_toggle_btn = Button.new()
	_toggle_btn.text = "💬"
	_toggle_btn.size = Vector2(32, 24)
	_toggle_btn.position = Vector2(0, 0)
	_toggle_btn.flat = true
	_toggle_btn.tooltip_text = "私訊"
	if _pixel_font:
		_toggle_btn.add_theme_font_override("font", _pixel_font)
		_toggle_btn.add_theme_font_size_override("font_size", 14)
	add_child(_toggle_btn)

	# 未讀徽章
	_unread_badge = Label.new()
	_unread_badge.position = Vector2(20, -4)
	_unread_badge.text = ""
	_unread_badge.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))
	if _pixel_font:
		_unread_badge.add_theme_font_override("font", _pixel_font)
		_unread_badge.add_theme_font_size_override("font_size", 9)
	_unread_badge.visible = false
	add_child(_unread_badge)

## 建立主面板
func _build_panel() -> void:
	_panel_bg = ColorRect.new()
	_panel_bg.position = Vector2(-PANEL_WIDTH + 32, 28)
	_panel_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_panel_bg.color = Color(0.04, 0.02, 0.14, 0.95)
	_panel_bg.visible = false
	add_child(_panel_bg)

	# 標題
	_title_label = Label.new()
	_title_label.position = Vector2(8, 4)
	_title_label.text = "💬 私訊"
	_title_label.add_theme_color_override("font_color", Color(0.6, 0.9, 1.0))
	if _pixel_font:
		_title_label.add_theme_font_override("font", _pixel_font)
		_title_label.add_theme_font_size_override("font_size", 12)
	_panel_bg.add_child(_title_label)

	# 關閉按鈕
	var close_btn := Button.new()
	close_btn.text = "✕"
	close_btn.position = Vector2(PANEL_WIDTH - 24, 2)
	close_btn.size = Vector2(20, 20)
	close_btn.flat = true
	close_btn.add_theme_color_override("font_color", Color(0.8, 0.4, 0.4))
	if _pixel_font:
		close_btn.add_theme_font_override("font", _pixel_font)
		close_btn.add_theme_font_size_override("font_size", 10)
	close_btn.pressed.connect(func(): _close_panel())
	_panel_bg.add_child(close_btn)

	# 分隔線
	var sep := ColorRect.new()
	sep.position = Vector2(4, 22)
	sep.size = Vector2(PANEL_WIDTH - 8, 1)
	sep.color = Color(0.3, 0.3, 0.5, 0.6)
	_panel_bg.add_child(sep)

	# 訊息容器（捲動區域）
	_messages_container = Node2D.new()
	_messages_container.position = Vector2(0, 26)
	_panel_bg.add_child(_messages_container)

	# 輸入框
	_input_field = LineEdit.new()
	_input_field.position = Vector2(4, PANEL_HEIGHT - 30)
	_input_field.size = Vector2(PANEL_WIDTH - 60, 24)
	_input_field.placeholder_text = "輸入訊息..."
	_input_field.max_length = 200
	if _pixel_font:
		_input_field.add_theme_font_override("font", _pixel_font)
		_input_field.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_input_field)

	# 發送按鈕
	_send_btn = Button.new()
	_send_btn.text = "發送"
	_send_btn.position = Vector2(PANEL_WIDTH - 54, PANEL_HEIGHT - 30)
	_send_btn.size = Vector2(50, 24)
	if _pixel_font:
		_send_btn.add_theme_font_override("font", _pixel_font)
		_send_btn.add_theme_font_size_override("font_size", 9)
	_send_btn.pressed.connect(_on_send_pressed)
	_panel_bg.add_child(_send_btn)

	# 輸入框 Enter 發送
	_input_field.text_submitted.connect(func(_t): _on_send_pressed())

## 連接訊號
func _connect_signals() -> void:
	_toggle_btn.pressed.connect(_on_toggle_pressed)

	if GameManager.has_signal("dm_received"):
		GameManager.dm_received.connect(_on_dm_received)
	if GameManager.has_signal("dm_sent"):
		GameManager.dm_sent.connect(_on_dm_sent)
	if GameManager.has_signal("dm_error"):
		GameManager.dm_error.connect(_on_dm_error)

## 開啟與特定好友的對話
func open_conversation(friend_id: String, friend_name: String) -> void:
	_current_friend_id = friend_id
	_current_friend_name = friend_name
	_title_label.text = "💬 與 %s 的對話" % friend_name
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
		# 顯示好友選擇提示
		_show_friend_select_hint()

func _show_friend_select_hint() -> void:
	var hint := Label.new()
	hint.text = "請從好友列表點擊好友名稱開啟對話"
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

## 收到私訊
func _on_dm_received(data: Dictionary) -> void:
	var from_id = data.get("from_id", "")
	var from_name = data.get("from_name", "好友")
	var content = data.get("content", "")
	var is_offline = data.get("is_offline", false)

	# 加入訊息記錄
	_messages.append({
		"from_id": from_id,
		"from_name": from_name,
		"content": content,
		"sent_at": data.get("sent_at", 0),
		"is_mine": false,
		"is_offline": is_offline
	})

	# 如果面板開著且是當前對話，刷新顯示
	if _is_open and from_id == _current_friend_id:
		_refresh_messages()
	else:
		# 增加未讀計數
		_unread_count += 1
		_update_unread_badge()
		# 顯示浮動通知
		var prefix = "📬 [離線]" if is_offline else "💬"
		_show_notification("%s %s：%s" % [prefix, from_name, content.substr(0, 30)],
			Color(0.4, 0.8, 1.0))

## 發送成功確認
func _on_dm_sent(data: Dictionary) -> void:
	var remaining = data.get("remaining", 0)
	# 加入自己的訊息到記錄（從 input_field 取得，但已清空，所以從 payload 取）
	# 這裡只更新計數顯示
	if remaining <= 10:
		_show_notification("💬 今日剩餘 %d 則訊息" % remaining, Color(0.8, 0.8, 0.6))

## 發送失敗
func _on_dm_error(data: Dictionary) -> void:
	var msg = data.get("message", "發送失敗")
	_show_notification("❌ %s" % msg, Color(1.0, 0.4, 0.4))

## 刷新訊息顯示
func _refresh_messages() -> void:
	for child in _messages_container.get_children():
		child.queue_free()

	# 只顯示與當前好友的對話
	var filtered = []
	for m in _messages:
		if m.get("from_id") == _current_friend_id or m.get("is_mine", false):
			filtered.append(m)

	# 最多顯示最新 MAX_DISPLAY_MESSAGES 則
	var start = max(0, filtered.size() - MAX_DISPLAY_MESSAGES)
	var y_offset = 0
	for i in range(start, filtered.size()):
		var m = filtered[i]
		y_offset = _build_message_bubble(y_offset, m)

	# 清除未讀
	_unread_count = 0
	_update_unread_badge()

## 建立訊息氣泡
func _build_message_bubble(y: int, msg_data: Dictionary) -> int:
	var is_mine = msg_data.get("is_mine", false)
	var from_name = msg_data.get("from_name", "?")
	var content = msg_data.get("content", "")
	var is_offline = msg_data.get("is_offline", false)

	# 名稱標籤
	var name_label := Label.new()
	name_label.position = Vector2(8 if not is_mine else PANEL_WIDTH - 120, y)
	name_label.text = ("你" if is_mine else from_name) + ("（離線）" if is_offline else "")
	name_label.add_theme_color_override("font_color",
		Color(0.6, 0.9, 1.0) if is_mine else Color(1.0, 0.85, 0.5))
	if _pixel_font:
		name_label.add_theme_font_override("font", _pixel_font)
		name_label.add_theme_font_size_override("font_size", 8)
	_messages_container.add_child(name_label)
	y += 12

	# 訊息氣泡背景
	var bubble_bg := ColorRect.new()
	var bubble_w = min(int(len(content) * 6 + 16), PANEL_WIDTH - 24)
	bubble_bg.position = Vector2(8 if not is_mine else PANEL_WIDTH - bubble_w - 8, y)
	bubble_bg.size = Vector2(bubble_w, 20)
	bubble_bg.color = Color(0.15, 0.1, 0.35) if is_mine else Color(0.1, 0.15, 0.3)
	_messages_container.add_child(bubble_bg)

	# 訊息文字
	var content_label := Label.new()
	content_label.position = Vector2(4, 2)
	content_label.text = content
	content_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	if _pixel_font:
		content_label.add_theme_font_override("font", _pixel_font)
		content_label.add_theme_font_size_override("font_size", 9)
	bubble_bg.add_child(content_label)

	return y + 26

## 更新未讀徽章
func _update_unread_badge() -> void:
	if not is_instance_valid(_unread_badge):
		return
	if _unread_count > 0:
		_unread_badge.text = str(_unread_count)
		_unread_badge.visible = true
	else:
		_unread_badge.visible = false

## 顯示浮動通知
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
