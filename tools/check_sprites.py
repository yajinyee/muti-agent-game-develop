from PIL import Image
import os

base = r"d:\Kiro\client\chiikawa-pixel\assets\sprites"

targets = ["targets/T001_grass.png", "targets/T002_bug_g.png", "targets/T005_pudding.png",
           "targets/T101_mimic.png", "targets/T105_coin_fish.png", "targets/B001_boss.png"]
chars = ["characters/chiikawa_idle.png", "characters/hachiware_idle.png", "characters/usagi_idle.png"]

print("=== Target Sprites ===")
for t in targets:
    path = os.path.join(base, t)
    if os.path.exists(path):
        img = Image.open(path)
        print(f"{t}: {img.size}, mode={img.mode}")
    else:
        print(f"{t}: NOT FOUND")

print("\n=== Character Sprites ===")
for c in chars:
    path = os.path.join(base, c)
    if os.path.exists(path):
        img = Image.open(path)
        print(f"{c}: {img.size}, mode={img.mode}")
    else:
        print(f"{c}: NOT FOUND")
