"""
qa_check_day323.py — DAY-323 品質驗證腳本
驗證 T206-T210 五個新 Lucky 魚系統的完整性
"""
import os
import sys

errors = []
passed = 0

def check(condition, msg):
    global passed
    if condition:
        passed += 1
    else:
        errors.append(f"❌ {msg}")

# ── 1. Server Handler 檔案存在 ────────────────────────────────
server_handlers = [
    "d:/Kiro/server/internal/game/lucky_fever_boost_handler.go",
    "d:/Kiro/server/internal/game/lucky_guild_battle_handler.go",
    "d:/Kiro/server/internal/game/lucky_path_fish_handler.go",
    "d:/Kiro/server/internal/game/lucky_chain_eel_handler.go",
    "d:/Kiro/server/internal/game/lucky_ultimate_miracle_handler.go",
]
for f in server_handlers:
    check(os.path.exists(f), f"Server handler 不存在：{os.path.basename(f)}")

# ── 2. tables.go 包含 T206-T210 ───────────────────────────────
tables_path = "d:/Kiro/server/internal/data/tables.go"
with open(tables_path, encoding="utf-8") as f:
    tables_content = f.read()
for tid in ["T206", "T207", "T208", "T209", "T210"]:
    check(f'ID: "{tid}"' in tables_content, f"tables.go 缺少 {tid}")

# ── 3. messages.go 包含新訊息類型 ─────────────────────────────
messages_path = "d:/Kiro/server/internal/protocol/messages.go"
with open(messages_path, encoding="utf-8") as f:
    messages_content = f.read()
new_msgs = [
    "MsgLuckyFeverBoost",
    "MsgLuckyGuildBattle",
    "MsgLuckyPathFish",
    "MsgLuckyChainEel",
    "MsgLuckyUltimateMiracle",
]
for msg in new_msgs:
    check(msg in messages_content, f"messages.go 缺少 {msg}")

# ── 4. game.go 包含新 manager 定義 ────────────────────────────
game_path = "d:/Kiro/server/internal/game/game.go"
with open(game_path, encoding="utf-8") as f:
    game_content = f.read()
new_managers = [
    "luckyFeverBoost",
    "luckyGuildBattle",
    "luckyPathFish",
    "luckyChainEel",
    "luckyUltimateMiracle",
]
for mgr in new_managers:
    check(mgr in game_content, f"game.go 缺少 manager：{mgr}")

# ── 5. game.go 包含新 switch case ─────────────────────────────
new_cases = [
    "isLuckyFeverBoostFish",
    "isLuckyGuildBattleFish",
    "isLuckyPathFish",
    "isLuckyChainEelFish",
    "isLuckyUltimateMiracleFish",
]
for case in new_cases:
    check(case in game_content, f"game.go 缺少 switch case：{case}")

# ── 6. Client Lucky Panel 檔案存在 ────────────────────────────
client_panels = [
    "d:/Kiro/client/chiikawa-pixel/scripts/ui/LuckyFeverBoostPanel.gd",
    "d:/Kiro/client/chiikawa-pixel/scripts/ui/LuckyGuildBattlePanel.gd",
    "d:/Kiro/client/chiikawa-pixel/scripts/ui/LuckyPathFishPanel.gd",
    "d:/Kiro/client/chiikawa-pixel/scripts/ui/LuckyChainEelPanel.gd",
    "d:/Kiro/client/chiikawa-pixel/scripts/ui/LuckyUltimateMiraclePanel.gd",
]
for f in client_panels:
    check(os.path.exists(f), f"Client Panel 不存在：{os.path.basename(f)}")

# ── 7. GameManager.gd 包含新訊號 ─────────────────────────────
gm_path = "d:/Kiro/client/chiikawa-pixel/scripts/game/GameManager.gd"
with open(gm_path, encoding="utf-8") as f:
    gm_content = f.read()
new_signals = [
    "lucky_fever_boost",
    "lucky_guild_battle",
    "lucky_path_fish",
    "lucky_chain_eel",
    "lucky_ultimate_miracle",
]
for sig in new_signals:
    check(f"signal {sig}" in gm_content, f"GameManager.gd 缺少訊號：{sig}")
    check(f'emit_signal("{sig}"' in gm_content, f"GameManager.gd 缺少 emit_signal：{sig}")

