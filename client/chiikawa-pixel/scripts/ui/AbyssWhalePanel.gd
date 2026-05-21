## AbyssWhalePanel.gd — 深淵巨鯨全服 Boss 挑戰面板（DAY-164）
## 深淵巨鯨出現時全服廣播，玩家合力攻擊，按傷害貢獻比例分配深淵寶藏
## 視覺：深海藍黑主題 + 全服 HP 進度條 + 貢獻排行榜 + 結果彈窗
## 業界依據：Fishing Frenzy Chapter 3 2026「Boss Fish endgame content」
extends Node2D

var _pixel_font: Font = null
var _hp_bar_bg: ColorRect = null
var _hp_bar_fill: ColorRect = null
var _hp_label: Label = null
var _whale_banner: Node = null
var _is_active: bool = false
var _total_hp: int = 500
var _current_hp: int = 500

func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

func _connect_signals() -> void:
	if GameManager.has_signal("abyss_whale"):
		GameManager.abyss_whale.connect(_on_abyss_whale)

## 處理深淵巨鯨事件
func _on_abyss_whale(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	match phase:
		"whale_spawn":
			_show_whale_spawn(data)
		"whale_hp_update":
			_update_hp(data)
		"whale_killed":
			_show_whale_killed(data)
		"whale_reward":
			_show_my_reward(data)

## 深淵巨鯨出現
func _show_whale_spawn(data: Dictionary) -> void:
	_is_active = true
	_total_hp = data.get("total_hp", 500)
	_current_hp = _total_hp

	# 全螢幕深藍閃光
	var flash := ColorRect.new()
	flash.position = Vector2(-640, -360)
	flash.size = Vector2(1280, 720)
	flash.color = Color(0.0, 0.1, 0.4, 0.5)
	add_child(flash)
	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", 0.0, 0.8)
	flash_tween.tween_callback(func():
		if is_instance_valid(flash): flash.queue_free()
	)

	# 頂部橫幅（深海主題）
	var banner_bg := ColorRect.new()
	banner_bg.name = "WhaleBanner"
	banner_bg.position = Vector2(-640, -360)
	banner_bg.size = Vector2(1280, 56)
	banner_bg.color = Color(0.0, 0.05, 0.2, 0.92)
	add_child(banner_bg)
	_whale_banner = banner_bg

	var banner_label := Label.new()
	banner_label.text = "🐋 深淵巨鯨出現！全服合力擊破，按貢獻分配深淵寶藏！"
	banner_label.position = Vector2(0, 8)
	banner_label.size = Vector2(1280, 40)
	banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	banner_label.add_theme_color_override("font_color", Color(0.4, 0.8, 1.0))
	if _pixel_font:
		banner_label.add_theme_font_override("font", _pixel_font)
		banner_label.add_theme_font_size_override("font_size", 18)
	banner_bg.add_child(banner_label)

	# 橫幅滑入動畫
	banner_bg.position.y = -420
	var banner_tween = banner_bg.create_tween()
	banner_tween.tween_property(banner_bg, "position:y", -360, 0.4).set_ease(Tween.EASE_OUT)

	# 底部 HP 進度條（常駐顯示）
	_create_hp_bar()

## 建立底部 HP 進度條
func _create_hp_bar() -> void:
	# 移除舊的
	if is_instance_valid(_hp_bar_bg):
		_hp_bar_bg.queue_free()

	var bar_bg := ColorRect.new()
	bar_bg.name = "WhaleHPBarBG"
	bar_bg.position = Vector2(-640, 300)
	bar_bg.size = Vector2(1280, 28)
	bar_bg.color = Color(0.0, 0.0, 0.1, 0.85)
	add_child(bar_bg)
	_hp_bar_bg = bar_bg

	# HP 填充條
	var bar_fill := ColorRect.new()
	bar_fill.name = "WhaleHPBarFill"
	bar_fill.position = Vector2(4, 4)
	bar_fill.size = Vector2(1272, 20)
	bar_fill.color = Color(0.1, 0.5, 1.0)
	bar_bg.add_child(bar_fill)
	_hp_bar_fill = bar_fill

	# HP 文字
	var hp_label := Label.new()
	hp_label.name = "WhaleHPLabel"
	hp_label.text = "🐋 深淵巨鯨 HP: %d / %d" % [_current_hp, _total_hp]
	hp_label.position = Vector2(0, 4)
	hp_label.size = Vector2(1280, 20)
	hp_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	hp_label.add_theme_color_override("font_color", Color(0.9, 0.95, 1.0))
	if _pixel_font:
		hp_label.add_theme_font_override("font", _pixel_font)
		hp_label.add_theme_font_size_override("font_size", 13)
	bar_bg.add_child(hp_label)
	_hp_label = hp_label

## 更新 HP 進度條
func _update_hp(data: Dictionary) -> void:
	_current_hp = data.get("current_hp", _current_hp)
	_total_hp = data.get("total_hp", _total_hp)
	var hp_percent: float = data.get("hp_percent", 1.0)

	if not is_instance_valid(_hp_bar_fill):
		return

	# 更新填充寬度
	var target_width = 1272.0 * hp_percent
	var tween = _hp_bar_fill.create_tween()
	tween.tween_property(_hp_bar_fill, "size:x", target_width, 0.2)

	# HP 顏色：高→藍，中→青，低→紅
	if hp_percent > 0.6:
		_hp_bar_fill.color = Color(0.1, 0.5, 1.0)
	elif hp_percent > 0.3:
		_hp_bar_fill.color = Color(0.0, 0.8, 0.8)
	else:
		_hp_bar_fill.color = Color(1.0, 0.2, 0.2)
		# 低血量閃爍警告
		var flash_tween = _hp_bar_fill.create_tween().set_loops(3)
		flash_tween.tween_property(_hp_bar_fill, "color:a", 0.4, 0.15)
		flash_tween.tween_property(_hp_bar_fill, "color:a", 1.0, 0.15)

	# 更新文字
	if is_instance_valid(_hp_label):
		_hp_label.text = "🐋 深淵巨鯨 HP: %d / %d" % [_current_hp, _total_hp]

	# 受擊浮動文字
	var attacker_id: String = data.get("attacker_id", "")
	var my_id = NetworkManager.get_player_id() if NetworkManager.has_method("get_player_id") else ""
	if attacker_id == my_id:
		_show_damage_text()

## 顯示我的傷害文字
func _show_damage_text() -> void:
	var dmg_label := Label.new()
	dmg_label.text = "⚔️ 命中！"
	dmg_label.position = Vector2(randf_range(-200, 200), 260)
	dmg_label.add_theme_color_override("font_color", Color(0.4, 0.8, 1.0))
	if _pixel_font:
		dmg_label.add_theme_font_override("font", _pixel_font)
		dmg_label.add_theme_font_size_override("font_size", 14)
	add_child(dmg_label)
	var tween = dmg_label.create_tween()
	tween.tween_property(dmg_label, "position:y", 220, 0.6)
	tween.parallel().tween_property(dmg_label, "modulate:a", 0.0, 0.6)
	tween.tween_callback(func():
		if is_instance_valid(dmg_label): dmg_label.queue_free()
	)

## 深淵巨鯨被擊破
func _show_whale_killed(data: Dictionary) -> void:
	_is_active = false
	var killer_name: String = data.get("killer_name", "")
	var entries: Array = data.get("entries", [])

	# 移除 HP 條
	if is_instance_valid(_hp_bar_bg):
		var tween = _hp_bar_bg.create_tween()
		tween.tween_property(_hp_bar_bg, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func():
			if is_instance_valid(_hp_bar_bg): _hp_bar_bg.queue_free()
		)

	# 移除橫幅
	if is_instance_valid(_whale_banner):
		var tween2 = _whale_banner.create_tween()
		tween2.tween_property(_whale_banner, "position:y", -420, 0.4)
		tween2.tween_callback(func():
			if is_instance_valid(_whale_banner): _whale_banner.queue_free()
		)

	# 全螢幕金色爆炸閃光
	var flash := ColorRect.new()
	flash.position = Vector2(-640, -360)
	flash.size = Vector2(1280, 720)
	flash.color = Color(0.0, 0.4, 1.0, 0.6)
	add_child(flash)
	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color", Color(1.0, 0.8, 0.0, 0.0), 1.0)
	flash_tween.tween_callback(func():
		if is_instance_valid(flash): flash.queue_free()
	)

	# 右側滑入結果彈窗
	_show_result_panel(killer_name, entries)

## 顯示結果彈窗
func _show_result_panel(killer_name: String, entries: Array) -> void:
	var panel_bg := ColorRect.new()
	panel_bg.position = Vector2(1400, -200)
	panel_bg.size = Vector2(340, 420)
	panel_bg.color = Color(0.0, 0.05, 0.2, 0.95)
	add_child(panel_bg)

	# 標題
	var title := Label.new()
	title.text = "🐋 深淵寶藏分配"
	title.position = Vector2(10, 10)
	title.size = Vector2(320, 30)
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_color_override("font_color", Color(0.4, 0.8, 1.0))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 16)
	panel_bg.add_child(title)

	# 擊破者
	var killer_label := Label.new()
	killer_label.text = "擊破者：%s" % killer_name
	killer_label.position = Vector2(10, 44)
	killer_label.size = Vector2(320, 22)
	killer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	killer_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	if _pixel_font:
		killer_label.add_theme_font_override("font", _pixel_font)
		killer_label.add_theme_font_size_override("font_size", 13)
	panel_bg.add_child(killer_label)

	# 分隔線
	var sep := ColorRect.new()
	sep.position = Vector2(10, 70)
	sep.size = Vector2(320, 2)
	sep.color = Color(0.2, 0.5, 0.8, 0.8)
	panel_bg.add_child(sep)

	# 貢獻者列表（最多顯示 8 名）
	var rank_colors = [Color(1.0, 0.85, 0.0), Color(0.8, 0.8, 0.8), Color(0.8, 0.5, 0.2)]
	var display_count = min(entries.size(), 8)
	for i in range(display_count):
		var entry = entries[i]
		var rank = entry.get("rank", i + 1)
		var name_str = entry.get("player_name", "???")
		var ratio = entry.get("ratio", 0.0)
		var bonus = entry.get("bonus", 0)

		var row := Label.new()
		var rank_icon = "🥇" if rank == 1 else ("🥈" if rank == 2 else ("🥉" if rank == 3 else "#%d" % rank))
		row.text = "%s %s  %.0f%%  +%d" % [rank_icon, name_str.left(8), ratio * 100, bonus]
		row.position = Vector2(10, 78 + i * 38)
		row.size = Vector2(320, 34)
		row.horizontal_alignment = HORIZONTAL_ALIGNMENT_LEFT
		var color = rank_colors[min(rank - 1, 2)] if rank <= 3 else Color(0.8, 0.9, 1.0)
		row.add_theme_color_override("font_color", color)
		if _pixel_font:
			row.add_theme_font_override("font", _pixel_font)
			row.add_theme_font_size_override("font_size", 13)
		panel_bg.add_child(row)

	# 滑入動畫
	var tween = panel_bg.create_tween()
	tween.tween_property(panel_bg, "position:x", 320, 0.5).set_ease(Tween.EASE_OUT)

	# 5 秒後淡出
	await get_tree().create_timer(5.0).timeout
	if is_instance_valid(panel_bg):
		var fade_tween = panel_bg.create_tween()
		fade_tween.tween_property(panel_bg, "modulate:a", 0.0, 0.6)
		fade_tween.tween_callback(func():
			if is_instance_valid(panel_bg): panel_bg.queue_free()
		)

