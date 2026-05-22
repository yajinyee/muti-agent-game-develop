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
const ChallengePanelScript = preload("res://scripts/ui/ChallengePanel.gd")
const MissionStreakPanelScript = preload("res://scripts/ui/MissionStreakPanel.gd")
const RoulettePanelScript = preload("res://scripts/ui/RoulettePanel.gd")
const BuyBonusPanelScript = preload("res://scripts/ui/BuyBonusPanel.gd")
const RaidPanelScript = preload("res://scripts/ui/RaidPanel.gd")
const FragmentPanelScript = preload("res://scripts/ui/FragmentPanel.gd")
const LuckyCatchPanelScript = preload("res://scripts/ui/LuckyCatchPanel.gd")
const RapidRespinPanelScript = preload("res://scripts/ui/RapidRespinPanel.gd")
const TreasureMapPanelScript = preload("res://scripts/ui/TreasureMapPanel.gd")
const FlashChallengePanelScript = preload("res://scripts/ui/FlashChallengePanel.gd")

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
	_init_friend_panel()      # 好友系統面板（DAY-073）+ 挑戰面板（DAY-102）+ 私訊面板（DAY-103）
	_init_dm_panel()          # 私訊面板（DAY-103）
	_init_guild_panel()       # 公會系統面板（DAY-074）
	_init_guild_war_panel()   # 公會戰面板（DAY-076）
	_init_daily_boss_panel()  # 每日 BOSS 挑戰面板（DAY-077）
	_init_vip_panel()         # VIP 等級面板（DAY-078）
	_init_event_panel()       # 限時活動面板（DAY-079）
	_init_codex_panel()       # 魚類圖鑑面板（DAY-081）
	_init_streak_panel()      # 連擊系統面板（DAY-083）
	_init_referral_panel()    # 推薦碼面板（DAY-082）
	_init_wheel_panel()       # 幸運轉盤面板（DAY-084）
	_init_challenge_panel()   # 隱藏挑戰面板（DAY-085）
	_init_mission_streak_panel() # 任務連續完成面板（DAY-086）
	_init_weather_panel()     # 天氣系統面板（DAY-087）
	_init_chain_panel()       # 連鎖爆炸面板（DAY-088）
	_init_special_weapon_panel() # 特殊武器面板（DAY-089）
	_init_mystery_box_panel()    # 神秘寶箱面板（DAY-090）
	_init_room_select_panel()    # 房間難度選擇面板（DAY-091）
	_init_daily_spin_panel()     # 每日簽到轉盤面板（DAY-092）
	_init_shop_panel()           # 商店面板（DAY-094）
	_init_player_stats_panel()   # 玩家統計面板（DAY-096）
	_init_announcement_panel()   # 全服公告面板（DAY-097）
	_init_activity_feed_panel()  # 成就動態牆面板（DAY-112）
	_init_player_card_panel()    # 玩家名片面板（DAY-106）
	_init_login_milestone_panel() # 登入里程碑面板（DAY-107）
	_init_player_journey_panel()  # 玩家旅程儀表板（DAY-108）
	_init_roulette_panel()        # 雙層倍率輪盤面板（DAY-113）
	_init_buy_bonus_panel()       # Buy Bonus 面板（DAY-114）
	_init_raid_panel()            # Co-op Boss Raid 面板（DAY-115）
	_init_fragment_panel()        # 碎片收集大獎面板（DAY-116）
	_init_lucky_catch_panel()     # 幸運捕獲通知面板（DAY-119）
	_init_rapid_respin_panel()    # Rapid Respin 通知面板（DAY-121）
	_init_treasure_map_panel()    # 寶藏地圖面板（DAY-122）
	_init_flash_challenge_panel() # 閃電挑戰面板（DAY-123）
	_init_rare_target_alert()     # 傳說目標警報（DAY-124）
	_init_golden_time_panel()     # 黃金時間面板（DAY-125）
	_init_rare_catch_panel()      # 稀有連擊面板（DAY-126）
	_init_weather_surge_panel()   # 天氣湧現事件面板（DAY-127）
	_init_dragon_wrath_panel()    # 龍怒蓄力大招面板（DAY-128）
	_init_immortal_boss_panel()   # 不死 BOSS 連勝面板（DAY-129）
	_init_awaken_boss_panel()     # 覺醒 BOSS 面板（DAY-130）
	_init_win_streak_panel()      # 連勝獎勵面板（DAY-131）
	_init_lightning_eel_panel()   # 閃電鰻連鎖攻擊面板（DAY-132）
	_init_fever_mode_panel()      # 狂熱模式面板（DAY-133）
	_init_unlucky_bonus_panel()   # 失敗補償面板（DAY-135）
	_init_speed_race_panel()      # 競速獵殺面板（DAY-136）
	_init_bounty_panel()          # 全服目標懸賞面板（DAY-137）
	_init_mult_storm_panel()      # 全服倍率風暴面板（DAY-138）
	_init_dual_roulette_panel()   # 雙環輪盤面板（DAY-139）
	_init_mega_catch_panel()      # Mega Catch 事件面板（DAY-140）
	_init_drill_lobster_panel()   # 鑽頭龍蝦連帶效果面板（DAY-142）
	_init_bomb_crab_panel()       # 炸彈蟹連環爆炸面板（DAY-143）
	_init_mega_octopus_panel()    # 巨型章魚轉盤面板（DAY-144）
	_init_anglerfish_panel()      # 巨型鮟鱇魚電擊寶箱面板（DAY-145）
	_init_crocodile_panel()       # 巨型鹹水鱷魚獵魚面板（DAY-146）
	_init_giant_prize_fish_panel() # 夢幻巨型獎勵魚面板（DAY-147）
	_init_chainlong_wheel_panel()  # 千龍王強化輪盤面板（DAY-148）
	_init_golden_jellyfish_panel() # 黃金水母全場電擊面板（DAY-149）
	_init_thunderbolt_lobster_panel() # 雷霆龍蝦免費射擊面板（DAY-150）
	_init_rainbow_phoenix_panel()     # 彩虹鳳凰 Power Up 面板（DAY-151）
	_init_vampire_panel()             # 吸血鬼成長倍率面板（DAY-152）
	_init_crystal_dragon_panel()      # 水晶龍收集大獎面板（DAY-153）
	_init_royal_chain_lightning_panel() # 皇家閃電鰻持續連鎖電擊面板（DAY-156）
	_init_golden_turtle_panel()       # 黃金海龜時間停止面板（DAY-159）
	_init_lucky_star_fish_panel()     # 幸運星魚全場倍率翻倍面板（DAY-160）
	_init_golden_shark_panel()        # 黃金鯊魚全服狂暴模式面板（DAY-161）
	_init_money_fish_panel()          # 金幣魚王即時獎勵面板（DAY-162）
	_init_captain_fish_panel()        # 船長魚全服競速模式面板（DAY-163）
	_init_abyss_whale_panel()         # 深淵巨鯨全服 Boss 挑戰面板（DAY-164）
	_init_black_hole_panel()          # 黑洞漩渦武器視覺效果面板（DAY-166）
	_init_roulette_crab_panel()       # 黃金輪盤螃蟹面板（DAY-167）
	_init_lion_dance_panel()          # 獅子舞大獎爆發面板（DAY-168）
	_init_vortex_fish_panel()         # 漩渦魚群吸引面板（DAY-169）
	_init_freeze_bomb_panel()         # 冰凍炸彈魚面板（DAY-170）
	_init_ice_fishing_panel()         # 冰釣幸運輪盤面板（DAY-171）
	_init_lucky_egg_panel()           # 幸運彩蛋魚面板（DAY-172）
	_init_rainbow_lucky_panel()       # 彩虹幸運魚面板（DAY-173）
	_init_sea_anemone_panel()         # 海葵觸手攻擊面板（DAY-174）
	_init_lucky_dice_panel()          # 幸運骰子魚面板（DAY-175）
	_init_fire_storm_panel()          # 火焰風暴魚面板（DAY-176）
	_init_golden_treasure_panel()     # 黃金寶藏魚面板（DAY-177）
	_init_mermaid_healing_panel()     # 美人魚治癒面板（DAY-178）
	_init_lucky_clover_panel()        # 幸運草魚面板（DAY-179）
	_init_rainbow_shark_panel()       # 彩虹鯊魚爆發面板（DAY-180）
	_init_thunder_shark_panel()       # 雷霆鯊魚連鎖閃電面板（DAY-181）
	_init_vampire_fish_panel()        # 吸血鬼魚累積倍率面板（DAY-182）
	_init_lightning_auto_chain_panel() # 閃電魚自動連鎖面板（DAY-183）
	_init_meteor_fish_panel()          # 隕石魚隕石雨面板（DAY-184）
	_init_phoenix_fish_panel()         # 鳳凰魚涅槃重生面板（DAY-185）
	_init_dragon_turtle_panel()        # 龍龜不死 Boss 面板（DAY-186）
	_init_chain_bomb_panel()           # 連鎖爆炸魚面板（DAY-187）
	_init_crocodile_hunter_panel()     # 巨型鱷魚獵食面板（DAY-188）
	_init_time_bomb_fish_panel()       # 時間炸彈魚面板（DAY-189）
	_init_triple_lucky_fish_panel()    # 三重幸運魚面板（DAY-190）
	_init_school_panic_panel()         # 魚群驚嚇連帶面板（DAY-191）
	_init_rock_skeleton_concert_panel() # 搖滾骷髏演唱會面板（DAY-192）
	_init_electric_jellyfish_panel()    # 電流水母電流網路面板（DAY-193）
	_init_chainlong_king_panel()        # 長龍王雙環輪盤面板（DAY-194）
	_init_drill_lobster_panel()         # 鑽頭龍蝦穿透爆炸面板（DAY-195）
	_init_anglerfish_electric_panel()   # 巨型鮟鱇魚電擊寶箱面板（DAY-196）
	_init_mystic_dragon_panel()         # 神秘龍魚八波攻擊面板（DAY-197）
	_init_ghost_fish_panel()            # 幽靈魚分身面板（DAY-198）
	_init_thunderbolt_lobster_panel()   # 雷霆龍蝦免費射擊面板（DAY-199）

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
	panel.size = Vector2(640, 66)  # DAY-118：高度從 36 → 66（加進度條）
	panel.z_index = 10
	add_child(panel)
	panel.setup(_pixel_font)
	_jackpot_panel_node = panel




# ── 排行榜面板（DAY-058 拆分到 LeaderboardPanel.gd）──────────────────────────

var _leaderboard_node = null  # LeaderboardPanelScript 實例

## 初始化排行榜面板（委派給 LeaderboardPanel.gd）
func _create_leaderboard_panel() -> void:
	_leaderboard_node = LeaderboardPanelScript.new()
	_leaderboard_node.setup(self, _pixel_font)

## 排行榜更新事件（委派給 LeaderboardPanel.gd）
func _on_leaderboard_updated(entries: Array) -> void:
	if _leaderboard_node:
		_leaderboard_node.update(entries, GameManager.get_player_id())

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
	# 初始化好友挑戰面板（DAY-102）
	_init_challenge_pvp_panel()

## ── 好友挑戰面板（DAY-102）──────────────────────────────────────────────────────
var _challenge_pvp_panel_node = null

