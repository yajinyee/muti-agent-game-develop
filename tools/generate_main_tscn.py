#!/usr/bin/env python3
"""
generate_main_tscn.py — 生成完整的 Main.tscn
包含所有 Lucky Panel 節點、LuckyEventSystem、LuckyPanelRegistry
DAY-321
"""

# 所有 Lucky Panel 的腳本名稱（按 DAY 順序）
LUCKY_PANELS = [
    # DAY-292
    "LuckyChainLightningPanel",
    "LuckyCrabTorpedoPanel",
    "LuckyVortexPanel",
    "LuckyGoldenDragonPanel",
    "LuckyThunderLobsterPanel",
    # DAY-293
    "LuckyAwakenedPhoenixPanel",
    "LuckyShockwaveBombPanel",
    # DAY-294
    "LuckyDrillTorpedoPanel",
    "LuckyTimeFreezePanel",
    "LuckyChainExplosionPanel",
    # DAY-295
    "LuckyChainLongKingPanel",
    "LuckyDragonShotgunPanel",
    "LuckyRocketCannonPanel",
    "LuckyDeepWhirlpoolPanel",
    "LuckyVampireMultPanel",
    # DAY-296
    "LuckyMirrorFishPanel",
    "LuckyGoldenRainPanel",
    "LuckyFreezeBombPanel",
    "LuckyThunderStormPanel",
    "LuckyLuckyWheelPanel",
    # DAY-301
    "LuckyJackpotFishPanel",
    "LuckyCoopFishPanel",
    "LuckyTimeWarpPanel",
    # DAY-302
    "LuckyChainMeteorPanel",
    # DAY-303
    "LuckyCrashFishPanel",
    # DAY-304
    "LuckyElectricEelPanel",
    "LuckyAnglerFishPanel",
    "LuckyBlackHolePanel",
    "LuckyBountyHunterPanel",
    "LuckyTsunamiPanel",
    # DAY-305
    "LuckyDragonWrathV2Panel",
    "LuckyHumpbackWhalePanel",
    "LuckyLegendDragonPanel",
    "LuckyGuildWarPanel",
    "LuckyQualityFishPanel",
    # DAY-306
    "LuckyTornadoPanel",
    "LuckyEarthquakePanel",
    "LuckyVolcanoPanel",
    "LuckyCosmicRayPanel",
    "LuckyDivineDragonPanel",
    # DAY-307
    "LuckyQuantumPanel",
    "LuckySupernovaPanel",
    "LuckyInfinitePanel",
    "LuckyGenesisPanel",
    "LuckyRebirthPanel",
    # DAY-308
    "LuckyAwakenedCrocPanel",
    "LuckyVampireV2Panel",
    "LuckySuperAwakenPanel",
    "LuckyGiantPrizePanel",
    "LuckyImmortalBossPanel",
    # DAY-309
    "LuckyIcePhoenixPanel",
    "LuckyDragonFuryPanel",
    "LuckyMultCascadePanel",
    "LuckyAwakenBossV2Panel",
    "LuckyUltimateJudgmentPanel",
    # DAY-310
    "LuckyComboBurstPanel",
    "LuckyTimeBombPanel",
    "LuckyElementalFusionPanel",
    "LuckyTreasureHunterPanel",
    "LuckyMythAwakenPanel",
    # DAY-312
    "LuckyStarPortalPanel",
    "LuckyDragonSoulPanel",
    "LuckySpacetimeRiftPanel",
    "LuckyHolyJudgmentPanel",
    "LuckyBigBangPanel",
    # DAY-313 Progressive Jackpot
    "LuckyJackpotMiniPanel",
    "LuckyJackpotMinorPanel",
    "LuckyJackpotMajorPanel",
    "LuckyJackpotGrandPanel",
    "LuckyJackpotTriggerPanel",
    # DAY-314
    "LuckyMultiversePanel",
    "LuckyTimeLoopPanel",
    "LuckyFateWheelPanel",
    "LuckyDivineRealmPanel",
    "LuckyFinalPowerPanel",
    # DAY-315
    "LuckyMutationPanel",
    "LuckyArcticStormPanel",
    "LuckyFisherWildPanel",
    "LuckyRiskLevelPanel",
    "LuckyCosmicPulsePanel",
    # DAY-316
    "LuckyMirrorUniversePanel",
    "LuckyGravityFieldPanel",
    "LuckyTimeAccelerationPanel",
    "LuckyNebulaVortexPanel",
    "LuckyCosmicJudgmentPanel",
    # DAY-317
    "LuckyPvpBattlePanel",
    "LuckySkillChainPanel",
    "LuckyGlobalExplosionPanel",
    "LuckySpacetimeFoldPanel",
    "LuckyCosmicEndPanel",
    # DAY-318
    "LuckyDragonKingPanel",
    "LuckyEternalCyclePanel",
    "LuckyChaosExplosionPanel",
    "LuckyDivineRevivalPanel",
    "LuckyGenesisEpochPanel",
    # DAY-319
    "LuckyEnergyStormPanel",
    "LuckyCrystalResonancePanel",
    "LuckyFateJudgmentPanel",
    "LuckyTimeReversalPanel",
    "LuckyCosmicSingularityPanel",
]

