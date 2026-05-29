## LuckyPanelRegistry.gd — 幸運面板統一管理器
## lucky-panel-agent 負責維護
## DAY-313：重構架構，統一管理所有 LuckyXxxPanel 的訊號連接
## 解決 HUD.gd 中 65+ 個訊號連接的架構問題
## 使用方式：在 Main.tscn 中加入此節點，自動連接所有 Lucky Panel
extends Node

## 所有已註冊的 Lucky Panel 實例
var _panels: Dictionary = {}

## Progressive Jackpot Panel 映射（tier → Panel 腳本名稱）
## 這些 Panel 統一由 lucky_jackpot_pool 訊號分發
const JACKPOT_TIER_TO_PANEL: Dictionary = {
	"mini":    "LuckyJackpotMiniPanel",
	"minor":   "LuckyJackpotMinorPanel",
	"major":   "LuckyJackpotMajorPanel",
	"grand":   "LuckyJackpotGrandPanel",
	"trigger": "LuckyJackpotTriggerPanel",
}
## Jackpot Panel 實例（tier → Panel）
var _jackpot_panels: Dictionary = {}

## 訊號→Panel 映射表（訊號名稱 → Panel 腳本名稱）
const SIGNAL_TO_PANEL: Dictionary = {
	# DAY-292
	"lucky_chain_lightning":   "LuckyChainLightningPanel",
	"lucky_crab_torpedo":      "LuckyCrabTorpedoPanel",
	"lucky_vortex":            "LuckyVortexPanel",
	"lucky_golden_dragon":     "LuckyGoldenDragonPanel",
	"lucky_thunder_lobster":   "LuckyThunderLobsterPanel",
	# DAY-293
	"lucky_awakened_phoenix":  "LuckyAwakenedPhoenixPanel",
	"lucky_shockwave_bomb":    "LuckyShockwaveBombPanel",
	# DAY-294
	"lucky_drill_torpedo":     "LuckyDrillTorpedoPanel",
	"lucky_time_freeze":       "LuckyTimeFreezePanel",
	"lucky_chain_explosion":   "LuckyChainExplosionPanel",
	# DAY-295
	"lucky_chain_long_king":   "LuckyChainLongKingPanel",
	"lucky_dragon_shotgun":    "LuckyDragonShotgunPanel",
	"lucky_rocket_cannon":     "LuckyRocketCannonPanel",
	"lucky_deep_whirlpool":    "LuckyDeepWhirlpoolPanel",
	"lucky_vampire_mult":      "LuckyVampireMultPanel",
	# DAY-296
	"lucky_mirror_fish":       "LuckyMirrorFishPanel",
	"lucky_golden_rain":       "LuckyGoldenRainPanel",
	"lucky_freeze_bomb":       "LuckyFreezeBombPanel",
	"lucky_thunder_storm":     "LuckyThunderStormPanel",
	"lucky_lucky_wheel":       "LuckyLuckyWheelPanel",
	# DAY-301
	"lucky_jackpot_fish":      "LuckyJackpotFishPanel",
	"lucky_coop_fish":         "LuckyCoopFishPanel",
	"lucky_time_warp":         "LuckyTimeWarpPanel",
	# DAY-302
	"lucky_chain_meteor":      "LuckyChainMeteorPanel",
	# DAY-303
	"lucky_crash_fish":        "LuckyCrashFishPanel",
	# DAY-304
	"lucky_electric_eel":      "LuckyElectricEelPanel",
	"lucky_angler_fish":       "LuckyAnglerFishPanel",
	"lucky_black_hole":        "LuckyBlackHolePanel",
	"lucky_bounty_hunter":     "LuckyBountyHunterPanel",
	"lucky_tsunami":           "LuckyTsunamiPanel",
	# DAY-305
	"lucky_dragon_wrath_v2":   "LuckyDragonWrathV2Panel",
	"lucky_humpback_whale":    "LuckyHumpbackWhalePanel",
	"lucky_legend_dragon":     "LuckyLegendDragonPanel",
	"lucky_guild_war":         "LuckyGuildWarPanel",
	"lucky_quality_fish":      "LuckyQualityFishPanel",
	# DAY-306
	"lucky_tornado":           "LuckyTornadoPanel",
	"lucky_earthquake":        "LuckyEarthquakePanel",
	"lucky_volcano":           "LuckyVolcanoPanel",
	"lucky_cosmic_ray":        "LuckyCosmicRayPanel",
	"lucky_divine_dragon":     "LuckyDivineDragonPanel",
	# DAY-307
	"lucky_quantum":           "LuckyQuantumPanel",
	"lucky_supernova":         "LuckySupernovaPanel",
	"lucky_infinite":          "LuckyInfinitePanel",
	"lucky_genesis":           "LuckyGenesisPanel",
	"lucky_rebirth":           "LuckyRebirthPanel",
	# DAY-308
	"lucky_awakened_croc":     "LuckyAwakenedCrocPanel",
	"lucky_vampire_v2":        "LuckyVampireV2Panel",
	"lucky_super_awaken":      "LuckySuperAwakenPanel",
	"lucky_giant_prize":       "LuckyGiantPrizePanel",
	"lucky_immortal_boss":     "LuckyImmortalBossPanel",
	# DAY-309
	"lucky_ice_phoenix":       "LuckyIcePhoenixPanel",
	"lucky_dragon_fury":       "LuckyDragonFuryPanel",
	"lucky_mult_cascade":      "LuckyMultCascadePanel",
	"lucky_awaken_boss_v2":    "LuckyAwakenBossV2Panel",
	"lucky_ultimate_judgment": "LuckyUltimateJudgmentPanel",
	# DAY-310
	"lucky_combo_burst":       "LuckyComboBurstPanel",
	"lucky_time_bomb":         "LuckyTimeBombPanel",
	"lucky_elemental_fusion":  "LuckyElementalFusionPanel",
	"lucky_treasure_hunter":   "LuckyTreasureHunterPanel",
	"lucky_myth_awaken":       "LuckyMythAwakenPanel",
	# DAY-312
	"lucky_star_portal":       "LuckyStarPortalPanel",
	"lucky_dragon_soul":       "LuckyDragonSoulPanel",
	"lucky_spacetime_rift":    "LuckySpacetimeRiftPanel",
	"lucky_holy_judgment":     "LuckyHolyJudgmentPanel",
	"lucky_big_bang":          "LuckyBigBangPanel",
	# DAY-313 Progressive Jackpot（統一由 lucky_jackpot_pool 訊號分發）
	# 注意：這 5 個 Panel 都監聽 lucky_jackpot_pool，由 _dispatch_jackpot_pool 分發
	# "lucky_jackpot_mini":   "LuckyJackpotMiniPanel",   # 由 lucky_jackpot_pool 分發
	# "lucky_jackpot_minor":  "LuckyJackpotMinorPanel",  # 由 lucky_jackpot_pool 分發
	# "lucky_jackpot_major":  "LuckyJackpotMajorPanel",  # 由 lucky_jackpot_pool 分發
	# "lucky_jackpot_grand":  "LuckyJackpotGrandPanel",  # 由 lucky_jackpot_pool 分發
	# "lucky_jackpot_trigger":"LuckyJackpotTriggerPanel",# 由 lucky_jackpot_pool 分發
	# DAY-314 新增
	"lucky_multiverse":        "LuckyMultiversePanel",
	"lucky_time_loop":         "LuckyTimeLoopPanel",
	"lucky_fate_wheel":        "LuckyFateWheelPanel",
	"lucky_divine_realm":      "LuckyDivineRealmPanel",
	"lucky_final_power":       "LuckyFinalPowerPanel",
	# DAY-315 新增
	"lucky_mutation":          "LuckyMutationPanel",
	"lucky_arctic_storm":      "LuckyArcticStormPanel",
	"lucky_fisher_wild":       "LuckyFisherWildPanel",
	"lucky_risk_level":        "LuckyRiskLevelPanel",
	"lucky_cosmic_pulse":      "LuckyCosmicPulsePanel",
	# DAY-316 新增
	"lucky_mirror_universe":   "LuckyMirrorUniversePanel",
	"lucky_gravity_field":     "LuckyGravityFieldPanel",
	"lucky_time_acceleration": "LuckyTimeAccelerationPanel",
	"lucky_nebula_vortex":     "LuckyNebulaVortexPanel",
	"lucky_cosmic_judgment":   "LuckyCosmicJudgmentPanel",
	# DAY-317 新增
	"lucky_pvp_battle":        "LuckyPvpBattlePanel",
	"lucky_skill_chain":       "LuckySkillChainPanel",
	"lucky_global_explosion":  "LuckyGlobalExplosionPanel",
	"lucky_spacetime_fold":    "LuckySpacetimeFoldPanel",
	"lucky_cosmic_end":        "LuckyCosmicEndPanel",
	# DAY-318 新增
	"lucky_dragon_king":       "LuckyDragonKingPanel",
	"lucky_eternal_cycle":     "LuckyEternalCyclePanel",
	"lucky_chaos_explosion":   "LuckyChaosExplosionPanel",
	"lucky_divine_revival":    "LuckyDivineRevivalPanel",
	"lucky_genesis_epoch":     "LuckyGenesisEpochPanel",
	# DAY-319 新增
	"lucky_energy_storm":      "LuckyEnergyStormPanel",
	"lucky_crystal_resonance": "LuckyCrystalResonancePanel",
	"lucky_fate_judgment":     "LuckyFateJudgmentPanel",
	"lucky_time_reversal":     "LuckyTimeReversalPanel",
	"lucky_cosmic_singularity":"LuckyCosmicSingularityPanel",
	# DAY-323 新增
	"lucky_fever_boost":       "LuckyFeverBoostPanel",
	"lucky_guild_battle":      "LuckyGuildBattlePanel",
	"lucky_path_fish":         "LuckyPathFishPanel",
	"lucky_chain_eel":         "LuckyChainEelPanel",
	"lucky_ultimate_miracle":  "LuckyUltimateMiraclePanel",
	# DAY-324 新增
	"lucky_avalanche":         "LuckyAvalanchePanel",
	"lucky_crash_multiplier":  "LuckyCrashMultiplierPanel",
	"lucky_multiplier_ladder": "LuckyMultiplierLadderPanel",
	"lucky_ice_fishing_wheel": "LuckyIceFishingWheelPanel",
	"lucky_global_avalanche":  "LuckyGlobalAvalanchePanel",
	# DAY-325 新增
	"lucky_fishing_net":      "LuckyFishingNetPanel",
	"lucky_tnt_bonus":        "LuckyTNTBonusPanel",
	"lucky_disturbance":      "LuckyDisturbancePanel",
	"lucky_pearl_multiplier": "LuckyPearlMultiplierPanel",
	"lucky_rapid_riches":     "LuckyRapidRichesPanel",
}

