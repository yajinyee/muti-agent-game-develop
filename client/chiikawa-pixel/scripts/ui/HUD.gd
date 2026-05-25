## HUD.gd — 核心 HUD
## hud-core-agent 負責維護
extends CanvasLayer

@onready var coins_label: Label = $TopBar/CoinsLabel
@onready var bet_label: Label = $TopBar/BetLabel
@onready var char_label: Label = $TopBar/CharLabel
@onready var labor_bar: ProgressBar = $TopBar/LaborBar
@onready var labor_label: Label = $TopBar/LaborLabel
@onready var state_label: Label = $TopBar/StateLabel
@onready var auto_btn: Button = $BottomBar/AutoBtn
@onready var lock_btn: Button = $BottomBar/LockBtn
@onready var bet_minus_btn: Button = $BottomBar/BetMinusBtn
@onready var bet_plus_btn: Button = $BottomBar/BetPlusBtn
@onready var boss_btn: Button = $BottomBar/BossBtn
@onready var bonus_btn: Button = $BottomBar/BonusBtn

var _reward_popup: Label = null
var _disconnect_overlay: Control = null
var _boss_timer_panel: Control = null
var _boss_time_left: float = 0.0
var _boss_active: bool = false
var _last_labor: int = 0
# DAY-292 幸運特殊魚 UI
var _lucky_banner: Control = null
var _announce_queue: Array = []
var _announce_showing: bool = false

func _ready() -> void:
	GameManager.player_updated.connect(_on_player_updated)
	GameManager.game_state_changed.connect(_on_state_changed)
	GameManager.reward_received.connect(_on_reward_received)
	GameManager.boss_event.connect(_on_boss_event)
	GameManager.bonus_event.connect(_on_bonus_event)
	NetworkManager.connected.connect(_on_reconnected)
	NetworkManager.disconnected.connect(_on_disconnected)

	auto_btn.pressed.connect(func(): NetworkManager.send_auto_toggle())
	lock_btn.pressed.connect(func(): NetworkManager.send_lock(""))
	bet_minus_btn.pressed.connect(func(): NetworkManager.send_bet_change(max(1, GameManager.get_bet_level() - 1)))
	bet_plus_btn.pressed.connect(func(): NetworkManager.send_bet_change(min(10, GameManager.get_bet_level() + 1)))
	boss_btn.pressed.connect(func(): NetworkManager.send_trigger_boss())
	bonus_btn.pressed.connect(func(): NetworkManager.send_trigger_bonus())

	_create_reward_popup()
	_create_disconnect_overlay()
	_create_lucky_banner()
	_update_ui()

	# DAY-292 幸運特殊魚訊號連接
	GameManager.lucky_chain_lightning.connect(_on_lucky_chain_lightning)
	GameManager.lucky_crab_torpedo.connect(_on_lucky_crab_torpedo)
	GameManager.lucky_vortex.connect(_on_lucky_vortex)
	GameManager.lucky_golden_dragon.connect(_on_lucky_golden_dragon)
	GameManager.lucky_thunder_lobster.connect(_on_lucky_thunder_lobster)
	GameManager.announce.connect(_on_announce)
	# DAY-293 新增幸運特殊魚訊號連接
	GameManager.lucky_awakened_phoenix.connect(_on_lucky_awakened_phoenix)
	GameManager.lucky_shockwave_bomb.connect(_on_lucky_shockwave_bomb)
	# DAY-294 新增幸運特殊魚訊號連接
	GameManager.lucky_drill_torpedo.connect(_on_lucky_drill_torpedo)
	GameManager.lucky_time_freeze.connect(_on_lucky_time_freeze)
	GameManager.lucky_chain_explosion.connect(_on_lucky_chain_explosion)

func _process(delta: float) -> void:
	if _boss_active and _boss_time_left > 0:
		_boss_time_left -= delta
		_update_boss_timer()

