## DailyQuestPanel.gd — 每日任務面板（DAY-345）
## 靈感來源：BGaming Quests（2026-05-27）
## 功能：顯示每日任務進度、完成通知、領取獎勵
extends CanvasLayer

# ── 節點引用 ──────────────────────────────────────────────────
var _panel: Panel
var _title_label: Label
var _quest_items: Array = []
var _coins_label: Label
var _reset_label: Label
var _close_btn: Button
var _toggle_btn: Button  # HUD 上的任務按鈕

# ── 狀態 ──────────────────────────────────────────────────────
var _is_visible: bool = false
var _quest_data: Array = []
var _quest_coins: int = 0
var _reset_at: int = 0

# ── 顏色常數 ──────────────────────────────────────────────────
const COLOR_BG = Color(0.08, 0.08, 0.15, 0.95)
const COLOR_PANEL = Color(0.12, 0.12, 0.22, 1.0)
const COLOR_QUEST_BG = Color(0.15, 0.15, 0.28, 1.0)
const COLOR_QUEST_DONE = Color(0.1, 0.25, 0.1, 1.0)
const COLOR_GOLD = Color(1.0, 0.85, 0.2, 1.0)
const COLOR_GREEN = Color(0.3, 0.9, 0.3, 1.0)
const COLOR_GRAY = Color(0.6, 0.6, 0.6, 1.0)
const COLOR_WHITE = Color(1.0, 1.0, 1.0, 1.0)
const COLOR_BLUE = Color(0.4, 0.7, 1.0, 1.0)

func _ready() -> void:
	layer = 20  # 在大部分 UI 上方
	_build_ui()
	_connect_signals()
	# 連線後自動請求任務狀態
	GameManager.daily_quest_update.connect(_on_quest_update)
	GameManager.daily_quest_complete.connect(_on_quest_complete)
	# 延遲請求，確保連線完成
	get_tree().create_timer(1.5).timeout.connect(func(): GameManager.request_daily_quests())

func _build_ui() -> void:
	# 主面板（右側滑出）
	_panel = Panel.new()
	_panel.set_anchors_preset(Control.PRESET_RIGHT_WIDE)
	_panel.size = Vector2(280, 400)
	_panel.position = Vector2(get_viewport().size.x - 290, 80)
	var style = StyleBoxFlat.new()
	style.bg_color = COLOR_PANEL
	style.corner_radius_top_left = 12
	style.corner_radius_bottom_left = 12
	style.border_color = COLOR_GOLD
	style.border_width_left = 2
	style.border_width_top = 2
	style.border_width_bottom = 2
	_panel.add_theme_stylebox_override("panel", style)
	add_child(_panel)

	# 標題
	_title_label = Label.new()
	_title_label.text = "🎯 每日任務"
	_title_label.position = Vector2(16, 12)
	_title_label.size = Vector2(200, 28)
	_title_label.add_theme_color_override("font_color", COLOR_GOLD)
	_title_label.add_theme_font_size_override("font_size", 18)
	_panel.add_child(_title_label)

	# 關閉按鈕
	_close_btn = Button.new()
	_close_btn.text = "✕"
	_close_btn.position = Vector2(248, 8)
	_close_btn.size = Vector2(28, 28)
	_close_btn.add_theme_color_override("font_color", COLOR_GRAY)
	_close_btn.pressed.connect(_toggle_panel)
	_panel.add_child(_close_btn)

	# 任務幣顯示
	_coins_label = Label.new()
	_coins_label.text = "🪙 任務幣：0"
	_coins_label.position = Vector2(16, 44)
	_coins_label.size = Vector2(248, 22)
	_coins_label.add_theme_color_override("font_color", COLOR_GOLD)
	_coins_label.add_theme_font_size_override("font_size", 14)
	_panel.add_child(_coins_label)

	# 重置時間
	_reset_label = Label.new()
	_reset_label.text = "重置：--:--:--"
	_reset_label.position = Vector2(16, 66)
	_reset_label.size = Vector2(248, 18)
	_reset_label.add_theme_color_override("font_color", COLOR_GRAY)
	_reset_label.add_theme_font_size_override("font_size", 11)
	_panel.add_child(_reset_label)

	# 分隔線
	var sep = ColorRect.new()
	sep.position = Vector2(8, 88)
	sep.size = Vector2(264, 1)
	sep.color = COLOR_GOLD * Color(1, 1, 1, 0.3)
	_panel.add_child(sep)

	# 任務項目（3個）
	for i in range(3):
		var item = _create_quest_item(i)
		_quest_items.append(item)
		_panel.add_child(item)

	# 預設隱藏
	_panel.visible = false

