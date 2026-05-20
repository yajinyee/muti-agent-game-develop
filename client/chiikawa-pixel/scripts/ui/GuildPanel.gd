## GuildPanel.gd — 公會系統面板（DAY-074/075）
## 顯示公會資訊、成員列表、公會任務進度、公會聊天室
## 位置：TopBar 右側（可折疊）
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 300
const PANEL_HEIGHT := 320
const MAX_CHAT_MESSAGES := 20

# ---- 節點引用 ----
var _pixel_font: Font = null
var _is_open: bool = false
var _toggle_btn: Button = null
var _panel_bg: ColorRect = null
var _content_container: Node2D = null
var _task_badge: Label = null
var _chat_container: Node2D = null
var _chat_input: LineEdit = null
var _chat_messages_node: Node2D = null

# ---- 公會資料 ----
var _guild_id: String = ""
var _guild_name: String = ""
var _guild_level: int = 0
var _guild_exp: int = 0
var _members: Array = []
var _tasks: Array = []
var _my_role: String = ""
var _total_kills: int = 0
var _total_coins: int = 0

# ---- 聊天記錄 ----
var _chat_history: Array = []  # Array of {name, role, message, timestamp}

# ---- 任務完成通知佇列 ----
var _task_complete_queue: Array = []

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
	_toggle_btn.text = "⚔️"
	_toggle_btn.size = Vector2(32, 24)
	_toggle_btn.position = Vector2(0, 0)
	_toggle_btn.flat = true
	_toggle_btn.tooltip_text = "公會"
	if _pixel_font:
		_toggle_btn.add_theme_font_override("font", _pixel_font)
		_toggle_btn.add_theme_font_size_override("font_size", 14)
	add_child(_toggle_btn)

	# 任務完成徽章
	_task_badge = Label.new()
	_task_badge.position = Vector2(20, -4)
	_task_badge.text = ""
	_task_badge.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	if _pixel_font:
		_task_badge.add_theme_font_override("font", _pixel_font)
		_task_badge.add_theme_font_size_override("font_size", 9)
	_task_badge.visible = false
	add_child(_task_badge)

## 建立主面板
func _build_panel() -> void:
	_panel_bg = ColorRect.new()
	_panel_bg.position = Vector2(-PANEL_WIDTH + 32, 28)
	_panel_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_panel_bg.color = Color(0.05, 0.08, 0.03, 0.92)
	_panel_bg.visible = false
	add_child(_panel_bg)

	# 標題
	var title := Label.new()
	title.position = Vector2(8, 4)
	title.text = "⚔️ 公會"
	title.add_theme_color_override("font_color", Color(0.9, 0.85, 0.3))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 12)
	_panel_bg.add_child(title)

	# 建立/加入公會按鈕（無公會時顯示）
	var create_btn := Button.new()
	create_btn.name = "CreateBtn"
	create_btn.text = "建立公會"
	create_btn.size = Vector2(90, 20)
	create_btn.position = Vector2(8, 22)
	create_btn.flat = false
	if _pixel_font:
		create_btn.add_theme_font_override("font", _pixel_font)
		create_btn.add_theme_font_size_override("font_size", 10)
	create_btn.pressed.connect(_on_create_guild_pressed)
	_panel_bg.add_child(create_btn)

	var join_btn := Button.new()
	join_btn.name = "JoinBtn"
	join_btn.text = "加入公會"
	join_btn.size = Vector2(90, 20)
	join_btn.position = Vector2(104, 22)
	join_btn.flat = false
	if _pixel_font:
		join_btn.add_theme_font_override("font", _pixel_font)
		join_btn.add_theme_font_size_override("font_size", 10)
	join_btn.pressed.connect(_on_join_guild_pressed)
	_panel_bg.add_child(join_btn)

	# 公會資訊區（有公會時顯示）
	_content_container = Node2D.new()
	_content_container.name = "ContentContainer"
	_content_container.visible = false
	_panel_bg.add_child(_content_container)

	# 退出公會按鈕
	var leave_btn := Button.new()
	leave_btn.name = "LeaveBtn"
	leave_btn.text = "退出公會"
	leave_btn.size = Vector2(80, 18)
	leave_btn.position = Vector2(PANEL_WIDTH - 90, 4)
	leave_btn.flat = false
	if _pixel_font:
		leave_btn.add_theme_font_override("font", _pixel_font)
		leave_btn.add_theme_font_size_override("font_size", 9)
	leave_btn.pressed.connect(_on_leave_guild_pressed)
	_content_container.add_child(leave_btn)

