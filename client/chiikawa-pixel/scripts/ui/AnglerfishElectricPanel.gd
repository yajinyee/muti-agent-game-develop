## AnglerfishElectricPanel.gd — 巨型鮟鱇魚電擊寶箱面板（DAY-196）
## 業界依據：JILI Mega Fishing「Giant Anglerfish can shoot electricity to open treasure chests,
## giant crocodiles awaken to hunt fish on the fish farm to accumulate big prizes!」
##
## 視覺設計：
##   - 深海藍紫主題（#4B0082 + #9400D3 + #00BFFF）
##   - anglerfish_appear：藍紫色雙閃光 + 頂部橫幅 + 電擊計數器 + 累積獎池顯示
##   - zap_N：電弧線（鋸齒形）+ 目標位置閃光 + 獎勵浮動文字
##   - 寶箱開箱：金色強閃光 + 「💰 寶箱開箱！×N倍」大字彈跳
##   - super_zap_start：全螢幕藍白強閃光 + 「⚡⚡⚡ 超級電擊！」橫幅
##   - super_zap_N：每個目標依序閃光（電流蔓延感）
##   - anglerfish_killed：金色閃光 + 右側滑入結算彈窗（電擊次數/獎池/加成/總獎勵）
##   - anglerfish_leave：淡出所有 UI
extends CanvasLayer

var _panel: Control
var _banner: Label
var _zap_counter: Label      # 電擊計數器
var _pool_label: Label       # 累積獎池顯示
var _result_popup: Control

func _ready() -> void:
	layer = 49
	_build_ui()

func _build_ui() -> void:
	_panel = Control.new()
	_panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	_panel.visible = false
	_panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_panel)

	# 頂部橫幅
	_banner = Label.new()
	_banner.text = "⚡ 巨型鮟鱇魚電擊！"
	_banner.add_theme_font_size_override("font_size", 22)
	_banner.add_theme_color_override("font_color", Color("#00BFFF"))
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.position = Vector2(0, 15)
	_banner.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.visible = false
	add_child(_banner)

	# 電擊計數器
	_zap_counter = Label.new()
	_zap_counter.text = "電擊: 0/8"
	_zap_counter.add_theme_font_size_override("font_size", 18)
	_zap_counter.add_theme_color_override("font_color", Color("#9400D3"))
	_zap_counter.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_zap_counter.position = Vector2(0, 45)
	_zap_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_zap_counter.visible = false
	add_child(_zap_counter)

	# 累積獎池
	_pool_label = Label.new()
	_pool_label.text = "獎池: 0"
	_pool_label.add_theme_font_size_override("font_size", 16)
	_pool_label.add_theme_color_override("font_color", Color("#FFD700"))
	_pool_label.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_pool_label.position = Vector2(0, 68)
	_pool_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_pool_label.visible = false
	add_child(_pool_label)

	# 結算彈窗
	_result_popup = Control.new()
	_result_popup.visible = false
	_result_popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_popup.size = Vector2(300, 200)
	_result_popup.position = Vector2(-320, -100)
	add_child(_result_popup)

	var popup_bg = ColorRect.new()
	popup_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_bg.color = Color(0.05, 0.0, 0.15, 0.93)
	_result_popup.add_child(popup_bg)

	var popup_label = Label.new()
	popup_label.name = "ResultLabel"
	popup_label.add_theme_font_size_override("font_size", 16)
	popup_label.add_theme_color_override("font_color", Color("#00BFFF"))
	popup_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_popup.add_child(popup_label)

## 處理巨型鮟鱇魚訊息
func handle_anglerfish_electric(payload: Dictionary) -> void:
	var phase = payload.get("phase", "")
	match phase:
		"anglerfish_appear":
			_on_appear(payload)
		"super_zap_start":
			_on_super_zap_start(payload)
		"super_zap_result":
			_on_super_zap_result(payload)
		"anglerfish_killed":
			_on_killed(payload)
		"anglerfish_leave":
			_on_leave(payload)
		_:
			if phase.begins_with("zap_"):
				_on_zap(payload)
			elif phase.begins_with("super_zap_"):
				_on_super_zap_target(payload)

