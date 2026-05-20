## LoginMilestonePanel.gd — 登入里程碑達成面板（DAY-107）
## 業界依據：ilogos.biz（2026）確認 gamified login streaks 讓留存率提升 75%
## 功能：
##   - 里程碑達成時顯示慶祝動畫（縮放彈入 + 金色光暈 + 獎勵列表）
##   - 登入進度面板（顯示所有里程碑進度條）
##   - 連續登入天數顯示（距離下一個里程碑還差幾天）
extends CanvasLayer

# ---- 節點引用 ----
var _overlay: ColorRect
var _panel: PanelContainer
var _title_label: Label
var _streak_label: Label
var _milestone_name: Label
var _milestone_icon: Label
var _rewards_container: VBoxContainer
var _close_btn: Button
var _progress_panel: PanelContainer
var _progress_container: VBoxContainer
var _glow_rect: ColorRect

# ---- 狀態 ----
var _is_showing := false
var _queue: Array = []  # 排隊等待顯示的里程碑

func _ready() -> void:
	layer = 95  # 在 HUD 之上，在 PlayerCard 之下
	_build_ui()
	visible = false

func _build_ui() -> void:
	# 半透明遮罩
	_overlay = ColorRect.new()
	_overlay.color = Color(0, 0, 0, 0.7)
	_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(_overlay)
	_overlay.gui_input.connect(_on_overlay_input)

	# 主面板（置中）
	_panel = PanelContainer.new()
	_panel.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	_panel.custom_minimum_size = Vector2(480, 520)
	_panel.offset_left = -240
	_panel.offset_top = -260
	_panel.offset_right = 240
	_panel.offset_bottom = 260
	add_child(_panel)

	# 金色光暈背景
	_glow_rect = ColorRect.new()
	_glow_rect.color = Color(1.0, 0.85, 0.0, 0.0)
	_glow_rect.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_panel.add_child(_glow_rect)

	var vbox := VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 12)
	_panel.add_child(vbox)

	# 標題
	_title_label = Label.new()
	_title_label.text = "🎉 里程碑達成！"
	_title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_title_label.add_theme_font_size_override("font_size", 22)
	_title_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	vbox.add_child(_title_label)

	# 連續天數
	_streak_label = Label.new()
	_streak_label.text = "連續登入 X 天"
	_streak_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_streak_label.add_theme_font_size_override("font_size", 16)
	_streak_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	vbox.add_child(_streak_label)

	# 里程碑圖示
	_milestone_icon = Label.new()
	_milestone_icon.text = "👑"
	_milestone_icon.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_milestone_icon.add_theme_font_size_override("font_size", 64)
	vbox.add_child(_milestone_icon)

	# 里程碑名稱
	_milestone_name = Label.new()
	_milestone_name.text = "月度傳說"
	_milestone_name.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_milestone_name.add_theme_font_size_override("font_size", 28)
	_milestone_name.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	vbox.add_child(_milestone_name)

	# 分隔線
	var sep := HSeparator.new()
	vbox.add_child(sep)

	# 獎勵標題
	var reward_title := Label.new()
	reward_title.text = "🎁 獎勵"
	reward_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_title.add_theme_font_size_override("font_size", 16)
	reward_title.add_theme_color_override("font_color", Color(0.9, 0.9, 0.9))
	vbox.add_child(reward_title)

	# 獎勵列表
	_rewards_container = VBoxContainer.new()
	_rewards_container.add_theme_constant_override("separation", 6)
	vbox.add_child(_rewards_container)

	# 關閉按鈕
	_close_btn = Button.new()
	_close_btn.text = "✓ 太棒了！"
	_close_btn.custom_minimum_size = Vector2(200, 44)
	_close_btn.add_theme_font_size_override("font_size", 16)
	_close_btn.pressed.connect(_on_close_pressed)
	var btn_container := HBoxContainer.new()
	btn_container.alignment = BoxContainer.ALIGNMENT_CENTER
	btn_container.add_child(_close_btn)
	vbox.add_child(btn_container)

## show_milestone 顯示里程碑達成面板
func show_milestone(data: Dictionary) -> void:
	if _is_showing:
		_queue.append(data)
		return
	_display_milestone(data)

