## FlashChallengePanel.gd — 閃電挑戰面板（DAY-123）
## 業界依據：Infingame（2026-05-19）確認 Challenges 工具是 2026 年最熱門留存機制
## 限時 90 秒的特殊目標挑戰，全服可見，完成獎勵豐厚
extends Control

# ---- 常數 ----
const PANEL_WIDTH := 320
const PANEL_HEIGHT := 200
const SLIDE_DURATION := 0.4
const AUTO_HIDE_DELAY := 5.0  # 挑戰結束後 5 秒自動隱藏

# ---- 節點引用 ----
var _bg: ColorRect
var _title_label: Label
var _desc_label: Label
var _progress_bar: ColorRect
var _progress_fill: ColorRect
var _progress_label: Label
var _timer_label: Label
var _reward_label: Label
var _leaderboard_container: VBoxContainer
var _result_label: Label

# ---- 狀態 ----
var _is_visible := false
var _current_target := 0
var _current_progress := 0
var _time_left := 0
var _challenge_color := Color(1.0, 0.85, 0.0)  # 預設金色
var _auto_hide_timer: SceneTreeTimer = null
var _tick_timer: SceneTreeTimer = null

func _ready() -> void:
	_build_ui()
	visible = false
	set_process(false)

func _build_ui() -> void:
	custom_minimum_size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	
	# 背景
	_bg = ColorRect.new()
	_bg.color = Color(0.05, 0.05, 0.15, 0.92)
	_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	add_child(_bg)
	
	# 頂部閃電條
	var top_bar := ColorRect.new()
	top_bar.color = Color(1.0, 0.85, 0.0, 1.0)
	top_bar.size = Vector2(PANEL_WIDTH, 4)
	top_bar.position = Vector2(0, 0)
	add_child(top_bar)
	
	# 標題
	_title_label = Label.new()
	_title_label.position = Vector2(10, 8)
	_title_label.size = Vector2(PANEL_WIDTH - 20, 28)
	_title_label.add_theme_font_size_override("font_size", 16)
	_title_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	add_child(_title_label)
	
	# 描述
	_desc_label = Label.new()
	_desc_label.position = Vector2(10, 36)
	_desc_label.size = Vector2(PANEL_WIDTH - 20, 20)
	_desc_label.add_theme_font_size_override("font_size", 11)
	_desc_label.add_theme_color_override("font_color", Color(0.85, 0.85, 0.85))
	add_child(_desc_label)
	
	# 計時器
	_timer_label = Label.new()
	_timer_label.position = Vector2(PANEL_WIDTH - 70, 8)
	_timer_label.size = Vector2(60, 28)
	_timer_label.add_theme_font_size_override("font_size", 18)
	_timer_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.2))
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	add_child(_timer_label)
	
	# 進度條背景
	_progress_bar = ColorRect.new()
	_progress_bar.color = Color(0.2, 0.2, 0.3)
	_progress_bar.position = Vector2(10, 62)
	_progress_bar.size = Vector2(PANEL_WIDTH - 20, 16)
	add_child(_progress_bar)
	
	# 進度條填充
	_progress_fill = ColorRect.new()
	_progress_fill.color = Color(1.0, 0.85, 0.0)
	_progress_fill.position = Vector2(10, 62)
	_progress_fill.size = Vector2(0, 16)
	add_child(_progress_fill)
	
	# 進度文字
	_progress_label = Label.new()
	_progress_label.position = Vector2(10, 62)
	_progress_label.size = Vector2(PANEL_WIDTH - 20, 16)
	_progress_label.add_theme_font_size_override("font_size", 11)
	_progress_label.add_theme_color_override("font_color", Color.WHITE)
	_progress_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	add_child(_progress_label)
	
	# 獎勵標籤
	_reward_label = Label.new()
	_reward_label.position = Vector2(10, 82)
	_reward_label.size = Vector2(PANEL_WIDTH - 20, 18)
	_reward_label.add_theme_font_size_override("font_size", 11)
	_reward_label.add_theme_color_override("font_color", Color(0.9, 0.8, 0.3))
	add_child(_reward_label)
	
	# 排行榜容器
	_leaderboard_container = VBoxContainer.new()
	_leaderboard_container.position = Vector2(10, 104)
	_leaderboard_container.size = Vector2(PANEL_WIDTH - 20, 80)
	add_child(_leaderboard_container)
	
	# 結果標籤（挑戰結束時顯示）
	_result_label = Label.new()
	_result_label.position = Vector2(0, 70)
	_result_label.size = Vector2(PANEL_WIDTH, 60)
	_result_label.add_theme_font_size_override("font_size", 20)
	_result_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_label.visible = false
	add_child(_result_label)