## 初始化好友挑戰面板（DAY-102）
## 位置：左下角，顯示進行中的挑戰分數
func _init_challenge_pvp_panel() -> void:
	var ChallengePvPPanelScript = load("res://scripts/ui/ChallengePvPPanel.gd")
	if ChallengePvPPanelScript == null:
		return
	var panel = ChallengePvPPanelScript.new()
	panel.position = Vector2(8, 540)
	panel.z_index = 50
	add_child(panel)
	if panel.has_method("setup"):
		panel.setup(_pixel_font)
	_challenge_pvp_panel_node = panel

## ── 私訊面板（DAY-103）──────────────────────────────────────────────────────────
var _dm_panel_node = null

## 初始化私訊面板（DAY-103）
## 位置：TopBar 右側（x=1248）
func _init_dm_panel() -> void:
	var DMPanelScript = load("res://scripts/ui/DMPanel.gd")
	if DMPanelScript == null:
		return
	var panel = DMPanelScript.new()
	panel.position = Vector2(1248, 4)
	panel.z_index = 10
	add_child(panel)
	if panel.has_method("setup"):
		panel.setup(_pixel_font)
	_dm_panel_node = panel
	# 連接開啟 DM 面板訊號
	if GameManager.has_signal("open_dm_panel"):
		GameManager.open_dm_panel.connect(func(friend_id: String, friend_name: String):
			if is_instance_valid(_dm_panel_node) and _dm_panel_node.has_method("open_conversation"):
				_dm_panel_node.open_conversation(friend_id, friend_name)
		)## ── 公會系統面板（DAY-074）──────────────────────────────────────────────────────
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

## ── 隱藏挑戰面板（DAY-085）──────────────────────────────────────────────────
var _challenge_panel_node = null

## 初始化隱藏挑戰面板（DAY-085）
## 位置：畫面上方中央（挑戰解鎖通知彈窗）
func _init_challenge_panel() -> void:
	var panel = ChallengePanelScript.new()
	panel.z_index = 60
	add_child(panel)
	panel.setup(_pixel_font)
	_challenge_panel_node = panel

## ── 任務連續完成面板（DAY-086）──────────────────────────────────────────────
var _mission_streak_panel_node = null

## 初始化任務連續完成面板（DAY-086）
func _init_mission_streak_panel() -> void:
	var panel = MissionStreakPanelScript.new()
	panel.z_index = 55
	add_child(panel)
	panel.setup(_pixel_font)
	_mission_streak_panel_node = panel

# ---- 天氣系統面板（DAY-087）----

const WeatherPanelScript = preload("res://scripts/ui/WeatherPanel.gd")
var _weather_panel: Control = null

func _init_weather_panel() -> void:
	var panel = WeatherPanelScript.new()
	panel.name = "WeatherPanel"
	add_child(panel)
	_weather_panel = panel

# ---- 連鎖爆炸面板（DAY-088）----

const ChainExplosionPanelScript = preload("res://scripts/ui/ChainExplosionPanel.gd")
var _chain_panel: Control = null

func _init_chain_panel() -> void:
	var panel = ChainExplosionPanelScript.new()
	panel.name = "ChainExplosionPanel"
	add_child(panel)
	_chain_panel = panel

# ---- 特殊武器面板（DAY-089）----

const SpecialWeaponPanelScript = preload("res://scripts/ui/SpecialWeaponPanel.gd")
var _special_weapon_panel: Control = null

func _init_special_weapon_panel() -> void:
	var panel = SpecialWeaponPanelScript.new()
	panel.name = "SpecialWeaponPanel"
	# 放在 WeaponPanel 右側（WeaponPanel x=10 寬200，SkinPanel x=215 寬200，SpecialWeapon x=420）
	# DAY-134：面板寬度從 240 升級到 320（加入龍捲風砲）
	# DAY-141：面板寬度從 320 升級到 400（加入追蹤飛彈）
	panel.position = Vector2(420, 540)
	panel.z_index = 8
	add_child(panel)
	panel.setup(_pixel_font)
	_special_weapon_panel = panel
	# 連接選擇訊號（讓 Cannon.gd 知道玩家選了特殊武器）
	panel.weapon_selected.connect(_on_special_weapon_selected)

func _on_special_weapon_selected(weapon_type: String) -> void:
	# 通知 GameManager 當前選中的特殊武器
	GameManager.set_meta("selected_special_weapon", weapon_type)

# ---- 神秘寶箱面板（DAY-090）----

const MysteryBoxPanelScript = preload("res://scripts/ui/MysteryBoxPanel.gd")
var _mystery_box_panel: Control = null

func _init_mystery_box_panel() -> void:
	var panel = MysteryBoxPanelScript.new()
	panel.name = "MysteryBoxPanel"
	# 放在 SpecialWeaponPanel 右側（x=420 寬400，所以 MysteryBox 在 x=825）
	# DAY-134：SpecialWeaponPanel 寬度從 240 升級到 320，MysteryBox 右移 80px
	# DAY-141：SpecialWeaponPanel 寬度從 320 升級到 400，MysteryBox 再右移 80px
	# DAY-154：SpecialWeaponPanel 寬度從 400 升級到 480，MysteryBox 再右移 80px
	# DAY-155：SpecialWeaponPanel 寬度從 480 升級到 560，MysteryBox 再右移 80px
	# DAY-157：SpecialWeaponPanel 寬度從 560 升級到 640，MysteryBox 再右移 80px
	# DAY-166：SpecialWeaponPanel 寬度從 640 升級到 720，MysteryBox 再右移 80px
	panel.position = Vector2(1145, 540)
	panel.z_index = 8
	add_child(panel)
	panel.setup(_pixel_font)
	_mystery_box_panel = panel

const RoomSelectPanelScript = preload("res://scripts/ui/RoomSelectPanel.gd")
var _room_select_panel: Control = null
var _room_btn: Button = null

func _init_room_select_panel() -> void:
	# 房間選擇面板（全螢幕覆蓋）
	var panel = RoomSelectPanelScript.new()
	panel.name = "RoomSelectPanel"
	panel.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 90
	add_child(panel)
	panel.setup(_pixel_font)
	_room_select_panel = panel

	# 房間切換按鈕（TopBar 右側）
	_room_btn = Button.new()
	_room_btn.name = "RoomBtn"
	_room_btn.text = "🏠 房間"
	_room_btn.position = Vector2(1100, 4)
	_room_btn.size = Vector2(80, 30)
	_room_btn.add_theme_color_override("font_color", Color(0.9, 0.85, 0.3))
	_room_btn.add_theme_font_size_override("font_size", 12)
	if is_instance_valid(_pixel_font):
		_room_btn.add_theme_font_override("font", _pixel_font)
	_room_btn.pressed.connect(_on_room_btn_pressed)
	add_child(_room_btn)

	# 連接房間切換成功訊號，更新按鈕文字
	if GameManager.has_signal("room_switched"):
		GameManager.room_switched.connect(_on_room_switched_hud)

func _on_room_btn_pressed() -> void:
	if is_instance_valid(_room_select_panel):
		_room_select_panel.show_panel()

func _on_room_switched_hud(data: Dictionary) -> void:
	var icon = data.get("room_icon", "🏠")
	var name_str = data.get("room_name", "房間")
	if is_instance_valid(_room_btn):
		_room_btn.text = icon + " " + name_str

const DailySpinPanelScript = preload("res://scripts/ui/DailySpinPanel.gd")
var _daily_spin_panel: Control = null
var _daily_spin_btn: Button = null

func _init_daily_spin_panel() -> void:
	# 每日轉盤面板（全螢幕覆蓋）
	var panel = DailySpinPanelScript.new()
	panel.name = "DailySpinPanel"
	panel.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 85
	add_child(panel)
	panel.setup(_pixel_font)
	_daily_spin_panel = panel

	# 每日轉盤按鈕（TopBar）
	_daily_spin_btn = Button.new()
	_daily_spin_btn.name = "DailySpinBtn"
	_daily_spin_btn.text = "🎡 轉盤"
	_daily_spin_btn.position = Vector2(1190, 4)
	_daily_spin_btn.size = Vector2(80, 30)
	_daily_spin_btn.add_theme_color_override("font_color", Color(0.3, 1.0, 0.5))
	_daily_spin_btn.add_theme_font_size_override("font_size", 12)
	if is_instance_valid(_pixel_font):
		_daily_spin_btn.add_theme_font_override("font", _pixel_font)
	_daily_spin_btn.pressed.connect(_on_daily_spin_btn_pressed)
	add_child(_daily_spin_btn)

	# 連接每日轉盤狀態訊號，更新按鈕提示
	if GameManager.has_signal("daily_spin_state"):
		GameManager.daily_spin_state.connect(_on_daily_spin_state_hud)

func _on_daily_spin_btn_pressed() -> void:
	if is_instance_valid(_daily_spin_panel):
		_daily_spin_panel.show_panel()

func _on_daily_spin_state_hud(data: Dictionary) -> void:
	var can_spin = data.get("can_spin", false)
	if is_instance_valid(_daily_spin_btn):
		if can_spin:
			_daily_spin_btn.text = "🎡 轉盤 ✓"
			_daily_spin_btn.add_theme_color_override("font_color", Color(0.3, 1.0, 0.3))
		else:
			_daily_spin_btn.text = "🎡 轉盤"
			_daily_spin_btn.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))

# ---- 商店系統（DAY-094）----
const ShopPanelScript = preload("res://scripts/ui/ShopPanel.gd")
var _shop_panel_node = null

func _init_shop_panel() -> void:
	var panel = ShopPanelScript.new()
	panel.name = "ShopPanel"
	panel.z_index = 50
	add_child(panel)
	if _pixel_font:
		panel.setup(_pixel_font)
	_shop_panel_node = panel

	# 商店按鈕（TopBar）
	var shop_btn := Button.new()
	shop_btn.name = "ShopBtn"
	shop_btn.text = "🛒 商店"
	shop_btn.position = Vector2(1280, 4)
	shop_btn.size = Vector2(80, 30)
	shop_btn.add_theme_color_override("font_color", Color(0.9, 0.6, 1.0))
	shop_btn.add_theme_font_size_override("font_size", 12)
	if is_instance_valid(_pixel_font):
		shop_btn.add_theme_font_override("font", _pixel_font)
	shop_btn.pressed.connect(_on_shop_btn_pressed)
	add_child(shop_btn)

func _on_shop_btn_pressed() -> void:
	if is_instance_valid(_shop_panel_node):
		_shop_panel_node._on_toggle_pressed()

# ---- 玩家統計面板（DAY-096）----
const PlayerStatsPanelScript = preload("res://scripts/ui/PlayerStatsPanel.gd")
var _player_stats_panel_node = null

func _init_player_stats_panel() -> void:
	var panel = PlayerStatsPanelScript.new()
	panel.name = "PlayerStatsPanel"
	panel.position = Vector2(440, 60)
	panel.z_index = 50
	panel.visible = false
	add_child(panel)
	panel.setup(_pixel_font)
	_player_stats_panel_node = panel

	# 統計按鈕（TopBar，x=1370）
	var stats_btn := Button.new()
	stats_btn.name = "StatsBtn"
	stats_btn.text = "📊 統計"
	stats_btn.position = Vector2(1370, 4)
	stats_btn.size = Vector2(72, 30)
	stats_btn.add_theme_color_override("font_color", Color(0.5, 1.0, 0.8))
	stats_btn.add_theme_font_size_override("font_size", 12)
	if is_instance_valid(_pixel_font):
		stats_btn.add_theme_font_override("font", _pixel_font)
	stats_btn.pressed.connect(_on_stats_btn_pressed)
	add_child(stats_btn)

