## HUD.gd
## 銝駁???UI嚗??潭 11蝡?
## DAY-053嚗???JackpotPanel / MissionPanel / SessionStatsPanel ?箇蝡??

extends CanvasLayer

const PixelTheme = preload("res://scripts/ui/PixelTheme.gd")
const JackpotPanelScript = preload("res://scripts/ui/JackpotPanel.gd")
const MissionPanelScript = preload("res://scripts/ui/MissionPanel.gd")
const SessionStatsPanelScript = preload("res://scripts/ui/SessionStatsPanel.gd")
const LeaderboardPanelScript = preload("res://scripts/ui/LeaderboardPanel.gd")
const TournamentPanelScript = preload("res://scripts/ui/TournamentPanel.gd")
const WeaponPanelScript = preload("res://scripts/ui/WeaponPanel.gd")
const TitlePanelScript = preload("res://scripts/ui/TitlePanel.gd")
const SkinPanelScript = preload("res://scripts/ui/SkinPanel.gd")
const SeasonPanelScript = preload("res://scripts/ui/SeasonPanel.gd")
const FriendPanelScript = preload("res://scripts/ui/FriendPanel.gd")
const GuildPanelScript = preload("res://scripts/ui/GuildPanel.gd")
const GuildWarPanelScript = preload("res://scripts/ui/GuildWarPanel.gd")
const DailyBossPanelScript = preload("res://scripts/ui/DailyBossPanel.gd")
const VIPPanelScript = preload("res://scripts/ui/VIPPanel.gd")
const EventPanelScript = preload("res://scripts/ui/EventPanel.gd")
const CodexPanelScript = preload("res://scripts/ui/CodexPanel.gd")
const StreakPanelScript = preload("res://scripts/ui/StreakPanel.gd")
const ReferralPanelScript = preload("res://scripts/ui/ReferralPanel.gd")
const WheelPanelScript = preload("res://scripts/ui/WheelPanel.gd")

@onready var coins_label: Label = $TopBar/CoinsLabel
@onready var bet_label: Label = $TopBar/BetLabel
@onready var character_label: Label = $TopBar/CharacterLabel
@onready var labor_bar: ProgressBar = $TopBar/LaborBar
@onready var labor_label: Label = $TopBar/LaborLabel
@onready var auto_button: Button = $BottomBar/AutoButton
@onready var lock_button: Button = $BottomBar/LockButton
@onready var bet_minus_button: Button = $BottomBar/BetMinusButton
@onready var bet_plus_button: Button = $BottomBar/BetPlusButton
@onready var boss_button: Button = $BottomBar/BossButton
@onready var bonus_button: Button = $BottomBar/BonusButton
@onready var reward_popup: Label = $RewardPopup
@onready var state_label: Label = $StateLabel
@onready var warning_overlay: Control = $WarningOverlay
@onready var bonus_overlay: Control = $BonusOverlay

var _reward_popup_base_y: float = 0.0
var _lock_active: bool = false

# BOSS 閮??剁?閬??28.3嚗＊蝷箏擗???撠???嚗?
var _boss_time_left: float = 0.0
var _boss_active: bool = false
var _boss_timer_node: Control = null

# ??摮?嚗??潭蝢?閬?嚗?
var _pixel_font: Font = null
const PIXEL_FONT_PATH = "res://assets/fonts/pixel8.fnt"

func _ready() -> void:
	# 憟??憸冽 Theme嚗?????? UI ?????渡???憸冽嚗?
	var pixel_theme = PixelTheme.create()
	# 憟??TopBar ??BottomBar ????蝭暺?
	var top_bar = get_node_or_null("TopBar")
	var bottom_bar = get_node_or_null("BottomBar")
	if is_instance_valid(top_bar):
		top_bar.theme = pixel_theme
		# TopBar ?嚗楛瘚瑁???嚗?
		var top_bg = ColorRect.new()
		top_bg.name = "PixelBG"
		top_bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		top_bg.color = Color(0.03, 0.06, 0.18, 0.88)
		top_bg.z_index = -1
		top_bar.add_child(top_bg)
		top_bar.move_child(top_bg, 0)
		# TopBar 摨??蝺??嚗?
		var top_line = ColorRect.new()
		top_line.name = "BottomLine"
		top_line.size = Vector2(1280, 2)
		top_line.position = Vector2(0, 38)
		top_line.color = Color(0.90, 0.75, 0.20, 0.60)
		top_bar.add_child(top_line)
		# 觀戰者計數標籤（DAY-055）
		var spectator_lbl = Label.new()
		spectator_lbl.name = "SpectatorCountLabel"
		spectator_lbl.text = ""
		spectator_lbl.position = Vector2(1180, 8)
		spectator_lbl.size = Vector2(90, 24)
		spectator_lbl.add_theme_font_size_override("font_size", 12)
		spectator_lbl.modulate = Color(0.7, 0.85, 1.0, 0.8)
		spectator_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
		top_bar.add_child(spectator_lbl)
	if is_instance_valid(bottom_bar):
		bottom_bar.theme = pixel_theme
		# BottomBar ?嚗楛瘚瑁???嚗?
		var bot_bg = ColorRect.new()
		bot_bg.name = "PixelBG"
		bot_bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		bot_bg.color = Color(0.03, 0.06, 0.18, 0.88)
		bot_bg.z_index = -1
		bottom_bar.add_child(bot_bg)
		bottom_bar.move_child(bot_bg, 0)
		# BottomBar ???蝺??嚗?
		var bot_line = ColorRect.new()
		bot_line.name = "TopLine"
		bot_line.size = Vector2(1280, 2)
		bot_line.position = Vector2(0, 0)
		bot_line.color = Color(0.90, 0.75, 0.20, 0.60)
		bottom_bar.add_child(bot_line)

	# 頛??摮?
	if ResourceLoader.exists(PIXEL_FONT_PATH):
		_pixel_font = load(PIXEL_FONT_PATH)
		_apply_pixel_font()

	GameManager.player_updated.connect(_on_player_updated)
	GameManager.game_state_changed.connect(_on_game_state_changed)
	GameManager.reward_received.connect(_on_reward_received)
	GameManager.boss_event.connect(_on_boss_event)
	GameManager.bonus_event.connect(_on_bonus_event)
	GameManager.leaderboard_updated.connect(_on_leaderboard_updated)
	GameManager.achievement_unlocked.connect(_on_achievement_unlocked)
	GameManager.combo_event.connect(_on_combo_event)  # ???鈭辣嚗AY-022嚗?
	# jackpot_updated / jackpot_won 撌脩宏??JackpotPanel.gd嚗AY-053嚗?
	GameManager.spectator_joined.connect(_on_spectator_joined)  # 觀戰者加入（DAY-054d）
	GameManager.spectator_left.connect(_on_spectator_left)      # 觀戰者離開（DAY-055）
	GameManager.daily_bonus_received.connect(_on_daily_bonus_received)  # 每日登入獎勵（DAY-065）

	# ?瑞?/???蝷?NetworkManager.disconnected.connect(_on_disconnected)
	NetworkManager.connected.connect(_on_reconnected)

	auto_button.pressed.connect(_on_auto_pressed)
	lock_button.pressed.connect(_on_lock_pressed)
	bet_minus_button.pressed.connect(_on_bet_minus)
	bet_plus_button.pressed.connect(_on_bet_plus)
	boss_button.pressed.connect(NetworkManager.send_trigger_boss)
	bonus_button.pressed.connect(NetworkManager.send_trigger_bonus)

	# UI ??暺??單?嚗??潭 audio-map.json嚗i.click = weed_pull.wav嚗?
	for btn in [auto_button, lock_button, bet_minus_button, bet_plus_button, boss_button, bonus_button]:
		if is_instance_valid(btn):
			btn.pressed.connect(func(): AudioManager.play_sfx(AudioManager.SFX.WEED_PULL))

	reward_popup.visible = false
	_reward_popup_base_y = reward_popup.position.y

	# WarningLabel ??憸冽嚗之摮??脯敶梧?
	var warning_label = get_node_or_null("WarningOverlay/WarningLabel")
	if is_instance_valid(warning_label):
		warning_label.add_theme_font_size_override("font_size", 72)
		warning_label.add_theme_color_override("font_color", Color(1.0, 0.15, 0.15))
		warning_label.add_theme_color_override("font_shadow_color", Color(0.5, 0.0, 0.0, 0.8))
		warning_label.add_theme_constant_override("shadow_offset_x", 3)
		warning_label.add_theme_constant_override("shadow_offset_y", 3)
		if is_instance_valid(_pixel_font):
			warning_label.add_theme_font_override("font", _pixel_font)

	# StateLabel ??憸冽嚗銝???＊蝷綽?
	var state_lbl = get_node_or_null("StateLabel")
	if is_instance_valid(state_lbl):
		state_lbl.add_theme_font_size_override("font_size", 11)
		state_lbl.add_theme_color_override("font_color", Color(0.6, 0.8, 1.0, 0.7))

	_update_ui()
	_create_disconnect_overlay()
	_create_leaderboard_panel()
	_create_achievement_queue()
	_create_lobby_overlay()  # 憭批輒 UI嚗AY-020嚗?
	_init_mission_panel()     # 瘥隞餃?蝟餌絞嚗AY-037嚗AY-053 ??嚗?
	_init_session_stats()     # Session Stats ?Ｘ嚗AY-046嚗AY-053 ??嚗?
	_init_jackpot_panel()     # Progressive Jackpot ?Ｘ嚗AY-048嚗AY-053 ??嚗?
	_init_tournament_panel()  # 週賽排名面板（DAY-066）
	_init_weapon_panel()      # 武器升級面板（DAY-067）
	_init_title_panel()       # 稱號面板（DAY-068）
	_init_skin_panel()        # 砲台外觀面板（DAY-071）
	_init_season_panel()      # 賽季通行證面板（DAY-072）
	_init_friend_panel()      # 好友系統面板（DAY-073）
	_init_guild_panel()       # 公會系統面板（DAY-074）
	_init_guild_war_panel()   # 公會戰面板（DAY-076）
	_init_daily_boss_panel()  # 每日 BOSS 挑戰面板（DAY-077）
	_init_vip_panel()         # VIP 等級面板（DAY-078）
	_init_event_panel()       # 限時活動面板（DAY-079）
	_init_codex_panel()       # 魚類圖鑑面板（DAY-081）
	_init_streak_panel()      # 連擊系統面板（DAY-083）
	_init_referral_panel()    # 推薦碼面板（DAY-082）
	_init_wheel_panel()       # 幸運轉盤面板（DAY-084）

