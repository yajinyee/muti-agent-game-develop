#!/usr/bin/env python3
"""掃描 GDScript 檔案的常見語法問題"""
import os
import re

scripts_dir = r"d:\Kiro\client\chiikawa-pixel\scripts"
errors = []

for root, dirs, files in os.walk(scripts_dir):
    for f in files:
        if not f.endswith(".gd"):
            continue
        path = os.path.join(root, f)
        try:
            with open(path, "r", encoding="utf-8") as fh:
                lines = fh.readlines()
        except Exception as e:
            errors.append(f"READ ERROR {f}: {e}")
            continue
        
        for i, line in enumerate(lines):
            stripped = line.rstrip("\n\r")
            
            # 1. 未閉合的字串（.text = " 但沒有結尾引號）
            if re.search(r'\.text\s*=\s*"[^"]*$', stripped):
                errors.append(f"{f}:{i+1}: unclosed string: {stripped[:60]}")
            
            # 2. float64 型別（GDScript 不支援）
            if "float64" in stripped:
                errors.append(f"{f}:{i+1}: invalid type float64: {stripped[:60]}")
            
            # 3. 函數體沒有縮排
            if re.match(r'^func\s+\w+', stripped):
                if i + 1 < len(lines):
                    next_line = lines[i+1].rstrip("\n\r")
                    if (next_line and 
                        not next_line.startswith("\t") and 
                        not next_line.startswith("#") and 
                        not next_line.startswith("func ") and
                        not next_line.startswith("var ") and
                        not next_line.startswith("const ") and
                        not next_line.startswith("@") and
                        not next_line.startswith("signal ") and
                        len(next_line.strip()) > 0):
                        errors.append(f"{f}:{i+2}: function body not indented: {next_line[:50]}")

print(f"Total issues: {len(errors)}")
for e in sorted(set(errors))[:50]:
    print(e)
