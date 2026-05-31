## TargetManager.gd — 目標物系統
## target-system-agent 負責維護
## DAY-322：加入 Shader 支援（sprite_outline / tier_glow / hit_flash）
##           視覺分層系統：依倍率等級自動套用不同視覺強度
##           目標：視覺清晰度 7.5/10 → 8.5/10
extends Node2D

# ── Shader 預載 ───────────────────────────────────────────────
const SHADER_OUTLINE = preload("res://assets/shaders/sprite_outline.gdshader")
const SHADER_TIER_GLOW = preload("res://assets/shaders/tier_glow.gdshader")
const SHADER_HIT_FLASH = preload("res://assets/shaders/hit_flash.gdshader")

# 倍率等級定義（視覺分層）
# Tier 0: <10x   → 無特效（基礎目標）
# Tier 1: 10-29x → 白色輪廓
# Tier 2: 30-99x → 金色輪廓 + 金色光暈
# Tier 3: 100-499x → 橙紅輪廓 + 橙色光暈 + 脈動
# Tier 4: 500-999x → 紫色輪廓 + 紫色光暈 + 快速脈動
# Tier 5: 1000x+  → 宇宙粉紅輪廓 + 彩虹光暈 + 最快脈動
const TIER_THRESHOLDS = [1000.0, 500.0, 100.0, 30.0, 10.0]
const TIER_OUTLINE_COLORS = [
	Color(1.0, 0.0, 0.5, 1.0),   # Tier 5: 宇宙粉紅
	Color(0.8, 0.0, 1.0, 1.0),   # Tier 4: 深紫
	Color(1.0, 0.4, 0.1, 1.0),   # Tier 3: 火橙
	Color(1.0, 0.85, 0.0, 1.0),  # Tier 2: 金色
	Color(0.9, 0.9, 0.9, 0.8),   # Tier 1: 白色
]
const TIER_GLOW_COLORS = [
	Color(1.0, 0.0, 0.5, 0.9),   # Tier 5
	Color(0.8, 0.0, 1.0, 0.8),   # Tier 4
	Color(1.0, 0.4, 0.1, 0.7),   # Tier 3
	Color(1.0, 0.85, 0.0, 0.6),  # Tier 2
	Color(0.9, 0.9, 0.9, 0.4),   # Tier 1
]
const TIER_PULSE_SPEEDS = [4.0, 3.0, 2.5, 1.5, 1.0]
const TIER_GLOW_RADII = [6.0, 5.0, 4.0, 3.0, 2.0]

