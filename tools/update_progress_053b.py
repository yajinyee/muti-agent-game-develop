#!/usr/bin/env python3
"""更新 progress.md 的最後更新行"""

path = r"d:\Kiro\docs\progress.md"

with open(path, "r", encoding="utf-8") as f:
    content = f.read()

old = "## 最後更新：2026-05-20（DAY-053 HUD.gd 大型腳本拆分 + 三個獨立面板腳本）"
new = "## 最後更新：2026-05-20（DAY-053b AudioManager 重構 + Audio Sync 100/100 達成）"

new_content = content.replace(old, new)

# 也更新架構成熟度
old2 = "HUD 模組化（JackpotPanel/MissionPanel/SessionStatsPanel 獨立腳本）**"
new2 = "HUD 模組化（JackpotPanel/MissionPanel/SessionStatsPanel 獨立腳本），AudioManager 重構（play_attack_by_character 統一走 play_sfx 路徑），Audio Sync 100/100**"

new_content = new_content.replace(old2, new2)

if new_content == content:
    print("WARNING: No changes made")
else:
    with open(path, "w", encoding="utf-8") as f:
        f.write(new_content)
    print("OK: progress.md updated")