def generate_main_tscn():
    # 計算 load_steps：7 個基礎腳本 + 3 個 UI 腳本 + 100 個 Lucky Panel 腳本
    # 基礎：TargetManager(1), Cannon(2), HUD(3), BackgroundManager(4), BonusGame(5), CharacterAnimator(6)
    # UI：LuckyEventSystem(7), LuckyPanelRegistry(8), BaseLuckyPanel(9)
    # Lucky Panels：10 ~ 109
    base_count = 9
    total_load_steps = base_count + len(LUCKY_PANELS)
    
    lines = []
    lines.append(f'[gd_scene load_steps={total_load_steps} format=3 uid="uid://main"]')
    lines.append('')
    
    # 基礎腳本 ext_resource
    lines.append('[ext_resource type="Script" path="res://scripts/game/TargetManager.gd" id="1"]')
    lines.append('[ext_resource type="Script" path="res://scripts/game/Cannon.gd" id="2"]')
    lines.append('[ext_resource type="Script" path="res://scripts/ui/HUD.gd" id="3"]')
    lines.append('[ext_resource type="Script" path="res://scripts/game/BackgroundManager.gd" id="4"]')
    lines.append('[ext_resource type="Script" path="res://scripts/game/BonusGame.gd" id="5"]')
    lines.append('[ext_resource type="Script" path="res://scripts/game/CharacterAnimator.gd" id="6"]')
    lines.append('[ext_resource type="Script" path="res://scripts/ui/LuckyEventSystem.gd" id="7"]')
    lines.append('[ext_resource type="Script" path="res://scripts/ui/LuckyPanelRegistry.gd" id="8"]')
    lines.append('[ext_resource type="Script" path="res://scripts/ui/BaseLuckyPanel.gd" id="9"]')
    lines.append('')
    
    # Lucky Panel ext_resource
    for i, panel_name in enumerate(LUCKY_PANELS):
        res_id = base_count + 1 + i
        lines.append(f'[ext_resource type="Script" path="res://scripts/ui/{panel_name}.gd" id="{res_id}"]')
    lines.append('')
    
    # 主場景節點
    lines.append('[node name="Main" type="Node2D"]')
    lines.append('')
    lines.append('[node name="Camera2D" type="Camera2D" parent="."]')
    lines.append('position = Vector2(640, 360)')
    lines.append('')
    lines.append('[node name="BackgroundManager" type="Node2D" parent="."]')
    lines.append('script = ExtResource("4")')
    lines.append('')
    lines.append('[node name="TargetManager" type="Node2D" parent="."]')
    lines.append('script = ExtResource("1")')
    lines.append('')
    lines.append('[node name="Cannon" type="Node2D" parent="."]')
    lines.append('position = Vector2(640, 630)')
    lines.append('script = ExtResource("2")')
    lines.append('')
    lines.append('[node name="CannonSprite" type="Sprite2D" parent="Cannon"]')
    lines.append('scale = Vector2(2, 2)')
    lines.append('')
    lines.append('[node name="CharLabel" type="Label" parent="Cannon"]')
    lines.append('offset_left = -40.0')
    lines.append('offset_top = -60.0')
    lines.append('offset_right = 40.0')
    lines.append('offset_bottom = -40.0')
    lines.append('text = "Chiikawa"')
    lines.append('horizontal_alignment = 1')
    lines.append('')
    lines.append('[node name="CharacterAnimator" type="Node2D" parent="Cannon"]')
    lines.append('position = Vector2(0, -40)')
    lines.append('script = ExtResource("6")')
    lines.append('')
    
    # HUD
    lines.append('[node name="HUD" type="CanvasLayer" parent="."]')
    lines.append('script = ExtResource("3")')
    lines.append('')
    lines.append('[node name="TopBar" type="Control" parent="HUD"]')
    lines.append('layout_mode = 3')
    lines.append('anchors_preset = 15')
    lines.append('anchor_right = 1.0')
    lines.append('offset_bottom = 40.0')
    lines.append('grow_horizontal = 2')
    lines.append('grow_vertical = 2')
    lines.append('')
    lines.append('[node name="CoinsLabel" type="Label" parent="HUD/TopBar"]')
    lines.append('offset_left = 8.0')
    lines.append('offset_top = 8.0')
    lines.append('offset_right = 160.0')
    lines.append('offset_bottom = 36.0')
    lines.append('text = "💰 10000"')
    lines.append('add_theme_font_size_override = {"font_size": 14}')
    lines.append('')
    lines.append('[node name="BetLabel" type="Label" parent="HUD/TopBar"]')
    lines.append('offset_left = 170.0')
    lines.append('offset_top = 8.0')
    lines.append('offset_right = 340.0')
    lines.append('offset_bottom = 36.0')
    lines.append('text = "BET LV1 (1)"')
    lines.append('add_theme_font_size_override = {"font_size": 14}')
    lines.append('')
    lines.append('[node name="CharLabel" type="Label" parent="HUD/TopBar"]')
    lines.append('offset_left = 350.0')
    lines.append('offset_top = 8.0')
    lines.append('offset_right = 500.0')
    lines.append('offset_bottom = 36.0')
    lines.append('text = "Chiikawa"')
    lines.append('modulate = Color(1, 0.6, 0.8, 1)')
    lines.append('add_theme_font_size_override = {"font_size": 14}')
    lines.append('')
    lines.append('[node name="LaborBar" type="ProgressBar" parent="HUD/TopBar"]')
    lines.append('offset_left = 510.0')
    lines.append('offset_top = 10.0')
    lines.append('offset_right = 710.0')
    lines.append('offset_bottom = 30.0')
    lines.append('max_value = 100.0')
    lines.append('value = 0.0')
    lines.append('')
    lines.append('[node name="LaborLabel" type="Label" parent="HUD/TopBar"]')
    lines.append('offset_left = 720.0')
    lines.append('offset_top = 8.0')
    lines.append('offset_right = 820.0')
    lines.append('offset_bottom = 36.0')
    lines.append('text = "0/100"')
    lines.append('add_theme_font_size_override = {"font_size": 13}')
    lines.append('')
    lines.append('[node name="StateLabel" type="Label" parent="HUD/TopBar"]')
    lines.append('offset_left = 1000.0')
    lines.append('offset_top = 8.0')
    lines.append('offset_right = 1270.0')
    lines.append('offset_bottom = 36.0')
    lines.append('text = "NORMAL PLAY"')
    lines.append('horizontal_alignment = 2')
    lines.append('add_theme_font_size_override = {"font_size": 13}')
    lines.append('')
    lines.append('[node name="BottomBar" type="Control" parent="HUD"]')
    lines.append('layout_mode = 3')
    lines.append('anchors_preset = 12')
    lines.append('anchor_top = 1.0')
    lines.append('anchor_right = 1.0')
    lines.append('anchor_bottom = 1.0')
    lines.append('offset_top = -50.0')
    lines.append('grow_horizontal = 2')
    lines.append('grow_vertical = 0')
    lines.append('')
    lines.append('[node name="BetMinusBtn" type="Button" parent="HUD/BottomBar"]')
    lines.append('offset_left = 8.0')
    lines.append('offset_top = 8.0')
    lines.append('offset_right = 108.0')
    lines.append('offset_bottom = 42.0')
    lines.append('text = "BET -"')
    lines.append('')
    lines.append('[node name="BetPlusBtn" type="Button" parent="HUD/BottomBar"]')
    lines.append('offset_left = 118.0')
    lines.append('offset_top = 8.0')
    lines.append('offset_right = 218.0')
    lines.append('offset_bottom = 42.0')
    lines.append('text = "BET +"')
    lines.append('')
    lines.append('[node name="AutoBtn" type="Button" parent="HUD/BottomBar"]')
    lines.append('offset_left = 228.0')
    lines.append('offset_top = 8.0')
    lines.append('offset_right = 328.0')
    lines.append('offset_bottom = 42.0')
    lines.append('text = "AUTO"')
    lines.append('')
    lines.append('[node name="LockBtn" type="Button" parent="HUD/BottomBar"]')
    lines.append('offset_left = 338.0')
    lines.append('offset_top = 8.0')
    lines.append('offset_right = 438.0')
    lines.append('offset_bottom = 42.0')
    lines.append('text = "🔓 LOCK"')
    lines.append('')
    lines.append('[node name="BossBtn" type="Button" parent="HUD/BottomBar"]')
    lines.append('offset_left = 448.0')
    lines.append('offset_top = 8.0')
    lines.append('offset_right = 548.0')
    lines.append('offset_bottom = 42.0')
    lines.append('text = "BOSS"')
    lines.append('')
    lines.append('[node name="BonusBtn" type="Button" parent="HUD/BottomBar"]')
    lines.append('offset_left = 558.0')
    lines.append('offset_top = 8.0')
    lines.append('offset_right = 658.0')
    lines.append('offset_bottom = 42.0')
    lines.append('text = "BONUS"')
    lines.append('')
    
    # BonusGame
    lines.append('[node name="BonusGame" type="CanvasLayer" parent="."]')
    lines.append('script = ExtResource("5")')
    lines.append('')
    
    # LuckyEventSystem（CanvasLayer，layer=62）
    lines.append('[node name="LuckyEventSystem" type="CanvasLayer" parent="."]')
    lines.append('layer = 62')
    lines.append('script = ExtResource("7")')
    lines.append('')
    
    # LuckyPanelRegistry（普通 Node）
    lines.append('[node name="LuckyPanelRegistry" type="Node" parent="."]')
    lines.append('script = ExtResource("8")')
    lines.append('')
    
    # 所有 Lucky Panel 節點（CanvasLayer，掛在 LuckyPanelRegistry 下）
    for i, panel_name in enumerate(LUCKY_PANELS):
        res_id = base_count + 1 + i
        lines.append(f'[node name="{panel_name}" type="CanvasLayer" parent="LuckyPanelRegistry"]')
        lines.append(f'script = ExtResource("{res_id}")')
        lines.append('')
    
    return '\n'.join(lines)


if __name__ == '__main__':
    content = generate_main_tscn()
    output_path = r'd:\Kiro\client\chiikawa-pixel\scenes\Main.tscn'
    with open(output_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f"✅ Main.tscn 生成完成：{output_path}")
    print(f"   Lucky Panel 數量：{len(LUCKY_PANELS)}")
    print(f"   總 load_steps：{9 + len(LUCKY_PANELS)}")
