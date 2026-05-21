# RaidPanel.gd — Co-op Boss Raid 面板（DAY-115）
# 顯示全服合作討伐狀態：BOSS HP、倒數計時、貢獻排名
extends CanvasLayer

# ---- 節點引用 ----
var _root: Control
var _warning_panel: Control
var _active_panel: Control
var _result_panel: Control

var _warning_label: Label
var _warning_countdown: Label

var _boss_name_label: Label
var _hp_bar: ProgressBar
var _hp_label: Label
var _time_label: Label
var _reward_pool_label: Label
var _rank_container: VBoxContainer

var _result_title: Label
var _result_reward_label: Label
var _result_rank_container: VBoxContainer
var _result_close_btn: Button

# ---- 狀態 ----
var _current_state: String = "idle"
var _my_player_id: String = ""
var _warning_timer: float = 30.0
var _time_left: float = 0.0

func _ready() -> void:
	layer = 70
	_build_ui()
	hide()

func _build_ui() -> void:
	_root = Control.new()
	_root.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(_root)

	# 半透明背景
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0, 0, 0, 0.6)
	_root.add_child(bg)

	# 主容器（置中）
	var center = CenterContainer.new()
	center.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_root.add_child(center)

	var main_box = VBoxContainer.new()
	main_box.custom_minimum_size = Vector2(600, 500)
	center.add_child(main_box)

	# ---- 警告面板 ----
	_warning_panel = _build_warning_panel()
	main_box.add_child(_warning_panel)

	# ---- 討伐進行面板 ----
	_active_panel = _build_active_panel()
	main_box.add_child(_active_panel)
	_active_panel.hide()

	# ---- 結算面板 ----
	_result_panel = _build_result_panel()
	main_box.add_child(_result_panel)
	_result_panel.hide()

func _build_warning_panel() -> Control:
	var panel = PanelContainer.new()
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.6, 0.1, 0.1, 0.95)
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	panel.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 16)
	panel.add_child(vbox)

	var title = Label.new()
	title.text = "⚔️ CO-OP BOSS RAID"
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_font_size_override("font_size", 28)
	title.add_theme_color_override("font_color", Color(1, 0.9, 0.2))
	vbox.add_child(title)

	_warning_label = Label.new()
	_warning_label.text = "全服合作討伐「吉伊卡哇大魔王」即將開始！"
	_warning_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_warning_label.add_theme_font_size_override("font_size", 16)
	_warning_label.add_theme_color_override("font_color", Color.WHITE)
	_warning_label.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	vbox.add_child(_warning_label)

	var reward_info = Label.new()
	reward_info.text = "🏆 獎勵池：200,000 金幣（依貢獻度分配）"
	reward_info.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_info.add_theme_font_size_override("font_size", 14)
	reward_info.add_theme_color_override("font_color", Color(1, 0.85, 0.3))
	vbox.add_child(reward_info)

	_warning_countdown = Label.new()
	_warning_countdown.text = "30 秒後開始"
	_warning_countdown.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_warning_countdown.add_theme_font_size_override("font_size", 22)
	_warning_countdown.add_theme_color_override("font_color", Color(1, 0.5, 0.5))
	vbox.add_child(_warning_countdown)

	return panel