## _display_milestone 實際顯示里程碑
func _display_milestone(data: Dictionary) -> void:
	_is_showing = true

	# 填入資料
	var days: int = data.get("days", 0)
	var name: String = data.get("name", "")
	var icon: String = data.get("icon", "🏆")
	var color_hex: String = data.get("color", "#FFD700")
	var rewards: Array = data.get("rewards", [])
	var coins_gained: int = data.get("coins_gained", 0)

	_streak_label.text = "連續登入 %d 天" % days
	_milestone_icon.text = icon
	_milestone_name.text = name

	# 設定顏色
	var color := Color.html(color_hex) if color_hex.begins_with("#") else Color.GOLD
	_milestone_name.add_theme_color_override("font_color", color)

	# 清空並填入獎勵
	for child in _rewards_container.get_children():
		child.queue_free()

	for reward in rewards:
		var row := HBoxContainer.new()
		row.alignment = BoxContainer.ALIGNMENT_CENTER
		row.add_theme_constant_override("separation", 8)

		var reward_label := Label.new()
		reward_label.add_theme_font_size_override("font_size", 15)

		var rtype: String = reward.get("type", "")
		var amount: int = reward.get("amount", 0)
		var rarity: String = reward.get("rarity", "")
		var title_id: String = reward.get("title_id", "")

		match rtype:
			"coins":
				reward_label.text = "🪙 +%s 金幣" % _format_number(amount)
				reward_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
			"mystery_box":
				var rarity_icon := _get_rarity_icon(rarity)
				reward_label.text = "%s %s 神秘寶箱 ×%d" % [rarity_icon, _get_rarity_name(rarity), amount]
				reward_label.add_theme_color_override("font_color", _get_rarity_color(rarity))
			"title":
				reward_label.text = "🏷️ 解鎖稱號"
				reward_label.add_theme_color_override("font_color", Color(0.8, 0.5, 1.0))

		row.add_child(reward_label)
		_rewards_container.add_child(row)

	# 顯示並播放動畫
	visible = true
	_panel.scale = Vector2(0.3, 0.3)
	_panel.modulate.a = 0.0

	var tween := create_tween()
	tween.set_parallel(true)
	tween.tween_property(_panel, "scale", Vector2(1.0, 1.0), 0.35).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.tween_property(_panel, "modulate:a", 1.0, 0.25)

	# 金色光暈閃爍
	var glow_tween := create_tween().set_loops(3)
	glow_tween.tween_property(_glow_rect, "color:a", 0.15, 0.3)
	glow_tween.tween_property(_glow_rect, "color:a", 0.0, 0.3)

## _on_close_pressed 關閉面板
func _on_close_pressed() -> void:
	_close_panel()

## _on_overlay_input 點擊遮罩關閉
func _on_overlay_input(event: InputEvent) -> void:
	if event is InputEventMouseButton and event.pressed:
		_close_panel()

## _close_panel 關閉並顯示下一個排隊的里程碑
func _close_panel() -> void:
	var tween := create_tween()
	tween.tween_property(_panel, "modulate:a", 0.0, 0.2)
	tween.tween_callback(func():
		visible = false
		_is_showing = false
		if _queue.size() > 0:
			var next = _queue.pop_front()
			_display_milestone(next)
	)

# ---- 進度面板 ----

## show_progress 顯示登入進度面板
func show_progress(data: Dictionary) -> void:
	# 建立進度面板（如果還沒建立）
	if _progress_panel == null:
		_build_progress_panel()
	_fill_progress_data(data)
	_progress_panel.visible = true

## _build_progress_panel 建立進度面板
func _build_progress_panel() -> void:
	_progress_panel = PanelContainer.new()
	_progress_panel.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	_progress_panel.custom_minimum_size = Vector2(520, 480)
	_progress_panel.offset_left = -260
	_progress_panel.offset_top = -240
	_progress_panel.offset_right = 260
	_progress_panel.offset_bottom = 240
	add_child(_progress_panel)

	var vbox := VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 10)
	_progress_panel.add_child(vbox)

	var header := Label.new()
	header.text = "📅 登入連續記錄"
	header.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	header.add_theme_font_size_override("font_size", 20)
	header.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	vbox.add_child(header)

	_progress_container = VBoxContainer.new()
	_progress_container.add_theme_constant_override("separation", 8)
	var scroll := ScrollContainer.new()
	scroll.custom_minimum_size = Vector2(0, 320)
	scroll.add_child(_progress_container)
	vbox.add_child(scroll)

	var close_btn := Button.new()
	close_btn.text = "✕ 關閉"
	close_btn.custom_minimum_size = Vector2(160, 40)
	close_btn.pressed.connect(func(): _progress_panel.visible = false)
	var btn_row := HBoxContainer.new()
	btn_row.alignment = BoxContainer.ALIGNMENT_CENTER
	btn_row.add_child(close_btn)
	vbox.add_child(btn_row)