func _ready() -> void:
	# 延遲初始化，確保場景樹完整
	call_deferred("_init_all_panels")

func _init_all_panels() -> void:
	## 掃描場景樹，找到所有 LuckyXxxPanel 並建立訊號連接
	var root = get_tree().get_root()
	_scan_and_register(root)
	print("[LuckyPanelRegistry] 已註冊 %d 個 Lucky Panel" % _panels.size())
	# DAY-320：掃描完成後自動連接所有訊號
	connect_all_signals()

func _scan_and_register(node: Node) -> void:
	## 遞迴掃描節點，找到所有 Lucky Panel
	if node.get_script() != null:
		var path = node.get_script().resource_path
		for signal_name in SIGNAL_TO_PANEL:
			var panel_class = SIGNAL_TO_PANEL[signal_name]
			if path.ends_with(panel_class + ".gd"):
				_panels[signal_name] = node
				break
		# 掃描 Jackpot Panel
		for tier in JACKPOT_TIER_TO_PANEL:
			var panel_class = JACKPOT_TIER_TO_PANEL[tier]
			if path.ends_with(panel_class + ".gd"):
				_jackpot_panels[tier] = node
				break
	for child in node.get_children():
		_scan_and_register(child)

## 手動註冊一個 Panel（供動態建立的 Panel 使用）
func register_panel(signal_name: String, panel: Node) -> void:
	_panels[signal_name] = panel
	print("[LuckyPanelRegistry] 手動註冊: %s → %s" % [signal_name, panel.name])