## 連接訊號
func _connect_signals() -> void:
	if _toggle_btn:
		_toggle_btn.pressed.connect(_on_toggle)
	# 連接 GameManager 訊號
	if GameManager:
		if not GameManager.guild_updated.is_connected(_on_guild_updated):
			GameManager.guild_updated.connect(_on_guild_updated)
		if not GameManager.guild_task_complete.is_connected(_on_guild_task_complete):
			GameManager.guild_task_complete.connect(_on_guild_task_complete)
		if not GameManager.guild_message_received.is_connected(_on_guild_message):
			GameManager.guild_message_received.connect(_on_guild_message)

## 切換面板顯示
func _on_toggle() -> void:
	_is_open = !_is_open
	_panel_bg.visible = _is_open
	if _is_open:
		# 請求最新公會資訊
		GameManager.send_message("get_guild_info", {})
		if _guild_id == "":
			# 同時請求公會列表
			GameManager.send_message("get_guild_list", {})

## 公會資訊更新
func _on_guild_updated(data: Dictionary) -> void:
	_guild_id = data.get("guild_id", "")
	_guild_name = data.get("name", "")
	_guild_level = data.get("level", 0)
	_guild_exp = data.get("exp", 0)
	_members = data.get("members", [])
	_tasks = data.get("tasks", [])
	_my_role = data.get("my_role", "")
	_total_kills = data.get("total_kills", 0)
	_total_coins = data.get("total_coins", 0)

	if _is_open:
		_refresh_panel()

## 公會任務完成通知
func _on_guild_task_complete(data: Dictionary) -> void:
	_task_complete_queue.append(data)
	_show_task_badge()
	_show_task_complete_popup(data)

## 公會聊天訊息接收（DAY-075）
func _on_guild_message(data: Dictionary) -> void:
	_chat_history.append(data)
	if _chat_history.size() > MAX_CHAT_MESSAGES:
		_chat_history.pop_front()
	if _is_open and _chat_messages_node:
		_refresh_chat_messages()
	# 未開啟時顯示徽章提示
	if not _is_open:
		_task_badge.text = "💬"
		_task_badge.visible = true

## 顯示任務完成徽章
func _show_task_badge() -> void:
	_task_badge.text = "✓"
	_task_badge.visible = true
	# 3 秒後隱藏
	var timer := get_tree().create_timer(3.0)
	timer.timeout.connect(func(): _task_badge.visible = false)

## 顯示任務完成彈窗
func _show_task_complete_popup(data: Dictionary) -> void:
	var task_name: String = data.get("task_name", "公會任務")
	var task_icon: String = data.get("task_icon", "⚔️")
	var reward: int = data.get("reward", 0)
	var new_balance: int = data.get("new_balance", 0)

	# 使用 HUD 的成就通知系統
	if HUD and HUD.has_method("show_achievement_notify"):
		HUD.show_achievement_notify(
			"guild_task",
			task_icon + " 公會任務完成！",
			task_name + "  +" + str(reward) + " 金幣",
			Color(0.9, 0.85, 0.3)
		)