# 目標物 Sprite 路徑
const TARGET_SPRITES = {
	"T001": "res://assets/sprites/targets/T001_grass.png",
	"T002": "res://assets/sprites/targets/T002_bug_g.png",
	"T003": "res://assets/sprites/targets/T003_bug_r.png",
	"T004": "res://assets/sprites/targets/T004_bug_b.png",
	"T005": "res://assets/sprites/targets/T005_pudding.png",
	"T006": "res://assets/sprites/targets/T006_mushroom.png",
	"T101": "res://assets/sprites/targets/T101_mimic.png",
	"T102": "res://assets/sprites/targets/T102_chest.png",
	"T103": "res://assets/sprites/targets/T103_meteor.png",
	"T104": "res://assets/sprites/targets/T104_gold_grass.png",
	"T105": "res://assets/sprites/targets/T105_coin_fish.png",
	"B001": "res://assets/sprites/targets/B001_boss.png",
	# DAY-292 新增特殊目標
	"T106": "res://assets/sprites/targets/T106_chain_lightning.png",
	"T107": "res://assets/sprites/targets/T107_crab_torpedo.png",
	"T108": "res://assets/sprites/targets/T108_vortex_anemone.png",
	"T109": "res://assets/sprites/targets/T109_golden_dragon.png",
	"T110": "res://assets/sprites/targets/T110_thunder_lobster.png",
	# DAY-293 新增特殊目標
	"T111": "res://assets/sprites/targets/T111_awakened_phoenix.png",
	"T112": "res://assets/sprites/targets/T112_shockwave_bomb.png",
	# DAY-294 新增特殊目標
	"T113": "res://assets/sprites/targets/T113_drill_torpedo.png",
	"T114": "res://assets/sprites/targets/T114_time_freeze.png",
	"T115": "res://assets/sprites/targets/T115_chain_explosion.png",
	# DAY-295 新增特殊目標
	"T116": "res://assets/sprites/targets/T116_chain_long_king.png",
	"T117": "res://assets/sprites/targets/T117_dragon_shotgun.png",
	"T118": "res://assets/sprites/targets/T118_rocket_cannon.png",
	"T119": "res://assets/sprites/targets/T119_deep_whirlpool.png",
	"T120": "res://assets/sprites/targets/T120_vampire_mult.png",
	# DAY-296 新增特殊目標
	"T121": "res://assets/sprites/targets/T121_mirror_fish.png",
	"T122": "res://assets/sprites/targets/T122_golden_rain.png",
	"T123": "res://assets/sprites/targets/T123_freeze_bomb.png",
	"T124": "res://assets/sprites/targets/T124_thunder_storm.png",
	"T125": "res://assets/sprites/targets/T125_lucky_wheel.png",
	# DAY-301 新增特殊目標
	"T126": "res://assets/sprites/targets/T126_jackpot_fish.png",
	"T127": "res://assets/sprites/targets/T127_coop_fish.png",
	"T128": "res://assets/sprites/targets/T128_time_warp.png",
	# DAY-302 新增特殊目標
	"T129": "res://assets/sprites/targets/T129_chain_meteor.png",
	# DAY-303 新增特殊目標
	"T130": "res://assets/sprites/targets/T130_crash_fish.png",
	# DAY-304 新增特殊目標
	"T131": "res://assets/sprites/targets/T131_electric_eel.png",
	"T132": "res://assets/sprites/targets/T132_angler_fish.png",
	"T133": "res://assets/sprites/targets/T133_black_hole.png",
	"T134": "res://assets/sprites/targets/T134_bounty_hunter.png",
	"T135": "res://assets/sprites/targets/T135_tsunami.png",
	# DAY-305 新增特殊目標
	"T136": "res://assets/sprites/targets/T136_dragon_wrath_v2.png",
	"T137": "res://assets/sprites/targets/T137_humpback_whale.png",
	"T138": "res://assets/sprites/targets/T138_legend_dragon.png",
	"T139": "res://assets/sprites/targets/T139_guild_war.png",
	"T140": "res://assets/sprites/targets/T140_quality_fish.png",
	# DAY-306 新增特殊目標
	"T141": "res://assets/sprites/targets/T141_tornado.png",
	"T142": "res://assets/sprites/targets/T142_earthquake.png",
	"T143": "res://assets/sprites/targets/T143_volcano.png",
	"T144": "res://assets/sprites/targets/T144_cosmic_ray.png",
	"T145": "res://assets/sprites/targets/T145_divine_dragon.png",
	# DAY-307 新增特殊目標
	"T146": "res://assets/sprites/targets/T146_quantum.png",
	"T147": "res://assets/sprites/targets/T147_supernova.png",
	"T148": "res://assets/sprites/targets/T148_infinite.png",
	"T149": "res://assets/sprites/targets/T149_genesis.png",
	"T150": "res://assets/sprites/targets/T150_rebirth.png",
	# DAY-308 新增特殊目標
	"T151": "res://assets/sprites/targets/T151_awakened_croc.png",
	"T152": "res://assets/sprites/targets/T152_vampire_v2.png",
	"T153": "res://assets/sprites/targets/T153_super_awaken.png",
	"T154": "res://assets/sprites/targets/T154_giant_prize.png",
	"T155": "res://assets/sprites/targets/T155_immortal_boss.png",
	# DAY-309 新增
	"T156": "res://assets/sprites/targets/T156_ice_phoenix.png",
	"T157": "res://assets/sprites/targets/T157_dragon_fury.png",
	"T158": "res://assets/sprites/targets/T158_mult_cascade.png",
	"T159": "res://assets/sprites/targets/T159_awaken_boss_v2.png",
	"T160": "res://assets/sprites/targets/T160_ultimate_judgment.png",
	# DAY-310 新增
	"T161": "res://assets/sprites/targets/T161_combo_burst.png",
	"T162": "res://assets/sprites/targets/T162_time_bomb.png",
	"T163": "res://assets/sprites/targets/T163_elemental_fusion.png",
	"T164": "res://assets/sprites/targets/T164_treasure_hunter.png",
	"T165": "res://assets/sprites/targets/T165_myth_awaken.png",
	# DAY-312 新增
	"T166": "res://assets/sprites/targets/T166_star_portal.png",
	"T167": "res://assets/sprites/targets/T167_dragon_soul.png",
	"T168": "res://assets/sprites/targets/T168_spacetime_rift.png",
	"T169": "res://assets/sprites/targets/T169_holy_judgment.png",
	"T170": "res://assets/sprites/targets/T170_big_bang.png",
	# DAY-313 新增（Progressive Jackpot 系列）
	"T171": "res://assets/sprites/targets/T171_jackpot_mini.png",
	"T172": "res://assets/sprites/targets/T172_jackpot_minor.png",
	"T173": "res://assets/sprites/targets/T173_jackpot_major.png",
	"T174": "res://assets/sprites/targets/T174_jackpot_grand.png",
	"T175": "res://assets/sprites/targets/T175_jackpot_trigger.png",
	# DAY-314 新增
	"T176": "res://assets/sprites/targets/T176_multiverse.png",
	"T177": "res://assets/sprites/targets/T177_time_loop.png",
	"T178": "res://assets/sprites/targets/T178_fate_wheel.png",
	"T179": "res://assets/sprites/targets/T179_divine_realm.png",
	"T180": "res://assets/sprites/targets/T180_final_power.png",
	# DAY-315 新增
	"T181": "res://assets/sprites/targets/T181_mutation.png",
	"T182": "res://assets/sprites/targets/T182_arctic_storm.png",
	"T183": "res://assets/sprites/targets/T183_fisher_wild.png",
	"T184": "res://assets/sprites/targets/T184_risk_level.png",
	"T185": "res://assets/sprites/targets/T185_cosmic_pulse.png",
	# DAY-316 新增
	"T186": "res://assets/sprites/targets/T186_mirror_universe.png",
	"T187": "res://assets/sprites/targets/T187_gravity_field.png",
	"T188": "res://assets/sprites/targets/T188_time_acceleration.png",
	"T189": "res://assets/sprites/targets/T189_nebula_vortex.png",
	"T190": "res://assets/sprites/targets/T190_cosmic_judgment.png",
	# DAY-317 新增
	"T191": "res://assets/sprites/targets/T191_pvp_battle.png",
	"T192": "res://assets/sprites/targets/T192_skill_chain.png",
	"T193": "res://assets/sprites/targets/T193_global_explosion.png",
	"T194": "res://assets/sprites/targets/T194_spacetime_fold.png",
	"T195": "res://assets/sprites/targets/T195_cosmic_end.png",
	"T196": "res://assets/sprites/targets/T196_dragon_king.png",
	"T197": "res://assets/sprites/targets/T197_eternal_cycle.png",
	"T198": "res://assets/sprites/targets/T198_chaos_explosion.png",
	"T199": "res://assets/sprites/targets/T199_divine_revival.png",
	"T200": "res://assets/sprites/targets/T200_genesis_epoch.png",
	# DAY-319 新增
	"T201": "res://assets/sprites/targets/T201_energy_storm.png",
	"T202": "res://assets/sprites/targets/T202_crystal_resonance.png",
	"T203": "res://assets/sprites/targets/T203_fate_judgment.png",
	"T204": "res://assets/sprites/targets/T204_time_reversal.png",
	"T205": "res://assets/sprites/targets/T205_cosmic_singularity.png",
	# DAY-323 新增
	"T206": "res://assets/sprites/targets/T206_fever_boost.png",
	"T207": "res://assets/sprites/targets/T207_guild_battle.png",
	"T208": "res://assets/sprites/targets/T208_path_fish.png",
	"T209": "res://assets/sprites/targets/T209_chain_eel.png",
	"T210": "res://assets/sprites/targets/T210_ultimate_miracle.png",
	# DAY-324 新增
	"T211": "res://assets/sprites/targets/T211_avalanche.png",
	"T212": "res://assets/sprites/targets/T212_crash_multiplier.png",
	"T213": "res://assets/sprites/targets/T213_multiplier_ladder.png",
	"T214": "res://assets/sprites/targets/T214_ice_fishing_wheel.png",
	"T215": "res://assets/sprites/targets/T215_global_avalanche.png",
	# DAY-325 新增
	"T216": "res://assets/sprites/targets/T216_fishing_net.png",
	"T217": "res://assets/sprites/targets/T217_tnt_bonus.png",
	"T218": "res://assets/sprites/targets/T218_disturbance.png",
	"T219": "res://assets/sprites/targets/T219_pearl_multiplier.png",
	"T220": "res://assets/sprites/targets/T220_rapid_riches.png",
	# DAY-326 新增
	"T221": "res://assets/sprites/targets/T221_dice_bonus.png",
	"T222": "res://assets/sprites/targets/T222_dual_bonus.png",
	"T223": "res://assets/sprites/targets/T223_coin_respin.png",
	# DAY-327 新增
	"T224": "res://assets/sprites/targets/T224_golden_pot.png",
	"T225": "res://assets/sprites/targets/T225_cascade_lock.png",
	"T226": "res://assets/sprites/targets/T226_legend_awaken.png",
	"T227": "res://assets/sprites/targets/T227_crash_harvest.png",
	"T228": "res://assets/sprites/targets/T228_cosmic_fusion.png",
	# DAY-328 新增
	"T229": "res://assets/sprites/targets/T229_magnetic_attraction.png",
	"T230": "res://assets/sprites/targets/T230_super_chain.png",
	"T231": "res://assets/sprites/targets/T231_holy_pillar.png",
	"T232": "res://assets/sprites/targets/T232_time_stop.png",
	"T233": "res://assets/sprites/targets/T233_cosmic_restart.png",
	# DAY-329 新增
	"T234": "res://assets/sprites/targets/T234_fever_boost_ultimate.png",
	"T235": "res://assets/sprites/targets/T235_rapid_riches_ultimate.png",
	"T236": "res://assets/sprites/targets/T236_ice_fishing_master.png",
	"T237": "res://assets/sprites/targets/T237_cosmic_miracle.png",
	"T238": "res://assets/sprites/targets/T238_genesis_ultimate.png",
	# DAY-331 新增
	"T239": "res://assets/sprites/targets/T239_shark_spark.png",
	"T240": "res://assets/sprites/targets/T240_winter_ice.png",
	"T241": "res://assets/sprites/targets/T241_atlantis_frenzy.png",
	"T242": "res://assets/sprites/targets/T242_fishing_time_wheel.png",
	"T243": "res://assets/sprites/targets/T243_ultimate_shark.png",
	# DAY-332 新增
	"T244": "res://assets/sprites/targets/T244_wild_collector.png",
	"T245": "res://assets/sprites/targets/T245_lightning_eel_ultra.png",
	"T246": "res://assets/sprites/targets/T246_domino_chain.png",
	"T247": "res://assets/sprites/targets/T247_immortal_boss_ultra.png",
	"T248": "res://assets/sprites/targets/T248_quad_fusion.png",
	"T249": "res://assets/sprites/targets/T249_electrical_frame.png",
	"T250": "res://assets/sprites/targets/T250_magnetic_respin.png",
	"T251": "res://assets/sprites/targets/T251_fisherman_trail.png",
	"T252": "res://assets/sprites/targets/T252_golden_gills.png",
	"T253": "res://assets/sprites/targets/T253_penta_fusion.png",
}