func _on_stats_btn_pressed() -> void:
	if is_instance_valid(_player_stats_panel_node):
		if _player_stats_panel_node.visible:
			_player_stats_panel_node.visible = false
		else:
			_player_stats_panel_node.show_panel()

# ---- 全服公告系統（DAY-097）----
const AnnouncementPanelScript = preload("res://scripts/ui/AnnouncementPanel.gd")

func _init_announcement_panel() -> void:
	var panel = AnnouncementPanelScript.new()
	panel.name = "AnnouncementPanel"
	panel.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 80
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	panel.setup(_pixel_font)

# ---- 成就動態牆系統（DAY-112）----
const ActivityFeedPanelScript = preload("res://scripts/ui/ActivityFeedPanel.gd")

func _init_activity_feed_panel() -> void:
	var panel = ActivityFeedPanelScript.new()
	panel.name = "ActivityFeedPanel"
	panel.z_index = 78  # 在公告面板下方，不遮擋重要通知
	add_child(panel)

## 玩家名片面板（DAY-106）
var _player_card_panel_node = null
const PlayerCardPanelScript = preload("res://scripts/ui/PlayerCardPanel.gd")

func _init_player_card_panel() -> void:
	var panel = PlayerCardPanelScript.new()
	panel.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 90
	add_child(panel)
	_player_card_panel_node = panel

## 顯示玩家名片（供排行榜/公會/好友面板呼叫）
func show_player_card(player_id: String) -> void:
	if is_instance_valid(_player_card_panel_node):
		_player_card_panel_node.show_card(player_id)

## 登入里程碑面板（DAY-107）
var _login_milestone_panel_node = null
const LoginMilestonePanelScript = preload("res://scripts/ui/LoginMilestonePanel.gd")

func _init_login_milestone_panel() -> void:
	var panel = LoginMilestonePanelScript.new()
	panel.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 95
	add_child(panel)
	_login_milestone_panel_node = panel
	# 連接里程碑訊號
	GameManager.login_milestone_reached.connect(_on_login_milestone_reached)
	GameManager.login_progress_received.connect(_on_login_progress_received)

## 里程碑達成通知（DAY-107）
func _on_login_milestone_reached(data: Dictionary) -> void:
	if is_instance_valid(_login_milestone_panel_node):
		_login_milestone_panel_node.show_milestone(data)

## 登入進度回應（DAY-107）
func _on_login_progress_received(data: Dictionary) -> void:
	if is_instance_valid(_login_milestone_panel_node):
		_login_milestone_panel_node.show_progress(data)

## 玩家旅程儀表板（DAY-108）
var _player_journey_panel_node = null
const PlayerJourneyPanelScript = preload("res://scripts/ui/PlayerJourneyPanel.gd")

func _init_player_journey_panel() -> void:
	var panel = PlayerJourneyPanelScript.new()
	panel.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 85
	add_child(panel)
	_player_journey_panel_node = panel

## 賽季節日活動面板（DAY-109）
var _festival_panel_node = null
const FestivalPanelScript = preload("res://scripts/ui/FestivalPanel.gd")

func _init_festival_panel() -> void:
	var panel = FestivalPanelScript.new()
	panel.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 87
	add_child(panel)
	_festival_panel_node = panel
	# 連接節日訊號
	GameManager.festival_updated.connect(_on_festival_updated)
	GameManager.festival_task_ready_signal.connect(_on_festival_task_ready)
	GameManager.festival_title_earned_signal.connect(_on_festival_title_earned)

## 節日狀態更新（DAY-109）
func _on_festival_updated(data: Dictionary) -> void:
	if is_instance_valid(_festival_panel_node):
		_festival_panel_node.update_from_data(data)

## 節日任務可領取通知（DAY-109）
func _on_festival_task_ready(task_id: String) -> void:
	# 顯示提示（可以加 toast 通知）
	print("[HUD] Festival task ready: %s" % task_id)

## 節日稱號獲得通知（DAY-109）
func _on_festival_title_earned(data: Dictionary) -> void:
	var title_name: String = data.get("title_name", "")
	print("[HUD] Festival title earned: %s" % title_name)

## 顯示節日面板（供 TopBar 按鈕呼叫）
func show_festival_panel() -> void:
	if is_instance_valid(_festival_panel_node):
		GameManager.request_festival()
		_festival_panel_node.show()

# ---- 名人堂面板（DAY-110）----
var _hall_of_fame_panel = null

func _init_hall_of_fame_panel() -> void:
	var HallOfFamePanelScript = load("res://scripts/ui/HallOfFamePanel.gd")
	if HallOfFamePanelScript == null:
		return
	_hall_of_fame_panel = HallOfFamePanelScript.new()
	add_child(_hall_of_fame_panel)

	# 連接訊號
	GameManager.hall_of_fame_updated.connect(func(data): 
		if is_instance_valid(_hall_of_fame_panel):
			_hall_of_fame_panel.update_records(data)
	)
	GameManager.hall_of_fame_new_record.connect(func(data):
		if is_instance_valid(_hall_of_fame_panel):
			_hall_of_fame_panel.show_new_record(data)
	)

	# 在 TopBar 加入名人堂按鈕
	var top_bar = get_node_or_null("TopBar")
	if is_instance_valid(top_bar):
		var hof_btn = Button.new()
		hof_btn.name = "HallOfFameButton"
		hof_btn.text = "🏆"
		hof_btn.custom_minimum_size = Vector2(36, 30)
		hof_btn.position = Vector2(1100, 5)
		hof_btn.pressed.connect(func():
			if is_instance_valid(_hall_of_fame_panel):
				_hall_of_fame_panel.show_panel()
		)
		top_bar.add_child(hof_btn)

## 公開方法：顯示名人堂面板
func show_hall_of_fame():
	if is_instance_valid(_hall_of_fame_panel):
		_hall_of_fame_panel.show_panel()

# ---- 智慧推薦面板（DAY-110）----
var _recommend_panel = null

func _init_recommend_panel() -> void:
	var RecommendPanelScript = load("res://scripts/ui/RecommendPanel.gd")
	if RecommendPanelScript == null:
		return
	_recommend_panel = RecommendPanelScript.new()
	add_child(_recommend_panel)

	# 連接訊號
	GameManager.recommendations_received.connect(func(data):
		if is_instance_valid(_recommend_panel):
			_recommend_panel.update_recommendations(data)
	)

	# 在 TopBar 加入推薦按鈕
	var top_bar = get_node_or_null("TopBar")
	if is_instance_valid(top_bar):
		var rec_btn = Button.new()
		rec_btn.name = "RecommendButton"
		rec_btn.text = "💡"
		rec_btn.custom_minimum_size = Vector2(36, 30)
		rec_btn.position = Vector2(1140, 5)
		rec_btn.pressed.connect(func():
			if is_instance_valid(_recommend_panel):
				_recommend_panel.show_panel()
		)
		top_bar.add_child(rec_btn)

## ── 雙層倍率輪盤面板（DAY-113）──────────────────────────────────────────────
var _roulette_panel_node = null

## 初始化雙層倍率輪盤面板（DAY-113）
## 位置：畫面中央（全螢幕覆蓋式彈窗，z_index=72 在 WheelPanel 之上）
func _init_roulette_panel() -> void:
	var panel = RoulettePanelScript.new()
	panel.position = Vector2(640, 360)
	panel.z_index = 72
	add_child(panel)
	panel.setup(_pixel_font)
	_roulette_panel_node = panel

## ── Buy Bonus 面板（DAY-114）──────────────────────────────────────────────────
var _buy_bonus_panel_node = null
var _buy_bonus_btn: Button = null

## 初始化 Buy Bonus 面板（DAY-114）
## 位置：畫面中央（全螢幕覆蓋式彈窗，z_index=65）
## 同時在 BottomBar 加入 Buy Bonus 按鈕
func _init_buy_bonus_panel() -> void:
	# 建立面板
	var panel = BuyBonusPanelScript.new()
	panel.position = Vector2(640, 360)
	panel.z_index = 65
	add_child(panel)
	panel.setup(_pixel_font)
	_buy_bonus_panel_node = panel

	# 在 BottomBar 加入 Buy Bonus 按鈕（Bonus 按鈕旁邊）
	_buy_bonus_btn = Button.new()
	_buy_bonus_btn.text = "💰 Buy"
	_buy_bonus_btn.size = Vector2(70, 30)
	_buy_bonus_btn.position = Vector2(1190, 4)
	_buy_bonus_btn.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
	if is_instance_valid(_pixel_font):
		_buy_bonus_btn.add_theme_font_override("font", _pixel_font)
		_buy_bonus_btn.add_theme_font_size_override("font_size", 12)
	_buy_bonus_btn.pressed.connect(_on_buy_bonus_btn_pressed)
	add_child(_buy_bonus_btn)

func _on_buy_bonus_btn_pressed() -> void:
	if is_instance_valid(_buy_bonus_panel_node):
		_buy_bonus_panel_node.show_panel()

## Co-op Boss Raid 面板（DAY-115）
## 位置：全螢幕覆蓋，z_index=68（在 BuyBonus 之上）
func _init_raid_panel() -> void:
	var panel = RaidPanelScript.new()
	panel.z_index = 68
	add_child(panel)

	# 連接 GameManager 訊號
	if GameManager.has_signal("raid_warning"):
		GameManager.raid_warning.connect(func(data): panel.show_warning(data))
	if GameManager.has_signal("raid_started"):
		GameManager.raid_started.connect(func(data): panel.show_raid_start(data))
	if GameManager.has_signal("raid_updated"):
		GameManager.raid_updated.connect(func(data): panel.update_raid(data))
	if GameManager.has_signal("raid_result"):
		GameManager.raid_result.connect(func(data):
			var my_id = GameManager.get_player_id() if GameManager.has_method("get_player_id") else ""
			panel.show_result(data, my_id)
		)

func _init_fragment_panel() -> void:
	var panel = FragmentPanelScript.new()
	panel.z_index = 62
	add_child(panel)

	# 連接 GameManager 訊號
	if GameManager.has_signal("fragment_dropped"):
		GameManager.fragment_dropped.connect(func(data): panel.on_fragment_drop(data))
	if GameManager.has_signal("fragment_completed"):
		GameManager.fragment_completed.connect(func(data): panel.on_fragment_complete(data))
	if GameManager.has_signal("fragment_status_received"):
		GameManager.fragment_status_received.connect(func(data): panel.update_status(data))

# ---- 幸運捕獲系統（DAY-119）----
var _lucky_catch_panel_node = null

func _init_lucky_catch_panel() -> void:
	var panel = LuckyCatchPanelScript.new()
	panel.name = "LuckyCatchPanel"
	panel.z_index = 75
	add_child(panel)
	_lucky_catch_panel_node = panel

# ---- Rapid Respin 系統（DAY-121）----
var _rapid_respin_panel_node = null

func _init_rapid_respin_panel() -> void:
	var panel = RapidRespinPanelScript.new()
	panel.name = "RapidRespinPanel"
	panel.z_index = 71
	add_child(panel)
	_rapid_respin_panel_node = panel

# ---- 寶藏地圖系統（DAY-122）----
var _treasure_map_panel_node = null

func _init_treasure_map_panel() -> void:
	var panel = TreasureMapPanelScript.new()
	panel.name = "TreasureMapPanel"
	panel.z_index = 80
	add_child(panel)
	_treasure_map_panel_node = panel

func show_treasure_map_panel() -> void:
	if is_instance_valid(_treasure_map_panel_node):
		_treasure_map_panel_node.show_panel()

