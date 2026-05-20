## ChallengePanel.gd - DAY-085
## 隱藏挑戰系統 UI：挑戰解鎖時顯示驚喜通知
## 設計原則：隱藏挑戰解鎖時要有「驚喜感」，比普通成就更誇張
extends Node2D

const PANEL_W := 280
const PANEL_H := 80

var _font: FontFile
var _queue: Array = []
var _is_showing: bool = false

func setup(font: FontFile) -> void:
	_font = font
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("challenge_unlocked"):
		GameManager.challenge_unlocked.connect(_on_challenge_unlocked)

func _on_challenge_unlocked(data: Dictionary) -> void:
	_queue.append(data)
	if not _is_showing:
		_show_next()

func _show_next() -> void:
	if _queue.is_empty():
		_is_showing = false
		return

	_is_showing = true
	var data: Dictionary = _queue.pop_front()
	_show_challenge_popup(data)

func _show_challenge_popup(data: Dictionary) -> void:
	var was_hidden: bool = data.get("was_hidden", false)
	var name_str: String = data.get("name", "")
	var desc_str: String = data.get("description", "")
	var icon_str: String = data.get("icon", "🏆")
	var reward: int = data.get("reward", 0)

	# 建立彈窗容器
	var popup = Node2D.new()
	popup.position = Vector2(640, 200)
	popup.scale = Vector2(0.0, 0.0)
	add_child(popup)

	# 背景（隱藏挑戰用金色，普通挑戰用藍色）
	var bg = ColorRect.new()
	bg.size = Vector2(PANEL_W, PANEL_H)
	bg.position = Vector2(-PANEL_W / 2, -PANEL_H / 2)
	if was_hidden:
		bg.color = Color(0.15, 0.10, 0.02, 0.97)  # 深金色
	else:
		bg.color = Color(0.02, 0.08, 0.18, 0.97)  # 深藍色
	popup.add_child(bg)

	# 頂部邊框
	var border = ColorRect.new()
	border.size = Vector2(PANEL_W, 3)
	border.position = Vector2(-PANEL_W / 2, -PANEL_H / 2)
	if was_hidden:
		border.color = Color(1.0, 0.85, 0.1, 1.0)  # 金色
	else:
		border.color = Color(0.3, 0.7, 1.0, 1.0)   # 藍色
	popup.add_child(border)

	# 標題（隱藏挑戰有特殊標記）
	var title_lbl = Label.new()
	if was_hidden:
		title_lbl.text = "🎉 隱藏挑戰解鎖！"
		title_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
	else:
		title_lbl.text = "✅ 挑戰完成！"
		title_lbl.add_theme_color_override("font_color", Color(0.5, 1.0, 0.5))
	title_lbl.position = Vector2(-PANEL_W / 2 + 8, -PANEL_H / 2 + 6)
	if _font:
		title_lbl.add_theme_font_override("font", _font)
		title_lbl.add_theme_font_size_override("font_size", 13)
	popup.add_child(title_lbl)

	# 挑戰名稱
	var name_lbl = Label.new()
	name_lbl.text = "%s %s" % [icon_str, name_str]
	name_lbl.position = Vector2(-PANEL_W / 2 + 8, -PANEL_H / 2 + 24)
	name_lbl.add_theme_color_override("font_color", Color.WHITE)
	if _font:
		name_lbl.add_theme_font_override("font", _font)
		name_lbl.add_theme_font_size_override("font_size", 16)
	popup.add_child(name_lbl)

	# 描述
	var desc_lbl = Label.new()
	desc_lbl.text = desc_str
	desc_lbl.position = Vector2(-PANEL_W / 2 + 8, -PANEL_H / 2 + 44)
	desc_lbl.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	if _font:
		desc_lbl.add_theme_font_override("font", _font)
		desc_lbl.add_theme_font_size_override("font_size", 11)
	popup.add_child(desc_lbl)

	# 獎勵
	if reward > 0:
		var reward_lbl = Label.new()
		reward_lbl.text = "+%d 🪙" % reward
		reward_lbl.position = Vector2(PANEL_W / 2 - 80, -PANEL_H / 2 + 28)
		reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
		if _font:
			reward_lbl.add_theme_font_override("font", _font)
			reward_lbl.add_theme_font_size_override("font_size", 16)
		popup.add_child(reward_lbl)

	# 動畫：彈出 → 停留 → 淡出
	var tween = create_tween()
	# 彈出
	tween.tween_property(popup, "scale", Vector2(1.1, 1.1), 0.15)
	tween.tween_property(popup, "scale", Vector2(1.0, 1.0), 0.1)
	# 停留（隱藏挑戰停留更久）
	var stay_time := 3.5 if was_hidden else 2.5
	tween.tween_interval(stay_time)
	# 淡出（向上移動）
	tween.tween_property(popup, "position:y", popup.position.y - 30, 0.4)
	tween.tween_property(popup, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func():
		popup.queue_free()
		_show_next()
	)

	# 隱藏挑戰額外特效：金色粒子
	if was_hidden:
		_spawn_gold_particles(popup)

func _spawn_gold_particles(parent: Node2D) -> void:
	# 生成 8 個金色粒子從中心散開
	for i in range(8):
		var particle = ColorRect.new()
		particle.size = Vector2(4, 4)
		particle.color = Color(1.0, 0.85, 0.1, 1.0)
		particle.position = Vector2(-2, -2)
		parent.add_child(particle)

		var angle := i * TAU / 8
		var dist := 60.0
		var target_pos = Vector2(cos(angle) * dist, sin(angle) * dist)

		var tween = parent.create_tween()
		tween.tween_property(particle, "position", target_pos, 0.5)
		tween.parallel().tween_property(particle, "modulate:a", 0.0, 0.5)
		tween.tween_callback(particle.queue_free)
