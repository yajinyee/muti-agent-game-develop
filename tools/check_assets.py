"""
check_assets.py
快速檢查所有美術資產的尺寸和非透明像素比例
"""
from PIL import Image
import os

def check_dir(label, path, ext=".png"):
    print(f"\n=== {label} ===")
    if not os.path.exists(path):
        print(f"  [MISSING] {path}")
        return
    files = sorted([f for f in os.listdir(path) if f.endswith(ext)])
    for f in files:
        fp = os.path.join(path, f)
        img = Image.open(fp).convert("RGBA")
        w, h = img.size
        pixels = list(img.getdata())
        non_transparent = sum(1 for p in pixels if p[3] > 10)
        total = w * h
        pct = non_transparent / total * 100
        flag = "✅" if pct > 30 else "⚠️"
        print(f"  {flag} {f}: {w}x{h}, {pct:.0f}% non-transparent ({non_transparent}/{total})")

BASE = "d:/Kiro/client/chiikawa-pixel/assets/sprites"

check_dir("Characters", f"{BASE}/characters")
check_dir("Targets", f"{BASE}/targets")
check_dir("Effects", f"{BASE}/effects")
check_dir("Backgrounds", f"{BASE}/backgrounds")
check_dir("UI", f"{BASE}/ui")

print("\n=== Sheets ===")
sheets_dir = f"{BASE}/sheets"
for f in sorted(os.listdir(sheets_dir)):
    if f.endswith(".png"):
        fp = os.path.join(sheets_dir, f)
        img = Image.open(fp)
        print(f"  {f}: {img.size}")

print("\n✅ 檢查完成")
