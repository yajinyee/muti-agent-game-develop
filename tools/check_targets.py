# -*- coding: utf-8 -*-
from PIL import Image
import os

TARGETS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"

targets = [f for f in os.listdir(TARGETS_DIR) if f.endswith(".png")]
targets.sort()

print("=== 目標物品質分析 ===")
for name in targets:
    path = os.path.join(TARGETS_DIR, name)
    img = Image.open(path).convert("RGBA")
    w, h = img.size
    total = w * h
    non_transparent = sum(1 for p in img.getdata() if p[3] > 10)
    pct = non_transparent / total * 100
    status = "✅" if pct >= 40 else "⚠️ "
    print(f"  {status} {name}: {w}x{h}, {non_transparent}/{total} ({pct:.0f}%)")
