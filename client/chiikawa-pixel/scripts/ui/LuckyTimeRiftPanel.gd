## LuckyTimeRiftPanel.gd — 幸運時間裂縫魚 UI 面板（DAY-278）
## 時間裂縫主題：#9B59B6 紫 + #FFD700 金 + #3498DB 藍 + #E8DAEF 淡紫白
## 業界原創「時間裂縫+最高倍率重現+裂縫複製體」機制
##
## 事件類型：
##   rift_start          — 時間裂縫觸發（個人，HasBestKill/ReplayName/ReplayMult/RiftMult/ImmediateReward/CloneInstanceID）
##   rift_broadcast      — 全服廣播橫幅（PlayerName/ReplayName/ReplayMult/RiftMult/ImmediateReward/CloneInstanceID）
##   rift_clone_kill     — 裂縫複製體被擊破（個人，PlayerName/ReplayName/CloneKillMult/CloneReward）
##   rift_clone_broadcast — 裂縫複製體擊破全服廣播（PlayerName/ReplayName/CloneReward）

extends CanvasLayer

const COLOR_PURPLE     = Color(0.608, 0.349, 0.714)  # #9B59B6 紫
const COLOR_GOLD       = Color(1.0,   0.843, 0.0)    # #FFD700 金
const COLOR_BLUE       = Color(0.204, 0.596, 0.859)  # #3498DB 藍
const COLOR_LIGHT_PURPLE = Color(0.91, 0.855, 0.937) # #E8DAEF 淡紫白
const COLOR_WHITE      = Color(1.0,   1.0,   1.0)
const COLOR_ORANGE     = Color(1.0,   0.647, 0.0)    # #FFA500 橙

var _banner: Control = null
var _rift_indicator: Control = null
var _rift_mult_label: Label = null

func _ready() -> void:
	layer = 51  # 比 LuckyLightningHammer（50）高一層

## handle 由 GameManager 呼叫
func handle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"rift_start":
			_on_rift_start(payload)
		"rift_broadcast":
			_on_rift_broadcast(payload)
		"rift_clone_kill":
			_on_rift_clone_kill(payload)
		"rift_clone_broadcast":
			_on_rift_clone_broadcast(payload)

# ── 時間裂縫觸發（個人）────────────────────────────────────────────────────────

func _on_rift_start(payload: Dictionary) -> void:
	var has_best_kill: bool = payload.get("has_best_kill", false)
	var replay_name: String = payload.get("replay_name", "???")
	var replay_mult: int = payload.get("replay_mult", 2)
	var rift_mult: float = payload.get("rift_mult", 2.5)
	var immediate_reward: int = payload.get("immediate_reward", 0)
	var clone_instance_id: String = payload.get("clone_instance_id", "")
	var clone_kill_mult: float = payload.get("clone_kill_mult", 3.0)

	# 紫色三次強閃光
	_flash_screen(COLOR_PURPLE, 3, 0.6)

	# 頂部橫幅
	if has_best_kill:
		_show_banner(
			"🌀 時間裂縫！",
			"最高倍率 %s（×%d）重現！即時獎勵 +%d！裂縫複製體出現！" % [replay_name, replay_mult, immediate_reward],
			COLOR_PURPLE
		)
	else:
		_show_banner(
			"🌀 時間裂縫！",
			"保底獎勵 +%d！裂縫複製體出現！" % immediate_reward,
			COLOR_BLUE
		)

	# 裂縫指示器（右上角）
	_show_rift_indicator(replay_name, replay_mult, rift_mult, clone_kill_mult)

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	if has_best_kill:
		_spawn_float_text(
			"🌀 時間裂縫！×%d 重現！+%d！" % [replay_mult, immediate_reward],
			Vector2(vp_size / 2),
			COLOR_PURPLE,
			40
		)
	else:
		_spawn_float_text(
			"🌀 時間裂縫！保底 +%d！" % immediate_reward,
			Vector2(vp_size / 2),
			COLOR_BLUE,
			36
		)

	# 即時獎勵浮動文字
	_spawn_float_text(
		"+%d 籌碼！" % immediate_reward,
		Vector2(vp_size.x * 0.5, vp_size.y * 0.6),
		COLOR_GOLD,
		28
	)

	# 裂縫複製體提示
	if clone_instance_id != "":
		_spawn_float_text(
			"🌀 裂縫複製體出現！擊破得 ×%.1f！" % clone_kill_mult,
			Vector2(vp_size.x * 0.5, vp_size.y * 0.7),
			COLOR_LIGHT_PURPLE,
			18
		)

	# 5 秒後清除橫幅和指示器
	var timer := get_tree().create_timer(5.0)
	timer.timeout.connect(func():
		_clear_banner()
		_clear_rift_indicator()
	)

# ── 全服廣播橫幅 ──────────────────────────────────────────────────────────────