## 憟??摮??唳???Label
func _apply_pixel_font() -> void:
	if not is_instance_valid(_pixel_font):
		return
	var labels = [coins_label, bet_label, character_label, labor_label, reward_popup, state_label]
	for label in labels:
		if is_instance_valid(label):
			label.add_theme_font_override("font", _pixel_font)
	# ??摮?
	var buttons = [auto_button, lock_button, bet_minus_button, bet_plus_button, boss_button, bonus_button]
	for btn in buttons:
		if is_instance_valid(btn):
			btn.add_theme_font_override("font", _pixel_font)

var _last_labor_value: int = 0  # 餈質馱銝活???潘??菜葫?遛閫貊

func _on_player_updated(_data: Dictionary) -> void:
	_update_ui()

func _update_ui() -> void:
	coins_label.text = "?? %d" % GameManager.get_coins()

	var lv = GameManager.get_bet_level()
	var cost = GameManager.get_bet_cost()
	bet_label.text = "BET LV%d  (%d/shot)" % [lv, cost]

	character_label.text = "??%s" % GameManager.get_character_name()
	character_label.modulate = GameManager.get_character_color()

	var labor = GameManager.get_labor_value()
	labor_bar.value = labor
	# ???潭餈遛????
	if labor >= 80:
		labor_label.text = "??%d/100" % labor
		labor_label.modulate = Color(1.0, 0.9, 0.2)
	else:
		labor_label.text = "? %d/100" % labor
		labor_label.modulate = Color.WHITE

	# ?菜葫???澆?? 100嚗孛?澆?蝝??
	if labor >= 100 and _last_labor_value < 100:
		var char_id = GameManager.player_data.get("character_id", "chiikawa")
		HitEffect.spawn_level_up(Vector2(640, 630), char_id)
		ScreenShake.add_trauma(0.3)
	_last_labor_value = labor

	# Auto ??
	if GameManager.is_auto():
		auto_button.modulate = Color(0.3, 1.0, 0.3)
		auto_button.text = "AUTO ON"
	else:
		auto_button.modulate = Color.WHITE
		auto_button.text = "AUTO"

	# Lock ??
	var lock_id = GameManager.get_lock_target_id()
	if lock_id != "":
		lock_button.modulate = Color(1.0, 0.8, 0.2)
		lock_button.text = "?? LOCK"
		_lock_active = true
	else:
		lock_button.modulate = Color(0.7, 0.7, 0.7)
		lock_button.text = "?? LOCK"
		_lock_active = false

func _on_game_state_changed(new_state: String) -> void:
	state_label.text = new_state.to_upper().replace("_", " ")

	match new_state:
		"boss_warning":
			_show_boss_warning()
		"boss_battle":
			warning_overlay.visible = false
		"bonus_game":
			bonus_overlay.visible = true
		"bonus_result", "normal_play", "boss_result":
			bonus_overlay.visible = false
			warning_overlay.visible = false

func _show_boss_warning() -> void:
	warning_overlay.visible = true
	var tween = create_tween().set_loops(8)
	tween.tween_property(warning_overlay, "modulate:a", 0.1, 0.18)
	tween.tween_property(warning_overlay, "modulate:a", 1.0, 0.18)

func _on_reward_received(reward: Dictionary) -> void:
	var amount = reward.get("amount", 0)
	var multiplier = reward.get("multiplier", 1.0)
	if amount <= 0:
		return
	_show_reward_popup(amount, multiplier)

func _show_reward_popup(amount: int, multiplier: float) -> void:
	# 靘?瘙箏?憿舐內?批捆
	var icon = "??"
	if multiplier >= 100:
		icon = "?"
		reward_popup.modulate = Color(1.0, 0.3, 0.1, 1.0)
	elif multiplier >= 20:
		icon = "?"
		reward_popup.modulate = Color(1.0, 0.85, 0.0, 1.0)
	elif multiplier >= 10:
		icon = "??"
		reward_popup.modulate = Color(1.0, 1.0, 0.4, 1.0)
	else:
		icon = "??"
		reward_popup.modulate = Color(1.0, 1.0, 1.0, 1.0)

	reward_popup.text = "%s +%d  ?%.0f" % [icon, amount, multiplier]
	reward_popup.position.y = _reward_popup_base_y

	reward_popup.visible = true

	var tween = create_tween()
	tween.tween_property(reward_popup, "position:y", _reward_popup_base_y - 70, 0.7)
	tween.parallel().tween_property(reward_popup, "modulate:a", 0.0, 0.7)
	tween.tween_callback(func():
		reward_popup.visible = false
		reward_popup.position.y = _reward_popup_base_y
	)

func _on_boss_event(event_data: Dictionary) -> void:
	match event_data.get("event", ""):
		"warning":
			AudioManager.stop_bgm_briefly()
			AudioManager.play_sfx(AudioManager.SFX.BOSS_WARNING)
			_show_boss_incoming_preview()  # BOSS 銵璇?閬賢???
		"spawn":
			AudioManager.play_bgm(AudioManager.BGM.BOSS_ENTER)
			_hide_boss_incoming_preview()  # ?梯??汗嚗＊蝷箸迤撘??
			_start_boss_timer()
		"phase_change":
			AudioManager.play_bgm(AudioManager.BGM.BOSS_RAGE)
		"kill", "timeout":
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			AudioManager.play_bgm(AudioManager.BGM.MAIN_GAME)
			_stop_boss_timer()
			# BOSS ?捏?嗥??寞?
			if event_data.get("event", "") == "kill":
				HitEffect.spawn_big_win(Vector2(640, 360), 100.0)
				ScreenShake.add_trauma(0.7)

# ?? BOSS 閮???UI嚗??潭 28.3嚗??????????????????????????????

func _start_boss_timer() -> void:
	_boss_time_left = 60.0
	_boss_active = true

	# 撱箇? BOSS 閮??券??
	if is_instance_valid(_boss_timer_node):
		_boss_timer_node.queue_free()

	var panel = Control.new()
	panel.name = "BossTimerPanel"
	panel.set_anchors_and_offsets_preset(Control.PRESET_TOP_RIGHT)
	panel.position = Vector2(900, 50)
	panel.size = Vector2(360, 80)
	add_child(panel)
	_boss_timer_node = panel

	# ?
	var bg = ColorRect.new()
	bg.size = Vector2(360, 80)
	bg.color = Color(0.1, 0.0, 0.0, 0.85)
	panel.add_child(bg)

	# BOSS 璅?
	var title = Label.new()
	title.name = "BossTitle"
	title.text = "??BOSS BATTLE"
	title.position = Vector2(10, 5)
	title.add_theme_font_size_override("font_size", 16)
	title.modulate = Color(1.0, 0.3, 0.3)
	if is_instance_valid(_pixel_font):
		title.add_theme_font_override("font", _pixel_font)
	panel.add_child(title)

	# ?拚???
	var timer_lbl = Label.new()
	timer_lbl.name = "BossTimeLabel"
	timer_lbl.text = "60.0s"
	timer_lbl.position = Vector2(10, 28)
	timer_lbl.add_theme_font_size_override("font_size", 28)
	timer_lbl.modulate = Color(1.0, 0.9, 0.2)
	if is_instance_valid(_pixel_font):
		timer_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(timer_lbl)

	# ???內
	var mult_lbl = Label.new()
	mult_lbl.name = "BossMultLabel"
	mult_lbl.text = "500x"
	mult_lbl.position = Vector2(200, 28)
	mult_lbl.add_theme_font_size_override("font_size", 28)
	mult_lbl.modulate = Color(1.0, 0.5, 0.0)
	if is_instance_valid(_pixel_font):
		mult_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(mult_lbl)

	# ??隤芣?
	var hint_lbl = Label.new()
	hint_lbl.name = "BossHintLabel"
	hint_lbl.text = "Kill faster = higher reward!"
	hint_lbl.position = Vector2(10, 60)
	hint_lbl.add_theme_font_size_override("font_size", 12)
	hint_lbl.modulate = Color(0.8, 0.8, 0.8)
	if is_instance_valid(_pixel_font):
		hint_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(hint_lbl)