# 目標物顏色（無 Sprite 時的備用顏色）
const TARGET_COLORS = {
	"T001": Color(0.2, 0.8, 0.2),   # 綠色雜草
	"T002": Color(0.3, 0.9, 0.3),   # 綠色小蟲
	"T003": Color(0.9, 0.3, 0.3),   # 紅色小蟲
	"T004": Color(0.3, 0.5, 0.9),   # 藍色小蟲
	"T005": Color(1.0, 0.9, 0.4),   # 黃色布丁
	"T006": Color(0.6, 0.4, 0.2),   # 棕色蘑菇
	"T101": Color(0.6, 0.6, 0.6),   # 灰色擬態
	"T102": Color(0.9, 0.7, 0.2),   # 金色寶箱
	"T103": Color(1.0, 1.0, 0.8),   # 白色流星
	"T104": Color(1.0, 0.85, 0.0),  # 金色雜草
	"T105": Color(1.0, 0.8, 0.2),   # 金色魚
	"B001": Color(0.8, 0.2, 0.2),   # 紅色 BOSS
	# DAY-292 新增特殊目標備用顏色
	"T106": Color(0.0, 0.9, 1.0),   # 青藍連鎖閃電
	"T107": Color(1.0, 0.4, 0.1),   # 橙紅螃蟹魚雷
	"T108": Color(0.5, 0.2, 0.8),   # 紫色渦旋海葵
	"T109": Color(1.0, 0.85, 0.0),  # 金色黃金龍魚
	"T110": Color(1.0, 0.3, 0.0),   # 火紅雷霆龍蝦
	# DAY-293 新增特殊目標備用顏色
	"T111": Color(1.0, 0.42, 0.21), # 火橙覺醒鳳凰
	"T112": Color(1.0, 0.27, 0.0),  # 深橙全場震盪
	# DAY-294 新增特殊目標備用顏色
	"T113": Color(1.0, 0.55, 0.15), # 橙色鑽頭魚雷
	"T114": Color(0.4, 0.85, 1.0),  # 冰藍時間凍結
	"T115": Color(0.9, 0.2, 0.15),  # 深紅連鎖爆炸
	# DAY-295 新增特殊目標備用顏色
	"T116": Color(1.0, 0.85, 0.0),  # 金色千龍王輪盤
	"T117": Color(0.8, 0.2, 0.9),   # 紫色龍力散彈
	"T118": Color(1.0, 0.3, 0.1),   # 火紅火箭砲
	"T119": Color(0.0, 0.6, 0.9),   # 深藍深海漩渦
	"T120": Color(0.5, 0.0, 0.5),   # 深紫吸血鬼
	# DAY-296 新增特殊目標備用顏色
	"T121": Color(0.88, 0.67, 1.0), # 淡紫鏡像魚
	"T122": Color(1.0, 0.85, 0.0),  # 金色黃金雨魚
	"T123": Color(0.0, 0.9, 1.0),   # 冰藍冰凍炸彈魚
	"T124": Color(1.0, 0.9, 0.2),   # 黃色雷暴魚
	"T125": Color(1.0, 0.42, 0.71), # 粉紅大轉盤魚
	# DAY-301 新增特殊目標備用顏色
	"T126": Color(1.0, 0.85, 0.0),  # 金色進階 Jackpot 魚
	"T127": Color(0.0, 0.9, 1.0),   # 青藍全服合作魚
	"T128": Color(0.55, 0.2, 0.86), # 紫色時間扭曲魚
	# DAY-302 新增特殊目標備用顏色
	"T129": Color(0.9, 0.4, 0.1),   # 火橙連鎖隕石魚
	# DAY-303 新增特殊目標備用顏色
	"T130": Color(0.8, 0.1, 0.1),   # 深紅崩潰魚
	# DAY-304 新增特殊目標備用顏色
	"T131": Color(1.0, 0.95, 0.0),  # 電黃電鰻魚
	"T132": Color(0.0, 0.9, 1.0),   # 青藍巨型安康魚
	"T133": Color(0.3, 0.0, 0.5),   # 深紫黑洞魚
	"T134": Color(1.0, 0.5, 0.0),   # 火橙賞金獵人魚
	"T135": Color(0.0, 0.5, 1.0),   # 深藍海嘯魚
	# DAY-305 新增特殊目標備用顏色
	"T136": Color(1.0, 0.27, 0.0),  # 火橙龍怒蓄積魚
	"T137": Color(0.0, 0.4, 0.8),   # 深藍座頭鯨魚
	"T138": Color(1.0, 0.5, 0.0),   # 金橙傳說龍魚
	"T139": Color(1.0, 0.85, 0.0),  # 金色公會戰魚
	"T140": Color(0.6, 0.0, 1.0),   # 紫色品質魚
	# DAY-306 新增特殊目標備用顏色
	"T141": Color(0.0, 0.9, 0.7),   # 青綠龍捲風魚
	"T142": Color(0.8, 0.4, 0.1),   # 棕橙地震魚
	"T143": Color(1.0, 0.27, 0.0),  # 火紅火山魚
	"T144": Color(0.6, 0.2, 1.0),   # 紫色星際魚
	"T145": Color(1.0, 0.85, 0.0),  # 金色神龍魚
	# DAY-307 新增特殊目標備用顏色
	"T146": Color(0.0, 0.9, 1.0),   # 青藍量子魚
	"T147": Color(1.0, 0.4, 0.1),   # 火橙超新星魚
	"T148": Color(0.6, 0.2, 1.0),   # 紫色無限魚
	"T149": Color(1.0, 0.85, 0.0),  # 金色創世魚
	"T150": Color(1.0, 0.27, 0.0),  # 火橙重生魚
	# DAY-308 新增特殊目標備用顏色
	"T151": Color(0.0, 0.8, 0.3),   # 深綠覺醒鱷魚
	"T152": Color(0.6, 0.0, 0.8),   # 深紫吸血鬼升級魚
	"T153": Color(1.0, 0.4, 0.0),   # 火橙超級覺醒魚
	"T154": Color(1.0, 0.85, 0.0),  # 金色巨型獎勵魚
	"T155": Color(0.9, 0.1, 0.1),   # 深紅不死 BOSS 魚
	# DAY-309 新增
	"T156": Color(0.0, 0.8, 1.0),   # 冰藍冰鳳凰魚
	"T157": Color(1.0, 0.3, 0.0),   # 火橙龍怒能量魚
	"T158": Color(0.1, 0.4, 1.0),   # 深藍倍率瀑布魚
	"T159": Color(1.0, 0.6, 0.0),   # 金橙覺醒 BOSS v2 魚
	"T160": Color(0.8, 0.0, 0.0),   # 深紅終極審判魚
	# DAY-310 新增
	"T161": Color(1.0, 0.4, 0.1),   # 火橙連擊爆發魚
	"T162": Color(0.9, 0.2, 0.1),   # 深紅時間炸彈魚
	"T163": Color(0.6, 0.2, 0.9),   # 深紫元素融合魚
	"T164": Color(1.0, 0.75, 0.1),  # 金色寶藏獵人魚
	"T165": Color(0.9, 0.8, 0.1),   # 神聖金色神話覺醒魚
	# DAY-312 新增
	"T166": Color(0.49, 0.11, 0.64), # 深紫星際門戶魚
	"T167": Color(0.83, 0.18, 0.18), # 龍紅龍魂融合魚
	"T168": Color(0.08, 0.27, 0.75), # 深藍時空裂縫魚
	"T169": Color(0.96, 0.50, 0.09), # 神聖橙金神聖審判魚
	"T170": Color(0.72, 0.07, 0.07), # 深紅宇宙大爆炸魚
	# DAY-313 新增（Progressive Jackpot 系列）
	"T171": Color(0.2, 0.8, 0.2),    # 綠色 Mini Jackpot 魚
	"T172": Color(0.13, 0.59, 0.95), # 藍色 Minor Jackpot 魚
	"T173": Color(1.0, 0.6, 0.0),    # 橙色 Major Jackpot 魚
	"T174": Color(1.0, 0.85, 0.0),   # 金色 Grand Jackpot 魚
	"T175": Color(0.9, 0.7, 0.1),    # 金黃色 Jackpot Trigger 魚
	# DAY-314 新增
	"T176": Color(0.49, 0.11, 0.64), # 深紫多重宇宙魚
	"T177": Color(0.08, 0.39, 0.75), # 深藍時間迴圈魚
	"T178": Color(0.97, 0.50, 0.09), # 火橙命運之輪魚
	"T179": Color(0.98, 0.66, 0.15), # 神聖橙金神域降臨魚
	"T180": Color(0.72, 0.07, 0.07), # 深紅終焉之力魚
	# DAY-315 新增
	"T181": Color(0.6, 0.0, 0.9),    # 深紫突變魚（150種突變）
	"T182": Color(0.0, 0.8, 1.0),    # 冰藍北極風暴魚（8波快速連擊）
	"T183": Color(0.2, 0.9, 0.3),    # 翠綠漁夫野生魚（3個Wild目標）
	"T184": Color(1.0, 0.2, 0.0),    # 深紅風險等級魚（最高×3000）
	"T185": Color(0.7, 0.0, 1.0),    # 深紫宇宙脈衝魚（全場HP-45%）
	# DAY-316 新增
	"T186": Color(0.0, 0.3, 0.9),    # 深藍鏡像宇宙魚（複製最強3個目標）
	"T187": Color(0.5, 0.0, 0.8),    # 深紫引力場魚（引力吸引+爆炸）
	"T188": Color(1.0, 0.4, 0.0),    # 火橙時間加速魚（射擊速度×3.0）
	"T189": Color(0.4, 0.0, 0.7),    # 深紫星雲漩渦魚（每秒HP-8%）
	"T190": Color(0.8, 0.0, 0.0),    # 深紅宇宙審判魚（全場清空×14.0）
	# DAY-317 新增
	"T191": Color(0.8, 0.1, 0.1),    # 紅藍對抗色（PvP 競技魚，劍盾圖案）
	"T192": Color(0.3, 0.5, 1.0),    # 彩虹漸層藍（技能連鎖魚，連鎖符文）
	"T193": Color(0.9, 0.05, 0.05),  # 深紅爆炸（全服大爆炸魚，多層光環）
	"T194": Color(0.45, 0.1, 0.65),  # 藍紫折疊（時空折疊魚，時空裂縫）
	"T195": Color(0.05, 0.0, 0.05),  # 黑色+金色（宇宙終焉魚，終焉符文）
	# DAY-318 新增特殊目標備用顏色
	"T196": Color(0.8, 0.35, 0.0),   # 深橙金（龍王輪盤魚，雙環輪盤）
	"T197": Color(0.0, 0.2, 0.6),    # 深藍（永恆循環魚，無限符號）
	"T198": Color(0.7, 0.1, 0.0),    # 深紅（混沌爆炸魚，混沌光環）
	"T199": Color(0.6, 0.6, 0.0),    # 深金（神聖復活魚，神聖光芒）
	"T200": Color(0.0, 0.0, 0.0),    # 純黑（創世紀元魚，里程碑最高階）
	# DAY-319 新增
	"T201": Color(0.0, 0.8, 0.8),    # 青藍（能量風暴魚，連鎖電擊）
	"T202": Color(0.7, 0.7, 1.0),    # 水晶白藍（水晶共鳴魚，全場共鳴）
	"T203": Color(0.8, 0.7, 0.0),    # 金色（命運審判魚，命運之輪）
	"T204": Color(0.4, 0.4, 0.9),    # 深藍紫（時間逆流魚，時間逆流）
	"T205": Color(0.8, 0.0, 0.8),    # 洋紅（宇宙奇點魚，史上最高）
	# DAY-323 新增
	"T206": Color(1.0, 0.4, 0.0),    # 火橙（Fever Boost 魚，Fever Boost™）
	"T207": Color(1.0, 0.85, 0.0),   # 金色（公會戰魚，公會戰旗幟）
	"T208": Color(0.0, 1.0, 1.0),    # 青藍（路徑魚，路徑光軌）
	"T209": Color(0.8, 0.0, 1.0),    # 深紫（連鎖電鰻魚，紫粉電鰻）
	"T210": Color(1.0, 1.0, 1.0),    # 純白（終極奇蹟魚，新史上最高）
	# DAY-324 新增
	"T211": Color(0.0, 0.75, 1.0),   # 冰藍（雪崩魚，Avalanche Cascade）
	"T212": Color(1.0, 0.27, 0.0),   # 火橙紅（崩潰倍率魚，Crash Multiplier）
	"T213": Color(1.0, 0.84, 0.0),   # 金色（倍率梯魚，Multiplier Ladder）
	"T214": Color(0.0, 0.81, 0.82),  # 冰青（冰釣輪盤魚，Ice Fishing Wheel）
	"T215": Color(0.53, 0.81, 0.98), # 天藍（全服雪崩魚，Global Avalanche 新史上最高）
	# DAY-325 新增
	"T216": Color(0.118, 0.565, 1.0), # 深海藍（漁網魚，Fishing Net ×60.0）
	"T217": Color(1.0, 0.271, 0.0),   # 火橙紅（TNT 爆炸魚，TNT Bonus ×100.0）
	"T218": Color(0.0, 0.808, 0.820), # 深青（擾動魚，Disturbance System 最高 ×50.0）
	"T219": Color(1.0, 0.843, 0.0),   # 金色（珍珠倍率魚，Pearl Multiplier 全服 ×40.0 里程碑）
	"T220": Color(1.0, 1.0, 0.0),     # 亮黃（快速暴富魚，Rapid Riches 全服 ×41.0 新史上最高）
	# DAY-326 新增
	"T221": Color(1.0, 0.6, 0.0),     # 橙金（骰子獎勵魚，Dice Bonus 全服 ×41.5）
	"T222": Color(0.8, 0.0, 0.8),     # 紫色（雙Bonus魚，Dual Bonus 全服 ×42.0）
	"T223": Color(1.0, 0.85, 0.0),    # 深金（Coin Respin 魚，全服 ×42.5 新史上最高）
	# DAY-327 新增
	"T224": Color(0.9, 0.65, 0.0),    # 深金（黃金鍋魚，Gold Blitz 全服 ×43.0 新史上最高）
	"T225": Color(0.0, 0.6, 1.0),     # 深藍（瀑布鎖定魚，Cascade Lock 全服 ×43.5）
	"T226": Color(1.0, 0.3, 0.0),     # 火橙（傳說覺醒魚，Legend Awaken 全服 ×44.0）
	"T227": Color(0.8, 0.1, 0.0),     # 深紅（崩潰收割魚，Crash Harvest 全服 ×44.5）
	"T228": Color(0.8, 0.0, 1.0),     # 宇宙紫（宇宙大融合魚，Cosmic Fusion 全服 ×45.0 新史上最高）
	# DAY-328 新增
	"T229": Color(1.0, 0.4, 0.0),     # 橙色（磁力吸引魚，Magnetic Attraction 全服 ×45.5）
	"T230": Color(0.0, 1.0, 1.0),     # 青色（超級連鎖魚，Super Chain 全服 ×46.0）
	"T231": Color(1.0, 1.0, 0.0),     # 黃色（神聖光柱魚，Holy Pillar 全服 ×46.5）
	"T232": Color(0.0, 0.8, 1.0),     # 冰藍（時間停止魚，Time Stop 全服 ×47.0）
	"T233": Color(1.0, 0.0, 1.0),     # 洋紅（宇宙重啟魚，Cosmic Restart 全服 ×47.5 新史上最高）
	# DAY-329 新增
	"T234": Color(1.0, 0.4, 0.0),     # 橙紅（Fever Boost升級魚，全服 ×48.0）
	"T235": Color(1.0, 0.85, 0.0),    # 金黃（快速暴富升級魚，全服 ×48.5）
	"T236": Color(0.0, 0.75, 1.0),    # 冰藍（冰釣大師魚，全服 ×49.0）
	"T237": Color(0.58, 0.0, 0.83),   # 深紫（宇宙奇蹟魚，全服 ×49.5）
	"T238": Color(1.0, 0.84, 0.0),    # 純金（創世終極魚，全服 ×50.0 里程碑）
	# DAY-331 新增
	"T239": Color(0.0, 0.75, 1.0),    # 深海藍（鯊魚閃電魚，Shark & Spark 全服 ×51.0 里程碑）
	"T240": Color(0.53, 0.81, 0.98),  # 冰藍（冬季冰釣魚，Winter Ice Fishing 全服 ×51.5）
	"T241": Color(0.12, 0.56, 1.0),   # 亞特蘭提斯藍（大西洋狂潮魚，全服 ×52.0）
	"T242": Color(1.0, 0.75, 0.0),    # 金橙（釣魚時間魚，Fishing Time Wheel 全服 ×52.5）
	"T243": Color(1.0, 0.27, 0.0),    # 鯊魚橙紅（終極鯊魚魚，全服 ×53.0 新史上最高）
	# DAY-332 新增
	"T244": Color(1.0, 0.85, 0.0),    # 黃金色（野生收集魚，Wild Collector 全服 ×54.0）
	"T245": Color(0.0, 1.0, 1.0),     # 青色（閃電鰻升級魚，Lightning Eel Ultra 全服 ×54.5）
	"T246": Color(1.0, 0.55, 0.0),    # 骨牌橙（骨牌連鎖魚，Domino Chain 全服 ×55.0 里程碑）
	"T247": Color(0.55, 0.0, 0.0),    # 深紅（不死BOSS升級魚，Immortal Boss Ultra 全服 ×55.5）
	"T248": Color(1.0, 0.0, 1.0),     # 洋紅（四重終極融合魚，Quad Fusion 全服 ×56.0 里程碑）
	# DAY-333 新增
	"T249": Color(0.0, 1.0, 1.0),     # 青色（電擊框架魚，Catfish Hunters 全服 ×56.5）
	"T250": Color(1.0, 0.84, 0.0),    # 黃金色（磁力連鎖魚，Golden Gills 全服 ×57.0）
	"T251": Color(1.0, 0.55, 0.0),    # 橙色（漁夫路徑魚，Bigger Bites 全服 ×57.5）
	"T252": Color(1.0, 0.84, 0.0),    # 黃金色（黃金鰓魚，Golden Gills Jackpot 全服 ×58.0）
	"T253": Color(1.0, 0.41, 0.71),   # 熱粉紅（五重終極魚，Penta Fusion 全服 ×58.5 里程碑）
}