func _on_player_updated(_data: Dictionary) -> void:
	_update_ui()

func _update_ui() -> void:
	coins_label.text = "💰 %d" % GameManager.get_coins()
	var lv = GameManager.get_bet_level()
	var cost = GameManager.get_bet_cost()
	bet_label.text = "BET LV%d (%d)" % [lv, cost]
	char_label.text = GameManager.get_character_name()
	char_label.modulate = GameManager.get_character_color()

	var labor = GameManager.get_labor_value()
	labor_bar.value = labor
	if labor >= 80:
		labor_label.text = "⚡%d/100" % labor
		labor_label.modulate = Color(1.0, 0.9, 0.2)
	else:
		labor_label.text = "%d/100" % labor
		labor_label.modulate = Color.WHITE
	if labor >= 100 and _last_labor < 100:
		ScreenShake.add_trauma(0.3)
	_last_labor = labor

	if GameManager.is_auto():
		auto_btn.text = "AUTO ON"
		auto_btn.modulate = Color(0.3, 1.0, 0.3)
	else:
		auto_btn.text = "AUTO"
		auto_btn.modulate = Color.WHITE

	var lock_id = GameManager.get_lock_target_id()
	if lock_id != "":
		lock_btn.text = "🔒 LOCK"
		lock_btn.modulate = Color(1.0, 0.8, 0.2)
	else:
		lock_btn.text = "🔓 LOCK"
		lock_btn.modulate = Color(0.7, 0.7, 0.7)

func _on_state_changed(new_state: String) -> void:
	state_label.text = new_state.to_upper().replace("_", " ")
	match new_state:
		"boss_battle":
			AudioManager.play_bgm(AudioManager.BGM.BOSS_ENTER)
		"boss_result":
			AudioManager.play_bgm(AudioManager.BGM.MAIN_GAME)
		"bonus_game":
			AudioManager.play_bgm(AudioManager.BGM.BONUS_GAME)
		"bonus_result", "normal_play":
			AudioManager.play_bgm(AudioManager.BGM.MAIN_GAME)

func _on_reward_received(reward: Dictionary) -> void:
	var amount = reward.get("amount", 0)
	var mult = reward.get("multiplier", 1.0)
	if amount <= 0:
		return
	_show_reward_popup(amount, mult)

func _show_reward_popup(amount: int, mult: float) -> void:
	if not is_instance_valid(_reward_popup):
		return
	var icon = "💰"
	if mult >= 100:
		icon = "🌟"
		_reward_popup.modulate = Color(1.0, 0.3, 0.1)
	elif mult >= 20:
		icon = "⭐"
		_reward_popup.modulate = Color(1.0, 0.85, 0.0)
	else:
		_reward_popup.modulate = Color.WHITE
	_reward_popup.text = "%s +%d  x%.0f" % [icon, amount, mult]
	_reward_popup.visible = true
	_reward_popup.modulate.a = 1.0
	_reward_popup.position.y = 350
	var tween = create_tween()
	tween.tween_property(_reward_popup, "position:y", 280.0, 0.7)
	tween.parallel().tween_property(_reward_popup, "modulate:a", 0.0, 0.7)
	tween.tween_callback(func(): if is_instance_valid(_reward_popup): _reward_popup.visible = false)

func _on_boss_event(event_data: Dictionary) -> void:
	match event_data.get("event", ""):
		"warning":
			AudioManager.play_sfx(AudioManager.SFX.BOSS_WARNING)
		"spawn":
			_start_boss_timer()
		"phase_change":
			AudioManager.play_bgm(AudioManager.BGM.BOSS_RAGE)
		"kill", "timeout":
			_stop_boss_timer()
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)

func _on_bonus_event(event_data: Dictionary) -> void:
	match event_data.get("event", ""):
		"start":
			AudioManager.play_sfx(AudioManager.SFX.BONUS_READY)

# ── BOSS 計時器 ───────────────────────────────────────────────

