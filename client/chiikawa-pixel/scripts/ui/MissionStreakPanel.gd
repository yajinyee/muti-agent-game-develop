## MissionStreakPanel.gd - DAY-086
## 每日任務連續完成獎勵 UI
## 全部任務完成後顯示連續天數和獎勵
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

func _on_mission_streak_bonus(data: Dictionary) -> void:
	var streak: int = data.get("streak", 1)
	var reward: int = data.get("reward", 0)
	var label: String = data.get("label", "完成所有任務")
	var new_balance: int = data.get("new_balance", 0)

	_show_popup(streak, reward, label)

func _show_popup(streak: int, reward: int, label: String) -> void:
	var popup = Node2D.new()
	popup.position = Vector2(640, 300)
	popup.scale = Vector2(0.0, 0.0)
	add_child(popup)

	# 背景（連續天數越高越金色）
	var bg = ColorRect.new()
	bg.size = Vector2(PANEL_W, PANEL_H)
	bg.position = Vector2(-PANEL_W / 2, -PANEL_H / 2)
	if streak >= 7:
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
	if streak >= 7:
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
	if streak >= 7:
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