# ---- 閃電挑戰系統（DAY-123）----
var _flash_challenge_panel_node = null

func _init_flash_challenge_panel() -> void:
	var panel = FlashChallengePanelScript.new()
	panel.name = "FlashChallengePanel"
	panel.z_index = 73
	# 右側中間位置
	panel.position = Vector2(1280 - 330, 200)
	add_child(panel)
	_flash_challenge_panel_node = panel
	
	# 連接 GameManager 訊號
	var gm = get_node_or_null("/root/GameManager")
	if gm:
		if gm.has_signal("flash_challenge_started"):
			gm.flash_challenge_started.connect(_on_flash_challenge_started)
		if gm.has_signal("flash_challenge_updated"):
			gm.flash_challenge_updated.connect(_on_flash_challenge_updated)
		if gm.has_signal("flash_challenge_ended"):
			gm.flash_challenge_ended.connect(_on_flash_challenge_ended)
		if gm.has_signal("flash_challenge_reward"):
			gm.flash_challenge_reward.connect(_on_flash_challenge_reward)

func _on_flash_challenge_started(data: Dictionary) -> void:
	if is_instance_valid(_flash_challenge_panel_node):
		_flash_challenge_panel_node.on_flash_challenge_start(data)

func _on_flash_challenge_updated(data: Dictionary) -> void:
	if is_instance_valid(_flash_challenge_panel_node):
		_flash_challenge_panel_node.on_flash_challenge_update(data)

func _on_flash_challenge_ended(data: Dictionary) -> void:
	if is_instance_valid(_flash_challenge_panel_node):
		_flash_challenge_panel_node.on_flash_challenge_end(data)

func _on_flash_challenge_reward(data: Dictionary) -> void:
	if is_instance_valid(_flash_challenge_panel_node):
		_flash_challenge_panel_node.on_flash_challenge_reward(data)

# ---- 傳說目標警報系統（DAY-124）----

func _init_rare_target_alert() -> void:
	var gm = get_node_or_null("/root/GameManager")
	if gm and gm.has_signal("rare_target_alerted"):
		gm.rare_target_alerted.connect(_on_rare_target_alerted)

func _on_rare_target_alerted(data: Dictionary) -> void:
	var quality: String = data.get("quality", "epic")
	var message: String = data.get("message", "")
	var color_str: String = data.get("color", "#FFD700")
	var icon: String = data.get("icon", "⭐")
	
	# 顯示頂部橫幅警報
	_show_rare_target_banner(icon, message, Color(color_str), quality == "legendary")

func _show_rare_target_banner(icon: String, message: String, color: Color, is_legendary: bool) -> void:
	# 建立橫幅
	var banner := PanelContainer.new()
	banner.z_index = 95
	
	var style := StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.05, 0.1, 0.9)
	style.border_color = color
	style.border_width_top = 2
	style.border_width_bottom = 2
	style.corner_radius_top_left = 4
	style.corner_radius_top_right = 4
	style.corner_radius_bottom_left = 4
	style.corner_radius_bottom_right = 4
	banner.add_theme_stylebox_override("panel", style)
	
	var label := Label.new()
	label.text = icon + " " + message
	label.add_theme_font_size_override("font_size", 14)
	label.add_theme_color_override("font_color", color)
	label.add_theme_constant_override("outline_size", 2)
	label.add_theme_color_override("font_outline_color", Color.BLACK)
	banner.add_child(label)
	
	# 置中頂部
	banner.position = Vector2(640 - 200, -50)
	banner.size = Vector2(400, 36)
	add_child(banner)
	
	# 滑入動畫
	var tween := create_tween()
	tween.tween_property(banner, "position:y", 10, 0.3).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	
	# 傳說品質：金色閃爍
	if is_legendary:
		tween.parallel().tween_property(banner, "modulate", Color(2.0, 1.8, 0.5, 1.0), 0.15)
		tween.tween_property(banner, "modulate", Color.WHITE, 0.15)
		tween.tween_property(banner, "modulate", Color(2.0, 1.8, 0.5, 1.0), 0.15)
		tween.tween_property(banner, "modulate", Color.WHITE, 0.15)
	
	# 3 秒後滑出
	tween.tween_interval(3.0)
	tween.tween_property(banner, "position:y", -50, 0.3).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
	tween.tween_callback(banner.queue_free)

# ---- 黃金時間系統（DAY-125）----

const GoldenTimePanelScript = preload("res://scripts/ui/GoldenTimePanel.gd")
var _golden_time_panel: Control = null

func _init_golden_time_panel() -> void:
	var panel = GoldenTimePanelScript.new()
	panel.name = "GoldenTimePanel"
	panel.z_index = 76
	add_child(panel)
	_golden_time_panel = panel

	var gm = get_node_or_null("/root/GameManager")
	if gm:
		if gm.has_signal("golden_time_started"):
			gm.golden_time_started.connect(_on_golden_time_started)
		if gm.has_signal("golden_time_ended"):
			gm.golden_time_ended.connect(_on_golden_time_ended)
		if gm.has_signal("golden_time_status"):
			gm.golden_time_status.connect(_on_golden_time_status)

func _on_golden_time_started(data: Dictionary) -> void:
	if is_instance_valid(_golden_time_panel):
		_golden_time_panel.show_golden_time_start(data)

func _on_golden_time_ended(data: Dictionary) -> void:
	if is_instance_valid(_golden_time_panel):
		_golden_time_panel.show_golden_time_end(data)

func _on_golden_time_status(data: Dictionary) -> void:
	if is_instance_valid(_golden_time_panel):
		_golden_time_panel.update_status(data)

# ---- 稀有連擊累積倍率系統（DAY-126）----

const RareCatchPanelScript = preload("res://scripts/ui/RareCatchPanel.gd")
var _rare_catch_panel: Control = null

func _init_rare_catch_panel() -> void:
	var panel = RareCatchPanelScript.new()
	panel.name = "RareCatchPanel"
	panel.z_index = 74
	# 右下角，在 ActivityFeed 上方
	panel.position = Vector2(1280 - 200, 720 - 200)
	add_child(panel)
	_rare_catch_panel = panel

	var gm = get_node_or_null("/root/GameManager")
	if gm:
		if gm.has_signal("rare_catch_updated"):
			gm.rare_catch_updated.connect(_on_rare_catch_updated)
		if gm.has_signal("rare_catch_broadcasted"):
			gm.rare_catch_broadcasted.connect(_on_rare_catch_broadcasted)
		if gm.has_signal("rare_catch_reset"):
			gm.rare_catch_reset.connect(_on_rare_catch_reset_signal)

func _on_rare_catch_updated(data: Dictionary) -> void:
	if is_instance_valid(_rare_catch_panel):
		_rare_catch_panel.on_rare_catch_update(data)

func _on_rare_catch_broadcasted(data: Dictionary) -> void:
	# 全服廣播：顯示頂部小橫幅
	var player_name: String = data.get("player_name", "玩家")
	var icon: String = data.get("icon", "💎")
	var level_name: String = data.get("level_name", "稀有連擊")
	var color_str: String = data.get("color", "#00BFFF")
	var mult_boost: float = data.get("mult_boost", 5.0)
	var mult_str := "%.0f" % mult_boost
	var msg := icon + " " + player_name + " 達成 " + level_name + " ×" + mult_str + "！"
	_show_rare_target_banner(icon, msg, Color(color_str), false)

func _on_rare_catch_reset_signal(data: Dictionary) -> void:
	if is_instance_valid(_rare_catch_panel):
		_rare_catch_panel.on_rare_catch_reset(data)

# ---- 天氣湧現事件（DAY-127）----

const WeatherSurgePanelScript = preload("res://scripts/ui/WeatherSurgePanel.gd")
var _weather_surge_panel: Control = null

func _init_weather_surge_panel() -> void:
	var panel = WeatherSurgePanelScript.new()
	panel.name = "WeatherSurgePanel"
	panel.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 77  # 在黃金時間面板（76）之上
	add_child(panel)
	panel.setup(_pixel_font)
	_weather_surge_panel = panel

# ---- 龍怒蓄力大招面板（DAY-128）----

const DragonWrathPanelScript = preload("res://scripts/ui/DragonWrathPanel.gd")
var _dragon_wrath_panel: Control = null

func _init_dragon_wrath_panel() -> void:
	var panel = DragonWrathPanelScript.new()
	panel.name = "DragonWrathPanel"
	# 左下角，在特殊武器面板上方
	panel.position = Vector2(8, 720 - 90)
	panel.z_index = 9
	add_child(panel)
	panel.setup(_pixel_font)
	_dragon_wrath_panel = panel

# ---- 不死 BOSS 連勝面板（DAY-129）----
const ImmortalBossPanelScript = preload("res://scripts/ui/ImmortalBossPanel.gd")
var _immortal_boss_panel: Control = null

func _init_immortal_boss_panel() -> void:
	var panel = ImmortalBossPanelScript.new()
	panel.name = "ImmortalBossPanel"
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 70
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_immortal_boss_panel = panel

	# 連接 GameManager 訊號
	GameManager.immortal_boss_spawned.connect(func(data): panel.on_immortal_boss_spawn(data))
	GameManager.immortal_boss_hit.connect(func(data): panel.on_immortal_boss_hit(data))
	GameManager.immortal_boss_left.connect(func(data): panel.on_immortal_boss_leave(data))
	GameManager.immortal_boss_status.connect(func(data): panel.on_immortal_boss_status(data))

# ---- 覺醒 BOSS 面板（DAY-130）----
const AwakenBossPanelScript = preload("res://scripts/ui/AwakenBossPanel.gd")
var _awaken_boss_panel: Control = null

func _init_awaken_boss_panel() -> void:
	var panel = AwakenBossPanelScript.new()
	panel.name = "AwakenBossPanel"
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 69  # 比不死 BOSS 低一層（不死 BOSS z=70）
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_awaken_boss_panel = panel

	GameManager.awaken_boss_spawned.connect(func(data): panel.on_awaken_boss_spawn(data))
	GameManager.awaken_boss_hit.connect(func(data): panel.on_awaken_boss_hit(data))
	GameManager.awaken_boss_powerup.connect(func(data): panel.on_awaken_boss_powerup(data))
	GameManager.awaken_boss_left.connect(func(data): panel.on_awaken_boss_leave(data))
	GameManager.awaken_boss_status.connect(func(data): panel.on_awaken_boss_status(data))

# ---- 連勝獎勵面板（DAY-131）----
const WinStreakPanelScript = preload("res://scripts/ui/WinStreakPanel.gd")
var _win_streak_panel: Control = null

func _init_win_streak_panel() -> void:
	var panel = WinStreakPanelScript.new()
	panel.name = "WinStreakPanel"
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 67
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_win_streak_panel = panel

	GameManager.win_streak_updated.connect(func(data): panel.on_win_streak_update(data))
	GameManager.win_streak_milestone.connect(func(data): panel.on_win_streak_milestone(data))
	GameManager.win_streak_reset.connect(func(data): panel.on_win_streak_reset(data))

# ---- 閃電鰻連鎖攻擊面板（DAY-132）----
const LightningEelPanelScript = preload("res://scripts/ui/LightningEelPanel.gd")
var _lightning_eel_panel: Control = null

