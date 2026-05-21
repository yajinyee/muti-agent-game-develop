## WinStreakPanel.gd — 連勝獎勵面板（DAY-131）
## 業界依據：BGaming Fishing Club 2026 Best Win/Best Catch 里程碑
## 顯示連勝次數、到下一個里程碑的進度、里程碑達成動畫
extends Control

# ---- 節點引用 ----
var _status_panel: PanelContainer
var _streak_label: Label          # 連勝次數
var _progress_bar_bg: ColorRect
var _progress_bar_fill: ColorRect
var _next_milestone_label: Label  # 下一個里程碑
var _timer_label: Label           # 超時倒數
var _milestone_popup: PanelContainer # 里程碑達成彈窗
var _milestone_icon_label: Label
var _milestone_name_label: Label
var _milestone_reward_label: Label
var _broadcast_banner: PanelContainer # 全服廣播橫幅
var _broadcast_label: Label

# ---- 狀態 ----
var _current_streak: int = 0
var _max_streak: int = 0
var _progress: float = 0.0
var _seconds_to_expiry: float = 0.0
var _popup_tween: Tween
var _banner_tween: Tween

const POPUP_DURATION := 3.0
const BANNER_DURATION := 4.0

func _ready() -> void:
	_build_ui()
	_status_panel.visible = false
	_milestone_popup.visible = false
	_broadcast_banner.visible = false

func _build_ui() -> void:
	# 左下角狀態面板（在龍怒面板上方）
	_status_panel = PanelContainer.new()
	_status_panel.position = Vector2(8, 720 - 200)
	_status_panel.custom_minimum_size = Vector2(130, 80)
	add_child(_status_panel)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.05, 0.1, 0.85)
	style.border_color = Color(0.5, 0.5, 1.0, 0.8)
	style.set_border_width_all(1)
	style.set_corner_radius_all(6)
	_status_panel.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 3)
	_status_panel.add_child(vbox)

	# 連勝次數
	_streak_label = Label.new()
	_streak_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_streak_label.add_theme_font_size_override("font_size", 22)
	_streak_label.add_theme_color_override("font_color", Color.WHITE)
	_streak_label.text = "🔥 0"
	vbox.add_child(_streak_label)

	# 下一個里程碑
	_next_milestone_label = Label.new()
	_next_milestone_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_next_milestone_label.add_theme_font_size_override("font_size", 10)
	_next_milestone_label.add_theme_color_override("font_color", Color(0.8, 0.8, 1.0))
	_next_milestone_label.text = "→ 🥉 10"
	vbox.add_child(_next_milestone_label)

	# 進度條
	var bar_container = Control.new()
	bar_container.custom_minimum_size = Vector2(110, 8)
	vbox.add_child(bar_container)

	_progress_bar_bg = ColorRect.new()
	_progress_bar_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	_progress_bar_bg.color = Color(0.2, 0.2, 0.3, 0.8)
	bar_container.add_child(_progress_bar_bg)

	_progress_bar_fill = ColorRect.new()
	_progress_bar_fill.anchor_top = 0.0
	_progress_bar_fill.anchor_bottom = 1.0
	_progress_bar_fill.anchor_left = 0.0
	_progress_bar_fill.anchor_right = 0.0
	_progress_bar_fill.color = Color(0.5, 0.5, 1.0)
	bar_container.add_child(_progress_bar_fill)

	# 超時倒數
	_timer_label = Label.new()
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_timer_label.add_theme_font_size_override("font_size", 9)
	_timer_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.8))
	vbox.add_child(_timer_label)

	# 里程碑達成彈窗（中央）
	_milestone_popup = PanelContainer.new()
	_milestone_popup.set_anchors_preset(Control.PRESET_CENTER)
	_milestone_popup.offset_left = -120
	_milestone_popup.offset_right = 120
	_milestone_popup.offset_top = -60
	_milestone_popup.offset_bottom = 60
	add_child(_milestone_popup)

	var popup_style = StyleBoxFlat.new()
	popup_style.bg_color = Color(0.05, 0.05, 0.15, 0.95)
	popup_style.border_color = Color.GOLD
	popup_style.set_border_width_all(3)
	popup_style.set_corner_radius_all(10)
	_milestone_popup.add_theme_stylebox_override("panel", popup_style)

	var popup_vbox = VBoxContainer.new()
	popup_vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	popup_vbox.add_theme_constant_override("separation", 4)
	_milestone_popup.add_child(popup_vbox)

	_milestone_icon_label = Label.new()
	_milestone_icon_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_milestone_icon_label.add_theme_font_size_override("font_size", 36)
	popup_vbox.add_child(_milestone_icon_label)

	_milestone_name_label = Label.new()
	_milestone_name_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_milestone_name_label.add_theme_font_size_override("font_size", 16)
	_milestone_name_label.add_theme_color_override("font_color", Color.GOLD)
	popup_vbox.add_child(_milestone_name_label)

	_milestone_reward_label = Label.new()
	_milestone_reward_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_milestone_reward_label.add_theme_font_size_override("font_size", 14)
	_milestone_reward_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	popup_vbox.add_child(_milestone_reward_label)

	# 全服廣播橫幅（頂部）
	_broadcast_banner = PanelContainer.new()
	_broadcast_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_broadcast_banner.custom_minimum_size = Vector2(0, 52)
	_broadcast_banner.position = Vector2(0, -56)
	add_child(_broadcast_banner)

	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.05, 0.05, 0.15, 0.92)
	banner_style.border_color = Color.GOLD
	banner_style.set_border_width_all(2)
	banner_style.set_corner_radius_all(6)
	_broadcast_banner.add_theme_stylebox_override("panel", banner_style)

	_broadcast_label = Label.new()
	_broadcast_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_broadcast_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_broadcast_label.add_theme_font_size_override("font_size", 16)
	_broadcast_label.add_theme_color_override("font_color", Color.GOLD)
	_broadcast_banner.add_child(_broadcast_label)

