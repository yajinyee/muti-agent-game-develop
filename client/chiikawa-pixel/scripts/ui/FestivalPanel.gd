## FestivalPanel.gd — 賽季節日活動面板（DAY-109）
## 顯示當前節日主題、限定目標物、節日任務進度、稱號獎勵
extends CanvasLayer

signal festival_task_claimed(task_id: String)

# 面板節點
var _panel: PanelContainer
var _title_label: Label
var _desc_label: Label
var _time_label: Label
var _bonus_label: Label
var _tasks_container: VBoxContainer
var _close_btn: Button
var _bg_overlay: ColorRect

# 當前節日資料
var _festival_data: Dictionary = {}

func _ready() -> void:
	layer = 88
	_build_ui()
	hide()

func _build_ui() -> void:
	# 半透明背景遮罩
	_bg_overlay = ColorRect.new()
	_bg_overlay.color = Color(0, 0, 0, 0.6)
	_bg_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(_bg_overlay)

	# 主面板
	_panel = PanelContainer.new()
	_panel.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	_panel.custom_minimum_size = Vector2(480, 560)
	_panel.position = Vector2(-240, -280)
	add_child(_panel)

	var vbox := VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 10)
	_panel.add_child(vbox)

	# 標題列
	var header := HBoxContainer.new()
	vbox.add_child(header)

	_title_label = Label.new()
	_title_label.text = "🎋 節日慶典"
	_title_label.add_theme_font_size_override("font_size", 20)
	_title_label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	header.add_child(_title_label)

	_close_btn = Button.new()
	_close_btn.text = "✕"
	_close_btn.custom_minimum_size = Vector2(32, 32)
	_close_btn.pressed.connect(_on_close)
	header.add_child(_close_btn)

	# 描述
	_desc_label = Label.new()
	_desc_label.text = ""
	_desc_label.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	_desc_label.add_theme_font_size_override("font_size", 13)
	vbox.add_child(_desc_label)

	# 加成資訊
	_bonus_label = Label.new()
	_bonus_label.text = ""
	_bonus_label.add_theme_font_size_override("font_size", 12)
	_bonus_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	vbox.add_child(_bonus_label)

	# 倒數計時
	_time_label = Label.new()
	_time_label.text = ""
	_time_label.add_theme_font_size_override("font_size", 12)
	_time_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	vbox.add_child(_time_label)

	# 分隔線
	var sep := HSeparator.new()
	vbox.add_child(sep)

	# 任務標題
	var task_title := Label.new()
	task_title.text = "📋 節日任務"
	task_title.add_theme_font_size_override("font_size", 15)
	vbox.add_child(task_title)

	# 任務列表（ScrollContainer）
	var scroll := ScrollContainer.new()
	scroll.custom_minimum_size = Vector2(0, 280)
	scroll.size_flags_vertical = Control.SIZE_EXPAND_FILL
	vbox.add_child(scroll)

	_tasks_container = VBoxContainer.new()
	_tasks_container.add_theme_constant_override("separation", 8)
	_tasks_container.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	scroll.add_child(_tasks_container)

func show_festival(data: Dictionary) -> void:
	_festival_data = data
	_refresh_ui()
	show()

func _refresh_ui() -> void:
	if _festival_data.is_empty():
		return

	var is_active: bool = _festival_data.get("is_active", false)
	if not is_active:
		_title_label.text = "目前沒有節日活動"
		_desc_label.text = "下次節日活動即將到來，敬請期待！"
		_bonus_label.text = ""
		_time_label.text = ""
		_clear_tasks()
		return

	var name: String = _festival_data.get("name", "節日慶典")
	var icon: String = _festival_data.get("icon", "🎉")
	var desc: String = _festival_data.get("description", "")
	var color_hex: String = _festival_data.get("color", "#FFFFFF")
	var jackpot_mult: float = _festival_data.get("jackpot_mult", 1.0)
	var reward_mult: float = _festival_data.get("reward_mult", 1.0)
	var bonus_chance: float = _festival_data.get("bonus_chance_add", 0.0)
	var time_left: float = _festival_data.get("time_left", 0.0)

	_title_label.text = "%s %s" % [icon, name]
	var title_color := Color.html(color_hex) if color_hex.begins_with("#") else Color.WHITE
	_title_label.add_theme_color_override("font_color", title_color)

	_desc_label.text = desc

	var bonus_parts: Array = []
	if reward_mult > 1.0:
		bonus_parts.append("獎勵 ×%.1f" % reward_mult)
	if jackpot_mult > 1.0:
		bonus_parts.append("Jackpot ×%.1f" % jackpot_mult)
	if bonus_chance > 0.0:
		bonus_parts.append("Bonus率 +%d%%" % int(bonus_chance * 100))
	_bonus_label.text = "✨ 加成：" + "、".join(bonus_parts) if bonus_parts.size() > 0 else ""

	# 倒數計時
	var days := int(time_left / 86400)
	var hours := int(fmod(time_left, 86400) / 3600)
	if days > 0:
		_time_label.text = "⏰ 剩餘 %d 天 %d 小時" % [days, hours]
	else:
		var mins := int(fmod(time_left, 3600) / 60)
		_time_label.text = "⏰ 剩餘 %d 小時 %d 分" % [hours, mins]

	# 更新任務列表
	_refresh_tasks()