func _stop_boss_timer() -> void:
	_boss_active = false
	if is_instance_valid(_boss_timer_node):
		var tween = create_tween()
		tween.tween_property(_boss_timer_node, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func():
			if is_instance_valid(_boss_timer_node):
				_boss_timer_node.queue_free()
				_boss_timer_node = null
		)

# ?? BOSS ?脣?汗 UI嚗郎??畾菟＊蝷?BOSS 銵璇? 0 憛急遛嚗??????????

var _boss_preview_node: Control = null

## 憿舐內 BOSS ?脣?汗嚗郎??畾?3 蝘?
## 銵璇? 0 蝺拇憛急遛??100%嚗???敺?
func _show_boss_incoming_preview() -> void:
	if is_instance_valid(_boss_preview_node):
		_boss_preview_node.queue_free()

	var panel = Control.new()
	panel.name = "BossIncomingPreview"
	panel.position = Vector2(320, 280)  # ?恍銝剖亢??
	panel.size = Vector2(640, 120)
	panel.z_index = 90
	panel.modulate.a = 0.0
	add_child(panel)
	_boss_preview_node = panel

	# ?
	var bg = ColorRect.new()
	bg.size = Vector2(640, 120)
	bg.color = Color(0.05, 0.0, 0.0, 0.88)
	panel.add_child(bg)

	# ?蝝??嚗???
	var top_bar = ColorRect.new()
	top_bar.name = "TopBar"
	top_bar.size = Vector2(640, 4)
	top_bar.color = Color(1.0, 0.1, 0.1, 1.0)
	panel.add_child(top_bar)

	# BOSS ?迂
	var name_lbl = Label.new()
	name_lbl.name = "BossNameLabel"
	name_lbl.text = "??酋摮?
	name_lbl.position = Vector2(20, 12)
	name_lbl.add_theme_font_size_override("font_size", 22)
	name_lbl.modulate = Color(1.0, 0.3, 0.3)
	if is_instance_valid(_pixel_font):
		name_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(name_lbl)

	# BOSS ?舀?憿?
	var sub_lbl = Label.new()
	sub_lbl.text = "BOSS  HP: 3000"
	sub_lbl.position = Vector2(20, 40)
	sub_lbl.add_theme_font_size_override("font_size", 13)
	sub_lbl.modulate = Color(0.8, 0.8, 0.8)
	if is_instance_valid(_pixel_font):
		sub_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(sub_lbl)

	# HP 璇???
	var hp_bg = ColorRect.new()
	hp_bg.size = Vector2(600, 20)
	hp_bg.position = Vector2(20, 65)
	hp_bg.color = Color(0.15, 0.0, 0.0, 1.0)
	panel.add_child(hp_bg)

	# HP 璇?敺?0 憛急遛嚗?
	var hp_bar = ColorRect.new()
	hp_bar.name = "BossHPBar"
	hp_bar.size = Vector2(0, 20)  # ??撖砍漲 0
	hp_bar.position = Vector2(20, 65)
	hp_bar.color = Color(0.9, 0.1, 0.1, 1.0)
	panel.add_child(hp_bar)

	# HP 璇????鈭桃?嚗?
	var hp_shine = ColorRect.new()
	hp_shine.name = "BossHPShine"
	hp_shine.size = Vector2(0, 4)
	hp_shine.position = Vector2(20, 65)
	hp_shine.color = Color(1.0, 0.5, 0.5, 0.6)
	panel.add_child(hp_shine)

	# ???內
	var mult_lbl = Label.new()
	mult_lbl.text = "MAX 500x"
	mult_lbl.position = Vector2(490, 12)
	mult_lbl.add_theme_font_size_override("font_size", 18)
	mult_lbl.modulate = Color(1.0, 0.6, 0.0)
	if is_instance_valid(_pixel_font):
		mult_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(mult_lbl)

	# ???
	var countdown_lbl = Label.new()
	countdown_lbl.name = "CountdownLabel"
	countdown_lbl.text = "3"
	countdown_lbl.position = Vector2(295, 88)
	countdown_lbl.add_theme_font_size_override("font_size", 20)
	countdown_lbl.modulate = Color(1.0, 0.9, 0.2)
	countdown_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	countdown_lbl.size = Vector2(50, 28)
	if is_instance_valid(_pixel_font):
		countdown_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(countdown_lbl)

	# ?摨?
	var tween = panel.create_tween()

	# 1. 瘛∪嚗?.2 蝘?
	tween.tween_property(panel, "modulate:a", 1.0, 0.2)

	# 2. HP 璇? 0 憛急遛嚗?.5 蝘?璅⊥ BOSS ?嚗?
	tween.parallel().tween_property(hp_bar, "size:x", 600.0, 2.5).set_ease(Tween.EASE_IN).set_trans(Tween.TRANS_QUAD)
	tween.parallel().tween_property(hp_shine, "size:x", 600.0, 2.5).set_ease(Tween.EASE_IN).set_trans(Tween.TRANS_QUAD)

	# 3. ? 3????
	tween.tween_callback(func():
		if is_instance_valid(countdown_lbl):
			countdown_lbl.text = "2"
			var t2 = countdown_lbl.create_tween()
			t2.tween_property(countdown_lbl, "scale", Vector2(1.5, 1.5), 0.1)
			t2.tween_property(countdown_lbl, "scale", Vector2(1.0, 1.0), 0.1)
	)
	tween.tween_interval(0.8)
	tween.tween_callback(func():
		if is_instance_valid(countdown_lbl):
			countdown_lbl.text = "1"
			countdown_lbl.modulate = Color(1.0, 0.4, 0.4)
			var t3 = countdown_lbl.create_tween()
			t3.tween_property(countdown_lbl, "scale", Vector2(1.8, 1.8), 0.1)
			t3.tween_property(countdown_lbl, "scale", Vector2(1.0, 1.0), 0.1)
	)
	tween.tween_interval(0.5)

	# 4. HP 璇????遛敺???3 甈∴?
	for _i in 3:
		tween.tween_property(hp_bar, "modulate", Color(2.0, 0.5, 0.5, 1.0), 0.06)
		tween.tween_property(hp_bar, "modulate", Color.WHITE, 0.06)

	# ??????嚗蝡?tween嚗?蝥?郎????
	var bar_tween = top_bar.create_tween().set_loops()
	bar_tween.tween_property(top_bar, "modulate:a", 0.3, 0.2)
	bar_tween.tween_property(top_bar, "modulate:a", 1.0, 0.2)

## ?梯? BOSS ?脣?汗嚗OSS 甇???箇??
func _hide_boss_incoming_preview() -> void:
	if not is_instance_valid(_boss_preview_node):
		return
	var tween = create_tween()
	tween.tween_property(_boss_preview_node, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(_boss_preview_node):
			_boss_preview_node.queue_free()
			_boss_preview_node = null
	)

func _get_boss_multiplier_text(time_left: float) -> String:
	if time_left <= 10:
		return "100x"
	elif time_left <= 20:
		return "150x"
	elif time_left <= 30:
		return "200x"
	elif time_left <= 40:
		return "300x"
	elif time_left <= 50:
		return "400x"
	else:
		return "500x"

func _get_boss_multiplier_color(time_left: float) -> Color:
	if time_left <= 10:
		return Color(0.6, 0.6, 0.6)  # ?堆??雿?嚗?
	elif time_left <= 20:
		return Color(0.4, 0.8, 1.0)  # ??
	elif time_left <= 30:
		return Color(0.4, 1.0, 0.4)  # 蝬?
	elif time_left <= 40:
		return Color(1.0, 0.9, 0.2)  # 暺?
	elif time_left <= 50:
		return Color(1.0, 0.5, 0.0)  # 璈?
	else:
		return Color(1.0, 0.2, 0.2)  # 蝝??擃?嚗?

func _process(delta: float) -> void:
	# Session Stats 自動彈出由 SessionStatsPanel._process 處理（DAY-053）

	if not _boss_active or not is_instance_valid(_boss_timer_node):
		return

	_boss_time_left = max(0.0, _boss_time_left - delta)

	var timer_lbl = _boss_timer_node.get_node_or_null("BossTimeLabel")
	var mult_lbl = _boss_timer_node.get_node_or_null("BossMultLabel")

	if timer_lbl:
		timer_lbl.text = "%.1fs" % _boss_time_left
		# ?敺?10 蝘??郎??
		if _boss_time_left <= 10.0:
			var flash = int(_boss_time_left * 4) % 2 == 0
			timer_lbl.modulate = Color.RED if flash else Color.WHITE
		else:
			timer_lbl.modulate = Color(1.0, 0.9, 0.2)

	if mult_lbl:
		var mult_text = _get_boss_multiplier_text(_boss_time_left)
		var mult_color = _get_boss_multiplier_color(_boss_time_left)
		mult_lbl.text = mult_text
		mult_lbl.modulate = mult_color

		# ??霈??憭扳?蝷?
		if mult_text != mult_lbl.get_meta("last_mult", ""):
			mult_lbl.set_meta("last_mult", mult_text)
			var tween = create_tween()
			tween.tween_property(mult_lbl, "scale", Vector2(1.4, 1.4), 0.1)
			tween.tween_property(mult_lbl, "scale", Vector2(1.0, 1.0), 0.1)

	# FPS 憿舐內嚗EBUG 璅∪?嚗?
	if OS.is_debug_build():
		_update_fps_display()

