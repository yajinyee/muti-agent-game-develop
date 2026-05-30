#!/usr/bin/env python3
"""
DAY-337: HUD.gd 重構腳本
目標：移除 HUD.gd 中重複的 Lucky 訊號連接和 Lucky 函數
這些已全部移入 HUDLuckySignals.gd
"""
import re

HUD_PATH = r"d:\Kiro\client\chiikawa-pixel\scripts\ui\HUD.gd"

with open(HUD_PATH, "r", encoding="utf-8") as f:
    content = f.read()
    lines = content.split("\n")

print(f"原始行數: {len(lines)}")

# 找到需要移除的範圍
# 1. 移除行 63-91 的重複訊號連接（0-indexed: 62-90）
# 2. 移除行 487 開始的所有 Lucky 函數

# 找到 "# DAY-292~303" 的行（0-indexed）
start_signals = None
end_signals = None
start_lucky_funcs = None

for i, line in enumerate(lines):
    if "# DAY-292~303 幸運特殊魚訊號連接" in line and start_signals is None:
        start_signals = i
    if "# DAY-304~319 訊號已移入 HUDLuckySignals.gd" in line and start_signals is not None:
        end_signals = i + 1  # 包含這行
    if "func _on_lucky_chain_lightning" in line and start_lucky_funcs is None:
        start_lucky_funcs = i

print(f"訊號連接範圍: 行 {start_signals+1} - {end_signals+1} (0-indexed: {start_signals}-{end_signals})")
print(f"Lucky 函數開始: 行 {start_lucky_funcs+1} (0-indexed: {start_lucky_funcs})")

# 建立新的行列表
new_lines = []

# 保留 start_signals 之前的所有行
for i in range(start_signals):
    new_lines.append(lines[i])

# 跳過 start_signals 到 end_signals（重複的訊號連接）
# 在這裡加一個注釋說明
new_lines.append("\t# DAY-337：Lucky 訊號連接已全部移入 HUDLuckySignals.gd（由 _init_lucky_signals 初始化）")

# 保留 end_signals+1 到 start_lucky_funcs-1 的行（核心 HUD 函數）
for i in range(end_signals + 1, start_lucky_funcs):
    new_lines.append(lines[i])

# 跳過 start_lucky_funcs 到末尾（所有 Lucky 函數）
# 加一個注釋說明
new_lines.append("")
new_lines.append("# ── DAY-337 重構完成 ────────────────────────────────────────")
new_lines.append("# 所有 Lucky 函數已移入 HUDLuckySignals.gd")
new_lines.append("# HUD.gd 只保留核心 HUD 功能")

print(f"新行數: {len(new_lines)}")
print(f"減少行數: {len(lines) - len(new_lines)}")

# 寫入新文件
new_content = "\n".join(new_lines)
with open(HUD_PATH, "w", encoding="utf-8") as f:
    f.write(new_content)

print("✅ HUD.gd 重構完成")

# 驗證
with open(HUD_PATH, "r", encoding="utf-8") as f:
    verify_lines = f.read().split("\n")
print(f"驗證行數: {len(verify_lines)}")

# 確認 Lucky 函數已移除
lucky_funcs = [l for l in verify_lines if l.strip().startswith("func _on_lucky_")]
print(f"剩餘 Lucky 函數數量: {len(lucky_funcs)}")
if lucky_funcs:
    print("⚠️ 仍有 Lucky 函數:")
    for f in lucky_funcs[:5]:
        print(f"  {f}")
else:
    print("✅ 所有 Lucky 函數已移除")

# 確認核心函數仍在
core_funcs = ["func _ready", "func _process", "func _update_ui", "func _on_player_updated",
              "func _on_state_changed", "func _on_reward_received", "func _on_boss_event",
              "func _show_lucky_banner", "func _show_lucky_event"]
for func in core_funcs:
    found = any(func in l for l in verify_lines)
    status = "✅" if found else "❌"
    print(f"{status} {func}")