## 取得已註冊的 Panel
func get_panel(signal_name: String) -> Node:
	return _panels.get(signal_name, null)

## 取得所有已註冊的 Panel 數量
func get_panel_count() -> int:
	return _panels.size()

## 連接 GameManager 的所有 Lucky 訊號到對應 Panel
## 這個方法取代 HUD.gd 中的 65+ 個 connect() 呼叫
## DAY-320：修復設計缺陷，真正連接訊號到 Panel 的 handle_event()
func connect_all_signals() -> void:
	var connected = 0
	var skipped = 0
	for signal_name in SIGNAL_TO_PANEL:
		if not GameManager.has_signal(signal_name):
			skipped += 1
			continue
		var panel = _panels.get(signal_name, null)
		if not is_instance_valid(panel):
			skipped += 1
			continue
		if not panel.has_method("handle_event"):
			skipped += 1
			continue
		# 連接 GameManager 訊號到 Panel 的 handle_event()
		# 使用 Callable 確保正確綁定
		var callable = Callable(panel, "handle_event")
		if not GameManager.is_connected(signal_name, callable):
			GameManager.connect(signal_name, callable)
			connected += 1
	# 連接 lucky_jackpot_pool 訊號（分發到各 Jackpot Panel）
	if GameManager.has_signal("lucky_jackpot_pool"):
		var callable = Callable(self, "_dispatch_jackpot_pool")
		if not GameManager.is_connected("lucky_jackpot_pool", callable):
			GameManager.connect("lucky_jackpot_pool", callable)
			connected += 1
	print("[LuckyPanelRegistry] 已連接 %d 個訊號，跳過 %d 個" % [connected, skipped])
	print("[LuckyPanelRegistry] Jackpot Panel 數量：%d" % _jackpot_panels.size())

