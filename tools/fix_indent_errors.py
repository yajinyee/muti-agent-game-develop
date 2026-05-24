#!/usr/bin/env python3
"""
修復 GDScript 中常見的縮排問題：
1. 函數體第一行縮排多了一層（應該和前後行同層）
2. 型別推斷問題（var x := array[i] 改為 var x: Type = array[i]）
"""
import os
import re
import subprocess
import sys

GODOT = r"C:\Users\yajinyee0306\AppData\Local\Microsoft\WinGet\Packages\GodotEngine.GodotEngine_Microsoft.Winget.Source_8wekyb3d8bbwe\Godot_v4.6.2-stable_win64_console.exe"
PROJECT_DIR = r"d:\Kiro\client\chiikawa-pixel"
SCRIPTS_DIR = os.path.join(PROJECT_DIR, "scripts", "ui")

def check_script(script_name):
    """用 Godot 檢查腳本是否有 parse error，返回 (has_error, line_num, error_msg)"""
    rel_path = f"res://scripts/ui/{script_name}"
    result = subprocess.run(
        [GODOT, "--headless", "--check-only", "--script", rel_path],
        capture_output=True, text=True, cwd=PROJECT_DIR, timeout=10
    )
    output = result.stdout + result.stderr
    if "Parse Error" in output:
        m = re.search(r'\.gd:(\d+)', output)
        line = int(m.group(1)) if m else 0
        # 取錯誤訊息
        m2 = re.search(r'Parse Error: (.+)', output)
        msg = m2.group(1) if m2 else "unknown"
        return True, line, msg
    return False, 0, ""

def fix_indent_at_line(lines, line_idx):
    """嘗試修復第 line_idx 行的縮排問題"""
    if line_idx <= 0 or line_idx >= len(lines):
        return False
    
    current = lines[line_idx]
    current_tabs = len(current) - len(current.lstrip('\t'))
    
    # 找前一個非空行的縮排
    prev_tabs = 0
    for i in range(line_idx - 1, -1, -1):
        if lines[i].strip():
            prev_tabs = len(lines[i]) - len(lines[i].lstrip('\t'))
            break
    
    # 找後一個非空行的縮排
    next_tabs = 0
    for i in range(line_idx + 1, len(lines)):
        if lines[i].strip():
            next_tabs = len(lines[i]) - len(lines[i].lstrip('\t'))
            break
    
    # 如果當前縮排比前後都多，嘗試減少
    if current_tabs > prev_tabs and current_tabs > next_tabs:
        target_tabs = max(prev_tabs, next_tabs)
        lines[line_idx] = '\t' * target_tabs + current.lstrip('\t')
        return True
    
    # 如果當前縮排比前後都少，嘗試增加
    if current_tabs < prev_tabs and current_tabs < next_tabs:
        target_tabs = min(prev_tabs, next_tabs)
        lines[line_idx] = '\t' * target_tabs + current.lstrip('\t')
        return True
    
    return False

def fix_type_inference(lines, line_idx):
    """修復型別推斷問題：var x := expr 改為 var x: Variant = expr"""
    if line_idx <= 0 or line_idx >= len(lines):
        return False
    
    line = lines[line_idx]
    # 找 var x := 模式
    m = re.match(r'^(\s*)var\s+(\w+)\s*:=\s*(.+)$', line)
    if m:
        indent = m.group(1)
        varname = m.group(2)
        expr = m.group(3)
        lines[line_idx] = f"{indent}var {varname}: Variant = {expr}\n"
        return True
    return False

# 掃描所有腳本
scripts = [f for f in os.listdir(SCRIPTS_DIR) if f.endswith(".gd")]
total_fixed = 0

for script in sorted(scripts):
    path = os.path.join(SCRIPTS_DIR, script)
    
    # 最多嘗試 10 次修復
    for attempt in range(10):
        try:
            has_error, line_num, error_msg = check_script(script)
        except subprocess.TimeoutExpired:
            break
        
        if not has_error:
            break
        
        # 讀取檔案
        with open(path, "r", encoding="utf-8") as f:
            lines = f.readlines()
        
        line_idx = line_num - 1  # 0-indexed
        fixed = False
        
        if "Unindent" in error_msg or "Indent" in error_msg:
            fixed = fix_indent_at_line(lines, line_idx)
        elif "infer the type" in error_msg or "Variant value" in error_msg:
            fixed = fix_type_inference(lines, line_idx)
        
        if fixed:
            with open(path, "w", encoding="utf-8") as f:
                f.writelines(lines)
            print(f"Fixed {script} line {line_num}: {error_msg[:50]}")
            total_fixed += 1
        else:
            print(f"CANNOT FIX {script} line {line_num}: {error_msg[:60]}")
            break

print(f"\nTotal fixes applied: {total_fixed}")