func _start_boss_timer() -> void:
	_boss_time_left = 60.0
	_boss_active = true
	if is_instance_valid(_boss_timer_panel):
		_boss_timer_panel.queue_free()

	var panel = Control.new()
	panel.name = "BossTimerPanel"
	panel.position = Vector2(900, 50)
	panel.size = Vector2(340, 75)
	add_child(panel)
	_boss_timer_panel = panel

	var bg = ColorRect.new()
	bg.size = panel.size
	bg.color = Color(0.1, 0.0, 0.0, 0.85)
	panel.add_child(bg)

	var title = Label.new()
	title.text = "⚔ BOSS BATTLE"
	title.position = Vector2(10, 5)
	title.add_theme_font_size_override("font_size", 15)
	title.modulate = Color(1.0, 0.3, 0.3)
	panel.add_child(title)

	var timer_lbl = Label.new()
	timer_lbl.name = "TimerLabel"
	timer_lbl.text = "60.0s"
	timer_lbl.position = Vector2(10, 26)
	timer_lbl.add_theme_font_size_override("font_size", 26)
	timer_lbl.modulate = Color(1.0, 0.9, 0.2)
	panel.add_child(timer_lbl)

	var mult_lbl = Label.new()
	mult_lbl.name = "MultLabel"
	mult_lbl.text = "500x"
	mult_lbl.position = Vector2(190, 26)
	mult_lbl.add_theme_font_size_override("font_size", 26)
	mult_lbl.modulate = Color(1.0, 0.5, 0.0)
	panel.add_child(mult_lbl)

	var hint = Label.new()
	hint.text = "Kill faster = higher reward!"
	hint.position = Vector2(10, 56)
	hint.add_theme_font_size_override("font_size", 11)
	hint.modulate = Color(0.8, 0.8, 0.8)
	panel.add_child(hint)

func _update_boss_timer() -> void:
	if not is_instance_valid(_boss_timer_panel):
		return
	var tl = _boss_timer_panel.get_node_or_null("TimerLabel")
	var ml = _boss_timer_panel.get_node_or_null("MultLabel")
	if is_instance_valid(tl):
		tl.text = "%.1fs" % max(0, _boss_time_left)
		tl.modulate = Color(1.0, 0.3, 0.3) if _boss_time_left <= 10 else Color(1.0, 0.9, 0.2)
	if is_instance_valid(ml):
		var m = 100
		if _boss_time_left > 50: m = 500
		elif _boss_time_left > 40: m = 400
		elif _boss_time_left > 30: m = 300
		elif _boss_time_left > 20: m = 200
		elif _boss_time_left > 10: m = 150
		ml.text = "%dx" % m

func _stop_boss_timer() -> void:
	_boss_active = false
	if is_instance_valid(_boss_timer_panel):
		var tween = create_tween()
		tween.tween_property(_boss_timer_panel, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func():
			if is_instance_valid(_boss_timer_panel):
				_boss_timer_panel.queue_free()
				_boss_timer_panel = null
		)

# ── 獎勵彈窗 ─────────────────────────────────────────────────

func _create_reward_popup() -> void:
	_reward_popup = Label.new()
	_reward_popup.visible = false
	_reward_popup.position = Vector2(540, 350)
	_reward_popup.add_theme_font_size_override("font_size", 24)
	_reward_popup.z_index = 50
	add_child(_reward_popup)

# ── 斷線提示 ─────────────────────────────────────────────────

func _create_disconnect_overlay() -> void:
	_disconnect_overlay = Control.new()
	_disconnect_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_disconnect_overlay.visible = false
	_disconnect_overlay.z_index = 100
	add_child(_disconnect_overlay)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0, 0, 0, 0.7)
	_disconnect_overlay.add_child(bg)

	var msg = Label.new()
	msg.text = "📡 DISCONNECTED\nReconnecting..."
	msg.position = Vector2(500, 330)
	msg.add_theme_font_size_override("font_size", 22)
	msg.modulate = Color(1.0, 0.4, 0.4)
	_disconnect_overlay.add_child(msg)

