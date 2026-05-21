## MissionStreakPanel.gd - DAY-086 / DAY-120
## 每日任務連續完成獎勵 UI
## 全部任務完成後顯示連續天數和獎勵
## DAY-120：加入寬限期保護通知
extends Node2D

const PANEL_W := 300
const PANEL_H := 100

var _font: FontFile

func setup(font: FontFile) -> void:
	_font = font
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("mission_streak_bonus"):
		GameManager.mission_streak_bonus.connect(_on_mission_streak_bonus)
	# 寬限期保護通知（DAY-120）
	if GameManager.has_signal("mission_mercy_protected"):
		GameManager.mission_mercy_protected.connect(_on_mission_mercy_protected)

func _on_mission_streak_bonus(data: Dictionary) -> void:
	var streak: int = data.get("streak", 1)
	var reward: int = data.get("reward", 0)
	var label: String = data.get("label", "完成所有任務")
	var mercy_used: bool = data.get("mercy_used", false)

	_show_popup(streak, reward, label, mercy_used)

# 寬限期保護通知（DAY-120）
func _on_mission_mercy_protected(data: Dictionary) -> void:
	var streak: int = data.get("streak", 0)
	var message: String = data.get("message", "🛡️ 連續記錄被保護了！")
	_show_mercy_popup(streak, message)

func _show_mercy_popup(streak: int, message: String) -> void:
	var popup = Node2D.new()
	popup.position = Vector2(640, 280)
	popup.scale = Vector2(0.0, 0.0)
	add_child(popup)

	# 背景（藍紫色，代表保護）
	var bg = ColorRect.new()
	bg.size = Vector2(PANEL_W, 80)
	bg.position = Vector2(-PANEL_W / 2, -40)
	bg.color = Color(0.08, 0.05, 0.20, 0.95)
	popup.add_child(bg)

	# 頂部邊框（紫色）
	var border = ColorRect.new()
	border.size = Vector2(PANEL_W, 3)
	border.position = Vector2(-PANEL_W / 2, -40)
	border.color = Color(0.6, 0.3, 1.0, 1.0)
	popup.add_child(border)

	# 保護圖示 + 訊息
	var msg_lbl = Label.new()
	msg_lbl.text = message
	msg_lbl.position = Vector2(-PANEL_W / 2 + 8, -30)
	msg_lbl.add_theme_color_override("font_color", Color(0.8, 0.6, 1.0))
	if _font:
		msg_lbl.add_theme_font_override("font", _font)
		msg_lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(msg_lbl)

	# 連續天數
	var streak_lbl = Label.new()
	streak_lbl.text = "連續 %d 天記錄保留 ✓" % streak
	streak_lbl.position = Vector2(-PANEL_W / 2 + 8, -8)
	streak_lbl.add_theme_color_override("font_color", Color.WHITE)
	if _font:
		streak_lbl.add_theme_font_override("font", _font)
		streak_lbl.add_theme_font_size_override("font_size", 13)
	popup.add_child(streak_lbl)

	# 動畫
	var tween = create_tween()
	tween.tween_property(popup, "scale", Vector2(1.05, 1.05), 0.2)
	tween.tween_property(popup, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(4.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.4)
	tween.tween_callback(popup.queue_free)

func _show_popup(streak: int, reward: int, label: String, mercy_used: bool = false) -> void:
	var popup = Node2D.new()
	popup.position = Vector2(640, 300)
	popup.scale = Vector2(0.0, 0.0)
	add_child(popup)

	# 背景（連續天數越高越金色；寬限期使用時用紫色）
	var bg = ColorRect.new()
	bg.size = Vector2(PANEL_W, PANEL_H)
	bg.position = Vector2(-PANEL_W / 2, -PANEL_H / 2)
	if mercy_used:
		bg.color = Color(0.10, 0.05, 0.18, 0.97)  # 深紫色（寬限期）
	elif streak >= 7:
		bg.color = Color(0.15, 0.10, 0.02, 0.97)  # 深金色
	elif streak >= 3:
		bg.color = Color(0.05, 0.12, 0.05, 0.97)  # 深綠色
	else:
		bg.color = Color(0.03, 0.06, 0.18, 0.97)  # 深藍色
	popup.add_child(bg)

	# 頂部邊框
	var border = ColorRect.new()
	border.size = Vector2(PANEL_W, 3)
	border.position = Vector2(-PANEL_W / 2, -PANEL_H / 2)
	if mercy_used:
		border.color = Color(0.6, 0.3, 1.0, 1.0)  # 紫色（寬限期）
	elif streak >= 7:
		border.color = Color(1.0, 0.85, 0.1, 1.0)
	elif streak >= 3:
		border.color = Color(0.3, 1.0, 0.3, 1.0)
	else:
		border.color = Color(0.3, 0.7, 1.0, 1.0)
	popup.add_child(border)

	# 標題
	var title_lbl = Label.new()
	title_lbl.text = "📋 所有任務完成！"
	title_lbl.position = Vector2(-PANEL_W / 2 + 8, -PANEL_H / 2 + 6)
	title_lbl.add_theme_color_override("font_color", Color(0.8, 1.0, 0.8))
	if _font:
		title_lbl.add_theme_font_override("font", _font)
		title_lbl.add_theme_font_size_override("font_size", 13)
	popup.add_child(title_lbl)

	# 連續天數
	var streak_lbl = Label.new()
	streak_lbl.text = "連續第 %d 天 🔥" % streak
	streak_lbl.position = Vector2(-PANEL_W / 2 + 8, -PANEL_H / 2 + 26)
	if mercy_used:
		streak_lbl.add_theme_color_override("font_color", Color(0.8, 0.6, 1.0))  # 紫色
	elif streak >= 7:
		streak_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
	else:
		streak_lbl.add_theme_color_override("font_color", Color.WHITE)
	if _font:
		streak_lbl.add_theme_font_override("font", _font)
		streak_lbl.add_theme_font_size_override("font_size", 18)
	popup.add_child(streak_lbl)

	# 標籤
	var label_lbl = Label.new()
	label_lbl.text = label
	label_lbl.position = Vector2(-PANEL_W / 2 + 8, -PANEL_H / 2 + 50)
	label_lbl.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	if _font:
		label_lbl.add_theme_font_override("font", _font)
		label_lbl.add_theme_font_size_override("font_size", 12)
	popup.add_child(label_lbl)

	# 獎勵
	var reward_lbl = Label.new()
	reward_lbl.text = "+%d 🪙" % reward
	reward_lbl.position = Vector2(PANEL_W / 2 - 90, -PANEL_H / 2 + 30)
	reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
	if _font:
		reward_lbl.add_theme_font_override("font", _font)
		reward_lbl.add_theme_font_size_override("font_size", 20)
	popup.add_child(reward_lbl)

	# 動畫
	var tween = create_tween()
	tween.tween_property(popup, "scale", Vector2(1.1, 1.1), 0.15)
	tween.tween_property(popup, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(3.0)
	tween.tween_property(popup, "position:y", popup.position.y - 25, 0.4)
	tween.tween_property(popup, "modulate:a", 0.0, 0.3)
	tween.tween_callback(popup.queue_free)