## 刷新面板內容
func _refresh_panel() -> void:
	if not _panel_bg:
		return

	var has_guild := _guild_id != ""

	# 顯示/隱藏建立加入按鈕
	var create_btn = _panel_bg.get_node_or_null("CreateBtn")
	var join_btn = _panel_bg.get_node_or_null("JoinBtn")
	if create_btn:
		create_btn.visible = not has_guild
	if join_btn:
		join_btn.visible = not has_guild

	# 顯示/隱藏公會內容
	if _content_container:
		_content_container.visible = has_guild

	if not has_guild:
		return

	# 清除舊內容（保留 LeaveBtn）
	for child in _content_container.get_children():
		if child.name != "LeaveBtn":
			child.queue_free()

	# 公會名稱 + 等級
	var name_label := Label.new()
	name_label.position = Vector2(8, 22)
	name_label.text = _guild_name + "  Lv." + str(_guild_level)
	name_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	if _pixel_font:
		name_label.add_theme_font_override("font", _pixel_font)
		name_label.add_theme_font_size_override("font_size", 11)
	_content_container.add_child(name_label)

	# 公會統計
	var stats_label := Label.new()
	stats_label.position = Vector2(8, 36)
	stats_label.text = "討伐: " + str(_total_kills) + "  金幣: " + str(_total_coins)
	stats_label.add_theme_color_override("font_color", Color(0.7, 0.9, 0.7))
	if _pixel_font:
		stats_label.add_theme_font_override("font", _pixel_font)
		stats_label.add_theme_font_size_override("font_size", 9)
	_content_container.add_child(stats_label)

	# 成員列表標題
	var member_title := Label.new()
	member_title.position = Vector2(8, 50)
	member_title.text = "成員（" + str(_members.size()) + "/20）"
	member_title.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	if _pixel_font:
		member_title.add_theme_font_override("font", _pixel_font)
		member_title.add_theme_font_size_override("font_size", 10)
	_content_container.add_child(member_title)

	# 成員列表（最多顯示 5 個）
	var display_count := mini(_members.size(), 5)
	for i in range(display_count):
		var m: Dictionary = _members[i]
		var role_icon := _get_role_icon(m.get("role", "member"))
		var online_dot := "🟢" if m.get("is_online", false) else "⚫"
		var name_str: String = m.get("display_name", "???")
		var contrib: int = m.get("contribution", 0)

		var member_label := Label.new()
		member_label.position = Vector2(8, 64 + i * 16)
		member_label.text = online_dot + role_icon + " " + name_str + "  +" + str(contrib)
		var color := Color(1.0, 0.9, 0.5) if m.get("role", "") == "leader" else Color(0.85, 0.85, 0.85)
		member_label.add_theme_color_override("font_color", color)
		if _pixel_font:
			member_label.add_theme_font_override("font", _pixel_font)
			member_label.add_theme_font_size_override("font_size", 9)
		_content_container.add_child(member_label)

	# 任務列表標題
	var task_y := 64 + display_count * 16 + 8
	var task_title := Label.new()
	task_title.position = Vector2(8, task_y)
	task_title.text = "公會任務"
	task_title.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	if _pixel_font:
		task_title.add_theme_font_override("font", _pixel_font)
		task_title.add_theme_font_size_override("font_size", 10)
	_content_container.add_child(task_title)

	# 任務列表
	for i in range(_tasks.size()):
		var t: Dictionary = _tasks[i]
		var completed: bool = t.get("completed", false)
		var current: int = t.get("current", 0)
		var target_val: int = t.get("target", 1)
		var progress_pct := float(current) / float(max(target_val, 1))

		var task_label := Label.new()
		task_label.position = Vector2(8, task_y + 14 + i * 28)
		var status_icon := "✅" if completed else "⬜"
		task_label.text = status_icon + " " + t.get("icon", "") + " " + t.get("name", "")
		var task_color := Color(0.5, 0.9, 0.5) if completed else Color(0.9, 0.9, 0.9)
		task_label.add_theme_color_override("font_color", task_color)
		if _pixel_font:
			task_label.add_theme_font_override("font", _pixel_font)
			task_label.add_theme_font_size_override("font_size", 9)
		_content_container.add_child(task_label)

		# 進度條
		var bar_bg := ColorRect.new()
		bar_bg.position = Vector2(8, task_y + 26 + i * 28)
		bar_bg.size = Vector2(PANEL_WIDTH - 20, 6)
		bar_bg.color = Color(0.2, 0.2, 0.2)
		_content_container.add_child(bar_bg)

		var bar_fill := ColorRect.new()
		bar_fill.position = Vector2(8, task_y + 26 + i * 28)
		bar_fill.size = Vector2((PANEL_WIDTH - 20) * progress_pct, 6)
		bar_fill.color = Color(0.3, 0.9, 0.3) if completed else Color(0.9, 0.7, 0.1)
		_content_container.add_child(bar_fill)

	# 聊天室（DAY-075）
	var chat_y := task_y + 14 + _tasks.size() * 28 + 8
	_build_chat_area(_content_container, chat_y)
	_refresh_chat_messages()

## 取得職位圖示
func _get_role_icon(role: String) -> String:
	match role:
		"leader":  return "👑"
		"officer": return "⭐"
		_:         return "👤"

