"""
qa_check_day335.py — DAY-335 品質驗證
驗證項目：
1. T001-T006 視覺升級（6 個 PNG 存在且尺寸正確）
2. HUDLuckySignals.gd 建立（拆分 HUD.gd 技術債）
3. Server 編譯狀態
4. 缺失 Agent 文件
"""
import os
import subprocess
import sys

PASS = 0
FAIL = 0

def check(name, condition, detail=""):
    global PASS, FAIL
    if condition:
        print(f"  ✅ {name}")
        PASS += 1
    else:
        print(f"  ❌ {name}" + (f" — {detail}" if detail else ""))
        FAIL += 1

def check_file_exists(path):
    return os.path.exists(path)

def check_file_size(path, min_bytes=100):
    return os.path.exists(path) and os.path.getsize(path) >= min_bytes

print("=== DAY-335 QA 驗證 ===\n")

# ── 1. T001-T006 視覺升級 ──────────────────────────────────────
print("【1】T001-T006 基礎目標物視覺升級")
targets_dir = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
basic_targets = ["T001_grass", "T002_bug_g", "T003_bug_r", "T004_bug_b", "T005_pudding", "T006_mushroom"]
for t in basic_targets:
    path = os.path.join(targets_dir, f"{t}.png")
    check(f"{t}.png 存在", check_file_exists(path))
    check(f"{t}.png 大小 >= 100 bytes", check_file_size(path, 100))

# 驗證圖片尺寸
try:
    from PIL import Image
    for t in basic_targets:
        path = os.path.join(targets_dir, f"{t}.png")
        if os.path.exists(path):
            img = Image.open(path)
            check(f"{t}.png 尺寸 64x64", img.size == (64, 64), f"實際: {img.size}")
except ImportError:
    print("  ⚠️ PIL 不可用，跳過尺寸驗證")

# ── 2. HUDLuckySignals.gd ──────────────────────────────────────
print("\n【2】HUDLuckySignals.gd 拆分")
hud_signals_path = r"d:\Kiro\client\chiikawa-pixel\scripts\ui\HUDLuckySignals.gd"
check("HUDLuckySignals.gd 存在", check_file_exists(hud_signals_path))
check("HUDLuckySignals.gd 大小 >= 1KB", check_file_size(hud_signals_path, 1024))

if os.path.exists(hud_signals_path):
    content = open(hud_signals_path, encoding="utf-8").read()
    check("包含 connect_all_lucky_signals", "connect_all_lucky_signals" in content)
    check("包含 DAY-292 訊號連接", "lucky_chain_lightning" in content)
    check("包含 DAY-296 訊號連接", "lucky_lucky_wheel" in content)
    check("包含 DAY-303 訊號連接", "lucky_crash_fish" in content)
    check("包含 _show_banner 委派", "_show_banner" in content)
    check("包含 _show_event 委派", "_show_event" in content)
    check("包含 _show_reward 委派", "_show_reward" in content)

# ── 3. Server 編譯 ──────────────────────────────────────────────
print("\n【3】Server 編譯狀態")
try:
    result = subprocess.run(
        ["go", "build", "./..."],
        cwd=r"d:\Kiro\server",
        capture_output=True, text=True, timeout=60
    )
    check("go build ./... 通過", result.returncode == 0, result.stderr[:200] if result.stderr else "")
    
    result2 = subprocess.run(
        ["go", "vet", "./..."],
        cwd=r"d:\Kiro\server",
        capture_output=True, text=True, timeout=60
    )
    check("go vet ./... 通過", result2.returncode == 0, result2.stderr[:200] if result2.stderr else "")
except Exception as e:
    check("Server 編譯", False, str(e))

# ── 4. 缺失 Agent 文件 ──────────────────────────────────────────
print("\n【4】Agent 文件完整性")
agents_dir = r"d:\Kiro\agents"
required_agents = [
    "target-design-agent.md",
    "visual-clarity-agent.md",
    "spec-architect.md",
    "server-core-agent.md",
    "server-combat-agent.md",
    "server-event-agent.md",
    "server-infra-agent.md",
    "cannon-agent.md",
    "target-system-agent.md",
    "boss-battle-agent.md",
    "bonus-game-agent.md",
    "game-state-agent.md",
    "hud-core-agent.md",
    "lucky-panel-agent.md",
    "hit-effect-agent.md",
    "screen-effect-agent.md",
    "environment-agent.md",
    "network-agent.md",
    "protocol-sync-agent.md",
    "integration-test-agent.md",
    "regression-guard-agent.md",
    "build-export-agent.md",
    "performance-agent.md",
    "qa-playtest-agent.md",
    "player-experience-agent.md",
    "balance-agent.md",
    "game-director.md",
]
for agent in required_agents:
    path = os.path.join(agents_dir, agent)
    check(f"{agent} 存在", check_file_exists(path))

# ── 5. 知識庫更新 ──────────────────────────────────────────────
print("\n【5】知識庫狀態")
knowhow_path = r"d:\Kiro\.kiro\skills\knowhow-log.md"
check("knowhow-log.md 存在", check_file_exists(knowhow_path))
if os.path.exists(knowhow_path):
    content = open(knowhow_path, encoding="utf-8").read()
    check("knowhow-log 有內容", len(content) > 500)

# ── 結果 ──────────────────────────────────────────────────────
print(f"\n=== 結果：{PASS} 通過 / {FAIL} 失敗 / 共 {PASS+FAIL} 項 ===")
if FAIL == 0:
    print("🎉 全部通過！DAY-335 品質驗證完成")
else:
    print(f"⚠️ {FAIL} 項需要修復")
    sys.exit(1)
