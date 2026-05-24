#!/usr/bin/env python3
"""
批量修復 GDScript 的所有常見問題：
1. 未閉合的字串（.text = "xxx 沒有結尾引號）
2. float64 型別
3. 重複的 const 宣告（檔案後半部分重複）
4. 兩個語句合併在同一行
"""
import os
import re
import subprocess

GODOT = r"C:\Users\yajinyee0306\AppData\Local\Microsoft\WinGet\Packages\GodotEngine.GodotEngine_Microsoft.Winget.Source_8wekyb3d8bbwe\Godot_v4.6.2-stable_win64_console.exe"
PROJECT_DIR = r"d:\Kiro\client\chiikawa-pixel"
SCRIPTS_DIR = os.path.join(PROJECT_DIR, "scripts", "ui")

def fix_unclosed_strings(lines):
    """修復未閉合的字串"""
    fixed = False
    for i, line in enumerate(lines):
        # 找 .text = "xxx 但沒有結尾引號的行
        if re.search(r'\.(text|tooltip_text|placeholder_text)\s*=\s*"[^"]*$', line.rstrip('\n')):
            lines[i] = line.rstrip('\n') + '"\n'
            fixed = True
            print(f"  Fixed unclosed string at line {i+1}")
    return fixed

def fix_float64(lines):
    """修復 float64 型別"""
    fixed = False
    for i, line in enumerate(lines):
        if 'float64' in line:
            lines[i] = line.replace('float64', 'float')
            fixed = True
    return fixed

def fix_duplicate_consts(lines):
    """修復重複的 const 宣告（截斷到第一個重複前）"""
    const_names = {}
    for i, line in enumerate(lines):
        stripped = line.strip()
        m = re.match(r'^const\s+(\w+)\s*[=:]', stripped)
        if m:
            name = m.group(1)
            if name in const_names:
                # 找到重複，截斷
                cut = i
                while cut > 0 and not lines[cut-1].strip():
                    cut -= 1
                print(f"  Truncated at line {cut+1} (duplicate const '{name}')")
                return lines[:cut], True
            else:
                const_names[name] = i
    return lines, False

def fix_merged_statements(lines):
    """修復兩個語句合併在同一行"""
    fixed = False
    new_lines = []
    for i, line in enumerate(lines):
        # 找 "= value   _var = " 模式
        m = re.match(r'^(\s*)(\w+\s*=\s*[^=\n]+?)\s{3,}(_\w+\s*=\s*.+)$', line.rstrip('\n'))
        if m:
            indent = m.group(1)
            stmt1 = m.group(2)
            stmt2 = m.group(3)
            new_lines.append(f"{indent}{stmt1}\n")
            new_lines.append(f"{indent}{stmt2}\n")
            fixed = True
            print(f"  Split merged statements at line {i+1}")
        else:
            new_lines.append(line)
    return new_lines, fixed

def process_file(path):
    """處理單個檔案"""
    fname = os.path.basename(path)
    
    with open(path, 'r', encoding='utf-8') as f:
        lines = f.readlines()
    
    original_len = len(lines)
    any_fixed = False
    
    # 1. 修復未閉合字串
    if fix_unclosed_strings(lines):
        any_fixed = True
    
    # 2. 修復 float64
    if fix_float64(lines):
        any_fixed = True
    
    # 3. 修復重複 const
    lines, fixed = fix_duplicate_consts(lines)
    if fixed:
        any_fixed = True
    
    # 4. 修復合併語句
    lines, fixed = fix_merged_statements(lines)
    if fixed:
        any_fixed = True
    
    if any_fixed:
        with open(path, 'w', encoding='utf-8') as f:
            f.writelines(lines)
        print(f"Fixed {fname} ({original_len} -> {len(lines)} lines)")
    
    return any_fixed

# 掃描所有腳本
scripts = [f for f in os.listdir(SCRIPTS_DIR) if f.endswith('.gd')]
total_fixed = 0

for script in sorted(scripts):
    path = os.path.join(SCRIPTS_DIR, script)
    print(f"\nProcessing {script}...")
    if process_file(path):
        total_fixed += 1

print(f"\n{'='*50}")
print(f"Total files fixed: {total_fixed}")
