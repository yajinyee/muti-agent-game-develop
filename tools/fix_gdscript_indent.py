#!/usr/bin/env python3
"""
自動修復 GDScript 縮排問題。
GDScript 使用 tab 縮排，規則：
- func/class/if/for/while/match 等關鍵字後的區塊需要縮排
- 縮排層級由上下文決定
"""
import os
import re
import sys

def fix_gdscript_indent(content: str) -> str:
    """修復 GDScript 的縮排"""
    lines = content.split('\n')
    result = []
    indent_level = 0
    
    # 需要增加縮排的關鍵字（這些行的下一行需要縮排）
    INDENT_INCREASE = re.compile(
        r'^\s*(func\s+\w+|class\s+\w+|if\s+|elif\s+|else\s*:|for\s+|while\s+|match\s+|'
        r'_ready\s*\(\)|_process\s*\(|_input\s*\(|_physics_process\s*\()'
    )
    # 這些行本身就是縮排塊的開始（以冒號結尾）
    BLOCK_START = re.compile(r':\s*$')
    
    i = 0
    while i < len(lines):
        line = lines[i]
        stripped = line.strip()
        
        # 空行保持
        if not stripped:
            result.append('')
            i += 1
            continue
        
        # 計算原始縮排
        orig_tabs = len(line) - len(line.lstrip('\t'))
        orig_spaces = len(line) - len(line.lstrip(' '))
        
        # 如果已經有 tab 縮排，保持原樣
        if orig_tabs > 0:
            result.append(line)
            i += 1
            continue
        
        # 如果有空格縮排，轉換為 tab
        if orig_spaces > 0:
            tab_count = orig_spaces // 4 if orig_spaces >= 4 else 1
            result.append('\t' * tab_count + stripped)
            i += 1
            continue
        
        # 沒有縮排的行
        result.append(stripped)
        i += 1
    
    return '\n'.join(result)


def smart_fix_indent(content: str) -> str:
    """
    智能修復縮排：
    1. 找到所有函數定義
    2. 確認函數體有縮排
    3. 如果沒有，加上縮排
    """
    lines = content.split('\n')
    result = list(lines)
    
    i = 0
    while i < len(lines):
        line = lines[i]
        stripped = line.strip()
        
        # 找函數/if/for/while/match 定義（以冒號結尾）
        is_block_start = (
            re.match(r'^(func|class|if|elif|else|for|while|match)\b', stripped) or
            re.match(r'^(var\s+\w+\s*=\s*func)', stripped)
        ) and stripped.endswith(':')
        
        if is_block_start:
            # 計算當前行的縮排
            current_tabs = len(line) - len(line.lstrip('\t'))
            expected_body_tabs = current_tabs + 1
            
            # 檢查下一個非空行是否有正確縮排
            j = i + 1
            while j < len(lines) and not lines[j].strip():
                j += 1
            
            if j < len(lines):
                next_line = lines[j]
                next_stripped = next_line.strip()
                next_tabs = len(next_line) - len(next_line.lstrip('\t'))
                
                # 如果下一行沒有縮排（且不是另一個頂層定義）
                if (next_tabs == 0 and next_stripped and 
                    not re.match(r'^(func|class|var|const|signal|@|#|extends|class_name)', next_stripped)):
                    
                    # 找到這個區塊的結束（下一個同層或更低層的行）
                    # 簡單策略：找到下一個 func/class 定義或空行後的頂層行
                    block_end = j
                    while block_end < len(lines):
                        bl = lines[block_end]
                        bl_stripped = bl.strip()
                        if not bl_stripped:
                            block_end += 1
                            continue
                        bl_tabs = len(bl) - len(bl.lstrip('\t'))
                        # 如果是頂層定義，停止
                        if bl_tabs == 0 and re.match(r'^(func|class|var|const|signal|@)', bl_stripped):
                            break
                        block_end += 1
                    
                    # 給 j 到 block_end 之間的行加縮排
                    for k in range(j, block_end):
                        if lines[k].strip():
                            result[k] = '\t' * expected_body_tabs + lines[k].lstrip('\t')
        
        i += 1
    
    return '\n'.join(result)


def process_file(path: str) -> bool:
    """處理單個檔案，返回是否有修改"""
    try:
        with open(path, 'r', encoding='utf-8') as f:
            original = f.read()
    except Exception as e:
        print(f"  READ ERROR: {e}")
        return False
    
    # 檢查是否有縮排問題（函數體沒有縮排）
    lines = original.split('\n')
    has_indent_issue = False
    
    for i, line in enumerate(lines):
        stripped = line.strip()
        if (re.match(r'^(func|if|elif|else|for|while|match)\b', stripped) and 
            stripped.endswith(':') and i + 1 < len(lines)):
            next_line = lines[i + 1]
            next_stripped = next_line.strip()
            if (next_stripped and 
                not next_line.startswith('\t') and
                not re.match(r'^(func|class|var|const|signal|@|#|extends|class_name)', next_stripped)):
                has_indent_issue = True
                break
    
    if not has_indent_issue:
        return False
    
    # 應用修復
    fixed = smart_fix_indent(original)
    
    if fixed != original:
        with open(path, 'w', encoding='utf-8') as f:
            f.write(fixed)
        return True
    
    return False


if __name__ == '__main__':
    scripts_dir = r"d:\Kiro\client\chiikawa-pixel\scripts"
    
    total = 0
    fixed = 0
    
    for root, dirs, files in os.walk(scripts_dir):
        for f in files:
            if not f.endswith('.gd'):
                continue
            path = os.path.join(root, f)
            total += 1
            if process_file(path):
                fixed += 1
                print(f"Fixed: {f}")
    
    print(f"\nTotal: {total}, Fixed: {fixed}")