func _init_lightning_eel_panel() -> void:
	var panel = LightningEelPanelScript.new()
	panel.name = "LightningEelPanel"
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 66
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_lightning_eel_panel = panel

	# 設定玩家 ID（用於判斷是否為自己觸發）
	if GameManager.local_player_id != "":
		panel.set_player_id(GameManager.local_player_id)

	GameManager.lightning_eel_chain.connect(func(data): panel.show_chain_result(data))
	GameManager.lightning_eel_status.connect(func(data): panel.update_status(data))

# ---- 狂熱模式面板（DAY-133）----
const FeverModePanelScript = preload("res://scripts/ui/FeverModePanel.gd")
var _fever_mode_panel: Control = null

func _init_fever_mode_panel() -> void:
	var panel = FeverModePanelScript.new()
	panel.name = "FeverModePanel"
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 65
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_fever_mode_panel = panel

	if GameManager.local_player_id != "":
		panel.set_player_id(GameManager.local_player_id)

	GameManager.fever_mode_started.connect(func(data): panel.on_fever_start(data))
	GameManager.fever_mode_ended.connect(func(data): panel.on_fever_end(data))
	GameManager.fever_mode_status.connect(func(data): panel.on_fever_status(data))

# ---- 失敗補償面板（DAY-135）----
const UnluckyBonusPanelScript = preload("res://scripts/ui/UnluckyBonusPanel.gd")
var _unlucky_bonus_panel: Control = null

func _init_unlucky_bonus_panel() -> void:
	var panel = UnluckyBonusPanelScript.new()
	panel.name = "UnluckyBonusPanel"
	panel.z_index = 78
	add_child(panel)
	panel.setup(_pixel_font)
	_unlucky_bonus_panel = panel

# ---- 競速獵殺面板（DAY-136）----
const SpeedRacePanelScript = preload("res://scripts/ui/SpeedRacePanel.gd")
var _speed_race_panel: Control = null

func _init_speed_race_panel() -> void:
	var panel = SpeedRacePanelScript.new()
	panel.name = "SpeedRacePanel"
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 79  # 在失敗補償面板（78）之上
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_speed_race_panel = panel

	GameManager.speed_race_started.connect(func(data): panel.on_speed_race_start(data))
	GameManager.speed_race_ended.connect(func(data): panel.on_speed_race_end(data))
	GameManager.speed_race_cancelled.connect(func(data): panel.on_speed_race_cancel(data))
	GameManager.speed_race_result.connect(func(data): panel.on_speed_race_result(data))

# ---- 全服目標懸賞面板（DAY-137）----
const BountyPanelScript = preload("res://scripts/ui/BountyPanel.gd")
var _bounty_panel: Control = null

func _init_bounty_panel() -> void:
	var panel = BountyPanelScript.new()
	panel.name = "BountyPanel"
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 80  # 在競速面板（79）之上
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_bounty_panel = panel

	GameManager.bounty_posted.connect(func(data): panel.on_bounty_posted(data))
	GameManager.bounty_claimed.connect(func(data): panel.on_bounty_claimed(data))
	GameManager.bounty_killed.connect(func(data): panel.on_bounty_killed(data))
	GameManager.bounty_expired.connect(func(data): panel.on_bounty_expired(data))

# ---- 全服倍率風暴面板（DAY-138）----
const MultStormPanelScript = preload("res://scripts/ui/MultStormPanel.gd")
var _mult_storm_panel: Control = null

func _init_mult_storm_panel() -> void:
	var panel = MultStormPanelScript.new()
	panel.name = "MultStormPanel"
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 81  # 在懸賞面板（80）之上
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_mult_storm_panel = panel

	GameManager.mult_storm_started.connect(func(data): panel.on_storm_start(data))
	GameManager.mult_storm_ended.connect(func(data): panel.on_storm_end(data))

# ---- 雙環輪盤面板（DAY-139）----
const DualRoulettePanelScript = preload("res://scripts/ui/DualRoulettePanel.gd")
var _dual_roulette_panel: Control = null

func _init_dual_roulette_panel() -> void:
	var panel = DualRoulettePanelScript.new()
	panel.name = "DualRoulettePanel"
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 82  # 在倍率風暴面板（81）之上
	add_child(panel)
	_dual_roulette_panel = panel

	GameManager.dual_roulette_started.connect(func(data):
		panel.show_roulette(
			data.get("target_mult", 30.0),
			data.get("base_reward", 0),
			data.get("spin_duration", 3.0),
			data.get("inner_ring", [2.0, 3.0, 5.0, 8.0, 10.0]),
			data.get("outer_ring", [2.0, 3.0, 5.0, 7.0, 10.0, 15.0])
		)
	)
	GameManager.dual_roulette_result.connect(func(data):
		panel.show_result(
			data.get("inner_result", 1.0),
			data.get("outer_result", 1.0),
			data.get("combined", 1.0),
			data.get("bonus_reward", 0),
			data.get("new_balance", 0)
		)
	)

# ---- Mega Catch 事件面板（DAY-140）----
const MegaCatchPanelScript = preload("res://scripts/ui/MegaCatchPanel.gd")
var _mega_catch_panel: Control = null

func _init_mega_catch_panel() -> void:
	var panel = MegaCatchPanelScript.new()
	panel.name = "MegaCatchPanel"
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.z_index = 83  # 在雙環輪盤面板（82）之上
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_mega_catch_panel = panel

	GameManager.mega_catch_started.connect(func(data): panel.on_mega_catch_start(data))
	GameManager.mega_catch_ended.connect(func(data): panel.on_mega_catch_end(data))

# ---- 鑽頭龍蝦連帶效果面板（DAY-142）----

const DrillLobsterPanelScript = preload("res://scripts/ui/DrillLobsterPanel.gd")
var _drill_lobster_panel: Control = null

func _init_drill_lobster_panel() -> void:
	var panel = DrillLobsterPanelScript.new()
	panel.name = "DrillLobsterPanel"
	panel.layer = 84  # 在 MegaCatchPanel(83) 之上
	add_child(panel)
	panel.setup(_pixel_font)
	_drill_lobster_panel = panel

# ---- 炸彈蟹連環爆炸面板（DAY-143）----

const BombCrabPanelScript = preload("res://scripts/ui/BombCrabPanel.gd")
var _bomb_crab_panel: Control = null

func _init_bomb_crab_panel() -> void:
	var panel = BombCrabPanelScript.new()
	panel.name = "BombCrabPanel"
	panel.layer = 85  # 在 DrillLobsterPanel(84) 之上
	add_child(panel)
	panel.setup(_pixel_font)
	_bomb_crab_panel = panel

# ---- 巨型章魚轉盤面板（DAY-144）----

const MegaOctopusPanelScript = preload("res://scripts/ui/MegaOctopusPanel.gd")
var _mega_octopus_panel: Control = null

func _init_mega_octopus_panel() -> void:
	var panel = MegaOctopusPanelScript.new()
	panel.name = "MegaOctopusPanel"
	panel.layer = 89  # 在 BombCrabPanel(85) 之上，轉盤需要高層級
	add_child(panel)
	panel.setup(_pixel_font)
	_mega_octopus_panel = panel

# ---- 巨型鮟鱇魚電擊寶箱面板（DAY-145）----

const AnglerfishPanelScript = preload("res://scripts/ui/AnglerfishPanel.gd")
var _anglerfish_panel: Control = null

func _init_anglerfish_panel() -> void:
	var panel = AnglerfishPanelScript.new()
	panel.name = "AnglerfishPanel"
	panel.layer = 86  # 在 BombCrabPanel(85) 之上
	add_child(panel)
	panel.setup(_pixel_font)
	_anglerfish_panel = panel

# ---- 巨型鹹水鱷魚獵魚面板（DAY-146）----

const CrocodilePanelScript = preload("res://scripts/ui/CrocodilePanel.gd")
var _crocodile_panel: Control = null

func _init_crocodile_panel() -> void:
	var panel = CrocodilePanelScript.new()
	panel.name = "CrocodilePanel"
	panel.layer = 87  # 在 AnglerfishPanel(86) 之上
	add_child(panel)
	panel.setup(_pixel_font)
	_crocodile_panel = panel

# ---- 夢幻巨型獎勵魚面板（DAY-147）----

const GiantPrizeFishPanelScript = preload("res://scripts/ui/GiantPrizeFishPanel.gd")
var _giant_prize_fish_panel: Control = null

func _init_giant_prize_fish_panel() -> void:
	var panel = GiantPrizeFishPanelScript.new()
	panel.name = "GiantPrizeFishPanel"
	panel.z_index = 88  # 在 CrocodilePanel(87) 之上
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	panel.setup(_pixel_font)
	_giant_prize_fish_panel = panel

# ---- 千龍王強化輪盤面板（DAY-148）----

const ChainLongWheelPanelScript = preload("res://scripts/ui/ChainLongWheelPanel.gd")
var _chainlong_wheel_panel: Control = null

func _init_chainlong_wheel_panel() -> void:
	var panel = ChainLongWheelPanelScript.new()
	panel.name = "ChainLongWheelPanel"
	panel.z_index = 90  # 在 GiantPrizeFishPanel(88) 之上，是最高優先級的個人機制
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_chainlong_wheel_panel = panel

	# 連接千龍王輪盤訊號
	var gm = get_node_or_null("/root/GameManager")
	if gm:
		if gm.has_signal("chainlong_wheel_start"):
			gm.chainlong_wheel_start.connect(_on_chainlong_wheel_start)
		if gm.has_signal("chainlong_wheel_result"):
			gm.chainlong_wheel_result.connect(_on_chainlong_wheel_result)

func _on_chainlong_wheel_start(data: Dictionary) -> void:
	if not is_instance_valid(_chainlong_wheel_panel):
		return
	var my_id = ""
	var gm = get_node_or_null("/root/GameManager")
	if gm and gm.has_method("get_player_id"):
		my_id = gm.get_player_id()
	var killer_id: String = data.get("killer_id", "")
	var is_personal: bool = (killer_id == my_id)
	_chainlong_wheel_panel.show_wheel(
		data.get("killer_name", ""),
		data.get("target_mult", 500.0),
		data.get("base_reward", 0),
		data.get("spin_secs", 4.0),
		data.get("inner_slots", [5.0, 10.0, 20.0, 30.0, 50.0]),
		data.get("outer_slots", [2.0, 3.0, 5.0, 7.0, 10.0, 20.0]),
		is_personal
	)

func _on_chainlong_wheel_result(data: Dictionary) -> void:
	if not is_instance_valid(_chainlong_wheel_panel):
		return
	var my_id = ""
	var gm = get_node_or_null("/root/GameManager")
	if gm and gm.has_method("get_player_id"):
		my_id = gm.get_player_id()
	var killer_id: String = data.get("killer_id", "")
	var is_personal: bool = (killer_id == my_id)
	_chainlong_wheel_panel.show_result(
		data.get("inner_result", 5.0),
		data.get("outer_result", 2.0),
		data.get("combined", 10.0),
		data.get("bonus_reward", 0),
		data.get("new_balance", 0),
		data.get("is_mega_win", false),
		is_personal
	)

# ---- 黃金水母全場電擊面板（DAY-149）----

const GoldenJellyfishPanelScript = preload("res://scripts/ui/GoldenJellyfishPanel.gd")
var _golden_jellyfish_panel: Control = null

func _init_golden_jellyfish_panel() -> void:
	var panel = GoldenJellyfishPanelScript.new()
	panel.name = "GoldenJellyfishPanel"
	panel.z_index = 86  # 在 AnglerfishPanel(86) 同層，但黃金水母是不同系統
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_golden_jellyfish_panel = panel

	# 連接黃金水母電擊訊號
	var gm = get_node_or_null("/root/GameManager")
	if gm and gm.has_signal("golden_jellyfish_shock"):
		gm.golden_jellyfish_shock.connect(_on_golden_jellyfish_shock)

