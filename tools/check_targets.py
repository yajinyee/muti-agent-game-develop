# -*- coding: utf-8 -*-
from PIL import Image
import os

d = "D:/Kiro/client/chiikawa-pixel/assets/sprites/targets"
print("=== 目標物品質報告 ===\n")
total_old = 0
total_new = 0
for f in sorted(os.listdir(d)):
    if not f.endswith(".png") or "B001" in f:
        continue
    img = Image.open(os.path.join(d, f)).convert("RGBA")
    non_t = sum(1 for px in img.getdata() if px[3] > 10)
    pct = non_t * 100 // (img.width * img.height)
    status = "✅" if non_t > 1500 else ("⚠️" if non_t > 500 else "❌")
    print(f"  {status} {f}: {img.size}, {non_t}px ({pct}%)")
    total_new += non_t

print(f"\n  平均非透明像素: {total_new // 11}px")