## _fill_progress_data 填入進度資料
func _fill_progress_data(data: Dictionary) -> void:
	for child in _progress_container.get_children():
		child.queue_free()

	var current_streak: int = data.get("current_streak", 0)
	var max_streak: int = data.get("max_streak", 0)
	var next_days: int = data.get("next_milestone_days", 0)
	var days_to_next: int = data.get("days_to_next", 0)
	var milestones: Array = data.get("milestones", [])

	# 當前連續天數
	var streak_row := HBoxContainer.new()
	streak_row.alignment = BoxContainer.ALIGNMENT_CENTER
	var streak_lbl := Label.new()
	streak_lbl.text = "🔥 當前連續：%d 天（最高 %d 天）" % [current_streak, max_streak]
	streak_lbl.add_theme_font_size_override("font_size", 16)
	streak_lbl.add_theme_color_override("font_color", Color(1.0, 0.7, 0.2))
	streak_row.add_child(streak_lbl)
	_progress_container.add_child(streak_row)

	# 下一個里程碑提示
	if next_days > 0:
		var next_row := HBoxContainer.new()
		next_row.alignment = BoxContainer.ALIGNMENT_CENTER
		var next_lbl := Label.new()
		next_lbl.text = "⏳ 距離下一個里程碑還差 %d 天" % days_to_next
		next_lbl.add_theme_font_size_override("font_size", 14)
		next_lbl.add_theme_color_override("font_color", Color(0.7, 0.9, 1.0))
		next_row.add_child(next_lbl)
		_progress_container.add_child(next_row)

	_progress_container.add_child(HSeparator.new())

	# 里程碑列表
	for m in milestones:
		var m_days: int = m.get("days", 0)
		var m_name: String = m.get("name", "")
		var m_icon: String = m.get("icon", "🏆")
		var m_color_hex: String = m.get("color", "#FFFFFF")
		var is_reached: bool = m.get("is_reached", false)

		var row := HBoxContainer.new()
		row.add_theme_constant_override("separation", 10)

		# 狀態圖示
		var status_lbl := Label.new()
		status_lbl.text = "✅" if is_reached else "⬜"
		status_lbl.add_theme_font_size_override("font_size", 18)
		row.add_child(status_lbl)

		# 里程碑圖示
		var icon_lbl := Label.new()
		icon_lbl.text = m_icon
		icon_lbl.add_theme_font_size_override("font_size", 20)
		row.add_child(icon_lbl)

		# 名稱和天數
		var info_vbox := VBoxContainer.new()
		info_vbox.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		var name_lbl := Label.new()
		name_lbl.text = m_name
		name_lbl.add_theme_font_size_override("font_size", 15)
		var m_color := Color.html(m_color_hex) if m_color_hex.begins_with("#") else Color.WHITE
		if not is_reached:
			m_color = m_color.darkened(0.4)
		name_lbl.add_theme_color_override("font_color", m_color)
		info_vbox.add_child(name_lbl)
		var days_lbl := Label.new()
		days_lbl.text = "第 %d 天" % m_days
		days_lbl.add_theme_font_size_override("font_size", 12)
		days_lbl.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
		info_vbox.add_child(days_lbl)
		row.add_child(info_vbox)

		# 進度條（未達成時顯示）
		if not is_reached:
			var progress := ProgressBar.new()
			progress.custom_minimum_size = Vector2(80, 16)
			progress.min_value = 0
			progress.max_value = m_days
			progress.value = current_streak
			progress.show_percentage = false
			row.add_child(progress)

		_progress_container.add_child(row)

# ---- 工具函數 ----

func _format_number(n: int) -> String:
	if n >= 1000:
		return "%dk" % (n / 1000)
	return str(n)

func _get_rarity_icon(rarity: String) -> String:
	match rarity:
		"common": return "📦"
		"rare": return "💎"
		"epic": return "🔮"
		"legendary": return "👑"
	return "📦"

func _get_rarity_name(rarity: String) -> String:
	match rarity:
		"common": return "普通"
		"rare": return "稀有"
		"epic": return "史詩"
		"legendary": return "傳說"
	return "普通"

func _get_rarity_color(rarity: String) -> Color:
	match rarity:
		"common": return Color(0.7, 0.7, 0.7)
		"rare": return Color(0.3, 0.5, 1.0)
		"epic": return Color(0.7, 0.3, 1.0)
		"legendary": return Color(1.0, 0.85, 0.0)
	return Color.WHITE
