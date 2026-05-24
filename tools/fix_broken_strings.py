#!/usr/bin/env python3
"""
修復 GDScript 中損壞的字串：
1. 找到包含損壞字元（?? 序列）的字串
2. 替換為安全的英文字串
3. 修復未閉合的括號
"""
import os
import re
import subprocess

GODOT = r"C:\Users\yajinyee0306\AppData\Local\Microsoft\WinGet\Packages\GodotEngine.GodotEngine_Microsoft.Winget.Source_8wekyb3d8bbwe\Godot_v4.6.2-stable_win64_console.exe"
PROJECT_DIR = r"d:\Kiro\client\chiikawa-pixel"
SCRIPTS_DIR = os.path.join(PROJECT_DIR, "scripts", "ui")

def check_script(script_name):
    """檢查腳本是否有 parse error"""
    rel_path = f"res://scripts/ui/{script_name}"
    try:
        result = subprocess.run(
            [GODOT, "--headless", "--check-only", "--script", rel_path],
            capture_output=True, text=True, cwd=PROJECT_DIR, timeout=8
        )
        output = result.stdout + result.stderr
        if "Parse Error" in output:
            m = re.search(r'\.gd:(\d+)', output)
            line = int(m.group(1)) if m else 0
            m2 = re.search(r'Parse Error: (.+)', output)
            msg = m2.group(1) if m2 else "unknown"
            return True, line, msg
    except subprocess.TimeoutExpired:
        pass
    return False, 0, ""

def fix_line_with_broken_chars(line: str) -> str:
    """修復包含損壞字元的行"""
    # 找到字串字面量中的損壞字元（? 後面跟著非 ASCII 字元）
    # 策略：把字串中的損壞部分替換為 "..."
    
    # 找到所有字串字面量
    result = []
    i = 0
    in_string = False
    string_start = -1
    
    while i < len(line):
        c = line[i]
        
        if not in_string and c == '"':
            in_string = True
            string_start = i
            result.append(c)
        elif in_string and c == '\\':
            result.append(c)
            if i + 1 < len(line):
                result.append(line[i+1])
                i += 2
            continue
        elif in_string and c == '"':
            in_string = False
            result.append(c)
        elif in_string:
            # 在字串中，檢查是否是損壞字元
            # 損壞字元通常是 ? 後面跟著奇怪的字元
            char_code = ord(c)
            if char_code > 127 and char_code < 160:
                # 損壞的控制字元，跳過
                pass
            else:
                result.append(c)
        else:
            result.append(c)
        
        i += 1
    
    return ''.join(result)

def fix_file_errors(path: str) -> bool:
    """修復單個檔案的所有錯誤"""
    fname = os.path.basename(path)
    
    # 先檢查是否有錯誤
    has_error, line_num, error_msg = check_script(fname)
    if not has_error:
        return False
    
    with open(path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed = False
    
    # 修復策略：根據錯誤類型
    if line_num > 0 and line_num <= len(lines):
        line_idx = line_num - 1
        original = lines[line_idx]
        
        if "closing" in error_msg or "arguments" in error_msg:
            # 括號不匹配，可能是字串中有括號
            # 嘗試修復：把損壞的 .get() 預設值替換為安全字串
            fixed_line = re.sub(
                r'\.get\("([^"]+)",\s*"[^"]*\)',
                r'.get("\1", "")',
                original
            )
            if fixed_line != original:
                lines[line_idx] = fixed_line
                fixed = True
                print(f"  Fixed {fname}:{line_num}: replaced broken default value")
        
        elif "Identifier" in error_msg or "statement" in error_msg:
            # 可能是字串中有損壞字元導致解析失敗
            fixed_line = fix_line_with_broken_chars(original)
            if fixed_line != original:
                lines[line_idx] = fixed_line
                fixed = True
                print(f"  Fixed {fname}:{line_num}: removed broken chars")
    
    if fixed:
        with open(path, 'w', encoding='utf-8') as f:
            f.writelines(lines)
    
    return fixed

# 掃描所有腳本
scripts = [f for f in os.listdir(SCRIPTS_DIR) if f.endswith('.gd')]
total_fixed = 0

for script in sorted(scripts):
    path = os.path.join(SCRIPTS_DIR, script)
    if fix_file_errors(path):
        total_fixed += 1

print(f"\nTotal files fixed: {total_fixed}")