func _on_golden_jellyfish_shock(data: Dictionary) -> void:
	if not is_instance_valid(_golden_jellyfish_panel):
		return
	_golden_jellyfish_panel.handle_shock_event(data)

# ---- 雷霆龍蝦免費射擊面板（DAY-150）----

const ThunderboltLobsterPanelScript = preload("res://scripts/ui/ThunderboltLobsterPanel.gd")
var _thunderbolt_lobster_panel: Control = null

func _init_thunderbolt_lobster_panel() -> void:
	var panel = ThunderboltLobsterPanelScript.new()
	panel.name = "ThunderboltLobsterPanel"
	panel.z_index = 87  # 在 GoldenJellyfishPanel(86) 之上
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_thunderbolt_lobster_panel = panel

	# 連接雷霆龍蝦訊號
	var gm = get_node_or_null("/root/GameManager")
	if gm:
		if gm.has_signal("thunderbolt_lobster_activate"):
			gm.thunderbolt_lobster_activate.connect(_on_thunderbolt_lobster_activate)
		if gm.has_signal("thunderbolt_lobster_shot"):
			gm.thunderbolt_lobster_shot.connect(_on_thunderbolt_lobster_shot)
		if gm.has_signal("thunderbolt_lobster_end"):
			gm.thunderbolt_lobster_end.connect(_on_thunderbolt_lobster_end)

func _on_thunderbolt_lobster_activate(data: Dictionary) -> void:
	if not is_instance_valid(_thunderbolt_lobster_panel):
		return
	_thunderbolt_lobster_panel.handle_activate(data)

func _on_thunderbolt_lobster_shot(data: Dictionary) -> void:
	if not is_instance_valid(_thunderbolt_lobster_panel):
		return
	_thunderbolt_lobster_panel.handle_shot(data)

func _on_thunderbolt_lobster_end(data: Dictionary) -> void:
	if not is_instance_valid(_thunderbolt_lobster_panel):
		return
	_thunderbolt_lobster_panel.handle_end(data)

# ---- 彩虹鳳凰 Power Up 面板（DAY-151）----

const RainbowPhoenixPanelScript = preload("res://scripts/ui/RainbowPhoenixPanel.gd")
var _rainbow_phoenix_panel: Control = null

func _init_rainbow_phoenix_panel() -> void:
	var panel = RainbowPhoenixPanelScript.new()
	panel.name = "RainbowPhoenixPanel"
	panel.z_index = 88  # 在 ThunderboltLobsterPanel(87) 之上
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_rainbow_phoenix_panel = panel

	# 連接彩虹鳳凰訊號
	var gm = get_node_or_null("/root/GameManager")
	if gm:
		if gm.has_signal("rainbow_phoenix_activate"):
			gm.rainbow_phoenix_activate.connect(_on_rainbow_phoenix_activate)
		if gm.has_signal("rainbow_phoenix_end"):
			gm.rainbow_phoenix_end.connect(_on_rainbow_phoenix_end)

func _on_rainbow_phoenix_activate(data: Dictionary) -> void:
	if not is_instance_valid(_rainbow_phoenix_panel):
		return
	_rainbow_phoenix_panel.handle_activate(data)

func _on_rainbow_phoenix_end(data: Dictionary) -> void:
	if not is_instance_valid(_rainbow_phoenix_panel):
		return
	_rainbow_phoenix_panel.handle_end(data)

# ---- 吸血鬼成長倍率面板（DAY-152）----

const VampirePanelScript = preload("res://scripts/ui/VampirePanel.gd")
var _vampire_panel: Control = null

func _init_vampire_panel() -> void:
	var panel = VampirePanelScript.new()
	panel.name = "VampirePanel"
	panel.z_index = 85  # 在 BombCrabPanel(85) 同層，吸血鬼是不同系統
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_vampire_panel = panel

	# 連接吸血鬼訊號
	var gm = get_node_or_null("/root/GameManager")
	if gm:
		if gm.has_signal("vampire_grow"):
			gm.vampire_grow.connect(_on_vampire_grow)
		if gm.has_signal("vampire_blood_moon"):
			gm.vampire_blood_moon.connect(_on_vampire_blood_moon)
		if gm.has_signal("vampire_killed"):
			gm.vampire_killed.connect(_on_vampire_killed)

func _on_vampire_grow(data: Dictionary) -> void:
	if not is_instance_valid(_vampire_panel):
		return
	_vampire_panel.handle_grow(data)

func _on_vampire_blood_moon(data: Dictionary) -> void:
	if not is_instance_valid(_vampire_panel):
		return
	_vampire_panel.handle_blood_moon(data)

func _on_vampire_killed(data: Dictionary) -> void:
	if not is_instance_valid(_vampire_panel):
		return
	_vampire_panel.handle_killed(data)

# ---- 水晶龍收集大獎面板（DAY-153）----

const CrystalDragonPanelScript = preload("res://scripts/ui/CrystalDragonPanel.gd")
var _crystal_dragon_panel: Control = null

func _init_crystal_dragon_panel() -> void:
	var panel = CrystalDragonPanelScript.new()
	panel.name = "CrystalDragonPanel"
	panel.z_index = 84  # 在 BombCrabPanel(85) 下方，水晶龍是常駐進度條
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	_crystal_dragon_panel = panel

	# 連接水晶龍訊號
	var gm = get_node_or_null("/root/GameManager")
	if gm:
		if gm.has_signal("crystal_dragon_drop"):
			gm.crystal_dragon_drop.connect(_on_crystal_dragon_drop)
		if gm.has_signal("crystal_dragon_reward"):
			gm.crystal_dragon_reward.connect(_on_crystal_dragon_reward)
		if gm.has_signal("crystal_dragon_status"):
			gm.crystal_dragon_status.connect(_on_crystal_dragon_status)

func _on_crystal_dragon_drop(data: Dictionary) -> void:
	if not is_instance_valid(_crystal_dragon_panel):
		return
	_crystal_dragon_panel.on_crystal_drop(
		data.get("killer_name", ""),
		data.get("crystals_gain", 0),
		data.get("total_crystals", 0),
		data.get("goal", 50),
		data.get("progress", 0.0)
	)

func _on_crystal_dragon_reward(data: Dictionary) -> void:
	if not is_instance_valid(_crystal_dragon_panel):
		return
	_crystal_dragon_panel.on_crystal_reward(
		data.get("contributors", []),
		data.get("total_reward", 0),
		data.get("message", "")
	)

func _on_crystal_dragon_status(data: Dictionary) -> void:
	if not is_instance_valid(_crystal_dragon_panel):
		return
	_crystal_dragon_panel.update_status(
		data.get("total_crystals", 0),
		data.get("goal", 50),
		data.get("progress", 0.0)
	)

# ---- 皇家閃電鰻持續連鎖電擊面板（DAY-156）----

const RoyalChainLightningPanelScript = preload("res://scripts/ui/RoyalChainLightningPanel.gd")
var _royal_chain_lightning_panel: Control = null

func _init_royal_chain_lightning_panel() -> void:
	var panel = RoyalChainLightningPanelScript.new()
	panel.name = "RoyalChainLightningPanel"
	panel.z_index = 83  # 在 CrystalDragonPanel(84) 下方
	panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	if panel.has_method("setup"):
		panel.setup(_pixel_font)
	_royal_chain_lightning_panel = panel

# ---- 黃金海龜時間停止面板（DAY-159）----

const GoldenTurtlePanelScript = preload("res://scripts/ui/GoldenTurtlePanel.gd")
var _golden_turtle_panel: Control = null

func _init_golden_turtle_panel() -> void:
	var panel = GoldenTurtlePanelScript.new()
	panel.name = "GoldenTurtlePanel"
	panel.z_index = 82  # 在 RoyalChainLightningPanel(83) 下方
	panel.position = Vector2(640, 360)  # 畫面中心（面板內部用相對座標）
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	if panel.has_method("setup"):
		panel.setup(_pixel_font)
	_golden_turtle_panel = panel

# ---- 幸運星魚全場倍率翻倍面板（DAY-160）----

const LuckyStarFishPanelScript = preload("res://scripts/ui/LuckyStarFishPanel.gd")
var _lucky_star_fish_panel: Control = null

func _init_lucky_star_fish_panel() -> void:
	var panel = LuckyStarFishPanelScript.new()
	panel.name = "LuckyStarFishPanel"
	panel.z_index = 81  # 在 GoldenTurtlePanel(82) 下方
	panel.position = Vector2(640, 360)  # 畫面中心
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	if panel.has_method("setup"):
		panel.setup(_pixel_font)
	_lucky_star_fish_panel = panel

# ---- 黃金鯊魚全服狂暴模式面板（DAY-161）----

const GoldenSharkPanelScript = preload("res://scripts/ui/GoldenSharkPanel.gd")
var _golden_shark_panel: Control = null

func _init_golden_shark_panel() -> void:
	var panel = GoldenSharkPanelScript.new()
	panel.name = "GoldenSharkPanel"
	panel.z_index = 80  # 在 LuckyStarFishPanel(81) 下方
	panel.position = Vector2(640, 360)  # 畫面中心
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	if panel.has_method("setup"):
		panel.setup(_pixel_font)
	_golden_shark_panel = panel

# ---- 金幣魚王即時獎勵面板（DAY-162）----

const MoneyFishPanelScript = preload("res://scripts/ui/MoneyFishPanel.gd")
var _money_fish_panel: Control = null

func _init_money_fish_panel() -> void:
	var panel = MoneyFishPanelScript.new()
	panel.name = "MoneyFishPanel"
	panel.z_index = 79  # 在 GoldenSharkPanel(80) 下方
	panel.position = Vector2(640, 360)  # 畫面中心
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	if panel.has_method("setup"):
		panel.setup(_pixel_font)
	_money_fish_panel = panel

# ---- 船長魚全服競速模式面板（DAY-163）----

const CaptainFishPanelScript = preload("res://scripts/ui/CaptainFishPanel.gd")
var _captain_fish_panel: Control = null

func _init_captain_fish_panel() -> void:
	var panel = CaptainFishPanelScript.new()
	panel.name = "CaptainFishPanel"
	panel.z_index = 78  # 在 MoneyFishPanel(79) 下方
	panel.position = Vector2(640, 360)  # 畫面中心
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	if panel.has_method("setup"):
		panel.setup(_pixel_font)
	_captain_fish_panel = panel

# ---- 深淵巨鯨全服 Boss 挑戰面板（DAY-164）----

const AbyssWhalePanelScript = preload("res://scripts/ui/AbyssWhalePanel.gd")
var _abyss_whale_panel: Control = null

func _init_abyss_whale_panel() -> void:
	var panel = AbyssWhalePanelScript.new()
	panel.name = "AbyssWhalePanel"
	panel.z_index = 77  # 在 CaptainFishPanel(78) 下方
	panel.position = Vector2(640, 360)  # 畫面中心
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	if panel.has_method("setup"):
		panel.setup(_pixel_font)
	_abyss_whale_panel = panel

# ---- 黑洞漩渦武器視覺效果面板（DAY-166）----

const BlackHolePanelScript = preload("res://scripts/ui/BlackHolePanel.gd")
var _black_hole_panel: Control = null