func _on_appear(payload: Dictionary) -> void:
	_panel.visible = true
	_panel.modulate.a = 1.0
	_banner.visible = true
	_zap_counter.visible = true
	_pool_label.visible = true
	_zap_counter.text = "電擊: 0/8"
	_pool_label.text = "獎池: 0 💰"

	# 藍紫色雙閃光
	_flash_screen(Color("#4B0082"), 0.4)
	await get_tree().create_timer(0.45).timeout
	_flash_screen(Color("#9400D3"), 0.3)

	# 橫幅滑入動畫
	_banner.position.y = -30
	var tween = create_tween()
	tween.tween_property(_banner, "position:y", 15, 0.3)

func _on_zap(payload: Dictionary) -> void:
	var zap_index = payload.get("zap_index", 0)
	var target_x = payload.get("target_x", 400.0)
	var target_y = payload.get("target_y", 300.0)
	var is_kill = payload.get("is_kill", false)
	var is_treasure = payload.get("is_treasure", false)
	var treasure_mult = payload.get("treasure_mult", 0.0)
	var zap_reward = payload.get("zap_reward", 0)
	var is_empty = payload.get("is_empty", false)

	# 更新計數器
	_zap_counter.text = "電擊: %d/8" % zap_index

	if is_empty:
		return

	# 電弧效果（目標位置閃光）
	_spawn_zap_effect(target_x, target_y, is_treasure)

	if is_treasure:
		# 寶箱開箱：金色強閃光 + 大字彈跳
		_flash_screen(Color("#FFD700"), 0.5)
		_spawn_treasure_text(target_x, target_y, treasure_mult, zap_reward)
		# 更新獎池
		_pool_label.text = "獎池: %d 💰" % zap_reward
	elif is_kill and zap_reward > 0:
		# 普通擊破：獎勵浮動文字
		_spawn_reward_text(target_x, target_y, zap_reward)
		# 更新獎池（累積）
		var current_pool = int(_pool_label.text.split(":")[1].strip_edges().split(" ")[0])
		_pool_label.text = "獎池: %d 💰" % (current_pool + zap_reward)

func _on_super_zap_start(payload: Dictionary) -> void:
	var target_count = payload.get("target_count", 0)

	# 全螢幕藍白強閃光
	_flash_screen(Color("#FFFFFF"), 0.6)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color("#00BFFF"), 0.5)

	# 超級電擊橫幅
	_banner.text = "⚡⚡⚡ 超級電擊！全場 %d 個目標！" % target_count
	_banner.add_theme_color_override("font_color", Color("#FFFFFF"))

	# 橫幅震動
	var tween = create_tween().set_loops(3)
	tween.tween_property(_banner, "position:x", 5, 0.05)
	tween.tween_property(_banner, "position:x", -5, 0.05)
	tween.tween_property(_banner, "position:x", 0, 0.05)

func _on_super_zap_target(payload: Dictionary) -> void:
	var target_x = payload.get("target_x", 400.0)
	var target_y = payload.get("target_y", 300.0)
	var is_kill = payload.get("is_kill", false)
	var is_treasure = payload.get("is_treasure", false)
	var zap_reward = payload.get("zap_reward", 0)

	# 電弧效果
	_spawn_zap_effect(target_x, target_y, is_treasure)

	if is_kill and zap_reward > 0:
		_spawn_reward_text(target_x, target_y, zap_reward)

func _on_super_zap_result(payload: Dictionary) -> void:
	var super_kills = payload.get("super_kills", 0)
	var super_reward = payload.get("super_reward", 0)

	# 恢復橫幅
	_banner.text = "⚡ 巨型鮟鱇魚電擊！"
	_banner.add_theme_color_override("font_color", Color("#00BFFF"))

	# 超級電擊結果浮動文字（螢幕中央）
	var vp_size = get_viewport().get_visible_rect().size
	_spawn_reward_text(vp_size.x / 2, vp_size.y / 2, super_reward)

	if super_kills >= 5:
		_flash_screen(Color("#FFD700"), 0.4)