func _on_bonus_event(event_data: Dictionary) -> void:
	match event_data.get("event", ""):
		"ready":
			AudioManager.play_sfx(AudioManager.SFX.BONUS_READY)
		"start":
			AudioManager.play_bgm(AudioManager.BGM.BONUS_GAME)
		"tick":
			var tl = event_data.get("time_left", 0.0)
			var timer_lbl = bonus_overlay.get_node_or_null("TimerLabel")
			if timer_lbl:
				timer_lbl.text = "%.1f" % tl
		"end":
			_show_reward_popup(event_data.get("reward", 0), event_data.get("multiplier", 50.0))
			AudioManager.play_bgm(AudioManager.BGM.MAIN_GAME)

func _on_auto_pressed() -> void:
	NetworkManager.send_auto_toggle()

func _on_lock_pressed() -> void:
	# 閫???
	NetworkManager.send_lock("")

func _on_bet_minus() -> void:
	NetworkManager.send_bet_change(max(1, GameManager.get_bet_level() - 1))

func _on_bet_plus() -> void:
	NetworkManager.send_bet_change(min(10, GameManager.get_bet_level() + 1))

# ?? ?瑞?/??UI ??????????????????????????????????????????????

var _disconnect_overlay: Control = null
var _reconnect_dots_timer: float = 0.0
var _reconnect_dots: int = 0
var _is_disconnected: bool = false

func _create_disconnect_overlay() -> void:
	var overlay = Control.new()
	overlay.name = "DisconnectOverlay"
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.visible = false
	overlay.z_index = 100
	add_child(overlay)
	_disconnect_overlay = overlay

	# ??暺?
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.0, 0.0, 0.0, 0.75)
	overlay.add_child(bg)

	# ?瑞??內
	var icon_label = Label.new()
	icon_label.name = "DisconnectIcon"
	icon_label.text = "?"
	icon_label.position = Vector2(580, 290)
	icon_label.add_theme_font_size_override("font_size", 48)
	overlay.add_child(icon_label)

	# ?瑞???
	var msg_label = Label.new()
	msg_label.name = "DisconnectMsg"
	msg_label.text = "???銝剜"
	msg_label.position = Vector2(540, 355)
	msg_label.add_theme_font_size_override("font_size", 24)
	msg_label.modulate = Color(1.0, 0.4, 0.4)
	if is_instance_valid(_pixel_font):
		msg_label.add_theme_font_override("font", _pixel_font)
	overlay.add_child(msg_label)

	# ??葉??嚗葆??暺?嚗?
	var reconnect_label = Label.new()
	reconnect_label.name = "ReconnectLabel"
	reconnect_label.text = "????銝?.."
	reconnect_label.position = Vector2(520, 390)
	reconnect_label.add_theme_font_size_override("font_size", 18)
	reconnect_label.modulate = Color(0.8, 0.8, 0.8)
	if is_instance_valid(_pixel_font):
		reconnect_label.add_theme_font_override("font", _pixel_font)
	overlay.add_child(reconnect_label)

func _on_disconnected() -> void:
	_is_disconnected = true
	if is_instance_valid(_disconnect_overlay):
		_disconnect_overlay.visible = true
		# ???
		var tween = create_tween().set_loops()
		tween.tween_property(_disconnect_overlay, "modulate:a", 0.7, 0.5)
		tween.tween_property(_disconnect_overlay, "modulate:a", 1.0, 0.5)

func _on_reconnected() -> void:
	_is_disconnected = false
	if is_instance_valid(_disconnect_overlay):
		# 憿舐內?歇?????敺楚??
		var msg = _disconnect_overlay.get_node_or_null("DisconnectMsg")
		if msg:
			msg.text = "撌脤??圈?? ??
			msg.modulate = Color(0.3, 1.0, 0.3)
		var reconnect = _disconnect_overlay.get_node_or_null("ReconnectLabel")
		if reconnect:
			reconnect.visible = false

		var tween = create_tween()
		tween.tween_interval(1.0)
		tween.tween_property(_disconnect_overlay, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func():
			if is_instance_valid(_disconnect_overlay):
				_disconnect_overlay.visible = false
				_disconnect_overlay.modulate.a = 1.0
				# ?蔭??
				var m = _disconnect_overlay.get_node_or_null("DisconnectMsg")
				if m:
					m.text = "???銝剜"
					m.modulate = Color(1.0, 0.4, 0.4)
				var r = _disconnect_overlay.get_node_or_null("ReconnectLabel")
				if r:
					r.visible = true
		)

# ?? ??璁?UI ??????????????????????????????????????????????????
var _fps_label: Label = null
var _fps_update_timer: float = 0.0
var _perf_panel: Control = null  # 摰????Ｘ

func _update_fps_display() -> void:
	_fps_update_timer += get_process_delta_time()
	if _fps_update_timer < 0.5:
		return
	_fps_update_timer = 0.0

	# 擐活撱箇?摰??Ｘ
	if _perf_panel == null:
		_create_perf_panel()

	if not is_instance_valid(_perf_panel):
		return

	var fps = Engine.get_frames_per_second()
	var mem_mb = PerformanceMonitor.snapshot_memory_mb
	var draw_calls = PerformanceMonitor.snapshot_draw_calls
	var nodes = PerformanceMonitor.snapshot_nodes
	var quality = PerformanceMonitor.current_quality
	var quality_str = ["HIGH", "MED", "LOW"][quality]

	# FPS 銵?
	var fps_lbl = _perf_panel.get_node_or_null("FPSLine")
	if fps_lbl:
		fps_lbl.text = "FPS: %d  [%s]" % [fps, quality_str]
		if fps < 30:
			fps_lbl.modulate = Color(1.0, 0.3, 0.3, 0.95)
		elif fps < 50:
			fps_lbl.modulate = Color(1.0, 0.8, 0.2, 0.9)
		else:
			fps_lbl.modulate = Color(0.4, 1.0, 0.5, 0.85)

	# 閮擃?
	var mem_lbl = _perf_panel.get_node_or_null("MemLine")
	if mem_lbl:
		mem_lbl.text = "MEM: %.1f MB" % mem_mb
		if mem_mb > 200.0:
			mem_lbl.modulate = Color(1.0, 0.5, 0.2, 0.9)
		else:
			mem_lbl.modulate = Color(0.7, 0.9, 1.0, 0.8)

	# Draw Calls 銵?
	var dc_lbl = _perf_panel.get_node_or_null("DCLine")
	if dc_lbl:
		dc_lbl.text = "DC: %d  Nodes: %d" % [draw_calls, nodes]
		if draw_calls > 500:
			dc_lbl.modulate = Color(1.0, 0.6, 0.2, 0.9)
		else:
			dc_lbl.modulate = Color(0.7, 0.9, 1.0, 0.8)

	# Ping 銵?DAY-036嚗?
	var ping_lbl = _perf_panel.get_node_or_null("PingLine")
	if ping_lbl:
		var ping_ms = NetworkManager.get_ping_ms()
		if ping_ms < 0:
			ping_lbl.text = "PING: --"
			ping_lbl.modulate = Color(0.6, 0.6, 0.6, 0.7)
		else:
			ping_lbl.text = "PING: %d ms" % ping_ms
			if ping_ms > 200:
				ping_lbl.modulate = Color(1.0, 0.3, 0.3, 0.9)  # 蝝?擃辣??
			elif ping_ms > 100:
				ping_lbl.modulate = Color(1.0, 0.8, 0.2, 0.9)  # 暺?銝剖辣??
			else:
				ping_lbl.modulate = Color(0.4, 1.0, 0.5, 0.85) # 蝬?雿辣??

	# Pool 蝯梯?銵?DAY-041嚗ulletPool + TargetPool嚗?
	var pool_lbl = _perf_panel.get_node_or_null("PoolLine")
	if pool_lbl:
		var b_stats = PerformanceMonitor.get_bullet_pool_stats()
		var t_stats = TargetPool.get_stats()
		pool_lbl.text = "POOL B:%d/%d T:%d/%d" % [
			b_stats.get("active", 0), b_stats.get("total", 0),
			t_stats.get("active", 0), t_stats.get("total", 0)
		]

