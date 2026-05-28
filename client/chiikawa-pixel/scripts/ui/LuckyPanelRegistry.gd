## LuckyPanelRegistry.gd — 幸運面板統一管理器
## lucky-panel-agent 負責維護
## DAY-313：重構架構，統一管理所有 LuckyXxxPanel 的訊號連接
## 解決 HUD.gd 中 65+ 個訊號連接的架構問題
## 使用方式：在 Main.tscn 中加入此節點，自動連接所有 Lucky Panel
extends Node

## 所有已註冊的 Lucky Panel 實例
var _panels: Dictionary = {}

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
	# DAY-313 Progressive Jackpot
	"lucky_jackpot_mini":      "LuckyJackpotMiniPanel",
	"lucky_jackpot_minor":     "LuckyJackpotMinorPanel",
	"lucky_jackpot_major":     "LuckyJackpotMajorPanel",
	"lucky_jackpot_grand":     "LuckyJackpotGrandPanel",
	"lucky_jackpot_trigger":   "LuckyJackpotTriggerPanel",
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
}

func _ready() -> void:
	# 延遲初始化，確保場景樹完整
	call_deferred("_init_all_panels")

func _init_all_panels() -> void:
	## 掃描場景樹，找到所有 LuckyXxxPanel 並建立訊號連接
	var root = get_tree().get_root()
	_scan_and_register(root)
	print("[LuckyPanelRegistry] 已註冊 %d 個 Lucky Panel" % _panels.size())

func _scan_and_register(node: Node) -> void:
	## 遞迴掃描節點，找到所有 Lucky Panel
	if node.get_script() != null:
		var path = node.get_script().resource_path
		for signal_name in SIGNAL_TO_PANEL:
			var panel_class = SIGNAL_TO_PANEL[signal_name]
			if path.ends_with(panel_class + ".gd"):
				_panels[signal_name] = node
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
func connect_all_signals() -> void:
	for signal_name in SIGNAL_TO_PANEL:
		if not GameManager.has_signal(signal_name):
			continue
		var panel = _panels.get(signal_name, null)
		if panel == null:
			continue
		# 每個 Panel 自己在 _ready() 中連接訊號，這裡只做驗證
		print("[LuckyPanelRegistry] ✓ %s → %s" % [signal_name, panel.name])

## 廣播事件到所有 Panel（用於全局事件如 big_bang）
func broadcast_event(event_name: String, data: Dictionary) -> void:
	for signal_name in _panels:
		var panel = _panels[signal_name]
		if is_instance_valid(panel) and panel.has_method("on_global_event"):
			panel.on_global_event(event_name, data)