# ---- 公開 API ----

## 閃電挑戰開始
func on_flash_challenge_start(data: Dictionary) -> void:
	_current_target = data.get("target", 10)
	_current_progress = 0
	_time_left = data.get("time_left", 90)
	
	var color_str: String = data.get("color", "#FFD700")
	_challenge_color = Color(color_str)
	
	_title_label.text = data.get("title", "⚡ 閃電挑戰")
	_desc_label.text = data.get("description", "")
	_reward_label.text = "完成獎勵：🪙 " + str(data.get("base_reward", 0) + data.get("bonus_reward", 0))
	_result_label.visible = false
	
	# 更新顏色
	_progress_fill.color = _challenge_color
	_title_label.add_theme_color_override("font_color", _challenge_color)
	
	_update_progress_bar()
	_update_leaderboard(data.get("top_players", []))
	_update_timer()
	
	_show_panel()
	_start_tick()

## 閃電挑戰進度更新
func on_flash_challenge_update(data: Dictionary) -> void:
	_time_left = data.get("time_left", _time_left)
	_update_leaderboard(data.get("top_players", []))
	_update_timer()
	
	# 如果是自己的進度
	var my_progress: int = data.get("progress", -1)
	if my_progress >= 0:
		_current_progress = my_progress
		_update_progress_bar()
		
		if data.get("completed", false):
			_show_completion_flash()

## 閃電挑戰結束
func on_flash_challenge_end(data: Dictionary) -> void:
	_stop_tick()
	
	var success: bool = data.get("success", false)
	var message: String = data.get("message", "挑戰結束")
	
	_result_label.text = ("🎉 " if success else "⏰ ") + message
	_result_label.visible = true
	_leaderboard_container.visible = false
	_timer_label.text = "結束"
	
	# 5 秒後自動隱藏
	if _auto_hide_timer:
		_auto_hide_timer = null
	_auto_hide_timer = get_tree().create_timer(AUTO_HIDE_DELAY)
	_auto_hide_timer.timeout.connect(_hide_panel)

## 閃電挑戰獎勵通知
func on_flash_challenge_reward(data: Dictionary) -> void:
	var reward: int = data.get("reward", 0)
	var completed: bool = data.get("completed", false)
	var message: String = data.get("message", "")
	
	if reward <= 0:
		return
	
	# 彈出獎勵通知
	_show_reward_popup(reward, completed, message)

# ---- 內部函數 ----

func _show_panel() -> void:
	if _is_visible:
		return
	_is_visible = true
	visible = true
	_leaderboard_container.visible = true
	
	# 從右側滑入
	var start_x := position.x + PANEL_WIDTH + 20
	var end_x := position.x
	position.x = start_x
	
	var tween := create_tween()
	tween.tween_property(self, "position:x", end_x, SLIDE_DURATION).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

func _hide_panel() -> void:
	if not _is_visible:
		return
	_is_visible = false
	
	var tween := create_tween()
	tween.tween_property(self, "position:x", position.x + PANEL_WIDTH + 20, SLIDE_DURATION).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
	tween.tween_callback(func(): visible = false)

func _update_progress_bar() -> void:
	if _current_target <= 0:
		return
	var ratio := float(_current_progress) / float(_current_target)
	ratio = clamp(ratio, 0.0, 1.0)
	var bar_width := (PANEL_WIDTH - 20) * ratio
	_progress_fill.size.x = bar_width
	_progress_label.text = str(_current_progress) + " / " + str(_current_target)