var _target_nodes: Dictionary = {}  # instance_id -> Node2D
var _cached_textures: Dictionary = {}
var _time_warp_speed_mult: float = 1.0  # 時間扭曲速度倍率（DAY-301）

func _ready() -> void:
	GameManager.target_spawned.connect(_on_target_spawned)
	GameManager.target_updated.connect(_on_target_updated)
	GameManager.target_killed.connect(_on_target_killed)
	GameManager.boss_event.connect(_on_boss_event)
	# DAY-301 時間扭曲訊號
	GameManager.lucky_time_warp.connect(_on_lucky_time_warp_for_speed)

func _process(delta: float) -> void:
	_update_positions(delta)

func _update_positions(delta: float) -> void:
	for instance_id in _target_nodes:
		var node = _target_nodes[instance_id]
		if not is_instance_valid(node):
			continue
		var behavior = node.get_meta("behavior", "linear")
		var speed = node.get_meta("speed", 0.0) * _time_warp_speed_mult
		if node.get_meta("is_fleeing", false):
			speed = node.get_meta("flee_speed", speed * 2.5)

		match behavior:
			"linear", "flee", "fast":
				node.position.x -= speed * delta
			"sink":
				node.position.y += speed * 0.3 * delta
				node.position.x -= 10 * delta
			# DAY-339 新增移動模式
			"wave":
				# 波浪移動：X 方向前進 + Y 方向正弦波
				node.position.x -= speed * delta
				var wave_amp = node.get_meta("wave_amp", 40.0)
				var wave_freq = node.get_meta("wave_freq", 2.0)
				var wave_phase = node.get_meta("wave_phase", 0.0)
				var base_y = node.get_meta("base_y", node.position.y)
				var elapsed = node.get_meta("wave_elapsed", 0.0) + delta
				node.set_meta("wave_elapsed", elapsed)
				node.position.y = base_y + sin(elapsed * wave_freq + wave_phase) * wave_amp
			"zigzag":
				# Z字形移動：X 方向前進 + Y 方向鋸齒波
				node.position.x -= speed * delta
				var zz_amp = node.get_meta("zz_amp", 60.0)
				var zz_period = node.get_meta("zz_period", 1.5)
				var base_y = node.get_meta("base_y", node.position.y)
				var elapsed = node.get_meta("zz_elapsed", 0.0) + delta
				node.set_meta("zz_elapsed", elapsed)
				var t_mod = fmod(elapsed, zz_period) / zz_period
				var zz_offset = (2.0 * abs(2.0 * t_mod - 1.0) - 1.0) * zz_amp
				node.position.y = clamp(base_y + zz_offset, 80.0, 640.0)
			"spiral":
				# 螺旋移動：X 方向前進 + Y 方向快速正弦波（模擬螺旋感）
				node.position.x -= speed * 0.7 * delta
				var sp_amp = node.get_meta("sp_amp", 80.0)
				var sp_freq = node.get_meta("sp_freq", 4.0)
				var base_y = node.get_meta("base_y", node.position.y)
				var elapsed = node.get_meta("sp_elapsed", 0.0) + delta
				node.set_meta("sp_elapsed", elapsed)
				node.position.y = base_y + sin(elapsed * sp_freq) * sp_amp * (1.0 - elapsed * 0.05)

		# 離開畫面則移除
		if node.position.x < -100:
			_target_nodes.erase(instance_id)
			node.queue_free()