func _init_black_hole_panel() -> void:
	var panel = BlackHolePanelScript.new()
	panel.name = "BlackHolePanel"
	panel.z_index = 76  # 在 AbyssWhalePanel(77) 下方
	panel.position = Vector2(0, 0)  # 左上角原點（面板內部用絕對座標）
	add_child(panel)
	_black_hole_panel = panel

# ---- 黃金輪盤螃蟹面板（DAY-167）----

const RouletteCrabPanelScript = preload("res://scripts/ui/RouletteCrabPanel.gd")
var _roulette_crab_panel: Control = null

func _init_roulette_crab_panel() -> void:
	var panel = RouletteCrabPanelScript.new()
	panel.name = "RouletteCrabPanel"
	panel.z_index = 75  # 在 BlackHolePanel(76) 下方
	panel.position = Vector2(0, 0)  # 左上角原點（面板內部用絕對座標）
	add_child(panel)
	_roulette_crab_panel = panel

## ---- 獅子舞大獎爆發面板（DAY-168）----
const LionDancePanelScript = preload("res://scripts/ui/LionDancePanel.gd")
var _lion_dance_panel: Control = null

func _init_lion_dance_panel() -> void:
	var panel = LionDancePanelScript.new()
	panel.name = "LionDancePanel"
	panel.z_index = 74  # 在 RouletteCrabPanel(75) 下方
	panel.position = Vector2(0, 0)  # 左上角原點（面板內部用絕對座標）
	add_child(panel)
	_lion_dance_panel = panel

## ---- 漩渦魚群吸引面板（DAY-169）----
const VortexFishPanelScript = preload("res://scripts/ui/VortexFishPanel.gd")
var _vortex_fish_panel: Control = null

func _init_vortex_fish_panel() -> void:
	var panel = VortexFishPanelScript.new()
	panel.name = "VortexFishPanel"
	panel.z_index = 73  # 在 LionDancePanel(74) 下方
	panel.position = Vector2(0, 0)  # 左上角原點（面板內部用絕對座標）
	add_child(panel)
	_vortex_fish_panel = panel

## ---- 冰凍炸彈魚面板（DAY-170）----
const FreezeBombPanelScript = preload("res://scripts/ui/FreezeBombPanel.gd")
var _freeze_bomb_panel: Control = null

func _init_freeze_bomb_panel() -> void:
	var panel = FreezeBombPanelScript.new()
	panel.name = "FreezeBombPanel"
	panel.z_index = 72  # 在 VortexFishPanel(73) 下方
	panel.position = Vector2(0, 0)  # 左上角原點（面板內部用絕對座標）
	add_child(panel)
	_freeze_bomb_panel = panel

## ---- 冰釣幸運輪盤面板（DAY-171）----
const IceFishingPanelScript = preload("res://scripts/ui/IceFishingPanel.gd")
var _ice_fishing_panel: Control = null

func _init_ice_fishing_panel() -> void:
	var panel = IceFishingPanelScript.new()
	panel.name = "IceFishingPanel"
	panel.z_index = 71  # 在 FreezeBombPanel(72) 下方
	panel.position = Vector2(0, 0)  # 左上角原點（面板內部用絕對座標）
	add_child(panel)
	_ice_fishing_panel = panel

## ---- 幸運彩蛋魚面板（DAY-172）----
const LuckyEggPanelScript = preload("res://scripts/ui/LuckyEggPanel.gd")
var _lucky_egg_panel: Control = null

func _init_lucky_egg_panel() -> void:
	var panel = LuckyEggPanelScript.new()
	panel.name = "LuckyEggPanel"
	panel.z_index = 70  # 在 IceFishingPanel(71) 下方
	panel.position = Vector2(0, 0)  # 左上角原點（面板內部用絕對座標）
	add_child(panel)
	_lucky_egg_panel = panel

## ---- 彩虹幸運魚面板（DAY-173）----
const RainbowLuckyPanelScript = preload("res://scripts/ui/RainbowLuckyPanel.gd")
var _rainbow_lucky_panel: Control = null

func _init_rainbow_lucky_panel() -> void:
	var panel = RainbowLuckyPanelScript.new()
	panel.name = "RainbowLuckyPanel"
	panel.z_index = 69  # 在 LuckyEggPanel(70) 下方
	panel.position = Vector2(0, 0)  # 左上角原點（面板內部用絕對座標）
	add_child(panel)
	_rainbow_lucky_panel = panel

## ---- 海葵觸手攻擊面板（DAY-174）----
const SeaAnemonePanelScript = preload("res://scripts/ui/SeaAnemonePanel.gd")
var _sea_anemone_panel: Control = null

func _init_sea_anemone_panel() -> void:
	var panel = SeaAnemonePanelScript.new()
	panel.name = "SeaAnemonePanel"
	panel.z_index = 68  # 在 RainbowLuckyPanel(69) 下方
	panel.position = Vector2(0, 0)  # 左上角原點（面板內部用絕對座標）
	add_child(panel)
	_sea_anemone_panel = panel

## ---- 幸運骰子魚面板（DAY-175）----
const LuckyDicePanelScript = preload("res://scripts/ui/LuckyDicePanel.gd")
var _lucky_dice_panel: Control = null

func _init_lucky_dice_panel() -> void:
	var panel = LuckyDicePanelScript.new()
	panel.name = "LuckyDicePanel"
	panel.z_index = 67  # 在 SeaAnemonePanel(68) 下方
	panel.position = Vector2(0, 0)  # 左上角原點（面板內部用絕對座標）
	add_child(panel)
	_lucky_dice_panel = panel

## ---- 火焰風暴魚面板（DAY-176）----
const FireStormPanelScript = preload("res://scripts/ui/FireStormPanel.gd")
var _fire_storm_panel: Control = null

func _init_fire_storm_panel() -> void:
	var panel = FireStormPanelScript.new()
	panel.name = "FireStormPanel"
	panel.z_index = 68  # 與 SeaAnemonePanel 同層，但 fire_storm 是全服事件
	panel.position = Vector2(0, 0)
	add_child(panel)
	_fire_storm_panel = panel
	# 連接 GameManager 訊號
	if GameManager.has_signal("fire_storm_fish"):
		GameManager.fire_storm_fish.connect(func(data):
			if is_instance_valid(_fire_storm_panel):
				_fire_storm_panel.handle_fire_storm(data)
		)

## ---- 黃金寶藏魚面板（DAY-177）----
const GoldenTreasurePanelScript = preload("res://scripts/ui/GoldenTreasurePanel.gd")
var _golden_treasure_panel: Control = null

func _init_golden_treasure_panel() -> void:
	var panel = GoldenTreasurePanelScript.new()
	panel.name = "GoldenTreasurePanel"
	panel.z_index = 67  # 在 FireStormPanel(68) 下方
	panel.position = Vector2(0, 0)
	add_child(panel)
	_golden_treasure_panel = panel
	# 連接 GameManager 訊號
	if GameManager.has_signal("golden_treasure_fish"):
		GameManager.golden_treasure_fish.connect(func(data):
			if is_instance_valid(_golden_treasure_panel):
				_golden_treasure_panel.handle_golden_treasure(data)
		)

## ---- 美人魚治癒面板（DAY-178）----
const MermaidHealingPanelScript = preload("res://scripts/ui/MermaidHealingPanel.gd")
var _mermaid_healing_panel: Control = null

func _init_mermaid_healing_panel() -> void:
	var panel = MermaidHealingPanelScript.new()
	panel.name = "MermaidHealingPanel"
	panel.z_index = 66
	panel.position = Vector2(0, 0)
	add_child(panel)
	_mermaid_healing_panel = panel
	if GameManager.has_signal("mermaid_healing"):
		GameManager.mermaid_healing.connect(func(data):
			if is_instance_valid(_mermaid_healing_panel):
				_mermaid_healing_panel.handle_mermaid_healing(data)
		)

## ---- 幸運草魚面板（DAY-179）----
const LuckyCloverPanelScript = preload("res://scripts/ui/LuckyCloverPanel.gd")
var _lucky_clover_panel: Control = null

func _init_lucky_clover_panel() -> void:
	var panel = LuckyCloverPanelScript.new()
	panel.name = "LuckyCloverPanel"
	panel.z_index = 65
	panel.position = Vector2(0, 0)
	add_child(panel)
	_lucky_clover_panel = panel
	if GameManager.has_signal("lucky_clover_fish"):
		GameManager.lucky_clover_fish.connect(func(data):
			if is_instance_valid(_lucky_clover_panel):
				_lucky_clover_panel.handle_lucky_clover(data)
		)

## ---- 彩虹鯊魚爆發面板（DAY-180）----
const RainbowSharkPanelScript = preload("res://scripts/ui/RainbowSharkPanel.gd")
var _rainbow_shark_panel: Control = null

func _init_rainbow_shark_panel() -> void:
	var panel = RainbowSharkPanelScript.new()
	panel.name = "RainbowSharkPanel"
	panel.z_index = 66
	panel.position = Vector2(0, 0)
	add_child(panel)
	_rainbow_shark_panel = panel
	if GameManager.has_signal("rainbow_shark_burst"):
		GameManager.rainbow_shark_burst.connect(func(data):
			if is_instance_valid(_rainbow_shark_panel):
				_rainbow_shark_panel.handle_rainbow_shark(data)
		)

## ---- 雷霆鯊魚連鎖閃電面板（DAY-181）----
const ThunderSharkPanelScript = preload("res://scripts/ui/ThunderSharkPanel.gd")
var _thunder_shark_panel: Control = null

func _init_thunder_shark_panel() -> void:
	var panel = ThunderSharkPanelScript.new()
	panel.name = "ThunderSharkPanel"
	panel.z_index = 64
	panel.position = Vector2(0, 0)
	add_child(panel)
	_thunder_shark_panel = panel
	if GameManager.has_signal("thunder_shark_chain"):
		GameManager.thunder_shark_chain.connect(func(data):
			if is_instance_valid(_thunder_shark_panel):
				_thunder_shark_panel.handle_thunder_shark(data)
		)

## ---- 吸血鬼魚累積倍率面板（DAY-182）----
const VampireFishPanelScript = preload("res://scripts/ui/VampireFishPanel.gd")
var _vampire_fish_panel: Control = null

func _init_vampire_fish_panel() -> void:
	var panel = VampireFishPanelScript.new()
	panel.name = "VampireFishPanel"
	panel.z_index = 63
	panel.position = Vector2(0, 0)
	add_child(panel)
	_vampire_fish_panel = panel
	if GameManager.has_signal("vampire_fish"):
		GameManager.vampire_fish.connect(func(data):
			if is_instance_valid(_vampire_fish_panel):
				var my_id: String = GameManager.get("_player_id") if GameManager.get("_player_id") != null else ""
				_vampire_fish_panel.handle_vampire_fish(data, my_id)
		)

## ---- 閃電魚自動連鎖面板（DAY-183）----
const LightningAutoChainPanelScript = preload("res://scripts/ui/LightningAutoChainPanel.gd")
var _lightning_auto_chain_panel: Control = null

func _init_lightning_auto_chain_panel() -> void:
	var panel = LightningAutoChainPanelScript.new()
	panel.name = "LightningAutoChainPanel"
	panel.z_index = 62
	panel.position = Vector2(0, 0)
	add_child(panel)
	_lightning_auto_chain_panel = panel
	if GameManager.has_signal("lightning_auto_chain"):
		GameManager.lightning_auto_chain.connect(func(data):
			if is_instance_valid(_lightning_auto_chain_panel):
				_lightning_auto_chain_panel.handle_lightning_auto_chain(data)
		)

