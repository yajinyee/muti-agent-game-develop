#!/usr/bin/env python3
"""
QA Check DAY-308
驗證 T146-T150 Lucky badge 修復 + Agent 文件補齊 + 整體狀態
"""
import os
import sys

ROOT = r"d:\Kiro"
CLIENT = os.path.join(ROOT, "client", "chiikawa-pixel")
SCRIPTS = os.path.join(CLIENT, "scripts")
SPRITES = os.path.join(CLIENT, "assets", "sprites", "targets")
AGENTS = os.path.join(ROOT, "agents")

results = []
passed = 0
failed = 0

def check(name, condition, detail=""):
    global passed, failed
    if condition:
        print(f"  ✅ {name}")
        passed += 1
    else:
        print(f"  ❌ {name}" + (f" — {detail}" if detail else ""))
        failed += 1
    results.append((name, condition))

print("=" * 60)
print("QA Check DAY-308")
print("=" * 60)

# 1. Server 編譯（用 go build 輸出確認）
print("\n[1] Server 編譯狀態")
import subprocess
r = subprocess.run(["go", "build", "./..."], cwd=os.path.join(ROOT, "server"), capture_output=True, text=True)
check("go build ./...", r.returncode == 0, r.stderr[:200] if r.returncode != 0 else "")
r2 = subprocess.run(["go", "vet", "./..."], cwd=os.path.join(ROOT, "server"), capture_output=True, text=True)
check("go vet ./...", r2.returncode == 0, r2.stderr[:200] if r2.returncode != 0 else "")

# 2. T146-T150 精靈圖
print("\n[2] T146-T150 精靈圖")
for tid in range(146, 151):
    files = [f for f in os.listdir(SPRITES) if f.startswith(f"T{tid}_")]
    check(f"T{tid} 精靈圖存在", len(files) > 0, f"找不到 T{tid}_*.png")

# 3. T146-T150 Lucky Panel 腳本
print("\n[3] T146-T150 Lucky Panel 腳本")
ui_dir = os.path.join(SCRIPTS, "ui")
panel_map = {
    146: "LuckyQuantumPanel.gd",
    147: "LuckySupernovaPanel.gd",
    148: "LuckyInfinitePanel.gd",
    149: "LuckyGenesisPanel.gd",
    150: "LuckyRebirthPanel.gd",
}
for tid, fname in panel_map.items():
    path = os.path.join(ui_dir, fname)
    check(f"T{tid} {fname}", os.path.exists(path))

# 4. TargetManager Lucky badge 範圍
print("\n[4] TargetManager Lucky badge 範圍")
tm_path = os.path.join(SCRIPTS, "game", "TargetManager.gd")
with open(tm_path, "r", encoding="utf-8") as f:
    tm_content = f.read()
check("Lucky badge 覆蓋 T106-T150", "tid_num >= 106 and tid_num <= 150" in tm_content)
check("T146-T150 Sprite 路徑存在", '"T150"' in tm_content)

# 5. GameManager T146-T150 訊號
print("\n[5] GameManager T146-T150 訊號")
gm_path = os.path.join(SCRIPTS, "game", "GameManager.gd")
with open(gm_path, "r", encoding="utf-8") as f:
    gm_content = f.read()
for sig in ["lucky_quantum", "lucky_supernova", "lucky_infinite", "lucky_genesis", "lucky_rebirth"]:
    check(f"訊號 {sig}", f"signal {sig}" in gm_content)

# 6. HUD T146-T150 事件處理
print("\n[6] HUD T146-T150 事件處理")
hud_path = os.path.join(SCRIPTS, "ui", "HUD.gd")
with open(hud_path, "r", encoding="utf-8") as f:
    hud_content = f.read()
for handler in ["_on_lucky_quantum", "_on_lucky_supernova", "_on_lucky_infinite", "_on_lucky_genesis", "_on_lucky_rebirth"]:
    check(f"HUD {handler}", handler in hud_content)

# 7. Agent 文件補齊
print("\n[7] Agent 文件補齊")
required_agents = [
    "target-design-agent.md",
    "spec-architect.md",
    "server-combat-agent.md",
    "server-event-agent.md",
    "server-infra-agent.md",
    "target-system-agent.md",
    "game-state-agent.md",
    "social-ui-agent.md",
    "screen-recorder-agent.md",
    "screen-effect-agent.md",
    "network-agent.md",
    "sfx-agent.md",
    "target-pixel-agent.md",
    "target-ai-agent.md",
    "ui-art-agent.md",
    "qa-playtest-agent.md",
    "video-analysis-agent.md",
    "research-agent.md",
    "skill-librarian.md",
]
for agent in required_agents:
    path = os.path.join(AGENTS, agent)
    check(f"Agent {agent}", os.path.exists(path))

# 8. 音效和 BGM 檔案
print("\n[8] 音效和 BGM 檔案")
sfx_dir = os.path.join(CLIENT, "assets", "audio", "sfx")
bgm_dir = os.path.join(CLIENT, "assets", "audio", "bgm")
sfx_files = os.listdir(sfx_dir) if os.path.exists(sfx_dir) else []
bgm_files = os.listdir(bgm_dir) if os.path.exists(bgm_dir) else []
check("SFX 檔案 >= 10 個", len(sfx_files) >= 10, f"只有 {len(sfx_files)} 個")
check("BGM 檔案 >= 4 個", len(bgm_files) >= 4, f"只有 {len(bgm_files)} 個")

# 9. 角色精靈圖
print("\n[9] 角色精靈圖")
char_dir = os.path.join(CLIENT, "assets", "sprites", "characters")
for char in ["chiikawa", "hachiware", "usagi"]:
    for state in ["idle", "attack", "bigwin"]:
        fname = f"{char}_{state}.png"
        path = os.path.join(char_dir, fname)
        check(f"{fname}", os.path.exists(path))

# 10. 目標物數量
print("\n[10] 目標物數量")
target_files = [f for f in os.listdir(SPRITES) if f.endswith(".png")]
check("目標物 >= 57 個", len(target_files) >= 57, f"只有 {len(target_files)} 個")

# 結果
print("\n" + "=" * 60)
print(f"結果：{passed} 通過 / {failed} 失敗 / {passed + failed} 總計")
print("=" * 60)

if failed == 0:
    print("🎉 全部通過！可以推送 GitHub。")
else:
    print(f"⚠️  有 {failed} 項失敗，請修復後再推送。")

sys.exit(0 if failed == 0 else 1)
