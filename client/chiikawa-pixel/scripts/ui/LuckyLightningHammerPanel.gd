## LuckyLightningHammerPanel.gd — 幸運閃電錘魚 UI 面板（DAY-277）
## 閃電錘主題：#FFD700 金 + #87CEEB 天藍 + #FFA500 橙 + #FFFFFF 白
## 業界依據：Battle of Luck「Lucky Slammer」機制（2026）進化版
##
## 事件類型：
##   hammer_start     — 閃電錘觸發（全服，PlayerID/PlayerName/HammerCount）
##   hammer_hit       — 錘擊目標（全服，InstanceID/DefID/HPDamage/Killed/AccumMult）
##   hammer_end       — 錘擊結算（全服，PlayerID/PlayerName/HitCount/KillCount/FinalMult/TotalReward）
##   hammer_broadcast — 全服廣播橫幅（PlayerName/HitCount/KillCount/FinalMult）

extends CanvasLayer

const COLOR_GOLD      = Color(1.0,   0.843, 0.0)    # #FFD700 金
const COLOR_CYAN      = Color(0.529, 0.808, 0.922)  # #87CEEB 天藍
const COLOR_ORANGE    = Color(1.0,   0.647, 0.0)    # #FFA500 橙
const COLOR_WHITE     = Color(1.0,   1.0,   1.0)
const COLOR_RED       = Color(1.0,   0.2,   0.2)
const COLOR_ELECTRIC  = Color(0.8,   0.9,   1.0)    # 電光藍白

var _banner: Control = null
var _hit_counter: Control = null
var _hit_label: Label = null
var _kill_label: Label = null
var _hit_count: int = 0
var _kill_count: int = 0

func _ready() -> void:
	layer = 50  # 比 LuckyGoldenHurricane（49）高一層

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"hammer_start":
			_on_hammer_start(payload)
		"hammer_hit":
			_on_hammer_hit(payload)
		"hammer_end":
			_on_hammer_end(payload)
		"hammer_broadcast":
			_on_hammer_broadcast(payload)

# ── 閃電錘觸發 ────────────────────────────────────────────────────────────────

func _on_hammer_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var hammer_count: int = payload.get("hammer_count", 3)
	_hit_count = 0
	_kill_count = 0

	# 電光白三次強閃光
	_flash_screen(COLOR_ELECTRIC, 3, 0.55)

	# 頂部橫幅
	_show_banner(
		"⚡ 閃電錘！",
		"%s 觸發！瞬間錘擊 %d 個目標！30%% 機率直接擊破！" % [player_name, hammer_count],
		COLOR_GOLD
	)

	# 錘擊計數器（右上角）
	_show_hit_counter()

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"⚡ 閃電錘！%d 個目標！" % hammer_count,
		Vector2(vp_size / 2),
		COLOR_GOLD,
		44
	)

# ── 錘擊目標 ──────────────────────────────────────────────────────────────────

func _on_hammer_hit(payload: Dictionary) -> void:
	var hp_damage: int = payload.get("hp_damage", 0)
	var killed: bool = payload.get("killed", false)
	var accum_mult: float = payload.get("accum_mult", 1.0)
	_hit_count += 1
	if killed:
		_kill_count += 1

	# 更新計數器
	_update_hit_counter(accum_mult)

	# 根據是否擊破決定閃光強度
	if killed:
		# 金色強閃光（直接擊破）
		_flash_screen(COLOR_GOLD, 2, 0.4)
		var vp_size := get_viewport().get_visible_rect().size
		var pos := Vector2(
			vp_size.x * (0.3 + randf() * 0.4),
			vp_size.y * (0.3 + randf() * 0.4)
		)
		_spawn_float_text(
			"⚡ 直接擊破！×%.1f" % accum_mult,
			pos,
			COLOR_GOLD,
			24
		)
	else:
		# 天藍閃光（普通錘擊）
		_flash_screen(COLOR_CYAN, 1, 0.2)
		var vp_size := get_viewport().get_visible_rect().size
		var pos := Vector2(
			vp_size.x * (0.3 + randf() * 0.4),
			vp_size.y * (0.3 + randf() * 0.4)
		)
		_spawn_float_text(
			"⚡ HP -%d" % hp_damage,
			pos,
			COLOR_CYAN,
			16
		)