func _process(delta: float) -> void:
	if _seconds_to_expiry > 0:
		_seconds_to_expiry -= delta
		if _seconds_to_expiry < 0:
			_seconds_to_expiry = 0
		_update_timer()

# ---- 外部呼叫 ----

func on_win_streak_update(data: Dictionary) -> void:
	_current_streak = data.get("current", 0)
	_max_streak = data.get("max_streak", 0)
	_progress = data.get("progress_to_next", 0.0)
	_seconds_to_expiry = data.get("seconds_to_expiry", 0.0)

	var next_milestone = data.get("next_milestone", 10)
	var next_name = data.get("next_milestone_name", "銅牌連勝")

	# 更新顯示
	_streak_label.text = "🔥 %d" % _current_streak
	_next_milestone_label.text = "→ %d (%s)" % [next_milestone, next_name]
	_progress_bar_fill.anchor_right = _progress

	# 顯示面板（第一次擊破後才顯示）
	if _current_streak > 0 and not _status_panel.visible:
		_status_panel.visible = true
		_status_panel.modulate.a = 0.0
		var tween = create_tween()
		tween.tween_property(_status_panel, "modulate:a", 1.0, 0.3)

	# 進度條顏色
	if _progress >= 0.8:
		_progress_bar_fill.color = Color.GOLD
	elif _progress >= 0.5:
		_progress_bar_fill.color = Color(0.7, 0.7, 1.0)
	else:
		_progress_bar_fill.color = Color(0.5, 0.5, 1.0)

func on_win_streak_milestone(data: Dictionary) -> void:
	var icon = data.get("icon", "🏆")
	var name = data.get("level_name", "里程碑")
	var reward = data.get("bonus_reward", 0)
	var color_str = data.get("color", "#FFD700")
	var is_broadcast = data.get("broadcast", false)
	var player_name = data.get("player_name", "")
	var streak = data.get("streak", 0)

	# 個人彈窗
	_show_milestone_popup(icon, name, reward, Color(color_str))

	# 全服廣播橫幅
	if is_broadcast:
		_show_broadcast_banner("%s %s 達成 %s！連勝 %d 次！" % [icon, player_name, name, streak])

func on_win_streak_reset(data: Dictionary) -> void:
	var final_streak = data.get("final_streak", 0)
	_current_streak = 0
	_progress = 0.0
	_seconds_to_expiry = 0.0

	if final_streak > 0:
		_streak_label.text = "🔥 0"
		_progress_bar_fill.anchor_right = 0.0
		# 淡出面板
		var tween = create_tween()
		tween.tween_property(_status_panel, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func(): _status_panel.visible = false)

# ---- 私有方法 ----

func _update_timer() -> void:
	if not is_instance_valid(_timer_label):
		return
	if _seconds_to_expiry > 0:
		_timer_label.text = "%.0fs" % _seconds_to_expiry
		if _seconds_to_expiry <= 5:
			_timer_label.add_theme_color_override("font_color", Color.RED)
		else:
			_timer_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.8))
	else:
		_timer_label.text = ""

func _show_milestone_popup(icon: String, name: String, reward: int, color: Color) -> void:
	_milestone_icon_label.text = icon
	_milestone_name_label.text = name
	_milestone_name_label.add_theme_color_override("font_color", color)
	_milestone_reward_label.text = "+%d 金幣" % reward

	_milestone_popup.visible = true
	_milestone_popup.modulate.a = 0.0
	_milestone_popup.scale = Vector2(0.5, 0.5)

	if _popup_tween:
		_popup_tween.kill()
	_popup_tween = create_tween()
	_popup_tween.tween_property(_milestone_popup, "modulate:a", 1.0, 0.2)
	_popup_tween.parallel().tween_property(_milestone_popup, "scale", Vector2(1.0, 1.0), 0.2).set_ease(Tween.EASE_OUT)
	_popup_tween.tween_interval(POPUP_DURATION)
	_popup_tween.tween_property(_milestone_popup, "modulate:a", 0.0, 0.3)
	_popup_tween.tween_callback(func(): _milestone_popup.visible = false)

func _show_broadcast_banner(text: String) -> void:
	_broadcast_label.text = text
	_broadcast_banner.visible = true

	if _banner_tween:
		_banner_tween.kill()
	_banner_tween = create_tween()
	_broadcast_banner.position.y = -56
	_banner_tween.tween_property(_broadcast_banner, "position:y", 0.0, 0.3).set_ease(Tween.EASE_OUT)
	_banner_tween.tween_interval(BANNER_DURATION)
	_banner_tween.tween_property(_broadcast_banner, "position:y", -56.0, 0.3).set_ease(Tween.EASE_IN)
	_banner_tween.tween_callback(func(): _broadcast_banner.visible = false)