## 分發 lucky_jackpot_pool 訊號到對應 Jackpot Panel
func _dispatch_jackpot_pool(data: Dictionary) -> void:
	var event = data.get("event", "")
	var tier = data.get("tier", "")
	# pool_update 廣播給所有 Jackpot Panel
	if event == "pool_update":
		for t in _jackpot_panels:
			var panel = _jackpot_panels[t]
			if is_instance_valid(panel) and panel.has_method("handle_event"):
				panel.handle_event(data)
		return
	# jackpot_win 只發給對應 tier 的 Panel
	if tier != "" and _jackpot_panels.has(tier):
		var panel = _jackpot_panels[tier]
		if is_instance_valid(panel) and panel.has_method("handle_event"):
			panel.handle_event(data)
	# T175 trigger 魚：發給 trigger Panel
	var target_id = data.get("target_id", "")
	if target_id == "T175" and _jackpot_panels.has("trigger"):
		var panel = _jackpot_panels["trigger"]
		if is_instance_valid(panel) and panel.has_method("handle_event"):
			panel.handle_event(data)

## 廣播事件到所有 Panel（用於全局事件如 big_bang）
func broadcast_event(event_name: String, data: Dictionary) -> void:
	for signal_name in _panels:
		var panel = _panels[signal_name]
		if is_instance_valid(panel) and panel.has_method("on_global_event"):
			panel.on_global_event(event_name, data)