## 顯示我的個人獎勵
func _show_my_reward(data: Dictionary) -> void:
	var my_rank: int = data.get("my_rank", 0)
	var my_bonus: int = data.get("my_bonus", 0)
	var my_damage: int = data.get("my_damage", 0)
	var my_ratio: float = data.get("my_ratio", 0.0)

	if my_bonus <= 0:
		return

	# 中央大獎彈窗
	var reward_bg := ColorRect.new()
	reward_bg.position = Vector2(-200, -120)
	reward_bg.size = Vector2(400, 240)
	reward_bg.color = Color(0.0, 0.05, 0.25, 0.95)
	add_child(reward_bg)

	# 邊框
	var border := ColorRect.new()
	border.position = Vector2(-2, -2)
	border.size = Vector2(404, 244)
	border.color = Color(0.2, 0.6, 1.0, 0.9)
	border.z_index = -1
	reward_bg.add_child(border)

	var rank_text = "🥇 第一名" if my_rank == 1 else ("🥈 第二名" if my_rank == 2 else ("🥉 第三名" if my_rank == 3 else "第 %d 名" % my_rank))
	var title_color = Color(1.0, 0.85, 0.0) if my_rank <= 3 else Color(0.4, 0.8, 1.0)

	var title := Label.new()
	title.text = "🐋 深淵寶藏"
	title.position = Vector2(0, 16)
	title.size = Vector2(400, 36)
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_color_override("font_color", Color(0.4, 0.8, 1.0))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 20)
	reward_bg.add_child(title)

	var rank_label := Label.new()
	rank_label.text = rank_text
	rank_label.position = Vector2(0, 58)
	rank_label.size = Vector2(400, 30)
	rank_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	rank_label.add_theme_color_override("font_color", title_color)
	if _pixel_font:
		rank_label.add_theme_font_override("font", _pixel_font)
		rank_label.add_theme_font_size_override("font_size", 16)
	reward_bg.add_child(rank_label)

	var contrib_label := Label.new()
	contrib_label.text = "貢獻傷害：%d（%.0f%%）" % [my_damage, my_ratio * 100]
	contrib_label.position = Vector2(0, 94)
	contrib_label.size = Vector2(400, 26)
	contrib_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	contrib_label.add_theme_color_override("font_color", Color(0.7, 0.9, 1.0))
	if _pixel_font:
		contrib_label.add_theme_font_override("font", _pixel_font)
		contrib_label.add_theme_font_size_override("font_size", 14)
	reward_bg.add_child(contrib_label)

	var bonus_label := Label.new()
	bonus_label.text = "+%d 金幣" % my_bonus
	bonus_label.position = Vector2(0, 128)
	bonus_label.size = Vector2(400, 48)
	bonus_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	bonus_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	if _pixel_font:
		bonus_label.add_theme_font_override("font", _pixel_font)
		bonus_label.add_theme_font_size_override("font_size", 28)
	reward_bg.add_child(bonus_label)

	# 彈跳動畫
	reward_bg.scale = Vector2(0.5, 0.5)
	var tween = reward_bg.create_tween()
	tween.tween_property(reward_bg, "scale", Vector2(1.1, 1.1), 0.25).set_ease(Tween.EASE_OUT)
	tween.tween_property(reward_bg, "scale", Vector2(1.0, 1.0), 0.1)

	# 高貢獻者雙閃光
	if my_ratio >= 0.3:
		for _i in range(2):
			var extra_flash := ColorRect.new()
			extra_flash.position = Vector2(-640, -360)
			extra_flash.size = Vector2(1280, 720)
			extra_flash.color = Color(0.0, 0.5, 1.0, 0.3)
			add_child(extra_flash)
			var ef_tween = extra_flash.create_tween()
			ef_tween.tween_property(extra_flash, "color:a", 0.0, 0.4)
			ef_tween.tween_callback(func():
				if is_instance_valid(extra_flash): extra_flash.queue_free()
			)

	# 4 秒後淡出
	await get_tree().create_timer(4.0).timeout
	if is_instance_valid(reward_bg):
		var fade_tween = reward_bg.create_tween()
		fade_tween.tween_property(reward_bg, "modulate:a", 0.0, 0.5)
		fade_tween.tween_callback(func():
			if is_instance_valid(reward_bg): reward_bg.queue_free()
		)