## 建立公會按鈕
func _on_create_guild_pressed() -> void:
	# 簡單輸入：使用玩家名稱作為公會名稱
	var player_name: String = GameManager.get_display_name() if GameManager.has_method("get_display_name") else "玩家"
	var guild_name := player_name + "的公會"
	GameManager.send_message("create_guild", {
		"name": guild_name,
		"description": "一起討伐！"
	})

## 加入公會按鈕（請求公會列表後選擇）
func _on_join_guild_pressed() -> void:
	GameManager.send_message("get_guild_list", {})
	# 顯示公會列表（簡化版：直接顯示第一個公會）
	# 完整版應該顯示選擇對話框

## 退出公會按鈕
func _on_leave_guild_pressed() -> void:
	if _guild_id == "":
		return
	GameManager.send_message("leave_guild", {})

## 建立聊天室區域（DAY-075）
func _build_chat_area(parent: Node, y_offset: int) -> void:
	_chat_container = Node2D.new()
	_chat_container.name = "ChatContainer"
	parent.add_child(_chat_container)

	# 聊天標題
	var chat_title := Label.new()
	chat_title.position = Vector2(8, y_offset)
	chat_title.text = "💬 公會聊天"
	chat_title.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	if _pixel_font:
		chat_title.add_theme_font_override("font", _pixel_font)
		chat_title.add_theme_font_size_override("font_size", 10)
	_chat_container.add_child(chat_title)

	# 聊天訊息區域背景
	var chat_bg := ColorRect.new()
	chat_bg.position = Vector2(4, y_offset + 14)
	chat_bg.size = Vector2(PANEL_WIDTH - 8, 80)
	chat_bg.color = Color(0.02, 0.02, 0.05, 0.8)
	_chat_container.add_child(chat_bg)

	# 聊天訊息顯示節點
	_chat_messages_node = Node2D.new()
	_chat_messages_node.name = "ChatMessages"
	_chat_messages_node.position = Vector2(4, y_offset + 14)
	_chat_container.add_child(_chat_messages_node)

	# 輸入框
	_chat_input = LineEdit.new()
	_chat_input.position = Vector2(4, y_offset + 98)
	_chat_input.size = Vector2(PANEL_WIDTH - 50, 18)
	_chat_input.placeholder_text = "輸入訊息..."
	_chat_input.max_length = 100
	if _pixel_font:
		_chat_input.add_theme_font_override("font", _pixel_font)
		_chat_input.add_theme_font_size_override("font_size", 9)
	_chat_input.text_submitted.connect(_on_chat_submitted)
	_chat_container.add_child(_chat_input)

	# 發送按鈕
	var send_btn := Button.new()
	send_btn.text = "發送"
	send_btn.position = Vector2(PANEL_WIDTH - 44, y_offset + 98)
	send_btn.size = Vector2(40, 18)
	if _pixel_font:
		send_btn.add_theme_font_override("font", _pixel_font)
		send_btn.add_theme_font_size_override("font_size", 9)
	send_btn.pressed.connect(func(): _on_chat_submitted(_chat_input.text))
	_chat_container.add_child(send_btn)

## 刷新聊天訊息顯示
func _refresh_chat_messages() -> void:
	if not _chat_messages_node:
		return

	# 清除舊訊息
	for child in _chat_messages_node.get_children():
		child.queue_free()

	# 顯示最近 5 條訊息
	var start_idx := max(0, _chat_history.size() - 5)
	for i in range(start_idx, _chat_history.size()):
		var chat_data: Dictionary = _chat_history[i]
		var role_icon := _get_role_icon(chat_data.get("role", "member"))
		var name_str: String = chat_data.get("display_name", "???")
		var msg_str: String = chat_data.get("message", "")
		var y_pos := (i - start_idx) * 15 + 3

		var msg_label := Label.new()
		msg_label.position = Vector2(4, y_pos)
		msg_label.text = role_icon + name_str + ": " + msg_str
		msg_label.add_theme_color_override("font_color", Color(0.9, 0.9, 0.9))
		if _pixel_font:
			msg_label.add_theme_font_override("font", _pixel_font)
			msg_label.add_theme_font_size_override("font_size", 9)
		_chat_messages_node.add_child(msg_label)

## 發送聊天訊息
func _on_chat_submitted(text: String) -> void:
	if text.strip_edges() == "":
		return
	if _guild_id == "":
		return
	GameManager.send_message("guild_chat", {"message": text.strip_edges()})
	if _chat_input:
		_chat_input.text = ""