func _create_perf_panel() -> void:
	var panel = Control.new()
	panel.name = "PerfPanel"
	panel.position = Vector2(8, 670)
	panel.size = Vector2(220, 74)  # 擃漲敺?56 憓???74嚗???ping 銵?
	panel.z_index = 200
	add_child(panel)
	_perf_panel = panel

	# ???
	var bg = ColorRect.new()
	bg.size = Vector2(220, 74)
	bg.color = Color(0.0, 0.0, 0.0, 0.55)
	panel.add_child(bg)

	# 撌血蝬??嚗EBUG 璅?嚗?
	var side = ColorRect.new()
	side.size = Vector2(3, 74)
	side.color = Color(0.2, 1.0, 0.4, 0.8)
	panel.add_child(side)

	# FPS 銵?
	var fps_lbl = Label.new()
	fps_lbl.name = "FPSLine"
	fps_lbl.position = Vector2(8, 4)
	fps_lbl.size = Vector2(210, 16)
	fps_lbl.add_theme_font_size_override("font_size", 12)
	fps_lbl.text = "FPS: --  [HIGH]"
	fps_lbl.modulate = Color(0.4, 1.0, 0.5, 0.85)
	if is_instance_valid(_pixel_font):
		fps_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(fps_lbl)

	# 閮擃?
	var mem_lbl = Label.new()
	mem_lbl.name = "MemLine"
	mem_lbl.position = Vector2(8, 22)
	mem_lbl.size = Vector2(210, 16)
	mem_lbl.add_theme_font_size_override("font_size", 12)
	mem_lbl.text = "MEM: -- MB"
	mem_lbl.modulate = Color(0.7, 0.9, 1.0, 0.8)
	if is_instance_valid(_pixel_font):
		mem_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(mem_lbl)

	# Draw Calls 銵?
	var dc_lbl = Label.new()
	dc_lbl.name = "DCLine"
	dc_lbl.position = Vector2(8, 40)
	dc_lbl.size = Vector2(210, 16)
	dc_lbl.add_theme_font_size_override("font_size", 12)
	dc_lbl.text = "DC: --  Nodes: --"
	dc_lbl.modulate = Color(0.7, 0.9, 1.0, 0.8)
	if is_instance_valid(_pixel_font):
		dc_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(dc_lbl)

	# Ping 銵?DAY-036嚗?
	var ping_lbl = Label.new()
	ping_lbl.name = "PingLine"
	ping_lbl.position = Vector2(8, 58)
	ping_lbl.size = Vector2(210, 16)
	ping_lbl.add_theme_font_size_override("font_size", 12)
	ping_lbl.text = "PING: --"
	ping_lbl.modulate = Color(0.6, 0.6, 0.6, 0.7)
	if is_instance_valid(_pixel_font):
		ping_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(ping_lbl)

	# Pool 蝯梯?銵?DAY-041嚗ulletPool + TargetPool嚗?
	var pool_lbl = Label.new()
	pool_lbl.name = "PoolLine"
	pool_lbl.position = Vector2(8, 76)
	pool_lbl.size = Vector2(210, 16)
	pool_lbl.add_theme_font_size_override("font_size", 11)
	pool_lbl.text = "POOL: B?/? T?/?"
	pool_lbl.modulate = Color(0.6, 0.8, 0.6, 0.75)
	if is_instance_valid(_pixel_font):
		pool_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(pool_lbl)

	# ?Ｘ擃漲隤踵嚗???pool 銵?敺?74 憓???96嚗?
	bg.size.y = 96
	panel.size.y = 96
	var side_bar = panel.get_node_or_null("ColorRect")  # 撌血??
	if is_instance_valid(side_bar):
		side_bar.size.y = 96

# ?? ?停?蝟餌絞 ??????????????????????????????????????????????

var _achievement_queue: Array = []   # 敺＊蝷箇??停雿?
var _achievement_showing: bool = false
var _achievement_panel: Control = null

func _create_achievement_queue() -> void:
	# ?停??Ｘ嚗銝?嚗?憪??
	var panel = Control.new()
	panel.name = "AchievementPanel"
	panel.position = Vector2(900, 650)  # ?喃?閫?
	panel.size = Vector2(360, 80)
	panel.z_index = 50
	panel.visible = false
	add_child(panel)
	_achievement_panel = panel

	# ?嚗楛?脣???嚗??脤?獢?嚗?
	var bg = ColorRect.new()
	bg.name = "AchBG"
	bg.size = Vector2(360, 80)
	bg.color = Color(0.08, 0.06, 0.02, 0.92)
	panel.add_child(bg)

	# ????
	var top_bar = ColorRect.new()
	top_bar.size = Vector2(360, 4)
	top_bar.color = Color(1.0, 0.85, 0.1, 1.0)
	panel.add_child(top_bar)

	# ?停?內嚗之 emoji嚗?
	var icon_lbl = Label.new()
	icon_lbl.name = "AchIcon"
	icon_lbl.text = "??"
	icon_lbl.position = Vector2(8, 18)
	icon_lbl.add_theme_font_size_override("font_size", 36)
	panel.add_child(icon_lbl)

	# ??撠梯圾????憿?
	var title_lbl = Label.new()
	title_lbl.name = "AchTitle"
	title_lbl.text = "?停閫??嚗?
	title_lbl.position = Vector2(58, 8)
	title_lbl.add_theme_font_size_override("font_size", 11)
	title_lbl.modulate = Color(1.0, 0.85, 0.1)
	if is_instance_valid(_pixel_font):
		title_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(title_lbl)

	# ?停?迂
	var name_lbl = Label.new()
	name_lbl.name = "AchName"
	name_lbl.text = ""
	name_lbl.position = Vector2(58, 26)
	name_lbl.add_theme_font_size_override("font_size", 16)
	name_lbl.modulate = Color.WHITE
	if is_instance_valid(_pixel_font):
		name_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(name_lbl)

	# ?停?膩
	var desc_lbl = Label.new()
	desc_lbl.name = "AchDesc"
	desc_lbl.text = ""
	desc_lbl.position = Vector2(58, 50)
	desc_lbl.add_theme_font_size_override("font_size", 11)
	desc_lbl.modulate = Color(0.8, 0.8, 0.8)
	if is_instance_valid(_pixel_font):
		desc_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(desc_lbl)

func _on_achievement_unlocked(achievement_data: Dictionary) -> void:
	_achievement_queue.append(achievement_data)
	if not _achievement_showing:
		_show_next_achievement()

func _show_next_achievement() -> void:
	if _achievement_queue.is_empty() or not is_instance_valid(_achievement_panel):
		_achievement_showing = false
		return

	_achievement_showing = true
	var data = _achievement_queue.pop_front()

	# ?湔?Ｘ?批捆
	var icon_lbl = _achievement_panel.get_node_or_null("AchIcon")
	var name_lbl = _achievement_panel.get_node_or_null("AchName")
	var desc_lbl = _achievement_panel.get_node_or_null("AchDesc")

	if icon_lbl:
		icon_lbl.text = data.get("icon", "??")
	if name_lbl:
		name_lbl.text = data.get("name", "")
	if desc_lbl:
		desc_lbl.text = data.get("description", "")

	# 靘?撠梢??身摰椰?游蔗?脤?璇???
	var side_bar = _achievement_panel.get_node_or_null("AchSideBar")
	if not is_instance_valid(side_bar):
		side_bar = ColorRect.new()
		side_bar.name = "AchSideBar"
		side_bar.size = Vector2(4, 80)
		side_bar.position = Vector2(0, 0)
		_achievement_panel.add_child(side_bar)
	var ach_type = data.get("type", "normal")
	match ach_type:
		"boss":    side_bar.color = Color(1.0, 0.2, 0.2, 1.0)   # 蝝 ??BOSS ?賊?
		"bonus":   side_bar.color = Color(0.2, 0.8, 0.2, 1.0)   # 蝬 ??Bonus ?賊?
		"special": side_bar.color = Color(0.6, 0.2, 1.0, 1.0)   # 蝝怨 ???寞??停
		_:         side_bar.color = Color(1.0, 0.85, 0.1, 1.0)  # ? ??銝?祆?撠?

	# ?剜?單?嚗 bonus_ready ?單?嚗?
	AudioManager.play_sfx(AudioManager.SFX.BONUS_READY)

	# ?嚗??喳皛 ??敶歲蝮格 ???? 3 蝘???瘛∪皛粥
	_achievement_panel.modulate.a = 1.0
	_achievement_panel.scale = Vector2(1.0, 1.0)
	_achievement_panel.position.x = 1300.0  # ?恍憭??
	_achievement_panel.visible = true

	var tween = create_tween().set_parallel(false)
	# 皛嚗?.35 蝘?BACK 敶改?
	tween.tween_property(_achievement_panel, "position:x", 900.0, 0.35).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	# 敶歲蝮格嚗?.15 蝘??曉之 ???迤嚗?
	var scale_tween = create_tween().set_parallel(true)
	scale_tween.tween_property(_achievement_panel, "scale", Vector2(1.05, 1.05), 0.08).set_ease(Tween.EASE_OUT)
	scale_tween.chain().tween_property(_achievement_panel, "scale", Vector2(1.0, 1.0), 0.1).set_ease(Tween.EASE_IN_OUT)
	# ?? 3 蝘?
	tween.tween_interval(3.0)
	# 瘛∪ + 皛嚗?.3 蝘?
	tween.tween_property(_achievement_panel, "modulate:a", 0.0, 0.3).set_ease(Tween.EASE_IN)
	tween.tween_callback(func():
		if is_instance_valid(_achievement_panel):
			_achievement_panel.visible = false
			_achievement_panel.modulate.a = 1.0
			_achievement_panel.scale = Vector2(1.0, 1.0)
		# 憿舐內銝???撠梧??乩???蝛綽?
		_show_next_achievement()
	)