func _create_quest_item(index: int) -> Control:
	var container = Control.new()
	container.position = Vector2(8, 96 + index * 96)
	container.size = Vector2(264, 88)

	# 背景
	var bg = ColorRect.new()
	bg.size = Vector2(264, 84)
	bg.color = COLOR_QUEST_BG
	bg.name = "bg"
	container.add_child(bg)

	# 任務名稱
	var name_label = Label.new()
	name_label.position = Vector2(8, 6)
	name_label.size = Vector2(200, 20)
	name_label.add_theme_font_size_override("font_size", 13)
	name_label.add_theme_color_override("font_color", COLOR_WHITE)
	name_label.name = "name_label"
	container.add_child(name_label)

	# 獎勵標籤
	var reward_label = Label.new()
	reward_label.position = Vector2(200, 6)
	reward_label.size = Vector2(60, 20)
	reward_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	reward_label.add_theme_font_size_override("font_size", 12)
	reward_label.add_theme_color_override("font_color", COLOR_GOLD)
	reward_label.name = "reward_label"
	container.add_child(reward_label)

	# 描述
	var desc_label = Label.new()
	desc_label.position = Vector2(8, 26)
	desc_label.size = Vector2(248, 18)
	desc_label.add_theme_font_size_override("font_size", 11)
	desc_label.add_theme_color_override("font_color", COLOR_GRAY)
	desc_label.name = "desc_label"
	container.add_child(desc_label)

	# 進度條背景
	var bar_bg = ColorRect.new()
	bar_bg.position = Vector2(8, 48)
	bar_bg.size = Vector2(248, 10)
	bar_bg.color = Color(0.2, 0.2, 0.3, 1.0)
	bar_bg.name = "bar_bg"
	container.add_child(bar_bg)

	# 進度條
	var bar = ColorRect.new()
	bar.position = Vector2(8, 48)
	bar.size = Vector2(0, 10)
	bar.color = COLOR_BLUE
	bar.name = "bar"
	container.add_child(bar)

	# 進度文字
	var progress_label = Label.new()
	progress_label.position = Vector2(8, 60)
	progress_label.size = Vector2(180, 18)
	progress_label.add_theme_font_size_override("font_size", 11)
	progress_label.add_theme_color_override("font_color", COLOR_GRAY)
	progress_label.name = "progress_label"
	container.add_child(progress_label)

	# 領取按鈕
	var claim_btn = Button.new()
	claim_btn.position = Vector2(188, 58)
	claim_btn.size = Vector2(68, 22)
	claim_btn.text = "領取"
	claim_btn.visible = false
	claim_btn.name = "claim_btn"
	container.add_child(claim_btn)

	return container

func _connect_signals() -> void:
	pass

func _toggle_panel() -> void:
	_is_visible = !_is_visible
	_panel.visible = _is_visible
	if _is_visible:
		GameManager.request_daily_quests()
		# 滑入動畫
		var tween = create_tween()
		_panel.modulate.a = 0
		tween.tween_property(_panel, "modulate:a", 1.0, 0.2)

func show_panel() -> void:
	_is_visible = true
	_panel.visible = true
	GameManager.request_daily_quests()

func hide_panel() -> void:
	_is_visible = false
	_panel.visible = false

func _on_quest_update(data: Dictionary) -> void:
	_quest_data = data.get("quests", [])
	_quest_coins = data.get("quest_coins", 0)
	_reset_at = data.get("reset_at", 0)
	_refresh_ui()

func _refresh_ui() -> void:
	# 更新任務幣
	_coins_label.text = "🪙 任務幣：%d" % _quest_coins

	# 更新重置時間
	if _reset_at > 0:
		var now_ms = Time.get_unix_time_from_system() * 1000
		var diff_sec = int((_reset_at - now_ms) / 1000)
		if diff_sec > 0:
			var h = diff_sec / 3600
			var m = (diff_sec % 3600) / 60
			var s = diff_sec % 60
			_reset_label.text = "重置：%02d:%02d:%02d" % [h, m, s]
		else:
			_reset_label.text = "重置：即將重置"

	# 更新任務項目
	for i in range(min(_quest_data.size(), _quest_items.size())):
		var quest = _quest_data[i]
		var item = _quest_items[i]
		_update_quest_item(item, quest)