func _on_disconnected() -> void:
	if is_instance_valid(_disconnect_overlay):
		_disconnect_overlay.visible = true

func _on_reconnected() -> void:
	if is_instance_valid(_disconnect_overlay):
		_disconnect_overlay.visible = false

# ── DAY-292 幸運特殊魚 UI ─────────────────────────────────────

func _create_lucky_banner() -> void:
	_lucky_banner = Control.new()
	_lucky_banner.name = "LuckyBanner"
	_lucky_banner.position = Vector2(0, 120)
	_lucky_banner.size = Vector2(1280, 80)
	_lucky_banner.visible = false
	_lucky_banner.z_index = 60
	add_child(_lucky_banner)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0, 0, 0, 0.75)
	_lucky_banner.add_child(bg)

	var lbl = Label.new()
	lbl.name = "BannerLabel"
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 28)
	lbl.modulate = Color(1.0, 0.9, 0.2)
	_lucky_banner.add_child(lbl)

func _show_lucky_banner(text: String, color: Color, duration: float = 2.5) -> void:
	if not is_instance_valid(_lucky_banner):
		return
	var lbl = _lucky_banner.get_node_or_null("BannerLabel")
	if is_instance_valid(lbl):
		lbl.text = text
		lbl.modulate = color
	_lucky_banner.visible = true
	_lucky_banner.modulate.a = 1.0
	var tween = create_tween()
	tween.tween_interval(duration - 0.5)
	tween.tween_property(_lucky_banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_lucky_banner):
			_lucky_banner.visible = false
	)