func _refresh_tasks() -> void:
	_clear_tasks()
	var tasks: Array = _festival_data.get("tasks", [])
	var title_id: String = _festival_data.get("title_id", "")
	var title_name: String = _festival_data.get("title_name", "")
	var title_claimed: bool = _festival_data.get("title_claimed", false)

	for task in tasks:
		_add_task_row(task)

	# 稱號獎勵行
	if title_id != "":
		_add_title_row(title_name, title_claimed, tasks)

func _add_task_row(task: Dictionary) -> void:
	var task_id: String = task.get("id", "")
	var desc: String = task.get("description", "")
	var target: int = task.get("target", 1)
	var progress: int = task.get("progress", 0)
	var done: bool = task.get("done", false)
	var reward_coins: int = task.get("reward_coins", 0)

	var row := HBoxContainer.new()
	row.add_theme_constant_override("separation", 8)
	_tasks_container.add_child(row)

	# 進度條容器
	var left := VBoxContainer.new()
	left.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	row.add_child(left)

	var desc_label := Label.new()
	desc_label.text = desc
	desc_label.add_theme_font_size_override("font_size", 12)
	if done:
		desc_label.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))
	left.add_child(desc_label)

	var progress_bar := ProgressBar.new()
	progress_bar.min_value = 0
	progress_bar.max_value = target
	progress_bar.value = min(progress, target)
	progress_bar.custom_minimum_size = Vector2(0, 14)
	left.add_child(progress_bar)

	var prog_label := Label.new()
	prog_label.text = "%d / %d" % [min(progress, target), target]
	prog_label.add_theme_font_size_override("font_size", 10)
	left.add_child(prog_label)

	# 領取按鈕
	var btn := Button.new()
	btn.custom_minimum_size = Vector2(80, 36)
	if done:
		btn.text = "✓ 已領取"
		btn.disabled = true
		btn.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))
	elif progress >= target:
		btn.text = "🎁 %d" % reward_coins
		btn.pressed.connect(func(): _on_claim_task(task_id))
	else:
		btn.text = "🎁 %d" % reward_coins
		btn.disabled = true
	row.add_child(btn)

func _add_title_row(title_name: String, claimed: bool, tasks: Array) -> void:
	var all_done := true
	for task in tasks:
		if not task.get("done", false):
			all_done = false
			break

	var sep := HSeparator.new()
	_tasks_container.add_child(sep)

	var row := HBoxContainer.new()
	_tasks_container.add_child(row)

	var label := Label.new()
	label.text = "🏆 完成全部任務獲得稱號：%s" % title_name
	label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	label.add_theme_font_size_override("font_size", 13)
	row.add_child(label)

	var btn := Button.new()
	btn.custom_minimum_size = Vector2(80, 36)
	if claimed:
		btn.text = "✓ 已獲得"
		btn.disabled = true
	elif all_done:
		btn.text = "領取稱號"
		btn.pressed.connect(func(): _on_claim_title())
	else:
		btn.text = "未完成"
		btn.disabled = true
	row.add_child(btn)

func _clear_tasks() -> void:
	for child in _tasks_container.get_children():
		child.queue_free()

func _on_claim_task(task_id: String) -> void:
	festival_task_claimed.emit(task_id)
	if GameManager:
		GameManager.send_claim_festival_task(task_id)

func _on_claim_title() -> void:
	# 稱號領取由 Server 在所有任務完成後自動處理
	# 這裡只是提示玩家
	pass

func _on_close() -> void:
	hide()

func update_from_data(data: Dictionary) -> void:
	_festival_data = data
	if is_visible():
		_refresh_ui()
