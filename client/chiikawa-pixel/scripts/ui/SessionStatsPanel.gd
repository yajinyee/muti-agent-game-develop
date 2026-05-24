## SessionStatsPanel.gd
## Session 蝯梯??Ｘ嚗AY-046嚗? HUD.gd ?? DAY-053嚗?
## 憿舐內?砍?蝯梯?嚗?畾箸??擃???蜇??OSS ?捏?onus 甈⊥?楊?嗥?
## 瘥?60 蝘???箔?甈∴?variable reinforcement嚗??拙振???圈脣漲嚗?

extends Control

# ??HUD.gd ?典遣蝡?閮剖?
var pixel_font: Font = null

const SESSION_AUTO_POPUP_INTERVAL = 60.0  # 瘥?60 蝘???箔?甈?

var _visible_flag: bool = false
var _auto_popup_timer: float = 0.0
var _start_coins: int = 0

# ?砍?蝯梯??豢?嚗 GameManager 閮??湔嚗?
var _kills: int = 0
var _max_combo: int = 0
var _total_reward: int = 0
var _boss_kills: int = 0
var _bonus_count: int = 0

## ??????HUD.gd ?澆嚗?
func setup(font: Font) -> void:
	pixel_font = font
	GameManager.reward_received.connect(_on_session_reward)
	GameManager.combo_event.connect(_on_session_combo)
	GameManager.boss_event.connect(_on_session_boss)
	GameManager.bonus_event.connect(_on_session_bonus_event)
	_start_coins = GameManager.get_coins()

## 撱箇?????砍?????TopBar嚗?
func create_button(top_bar: Control) -> void:
	if not is_instance_valid(top_bar):
		return
	var btn = Button.new()
	btn.name = "SessionStatsButton"
	btn.text = "?? ?砍?"
	btn.position = Vector2(840, 4)
	btn.size = Vector2(80, 32)
	btn.add_theme_font_size_override("font_size", 12)
	if is_instance_valid(pixel_font):
		btn.add_theme_font_override("font", pixel_font)
	btn.pressed.connect(toggle)
	top_bar.add_child(btn)

## 瘥??湔嚗???箄???
func _process(delta: float) -> void:
	_auto_popup_timer += delta
	if _auto_popup_timer >= SESSION_AUTO_POPUP_INTERVAL:
		_auto_popup_timer = 0.0
		var state = GameManager.current_state
		if state == "normal_play" or state == "special_target_event":
			show_popup()

## ??憿舐內
func toggle() -> void:
	_visible_flag = not _visible_flag
	visible = _visible_flag
	if _visible_flag:
		_refresh()

## 敶憿舐內嚗? 蝘??芸??嗉絲嚗?
func show_popup() -> void:
	_visible_flag = true
	visible = true
	_refresh()
	var t = get_tree().create_timer(3.0)
	t.timeout.connect(func():
		if is_instance_valid(self):
			visible = false
			_visible_flag = false
	)

## 閮???
func _on_session_reward(reward: Dictionary) -> void:
	_total_reward += reward.get("amount", 0)

func _on_session_combo(combo_data: Dictionary) -> void:
	var count = combo_data.get("combo_count", 0)
	if count > _max_combo:
		_max_combo = count

func _on_session_boss(boss_data: Dictionary) -> void:
	var event = boss_data.get("event", "")
	if event == "kill":
		_boss_kills += 1

func _on_session_bonus_event(bonus_data: Dictionary) -> void:
	var event = bonus_data.get("event", "")
	if event == "end":
		_bonus_count += 1

## ?瑟蝯梯??豢?憿舐內
func _refresh() -> void:
	var player_data = GameManager.player_data
	var kills = player_data.get("kill_count", _kills)
	var reward = player_data.get("session_score", _total_reward)
	var current_coins = GameManager.get_coins()
	var net_profit = current_coins - _start_coins

	var rows = {
		"KillsRow":  str(kills),
		"ComboRow":  ("?%d" % _max_combo) if _max_combo > 0 else "??,"
		"RewardRow": ("??%d" % reward) if reward > 0 else "0",
		"BossRow":   str(_boss_kills),
		"BonusRow":  str(_bonus_count),
		"ProfitRow": ("%+d" % net_profit),
	}
	for row_name in rows:
		var row = get_node_or_null(row_name)
		if is_instance_valid(row):
			var val_lbl = row.get_node_or_null("Value")
			if is_instance_valid(val_lbl):
				val_lbl.text = rows[row_name]
				# 憿擃漁
				if row_name == "ComboRow" and _max_combo >= 5:
					val_lbl.add_theme_color_override("font_color", Color(1.0, 0.7, 0.1))
				elif row_name == "BossRow" and _boss_kills > 0:
					val_lbl.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))
				elif row_name == "BonusRow" and _bonus_count > 0:
					val_lbl.add_theme_color_override("font_color", Color(0.4, 1.0, 0.5))
				elif row_name == "ProfitRow":
					if net_profit > 0:
						val_lbl.add_theme_color_override("font_color", Color(0.3, 1.0, 0.4))
					elif net_profit < 0:
						val_lbl.add_theme_color_override("font_color", Color(1.0, 0.4, 0.4))
					else:
						val_lbl.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))