# ── 目標物生成 ────────────────────────────────────────────────

func _on_target_spawned(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	if _target_nodes.has(instance_id):
		return

	var node = _create_target_node(data)
	_target_nodes[instance_id] = node

	# 進場動畫
	node.scale = Vector2.ZERO
	var tween = node.create_tween()
	var target_type = data.get("type", "basic")
	if target_type == "boss":
		tween.tween_property(node, "scale", Vector2(1.0, 1.0), 0.4).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
	else:
		tween.tween_property(node, "scale", Vector2(1.0, 1.0), 0.15).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)

func _create_target_node(data: Dictionary) -> Node2D:
	var container = Node2D.new()
	container.position = Vector2(data.get("x", 0), data.get("y", 0))
	container.name = "Target_" + data.get("instance_id", "")
	add_child(container)

	var def_id = data.get("def_id", "T001")
	var target_type = data.get("type", "basic")
	var multiplier = data.get("multiplier", 2.0)

	# Sprite
	var sprite = Sprite2D.new()
	var tex = _get_texture(def_id)
	if tex != null:
		sprite.texture = tex
		sprite.texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST
		# 目標物放大：基礎目標 2.5x，特殊目標 2.8x，確保在畫面上清楚可見
		if target_type == "boss":
			sprite.scale = Vector2(3.0, 3.0)
		elif target_type == "special":
			sprite.scale = Vector2(2.8, 2.8)
		else:
			sprite.scale = Vector2(2.5, 2.5)

		# DAY-322：套用 Hit Flash Shader（受擊閃白更精確）
		var hit_mat = ShaderMaterial.new()
		hit_mat.shader = SHADER_HIT_FLASH
		hit_mat.set_shader_parameter("flash_intensity", 0.0)
		sprite.material = hit_mat
		sprite.name = "Sprite"
	else:
		# 備用：ColorRect
		var rect = ColorRect.new()
		var size = 64.0 if target_type != "boss" else 128.0
		rect.size = Vector2(size, size)
		rect.position = -rect.size / 2
		rect.color = TARGET_COLORS.get(def_id, Color(0.8, 0.2, 0.8))
		container.add_child(rect)
	container.add_child(sprite)

	# DAY-322：依倍率等級套用 Outline + Tier Glow Shader
	var tier = _get_multiplier_tier(multiplier)
	if tier >= 0 and target_type != "boss":
		_apply_tier_shaders(container, sprite, multiplier, tier)

	# HP 條
	var hp_max = data.get("max_hp", 1)
	var hp_bg = ColorRect.new()
	hp_bg.size = Vector2(64, 6)
	hp_bg.position = Vector2(-32, -52)
	hp_bg.color = Color(0.2, 0.2, 0.2, 0.8)
	hp_bg.name = "HPBarBG"
	container.add_child(hp_bg)

	var hp_bar = ColorRect.new()
	hp_bar.size = Vector2(64, 6)
	hp_bar.position = Vector2(-32, -52)
	hp_bar.color = Color(0.2, 0.9, 0.2)
	hp_bar.name = "HPBar"
	container.add_child(hp_bar)

	# 倍率標籤
	if target_type != "boss":
		var mult_label = Label.new()
		mult_label.text = "x%.0f" % multiplier
		mult_label.position = Vector2(-20, 36)
		# 視覺清晰度改善（DAY-317）：高倍率用更大字體
		var font_size = 14
		if multiplier >= 500:
			font_size = 20
		elif multiplier >= 100:
			font_size = 17
		elif multiplier >= 30:
			font_size = 15
		mult_label.add_theme_font_size_override("font_size", font_size)
		mult_label.modulate = _mult_label_color(multiplier)
		mult_label.name = "MultLabel"
		container.add_child(mult_label)

	# 高倍率光暈
	if multiplier >= 30.0:
		_add_glow(container, multiplier)

	# DAY-340 游泳動畫：讓目標物有生命感（左右搖擺 + 輕微縮放）
	if target_type != "boss":
		_add_swim_animation(container, def_id, multiplier)

	# Lucky 特殊魚標記（T106-T253）— DAY-338 擴展到 T253
	if def_id.begins_with("T") and def_id.length() >= 4:
		var tid_num = int(def_id.substr(1))
		if tid_num >= 106 and tid_num <= 253:
			_add_lucky_badge(container, def_id)
	if def_id in ["T103", "T104"]:
		var wobble = container.create_tween().set_loops()
		var deg = 5.0 if def_id == "T103" else 3.0
		var dur = 0.15 if def_id == "T103" else 0.4
		wobble.tween_property(container, "rotation_degrees", deg, dur)
		wobble.tween_property(container, "rotation_degrees", -deg, dur)

	# 儲存 meta
	container.set_meta("instance_id", data.get("instance_id", ""))
	container.set_meta("def_id", def_id)
	container.set_meta("speed", data.get("speed", 0.0))
	container.set_meta("behavior", data.get("behavior", "linear"))
	container.set_meta("target_type", target_type)
	container.set_meta("multiplier", multiplier)
	container.set_meta("is_fleeing", false)

	# DAY-339 波浪/Z字形/螺旋移動 meta 初始化
	var behavior = data.get("behavior", "linear")
	var spawn_y = data.get("y", 300.0)
	container.set_meta("base_y", spawn_y)
	match behavior:
		"wave":
			container.set_meta("wave_amp", randf_range(25.0, 55.0))
			container.set_meta("wave_freq", randf_range(1.5, 3.0))
			container.set_meta("wave_phase", randf_range(0.0, TAU))
			container.set_meta("wave_elapsed", 0.0)
		"zigzag":
			container.set_meta("zz_amp", randf_range(40.0, 80.0))
			container.set_meta("zz_period", randf_range(1.0, 2.0))
			container.set_meta("zz_elapsed", 0.0)
		"spiral":
			container.set_meta("sp_amp", randf_range(50.0, 90.0))
			container.set_meta("sp_freq", randf_range(3.0, 5.0))
			container.set_meta("sp_elapsed", 0.0)

	return container