# ?? 憭批輒 UI嚗AY-020嚗??????????????????????????????????????????

var _lobby_overlay: Control = null
var _lobby_manager: Control = null

## 撱箇?憭批輒 overlay嚗?憪???舐???????恬?
func _create_lobby_overlay() -> void:
	# 撱箇??刻撟?overlay 摰孵
	var overlay = Control.new()
	overlay.name = "LobbyOverlay"
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.visible = false
	overlay.z_index = 150
	add_child(overlay)
	_lobby_overlay = overlay

	# 撱箇? LobbyManager UI
	var lobby_script = load("res://scripts/ui/LobbyManager.gd")
	if lobby_script:
		_lobby_manager = lobby_script.new()
		_lobby_manager.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		overlay.add_child(_lobby_manager)
		# ???輸??豢?閮?
		_lobby_manager.room_selected.connect(_on_lobby_room_selected)

	# ??TopBar ????????
	var top_bar = get_node_or_null("TopBar")
	if is_instance_valid(top_bar):
		var switch_btn = Button.new()
		switch_btn.name = "SwitchRoomBtn"
		switch_btn.text = "??"
		switch_btn.position = Vector2(1240, 6)
		switch_btn.size = Vector2(32, 28)
		switch_btn.add_theme_font_size_override("font_size", 16)
		switch_btn.pressed.connect(_show_lobby)
		switch_btn.tooltip_text = "???輸?"
		if is_instance_valid(_pixel_font):
			switch_btn.add_theme_font_override("font", _pixel_font)
		top_bar.add_child(switch_btn)

		# ?迂閮剖???嚗AY-021嚗?
		var name_btn = Button.new()
		name_btn.name = "SetNameBtn"
		name_btn.text = "??
		name_btn.position = Vector2(1204, 6)
		name_btn.size = Vector2(32, 28)
		name_btn.add_theme_font_size_override("font_size", 16)
		name_btn.pressed.connect(show_name_dialog)
		name_btn.tooltip_text = "閮剖??迂"
		if is_instance_valid(_pixel_font):
			name_btn.add_theme_font_override("font", _pixel_font)
		top_bar.add_child(name_btn)

## 憿舐內憭批輒
func _show_lobby() -> void:
	if not is_instance_valid(_lobby_overlay):
		return
	_lobby_overlay.visible = true
	_lobby_overlay.modulate.a = 0.0
	var tween = create_tween()
	tween.tween_property(_lobby_overlay, "modulate:a", 1.0, 0.3)
	# 閫貊?輸??”?瑟
	if is_instance_valid(_lobby_manager) and _lobby_manager.has_method("show_lobby"):
		_lobby_manager.show_lobby()

## 憭批輒?豢??輸?敺??矽
func _on_lobby_room_selected(room_id: String) -> void:
	print("[HUD] Room selected: ", room_id)
	# 瘛∪憭批輒
	var tween = create_tween()
	tween.tween_property(_lobby_overlay, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(_lobby_overlay):
			_lobby_overlay.visible = false
	)

	# 閫?唳芋撘?憿舐內???唬葉??蝐歹?DAY-024嚗?
	if NetworkManager.is_spectator():
		_show_spectator_badge()
		var notify_lbl = Label.new()
		notify_lbl.text = "?? 閫??%s 銝?.." % room_id
		notify_lbl.position = Vector2(440, 360)
		notify_lbl.size = Vector2(400, 40)
		notify_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		notify_lbl.add_theme_font_size_override("font_size", 18)
		notify_lbl.modulate = Color(0.5, 0.8, 1.0)
		if is_instance_valid(_pixel_font):
			notify_lbl.add_theme_font_override("font", _pixel_font)
		add_child(notify_lbl)
		var t2 = create_tween()
		t2.tween_interval(2.0)
		t2.tween_property(notify_lbl, "modulate:a", 0.0, 0.5)
		t2.tween_callback(func():
			if is_instance_valid(notify_lbl):
				notify_lbl.queue_free()
		)
		return

	# 銝?砍??交??憿舐內???輸??內
	var notify_lbl = Label.new()
	notify_lbl.text = "????%s..." % room_id
	notify_lbl.position = Vector2(440, 360)
	notify_lbl.size = Vector2(400, 40)
	notify_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	notify_lbl.add_theme_font_size_override("font_size", 18)
	notify_lbl.modulate = Color(0.4, 1.0, 0.5)
	if is_instance_valid(_pixel_font):
		notify_lbl.add_theme_font_override("font", _pixel_font)
	add_child(notify_lbl)
	var t2 = create_tween()
	t2.tween_interval(1.5)
	t2.tween_property(notify_lbl, "modulate:a", 0.0, 0.5)
	t2.tween_callback(func():
		if is_instance_valid(notify_lbl):
			notify_lbl.queue_free()
	)

## 憿舐內閫?唳?蝐歹?DAY-024嚗??喃?閫??脯??閫?唬葉??蝐?
func _show_spectator_badge() -> void:
	# ?踹???撱箇?
	if get_node_or_null("SpectatorBadge") != null:
		return
	var badge = Label.new()
	badge.name = "SpectatorBadge"
	badge.text = "?? 閫?唬葉"
	badge.position = Vector2(1050, 8)
	badge.size = Vector2(180, 24)
	badge.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	badge.add_theme_font_size_override("font_size", 14)
	badge.modulate = Color(0.5, 0.8, 1.0)
	if is_instance_valid(_pixel_font):
		badge.add_theme_font_override("font", _pixel_font)
	# ? TopBar
	var top_bar = get_node_or_null("TopBar")
	if is_instance_valid(top_bar):
		top_bar.add_child(badge)
	else:
		add_child(badge)

# ?? ?拙振?迂閮剖?嚗AY-021嚗??????????????????????????????????????

var _name_dialog: Control = null

## 憿舐內?迂閮剖?撠店獢?
func show_name_dialog() -> void:
	if is_instance_valid(_name_dialog):
		_name_dialog.queue_free()

	var dialog = Control.new()
	dialog.name = "NameDialog"
	dialog.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	dialog.z_index = 200
	add_child(dialog)
	_name_dialog = dialog

	# ???嚗?????
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.0, 0.0, 0.0, 0.6)
	bg.gui_input.connect(func(event):
		if event is InputEventMouseButton and event.pressed:
			if is_instance_valid(_name_dialog):
				_name_dialog.queue_free()
				_name_dialog = null
	)
	dialog.add_child(bg)

	# 撠店獢?選??恍銝剖亢嚗?
	var panel = Control.new()
	panel.name = "Panel"
	panel.position = Vector2(390, 300)
	panel.size = Vector2(500, 160)
	dialog.add_child(panel)

	# ?Ｘ?
	var panel_bg = ColorRect.new()
	panel_bg.size = Vector2(500, 160)
	panel_bg.color = Color(0.05, 0.08, 0.2, 0.97)
	panel.add_child(panel_bg)

	# ????
	var top_line = ColorRect.new()
	top_line.size = Vector2(500, 3)
	top_line.color = Color(0.9, 0.75, 0.2, 0.9)
	panel.add_child(top_line)

	# 璅?
	var title = Label.new()
	title.text = "??閮剖?憿舐內?迂"
	title.position = Vector2(16, 12)
	title.add_theme_font_size_override("font_size", 18)
	title.modulate = Color(1.0, 0.9, 0.3)
	if is_instance_valid(_pixel_font):
		title.add_theme_font_override("font", _pixel_font)
	panel.add_child(title)

	# 隤芣???
	var hint = Label.new()
	hint.text = "1-16 摮?嚗＊蝷箏??璁?"
	hint.position = Vector2(16, 40)
	hint.add_theme_font_size_override("font_size", 12)
	hint.modulate = Color(0.7, 0.7, 0.7)
	if is_instance_valid(_pixel_font):
		hint.add_theme_font_override("font", _pixel_font)
	panel.add_child(hint)

	# 頛詨獢?
	var line_edit = LineEdit.new()
	line_edit.name = "NameInput"
	line_edit.position = Vector2(16, 65)
	line_edit.size = Vector2(360, 36)
	line_edit.placeholder_text = "頛詨?迂..."
	line_edit.max_length = 16
	line_edit.text = GameManager.player_data.get("display_name", "")
	line_edit.add_theme_font_size_override("font_size", 16)
	if is_instance_valid(_pixel_font):
		line_edit.add_theme_font_override("font", _pixel_font)
	panel.add_child(line_edit)

	# 蝣箄???
	var confirm_btn = Button.new()
	confirm_btn.name = "ConfirmBtn"
	confirm_btn.text = "蝣箄?"
	confirm_btn.position = Vector2(390, 65)
	confirm_btn.size = Vector2(96, 36)
	confirm_btn.add_theme_font_size_override("font_size", 14)
	if is_instance_valid(_pixel_font):
		confirm_btn.add_theme_font_override("font", _pixel_font)
	confirm_btn.pressed.connect(func():
		var name_input = panel.get_node_or_null("NameInput")
		if not name_input:
			return
		var new_name = name_input.text.strip_edges()
		if new_name.length() == 0 or new_name.length() > 16:
			var err_tween = create_tween().set_loops(3)
			err_tween.tween_property(name_input, "modulate", Color(1.0, 0.3, 0.3), 0.1)
			err_tween.tween_property(name_input, "modulate", Color.WHITE, 0.1)
			return
		NetworkManager.send_set_display_name(new_name)
		AudioManager.play_sfx(AudioManager.SFX.WEED_PULL)
		if is_instance_valid(_name_dialog):
			_name_dialog.queue_free()
			_name_dialog = null
	)
	panel.add_child(confirm_btn)

	# ????
	var cancel_btn = Button.new()
	cancel_btn.text = "??"
	cancel_btn.position = Vector2(16, 115)
	cancel_btn.size = Vector2(96, 32)
	cancel_btn.add_theme_font_size_override("font_size", 13)
	cancel_btn.modulate = Color(0.7, 0.7, 0.7)
	if is_instance_valid(_pixel_font):
		cancel_btn.add_theme_font_override("font", _pixel_font)
	cancel_btn.pressed.connect(func():
		if is_instance_valid(_name_dialog):
			_name_dialog.queue_free()
			_name_dialog = null
	)
	panel.add_child(cancel_btn)

	# 瘛∪?
	dialog.modulate.a = 0.0
	var tween = create_tween()
	tween.tween_property(dialog, "modulate:a", 1.0, 0.15)
	line_edit.grab_focus()
	line_edit.select_all()


