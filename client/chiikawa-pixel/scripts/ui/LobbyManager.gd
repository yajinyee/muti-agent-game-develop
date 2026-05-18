## LobbyManager.gd
## 遊戲大廳 UI 管理（DAY-020）
## 顯示房間列表，讓玩家選擇房間後進入遊戲
## 由 Main.tscn 在遊戲啟動時顯示

extends Control

# 像素字體
var _pixel_font: Font = null
const PIXEL_FONT_PATH = "res://assets/fonts/pixel8.fnt"

# UI 節點
var _room_list_container: Control = null
var _status_label: Label = null
var _refresh_btn: Button = null
var _quick_join_btn: Button = null
var _title_label: Label = null

# 房間資料
var _rooms: Array = []
var _selected_room_id: String = ""

# 訊號
signal room_selected(room_id: String)

func _ready() -> void:
	# 載入像素字體
	if ResourceLoader.exists(PIXEL_FONT_PATH):
		_pixel_font = load(PIXEL_FONT_PATH)

	_build_ui()
	NetworkManager.rooms_fetched.connect(_on_rooms_fetched)

	# 啟動時自動查詢房間列表
	_refresh_rooms()

## 建立大廳 UI
func _build_ui() -> void:
	# 全螢幕背景
	set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)

	# 海底漸層背景
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.02, 0.05, 0.15, 1.0)
	add_child(bg)

	# 標題區域
	var title_panel = Control.new()
	title_panel.position = Vector2(0, 40)
	title_panel.size = Vector2(1280, 100)
	add_child(title_panel)

	_title_label = Label.new()
	_title_label.text = "吉伊卡哇：像素大討伐"
	_title_label.position = Vector2(0, 0)
	_title_label.size = Vector2(1280, 60)
	_title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_title_label.add_theme_font_size_override("font_size", 42)
	_title_label.modulate = Color(1.0, 0.9, 0.3)
	if is_instance_valid(_pixel_font):
		_title_label.add_theme_font_override("font", _pixel_font)
	title_panel.add_child(_title_label)

	var subtitle = Label.new()
	subtitle.text = "選擇房間進入遊戲"
	subtitle.position = Vector2(0, 62)
	subtitle.size = Vector2(1280, 30)
	subtitle.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	subtitle.add_theme_font_size_override("font_size", 16)
	subtitle.modulate = Color(0.7, 0.85, 1.0, 0.8)
	if is_instance_valid(_pixel_font):
		subtitle.add_theme_font_override("font", _pixel_font)
	title_panel.add_child(subtitle)

	# 金色分隔線
	var divider = ColorRect.new()
	divider.position = Vector2(240, 148)
	divider.size = Vector2(800, 2)
	divider.color = Color(0.9, 0.75, 0.2, 0.6)
	add_child(divider)

	# 房間列表容器（中央）
	var list_bg = ColorRect.new()
	list_bg.position = Vector2(240, 160)
	list_bg.size = Vector2(800, 420)
	list_bg.color = Color(0.03, 0.07, 0.2, 0.85)
	add_child(list_bg)

	# 列表標題列
	var header_bg = ColorRect.new()
	header_bg.position = Vector2(240, 160)
	header_bg.size = Vector2(800, 36)
	header_bg.color = Color(0.05, 0.15, 0.4, 0.95)
	add_child(header_bg)

	var header_room = Label.new()
	header_room.text = "房間名稱"
	header_room.position = Vector2(260, 166)
	header_room.add_theme_font_size_override("font_size", 13)
	header_room.modulate = Color(1.0, 0.9, 0.3)
	if is_instance_valid(_pixel_font):
		header_room.add_theme_font_override("font", _pixel_font)
	add_child(header_room)

	var header_players = Label.new()
	header_players.text = "玩家"
	header_players.position = Vector2(620, 166)
	header_players.add_theme_font_size_override("font_size", 13)
	header_players.modulate = Color(1.0, 0.9, 0.3)
	if is_instance_valid(_pixel_font):
		header_players.add_theme_font_override("font", _pixel_font)
	add_child(header_players)

	var header_bet = Label.new()
	header_bet.text = "投注等級"
	header_bet.position = Vector2(720, 166)
	header_bet.add_theme_font_size_override("font_size", 13)
	header_bet.modulate = Color(1.0, 0.9, 0.3)
	if is_instance_valid(_pixel_font):
		header_bet.add_theme_font_override("font", _pixel_font)
	add_child(header_bet)

	# 房間列表滾動容器
	_room_list_container = Control.new()
	_room_list_container.name = "RoomListContainer"
	_room_list_container.position = Vector2(240, 196)
	_room_list_container.size = Vector2(800, 384)
	add_child(_room_list_container)

	# 狀態文字（載入中/錯誤）
	_status_label = Label.new()
	_status_label.name = "StatusLabel"
	_status_label.text = "載入房間列表中..."
	_status_label.position = Vector2(240, 340)
	_status_label.size = Vector2(800, 40)
	_status_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_status_label.add_theme_font_size_override("font_size", 16)
	_status_label.modulate = Color(0.7, 0.7, 0.7)
	if is_instance_valid(_pixel_font):
		_status_label.add_theme_font_override("font", _pixel_font)
	add_child(_status_label)

	# 底部按鈕區
	_refresh_btn = _make_button("🔄 重新整理", Vector2(340, 610), Vector2(200, 44))
	_refresh_btn.pressed.connect(_refresh_rooms)
	add_child(_refresh_btn)

	_quick_join_btn = _make_button("⚡ 快速加入", Vector2(560, 610), Vector2(200, 44))
	_quick_join_btn.pressed.connect(_quick_join)
	_quick_join_btn.modulate = Color(0.3, 1.0, 0.5)
	add_child(_quick_join_btn)

	# 版本資訊
	var ver_lbl = Label.new()
	ver_lbl.text = "v1.0  |  Go + Godot 4  |  Port 7777"
	ver_lbl.position = Vector2(0, 700)
	ver_lbl.size = Vector2(1280, 20)
	ver_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	ver_lbl.add_theme_font_size_override("font_size", 11)
	ver_lbl.modulate = Color(0.4, 0.5, 0.6, 0.6)
	if is_instance_valid(_pixel_font):
		ver_lbl.add_theme_font_override("font", _pixel_font)
	add_child(ver_lbl)