## DAY-340 游泳動畫：讓目標物有生命感
## 基礎目標：左右搖擺（rotation）
## 特殊目標：搖擺 + 輕微縮放脈動
## 高倍率目標：更快更明顯的動畫
func _add_swim_animation(node: Node2D, def_id: String, multiplier: float) -> void:
	var sprite = node.get_node_or_null("Sprite")
	if not is_instance_valid(sprite):
		return

	# 隨機化動畫參數，讓每個目標物看起來不一樣
	var base_duration = randf_range(0.6, 1.2)
	var rot_deg = 5.0  # 搖擺角度

	# 依倍率調整動畫強度
	if multiplier >= 100:
		rot_deg = 8.0
		base_duration = randf_range(0.3, 0.6)  # 高倍率更活躍
	elif multiplier >= 30:
		rot_deg = 6.0
		base_duration = randf_range(0.5, 0.9)
	elif multiplier >= 10:
		rot_deg = 5.0

	# 隨機起始方向
	var start_rot = randf_range(-rot_deg * 0.5, rot_deg * 0.5)
	sprite.rotation_degrees = start_rot

	# 搖擺動畫（循環）
	var swim_tween = sprite.create_tween().set_loops()
	swim_tween.tween_property(sprite, "rotation_degrees", rot_deg, base_duration).set_ease(Tween.EASE_IN_OUT).set_trans(Tween.TRANS_SINE)
	swim_tween.tween_property(sprite, "rotation_degrees", -rot_deg, base_duration).set_ease(Tween.EASE_IN_OUT).set_trans(Tween.TRANS_SINE)

	# 特殊目標（T101+）加入縮放脈動
	if def_id.begins_with("T"):
		var tid_num = int(def_id.substr(1))
		if tid_num >= 101:
			var scale_tween = sprite.create_tween().set_loops()
			var scale_dur = base_duration * 0.8
			scale_tween.tween_property(sprite, "scale", sprite.scale * 1.05, scale_dur).set_ease(Tween.EASE_IN_OUT)
			scale_tween.tween_property(sprite, "scale", sprite.scale * 0.97, scale_dur).set_ease(Tween.EASE_IN_OUT)

func _add_glow(node: Node2D, multiplier: float) -> void:
	var glow = ColorRect.new()
	# 視覺清晰度改善（DAY-317）：依倍率調整光暈大小
	var glow_size = 80.0
	if multiplier >= 1000:
		glow_size = 120.0
	elif multiplier >= 500:
		glow_size = 100.0
	elif multiplier >= 100:
		glow_size = 90.0
	glow.size = Vector2(glow_size, glow_size)
	glow.position = -Vector2(glow_size / 2, glow_size / 2)
	glow.z_index = -1
	if multiplier >= 1000:
		glow.color = Color(1.0, 0.0, 0.5, 0.45)   # 宇宙粉紅（超高倍率）
	elif multiplier >= 500:
		glow.color = Color(0.8, 0.0, 1.0, 0.40)   # 深紫（極高倍率）
	elif multiplier >= 100:
		glow.color = Color(1.0, 0.4, 0.1, 0.35)   # 火橙（高倍率）
	else:
		glow.color = Color(1.0, 0.85, 0.0, 0.25)  # 金色（中高倍率）
	node.add_child(glow)
	var tween = glow.create_tween().set_loops()
	tween.tween_property(glow, "modulate:a", 0.1, 0.4)
	tween.tween_property(glow, "modulate:a", 1.0, 0.4)

# Lucky 特殊魚標記（T106-T190）— 左上角 LUCKY 徽章 + 脈動光環
func _add_lucky_badge(node: Node2D, def_id: String) -> void:
	# 脈動光環（比普通光暈更大更亮）
	var ring = ColorRect.new()
	ring.size = Vector2(96, 96)
	ring.position = Vector2(-48, -48)
	ring.z_index = -1
	# 依倍率範圍選顏色
	var tid_num = int(def_id.substr(1))
	var ring_color: Color
	if tid_num >= 246:
		ring_color = Color(1.0, 0.41, 0.71, 1.0)   # 熱粉紅（T246+，里程碑最高階）
	elif tid_num >= 241:
		ring_color = Color(1.0, 0.0, 1.0, 1.0)     # 洋紅（T241+）
	elif tid_num >= 236:
		ring_color = Color(0.0, 1.0, 1.0, 1.0)     # 青色（T236+）
	elif tid_num >= 231:
		ring_color = Color(1.0, 0.84, 0.0, 1.0)    # 純金（T231+）
	elif tid_num >= 226:
		ring_color = Color(1.0, 0.4, 0.0, 1.0)     # 橙紅（T226+）
	elif tid_num >= 221:
		ring_color = Color(0.8, 0.0, 0.8, 1.0)     # 紫色（T221+）
	elif tid_num >= 216:
		ring_color = Color(1.0, 0.85, 0.0, 1.0)    # 金色（T216+）
	elif tid_num >= 211:
		ring_color = Color(0.0, 0.75, 1.0, 1.0)    # 冰藍（T211+）
	elif tid_num >= 206:
		ring_color = Color(1.0, 0.4, 0.0, 1.0)     # 火橙（T206+）
	elif tid_num >= 201:
		ring_color = Color(0.8, 0.0, 0.8, 1.0)     # 洋紅（T201+）
	elif tid_num >= 196:
		ring_color = Color(1.0, 1.0, 1.0, 1.0)     # 純白（T196+，DAY-318 里程碑最高階）
	elif tid_num >= 191:
		ring_color = Color(1.0, 0.85, 0.0, 1.0)    # 最亮金（T191+，DAY-317 最高階）
	elif tid_num >= 186:
		ring_color = Color(1.0, 0.0, 0.5, 1.0)     # 宇宙粉紅（T186+，DAY-316 最高階）
	elif tid_num >= 181:
		ring_color = Color(1.0, 0.85, 0.0, 0.95)   # 最亮金（T181+，DAY-315 最高階）
	elif tid_num >= 171:
		ring_color = Color(1.0, 0.85, 0.0, 0.85)   # 超亮金（T171+，Progressive Jackpot）
	elif tid_num >= 166:
		ring_color = Color(1.0, 1.0, 0.8, 0.70)    # 極亮白金（T166+，DAY-312 最高階）
	elif tid_num >= 141:
		ring_color = Color(1.0, 1.0, 0.5, 0.60)    # 超亮金（T141+，最高階）
	elif tid_num >= 131:
		ring_color = Color(1.0, 0.95, 0.0, 0.50)   # 亮金（T131-T140）
	elif tid_num >= 126:
		ring_color = Color(1.0, 0.85, 0.0, 0.40)   # 金色（T126-T130）
	elif tid_num >= 121:
		ring_color = Color(0.88, 0.67, 1.0, 0.35)  # 淡紫（T121-T125）
	elif tid_num >= 116:
		ring_color = Color(1.0, 0.85, 0.0, 0.35)   # 金色（T116-T120）
	elif tid_num >= 111:
		ring_color = Color(1.0, 0.42, 0.21, 0.35)  # 火橙（T111-T115）
	else:
		ring_color = Color(0.0, 0.9, 1.0, 0.35)    # 青藍（T106-T110）
	ring.color = ring_color
	node.add_child(ring)

	# 脈動動畫（比普通光暈快）
	var tween = ring.create_tween().set_loops()
	tween.tween_property(ring, "modulate:a", 0.05, 0.25)
	tween.tween_property(ring, "modulate:a", 1.0, 0.25)

	# LUCKY 徽章（左上角小標籤）
	var badge = Label.new()
	# DAY-338 擴展到 T253
	if tid_num >= 246:
		badge.text = "🌸"  # 五重終極（T246+，里程碑最高階）
	elif tid_num >= 241:
		badge.text = "🔮"  # 宇宙（T241+）
	elif tid_num >= 236:
		badge.text = "⚡"  # 閃電（T236+）
	elif tid_num >= 231:
		badge.text = "🌟"  # 星光（T231+）
	elif tid_num >= 226:
		badge.text = "🔥"  # 火焰（T226+）
	elif tid_num >= 221:
		badge.text = "💎"  # 鑽石（T221+）
	elif tid_num >= 216:
		badge.text = "🌊"  # 海浪（T216+）
	elif tid_num >= 211:
		badge.text = "❄️"  # 冰（T211+）
	elif tid_num >= 206:
		badge.text = "🎯"  # 目標（T206+）
	elif tid_num >= 201:
		badge.text = "🌌"  # 宇宙（T201+）
	elif tid_num >= 196:
		badge.text = "🌌"  # 宇宙（T196+，DAY-318 里程碑最高階）
	elif tid_num >= 191:
		badge.text = "💀"  # 終焉（T191+，DAY-317 最高階）
	elif tid_num >= 186:
		badge.text = "🌌"  # 宇宙（T186+，DAY-316 最高階）
	elif tid_num >= 181:
		badge.text = "💫"  # 宇宙星（T181+，DAY-315 最高階）
	elif tid_num >= 171:
		badge.text = "🎰"  # 老虎機（T171+，Progressive Jackpot）
	else:
		badge.text = "✨"  # 閃光（T106-T170）
	badge.position = Vector2(-48, -68)
	badge.add_theme_font_size_override("font_size", 14)
	badge.z_index = 10
	node.add_child(badge)

	# 徽章浮動動畫
	var badge_tween = badge.create_tween().set_loops()
	badge_tween.tween_property(badge, "position:y", -72.0, 0.5)
	badge_tween.tween_property(badge, "position:y", -68.0, 0.5)

