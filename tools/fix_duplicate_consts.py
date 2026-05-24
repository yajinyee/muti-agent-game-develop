#!/usr/bin/env python3
"""找出並修復 GDScript 中重複的 const/var 宣告（通常是檔案後半部分重複）"""
import os
import re

scripts_dir = r"d:\Kiro\client\chiikawa-pixel\scripts\ui"
fixed_files = []

for f in sorted(os.listdir(scripts_dir)):
    if not f.endswith(".gd"):
        continue
    path = os.path.join(scripts_dir, f)
    try:
        with open(path, "r", encoding="utf-8") as fh:
            lines = fh.readlines()
    except Exception as e:
        print(f"READ ERROR {f}: {e}")
        continue
    
    # 找重複的 const 宣告
    const_names = {}
    first_dup_line = None
    
    for i, line in enumerate(lines):
        stripped = line.strip()
        m = re.match(r'^const\s+(\w+)\s*[=:]', stripped)
        if m:
            name = m.group(1)
            if name in const_names:
                # 找到重複！記錄第一個重複的行號
                if first_dup_line is None or i < first_dup_line:
                    first_dup_line = i
                    print(f"{f}: duplicate const '{name}' at line {i+1} (first at {const_names[name]+1})")
            else:
                const_names[name] = i
    
    if first_dup_line is not None:
        # 截斷到重複開始前的位置
        # 找到重複前的最後一個非空行
        cut_line = first_dup_line
        while cut_line > 0 and lines[cut_line-1].strip() == "":
            cut_line -= 1
        
        trimmed = lines[:cut_line]
        with open(path, "w", encoding="utf-8") as fh:
            fh.writelines(trimmed)
        
        print(f"  -> Trimmed {f} from {len(lines)} to {len(trimmed)} lines")
        fixed_files.append(f)

print(f"\nFixed {len(fixed_files)} files:")
for f in fixed_files:
    print(f"  {f}")