# ?? ???鈭辣嚗AY-022嚗??????????????????????????????????????????

func _on_combo_event(combo_data: Dictionary) -> void:
	var combo_count = combo_data.get("combo_count", 1)
	if combo_count < 2:
		return
	# ?函?唬?蝵桅＊蝷粹???寞?
	HitEffect.spawn_combo(combo_count, Vector2(640, 580))
	# ????單?嚗 kill.wav嚗靽???
	AudioManager.play_sfx(AudioManager.SFX.KILL)

# ?? 瘥隞餃??Ｘ嚗AY-037嚗??????????????????????????????????????

# ?? 瘥隞餃??Ｘ嚗AY-037嚗AY-053 ????MissionPanel.gd嚗??????????????????

var _mission_panel_node: MissionPanelScript = null  # ?函??Ｘ蝭暺?

## ???遙??選?DAY-053嚗?
func _init_mission_panel() -> void:
	var panel = MissionPanelScript.new()
	add_child(panel)
	panel.setup(_pixel_font)
	panel.mission_completed_notify.connect(_on_mission_completed_notify)
	var top_bar = get_node_or_null("TopBar")
	panel.create_button(top_bar)
	_mission_panel_node = panel

## 隞餃?摰??嚗? MissionPanel 頧?唳?撠梢蝟餌絞嚗?
func _on_mission_completed_notify(mission_data: Dictionary) -> void:
	var name_str = mission_data.get("name", "隞餃?摰?")
	var icon = mission_data.get("icon", "??")
	var reward = mission_data.get("reward", 0)
	_achievement_queue.append({
		"name": "%s %s" % [icon, name_str],
		"description": "摰?嚗?????%d" % reward,
		"icon": icon,
		"type": "special"
	})
	if not _achievement_showing:
		_show_next_achievement()

## 閮剖?隞餃??蔭??嚗 GameManager ?澆嚗AY-038嚗?
func set_mission_reset_at(reset_at_ms: int) -> void:
	if is_instance_valid(_mission_panel_node):
		_mission_panel_node.set_mission_reset_at(reset_at_ms)



# ── Session Stats 面板（DAY-046，DAY-053 拆分為 SessionStatsPanel.gd）──────

var _session_stats_node: SessionStatsPanelScript = null  # 獨立面板節點

## 初始化 Session Stats 面板（DAY-053）
func _init_session_stats() -> void:
	var panel = SessionStatsPanelScript.new()
	add_child(panel)
	panel.setup(_pixel_font)
	var top_bar = get_node_or_null("TopBar")
	panel.create_button(top_bar)
	_session_stats_node = panel


# ── Progressive Jackpot 面板（DAY-048，DAY-053 拆分為 JackpotPanel.gd）──────

var _jackpot_panel_node: JackpotPanelScript = null  # 獨立面板節點

## 初始化 Jackpot 面板（DAY-053）
func _init_jackpot_panel() -> void:
	var panel = JackpotPanelScript.new()
	panel.position = Vector2(320, 42)  # TopBar 下方，畫面中央
	panel.size = Vector2(640, 36)
	panel.z_index = 10
	add_child(panel)
	panel.setup(_pixel_font)
	_jackpot_panel_node = panel


# ── ESC 快捷鍵（DAY-049，DAY-053 更新）──────────────────────────

## 覆寫 _input 處理 ESC 快捷鍵
func _input(event: InputEvent) -> void:
	if event is InputEventKey and event.pressed and not event.echo:
		if event.keycode == KEY_ESCAPE:
			if is_instance_valid(_session_stats_node):
				_session_stats_node.toggle()
			get_viewport().set_input_as_handled()

# ── 觀戰者加入通知（DAY-054d）──────────────────────────────────────────────

func _on_spectator_joined(spectator_data: Dictionary) -> void:
	var count = spectator_data.get("spectator_count", 1)
	# 更新 TopBar 觀戰者計數
	_update_spectator_count_label(count)
	# 用成就通知系統顯示觀戰者加入（複用現有 UI）
	_achievement_queue.append({
		"icon": "👁️",
		"name": "有人在觀戰！",
		"description": "目前 %d 位觀戰者" % count,
		"type": "special"
	})
	if not _achievement_showing:
		_show_next_achievement()

# ── 觀戰者離開通知（DAY-055）──────────────────────────────────────────────

func _on_spectator_left(spectator_data: Dictionary) -> void:
	var count = spectator_data.get("spectator_count", 0)
	# 更新 TopBar 觀戰者計數
	_update_spectator_count_label(count)
	if count > 0:
		# 還有觀戰者，靜默更新（不打擾玩家）
		return
	# 最後一位觀戰者離開，顯示通知
	_achievement_queue.append({
		"icon": "👋",
		"name": "觀戰者離開了",
		"description": "目前無觀戰者",
		"type": "normal"
	})
	if not _achievement_showing:
		_show_next_achievement()

## 更新 TopBar 觀戰者計數標籤（DAY-055）
func _update_spectator_count_label(count: int) -> void:
	var top_bar = get_node_or_null("TopBar")
	if not is_instance_valid(top_bar):
		return
	var lbl = top_bar.get_node_or_null("SpectatorCountLabel")
	if not is_instance_valid(lbl):
		return
	if count > 0:
		lbl.text = "👁️ %d" % count
		lbl.modulate = Color(0.7, 0.85, 1.0, 0.9)
	else:
		lbl.text = ""

# ── 每日登入獎勵彈窗（DAY-065）──────────────────────────────────────────────

func _on_daily_bonus_received(bonus_data: Dictionary) -> void:
	var streak = bonus_data.get("streak", 1)
	var reward = bonus_data.get("reward", 0)
	var max_streak = bonus_data.get("max_streak", streak)
	_show_daily_bonus_popup(streak, reward, max_streak)