func _build_active_panel() -> Control:
	var panel = PanelContainer.new()
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.05, 0.2, 0.95)
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	panel.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 10)
	panel.add_child(vbox)

	# BOSS 名稱
	_boss_name_label = Label.new()
	_boss_name_label.text = "⚔️ 吉伊卡哇大魔王"
	_boss_name_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_boss_name_label.add_theme_font_size_override("font_size", 22)
	_boss_name_label.add_theme_color_override("font_color", Color(1, 0.4, 0.4))
	vbox.add_child(_boss_name_label)

	# HP 條
	_hp_bar = ProgressBar.new()
	_hp_bar.custom_minimum_size = Vector2(560, 28)
	_hp_bar.max_value = 50000
	_hp_bar.value = 50000
	_hp_bar.show_percentage = false
	vbox.add_child(_hp_bar)

	_hp_label = Label.new()
	_hp_label.text = "HP: 50000 / 50000"
	_hp_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_hp_label.add_theme_font_size_override("font_size", 13)
	_hp_label.add_theme_color_override("font_color", Color(0.9, 0.9, 0.9))
	vbox.add_child(_hp_label)

	# 時間 + 獎勵池
	var info_row = HBoxContainer.new()
	info_row.alignment = BoxContainer.ALIGNMENT_CENTER
	info_row.add_theme_constant_override("separation", 30)
	vbox.add_child(info_row)

	_time_label = Label.new()
	_time_label.text = "⏱ 5:00"
	_time_label.add_theme_font_size_override("font_size", 18)
	_time_label.add_theme_color_override("font_color", Color(0.5, 1, 0.5))
	info_row.add_child(_time_label)

	_reward_pool_label = Label.new()
	_reward_pool_label.text = "🏆 200,000 金幣"
	_reward_pool_label.add_theme_font_size_override("font_size", 18)
	_reward_pool_label.add_theme_color_override("font_color", Color(1, 0.85, 0.3))
	info_row.add_child(_reward_pool_label)

	# 貢獻排名標題
	var rank_title = Label.new()
	rank_title.text = "📊 貢獻排名"
	rank_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	rank_title.add_theme_font_size_override("font_size", 14)
	rank_title.add_theme_color_override("font_color", Color(0.8, 0.8, 1))
	vbox.add_child(rank_title)

	# 貢獻排名列表
	var scroll = ScrollContainer.new()
	scroll.custom_minimum_size = Vector2(560, 180)
	vbox.add_child(scroll)

	_rank_container = VBoxContainer.new()
	_rank_container.add_theme_constant_override("separation", 4)
	scroll.add_child(_rank_container)

	return panel

func _build_result_panel() -> Control:
	var panel = PanelContainer.new()
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.15, 0.05, 0.95)
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	panel.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 12)
	panel.add_child(vbox)

	_result_title = Label.new()
	_result_title.text = "🏆 RAID 勝利！"
	_result_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_title.add_theme_font_size_override("font_size", 26)
	_result_title.add_theme_color_override("font_color", Color(1, 0.9, 0.2))
	vbox.add_child(_result_title)

	_result_reward_label = Label.new()
	_result_reward_label.text = "你的獎勵：0 金幣（排名 #0）"
	_result_reward_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_reward_label.add_theme_font_size_override("font_size", 18)
	_result_reward_label.add_theme_color_override("font_color", Color(1, 0.85, 0.3))
	vbox.add_child(_result_reward_label)

	var rank_title = Label.new()
	rank_title.text = "📊 最終貢獻排名"
	rank_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	rank_title.add_theme_font_size_override("font_size", 14)
	rank_title.add_theme_color_override("font_color", Color(0.8, 0.8, 1))
	vbox.add_child(rank_title)

	var scroll = ScrollContainer.new()
	scroll.custom_minimum_size = Vector2(560, 200)
	vbox.add_child(scroll)

	_result_rank_container = VBoxContainer.new()
	_result_rank_container.add_theme_constant_override("separation", 4)
	scroll.add_child(_result_rank_container)

	_result_close_btn = Button.new()
	_result_close_btn.text = "關閉"
	_result_close_btn.custom_minimum_size = Vector2(120, 36)
	_result_close_btn.pressed.connect(_on_close_pressed)
	vbox.add_child(_result_close_btn)

	return panel

func _process(delta: float) -> void:
	if not visible:
		return

	# 警告倒數
	if _current_state == "warning":
		_warning_timer -= delta
		if _warning_timer > 0:
			_warning_countdown.text = "%.0f 秒後開始" % _warning_timer
		else:
			_warning_countdown.text = "開始！"

	# 討伐倒數（本地計時，Server 每 3 秒同步一次）
	if _current_state == "active" and _time_left > 0:
		_time_left -= delta
		var mins = int(_time_left) / 60
		var secs = int(_time_left) % 60
		_time_label.text = "⏱ %d:%02d" % [mins, secs]
		# 最後 30 秒變紅色
		if _time_left <= 30:
			_time_label.add_theme_color_override("font_color", Color(1, 0.3, 0.3))
		else:
			_time_label.add_theme_color_override("font_color", Color(0.5, 1, 0.5))

# ---- 外部呼叫 API ----

func show_warning(data: Dictionary) -> void:
	_current_state = "warning"
	_warning_timer = float(data.get("starts_in", 30))
	_warning_label.text = "全服合作討伐「%s」即將開始！" % data.get("boss_name", "大魔王")
	_warning_panel.show()
	_active_panel.hide()
	_result_panel.hide()
	show()
	# 閃爍動畫
	_animate_warning()

