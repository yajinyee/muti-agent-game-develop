#!/usr/bin/env python3
"""掃描 UI GDScript 腳本的常見問題"""
import os
import re

scripts_dir = r"d:\Kiro\client\chiikawa-pixel\scripts\ui"
all_errors = {}

for f in sorted(os.listdir(scripts_dir)):
    if not f.endswith(".gd"):
        continue
    path = os.path.join(scripts_dir, f)
    try:
        with open(path, "r", encoding="utf-8") as fh:
            lines = fh.readlines()
    except Exception as e:
        all_errors[f] = [f"READ ERROR: {e}"]
        continue
    
    issues = []
    for i, line in enumerate(lines):
        stripped = line.rstrip("\n\r")
        # 找兩個語句合併在同一行（用多個空格分隔）
        if re.search(r'\)\s{3,}_\w+\s*=', stripped):
            issues.append(f"Line {i+1}: merged statements")
        # 找未閉合字串
        if re.search(r'\.text\s*=\s*"[^"]*$', stripped):
            issues.append(f"Line {i+1}: unclosed string")
        # 找重複的 const 宣告（簡單版）
    
    if issues:
        all_errors[f] = issues

print(f"Files with issues: {len(all_errors)}")
for fname, issues in all_errors.items():
    print(f"\n{fname}:")
    for issue in issues:
        print(f"  {issue}")
