#!/usr/bin/env python3
"""
在 HUD.gd 中禁用所有有問題的 Panel preload。
策略：把 const XxxPanelScript = preload(...) 改為注釋。
同時把對應的 _init_xxx_panel() 函數體替換為 pass。
"""
import os
import re

HUD_PATH = r"d:\Kiro\client\chiikawa-pixel\scripts\ui\HUD.gd"

# 有問題的腳本列表（不含 .gd 後綴）
BROKEN_SCRIPTS = [
    "AbyssWhalePanel", "ActivityFeedPanel", "AnglerfishPanel", "AnnouncementPanel",
    "BlackHolePanel", "BombCrabPanel", "BuyBonusPanel", "CaptainFishPanel",
    "ChainExplosionPanel", "ChallengePanel", "ChallengePvPPanel", "CodexPanel",
    "CrocodilePanel", "DailyBossPanel", "DailySpinPanel", "DMPanel",
    "DragonWrathPanel", "ElectricJellyfishPanel", "EventPanel", "FestivalPanel",
    "FireStormPanel", "FreezeBombPanel", "FriendPanel", "GiantPrizeFishPanel",
    "GoldenSharkPanel", "GoldenTreasurePanel", "GoldenTurtlePanel",
    "GuildPanel", "GuildWarPanel", "HallOfFamePanel",
    "IceFishingPanel", "JackpotPanel", "LeaderboardPanel", "LionDancePanel",
    "LobbyManager", "LuckyCatchPanel", "LuckyDicePanel", "LuckyEggPanel",
    "LuckyFlagFishPanel", "LuckyGravityFlipPanel", "LuckyProphecyFishPanel",
    "LuckyRicochetFishPanel", "LuckyStarFishPanel", "LuckyTridentPanel",
    "MegaOctopusPanel", "MissionPanel", "MissionStreakPanel", "MoneyFishPanel",
    "MysteryBoxPanel", "PlayerCardPanel", "PlayerJourneyPanel", "PlayerStatsPanel",
    "RainbowLuckyPanel", "RainbowSharkPanel", "RapidRespinPanel", "RecommendPanel",
    "ReferralPanel", "RockSkeletonConcertPanel", "RoomSelectPanel", "RouletteCrabPanel",
    "RoulettePanel", "RoyalChainLightningPanel", "SchoolPanicPanel", "SeaAnemonePanel",
    "SeasonPanel", "SessionStatsPanel", "ShopPanel", "SkinPanel",
    "SpecialWeaponPanel", "StreakPanel", "TitlePanel", "TournamentPanel",
    "TreasureMapPanel", "TripleLuckyFishPanel", "UnluckyBonusPanel", "VIPPanel",
    "VortexFishPanel", "WeaponPanel", "WeatherPanel", "WeatherSurgePanel", "WheelPanel",
]

with open(HUD_PATH, 'r', encoding='utf-8') as f:
    lines = f.readlines()

modified = 0
i = 0
while i < len(lines):
    line = lines[i]
    
    # 找到 const XxxPanelScript = preload(...) 行
    for script_name in BROKEN_SCRIPTS:
        const_name = f"{script_name}Script"
        if f"const {const_name}" in line and "preload" in line:
            lines[i] = f"# DISABLED: {line.rstrip()}\n"
            modified += 1
            print(f"Disabled preload: {const_name} at line {i+1}")
            break
    
    i += 1

# 現在找到所有使用這些 const 的 func 並替換為 pass
# 策略：找到 func _init_xxx_panel() 並把函數體替換為 pass
i = 0
in_broken_func = False
func_indent = 0

while i < len(lines):
    line = lines[i]
    stripped = line.strip()
    
    # 找到 func _init_xxx_panel() 定義
    if stripped.startswith("func ") and stripped.endswith(":"):
        func_name_match = re.match(r'func\s+(\w+)\s*\(', stripped)
        if func_name_match:
            func_name = func_name_match.group(1)
            # 檢查這個函數是否使用了損壞的 Panel
            # 往後看幾行
            uses_broken = False
            for j in range(i+1, min(i+5, len(lines))):
                for script_name in BROKEN_SCRIPTS:
                    const_name = f"{script_name}Script"
                    if const_name in lines[j]:
                        uses_broken = True
                        break
                if uses_broken:
                    break
            
            if uses_broken:
                # 找到函數體的範圍
                func_start = i
                func_end = i + 1
                while func_end < len(lines):
                    next_line = lines[func_end]
                    next_stripped = next_line.strip()
                    if next_stripped and not next_line.startswith('\t') and not next_line.startswith(' '):
                        break
                    func_end += 1
                
                # 替換函數體為 pass
                func_header = lines[func_start]
                new_func = [func_header, '\tpass  # disabled: broken panel\n']
                for j in range(func_start + 1, func_end):
                    lines[j] = ''
                lines[func_start] = func_header
                if func_start + 1 < len(lines):
                    lines[func_start + 1] = '\tpass  # disabled: broken panel\n'
                
                modified += 1
                print(f"Disabled func: {func_name} at line {func_start+1}")
    
    i += 1

with open(HUD_PATH, 'w', encoding='utf-8') as f:
    f.writelines(lines)

print(f"\nTotal modifications: {modified}")