func _update_quest_item(item: Control, quest: Dictionary) -> void:
	var name_label = item.get_node("name_label")
	var desc_label = item.get_node("desc_label")
	var reward_label = item.get_node("reward_label")
	var bar = item.get_node("bar")
	var progress_label = item.get_node("progress_label")
	var claim_btn = item.get_node("claim_btn")
	var bg = item.get_node("bg")

	var quest_name = quest.get("name", "")
	var description = quest.get("description", "")
	var target = quest.get("target", 1)
	var progress = quest.get("progress", 0)
	var completed = quest.get("completed", false)
	var claimed = quest.get("claimed", false)
	var reward = quest.get("reward", 0)
	var quest_id = quest.get("id", "")

	name_label.text = quest_name
	desc_label.text = description
	reward_label.text = "🪙%d" % reward

	# 進度條
	var pct = clamp(float(progress) / float(max(target, 1)), 0.0, 1.0)
	bar.size.x = 248 * pct

	if completed and not claimed:
		# 可領取
		bg.color = COLOR_QUEST_DONE
		bar.color = COLOR_GREEN
		progress_label.text = "✅ 完成！"
		claim_btn.visible = true
		claim_btn.text = "領取"
		# 連接按鈕（避免重複連接）
		if not claim_btn.pressed.is_connected(func(): _claim_quest(quest_id)):
			# 清除舊連接
			for conn in claim_btn.pressed.get_connections():
				claim_btn.pressed.disconnect(conn["callable"])
			claim_btn.pressed.connect(func(): _claim_quest(quest_id))
	elif claimed:
		# 已領取
		bg.color = COLOR_QUEST_BG * Color(0.8, 0.8, 0.8, 1.0)
		bar.color = COLOR_GRAY
		progress_label.text = "✅ 已領取"
		claim_btn.visible = false
	else:
		# 進行中
		bg.color = COLOR_QUEST_BG
		bar.color = COLOR_BLUE
		progress_label.text = "%d / %d" % [progress, target]
		claim_btn.visible = false

func _claim_quest(quest_id: String) -> void:
	GameManager.claim_daily_quest(quest_id)
	# 播放音效
	if AudioManager:
		AudioManager.play_sfx(AudioManager.SFX.COIN_DROP)

func _on_quest_complete(data: Dictionary) -> void:
	var quest_name = data.get("quest_name", "任務")
	var reward = data.get("reward", 0)
	var message = data.get("message", "")
	_show_complete_notification(quest_name, reward)

func _show_complete_notification(quest_name: String, reward: int) -> void:
	# 建立通知浮窗
	var notif = Panel.new()
	notif.size = Vector2(280, 60)
	notif.position = Vector2(
		get_viewport().size.x / 2 - 140,
		get_viewport().size.y - 120
	)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.3, 0.1, 0.95)
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	style.border_color = COLOR_GREEN
	style.border_width_left = 2
	style.border_width_right = 2
	style.border_width_top = 2
	style.border_width_bottom = 2
	notif.add_theme_stylebox_override("panel", style)
	add_child(notif)

	var label = Label.new()
	label.text = "🎯 任務完成：%s\n點擊任務面板領取 🪙%d" % [quest_name, reward]
	label.position = Vector2(8, 8)
	label.size = Vector2(264, 44)
	label.add_theme_font_size_override("font_size", 12)
	label.add_theme_color_override("font_color", COLOR_GREEN)
	notif.add_child(label)

	# 動畫：上升 + 淡出
	var tween = create_tween()
	tween.tween_property(notif, "position:y", notif.position.y - 40, 2.0)
	tween.parallel().tween_property(notif, "modulate:a", 0.0, 2.0).set_delay(1.5)
	tween.tween_callback(notif.queue_free)

func _process(_delta: float) -> void:
	# 每秒更新重置倒數
	if _is_visible and _reset_at > 0:
		var now_ms = Time.get_unix_time_from_system() * 1000
		var diff_sec = int((_reset_at - now_ms) / 1000)
		if diff_sec > 0:
			var h = diff_sec / 3600
			var m = (diff_sec % 3600) / 60
			var s = diff_sec % 60
			_reset_label.text = "重置：%02d:%02d:%02d" % [h, m, s]