## 建立按鈕
func _make_button(text: String, pos: Vector2, size: Vector2) -> Button:
	var btn = Button.new()
	btn.text = text
	btn.position = pos
	btn.size = size
	btn.add_theme_font_size_override("font_size", 16)
	if is_instance_valid(_pixel_font):
		btn.add_theme_font_override("font", _pixel_font)
	return btn

## 查詢房間列表
func _refresh_rooms() -> void:
	_status_label.text = "載入房間列表中..."
	_status_label.visible = true
	_clear_room_list()
	NetworkManager.fetch_rooms()

## 清除房間列表
func _clear_room_list() -> void:
	if not is_instance_valid(_room_list_container):
		return
	for child in _room_list_container.get_children():
		child.queue_free()

## 收到房間列表
func _on_rooms_fetched(rooms: Array) -> void:
	_rooms = rooms
	_clear_room_list()

	if rooms.is_empty():
		_status_label.text = "無法連線到 Server，請確認 Server 已啟動"
		_status_label.modulate = Color(1.0, 0.4, 0.4)
		_status_label.visible = true
		return

	_status_label.visible = false

	for i in range(rooms.size()):
		_create_room_row(rooms[i], i)

## 建立房間列表行
func _create_room_row(room_data: Dictionary, index: int) -> void:
	var row_height = 72
	var row = Control.new()
	row.name = "RoomRow_%s" % room_data.get("id", str(index))
	row.position = Vector2(0, index * row_height)
	row.size = Vector2(800, row_height - 4)
	_room_list_container.add_child(row)

	var is_full = room_data.get("is_full", false)
	var player_count = room_data.get("player_count", 0)
	var max_players = room_data.get("max_players", 16)
	var room_id = room_data.get("id", "room-001")

	# 行背景（交替色）
	var row_bg = ColorRect.new()
	row_bg.size = Vector2(800, row_height - 4)
	if index % 2 == 0:
		row_bg.color = Color(0.04, 0.09, 0.25, 0.7)
	else:
		row_bg.color = Color(0.03, 0.07, 0.20, 0.7)
	row.add_child(row_bg)

	# 滿員時加紅色遮罩
	if is_full:
		var full_overlay = ColorRect.new()
		full_overlay.size = Vector2(800, row_height - 4)
		full_overlay.color = Color(0.3, 0.0, 0.0, 0.3)
		row.add_child(full_overlay)

	# 房間名稱
	var name_lbl = Label.new()
	name_lbl.text = room_data.get("name", room_id)
	name_lbl.position = Vector2(16, 8)
	name_lbl.add_theme_font_size_override("font_size", 18)
	name_lbl.modulate = Color(0.5, 0.5, 0.5) if is_full else Color.WHITE
	if is_instance_valid(_pixel_font):
		name_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(name_lbl)

	# 房間 ID（小字）
	var id_lbl = Label.new()
	id_lbl.text = room_id
	id_lbl.position = Vector2(16, 36)
	id_lbl.add_theme_font_size_override("font_size", 11)
	id_lbl.modulate = Color(0.5, 0.6, 0.8, 0.7)
	if is_instance_valid(_pixel_font):
		id_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(id_lbl)

	# 玩家數量（帶進度條）
	var player_lbl = Label.new()
	player_lbl.text = "%d/%d" % [player_count, max_players]
	player_lbl.position = Vector2(380, 12)
	player_lbl.add_theme_font_size_override("font_size", 16)
	var fill_ratio = float(player_count) / float(max_players)
	if is_full:
		player_lbl.modulate = Color(1.0, 0.3, 0.3)
	elif fill_ratio > 0.7:
		player_lbl.modulate = Color(1.0, 0.8, 0.2)
	else:
		player_lbl.modulate = Color(0.4, 1.0, 0.5)
	if is_instance_valid(_pixel_font):
		player_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(player_lbl)

	# 玩家數量進度條
	var bar_bg = ColorRect.new()
	bar_bg.position = Vector2(380, 38)
	bar_bg.size = Vector2(100, 8)
	bar_bg.color = Color(0.1, 0.1, 0.2, 0.8)
	row.add_child(bar_bg)

	var bar_fill = ColorRect.new()
	bar_fill.position = Vector2(380, 38)
	bar_fill.size = Vector2(100.0 * fill_ratio, 8)
	bar_fill.color = Color(1.0, 0.3, 0.3) if is_full else Color(0.3, 0.9, 0.4)
	row.add_child(bar_fill)

	# 投注等級範圍
	var bet_lbl = Label.new()
	var min_bet = room_data.get("min_bet_level", 1)
	var max_bet = room_data.get("max_bet_level", 10)
	bet_lbl.text = "LV%d - LV%d" % [min_bet, max_bet]
	bet_lbl.position = Vector2(510, 12)
	bet_lbl.add_theme_font_size_override("font_size", 14)
	bet_lbl.modulate = Color(0.8, 0.85, 1.0)
	if is_instance_valid(_pixel_font):
		bet_lbl.add_theme_font_override("font", _pixel_font)
	row.add_child(bet_lbl)

	# 加入按鈕（滿員時禁用）
	var join_btn = Button.new()
	join_btn.text = "滿員" if is_full else "加入"
	join_btn.position = Vector2(680, 16)
	join_btn.size = Vector2(100, 36)
	join_btn.add_theme_font_size_override("font_size", 14)
	join_btn.disabled = is_full
	if not is_full:
		join_btn.modulate = Color(0.3, 1.0, 0.5)
		join_btn.pressed.connect(func(): _join_room(room_id))
	else:
		join_btn.modulate = Color(0.5, 0.5, 0.5)
	if is_instance_valid(_pixel_font):
		join_btn.add_theme_font_override("font", _pixel_font)
	row.add_child(join_btn)

## 加入指定房間
func _join_room(room_id: String) -> void:
	_selected_room_id = room_id
	print("[Lobby] Joining room: ", room_id)
	# 切換到指定房間（重新連線）
	NetworkManager.connect_to_room(room_id)
	emit_signal("room_selected", room_id)
	# 淡出大廳 UI
	var tween = create_tween()
	tween.tween_property(self, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func(): visible = false)

## 快速加入（人數最少的房間）
func _quick_join() -> void:
	if _rooms.is_empty():
		_refresh_rooms()
		return

	# 找人數最少且未滿的房間
	var best: Dictionary = {}
	for r in _rooms:
		if r.get("is_full", false):
			continue
		if best.is_empty() or r.get("player_count", 0) < best.get("player_count", 999):
			best = r

	if best.is_empty():
		_status_label.text = "所有房間已滿，請稍後再試"
		_status_label.modulate = Color(1.0, 0.5, 0.2)
		_status_label.visible = true
		return

	_join_room(best.get("id", "room-001"))

## 顯示大廳（從遊戲中呼叫）
func show_lobby() -> void:
	visible = true
	modulate.a = 0.0
	var tween = create_tween()
	tween.tween_property(self, "modulate:a", 1.0, 0.3)
	_refresh_rooms()
