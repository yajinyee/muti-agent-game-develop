## AchievementPanel.gd — 成就系統 UI（DAY-349）
## 顯示成就解鎖通知（右下角彈出）+ 成就列表面板
## achievement-agent 負責維護
extends CanvasLayer

# ── 節點引用 ──────────────────────────────────────────────────
var _list_panel: PanelContainer
var _list_container: VBoxContainer
var _title_label: Label
var _progress_label: Label
var _close_btn: Button
var _refresh_btn: Button

# 通知佇列（避免多個成就同時彈出）
var _notify_queue: Array = []
var _is_showing_notify: bool = false

# ── 常數 ──────────────────────────────────────────────────────
const PANEL_WIDTH = 520
const PANEL_HEIGHT = 580
const NOTIFY_WIDTH = 340
const NOTIFY_HEIGHT = 80

# 稀有度顏色
const RARITY_COLORS = {
	"common":    Color(0.8, 0.8, 0.8),
	"rare":      Color(0.3, 0.6, 1.0),
	"epic":      Color(0.7, 0.2, 0.9),
	"legendary": Color(1.0, 0.75, 0.0),
}

const RARITY_LABELS = {
	"common":    "普通",
	"rare":      "稀有",
	"epic":      "史詩",
	"legendary": "傳說",
}

func _ready() -> void:
	layer = 24  # 在賽季排行榜上方
	_build_list_ui()
	_connect_signals()
	visible = false

func _build_list_ui() -> void:
	# 半透明背景遮罩
	var overlay = ColorRect.new()
	overlay.color = Color(0, 0, 0, 0.65)
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)

	# 主面板
	_list_panel = PanelContainer.new()
	_list_panel.set_anchors_preset(Control.PRESET_CENTER)
	_list_panel.custom_minimum_size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_list_panel.position = Vector2(-PANEL_WIDTH / 2, -PANEL_HEIGHT / 2)
	add_child(_list_panel)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.05, 0.12, 0.97)
	style.border_color = Color(1.0, 0.75, 0.0, 1.0)
	style.border_width_left = 2
	style.border_width_right = 2
	style.border_width_top = 2
	style.border_width_bottom = 2
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	_list_panel.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 6)
	_list_panel.add_child(vbox)

	# ── 標題列 ──────────────────────────────────────────────────
	var header = HBoxContainer.new()
	vbox.add_child(header)

	_title_label = Label.new()
	_title_label.text = "🏅 成就系統"
	_title_label.add_theme_font_size_override("font_size", 20)
	_title_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	_title_label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	header.add_child(_title_label)

	_refresh_btn = Button.new()
	_refresh_btn.text = "🔄"
	_refresh_btn.custom_minimum_size = Vector2(36, 36)
	header.add_child(_refresh_btn)

	_close_btn = Button.new()
	_close_btn.text = "✕"
	_close_btn.custom_minimum_size = Vector2(36, 36)
	header.add_child(_close_btn)

	# 進度標籤
	_progress_label = Label.new()
	_progress_label.text = "已解鎖：0 / 25"
	_progress_label.add_theme_font_size_override("font_size", 14)
	_progress_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	vbox.add_child(_progress_label)

	# 分隔線
	var sep = HSeparator.new()
	vbox.add_child(sep)

	# 成就列表（可捲動）
	var scroll = ScrollContainer.new()
	scroll.size_flags_vertical = Control.SIZE_EXPAND_FILL
	scroll.custom_minimum_size = Vector2(0, PANEL_HEIGHT - 120)
	vbox.add_child(scroll)

	_list_container = VBoxContainer.new()
	_list_container.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	_list_container.add_theme_constant_override("separation", 4)
	scroll.add_child(_list_container)

func _connect_signals() -> void:
	if is_instance_valid(_close_btn):
		_close_btn.pressed.connect(_on_close)
	if is_instance_valid(_refresh_btn):
		_refresh_btn.pressed.connect(_on_refresh)
	# 連接 GameManager 訊號
	if GameManager.has_signal("achievement_unlock"):
		GameManager.achievement_unlock.connect(_on_achievement_unlock)
	if GameManager.has_signal("achievement_list"):
		GameManager.achievement_list.connect(_on_achievement_list)

func _on_close() -> void:
	visible = false

func _on_refresh() -> void:
	GameManager.request_achievement_list()

func show_panel() -> void:
	visible = true
	_on_refresh()

# ── 成就解鎖通知（右下角彈出）────────────────────────────────

func _on_achievement_unlock(data: Dictionary) -> void:
	_notify_queue.append(data)
	if not _is_showing_notify:
		_show_next_notify()

func _show_next_notify() -> void:
	if _notify_queue.is_empty():
		_is_showing_notify = false
		return
	_is_showing_notify = true
	var data = _notify_queue.pop_front()
	_spawn_notify_popup(data)

