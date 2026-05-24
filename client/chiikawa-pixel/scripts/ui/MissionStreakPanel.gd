## MissionStreakPanel.gd - DAY-086 / DAY-120
## 瘥隞餃????摰?? UI
## ?券隞餃?摰?敺＊蝷粹??憭拇????
## DAY-120嚗??亙祝??靽風?
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
	# 撖祇???霅琿嚗AY-120嚗?
	if GameManager.has_signal("mission_mercy_protected"):
		GameManager.mission_mercy_protected.connect(_on_mission_mercy_protected)

func _on_mission_streak_bonus(data: Dictionary) -> void:
	var streak: int = data.get("streak", 1)
	var reward: int = data.get("reward", 0)
	var label: String = data.get("label", "")"
	var mercy_used: bool = data.get("mercy_used", false)

	_show_popup(streak, reward, label, mercy_used)

# 撖祇???霅琿嚗AY-120嚗?
func _on_mission_mercy_protected(data: Dictionary) -> void:
	var streak: int = data.get("streak", 0)
	var message: String = data.get("message", "?儭????閮?鋡思?霅瑚?嚗?)"
	_show_mercy_popup(streak, message)

func _show_mercy_popup(streak: int, message: String) -> void:
	var popup = Node2D.new()
	popup.position = Vector2(640, 280)
	popup.scale = Vector2(0.0, 0.0)
	add_child(popup)

	# ?嚗?蝝怨嚗誨銵其?霅瘀?
	var bg = ColorRect.new()
	bg.size = Vector2(PANEL_W, 80)
	bg.position = Vector2(-PANEL_W / 2, -40)
	bg.color = Color(0.08, 0.05, 0.20, 0.95)
	popup.add_child(bg)

	# ???嚗換?莎?
	var border = ColorRect.new()
	border.size = Vector2(PANEL_W, 3)
	border.position = Vector2(-PANEL_W / 2, -40)
	border.color = Color(0.6, 0.3, 1.0, 1.0)
	popup.add_child(border)

	# 靽風?內 + 閮
	var msg_lbl = Label.new()
	msg_lbl.text = message
	msg_lbl.position = Vector2(-PANEL_W / 2 + 8, -30)
	msg_lbl.add_theme_color_override("font_color", Color(0.8, 0.6, 1.0))
	if _font:
		msg_lbl.add_theme_font_override("font", _font)
		msg_lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(msg_lbl)

	# ???憭拇
	var streak_lbl = Label.new()
	streak_lbl.text = "??? %d 憭抵??????? % streak"
	streak_lbl.position = Vector2(-PANEL_W / 2 + 8, -8)
	streak_lbl.add_theme_color_override("font_color", Color.WHITE)
	if _font:
		streak_lbl.add_theme_font_override("font", _font)
		streak_lbl.add_theme_font_size_override("font_size", 13)
	popup.add_child(streak_lbl)

	# ?
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

	# ?嚗??憭拇頞?頞??莎?撖祇??蝙?冽??函換?莎?
	var bg = ColorRect.new()
	bg.size = Vector2(PANEL_W, PANEL_H)
	bg.position = Vector2(-PANEL_W / 2, -PANEL_H / 2)
	if mercy_used:
		bg.color = Color(0.10, 0.05, 0.18, 0.97)  # 瘛梁換?莎?撖祇???
	elif streak >= 7:
		bg.color = Color(0.15, 0.10, 0.02, 0.97)  # 瘛梢???
	elif streak >= 3:
		bg.color = Color(0.05, 0.12, 0.05, 0.97)  # 瘛梁???
	else:
		bg.color = Color(0.03, 0.06, 0.18, 0.97)  # 瘛梯???
	popup.add_child(bg)

	# ???
	var border = ColorRect.new()
	border.size = Vector2(PANEL_W, 3)
	border.position = Vector2(-PANEL_W / 2, -PANEL_H / 2)
	if mercy_used:
		border.color = Color(0.6, 0.3, 1.0, 1.0)  # 蝝怨嚗祝??嚗?
	elif streak >= 7:
		border.color = Color(1.0, 0.85, 0.1, 1.0)
	elif streak >= 3:
		border.color = Color(0.3, 1.0, 0.3, 1.0)
	else:
		border.color = Color(0.3, 0.7, 1.0, 1.0)
	popup.add_child(border)

	# 璅?
	var title_lbl = Label.new()
	title_lbl.text = "?? ??遙????"
	title_lbl.position = Vector2(-PANEL_W / 2 + 8, -PANEL_H / 2 + 6)
	title_lbl.add_theme_color_override("font_color", Color(0.8, 1.0, 0.8))
	if _font:
		title_lbl.add_theme_font_override("font", _font)
		title_lbl.add_theme_font_size_override("font_size", 13)
	popup.add_child(title_lbl)

	# ???憭拇
	var streak_lbl = Label.new()
	streak_lbl.text = "???蝚?%d 憭??" % streak
	streak_lbl.position = Vector2(-PANEL_W / 2 + 8, -PANEL_H / 2 + 26)
	if mercy_used:
		streak_lbl.add_theme_color_override("font_color", Color(0.8, 0.6, 1.0))  # 蝝怨
	elif streak >= 7:
		streak_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
	else:
		streak_lbl.add_theme_color_override("font_color", Color.WHITE)
	if _font:
		streak_lbl.add_theme_font_override("font", _font)
		streak_lbl.add_theme_font_size_override("font_size", 18)
	popup.add_child(streak_lbl)

	# 璅惜
	var label_lbl = Label.new()
	label_lbl.text = label
	label_lbl.position = Vector2(-PANEL_W / 2 + 8, -PANEL_H / 2 + 50)
	label_lbl.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	if _font:
		label_lbl.add_theme_font_override("font", _font)
		label_lbl.add_theme_font_size_override("font_size", 12)
	popup.add_child(label_lbl)

	# ?
	var reward_lbl = Label.new()
	reward_lbl.text = "+%d ??" % reward
	reward_lbl.position = Vector2(PANEL_W / 2 - 90, -PANEL_H / 2 + 30)
	reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
	if _font:
		reward_lbl.add_theme_font_override("font", _font)
		reward_lbl.add_theme_font_size_override("font_size", 20)
	popup.add_child(reward_lbl)

	# ?
	var tween = create_tween()
	tween.tween_property(popup, "scale", Vector2(1.1, 1.1), 0.15)
	tween.tween_property(popup, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(3.0)
	tween.tween_property(popup, "position:y", popup.position.y - 25, 0.4)
	tween.tween_property(popup, "modulate:a", 0.0, 0.3)
	tween.tween_callback(popup.queue_free)