func _update_timer() -> void:
	if _time_left <= 0:
		_timer_label.text = "⏰ 0"
		_timer_label.add_theme_color_override("font_color", Color(1.0, 0.2, 0.2))
		return
	
	_timer_label.text = "⏱ " + str(_time_left)
	
	# 最後 10 秒變紅閃爍
	if _time_left <= 10:
		_timer_label.add_theme_color_override("font_color", Color(1.0, 0.2, 0.2))
		var tween := create_tween()
		tween.tween_property(_timer_label, "modulate:a", 0.3, 0.2)
		tween.tween_property(_timer_label, "modulate:a", 1.0, 0.2)
	else:
		_timer_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.2))

func _update_leaderboard(top_players: Array) -> void:
	# 清除舊排行榜
	for child in _leaderboard_container.get_children():
		child.queue_free()
	
	if top_players.is_empty():
		return
	
	var rank_icons := ["🥇", "🥈", "🥉", "4️⃣", "5️⃣"]
	
	for i in range(min(top_players.size(), 3)):
		var player_data: Dictionary = top_players[i]
		var row := HBoxContainer.new()
		row.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		
		var rank_label := Label.new()
		rank_label.text = rank_icons[i] if i < rank_icons.size() else str(i + 1)
		rank_label.add_theme_font_size_override("font_size", 11)
		rank_label.custom_minimum_size = Vector2(24, 0)
		row.add_child(rank_label)
		
		var name_label := Label.new()
		name_label.text = player_data.get("player_name", "?")
		name_label.add_theme_font_size_override("font_size", 11)
		name_label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		var name_color := Color(1.0, 0.9, 0.2) if player_data.get("completed", false) else Color(0.85, 0.85, 0.85)
		name_label.add_theme_color_override("font_color", name_color)
		row.add_child(name_label)
		
		var progress_label := Label.new()
		var prog: int = player_data.get("progress", 0)
		progress_label.text = str(prog) + "/" + str(_current_target)
		progress_label.add_theme_font_size_override("font_size", 11)
		progress_label.add_theme_color_override("font_color", Color(0.7, 0.9, 0.7) if player_data.get("completed", false) else Color(0.7, 0.7, 0.7))
		row.add_child(progress_label)
		
		_leaderboard_container.add_child(row)

func _show_completion_flash() -> void:
	# 金色閃光效果
	var flash := ColorRect.new()
	flash.color = Color(_challenge_color.r, _challenge_color.g, _challenge_color.b, 0.4)
	flash.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	flash.position = Vector2(0, 0)
	add_child(flash)
	
	var tween := create_tween()
	tween.tween_property(flash, "color:a", 0.0, 0.5)
	tween.tween_callback(flash.queue_free)

func _show_reward_popup(reward: int, completed: bool, message: String) -> void:
	var popup := Label.new()
	popup.text = ("✅ " if completed else "💪 ") + message + "\n🪙 +" + str(reward)
	popup.add_theme_font_size_override("font_size", 14)
	popup.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2) if completed else Color(0.8, 0.8, 0.8))
	popup.position = Vector2(PANEL_WIDTH / 2 - 80, PANEL_HEIGHT / 2 - 20)
	popup.size = Vector2(160, 50)
	popup.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	add_child(popup)
	
	var tween := create_tween()
	tween.tween_property(popup, "position:y", popup.position.y - 30, 1.5)
	tween.parallel().tween_property(popup, "modulate:a", 0.0, 1.5)
	tween.tween_callback(popup.queue_free)

func _start_tick() -> void:
	_stop_tick()
	_tick_timer = get_tree().create_timer(1.0)
	_tick_timer.timeout.connect(_on_tick)

func _on_tick() -> void:
	if _time_left > 0:
		_time_left -= 1
		_update_timer()
	if _time_left > 0 and _is_visible:
		_start_tick()

func _stop_tick() -> void:
	_tick_timer = null