## 撱箇??Ｘ UI
func _ready() -> void:
	name = "SessionStatsPanel"
	position = Vector2(1050, 50)
	size = Vector2(220, 200)
	z_index = 120
	visible = false
	_build_panel_ui()

func _build_panel_ui() -> void:
	# ?
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.02, 0.05, 0.15, 0.92)
	add_child(bg)

	# ??嚗??莎?
	for border_data in [
		[Vector2(0, 0), Vector2(220, 2)],
		[Vector2(0, 198), Vector2(220, 2)],
		[Vector2(0, 0), Vector2(2, 200)],
		[Vector2(218, 0), Vector2(2, 200)],
	]:
		var border = ColorRect.new()
		border.position = border_data[0]
		border.size = border_data[1]
		border.color = Color(0.90, 0.75, 0.20, 0.80)
		add_child(border)

	# 璅?
	var title = Label.new()
	title.name = "Title"
	title.text = "?? ?砍?蝯梯?"
	title.position = Vector2(10, 8)
	title.size = Vector2(200, 24)
	title.add_theme_font_size_override("font_size", 14)
	title.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if is_instance_valid(pixel_font):
		title.add_theme_font_override("font", pixel_font)
	add_child(title)

	# ??蝺?
	var sep = ColorRect.new()
	sep.position = Vector2(8, 34)
	sep.size = Vector2(204, 1)
	sep.color = Color(0.90, 0.75, 0.20, 0.40)
	add_child(sep)

	# 蝯梯?銵?6銵?
	var stats_data = [
		["KillsRow",   "?? ?捏",     "0"],
		["ComboRow",   "? ?擃??",  "0"],
		["RewardRow",  "?? 蝮賜???,    "0"],"
		["BossRow",    "? BOSS ?捏", "0"],
		["BonusRow",   "? Bonus 甈⊥","0"],
		["ProfitRow",  "?? 瘛冽??,    "0"],"
	]
	for i in range(stats_data.size()):
		var row = Control.new()
		row.name = stats_data[i][0]
		row.position = Vector2(8, 40 + i * 26)
		row.size = Vector2(204, 24)
		add_child(row)

		var key_lbl = Label.new()
		key_lbl.name = "Key"
		key_lbl.text = stats_data[i][1]
		key_lbl.position = Vector2(0, 3)
		key_lbl.size = Vector2(130, 20)
		key_lbl.add_theme_font_size_override("font_size", 11)
		key_lbl.add_theme_color_override("font_color", Color(0.7, 0.85, 1.0))
		if is_instance_valid(pixel_font):
			key_lbl.add_theme_font_override("font", pixel_font)
		row.add_child(key_lbl)

		var val_lbl = Label.new()
		val_lbl.name = "Value"
		val_lbl.text = stats_data[i][2]
		val_lbl.position = Vector2(130, 3)
		val_lbl.size = Vector2(74, 20)
		val_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
		val_lbl.add_theme_font_size_override("font_size", 12)
		val_lbl.add_theme_color_override("font_color", Color(1.0, 0.95, 0.5))
		if is_instance_valid(pixel_font):
			val_lbl.add_theme_font_override("font", pixel_font)
		row.add_child(val_lbl)

	# ESC ?內嚗??剁?
	var esc_hint = Label.new()
	esc_hint.name = "EscHint"
	esc_hint.text = "[ESC] ??"
	esc_hint.position = Vector2(8, 182)
	esc_hint.size = Vector2(204, 14)
	esc_hint.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	esc_hint.add_theme_font_size_override("font_size", 9)
	esc_hint.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5, 0.6))
	if is_instance_valid(pixel_font):
		esc_hint.add_theme_font_override("font", pixel_font)
	add_child(esc_hint)
