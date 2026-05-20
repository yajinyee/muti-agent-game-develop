## TitlePanel.gd — 稱號顯示面板（DAY-068）
## 顯示玩家當前稱號，解鎖新稱號時彈出通知
extends Node2D

# ── 常數 ──────────────────────────────────────────────────────────────────────
const PANEL_WIDTH  := 200
const PANEL_HEIGHT := 28
const NOTIFY_DURATION := 3.5  # 稱號解鎖通知顯示秒數

# ── 節點引用 ──────────────────────────────────────────────────────────────────
var _font: FontFile
var _bg: ColorRect
var _title_label: Label
var _notify_container: Control  # 稱號解鎖通知容器

# ── 狀態 ──────────────────────────────────────────────────────────────────────
var _current_title_id: String = "novice"
var _current_title_name: String = "新手討伐者"
var _current_title_icon: String = "🌱"
var _current_title_color: String = "#AAAAAA"
var _notify_queue: Array = []
var _is_showing_notify: bool = false

# ── 初始化 ────────────────────────────────────────────────────────────────────
func setup(font: FontFile) -> void:
	_font = font
	_build_ui()
	_connect_signals()

func _build_ui() -> void:
	# 背景
	_bg = ColorRect.new()
	_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_bg.color = Color(0.05, 0.05, 0.15, 0.75)
	add_child(_bg)

	# 稱號標籤
	_title_label = Label.new()
	_title_label.position = Vector2(6, 4)
	_title_label.size = Vector2(PANEL_WIDTH - 12, PANEL_HEIGHT - 8)
	_title_label.text = _current_title_icon + " " + _current_title_name
	_title_label.add_theme_color_override("font_color", Color.html(_current_title_color))
	if _font:
		_title_label.add_theme_font_override("font", _font)
		_title_label.add_theme_font_size_override("font_size", 13)
	add_child(_title_label)

	# 通知容器（CanvasLayer 上）
	_notify_container = Control.new()
	_notify_container.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(_notify_container)

func _connect_signals() -> void:
	if GameManager.has_signal("title_unlocked"):
		GameManager.title_unlocked.connect(_on_title_unlocked)
	if GameManager.has_signal("player_updated"):
		GameManager.player_updated.connect(_on_player_updated)

# ── 更新稱號顯示 ──────────────────────────────────────────────────────────────
func update_title(title_id: String, title_name: String, title_icon: String, title_color: String) -> void:
	_current_title_id = title_id
	_current_title_name = title_name
	_current_title_icon = title_icon
	_current_title_color = title_color
	if _title_label:
		_title_label.text = title_icon + " " + title_name
		_title_label.add_theme_color_override("font_color", Color.html(title_color))

# ── 稱號解鎖通知 ──────────────────────────────────────────────────────────────
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
	var title_name: String = data.get("title_name", "新稱號")
	var title_icon: String = data.get("title_icon", "🏆")
	var title_color: String = data.get("title_color", "#FFD700")
	var description: String = data.get("description", "")

	# 通知容器（固定在畫面右下角）
	var notify := Control.new()
	notify.set_anchors_and_offsets_preset(Control.PRESET_BOTTOM_RIGHT)
	notify.offset_left = -280
	notify.offset_top = -90
	notify.offset_right = -10
	notify.offset_bottom = -10
	_notify_container.add_child(notify)

	# 背景
	var bg := ColorRect.new()
	bg.size = Vector2(270, 80)
	bg.color = Color(0.05, 0.05, 0.15, 0.92)
	notify.add_child(bg)

	# 左側彩色邊條
	var bar := ColorRect.new()
	bar.size = Vector2(4, 80)
	bar.color = Color.html(title_color)
	notify.add_child(bar)

	# 標題行：「🏆 稱號解鎖！」
	var header := Label.new()
	header.position = Vector2(12, 6)
	header.text = "🏆 稱號解鎖！"
	header.add_theme_color_override("font_color", Color.html(title_color))
	if _font:
		header.add_theme_font_override("font", _font)
		header.add_theme_font_size_override("font_size", 12)
	notify.add_child(header)

	# 稱號名稱
	var name_label := Label.new()
	name_label.position = Vector2(12, 26)
	name_label.text = title_icon + " " + title_name
	name_label.add_theme_color_override("font_color", Color.WHITE)
	if _font:
		name_label.add_theme_font_override("font", _font)
		name_label.add_theme_font_size_override("font_size", 15)
	notify.add_child(name_label)

	# 描述
	if description != "":
		var desc_label := Label.new()
		desc_label.position = Vector2(12, 52)
		desc_label.text = description
		desc_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
		if _font:
			desc_label.add_theme_font_override("font", _font)
			desc_label.add_theme_font_size_override("font_size", 11)
		notify.add_child(desc_label)

	# 滑入動畫
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

# ── 訊號處理 ──────────────────────────────────────────────────────────────────
func _on_player_updated(snapshot: Dictionary) -> void:
	var title_id: String = snapshot.get("title_id", "novice")
	var title_name: String = snapshot.get("title_name", "新手討伐者")
	var title_icon: String = snapshot.get("title_icon", "🌱")
	var title_color: String = snapshot.get("title_color", "#AAAAAA")
	update_title(title_id, title_name, title_icon, title_color)
