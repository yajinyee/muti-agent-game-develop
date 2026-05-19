#!/usr/bin/env python3
"""修復 HUD.gd 中的 _process 函數，移除舊的 _process_session_stats 呼叫"""
import re

path = r"d:\Kiro\client\chiikawa-pixel\scripts\ui\HUD.gd"

with open(path, "r", encoding="utf-8") as f:
    content = f.read()

# 找到 _process 函數中的舊呼叫並替換
old_pattern = r'func _process\(delta: float\) -> void:\n\t# Session Stats.*?\n\t_process_session_stats\(delta\)\n'
new_text = 'func _process(delta: float) -> void:\n\t# Session Stats 自動彈出由 SessionStatsPanel._process 處理（DAY-053）\n'

new_content = re.sub(old_pattern, new_text, content, flags=re.DOTALL)

if new_content == content:
    print("WARNING: Pattern not found, trying line-by-line replacement")
    lines = content.split('\n')
    result = []
    i = 0
    while i < len(lines):
        line = lines[i]
        if '_process_session_stats(delta)' in line:
            print(f"Found at line {i+1}: {line!r}")
            # Skip this line
            i += 1
            continue
        # Also replace the comment line before it
        if i + 1 < len(lines) and '_process_session_stats(delta)' in lines[i+1]:
            if '# Session Stats' in line:
                result.append('\t# Session Stats 自動彈出由 SessionStatsPanel._process 處理（DAY-053）')
                i += 2  # skip both comment and call
                continue
        result.append(line)
        i += 1
    new_content = '\n'.join(result)
else:
    print("Pattern found and replaced via regex")

with open(path, "w", encoding="utf-8") as f:
    f.write(new_content)

# Verify
with open(path, "r", encoding="utf-8") as f:
    verify = f.read()

if '_process_session_stats' in verify:
    print("ERROR: Still found _process_session_stats in file")
else:
    print("OK: _process_session_stats removed successfully")

lines = verify.split('\n')
print(f"Total lines: {len(lines)}")