# ── 目標物更新 ────────────────────────────────────────────────

func _on_target_updated(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	if not _target_nodes.has(instance_id):
		return
	var node = _target_nodes[instance_id]
	if not is_instance_valid(node):
		return

	# 更新 HP 條
	update_target_hp(instance_id, data.get("hp", 0), data.get("max_hp", 1))

	# 受擊閃白
	_flash_hit(node)

	# T102 逃跑
	if data.get("is_fleeing", false):
		node.set_meta("is_fleeing", true)
		node.set_meta("flee_speed", node.get_meta("speed", 70.0) * 2.5)
		var tween = node.create_tween()
		tween.tween_property(node, "modulate", Color(2.0, 0.5, 0.5), 0.06)
		tween.tween_property(node, "modulate", Color.WHITE, 0.06)
		tween.tween_property(node, "modulate", Color(2.0, 0.5, 0.5), 0.06)
		tween.tween_property(node, "modulate", Color.WHITE, 0.06)

func update_target_hp(instance_id: String, hp: int, max_hp: int) -> void:
	if not _target_nodes.has(instance_id):
		return
	var node = _target_nodes[instance_id]
	if not is_instance_valid(node):
		return
	var hp_bar = node.get_node_or_null("HPBar")
	var hp_bg = node.get_node_or_null("HPBarBG")
	if not is_instance_valid(hp_bar) or not is_instance_valid(hp_bg):
		return
	var pct = float(hp) / float(max_hp) if max_hp > 0 else 0.0
	hp_bar.size.x = hp_bg.size.x * pct
	if pct > 0.6:
		hp_bar.color = Color(0.2, 0.9, 0.2)
	elif pct > 0.3:
		hp_bar.color = Color(1.0, 0.8, 0.1)
	else:
		hp_bar.color = Color(1.0, 0.2, 0.2)

func _flash_hit(node: Node2D) -> void:
	# DAY-322：優先使用 Hit Flash Shader（更精確，只閃有像素的地方）
	var sprite = node.get_node_or_null("Sprite")
	if is_instance_valid(sprite) and sprite.material is ShaderMaterial:
		var mat = sprite.material as ShaderMaterial
		var tween = node.create_tween()
		tween.tween_method(func(v: float): mat.set_shader_parameter("flash_intensity", v), 0.0, 1.0, 0.04)
		tween.tween_method(func(v: float): mat.set_shader_parameter("flash_intensity", v), 1.0, 0.0, 0.08)
		return
	# 備用：舊的 modulate 方式
	for child in node.get_children():
		if child is Sprite2D:
			var tween = node.create_tween()
			tween.tween_property(child, "modulate", Color(3.0, 3.0, 3.0), 0.04)
			tween.tween_property(child, "modulate", Color.WHITE, 0.08)
			break

# ── 目標物擊破 ────────────────────────────────────────────────

func _on_target_killed(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	if not _target_nodes.has(instance_id):
		return
	var node = _target_nodes[instance_id]
	_target_nodes.erase(instance_id)

	if not is_instance_valid(node):
		return

	var reward = data.get("reward", 0)
	var multiplier = data.get("multiplier", 1.0)

	# 擊破特效
	HitEffect.spawn_kill(node.position, multiplier)
	if reward > 0:
		HitEffect.spawn_reward_text(node.position, reward, multiplier)

	# 大獎演出
	if multiplier >= 20:
		HitEffect.spawn_big_win(node.position, multiplier)
		AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
	else:
		AudioManager.play_sfx(AudioManager.SFX.KILL)
		# DAY-340：擊破後播放金幣音效（增加爽感）
		if reward > 0:
			AudioManager.play_sfx(AudioManager.SFX.COIN_DROP)

	# T105 金幣魚：金幣雨
	if data.get("def_id", "") == "T105":
		_spawn_coin_rain(node.position)

	# 消失動畫（DAY-340 升級：更有爽感的爆炸消失）
	var tween = node.create_tween()
	if multiplier >= 50:
		# 高倍率：先放大再消失（爆炸感）
		tween.tween_property(node, "scale", Vector2(1.8, 1.8), 0.08).set_ease(Tween.EASE_OUT)
		tween.parallel().tween_property(node, "modulate", Color(3.0, 3.0, 3.0, 1.0), 0.08)
		tween.tween_property(node, "scale", Vector2(0.1, 0.1), 0.12).set_ease(Tween.EASE_IN)
		tween.parallel().tween_property(node, "modulate:a", 0.0, 0.12)
	elif multiplier >= 10:
		# 中倍率：放大後消失
		tween.tween_property(node, "scale", Vector2(1.5, 1.5), 0.1)
		tween.parallel().tween_property(node, "modulate:a", 0.0, 0.15)
	else:
		# 低倍率：直接消失
		tween.tween_property(node, "scale", Vector2(1.3, 1.3), 0.08)
		tween.parallel().tween_property(node, "modulate:a", 0.0, 0.12)
	tween.tween_callback(func(): if is_instance_valid(node): node.queue_free())

func _spawn_coin_rain(pos: Vector2) -> void:
	for i in 10:
		var coin = ColorRect.new()
		coin.size = Vector2(12, 12)
		coin.color = Color(1.0, 0.85, 0.0)
		coin.position = pos
		coin.z_index = 45
		get_parent().add_child(coin)
		var angle = randf_range(-PI/2 - 0.5, -PI/2 + 0.5)
		var dist = randf_range(40, 100)
		var target = pos + Vector2(cos(angle), sin(angle)) * dist
		var tween = coin.create_tween()
		tween.tween_property(coin, "position", target, 0.4)
		tween.tween_property(coin, "position:y", target.y + 60, 0.3)
		tween.parallel().tween_property(coin, "modulate:a", 0.0, 0.3)
		tween.tween_callback(func(): if is_instance_valid(coin): coin.queue_free())

# ── BOSS 事件 ─────────────────────────────────────────────────

func _on_boss_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	if event == "phase_change":
		var phase = event_data.get("phase", 2)
		var instance_id = event_data.get("instance_id", "")
		if _target_nodes.has(instance_id):
			var node = _target_nodes[instance_id]
			if is_instance_valid(node):
				if phase == 3:
					# Phase 3 絕望模式：更快閃爍 + 更大縮放 + 深紅色調
					var tween = node.create_tween()
					# 5次快速閃爍（比 Phase 2 更快，0.04s vs 0.06s）
					for i in 5:
						tween.tween_property(node, "modulate", Color(5.0, 0.1, 0.1), 0.04)
						tween.tween_property(node, "modulate", Color(0.3, 0.05, 0.05), 0.04)
					# 放大到 1.3x（比 Phase 2 的 1.2x 更大）
					tween.tween_property(node, "scale", Vector2(1.3, 1.3), 0.4).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
					# 持續深紅色調（比 Phase 2 更深）
					tween.tween_property(node, "modulate", Color(2.5, 0.2, 0.2), 0.2)
					# 加入 Phase 3 標記
					_add_phase3_indicator(node)
					ScreenShake.add_trauma(1.0)
					AudioManager.play_bgm(AudioManager.BGM.BOSS_RAGE)
				else:
					# Phase 2 進入：強烈閃爍 + 放大 + 持續紅色調
					var tween = node.create_tween()
					# 5次強烈閃爍
					for i in 5:
						tween.tween_property(node, "modulate", Color(4.0, 0.2, 0.2), 0.06)
						tween.tween_property(node, "modulate", Color(0.5, 0.1, 0.1), 0.06)
					# 放大到 1.2x（Phase 2 BOSS 更大更威脅）
					tween.tween_property(node, "scale", Vector2(1.2, 1.2), 0.4).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
					# 持續紅色調（Phase 2 標誌）
					tween.tween_property(node, "modulate", Color(1.8, 0.4, 0.4), 0.2)
					# 加入 Phase 2 標記
					_add_phase2_indicator(node)
					ScreenShake.add_trauma(0.8)
					AudioManager.play_bgm(AudioManager.BGM.BOSS_RAGE)

func _add_phase2_indicator(boss_node: Node2D) -> void:
	# 移除舊的 Phase 2 標記（如果有）
	var old = boss_node.get_node_or_null("Phase2Label")
	if is_instance_valid(old):
		old.queue_free()
	# 新增 PHASE 2 標籤
	var lbl = Label.new()
	lbl.name = "Phase2Label"
	lbl.text = "⚠ PHASE 2"
	lbl.position = Vector2(-40, -80)
	lbl.add_theme_font_size_override("font_size", 14)
	lbl.modulate = Color(1.0, 0.2, 0.2)
	lbl.z_index = 15
	boss_node.add_child(lbl)
	# 脈動動畫
	var tween = lbl.create_tween().set_loops()
	tween.tween_property(lbl, "modulate:a", 0.3, 0.3)
	tween.tween_property(lbl, "modulate:a", 1.0, 0.3)

func _add_phase3_indicator(boss_node: Node2D) -> void:
	# 移除舊的 Phase 2 / Phase 3 標記（如果有）
	for label_name in ["Phase2Label", "Phase3Label"]:
		var old = boss_node.get_node_or_null(label_name)
		if is_instance_valid(old):
			old.queue_free()
	# 新增 PHASE 3 標籤（紅色，更大字體）
	var lbl = Label.new()
	lbl.name = "Phase3Label"
	lbl.text = "💀 PHASE 3 絕望模式！"
	lbl.position = Vector2(-72, -96)
	lbl.add_theme_font_size_override("font_size", 18)
	lbl.modulate = Color(1.0, 0.0, 0.0)
	lbl.z_index = 20
	boss_node.add_child(lbl)
	# 快速脈動動畫（比 Phase 2 更快）
	var tween = lbl.create_tween().set_loops()
	tween.tween_property(lbl, "modulate:a", 0.1, 0.15)
	tween.tween_property(lbl, "modulate:a", 1.0, 0.15)

# ── 點擊目標 ──────────────────────────────────────────────────

func try_click_target(click_pos: Vector2) -> String:
	var best_id = ""
	var best_dist = 80.0  # 點擊半徑
	for instance_id in _target_nodes:
		var node = _target_nodes[instance_id]
		if not is_instance_valid(node):
			continue
		var dist = node.position.distance_to(click_pos)
		if dist < best_dist:
			best_dist = dist
			best_id = instance_id
	return best_id

# ── 輔助 ──────────────────────────────────────────────────────

func _get_texture(def_id: String) -> Texture2D:
	var path = TARGET_SPRITES.get(def_id, "")
	if path == "":
		return null
	if _cached_textures.has(path):
		return _cached_textures[path]
	if ResourceLoader.exists(path):
		var tex = load(path)
		_cached_textures[path] = tex
		return tex
	return null

func _mult_label_color(mult: float) -> Color:
	if mult >= 1000: return Color(1.0, 0.0, 0.5)  # 宇宙粉紅（Tier 5）
	if mult >= 500:  return Color(0.8, 0.0, 1.0)  # 深紫（Tier 4）
	if mult >= 100:  return Color(1.0, 0.3, 0.1)  # 火紅（Tier 3）
	if mult >= 50:   return Color(1.0, 0.4, 0.1)  # 橙紅
	if mult >= 30:   return Color(1.0, 0.85, 0.0) # 金色（Tier 2）
	if mult >= 15:   return Color(0.8, 0.9, 1.0)  # 淡藍
	return Color(1.0, 1.0, 1.0)                   # 白色

# ── DAY-322 視覺分層 Shader 系統 ──────────────────────────────

## 取得倍率等級（-1 = 無特效，0-4 = Tier 1-5）
func _get_multiplier_tier(multiplier: float) -> int:
	for i in TIER_THRESHOLDS.size():
		if multiplier >= TIER_THRESHOLDS[i]:
			return i
	return -1  # <10x，無特效

## 套用 Outline + Tier Glow Shader
func _apply_tier_shaders(container: Node2D, sprite: Sprite2D, multiplier: float, tier: int) -> void:
	if tier < 0 or tier >= TIER_OUTLINE_COLORS.size():
		return

	var outline_color = TIER_OUTLINE_COLORS[tier]
	var glow_color = TIER_GLOW_COLORS[tier]
	var pulse_speed = TIER_PULSE_SPEEDS[tier]
	var glow_radius = TIER_GLOW_RADII[tier]

	# 建立 Outline Sprite（在原始 Sprite 後面一層）
	var outline_sprite = Sprite2D.new()
	outline_sprite.texture = sprite.texture
	outline_sprite.texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST
	outline_sprite.scale = sprite.scale
	outline_sprite.z_index = sprite.z_index - 1
	outline_sprite.name = "OutlineSprite"

	var outline_mat = ShaderMaterial.new()
	outline_mat.shader = SHADER_OUTLINE
	outline_mat.set_shader_parameter("outline_color", outline_color)
	outline_mat.set_shader_parameter("outline_width", 1.5 + float(tier) * 0.3)
	outline_mat.set_shader_parameter("pulse_speed", pulse_speed)
	outline_mat.set_shader_parameter("pulse_min", 0.3)
	outline_mat.set_shader_parameter("enable_pulse", tier >= 2)  # Tier 2+ 才脈動
	outline_sprite.material = outline_mat
	container.add_child(outline_sprite)
	container.move_child(outline_sprite, 0)  # 移到最底層

	# Tier 3+ 加入 Tier Glow（光暈效果）
	if tier >= 2:
		var glow_sprite = Sprite2D.new()
		glow_sprite.texture = sprite.texture
		glow_sprite.texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST
		glow_sprite.scale = sprite.scale * (1.0 + float(tier) * 0.05)
		glow_sprite.z_index = sprite.z_index - 2
		glow_sprite.name = "GlowSprite"

		var glow_mat = ShaderMaterial.new()
		glow_mat.shader = SHADER_TIER_GLOW
		glow_mat.set_shader_parameter("glow_color", glow_color)
		glow_mat.set_shader_parameter("glow_radius", glow_radius)
		glow_mat.set_shader_parameter("glow_intensity", 1.0 + float(tier) * 0.3)
		glow_mat.set_shader_parameter("pulse_speed", pulse_speed * 0.7)
		glow_mat.set_shader_parameter("pulse_amplitude", 0.15 + float(tier) * 0.05)
		glow_sprite.material = glow_mat
		container.add_child(glow_sprite)
		container.move_child(glow_sprite, 0)  # 移到最底層

	# Tier 5（1000x+）：額外加入彩虹旋轉光環
	if tier == 0:
		_add_rainbow_ring(container, multiplier)

## 彩虹旋轉光環（Tier 5 / 1000x+ 專用）
func _add_rainbow_ring(container: Node2D, multiplier: float) -> void:
	var ring_count = 3
	for i in ring_count:
		var ring = ColorRect.new()
		var ring_size = 110.0 + float(i) * 15.0
		ring.size = Vector2(ring_size, ring_size)
		ring.position = -Vector2(ring_size / 2, ring_size / 2)
		ring.z_index = -3 - i
		ring.name = "RainbowRing_%d" % i
		# 彩虹顏色循環
		var hue = float(i) / float(ring_count)
		ring.color = Color.from_hsv(hue, 1.0, 1.0, 0.3)
		container.add_child(ring)
		container.move_child(ring, 0)
		# 旋轉動畫（每個環速度不同）
		var rot_tween = ring.create_tween().set_loops()
		var rot_speed = (0.8 + float(i) * 0.3) * (1.0 if i % 2 == 0 else -1.0)
		rot_tween.tween_property(ring, "rotation", TAU * rot_speed, 2.0)
		# 顏色循環動畫
		var color_tween = ring.create_tween().set_loops()
		color_tween.tween_method(
			func(h: float): ring.color = Color.from_hsv(h, 1.0, 1.0, 0.3),
			float(i) / float(ring_count),
			float(i) / float(ring_count) + 1.0,
			3.0
		)

# ── DAY-301 時間扭曲速度效果 ──────────────────────────────────

func _on_lucky_time_warp_for_speed(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"warp_start":
			_time_warp_speed_mult = data.get("speed_mult", 0.3)
			# 視覺提示：所有目標物變藍色調
			for instance_id in _target_nodes:
				var node = _target_nodes[instance_id]
				if is_instance_valid(node):
					var tween = node.create_tween()
					tween.tween_property(node, "modulate", Color(0.7, 0.8, 1.2), 0.3)
		"warp_end", "time_collapse", "collapse_end":
			_time_warp_speed_mult = 1.0
			# 恢復正常顏色
			for instance_id in _target_nodes:
				var node = _target_nodes[instance_id]
				if is_instance_valid(node):
					var tween = node.create_tween()
					tween.tween_property(node, "modulate", Color.WHITE, 0.3)