## ---- 隕石魚隕石雨面板（DAY-184）----
const MeteorFishPanelScript = preload("res://scripts/ui/MeteorFishPanel.gd")
var _meteor_fish_panel: Control = null

func _init_meteor_fish_panel() -> void:
	var panel = MeteorFishPanelScript.new()
	panel.name = "MeteorFishPanel"
	panel.z_index = 61
	panel.position = Vector2(0, 0)
	add_child(panel)
	_meteor_fish_panel = panel
	if GameManager.has_signal("meteor_fish"):
		GameManager.meteor_fish.connect(func(data):
			if is_instance_valid(_meteor_fish_panel):
				_meteor_fish_panel.handle_meteor_fish(data)
		)

## ---- 鳳凰魚涅槃重生面板（DAY-185）----
const PhoenixFishPanelScript = preload("res://scripts/ui/PhoenixFishPanel.gd")
var _phoenix_fish_panel: Control = null

func _init_phoenix_fish_panel() -> void:
	var panel = PhoenixFishPanelScript.new()
	panel.name = "PhoenixFishPanel"
	panel.z_index = 60
	panel.position = Vector2(0, 0)
	add_child(panel)
	_phoenix_fish_panel = panel
	if GameManager.has_signal("phoenix_fish"):
		GameManager.phoenix_fish.connect(func(data):
			if is_instance_valid(_phoenix_fish_panel):
				_phoenix_fish_panel.handle_phoenix_fish(data)
		)

## ---- 龍龜不死 Boss 面板（DAY-186）----
const DragonTurtlePanelScript = preload("res://scripts/ui/DragonTurtlePanel.gd")
var _dragon_turtle_panel: Control = null

func _init_dragon_turtle_panel() -> void:
	var panel = DragonTurtlePanelScript.new()
	panel.name = "DragonTurtlePanel"
	panel.z_index = 59
	panel.position = Vector2(0, 0)
	add_child(panel)
	_dragon_turtle_panel = panel
	if GameManager.has_signal("dragon_turtle"):
		GameManager.dragon_turtle.connect(func(data):
			if is_instance_valid(_dragon_turtle_panel):
				_dragon_turtle_panel.handle_dragon_turtle(data)
		)

## ---- 連鎖爆炸魚面板（DAY-187）----
const ChainBombPanelScript = preload("res://scripts/ui/ChainBombPanel.gd")
var _chain_bomb_panel: Control = null

func _init_chain_bomb_panel() -> void:
	var panel = ChainBombPanelScript.new()
	panel.name = "ChainBombPanel"
	panel.z_index = 58
	panel.position = Vector2(0, 0)
	add_child(panel)
	_chain_bomb_panel = panel
	if GameManager.has_signal("chain_bomb"):
		GameManager.chain_bomb.connect(func(data):
			if is_instance_valid(_chain_bomb_panel):
				_chain_bomb_panel.handle_chain_bomb(data)
		)

## ---- 巨型鱷魚獵食面板（DAY-188）----
const CrocodileHunterPanelScript = preload("res://scripts/ui/CrocodileHunterPanel.gd")
var _crocodile_hunter_panel: Control = null

func _init_crocodile_hunter_panel() -> void:
	var panel = CrocodileHunterPanelScript.new()
	panel.name = "CrocodileHunterPanel"
	panel.z_index = 57
	panel.position = Vector2(0, 0)
	add_child(panel)
	_crocodile_hunter_panel = panel
	if GameManager.has_signal("crocodile_hunter"):
		GameManager.crocodile_hunter.connect(func(data):
			if is_instance_valid(_crocodile_hunter_panel):
				_crocodile_hunter_panel.handle(data)
		)

## ---- 時間炸彈魚面板（DAY-189）----
const TimeBombFishPanelScript = preload("res://scripts/ui/TimeBombFishPanel.gd")
var _time_bomb_fish_panel: Control = null

func _init_time_bomb_fish_panel() -> void:
	var panel = TimeBombFishPanelScript.new()
	panel.name = "TimeBombFishPanel"
	panel.z_index = 56
	panel.position = Vector2(0, 0)
	add_child(panel)
	_time_bomb_fish_panel = panel
	if GameManager.has_signal("time_bomb_fish"):
		GameManager.time_bomb_fish.connect(func(data):
			if is_instance_valid(_time_bomb_fish_panel):
				_time_bomb_fish_panel.handle(data)
		)

## ---- 三重幸運魚面板（DAY-190）----
const TripleLuckyFishPanelScript = preload("res://scripts/ui/TripleLuckyFishPanel.gd")
var _triple_lucky_fish_panel: Control = null

func _init_triple_lucky_fish_panel() -> void:
	var panel = TripleLuckyFishPanelScript.new()
	panel.name = "TripleLuckyFishPanel"
	panel.z_index = 55  # 在 TimeBombFishPanel(56) 之下
	panel.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(panel)
	_triple_lucky_fish_panel = panel

## ---- 魚群驚嚇連帶面板（DAY-191）----
const SchoolPanicPanelScript = preload("res://scripts/ui/SchoolPanicPanel.gd")
var _school_panic_panel: Control = null

func _init_school_panic_panel() -> void:
	var panel = SchoolPanicPanelScript.new()
	panel.name = "SchoolPanicPanel"
	panel.z_index = 54  # 在 TripleLuckyFishPanel(55) 之下
	panel.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(panel)
	_school_panic_panel = panel

## ---- 搖滾骷髏演唱會面板（DAY-192）----
const RockSkeletonConcertPanelScript = preload("res://scripts/ui/RockSkeletonConcertPanel.gd")
var _rock_skeleton_concert_panel: Control = null

func _init_rock_skeleton_concert_panel() -> void:
	var panel = RockSkeletonConcertPanelScript.new()
	panel.name = "RockSkeletonConcertPanel"
	panel.z_index = 53  # 在 SchoolPanicPanel(54) 之下
	panel.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(panel)
	_rock_skeleton_concert_panel = panel

## ---- 電流水母電流網路面板（DAY-193）----
const ElectricJellyfishPanelScript = preload("res://scripts/ui/ElectricJellyfishPanel.gd")
var _electric_jellyfish_panel: Control = null

func _init_electric_jellyfish_panel() -> void:
	var panel = ElectricJellyfishPanelScript.new()
	panel.name = "ElectricJellyfishPanel"
	panel.z_index = 52  # 在 RockSkeletonConcertPanel(53) 之下
	panel.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(panel)
	_electric_jellyfish_panel = panel

## ---- 長龍王雙環輪盤面板（DAY-194）----
const ChainLongKingPanelScript = preload("res://scripts/ui/ChainLongKingPanel.gd")
var _chainlong_king_panel = null

func _init_chainlong_king_panel() -> void:
	var panel = ChainLongKingPanelScript.new()
	panel.name = "ChainLongKingPanel"
	panel.layer = 51  # 在 ElectricJellyfishPanel(52) 之下
	add_child(panel)
	_chainlong_king_panel = panel
	# 連接停止按鈕訊號
	if panel.has_signal("chainlong_king_stop_pressed"):
		panel.chainlong_king_stop_pressed.connect(_on_chainlong_king_stop_pressed)
	# 連接 GameManager 訊號
	if GameManager.has_signal("chainlong_king"):
		GameManager.chainlong_king.connect(_on_chainlong_king)

func _on_chainlong_king(data: Dictionary) -> void:
	if is_instance_valid(_chainlong_king_panel):
		_chainlong_king_panel.handle_chainlong_king(data)

func _on_chainlong_king_stop_pressed(instance_id: String) -> void:
	GameManager.send_chainlong_king_stop(instance_id)

## ---- 鑽頭龍蝦穿透爆炸面板（DAY-195）----
const DrillLobsterPanelScript = preload("res://scripts/ui/DrillLobsterPanel.gd")
var _drill_lobster_panel = null

func _init_drill_lobster_panel() -> void:
	var panel = DrillLobsterPanelScript.new()
	panel.name = "DrillLobsterPanel"
	panel.layer = 50
	add_child(panel)
	_drill_lobster_panel = panel
	if GameManager.has_signal("drill_lobster"):
		GameManager.drill_lobster.connect(_on_drill_lobster)

func _on_drill_lobster(data: Dictionary) -> void:
	if is_instance_valid(_drill_lobster_panel):
		_drill_lobster_panel.handle_drill_lobster(data)

## ---- 巨型鮟鱇魚電擊寶箱面板（DAY-196）----
const AnglerfishElectricPanelScript = preload("res://scripts/ui/AnglerfishElectricPanel.gd")
var _anglerfish_electric_panel = null

func _init_anglerfish_electric_panel() -> void:
	var panel = AnglerfishElectricPanelScript.new()
	panel.name = "AnglerfishElectricPanel"
	panel.layer = 49
	add_child(panel)
	_anglerfish_electric_panel = panel
	if GameManager.has_signal("anglerfish_electric"):
		GameManager.anglerfish_electric.connect(_on_anglerfish_electric)

func _on_anglerfish_electric(data: Dictionary) -> void:
	if is_instance_valid(_anglerfish_electric_panel):
		_anglerfish_electric_panel.handle_anglerfish_electric(data)

## ---- 神秘龍魚八波攻擊面板（DAY-197）----
const MysticDragonPanelScript = preload("res://scripts/ui/MysticDragonPanel.gd")
var _mystic_dragon_panel = null

func _init_mystic_dragon_panel() -> void:
	var panel = MysticDragonPanelScript.new()
	panel.name = "MysticDragonPanel"
	panel.layer = 48
	add_child(panel)
	_mystic_dragon_panel = panel
	if GameManager.has_signal("mystic_dragon"):
		GameManager.mystic_dragon.connect(_on_mystic_dragon)

func _on_mystic_dragon(data: Dictionary) -> void:
	if is_instance_valid(_mystic_dragon_panel):
		_mystic_dragon_panel.handle_mystic_dragon(data)

## ---- 幽靈魚分身面板（DAY-198）----
const GhostFishPanelScript = preload("res://scripts/ui/GhostFishPanel.gd")
var _ghost_fish_panel = null

func _init_ghost_fish_panel() -> void:
	var panel = GhostFishPanelScript.new()
	panel.name = "GhostFishPanel"
	panel.layer = 47
	add_child(panel)
	_ghost_fish_panel = panel
	if GameManager.has_signal("ghost_fish"):
		GameManager.ghost_fish.connect(_on_ghost_fish)

func _on_ghost_fish(data: Dictionary) -> void:
	if is_instance_valid(_ghost_fish_panel):
		_ghost_fish_panel.handle_ghost_fish(data)

## ---- 雷霆龍蝦免費射擊面板（DAY-199）----
const ThunderboltLobsterPanelScript = preload("res://scripts/ui/ThunderboltLobsterPanel.gd")
var _thunderbolt_lobster_panel = null

func _init_thunderbolt_lobster_panel() -> void:
	var panel = ThunderboltLobsterPanelScript.new()
	panel.name = "ThunderboltLobsterPanel"
	panel.layer = 46
	add_child(panel)
	_thunderbolt_lobster_panel = panel
	if GameManager.has_signal("thunderbolt_lobster"):
		GameManager.thunderbolt_lobster.connect(_on_thunderbolt_lobster)

func _on_thunderbolt_lobster(data: Dictionary) -> void:
	if is_instance_valid(_thunderbolt_lobster_panel):
		_thunderbolt_lobster_panel.handle_thunderbolt_lobster(data)