## 顯示每日登入獎勵彈窗
func _show_daily_bonus_popup(streak: int, reward: int, max_streak: int) -> void:
	# 建立彈窗容器（CanvasLayer 確保在最上層）
	var canvas = CanvasLayer.new()
	canvas.layer = 200
	add_child(canvas)

	# 半透明背景遮罩
	var overlay = ColorRect.new()
	overlay.color = Color(0, 0, 0, 0.6)
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.mouse_filter = Control.MOUSE_FILTER_STOP
	canvas.add_child(overlay)

	# 彈窗主體
	var panel = PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_CENTER)
	panel.custom_minimum_size = Vector2(380, 280)
	panel.position = Vector2(-190, -140)
	canvas.add_child(panel)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 12)
	panel.add_child(vbox)

	# 標題
	var title_lbl = Label.new()
	title_lbl.text = "🌟 每日登入獎勵"
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 22)
	title_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	vbox.add_child(title_lbl)

	# 連續天數
	var streak_lbl = Label.new()
	var streak_text = "連續登入 %d 天" % streak
	if streak >= 7:
		streak_text += " 🔥"
	streak_lbl.text = streak_text
	streak_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	streak_lbl.add_theme_font_size_override("font_size", 16)
	streak_lbl.add_theme_color_override("font_color", Color(1.0, 0.7, 0.3))
	vbox.add_child(streak_lbl)

	# 獎勵金額（大字）
	var reward_lbl = Label.new()
	reward_lbl.text = "🪙 +%d 金幣" % reward
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_lbl.add_theme_font_size_override("font_size", 32)
	reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.95, 0.4))
	vbox.add_child(reward_lbl)

	# 7天獎勵預覽條
	var preview_lbl = Label.new()
	preview_lbl.text = "7天獎勵：500 → 800 → 1200 → 1800 → 2500 → 3500 → 5000"
	preview_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	preview_lbl.add_theme_font_size_override("font_size", 11)
	preview_lbl.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	preview_lbl.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	vbox.add_child(preview_lbl)

	# 最高連續天數
	if max_streak > 1:
		var max_lbl = Label.new()
		max_lbl.text = "最高連續：%d 天" % max_streak
		max_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		max_lbl.add_theme_font_size_override("font_size", 12)
		max_lbl.add_theme_color_override("font_color", Color(0.6, 0.9, 0.6))
		vbox.add_child(max_lbl)

	# 確認按鈕
	var btn = Button.new()
	btn.text = "太棒了！"
	btn.custom_minimum_size = Vector2(160, 40)
	btn.add_theme_font_size_override("font_size", 16)
	var btn_container = CenterContainer.new()
	btn_container.add_child(btn)
	vbox.add_child(btn_container)

	# 彈入動畫
	panel.scale = Vector2(0.5, 0.5)
	panel.modulate.a = 0.0
	var tween = panel.create_tween()
	tween.set_parallel(true)
	tween.tween_property(panel, "scale", Vector2(1.0, 1.0), 0.3).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.tween_property(panel, "modulate:a", 1.0, 0.2)

	# 按鈕關閉
	btn.pressed.connect(func():
		var close_tween = panel.create_tween()
		close_tween.set_parallel(true)
		close_tween.tween_property(panel, "scale", Vector2(0.8, 0.8), 0.15)
		close_tween.tween_property(panel, "modulate:a", 0.0, 0.15)
		close_tween.tween_callback(canvas.queue_free).set_delay(0.15)
	)

	# 5 秒後自動關閉
	get_tree().create_timer(5.0).timeout.connect(func():
		if is_instance_valid(canvas):
			canvas.queue_free()
	)


# ── 週賽排名面板（DAY-066）──────────────────────────────────────────────

var _tournament_panel_node = null

## 初始化週賽排名面板（DAY-066）
func _init_tournament_panel() -> void:
	var panel = TournamentPanelScript.new()
	panel.position = Vector2(0, 0)
	panel.z_index = 8
	add_child(panel)
	panel.setup(_pixel_font)
	_tournament_panel_node = panel


# ── 武器升級面板（DAY-067）──────────────────────────────────────────────

var _weapon_panel_node = null

## 初始化武器升級面板（DAY-067）
## 位置：BottomBar 左側（x=10, y=畫面高度-90）
func _init_weapon_panel() -> void:
	var panel = WeaponPanelScript.new()
	# 放在畫面左下角，BottomBar 上方
	panel.position = Vector2(10, 540)
	panel.z_index = 8
	add_child(panel)
	panel.setup(_pixel_font)
	_weapon_panel_node = panel

## ── 稱號面板（DAY-068）──────────────────────────────────────────────────────
var _title_panel_node = null

## 初始化稱號面板（DAY-068）
## 位置：TopBar 下方左側，顯示玩家當前稱號
func _init_title_panel() -> void:
	var panel = TitlePanelScript.new()
	# 放在 TopBar 下方，金幣顯示旁邊
	panel.position = Vector2(10, 44)
	panel.z_index = 7
	add_child(panel)
	panel.setup(_pixel_font)
	_title_panel_node = panel

## ── 砲台外觀面板（DAY-071）──────────────────────────────────────────────────────
var _skin_panel_node = null

## 初始化砲台外觀面板（DAY-071）
## 位置：BottomBar 右側（WeaponPanel 右邊）
func _init_skin_panel() -> void:
	var panel = SkinPanelScript.new()
	# 放在 WeaponPanel 右側（WeaponPanel 在 x=10，寬 200，所以 SkinPanel 在 x=215）
	panel.position = Vector2(215, 540)
	panel.z_index = 8
	add_child(panel)
	panel.setup(_pixel_font)
	_skin_panel_node = panel

## ── 賽季通行證面板（DAY-072）──────────────────────────────────────────────────────
var _season_panel_node = null

## 初始化賽季通行證面板（DAY-072）
## 位置：TopBar 右側（從右往左排列，x=1248）
func _init_season_panel() -> void:
	var panel = SeasonPanelScript.new()
	panel.position = Vector2(1248, 4)
	panel.z_index = 10
	add_child(panel)
	panel.setup(_pixel_font)
	_season_panel_node = panel

## ── 好友系統面板（DAY-073）──────────────────────────────────────────────────────
var _friend_panel_node = null

## 初始化好友系統面板（DAY-073）
## 位置：TopBar 右側（從右往左排列，x=1216）
func _init_friend_panel() -> void:
	var panel = FriendPanelScript.new()
	panel.position = Vector2(1216, 4)
	panel.z_index = 10
	add_child(panel)
	panel.setup(_pixel_font)
	_friend_panel_node = panel

## ── 公會系統面板（DAY-074）──────────────────────────────────────────────────────
var _guild_panel_node = null

## 初始化公會系統面板（DAY-074）
## 位置：TopBar 右側（從右往左排列，x=1184）
func _init_guild_panel() -> void:
	var panel = GuildPanelScript.new()
	panel.position = Vector2(1184, 4)
	panel.z_index = 10
	add_child(panel)
	panel.setup(_pixel_font)
	_guild_panel_node = panel

## ── 公會戰面板（DAY-076）──────────────────────────────────────────────────────
var _guild_war_panel_node = null

## 初始化公會戰面板（DAY-076）
## 位置：TopBar 右側（從右往左排列，x=1152）
func _init_guild_war_panel() -> void:
	var panel = GuildWarPanelScript.new()
	panel.position = Vector2(1152, 4)
	panel.z_index = 10
	add_child(panel)
	panel.setup(_pixel_font)
	_guild_war_panel_node = panel

## ── 每日 BOSS 挑戰面板（DAY-077）──────────────────────────────────────────────
var _daily_boss_panel_node = null

## 初始化每日 BOSS 挑戰面板（DAY-077）
## 位置：TopBar 右側（從右往左排列，x=1120）
func _init_daily_boss_panel() -> void:
	var panel = DailyBossPanelScript.new()
	panel.position = Vector2(1120, 4)
	panel.z_index = 10
	add_child(panel)
	panel.setup(_pixel_font)
	_daily_boss_panel_node = panel

## ── VIP 等級面板（DAY-078）──────────────────────────────────────────────────
var _vip_panel_node = null

## 初始化 VIP 等級面板（DAY-078）
## 位置：TopBar 右側（從右往左排列，x=1088）
func _init_vip_panel() -> void:
	var panel = VIPPanelScript.new()
	panel.position = Vector2(1088, 4)
	panel.z_index = 10
	add_child(panel)
	panel.setup(_pixel_font)
	_vip_panel_node = panel

## ── 限時活動面板（DAY-079）──────────────────────────────────────────────────
var _event_panel_node = null

## 初始化限時活動面板（DAY-079）
## 位置：TopBar 右側（從右往左排列，x=1056）
func _init_event_panel() -> void:
	var panel = EventPanelScript.new()
	panel.position = Vector2(1056, 4)
	panel.z_index = 10
	add_child(panel)
	panel.setup(_pixel_font)
	_event_panel_node = panel

## ── 魚類圖鑑面板（DAY-081）──────────────────────────────────────────────────
var _codex_panel_node = null

## 初始化魚類圖鑑面板（DAY-081）
## 位置：TopBar 右側（從右往左排列，x=1024）
func _init_codex_panel() -> void:
	var panel = CodexPanelScript.new()
	panel.position = Vector2(1024, 4)
	panel.z_index = 10
	add_child(panel)
	panel.setup(_pixel_font)
	_codex_panel_node = panel

## ── 連擊系統面板（DAY-083）──────────────────────────────────────────────────
var _streak_panel_node = null

## 初始化連擊系統面板（DAY-083）
## 位置：遊戲畫面中央上方（顯示當前連擊數）
func _init_streak_panel() -> void:
	var panel = StreakPanelScript.new()
	# 放在畫面中央上方，TopBar 下方
	panel.position = Vector2(740, 48)
	panel.z_index = 12
	add_child(panel)
	panel.setup(_pixel_font)
	_streak_panel_node = panel

## ── 推薦碼面板（DAY-082）──────────────────────────────────────────────────
var _referral_panel_node = null

## 初始化推薦碼面板（DAY-082）
## 位置：TopBar 右側（從右往左排列，x=992）
func _init_referral_panel() -> void:
	var panel = ReferralPanelScript.new()
	panel.position = Vector2(992, 4)
	panel.z_index = 10
	add_child(panel)
	panel.setup(_pixel_font)
	_referral_panel_node = panel

## ── 幸運轉盤面板（DAY-084）──────────────────────────────────────────────────
var _wheel_panel_node = null

## 初始化幸運轉盤面板（DAY-084）
## 位置：畫面中央（全螢幕覆蓋式彈窗）
func _init_wheel_panel() -> void:
	var panel = WheelPanelScript.new()
	panel.position = Vector2(640, 360)
	panel.z_index = 50
	add_child(panel)
	panel.setup(_pixel_font)
	_wheel_panel_node = panel