# ── 錘擊結算 ──────────────────────────────────────────────────────────────────

func _on_hammer_end(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var hit_count: int = payload.get("hit_count", 0)
	var kill_count: int = payload.get("kill_count", 0)
	var final_mult: float = payload.get("final_mult", 1.0)
	var total_reward: int = payload.get("total_reward", 0)

	# 清除橫幅和計數器
	_clear_banner()
	_clear_hit_counter()

	# 根據擊破數決定閃光強度
	if kill_count >= 3:
		# 金色三次強閃光（大爆發）
		_flash_screen(COLOR_GOLD, 3, 0.6)
		var vp_size := get_viewport().get_visible_rect().size
		_spawn_float_text(
			"⚡ 閃電錘大爆發！擊破 %d 個！×%.1f！" % [kill_count, final_mult],
			Vector2(vp_size / 2),
			COLOR_GOLD,
			40
		)
	elif kill_count >= 1:
		_flash_screen(COLOR_ORANGE, 2, 0.35)
	else:
		_flash_screen(COLOR_CYAN, 1, 0.2)

	# 結算彈窗
	_show_end_popup(player_name, hit_count, kill_count, final_mult, total_reward)

# ── 全服廣播橫幅 ──────────────────────────────────────────────────────────────

func _on_hammer_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var hit_count: int = payload.get("hit_count", 0)
	var kill_count: int = payload.get("kill_count", 0)
	var final_mult: float = payload.get("final_mult", 1.0)

	var color := COLOR_GOLD if kill_count >= 3 else (COLOR_ORANGE if kill_count >= 1 else COLOR_CYAN)
	_show_mini_banner(
		"⚡ %s 的閃電錘！錘擊 %d 個，擊破 %d 個，倍率 ×%.1f！" % [player_name, hit_count, kill_count, final_mult],
		color
	)

# ── 內部 UI 工具 ──────────────────────────────────────────────────────────────

func _flash_screen(color: Color, times: int, alpha: float = 0.4) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var flash := ColorRect.new()
	flash.size = vp_size
	flash.color = Color(color.r, color.g, color.b, 0.0)
	add_child(flash)

	var tween := create_tween()
	for i in range(times):
		tween.tween_property(flash, "color:a", alpha, 0.06)
		tween.tween_property(flash, "color:a", 0.0, 0.09)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)

func _show_banner(title: String, subtitle: String, color: Color) -> void:
	_clear_banner()
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(0, -80)
	panel.size = Vector2(vp_size.x, 72)
	panel.modulate = Color(0.1, 0.08, 0.0, 0.92)
	add_child(panel)
	_banner = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = title
	title_lbl.add_theme_color_override("font_color", color)
	title_lbl.add_theme_font_size_override("font_size", 22)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var sub_lbl := Label.new()
	sub_lbl.text = subtitle
	sub_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	sub_lbl.add_theme_font_size_override("font_size", 14)
	sub_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(sub_lbl)

	var tween := create_tween()
	tween.tween_property(panel, "position:y", 0.0, 0.22).set_ease(Tween.EASE_OUT)

func _clear_banner() -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()
	_banner = null

func _show_mini_banner(text: String, color: Color) -> void:
	var vp_size := get_viewport().get_visible_rect().size
	var lbl := Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 16)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.position = Vector2(0, 4)
	lbl.size = Vector2(vp_size.x, 28)
	add_child(lbl)

	var tween := create_tween()
	tween.tween_interval(3.5)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)