# ── 8. LuckyPanelRegistry.gd 包含新 Panel 映射 ───────────────
registry_path = "d:/Kiro/client/chiikawa-pixel/scripts/ui/LuckyPanelRegistry.gd"
with open(registry_path, encoding="utf-8") as f:
    registry_content = f.read()
new_panel_mappings = [
    "LuckyFeverBoostPanel",
    "LuckyGuildBattlePanel",
    "LuckyPathFishPanel",
    "LuckyChainEelPanel",
    "LuckyUltimateMiraclePanel",
]
for panel in new_panel_mappings:
    check(panel in registry_content, f"LuckyPanelRegistry.gd 缺少 Panel 映射：{panel}")

# ── 9. TargetManager.gd 包含 T206-T210 ───────────────────────
tm_path = "d:/Kiro/client/chiikawa-pixel/scripts/game/TargetManager.gd"
with open(tm_path, encoding="utf-8") as f:
    tm_content = f.read()
for tid in ["T206", "T207", "T208", "T209", "T210"]:
    check(tid in tm_content, f"TargetManager.gd 缺少 {tid}")

# ── 10. 精靈圖存在 ────────────────────────────────────────────
sprites = [
    "d:/Kiro/client/chiikawa-pixel/assets/sprites/targets/T206_fever_boost.png",
    "d:/Kiro/client/chiikawa-pixel/assets/sprites/targets/T207_guild_battle.png",
    "d:/Kiro/client/chiikawa-pixel/assets/sprites/targets/T208_path_fish.png",
    "d:/Kiro/client/chiikawa-pixel/assets/sprites/targets/T209_chain_eel.png",
    "d:/Kiro/client/chiikawa-pixel/assets/sprites/targets/T210_ultimate_miracle.png",
]
for f in sprites:
    check(os.path.exists(f), f"精靈圖不存在：{os.path.basename(f)}")

# ── 11. Handler 函數名稱正確 ──────────────────────────────────
handler_checks = [
    ("d:/Kiro/server/internal/game/lucky_fever_boost_handler.go", "isLuckyFeverBoostFish"),
    ("d:/Kiro/server/internal/game/lucky_guild_battle_handler.go", "isLuckyGuildBattleFish"),
    ("d:/Kiro/server/internal/game/lucky_path_fish_handler.go", "isLuckyPathFish"),
    ("d:/Kiro/server/internal/game/lucky_chain_eel_handler.go", "isLuckyChainEelFish"),
    ("d:/Kiro/server/internal/game/lucky_ultimate_miracle_handler.go", "isLuckyUltimateMiracleFish"),
]
for filepath, func_name in handler_checks:
    with open(filepath, encoding="utf-8") as f:
        content = f.read()
    check(func_name in content, f"Handler 缺少函數：{func_name}")

# ── 12. 倍率設計驗證 ──────────────────────────────────────────
# T206-T210 倍率應該遞增
multipliers = {
    "T206": 9500,
    "T207": 10000,
    "T208": 11000,
    "T209": 12000,
    "T210": 16888,
}
for tid, expected_mult in multipliers.items():
    check(f'Multiplier: {expected_mult}' in tables_content, f"tables.go {tid} 倍率不正確（應為 {expected_mult}）")

# ── 結果 ──────────────────────────────────────────────────────
total = passed + len(errors)
print(f"\n{'='*50}")
print(f"DAY-323 QA 驗證結果：{passed}/{total} 通過")
print(f"{'='*50}")
if errors:
    print("\n❌ 失敗項目：")
    for e in errors:
        print(f"  {e}")
    sys.exit(1)
else:
    print("\n✅ 全部通過！DAY-323 T206-T210 五個新 Lucky 魚系統完整！")
    print(f"\n📊 DAY-323 統計：")
    print(f"  - 新增目標物：T206-T210（5 個）")
    print(f"  - 總目標物數量：117 種（T001-T006 + T101-T210 + B001）")
    print(f"  - Lucky 系統數量：105 個（T106-T210）")
    print(f"  - 最高全服倍率：T210 終極奇蹟 ×35.0（新史上最高）")
    print(f"  - 最高個人倍率：T184 風險等級 ×3000")
    print(f"  - 最高 Jackpot：T174 Grand Jackpot 5000x 起跳")
    print(f"  - 業界依據：Games Global Fever Boost™ + Fishing Frenzy Guild Wars + Fish Road + Royal Fishing Chain Eel")
