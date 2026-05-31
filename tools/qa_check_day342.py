"""
qa_check_day342.py — DAY-342 QA 驗證腳本
驗證 在線玩家數顯示 + 高倍率全服公告 + Server 編譯
"""
import os
import sys
import subprocess

PASS = "✅"
FAIL = "❌"
results = []

def check(name, condition, detail=""):
    status = PASS if condition else FAIL
    results.append((status, name, detail))
    if not condition:
        print(f"  {FAIL} {name}: {detail}")

# ── 1. Server protocol 包含 OnlineCount ──────────────────────
MESSAGES_GO = r"d:\Kiro\server\internal\protocol\messages.go"
with open(MESSAGES_GO, 'r', encoding='utf-8') as f:
    msg_content = f.read()

check("Protocol: online_count 欄位", "online_count" in msg_content)
check("Protocol: OnlineCount 欄位", "OnlineCount" in msg_content)

# ── 2. game.go 包含 OnlineCount 填入 ─────────────────────────
GAME_GO = r"d:\Kiro\server\internal\game\game.go"
with open(GAME_GO, 'r', encoding='utf-8') as f:
    game_content = f.read()

check("game.go: OnlineCount 填入", "OnlineCount:" in game_content)
check("game.go: len(g.players)", "len(g.players)" in game_content)
check("game.go: 高倍率全服公告", "高倍率擊破全服公告" in game_content)
check("game.go: 50x 公告門檻", "t.Multiplier >= 50" in game_content)
check("game.go: 100x 公告門檻", "t.Multiplier >= 100" in game_content)
check("game.go: 1000x 公告門檻", "t.Multiplier >= 1000" in game_content)
check("game.go: fmt import", '"fmt"' in game_content)

# ── 3. GameManager.gd 包含 get_online_count ──────────────────
GM_GD = r"d:\Kiro\client\chiikawa-pixel\scripts\game\GameManager.gd"
with open(GM_GD, 'r', encoding='utf-8') as f:
    gm_content = f.read()

check("GameManager: get_online_count 方法", "get_online_count" in gm_content)
check("GameManager: online_count 欄位讀取", "online_count" in gm_content)

# ── 4. HUD.gd 包含在線玩家數顯示 ─────────────────────────────
HUD_GD = r"d:\Kiro\client\chiikawa-pixel\scripts\ui\HUD.gd"
with open(HUD_GD, 'r', encoding='utf-8') as f:
    hud_content = f.read()

check("HUD: _online_label 變數", "_online_label" in hud_content)
check("HUD: _create_online_label 函數", "_create_online_label" in hud_content)
check("HUD: _update_online_display 函數", "_update_online_display" in hud_content)
check("HUD: 在線顯示文字", "在線" in hud_content)
check("HUD: 多人在線火焰", "🔥" in hud_content)
check("HUD: _update_ui 呼叫 _update_online_display", "_update_online_display()" in hud_content)
check("HUD: _create_online_label 在 _ready 中呼叫", "_create_online_label()" in hud_content)

# ── 5. Server 編譯狀態 ────────────────────────────────────────
result = subprocess.run(
    ["go", "build", "./..."],
    cwd=r"d:\Kiro\server",
    capture_output=True, text=True
)
check("Server: go build 通過", result.returncode == 0, result.stderr[:200] if result.stderr else "")

result2 = subprocess.run(
    ["go", "vet", "./..."],
    cwd=r"d:\Kiro\server",
    capture_output=True, text=True
)
check("Server: go vet 通過", result2.returncode == 0, result2.stderr[:200] if result2.stderr else "")

# ── 結果統計 ──────────────────────────────────────────────────
total = len(results)
passed = sum(1 for r in results if r[0] == PASS)
failed = total - passed

print(f"\n{'='*50}")
print(f"DAY-342 QA 結果：{passed}/{total} 通過")
print(f"{'='*50}")

if failed > 0:
    print("\n失敗項目：")
    for status, name, detail in results:
        if status == FAIL:
            print(f"  {FAIL} {name}: {detail}")
else:
    print("🎉 全部通過！")

sys.exit(0 if failed == 0 else 1)