func _on_killed(payload: Dictionary) -> void:
	var zap_count = payload.get("zap_count", 0)
	var total_pool = payload.get("total_pool", 0)
	var pool_bonus = payload.get("pool_bonus", 0)
	var base_reward = payload.get("base_reward", 0)
	var total_reward = payload.get("total_reward", 0)
	var killer_name = payload.get("killer_name", "")

	# 金色強閃光
	_flash_screen(Color("#FFD700"), 0.5)

	# 顯示結算彈窗（右側滑入）
	_result_popup.visible = true
	_result_popup.modulate.a = 0.0
	var vp_size = get_viewport().get_visible_rect().size
	_result_popup.position.x = vp_size.x

	var label = _result_popup.get_node("ResultLabel")
	label.text = "⚡ 鮟鱇魚擊破！\n電擊: %d 次\n獎池: %d 金幣\n獎池加成: +%d\n基礎獎勵: %d\n🏆 總獎勵: %d" % [
		zap_count, total_pool, pool_bonus, base_reward, total_reward
	]

	var tween = create_tween()
	tween.tween_property(_result_popup, "modulate:a", 1.0, 0.3)
	tween.parallel().tween_property(_result_popup, "position:x", vp_size.x - 320, 0.3)

	# 依獎勵決定特效
	if total_reward >= 500:
		await get_tree().create_timer(0.3).timeout
		_flash_screen(Color("#FFD700"), 0.4)
	elif total_reward >= 200:
		await get_tree().create_timer(0.3).timeout
		_flash_screen(Color("#9400D3"), 0.3)

	# 5 秒後淡出
	await get_tree().create_timer(5.0).timeout
	_fade_out()

func _on_leave(payload: Dictionary) -> void:
	var zap_count = payload.get("zap_count", 0)
	# 鮟鱇魚離開，淡出 UI
	_fade_out()

## 電弧效果（目標位置）
func _spawn_zap_effect(x: float, y: float, is_treasure: bool) -> void:
	var effect = Control.new()
	effect.position = Vector2(x, y)
	effect.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_panel.add_child(effect)

	var color = Color("#FFD700") if is_treasure else Color("#00BFFF")
	var radius = 50.0 if is_treasure else 35.0

	# 電弧圓圈
	var circle = ColorRect.new()
	circle.size = Vector2(radius * 2, radius * 2)
	circle.position = Vector2(-radius, -radius)
	circle.color = Color(color.r, color.g, color.b, 0.6)
	circle.mouse_filter = Control.MOUSE_FILTER_IGNORE
	effect.add_child(circle)

	# ⚡ 符號
	var zap_label = Label.new()
	zap_label.text = "⚡"
	zap_label.add_theme_font_size_override("font_size", 28 if is_treasure else 20)
	zap_label.position = Vector2(-14, -18)
	zap_label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	effect.add_child(zap_label)

	var tween = create_tween()
	tween.tween_property(circle, "scale", Vector2(1.8, 1.8), 0.35)
	tween.parallel().tween_property(circle, "modulate:a", 0.0, 0.35)
	tween.parallel().tween_property(zap_label, "modulate:a", 0.0, 0.35)
	await tween.finished
	effect.queue_free()

## 寶箱開箱大字彈跳
func _spawn_treasure_text(x: float, y: float, mult: float, reward: int) -> void:
	var label = Label.new()
	label.text = "💰 寶箱開箱！×%.1f\n+%d 金幣" % [mult, reward]
	label.add_theme_font_size_override("font_size", 26)
	label.add_theme_color_override("font_color", Color("#FFD700"))
	label.position = Vector2(x - 80, y - 40)
	label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(label)

	# 彈跳動畫
	var tween = create_tween()
	tween.tween_property(label, "scale", Vector2(1.3, 1.3), 0.15)
	tween.tween_property(label, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_property(label, "position:y", y - 100, 0.8)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.8)
	await tween.finished
	label.queue_free()

## 獎勵浮動文字
func _spawn_reward_text(x: float, y: float, reward: int) -> void:
	var label = Label.new()
	label.text = "+%d" % reward
	label.add_theme_font_size_override("font_size", 20)
	label.add_theme_color_override("font_color", Color("#00BFFF"))
	label.position = Vector2(x - 25, y - 20)
	label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_panel.add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "position:y", y - 65, 0.8)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.8)
	await tween.finished
	label.queue_free()

func _fade_out() -> void:
	var tween = create_tween()
	tween.tween_property(_panel, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_banner, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_zap_counter, "modulate:a", 0.0, 0.5)
	tween.parallel().tween_property(_pool_label, "modulate:a", 0.0, 0.5)
	await tween.finished
	_panel.visible = false
	_panel.modulate.a = 1.0
	_result_popup.visible = false
	_banner.visible = false
	_banner.modulate.a = 1.0
	_zap_counter.visible = false
	_zap_counter.modulate.a = 1.0
	_pool_label.visible = false
	_pool_label.modulate.a = 1.0

func _flash_screen(color: Color, intensity: float) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, intensity)
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.28)
	await tween.finished
	flash.queue_free()
