#!/usr/bin/env python3
"""
修復 GDScript 中所有未閉合的字串字面量。
策略：找到行尾沒有閉合引號的字串，加上閉合引號。
"""
import os
import re

SCRIPTS_DIR = r"d:\Kiro\client\chiikawa-pixel\scripts"

def has_unclosed_string(line: str) -> bool:
    """檢查一行是否有未閉合的字串"""
    # 移除行尾換行
    line = line.rstrip('\n\r')
    
    # 計算引號數量（簡單版：不處理轉義）
    # 找到所有 " 的位置
    in_string = False
    i = 0
    while i < len(line):
        c = line[i]
        if c == '\\' and in_string:
            i += 2  # 跳過轉義字元
            continue
        if c == '"':
            in_string = not in_string
        i += 1
    
    return in_string  # 如果還在字串中，表示未閉合

def fix_file(path: str) -> int:
    """修復單個檔案，返回修復的行數"""
    fname = os.path.basename(path)
    
    with open(path, 'r', encoding='utf-8') as f:
        lines = f.readlines()
    
    fixed_count = 0
    for i, line in enumerate(lines):
        stripped = line.rstrip('\n\r')
        if has_unclosed_string(stripped):
            # 加上閉合引號
            lines[i] = stripped + '"\n'
            fixed_count += 1
            print(f"  {fname}:{i+1}: fixed unclosed string")
    
    if fixed_count > 0:
        with open(path, 'w', encoding='utf-8') as f:
            f.writelines(lines)
    
    return fixed_count

# 掃描所有腳本
total_fixed = 0
for root, dirs, files in os.walk(SCRIPTS_DIR):
    for f in sorted(files):
        if not f.endswith('.gd'):
            continue
        path = os.path.join(root, f)
        count = fix_file(path)
        total_fixed += count

print(f"\nTotal lines fixed: {total_fixed}")