func show_raid_start(data: Dictionary) -> void:
	_current_state = "active"
	_time_left = float(data.get("duration", 300))
	_boss_name_label.text = "⚔️ " + data.get("boss_name", "大魔王")
	_hp_bar.max_value = data.get("max_hp", 50000)
	_hp_bar.value = data.get("hp", 50000)
	_hp_label.text = "HP: %d / %d" % [data.get("hp", 50000), data.get("max_hp", 50000)]
	_reward_pool_label.text = "🏆 %s 金幣" % _format_coins(data.get("reward_pool", 200000))
	_warning_panel.hide()
	_active_panel.show()
	_result_panel.hide()
	show()

func update_raid(data: Dictionary) -> void:
	if _current_state != "active":
		return
	var hp = data.get("hp", 0)
	var max_hp = data.get("max_hp", 50000)
	_hp_bar.value = hp
	_hp_label.text = "HP: %d / %d" % [hp, max_hp]
	_time_left = float(data.get("time_left", _time_left))
	_update_rank_list(_rank_container, data.get("contributors", []), false)

func show_result(data: Dictionary, my_player_id: String) -> void:
	_current_state = "result"
	var defeated = data.get("defeated", false)
	if defeated:
		_result_title.text = "🏆 RAID 勝利！"
		_result_title.add_theme_color_override("font_color", Color(1, 0.9, 0.2))
	else:
		_result_title.text = "💀 RAID 失敗..."
		_result_title.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))

	# 找到自己的獎勵
	var my_reward = data.get("my_reward", 0)
	var my_rank = data.get("my_rank", 0)
	if my_reward > 0:
		_result_reward_label.text = "你的獎勵：%s 金幣（排名 #%d）" % [_format_coins(my_reward), my_rank]
		_result_reward_label.add_theme_color_override("font_color", Color(1, 0.85, 0.3))
	else:
		_result_reward_label.text = "本次未參與討伐"
		_result_reward_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))

	_update_rank_list(_result_rank_container, data.get("contributors", []), true)
	_warning_panel.hide()
	_active_panel.hide()
	_result_panel.show()
	show()

	# 10 秒後自動關閉
	var timer = get_tree().create_timer(10.0)
	timer.timeout.connect(_on_close_pressed)

func _update_rank_list(container: VBoxContainer, contributors: Array, show_reward: bool) -> void:
	# 清空
	for child in container.get_children():
		child.queue_free()

	var rank_icons = ["🥇", "🥈", "🥉"]
	for i in range(min(contributors.size(), 10)):
		var entry = contributors[i]
		var row = HBoxContainer.new()
		row.add_theme_constant_override("separation", 8)
		container.add_child(row)

		# 排名圖示
		var rank_icon = Label.new()
		rank_icon.text = rank_icons[i] if i < 3 else "#%d" % (i + 1)
		rank_icon.custom_minimum_size = Vector2(36, 0)
		rank_icon.add_theme_font_size_override("font_size", 14)
		row.add_child(rank_icon)

		# 玩家名稱
		var name_label = Label.new()
		name_label.text = entry.get("display_name", "玩家")
		name_label.custom_minimum_size = Vector2(180, 0)
		name_label.add_theme_font_size_override("font_size", 13)
		name_label.add_theme_color_override("font_color", Color.WHITE)
		row.add_child(name_label)

		# 傷害
		var dmg_label = Label.new()
		dmg_label.text = "傷害: %s" % _format_coins(entry.get("damage", 0))
		dmg_label.custom_minimum_size = Vector2(120, 0)
		dmg_label.add_theme_font_size_override("font_size", 13)
		dmg_label.add_theme_color_override("font_color", Color(1, 0.7, 0.7))
		row.add_child(dmg_label)

		# 獎勵（結算後才顯示）
		if show_reward:
			var reward_label = Label.new()
			reward_label.text = "+%s 💰" % _format_coins(entry.get("reward", 0))
			reward_label.add_theme_font_size_override("font_size", 13)
			reward_label.add_theme_color_override("font_color", Color(1, 0.85, 0.3))
			row.add_child(reward_label)

func _animate_warning() -> void:
	# 閃爍警告效果
	var tween = create_tween().set_loops(6)
	tween.tween_property(_warning_panel, "modulate", Color(1, 0.5, 0.5), 0.3)
	tween.tween_property(_warning_panel, "modulate", Color.WHITE, 0.3)

func _on_close_pressed() -> void:
	_current_state = "idle"
	hide()

func _format_coins(n: int) -> String:
	if n >= 10000:
		return "%d萬" % (n / 10000)
	return str(n)
