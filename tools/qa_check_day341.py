"""
qa_check_day341.py — DAY-341 QA 驗證腳本
驗證 Combo 里程碑音效 + Combo UI 升級
"""
import os
import sys

PASS = "✅"
FAIL = "❌"
results = []

def check(name, condition, detail=""):
    status = PASS if condition else FAIL
    results.append((status, name, detail))
    if not condition:
        print(f"  {FAIL} {name}: {detail}")

# ── 1. Combo 音效檔案存在 ──────────────────────────────────────
SFX_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\audio\sfx"
combo_sfx = ["combo_5.wav", "combo_10.wav", "combo_20.wav", "combo_30.wav"]
for sfx in combo_sfx:
    path = os.path.join(SFX_DIR, sfx)
    check(f"音效存在: {sfx}", os.path.exists(path), f"路徑: {path}")

# ── 2. Combo 音效 .import 存在 ────────────────────────────────
for sfx in combo_sfx:
    path = os.path.join(SFX_DIR, sfx + ".import")
    check(f".import 存在: {sfx}.import", os.path.exists(path))

# ── 3. Combo 音效大小合理（> 1KB）────────────────────────────
for sfx in combo_sfx:
    path = os.path.join(SFX_DIR, sfx)
    if os.path.exists(path):
        size = os.path.getsize(path)
        check(f"音效大小合理: {sfx}", size > 1000, f"{size} bytes")

# ── 4. AudioManager.gd 包含 Combo 音效枚舉 ───────────────────
AUDIO_MGR = r"d:\Kiro\client\chiikawa-pixel\scripts\game\AudioManager.gd"
with open(AUDIO_MGR, 'r', encoding='utf-8') as f:
    audio_content = f.read()

check("AudioManager: COMBO_5 枚舉", "COMBO_5" in audio_content)
check("AudioManager: COMBO_10 枚舉", "COMBO_10" in audio_content)
check("AudioManager: COMBO_20 枚舉", "COMBO_20" in audio_content)
check("AudioManager: COMBO_30 枚舉", "COMBO_30" in audio_content)
check("AudioManager: combo_5.wav 路徑", "combo_5.wav" in audio_content)
check("AudioManager: combo_10.wav 路徑", "combo_10.wav" in audio_content)
check("AudioManager: combo_20.wav 路徑", "combo_20.wav" in audio_content)
check("AudioManager: combo_30.wav 路徑", "combo_30.wav" in audio_content)
check("AudioManager: play_combo_milestone 函數", "play_combo_milestone" in audio_content)

# ── 5. HUD.gd 包含 Combo 里程碑功能 ──────────────────────────
HUD_GD = r"d:\Kiro\client\chiikawa-pixel\scripts\ui\HUD.gd"
with open(HUD_GD, 'r', encoding='utf-8') as f:
    hud_content = f.read()

check("HUD: play_combo_milestone 呼叫", "play_combo_milestone" in hud_content)
check("HUD: _spawn_combo_milestone_effect 函數", "_spawn_combo_milestone_effect" in hud_content)
check("HUD: ScreenShake.add_trauma 在 combo 里程碑", "ScreenShake.add_trauma(0.2 + (combo / 30.0)" in hud_content)
check("HUD: Combo 面板背景", "ComboPanel" in hud_content)
check("HUD: 里程碑文字 MAX COMBO", "MAX COMBO" in hud_content)
check("HUD: 里程碑文字 COMBO x20", "COMBO x20" in hud_content)
check("HUD: combo 粒子特效", "particle_count" in hud_content)
check("HUD: _last_combo 重置（combo < 5 時）", "_last_combo = combo" in hud_content)

# ── 6. Server 端 Combo 系統完整 ──────────────────────────────
PLAYER_GO = r"d:\Kiro\server\internal\game\player.go"
with open(PLAYER_GO, 'r', encoding='utf-8') as f:
    player_content = f.read()

check("Server: ComboLevels 定義", "ComboLevels" in player_content)
check("Server: AddCombo 函數", "AddCombo" in player_content)
check("Server: Combo 5 等級", "{5," in player_content)
check("Server: Combo 10 等級", "{10," in player_content)
check("Server: Combo 20 等級", "{20," in player_content)
check("Server: Combo 30 等級", "{30," in player_content)

# ── 7. Protocol 包含 combo_count ─────────────────────────────
MESSAGES_GO = r"d:\Kiro\server\internal\protocol\messages.go"
with open(MESSAGES_GO, 'r', encoding='utf-8') as f:
    msg_content = f.read()

check("Protocol: combo_count 欄位", "combo_count" in msg_content)
check("Protocol: combo_mult_bonus 欄位", "combo_mult_bonus" in msg_content)

# ── 8. GameManager.gd 包含 combo 方法 ────────────────────────
GM_GD = r"d:\Kiro\client\chiikawa-pixel\scripts\game\GameManager.gd"
with open(GM_GD, 'r', encoding='utf-8') as f:
    gm_content = f.read()

check("GameManager: get_combo_count 方法", "get_combo_count" in gm_content)
check("GameManager: get_combo_mult_bonus 方法", "get_combo_mult_bonus" in gm_content)

# ── 9. Server 編譯狀態 ────────────────────────────────────────
import subprocess
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
print(f"DAY-341 QA 結果：{passed}/{total} 通過")
print(f"{'='*50}")

if failed > 0:
    print("\n失敗項目：")
    for status, name, detail in results:
        if status == FAIL:
            print(f"  {FAIL} {name}: {detail}")
else:
    print("🎉 全部通過！")

sys.exit(0 if failed == 0 else 1)