func _spawn_notify_popup(data: Dictionary) -> void:
	var icon = data.get("icon", "🏅")
	var name_text = data.get("name", "成就解鎖")
	var rarity = data.get("rarity", "common")
	var reward = data.get("reward", 0)

	var rarity_color = RARITY_COLORS.get(rarity, Color.WHITE)
	var rarity_label = RARITY_LABELS.get(rarity, "普通")

	# 建立通知面板（右下角）
	var notify = Control.new()
	notify.size = Vector2(NOTIFY_WIDTH, NOTIFY_HEIGHT)
	notify.position = Vector2(1280 - NOTIFY_WIDTH - 16, 720 - NOTIFY_HEIGHT - 16)
	notify.z_index = 200
	add_child(notify)

	# 背景
	var bg = ColorRect.new()
	bg.size = Vector2(NOTIFY_WIDTH, NOTIFY_HEIGHT)
	bg.color = Color(0.05, 0.05, 0.12, 0.95)
	notify.add_child(bg)

	# 稀有度邊框
	var border = ColorRect.new()
	border.size = Vector2(4, NOTIFY_HEIGHT)
	border.color = rarity_color
	notify.add_child(border)

	# 圖示
	var icon_lbl = Label.new()
	icon_lbl.text = icon
	icon_lbl.position = Vector2(12, 10)
	icon_lbl.add_theme_font_size_override("font_size", 32)
	notify.add_child(icon_lbl)

	# 標題
	var title_lbl = Label.new()
	title_lbl.text = "🏅 成就解鎖！"
	title_lbl.position = Vector2(56, 6)
	title_lbl.add_theme_font_size_override("font_size", 12)
	title_lbl.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	notify.add_child(title_lbl)

	# 成就名稱
	var name_lbl = Label.new()
	name_lbl.text = name_text
	name_lbl.position = Vector2(56, 22)
	name_lbl.add_theme_font_size_override("font_size", 16)
	name_lbl.add_theme_color_override("font_color", rarity_color)
	notify.add_child(name_lbl)

	# 稀有度 + 獎勵
	var sub_lbl = Label.new()
	sub_lbl.text = "[%s] +%d 金幣" % [rarity_label, reward]
	sub_lbl.position = Vector2(56, 44)
	sub_lbl.add_theme_font_size_override("font_size", 12)
	sub_lbl.add_theme_color_override("font_color", Color(0.8, 0.8, 0.5))
	notify.add_child(sub_lbl)

	# 進場動畫：從右側滑入
	notify.modulate.a = 0.0
	notify.position.x = 1280
	var tween = notify.create_tween()
	tween.tween_property(notify, "position:x", 1280 - NOTIFY_WIDTH - 16, 0.3).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(notify, "modulate:a", 1.0, 0.3)
	# 停留 2.5 秒後淡出
	tween.tween_interval(2.5)
	tween.tween_property(notify, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(notify):
			notify.queue_free()
		_show_next_notify()
	)

# ── 成就列表更新 ─────────────────────────────────────────────

func _on_achievement_list(data: Dictionary) -> void:
	var achievements = data.get("achievements", [])
	var total = data.get("total_count", 0)
	var unlocked_count = data.get("unlocked_count", 0)

	_progress_label.text = "已解鎖：%d / %d" % [unlocked_count, total]

	# 清空舊列表
	for child in _list_container.get_children():
		child.queue_free()

	# 依稀有度排序：legendary > epic > rare > common
	var rarity_order = {"legendary": 0, "epic": 1, "rare": 2, "common": 3}
	achievements.sort_custom(func(a, b):
		var ra = rarity_order.get(a.get("rarity", "common"), 3)
		var rb = rarity_order.get(b.get("rarity", "common"), 3)
		if ra != rb:
			return ra < rb
		# 已解鎖的排前面
		return a.get("unlocked", false) and not b.get("unlocked", false)
	)

	for ach in achievements:
		_add_achievement_row(ach)

func _add_achievement_row(ach: Dictionary) -> void:
	var icon = ach.get("icon", "🏅")
	var name_text = ach.get("name", "")
	var desc = ach.get("description", "")
	var rarity = ach.get("rarity", "common")
	var reward = ach.get("reward", 0)
	var unlocked = ach.get("unlocked", false)

	var rarity_color = RARITY_COLORS.get(rarity, Color.WHITE)

	var row = HBoxContainer.new()
	row.custom_minimum_size = Vector2(0, 52)
	row.add_theme_constant_override("separation", 8)
	_list_container.add_child(row)

	# 背景
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	if unlocked:
		bg.color = Color(0.08, 0.12, 0.22, 0.9)
	else:
		bg.color = Color(0.05, 0.05, 0.08, 0.7)
	row.add_child(bg)

	# 稀有度條
	var rarity_bar = ColorRect.new()
	rarity_bar.custom_minimum_size = Vector2(4, 0)
	rarity_bar.size_flags_vertical = Control.SIZE_EXPAND_FILL
	if unlocked:
		rarity_bar.color = rarity_color
	else:
		rarity_bar.color = Color(rarity_color.r, rarity_color.g, rarity_color.b, 0.3)
	row.add_child(rarity_bar)

	# 圖示
	var icon_lbl = Label.new()
	icon_lbl.text = icon if unlocked else "🔒"
	icon_lbl.custom_minimum_size = Vector2(36, 0)
	icon_lbl.add_theme_font_size_override("font_size", 24)
	icon_lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	if not unlocked:
		icon_lbl.modulate = Color(0.4, 0.4, 0.4)
	row.add_child(icon_lbl)

	# 文字區
	var text_vbox = VBoxContainer.new()
	text_vbox.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	text_vbox.add_theme_constant_override("separation", 2)
	row.add_child(text_vbox)

	var name_lbl = Label.new()
	name_lbl.text = name_text
	name_lbl.add_theme_font_size_override("font_size", 14)
	if unlocked:
		name_lbl.add_theme_color_override("font_color", rarity_color)
	else:
		name_lbl.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))
	text_vbox.add_child(name_lbl)

	var desc_lbl = Label.new()
	desc_lbl.text = desc
	desc_lbl.add_theme_font_size_override("font_size", 11)
	desc_lbl.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	text_vbox.add_child(desc_lbl)

	# 獎勵
	var reward_lbl = Label.new()
	reward_lbl.text = "+%d 💰" % reward
	reward_lbl.custom_minimum_size = Vector2(60, 0)
	reward_lbl.add_theme_font_size_override("font_size", 12)
	reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2) if unlocked else Color(0.4, 0.4, 0.4))
	reward_lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	row.add_child(reward_lbl)