func _on_lucky_chain_lightning(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_banner("⚡ %s 觸發連鎖閃電！" % name, Color(0.0, 0.9, 1.0))
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"chain_hit":
			var chain = data.get("chain_count", 0)
			var mult = data.get("multiplier", 1.0)
			_show_lucky_banner("⚡ 連鎖 %d！×%.1f" % [chain, mult], Color(0.0, 0.9, 1.0), 1.0)
		"settle":
			var reward = data.get("total_reward", 0)
			if reward > 0:
				_show_reward_popup(reward, data.get("multiplier", 1.0))

func _on_lucky_crab_torpedo(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_banner("🦀 %s 發射螃蟹魚雷！" % name, Color(1.0, 0.4, 0.1))
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.35)
		"explosion":
			var no = data.get("explosion_no", 1)
			_show_lucky_banner("💥 魚雷爆炸 %d/3！" % no, Color(1.0, 0.6, 0.2), 0.8)
			ScreenShake.add_trauma(0.5)
		"settle":
			var reward = data.get("total_reward", 0)
			if reward > 0:
				_show_reward_popup(reward, 3.0)

func _on_lucky_vortex(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_banner("🌀 %s 召喚渦旋海葵！" % name, Color(0.5, 0.2, 0.8))
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"pull":
			var tl = data.get("time_left", 0.0)
			_show_lucky_banner("🌀 渦旋中... %.0fs" % tl, Color(0.7, 0.4, 1.0), 0.9)
		"end":
			_show_lucky_banner("🌀 渦旋爆炸！全場 HP -20%！", Color(0.8, 0.5, 1.0))
			ScreenShake.add_trauma(0.6)

func _on_lucky_golden_dragon(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_banner("🐉 %s 觸發黃金龍魚輪盤！" % name, Color(1.0, 0.85, 0.0))
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"spin":
			var inner = data.get("inner_mult", 1.0)
			var outer = data.get("outer_mult", 1.0)
			var final_m = data.get("final_mult", 1.0)
			_show_lucky_banner("🐉 內環 ×%.0f × 外環 ×%.0f = ×%.0f！" % [inner, outer, final_m], Color(1.0, 0.85, 0.0), 3.0)
		"result":
			var reward = data.get("reward", 0)
			var final_m = data.get("final_mult", 1.0)
			if reward > 0:
				_show_reward_popup(reward, final_m)
			if final_m >= 100:
				ScreenShake.add_trauma(0.8)

func _on_lucky_thunder_lobster(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_banner("🦞⚡ %s 觸發雷霆龍蝦！15 秒免費射擊！" % name, Color(1.0, 0.3, 0.0))
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"auto_fire":
			var tl = data.get("time_left", 0.0)
			var kills = data.get("kill_count", 0)
			_show_lucky_banner("🦞⚡ 雷霆模式 %.0fs | 擊破 %d 條" % [tl, kills], Color(1.0, 0.5, 0.2), 0.8)
		"end":
			var reward = data.get("total_reward", 0)
			var kills = data.get("kill_count", 0)
			_show_lucky_banner("🦞 雷霆結束！擊破 %d 條，獎勵 %d！" % [kills, reward], Color(1.0, 0.7, 0.3))
			if reward > 0:
				_show_reward_popup(reward, float(kills))

func _on_announce(data: Dictionary) -> void:
	var msg = data.get("message", "")
	var priority = data.get("priority", "normal")
	var color_str = data.get("color", "#FFFFFF")
	var color = Color.WHITE
	# 解析 hex 顏色
	if color_str.begins_with("#") and color_str.length() == 7:
		var r = color_str.substr(1, 2).hex_to_int() / 255.0
		var g = color_str.substr(3, 2).hex_to_int() / 255.0
		var b = color_str.substr(5, 2).hex_to_int() / 255.0
		color = Color(r, g, b)

	var duration = 2.0
	match priority:
		"high": duration = 3.0
		"critical": duration = 4.0

	_announce_queue.append({"msg": msg, "color": color, "duration": duration})
	if not _announce_showing:
		_process_announce_queue()

func _process_announce_queue() -> void:
	if _announce_queue.is_empty():
		_announce_showing = false
		return
	_announce_showing = true
	var item = _announce_queue.pop_front()
	_show_lucky_banner(item["msg"], item["color"], item["duration"])
	var tween = create_tween()
	tween.tween_interval(item["duration"] + 0.2)
	tween.tween_callback(_process_announce_queue)

# ── DAY-293 新增幸運特殊魚事件處理 ───────────────────────────

func _on_lucky_awakened_phoenix(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"awaken_start":
			_show_lucky_banner("🔥 %s 觸發覺醒鳳凰！下 5 次攻擊 Power Up！" % name, Color(1.0, 0.42, 0.21))
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.35)
		"power_up":
			var mult = data.get("power_up_mult", 6.0)
			var shots = data.get("shots_left", 0)
			_show_lucky_banner("🔥 Power Up ×%.0f！剩餘 %d 次" % [mult, shots], Color(1.0, 0.6, 0.2), 1.0)
		"perfect_awaken":
			_show_lucky_banner("🔥✨ 完美覺醒！%s 全服 ×2.0 加成 8 秒！" % name, Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.6)
		"perfect_end":
			_show_lucky_banner("🔥 完美覺醒加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"awaken_end":
			var reward = data.get("total_reward", 0)
			var hits = data.get("hit_count", 0)
			_show_lucky_banner("🔥 覺醒結束！命中 %d 次，獎勵 %d！" % [hits, reward], Color(1.0, 0.7, 0.3))
			if reward > 0:
				_show_reward_popup(reward, float(hits))

func _on_lucky_shockwave_bomb(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"shockwave_start":
			_show_lucky_banner("💥 %s 觸發全場震盪！全場 HP -35%！" % name, Color(1.0, 0.27, 0.0))
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.7)
		"shockwave_hit":
			var hits = data.get("hit_count", 0)
			var reward = data.get("total_reward", 0)
			_show_lucky_banner("💥 震盪命中 %d 個目標！獎勵 %d！" % [hits, reward], Color(1.0, 0.5, 0.2))
			if reward > 0:
				_show_reward_popup(reward, float(hits) * 0.5)
		"super_shockwave":
			_show_lucky_banner("💥🌊 超級震盪！%s 全服 ×1.8 加成 6 秒！" % name, Color(1.0, 0.42, 0.21), 3.5)
			ScreenShake.add_trauma(0.8)
		"super_end":
			_show_lucky_banner("💥 超級震盪加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"power_end":
			_show_lucky_banner("💥 %s 的震盪強化結束" % name, Color(0.6, 0.6, 0.6), 1.5)

# ── DAY-294 新增幸運特殊魚事件處理 ───────────────────────────

func _on_lucky_drill_torpedo(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_lucky_banner("🚀 %s 發射鑽頭魚雷！穿透最多 5 個目標！" % name, Color(1.0, 0.55, 0.15))
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"penetrate":
			var cnt = data.get("penetrate_cnt", 0)
			var mult = data.get("accum_mult", 1.0)
			_show_lucky_banner("🚀 穿透 %d 個！累積 ×%.1f" % [cnt, mult], Color(1.0, 0.7, 0.3), 0.8)
		"explode":
			_show_lucky_banner("💥 魚雷終點爆炸！AOE 傷害！", Color(1.0, 0.4, 0.1))
			ScreenShake.add_trauma(0.55)
		"perfect":
			_show_lucky_banner("🚀💥 完美穿透！%s 全服 ×2.2 加成 6 秒！" % name, Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.65)
		"perfect_end":
			_show_lucky_banner("🚀 完美穿透加成結束", Color(0.7, 0.7, 0.7), 1.5)

func _on_lucky_time_freeze(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"freeze_start":
			_show_lucky_banner("❄️ %s 觸發時間凍結！全場凍結 8 秒！傷害 ×1.8！" % name, Color(0.4, 0.85, 1.0))
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.3)
		"freeze_end":
			_show_lucky_banner("❄️💥 冰裂爆炸！全場 HP -25%！", Color(0.6, 0.9, 1.0))
			ScreenShake.add_trauma(0.5)
		"perfect_freeze":
			var kills = data.get("kill_count", 0)
			_show_lucky_banner("❄️✨ 完美凍結！%s 擊破 %d 條！全服 ×2.0 加成 5 秒！" % [name, kills], Color(0.0, 0.9, 1.0), 3.5)
			ScreenShake.add_trauma(0.6)
		"perfect_end":
			_show_lucky_banner("❄️ 完美凍結加成結束", Color(0.7, 0.7, 0.7), 1.5)

func _on_lucky_chain_explosion(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"chain_start":
			_show_lucky_banner("💥 %s 觸發連鎖爆炸！12 秒連鎖模式！" % name, Color(0.9, 0.2, 0.15))
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"chain_explode":
			var cnt = data.get("chain_count", 0)
			var mult = data.get("accum_mult", 1.0)
			_show_lucky_banner("💥 連鎖 %d！累積 ×%.1f" % [cnt, mult], Color(1.0, 0.5, 0.2), 0.7)
			ScreenShake.add_trauma(0.25)
		"chain_burst":
			_show_lucky_banner("💥🔥 連鎖爆發！%s 全服 ×2.5 加成 6 秒！" % name, Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.7)
		"burst_end":
			_show_lucky_banner("💥 連鎖爆發加成結束", Color(0.7, 0.7, 0.7), 1.5)
		"chain_end":
			var cnt = data.get("chain_count", 0)
			var reward = data.get("total_reward", 0)
			_show_lucky_banner("💥 連鎖結束！%d 次連鎖，獎勵 %d！" % [cnt, reward], Color(1.0, 0.6, 0.3))
			if reward > 0:
				_show_reward_popup(reward, float(cnt) * 0.5)