func _show_hit_counter() -> void:
	_clear_hit_counter()
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x - 130, 80)
	panel.size = Vector2(120, 60)
	panel.modulate = Color(0.1, 0.08, 0.0, 0.9)
	add_child(panel)
	_hit_counter = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	_hit_label = Label.new()
	_hit_label.text = "⚡ 錘擊 0"
	_hit_label.add_theme_color_override("font_color", COLOR_CYAN)
	_hit_label.add_theme_font_size_override("font_size", 15)
	_hit_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(_hit_label)

	_kill_label = Label.new()
	_kill_label.text = "💀 擊破 0"
	_kill_label.add_theme_color_override("font_color", COLOR_GOLD)
	_kill_label.add_theme_font_size_override("font_size", 15)
	_kill_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(_kill_label)

	# 脈衝動畫
	var tween := panel.create_tween().set_loops()
	tween.tween_property(panel, "modulate:a", 0.75, 0.3)
	tween.tween_property(panel, "modulate:a", 1.0, 0.3)

func _update_hit_counter(accum_mult: float) -> void:
	if is_instance_valid(_hit_label):
		_hit_label.text = "⚡ 錘擊 %d  ×%.1f" % [_hit_count, accum_mult]
	if is_instance_valid(_kill_label):
		_kill_label.text = "💀 擊破 %d" % _kill_count
		var kill_color := COLOR_GOLD if _kill_count >= 3 else (COLOR_ORANGE if _kill_count >= 1 else COLOR_CYAN)
		_kill_label.add_theme_color_override("font_color", kill_color)

func _clear_hit_counter() -> void:
	if is_instance_valid(_hit_counter):
		_hit_counter.queue_free()
	_hit_counter = null
	_hit_label = null
	_kill_label = null

func _show_end_popup(player_name: String, hit_count: int, kill_count: int, final_mult: float, total_reward: int) -> void:
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x + 10, vp_size.y * 0.3)
	panel.size = Vector2(270, 160)
	panel.modulate = Color(0.1, 0.08, 0.0, 0.95)
	add_child(panel)

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "⚡ 閃電錘結算"
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var trigger_lbl := Label.new()
	trigger_lbl.text = "觸發者：%s" % player_name
	trigger_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	trigger_lbl.add_theme_font_size_override("font_size", 13)
	trigger_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(trigger_lbl)

	var hit_lbl := Label.new()
	hit_lbl.text = "錘擊目標：%d 個" % hit_count
	hit_lbl.add_theme_color_override("font_color", COLOR_CYAN)
	hit_lbl.add_theme_font_size_override("font_size", 14)
	hit_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(hit_lbl)

	var kill_lbl := Label.new()
	kill_lbl.text = "直接擊破：%d 個" % kill_count
	var kill_color := COLOR_GOLD if kill_count >= 3 else (COLOR_ORANGE if kill_count >= 1 else COLOR_CYAN)
	kill_lbl.add_theme_color_override("font_color", kill_color)
	kill_lbl.add_theme_font_size_override("font_size", 16)
	kill_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(kill_lbl)

	var mult_lbl := Label.new()
	mult_lbl.text = "累積倍率：×%.1f" % final_mult
	mult_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	mult_lbl.add_theme_font_size_override("font_size", 18)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	if total_reward > 0:
		var reward_lbl := Label.new()
		reward_lbl.text = "擊破獎勵：+%d" % total_reward
		reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
		reward_lbl.add_theme_font_size_override("font_size", 14)
		reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		vbox.add_child(reward_lbl)

	# 右側滑入動畫
	var tween := create_tween()
	tween.tween_property(panel, "position:x", vp_size.x - 280.0, 0.3).set_ease(Tween.EASE_OUT)
	tween.tween_interval(4.5)
	tween.tween_property(panel, "position:x", vp_size.x + 10.0, 0.3).set_ease(Tween.EASE_IN)
	tween.tween_callback(func():
		if is_instance_valid(panel):
			panel.queue_free()
	)

func _spawn_float_text(text: String, pos: Vector2, color: Color, font_size: int = 28) -> void:
	var lbl := Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", font_size)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.position = pos - Vector2(300, font_size * 0.5)
	lbl.size = Vector2(600, font_size * 2)
	add_child(lbl)

	var tween := create_tween()
	tween.tween_property(lbl, "position:y", pos.y - 60, 0.8).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.8).set_delay(0.5)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)