func _on_rift_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var replay_name: String = payload.get("replay_name", "???")
	var replay_mult: int = payload.get("replay_mult", 2)
	var rift_mult: float = payload.get("rift_mult", 2.5)
	var immediate_reward: int = payload.get("immediate_reward", 0)

	var color := COLOR_PURPLE if rift_mult >= 2.0 else COLOR_BLUE
	_show_mini_banner(
		"🌀 %s 觸發時間裂縫！%s（×%d）重現！即時獎勵 +%d！裂縫複製體出現！" % [player_name, replay_name, replay_mult, immediate_reward],
		color
	)

# ── 裂縫複製體被擊破（個人）──────────────────────────────────────────────────

func _on_rift_clone_kill(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var replay_name: String = payload.get("replay_name", "???")
	var clone_kill_mult: float = payload.get("clone_kill_mult", 3.0)
	var clone_reward: int = payload.get("clone_reward", 0)

	# 金色三次強閃光
	_flash_screen(COLOR_GOLD, 3, 0.55)

	# 清除指示器
	_clear_rift_indicator()

	# 浮動大字
	var vp_size := get_viewport().get_visible_rect().size
	_spawn_float_text(
		"🌀 裂縫複製體擊破！×%.1f 大獎！" % clone_kill_mult,
		Vector2(vp_size / 2),
		COLOR_GOLD,
		40
	)

	# 結算彈窗
	_show_clone_kill_popup(player_name, replay_name, clone_kill_mult, clone_reward)

# ── 裂縫複製體擊破全服廣播 ────────────────────────────────────────────────────

func _on_rift_clone_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "???")
	var replay_name: String = payload.get("replay_name", "???")
	var clone_reward: int = payload.get("clone_reward", 0)

	_show_mini_banner(
		"🌀 %s 擊破裂縫複製體 %s！獲得 +%d 籌碼！" % [player_name, replay_name, clone_reward],
		COLOR_GOLD
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
		tween.tween_property(flash, "color:a", alpha, 0.07)
		tween.tween_property(flash, "color:a", 0.0, 0.10)
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
	panel.modulate = Color(0.08, 0.04, 0.12, 0.92)
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

func _show_rift_indicator(replay_name: String, replay_mult: int, rift_mult: float, clone_kill_mult: float) -> void:
	_clear_rift_indicator()
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x - 145, 80)
	panel.size = Vector2(135, 80)
	panel.modulate = Color(0.08, 0.04, 0.12, 0.92)
	add_child(panel)
	_rift_indicator = panel

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "🌀 時間裂縫"
	title_lbl.add_theme_color_override("font_color", COLOR_PURPLE)
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	_rift_mult_label = Label.new()
	_rift_mult_label.text = "重現 ×%d  ×%.1f" % [replay_mult, rift_mult]
	_rift_mult_label.add_theme_color_override("font_color", COLOR_GOLD)
	_rift_mult_label.add_theme_font_size_override("font_size", 14)
	_rift_mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(_rift_mult_label)

	var clone_lbl := Label.new()
	clone_lbl.text = "裂縫複製體 ×%.1f" % clone_kill_mult
	clone_lbl.add_theme_color_override("font_color", COLOR_LIGHT_PURPLE)
	clone_lbl.add_theme_font_size_override("font_size", 12)
	clone_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(clone_lbl)

	# 脈衝動畫
	var tween := panel.create_tween().set_loops()
	tween.tween_property(panel, "modulate:a", 0.7, 0.4)
	tween.tween_property(panel, "modulate:a", 1.0, 0.4)

func _clear_rift_indicator() -> void:
	if is_instance_valid(_rift_indicator):
		_rift_indicator.queue_free()
	_rift_indicator = null
	_rift_mult_label = null

func _show_clone_kill_popup(player_name: String, replay_name: String, clone_kill_mult: float, clone_reward: int) -> void:
	var vp_size := get_viewport().get_visible_rect().size

	var panel := PanelContainer.new()
	panel.position = Vector2(vp_size.x + 10, vp_size.y * 0.3)
	panel.size = Vector2(270, 150)
	panel.modulate = Color(0.08, 0.04, 0.12, 0.95)
	add_child(panel)

	var vbox := VBoxContainer.new()
	vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	panel.add_child(vbox)

	var title_lbl := Label.new()
	title_lbl.text = "🌀 裂縫複製體擊破！"
	title_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var name_lbl := Label.new()
	name_lbl.text = "目標：%s" % replay_name
	name_lbl.add_theme_color_override("font_color", COLOR_LIGHT_PURPLE)
	name_lbl.add_theme_font_size_override("font_size", 14)
	name_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(name_lbl)

	var mult_lbl := Label.new()
	mult_lbl.text = "裂縫倍率：×%.1f" % clone_kill_mult
	mult_lbl.add_theme_color_override("font_color", COLOR_PURPLE)
	mult_lbl.add_theme_font_size_override("font_size", 18)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	var reward_lbl := Label.new()
	reward_lbl.text = "裂縫大獎：+%d 籌碼！" % clone_reward
	reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	reward_lbl.add_theme_font_size_override("font_size", 20)
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
