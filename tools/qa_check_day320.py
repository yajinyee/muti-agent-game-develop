"""
qa_check_day320.py — DAY-320 QA 驗證腳本
驗證 HUD.gd 訊號連接完整性（DAY-318/319 補齊）+ Agent 文件更新
"""
import os
import re

PASS = 0
FAIL = 0

def check(name, condition, detail=""):
    global PASS, FAIL
    if condition:
        PASS += 1
        print(f"  ✅ {name}")
    else:
        FAIL += 1
        print(f"  ❌ {name}" + (f" — {detail}" if detail else ""))

BASE = r"d:\Kiro"
CLIENT_UI = os.path.join(BASE, "client", "chiikawa-pixel", "scripts", "ui")
CLIENT_GAME = os.path.join(BASE, "client", "chiikawa-pixel", "scripts", "game")
AGENTS = os.path.join(BASE, "agents")

print("=== DAY-320 QA 驗證 ===\n")

# ── HUD.gd 訊號連接驗證 ──────────────────────────────────────
print("【HUD.gd 訊號連接 — DAY-318 補齊】")
hud_path = os.path.join(CLIENT_UI, "HUD.gd")
hud_content = open(hud_path, encoding="utf-8").read()

day318_signals = [
    "lucky_dragon_king",
    "lucky_eternal_cycle",
    "lucky_chaos_explosion",
    "lucky_divine_revival",
    "lucky_genesis_epoch",
]
for sig in day318_signals:
    check(f"HUD.gd 連接 {sig}", f"GameManager.{sig}.connect" in hud_content)

print("\n【HUD.gd 訊號連接 — DAY-319 補齊】")
day319_signals = [
    "lucky_energy_storm",
    "lucky_crystal_resonance",
    "lucky_fate_judgment",
    "lucky_time_reversal",
    "lucky_cosmic_singularity",
]
for sig in day319_signals:
    check(f"HUD.gd 連接 {sig}", f"GameManager.{sig}.connect" in hud_content)

# ── HUD.gd 處理函數驗證 ──────────────────────────────────────
print("\n【HUD.gd 處理函數 — DAY-318 補齊】")
for sig in day318_signals:
    check(f"HUD.gd 有 _on_{sig}", f"func _on_{sig}" in hud_content)

print("\n【HUD.gd 處理函數 — DAY-319 補齊】")
for sig in day319_signals:
    check(f"HUD.gd 有 _on_{sig}", f"func _on_{sig}" in hud_content)

# ── HUD.gd 統計 ──────────────────────────────────────────────
print("\n【HUD.gd 統計】")
connect_count = len(re.findall(r"GameManager\.lucky_\w+\.connect", hud_content))
handler_count = len(re.findall(r"func _on_lucky_\w+", hud_content))
line_count = hud_content.count("\n")
check(f"HUD.gd 訊號連接數量 >= 90", connect_count >= 90, f"實際: {connect_count}")
check(f"HUD.gd 處理函數數量 >= 90", handler_count >= 90, f"實際: {handler_count}")
print(f"  ℹ️  HUD.gd 行數: {line_count}")
print(f"  ℹ️  訊號連接數: {connect_count}")
print(f"  ℹ️  處理函數數: {handler_count}")

# ── GameManager.gd 訊號定義驗證 ──────────────────────────────
print("\n【GameManager.gd 訊號定義】")
gm_path = os.path.join(CLIENT_GAME, "GameManager.gd")
gm_content = open(gm_path, encoding="utf-8").read()

all_signals = day318_signals + day319_signals
for sig in all_signals:
    check(f"GameManager 有訊號 {sig}", f"signal {sig}" in gm_content)
    check(f"GameManager emit {sig}", f'emit_signal("{sig}"' in gm_content)

# ── LuckyPanelRegistry.gd 映射驗證 ──────────────────────────
print("\n【LuckyPanelRegistry.gd 映射】")
registry_path = os.path.join(CLIENT_UI, "LuckyPanelRegistry.gd")
registry_content = open(registry_path, encoding="utf-8").read()

for sig in all_signals:
    check(f"Registry 有映射 {sig}", f'"{sig}"' in registry_content)

# ── Panel 文件存在驗證 ────────────────────────────────────────
print("\n【Lucky Panel 文件存在】")
panel_map = {
    "lucky_dragon_king": "LuckyDragonKingPanel.gd",
    "lucky_eternal_cycle": "LuckyEternalCyclePanel.gd",
    "lucky_chaos_explosion": "LuckyChaosExplosionPanel.gd",
    "lucky_divine_revival": "LuckyDivineRevivalPanel.gd",
    "lucky_genesis_epoch": "LuckyGenesisEpochPanel.gd",
    "lucky_energy_storm": "LuckyEnergyStormPanel.gd",
    "lucky_crystal_resonance": "LuckyCrystalResonancePanel.gd",
    "lucky_fate_judgment": "LuckyFateJudgmentPanel.gd",
    "lucky_time_reversal": "LuckyTimeReversalPanel.gd",
    "lucky_cosmic_singularity": "LuckyCosmicSingularityPanel.gd",
}
for sig, panel_file in panel_map.items():
    path = os.path.join(CLIENT_UI, panel_file)
    check(f"Panel 文件存在: {panel_file}", os.path.exists(path))

# ── Agent 文件更新驗證 ────────────────────────────────────────
print("\n【Agent 文件更新】")
server_event_path = os.path.join(AGENTS, "server-event-agent.md")
lucky_panel_path = os.path.join(AGENTS, "lucky-panel-agent.md")

se_content = open(server_event_path, encoding="utf-8").read()
lp_content = open(lucky_panel_path, encoding="utf-8").read()

check("server-event-agent.md 提到 100 個 Lucky 系統", "100 個" in se_content or "T106-T205" in se_content)
check("server-event-agent.md 提到 T205 宇宙奇點", "T205" in se_content)
check("lucky-panel-agent.md 提到 LuckyPanelRegistry", "LuckyPanelRegistry" in lp_content)
check("lucky-panel-agent.md 有 Checklist", "Checklist" in lp_content)

# ── knowhow-log 更新驗證 ──────────────────────────────────────
print("\n【knowhow-log 更新】")
knowhow_path = os.path.join(BASE, ".kiro", "skills", "knowhow-log.md")
kh_content = open(knowhow_path, encoding="utf-8").read()
check("knowhow-log 有條目 30（Lucky 訊號連接遺漏）", "Lucky 訊號連接遺漏" in kh_content or "30." in kh_content)
check("knowhow-log 有條目 31（三層架構同步）", "三層架構同步" in kh_content or "31." in kh_content)

# ── 總結 ──────────────────────────────────────────────────────
print(f"\n{'='*40}")
print(f"DAY-320 QA 結果：{PASS} 通過 / {FAIL} 失敗")
total = PASS + FAIL
pct = PASS / total * 100 if total > 0 else 0
print(f"通過率：{pct:.1f}%")
if FAIL == 0:
    print("🎉 全部通過！DAY-320 整合完整性驗證成功！")
else:
    print(f"⚠️  有 {FAIL} 項需要修復")
